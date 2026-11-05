package catalog

import nodes "github.com/bytebase/omni/mariadb/ast"

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
			Length:          ic.Length,
			Descending:      ic.Desc,
			WithoutOverlaps: ic.WithoutOverlaps,
		}
		if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && !ic.Functional {
			col.Name = cr.Column
		} else {
			col.Expr = nodeToSQL(ic.Expr)
			col.ExprNode = ic.Expr
		}
		idxCols = append(idxCols, col)
	}
	if err := validateKeyPartPrefixes(stmt.Columns, stmt.Spatial); err != nil {
		return err
	}
	if err := validateIndexColumns(tbl, idxCols, stmt.Fulltext, stmt.Spatial); err != nil {
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
		if err := validateVectorIndex(tbl, idxCols); err != nil {
			return err
		}
		idx.Vector = true
		idx.IndexType = "VECTOR"
	default:
		// CREATE INDEX may declare the type before ON or among the options
		// after the column list; reading only the first dropped the second.
		idx.IndexType = resolveIndexUsingType(stmt.IndexType, stmt.Options)
		idx.Unique = stmt.Unique
	}

	applyIndexOptions(idx, stmt.Options)
	if err := validateNewIndexCoversPartitioning(tbl, idx); err != nil {
		return err
	}
	if err := synthesizeFunctionalIndexColumns(tbl, idx); err != nil {
		return err
	}
	if err := validateColsWithoutOverlaps(tbl, idx.Name, idx.Columns); err != nil {
		return err
	}

	tbl.Indexes = append(tbl.Indexes, idx)
	// An IGNORED UNIQUE index can become the implicit primary key.
	if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
		tbl.Indexes = tbl.Indexes[:len(tbl.Indexes)-1]
		return err
	}

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

	// Dropping a key can promote an IGNORED one to the implicit primary key,
	// which MariaDB refuses (4174, verified on 11.8.8), so keep the table to
	// put back.
	var snapshot *Table
	if hasIgnoredIndex(tbl) {
		s := cloneTable(tbl)
		snapshot = &s
	}

	// Find and remove index.
	key := toLower(stmt.Name)
	found := false
	var dropped *Index
	for i, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			dropped = idx
			tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		// IF EXISTS turns the missing index into a warning and leaves the table
		// as it was: mariadb:11.8.8 answers 1091 for DROP INDEX nope ON t and 0
		// for DROP INDEX IF EXISTS nope ON t. A missing table is still an error
		// either way (1146), which the lookup above has already raised.
		if stmt.IfExists {
			return nil
		}
		return errCantDropKey(stmt.Name)
	}
	removeFunctionalIndexHiddenColumns(tbl, dropped)

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if toLower(con.IndexName) == key || toLower(con.Name) == key {
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			break
		}
	}

	if snapshot != nil {
		if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
			*tbl = *snapshot
			return err
		}
	}
	return nil
}
