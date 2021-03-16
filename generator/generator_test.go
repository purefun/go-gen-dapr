package generator

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name  string
		iface string
		path  string
		err   error
	}{
		{
			name:  "no service",
			iface: "InvalidService",
			path:  "./testdata/no_service",
			err:   ErrServiceNotFound,
		},
		{
			name:  "empty service",
			iface: "InvalidService",
			path:  "./testdata/empty_service",
			err:   ErrEmptyService,
		},
		{
			name:  "not an interface",
			iface: "InvalidService",
			path:  "./testdata/not_interface",
			err:   ErrNotInterface,
		},
		{
			name:  "first param should be context.Context",
			iface: "InvalidService",
			path:  "./testdata/no_ctx_param",
			err:   ErrNoCtxParam,
		},
		{
			name:  "invalid results",
			iface: "InvalidService",
			path:  "./testdata/invalid_results",
			err:   ErrInvalidResults,
		},
		{
			name:  "no param no response",
			iface: "Example",
			path:  "./testdata/no_param_no_response",
		},
		{
			name:  "basic type params no response",
			iface: "Example",
			path:  "./testdata/basic_params_no_response",
		},
		{
			name:  "no param with basic type response",
			iface: "Example",
			path:  "./testdata/no_param_with_basic_response",
		},
		{
			name:  "basic params with basic response",
			iface: "Example",
			path:  "./testdata/basic_params_with_basic_response",
		},
		{
			name:  "struct params with struct response",
			iface: "Example",
			path:  "./testdata/struct_params_with_struct_response",
		},
		{
			name:  "external types",
			iface: "Example",
			path:  "./testdata/external_types",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := "github.com/purefun/go-gen-dapr/generator/" + strings.TrimPrefix(tt.path, "./")

			g := NewGenerator(Options{
				PackageName: filepath.Base(pkg),
				ServicePkg:  pkg,
				ServiceType: tt.iface,
				GenComment:  false,
			})

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
