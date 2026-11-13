package grammar

import (
	"strings"
	"testing"
)

func antlrRows(t *testing.T, src string, forms map[string]string) ([]Row, []string) {
	t.Helper()
	g, err := ReadANTLR(src, tokens(forms))
	if err != nil {
		t.Fatal(err)
	}
	return g.Rows()
}

const warehouseG4 = `
grammar Warehouse;
options { caseInsensitive = true; }
@header { package x; }

// https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse
alterWarehouse
    : ALTER WAREHOUSE (IF EXISTS)? name? (SUSPEND | RESUME (IF SUSPENDED)?)   # suspend
    | ALTER WAREHOUSE n=name SET TAG name '=' STRING (',' name '=' STRING)*  # setTag
    ;
name returns [String v] : ID {$v = $ID.text;} ;

ALTER : 'ALTER' ;
WAREHOUSE : W A R E H O U S E ;
IF : 'IF' ; EXISTS : 'EXISTS' ; SUSPEND : 'SUSPEND' ; SUSPENDED : 'SUSPENDED' ;
RESUME : 'RESUME' ; SET : 'SET' ; TAG : 'TAG' ;
ID : [a-z_] [a-z0-9_]* ;
STRING : '\'' ( ~'\'' | '\'\'' )* '\'' ;
WS : [ \t\r\n]+ -> channel(HIDDEN) ;
fragment A : [aA] ; fragment E : [eE] ; fragment H : [hH] ; fragment O : [oO] ;
fragment R : [rR] ; fragment S : [sS] ; fragment U : [uU] ; fragment W : [wW] ;
`

// TestANTLRGroupsBecomeRulesNamedByTheirText pins how EBNF is lowered: each
// group or operator is a rule named after the rule that holds it and the
// group as written, so each optional part and each choice has rows of its
// own, and labels stay stable while the group does.
func TestANTLRGroupsBecomeRulesNamedByTheirText(t *testing.T) {
	got, uncovered := antlrRows(t, warehouseG4, map[string]string{"ID": "a"})
	labels := map[string]string{}
	for _, r := range got {
		if !IsNearMiss(r.Rule) {
			labels[r.Rule] = r.Statement
		}
	}
	for label, want := range map[string]string{
		// Over 64 characters, the symbols are cut and end in a hash.
		"alterWarehouse/ALTER_WAREHOUSE_(IF_EXISTS)?_(name)?_(SUSPEND|RESUME_(I~18f316a7": "ALTER WAREHOUSE SUSPEND",
		"alterWarehouse(IF_EXISTS)?/%empty":                                               "ALTER WAREHOUSE SUSPEND",
		"alterWarehouse(IF_EXISTS)?/IF_EXISTS":                                            "ALTER WAREHOUSE IF EXISTS SUSPEND",
		"alterWarehouse(name)?/name":                                                      "ALTER WAREHOUSE a SUSPEND",
		"alterWarehouse(SUSPEND|RESUME_(IF_SUSPENDED)?)/RESUME_(IF_SUSPENDED)?":           "ALTER WAREHOUSE RESUME",
		"alterWarehouse(IF_SUSPENDED)?/IF_SUSPENDED":                                      "ALTER WAREHOUSE RESUME IF SUSPENDED",
		"alterWarehouse(,_name_=_STRING)*/%empty":                                         "ALTER WAREHOUSE a SET TAG a = ''",
		"alterWarehouse(,_name_=_STRING)*/,_name_=_STRING_(,_name_=_STRING)*":             "ALTER WAREHOUSE a SET TAG a = '', a = ''",
		"name/ID": "ALTER WAREHOUSE a SUSPEND",
	} {
		if s, ok := labels[label]; !ok {
			t.Errorf("no row %s; rows:\n%s", label, render(got))
		} else if s != want {
			t.Errorf("%s is %q, want %q", label, s, want)
		}
	}
	if len(uncovered) != 0 {
		t.Errorf("uncovered: %v", uncovered)
	}
}

// TestANTLRLexerRulesSampleTheTextTheySpell pins where a token's text comes
// from: a literal, letters spelled through case-insensitive fragments, the
// shortest text of a pattern, or the token table, which wins.
func TestANTLRLexerRulesSampleTheTextTheySpell(t *testing.T) {
	g, err := ReadANTLR(warehouseG4, tokens(map[string]string{"STRING": "'s'"}))
	if err != nil {
		t.Fatal(err)
	}
	for tok, want := range map[string]string{"ALTER": "ALTER", "WAREHOUSE": "WAREHOUSE", "ID": "a", "STRING": "'s'", "EOF": ""} {
		if got, ok := g.termText[tok]; !ok || got != want {
			t.Errorf("%s is sampled %q (%v), want %q", tok, got, ok, want)
		}
	}
	if got, ok := g.termText["W"]; ok {
		t.Errorf("the fragment W is a token sampled %q", got)
	}
}

// TestANTLRRepeatedChoicesAreAChoiceRepeated pins (a | b)+: the choice is a
// rule with a row per choice, and the repetition one with a row for once
// and one for twice, so a list of n properties costs n+2 rows, not 2n.
func TestANTLRRepeatedChoicesAreAChoiceRepeated(t *testing.T) {
	got, _ := antlrRows(t, "grammar P; s : 'X' ('Y' | 'Z')+ 'W'* ;", nil)
	labels := map[string]string{}
	var alts []string
	for _, r := range got {
		labels[r.Rule] = r.Statement
		if !IsNearMiss(r.Rule) {
			alts = append(alts, r.Rule)
		}
	}
	want := "s/X_(Y|Z)+_(W)*,s(Y|Z)+/(Y|Z),s(Y|Z)+/(Y|Z)_(Y|Z)+,s(Y|Z)/Y,s(Y|Z)/Z,s(W)*/%empty,s(W)*/W_(W)*"
	if strings.Join(alts, ",") != want {
		t.Errorf("alternatives:\n got %s\nwant %s", strings.Join(alts, ","), want)
	}
	for label, stmt := range map[string]string{
		"s/X_(Y|Z)+_(W)*":        "X Y",
		"s(Y|Z)+/(Y|Z)_(Y|Z)+":   "X Y Y",
		"s(Y|Z)/Z":               "X Z",
		"s(W)*/W_(W)*":           "X Y W",
		"s/X_(Y|Z)+_(W)*!drop-2": "X",
	} {
		if labels[label] != stmt {
			t.Errorf("%s is %q, want %q", label, labels[label], stmt)
		}
	}
}

// TestANTLRLongGroupsAreCutAndHashed pins that a group's name stays short
// and holds nothing a label may not: no space, colon or slash.
func TestANTLRLongGroupsAreCutAndHashed(t *testing.T) {
	src := "grammar L; s : 'A' ('B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P')? ('::' 'a b' '/')? ;"
	got, _ := antlrRows(t, src, nil)
	for _, r := range got {
		rule, _, _ := strings.Cut(r.Rule, "/")
		if strings.ContainsAny(rule, " :\t") || len(rule) > 1+maxGroupName+2 {
			t.Errorf("label %q", r.Rule)
		}
	}
	if !strings.Contains(render(got), "s(::_a_b_/)?") && !strings.Contains(render(got), "s(COLONCOLON_a_b_SLASH)?") {
		t.Errorf("rows:\n%s", render(got))
	}
}

func TestANTLRRefusesWhatHasNoShortestText(t *testing.T) {
	for src, want := range map[string]string{
		"grammar X; s : ~'A' ;":         "~ in a parser rule",
		"grammar X; s : . ;":            ". in a parser rule",
		"grammar X; import Y; s : 'A';": "import",
		"grammar X; s : ID ;":           "ID",
	} {
		if _, err := ReadANTLR(src, tokens(nil)); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want an error naming %q", src, err, want)
		}
	}
}

// TestRecognizerDecidesWhatTheGrammarDerives pins the docs judge's reading:
// every generated alternative is derived, a near-miss is derived only when
// it lands in the language, and a token class takes any token of its kind
// but no reserved word.
func TestRecognizerDecidesWhatTheGrammarDerives(t *testing.T) {
	g, err := ReadANTLR(warehouseG4, tokens(map[string]string{"ID": "a", "STRING": "'s'", "ALTER": Reserved}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	rows, _ := g.Rows()
	for _, r := range rows {
		if !IsNearMiss(r.Rule) && !rec.Accepts(r.Statement) {
			t.Errorf("%s: the grammar's own statement is not derived: %s", r.Rule, r.Statement)
		}
	}
	for sql, want := range map[string]bool{
		"alter warehouse resume if suspended":                  true,
		`ALTER WAREHOUSE "My WH" SUSPEND`:                      true,
		"ALTER WAREHOUSE IF EXISTS SUSPEND":                    true,
		"ALTER WAREHOUSE b SET TAG c = 'x', d = 'it''s'":       true,
		"ALTER WAREHOUSE b SET TAG c = 'x' d = 'y'":            false,
		"ALTER WAREHOUSE alter SUSPEND":                        false, // reserved
		`ALTER WAREHOUSE "ALTER" SUSPEND`:                      true,
		"ALTER WAREHOUSE RESUME IF":                            false,
		"ALTER WAREHOUSE a SUSPEND -- a comment":               true,
		"ALTER WAREHOUSE a SUSPEND zzz":                        false,
		"ALTER WAREHOUSE SET TAG a = 's'":                      false,
		"":                                                     false,
		"ALTER WAREHOUSE a SET TAG a = $$it's$$, b = 's'":      true,
		"ALTER WAREHOUSE a SET TAG a = 1":                      false,
		"ALTER WAREHOUSE /* x */ a /* y */ SUSPEND":            true,
		"ALTER WAREHOUSE a b SUSPEND":                          false,
		"ALTER WAREHOUSE IF EXISTS IF EXISTS SUSPEND":          false,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED":      true,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED zzz":  false,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED,":     false,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED IF":   false,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED 'x'":  false,
		"ALTER WAREHOUSE IF EXISTS a RESUME IF SUSPENDED \x00": false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestRecognizerReadsTheGrammarsPunctuationWhole pins that => is one token
// where the grammar spells it, and = then > where it does not.
func TestRecognizerReadsTheGrammarsPunctuationWhole(t *testing.T) {
	g, err := ReadANTLR("grammar P; s : 'F' '(' 'a' '=>' 'b' ')' | 'G' 'a' '=' '>' 'b' ;", tokens(nil))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"F(a => b)":  true,
		"F(a = > b)": false,
		"G a = > b":  true,
		"G a => b":   false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestRecognizerHandlesNullableAndLeftRecursiveRules pins the Earley parts
// a recursive-descent reading would get wrong.
func TestRecognizerHandlesNullableAndLeftRecursiveRules(t *testing.T) {
	g, err := ReadANTLR(`grammar E; s : 'SELECT' e opt opt 'END' ; e : e '+' e | '1' | '(' e ')' ; opt : | 'X' ;`, tokens(nil))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT 1 END":           true,
		"SELECT 1+1+(1+1) X END": true,
		"SELECT 1 X X END":       true,
		"SELECT 1 X X X END":     false,
		"SELECT 1+ END":          false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestANTLRLiteralsAreTheirOwnText pins that a quoted literal is written as
// itself even where a token of the same name has another sample: 'STAGE'
// is the keyword, STAGE the token.
func TestANTLRLiteralsAreTheirOwnText(t *testing.T) {
	g, err := ReadANTLR("grammar S; s : 'CREATE' 'STAGE' STAGE ;", tokens(map[string]string{"STAGE": "@a"}))
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := g.Rows()
	if rows[0].Statement != "CREATE STAGE @a" {
		t.Errorf("got %q, want CREATE STAGE @a", rows[0].Statement)
	}
	if !g.Recognizer().Accepts("create stage @b") {
		t.Errorf("the recognizer does not read the keyword as itself")
	}
}

// TestANTLRGroupsThatReadAlikeStayApart pins that two groups whose labels
// read the same, a literal holding _ and two literals, are two rules.
func TestANTLRGroupsThatReadAlikeStayApart(t *testing.T) {
	g, err := ReadANTLR("grammar S; s : ('A_B')? 'X' ('A' 'B')? ;", tokens(nil))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{"X A B": true, "A_B X": true, "X A_B": false, "A B X": false} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
	rows, _ := g.Rows()
	labels := map[string]string{}
	for _, r := range rows {
		labels[r.Rule] = r.Statement
	}
	if labels["s(A_B)?/A_B"] != "A_B X" || labels["s(A_B)?#2/A_B"] != "X A B" {
		t.Errorf("rows:\n%s", render(rows))
	}
}

// TestRecognizerRefusesWhatTheLexerCannotRead pins three tokens no
// statement may match: a quote never closed, a token the lexer never
// produces, and a keyword a lexer rule spells in lowercase, which is that
// word and not a class of names.
func TestRecognizerRefusesWhatTheLexerCannotRead(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
options { caseInsensitive = true; }
s : SELECT NAME STR? | 'GET' NEVER | 'SHOW' QID ;
SELECT : 'select' ;
NAME : [a-z]+ ;
STR : '\'' (~'\'')* '\'' ;
QID : '"' (~'"')* '"' ;
`, tokens(map[string]string{"NEVER": Unreachable}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT a":        true,
		"select a 'it'":   true,
		"SELECT a 'oops":  false,
		"SELECT a $$oops": false,
		`SHOW "a"`:        true,
		`SHOW "a`:         false,
		"delete a":        false,
		`"other" a`:       false,
		"GET NEVER":       false,
		"GET":             false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestRecognizerKeepsClassesApart pins three readings: an identifier rule
// sampled in capitals is still a class of names, a stage @a and a column $1
// are different classes, and a comment never closed is not a comment.
func TestRecognizerKeepsClassesApart(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' ID | 'LIST' STAGE | 'SELECT' 'COL' POSITION | 'SELECT' '1' ;
ID : [A-Z_] [A-Z0-9_]* ;
`, tokens(map[string]string{"STAGE": "@a", "POSITION": "$1"}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT FOO":       true,
		"LIST @s":          true,
		"LIST $1":          false,
		"SELECT COL $2":    true,
		"SELECT COL @a":    false,
		"SELECT 1 /* x */": true,
		"SELECT 1 /*":      false,
		"SELECT 1 /*oops":  false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestRecognizerReadsNumbersAndQuotedPunctuation pins three lexer details:
// 2. is a number, 0x with no digit is not, and punctuation inside a quoted
// sample does not become a token of the grammar's.
func TestRecognizerReadsNumbersAndQuotedPunctuation(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' FLOAT | 'HEX' NUM | 'SELECT' '(' ')' | 'TEXT' STR ;
FLOAT : [0-9]+ '.' [0-9]* ;
`, tokens(map[string]string{"NUM": "1", "STR": "'()'"}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT 1.5": true,
		"SELECT 2.":  true,
		"HEX 0x1F":   true,
		"HEX 0x":     false,
		"SELECT()":   true,
		"SELECT ( )": true,
		"TEXT '(()'": true,
		"TEXT ( )":   false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestRecognizerMatchesEachSpellingOfAKeywordToken pins a lexer rule that
// chooses among keywords: it matches those words, and is no class of names.
func TestRecognizerMatchesEachSpellingOfAKeywordToken(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' DIRECTION ;
DIRECTION : 'ASC' | 'DESC' | D E S C E N D I N G ;
fragment D : [dD] ; fragment E : [eE] ; fragment S : [sS] ; fragment C : [cC] ;
fragment N : [nN] ; fragment I : [iI] ; fragment G : [gG] ;
`, tokens(nil))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT ASC":        true,
		"select desc":       true,
		"SELECT DESCENDING": true,
		"SELECT zzz":        false,
		`SELECT "ASC"`:      false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestReadANTLRReservesTheWordsItsLexerSpells pins what a vendor's lexer
// rules say of keywords. ANTLR gives a text two rules match to the one
// written first, so a word a token spells before the identifier rule is
// never a name, except where a parser rule lets the token stand for one, as
// name does for FIRST. A grammar with no lexer rules, as one converted from
// docs is, reserves nothing this way.
func TestReadANTLRReservesTheWordsItsLexerSpells(t *testing.T) {
	const parser = `grammar S;
s : K_SELECT name (K_FROM name)? ('LIMIT' NUMBER)? ;
name : ID | K_FIRST ;
`
	const lexer = `K_SELECT : S E L E C T ;
K_FROM : F R O M ;
K_FIRST : 'FIRST' ;
ORDER : 'ASC' | 'DESC' ;
ID : [a-zA-Z_] [a-zA-Z_0-9]* ;
NUMBER : [0-9]+ ;
LATE : 'late' ;
fragment S : [sS] ; fragment E : [eE] ; fragment L : [lL] ; fragment C : [cC] ; fragment T : [tT] ;
fragment F : ('f'|'F') ; fragment R : 'r' | 'R' ; fragment O : ('O' | 'o') ; fragment M : 'm' ;
`
	g, err := ReadANTLR("\ufeff"+parser+lexer, tokens(map[string]string{"ID": "a", "NUMBER": "1"}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT a":          true,
		"SELECT a FROM a":   true,
		"SELECT first":      true,  // the parser rule admits the keyword
		"SELECT from":       false, // a keyword spelled ('f'|'F') ('r'|'R') ..., as ANTLR 3 grammars do
		"SELECT a from a":   true,
		"SELECT a frox a":   false, // and it is no class of names
		"SELECT select":     false,
		"SELECT a FROM asc": false, // a word of a token of several spellings
		"SELECT limit":      false, // a literal of a parser rule
		"SELECT a LIMIT 1":  true,
		"SELECT late":       true, // spelled after the last class token: ID takes it
		"SELECT s":          true, // a fragment is no token
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}

	docs, err := ReadANTLR(`grammar D;
s : 'SELECT' ID ('FROM' ID)? ;
`, tokens(map[string]string{"ID": "a"}))
	if err != nil {
		t.Fatal(err)
	}
	if !docs.Recognizer().Accepts("SELECT from") {
		t.Error("a grammar with no lexer rules reserved a word: only its token table can")
	}
}

// TestRecognizerTellsQuotedNamesByTheirQuotes pins the classes of names: a
// token sampled by a bare name takes a name however it is quoted, and one
// sampled in quotes takes only names in those quotes, as PartiQL's lexer
// reads "a" as an identifier and `a` as an Ion value.
func TestRecognizerTellsQuotedNamesByTheirQuotes(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' ( ID | QUOTED | ION 'ION' ) ;
`, tokens(map[string]string{"ID": "a", "QUOTED": `"a"`, "ION": "`1`"}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT b":           true,
		`SELECT "b"`:         true,
		"SELECT `b`":         true, // ID is sampled bare: a name in any quotes
		`SELECT "b" ION`:     false,
		"SELECT `{a:1}` ION": true,
		"SELECT b ION":       false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestReadANTLRSamplesARuleFromTheTokenTable pins a token table row that
// names a rule of a .g4, as one may name a rule of a bison grammar: other
// rules' rows write its text where it stands, and its own alternatives keep
// their rows.
func TestReadANTLRSamplesARuleFromTheTokenTable(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' term ('+' term)? ;
term : 'CURRENT_USER' | ID ;
`, tokens(map[string]string{"ID": "a", "term": "a"}))
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := g.Rows()
	got := map[string]string{}
	for _, r := range rows {
		got[r.Rule] = r.Statement
	}
	for rule, want := range map[string]string{
		"s(+_term)?/+_term": "SELECT a + a",
		"term/CURRENT_USER": "SELECT CURRENT_USER",
		"term/ID":           "SELECT a",
	} {
		if got[rule] != want {
			t.Errorf("%s = %q, want %q", rule, got[rule], want)
		}
	}
}

// TestReadANTLRReservesOnlyWhatNoEarlierClassTakes pins precedence against
// the class token that matches the word, not the last class token: a
// keyword written after an identifier rule that matches it is never
// produced, so its word stays a name, whatever class rules follow.
func TestReadANTLRReservesOnlyWhatNoEarlierClassTakes(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' ID ;
ID : [a-zA-Z_][a-zA-Z_0-9]* ;
LATE : 'late' ;
EARLY_NUMBER : [0-9]+ ;
FIRST : 'first' ;
WS : [ \t\r\n]+ -> skip ;
`, tokens(map[string]string{}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT late":  true, // ID is written first and matches it
		"SELECT first": true,
		"SELECT a":     true,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}

	// A class written first that does not match the word leaves it reserved,
	// and a caseInsensitive grammar's [A-Z] matches either case.
	g, err = ReadANTLR(`grammar S;
options { caseInsensitive = true; }
s : 'SELECT' ID ;
NUMBER : [0-9]+ ;
FROM : 'FROM' ;
ID : [A-Z_][A-Z_0-9]* ;
LATE : 'late' ;
`, tokens(map[string]string{}))
	if err != nil {
		t.Fatal(err)
	}
	rec = g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT from": false, // NUMBER does not match it: FROM takes it
		"SELECT late": true,  // ID matches it in either case
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("caseInsensitive: Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestReadANTLRKeepsTheCaseOfStringSpellings pins what case folding
// touches: a keyword spelled in either case is one word, but a token whose
// spellings are strings that differ in case keeps each, since the
// recognizer compares a string exactly.
func TestReadANTLRKeepsTheCaseOfStringSpellings(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' VALUE ;
VALUE : '\'a\'' | '\'A\'' ;
`, tokens(map[string]string{}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT 'a'": true,
		"SELECT 'A'": true,
		"SELECT 'b'": false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}

	// A keyword spelled letter by letter, longer than maxSpellings would
	// allow were each case a text of its own, is still one keyword.
	g, err = ReadANTLR(`grammar S;
s : 'SELECT' K_KEYSPACES ID ;
K_KEYSPACES : K E Y S P A C E S ;
ID : [a-z]+ ;
fragment K : ('k'|'K') ; fragment E : ('e'|'E') ; fragment Y : ('y'|'Y') ; fragment S : ('s'|'S') ;
fragment P : ('p'|'P') ; fragment A : ('a'|'A') ; fragment C : ('c'|'C') ;
`, tokens(map[string]string{"ID": "a"}))
	if err != nil {
		t.Fatal(err)
	}
	rec = g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT keyspaces a": true,
		"SELECT KeySpaces a": true,
		"SELECT keyspace a":  false,
		"SELECT a keyspaces": false, // reserved: no name
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("long keyword: Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestTokenTablePicksASpelling pins a token table row for a token of a few
// spellings: rows write the spelling it names, and the recognizer still
// matches every spelling the lexer rule has.
func TestTokenTablePicksASpelling(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'DROP' K_TABLE ID ;
K_TABLE : 'COLUMNFAMILY' | 'TABLE' ;
ID : [a-z]+ ;
`, tokens(map[string]string{"ID": "a", "K_TABLE": "TABLE"}))
	if err != nil {
		t.Fatal(err)
	}
	rows, _ := g.Rows()
	if rows[0].Statement != "DROP TABLE a" {
		t.Errorf("first row %q, want DROP TABLE a", rows[0].Statement)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"DROP TABLE a":        true,
		"DROP columnfamily a": true,
		"DROP VIEW a":         false,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestParserLiteralTakesItsLexerRulesPrecedence pins a parser literal that
// a lexer rule spells alone: it is that rule's token, so an identifier rule
// written before it takes the word, while a literal no rule spells is an
// implicit token defined first.
func TestParserLiteralTakesItsLexerRulesPrecedence(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' ID | 'late' | 'early' ;
ID : [a-z]+ ;
LATE : 'late' ;
`, tokens(map[string]string{}))
	if err != nil {
		t.Fatal(err)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{
		"SELECT late":  true, // ID is written before LATE, which 'late' names
		"late":         true,
		"SELECT early": false, // 'early' is an implicit token, ahead of ID
		"early":        true,
	} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}

// TestCaseIsKeptInsideAQuotedToken pins folding against a string token
// built from pieces: '\” ('a' | 'A') '\” is two texts, not one.
func TestCaseIsKeptInsideAQuotedToken(t *testing.T) {
	for _, value := range []string{`'\'' ('a' | 'A') '\''`, `'\'' [aA] '\''`, `'\'' LETTER '\'' ; fragment LETTER : 'a' | 'A'`} {
		g, err := ReadANTLR("grammar S;\ns : 'SELECT' VALUE ;\nVALUE : "+value+" ;\n", tokens(map[string]string{}))
		if err != nil {
			t.Fatal(err)
		}
		rec := g.Recognizer()
		for sql, want := range map[string]bool{"SELECT 'a'": true, "SELECT 'A'": true, "SELECT 'b'": false} {
			if got := rec.Accepts(sql); got != want {
				t.Errorf("VALUE : %s: Accepts(%q) = %v, want %v", value, sql, got, want)
			}
		}
	}
}

// TestCaseIsFoldedOutsideQuotes pins a keyword with punctuation, as
// Cassandra's '-' N A N: it is one text in capitals, matched in any case.
func TestCaseIsFoldedOutsideQuotes(t *testing.T) {
	g, err := ReadANTLR(`grammar S;
s : 'SELECT' K_NEGATIVE_NAN ;
K_NEGATIVE_NAN : '-' N A N ;
fragment N : ('n'|'N') ; fragment A : ('a'|'A') ;
`, tokens(map[string]string{}))
	if err != nil {
		t.Fatal(err)
	}
	if rows, _ := g.Rows(); rows[0].Statement != "SELECT -NAN" {
		t.Errorf("first row %q, want SELECT -NAN", rows[0].Statement)
	}
	rec := g.Recognizer()
	for sql, want := range map[string]bool{"SELECT -nan": true, "SELECT -NaN": true, "SELECT -na": false} {
		if got := rec.Accepts(sql); got != want {
			t.Errorf("Accepts(%q) = %v, want %v", sql, got, want)
		}
	}
}
