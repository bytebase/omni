package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
}

// priorBackupFacts is a prior backup that can take 100 bytes and finds its
// database.
func priorBackupFacts(size int) v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: size, MaxSize: 100, Database: "bbdataarchive", DatabaseExists: true}
}

// priorBackupSchema is the pre-change schema of one database holding table t.
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
	tooLarge := priorBackupFinding{Code: v1review.PriorBackupScriptTooLarge, Value: "100"}
	mixedDDL := func(stmt int) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDDLAndDML, Statement: stmt}
	}
	missingDatabase := priorBackupFinding{Code: v1review.PriorBackupDatabaseMissing, Name: "bbdataarchive"}
	mixedDML := priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: "test.t"}

	for _, tc := range []struct {
		name string
		ddl  string
	}{
		{name: "create table", ddl: "CREATE TABLE u (a INT)"},
		{name: "create sequence", ddl: "CREATE SEQUENCE s START WITH 1"},
		{name: "alter sequence", ddl: "ALTER SEQUENCE s RESTART"},
		{name: "drop sequence", ddl: "DROP SEQUENCE s"},
		{name: "create or replace table", ddl: "CREATE OR REPLACE TABLE u (a INT)"},
	} {
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
		{name: "mariadb returning", sql: "DELETE FROM t WHERE id = 2 RETURNING id;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "size over the limit", sql: "UPDATE t SET a = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(101), want: []priorBackupFinding{tooLarge}},
		{name: "size at the limit", sql: "UPDATE t SET a = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(100)},
		{name: "database missing", sql: "UPDATE t SET a = 1 WHERE id = 1;", schema: priorBackupSchema(), facts: v1review.PriorBackupFacts{TotalSize: 10, MaxSize: 100, Database: "bbdataarchive"}, want: []priorBackupFinding{missingDatabase}},
		{
			name:   "mixed ddl and dml",
			sql:    "CREATE TABLE u (a INT);\nDELETE FROM t WHERE id = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDDL(0)},
		},
		{name: "create user is not ddl", sql: "CREATE USER u;\nDELETE FROM t WHERE id = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "mixed dml on one table",
			sql:    "UPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML},
		},
		{name: "one kind of dml twice", sql: "UPDATE t SET a = 1 WHERE id = 1;\nUPDATE t SET b = 1 WHERE id = 2;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "two databases are two tables",
			sql:    "UPDATE db1.t SET a = 1 WHERE id = 1;\nDELETE FROM db2.t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a qualified and an unqualified reference are one table",
			sql:    "UPDATE test.t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML},
		},
		{
			name:   "every finding in order",
			sql:    "CREATE SEQUENCE s START WITH 1;\nUPDATE t SET a = 1 WHERE id = 1;\nDELETE FROM t WHERE id = 2;",
			schema: priorBackupSchema(),
			facts:  v1review.PriorBackupFacts{TotalSize: 101, MaxSize: 100, Database: "bbdataarchive"},
			want:   []priorBackupFinding{tooLarge, mixedDDL(0), missingDatabase, mixedDML},
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
	if _, err := judgePriorBackup(t, "UPDATE t SET x.a = 1 WHERE id = 1;", priorBackupSchema(), priorBackupFacts(10)); err == nil {
		t.Fatal("expected an error for an UPDATE whose SET names a table it does not read")
	}
}
