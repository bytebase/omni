package driver

import "strings"

// explainStatement asks Hive for a plan of statement: FORMATTED when the
// caller wants a document, which Hive answers as one JSON object whose CBOPlan
// is the plan Calcite built, and the plan tree of lines otherwise.
//
// A statement that already asks for a plan runs as the caller wrote it, so a
// typed EXPLAIN is never planned twice. With no parser to ask, that is decided
// by reading the text, the way Bytebase's standard explain rewrite read it.
func explainStatement(statement, format string) string {
	if _, ok := planFormat(statement); ok {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN FORMATTED " + statement
	}
	return "EXPLAIN " + statement
}

// planFormat reports the format of a statement that is already a plan:
// EXPLAIN FORMATTED answers with one JSON document, and every other EXPLAIN
// with a tree of lines. A plan does not run what it plans, ANALYZE included —
// the pinned HiveServer2 answers EXPLAIN ANALYZE INSERT with a plan and inserts
// nothing — so the plan Result carries no Executed.
func planFormat(statement string) (string, bool) {
	s, ok := pastLeadingNoise(statement)
	if !ok || !hasKeyword(s, "EXPLAIN") {
		return "", false
	}
	s = strings.TrimLeft(s[len("EXPLAIN"):], " \t\r\n\v\f")
	if hasKeyword(s, "FORMATTED") {
		return "JSON", true
	}
	return "TEXT", true
}

// pastLeadingNoise returns statement after its leading whitespace and
// comments, and reports whether it could be read at all. A block comment left
// open cannot be, and neither can one holding another `/*`: whether an engine
// ends a block comment at the first `*/` or nests them is a question about the
// engine, and Hive's rewrite has no parser to answer it, so the statement's
// first keyword is read as one the text does not settle.
//
// Whichever delimiter comes first decides, one position at a time, because the
// two can overlap: the second `/*` of `/*/*/` ends where the first `*/` begins,
// so searching for the closer and then looking back for an opener finds none.
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

// pastBlockComment returns what follows the block comment opening s, and
// reports whether it could be read at all.
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

// hasKeyword reports whether s starts with keyword, which ends there: the
// character after it cannot continue an identifier, so EXPLAINED is not an
// EXPLAIN.
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
