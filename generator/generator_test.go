package generator

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name string
		s    string
		path string
		err  error
	}{
		{
			name: "no service",
			s:    "InvalidService",
			path: "./testdata/no_service",
			err:  ErrServiceNotFound,
		},
		{
			name: "empty service",
			s:    "InvalidService",
			path: "./testdata/empty_service",
			err:  ErrEmptyService,
		},
		{
			name: "not an interface",
			s:    "InvalidService",
			path: "./testdata/not_interface",
			err:  ErrNotInterface,
		},
		{
			name: "first param should be context.Context",
			s:    "InvalidService",
			path: "./testdata/no_ctx_param",
			err:  ErrNoCtxParam,
		},
		{
			name: "invalid results",
			s:    "InvalidService",
			path: "./testdata/invalid_results",
			err:  ErrInvalidResults,
		},
		{
			name: "no param no response",
			s:    "Example",
			path: "./testdata/no_param_no_response",
		},
		{
			name: "basic type params no response",
			s:    "Example",
			path: "./testdata/basic_params_no_response",
		},
		{
			name: "no param with basic type response",
			s:    "Example",
			path: "./testdata/no_param_with_basic_response",
		},
		{
			name: "basic type params with basic response",
			s:    "Example",
			path: "./testdata/basic_params_with_basic_response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := "github.com/purefun/go-gen-dapr/generator/" + strings.TrimPrefix(tt.path, "./")

			g := NewGenerator(Options{PackageName: filepath.Base(pkg), ServiceName: tt.s})

			err := g.Load(&packages.Config{Mode: LoadMode}, pkg)
			require.NoError(t, err)

			got, err := g.Generate()

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
