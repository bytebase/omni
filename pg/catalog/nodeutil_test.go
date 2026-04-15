package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestTypeNamePartsThreeComponent locks down typeNameParts behavior for
// 1/2/3/4-component qualified type names.
//
// Background: prior to this fix, len > 2 returned ("", lastItem),
// silently dropping ALL qualification. That broke 3-component names in
// flows like %TYPE references and parseAnyName-driven CREATE TABLE OF.
// Fix: 3-component is treated as catalog.schema.name → (schema, name);
// 4+ falls back to ("", lastItem) for backward compatibility.
//
// See docs/plans/2026-04-14-pg-followups.md for the audit that
// motivated this change.
func TestTypeNamePartsThreeComponent(t *testing.T) {
	cases := []struct {
		name       string
		items      []string
		wantSchema string
		wantName   string
	}{
		{
			name:       "0-component (nil-equivalent edge)",
			items:      []string{},
			wantSchema: "",
			wantName:   "",
		},
		{
			name:       "1-component bare",
			items:      []string{"int4"},
			wantSchema: "",
			wantName:   "int4",
		},
		{
			name:       "2-component schema-qualified",
			items:      []string{"pg_catalog", "int4"},
			wantSchema: "pg_catalog",
			wantName:   "int4",
		},
		{
			name:       "3-component catalog.schema.name (the fix)",
			items:      []string{"db", "schema", "mytype"},
			wantSchema: "schema",
			wantName:   "mytype",
		},
		{
			name:       "3-component pg_catalog catalog form",
			items:      []string{"localdb", "pg_catalog", "int4"},
			wantSchema: "pg_catalog",
			wantName:   "int4",
		},
		{
			// 4+ components: return empty so downstream resolution
			// fails. We deliberately do NOT fall back to ("", lastItem)
			// — that would silently resolve invalid SQL like `CREATE
			// TABLE t (c a.b.c.int4)` as the local int4 type. See
			// typeNameParts doc comment for the full rationale.
			name:       "4-component returns empty (silent-success guard)",
			items:      []string{"a", "b", "c", "d"},
			wantSchema: "",
			wantName:   "",
		},
		{
			name:       "5-component returns empty",
			items:      []string{"a", "b", "c", "d", "e"},
			wantSchema: "",
			wantName:   "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tn := buildTypeName(tc.items)
			gotSchema, gotName := typeNameParts(tn)
			if gotSchema != tc.wantSchema {
				t.Errorf("schema: got %q, want %q", gotSchema, tc.wantSchema)
			}
			if gotName != tc.wantName {
				t.Errorf("name: got %q, want %q", gotName, tc.wantName)
			}
		})
	}
}

// TestTypeNamePartsNilNames covers the early-return branch.
func TestTypeNamePartsNilNames(t *testing.T) {
	tn := &nodes.TypeName{Names: nil}
	schema, name := typeNameParts(tn)
	if schema != "" || name != "" {
		t.Errorf("nil Names: got (%q, %q), want (\"\", \"\")", schema, name)
	}
}

// buildTypeName constructs a *nodes.TypeName whose Names list contains
// the given string components.
func buildTypeName(items []string) *nodes.TypeName {
	nodeItems := make([]nodes.Node, len(items))
	for i, s := range items {
		nodeItems[i] = &nodes.String{Str: s}
	}
	return &nodes.TypeName{
		Names: &nodes.List{Items: nodeItems},
	}
}
