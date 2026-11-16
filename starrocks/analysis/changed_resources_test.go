package analysis

import (
	"testing"

	"github.com/bytebase/omni/starrocks/parser"
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

func TestExtractChangedResourcesInsertIntoFiles(t *testing.T) {
	// INSERT INTO FILES(...) writes to external storage and leaves Target nil,
	// so it records no table while still counting as a DML statement.
	changes := changesFor(t, "INSERT INTO FILES('path'='s3://x', 'format'='parquet') SELECT a FROM s;", "db")
	if len(changes.Resources) != 0 {
		t.Errorf("resources = %+v, want none", changes.Resources)
	}
	if changes.DMLCount != 1 {
		t.Errorf("DMLCount = %d, want 1", changes.DMLCount)
	}
}

func TestExtractChangedResourcesLoadLabel(t *testing.T) {
	// A LOAD LABEL's database qualifier is the database its job runs in, and
	// the data descriptions' target is a bare identifier (both grammars write
	// `INTO TABLE identifier`), so it loads into that database, not the
	// session's (engine-verified: with the session in another database, LOAD
	// LABEL db_b.job (... INTO TABLE t ...) loads db_b.t).
	for _, tc := range []struct {
		name      string
		statement string
		want      []string
	}{
		{
			name:      "a qualified label is the target's database",
			statement: "LOAD LABEL other.job (DATA INFILE('s3://b/f.csv') INTO TABLE t COLUMNS TERMINATED BY ',')",
			want:      []string{"other.t"},
		},
		{
			name:      "an unqualified label keeps the request database",
			statement: "LOAD LABEL job (DATA INFILE('s3://b/f.csv') INTO TABLE t COLUMNS TERMINATED BY ',')",
			want:      []string{"db.t"},
		},
		{
			name:      "each data description is a target",
			statement: "LOAD LABEL other.job (DATA INFILE('s3://b/f.csv') INTO TABLE t, DATA INFILE('s3://b/g.csv') INTO TABLE u)",
			want:      []string{"other.t", "other.u"},
		},
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

func TestExtractChangedResourcesLoadLabelQuotedParts(t *testing.T) {
	// Both pinned engines take a LOAD LABEL as [db.]label and reject a dot in
	// either part (their label check is ^[\-_A-Za-z0-9:]{1,128}$, and a
	// database name cannot hold one). omni still parses the quoted
	// identifiers, so the analysis must read the parser's parts rather than
	// turn a dot inside one part into a qualifier.
	for _, tc := range []struct {
		name      string
		statement string
		want      []string
	}{
		{
			name:      "a quoted label part keeps its dot out of the database",
			statement: "LOAD LABEL other.`job.with.dot` (DATA INFILE('s3://b/f.csv') INTO TABLE t)",
			want:      []string{"other.t"},
		},
		{
			name:      "a quoted unqualified label is not qualified",
			statement: "LOAD LABEL `job.with.dot` (DATA INFILE('s3://b/f.csv') INTO TABLE t)",
			want:      []string{"db.t"},
		},
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

func TestExtractChangedResourcesSwapWith(t *testing.T) {
	// SWAP WITH exchanges the two tables' data, so the partner is a changed
	// table too, and it resolves in the altered table's database
	// (engine-verified: a session in another database still finds it there).
	for _, tc := range []struct {
		name      string
		statement string
		want      []string
	}{
		{
			name:      "swap",
			statement: "ALTER TABLE t SWAP WITH other",
			want:      []string{"db.t", "db.other"},
		},
		{
			name:      "the partner follows the altered table's database",
			statement: "ALTER TABLE other.t SWAP WITH partner",
			want:      []string{"other.t", "other.partner"},
		},
		{
			name:      "an action with no partner adds none",
			statement: "ALTER TABLE t ADD COLUMN b INT",
			want:      []string{"db.t"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changes := changesFor(t, tc.statement, "db")
			got := resourceNames(changes)
			if len(got) != len(tc.want) {
				t.Fatalf("resources = %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("resources = %v, want %v", got, tc.want)
					break
				}
			}
			for _, resource := range changes.Resources {
				if !resource.AffectsData {
					t.Errorf("resource %+v = 0 AffectsData, want a data change", resource)
				}
			}
		})
	}
}
