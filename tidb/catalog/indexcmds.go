package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// isBlobOrTextKeyType reports whether a column of the given normalized data
// type needs a prefix length to appear in a key. TiDB answers 1170 for a key
// part over one without a length; a JSON column is 3152 instead, so it is not
// listed here.
func isBlobOrTextKeyType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "tinyblob", "blob", "mediumblob", "longblob",
		"tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

// validateKeyColumnLengths rejects a key part over a BLOB or TEXT column that
// names no prefix length, which TiDB answers 1170 for (read against
// pingcap/tidb:v8.5.5 on 2026-11-09). It covers PRIMARY KEY, UNIQUE and plain
// KEY alike; FULLTEXT indexes index the whole column and never reach here, and
// a functional key part (an expression, so no column name) is exempt.
func validateKeyColumnLengths(tbl *Table, cols []*IndexColumn) error {
	for _, ic := range cols {
		if ic.Name == "" {
			continue
		}
		col := tbl.GetColumn(ic.Name)
		if col == nil {
			continue
		}
		if ic.Length == 0 {
			if isBlobOrTextKeyType(col.DataType) {
				return errBlobKeyWithoutLength(ic.Name)
			}
			continue
		}
		// A prefix longer than the column's own character length is 1089
		// (ER_WRONG_SUB_KEY): KEY (c(768)) over VARCHAR(200) is rejected.
		if n := columnCharLength(col); n > 0 && ic.Length > n {
			return errWrongSubKey(ic.Name)
		}
	}
	return nil
}

// columnCharLength returns the declared character length of a CHAR/VARCHAR/
// BINARY/VARBINARY column, parsed from its ColumnType, or 0 for a type whose
// prefix length is not bounded by a declared character length.
func columnCharLength(col *Column) int {
	switch strings.ToLower(col.DataType) {
	case "char", "varchar", "binary", "varbinary":
	default:
		return 0
	}
	ct := col.ColumnType
	open := strings.IndexByte(ct, '(')
	if open < 0 {
		return 0
	}
	close := strings.IndexByte(ct[open:], ')')
	if close < 0 {
		return 0
	}
	n := 0
	for _, r := range ct[open+1 : open+close] {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// validateIndexColumnsExist rejects a key part naming a column the table does
// not have, which TiDB v8.5.5 answers 1072 for (measured 2026-11-10). It covers
// PRIMARY KEY, UNIQUE and plain KEY in CREATE TABLE, CREATE INDEX and ALTER
// TABLE ADD alike. A functional key part carries an expression and no column
// name, so it is exempt: v8.5.5 answers 1054 for `INDEX ((c+1))` over a missing
// c, a different error this check does not model. A FULLTEXT key is exempt too,
// because v8.5.5 reads one over a missing column and keeps nothing rather than
// failing.
func validateIndexColumnsExist(tbl *Table, cols []*IndexColumn) error {
	for _, ic := range cols {
		if ic.Name == "" {
			continue
		}
		if tbl.GetColumn(ic.Name) == nil {
			return errKeyColumnDoesNotExist(ic.Name)
		}
	}
	return nil
}

// validateCreateTablePKColumnsExist is validateIndexColumnsExist with the
// wording v8.5.5 uses for a table-level PRIMARY KEY in CREATE TABLE.
func validateCreateTablePKColumnsExist(tbl *Table, cols []*IndexColumn) error {
	for _, ic := range cols {
		if ic.Name == "" {
			continue
		}
		if tbl.GetColumn(ic.Name) == nil {
			return errCreateTablePKColumnDoesNotExist(ic.Name)
		}
	}
	return nil
}

// createIndex handles standalone CREATE INDEX statements.
func (c *Catalog) createIndex(stmt *nodes.CreateIndexStmt) error {
	// Resolve database.
	dbName := ""
	if stmt.Table != nil {
		dbName = stmt.Table.Schema
	}
	db, err := c.resolveDatabase(dbName)
	if err != nil {
		return err
	}

	tableName := stmt.Table.Name
	tbl := db.Tables[toLower(tableName)]
	if tbl == nil {
		return errNoSuchTable(db.Name, tableName)
	}

	// CREATE FULLTEXT INDEX and CREATE SPATIAL INDEX parse, and v8.5.5 refuses
	// both once the table is found: 8200 "FULLTEXT and SPATIAL index is not
	// supported" (read on 2026-11-06).
	if stmt.Fulltext || stmt.Spatial {
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL),
			Message: "FULLTEXT and SPATIAL index is not supported"}
	}

	// A non-fulltext secondary index may not be named PRIMARY (1280). FULLTEXT
	// and SPATIAL were already refused above.
	if !stmt.Fulltext && !stmt.Spatial && !stmt.Vector {
		if err := checkIndexName(stmt.IndexName); err != nil {
			return err
		}
	}

	// Check for duplicate key name.
	if indexNameExists(tbl, stmt.IndexName) {
		if stmt.IfNotExists {
			return nil
		}
		return errDupKeyName(stmt.IndexName)
	}

	// Build index columns from AST columns.
	idxCols := make([]*IndexColumn, 0, len(stmt.Columns))
	for _, ic := range stmt.Columns {
		col := &IndexColumn{
			Length:     ic.Length,
			Descending: ic.Desc,
		}
		// `((v))` is a functional key part over a column, which a vector index requires.
		if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && !ic.Functional {
			col.Name = cr.Column
		} else {
			col.Expr = nodeToSQL(ic.Expr)
		}
		idxCols = append(idxCols, col)
	}

	// Reject functional key parts TiDB refuses (bare column, JSON/GEOMETRY,
	// BLOB/TEXT, non-deterministic function). A VECTOR index is exempt: its key
	// part is a functional part over a column ((v)), which TiDB requires.
	if !stmt.Vector {
		if err := validateFunctionalIndexColumns(stmt.Columns); err != nil {
			return err
		}
	}

	// Every key part must name a column the table has.
	if err := validateIndexColumnsExist(tbl, idxCols); err != nil {
		return err
	}

	// A plain or unique key over a BLOB/TEXT column needs a prefix length;
	// FULLTEXT and SPATIAL were already refused above, so every key that
	// reaches here is held to the rule.
	if err := validateKeyColumnLengths(tbl, idxCols); err != nil {
		return err
	}

	// Determine index type: Fulltext/Spatial override IndexType.
	idx := &Index{
		Name:    stmt.IndexName,
		Table:   tbl,
		Columns: idxCols,
		Visible: true,
	}

	switch {
	case stmt.Fulltext:
		idx.Fulltext = true
		idx.IndexType = "FULLTEXT"
	case stmt.Spatial:
		idx.Spatial = true
		idx.IndexType = "SPATIAL"
	case stmt.Vector:
		idx.Vector = true
		idx.IndexType = vectorIndexHNSW
	default:
		// CREATE INDEX may declare the type before ON or among the options
		// after the column list; reading only the first dropped the second.
		idx.IndexType = resolveIndexUsingType(stmt.IndexType, stmt.Options)
		idx.Unique = stmt.Unique
	}

	applyIndexOptions(idx, stmt.Options)
	if idx.Unique {
		if err := validateNewUniqueKey(tbl, idx); err != nil {
			return err
		}
	}

	tbl.Indexes = append(tbl.Indexes, idx)

	// If unique, also add a UniqueKey constraint.
	if stmt.Unique {
		cols := make([]string, 0, len(idxCols))
		for _, ic := range idxCols {
			if ic.Name != "" {
				cols = append(cols, ic.Name)
			}
		}
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      stmt.IndexName,
			Type:      ConUniqueKey,
			Table:     tbl,
			Columns:   cols,
			IndexName: stmt.IndexName,
		})
	}

	return nil
}

// dropIndex handles standalone DROP INDEX statements.
func (c *Catalog) dropIndex(stmt *nodes.DropIndexStmt) error {
	// DROP HYPO INDEX removes a hypothetical index the optimizer holds for
	// what-if planning; it never touches a persistent index. This catalog
	// models no hypothetical indexes, so the statement changes nothing here.
	// Routing it through the physical removal below would delete a real index
	// (and its constraint) that happens to share the name.
	if stmt.IsHypo {
		return nil
	}

	// Resolve database.
	dbName := ""
	if stmt.Table != nil {
		dbName = stmt.Table.Schema
	}
	db, err := c.resolveDatabase(dbName)
	if err != nil {
		return err
	}

	tableName := stmt.Table.Name
	tbl := db.Tables[toLower(tableName)]
	if tbl == nil {
		return errNoSuchTable(db.Name, tableName)
	}

	// Find and remove index.
	key := toLower(stmt.Name)
	found := false
	for i, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		// IF EXISTS turns the missing index into a warning and leaves the table
		// as it was: pingcap/tidb:v8.5.5 answers 1091 for DROP INDEX nope ON t
		// and 0 for DROP INDEX IF EXISTS nope ON t. A missing table is still an
		// error either way (1146), which the lookup above has already raised.
		if stmt.IfExists {
			return nil
		}
		return errIndexDoesntExist(stmt.Name)
	}

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if toLower(con.IndexName) == key || toLower(con.Name) == key {
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			break
		}
	}

	return nil
}
