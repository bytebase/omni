package analysis

import (
	"strings"
	"unicode"
)

// Position is a place in a script as an LSP editor counts it: a 0-based line
// and a 0-based character offset in UTF-16 code units.
type Position struct {
	Line      uint32
	Character uint32
}

// Range is the half-open range [Start, End) of a statement.
type Range struct {
	Start Position
	End   Position
}

// StatementRanges returns the range of each statement in a script, for
// editors that run the statement under the cursor. It splits at every ';',
// including one inside a string, comment or dollar-quoted body. A range
// starts at the statement's first non-whitespace character and includes its
// ';'; when a newline follows the ';' it ends at the start of the next line.
// The last statement's range ends at the start of the line after it. A script
// with no statements returns nil.
func StatementRanges(statement string) []Range {
	s := strings.TrimRightFunc(statement, unicode.IsSpace)
	// Assuming the content is UTF-8 encoded.
	statements := strings.Split(s, ";")

	var ranges []Range
	// 0-based UTF-16 encoded line and character.
	line, character := 0, 0

	for i, statement := range statements {
		// Trim left space to provide more accurate range.
		statement = strings.TrimLeftFunc(statement, func(r rune) bool {
			if !unicode.IsSpace(r) {
				return false
			}
			if r == '\n' {
				line++
				character = 0
			} else {
				// Check rune utf16 length by BMP.
				if r <= 0xFFFF {
					character++
				} else {
					character += 2
				}
			}
			return true
		})

		// If the statement is empty, skip it.
		if statement == "" {
			continue
		}

		begin := Position{Line: uint32(line), Character: uint32(character)}
		for _, r := range statement {
			if r == '\n' {
				line++
				character = 0
			} else {
				// Check rune utf16 length by BMP.
				if r <= 0xFFFF {
					character++
				} else {
					character += 2
				}
			}
		}

		endLine, endCharacter := line, character
		// End is exclusive, so we check the next byte.
		if i == len(statements)-1 {
			// End of the content.
			endLine++
			endCharacter = 0
		} else {
			// Next byte is ';', include it.
			character++
			endCharacter++
			if nextStatement := statements[i+1]; len(nextStatement) > 0 && nextStatement[0] == '\n' {
				endLine++
				endCharacter = 0
			} else {
				endCharacter++
			}
		}
		end := Position{Line: uint32(endLine), Character: uint32(endCharacter)}
		ranges = append(ranges, Range{Start: begin, End: end})
	}
	return ranges
}
