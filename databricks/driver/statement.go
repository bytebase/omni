package driver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/bytebase/omni/driver"
)

// splitCommands cuts script into the statements a request runs, in order, with
// the text and byte range a task log reports.
//
// omni has no databricks tree, and Bytebase split Databricks SQL with
// parser/standard's applyMultiStatements, so this is that splitter moved here
// as clickhouse/driver and hive/driver carry it: a statement ends at a line
// ending in a semicolon, `--` and `/* */` comments and blank lines between
// statements belong to none, and a script in a `DELIMITER ;;` block is cut at
// the block's own terminator. It keeps that splitter's reading of a block
// comment that starts a line and of a line naming DELIMITER, including the
// error it answers for a `*/` that does not end its line. When a databricks
// tree lands, the splitter belongs there and this package should call it.
func splitCommands(script string) ([]driver.Command, error) {
	var commands []driver.Command
	err := applyMultiStatements(script, func(text string, start, end int) error {
		commands = append(commands, driver.Command{Text: text, ByteStart: start, ByteEnd: end})
		return nil
	})
	return commands, err
}

// applyMultiStatements splits the statement by semicolons and invokes f for
// each sub-statement with the text slice from the original statement and its
// [start, end) byte offsets. The text for each statement includes any leading
// whitespace from the end of the previous statement, matching the original
// behavior needed for position tracking.
func applyMultiStatements(statement string, f func(text string, start, end int) error) error {
	reader := bufio.NewReader(strings.NewReader(statement))
	delimiter := false
	comment := false
	inlineComment := false
	quote := rune(0)
	done := false
	hasContent := false
	// byteOffset tracks our read position in the original statement.
	byteOffset := 0
	// prevEnd tracks where the last emitted statement ended; the next
	// statement's text will start here (to include inter-statement whitespace).
	prevEnd := 0
	// contentEnd is the byte offset just past the last content character of
	// the statement being accumulated (excludes trailing whitespace/newlines).
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
			if scanLineQuotes(line, &quote, &inlineComment) {
				execute = true
			}
		} else {
			endsWithSemicolon := false
			if !comment && (inlineComment || !strings.HasPrefix(line, "/*")) {
				endsWithSemicolon = scanLineQuotes(line, &quote, &inlineComment)
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
				// skip line in comment mode
			case comment && strings.Contains(line, "*/"):
				if !strings.HasSuffix(line, "*/") {
					return errCommentEnd
				}
				comment = false
			case !hasContent && line == "":
				// skip leading blank lines
			case strings.HasPrefix(line, "--"):
				// skip comment lines
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
	// Apply the remaining content.
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
func scanLineQuotes(line string, quote *rune, inlineComment *bool) bool {
	last := byte(0)
	for i := 0; i < len(line); {
		c := line[i]
		if *inlineComment {
			if c == '*' && i+1 < len(line) && line[i+1] == '/' {
				*inlineComment = false
				i += 2
				continue
			}
			i++
			continue
		}
		if *quote != 0 {
			if c == '\\' && *quote == '\'' && i+1 < len(line) {
				i += 2
				last = 'x'
				continue
			}
			if c == byte(*quote) {
				if i+1 < len(line) && line[i+1] == c {
					i += 2
					last = c
					continue
				}
				*quote = 0
			}
			i++
			last = c
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
			break
		}
		if c == '\'' || c == '"' || c == '`' {
			*quote = rune(c)
		}
		i++
		last = c
	}
	return *quote == 0 && !*inlineComment && last == ';'
}

// errCommentEnd is the splitter's answer to a block comment whose closing does
// not end its line, as Bytebase's splitter wrote it.
var errCommentEnd = errors.New("`*/` must be the end of the line; new statement should start as a new line")

// trimStatement drops the whitespace and semicolons around a statement, so a
// split segment's terminator does not reach the server.
func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

// classifyStatement reports whether statement is expected to return rows. It
// is a lexical reading, because there is no databricks parser to ask: SELECT,
// WITH (unless it writes), SHOW, DESCRIBE, DESC, EXPLAIN, TABLE and VALUES are
// reads, and everything else is a command. A statement the reader cannot scan
// is a command; the server still sees it and reports its own error.
func classifyStatement(statement string) bool {
	text, err := sanitizeSQL(statement)
	if err != nil {
		return false
	}
	fields := strings.Fields(strings.ToUpper(strings.TrimSpace(text)))
	if len(fields) == 0 {
		return false
	}
	first := strings.TrimLeft(fields[0], "(")
	switch first {
	case "SELECT", "SHOW", "DESCRIBE", "DESC", "EXPLAIN", "TABLE", "VALUES":
		return true
	case "WITH":
		for _, keyword := range []string{"INSERT", "UPDATE", "DELETE", "MERGE"} {
			for _, field := range fields[1:] {
				if strings.Trim(field, "(),;") == keyword {
					return false
				}
			}
		}
		return true
	default:
		return false
	}
}

// sanitizeSQL replaces string literals and comments with spaces, keeping the
// bytes that carry statement keywords. It reports an unterminated literal or
// comment rather than guessing.
func sanitizeSQL(statement string) (string, error) {
	var out strings.Builder
	i := 0
	for i < len(statement) {
		switch {
		case strings.HasPrefix(statement[i:], "--"):
			i += 2
			for i < len(statement) && statement[i] != '\n' {
				i++
			}
			out.WriteByte(' ')
		case strings.HasPrefix(statement[i:], "/*"):
			i += 2
			for {
				if i >= len(statement) {
					return "", errors.New("invalid comment: not found */, but found EOF")
				}
				if strings.HasPrefix(statement[i:], "*/") {
					i += 2
					break
				}
				i++
			}
			out.WriteByte(' ')
		case statement[i] == '\'' || statement[i] == '"' || statement[i] == '`':
			quote := statement[i]
			i++
			for {
				if i >= len(statement) {
					return "", fmt.Errorf("invalid quoted text: not found %c, but found EOF", quote)
				}
				if statement[i] == '\\' && quote == '\'' {
					i += 2
					continue
				}
				if statement[i] == quote {
					if i+1 < len(statement) && statement[i+1] == quote {
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
			out.WriteString(" x ")
		default:
			out.WriteByte(statement[i])
			i++
		}
	}
	return out.String(), nil
}

// isExplainStatement reports whether the first keyword of statement is
// EXPLAIN, ignoring leading comments. It stands in for an AST check on an
// engine with no omni tree.
func isExplainStatement(statement string) bool {
	text, err := sanitizeSQL(statement)
	if err != nil {
		return false
	}
	fields := strings.Fields(strings.ToUpper(text))
	return len(fields) > 0 && strings.TrimLeft(fields[0], "(") == "EXPLAIN"
}

// explainStatement asks for Spark's FORMATTED plan, which is a tree of numbered
// operators followed by a block of detail for each. A statement that already
// asks for a plan runs as the caller wrote it, keeping the mode it named.
//
// Databricks prints this one plan whatever the caller asks for, as Bytebase's
// standard explain rewrite did.
func explainStatement(statement, _ string) string {
	if isExplainStatement(statement) {
		return statement
	}
	return "EXPLAIN FORMATTED " + statement
}

// planExecuted reports whether statement's EXPLAIN is the ANALYZE kind, whose
// output is a plan of an execution that happened.
func planExecuted(statement string) bool {
	text, err := sanitizeSQL(statement)
	if err != nil {
		return false
	}
	fields := strings.Fields(strings.ToUpper(text))
	return len(fields) > 1 && strings.TrimLeft(fields[0], "(") == "EXPLAIN" && strings.Trim(fields[1], "(),;") == "ANALYZE"
}
