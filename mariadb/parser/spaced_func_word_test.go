package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/deparse"
)

// ============================================================================
// A space before an aggregate or scalar builtin's '(' (see spacedFuncWords in
// temporal.go).
//
// These builtins reach their special argument syntax — COUNT(*), the
// DISTINCT/ALL quantifier, GROUP_CONCAT's SEPARATOR, SUBSTRING/TRIM's FROM,
// POSITION's IN — through a token the lexer folds out of `name(` only when the
// '(' sits directly after the name. With a space the word is a plain identifier
// and the call is an ordinary generic one: it takes expression arguments and no
// '*', quantifier or keyword clause. So spaced the special forms are a 1064, and
// the completeness gap closes the other way too — count (a, a) and max ( ) run as
// generic calls where the aggregate productions reject them.
//
// AVG, WEIGHT_STRING, GET_FORMAT, TIMESTAMPDIFF and TIMESTAMPADD are
// reserved-keyword productions, not folded tokens, so a space changes nothing
// for them.
//
// Every row is held to a live mariadb:11.8.8 by
// TestMariaDBSpacedFuncWordsMatchEngine.
// ============================================================================

var spacedFuncWordCases = []struct {
	sql    string
	accept bool
}{
	// --- COUNT: spaced it is a generic call ---
	{"select count(*)", true},
	{"select count ( * )", false},
	{"select count(distinct a)", true},
	{"select count ( distinct a )", false},
	{"select count ( all a )", false},
	{"select count(a)", true},
	{"select count ( a )", true},
	{"select count(a,a)", false},
	{"select count ( a, a )", true},
	{"select count ( )", true},

	// --- SUM/MIN/MAX: spaced they are generic calls, so the quantifier goes but
	//     a second argument and an empty list are now allowed ---
	{"select max(a)", true},
	{"select max ( a )", true},
	{"select max(distinct a)", true},
	{"select max ( distinct a )", false},
	{"select max(a,a)", false},
	{"select max ( a, a )", true},
	{"select max ( )", true},
	{"select sum(a)", true},
	{"select sum ( a )", true},
	{"select sum ( * )", false},
	{"select sum ( distinct a )", false},
	{"select min(a)", true},
	{"select min ( a )", true},

	// --- GROUP_CONCAT ---
	{"select group_concat(a)", true},
	{"select group_concat ( a )", true},
	{"select group_concat(distinct a)", true},
	{"select group_concat ( distinct a )", false},
	{"select group_concat(a separator ',')", true},
	{"select group_concat ( a separator ',' )", false},
	{"select group_concat ( a, a )", true},

	// --- STD / STDDEV* / VAR* / BIT_* ---
	{"select std(a)", true},
	{"select std ( a )", true},
	{"select std ( * )", false},
	{"select stddev ( a )", true},
	{"select stddev_pop ( a )", true},
	{"select stddev_samp ( a )", true},
	{"select var_pop ( a )", true},
	{"select var_samp ( a )", true},
	{"select variance ( a )", true},
	{"select bit_and ( a )", true},
	{"select bit_or ( a )", true},
	{"select bit_xor ( a )", true},
	{"select bit_and ( * )", false},

	// --- JSON aggregates: spaced they are generic, so arity opens up ---
	{"select json_arrayagg(a)", true},
	{"select json_arrayagg ( a )", true},
	{"select json_arrayagg(a,a)", false},
	{"select json_arrayagg ( a, a )", true},
	{"select json_objectagg(a, b)", true},
	{"select json_objectagg ( a, b )", true},
	{"select json_objectagg(a)", false},
	{"select json_objectagg ( a )", true},

	// --- SUBSTR / SUBSTRING ---
	{"select substr(a, 1)", true},
	{"select substr ( a, 1 )", true},
	{"select substr(a from 1)", true},
	{"select substr ( a from 1 )", false},
	{"select substring(a from 1)", true},
	{"select substring ( a from 1 )", false},
	{"select substring ( a )", true},

	// --- TRIM ---
	{"select trim(a)", true},
	{"select trim ( a )", true},
	{"select trim(leading 'x' from a)", true},
	{"select trim ( leading 'x' from a )", false},

	// --- POSITION ---
	{"select position('a' in b)", true},
	{"select position ( 'a' in b )", false},
	{"select position(a)", false},
	{"select position ( a )", true},

	// --- AVG: a reserved-keyword production, so a space is harmless ---
	{"select avg(a)", true},
	{"select avg ( a )", true},
	{"select avg(distinct a)", true},
	{"select avg ( distinct a )", true},
	{"select avg(*)", false},
	{"select avg ( * )", false},

	// --- a non-keyword name is a generic call either way; only the illegal '*'
	//     argument is refused, which is what gates the spaced builtins above ---
	{"select abc(a)", true},
	{"select abc ( a )", true},
	{"select abc ( * )", false},
	{"select foo (x)", true},
}

func TestMariaDBSpacedFuncWords(t *testing.T) {
	for _, c := range spacedFuncWordCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.accept {
			t.Errorf("Parse(%q) err = %v, want accept = %v", c.sql, err, c.accept)
		}
	}
}

// TestMariaDBSpacedFuncWordsRoundTrip pins that a spaced generic call whose name
// is a gated builtin keyword survives deparse and parses again. Without the
// SpacedParen marker the renderer would drop the space and the output would
// re-lex as the builtin: max(a, a), count(a, a), count() and the special-form
// arities are all 1064s. The parser sets SpacedParen on exactly these calls.
func TestMariaDBSpacedFuncWordsRoundTrip(t *testing.T) {
	for _, sql := range []string{
		"select count (a)",
		"select count (a, a)",
		"select count ()",
		"select max (a, a)",
		"select max ()",
		"select sum (a)",
		"select group_concat (a, a)",
		"select position (a)",
		"select substring (a)",
		"select json_objectagg (a)",
	} {
		fc := firstFuncCall(t, sql)
		if !fc.SpacedParen {
			t.Errorf("%q: want SpacedParen set on the FuncCallExpr", sql)
		}
		out := deparse.Deparse(fc)
		if _, err := Parse("select " + out); err != nil {
			t.Errorf("%q: deparse produced %q, which does not re-parse: %v", sql, out, err)
		}
	}
}
