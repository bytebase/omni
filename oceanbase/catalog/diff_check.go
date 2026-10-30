package catalog

import (
	"sort"
	"strconv"
)

// MySQL SDL diff — CHECK constraints (8.0.16+ only).
//
// This is the MySQL analog of PG's CHECK differ. It compares the CHECK-constraint set of two
// versions of a table and reports the per-check changes via TableDiffEntry.Checks. It is wired
// into compareTable (diff_table.go); tableSubdiffsChanged folds a non-empty Checks slice into the
// "is this table modified?" decision, so this node only fills its own slice.
//
// VERSION GATE (the headline correctness property). CHECK constraints are 8.0.16+ only. On MySQL
// 5.7 a CHECK clause is PARSED AND SILENTLY DROPPED — the engine never stores it, so SHOW CREATE
// echoes nothing (normalize.go CheckSupported, entry check-constraint). The omni loader, however,
// stores a ConCheck constraint UNCONDITIONALLY (tablecmds.go appends it regardless of version), so
// a 5.7-targeted catalog loaded from user DDL DOES carry ConCheck entries while its engine readback
// carries none. Comparing them would phantom-diff forever. So when the Normalizer targets a version
// that does not represent CHECK in its stored form (5.7), this differ returns nil: no check exists
// in the stored form, so there is nothing to diff. The gate is CheckSupported(), the single owner
// of the version decision — this node never re-derives "is 8.0".
//
// Identity is the constraint name, lower-cased. MySQL auto-names an unnamed CHECK `<table>_chk_<N>`
// at parse time, and the omni loader reproduces that same auto-name (tablecmds.go nextCheckNumber),
// so by the time a check reaches the catalog it is ALWAYS named — both user-named and engine/loader
// auto-named — exactly like the index node keys on the (auto-)name. Equality is decided by a
// canonical key (canonicalCheck) folding the two stored-form-significant attributes:
//   - the CHECK expression, routed through normalize-core's CanonicalCheckExpr (backtick idents,
//     operator spacing, function casing, the 8.0 `_latin1` string introducer — all canonicalized
//     there, never re-implemented here); and
//   - the ENFORCED / NOT ENFORCED state (NotEnforced), a version-independent boolean.
//
// An expression whose surface form differs from its SHOW CREATE readback but whose stored form is
// identical (e.g. `a>0` vs the readback's `(`a` > 0)`) therefore produces no phantom diff.
func diffChecks(from, to *Table, n *Normalizer) []CheckDiffEntry {
	// Version gate: CHECK is unrepresentable on 5.7 (parsed-and-dropped), so there is no stored
	// check to diff. Returning nil keeps a 5.7 self/round-trip diff empty even when the loaded
	// catalog carries ConCheck entries the engine would have discarded.
	if !n.CheckSupported() {
		return nil
	}

	fromMap := checkMap(from)
	toMap := checkMap(to)

	var result []CheckDiffEntry

	// Dropped: in from but not in to.
	for name, fromChk := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, CheckDiffEntry{
				Action: DiffDrop,
				Name:   fromChk.Name,
				From:   fromChk,
			})
		}
	}

	// Added or modified: in to.
	for name, toChk := range toMap {
		fromChk, ok := fromMap[name]
		if !ok {
			result = append(result, CheckDiffEntry{
				Action: DiffAdd,
				Name:   toChk.Name,
				To:     toChk,
			})
			continue
		}
		if checksChanged(fromChk, toChk, n) {
			result = append(result, CheckDiffEntry{
				Action: DiffModify,
				Name:   toChk.Name,
				From:   fromChk,
				To:     toChk,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if a, b := toLower(result[i].Name), toLower(result[j].Name); a != b {
			return a < b
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// checkMap indexes a table's CHECK constraints by lower-cased name. Only ConCheck constraints
// participate — PRIMARY KEY / UNIQUE are owned by the index/constraint nodes and FOREIGN KEY by the
// FK node (the same constraint-kind split show.go uses). A check with an empty name is skipped
// defensively (the loader always names them, so this never fires in practice).
func checkMap(t *Table) map[string]*Constraint {
	m := make(map[string]*Constraint)
	if t == nil {
		return m
	}
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConCheck {
			continue
		}
		name := toLower(con.Name)
		if name == "" {
			continue
		}
		m[name] = con
	}
	return m
}

// checksChanged reports whether two same-name CHECK constraints differ, comparing their canonical
// keys. canonicalCheck folds every stored-form-significant attribute into one key, so this differ
// never re-implements a per-attribute comparison (mirroring indexesChanged).
func checksChanged(a, b *Constraint, n *Normalizer) bool {
	return canonicalCheck(a, n) != canonicalCheck(b, n)
}

// canonicalCheck returns a single stable comparison key for a CHECK constraint, folding the
// attributes that survive into MySQL's stored form. The expression is routed through normalize-core
// (CanonicalCheckExpr) for all canonicalization (backtick idents, operator spacing, function
// casing, the 8.0 `_latin1` introducer); the ENFORCED/NOT ENFORCED state is a version-independent
// boolean compared directly. The constraint NAME is the identity key (handled by the caller's map),
// not part of this content key.
func canonicalCheck(con *Constraint, n *Normalizer) string {
	return encodeKeyFields(
		"expr", n.CanonicalCheckExpr(con.CheckExpr),
		"enforced", strconv.FormatBool(!con.NotEnforced),
	)
}
