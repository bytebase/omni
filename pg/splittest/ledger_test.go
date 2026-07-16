package splittest

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/pg"
)

// TestLedgerEntries executes every SPLIT_COVERAGE.json entry: exact
// segmentation match plus the truth-free invariants. The ledger is the
// enumeration layer's executable account — an entry that stops passing
// is a splitter regression (or an intentional contract change, which
// must update the ledger in the same PR).
func TestLedgerEntries(t *testing.T) {
	entries, err := LoadLedger()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		t.Run(e.ID, func(t *testing.T) {
			if err := CheckScript(Script{SQL: e.SQL, Want: e.Want}); err != nil {
				t.Errorf("%v (class=%s form=%s truth=%s notes=%s)", err, e.Class, e.Form, e.TruthSource, e.Notes)
			}
		})
	}
}

// TestLedgerLint enforces the ledger's structural rules so coverage
// cannot silently rot (PAREN_AUDIT mechanism):
//
//	L1 every taxonomy class meets its in/boundary/out quota
//	L2 ids are unique, classes and forms are from the taxonomy
//	L3 truth_source is one of the accepted sources
//	L4 want concatenation reproduces sql (entries are lossless by
//	   construction, catching hand-editing slips at lint time)
func TestLedgerLint(t *testing.T) {
	entries, err := LoadLedger()
	if err != nil {
		t.Fatal(err)
	}

	validSource := map[string]bool{
		"construct": true, "invariant": true,
		"spec:psqlscan": true, "spec:scan.l": true, "container": true,
	}
	classIdx := map[string]int{}
	for i, c := range LedgerClasses {
		classIdx[c.Class] = i
	}

	seen := map[string]bool{}
	type count struct{ in, bd, out int }
	counts := map[string]*count{}

	for _, e := range entries {
		if seen[e.ID] {
			t.Errorf("L2: duplicate id %q", e.ID)
		}
		seen[e.ID] = true
		if _, ok := classIdx[e.Class]; !ok {
			t.Errorf("L2: entry %s has unknown class %q", e.ID, e.Class)
			continue
		}
		if counts[e.Class] == nil {
			counts[e.Class] = &count{}
		}
		switch e.Form {
		case "in":
			counts[e.Class].in++
		case "boundary":
			counts[e.Class].bd++
		case "out":
			counts[e.Class].out++
		default:
			t.Errorf("L2: entry %s has unknown form %q", e.ID, e.Form)
		}
		if !validSource[e.TruthSource] {
			t.Errorf("L3: entry %s has invalid truth_source %q", e.ID, e.TruthSource)
		}
		if got := strings.Join(e.Want, ""); got != e.SQL {
			t.Errorf("L4: entry %s: concat(want) != sql\n concat: %q\n    sql: %q", e.ID, got, e.SQL)
		}
	}

	for _, c := range LedgerClasses {
		got := counts[c.Class]
		if got == nil {
			t.Errorf("L1: class %s has no ledger entries", c.Class)
			continue
		}
		if got.in < c.MinIn || got.bd < c.MinBd || got.out < c.MinOut {
			t.Errorf("L1: class %s below quota: have in=%d bd=%d out=%d, need in>=%d bd>=%d out>=%d",
				c.Class, got.in, got.bd, got.out, c.MinIn, c.MinBd, c.MinOut)
		}
	}
}

// TestLedgerCoversKnownBetter pins that every KnownBetterThanPsql input
// is (a) registered in the coverage ledger by exact SQL, so a new
// whitelist entry cannot be added without a documented ledger row, and
// (b) still splits to a single non-empty statement. Without the ledger
// lookup this would just duplicate TestKnownBetterThanPsqlSplitAsSingle-
// Statements and the advertised anti-rot guarantee would be hollow.
func TestLedgerCoversKnownBetter(t *testing.T) {
	entries, err := LoadLedger()
	if err != nil {
		t.Fatal(err)
	}
	ledgerSQL := make(map[string]bool, len(entries))
	for _, e := range entries {
		ledgerSQL[e.SQL] = true
	}

	for i, sql := range KnownBetterThanPsql {
		if !ledgerSQL[sql] {
			t.Errorf("known-better[%d] is not registered in SPLIT_COVERAGE.json — add a ledger entry: %q", i, sql)
		}
		segs := pg.Split(sql)
		nonEmpty := 0
		for _, s := range segs {
			if !s.Empty() {
				nonEmpty++
			}
		}
		if nonEmpty != 1 {
			t.Errorf("known-better[%d]: %d non-empty segments, want 1: %q", i, nonEmpty, sql)
		}
	}
}
