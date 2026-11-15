package analysis

import "strings"

func SearchPathForSelectedSchema(selectedSchema string, defaultSearchPath []string) []string {
	if selectedSchema == "" {
		return defaultSearchPath
	}
	if selectedSchema == "public" || strings.ContainsRune(selectedSchema, '\x00') {
		return []string{selectedSchema}
	}
	return []string{selectedSchema, "public"}
}

// searchPathPrefix starts the statement WithSearchPath puts before a DML
// statement.
const searchPathPrefix = "SET LOCAL search_path TO "

// WithSearchPath returns statement preceded by a SET LOCAL search_path
// statement for schemas.
func WithSearchPath(statement string, schemas []string) string {
	quoted := make([]string, len(schemas))
	for i, schema := range schemas {
		quoted[i] = `"` + strings.ReplaceAll(schema, `"`, `""`) + `"`
	}
	return searchPathPrefix + strings.Join(quoted, ", ") + ";\n" + statement
}

// SplitSearchPath returns the SET LOCAL search_path statement that
// WithSearchPath put before statement, or "" when there is none, and the
// statement after it.
func SplitSearchPath(statement string) (string, string) {
	if !strings.HasPrefix(statement, searchPathPrefix) {
		return "", statement
	}
	setup, rest, ok := strings.Cut(statement, ";\n")
	if !ok {
		return "", statement
	}
	return setup, rest
}
