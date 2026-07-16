package splittest

import (
	"testing"

	pg "github.com/bytebase/omni/pg"
)

// TestKnownBetterThanPsqlSplitAsSingleStatements pins that every
// known-better input splits as exactly one non-empty segment: these are
// the shapes where omni intentionally exceeds psql's client scanner.
func TestKnownBetterThanPsqlSplitAsSingleStatements(t *testing.T) {
	for _, sql := range KnownBetterThanPsql {
		segs := pg.Split(sql)
		n := 0
		for _, s := range segs {
			if !s.Empty() {
				n++
			}
		}
		if n != 1 {
			t.Errorf("expected 1 statement, got %d for %q", n, sql)
		}
	}
}
