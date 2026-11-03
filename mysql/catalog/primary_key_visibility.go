package catalog

import "sort"

// declaredUniqueKey is a UNIQUE key of a CREATE TABLE with the byte offset its
// definition starts at, -1 when unknown.
type declaredUniqueKey struct {
	idx *Index
	pos int
}

// checkImplicitPrimaryKeyVisible refuses a table without a primary key whose
// implicit one is invisible. MySQL 8.0 takes the first UNIQUE key on NOT NULL
// parts, in the order the statement declares the keys, for the primary key and
// answers 3522 when it is invisible:
//
//	CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)                  -- 3522
//	CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE (b), UNIQUE KEY u (a) INVISIBLE) -- 3522
//	CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE) -- ok
//
// A prefix part counts. A key part that is an expression is NOT NULL when the
// expression is, (a + 1) on a NOT NULL a or (a IS NULL) on any a, which the
// catalog does not work out; a key with one ends the search without a verdict,
// so the catalog never refuses a table MySQL accepts.
//
// A column definition's UNIQUE (or SERIAL) declares a key where the column
// stands, and such a key is always visible. When an offset is unknown the keys
// stay in the catalog's order, column keys first, which can only miss a refusal.
func checkImplicitPrimaryKeyVisible(tbl *Table, keys []declaredUniqueKey) error {
	ordered := keys
	if positionsKnown(keys) {
		ordered = append([]declaredUniqueKey(nil), keys...)
		sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].pos < ordered[j].pos })
	}
	for _, k := range ordered {
		switch uniqueKeyNullability(tbl, k.idx) {
		case keyHasNullablePart:
			continue
		case keyNullabilityUnknown:
			return nil
		}
		if !k.idx.Visible {
			return &Error{Code: 3522, SQLState: "HY000", Message: "A primary key index cannot be invisible"}
		}
		return nil
	}
	return nil
}

func positionsKnown(keys []declaredUniqueKey) bool {
	for _, k := range keys {
		if k.pos < 0 {
			return false
		}
	}
	return true
}

type keyNullability int

const (
	keyNotNull keyNullability = iota
	keyHasNullablePart
	keyNullabilityUnknown
)

// uniqueKeyNullability classifies idx for checkImplicitPrimaryKeyVisible. One
// nullable column makes the key nullable whatever its other parts are.
func uniqueKeyNullability(tbl *Table, idx *Index) keyNullability {
	result := keyNotNull
	if len(idx.Columns) == 0 {
		result = keyNullabilityUnknown
	}
	for _, part := range idx.Columns {
		if part.Expr != "" {
			result = keyNullabilityUnknown
			continue
		}
		col := tbl.GetColumn(part.Name)
		if col == nil {
			result = keyNullabilityUnknown
			continue
		}
		if col.Nullable {
			return keyHasNullablePart
		}
	}
	return result
}
