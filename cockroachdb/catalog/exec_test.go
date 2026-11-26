package catalog_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
)

// execLast runs a script and returns the refusal of its last statement,
// failing on any refusal before it.
func execLast(t *testing.T, script string) (*catalog.Catalog, error) {
	t.Helper()
	c := catalog.New()
	results, err := c.Exec(script)
	require.NoError(t, err)
	for i, r := range results {
		if r.Error != nil && i < len(results)-1 {
			t.Fatalf("line %d: %s: %v", r.Line, r.SQL, r.Error)
		}
	}
	return c, results[len(results)-1].Error
}

func code(err error) string {
	var refusal *catalog.Error
	if errors.As(err, &refusal) {
		return refusal.Code
	}
	return ""
}

func TestExecRefusals(t *testing.T) {
	for _, tc := range []struct{ script, code string }{
		{"ALTER TABLE missing.t ADD COLUMN x INT8", "42P01"},
		{"COMMENT ON TABLE missing.t IS 'x'", "42P01"},
		{"CREATE TABLE t (a INT8); CREATE INDEX i ON t (a) STORING (missing)", "42703"},
		{"CREATE TABLE t (a INT8); CREATE INDEX i ON t (a) WHERE missing > 0", "42703"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (id STRING REFERENCES p (id))", "42804"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (s STRING); ALTER TABLE c ADD FOREIGN KEY (s) REFERENCES p", "42804"},
		{"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t DROP COLUMN a", "42P10"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY, u INT8); CREATE UNIQUE INDEX i ON p (u); CREATE TABLE c (u INT8 REFERENCES p (u)); DROP INDEX p@i", "2BP01"},
		{"CREATE TABLE x (a INT8, b INT8); CREATE TABLE y (a INT8); CREATE VIEW v AS SELECT a, b FROM x; CREATE OR REPLACE VIEW v AS SELECT * FROM y", "42P16"},
		{"CREATE SEQUENCE sq; ALTER SEQUENCE sq MINVALUE 10 MAXVALUE 5", "22023"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (a INT8, b INT8, FOREIGN KEY (a, b) REFERENCES p (id))", "42601"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (a INT8, b INT8, FOREIGN KEY (a, b) REFERENCES p)", "42601"},
		{"CREATE TABLE t (a INT8, b INT8 AS (missing * 2) STORED)", "42703"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT a FROM t WHERE missing > 0", "42703"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (id INT8 REFERENCES p (id)); ALTER TABLE c ALTER COLUMN id TYPE STRING", "0A000"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (id INT8 REFERENCES p (id)); ALTER TABLE p ALTER COLUMN id TYPE STRING", "0A000"},
		{"CREATE FUNCTION missing.f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$", "3F000"},
		{"CREATE TABLE t (a INT8 PRIMARY KEY, b INT8 PRIMARY KEY)", "XXUUU"},
		{"CREATE TABLE t (a INT8 PRIMARY KEY); ALTER TABLE t ADD COLUMN b INT8 PRIMARY KEY", "42611"},
		{"CREATE TABLE t (a INT8, CONSTRAINT ck CHECK (a > 0), CONSTRAINT ck CHECK (a < 10))", "42710"},
		{"CREATE TABLE t (a INT8 CONSTRAINT ck CHECK (a > 0), b INT8 CONSTRAINT ck UNIQUE)", "42710"},
		{"CREATE TABLE t (a INT8, CONSTRAINT ck UNIQUE (a), CONSTRAINT ck UNIQUE (a))", "42P07"},
		{"CREATE TABLE t (a INT8); CREATE TABLE u (a INT8 DEFAULT nextval('t'))", "42809"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (a INT8); ALTER TABLE c ADD COLUMN b STRING REFERENCES p (id)", "42804"},
		{"CREATE TABLE t (a INT8 CHECK (a > 0)); ALTER TABLE t ALTER COLUMN a TYPE STRING USING a::STRING", "0A000"},
		{"CREATE TABLE t (a INT8, INDEX (a)); ALTER TABLE t ALTER COLUMN a TYPE STRING USING a::STRING", "0A000"},
		{"CREATE TYPE mood AS ENUM ('sad'); CREATE FUNCTION f(x mood) RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; DROP TYPE mood", "2BP01"},
		{"CREATE SEQUENCE sq INCREMENT BY 0", "22023"},
		{"CREATE SEQUENCE sq; ALTER SEQUENCE sq INCREMENT BY 0", "22023"},
		{"CREATE TABLE t (a INT8, b INT8 DEFAULT a)", "42601"},
		{"CREATE TABLE t (a INT8); ALTER TABLE t ADD COLUMN b INT8 DEFAULT a", "42601"},
		{"CREATE TABLE t (a INT8); ALTER TABLE t ALTER COLUMN a SET DEFAULT a", "42601"},
		{"CREATE TABLE t (id INT8 PRIMARY KEY, parent INT8 REFERENCES t (missing))", "42703"},
		{"CREATE TABLE t (id INT8 PRIMARY KEY, u INT8, parent INT8 REFERENCES t (u))", "23503"},
		{"CREATE TABLE p (id INT8 PRIMARY KEY, u INT8); CREATE UNIQUE INDEX ux ON p (u); ALTER TABLE p ADD COLUMN parent INT8 REFERENCES p (u); DROP INDEX p@ux", "2BP01"},
		{"CREATE SCHEMA s; CREATE TYPE s.t AS ENUM ('a'); CREATE TABLE public.t (a INT8); ALTER TABLE public.t SET SCHEMA s", "42P07"},
		{"CREATE SCHEMA s; CREATE TYPE s.m AS ENUM ('a'); CREATE TABLE t (m s.m); DROP SCHEMA s CASCADE", "0A000"},
	} {
		_, err := execLast(t, tc.script)
		require.Equal(t, tc.code, code(err), "%s: %v", tc.script, err)
	}
	_, err := execLast(t, "CREATE SCHEMA s; SET search_path = s")
	require.Error(t, err)
}

func TestExecRenameColumnRewritesComputedColumns(t *testing.T) {
	c, err := execLast(t, "CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t RENAME COLUMN a TO x")
	require.NoError(t, err)
	require.Equal(t, "x * 2", c.Metadata().Schemas[0].Tables[0].Columns[1].Generation.Expression)
}

func TestExecDropIndexCascadeDropsForeignKeys(t *testing.T) {
	c, err := execLast(t, "CREATE TABLE p (id INT8 PRIMARY KEY, u INT8 UNIQUE); CREATE TABLE c (u INT8 REFERENCES p (u)); DROP INDEX p@p_u_key CASCADE")
	require.NoError(t, err)
	for _, table := range c.Metadata().Schemas[0].Tables {
		require.Empty(t, table.ForeignKeys, table.Name)
	}
}

// TestSDLComputedAnnotations keeps an annotation that picks another overload
// and drops one CockroachDB writes on a literal's own type.
func TestSDLComputedAnnotations(t *testing.T) {
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 AS (length('é':::BYTES)) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 AS (length('é':::STRING)) STORED);",
		catalog.SDLMigrationOptions{})
	require.Error(t, err, "%v", plan)
	plan, err = catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100:::DECIMAL) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100) STORED);",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)
	_, err = catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL AS ((9007199254740993:::FLOAT8)::DECIMAL) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, a DECIMAL AS ((9007199254740993:::DECIMAL)::DECIMAL) STORED);",
		catalog.SDLMigrationOptions{})
	require.Error(t, err)
}

func TestExecRenamesFollowIDs(t *testing.T) {
	c, err := execLast(t, "CREATE SCHEMA s; CREATE TYPE s.m AS ENUM ('a'); CREATE TABLE t (m s.m); ALTER SCHEMA s RENAME TO s2")
	require.NoError(t, err)
	require.Equal(t, "s2.m", c.Metadata().Schemas[0].Tables[0].Columns[0].Type)
	c, err = execLast(t, "CREATE SEQUENCE sq; CREATE TABLE t (v STRING DEFAULT nextval('sq')::STRING || 'public.sq'); ALTER SEQUENCE sq RENAME TO sq2")
	require.NoError(t, err)
	require.Equal(t, "nextval('public.sq2')::STRING || 'public.sq'", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
}

func TestExecDropsDependents(t *testing.T) {
	c, err := execLast(t, "CREATE SCHEMA s; CREATE VIEW s.v AS SELECT 1 AS a; CREATE VIEW public.w AS SELECT a FROM s.v; DROP SCHEMA s CASCADE")
	require.NoError(t, err)
	require.Len(t, c.Metadata().Schemas, 1)
	require.Empty(t, c.Metadata().Schemas[0].Views)
	c, err = execLast(t, "CREATE TABLE t (id INT8 PRIMARY KEY, a INT8); CREATE SEQUENCE sq OWNED BY t.a; ALTER TABLE t DROP COLUMN a")
	require.NoError(t, err)
	require.Empty(t, c.Metadata().Schemas[0].Sequences)
}

func TestExecRetypesDefaultsAndDropsOwnedSequences(t *testing.T) {
	c, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (m mood DEFAULT 'sad'::mood); ALTER TYPE mood RENAME TO feeling")
	require.NoError(t, err)
	require.Equal(t, "'sad'::public.feeling", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	c, err = execLast(t, "CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; CREATE TABLE u (a INT8 DEFAULT nextval('sq')); DROP TABLE t CASCADE")
	require.NoError(t, err)
	require.Empty(t, c.Metadata().Schemas[0].Sequences)
	require.Empty(t, c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	_, err = execLast(t, "CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; CREATE TABLE u (a INT8 DEFAULT nextval('sq')); DROP TABLE t")
	require.Equal(t, "2BP01", code(err))
	c, err = execLast(t, "CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; ALTER TABLE t ALTER COLUMN a SET DEFAULT nextval('sq'); DROP TABLE t")
	require.NoError(t, err)
	require.Empty(t, c.Metadata().Schemas[0].Sequences)
}

func TestExecFollowsLabelAndSchemaRenames(t *testing.T) {
	c, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (m mood DEFAULT 'sad'::mood); ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.NoError(t, err)
	require.Equal(t, "'blue'::mood", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	c, err = execLast(t, "CREATE SCHEMA s; CREATE SEQUENCE s.sq; CREATE TABLE t (a INT8 DEFAULT nextval('s.sq')); ALTER SCHEMA s RENAME TO s2")
	require.NoError(t, err)
	require.Equal(t, "nextval('s2.sq')", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	_, err = execLast(t, "CREATE SCHEMA s; CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; ALTER TABLE t SET SCHEMA s")
	require.Error(t, err)
}

func TestExecRelabelsOnlyTheEnum(t *testing.T) {
	c, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad', 'happy'); CREATE TABLE t (m mood DEFAULT CASE WHEN current_user = 'sad' THEN 'sad'::mood ELSE 'happy'::mood END); ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.NoError(t, err)
	require.Equal(t, "CASE WHEN current_user() = 'sad' THEN 'blue'::mood ELSE 'happy'::mood END", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
}

func TestExecLineSkipsComments(t *testing.T) {
	results, err := catalog.New().Exec("CREATE TABLE t (a INT8);\n-- migration\n/* two\nlines */\nALTER TABLE missing ADD COLUMN a INT8;")
	require.NoError(t, err)
	require.Equal(t, 5, results[len(results)-1].Line)
}

// TestExecRefusesUnmodeledForms refuses, with no code, what the metadata
// cannot hold afterwards.
func TestExecRefusesUnmodeledForms(t *testing.T) {
	for _, script := range []string{
		"CREATE SCHEMA s; CREATE TABLE s.t (a INT8); CREATE VIEW public.v AS SELECT a FROM s.t; ALTER SCHEMA s RENAME TO s2",
		"CREATE SCHEMA s; CREATE TABLE public.t (a INT8); CREATE SEQUENCE s.sq OWNED BY public.t.a",
		"CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t ALTER COLUMN a TYPE STRING USING a::STRING",
		"CREATE TABLE p (id INT8 PRIMARY KEY, u INT8, x INT8); CREATE UNIQUE INDEX ux ON p (u) STORING (x); CREATE TABLE c (u INT8 REFERENCES p (u)); ALTER TABLE p DROP COLUMN x",
		"CREATE TYPE mood AS ENUM ('sad'); CREATE FUNCTION f(x mood) RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; ALTER TYPE mood RENAME TO feeling",
		"CREATE SCHEMA s; CREATE FUNCTION s.f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; ALTER SCHEMA s RENAME TO s2",
	} {
		_, err := execLast(t, script)
		require.Error(t, err, script)
		require.Empty(t, code(err), script)
	}
}

func TestExecSelfReference(t *testing.T) {
	_, err := execLast(t, "CREATE TABLE t (id INT8 PRIMARY KEY, parent INT8 REFERENCES t (id))")
	require.NoError(t, err)
	_, err = execLast(t, "CREATE TABLE t (id INT8 PRIMARY KEY, parent INT8 REFERENCES t)")
	require.NoError(t, err)
}

func TestExecRound12(t *testing.T) {
	for _, script := range []string{
		"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (id INT8 REFERENCES p (id)); TRUNCATE p, c",
		"CREATE TABLE a (x INT8); CREATE TABLE b (x INT8); CREATE INDEX i ON a (x); CREATE INDEX i ON b (x)",
		"CREATE SCHEMA s; CREATE TABLE s.t (id INT8 PRIMARY KEY, parent INT8 REFERENCES s.t (id))",
		"CREATE SEQUENCE sq; ALTER SEQUENCE sq AS INT4 NO MAXVALUE",
	} {
		_, err := execLast(t, script)
		require.NoError(t, err, script)
	}
	for _, tc := range []struct{ script, code string }{
		{"CREATE TABLE t (a INT8, INDEX dup (a), INDEX dup (a))", "42P07"},
		{"CREATE TYPE mood AS ENUM ('a'); CREATE VIEW mood AS SELECT 1 AS x", "42P07"},
		{"CREATE TYPE mood AS ENUM ('a'); CREATE SEQUENCE mood", "42P07"},
		{"CREATE TABLE t (a INT8); ALTER SEQUENCE t INCREMENT BY 5", "42809"},
		{"CREATE SCHEMA s; CREATE TABLE t (a INT8); ALTER VIEW t SET SCHEMA s", "42809"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT a FROM t; CREATE OR REPLACE VIEW v AS SELECT * FROM t WHERE missing > 0", "42703"},
		{"CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (a STRING DEFAULT ('sad'::mood)::STRING); DROP TYPE mood", "2BP01"},
		{"CREATE SCHEMA s; CREATE TYPE s.m AS ENUM ('a'); CREATE TABLE s.t (x INT8); CREATE TABLE public.t (m s.m); DROP SCHEMA s CASCADE", "0A000"},
		{"CREATE TABLE t (a INT8, b INT8); CREATE VIEW v AS SELECT x FROM t AS u (x, y); ALTER TABLE t DROP COLUMN a", "2BP01"},
		{"ALTER TABLE missing OWNER TO root", "42P01"},
		{"CREATE TABLE t (a INT8, INDEX i (a) STORING (missing))", "42703"},
		{"CREATE TABLE t (a INT8, UNIQUE (a) WHERE missing > 0)", "42703"},
		{"CREATE TABLE t (a INT8, CONSTRAINT ck CHECK (a > 0)); ALTER TABLE t ADD COLUMN b INT8 CONSTRAINT ck CHECK (b > 0)", "42710"},
		{"CREATE TABLE t (a INT8 DEFAULT nextval(('missing')))", "42P01"},
	} {
		_, err := execLast(t, tc.script)
		require.Equal(t, tc.code, code(err), "%s: %v", tc.script, err)
	}
	for _, script := range []string{
		"SELECT set_config('search_path', 's', false)",
		"CREATE SCHEMA s; CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; ALTER SEQUENCE sq SET SCHEMA s",
		"CREATE SCHEMA s; CREATE TABLE t (a INT8); CREATE SEQUENCE sq OWNED BY t.a; ALTER TABLE public.t RENAME TO s.t2",
		"ALTER SCHEMA public RENAME TO p2",
		"CREATE TABLE t (a INT8); ALTER TABLE t ADD COLUMN b INT8 CREATE FAMILY f2",
	} {
		_, err := execLast(t, script)
		require.Error(t, err, script)
		require.Empty(t, code(err), script)
	}
	results, err := catalog.New().Exec("/* a /* nested */ comment */\nALTER TABLE missing ADD COLUMN a INT8;")
	require.NoError(t, err)
	require.Equal(t, 2, results[0].Line)
	c, err := execLast(t, `CREATE SCHEMA "Mixed Case"; CREATE SEQUENCE "Mixed Case"."My Seq"; CREATE TABLE t (a INT8 DEFAULT nextval('"Mixed Case"."My Seq"')); ALTER SEQUENCE "Mixed Case"."My Seq" RENAME TO "Mixed Case"."Other Seq"`)
	require.NoError(t, err)
	require.Equal(t, `nextval('"Mixed Case"."Other Seq"')`, c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	c, err = execLast(t, "CREATE TABLE t (a STRING); CREATE INDEX i ON t (lower(a)); ALTER TABLE t RENAME COLUMN a TO b")
	require.NoError(t, err)
	require.Equal(t, "lower(b)", c.Metadata().Schemas[0].Tables[0].Indexes[0].Expressions[0])
}

// TestExecReviewRound13 holds the codes the pinned server answered for the
// review's thirteenth round.
func TestExecReviewRound13(t *testing.T) {
	for _, tc := range []struct{ script, code string }{
		{"CREATE FUNCTION f(x missing_type) RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$", "42704"},
		{"CREATE FUNCTION f() RETURNS missing_type LANGUAGE SQL AS $$ SELECT 1 $$", "42704"},
		{"CREATE FUNCTION f(x INT8 DEFAULT nextval('missing')) RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$", "42P01"},
		{"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$", "42723"},
		{"CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; CREATE OR REPLACE FUNCTION f() RETURNS STRING LANGUAGE SQL AS $$ SELECT 'a' $$", "42P13"},
		{"CREATE PROCEDURE f() LANGUAGE SQL AS $$ SELECT 1 $$; CREATE OR REPLACE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$", "42809"},
		{"CREATE TABLE t (a INT8 AS (a + 1) STORED)", "42P16"},
		{"CREATE TABLE t (a INT8, b INT8 AS (a + 1) STORED, c INT8 AS (b + 1) STORED)", "42P16"},
		{"CREATE TABLE t (a INT8, b INT8 DEFAULT 1 AS (a + 1) STORED)", "42P16"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT a FROM t; CREATE OR REPLACE VIEW v AS SELECT a FROM v", "42P17"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v (x, y) AS SELECT a FROM t", "42601"},
		{"CREATE VIEW v (x, x) AS SELECT 1, 2", "42701"},
		{"CREATE VIEW v AS SELECT 1 AS x, 2 AS x", "42701"},
		{"CREATE MATERIALIZED VIEW mv (a, rowid) AS SELECT 1", "42601"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT a FROM t; CREATE VIEW w AS SELECT missing FROM v", "42703"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT a FROM t ORDER BY missing", "42703"},
		{"CREATE TABLE t (a INT8); CREATE VIEW v AS SELECT t.missing FROM t", "42703"},
		{"CREATE TABLE t (a INT8, CHECK (t.missing > 0))", "42703"},
		{"CREATE SEQUENCE sq CACHE 0", "22023"},
		{"CREATE SEQUENCE sq; ALTER SEQUENCE sq CACHE -1", "22023"},
		{"CREATE SEQUENCE sq AS INT2 MAXVALUE 100000", "22023"},
		{"CREATE SEQUENCE sq MAXVALUE 100000; ALTER SEQUENCE sq AS INT2", "22023"},
		{"CREATE SEQUENCE sq MAXVALUE 10; ALTER SEQUENCE sq RESTART WITH 50", "22023"},
		{"CREATE TABLE a (x INT8); CREATE TABLE b (x INT8); CREATE INDEX i ON a (x); CREATE INDEX i ON b (x); DROP INDEX i", "42P08"},
		{"CREATE TYPE pair AS (x INT8, x STRING)", "42P17"},
		{"CREATE TYPE m1 AS ENUM ('a'); CREATE TABLE p (id m1 PRIMARY KEY); CREATE TABLE c (x INT8 REFERENCES p (id))", "42804"},
		{"CREATE TABLE p (id INT8[] PRIMARY KEY); CREATE TABLE c (x STRING[] REFERENCES p (id))", "42804"},
		{"CREATE TYPE m1 AS ENUM ('a'); CREATE TYPE m2 AS ENUM ('a'); CREATE TABLE t (a m1 PRIMARY KEY); ALTER TABLE t ALTER COLUMN a TYPE m2 USING a::STRING::m2", "0A000"},
		{"CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t ALTER COLUMN b TYPE STRING", "XXUUU"},
		{"CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t ALTER COLUMN b SET DEFAULT 5", "42601"},
		{"CREATE TABLE t (a INT8, b INT8 AS (a * 2) STORED); ALTER TABLE t ALTER COLUMN b DROP DEFAULT", "42601"},
		{"CREATE TABLE t (a INT8); ALTER TABLE t ALTER COLUMN a TYPE INT8 COLLATE \"en_US\"", "42601"},
		{"CREATE TABLE t (a STRING, INDEX (a)); ALTER TABLE t ALTER COLUMN a TYPE STRING COLLATE \"en_US\"", "0A000"},
		{"CREATE TABLE t (b INT8); ALTER TABLE t ALTER COLUMN b TYPE STRING USING missing::STRING", "42703"},
		{"CREATE TABLE t (a INT8, b INT8 UNIQUE, INDEX t_a_idx (a)); ALTER TABLE t RENAME CONSTRAINT t_b_key TO t_a_idx", "42P07"},
		{"CREATE TABLE t (a INT8, b INT8); CREATE UNIQUE INDEX ux ON t (b); ALTER TABLE t VALIDATE CONSTRAINT ux", "42809"},
		{"CREATE TABLE t (a INT8); CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) FROM t $$; DROP TABLE t", "2BP01"},
		{"CREATE TABLE t (a INT8); CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) FROM t $$; ALTER TABLE t RENAME TO t2", "2BP01"},
		{"CREATE SCHEMA s; CREATE TABLE s.t (a INT8); CREATE FUNCTION public.f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) FROM s.t $$; ALTER SCHEMA s RENAME TO s2", "2BP01"},
	} {
		_, err := execLast(t, tc.script)
		require.Equal(t, tc.code, code(err), "%s: %v", tc.script, err)
	}
	for _, script := range []string{
		"CREATE EXTENSION fuzzystrmatch; CREATE EXTENSION fuzzystrmatch",
		"CREATE TABLE p (id INT8 PRIMARY KEY); CREATE TABLE c (x INT4 REFERENCES p (id))",
		"CREATE TABLE t (a INT8, b INT8); CREATE UNIQUE INDEX ux ON t (b); ALTER TABLE t RENAME CONSTRAINT ux TO uy",
	} {
		_, err := execLast(t, script)
		require.NoError(t, err, script)
	}
	c, err := execLast(t, "CREATE TABLE t (a INT8 CONSTRAINT k PRIMARY KEY UNIQUE)")
	require.NoError(t, err)
	require.Len(t, c.Metadata().Schemas[0].Tables[0].Indexes, 1)
	c, err = execLast(t, "CREATE TABLE t (a INT8 PRIMARY KEY, b INT8 NOT NULL, c INT8); CREATE INDEX t_a_key ON t (c); ALTER TABLE t ALTER PRIMARY KEY USING COLUMNS (b)")
	require.NoError(t, err)
	var names []string
	for _, index := range c.Metadata().Schemas[0].Tables[0].Indexes {
		names = append(names, index.Name)
	}
	require.ElementsMatch(t, []string{"t_pkey", "t_a_key", "t_a_key1"}, names)
	c, err = execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (m mood, CHECK (m = 'sad'::mood)); ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.NoError(t, err)
	require.Equal(t, "((m = 'blue'::mood))", c.Metadata().Schemas[0].Tables[0].CheckConstraints[0].Expression)
	c, err = execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (m mood[] DEFAULT ARRAY['sad']::mood[]); ALTER TYPE mood RENAME TO feeling")
	require.NoError(t, err)
	require.Equal(t, "ARRAY['sad']::public.feeling[]", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	_, err = execLast(t, "CREATE TYPE mood AS ENUM ('sad', 'x'); CREATE TABLE t (m mood DEFAULT COALESCE(NULL, 'sad')); ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.Error(t, err)
	_, err = execLast(t, `CREATE TYPE "Mood" AS ENUM ('a'); CREATE FUNCTION f(x "Mood") RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$; DROP TYPE "Mood"`)
	require.Equal(t, "2BP01", code(err))
}

func TestSDLComputedAnnotationOverloads(t *testing.T) {
	_, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING AS (pg_typeof(1:::INT8)::STRING) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, a STRING AS (pg_typeof(1:::DECIMAL)::STRING) STORED);",
		catalog.SDLMigrationOptions{})
	require.Error(t, err)
}

func TestExecRound14(t *testing.T) {
	c, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (a STRING DEFAULT ('sad'::mood)::STRING); ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.NoError(t, err)
	require.Equal(t, "('blue'::mood)::STRING", c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
	c, err = execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE TABLE t (a STRING CHECK (a::mood = 'sad'::mood)); ALTER TYPE mood RENAME TO feeling")
	require.NoError(t, err)
	require.Contains(t, c.Metadata().Schemas[0].Tables[0].CheckConstraints[0].Expression, "public.feeling")
	_, err = execLast(t, `CREATE SCHEMA "a.b"; CREATE SEQUENCE "a.b".sq; CREATE TABLE t (id INT8 DEFAULT nextval('"a.b".sq'))`)
	require.NoError(t, err)
}

func TestExecRound15(t *testing.T) {
	_, err := execLast(t, "CREATE TABLE t (a INT8); CREATE FUNCTION f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT a FROM t LIMIT 1 $$; ALTER TABLE t RENAME COLUMN a TO b")
	require.Error(t, err)
	_, err = execLast(t, "CREATE TABLE t (a INT8, b INT8); CREATE TABLE u (a INT8); CREATE VIEW v AS SELECT t.b FROM t JOIN u USING (a); ALTER TABLE t DROP COLUMN a")
	require.Equal(t, "2BP01", code(err))
}

func TestExecRound16(t *testing.T) {
	_, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE VIEW v AS SELECT 'sad'::mood AS m; DROP TYPE mood")
	require.Equal(t, "2BP01", code(err))
	_, err = execLast(t, "CREATE TABLE t (x INT8); ALTER TABLE t ADD COLUMN a INT8 AS (a + 1) STORED")
	require.Equal(t, "42P16", code(err))
	_, err = execLast(t, "CREATE VIEW v AS SELECT 1 AS a; CREATE VIEW w AS SELECT a FROM v; CREATE OR REPLACE VIEW v AS SELECT a FROM w")
	require.Equal(t, "42P17", code(err))
}

func TestExecRound17(t *testing.T) {
	_, err := execLast(t, "CREATE TYPE mood AS ENUM ('sad'); CREATE VIEW v AS SELECT 'sad'::mood AS m; ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'")
	require.Error(t, err)
	c, err := execLast(t, `CREATE SCHEMA "a.b"; CREATE SEQUENCE "a.b".sq; CREATE TABLE t (id INT8 DEFAULT nextval('"a.b".sq')); ALTER SEQUENCE "a.b".sq RENAME TO sq2`)
	require.NoError(t, err)
	require.Equal(t, `nextval('"a.b".sq2')`, c.Metadata().Schemas[0].Tables[0].Columns[0].Default)
}

func TestSDLComputedAnnotationContext(t *testing.T) {
	_, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, x INT8, a STRING AS (pg_typeof(x + 1:::DECIMAL)::STRING) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, x INT8, a STRING AS (pg_typeof(x + 1)::STRING) STORED);",
		catalog.SDLMigrationOptions{})
	require.Error(t, err)
	plan, err := catalog.BuildSDLMigrationPlan(
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100:::DECIMAL) STORED);",
		"CREATE TABLE t (id INT8 PRIMARY KEY, p DECIMAL, a DECIMAL AS (p * 100) STORED);",
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)
}

func TestExecRound18(t *testing.T) {
	for _, script := range []string{
		"CREATE SEQUENCE sq; CREATE VIEW v AS SELECT nextval('sq') AS id; ALTER SEQUENCE sq RENAME TO sq2",
		"CREATE SEQUENCE sq; CREATE VIEW v AS SELECT nextval('sq') AS id; DROP SEQUENCE sq",
		"CREATE TYPE mood AS ENUM ('sad', 'happy'); CREATE VIEW v AS SELECT 'sad'::mood AS m; ALTER TYPE mood DROP VALUE 'sad'",
	} {
		_, err := execLast(t, script)
		require.Error(t, err, script)
	}
}

func TestExecRound19(t *testing.T) {
	_, err := execLast(t, "CREATE SCHEMA s; CREATE SEQUENCE s.sq; CREATE VIEW v AS SELECT nextval('s.sq') AS id; ALTER SCHEMA s RENAME TO s2")
	require.Error(t, err)
	_, err = execLast(t, "CREATE TABLE t (a INT8); ALTER TABLE t ADD COLUMN b INT8 CONSTRAINT ck CHECK (b > 0) CONSTRAINT ck CHECK (b < 10)")
	require.Equal(t, "42710", code(err))
	plan, err := catalog.BuildSDLMigrationPlan(
		`CREATE SCHEMA "a.b"; CREATE SCHEMA a; CREATE SEQUENCE "a.b".c; CREATE SEQUENCE a."b.c"; CREATE TABLE t (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('"a.b".c'));`,
		`CREATE SCHEMA "a.b"; CREATE SCHEMA a; CREATE SEQUENCE "a.b".c; CREATE SEQUENCE a."b.c"; CREATE TABLE t (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('a."b.c"'));`,
		catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.NotNil(t, plan)
}
