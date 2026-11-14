package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
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
		{"CREATE TABLE other_db.t1 LIKE t2", StatementChanges{Resources: []ChangedResource{table("other_db", "t1", false)}}},
		// CREATE OR REPLACE drops the existing table, so its rows are affected.
		{"CREATE OR REPLACE TABLE other_db.u (a INT)", StatementChanges{Resources: []ChangedResource{table("other_db", "u", true)}}},
		{"CREATE OR REPLACE TABLE t LIKE s", StatementChanges{Resources: []ChangedResource{table("db", "t", true)}}},
		{"CREATE OR REPLACE TABLE t AS SELECT * FROM s", StatementChanges{Resources: []ChangedResource{table("db", "t", true)}}},
		{"ALTER TABLE t DROP COLUMN IF EXISTS b", StatementChanges{Resources: []ChangedResource{table("db", "t", true)}}},
		{"CREATE SEQUENCE other_db.s START WITH 1", StatementChanges{}},
		{"DELETE FROM d WHERE a = 1 RETURNING a", StatementChanges{
			Resources:     []ChangedResource{table("db", "d", false)},
			DMLStatements: []string{"DELETE FROM d WHERE a = 1 RETURNING a;"},
			DMLCount:      1,
		}},
		{"DROP TABLE t1, other_db.t2", StatementChanges{Resources: []ChangedResource{table("db", "t1", true), table("other_db", "t2", true)}}},
		{"TRUNCATE TABLE t", StatementChanges{Resources: []ChangedResource{table("db", "t", true)}}},
		{"ALTER TABLE t1 ADD COLUMN c1 INT", StatementChanges{Resources: []ChangedResource{table("db", "t1", true)}}},
		{"RENAME TABLE t1 TO t2, other_db.a TO b", StatementChanges{Resources: []ChangedResource{
			table("db", "t1", false), table("db", "t2", false), table("other_db", "a", false), table("db", "b", false),
		}}},
		{"CREATE INDEX i ON other_db.t (a)", StatementChanges{Resources: []ChangedResource{table("other_db", "t", false)}}},
		{"DROP INDEX i ON t", StatementChanges{Resources: []ChangedResource{table("db", "t", false)}}},
		{"LOAD DATA INFILE '/tmp/f.csv' INTO TABLE other_db.t (a, b)", StatementChanges{Resources: []ChangedResource{table("other_db", "t", false)}}},
		{"INSERT INTO t1 (c1) VALUES (1), (5)", StatementChanges{Resources: []ChangedResource{table("db", "t1", false)}, InsertCount: 2}},
		// INSERT ... SET inserts one row.
		{"INSERT INTO t1 SET c1 = 7", StatementChanges{Resources: []ChangedResource{table("db", "t1", false)}, InsertCount: 1}},
		{"INSERT INTO t1 SELECT * FROM t2", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false)},
			DMLStatements: []string{"INSERT INTO t1 SELECT * FROM t2;"},
			DMLCount:      1,
		}},
		{"  UPDATE t1 SET c1 = 5  ", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false)},
			DMLStatements: []string{"UPDATE t1 SET c1 = 5;"},
			DMLCount:      1,
		}},
		{"UPDATE t1 JOIN other_db.t2 ON t1.a = t2.a SET t1.b = 1", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false), table("other_db", "t2", false)},
			DMLStatements: []string{"UPDATE t1 JOIN other_db.t2 ON t1.a = t2.a SET t1.b = 1;"},
			DMLCount:      1,
		}},
		{"DELETE FROM t1 WHERE a = 1;", StatementChanges{
			Resources:     []ChangedResource{table("db", "t1", false)},
			DMLStatements: []string{"DELETE FROM t1 WHERE a = 1;"},
			DMLCount:      1,
		}},
		// A derived table names no base table.
		{"UPDATE (SELECT * FROM t) AS x SET a = 1", StatementChanges{
			DMLStatements: []string{"UPDATE (SELECT * FROM t) AS x SET a = 1;"},
			DMLCount:      1,
		}},
		{"CREATE PROCEDURE getUser(id INT) SELECT * FROM hello WHERE uid = id", StatementChanges{}},
		// A qualified non-table object records its database only.
		{"CREATE VIEW other_db.v AS SELECT 1", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"ALTER VIEW other_db.v AS SELECT 1", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"DROP VIEW v, other_db.v", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"CREATE FUNCTION other_db.f() RETURNS INT RETURN 1", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"ALTER PROCEDURE other_db.p COMMENT 'c'", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"DROP PROCEDURE other_db.p", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		{"DROP TRIGGER other_db.tr", StatementChanges{Resources: []ChangedResource{databaseOnly("other_db")}}},
		// An unqualified object falls back to the current database, so it records nothing.
		{"CREATE VIEW v AS SELECT 1", StatementChanges{}},
		// CREATE TRIGGER's name carries no qualifier in the AST.
		{"CREATE TRIGGER other_db.tr BEFORE INSERT ON other_db.t FOR EACH ROW SET NEW.a = 1", StatementChanges{}},
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
