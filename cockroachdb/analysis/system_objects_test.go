package analysis

import "testing"

func TestSystemObjects(t *testing.T) {
	for _, tc := range []struct {
		name   string
		got    bool
		wanted bool
	}{
		{name: "system database postgres", got: IsSystemDatabase("postgres"), wanted: true},
		{name: "system database system", got: IsSystemDatabase("system"), wanted: true},
		{name: "user database", got: IsSystemDatabase("defaultdb"), wanted: false},
		{name: "crdb_internal schema", got: IsSystemSchema("crdb_internal"), wanted: true},
		{name: "information_schema schema", got: IsSystemSchema("information_schema"), wanted: true},
		{name: "pg_catalog schema", got: IsSystemSchema("pg_catalog"), wanted: true},
		{name: "pg_extension schema", got: IsSystemSchema("pg_extension"), wanted: true},
		{name: "user schema", got: IsSystemSchema("public"), wanted: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.wanted {
				t.Errorf("got %v, want %v", tc.got, tc.wanted)
			}
		})
	}
}

func TestSystemSchemaWhereClause(t *testing.T) {
	const want = "'crdb_internal','information_schema','pg_catalog','pg_extension'"
	if SystemSchemaWhereClause != want {
		t.Errorf("SystemSchemaWhereClause = %q, want %q", SystemSchemaWhereClause, want)
	}
}
