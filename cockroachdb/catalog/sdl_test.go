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
		"syntax":             "CREATE TABLE (",
		"data":               "INSERT INTO t VALUES (1)",
		"missing table":      "CREATE INDEX i ON app.t (a)",
		"missing column":     "CREATE TABLE t (a INT8); COMMENT ON COLUMN t.b IS 'x'",
		"table as select":    "CREATE TABLE t AS SELECT 1",
		"unknown sequence":   "ALTER SEQUENCE s INCREMENT BY 2",
		"locality":           "CREATE TABLE t (id INT8 PRIMARY KEY) LOCALITY GLOBAL",
		"temporary":          "CREATE TEMP TABLE t (id INT8 PRIMARY KEY)",
		"temporary sequence": "CREATE TEMP SEQUENCE s",
		"identity":           "CREATE TABLE t (id INT8 GENERATED ALWAYS AS IDENTITY PRIMARY KEY)",
		"unique no index":    "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, UNIQUE WITHOUT INDEX (a))",
		"column no index":    "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 UNIQUE WITHOUT INDEX)",
		"partition":          "CREATE TABLE t (id INT8 PRIMARY KEY) PARTITION BY LIST (id) (PARTITION p VALUES IN (1))",
		"families":           "CREATE TABLE t (id INT8 PRIMARY KEY, a STRING, b STRING, FAMILY f1 (id, a), FAMILY f2 (b))",
		"column family":      "CREATE TABLE t (id INT8 PRIMARY KEY, a STRING, b STRING FAMILY f2, FAMILY f1 (id, a))",
		"temporary view":     "CREATE TEMP VIEW v AS SELECT 1 AS n",
		"schema rename":      "CREATE SCHEMA old; CREATE TABLE old.t (id INT8 PRIMARY KEY); ALTER SCHEMA old RENAME TO new",
		"table twice":        "CREATE TABLE t (id INT8 PRIMARY KEY); CREATE TABLE t (id INT8 PRIMARY KEY)",
		"column twice":       "CREATE TABLE t (id INT8 PRIMARY KEY, a STRING); ALTER TABLE t ADD COLUMN a INT8",
		"index twice":        "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE INDEX i ON t (a); CREATE INDEX i ON t (a)",
		"view twice":         "CREATE VIEW v AS SELECT 1 AS n; CREATE VIEW v AS SELECT 1 AS n",
		"view over table":    "CREATE TABLE v (id INT8 PRIMARY KEY); CREATE VIEW v AS SELECT 1 AS n",
		"sequence twice":     "CREATE SEQUENCE s; CREATE SEQUENCE s",
		"type twice":         "CREATE TYPE e AS ENUM ('a'); CREATE TYPE e AS ENUM ('b')",
		"function twice":     "CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 2 $$",
		"not valid":          "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); ALTER TABLE t ADD CONSTRAINT positive CHECK (a > 0) NOT VALID",
		"rowid reference":    "CREATE TABLE p (a INT8); CREATE TABLE c (id INT8 PRIMARY KEY, pid INT8 REFERENCES p)",
		"storage params":     "CREATE TABLE t (id INT8 PRIMARY KEY, expires_at TIMESTAMPTZ) WITH (ttl_expiration_expression = 'expires_at')",
		"on update":          "CREATE TABLE t (id INT8 PRIMARY KEY, updated TIMESTAMPTZ DEFAULT now() ON UPDATE now())",
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
		"decimal scale": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.25::DECIMAL(3,1));",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.25);",
		},
		"float4": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a FLOAT8 DEFAULT 1.1::FLOAT4);",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a FLOAT8 DEFAULT 1.1);",
		},
		"string length": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'abc'::VARCHAR(1));",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'abc'::STRING);",
		},
		"time precision": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a TIMESTAMP DEFAULT '2020-01-01 10:00:00.5'::TIMESTAMP(0));",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a TIMESTAMP DEFAULT '2020-01-01 10:00:00.5'::TIMESTAMP);",
		},
		"cast through a type": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT '01'::INT8::STRING);",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT '01'::STRING);",
		},
		"grouping": {
			"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 DEFAULT (1 + 2) * 3);",
			"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 DEFAULT 1 + 2 * 3);",
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
		"stored cast":           {"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'x'::STRING);", "CREATE TABLE t (id INT8 PRIMARY KEY, a STRING DEFAULT 'x');"},
		"redundant parentheses": {"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, CHECK (((a > 0) AND (a < 9))));", "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, CHECK (a > 0 AND a < 9));"},
		"lossless cast":         {"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.25::DECIMAL, b FLOAT8 DEFAULT 1.1::FLOAT8);", "CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.25, b FLOAT8 DEFAULT 1.1);"},
		"stored annotation":     {"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5:::DECIMAL);", "CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL DEFAULT 1.5);"},
		"annotated body":        {"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1:::INT8 + 2:::DECIMAL(10,2)::INT8 $$;", "CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 + 2::INT8 $$;"},
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

// TestLoadSDLNamesForeignKeys names unnamed foreign keys as v25.2.3 named
// them on this schema: a table's own before its columns', each the table, its
// columns and _fkey, then _1, _2 and on past the constraint names taken.
func TestLoadSDLNamesForeignKeys(t *testing.T) {
	loaded, err := catalog.LoadSDL(`CREATE TABLE p (id INT8 PRIMARY KEY);
CREATE TABLE q (id INT8 PRIMARY KEY);
CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 REFERENCES p (id), FOREIGN KEY (a) REFERENCES q (id), CONSTRAINT t_a_fkey2 FOREIGN KEY (a) REFERENCES p (id), FOREIGN KEY (a) REFERENCES p (id));
ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES q (id);`)
	require.NoError(t, err)
	names := make(map[string]bool)
	for _, table := range loaded.Schemas[0].Tables {
		if table.Name != "t" {
			continue
		}
		var got []string
		for _, fk := range table.ForeignKeys {
			got = append(got, fk.Name+" "+fk.ReferencedTable)
			names[fk.Name] = true
		}
		require.Equal(t, []string{"t_a_fkey q", "t_a_fkey2 p", "t_a_fkey_1 p", "t_a_fkey_2 p", "t_a_fkey_3 q"}, got)
	}
	require.Len(t, names, 5)
}

// TestSDLNaturalJoinRebuildsView rebuilds a view whose NATURAL join reads a
// column the migration drops, though the query never names it.
func TestSDLNaturalJoinRebuildsView(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(`CREATE TABLE p (id INT8 PRIMARY KEY, a INT8);
CREATE TABLE q (other INT8 PRIMARY KEY, a INT8);
CREATE VIEW v AS SELECT p.id FROM p NATURAL JOIN q;`, `CREATE TABLE p (id INT8 PRIMARY KEY);
CREATE TABLE q (other INT8 PRIMARY KEY);
CREATE VIEW v AS SELECT p.id FROM p NATURAL JOIN q;`, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	sql := plan.SQL()
	require.Contains(t, sql, `DROP VIEW "public"."v"`)
	require.Less(t, strings.Index(sql, `DROP VIEW "public"."v"`), strings.Index(sql, `DROP COLUMN "a"`), sql)
}

// TestLoadSDLNamesForeignKeyPastCheck skips a check's name, as v25.2.3 did.
func TestLoadSDLNamesForeignKeyPastCheck(t *testing.T) {
	loaded, err := catalog.LoadSDL(`CREATE TABLE p (id INT8 PRIMARY KEY);
CREATE TABLE u (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0), CONSTRAINT u_a_fkey CHECK (a < 9), FOREIGN KEY (a) REFERENCES p (id));`)
	require.NoError(t, err)
	require.Equal(t, "u_a_fkey_1", loaded.Schemas[0].Tables[1].ForeignKeys[0].Name)
}

// TestSDLViewQualifiedRelations compares a view naming its relations
// unqualified with one naming them in public.
func TestSDLViewQualifiedRelations(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE public.t (id INT8 PRIMARY KEY); CREATE VIEW public.v AS SELECT id FROM public.t; CREATE MATERIALIZED VIEW public.m AS WITH c AS (SELECT id FROM public.t) SELECT id FROM c;",
		"CREATE TABLE t (id INT8 PRIMARY KEY); CREATE VIEW v AS SELECT id FROM t; CREATE MATERIALIZED VIEW m AS WITH c AS (SELECT id FROM t) SELECT id FROM c;",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}

// TestSDLOrReplace holds CREATE OR REPLACE to the last definition.
func TestSDLOrReplace(t *testing.T) {
	for name, pair := range map[string][2]string{
		"view": {"CREATE VIEW v AS SELECT 1 AS n;", "CREATE VIEW v AS SELECT 1 AS n; CREATE OR REPLACE VIEW v AS SELECT 2 AS n;"},
		"function": {
			"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;",
			"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; CREATE OR REPLACE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 2 $$;",
		},
	} {
		t.Run(name, func(t *testing.T) {
			plan, err := catalog.BuildSDLMigrationPlan(pair[0], pair[1], catalog.SDLMigrationOptions{})
			require.NoError(t, err)
			require.NotNil(t, plan)
			require.Contains(t, plan.SQL(), "2")
			loaded, err := catalog.LoadSDL(pair[1])
			require.NoError(t, err)
			require.Equal(t, 1, len(loaded.Schemas[0].Views)+len(loaded.Schemas[0].Functions))
		})
	}
}

// TestSDLUnqualifiedRoutine compares an unqualified routine with the one in
// public.
func TestSDLUnqualifiedRoutine(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE FUNCTION public.f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;",
		"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}

// TestSDLViewCTEScope reads a CTE's own name inside its body as the stored
// relation, and as the CTE in the query the WITH heads.
func TestSDLViewCTEScope(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan("", `CREATE TABLE public.t (id INT8 PRIMARY KEY);
CREATE VIEW public.a AS WITH z AS (SELECT id FROM z) SELECT id FROM z;
CREATE VIEW public.z AS SELECT id FROM public.t;`, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	sql := plan.SQL()
	require.Less(t, strings.Index(sql, `CREATE VIEW "public"."z"`), strings.Index(sql, `CREATE VIEW "public"."a"`), sql)
}

// TestLoadSDLInlinePrimaryKeyName keeps an inline primary key's CONSTRAINT
// name, as v25.2.3 does.
func TestLoadSDLInlinePrimaryKeyName(t *testing.T) {
	loaded, err := catalog.LoadSDL("CREATE TABLE t (id INT8 CONSTRAINT custom_pk PRIMARY KEY); COMMENT ON CONSTRAINT custom_pk ON t IS 'key';")
	require.NoError(t, err)
	require.Equal(t, "custom_pk", loaded.Schemas[0].Tables[0].Indexes[0].Name)
}

// TestSDLRoutineUserTypes reads a routine's unqualified user-defined types as
// the ones in public, and skips a CREATE TABLE IF NOT EXISTS of a table the
// SDL created.
func TestSDLRoutineUserTypes(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TYPE public.mood AS ENUM ('ok'); CREATE FUNCTION public.f(x public.mood) RETURNS public.mood LANGUAGE SQL AS $$ SELECT x $$; CREATE TABLE t (id INT8 PRIMARY KEY);",
		"CREATE TYPE public.mood AS ENUM ('ok'); CREATE FUNCTION f(x mood) RETURNS mood LANGUAGE SQL AS $$ SELECT x $$; CREATE TABLE t (id INT8 PRIMARY KEY); CREATE TABLE IF NOT EXISTS t (id INT8 PRIMARY KEY, extra STRING);",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}

// TestSDLIfNotExists skips a CREATE ... IF NOT EXISTS of an object the SDL
// created, as the engine does.
func TestSDLIfNotExists(t *testing.T) {
	source := "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, b INT8); CREATE INDEX i ON t (a); CREATE SEQUENCE s; CREATE TYPE e AS ENUM ('a'); CREATE VIEW v AS SELECT 1 AS n;"
	plan, err := catalog.BuildSDLMigrationPlan(source, source+
		"ALTER TABLE t ADD COLUMN IF NOT EXISTS a STRING NOT NULL DEFAULT ''; CREATE INDEX IF NOT EXISTS i ON t (b); CREATE SEQUENCE IF NOT EXISTS s INCREMENT 5; CREATE TYPE IF NOT EXISTS e AS ENUM ('b'); CREATE VIEW IF NOT EXISTS v AS SELECT 2 AS n;",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}

// TestLoadSDLOneFamily loads a table that names its one family.
func TestLoadSDLOneFamily(t *testing.T) {
	_, err := catalog.LoadSDL("CREATE TABLE t (id INT8 PRIMARY KEY, a STRING, FAMILY f1 (id, a))")
	require.NoError(t, err)
}

// TestSDLNestedCTEScope keeps a nested WITH's names inside the subquery that
// holds it.
func TestSDLNestedCTEScope(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan("", `CREATE TABLE public.t (id INT8 PRIMARY KEY);
CREATE VIEW public.a AS SELECT z.id FROM z CROSS JOIN (WITH z AS (SELECT 1 AS n) SELECT n FROM z) AS q;
CREATE VIEW public.z AS SELECT id FROM public.t;`, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	sql := plan.SQL()
	require.Less(t, strings.Index(sql, `CREATE VIEW "public"."z"`), strings.Index(sql, `CREATE VIEW "public"."a"`), sql)

	// The same query, written with the stored view qualified, is the same view.
	plan, err = catalog.BuildSDLMigrationPlan(`CREATE TABLE public.t (id INT8 PRIMARY KEY);
CREATE VIEW public.z AS SELECT id FROM public.t;
CREATE VIEW public.a AS SELECT z.id FROM public.z CROSS JOIN (WITH z AS (SELECT 1 AS n) SELECT n FROM z) AS q;`, `CREATE TABLE public.t (id INT8 PRIMARY KEY);
CREATE VIEW public.z AS SELECT id FROM public.t;
CREATE VIEW public.a AS SELECT z.id FROM z CROSS JOIN (WITH z AS (SELECT 1 AS n) SELECT n FROM z) AS q;`, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "%v", plan)
}
