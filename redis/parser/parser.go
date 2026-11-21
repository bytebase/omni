// Package parser parses and validates Redis command-line scripts against the
// pinned Redis 7.4.12 command and subcommand table.
//
// Redis has no SQL statement grammar; its command-line syntax consists of
// newline-delimited commands whose arguments follow shell-like quoting rules
// (matching redis/analysis and redis/driver) plus command name, subcommand,
// and exact/minimum arity validation from Redis's own COMMAND metadata.
package parser

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Command is one parsed Redis command line.
type Command struct {
	// Name is the normalized lowercase top-level command name (e.g. "get", "acl").
	Name string
	// Subcommand is the normalized lowercase subcommand name for container
	// commands (e.g. "list" for "ACL LIST"), or empty otherwise.
	Subcommand string
	// Args holds the tokenized arguments as written, with Args[0] being the
	// raw command token and Args[1] (when Subcommand != "") the subcommand token.
	Args []string
	// Text is the raw source line the command was parsed from.
	Text string
	// Start and End are the byte offsets of Text in the input script.
	Start int
	End   int
}

// ParseError reports a lexical, unknown-command, unknown-subcommand, or arity
// error while parsing a Redis command script.
type ParseError struct {
	Line    int // 1-based line number
	Offset  int // byte offset of the offending line in the input
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Message)
}

// Parse parses a newline-delimited Redis command script. Blank and
// whitespace-only lines are skipped. Each non-empty line is tokenized and
// validated against Redis 7.4.12's command/subcommand names and arity rules.
func Parse(script string) ([]Command, error) {
	var commands []Command
	start := 0
	lineNo := 1
	for line := range strings.SplitSeq(script, "\n") {
		end := start + len(line)
		fields, err := Tokenize(line)
		if err != nil {
			return nil, &ParseError{
				Line:    lineNo,
				Offset:  start,
				Message: err.Error(),
			}
		}
		if len(fields) > 0 {
			cmd, err := validateCommand(fields, line, start, end)
			if err != nil {
				return nil, &ParseError{
					Line:    lineNo,
					Offset:  start,
					Message: err.Error(),
				}
			}
			commands = append(commands, cmd)
		}
		start = end + 1
		lineNo++
	}
	return commands, nil
}

// ParseCommand parses and validates a single Redis command line. Empty lines
// and multi-line strings (containing '\n') are rejected.
func ParseCommand(line string) (*Command, error) {
	if strings.ContainsRune(line, '\n') {
		return nil, &ParseError{
			Line:    1,
			Offset:  strings.IndexByte(line, '\n'),
			Message: "unexpected newline in single command",
		}
	}
	fields, err := Tokenize(line)
	if err != nil {
		return nil, &ParseError{
			Line:    1,
			Offset:  0,
			Message: err.Error(),
		}
	}
	if len(fields) == 0 {
		return nil, &ParseError{
			Line:    1,
			Offset:  0,
			Message: "empty command",
		}
	}
	cmd, err := validateCommand(fields, line, 0, len(line))
	if err != nil {
		return nil, &ParseError{
			Line:    1,
			Offset:  0,
			Message: err.Error(),
		}
	}
	return &cmd, nil
}

// Lookup returns the Redis 7.4.12 CommandSpec for a top-level command name
// (case-insensitive). The returned spec's Subcommands map is a fresh copy.
func Lookup(name string) (CommandSpec, bool) {
	spec, ok := commandRegistry[strings.ToLower(name)]
	if !ok {
		return CommandSpec{}, false
	}
	if spec.Subcommands != nil {
		spec.Subcommands = maps.Clone(spec.Subcommands)
	}
	return spec, true
}

// Commands returns all registered Redis 7.4.12 top-level command names in
// sorted lowercase order.
func Commands() []string {
	names := make([]string, 0, len(commandRegistry))
	for name := range commandRegistry {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func validateCommand(fields []string, text string, start, end int) (Command, error) {
	cmdName := strings.ToLower(fields[0])
	spec, ok := commandRegistry[cmdName]
	if !ok {
		return Command{}, fmt.Errorf("unknown command %q", fields[0])
	}
	if len(spec.Subcommands) > 0 {
		if len(fields) == 1 {
			if spec.StandaloneArity == 0 || !matchArity(spec.StandaloneArity, 1) {
				return Command{}, fmt.Errorf("wrong number of arguments for '%s' command", cmdName)
			}
			return Command{
				Name:  cmdName,
				Args:  fields,
				Text:  text,
				Start: start,
				End:   end,
			}, nil
		}
		subName := strings.ToLower(fields[1])
		subArity, ok := spec.Subcommands[subName]
		if !ok {
			return Command{}, fmt.Errorf("unknown subcommand %q for '%s'", fields[1], cmdName)
		}
		if !matchArity(subArity, len(fields)) {
			return Command{}, fmt.Errorf("wrong number of arguments for '%s|%s' command", cmdName, subName)
		}
		return Command{
			Name:       cmdName,
			Subcommand: subName,
			Args:       fields,
			Text:       text,
			Start:      start,
			End:        end,
		}, nil
	}

	if !matchArity(spec.StandaloneArity, len(fields)) {
		return Command{}, fmt.Errorf("wrong number of arguments for '%s' command", cmdName)
	}
	return Command{
		Name:  cmdName,
		Args:  fields,
		Text:  text,
		Start: start,
		End:   end,
	}, nil
}

func matchArity(arity, argc int) bool {
	if arity > 0 {
		return argc == arity
	}
	if arity < 0 {
		return argc >= -arity
	}
	return false
}

// Tokenize splits one Redis command line into its arguments using Redis's
// shell-like command-line rules (matching redis/analysis and redis/driver):
//   - spaces, tabs, and carriage returns separate arguments;
//   - single quotes keep every character verbatim (including backslashes);
//   - double quotes interpret backslash escapes (\xNN, \n, \r, \t, \b, \a,
//     \\, \", and \X -> X);
//   - outside quotes, a backslash escapes the next character;
//   - `#` is ordinary argument text because Redis has no comment syntax;
//   - an unclosed quote or a trailing backslash returns an error.
func Tokenize(line string) ([]string, error) {
	var fields []string
	var field strings.Builder
	// started marks that field is an argument even when empty, so "" or ''
	// produces an empty-string argument.
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
					switch line[i] {
					case 'n':
						field.WriteByte('\n')
					case 'r':
						field.WriteByte('\r')
					case 't':
						field.WriteByte('\t')
					case 'b':
						field.WriteByte('\b')
					case 'a':
						field.WriteByte('\a')
					case 'x':
						if i+2 < len(line) {
							if b, ok := parseHexByte(line[i+1], line[i+2]); ok {
								field.WriteByte(b)
								i += 2
								continue
							}
						}
						field.WriteByte('x')
					default:
						field.WriteByte(line[i])
					}
					continue
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

func parseHexByte(hi, lo byte) (byte, bool) {
	h, ok1 := hexNibble(hi)
	l, ok2 := hexNibble(lo)
	if !ok1 || !ok2 {
		return 0, false
	}
	return (h << 4) | l, true
}

func hexNibble(b byte) (byte, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}
