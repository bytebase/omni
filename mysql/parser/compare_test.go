package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// ParseAndCheck parses sql and verifies basic properties:
// - Parse succeeds for valid SQL
// - All Loc.Start fields on AST nodes are non-negative
func ParseAndCheck(t *testing.T, sql string) *ast.List {
	t.Helper()
	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if result == nil {
		t.Fatalf("Parse(%q) returned nil", sql)
	}
	// Walk AST and check all Loc.Start fields >= 0
	for i, item := range result.Items {
		checkLocations(t, sql, i, item)
	}
	return result
}

// ParseExpectError verifies that parsing the given SQL produces an error.
func ParseExpectError(t *testing.T, sql string) {
	t.Helper()
	_, err := Parse(sql)
	if err == nil {
		t.Fatalf("Parse(%q) expected error, got nil", sql)
	}
}

// ParseAndCompare parses sql and compares the AST string representation
// against the expected string.
func ParseAndCompare(t *testing.T, sql string, expected string) {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() == 0 {
		t.Fatalf("Parse(%q) returned empty list", sql)
	}
	got := ast.NodeToString(result.Items[0])
	if got != expected {
		t.Errorf("Parse(%q):\n  got:  %s\n  want: %s", sql, got, expected)
	}
}

// checkLocations recursively checks that Loc.Start fields are non-negative.
func checkLocations(t *testing.T, sql string, idx int, node ast.Node) {
	t.Helper()
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.SelectStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] SelectStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.InsertStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] InsertStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.UpdateStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] UpdateStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.DeleteStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] DeleteStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.CreateTableStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] CreateTableStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.IntLit:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] IntLit.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.StringLit:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] StringLit.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.BinaryExpr:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] BinaryExpr.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
		checkLocations(t, sql, idx, n.Left)
		checkLocations(t, sql, idx, n.Right)
	case *ast.UnaryExpr:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] UnaryExpr.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
		checkLocations(t, sql, idx, n.Operand)
	case *ast.ColumnRef:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] ColumnRef.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.FuncCallExpr:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] FuncCallExpr.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.CaseExpr:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] CaseExpr.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
	case *ast.RawStmt:
		if n.Loc.Start < 0 {
			t.Errorf("Parse(%q) stmt[%d] RawStmt.Loc.Start = %d, want >= 0", sql, idx, n.Loc.Start)
		}
		checkLocations(t, sql, idx, n.Stmt)
	}
}

// TestLexerBasic tests the lexer with basic tokens.
func TestLexerBasic(t *testing.T) {
	tests := []struct {
		input string
		types []int
	}{
		{"SELECT 1", []int{kwSELECT, tokICONST}},
		{"SELECT 'hello'", []int{kwSELECT, tokSCONST}},
		{"SELECT `col`", []int{kwSELECT, tokIDENT}},
		{"SELECT @var", []int{kwSELECT, tokIDENT}},
		{"SELECT @@global.var", []int{kwSELECT, tokAt2, kwGLOBAL, '.', tokIDENT}},
		{"1 + 2", []int{tokICONST, '+', tokICONST}},
		{"1 <= 2", []int{tokICONST, tokLessEq, tokICONST}},
		{"1 >= 2", []int{tokICONST, tokGreaterEq, tokICONST}},
		{"1 <> 2", []int{tokICONST, tokNotEq, tokICONST}},
		{"1 != 2", []int{tokICONST, tokNotEq, tokICONST}},
		{"1 <=> 2", []int{tokICONST, tokNullSafeEq, tokICONST}},
		{"1 << 2", []int{tokICONST, tokShiftLeft, tokICONST}},
		{"1 >> 2", []int{tokICONST, tokShiftRight, tokICONST}},
		{"a := 1", []int{tokIDENT, tokAssign, tokICONST}},
		{"0xFF", []int{tokXCONST}},
		{"0b101", []int{tokBCONST}},
		{"X'FF'", []int{tokXCONST}},
		{"b'101'", []int{tokBCONST}},
		{"3.14", []int{tokFCONST}},
		{".5", []int{tokFCONST}},
		{"1e10", []int{tokFCONST}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := NewLexer(tt.input)
			for i, want := range tt.types {
				tok := lex.NextToken()
				if tok.Type != want {
					t.Errorf("token[%d] type = %d, want %d (input=%q)", i, tok.Type, want, tt.input)
				}
			}
			// Should be at EOF
			tok := lex.NextToken()
			if tok.Type != tokEOF {
				t.Errorf("expected EOF, got type=%d (input=%q)", tok.Type, tt.input)
			}
		})
	}
}

// TestLexerComments tests that comments are properly skipped.
func TestLexerComments(t *testing.T) {
	tests := []struct {
		input string
		types []int
	}{
		{"SELECT -- comment\n1", []int{kwSELECT, tokICONST}},
		{"SELECT # comment\n1", []int{kwSELECT, tokICONST}},
		{"SELECT /* comment */ 1", []int{kwSELECT, tokICONST}},
		{"SELECT /* nested /* comment */ */ 1", []int{kwSELECT, tokICONST}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := NewLexer(tt.input)
			for i, want := range tt.types {
				tok := lex.NextToken()
				if tok.Type != want {
					t.Errorf("token[%d] type = %d, want %d (input=%q)", i, tok.Type, want, tt.input)
				}
			}
		})
	}
}

// TestLexerStrings tests string literal scanning.
func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input string
		value string
	}{
		{`'hello'`, "hello"},
		{`'it''s'`, "it's"},
		{`"double"`, "double"},
		{`'back\\slash'`, "back\\slash"},
		{`'new\nline'`, "new\nline"},
		{`'tab\there'`, "tab\there"},
		{`'cr\rhere'`, "cr\rhere"},
		{`'null\0here'`, "null\x00here"},
		{`'bs\bhere'`, "bs\x08here"},
		{`'ctrlz\Zhere'`, "ctrlz\x1ahere"},
		{`'like\_pattern'`, `like\_pattern`},
		{`'like\%pattern'`, `like\%pattern`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := NewLexer(tt.input)
			tok := lex.NextToken()
			if tok.Type != tokSCONST {
				t.Errorf("type = %d, want tokSCONST", tok.Type)
			}
			if tok.Str != tt.value {
				t.Errorf("value = %q, want %q", tok.Str, tt.value)
			}
		})
	}
}

// TestLexerBacktickIdent tests backtick-quoted identifier scanning.
func TestLexerBacktickIdent(t *testing.T) {
	tests := []struct {
		input string
		value string
	}{
		{"`column`", "column"},
		{"`col``name`", "col`name"},
		{"`table name`", "table name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := NewLexer(tt.input)
			tok := lex.NextToken()
			if tok.Type != tokIDENT {
				t.Errorf("type = %d, want tokIDENT", tok.Type)
			}
			if tok.Str != tt.value {
				t.Errorf("value = %q, want %q", tok.Str, tt.value)
			}
		})
	}
}

// TestLexerKeywords tests keyword recognition.
func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input string
		typ   int
	}{
		{"SELECT", kwSELECT},
		{"select", kwSELECT},
		{"SeLeCt", kwSELECT},
		{"INSERT", kwINSERT},
		{"UPDATE", kwUPDATE},
		{"DELETE", kwDELETE},
		{"CREATE", kwCREATE},
		{"ALTER", kwALTER},
		{"DROP", kwDROP},
		{"TABLE", kwTABLE},
		{"FROM", kwFROM},
		{"WHERE", kwWHERE},
		{"AND", kwAND},
		{"OR", kwOR},
		{"NOT", kwNOT},
		{"NULL", kwNULL},
		{"TRUE", kwTRUE},
		{"FALSE", kwFALSE},
		{"DIV", kwDIV},
		{"MOD", kwMOD},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lex := NewLexer(tt.input)
			tok := lex.NextToken()
			if tok.Type != tt.typ {
				t.Errorf("keyword %q: type = %d, want %d", tt.input, tok.Type, tt.typ)
			}
		})
	}
}

// TestLexerDotContext verifies that keyword lookup is suppressed after '.',
// matching MySQL 8.0 behaviour (sql_lex.cc MY_LEX_IDENT_START state).
// Any word after a dot must be returned as tokIDENT, even if it is a reserved word.
func TestLexerDotContext(t *testing.T) {
	// Lexer-level: after '.', the next identifier must be tokIDENT.
	reservedWords := []string{
		"select", "from", "table", "where", "key", "index", "group", "order",
	}
	for _, kw := range reservedWords {
		t.Run("lex_dot_"+kw, func(t *testing.T) {
			lex := NewLexer("t." + kw)
			tok1 := lex.NextToken() // t
			if tok1.Type != tokIDENT {
				t.Fatalf("expected tokIDENT for 't', got %d", tok1.Type)
			}
			tok2 := lex.NextToken() // .
			if tok2.Type != '.' {
				t.Fatalf("expected '.' token, got %d", tok2.Type)
			}
			tok3 := lex.NextToken() // keyword-as-ident
			if tok3.Type != tokIDENT {
				t.Errorf("after dot, %q should be tokIDENT (%d), got %d", kw, tokIDENT, tok3.Type)
			}
		})
	}

	// Parser-level: full SELECT statements with reserved words after dot.
	parserTests := []struct {
		name  string
		input string
	}{
		{"select", "SELECT t.select FROM t"},
		{"from", "SELECT t.from FROM t"},
		{"table", "SELECT t.table FROM t"},
		{"where", "SELECT t.where FROM t"},
		{"key", "SELECT t.key FROM t"},
		{"index", "SELECT t.index FROM t"},
		{"group", "SELECT t.group FROM t"},
		{"order", "SELECT t.order FROM t"},
		{"schema_qual", "SELECT db.select FROM db.select"},
		{"three_part", "SELECT a.b.c FROM t a"},
		{"backtick_after_dot", "SELECT t.`select` FROM t"},
	}
	for _, tt := range parserTests {
		t.Run("parse_"+tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tt.input, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 1: names — identifiers, qualified names, variable refs
// ---------------------------------------------------------------------------

// TestParseIdentifiers tests simple and backtick-quoted identifier parsing.
func TestParseIdentifiers(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"col1", "col1"},
		{"`my col`", "my col"},
		{"`select`", "select"},
		{"_myvar", "_myvar"},
		{"$dollar", "$dollar"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			name, _, err := p.parseIdentifier()
			if err != nil {
				t.Fatalf("parseIdentifier(%q) error: %v", tt.input, err)
			}
			if name != tt.want {
				t.Errorf("parseIdentifier(%q) = %q, want %q", tt.input, name, tt.want)
			}
		})
	}
}

// TestParseQualifiedNames tests column refs (unqualified, table.col, schema.table.col, table.*, schema.table.*).
func TestParseQualifiedNames(t *testing.T) {
	tests := []struct {
		input string
		want  string // expected NodeToString
	}{
		{"col1", "{COLREF :loc 0 :col col1}"},
		{"`col`", "{COLREF :loc 0 :col col}"},
		{"t.col", "{COLREF :loc 0 :table t :col col}"},
		{"`db`.`tbl`.`col`", "{COLREF :loc 0 :schema db :table tbl :col col}"},
		{"t.*", "{COLREF :loc 0 :table t :star true}"},
		{"db.t.*", "{COLREF :loc 0 :schema db :table t :star true}"},
		// Keywords as identifiers in qualified names
		{"t.status", "{COLREF :loc 0 :table t :col status}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			ref, err := p.parseColumnRef()
			if err != nil {
				t.Fatalf("parseColumnRef(%q) error: %v", tt.input, err)
			}
			got := ast.NodeToString(ref)
			if got != tt.want {
				t.Errorf("parseColumnRef(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseTableRef tests table reference parsing (simple and schema-qualified).
func TestParseTableRef(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"users", "{TABLEREF :loc 0 :name users}"},
		{"`users`", "{TABLEREF :loc 0 :name users}"},
		{"mydb.users", "{TABLEREF :loc 0 :schema mydb :name users}"},
		{"`my db`.`my table`", "{TABLEREF :loc 0 :schema my db :name my table}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				t.Fatalf("parseTableRef(%q) error: %v", tt.input, err)
			}
			got := ast.NodeToString(ref)
			if got != tt.want {
				t.Errorf("parseTableRef(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseTableRefWithAlias tests table reference with alias parsing.
func TestParseTableRefWithAlias(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"users AS u", "{TABLEREF :loc 0 :name users :alias u}"},
		{"users u", "{TABLEREF :loc 0 :name users :alias u}"},
		{"mydb.users AS u", "{TABLEREF :loc 0 :schema mydb :name users :alias u}"},
		{"users", "{TABLEREF :loc 0 :name users}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			ref, err := p.parseTableRefWithAlias()
			if err != nil {
				t.Fatalf("parseTableRefWithAlias(%q) error: %v", tt.input, err)
			}
			got := ast.NodeToString(ref)
			if got != tt.want {
				t.Errorf("parseTableRefWithAlias(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseVariableRefs tests user and system variable reference parsing.
func TestParseVariableRefs(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"@myvar", "{VAR :name myvar :loc 0}"},
		{"@@version", "{VAR :name version :loc 0 :system true}"},
		{"@@global.max_connections", "{VAR :name max_connections :loc 0 :system true :scope GLOBAL}"},
		{"@@session.sql_mode", "{VAR :name sql_mode :loc 0 :system true :scope SESSION}"},
		{"@@local.wait_timeout", "{VAR :name wait_timeout :loc 0 :system true :scope LOCAL}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			ref, err := p.parseVariableRef()
			if err != nil {
				t.Fatalf("parseVariableRef(%q) error: %v", tt.input, err)
			}
			got := ast.NodeToString(ref)
			if got != tt.want {
				t.Errorf("parseVariableRef(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseIdentifierError tests that non-identifier tokens produce errors.
func TestParseIdentifierError(t *testing.T) {
	inputs := []string{
		"123",   // number
		"'str'", // string literal
		"(",     // punctuation
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(input)}
			p.advance()
			_, _, err := p.parseIdentifier()
			if err == nil {
				t.Errorf("parseIdentifier(%q) expected error, got nil", input)
			}
		})
	}
}

// TestParseVariableRefError tests that non-variable tokens produce errors.
func TestParseVariableRefError(t *testing.T) {
	inputs := []string{
		"col",   // plain identifier
		"123",   // number
		"'str'", // string literal
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(input)}
			p.advance()
			_, err := p.parseVariableRef()
			if err == nil {
				t.Errorf("parseVariableRef(%q) expected error, got nil", input)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 2: types — data type parsing
// ---------------------------------------------------------------------------

// TestParseDataTypes tests parsing of MySQL data types.
func TestParseDataTypes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Integer types
		{"INT", "{DATATYPE :loc 0 :name INT}"},
		{"INTEGER", "{DATATYPE :loc 0 :name INT}"},
		{"INT(11)", "{DATATYPE :loc 0 :name INT :len 11}"},
		{"INT UNSIGNED", "{DATATYPE :loc 0 :name INT :unsigned true}"},
		{"INT(11) UNSIGNED ZEROFILL", "{DATATYPE :loc 0 :name INT :len 11 :unsigned true :zerofill true}"},
		{"TINYINT", "{DATATYPE :loc 0 :name TINYINT}"},
		{"SMALLINT", "{DATATYPE :loc 0 :name SMALLINT}"},
		{"MEDIUMINT", "{DATATYPE :loc 0 :name MEDIUMINT}"},
		{"BIGINT", "{DATATYPE :loc 0 :name BIGINT}"},
		{"BIGINT(20) UNSIGNED", "{DATATYPE :loc 0 :name BIGINT :len 20 :unsigned true}"},

		// SIGNED modifier (accepted and ignored — SIGNED is the default)
		{"INT SIGNED", "{DATATYPE :loc 0 :name INT}"},
		{"BIGINT SIGNED", "{DATATYPE :loc 0 :name BIGINT}"},
		{"TINYINT SIGNED", "{DATATYPE :loc 0 :name TINYINT}"},
		{"SMALLINT SIGNED", "{DATATYPE :loc 0 :name SMALLINT}"},
		{"MEDIUMINT SIGNED", "{DATATYPE :loc 0 :name MEDIUMINT}"},
		{"DECIMAL SIGNED", "{DATATYPE :loc 0 :name DECIMAL}"},

		// Float types
		{"FLOAT", "{DATATYPE :loc 0 :name FLOAT}"},
		{"FLOAT(10)", "{DATATYPE :loc 0 :name FLOAT :len 10}"},
		{"FLOAT(10,2)", "{DATATYPE :loc 0 :name FLOAT :len 10 :scale 2}"},
		{"DOUBLE", "{DATATYPE :loc 0 :name DOUBLE}"},
		{"DOUBLE(10,2)", "{DATATYPE :loc 0 :name DOUBLE :len 10 :scale 2}"},
		{"DOUBLE UNSIGNED", "{DATATYPE :loc 0 :name DOUBLE :unsigned true}"},

		// Decimal types
		{"DECIMAL", "{DATATYPE :loc 0 :name DECIMAL}"},
		{"DECIMAL(10)", "{DATATYPE :loc 0 :name DECIMAL :len 10}"},
		{"DECIMAL(10,2)", "{DATATYPE :loc 0 :name DECIMAL :len 10 :scale 2}"},
		{"NUMERIC(10,2)", "{DATATYPE :loc 0 :name NUMERIC :len 10 :scale 2}"},

		// Boolean
		{"BOOL", "{DATATYPE :loc 0 :name BOOLEAN}"},
		{"BOOLEAN", "{DATATYPE :loc 0 :name BOOLEAN}"},

		// String types
		{"CHAR", "{DATATYPE :loc 0 :name CHAR}"},
		{"CHAR(10)", "{DATATYPE :loc 0 :name CHAR :len 10}"},
		{"CHAR(10) CHARACTER SET utf8", "{DATATYPE :loc 0 :name CHAR :len 10 :charset utf8}"},
		{"VARCHAR(255)", "{DATATYPE :loc 0 :name VARCHAR :len 255}"},
		{"VARCHAR(255) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci", "{DATATYPE :loc 0 :name VARCHAR :len 255 :charset utf8mb4 :collate utf8mb4_unicode_ci}"},
		{"TEXT", "{DATATYPE :loc 0 :name TEXT}"},
		{"TINYTEXT", "{DATATYPE :loc 0 :name TINYTEXT}"},
		{"MEDIUMTEXT", "{DATATYPE :loc 0 :name MEDIUMTEXT}"},
		{"LONGTEXT", "{DATATYPE :loc 0 :name LONGTEXT}"},

		// Date/time types
		{"DATE", "{DATATYPE :loc 0 :name DATE}"},
		{"TIME", "{DATATYPE :loc 0 :name TIME}"},
		{"TIME(3)", "{DATATYPE :loc 0 :name TIME :len 3}"},
		{"DATETIME", "{DATATYPE :loc 0 :name DATETIME}"},
		{"DATETIME(6)", "{DATATYPE :loc 0 :name DATETIME :len 6}"},
		{"TIMESTAMP", "{DATATYPE :loc 0 :name TIMESTAMP}"},
		{"TIMESTAMP(6)", "{DATATYPE :loc 0 :name TIMESTAMP :len 6}"},
		{"YEAR", "{DATATYPE :loc 0 :name YEAR}"},

		// Blob types
		{"BLOB", "{DATATYPE :loc 0 :name BLOB}"},
		{"BLOB(1000)", "{DATATYPE :loc 0 :name BLOB :len 1000}"},
		{"TINYBLOB", "{DATATYPE :loc 0 :name TINYBLOB}"},
		{"MEDIUMBLOB", "{DATATYPE :loc 0 :name MEDIUMBLOB}"},
		{"LONGBLOB", "{DATATYPE :loc 0 :name LONGBLOB}"},

		// Binary types
		{"BINARY", "{DATATYPE :loc 0 :name BINARY}"},
		{"BINARY(16)", "{DATATYPE :loc 0 :name BINARY :len 16}"},
		{"VARBINARY(255)", "{DATATYPE :loc 0 :name VARBINARY :len 255}"},

		// Bit type
		{"BIT", "{DATATYPE :loc 0 :name BIT}"},
		{"BIT(8)", "{DATATYPE :loc 0 :name BIT :len 8}"},

		// ENUM and SET
		{"ENUM('a', 'b', 'c')", "{DATATYPE :loc 0 :name ENUM :enum_values a, b, c}"},
		{"ENUM('x') CHARACTER SET utf8", "{DATATYPE :loc 0 :name ENUM :charset utf8 :enum_values x}"},

		// JSON
		{"JSON", "{DATATYPE :loc 0 :name JSON}"},

		// Type synonyms — Numeric (section 1.1)
		// REAL → DOUBLE
		{"REAL", "{DATATYPE :loc 0 :name DOUBLE}"},
		{"REAL(10,2)", "{DATATYPE :loc 0 :name DOUBLE :len 10 :scale 2}"},
		{"REAL UNSIGNED", "{DATATYPE :loc 0 :name DOUBLE :unsigned true}"},
		// DEC → DECIMAL
		{"DEC", "{DATATYPE :loc 0 :name DECIMAL}"},
		{"DEC(10,2)", "{DATATYPE :loc 0 :name DECIMAL :len 10 :scale 2}"},
		{"DEC UNSIGNED", "{DATATYPE :loc 0 :name DECIMAL :unsigned true}"},
		// FIXED → DECIMAL
		{"FIXED", "{DATATYPE :loc 0 :name DECIMAL}"},
		{"FIXED(10,2)", "{DATATYPE :loc 0 :name DECIMAL :len 10 :scale 2}"},
		{"FIXED UNSIGNED ZEROFILL", "{DATATYPE :loc 0 :name DECIMAL :unsigned true :zerofill true}"},
		// INT1..INT8, MIDDLEINT, FLOAT4, FLOAT8
		{"INT1", "{DATATYPE :loc 0 :name TINYINT}"},
		{"INT2", "{DATATYPE :loc 0 :name SMALLINT}"},
		{"INT3", "{DATATYPE :loc 0 :name MEDIUMINT}"},
		{"INT4", "{DATATYPE :loc 0 :name INT}"},
		{"INT8", "{DATATYPE :loc 0 :name BIGINT}"},
		{"MIDDLEINT", "{DATATYPE :loc 0 :name MEDIUMINT}"},
		{"FLOAT4", "{DATATYPE :loc 0 :name FLOAT}"},
		{"FLOAT8", "{DATATYPE :loc 0 :name DOUBLE}"},
		// DOUBLE PRECISION (regression check)
		{"DOUBLE PRECISION", "{DATATYPE :loc 0 :name DOUBLE}"},

		// Type synonyms — String & Binary (section 1.2)
		// LONG → MEDIUMTEXT
		{"LONG", "{DATATYPE :loc 0 :name MEDIUMTEXT}"},
		// LONG VARCHAR → MEDIUMTEXT
		{"LONG VARCHAR", "{DATATYPE :loc 0 :name MEDIUMTEXT}"},
		// LONG VARBINARY → MEDIUMBLOB
		{"LONG VARBINARY", "{DATATYPE :loc 0 :name MEDIUMBLOB}"},
		// LONG with CHARACTER SET
		{"LONG CHARACTER SET utf8mb4", "{DATATYPE :loc 0 :name MEDIUMTEXT :charset utf8mb4}"},
		// LONG VARCHAR with COLLATE
		{"LONG VARCHAR COLLATE utf8mb4_general_ci", "{DATATYPE :loc 0 :name MEDIUMTEXT :collate utf8mb4_general_ci}"},
		// TEXT with length (regression check)
		{"TEXT(1000)", "{DATATYPE :loc 0 :name TEXT :len 1000}"},
		// NATIONAL CHAR (regression check)
		{"NATIONAL CHAR(10)", "{DATATYPE :loc 0 :name CHAR :len 10 :charset utf8mb3}"},
		// NCHAR (regression check)
		{"NCHAR(10)", "{DATATYPE :loc 0 :name CHAR :len 10 :charset utf8mb3}"},
		// NVARCHAR (regression check)
		{"NVARCHAR(100)", "{DATATYPE :loc 0 :name VARCHAR :len 100 :charset utf8mb3}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.input)}
			p.advance()
			dt, err := p.parseDataType()
			if err != nil {
				t.Fatalf("parseDataType(%q) error: %v", tt.input, err)
			}
			got := ast.NodeToString(dt)
			if got != tt.want {
				t.Errorf("parseDataType(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseDataTypeError tests that invalid data type inputs produce errors.
func TestParseDataTypeError(t *testing.T) {
	inputs := []string{
		"123",
		"'string'",
		"(",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(input)}
			p.advance()
			_, err := p.parseDataType()
			if err == nil {
				t.Errorf("parseDataType(%q) expected error, got nil", input)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 3: expressions
// ---------------------------------------------------------------------------

// helper to parse an expression from a string
func parseExpr(t *testing.T, input string) ast.ExprNode {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	expr, err := p.parseExpr()
	if err != nil {
		t.Fatalf("parseExpr(%q) error: %v", input, err)
	}
	return expr
}

// TestParseArithmetic tests arithmetic expression parsing.
func TestParseArithmetic(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 + 2", "{BINEXPR :op + :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 - 2", "{BINEXPR :op - :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 * 2", "{BINEXPR :op * :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 / 2", "{BINEXPR :op / :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"10 DIV 3", "{BINEXPR :op DIV :loc 3 :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 7}}"},
		{"10 % 3", "{BINEXPR :op %% :loc 3 :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 5}}"},
		{"10 MOD 3", "{BINEXPR :op %% :loc 3 :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 7}}"},
		{"-1", "{UNARY :op - :loc 0 :operand {INT_LIT :val 1 :loc 1}}"},
		{"~5", "{UNARY :op ~ :loc 0 :operand {INT_LIT :val 5 :loc 1}}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			got := ast.NodeToString(expr)
			if got != tt.want {
				t.Errorf("parseExpr(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseComparison tests comparison expression parsing.
func TestParseComparison(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 = 2", "{BINEXPR :op = :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 <> 2", "{BINEXPR :op <> :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 != 2", "{BINEXPR :op <> :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 < 2", "{BINEXPR :op < :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 > 2", "{BINEXPR :op > :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 <= 2", "{BINEXPR :op <= :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 >= 2", "{BINEXPR :op >= :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 <=> 2", "{BINEXPR :op <=> :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 6}}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			got := ast.NodeToString(expr)
			if got != tt.want {
				t.Errorf("parseExpr(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseLogical tests logical expression parsing.
func TestParseLogical(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 AND 2", "{BINEXPR :op AND :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 6}}"},
		{"1 OR 2", "{BINEXPR :op OR :loc 2 :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"NOT 1", "{UNARY :op NOT :loc 0 :operand {INT_LIT :val 1 :loc 4}}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			got := ast.NodeToString(expr)
			if got != tt.want {
				t.Errorf("parseExpr(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseCaseExpr tests CASE expression parsing.
func TestParseCaseExpr(t *testing.T) {
	// Use a simpler comparison approach: just verify parsing succeeds
	// and spot-check key properties
	cases := []string{
		"CASE WHEN 1 THEN 'one' ELSE 'other' END",
		"CASE x WHEN 1 THEN 'one' END",
		"CASE WHEN 1 THEN 'a' WHEN 2 THEN 'b' ELSE 'c' END",
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			expr := parseExpr(t, input)
			ce, ok := expr.(*ast.CaseExpr)
			if !ok {
				t.Fatalf("parseExpr(%q) expected *ast.CaseExpr, got %T", input, expr)
			}
			if len(ce.Whens) == 0 {
				t.Errorf("parseExpr(%q) has no WHEN clauses", input)
			}
		})
	}
}

// TestParseBetween tests BETWEEN expression parsing.
func TestParseBetween(t *testing.T) {
	cases := []struct {
		input string
		not   bool
	}{
		{"x BETWEEN 1 AND 10", false},
		{"x NOT BETWEEN 1 AND 10", true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			be, ok := expr.(*ast.BetweenExpr)
			if !ok {
				t.Fatalf("expected *ast.BetweenExpr, got %T", expr)
			}
			if be.Not != tt.not {
				t.Errorf("Not = %v, want %v", be.Not, tt.not)
			}
		})
	}
}

// TestParseIn tests IN expression parsing.
func TestParseIn(t *testing.T) {
	cases := []struct {
		input string
		not   bool
		count int
	}{
		{"x IN (1, 2, 3)", false, 3},
		{"x NOT IN (1)", true, 1},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			ie, ok := expr.(*ast.InExpr)
			if !ok {
				t.Fatalf("expected *ast.InExpr, got %T", expr)
			}
			if ie.Not != tt.not {
				t.Errorf("Not = %v, want %v", ie.Not, tt.not)
			}
			if len(ie.List) != tt.count {
				t.Errorf("list len = %d, want %d", len(ie.List), tt.count)
			}
		})
	}
}

// TestParseLike tests LIKE expression parsing.
func TestParseLike(t *testing.T) {
	cases := []struct {
		input string
		not   bool
	}{
		{"x LIKE '%test%'", false},
		{"x NOT LIKE '%test%'", true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			le, ok := expr.(*ast.LikeExpr)
			if !ok {
				t.Fatalf("expected *ast.LikeExpr, got %T", expr)
			}
			if le.Not != tt.not {
				t.Errorf("Not = %v, want %v", le.Not, tt.not)
			}
		})
	}
}

// TestParseIs tests IS expression parsing.
func TestParseIs(t *testing.T) {
	cases := []struct {
		input string
		not   bool
		test  ast.IsTestType
	}{
		{"x IS NULL", false, ast.IsNull},
		{"x IS NOT NULL", true, ast.IsNull},
		{"x IS TRUE", false, ast.IsTrue},
		{"x IS FALSE", false, ast.IsFalse},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			ie, ok := expr.(*ast.IsExpr)
			if !ok {
				t.Fatalf("expected *ast.IsExpr, got %T", expr)
			}
			if ie.Not != tt.not {
				t.Errorf("Not = %v, want %v", ie.Not, tt.not)
			}
			if ie.Test != tt.test {
				t.Errorf("Test = %v, want %v", ie.Test, tt.test)
			}
		})
	}
}

// TestParseCast tests CAST expression parsing.
func TestParseCast(t *testing.T) {
	cases := []struct {
		input    string
		typeName string
	}{
		{"CAST(x AS INT)", "INT"},
		{"CAST('123' AS DECIMAL(10,2))", "DECIMAL"},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			ce, ok := expr.(*ast.CastExpr)
			if !ok {
				t.Fatalf("expected *ast.CastExpr, got %T", expr)
			}
			if ce.TypeName.Name != tt.typeName {
				t.Errorf("TypeName = %s, want %s", ce.TypeName.Name, tt.typeName)
			}
		})
	}
}

// TestParseExtract tests EXTRACT(unit FROM expr) parsing.
func TestParseExtract(t *testing.T) {
	cases := []struct {
		sql  string
		unit string
	}{
		{"SELECT EXTRACT(HOUR FROM NOW())", "HOUR"},
		{"SELECT EXTRACT(DAY FROM '2024-01-01')", "DAY"},
		{"SELECT EXTRACT(YEAR FROM created_at) FROM t", "YEAR"},
		{"SELECT EXTRACT(MINUTE FROM ts)", "MINUTE"},
		{"SELECT EXTRACT(SECOND FROM ts)", "SECOND"},
		{"SELECT EXTRACT(MONTH FROM ts)", "MONTH"},
		// Keyword-token units: YEAR is kwYEAR (registered keyword), must be
		// accepted via parseKeywordOrIdent.
		{"SELECT EXTRACT(YEAR FROM ts)", "YEAR"},
		// MICROSECOND and QUARTER are plain identifiers (tokIDENT).
		{"SELECT EXTRACT(MICROSECOND FROM ts)", "MICROSECOND"},
		{"SELECT EXTRACT(QUARTER FROM ts)", "QUARTER"},
		{"SELECT EXTRACT(WEEK FROM ts)", "WEEK"},
	}
	for _, tt := range cases {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if result.Len() == 0 {
				t.Fatalf("Parse(%q) returned empty list", tt.sql)
			}
			// Extract the first target expression from the SELECT.
			sel, ok := result.Items[0].(*ast.SelectStmt)
			if !ok {
				t.Fatalf("expected *ast.SelectStmt, got %T", result.Items[0])
			}
			if len(sel.TargetList) == 0 {
				t.Fatalf("no targets in SELECT")
			}
			// TargetList elements may be ResTarget wrapping the actual expr.
			var exprNode ast.ExprNode = sel.TargetList[0]
			if rt, ok := exprNode.(*ast.ResTarget); ok {
				exprNode = rt.Val
			}
			ext, ok := exprNode.(*ast.ExtractExpr)
			if !ok {
				t.Fatalf("expected *ast.ExtractExpr, got %T", exprNode)
			}
			if ext.Unit != tt.unit {
				t.Errorf("Unit = %s, want %s", ext.Unit, tt.unit)
			}
		})
	}
}

// TestParseFuncCall tests function call parsing.
func TestParseFuncCall(t *testing.T) {
	cases := []struct {
		input    string
		name     string
		star     bool
		distinct bool
		argCount int
	}{
		{"NOW()", "NOW", false, false, 0},
		{"COUNT(*)", "COUNT", true, false, 0},
		{"COUNT(DISTINCT x)", "COUNT", false, true, 1},
		{"COALESCE(a, b, c)", "COALESCE", false, false, 3},
		{"CONCAT('a', 'b')", "CONCAT", false, false, 2},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			fc, ok := expr.(*ast.FuncCallExpr)
			if !ok {
				t.Fatalf("expected *ast.FuncCallExpr, got %T", expr)
			}
			if fc.Name != tt.name {
				t.Errorf("Name = %s, want %s", fc.Name, tt.name)
			}
			if fc.Star != tt.star {
				t.Errorf("Star = %v, want %v", fc.Star, tt.star)
			}
			if fc.Distinct != tt.distinct {
				t.Errorf("Distinct = %v, want %v", fc.Distinct, tt.distinct)
			}
			if len(fc.Args) != tt.argCount {
				t.Errorf("arg count = %d, want %d", len(fc.Args), tt.argCount)
			}
		})
	}
}

// TestParseLiterals tests literal parsing.
func TestParseLiterals(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"42", "{INT_LIT :val 42 :loc 0}"},
		{"3.14", "{FLOAT_LIT :val 3.14 :loc 0}"},
		{"'hello'", "{STRING_LIT :val \"hello\" :loc 0}"},
		{"TRUE", "{BOOL_LIT :val true :loc 0}"},
		{"FALSE", "{BOOL_LIT :val false :loc 0}"},
		{"NULL", "{NULL_LIT :loc 0}"},
		{"DEFAULT", "{DEFAULT :loc 0}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			got := ast.NodeToString(expr)
			if got != tt.want {
				t.Errorf("parseExpr(%q):\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseParenExpr tests parenthesized expression parsing.
func TestParseParenExpr(t *testing.T) {
	expr := parseExpr(t, "(1 + 2)")
	pe, ok := expr.(*ast.ParenExpr)
	if !ok {
		t.Fatalf("expected *ast.ParenExpr, got %T", expr)
	}
	be, ok := pe.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr inside paren, got %T", pe.Expr)
	}
	if be.Op != ast.BinOpAdd {
		t.Errorf("Op = %v, want BinOpAdd", be.Op)
	}
}

// TestParsePrecedence tests operator precedence.
func TestParsePrecedence(t *testing.T) {
	// 1 + 2 * 3 should parse as 1 + (2 * 3)
	expr := parseExpr(t, "1 + 2 * 3")
	be, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", expr)
	}
	if be.Op != ast.BinOpAdd {
		t.Errorf("outer Op = %v, want BinOpAdd", be.Op)
	}
	inner, ok := be.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("right side: expected *ast.BinaryExpr, got %T", be.Right)
	}
	if inner.Op != ast.BinOpMul {
		t.Errorf("inner Op = %v, want BinOpMul", inner.Op)
	}
}

// TestParseConvert tests CONVERT expression parsing.
func TestParseConvert(t *testing.T) {
	// CONVERT with USING
	expr1 := parseExpr(t, "CONVERT(x USING utf8)")
	cv1, ok := expr1.(*ast.ConvertExpr)
	if !ok {
		t.Fatalf("expected *ast.ConvertExpr, got %T", expr1)
	}
	if cv1.Charset != "utf8" {
		t.Errorf("Charset = %s, want utf8", cv1.Charset)
	}

	// CONVERT with type
	expr2 := parseExpr(t, "CONVERT('123', DECIMAL(10,2))")
	cv2, ok := expr2.(*ast.ConvertExpr)
	if !ok {
		t.Fatalf("expected *ast.ConvertExpr, got %T", expr2)
	}
	if cv2.TypeName == nil || cv2.TypeName.Name != "DECIMAL" {
		t.Errorf("TypeName.Name = %v, want DECIMAL", cv2.TypeName)
	}
}

// ---------------------------------------------------------------------------
// Batch 4: select
// ---------------------------------------------------------------------------

// helper to parse a SELECT statement directly
func parseSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	sel, err := p.parseSelectStmt()
	if err != nil {
		t.Fatalf("parseSelectStmt(%q) error: %v", input, err)
	}
	return sel
}

// TestParseSimpleSelect tests simple SELECT statements.
func TestParseSimpleSelect(t *testing.T) {
	cases := []struct {
		input     string
		targetLen int
		hasFrom   bool
		hasWhere  bool
	}{
		{"SELECT 1", 1, false, false},
		{"SELECT 1, 2, 3", 3, false, false},
		{"SELECT *", 1, false, false},
		{"SELECT a, b FROM t", 2, true, false},
		{"SELECT a FROM t WHERE a > 1", 1, true, true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			sel := parseSelect(t, tt.input)
			if len(sel.TargetList) != tt.targetLen {
				t.Errorf("target count = %d, want %d", len(sel.TargetList), tt.targetLen)
			}
			if (len(sel.From) > 0) != tt.hasFrom {
				t.Errorf("hasFrom = %v, want %v", len(sel.From) > 0, tt.hasFrom)
			}
			if (sel.Where != nil) != tt.hasWhere {
				t.Errorf("hasWhere = %v, want %v", sel.Where != nil, tt.hasWhere)
			}
		})
	}
}

// TestParseSelectFrom tests SELECT with FROM clause.
func TestParseSelectFrom(t *testing.T) {
	sel := parseSelect(t, "SELECT a, b FROM users AS u")
	if len(sel.From) != 1 {
		t.Fatalf("from count = %d, want 1", len(sel.From))
	}
	tr, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("expected *ast.TableRef, got %T", sel.From[0])
	}
	if tr.Name != "users" {
		t.Errorf("table name = %s, want users", tr.Name)
	}
	if tr.Alias != "u" {
		t.Errorf("alias = %s, want u", tr.Alias)
	}
}

// TestParseSelectJoin tests SELECT with JOIN clauses.
func TestParseSelectJoin(t *testing.T) {
	cases := []struct {
		input    string
		joinType ast.JoinType
	}{
		{"SELECT * FROM a JOIN b ON a.id = b.id", ast.JoinInner},
		{"SELECT * FROM a LEFT JOIN b ON a.id = b.id", ast.JoinLeft},
		{"SELECT * FROM a RIGHT JOIN b ON a.id = b.id", ast.JoinRight},
		{"SELECT * FROM a CROSS JOIN b", ast.JoinCross},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			sel := parseSelect(t, tt.input)
			if len(sel.From) != 1 {
				t.Fatalf("from count = %d, want 1", len(sel.From))
			}
			jc, ok := sel.From[0].(*ast.JoinClause)
			if !ok {
				t.Fatalf("expected *ast.JoinClause, got %T", sel.From[0])
			}
			if jc.Type != tt.joinType {
				t.Errorf("join type = %d, want %d", jc.Type, tt.joinType)
			}
		})
	}
}

// TestParseSelectGroupBy tests SELECT with GROUP BY clause.
func TestParseSelectGroupBy(t *testing.T) {
	sel := parseSelect(t, "SELECT dept, COUNT(*) FROM emp GROUP BY dept HAVING COUNT(*) > 1")
	if len(sel.GroupBy) != 1 {
		t.Errorf("group_by count = %d, want 1", len(sel.GroupBy))
	}
	if sel.Having == nil {
		t.Error("expected HAVING clause")
	}
}

// TestParseSelectOrderBy tests SELECT with ORDER BY clause.
func TestParseSelectOrderBy(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t ORDER BY a ASC, b DESC")
	if len(sel.OrderBy) != 2 {
		t.Fatalf("order_by count = %d, want 2", len(sel.OrderBy))
	}
	if sel.OrderBy[0].Desc {
		t.Error("first ORDER BY item should be ASC")
	}
	if !sel.OrderBy[1].Desc {
		t.Error("second ORDER BY item should be DESC")
	}
}

// TestParseSelectLimit tests SELECT with LIMIT clause.
func TestParseSelectLimit(t *testing.T) {
	cases := []struct {
		input     string
		hasOffset bool
	}{
		{"SELECT * FROM t LIMIT 10", false},
		{"SELECT * FROM t LIMIT 10 OFFSET 5", true},
		{"SELECT * FROM t LIMIT 5, 10", true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			sel := parseSelect(t, tt.input)
			if sel.Limit == nil {
				t.Fatal("expected LIMIT clause")
			}
			if (sel.Limit.Offset != nil) != tt.hasOffset {
				t.Errorf("hasOffset = %v, want %v", sel.Limit.Offset != nil, tt.hasOffset)
			}
		})
	}
}

// TestParseSelectUnion tests SELECT with UNION.
func TestParseSelectUnion(t *testing.T) {
	sel := parseSelect(t, "SELECT 1 UNION SELECT 2")
	if sel.SetOp != ast.SetOpUnion {
		t.Errorf("SetOp = %d, want SetOpUnion", sel.SetOp)
	}
	if sel.Left == nil || sel.Right == nil {
		t.Error("expected Left and Right in UNION")
	}
}

// TestParseSelectForUpdate tests SELECT with FOR UPDATE.
func TestParseSelectForUpdate(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t FOR UPDATE")
	if sel.ForUpdate == nil {
		t.Fatal("expected FOR UPDATE clause")
	}
	if sel.ForUpdate.Share {
		t.Error("expected FOR UPDATE, not FOR SHARE")
	}

	sel2 := parseSelect(t, "SELECT * FROM t FOR SHARE")
	if sel2.ForUpdate == nil {
		t.Fatal("expected FOR SHARE clause")
	}
	if !sel2.ForUpdate.Share {
		t.Error("expected FOR SHARE")
	}
}

// TestParseSelectDistinct tests SELECT DISTINCT.
func TestParseSelectDistinct(t *testing.T) {
	sel := parseSelect(t, "SELECT DISTINCT a, b FROM t")
	if sel.DistinctKind != ast.DistinctOn {
		t.Errorf("DistinctKind = %d, want DistinctOn", sel.DistinctKind)
	}
}

// TestParseSelectAlias tests SELECT with column aliases.
func TestParseSelectAlias(t *testing.T) {
	sel := parseSelect(t, "SELECT a AS col1, b col2 FROM t")
	if len(sel.TargetList) != 2 {
		t.Fatalf("target count = %d, want 2", len(sel.TargetList))
	}
	rt1, ok := sel.TargetList[0].(*ast.ResTarget)
	if !ok {
		t.Fatalf("expected *ast.ResTarget, got %T", sel.TargetList[0])
	}
	if rt1.Name != "col1" {
		t.Errorf("alias = %s, want col1", rt1.Name)
	}
	rt2, ok := sel.TargetList[1].(*ast.ResTarget)
	if !ok {
		t.Fatalf("expected *ast.ResTarget, got %T", sel.TargetList[1])
	}
	if rt2.Name != "col2" {
		t.Errorf("alias = %s, want col2", rt2.Name)
	}
}

// TestParseSelectSubquery tests SELECT with subquery.
func TestParseSelectSubquery(t *testing.T) {
	sel := parseSelect(t, "SELECT * FROM t WHERE id IN (SELECT id FROM t2)")
	if sel.Where == nil {
		t.Fatal("expected WHERE clause")
	}
	inExpr, ok := sel.Where.(*ast.InExpr)
	if !ok {
		t.Fatalf("expected *ast.InExpr, got %T", sel.Where)
	}
	if inExpr.Select == nil {
		t.Error("expected subquery in IN")
	}
}

// TestParseEmpty tests parsing empty input.
func TestParseEmpty(t *testing.T) {
	result, err := Parse("")
	if err != nil {
		t.Fatalf("Parse(\"\") error: %v", err)
	}
	if result.Len() != 0 {
		t.Errorf("Parse(\"\") returned %d items, want 0", result.Len())
	}
}

// TestParseSemicolons tests parsing semicolons only.
func TestParseSemicolons(t *testing.T) {
	result, err := Parse(";;;")
	if err != nil {
		t.Fatalf("Parse(\";;;\") error: %v", err)
	}
	if result.Len() != 0 {
		t.Errorf("Parse(\";;;\") returned %d items, want 0", result.Len())
	}
}

// ---------------------------------------------------------------------------
// Batch 5: insert
// ---------------------------------------------------------------------------

// helper to parse an INSERT statement directly
func parseInsert(t *testing.T, input string) *ast.InsertStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseInsertStmt()
	if err != nil {
		t.Fatalf("parseInsertStmt(%q) error: %v", input, err)
	}
	return stmt
}

// helper to parse a REPLACE statement directly
func parseReplace(t *testing.T, input string) *ast.InsertStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseReplaceStmt()
	if err != nil {
		t.Fatalf("parseReplaceStmt(%q) error: %v", input, err)
	}
	return stmt
}

// TestParseInsertValues tests INSERT ... VALUES parsing.
func TestParseInsertValues(t *testing.T) {
	cases := []struct {
		input    string
		table    string
		colCount int
		rowCount int
		ignore   bool
		priority ast.InsertPriority
	}{
		{
			input:    "INSERT INTO users (id, name) VALUES (1, 'alice')",
			table:    "users",
			colCount: 2,
			rowCount: 1,
		},
		{
			input:    "INSERT INTO users VALUES (1, 'alice'), (2, 'bob')",
			table:    "users",
			colCount: 0,
			rowCount: 2,
		},
		{
			input:    "INSERT users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
		},
		{
			input:    "INSERT LOW_PRIORITY INTO users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
			priority: ast.InsertPriorityLow,
		},
		{
			input:    "INSERT HIGH_PRIORITY IGNORE INTO users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
			ignore:   true,
			priority: ast.InsertPriorityHigh,
		},
		{
			input:    "INSERT DELAYED INTO users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
			priority: ast.InsertPriorityDelayed,
		},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			stmt := parseInsert(t, tt.input)
			if stmt.IsReplace {
				t.Error("IsReplace should be false for INSERT")
			}
			if stmt.Table.Name != tt.table {
				t.Errorf("Table.Name = %s, want %s", stmt.Table.Name, tt.table)
			}
			if len(stmt.Columns) != tt.colCount {
				t.Errorf("column count = %d, want %d", len(stmt.Columns), tt.colCount)
			}
			if len(stmt.Values) != tt.rowCount {
				t.Errorf("row count = %d, want %d", len(stmt.Values), tt.rowCount)
			}
			if stmt.Ignore != tt.ignore {
				t.Errorf("Ignore = %v, want %v", stmt.Ignore, tt.ignore)
			}
			if stmt.Priority != tt.priority {
				t.Errorf("Priority = %d, want %d", stmt.Priority, tt.priority)
			}
			if stmt.Select != nil {
				t.Error("Select should be nil for VALUES insert")
			}
			if len(stmt.SetList) != 0 {
				t.Error("SetList should be empty for VALUES insert")
			}
			if stmt.Loc.Start < 0 {
				t.Error("Loc.Start should be >= 0")
			}
		})
	}
}

// TestParseInsertValuesRow tests INSERT ... VALUES ROW(...) row constructor syntax.
func TestParseInsertValuesRow(t *testing.T) {
	ParseAndCheck(t, "INSERT INTO t VALUES ROW(1, 2, 3)")
}

// TestParseInsertValuesRowMultiple tests multiple ROW constructors in INSERT.
func TestParseInsertValuesRowMultiple(t *testing.T) {
	ParseAndCheck(t, "INSERT INTO t VALUES ROW(1, 'a'), ROW(2, 'b'), ROW(3, 'c')")
}

// TestParseInsertSelect tests INSERT ... SELECT parsing.
func TestParseInsertSelect(t *testing.T) {
	cases := []struct {
		input    string
		table    string
		colCount int
	}{
		{
			input:    "INSERT INTO users (id, name) SELECT id, name FROM tmp_users",
			table:    "users",
			colCount: 2,
		},
		{
			input:    "INSERT INTO users SELECT * FROM tmp_users",
			table:    "users",
			colCount: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			stmt := parseInsert(t, tt.input)
			if stmt.Table.Name != tt.table {
				t.Errorf("Table.Name = %s, want %s", stmt.Table.Name, tt.table)
			}
			if len(stmt.Columns) != tt.colCount {
				t.Errorf("column count = %d, want %d", len(stmt.Columns), tt.colCount)
			}
			if stmt.Select == nil {
				t.Fatal("Select should not be nil for INSERT ... SELECT")
			}
			if len(stmt.Values) != 0 {
				t.Error("Values should be empty for INSERT ... SELECT")
			}
			if len(stmt.Select.TargetList) == 0 {
				t.Error("Select.TargetList should not be empty")
			}
		})
	}
}

// TestParseInsertSet tests INSERT ... SET parsing.
func TestParseInsertSet(t *testing.T) {
	cases := []struct {
		input    string
		table    string
		setCount int
	}{
		{
			input:    "INSERT INTO users SET id = 1, name = 'alice'",
			table:    "users",
			setCount: 2,
		},
		{
			input:    "INSERT INTO users SET id = 1",
			table:    "users",
			setCount: 1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			stmt := parseInsert(t, tt.input)
			if stmt.Table.Name != tt.table {
				t.Errorf("Table.Name = %s, want %s", stmt.Table.Name, tt.table)
			}
			if len(stmt.SetList) != tt.setCount {
				t.Errorf("SetList count = %d, want %d", len(stmt.SetList), tt.setCount)
			}
			if len(stmt.Values) != 0 {
				t.Error("Values should be empty for SET insert")
			}
			if stmt.Select != nil {
				t.Error("Select should be nil for SET insert")
			}
			// Verify assignment structure
			for i, a := range stmt.SetList {
				if a.Column == nil {
					t.Errorf("SetList[%d].Column is nil", i)
				}
				if a.Value == nil {
					t.Errorf("SetList[%d].Value is nil", i)
				}
			}
		})
	}
}

// TestParseReplace tests REPLACE statement parsing.
func TestParseReplace(t *testing.T) {
	cases := []struct {
		input    string
		table    string
		colCount int
		rowCount int
		priority ast.InsertPriority
	}{
		{
			input:    "REPLACE INTO users (id, name) VALUES (1, 'alice')",
			table:    "users",
			colCount: 2,
			rowCount: 1,
		},
		{
			input:    "REPLACE users VALUES (1, 'alice')",
			table:    "users",
			colCount: 0,
			rowCount: 1,
		},
		{
			input:    "REPLACE LOW_PRIORITY INTO users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
			priority: ast.InsertPriorityLow,
		},
		{
			input:    "REPLACE DELAYED INTO users (id) VALUES (1)",
			table:    "users",
			colCount: 1,
			rowCount: 1,
			priority: ast.InsertPriorityDelayed,
		},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			stmt := parseReplace(t, tt.input)
			if !stmt.IsReplace {
				t.Error("IsReplace should be true for REPLACE")
			}
			if stmt.Table.Name != tt.table {
				t.Errorf("Table.Name = %s, want %s", stmt.Table.Name, tt.table)
			}
			if len(stmt.Columns) != tt.colCount {
				t.Errorf("column count = %d, want %d", len(stmt.Columns), tt.colCount)
			}
			if len(stmt.Values) != tt.rowCount {
				t.Errorf("row count = %d, want %d", len(stmt.Values), tt.rowCount)
			}
			if stmt.Priority != tt.priority {
				t.Errorf("Priority = %d, want %d", stmt.Priority, tt.priority)
			}
			// REPLACE should never have ON DUPLICATE KEY UPDATE
			if len(stmt.OnDuplicateKey) != 0 {
				t.Error("OnDuplicateKey should be empty for REPLACE")
			}
		})
	}
}

// TestParseInsertOnDupKey tests INSERT ... ON DUPLICATE KEY UPDATE parsing.
func TestParseInsertOnDupKey(t *testing.T) {
	cases := []struct {
		input    string
		table    string
		rowCount int
		dupCount int
	}{
		{
			input:    "INSERT INTO users (id, name) VALUES (1, 'alice') ON DUPLICATE KEY UPDATE name = 'alice'",
			table:    "users",
			rowCount: 1,
			dupCount: 1,
		},
		{
			input:    "INSERT INTO counters (id, cnt) VALUES (1, 1) ON DUPLICATE KEY UPDATE cnt = cnt + 1, updated_at = NOW()",
			table:    "counters",
			rowCount: 1,
			dupCount: 2,
		},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			stmt := parseInsert(t, tt.input)
			if stmt.Table.Name != tt.table {
				t.Errorf("Table.Name = %s, want %s", stmt.Table.Name, tt.table)
			}
			if len(stmt.Values) != tt.rowCount {
				t.Errorf("row count = %d, want %d", len(stmt.Values), tt.rowCount)
			}
			if len(stmt.OnDuplicateKey) != tt.dupCount {
				t.Fatalf("OnDuplicateKey count = %d, want %d", len(stmt.OnDuplicateKey), tt.dupCount)
			}
			// Verify assignment structure
			for i, a := range stmt.OnDuplicateKey {
				if a.Column == nil {
					t.Errorf("OnDuplicateKey[%d].Column is nil", i)
				}
				if a.Value == nil {
					t.Errorf("OnDuplicateKey[%d].Value is nil", i)
				}
				if a.Loc.Start < 0 {
					t.Errorf("OnDuplicateKey[%d].Loc.Start should be >= 0", i)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 99 (review_dml): INSERT ... TABLE source
// ---------------------------------------------------------------------------

// TestParseInsertTableSource tests INSERT INTO t1 TABLE t2 (MySQL 8.0.19+).
func TestParseInsertTableSource(t *testing.T) {
	stmt := parseInsert(t, "INSERT INTO t1 TABLE t2")
	if stmt.TableSource == nil {
		t.Fatal("expected TableSource to be non-nil")
	}
	if stmt.TableSource.Table.Name != "t2" {
		t.Errorf("TableSource.Table.Name = %q, want %q", stmt.TableSource.Table.Name, "t2")
	}

	// REPLACE ... TABLE
	rep := parseReplace(t, "REPLACE INTO t1 TABLE t2")
	if rep.TableSource == nil {
		t.Fatal("expected TableSource to be non-nil for REPLACE")
	}

	// INSERT with ON DUPLICATE KEY UPDATE after TABLE source
	ParseAndCheck(t, "INSERT INTO t1 TABLE t2 ON DUPLICATE KEY UPDATE col = 1")
}

// ---------------------------------------------------------------------------
// Batch 6: update_delete
// ---------------------------------------------------------------------------

// helper to parse an UPDATE statement directly
func parseUpdate(t *testing.T, input string) *ast.UpdateStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseUpdateStmt()
	if err != nil {
		t.Fatalf("parseUpdateStmt(%q) error: %v", input, err)
	}
	return stmt
}

// helper to parse a DELETE statement directly
func parseDelete(t *testing.T, input string) *ast.DeleteStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseDeleteStmt()
	if err != nil {
		t.Fatalf("parseDeleteStmt(%q) error: %v", input, err)
	}
	return stmt
}

// TestParseUpdate tests single-table UPDATE statements.
func TestParseUpdate(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users SET name = 'alice'")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1", len(stmt.Tables))
		}
		tr, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef, got %T", stmt.Tables[0])
		}
		if tr.Name != "users" {
			t.Errorf("table name = %s, want users", tr.Name)
		}
		if len(stmt.SetList) != 1 {
			t.Fatalf("set list count = %d, want 1", len(stmt.SetList))
		}
		if stmt.SetList[0].Column.Column != "name" {
			t.Errorf("set column = %s, want name", stmt.SetList[0].Column.Column)
		}
		if stmt.Where != nil {
			t.Error("expected no WHERE clause")
		}
		if stmt.Loc.Start < 0 {
			t.Error("Loc.Start should be >= 0")
		}
	})

	t.Run("with WHERE", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users SET age = 30 WHERE id = 1")
		if stmt.Where == nil {
			t.Fatal("expected WHERE clause")
		}
		if len(stmt.SetList) != 1 {
			t.Fatalf("set list count = %d, want 1", len(stmt.SetList))
		}
	})

	t.Run("with ORDER BY and LIMIT", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users SET score = 0 ORDER BY score ASC LIMIT 10")
		if len(stmt.OrderBy) != 1 {
			t.Fatalf("order_by count = %d, want 1", len(stmt.OrderBy))
		}
		if stmt.Limit == nil {
			t.Fatal("expected LIMIT clause")
		}
	})

	t.Run("multiple assignments", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users SET name = 'bob', age = 25, active = TRUE WHERE id = 1")
		if len(stmt.SetList) != 3 {
			t.Fatalf("set list count = %d, want 3", len(stmt.SetList))
		}
		if stmt.SetList[0].Column.Column != "name" {
			t.Errorf("first set column = %s, want name", stmt.SetList[0].Column.Column)
		}
		if stmt.SetList[1].Column.Column != "age" {
			t.Errorf("second set column = %s, want age", stmt.SetList[1].Column.Column)
		}
		if stmt.SetList[2].Column.Column != "active" {
			t.Errorf("third set column = %s, want active", stmt.SetList[2].Column.Column)
		}
	})

	t.Run("LOW_PRIORITY IGNORE", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE LOW_PRIORITY IGNORE users SET name = 'test'")
		if !stmt.LowPriority {
			t.Error("expected LowPriority = true")
		}
		if !stmt.Ignore {
			t.Error("expected Ignore = true")
		}
	})

	t.Run("qualified column", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users SET users.name = 'test'")
		if len(stmt.SetList) != 1 {
			t.Fatalf("set list count = %d, want 1", len(stmt.SetList))
		}
		if stmt.SetList[0].Column.Table != "users" {
			t.Errorf("set column table = %s, want users", stmt.SetList[0].Column.Table)
		}
		if stmt.SetList[0].Column.Column != "name" {
			t.Errorf("set column name = %s, want name", stmt.SetList[0].Column.Column)
		}
	})
}

// TestParseMultiTableUpdate tests multi-table UPDATE statements.
func TestParseMultiTableUpdate(t *testing.T) {
	t.Run("two tables with join", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE users JOIN orders ON users.id = orders.user_id SET users.name = 'test' WHERE orders.total > 100")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1 (join clause)", len(stmt.Tables))
		}
		jc, ok := stmt.Tables[0].(*ast.JoinClause)
		if !ok {
			t.Fatalf("expected *ast.JoinClause, got %T", stmt.Tables[0])
		}
		if jc.Type != ast.JoinInner {
			t.Errorf("join type = %d, want JoinInner", jc.Type)
		}
		if len(stmt.SetList) != 1 {
			t.Fatalf("set list count = %d, want 1", len(stmt.SetList))
		}
		if stmt.Where == nil {
			t.Fatal("expected WHERE clause")
		}
	})

	t.Run("comma-separated tables", func(t *testing.T) {
		stmt := parseUpdate(t, "UPDATE t1, t2 SET t1.a = t2.b WHERE t1.id = t2.id")
		if len(stmt.Tables) != 2 {
			t.Fatalf("tables count = %d, want 2", len(stmt.Tables))
		}
		tr1, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef for tables[0], got %T", stmt.Tables[0])
		}
		if tr1.Name != "t1" {
			t.Errorf("table[0] name = %s, want t1", tr1.Name)
		}
		tr2, ok := stmt.Tables[1].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef for tables[1], got %T", stmt.Tables[1])
		}
		if tr2.Name != "t2" {
			t.Errorf("table[1] name = %s, want t2", tr2.Name)
		}
	})
}

// TestParseDelete tests single-table DELETE statements.
func TestParseDelete(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE FROM users")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1", len(stmt.Tables))
		}
		tr, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef, got %T", stmt.Tables[0])
		}
		if tr.Name != "users" {
			t.Errorf("table name = %s, want users", tr.Name)
		}
		if stmt.Where != nil {
			t.Error("expected no WHERE clause")
		}
		if stmt.Loc.Start < 0 {
			t.Error("Loc.Start should be >= 0")
		}
	})

	t.Run("with WHERE", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE FROM users WHERE id = 1")
		if stmt.Where == nil {
			t.Fatal("expected WHERE clause")
		}
	})

	t.Run("with ORDER BY and LIMIT", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE FROM users ORDER BY created_at ASC LIMIT 100")
		if len(stmt.OrderBy) != 1 {
			t.Fatalf("order_by count = %d, want 1", len(stmt.OrderBy))
		}
		if stmt.Limit == nil {
			t.Fatal("expected LIMIT clause")
		}
	})

	t.Run("modifiers", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE LOW_PRIORITY QUICK IGNORE FROM users WHERE id = 1")
		if !stmt.LowPriority {
			t.Error("expected LowPriority = true")
		}
		if !stmt.Quick {
			t.Error("expected Quick = true")
		}
		if !stmt.Ignore {
			t.Error("expected Ignore = true")
		}
	})
}

// TestParseMultiTableDelete tests multi-table DELETE statements.
func TestParseMultiTableDelete(t *testing.T) {
	t.Run("syntax1: DELETE t1 FROM t1 JOIN t2", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id WHERE t2.status = 0")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1", len(stmt.Tables))
		}
		tr, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef, got %T", stmt.Tables[0])
		}
		if tr.Name != "t1" {
			t.Errorf("table name = %s, want t1", tr.Name)
		}
		if len(stmt.Using) != 1 {
			t.Fatalf("using count = %d, want 1 (join clause)", len(stmt.Using))
		}
		_, ok = stmt.Using[0].(*ast.JoinClause)
		if !ok {
			t.Fatalf("expected *ast.JoinClause in Using, got %T", stmt.Using[0])
		}
		if stmt.Where == nil {
			t.Fatal("expected WHERE clause")
		}
	})

	t.Run("syntax1: DELETE t1, t2 FROM ...", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE t1, t2 FROM t1 JOIN t2 ON t1.id = t2.ref_id")
		if len(stmt.Tables) != 2 {
			t.Fatalf("tables count = %d, want 2", len(stmt.Tables))
		}
		if len(stmt.Using) != 1 {
			t.Fatalf("using count = %d, want 1", len(stmt.Using))
		}
	})

	t.Run("syntax2: DELETE FROM t1 USING t1 JOIN t2", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE FROM t1 USING t1 JOIN t2 ON t1.id = t2.id WHERE t2.active = 0")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1", len(stmt.Tables))
		}
		tr, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef, got %T", stmt.Tables[0])
		}
		if tr.Name != "t1" {
			t.Errorf("table name = %s, want t1", tr.Name)
		}
		if len(stmt.Using) != 1 {
			t.Fatalf("using count = %d, want 1", len(stmt.Using))
		}
		if stmt.Where == nil {
			t.Fatal("expected WHERE clause")
		}
	})

	t.Run("with .* suffix", func(t *testing.T) {
		stmt := parseDelete(t, "DELETE t1.* FROM t1 JOIN t2 ON t1.id = t2.id")
		if len(stmt.Tables) != 1 {
			t.Fatalf("tables count = %d, want 1", len(stmt.Tables))
		}
		tr, ok := stmt.Tables[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("expected *ast.TableRef, got %T", stmt.Tables[0])
		}
		if tr.Name != "t1" {
			t.Errorf("table name = %s, want t1", tr.Name)
		}
	})
}

// ============================================================================
// Batch 11: CREATE/ALTER/DROP DATABASE
// ============================================================================

func parseCreateDatabase(t *testing.T, input string) *ast.CreateDatabaseStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume CREATE
	stmt, err := p.parseCreateDatabaseStmt()
	if err != nil {
		t.Fatalf("parseCreateDatabaseStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseCreateDatabase(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseCreateDatabase(t, "CREATE DATABASE mydb")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if stmt.IfNotExists {
			t.Error("IfNotExists should be false")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := parseCreateDatabase(t, "CREATE DATABASE IF NOT EXISTS mydb")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if !stmt.IfNotExists {
			t.Error("IfNotExists should be true")
		}
	})

	t.Run("with charset", func(t *testing.T) {
		stmt := parseCreateDatabase(t, "CREATE DATABASE mydb CHARACTER SET utf8mb4")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "CHARACTER SET" {
			t.Errorf("option name = %s, want CHARACTER SET", stmt.Options[0].Name)
		}
		if stmt.Options[0].Value != "utf8mb4" {
			t.Errorf("option value = %s, want utf8mb4", stmt.Options[0].Value)
		}
	})

	t.Run("schema keyword", func(t *testing.T) {
		stmt := parseCreateDatabase(t, "CREATE SCHEMA mydb")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
	})
}

func parseAlterDatabase(t *testing.T, input string) *ast.AlterDatabaseStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume ALTER
	stmt, err := p.parseAlterDatabaseStmt()
	if err != nil {
		t.Fatalf("parseAlterDatabaseStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseAlterDatabase(t *testing.T) {
	t.Run("charset", func(t *testing.T) {
		stmt := parseAlterDatabase(t, "ALTER DATABASE mydb CHARACTER SET utf8mb4")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
	})

	t.Run("collate", func(t *testing.T) {
		stmt := parseAlterDatabase(t, "ALTER DATABASE mydb COLLATE utf8mb4_general_ci")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "COLLATE" {
			t.Errorf("option name = %s, want COLLATE", stmt.Options[0].Name)
		}
	})

	t.Run("read only", func(t *testing.T) {
		stmt := parseAlterDatabase(t, "ALTER DATABASE mydb READ ONLY = 1")
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "READ ONLY" {
			t.Errorf("option name = %s, want READ ONLY", stmt.Options[0].Name)
		}
		if stmt.Options[0].Value != "1" {
			t.Errorf("option value = %s, want 1", stmt.Options[0].Value)
		}
	})

	t.Run("read only default", func(t *testing.T) {
		stmt := parseAlterDatabase(t, "ALTER DATABASE mydb READ ONLY DEFAULT")
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "READ ONLY" {
			t.Errorf("option name = %s, want READ ONLY", stmt.Options[0].Name)
		}
	})

	t.Run("encryption", func(t *testing.T) {
		stmt := parseAlterDatabase(t, "ALTER DATABASE mydb ENCRYPTION = 'Y'")
		if len(stmt.Options) != 1 {
			t.Fatalf("Options count = %d, want 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "ENCRYPTION" {
			t.Errorf("option name = %s, want ENCRYPTION", stmt.Options[0].Name)
		}
		if stmt.Options[0].Value != "Y" {
			t.Errorf("option value = %s, want Y", stmt.Options[0].Value)
		}
	})
}

func TestParseCreateDatabaseEncryption(t *testing.T) {
	stmt := parseCreateDatabase(t, "CREATE DATABASE mydb DEFAULT ENCRYPTION = 'N'")
	if len(stmt.Options) != 1 {
		t.Fatalf("Options count = %d, want 1", len(stmt.Options))
	}
	if stmt.Options[0].Name != "ENCRYPTION" {
		t.Errorf("option name = %s, want ENCRYPTION", stmt.Options[0].Name)
	}
	if stmt.Options[0].Value != "N" {
		t.Errorf("option value = %s, want N", stmt.Options[0].Value)
	}
}

func parseDropDatabase(t *testing.T, input string) *ast.DropDatabaseStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume DROP
	stmt, err := p.parseDropDatabaseStmt()
	if err != nil {
		t.Fatalf("parseDropDatabaseStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseDropDatabase(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseDropDatabase(t, "DROP DATABASE mydb")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if stmt.IfExists {
			t.Error("IfExists should be false")
		}
	})

	t.Run("if exists", func(t *testing.T) {
		stmt := parseDropDatabase(t, "DROP DATABASE IF EXISTS mydb")
		if stmt.Name != "mydb" {
			t.Errorf("Name = %s, want mydb", stmt.Name)
		}
		if !stmt.IfExists {
			t.Error("IfExists should be true")
		}
	})
}

// ============================================================================
// Batch 12: DROP TABLE/INDEX/VIEW, TRUNCATE, RENAME
// ============================================================================

func parseDropTable(t *testing.T, input string) *ast.DropTableStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // p.cur = DROP
	p.advance() // p.cur = TABLE or TEMPORARY
	temporary := false
	if p.cur.Type == kwTEMPORARY {
		temporary = true
		p.advance() // p.cur = TABLE
	}
	stmt, err := p.parseDropTableStmt(temporary)
	if err != nil {
		t.Fatalf("parseDropTableStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseDropTable(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseDropTable(t, "DROP TABLE users")
		if len(stmt.Tables) != 1 {
			t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
		}
		if stmt.Tables[0].Name != "users" {
			t.Errorf("table = %s, want users", stmt.Tables[0].Name)
		}
	})

	t.Run("multiple if exists cascade", func(t *testing.T) {
		stmt := parseDropTable(t, "DROP TABLE IF EXISTS t1, t2 CASCADE")
		if !stmt.IfExists {
			t.Error("IfExists should be true")
		}
		if len(stmt.Tables) != 2 {
			t.Fatalf("Tables count = %d, want 2", len(stmt.Tables))
		}
		if !stmt.Cascade {
			t.Error("Cascade should be true")
		}
	})
}

func parseDropIndex(t *testing.T, input string) *ast.DropIndexStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume DROP, now at INDEX
	stmt, err := p.parseDropIndexStmt()
	if err != nil {
		t.Fatalf("parseDropIndexStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseDropIndex(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseDropIndex(t, "DROP INDEX idx_name ON users")
		if stmt.Name != "idx_name" {
			t.Errorf("Name = %s, want idx_name", stmt.Name)
		}
		if stmt.Table.Name != "users" {
			t.Errorf("Table = %s, want users", stmt.Table.Name)
		}
	})
}

func parseDropView(t *testing.T, input string) *ast.DropViewStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume DROP, now at VIEW
	stmt, err := p.parseDropViewStmt()
	if err != nil {
		t.Fatalf("parseDropViewStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseDropView(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseDropView(t, "DROP VIEW my_view")
		if len(stmt.Views) != 1 {
			t.Fatalf("Views count = %d, want 1", len(stmt.Views))
		}
		if stmt.Views[0].Name != "my_view" {
			t.Errorf("view = %s, want my_view", stmt.Views[0].Name)
		}
	})

	t.Run("if exists", func(t *testing.T) {
		stmt := parseDropView(t, "DROP VIEW IF EXISTS v1, v2")
		if !stmt.IfExists {
			t.Error("IfExists should be true")
		}
		if len(stmt.Views) != 2 {
			t.Fatalf("Views count = %d, want 2", len(stmt.Views))
		}
	})
}

func parseTruncate(t *testing.T, input string) *ast.TruncateStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // now at TRUNCATE
	stmt, err := p.parseTruncateStmt()
	if err != nil {
		t.Fatalf("parseTruncateStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseTruncate(t *testing.T) {
	t.Run("with TABLE keyword", func(t *testing.T) {
		stmt := parseTruncate(t, "TRUNCATE TABLE users")
		if len(stmt.Tables) != 1 {
			t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
		}
		if stmt.Tables[0].Name != "users" {
			t.Errorf("table = %s, want users", stmt.Tables[0].Name)
		}
	})

	t.Run("without TABLE keyword", func(t *testing.T) {
		stmt := parseTruncate(t, "TRUNCATE users")
		if len(stmt.Tables) != 1 {
			t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
		}
		if stmt.Tables[0].Name != "users" {
			t.Errorf("table = %s, want users", stmt.Tables[0].Name)
		}
	})
}

func parseRenameTable(t *testing.T, input string) *ast.RenameTableStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume RENAME, now at TABLE
	stmt, err := p.parseRenameTableStmt()
	if err != nil {
		t.Fatalf("parseRenameTableStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseRenameTable(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		stmt := parseRenameTable(t, "RENAME TABLE old_t TO new_t")
		if len(stmt.Pairs) != 1 {
			t.Fatalf("Pairs count = %d, want 1", len(stmt.Pairs))
		}
		if stmt.Pairs[0].Old.Name != "old_t" {
			t.Errorf("Old = %s, want old_t", stmt.Pairs[0].Old.Name)
		}
		if stmt.Pairs[0].New.Name != "new_t" {
			t.Errorf("New = %s, want new_t", stmt.Pairs[0].New.Name)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		stmt := parseRenameTable(t, "RENAME TABLE a TO b, c TO d")
		if len(stmt.Pairs) != 2 {
			t.Fatalf("Pairs count = %d, want 2", len(stmt.Pairs))
		}
	})
}

// ============================================================================
// Batch 14: Transaction statements
// ============================================================================

func parseBegin(t *testing.T, input string) *ast.BeginStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseBeginStmt()
	if err != nil {
		t.Fatalf("parseBeginStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseBegin(t *testing.T) {
	t.Run("simple BEGIN", func(t *testing.T) {
		stmt := parseBegin(t, "BEGIN")
		if stmt.ReadOnly || stmt.ReadWrite || stmt.WithConsistentSnapshot {
			t.Error("all flags should be false for simple BEGIN")
		}
	})

	t.Run("BEGIN WORK", func(t *testing.T) {
		stmt := parseBegin(t, "BEGIN WORK")
		if stmt.ReadOnly || stmt.ReadWrite || stmt.WithConsistentSnapshot {
			t.Error("all flags should be false for BEGIN WORK")
		}
	})

	t.Run("START TRANSACTION READ ONLY", func(t *testing.T) {
		stmt := parseBegin(t, "START TRANSACTION READ ONLY")
		if !stmt.ReadOnly {
			t.Error("ReadOnly should be true")
		}
	})

	t.Run("START TRANSACTION WITH CONSISTENT SNAPSHOT", func(t *testing.T) {
		stmt := parseBegin(t, "START TRANSACTION WITH CONSISTENT SNAPSHOT")
		if !stmt.WithConsistentSnapshot {
			t.Error("WithConsistentSnapshot should be true")
		}
	})
}

func parseCommit(t *testing.T, input string) *ast.CommitStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseCommitStmt()
	if err != nil {
		t.Fatalf("parseCommitStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseCommit(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseCommit(t, "COMMIT")
		if stmt.Chain || stmt.Release {
			t.Error("Chain and Release should be false")
		}
	})

	t.Run("and chain", func(t *testing.T) {
		stmt := parseCommit(t, "COMMIT AND CHAIN")
		if !stmt.Chain {
			t.Error("Chain should be true")
		}
	})

	t.Run("release", func(t *testing.T) {
		stmt := parseCommit(t, "COMMIT RELEASE")
		if !stmt.Release {
			t.Error("Release should be true")
		}
	})
}

func parseRollback(t *testing.T, input string) *ast.RollbackStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseRollbackStmt()
	if err != nil {
		t.Fatalf("parseRollbackStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseRollback(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseRollback(t, "ROLLBACK")
		if stmt.Savepoint != "" {
			t.Error("Savepoint should be empty")
		}
	})

	t.Run("to savepoint", func(t *testing.T) {
		stmt := parseRollback(t, "ROLLBACK TO SAVEPOINT sp1")
		if stmt.Savepoint != "sp1" {
			t.Errorf("Savepoint = %s, want sp1", stmt.Savepoint)
		}
	})

	t.Run("and chain", func(t *testing.T) {
		stmt := parseRollback(t, "ROLLBACK AND CHAIN")
		if !stmt.Chain {
			t.Error("Chain should be true")
		}
	})
}

func parseSavepoint(t *testing.T, input string) *ast.SavepointStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseSavepointStmt()
	if err != nil {
		t.Fatalf("parseSavepointStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseSavepoint(t *testing.T) {
	stmt := parseSavepoint(t, "SAVEPOINT sp1")
	if stmt.Name != "sp1" {
		t.Errorf("Name = %s, want sp1", stmt.Name)
	}
}

func parseLockTables(t *testing.T, input string) *ast.LockTablesStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance() // consume LOCK, now at TABLES
	stmt, err := p.parseLockTablesStmt()
	if err != nil {
		t.Fatalf("parseLockTablesStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseLockTables(t *testing.T) {
	t.Run("read lock", func(t *testing.T) {
		stmt := parseLockTables(t, "LOCK TABLES users READ")
		if len(stmt.Tables) != 1 {
			t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
		}
		if stmt.Tables[0].Table.Name != "users" {
			t.Errorf("table = %s, want users", stmt.Tables[0].Table.Name)
		}
		if stmt.Tables[0].LockType != "READ" {
			t.Errorf("LockType = %s, want READ", stmt.Tables[0].LockType)
		}
	})

	t.Run("write lock", func(t *testing.T) {
		stmt := parseLockTables(t, "LOCK TABLES users WRITE")
		if len(stmt.Tables) != 1 {
			t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
		}
		if stmt.Tables[0].LockType != "WRITE" {
			t.Errorf("LockType = %s, want WRITE", stmt.Tables[0].LockType)
		}
	})

	t.Run("multiple tables", func(t *testing.T) {
		stmt := parseLockTables(t, "LOCK TABLES t1 READ, t2 WRITE")
		if len(stmt.Tables) != 2 {
			t.Fatalf("Tables count = %d, want 2", len(stmt.Tables))
		}
		if stmt.Tables[0].LockType != "READ" {
			t.Errorf("table[0] LockType = %s, want READ", stmt.Tables[0].LockType)
		}
		if stmt.Tables[1].LockType != "WRITE" {
			t.Errorf("table[1] LockType = %s, want WRITE", stmt.Tables[1].LockType)
		}
	})
}

// ============================================================================
// Batch 10: CREATE VIEW
// ============================================================================

func parseCreateView(t *testing.T, input string, orReplace bool) *ast.CreateViewStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime lexer
	p.advance() // skip CREATE
	if orReplace {
		p.advance() // OR
		p.advance() // REPLACE
	}
	stmt, err := p.parseCreateViewStmt(orReplace)
	if err != nil {
		t.Fatalf("parseCreateViewStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseCreateView(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseCreateView(t, "CREATE VIEW v AS SELECT 1", false)
		if stmt.Name == nil || stmt.Name.Name != "v" {
			t.Errorf("Name = %v, want v", stmt.Name)
		}
		if stmt.OrReplace {
			t.Errorf("OrReplace = true, want false")
		}
		if stmt.Select == nil {
			t.Fatal("Select is nil")
		}
		if len(stmt.Columns) != 0 {
			t.Errorf("Columns = %v, want empty", stmt.Columns)
		}
	})

	t.Run("with columns", func(t *testing.T) {
		stmt := parseCreateView(t, "CREATE VIEW v (a, b, c) AS SELECT 1, 2, 3", false)
		if stmt.Name == nil || stmt.Name.Name != "v" {
			t.Errorf("Name = %v, want v", stmt.Name)
		}
		if len(stmt.Columns) != 3 {
			t.Fatalf("Columns count = %d, want 3", len(stmt.Columns))
		}
		if stmt.Columns[0] != "a" || stmt.Columns[1] != "b" || stmt.Columns[2] != "c" {
			t.Errorf("Columns = %v, want [a b c]", stmt.Columns)
		}
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := parseCreateView(t, "CREATE OR REPLACE VIEW v AS SELECT 1", true)
		if !stmt.OrReplace {
			t.Errorf("OrReplace = false, want true")
		}
		if stmt.Name == nil || stmt.Name.Name != "v" {
			t.Errorf("Name = %v, want v", stmt.Name)
		}
	})
}

func TestParseCreateOrReplaceView(t *testing.T) {
	stmt := parseCreateView(t, "CREATE OR REPLACE VIEW v AS SELECT 1", true)
	if !stmt.OrReplace {
		t.Errorf("OrReplace = false, want true")
	}
	if stmt.Name == nil || stmt.Name.Name != "v" {
		t.Errorf("Name = %v, want v", stmt.Name)
	}
	if stmt.Select == nil {
		t.Fatal("Select is nil")
	}
}

// ============================================================================
// Batch 7: CREATE TABLE
// ============================================================================

// helper to parse a CREATE TABLE statement directly
func parseCreateTable(t *testing.T, input string) *ast.CreateTableStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime lexer
	p.advance() // skip CREATE
	temporary := false
	// handle CREATE TEMPORARY TABLE
	if p.cur.Type == kwTEMPORARY {
		p.advance()
		temporary = true
	}
	stmt, err := p.parseCreateTableStmt(temporary)
	if err != nil {
		t.Fatalf("parseCreateTableStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseCreateTable(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT)")
		if stmt.Table == nil || stmt.Table.Name != "t" {
			t.Errorf("Table = %v, want t", stmt.Table)
		}
		if len(stmt.Columns) != 1 {
			t.Fatalf("Columns count = %d, want 1", len(stmt.Columns))
		}
		if stmt.Columns[0].Name != "id" {
			t.Errorf("Column name = %s, want id", stmt.Columns[0].Name)
		}
		if stmt.Columns[0].TypeName == nil || stmt.Columns[0].TypeName.Name != "INT" {
			t.Errorf("Column type = %v, want INT", stmt.Columns[0].TypeName)
		}
	})

	t.Run("multiple columns", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(255))")
		if len(stmt.Columns) != 3 {
			t.Fatalf("Columns count = %d, want 3", len(stmt.Columns))
		}
		if stmt.Columns[0].Name != "id" {
			t.Errorf("Col[0] = %s, want id", stmt.Columns[0].Name)
		}
		if stmt.Columns[1].Name != "name" {
			t.Errorf("Col[1] = %s, want name", stmt.Columns[1].Name)
		}
		if stmt.Columns[2].Name != "email" {
			t.Errorf("Col[2] = %s, want email", stmt.Columns[2].Name)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE IF NOT EXISTS t (id INT)")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})

	t.Run("temporary", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TEMPORARY TABLE t (id INT)")
		if !stmt.Temporary {
			t.Errorf("Temporary = false, want true")
		}
	})

	t.Run("schema qualified", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE mydb.t (id INT)")
		if stmt.Table == nil || stmt.Table.Schema != "mydb" || stmt.Table.Name != "t" {
			t.Errorf("Table = %v, want mydb.t", stmt.Table)
		}
	})

	t.Run("column not null", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT NOT NULL)")
		if len(stmt.Columns) != 1 {
			t.Fatalf("Columns count = %d, want 1", len(stmt.Columns))
		}
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrNotNull {
				found = true
			}
		}
		if !found {
			t.Errorf("NOT NULL constraint not found")
		}
	})

	t.Run("column null", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT NULL)")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrNull {
				found = true
			}
		}
		if !found {
			t.Errorf("NULL constraint not found")
		}
	})

	t.Run("column default", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (status INT DEFAULT 0)")
		col := stmt.Columns[0]
		if col.DefaultValue == nil {
			t.Fatal("DefaultValue is nil")
		}
	})

	t.Run("column auto increment", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT AUTO_INCREMENT)")
		if !stmt.Columns[0].AutoIncrement {
			t.Errorf("AutoIncrement = false, want true")
		}
	})

	t.Run("column primary key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT PRIMARY KEY)")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrPrimaryKey {
				found = true
			}
		}
		if !found {
			t.Errorf("PRIMARY KEY constraint not found")
		}
	})

	t.Run("column unique", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (email VARCHAR(255) UNIQUE)")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrUnique {
				found = true
			}
		}
		if !found {
			t.Errorf("UNIQUE constraint not found")
		}
	})

	t.Run("column comment", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT COMMENT 'primary key')")
		if stmt.Columns[0].Comment != "primary key" {
			t.Errorf("Comment = %q, want 'primary key'", stmt.Columns[0].Comment)
		}
	})

	t.Run("column on update", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
		col := stmt.Columns[0]
		if col.OnUpdate == nil {
			t.Fatal("OnUpdate is nil")
		}
	})

	t.Run("column references", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (user_id INT REFERENCES users(id))")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrReferences {
				found = true
				if c.RefTable == nil || c.RefTable.Name != "users" {
					t.Errorf("RefTable = %v, want users", c.RefTable)
				}
				if len(c.RefColumns) != 1 || c.RefColumns[0] != "id" {
					t.Errorf("RefColumns = %v, want [id]", c.RefColumns)
				}
			}
		}
		if !found {
			t.Errorf("REFERENCES constraint not found")
		}
	})

	t.Run("column check", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (age INT CHECK (age > 0))")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrCheck {
				found = true
			}
		}
		if !found {
			t.Errorf("CHECK constraint not found")
		}
	})

	t.Run("generated column virtual", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) VIRTUAL)")
		col := stmt.Columns[2]
		if col.Generated == nil {
			t.Fatal("Generated is nil")
		}
		if col.Generated.Stored {
			t.Errorf("Stored = true, want false")
		}
	})

	t.Run("generated column stored", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) STORED)")
		col := stmt.Columns[2]
		if col.Generated == nil {
			t.Fatal("Generated is nil")
		}
		if !col.Generated.Stored {
			t.Errorf("Stored = false, want true")
		}
	})

	t.Run("column collate on varchar", func(t *testing.T) {
		// COLLATE on VARCHAR is consumed by the data type parser, not as a column constraint
		stmt := parseCreateTable(t, "CREATE TABLE t (name VARCHAR(100) COLLATE utf8mb4_unicode_ci)")
		col := stmt.Columns[0]
		if col.TypeName == nil || col.TypeName.Collate != "utf8mb4_unicode_ci" {
			t.Errorf("TypeName.Collate = %v, want utf8mb4_unicode_ci", col.TypeName)
		}
	})

	t.Run("column collate on int", func(t *testing.T) {
		// COLLATE on non-string type is parsed as a column constraint
		stmt := parseCreateTable(t, "CREATE TABLE t (name INT COLLATE utf8mb4_unicode_ci)")
		col := stmt.Columns[0]
		found := false
		for _, c := range col.Constraints {
			if c.Type == ast.ColConstrCollate {
				found = true
				if c.Name != "utf8mb4_unicode_ci" {
					t.Errorf("Collation = %s, want utf8mb4_unicode_ci", c.Name)
				}
			}
		}
		if !found {
			t.Errorf("COLLATE constraint not found")
		}
	})

	t.Run("full column definition", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'pk')")
		col := stmt.Columns[0]
		if col.TypeName == nil || col.TypeName.Name != "BIGINT" {
			t.Errorf("TypeName = %v, want BIGINT", col.TypeName)
		}
		if !col.TypeName.Unsigned {
			t.Errorf("Unsigned = false, want true")
		}
		if !col.AutoIncrement {
			t.Errorf("AutoIncrement = false, want true")
		}
		if col.Comment != "pk" {
			t.Errorf("Comment = %q, want 'pk'", col.Comment)
		}
	})
}

func TestParseCreateTableConstraints(t *testing.T) {
	t.Run("primary key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, PRIMARY KEY (id))")
		if len(stmt.Constraints) != 1 {
			t.Fatalf("Constraints count = %d, want 1", len(stmt.Constraints))
		}
		if stmt.Constraints[0].Type != ast.ConstrPrimaryKey {
			t.Errorf("Type = %d, want ConstrPrimaryKey", stmt.Constraints[0].Type)
		}
		if len(stmt.Constraints[0].Columns) != 1 || stmt.Constraints[0].Columns[0] != "id" {
			t.Errorf("Columns = %v, want [id]", stmt.Constraints[0].Columns)
		}
	})

	t.Run("composite primary key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b))")
		if len(stmt.Constraints[0].Columns) != 2 {
			t.Fatalf("PK columns = %d, want 2", len(stmt.Constraints[0].Columns))
		}
	})

	t.Run("named constraint", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, CONSTRAINT pk_t PRIMARY KEY (id))")
		if stmt.Constraints[0].Name != "pk_t" {
			t.Errorf("Name = %s, want pk_t", stmt.Constraints[0].Name)
		}
	})

	t.Run("unique index", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (email VARCHAR(255), UNIQUE KEY idx_email (email))")
		if len(stmt.Constraints) != 1 {
			t.Fatalf("Constraints count = %d, want 1", len(stmt.Constraints))
		}
		if stmt.Constraints[0].Type != ast.ConstrUnique {
			t.Errorf("Type = %d, want ConstrUnique", stmt.Constraints[0].Type)
		}
		if stmt.Constraints[0].Name != "idx_email" {
			t.Errorf("Name = %s, want idx_email", stmt.Constraints[0].Name)
		}
	})

	t.Run("index", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (name VARCHAR(100), INDEX idx_name (name))")
		if stmt.Constraints[0].Type != ast.ConstrIndex {
			t.Errorf("Type = %d, want ConstrIndex", stmt.Constraints[0].Type)
		}
	})

	t.Run("key alias", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (name VARCHAR(100), KEY idx_name (name))")
		if stmt.Constraints[0].Type != ast.ConstrIndex {
			t.Errorf("Type = %d, want ConstrIndex", stmt.Constraints[0].Type)
		}
	})

	t.Run("fulltext index", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (content TEXT, FULLTEXT INDEX ft_content (content))")
		if stmt.Constraints[0].Type != ast.ConstrFulltextIndex {
			t.Errorf("Type = %d, want ConstrFulltextIndex", stmt.Constraints[0].Type)
		}
	})

	t.Run("spatial index", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (geo GEOMETRY, SPATIAL INDEX sp_geo (geo))")
		if stmt.Constraints[0].Type != ast.ConstrSpatialIndex {
			t.Errorf("Type = %d, want ConstrSpatialIndex", stmt.Constraints[0].Type)
		}
	})

	t.Run("foreign key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE orders (id INT, user_id INT, FOREIGN KEY (user_id) REFERENCES users(id))")
		fk := stmt.Constraints[0]
		if fk.Type != ast.ConstrForeignKey {
			t.Errorf("Type = %d, want ConstrForeignKey", fk.Type)
		}
		if len(fk.Columns) != 1 || fk.Columns[0] != "user_id" {
			t.Errorf("Columns = %v, want [user_id]", fk.Columns)
		}
		if fk.RefTable == nil || fk.RefTable.Name != "users" {
			t.Errorf("RefTable = %v, want users", fk.RefTable)
		}
		if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
			t.Errorf("RefColumns = %v, want [id]", fk.RefColumns)
		}
	})

	t.Run("foreign key with actions", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE orders (id INT, user_id INT, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE SET NULL)")
		fk := stmt.Constraints[0]
		if fk.OnDelete != ast.RefActCascade {
			t.Errorf("OnDelete = %d, want RefActCascade", fk.OnDelete)
		}
		if fk.OnUpdate != ast.RefActSetNull {
			t.Errorf("OnUpdate = %d, want RefActSetNull", fk.OnUpdate)
		}
	})

	t.Run("named foreign key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, pid INT, CONSTRAINT fk_parent FOREIGN KEY (pid) REFERENCES parent(id))")
		if stmt.Constraints[0].Name != "fk_parent" {
			t.Errorf("Name = %s, want fk_parent", stmt.Constraints[0].Name)
		}
	})

	t.Run("check constraint", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (age INT, CHECK (age >= 0))")
		if stmt.Constraints[0].Type != ast.ConstrCheck {
			t.Errorf("Type = %d, want ConstrCheck", stmt.Constraints[0].Type)
		}
		if stmt.Constraints[0].Expr == nil {
			t.Fatal("Check expr is nil")
		}
	})

	t.Run("using btree", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, INDEX idx_id USING BTREE (id))")
		if stmt.Constraints[0].IndexType != "BTREE" {
			t.Errorf("IndexType = %s, want BTREE", stmt.Constraints[0].IndexType)
		}
	})

	t.Run("multiple constraints", func(t *testing.T) {
		stmt := parseCreateTable(t, `CREATE TABLE t (
			id INT,
			name VARCHAR(100),
			email VARCHAR(255),
			PRIMARY KEY (id),
			UNIQUE KEY idx_email (email),
			INDEX idx_name (name)
		)`)
		if len(stmt.Constraints) != 3 {
			t.Fatalf("Constraints count = %d, want 3", len(stmt.Constraints))
		}
	})
}

func TestParseCreateTableOptions(t *testing.T) {
	t.Run("engine", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) ENGINE=InnoDB")
		if len(stmt.Options) < 1 {
			t.Fatalf("Options count = %d, want >= 1", len(stmt.Options))
		}
		if stmt.Options[0].Name != "ENGINE" || stmt.Options[0].Value != "InnoDB" {
			t.Errorf("Option = %s=%s, want ENGINE=InnoDB", stmt.Options[0].Name, stmt.Options[0].Value)
		}
	})

	t.Run("charset", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) DEFAULT CHARSET=utf8mb4")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "CHARSET" && opt.Value == "utf8mb4" {
				found = true
			}
		}
		if !found {
			t.Errorf("CHARSET=utf8mb4 option not found, opts=%v", stmt.Options)
		}
	})

	t.Run("collate", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) COLLATE=utf8mb4_unicode_ci")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "COLLATE" && opt.Value == "utf8mb4_unicode_ci" {
				found = true
			}
		}
		if !found {
			t.Errorf("COLLATE option not found")
		}
	})

	t.Run("comment", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) COMMENT='test table'")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "COMMENT" && opt.Value == "test table" {
				found = true
			}
		}
		if !found {
			t.Errorf("COMMENT option not found")
		}
	})

	t.Run("auto increment value", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) AUTO_INCREMENT=100")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "AUTO_INCREMENT" {
				found = true
			}
		}
		if !found {
			t.Errorf("AUTO_INCREMENT option not found")
		}
	})

	t.Run("row format", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) ROW_FORMAT=DYNAMIC")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "ROW_FORMAT" && opt.Value == "DYNAMIC" {
				found = true
			}
		}
		if !found {
			t.Errorf("ROW_FORMAT option not found")
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='test'")
		if len(stmt.Options) < 3 {
			t.Errorf("Options count = %d, want >= 3", len(stmt.Options))
		}
	})

	t.Run("options with comma separator", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) ENGINE=InnoDB, CHARSET=utf8mb4")
		if len(stmt.Options) < 2 {
			t.Errorf("Options count = %d, want >= 2", len(stmt.Options))
		}
	})

	t.Run("character set", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) CHARACTER SET utf8mb4")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "CHARACTER SET" && opt.Value == "utf8mb4" {
				found = true
			}
		}
		if !found {
			t.Errorf("CHARACTER SET option not found")
		}
	})

	t.Run("start transaction", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) ENGINE=InnoDB START TRANSACTION")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "START TRANSACTION" {
				found = true
			}
		}
		if !found {
			t.Errorf("START TRANSACTION option not found in options: %v", stmt.Options)
		}
	})
}

func TestParseCreateTableLike(t *testing.T) {
	t.Run("simple like", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t2 LIKE t1")
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
		if stmt.Like.Name != "t1" {
			t.Errorf("Like.Name = %s, want t1", stmt.Like.Name)
		}
	})

	t.Run("like with schema", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE mydb.t2 LIKE mydb.t1")
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
		if stmt.Like.Schema != "mydb" || stmt.Like.Name != "t1" {
			t.Errorf("Like = %v, want mydb.t1", stmt.Like)
		}
	})

	t.Run("like if not exists", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE IF NOT EXISTS t2 LIKE t1")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("parenthesized like", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t2 (LIKE t1)")
		if stmt.Like == nil {
			t.Fatal("Like is nil for parenthesized LIKE")
		}
		if stmt.Like.Name != "t1" {
			t.Errorf("Like.Name = %q, want %q", stmt.Like.Name, "t1")
		}
	})

	t.Run("create as select", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t2 AS SELECT * FROM t1")
		if stmt.Select == nil {
			t.Fatal("Select is nil")
		}
	})

	t.Run("create as select with columns", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t2 (id INT) AS SELECT id FROM t1")
		if stmt.Select == nil {
			t.Fatal("Select is nil")
		}
		if len(stmt.Columns) != 1 {
			t.Errorf("Columns count = %d, want 1", len(stmt.Columns))
		}
	})

	t.Run("partition by hash", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4")
		if stmt.Partitions == nil {
			t.Fatal("Partitions is nil")
		}
		if stmt.Partitions.Type != ast.PartitionHash {
			t.Errorf("Partition type = %d, want PartitionHash", stmt.Partitions.Type)
		}
		if stmt.Partitions.NumParts != 4 {
			t.Errorf("NumParts = %d, want 4", stmt.Partitions.NumParts)
		}
	})

	t.Run("partition by key", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT) PARTITION BY KEY(id)")
		if stmt.Partitions == nil {
			t.Fatal("Partitions is nil")
		}
		if stmt.Partitions.Type != ast.PartitionKey {
			t.Errorf("Partition type = %d, want PartitionKey", stmt.Partitions.Type)
		}
	})

	t.Run("partition by range", func(t *testing.T) {
		stmt := parseCreateTable(t, `CREATE TABLE t (id INT, created DATE) PARTITION BY RANGE(id) (
			PARTITION p0 VALUES LESS THAN (100),
			PARTITION p1 VALUES LESS THAN (200),
			PARTITION pmax VALUES LESS THAN MAXVALUE
		)`)
		if stmt.Partitions == nil {
			t.Fatal("Partitions is nil")
		}
		if stmt.Partitions.Type != ast.PartitionRange {
			t.Errorf("Partition type = %d, want PartitionRange", stmt.Partitions.Type)
		}
		if len(stmt.Partitions.Partitions) != 3 {
			t.Fatalf("Partition defs = %d, want 3", len(stmt.Partitions.Partitions))
		}
	})

	t.Run("partition by list", func(t *testing.T) {
		stmt := parseCreateTable(t, `CREATE TABLE t (id INT, region INT) PARTITION BY LIST(region) (
			PARTITION p_east VALUES IN (1, 2, 3),
			PARTITION p_west VALUES IN (4, 5, 6)
		)`)
		if stmt.Partitions == nil {
			t.Fatal("Partitions is nil")
		}
		if stmt.Partitions.Type != ast.PartitionList {
			t.Errorf("Partition type = %d, want PartitionList", stmt.Partitions.Type)
		}
		if len(stmt.Partitions.Partitions) != 2 {
			t.Fatalf("Partition defs = %d, want 2", len(stmt.Partitions.Partitions))
		}
	})

	t.Run("partition by range columns", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b)")
		if stmt.Partitions == nil {
			t.Fatal("Partitions is nil")
		}
		if len(stmt.Partitions.Columns) != 2 {
			t.Errorf("Partition columns = %d, want 2", len(stmt.Partitions.Columns))
		}
	})

	t.Run("foreign key restrict no action", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, pid INT, FOREIGN KEY (pid) REFERENCES p(id) ON DELETE RESTRICT ON UPDATE NO ACTION)")
		fk := stmt.Constraints[0]
		if fk.OnDelete != ast.RefActRestrict {
			t.Errorf("OnDelete = %d, want RefActRestrict", fk.OnDelete)
		}
		if fk.OnUpdate != ast.RefActNoAction {
			t.Errorf("OnUpdate = %d, want RefActNoAction", fk.OnUpdate)
		}
	})

	t.Run("foreign key set default", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, pid INT, FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET DEFAULT)")
		fk := stmt.Constraints[0]
		if fk.OnDelete != ast.RefActSetDefault {
			t.Errorf("OnDelete = %d, want RefActSetDefault", fk.OnDelete)
		}
	})
}

// ============================================================================
// Batch 59: CREATE TABLE depth fixes
// ============================================================================

func TestParseCreateTableColumnInvisible(t *testing.T) {
	t.Run("column visible", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT VISIBLE, name VARCHAR(100) INVISIBLE)")
		if len(stmt.Columns) != 2 {
			t.Fatalf("Columns count = %d, want 2", len(stmt.Columns))
		}
		// Check id column has VISIBLE constraint
		found := false
		for _, c := range stmt.Columns[0].Constraints {
			if c.Type == ast.ColConstrVisible {
				found = true
			}
		}
		if !found {
			t.Error("expected ColConstrVisible for id column")
		}
		// Check name column has INVISIBLE constraint
		found = false
		for _, c := range stmt.Columns[1].Constraints {
			if c.Type == ast.ColConstrInvisible {
				found = true
			}
		}
		if !found {
			t.Error("expected ColConstrInvisible for name column")
		}
	})
}

func TestParseCreateTableCheckNotEnforced(t *testing.T) {
	t.Run("column check not enforced", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, CONSTRAINT chk_id CHECK (id > 0) NOT ENFORCED)")
		if len(stmt.Constraints) != 1 {
			t.Fatalf("Constraints count = %d, want 1", len(stmt.Constraints))
		}
		chk := stmt.Constraints[0]
		if chk.Type != ast.ConstrCheck {
			t.Errorf("Type = %d, want ConstrCheck", chk.Type)
		}
		if !chk.NotEnforced {
			t.Error("expected NotEnforced = true")
		}
	})

	t.Run("check enforced (default)", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT, CHECK (id > 0) ENFORCED)")
		chk := stmt.Constraints[0]
		if chk.NotEnforced {
			t.Error("expected NotEnforced = false for ENFORCED")
		}
	})

	t.Run("column-level check not enforced", func(t *testing.T) {
		stmt := parseCreateTable(t, "CREATE TABLE t (id INT CHECK (id > 0) NOT ENFORCED)")
		found := false
		for _, c := range stmt.Columns[0].Constraints {
			if c.Type == ast.ColConstrCheck && c.NotEnforced {
				found = true
			}
		}
		if !found {
			t.Error("expected column-level CHECK NOT ENFORCED")
		}
	})
}

func TestParseCreateTableFunctionalIndex(t *testing.T) {
	t.Run("expression-based index", func(t *testing.T) {
		ParseAndCheck(t, "CREATE TABLE t (id INT, name VARCHAR(100), INDEX idx_upper ((UPPER(name))))")
	})

	t.Run("functional key in unique index", func(t *testing.T) {
		ParseAndCheck(t, "CREATE TABLE t (a INT, b INT, UNIQUE INDEX idx_expr ((a + b)))")
	})
}

func TestParseCreateTableSubPartition(t *testing.T) {
	t.Run("subpartition by hash", func(t *testing.T) {
		ParseAndCheck(t, "CREATE TABLE t (id INT, ts DATE) PARTITION BY RANGE (YEAR(ts)) SUBPARTITION BY HASH (id) (PARTITION p0 VALUES LESS THAN (2000), PARTITION p1 VALUES LESS THAN MAXVALUE)")
	})

	t.Run("subpartition definitions", func(t *testing.T) {
		ParseAndCheck(t, "CREATE TABLE t (id INT, ts DATE) PARTITION BY RANGE (YEAR(ts)) SUBPARTITION BY HASH (id) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (2000) (SUBPARTITION s0, SUBPARTITION s1), PARTITION p1 VALUES LESS THAN MAXVALUE (SUBPARTITION s2, SUBPARTITION s3))")
	})
}

// ============================================================================
// Batch 8: ALTER TABLE
// ============================================================================

func parseAlterTable(t *testing.T, input string) *ast.AlterTableStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime lexer
	p.advance() // skip ALTER
	stmt, err := p.parseAlterTableStmt()
	if err != nil {
		t.Fatalf("parseAlterTableStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseAlterTableAddColumn(t *testing.T) {
	t.Run("add column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN name VARCHAR(100)")
		if len(stmt.Commands) != 1 {
			t.Fatalf("Commands count = %d, want 1", len(stmt.Commands))
		}
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddColumn {
			t.Errorf("Type = %d, want ATAddColumn", cmd.Type)
		}
		if cmd.Column == nil || cmd.Column.Name != "name" {
			t.Errorf("Column = %v, want name", cmd.Column)
		}
	})

	t.Run("add column without column keyword", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD name VARCHAR(100)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddColumn {
			t.Errorf("Type = %d, want ATAddColumn", cmd.Type)
		}
	})

	t.Run("add column first", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN id INT FIRST")
		cmd := stmt.Commands[0]
		if !cmd.First {
			t.Errorf("First = false, want true")
		}
	})

	t.Run("add column after", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN email VARCHAR(255) AFTER name")
		cmd := stmt.Commands[0]
		if cmd.After != "name" {
			t.Errorf("After = %s, want name", cmd.After)
		}
	})

	t.Run("add constraint primary key", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD PRIMARY KEY (id)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddConstraint {
			t.Errorf("Type = %d, want ATAddConstraint", cmd.Type)
		}
		if cmd.Constraint == nil || cmd.Constraint.Type != ast.ConstrPrimaryKey {
			t.Errorf("Constraint type = %v, want ConstrPrimaryKey", cmd.Constraint)
		}
	})

	t.Run("add primary key with name", func(t *testing.T) {
		// MySQL 8.0 accepts PRIMARY KEY index_name (cols) — name is silently ignored
		stmt := parseAlterTable(t, "ALTER TABLE t ADD PRIMARY KEY pk_a (a)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddConstraint {
			t.Errorf("Type = %d, want ATAddConstraint", cmd.Type)
		}
		if cmd.Constraint == nil || cmd.Constraint.Type != ast.ConstrPrimaryKey {
			t.Errorf("Constraint type = %v, want ConstrPrimaryKey", cmd.Constraint)
		}
	})

	t.Run("add constraint primary key with name", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t ADD CONSTRAINT pk_a PRIMARY KEY (a)")
	})

	t.Run("add primary key with name in create table", func(t *testing.T) {
		ParseAndCheck(t, "CREATE TABLE t (a INT, PRIMARY KEY pk_a (a))")
	})

	t.Run("add constraint unique", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD UNIQUE KEY idx_email (email)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddConstraint {
			t.Errorf("Type = %d, want ATAddConstraint", cmd.Type)
		}
		if cmd.Constraint.Type != ast.ConstrUnique {
			t.Errorf("Constraint type = %d, want ConstrUnique", cmd.Constraint.Type)
		}
	})

	t.Run("add index", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD INDEX idx_name (name)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddConstraint {
			t.Errorf("Type = %d, want ATAddConstraint", cmd.Type)
		}
		if cmd.Constraint.Type != ast.ConstrIndex {
			t.Errorf("Constraint type = %d, want ConstrIndex", cmd.Constraint.Type)
		}
	})

	t.Run("add foreign key", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)")
		cmd := stmt.Commands[0]
		if cmd.Constraint.Name != "fk_user" {
			t.Errorf("Name = %s, want fk_user", cmd.Constraint.Name)
		}
		if cmd.Constraint.Type != ast.ConstrForeignKey {
			t.Errorf("Type = %d, want ConstrForeignKey", cmd.Constraint.Type)
		}
	})

	t.Run("add check", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD CHECK (age >= 0)")
		cmd := stmt.Commands[0]
		if cmd.Constraint.Type != ast.ConstrCheck {
			t.Errorf("Type = %d, want ConstrCheck", cmd.Constraint.Type)
		}
	})
}

func TestParseAlterTableDropColumn(t *testing.T) {
	t.Run("drop column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP COLUMN name")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropColumn {
			t.Errorf("Type = %d, want ATDropColumn", cmd.Type)
		}
		if cmd.Name != "name" {
			t.Errorf("Name = %s, want name", cmd.Name)
		}
	})

	t.Run("drop column without keyword", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP name")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropColumn {
			t.Errorf("Type = %d, want ATDropColumn", cmd.Type)
		}
	})

	t.Run("drop primary key", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP PRIMARY KEY")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropConstraint {
			t.Errorf("Type = %d, want ATDropConstraint", cmd.Type)
		}
		if cmd.Name != "PRIMARY" {
			t.Errorf("Name = %s, want PRIMARY", cmd.Name)
		}
	})

	t.Run("drop index", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP INDEX idx_name")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropIndex {
			t.Errorf("Type = %d, want ATDropIndex", cmd.Type)
		}
		if cmd.Name != "idx_name" {
			t.Errorf("Name = %s, want idx_name", cmd.Name)
		}
	})

	t.Run("drop foreign key", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP FOREIGN KEY fk_user")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropConstraint {
			t.Errorf("Type = %d, want ATDropConstraint", cmd.Type)
		}
		if cmd.Name != "fk_user" {
			t.Errorf("Name = %s, want fk_user", cmd.Name)
		}
	})
}

func TestParseAlterTableModifyColumn(t *testing.T) {
	t.Run("modify column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t MODIFY COLUMN name VARCHAR(200)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATModifyColumn {
			t.Errorf("Type = %d, want ATModifyColumn", cmd.Type)
		}
		if cmd.Name != "name" {
			t.Errorf("Name = %s, want name", cmd.Name)
		}
	})

	t.Run("modify without column keyword", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t MODIFY name VARCHAR(200)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATModifyColumn {
			t.Errorf("Type = %d, want ATModifyColumn", cmd.Type)
		}
	})

	t.Run("modify with first", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t MODIFY COLUMN id INT FIRST")
		cmd := stmt.Commands[0]
		if !cmd.First {
			t.Errorf("First = false, want true")
		}
	})

	t.Run("modify with after", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t MODIFY COLUMN email VARCHAR(255) AFTER name")
		cmd := stmt.Commands[0]
		if cmd.After != "name" {
			t.Errorf("After = %s, want name", cmd.After)
		}
	})

	t.Run("change column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t CHANGE COLUMN old_name new_name VARCHAR(200)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATChangeColumn {
			t.Errorf("Type = %d, want ATChangeColumn", cmd.Type)
		}
		if cmd.Name != "old_name" {
			t.Errorf("Name = %s, want old_name", cmd.Name)
		}
		if cmd.NewName != "new_name" {
			t.Errorf("NewName = %s, want new_name", cmd.NewName)
		}
	})
}

func TestParseAlterTableAddIndex(t *testing.T) {
	t.Run("add fulltext index", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD FULLTEXT INDEX ft_content (content)")
		cmd := stmt.Commands[0]
		if cmd.Constraint.Type != ast.ConstrFulltextIndex {
			t.Errorf("Type = %d, want ConstrFulltextIndex", cmd.Constraint.Type)
		}
	})

	t.Run("add spatial index", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD SPATIAL INDEX sp_geo (geo)")
		cmd := stmt.Commands[0]
		if cmd.Constraint.Type != ast.ConstrSpatialIndex {
			t.Errorf("Type = %d, want ConstrSpatialIndex", cmd.Constraint.Type)
		}
	})

	t.Run("multiple commands", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN a INT, ADD COLUMN b INT, DROP COLUMN c")
		if len(stmt.Commands) != 3 {
			t.Fatalf("Commands count = %d, want 3", len(stmt.Commands))
		}
		if stmt.Commands[0].Type != ast.ATAddColumn {
			t.Errorf("cmd[0] = %d, want ATAddColumn", stmt.Commands[0].Type)
		}
		if stmt.Commands[1].Type != ast.ATAddColumn {
			t.Errorf("cmd[1] = %d, want ATAddColumn", stmt.Commands[1].Type)
		}
		if stmt.Commands[2].Type != ast.ATDropColumn {
			t.Errorf("cmd[2] = %d, want ATDropColumn", stmt.Commands[2].Type)
		}
	})
}

func TestParseAlterTableRename(t *testing.T) {
	t.Run("rename table", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t RENAME TO new_t")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRenameTable {
			t.Errorf("Type = %d, want ATRenameTable", cmd.Type)
		}
		if cmd.NewName != "new_t" {
			t.Errorf("NewName = %s, want new_t", cmd.NewName)
		}
	})

	t.Run("rename table as", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t RENAME AS new_t")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRenameTable {
			t.Errorf("Type = %d, want ATRenameTable", cmd.Type)
		}
	})

	t.Run("rename column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t RENAME COLUMN old_col TO new_col")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRenameColumn {
			t.Errorf("Type = %d, want ATRenameColumn", cmd.Type)
		}
		if cmd.Name != "old_col" {
			t.Errorf("Name = %s, want old_col", cmd.Name)
		}
		if cmd.NewName != "new_col" {
			t.Errorf("NewName = %s, want new_col", cmd.NewName)
		}
	})

	t.Run("rename index", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t RENAME INDEX old_idx TO new_idx")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRenameIndex {
			t.Errorf("Type = %d, want ATRenameIndex", cmd.Type)
		}
		if cmd.Name != "old_idx" {
			t.Errorf("Name = %s, want old_idx", cmd.Name)
		}
		if cmd.NewName != "new_idx" {
			t.Errorf("NewName = %s, want new_idx", cmd.NewName)
		}
	})

	t.Run("convert charset", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATConvertCharset {
			t.Errorf("Type = %d, want ATConvertCharset", cmd.Type)
		}
		if cmd.Name != "utf8mb4" {
			t.Errorf("Name = %s, want utf8mb4", cmd.Name)
		}
		if cmd.NewName != "utf8mb4_unicode_ci" {
			t.Errorf("NewName = %s, want utf8mb4_unicode_ci", cmd.NewName)
		}
	})

	t.Run("algorithm option", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN a INT, ALGORITHM=INPLACE")
		if len(stmt.Commands) != 2 {
			t.Fatalf("Commands count = %d, want 2", len(stmt.Commands))
		}
		cmd := stmt.Commands[1]
		if cmd.Type != ast.ATAlgorithm {
			t.Errorf("Type = %d, want ATAlgorithm", cmd.Type)
		}
	})

	t.Run("lock option", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN a INT, LOCK=NONE")
		cmd := stmt.Commands[1]
		if cmd.Type != ast.ATLock {
			t.Errorf("Type = %d, want ATLock", cmd.Type)
		}
	})

	t.Run("alter column set default", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN status SET DEFAULT 0")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterColumnDefault {
			t.Errorf("Type = %d, want ATAlterColumnDefault", cmd.Type)
		}
		if cmd.Name != "status" {
			t.Errorf("Name = %s, want status", cmd.Name)
		}
	})

	t.Run("alter column drop default", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN status DROP DEFAULT")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterColumnDefault {
			t.Errorf("Type = %d, want ATAlterColumnDefault", cmd.Type)
		}
	})

	t.Run("table option engine", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ENGINE=InnoDB")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATTableOption {
			t.Errorf("Type = %d, want ATTableOption", cmd.Type)
		}
		if cmd.Option == nil || cmd.Option.Name != "ENGINE" {
			t.Errorf("Option = %v, want ENGINE", cmd.Option)
		}
	})
}

// ============================================================================
// Batch 57: ALTER TABLE partition operations
// ============================================================================

func TestParseAlterTableAddPartition(t *testing.T) {
	t.Run("add partition range", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD PARTITION (PARTITION p3 VALUES LESS THAN (2000))")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddPartition {
			t.Errorf("Type = %d, want ATAddPartition", cmd.Type)
		}
		if len(cmd.PartitionDefs) != 1 {
			t.Fatalf("PartitionDefs count = %d, want 1", len(cmd.PartitionDefs))
		}
		if cmd.PartitionDefs[0].Name != "p3" {
			t.Errorf("PartitionDefs[0].Name = %s, want p3", cmd.PartitionDefs[0].Name)
		}
	})

	t.Run("add partition list", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD PARTITION (PARTITION p4 VALUES IN (10, 20, 30))")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddPartition {
			t.Errorf("Type = %d, want ATAddPartition", cmd.Type)
		}
		if len(cmd.PartitionDefs) != 1 {
			t.Fatalf("PartitionDefs count = %d, want 1", len(cmd.PartitionDefs))
		}
	})

	t.Run("add multiple partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD PARTITION (PARTITION p3 VALUES LESS THAN (2000), PARTITION p4 VALUES LESS THAN MAXVALUE)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddPartition {
			t.Errorf("Type = %d, want ATAddPartition", cmd.Type)
		}
		if len(cmd.PartitionDefs) != 2 {
			t.Fatalf("PartitionDefs count = %d, want 2", len(cmd.PartitionDefs))
		}
	})
}

func TestParseAlterTableDropPartition(t *testing.T) {
	t.Run("drop single partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP PARTITION p1")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropPartition {
			t.Errorf("Type = %d, want ATDropPartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 1 || cmd.PartitionNames[0] != "p1" {
			t.Errorf("PartitionNames = %v, want [p1]", cmd.PartitionNames)
		}
	})

	t.Run("drop multiple partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DROP PARTITION p1, p2, p3")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropPartition {
			t.Errorf("Type = %d, want ATDropPartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 3 {
			t.Errorf("PartitionNames count = %d, want 3", len(cmd.PartitionNames))
		}
	})
}

func TestParseAlterTableCoalescePartition(t *testing.T) {
	t.Run("coalesce partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t COALESCE PARTITION 4")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATCoalescePartition {
			t.Errorf("Type = %d, want ATCoalescePartition", cmd.Type)
		}
		if cmd.Number != 4 {
			t.Errorf("Number = %d, want 4", cmd.Number)
		}
	})
}

func TestParseAlterTableReorganizePartition(t *testing.T) {
	t.Run("reorganize partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REORGANIZE PARTITION p1, p2 INTO (PARTITION p1 VALUES LESS THAN (1000), PARTITION p2 VALUES LESS THAN (2000))")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATReorganizePartition {
			t.Errorf("Type = %d, want ATReorganizePartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 2 {
			t.Errorf("PartitionNames count = %d, want 2", len(cmd.PartitionNames))
		}
		if len(cmd.PartitionDefs) != 2 {
			t.Errorf("PartitionDefs count = %d, want 2", len(cmd.PartitionDefs))
		}
	})
}

func TestParseAlterTableExchangePartition(t *testing.T) {
	t.Run("exchange partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t EXCHANGE PARTITION p2 WITH TABLE t2")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATExchangePartition {
			t.Errorf("Type = %d, want ATExchangePartition", cmd.Type)
		}
		if cmd.Name != "p2" {
			t.Errorf("Name = %s, want p2", cmd.Name)
		}
		if cmd.ExchangeTable == nil || cmd.ExchangeTable.Name != "t2" {
			t.Errorf("ExchangeTable = %v, want t2", cmd.ExchangeTable)
		}
	})

	t.Run("exchange partition with validation", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t EXCHANGE PARTITION p2 WITH TABLE t2 WITH VALIDATION")
		cmd := stmt.Commands[0]
		if cmd.WithValidation == nil || !*cmd.WithValidation {
			t.Errorf("WithValidation = %v, want true", cmd.WithValidation)
		}
	})

	t.Run("exchange partition without validation", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t EXCHANGE PARTITION p2 WITH TABLE t2 WITHOUT VALIDATION")
		cmd := stmt.Commands[0]
		if cmd.WithValidation == nil || *cmd.WithValidation {
			t.Errorf("WithValidation = %v, want false", cmd.WithValidation)
		}
	})
}

func TestParseAlterTableTruncatePartition(t *testing.T) {
	t.Run("truncate specific partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t TRUNCATE PARTITION p1, p2")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATTruncatePartition {
			t.Errorf("Type = %d, want ATTruncatePartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 2 {
			t.Errorf("PartitionNames count = %d, want 2", len(cmd.PartitionNames))
		}
	})

	t.Run("truncate all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t TRUNCATE PARTITION ALL")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATTruncatePartition {
			t.Errorf("Type = %d, want ATTruncatePartition", cmd.Type)
		}
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableAnalyzePartition(t *testing.T) {
	t.Run("analyze partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ANALYZE PARTITION p1")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAnalyzePartition {
			t.Errorf("Type = %d, want ATAnalyzePartition", cmd.Type)
		}
	})

	t.Run("analyze all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ANALYZE PARTITION ALL")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAnalyzePartition {
			t.Errorf("Type = %d, want ATAnalyzePartition", cmd.Type)
		}
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableCheckPartition(t *testing.T) {
	t.Run("check partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t CHECK PARTITION p1, p2")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATCheckPartition {
			t.Errorf("Type = %d, want ATCheckPartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 2 {
			t.Errorf("PartitionNames count = %d, want 2", len(cmd.PartitionNames))
		}
	})

	t.Run("check all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t CHECK PARTITION ALL")
		cmd := stmt.Commands[0]
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableOptimizePartition(t *testing.T) {
	t.Run("optimize partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t OPTIMIZE PARTITION p1")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATOptimizePartition {
			t.Errorf("Type = %d, want ATOptimizePartition", cmd.Type)
		}
	})

	t.Run("optimize all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t OPTIMIZE PARTITION ALL")
		cmd := stmt.Commands[0]
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableRebuildPartition(t *testing.T) {
	t.Run("rebuild partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REBUILD PARTITION p1, p2")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRebuildPartition {
			t.Errorf("Type = %d, want ATRebuildPartition", cmd.Type)
		}
		if len(cmd.PartitionNames) != 2 {
			t.Errorf("PartitionNames count = %d, want 2", len(cmd.PartitionNames))
		}
	})

	t.Run("rebuild all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REBUILD PARTITION ALL")
		cmd := stmt.Commands[0]
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableRepairPartition(t *testing.T) {
	t.Run("repair partition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REPAIR PARTITION p1")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRepairPartition {
			t.Errorf("Type = %d, want ATRepairPartition", cmd.Type)
		}
	})

	t.Run("repair all partitions", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REPAIR PARTITION ALL")
		cmd := stmt.Commands[0]
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableDiscardImportPartition(t *testing.T) {
	t.Run("discard partition tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DISCARD PARTITION p1 TABLESPACE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDiscardPartitionTablespace {
			t.Errorf("Type = %d, want ATDiscardPartitionTablespace", cmd.Type)
		}
		if len(cmd.PartitionNames) != 1 || cmd.PartitionNames[0] != "p1" {
			t.Errorf("PartitionNames = %v, want [p1]", cmd.PartitionNames)
		}
	})

	t.Run("discard all partition tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DISCARD PARTITION ALL TABLESPACE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDiscardPartitionTablespace {
			t.Errorf("Type = %d, want ATDiscardPartitionTablespace", cmd.Type)
		}
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})

	t.Run("import partition tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t IMPORT PARTITION p1 TABLESPACE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATImportPartitionTablespace {
			t.Errorf("Type = %d, want ATImportPartitionTablespace", cmd.Type)
		}
	})

	t.Run("import all partition tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t IMPORT PARTITION ALL TABLESPACE")
		cmd := stmt.Commands[0]
		if !cmd.AllPartitions {
			t.Errorf("AllPartitions = false, want true")
		}
	})
}

func TestParseAlterTableRemovePartitioning(t *testing.T) {
	t.Run("remove partitioning", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REMOVE PARTITIONING")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATRemovePartitioning {
			t.Errorf("Type = %d, want ATRemovePartitioning", cmd.Type)
		}
	})
}

// ============================================================================
// Batch 58: ALTER TABLE visibility/misc
// ============================================================================

func TestParseAlterTableColumnVisible(t *testing.T) {
	t.Run("set column visible", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN name SET VISIBLE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterColumnVisible {
			t.Errorf("Type = %d, want ATAlterColumnVisible", cmd.Type)
		}
		if cmd.Name != "name" {
			t.Errorf("Name = %s, want name", cmd.Name)
		}
	})

	t.Run("set column invisible", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN name SET INVISIBLE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterColumnInvisible {
			t.Errorf("Type = %d, want ATAlterColumnInvisible", cmd.Type)
		}
	})
}

func TestParseAlterTableIndexVisible(t *testing.T) {
	t.Run("alter index visible", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER INDEX idx_name VISIBLE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterIndexVisible {
			t.Errorf("Type = %d, want ATAlterIndexVisible", cmd.Type)
		}
		if cmd.Name != "idx_name" {
			t.Errorf("Name = %s, want idx_name", cmd.Name)
		}
	})

	t.Run("alter index invisible", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER INDEX idx_name INVISIBLE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAlterIndexInvisible {
			t.Errorf("Type = %d, want ATAlterIndexInvisible", cmd.Type)
		}
	})
}

func TestParseAlterTableForce(t *testing.T) {
	t.Run("force", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t FORCE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATForce {
			t.Errorf("Type = %d, want ATForce", cmd.Type)
		}
	})

	t.Run("force with algorithm", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t FORCE, ALGORITHM=INPLACE")
		if len(stmt.Commands) != 2 {
			t.Fatalf("Commands count = %d, want 2", len(stmt.Commands))
		}
		if stmt.Commands[0].Type != ast.ATForce {
			t.Errorf("Commands[0].Type = %d, want ATForce", stmt.Commands[0].Type)
		}
		if stmt.Commands[1].Type != ast.ATAlgorithm {
			t.Errorf("Commands[1].Type = %d, want ATAlgorithm", stmt.Commands[1].Type)
		}
	})
}

func TestParseAlterTableOrderBy(t *testing.T) {
	t.Run("order by single column", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ORDER BY name")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATOrderBy {
			t.Errorf("Type = %d, want ATOrderBy", cmd.Type)
		}
		if len(cmd.OrderByItems) != 1 {
			t.Fatalf("OrderByItems count = %d, want 1", len(cmd.OrderByItems))
		}
	})

	t.Run("order by multiple columns", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ORDER BY name ASC, age DESC")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATOrderBy {
			t.Errorf("Type = %d, want ATOrderBy", cmd.Type)
		}
		if len(cmd.OrderByItems) != 2 {
			t.Fatalf("OrderByItems count = %d, want 2", len(cmd.OrderByItems))
		}
	})
}

func TestParseAlterTableEnableKeys(t *testing.T) {
	t.Run("enable keys", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ENABLE KEYS")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATEnableKeys {
			t.Errorf("Type = %d, want ATEnableKeys", cmd.Type)
		}
	})

	t.Run("disable keys", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DISABLE KEYS")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDisableKeys {
			t.Errorf("Type = %d, want ATDisableKeys", cmd.Type)
		}
	})
}

func TestParseAlterTableDiscardTablespace(t *testing.T) {
	t.Run("discard tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t DISCARD TABLESPACE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDiscardTablespace {
			t.Errorf("Type = %d, want ATDiscardTablespace", cmd.Type)
		}
	})

	t.Run("import tablespace", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t IMPORT TABLESPACE")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATImportTablespace {
			t.Errorf("Type = %d, want ATImportTablespace", cmd.Type)
		}
	})
}

func TestParseAlterTableValidation(t *testing.T) {
	t.Run("with validation", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN v SET INVISIBLE, WITH VALIDATION")
		if len(stmt.Commands) != 2 {
			t.Fatalf("Commands count = %d, want 2", len(stmt.Commands))
		}
		cmd := stmt.Commands[1]
		if cmd.Type != ast.ATWithValidation {
			t.Errorf("Type = %d, want ATWithValidation", cmd.Type)
		}
	})

	t.Run("without validation", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ALTER COLUMN v SET INVISIBLE, WITHOUT VALIDATION")
		cmd := stmt.Commands[1]
		if cmd.Type != ast.ATWithoutValidation {
			t.Errorf("Type = %d, want ATWithoutValidation", cmd.Type)
		}
	})
}

// ============================================================================
// Batch 9: CREATE INDEX
// ============================================================================

func parseCreateIndex(t *testing.T, input string) *ast.CreateIndexStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime lexer
	p.advance() // skip CREATE
	unique := false
	fulltext := false
	spatial := false
	if p.cur.Type == kwUNIQUE {
		p.advance()
		unique = true
	} else if p.cur.Type == kwFULLTEXT {
		p.advance()
		fulltext = true
	} else if p.cur.Type == kwSPATIAL {
		p.advance()
		spatial = true
	}
	stmt, err := p.parseCreateIndexStmt(unique, fulltext, spatial)
	if err != nil {
		t.Fatalf("parseCreateIndexStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseCreateIndex(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_name ON t (name)")
		if stmt.IndexName != "idx_name" {
			t.Errorf("IndexName = %s, want idx_name", stmt.IndexName)
		}
		if stmt.Table == nil || stmt.Table.Name != "t" {
			t.Errorf("Table = %v, want t", stmt.Table)
		}
		if len(stmt.Columns) != 1 {
			t.Fatalf("Columns count = %d, want 1", len(stmt.Columns))
		}
	})

	t.Run("multiple columns", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_ab ON t (a, b)")
		if len(stmt.Columns) != 2 {
			t.Fatalf("Columns count = %d, want 2", len(stmt.Columns))
		}
	})

	t.Run("prefix length", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_name ON t (name(10))")
		if stmt.Columns[0].Length != 10 {
			t.Errorf("Length = %d, want 10", stmt.Columns[0].Length)
		}
	})

	t.Run("desc order", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_id ON t (id DESC)")
		if !stmt.Columns[0].Desc {
			t.Errorf("Desc = false, want true")
		}
	})

	t.Run("using btree", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_id USING BTREE ON t (id)")
		if stmt.IndexType != "BTREE" {
			t.Errorf("IndexType = %s, want BTREE", stmt.IndexType)
		}
	})

	t.Run("schema qualified table", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_id ON mydb.t (id)")
		if stmt.Table.Schema != "mydb" {
			t.Errorf("Schema = %s, want mydb", stmt.Table.Schema)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX IF NOT EXISTS idx_id ON t (id)")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})

	t.Run("with comment", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_id ON t (id) COMMENT 'test index'")
		found := false
		for _, opt := range stmt.Options {
			if opt.Name == "COMMENT" {
				found = true
			}
		}
		if !found {
			t.Errorf("COMMENT option not found")
		}
	})

	t.Run("algorithm lock", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE INDEX idx_id ON t (id) ALGORITHM=INPLACE LOCK=NONE")
		if stmt.Algorithm != "INPLACE" {
			t.Errorf("Algorithm = %s, want INPLACE", stmt.Algorithm)
		}
		if stmt.Lock != "NONE" {
			t.Errorf("Lock = %s, want NONE", stmt.Lock)
		}
	})
}

func TestParseCreateUniqueIndex(t *testing.T) {
	stmt := parseCreateIndex(t, "CREATE UNIQUE INDEX idx_email ON users (email)")
	if !stmt.Unique {
		t.Errorf("Unique = false, want true")
	}
	if stmt.IndexName != "idx_email" {
		t.Errorf("IndexName = %s, want idx_email", stmt.IndexName)
	}
}

func TestParseCreateFulltextIndex(t *testing.T) {
	stmt := parseCreateIndex(t, "CREATE FULLTEXT INDEX ft_content ON articles (content)")
	if !stmt.Fulltext {
		t.Errorf("Fulltext = false, want true")
	}

	t.Run("spatial", func(t *testing.T) {
		stmt := parseCreateIndex(t, "CREATE SPATIAL INDEX sp_geo ON places (location)")
		if !stmt.Spatial {
			t.Errorf("Spatial = false, want true")
		}
	})
}

// ============================================================================
// Batch 13: SET, SHOW, USE, EXPLAIN
// ============================================================================

func parseSet(t *testing.T, input string) *ast.SetStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	node, err := p.parseSetStmt()
	if err != nil {
		t.Fatalf("parseSetStmt(%q) error: %v", input, err)
	}
	stmt, ok := node.(*ast.SetStmt)
	if !ok {
		t.Fatalf("parseSetStmt(%q) returned %T, want *ast.SetStmt", input, node)
	}
	return stmt
}

func parseShow(t *testing.T, input string) *ast.ShowStmt {
	t.Helper()
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	stmt, err := p.parseShowStmt()
	if err != nil {
		t.Fatalf("parseShowStmt(%q) error: %v", input, err)
	}
	return stmt
}

func TestParseSet(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		stmt := parseSet(t, "SET x = 1")
		if len(stmt.Assignments) != 1 {
			t.Fatalf("Assignments count = %d, want 1", len(stmt.Assignments))
		}
		if stmt.Assignments[0].Column.Column != "x" {
			t.Errorf("Column = %s, want x", stmt.Assignments[0].Column.Column)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		stmt := parseSet(t, "SET x = 1, y = 2")
		if len(stmt.Assignments) != 2 {
			t.Fatalf("Assignments count = %d, want 2", len(stmt.Assignments))
		}
	})

	t.Run("global scope", func(t *testing.T) {
		stmt := parseSet(t, "SET GLOBAL max_connections = 100")
		if stmt.Scope != "GLOBAL" {
			t.Errorf("Scope = %s, want GLOBAL", stmt.Scope)
		}
	})

	t.Run("session scope", func(t *testing.T) {
		stmt := parseSet(t, "SET SESSION sql_mode = 'STRICT_TRANS_TABLES'")
		if stmt.Scope != "SESSION" {
			t.Errorf("Scope = %s, want SESSION", stmt.Scope)
		}
	})

	t.Run("names", func(t *testing.T) {
		stmt := parseSet(t, "SET NAMES utf8mb4")
		if len(stmt.Assignments) < 1 {
			t.Fatalf("Assignments count = %d, want >= 1", len(stmt.Assignments))
		}
		if stmt.Assignments[0].Column.Column != "NAMES" {
			t.Errorf("Column = %s, want NAMES", stmt.Assignments[0].Column.Column)
		}
	})

	t.Run("names with collate", func(t *testing.T) {
		stmt := parseSet(t, "SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")
		if len(stmt.Assignments) != 2 {
			t.Fatalf("Assignments count = %d, want 2", len(stmt.Assignments))
		}
	})

	t.Run("character set", func(t *testing.T) {
		stmt := parseSet(t, "SET CHARACTER SET utf8mb4")
		if len(stmt.Assignments) < 1 {
			t.Fatalf("Assignments count = %d, want >= 1", len(stmt.Assignments))
		}
	})
}

func TestParseShow(t *testing.T) {
	t.Run("databases", func(t *testing.T) {
		stmt := parseShow(t, "SHOW DATABASES")
		if stmt.Type != "DATABASES" {
			t.Errorf("Type = %s, want DATABASES", stmt.Type)
		}
	})

	t.Run("tables", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLES")
		if stmt.Type != "TABLES" {
			t.Errorf("Type = %s, want TABLES", stmt.Type)
		}
	})

	t.Run("tables from db", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLES FROM mydb")
		if stmt.Type != "TABLES" {
			t.Errorf("Type = %s, want TABLES", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mydb" {
			t.Errorf("From = %v, want mydb", stmt.From)
		}
	})

	t.Run("tables like", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLES LIKE 'user%'")
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("columns from table", func(t *testing.T) {
		stmt := parseShow(t, "SHOW COLUMNS FROM users")
		if stmt.Type != "COLUMNS" {
			t.Errorf("Type = %s, want COLUMNS", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "users" {
			t.Errorf("From = %v, want users", stmt.From)
		}
	})

	t.Run("create table", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CREATE TABLE users")
		if stmt.Type != "CREATE TABLE" {
			t.Errorf("Type = %s, want CREATE TABLE", stmt.Type)
		}
	})

	t.Run("variables", func(t *testing.T) {
		stmt := parseShow(t, "SHOW VARIABLES")
		if stmt.Type != "VARIABLES" {
			t.Errorf("Type = %s, want VARIABLES", stmt.Type)
		}
	})

	t.Run("global variables", func(t *testing.T) {
		stmt := parseShow(t, "SHOW GLOBAL VARIABLES")
		if stmt.Type != "GLOBAL VARIABLES" {
			t.Errorf("Type = %s, want GLOBAL VARIABLES", stmt.Type)
		}
	})

	t.Run("status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW STATUS")
		if stmt.Type != "STATUS" {
			t.Errorf("Type = %s, want STATUS", stmt.Type)
		}
	})

	t.Run("warnings", func(t *testing.T) {
		stmt := parseShow(t, "SHOW WARNINGS")
		if stmt.Type != "WARNINGS" {
			t.Errorf("Type = %s, want WARNINGS", stmt.Type)
		}
	})

	t.Run("errors", func(t *testing.T) {
		stmt := parseShow(t, "SHOW ERRORS")
		if stmt.Type != "ERRORS" {
			t.Errorf("Type = %s, want ERRORS", stmt.Type)
		}
	})

	t.Run("index from table", func(t *testing.T) {
		stmt := parseShow(t, "SHOW INDEX FROM users")
		if stmt.Type != "INDEX" {
			t.Errorf("Type = %s, want INDEX", stmt.Type)
		}
	})

	// Batch 32: Extended SHOW variants
	t.Run("engines", func(t *testing.T) {
		stmt := parseShow(t, "SHOW ENGINES")
		if stmt.Type != "ENGINES" {
			t.Errorf("Type = %s, want ENGINES", stmt.Type)
		}
	})

	t.Run("plugins", func(t *testing.T) {
		stmt := parseShow(t, "SHOW PLUGINS")
		if stmt.Type != "PLUGINS" {
			t.Errorf("Type = %s, want PLUGINS", stmt.Type)
		}
	})

	t.Run("master status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW MASTER STATUS")
		if stmt.Type != "MASTER STATUS" {
			t.Errorf("Type = %s, want MASTER STATUS", stmt.Type)
		}
	})

	t.Run("slave status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW SLAVE STATUS")
		if stmt.Type != "SLAVE STATUS" {
			t.Errorf("Type = %s, want SLAVE STATUS", stmt.Type)
		}
	})

	t.Run("replica status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW REPLICA STATUS")
		if stmt.Type != "REPLICA STATUS" {
			t.Errorf("Type = %s, want REPLICA STATUS", stmt.Type)
		}
	})

	t.Run("binary logs", func(t *testing.T) {
		stmt := parseShow(t, "SHOW BINARY LOGS")
		if stmt.Type != "BINARY LOGS" {
			t.Errorf("Type = %s, want BINARY LOGS", stmt.Type)
		}
	})

	t.Run("binlog events", func(t *testing.T) {
		stmt := parseShow(t, "SHOW BINLOG EVENTS")
		if stmt.Type != "BINLOG EVENTS" {
			t.Errorf("Type = %s, want BINLOG EVENTS", stmt.Type)
		}
	})

	t.Run("binlog events in", func(t *testing.T) {
		stmt := parseShow(t, "SHOW BINLOG EVENTS IN 'mysql-bin.000001'")
		if stmt.Type != "BINLOG EVENTS" {
			t.Errorf("Type = %s, want BINLOG EVENTS", stmt.Type)
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil, expected log name")
		}
	})

	t.Run("binlog events limit", func(t *testing.T) {
		stmt := parseShow(t, "SHOW BINLOG EVENTS LIMIT 10")
		if stmt.Type != "BINLOG EVENTS" {
			t.Errorf("Type = %s, want BINLOG EVENTS", stmt.Type)
		}
	})

	t.Run("table status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLE STATUS")
		if stmt.Type != "TABLE STATUS" {
			t.Errorf("Type = %s, want TABLE STATUS", stmt.Type)
		}
	})

	t.Run("table status from db", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLE STATUS FROM mydb")
		if stmt.Type != "TABLE STATUS" {
			t.Errorf("Type = %s, want TABLE STATUS", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mydb" {
			t.Errorf("From = %v, want mydb", stmt.From)
		}
	})

	t.Run("table status like", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TABLE STATUS LIKE 'user%'")
		if stmt.Type != "TABLE STATUS" {
			t.Errorf("Type = %s, want TABLE STATUS", stmt.Type)
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("triggers", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TRIGGERS")
		if stmt.Type != "TRIGGERS" {
			t.Errorf("Type = %s, want TRIGGERS", stmt.Type)
		}
	})

	t.Run("triggers from db", func(t *testing.T) {
		stmt := parseShow(t, "SHOW TRIGGERS FROM mydb")
		if stmt.Type != "TRIGGERS" {
			t.Errorf("Type = %s, want TRIGGERS", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mydb" {
			t.Errorf("From = %v, want mydb", stmt.From)
		}
	})

	t.Run("events", func(t *testing.T) {
		stmt := parseShow(t, "SHOW EVENTS")
		if stmt.Type != "EVENTS" {
			t.Errorf("Type = %s, want EVENTS", stmt.Type)
		}
	})

	t.Run("events from db like", func(t *testing.T) {
		stmt := parseShow(t, "SHOW EVENTS FROM mydb LIKE 'ev%'")
		if stmt.Type != "EVENTS" {
			t.Errorf("Type = %s, want EVENTS", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mydb" {
			t.Errorf("From = %v, want mydb", stmt.From)
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("procedure status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW PROCEDURE STATUS")
		if stmt.Type != "PROCEDURE STATUS" {
			t.Errorf("Type = %s, want PROCEDURE STATUS", stmt.Type)
		}
	})

	t.Run("procedure status like", func(t *testing.T) {
		stmt := parseShow(t, "SHOW PROCEDURE STATUS LIKE 'proc%'")
		if stmt.Type != "PROCEDURE STATUS" {
			t.Errorf("Type = %s, want PROCEDURE STATUS", stmt.Type)
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("function status", func(t *testing.T) {
		stmt := parseShow(t, "SHOW FUNCTION STATUS")
		if stmt.Type != "FUNCTION STATUS" {
			t.Errorf("Type = %s, want FUNCTION STATUS", stmt.Type)
		}
	})

	t.Run("create procedure", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CREATE PROCEDURE myproc")
		if stmt.Type != "CREATE PROCEDURE" {
			t.Errorf("Type = %s, want CREATE PROCEDURE", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "myproc" {
			t.Errorf("From = %v, want myproc", stmt.From)
		}
	})

	t.Run("create function", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CREATE FUNCTION myfunc")
		if stmt.Type != "CREATE FUNCTION" {
			t.Errorf("Type = %s, want CREATE FUNCTION", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "myfunc" {
			t.Errorf("From = %v, want myfunc", stmt.From)
		}
	})

	t.Run("create trigger", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CREATE TRIGGER mytrigger")
		if stmt.Type != "CREATE TRIGGER" {
			t.Errorf("Type = %s, want CREATE TRIGGER", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mytrigger" {
			t.Errorf("From = %v, want mytrigger", stmt.From)
		}
	})

	t.Run("create event", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CREATE EVENT myevent")
		if stmt.Type != "CREATE EVENT" {
			t.Errorf("Type = %s, want CREATE EVENT", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "myevent" {
			t.Errorf("From = %v, want myevent", stmt.From)
		}
	})

	t.Run("open tables", func(t *testing.T) {
		stmt := parseShow(t, "SHOW OPEN TABLES")
		if stmt.Type != "OPEN TABLES" {
			t.Errorf("Type = %s, want OPEN TABLES", stmt.Type)
		}
	})

	t.Run("open tables from db", func(t *testing.T) {
		stmt := parseShow(t, "SHOW OPEN TABLES FROM mydb")
		if stmt.Type != "OPEN TABLES" {
			t.Errorf("Type = %s, want OPEN TABLES", stmt.Type)
		}
		if stmt.From == nil || stmt.From.Name != "mydb" {
			t.Errorf("From = %v, want mydb", stmt.From)
		}
	})

	t.Run("privileges", func(t *testing.T) {
		stmt := parseShow(t, "SHOW PRIVILEGES")
		if stmt.Type != "PRIVILEGES" {
			t.Errorf("Type = %s, want PRIVILEGES", stmt.Type)
		}
	})

	t.Run("profiles", func(t *testing.T) {
		stmt := parseShow(t, "SHOW PROFILES")
		if stmt.Type != "PROFILES" {
			t.Errorf("Type = %s, want PROFILES", stmt.Type)
		}
	})

	t.Run("relaylog events", func(t *testing.T) {
		stmt := parseShow(t, "SHOW RELAYLOG EVENTS")
		if stmt.Type != "RELAYLOG EVENTS" {
			t.Errorf("Type = %s, want RELAYLOG EVENTS", stmt.Type)
		}
	})

	t.Run("character set", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CHARACTER SET")
		if stmt.Type != "CHARACTER SET" {
			t.Errorf("Type = %s, want CHARACTER SET", stmt.Type)
		}
	})

	t.Run("character set like", func(t *testing.T) {
		stmt := parseShow(t, "SHOW CHARACTER SET LIKE 'utf8%'")
		if stmt.Type != "CHARACTER SET" {
			t.Errorf("Type = %s, want CHARACTER SET", stmt.Type)
		}
		if stmt.Like == nil {
			t.Fatal("Like is nil")
		}
	})

	t.Run("collation", func(t *testing.T) {
		stmt := parseShow(t, "SHOW COLLATION")
		if stmt.Type != "COLLATION" {
			t.Errorf("Type = %s, want COLLATION", stmt.Type)
		}
	})

	t.Run("collation where", func(t *testing.T) {
		stmt := parseShow(t, "SHOW COLLATION WHERE Charset = 'utf8mb4'")
		if stmt.Type != "COLLATION" {
			t.Errorf("Type = %s, want COLLATION", stmt.Type)
		}
		if stmt.Where == nil {
			t.Fatal("Where is nil")
		}
	})
}

func TestParseUse(t *testing.T) {
	p := &Parser{lexer: NewLexer("USE mydb")}
	p.advance()
	stmt, err := p.parseUseStmt()
	if err != nil {
		t.Fatalf("parseUseStmt error: %v", err)
	}
	if stmt.Database != "mydb" {
		t.Errorf("Database = %s, want mydb", stmt.Database)
	}
}

func TestParseExplain(t *testing.T) {
	t.Run("explain select", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN SELECT 1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Stmt == nil {
			t.Fatal("Stmt is nil")
		}
		if _, ok := stmt.Stmt.(*ast.SelectStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.SelectStmt", stmt.Stmt)
		}
	})

	t.Run("explain format json", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN FORMAT=JSON SELECT 1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Format != "JSON" {
			t.Errorf("Format = %s, want JSON", stmt.Format)
		}
	})

	t.Run("describe table", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("DESCRIBE users")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Stmt == nil {
			t.Fatal("Stmt is nil")
		}
	})

	t.Run("describe table with column", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("DESCRIBE users name")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Stmt == nil {
			t.Fatal("Stmt is nil")
		}
		show, ok := stmt.Stmt.(*ast.ShowStmt)
		if !ok {
			t.Fatalf("Stmt type = %T, want *ast.ShowStmt", stmt.Stmt)
		}
		if show.Like == nil {
			t.Error("Like is nil, want column filter")
		}
	})

	t.Run("explain analyze select", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN ANALYZE SELECT * FROM t1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if !stmt.Analyze {
			t.Error("Analyze = false, want true")
		}
		if _, ok := stmt.Stmt.(*ast.SelectStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.SelectStmt", stmt.Stmt)
		}
	})

	t.Run("explain extended select", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN EXTENDED SELECT * FROM t1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if !stmt.Extended {
			t.Error("Extended = false, want true")
		}
	})

	t.Run("explain partitions select", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN PARTITIONS SELECT * FROM t1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if !stmt.Partitions {
			t.Error("Partitions = false, want true")
		}
	})

	t.Run("explain format tree select", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN FORMAT=TREE SELECT * FROM t1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Format != "TREE" {
			t.Errorf("Format = %s, want TREE", stmt.Format)
		}
	})

	t.Run("explain insert", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN INSERT INTO t1 (a) VALUES (1)")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if _, ok := stmt.Stmt.(*ast.InsertStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.InsertStmt", stmt.Stmt)
		}
	})

	t.Run("explain update", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN UPDATE t1 SET a = 1 WHERE id = 2")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if _, ok := stmt.Stmt.(*ast.UpdateStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.UpdateStmt", stmt.Stmt)
		}
	})

	t.Run("explain delete", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN DELETE FROM t1 WHERE id = 1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if _, ok := stmt.Stmt.(*ast.DeleteStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.DeleteStmt", stmt.Stmt)
		}
	})

	t.Run("explain replace", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN REPLACE INTO t1 (a) VALUES (1)")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if _, ok := stmt.Stmt.(*ast.InsertStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.InsertStmt", stmt.Stmt)
		}
	})

	t.Run("explain format json insert", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN FORMAT=JSON INSERT INTO t1 (a) VALUES (1)")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if stmt.Format != "JSON" {
			t.Errorf("Format = %s, want JSON", stmt.Format)
		}
		if _, ok := stmt.Stmt.(*ast.InsertStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.InsertStmt", stmt.Stmt)
		}
	})

	t.Run("explain table name", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXPLAIN t1")}
		p.advance()
		stmt, err := p.parseExplainStmt()
		if err != nil {
			t.Fatalf("parseExplainStmt error: %v", err)
		}
		if _, ok := stmt.Stmt.(*ast.ShowStmt); !ok {
			t.Errorf("Stmt type = %T, want *ast.ShowStmt", stmt.Stmt)
		}
	})
}

// ============================================================================
// Batch 15: GRANT, REVOKE, CREATE/DROP USER
// ============================================================================

func TestParseGrant(t *testing.T) {
	t.Run("all privileges", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT ALL PRIVILEGES ON *.* TO root")}
		p.advance()
		node, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		stmt := node.(*ast.GrantStmt)
		if !stmt.AllPriv {
			t.Errorf("AllPriv = false, want true")
		}
		if stmt.On == nil || stmt.On.Name == nil {
			t.Fatal("On is nil")
		}
	})

	t.Run("specific privileges", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT SELECT, INSERT ON mydb.* TO app_user")}
		p.advance()
		node, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		stmt := node.(*ast.GrantStmt)
		if len(stmt.Privileges) != 2 {
			t.Fatalf("Privileges count = %d, want 2", len(stmt.Privileges))
		}
	})

	t.Run("with grant option", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT ALL ON *.* TO admin WITH GRANT OPTION")}
		p.advance()
		node, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		stmt := node.(*ast.GrantStmt)
		if !stmt.WithGrant {
			t.Errorf("WithGrant = false, want true")
		}
	})
}

func TestParseRevoke(t *testing.T) {
	p := &Parser{lexer: NewLexer("REVOKE SELECT ON mydb.* FROM app_user")}
	p.advance()
	node, err := p.parseRevokeStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	stmt := node.(*ast.RevokeStmt)
	if len(stmt.Privileges) != 1 {
		t.Fatalf("Privileges count = %d, want 1", len(stmt.Privileges))
	}
	if len(stmt.From) != 1 {
		t.Fatalf("From count = %d, want 1", len(stmt.From))
	}
}

func TestParseCreateUser(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE USER testuser")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateUserStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(stmt.Users) != 1 {
			t.Fatalf("Users count = %d, want 1", len(stmt.Users))
		}
		if stmt.Users[0].Name != "testuser" {
			t.Errorf("Name = %s, want testuser", stmt.Users[0].Name)
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE USER IF NOT EXISTS testuser")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateUserStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})

	t.Run("identified by", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE USER testuser IDENTIFIED BY 'secret'")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateUserStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Users[0].Password != "secret" {
			t.Errorf("Password = %s, want secret", stmt.Users[0].Password)
		}
	})
}

func TestParseDropUser(t *testing.T) {
	p := &Parser{lexer: NewLexer("DROP USER testuser")}
	p.advance()
	p.advance() // skip DROP
	stmt, err := p.parseDropUserStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(stmt.Users) != 1 {
		t.Fatalf("Users count = %d, want 1", len(stmt.Users))
	}
}

// ============================================================================
// Batch 16: CREATE FUNCTION/PROCEDURE
// ============================================================================

func TestParseCreateFunction(t *testing.T) {
	t.Run("simple function", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE FUNCTION myfunc(a INT) RETURNS INT RETURN a + 1")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateFunctionStmt(false)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Name == nil || stmt.Name.Name != "myfunc" {
			t.Errorf("Name = %v, want myfunc", stmt.Name)
		}
		if len(stmt.Params) != 1 {
			t.Fatalf("Params count = %d, want 1", len(stmt.Params))
		}
		if stmt.Params[0].Name != "a" {
			t.Errorf("Param name = %s, want a", stmt.Params[0].Name)
		}
		if stmt.Returns == nil {
			t.Fatal("Returns is nil")
		}
	})

	t.Run("function no params", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE FUNCTION myfunc() RETURNS INT RETURN 1")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateFunctionStmt(false)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(stmt.Params) != 0 {
			t.Errorf("Params count = %d, want 0", len(stmt.Params))
		}
	})
}

func TestParseCreateProcedure(t *testing.T) {
	t.Run("simple procedure", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE PROCEDURE myproc(IN x INT, OUT y INT) BEGIN SELECT x INTO y; END")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateFunctionStmt(true)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.IsProcedure {
			t.Errorf("IsProcedure = false, want true")
		}
		if len(stmt.Params) != 2 {
			t.Fatalf("Params count = %d, want 2", len(stmt.Params))
		}
		if stmt.Params[0].Direction != "IN" {
			t.Errorf("Param[0].Direction = %s, want IN", stmt.Params[0].Direction)
		}
		if stmt.Params[1].Direction != "OUT" {
			t.Errorf("Param[1].Direction = %s, want OUT", stmt.Params[1].Direction)
		}
	})
}

// TestParseRoutineCompoundBody verifies that routine/trigger/event body scanners
// correctly balance inner compound statements (IF/CASE/WHILE/LOOP/REPEAT) against
// their END so that END IF, END CASE, etc. do not prematurely close the outer
// BEGIN...END wrapper.
//
// Each case specifies the SQL plus the expected number of top-level statements
// and the expected body-text suffix. The suffix check guarantees the body
// scanner did not terminate early at an inner END (which would leave e.g.
// " IF" trailing outside the captured body).
func TestParseRoutineCompoundBody(t *testing.T) {
	cases := []struct {
		name       string
		sql        string
		wantStmts  int
		bodySuffix string // expected suffix of the captured routine body (raw text)
	}{
		{
			name: "procedure with IF...END IF",
			sql: `CREATE PROCEDURE p1(IN uid INT)
BEGIN
    IF uid = 0 THEN
        SET @statement = 'x';
    END IF;
    SET @user_id = uid;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with DELIMITER and ELSEIF/ELSE chain",
			sql: `DELIMITER ;;
CREATE PROCEDURE p2(IN uid INT)
BEGIN
    IF uid = 0 THEN
        SET @a = 1;
    ELSEIF uid = 1 THEN
        SET @a = 2;
    ELSE
        SET @a = 3;
    END IF;
    SET @b = uid;
END ;;
DELIMITER ;`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with WHILE...END WHILE",
			sql: `CREATE PROCEDURE p3()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 10 DO
        SET i = i + 1;
    END WHILE;
    SET @done = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with simple CASE statement",
			sql: `CREATE PROCEDURE p4(IN x INT)
BEGIN
    CASE x
        WHEN 1 THEN SET @r = 'one';
        WHEN 2 THEN SET @r = 'two';
        ELSE SET @r = 'other';
    END CASE;
    SET @y = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with searched CASE statement",
			sql: `CREATE PROCEDURE p4s(IN x INT)
BEGIN
    CASE
        WHEN x < 0 THEN SET @r = 'neg';
        WHEN x = 0 THEN SET @r = 'zero';
        ELSE SET @r = 'pos';
    END CASE;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with labeled LOOP + LEAVE + ITERATE",
			sql: `CREATE PROCEDURE p5()
BEGIN
    DECLARE i INT DEFAULT 0;
    outer_loop: LOOP
        SET i = i + 1;
        IF i = 3 THEN
            ITERATE outer_loop;
        END IF;
        IF i >= 5 THEN
            LEAVE outer_loop;
        END IF;
    END LOOP outer_loop;
    SET @done = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with REPEAT...UNTIL...END REPEAT",
			sql: `CREATE PROCEDURE p6()
BEGIN
    DECLARE i INT DEFAULT 0;
    REPEAT
        SET i = i + 1;
    UNTIL i >= 10
    END REPEAT;
    SET @done = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "function with CASE expression inside BEGIN...END",
			sql: `CREATE FUNCTION f1(x INT) RETURNS VARCHAR(20)
BEGIN
    DECLARE r VARCHAR(20);
    SET r = CASE x WHEN 1 THEN 'one' ELSE 'other' END;
    RETURN r;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with nested BEGIN inside IF",
			sql: `CREATE PROCEDURE p7(IN x INT)
BEGIN
    IF x > 0 THEN
        BEGIN
            DECLARE v INT;
            SET v = x * 2;
            SET @r = v;
        END;
    END IF;
    SET @done = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "labeled BEGIN...END with matching end-label",
			sql: `CREATE PROCEDURE p8()
myblock: BEGIN
    SET @a = 1;
    IF @a = 1 THEN
        SET @b = 2;
    END IF;
END myblock`,
			wantStmts:  1,
			bodySuffix: "END myblock",
		},
		{
			name: "labeled WHILE with end-label",
			sql: `CREATE FUNCTION CalcIncome(starting_value INT) RETURNS INT DETERMINISTIC
BEGIN
   DECLARE income INT;
   SET income = 0;
   label1: WHILE income <= 3000 DO
     SET income = income + starting_value;
   END WHILE label1;
   RETURN income;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with cursor + DECLARE HANDLER + LOOP",
			sql: `CREATE PROCEDURE p9()
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE v INT;
    DECLARE cur CURSOR FOR SELECT id FROM t;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO v;
        IF done THEN
            LEAVE read_loop;
        END IF;
        SET @last = v;
    END LOOP read_loop;
    CLOSE cur;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with named condition + handler for SQLSTATE",
			sql: `CREATE PROCEDURE p10()
BEGIN
    DECLARE dup_key CONDITION FOR SQLSTATE '23000';
    DECLARE EXIT HANDLER FOR dup_key
    BEGIN
        SET @err = 'duplicate';
    END;
    INSERT INTO t(id) VALUES (1);
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with SIGNAL and RESIGNAL",
			sql: `CREATE PROCEDURE p11(IN x INT)
BEGIN
    IF x < 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'negative';
    END IF;
    BEGIN
        DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            RESIGNAL;
        END;
        SET @y = x * 2;
    END;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with DROP TABLE IF EXISTS inside body",
			sql: `CREATE PROCEDURE p12()
BEGIN
    DROP TABLE IF EXISTS tmp;
    CREATE TABLE IF NOT EXISTS tmp (id INT);
    SET @done = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with IF EXISTS (subquery) THEN flow control",
			sql: `CREATE PROCEDURE p13()
BEGIN
    IF EXISTS (SELECT 1 FROM t WHERE id = 1) THEN
        DELETE FROM t WHERE id = 1;
    END IF;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with end-label containing keyword tokens in string",
			sql: `CREATE PROCEDURE p14()
BEGIN
    SET @msg = 'END IF; END WHILE; END';
    IF LENGTH(@msg) > 0 THEN
        SET @ok = 1;
    END IF;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "procedure with comment containing fake END IF",
			sql: `CREATE PROCEDURE p15()
BEGIN
    /* END IF; END WHILE; END */
    -- END IF
    SET @x = 1;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "deeply nested: IF inside WHILE inside IF inside BEGIN",
			sql: `CREATE PROCEDURE p16()
BEGIN
    DECLARE i INT DEFAULT 0;
    IF @cond = 1 THEN
        WHILE i < 10 DO
            IF i % 2 = 0 THEN
                CASE i
                    WHEN 0 THEN SET @a = 'zero';
                    ELSE SET @a = 'even';
                END CASE;
            END IF;
            SET i = i + 1;
        END WHILE;
    END IF;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "trigger with IF...END IF",
			sql: `CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW
BEGIN
    IF NEW.x IS NULL THEN
        SET NEW.x = 0;
    END IF;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "trigger with WHILE loop",
			sql: `CREATE TRIGGER trg_w BEFORE UPDATE ON t FOR EACH ROW
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 3 DO
        SET i = i + 1;
    END WHILE;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "event with IF EXISTS (subquery)",
			sql: `CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO
BEGIN
    IF EXISTS (SELECT 1 FROM t) THEN
        DELETE FROM t;
    END IF;
END`,
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "multiple procedures with DELIMITER in one script",
			sql: `DELIMITER ;;
CREATE PROCEDURE p_a()
BEGIN
    IF @x = 1 THEN SET @y = 2; END IF;
END ;;
CREATE PROCEDURE p_b()
BEGIN
    WHILE @i < 3 DO SET @i = @i + 1; END WHILE;
END ;;
DELIMITER ;`,
			wantStmts:  2,
			bodySuffix: "END",
		},
		{
			name: "procedure with DEFINER and characteristics",
			sql: "CREATE DEFINER=`root`@`%` PROCEDURE p_def(IN x INT)\n" +
				"    READS SQL DATA\n" +
				"    DETERMINISTIC\n" +
				"    SQL SECURITY DEFINER\n" +
				"    COMMENT 'test'\n" +
				"BEGIN\n" +
				"    IF x > 0 THEN SET @y = 1; END IF;\n" +
				"END",
			wantStmts:  1,
			bodySuffix: "END",
		},
		{
			name: "function RETURNS with single RETURN body (no BEGIN)",
			sql:        `CREATE FUNCTION f_simple(a INT) RETURNS INT RETURN a + 1`,
			wantStmts:  1,
			bodySuffix: "a + 1",
		},
		{
			name: "procedure body is a single IF (no BEGIN wrapper)",
			sql: `CREATE PROCEDURE p_single(IN x INT)
IF x > 0 THEN
    SET @y = 1;
ELSE
    SET @y = 0;
END IF`,
			wantStmts:  1,
			bodySuffix: "END IF",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if got := len(list.Items); got != tc.wantStmts {
				t.Fatalf("got %d top-level stmts, want %d", got, tc.wantStmts)
			}
			// Verify the first routine/trigger/event body ends at the outermost END.
			body := firstRoutineBody(t, list)
			trimmed := strings.TrimSpace(body)
			if !strings.HasSuffix(trimmed, tc.bodySuffix) {
				t.Errorf("body suffix mismatch:\n  got:  %q\n  want suffix: %q", trimmed, tc.bodySuffix)
			}
		})
	}
}

// firstRoutineBody returns the Body field of the first CreateFunctionStmt,
// CreateTriggerStmt, or CreateEventStmt in the statement list.
func firstRoutineBody(t *testing.T, list *ast.List) string {
	t.Helper()
	for _, n := range list.Items {
		switch s := n.(type) {
		case *ast.CreateFunctionStmt:
			return s.Body
		case *ast.CreateTriggerStmt:
			return s.Body
		case *ast.CreateEventStmt:
			return s.Body
		}
	}
	t.Fatalf("no routine/trigger/event statement in list")
	return ""
}

// ============================================================================
// Batch 17: CREATE TRIGGER, CREATE EVENT
// ============================================================================

func TestParseCreateTrigger(t *testing.T) {
	t.Run("simple trigger", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET @a = 1")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateTriggerStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Name != "trg" {
			t.Errorf("Name = %s, want trg", stmt.Name)
		}
		if stmt.Timing != "BEFORE" {
			t.Errorf("Timing = %s, want BEFORE", stmt.Timing)
		}
		if stmt.Event != "INSERT" {
			t.Errorf("Event = %s, want INSERT", stmt.Event)
		}
		if stmt.Table == nil || stmt.Table.Name != "t" {
			t.Errorf("Table = %v, want t", stmt.Table)
		}
	})

	t.Run("after update", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE TRIGGER trg AFTER UPDATE ON t FOR EACH ROW SET @a = 1")}
		p.advance()
		p.advance()
		stmt, err := p.parseCreateTriggerStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Timing != "AFTER" {
			t.Errorf("Timing = %s, want AFTER", stmt.Timing)
		}
		if stmt.Event != "UPDATE" {
			t.Errorf("Event = %s, want UPDATE", stmt.Event)
		}
	})
}

func TestParseCreateEvent(t *testing.T) {
	t.Run("at schedule", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE EVENT myevent ON SCHEDULE AT '2025-01-01 00:00:00' DO SELECT 1")}
		p.advance()
		p.advance() // skip CREATE
		stmt, err := p.parseCreateEventStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Name != "myevent" {
			t.Errorf("Name = %s, want myevent", stmt.Name)
		}
		if stmt.Schedule == nil {
			t.Fatal("Schedule is nil")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("CREATE EVENT IF NOT EXISTS myevent ON SCHEDULE AT '2025-01-01' DO SELECT 1")}
		p.advance()
		p.advance()
		stmt, err := p.parseCreateEventStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})
}

// ============================================================================
// Batch 18: LOAD DATA
// ============================================================================

func TestParseLoadData(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Infile != "/tmp/data.csv" {
			t.Errorf("Infile = %s, want /tmp/data.csv", stmt.Infile)
		}
		if stmt.Table == nil || stmt.Table.Name != "t" {
			t.Errorf("Table = %v, want t", stmt.Table)
		}
	})

	t.Run("local", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA LOCAL INFILE '/tmp/data.csv' INTO TABLE t")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.Local {
			t.Errorf("Local = false, want true")
		}
	})

	t.Run("replace", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' REPLACE INTO TABLE t")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.Replace {
			t.Errorf("Replace = false, want true")
		}
	})

	t.Run("fields terminated", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t FIELDS TERMINATED BY ','")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.FieldsTerminatedBy != "," {
			t.Errorf("FieldsTerminatedBy = %s, want ,", stmt.FieldsTerminatedBy)
		}
	})
}

// ============================================================================
// Batch 63: LOAD DATA LINES STARTING BY fix
// ============================================================================

func TestParseLoadDataLinesStartingBy(t *testing.T) {
	t.Run("lines starting by", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t LINES STARTING BY '>>'")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.LinesStartingBy != ">>" {
			t.Errorf("LinesStartingBy = %q, want \">>\"", stmt.LinesStartingBy)
		}
	})

	t.Run("lines starting by and terminated by", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t LINES STARTING BY '>' TERMINATED BY '\\n'")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.LinesStartingBy != ">" {
			t.Errorf("LinesStartingBy = %q, want \">\"", stmt.LinesStartingBy)
		}
		if stmt.LinesTerminatedBy == "" {
			t.Errorf("LinesTerminatedBy is empty, want non-empty")
		}
	})

	t.Run("lines terminated by only", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t LINES TERMINATED BY '\\r\\n'")}
		p.advance()
		start := p.pos()
		p.advance() // consume LOAD
		stmt, err := p.parseLoadDataStmt(start)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.LinesStartingBy != "" {
			t.Errorf("LinesStartingBy = %q, want empty", stmt.LinesStartingBy)
		}
		if stmt.LinesTerminatedBy == "" {
			t.Errorf("LinesTerminatedBy is empty, want non-empty")
		}
	})
}

// ============================================================================
// Batch 19: PREPARE, EXECUTE, DEALLOCATE
// ============================================================================

func TestParsePrepare(t *testing.T) {
	p := &Parser{lexer: NewLexer("PREPARE stmt FROM 'SELECT * FROM t WHERE id = ?'")}
	p.advance()
	stmt, err := p.parsePrepareStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if stmt.Name != "stmt" {
		t.Errorf("Name = %s, want stmt", stmt.Name)
	}
	if stmt.Stmt != "SELECT * FROM t WHERE id = ?" {
		t.Errorf("Stmt = %s, want 'SELECT * FROM t WHERE id = ?'", stmt.Stmt)
	}
}

func TestParseExecute(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXECUTE stmt")}
		p.advance()
		stmt, err := p.parseExecuteStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.Name != "stmt" {
			t.Errorf("Name = %s, want stmt", stmt.Name)
		}
	})

	t.Run("with using", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("EXECUTE stmt USING @a, @b")}
		p.advance()
		stmt, err := p.parseExecuteStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(stmt.Params) != 2 {
			t.Fatalf("Params count = %d, want 2", len(stmt.Params))
		}
	})
}

func TestParseDeallocate(t *testing.T) {
	p := &Parser{lexer: NewLexer("DEALLOCATE PREPARE stmt")}
	p.advance()
	stmt, err := p.parseDeallocateStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if stmt.Name != "stmt" {
		t.Errorf("Name = %s, want stmt", stmt.Name)
	}
}

func TestParseDropPrepare(t *testing.T) {
	// DROP PREPARE is an alias for DEALLOCATE PREPARE
	ParseAndCheck(t, "DROP PREPARE stmt1")
}

// ============================================================================
// Batch 20: ANALYZE, OPTIMIZE, KILL, DO
// ============================================================================

func TestParseAnalyze(t *testing.T) {
	p := &Parser{lexer: NewLexer("ANALYZE TABLE t1, t2")}
	p.advance()
	stmt, err := p.parseAnalyzeTableStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(stmt.Tables) != 2 {
		t.Fatalf("Tables count = %d, want 2", len(stmt.Tables))
	}
}

func TestParseOptimize(t *testing.T) {
	p := &Parser{lexer: NewLexer("OPTIMIZE TABLE t1")}
	p.advance()
	stmt, err := p.parseOptimizeTableStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(stmt.Tables) != 1 {
		t.Fatalf("Tables count = %d, want 1", len(stmt.Tables))
	}
}

// Old TestParseKill and TestParseDo moved to batch 82 section below.

// ─── Batch 21: stmt_dispatch integration tests ────────────────────────────

func TestParseEmptyInput(t *testing.T) {
	result := ParseAndCheck(t, "")
	if len(result.Items) != 0 {
		t.Fatalf("expected 0 items for empty input, got %d", len(result.Items))
	}

	result = ParseAndCheck(t, "   ")
	if len(result.Items) != 0 {
		t.Fatalf("expected 0 items for whitespace input, got %d", len(result.Items))
	}

	result = ParseAndCheck(t, ";;;")
	if len(result.Items) != 0 {
		t.Fatalf("expected 0 items for semicolons only, got %d", len(result.Items))
	}
}

func TestParseMultipleStatements(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		count int
	}{
		{
			name:  "two selects",
			sql:   "SELECT 1; SELECT 2",
			count: 2,
		},
		{
			name:  "three statements",
			sql:   "SELECT 1; INSERT INTO t VALUES (1); UPDATE t SET a = 1",
			count: 3,
		},
		{
			name:  "mixed DDL and DML",
			sql:   "CREATE TABLE t (id INT); INSERT INTO t VALUES (1); SELECT * FROM t; DROP TABLE t",
			count: 4,
		},
		{
			name:  "trailing semicolons",
			sql:   "SELECT 1;;;SELECT 2;;;",
			count: 2,
		},
		{
			name: "all statement types via Parse()",
			sql: `SELECT 1;
INSERT INTO t VALUES (1);
REPLACE INTO t VALUES (2);
UPDATE t SET a = 1;
DELETE FROM t;
CREATE TABLE t2 (id INT);
ALTER TABLE t2 ADD COLUMN name VARCHAR(100);
DROP TABLE t2;
TRUNCATE TABLE t;
RENAME TABLE t TO t3;
SET @a = 1;
SHOW TABLES;
USE mydb;
EXPLAIN SELECT 1;
BEGIN;
COMMIT;
ROLLBACK;
SAVEPOINT sp1;
RELEASE SAVEPOINT sp1;
LOCK TABLES t READ;
UNLOCK TABLES;
GRANT SELECT ON *.* TO 'user'@'localhost';
REVOKE SELECT ON *.* FROM 'user'@'localhost';
LOAD DATA INFILE 'file.csv' INTO TABLE t;
PREPARE stmt FROM 'SELECT 1';
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
ANALYZE TABLE t;
OPTIMIZE TABLE t;
CHECK TABLE t;
REPAIR TABLE t;
FLUSH TABLES;
RESET MASTER;
KILL 42;
DO 1;
CHECKSUM TABLE t1;
SHUTDOWN;
RESTART`,
			count: 38,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			if len(result.Items) != tt.count {
				t.Fatalf("expected %d statements, got %d", tt.count, len(result.Items))
			}
		})
	}
}

func TestParseStmtDispatchCreateVariants(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT)",
		"CREATE TEMPORARY TABLE t (id INT)",
		"CREATE INDEX idx ON t (col)",
		"CREATE UNIQUE INDEX idx ON t (col)",
		"CREATE FULLTEXT INDEX idx ON t (col)",
		"CREATE SPATIAL INDEX idx ON t (col)",
		"CREATE VIEW v AS SELECT 1",
		"CREATE OR REPLACE VIEW v AS SELECT 1",
		"CREATE DATABASE mydb",
		"CREATE SCHEMA mydb",
		"CREATE FUNCTION f() RETURNS INT BEGIN RETURN 1; END",
		"CREATE PROCEDURE p() BEGIN SELECT 1; END",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1",
		"CREATE EVENT ev ON SCHEDULE AT CURRENT_TIMESTAMP DO SELECT 1",
		"CREATE USER 'user'@'localhost'",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseStmtDispatchDropVariants(t *testing.T) {
	tests := []string{
		"DROP TABLE t",
		"DROP TABLE IF EXISTS t",
		"DROP TEMPORARY TABLE t",
		"DROP INDEX idx ON t",
		"DROP VIEW v",
		"DROP DATABASE mydb",
		"DROP USER 'user'@'localhost'",
		"DROP FUNCTION IF EXISTS f",
		"DROP PROCEDURE IF EXISTS p",
		"DROP TRIGGER IF EXISTS tr",
		"DROP EVENT IF EXISTS ev",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 22: WITH (Common Table Expressions)
// -----------------------------------------------------------------------

func TestParseWithCTE(t *testing.T) {
	tests := []string{
		// Simple CTE
		"WITH cte AS (SELECT 1) SELECT * FROM cte",
		// Multiple CTEs
		"WITH cte1 AS (SELECT 1), cte2 AS (SELECT 2) SELECT * FROM cte1, cte2",
		// CTE with column list
		"WITH cte (a, b) AS (SELECT 1, 2) SELECT a, b FROM cte",
		// CTE with complex subquery
		"WITH cte AS (SELECT id, name FROM users WHERE active = 1) SELECT * FROM cte WHERE id > 10",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWithRecursiveCTE(t *testing.T) {
	tests := []string{
		// Basic recursive CTE
		"WITH RECURSIVE cte (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM cte WHERE n < 5) SELECT * FROM cte",
		// Recursive CTE without column list
		"WITH RECURSIVE cte AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM cte WHERE n < 10) SELECT * FROM cte",
		// Recursive CTE with UNION DISTINCT
		"WITH RECURSIVE cte (n) AS (SELECT 1 UNION DISTINCT SELECT n + 1 FROM cte WHERE n < 5) SELECT * FROM cte",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseInsertWithCTE(t *testing.T) {
	tests := []string{
		"INSERT INTO t1 WITH cte AS (SELECT 1 AS a) SELECT * FROM cte",
		"INSERT INTO t1 (col1) WITH cte AS (SELECT 1) SELECT * FROM cte",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseUpdateWithCTE(t *testing.T) {
	tests := []string{
		"WITH cte AS (SELECT id FROM t2) UPDATE t1 SET col1 = 1 WHERE id IN (SELECT id FROM cte)",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDeleteWithCTE(t *testing.T) {
	tests := []string{
		"WITH cte AS (SELECT id FROM t2) DELETE FROM t1 WHERE id IN (SELECT id FROM cte)",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 26: SET TRANSACTION
// -----------------------------------------------------------------------

func TestParseSetTransactionReadUncommitted(t *testing.T) {
	ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
}

func TestParseSetTransactionReadCommitted(t *testing.T) {
	ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
}

func TestParseSetTransactionRepeatableRead(t *testing.T) {
	ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
}

func TestParseSetTransactionSerializable(t *testing.T) {
	ParseAndCheck(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
}

func TestParseSetTransactionReadOnly(t *testing.T) {
	ParseAndCheck(t, "SET TRANSACTION READ ONLY")
}

func TestParseSetGlobalTransaction(t *testing.T) {
	tests := []string{
		"SET GLOBAL TRANSACTION ISOLATION LEVEL REPEATABLE READ",
		"SET SESSION TRANSACTION READ WRITE",
		"SET GLOBAL TRANSACTION ISOLATION LEVEL SERIALIZABLE, READ ONLY",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 25: XA Transactions
// -----------------------------------------------------------------------

func TestParseXAStart(t *testing.T) {
	tests := []string{
		"XA START 'xid1'",
		"XA BEGIN 'xid1'",
		"XA START 'xid1' JOIN",
		"XA START 'xid1' RESUME",
		"XA START 'gtrid', 'bqual'",
		"XA START 'gtrid', 'bqual', 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseXAEnd(t *testing.T) {
	tests := []string{
		"XA END 'xid1'",
		"XA END 'xid1' SUSPEND",
		"XA END 'xid1' SUSPEND FOR MIGRATE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseXAPrepare(t *testing.T) {
	ParseAndCheck(t, "XA PREPARE 'xid1'")
}

func TestParseXACommit(t *testing.T) {
	tests := []string{
		"XA COMMIT 'xid1'",
		"XA COMMIT 'xid1' ONE PHASE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseXARollback(t *testing.T) {
	ParseAndCheck(t, "XA ROLLBACK 'xid1'")
}

func TestParseXARecover(t *testing.T) {
	tests := []string{
		"XA RECOVER",
		"XA RECOVER CONVERT XID",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 24: JSON Operators
// -----------------------------------------------------------------------

func TestParseJsonExtract(t *testing.T) {
	tests := []string{
		"SELECT c->'$.name' FROM t",
		"SELECT t.c->'$.id' FROM t",
		"SELECT c->'$.a.b' FROM t WHERE c->'$.id' > 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseJsonUnquoteExtract(t *testing.T) {
	tests := []string{
		"SELECT c->>'$.name' FROM t",
		"SELECT t.c->>'$.id' FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseMemberOf(t *testing.T) {
	tests := []string{
		"SELECT 17 MEMBER OF('[23, 17, 10]')",
		"SELECT * FROM t WHERE val MEMBER OF(json_col)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseJsonTable(t *testing.T) {
	tests := []string{
		"SELECT * FROM JSON_TABLE('[1,2,3]', '$[*]' COLUMNS(id INT PATH '$')) AS jt",
		"SELECT * FROM JSON_TABLE('{\"a\":1}', '$' COLUMNS(a INT PATH '$.a', b VARCHAR(10) PATH '$.b' NULL ON EMPTY NULL ON ERROR)) AS jt",
		"SELECT * FROM JSON_TABLE('[1,2]', '$[*]' COLUMNS(rowid FOR ORDINALITY, val INT PATH '$')) AS jt",
		"SELECT * FROM JSON_TABLE('{\"a\":1}', '$' COLUMNS(a INT EXISTS PATH '$.a')) AS jt",
		"SELECT * FROM JSON_TABLE('[{\"a\":[1,2]}]', '$[*]' COLUMNS(NESTED PATH '$.a[*]' COLUMNS(val INT PATH '$'))) AS jt",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 23: Window OVER clause
// -----------------------------------------------------------------------

func TestParseWindowRowNumber(t *testing.T) {
	tests := []string{
		"SELECT ROW_NUMBER() OVER () FROM t",
		"SELECT ROW_NUMBER() OVER (ORDER BY id) FROM t",
		"SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY id) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWindowRank(t *testing.T) {
	tests := []string{
		"SELECT RANK() OVER (ORDER BY score DESC) FROM t",
		"SELECT DENSE_RANK() OVER (PARTITION BY dept ORDER BY score DESC) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWindowFrameRows(t *testing.T) {
	tests := []string{
		"SELECT SUM(x) OVER (ORDER BY id ROWS UNBOUNDED PRECEDING) FROM t",
		"SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t",
		"SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) FROM t",
		"SELECT SUM(x) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWindowFrameRange(t *testing.T) {
	tests := []string{
		"SELECT SUM(x) OVER (ORDER BY id RANGE UNBOUNDED PRECEDING) FROM t",
		"SELECT SUM(x) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWindowFrameGroups(t *testing.T) {
	tests := []string{
		"SELECT SUM(x) OVER (ORDER BY id GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t",
		"SELECT SUM(x) OVER (ORDER BY id GROUPS CURRENT ROW) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWindowPartitionOrderBy(t *testing.T) {
	tests := []string{
		"SELECT COUNT(*) OVER (PARTITION BY dept) FROM t",
		"SELECT AVG(salary) OVER (PARTITION BY dept ORDER BY hire_date) FROM t",
		"SELECT SUM(x) OVER (PARTITION BY a, b ORDER BY c ASC, d DESC) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseNamedWindow(t *testing.T) {
	tests := []string{
		"SELECT SUM(x) OVER w FROM t WINDOW w AS (ORDER BY id)",
		"SELECT SUM(x) OVER w, AVG(x) OVER w FROM t WINDOW w AS (PARTITION BY dept ORDER BY id)",
		"SELECT SUM(x) OVER w1, AVG(x) OVER w2 FROM t WINDOW w1 AS (ORDER BY id), w2 AS (ORDER BY name)",
		// OVER referencing a named window by name
		"SELECT SUM(x) OVER w FROM t WINDOW w AS (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 27: CALL and HANDLER
// -----------------------------------------------------------------------

func TestParseCall(t *testing.T) {
	tests := []string{
		"CALL my_proc()",
		"CALL my_proc",
		"CALL mydb.my_proc()",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCallWithArgs(t *testing.T) {
	tests := []string{
		"CALL my_proc(1, 2, 3)",
		"CALL my_proc('hello', @var)",
		"CALL mydb.my_proc(1 + 2, 'abc')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseHandlerOpen(t *testing.T) {
	tests := []string{
		"HANDLER t1 OPEN",
		"HANDLER t1 OPEN AS h1",
		"HANDLER mydb.t1 OPEN",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseHandlerRead(t *testing.T) {
	tests := []string{
		"HANDLER t1 READ FIRST",
		"HANDLER t1 READ NEXT",
		"HANDLER t1 READ PREV",
		"HANDLER t1 READ LAST",
		"HANDLER t1 READ idx FIRST",
		"HANDLER t1 READ idx NEXT",
		"HANDLER t1 READ NEXT LIMIT 5",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseHandlerReadWhere(t *testing.T) {
	tests := []string{
		"HANDLER t1 READ FIRST WHERE id > 10",
		"HANDLER t1 READ idx NEXT WHERE name = 'foo'",
		"HANDLER t1 READ NEXT WHERE id > 5 LIMIT 10",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseHandlerClose(t *testing.T) {
	tests := []string{
		"HANDLER t1 CLOSE",
		"HANDLER mydb.t1 CLOSE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSignal(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SIGNAL SQLSTATE '45000'",
			want: "{SIGNAL :loc 0 :condition 45000}",
		},
		{
			sql:  "SIGNAL SQLSTATE VALUE '45000'",
			want: "{SIGNAL :loc 0 :condition 45000}",
		},
		{
			sql:  "SIGNAL my_error",
			want: "{SIGNAL :loc 0 :condition my_error}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseSignalStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseSignalSet(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'An error occurred'",
			want: "{SIGNAL :loc 0 :condition 45000 :set {SIGNAL_INFO :loc 28 :name MESSAGE_TEXT :val {STRING_LIT :val \"An error occurred\" :loc 43}}}",
		},
		{
			sql:  "SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'oops', MYSQL_ERRNO = 1234",
			want: "{SIGNAL :loc 0 :condition 45000 :set {SIGNAL_INFO :loc 28 :name MESSAGE_TEXT :val {STRING_LIT :val \"oops\" :loc 43}} {SIGNAL_INFO :loc 51 :name MYSQL_ERRNO :val {INT_LIT :val 1234 :loc 65}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseSignalStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseResignal(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "RESIGNAL",
			want: "{RESIGNAL :loc 0}",
		},
		{
			sql:  "RESIGNAL SQLSTATE '45000'",
			want: "{RESIGNAL :loc 0 :condition 45000}",
		},
		{
			sql:  "RESIGNAL SET MESSAGE_TEXT = 'Error handled'",
			want: "{RESIGNAL :loc 0 :set {SIGNAL_INFO :loc 13 :name MESSAGE_TEXT :val {STRING_LIT :val \"Error handled\" :loc 28}}}",
		},
		{
			sql:  "RESIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 5678",
			want: "{RESIGNAL :loc 0 :condition 45000 :set {SIGNAL_INFO :loc 30 :name MYSQL_ERRNO :val {INT_LIT :val 5678 :loc 44}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseResignalStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseGetDiagnostics(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "GET DIAGNOSTICS @cnt = NUMBER",
			want: "{GET_DIAGNOSTICS :loc 0 :stmt_info true :items {DIAG_ITEM :loc 16 :target {VAR :name cnt :loc 16} :name NUMBER}}",
		},
		{
			sql:  "GET DIAGNOSTICS @cnt = NUMBER, @rc = ROW_COUNT",
			want: "{GET_DIAGNOSTICS :loc 0 :stmt_info true :items {DIAG_ITEM :loc 16 :target {VAR :name cnt :loc 16} :name NUMBER} {DIAG_ITEM :loc 31 :target {VAR :name rc :loc 31} :name ROW_COUNT}}",
		},
		{
			sql:  "GET CURRENT DIAGNOSTICS @cnt = NUMBER",
			want: "{GET_DIAGNOSTICS :loc 0 :stmt_info true :items {DIAG_ITEM :loc 24 :target {VAR :name cnt :loc 24} :name NUMBER}}",
		},
		{
			sql:  "GET DIAGNOSTICS cnt = NUMBER",
			want: "{GET_DIAGNOSTICS :loc 0 :stmt_info true :items {DIAG_ITEM :loc 16 :target {COLREF :loc 16 :col cnt} :name NUMBER}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseGetDiagnosticsStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseGetDiagnosticsCondition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "GET DIAGNOSTICS CONDITION 1 @msg = MESSAGE_TEXT",
			want: "{GET_DIAGNOSTICS :loc 0 :condition_number {INT_LIT :val 1 :loc 26} :items {DIAG_ITEM :loc 28 :target {VAR :name msg :loc 28} :name MESSAGE_TEXT}}",
		},
		{
			sql:  "GET DIAGNOSTICS CONDITION 1 @msg = MESSAGE_TEXT, @errno = MYSQL_ERRNO",
			want: "{GET_DIAGNOSTICS :loc 0 :condition_number {INT_LIT :val 1 :loc 26} :items {DIAG_ITEM :loc 28 :target {VAR :name msg :loc 28} :name MESSAGE_TEXT} {DIAG_ITEM :loc 49 :target {VAR :name errno :loc 49} :name MYSQL_ERRNO}}",
		},
		{
			sql:  "GET STACKED DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE",
			want: "{GET_DIAGNOSTICS :loc 0 :stacked true :condition_number {INT_LIT :val 1 :loc 34} :items {DIAG_ITEM :loc 36 :target {VAR :name sqlstate :loc 36} :name RETURNED_SQLSTATE}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseGetDiagnosticsStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseStmtDispatchAlterVariants(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ADD COLUMN c INT",
		"ALTER DATABASE mydb CHARACTER SET utf8mb4",
		"ALTER USER 'user'@'localhost' IDENTIFIED BY 'newpass'",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 29: Compound Statements
// -----------------------------------------------------------------------

func TestParseBeginEndBlock(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "BEGIN END",
			want: "{BEGIN_END :loc 0}",
		},
		{
			sql:  "BEGIN SELECT 1; END",
			want: "{BEGIN_END :loc 0 :stmts {SELECT :loc 6 :targets ({INT_LIT :val 1 :loc 13})}}",
		},
		{
			sql:  "BEGIN SELECT 1; SELECT 2; END",
			want: "{BEGIN_END :loc 0 :stmts {SELECT :loc 6 :targets ({INT_LIT :val 1 :loc 13})} {SELECT :loc 16 :targets ({INT_LIT :val 2 :loc 23})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseBeginEndBlock("", 0)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseBeginEndBlockLabeled(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "myblock: BEGIN END myblock",
			want: "{BEGIN_END :loc 0 :label myblock :end_label myblock}",
		},
		{
			sql:  "lbl: BEGIN SELECT 1; END lbl",
			want: "{BEGIN_END :loc 0 :label lbl :end_label lbl :stmts {SELECT :loc 11 :targets ({INT_LIT :val 1 :loc 18})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCompoundStmtOrStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseDeclareVar(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DECLARE x INT",
			want: "{DECLARE_VAR :loc 0 :names x :type {DATATYPE :loc 10 :name INT}}",
		},
		{
			sql:  "DECLARE a, b, c VARCHAR(100)",
			want: "{DECLARE_VAR :loc 0 :names a b c :type {DATATYPE :loc 16 :name VARCHAR :len 100}}",
		},
		{
			sql:  "DECLARE x INT DEFAULT 0",
			want: "{DECLARE_VAR :loc 0 :names x :type {DATATYPE :loc 10 :name INT} :default {INT_LIT :val 0 :loc 22}}",
		},
		{
			sql:  "DECLARE msg VARCHAR(255) DEFAULT 'hello'",
			want: "{DECLARE_VAR :loc 0 :names msg :type {DATATYPE :loc 12 :name VARCHAR :len 255} :default {STRING_LIT :val \"hello\" :loc 33}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseDeclareStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseDeclareCondition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DECLARE my_error CONDITION FOR SQLSTATE '45000'",
			want: "{DECLARE_CONDITION :loc 0 :name my_error :value 45000}",
		},
		{
			sql:  "DECLARE my_error CONDITION FOR SQLSTATE VALUE '45000'",
			want: "{DECLARE_CONDITION :loc 0 :name my_error :value 45000}",
		},
		{
			sql:  "DECLARE my_error CONDITION FOR 1051",
			want: "{DECLARE_CONDITION :loc 0 :name my_error :value 1051}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseDeclareStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseDeclareHandler(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SELECT 1",
			want: "{DECLARE_HANDLER :loc 0 :action CONTINUE :conditions SQLEXCEPTION :stmt {SELECT :loc 42 :targets ({INT_LIT :val 1 :loc 49})}}",
		},
		{
			sql:  "DECLARE EXIT HANDLER FOR NOT FOUND SELECT 1",
			want: "{DECLARE_HANDLER :loc 0 :action EXIT :conditions NOT FOUND :stmt {SELECT :loc 35 :targets ({INT_LIT :val 1 :loc 42})}}",
		},
		{
			sql:  "DECLARE CONTINUE HANDLER FOR SQLSTATE '23000' SELECT 1",
			want: "{DECLARE_HANDLER :loc 0 :action CONTINUE :conditions 23000 :stmt {SELECT :loc 46 :targets ({INT_LIT :val 1 :loc 53})}}",
		},
		{
			sql:  "DECLARE CONTINUE HANDLER FOR SQLWARNING, SQLEXCEPTION SELECT 1",
			want: "{DECLARE_HANDLER :loc 0 :action CONTINUE :conditions SQLWARNING SQLEXCEPTION :stmt {SELECT :loc 54 :targets ({INT_LIT :val 1 :loc 61})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseDeclareStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseDeclareCursor(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DECLARE cur1 CURSOR FOR SELECT id FROM t1",
			want: "{DECLARE_CURSOR :loc 0 :name cur1 :select {SELECT :loc 24 :targets ({COLREF :loc 31 :col id}) :from ({TABLEREF :loc 39 :name t1})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseDeclareStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseBeginEndWithDeclare(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "BEGIN DECLARE x INT DEFAULT 0; SELECT x; END",
			want: "{BEGIN_END :loc 0 :stmts {DECLARE_VAR :loc 6 :names x :type {DATATYPE :loc 16 :name INT} :default {INT_LIT :val 0 :loc 28}} {SELECT :loc 31 :targets ({COLREF :loc 38 :col x})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseBeginEndBlock("", 0)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Batch 31: Role Management
// ============================================================================

func TestParseCreateRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE ROLE myrole",
			want: "{CREATE_ROLE :loc 7 :roles myrole}",
		},
		{
			sql:  "CREATE ROLE IF NOT EXISTS myrole",
			want: "{CREATE_ROLE :loc 7 :if_not_exists true :roles myrole}",
		},
		{
			sql:  "CREATE ROLE role1, role2, role3",
			want: "{CREATE_ROLE :loc 7 :roles role1, role2, role3}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseDropRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP ROLE myrole",
			want: "{DROP_ROLE :loc 5 :roles myrole}",
		},
		{
			sql:  "DROP ROLE IF EXISTS myrole",
			want: "{DROP_ROLE :loc 5 :if_exists true :roles myrole}",
		},
		{
			sql:  "DROP ROLE role1, role2",
			want: "{DROP_ROLE :loc 5 :roles role1, role2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseSetDefaultRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SET DEFAULT ROLE NONE TO root",
			want: "{SET_DEFAULT_ROLE :loc 0 :none true :to root}",
		},
		{
			sql:  "SET DEFAULT ROLE ALL TO root",
			want: "{SET_DEFAULT_ROLE :loc 0 :all true :to root}",
		},
		{
			sql:  "SET DEFAULT ROLE myrole TO root",
			want: "{SET_DEFAULT_ROLE :loc 0 :roles myrole :to root}",
		},
		{
			sql:  "SET DEFAULT ROLE role1, role2 TO user1, user2",
			want: "{SET_DEFAULT_ROLE :loc 0 :roles role1, role2 :to user1, user2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseSetRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SET ROLE DEFAULT",
			want: "{SET_ROLE :loc 0 :default true}",
		},
		{
			sql:  "SET ROLE NONE",
			want: "{SET_ROLE :loc 0 :none true}",
		},
		{
			sql:  "SET ROLE ALL",
			want: "{SET_ROLE :loc 0 :all true}",
		},
		{
			sql:  "SET ROLE ALL EXCEPT role1, role2",
			want: "{SET_ROLE :loc 0 :all true :all_except role1, role2}",
		},
		{
			sql:  "SET ROLE myrole",
			want: "{SET_ROLE :loc 0 :roles myrole}",
		},
		{
			sql:  "SET ROLE role1, role2",
			want: "{SET_ROLE :loc 0 :roles role1, role2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseGrantRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "GRANT myrole TO myuser",
			want: "{GRANT_ROLE :loc 0 :roles MYROLE :to myuser}",
		},
		{
			sql:  "GRANT role1, role2 TO user1, user2",
			want: "{GRANT_ROLE :loc 0 :roles ROLE1, ROLE2 :to user1, user2}",
		},
		{
			sql:  "GRANT myrole TO myuser WITH ADMIN OPTION",
			want: "{GRANT_ROLE :loc 0 :roles MYROLE :to myuser :with_admin true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseRevokeRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "REVOKE myrole FROM myuser",
			want: "{REVOKE_ROLE :loc 0 :roles MYROLE :from myuser}",
		},
		{
			sql:  "REVOKE role1, role2 FROM user1, user2",
			want: "{REVOKE_ROLE :loc 0 :roles ROLE1, ROLE2 :from user1, user2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseChecksumTable(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CHECKSUM TABLE t1",
			want: "{CHECKSUM_TABLE :loc 0 :tables {TABLEREF :loc 15 :name t1}}",
		},
		{
			sql:  "CHECKSUM TABLE t1, t2",
			want: "{CHECKSUM_TABLE :loc 0 :tables {TABLEREF :loc 15 :name t1} {TABLEREF :loc 19 :name t2}}",
		},
		{
			sql:  "CHECKSUM TABLE t1 QUICK",
			want: "{CHECKSUM_TABLE :loc 0 :tables {TABLEREF :loc 15 :name t1} :quick true}",
		},
		{
			sql:  "CHECKSUM TABLE t1 EXTENDED",
			want: "{CHECKSUM_TABLE :loc 0 :tables {TABLEREF :loc 15 :name t1} :extended true}",
		},
		{
			sql:  "CHECKSUM TABLE db.t1, db.t2 QUICK",
			want: "{CHECKSUM_TABLE :loc 0 :tables {TABLEREF :loc 15 :schema db :name t1} {TABLEREF :loc 22 :schema db :name t2} :quick true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseShutdown(t *testing.T) {
	ParseAndCompare(t, "SHUTDOWN", "{SHUTDOWN :loc 0}")
}

func TestParseRestart(t *testing.T) {
	ParseAndCompare(t, "RESTART", "{RESTART :loc 0}")
}

// ============================================================================
// Batch 36: CLONE
// ============================================================================

func TestParseCloneStmt(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CLONE LOCAL DATA DIRECTORY '/tmp/clone'",
			want: `{CLONE :loc 0 :local true :directory "/tmp/clone"}`,
		},
		{
			sql:  "CLONE LOCAL DATA DIRECTORY = '/tmp/clone'",
			want: `{CLONE :loc 0 :local true :directory "/tmp/clone"}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'root'@'10.0.0.1':3306 IDENTIFIED BY 'secret'",
			want: `{CLONE :loc 0 :user "root" :host "10.0.0.1" :port 3306 :password "secret"}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'admin'@'db.example.com':3307 IDENTIFIED BY 'pass' DATA DIRECTORY '/backup/clone'",
			want: `{CLONE :loc 0 :directory "/backup/clone" :user "admin" :host "db.example.com" :port 3307 :password "pass"}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'admin'@'db.example.com':3307 IDENTIFIED BY 'pass' DATA DIRECTORY = '/backup/clone'",
			want: `{CLONE :loc 0 :directory "/backup/clone" :user "admin" :host "db.example.com" :port 3307 :password "pass"}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'root'@'10.0.0.1':3306 IDENTIFIED BY 'secret' REQUIRE SSL",
			want: `{CLONE :loc 0 :user "root" :host "10.0.0.1" :port 3306 :password "secret" :require_ssl true}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'root'@'10.0.0.1':3306 IDENTIFIED BY 'secret' REQUIRE NO SSL",
			want: `{CLONE :loc 0 :user "root" :host "10.0.0.1" :port 3306 :password "secret" :require_ssl false}`,
		},
		{
			sql:  "CLONE INSTANCE FROM 'root'@'10.0.0.1':3306 IDENTIFIED BY 'secret' DATA DIRECTORY '/data' REQUIRE SSL",
			want: `{CLONE :loc 0 :directory "/data" :user "root" :host "10.0.0.1" :port 3306 :password "secret" :require_ssl true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseInstallPlugin(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "INSTALL PLUGIN myplugin SONAME 'ha_example.so'",
			want: `{INSTALL_PLUGIN :loc 0 :name myplugin :soname "ha_example.so"}`,
		},
		{
			sql:  "INSTALL PLUGIN validate_password SONAME 'validate_password.so'",
			want: `{INSTALL_PLUGIN :loc 0 :name validate_password :soname "validate_password.so"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseUninstallPlugin(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "UNINSTALL PLUGIN myplugin",
			want: "{UNINSTALL_PLUGIN :loc 0 :name myplugin}",
		},
		{
			sql:  "UNINSTALL PLUGIN validate_password",
			want: "{UNINSTALL_PLUGIN :loc 0 :name validate_password}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseInstallComponent(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "INSTALL COMPONENT 'file://component_validate_password'",
			want: `{INSTALL_COMPONENT :loc 0 :components "file://component_validate_password"}`,
		},
		{
			sql:  "INSTALL COMPONENT 'file://component1', 'file://component2'",
			want: `{INSTALL_COMPONENT :loc 0 :components "file://component1" "file://component2"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseInstallComponentSet(t *testing.T) {
	// INSTALL COMPONENT with SET clause (MySQL 8.0)
	tests := []string{
		"INSTALL COMPONENT 'file://component_validate_password' SET @@GLOBAL.validate_password.length = 10",
		"INSTALL COMPONENT 'file://component1' SET @@GLOBAL.comp.var1 = 1, @@PERSIST.comp.var2 = 'abc'",
		"INSTALL COMPONENT 'file://c1', 'file://c2' SET @@GLOBAL.c1.x = 100",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseUninstallComponent(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "UNINSTALL COMPONENT 'file://component_validate_password'",
			want: `{UNINSTALL_COMPONENT :loc 0 :components "file://component_validate_password"}`,
		},
		{
			sql:  "UNINSTALL COMPONENT 'file://component1', 'file://component2'",
			want: `{UNINSTALL_COMPONENT :loc 0 :components "file://component1" "file://component2"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

// ============================================================================
// Batch 30: Flow Control
// ============================================================================

func TestParseIfElse(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "IF x > 0 THEN SELECT 1; END IF",
			want: "{IF_STMT :loc 0 :cond {BINEXPR :op > :loc 5 :left {COLREF :loc 3 :col x} :right {INT_LIT :val 0 :loc 7}} :then ({SELECT :loc 14 :targets ({INT_LIT :val 1 :loc 21})})}",
		},
		{
			sql:  "IF x > 0 THEN SELECT 1; ELSE SELECT 2; END IF",
			want: "{IF_STMT :loc 0 :cond {BINEXPR :op > :loc 5 :left {COLREF :loc 3 :col x} :right {INT_LIT :val 0 :loc 7}} :then ({SELECT :loc 14 :targets ({INT_LIT :val 1 :loc 21})}) :else ({SELECT :loc 29 :targets ({INT_LIT :val 2 :loc 36})})}",
		},
		{
			sql:  "IF x > 0 THEN SELECT 1; ELSEIF x > 1 THEN SELECT 2; END IF",
			want: "{IF_STMT :loc 0 :cond {BINEXPR :op > :loc 5 :left {COLREF :loc 3 :col x} :right {INT_LIT :val 0 :loc 7}} :then ({SELECT :loc 14 :targets ({INT_LIT :val 1 :loc 21})}) :elseifs ({ELSEIF :loc 24 :cond {BINEXPR :op > :loc 33 :left {COLREF :loc 31 :col x} :right {INT_LIT :val 1 :loc 35}} :then ({SELECT :loc 42 :targets ({INT_LIT :val 2 :loc 49})})})}",
		},
		{
			sql:  "IF x > 0 THEN SELECT 1; ELSEIF x > 1 THEN SELECT 2; ELSE SELECT 3; END IF",
			want: "{IF_STMT :loc 0 :cond {BINEXPR :op > :loc 5 :left {COLREF :loc 3 :col x} :right {INT_LIT :val 0 :loc 7}} :then ({SELECT :loc 14 :targets ({INT_LIT :val 1 :loc 21})}) :elseifs ({ELSEIF :loc 24 :cond {BINEXPR :op > :loc 33 :left {COLREF :loc 31 :col x} :right {INT_LIT :val 1 :loc 35}} :then ({SELECT :loc 42 :targets ({INT_LIT :val 2 :loc 49})})}) :else ({SELECT :loc 57 :targets ({INT_LIT :val 3 :loc 64})})}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseIfStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseCaseStmt(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CASE x WHEN 1 THEN SELECT 1; END CASE",
			want: "{CASE_STMT :loc 0 :operand {COLREF :loc 5 :col x} :whens ({WHEN :loc 7 :cond {INT_LIT :val 1 :loc 12} :then ({SELECT :loc 19 :targets ({INT_LIT :val 1 :loc 26})})})}",
		},
		{
			sql:  "CASE WHEN x > 0 THEN SELECT 1; WHEN x > 1 THEN SELECT 2; END CASE",
			want: "{CASE_STMT :loc 0 :whens ({WHEN :loc 5 :cond {BINEXPR :op > :loc 12 :left {COLREF :loc 10 :col x} :right {INT_LIT :val 0 :loc 14}} :then ({SELECT :loc 21 :targets ({INT_LIT :val 1 :loc 28})})} {WHEN :loc 31 :cond {BINEXPR :op > :loc 38 :left {COLREF :loc 36 :col x} :right {INT_LIT :val 1 :loc 40}} :then ({SELECT :loc 47 :targets ({INT_LIT :val 2 :loc 54})})})}",
		},
		{
			sql:  "CASE WHEN x > 0 THEN SELECT 1; ELSE SELECT 2; END CASE",
			want: "{CASE_STMT :loc 0 :whens ({WHEN :loc 5 :cond {BINEXPR :op > :loc 12 :left {COLREF :loc 10 :col x} :right {INT_LIT :val 0 :loc 14}} :then ({SELECT :loc 21 :targets ({INT_LIT :val 1 :loc 28})})}) :else ({SELECT :loc 36 :targets ({INT_LIT :val 2 :loc 43})})}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCaseStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseWhile(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "WHILE x > 0 DO SELECT 1; END WHILE",
			want: "{WHILE :loc 0 :cond {BINEXPR :op > :loc 8 :left {COLREF :loc 6 :col x} :right {INT_LIT :val 0 :loc 10}} :stmts ({SELECT :loc 15 :targets ({INT_LIT :val 1 :loc 22})})}",
		},
		{
			sql:  "lbl: WHILE x > 0 DO SELECT 1; END WHILE lbl",
			want: "{WHILE :loc 0 :label lbl :end_label lbl :cond {BINEXPR :op > :loc 13 :left {COLREF :loc 11 :col x} :right {INT_LIT :val 0 :loc 15}} :stmts ({SELECT :loc 20 :targets ({INT_LIT :val 1 :loc 27})})}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCompoundStmtOrStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseRepeat(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "REPEAT SELECT 1; UNTIL x > 0 END REPEAT",
			want: "{REPEAT :loc 0 :stmts ({SELECT :loc 7 :targets ({INT_LIT :val 1 :loc 14})}) :until {BINEXPR :op > :loc 25 :left {COLREF :loc 23 :col x} :right {INT_LIT :val 0 :loc 27}}}",
		},
		{
			sql:  "lbl: REPEAT SELECT 1; UNTIL x > 0 END REPEAT lbl",
			want: "{REPEAT :loc 0 :label lbl :end_label lbl :stmts ({SELECT :loc 12 :targets ({INT_LIT :val 1 :loc 19})}) :until {BINEXPR :op > :loc 30 :left {COLREF :loc 28 :col x} :right {INT_LIT :val 0 :loc 32}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCompoundStmtOrStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseLoop(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "LOOP SELECT 1; END LOOP",
			want: "{LOOP :loc 0 :stmts ({SELECT :loc 5 :targets ({INT_LIT :val 1 :loc 12})})}",
		},
		{
			sql:  "lbl: LOOP SELECT 1; END LOOP lbl",
			want: "{LOOP :loc 0 :label lbl :end_label lbl :stmts ({SELECT :loc 10 :targets ({INT_LIT :val 1 :loc 17})})}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCompoundStmtOrStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseLeave(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "LEAVE myloop",
			want: "{LEAVE :loc 0 :label myloop}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseLeaveStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseIterate(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "ITERATE myloop",
			want: "{ITERATE :loc 0 :label myloop}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseIterateStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// TestLabelIdentRejectsAmbiguous2 verifies that ambiguous_2 keywords (like BEGIN)
// cannot be used as SP labels. BEGIN is allowed as a general identifier (e.g.,
// CREATE TABLE begin) but NOT as a label_ident.
func TestLabelIdentRejectsAmbiguous2(t *testing.T) {
	// BEGIN as a table name should work (general ident context).
	ParseAndCheck(t, "CREATE TABLE begin (a INT)")

	// BEGIN as an SP label should fail — ambiguous_2 is excluded from label_ident.
	t.Run("begin as LEAVE label rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("LEAVE begin")}
		p.advance()
		_, err := p.parseLeaveStmt()
		if err == nil {
			t.Fatal("expected error for LEAVE begin (BEGIN is ambiguous_2, not allowed as label)")
		}
	})

	t.Run("begin as ITERATE label rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("ITERATE begin")}
		p.advance()
		_, err := p.parseIterateStmt()
		if err == nil {
			t.Fatal("expected error for ITERATE begin (BEGIN is ambiguous_2, not allowed as label)")
		}
	})

	t.Run("begin as block label rejected", func(t *testing.T) {
		// In a compound context, `begin: BEGIN ... END begin` should fail
		// because BEGIN cannot be used as a label.
		p := &Parser{lexer: NewLexer("begin: BEGIN END begin")}
		p.advance()
		_, err := p.parseCompoundStmtOrStmt()
		if err == nil {
			t.Fatal("expected error for begin: BEGIN ... END (BEGIN is ambiguous_2, not allowed as label)")
		}
	})
}

func TestParseReturn(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "RETURN 42",
			want: "{RETURN :loc 0 :expr {INT_LIT :val 42 :loc 7}}",
		},
		{
			sql:  "RETURN x + 1",
			want: "{RETURN :loc 0 :expr {BINEXPR :op + :loc 9 :left {COLREF :loc 7 :col x} :right {INT_LIT :val 1 :loc 11}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseReturnStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseOpenCursor(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "OPEN cur1",
			want: "{OPEN_CURSOR :loc 0 :name cur1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseOpenCursorStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseFetchCursor(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "FETCH cur1 INTO a",
			want: "{FETCH_CURSOR :loc 0 :name cur1 :into a}",
		},
		{
			sql:  "FETCH cur1 INTO a, b, c",
			want: "{FETCH_CURSOR :loc 0 :name cur1 :into a b c}",
		},
		{
			sql:  "FETCH NEXT FROM cur1 INTO a",
			want: "{FETCH_CURSOR :loc 0 :name cur1 :into a}",
		},
		{
			sql:  "FETCH FROM cur1 INTO a",
			want: "{FETCH_CURSOR :loc 0 :name cur1 :into a}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseFetchCursorStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseCloseCursor(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CLOSE cur1",
			want: "{CLOSE_CURSOR :loc 0 :name cur1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(tt.sql)}
			p.advance()
			stmt, err := p.parseCloseCursorStmt()
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Batch 33: TABLE and VALUES statements
// ============================================================================

func TestParseTableStmt(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "TABLE t1",
			want: "{TABLE_STMT :loc 0 :table {TABLEREF :loc 6 :name t1}}",
		},
		{
			sql:  "TABLE mydb.t1",
			want: "{TABLE_STMT :loc 0 :table {TABLEREF :loc 6 :schema mydb :name t1}}",
		},
		{
			sql:  "TABLE t1 ORDER BY id",
			want: "{TABLE_STMT :loc 0 :table {TABLEREF :loc 6 :name t1} :order_by ({ORDER_BY :loc 18 :expr {COLREF :loc 18 :col id}})}",
		},
		{
			sql:  "TABLE t1 LIMIT 10",
			want: "{TABLE_STMT :loc 0 :table {TABLEREF :loc 6 :name t1} :limit {LIMIT :loc 15 :count {INT_LIT :val 10 :loc 15}}}",
		},
		{
			sql:  "TABLE t1 ORDER BY name LIMIT 5 OFFSET 10",
			want: "{TABLE_STMT :loc 0 :table {TABLEREF :loc 6 :name t1} :order_by ({ORDER_BY :loc 18 :expr {COLREF :loc 18 :col name}}) :limit {LIMIT :loc 29 :count {INT_LIT :val 5 :loc 29} :offset {INT_LIT :val 10 :loc 38}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseValuesStmt(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "VALUES ROW(1, 2, 3)",
			want: "{VALUES_STMT :loc 0 :rows (({INT_LIT :val 1 :loc 11} {INT_LIT :val 2 :loc 14} {INT_LIT :val 3 :loc 17}))}",
		},
		{
			sql:  "VALUES ROW(1, 2), ROW(3, 4)",
			want: "{VALUES_STMT :loc 0 :rows (({INT_LIT :val 1 :loc 11} {INT_LIT :val 2 :loc 14}) ({INT_LIT :val 3 :loc 22} {INT_LIT :val 4 :loc 25}))}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseValuesRow(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "VALUES ROW('a', 'b')",
			want: `{VALUES_STMT :loc 0 :rows (({STRING_LIT :val "a" :loc 11} {STRING_LIT :val "b" :loc 16}))}`,
		},
		{
			sql:  "VALUES ROW(1, 'hello', NULL)",
			want: `{VALUES_STMT :loc 0 :rows (({INT_LIT :val 1 :loc 11} {STRING_LIT :val "hello" :loc 14} {NULL_LIT :loc 23}))}`,
		},
		{
			sql:  "VALUES ROW(1) LIMIT 1",
			want: "{VALUES_STMT :loc 0 :rows (({INT_LIT :val 1 :loc 11})) :limit {LIMIT :loc 20 :count {INT_LIT :val 1 :loc 20}}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			result := ParseAndCheck(t, tt.sql)
			got := ast.NodeToString(result.Items[0])
			if got != tt.want {
				t.Errorf("got:  %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestParseTableUnion(t *testing.T) {
	// TABLE statement can be parsed standalone
	ParseAndCheck(t, "TABLE t1")
	// VALUES statement can be parsed standalone
	ParseAndCheck(t, "VALUES ROW(1, 2)")
	// TABLE with ORDER BY and LIMIT
	ParseAndCheck(t, "TABLE t1 ORDER BY id LIMIT 10")
	// VALUES with multiple rows
	ParseAndCheck(t, "VALUES ROW(1, 2), ROW(3, 4), ROW(5, 6)")
	// VALUES with ORDER BY
	ParseAndCheck(t, "VALUES ROW(1, 'a'), ROW(2, 'b') ORDER BY column_0")
}

func TestParseCreateTablespace(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE TABLESPACE ts1",
			want: "{CREATE_TABLESPACE :loc 0 :name ts1}",
		},
		{
			sql:  "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd'",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :datafile "ts1.ibd"}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd' ENGINE = InnoDB",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :datafile "ts1.ibd" :engine InnoDB}`,
		},
		{
			sql:  "CREATE UNDO TABLESPACE ts1 ADD DATAFILE 'undo_ts1.ibu'",
			want: `{CREATE_TABLESPACE :loc 0 :undo true :name ts1 :datafile "undo_ts1.ibu"}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd' FILE_BLOCK_SIZE = 16384 ENCRYPTION = 'Y' ENGINE InnoDB",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :datafile "ts1.ibd" :file_block_size 16384 :encryption "Y" :engine InnoDB}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd' AUTOEXTEND_SIZE = 4M",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :datafile "ts1.ibd" :autoextend_size 4M}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 USE LOGFILE GROUP lg1",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :use_logfile_group lg1}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 EXTENT_SIZE = 1M INITIAL_SIZE = 256M MAX_SIZE = 512M",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :extent_size 1M :initial_size 256M :max_size 512M}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 NODEGROUP = 1 WAIT COMMENT = 'test ts' ENGINE NDB",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :nodegroup 1 :wait true :comment "test ts" :engine NDB}`,
		},
		{
			sql:  "CREATE TABLESPACE ts1 ENGINE_ATTRIBUTE = '{\"key\": \"val\"}'",
			want: `{CREATE_TABLESPACE :loc 0 :name ts1 :engine_attribute "{\"key\": \"val\"}"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseAlterTablespace(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "ALTER TABLESPACE ts1 ADD DATAFILE 'ts1_2.ibd'",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :add_datafile "ts1_2.ibd"}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 DROP DATAFILE 'ts1_2.ibd'",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :drop_datafile "ts1_2.ibd"}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 ADD DATAFILE 'ts1.ibd' ENGINE = NDB",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :add_datafile "ts1.ibd" :engine NDB}`,
		},
		{
			sql:  "ALTER UNDO TABLESPACE ts1 ADD DATAFILE 'undo.ibu'",
			want: `{ALTER_TABLESPACE :loc 0 :undo true :name ts1 :add_datafile "undo.ibu"}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 RENAME TO ts2",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :rename_to ts2}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 AUTOEXTEND_SIZE = 4M",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :autoextend_size 4M}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 ENCRYPTION = 'Y'",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :encryption "Y"}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 ENGINE_ATTRIBUTE = '{\"key\": \"val\"}'",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :engine_attribute "{\"key\": \"val\"}"}`,
		},
		{
			sql:  "ALTER TABLESPACE ts1 ADD DATAFILE 'f.ibd' WAIT ENGINE NDB",
			want: `{ALTER_TABLESPACE :loc 0 :name ts1 :add_datafile "f.ibd" :wait true :engine NDB}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropTablespace(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP TABLESPACE ts1",
			want: "{DROP_TABLESPACE :loc 0 :name ts1}",
		},
		{
			sql:  "DROP TABLESPACE ts1 ENGINE = InnoDB",
			want: "{DROP_TABLESPACE :loc 0 :name ts1 :engine InnoDB}",
		},
		{
			sql:  "DROP UNDO TABLESPACE ts1",
			want: "{DROP_TABLESPACE :loc 0 :undo true :name ts1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseCreateServer(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (HOST 'remote_host', DATABASE 'test', USER 'remote_user')",
			want: `{CREATE_SERVER :loc 0 :name s1 :wrapper mysql :options ("HOST remote_host" "DATABASE test" "USER remote_user")}`,
		},
		{
			sql:  "CREATE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')",
			want: `{CREATE_SERVER :loc 0 :name s1 :wrapper mysql :options ("USER root")}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropServer(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP SERVER s1",
			want: "{DROP_SERVER :loc 0 :name s1}",
		},
		{
			sql:  "DROP SERVER IF EXISTS s1",
			want: "{DROP_SERVER :loc 0 :if_exists true :name s1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseCreateSpatialRefSys(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[\"Greek\",DATUM[\"Greek\"]]'",
			want: `{CREATE_SPATIAL_REF_SYS :loc 0 :srid 4120 :name "Greek" :definition "GEOGCS[\"Greek\",DATUM[\"Greek\"]]"}`,
		},
		{
			sql:  "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek'",
			want: `{CREATE_SPATIAL_REF_SYS :loc 0 :or_replace true :srid 4120 :name "Greek"}`,
		},
		{
			sql:  "CREATE SPATIAL REFERENCE SYSTEM IF NOT EXISTS 4120 NAME 'Greek'",
			want: `{CREATE_SPATIAL_REF_SYS :loc 0 :if_not_exists true :srid 4120 :name "Greek"}`,
		},
		{
			sql:  "CREATE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' ORGANIZATION 'EPSG' IDENTIFIED BY 4120 DESCRIPTION 'Greek coordinate system'",
			want: `{CREATE_SPATIAL_REF_SYS :loc 0 :srid 4120 :name "Greek" :organization "EPSG" :org_srid 4120 :description "Greek coordinate system"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropSpatialRefSys(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP SPATIAL REFERENCE SYSTEM 4120",
			want: "{DROP_SPATIAL_REF_SYS :loc 0 :srid 4120}",
		},
		{
			sql:  "DROP SPATIAL REFERENCE SYSTEM IF EXISTS 4120",
			want: "{DROP_SPATIAL_REF_SYS :loc 0 :if_exists true :srid 4120}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseCreateResourceGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE RESOURCE GROUP rg1 TYPE = USER",
			want: "{CREATE_RESOURCE_GROUP :loc 0 :name rg1 :type USER}",
		},
		{
			sql:  "CREATE RESOURCE GROUP rg1 TYPE = SYSTEM VCPU = 0-3 THREAD_PRIORITY = 0 ENABLE",
			want: "{CREATE_RESOURCE_GROUP :loc 0 :name rg1 :type SYSTEM :vcpus (0-3) :thread_priority 0 :enable true}",
		},
		{
			sql:  "CREATE RESOURCE GROUP rg1 TYPE = USER VCPU = 0, 2, 4-7 DISABLE",
			want: "{CREATE_RESOURCE_GROUP :loc 0 :name rg1 :type USER :vcpus (0 2 4-7) :enable false}",
		},
		{
			sql:  "CREATE RESOURCE GROUP rg1 TYPE = SYSTEM THREAD_PRIORITY = -10",
			want: "{CREATE_RESOURCE_GROUP :loc 0 :name rg1 :type SYSTEM :thread_priority -10}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseAlterResourceGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "ALTER RESOURCE GROUP rg1 VCPU = 0-3",
			want: "{ALTER_RESOURCE_GROUP :loc 0 :name rg1 :vcpus (0-3)}",
		},
		{
			sql:  "ALTER RESOURCE GROUP rg1 THREAD_PRIORITY = 5 ENABLE",
			want: "{ALTER_RESOURCE_GROUP :loc 0 :name rg1 :thread_priority 5 :enable true}",
		},
		{
			sql:  "ALTER RESOURCE GROUP rg1 DISABLE FORCE",
			want: "{ALTER_RESOURCE_GROUP :loc 0 :name rg1 :enable false :force true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropResourceGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP RESOURCE GROUP rg1",
			want: "{DROP_RESOURCE_GROUP :loc 0 :name rg1}",
		},
		{
			sql:  "DROP RESOURCE GROUP rg1 FORCE",
			want: "{DROP_RESOURCE_GROUP :loc 0 :name rg1 :force true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 39: ALTER VIEW/EVENT/FUNCTION/PROCEDURE, DROP FUNCTION/PROCEDURE/TRIGGER/EVENT
// ---------------------------------------------------------------------------

func TestParseAlterView(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "ALTER VIEW v AS SELECT 1",
			want: "{ALTER_VIEW :loc 6 :name {TABLEREF :loc 11 :name v} :select {SELECT :loc 16 :targets ({INT_LIT :val 1 :loc 23})}}",
		},
		{
			sql:  "ALTER ALGORITHM = MERGE VIEW v AS SELECT 1",
			want: "{ALTER_VIEW :loc 6 :algorithm MERGE :name {TABLEREF :loc 29 :name v} :select {SELECT :loc 34 :targets ({INT_LIT :val 1 :loc 41})}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseAlterViewParseAndCheck(t *testing.T) {
	tests := []string{
		"ALTER VIEW v AS SELECT 1",
		"ALTER ALGORITHM = MERGE VIEW v AS SELECT 1",
		"ALTER VIEW v (a, b) AS SELECT 1, 2",
		"ALTER VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterEvent(t *testing.T) {
	tests := []string{
		"ALTER EVENT ev1 ENABLE",
		"ALTER EVENT ev1 DISABLE",
		"ALTER EVENT ev1 RENAME TO ev2",
		"ALTER EVENT ev1 COMMENT 'test'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterRoutine(t *testing.T) {
	tests := []string{
		"ALTER FUNCTION f1 COMMENT 'test'",
		"ALTER PROCEDURE p1 COMMENT 'test'",
		"ALTER FUNCTION f1 SQL SECURITY INVOKER",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDropRoutineStmt(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP FUNCTION f1",
			want: "{DROP_ROUTINE :loc 5 :name {TABLEREF :loc 14 :name f1}}",
		},
		{
			sql:  "DROP FUNCTION IF EXISTS f1",
			want: "{DROP_ROUTINE :loc 5 :if_exists true :name {TABLEREF :loc 24 :name f1}}",
		},
		{
			sql:  "DROP PROCEDURE p1",
			want: "{DROP_ROUTINE :loc 5 :is_procedure true :name {TABLEREF :loc 15 :name p1}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropTriggerStmtBatch39(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP TRIGGER tr1",
			want: "{DROP_TRIGGER :loc 5 :name {TABLEREF :loc 13 :name tr1}}",
		},
		{
			sql:  "DROP TRIGGER IF EXISTS mydb.tr1",
			want: "{DROP_TRIGGER :loc 5 :if_exists true :name {TABLEREF :loc 23 :schema mydb :name tr1}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropEventStmtBatch39(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP EVENT ev1",
			want: "{DROP_EVENT :loc 5 :name ev1}",
		},
		{
			sql:  "DROP EVENT IF EXISTS ev1",
			want: "{DROP_EVENT :loc 5 :if_exists true :name ev1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 41: GROUP BY WITH ROLLUP, LATERAL derived tables, GROUPING()
// ---------------------------------------------------------------------------

// TestParseGroupByWithRollup tests GROUP BY ... WITH ROLLUP.
func TestParseGroupByWithRollup(t *testing.T) {
	// Verify the WithRollup flag via parseSelect
	sel := parseSelect(t, "SELECT dept, SUM(salary) FROM emp GROUP BY dept WITH ROLLUP")
	if !sel.WithRollup {
		t.Error("expected WithRollup = true")
	}
	if len(sel.GroupBy) != 1 {
		t.Errorf("group_by count = %d, want 1", len(sel.GroupBy))
	}

	// Multi-column GROUP BY WITH ROLLUP
	sel2 := parseSelect(t, "SELECT col1, col2, COUNT(*) FROM t GROUP BY col1, col2 WITH ROLLUP")
	if !sel2.WithRollup {
		t.Error("expected WithRollup = true for multi-column")
	}
	if len(sel2.GroupBy) != 2 {
		t.Errorf("group_by count = %d, want 2", len(sel2.GroupBy))
	}

	// Without ROLLUP, flag should be false
	sel3 := parseSelect(t, "SELECT dept FROM emp GROUP BY dept")
	if sel3.WithRollup {
		t.Error("expected WithRollup = false for plain GROUP BY")
	}

	// WITH ROLLUP followed by HAVING
	sel4 := parseSelect(t, "SELECT dept, SUM(salary) FROM emp GROUP BY dept WITH ROLLUP HAVING SUM(salary) > 1000")
	if !sel4.WithRollup {
		t.Error("expected WithRollup = true before HAVING")
	}
	if sel4.Having == nil {
		t.Error("expected HAVING clause after WITH ROLLUP")
	}
}

// TestParseGroupingSets tests GROUPING() function (parsed as regular function call).
func TestParseGroupingSets(t *testing.T) {
	// GROUPING() is a regular function call — parsed by parseFuncCall
	ParseAndCheck(t, "SELECT col1, GROUPING(col1) FROM t GROUP BY col1 WITH ROLLUP")
	ParseAndCheck(t, "SELECT col1, col2, GROUPING(col1, col2) FROM t GROUP BY col1, col2 WITH ROLLUP")

	// GROUPING() in HAVING clause
	ParseAndCheck(t, "SELECT col1, SUM(col2) FROM t GROUP BY col1 WITH ROLLUP HAVING GROUPING(col1) = 1")
}

// TestParseCube is a placeholder — MySQL 8.0 does not support CUBE syntax directly.
// GROUP BY ... WITH ROLLUP is the MySQL equivalent.
func TestParseCube(t *testing.T) {
	// MySQL does not support CUBE; WITH ROLLUP is the supported syntax.
	// Just verify WITH ROLLUP parses correctly.
	ParseAndCheck(t, "SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP")
}

// TestParseRollup tests GROUP BY WITH ROLLUP (same as TestParseGroupByWithRollup).
func TestParseRollup(t *testing.T) {
	sel := parseSelect(t, "SELECT year, SUM(amount) FROM sales GROUP BY year WITH ROLLUP")
	if !sel.WithRollup {
		t.Error("expected WithRollup = true")
	}
}

// TestParseOrderByWithRollup tests ORDER BY ... WITH ROLLUP (MySQL 8.0.12+).
func TestParseOrderByWithRollup(t *testing.T) {
	sel := parseSelect(t, "SELECT year, SUM(amount) FROM sales GROUP BY year ORDER BY year WITH ROLLUP")
	if !sel.OrderByWithRollup {
		t.Error("expected OrderByWithRollup = true")
	}
	if len(sel.OrderBy) != 1 {
		t.Errorf("order_by count = %d, want 1", len(sel.OrderBy))
	}

	// ORDER BY multiple columns WITH ROLLUP
	ParseAndCheck(t, "SELECT a, b, SUM(c) FROM t ORDER BY a, b WITH ROLLUP")

	// ORDER BY WITH ROLLUP followed by LIMIT
	ParseAndCheck(t, "SELECT a, SUM(b) FROM t ORDER BY a WITH ROLLUP LIMIT 10")
}

// TestParseLateralDerivedTable tests LATERAL derived tables in FROM clause.
func TestParseLateralDerivedTable(t *testing.T) {
	// LATERAL with AS alias
	ParseAndCheck(t, "SELECT * FROM t1, LATERAL (SELECT * FROM t2 WHERE t2.id = t1.id) AS lat")

	// LATERAL with implicit alias
	ParseAndCheck(t, "SELECT * FROM t1, LATERAL (SELECT * FROM t2) lat")

	// LATERAL without alias
	ParseAndCheck(t, "SELECT * FROM t1, LATERAL (SELECT 1)")

	// LATERAL in JOIN
	ParseAndCheck(t, "SELECT * FROM t1 JOIN LATERAL (SELECT * FROM t2 WHERE t2.a = t1.a) AS sub ON 1=1")

	// LATERAL with correlated subquery
	ParseAndCheck(t, "SELECT * FROM t1 LEFT JOIN LATERAL (SELECT * FROM t2 WHERE t2.fk = t1.pk LIMIT 3) AS top3 ON TRUE")

	// Verify SubqueryExpr has Lateral flag set
	result := ParseAndCheck(t, "SELECT * FROM t1, LATERAL (SELECT * FROM t2 WHERE t2.id = t1.id) AS lat")
	sel := result.Items[0].(*ast.SelectStmt)
	if len(sel.From) != 2 {
		t.Fatalf("from count = %d, want 2", len(sel.From))
	}
	sub, ok := sel.From[1].(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("from[1] type = %T, want *SubqueryExpr", sel.From[1])
	}
	if !sub.Lateral {
		t.Error("expected Lateral = true")
	}
	if sub.Alias != "lat" {
		t.Errorf("alias = %q, want %q", sub.Alias, "lat")
	}
}

// ---------------------------------------------------------------------------
// Batch 42: Phase 2 Statement Dispatch — comprehensive test
// ---------------------------------------------------------------------------

func TestParsePhase2Dispatch(t *testing.T) {
	tests := []string{
		// Signal/Resignal/Get Diagnostics (batch 28)
		"SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'error'",
		"RESIGNAL",
		"GET DIAGNOSTICS @p1 = ROW_COUNT",

		// Compound statements (batches 29, 30) tested separately with delimiter-aware parsing

		// Role management (batch 31)
		"CREATE ROLE r1",
		"DROP ROLE r1",
		"SET ROLE DEFAULT",
		"GRANT r1 TO u1",

		// TABLE / VALUES (batch 33)
		"TABLE t1",
		"VALUES ROW(1, 2), ROW(3, 4)",

		// Utility (batch 34)
		"CHECKSUM TABLE t1",
		"SHUTDOWN",
		"RESTART",

		// Plugin/Component (batch 35)
		"INSTALL PLUGIN p SONAME 'p.so'",
		"UNINSTALL PLUGIN p",
		"INSTALL COMPONENT 'file://c1'",
		"UNINSTALL COMPONENT 'file://c1'",

		// Clone (batch 36)
		"CLONE LOCAL DATA DIRECTORY = '/tmp/clone'",

		// Tablespace/Server (batch 37)
		"CREATE TABLESPACE ts1 ADD DATAFILE 'ts1.ibd'",
		"DROP TABLESPACE ts1",
		"CREATE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (HOST '127.0.0.1')",
		"DROP SERVER s1",

		// SRS/Resource Group (batch 38)
		"DROP SPATIAL REFERENCE SYSTEM 4326",
		"CREATE RESOURCE GROUP rg1 TYPE = USER",
		"DROP RESOURCE GROUP rg1",

		// ALTER/DROP misc (batch 39)
		"ALTER VIEW v AS SELECT 1",
		"ALTER EVENT ev1 ENABLE",
		"ALTER FUNCTION f1 COMMENT 'test'",
		"ALTER PROCEDURE p1 COMMENT 'test'",
		"DROP FUNCTION f1",
		"DROP PROCEDURE p1",
		"DROP TRIGGER tr1",
		"DROP EVENT ev1",

		// EXPLAIN (batch 40)
		"EXPLAIN SELECT 1",
		"EXPLAIN FORMAT=JSON SELECT 1",

		// GROUP BY WITH ROLLUP, LATERAL (batch 41)
		"SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP",
		"SELECT * FROM t1, LATERAL (SELECT * FROM t2) AS sub",

		// XA, CALL, HANDLER
		"XA START 'xid1'",
		"CALL my_proc(1, 2)",
		"HANDLER t1 OPEN",

		// Replication (batch 43)
		"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1'",
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1)",

		// Replication control (batch 44)
		"START REPLICA",
		"STOP REPLICA",
		"RESET REPLICA",
		"PURGE BINARY LOGS TO 'mysql-bin.010'",
		"RESET MASTER",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------- Batch 43: Replication ----------

func TestParseChangeReplicationSource(t *testing.T) {
	tests := []string{
		// Basic single option
		"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1'",
		// Multiple options
		"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1', SOURCE_PORT = 3306",
		// String options
		"CHANGE REPLICATION SOURCE TO SOURCE_USER = 'repl', SOURCE_PASSWORD = 'secret'",
		// Numeric options
		"CHANGE REPLICATION SOURCE TO SOURCE_PORT = 3306, SOURCE_LOG_POS = 1234",
		// Boolean-like options (0/1)
		"CHANGE REPLICATION SOURCE TO SOURCE_AUTO_POSITION = 1",
		"CHANGE REPLICATION SOURCE TO SOURCE_SSL = 0",
		// NULL value
		"CHANGE REPLICATION SOURCE TO PRIVILEGE_CHECKS_USER = NULL",
		// Identifier value options
		"CHANGE REPLICATION SOURCE TO REQUIRE_TABLE_PRIMARY_KEY_CHECK = STREAM",
		"CHANGE REPLICATION SOURCE TO REQUIRE_TABLE_PRIMARY_KEY_CHECK = ON",
		"CHANGE REPLICATION SOURCE TO REQUIRE_TABLE_PRIMARY_KEY_CHECK = OFF",
		// IGNORE_SERVER_IDS
		"CHANGE REPLICATION SOURCE TO IGNORE_SERVER_IDS = (1, 2, 3)",
		"CHANGE REPLICATION SOURCE TO IGNORE_SERVER_IDS = ()",
		// FOR CHANNEL
		"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1' FOR CHANNEL ch1",
		// SSL options
		"CHANGE REPLICATION SOURCE TO SOURCE_SSL = 1, SOURCE_SSL_CA = '/etc/ssl/ca.pem', SOURCE_SSL_CERT = '/etc/ssl/cert.pem'",
		// Many options combined
		"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1', SOURCE_PORT = 3306, SOURCE_USER = 'repl', SOURCE_PASSWORD = 'secret', SOURCE_AUTO_POSITION = 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseChangeReplicationSourceOptions(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1'",
			want: "{CHANGE_REPLICATION_SOURCE :loc 0 :options ({REPL_OPT :loc 29 :name SOURCE_HOST :val host1})}",
		},
		{
			sql:  "CHANGE REPLICATION SOURCE TO SOURCE_PORT = 3306",
			want: "{CHANGE_REPLICATION_SOURCE :loc 0 :options ({REPL_OPT :loc 29 :name SOURCE_PORT :val 3306})}",
		},
		{
			sql:  "CHANGE REPLICATION SOURCE TO PRIVILEGE_CHECKS_USER = NULL",
			want: "{CHANGE_REPLICATION_SOURCE :loc 0 :options ({REPL_OPT :loc 29 :name PRIVILEGE_CHECKS_USER :val NULL})}",
		},
		{
			sql:  "CHANGE REPLICATION SOURCE TO IGNORE_SERVER_IDS = (1, 2)",
			want: "{CHANGE_REPLICATION_SOURCE :loc 0 :options ({REPL_OPT :loc 29 :name IGNORE_SERVER_IDS :ids (1 2)})}",
		},
		{
			sql:  "CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'host1' FOR CHANNEL ch1",
			want: "{CHANGE_REPLICATION_SOURCE :loc 0 :options ({REPL_OPT :loc 29 :name SOURCE_HOST :val host1}) :channel ch1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseChangeReplicationFilter(t *testing.T) {
	tests := []string{
		// Single filter
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1)",
		// Multiple databases
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1, db2, db3)",
		// Ignore DB
		"CHANGE REPLICATION FILTER REPLICATE_IGNORE_DB = (db1)",
		// Table filter
		"CHANGE REPLICATION FILTER REPLICATE_DO_TABLE = (db1.t1)",
		// Multiple tables
		"CHANGE REPLICATION FILTER REPLICATE_DO_TABLE = (db1.t1, db2.t2)",
		// Ignore table
		"CHANGE REPLICATION FILTER REPLICATE_IGNORE_TABLE = (db1.t1)",
		// Wild do table
		"CHANGE REPLICATION FILTER REPLICATE_WILD_DO_TABLE = ('db%.t%')",
		// Wild ignore table
		"CHANGE REPLICATION FILTER REPLICATE_WILD_IGNORE_TABLE = ('db%.t%')",
		// Rewrite DB
		"CHANGE REPLICATION FILTER REPLICATE_REWRITE_DB = ((db1, db2))",
		// Multiple rewrite pairs
		"CHANGE REPLICATION FILTER REPLICATE_REWRITE_DB = ((db1, db2), (db3, db4))",
		// Multiple filters
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1), REPLICATE_IGNORE_DB = (db2)",
		// Empty filter list
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = ()",
		// FOR CHANNEL
		"CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1) FOR CHANNEL ch1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------- Batch 44: Replication Control ----------

func TestParseStartReplica(t *testing.T) {
	tests := []string{
		"START REPLICA",
		"START SLAVE",
		"START REPLICA IO_THREAD",
		"START REPLICA SQL_THREAD",
		"START REPLICA IO_THREAD, SQL_THREAD",
		"START REPLICA FOR CHANNEL ch1",
		"START REPLICA IO_THREAD FOR CHANNEL ch1",
		"START REPLICA UNTIL SOURCE_LOG_FILE = 'binlog.000002', SOURCE_LOG_POS = 1234",
		"START REPLICA UNTIL SQL_BEFORE_GTIDS = '3E11FA47-71CA-11E1-9E33-C80AA9429562:1-5'",
		"START REPLICA UNTIL SQL_AFTER_MTS_GAPS",
		"START REPLICA USER = 'repl' PASSWORD = 'secret'",
		"START REPLICA USER = 'repl' PASSWORD = 'secret' DEFAULT_AUTH = 'mysql_native_password'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseStopReplica(t *testing.T) {
	tests := []string{
		"STOP REPLICA",
		"STOP SLAVE",
		"STOP REPLICA IO_THREAD",
		"STOP REPLICA SQL_THREAD",
		"STOP REPLICA IO_THREAD, SQL_THREAD",
		"STOP REPLICA FOR CHANNEL ch1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseResetReplica(t *testing.T) {
	tests := []string{
		"RESET REPLICA",
		"RESET SLAVE",
		"RESET REPLICA ALL",
		"RESET REPLICA FOR CHANNEL ch1",
		"RESET REPLICA ALL FOR CHANNEL ch1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParsePurgeBinaryLogs(t *testing.T) {
	tests := []string{
		"PURGE BINARY LOGS TO 'mysql-bin.010'",
		"PURGE BINARY LOGS BEFORE '2019-04-02 22:46:26'",
		"PURGE MASTER LOGS TO 'mysql-bin.010'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseResetMaster(t *testing.T) {
	tests := []string{
		"RESET MASTER",
		"RESET MASTER TO 1234",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ============================================================================
// Batch 45: GROUP REPLICATION
// ============================================================================

func TestParseStartGroupReplication(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			"START GROUP_REPLICATION",
			"{START_GROUP_REPLICATION :loc 0}",
		},
		{
			"START GROUP_REPLICATION USER='repl', PASSWORD='secret'",
			"{START_GROUP_REPLICATION :loc 0 :user repl :password ***}",
		},
		{
			"START GROUP_REPLICATION USER='repl', PASSWORD='secret', DEFAULT_AUTH='mysql_native_password'",
			"{START_GROUP_REPLICATION :loc 0 :user repl :password *** :default_auth mysql_native_password}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

func TestParseStopGroupReplication(t *testing.T) {
	ParseAndCompare(t, "STOP GROUP_REPLICATION", "{STOP_GROUP_REPLICATION :loc 0}")
}

// ============================================================================
// Batch 46: INSTANCE ADMIN
// ============================================================================

func TestParseAlterInstanceRotateKey(t *testing.T) {
	ParseAndCheck(t, "ALTER INSTANCE ROTATE INNODB MASTER KEY")
}

func TestParseAlterInstanceRedoLog(t *testing.T) {
	tests := []string{
		"ALTER INSTANCE ENABLE INNODB REDO_LOG",
		"ALTER INSTANCE DISABLE INNODB REDO_LOG",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterInstanceReloadTLS(t *testing.T) {
	tests := []string{
		"ALTER INSTANCE RELOAD TLS",
		"ALTER INSTANCE RELOAD TLS NO ROLLBACK ON ERROR",
		"ALTER INSTANCE RELOAD TLS FOR CHANNEL mysql_main",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseLockInstance(t *testing.T) {
	ParseAndCompare(t, "LOCK INSTANCE FOR BACKUP", "{LOCK_INSTANCE :loc 0}")
}

func TestParseUnlockInstance(t *testing.T) {
	ParseAndCompare(t, "UNLOCK INSTANCE", "{UNLOCK_INSTANCE :loc 0}")
}

func TestParseImportTable(t *testing.T) {
	ParseAndCheck(t, "IMPORT TABLE FROM '/tmp/t1.sdi', '/tmp/t2.sdi'")
}

// ============================================================================
// Batch 47: CACHE/BINLOG
// ============================================================================

func TestParseBinlog(t *testing.T) {
	ParseAndCompare(t, "BINLOG 'base64str'", "{BINLOG :loc 0 :str \"base64str\"}")
}

func TestParseCacheIndex(t *testing.T) {
	ParseAndCheck(t, "CACHE INDEX t1 IN hot_cache")
}

func TestParseLoadIndexIntoCache(t *testing.T) {
	ParseAndCheck(t, "LOAD INDEX INTO CACHE t1")
}

func TestParseResetPersist(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			"RESET PERSIST",
			"{RESET_PERSIST :loc 0}",
		},
		{
			"RESET PERSIST max_connections",
			"{RESET_PERSIST :loc 0 :variable max_connections}",
		},
		{
			"RESET PERSIST IF EXISTS max_connections",
			"{RESET_PERSIST :loc 0 :if_exists true :variable max_connections}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

// ============================================================================
// Batch 48: ACCOUNT MISC
// ============================================================================

func TestParseRenameUser(t *testing.T) {
	tests := []string{
		"RENAME USER foo TO bar",
		"RENAME USER foo TO bar, baz TO qux",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSetResourceGroup(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			"SET RESOURCE GROUP rg1",
			"{SET_RESOURCE_GROUP :loc 0 :name rg1}",
		},
		{
			"SET RESOURCE GROUP rg1 FOR 1, 2, 3",
			"{SET_RESOURCE_GROUP :loc 0 :name rg1 :threads (1 2 3)}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

func TestParseHelp(t *testing.T) {
	ParseAndCompare(t, "HELP 'contents'", "{HELP :loc 0 :topic \"contents\"}")
}

// ============================================================================
// Batch 49: VCPUSpec depth fix (tested via CREATE RESOURCE GROUP)
// ============================================================================

func TestDeparseCreateResourceGroup(t *testing.T) {
	ParseAndCheck(t, "CREATE RESOURCE GROUP rg1 TYPE = USER VCPU = 0-3")
}

// ============================================================================
// Batch 51: SHOW REMAINING
// ============================================================================

func TestParseShowCreateUser(t *testing.T) {
	tests := []string{
		"SHOW CREATE USER root",
		"SHOW CREATE USER current_user",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowReplicas(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			"SHOW REPLICAS",
			`{SHOW :loc 0 :type REPLICAS}`,
		},
		{
			"SHOW SLAVE HOSTS",
			`{SHOW :loc 0 :type REPLICAS}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

func TestParseShowEngineInnodbStatus(t *testing.T) {
	tests := []string{
		"SHOW ENGINE INNODB STATUS",
		"SHOW ENGINE INNODB MUTEX",
		"SHOW ENGINE PERFORMANCE_SCHEMA STATUS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowProfile(t *testing.T) {
	tests := []string{
		"SHOW PROFILE",
		"SHOW PROFILE CPU",
		"SHOW PROFILE CPU, BLOCK IO",
		"SHOW PROFILE ALL FOR QUERY 2",
		"SHOW PROFILE CPU FOR QUERY 1 LIMIT 5",
		"SHOW PROFILE SOURCE, SWAPS FOR QUERY 3",
		"SHOW PROFILE CONTEXT SWITCHES",
		"SHOW PROFILE PAGE FAULTS, IPC, MEMORY",
		"SHOW PROFILE ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ============================================================================
// Batch 88: SHOW PROFILE LIMIT/OFFSET expression parsing
// ============================================================================

func TestParseShowProfileLimitOffset(t *testing.T) {
	tests := []string{
		"SHOW PROFILE CPU FOR QUERY 1 LIMIT 5",
		"SHOW PROFILE CPU FOR QUERY 1 LIMIT 10 OFFSET 5",
		"SHOW PROFILE ALL FOR QUERY 2 LIMIT 100",
		"SHOW PROFILE LIMIT 3",
		"SHOW PROFILE ALL LIMIT 5 OFFSET 10",
		"SHOW PROFILE CPU, BLOCK IO FOR QUERY 1 LIMIT 20 OFFSET 0",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowFunctionCode(t *testing.T) {
	tests := []string{
		"SHOW FUNCTION CODE my_func",
		"SHOW FUNCTION CODE db1.my_func",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ============================================================================
// Batch 52: WITH CHECK OPTION depth fix
// ============================================================================

func TestParseCreateViewWithCheckOption(t *testing.T) {
	tests := []string{
		"CREATE VIEW v AS SELECT 1 WITH CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION",
		"CREATE OR REPLACE VIEW v AS SELECT * FROM t WITH CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterViewWithCheckOption(t *testing.T) {
	tests := []string{
		"ALTER VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION",
		"ALTER VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowProcedureCode(t *testing.T) {
	tests := []string{
		"SHOW PROCEDURE CODE my_proc",
		"SHOW PROCEDURE CODE db1.my_proc",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseChangeMasterTo(t *testing.T) {
	tests := []string{
		// Basic single option
		"CHANGE MASTER TO MASTER_HOST = 'host1'",
		// Multiple options
		"CHANGE MASTER TO MASTER_HOST = 'host1', MASTER_PORT = 3306",
		// String options
		"CHANGE MASTER TO MASTER_USER = 'repl', MASTER_PASSWORD = 'secret'",
		// Numeric options
		"CHANGE MASTER TO MASTER_PORT = 3306, MASTER_LOG_POS = 1234",
		// Boolean-like options (0/1)
		"CHANGE MASTER TO MASTER_AUTO_POSITION = 1",
		"CHANGE MASTER TO MASTER_SSL = 0",
		// NULL value
		"CHANGE MASTER TO PRIVILEGE_CHECKS_USER = NULL",
		// Identifier value options
		"CHANGE MASTER TO REQUIRE_TABLE_PRIMARY_KEY_CHECK = STREAM",
		// IGNORE_SERVER_IDS
		"CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 2, 3)",
		// Empty IGNORE_SERVER_IDS
		"CHANGE MASTER TO IGNORE_SERVER_IDS = ()",
		// FOR CHANNEL
		"CHANGE MASTER TO MASTER_HOST = 'host1' FOR CHANNEL ch1",
		// SSL options
		"CHANGE MASTER TO MASTER_SSL = 1, MASTER_SSL_CA = '/ca.pem', MASTER_SSL_CERT = '/cert.pem'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseChangeMasterToOptions(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CHANGE MASTER TO MASTER_HOST = 'host1'",
			want: "{CHANGE_MASTER :loc 0 :options ({REPL_OPT :loc 17 :name MASTER_HOST :val host1})}",
		},
		{
			sql:  "CHANGE MASTER TO MASTER_PORT = 3306",
			want: "{CHANGE_MASTER :loc 0 :options ({REPL_OPT :loc 17 :name MASTER_PORT :val 3306})}",
		},
		{
			sql:  "CHANGE MASTER TO PRIVILEGE_CHECKS_USER = NULL",
			want: "{CHANGE_MASTER :loc 0 :options ({REPL_OPT :loc 17 :name PRIVILEGE_CHECKS_USER :val NULL})}",
		},
		{
			sql:  "CHANGE MASTER TO MASTER_HOST = 'host1' FOR CHANNEL ch1",
			want: "{CHANGE_MASTER :loc 0 :options ({REPL_OPT :loc 17 :name MASTER_HOST :val host1}) :channel ch1}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseSetPassword(t *testing.T) {
	tests := []string{
		"SET PASSWORD = 'newpass'",
		"SET PASSWORD = PASSWORD('newpass')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSetPasswordForUser(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SET PASSWORD = 'newpass'",
			want: "{SET_PASSWORD :loc 0 :password newpass}",
		},
		{
			sql:  "SET PASSWORD FOR 'admin'@'localhost' = 'newpass'",
			want: "{SET_PASSWORD :loc 0 :user {USER_SPEC :loc 17 :name admin :host localhost} :password newpass}",
		},
		{
			sql:  "SET PASSWORD FOR root = 'newpass'",
			want: "{SET_PASSWORD :loc 0 :user {USER_SPEC :loc 17 :name root} :password newpass}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseGrantAs(t *testing.T) {
	tests := []string{
		"GRANT SELECT ON db.* TO 'user1'@'localhost' AS 'admin'@'localhost'",
		"GRANT INSERT ON db.t TO user1 AS admin",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseGrantAsWithRole(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "GRANT SELECT ON db.* TO user1 AS admin WITH ROLE DEFAULT",
			want: "{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_role_type DEFAULT}",
		},
		{
			sql:  "GRANT SELECT ON db.* TO user1 AS admin WITH ROLE NONE",
			want: "{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_role_type NONE}",
		},
		{
			sql:  "GRANT SELECT ON db.* TO user1 AS admin WITH ROLE ALL",
			want: "{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_role_type ALL}",
		},
		{
			sql:  "GRANT SELECT ON db.* TO user1 AS admin WITH ROLE ALL EXCEPT role1, role2",
			want: "{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_role_type ALL EXCEPT :with_roles role1, role2}",
		},
		{
			sql:  "GRANT SELECT ON db.* TO user1 AS admin WITH ROLE role1, role2",
			want: "{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_roles role1, role2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

// TestRoleIdentRejectsAmbiguous3 verifies that ambiguous_3 keywords (like EVENT)
// cannot be used as role names. EVENT is allowed as a general identifier (e.g.,
// CREATE TABLE event) but NOT as a role_ident.
func TestRoleIdentRejectsAmbiguous3(t *testing.T) {
	// EVENT as a role in GRANT WITH ROLE should fail — ambiguous_3 is excluded from role_ident.
	t.Run("event as WITH ROLE name rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT SELECT ON db.* TO user1 AS admin WITH ROLE event")}
		p.advance()
		_, err := p.parseGrantStmt()
		if err == nil {
			t.Fatal("expected error for WITH ROLE event (EVENT is ambiguous_3, not allowed as role)")
		}
	})

	t.Run("event as WITH ROLE ALL EXCEPT name rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT SELECT ON db.* TO user1 AS admin WITH ROLE ALL EXCEPT event")}
		p.advance()
		_, err := p.parseGrantStmt()
		if err == nil {
			t.Fatal("expected error for WITH ROLE ALL EXCEPT event (EVENT is ambiguous_3, not allowed as role)")
		}
	})

	// BEGIN (ambiguous_2) should be accepted as a role name.
	t.Run("begin as WITH ROLE name accepted", func(t *testing.T) {
		ParseAndCompare(t,
			"GRANT SELECT ON db.* TO user1 AS admin WITH ROLE begin",
			"{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_roles begin}",
		)
	})

	t.Run("begin as WITH ROLE ALL EXCEPT name accepted", func(t *testing.T) {
		ParseAndCompare(t,
			"GRANT SELECT ON db.* TO user1 AS admin WITH ROLE ALL EXCEPT begin",
			"{GRANT :loc 0 :privileges SELECT :on {GRANT_TARGET :loc 16 :name {TABLEREF :loc 16 :schema db :name *}} :to user1 :as admin :with_role_type ALL EXCEPT :with_roles begin}",
		)
	})
}

// ============================================================================
// Batch 60: GRANT/REVOKE REQUIRE and resource limits
// ============================================================================

func TestParseGrantRequireSSL(t *testing.T) {
	t.Run("require ssl", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE SSL")
	})
	t.Run("require none", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE NONE")
	})
}

func TestParseGrantRequireX509(t *testing.T) {
	t.Run("require x509", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE X509")
	})
	t.Run("require cipher", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE CIPHER 'EDH-RSA-DES-CBC3-SHA'")
	})
	t.Run("require issuer", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE ISSUER '/C=SE/ST=Stockholm'")
	})
	t.Run("require subject", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE SUBJECT '/C=SE/ST=Stockholm'")
	})
	t.Run("require multiple tls options", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE ISSUER '/C=SE' AND SUBJECT '/C=SE' AND CIPHER 'EDH-RSA'")
	})
}

func TestParseGrantWithResourceLimits(t *testing.T) {
	t.Run("max queries per hour", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_QUERIES_PER_HOUR 500")
	})
	t.Run("max updates per hour", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_UPDATES_PER_HOUR 100")
	})
	t.Run("max connections per hour", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_CONNECTIONS_PER_HOUR 50")
	})
	t.Run("max user connections", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_USER_CONNECTIONS 10")
	})
	t.Run("multiple resource limits", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_QUERIES_PER_HOUR 500 MAX_UPDATES_PER_HOUR 100")
	})
	t.Run("grant option with resource limits", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH GRANT OPTION MAX_QUERIES_PER_HOUR 500")
	})
	t.Run("require ssl with resource limits", func(t *testing.T) {
		ParseAndCheck(t, "GRANT ALL ON *.* TO 'jeffrey'@'localhost' REQUIRE SSL WITH MAX_QUERIES_PER_HOUR 90")
	})
}

func TestParseCreateUserRequire(t *testing.T) {
	t.Run("require ssl", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY 'password' REQUIRE SSL")
	})
	t.Run("require x509", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' REQUIRE X509")
	})
	t.Run("require with resource limits", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' REQUIRE SSL WITH MAX_QUERIES_PER_HOUR 100")
	})
}

// ============================================================================
// Batch 61: ALTER USER extended options
// ============================================================================

func TestParseAlterUserPasswordExpire(t *testing.T) {
	t.Run("password expire", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD EXPIRE")
	})
	t.Run("password expire default", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD EXPIRE DEFAULT")
	})
	t.Run("password expire never", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD EXPIRE NEVER")
	})
	t.Run("password expire interval", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD EXPIRE INTERVAL 180 DAY")
	})
	t.Run("password history", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD HISTORY 5")
	})
	t.Run("password history default", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD HISTORY DEFAULT")
	})
	t.Run("password reuse interval", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD REUSE INTERVAL 365 DAY")
	})
	t.Run("password reuse interval default", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD REUSE INTERVAL DEFAULT")
	})
	t.Run("password require current", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD REQUIRE CURRENT")
	})
	t.Run("password require current optional", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD REQUIRE CURRENT OPTIONAL")
	})
}

func TestParseAlterUserAccountLock(t *testing.T) {
	t.Run("account lock", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' ACCOUNT LOCK")
	})
	t.Run("account unlock", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' ACCOUNT UNLOCK")
	})
}

func TestParseAlterUserFailedLoginAttempts(t *testing.T) {
	t.Run("failed login attempts", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' FAILED_LOGIN_ATTEMPTS 3")
	})
	t.Run("password lock time", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD_LOCK_TIME 2")
	})
	t.Run("password lock time unbounded", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' PASSWORD_LOCK_TIME UNBOUNDED")
	})
	t.Run("failed login with lock time", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 2")
	})
}

func TestParseAlterUserComment(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' COMMENT 'this is a test user'")
	})
	t.Run("attribute", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' ATTRIBUTE '{\"department\": \"engineering\"}'")
	})
}

func TestParseAlterUserRequireSSL(t *testing.T) {
	t.Run("require ssl with password expire", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' REQUIRE SSL PASSWORD EXPIRE INTERVAL 90 DAY")
	})
	t.Run("full alter user", func(t *testing.T) {
		ParseAndCheck(t, "ALTER USER 'jeffrey'@'localhost' IDENTIFIED BY 'newpass' REQUIRE X509 WITH MAX_QUERIES_PER_HOUR 100 PASSWORD EXPIRE INTERVAL 180 DAY ACCOUNT LOCK")
	})
}

func TestParseCreateUserPasswordExpire(t *testing.T) {
	t.Run("create user with password expire", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY 'password' PASSWORD EXPIRE INTERVAL 180 DAY")
	})
	t.Run("create user with account lock", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' ACCOUNT LOCK")
	})
	t.Run("create user with failed login", func(t *testing.T) {
		ParseAndCheck(t, "CREATE USER 'jeffrey'@'localhost' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 2")
	})
}

func TestParseShowGrants(t *testing.T) {
	tests := []string{
		"SHOW GRANTS",
		"SHOW GRANTS FOR root",
		"SHOW GRANTS FOR 'jeffrey'@'localhost'",
		"SHOW GRANTS FOR CURRENT_USER",
		"SHOW GRANTS FOR CURRENT_USER()",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowGrantsForUser(t *testing.T) {
	tests := []string{
		"SHOW GRANTS FOR 'u1'@'localhost' USING 'r1'",
		"SHOW GRANTS FOR 'u1'@'localhost' USING 'r1', 'r2'",
		"SHOW GRANTS FOR admin USING role1, role2, role3",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowIndex(t *testing.T) {
	tests := []string{
		// INDEX keyword
		"SHOW INDEX FROM users",
		"SHOW INDEX FROM users FROM mydb",
		"SHOW INDEX IN users",
		"SHOW INDEX IN users IN mydb",
		"SHOW INDEX FROM users WHERE Key_name = 'PRIMARY'",
		// INDEXES synonym (identifier)
		"SHOW INDEXES FROM users",
		"SHOW INDEXES IN users",
		// KEYS synonym
		"SHOW KEYS FROM users",
		"SHOW KEYS IN users IN mydb",
		// EXTENDED
		"SHOW EXTENDED INDEX FROM users",
		"SHOW EXTENDED INDEXES FROM users",
		"SHOW EXTENDED KEYS IN users",
		"SHOW EXTENDED INDEX FROM users FROM mydb",
		"SHOW EXTENDED INDEX FROM users WHERE Key_name = 'idx_name'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectDistinctrow(t *testing.T) {
	tests := []string{
		"SELECT DISTINCTROW a FROM t",
		"SELECT DISTINCTROW a, b, c FROM t WHERE a > 1",
		"SELECT DISTINCTROW * FROM t ORDER BY a",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateLoadableFunction(t *testing.T) {
	tests := []string{
		"CREATE FUNCTION myfunc RETURNS STRING SONAME 'mylib.so'",
		"CREATE FUNCTION myfunc RETURNS INTEGER SONAME 'mylib.so'",
		"CREATE FUNCTION myfunc RETURNS REAL SONAME 'mylib.so'",
		"CREATE FUNCTION myfunc RETURNS DECIMAL SONAME 'mylib.so'",
		"CREATE FUNCTION IF NOT EXISTS myfunc RETURNS STRING SONAME 'mylib.so'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateAggregateFunction(t *testing.T) {
	tests := []string{
		"CREATE AGGREGATE FUNCTION myagg RETURNS STRING SONAME 'myagg.so'",
		"CREATE AGGREGATE FUNCTION myagg RETURNS INTEGER SONAME 'myagg.so'",
		"CREATE AGGREGATE FUNCTION IF NOT EXISTS myagg RETURNS REAL SONAME 'myagg.so'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 65: depth_fix_select_into_outfile ---

func TestParseSelectIntoOutfileFields(t *testing.T) {
	tests := []string{
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' FIELDS TERMINATED BY ','",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '\"'",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' ESCAPED BY '\\\\'",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' COLUMNS TERMINATED BY '\\t' ENCLOSED BY '\"'",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' FIELDS TERMINATED BY ',' LINES TERMINATED BY '\\n'",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' FIELDS TERMINATED BY '|' ENCLOSED BY '\"' ESCAPED BY '\\\\' LINES STARTING BY 'X' TERMINATED BY '\\n'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectIntoOutfileCharset(t *testing.T) {
	tests := []string{
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' CHARACTER SET utf8",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' CHARACTER SET utf8mb4 FIELDS TERMINATED BY ','",
		"SELECT * FROM t INTO OUTFILE '/tmp/data.csv' CHARSET latin1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectIntoBeforeFrom(t *testing.T) {
	tests := []string{
		"SELECT * INTO OUTFILE '/tmp/data.csv' FROM t",
		"SELECT * INTO DUMPFILE '/tmp/data.bin' FROM t",
		"SELECT a, b INTO @x, @y FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseLockInShareMode(t *testing.T) {
	tests := []string{
		"SELECT * FROM t LOCK IN SHARE MODE",
		"SELECT * FROM t WHERE id = 1 LOCK IN SHARE MODE",
		"SELECT * FROM t ORDER BY id LOCK IN SHARE MODE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectSqlSmallResult(t *testing.T) {
	tests := []string{
		"SELECT SQL_SMALL_RESULT * FROM t GROUP BY c",
		"SELECT SQL_BIG_RESULT * FROM t GROUP BY c",
		"SELECT SQL_SMALL_RESULT SQL_BIG_RESULT * FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectSqlBufferResult(t *testing.T) {
	tests := []string{
		"SELECT SQL_BUFFER_RESULT * FROM t",
		"SELECT SQL_NO_CACHE * FROM t",
		"SELECT SQL_BUFFER_RESULT SQL_NO_CACHE * FROM t",
		"SELECT HIGH_PRIORITY SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_NO_CACHE SQL_CALC_FOUND_ROWS * FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 66: depth_fix_select_partition_with_insert ---

func TestParseSelectPartition(t *testing.T) {
	tests := []string{
		"SELECT * FROM t PARTITION (p0)",
		"SELECT * FROM t PARTITION (p0, p1, p2)",
		"SELECT * FROM t PARTITION (p0) WHERE id = 1",
		"SELECT * FROM t PARTITION (p0, p1) AS a",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWithInsert(t *testing.T) {
	tests := []string{
		"WITH cte AS (SELECT 1 AS x) INSERT INTO t SELECT * FROM cte",
		"WITH cte AS (SELECT 1 AS x) REPLACE INTO t SELECT * FROM cte",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseIntersectPrecedence(t *testing.T) {
	tests := []string{
		// INTERSECT has higher precedence than UNION
		"SELECT 1 UNION SELECT 2 INTERSECT SELECT 3",
		"SELECT 1 INTERSECT SELECT 2 UNION SELECT 3",
		"SELECT 1 INTERSECT ALL SELECT 2",
		"SELECT 1 EXCEPT SELECT 2 INTERSECT SELECT 3",
		"SELECT 1 INTERSECT SELECT 2 INTERSECT SELECT 3",
		"SELECT 1 UNION SELECT 2 EXCEPT SELECT 3",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 67: depth_fix_create_table_constraints_ctas ---

func TestParseCreateTableConstraintIndexOption(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT, PRIMARY KEY (id) USING BTREE)",
		"CREATE TABLE t (id INT, PRIMARY KEY (id) KEY_BLOCK_SIZE 1024)",
		"CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT 'pk')",
		"CREATE TABLE t (id INT, UNIQUE KEY idx (name) USING HASH COMMENT 'uniq')",
		"CREATE TABLE t (id INT, INDEX idx (name) KEY_BLOCK_SIZE 512 COMMENT 'idx')",
		"CREATE TABLE t (id INT, FULLTEXT KEY ft_idx (content) COMMENT 'fulltext')",
		"CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE)",
		"CREATE TABLE t (id INT, UNIQUE KEY idx (name) INVISIBLE)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableIgnoreAsSelect(t *testing.T) {
	tests := []string{
		"CREATE TABLE t IGNORE SELECT * FROM s",
		"CREATE TABLE t REPLACE SELECT * FROM s",
		"CREATE TABLE t IGNORE AS SELECT * FROM s",
		"CREATE TABLE t REPLACE AS SELECT * FROM s",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableBareSelect(t *testing.T) {
	tests := []string{
		"CREATE TABLE t SELECT * FROM s",
		"CREATE TABLE t (id INT) SELECT * FROM s",
		"CREATE TABLE t ENGINE=InnoDB SELECT * FROM s",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableLinearHash(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) PARTITION BY LINEAR HASH(id) PARTITIONS 4",
		"CREATE TABLE t (id INT) PARTITION BY LINEAR KEY(id) PARTITIONS 4",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableFKMatch(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT, FOREIGN KEY (pid) REFERENCES p(id) MATCH FULL)",
		"CREATE TABLE t (id INT, FOREIGN KEY (pid) REFERENCES p(id) MATCH PARTIAL ON DELETE CASCADE)",
		"CREATE TABLE t (id INT, FOREIGN KEY (pid) REFERENCES p(id) MATCH SIMPLE ON DELETE SET NULL ON UPDATE CASCADE)",
		"CREATE TABLE t (id INT, pid INT REFERENCES p(id) MATCH FULL)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 68: depth_fix_alter_table_default_constraint ---

func TestParseAlterColumnSetDefault(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ALTER COLUMN x SET DEFAULT 0",
		"ALTER TABLE t ALTER COLUMN x SET DEFAULT 'hello'",
		"ALTER TABLE t ALTER COLUMN x SET DEFAULT (NOW())",
		"ALTER TABLE t ALTER x SET DEFAULT 42",
		"ALTER TABLE t ALTER COLUMN x DROP DEFAULT",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterTableAddMultiColumn(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ADD (a INT, b VARCHAR(100))",
		"ALTER TABLE t ADD (x INT)",
		"ALTER TABLE t ADD COLUMN (a INT, b INT, c INT)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterCheckEnforced(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ALTER CHECK chk1 ENFORCED",
		"ALTER TABLE t ALTER CHECK chk1 NOT ENFORCED",
		"ALTER TABLE t ALTER CONSTRAINT con1 ENFORCED",
		"ALTER TABLE t ALTER CONSTRAINT con1 NOT ENFORCED",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterTableDropConstraint(t *testing.T) {
	tests := []string{
		"ALTER TABLE t DROP CONSTRAINT con1",
		"ALTER TABLE t DROP CHECK chk1",
		"ALTER TABLE t DROP FOREIGN KEY fk1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 72: depth_fix_expr_group_concat ---

func TestParseGroupConcatOrderByExpr(t *testing.T) {
	tests := []string{
		"SELECT GROUP_CONCAT(name ORDER BY id) FROM t",
		"SELECT GROUP_CONCAT(name ORDER BY id DESC) FROM t",
		"SELECT GROUP_CONCAT(name ORDER BY UPPER(name)) FROM t",
		"SELECT GROUP_CONCAT(name ORDER BY id + 1 DESC) FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseGroupConcatOrderByFuncCall(t *testing.T) {
	tests := []string{
		"SELECT GROUP_CONCAT(name ORDER BY CONCAT(a, b)) FROM t",
		"SELECT GROUP_CONCAT(DISTINCT name ORDER BY name SEPARATOR ',') FROM t",
		"SELECT GROUP_CONCAT(name ORDER BY COALESCE(x, 0) DESC, name ASC SEPARATOR '; ') FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 71: depth_fix_analyze_histogram_flush ---

func TestParseAnalyzeTableUpdateHistogram(t *testing.T) {
	tests := []string{
		"ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1 WITH 10 BUCKETS",
		"ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1, col2 WITH 256 BUCKETS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAnalyzeTableDropHistogram(t *testing.T) {
	tests := []string{
		"ANALYZE TABLE t1 DROP HISTOGRAM ON col1",
		"ANALYZE TABLE t1 DROP HISTOGRAM ON col1, col2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseFlushTablesWithReadLock(t *testing.T) {
	tests := []string{
		"FLUSH TABLES WITH READ LOCK",
		"FLUSH TABLES t1, t2 WITH READ LOCK",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseFlushTablesForExport(t *testing.T) {
	tests := []string{
		"FLUSH TABLES t1 FOR EXPORT",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseFlushBinaryLogs(t *testing.T) {
	tests := []string{
		"FLUSH BINARY LOGS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUndoTablespaceSetActive(t *testing.T) {
	tests := []string{
		"ALTER UNDO TABLESPACE myts SET ACTIVE",
		"ALTER UNDO TABLESPACE myts SET INACTIVE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 70: depth_fix_set_persist_show_full ---

func TestParseSetPersist(t *testing.T) {
	tests := []string{
		"SET PERSIST max_connections = 100",
		"SET PERSIST innodb_buffer_pool_size = 1073741824",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSetPersistOnly(t *testing.T) {
	tests := []string{
		"SET PERSIST_ONLY back_log = 100",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowFullTables(t *testing.T) {
	tests := []string{
		"SHOW FULL TABLES",
		"SHOW FULL TABLES FROM mydb",
		"SHOW FULL TABLES LIKE 'test%'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowSchemas(t *testing.T) {
	tests := []string{
		"SHOW SCHEMAS",
		"SHOW SCHEMAS LIKE 'my%'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowFields(t *testing.T) {
	tests := []string{
		"SHOW FIELDS FROM t1",
		"SHOW FULL FIELDS FROM t1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowExtendedColumns(t *testing.T) {
	tests := []string{
		"SHOW EXTENDED COLUMNS FROM t1",
		"SHOW EXTENDED FULL COLUMNS FROM t1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseExplainForConnection(t *testing.T) {
	tests := []string{
		"EXPLAIN FOR CONNECTION 42",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 69: depth_fix_grant_proxy_privileges ---

func TestParseGrantProxy(t *testing.T) {
	tests := []string{
		"GRANT PROXY ON 'admin'@'localhost' TO 'user1'@'localhost'",
		"GRANT PROXY ON 'admin'@'localhost' TO 'user1'@'localhost' WITH GRANT OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRevokeProxy(t *testing.T) {
	tests := []string{
		"REVOKE PROXY ON 'admin'@'localhost' FROM 'user1'@'localhost'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseGrantColumnPrivilege(t *testing.T) {
	tests := []string{
		"GRANT SELECT (col1, col2), INSERT (col1) ON mydb.mytbl TO 'user1'@'localhost'",
		"GRANT UPDATE (col1) ON mytbl TO 'user1'@'localhost'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseGrantMultiWordPrivilege(t *testing.T) {
	tests := []string{
		"GRANT CREATE VIEW ON mydb.* TO 'user1'@'localhost'",
		"GRANT ALTER ROUTINE ON mydb.* TO 'user1'@'localhost'",
		"GRANT SHOW DATABASES ON *.* TO 'user1'@'localhost'",
		"GRANT CREATE TEMPORARY TABLES ON mydb.* TO 'user1'@'localhost'",
		"GRANT LOCK TABLES ON mydb.* TO 'user1'@'localhost'",
		"GRANT REPLICATION CLIENT ON *.* TO 'user1'@'localhost'",
		"GRANT REPLICATION SLAVE ON *.* TO 'user1'@'localhost'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRevokeAllPrivileges(t *testing.T) {
	tests := []string{
		"REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'user1'@'localhost'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseIdentifiedWithAs(t *testing.T) {
	tests := []string{
		"CREATE USER 'user1'@'localhost' IDENTIFIED WITH mysql_native_password AS '*hash123'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUserIfExists(t *testing.T) {
	tests := []string{
		"ALTER USER IF EXISTS 'user1'@'localhost' IDENTIFIED BY 'newpass'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUserDefaultRole(t *testing.T) {
	tests := []string{
		"ALTER USER 'user1'@'localhost' DEFAULT ROLE ALL",
		"ALTER USER 'user1'@'localhost' DEFAULT ROLE NONE",
		"ALTER USER 'user1'@'localhost' DEFAULT ROLE role1, role2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 73: depth_fix_insert_update_delete_partition ---

func TestParseInsertPartition(t *testing.T) {
	tests := []string{
		"INSERT INTO t1 PARTITION (p0) VALUES (1, 2)",
		"INSERT INTO t1 PARTITION (p0, p1, p2) (col1, col2) VALUES (1, 2), (3, 4)",
		"REPLACE INTO t1 PARTITION (p0) VALUES (1)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseInsertAsRowAlias(t *testing.T) {
	tests := []string{
		"INSERT INTO t1 VALUES (1, 2) AS new ON DUPLICATE KEY UPDATE col1 = new.col1",
		"INSERT INTO t1 VALUES (1, 2) AS new(a, b) ON DUPLICATE KEY UPDATE col1 = new.a",
		"INSERT INTO t1 SET col1 = 1 AS new ON DUPLICATE KEY UPDATE col1 = new.col1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDeletePartition(t *testing.T) {
	tests := []string{
		"DELETE FROM t1 PARTITION (p0) WHERE id = 1",
		"DELETE FROM t1 PARTITION (p0, p1) WHERE id > 10",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDeleteTableAlias(t *testing.T) {
	tests := []string{
		"DELETE FROM t1 AS a WHERE a.id = 1",
		"DELETE FROM t1 a WHERE a.id = 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWithUpdate(t *testing.T) {
	tests := []string{
		"WITH cte AS (SELECT id, val FROM t2) UPDATE t1 JOIN cte ON t1.id = cte.id SET t1.val = cte.val",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseWithDelete(t *testing.T) {
	tests := []string{
		"WITH cte AS (SELECT id FROM t2) DELETE FROM t1 WHERE id IN (SELECT id FROM cte)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Batch 74 – depth_fix_replication_misc
// -----------------------------------------------------------------------

func TestParsePurgeBinaryLogsBeforeDatetime(t *testing.T) {
	tests := []string{
		"PURGE BINARY LOGS BEFORE '2023-01-01 00:00:00'",
		"PURGE BINARY LOGS BEFORE NOW() - INTERVAL 3 DAY",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCacheIndexPartition(t *testing.T) {
	tests := []string{
		"CACHE INDEX t1 PARTITION (p0, p1) IN hot_cache",
		"CACHE INDEX t1 PARTITION ALL IN hot_cache",
		"CACHE INDEX t1 INDEX (idx1, idx2) IN hot_cache",
		"CACHE INDEX t1 KEY (idx1) PARTITION (p0) IN hot_cache",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseShowProcesslist(t *testing.T) {
	tests := []string{
		"SHOW PROCESSLIST",
		"SHOW FULL PROCESSLIST",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Batch 75: SOUNDS LIKE, ROW(), DEFAULT(col), VALUES(col)
// ---------------------------------------------------------------------------

func TestParseSoundsLike(t *testing.T) {
	tests := []string{
		// Basic SOUNDS LIKE
		"SELECT * FROM t WHERE name SOUNDS LIKE 'John'",
		// SOUNDS LIKE with column references
		"SELECT * FROM t WHERE a SOUNDS LIKE b",
		// SOUNDS LIKE in complex expression
		"SELECT * FROM t WHERE name SOUNDS LIKE 'John' AND age > 20",
		// SOUNDS LIKE with function
		"SELECT * FROM t WHERE CONCAT(first, last) SOUNDS LIKE 'JohnDoe'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRowConstructor(t *testing.T) {
	tests := []string{
		// Basic ROW() constructor
		"SELECT ROW(1, 2, 3)",
		// ROW in WHERE clause comparison
		"SELECT * FROM t WHERE ROW(a, b) = ROW(1, 2)",
		// ROW with expressions
		"SELECT * FROM t WHERE ROW(a + 1, b * 2) = ROW(10, 20)",
		// ROW with single value (MySQL requires at least 2 for row comparisons but ROW(1) is valid syntax)
		"SELECT ROW(1)",
		// Nested ROW is valid syntax
		"SELECT ROW(1, ROW(2, 3))",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRowComparison(t *testing.T) {
	tests := []string{
		// ROW comparison with =
		"SELECT * FROM t1 WHERE ROW(col1, col2) = (SELECT col3, col4 FROM t2 WHERE id = 10)",
		// ROW comparison with <>
		"SELECT * FROM t WHERE ROW(a, b) <> ROW(1, 2)",
		// ROW comparison with <=
		"SELECT * FROM t WHERE ROW(a, b) <= ROW(1, 2)",
		// ROW comparison with >=
		"SELECT * FROM t WHERE ROW(a, b) >= ROW(1, 2)",
		// ROW comparison with <
		"SELECT * FROM t WHERE ROW(a, b) < ROW(1, 2)",
		// ROW comparison with >
		"SELECT * FROM t WHERE ROW(a, b) > ROW(1, 2)",
		// ROW comparison with <=>
		"SELECT * FROM t WHERE ROW(a, b) <=> ROW(1, 2)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDefaultColName(t *testing.T) {
	tests := []string{
		// Bare DEFAULT in INSERT
		"INSERT INTO t VALUES (DEFAULT)",
		// DEFAULT(col_name) function form in UPDATE
		"UPDATE t SET i = DEFAULT(i) + 1 WHERE id < 100",
		// DEFAULT(col) in SELECT
		"SELECT DEFAULT(name) FROM t",
		// Multiple DEFAULT(col)
		"UPDATE t SET a = DEFAULT(a), b = DEFAULT(b)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseValuesFunc(t *testing.T) {
	tests := []string{
		// VALUES(col) in INSERT ON DUPLICATE KEY UPDATE
		"INSERT INTO t (a, b, c) VALUES (1, 2, 3) ON DUPLICATE KEY UPDATE c = VALUES(a) + VALUES(b)",
		// VALUES(col) simple usage
		"INSERT INTO t (a) VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseMatchWithQueryExpansion(t *testing.T) {
	tests := []string{
		// Standalone WITH QUERY EXPANSION modifier
		"SELECT * FROM articles WHERE MATCH (title) AGAINST ('database' WITH QUERY EXPANSION)",
		// WITH QUERY EXPANSION with multiple columns
		"SELECT * FROM articles WHERE MATCH (title, body) AGAINST ('MySQL' WITH QUERY EXPANSION)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseMatchNaturalLanguageWithQueryExpansion(t *testing.T) {
	tests := []string{
		// Combined IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION
		"SELECT * FROM articles WHERE MATCH (title) AGAINST ('database' IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION)",
		// Combined modifier with multiple columns
		"SELECT * FROM articles WHERE MATCH (title, body) AGAINST ('MySQL' IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION)",
		// Basic IN NATURAL LANGUAGE MODE (regression check)
		"SELECT * FROM articles WHERE MATCH (title) AGAINST ('database' IN NATURAL LANGUAGE MODE)",
		// IN BOOLEAN MODE (regression check)
		"SELECT * FROM articles WHERE MATCH (title) AGAINST ('+MySQL -YourSQL' IN BOOLEAN MODE)",
		// No modifier (regression check)
		"SELECT * FROM articles WHERE MATCH (title, body) AGAINST ('database')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseNaturalLeftJoin(t *testing.T) {
	tests := []string{
		"SELECT * FROM t1 NATURAL LEFT JOIN t2",
		"SELECT * FROM t1 NATURAL LEFT JOIN t2 WHERE t1.id > 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseNaturalRightJoin(t *testing.T) {
	tests := []string{
		"SELECT * FROM t1 NATURAL RIGHT JOIN t2",
		"SELECT * FROM t1 NATURAL RIGHT JOIN t2 WHERE t1.id > 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseNaturalLeftOuterJoin(t *testing.T) {
	tests := []string{
		"SELECT * FROM t1 NATURAL LEFT OUTER JOIN t2",
		"SELECT * FROM t1 NATURAL RIGHT OUTER JOIN t2",
		// Plain NATURAL JOIN (regression check)
		"SELECT * FROM t1 NATURAL JOIN t2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectIntoAfterHaving(t *testing.T) {
	tests := []string{
		// INTO at position 2: after HAVING, before ORDER BY
		"SELECT a, COUNT(*) AS cnt FROM t GROUP BY a HAVING cnt > 1 INTO OUTFILE '/tmp/result.txt'",
		// INTO at position 2: after HAVING, before ORDER BY with ORDER BY
		"SELECT a, COUNT(*) AS cnt FROM t GROUP BY a HAVING cnt > 1 INTO OUTFILE '/tmp/result.txt' ORDER BY a",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseSelectIntoAfterHavingBeforeOrderBy(t *testing.T) {
	tests := []string{
		// INTO with DUMPFILE at position 2
		"SELECT a FROM t GROUP BY a HAVING a > 0 INTO DUMPFILE '/tmp/dump.dat' ORDER BY a",
		// INTO with variable at position 2
		"SELECT a, b FROM t GROUP BY a HAVING a > 0 INTO @x, @y ORDER BY a",
		// Regression: INTO at position 1 still works
		"SELECT a INTO @x FROM t",
		// Regression: INTO at position 3 still works
		"SELECT a FROM t FOR UPDATE INTO OUTFILE '/tmp/result.txt'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateUserRandomPassword(t *testing.T) {
	tests := []string{
		// IDENTIFIED BY RANDOM PASSWORD
		"CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY RANDOM PASSWORD",
		// IDENTIFIED WITH plugin BY RANDOM PASSWORD
		"CREATE USER 'jeffrey'@'localhost' IDENTIFIED WITH caching_sha2_password BY RANDOM PASSWORD",
		// Multiple users with RANDOM PASSWORD
		"CREATE USER 'u1'@'localhost' IDENTIFIED BY RANDOM PASSWORD, 'u2'@'localhost' IDENTIFIED BY RANDOM PASSWORD",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUserRetainCurrentPassword(t *testing.T) {
	tests := []string{
		// RETAIN CURRENT PASSWORD
		"ALTER USER 'jeffrey'@'localhost' IDENTIFIED BY 'new_password' RETAIN CURRENT PASSWORD",
		// RETAIN CURRENT PASSWORD with auth plugin
		"ALTER USER 'jeffrey'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'new_password' RETAIN CURRENT PASSWORD",
		// RANDOM PASSWORD with RETAIN CURRENT PASSWORD
		"ALTER USER 'jeffrey'@'localhost' IDENTIFIED BY RANDOM PASSWORD RETAIN CURRENT PASSWORD",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterUserDiscardOldPassword(t *testing.T) {
	tests := []string{
		// DISCARD OLD PASSWORD
		"ALTER USER 'jeffrey'@'localhost' DISCARD OLD PASSWORD",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTablePassword(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) PASSWORD = 'secret'",
		"CREATE TABLE t (id INT) PASSWORD 'secret'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableSecondaryEngine(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) SECONDARY_ENGINE = RAPID",
		"CREATE TABLE t (id INT) SECONDARY_ENGINE RAPID",
		"CREATE TABLE t (id INT) ENGINE = InnoDB SECONDARY_ENGINE = RAPID",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableSecondaryEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE = '{\"key\": \"value\"}'",
		"CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE '{\"key\": \"value\"}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseColumnFormat(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT COLUMN_FORMAT FIXED)",
		"CREATE TABLE t (id INT COLUMN_FORMAT DYNAMIC)",
		"CREATE TABLE t (id INT COLUMN_FORMAT DEFAULT)",
		"CREATE TABLE t (id INT NOT NULL COLUMN_FORMAT FIXED AUTO_INCREMENT)",
		"CREATE TABLE t (a INT COLUMN_FORMAT FIXED, b VARCHAR(100) COLUMN_FORMAT DYNAMIC)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseColumnStorage(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT STORAGE DISK)",
		"CREATE TABLE t (id INT STORAGE MEMORY)",
		"CREATE TABLE t (id INT NOT NULL STORAGE DISK AUTO_INCREMENT)",
		"CREATE TABLE t (a INT STORAGE DISK, b VARCHAR(100) STORAGE MEMORY)",
		"CREATE TABLE t (id INT COLUMN_FORMAT FIXED STORAGE DISK)",
		"CREATE TABLE t (id INT STORAGE MEMORY COLUMN_FORMAT DYNAMIC)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ─── Batch 82: depth_fix_utility_tests_batch1 ─────────────────────────────

func TestParseAnalyzeTable(t *testing.T) {
	tests := []string{
		"ANALYZE TABLE t1",
		"ANALYZE TABLE t1, t2, t3",
		"ANALYZE TABLE db1.t1",
		"ANALYZE LOCAL TABLE t1",
		"ANALYZE NO_WRITE_TO_BINLOG TABLE t1",
		"ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1 WITH 10 BUCKETS",
		"ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1, col2 WITH 50 BUCKETS",
		"ANALYZE TABLE t1 DROP HISTOGRAM ON col1",
		"ANALYZE TABLE t1 DROP HISTOGRAM ON col1, col2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseOptimizeTable(t *testing.T) {
	tests := []string{
		"OPTIMIZE TABLE t1",
		"OPTIMIZE TABLE t1, t2",
		"OPTIMIZE TABLE db1.t1",
		"OPTIMIZE LOCAL TABLE t1",
		"OPTIMIZE NO_WRITE_TO_BINLOG TABLE t1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCheckTable(t *testing.T) {
	tests := []string{
		"CHECK TABLE t1",
		"CHECK TABLE t1, t2",
		"CHECK TABLE db1.t1",
		"CHECK TABLE t1 FOR UPGRADE",
		"CHECK TABLE t1 QUICK",
		"CHECK TABLE t1 FAST",
		"CHECK TABLE t1 MEDIUM",
		"CHECK TABLE t1 EXTENDED",
		"CHECK TABLE t1 CHANGED",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRepairTable(t *testing.T) {
	tests := []string{
		"REPAIR TABLE t1",
		"REPAIR TABLE t1, t2",
		"REPAIR TABLE db1.t1",
		"REPAIR LOCAL TABLE t1",
		"REPAIR NO_WRITE_TO_BINLOG TABLE t1",
		"REPAIR TABLE t1 QUICK",
		"REPAIR TABLE t1 EXTENDED",
		"REPAIR TABLE t1 USE_FRM",
		"REPAIR TABLE t1 QUICK EXTENDED",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseFlush(t *testing.T) {
	tests := []string{
		"FLUSH TABLES",
		"FLUSH TABLES t1",
		"FLUSH TABLES t1, t2",
		"FLUSH TABLES WITH READ LOCK",
		"FLUSH TABLES t1, t2 WITH READ LOCK",
		"FLUSH TABLES t1 FOR EXPORT",
		"FLUSH LOCAL TABLES",
		"FLUSH NO_WRITE_TO_BINLOG TABLES",
		"FLUSH PRIVILEGES",
		"FLUSH STATUS",
		"FLUSH BINARY LOGS",
		"FLUSH HOSTS",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseReset(t *testing.T) {
	tests := []string{
		"RESET QUERY CACHE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseKill(t *testing.T) {
	tests := []string{
		"KILL 123",
		"KILL QUERY 123",
		"KILL CONNECTION 123",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseDo(t *testing.T) {
	tests := []string{
		"DO 1",
		"DO 1 + 2",
		"DO 1 + 2, 3 + 4",
		"DO SLEEP(1)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestParseChecksumTable, TestParseShutdown, TestParseRestart are defined
// earlier in batch 34 with ParseAndCompare assertions.

// ─── Batch 83: depth_fix_utility_tests_batch2 ─────────────────────────────

func TestParseAlterServer(t *testing.T) {
	tests := []string{
		"ALTER SERVER s1 OPTIONS (USER 'admin')",
		"ALTER SERVER s1 OPTIONS (HOST '10.0.0.1', DATABASE 'db1')",
		"ALTER SERVER s1 OPTIONS (USER 'root', HOST 'localhost', DATABASE 'test', PORT 3306)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseHandlerOpenReadClose(t *testing.T) {
	tests := []string{
		"HANDLER t1 OPEN",
		"HANDLER t1 OPEN AS h1",
		"HANDLER h1 READ FIRST",
		"HANDLER h1 READ NEXT",
		"HANDLER h1 READ idx1 = (1)",
		"HANDLER h1 READ idx1 FIRST",
		"HANDLER h1 READ idx1 NEXT",
		"HANDLER h1 READ idx1 PREV",
		"HANDLER h1 READ idx1 LAST",
		"HANDLER h1 READ NEXT WHERE id > 10",
		"HANDLER h1 READ NEXT LIMIT 5",
		"HANDLER h1 CLOSE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseLockUnlockInstance(t *testing.T) {
	tests := []string{
		"LOCK INSTANCE FOR BACKUP",
		"UNLOCK INSTANCE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseIndexHints(t *testing.T) {
	tests := []string{
		// USE INDEX with index list
		"SELECT * FROM t1 USE INDEX (idx1)",
		"SELECT * FROM t1 USE INDEX (idx1, idx2)",
		// USE KEY (synonym)
		"SELECT * FROM t1 USE KEY (idx1)",
		// USE INDEX with empty list (use no indexes)
		"SELECT * FROM t1 USE INDEX ()",
		// USE INDEX with FOR clause
		"SELECT * FROM t1 USE INDEX FOR JOIN (idx1)",
		"SELECT * FROM t1 USE INDEX FOR ORDER BY (idx1)",
		"SELECT * FROM t1 USE INDEX FOR GROUP BY (idx1)",
		// FORCE INDEX
		"SELECT * FROM t1 FORCE INDEX (idx1)",
		"SELECT * FROM t1 FORCE KEY (idx1, idx2)",
		// FORCE INDEX with FOR clause
		"SELECT * FROM t1 FORCE INDEX FOR JOIN (idx1)",
		"SELECT * FROM t1 FORCE INDEX FOR ORDER BY (idx1)",
		"SELECT * FROM t1 FORCE INDEX FOR GROUP BY (idx1)",
		// IGNORE INDEX
		"SELECT * FROM t1 IGNORE INDEX (idx1)",
		"SELECT * FROM t1 IGNORE KEY (idx1, idx2)",
		// IGNORE INDEX with FOR clause
		"SELECT * FROM t1 IGNORE INDEX FOR JOIN (idx1)",
		"SELECT * FROM t1 IGNORE INDEX FOR ORDER BY (idx1)",
		"SELECT * FROM t1 IGNORE INDEX FOR GROUP BY (idx1)",
		// Multiple index hints
		"SELECT * FROM t1 USE INDEX (i1) IGNORE INDEX FOR ORDER BY (i2) ORDER BY a",
		// Index hints with alias
		"SELECT * FROM t1 AS a USE INDEX (idx1)",
		// Index hints with schema-qualified table
		"SELECT * FROM db1.t1 USE INDEX (idx1)",
		// Index hints in JOIN
		"SELECT * FROM t1 USE INDEX (i1) JOIN t2 FORCE INDEX (i2) ON t1.id = t2.id",
		// Multiple FORCE INDEX hints
		"SELECT * FROM t1 FORCE INDEX (i1) FORCE INDEX FOR ORDER BY (i2)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterInstance(t *testing.T) {
	tests := []string{
		"ALTER INSTANCE ROTATE INNODB MASTER KEY",
		"ALTER INSTANCE ENABLE INNODB REDO_LOG",
		"ALTER INSTANCE DISABLE INNODB REDO_LOG",
		"ALTER INSTANCE RELOAD TLS",
		"ALTER INSTANCE RELOAD TLS NO ROLLBACK ON ERROR",
		"ALTER INSTANCE RELOAD TLS FOR CHANNEL mysql_main",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// --- Batch 85: depth_fix_show_event_value_capture ---

func TestParseShowBinlogEventsFromLimit(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			sql:      "SHOW BINLOG EVENTS",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS IN 'binlog.000001'",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :like {STRING_LIT :val "binlog.000001" :loc 22}}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS FROM 4",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :from_pos {INT_LIT :val 4 :loc 24}}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS IN 'binlog.000001' FROM 4",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :like {STRING_LIT :val "binlog.000001" :loc 22} :from_pos {INT_LIT :val 4 :loc 43}}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS LIMIT 10",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :limit_count {INT_LIT :val 10 :loc 25}}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS LIMIT 5, 10",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :limit_count {INT_LIT :val 10 :loc 28} :limit_offset {INT_LIT :val 5 :loc 25}}`,
		},
		{
			sql:      "SHOW BINLOG EVENTS IN 'binlog.000001' FROM 4 LIMIT 5, 10",
			expected: `{SHOW :loc 0 :type BINLOG EVENTS :like {STRING_LIT :val "binlog.000001" :loc 22} :from_pos {INT_LIT :val 4 :loc 43} :limit_count {INT_LIT :val 10 :loc 54} :limit_offset {INT_LIT :val 5 :loc 51}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

func TestParseShowRelaylogEventsFromLimit(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			sql:      "SHOW RELAYLOG EVENTS",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS}`,
		},
		{
			sql:      "SHOW RELAYLOG EVENTS IN 'relay-bin.000001'",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS :like {STRING_LIT :val "relay-bin.000001" :loc 24}}`,
		},
		{
			sql:      "SHOW RELAYLOG EVENTS FROM 4",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS :from_pos {INT_LIT :val 4 :loc 26}}`,
		},
		{
			sql:      "SHOW RELAYLOG EVENTS LIMIT 10",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS :limit_count {INT_LIT :val 10 :loc 27}}`,
		},
		{
			sql:      "SHOW RELAYLOG EVENTS LIMIT 5, 10",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS :limit_count {INT_LIT :val 10 :loc 30} :limit_offset {INT_LIT :val 5 :loc 27}}`,
		},
		{
			sql:      "SHOW RELAYLOG EVENTS IN 'relay-bin.000001' FROM 100 LIMIT 20",
			expected: `{SHOW :loc 0 :type RELAYLOG EVENTS :like {STRING_LIT :val "relay-bin.000001" :loc 24} :from_pos {INT_LIT :val 100 :loc 48} :limit_count {INT_LIT :val 20 :loc 58}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

func TestParseShowWarningsLimit(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			sql:      "SHOW WARNINGS",
			expected: `{SHOW :loc 0 :type WARNINGS}`,
		},
		{
			sql:      "SHOW WARNINGS LIMIT 10",
			expected: `{SHOW :loc 0 :type WARNINGS :limit_count {INT_LIT :val 10 :loc 20}}`,
		},
		{
			sql:      "SHOW WARNINGS LIMIT 5, 10",
			expected: `{SHOW :loc 0 :type WARNINGS :limit_count {INT_LIT :val 10 :loc 23} :limit_offset {INT_LIT :val 5 :loc 20}}`,
		},
		{
			sql:      "SHOW ERRORS",
			expected: `{SHOW :loc 0 :type ERRORS}`,
		},
		{
			sql:      "SHOW ERRORS LIMIT 10",
			expected: `{SHOW :loc 0 :type ERRORS :limit_count {INT_LIT :val 10 :loc 18}}`,
		},
		{
			sql:      "SHOW ERRORS LIMIT 5, 10",
			expected: `{SHOW :loc 0 :type ERRORS :limit_count {INT_LIT :val 10 :loc 21} :limit_offset {INT_LIT :val 5 :loc 18}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

// --- Batch 86: depth_fix_show_create_user_host ---

func TestParseShowCreateUserHost(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{
			sql:      "SHOW CREATE USER root",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name root}}`,
		},
		{
			sql:      "SHOW CREATE USER 'root'@'localhost'",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name root :host localhost}}`,
		},
		{
			sql:      "SHOW CREATE USER 'admin'@'%'",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name admin :host %}}`,
		},
		{
			sql:      "SHOW CREATE USER 'user1'@'192.168.1.1'",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name user1 :host 192.168.1.1}}`,
		},
		{
			sql:      "SHOW CREATE USER CURRENT_USER",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name CURRENT_USER}}`,
		},
		{
			sql:      "SHOW CREATE USER CURRENT_USER()",
			expected: `{SHOW :loc 0 :type CREATE USER :for_user {USER_SPEC :loc 17 :name CURRENT_USER}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.expected)
		})
	}
}

// --- Batch 87: depth_fix_partition_algorithm_value ---

func TestParsePartitionKeyAlgorithm(t *testing.T) {
	tests := []string{
		"CREATE TABLE t1 (id INT) PARTITION BY KEY ALGORITHM=1 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY KEY ALGORITHM=2 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY LINEAR KEY ALGORITHM=1 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY KEY (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY KEY ALGORITHM=1 (id) SUBPARTITION BY KEY ALGORITHM=2 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY HASH(id) SUBPARTITION BY KEY ALGORITHM=1 (id)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify algorithm value is captured in AST
	t.Run("algorithm_value_captured", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t1 (id INT) PARTITION BY KEY ALGORITHM=2 (id)")
		s := ast.NodeToString(result.Items[0])
		if !strings.Contains(s, ":algorithm 2") {
			t.Errorf("expected :algorithm 2 in AST, got: %s", s)
		}
	})
	t.Run("no_algorithm_when_unset", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t1 (id INT) PARTITION BY KEY (id)")
		s := ast.NodeToString(result.Items[0])
		if strings.Contains(s, ":algorithm") {
			t.Errorf("expected no :algorithm in AST, got: %s", s)
		}
	})
	t.Run("sub_algorithm_value_captured", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t1 (id INT) PARTITION BY HASH(id) SUBPARTITION BY KEY ALGORITHM=1 (id)")
		s := ast.NodeToString(result.Items[0])
		if !strings.Contains(s, ":sub_algorithm 1") {
			t.Errorf("expected :sub_algorithm 1 in AST, got: %s", s)
		}
	})
}

func TestParseSubPartitionKeyAlgorithm(t *testing.T) {
	tests := []string{
		"CREATE TABLE t1 (id INT) PARTITION BY HASH(id) SUBPARTITION BY LINEAR KEY ALGORITHM=1 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY HASH(id) SUBPARTITION BY LINEAR KEY ALGORITHM=2 (id)",
		"CREATE TABLE t1 (id INT) PARTITION BY KEY(id) SUBPARTITION BY LINEAR KEY ALGORITHM=1 (id)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify LINEAR KEY algorithm value is captured in AST
	t.Run("linear_sub_algorithm_value_captured", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t1 (id INT) PARTITION BY HASH(id) SUBPARTITION BY LINEAR KEY ALGORITHM=2 (id)")
		s := ast.NodeToString(result.Items[0])
		if !strings.Contains(s, ":sub_algorithm 2") {
			t.Errorf("expected :sub_algorithm 2 in AST, got: %s", s)
		}
	})
}

func TestParseCreateTableAutoextendSize(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) AUTOEXTEND_SIZE = 4194304",
		"CREATE TABLE t (id INT) AUTOEXTEND_SIZE 4194304",
		"CREATE TABLE t (id INT) autoextend_size = 67108864",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) ENGINE_ATTRIBUTE = '{\"key\": \"value\"}'",
		"CREATE TABLE t (id INT) ENGINE_ATTRIBUTE '{\"key\": \"value\"}'",
		"CREATE TABLE t (id INT) engine_attribute = '{}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseColumnEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT ENGINE_ATTRIBUTE = '{\"key\": \"val\"}')",
		"CREATE TABLE t (id INT ENGINE_ATTRIBUTE '{\"key\": \"val\"}')",
		"CREATE TABLE t (id INT engine_attribute = '{}')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseColumnSecondaryEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE = '{\"key\": \"val\"}')",
		"CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE '{\"key\": \"val\"}')",
		"CREATE TABLE t (id INT secondary_engine_attribute = '{}')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateIndexEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE INDEX idx ON t (col) ENGINE_ATTRIBUTE = '{\"key\": \"val\"}'",
		"CREATE INDEX idx ON t (col) ENGINE_ATTRIBUTE '{\"key\": \"val\"}'",
		"CREATE INDEX idx ON t (col) engine_attribute = '{}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateIndexSecondaryEngineAttribute(t *testing.T) {
	tests := []string{
		"CREATE INDEX idx ON t (col) SECONDARY_ENGINE_ATTRIBUTE = '{\"key\": \"val\"}'",
		"CREATE INDEX idx ON t (col) SECONDARY_ENGINE_ATTRIBUTE '{\"key\": \"val\"}'",
		"CREATE INDEX idx ON t (col) secondary_engine_attribute = '{}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseAlterTableSecondaryLoad(t *testing.T) {
	tests := []string{
		"ALTER TABLE t SECONDARY_LOAD",
		"ALTER TABLE mydb.t secondary_load",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseRepairTableUseFrm(t *testing.T) {
	tests := []string{
		"REPAIR TABLE t USE_FRM",
		"REPAIR TABLE t QUICK EXTENDED USE_FRM",
		"REPAIR TABLE t1, t2 USE_FRM",
		"REPAIR LOCAL TABLE t USE_FRM",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify USE_FRM is captured in AST
	t.Run("use_frm_value_captured", func(t *testing.T) {
		result := ParseAndCheck(t, "REPAIR TABLE t USE_FRM")
		s := ast.NodeToString(result.Items[0])
		if !strings.Contains(s, ":use_frm true") {
			t.Errorf("expected :use_frm true in AST, got: %s", s)
		}
	})
}

func TestParseChecksumTableExtended(t *testing.T) {
	tests := []string{
		"CHECKSUM TABLE t EXTENDED",
		"CHECKSUM TABLE t1, t2 EXTENDED",
		"CHECKSUM TABLE t QUICK",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}

	// Verify EXTENDED is captured in AST
	t.Run("extended_value_captured", func(t *testing.T) {
		result := ParseAndCheck(t, "CHECKSUM TABLE t EXTENDED")
		s := ast.NodeToString(result.Items[0])
		if !strings.Contains(s, ":extended true") {
			t.Errorf("expected :extended true in AST, got: %s", s)
		}
	})
}

func TestParseCreateLogfileGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' ENGINE = NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' INITIAL_SIZE = 16M ENGINE NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :initial_size 16M :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' UNDO_BUFFER_SIZE = 8M ENGINE = NDBCLUSTER",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :undo_buffer_size 8M :engine NDBCLUSTER}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' REDO_BUFFER_SIZE 32M ENGINE = NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :redo_buffer_size 32M :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' NODEGROUP = 0 ENGINE = NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :nodegroup 0 :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' WAIT ENGINE = NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :wait true :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' COMMENT = 'test group' ENGINE = NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :comment "test group" :engine NDB}`,
		},
		{
			sql:  "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' INITIAL_SIZE 16M UNDO_BUFFER_SIZE 8M REDO_BUFFER_SIZE 32M NODEGROUP 0 WAIT COMMENT 'full opts' ENGINE NDB",
			want: `{CREATE_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo1.dat" :initial_size 16M :undo_buffer_size 8M :redo_buffer_size 32M :nodegroup 0 :wait true :comment "full opts" :engine NDB}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseAlterLogfileGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' ENGINE = NDB",
			want: `{ALTER_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo2.dat" :engine NDB}`,
		},
		{
			sql:  "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' INITIAL_SIZE = 32M ENGINE NDB",
			want: `{ALTER_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo2.dat" :initial_size 32M :engine NDB}`,
		},
		{
			sql:  "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' WAIT ENGINE = NDBCLUSTER",
			want: `{ALTER_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo2.dat" :wait true :engine NDBCLUSTER}`,
		},
		{
			sql:  "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' INITIAL_SIZE 16M WAIT ENGINE NDB",
			want: `{ALTER_LOGFILE_GROUP :loc 0 :name lg1 :undo_file "undo2.dat" :initial_size 16M :wait true :engine NDB}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseDropLogfileGroup(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "DROP LOGFILE GROUP lg1 ENGINE = NDB",
			want: "{DROP_LOGFILE_GROUP :loc 0 :name lg1 :engine NDB}",
		},
		{
			sql:  "DROP LOGFILE GROUP lg1 ENGINE NDBCLUSTER",
			want: "{DROP_LOGFILE_GROUP :loc 0 :name lg1 :engine NDBCLUSTER}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseShowCreateSchema(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SHOW CREATE SCHEMA mydb",
			want: "{SHOW :loc 0 :type CREATE DATABASE :from {TABLEREF :loc 19 :name mydb}}",
		},
		{
			sql:  "SHOW CREATE SCHEMA `my_db`",
			want: "{SHOW :loc 0 :type CREATE DATABASE :from {TABLEREF :loc 19 :name my_db}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseShowCreateDatabaseIfNotExists(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{
			sql:  "SHOW CREATE DATABASE IF NOT EXISTS mydb",
			want: "{SHOW :loc 0 :type CREATE DATABASE :from {TABLEREF :loc 35 :name mydb} :if_not_exists true}",
		},
		{
			sql:  "SHOW CREATE SCHEMA IF NOT EXISTS mydb",
			want: "{SHOW :loc 0 :type CREATE DATABASE :from {TABLEREF :loc 33 :name mydb} :if_not_exists true}",
		},
		{
			sql:  "SHOW CREATE DATABASE mydb",
			want: "{SHOW :loc 0 :type CREATE DATABASE :from {TABLEREF :loc 21 :name mydb}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			ParseAndCompare(t, tt.sql, tt.want)
		})
	}
}

func TestParseAlterTableSecondaryUnload(t *testing.T) {
	tests := []string{
		"ALTER TABLE t SECONDARY_UNLOAD",
		"ALTER TABLE mydb.t secondary_unload",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// Batch 102: BNF review — CREATE/DROP INDEX, CREATE/ALTER/DROP VIEW, CREATE/DROP TRIGGER

func TestBatch102_CreateTriggerIfNotExists(t *testing.T) {
	tests := []string{
		"CREATE TRIGGER IF NOT EXISTS trg BEFORE INSERT ON t FOR EACH ROW SET @a = 1",
		"CREATE TRIGGER IF NOT EXISTS trg2 AFTER DELETE ON t FOR EACH ROW SET @b = 2",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_CreateIndexFull(t *testing.T) {
	tests := []string{
		"CREATE UNIQUE INDEX idx1 ON t (col1)",
		"CREATE FULLTEXT INDEX ft_idx ON t (col1) WITH PARSER ngram",
		"CREATE SPATIAL INDEX sp_idx ON t (geo_col)",
		"CREATE INDEX idx1 USING BTREE ON t (col1, col2 DESC)",
		"CREATE INDEX idx1 ON t (col1) KEY_BLOCK_SIZE = 1024 COMMENT 'test' VISIBLE",
		"CREATE INDEX idx1 ON t (col1) ALGORITHM = INPLACE LOCK = NONE",
		"CREATE INDEX idx1 ON t ((col1 + col2))",
		"CREATE INDEX idx1 ON t (col1) ENGINE_ATTRIBUTE = '{}'",
		"CREATE INDEX idx1 ON t (col1) SECONDARY_ENGINE_ATTRIBUTE = '{}'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_DropIndexFull(t *testing.T) {
	tests := []string{
		"DROP INDEX idx1 ON t",
		"DROP INDEX idx1 ON t ALGORITHM = INPLACE",
		"DROP INDEX idx1 ON t LOCK = SHARED",
		"DROP INDEX idx1 ON t ALGORITHM = COPY LOCK = EXCLUSIVE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_CreateViewFull(t *testing.T) {
	tests := []string{
		"CREATE VIEW v AS SELECT 1",
		"CREATE OR REPLACE VIEW v AS SELECT 1",
		"CREATE ALGORITHM = MERGE VIEW v AS SELECT 1",
		"CREATE SQL SECURITY INVOKER VIEW v (a, b) AS SELECT 1, 2",
		"CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION",
		"CREATE VIEW v AS SELECT 1 WITH CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_AlterViewFull(t *testing.T) {
	tests := []string{
		"ALTER VIEW v AS SELECT 1",
		"ALTER ALGORITHM = TEMPTABLE VIEW v AS SELECT 1",
		"ALTER SQL SECURITY DEFINER VIEW v (a, b) AS SELECT 1, 2",
		"ALTER VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_DropViewFull(t *testing.T) {
	tests := []string{
		"DROP VIEW v1",
		"DROP VIEW IF EXISTS v1",
		"DROP VIEW v1, v2, v3",
		"DROP VIEW IF EXISTS v1 RESTRICT",
		"DROP VIEW IF EXISTS v1 CASCADE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch102_DropTriggerFull(t *testing.T) {
	tests := []string{
		"DROP TRIGGER trg1",
		"DROP TRIGGER IF EXISTS trg1",
		"DROP TRIGGER IF EXISTS mydb.trg1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// Batch 103: BNF review — routines and events
func TestBatch103_CreateProcedureIfNotExists(t *testing.T) {
	tests := []string{
		"CREATE PROCEDURE IF NOT EXISTS myproc() SELECT 1",
		"CREATE PROCEDURE IF NOT EXISTS mydb.myproc(IN x INT) SELECT x",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestBatch103_CreateEventEnableDisable(t *testing.T) {
	tests := []string{
		"CREATE EVENT ev1 ON SCHEDULE AT '2024-01-01 00:00:00' ENABLE DO SELECT 1",
		"CREATE EVENT ev1 ON SCHEDULE AT '2024-01-01 00:00:00' DISABLE DO SELECT 1",
		"CREATE EVENT ev1 ON SCHEDULE AT '2024-01-01 00:00:00' DISABLE ON SLAVE DO SELECT 1",
		"CREATE EVENT ev1 ON SCHEDULE EVERY 1 DAY ENABLE COMMENT 'daily' DO DELETE FROM t WHERE created < NOW() - INTERVAL 30 DAY",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// TestBNFReviewSecurity tests BNF review batch 105 - security statements.
func TestBNFReviewSecurity(t *testing.T) {
	tests := []string{
		// REVOKE IF EXISTS
		"REVOKE IF EXISTS SELECT ON db1.* FROM 'user1'@'localhost'",
		"REVOKE IF EXISTS ALL PRIVILEGES, GRANT OPTION FROM 'user1'@'localhost'",
		"REVOKE IF EXISTS INSERT, UPDATE ON *.* FROM user1",
		// REVOKE IGNORE UNKNOWN USER
		"REVOKE SELECT ON db1.* FROM 'user1'@'localhost' IGNORE UNKNOWN USER",
		"REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'user1'@'localhost' IGNORE UNKNOWN USER",
		"REVOKE IF EXISTS SELECT ON db1.t1 FROM user1 IGNORE UNKNOWN USER",
		// REVOKE IF EXISTS role
		"REVOKE IF EXISTS role1 FROM user1",
		// REVOKE PROXY IF EXISTS
		"REVOKE IF EXISTS PROXY ON 'admin'@'localhost' FROM 'user1'@'localhost'",
		// SET PASSWORD TO RANDOM
		"SET PASSWORD TO RANDOM",
		"SET PASSWORD FOR 'user1'@'localhost' TO RANDOM",
		// SET PASSWORD with REPLACE
		"SET PASSWORD = 'newpass' REPLACE 'oldpass'",
		"SET PASSWORD FOR user1 = 'newpass' REPLACE 'oldpass'",
		// SET PASSWORD with RETAIN CURRENT PASSWORD
		"SET PASSWORD = 'newpass' RETAIN CURRENT PASSWORD",
		"SET PASSWORD FOR 'user1'@'localhost' = 'newpass' REPLACE 'oldpass' RETAIN CURRENT PASSWORD",
		// CREATE USER DEFAULT ROLE
		"CREATE USER 'user1'@'localhost' IDENTIFIED BY 'pass' DEFAULT ROLE role1, role2",
		// CREATE USER with 2FA chaining
		"CREATE USER 'user1'@'localhost' IDENTIFIED BY 'pass1' AND IDENTIFIED BY 'pass2'",
		"CREATE USER 'user1'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'pass1' AND IDENTIFIED BY 'pass2' AND IDENTIFIED BY 'pass3'",
		// CREATE USER with INITIAL AUTHENTICATION
		"CREATE USER 'user1'@'localhost' IDENTIFIED WITH caching_sha2_password INITIAL AUTHENTICATION IDENTIFIED BY RANDOM PASSWORD",
		"CREATE USER 'user1'@'localhost' IDENTIFIED WITH caching_sha2_password INITIAL AUTHENTICATION IDENTIFIED WITH sha256_password AS 'hash'",
		"CREATE USER 'user1'@'localhost' IDENTIFIED WITH caching_sha2_password INITIAL AUTHENTICATION IDENTIFIED BY 'initpass'",
		// ALTER USER with REPLACE
		"ALTER USER 'user1'@'localhost' IDENTIFIED BY 'newpass' REPLACE 'oldpass'",
		"ALTER USER 'user1'@'localhost' IDENTIFIED BY 'newpass' REPLACE 'oldpass' RETAIN CURRENT PASSWORD",
		// ALTER USER with IDENTIFIED WITH ... BY ... REPLACE
		"ALTER USER 'user1'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'newpass' REPLACE 'oldpass'",
		// RENAME USER with quoted names
		"RENAME USER 'olduser'@'localhost' TO 'newuser'@'localhost'",
		"RENAME USER 'u1'@'h1' TO 'u2'@'h2', 'u3'@'h3' TO 'u4'@'h4'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// -----------------------------------------------------------------------
// Review batch 107: transaction_xa_locking BNF review
// -----------------------------------------------------------------------

func TestReviewBatch107_CommitWorkAndVariants(t *testing.T) {
	tests := []string{
		"COMMIT WORK",
		"COMMIT WORK AND CHAIN",
		"COMMIT AND NO CHAIN",
		"COMMIT WORK AND NO CHAIN",
		"COMMIT NO RELEASE",
		"COMMIT WORK NO RELEASE",
		"COMMIT WORK AND CHAIN RELEASE",
		"COMMIT WORK AND CHAIN NO RELEASE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestReviewBatch107_RollbackWorkAndVariants(t *testing.T) {
	tests := []string{
		"ROLLBACK WORK",
		"ROLLBACK WORK TO SAVEPOINT sp1",
		"ROLLBACK WORK TO sp1",
		"ROLLBACK AND NO CHAIN",
		"ROLLBACK WORK AND NO CHAIN",
		"ROLLBACK NO RELEASE",
		"ROLLBACK WORK NO RELEASE",
		"ROLLBACK WORK AND CHAIN RELEASE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestReviewBatch107_ReleaseSavepoint(t *testing.T) {
	ParseAndCompare(t, "RELEASE SAVEPOINT sp1", "{RELEASE_SAVEPOINT :loc 0 :name sp1}")
}

func TestReviewBatch107_LockTableSingular(t *testing.T) {
	tests := []string{
		"LOCK TABLE t READ",
		"LOCK TABLE t WRITE",
		"LOCK TABLE t READ LOCAL",
		"LOCK TABLE t LOW_PRIORITY WRITE",
		"LOCK TABLE t AS a READ, t2 WRITE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestReviewBatch107_UnlockTableSingular(t *testing.T) {
	ParseAndCheck(t, "UNLOCK TABLE")
}

// === Review Batch 108: SHOW BNF gaps ===

func TestReviewBatch108_ShowTablesIn(t *testing.T) {
	// SHOW TABLES {FROM | IN} db_name - IN variant
	stmt := parseShow(t, "SHOW TABLES IN mydb")
	if stmt.Type != "TABLES" {
		t.Errorf("expected TABLES, got %s", stmt.Type)
	}
	if stmt.From == nil || stmt.From.Name != "mydb" {
		t.Errorf("expected From=mydb, got %v", stmt.From)
	}
}

func TestReviewBatch108_ShowFullTablesIn(t *testing.T) {
	stmt := parseShow(t, "SHOW FULL TABLES IN mydb")
	if stmt.Type != "FULL TABLES" {
		t.Errorf("expected FULL TABLES, got %s", stmt.Type)
	}
	if stmt.From == nil || stmt.From.Name != "mydb" {
		t.Errorf("expected From=mydb, got %v", stmt.From)
	}
}

func TestReviewBatch108_ShowExtendedTables(t *testing.T) {
	ParseAndCheck(t, "SHOW EXTENDED TABLES")
	ParseAndCheck(t, "SHOW EXTENDED TABLES FROM mydb")
	ParseAndCheck(t, "SHOW EXTENDED TABLES LIKE 'test%'")
}

func TestReviewBatch108_ShowExtendedFullTables(t *testing.T) {
	ParseAndCheck(t, "SHOW EXTENDED FULL TABLES")
	ParseAndCheck(t, "SHOW EXTENDED FULL TABLES IN mydb")
	ParseAndCheck(t, "SHOW EXTENDED FULL TABLES LIKE 'user%'")
}

func TestReviewBatch108_ShowCountErrors(t *testing.T) {
	stmt := parseShow(t, "SHOW COUNT(*) ERRORS")
	if stmt.Type != "COUNT ERRORS" {
		t.Errorf("expected COUNT ERRORS, got %s", stmt.Type)
	}
}

func TestReviewBatch108_ShowCountWarnings(t *testing.T) {
	stmt := parseShow(t, "SHOW COUNT(*) WARNINGS")
	if stmt.Type != "COUNT WARNINGS" {
		t.Errorf("expected COUNT WARNINGS, got %s", stmt.Type)
	}
}

func TestReviewBatch108_ShowStorageEngines(t *testing.T) {
	stmt := parseShow(t, "SHOW STORAGE ENGINES")
	if stmt.Type != "ENGINES" {
		t.Errorf("expected ENGINES, got %s", stmt.Type)
	}
}

func TestReviewBatch108_ShowMasterLogs(t *testing.T) {
	stmt := parseShow(t, "SHOW MASTER LOGS")
	if stmt.Type != "BINARY LOGS" {
		t.Errorf("expected BINARY LOGS, got %s", stmt.Type)
	}
}

func TestReviewBatch108_ShowReplicaStatusForChannel(t *testing.T) {
	stmt := parseShow(t, "SHOW REPLICA STATUS FOR CHANNEL ch1")
	if stmt.Type != "REPLICA STATUS" {
		t.Errorf("expected REPLICA STATUS, got %s", stmt.Type)
	}
	if stmt.Channel != "ch1" {
		t.Errorf("expected channel ch1, got %s", stmt.Channel)
	}
}

func TestReviewBatch108_ShowSlaveStatusForChannel(t *testing.T) {
	stmt := parseShow(t, "SHOW SLAVE STATUS FOR CHANNEL ch1")
	if stmt.Type != "SLAVE STATUS" {
		t.Errorf("expected SLAVE STATUS, got %s", stmt.Type)
	}
	if stmt.Channel != "ch1" {
		t.Errorf("expected channel ch1, got %s", stmt.Channel)
	}
}

func TestReviewBatch108_ShowRelaylogEventsForChannel(t *testing.T) {
	stmt := parseShow(t, "SHOW RELAYLOG EVENTS FOR CHANNEL ch1")
	if stmt.Type != "RELAYLOG EVENTS" {
		t.Errorf("expected RELAYLOG EVENTS, got %s", stmt.Type)
	}
	if stmt.Channel != "ch1" {
		t.Errorf("expected channel ch1, got %s", stmt.Channel)
	}
}

func TestReviewBatch108_ShowCharset(t *testing.T) {
	stmt := parseShow(t, "SHOW CHARSET")
	if stmt.Type != "CHARACTER SET" {
		t.Errorf("expected CHARACTER SET, got %s", stmt.Type)
	}
}

func TestReviewBatch108_ShowCharsetLike(t *testing.T) {
	stmt := parseShow(t, "SHOW CHARSET LIKE 'utf8%'")
	if stmt.Type != "CHARACTER SET" {
		t.Errorf("expected CHARACTER SET, got %s", stmt.Type)
	}
	if stmt.Like == nil {
		t.Error("expected Like to be set")
	}
}

// ============================================================================
// Batch 109: review_set_explain
// ============================================================================

func TestReviewBatch109_SetNamesDefault(t *testing.T) {
	stmt := parseSet(t, "SET NAMES DEFAULT")
	if len(stmt.Assignments) != 1 {
		t.Fatalf("Assignments count = %d, want 1", len(stmt.Assignments))
	}
	if stmt.Assignments[0].Column.Column != "NAMES" {
		t.Errorf("Column = %s, want NAMES", stmt.Assignments[0].Column.Column)
	}
	val, ok := stmt.Assignments[0].Value.(*ast.StringLit)
	if !ok || val.Value != "DEFAULT" {
		t.Errorf("Value = %v, want DEFAULT", stmt.Assignments[0].Value)
	}
}

func TestReviewBatch109_SetCharsetDefault(t *testing.T) {
	stmt := parseSet(t, "SET CHARACTER SET DEFAULT")
	if len(stmt.Assignments) != 1 {
		t.Fatalf("Assignments count = %d, want 1", len(stmt.Assignments))
	}
	if stmt.Assignments[0].Column.Column != "CHARACTER SET" {
		t.Errorf("Column = %s, want CHARACTER SET", stmt.Assignments[0].Column.Column)
	}
	val, ok := stmt.Assignments[0].Value.(*ast.StringLit)
	if !ok || val.Value != "DEFAULT" {
		t.Errorf("Value = %v, want DEFAULT", stmt.Assignments[0].Value)
	}
}

func TestReviewBatch109_SetCharsetSynonym(t *testing.T) {
	stmt := parseSet(t, "SET CHARSET utf8mb4")
	if len(stmt.Assignments) != 1 {
		t.Fatalf("Assignments count = %d, want 1", len(stmt.Assignments))
	}
	if stmt.Assignments[0].Column.Column != "CHARACTER SET" {
		t.Errorf("Column = %s, want CHARACTER SET", stmt.Assignments[0].Column.Column)
	}
	val, ok := stmt.Assignments[0].Value.(*ast.StringLit)
	if !ok || val.Value != "utf8mb4" {
		t.Errorf("Value = %v, want utf8mb4", stmt.Assignments[0].Value)
	}
}

func TestReviewBatch109_SetCharsetSynonymDefault(t *testing.T) {
	stmt := parseSet(t, "SET CHARSET DEFAULT")
	if len(stmt.Assignments) != 1 {
		t.Fatalf("Assignments count = %d, want 1", len(stmt.Assignments))
	}
	val, ok := stmt.Assignments[0].Value.(*ast.StringLit)
	if !ok || val.Value != "DEFAULT" {
		t.Errorf("Value = %v, want DEFAULT", stmt.Assignments[0].Value)
	}
}

func TestReviewBatch109_SetAssignOperator(t *testing.T) {
	ParseAndCheck(t, "SET @x := 42")
}

func TestReviewBatch109_DescTable(t *testing.T) {
	ParseAndCheck(t, "DESC users")
}

func TestReviewBatch109_DescTableColumn(t *testing.T) {
	ParseAndCheck(t, "DESC users name")
}

func TestReviewBatch109_ExplainTableColumn(t *testing.T) {
	ParseAndCheck(t, "EXPLAIN users name")
}

func TestReviewBatch109_ExplainTableWild(t *testing.T) {
	ParseAndCheck(t, "EXPLAIN users '%name%'")
}

func TestReviewBatch109_ExplainTableStmt(t *testing.T) {
	ParseAndCheck(t, "EXPLAIN TABLE t")
}

func TestReviewBatch109_ExplainAnalyze(t *testing.T) {
	ParseAndCheck(t, "EXPLAIN ANALYZE SELECT * FROM t")
}

func TestReviewBatch109_ExplainFormatTree(t *testing.T) {
	ParseAndCheck(t, "EXPLAIN FORMAT = TREE SELECT * FROM t")
}

func TestReviewBatch109_HelpString(t *testing.T) {
	ParseAndCheck(t, "HELP 'SELECT'")
}

func TestReviewBatch109_SetResourceGroupFor(t *testing.T) {
	ParseAndCheck(t, "SET RESOURCE GROUP rg1 FOR 10, 20")
}

// ============================================================================
// Batch 110: review_load_import
// ============================================================================

func TestReviewBatch110_LoadDataLowPriority(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA LOW_PRIORITY INFILE '/tmp/data.csv' INTO TABLE t")
}

func TestReviewBatch110_LoadDataConcurrent(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA CONCURRENT INFILE '/tmp/data.csv' INTO TABLE t")
}

func TestReviewBatch110_LoadDataPartition(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t PARTITION (p0, p1)")
}

func TestReviewBatch110_LoadDataCharacterSet(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t CHARACTER SET utf8mb4")
}

func TestReviewBatch110_LoadDataOptionallyEnclosed(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA INFILE '/tmp/data.csv' INTO TABLE t FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"'")
}

func TestReviewBatch110_LoadXmlBasic(t *testing.T) {
	ParseAndCheck(t, "LOAD XML INFILE '/tmp/data.xml' INTO TABLE t")
}

func TestReviewBatch110_LoadXmlRowsIdentifiedBy(t *testing.T) {
	ParseAndCheck(t, "LOAD XML INFILE '/tmp/data.xml' INTO TABLE t ROWS IDENTIFIED BY '<row>'")
}

func TestReviewBatch110_LoadXmlCharsetSetClause(t *testing.T) {
	ParseAndCheck(t, "LOAD XML LOCAL INFILE '/tmp/data.xml' REPLACE INTO TABLE t CHARACTER SET utf8 SET col1 = UPPER(col1)")
}

func TestReviewBatch110_LoadDataFullSyntax(t *testing.T) {
	ParseAndCheck(t, "LOAD DATA LOW_PRIORITY LOCAL INFILE '/tmp/data.csv' REPLACE INTO TABLE db1.t1 PARTITION (p0) CHARACTER SET utf8mb4 FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' ESCAPED BY '\\\\' LINES STARTING BY '>>' TERMINATED BY '\\n' IGNORE 1 LINES (col1, col2) SET col3 = col1 + col2")
}

func TestReviewBatch110_LoadIndexIntoCachePartition(t *testing.T) {
	ParseAndCheck(t, "LOAD INDEX INTO CACHE t1 PARTITION (p0, p1) IGNORE LEAVES")
}

func TestReviewBatch110_LoadIndexIntoCacheIndex(t *testing.T) {
	ParseAndCheck(t, "LOAD INDEX INTO CACHE t1 INDEX (idx1, idx2)")
}

func TestReviewBatch110_LoadIndexIntoCachePartitionAll(t *testing.T) {
	ParseAndCheck(t, "LOAD INDEX INTO CACHE t1 PARTITION (ALL)")
}

func TestReviewBatch110_LoadIndexIntoCacheMultiTable(t *testing.T) {
	ParseAndCheck(t, "LOAD INDEX INTO CACHE t1 IGNORE LEAVES, t2 INDEX (idx1)")
}

func TestReviewBatch110_TableIntoOutfile(t *testing.T) {
	ParseAndCheck(t, "TABLE t1 INTO OUTFILE '/tmp/out.csv'")
}

func TestReviewBatch110_TableIntoDumpfile(t *testing.T) {
	ParseAndCheck(t, "TABLE t1 INTO DUMPFILE '/tmp/out.dat'")
}

func TestReviewBatch110_TableIntoVar(t *testing.T) {
	ParseAndCheck(t, "TABLE t1 INTO @v1, @v2")
}

func TestReviewBatch110_ValuesBasic(t *testing.T) {
	ParseAndCheck(t, "VALUES ROW(1, 2), ROW(3, 4)")
}

func TestReviewBatch110_ValuesOrderByLimit(t *testing.T) {
	ParseAndCheck(t, "VALUES ROW(1, 2), ROW(3, 4) ORDER BY column_0 LIMIT 1")
}

// --- Batch 111: review_maintenance_utility ---

func TestReviewBatch111_AnalyzeTableNoWriteToBinlog(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE NO_WRITE_TO_BINLOG TABLE t1")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_AnalyzeTableLocal(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE LOCAL TABLE t1")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_AnalyzeTableUsingData(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1 USING DATA '{\"key\": \"value\"}'")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":using_data") {
		t.Errorf("expected :using_data, got %s", s)
	}
}

func TestReviewBatch111_AnalyzeTableHistogramBuckets(t *testing.T) {
	result := ParseAndCheck(t, "ANALYZE TABLE t1 UPDATE HISTOGRAM ON col1, col2 WITH 100 BUCKETS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":buckets 100") {
		t.Errorf("expected :buckets 100, got %s", s)
	}
}

func TestReviewBatch111_OptimizeTableNoWriteToBinlog(t *testing.T) {
	result := ParseAndCheck(t, "OPTIMIZE NO_WRITE_TO_BINLOG TABLE t1")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_OptimizeTableLocal(t *testing.T) {
	result := ParseAndCheck(t, "OPTIMIZE LOCAL TABLE t1, t2")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_RepairTableNoWriteToBinlog(t *testing.T) {
	result := ParseAndCheck(t, "REPAIR NO_WRITE_TO_BINLOG TABLE t1")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_RepairTableLocal(t *testing.T) {
	result := ParseAndCheck(t, "REPAIR LOCAL TABLE t1 QUICK EXTENDED USE_FRM")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_CheckTableFastUpperCase(t *testing.T) {
	result := ParseAndCheck(t, "CHECK TABLE t1 FAST")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "FAST") {
		t.Errorf("expected FAST (uppercase), got %s", s)
	}
}

func TestReviewBatch111_CheckTableMultipleOptions(t *testing.T) {
	ParseAndCheck(t, "CHECK TABLE t1 QUICK FAST MEDIUM EXTENDED CHANGED")
}

func TestReviewBatch111_FlushEngineLogs(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH ENGINE LOGS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "ENGINE LOGS") {
		t.Errorf("expected ENGINE LOGS, got %s", s)
	}
}

func TestReviewBatch111_FlushErrorLogs(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH ERROR LOGS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "ERROR LOGS") {
		t.Errorf("expected ERROR LOGS, got %s", s)
	}
}

func TestReviewBatch111_FlushGeneralLogs(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH GENERAL LOGS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "GENERAL LOGS") {
		t.Errorf("expected GENERAL LOGS, got %s", s)
	}
}

func TestReviewBatch111_FlushSlowLogs(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH SLOW LOGS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "SLOW LOGS") {
		t.Errorf("expected SLOW LOGS, got %s", s)
	}
}

func TestReviewBatch111_FlushRelayLogs(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH RELAY LOGS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "RELAY LOGS") {
		t.Errorf("expected RELAY LOGS, got %s", s)
	}
}

func TestReviewBatch111_FlushRelayLogsForChannel(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH RELAY LOGS FOR CHANNEL 'channel_1'")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "RELAY LOGS") {
		t.Errorf("expected RELAY LOGS, got %s", s)
	}
	if !strings.Contains(s, ":relay_channel") {
		t.Errorf("expected :relay_channel, got %s", s)
	}
}

func TestReviewBatch111_FlushMultipleOptions(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH HOSTS, STATUS")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, "HOSTS") || !strings.Contains(s, "STATUS") {
		t.Errorf("expected HOSTS and STATUS, got %s", s)
	}
}

func TestReviewBatch111_FlushNoWriteToBinlog(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH NO_WRITE_TO_BINLOG TABLES")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_FlushLocalTables(t *testing.T) {
	result := ParseAndCheck(t, "FLUSH LOCAL TABLES")
	s := ast.NodeToString(result.Items[0])
	if !strings.Contains(s, ":no_write_to_binlog true") {
		t.Errorf("expected :no_write_to_binlog true, got %s", s)
	}
}

func TestReviewBatch111_FlushLogs(t *testing.T) {
	ParseAndCheck(t, "FLUSH LOGS")
}

func TestReviewBatch111_FlushOptimizerCosts(t *testing.T) {
	ParseAndCheck(t, "FLUSH OPTIMIZER_COSTS")
}

func TestReviewBatch111_FlushUserResources(t *testing.T) {
	ParseAndCheck(t, "FLUSH USER_RESOURCES")
}

func TestReviewBatch111_BinlogStr(t *testing.T) {
	ParseAndCheck(t, "BINLOG 'base64str'")
}

func TestReviewBatch111_CacheIndexBasic(t *testing.T) {
	ParseAndCheck(t, "CACHE INDEX t1 IN hot_cache")
}

func TestReviewBatch111_KillConnection(t *testing.T) {
	ParseAndCheck(t, "KILL CONNECTION 42")
}

func TestReviewBatch111_KillQuery(t *testing.T) {
	ParseAndCheck(t, "KILL QUERY 42")
}

func TestReviewBatch111_DoMultiExpr(t *testing.T) {
	ParseAndCheck(t, "DO 1 + 2, SLEEP(1)")
}

func TestReviewBatch111_ChecksumTableQuick(t *testing.T) {
	ParseAndCheck(t, "CHECKSUM TABLE t1, t2 QUICK")
}

func TestReviewBatch111_ChecksumTableExtended(t *testing.T) {
	ParseAndCheck(t, "CHECKSUM TABLE t1 EXTENDED")
}

func TestReviewBatch111_Shutdown(t *testing.T) {
	ParseAndCheck(t, "SHUTDOWN")
}

// ============================================================================
// Batch 116: Expression deep review
// ============================================================================

func TestReviewBatch116_BinaryPrefix(t *testing.T) {
	tests := []string{
		"SELECT BINARY 'hello'",
		"SELECT BINARY col FROM t",
		"SELECT BINARY col = 'abc' FROM t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---- Batch 117: DEFINER propagation ----

func TestDefinerPropagation(t *testing.T) {
	tests := []string{
		// CREATE VIEW with DEFINER at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' VIEW v1 AS SELECT 1",
		"CREATE DEFINER = admin@localhost VIEW v1 AS SELECT 1",
		"CREATE DEFINER = CURRENT_USER VIEW v1 AS SELECT 1",
		"CREATE DEFINER = CURRENT_USER() VIEW v1 AS SELECT 1",
		// CREATE VIEW with DEFINER + SQL SECURITY at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' SQL SECURITY INVOKER VIEW v1 AS SELECT 1",
		"CREATE DEFINER = root@localhost SQL SECURITY DEFINER VIEW v1 AS SELECT 1",
		// CREATE FUNCTION with DEFINER at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1",
		"CREATE DEFINER = CURRENT_USER FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1",
		// CREATE PROCEDURE with DEFINER at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' PROCEDURE p1() SELECT 1",
		"CREATE DEFINER = CURRENT_USER() PROCEDURE p1() SELECT 1",
		// CREATE TRIGGER with DEFINER at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' TRIGGER tr1 BEFORE INSERT ON t1 FOR EACH ROW SET @x = 1",
		"CREATE DEFINER = CURRENT_USER TRIGGER tr1 AFTER UPDATE ON t1 FOR EACH ROW SET @x = 1",
		// CREATE EVENT with DEFINER at dispatch layer
		"CREATE DEFINER = 'admin'@'localhost' EVENT ev1 ON SCHEDULE EVERY 1 HOUR DO SELECT 1",
		"CREATE DEFINER = CURRENT_USER EVENT ev1 ON SCHEDULE AT '2026-01-01 00:00:00' DO SELECT 1",
		// ALTER VIEW with DEFINER at dispatch layer
		"ALTER DEFINER = 'admin'@'localhost' VIEW v1 AS SELECT 2",
		"ALTER DEFINER = CURRENT_USER() SQL SECURITY INVOKER VIEW v1 AS SELECT 2",
		// ALTER EVENT with DEFINER at dispatch layer
		"ALTER DEFINER = 'admin'@'localhost' EVENT ev1 ENABLE",
		"ALTER DEFINER = CURRENT_USER EVENT ev1 ENABLE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestDefinerPropagation_Values(t *testing.T) {
	// Verify DEFINER value is actually captured in AST
	result, err := Parse("CREATE DEFINER = 'admin'@'localhost' VIEW v1 AS SELECT 1")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*ast.CreateViewStmt)
	if !ok {
		t.Fatalf("expected CreateViewStmt, got %T", result.Items[0])
	}
	if stmt.Definer != "'admin'@'localhost'" {
		t.Errorf("expected Definer='admin'@'localhost', got %q", stmt.Definer)
	}

	// Verify SQL SECURITY is captured
	result2, err := Parse("CREATE DEFINER = root@localhost SQL SECURITY INVOKER VIEW v2 AS SELECT 2")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt2, ok := result2.Items[0].(*ast.CreateViewStmt)
	if !ok {
		t.Fatalf("expected CreateViewStmt, got %T", result2.Items[0])
	}
	if stmt2.Definer != "root@localhost" {
		t.Errorf("expected Definer=root@localhost, got %q", stmt2.Definer)
	}
	if stmt2.SqlSecurity != "INVOKER" && stmt2.SqlSecurity != "invoker" {
		t.Errorf("expected SqlSecurity=INVOKER, got %q", stmt2.SqlSecurity)
	}

	// Verify DEFINER on CREATE TRIGGER
	result3, err := Parse("CREATE DEFINER = 'admin'@'%' TRIGGER tr1 BEFORE INSERT ON t1 FOR EACH ROW SET @x = 1")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt3, ok := result3.Items[0].(*ast.CreateTriggerStmt)
	if !ok {
		t.Fatalf("expected CreateTriggerStmt, got %T", result3.Items[0])
	}
	if stmt3.Definer != "'admin'@'%'" {
		t.Errorf("expected Definer='admin'@'%%', got %q", stmt3.Definer)
	}

	// Verify DEFINER on ALTER EVENT
	result4, err := Parse("ALTER DEFINER = 'root'@'localhost' EVENT ev1 ENABLE")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt4, ok := result4.Items[0].(*ast.AlterEventStmt)
	if !ok {
		t.Fatalf("expected AlterEventStmt, got %T", result4.Items[0])
	}
	if stmt4.Definer != "'root'@'localhost'" {
		t.Errorf("expected Definer='root'@'localhost', got %q", stmt4.Definer)
	}

	// Verify CURRENT_USER()
	result5, err := Parse("CREATE DEFINER = CURRENT_USER() FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt5, ok := result5.Items[0].(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("expected CreateFunctionStmt, got %T", result5.Items[0])
	}
	if stmt5.Definer != "CURRENT_USER()" {
		t.Errorf("expected Definer=CURRENT_USER(), got %q", stmt5.Definer)
	}
}

// ---- Batch 118: ALTER USER registration and factor auth ----

func TestAlterUserRegistration(t *testing.T) {
	tests := []string{
		// registration_option with user
		"ALTER USER 'admin'@'localhost' 2 FACTOR INITIATE REGISTRATION",
		"ALTER USER 'admin'@'localhost' 3 FACTOR INITIATE REGISTRATION",
		"ALTER USER 'admin'@'localhost' 2 FACTOR FINISH REGISTRATION SET CHALLENGE_RESPONSE AS 'response_data'",
		"ALTER USER 'admin'@'localhost' 3 FACTOR UNREGISTER",
		// registration_option with USER()
		"ALTER USER USER() 2 FACTOR INITIATE REGISTRATION",
		"ALTER USER USER() 3 FACTOR FINISH REGISTRATION SET CHALLENGE_RESPONSE AS 'response'",
		"ALTER USER USER() 2 FACTOR UNREGISTER",
		// IF EXISTS
		"ALTER USER IF EXISTS 'admin'@'localhost' 2 FACTOR INITIATE REGISTRATION",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestAlterUserFactorAuth(t *testing.T) {
	tests := []string{
		// ADD factor
		"ALTER USER 'admin'@'localhost' ADD 2 FACTOR IDENTIFIED BY 'password'",
		"ALTER USER 'admin'@'localhost' ADD 2 FACTOR IDENTIFIED BY RANDOM PASSWORD",
		"ALTER USER 'admin'@'localhost' ADD 3 FACTOR IDENTIFIED WITH caching_sha2_password BY 'pass'",
		"ALTER USER 'admin'@'localhost' ADD 2 FACTOR IDENTIFIED WITH auth_plugin AS 'hash'",
		// MODIFY factor
		"ALTER USER 'admin'@'localhost' MODIFY 2 FACTOR IDENTIFIED BY 'newpass'",
		"ALTER USER 'admin'@'localhost' MODIFY 3 FACTOR IDENTIFIED WITH caching_sha2_password BY RANDOM PASSWORD",
		// DROP factor
		"ALTER USER 'admin'@'localhost' DROP 2 FACTOR",
		"ALTER USER 'admin'@'localhost' DROP 3 FACTOR",
		// Multiple factor ops
		"ALTER USER 'admin'@'localhost' ADD 2 FACTOR IDENTIFIED BY 'pass2' ADD 3 FACTOR IDENTIFIED BY 'pass3'",
		"ALTER USER 'admin'@'localhost' DROP 2 FACTOR DROP 3 FACTOR",
		// USER() with user_func_auth_option
		"ALTER USER USER() IDENTIFIED BY 'newpass'",
		"ALTER USER USER() IDENTIFIED BY 'newpass' REPLACE 'oldpass'",
		"ALTER USER USER() IDENTIFIED BY 'newpass' REPLACE 'oldpass' RETAIN CURRENT PASSWORD",
		"ALTER USER USER() DISCARD OLD PASSWORD",
		// IF EXISTS
		"ALTER USER IF EXISTS 'admin'@'localhost' ADD 2 FACTOR IDENTIFIED BY 'pass'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestAlterUserFactorAuth_Values(t *testing.T) {
	// Verify USER() is captured
	result, err := Parse("ALTER USER USER() IDENTIFIED BY 'newpass'")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*ast.AlterUserStmt)
	if !ok {
		t.Fatalf("expected AlterUserStmt, got %T", result.Items[0])
	}
	if !stmt.IsUserFunc {
		t.Error("expected IsUserFunc=true")
	}
	if len(stmt.Users) != 1 || stmt.Users[0].Name != "USER()" {
		t.Errorf("expected user name USER(), got %v", stmt.Users)
	}

	// Verify registration op
	result2, err := Parse("ALTER USER 'admin'@'localhost' 2 FACTOR INITIATE REGISTRATION")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt2, ok := result2.Items[0].(*ast.AlterUserStmt)
	if !ok {
		t.Fatalf("expected AlterUserStmt, got %T", result2.Items[0])
	}
	if stmt2.RegistrationOp == nil {
		t.Fatal("expected RegistrationOp to be set")
	}
	if stmt2.RegistrationOp.Factor != 2 {
		t.Errorf("expected Factor=2, got %d", stmt2.RegistrationOp.Factor)
	}
	if stmt2.RegistrationOp.Action != "INITIATE" {
		t.Errorf("expected Action=INITIATE, got %s", stmt2.RegistrationOp.Action)
	}

	// Verify factor op
	result3, err := Parse("ALTER USER 'admin'@'localhost' ADD 3 FACTOR IDENTIFIED WITH caching_sha2_password BY 'pass'")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt3, ok := result3.Items[0].(*ast.AlterUserStmt)
	if !ok {
		t.Fatalf("expected AlterUserStmt, got %T", result3.Items[0])
	}
	if len(stmt3.FactorOps) != 1 {
		t.Fatalf("expected 1 FactorOp, got %d", len(stmt3.FactorOps))
	}
	fop := stmt3.FactorOps[0]
	if fop.Action != "ADD" || fop.Factor != 3 || fop.AuthPlugin != "caching_sha2_password" || fop.Password != "pass" {
		t.Errorf("unexpected FactorOp: %+v", fop)
	}
}

// ============================================================================
// Batch 120: depth_fix_outfuncs_expr_loc
// ============================================================================

func TestDeparseBinaryExprLoc(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 + 2")
	got := ast.NodeToString(result.Items[0])
	if !strings.Contains(got, "BINEXPR :op + :loc") {
		t.Errorf("BinaryExpr missing :loc in deparse output: %s", got)
	}
}

func TestDeparseUnaryExprLoc(t *testing.T) {
	result := ParseAndCheck(t, "SELECT -1")
	got := ast.NodeToString(result.Items[0])
	if !strings.Contains(got, "UNARY :op - :loc") {
		t.Errorf("UnaryExpr missing :loc in deparse output: %s", got)
	}
}

func TestDeparseIsExprLoc(t *testing.T) {
	result := ParseAndCheck(t, "SELECT 1 IS NULL")
	got := ast.NodeToString(result.Items[0])
	if !strings.Contains(got, "IS :loc") {
		t.Errorf("IsExpr missing :loc in deparse output: %s", got)
	}
}

func TestParseCreateTableChecksum(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) CHECKSUM = 1",
		"CREATE TABLE t (id INT) CHECKSUM = 0",
		"CREATE TABLE t (id INT) CHECKSUM 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableTablespace(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) TABLESPACE ts1",
		"CREATE TABLE t (id INT) TABLESPACE = ts1",
		"CREATE TABLE t (id INT) TABLESPACE ts1 STORAGE DISK",
		"CREATE TABLE t (id INT) TABLESPACE ts1 STORAGE MEMORY",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableEncryption(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) ENCRYPTION = 'Y'",
		"CREATE TABLE t (id INT) ENCRYPTION = 'N'",
		"CREATE TABLE t (id INT) ENCRYPTION 'Y'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableUnion(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) UNION = (t1, t2, t3)",
		"CREATE TABLE t (id INT) UNION = (t1)",
		"CREATE TABLE t (id INT) UNION (t1, t2)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableDataDirectory(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) DATA DIRECTORY = '/var/data'",
		"CREATE TABLE t (id INT) DATA DIRECTORY '/var/data'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

func TestParseCreateTableIndexDirectory(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id INT) INDEX DIRECTORY = '/var/idx'",
		"CREATE TABLE t (id INT) INDEX DIRECTORY '/var/idx'",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Section 1.4: Reserved Word Registration
// ---------------------------------------------------------------------------

// TestParseReservedWordDUAL tests DUAL as keyword in FROM clause.
func TestParseReservedWordDUAL(t *testing.T) {
	ParseAndCheck(t, "SELECT 1 FROM DUAL")
	ParseAndCheck(t, "SELECT 1 + 2 FROM DUAL")
	ParseAndCheck(t, "SELECT NOW() FROM DUAL")
}

// TestParseReservedWordMAXVALUE tests MAXVALUE recognized in partition definitions.
func TestParseReservedWordMAXVALUE(t *testing.T) {
	// MAXVALUE in parentheses
	ParseAndCheck(t, "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN (MAXVALUE))")
	// MAXVALUE without parentheses (bare form)
	ParseAndCheck(t, "CREATE TABLE t (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN MAXVALUE)")
}

// TestParseReservedWordGROUPING tests GROUPING function with ROLLUP.
func TestParseReservedWordGROUPING(t *testing.T) {
	ParseAndCheck(t, "SELECT GROUPING(a) FROM t GROUP BY a WITH ROLLUP")
	ParseAndCheck(t, "SELECT a, b, GROUPING(a, b) FROM t GROUP BY a, b WITH ROLLUP")
}

// TestParseReservedWordVALUE tests VALUE as a column name in SELECT.
func TestParseReservedWordVALUE(t *testing.T) {
	ParseAndCheck(t, "SELECT VALUE FROM t")
	ParseAndCheck(t, "SELECT t.VALUE FROM t")
}

// --- Section 1.3: SIGNED modifier ---

// TestParseSIGNEDModifier tests that SIGNED is accepted on numeric types in CREATE TABLE.
func TestParseSIGNEDModifier(t *testing.T) {
	// SIGNED on various integer types in CREATE TABLE context
	tests := []string{
		"CREATE TABLE t (a INT SIGNED)",
		"CREATE TABLE t (a BIGINT SIGNED)",
		"CREATE TABLE t (a TINYINT SIGNED)",
		"CREATE TABLE t (a SMALLINT SIGNED)",
		"CREATE TABLE t (a MEDIUMINT SIGNED)",
		"CREATE TABLE t (a TINYINT SIGNED NOT NULL)",
		"CREATE TABLE t (a DECIMAL SIGNED)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// Section 1.5: Numeric Type Shorthands in Combinations
// ============================================================================

func TestNumericShorthandCombinations(t *testing.T) {
	// Scenario 1: INT1 UNSIGNED ZEROFILL — full modifier chain on shorthand
	t.Run("INT1 UNSIGNED ZEROFILL", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("INT1 UNSIGNED ZEROFILL")}
		p.advance()
		dt, err := p.parseDataType()
		if err != nil {
			t.Fatalf("parseDataType error: %v", err)
		}
		got := ast.NodeToString(dt)
		want := "{DATATYPE :loc 0 :name TINYINT :unsigned true :zerofill true}"
		if got != want {
			t.Errorf("got:  %s\nwant: %s", got, want)
		}
	})

	// Scenario 2: INT4(11) UNSIGNED — display width on shorthand
	t.Run("INT4(11) UNSIGNED", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("INT4(11) UNSIGNED")}
		p.advance()
		dt, err := p.parseDataType()
		if err != nil {
			t.Fatalf("parseDataType error: %v", err)
		}
		got := ast.NodeToString(dt)
		want := "{DATATYPE :loc 0 :name INT :len 11 :unsigned true}"
		if got != want {
			t.Errorf("got:  %s\nwant: %s", got, want)
		}
	})

	// Scenario 3: ALTER TABLE t ADD COLUMN b REAL DEFAULT 0.0 — type synonym in ALTER TABLE
	t.Run("ALTER TABLE ADD COLUMN REAL", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t ADD COLUMN b REAL DEFAULT 0.0")
		if len(stmt.Commands) != 1 {
			t.Fatalf("Commands count = %d, want 1", len(stmt.Commands))
		}
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddColumn {
			t.Errorf("Type = %d, want ATAddColumn", cmd.Type)
		}
		if cmd.Column == nil {
			t.Fatal("Column is nil")
		}
		if cmd.Column.Name != "b" {
			t.Errorf("Column.Name = %s, want b", cmd.Column.Name)
		}
		if cmd.Column.TypeName == nil || cmd.Column.TypeName.Name != "DOUBLE" {
			t.Errorf("TypeName.Name = %v, want DOUBLE", cmd.Column.TypeName)
		}
	})

	// Scenario 4: ALTER TABLE t MODIFY c DEC(10,2) NOT NULL — type synonym in MODIFY COLUMN
	t.Run("ALTER TABLE MODIFY DEC", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t MODIFY c DEC(10,2) NOT NULL")
		if len(stmt.Commands) != 1 {
			t.Fatalf("Commands count = %d, want 1", len(stmt.Commands))
		}
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATModifyColumn {
			t.Errorf("Type = %d, want ATModifyColumn", cmd.Type)
		}
		if cmd.Name != "c" {
			t.Errorf("Name = %s, want c", cmd.Name)
		}
		if cmd.Column == nil {
			t.Fatal("Column is nil")
		}
		if cmd.Column.TypeName == nil || cmd.Column.TypeName.Name != "DECIMAL" {
			t.Errorf("TypeName.Name = %v, want DECIMAL", cmd.Column.TypeName)
		}
		if cmd.Column.TypeName.Length != 10 {
			t.Errorf("TypeName.Length = %d, want 10", cmd.Column.TypeName.Length)
		}
		if cmd.Column.TypeName.Scale != 2 {
			t.Errorf("TypeName.Scale = %d, want 2", cmd.Column.TypeName.Scale)
		}
	})

	// Scenario 5: all shorthands in one table
	t.Run("all shorthands in one table", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t (a FLOAT4, b FLOAT8, c INT1, d INT2, e INT3, f INT4, g INT8)")
		if result.Len() == 0 {
			t.Fatal("empty result")
		}
		ct, ok := result.Items[0].(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("expected *ast.CreateTableStmt, got %T", result.Items[0])
		}
		expected := []struct {
			name     string
			typeName string
		}{
			{"a", "FLOAT"},
			{"b", "DOUBLE"},
			{"c", "TINYINT"},
			{"d", "SMALLINT"},
			{"e", "MEDIUMINT"},
			{"f", "INT"},
			{"g", "BIGINT"},
		}
		if len(ct.Columns) != len(expected) {
			t.Fatalf("Columns count = %d, want %d", len(ct.Columns), len(expected))
		}
		for i, exp := range expected {
			col := ct.Columns[i]
			if col.Name != exp.name {
				t.Errorf("Columns[%d].Name = %s, want %s", i, col.Name, exp.name)
			}
			if col.TypeName == nil || col.TypeName.Name != exp.typeName {
				t.Errorf("Columns[%d].TypeName.Name = %v, want %s", i, col.TypeName, exp.typeName)
			}
		}
	})

	// Scenario 6: SERIAL type (existing, regression check)
	t.Run("SERIAL regression", func(t *testing.T) {
		result := ParseAndCheck(t, "CREATE TABLE t (a SERIAL)")
		if result.Len() == 0 {
			t.Fatal("empty result")
		}
		ct, ok := result.Items[0].(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("expected *ast.CreateTableStmt, got %T", result.Items[0])
		}
		if len(ct.Columns) != 1 {
			t.Fatalf("Columns count = %d, want 1", len(ct.Columns))
		}
		if ct.Columns[0].TypeName == nil || ct.Columns[0].TypeName.Name != "SERIAL" {
			t.Errorf("TypeName.Name = %v, want SERIAL", ct.Columns[0].TypeName)
		}
	})
}

// TestParseCastSIGNEDRegression tests that CAST(x AS SIGNED) still works.
func TestParseCastSIGNEDRegression(t *testing.T) {
	tests := []struct {
		input    string
		typeName string
	}{
		{"CAST(x AS SIGNED)", "SIGNED"},
		{"CAST(x AS SIGNED INTEGER)", "SIGNED"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr := parseExpr(t, tt.input)
			ce, ok := expr.(*ast.CastExpr)
			if !ok {
				t.Fatalf("expected *ast.CastExpr, got %T", expr)
			}
			if ce.TypeName.Name != tt.typeName {
				t.Errorf("TypeName = %s, want %s", ce.TypeName.Name, tt.typeName)
			}
		})
	}
}

// ============================================================================
// Section 2.1: ALTER TABLE PARTITION BY
// ============================================================================

func TestAlterTablePartitionBy(t *testing.T) {
	t.Run("basic hash repartition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4")
		if len(stmt.Commands) != 1 {
			t.Fatalf("Commands count = %d, want 1", len(stmt.Commands))
		}
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATPartitionBy {
			t.Errorf("Type = %d, want ATPartitionBy", cmd.Type)
		}
		if cmd.PartitionBy == nil {
			t.Fatal("PartitionBy is nil")
		}
		if cmd.PartitionBy.Type != ast.PartitionHash {
			t.Errorf("PartitionBy.Type = %d, want PartitionHash", cmd.PartitionBy.Type)
		}
		if cmd.PartitionBy.NumParts != 4 {
			t.Errorf("NumParts = %d, want 4", cmd.PartitionBy.NumParts)
		}
	})

	t.Run("key repartition", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t PARTITION BY KEY(id) PARTITIONS 4")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATPartitionBy {
			t.Errorf("Type = %d, want ATPartitionBy", cmd.Type)
		}
		if cmd.PartitionBy.Type != ast.PartitionKey {
			t.Errorf("PartitionBy.Type = %d, want PartitionKey", cmd.PartitionBy.Type)
		}
	})

	t.Run("range repartition", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t PARTITION BY RANGE(id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN (MAXVALUE))")
	})

	t.Run("list repartition", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t PARTITION BY LIST(status) (PARTITION p0 VALUES IN (1,2), PARTITION p1 VALUES IN (3,4))")
	})

	t.Run("range columns", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t PARTITION BY RANGE COLUMNS(created_at) (PARTITION p0 VALUES LESS THAN ('2024-01-01'))")
	})

	t.Run("linear hash", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t PARTITION BY LINEAR HASH(id) PARTITIONS 8")
		cmd := stmt.Commands[0]
		if cmd.PartitionBy == nil {
			t.Fatal("PartitionBy is nil")
		}
		if !cmd.PartitionBy.Linear {
			t.Errorf("Linear = false, want true")
		}
		if cmd.PartitionBy.NumParts != 8 {
			t.Errorf("NumParts = %d, want 8", cmd.PartitionBy.NumParts)
		}
	})

	t.Run("with subpartition", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4 SUBPARTITION BY KEY(name) SUBPARTITIONS 2")
	})

	t.Run("KEY with algorithm", func(t *testing.T) {
		ParseAndCheck(t, "ALTER TABLE t PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 4")
	})

	t.Run("combined with ALTER TABLE options", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, ALGORITHM=INPLACE")
		if len(stmt.Commands) != 2 {
			t.Fatalf("Commands count = %d, want 2", len(stmt.Commands))
		}
		if stmt.Commands[0].Type != ast.ATPartitionBy {
			t.Errorf("first cmd Type = %d, want ATPartitionBy", stmt.Commands[0].Type)
		}
		if stmt.Commands[1].Type != ast.ATAlgorithm {
			t.Errorf("second cmd Type = %d, want ATAlgorithm", stmt.Commands[1].Type)
		}
	})

	t.Run("rejected: missing partition specification", func(t *testing.T) {
		ParseExpectError(t, "ALTER TABLE t PARTITION BY")
	})

	t.Run("REMOVE PARTITIONING regression", func(t *testing.T) {
		stmt := parseAlterTable(t, "ALTER TABLE t REMOVE PARTITIONING")
		if len(stmt.Commands) != 1 {
			t.Fatalf("Commands count = %d, want 1", len(stmt.Commands))
		}
		if stmt.Commands[0].Type != ast.ATRemovePartitioning {
			t.Errorf("Type = %d, want ATRemovePartitioning", stmt.Commands[0].Type)
		}
	})
}

// ---------------------------------------------------------------------------
// Section 2.2: Window Function Names as Reserved Words
// ---------------------------------------------------------------------------

func TestWindowFunctionNames(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"RANK", "SELECT RANK() OVER (ORDER BY score DESC) FROM t"},
		{"DENSE_RANK", "SELECT DENSE_RANK() OVER (ORDER BY score DESC) FROM t"},
		{"ROW_NUMBER", "SELECT ROW_NUMBER() OVER (ORDER BY id) FROM t"},
		{"NTILE", "SELECT NTILE(4) OVER (ORDER BY id) FROM t"},
		{"LAG", "SELECT LAG(val, 1) OVER (ORDER BY id) FROM t"},
		{"LEAD", "SELECT LEAD(val, 1, 0) OVER (ORDER BY id) FROM t"},
		{"FIRST_VALUE", "SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t"},
		{"LAST_VALUE", "SELECT LAST_VALUE(val) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) FROM t"},
		{"NTH_VALUE", "SELECT NTH_VALUE(val, 2) OVER (ORDER BY id) FROM t"},
		{"PERCENT_RANK", "SELECT PERCENT_RANK() OVER (ORDER BY score) FROM t"},
		{"CUME_DIST", "SELECT CUME_DIST() OVER (ORDER BY score) FROM t"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseAndCheck(t, tt.sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Section 2.3: Interval Unit Keywords
// ---------------------------------------------------------------------------

func TestIntervalUnits(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"DAY", "SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM t"},
		{"HOUR", "SELECT DATE_ADD(d, INTERVAL 1 HOUR) FROM t"},
		{"MINUTE", "SELECT DATE_ADD(d, INTERVAL 1 MINUTE) FROM t"},
		{"SECOND", "SELECT DATE_ADD(d, INTERVAL 1 SECOND) FROM t"},
		{"MONTH", "SELECT DATE_ADD(d, INTERVAL 1 MONTH) FROM t"},
		{"WEEK", "SELECT DATE_ADD(d, INTERVAL 1 WEEK) FROM t"},
		{"QUARTER", "SELECT DATE_ADD(d, INTERVAL 1 QUARTER) FROM t"},
		{"YEAR", "SELECT DATE_ADD(d, INTERVAL 1 YEAR) FROM t"},
		{"MICROSECOND", "SELECT DATE_ADD(d, INTERVAL 1 MICROSECOND) FROM t"},
		{"HOUR_MINUTE", "SELECT DATE_ADD(d, INTERVAL '1:30' HOUR_MINUTE) FROM t"},
		{"DAY_SECOND", "SELECT DATE_ADD(d, INTERVAL '1 1:30:00' DAY_SECOND) FROM t"},
		{"YEAR_MONTH", "SELECT DATE_ADD(d, INTERVAL '1-6' YEAR_MONTH) FROM t"},
		{"HOUR_MICROSECOND", "SELECT DATE_ADD(d, INTERVAL '1:30:00.5' HOUR_MICROSECOND) FROM t"},
		{"DAY_HOUR", "SELECT DATE_ADD(d, INTERVAL '1 12' DAY_HOUR) FROM t"},
		{"DAY_MINUTE", "SELECT DATE_ADD(d, INTERVAL '1 12:30' DAY_MINUTE) FROM t"},
		{"MINUTE_SECOND", "SELECT DATE_ADD(d, INTERVAL '30:00' MINUTE_SECOND) FROM t"},
		{"MINUTE_MICROSECOND", "SELECT DATE_ADD(d, INTERVAL '30:00.5' MINUTE_MICROSECOND) FROM t"},
		{"SECOND_MICROSECOND", "SELECT DATE_ADD(d, INTERVAL '59.999999' SECOND_MICROSECOND) FROM t"},
		{"HOUR_SECOND", "SELECT DATE_ADD(d, INTERVAL '1:30:59' HOUR_SECOND) FROM t"},
		{"interval_add", "SELECT d + INTERVAL 1 DAY FROM t"},
		{"interval_sub", "SELECT d - INTERVAL 1 MONTH FROM t"},
		// YEAR is a registered keyword token (kwYEAR) — parseKeywordOrIdent accepts it.
		{"keyword_YEAR", "SELECT DATE_ADD(d, INTERVAL 1 YEAR) FROM t"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseAndCheck(t, tt.sql)
		})
	}
}

// TestIntervalUnitValidation verifies that invalid INTERVAL unit names are rejected.
func TestIntervalUnitValidation(t *testing.T) {
	cases := []string{
		"SELECT DATE_ADD(d, INTERVAL 1 BOGUS) FROM t",
		"SELECT DATE_ADD(d, INTERVAL 1 FOOBAR) FROM t",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			if err == nil {
				t.Fatalf("expected error for invalid interval unit, got nil")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Section 2.4: UTC Temporal Functions
// ---------------------------------------------------------------------------

func TestUTCTemporalFunctions(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"UTC_DATE_no_parens", "SELECT UTC_DATE"},
		{"UTC_DATE_parens", "SELECT UTC_DATE()"},
		{"UTC_TIME_no_parens", "SELECT UTC_TIME"},
		{"UTC_TIME_parens", "SELECT UTC_TIME()"},
		{"UTC_TIME_fsp", "SELECT UTC_TIME(3)"},
		{"UTC_TIMESTAMP_no_parens", "SELECT UTC_TIMESTAMP"},
		{"UTC_TIMESTAMP_parens", "SELECT UTC_TIMESTAMP()"},
		{"UTC_TIMESTAMP_fsp", "SELECT UTC_TIMESTAMP(6)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseAndCheck(t, tt.sql)
		})
	}
}

// TestLvalueIdentRejectsAmbiguous4 verifies that ambiguous_4 keywords (like GLOBAL)
// cannot be used as SET variable names. GLOBAL is allowed as a general identifier
// (e.g., CREATE TABLE global) but NOT as a lvalue_ident (SET target).
func TestLvalueIdentRejectsAmbiguous4(t *testing.T) {
	// GLOBAL as a table name should work (general ident context).
	ParseAndCheck(t, "CREATE TABLE global (a INT)")

	// SET global = 1 should fail — GLOBAL is ambiguous_4 (scope modifier, not variable).
	// In MySQL, `SET global = 1` is parsed as `SET GLOBAL = 1` which is incomplete
	// (missing variable name after scope), so it's an error.
	t.Run("SET global = 1 rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("SET global = 1")}
		p.advance()
		_, err := p.parseSetStmt()
		if err == nil {
			t.Fatal("expected error for SET global = 1 (GLOBAL is ambiguous_4, consumed as scope modifier)")
		}
	})

	// SET begin = 1 should succeed — BEGIN is ambiguous_2 (allowed as lvalue).
	t.Run("SET begin = 1 accepted", func(t *testing.T) {
		ParseAndCheck(t, "SET begin = 1")
	})

	// RESET PERSIST global should fail — GLOBAL is ambiguous_4 (not lvalue).
	t.Run("RESET PERSIST global not consumed as variable", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("RESET PERSIST global")}
		p.advance()
		stmt, err := p.parseStmt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// The RESET PERSIST should parse with no variable name since GLOBAL
		// is not accepted by isLvalueIdentToken.
		rp, ok := stmt.(*ast.ResetPersistStmt)
		if !ok {
			t.Fatalf("expected ResetPersistStmt, got %T", stmt)
		}
		if rp.Variable != "" {
			t.Fatalf("expected empty variable (GLOBAL should not be consumed as lvalue), got %q", rp.Variable)
		}
	})

	// RESET PERSIST begin should succeed — BEGIN is ambiguous_2 (allowed as lvalue).
	t.Run("RESET PERSIST begin accepted", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("RESET PERSIST begin")}
		p.advance()
		stmt, err := p.parseStmt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		rp, ok := stmt.(*ast.ResetPersistStmt)
		if !ok {
			t.Fatalf("expected ResetPersistStmt, got %T", stmt)
		}
		if rp.Variable != "begin" {
			t.Fatalf("expected variable 'begin', got %q", rp.Variable)
		}
	})

	// FETCH INTO with lvalue context — BEGIN (ambiguous_2) is allowed.
	t.Run("FETCH INTO begin accepted", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("FETCH cur INTO begin")}
		p.advance()
		stmt, err := p.parseFetchCursorStmt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(stmt.Into) != 1 || stmt.Into[0] != "begin" {
			t.Fatalf("expected INTO [begin], got %v", stmt.Into)
		}
	})

	// FETCH INTO with lvalue context — GLOBAL (ambiguous_4) is NOT allowed.
	t.Run("FETCH INTO global rejected", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("FETCH cur INTO global")}
		p.advance()
		_, err := p.parseFetchCursorStmt()
		if err == nil {
			t.Fatal("expected error for FETCH INTO global (GLOBAL is ambiguous_4, not allowed as lvalue)")
		}
	})
}
