package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// The walk takes a refusal from pg/catalog for a routine or a view when
// the synced schema loads in full, and only for a code no unsynced object
// can cause.
func TestWalkThroughCatalogRefusals(t *testing.T) {
	target := review.Target{Schema: &metadata.DatabaseSchemaMetadata{
		Name:       "db",
		SearchPath: `"$user", public`,
		Schemas: []*metadata.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "n", Type: "integer", Nullable: true}}}},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "f_double",
				Definition: "CREATE OR REPLACE FUNCTION public.f_double(x integer)\n RETURNS integer\n LANGUAGE sql\nAS $function$SELECT x * 2$function$\n",
			}},
		}},
	}}
	for _, tc := range []struct {
		sql, message string
		statement    int
	}{
		{sql: "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;", message: `PostgreSQL refuses this statement: syntax error at or near "end of input", expected END (SQLSTATE 42601).`},
		{sql: "CREATE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y';", message: `PostgreSQL refuses this statement: function "f_double" already exists (SQLSTATE 42723).`},
		{sql: "CALL f_double(1);", message: "PostgreSQL refuses this statement: f_double(integer) is not a procedure (SQLSTATE 42809)."},
		{sql: "CREATE TABLE u (a int);\nCREATE FUNCTION g() RETURNS void LANGUAGE plpgsql AS $$ BEGIN RETURN 1; END $$;", statement: 1, message: "PostgreSQL refuses this statement: RETURN cannot have a parameter in function returning void (SQLSTATE 42804)."},
	} {
		got := walkFindings(t, tc.sql, target)
		if len(got) != 1 || got[0].Statement != tc.statement || got[0].Message != tc.message {
			t.Errorf("%q: findings = %+v, want %q at %d", tc.sql, got, tc.message, tc.statement)
		}
	}
	for _, sql := range []string{
		// A missing name may be an object the sync does not list.
		"CREATE FUNCTION f() RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM missing_t';",
		"CREATE VIEW v AS SELECT nope_fn(n) FROM t;",
		// The change turns the body check off.
		"SET check_function_bodies = off;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// After DML, which may run user code, the catalog is not asked.
		"INSERT INTO t VALUES (1);\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// A refusal of another statement ends the replay.
		"ALTER TABLE t ALTER n TYPE boolean;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
	// Another spec may change the database first.
	other := target
	other.OtherSpecs = true
	if got := walkFindings(t, "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;", other); len(got) != 0 {
		t.Errorf("OtherSpecs: findings = %+v, want none", got)
	}
}
