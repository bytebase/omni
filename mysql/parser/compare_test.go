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
