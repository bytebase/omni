package driver

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/bytebase/omni/driver"
)

// splitCommands cuts script into the statements a request runs, in order, with
// the text and byte range a task log reports.
//
// Hive has no omni parser to ask, and HiveQL's statement boundaries are
// reached by reading the text, so this is Bytebase's standard splitter for
// Hive, parser/standard's applyMultiStatements, moved here: a statement ends at
// a line ending in a semicolon, `--` and `/* */` comments and blank lines
// between statements belong to none, and a script in a `DELIMITER ;;` block is
// cut at the block's own terminator. It keeps that splitter's reading of a
// block comment that starts a line and of a line naming DELIMITER, including
// the errors it answers for a `*/` that does not end its line.
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
	// Apply the remaining content.
	if hasContent {
		text := statement[prevEnd:contentEnd]
		if err := f(text, prevEnd, contentEnd); err != nil {
			return err
		}
	}
	return nil
}

// errCommentEnd is the splitter's answer to a block comment whose closing does
// not end its line, as Bytebase's splitter wrote it.
var errCommentEnd = errors.New("`*/` must be the end of the line; new statement should start as a new line")
