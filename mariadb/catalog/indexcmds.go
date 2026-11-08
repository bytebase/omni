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
		if stmt.OrReplace {
			return c.replaceIndex(tbl, stmt)
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
	// Hidden columns go on the table before the checks that can still fail, and a statement that
	// fails leaves the table as it was, or a retry would allocate another set beside them.
	snapshot := cloneTable(tbl)
	if err := synthesizeFunctionalIndexColumns(tbl, idx); err != nil {
		*tbl = snapshot
		return err
	}
	if err := validateColsWithoutOverlaps(tbl, idx.Name, idx.Columns); err != nil {
		*tbl = snapshot
		return err
	}
	// mariadb:11.8.8 answers 1072 for CREATE INDEX k ON t (nosuch), which a CREATE OR REPLACE
	// must not read as leave to drop the old index.
	for _, ic := range idxCols {
		if ic.Name != "" && !ic.WithoutOverlaps && tbl.GetColumn(ic.Name) == nil {
			*tbl = snapshot
			return errKeyColumnDoesNotExist(ic.Name)
		}
	}

	tbl.Indexes = append(tbl.Indexes, idx)
	// An IGNORED UNIQUE index can become the implicit primary key.
	if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
		*tbl = snapshot
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

// replaceIndex runs CREATE OR REPLACE INDEX over an index that exists: the old index goes, and
// the new one is held to every check as though it had never been there, so a vector index may
// replace the table's only vector index. A statement that fails leaves the old index in place:
// mariadb:11.8.8 answers 1072 for a missing key column and keeps what the table had.
func (c *Catalog) replaceIndex(tbl *Table, stmt *nodes.CreateIndexStmt) error {
	snapshot := cloneTable(tbl)
	drop := &nodes.DropIndexStmt{Name: stmt.IndexName, Table: stmt.Table}
	if err := c.dropIndex(drop); err != nil {
		*tbl = snapshot
		return err
	}
	create := *stmt
	create.OrReplace = false
	if err := c.createIndex(&create); err != nil {
		*tbl = snapshot
		return err
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
		return errCantDropKey("INDEX", stmt.Name)
	}
	removeFunctionalIndexHiddenColumns(tbl, dropped)

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if constraintOfIndex(con, key) {
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

// constraintOfIndex reports whether con is what the index named key holds up: its primary or unique
// key. A CHECK or a foreign key may have the same name without being tied to the index, and the
// server keeps it when the index goes.
func constraintOfIndex(con *Constraint, key string) bool {
	switch con.Type {
	case ConPrimaryKey, ConUniqueKey:
		return toLower(con.IndexName) == key || toLower(con.Name) == key
	}
	return false
}
