package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTSql160G4IsTheConversion holds the committed TSql160.g4 to what the
// converter writes from the committed vendor files, so neither is edited
// alone.
func TestTSql160G4IsTheConversion(t *testing.T) {
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
		t.Errorf("mssql/testdata/%s is not the conversion of %s, %s and %s. Rerun it:\n  go run ./mssql/parser/cmd/scriptdomg4", output, grammarFile, tokenTypesFile, supporterFile)
	}
}

// TestConvert pins what the conversion keeps, drops, rewrites and notes.
func TestConvert(t *testing.T) {
	const g = `options { language = "CSharp"; }
{ using System; }
class P extends Parser("B");
options { k = 2; }

script returns [S vResult = null]
{ int n; }
    : (vResult = stmt)* EOF
    ;

stmt returns [S vResult = null]
    : tSet:Identifier { Match(tSet, CodeGenerationSupporter.Set); } opt
    | (Select Star)=> Select Star
    | tUnit:Identifier { Match(tUnit, CodeGenerationSupporter.GB, CodeGenerationSupporter.MB); }
    | tKind:Identifier { if (TryMatch(tKind, CodeGenerationSupporter.Row)) n++; else Match(tKind, CodeGenerationSupporter.Range); }
    | Select tBad:Identifier { ThrowParseErrorException("SQL1", tBad, "no"); }
    | odbc
    ;
    exception
    catch[TSqlParseErrorException e] { Recover(); }

opt
    : {NextTokenMatches(CodeGenerationSupporter.Gb)}? sizeOpt
    | {LA(2) == Star}? Star
    | ( { if (!SkipGuessing(t)) }: Select )
    ;

sizeOpt : tName:Identifier EqualsSign Integer ;

odbc : tO:Odbc { ThrowParseErrorException("SQL2", tO, "no"); } ;

class L extends Lexer("B");
tokens {
    Odbc = "odbc";
}

protected
Integer :;

Number : (Digit)+ { $setType(Integer); } ;

protected
Digit : '0'..'9' ;

WhiteSpace : (' ' | '\t')+ ;
`
	const tokenTypes = `class TSql extends Lexer;
tokens {
    Select = "select";
    Star;
}
`
	const supporter = `internal const string Set = "SET";
internal const string GB = "GB";
internal const string MB = "MB";
internal const string Gb = "GB";
internal const string Row = "ROW";
internal const string Range = "RANGE";`
	got, err := convert(source{"T.g", g}, source{"K.g", tokenTypes}, source{"C.cs", supporter})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"grammar TSql160;\n\noptions { caseInsensitive = true; }\n\n// T.g:6\nscript\n    : ( stmt )* EOF\n",
		// A Match at an action's top level writes the word.
		"    : 'SET' opt\n",
		"    | Select Star\n",
		// Two words, and the words of a TryMatch chain.
		"    | ( 'GB' | 'MB' )\n",
		"    | ( 'ROW' | 'RANGE' )\n",
		"// T.g:16: alternative dropped, its action refuses every input: Select tBad : Identifier {ThrowParseErrorException(\"SQL1\", tBad, \"no\");}\n",
		// A predicate on the next word writes the rule it gates as that word.
		"// T.g:23: predicate written as its words: {NextTokenMatches(CodeGenerationSupporter.Gb)}?\n",
		"    : sizeOpt_GB\n",
		"// T.g:24: predicate dropped: {LA(2) == Star}?\n",
		"sizeOpt_GB\n    : 'GB' EqualsSign Integer\n",
		"// Keywords:",
		"// K.g:3\nSelect\n    : 'SELECT'\n",
		"// T.g:34\nOdbc\n    : 'ODBC'\n",
		"// T.g:40: this alternative may emit Integer instead ($setType)\nNumber\n    : ( Digit )+\n",
		"fragment Digit\n    : '0'..'9'\n",
		"WhiteSpace\n    : ( ' ' | '\\t' )+\n      -> channel(HIDDEN)\n",
		"// Integer: T.g:38\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the conversion lacks %q:\n%s", want, got)
		}
	}
	for _, gone := range []string{
		// odbc refuses every input, so is a rule no more, and the
		// alternative of stmt that needs it goes with it; sizeOpt is
		// reached only as its copy.
		"\nodbc\n", "    | odbc\n", "\nsizeOpt\n",
		"Star\n    :",
	} {
		if strings.Contains(got, gone) {
			t.Errorf("the conversion holds %q:\n%s", gone, got)
		}
	}
	if _, err := convert(source{"T.g", "class P extends Parser;\nstmt : A ;"}, source{"K.g", tokenTypes}, source{"C.cs", ""}); err == nil || !strings.Contains(err.Error(), "no rule script") {
		t.Errorf("a grammar with no start rule: %v", err)
	}
	if _, err := convert(source{"T.g", "class P extends Parser;\nscript : A { refuse( ;"}, source{"K.g", tokenTypes}, source{"C.cs", ""}); err == nil {
		t.Error("an action never closed converted")
	}
}
