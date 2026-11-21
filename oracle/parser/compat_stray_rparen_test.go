package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestStrayRParenDeferralAndReadings pins the delimiter scan's stray-')'
// deferral, the grammar readings that unblocked it without regressing any
// previously agreeing row, and the accompanying omni_first grammar families,
// together with accepted valid controls beside each refusal.
func TestStrayRParenDeferralAndReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// Stray ')' deferral: report the statement parser's earlier token when
		// it stops before the stray ')', while keeping ')' when the parser
		// reaches the ')'.
		{`ALTER TABLE a MODIFY PARTITIONSET BY RANGE(a) PARTITION BY CONSISTENT HASH(a) PARTITIONS AUTO(PARTITIONSET a VALUES LESS THAN(1) LOB(a) STORE AS LOGGING))`, 21},
		{`CREATE TABLE a PARTITION BY CONSISTENT HASH(a) SUBPARTITION BY RANGE a)`, 43},
		{`CREATE TABLE a ORGANIZATION EXTERNAL TYPE a)`, 37},
		{`CREATE TABLE a ORGANIZATION EXTERNAL(TYPE a LOCATION 's'))`, 13},
		{`SELECT 1 FROM dual)`, 18},

		// Readings that unblocked the stray-')' deferral (13 previously
		// agreeing rows where Oracle reads past the token before ')').
		{`CREATE VECTOR INDEX a ON a a)`, 28},
		{`ALTER MATERIALIZED VIEW a MODIFY LOB a)()`, 38},
		{`ALTER TABLE a ADD PERIOD FOR a)`, 30},
		{`ALTER TABLE a UPGRADE INCLUDING DATA LOB(a) STORE AS LOGGING)`, 60},
		{`ALTER TABLE a UPGRADE INCLUDING DATA VARRAY a STORE AS LOB a(PARTITION a LOB(a) STORE AS LOGGING))`, 97},
		{`CREATE MATERIALIZED VIEW a LOB(a) STORE AS LOGGING) AS SELECT 1 FROM dual`, 50},
		{`CREATE TABLE a LOB(a) STORE AS LOGGING)`, 38},
		{`ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY NOMAXVALUE))`, 63},
		{`CREATE TABLE a XMLTYPE a STORE AS CLOB LOGGING)`, 46},
		{`CREATE TABLE a XMLTYPE a STORE AS CLOB a LOGGING)`, 48},
		{`INSERT ALL INTO a VALUES 1) SELECT 1 FROM dual`, 26},
		{`INSERT WHEN 1 = 1 THEN INTO a VALUES 1) SELECT 1 FROM dual`, 38},
		{`MERGE INTO a USING a ON(1 = 1) WHEN NOT MATCHED THEN INSERT VALUES 1)`, 68},

		// Non-parenthesized VALUES in multi-table INSERT / MERGE must still be
		// rejected even without a trailing ')'.
		{`INSERT ALL INTO a VALUES 1 SELECT 1 FROM dual`, 27},
		{`MERGE INTO a USING a ON(1 = 1) WHEN NOT MATCHED THEN INSERT VALUES 1`, 68},

		// CREATE TABLE XMLTYPE column properties.
		{`CREATE TABLE a XMLTYPE`, 22},
		{`CREATE TABLE a XMLTYPE a`, 15},
		{`CREATE TABLE a XMLTYPE a zzz`, 25},
		{`CREATE TABLE a XMLTYPE a STORE AS CLOB(LOGGING)`, 39},
		{`CREATE TABLE a XMLTYPE a STORE AS CLOB a zzz`, 41},
		{`CREATE TABLE a XMLTYPE a XMLSCHEMA ELEMENT a`, 43},
		{`CREATE TABLE a XMLTYPE a ELEMENT a ALLOW NONSCHEMA`, 35},

		// SUBPARTITION missing BY / TEMPLATE / non-hash template quantity.
		{`CREATE TABLE a PARTITION BY RANGE(a) SUBPARTITION RANGE(a)(PARTITION VALUES LESS THAN(1))`, 50},
		{`CREATE MATERIALIZED VIEW a PARTITION BY LIST(a) SUBPARTITION HASH(a)(PARTITION VALUES(1)) AS SELECT 1 FROM dual`, 61},
		{`CREATE MATERIALIZED VIEW LOG ON a PARTITION BY HASH(a) SUBPARTITION LIST(a) PARTITIONS 1`, 68},
		{`CREATE TABLE a PARTITION BY HASH(a) SUBPARTITION BY RANGE(a) SUBPARTITION 1 PARTITIONS 1`, 74},
		{`CREATE TABLE a PARTITION BY HASH(a) SUBPARTITION BY RANGE(a) SUBPARTITION(SUBPARTITION) PARTITIONS 1`, 73},
		{`CREATE TABLE a PARTITION BY HASH(a) SUBPARTITION BY RANGE(a) SUBPARTITION TEMPLATE 1 PARTITIONS 1`, 85},
		{`ALTER TABLE a MODIFY PARTITION BY HASH(a) SUBPARTITION BY LIST(a) SUBPARTITION TEMPLATE 1 PARTITIONS 1`, 90},

		// Partition maintenance verbs, DEPENDENT TABLES, SPLIT INCLUDING ROWS,
		// and MODIFY [SUB]PARTITION LOB.
		{`ALTER TABLE a SPLIT a AT(1)`, 20},
		{`ALTER TABLE a MERGE a, a`, 20},
		{`ALTER TABLE a EXCHANGE a WITH TABLE a`, 23},
		{`ALTER MATERIALIZED VIEW a EXCHANGE a WITH TABLE a`, 26},
		{`ALTER MATERIALIZED VIEW a MOVE a`, 31},
		{`ALTER MATERIALIZED VIEW a RENAME a TO a`, 33},
		{`ALTER MATERIALIZED VIEW a SET()`, 29},
		{`ALTER TABLE a SPLIT PARTITION a AT(1) DEPENDENT(a(PARTITION))`, 47},
		{`ALTER TABLE a SPLIT PARTITION a AT(1) INCLUDING WHERE 1 = 1`, 48},
		{`ALTER TABLE a MOVE PARTITION a INCLUDING WHERE 1 = 1`, 31},
		{`ALTER TABLE a MODIFY PARTITION a LOB a)(RETENTION)`, 37},
		{`ALTER MATERIALIZED VIEW a MODIFY SUBPARTITION a LOB a)()`, 52},

		// ASSOCIATE STATISTICS ... USING ... WITH.
		{`ASSOCIATE STATISTICS WITH FUNCTIONS a USING a WITH SYSTEM MANAGED STORAGE TABLES`, 51},
		{`ASSOCIATE STATISTICS WITH FUNCTIONS a USING a WITH USER MANAGED STORAGE TABLES`, 51},

		// ILM ADD POLICY and INMEMORY.
		{`ALTER TABLE a ILM ADD POLICY ROW COMPRESS SEGMENT AFTER 1 DAY OF CREATION`, 33},
		{`ALTER TABLE a ILM ADD POLICY COLUMN COMPRESS SEGMENT AFTER 1 DAY OF CREATION`, 36},
		{`ALTER TABLE a ILM ADD POLICY MODIFY NO MEMCOMPRESS SEGMENT AFTER 1 DAY OF CREATION`, 36},
		{`ALTER MATERIALIZED VIEW a ADD PARTITION ILM ADD POLICY COMPRESS ROW AFTER 1 DAY OF CREATION`, 68},
		{`CREATE TABLE a ILM ADD POLICY NO SEGMENT ON a`, 33},
		{`ALTER MATERIALIZED VIEW a INMEMORY DUPLICATE FOR SERVICE a`, 49},
		{`CREATE MATERIALIZED VIEW a INMEMORY(a) INMEMORY,(a) INMEMORY AS SELECT 1 FROM dual`, 48},
		{`CREATE MATERIALIZED VIEW a INMEMORY(a) INMEMORY, a) INMEMORY AS SELECT 1 FROM dual`, 49},

		// CREATE PLUGGABLE DATABASE, ALTER DISKGROUP REPAIR, SET.
		{`CREATE PLUGGABLE DATABASE a USER a IDENTIFIED BY a`, 33},
		{`CREATE PLUGGABLE DATABASE a AS PROXY FROM a@a HOST 's'`, 51},
		{`CREATE PLUGGABLE DATABASE a AS PROXY FROM a@a PORT 1`, 51},
		{`ALTER DISKGROUP a REPAIR DISKS ALL`, 25},
		{`SET`, 3},
		{`SET a`, 4},
		{`SET READ ONLY`, 4},

		// CREATE/ALTER TABLE/MVIEW JSON, LOB, ROW, READ, TABLESPACE, ENABLE,
		// DEFAULT ON, virtual AS, ENCRYPT NO.
		{`CREATE TABLE a JSON a) STORE AS`, 20},
		{`CREATE TABLE a JSON(a) STORE AS zzz`, 15},
		{`CREATE TABLE a ROW STORE`, 15},
		{`CREATE TABLE a ROW MOVEMENT`, 19},
		{`CREATE TABLE a ROW COMPRESS`, 19},
		{`ALTER TABLE a ROW MOVEMENT`, 18},
		{`CREATE TABLE a READ`, 19},
		{`CREATE TABLE a TABLESPACE`, 25},
		{`CREATE TABLE a ENABLE LOGICAL`, 29},
		{`CREATE TABLE a ENABLE LOGICAL REPLICATION zzz`, 42},
		{`CREATE TABLE a DISABLE LOGICAL REPLICATION`, 15},
		{`ALTER TABLE a ENABLE PRIMARY KEY USING a`, 39},
		{`ALTER TABLE a ENABLE PRIMARY KEY EXCEPTIONS a`, 44},
		{`CREATE TABLE a(a DEFAULT ON 1)`, 28},
		{`CREATE TABLE a(a GENERATED BY DEFAULT ON AS IDENTITY)`, 41},
		{`CREATE TABLE a(a AS 1))`, 20},
		{`CREATE TABLE a(a ENCRYPT NO)`, 27},
		{`ALTER TABLE a MODIFY(a ENCRYPT NO zzz SALT)`, 34},

		// SELECT/DELETE/UPDATE APPLY alias, TABLE collection, inline analytic view ADD.
		{`SELECT * FROM a APPLY a`, 22},
		{`DELETE a FROM a APPLY a`, 22},
		{`UPDATE a SET a = 1 FROM a APPLY a`, 32},
		{`SELECT * FROM TABLE 1)`, 20},
		{`DELETE a FROM TABLE 1)`, 20},
		{`SELECT * FROM ANALYTIC VIEW(USING a ADD(a))`, 39},
		{`DELETE a FROM ANALYTIC VIEW(USING a ADD(a))`, 39},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// Controls: valid Oracle forms beside those refusals.
	controls := []string{
		`CREATE VECTOR INDEX vi ON t tbl (vec) ORGANIZATION INMEMORY NEIGHBOR GRAPH`,
		`ALTER TABLE t ADD PERIOD FOR valid_time`,
		`ALTER TABLE t ADD PERIOD FOR valid_time (start_col, end_col)`,
		`ALTER TABLE t ADD (PERIOD FOR valid_time (start_col, end_col))`,
		`ALTER TABLE t ADD PERIOD VARCHAR2(10)`,
		`CREATE TABLE t (c CLOB) LOB (c) STORE AS logging_seg`,
		`CREATE TABLE t (c CLOB) LOB (c) STORE AS LOGGING`,
		`CREATE TABLE t (id NUMBER GENERATED ALWAYS AS IDENTITY NOMAXVALUE)`,
		`ALTER TABLE t MODIFY (id GENERATED ALWAYS AS IDENTITY START WITH 10 INCREMENT BY 2)`,
		`CREATE TABLE t (x XMLTYPE) XMLTYPE COLUMN x STORE AS CLOB (CACHE)`,
		`CREATE TABLE t (x XMLTYPE) XMLTYPE x STORE AS SECUREFILE BINARY XML`,
		`CREATE TABLE t (j JSON) JSON (j) STORE AS (CACHE)`,
		`CREATE TABLE t (a NUMBER) ENABLE LOGICAL REPLICATION`,
		`CREATE TABLE t (a NUMBER) DISABLE LOGICAL REPLICATION`,
		`SELECT * FROM dual APPLY`,
		`SELECT * FROM a CROSS APPLY (SELECT 1 FROM dual)`,
		`SELECT * FROM a OUTER APPLY (SELECT 1 FROM dual)`,
		`SELECT * FROM TABLE(CAST(NULL AS t_tab))`,
		`INSERT ALL INTO t VALUES (1) SELECT 1 FROM dual`,
		`MERGE INTO t USING s ON (t.id = s.id) WHEN NOT MATCHED THEN INSERT VALUES (s.id)`,
	}
	for _, sql := range controls {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}

	// Verify AST fidelity for CreateVectorIndexStmt.TableAlias and unparenthesized ADD PERIOD.
	list, err := Parse(`CREATE VECTOR INDEX vi ON t tbl (vec) ORGANIZATION INMEMORY NEIGHBOR GRAPH`)
	if err != nil {
		t.Fatalf("Parse vector index with alias: %v", err)
	}
	vi, ok := list.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateVectorIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateVectorIndexStmt, got %T", list.Items[0].(*ast.RawStmt).Stmt)
	}
	if vi.TableAlias != "TBL" {
		t.Fatalf("expected TableAlias TBL, got %q", vi.TableAlias)
	}
	if rendered := ast.NodeToString(vi); !strings.Contains(rendered, `:table_alias "TBL"`) {
		t.Fatalf("unexpected NodeToString: %s", rendered)
	}

	atList, err := Parse(`ALTER TABLE t ADD PERIOD FOR valid_time (start_col, end_col)`)
	if err != nil {
		t.Fatalf("Parse bare ADD PERIOD: %v", err)
	}
	atStmt := atList.Items[0].(*ast.RawStmt).Stmt.(*ast.AlterTableStmt)
	cmd := atStmt.Actions.Items[0].(*ast.AlterTableCmd)
	if cmd.Action != ast.AT_ADD_PERIOD || cmd.ColumnName != "VALID_TIME" {
		t.Fatalf("unexpected AlterTableCmd: %+v", cmd)
	}
}
