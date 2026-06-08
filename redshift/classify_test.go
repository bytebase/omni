package redshift

import "testing"

func TestGetStatementTypes(t *testing.T) {
	sql := `
CREATE TABLE t (id INT);
INSERT INTO t SELECT 1;
UPDATE t SET id = 2;
DELETE FROM t WHERE id = 2;
MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET id = s.id;
SHOW DATABASES;
COPY t FROM 's3://bucket/file' IAM_ROLE DEFAULT;
UNLOAD ('SELECT * FROM t') TO 's3://bucket/out' IAM_ROLE DEFAULT;
SELECT * INTO copied FROM t;
`
	types, err := GetStatementTypes(sql)
	if err != nil {
		t.Fatalf("GetStatementTypes returned error: %v", err)
	}
	want := []StatementType{
		StatementTypeDDL,
		StatementTypeDML,
		StatementTypeDML,
		StatementTypeDML,
		StatementTypeDML,
		StatementTypeShow,
		StatementTypeCopy,
		StatementTypeUnload,
		StatementTypeDDL,
	}
	if len(types) != len(want) {
		t.Fatalf("got %d types, want %d: %#v", len(types), len(want), types)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("type[%d] = %q, want %q", i, types[i], want[i])
		}
	}
}
