package cassandra_util

import (
	"testing"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
)

func _Test_create_statement(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		params *out.GetParams
		result string
	}{
		"default": {
			params: &out.GetParams{
				Limit: nil,
				Page:  nil,
				Tag:   nil,
				Name:  nil,
			},
			result: "SELECT id, name, tags FROM cv_experiences LIMIT 100",
		},
		"limit": {
			params: &out.GetParams{
				Limit: Pointer(int32(10)),
				Page:  nil,
				Tag:   nil,
				Name:  nil,
			},
			result: "SELECT id, name, tags FROM cv_experiences LIMIT 10",
		},
		"page": {
			params: &out.GetParams{
				Limit: nil,
				Page:  Pointer("page"),
				Tag:   nil,
				Name:  nil,
			},
			result: "SELECT id, name, tags FROM cv_experiences WHERE page > 'page' LIMIT 100",
		},
		"tag": {
			params: &out.GetParams{
				Limit: nil,
				Page:  nil,
				Tag:   Pointer("tag"),
				Name:  nil,
			},
			result: "SELECT id, name, tags FROM cv_experiences WHERE tags CONTAINS 'tag' LIMIT 100",
		},
		"name": {
			params: &out.GetParams{
				Limit: nil,
				Page:  nil,
				Tag:   Pointer("tag"),
				Name:  nil,
			},
			result: "SELECT id, name, tags FROM cv_experiences WHERE name CONTAINS 'name' LIMIT 100",
		},
	}

	for name, params := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, expected := GetStatement(params.params), params.result; got != expected {
				t.Fatalf("Test %s - GetStatement(%+v) returned %q;\n expected %q\n", name, params.params, got, expected)
			}
		})
	}

}

func Pointer[K any](val K) *K {
	return &val
}
