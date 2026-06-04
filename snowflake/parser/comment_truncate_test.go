package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// COMMENT ON
// ---------------------------------------------------------------------------

func TestParseComment(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantType     string // ObjectType
		wantColumn   bool
		wantName     string // ObjectName.String()
		wantIfExists bool
		wantComment  string // raw string content (without quotes)
		wantSig      []string
	}{
		{"table", "COMMENT ON TABLE t IS 't'", "TABLE", false, "t", false, "t", nil},
		{"database", "COMMENT ON DATABASE d IS 'd'", "DATABASE", false, "d", false, "d", nil},
		{"file format two-word", "COMMENT ON FILE FORMAT f IS 'f'", "FILE FORMAT", false, "f", false, "f", nil},
		{"masking policy", "COMMENT ON MASKING POLICY m IS 'm'", "MASKING POLICY", false, "m", false, "m", nil},
		{"row access policy", "COMMENT ON ROW ACCESS POLICY r IS 'r'", "ROW ACCESS POLICY", false, "r", false, "r", nil},
		{"session policy", "COMMENT ON SESSION POLICY s IS 's'", "SESSION POLICY", false, "s", false, "s", nil},
		{"warehouse", "COMMENT ON WAREHOUSE w IS 'w'", "WAREHOUSE", false, "w", false, "w", nil},
		{"if exists", "COMMENT IF EXISTS ON TABLE t IS 't'", "TABLE", false, "t", true, "t", nil},
		{"column 1-part", "COMMENT ON COLUMN c IS 'c'", "COLUMN", true, "c", false, "c", nil},
		{"column 2-part", "COMMENT ON COLUMN t.c IS 't.c'", "COLUMN", true, "t.c", false, "t.c", nil},
		{"column 3-part", "COMMENT ON COLUMN s.t.c IS 's.t.c'", "COLUMN", true, "s.t.c", false, "s.t.c", nil},
		{"column 4-part", "COMMENT ON COLUMN d.s.t.c IS 'd.s.t.c'", "COLUMN", true, "d.s.t.c", false, "d.s.t.c", nil},
		{"qualified name", "COMMENT ON TABLE db.sch.t IS 'x'", "TABLE", false, "db.sch.t", false, "x", nil},
		{"function with sig", "COMMENT ON FUNCTION f(INT, STRING) IS 'fn'", "FUNCTION", false, "f", false, "fn", []string{"INT", "STRING"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.CommentStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.CommentStmt", node)
			}
			if stmt.ObjectType != tt.wantType {
				t.Errorf("ObjectType = %q, want %q", stmt.ObjectType, tt.wantType)
			}
			if stmt.IsColumn != tt.wantColumn {
				t.Errorf("IsColumn = %v, want %v", stmt.IsColumn, tt.wantColumn)
			}
			if stmt.IfExists != tt.wantIfExists {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tt.wantIfExists)
			}
			if tt.wantColumn {
				if stmt.Name != nil {
					t.Errorf("Name = %v, want nil for column form", stmt.Name)
				}
				if got := columnRefString(stmt.Column); got != tt.wantName {
					t.Errorf("Column = %q, want %q", got, tt.wantName)
				}
			} else {
				if stmt.Column != nil {
					t.Errorf("Column = %v, want nil for object form", stmt.Column)
				}
				if stmt.Name == nil || stmt.Name.String() != tt.wantName {
					t.Errorf("Name = %v, want %q", stmt.Name, tt.wantName)
				}
			}
			if got := stringContent(stmt.Comment); got != tt.wantComment {
				t.Errorf("Comment = %q (content %q), want %q", stmt.Comment, got, tt.wantComment)
			}
			if !eqStrings(sigNames(stmt.Signature), tt.wantSig) {
				t.Errorf("Signature = %v, want %v", sigNames(stmt.Signature), tt.wantSig)
			}
		})
	}
}

// columnRefString renders a ColumnRef's parts as a dotted source string for
// comparison in tests. Returns "" for a nil ref.
func columnRefString(c *ast.ColumnRef) string {
	if c == nil {
		return ""
	}
	parts := make([]string, len(c.Parts))
	for i, p := range c.Parts {
		parts[i] = p.String()
	}
	return strings.Join(parts, ".")
}

// stringContent strips a single pair of surrounding single quotes from a raw
// SQL string literal token, for comparison in tests.
func stringContent(raw string) string {
	if len(raw) >= 2 && raw[0] == '\'' && raw[len(raw)-1] == '\'' {
		return raw[1 : len(raw)-1]
	}
	return raw
}

// ---------------------------------------------------------------------------
// TRUNCATE
// ---------------------------------------------------------------------------

func TestParseTruncate(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantName     string
		wantIfExists bool
		wantMatView  bool
	}{
		{"table keyword", "TRUNCATE TABLE t", "t", false, false},
		{"bare", "TRUNCATE t", "t", false, false},
		{"if exists", "TRUNCATE TABLE IF EXISTS t", "t", true, false},
		{"bare if exists", "TRUNCATE IF EXISTS t", "t", true, false},
		{"qualified", "TRUNCATE TABLE db.sch.t", "db.sch.t", false, false},
		{"materialized view", "TRUNCATE MATERIALIZED VIEW mv", "mv", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.TruncateStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.TruncateStmt", node)
			}
			if stmt.MaterializedView != tt.wantMatView {
				t.Errorf("MaterializedView = %v, want %v", stmt.MaterializedView, tt.wantMatView)
			}
			if stmt.IfExists != tt.wantIfExists {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tt.wantIfExists)
			}
			if stmt.ErrorTable {
				t.Errorf("ErrorTable = true, want false")
			}
			if stmt.Name == nil || stmt.Name.String() != tt.wantName {
				t.Errorf("Name = %v, want %q", stmt.Name, tt.wantName)
			}
		})
	}
}

// TestParseTruncateErrorTable covers the documented (docs = truth1) form
// TRUNCATE [TABLE] [IF EXISTS] ERROR_TABLE(<base_table_name>), which the legacy
// ANTLR grammar omits.
func TestParseTruncateErrorTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantIf   bool
	}{
		{"basic", "TRUNCATE TABLE ERROR_TABLE(my_base)", "my_base", false},
		{"no table kw", "TRUNCATE ERROR_TABLE(my_base)", "my_base", false},
		{"if exists", "TRUNCATE TABLE IF EXISTS ERROR_TABLE(db.sch.t)", "db.sch.t", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.TruncateStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.TruncateStmt", node)
			}
			if !stmt.ErrorTable {
				t.Errorf("ErrorTable = false, want true")
			}
			if stmt.IfExists != tt.wantIf {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tt.wantIf)
			}
			if stmt.Name == nil || stmt.Name.String() != tt.wantName {
				t.Errorf("Name = %v, want %q", stmt.Name, tt.wantName)
			}
		})
	}
}

func TestParseCommentTruncate_Errors(t *testing.T) {
	cases := []string{
		"COMMENT",                            // nothing
		"COMMENT ON TABLE t",                 // missing IS '...'
		"COMMENT ON TABLE t IS",              // missing string
		"COMMENT ON TABLE IS 'x'",            // missing object type/name
		"COMMENT ON COLUMN IS 'x'",           // COLUMN without a name
		"COMMENT ON COLUMN a.b.c.d.e IS 'x'", // 5-part column name (max is 4)
		"TRUNCATE",                           // nothing
		"TRUNCATE TABLE",                     // missing name
		"TRUNCATE MATERIALIZED VIEW",         // missing name
		"TRUNCATE MATERIALIZED t",            // MATERIALIZED without VIEW
		"TRUNCATE TABLE ERROR_TABLE(",        // ERROR_TABLE missing arg + close paren
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			result := ParseBestEffort(c)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", c, len(result.File.Stmts))
			}
		})
	}
}
