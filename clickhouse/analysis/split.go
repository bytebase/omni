// Package analysis provides lexical SQL analysis for ClickHouse: statement
// splitting, EXPLAIN rewriting and inspection, result-limit wrapping, and
// read-only query validation.
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
	var list []Statement
	if err := applyMultiStatements(script, func(text string, start, end int) error {
		list = append(list, Statement{
			Text:  text,
			Empty: isEmptySQL(text),
			Start: start,
			End:   end,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
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
