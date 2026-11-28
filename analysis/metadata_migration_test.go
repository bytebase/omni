package analysis

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestWithoutSchemas(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Name: "a"}, {Name: "b"}, {Name: "c"}}}
	if got := WithoutSchemas(meta, nil); got != meta {
		t.Fatalf("no names: got a copy")
	}
	if got := WithoutSchemas(nil, []string{"a"}); got != nil {
		t.Fatalf("nil metadata = %v", got)
	}
	got := WithoutSchemas(meta, []string{"b", "missing"})
	if len(got.GetSchemas()) != 2 || got.GetSchemas()[0].GetName() != "a" || got.GetSchemas()[1].GetName() != "c" {
		t.Fatalf("schemas = %v", got.GetSchemas())
	}
	if len(meta.GetSchemas()) != 3 {
		t.Fatalf("input changed: %v", meta.GetSchemas())
	}
}
