package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// mergeStmt is a complete MERGE without its terminating semicolon.
const mergeStmt = "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE"

// mergeTerminatorCases covers the semicolon SQL Server requires after a
// MERGE, wherever the MERGE stands, and the semicolons it takes after the
// THEN statement of an IF before ELSE, which a MERGE there needs.
var mergeTerminatorCases = []struct {
	name    string
	sql     string
	accepts bool
}{
	// The semicolon after MERGE.
	{"semicolon", mergeStmt + ";", true},
	{"semicolon_then_statement", mergeStmt + "; SELECT 1", true},
	{"two_semicolons", mergeStmt + ";;", true},
	{"comment_before_semicolon", mergeStmt + " /* c */ ;", true},
	{"newline_before_semicolon", mergeStmt + "\n;", true},
	{"line_comment_before_semicolon", mergeStmt + " -- c\n;", true},
	{"output_semicolon", mergeStmt + " OUTPUT $action;", true},
	{"option_semicolon", mergeStmt + " OPTION (RECOMPILE);", true},
	{"cte_semicolon", "WITH c AS (SELECT 1 AS a) " + mergeStmt + ";", true},
	{"begin_end_semicolon", "BEGIN " + mergeStmt + "; END", true},
	{"if_semicolon", "IF 1 = 1 " + mergeStmt + ";", true},
	{"if_else_semicolon", "IF 1 = 1 " + mergeStmt + "; ELSE SELECT 1", true},
	{"try_semicolon", "BEGIN TRY " + mergeStmt + "; END TRY BEGIN CATCH SELECT 1 END CATCH", true},
	{"while_semicolon", "WHILE 1 = 0 " + mergeStmt + ";", true},
	{"after_select", "SELECT 1 " + mergeStmt + ";", true},
	{"two_merges", mergeStmt + "; " + mergeStmt + ";", true},
	{"update_semicolon", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN UPDATE SET a = 1;", true},
	{"procedure_semicolon", "CREATE PROCEDURE p AS " + mergeStmt + ";", true},
	{"end_of_batch", mergeStmt, false},
	{"then_statement", mergeStmt + " SELECT 1", false},
	{"output", mergeStmt + " OUTPUT $action", false},
	{"option", mergeStmt + " OPTION (RECOMPILE)", false},
	{"cte", "WITH c AS (SELECT 1 AS a) " + mergeStmt, false},
	{"begin_end", "BEGIN " + mergeStmt + " END", false},
	{"if", "IF 1 = 1 " + mergeStmt, false},
	{"if_else", "IF 1 = 1 " + mergeStmt + " ELSE SELECT 1", false},
	{"try", "BEGIN TRY " + mergeStmt + " END TRY BEGIN CATCH SELECT 1 END CATCH", false},
	{"while", "WHILE 1 = 0 " + mergeStmt, false},
	{"two_merges_first_unterminated", mergeStmt + " " + mergeStmt + ";", false},
	{"update", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN UPDATE SET a = 1", false},
	{"insert", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN NOT MATCHED THEN INSERT (a) VALUES (1)", false},
	{"insert_default_values", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN NOT MATCHED THEN INSERT DEFAULT VALUES", false},
	{"procedure", "CREATE PROCEDURE p AS " + mergeStmt, false},
	{"procedure_begin_end", "CREATE PROCEDURE p AS BEGIN " + mergeStmt + " END", false},

	// Semicolons before ELSE.
	{"if_select_semicolon_else", "IF 1 = 1 SELECT 1; ELSE SELECT 2", true},
	{"if_select_else", "IF 1 = 1 SELECT 1 ELSE SELECT 2", true},
	{"if_two_semicolons_else", "IF 1 = 1 SELECT 1;; ELSE SELECT 2", true},
	{"if_else_both_semicolons", "IF 1 = 1 SELECT 1; ELSE SELECT 2;", true},
	{"if_block_semicolon_else", "IF 1 = 1 BEGIN SELECT 1 END; ELSE SELECT 2", true},
	{"nested_if_semicolons", "IF 1 = 1 IF 2 = 2 SELECT 1; ELSE SELECT 2; ELSE SELECT 3", true},
	{"if_semicolon_then_statement", "IF 1 = 1 SELECT 1; SELECT 2", true},
	{"if_two_statements_else", "IF 1 = 1 SELECT 1; SELECT 3; ELSE SELECT 2", false},
}

// TestMergeTerminator checks that omni accepts the MERGE and IF forms SQL
// Server accepts, rejects the others, reads an ELSE after semicolons as the
// IF's, and keeps every Loc inside its parent.
func TestMergeTerminator(t *testing.T) {
	for _, tc := range mergeTerminatorCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error:\n%s", tc.sql, ast.NodeToString(list))
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			if strings.Contains(tc.sql, "ELSE") {
				ifStmt, ok := nthNodeByType(list, "IfStmt", 0).(*ast.IfStmt)
				if !ok || ifStmt.Else == nil {
					t.Errorf("Parse(%q): the first IF has no ELSE:\n%s", tc.sql, ast.NodeToString(list))
				}
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestMergeTerminatorOracle checks accept/reject parity with SQL Server for
// every case in mergeTerminatorCases.
func TestMergeTerminatorOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range mergeTerminatorCases {
		t.Run(tc.name, func(t *testing.T) {
			ssErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, tc.sql)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

// TestMergeTerminatorError checks that a missing semicolon after MERGE is
// reported at the token that stands in its place, with SQL Server's message.
func TestMergeTerminatorError(t *testing.T) {
	tests := []struct {
		sql  string
		want string // Error()
		pos  int
	}{
		{mergeStmt, "syntax error at end of input", len(mergeStmt)},
		{mergeStmt + " SELECT 1", `syntax error at or near "SELECT"`, len(mergeStmt) + 1},
		{"IF 1 = 1 " + mergeStmt + " ELSE SELECT 1", `syntax error at or near "ELSE"`, len("IF 1 = 1 " + mergeStmt + " ")},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			if pe.Message != "A MERGE statement must be terminated by a semi-colon (;)" {
				t.Errorf("Message = %q", pe.Message)
			}
			if got := pe.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
			if pe.Position != tc.pos {
				t.Errorf("Position = %d, want %d", pe.Position, tc.pos)
			}
		})
	}
}

// TestMergeTerminatorCompletion checks that completion does not need the
// semicolon after a MERGE the user has not typed yet, inside the MERGE or in
// a statement after it.
func TestMergeTerminatorCompletion(t *testing.T) {
	tests := []struct {
		sql    string
		cursor string // the text before the cursor
		rule   string
	}{
		{mergeStmt, "MERGE INTO dbo.t USING t1 ON ", "columnref"},
		{mergeStmt + "\nSELECT * FROM t WHERE ", mergeStmt + "\nSELECT * FROM t WHERE ", "columnref"},
		{mergeStmt + "\nSELECT * FROM ", mergeStmt + "\nSELECT * FROM ", "table_ref"},
	}
	for _, tc := range tests {
		t.Run(tc.cursor+"|", func(t *testing.T) {
			cs := Collect(tc.sql, len(tc.cursor))
			if cs == nil || !cs.HasRule(tc.rule) {
				t.Errorf("missing rule candidate %q", tc.rule)
			}
		})
	}
}
