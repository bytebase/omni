package parser_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/parser"
)

func TestCompletionPublicTokenize(t *testing.T) {
	sql := "SELECT [select], dbo.Users FROM [Order Details] WHERE id = @id"

	tokens := parser.Tokenize(sql)
	if len(tokens) == 0 {
		t.Fatal("Tokenize returned no tokens")
	}
	if got := parser.TokenName(tokens[0].Type); got != "SELECT" {
		t.Fatalf("first token = %q, want SELECT", got)
	}
	if tokens[len(tokens)-1].Str != "@id" {
		t.Fatalf("last token string = %q, want @id", tokens[len(tokens)-1].Str)
	}
	for _, tok := range tokens {
		if tok.Type == 0 {
			t.Fatalf("Tokenize returned EOF token: %+v", tok)
		}
		if tok.End <= tok.Loc {
			t.Fatalf("token has invalid location: %+v", tok)
		}
	}
}

func TestCompletionPublicIdentifierTokenPredicate(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{name: "regular identifier", sql: "customer", want: true},
		{name: "bracketed reserved word", sql: "[FROM]", want: true},
		{name: "double quoted identifier", sql: `"FROM"`, want: true},
		{name: "context keyword as identifier", sql: "LOGON", want: true},
		{name: "core keyword not bare identifier", sql: "FROM", want: false},
		{name: "variable is not identifier", sql: "@id", want: false},
		{name: "string literal is not identifier", sql: "'id'", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := parser.Tokenize(tt.sql)
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) returned %d tokens, want 1", tt.sql, len(tokens))
			}
			if got := parser.IsIdentTokenType(tokens[0].Type); got != tt.want {
				t.Fatalf("IsIdentTokenType(%q/%d) = %v, want %v", tokens[0].Str, tokens[0].Type, got, tt.want)
			}
		})
	}
}

func TestCompletionPublicTokenConstants(t *testing.T) {
	required := []struct {
		token int
		name  string
	}{
		{parser.SELECT, "SELECT"},
		{parser.INSERT, "INSERT"},
		{parser.UPDATE, "UPDATE"},
		{parser.DELETE, "DELETE"},
		{parser.FROM, "FROM"},
		{parser.WHERE, "WHERE"},
		{parser.SET, "SET"},
		{parser.INTO, "INTO"},
		{parser.VALUES, "VALUES"},
		{parser.AS, "AS"},
		{parser.ON, "ON"},
		{parser.JOIN, "JOIN"},
		{parser.INNER, "INNER"},
		{parser.LEFT, "LEFT"},
		{parser.RIGHT, "RIGHT"},
		{parser.CROSS, "CROSS"},
		{parser.FULL, "FULL"},
		{parser.OUTER, "OUTER"},
		{parser.ORDER, "ORDER"},
		{parser.GROUP, "GROUP"},
		{parser.HAVING, "HAVING"},
		{parser.UNION, "UNION"},
		{parser.FOR, "FOR"},
		{parser.WITH, "WITH"},
		{parser.GO, "GO"},
		{parser.EXEC, "EXEC"},
		{parser.EXECUTE, "EXECUTE"},
		{parser.DATABASE, "DATABASE"},
		{parser.SCHEMA, "SCHEMA"},
		{parser.TABLE, "TABLE"},
		{parser.VIEW, "VIEW"},
		{parser.SEQUENCE, "SEQUENCE"},
		{parser.COLUMN, "COLUMN"},
		{parser.INDEX, "INDEX"},
		{parser.TRIGGER, "TRIGGER"},
		{parser.FUNCTION, "FUNCTION"},
		{parser.PROCEDURE, "PROCEDURE"},
		{parser.BY, "BY"},
		{parser.ASC, "ASC"},
		{parser.DESC, "DESC"},
		{parser.OFFSET, "OFFSET"},
		{parser.FETCH, "FETCH"},
		{parser.NEXT, "NEXT"},
		{parser.FIRST, "FIRST"},
		{parser.ROW, "ROW"},
		{parser.ROWS, "ROWS"},
		{parser.ONLY, "ONLY"},
		{parser.NULL, "NULL"},
		{parser.NOT, "NOT"},
		{parser.AND, "AND"},
		{parser.OR, "OR"},
		{parser.CASE, "CASE"},
		{parser.WHEN, "WHEN"},
		{parser.THEN, "THEN"},
		{parser.ELSE, "ELSE"},
		{parser.END, "END"},
		{parser.OVER, "OVER"},
		{parser.PARTITION, "PARTITION"},
		{parser.XML, "XML"},
		{parser.JSON, "JSON"},
		{parser.BROWSE, "BROWSE"},
		{parser.PATH, "PATH"},
		{parser.AUTO, "AUTO"},
		{parser.RAW, "RAW"},
		{parser.RECOMPILE, "RECOMPILE"},
		{parser.OPTIMIZE, "OPTIMIZE"},
		{parser.MAXDOP, "MAXDOP"},
		{parser.NOLOCK, "NOLOCK"},
	}

	for _, tt := range required {
		if got := parser.TokenName(tt.token); got != tt.name {
			t.Errorf("TokenName(%d) = %q, want %q", tt.token, got, tt.name)
		}
	}
}

func TestCompletionPrefixRetryAtIdentifierStart(t *testing.T) {
	sql := "SELECT * FROM dbo.Us"
	cursor := len(sql)
	tokens := parser.Tokenize(sql)

	var prefix parser.Token
	found := false
	for _, tok := range tokens {
		if tok.End == cursor && parser.IsIdentTokenType(tok.Type) {
			prefix = tok
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("failed to find identifier prefix token at cursor in %q", sql)
	}

	candidates := parser.Collect(sql, prefix.Loc)
	if !candidates.HasRule("table_ref") {
		t.Fatalf("Collect at prefix start did not return table_ref for %q", sql)
	}
}

func TestCompletionRecordsCTEPositionsForBytebase(t *testing.T) {
	sql := "WITH cte AS (SELECT id FROM dbo.Users) SELECT * FROM cte WHERE "
	candidates := parser.Collect(sql, len(sql))

	want := strings.Index(sql, "WITH")
	if !slices.Contains(candidates.CTEPositions, want) {
		t.Fatalf("CTEPositions = %v, want to contain %d", candidates.CTEPositions, want)
	}
}

func TestCompletionRecordsSelectAliasPositionsForBytebase(t *testing.T) {
	sql := "SELECT id AS c_id, name implicit_name, eq_name = id FROM dbo.Users ORDER BY "
	candidates := parser.Collect(sql, len(sql))

	want := []int{
		strings.Index(sql, "c_id"),
		strings.Index(sql, "implicit_name"),
		strings.Index(sql, "eq_name"),
	}
	for _, pos := range want {
		if !slices.Contains(candidates.SelectAliasPositions, pos) {
			t.Fatalf("SelectAliasPositions = %v, want to contain %d", candidates.SelectAliasPositions, pos)
		}
	}
}

func TestCompletionAfterGoBatchSeparator(t *testing.T) {
	sql := "SELECT 1\nGO\n"
	candidates := parser.Collect(sql, len(sql))

	required := []int{parser.SELECT, parser.INSERT, parser.UPDATE, parser.DELETE}
	for _, tok := range required {
		if !candidates.HasToken(tok) {
			t.Fatalf("after GO: missing top-level keyword %s", parser.TokenName(tok))
		}
	}
}
