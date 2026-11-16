package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// priorBackupFinding is a prior-backup finding as the tests compare it.
type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
}

// priorBackupFacts is a prior backup that can take 100 bytes and finds its
// database.
func priorBackupFacts() v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: 100, MaxSize: 100, Database: "bbdataarchive", DatabaseExists: true}
}

// judgePriorBackup parses sql the way Bytebase does and runs PriorBackupCheck
// over the statements.
func judgePriorBackup(t *testing.T, sql string, facts v1review.PriorBackupFacts) []priorBackupFinding {
	t.Helper()
	stmts, err := pg.Parse(sql)
	if err != nil {
		t.Fatalf("Parse %q: %v", sql, err)
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
		if stmt.AST == nil {
			continue
		}
		nodes = append(nodes, stmt.AST)
	}
	var findings []priorBackupFinding
	for _, f := range PriorBackupCheck(nodes, facts) {
		findings = append(findings, priorBackupFinding{Code: f.Code, Name: f.Name, Table: f.Table, Value: f.Value, Statement: f.Statement})
	}
	return findings
}

func TestPriorBackupCheckDDLForms(t *testing.T) {
	// Every DDL form of the V1 advisor's list, one per node type, so a type
	// dropped from the list fails here.
	ddl := []struct{ name, sql string }{
		{"create table", "CREATE TABLE t (a int)"},
		{"create table as", "CREATE TABLE t AS SELECT 1"},
		{"create domain", "CREATE DOMAIN d AS int"},
		{"create extension", "CREATE EXTENSION plpgsql"},
		{"create fdw", "CREATE FOREIGN DATA WRAPPER w"},
		{"create foreign server", "CREATE SERVER s FOREIGN DATA WRAPPER w"},
		{"create foreign table", "CREATE FOREIGN TABLE ft (a int) SERVER s"},
		{"create function", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1'"},
		{"create role", "CREATE ROLE r"},
		{"create user", "CREATE USER u"},
		{"create policy", "CREATE POLICY p ON t USING (true)"},
		{"create language", "CREATE LANGUAGE plpgsql"},
		{"create schema", "CREATE SCHEMA s"},
		{"create sequence", "CREATE SEQUENCE sq"},
		{"create subscription", "CREATE SUBSCRIPTION sub CONNECTION 'host=x' PUBLICATION p"},
		{"create statistics", "CREATE STATISTICS st ON a, b FROM t"},
		{"create tablespace", "CREATE TABLESPACE ts LOCATION '/x'"},
		{"create transform", "CREATE TRANSFORM FOR int LANGUAGE sql (FROM SQL WITH FUNCTION f(internal), TO SQL WITH FUNCTION f(internal))"},
		{"create trigger", "CREATE TRIGGER tr AFTER INSERT ON t EXECUTE FUNCTION f()"},
		{"create event trigger", "CREATE EVENT TRIGGER et ON ddl_command_start EXECUTE FUNCTION f()"},
		{"create user mapping", "CREATE USER MAPPING FOR u SERVER s"},
		{"create database", "CREATE DATABASE d"},
		{"create access method", "CREATE ACCESS METHOD am TYPE INDEX HANDLER h"},
		{"create publication", "CREATE PUBLICATION p"},
		{"create operator class", "CREATE OPERATOR CLASS oc FOR TYPE int USING btree AS OPERATOR 1 <"},
		{"create operator family", "CREATE OPERATOR FAMILY of USING btree"},
		{"create cast", "CREATE CAST (int AS text) WITH INOUT"},
		{"create conversion", "CREATE CONVERSION c FOR 'a' TO 'b' FROM f"},
		{"create range", "CREATE TYPE r AS RANGE (subtype = int)"},
		{"create composite type", "CREATE TYPE ct AS (a int)"},
		{"create enum", "CREATE TYPE e AS ENUM ('a')"},
		{"create index", "CREATE INDEX i ON t (a)"},
		{"create view", "CREATE VIEW v AS SELECT 1"},
		{"create collation", "CREATE COLLATION c (locale = 'x')"},
		{"create operator", "CREATE OPERATOR + (leftarg = int, rightarg = int, function = f)"},
		{"alter event trigger", "ALTER EVENT TRIGGER et DISABLE"},
		{"alter collation", "ALTER COLLATION c REFRESH VERSION"},
		{"alter database", "ALTER DATABASE d CONNECTION LIMIT 5"},
		{"alter database set", "ALTER DATABASE d SET x = 1"},
		{"alter default privileges", "ALTER DEFAULT PRIVILEGES GRANT SELECT ON TABLES TO r"},
		{"alter domain", "ALTER DOMAIN d SET NOT NULL"},
		{"alter enum", "ALTER TYPE e ADD VALUE 'b'"},
		{"alter extension", "ALTER EXTENSION e UPDATE"},
		{"alter extension contents", "ALTER EXTENSION e ADD TABLE t"},
		{"alter fdw", "ALTER FOREIGN DATA WRAPPER w OPTIONS (x '1')"},
		{"alter foreign server", "ALTER SERVER s OPTIONS (x '1')"},
		{"alter function", "ALTER FUNCTION f() IMMUTABLE"},
		{"alter object depends", "ALTER FUNCTION f() DEPENDS ON EXTENSION e"},
		{"alter object schema", "ALTER TABLE t SET SCHEMA s"},
		{"alter owner", "ALTER SCHEMA s OWNER TO r"},
		{"alter operator", "ALTER OPERATOR + (int, int) SET (x = 1)"},
		{"alter type", "ALTER TYPE ct SET (x = 1)"},
		{"alter policy", "ALTER POLICY p ON t USING (false)"},
		{"alter sequence", "ALTER SEQUENCE sq RESTART"},
		{"alter system", "ALTER SYSTEM SET x = 1"},
		{"alter table", "ALTER TABLE t ADD COLUMN b int"},
		{"alter tablespace options", "ALTER TABLESPACE ts SET (x = 1)"},
		{"alter statistics", "ALTER STATISTICS st SET STATISTICS 100"},
		{"alter publication", "ALTER PUBLICATION p SET TABLE t"},
		{"alter role set", "ALTER ROLE r SET x = 1"},
		{"alter role", "ALTER ROLE r LOGIN"},
		{"alter subscription", "ALTER SUBSCRIPTION sub SET (x = 1)"},
		{"alter ts configuration", "ALTER TEXT SEARCH CONFIGURATION c ADD MAPPING FOR a WITH b"},
		{"alter ts dictionary", "ALTER TEXT SEARCH DICTIONARY d (x = 1)"},
		{"alter user mapping", "ALTER USER MAPPING FOR u SERVER s OPTIONS (x '1')"},
		{"alter operator family", "ALTER OPERATOR FAMILY of USING btree ADD OPERATOR 1 < (int, int)"},
		{"drop table", "DROP TABLE t"},
		{"drop owned", "DROP OWNED BY r"},
		{"drop subscription", "DROP SUBSCRIPTION sub"},
		{"drop tablespace", "DROP TABLESPACE ts"},
		{"drop role", "DROP ROLE r"},
		{"drop user mapping", "DROP USER MAPPING FOR u SERVER s"},
		{"drop database", "DROP DATABASE d"},
		{"truncate", "TRUNCATE t"},
		{"rename", "ALTER TABLE t RENAME TO t2"},
		{"comment", "COMMENT ON TABLE t IS 'x'"},
		{"reindex", "REINDEX TABLE t"},
		{"cluster", "CLUSTER t USING i"},
		{"refresh materialized view", "REFRESH MATERIALIZED VIEW mv"},
		{"rule", "CREATE RULE r AS ON INSERT TO t DO NOTHING"},
		{"security label", "SECURITY LABEL ON TABLE t IS 'x'"},
		{"reassign owned", "REASSIGN OWNED BY r TO r2"},
	}
	want := []priorBackupFinding{{Code: v1review.PriorBackupMixedDDLAndDML}}
	for _, tt := range ddl {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(want, judgePriorBackup(t, tt.sql, priorBackupFacts())); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}

	notDDL := []struct{ name, sql string }{
		{"insert", "INSERT INTO t VALUES (1)"},
		{"update", "UPDATE t SET a = 1"},
		{"delete", "DELETE FROM t"},
		{"select", "SELECT 1"},
		{"merge", "MERGE INTO t USING u ON t.a = u.a WHEN MATCHED THEN UPDATE SET a = 1"},
		{"grant", "GRANT SELECT ON t TO r"},
		{"revoke", "REVOKE SELECT ON t FROM r"},
		{"set role", "SET ROLE r"},
		{"begin", "BEGIN"},
		{"commit", "COMMIT"},
		{"copy", "COPY t TO STDOUT"},
	}
	for _, tt := range notDDL {
		t.Run("not ddl "+tt.name, func(t *testing.T) {
			if diff := cmp.Diff([]priorBackupFinding(nil), judgePriorBackup(t, tt.sql, priorBackupFacts())); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckDataModifyingCTE(t *testing.T) {
	cte := func(value string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupDataModifyingCTE, Value: value}
	}
	tests := []struct {
		name string
		sql  string
		want []priorBackupFinding
	}{
		{
			name: "a select CTE is not data-modifying",
			sql:  "WITH c AS (SELECT 1) UPDATE t SET a = 1 FROM c",
			want: nil,
		},
		{
			name: "an update",
			sql:  "WITH d AS (DELETE FROM orders RETURNING order_id) UPDATE t SET a = 1 WHERE a IN (SELECT order_id FROM d)",
			want: []priorBackupFinding{cte(`"d"`)},
		},
		{
			name: "a delete",
			sql:  "WITH d AS (DELETE FROM orders RETURNING order_id) DELETE FROM t WHERE a IN (SELECT order_id FROM d)",
			want: []priorBackupFinding{cte(`"d"`)},
		},
		{
			name: "an insert",
			sql:  "WITH d AS (DELETE FROM orders RETURNING order_id) INSERT INTO t (a) SELECT order_id FROM d",
			want: []priorBackupFinding{cte(`"d"`)},
		},
		{
			name: "a select",
			sql:  "WITH d AS (DELETE FROM orders RETURNING order_id) SELECT order_id FROM d",
			want: []priorBackupFinding{cte(`"d"`)},
		},
		{
			name: "a merge",
			sql:  "WITH d AS (DELETE FROM orders RETURNING order_id) MERGE INTO t USING d ON t.a = d.order_id WHEN MATCHED THEN DO NOTHING",
			want: []priorBackupFinding{cte(`"d"`)},
		},
		{
			name: "two modifying CTEs in order",
			sql:  "WITH s AS (SELECT 1), i AS (INSERT INTO orders VALUES (1) RETURNING order_id), u AS (UPDATE orders SET order_id = 2 RETURNING order_id) DELETE FROM t WHERE a IN (SELECT order_id FROM i)",
			want: []priorBackupFinding{cte(`"i", "u"`)},
		},
		{
			name: "a CTE nested in a subquery is not read",
			sql:  "WITH outer_cte AS (SELECT id FROM (WITH d AS (DELETE FROM orders RETURNING order_id) SELECT order_id AS id FROM d) x) SELECT 1",
			want: nil,
		},
		{
			name: "one finding per statement",
			sql:  "WITH a AS (DELETE FROM orders RETURNING order_id) SELECT 1; WITH b AS (DELETE FROM orders RETURNING order_id) SELECT 2",
			want: []priorBackupFinding{cte(`"a"`), {Code: v1review.PriorBackupDataModifyingCTE, Value: `"b"`, Statement: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, judgePriorBackup(t, tt.sql, priorBackupFacts())); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckMixedDML(t *testing.T) {
	mixed := func(table string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: table}
	}
	tests := []struct {
		name string
		sql  string
		want []priorBackupFinding
	}{
		{
			name: "update then delete",
			sql:  "UPDATE t SET a = 1; DELETE FROM t",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "qualified and unqualified group together",
			sql:  "UPDATE public.t SET a = 1; DELETE FROM t",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "an unqualified name resolves to public",
			sql:  "UPDATE t SET a = 1; DELETE FROM public.t WHERE a = 1",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "another schema is another table",
			sql:  "UPDATE s.t SET a = 1; DELETE FROM public.t",
			want: nil,
		},
		{
			name: "another catalog is the same table",
			sql:  "UPDATE db1.public.t SET a = 1; DELETE FROM db2.public.t",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "an insert is not grouped",
			sql:  "INSERT INTO t VALUES (1); UPDATE t SET a = 1",
			want: nil,
		},
		{
			name: "the same kind twice is not mixed",
			sql:  "DELETE FROM t; DELETE FROM t WHERE a = 1",
			want: nil,
		},
		{
			name: "an alias does not change the table",
			sql:  "UPDATE t AS x SET a = 1; DELETE FROM t AS y WHERE y.a = 2",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "a FROM clause does not add a table",
			sql:  "UPDATE t SET a = u.a FROM u WHERE t.a = u.a; DELETE FROM t",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "three statements on one table",
			sql:  "UPDATE t SET a = 1; DELETE FROM t WHERE a = 1; UPDATE t SET a = 2",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "tables in the order first named",
			sql:  "UPDATE b SET a = 1; DELETE FROM b; UPDATE a SET a = 1; DELETE FROM a",
			want: []priorBackupFinding{mixed("public.b"), mixed("public.a")},
		},
		{
			name: "an unquoted name folds to lower case",
			sql:  "UPDATE T SET a = 1; DELETE FROM t",
			want: []priorBackupFinding{mixed("public.t")},
		},
		{
			name: "quoted names keep their case",
			sql:  `UPDATE "T" SET a = 1; DELETE FROM "T"`,
			want: []priorBackupFinding{mixed("public.T")},
		},
		{
			name: "quoted and unquoted names differ",
			sql:  `UPDATE T SET a = 1; DELETE FROM "T"`,
			want: nil,
		},
		{
			name: "a quoted lower-case name matches",
			sql:  `UPDATE "t" SET a = 1; DELETE FROM t`,
			want: []priorBackupFinding{mixed("public.t")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, judgePriorBackup(t, tt.sql, priorBackupFacts())); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckDatabase(t *testing.T) {
	missing := priorBackupFinding{Code: v1review.PriorBackupDatabaseMissing, Name: "bbdataarchive"}
	ddl := priorBackupFinding{Code: v1review.PriorBackupMixedDDLAndDML}
	mixed := priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: "public.t"}

	absent := priorBackupFacts()
	absent.DatabaseExists = false

	tests := []struct {
		name  string
		sql   string
		facts v1review.PriorBackupFacts
		want  []priorBackupFinding
	}{
		{
			name:  "the database exists",
			sql:   "UPDATE t SET a = 1; DELETE FROM t",
			facts: priorBackupFacts(),
			want:  []priorBackupFinding{mixed},
		},
		{
			name:  "the database is missing",
			sql:   "UPDATE t SET a = 1",
			facts: absent,
			want:  []priorBackupFinding{missing},
		},
		{
			name:  "a missing database does not end the check",
			sql:   "CREATE TABLE u (a int); UPDATE t SET a = 1; DELETE FROM t",
			facts: absent,
			want:  []priorBackupFinding{ddl, missing, mixed},
		},
		{
			name:  "DDL and CTE findings come in statement order",
			sql:   "CREATE TABLE u (a int);\nWITH d AS (DELETE FROM orders RETURNING order_id) UPDATE t SET a = 1;\nDELETE FROM t",
			facts: absent,
			want: []priorBackupFinding{
				ddl,
				{Code: v1review.PriorBackupDataModifyingCTE, Value: `"d"`, Statement: 1},
				missing,
				mixed,
			},
		},
		{
			name:  "an empty script still reports the missing database",
			sql:   "",
			facts: absent,
			want:  []priorBackupFinding{missing},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, judgePriorBackup(t, tt.sql, tt.facts)); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckSize(t *testing.T) {
	// V1 had no size cap in this rule: the backup task enforces one, and this
	// rule does not read the size pair.
	facts := priorBackupFacts()
	facts.TotalSize = 3 * 1024 * 1024
	facts.MaxSize = 2 * 1024 * 1024
	if diff := cmp.Diff([]priorBackupFinding(nil), judgePriorBackup(t, "UPDATE t SET a = 1", facts)); diff != "" {
		t.Fatalf("finding mismatch (-want +got):\n%s", diff)
	}
	facts.TotalSize = facts.MaxSize
	if diff := cmp.Diff([]priorBackupFinding(nil), judgePriorBackup(t, "UPDATE t SET a = 1", facts)); diff != "" {
		t.Fatalf("finding mismatch (-want +got):\n%s", diff)
	}
}
