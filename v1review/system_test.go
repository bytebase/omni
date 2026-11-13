package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestSchemaDatabaseName(t *testing.T) {
	if got := (Schema{}).DatabaseName(); got != "" {
		t.Errorf("no snapshot: got %q, want none", got)
	}
	if got := (Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: "Shop"}}).DatabaseName(); got != "Shop" {
		t.Errorf("got %q, want Shop", got)
	}
}

func TestSchemaHasNoTable(t *testing.T) {
	table := []*metadata.TableMetadata{{Name: "t"}}
	for _, tt := range []struct {
		name          string
		caseSensitive bool
		schemas       []*metadata.SchemaMetadata
		want          bool
	}{
		{name: "no snapshot", want: true},
		{name: "no table", schemas: []*metadata.SchemaMetadata{{Name: ""}, {Name: "other"}}, want: true},
		{name: "a table in the only schema", schemas: []*metadata.SchemaMetadata{{Name: "", Tables: table}}},
		{name: "a table in a later schema", schemas: []*metadata.SchemaMetadata{{Name: "a"}, {Name: "b", Tables: table}}},
		{name: "a view alone is no table", schemas: []*metadata.SchemaMetadata{{Name: "", Views: []*metadata.ViewMetadata{{Name: "v"}}}}, want: true},
		{
			name:    "a later empty schema of the same name hides the earlier",
			schemas: []*metadata.SchemaMetadata{{Name: "a", Tables: table}, {Name: "a"}},
			want:    true,
		},
		{
			name:    "names differing in case are one schema when case insensitive",
			schemas: []*metadata.SchemaMetadata{{Name: "A", Tables: table}, {Name: "a"}},
			want:    true,
		},
		{
			name:          "names differing in case are two schemas when case sensitive",
			caseSensitive: true,
			schemas:       []*metadata.SchemaMetadata{{Name: "A", Tables: table}, {Name: "a"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := Schema{ObjectCaseSensitive: tt.caseSensitive}
			if tt.schemas != nil {
				s.Metadata = &metadata.DatabaseSchemaMetadata{Schemas: tt.schemas}
			}
			if got := s.HasNoTable(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
