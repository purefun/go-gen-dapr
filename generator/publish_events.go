package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"strings"

	"github.com/purefun/go-gen-dapr/generator/box"
)

var (
	ErrEventNotFound    = errors.New("event not found")
	ErrEventDocNotFound = errors.New("event doc not found")
)

type Event struct {
	Name  string
	Topic string
}

type PublishEventsOptions struct {
	Pkg        string
	GenComment bool
}

func GeneratePublishEvents(o PublishEventsOptions) (string, error) {
	pkg, err := LoadSource(o.Pkg)
	if err != nil {
		return "", err
	}
	pkgOut, err := GenPackage(GenPackageData{
		PackageName: pkg.Name,
		GenComment:  o.GenComment,
		Version:     Version,
		SourceType:  pkg.PkgPath,
	})
	if err != nil {
		return "", err
	}

	importsOut, err := GenImports(GenImportsData{
		Imports: map[string]string{
			"context":                       "",
			"encoding/json":                 "",
			"github.com/dapr/go-sdk/client": "",
		},
	})
	if err != nil {
		return "", nil
	}

	var events []Event

	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						name := typeSpec.Name.String()
						if strings.HasSuffix(name, "Event") {
							if typeSpec.Doc == nil && len(genDecl.Specs) > 1 {
								return "", fmt.Errorf("%w, event: %s", ErrEventDocNotFound, name)
							}
							events = append(events, Event{Name: name, Topic: name})
						}
					}
				}
			}
		}
	}

	if len(events) == 0 {
		return "", ErrEventNotFound
	}

	data := struct{ Events []Event }{
		Events: events,
	}

	pubsubOut, err := box.Template.Execute("pubsub.tmpl", data)

	if err != nil {
		return "", err
	}

	out := pkgOut + importsOut + pubsubOut
	formatted, err := format.Source([]byte(out))

	if err != nil {
		return "", fmt.Errorf("format source failed: %w, source: \n%s", err, out)
	}
	return string(formatted), nil
}
