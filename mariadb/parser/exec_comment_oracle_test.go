package parser

import "testing"

// TestExecCommentBoundary_Oracle proves against a live MariaDB 11.8 container
// how the syntax of an executable comment bounds tokens. The lexer reads the
// body in place, so "/*!" and "*/" end the token before them the way
// whitespace would.
func TestExecCommentBoundary_Oracle(t *testing.T) {
	o := startMariaDB(t)

	for _, sql := range []string{
		// "*/" ends the number and the identifier before it: SELECT 1 2 and
		// SELECT 'v' x y, both rejected.
		"/*!SELECT 1*/2",
		"/*!SELECT 'v' x*/y",
		"SELECT /*!50000 1, */ 2",
		// A token that swallows the closing "*/" leaves the comment open, which
		// MariaDB rejects. TiDB accepts both of these, so its fork of the lexer
		// does not treat the shape as unterminated.
		`SELECT 'x' /*! , 'a*/'`,
		"SELECT 1 /*!50000 + 1 -- c */\n",
	} {
		v, _, msg := o.classify(sql)
		if v == vErrored {
			t.Fatalf("server: %q: %s", sql, msg)
		}
		_, err := Parse(sql)
		if (err == nil) != (v == vAccepted) {
			t.Errorf("%q: omni accepts=%v (err=%v), MariaDB accepts=%v", sql, err == nil, err, v == vAccepted)
		}
	}
}

// TestExecCommentQuotedDelimiters_Oracle proves against a live MariaDB 11.8
// container that omni ends an executable comment where the server does when a
// "/*" or "*/" sits inside a string, a quoted identifier or a line comment.
func TestExecCommentQuotedDelimiters_Oracle(t *testing.T) {
	o := startMariaDB(t)

	for _, list := range [][]string{execCommentQuotedAccepted, execCommentQuotedHeldClose} {
		for _, sql := range list {
			v, _, msg := o.classify(sql)
			if v == vErrored {
				t.Fatalf("server: %q: %s", sql, msg)
			}
			_, err := Parse(sql)
			if (err == nil) != (v == vAccepted) {
				t.Errorf("%q: omni accepts=%v (err=%v), MariaDB accepts=%v", sql, err == nil, err, v == vAccepted)
			}
		}
	}
}
