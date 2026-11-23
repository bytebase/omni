package review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// A sequence's bounds are one key of its own under its path, read only when
// Present, with "start min max" as its kind. The walk knows them from the
// sync (information_schema.SEQUENCES) and follows ALTER SEQUENCE; a
// sequence it creates or drops has bounds it leaves Unknown.
const (
	spaceSequenceBounds = "sequence bounds"
	sequenceBoundsName  = "are"
)

func sequenceBoundsKey(db, seq string) walk.Key {
	return walk.Key{Space: spaceSequenceBounds, Parent: tablePath(db, seq), Name: sequenceBoundsName}
}

type sequenceBounds struct{ start, min, max int64 }

func (s sequenceBounds) kind() string { return fmt.Sprintf("%d %d %d", s.start, s.min, s.max) }

func parseSequenceBounds(kind string) (sequenceBounds, bool) {
	f := strings.Fields(kind)
	if len(f) != 3 {
		return sequenceBounds{}, false
	}
	var s sequenceBounds
	for i, p := range []*int64{&s.start, &s.min, &s.max} {
		v, err := strconv.ParseInt(f[i], 10, 64)
		if err != nil {
			return sequenceBounds{}, false
		}
		*p = v
	}
	return s, true
}

// syncedSequenceBounds reads a synced sequence's bounds, false when the
// sync gives one the walk cannot read.
func syncedSequenceBounds(q *metadata.SequenceMetadata) (string, bool) {
	var s sequenceBounds
	for _, p := range []struct {
		text string
		v    *int64
	}{{q.GetStart(), &s.start}, {q.GetMinValue(), &s.min}, {q.GetMaxValue(), &s.max}} {
		v, err := strconv.ParseInt(strings.TrimSpace(p.text), 10, 64)
		if err != nil {
			return "", false
		}
		*p.v = v
	}
	return s.kind(), true
}

// sequenceBoundsFail fails an ALTER SEQUENCE that leaves the sequence's
// start outside its bounds, or its maximum no greater than its minimum
// (4136, verified on v8.5.5: MAXVALUE 1 on a sequence of the default
// bounds, MINVALUE 1 on one that starts at -5). It records the bounds the
// statement leaves.
func (b *builder) sequenceBoundsFail(r tableRef, opts []*ast.SequenceOption) bool {
	k := sequenceBoundsKey(r.db, r.name)
	changes, follows := false, true
	for _, o := range opts {
		switch o.Type {
		case ast.SequenceOptionStartWith, ast.SequenceOptionMinValue, ast.SequenceOptionMaxValue:
			changes = true
		case ast.SequenceOptionNoMinValue, ast.SequenceOptionNoMaxValue,
			ast.SequenceOptionRestart, ast.SequenceOptionRestartWith:
			// Their bounds depend on the increment, or the walk does not
			// follow what they do to the start.
			changes, follows = true, false
		}
	}
	if !changes {
		return false
	}
	f, kind := b.l.Lookup(k)
	cur, ok := parseSequenceBounds(kind)
	if f != walk.Present || !ok || !follows {
		b.Set(k, walk.Unknown, "")
		return false
	}
	for _, o := range opts {
		switch o.Type {
		case ast.SequenceOptionStartWith:
			cur.start = o.Value
		case ast.SequenceOptionMinValue:
			cur.min = o.Value
		case ast.SequenceOptionMaxValue:
			cur.max = o.Value
		}
	}
	if cur.max <= cur.min || cur.start < cur.min || cur.start > cur.max {
		b.fail(b.anchor(), nil, review.UnsupportedForm, "Sequence "+review.Code(r.written)+" has conflicting options.")
		return true
	}
	b.Set(k, walk.Present, cur.kind())
	return false
}
