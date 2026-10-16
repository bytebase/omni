package parser

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestParseSetOptionList checks the AST of SET statements that list several
// ON/OFF options sharing one value, and that single-option statements report
// their one name in Options too.
func TestParseSetOptionList(t *testing.T) {
	tests := []struct {
		sql     string
		option  string
		options []string
		value   string // ON/OFF text of Value, "" to skip
		end     int    // Loc.End of the statement, -1 for len(sql)
	}{
		{
			sql:     "SET ANSI_NULLS, QUOTED_IDENTIFIER OFF",
			option:  "ANSI_NULLS",
			options: []string{"ANSI_NULLS", "QUOTED_IDENTIFIER"},
			value:   "OFF",
			end:     -1,
		},
		{
			sql:     "SET ANSI_PADDING, ANSI_WARNINGS, CONCAT_NULL_YIELDS_NULL, ARITHABORT, QUOTED_IDENTIFIER, ANSI_NULLS ON",
			option:  "ANSI_PADDING",
			options: []string{"ANSI_PADDING", "ANSI_WARNINGS", "CONCAT_NULL_YIELDS_NULL", "ARITHABORT", "QUOTED_IDENTIFIER", "ANSI_NULLS"},
			value:   "ON",
			end:     -1,
		},
		{
			sql:     "set nocount /* c */ , xact_abort on",
			option:  "NOCOUNT",
			options: []string{"NOCOUNT", "XACT_ABORT"},
			value:   "ON",
			end:     -1,
		},
		{
			sql:     "SET STATISTICS IO, TIME ON",
			option:  "STATISTICS IO",
			options: []string{"STATISTICS IO", "STATISTICS TIME"},
			value:   "ON",
			end:     -1,
		},
		{
			sql:     "SET STATISTICS PROFILE, XML OFF",
			option:  "STATISTICS PROFILE",
			options: []string{"STATISTICS PROFILE", "STATISTICS XML"},
			value:   "OFF",
			end:     -1,
		},
		{
			sql:     "SET ANSI_NULLS, QUOTED_IDENTIFIER OFF;",
			option:  "ANSI_NULLS",
			options: []string{"ANSI_NULLS", "QUOTED_IDENTIFIER"},
			value:   "OFF",
			end:     len("SET ANSI_NULLS, QUOTED_IDENTIFIER OFF"),
		},
		{
			sql:     "SET ANSI_NULLS ON",
			option:  "ANSI_NULLS",
			options: []string{"ANSI_NULLS"},
			value:   "ON",
			end:     -1,
		},
		{
			sql:     "SET STATISTICS IO ON",
			option:  "STATISTICS IO",
			options: []string{"STATISTICS IO"},
			value:   "ON",
			end:     -1,
		},
		{
			sql:     "SET DATEFORMAT mdy",
			option:  "DATEFORMAT",
			options: []string{"DATEFORMAT"},
			end:     -1,
		},
		{
			sql:     "SET TRANSACTION ISOLATION LEVEL SNAPSHOT",
			option:  "TRANSACTION ISOLATION LEVEL",
			options: []string{"TRANSACTION ISOLATION LEVEL"},
			end:     -1,
		},
		{
			sql:     "SET IDENTITY_INSERT dbo.t ON",
			option:  "IDENTITY_INSERT",
			options: []string{"IDENTITY_INSERT"},
			end:     -1,
		},
		{
			sql:     "SET OFFSETS SELECT, FROM ON",
			option:  "OFFSETS SELECT, FROM",
			options: []string{"OFFSETS SELECT, FROM"},
			end:     -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if len(result.Items) != 1 {
				t.Fatalf("got %d statements, want 1", len(result.Items))
			}
			stmt, ok := result.Items[0].(*ast.SetOptionStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.SetOptionStmt", result.Items[0])
			}
			if stmt.Option != tt.option {
				t.Errorf("Option = %q, want %q", stmt.Option, tt.option)
			}
			if !slices.Equal(stmt.Options, tt.options) {
				t.Errorf("Options = %q, want %q", stmt.Options, tt.options)
			}
			end := tt.end
			if end < 0 {
				end = len(tt.sql)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != end {
				t.Errorf("Loc = %+v, want {0 %d}", stmt.Loc, end)
			}
			if tt.value == "" {
				return
			}
			ref, ok := stmt.Value.(*ast.ColumnRef)
			if !ok {
				t.Fatalf("Value = %T, want *ast.ColumnRef", stmt.Value)
			}
			if ref.Column != tt.value {
				t.Errorf("Value = %q, want %q", ref.Column, tt.value)
			}
			start := strings.LastIndex(strings.ToUpper(tt.sql), tt.value)
			if ref.Loc.Start != start || ref.Loc.End != start+len(tt.value) {
				t.Errorf("Value.Loc = %+v, want {%d %d}", ref.Loc, start, start+len(tt.value))
			}
		})
	}
}

// TestParseSetOptionListStatementEnd checks that a list ends at its ON/OFF,
// so a statement may follow without a semicolon.
func TestParseSetOptionListStatementEnd(t *testing.T) {
	sql := "SET ANSI_NULLS, QUOTED_IDENTIFIER OFF SELECT 1"
	result := ParseAndCheck(t, sql)
	if len(result.Items) != 2 {
		t.Fatalf("got %d statements, want 2", len(result.Items))
	}
	stmt, ok := result.Items[0].(*ast.SetOptionStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.SetOptionStmt", result.Items[0])
	}
	if want := len("SET ANSI_NULLS, QUOTED_IDENTIFIER OFF"); stmt.Loc.End != want {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, want)
	}
	if _, ok := result.Items[1].(*ast.SelectStmt); !ok {
		t.Errorf("second statement = %T, want *ast.SelectStmt", result.Items[1])
	}
}

// TestParseSetOptionListErrors checks where a malformed list is reported.
func TestParseSetOptionListErrors(t *testing.T) {
	tests := []struct {
		sql string
		at  string // the error is at the first occurrence of this text, "" for end of input
	}{
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER", ""},
		{"SET ANSI_NULLS, ON", "ON"},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER 1", "1"},
		{"SET ANSI_NULLS,, QUOTED_IDENTIFIER ON", ", QUOTED"},
		{"SET ANSI_NULLS, FOO ON", "FOO"},
		{"SET ANSI_NULLS, [QUOTED_IDENTIFIER] ON", "[QUOTED"},
		{"SET [ANSI_NULLS], QUOTED_IDENTIFIER ON", "[ANSI"},
		{"SET ANSI_NULLS, SHOWPLAN_XML OFF", "SHOWPLAN"},
		{"SET SHOWPLAN_TEXT, NOCOUNT OFF", "SHOWPLAN"},
		{"SET ANSI_NULLS, DATEFORMAT ON", "DATEFORMAT"},
		{"SET DATEFORMAT, DATEFIRST mdy", "DATEFORMAT"},
		{"SET ANSI_NULLS, ROWCOUNT ON", "ROWCOUNT"},
		{"SET NOCOUNT, STATISTICS IO ON", "STATISTICS"},
		{"SET STATISTICS IO, NOCOUNT ON", "NOCOUNT"},
		{"SET STATISTICS IO, STATISTICS TIME ON", "STATISTICS TIME"},
		{"SET STATISTICS FOO, IO ON", "FOO"},
		{"SET STATISTICS, IO ON", "STATISTICS"},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want a *ParseError", err)
			}
			want := len(tt.sql)
			if tt.at != "" {
				want = strings.Index(tt.sql, tt.at)
			}
			if pe.Position != want {
				t.Errorf("Position = %d, want %d (%v)", pe.Position, want, err)
			}
		})
	}
}

// TestSetOptionListOracle checks omni against SQL Server 2022 on SET
// statements that list options. Only ON/OFF options take the shared-value
// list; SHOWPLAN_* may not be listed, and STATISTICS lists only its own
// options.
func TestSetOptionListOracle(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		sql    string
		accept bool
	}{
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER OFF", true},
		{"SET ANSI_PADDING, ANSI_WARNINGS, CONCAT_NULL_YIELDS_NULL, ARITHABORT, QUOTED_IDENTIFIER, ANSI_NULLS ON", true},
		{"SET NUMERIC_ROUNDABORT, ARITHIGNORE, CURSOR_CLOSE_ON_COMMIT, IMPLICIT_TRANSACTIONS, ANSI_NULL_DFLT_ON OFF", true},
		{"SET ANSI_NULL_DFLT_OFF, FORCEPLAN, REMOTE_PROC_TRANSACTIONS, NO_BROWSETABLE OFF", true},
		{"SET ANSI_DEFAULTS, FMTONLY, NOEXEC OFF", true},
		{"SET NOCOUNT, XACT_ABORT ON", true},
		{"SET PARSEONLY, NOCOUNT ON", true},
		{"SET ANSI_NULLS, ANSI_NULLS ON", true},
		{"set ansi_nulls /* c */ , quoted_identifier on", true},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER OFF SELECT 1", true},
		{"IF 1 = 1 SET ANSI_NULLS, ANSI_WARNINGS ON ELSE SET ANSI_NULLS, ANSI_WARNINGS OFF", true},
		{"CREATE PROCEDURE dbo.p_set_list AS SET NOCOUNT, XACT_ABORT ON; SELECT 1", true},
		{"SET STATISTICS IO, TIME ON", true},
		{"SET STATISTICS PROFILE, XML, IO, TIME OFF", true},
		{"SET STATISTICS IO, IO ON", true},

		{"SET ANSI_NULLS, ON", false},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER", false},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER 1", false},
		{"SET ANSI_NULLS ON, QUOTED_IDENTIFIER ON", false},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER ON, ARITHABORT ON", false},
		{"SET ANSI_NULLS,, QUOTED_IDENTIFIER ON", false},
		{"SET , ANSI_NULLS ON", false},
		{"SET ANSI_NULLS QUOTED_IDENTIFIER ON", false},
		{"SET ANSI_NULLS, FOO ON", false},
		{"SET ANSI_NULLS, [QUOTED_IDENTIFIER] ON", false},
		{"SET [ANSI_NULLS], QUOTED_IDENTIFIER ON", false},
		{"SET ANSI_NULLS, SHOWPLAN_XML OFF", false},
		{"SET SHOWPLAN_TEXT, NOCOUNT OFF", false},
		{"SET DISABLE_DEF_CNST_CHK, NOCOUNT OFF", false},
		{"SET ANSI_NULLS, DATEFORMAT mdy", false},
		{"SET DATEFORMAT, DATEFIRST mdy", false},
		{"SET LANGUAGE us_english, NOCOUNT ON", false},
		{"SET DATEFIRST 1, NOCOUNT ON", false},
		{"SET ROWCOUNT 10, NOCOUNT ON", false},
		{"SET ANSI_NULLS, ROWCOUNT ON", false},
		{"SET ANSI_NULLS, TEXTSIZE ON", false},
		{"SET ANSI_NULLS, QUOTED_IDENTIFIER OFF, DATEFORMAT mdy", false},
		{"SET NOCOUNT, TRANSACTION ISOLATION LEVEL READ COMMITTED", false},
		{"SET TRANSACTION ISOLATION LEVEL READ COMMITTED, NOCOUNT ON", false},
		{"SET NOCOUNT, IDENTITY_INSERT dbo.t ON", false},
		{"SET IDENTITY_INSERT dbo.t, NOCOUNT ON", false},
		{"SET NOCOUNT, OFFSETS SELECT ON", false},
		{"SET NOCOUNT, STATISTICS IO ON", false},
		{"SET STATISTICS IO, STATISTICS TIME ON", false},
		{"SET STATISTICS IO, NOCOUNT ON", false},
		{"SET STATISTICS IO, FOO ON", false},
		{"SET STATISTICS FOO, IO ON", false},
		{"SET STATISTICS IO, TIME", false},
		{"SET STATISTICS, IO ON", false},
		{"SET STATISTICS IO, [TIME] ON", false},
	}

	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			// QUOTED_IDENTIFIER takes effect while the batch is parsed, even
			// under SET PARSEONLY ON, so restore it for the tests that follow.
			if _, err := oracle.conn.ExecContext(oracle.ctx, "SET QUOTED_IDENTIFIER ON"); err != nil {
				t.Fatalf("restore QUOTED_IDENTIFIER: %v", err)
			}
			if ssAccepts != tc.accept {
				t.Errorf("SQL Server %s, test expects %s\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(tc.accept), tc.sql)
			}
			_, omniErr := Parse(tc.sql)
			if omniAccepts := omniErr == nil; omniAccepts != ssAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}
