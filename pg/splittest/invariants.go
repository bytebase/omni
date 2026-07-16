package splittest

import (
	"fmt"

	"github.com/bytebase/omni/pg"
)

// CheckInvariants runs the truth-free properties against one input.
// These hold for ANY input, valid SQL or not, so they can absorb
// unlimited volume (corpus files, byte mutations, generated scripts):
//
//	I1 lossless: concatenating all segment texts reproduces the input
//	I2 ranges:   each segment's [ByteStart,ByteEnd) slices to its Text,
//	             ranges are adjacent and cover the whole input
//	I3 re-split: splitting a segment's text yields that segment back
//	             (idempotence)
//
// It returns a descriptive error for the first violated property.
func CheckInvariants(sql string) error {
	segs := pg.Split(sql)

	// I1 + I2.
	pos := 0
	var concat []byte
	for i, s := range segs {
		if s.ByteStart != pos {
			return fmt.Errorf("I2: segment %d starts at %d, want %d", i, s.ByteStart, pos)
		}
		if s.ByteEnd < s.ByteStart || s.ByteEnd > len(sql) {
			return fmt.Errorf("I2: segment %d has range [%d,%d) outside input len %d", i, s.ByteStart, s.ByteEnd, len(sql))
		}
		if sql[s.ByteStart:s.ByteEnd] != s.Text {
			return fmt.Errorf("I2: segment %d Text %q != sql[%d:%d] %q", i, s.Text, s.ByteStart, s.ByteEnd, sql[s.ByteStart:s.ByteEnd])
		}
		concat = append(concat, s.Text...)
		pos = s.ByteEnd
	}
	if len(sql) > 0 && pos != len(sql) {
		return fmt.Errorf("I2: segments cover [0,%d), input len %d", pos, len(sql))
	}
	if string(concat) != sql {
		return fmt.Errorf("I1: concat(split(sql)) != sql")
	}

	// I3: re-splitting one segment must not find new boundaries.
	for i, s := range segs {
		re := pg.Split(s.Text)
		if len(re) != 1 {
			return fmt.Errorf("I3: segment %d re-splits into %d segments: %q", i, len(re), s.Text)
		}
		if re[0].Text != s.Text {
			return fmt.Errorf("I3: segment %d re-split text %q != %q", i, re[0].Text, s.Text)
		}
	}
	return nil
}

// CheckScript verifies a generated script against its constructive
// expectation: exact segment texts in order, plus all invariants.
func CheckScript(s Script) error {
	if err := CheckInvariants(s.SQL); err != nil {
		return err
	}
	segs := pg.Split(s.SQL)
	if len(segs) != len(s.Want) {
		return fmt.Errorf("constructive: got %d segments, want %d\ninput: %q", len(segs), len(s.Want), s.SQL)
	}
	for i := range segs {
		if segs[i].Text != s.Want[i] {
			return fmt.Errorf("constructive: segment %d = %q, want %q\ninput: %q", i, segs[i].Text, s.Want[i], s.SQL)
		}
	}
	return nil
}
