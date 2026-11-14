package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// TestDataMetricFunctionDescDrop covers DESC and DROP FUNCTION with a data
// metric function's TABLE ( ... ) signature.
// https://docs.snowflake.com/en/sql-reference/sql/desc-function-dmf
// https://docs.snowflake.com/en/sql-reference/sql/drop-function-dmf
func TestDataMetricFunctionDescDrop(t *testing.T) {
	sql := "DROP FUNCTION IF EXISTS governance.dmfs.count_positive_numbers(TABLE(NUMBER, NUMBER), TABLE(VARCHAR))"
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	drop := file.Stmts[0].(*ast.DropStmt)
	if !drop.IfExists || drop.HasArgs || len(drop.TableArgs) != 2 || len(drop.TableArgs[0].Types) != 2 ||
		drop.Loc.End != len(sql) {
		t.Fatalf("drop: %+v", drop)
	}

	for _, sql := range []string{
		"DESC FUNCTION governance.dmfs.count_positive_numbers(TABLE(NUMBER, NUMBER, NUMBER))",
		"DESCRIBE FUNCTION IF EXISTS f(TABLE(INT), TABLE(INT))",
	} {
		file, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		desc := file.Stmts[0].(*ast.DescribeStmt)
		if desc.TableArgs == nil || desc.Signature != nil || desc.Loc.End != len(sql) {
			t.Errorf("%q: %+v", sql, desc)
		}
	}

	for _, sql := range []string{
		// A DMF takes at most two tables.
		"DROP FUNCTION f(TABLE(INT), TABLE(INT), TABLE(INT))",
		// Each TABLE lists one or more types.
		"DROP FUNCTION f(TABLE())",
		"DESC FUNCTION f(TABLE(INT,))",
		// TABLE arguments are a function's only.
		"DROP PROCEDURE p(TABLE(INT))",
		"DESC PROCEDURE p(TABLE(INT))",
		// The signature does not mix TABLE and plain types.
		"DROP FUNCTION f(TABLE(INT), INT)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}
