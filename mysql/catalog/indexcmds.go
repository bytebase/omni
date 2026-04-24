package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

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

	if err := validateNonPrimaryIndexName(stmt.IndexName); err != nil {
		return err
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
		if cr, ok := ic.Expr.(*nodes.ColumnRef); ok {
			col.Name = cr.Column
		} else {
			col.Expr = nodeToSQL(ic.Expr)
		}
		idxCols = append(idxCols, col)
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
	default:
		idx.IndexType = stmt.IndexType
		idx.Unique = stmt.Unique
	}

	applyIndexOptions(idx, stmt.Options)
	synthesizeFunctionalIndexColumns(tbl, idx)

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
		return errCantDropKey(stmt.Name)
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
