package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestSchemaTableIndexesByName(t *testing.T) {
	a := &metadata.IndexMetadata{Name: "a"}
	b := &metadata.IndexMetadata{Name: "b"}
	upperA := &metadata.IndexMetadata{Name: "A"}
	db := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{
			Name:       "t",
			Indexes:    []*metadata.IndexMetadata{a, b, upperA},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p"}},
		}},
	}}}
	s := Schema{Metadata: db}
	if got := s.TableIndexesByName("", "t", false); len(got) != 2 || got[0] != upperA || got[1] != b {
		t.Errorf("case-insensitive: got %v", got)
	}
	if got := s.TableIndexesByName("", "t", true); len(got) != 3 {
		t.Errorf("case-sensitive: got %v", got)
	}
	if got := s.TableIndexesByName("", "p", false); got != nil {
		t.Errorf("partition: got %v", got)
	}
	if got := s.TableIndexesByName("", "missing", false); got != nil {
		t.Errorf("missing: got %v", got)
	}
}
