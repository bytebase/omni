package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSDLStyle(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []SDLFinding
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
			name: "table without schema at the statement's first line, leading comment included",
			sql:  "\n-- users\nCREATE TABLE users (id int);",
			want: []SDLFinding{{Code: SDLTableWithoutSchema, Line: 2, Table: "users"}},
		},
		{
			name: "column constraints: references finding first, column finding at its first disallowed constraint",
			sql:  "CREATE TABLE public.t (\n  a int NOT NULL\n    CHECK (a > 0)\n    REFERENCES u(id)\n);",
			want: []SDLFinding{
				{Code: SDLColumnForeignKeyWithoutSchema, Line: 4, Table: "t", Column: "a", RefTable: "u"},
				{Code: SDLColumnConstraint, Line: 3, Table: "t", Column: "a", ConstraintTypes: []string{"CHECK", "FOREIGN KEY"}},
			},
		},
		{
			name: "column references with a schema is still a column constraint",
			sql:  "CREATE TABLE public.t (a int REFERENCES public.u(id));",
			want: []SDLFinding{{Code: SDLColumnConstraint, Line: 1, Table: "t", Column: "a", ConstraintTypes: []string{"FOREIGN KEY"}}},
		},
		{
			name: "unnamed table constraints, and a foreign key without schema",
			sql:  "CREATE TABLE public.t (\n  a int,\n  PRIMARY KEY (a),\n  FOREIGN KEY (a) REFERENCES u(id),\n  EXCLUDE USING gist (a WITH =)\n);",
			want: []SDLFinding{
				{Code: SDLUnnamedConstraint, Line: 3, Table: "t", ConstraintType: "PRIMARY KEY"},
				{Code: SDLUnnamedConstraint, Line: 4, Table: "t", ConstraintType: "FOREIGN KEY"},
				{Code: SDLForeignKeyWithoutSchema, Line: 4, Table: "t", RefTable: "u"},
				{Code: SDLUnnamedConstraint, Line: 5, Table: "t", ConstraintType: "CONSTRAINT"},
			},
		},
		{
			name: "index without name and schema",
			sql:  "CREATE INDEX ON t (a);\nCREATE INDEX i ON t (a);\nCREATE INDEX i2 ON public.t (a);",
			want: []SDLFinding{
				{Code: SDLUnnamedIndex, Line: 1},
				{Code: SDLIndexWithoutSchema, Line: 1, Table: "t"},
				{Code: SDLIndexWithoutSchema, Line: 2, Name: "i", Table: "t"},
			},
		},
		{
			name: "views, sequences and functions without schema",
			sql: "CREATE VIEW v AS SELECT 1;\nCREATE SEQUENCE seq;\nALTER SEQUENCE seq RESTART;\n" +
				"CREATE FUNCTION f() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;\n" +
				"CREATE VIEW public.v2 AS SELECT 1;\nCREATE SEQUENCE public.s2;\nALTER SEQUENCE public.s2 RESTART;\n" +
				"CREATE FUNCTION public.f2() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;",
			want: []SDLFinding{
				{Code: SDLViewWithoutSchema, Line: 1, Name: "v"},
				{Code: SDLSequenceWithoutSchema, Line: 2, Name: "seq"},
				{Code: SDLAlterSequenceWithoutSchema, Line: 3, Name: "seq"},
				{Code: SDLFunctionWithoutSchema, Line: 4, Name: "f"},
			},
		},
		{
			name: "a three-part function name takes its catalog as the schema",
			sql:  "CREATE FUNCTION c.s.f() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql;",
		},
		{
			name: "a CreateStmt's separate Constraints list and other statements are not judged",
			sql:  "CREATE MATERIALIZED VIEW mv AS SELECT 1;\nDROP TABLE t;\nALTER TABLE t ADD COLUMN x int PRIMARY KEY;",
		},
		{
			name: "statement lines count only LF, constraint lines count CR, LF and CRLF",
			sql:  "CREATE TABLE public.a (x int);\rCREATE TABLE public.t (\r\n  a int,\r  PRIMARY KEY (a)\n);\nCREATE TABLE b (x int);",
			want: []SDLFinding{
				{Code: SDLUnnamedConstraint, Line: 4, Table: "t", ConstraintType: "PRIMARY KEY"},
				{Code: SDLTableWithoutSchema, Line: 4, Table: "b"},
			},
		},
		{name: "empty", sql: ""},
		{name: "comment only", sql: "-- nothing\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, CheckSDLStyle(tc.sql))
		})
	}
}

func TestCheckSDLStyleSyntaxError(t *testing.T) {
	tests := []struct {
		sql    string
		detail string
	}{
		{"CREATE TABLE public.t (a int;", "syntax error: unmatched opening parenthesis"},
		{"CREATE TABLE public.t (a int));", "syntax error: unmatched closing parenthesis"},
		// The balance check skips strings, dollar quotes and comments, but
		// ignores backslash escapes and does not nest block comments.
		{"CREATE TABLE public.t (a text DEFAULT E'\\'(');", "syntax error: unmatched opening parenthesis"},
		{"/* /* */ ( */ CREATE TABLE public.t (a int);", "syntax error: unmatched opening parenthesis"},
		{"CREATE TABLE public.t (a text DEFAULT '(', b text DEFAULT $x$)$x$); -- (", ""},
		{"CREATE TABLEX public.t (a int);", `syntax error at or near "TABLEX"`},
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
		require.Zero(t, got[0].Line, tc.sql)
	}
}
