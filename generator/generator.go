package generator

import (
	"errors"
	"fmt"
	"go/format"
	"go/types"
	"strings"

	"github.com/purefun/go-gen-dapr/generator/box"
	"golang.org/x/tools/go/packages"
)

const Version = "v0.1.0"

var LoadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

var (
	ErrServiceNotFound = errors.New("service not found")
	ErrEmptyService    = errors.New("empty service")
	ErrNotInterface    = errors.New("service is not an interface")
	ErrNoCtxParam      = errors.New(`first param is not "context.Context"`)
	ErrInvalidResults  = errors.New(`method results are not "error", neither "(*SomeType, error)"`)
)

type Options struct {
	PackageName string
	ServiceName string
	GenComment  bool
}

func NewGenerator(o Options) *Generator {
	return &Generator{
		Version:     Version,
		PackageName: o.PackageName,
		ServiceName: o.ServiceName,
		Imports:     make(map[string]string),
	}
}

type Param struct {
	Name string
	Type string
}

type Response struct {
	Type string
}

type Method struct {
	Name     string
	Params   []*Param
	Response *Response
}

type Generator struct {
	Version     string
	SourceType  string
	Packages    []*packages.Package
	ServiceName string
	PackageName string
	GenComment  bool
	Imports     map[string]string // package->alias
	Methods     []*Method
}

func (g *Generator) Generate() (string, error) {
	err := g.Build()
	if err != nil {
		return "", err
	}

	pkgOut, err := g.genPackage()
	if err != nil {
		return "", err
	}

	serviceOut, err := g.genService()
	if err != nil {
		return "", err
	}

	// should be the last genXXX
	importsOut, err := g.genImports()
	if err != nil {
		return "", err
	}

	out := pkgOut + importsOut + serviceOut

	formatted, err := format.Source([]byte(out))
	if err != nil {
		return "", fmt.Errorf("format source failed: %w, source: \n%s", err, out)
	}

	return string(formatted), nil
}

func (g *Generator) Load(cfg *packages.Config, patterns string) error {
	pkgs, err := packages.Load(cfg, patterns)
	if err != nil {
		return err
	}
	g.Packages = pkgs
	return nil
}

func (g *Generator) Build() error {
	var s types.Object

	for _, pkg := range g.Packages {
		s = pkg.Types.Scope().Lookup(g.ServiceName)
	}

	if s == nil {
		return fmt.Errorf("%w, name: %s", ErrServiceNotFound, g.ServiceName)
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
	if !g.validateParams(params) {
		return fmt.Errorf("%w, method: %s", ErrNoCtxParam, name)
	}
	// skip the first param: ctx context.Context
	for i := 1; i < params.Len(); i++ {
		p := params.At(i)
		method.Params = append(method.Params, &Param{
			Name: p.Name(),
			Type: p.Type().String(),
		})
	}

	results := sig.Results()
	if !g.validateResults(results) {
		return fmt.Errorf("%w, name: %s", ErrInvalidResults, name)
	}
	if results.Len() == 2 {
		method.Response = &Response{Type: results.At(0).Type().String()}
	}

	g.Methods = append(g.Methods, method)
	return nil
}

func (g *Generator) validateParams(ps *types.Tuple) bool {
	return ps.Len() > 0 && ps.At(0).Type().String() == "context.Context"
}

func (g *Generator) validateResults(rs *types.Tuple) bool {
	l := rs.Len()
	if l == 1 && rs.At(0).Type().String() == "error" {
		return true
	}
	if l == 2 &&
		strings.HasPrefix(rs.At(0).Type().String(), "*") &&
		rs.At(1).Type().String() == "error" {
		return true
	}
	return false
}

func (g *Generator) AddImport(pkg, alias string) {
	g.Imports[pkg] = alias
}

func (g *Generator) genPackage() (string, error) {
	out, err := box.Template.Execute("package.tmpl", g)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (g *Generator) genImports() (string, error) {
	out, err := box.Template.Execute("imports.tmpl", g)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (g *Generator) genService() (string, error) {
	g.AddImport("context", "")
	g.AddImport("encoding/json", "")
	g.AddImport("fmt", "")
	g.AddImport("github.com/dapr/go-sdk/client", "")
	g.AddImport("github.com/dapr/go-sdk/service/common", "")
	g.AddImport("github.com/dapr/go-sdk/service/grpc", "")

	out, err := box.Template.Execute("service.tmpl", g)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
