package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// The trees of the keyword-choice forms compat.tsv's (every alternative) rows
// found omni refusing or taking wrongly; the rows hold the verdicts, these
// hold what the parser builds.

func parseOne(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

// A hex or bit literal is a text_string, decoded to its bytes as LOAD DATA's
// terminators are.
func TestTextStringHexAndBitLiterals(t *testing.T) {
	for _, sql := range []string{"PREPARE p FROM 0x41", "PREPARE p FROM X'41'", "PREPARE p FROM b'1000001'"} {
		stmt, ok := parseOne(t, sql).(*nodes.PrepareStmt)
		if !ok || stmt.Stmt != "A" {
			t.Errorf("%s: got %#v, want a PrepareStmt of %q", sql, stmt, "A")
		}
	}
	sel, ok := parseOne(t, "SELECT a INTO OUTFILE 's' COLUMNS TERMINATED BY 0x2C ENCLOSED BY b'100010' LINES STARTING BY X'3E' TERMINATED BY 'x'").(*nodes.SelectStmt)
	if !ok || sel.Into == nil {
		t.Fatalf("got %#v, want a SELECT ... INTO OUTFILE", sel)
	}
	into := sel.Into
	if into.FieldsTerminatedBy != "," || into.FieldsEnclosedBy != `"` || into.LinesStartingBy != ">" || into.LinesTerminatedBy != "x" {
		t.Errorf("got terminated %q, enclosed %q, starting %q, lines terminated %q", into.FieldsTerminatedBy, into.FieldsEnclosedBy, into.LinesStartingBy, into.LinesTerminatedBy)
	}
}

// PERCENTILE_CONT takes ALL, DISTINCT or UNIQUE before two or more arguments
// (sql_parser_mysql_mode.y:2791); with one argument only the unquantified
// ordinary call remains.
func TestPercentileContQuantifier(t *testing.T) {
	for sql, want := range map[string]nodes.AggQuantifier{
		"SELECT PERCENTILE_CONT(ALL a, b)":         nodes.AggQuantAll,
		"SELECT PERCENTILE_CONT(DISTINCT a, b, c)": nodes.AggQuantDistinct,
		"SELECT PERCENTILE_CONT(UNIQUE a, b)":      nodes.AggQuantUnique,
	} {
		sel := parseOne(t, sql).(*nodes.SelectStmt)
		item := sel.TargetList[0]
		if rt, ok := item.(*nodes.ResTarget); ok {
			item = rt.Val
		}
		fc, ok := item.(*nodes.FuncCallExpr)
		if !ok || fc.Quantifier != want || len(fc.Args) < 2 {
			t.Errorf("%s: got %#v", sql, fc)
		}
	}
	for _, sql := range []string{"SELECT PERCENTILE_CONT(DISTINCT a)", "SELECT PERCENTILE_CONT(DISTINCT ALL a, b)"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted; OceanBase's grammar has no such form", sql)
		}
	}
}

func TestDropTablesKeyword(t *testing.T) {
	for _, sql := range []string{"DROP TABLES a, b", "DROP TEMPORARY TABLES IF EXISTS a"} {
		if _, ok := parseOne(t, sql).(*nodes.DropTableStmt); !ok {
			t.Errorf("%s: not a DropTableStmt", sql)
		}
	}
}

// opt_check_option and check_state write every word.
func TestCheckOptionAndStateAreWhole(t *testing.T) {
	for _, sql := range []string{
		"ALTER VIEW a AS SELECT 1 WITH CHECK",
		"ALTER VIEW a AS SELECT 1 WITH LOCAL CHECK",
		"ALTER VIEW a AS SELECT 1 WITH",
		"ALTER TABLE a ADD CHECK (a > 0) NOT",
		"CREATE TABLE a (b INT, CHECK (b > 0) NOT)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
	for _, sql := range []string{
		"ALTER VIEW a AS SELECT 1 WITH LOCAL CHECK OPTION",
		"ALTER TABLE a ADD CHECK (a > 0) NOT ENFORCED",
	} {
		parseOne(t, sql)
	}
}
