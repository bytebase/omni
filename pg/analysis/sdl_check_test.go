package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSDLStyle(t *testing.T) {
	// node and stmt name the text a finding's Start and StatementStart point
	// at (its first occurrence); an empty node means Start is -1.
	type want struct {
		SDLFinding
		node, stmt string
	}
	tests := []struct {
		name string
		sql  string
		want []want
	}{
		{
			name: "conforming table",
			sql: `CREATE TABLE public.users (
				id INTEGER NOT NULL DEFAULT 0,
				email TEXT,
				email_lower TEXT GENERATED ALWAYS AS (LOWER(email)) STORED,
				CONSTRAINT pk_users PRIMARY KEY (id),
				CONSTRAINT uk_users_email UNIQUE (email)
			);`,
		},
		{
			name: "a statement starts at its leading comment",
			sql:  "SELECT 1;\n-- users\nCREATE TABLE users (id int);",
			want: []want{{SDLFinding{Code: SDLTableWithoutSchema, Table: "users"}, "", "-- users"}},
		},
		{
			name: "column constraints: references finding first, column finding at its first disallowed constraint",
			sql:  "CREATE TABLE public.t (\n  a int NOT NULL\n    CHECK (a > 0)\n    REFERENCES u(id)\n);",
			want: []want{
				{SDLFinding{Code: SDLColumnForeignKeyWithoutSchema, Table: "t", Column: "a", RefTable: "u"}, "REFERENCES", "CREATE"},
				{SDLFinding{Code: SDLColumnConstraint, Table: "t", Column: "a", ConstraintTypes: []string{"CHECK", "FOREIGN KEY"}}, "CHECK", "CREATE"},
			},
		},
		{
			name: "column references with a schema is still a column constraint",
			sql:  "CREATE TABLE public.t (a int REFERENCES public.u(id));",
			want: []want{{SDLFinding{Code: SDLColumnConstraint, Table: "t", Column: "a", ConstraintTypes: []string{"FOREIGN KEY"}}, "REFERENCES", "CREATE"}},
		},
		{
			name: "unnamed table constraints, and a foreign key without schema",
			sql:  "CREATE TABLE public.t (\n  a int,\n  PRIMARY KEY (a),\n  FOREIGN KEY (a) REFERENCES u(id),\n  EXCLUDE USING gist (a WITH =)\n);",
			want: []want{
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "PRIMARY KEY"}, "PRIMARY", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "FOREIGN KEY"}, "FOREIGN", "CREATE"},
				{SDLFinding{Code: SDLForeignKeyWithoutSchema, Table: "t", RefTable: "u"}, "FOREIGN", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "CONSTRAINT"}, "EXCLUDE", "CREATE"},
			},
		},
		{
			name: "index without name and schema",
			sql:  "CREATE INDEX ON t (a);\nCREATE INDEX i ON t (a);\nCREATE INDEX i2 ON public.t (a);",
			want: []want{
				{SDLFinding{Code: SDLUnnamedIndex}, "", "CREATE INDEX ON"},
				{SDLFinding{Code: SDLIndexWithoutSchema, Table: "t"}, "", "CREATE INDEX ON"},
				{SDLFinding{Code: SDLIndexWithoutSchema, Name: "i", Table: "t"}, "", "CREATE INDEX i "},
			},
		},
		{
			name: "views, sequences and functions without schema",
			sql: "CREATE VIEW v AS SELECT 1;\nCREATE SEQUENCE seq;\nALTER SEQUENCE seq RESTART;\n" +
				"CREATE FUNCTION f() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;\n" +
				"CREATE VIEW public.v2 AS SELECT 1;\nCREATE SEQUENCE public.s2;\nALTER SEQUENCE public.s2 RESTART;\n" +
				"CREATE FUNCTION public.f2() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;",
			want: []want{
				{SDLFinding{Code: SDLViewWithoutSchema, Name: "v"}, "", "CREATE VIEW v "},
				{SDLFinding{Code: SDLSequenceWithoutSchema, Name: "seq"}, "", "CREATE SEQUENCE seq"},
				{SDLFinding{Code: SDLAlterSequenceWithoutSchema, Name: "seq"}, "", "ALTER SEQUENCE seq"},
				{SDLFinding{Code: SDLFunctionWithoutSchema, Name: "f"}, "", "CREATE FUNCTION f("},
			},
		},
		{
			name: "a three-part function name takes its catalog as the schema",
			sql:  "CREATE FUNCTION c.s.f() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;",
		},
		{
			name: "other statements are not judged",
			sql:  "CREATE MATERIALIZED VIEW mv AS SELECT 1;\nDROP TABLE t;\nALTER TABLE t ADD COLUMN x int PRIMARY KEY;",
		},
		{name: "empty", sql: ""},
		{name: "comment only", sql: "-- nothing\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var expected []SDLFinding
			for _, w := range tc.want {
				f := w.SDLFinding
				f.Start = -1
				if w.node != "" {
					f.Start = strings.Index(tc.sql, w.node)
					require.GreaterOrEqual(t, f.Start, 0, w.node)
				}
				f.StatementStart = strings.Index(tc.sql, w.stmt)
				require.GreaterOrEqual(t, f.StatementStart, 0, w.stmt)
				expected = append(expected, f)
			}
			require.Equal(t, expected, CheckSDLStyle(tc.sql))
		})
	}
}

func TestCheckSDLStyleSyntaxError(t *testing.T) {
	tests := []struct {
		sql    string
		detail string
		start  int
	}{
		{"CREATE TABLE public.t (a int;", "syntax error: unmatched opening parenthesis", 22},
		{"CREATE TABLE public.t (a int));", "syntax error: unmatched closing parenthesis", 29},
		// The balance check skips strings, dollar quotes and comments, but
		// ignores backslash escapes and does not nest block comments.
		{"CREATE TABLE public.t (a text DEFAULT E'\\'(');", "syntax error: unmatched opening parenthesis", 42},
		{"/* /* */ ( */ CREATE TABLE public.t (a int);", "syntax error: unmatched opening parenthesis", 9},
		{`CREATE TABLE public."t(" (id int);`, "syntax error: unmatched opening parenthesis", 22},
		{"CREATE TABLE public.t (a text DEFAULT '(', b text DEFAULT $x$)$x$); -- (", "", 0},
		{"CREATE TABLEX public.t (a int);", `syntax error at or near "TABLEX"`, 7},
	}
	for _, tc := range tests {
		got := CheckSDLStyle(tc.sql)
		if tc.detail == "" {
			require.Empty(t, got, tc.sql)
			continue
		}
		require.Len(t, got, 1, tc.sql)
		require.Equal(t, SDLSyntaxError, got[0].Code, tc.sql)
		require.Contains(t, got[0].Detail, tc.detail, tc.sql)
		require.Equal(t, tc.start, got[0].Start, tc.sql)
	}
}
