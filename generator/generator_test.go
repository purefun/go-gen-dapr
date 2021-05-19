package generator

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var snapshot bool

func TestMain(m *testing.M) {
	flag.BoolVar(&snapshot, "snapshot", false, "")
	os.Exit(m.Run())
}

func TestGenerator(t *testing.T) {
	tests := []struct {
		name    string
		iface   string
		path    string
		genType GenerateType
		err     error
	}{
		{
			name:    "no service",
			iface:   "InvalidService",
			path:    "./testdata/services/no_service",
			genType: GenerateTypeService,
			err:     ErrServiceNotFound,
		},
		{
			name:    "empty service",
			iface:   "InvalidService",
			path:    "./testdata/services/empty_service",
			genType: GenerateTypeService,
			err:     ErrEmptyService,
		},
		{
			name:    "not an interface",
			iface:   "InvalidService",
			path:    "./testdata/services/not_interface",
			genType: GenerateTypeService,
			err:     ErrNotInterface,
		},
		{
			name:    "first param should be context.Context",
			iface:   "InvalidService",
			path:    "./testdata/services/no_ctx_param",
			genType: GenerateTypeService,
			err:     ErrNoCtxParam,
		},
		{
			name:    "invalid results",
			iface:   "InvalidService",
			path:    "./testdata/services/invalid_results",
			genType: GenerateTypeService,
			err:     ErrInvalidResults,
		},
		{
			name:    "no param no response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/no_param_no_response",
		},
		{
			name:    "basic type params no response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/basic_params_no_response",
		},
		{
			name:    "no param with basic type response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/no_param_with_basic_response",
		},
		{
			name:    "basic params with basic response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/basic_params_with_basic_response",
		},
		{
			name:    "struct params with struct response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/struct_params_with_struct_response",
		},
		{
			name:    "slice response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/slice_response",
		},
		{
			name:    "external types",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/external_types",
		},
		{
			name:    "variadic",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/variadic",
		},
		{
			name:    "interface{} response",
			iface:   "Example",
			genType: GenerateTypeService,
			path:    "./testdata/services/interface_response",
		},
		{
			name:    "subscriptions basic",
			iface:   "Subscriber",
			genType: GenerateTypeSubscriber,
			path:    "./testdata/subscriptions/basic",
		},
		{
			name:    "subscriptions external types",
			iface:   "Subscriptions",
			genType: GenerateTypeSubscriber,
			path:    "./testdata/subscriptions/external_types",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := "github.com/purefun/go-gen-dapr/generator/" + strings.TrimPrefix(tt.path, "./")

			g := NewGenerator(Options{
				ServicePkg:   pkg,
				ServiceType:  tt.iface,
				GenComment:   false,
				GenerateType: tt.genType,
			})

			got, err := g.Generate()

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
				return
			}

			outFile := tt.path + "/output.go"

			outContent, err := ioutil.ReadFile(outFile)
			require.NoError(t, err)

			want := string(outContent)
			assert.Equal(t, want, got)

			if want != got && snapshot {
				err := ioutil.WriteFile(outFile, []byte(got), 0644)
				require.NoError(t, err)
			}

		})
	}
}
