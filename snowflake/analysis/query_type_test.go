package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/analysis"
	"github.com/bytebase/omni/snowflake/parser"
)

// TestGetQueryType_CallExecuteAreDML locks the classification of
// CALL / EXECUTE IMMEDIATE / EXECUTE TASK as DML (procedures can mutate data;
// the DDL fallback would map them to the wrong ACL bucket).
func TestGetQueryType_CallExecuteAreDML(t *testing.T) {
	for _, sql := range []string{
		"CALL P(1);",
		"EXECUTE IMMEDIATE 'SELECT 1';",
		"EXECUTE TASK T1;",
	} {
		file, err := parser.Parse(sql)
		if err != nil {
			t.Fatalf("parse %q: %v", sql, err)
		}
		if got := analysis.QueryTypeOf(file.Stmts[0]); got != analysis.DML {
			t.Fatalf("QueryTypeOf(%q) = %v, want analysis.DML", sql, got)
		}
	}
}

// TestGetQueryType_UnmappedStaysUnknown locks the fail-closed fallback: parsed
// statements the classifier does not deliberately map (legacy left them
// QueryTypeUnknown, which the SQL service denies) must NOT default to DDL.
func TestGetQueryType_UnmappedStaysUnknown(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want analysis.QueryType
	}{
		{"TRUNCATE TABLE T1;", analysis.QueryTypeUnknown},
		{"GRANT SELECT ON TABLE T1 TO ROLE R1;", analysis.QueryTypeUnknown},
		{"CREATE TABLE T1 (A INT);", analysis.DDL},
		{"ALTER TABLE T1 ADD COLUMN B INT;", analysis.DDL},
		{"DROP TABLE T1;", analysis.DDL},
	} {
		file, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.sql, err)
		}
		if got := analysis.QueryTypeOf(file.Stmts[0]); got != tc.want {
			t.Fatalf("QueryTypeOf(%q) = %v, want %v", tc.sql, got, tc.want)
		}
	}
}

// TestGetQueryType_ShowPipeClassifiesByQuery locks that SHOW ... ->> <query>
// classifies as the piped query's type (so the trailing SELECT is
// permission-checked and masked as a query, not as info-schema metadata).
func TestGetQueryType_ShowPipeClassifiesByQuery(t *testing.T) {
	file, err := parser.Parse("SHOW TABLES ->> SELECT * FROM SENSITIVE_T;")
	if err != nil {
		t.Fatal(err)
	}
	if got := analysis.QueryTypeOf(file.Stmts[0]); got != analysis.Select {
		t.Fatalf("got %v, want analysis.Select", got)
	}
	file, err = parser.Parse("SHOW TABLES;")
	if err != nil {
		t.Fatal(err)
	}
	if got := analysis.QueryTypeOf(file.Stmts[0]); got != analysis.SelectInfoSchema {
		t.Fatalf("got %v, want analysis.SelectInfoSchema", got)
	}
}
