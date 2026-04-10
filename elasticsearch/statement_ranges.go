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
	parseResult, _ := ParseElasticsearchREST(statement)
	if parseResult == nil {
		return []LSPRange{}, nil
	}
	bs := []byte(statement)
	var ranges []LSPRange
	for _, r := range parseResult.Requests {
		if r == nil {
			continue
		}
		if r.EndOffset <= r.StartOffset {
			continue
		}

		startPosition := getPositionByByteOffset(r.StartOffset, bs)
		endPosition := getPositionByByteOffset(r.EndOffset, bs)
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

// getPositionByByteOffset converts a byte offset into an LSPPosition.
// It counts newlines to determine the line number and accumulates UTF-16
// code units for the character column. Characters outside the BMP (code points
// above U+FFFF, encoded as 4-byte UTF-8 sequences) consume two UTF-16 code
// units and therefore advance the character counter by 2.
func getPositionByByteOffset(byteOffset int, bs []byte) *LSPPosition {
	var position LSPPosition
	for i := 0; ; {
		if i >= byteOffset || i > len(bs) {
			break
		}
		if bs[i] == '\n' {
			position.Line++
			position.Character = 0
			i++
			continue
		}
		r, size := utf8.DecodeRune(bs[i:])
		if r == utf8.RuneError {
			return nil
		}
		position.Character++
		if r > 0xFFFF {
			// Outside the BMP — requires a surrogate pair in UTF-16.
			position.Character++
		}
		i += size
	}
	return &position
}
