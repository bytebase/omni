package parser

// execOpener reports whether s[i:], which starts with "/*", opens an
// executable comment, and the offset its body starts at: past "/*!" and past
// the five digits after it, when five digits follow.
//
// OceanBase's lexer has two openers, `\/\*\![0-9]{5}` and `\/\*\!`, and takes
// the longer one that matches, so a version is exactly five digits: a sixth
// digit, or a run shorter than five, is the first token of the body. It never
// compares the version with its own, so the body is read whatever the version
// says.
//
// Ref: OceanBase CE v4.3.5_CE_BP6 src/sql/parser/sql_parser_mysql_mode.l:95-96
// (mysql_compatible_comment_with_version, _without_version), :1412 and :1844.
//
// OceanBase 4.3.5.6 (reports 5.7.25) returns 2 for "SELECT 1 /*!99999 + 1 */",
// 1234 for "SELECT /*!1234 */", 6 for "SELECT /*!123456 */" and 10 for
// "SELECT 2 * /*!100005 */", measured on 2026-11-05 (execVersionCases).
func execOpener(s string, i int) (int, bool) {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' || s[i+2] != '!' {
		return 0, false
	}
	j := i + 3
	for k := j; k < j+5; k++ {
		if k >= len(s) || s[k] < '0' || s[k] > '9' {
			return j, true
		}
	}
	return j + 5, true
}
