package analysis

import "testing"

func TestValidateQuery(t *testing.T) {
	for _, tc := range []struct {
		script      string
		readOnly    bool
		returnsData bool
	}{
		{script: "SELECT * FROM t", readOnly: true, returnsData: true},
		{script: "SELECT * FROM t AS OF SYSTEM TIME '-10s'", readOnly: true, returnsData: true},
		{script: "SELECT * FROM t@idx", readOnly: true, returnsData: true},
		{script: "VALUES (1)", readOnly: true, returnsData: true},
		{script: "TABLE t", readOnly: true, returnsData: true},
		{script: "(SELECT 1)", readOnly: true, returnsData: true},
		{script: "SELECT * FROM t FOR UPDATE", readOnly: true, returnsData: true},
		{script: "SELECT * FROM [SHOW TABLES]", readOnly: true, returnsData: true},
		{script: "SHOW RANGES FROM TABLE t", readOnly: true, returnsData: true},
		{script: "SHOW CREATE TABLE t", readOnly: true, returnsData: true},
		{script: "EXPLAIN DELETE FROM t", readOnly: true, returnsData: true},
		{script: "EXPLAIN ANALYZE SELECT 1", readOnly: true, returnsData: true},
		{script: "SET work_mem = '64MB'; SELECT 1", readOnly: true},
		{script: "RESET application_name", readOnly: true},
		{script: "SET TRANSACTION PRIORITY HIGH", readOnly: true},
		{script: "EXPLAIN ANALYZE DELETE FROM t"},
		{script: "WITH x AS (DELETE FROM t RETURNING *) SELECT * FROM x"},
		{script: "SELECT * FROM [DELETE FROM t RETURNING *]"},
		{script: "SELECT (SELECT 1 FROM [UPSERT INTO t VALUES (1) RETURNING 1])"},
		{script: "EXPLAIN ANALYZE SELECT * FROM [INSERT INTO t VALUES (1) RETURNING *]"},
		{script: "SELECT * FROM [EXPLAIN ANALYZE DELETE FROM t]"},
		{script: "SELECT * FROM [EXPLAIN ANALYZE SELECT * FROM [DELETE FROM t RETURNING *]]"},
		{script: "SELECT * FROM [EXPLAIN DELETE FROM t]", readOnly: true, returnsData: true},
		{script: "SELECT * FROM [EXPLAIN ANALYZE SELECT 1]", readOnly: true, returnsData: true},
		{script: "UPSERT INTO t VALUES (1)"},
		{script: "DELETE FROM t"},
		{script: "SELECT 1; CREATE TABLE x (a INT)"},
		{script: "BEGIN"},
	} {
		t.Run(tc.script, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.script)
			if err != nil {
				t.Fatalf("ValidateQuery() error: %v", err)
			}
			if readOnly != tc.readOnly || returnsData != tc.returnsData {
				t.Errorf("ValidateQuery() = (%v, %v), want (%v, %v)", readOnly, returnsData, tc.readOnly, tc.returnsData)
			}
		})
	}

	if _, _, err := ValidateQuery("SELECT * FROM"); err == nil {
		t.Error("ValidateQuery(syntax error) returned no error")
	}
}
