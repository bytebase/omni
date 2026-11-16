package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
}

// facts is a prior backup that can take 100 bytes and finds its database.
func priorBackupFacts(size int) v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: size, MaxSize: 100, Database: "bbdataarchive", DatabaseExists: true}
}

// snapshot is the pre-change schema of one database holding table t.
func priorBackupSchema() v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Name:    "test",
		Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{Name: "t"}}}},
	}}
}

// judgePriorBackup splits sql the way Bytebase does, parses each statement on
// its own, and runs PriorBackupCheck over them.
func judgePriorBackup(t *testing.T, sql string, schema v1review.Schema, facts v1review.PriorBackupFacts) ([]priorBackupFinding, error) {
	t.Helper()
	var stmts []PriorBackupStatement
	for _, segment := range parser.Split(sql) {
		if segment.Empty() {
			continue
		}
		list, err := parser.Parse(segment.Text)
		if err != nil {
			t.Fatalf("Parse %q: %v", segment.Text, err)
		}
		for _, node := range list.Items {
			stmts = append(stmts, PriorBackupStatement{Node: node, Text: segment.Text})
		}
	}
	findings, err := PriorBackupCheck(stmts, schema, facts)
	if err != nil {
		return nil, err
	}
	var got []priorBackupFinding
	for _, f := range findings {
		got = append(got, priorBackupFinding{Code: f.Code, Name: f.Name, Table: f.Table, Value: f.Value, Statement: f.Statement})
	}
	return got, nil
}

func TestPriorBackupCheck(t *testing.T) {
	tooLarge := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupScriptTooLarge, Value: "100"}
	}
	mixedDDL := func(stmt int) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDDLAndDML, Statement: stmt}
	}
	missingDatabase := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupDatabaseMissing, Name: "bbdataarchive"}
	}
	mixedDML := func(table string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: table}
	}
	ddlCases := []struct {
		name string
		ddl  string
	}{
		{name: "create table", ddl: "CREATE TABLE u (a INT)"},
		{name: "create database", ddl: "CREATE DATABASE d"},
		{name: "create index", ddl: "CREATE INDEX i ON t (a)"},
		{name: "create view", ddl: "CREATE VIEW v AS SELECT 1"},
		{name: "create trigger", ddl: "CREATE TRIGGER g BEFORE INSERT ON t FOR EACH ROW SET @a = 1"},
		{name: "create event", ddl: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1"},
		{name: "create function", ddl: "CREATE FUNCTION f() RETURNS INT RETURN 1"},
		{name: "alter table", ddl: "ALTER TABLE t ADD b INT"},
		{name: "alter database", ddl: "ALTER DATABASE test CHARACTER SET utf8mb4"},
		{name: "alter view", ddl: "ALTER VIEW v AS SELECT 1"},
		{name: "alter event", ddl: "ALTER EVENT e DISABLE"},
		{name: "drop table", ddl: "DROP TABLE t"},
		{name: "drop database", ddl: "DROP DATABASE d"},
		{name: "drop index", ddl: "DROP INDEX i ON t"},
		{name: "drop view", ddl: "DROP VIEW v"},
		{name: "drop trigger", ddl: "DROP TRIGGER g"},
		{name: "drop event", ddl: "DROP EVENT e"},
		{name: "drop procedure", ddl: "DROP PROCEDURE p"},
		{name: "rename table", ddl: "RENAME TABLE t TO u"},
		{name: "truncate", ddl: "TRUNCATE TABLE t"},
	}
	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := judgePriorBackup(t, tc.ddl+";\nDELETE FROM u WHERE id = 1;", priorBackupSchema(), priorBackupFacts(10))
			if err != nil {
				t.Fatalf("PriorBackupCheck: %v", err)
			}
			if diff := cmp.Diff([]priorBackupFinding{mixedDDL(0)}, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}

	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		facts  v1review.PriorBackupFacts
		want   []priorBackupFinding
	}{
		{name: "dml only", sql: "UPDATE t SET a = 1 WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "nothing", sql: "SELECT 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "size over the limit", sql: "UPDATE t SET a = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(101), want: []priorBackupFinding{tooLarge()}},
		{name: "size at the limit", sql: "UPDATE t SET a = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(100)},
		{name: "database missing", sql: "UPDATE t SET a = 1 WHERE id = 1;", schema: priorBackupSchema(), facts: v1review.PriorBackupFacts{TotalSize: 10, MaxSize: 100, Database: "bbdataarchive"}, want: []priorBackupFinding{missingDatabase()}},
		{
			name:   "mixed ddl and dml",
			sql:    "CREATE TABLE u (a INT);\nUPDATE t SET a = 1 WHERE id = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDDL(0)},
		},
		{
			name:   "every ddl statement is reported",
			sql:    "CREATE TABLE u (a INT); ALTER TABLE u ADD b INT;\nDELETE FROM t WHERE id = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDDL(0), mixedDDL(1)},
		},
		{
			// V1 reported every DDL statement, whether or not the batch
			// changes rows as well.
			name:   "ddl alone is reported as mixed",
			sql:    "CREATE TABLE u (a INT);\nALTER TABLE u ADD b INT;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDDL(0), mixedDDL(1)},
		},
		{name: "create user is not ddl", sql: "CREATE USER u;\nDELETE FROM t WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "grant is not ddl", sql: "GRANT ALL ON *.* TO u;\nDELETE FROM t WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "mixed dml on one table",
			sql:    "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.t")},
		},
		{name: "one kind of dml twice", sql: "UPDATE t SET a = 1 WHERE id = 1;\nUPDATE t SET b = 1 WHERE id = 2;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "an insert is not grouped",
			sql:    "INSERT INTO t (a) VALUES (1);\nUPDATE t SET a = 1 WHERE id = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{name: "two databases are two tables", sql: "UPDATE db1.t SET a = 1 WHERE id = 1;\nDELETE FROM db2.t WHERE id = 2;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "a qualified and an unqualified reference are one table",
			sql:    "UPDATE test.t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.t")},
		},
		{
			name:  "without a database name",
			sql:   "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			facts: priorBackupFacts(10),
			want:  []priorBackupFinding{mixedDML(".t")},
		},
		{
			name:   "mixed tables in statement order",
			sql:    "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;\nUPDATE u SET a = 1 WHERE id = 1;\nDELETE FROM u WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.t"), mixedDML("test.u")},
		},
		{
			name:   "a join's other table is not changed",
			sql:    "UPDATE t JOIN u ON t.id = u.id SET t.a = 1 WHERE u.b = 2;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.t")},
		},
		{
			name:   "a cte is not a table",
			sql:    "WITH c AS (SELECT 1) UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.t")},
		},
		{
			name:   "every finding in order",
			sql:    "CREATE TABLE u (a INT);\nUPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  v1review.PriorBackupFacts{TotalSize: 101, MaxSize: 100, Database: "bbdataarchive"},
			want:   []priorBackupFinding{tooLarge(), mixedDDL(0), missingDatabase(), mixedDML("test.t")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := judgePriorBackup(t, tt.sql, tt.schema, tt.facts)
			if err != nil {
				t.Fatalf("PriorBackupCheck: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckTableError(t *testing.T) {
	if _, err := judgePriorBackup(t, "DELETE x FROM t;", priorBackupSchema(), priorBackupFacts(10)); err == nil {
		t.Fatal("expected an error for a DELETE target that names no table")
	}
}
