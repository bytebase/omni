package parser

// execOpener reports whether s[i:], which starts with "/*", opens an
// executable comment, "/*!" and any version digits after it, and the offset
// its body starts at. A comment MariaDB skips (skippedComment) opens none.
func execOpener(s string, i int) (int, bool) {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' || s[i+2] != '!' {
		return 0, false
	}
	if skippedComment(s, i) {
		return 0, false
	}
	j := i + 3
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	return j, true
}

// skippedComment reports whether s[i:] opens a versioned comment MariaDB skips
// unread: "/*!" with a version from 50700 to 99999. MariaDB reads five digits,
// or six when a sixth follows, and takes those versions for MySQL 5.7 and 8.0
// releases ("Ignore all versioned comments with MySQL versions in the range
// 50700-999999", sql/sql_lex.cc, mariadb-11.8.8), so MariaDB 11.8.8 returns 4
// for "SELECT 1 /*!80023 + 1 */ + 3" and 5 for "SELECT 1 /*!50699 + 1 */ + 3".
func skippedComment(s string, i int) bool {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' || s[i+2] != '!' {
		return false
	}
	j := i + 3
	v, n := 0, 0
	for j < len(s) && n < 6 && s[j] >= '0' && s[j] <= '9' {
		v = v*10 + int(s[j]-'0')
		j++
		n++
	}
	return n >= 5 && v >= 50700 && v <= 99999
}

// skippedCommentEnd returns the offset after the "*/" that ends the skipped
// comment opening at s[i:], or -1 when the input ends first. MariaDB counts
// nesting here and reads raw bytes, not tokens: every "/*" opens a level and
// every "*/" closes one, a quoted one included. MariaDB 11.8.8 returns 4 for
// "SELECT 1 /*!80023 + 1 /* c */ + 2 */ + 3" and for
// "SELECT 1 /*!80023 'x */ + 3".
func skippedCommentEnd(s string, i int) int {
	depth := 0
	for i+1 < len(s) {
		switch {
		case s[i] == '/' && s[i+1] == '*':
			depth++
			i += 2
		case s[i] == '*' && s[i+1] == '/':
			depth--
			i += 2
			if depth == 0 {
				return i
			}
		default:
			i++
		}
	}
	return -1
}
