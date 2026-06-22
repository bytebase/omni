package cassandra

import (
	"github.com/bytebase/omni/cassandra/parser"
)

// Segment represents a single SQL segment from splitting at top-level semicolons.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
	Empty     bool
}

// Split splits a CQL input into segments at top-level semicolons.
func Split(sql string) []Segment {
	internal := parser.Split(sql)
	result := make([]Segment, len(internal))
	for i, seg := range internal {
		result[i] = Segment{
			Text:      seg.Text,
			ByteStart: seg.ByteStart,
			ByteEnd:   seg.ByteEnd,
			Empty:     seg.Empty,
		}
	}
	return result
}
