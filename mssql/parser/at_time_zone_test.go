package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// atTimeZoneCases pin where SQL Server 2022 takes AT after an expression,
// as SET PARSEONLY ON reports it. AT starts AT TIME ZONE only when TIME ZONE
// follows it; otherwise it is a column alias, so SELECT x AT FROM t names
// the column AT. omni used to consume AT, and TIME after it, whatever came
// next, and drop them: it accepted SELECT x AT y FROM t with the alias y.
var atTimeZoneCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT d AT TIME ZONE 'UTC' FROM t", true},
	{"SELECT d AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time' FROM t", true},
	{"SELECT x AT FROM t", true},
	{"SELECT x AT", true},
	{"SELECT d AT TIME ZONE 'UTC' AT FROM t", true},

	{"SELECT x AT y FROM t", false},
	{"SELECT x AT TIME y FROM t", false},
	{"SELECT x AT TIME FROM t", false},
	{"SELECT x AT TIME", false},
	{"SELECT x AT ZONE 'UTC' FROM t", false},
	{"SELECT x AT TIME ZONE FROM t", false},
	{"SELECT x COLLATE Latin1_General_CI_AS AT y FROM t", false},
	{"SELECT d AT TIME ZONE 'UTC' AT y FROM t", false},
	{"SELECT 1 FROM t WHERE x AT y = 1", false},
	{"SELECT 1 FROM t WHERE x AT TIME y = 1", false},
	{"UPDATE t SET x = y AT TIME 'UTC'", false},
	{"DECLARE @v DATETIMEOFFSET = SYSDATETIMEOFFSET() AT 'UTC'", false},
}

// TestAtTimeZoneOracle: omni accepts AT after an expression exactly where
// SQL Server does, and reports a rejected one where SqlScriptDOM does.
func TestAtTimeZoneOracle(t *testing.T) {
	oracle := startParserOracle(t)
	h := getHarness(t)
	for _, tc := range atTimeZoneCases {
		t.Run(tc.sql, func(t *testing.T) {
			serverErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			if (serverErr == nil) != tc.accept {
				t.Fatalf("SQL Server: %v, want accept = %v", serverErr, tc.accept)
			}
			_, omniErr := Parse(tc.sql)
			if (omniErr == nil) != tc.accept {
				t.Fatalf("omni: %v, want accept = %v", omniErr, tc.accept)
			}
			if tc.accept {
				return
			}
			errs, _ := h.Shape(t, tc.sql)["errors"].([]any)
			if len(errs) == 0 {
				t.Fatal("SqlScriptDOM accepted it")
			}
			first, _ := errs[0].(map[string]any)
			line, _ := first["line"].(float64)
			column, _ := first["column"].(float64)
			var pe *ParseError
			if !errors.As(omniErr, &pe) {
				t.Fatalf("omni error %v is not a *ParseError", omniErr)
			}
			if want := lineColumnOffset(tc.sql, int(line), int(column)); pe.Position != want {
				t.Errorf("omni reports offset %d, SqlScriptDOM line %d column %d (offset %d)",
					pe.Position, int(line), int(column), want)
			}
		})
	}
}

// TestAtAlias: AT that TIME ZONE does not follow is the alias of the select
// item before it.
func TestAtAlias(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		alias string
		val   string
	}{
		{"SELECT x AT FROM t", "AT", "x"},
		{"SELECT x at", "at", "x"},
		{"SELECT d AT TIME ZONE 'UTC' AT FROM t", "AT", "d AT TIME ZONE 'UTC'"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			rt := list.Items[0].(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget)
			if rt.Name != tc.alias {
				t.Errorf("alias %q, want %q", rt.Name, tc.alias)
			}
			loc := ast.NodeLoc(rt.Val)
			if got := tc.sql[loc.Start:loc.End]; got != tc.val {
				t.Errorf("value %q, want %q", got, tc.val)
			}
		})
	}
}
