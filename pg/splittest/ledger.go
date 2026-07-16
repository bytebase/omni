package splittest

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// The enumeration-layer coverage ledger (framework design §2.1).
//
// SPLIT_COVERAGE.json is the resident, human-auditable account of what
// the deterministic layer tests per construct class. Every entry is
// executed by TestLedgerEntries (exact segmentation + invariants), and
// TestLedgerLint enforces the per-class coverage quota so the ledger
// cannot silently rot: a new construct class without registered forms
// fails CI, mirroring the PAREN_AUDIT mechanism.
//
//go:embed SPLIT_COVERAGE.json
var ledgerJSON []byte

// LedgerEntry is one deterministic, human-reviewed case.
type LedgerEntry struct {
	// ID is stable and unique, referenced from fix PRs and defer
	// artifacts (anchor format: ledger:<id>).
	ID string `json:"id"`
	// Class is the construct class (C1..C18 from the design taxonomy;
	// C19 real-corpus coverage lives in dedicated file-level tests).
	Class string `json:"class"`
	// Form positions the case in the class's boundary matrix:
	// "in" (界内), "boundary" (边界), "out" (界外/negative).
	Form string `json:"form"`
	// SQL is the input script.
	SQL string `json:"sql"`
	// Want holds the exact expected segment texts, in order.
	Want []string `json:"want"`
	// TruthSource records where the expectation comes from:
	// construct | invariant | spec:psqlscan | spec:scan.l | container.
	TruthSource string `json:"truth_source"`
	// Notes carries provenance (audit defect ids, PR numbers,
	// documented divergences).
	Notes string `json:"notes,omitempty"`
}

// LoadLedger parses the embedded ledger.
func LoadLedger() ([]LedgerEntry, error) {
	var entries []LedgerEntry
	if err := json.Unmarshal(ledgerJSON, &entries); err != nil {
		return nil, fmt.Errorf("SPLIT_COVERAGE.json: %w", err)
	}
	return entries, nil
}

// LedgerClasses is the construct-class taxonomy the lint enforces.
// Classes marked minOut=0 have no meaningful negative form of their
// own (their negatives live in adjacent classes).
var LedgerClasses = []struct {
	Class                string
	MinIn, MinBd, MinOut int
}{
	{"C1-plain-string", 1, 1, 1},
	{"C2-escape-string", 1, 1, 1},
	{"C3-unicode-string", 1, 1, 0},
	{"C4-bit-string", 1, 1, 0},
	{"C5-dollar-quote", 1, 1, 1},
	{"C6-ident-dollar", 1, 1, 1},
	{"C7-line-comment", 1, 1, 0},
	{"C8-block-comment", 1, 1, 1},
	{"C9-quoted-ident", 1, 1, 0},
	{"C10-paren-depth", 1, 1, 1},
	{"C11-begin-atomic", 1, 1, 1},
	{"C12-copy-stdin", 1, 1, 1},
	{"C13-semicolon-variants", 1, 1, 0},
	{"C14-encoding", 1, 1, 0},
	{"C15-resilience", 1, 1, 0},
	{"C16-boundary-misc", 1, 1, 0},
	{"C17-keyword-boundary", 1, 1, 0},
	{"C18-metacommand", 1, 1, 1},
}
