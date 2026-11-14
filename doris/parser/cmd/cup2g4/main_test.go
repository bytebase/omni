package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDorisLegacyParserIsTheConversion holds the committed .g4 to what the
// converter writes from the committed CUP and JFlex files, so neither is
// edited alone.
func TestDorisLegacyParserIsTheConversion(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "testdata")
	want, err := convertDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, output))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("doris/testdata/%s is not the conversion of %s and %s. Rerun it:\n  go run ./doris/parser/cmd/cup2g4", output, cupFile, flexFile)
	}
}

// TestConvert pins what the conversion keeps and drops.
func TestConvert(t *testing.T) {
	const cup = `package x;
import java.util.List;
parser code {: int n; :};
terminal String KW_SET, KW_GLOBAL, KW_AND, KW_TREE, IDENT;
terminal COMMA, EQUAL;
nonterminal List<String> stmts, set_stmt, opt_type, unused;
precedence left KW_AND;
start with stmts;

stmts ::=
    set_stmt:s
    {: RESULT = Lists.newArrayList(s); :}
    | stmts:l COMMA set_stmt:s %prec KW_AND
    | error
    ;
set_stmt ::= KW_SET opt_type IDENT:v EQUAL IDENT KW_AND KW_TREE ;
opt_type ::=
    /* empty */
    | KW_GLOBAL
    ;
unused ::= KW_SET ;
`
	const flex = `
        keywordMap.put("&&", new Integer(SqlParserSymbols.KW_AND));
        keywordMap.put("and", new Integer(SqlParserSymbols.KW_AND));
        keywordMap.put("global", new Integer(SqlParserSymbols.KW_GLOBAL));
        keywordMap.put("set", new Integer(SqlParserSymbols.KW_SET));
"," { return newToken(SqlParserSymbols.COMMA, null); }
"=" { return newToken(SqlParserSymbols.EQUAL, null); }
`
	got, err := convert("// header\n", cup, flex)
	if err != nil {
		t.Fatal(err)
	}
	want := `// header
// sql_parser.cup declares terminals sql_scanner.flex never makes: KW_TREE.

parser grammar DorisLegacyParser;

// sql_parser.cup:10
stmts
    : set_stmt
    | stmts ',' set_stmt
    ;

// sql_parser.cup:16
set_stmt
    : 'SET' opt_type IDENTIFIER '=' IDENTIFIER 'AND' KW_TREE
    ;

// sql_parser.cup:17
opt_type
    : /* empty */
    | 'GLOBAL'
    ;
`
	if got != want {
		t.Errorf("convert:\n%s\nwant:\n%s", got, want)
	}
	if !strings.Contains(got, "stmts") || strings.Contains(got, "unused") {
		t.Errorf("an unreached rule was kept")
	}
}
