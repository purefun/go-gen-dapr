package generator

import (
	"errors"
	"fmt"
	"go/format"
	"go/types"

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
	ErrNoCtxParam      = errors.New("first param of method is not context.Context")
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

type Result struct {
	Name string
	Type string
}

type Method struct {
	Name    string
	Params  []*Param
	Results []*Result
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
		m := iface.Method(i)
		methodName := m.Id()

		g.Methods = append(g.Methods, &Method{Name: methodName})

		sig, _ := m.Type().(*types.Signature)

		if !m.Exported() {
			fmt.Println("WARN: unexported method will be ignored, method:", methodName)
			continue
		}

		params := sig.Params()
		if params.Len() == 0 {
			return fmt.Errorf("%w, method: %s", ErrNoCtxParam, methodName)
		}
		if p := params.At(0); p.Type().String() != "context.Context" {
			return fmt.Errorf("%w, method: %s", ErrNoCtxParam, methodName)
		}
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)
			fmt.Println("param =>", p.Id(), p.Name(), p.Type().String())
		}

		// results := sig.Results()
	}

	return nil
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
	g.AddImport("github.com/dapr/go-sdk/client", "")
	g.AddImport("github.com/dapr/go-sdk/service/common", "")
	g.AddImport("github.com/dapr/go-sdk/service/grpc", "")

	out, err := box.Template.Execute("service.tmpl", g)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
