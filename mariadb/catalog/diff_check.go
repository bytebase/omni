package catalog

import (
	"sort"
	"strconv"
)

// diffChecks compares the table-level CHECK constraints of two versions of a
// table. Identity is the lower-cased constraint name; two same-name checks are
// equal when their expression and ENFORCED state match (canonicalCheck). A
// column's own CHECK is part of the column's definition and is compared with the
// column, not here.
func diffChecks(from, to *Table, n *Normalizer) []CheckDiffEntry {
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

// checkMap indexes the CHECK constraints a table writes with its constraints by
// lower-cased name.
func checkMap(t *Table) map[string]*Constraint {
	m := make(map[string]*Constraint)
	if t == nil {
		return m
	}
	// A column's check, a JSON column's json_valid among them, is part of the
	// column's definition where SHOW CREATE TABLE writes it on the column (see
	// CanonicalColumn). A generated column's, which it writes with the table's
	// constraints, is compared here.
	onColumn := columnChecks(t)
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConCheck || (isColumnCheck(con) && onColumn[toLower(con.Column)] == con) {
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

// checksChanged reports whether two same-name checks differ.
func checksChanged(a, b *Constraint, n *Normalizer) bool {
	return canonicalCheck(a, n) != canonicalCheck(b, n)
}

// canonicalCheck is the comparison key of a check: its expression and its
// ENFORCED state. The name is the identity key, not part of this key.
func canonicalCheck(con *Constraint, n *Normalizer) string {
	return encodeKeyFields(
		"expr", n.CanonicalCheckExpr(con.CheckExpr),
		"enforced", strconv.FormatBool(!con.NotEnforced),
	)
}
