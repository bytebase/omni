package analysis

import "strings"

// ExplainStatement returns the ClickHouse EXPLAIN statement for statement.
// When format is "JSON" (case-insensitive), it requests `json = 1, indexes = 1`
// so reads report part and granule index pruning. A statement that already
// begins with EXPLAIN is returned unchanged so a typed EXPLAIN keeps its kind.
func ExplainStatement(statement, format string) string {
	if StartsWithExplain(statement) {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN json = 1, indexes = 1 " + statement
	}
	return "EXPLAIN " + statement
}

// StartsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading comments and whitespace. A block comment left open or
// holding another `/*` cannot be settled without a parser and returns false.
func StartsWithExplain(statement string) bool {
	return startsWithKeyword(statementBody(statement), "EXPLAIN")
}

// DescribeExplain returns the plan format ("JSON" or "TEXT") and whether
// statement executes what it plans (EXPLAIN ANALYZE). ok is false when
// statement does not start with EXPLAIN.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	if !StartsWithExplain(statement) {
		return "", false, false
	}
	return explainFormatOf(statement), explainAnalyzes(statement), true
}

func explainAnalyzes(statement string) bool {
	body, err := removeQuotedTextAndComment(statement)
	if err != nil {
		return false
	}
	tokens := keywordTokens(body)
	return len(tokens) > 1 && tokens[1].word && strings.EqualFold(tokens[1].text, "ANALYZE")
}

func explainFormatOf(statement string) string {
	body, err := removeQuotedTextAndComment(statement)
	if err != nil {
		return "TEXT"
	}
	tokens := keywordTokens(body)
	// Skip EXPLAIN and its optional kind, unless the first word is already a
	// setting (json = 1).
	i := 1
	if i < len(tokens) && tokens[i].word && !(i+1 < len(tokens) && tokens[i+1].text == "=") {
		kind := tokens[i].text
		i++
		if strings.EqualFold(kind, "QUERY") || strings.EqualFold(kind, "TABLE") {
			i++
		}
		if strings.EqualFold(kind, "JSON") {
			return "JSON"
		}
	}
	// The settings are name = value pairs; the query body starts at the first
	// word that is not one.
	for i < len(tokens) {
		if !tokens[i].word || i+1 >= len(tokens) || tokens[i+1].text != "=" || i+2 >= len(tokens) {
			break
		}
		if strings.EqualFold(tokens[i].text, "JSON") {
			if settingIsOff(tokens[i+2].text) {
				return "TEXT"
			}
			return "JSON"
		}
		i += 3
		if i >= len(tokens) || tokens[i].text != "," {
			break
		}
		i++
	}
	return "TEXT"
}

func settingIsOff(value string) bool {
	switch strings.ToUpper(value) {
	case "0", "FALSE", "OFF", "NO":
		return true
	}
	return false
}

// statementBody returns statement without its leading comments and whitespace,
// or "" when a block comment is unclosed or contains a nested `/*`.
func statementBody(statement string) string {
	s := statement
	for {
		s = strings.TrimLeft(s, " \t\r\n\v\f")
		switch {
		case strings.HasPrefix(s, "--"), strings.HasPrefix(s, "#"):
			_, rest, found := strings.Cut(s, "\n")
			if !found {
				return ""
			}
			s = rest
		case strings.HasPrefix(s, "/*"):
			rest, closed := pastBlockComment(s)
			if !closed {
				return ""
			}
			s = rest
		default:
			return s
		}
	}
}

// pastBlockComment returns what follows the block comment opening s, and
// reports whether it could be read at all. A comment left open cannot be, and
// neither can one holding another `/*`.
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

// startsWithKeyword reports whether s begins with keyword as a whole word.
func startsWithKeyword(s, keyword string) bool {
	if len(s) < len(keyword) || !strings.EqualFold(s[:len(keyword)], keyword) {
		return false
	}
	if len(s) == len(keyword) {
		return true
	}
	c := s[len(keyword)]
	return c != '_' && c != '$' && (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && c < 0x80
}
