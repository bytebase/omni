package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateWarehouse(t *testing.T, input string) *ast.CreateWarehouseStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateWarehouseStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateWarehouseStmt", input, node)
	}
	return stmt
}

func mustAlterWarehouse(t *testing.T, input string) *ast.AlterWarehouseStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterWarehouseStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterWarehouseStmt", input, node)
	}
	return stmt
}

// optWords / optLitString / optLitInt are small accessors for asserting on the
// open-ended CopyOption value shapes that warehouse properties use.
func optWords(opts []*ast.CopyOption, name string) (string, bool) {
	o := findOption(opts, name)
	if o == nil {
		return "", false
	}
	return o.Words, true
}

func optLitString(opts []*ast.CopyOption, name string) (string, bool) {
	o := findOption(opts, name)
	if o == nil || o.Lit == nil || o.Lit.Kind != ast.LitString {
		return "", false
	}
	return o.Lit.Value, true
}

func optLitInt(opts []*ast.CopyOption, name string) (int64, bool) {
	o := findOption(opts, name)
	if o == nil || o.Lit == nil || o.Lit.Kind != ast.LitInt {
		return 0, false
	}
	return o.Lit.Ival, true
}

// ---------------------------------------------------------------------------
// CREATE WAREHOUSE — modifiers, name, IF NOT EXISTS
// ---------------------------------------------------------------------------

func TestParseCreateWarehouse_Modifiers(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE my_wh")
		if stmt.Name.String() != "my_wh" {
			t.Errorf("Name = %q, want my_wh", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.OrAlter || stmt.IfNotExists {
			t.Errorf("unexpected modifier set: %+v", stmt)
		}
		if len(stmt.Options) != 0 || len(stmt.Tags) != 0 {
			t.Errorf("expected no options/tags, got %+v / %+v", stmt.Options, stmt.Tags)
		}
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR REPLACE WAREHOUSE w")
		if !stmt.OrReplace || stmt.OrAlter {
			t.Errorf("OrReplace=%v OrAlter=%v, want true/false", stmt.OrReplace, stmt.OrAlter)
		}
	})

	t.Run("or alter", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR ALTER WAREHOUSE w")
		if stmt.OrReplace || !stmt.OrAlter {
			t.Errorf("OrReplace=%v OrAlter=%v, want false/true", stmt.OrReplace, stmt.OrAlter)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE IF NOT EXISTS w")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE db.sch.w")
		if stmt.Name.String() != "db.sch.w" {
			t.Errorf("Name = %q, want db.sch.w", stmt.Name.String())
		}
	})

	t.Run("quoted name", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, `CREATE WAREHOUSE "My WH"`)
		if stmt.Name.String() != `"My WH"` {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE WAREHOUSE — official docs corpus shapes
// ---------------------------------------------------------------------------

func TestParseCreateWarehouse_CorpusShapes(t *testing.T) {
	// example_01: CREATE OR REPLACE WAREHOUSE my_wh WITH WAREHOUSE_SIZE = 'X-LARGE'
	t.Run("with leading WITH + string size", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR REPLACE WAREHOUSE my_wh WITH WAREHOUSE_SIZE = 'X-LARGE'")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
		if v, ok := optLitString(stmt.Options, "WAREHOUSE_SIZE"); !ok || v != "X-LARGE" {
			t.Errorf("WAREHOUSE_SIZE = %q (ok=%v), want 'X-LARGE'", v, ok)
		}
	})

	// example_02: WAREHOUSE_SIZE = LARGE (word) INITIALLY_SUSPENDED = TRUE (word), no WITH
	t.Run("bare word values, no WITH", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR REPLACE WAREHOUSE my_wh WAREHOUSE_SIZE = LARGE INITIALLY_SUSPENDED = TRUE")
		if v, ok := optWords(stmt.Options, "WAREHOUSE_SIZE"); !ok || v != "LARGE" {
			t.Errorf("WAREHOUSE_SIZE words = %q (ok=%v), want LARGE", v, ok)
		}
		if v, ok := optWords(stmt.Options, "INITIALLY_SUSPENDED"); !ok || v != "TRUE" {
			t.Errorf("INITIALLY_SUSPENDED words = %q (ok=%v), want TRUE", v, ok)
		}
	})

	// example_03: WITH + WAREHOUSE_TYPE (string) + WAREHOUSE_SIZE (word) + RESOURCE_CONSTRAINT (string)
	t.Run("snowpark-optimized multi-prop", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE so_warehouse WITH WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAREHOUSE_SIZE = XLARGE RESOURCE_CONSTRAINT = 'MEMORY_16X_x86'")
		if len(stmt.Options) != 3 {
			t.Fatalf("got %d options, want 3: %+v", len(stmt.Options), stmt.Options)
		}
		if v, ok := optLitString(stmt.Options, "WAREHOUSE_TYPE"); !ok || v != "SNOWPARK-OPTIMIZED" {
			t.Errorf("WAREHOUSE_TYPE = %q (ok=%v)", v, ok)
		}
		if v, ok := optWords(stmt.Options, "WAREHOUSE_SIZE"); !ok || v != "XLARGE" {
			t.Errorf("WAREHOUSE_SIZE = %q (ok=%v)", v, ok)
		}
		if v, ok := optLitString(stmt.Options, "RESOURCE_CONSTRAINT"); !ok || v != "MEMORY_16X_x86" {
			t.Errorf("RESOURCE_CONSTRAINT = %q (ok=%v)", v, ok)
		}
	})

	// example_04: GENERATION = '2'
	t.Run("generation string", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE gen2_wh WITH WAREHOUSE_SIZE = LARGE GENERATION = '2'")
		if v, ok := optLitString(stmt.Options, "GENERATION"); !ok || v != "2" {
			t.Errorf("GENERATION = %q (ok=%v), want '2'", v, ok)
		}
	})

	// example_05: CREATE OR ALTER + AUTO_RESUME (word) + COMMENT (string)
	t.Run("or alter with auto_resume + comment", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR ALTER WAREHOUSE so_warehouse WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' AUTO_RESUME = TRUE COMMENT = 'Snowpark warehouse for ingestion'")
		if !stmt.OrAlter {
			t.Error("OrAlter not set")
		}
		if v, ok := optWords(stmt.Options, "AUTO_RESUME"); !ok || v != "TRUE" {
			t.Errorf("AUTO_RESUME = %q (ok=%v), want TRUE", v, ok)
		}
		if v, ok := optLitString(stmt.Options, "COMMENT"); !ok || v != "Snowpark warehouse for ingestion" {
			t.Errorf("COMMENT = %q (ok=%v)", v, ok)
		}
	})

	// example_07: GENERATION='1' AUTO_SUSPEND=60 (int) INITIALLY_SUSPENDED=TRUE, with leading WITH
	t.Run("auto_suspend integer value", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE OR ALTER WAREHOUSE test_gen_warehouse WITH WAREHOUSE_SIZE = XSMALL GENERATION = '1' AUTO_SUSPEND = 60 INITIALLY_SUSPENDED = TRUE")
		if v, ok := optLitInt(stmt.Options, "AUTO_SUSPEND"); !ok || v != 60 {
			t.Errorf("AUTO_SUSPEND = %d (ok=%v), want 60", v, ok)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE WAREHOUSE — trailing WITH TAG clause
// ---------------------------------------------------------------------------

func TestParseCreateWarehouse_Tags(t *testing.T) {
	t.Run("with tag", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE w WAREHOUSE_SIZE = SMALL WITH TAG (cost_center = 'eng', env = 'prod')")
		if len(stmt.Tags) != 2 {
			t.Fatalf("got %d tags, want 2: %+v", len(stmt.Tags), stmt.Tags)
		}
		if stmt.Tags[0].Name.String() != "cost_center" || stmt.Tags[0].Value != "eng" {
			t.Errorf("tag[0] = %+v", stmt.Tags[0])
		}
		if stmt.Tags[1].Name.String() != "env" || stmt.Tags[1].Value != "prod" {
			t.Errorf("tag[1] = %+v", stmt.Tags[1])
		}
		// The property before the tag clause must still be captured.
		if _, ok := optWords(stmt.Options, "WAREHOUSE_SIZE"); !ok {
			t.Error("WAREHOUSE_SIZE option lost when a TAG clause trails it")
		}
	})

	t.Run("bare tag without WITH", func(t *testing.T) {
		stmt := mustCreateWarehouse(t, "CREATE WAREHOUSE w TAG (t = 'v')")
		if len(stmt.Tags) != 1 || stmt.Tags[0].Value != "v" {
			t.Errorf("tags = %+v, want one tag t='v'", stmt.Tags)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER WAREHOUSE — every action
// ---------------------------------------------------------------------------

func TestParseAlterWarehouse_Actions(t *testing.T) {
	t.Run("suspend", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh SUSPEND")
		if stmt.Action != ast.AlterWarehouseSuspend {
			t.Errorf("Action = %d, want Suspend", stmt.Action)
		}
	})

	t.Run("resume", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh RESUME")
		if stmt.Action != ast.AlterWarehouseResume || stmt.ResumeIfSuspended {
			t.Errorf("Action=%d ResumeIfSuspended=%v, want Resume/false", stmt.Action, stmt.ResumeIfSuspended)
		}
	})

	t.Run("resume if suspended", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh RESUME IF SUSPENDED")
		if stmt.Action != ast.AlterWarehouseResume || !stmt.ResumeIfSuspended {
			t.Errorf("Action=%d ResumeIfSuspended=%v, want Resume/true", stmt.Action, stmt.ResumeIfSuspended)
		}
	})

	t.Run("abort all queries", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh ABORT ALL QUERIES")
		if stmt.Action != ast.AlterWarehouseAbort {
			t.Errorf("Action = %d, want Abort", stmt.Action)
		}
	})

	t.Run("rename to", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE IF EXISTS wh1 RENAME TO wh2")
		if !stmt.IfExists {
			t.Error("IfExists not set")
		}
		if stmt.Action != ast.AlterWarehouseRename {
			t.Fatalf("Action = %d, want Rename", stmt.Action)
		}
		if stmt.NewName == nil || stmt.NewName.String() != "wh2" {
			t.Errorf("NewName = %+v, want wh2", stmt.NewName)
		}
	})

	t.Run("set word value", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh SET warehouse_size=MEDIUM")
		if stmt.Action != ast.AlterWarehouseSet {
			t.Fatalf("Action = %d, want Set", stmt.Action)
		}
		if v, ok := optWords(stmt.Options, "WAREHOUSE_SIZE"); !ok || v != "MEDIUM" {
			t.Errorf("WAREHOUSE_SIZE = %q (ok=%v), want MEDIUM", v, ok)
		}
	})

	t.Run("set string value", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE so_warehouse SET RESOURCE_CONSTRAINT = 'MEMORY_16X_x86'")
		if v, ok := optLitString(stmt.Options, "RESOURCE_CONSTRAINT"); !ok || v != "MEMORY_16X_x86" {
			t.Errorf("RESOURCE_CONSTRAINT = %q (ok=%v)", v, ok)
		}
	})

	t.Run("set multiple props", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh SET GENERATION = '2' AUTO_SUSPEND = 120")
		if len(stmt.Options) != 2 {
			t.Fatalf("got %d options, want 2", len(stmt.Options))
		}
		if v, ok := optLitInt(stmt.Options, "AUTO_SUSPEND"); !ok || v != 120 {
			t.Errorf("AUTO_SUSPEND = %d (ok=%v), want 120", v, ok)
		}
	})

	t.Run("unset single", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh UNSET STATEMENT_TIMEOUT_IN_SECONDS")
		if stmt.Action != ast.AlterWarehouseUnset {
			t.Fatalf("Action = %d, want Unset", stmt.Action)
		}
		if len(stmt.UnsetKeys) != 1 || stmt.UnsetKeys[0] != "STATEMENT_TIMEOUT_IN_SECONDS" {
			t.Errorf("UnsetKeys = %v, want [STATEMENT_TIMEOUT_IN_SECONDS]", stmt.UnsetKeys)
		}
	})

	t.Run("unset multiple", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE my_wh UNSET auto_suspend, comment")
		want := []string{"AUTO_SUSPEND", "COMMENT"}
		if len(stmt.UnsetKeys) != len(want) {
			t.Fatalf("UnsetKeys = %v, want %v", stmt.UnsetKeys, want)
		}
		for i, w := range want {
			if stmt.UnsetKeys[i] != w {
				t.Errorf("UnsetKeys[%d] = %q, want %q", i, stmt.UnsetKeys[i], w)
			}
		}
	})

	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE w SET TAG cost_center = 'eng', env = 'prod'")
		if stmt.Action != ast.AlterWarehouseSetTag {
			t.Fatalf("Action = %d, want SetTag", stmt.Action)
		}
		if len(stmt.Tags) != 2 || stmt.Tags[0].Name.String() != "cost_center" || stmt.Tags[0].Value != "eng" {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		// Regression: UNSET TAG t must NOT be mis-parsed as UNSET keys "TAG", "t".
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE w UNSET TAG cost_center, env")
		if stmt.Action != ast.AlterWarehouseUnsetTag {
			t.Fatalf("Action = %d, want UnsetTag", stmt.Action)
		}
		if len(stmt.UnsetTags) != 2 || stmt.UnsetTags[0].String() != "cost_center" || stmt.UnsetTags[1].String() != "env" {
			t.Errorf("UnsetTags = %v", objNames(stmt.UnsetTags))
		}
		if len(stmt.UnsetKeys) != 0 {
			t.Errorf("UnsetKeys = %v, want empty (TAG must not leak into the key list)", stmt.UnsetKeys)
		}
	})

	t.Run("unset plain key named differently is not a tag", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE w UNSET COMMENT")
		if stmt.Action != ast.AlterWarehouseUnset {
			t.Fatalf("Action = %d, want Unset", stmt.Action)
		}
		if len(stmt.UnsetKeys) != 1 || stmt.UnsetKeys[0] != "COMMENT" {
			t.Errorf("UnsetKeys = %v, want [COMMENT]", stmt.UnsetKeys)
		}
	})

	t.Run("add tables", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE interactive_demo ADD TABLES (orders, customers)")
		if stmt.Action != ast.AlterWarehouseAddTables {
			t.Fatalf("Action = %d, want AddTables", stmt.Action)
		}
		if len(stmt.Tables) != 2 || stmt.Tables[0].String() != "orders" || stmt.Tables[1].String() != "customers" {
			t.Errorf("Tables = %v", objNames(stmt.Tables))
		}
	})

	t.Run("remove tables", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE interactive_demo REMOVE TABLES (orders, customers)")
		if stmt.Action != ast.AlterWarehouseRemoveTables {
			t.Fatalf("Action = %d, want RemoveTables", stmt.Action)
		}
		if len(stmt.Tables) != 2 {
			t.Errorf("Tables = %v, want 2", objNames(stmt.Tables))
		}
	})

	t.Run("drop tables aliases to remove", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE interactive_demo DROP TABLES (orders, customers)")
		if stmt.Action != ast.AlterWarehouseRemoveTables {
			t.Errorf("Action = %d, want RemoveTables (DROP TABLES is the same action)", stmt.Action)
		}
	})

	t.Run("add tables qualified", func(t *testing.T) {
		stmt := mustAlterWarehouse(t, "ALTER WAREHOUSE w ADD TABLES (db.sch.t1, sch.t2)")
		if len(stmt.Tables) != 2 || stmt.Tables[0].String() != "db.sch.t1" || stmt.Tables[1].String() != "sch.t2" {
			t.Errorf("Tables = %v", objNames(stmt.Tables))
		}
	})
}

func objNames(ns []*ast.ObjectName) []string {
	out := make([]string, len(ns))
	for i, n := range ns {
		out[i] = n.String()
	}
	return out
}

// ---------------------------------------------------------------------------
// Loc accuracy
// ---------------------------------------------------------------------------

func TestParseWarehouse_Loc(t *testing.T) {
	input := "CREATE WAREHOUSE my_wh"
	stmt := mustCreateWarehouse(t, input)
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(input) {
		t.Errorf("CREATE Loc = %+v, want {0, %d}", stmt.Loc, len(input))
	}

	// CREATE OR REPLACE: Loc.Start anchors at the CREATE keyword (offset 0).
	rinput := "CREATE OR REPLACE WAREHOUSE w WAREHOUSE_SIZE = SMALL"
	rstmt := mustCreateWarehouse(t, rinput)
	if rstmt.Loc.Start != 0 || rstmt.Loc.End != len(rinput) {
		t.Errorf("CREATE OR REPLACE Loc = %+v, want {0, %d}", rstmt.Loc, len(rinput))
	}

	// ALTER sub-parsers set Loc.Start at the object-type keyword (WAREHOUSE), not
	// the ALTER keyword, matching the established ALTER TABLE/STAGE convention.
	ainput := "ALTER WAREHOUSE w RENAME TO w2"
	const whKwOff = len("ALTER ")
	astmt := mustAlterWarehouse(t, ainput)
	if astmt.Loc.Start != whKwOff || astmt.Loc.End != len(ainput) {
		t.Errorf("ALTER Loc = %+v, want {%d, %d}", astmt.Loc, whKwOff, len(ainput))
	}
}

// ---------------------------------------------------------------------------
// Negative cases — malformed statements must produce a parse error, not hang.
// ---------------------------------------------------------------------------

func TestParseWarehouse_Errors(t *testing.T) {
	bad := []string{
		"CREATE WAREHOUSE",                       // missing name
		"CREATE WAREHOUSE w WAREHOUSE_SIZE =",    // option missing value
		"ALTER WAREHOUSE",                        // missing name + action
		"ALTER WAREHOUSE w",                      // missing action
		"ALTER WAREHOUSE w RENAME",               // missing TO
		"ALTER WAREHOUSE w RENAME TO",            // missing new name
		"ALTER WAREHOUSE w SET",                  // nothing to set
		"ALTER WAREHOUSE w SET WAREHOUSE_SIZE =", // set value missing
		"ALTER WAREHOUSE w SET TAG",              // missing tag assignment
		"ALTER WAREHOUSE w SET TAG t =",          // tag missing value
		"ALTER WAREHOUSE w UNSET TAG",            // missing tag name
		"ALTER WAREHOUSE w UNSET",                // missing key
		"ALTER WAREHOUSE w UNSET ,",              // empty key
		"ALTER WAREHOUSE w ABORT",                // missing ALL QUERIES
		"ALTER WAREHOUSE w ABORT ALL",            // missing QUERIES
		"ALTER WAREHOUSE w ADD",                  // missing TABLES
		"ALTER WAREHOUSE w ADD TABLES",           // missing ( ... )
		"ALTER WAREHOUSE w ADD TABLES ()",        // empty id list
		"ALTER WAREHOUSE w ADD TABLES (a,",       // unterminated list
		"ALTER WAREHOUSE w REMOVE TABLES",        // missing ( ... )
		"ALTER WAREHOUSE w FROBNICATE",           // unknown action
	}
	for _, in := range bad {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", in, len(result.File.Stmts))
			}
		})
	}
}
