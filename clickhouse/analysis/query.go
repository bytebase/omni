package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	explainRe      = regexp.MustCompile(`^EXPLAIN\s+?`)
	explainAnalyze = regexp.MustCompile(`^EXPLAIN\s+ANALYZE\s+?`)
	cteRe          = regexp.MustCompile(`^WITH\s+?`)
)

// ValidateQuery reports whether statement is a read-only query that returns
// rows. It strips quoted text and comments, then classifies the leading
// keyword: SELECT, parenthesized SELECT, SHOW/DESCRIBE/DESC/EXISTS/CHECK,
// EXPLAIN (excluding EXPLAIN ANALYZE), or a WITH query without top-level DML.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	text, err := removeQuotedTextAndComment(statement)
	if err != nil {
		return false, false, err
	}
	ok := CheckStatementWithoutQuotedTextAndComment(text)
	return ok, ok, nil
}

// StatementReadsRows reports whether the server answers statement with rows.
// EXPLAIN ANALYZE is classified as non-read-only by ValidateQuery because it
// executes what it plans, yet the server still returns the plan as rows.
func StatementReadsRows(statement string) (bool, error) {
	readOnly, _, err := ValidateQuery(statement)
	if err != nil {
		return false, err
	}
	return readOnly || StartsWithExplain(statement), nil
}

// CheckStatementWithoutQuotedTextAndComment classifies a single statement
// whose quoted text and comments have already been removed.
func CheckStatementWithoutQuotedTextAndComment(statement string) bool {
	formatted := strings.ToUpper(strings.TrimSpace(statement))
	// Leading-keyword classification speaks for one statement only, and the
	// line-oriented splitter leaves a one-line batch whole: "SELECT 1; DROP
	// TABLE t" would otherwise pass on its SELECT.
	if strings.Contains(strings.TrimSuffix(formatted, ";"), ";") {
		return false
	}
	for _, keyword := range []string{"SHOW", "DESCRIBE", "DESC", "EXISTS", "CHECK"} {
		if startsWithKeyword(formatted, keyword) {
			return true
		}
	}
	if startsWithKeyword(formatted, "SELECT") || strings.HasPrefix(formatted, "(") {
		return true
	}
	if explainRe.MatchString(formatted) {
		return !explainAnalyze.MatchString(formatted)
	}
	if cteRe.MatchString(formatted) {
		return withStatementIsQuery(formatted)
	}
	return false
}

// withStatementIsQuery reads the statement a WITH clause runs, ignoring words
// inside its CTE bodies and the aliases they define.
func withStatementIsQuery(statement string) bool {
	tokens := keywordTokens(statement)
	depth := 0
	for i, token := range tokens {
		switch token.text {
		case "(":
			depth++
			continue
		case ")":
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 || !token.word {
			continue
		}
		switch token.text {
		case "SELECT":
			if previousWord(tokens, i) == "AS" {
				continue
			}
			if nextWord(tokens, i) == "AS" {
				switch previousToken(tokens, i) {
				case "WITH", ",", "(":
					continue
				}
			}
			return true
		case "INSERT":
			if nextWord(tokens, i) == "INTO" {
				return false
			}
		case "UPDATE":
			next := nextWord(tokens, i)
			if next != "" && next != "SELECT" && next != "SET" && next != "AS" {
				return false
			}
		case "DELETE":
			if nextWord(tokens, i) == "FROM" {
				return false
			}
		}
	}
	return false
}

type keywordToken struct {
	text string
	word bool
}

func keywordTokens(s string) []keywordToken {
	runes := []rune(s)
	var tokens []keywordToken
	for i := 0; i < len(runes); {
		c := runes[i]
		switch {
		case isWordRune(c):
			start := i
			for i < len(runes) && isWordRune(runes[i]) {
				i++
			}
			tokens = append(tokens, keywordToken{text: strings.ToUpper(string(runes[start:i])), word: true})
		case unicode.IsSpace(c):
			i++
		default:
			tokens = append(tokens, keywordToken{text: string(c)})
			i++
		}
	}
	return tokens
}

func isWordRune(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

func previousWord(tokens []keywordToken, i int) string {
	for j := i - 1; j >= 0; j-- {
		if tokens[j].word {
			return tokens[j].text
		}
	}
	return ""
}

func previousToken(tokens []keywordToken, i int) string {
	if i == 0 {
		return ""
	}
	return tokens[i-1].text
}

func nextWord(tokens []keywordToken, i int) string {
	for j := i + 1; j < len(tokens); j++ {
		if tokens[j].word {
			return tokens[j].text
		}
	}
	return ""
}

func removeQuotedTextAndComment(statement string) (string, error) {
	runes := []rune(statement)
	var buf strings.Builder
	i := 0
	for i < len(runes) && unicode.IsSpace(runes[i]) {
		i++
	}
	start := i
	for {
		switch {
		case i+1 < len(runes) && runes[i] == '/' && runes[i+1] == '*':
			buf.WriteString(string(runes[start:i]))
			if err := scanComment(runes, &i); err != nil {
				return "", err
			}
			buf.WriteByte(' ')
			start = i
		case i+1 < len(runes) && runes[i] == '-' && runes[i+1] == '-':
			buf.WriteString(string(runes[start:i]))
			if err := scanComment(runes, &i); err != nil {
				return "", err
			}
			buf.WriteByte(' ')
			start = i
		case i < len(runes) && runes[i] == '#':
			buf.WriteString(string(runes[start:i]))
			if err := scanComment(runes, &i); err != nil {
				return "", err
			}
			buf.WriteByte(' ')
			start = i
		case i < len(runes) && runes[i] == '\'':
			buf.WriteString(string(runes[start:i]))
			if err := scanString(runes, &i); err != nil {
				return "'", err
			}
			buf.WriteString(" x ")
			start = i
		case i < len(runes) && runes[i] == '"':
			buf.WriteString(string(runes[start:i]))
			if err := scanIdentifier(runes, &i, '"'); err != nil {
				return "", err
			}
			buf.WriteString(" x ")
			start = i
		case i < len(runes) && runes[i] == '`':
			buf.WriteString(string(runes[start:i]))
			if err := scanIdentifier(runes, &i, '`'); err != nil {
				return "", err
			}
			buf.WriteString(" x ")
			start = i
		case i < len(runes) && runes[i] == '$' && isHeredocStart(runes, i):
			buf.WriteString(string(runes[start:i]))
			if err := scanHeredoc(runes, &i); err != nil {
				return "", err
			}
			buf.WriteString(" x ")
			start = i
		case i >= len(runes):
			buf.WriteString(string(runes[start:]))
			return buf.String(), nil
		default:
			i++
		}
	}
}

func isHeredocStart(runes []rune, i int) bool {
	_, ok := heredocTag(runes, i)
	return ok
}

func heredocTag(runes []rune, i int) (string, bool) {
	if i >= len(runes) || runes[i] != '$' {
		return "", false
	}
	j := i + 1
	for j < len(runes) && isWordRune(runes[j]) {
		j++
	}
	if j >= len(runes) || runes[j] != '$' {
		return "", false
	}
	return string(runes[i : j+1]), true
}

func scanHeredoc(runes []rune, i *int) error {
	tag, ok := heredocTag(runes, *i)
	if !ok {
		return nil
	}
	tagRunes := []rune(tag)
	*i += len(tagRunes)
	for *i < len(runes) {
		if *i+len(tagRunes) <= len(runes) && string(runes[*i:*i+len(tagRunes)]) == tag {
			*i += len(tagRunes)
			return nil
		}
		*i++
	}
	return fmt.Errorf("invalid heredoc: not found delimiter: %s, but found EOF", tag)
}

func scanComment(runes []rune, i *int) error {
	switch {
	case *i+1 < len(runes) && runes[*i] == '/' && runes[*i+1] == '*':
		*i += 2
		for {
			switch {
			case *i+1 < len(runes) && runes[*i] == '*' && runes[*i+1] == '/':
				*i += 2
				return nil
			case *i >= len(runes):
				return errors.New("invalid comment: not found */, but found EOF")
			default:
				*i++
			}
		}
	case *i+1 < len(runes) && runes[*i] == '-' && runes[*i+1] == '-':
		*i += 2
		for *i < len(runes) && runes[*i] != '\n' {
			*i++
		}
		if *i < len(runes) {
			*i++
		}
		return nil
	case *i < len(runes) && runes[*i] == '#':
		*i++
		for *i < len(runes) && runes[*i] != '\n' {
			*i++
		}
		if *i < len(runes) {
			*i++
		}
		return nil
	default:
		return errors.New("no comment found")
	}
}

func scanString(runes []rune, i *int) error {
	*i++
	for {
		switch {
		case *i >= len(runes):
			return errors.New("invalid string: not found delimiter: ', but found EOF")
		case runes[*i] == '\'':
			*i++
			return nil
		case runes[*i] == '\\':
			*i += 2
		default:
			*i++
		}
	}
}

func scanIdentifier(runes []rune, i *int, delimiter rune) error {
	*i++
	for {
		switch {
		case *i >= len(runes):
			return fmt.Errorf("invalid indentifier: not found delimiter: %c, but found EOF", delimiter)
		case runes[*i] == '\\':
			*i += 2
		case runes[*i] == delimiter:
			*i++
			return nil
		default:
			*i++
		}
	}
}
