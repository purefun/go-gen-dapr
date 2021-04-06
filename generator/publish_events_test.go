package generator

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePublishEvents(t *testing.T) {
	tests := []struct {
		name string
		path string
		err  error
	}{
		{
			name: "basic",
			path: "./testdata/events/basic",
		},
	}

	for _, tt := range tests {
		pkg := "github.com/purefun/go-gen-dapr/generator/" + strings.TrimPrefix(tt.path, "./")
		options := PublishEventsOptions{
			Pkg: pkg,
		}
		got, err := GeneratePublishEvents(options)
		require.NoError(t, err)

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

		require.Equal(t, want, got)
	}

}
