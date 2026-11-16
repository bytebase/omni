package driver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

// statement is one statement of a script, with its byte range in the script.
type statement struct {
	Text      string
	ByteStart int
	ByteEnd   int
	Empty     bool
}

// splitSQL splits a script the way Bytebase's standard splitter did: a line
// holding the terminator, or a line ending in one, ends a statement. It is a
// port of backend/plugin/parser/standard's applyMultiStatements, which ClickHouse
// and Hive share; the newline-shaped rule is the defect that came with it.
func splitSQL(script string) ([]statement, error) {
	var list []statement
	if err := applyMultiStatements(script, func(text string, start, end int) error {
		list = append(list, statement{
			Text:      text,
			ByteStart: start,
			ByteEnd:   end,
			Empty:     isEmptySQL(text),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

func filterEmptyStatements(list []statement) []statement {
	var out []statement
	for _, s := range list {
		if !s.Empty {
			out = append(out, s)
		}
	}
	return out
}

// isEmptySQL reports whether the text holds only whitespace and comments.
func isEmptySQL(sql string) bool {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return true
	}
	if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") {
		lines := strings.Split(trimmed, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "--") {
				continue
			}
			if strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") {
				continue
			}
			return false
		}
		return true
	}
	return false
}

// applyMultiStatements splits the statement by semicolons and invokes f for
// each sub-statement with its byte range. The text of a statement includes the
// whitespace since the end of the previous one, which the execution log's
// ranges are drawn from.
func applyMultiStatements(statement string, f func(text string, start, end int) error) error {
	reader := bufio.NewReader(strings.NewReader(statement))
	delimiter := false
	comment := false
	done := false
	hasContent := false
	byteOffset := 0
	prevEnd := 0
	contentEnd := 0
	for !done {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			done = true
		}
		lineLen := len(line)
		line = strings.TrimRightFunc(line, unicode.IsSpace)

		execute := false
		switch {
		case strings.HasPrefix(line, "/*"):
			if strings.Contains(line, "*/") {
				if !strings.HasSuffix(line, "*/") {
					return errors.New("`*/` must be the end of the line; new statement should start as a new line")
				}
			} else {
				comment = true
			}
		case comment && !strings.Contains(line, "*/"):
		case comment && strings.Contains(line, "*/"):
			if !strings.HasSuffix(line, "*/") {
				return errors.New("`*/` must be the end of the line; new statement should start as a new line")
			}
			comment = false
		case !hasContent && line == "":
		case strings.HasPrefix(line, "--"):
		case line == "DELIMITER ;;":
			delimiter = true
		case line == "DELIMITER ;" && delimiter:
			delimiter = false
			execute = true
		case strings.HasSuffix(line, ";"):
			hasContent = true
			contentEnd = byteOffset + len(line)
			if !delimiter {
				execute = true
			}
		default:
			if line != "" {
				hasContent = true
				contentEnd = byteOffset + len(line)
			}
		}
		if execute && hasContent {
			text := statement[prevEnd:contentEnd]
			if err := f(text, prevEnd, contentEnd); err != nil {
				return fmt.Errorf("execute query %q failed: %w", text, err)
			}
			hasContent = false
			prevEnd = contentEnd
		}
		byteOffset += lineLen
	}
	if hasContent {
		text := statement[prevEnd:contentEnd]
		if err := f(text, prevEnd, contentEnd); err != nil {
			return fmt.Errorf("execute query %q failed: %w", text, err)
		}
	}
	return nil
}

// statementWithResultLimit bounds a statement with ClickHouse's CTE wrapper.
// ClickHouse is read as text only here, so there is no AST to place a LIMIT
// in, and every statement takes the wrapper.
func statementWithResultLimit(statement string, limit int) string {
	wrapped := trimStatementForWrapper(statement)
	// Only a statement that can stand as a subquery takes the wrapper. A
	// SHOW, DESCRIBE, INSERT or EXPLAIN is run as written.
	if !isSubqueryable(wrapped) {
		return statement
	}
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result LIMIT %d;", wrapped, limit)
}

// trimStatementForWrapper removes trailing comments and the terminators before
// them, which a CTE cannot hold.
func trimStatementForWrapper(statement string) string {
	trimmed := trimStatement(statement)
	for {
		start := trailingCommentStart(trimmed)
		if start < 0 {
			return trimmed
		}
		trimmed = trimStatement(strings.TrimSpace(trimmed[:start]))
	}
}

// isSubqueryable reports whether statement can stand inside a CTE.
func isSubqueryable(statement string) bool {
	body := statementBody(statement)
	if strings.HasPrefix(body, "(") {
		return true
	}
	upper := strings.ToUpper(body)
	for _, keyword := range []string{"SELECT", "WITH", "TABLE", "VALUES"} {
		if startsWithKeyword(upper, keyword) {
			return true
		}
	}
	return false
}

// trailingCommentStart is the index where a comment that reaches the end of
// statement starts, or -1. It reads strings, quoted identifiers and comments
// so a `--` inside one does not count.
func trailingCommentStart(statement string) int {
	runes := []rune(statement)
	last := -1
	line, block, str := false, false, false
	var ident rune
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case line:
			if c == '\n' {
				line = false
			}
		case block:
			if c == '*' && i+1 < len(runes) && runes[i+1] == '/' {
				block = false
				i++
			}
		case str:
			if c == '\\' {
				i++
			} else if c == '\'' {
				str = false
			}
		case ident != 0:
			if c == ident {
				ident = 0
			}
		default:
			switch {
			case c == '-' && i+1 < len(runes) && runes[i+1] == '-':
				line = true
				i++
				last = i - 1
			case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
				block = true
				i++
				last = i - 1
			case unicode.IsSpace(c):
			case c == '\'':
				str = true
				last = -1
			case c == '"' || c == '`':
				ident = c
				last = -1
			default:
				last = -1
			}
		}
	}
	return last
}

// trimStatement trims the whitespace and trailing semicolons a limit rewrite
// wraps.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
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

// explainStatement asks ClickHouse for the structured plan the query-plan
// viewer draws. `indexes = 1` is what makes each read report how many parts and
// granules its indexes let it skip, which is what a ClickHouse plan is read
// for. A statement that already asks for a plan runs as written, so a typed
// EXPLAIN keeps whichever of ClickHouse's plan kinds it named.
func explainStatement(statement, format string) string {
	if startsWithExplain(statement) {
		return statement
	}
	if strings.EqualFold(format, "JSON") {
		return "EXPLAIN json = 1, indexes = 1 " + statement
	}
	return "EXPLAIN " + statement
}

// startsWithExplain reports whether the first keyword of statement is EXPLAIN,
// ignoring leading comments and whitespace. It stands in for an AST check on an
// engine whose parser lives outside this tree.
func startsWithExplain(statement string) bool {
	return startsWithKeyword(statementBody(statement), "EXPLAIN")
}

// statementBody is statement without its leading comments and whitespace, or
// "" when a block comment is not closed.
func statementBody(statement string) string {
	s := statement
	for {
		s = strings.TrimLeft(s, " \t\r\n\v\f")
		switch {
		case strings.HasPrefix(s, "--"):
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

// validateSQLForEditor reports whether statement is a read-only query the
// editor may run with Query rather than Exec. It removes quoted text and
// comments, then classifies the leading keyword: SELECT, EXPLAIN (not EXPLAIN
// ANALYZE), or a WITH that does not hold INSERT, UPDATE or DELETE.
func validateSQLForEditor(statement string) (bool, error) {
	text, err := removeQuotedTextAndComment(statement)
	if err != nil {
		return false, err
	}
	return checkStatementWithoutQuotedTextAndComment(text), nil
}

var (
	selectRe        = regexp.MustCompile(`^SELECT\s+?`)
	selectStarRe    = regexp.MustCompile(`^SELECT\*\s+?`)
	explainRe       = regexp.MustCompile(`^EXPLAIN\s+?`)
	explainAnalyze  = regexp.MustCompile(`^EXPLAIN\s+ANALYZE\s+?`)
	cteRe           = regexp.MustCompile(`^WITH\s+?`)
	withDMLPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\bINSERT\b`),
		regexp.MustCompile(`\bUPDATE\b`),
		regexp.MustCompile(`\bDELETE\b`),
	}
)

func checkStatementWithoutQuotedTextAndComment(statement string) bool {
	formatted := strings.ToUpper(strings.TrimSpace(statement))
	// Leading-keyword classification speaks for one statement only, and the
	// ClickHouse splitter breaks on newlines rather than terminators, so a
	// one-line batch arrives whole: "SELECT 1; DROP TABLE t" would pass on its
	// SELECT.
	if strings.Contains(strings.TrimSuffix(formatted, ";"), ";") {
		return false
	}
	// Statements that answer with rows but do not start with SELECT or EXPLAIN.
	for _, keyword := range []string{"SHOW", "DESCRIBE", "DESC", "EXISTS", "CHECK"} {
		if startsWithKeyword(formatted, keyword) {
			return true
		}
	}
	if selectRe.MatchString(formatted) || selectStarRe.MatchString(formatted) {
		return true
	}
	if explainRe.MatchString(formatted) {
		return !explainAnalyze.MatchString(formatted)
	}
	if cteRe.MatchString(formatted) {
		for _, pattern := range withDMLPatterns {
			if pattern.MatchString(formatted) {
				return false
			}
		}
		return true
	}
	return false
}

// removeQuotedTextAndComment is Bytebase's
// tokenizer.StandardRemoveQuotedTextAndComment over the same grammar: a block
// comment runs to `*/`, a line comment to the newline, a string to its closing
// quote with backslash escaping, and a quoted identifier to its closing double
// quote. Each becomes one space.
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
		case i < len(runes) && runes[i] == '\'':
			buf.WriteString(string(runes[start:i]))
			if err := scanString(runes, &i); err != nil {
				return "'", err
			}
			buf.WriteByte(' ')
			start = i
		case i < len(runes) && runes[i] == '"':
			buf.WriteString(string(runes[start:i]))
			if err := scanIdentifier(runes, &i, '"'); err != nil {
				return "", err
			}
			buf.WriteByte(' ')
			start = i
		case i < len(runes) && runes[i] == '`':
			buf.WriteString(string(runes[start:i]))
			if err := scanIdentifier(runes, &i, '`'); err != nil {
				return "", err
			}
			buf.WriteByte(' ')
			start = i
		case i >= len(runes):
			buf.WriteString(string(runes[start:]))
			return buf.String(), nil
		default:
			i++
		}
	}
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
		case runes[*i] == delimiter:
			*i++
			return nil
		default:
			*i++
		}
	}
}
