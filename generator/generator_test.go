package generator

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			g:    g("no_service", "InvalidService"),
			path: "./testdata/no_service",
			err:  ErrServiceNotFound,
		},
		{
			name: "empty service",
			g:    g("empty_service", "InvalidService"),
			path: "./testdata/empty_service",
			err:  ErrEmptyService,
		},
		{
			name: "not an interface",
			g:    g("not_interface", "InvalidService"),
			path: "./testdata/not_interface",
			err:  ErrNotInterface,
		},
		{
			name: "first param should be context.Context",
			g:    g("no_ctx_param", "InvalidService"),
			path: "./testdata/no_ctx_param",
			err:  ErrNoCtxParam,
		},
		{
			name: "invalid results",
			g:    g("invalid_results", "InvalidService"),
			path: "./testdata/invalid_results",
			err:  ErrInvalidResults,
		},
		{
			name: "no param no response",
			g:    g("no_param_no_response", "Example"),
			path: "./testdata/no_param_no_response",
		},
		{
			name: "params no response",
			g:    g("params_no_response", "Example"),
			path: "./testdata/params_no_response",
		},
		{
			name: "no param with response",
			g:    g("no_param_with_response", "Example"),
			path: "./testdata/no_param_with_response",
		},
		{
			name: "params with response",
			g:    g("params_with_response", "Example"),
			path: "./testdata/params_with_response",
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
			assert.Equal(t, want, got)

		})
	}
}
