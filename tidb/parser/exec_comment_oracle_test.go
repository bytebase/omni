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
		tidbAccepts := !tidbRejectedSyntax(tc.canExecute(sql))
		_, err := Parse(sql)
		if (err == nil) != tidbAccepts {
			t.Errorf("%q: omni accepts=%v (err=%v), TiDB accepts=%v", sql, err == nil, err, tidbAccepts)
		}
	}
}
