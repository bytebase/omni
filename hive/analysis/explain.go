package analysis

import "strings"

// ExplainStatement returns the Hive EXPLAIN statement for statement:
// EXPLAIN FORMATTED when format is "JSON" (case-insensitive), which Hive
// answers as one JSON object, and EXPLAIN otherwise. A statement that already
// starts with EXPLAIN is returned unchanged.
func ExplainStatement(statement, format string) string {
	if _, ok := PlanFormat(statement); ok {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN FORMATTED " + statement
	}
	return "EXPLAIN " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading whitespace and comments.
func StartsWithExplain(statement string) bool {
	_, ok := PlanFormat(statement)
	return ok
}

// PlanFormat reports the format ("JSON" or "TEXT") of a statement that is
// already an EXPLAIN: EXPLAIN FORMATTED answers with JSON, and every other
// EXPLAIN with a text tree.
func PlanFormat(statement string) (string, bool) {
	s, ok := pastLeadingNoise(statement)
	if !ok || !hasKeyword(s, "EXPLAIN") {
		return "", false
	}
	if rest, ok := pastLeadingNoise(s[len("EXPLAIN"):]); ok && hasKeyword(rest, "FORMATTED") {
		return "JSON", true
	}
	return "TEXT", true
}

// pastLeadingNoise returns statement after its leading whitespace and
// comments, and reports whether it could be read at all. A block comment left
// open or holding another `/*` cannot be settled without a parser.
func pastLeadingNoise(statement string) (string, bool) {
	s := statement
	for {
		s = strings.TrimLeft(s, " \t\r\n\v\f")
		switch {
		case strings.HasPrefix(s, "--"):
			_, rest, found := strings.Cut(s, "\n")
			if !found {
				return "", false
			}
			s = rest
		case strings.HasPrefix(s, "/*"):
			rest, closed := pastBlockComment(s)
			if !closed {
				return "", false
			}
			s = rest
		default:
			return s, true
		}
	}
}

func pastBlockComment(s string) (string, bool) {
	for at := len("/*"); at < len(s); at++ {
		switch {
		case strings.HasPrefix(s[at:], "/*"):
			return "", false
		case strings.HasPrefix(s[at:], "*/"):
			return s[at+len("*/"):], true
		}
	}
	return "", false
}

func hasKeyword(s, keyword string) bool {
	if len(s) < len(keyword) || !strings.EqualFold(s[:len(keyword)], keyword) {
		return false
	}
	if len(s) == len(keyword) {
		return true
	}
	c := s[len(keyword)]
	return c != '_' && c != '$' && (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && c < 0x80
}
