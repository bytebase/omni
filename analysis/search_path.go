package analysis

import "strings"

// SearchPath returns the schemas a PostgreSQL-compatible search_path setting
// names, in order: each entry unquoted (a doubled quote undone) or
// lower-cased, with $user, system schemas and empty entries left out.
func SearchPath(setting string) []string {
	var named []string
	configured := 0
	for _, token := range splitSearchPath(setting) {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		schema := strings.ToLower(token)
		switch {
		case len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"':
			schema = strings.ReplaceAll(token[1:len(token)-1], `""`, `"`)
		case len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'':
			schema = strings.ReplaceAll(token[1:len(token)-1], `''`, `'`)
		default:
		}
		if schema == "$user" {
			configured++
			continue
		}
		if schema == "" || isSystemSearchPathSchema(schema) {
			continue
		}
		configured++
		named = append(named, schema)
	}
	// A setting with no entry is an empty path; one whose entries all name the
	// unset $user is no path at all.
	if configured == 0 {
		return []string{}
	}
	return named
}

// splitSearchPath splits a search_path setting at the commas outside quotes.
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
		default:
		}
	}
	return append(parts, setting[start:])
}

func isSystemSearchPathSchema(schema string) bool {
	for _, system := range []string{"pg_catalog", "information_schema", "pg_toast", "pg_temp_1", "pg_temp_2", "pg_global"} {
		if strings.EqualFold(schema, system) {
			return true
		}
	}
	return false
}
