package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func parseOneUndrop(t *testing.T, sql string) *ast.UndropStmt {
	t.Helper()
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	stmt, ok := file.Stmts[0].(*ast.UndropStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.UndropStmt", sql, file.Stmts[0])
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Parse(%q): Loc.End = %d, want %d", sql, stmt.Loc.End, len(sql))
	}
	return stmt
}

// TestUndropMore covers the UNDROP forms undrop.go reads.
func TestUndropMore(t *testing.T) {
	for _, tt := range []struct {
		sql  string
		kind ast.UndropObjectKind
		name string
	}{
		{"UNDROP ACCOUNT my_account", ast.UndropAccount, "MY_ACCOUNT"},
		{"UNDROP AGENT db.s.agent1", ast.UndropAgent, "DB.S.AGENT1"},
		{"UNDROP EXTERNAL VOLUME vol1", ast.UndropExternalVolume, "VOL1"},
		{"UNDROP NOTEBOOK db.s.nb", ast.UndropNotebook, "DB.S.NB"},
		{"UNDROP SNAPSHOT example_snapshot", ast.UndropSnapshot, "EXAMPLE_SNAPSHOT"},
		{"UNDROP STAGE s.stg", ast.UndropStage, "S.STG"},
		{"UNDROP STREAMLIT app", ast.UndropStreamlit, "APP"},
		{"UNDROP TYPE db.s.t", ast.UndropType, "DB.S.T"},
		{"UNDROP VARIABLE v", ast.UndropVariable, "V"},
	} {
		stmt := parseOneUndrop(t, tt.sql)
		if stmt.Kind != tt.kind || stmt.Name.Normalize() != tt.name || stmt.NewName != nil {
			t.Errorf("%q: kind=%v name=%q new=%v", tt.sql, stmt.Kind, stmt.Name.Normalize(), stmt.NewName)
		}
	}

	stmt := parseOneUndrop(t, "UNDROP SNAPSHOT example_snapshot RENAME TO restored")
	if stmt.NewName == nil || stmt.NewName.Normalize() != "RESTORED" {
		t.Errorf("RENAME TO: %+v", stmt)
	}
	parseOneUndrop(t, "UNDROP SNAPSHOT IDENTIFIER(1)")

	for _, sql := range []string{
		// ACCOUNT and EXTERNAL VOLUME take one part.
		"UNDROP ACCOUNT org.acct",
		"UNDROP EXTERNAL VOLUME db.vol",
		// RENAME TO is UNDROP SNAPSHOT's, and takes one part.
		"UNDROP STAGE s RENAME TO t",
		"UNDROP SNAPSHOT s RENAME TO db.t",
		"UNDROP SNAPSHOT s RENAME t",
		"UNDROP EXTERNAL v",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}
