package analysis

import (
	"errors"
	"fmt"
	"strings"
)

// ClassifyStatement reports whether statement is expected to return rows.
// SELECT, WITH (unless it writes), SHOW, DESCRIBE, DESC, EXPLAIN, TABLE, and
// VALUES are reads; everything else (or an unlexable statement) is a command.
func ClassifyStatement(statement string) bool {
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

// sanitizeSQL replaces string literals, quoted identifiers, and comments with
// placeholders/spaces, keeping the bytes that carry statement keywords.
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
