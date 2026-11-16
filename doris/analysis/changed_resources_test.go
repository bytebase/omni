package analysis

import (
	"testing"

	"github.com/bytebase/omni/doris/parser"
)

// changesFor scripts the same way a caller does: split, parse each segment,
// and merge one statement's changes at a time.
func changesFor(t *testing.T, script, database string) StatementChanges {
	t.Helper()
	var merged StatementChanges
	for _, statement := range SplitStatements(script) {
		if statement.Empty {
			continue
		}
		file, err := parser.Parse(statement.Text)
		if err != nil {
			t.Fatalf("parse %q: %v", statement.Text, err)
		}
		changes := ExtractChangedResources(file.Stmts[0], database, statement.Text)
		merged.Resources = append(merged.Resources, changes.Resources...)
		merged.DMLStatements = append(merged.DMLStatements, changes.DMLStatements...)
		merged.DMLCount += changes.DMLCount
		merged.InsertCount += changes.InsertCount
	}
	return merged
}

func resourceNames(changes StatementChanges) []string {
	var names []string
	for _, resource := range changes.Resources {
		if resource.DatabaseOnly {
			names = append(names, resource.Database)
			continue
		}
		names = append(names, resource.Database+"."+resource.Table)
	}
	return names
}

func TestExtractChangedResourcesTables(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      []string
	}{
		{"create table takes the request database", "CREATE TABLE t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;", []string{"db.t"}},
		{"a qualified name keeps its own database", "CREATE TABLE other.t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;", []string{"other.t"}},
		{"a catalog qualifier is dropped", "CREATE TABLE cat.other.t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;", []string{"other.t"}},
		{"drop table", "DROP TABLE t;", []string{"db.t"}},
		{"truncate table", "TRUNCATE TABLE t;", []string{"db.t"}},
		{"alter table", "ALTER TABLE t ADD COLUMN b INT;", []string{"db.t"}},
		{"create index names its table", "CREATE INDEX idx ON t (a);", []string{"db.t"}},
		{"drop index names its table", "DROP INDEX idx ON t;", []string{"db.t"}},
		{"insert names its target", "INSERT INTO t VALUES (1);", []string{"db.t"}},
		{"update names its target, not its source", "UPDATE t SET a = 1 WHERE a IN (SELECT a FROM s);", []string{"db.t"}},
		{"delete names its target, not its source", "DELETE FROM t WHERE a = 1;", []string{"db.t"}},
		{"a view records no table", "CREATE VIEW v AS SELECT 1;", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := resourceNames(changesFor(t, tc.statement, "db"))
			if len(got) != len(tc.want) {
				t.Fatalf("resources = %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("resources = %v, want %v", got, tc.want)
					break
				}
			}
		})
	}
}

func TestExtractChangedResourcesDMLAccounting(t *testing.T) {
	for _, tc := range []struct {
		name        string
		statement   string
		wantDML     int
		wantInserts int
		wantTexts   []string
	}{
		{name: "insert values counts rows, not statements", statement: "INSERT INTO t VALUES (1), (2), (3);", wantInserts: 3},
		{name: "insert select is a DML statement", statement: "INSERT INTO t SELECT a FROM s;", wantDML: 1, wantTexts: []string{"INSERT INTO t SELECT a FROM s;"}},
		{name: "update is a DML statement", statement: "UPDATE t SET a = 1 WHERE a = 2;", wantDML: 1, wantTexts: []string{"UPDATE t SET a = 1 WHERE a = 2;"}},
		{name: "delete is a DML statement", statement: "DELETE FROM t WHERE a = 1;", wantDML: 1, wantTexts: []string{"DELETE FROM t WHERE a = 1;"}},
		{name: "DDL is neither", statement: "CREATE TABLE t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;"},
		{
			name:        "a script accumulates both",
			statement:   "INSERT INTO t VALUES (1), (2);\nUPDATE t SET a = 1;\nDELETE FROM t;",
			wantDML:     2,
			wantInserts: 2,
			wantTexts:   []string{"UPDATE t SET a = 1;", "DELETE FROM t;"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changes := changesFor(t, tc.statement, "db")
			if changes.DMLCount != tc.wantDML {
				t.Errorf("DMLCount = %d, want %d", changes.DMLCount, tc.wantDML)
			}
			if changes.InsertCount != tc.wantInserts {
				t.Errorf("InsertCount = %d, want %d", changes.InsertCount, tc.wantInserts)
			}
			if len(changes.DMLStatements) != len(tc.wantTexts) {
				t.Fatalf("DMLStatements = %v, want %v", changes.DMLStatements, tc.wantTexts)
			}
			for i := range changes.DMLStatements {
				if changes.DMLStatements[i] != tc.wantTexts[i] {
					t.Errorf("DMLStatements = %v, want %v", changes.DMLStatements, tc.wantTexts)
				}
			}
		})
	}
}

func TestExtractChangedResourcesDatabaseOnlyTargets(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      []string
	}{
		{"a qualified view records its database", "CREATE VIEW other.v AS SELECT 1;", []string{"other"}},
		{"an unqualified view records nothing", "CREATE VIEW v AS SELECT 1;", nil},
		{"alter view", "ALTER VIEW other.v AS SELECT 1;", []string{"other"}},
		{"drop view", "DROP VIEW other.v;", []string{"other"}},
		{"a qualified materialized view records its database", "DROP MATERIALIZED VIEW other.mv;", []string{"other"}},
		{"drop procedure", "DROP PROCEDURE other.p;", []string{"other"}},
		{"a table DDL records no database-only target", "CREATE TABLE other.t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for _, resource := range changesFor(t, tc.statement, "db").Resources {
				if resource.DatabaseOnly {
					got = append(got, resource.Database)
				}
			}
			if len(got) != len(tc.want) {
				t.Fatalf("database-only targets = %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("database-only targets = %v, want %v", got, tc.want)
					break
				}
			}
		})
	}
}

func TestExtractChangedResourcesAffectsData(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      bool
	}{
		{"drop discards existing rows", "DROP TABLE t;", true},
		{"truncate discards existing rows", "TRUNCATE TABLE t;", true},
		{"alter rewrites existing rows", "ALTER TABLE t ADD COLUMN b INT;", true},
		{"create has no existing rows", "CREATE TABLE t (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1;", false},
		{"insert leaves existing rows alone", "INSERT INTO t VALUES (1);", false},
		{"create index leaves existing rows alone", "CREATE INDEX idx ON t (a);", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resources := changesFor(t, tc.statement, "db").Resources
			if len(resources) != 1 {
				t.Fatalf("resources = %+v, want one", resources)
			}
			if resources[0].AffectsData != tc.want {
				t.Errorf("AffectsData = %v, want %v", resources[0].AffectsData, tc.want)
			}
		})
	}
}

func TestExtractChangedResourcesNilNode(t *testing.T) {
	if changes := ExtractChangedResources(nil, "db", "SELECT 1;"); len(changes.Resources) != 0 || changes.DMLCount != 0 {
		t.Errorf("nil node changes = %+v, want zero", changes)
	}
}
