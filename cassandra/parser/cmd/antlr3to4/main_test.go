package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCqlG4IsTheConversion holds the committed Cql.g4 to what the converter
// writes from the committed ANTLR 3 files, so neither is edited alone.
func TestCqlG4IsTheConversion(t *testing.T) {
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
		t.Errorf("cassandra/testdata/%s is not the conversion of %s. Rerun it:\n  go run ./cassandra/parser/cmd/antlr3to4", output, strings.Join(vendor, ", "))
	}
}

// TestConvert pins what the conversion keeps, drops and notes.
func TestConvert(t *testing.T) {
	const parser = `parser grammar P;
options { language = Java; }
@members { int n; }

stmt returns [Stmt s]
    @init { n = 0; }
    : K_SELECT (K_JSON sel)=> K_JSON? s=sel { $s = new Stmt(); }
    | K_DROP QMARK { addRecognitionError("no bind variables here"); }
    | K_DROP n=name[true] { if (n == null) addRecognitionError("bad"); }
    ;

sel : '\*' | {input.LT(1) != null}? name (',' name)* ;
name[boolean b] : IDENT | K_JSON ;
`
	const lexer = `lexer grammar L;
K_SELECT : S E L E C T ;
K_JSON : 'json' ;
K_DROP : 'DROP' ;
QMARK : '?' ;
IDENT : ('a'..'z')+ ;
F : INT '.' | INT {$type = INT;} ;
fragment INT : '0'..'9'+ ;
WS : (' ')+ { $channel = HIDDEN; } ;
COMMENT : '--' .* '\n' { $channel = HIDDEN; } ;
fragment S: ('s'|'S'); fragment E: ('e'|'E'); fragment L: ('l'|'L'); fragment C: ('c'|'C'); fragment T: ('t'|'T');
`
	got, err := convert("G", "// header\n", source{"P.g", parser}, source{"L.g", lexer})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"// header\ngrammar G;\n\n// P.g:5\n// P.g:8: alternative dropped",
		"stmt\n    : K_SELECT K_JSON? sel\n",
		"// P.g:8: alternative dropped, its action refuses every input: K_DROP QMARK {addRecognitionError(\"no bind variables here\");}\n",
		"// P.g:9: action may refuse the input: {if (n == null) addRecognitionError(\"bad\");}\n",
		"    | K_DROP name\n    ;\n",
		"// P.g:12\n// P.g:12: predicate dropped: {input.LT(1) != null}?\n",
		"sel\n    : '*'\n    | name ( ',' name )*\n    ;\n",
		"name\n    : IDENT\n    | K_JSON\n    ;\n",
		"// L.g:7: action dropped, it emits another token type: {$type = INT;}\n",
		"fragment INT\n    : '0'..'9'+\n",
		"WS\n    : ( ' ' )+\n      -> channel(HIDDEN)\n",
		"COMMENT\n    : '--' .*? '\\n'\n      -> channel(HIDDEN)\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the conversion lacks %q:\n%s", want, got)
		}
	}
	// Parser rules come before lexer rules, whatever file holds them.
	if strings.Index(got, "\nstmt\n") > strings.Index(got, "\nK_SELECT\n") {
		t.Errorf("a lexer rule before the parser rules:\n%s", got)
	}
	if _, err := convert("G", "", source{"L.g", lexer}); err == nil || !strings.Contains(err.Error(), "no parser rule") {
		t.Errorf("a lexer grammar alone: %v", err)
	}
	if _, err := convert("G", "", source{"P.g", "p : A { refuse( ;"}); err == nil {
		t.Error("an action never closed converted")
	}
}
