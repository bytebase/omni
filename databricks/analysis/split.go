// Package analysis provides lexical SQL analysis for Databricks SQL:
// statement splitting, EXPLAIN rewriting and inspection, and row-returning
// statement classification.
package analysis

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/bytebase/omni/analysis"
)

// Statement is one segment of a split SQL script.
type Statement = analysis.Statement

var errCommentEnd = errors.New("`*/` must be the end of the line; new statement should start as a new line")

// Split splits script into statements by line-ending semicolons (ignoring
// semicolons inside quoted strings and inline comments) and `DELIMITER ;;`
// blocks. Each segment's Text includes the whitespace since the previous
// segment so byte offsets stay contiguous for position tracking.
func Split(script string) ([]Statement, error) {
	var statements []Statement
	err := applyMultiStatements(script, func(text string, start, end int) error {
		statements = append(statements, Statement{
			Text:  text,
			Empty: isEmptySQL(text),
			Start: start,
			End:   end,
		})
		return nil
	})
	return statements, err
}

// TrimStatement drops surrounding whitespace and trailing semicolons from
// statement so a split segment's terminator does not reach the server.
func TrimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

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

func applyMultiStatements(statement string, f func(text string, start, end int) error) error {
	reader := bufio.NewReader(strings.NewReader(statement))
	delimiter := false
	comment := false
	inlineComment := false
	pendingTerminator := false
	quote := rune(0)
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
		if quote != 0 {
			// The previous line left a quoted string open, so this line is
			// its content, not a comment line.
			hasContent = true
			if line != "" {
				contentEnd = byteOffset + len(line)
			}
			if scanLineQuotes(line, &quote, &inlineComment, &pendingTerminator) {
				execute = true
			}
		} else {
			endsWithSemicolon := false
			if !comment && (inlineComment || !strings.HasPrefix(line, "/*")) {
				endsWithSemicolon = scanLineQuotes(line, &quote, &inlineComment, &pendingTerminator)
			}
			switch {
			case strings.HasPrefix(line, "/*"):
				if strings.Contains(line, "*/") {
					if !strings.HasSuffix(line, "*/") {
						return errCommentEnd
					}
				} else {
					comment = true
				}
			case comment && !strings.Contains(line, "*/"):
			case comment && strings.Contains(line, "*/"):
				if !strings.HasSuffix(line, "*/") {
					return errCommentEnd
				}
				comment = false
			case !hasContent && line == "":
			case strings.HasPrefix(line, "--"):
			case line == "DELIMITER ;;":
				delimiter = true
			case line == "DELIMITER ;" && delimiter:
				delimiter = false
				execute = true
			case endsWithSemicolon:
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
		}
		if execute && hasContent {
			text := statement[prevEnd:contentEnd]
			if err := f(text, prevEnd, contentEnd); err != nil {
				return err
			}
			hasContent = false
			prevEnd = contentEnd
		}
		byteOffset += lineLen
	}
	if hasContent {
		text := statement[prevEnd:contentEnd]
		if err := f(text, prevEnd, contentEnd); err != nil {
			return err
		}
	}
	return nil
}

// scanLineQuotes advances quote and inlineComment across one line, ignoring
// comment text outside a quoted string, and reports whether the last content
// outside both is a semicolon. A block comment that opens after SQL text keeps
// its state on the next line.
func scanLineQuotes(line string, quote *rune, inlineComment, pendingTerminator *bool) bool {
	last := byte(0)
	for i := 0; i < len(line); {
		c := line[i]
		if *inlineComment {
			if c == '*' && i+1 < len(line) && line[i+1] == '/' {
				*inlineComment = false
				if *pendingTerminator {
					last = ';'
					*pendingTerminator = false
				}
				i += 2
				continue
			}
			i++
			continue
		}
		if *quote != 0 {
			if c == '\\' && *quote == '\'' && i+1 < len(line) {
				i += 2
				continue
			}
			if c == byte(*quote) {
				if i+1 < len(line) && line[i+1] == c {
					i += 2
					continue
				}
				*quote = 0
			}
			i++
			continue
		}
		if c == '-' && i+1 < len(line) && line[i+1] == '-' {
			break
		}
		if c == '/' && i+1 < len(line) && line[i+1] == '*' {
			if end := strings.Index(line[i+2:], "*/"); end >= 0 {
				i += 2 + end + 2
				continue
			}
			*inlineComment = true
			*pendingTerminator = last == ';'
			break
		}
		if c == '\'' || c == '"' || c == '`' {
			*quote = rune(c)
			i++
			continue
		}
		if c != ' ' && c != '\t' && c != '\r' {
			last = c
		}
		i++
	}
	return *quote == 0 && !*inlineComment && last == ';'
}
