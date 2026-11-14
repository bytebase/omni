package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGoogleSQLG4IsTheConversion holds the committed .g4 to what the
// converter writes from the committed googlesql.tm, so neither is edited
// alone.
func TestGoogleSQLG4IsTheConversion(t *testing.T) {
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
		t.Errorf("googlesql/testdata/%s is not the conversion of %s. Rerun it:\n  go run ./googlesql/parser/cmd/tm2g4", output, tmFile)
	}
}

// TestConvert pins what the conversion keeps and drops.
func TestConvert(t *testing.T) {
	const tm = `language sql(cc);

:: lexer

decimal_digits = /[0-9]+/
IDENTIFIER {absl::string_view}: /[a-z_]+/ -1
INTEGER_LITERAL {absl::string_view}:
"," (COMMA): /,/
"(" (LPAREN): /\(/
")" (RPAREN): /\)/
"<>" (KW_NOT_EQUALS_SQL_STYLE):
"SELECT" (KW_SELECT): /select/
"UNION" (KW_UNION): /union/
"ALL" (KW_ALL): /all/
"NOT_SPECIAL" (KW_NOT_SPECIAL):
KW_OPEN_HINT:
whitespace: /[ \n]+/ (space)

:: parser

%left "UNION";
%input stmt;

stmt {ASTNode*}:
    "SELECT" (item separator ",")+[items] tail_opt
    {
      // A } in a comment, and "}" in a string: '}'.
      $$ = MakeNode<ASTSelect>(@$, $items);
    }
  | stmt "UNION" "ALL" stmt %prec "UNION"
  | KW_OPEN_HINT INTEGER_LITERAL stmt
;

# A comment.
item:
    IDENTIFIER[name] ("<>" | "NOT_SPECIAL")? INTEGER_LITERAL*
;

tail {std::vector<int>}:
    %empty
  | "(" (item separator "UNION" "ALL")* ")"
;

options: "ALL" ;

unused: "ALL" ;

%%
{{define "parserPrivateDecls"}}{{end}}
`
	got, classes, err := convert("// header\n", tm, "stmt")
	if err != nil {
		t.Fatal(err)
	}
	want := `// header
grammar GoogleSQL;

// googlesql.tm:24
stmt
    : 'SELECT' ((item) (',' (item))*) tail?
    | stmt 'UNION' 'ALL' stmt
    | KW_OPEN_HINT INTEGER_LITERAL stmt
    ;

// googlesql.tm:35
item
    : IDENTIFIER ('<>' | KW_NOT_SPECIAL)? INTEGER_LITERAL*
    ;

// googlesql.tm:39
tail
    : /* empty */
    | '(' ((item) ('UNION' 'ALL' (item))*)? ')'
    ;
`
	if got != want {
		t.Errorf("convert:\n%s\nwant:\n%s", got, want)
	}
	if want := []string{"IDENTIFIER", "INTEGER_LITERAL", "KW_NOT_SPECIAL", "KW_OPEN_HINT"}; len(classes) != len(want) || classes[0] != want[0] || classes[3] != want[3] {
		t.Errorf("classes = %v, want %v", classes, want)
	}
}
