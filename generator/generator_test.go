package generator

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func g(pkg, s string) *Generator {
	return NewGenerator(Options{PackageName: pkg, ServiceName: s})
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name string
		g    *Generator
		path string
		err  error
	}{
		{
			name: "no service",
			g:    g("no_service", "Service"),
			path: "./testdata/no_service",
			err:  ErrServiceNotFound,
		},
		{
			name: "empty service",
			g:    g("empty_service", "Service"),
			path: "./testdata/empty_service",
			err:  ErrEmptyService,
		},
		{
			name: "not an interface",
			g:    g("not_interface", "Service"),
			path: "./testdata/not_interface",
			err:  ErrNotInterface,
		},
		{
			name: "first param should be context.Context",
			g:    g("no_ctx_param", "Service"),
			path: "./testdata/no_ctx_param",
			err:  ErrNoCtxParam,
		},
		{
			name: "invalid results",
			g:    g("invalid_results", "Service"),
			path: "./testdata/invalid_results",
			err:  ErrInvalidResults,
		},
		{
			name: "example service",
			g:    g("example_service", "Example"),
			path: "./testdata/example_service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := "github.com/purefun/go-gen-dapr/generator/" + strings.TrimPrefix(tt.path, "./")

			t.Log(pkg)

			err := tt.g.Load(&packages.Config{Mode: LoadMode}, pkg)
			require.NoError(t, err)

			got, err := tt.g.Generate()

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
				return
			}

			outContent, err := ioutil.ReadFile(tt.path + "/output.go")
			require.NoError(t, err)

			want := string(outContent)
			require.Equal(t, want, got)

		})
	}
}
