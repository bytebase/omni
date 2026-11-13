package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/v1review"
)

func parse(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return list.Items
}

// anchored renders findings as the tests compare them, and fails on a field
// the keyword and letter-case rules never set.
func anchored(t *testing.T, sql string, got []v1review.Finding) []finding {
	t.Helper()
	var findings []finding
	for _, f := range got {
		if f.Table != "" || f.Pattern != "" || f.Detail != "" || f.Search != "" {
			t.Errorf("finding %+v sets Table, Pattern, Detail or Search", f)
		}
		findings = append(findings, finding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return findings
}

func TestIdentifiers(t *testing.T) {
	// named is a collected name and the text its Loc covers.
	type named struct{ Name, Text string }
	for _, tt := range []struct {
		name string
		sql  string
		want []named
	}{
		{
			name: "table, columns, named constraints, partitions",
			sql:  `CREATE TABLE shop."orderItem" (id NUMBER CONSTRAINT c1 PRIMARY KEY, CONSTRAINT uq1 UNIQUE (id)) PARTITION BY RANGE (id) (PARTITION p1 VALUES LESS THAN (10))`,
			want: []named{
				{"orderItem", `shop."orderItem"`},
				{"UQ1", "CONSTRAINT uq1 UNIQUE (id)"},
				{"ID", "id NUMBER CONSTRAINT c1 PRIMARY KEY"},
				{"C1", "CONSTRAINT c1 PRIMARY KEY"},
				{"P1", "PARTITION p1 VALUES LESS THAN (10)"},
			},
		},
		{
			name: "named NOT NULL constraint and subpartitions not collected",
			sql:  "CREATE TABLE t (a NUMBER CONSTRAINT nn1 NOT NULL) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p1 VALUES LESS THAN (10) (SUBPARTITION sp1))",
			want: []named{
				{"T", "t"},
				{"A", "a NUMBER CONSTRAINT nn1 NOT NULL"},
				{"P1", "PARTITION p1 VALUES LESS THAN (10) (SUBPARTITION sp1)"},
			},
		},
		{
			name: "view columns and check option constraint anchor the statement",
			sql:  "CREATE VIEW v1 (id) AS SELECT 1 AS id FROM dual WITH CHECK OPTION CONSTRAINT vc1",
			want: []named{
				{"V1", "v1"},
				{"ID", "CREATE VIEW v1 (id) AS SELECT 1 AS id FROM dual WITH CHECK OPTION CONSTRAINT vc1"},
				{"VC1", "CREATE VIEW v1 (id) AS SELECT 1 AS id FROM dual WITH CHECK OPTION CONSTRAINT vc1"},
				{"ID", "1 AS id"},
			},
		},
		{
			name: "select aliases, references not collected",
			sql:  "SELECT o.a AS aa, COUNT(*) cnt FROM orders o JOIN (SELECT 1 one FROM dual) sub ON o.id = sub.one",
			want: []named{
				{"AA", "o.a AS aa"},
				{"CNT", "COUNT(*) cnt"},
				{"O", "o"},
				{"ONE", "1 one"},
				{"SUB", "sub"},
			},
		},
		{
			name: "CTE and its columns share the CTE's Loc",
			sql:  "WITH cte1 (c1, c2) AS (SELECT 1, 2 FROM dual) SELECT * FROM cte1",
			want: []named{
				{"CTE1", "cte1 (c1, c2) AS (SELECT 1, 2 FROM dual)"},
				{"C1", "cte1 (c1, c2) AS (SELECT 1, 2 FROM dual)"},
				{"C2", "cte1 (c1, c2) AS (SELECT 1, 2 FROM dual)"},
			},
		},
		{
			name: "a name repeated at one Loc is collected once",
			sql:  "WITH a (a) AS (SELECT 1 FROM dual) SELECT * FROM a",
			want: []named{{"A", "a (a) AS (SELECT 1 FROM dual)"}},
		},
		{name: "rename statement", sql: "RENAME old_name TO new_name", want: []named{{"NEW_NAME", "new_name"}}},
		{name: "alter table rename", sql: "ALTER TABLE t RENAME TO t2", want: []named{{"T2", "RENAME TO t2"}}},
		{name: "alter table rename column", sql: "ALTER TABLE t RENAME COLUMN a TO b", want: []named{{"B", "RENAME COLUMN a TO b"}}},
		{name: "alter table rename constraint", sql: "ALTER TABLE t RENAME CONSTRAINT a TO b", want: []named{{"B", "RENAME CONSTRAINT a TO b"}}},
		{name: "alter table rename partition", sql: "ALTER TABLE t RENAME PARTITION p1 TO p2", want: []named{{"P2", "RENAME PARTITION p1 TO p2"}}},
		{name: "alter table add constraint", sql: "ALTER TABLE t ADD CONSTRAINT ck1 CHECK (a > 0)", want: []named{{"CK1", "CONSTRAINT ck1 CHECK (a > 0)"}}},
		{name: "alter table add column", sql: "ALTER TABLE t ADD (c1 NUMBER)", want: []named{{"C1", "c1 NUMBER"}}},
		{name: "alter table drop constraint not collected", sql: "ALTER TABLE t DROP CONSTRAINT c1"},
		{name: "alter index rename anchors the statement", sql: "ALTER INDEX i1 RENAME TO i2", want: []named{{"I2", "ALTER INDEX i1 RENAME TO i2"}}},
		{name: "alter trigger rename anchors the statement", sql: "ALTER TRIGGER tr1 RENAME TO tr2", want: []named{{"TR2", "ALTER TRIGGER tr1 RENAME TO tr2"}}},
		{name: "create index", sql: "CREATE INDEX shop.idx1 ON t (a)", want: []named{{"IDX1", "shop.idx1"}}},
		{name: "create cluster, its columns not collected", sql: "CREATE CLUSTER cl1 (id NUMBER)", want: []named{{"CL1", "cl1"}}},
		{
			name: "string-named create anchors the statement",
			sql:  "CREATE DATABASE LINK link1 CONNECT TO u IDENTIFIED BY p USING 'x'",
			want: []named{{"LINK1", "CREATE DATABASE LINK link1 CONNECT TO u IDENTIFIED BY p USING 'x'"}},
		},
		{name: "admin create", sql: "CREATE DIRECTORY dir1 AS '/tmp'", want: []named{{"DIR1", "dir1"}}},
		{name: "admin drop not collected", sql: "DROP DIRECTORY dir1"},
		{name: "unnamed create", sql: "CREATE MATERIALIZED VIEW LOG ON t"},
		{
			name: "routine parameters and declarations",
			sql:  "CREATE OR REPLACE PROCEDURE p1 (a IN NUMBER) AS\n  v NUMBER;\n  CURSOR c IS SELECT 1 FROM dual;\n  TYPE t_n IS TABLE OF NUMBER;\n  SUBTYPE st IS NUMBER;\nBEGIN\n  FOR i IN 1..2 LOOP NULL; END LOOP;\n  FORALL j IN 1..2 DELETE FROM t WHERE id = j;\nEND;",
			want: []named{
				{"P1", "p1"},
				{"A", "a IN NUMBER"},
				{"V", "v NUMBER;"},
				{"C", "CURSOR c IS SELECT 1 FROM dual;"},
				{"T_N", "TYPE t_n IS TABLE OF NUMBER;"},
				{"ST", "SUBTYPE st IS NUMBER;"},
				{"I", "FOR i IN 1..2 LOOP NULL; END LOOP;"},
				{"J", "FORALL j IN 1..2 DELETE FROM t WHERE id = j;"},
			},
		},
		{
			name: "nested routines, iterands, statements in a block",
			sql:  "DECLARE\n  PROCEDURE inner_p (x NUMBER) IS BEGIN NULL; END;\nBEGIN\n  <<lbl>>\n  FOR i, j IN PAIRS OF v LOOP\n    INSERT INTO t SELECT 1 one FROM dual;\n  END LOOP;\nEND;",
			want: []named{
				{"INNER_P", "inner_p"},
				{"X", "x NUMBER"},
				{"I", "i"},
				{"J", "j"},
				{"ONE", "1 one"},
			},
		},
		{name: "explain plan", sql: "EXPLAIN PLAN FOR SELECT a AS b FROM t", want: []named{{"B", "a AS b"}}},
		{
			name: "create schema nested table",
			sql:  "CREATE SCHEMA AUTHORIZATION shop CREATE TABLE nested1 (id NUMBER)",
			want: []named{{"NESTED1", "nested1"}, {"ID", "id NUMBER"}},
		},
		{name: "dml without aliases", sql: "UPDATE t SET name = 1 WHERE id = 1"},
		{name: "grant", sql: "GRANT SELECT ON my_table TO admin"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got []named
			for _, stmt := range parse(t, tt.sql) {
				for _, ident := range identifiers(stmt) {
					start, end := span(ident.loc)
					got = append(got, named{ident.name, tt.sql[start:end]})
				}
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("identifiers mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
