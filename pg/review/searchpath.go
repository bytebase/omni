package review

import "strings"

// parseSearchPath splits a search_path setting as SHOW search_path prints
// it: comma-separated entries, each an unquoted name (case-folded), a
// double-quoted identifier, or a single-quoted string. "$user" is kept as
// is; the catalog expands it to the current user. Nothing is dropped: a
// schema that does not exist is skipped at lookup time, as in PostgreSQL.
func parseSearchPath(setting string) []string {
	var path []string
	for _, token := range splitSearchPath(setting) {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		switch {
		case len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"':
			path = append(path, strings.ReplaceAll(token[1:len(token)-1], `""`, `"`))
		case len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'':
			path = append(path, strings.ReplaceAll(token[1:len(token)-1], `''`, `'`))
		default:
			path = append(path, strings.ToLower(token))
		}
	}
	return path
}

// splitSearchPath splits on the commas outside quotes.
func splitSearchPath(setting string) []string {
	var parts []string
	start := 0
	var quote byte
	for i := 0; i < len(setting); i++ {
		switch {
		case quote != 0:
			if setting[i] != quote {
				continue
			}
			if i+1 < len(setting) && setting[i+1] == quote {
				i++
				continue
			}
			quote = 0
		case setting[i] == '"' || setting[i] == '\'':
			quote = setting[i]
		case setting[i] == ',':
			parts = append(parts, setting[start:i])
			start = i + 1
		}
	}
	return append(parts, setting[start:])
}
