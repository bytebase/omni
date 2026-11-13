package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// unqualifiedFinding is a NamingFullyQualified finding as the tests compare
// it.
type unqualifiedFinding struct {
	Code v1review.Code
	Name string
	Text string
}

func runNamingFullyQualified(t *testing.T, sql string, schema v1review.Schema) []unqualifiedFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var findings []unqualifiedFinding
	for _, f := range NamingFullyQualified(list.Items, schema) {
		if f.Table != "" || f.Pattern != "" || f.Detail != "" || f.Search != "" {
			t.Errorf("finding sets more than its code, name, and range: %+v", f)
		}
		findings = append(findings, unqualifiedFinding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingFullyQualified(t *testing.T) {
	unqualified := func(name, text string) unqualifiedFinding {
		return unqualifiedFinding{Code: v1review.NameUnqualified, Name: name, Text: text}
	}
	// A snapshot with tables in two schemas and an external table.
	snapshot := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{
			Name:              "public",
			Tables:            []*metadata.TableMetadata{{Name: "tech_book"}, {Name: "Book"}, {Name: "a.b"}},
			ExternalTables:    []*metadata.ExternalTableMetadata{{Name: "ext"}},
			Views:             []*metadata.ViewMetadata{{Name: "v"}},
			MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "mv"}},
		},
		{Name: "other", Tables: []*metadata.TableMetadata{{Name: "elsewhere"}}},
	}}}
	// A snapshot with no schemas at all, and one with a schema list that is
	// empty rather than absent.
	noSchemas := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{}}
	emptySchemas := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{}}}

	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		want   []unqualifiedFinding
	}{
		{
			name: "relations of the judged statements",
			sql:  "CREATE TABLE t1 (a int); CREATE TABLE t2 PARTITION OF p DEFAULT; CREATE INDEX i ON t3 (a); CREATE TRIGGER trg AFTER INSERT ON t4 EXECUTE FUNCTION f(); CREATE SEQUENCE s5 OWNED BY t.a; ALTER TABLE t6 ADD COLUMN a int; ALTER SEQUENCE s7 RESTART; INSERT INTO t8 SELECT * FROM u; UPDATE t9 SET a = 1 FROM u;",
			want: []unqualifiedFinding{
				unqualified("t1", "CREATE TABLE t1 (a int)"),
				unqualified("t2", "CREATE TABLE t2 PARTITION OF p DEFAULT"),
				unqualified("t3", "CREATE INDEX i ON t3 (a)"),
				unqualified("t4", "CREATE TRIGGER trg AFTER INSERT ON t4 EXECUTE FUNCTION f()"),
				unqualified("s5", "CREATE SEQUENCE s5 OWNED BY t.a"),
				unqualified("t6", "ALTER TABLE t6 ADD COLUMN a int"),
				unqualified("s7", "ALTER SEQUENCE s7 RESTART"),
				unqualified("t8", "INSERT INTO t8 SELECT * FROM u"),
				unqualified("t9", "UPDATE t9 SET a = 1 FROM u"),
			},
		},
		{
			name: "qualified relations pass, with a database too",
			sql:  "CREATE TABLE s.t (a int); CREATE INDEX i ON s.t (a); CREATE TRIGGER trg AFTER INSERT ON s.t EXECUTE FUNCTION f(); CREATE SEQUENCE s.q; ALTER TABLE s.t ADD COLUMN a int; ALTER SEQUENCE s.q RESTART; INSERT INTO db.s.t VALUES (1); UPDATE s.t SET a = 1; ALTER TABLE s.t RENAME TO u; DROP TABLE s.t, db.s.t;",
		},
		{
			name: "alters that share the alter table node and renames that have a relation",
			sql:  "ALTER INDEX i SET (fillfactor = 50); ALTER VIEW v OWNER TO x; ALTER SEQUENCE s OWNER TO x; ALTER INDEX i RENAME TO j; ALTER TABLE t RENAME COLUMN a TO b; ALTER TABLE t RENAME CONSTRAINT a TO b; ALTER VIEW v RENAME TO w; ALTER SEQUENCE s RENAME TO q; ALTER TYPE ty RENAME ATTRIBUTE a TO b;",
			want: []unqualifiedFinding{
				unqualified("i", "ALTER INDEX i SET (fillfactor = 50)"),
				unqualified("v", "ALTER VIEW v OWNER TO x"),
				unqualified("s", "ALTER SEQUENCE s OWNER TO x"),
				unqualified("i", "ALTER INDEX i RENAME TO j"),
				unqualified("t", "ALTER TABLE t RENAME COLUMN a TO b"),
				unqualified("t", "ALTER TABLE t RENAME CONSTRAINT a TO b"),
				unqualified("v", "ALTER VIEW v RENAME TO w"),
				unqualified("s", "ALTER SEQUENCE s RENAME TO q"),
				unqualified("ty", "ALTER TYPE ty RENAME ATTRIBUTE a TO b"),
			},
		},
		{
			name: "a rename the parser keeps no range for anchors nothing",
			sql:  "ALTER TRIGGER trg ON t RENAME TO b;",
			want: []unqualifiedFinding{unqualified("t", "")},
		},
		{
			name: "renames without a relation and other alters not judged",
			sql:  "ALTER SCHEMA a RENAME TO b; ALTER TYPE ty RENAME TO b; ALTER FUNCTION f(int) RENAME TO g; ALTER DOMAIN d RENAME CONSTRAINT a TO b; ALTER DATABASE d RENAME TO e; ALTER TABLE t SET SCHEMA s;",
		},
		{
			name: "drop reports each unqualified object, whatever it drops",
			sql:  "DROP TABLE a, s.b, c; DROP INDEX CONCURRENTLY i; DROP VIEW v; DROP SEQUENCE q; DROP TYPE ty; DROP SCHEMA sch; DROP EXTENSION ext;",
			want: []unqualifiedFinding{
				unqualified("a", "DROP TABLE a, s.b, c"),
				unqualified("c", "DROP TABLE a, s.b, c"),
				unqualified("i", "DROP INDEX CONCURRENTLY i"),
				unqualified("v", "DROP VIEW v"),
				unqualified("q", "DROP SEQUENCE q"),
				unqualified("ty", "DROP TYPE ty"),
				unqualified("sch", "DROP SCHEMA sch"),
				unqualified("ext", "DROP EXTENSION ext"),
			},
		},
		{
			name: "dropped trigger, rule, and policy count as qualified by their table; functions, roles, and databases not judged",
			sql:  "DROP TRIGGER trg ON t; DROP RULE r ON t; DROP POLICY p ON t; DROP FUNCTION f(int); DROP ROLE r; DROP DATABASE d; DROP OWNED BY r;",
		},
		{
			name: "a quoted name holding a dot passes",
			sql:  `CREATE TABLE "a.b" (a int); DROP TABLE "a.b"; INSERT INTO "a.b" VALUES (1);`,
		},
		{
			name:   "select without a snapshot not judged",
			sql:    "SELECT * FROM tech_book;",
			schema: v1review.Schema{},
		},
		{
			name:   "select reports the from relations the snapshot holds as tables, in any schema",
			sql:    `SELECT * FROM tech_book, public.tech_book, nope, ext, elsewhere, v, mv, "Book", book, "TECH_BOOK";`,
			schema: snapshot,
			want: []unqualifiedFinding{
				unqualified("tech_book", `SELECT * FROM tech_book, public.tech_book, nope, ext, elsewhere, v, mv, "Book", book, "TECH_BOOK"`),
				unqualified("ext", `SELECT * FROM tech_book, public.tech_book, nope, ext, elsewhere, v, mv, "Book", book, "TECH_BOOK"`),
				unqualified("elsewhere", `SELECT * FROM tech_book, public.tech_book, nope, ext, elsewhere, v, mv, "Book", book, "TECH_BOOK"`),
				unqualified("Book", `SELECT * FROM tech_book, public.tech_book, nope, ext, elsewhere, v, mv, "Book", book, "TECH_BOOK"`),
			},
		},
		{
			name:   "select matches the snapshot case sensitively whatever the instance",
			sql:    `SELECT * FROM book, "TECH_BOOK";`,
			schema: v1review.Schema{Metadata: snapshot.Metadata, ObjectCaseSensitive: false},
		},
		{
			name:   "select reaches joins and subqueries in from, in order",
			sql:    "SELECT * FROM ONLY tech_book t JOIN (ext CROSS JOIN s.x) ON true, LATERAL (SELECT * FROM elsewhere, (SELECT 1 FROM tech_book) q) sub;",
			schema: snapshot,
			want: []unqualifiedFinding{
				unqualified("tech_book", "SELECT * FROM ONLY tech_book t JOIN (ext CROSS JOIN s.x) ON true, LATERAL (SELECT * FROM elsewhere, (SELECT 1 FROM tech_book) q) sub"),
				unqualified("ext", "SELECT * FROM ONLY tech_book t JOIN (ext CROSS JOIN s.x) ON true, LATERAL (SELECT * FROM elsewhere, (SELECT 1 FROM tech_book) q) sub"),
				unqualified("elsewhere", "SELECT * FROM ONLY tech_book t JOIN (ext CROSS JOIN s.x) ON true, LATERAL (SELECT * FROM elsewhere, (SELECT 1 FROM tech_book) q) sub"),
				unqualified("tech_book", "SELECT * FROM ONLY tech_book t JOIN (ext CROSS JOIN s.x) ON true, LATERAL (SELECT * FROM elsewhere, (SELECT 1 FROM tech_book) q) sub"),
			},
		},
		{
			name:   "select forms that reach the from clause",
			sql:    "TABLE tech_book; (SELECT * FROM tech_book); SELECT * INTO n FROM tech_book; WITH ext AS (SELECT 1) SELECT * FROM ext;",
			schema: snapshot,
			want: []unqualifiedFinding{
				unqualified("tech_book", "TABLE tech_book"),
				unqualified("tech_book", "(SELECT * FROM tech_book)"),
				unqualified("tech_book", "SELECT * INTO n FROM tech_book"),
				unqualified("ext", "WITH ext AS (SELECT 1) SELECT * FROM ext"),
			},
		},
		{
			name:   "select relations outside the from clause not judged",
			sql:    "SELECT * FROM tech_book UNION SELECT * FROM ext; WITH c AS (SELECT * FROM tech_book) SELECT (SELECT 1 FROM tech_book) FROM c WHERE EXISTS (SELECT 1 FROM tech_book); SELECT * FROM tech_book TABLESAMPLE SYSTEM (10); SELECT * FROM (SELECT 1 FROM tech_book UNION SELECT 1 FROM ext) q; SELECT * FROM generate_series(1, 2) tech_book;",
			schema: snapshot,
		},
		{
			name:   "a name holding a dot is still reported in a select",
			sql:    `SELECT * FROM "a.b";`,
			schema: snapshot,
			want:   []unqualifiedFinding{unqualified("a.b", `SELECT * FROM "a.b"`)},
		},
		{
			name:   "snapshot with no schemas reports every unqualified from relation",
			sql:    "SELECT * FROM anything JOIN s.other ON true, more;",
			schema: noSchemas,
			want: []unqualifiedFinding{
				unqualified("anything", "SELECT * FROM anything JOIN s.other ON true, more"),
				unqualified("more", "SELECT * FROM anything JOIN s.other ON true, more"),
			},
		},
		{
			name:   "snapshot with an empty schema list reports none",
			sql:    "SELECT * FROM anything;",
			schema: emptySchemas,
		},
		{
			name:   "other statements read the snapshot not at all",
			sql:    "CREATE TABLE tech_book (a int); INSERT INTO nope VALUES (1); DROP TABLE nope;",
			schema: emptySchemas,
			want: []unqualifiedFinding{
				unqualified("tech_book", "CREATE TABLE tech_book (a int)"),
				unqualified("nope", "INSERT INTO nope VALUES (1)"),
				unqualified("nope", "DROP TABLE nope"),
			},
		},
		{
			name:   "statements not judged",
			sql:    "DELETE FROM tech_book; MERGE INTO tech_book USING ext ON true WHEN MATCHED THEN DELETE; CREATE VIEW v AS SELECT * FROM tech_book; CREATE MATERIALIZED VIEW mv AS SELECT 1; CREATE TABLE t AS SELECT * FROM tech_book; CREATE FOREIGN TABLE ft (a int) SERVER srv; TRUNCATE tech_book; COMMENT ON TABLE tech_book IS 'x'; GRANT SELECT ON tech_book TO r; EXPLAIN SELECT * FROM tech_book; CREATE SCHEMA s CREATE TABLE t (a int); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC INSERT INTO tech_book VALUES (1); END;",
			schema: noSchemas,
		},
		{
			name: "unquoted names fold to lower, quoted names keep their case",
			sql:  `CREATE TABLE Tbl (a int); CREATE TABLE "Tbl" (a int);`,
			want: []unqualifiedFinding{
				unqualified("tbl", "CREATE TABLE Tbl (a int)"),
				unqualified("Tbl", `CREATE TABLE "Tbl" (a int)`),
			},
		},
		{
			name: "anchor excludes a leading comment and the semicolon",
			sql:  "-- c\nCREATE TABLE\n  t (a int);",
			want: []unqualifiedFinding{unqualified("t", "CREATE TABLE\n  t (a int)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runNamingFullyQualified(t, tt.sql, tt.schema)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
