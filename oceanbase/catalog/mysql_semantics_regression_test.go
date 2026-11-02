package catalog

import "testing"

func TestMySQLSemanticsRegression_InvalidColumnAndIndexDDL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		code int
	}{
		{
			name: "create_spatial_index_requires_not_null_column",
			sql: `CREATE TABLE t (g GEOMETRY);
CREATE SPATIAL INDEX sp_g ON t (g);`,
			code: 1252,
		},
		{
			name: "alter_column_set_default_null_rejects_not_null_column",
			sql: `CREATE TABLE t (a INT NOT NULL);
ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL;`,
			code: ErrInvalidDefault,
		},
		{
			name: "alter_column_set_default_rejects_lob_literal_default",
			sql: `CREATE TABLE t (a INT);
ALTER TABLE t CHANGE COLUMN a a BLOB;
ALTER TABLE t ALTER COLUMN a SET DEFAULT 'x';`,
			code: 1101,
		},
		{
			name: "create_table_rejects_duplicate_auto_increment",
			sql:  `CREATE TABLE t (a INT AUTO_INCREMENT KEY, b INT AUTO_INCREMENT KEY);`,
			code: 1075,
		},
		{
			name: "create_table_rejects_not_null_default_null",
			sql:  `CREATE TABLE t (a INT NOT NULL DEFAULT NULL);`,
			code: ErrInvalidDefault,
		},
		{
			name: "create_table_rejects_non_temporal_on_update",
			sql:  `CREATE TABLE t (a INT ON UPDATE NOW());`,
			code: 1294,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			mustExec(t, c, "CREATE DATABASE testdb; USE testdb")
			results, err := c.Exec(tt.sql, &ExecOptions{ContinueOnError: true})
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			catErr := lastExecError(results)
			if catErr == nil {
				t.Fatalf("expected catalog error code %d, got nil", tt.code)
			}
			if catErr.Code != tt.code {
				t.Fatalf("expected catalog error code %d, got %d: %v", tt.code, catErr.Code, catErr)
			}
		})
	}
}

func lastExecError(results []ExecResult) *Error {
	for i := len(results) - 1; i >= 0; i-- {
		if results[i].Error == nil {
			continue
		}
		if catErr, ok := results[i].Error.(*Error); ok {
			return catErr
		}
		return nil
	}
	return nil
}
