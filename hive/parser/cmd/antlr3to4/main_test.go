package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestHiveG4IsTheConversion holds the committed Hive.g4 to what the converter
// writes from the committed ANTLR 3 files, so neither is edited alone.
func TestHiveG4IsTheConversion(t *testing.T) {
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
		t.Errorf("hive/testdata/%s is not the conversion of %s. Rerun it:\n  go run ./hive/parser/cmd/antlr3to4", output, strings.Join(vendor, ", "))
	}
}

// TestConvert pins what the conversion keeps, drops, and notes.
func TestConvert(t *testing.T) {
	const parser = `parser grammar P;
options { output = AST; }
tokens { TOK_SELECT; }
@members { int n; }
@rulecatch { catch (RecognitionException e) { throw e; } }

statement
@init { pushMsg("stmt"); }
@after { popMsg(); }
    : KW_SELECT (KW_DISTINCT)=> dist=KW_DISTINCT? col=name[true]
      -> {$dist != null}? ^(TOK_SELECTDI $col)
      -> ^(TOK_SELECT $col)
    | KW_DROP (KW_TABLE name -> ^(TOK_DROP name))
    ;

name[boolean b]
    : Identifier
    | {b}? KW_SELECT
    ;
`
	const lexerSub = `lexer grammar LSub;
import LParent;
StringLiteral : ('\'' (~('\''|'\\') | '\\' .)* '\'')+ ;
`
	const lexerParent = `lexer grammar LParent;
KW_SELECT : 'SELECT' ;
KW_DISTINCT : 'DISTINCT' ;
KW_DROP : 'DROP' ;
KW_TABLE : 'TABLE' ;
StringLiteral : ;
Identifier : ('a'..'z')+ ;
WS : (' ')+ { $channel = HIDDEN; } ;
QUERY_HINT : '/*' (options { greedy=false; } : QUERY_HINT | .)* '*/' { if (getText().charAt(2) != '+') { $channel = HIDDEN; } } ;
`
	got, err := convert("H", "// header\n", source{"P.g", parser}, source{"LSub.g", lexerSub}, source{"LParent.g", lexerParent})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"// header\ngrammar H;\n\n// P.g:7\nstatement\n    : KW_SELECT KW_DISTINCT? name\n    | KW_DROP ( KW_TABLE name )\n    ;\n",
		"// P.g:18: predicate dropped: {b}?\nname\n    : Identifier\n    | KW_SELECT\n    ;\n",
		"// LSub.g:3\nStringLiteral\n",
		"WS\n    : ( ' ' )+\n      -> channel(HIDDEN)\n",
		"// LParent.g:9: conditional hidden-channel action dropped:",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the conversion lacks %q:\n%s", want, got)
		}
	}
	// Empty StringLiteral stub from LParent.g must not appear twice.
	if strings.Count(got, "\nStringLiteral\n") != 1 {
		t.Errorf("expected StringLiteral once, got:\n%s", got)
	}
	// QUERY_HINT has conditional $channel=HIDDEN and must NOT get -> channel(HIDDEN).
	if strings.Contains(got, "QUERY_HINT\n    : '/*'") && strings.Contains(got[strings.Index(got, "QUERY_HINT"):], "-> channel(HIDDEN)") {
		t.Errorf("QUERY_HINT must not become unconditionally hidden:\n%s", got)
	}
}
