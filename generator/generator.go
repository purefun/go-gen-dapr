package generator

import (
	"errors"
	"fmt"
	"go/format"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/purefun/go-gen-dapr/generator/box"
	"golang.org/x/tools/go/packages"
)

const Version = "v0.7.3"

var (
	ErrServiceNotFound     = errors.New("service not found")
	ErrEmptyService        = errors.New("empty service")
	ErrNotInterface        = errors.New("service is not an interface")
	ErrNoCtxParam          = errors.New(`first param is not "context.Context"`)
	ErrInvalidResults      = errors.New(`method results are not "error", neither "(*SomeType, error)"`)
	ErrInvalidGenerateType = errors.New("invalid generate type")
)

type GenerateType string

const (
	GenerateTypeService    GenerateType = "service"
	GenerateTypeSubscriber GenerateType = "subscriber"
)

func GenerateTypeFromString(t string) (GenerateType, error) {
	gt := GenerateType(t)
	switch gt {
	case GenerateTypeService, GenerateTypeSubscriber:
		return gt, nil
	default:
		return GenerateType(""), fmt.Errorf("invalid generate type: %s", t)
	}
}

type Options struct {
	ServicePkg   string
	ServiceType  string
	GenComment   bool
	GenerateType GenerateType
}

type Id struct {
	Pkg  string
	Name string
}

func NewGenerator(o Options) *Generator {
	return &Generator{
		Version:      Version,
		GenerateType: o.GenerateType,
		ServicePkg:   o.ServicePkg,
		ServiceType:  o.ServiceType,
		GenComment:   o.GenComment,
		Imports:      make(map[string]string),
	}
}

type Param struct {
	Name string
	Type string
}

type Response struct {
	Name    string
	Type    string
	IsSlice bool
}

type Method struct {
	Name      string
	Params    []*Param
	Responses []*Response
	Response  *Response
}

type Generator struct {
	Version    string
	SourceType string
	Package    *packages.Package

	GenerateType GenerateType

	ServicePkg  string
	ServiceType string

	PackageName string // package {{.PackageName}}
	GenComment  bool
	Imports     map[string]string // package->alias
	Methods     []*Method
}

func (g *Generator) Generate() (string, error) {
	pkg, err := LoadSource(g.ServicePkg)
	if err != nil {
		return "", nil
	}
	g.Package = pkg
	g.PackageName = pkg.Name

	err = g.Build()
	if err != nil {
		return "", err
	}

	pkgOut, err := g.genPackage()
	if err != nil {
		return "", err
	}

	var mainOut string

	switch g.GenerateType {
	case GenerateTypeService:
		mainOut, err = g.genService()
	case GenerateTypeSubscriber:
		mainOut, err = g.genSubscriber()
	default:
		err = ErrInvalidGenerateType
	}

	if err != nil {
		return "", err
	}

	// should be the last genXXX
	importsOut, err := g.genImports()
	if err != nil {
		return "", err
	}

	out := pkgOut + importsOut + mainOut

	formatted, err := format.Source([]byte(out))
	if err != nil {
		return "", fmt.Errorf("format source failed: %w, source: \n%s", err, out)
	}

	return string(formatted), nil
}

func (g *Generator) Build() error {
	var s types.Object

	s = g.Package.Types.Scope().Lookup(g.ServiceType)

	if s == nil {
		return fmt.Errorf("%w, name: %s", ErrServiceNotFound, g.ServiceType)
	}

	iface, ok := s.Type().Underlying().(*types.Interface)
	if !ok {
		return fmt.Errorf("%w, name: %s", ErrNotInterface, g.PackageName)
	}

	if iface.Empty() {
		return fmt.Errorf("%w, name: %s", ErrEmptyService, g.PackageName)
	}

	g.SourceType = s.Type().String()

	for i := 0; i < iface.NumMethods(); i++ {
		err := g.BuildMethod(iface.Method(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) BuildMethod(m *types.Func) error {
	name := m.Id()

	if !m.Exported() {
		fmt.Println("WARN: unexported method will be ignored, method:", name)
		return nil
	}

	method := &Method{Name: name}
	sig, _ := m.Type().(*types.Signature)

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		method.Params = append(method.Params, &Param{
			Name: p.Name(),
			Type: g.typeName(p.Type()),
		})
	}

	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		typeName := g.typeName(result.Type())
		method.Responses = append(method.Responses, &Response{
			Name:    result.Name(),
			Type:    typeName,
			IsSlice: strings.HasPrefix(typeName, "[]"),
		})
	}

	g.Methods = append(g.Methods, method)
	return nil
}

func (g *Generator) typeName(t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		pkg := p.Path()
		if pkg != g.Package.PkgPath {
			g.AddImport(p.Path(), "")
			return filepath.Base(pkg)
		}

		return ""
	})
}

func (g *Generator) validateServiceParams() error {
	for _, method := range g.Methods {
		if len(method.Params) == 0 || method.Params[0].Type != "context.Context" {
			return fmt.Errorf("%w, method: %s", ErrNoCtxParam, method.Name)
		}
	}
	return nil
}

func (g *Generator) validateServiceResponses() error {
	for _, method := range g.Methods {
		err := fmt.Errorf("%w, method: %s", ErrInvalidResults, method.Name)
		l := len(method.Responses)
		if l == 0 || l > 2 {
			return err
		}
		if l == 1 {
			if method.Responses[0].Type != "error" {
				return err
			}
		}
		if l == 2 {
			respType := method.Responses[0].Type

			if !strings.HasPrefix(respType, "*") &&
				!strings.HasPrefix(respType, "[]") ||
				method.Responses[1].Type != "error" {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) AddImport(pkg, alias string) {
	g.Imports[pkg] = alias
}

func (g *Generator) genPackage() (string, error) {
	return GenPackage(GenPackageData{
		PackageName: g.PackageName,
		GenComment:  g.GenComment,
		Version:     g.Version,
		SourceType:  g.SourceType,
	})
}

func (g *Generator) genImports() (string, error) {
	return GenImports(GenImportsData{Imports: g.Imports})
}

func (g *Generator) genService() (string, error) {
	if err := g.validateServiceParams(); err != nil {
		return "", err
	}
	if err := g.validateServiceResponses(); err != nil {
		return "", err
	}

	needJSON := false

	// shift ctx param, response
	for _, m := range g.Methods {
		m.Params = m.Params[1:]
		if len(m.Responses) == 2 {
			m.Response = m.Responses[0]
		}
		if len(m.Params) > 0 || m.Response != nil {
			needJSON = true
		}
	}

	g.AddImport("context", "")
	if needJSON {
		g.AddImport("encoding/json", "")
	}
	g.AddImport("github.com/dapr/go-sdk/client", "\n")
	g.AddImport("github.com/dapr/go-sdk/service/common", "")
	g.AddImport("github.com/dapr/go-sdk/service/grpc", "")
	g.AddImport("github.com/purefun/go-gen-dapr/pkg/errors", "errorHandlers")
	g.AddImport("github.com/purefun/go-gen-dapr/pkg/dapr", "")
	g.AddImport("github.com/pkg/errors", "")

	out, err := box.Template.Execute("service.tmpl", g)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type Subscriber struct {
	EventName   string
	HandlerName string
}

func (g *Generator) genSubscriber() (string, error) {
	g.AddImport("context", "")
	g.AddImport("encoding/json", "")
	g.AddImport("github.com/purefun/go-gen-dapr/pkg/dapr", "")
	g.AddImport("github.com/dapr/go-sdk/service/common", "")
	g.AddImport("github.com/purefun/go-gen-dapr/pkg/errors", "errorHandlers")

	var subs []Subscriber
	for _, m := range g.Methods {
		subs = append(subs, Subscriber{
			HandlerName: m.Name,
			EventName:   m.Params[1].Type,
		})
	}

	type Data struct {
		Handlers    []Subscriber
		ServiceType string
	}

	data := Data{
		Handlers:    subs,
		ServiceType: g.ServiceType,
	}

	return box.Template.Execute("subscriber.tmpl", data)
}
