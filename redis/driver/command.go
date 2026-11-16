package driver

import (
	"errors"
	"fmt"
	"strings"
)

// command is one command of a script: its arguments, the line it was written on,
// and that line's byte range in the script.
type command struct {
	args  []any
	text  string
	start int
	end   int
}

// splitCommands reads one command per line and skips lines that hold none. A line
// is read as a shell reads one, so a quoted argument may hold blanks, but Redis
// has no comment syntax: a `#` is an argument's text.
func splitCommands(statement string) ([]command, error) {
	var commands []command
	start := 0
	for _, line := range strings.Split(statement, "\n") {
		end := start + len(line)
		fields, err := splitLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to split command %s: %w", line, err)
		}
		if len(fields) > 0 {
			args := make([]any, len(fields))
			for i, field := range fields {
				args[i] = field
			}
			commands = append(commands, command{args: args, text: line, start: start, end: end})
		}
		start = end + 1
	}
	return commands, nil
}

// splitLine reads one command line into its arguments. Blanks separate them; a
// single quote keeps its text as written; a double quote keeps its text with a
// backslash escaping the next character; outside quotes a backslash escapes the
// next character too. A quote or an escape the line does not close is an error.
func splitLine(line string) ([]string, error) {
	var fields []string
	var field strings.Builder
	// started reports that field is an argument even when it is empty, so a
	// quoted empty string stays one.
	started := false
	appendField := func() {
		fields = append(fields, field.String())
		field.Reset()
		started = false
	}
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case ' ', '\t', '\r':
			if started {
				appendField()
			}
		case '\'':
			started = true
			// A single quote keeps every character as written, the backslash
			// included.
			closing := strings.IndexByte(line[i+1:], '\'')
			if closing < 0 {
				return nil, errors.New("EOF found when expecting closing quote")
			}
			field.WriteString(line[i+1 : i+1+closing])
			i += closing + 1
		case '"':
			started = true
			i++
			for ; i < len(line) && line[i] != '"'; i++ {
				if line[i] == '\\' {
					if i+1 >= len(line) {
						return nil, errors.New("EOF found after escape character")
					}
					i++
				}
				field.WriteByte(line[i])
			}
			if i >= len(line) {
				return nil, errors.New("EOF found when expecting closing quote")
			}
		case '\\':
			if i+1 >= len(line) {
				return nil, errors.New("EOF found after escape character")
			}
			i++
			started = true
			field.WriteByte(line[i])
		default:
			started = true
			field.WriteByte(line[i])
		}
	}
	if started {
		appendField()
	}
	return fields, nil
}
