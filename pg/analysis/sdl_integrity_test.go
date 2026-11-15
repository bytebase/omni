package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func checkOneSDLFile(text string) []SDLFinding {
	return CheckSDLIntegrity([]SDLFile{{Path: "f.sql", Text: text}})["f.sql"]
}

func TestCheckSDLIntegrityFile(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []SDLFinding
	}{
		{
			name: "valid foreign key and view",
			sql: "CREATE TABLE public.users (id int, CONSTRAINT pk_users PRIMARY KEY (id));\n" +
				"CREATE TABLE orders (uid int, CONSTRAINT fk FOREIGN KEY (uid) REFERENCES users(id));\n" +
				"CREATE VIEW public.v AS WITH c AS (SELECT 1) SELECT * FROM c, users, pg_catalog.pg_class, pg_tables;",
		},
		{
			name: "duplicates within a file, in statement order; the duplicate table is dropped whole",
			sql: "CREATE TABLE t (a int, a text, CONSTRAINT c CHECK (a > 0), CONSTRAINT c CHECK (a < 9), CONSTRAINT k UNIQUE (a));\n" +
				"CREATE TABLE public.t (zz int, zz int);\n" +
				"CREATE TABLE u (a int, CONSTRAINT k PRIMARY KEY (a));\n" +
				"CREATE INDEX i ON t (a);\nCREATE INDEX i ON public.u (a);\nCREATE INDEX ON t (a);\nCREATE INDEX ON u (a);",
			want: []SDLFinding{
				{Code: SDLDuplicateColumn, Line: 1, Schema: "public", Table: "t", Column: "a"},
				{Code: SDLDuplicateConstraint, Line: 1, Schema: "public", Table: "t", Name: "c", ConstraintType: "CHECK", First: &SDLDefinition{Table: "t", ConstraintType: "CHECK", Line: 1}},
				{Code: SDLDuplicateTable, Line: 2, Schema: "public", Table: "t", First: &SDLDefinition{Line: 1}},
				{Code: SDLDuplicateConstraint, Line: 3, Schema: "public", Table: "u", Name: "k", ConstraintType: "PRIMARY KEY", First: &SDLDefinition{Table: "t", ConstraintType: "UNIQUE", Line: 1}},
				{Code: SDLDuplicateIndex, Line: 5, Schema: "public", Name: "i", First: &SDLDefinition{Line: 4}},
			},
		},
		{
			name: "multiple primary keys, reported at the second; column-level keys are not counted",
			sql:  "CREATE TABLE s.t (\n  a int PRIMARY KEY,\n  PRIMARY KEY (a),\n  CONSTRAINT p2 PRIMARY KEY (a)\n);",
			want: []SDLFinding{{
				Code: SDLMultiplePrimaryKeys, Line: 4, Schema: "s", Table: "t",
				PrimaryKeys: []SDLDefinition{{Line: 3}, {Name: "p2", Line: 4}},
			}},
		},
		{
			name: "foreign keys",
			sql: "CREATE TABLE s.u (x int, y int);\n" +
				"CREATE TABLE s.t (\n  a int,\n  b int,\n" +
				"  CONSTRAINT f1 FOREIGN KEY (a) REFERENCES missing(x),\n" +
				"  FOREIGN KEY (a, b) REFERENCES s.u(x, zz),\n" +
				"  FOREIGN KEY (a, b) REFERENCES s.u(zz),\n" +
				"  FOREIGN KEY (nope) REFERENCES s.u(zz)\n);",
			want: []SDLFinding{
				{Code: SDLForeignKeyTableNotFound, Line: 5, Schema: "s", Table: "t", Name: "f1", RefSchema: "public", RefTable: "missing", Columns: []string{"a"}},
				{Code: SDLForeignKeyColumnNotFound, Line: 6, Schema: "s", Table: "t", Column: "b", RefSchema: "s", RefTable: "u", RefColumn: "zz", Columns: []string{"x", "y"}},
			},
		},
		{
			name: "view dependencies resolve against the view's schema; nested CTEs are not excluded",
			sql: "CREATE TABLE s.t (a int);\n" +
				"CREATE VIEW s.v AS SELECT * FROM t, public.t, (WITH d AS (SELECT 1) SELECT * FROM d) x;\n" +
				"CREATE VIEW w AS SELECT * FROM v3;\nCREATE VIEW w AS SELECT * FROM s.v;",
			want: []SDLFinding{
				{Code: SDLViewDependencyNotFound, Line: 2, Schema: "s", Name: "v", RefSchema: "public", RefTable: "t"},
				{Code: SDLViewDependencyNotFound, Line: 2, Schema: "s", Name: "v", RefSchema: "s", RefTable: "d"},
			},
		},
		{
			name: "system schemas, and system table and view names in any schema, are skipped",
			sql:  "CREATE VIEW v AS SELECT * FROM s.pg_class, s.pg_tables, pg_temp_3.x, pg_toast.y, information_schema.tables;",
		},
		{name: "empty", sql: "", want: nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkOneSDLFile(tc.sql)
			if len(tc.want) == 0 {
				require.Empty(t, got)
				return
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestCheckSDLIntegrityFiles(t *testing.T) {
	files := []SDLFile{
		{Path: "a.sql", Text: "CREATE TABLE s.t (a int, CONSTRAINT pk PRIMARY KEY (a));\nCREATE INDEX i ON s.t (a);\nCREATE VIEW s.v AS SELECT * FROM s.t;"},
		{Path: "b.sql", Text: "\nCREATE TABLE s.t (b int, CONSTRAINT pk UNIQUE (b), CONSTRAINT c CHECK (b > 0));\nCREATE INDEX i ON s.t (b);\nCREATE VIEW s.v AS SELECT * FROM s.t;"},
		{Path: "c.sql", Text: "CREATE TABLE s.r (x int, CONSTRAINT fk FOREIGN KEY (x) REFERENCES s.t(b), CONSTRAINT c CHECK (x > 0));"},
	}
	first := &SDLDefinition{Path: "a.sql", Line: 1}
	require.Equal(t, map[string][]SDLFinding{
		"a.sql": {},
		"b.sql": {
			{Code: SDLDuplicateTableAcrossFiles, Line: 2, Schema: "s", Table: "t", First: first},
			{Code: SDLDuplicateIndexAcrossFiles, Line: 3, Schema: "s", Name: "i", First: &SDLDefinition{Path: "a.sql", Line: 2}},
			{Code: SDLDuplicateViewAcrossFiles, Line: 4, Schema: "s", Name: "v", First: &SDLDefinition{Path: "a.sql", Line: 3}},
			{Code: SDLDuplicateConstraintAcrossFiles, Line: 2, Schema: "s", Table: "t", Name: "pk", ConstraintType: "UNIQUE", First: &SDLDefinition{Path: "a.sql", Table: "t", ConstraintType: "PRIMARY KEY", Line: 1}},
		},
		// The foreign key resolves against a.sql's s.t, the first definition.
		"c.sql": {
			{Code: SDLForeignKeyColumnNotFound, Line: 1, Schema: "s", Table: "r", Name: "fk", Column: "x", RefSchema: "s", RefTable: "t", RefColumn: "b", Columns: []string{"a"}},
		},
	}, CheckSDLIntegrity(files))

	// A syntax error returns only the first failing file's finding.
	got := CheckSDLIntegrity([]SDLFile{
		{Path: "ok.sql", Text: "CREATE TABLE t (a int);"},
		{Path: "bad1.sql", Text: "CREATE TABLE t (a int;"},
		{Path: "bad2.sql", Text: "CREATE TABLEX t (a int);"},
	})
	require.Equal(t, map[string][]SDLFinding{
		"bad1.sql": {{Code: SDLSyntaxError, Detail: "syntax error: unmatched opening parenthesis"}},
	}, got)

	require.Empty(t, CheckSDLIntegrity(nil))
}
