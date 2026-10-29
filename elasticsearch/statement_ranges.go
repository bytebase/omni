package elasticsearch

import "unicode/utf8"

// LSPPosition is a 0-based line/character position using UTF-16 code units,
// following the Language Server Protocol convention.
type LSPPosition struct {
	Line      uint32 `yaml:"line"`
	Character uint32 `yaml:"character"`
}

// LSPRange is a half-open [Start, End) range of LSP positions.
type LSPRange struct {
	Start LSPPosition `yaml:"start"`
	End   LSPPosition `yaml:"end"`
}

// GetStatementRanges returns one LSPRange per parsed request in statement.
// Positions use 0-based lines and UTF-16 code-unit columns, matching the LSP
// specification (BMP characters count as 1, non-BMP characters count as 2).
func GetStatementRanges(statement string) ([]LSPRange, error) {
	parseResult, _, _ := parseREST(statement, false)
	if parseResult == nil {
		return []LSPRange{}, nil
	}
	positions := &lspPositioner{bs: []byte(statement)}
	var ranges []LSPRange
	for _, r := range parseResult.Requests {
		if r == nil {
			continue
		}
		if r.EndOffset <= r.StartOffset {
			continue
		}

		startPosition := positions.position(r.StartOffset)
		endPosition := positions.position(r.EndOffset)
		if startPosition == nil || endPosition == nil {
			continue
		}
		ranges = append(ranges, LSPRange{
			Start: *startPosition,
			End:   *endPosition,
		})
	}
	return ranges, nil
}

// lspPositioner converts byte offsets into LSPPositions. It counts newlines
// to determine the line number and accumulates UTF-16 code units for the
// character column. Characters outside the BMP (code points above U+FFFF,
// encoded as 4-byte UTF-8 sequences) consume two UTF-16 code units and
// therefore advance the character counter by 2. It resumes from the offset it
// was last asked for instead of scanning from the start for every offset,
// which made the ranges of a script of n requests take time quadratic in n.
type lspPositioner struct {
	bs []byte
	// i is the byte offset at describes.
	i  int
	at LSPPosition
	// invalid is set once the scan reaches U+FFFD or a byte that is not
	// UTF-8, at offset i; no position past it is reported.
	invalid bool
	// last is the offset last asked for.
	last int
}

// position returns the position at byteOffset, or nil when the text before
// it holds U+FFFD or a byte that is not UTF-8.
func (p *lspPositioner) position(byteOffset int) *LSPPosition {
	if byteOffset < p.last {
		// Offsets come in order; start over for one that does not.
		*p = lspPositioner{bs: p.bs}
	}
	p.last = byteOffset
	for !p.invalid && p.i < byteOffset && p.i < len(p.bs) {
		if p.bs[p.i] == '\n' {
			p.at.Line++
			p.at.Character = 0
			p.i++
			continue
		}
		r, size := utf8.DecodeRune(p.bs[p.i:])
		if r == utf8.RuneError {
			p.invalid = true
			break
		}
		p.at.Character++
		if r > 0xFFFF {
			// Outside the BMP — requires a surrogate pair in UTF-16.
			p.at.Character++
		}
		p.i += size
	}
	if p.invalid && p.i < byteOffset {
		return nil
	}
	position := p.at
	return &position
}
