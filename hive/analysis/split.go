// Package analysis provides lexical SQL analysis for Hive: statement
// splitting and EXPLAIN rewriting and plan format detection.
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

// Split splits script into statements by line-ending semicolons and
// `DELIMITER ;;` blocks. Each segment's Text includes the whitespace since the
// previous segment so byte offsets stay contiguous for position tracking.
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
