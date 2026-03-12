package parser

import (
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
		{"SELECT @@global.var", []int{kwSELECT, tokIDENT}},
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
		{"1 + 2", "{BINEXPR :op + :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 - 2", "{BINEXPR :op - :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 * 2", "{BINEXPR :op * :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 / 2", "{BINEXPR :op / :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"10 DIV 3", "{BINEXPR :op DIV :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 7}}"},
		{"10 % 3", "{BINEXPR :op %% :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 5}}"},
		{"10 MOD 3", "{BINEXPR :op %% :left {INT_LIT :val 10 :loc 0} :right {INT_LIT :val 3 :loc 7}}"},
		{"-1", "{UNARY :op - :operand {INT_LIT :val 1 :loc 1}}"},
		{"~5", "{UNARY :op ~ :operand {INT_LIT :val 5 :loc 1}}"},
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
		{"1 = 2", "{BINEXPR :op = :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 <> 2", "{BINEXPR :op <> :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 != 2", "{BINEXPR :op <> :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 < 2", "{BINEXPR :op < :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 > 2", "{BINEXPR :op > :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 4}}"},
		{"1 <= 2", "{BINEXPR :op <= :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 >= 2", "{BINEXPR :op >= :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"1 <=> 2", "{BINEXPR :op <=> :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 6}}"},
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
		{"1 AND 2", "{BINEXPR :op AND :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 6}}"},
		{"1 OR 2", "{BINEXPR :op OR :left {INT_LIT :val 1 :loc 0} :right {INT_LIT :val 2 :loc 5}}"},
		{"NOT 1", "{UNARY :op NOT :operand {INT_LIT :val 1 :loc 4}}"},
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
	stmt, err := p.parseSetStmt()
	if err != nil {
		t.Fatalf("parseSetStmt(%q) error: %v", input, err)
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
}

// ============================================================================
// Batch 15: GRANT, REVOKE, CREATE/DROP USER
// ============================================================================

func TestParseGrant(t *testing.T) {
	t.Run("all privileges", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT ALL PRIVILEGES ON *.* TO root")}
		p.advance()
		stmt, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
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
		stmt, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(stmt.Privileges) != 2 {
			t.Fatalf("Privileges count = %d, want 2", len(stmt.Privileges))
		}
	})

	t.Run("with grant option", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("GRANT ALL ON *.* TO admin WITH GRANT OPTION")}
		p.advance()
		stmt, err := p.parseGrantStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.WithGrant {
			t.Errorf("WithGrant = false, want true")
		}
	})
}

func TestParseRevoke(t *testing.T) {
	p := &Parser{lexer: NewLexer("REVOKE SELECT ON mydb.* FROM app_user")}
	p.advance()
	stmt, err := p.parseRevokeStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
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
		stmt, err := p.parseLoadDataStmt()
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
		stmt, err := p.parseLoadDataStmt()
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
		stmt, err := p.parseLoadDataStmt()
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
		stmt, err := p.parseLoadDataStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.FieldsTerminatedBy != "," {
			t.Errorf("FieldsTerminatedBy = %s, want ,", stmt.FieldsTerminatedBy)
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

func TestParseKill(t *testing.T) {
	t.Run("connection", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("KILL 123")}
		p.advance()
		stmt, err := p.parseKillStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if stmt.ConnectionID == nil {
			t.Fatal("ConnectionID is nil")
		}
	})

	t.Run("query", func(t *testing.T) {
		p := &Parser{lexer: NewLexer("KILL QUERY 123")}
		p.advance()
		stmt, err := p.parseKillStmt()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if !stmt.Query {
			t.Errorf("Query = false, want true")
		}
	})
}

func TestParseDo(t *testing.T) {
	p := &Parser{lexer: NewLexer("DO 1 + 2, 3 + 4")}
	p.advance()
	stmt, err := p.parseDoStmt()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(stmt.Exprs) != 2 {
		t.Fatalf("Exprs count = %d, want 2", len(stmt.Exprs))
	}
}

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
			name:  "all statement types via Parse()",
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
DO 1`,
			count: 35,
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
