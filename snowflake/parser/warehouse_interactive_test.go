package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// CREATE INTERACTIVE WAREHOUSE and the nameless / adaptive ALTER WAREHOUSE
// actions (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-warehouse,
// https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse).

func TestParseCreateInteractiveWarehouse(t *testing.T) {
	in := "CREATE OR REPLACE INTERACTIVE WAREHOUSE sales_wh TABLES (orders, db.s.customers) " +
		"WAREHOUSE_SIZE = 'MEDIUM' AUTO_SUSPEND = NULL WITH TAG (cost = 'eng') MAX_CONCURRENCY_LEVEL = 4"
	stmt := mustCreateWarehouse(t, in)
	if !stmt.Interactive || !stmt.OrReplace || stmt.IfNotExists || stmt.Name.Name.Name != "sales_wh" {
		t.Fatalf("got %+v", stmt)
	}
	if got := objNames(stmt.Tables); len(got) != 2 || got[1] != "db.s.customers" {
		t.Errorf("Tables = %v", got)
	}
	if v, _ := optLitString(stmt.Options, "WAREHOUSE_SIZE"); v != "MEDIUM" {
		t.Errorf("WAREHOUSE_SIZE = %q", v)
	}
	if v, _ := optWords(stmt.Options, "AUTO_SUSPEND"); v != "NULL" {
		t.Errorf("AUTO_SUSPEND = %q", v)
	}
	if v, _ := optLitInt(stmt.Options, "MAX_CONCURRENCY_LEVEL"); v != 4 {
		t.Errorf("MAX_CONCURRENCY_LEVEL = %d", v)
	}
	if len(stmt.Tags) != 1 || stmt.Loc.Start != 0 || stmt.Loc.End != len(in) {
		t.Errorf("Tags = %d, Loc = %+v", len(stmt.Tags), stmt.Loc)
	}

	for _, ok := range []string{
		"CREATE INTERACTIVE WAREHOUSE a",
		"CREATE INTERACTIVE WAREHOUSE IF NOT EXISTS a WITH WAREHOUSE_SIZE = XSMALL MIN_CLUSTER_COUNT = 1",
		"CREATE INTERACTIVE WAREHOUSE a RESOURCE_MONITOR = m FALLBACK_WAREHOUSE = w INITIALLY_SUSPENDED = TRUE",
		"CREATE INTERACTIVE WAREHOUSE a TAG (t = 's') STATEMENT_TIMEOUT_IN_SECONDS = 10",
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = X2LARGE",
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = X3LARGE",
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = '2X-LARGE'",
		"CREATE INTERACTIVE WAREHOUSE a COMMENT = 'c' STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1 MAX_CONCURRENCY_LEVEL = 2",
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = XSMALL WAREHOUSE_SIZE = XSMALL",
		"CREATE INTERACTIVE WAREHOUSE a COMMENT = 'x' COMMENT = 'y'",
		"CREATE INTERACTIVE WAREHOUSE a MAX_CONCURRENCY_LEVEL = 1 MAX_CONCURRENCY_LEVEL = 2",
	} {
		if s := mustCreateWarehouse(t, ok); !s.Interactive {
			t.Errorf("%q: not interactive", ok)
		}
	}
	for _, bad := range []string{
		"CREATE OR ALTER INTERACTIVE WAREHOUSE a",                                // no OR ALTER
		"CREATE OR REPLACE INTERACTIVE WAREHOUSE IF NOT EXISTS a",                // IF NOT EXISTS only without OR REPLACE
		"CREATE INTERACTIVE WAREHOUSE db.a",                                      // a warehouse name has one part
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_TYPE = STANDARD",               // not an interactive property
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = HUGE",                   // not a size
		"CREATE INTERACTIVE WAREHOUSE a AUTO_SUSPEND = 'x'",                      // a number or NULL
		"CREATE INTERACTIVE WAREHOUSE a RESOURCE_MONITOR = db.m",                 // one-part name
		"CREATE INTERACTIVE WAREHOUSE a MAX_CONCURRENCY_LEVEL = 1 COMMENT = 'x'", // properties before parameters
		"CREATE INTERACTIVE WAREHOUSE a MAX_CONCURRENCY_LEVEL = 1 TAG (t = 's')", // tags before parameters
		"CREATE INTERACTIVE WAREHOUSE a TAG (t = 's') COMMENT = 'x'",             // properties before tags
		"CREATE INTERACTIVE WAREHOUSE a TAG (t = 's') WITH TAG (u = 's')",        // tags twice
		"CREATE INTERACTIVE WAREHOUSE a WITH MAX_CONCURRENCY_LEVEL = 1",          // WITH introduces properties
		"CREATE INTERACTIVE WAREHOUSE a COMMENT = 'x' WITH AUTO_RESUME = TRUE",   // WITH leads the run
		"CREATE INTERACTIVE WAREHOUSE a TABLES ()",                               // empty list
		"CREATE INTERACTIVE WAREHOUSE a TABLES (t,)",                             // trailing comma
		"CREATE INTERACTIVE WAREHOUSE a \"WAREHOUSE_SIZE\" = XSMALL",             // quoted key
		"CREATE INTERACTIVE WAREHOUSE a WAREHOUSE_SIZE = XSMALL TABLES (t)",      // TABLES comes first
		"CREATE TRANSIENT INTERACTIVE WAREHOUSE a",                               // no modifiers
		"CREATE \"INTERACTIVE\" WAREHOUSE a",                                     // quoted keyword
	} {
		mustRejectObject(t, bad)
	}
}

func TestParseAlterWarehouse_NamelessAndAdaptive(t *testing.T) {
	cases := []struct {
		in     string
		action ast.AlterWarehouseAction
		name   string
	}{
		{"ALTER WAREHOUSE SUSPEND", ast.AlterWarehouseSuspend, ""},
		{"ALTER WAREHOUSE IF EXISTS RESUME", ast.AlterWarehouseResume, ""},
		{"ALTER WAREHOUSE RESUME IF SUSPENDED", ast.AlterWarehouseResume, ""},
		{"ALTER WAREHOUSE ABORT ALL QUERIES", ast.AlterWarehouseAbort, ""},
		{"ALTER WAREHOUSE suspend SUSPEND", ast.AlterWarehouseSuspend, "suspend"},
		{"ALTER WAREHOUSE resume RESUME IF SUSPENDED", ast.AlterWarehouseResume, "resume"},
		{"ALTER WAREHOUSE abort ABORT ALL QUERIES", ast.AlterWarehouseAbort, "abort"},
		{"ALTER WAREHOUSE resume SET COMMENT = 'x'", ast.AlterWarehouseSet, "resume"},
		{"ALTER WAREHOUSE a ENABLE", ast.AlterWarehouseEnable, "a"},
		{"ALTER WAREHOUSE IF EXISTS a DISABLE", ast.AlterWarehouseDisable, "a"},
	}
	for _, c := range cases {
		stmt := mustAlterWarehouse(t, c.in)
		if stmt.Action != c.action {
			t.Errorf("%q: Action = %v, want %v", c.in, stmt.Action, c.action)
		}
		got := ""
		if stmt.Name != nil {
			got = stmt.Name.Name.Name
		}
		if got != c.name {
			t.Errorf("%q: Name = %q, want %q", c.in, got, c.name)
		}
	}
	if stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE RESUME IF SUSPENDED"); !stmt.ResumeIfSuspended {
		t.Error("RESUME IF SUSPENDED not recorded")
	}
	for _, bad := range []string{
		"ALTER WAREHOUSE ENABLE",      // ENABLE needs a name
		"ALTER WAREHOUSE DISABLE",     // DISABLE needs a name
		"ALTER WAREHOUSE RESUME IF",   // IF SUSPENDED is whole
		"ALTER WAREHOUSE ABORT ALL",   // ABORT ALL QUERIES is whole
		"ALTER WAREHOUSE a ENABLE b",  // nothing follows
		"ALTER WAREHOUSE IF EXISTS",   // no action
		"ALTER WAREHOUSE \"SUSPEND\"", // a quoted name is a name
	} {
		mustRejectObject(t, bad)
	}
}

func TestStartsNamelessWarehouseAction(t *testing.T) {
	for in, want := range map[string]bool{
		"SUSPEND":             true,
		"SUSPEND;":            true,
		"SUSPEND SUSPEND":     false,
		"RESUME":              true,
		"RESUME IF SUSPENDED": true,
		"RESUME IF EXISTS":    false,
		"RESUME RESUME":       false,
		"ABORT ALL QUERIES":   true,
		"ABORT ABORT":         false,
		"w SUSPEND":           false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsNamelessWarehouseAction(); got != want || p.cur.Loc.Start != 0 {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}
