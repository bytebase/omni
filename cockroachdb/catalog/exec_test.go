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
