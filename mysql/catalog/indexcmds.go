package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

// createIndex handles standalone CREATE INDEX statements.
func (c *Catalog) createIndex(stmt *nodes.CreateIndexStmt) error {
	// A FULLTEXT key's parser is looked up while MySQL prepares the statement,
	// so 1128 answers ahead of the database and the table it names
	// (mysql:8.0.46).
	if err := validateFulltextParserOption(stmt.Fulltext, stmt.Options); err != nil {
		return err
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
	tableKey := toLower(tableName)
	tbl := db.Tables[tableKey]
	if tbl == nil {
		if db.Views[tableKey] != nil {
			return errWrongObject(db.Name, tableName, "BASE TABLE")
		}
		return errNoSuchTable(db.Name, tableName)
	}

	if err := validateNonPrimaryIndexName(stmt.IndexName); err != nil {
		return err
	}
	if err := validateNoColumnFunctionalParts(stmt.Columns); err != nil {
		return err
	}

	// MySQL drops implicitly generated FK backing indexes (HA_GENERATED_KEY)
	// before checking duplicate index names when adding a new index, and
	// re-creates one at the end only if no remaining index covers the FK.
	snapshot := cloneTable(tbl)
	strippedGen := stripGeneratedFKIndexes(tbl, nil)
	rollback := func(err error) error {
		undoAlter(db, tableKey, tbl, snapshot)
		return err
	}

	// Check for duplicate key name.
	if indexNameExists(tbl, stmt.IndexName) {
		if stmt.IfNotExists {
			undoAlter(db, tableKey, tbl, snapshot)
			return nil
		}
		return rollback(errDupKeyName(stmt.IndexName))
	}

	// A key part that names a column the table does not have is 1072, checked
	// after the key's name: mysql:8.0.46 answers 1061 for a name the table
	// already has and 1072 only for one it does not.
	if err := validateKeyPartsExist(tbl, nil, stmt.Columns); err != nil {
		return rollback(err)
	}

	// Partitioned tables do not support FULLTEXT indexes; MySQL 8.0 answers
	// 1214 ahead of the column-type check (1283).
	if stmt.Fulltext && tbl.Partitioning != nil {
		return rollback(errTableCantHandleFt())
	}

	// Build index columns from AST columns.
	idxCols := make([]*IndexColumn, 0, len(stmt.Columns))
	for _, ic := range stmt.Columns {
		col := &IndexColumn{
			Length:     ic.Length,
			Descending: ic.Desc,
		}
		if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && !ic.Functional {
			col.Name = cr.Column
		} else {
			expr := stripColumnQualifiers(ic.Expr)
			col.Expr = nodeToSQL(expr)
			col.ExprNode = expr
		}
		idxCols = append(idxCols, col)
	}
	// USING RTREE makes a plain CREATE INDEX a SPATIAL one, as it does in a table
	// definition (mysql:8.0.46 stores `SPATIAL KEY g (g)`). On a UNIQUE index,
	// MySQL checks that every key part is geometrical (1687) and then refuses a
	// unique spatial index with 3728.
	effUsingType := resolveIndexUsingType(stmt.IndexType, stmt.Options)
	if stmt.Unique {
		if err := validateUniqueOrPrimaryNotSpatial(tbl, idxCols, effUsingType); err != nil {
			return rollback(err)
		}
	}
	usingSpatial := !stmt.Unique && isSpatialIndexType(effUsingType)
	if err := c.validateIndexColumns(tbl, idxCols, stmt.Fulltext, stmt.Spatial || usingSpatial, stmt.Unique); err != nil {
		return rollback(err)
	}
	seenCols := make(map[string]bool, len(stmt.Columns))
	for _, ic := range stmt.Columns {
		if ic == nil || ic.Functional {
			continue
		}
		if cr, ok := ic.Expr.(*nodes.ColumnRef); ok && cr.Column != "" {
			k := toLower(cr.Column)
			if seenCols[k] {
				return rollback(errDupColumn(cr.Column))
			}
			seenCols[k] = true
		}
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
	case stmt.Spatial || usingSpatial:
		idx.Spatial = true
		idx.IndexType = "SPATIAL"
	default:
		// `USING {BTREE|HASH}` may appear in either grammar position — before ON
		// (stmt.IndexType) or as a trailing index_option (stmt.Options). Resolve both via the
		// shared helper (same as the constraint paths) so the trailing form is not silently
		// dropped (applyIndexOptions ignores USING).
		idx.IndexType = resolveIndexUsingType(stmt.IndexType, stmt.Options)
		idx.Unique = stmt.Unique
	}
	// Fold a meaningless `USING HASH` on InnoDB to the no-explicit-type form so a HASH index
	// authored via CREATE [UNIQUE] INDEX round-trips empty (Fulltext/Spatial are never HASH, so
	// the fold is a no-op for them). Mirrors the CREATE TABLE / ALTER ADD paths.
	coerceInnoDBHashIndex(tbl, idx)

	applyIndexOptions(idx, stmt.Options)
	if err := validateNewIndexCoversPartitioning(tbl, idx); err != nil {
		return rollback(err)
	}
	if err := synthesizeFunctionalIndexColumns(tbl, idx, c.isIgnoreSpaceSQLMode()); err != nil {
		return rollback(err)
	}
	if !idx.Fulltext && !idx.Spatial && hasFunctionalIndexColumns(idx.Columns) {
		if err := c.validateIndexTotalKeyWidth(tbl, idx.Columns); err != nil {
			return rollback(err)
		}
	}

	tbl.Indexes = append(tbl.Indexes, idx)
	restoreGeneratedFKIndexes(tbl, strippedGen)

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
	tableKey := toLower(tableName)
	tbl := db.Tables[tableKey]
	if tbl == nil {
		if db.Views[tableKey] != nil {
			return errWrongObject(db.Name, tableName, "BASE TABLE")
		}
		return errNoSuchTable(db.Name, tableName)
	}

	snapshot := cloneTable(tbl)

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
		return errCantDropKey(stmt.Name)
	}
	removeFunctionalIndexHiddenColumns(tbl, dropped)

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if constraintOfIndex(con, key) {
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			break
		}
	}
	if err := c.validateDroppedIndexesForFK(tbl, []*Index{dropped}); err != nil {
		undoAlter(db, tableKey, tbl, snapshot)
		return err
	}
	if err := validateAutoIncrementKeyed(tbl); err != nil {
		undoAlter(db, tableKey, tbl, snapshot)
		return err
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
