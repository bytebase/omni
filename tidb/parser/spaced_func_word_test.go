package parser

import "testing"

// ============================================================================
// A space before a builtin's '(' (see checkSpacedFuncWord in temporal.go).
//
// The aggregate and datetime builtins reach the grammar as a per-name token the
// lexer only folds out of `name(` when the '(' sits directly after the name
// (btFuncTokenMap, misc.go isTokenIdentifier). With a space the word is no longer
// that token, so SumExpr's special forms — '*', the DISTINCT/ALL quantifiers,
// GROUP_CONCAT's SEPARATOR, SUBSTRING/TRIM's FROM, a precision — are unreachable.
// A word that is also a keyword (SUM, MAX, GROUP_CONCAT, SUBSTRING, ...) then
// cannot name a generic call either, so spaced it is a 1064 whatever follows; a
// word that is only the builtin token (COUNT, SYSDATE) reverts to a plain
// identifier and reads as an ordinary generic call. AVG and JSON_ARRAYAGG/
// JSON_OBJECTAGG are keyword literals in SumExpr, not btFuncTokenMap tokens, so a
// space changes nothing for them.
//
// Every row is held to a live TiDB v8.5.5 by TestTiDBSpacedFuncWordsMatchEngine.
// ============================================================================

var spacedFuncWordCases = []struct {
	sql    string
	accept bool
}{
	// --- COUNT: only the builtin token, so spaced it is a generic call ---
	{"select count(*)", true},
	{"select count ( * )", false},
	{"select count(distinct a)", true},
	{"select count ( distinct a )", false},
	{"select count ( all a )", false},
	{"select count(a)", true},
	{"select count ( a )", true},
	{"select count(a,a)", false},
	{"select count ( a, a )", true},
	{"select count ( 1 )", true},

	// --- SUM/MIN/MAX: also keywords, so spaced is a syntax error whatever follows ---
	{"select max(a)", true},
	{"select max ( a )", false},
	{"select max('bbb')", true},
	{"select max ( 'bbb' )", false},
	{"select max(distinct a)", true},
	{"select max ( distinct a )", false},
	{"select sum(a)", true},
	{"select sum ( a )", false},
	{"select min(a)", true},
	{"select min ( a )", false},

	// --- GROUP_CONCAT ---
	{"select group_concat(a)", true},
	{"select group_concat ( a )", false},
	{"select group_concat(a separator ',')", true},
	{"select group_concat ( a separator ',' )", false},

	// --- STD / STDDEV* ---
	{"select std(a)", true},
	{"select std ( a )", false},
	{"select stddev ( a )", false},
	{"select stddev_pop ( a )", false},
	{"select stddev_samp(a)", true},
	{"select stddev_samp ( a )", false},

	// --- VAR_POP / VAR_SAMP / VARIANCE ---
	{"select var_pop(a)", true},
	{"select var_pop ( a )", false},
	{"select var_samp ( a )", false},
	{"select variance ( a )", false},

	// --- BIT_AND / BIT_OR / BIT_XOR ---
	{"select bit_and(a)", true},
	{"select bit_and ( a )", false},
	{"select bit_or ( a )", false},
	{"select bit_xor ( a )", false},

	// --- APPROX_COUNT_DISTINCT / APPROX_PERCENTILE ---
	{"select approx_count_distinct(a)", true},
	{"select approx_count_distinct ( a )", false},
	{"select approx_percentile(a, 50)", true},
	{"select approx_percentile ( a, 50 )", false},

	// --- the scalar builtins the same rule gates ---
	{"select substr(a, 1)", true},
	{"select substr ( a, 1 )", false},
	{"select substring(a from 1)", true},
	{"select substring ( a from 1 )", false},
	{"select trim(a)", true},
	{"select trim ( a )", false},
	{"select trim(leading 'x' from a)", true},
	{"select trim ( leading 'x' from a )", false},
	{"select position('a' in b)", true},
	{"select position ( 'a' in b )", false},

	// --- AVG and JSON_*: keyword literals in SumExpr, so a space is harmless ---
	{"select avg(a)", true},
	{"select avg ( a )", true},
	{"select avg(distinct a)", true},
	{"select avg ( distinct a )", true},
	{"select avg(*)", false},
	{"select avg ( * )", false},
	{"select json_arrayagg(a)", true},
	{"select json_arrayagg (a)", true},
	{"select json_objectagg(a, a)", true},
	{"select json_objectagg (a, a)", true},

	// --- SYSDATE: only the builtin token, so spaced it is a generic call; a
	//     precision argument then becomes an ordinary one (sysdate(b) wants a
	//     precision, sysdate (b) does not) ---
	{"select sysdate(b)", false},
	{"select sysdate ( b )", true},

	// --- temporal keywords still reject spaced, unaffected by this change ---
	{"select now()", true},
	{"select now ( )", false},
	{"select curdate(1)", false},
	{"select curdate ( 1 )", false},

	// --- a non-keyword name is a generic call either way; only the illegal '*'
	//     argument is refused, which is what gates spaced COUNT above ---
	{"select abc(a)", true},
	{"select abc ( a )", true},
	{"select abc ( * )", false},
	{"select foo (x)", true},
}

func TestTiDBSpacedFuncWords(t *testing.T) {
	for _, c := range spacedFuncWordCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.accept {
			t.Errorf("Parse(%q) err = %v, want accept = %v", c.sql, err, c.accept)
		}
	}
}
