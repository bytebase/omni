package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// checkSDLIntegrityLines runs CheckSDLIntegrity and turns every Start into a
// 1-based line of the file it belongs to, so expectations read as lines.
func checkSDLIntegrityLines(files []SDLFile) map[string][]SDLFinding {
	texts := make(map[string]string)
	for _, f := range files {
		texts[f.Path] = f.Text
	}
	line := func(path string, start int) int {
		if start < 0 {
			return start
		}
		return strings.Count(texts[path][:start], "\n") + 1
	}
	got := CheckSDLIntegrity(files)
	for path, findings := range got {
		for i := range findings {
			f := &findings[i]
			if f.Code == SDLSyntaxError {
				continue
			}
			f.Start = line(path, f.Start)
			if f.First != nil {
				firstPath := f.First.Path
				if firstPath == "" {
					firstPath = path
				}
				f.First.Start = line(firstPath, f.First.Start)
			}
			for j := range f.PrimaryKeys {
				f.PrimaryKeys[j].Start = line(path, f.PrimaryKeys[j].Start)
			}
		}
	}
	return got
}

func checkOneSDLFile(text string) []SDLFinding {
	return checkSDLIntegrityLines([]SDLFile{{Path: "f.sql", Text: text}})["f.sql"]
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
			name: "valid composite PK, composite FK, self-referencing FK, unnamed FK, and type-varying FKs",
			sql: "CREATE TABLE public.users (\n" +
				"  id int NOT NULL,\n  tenant_id int NOT NULL,\n  email VARCHAR(255) NOT NULL,\n  bio TEXT NOT NULL,\n" +
				"  CONSTRAINT pk_users PRIMARY KEY (id, tenant_id)\n);\n" +
				"CREATE TABLE public.orders (\n" +
				"  id bigint NOT NULL,\n  user_id int NOT NULL,\n  tenant_id int NOT NULL,\n" +
				"  user_email VARCHAR(100) NOT NULL,\n  user_bio VARCHAR(255) NOT NULL,\n" +
				"  CONSTRAINT pk_orders PRIMARY KEY (id),\n" +
				"  CONSTRAINT fk_orders_user FOREIGN KEY (user_id, tenant_id) REFERENCES public.users(id, tenant_id),\n" +
				"  CONSTRAINT fk_orders_email FOREIGN KEY (user_email) REFERENCES public.users(email),\n" +
				"  FOREIGN KEY (user_bio) REFERENCES public.users(bio)\n);\n" +
				"CREATE TABLE public.employees (\n" +
				"  id int NOT NULL,\n  manager_id int,\n" +
				"  CONSTRAINT pk_employees PRIMARY KEY (id),\n" +
				"  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES public.employees(id)\n);\n" +
				"CREATE TABLE public.logs (id bigint, message text);",
		},
		{
			name: "same table name in different schemas, and same table-scoped constraint names on different tables",
			sql: "CREATE TABLE public.users (id int, CONSTRAINT pk_users PRIMARY KEY (id));\n" +
				"CREATE TABLE analytics.users (id bigint, CONSTRAINT pk_users PRIMARY KEY (id));\n" +
				"CREATE TABLE public.orders (\n" +
				"  id int, user_id int, amount numeric,\n" +
				"  CONSTRAINT pk_orders PRIMARY KEY (id),\n" +
				"  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id),\n" +
				"  CONSTRAINT chk_positive CHECK (amount > 0)\n);\n" +
				"CREATE TABLE public.invoices (\n" +
				"  id int, user_id int, total numeric,\n" +
				"  CONSTRAINT pk_invoices PRIMARY KEY (id),\n" +
				"  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id),\n" +
				"  CONSTRAINT chk_positive CHECK (total > 0)\n);\n" +
				"CREATE VIEW public.user_orders AS SELECT u.id, o.id AS order_id FROM public.users u JOIN public.orders o ON u.id = o.user_id;",
		},
		{
			name: "duplicates within a file, in statement order; the duplicate table is dropped whole",
			sql: "CREATE TABLE t (a int, a text, CONSTRAINT c CHECK (a > 0), CONSTRAINT c CHECK (a < 9), CONSTRAINT k UNIQUE (a));\n" +
				"CREATE TABLE public.t (zz int, zz int);\n" +
				"CREATE TABLE u (a int, CONSTRAINT k PRIMARY KEY (a));\n" +
				"CREATE INDEX i ON t (a);\nCREATE INDEX i ON public.u (a);\nCREATE INDEX ON t (a);\nCREATE INDEX ON u (a);",
			want: []SDLFinding{
				{Code: SDLDuplicateColumn, Start: 1, Schema: "public", Table: "t", Column: "a"},
				{Code: SDLDuplicateConstraint, Start: 1, Schema: "public", Table: "t", Name: "c", ConstraintType: "CHECK", First: &SDLDefinition{Table: "t", ConstraintType: "CHECK", Start: 1}},
				{Code: SDLDuplicateTable, Start: 2, Schema: "public", Table: "t", First: &SDLDefinition{Start: 1}},
				{Code: SDLDuplicateConstraint, Start: 3, Schema: "public", Table: "u", Name: "k", ConstraintType: "PRIMARY KEY", First: &SDLDefinition{Table: "t", ConstraintType: "UNIQUE", Start: 1}},
				{Code: SDLDuplicateIndex, Start: 5, Schema: "public", Name: "i", First: &SDLDefinition{Start: 4}},
			},
		},
		{
			name: "duplicate constraint scoping: same-table FK, cross-table PK, cross-table UNIQUE",
			sql: "CREATE TABLE public.users (id int, email text, CONSTRAINT pk_common PRIMARY KEY (id), CONSTRAINT uk_email UNIQUE (email));\n" +
				"CREATE TABLE public.posts (\n" +
				"  id int, email text, author_id int, editor_id int,\n" +
				"  CONSTRAINT pk_common PRIMARY KEY (id),\n" +
				"  CONSTRAINT uk_email UNIQUE (email),\n" +
				"  CONSTRAINT fk_user FOREIGN KEY (author_id) REFERENCES public.users(id),\n" +
				"  CONSTRAINT fk_user FOREIGN KEY (editor_id) REFERENCES public.users(id)\n);",
			want: []SDLFinding{
				{Code: SDLDuplicateConstraint, Start: 4, Schema: "public", Table: "posts", Name: "pk_common", ConstraintType: "PRIMARY KEY", First: &SDLDefinition{Table: "users", ConstraintType: "PRIMARY KEY", Start: 1}},
				{Code: SDLDuplicateConstraint, Start: 5, Schema: "public", Table: "posts", Name: "uk_email", ConstraintType: "UNIQUE", First: &SDLDefinition{Table: "users", ConstraintType: "UNIQUE", Start: 1}},
				{Code: SDLDuplicateConstraint, Start: 7, Schema: "public", Table: "posts", Name: "fk_user", ConstraintType: "FOREIGN KEY", First: &SDLDefinition{Table: "posts", ConstraintType: "FOREIGN KEY", Start: 6}},
			},
		},
		{
			name: "multiple primary keys, reported at the second; column-level keys are not counted",
			sql:  "CREATE TABLE s.t (\n  a int PRIMARY KEY,\n  PRIMARY KEY (a),\n  CONSTRAINT p2 PRIMARY KEY (a)\n);",
			want: []SDLFinding{{
				Code: SDLMultiplePrimaryKeys, Start: 4, Schema: "s", Table: "t",
				PrimaryKeys: []SDLDefinition{{Start: 3}, {Name: "p2", Start: 4}},
			}},
		},
		{
			name: "three primary keys and duplicate constraint ordering (duplicates reported before multiple PKs)",
			sql: "CREATE TABLE public.users (\n" +
				"  id int, email text, username text,\n" +
				"  CONSTRAINT pk1 PRIMARY KEY (id),\n" +
				"  CONSTRAINT pk2 PRIMARY KEY (email),\n" +
				"  CONSTRAINT pk3 PRIMARY KEY (username)\n);\n" +
				"CREATE TABLE public.orders (\n" +
				"  id bigint,\n" +
				"  CONSTRAINT pk1 PRIMARY KEY (id)\n);",
			want: []SDLFinding{
				{Code: SDLDuplicateConstraint, Start: 9, Schema: "public", Table: "orders", Name: "pk1", ConstraintType: "PRIMARY KEY", First: &SDLDefinition{Table: "users", ConstraintType: "PRIMARY KEY", Start: 3}},
				{
					Code: SDLMultiplePrimaryKeys, Start: 4, Schema: "public", Table: "users",
					PrimaryKeys: []SDLDefinition{{Name: "pk1", Start: 3}, {Name: "pk2", Start: 4}, {Name: "pk3", Start: 5}},
				},
			},
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
				{Code: SDLForeignKeyTableNotFound, Start: 5, Schema: "s", Table: "t", Name: "f1", RefSchema: "public", RefTable: "missing", Columns: []string{"a"}},
				{Code: SDLForeignKeyColumnNotFound, Start: 6, Schema: "s", Table: "t", Column: "b", RefSchema: "s", RefTable: "u", RefColumn: "zz", Columns: []string{"x", "y"}},
			},
		},
		{
			name: "view dependencies resolve against the view's schema; nested CTEs are not excluded",
			sql: "CREATE TABLE s.t (a int);\n" +
				"CREATE VIEW s.v AS SELECT * FROM t, public.t, (WITH d AS (SELECT 1) SELECT * FROM d) x;\n" +
				"CREATE VIEW w AS SELECT * FROM v3;\nCREATE VIEW w AS SELECT * FROM s.v;",
			want: []SDLFinding{
				{Code: SDLViewDependencyNotFound, Start: 2, Schema: "s", Name: "v", RefSchema: "public", RefTable: "t"},
				{Code: SDLViewDependencyNotFound, Start: 2, Schema: "s", Name: "v", RefSchema: "s", RefTable: "d"},
			},
		},
		{
			name: "view JOIN with one existing and one missing table reports only the missing table",
			sql: "CREATE TABLE public.users (id int);\n" +
				"CREATE VIEW public.user_orders AS SELECT u.id, o.id AS order_id FROM public.users u JOIN public.orders o ON u.id = o.user_id;",
			want: []SDLFinding{
				{Code: SDLViewDependencyNotFound, Start: 2, Schema: "public", Name: "user_orders", RefSchema: "public", RefTable: "orders"},
			},
		},
		{
			name: "view with multiple chained CTEs and CTE referencing missing table",
			sql: "CREATE TABLE public.orders (id int, product_id int, total numeric);\n" +
				"CREATE TABLE public.products (id int, name text);\n" +
				"CREATE VIEW public.sales_summary AS\n" +
				"  WITH high_value AS (SELECT id, product_id, total FROM public.orders WHERE total > 1000),\n" +
				"       product_stats AS (SELECT p.name, COUNT(h.id) FROM high_value h JOIN public.products p ON h.product_id = p.id GROUP BY p.name)\n" +
				"  SELECT * FROM product_stats;\n" +
				"CREATE VIEW public.broken_cte AS\n" +
				"  WITH valid_cte AS (SELECT id FROM public.missing_orders)\n" +
				"  SELECT * FROM valid_cte;",
			want: []SDLFinding{
				{Code: SDLViewDependencyNotFound, Start: 7, Schema: "public", Name: "broken_cte", RefSchema: "public", RefTable: "missing_orders"},
			},
		},
		{
			name: "system schemas, and system table and view names in any schema, are skipped",
			sql:  "CREATE VIEW v AS SELECT * FROM s.pg_class, s.pg_tables, pg_temp_3.x, pg_toast.y, information_schema.tables, information_schema.columns, pg_catalog.pg_index;",
		},
		{
			name: "qualified names are joined with a dot, so quoted dots collide",
			sql:  "CREATE TABLE \"a.b\".c (id int);\nCREATE TABLE a.\"b.c\" (id int);",
			want: []SDLFinding{{Code: SDLDuplicateTable, Start: 2, Schema: "a", Table: "b.c", First: &SDLDefinition{Start: 1}}},
		},
		{name: "empty", sql: "", want: nil},
		{name: "comments only", sql: "-- first comment\n-- second comment", want: nil},
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
	t.Run("cross-file duplicates and first-definition FK resolution", func(t *testing.T) {
		files := []SDLFile{
			{Path: "a.sql", Text: "CREATE TABLE s.t (a int, CONSTRAINT pk PRIMARY KEY (a));\nCREATE INDEX i ON s.t (a);\nCREATE VIEW s.v AS SELECT * FROM s.t;"},
			{Path: "b.sql", Text: "\nCREATE TABLE s.t (b int, CONSTRAINT pk UNIQUE (b), CONSTRAINT c CHECK (b > 0));\nCREATE INDEX i ON s.t (b);\nCREATE VIEW s.v AS SELECT * FROM s.t;"},
			{Path: "c.sql", Text: "CREATE TABLE s.r (x int, CONSTRAINT fk FOREIGN KEY (x) REFERENCES s.t(b), CONSTRAINT c CHECK (x > 0));"},
		}
		first := &SDLDefinition{Path: "a.sql", Start: 1}
		require.Equal(t, map[string][]SDLFinding{
			"a.sql": {},
			"b.sql": {
				{Code: SDLDuplicateTableAcrossFiles, Start: 2, Schema: "s", Table: "t", First: first},
				{Code: SDLDuplicateIndexAcrossFiles, Start: 3, Schema: "s", Name: "i", First: &SDLDefinition{Path: "a.sql", Start: 2}},
				{Code: SDLDuplicateViewAcrossFiles, Start: 4, Schema: "s", Name: "v", First: &SDLDefinition{Path: "a.sql", Start: 3}},
				{Code: SDLDuplicateConstraintAcrossFiles, Start: 2, Schema: "s", Table: "t", Name: "pk", ConstraintType: "UNIQUE", First: &SDLDefinition{Path: "a.sql", Table: "t", ConstraintType: "PRIMARY KEY", Start: 1}},
			},
			// The foreign key resolves against a.sql's s.t, the first definition.
			"c.sql": {
				{Code: SDLForeignKeyColumnNotFound, Start: 1, Schema: "s", Table: "r", Name: "fk", Column: "x", RefSchema: "s", RefTable: "t", RefColumn: "b", Columns: []string{"a"}},
			},
		}, checkSDLIntegrityLines(files))
	})

	t.Run("valid multi-file FKs, views, chained views, CTEs, system tables, table-scoped names, and empty files", func(t *testing.T) {
		files := []SDLFile{
			{Path: "empty1.sql", Text: ""},
			{Path: "empty2.sql", Text: "   \n\n  "},
			{Path: "01_base.sql", Text: `
				CREATE TABLE public.organizations (
					id SERIAL PRIMARY KEY,
					name TEXT NOT NULL
				);
				CREATE TABLE public.users (
					id SERIAL PRIMARY KEY,
					org_id INTEGER NOT NULL,
					email TEXT UNIQUE,
					active BOOLEAN NOT NULL DEFAULT true,
					CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES public.organizations(id),
					CONSTRAINT chk_nonempty CHECK (email <> '')
				);
			`},
			{Path: "02_business.sql", Text: `
				CREATE TABLE public.projects (
					id SERIAL PRIMARY KEY,
					org_id INTEGER NOT NULL,
					owner_id INTEGER,
					name TEXT NOT NULL,
					CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES public.organizations(id),
					FOREIGN KEY (owner_id) REFERENCES public.users(id),
					CONSTRAINT chk_nonempty CHECK (name <> '')
				);
				CREATE TABLE public.tasks (
					id SERIAL PRIMARY KEY,
					project_id INTEGER NOT NULL,
					assignee_id INTEGER,
					FOREIGN KEY (project_id) REFERENCES public.projects(id),
					FOREIGN KEY (assignee_id) REFERENCES public.users(id)
				);
			`},
			{Path: "03_views.sql", Text: `
				CREATE VIEW public.active_users AS
					SELECT id, email FROM public.users WHERE active = true;
				CREATE VIEW public.active_user_summary AS
					SELECT id, email FROM public.active_users;
				CREATE VIEW public.project_stats AS
					WITH open_tasks AS (SELECT id, project_id FROM public.tasks)
					SELECT p.id, o.name, COUNT(t.id)
					FROM public.projects p
					JOIN public.organizations o ON p.org_id = o.id
					LEFT JOIN open_tasks t ON t.project_id = p.id
					CROSS JOIN pg_tables pt
					GROUP BY p.id, o.name;
			`},
		}
		got := checkSDLIntegrityLines(files)
		for _, f := range files {
			require.Empty(t, got[f.Path], f.Path)
		}
	})

	t.Run("multiple cross-file and local errors in one file, plus mixed system functions in view", func(t *testing.T) {
		files := []SDLFile{
			{Path: "base.sql", Text: "CREATE TABLE public.users (id int);\nCREATE TABLE public.products (id int, name text);"},
			{Path: "errors.sql", Text: "CREATE TABLE public.users (id int);\n" +
				"CREATE TABLE public.orders (\n  id int, customer_id int,\n" +
				"  FOREIGN KEY (customer_id) REFERENCES public.customers(id)\n);\n" +
				"CREATE VIEW public.order_analysis AS\n" +
				"  SELECT o.id, p.name, pg_size_pretty(pg_table_size('orders'))\n" +
				"  FROM public.missing_orders o JOIN public.products p ON o.product_id = p.id;"},
		}
		require.Equal(t, map[string][]SDLFinding{
			"base.sql": {},
			"errors.sql": {
				{Code: SDLDuplicateTableAcrossFiles, Start: 1, Schema: "public", Table: "users", First: &SDLDefinition{Path: "base.sql", Start: 1}},
				{Code: SDLForeignKeyTableNotFound, Start: 4, Schema: "public", Table: "orders", RefSchema: "public", RefTable: "customers", Columns: []string{"customer_id"}},
				{Code: SDLViewDependencyNotFound, Start: 6, Schema: "public", Name: "order_analysis", RefSchema: "public", RefTable: "missing_orders"},
			},
		}, checkSDLIntegrityLines(files))
	})

	t.Run("syntax errors stop at the first failing file", func(t *testing.T) {
		got := CheckSDLIntegrity([]SDLFile{
			{Path: "ok.sql", Text: "CREATE TABLE t (a int);"},
			{Path: "bad1.sql", Text: "CREATE TABLE t (a int;"},
			{Path: "bad2.sql", Text: "CREATE TABLEX t (a int);"},
		})
		require.Equal(t, map[string][]SDLFinding{
			"bad1.sql": {{Code: SDLSyntaxError, Start: 15, Detail: "syntax error: unmatched opening parenthesis"}},
		}, got)

		gotParse := CheckSDLIntegrity([]SDLFile{
			{Path: "ok.sql", Text: "CREATE TABLE t (a int);"},
			{Path: "bad.sql", Text: "CREATE TABLEX t (a int);"},
		})
		require.Len(t, gotParse, 1)
		require.Len(t, gotParse["bad.sql"], 1)
		require.Equal(t, SDLSyntaxError, gotParse["bad.sql"][0].Code)
		require.Contains(t, gotParse["bad.sql"][0].Detail, `syntax error at or near "TABLEX"`)
	})

	require.Empty(t, CheckSDLIntegrity(nil))
}
