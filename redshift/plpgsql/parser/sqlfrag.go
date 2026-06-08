package parser

import (
	"strings"

	pgparser "github.com/bytebase/omni/redshift/parser"
)

// collectUntil scans tokens tracking paren depth and string boundaries
// until one of the terminator keywords is found at depth 0.
// Returns the collected text (trimmed) from the source.
// The terminator keyword is NOT consumed.
func (p *Parser) collectUntil(terminators ...string) (string, error) {
	start := p.pos()
	depth := 0

	for !p.isEOF() {
		// Track parenthesis depth
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}

		// Only check terminators at depth 0
		if depth == 0 {
			for _, term := range terminators {
				if p.isKeyword(term) {
					text := strings.TrimSpace(p.source[start:p.cur.Loc])
					return text, nil
				}
			}
			// Also check for semicolon as implicit terminator at depth 0
			if p.cur.Type == ';' {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				return text, nil
			}
		}

		p.advance()
	}

	// Hit EOF — return what we have
	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}

// collectUntilSemicolon scans tokens tracking paren depth until
// a semicolon is found at depth 0.
// Returns the collected text (trimmed). The semicolon is NOT consumed.
func (p *Parser) collectUntilSemicolon() (string, error) {
	start := p.pos()
	depth := 0

	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}
		if depth == 0 && p.cur.Type == ';' {
			text := strings.TrimSpace(p.source[start:p.cur.Loc])
			return text, nil
		}
		p.advance()
	}

	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}

// Ensure pgparser import is used (Token type reference).
var _ = pgparser.Token{}
