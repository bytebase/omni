package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/tidb/parser"
)

func TestExtractChangedResources(t *testing.T) {
	table := func(database, name string, affectsData bool) ChangedResource {
		return ChangedResource{Database: database, Table: name, AffectsData: affectsData}
	}
	databaseOnly := func(database string) ChangedResource {
		return ChangedResource{Database: database, DatabaseOnly: true}
	}
	tests := []struct {
		statement string
		want      StatementChanges
	}{
		{"CREATE TABLE t1 (c1 INT)", StatementChanges{Resources: []ChangedResource{table("db", "t1", false)}}},
		{"DROP TABLE t1, other_db.t2", StatementChanges{Resources: []ChangedResource{table("db", "t1", true), table("other_db", "t2", true)}}},
		{"TRUNCATE TABLE t", StatementChanges{Resources: []ChangedResource{table("db", "t", true)}}},
		{"ALTER TABLE t1 ADD COLUMN c1 INT", StatementChanges{Resources: []ChangedResource{table("db", "t1", true)}}},
		{"RENAME TABLE t1 TO t2, other_db.a TO b", StatementChanges{Resources: []ChangedResource{
			table("db", "t1", false), table("db", "t2", false), table("other_db", "a", false), table("db", "b", false),
		}}},
		{"CREATE INDEX i ON other_db.t (a)", StatementChanges{Resources: []ChangedResource{table("other_db", "t", false)}}},
		{"DROP INDEX i ON t", StatementChanges{Resources: []ChangedResource{table("db", "t", false)}}},
		{"LOAD DATA INFILE '/tmp/f.csv' INTO TABLE other_db.t", StatementChanges{Resources: []ChangedResource{table("other_db", "t", false)}}},
		{"IMPORT INTO other_db.t FROM '/tmp/f.csv'", StatementChanges{Resources: []ChangedResource{table("other_db", "t", false)}}},
		{"INSERT INTO t1 (c1) VALUES (1), (5)", StatementChanges{Resources: []ChangedResource{table("db", "t1", false)}, InsertCount: 2}},
		{"INSERT INTO t SET c = 1, d = 2", StatementChanges{Resources: []ChangedResource{table("db", "t", false)}, InsertCount: 1}},
		{"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE c = 3", StatementChanges{Resources: []ChangedResource{table("db", "t", false)}, InsertCount: 1}},
		{"INSERT INTO t1 SELECT * FROM t2;", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false)},
			DMLStatements: []string{"INSERT INTO t1 SELECT * FROM t2"},
			DMLCount:      1,
		}},
		{" UPDATE t1 SET c1 = 5 ; ", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false)},
			DMLStatements: []string{"UPDATE t1 SET c1 = 5"},
			DMLCount:      1,
		}},
		{"DELETE t1 FROM t1 JOIN other_db.t2 ON t1.a = t2.a", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false), table("db", "t1", false), table("other_db", "t2", false)},
			DMLStatements: []string{"DELETE t1 FROM t1 JOIN other_db.t2 ON t1.a = t2.a"},
			DMLCount:      1,
		}},
		// EXPLAIN rejects BATCH, so the DML inside it is sampled.
		{"BATCH ON id LIMIT 1000 DELETE FROM other_db.t WHERE c = 1;", StatementChanges{
			Resources:     []ChangedResource{table("other_db", "t", false)},
			DMLStatements: []string{"DELETE FROM other_db.t WHERE c = 1"},
			DMLCount:      1,
		}},
		// The DML's location can fall inside an executable comment, so it is counted without a sample.
		{"BATCH ON id LIMIT 1000 DELETE FROM t WHERE /*T![clustered_index] c = 0 AND */ id < 10000000000", StatementChanges{
			Resources: []ChangedResource{table("db", "t", false)},
			DMLCount:  1,
		}},
		{"BATCH ON id LIMIT 1000 DRY RUN DELETE FROM t WHERE c = 1", StatementChanges{}},
		{"CREATE PROCEDURE p() SELECT 1", StatementChanges{}},
		// A qualified non-table object records its database only.
		{"CREATE VIEW other_db.v AS SELECT 1", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"DROP VIEW v, other_db.v", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"CREATE SEQUENCE other_db.s", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"ALTER SEQUENCE other_db.s RESTART", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"DROP SEQUENCE s, other_db.s", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		// An unqualified object falls back to the current database, so it records nothing.
		{"CREATE VIEW v AS SELECT 1", StatementChanges{}},
		{"SELECT * FROM t", StatementChanges{}},
	}
	for _, tc := range tests {
		t.Run(tc.statement, func(t *testing.T) {
			list, err := parser.Parse(tc.statement)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("got %d statements, want 1", len(list.Items))
			}
			got := ExtractChangedResources(list.Items[0], "db", tc.statement)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestExtractChangedResourcesNilNode(t *testing.T) {
	if got := ExtractChangedResources(nil, "db", ""); !reflect.DeepEqual(got, StatementChanges{}) {
		t.Errorf("got %+v, want nothing", got)
	}
}
