package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"regexp"
	"strings"

	"github.com/purefun/go-gen-dapr/generator/box"
)

var (
	ErrEventDocNotFound = errors.New("event doc not found")
	ErrTopicNotFound    = errors.New("event topic not found")
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

	// for _, name := range pkg.Types.Scope().Names() {
	// if strings.HasSuffix(name, "Event") {
	// 	eventType := pkg.Types.Scope().Lookup(name)
	// }
	// }

	re := regexp.MustCompile(`topic:(\w+)`)

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
							doc := typeSpec.Doc
							if doc == nil {
								doc = genDecl.Doc
							}
							match := re.FindStringSubmatch(doc.Text())
							if len(match) != 2 {
								return "", fmt.Errorf("%w, event: %s", ErrTopicNotFound, name)
							}
							topic := match[1]
							events = append(events, Event{Name: name, Topic: topic})
						}
					}
				}
			}
		}
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
