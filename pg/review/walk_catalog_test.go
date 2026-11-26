package review

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// The walk takes a refusal from pg/catalog when the synced schema loads in
// full, and only for a code no unsynced object can cause or cure.
func TestWalkThroughCatalogRefusals(t *testing.T) {
	target := review.Target{SessionUser: "u", Schema: &metadata.DatabaseSchemaMetadata{
		Name:       "db",
		SearchPath: `"$user", public`,
		Schemas: []*metadata.SchemaMetadata{{
			Name: "public",
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "n", Type: "integer", Nullable: true}}}, {
				Name:    "k",
				Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "integer"}, {Name: "b", Type: "integer"}},
				Indexes: []*metadata.IndexMetadata{{Name: "k_pkey", Expressions: []string{"a", "b"}, Primary: true, Unique: true, IsConstraint: true, Type: "btree", Definition: "CREATE UNIQUE INDEX k_pkey ON public.k USING btree (a, b)"}},
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name:       "f_double",
				Definition: "CREATE OR REPLACE FUNCTION public.f_double(x integer)\n RETURNS integer\n LANGUAGE sql\nAS $function$SELECT x * 2$function$\n",
			}, {
				Name:       "p",
				Definition: "CREATE OR REPLACE PROCEDURE public.p()\n LANGUAGE plpgsql\nAS $procedure$BEGIN CREATE VIEW v AS SELECT p FROM now(); END$procedure$\n",
			}, {
				Name:       "r",
				Definition: "CREATE OR REPLACE PROCEDURE public.r()\n LANGUAGE plpgsql\nAS $procedure$BEGIN CALL r('1'::text); END$procedure$\n",
			}, {
				Name:       "r",
				Definition: "CREATE OR REPLACE PROCEDURE public.r(a integer)\n LANGUAGE plpgsql\nAS $procedure$BEGIN NULL; END$procedure$\n",
			}, {
				Name:       "w",
				Definition: "CREATE OR REPLACE PROCEDURE public.w()\n LANGUAGE plpgsql\nAS $procedure$BEGIN CALL nope(); END$procedure$\n",
			}, {
				Name:       "f_stable",
				Definition: "CREATE OR REPLACE FUNCTION public.f_stable(x integer)\n RETURNS integer\n LANGUAGE sql\n STABLE\nAS $function$SELECT x$function$\n",
			}, {
				Name:       "trg",
				Definition: "CREATE OR REPLACE FUNCTION public.trg()\n RETURNS trigger\n LANGUAGE plpgsql\nAS $function$BEGIN RETURN NEW; END$function$\n",
			}},
		}},
	}}
	for _, tc := range []struct {
		sql, message string
		statement    int
	}{
		{sql: "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;", message: `PostgreSQL refuses this statement: syntax error at or near "end of input", expected END (SQLSTATE 42601).`},
		{sql: "CREATE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y';", message: `PostgreSQL refuses this statement: function "f_double" already exists (SQLSTATE 42723).`},
		{sql: "CALL f_double('1');", message: "PostgreSQL refuses this statement: f_double(unknown) is not a procedure (SQLSTATE 42809)."},
		{sql: "CALL nope('x', NULL);", message: "PostgreSQL refuses this statement: procedure nope(unknown, unknown) does not exist (SQLSTATE 42883)."},
		{sql: "CALL nope();", message: "PostgreSQL refuses this statement: procedure nope() does not exist (SQLSTATE 42883)."},
		{sql: "CREATE TRIGGER tr BEFORE UPDATE OF nope ON t FOR EACH ROW EXECUTE FUNCTION f_double();", message: `PostgreSQL refuses this statement: function "f_double" does not exist (SQLSTATE 42883).`},
		{sql: "ALTER TABLE t ALTER n TYPE text USING missing::text;", message: `PostgreSQL refuses this statement: column "missing" does not exist (SQLSTATE 42703).`},
		{sql: "ALTER TABLE t ALTER n TYPE json USING JSON_ARRAYAGG(n);", message: "PostgreSQL refuses this statement: aggregate functions are not allowed in transform expressions (SQLSTATE 42803)."},
		{sql: "DO $$ BEGIN RAISE EXCEPTION 'boom'; END $$;", message: "PostgreSQL refuses this statement: boom (SQLSTATE P0001)."},
		{sql: "CREATE INDEX ON t (n) WITH (nope);", message: `PostgreSQL refuses this statement: unrecognized parameter "nope" (SQLSTATE 22023).`},
		{sql: "ALTER TABLE t ADD UNIQUE (n) WITH (nope);", message: `PostgreSQL refuses this statement: unrecognized parameter "nope" (SQLSTATE 22023).`},
		{sql: "ALTER TABLE t ADD FOREIGN KEY (n) REFERENCES t (n);", message: `PostgreSQL refuses this statement: there is no unique constraint matching given keys for referenced table "t" (SQLSTATE 42830).`},
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
		// A query may call an aggregate, which the sync does not list.
		"CREATE TABLE u AS SELECT nope_agg(n) FROM t;",
		// PL/pgSQL's %TYPE takes a system column, and a.b%TYPE may name a
		// domain b in a schema a.
		"CREATE FUNCTION g() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x t.nope%TYPE; BEGIN RETURN 1; END $$;",
		// A qualified column may be a function called in attribute notation.
		"ALTER TABLE t ALTER n TYPE integer USING t.f_double;",
		"CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW WHEN (OLD.tableoid IS NOT NULL) EXECUTE FUNCTION trg();",
		// A refusal in an expression is not the lookup of the trigger's function.
		"CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW WHEN (trg(NEW.n) > 0) EXECUTE FUNCTION trg();",
		// A schema the sync leaves out may hold what a name resolves to.
		"SET search_path = cron, public;\nCALL nope();",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION partman.nope();",
		"ALTER FUNCTION f_double(int) SUPPORT partman.f_double;",
		// A refusal of what a procedure's body runs.
		"CALL w();",
		// ALTER ROUTINE also finds an aggregate.
		"ALTER ROUTINE nope_agg(int) PARALLEL SAFE;",
		"CREATE FUNCTION g() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x t.ctid%TYPE; BEGIN RETURN 1; END $$;",
		// The sync lists no operator.
		"CREATE VIEW v AS SELECT 1 OPERATOR(public.##) 2 AS n;",
		// The sync leaves out a routine named as a system one, which a body may call.
		"DO $$ BEGIN CALL public.abs(); END $$;",
		// An implicit cast the sync does not list may let an argument reach a procedure.
		"CALL f_double('1'::text);",
		"CALL f_double(1);",
		"CALL nope(1);",
		// A user's sum may stand before pg_catalog's, as a DO block may set.
		"DO $$ BEGIN SET search_path = public, pg_catalog; ALTER TABLE t ALTER n TYPE integer USING sum(n); END $$;",
		"SET search_path = public, pg_catalog;\nALTER TABLE t ALTER n TYPE json USING JSON_ARRAYAGG(n);",
		// A DO block may call a routine the sync leaves out.
		"DO $$ BEGIN CREATE VIEW v AS SELECT n FROM now(); END $$;",
		"CREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT n FROM now()';",
		"CALL p();",
		// A nested CALL may reach its procedure through an unsynced cast.
		"CALL r();",
		"CREATE FUNCTION g() RETURNS int LANGUAGE plpgsql AS $$ BEGIN CALL now(); RETURN 1; END $$;",
		// DROP ROUTINE also drops an aggregate, which the sync does not list.
		"DROP ROUTINE nope_agg(int);",
		// A call may be a function-style cast to a domain the sync does not list.
		"CALL f_double(nope_domain(1));",
		"CREATE PROCEDURE q() LANGUAGE sql AS 'CALL f_double(nope_domain(1))';",
		// An extension's access method takes its own storage parameters.
		"CREATE INDEX ON t USING bloom (n) WITH (length = 80);",
		// A foreign key may list a unique key's columns in any order.
		"ALTER TABLE k ADD FOREIGN KEY (b, a) REFERENCES k (b, a);",
		// PostgreSQL inlines a SQL function before it checks an index's volatility.
		"CREATE INDEX ON t (f_stable(n));",
		// A storage parameter's value is parsed as PostgreSQL's parse_int does.
		"CREATE INDEX ON t (n) WITH (fillfactor = '0x50');",
		// An enum's labels are not judged.
		"CREATE TYPE e AS ENUM ('');\nALTER TYPE e RENAME VALUE '' TO 'x';\nALTER TYPE e ADD VALUE '';",
		// A routine's options are not judged by their 22023.
		"CREATE FUNCTION g() RETURNS int LANGUAGE sql COST -1 AS 'SELECT 1';",
		// The catalog does not know which index a foreign key uses.
		"ALTER TABLE t DROP CONSTRAINT t_pkey;",
		// The change turns the body check off.
		"SET check_function_bodies = off;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// After a statement that runs user code, the catalog is not asked.
		"SELECT f_double(1);\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// Nor where the walk stops: RESET restores a default it does not know.
		"RESET search_path;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		"CREATE TABLE x AS SELECT f_double(1);\nCREATE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y';",
		// A rollback undoes what the catalog does not.
		"BEGIN;\nCREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nROLLBACK;\nCREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1';",
		// The replay starts on the search path the change starts on.
		"CREATE SCHEMA s;\nCREATE FUNCTION s.g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nSET search_path = s;",
		// A SET LOCAL ends with its transaction.
		"BEGIN;\nSET LOCAL search_path = s;\nCOMMIT;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// set_config moves the search path as SET does.
		"CREATE SCHEMA s;\nSELECT pg_catalog.set_config('search_path', 's', false);\nCREATE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y';",
		// The sync leaves out a routine named as a system one may be.
		"CALL now();",
		// SELECT INTO creates a table the replay does not.
		"CREATE SCHEMA s;\nSELECT 1 AS a INTO s.t2;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
		// A SET LOCAL of check_function_bodies ends with its transaction.
		"SET check_function_bodies = off;\nBEGIN;\nSET LOCAL check_function_bodies = on;\nCOMMIT;\nCREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'not SQL';",
		// A routine's own SET applies while PostgreSQL checks its body.
		"CREATE FUNCTION f() RETURNS int LANGUAGE sql SET search_path = public AS 'SELECT n, n FROM t';",
		// A cast the sync does not list may allow the coercion.
		"CREATE FUNCTION f() RETURNS date LANGUAGE sql AS 'SELECT 20260101';",
		// A refusal a cast may cure ends the replay.
		"ALTER TABLE t ALTER n TYPE boolean;\nCREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
	// "$user" names a schema of a user the walk does not know.
	unknown := target
	unknown.SessionUser = ""
	if got := walkFindings(t, "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;", unknown); len(got) != 0 {
		t.Errorf("no session user: findings = %+v, want none", got)
	}
	// An event trigger may fire on the statement.
	triggered := target
	triggered.Schema = proto.Clone(target.Schema).(*metadata.DatabaseSchemaMetadata)
	triggered.Schema.EventTriggers = []*metadata.EventTriggerMetadata{{Name: "e", Event: "ddl_command_end", Tags: []string{"CREATE TABLE"}, FunctionSchema: "public", FunctionName: "evt", Enabled: true}}
	if got := walkFindings(t, "CREATE TABLE u (a int);\nCREATE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y';", triggered); len(got) != 0 {
		t.Errorf("event trigger: findings = %+v, want none", got)
	}
	// An extension may add an access method with storage parameters of its
	// own.
	extended := target
	extended.Schema = proto.Clone(target.Schema).(*metadata.DatabaseSchemaMetadata)
	extended.Schema.Extensions = []*metadata.ExtensionMetadata{{Name: "bloom", Schema: "public", Version: "1.0"}}
	if got := walkFindings(t, "CREATE INDEX ON t USING bloom (n) WITH (length = 80);\nCREATE INDEX ON t (n) WITH (nope);", extended); len(got) != 0 {
		t.Errorf("bloom: findings = %+v, want none", got)
	}
	// The sync does not list a range type's constructor.
	ranged := target
	ranged.Schema = proto.Clone(target.Schema).(*metadata.DatabaseSchemaMetadata)
	ranged.Schema.Schemas[0].RangeTypes = []*metadata.RangeTypeMetadata{{Name: "floatrange", Subtype: "double precision"}}
	if got := walkFindings(t, "ALTER FUNCTION floatrange(double precision, double precision) COST 10;", ranged); len(got) != 0 {
		t.Errorf("range constructor: findings = %+v, want none", got)
	}
	// Another spec may change the database first.
	other := target
	other.OtherSpecs = true
	if got := walkFindings(t, "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$;", other); len(got) != 0 {
		t.Errorf("OtherSpecs: findings = %+v, want none", got)
	}
}
