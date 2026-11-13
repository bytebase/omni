package catalog

import (
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// MySQL SDL diff — CHECK constraints (8.0.16+ only).
//
// This is the MySQL analog of PG's CHECK differ. It compares the CHECK-constraint set of two
// versions of a table and reports the per-check changes via TableDiffEntry.Checks. It is wired
// into compareTable (diff_table.go); tableSubdiffsChanged folds a non-empty Checks slice into the
// "is this table modified?" decision, so this node only fills its own slice.
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
//
// One case is not decided by the name. OceanBase names a CHECK written without a name
// `<table>_OBCHECK_<time>`, which this catalog's loader cannot reproduce: it names the same
// constraint `<table>_chk_<N>`. Two CHECKs that differ only in such names are one constraint
// (forgetMadeUpRenames), or every unnamed CHECK would be dropped and added again on each release.
func diffChecks(from, to *Table, n *Normalizer) []CheckDiffEntry {
	fromMap := checkMap(from)
	toMap := checkMap(to)
	forgetMadeUpRenames(fromMap, toMap, madeUpCheckName, func(con *Constraint) string { return canonicalCheck(con, n) })

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

// madeUpCheckName reports whether a CHECK's name is one nobody wrote: the name this catalog's
// loader gives a CHECK written without one, `<table>_chk_<N>` as MySQL does, or the name
// OceanBase gives it, `<table>_OBCHECK_<time in microseconds>`.
func madeUpCheckName(con *Constraint) bool {
	return madeUpName(con, "_chk_", "_obcheck_")
}

// madeUpName reports whether con's name is its table's name, loaderInfix and digits, or any
// prefix, engineInfix and digits. The engine's own names are told by the infix alone: nobody
// writes one, and a renamed table keeps the names made under its old name.
func madeUpName(con *Constraint, loaderInfix, engineInfix string) bool {
	name := toLower(con.Name)
	if con.Table != nil {
		if n, ok := strings.CutPrefix(name, toLower(con.Table.Name)+loaderInfix); ok && n != "" && isAllDigits(n) {
			return true
		}
	}
	if i := strings.LastIndex(name, engineInfix); i >= 0 {
		n := name[i+len(engineInfix):]
		return n != "" && isAllDigits(n)
	}
	return false
}

// forgetMadeUpRenames deletes from both maps every pair of constraints that is one constraint
// under two names nobody wrote: a name only one side has, made up on both sides, over the same
// definition. A constraint both sides name alike is left for the caller to compare, and so is
// one whose name somebody wrote: giving an unnamed constraint a name is a change.
//
// Definitions can repeat, so the pairing counts: two unnamed copies on one side and one on the
// other leave one over.
func forgetMadeUpRenames(from, to map[string]*Constraint, madeUp func(*Constraint) bool, definition func(*Constraint) string) {
	unpaired := map[string][]string{} // definition to the names `to` alone has for it
	for _, name := range slices.Sorted(maps.Keys(to)) {
		if _, both := from[name]; !both && madeUp(to[name]) {
			def := definition(to[name])
			unpaired[def] = append(unpaired[def], name)
		}
	}
	for _, name := range slices.Sorted(maps.Keys(from)) {
		if _, both := to[name]; both || !madeUp(from[name]) {
			continue
		}
		def := definition(from[name])
		if names := unpaired[def]; len(names) > 0 {
			delete(from, name)
			delete(to, names[0])
			unpaired[def] = names[1:]
		}
	}
}
