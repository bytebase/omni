package parser

import "testing"

// Function calls in expression position (reached here through the DO statement, parseDoStmt,
// utility.go), as TiDB v8.5.5 reads them. Two over-acceptance families the §3.6 negatives
// corpus surfaced, both in the shared function-call path (expr.go):
//
//   - A keyword TOKEN is not the grammar's `identifier`, so it cannot be a generic function
//     name. QUERIES, REPLAYER, TSO, DEFINED and LAST_BACKUP are keyword tokens on TiDB (all in
//     tidbLexerWords) with no function-call production, so `DO queries()` is a 1064 where the
//     plain `DO foo()` parses. This parser's forked lexer scans them as identifiers, so
//     tidbKeywordTokensNotFunctionNames carries them for refusesFunctionName.
//   - APPROX_COUNT_DISTINCT and APPROX_PERCENTILE (SumExpr, parser.y:8470,8474) take a non-empty
//     ExpressionList, and LASTVAL (FunctionNameSequence, parser.y:8438) exactly one TableName, so
//     a zero-argument call is a 1064 that the generic variable-arity path would take.
//
// The valid forms are pinned beside the rejections so a widened or narrowed grammar is reported
// without a server. The reject verdicts are the negatives corpus's measured engine column
// (compat.tsv); the accept forms are measured in the acceptance corpora.

// TestDoFuncNameArityForms holds omni to the recorded verdicts without a server.
func TestDoFuncNameArityForms(t *testing.T) {
	for _, c := range doFuncNameArityCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

var doFuncNameArityCases = []tableSampleCase{
	// A keyword token is not a generic function name.
	{"DO QUERIES()", false},
	{"DO REPLAYER()", false},
	{"DO TSO()", false},
	{"DO DEFINED(a)", false},
	{"DO LAST_BACKUP(a)", false},
	{"DO QUERIES(a)", false},
	{"DO TSO(a, b)", false},
	// A quoted name is an ordinary identifier, so the generic call stands, and a qualified name
	// is a stored-function reference.
	{"DO `queries`()", true},
	{"DO `tso`(a, b)", true},
	{"DO `last_backup`(a)", true},
	{"DO d.queries()", true},
	{"DO d.tso(a, b)", true},
	// A plain identifier still names a generic call of any arity.
	{"DO foo()", true},
	{"DO myfunc(a, b, c)", true},
	// Keyword functions kept through their dedicated productions must still parse.
	{"DO IF(a, b, c)", true},
	{"DO REPLACE(a, b, c)", true},
	{"DO LEFT(a, 1)", true},
	{"DO ASCII(a)", true},
	{"DO COALESCE(a, b)", true},

	// APPROX_COUNT_DISTINCT / APPROX_PERCENTILE take a non-empty argument list.
	{"DO APPROX_COUNT_DISTINCT()", false},
	{"DO APPROX_PERCENTILE()", false},
	{"DO APPROX_COUNT_DISTINCT(a)", true},
	{"DO APPROX_COUNT_DISTINCT(a, b)", true},
	{"DO APPROX_PERCENTILE(a, 50)", true},
	{"select approx_count_distinct(c) from t", true},
	// A quoted name is an ordinary call with no such floor.
	{"DO `approx_count_distinct`()", true},

	// LASTVAL takes exactly one TableName.
	{"DO LASTVAL()", false},
	{"select lastval()", false},
	{"select lastval(a, b)", false},
	{"select lastval(1)", false},
	{"select lastval(a)", true},
	{"select lastval(a.b)", true},
	{"select lastval(test.seq)", true},
	{"select `lastval`()", true},
	{"select `lastval`(a, b)", true},
}
