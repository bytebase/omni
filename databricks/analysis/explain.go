package analysis

import "strings"

// ExplainStatement returns the Databricks EXPLAIN statement for statement.
// Databricks prints Spark's FORMATTED plan regardless of the requested format;
// a statement that already asks for a plan is returned unchanged.
func ExplainStatement(statement, _ string) string {
	if StartsWithExplain(statement) {
		return statement
	}
	return "EXPLAIN FORMATTED " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading whitespace and comments. Because Spark nests block comments
// and this package has no parser, a block comment left open or containing
// another `/*` returns false so a write hidden behind a nested comment is
// never mistaken for an EXPLAIN.
func StartsWithExplain(statement string) bool {
	s, ok := pastLeadingNoise(statement)
	return ok && hasKeyword(s, "EXPLAIN")
}

// PlanExecuted reports whether statement's EXPLAIN is the ANALYZE kind, whose
// output is a plan of an execution that happened.
func PlanExecuted(statement string) bool {
	s, ok := pastLeadingNoise(statement)
	if !ok || !hasKeyword(s, "EXPLAIN") {
		return false
	}
	rest, ok := pastLeadingNoise(s[len("EXPLAIN"):])
	return ok && hasKeyword(rest, "ANALYZE")
}

// pastLeadingNoise returns statement after its leading whitespace and
// comments, and reports whether it could be read at all.
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
