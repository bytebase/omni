package redshift

import redshiftparser "github.com/bytebase/omni/redshift/parser"

// LSPPosition is a zero-based UTF-16 LSP position.
type LSPPosition struct {
	Line      int
	Character int
}

// LSPRange is a zero-based UTF-16 LSP range.
type LSPRange struct {
	Start LSPPosition
	End   LSPPosition
}

// Diagnostic reports a parser diagnostic over an LSP range.
type Diagnostic struct {
	Message string
	Range   LSPRange
}

// Diagnose returns syntax diagnostics for Redshift SQL.
func Diagnose(sql string) []Diagnostic {
	if _, err := Parse(sql); err != nil {
		pos := 0
		if pe, ok := err.(*redshiftparser.ParseError); ok && pe.Position >= 0 {
			pos = pe.Position
		}
		start := offsetToLSPPosition(sql, pos)
		end := start
		end.Character++
		return []Diagnostic{{
			Message: err.Error(),
			Range:   LSPRange{Start: start, End: end},
		}}
	}
	return nil
}
