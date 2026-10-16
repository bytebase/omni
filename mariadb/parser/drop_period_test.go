package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// dropPeriodOracleCases are the DROP PERIOD arms of alter_list_item in
// MariaDB 11.8 sql/sql_yacc.yy, with the code mariadb:11.8.8 gives each:
//
//	DROP PERIOD_SYM FOR_SYSTEM_TIME_SYM
//	DROP PERIOD_SYM opt_if_exists_table_element FOR_SYM ident
var dropPeriodOracleCases = []parseCodeCase{
	// IF EXISTS before an application-time period name.
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR `app`", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR period", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR SYSTEM_TIME_X", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR `if`", 0},
	{"ALTER TABLE t DROP PERIOD /* c */ IF /* c */ EXISTS /* c */ FOR app", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS /*!100000 FOR app */", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app, DROP PERIOD IF EXISTS FOR app", 0},
	// The lexer merges FOR and an unquoted SYSTEM_TIME, in any case and across
	// a comment, into FOR_SYSTEM_TIME_SYM, which the IF EXISTS arm does not
	// take. Quoted, SYSTEM_TIME is an ident and names an application-time
	// period.
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR SYSTEM_TIME", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR system_time", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR /* c */ SYSTEM_TIME", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS /*!FOR SYSTEM_TIME */", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", 0},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR \"SYSTEM_TIME\"", 1064},
	// Without IF EXISTS.
	{"ALTER TABLE t DROP PERIOD FOR app", 0},
	{"ALTER TABLE t DROP PERIOD FOR SYSTEM_TIME", 0},
	{"ALTER TABLE t DROP PERIOD FOR `SYSTEM_TIME`", 0},
	// Malformed IF EXISTS and name.
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS app", 1064},
	{"ALTER TABLE t DROP PERIOD IF FOR app", 1064},
	{"ALTER TABLE t DROP PERIOD EXISTS FOR app", 1064},
	{"ALTER TABLE t DROP PERIOD IF NOT EXISTS FOR app", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR if", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR FOR", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR t.app", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app RESTRICT", 1064},
	{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app CASCADE", 1064},
	{"ALTER TABLE t DROP PERIOD FOR app RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN PERIOD IF EXISTS FOR app", 1064},
	// An unquoted PERIOD right after DROP always starts a DROP PERIOD arm; a
	// quoted `PERIOD` is a column name.
	{"ALTER TABLE t DROP period", 1064},
	{"ALTER TABLE t DROP PERIOD", 1064},
	{"ALTER TABLE t DROP period, DROP id", 1064},
	{"ALTER TABLE t DROP period RESTRICT", 1064},
	{"ALTER TABLE t DROP `period`", 0},
	{"ALTER TABLE t DROP COLUMN period", 0},
	{"ALTER TABLE t DROP IF EXISTS period", 0},
	{"ALTER TABLE t DROP `PERIOD` FOR app", 1064},
	{"ALTER TABLE t DROP `PERIOD` IF EXISTS FOR app", 1064},
}

// TestMariaDBDropPeriodOracle asserts omni agrees with a live mariadb:11.8.8
// on the DROP PERIOD forms in dropPeriodOracleCases (container test).
func TestMariaDBDropPeriodOracle(t *testing.T) {
	assertParseCodes(t, dropPeriodOracleCases)
}

// TestDropPeriodAST checks the command each DROP PERIOD form parses into and
// the text its Loc spans.
func TestDropPeriodAST(t *testing.T) {
	cases := []struct {
		sql      string
		period   string
		ifExists bool
		text     string
	}{
		{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app", "app", true, "DROP PERIOD IF EXISTS FOR app"},
		{"ALTER TABLE t DROP PERIOD if exists FOR `app`", "app", true, "DROP PERIOD if exists FOR `app`"},
		{"ALTER TABLE t DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", "SYSTEM_TIME", true, "DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`"},
		{"ALTER TABLE t DROP PERIOD IF EXISTS FOR app, DROP COLUMN c", "app", true, "DROP PERIOD IF EXISTS FOR app"},
		{"ALTER TABLE t DROP PERIOD FOR app", "app", false, "DROP PERIOD FOR app"},
		{"ALTER TABLE t DROP PERIOD FOR SYSTEM_TIME", "SYSTEM_TIME", false, "DROP PERIOD FOR SYSTEM_TIME"},
		// A later segment: the quoted check reads the segment's own text.
		{"SELECT 1;\nALTER TABLE t DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`", "SYSTEM_TIME", true, "DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			list := ParseAndCheck(t, tc.sql)
			stmt, ok := list.Items[len(list.Items)-1].(*ast.AlterTableStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.AlterTableStmt", list.Items[len(list.Items)-1])
			}
			cmd := stmt.Commands[0]
			if cmd.Type != ast.ATDropPeriod {
				t.Fatalf("Type = %v, want ATDropPeriod", cmd.Type)
			}
			if cmd.PeriodName != tc.period || cmd.IfExists != tc.ifExists {
				t.Errorf("PeriodName, IfExists = %q, %v; want %q, %v", cmd.PeriodName, cmd.IfExists, tc.period, tc.ifExists)
			}
			if got := tc.sql[cmd.Loc.Start:cmd.Loc.End]; got != tc.text {
				t.Errorf("Loc text = %q, want %q", got, tc.text)
			}
			if tc.ifExists && !strings.Contains(ast.NodeToString(cmd), ":if_exists true") {
				t.Errorf("NodeToString missing :if_exists:\n%s", ast.NodeToString(cmd))
			}
		})
	}

	t.Run("quoted period is a column", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE t DROP `period`")
		if cmd := stmt.Commands[0]; cmd.Type != ast.ATDropColumn || cmd.Name != "period" {
			t.Errorf("Type, Name = %v, %q; want ATDropColumn, period", cmd.Type, cmd.Name)
		}
	})
}

// TestDropPeriodErrorPosition checks that a rejected DROP PERIOD reports the
// byte offset of the token MariaDB names, in the whole script.
func TestDropPeriodErrorPosition(t *testing.T) {
	cases := []struct {
		sql  string
		near string
	}{
		{"ALTER TABLE t DROP PERIOD IF EXISTS FOR SYSTEM_TIME", "SYSTEM_TIME"},
		{"SELECT 1;\nALTER TABLE t DROP PERIOD IF EXISTS FOR /* c */ SYSTEM_TIME", "SYSTEM_TIME"},
		{"ALTER TABLE t DROP PERIOD IF EXISTS app", "app"},
		{"ALTER TABLE t DROP `PERIOD` FOR app", "FOR app"},
		{"ALTER TABLE t DROP period, DROP id", ", DROP id"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want *ParseError", err)
			}
			if want := strings.LastIndex(tc.sql, tc.near); pe.Position != want {
				t.Errorf("Position = %d, want %d (%q)", pe.Position, want, tc.near)
			}
		})
	}
}
