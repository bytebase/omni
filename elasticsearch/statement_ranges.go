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
// Each is review.Index's position made 0-based, with a character outside the
// BMP counted twice.
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
		ranges = append(ranges, LSPRange{
			Start: positions.position(r.StartOffset),
			End:   positions.position(r.EndOffset),
		})
	}
	return ranges, nil
}

// lspPositioner converts byte offsets into LSPPositions. It counts newlines
// to determine the line number and accumulates UTF-16 code units for the
// character column. Characters outside the BMP (code points above U+FFFF,
// encoded as 4-byte UTF-8 sequences) consume two UTF-16 code units and
// therefore advance the character counter by 2. A byte that is not UTF-8
// advances it by 1, as review.Index counts it as one code point and an editor
// shows it as U+FFFD; the scan used to stop there, or at U+FFFD itself, and
// GetStatementRanges dropped every range after it. It resumes from the offset
// it was last asked for instead of scanning from the start for every offset,
// which made the ranges of a script of n requests take time quadratic in n.
type lspPositioner struct {
	bs []byte
	// i is the byte offset at describes.
	i  int
	at LSPPosition
	// last is the offset last asked for.
	last int
}

// position returns the position at byteOffset.
func (p *lspPositioner) position(byteOffset int) LSPPosition {
	if byteOffset < p.last {
		// Offsets come in order; start over for one that does not.
		*p = lspPositioner{bs: p.bs}
	}
	p.last = byteOffset
	for p.i < byteOffset && p.i < len(p.bs) {
		if p.bs[p.i] == '\n' {
			p.at.Line++
			p.at.Character = 0
			p.i++
			continue
		}
		r, size := utf8.DecodeRune(p.bs[p.i:])
		p.at.Character++
		if r > 0xFFFF {
			// Outside the BMP — requires a surrogate pair in UTF-16.
			p.at.Character++
		}
		p.i += size
	}
	return p.at
}
