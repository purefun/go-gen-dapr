package generator

import (
	"github.com/purefun/go-gen-dapr/generator/box"
	"golang.org/x/tools/go/packages"
)

const LoadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

func LoadSource(pkg string) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: LoadMode}, pkg)
	if err != nil {
		return nil, err
	}
	return pkgs[0], nil
}

type GenPackageData struct {
	PackageName string
	GenComment  bool
	Version     string
	SourceType  string
}

func GenPackage(data GenPackageData) (string, error) {
	return box.Template.Execute("package.tmpl", data)
}

type ImportPkg = string
type ImportAlias = string

type GenImportsData struct {
	Imports map[ImportPkg]ImportAlias
}

func GenImports(data GenImportsData) (string, error) {
	return box.Template.Execute("imports.tmpl", data)
}
