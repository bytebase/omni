package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateSequence(t *testing.T, input string) *ast.CreateSequenceStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateSequenceStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateSequenceStmt", input, node)
	}
	return stmt
}

func mustAlterSequence(t *testing.T, input string) *ast.AlterSequenceStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterSequenceStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterSequenceStmt", input, node)
	}
	return stmt
}

func i64(v int64) *int64 { return &v }

// ---------------------------------------------------------------------------
// CREATE SEQUENCE
//
// Docs (truth1, authoritative):
//
//	CREATE [ OR REPLACE ] SEQUENCE [ IF NOT EXISTS ] <name>
//	  [ WITH ]
//	  [ START [ WITH ] [ = ] <initial_value> ]
//	  [ INCREMENT [ BY ] [ = ] <sequence_interval> ]
//	  [ { ORDER | NOORDER } ]
//	  [ COMMENT = '<string_literal>' ]
//
// Legacy ANTLR (truth2) — corroborates the optional connectors:
//
//	create_sequence: CREATE or_replace? SEQUENCE if_not_exists? object_name
//	  WITH? start_with? increment_by? order_noorder? comment_clause?
//	start_with:    START WITH? EQ? num
//	increment_by:  INCREMENT BY? EQ? num
// ---------------------------------------------------------------------------

func TestParseCreateSequence(t *testing.T) {
	t.Run("minimal name only", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE seq90")
		if stmt.Name.String() != "seq90" {
			t.Errorf("Name = %q, want seq90", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.IfNotExists {
			t.Errorf("unexpected modifier: %+v", stmt)
		}
		if stmt.Start != nil || stmt.Increment != nil || stmt.Order != nil || stmt.Comment != nil {
			t.Errorf("expected all clauses nil: %+v", stmt)
		}
	})

	t.Run("or replace with start and increment EQ", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE OR REPLACE SEQUENCE seq_01 START = 1 INCREMENT = 1")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
		if stmt.Start == nil || *stmt.Start != 1 {
			t.Errorf("Start = %v, want 1", stmt.Start)
		}
		if stmt.Increment == nil || *stmt.Increment != 1 {
			t.Errorf("Increment = %v, want 1", stmt.Increment)
		}
	})

	t.Run("increment 5", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE OR REPLACE SEQUENCE seq_5 START = 1 INCREMENT = 5")
		if stmt.Increment == nil || *stmt.Increment != 5 {
			t.Errorf("Increment = %v, want 5", stmt.Increment)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE IF NOT EXISTS s")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("leading WITH connector", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s WITH START 100 INCREMENT 2")
		if stmt.Start == nil || *stmt.Start != 100 {
			t.Errorf("Start = %v, want 100", stmt.Start)
		}
		if stmt.Increment == nil || *stmt.Increment != 2 {
			t.Errorf("Increment = %v, want 2", stmt.Increment)
		}
	})

	t.Run("START WITH spelling", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START WITH 42")
		if stmt.Start == nil || *stmt.Start != 42 {
			t.Errorf("Start = %v, want 42", stmt.Start)
		}
	})

	t.Run("INCREMENT BY spelling", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s INCREMENT BY 7")
		if stmt.Increment == nil || *stmt.Increment != 7 {
			t.Errorf("Increment = %v, want 7", stmt.Increment)
		}
	})

	t.Run("bare values no connector", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START 10 INCREMENT 3")
		if stmt.Start == nil || *stmt.Start != 10 || stmt.Increment == nil || *stmt.Increment != 3 {
			t.Errorf("Start=%v Increment=%v, want 10/3", stmt.Start, stmt.Increment)
		}
	})

	t.Run("negative increment", func(t *testing.T) {
		// A negative INCREMENT is documented and valid (a descending sequence).
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START WITH 100 INCREMENT BY -5")
		if stmt.Increment == nil || *stmt.Increment != -5 {
			t.Errorf("Increment = %v, want -5", stmt.Increment)
		}
	})

	t.Run("negative start", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START = -1 INCREMENT = 1")
		if stmt.Start == nil || *stmt.Start != -1 {
			t.Errorf("Start = %v, want -1", stmt.Start)
		}
	})

	t.Run("explicit positive sign", func(t *testing.T) {
		// A leading '+' is tolerated on START / INCREMENT values.
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START = +5 INCREMENT = +2")
		if stmt.Start == nil || *stmt.Start != 5 || stmt.Increment == nil || *stmt.Increment != 2 {
			t.Errorf("Start=%v Increment=%v, want 5/2", stmt.Start, stmt.Increment)
		}
	})

	t.Run("order", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START 1 INCREMENT 1 ORDER")
		if stmt.Order == nil || *stmt.Order != true {
			t.Errorf("Order = %v, want true", stmt.Order)
		}
	})

	t.Run("noorder", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s NOORDER")
		if stmt.Order == nil || *stmt.Order != false {
			t.Errorf("Order = %v, want false", stmt.Order)
		}
	})

	t.Run("comment", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE SEQUENCE s START 1 INCREMENT 1 ORDER COMMENT = 'a seq'")
		if stmt.Comment == nil || *stmt.Comment != "a seq" {
			t.Errorf("Comment = %v, want 'a seq'", stmt.Comment)
		}
	})

	t.Run("all clauses qualified name", func(t *testing.T) {
		stmt := mustCreateSequence(t, "CREATE OR REPLACE SEQUENCE db.sch.s WITH START WITH = 5 INCREMENT BY = 10 NOORDER COMMENT = 'x'")
		if stmt.Name.String() != "db.sch.s" {
			t.Errorf("Name = %q, want db.sch.s", stmt.Name.String())
		}
		if stmt.Start == nil || *stmt.Start != 5 || stmt.Increment == nil || *stmt.Increment != 10 {
			t.Errorf("Start=%v Increment=%v, want 5/10", stmt.Start, stmt.Increment)
		}
		if stmt.Order == nil || *stmt.Order != false || stmt.Comment == nil || *stmt.Comment != "x" {
			t.Errorf("Order=%v Comment=%v", stmt.Order, stmt.Comment)
		}
	})

	// Negatives.
	t.Run("reject: missing name", func(t *testing.T) {
		// A bare CREATE SEQUENCE with no name token fails (expected identifier).
		// (Note: START / INCREMENT are NOT reserved, so `CREATE SEQUENCE START`
		// parses with START as the sequence name — they are valid object names.)
		mustNotParse(t, "CREATE SEQUENCE")
	})
	t.Run("reject: numeric name", func(t *testing.T) {
		mustNotParse(t, "CREATE SEQUENCE 123")
	})
	t.Run("reject: START without value", func(t *testing.T) {
		mustNotParse(t, "CREATE SEQUENCE s START")
	})
	t.Run("reject: INCREMENT without value", func(t *testing.T) {
		mustNotParse(t, "CREATE SEQUENCE s INCREMENT")
	})
	t.Run("reject: START with non-integer", func(t *testing.T) {
		mustNotParse(t, "CREATE SEQUENCE s START = 'abc'")
	})
	t.Run("reject: COMMENT without =", func(t *testing.T) {
		// COMMENT in CREATE SEQUENCE requires '=' (comment_clause: COMMENT EQ string).
		mustNotParse(t, "CREATE SEQUENCE s COMMENT 'x'")
	})
}

// ---------------------------------------------------------------------------
// ALTER SEQUENCE
//
// Legacy ANTLR (truth2):
//
//	alter_sequence
//	  : ALTER SEQUENCE if_exists? object_name RENAME TO object_name
//	  | ALTER SEQUENCE if_exists? object_name SET? ( INCREMENT BY? EQ? num)?
//	  | ALTER SEQUENCE if_exists? object_name SET (order_noorder? comment_clause | order_noorder)
//	  | ALTER SEQUENCE if_exists? object_name UNSET COMMENT
// ---------------------------------------------------------------------------

func TestParseAlterSequence(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s RENAME TO s2")
		if stmt.Action != ast.AlterSequenceRename {
			t.Errorf("Action = %v, want AlterSequenceRename", stmt.Action)
		}
		if stmt.NewName == nil || stmt.NewName.String() != "s2" {
			t.Errorf("NewName = %v, want s2", stmt.NewName)
		}
	})

	t.Run("if exists rename", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE IF EXISTS s RENAME TO s2")
		if !stmt.IfExists || stmt.Action != ast.AlterSequenceRename {
			t.Errorf("IfExists=%v Action=%v", stmt.IfExists, stmt.Action)
		}
	})

	t.Run("bare increment by", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s INCREMENT BY 5")
		if stmt.Action != ast.AlterSequenceSetIncrement {
			t.Errorf("Action = %v, want AlterSequenceSetIncrement", stmt.Action)
		}
		if stmt.Increment == nil || *stmt.Increment != 5 {
			t.Errorf("Increment = %v, want 5", stmt.Increment)
		}
	})

	t.Run("bare increment negative", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s INCREMENT -2")
		if stmt.Increment == nil || *stmt.Increment != -2 {
			t.Errorf("Increment = %v, want -2", stmt.Increment)
		}
	})

	t.Run("set increment", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s SET INCREMENT BY = 10")
		if stmt.Action != ast.AlterSequenceSetIncrement {
			t.Errorf("Action = %v, want AlterSequenceSetIncrement", stmt.Action)
		}
		if stmt.Increment == nil || *stmt.Increment != 10 {
			t.Errorf("Increment = %v, want 10", stmt.Increment)
		}
	})

	t.Run("set order", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s SET ORDER")
		if stmt.Action != ast.AlterSequenceSet {
			t.Errorf("Action = %v, want AlterSequenceSet", stmt.Action)
		}
		if stmt.Order == nil || *stmt.Order != true {
			t.Errorf("Order = %v, want true", stmt.Order)
		}
	})

	t.Run("set noorder", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s SET NOORDER")
		if stmt.Order == nil || *stmt.Order != false {
			t.Errorf("Order = %v, want false", stmt.Order)
		}
	})

	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s SET COMMENT = 'hi'")
		if stmt.Action != ast.AlterSequenceSet {
			t.Errorf("Action = %v, want AlterSequenceSet", stmt.Action)
		}
		if stmt.Comment == nil || *stmt.Comment != "hi" {
			t.Errorf("Comment = %v, want 'hi'", stmt.Comment)
		}
	})

	t.Run("set order and comment", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s SET ORDER COMMENT = 'hi'")
		if stmt.Order == nil || *stmt.Order != true {
			t.Errorf("Order = %v, want true", stmt.Order)
		}
		if stmt.Comment == nil || *stmt.Comment != "hi" {
			t.Errorf("Comment = %v, want 'hi'", stmt.Comment)
		}
	})

	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterSequence(t, "ALTER SEQUENCE s UNSET COMMENT")
		if stmt.Action != ast.AlterSequenceUnsetComment {
			t.Errorf("Action = %v, want AlterSequenceUnsetComment", stmt.Action)
		}
	})

	// Negatives.
	t.Run("reject: SET nothing", func(t *testing.T) {
		mustNotParse(t, "ALTER SEQUENCE s SET")
	})
	t.Run("reject: RENAME without TO", func(t *testing.T) {
		mustNotParse(t, "ALTER SEQUENCE s RENAME s2")
	})
	t.Run("reject: UNSET non-comment", func(t *testing.T) {
		mustNotParse(t, "ALTER SEQUENCE s UNSET ORDER")
	})
	t.Run("reject: bad action", func(t *testing.T) {
		mustNotParse(t, "ALTER SEQUENCE s FROBNICATE")
	})
	t.Run("reject: SET INCREMENT without value", func(t *testing.T) {
		mustNotParse(t, "ALTER SEQUENCE s SET INCREMENT BY")
	})
}

// ---------------------------------------------------------------------------
// Official docs corpus — every CREATE SEQUENCE statement in the create-sequence
// corpus directory must parse with zero errors to *ast.CreateSequenceStmt. The
// official docs are the authoritative oracle (truth1). Interleaved statements
// owned by other DAG nodes (SELECT / INSERT / CREATE TABLE) are skipped.
// ---------------------------------------------------------------------------

func TestSequence_OfficialCorpus(t *testing.T) {
	dir := "testdata/official/create-sequence"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read corpus dir %s: %v", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			assertOwnedCreateParses(t, string(data), "SEQUENCE", func(n ast.Node) bool {
				_, ok := n.(*ast.CreateSequenceStmt)
				return ok
			})
		})
	}
}

// ---------------------------------------------------------------------------
// Walker integration — the sequence nodes must be reachable by ast.Inspect and
// visit their name (+ rename target) children.
// ---------------------------------------------------------------------------

func TestSequence_WalkerVisitsChildren(t *testing.T) {
	cases := []string{
		"CREATE SEQUENCE s START 1 INCREMENT 1",
		"ALTER SEQUENCE s RENAME TO s2",
		"ALTER SEQUENCE s SET COMMENT = 'x'",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			node := mustParseOne(t, input)
			count := 0
			ast.Inspect(node, func(n ast.Node) bool {
				if n != nil {
					count++
				}
				return true
			})
			if count < 2 {
				t.Errorf("Inspect visited %d nodes, want >= 2 (root + name)", count)
			}
		})
	}
}
