package completion

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/pg/yacc"
)

// QuoteIdentifier returns a properly quoted PostgreSQL identifier.
// It adds double quotes when the identifier:
// - is a reserved keyword
// - contains uppercase letters
// - starts with a digit
// - contains characters other than lowercase letters, digits, and underscores
// - is empty
func QuoteIdentifier(name string) string {
	if name == "" {
		return `""`
	}
	if needsQuoting(name) {
		// Double any existing double quotes inside the name
		escaped := strings.ReplaceAll(name, `"`, `""`)
		return `"` + escaped + `"`
	}
	return name
}

// needsQuoting checks whether a PostgreSQL identifier requires double-quoting.
func needsQuoting(name string) bool {
	// Check if it's a reserved keyword
	if kw := yacc.LookupKeyword(name); kw != nil && kw.Category == yacc.ReservedKeyword {
		return true
	}

	for i, r := range name {
		// First character must be a letter or underscore
		if i == 0 {
			if r >= '0' && r <= '9' {
				return true
			}
		}
		// Only lowercase letters, digits, and underscore are safe unquoted
		if unicode.IsUpper(r) {
			return true
		}
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return true
		}
	}
	return false
}
