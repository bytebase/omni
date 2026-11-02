package catalog

// errPrimaryKeyInvisible is ER_PK_INDEX_CANT_BE_INVISIBLE, which TiDB v8.5.5
// answers for a table whose primary key, explicit or implicit, is invisible.
func errPrimaryKeyInvisible() error {
	return &Error{Code: 3522, SQLState: "HY000", Message: "A primary key index cannot be invisible"}
}

// checkImplicitPrimaryKeyVisible refuses a table without a primary key whose
// implicit one is invisible. TiDB takes the first UNIQUE index whose columns
// are all NOT NULL for the primary key (TableInfo.GetPrimaryKey), passing over
// an expression index, whose column is hidden, and refuses CREATE TABLE with
// 3522 when that index is invisible.
//
// "First" is in TiDB's index order, not the statement's: TiDB builds the
// indexes of the table-level definitions before those a column definition
// declares (UNIQUE, SERIAL), so
//
//	CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
//
// takes u for the primary key and fails, where MySQL 8.0 takes a and succeeds.
// Callers pass the indexes in that order, one group after another.
func checkImplicitPrimaryKeyVisible(tbl *Table, groups ...[]*Index) error {
	for _, group := range groups {
		for _, idx := range group {
			if idx.Primary || !idx.Unique || !columnsAllNotNull(tbl, idx) {
				continue
			}
			if !idx.Visible {
				return errPrimaryKeyInvisible()
			}
			return nil
		}
	}
	return nil
}

// columnsAllNotNull reports whether every part of idx is a column declared NOT
// NULL. A prefix part counts; an expression part does not, since TiDB keys it on
// a hidden column it never takes for a primary key.
func columnsAllNotNull(tbl *Table, idx *Index) bool {
	if len(idx.Columns) == 0 {
		return false
	}
	for _, part := range idx.Columns {
		if part.Expr != "" {
			return false
		}
		col := tbl.GetColumn(part.Name)
		if col == nil || col.Nullable {
			return false
		}
	}
	return true
}
