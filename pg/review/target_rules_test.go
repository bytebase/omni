package review

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

func database(searchPath string, schemas ...*metadata.SchemaMetadata) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: "db", SearchPath: searchPath, Schemas: schemas}
}

func schema(name string, tables ...*metadata.TableMetadata) *metadata.SchemaMetadata {
	return &metadata.SchemaMetadata{Name: name, Tables: tables}
}

// table builds a table from "name type" column specs.
func table(name string, columns ...string) *metadata.TableMetadata {
	t := &metadata.TableMetadata{Name: name}
	for _, c := range columns {
		nameType := strings.SplitN(c, " ", 2)
		t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: nameType[0], Type: nameType[1], Nullable: true})
	}
	return t
}

func withPrimaryKey(t *metadata.TableMetadata, name string, columns ...string) *metadata.TableMetadata {
	for _, c := range t.Columns {
		for _, k := range columns {
			if c.Name == k {
				c.Nullable = false
			}
		}
	}
	t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Expressions: columns, Primary: true, Unique: true, IsConstraint: true, Type: "btree"})
	return t
}

func withUnique(t *metadata.TableMetadata, name string, columns ...string) *metadata.TableMetadata {
	t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: name, Expressions: columns, Unique: true, IsConstraint: true, Type: "btree"})
	return t
}

func withCheck(t *metadata.TableMetadata, name, expr string) *metadata.TableMetadata {
	t.CheckConstraints = append(t.CheckConstraints, &metadata.CheckConstraintMetadata{Name: name, Expression: expr})
	return t
}

func withFunctionReturning(db *metadata.DatabaseSchemaMetadata, schema, table string) *metadata.DatabaseSchemaMetadata {
	for _, s := range db.Schemas {
		if s.Name == schema {
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: "f", DependencyTables: []*metadata.DependencyTable{{Schema: schema, Table: table}}})
		}
	}
	return db
}

func withEventTrigger(db *metadata.DatabaseSchemaMetadata, enabled bool, tags ...string) *metadata.DatabaseSchemaMetadata {
	db.EventTriggers = append(db.EventTriggers, &metadata.EventTriggerMetadata{Name: "et", Event: "DDL_COMMAND_END", Tags: tags, Enabled: enabled})
	return db
}

func withForeignKey(t *metadata.TableMetadata, name, column, refTable, refColumn string) *metadata.TableMetadata {
	t.ForeignKeys = append(t.ForeignKeys, &metadata.ForeignKeyMetadata{Name: name, Columns: []string{column}, ReferencedSchema: "public", ReferencedTable: refTable, ReferencedColumns: []string{refColumn}})
	return t
}

// shop is a small database: parent t with every constraint kind, a table
// without a primary key, and a second schema holding another t.
func shop() *metadata.DatabaseSchemaMetadata {
	return database(`"$user", public`,
		schema("public",
			withPrimaryKey(table("p", "id integer", "note text"), "p_pkey", "id"),
			withCheck(withForeignKey(withUnique(withPrimaryKey(table("t", "id integer", "code text", "p_id integer", "n integer"), "t_pkey", "id"), "t_code_key", "code"), "t_p_fk", "p_id", "p", "id"), "t_n_check", "(n > 0)"),
			table("nokey", "a integer"),
		),
		schema("s", withUnique(table("t", "id integer"), "t_pkey", "id")),
	)
}

type targetFinding struct {
	statement int
	text      string // the anchored text; "" for a zero range
	message   string
	targets   []int
}

func runTargetRule(t *testing.T, rule review.Rule, sql string, change review.Change, targets []review.Target) []targetFinding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{rule}, Change: change}, targets)
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	if len(result.Failures) != 0 {
		t.Fatalf("Failures = %v", result.Failures)
	}
	var got []targetFinding
	for _, f := range result.Findings {
		if f.Rule == review.Syntax {
			t.Fatalf("syntax finding: %s", f.Message)
		}
		if f.Rule != rule {
			continue
		}
		got = append(got, targetFinding{f.Statement, sql[f.Range.Start:f.Range.End], f.Message, f.Targets})
	}
	return got
}

func TestDisallowDropConstraint(t *testing.T) {
	one := []review.Target{{Schema: shop(), SessionUser: "alice"}}
	tests := []struct {
		name    string
		sql     string
		targets []review.Target
		want    []targetFinding
	}{
		{
			name:    "every kind the target had, at its subcommand",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE, DROP CONSTRAINT t_code_key, DROP CONSTRAINT t_p_fk, DROP CONSTRAINT t_n_check;",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_pkey CASCADE", "drops primary key t_pkey of t", []int{0}},
				{0, "DROP CONSTRAINT t_code_key", "drops unique constraint t_code_key of t", []int{0}},
				{0, "DROP CONSTRAINT t_p_fk", "drops foreign key t_p_fk of t", []int{0}},
				{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}},
			},
		},
		{
			name:    "the type comes from the table the name resolves to",
			sql:     "ALTER TABLE s.t DROP CONSTRAINT t_pkey;\nSET search_path = s;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_pkey;",
			targets: []review.Target{{Schema: shop(), SessionUser: "alice"}},
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of s.t", []int{0}},
			},
		},
		{
			name:    "a renamed table is not looked up under either name",
			sql:     "ALTER TABLE t RENAME TO u;\nALTER TABLE u DROP CONSTRAINT t_n_check;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_n_check;",
			targets: one,
		},
		{
			name:    "a column the change added",
			sql:     "ALTER TABLE t ADD COLUMN tmp int;\nALTER TABLE t DROP COLUMN tmp;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
		},
		{
			name:    "a constraint the change added is not reported",
			sql:     "ALTER TABLE nokey ADD CONSTRAINT nokey_a_check CHECK (a > 0);\nALTER TABLE nokey DROP CONSTRAINT nokey_a_check;\nALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_n_check CHECK (n > 1);\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{2, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "a missing constraint under IF EXISTS drops nothing",
			sql:     "ALTER TABLE t DROP CONSTRAINT IF EXISTS nope;",
			targets: one,
		},
		{
			name:    "a statement the server refuses is not reported, and the scan ends",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP COLUMN nope;\nALTER TABLE t DROP CONSTRAINT t_code_key;",
			targets: one,
		},
		{
			name:    "dropping a missing constraint is refused",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT nope;",
			targets: one,
		},
		{
			name:    "dropping a referenced key is refused without CASCADE",
			sql:     "ALTER TABLE p DROP CONSTRAINT p_pkey;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
		},
		{
			name:    "CASCADE drops the foreign keys on the key",
			sql:     "ALTER TABLE p DROP CONSTRAINT p_pkey CASCADE;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_p_fk;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT p_pkey CASCADE", "drops primary key p_pkey of p", []int{0}},
				{2, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}},
			},
		},
		{
			name:    "DROP TABLE CASCADE drops the foreign keys on it",
			sql:     "DROP TABLE p CASCADE;\nALTER TABLE t DROP CONSTRAINT IF EXISTS t_p_fk;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{2, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "the same names in two schemas are two constraints",
			sql:     "ALTER TABLE s.t DROP CONSTRAINT t_pkey;\nALTER TABLE public.t DROP CONSTRAINT t_pkey;",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of s.t", []int{0}},
				{1, "DROP CONSTRAINT t_pkey", "drops primary key t_pkey of public.t", []int{0}},
			},
		},
		{
			name:    "a dropped table no longer holds its foreign keys",
			sql:     "DROP TABLE t;\nALTER TABLE p DROP CONSTRAINT p_pkey;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT p_pkey", "drops primary key p_pkey of p", []int{0}}},
		},
		{
			name:    "dropping a table a foreign key references is refused without CASCADE",
			sql:     "DROP TABLE p;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
		},
		{
			name:    "dropping both ends of a foreign key together",
			sql:     "DROP TABLE p, t;\nALTER TABLE nokey ADD CONSTRAINT nokey_a CHECK (a > 0);\nALTER TABLE s.t DROP CONSTRAINT t_pkey;",
			targets: one,
			want:    []targetFinding{{2, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of s.t", []int{0}}},
		},
		{
			name:    "a subcommand certain to fail refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD COLUMN id int;\nALTER TABLE t DROP CONSTRAINT t_code_key;",
			targets: one,
		},
		{
			name:    "altering a missing column refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ALTER COLUMN nope SET DEFAULT 1;",
			targets: one,
		},
		{
			name:    "adding a constraint under a name the table uses refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_code_key CHECK (n > 1);",
			targets: one,
		},
		{
			name:    "a second primary key refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD PRIMARY KEY (code);",
			targets: one,
		},
		{
			name:    "the drops run first, so a name or a column they free can be used again",
			sql:     "ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (n > 1), DROP CONSTRAINT t_n_check, DROP COLUMN code, ADD COLUMN code int, ADD COLUMN IF NOT EXISTS id int;\nALTER TABLE s.t DROP CONSTRAINT t_pkey, ADD COLUMN x int, ALTER COLUMN id SET NOT NULL;",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}},
				{1, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of s.t", []int{0}},
			},
		},
		{
			name:    "a cascading drop of an unrelated column keeps the foreign keys on the key",
			sql:     "ALTER TABLE p DROP COLUMN note CASCADE;\nALTER TABLE t DROP CONSTRAINT t_p_fk;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_p_fk", "drops foreign key t_p_fk of t", []int{0}}},
		},
		{
			name: "a cascading key drop takes only the foreign keys on that table",
			sql:  "ALTER TABLE s.p DROP CONSTRAINT p_pkey CASCADE;\nALTER TABLE public.c DROP CONSTRAINT c_fk;",
			targets: []review.Target{{Schema: database("public",
				schema("public", withPrimaryKey(table("p", "id integer"), "p_pkey", "id"), withForeignKey(table("c", "p_id integer"), "c_fk", "p_id", "p", "id")),
				schema("s", withPrimaryKey(table("p", "id integer"), "p_pkey", "id")),
			)}},
			want: []targetFinding{
				{0, "DROP CONSTRAINT p_pkey CASCADE", "drops primary key p_pkey of s.p", []int{0}},
				{1, "DROP CONSTRAINT c_fk", "drops foreign key c_fk of public.c", []int{0}},
			},
		},
		{
			name:    "a second drop of the same constraint refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT t_n_check;\nALTER TABLE s.t DROP CONSTRAINT t_pkey;",
			targets: one,
		},
		{
			name:    "a drop of a constraint a column drop may have taken withholds the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_code_key, DROP COLUMN n, DROP CONSTRAINT t_n_check;\nALTER TABLE s.t DROP CONSTRAINT t_pkey;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of s.t", []int{0}}},
		},
		{
			name:    "a new foreign key blocks only the key it references",
			sql:     "CREATE TABLE c (t_id int REFERENCES t (id));\nALTER TABLE t DROP CONSTRAINT t_code_key;\nCREATE TABLE c2 (t_id int REFERENCES t);\nALTER TABLE t DROP CONSTRAINT t_pkey;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_code_key", "drops unique constraint t_code_key of t", []int{0}}},
		},
		{
			name: "an enabled event trigger on the command leaves its effect unknown",
			sql:  "ALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: []review.Target{
				{Schema: withEventTrigger(shop(), true), SessionUser: "alice"},
				{Schema: withEventTrigger(shop(), true, "ALTER TABLE"), SessionUser: "alice"},
				{Schema: withEventTrigger(shop(), true, "CREATE FUNCTION"), SessionUser: "alice"},
				{Schema: withEventTrigger(shop(), false), SessionUser: "alice"},
			},
			want: []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{2, 3}}},
		},
		{
			name:    "a new foreign key goes with its table or its name",
			sql:     "CREATE TABLE c (t_id int REFERENCES t (id));\nDROP TABLE c;\nCREATE TABLE c2 (t_id int CONSTRAINT c2_fk REFERENCES t (id));\nALTER TABLE c2 DROP CONSTRAINT c2_fk;\nALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: one,
			want:    []targetFinding{{4, "DROP CONSTRAINT t_pkey", "drops primary key t_pkey of t", []int{0}}},
		},
		{
			name:    "a column added twice refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD COLUMN x int, ADD COLUMN x int;",
			targets: one,
		},
		{
			name:    "a cascading drop of anything but a relation ends the scan",
			sql:     "DROP FUNCTION f() CASCADE;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
		},
		{
			name:    "procedural code ends the scan",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check;\nDO $$ BEGIN EXECUTE 'ALTER TABLE t RENAME TO old_t'; END $$;\nALTER TABLE t DROP CONSTRAINT t_code_key;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "after RESET search_path only a qualified name is looked up",
			sql:     "RESET search_path;\nALTER TABLE t DROP CONSTRAINT t_n_check;\nALTER TABLE public.t DROP CONSTRAINT t_code_key;",
			targets: one,
			want:    []targetFinding{{2, "DROP CONSTRAINT t_code_key", "drops unique constraint t_code_key of public.t", []int{0}}},
		},
		{
			name:    "a quoted SET search_path value leaves the path unknown",
			sql:     "SET search_path = 'S';\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
		},
		{
			name:    "a name that may be a system catalog",
			sql:     "ALTER TABLE pg_class DROP CONSTRAINT x;",
			targets: []review.Target{{Schema: database("public", schema("public", withCheck(table("pg_class", "a integer"), "x", "(a > 0)")))}},
		},
		{
			name:    "no schema, no constraint to know",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: []review.Target{{}},
		},
		{
			name:    "a schema synced without its search path",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: []review.Target{{Schema: database("", shop().Schemas...), SessionUser: "alice"}},
		},
		{
			name:    "$user with no session user while a schema could be the user's",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: []review.Target{{Schema: shop()}},
		},
		{
			name:    "$user with no session user and no schema it could name",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: []review.Target{{Schema: database(`"$user", public`, shop().Schemas[0])}},
			want:    []targetFinding{{0, "DROP CONSTRAINT t_pkey", "drops primary key t_pkey of t", []int{0}}},
		},
		{
			name: "targets with equal findings merge, different ones do not",
			sql:  "ALTER TABLE t DROP CONSTRAINT t_pkey;",
			targets: []review.Target{
				{Schema: shop(), SessionUser: "alice"},
				{Schema: database(`public`, schema("public", withUnique(table("t", "id integer"), "t_pkey", "id")))},
				{Schema: shop(), SessionUser: "alice"},
				{Schema: database(`public`, schema("public", table("t", "id integer")))},
			},
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_pkey", "drops primary key t_pkey of t", []int{0, 2}},
				{0, "DROP CONSTRAINT t_pkey", "drops unique constraint t_pkey of t", []int{1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runTargetRule(t, review.DisallowDropConstraint, tt.sql, review.Change{}, tt.targets)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(targetFinding{})); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRequirePrimaryKey(t *testing.T) {
	one := []review.Target{{Schema: shop(), SessionUser: "alice"}}
	tests := []struct {
		name    string
		sql     string
		targets []review.Target
		want    []targetFinding
	}{
		{
			name:    "a new table without a primary key",
			sql:     "CREATE TABLE n (id int);\nCREATE TABLE k (id int PRIMARY KEY);\nCREATE TABLE k2 (a int, b int, CONSTRAINT k2_pk PRIMARY KEY (a, b));",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a key added later in the change counts, under any name the table has by then",
			sql:     "CREATE TABLE n (id int);\nALTER TABLE n RENAME TO m;\nALTER TABLE m ADD PRIMARY KEY (id);\nCREATE TABLE o (id int NOT NULL);\nCREATE UNIQUE INDEX o_id ON o (id);\nALTER TABLE o ADD CONSTRAINT o_pkey PRIMARY KEY USING INDEX o_id;",
			targets: one,
		},
		{
			name:    "a key the table was created with is not followed",
			sql:     "CREATE TABLE n (id int PRIMARY KEY);\nALTER TABLE n DROP CONSTRAINT n_pkey;",
			targets: one,
		},
		{
			name:    "a key added in the same statement that drops one",
			sql:     "ALTER TABLE t ADD PRIMARY KEY (code), DROP CONSTRAINT t_pkey CASCADE;",
			targets: one,
		},
		{
			name:    "a key that may be added under another name",
			sql:     "CREATE TABLE n (id int);\nSET search_path = s, public;\nALTER TABLE n ADD PRIMARY KEY (id);",
			targets: one,
		},
		{
			name:    "a table attached as a partition",
			sql:     "CREATE TABLE pt (id int PRIMARY KEY) PARTITION BY RANGE (id);\nCREATE TABLE n (id int NOT NULL);\nALTER TABLE pt ATTACH PARTITION n FOR VALUES FROM (1) TO (10);",
			targets: one,
		},
		{
			name:    "a dropped table takes its indexes, whose names are then free",
			sql:     "DROP TABLE t;\nCREATE TABLE IF NOT EXISTS t_code_key (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE IF NOT EXISTS t_code_key (id int)", "creates table t_code_key without a primary key", []int{0}}},
		},
		{
			name: "a cascading drop takes the views that read the table, in their own schema",
			sql:  "DROP TABLE public.base CASCADE;\nCREATE TABLE IF NOT EXISTS s.reader (id int);",
			targets: []review.Target{{Schema: database("public", schema("public", table("base", "id integer")), &metadata.SchemaMetadata{
				Name:  "s",
				Views: []*metadata.ViewMetadata{{Name: "reader", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}}},
			})}},
			want: []targetFinding{{1, "CREATE TABLE IF NOT EXISTS s.reader (id int)", "creates table s.reader without a primary key", []int{0}}},
		},
		{
			name:    "a qualified move keeps the table's schema exact",
			sql:     "CREATE TABLE public.n (id int);\nALTER TABLE public.n SET SCHEMA s;\nALTER TABLE x.n ADD PRIMARY KEY (id);\nCREATE TABLE m (id int);\nALTER TABLE public.m SET SCHEMA s;\nALTER TABLE x.m ADD PRIMARY KEY (id);",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE public.n (id int)", "creates table public.n without a primary key", []int{0}}},
		},
		{
			name: "a dropped materialized view takes its indexes",
			sql:  "DROP MATERIALIZED VIEW mv;\nCREATE TABLE IF NOT EXISTS mv_idx (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:              "public",
				MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "mv", Indexes: []*metadata.IndexMetadata{{Name: "mv_idx"}}}},
			})}},
			want: []targetFinding{{1, "CREATE TABLE IF NOT EXISTS mv_idx (id int)", "creates table mv_idx without a primary key", []int{0}}},
		},
		{
			name: "dropping a view another view reads is refused without CASCADE",
			sql:  "DROP VIEW v1;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name: "public",
				Views: []*metadata.ViewMetadata{
					{Name: "v1"},
					{Name: "v2", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v1", Column: "a"}}},
				},
			})}},
		},
		{
			name: "dropping a view with the view that reads it",
			sql:  "DROP VIEW v1, v2;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name: "public",
				Views: []*metadata.ViewMetadata{
					{Name: "v1"},
					{Name: "v2", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v1", Column: "a"}}},
				},
			})}},
			want: []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "the tables of CREATE SCHEMA",
			sql:     "CREATE SCHEMA z CREATE TABLE t (id int) CREATE TABLE k (id int PRIMARY KEY) CREATE TABLE s.q (id int);",
			targets: one,
			want: []targetFinding{
				{0, "CREATE TABLE t (id int)", "creates table z.t without a primary key", []int{0}},
				{0, "CREATE TABLE s.q (id int)", "creates table s.q without a primary key", []int{0}},
			},
		},
		{
			name:    "a new view of a table in another schema does not block a drop",
			sql:     "CREATE VIEW v AS SELECT * FROM s.t;\nDROP TABLE public.t;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "dropping a table the target lacks is refused without IF EXISTS",
			sql:     "DROP TABLE public.nope;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "IF EXISTS drops nothing",
			sql:     "DROP TABLE IF EXISTS public.nope, nope2;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "dropping a table a function returns is refused without CASCADE",
			sql:     "DROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withFunctionReturning(shop(), "public", "nokey"), SessionUser: "alice"}},
		},
		{
			name:    "a CTE the view defines is not a relation it reads",
			sql:     "CREATE VIEW v AS WITH t AS (SELECT 1 AS id) SELECT * FROM t;\nDROP TABLE public.t;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a view the change drops no longer reads the table",
			sql:     "CREATE VIEW v AS SELECT * FROM public.t;\nDROP VIEW v;\nDROP TABLE public.t;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a function the change drops no longer returns the table",
			sql:     "DROP FUNCTION f();\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withFunctionReturning(shop(), "public", "nokey"), SessionUser: "alice"}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a function the change creates returning the table blocks its drop",
			sql:     "CREATE FUNCTION g() RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "DROP TABLE of a view is refused",
			sql:  "DROP TABLE v;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:  "public",
				Views: []*metadata.ViewMetadata{{Name: "v"}},
			})}},
		},
		{
			name:    "DROP INDEX of a constraint's index is refused",
			sql:     "DROP INDEX t_pkey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE TABLE of a name the schema holds is refused",
			sql:     "CREATE TABLE nokey (a int);\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a cascading drop of a parent may drop a child",
			sql:     "CREATE TABLE child (x int) INHERITS (nokey);\nCREATE TABLE other (x int);\nDROP TABLE nokey CASCADE;",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE other (x int)", "creates table other without a primary key", []int{0}}},
		},
		{
			name:    "a table dropped again is not reported",
			sql:     "CREATE TABLE n (id int);\nDROP TABLE n;",
			targets: one,
		},
		{
			name:    "CREATE TABLE AS, SELECT INTO, an existing table under IF NOT EXISTS, and a partition are not reported",
			sql:     "CREATE TABLE a AS SELECT 1 AS x;\nSELECT 1 AS x INTO b;\nCREATE TABLE IF NOT EXISTS nokey (a int);\nCREATE TABLE pt (id int PRIMARY KEY) PARTITION BY RANGE (id);\nCREATE TABLE pt1 PARTITION OF pt FOR VALUES FROM (1) TO (10);",
			targets: one,
		},
		{
			name:    "a partitioned table is a table",
			sql:     "CREATE TABLE pt (id int) PARTITION BY RANGE (id);",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE pt (id int) PARTITION BY RANGE (id)", "creates table pt without a primary key", []int{0}}},
		},
		{
			name:    "LIKE copies the key only when it includes indexes",
			sql:     "CREATE TABLE a (LIKE t INCLUDING ALL);\nCREATE TABLE b (LIKE t INCLUDING INDEXES);\nCREATE TABLE c (LIKE t);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE c (LIKE t)", "creates table c without a primary key", []int{0}}},
		},
		{
			name:    "a table that had a key and lost it",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE p DROP COLUMN id CASCADE;",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_pkey CASCADE", "removes the primary key of t, and the change adds none back", []int{0}},
				{1, "DROP COLUMN id CASCADE", "removes the primary key of p, and the change adds none back", []int{0}},
			},
		},
		{
			name:    "a key dropped and added back",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE t ADD PRIMARY KEY (id);",
			targets: one,
		},
		{
			name:    "a table without a key before the change is never reported",
			sql:     "ALTER TABLE nokey ADD PRIMARY KEY (a);\nALTER TABLE nokey DROP CONSTRAINT nokey_pkey;",
			targets: one,
		},
		{
			name:    "pg_dump's search path, keys added at the end",
			sql:     "SELECT pg_catalog.set_config('search_path', '', false);\nCREATE TABLE public.n (id integer NOT NULL);\nCREATE TABLE public.m (id integer NOT NULL);\nALTER TABLE ONLY public.n ADD CONSTRAINT n_pkey PRIMARY KEY (id);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE public.m (id integer NOT NULL)", "creates table public.m without a primary key", []int{0}}},
		},
		{
			name:    "an unknown search path does not hide a new table",
			sql:     "SELECT set_config('search_path', 'public', true);\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "set_config of another setting changes nothing",
			sql:     "SELECT set_config('lock_timeout', '1s', false);\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "procedural code anywhere leaves the end unknown",
			sql:     "CREATE TABLE n (id int);\nDO $$ BEGIN ALTER TABLE n ADD PRIMARY KEY (id); END $$;",
			targets: one,
		},
		{
			name:    "a rollback leaves the end unknown",
			sql:     "BEGIN;\nCREATE TABLE n (id int);\nROLLBACK;",
			targets: one,
		},
		{
			name:    "a temporary table is not reported",
			sql:     "CREATE TEMP TABLE n (id int);\nCREATE TABLE m (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE m (id int)", "creates table m without a primary key", []int{0}}},
		},
		{
			name:    "a statement on another table leaves a new one as it is",
			sql:     "CREATE TABLE n (id int);\nALTER TABLE nope ADD COLUMN a int;\nALTER TABLE n ADD COLUMN b int, ADD CONSTRAINT n_b UNIQUE (b);",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "no schema: a change that stands on its own is still known",
			sql:     "CREATE TABLE n (id int);\nCREATE TABLE m (id int);\nALTER TABLE m ADD PRIMARY KEY (id);",
			targets: []review.Target{{}},
			want:    []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runTargetRule(t, review.RequirePrimaryKey, tt.sql, review.Change{}, tt.targets)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(targetFinding{})); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackup(t *testing.T) {
	withArchive := func(db *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
		db.Schemas = append(db.Schemas, schema(backupSchema))
		return db
	}
	on := review.Change{PriorBackup: true, MaxBackupSize: 2 << 20}
	one := []review.Target{{Schema: withArchive(shop()), SessionUser: "alice"}}
	tests := []struct {
		name    string
		sql     string
		change  review.Change
		targets []review.Target
		want    []targetFinding
	}{
		{
			name:    "backup off",
			sql:     "UPDATE t SET n = 1 WHERE id = 1; DELETE FROM t WHERE id = 2;",
			change:  review.Change{MaxBackupSize: 10},
			targets: []review.Target{{Schema: shop()}},
		},
		{
			name:    "over the size limit, with or without a target",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;",
			change:  review.Change{PriorBackup: true, MaxBackupSize: 10},
			targets: []review.Target{{}, {Schema: withArchive(shop())}},
			want:    []targetFinding{{-1, "", "the change is 32 bytes, over the 10-byte limit of prior backup", []int{0, 1}}},
		},
		{
			name:    "no backup schema",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;",
			change:  on,
			targets: []review.Target{{Schema: shop()}, {Schema: withArchive(shop())}},
			want:    []targetFinding{{-1, "", "prior backup needs schema bbdataarchive, which does not exist", []int{0}}},
		},
		{
			name:    "no backup schema, nothing to back up",
			sql:     "ALTER TABLE t ADD COLUMN x int; INSERT INTO t (id) VALUES (1); WITH d AS (DELETE FROM t RETURNING *) SELECT 1;",
			change:  on,
			targets: []review.Target{{Schema: shop()}},
		},
		{
			name:    "UPDATE and DELETE on one table, however it is named",
			sql:     "UPDATE t SET n = 1 WHERE id = 1;\nUPDATE public.t SET n = 2 WHERE id = 2;\nDELETE FROM public.t WHERE id = 3;\nDELETE FROM t WHERE id = 4;\nDELETE FROM s.t WHERE id = 5;\nUPDATE s.t SET id = 1 WHERE id = 6;",
			change:  on,
			targets: one,
			want: []targetFinding{
				{2, "DELETE FROM public.t WHERE id = 3", "prior backup cannot back up both UPDATE and DELETE on public.t", []int{0}},
				{5, "UPDATE s.t SET id = 1 WHERE id = 6", "prior backup cannot back up both UPDATE and DELETE on s.t", []int{0}},
			},
		},
		{
			name:    "SET search_path moves the lookup",
			sql:     "UPDATE t SET id = 1 WHERE id = 1;\nSET search_path = s;\nDELETE FROM t WHERE id = 2;",
			change:  on,
			targets: one,
		},
		{
			name:    "the lookup follows the synced path, not the session user",
			sql:     "UPDATE t SET id = 1 WHERE id = 1;\nDELETE FROM s.t WHERE id = 2;",
			change:  on,
			targets: []review.Target{{Schema: withArchive(database(`"$user", s, public`, shop().Schemas...)), SessionUser: "public"}},
			want:    []targetFinding{{1, "DELETE FROM s.t WHERE id = 2", "prior backup cannot back up both UPDATE and DELETE on s.t", []int{0}}},
		},
		{
			name:    "a table the change creates",
			sql:     "CREATE TABLE n (id int);\nUPDATE n SET id = 1 WHERE id = 2;\nCREATE TEMP TABLE tmp (id int);\nDELETE FROM tmp WHERE id = 1;\nCREATE TABLE s.m (id int);\nDELETE FROM s.m WHERE id = 1;",
			change:  on,
			targets: one,
			want: []targetFinding{
				{1, "UPDATE n SET id = 1 WHERE id = 2", "prior backup runs before the change, when n does not exist yet", []int{0}},
				{3, "DELETE FROM tmp WHERE id = 1", "prior backup runs before the change, when tmp does not exist yet", []int{0}},
				{5, "DELETE FROM s.m WHERE id = 1", "prior backup runs before the change, when s.m does not exist yet", []int{0}},
			},
		},
		{
			name:    "a table the change moves into the schema",
			sql:     "ALTER TABLE nokey SET SCHEMA s;\nUPDATE s.nokey SET a = 1 WHERE a = 2;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE s.nokey SET a = 1 WHERE a = 2", "prior backup runs before the change, when s.nokey does not exist yet", []int{0}}},
		},
		{
			name:    "a table CREATE SCHEMA creates",
			sql:     "CREATE SCHEMA z CREATE TABLE t (id int) CREATE VIEW v AS SELECT 1 AS a;\nUPDATE z.t SET id = 1 WHERE id = 2;\nDELETE FROM z.v WHERE a = 1;",
			change:  on,
			targets: one,
			want: []targetFinding{
				{1, "UPDATE z.t SET id = 1 WHERE id = 2", "prior backup runs before the change, when z.t does not exist yet", []int{0}},
				{2, "DELETE FROM z.v WHERE a = 1", "prior backup runs before the change, when z.v does not exist yet", []int{0}},
			},
		},
		{
			name:    "an unqualified creation named with its schema",
			sql:     "CREATE TABLE n (id int);\nUPDATE public.n SET id = 1 WHERE id = 2;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE public.n SET id = 1 WHERE id = 2", "prior backup runs before the change, when public.n does not exist yet", []int{0}}},
		},
		{
			name:    "a table the change creates after the statement",
			sql:     "UPDATE n SET id = 1 WHERE id = 2;\nCREATE TABLE n (id int);",
			change:  on,
			targets: one,
			want:    []targetFinding{{0, "UPDATE n SET id = 1 WHERE id = 2", "prior backup runs before the change, when n does not exist yet", []int{0}}},
		},
		{
			name:    "a table of a schema named after its owner",
			sql:     "CREATE SCHEMA AUTHORIZATION bob CREATE TABLE t (id int);\nUPDATE bob.t SET id = 1 WHERE id = 2;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE bob.t SET id = 1 WHERE id = 2", "prior backup runs before the change, when bob.t does not exist yet", []int{0}}},
		},
		{
			name:    "a table the change recreates still exists when the backup runs",
			sql:     "DROP TABLE t;\nCREATE TABLE t (id int);\nDELETE FROM t WHERE id = 1;",
			change:  on,
			targets: one,
		},
		{
			name:    "a table the target lacks that the change does not create is left to the server",
			sql:     "UPDATE nope SET a = 1 WHERE a = 2;\nCREATE TABLE IF NOT EXISTS maybe (a int);\nDELETE FROM maybe WHERE a = 1;",
			change:  on,
			targets: one,
		},
		{
			name:    "no schema, nothing known about the target",
			sql:     "UPDATE t SET n = 1 WHERE id = 1; DELETE FROM t WHERE id = 2;",
			change:  on,
			targets: []review.Target{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runTargetRule(t, review.PriorBackup, tt.sql, tt.change, tt.targets)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(targetFinding{})); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestTargetsScanTheSameStatements reviews one change against targets that
// differ only in session user, so each scans the shared statements at the
// same time. Every finding must list every target: a scan that kept state
// in a statement's nodes would make one target see another's. Run with
// -race to catch the write itself.
func TestTargetsScanTheSameStatements(t *testing.T) {
	sql := `CREATE TABLE a (id serial, b int GENERATED ALWAYS AS (id * 2) STORED, c text DEFAULT 'x' CHECK (c <> ''), p_id int REFERENCES p (id));
CREATE TABLE b (LIKE t INCLUDING ALL, extra int);
CREATE TABLE c (LIKE t);
CREATE INDEX a_c ON a (lower(c)) WHERE id > 0;
CREATE VIEW v AS SELECT a.id, t.code FROM a JOIN t ON t.id = a.id WHERE a.c = 'y';
CREATE FUNCTION f(x int) RETURNS int LANGUAGE sql AS 'SELECT x + 1';
SELECT id INTO d FROM a WHERE id > 1;
CREATE TABLE e AS SELECT f(id) AS id FROM a;
ALTER TABLE a ADD COLUMN z int DEFAULT 1, ALTER COLUMN c SET NOT NULL, ADD CONSTRAINT a_z CHECK (z > 0);
ALTER TABLE t DROP CONSTRAINT t_n_check, DROP CONSTRAINT t_pkey CASCADE, ADD COLUMN w int;
ALTER TABLE c ADD PRIMARY KEY (id);
COMMENT ON TABLE a IS 'a';
CREATE TABLE pt (id int, k int) PARTITION BY RANGE (id);
CREATE TABLE pt1 PARTITION OF pt FOR VALUES FROM (1) TO (10);`
	var targets []review.Target
	var all []int
	for i := range 8 {
		targets = append(targets, review.Target{Schema: shop(), SessionUser: "user" + string(rune('a'+i))})
		all = append(all, i)
	}
	result, err := Review(context.Background(), sql, review.Options{}, targets)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Failures) != 0 {
		t.Fatalf("Failures = %v", result.Failures)
	}
	var walked int
	for _, f := range result.Findings {
		if f.Rule == review.DisallowDropConstraint || f.Rule == review.RequirePrimaryKey {
			walked++
		}
		if diff := cmp.Diff(all, f.Targets); diff != "" {
			t.Errorf("%s %q: targets (-want +got):\n%s", f.Rule, f.Message, diff)
		}
	}
	if walked != 5 {
		t.Errorf("got %d findings from the scan, want 5: %+v", walked, result.Findings)
	}
}
