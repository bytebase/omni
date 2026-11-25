package catalog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
)

func TestLoadSDLRefuses(t *testing.T) {
	for name, sdl := range map[string]string{
		"syntax":           "CREATE TABLE (",
		"data":             "INSERT INTO t VALUES (1)",
		"missing table":    "CREATE INDEX i ON app.t (a)",
		"missing column":   "CREATE TABLE t (a INT8); COMMENT ON COLUMN t.b IS 'x'",
		"table as select":  "CREATE TABLE t AS SELECT 1",
		"unknown sequence": "ALTER SEQUENCE s INCREMENT BY 2",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := catalog.LoadSDL(sdl)
			require.Error(t, err)
		})
	}
}

func TestBuildSDLMigrationPlan(t *testing.T) {
	source := `SET statement_timeout = 0;
SELECT pg_catalog.set_config('search_path', '', false);
CREATE TABLE public.users (id INT8 PRIMARY KEY, team INT8 REFERENCES public.teams);
CREATE TABLE public.teams (id INT8 PRIMARY KEY);`
	plan, err := catalog.BuildSDLMigrationPlan(source, source, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)

	_, err = catalog.BuildSDLMigrationPlan("CREATE TABLE (", source, catalog.SDLMigrationOptions{})
	require.ErrorContains(t, err, "failed to load source schema")
	_, err = catalog.BuildSDLMigrationPlan(source, "CREATE TABLE (", catalog.SDLMigrationOptions{})
	require.ErrorContains(t, err, "failed to load target schema")

	target := `CREATE TABLE public.users (id INT8 PRIMARY KEY, email STRING);
CREATE TABLE public.teams (id INT8 PRIMARY KEY);`
	plan, err = catalog.BuildSDLMigrationPlan(source, target, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Equal(t, "ALTER TABLE \"public\".\"users\" DROP CONSTRAINT \"users_team_fkey\";\n"+
		"ALTER TABLE \"public\".\"users\" DROP COLUMN \"team\";\n"+
		"ALTER TABLE \"public\".\"users\" ADD COLUMN \"email\" STRING", plan.SQL())
	var types []catalog.MigrationOpType
	for _, op := range plan.Ops {
		types = append(types, op.Type)
	}
	require.Equal(t, []catalog.MigrationOpType{catalog.OpDropConstraint, catalog.OpDropColumn, catalog.OpAddColumn}, types)
}

// TestLoadSDLReadsDump reads the DDL cockroachdb/driver dumps (SHOW CREATE),
// which Bytebase hands the SDL diff as a raw schema, and requires the SDL
// rendered from it to hold the same schema.
func TestLoadSDLReadsDump(t *testing.T) {
	dump, err := os.ReadFile(filepath.Join("..", "driver", "testdata", "dump", "default.sql"))
	require.NoError(t, err)
	loaded, err := catalog.LoadSDL(string(dump))
	require.NoError(t, err)
	sdl, err := catalog.DatabaseDefinition(loaded, catalog.DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	require.NotContains(t, sdl, "omni_testdump")
	plan, err := catalog.BuildSDLMigrationPlan(sdl, string(dump), catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "SDL:\n%s\ndiffers from the dump by:\n%s", sdl, plan.SQL())
}

// TestLoadSDLNamesChecks names unnamed checks as v25.2.3 does: a table's own
// checks before its columns' checks, each check_ and its columns, then 1, 2
// and on past the names taken.
func TestLoadSDLNamesChecks(t *testing.T) {
	loaded, err := catalog.LoadSDL(`CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0) CHECK (a < 10), CHECK (a <> 5), b INT8 CHECK (b > 0), CONSTRAINT check_b1 CHECK (b < 9));
ALTER TABLE t ADD CHECK (a <> 7);
ALTER TABLE t ADD CHECK (b <> 7);`)
	require.NoError(t, err)
	names := make(map[string]string)
	for _, check := range loaded.Schemas[0].Tables[0].CheckConstraints {
		names[check.Expression] = check.Name
	}
	require.Equal(t, map[string]string{
		"((a != 5))": "check_a",
		"((a > 0))":  "check_a1",
		"((a < 10))": "check_a2",
		"((b > 0))":  "check_b",
		"((b < 9))":  "check_b1",
		"((a != 7))": "check_a3",
		"((b != 7))": "check_b2",
	}, names)
}

// TestSDLViewWithCTE orders a view after the stored view it reads, not after
// a stored view its CTE shares a name with.
func TestSDLViewWithCTE(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan("", `CREATE TABLE public.t (id INT8 PRIMARY KEY);
CREATE VIEW public.z AS WITH a AS (SELECT id FROM public.t) SELECT id FROM a;
CREATE VIEW public.a AS SELECT id FROM public.z;`, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	sql := plan.SQL()
	require.Less(t, strings.Index(sql, `CREATE VIEW "public"."z"`), strings.Index(sql, `CREATE VIEW "public"."a"`), sql)
}

// TestSDLSequenceAlterBounds holds ALTER SEQUENCE's bounds to v25.2.3: NO
// MINVALUE and NO MAXVALUE reset to INT8's defaults, and a bound left at its
// type's limit widens to INT8's.
func TestSDLSequenceAlterBounds(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan("CREATE SEQUENCE s MINVALUE 5;", "CREATE SEQUENCE s MINVALUE 5; ALTER SEQUENCE s NO MINVALUE;", catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.NotNil(t, plan)
	require.Contains(t, plan.SQL(), "MINVALUE 1")

	loaded, err := catalog.LoadSDL(`CREATE SEQUENCE m AS INT2 INCREMENT 3 MINVALUE 5 START 7;
ALTER SEQUENCE m INCREMENT -2 NO MINVALUE;
CREATE SEQUENCE d AS INT4 INCREMENT -1 MINVALUE -50 MAXVALUE -5 START -10;
ALTER SEQUENCE d NO MINVALUE NO MAXVALUE;
CREATE SEQUENCE p AS INT4;
ALTER SEQUENCE p MAXVALUE 2147483647 CACHE 2;`)
	require.NoError(t, err)
	bounds := make(map[string][2]string)
	for _, sequence := range loaded.Schemas[0].Sequences {
		bounds[sequence.Name] = [2]string{sequence.MinValue, sequence.MaxValue}
	}
	require.Equal(t, map[string][2]string{
		"m": {"-9223372036854775808", "9223372036854775807"},
		"d": {"-9223372036854775808", "-1"},
		"p": {"1", "2147483647"},
	}, bounds)
}

// TestSDLKeysKeepValues requires the comparison keys to keep what changes a
// value: a cast that rounds a literal, and annotation-like text in a string.
func TestSDLKeysKeepValues(t *testing.T) {
	for name, pair := range map[string][2]string{
		"rounding cast": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5::INT8);",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5::DECIMAL);",
		},
		"string in a body": {
			"CREATE FUNCTION f() RETURNS STRING LANGUAGE SQL AS $$ SELECT 'a:::INT8' $$;",
			"CREATE FUNCTION f() RETURNS STRING LANGUAGE SQL AS $$ SELECT 'a' $$;",
		},
	} {
		t.Run(name, func(t *testing.T) {
			plan, err := catalog.BuildSDLMigrationPlan(pair[0], pair[1], catalog.SDLMigrationOptions{})
			require.NoError(t, err)
			require.NotNil(t, plan)
		})
	}
	for name, pair := range map[string][2]string{
		"stored cast":       {"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'x'::STRING);", "CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'x');"},
		"stored annotation": {"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5:::DECIMAL);", "CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5);"},
		"annotated body":    {"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1:::INT8 + 2:::DECIMAL(10,2)::INT8 $$;", "CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 + 2::INT8 $$;"},
	} {
		t.Run(name, func(t *testing.T) {
			plan, err := catalog.BuildSDLMigrationPlan(pair[0], pair[1], catalog.SDLMigrationOptions{})
			require.NoError(t, err)
			require.Nil(t, plan, "%v", plan)
		})
	}
}

// TestLoadSDLNamesIndexes names unnamed indexes as v25.2.3 does: the table,
// the keys and _idx or _key, then 1, 2 and on past the names taken.
func TestLoadSDLNamesIndexes(t *testing.T) {
	loaded, err := catalog.LoadSDL(`CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, b INT8, INDEX t_a_idx1 (a, b), INDEX (a), UNIQUE (a));
CREATE INDEX ON t (a) STORING (b);
CREATE INDEX ON t (a);
CREATE UNIQUE INDEX ON t (a);`)
	require.NoError(t, err)
	var names []string
	for _, index := range loaded.Schemas[0].Tables[0].Indexes {
		names = append(names, index.Name)
	}
	require.Equal(t, []string{"t_pkey", "t_a_idx1", "t_a_idx", "t_a_key", "t_a_idx2", "t_a_idx3", "t_a_key1"}, names)
}

// TestSDLViewColumns compares a view's output names when one SDL lists them
// and the other leaves them to the query.
func TestSDLViewColumns(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan("CREATE VIEW v (a) AS SELECT 1 AS b;", "CREATE VIEW v AS SELECT 1 AS b;", catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.NotNil(t, plan)
	for name, pair := range map[string][2]string{
		"alias": {"CREATE VIEW v (b) AS SELECT 1 AS b;", "CREATE VIEW v AS SELECT 1 AS b;"},
		"names": {
			`CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE VIEW v ("?column?", a, id, lower, "case", bool) AS SELECT 1, t.a, (id)::STRING, lower('x'), CASE WHEN true THEN 1 END, true FROM t;`,
			`CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE VIEW v AS SELECT 1, t.a, (id)::STRING, lower('x'), CASE WHEN true THEN 1 END, true FROM t;`,
		},
		"star": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE VIEW v (x, y) AS SELECT * FROM t;",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE VIEW v AS SELECT * FROM t;",
		},
		"cast constant": {"CREATE VIEW v (x) AS SELECT 'x'::STRING;", "CREATE VIEW v AS SELECT 'x'::STRING;"},
	} {
		t.Run(name, func(t *testing.T) {
			plan, err := catalog.BuildSDLMigrationPlan(pair[0], pair[1], catalog.SDLMigrationOptions{})
			require.NoError(t, err)
			require.Nil(t, plan, "%v", plan)
		})
	}
}

// TestSDLUnqualifiedType reads an unqualified user-defined type as the one in
// public, which is how the sync names it.
func TestSDLUnqualifiedType(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TYPE public.mood AS ENUM ('ok'); CREATE TABLE t (id INT8 PRIMARY KEY, m public.mood, ms public.mood[]);",
		"CREATE TYPE public.mood AS ENUM ('ok'); CREATE TABLE t (id INT8 PRIMARY KEY, m mood, ms mood[]);",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}
