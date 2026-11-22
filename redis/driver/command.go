package driver

import (
	"fmt"

	redisparser "github.com/bytebase/omni/redis/parser"
)

// command is one command of a script: its arguments, the line it was written on,
// and that line's byte range in the script.
type command struct {
	args  []any
	text  string
	start int
	end   int
}

// splitCommands reads one command per line and skips lines that hold none,
// delegating tokenization to redis/parser so driver and analysis share one
// lexer.
func splitCommands(statement string) ([]command, error) {
	stmts, err := redisparser.Split(statement)
	if err != nil {
		return nil, fmt.Errorf("failed to split command: %w", err)
	}
	if len(stmts) == 0 {
		return nil, nil
	}
	commands := make([]command, 0, len(stmts))
	for _, stmt := range stmts {
		args := make([]any, len(stmt.Args))
		for i, field := range stmt.Args {
			args[i] = field
		}
		commands = append(commands, command{
			args:  args,
			text:  stmt.Text,
			start: stmt.Start,
			end:   stmt.End,
		})
	}
	return commands, nil
}
