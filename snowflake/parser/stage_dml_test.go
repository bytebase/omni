package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// PUT
// ---------------------------------------------------------------------------

func mustPut(t *testing.T, input string) *ast.PutStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.PutStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.PutStmt", input, node)
	}
	return stmt
}

func TestParsePut(t *testing.T) {
	t.Run("basic linux path to named stage", func(t *testing.T) {
		stmt := mustPut(t, "PUT file:///tmp/data/mydata.csv @my_int_stage")
		if stmt.File.Kind != ast.StageLocalFile || stmt.File.Raw != "file:///tmp/data/mydata.csv" {
			t.Fatalf("File = %+v", stmt.File)
		}
		if stmt.Stage.Kind != ast.StageRef || stmt.Stage.Raw != "@my_int_stage" {
			t.Fatalf("Stage = %+v", stmt.Stage)
		}
	})

	t.Run("table stage destination with auto_compress", func(t *testing.T) {
		stmt := mustPut(t, "PUT file:///tmp/data/orders_001.csv @%orderstiny_ext AUTO_COMPRESS = FALSE")
		if stmt.Stage.Raw != "@%orderstiny_ext" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
		opt := findOption(stmt.Options, "AUTO_COMPRESS")
		if opt == nil || opt.Words != "FALSE" {
			t.Errorf("AUTO_COMPRESS = %+v", opt)
		}
	})

	t.Run("glob source", func(t *testing.T) {
		stmt := mustPut(t, "PUT file:///tmp/data/orders_*01.csv @my_int_stage AUTO_COMPRESS = FALSE")
		if stmt.File.Raw != "file:///tmp/data/orders_*01.csv" {
			t.Errorf("File.Raw = %q", stmt.File.Raw)
		}
	})

	t.Run("double-star glob no spaces around equals", func(t *testing.T) {
		stmt := mustPut(t, "PUT file:///tmp/data/** @my_int_stage AUTO_COMPRESS=FALSE")
		if stmt.File.Raw != "file:///tmp/data/**" {
			t.Errorf("File.Raw = %q", stmt.File.Raw)
		}
		if findOption(stmt.Options, "AUTO_COMPRESS") == nil {
			t.Errorf("missing AUTO_COMPRESS")
		}
	})

	t.Run("quoted path with space", func(t *testing.T) {
		stmt := mustPut(t, "PUT 'file:///tmp/data/orders 001.csv' @my_int_stage AUTO_COMPRESS = FALSE")
		if stmt.File.Raw != "file:///tmp/data/orders 001.csv" {
			t.Errorf("File.Raw = %q", stmt.File.Raw)
		}
	})

	t.Run("windows path user stage", func(t *testing.T) {
		stmt := mustPut(t, `PUT file://C:\\temp\\data\\mydata.csv @~ AUTO_COMPRESS = TRUE`)
		if stmt.Stage.Raw != "@~" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
		if stmt.File.Kind != ast.StageLocalFile {
			t.Errorf("File.Kind = %v", stmt.File.Kind)
		}
	})

	t.Run("all transfer options", func(t *testing.T) {
		stmt := mustPut(t, "PUT file:///x.csv @s PARALLEL = 8 AUTO_COMPRESS = TRUE SOURCE_COMPRESSION = GZIP OVERWRITE = FALSE")
		for _, name := range []string{"PARALLEL", "AUTO_COMPRESS", "SOURCE_COMPRESSION", "OVERWRITE"} {
			if findOption(stmt.Options, name) == nil {
				t.Errorf("missing option %s", name)
			}
		}
		if opt := findOption(stmt.Options, "PARALLEL"); opt.Lit == nil || opt.Lit.Ival != 8 {
			t.Errorf("PARALLEL = %+v", opt)
		}
		if opt := findOption(stmt.Options, "SOURCE_COMPRESSION"); opt.Words != "GZIP" {
			t.Errorf("SOURCE_COMPRESSION = %+v", opt)
		}
	})
}

// ---------------------------------------------------------------------------
// GET
// ---------------------------------------------------------------------------

func mustGet(t *testing.T, input string) *ast.GetStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.GetStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.GetStmt", input, node)
	}
	return stmt
}

func TestParseGet(t *testing.T) {
	t.Run("table stage to linux dir", func(t *testing.T) {
		stmt := mustGet(t, "GET @%mytable file:///tmp/data/")
		if stmt.Stage.Raw != "@%mytable" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
		if stmt.Target.Kind != ast.StageLocalFile || stmt.Target.Raw != "file:///tmp/data/" {
			t.Fatalf("Target = %+v", stmt.Target)
		}
	})

	t.Run("user stage path", func(t *testing.T) {
		stmt := mustGet(t, "GET @~/myfiles file:///tmp/data/")
		if stmt.Stage.Raw != "@~/myfiles" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("bare target path with double-quoted pattern", func(t *testing.T) {
		// Docs example: GET @my_int_stage my_target_path PATTERN = "tmp.parquet".
		// The target is a bare path (not file://) and PATTERN's value is a
		// double-quoted identifier.
		stmt := mustGet(t, `GET @my_int_stage my_target_path PATTERN = "tmp.parquet"`)
		if stmt.Target.Raw != "my_target_path" {
			t.Errorf("Target.Raw = %q", stmt.Target.Raw)
		}
		opt := findOption(stmt.Options, "PATTERN")
		if opt == nil || opt.Lit == nil || opt.Lit.Value != "tmp.parquet" {
			t.Fatalf("PATTERN = %+v", opt)
		}
	})

	t.Run("parallel option", func(t *testing.T) {
		stmt := mustGet(t, "GET @s file:///tmp/ PARALLEL = 10")
		if opt := findOption(stmt.Options, "PARALLEL"); opt == nil || opt.Lit.Ival != 10 {
			t.Errorf("PARALLEL = %+v", opt)
		}
	})
}

// ---------------------------------------------------------------------------
// LIST / LS
// ---------------------------------------------------------------------------

func mustList(t *testing.T, input string) *ast.ListStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ListStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ListStmt", input, node)
	}
	return stmt
}

func TestParseList(t *testing.T) {
	t.Run("named stage", func(t *testing.T) {
		stmt := mustList(t, "LIST @my_gcs_stage")
		if stmt.Short {
			t.Errorf("Short = true, want false")
		}
		if stmt.Stage.Raw != "@my_gcs_stage" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("table stage path", func(t *testing.T) {
		stmt := mustList(t, "LIST @%mytable/path")
		if stmt.Stage.Raw != "@%mytable/path" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("user stage", func(t *testing.T) {
		stmt := mustList(t, "LIST @~")
		if stmt.Stage.Raw != "@~" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("ls alias", func(t *testing.T) {
		stmt := mustList(t, "LS @%t1")
		if !stmt.Short {
			t.Errorf("Short = false, want true (LS)")
		}
		if stmt.Stage.Raw != "@%t1" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("with pattern", func(t *testing.T) {
		stmt := mustList(t, "LIST @mystage PATTERN = '.*[.]csv'")
		if stmt.Pattern == nil || stmt.Pattern.Value != ".*[.]csv" {
			t.Fatalf("Pattern = %+v", stmt.Pattern)
		}
	})

	t.Run("qualified stage with pattern", func(t *testing.T) {
		stmt := mustList(t, "LIST @db.schema.mystage/p PATTERN = 'a.*'")
		if stmt.Stage.Raw != "@db.schema.mystage/p" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
		if stmt.Pattern == nil || stmt.Pattern.Value != "a.*" {
			t.Errorf("Pattern = %+v", stmt.Pattern)
		}
	})

	t.Run("double-slash in stage path is truncated by the lexer", func(t *testing.T) {
		// KNOWN lexer limitation (flagged divergence): omni's lexer treats '//'
		// as a line comment, so a stage path containing '//' is truncated at the
		// '//'. This is a lexer-node concern, not a dml-copy parser bug — the
		// token-contiguity stage reader faithfully captures what the lexer emits.
		// Documented here so the behavior is tracked rather than silent; '//' in a
		// stage path is not present in the official corpus or legacy grammar.
		stmt := mustList(t, "LIST @s/a//b")
		if stmt.Stage.Raw != "@s/a" {
			t.Errorf("Stage.Raw = %q; lexer-truncation behavior changed — re-evaluate", stmt.Stage.Raw)
		}
	})
}

// ---------------------------------------------------------------------------
// REMOVE / RM
// ---------------------------------------------------------------------------

func mustRemove(t *testing.T, input string) *ast.RemoveStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.RemoveStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.RemoveStmt", input, node)
	}
	return stmt
}

func TestParseRemove(t *testing.T) {
	t.Run("named stage", func(t *testing.T) {
		stmt := mustRemove(t, "REMOVE @mystage")
		if stmt.Short {
			t.Errorf("Short = true, want false")
		}
		if stmt.Stage.Raw != "@mystage" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("rm alias table stage", func(t *testing.T) {
		stmt := mustRemove(t, "RM @%mytable")
		if !stmt.Short {
			t.Errorf("Short = false, want true (RM)")
		}
		if stmt.Stage.Raw != "@%mytable" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
	})

	t.Run("user stage with path and pattern", func(t *testing.T) {
		stmt := mustRemove(t, "REMOVE @~/staged PATTERN = '.*[.]csv'")
		if stmt.Stage.Raw != "@~/staged" {
			t.Errorf("Stage.Raw = %q", stmt.Stage.Raw)
		}
		if stmt.Pattern == nil || stmt.Pattern.Value != ".*[.]csv" {
			t.Errorf("Pattern = %+v", stmt.Pattern)
		}
	})
}

// ---------------------------------------------------------------------------
// Negative tests
// ---------------------------------------------------------------------------

func TestParseStageDML_Negative(t *testing.T) {
	bad := []string{
		"PUT @stage @other",        // PUT source must be a file, not a stage
		"PUT file:///x.csv",        // PUT missing destination stage
		"GET file:///x.csv @stage", // GET source must be a stage, not a file
		"GET @stage",               // GET missing destination
		"LIST",                     // LIST missing stage
		"LIST mytable",             // LIST operand must be a stage (@...)
		"REMOVE",                   // REMOVE missing stage
		"LIST @s PATTERN =",        // PATTERN missing value
		"LIST @s PATTERN = 123",    // PATTERN value must be a string/quoted ident
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

// ---------------------------------------------------------------------------
// Loc accuracy with a non-zero base offset (raw file-path scanning must respect
// the segment base offset).
// ---------------------------------------------------------------------------

func TestParsePut_BaseOffset(t *testing.T) {
	multi := "SELECT 1;\nPUT file:///tmp/x.csv @s AUTO_COMPRESS = TRUE"
	res := ParseBestEffort(multi)
	if len(res.Errors) != 0 {
		t.Fatalf("multi parse errors: %v", res.Errors)
	}
	put, ok := res.File.Stmts[1].(*ast.PutStmt)
	if !ok {
		t.Fatalf("stmt[1] = %T", res.File.Stmts[1])
	}
	if put.File.Raw != "file:///tmp/x.csv" {
		t.Errorf("File.Raw = %q, want file:///tmp/x.csv (base offset wrong)", put.File.Raw)
	}
	if put.Stage.Raw != "@s" {
		t.Errorf("Stage.Raw = %q, want @s", put.Stage.Raw)
	}
	// Loc must point at the absolute position in the multi-statement input.
	wantStart := len("SELECT 1;\n")
	if put.Loc.Start != wantStart {
		t.Errorf("Loc.Start = %d, want %d", put.Loc.Start, wantStart)
	}
}
