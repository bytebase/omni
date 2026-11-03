package parser

import "testing"

// TestExecCommentBoundary_Oracle proves against a live TiDB v8.5.5 container
// how the syntax of an executable comment bounds tokens. The lexer reads the
// body in place, so "/*!", "/*T!", and "*/" end the token before them the way
// whitespace would.
func TestExecCommentBoundary_Oracle(t *testing.T) {
	tc := startTiDB(t)

	for _, sql := range []string{
		// "*/" ends the number and the identifier before it: SELECT 1 2 and
		// SELECT 'v' x y, both rejected.
		"/*!SELECT 1*/2",
		"/*!SELECT 'v' x*/y",
		"/*T!SELECT 1*/2",
		"SELECT /*!50000 1, */ 2",
		"SELECT /*T! 1, */ 2",
		// TiDB accepts a string or a line comment that holds the closing "*/".
		"SELECT 'x' /*! , 'a*/'",
		"SELECT 1 /*!50000 + 1 -- c */\n",
	} {
		tidbAccepts := !tidbRejectedSyntax(t, tc.canExecute(sql))
		_, err := Parse(sql)
		if (err == nil) != tidbAccepts {
			t.Errorf("%q: omni accepts=%v (err=%v), TiDB accepts=%v", sql, err == nil, err, tidbAccepts)
		}
	}
}

// TestExecCommentQuotedDelimiters_Oracle proves against a live TiDB v8.5.5
// container that omni ends an executable comment where the server does when a
// "/*" or "*/" sits inside a string, a quoted identifier or a line comment,
// and that TiDB accepts a comment whose close is held inside a token.
func TestExecCommentQuotedDelimiters_Oracle(t *testing.T) {
	tc := startTiDB(t)

	for _, list := range [][]string{execCommentQuotedAccepted, execCommentQuotedHeldClose} {
		for _, sql := range list {
			tidbAccepts := !tidbRejectedSyntax(t, tc.canExecute(sql))
			_, err := Parse(sql)
			if (err == nil) != tidbAccepts {
				t.Errorf("%q: omni accepts=%v (err=%v), TiDB accepts=%v", sql, err == nil, err, tidbAccepts)
			}
		}
	}
}
