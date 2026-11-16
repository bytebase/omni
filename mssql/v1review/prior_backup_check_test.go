package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

// priorBackupFinding is a prior-backup finding as the tests compare it: every
// field a finding can set, the range as the text it anchors.
type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Schema    string
	Table     string
	Value     string
	Statement int
	Text      string
}

// priorBackupFacts is a prior backup that can take 100 bytes and finds its
// database.
func priorBackupFacts(size int) v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: size, MaxSize: 100, Database: "bbdataarchive", DatabaseExists: true}
}

// priorBackupSchema is the pre-change snapshot of database test.
func priorBackupSchema() v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: "test"}}
}

// judgePriorBackup parses sql the way Bytebase does and runs PriorBackupCheck
// over its statements.
func judgePriorBackup(t *testing.T, sql string, schema v1review.Schema, facts v1review.PriorBackupFacts) []priorBackupFinding {
	t.Helper()
	var findings []priorBackupFinding
	for _, f := range PriorBackupCheck(parseStatements(t, sql), schema, facts) {
		findings = append(findings, priorBackupFinding{f.Code, f.Name, f.Schema, f.Table, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return findings
}

func TestPriorBackupCheck(t *testing.T) {
	tooLarge := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupScriptTooLarge, Value: "100"}
	}
	ddl := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupDDLInBatch}
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
		{name: "alter table", ddl: "ALTER TABLE t ADD b INT"},
		{name: "drop table", ddl: "DROP TABLE u"},
		{name: "truncate", ddl: "TRUNCATE TABLE t"},
		{name: "create index", ddl: "CREATE INDEX i ON t (a)"},
		{name: "create view", ddl: "CREATE VIEW v AS SELECT 1 AS a"},
		{name: "create schema", ddl: "CREATE SCHEMA s"},
		{name: "alter schema", ddl: "ALTER SCHEMA s TRANSFER dbo.t"},
		{name: "create sequence", ddl: "CREATE SEQUENCE sq START WITH 1"},
		{name: "alter database", ddl: "ALTER DATABASE test SET RECOVERY SIMPLE"},
	}
	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			got := judgePriorBackup(t, tc.ddl+";\nDELETE FROM t WHERE id = 1;", priorBackupSchema(), priorBackupFacts(10))
			if diff := cmp.Diff([]priorBackupFinding{ddl()}, got); diff != "" {
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
		{
			name:   "database missing ends the check",
			sql:    "CREATE TABLE u (a INT);\nUPDATE t SET a = 1 WHERE id = 2;\nDELETE FROM t WHERE id = 3;",
			schema: priorBackupSchema(),
			facts:  v1review.PriorBackupFacts{TotalSize: 101, MaxSize: 100, Database: "bbdataarchive"},
			want:   []priorBackupFinding{tooLarge(), missingDatabase()},
		},
		{
			name:   "ddl is reported once for the batch",
			sql:    "CREATE TABLE u (a INT);\nALTER TABLE u ADD b INT;\nDELETE FROM t WHERE id = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{ddl()},
		},
		{
			// V1 reported the DDL whether or not the batch changes rows.
			name:   "ddl alone is reported as mixed",
			sql:    "CREATE TABLE u (a INT);\nDROP TABLE u;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{ddl()},
		},
		{name: "create user is not ddl", sql: "CREATE USER u;\nDELETE FROM t WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "grant is not ddl", sql: "GRANT SELECT ON t TO u;\nDELETE FROM t WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "mixed dml on one table",
			sql:    "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t")},
		},
		{name: "one kind of dml twice", sql: "UPDATE t SET a = 1 WHERE id = 1;\nUPDATE t SET b = 1 WHERE id = 2;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "an insert is not grouped",
			sql:    "INSERT INTO t (a) VALUES (1);\nUPDATE t SET a = 1 WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a qualified and an unqualified reference are one table",
			sql:    "UPDATE test.dbo.t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t")},
		},
		{
			name:   "two databases are two tables",
			sql:    "UPDATE db1.dbo.t SET a = 1 WHERE id = 1;\nDELETE FROM db2.dbo.t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "two schemas are two tables",
			sql:    "UPDATE dbo.t SET a = 1 WHERE id = 1;\nDELETE FROM other.t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "an update through an alias changes the physical table",
			sql:    "UPDATE x SET a = 1 FROM t AS x WHERE x.id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t")},
		},
		{
			name:   "a delete through an alias changes the physical table",
			sql:    "DELETE x FROM t AS x WHERE x.id = 1;\nUPDATE t SET a = 1 WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t")},
		},
		{
			name:   "a joined alias is followed",
			sql:    "UPDATE x SET a = 1 FROM t AS x JOIN u AS y ON x.id = y.id;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t")},
		},
		{
			// V1 followed the alias to the physical table with the defaults
			// it had: an alias on another database's table keeps that
			// database.
			name:   "an alias on another database keeps that database",
			sql:    "UPDATE x SET a = 1 FROM db2.dbo.t AS x WHERE x.id = 1;\nDELETE FROM db2.dbo.t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("db2.dbo.t")},
		},
		{
			name:   "an alias V1 could not resolve keeps the name it was given",
			sql:    "UPDATE x SET a = 1 FROM u AS y WHERE y.id = 1;\nDELETE FROM x WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.x")},
		},
		{
			// Prior backup never copies a temp table, so two kinds of DML on
			// one cannot conflict.
			name:   "a temp table is not grouped",
			sql:    "UPDATE #t SET a = 1;\nDELETE FROM #t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a temp table reached through an alias is not grouped",
			sql:    "UPDATE x SET a = 1 FROM #t AS x WHERE x.id = 1;\nDELETE FROM #t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a table variable is not grouped",
			sql:    "UPDATE @tv SET a = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:  "without a database name",
			sql:   "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			facts: priorBackupFacts(10),
			want:  []priorBackupFinding{mixedDML(".dbo.t")},
		},
		{
			name:   "mixed tables in statement order",
			sql:    "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM u WHERE id = 2;\nUPDATE u SET a = 1 WHERE id = 3;\nDELETE FROM t WHERE id = 4;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("test.dbo.t"), mixedDML("test.dbo.u")},
		},
		{
			name:   "every finding in order",
			sql:    "CREATE TABLE u (a INT);\nUPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{ddl(), mixedDML("test.dbo.t")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, tt.schema, tt.facts)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
