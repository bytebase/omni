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

func withNoArgFunctionReturning(db *metadata.DatabaseSchemaMetadata, schema, table string) *metadata.DatabaseSchemaMetadata {
	for _, s := range db.Schemas {
		if s.Name == schema {
			s.Functions = append(s.Functions, &metadata.FunctionMetadata{Name: "f", Signature: "f()", DependencyTables: []*metadata.DependencyTable{{Schema: schema, Table: table}}})
		}
	}
	return db
}

func withPlainIndex(db *metadata.DatabaseSchemaMetadata, tableName, index, column string, unique bool) *metadata.DatabaseSchemaMetadata {
	for _, t := range db.Schemas[0].Tables {
		if t.Name == tableName {
			t.Indexes = append(t.Indexes, &metadata.IndexMetadata{Name: index, Expressions: []string{column}, Unique: unique})
		}
	}
	return db
}

func withView(db *metadata.DatabaseSchemaMetadata, schema, view string) *metadata.DatabaseSchemaMetadata {
	for _, s := range db.Schemas {
		if s.Name == schema {
			s.Views = append(s.Views, &metadata.ViewMetadata{Name: view})
		}
	}
	return db
}

func withEventTrigger(db *metadata.DatabaseSchemaMetadata, enabled bool, tags ...string) *metadata.DatabaseSchemaMetadata {
	db.EventTriggers = append(db.EventTriggers, &metadata.EventTriggerMetadata{Name: "et", Event: "DDL_COMMAND_END", Tags: tags, Enabled: enabled})
	return db
}

func withEvent(db *metadata.DatabaseSchemaMetadata, event string, tags ...string) *metadata.DatabaseSchemaMetadata {
	db.EventTriggers = append(db.EventTriggers, &metadata.EventTriggerMetadata{Name: "et", Event: event, Tags: tags, Enabled: true})
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
			name:    "dropping a constraint the snapshot lacks withholds the statement",
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
			name: "an unqualified new foreign key references the table the path finds",
			sql:  "CREATE TABLE c (x int REFERENCES p (id));\nALTER TABLE public.p DROP CONSTRAINT p_pkey;",
			targets: []review.Target{{Schema: database("s, public",
				schema("public", withPrimaryKey(table("p", "id integer"), "p_pkey", "id")),
				schema("s", withPrimaryKey(table("p", "id integer"), "p_pkey", "id")),
			)}},
			want: []targetFinding{{1, "DROP CONSTRAINT p_pkey", "drops primary key p_pkey of public.p", []int{0}}},
		},
		{
			name:    "a no-op CREATE TABLE IF NOT EXISTS adds no foreign key",
			sql:     "CREATE TABLE IF NOT EXISTS public.nokey (a text REFERENCES t (code));\nALTER TABLE t DROP CONSTRAINT t_code_key;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_code_key", "drops unique constraint t_code_key of t", []int{0}}},
		},
		{
			name:    "an event trigger on function DDL",
			sql:     "CREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: []review.Target{{Schema: withEventTrigger(shop(), true, "CREATE FUNCTION"), SessionUser: "alice"}},
		},
		{
			name:    "DROP COLUMN IF EXISTS of a missing column changes nothing",
			sql:     "ALTER TABLE t DROP COLUMN IF EXISTS nope;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "a constraint the snapshot lacks does not end the scan",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_nn;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "a foreign key a cascade took is not there to drop",
			sql:     "ALTER TABLE p DROP CONSTRAINT p_pkey CASCADE;\nALTER TABLE t DROP CONSTRAINT t_p_fk;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT p_pkey CASCADE", "drops primary key p_pkey of p", []int{0}}},
		},
		{
			name:    "an event trigger on SELECT INTO in a set operation",
			sql:     "SELECT 1 AS a INTO x UNION SELECT 2;\nALTER TABLE t DROP CONSTRAINT t_n_check;",
			targets: []review.Target{{Schema: withEventTrigger(shop(), true, "SELECT INTO"), SessionUser: "alice"}},
		},
		{
			name:    "ADD CONSTRAINT USING INDEX of a missing index refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE USING INDEX missing_idx;",
			targets: one,
		},
		{
			name:    "ADD CONSTRAINT USING INDEX of a constraint's index refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE USING INDEX t_code_key;",
			targets: one,
		},
		{
			name:    "ADD CONSTRAINT USING INDEX of a unique index of the table",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_n_key UNIQUE USING INDEX t_n_uidx;",
			targets: []review.Target{{Schema: withPlainIndex(shop(), "t", "t_n_uidx", "n", true), SessionUser: "alice"}},
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "ADD CONSTRAINT USING INDEX of an index that is not unique refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT t_n_key UNIQUE USING INDEX t_n_uidx;",
			targets: []review.Target{{Schema: withPlainIndex(shop(), "t", "t_n_uidx", "n", false), SessionUser: "alice"}},
		},
		{
			name:    "a key named after a table refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT nokey UNIQUE (n);",
			targets: one,
		},
		{
			name:    "a key named after another table's index refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT p_pkey UNIQUE (n);",
			targets: one,
		},
		{
			name:    "a key named after the one the statement drops",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_code_key, ADD CONSTRAINT t_code_key UNIQUE (code);",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_code_key", "drops unique constraint t_code_key of t", []int{0}}},
		},
		{
			name:    "an added foreign key to a column twice refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id, n) REFERENCES p (id, id);",
			targets: one,
		},
		{
			name:    "an added foreign key to a missing column refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id) REFERENCES p (nope);",
			targets: one,
		},
		{
			name:    "an added foreign key to a table without a primary key refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (n) REFERENCES nokey;",
			targets: one,
		},
		{
			name:    "an added foreign key to columns without a unique index refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (code) REFERENCES p (note);",
			targets: one,
		},
		{
			name:    "an added foreign key of another column count refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id, n) REFERENCES p (id);",
			targets: one,
		},
		{
			name:    "an added foreign key to a missing table refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id) REFERENCES missing_table;",
			targets: one,
		},
		{
			name:    "an added foreign key to a key the table has",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT good FOREIGN KEY (p_id) REFERENCES p (id), ADD CONSTRAINT good2 FOREIGN KEY (n) REFERENCES p;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			// The scan does not follow the keys the change makes.
			name:    "an added foreign key to a unique index the change made withholds the statement's findings",
			sql:     "CREATE UNIQUE INDEX ON p (note);\nALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT good FOREIGN KEY (code) REFERENCES p (note);\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "an added foreign key to a table the change made withholds the statement's findings",
			sql:     "CREATE TABLE k (id int);\nALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id) REFERENCES k (id);",
			targets: one,
		},
		{
			name:    "an added foreign key to its own table without a unique key refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (n) REFERENCES t (n);",
			targets: one,
		},
		{
			name:    "an added foreign key to its own table's key the statement adds",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT good FOREIGN KEY (n) REFERENCES t (n), ADD CONSTRAINT t_n_key UNIQUE (n);",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "an added foreign key to its own table's key the statement drops refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_code_key, ADD CONSTRAINT bad FOREIGN KEY (code) REFERENCES t (code);",
			targets: one,
		},
		{
			name:    "the server adds a column before a constraint written ahead of it",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE (x), ADD COLUMN x int;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "an added foreign key to a key the change dropped refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_p_fk;\nALTER TABLE p DROP CONSTRAINT p_pkey;\nALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT bad FOREIGN KEY (p_id) REFERENCES p (id);",
			targets: one,
			want: []targetFinding{
				{0, "DROP CONSTRAINT t_p_fk", "drops foreign key t_p_fk of t", []int{0}},
				{1, "DROP CONSTRAINT p_pkey", "drops primary key p_pkey of p", []int{0}},
			},
		},
		{
			name:    "an added foreign key to its own table",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT self FOREIGN KEY (p_id) REFERENCES t (id);",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_n_check", "drops check constraint t_n_check of t", []int{0}}},
		},
		{
			name:    "an added key on a missing column refuses the statement",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD CONSTRAINT replacement UNIQUE (nope);",
			targets: one,
		},
		{
			name:    "a no-op column drop frees no constraint name",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_code_key, DROP COLUMN IF EXISTS nope, ADD CONSTRAINT t_n_check CHECK (n > 0);",
			targets: one,
		},
		{
			name:    "a dropped column takes its table's foreign key on it",
			sql:     "ALTER TABLE t DROP COLUMN p_id;\nALTER TABLE p DROP CONSTRAINT p_pkey;",
			targets: one,
			want:    []targetFinding{{1, "DROP CONSTRAINT p_pkey", "drops primary key p_pkey of p", []int{0}}},
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
			sql:     "CREATE SCHEMA z CREATE TABLE t (id int) CREATE TABLE k (id int PRIMARY KEY) CREATE TABLE z.q (id int);",
			targets: one,
			want: []targetFinding{
				{0, "CREATE TABLE t (id int)", "creates table z.t without a primary key", []int{0}},
				{0, "CREATE TABLE z.q (id int)", "creates table z.q without a primary key", []int{0}},
			},
		},
		{
			name:    "CREATE SCHEMA with an element in another schema is refused",
			sql:     "CREATE SCHEMA z CREATE TABLE t (id int) CREATE TABLE s.q (id int);",
			targets: one,
		},
		{
			name:    "CREATE SCHEMA with two elements of one name is refused",
			sql:     "CREATE SCHEMA z CREATE TABLE x (id int) CREATE VIEW x AS SELECT 1;",
			targets: one,
		},
		{
			name:    "CREATE SCHEMA with an element named like a generated name",
			sql:     "CREATE SCHEMA z CREATE TABLE x (id int UNIQUE) CREATE TABLE x_id_key (id int);",
			targets: one,
		},
		{
			name:    "a cascading drop takes a view the change made",
			sql:     "CREATE VIEW w AS SELECT * FROM nokey;\nDROP TABLE nokey CASCADE;\nCREATE TABLE w (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE w (id int)", "creates table w without a primary key", []int{0}}},
		},
		{
			name:    "a cascading drop takes the views reading a view the change made",
			sql:     "CREATE VIEW w AS SELECT * FROM nokey;\nCREATE VIEW w2 AS SELECT * FROM public.w;\nDROP TABLE nokey CASCADE;\nCREATE TABLE w2 (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE w2 (id int)", "creates table w2 without a primary key", []int{0}}},
		},
		{
			name:    "a view renamed after it was made keeps its name uncertain",
			sql:     "CREATE VIEW w AS SELECT * FROM nokey;\nALTER VIEW public.w RENAME TO w2;\nCREATE TABLE w (id int);\nDROP TABLE nokey CASCADE;\nCREATE TABLE w (id int);",
			targets: one,
		},
		{
			name:    "CREATE EXTENSION ends the scan",
			sql:     "CREATE EXTENSION e;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "CREATE EXTENSION IF NOT EXISTS of an installed extension does nothing",
			sql:  "CREATE EXTENSION IF NOT EXISTS pgcrypto;\nCREATE TABLE n (id int);\nCREATE EXTENSION IF NOT EXISTS citext;\nCREATE TABLE m (id int);",
			targets: []review.Target{{Schema: func() *metadata.DatabaseSchemaMetadata {
				db := shop()
				db.Extensions = []*metadata.ExtensionMetadata{{Name: "pgcrypto", Schema: "public"}}
				return db
			}(), SessionUser: "alice"}},
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
			sql:     "DROP FUNCTION f;\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
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
			name:    "a CTE body reads the table its name shadows",
			sql:     "CREATE VIEW v AS WITH t AS (SELECT * FROM t) SELECT * FROM t;\nDROP TABLE public.t;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "a cascading drop takes the views that read its views",
			sql:  "DROP TABLE base CASCADE;\nCREATE TABLE IF NOT EXISTS v2 (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer")},
				Views: []*metadata.ViewMetadata{
					{Name: "v1", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}},
					{Name: "v2", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v1", Column: "id"}}},
				},
			})}},
			want: []targetFinding{{1, "CREATE TABLE IF NOT EXISTS v2 (id int)", "creates table v2 without a primary key", []int{0}}},
		},
		{
			name:    "dropping one of two new overloads leaves the other",
			sql:     "CREATE FUNCTION g(int) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nCREATE FUNCTION g(text) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(int);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "DROP VIEW of a materialized view is refused",
			sql:  "DROP VIEW mv;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:              "public",
				MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "mv"}},
			})}},
		},
		{
			name:    "ADD COLUMN IF NOT EXISTS of an existing column adds no key",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE t ADD COLUMN IF NOT EXISTS id int PRIMARY KEY;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_pkey CASCADE", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name: "a replaced view no longer reads what it read",
			sql:  "CREATE OR REPLACE VIEW v AS SELECT 1 AS id;\nDROP TABLE base;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer")},
				Views:  []*metadata.ViewMetadata{{Name: "v", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}}},
			})}},
			want: []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name: "a sequence another owner took stays",
			sql:  "ALTER SEQUENCE q OWNED BY NONE;\nDROP TABLE base;\nCREATE TABLE IF NOT EXISTS q (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:      "public",
				Tables:    []*metadata.TableMetadata{table("base", "id integer")},
				Sequences: []*metadata.SequenceMetadata{{Name: "q", OwnerTable: "base", OwnerColumn: "id"}},
			})}},
		},
		{
			name:    "a table named like a relation an earlier statement made",
			sql:     "CREATE INDEX ON t (n);\nCREATE TABLE t_n_idx1 (x int);\nCREATE TABLE a (id int CONSTRAINT foo PRIMARY KEY);\nCREATE TABLE foo (x int);",
			targets: one,
		},
		{
			name:    "a table in a missing schema",
			sql:     "CREATE TABLE missing.n (id int);\nCREATE TABLE m (id int);",
			targets: one,
		},
		{
			name:    "IF EXISTS of a missing table leaves the name free",
			sql:     "DROP TABLE IF EXISTS gone;\nCREATE TABLE gone (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE gone (id int)", "creates table gone without a primary key", []int{0}}},
		},
		{
			name:    "DROP FUNCTION IF EXISTS may drop nothing",
			sql:     "CREATE FUNCTION g(integer) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION IF EXISTS g(text);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a view the change created and replaced no longer reads the table",
			sql:     "CREATE VIEW v AS SELECT * FROM nokey;\nCREATE OR REPLACE VIEW v AS SELECT 1 AS a;\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a table in another database",
			sql:     "CREATE TABLE other.public.n (id int);\nCREATE TABLE m (id int);",
			targets: one,
		},
		{
			name:    "an ALTER refused for a subcommand, with only this rule on",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_n_check, ADD COLUMN id int;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "ALTER TABLE of a table the target lacks",
			sql:     "ALTER TABLE nope DROP CONSTRAINT x;\nCREATE TABLE n (id int);\nALTER TABLE IF EXISTS nope DROP CONSTRAINT x;\nCREATE TABLE m (id int);",
			targets: one,
		},
		{
			name:    "ALTER TABLE IF EXISTS of a table the target lacks",
			sql:     "ALTER TABLE IF EXISTS nope DROP CONSTRAINT x;\nCREATE TABLE m (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE m (id int)", "creates table m without a primary key", []int{0}}},
		},
		{
			name:    "CREATE SCHEMA of an existing schema",
			sql:     "CREATE SCHEMA s CREATE TABLE t2 (id int);",
			targets: one,
		},
		{
			name: "a dropped partition takes its own indexes",
			sql:  "DROP TABLE p1;\nCREATE TABLE IF NOT EXISTS p1_idx (id int);",
			targets: []review.Target{{Schema: database("public", schema("public", &metadata.TableMetadata{
				Name:       "pt",
				Columns:    []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadata.TablePartitionMetadata{{Name: "p1", Indexes: []*metadata.IndexMetadata{{Name: "p1_idx"}}}},
			}))}},
			want: []targetFinding{{1, "CREATE TABLE IF NOT EXISTS p1_idx (id int)", "creates table p1_idx without a primary key", []int{0}}},
		},
		{
			name:    "DROP SCHEMA IF EXISTS of a missing schema leaves it missing",
			sql:     "DROP SCHEMA IF EXISTS missing;\nCREATE TABLE missing.n (id int);",
			targets: one,
		},
		{
			name:    "DROP SCHEMA of a missing schema is refused",
			sql:     "DROP SCHEMA missing;\nCREATE TABLE m (id int);",
			targets: one,
		},
		{
			name:    "a skipped ADD COLUMN adds no second key",
			sql:     "ALTER TABLE t ADD COLUMN IF NOT EXISTS id int PRIMARY KEY;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "every column drop of a statement sees the table",
			sql:     "ALTER TABLE t DROP COLUMN code, DROP COLUMN id;",
			targets: one,
			want:    []targetFinding{{0, "DROP COLUMN id", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name:    "EXPLAIN ANALYZE of CREATE MATERIALIZED VIEW makes the view",
			sql:     "EXPLAIN ANALYZE CREATE MATERIALIZED VIEW mv AS SELECT * FROM nokey;\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE SCHEMA IF NOT EXISTS with elements is refused",
			sql:     "CREATE SCHEMA IF NOT EXISTS z CREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE VIEW of a name its schema holds is refused",
			sql:     "CREATE VIEW public.nokey AS SELECT 1 AS a;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE OR REPLACE VIEW of a table is refused",
			sql:     "CREATE OR REPLACE VIEW public.t AS SELECT 1 AS a;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a new function taking a table's row type blocks its drop",
			sql:     "CREATE FUNCTION g(x nokey) RETURNS integer LANGUAGE sql AS 'SELECT 1';\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "an event trigger on ALTER SEQUENCE",
			sql:     "ALTER SEQUENCE q OWNED BY NONE;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEventTrigger(shop(), true, "ALTER SEQUENCE"), SessionUser: "alice"}},
		},
		{
			name:    "a cascading view drop keeps a table's lost key",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nDROP VIEW v CASCADE;",
			targets: []review.Target{{Schema: withView(shop(), "public", "v"), SessionUser: "alice"}},
			want:    []targetFinding{{0, "DROP CONSTRAINT t_pkey CASCADE", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name:    "a definite rename leaves the old name to another table",
			sql:     "CREATE TABLE public.a (id int);\nALTER TABLE public.a RENAME TO b;\nCREATE TABLE public.a (id int PRIMARY KEY);\nDROP TABLE public.a;",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE public.a (id int)", "creates table public.a without a primary key", []int{0}}},
		},
		{
			name:    "a column drop frees the names of its constraints",
			sql:     "ALTER TABLE t DROP COLUMN code, ADD CONSTRAINT t_code_key UNIQUE (n);\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "CREATE SCHEMA IF NOT EXISTS of an existing schema changes nothing",
			sql:     "CREATE SCHEMA IF NOT EXISTS s;\nCREATE TABLE s.t (id int);",
			targets: one,
		},
		{
			name:    "CREATE SEQUENCE of a name its schema holds is refused",
			sql:     "CREATE SEQUENCE public.nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE SEQUENCE IF NOT EXISTS of an existing name does nothing",
			sql:     "CREATE SEQUENCE IF NOT EXISTS public.nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a dropped schema is gone",
			sql:     "DROP SCHEMA s;\nCREATE TABLE s.n (id int);",
			targets: []review.Target{{Schema: database("public", schema("public"), schema("s"))}},
		},
		{
			name:    "a renamed new view keeps its reads under the new name",
			sql:     "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER VIEW public.v RENAME TO w;\nDROP VIEW public.w;\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{4, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a new view dropped with the view it reads",
			sql:     "CREATE VIEW w AS SELECT * FROM v;\nDROP VIEW v, w;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withView(shop(), "public", "v"), SessionUser: "alice"}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a column renamed twice is gone under the middle name",
			sql:     "ALTER TABLE t RENAME COLUMN id TO x;\nALTER TABLE t RENAME COLUMN x TO y;\nALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nALTER TABLE t ADD COLUMN IF NOT EXISTS x int PRIMARY KEY;",
			targets: one,
		},
		{
			name:    "a column the statement adds twice gets one definition",
			sql:     "CREATE TABLE n (id int);\nALTER TABLE n ADD COLUMN x int, ADD COLUMN IF NOT EXISTS x int PRIMARY KEY;",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a new SQL-standard function body reading a table blocks its drop",
			sql:     "CREATE FUNCTION g() RETURNS bigint LANGUAGE SQL RETURN (SELECT count(*) FROM nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "a cascading column drop takes the views of the column and their views",
			sql:  "ALTER TABLE base DROP COLUMN note CASCADE;\nCREATE TABLE w (id int);\nCREATE TABLE IF NOT EXISTS other (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer", "note text")},
				Views: []*metadata.ViewMetadata{
					{Name: "v", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "note"}}},
					{Name: "w", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v", Column: "note"}}},
					{Name: "other", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}},
				},
			})}},
			want: []targetFinding{{1, "CREATE TABLE w (id int)", "creates table w without a primary key", []int{0}}},
		},
		{
			name:    "DROP of a name in another database is refused",
			sql:     "DROP TABLE other.public.nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "DROP FUNCTION with arguments is not matched against a synced function",
			sql:     "DROP FUNCTION f(text);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withFunctionReturning(shop(), "public", "nokey"), SessionUser: "alice"}},
		},
		{
			name:    "DROP FUNCTION of a new function by its arguments",
			sql:     "CREATE FUNCTION g(integer) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(text);\nDROP TABLE nokey;\nCREATE TABLE n (id int);\nDROP FUNCTION IF EXISTS g(integer);",
			targets: one,
		},
		{
			name:    "DROP FUNCTION of a new function, matching arguments",
			sql:     "CREATE FUNCTION g(integer) RETURNS SETOF nokey LANGUAGE sql AS 'SELECT * FROM nokey';\nDROP FUNCTION g(integer);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP FUNCTION of a name the target lacks is refused",
			sql:     "DROP FUNCTION public.missing();\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "DROP FUNCTION of an unqualified name the target lacks is refused",
			sql:     "DROP FUNCTION missing;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "DROP FUNCTION IF EXISTS of a name the target lacks does nothing",
			sql:     "DROP FUNCTION IF EXISTS public.missing();\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP FUNCTION of a name the change gave a function",
			sql:     "CREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION g() RENAME TO h;\nDROP FUNCTION h();\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP PROCEDURE retires the procedure's dependencies",
			sql:     "CREATE PROCEDURE pr(a nokey) LANGUAGE sql AS 'SELECT 1';\nDROP PROCEDURE pr(nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a renamed function keeps its dependencies",
			sql:     "CREATE FUNCTION g(nokey) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION g(nokey) RENAME TO h;\nDROP FUNCTION h(nokey);\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{4, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a renamed function still holds the table under its new name",
			sql:     "CREATE FUNCTION g(nokey) RETURNS int LANGUAGE sql AS 'SELECT 1';\nALTER FUNCTION g(nokey) RENAME TO h;\nDROP TABLE nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a cascading drop takes the synced function on the table",
			sql:     "DROP TABLE s.t CASCADE;\nDROP SCHEMA s;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withFunctionReturning(shop(), "s", "t"), SessionUser: "alice"}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a cascading drop takes a function the change made on the table",
			sql:     "CREATE FUNCTION g(a) RETURNS SETOF b LANGUAGE sql AS 'SELECT * FROM b';\nDROP TABLE a CASCADE;\nDROP TABLE b;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", schema("public", table("a", "id integer"), table("b", "id integer")))}},
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a plain table generates no names",
			sql:     "CREATE TABLE public.a (id int);\nCREATE TABLE public.n_idx (id int);\nDROP TABLE public.a;",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE public.n_idx (id int)", "creates table public.n_idx without a primary key", []int{0}}},
		},
		{
			name:    "a serial column generates a name",
			sql:     "CREATE TABLE a (id serial);\nCREATE TABLE n_idx (id int);",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE a (id serial)", "creates table a without a primary key", []int{0}}},
		},
		{
			name:    "LIKE of a table the change made without a key copies none",
			sql:     "CREATE TABLE public.a (id int);\nCREATE TABLE public.b (LIKE public.a INCLUDING INDEXES);\nDROP TABLE public.a;",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE public.b (LIKE public.a INCLUDING INDEXES)", "creates table public.b without a primary key", []int{0}}},
		},
		{
			name:    "LIKE of a table the change made with a key",
			sql:     "CREATE TABLE public.a (id int PRIMARY KEY);\nCREATE TABLE public.b (LIKE public.a INCLUDING INDEXES);",
			targets: one,
		},
		{
			name:    "CREATE TABLE IF NOT EXISTS in a schema the change created",
			sql:     "CREATE SCHEMA z;\nCREATE TABLE IF NOT EXISTS z.n (id int);\nCREATE TABLE z.k (id int PRIMARY KEY);\nCREATE TABLE IF NOT EXISTS z.k (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE IF NOT EXISTS z.n (id int)", "creates table z.n without a primary key", []int{0}}},
		},
		{
			name: "CREATE TABLE of a type's name is refused",
			sql:  "CREATE TABLE IF NOT EXISTS mood (id int);",
			targets: []review.Target{{Schema: func() *metadata.DatabaseSchemaMetadata {
				db := shop()
				db.Schemas[0].EnumTypes = []*metadata.EnumTypeMetadata{{Name: "mood"}}
				return db
			}(), SessionUser: "alice"}},
		},
		{
			name:    "CREATE TABLE of a type the change made",
			sql:     "CREATE DOMAIN d AS int;\nCREATE TABLE d (id int);",
			targets: one,
		},
		{
			name:    "a table_rewrite event trigger does not fire on ADD COLUMN of a plain column",
			sql:     "ALTER TABLE t ADD COLUMN x int, ADD COLUMN y text[];\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "TABLE_REWRITE"), SessionUser: "alice"}},
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a table_rewrite event trigger fires on ADD COLUMN with a default",
			sql:     "ALTER TABLE t ADD COLUMN x int DEFAULT random()::int;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "TABLE_REWRITE"), SessionUser: "alice"}},
		},
		{
			name:    "a table_rewrite event trigger fires on a column type change",
			sql:     "ALTER TABLE t ALTER COLUMN n TYPE bigint;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "TABLE_REWRITE"), SessionUser: "alice"}},
		},
		{
			name:    "DROP FUNCTION IF EXISTS of a function the change made, by its arguments",
			sql:     "CREATE FUNCTION public.f(public.nokey) RETURNS integer LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION IF EXISTS public.f(public.nokey);\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{3, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a view the change made follows the table it reads through a rename",
			sql:     "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER TABLE public.nokey RENAME TO u;\nDROP TABLE public.u;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "a synced view follows the table it reads through a rename",
			sql:  "ALTER TABLE base RENAME TO u;\nDROP TABLE u;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer")},
				Views:  []*metadata.ViewMetadata{{Name: "v", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}}},
			})}},
		},
		{
			name:    "a renamed table without dependents drops",
			sql:     "ALTER TABLE nokey RENAME TO u;\nDROP TABLE u;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a cascading DROP TABLE IF EXISTS of a missing table settles nothing",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_pkey CASCADE;\nDROP TABLE IF EXISTS public.missing CASCADE;",
			targets: one,
			want:    []targetFinding{{0, "DROP CONSTRAINT t_pkey CASCADE", "removes the primary key of t, and the change adds none back", []int{0}}},
		},
		{
			name:    "CREATE OR REPLACE FUNCTION redefines a synced function without arguments",
			sql:     "CREATE OR REPLACE FUNCTION public.f() RETURNS integer LANGUAGE sql RETURN 1;\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withNoArgFunctionReturning(shop(), "public", "nokey"), SessionUser: "alice"}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "CREATE OR REPLACE FUNCTION with arguments may add an overload",
			sql:     "CREATE OR REPLACE FUNCTION public.f(int) RETURNS integer LANGUAGE sql RETURN 1;\nDROP TABLE public.nokey;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withNoArgFunctionReturning(shop(), "public", "nokey"), SessionUser: "alice"}},
		},
		{
			name:    "a CREATE PROCEDURE tag filter does not fire on CREATE FUNCTION",
			sql:     "CREATE FUNCTION g() RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "DDL_COMMAND_END", "CREATE PROCEDURE"), SessionUser: "alice"}},
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a CREATE PROCEDURE tag filter fires on CREATE PROCEDURE",
			sql:     "CREATE PROCEDURE pr() LANGUAGE sql AS 'SELECT 1';\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "DDL_COMMAND_END", "CREATE PROCEDURE"), SessionUser: "alice"}},
		},
		{
			name:    "CREATE TABLE with a column twice is refused",
			sql:     "CREATE TABLE n (id int, id text);",
			targets: one,
		},
		{
			name:    "a rename to a name the schema holds is refused",
			sql:     "ALTER TABLE t RENAME TO nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a view the change made follows the table it reads to another schema",
			sql:     "CREATE VIEW public.v AS SELECT * FROM public.nokey;\nALTER TABLE public.nokey SET SCHEMA s;\nDROP TABLE s.nokey;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name: "a synced view follows the table it reads to another schema",
			sql:  "ALTER TABLE base SET SCHEMA s;\nDROP TABLE s.base;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer")},
				Views:  []*metadata.ViewMetadata{{Name: "v", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}}},
			}, &metadata.SchemaMetadata{Name: "s"})}},
		},
		{
			name:    "SET SCHEMA to a schema the target lacks is refused",
			sql:     "ALTER TABLE nokey SET SCHEMA nope;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "SET SCHEMA to a schema holding the name is refused",
			sql:     "ALTER TABLE public.t SET SCHEMA s;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "dropping a key frees its index's name",
			sql:     "ALTER TABLE t DROP CONSTRAINT t_code_key;\nCREATE TABLE IF NOT EXISTS t_code_key (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE IF NOT EXISTS t_code_key (id int)", "creates table t_code_key without a primary key", []int{0}}},
		},
		{
			name: "dropping a column frees the sequence it owns",
			sql:  "ALTER TABLE base DROP COLUMN id;\nCREATE TABLE q (x int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:      "public",
				Tables:    []*metadata.TableMetadata{table("base", "id integer", "x integer")},
				Sequences: []*metadata.SequenceMetadata{{Name: "q", OwnerTable: "base", OwnerColumn: "id"}},
			})}},
			want: []targetFinding{{1, "CREATE TABLE q (x int)", "creates table q without a primary key", []int{0}}},
		},
		{
			name:    "a qualified DROP of a table the change made frees its name",
			sql:     "CREATE TABLE public.n (id int PRIMARY KEY);\nDROP TABLE public.n;\nCREATE TABLE public.n (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE public.n (id int)", "creates table public.n without a primary key", []int{0}}},
		},
		{
			name:    "a DROP of a table the change made that a new view reads",
			sql:     "CREATE TABLE public.n (id int PRIMARY KEY);\nCREATE VIEW public.w AS SELECT * FROM public.n;\nDROP TABLE public.n;\nCREATE TABLE public.n (id int);",
			targets: one,
		},
		{
			name:    "a rename of a relation the target lacks is refused",
			sql:     "ALTER TABLE public.missing RENAME TO x;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a rename under IF EXISTS of a relation the target lacks does nothing",
			sql:     "ALTER TABLE IF EXISTS public.missing RENAME TO x;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "SET SCHEMA of a relation the target lacks is refused",
			sql:     "ALTER TABLE public.missing SET SCHEMA s;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a login event trigger fires on no statement",
			sql:     "CREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "LOGIN"), SessionUser: "alice"}},
			want:    []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP SCHEMA of a schema the change put a table in is refused",
			sql:     "CREATE SCHEMA z;\nCREATE TABLE z.x (id int PRIMARY KEY);\nDROP SCHEMA z;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "DROP SCHEMA of a schema the change emptied again",
			sql:     "CREATE SCHEMA z;\nCREATE TABLE z.x (id int PRIMARY KEY);\nDROP TABLE z.x;\nDROP SCHEMA z;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{4, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP SCHEMA of a synced schema a table was renamed in is refused",
			sql:     "ALTER TABLE s.t RENAME TO t2;\nDROP SCHEMA s;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "a foreign key of a table the same DROP drops does not hold it",
			sql:     "CREATE TABLE public.c (id int PRIMARY KEY, n_a int REFERENCES public.k);\nDROP TABLE public.c, public.k;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", schema("public", withPrimaryKey(table("k", "id integer"), "k_pkey", "id")))}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "a foreign key of a table the change made still holds the table it references",
			sql:     "CREATE TABLE public.c (id int PRIMARY KEY, n_a int REFERENCES public.k);\nDROP TABLE public.k;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", schema("public", withPrimaryKey(table("k", "id integer"), "k_pkey", "id")))}},
		},
		{
			name:    "CREATE TABLE INHERITS of a parent the target lacks is refused",
			sql:     "CREATE TABLE child (id int) INHERITS (public.missing);",
			targets: one,
		},
		{
			name:    "CREATE TABLE INHERITS of an index is refused",
			sql:     "CREATE TABLE child (id int) INHERITS (public.t_pkey);",
			targets: one,
		},
		{
			name:    "CREATE TABLE LIKE of a relation the target lacks is refused",
			sql:     "CREATE TABLE child (LIKE public.missing);",
			targets: one,
		},
		{
			name:    "CREATE TABLE INHERITS of a table",
			sql:     "CREATE TABLE child (id int) INHERITS (public.nokey);",
			targets: one,
			want:    []targetFinding{{0, "CREATE TABLE child (id int) INHERITS (public.nokey)", "creates table child without a primary key", []int{0}}},
		},
		{
			name:    "a second CREATE TABLE of a name the change made is refused",
			sql:     "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE TABLE public.x (id int);\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "CREATE TABLE IF NOT EXISTS of a name the change made does nothing",
			sql:     "CREATE TABLE public.x (id int PRIMARY KEY);\nCREATE TABLE IF NOT EXISTS public.x (id int);\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "ALTER TABLE of an index is refused",
			sql:     "ALTER TABLE public.t_pkey DROP CONSTRAINT x;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "TRUNCATE fires no event trigger",
			sql:     "TRUNCATE t;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "DDL_COMMAND_END", "CREATE INDEX"), SessionUser: "alice"}},
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "SET SCHEMA vacates the schema the relation leaves",
			sql:     "ALTER TABLE s.m SET SCHEMA public;\nDROP SCHEMA s;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: database("public", schema("public"), schema("s", withUnique(table("m", "id integer"), "m_key", "id")))}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP SCHEMA after the change dropped the function it made there",
			sql:     "CREATE SCHEMA z;\nCREATE FUNCTION z.f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION z.f();\nDROP SCHEMA z;\nCREATE TABLE n (id int);",
			targets: one,
			want:    []targetFinding{{4, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "DROP SCHEMA of a schema holding a function the change made is refused",
			sql:     "CREATE SCHEMA z;\nCREATE FUNCTION z.f() RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP SCHEMA z;\nCREATE TABLE n (id int);",
			targets: one,
		},
		{
			name:    "an unqualified name the search path resolves elsewhere does not settle a table",
			sql:     "CREATE TABLE public.n (id int);\nSET search_path = s, public;\nALTER TABLE n ADD PRIMARY KEY (id);",
			targets: []review.Target{{Schema: database(`"$user", public`, schema("public"), schema("s", table("n", "id integer")))}},
			want:    []targetFinding{{0, "CREATE TABLE public.n (id int)", "creates table public.n without a primary key", []int{0}}},
		},
		{
			name: "a cascade does not go past a view the change replaced",
			sql:  "CREATE OR REPLACE VIEW v AS SELECT 1 AS id;\nDROP TABLE base CASCADE;\nCREATE TABLE w (id int);",
			targets: []review.Target{{Schema: database("public", &metadata.SchemaMetadata{
				Name:   "public",
				Tables: []*metadata.TableMetadata{table("base", "id integer")},
				Views: []*metadata.ViewMetadata{
					{Name: "v", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "base", Column: "id"}}},
					{Name: "w", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v", Column: "id"}}},
				},
			})}},
		},
		{
			name:    "DROP SCHEMA of a schema that holds a table is refused",
			sql:     "CREATE TABLE n (id int);\nDROP SCHEMA s;",
			targets: one,
		},
		{
			name: "reowning a sequence in one schema leaves another's",
			sql:  "ALTER SEQUENCE s.q OWNED BY NONE;\nDROP TABLE public.base;\nCREATE TABLE public.q (id int);",
			targets: []review.Target{{Schema: database("public",
				&metadata.SchemaMetadata{Name: "public", Tables: []*metadata.TableMetadata{table("base", "id integer")}, Sequences: []*metadata.SequenceMetadata{{Name: "q", OwnerTable: "base", OwnerColumn: "id"}}},
				&metadata.SchemaMetadata{Name: "s", Sequences: []*metadata.SequenceMetadata{{Name: "q"}}},
			)}},
			want: []targetFinding{{2, "CREATE TABLE public.q (id int)", "creates table public.q without a primary key", []int{0}}},
		},
		{
			name:    "an sql_drop event trigger does not fire on an ALTER TABLE that drops nothing",
			sql:     "ALTER TABLE t ADD COLUMN x int, ALTER COLUMN n SET DEFAULT 1;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "SQL_DROP"), SessionUser: "alice"}},
			want:    []targetFinding{{1, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name:    "an sql_drop event trigger fires on DROP DEFAULT",
			sql:     "ALTER TABLE t ALTER COLUMN n DROP DEFAULT;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "SQL_DROP"), SessionUser: "alice"}},
		},
		{
			name:    "an sql_drop event trigger fires on a DROP the scan does not name",
			sql:     "DROP DOMAIN d;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "SQL_DROP", "DROP TABLE"), SessionUser: "alice"}},
		},
		{
			name:    "a tag filter fires on DDL the scan has no tag for",
			sql:     "CREATE POLICY pol ON t USING (true);\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "DDL_COMMAND_END", "CREATE INDEX"), SessionUser: "alice"}},
		},
		{
			name:    "a tag filter matches the tags the scan names",
			sql:     "COMMENT ON TABLE t IS 'x';\nGRANT SELECT ON t TO PUBLIC;\nCREATE TABLE n (id int);",
			targets: []review.Target{{Schema: withEvent(shop(), "DDL_COMMAND_END", "CREATE INDEX"), SessionUser: "alice"}},
			want:    []targetFinding{{2, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
		},
		{
			name: "an sql_drop event trigger does not fire on CREATE TABLE",
			sql:  "CREATE TABLE n (id int);",
			targets: []review.Target{{Schema: func() *metadata.DatabaseSchemaMetadata {
				db := shop()
				db.EventTriggers = []*metadata.EventTriggerMetadata{{Name: "et", Event: "SQL_DROP", Enabled: true}}
				return db
			}(), SessionUser: "alice"}},
			want: []targetFinding{{0, "CREATE TABLE n (id int)", "creates table n without a primary key", []int{0}}},
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
			sql:     "CREATE TABLE n (id int);\nALTER TABLE IF EXISTS nope ADD COLUMN a int;\nALTER TABLE n ADD COLUMN b int, ADD CONSTRAINT n_b UNIQUE (b);",
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
			name:    "a name qualified with the current database",
			sql:     "CREATE TABLE db.public.n (id int);\nUPDATE db.public.n SET id = 1 WHERE id = 2;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE db.public.n SET id = 1 WHERE id = 2", "prior backup runs before the change, when db.public.n does not exist yet", []int{0}}},
		},
		{
			name:    "SELECT INTO on a branch of a set operation",
			sql:     "SELECT 1 AS id INTO n UNION SELECT 2;\nUPDATE n SET id = 3 WHERE id = 1;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE n SET id = 3 WHERE id = 1", "prior backup runs before the change, when n does not exist yet", []int{0}}},
		},
		{
			name:    "EXPLAIN ANALYZE of CREATE TABLE AS",
			sql:     "EXPLAIN ANALYZE CREATE TABLE n AS SELECT 1 AS id;\nUPDATE n SET id = 2 WHERE id = 1;",
			change:  on,
			targets: one,
			want:    []targetFinding{{1, "UPDATE n SET id = 2 WHERE id = 1", "prior backup runs before the change, when n does not exist yet", []int{0}}},
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
