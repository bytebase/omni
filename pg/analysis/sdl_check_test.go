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
				email TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				email_lower TEXT GENERATED ALWAYS AS (LOWER(email)) STORED,
				CONSTRAINT pk_users PRIMARY KEY (id),
				CONSTRAINT uk_users_email UNIQUE (email)
			);`,
		},
		{
			name: "conforming complex schema with named constraints, index, sequence, view, function",
			sql: `CREATE TABLE public.users (
				id INTEGER,
				email TEXT,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);
			CREATE TABLE public.products (
				id INTEGER,
				CONSTRAINT pk_products PRIMARY KEY (id)
			);
			CREATE TABLE public.orders (
				id INTEGER,
				user_id INTEGER,
				product_id INTEGER,
				quantity INTEGER,
				total_price NUMERIC,
				CONSTRAINT pk_orders PRIMARY KEY (id),
				CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id),
				CONSTRAINT fk_orders_product FOREIGN KEY (product_id) REFERENCES public.products(id),
				CONSTRAINT uk_orders_composite UNIQUE (user_id, product_id),
				CONSTRAINT chk_orders_quantity CHECK (quantity > 0),
				CONSTRAINT chk_orders_price CHECK (total_price >= 0)
			);
			CREATE UNIQUE INDEX idx_users_email ON public.users (email);
			CREATE SEQUENCE public.user_id_seq START 1;
			ALTER SEQUENCE public.user_id_seq OWNED BY public.users.id;
			ALTER SEQUENCE public.user_id_seq INCREMENT BY 5 MINVALUE 1 MAXVALUE 1000;`,
		},
		{
			name: "a statement starts at its leading comment",
			sql:  "SELECT 1;\n-- users\nCREATE TABLE users (id int);",
			want: []want{{SDLFinding{Code: SDLTableWithoutSchema, Table: "users"}, "", "-- users"}},
		},
		{
			name: "column constraints: primary key, unique, multiple on one column, and multiple columns",
			sql: "CREATE TABLE public.t1 (\n  id int PRIMARY KEY,\n  email text UNIQUE\n);\n" +
				"CREATE TABLE public.t2 (\n  id int PRIMARY KEY UNIQUE CHECK (id > 0)\n);",
			want: []want{
				{SDLFinding{Code: SDLColumnConstraint, Table: "t1", Column: "id", ConstraintTypes: []string{"PRIMARY KEY"}}, "PRIMARY KEY,", "CREATE TABLE public.t1"},
				{SDLFinding{Code: SDLColumnConstraint, Table: "t1", Column: "email", ConstraintTypes: []string{"UNIQUE"}}, "UNIQUE\n", "CREATE TABLE public.t1"},
				{SDLFinding{Code: SDLColumnConstraint, Table: "t2", Column: "id", ConstraintTypes: []string{"PRIMARY KEY", "UNIQUE", "CHECK"}}, "PRIMARY KEY UNIQUE", "CREATE TABLE public.t2"},
			},
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
			name: "unnamed table constraints, and foreign keys without schema",
			sql: "CREATE TABLE public.t (\n  a int,\n  b int,\n" +
				"  PRIMARY KEY (a),\n  UNIQUE (b),\n  CHECK (a > 0),\n" +
				"  FOREIGN KEY (a) REFERENCES u(id),\n  FOREIGN KEY (b) REFERENCES public.u(id),\n" +
				"  EXCLUDE USING gist (a WITH =)\n);",
			want: []want{
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "PRIMARY KEY"}, "PRIMARY", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "UNIQUE"}, "UNIQUE", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "CHECK"}, "CHECK", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "FOREIGN KEY"}, "FOREIGN KEY (a)", "CREATE"},
				{SDLFinding{Code: SDLForeignKeyWithoutSchema, Table: "t", RefTable: "u"}, "FOREIGN KEY (a)", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "FOREIGN KEY"}, "FOREIGN KEY (b)", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "t", ConstraintType: "CONSTRAINT"}, "EXCLUDE", "CREATE"},
			},
		},
		{
			name: "named table foreign keys without schema",
			sql: "CREATE TABLE public.orders (\n  user_id int,\n  product_id int,\n" +
				"  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),\n" +
				"  CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)\n);",
			want: []want{
				{SDLFinding{Code: SDLForeignKeyWithoutSchema, Table: "orders", RefTable: "users"}, "CONSTRAINT fk_user", "CREATE"},
				{SDLFinding{Code: SDLForeignKeyWithoutSchema, Table: "orders", RefTable: "products"}, "CONSTRAINT fk_product", "CREATE"},
			},
		},
		{
			name: "table with missing schema, disallowed column constraint, and unnamed table constraint",
			sql:  "CREATE TABLE users (\n  id int PRIMARY KEY,\n  UNIQUE (id)\n);",
			want: []want{
				{SDLFinding{Code: SDLTableWithoutSchema, Table: "users"}, "", "CREATE"},
				{SDLFinding{Code: SDLColumnConstraint, Table: "users", Column: "id", ConstraintTypes: []string{"PRIMARY KEY"}}, "PRIMARY KEY", "CREATE"},
				{SDLFinding{Code: SDLUnnamedConstraint, Table: "users", ConstraintType: "UNIQUE"}, "UNIQUE", "CREATE"},
			},
		},
		{
			name: "index without name, without schema, unique, and quoted ONLY table",
			sql: "CREATE INDEX ON t (a);\nCREATE INDEX ON public.t (a);\nCREATE INDEX i ON t (a);\n" +
				"CREATE UNIQUE INDEX ui ON t (a);\nCREATE INDEX \"idx_dept\" ON ONLY \"departments\" (budget);\n" +
				"CREATE INDEX i2 ON public.t (a);\nCREATE UNIQUE INDEX ui2 ON public.t (a);",
			want: []want{
				{SDLFinding{Code: SDLUnnamedIndex}, "", "CREATE INDEX ON t"},
				{SDLFinding{Code: SDLIndexWithoutSchema, Table: "t"}, "", "CREATE INDEX ON t"},
				{SDLFinding{Code: SDLUnnamedIndex}, "", "CREATE INDEX ON public.t"},
				{SDLFinding{Code: SDLIndexWithoutSchema, Name: "i", Table: "t"}, "", "CREATE INDEX i "},
				{SDLFinding{Code: SDLIndexWithoutSchema, Name: "ui", Table: "t"}, "", "CREATE UNIQUE INDEX ui "},
				{SDLFinding{Code: SDLIndexWithoutSchema, Name: "idx_dept", Table: "departments"}, "", `CREATE INDEX "idx_dept"`},
			},
		},
		{
			name: "views, sequences and functions without schema",
			sql: "CREATE VIEW v AS SELECT 1;\nCREATE SEQUENCE seq;\nALTER SEQUENCE seq OWNED BY public.t.id;\n" +
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
			name: "multiple statements with mixed errors preserve statement order",
			sql: "CREATE TABLE users (\n  id int PRIMARY KEY\n);\n" +
				"CREATE INDEX idx_users_email ON users (email);\n" +
				"CREATE VIEW active_users AS SELECT * FROM users;",
			want: []want{
				{SDLFinding{Code: SDLTableWithoutSchema, Table: "users"}, "", "CREATE TABLE users"},
				{SDLFinding{Code: SDLColumnConstraint, Table: "users", Column: "id", ConstraintTypes: []string{"PRIMARY KEY"}}, "PRIMARY KEY", "CREATE TABLE users"},
				{SDLFinding{Code: SDLIndexWithoutSchema, Name: "idx_users_email", Table: "users"}, "", "CREATE INDEX idx_users_email"},
				{SDLFinding{Code: SDLViewWithoutSchema, Name: "active_users"}, "", "CREATE VIEW active_users"},
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
