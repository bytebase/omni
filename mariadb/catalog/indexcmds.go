package catalog

import nodes "github.com/bytebase/omni/mariadb/ast"

// createIndex handles standalone CREATE INDEX statements.
func (c *Catalog) createIndex(stmt *nodes.CreateIndexStmt) error {
	defer c.ignoreDescendingKeys(stmt.Table)()
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
		if db.Views[toLower(tableName)] != nil {
			return errWrongObject(db.Name, tableName, "BASE TABLE")
		}
		return errNoSuchTable(db.Name, tableName)
	}

	if err := validateNonPrimaryIndexName(stmt.IndexName); err != nil {
		return err
	}

	// A name the table holds is a note under IF NOT EXISTS, which MariaDB reads
	// before the key: `CREATE INDEX IF NOT EXISTS k (nosuch)` over a table that
	// has k runs. Otherwise the name waits for the key's own errors below (1061,
	// or another one).
	nameExists := indexNameExists(tbl, stmt.IndexName)
	if nameExists && stmt.IfNotExists {
		return nil
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
	// A FULLTEXT key's parser is looked up before its columns, as the ALTER
	// path's is.
	if stmt.Fulltext {
		if err := validateFulltextParser(stmt.Options); err != nil {
			return err
		}
	}
	// A WITHOUT OVERLAPS part's period is looked up before the key's columns,
	// type, and colliding name: measured on mariadb:11.8.8, `CREATE UNIQUE
	// INDEX u ON a (missing, missing WITHOUT OVERLAPS)` answers 4156.
	if err := validateDeclaredWithoutOverlapsParts(tbl.AppPeriodName, stmt.Columns); err != nil {
		return err
	}
	// The key's own validity is judged before a name it collides with:
	// measured on mariadb:11.8.8, `CREATE SPATIAL INDEX idx_d ON a (b)` over an
	// INT b answers 1210 and `CREATE FULLTEXT INDEX idx_d ON a (b)` 1283 where
	// the table already holds idx_d (1061), and a column that is no column of
	// the table answers 1072 and a length over an INT column 1089, where a plain
	// repeated column answers 1061.
	for _, ic := range idxCols {
		if ic.Name != "" && !ic.WithoutOverlaps && tbl.GetColumn(ic.Name) == nil {
			return errKeyColumnDoesNotExist(ic.Name)
		}
	}
	// The key's columns are looked up before its type: measured on
	// mariadb:11.8.8, `CREATE INDEX k USING RTREE ON t(missing)` answers 1072,
	// and `ON t(b)` over an INT column 1210.
	if err := validateKeyTypeArguments(resolveIndexUsingType(stmt.IndexType, stmt.Options)); err != nil {
		return err
	}
	long := false
	switch {
	case stmt.Vector:
		if err := validateVectorColumns(tbl, idxCols); err != nil {
			return err
		}
	case stmt.Unique && !stmt.Fulltext && !stmt.Spatial:
		var err error
		if long, err = validateUniqueIndexColumns(tbl, idxCols); err != nil {
			return err
		}
	default:
		if err := validateIndexColumns(tbl, idxCols, stmt.Fulltext, stmt.Spatial); err != nil {
			return err
		}
	}

	// Snapshot before dropping any covered implicit FK backing index or adding
	// hidden columns so any subsequent failure leaves the table as it was.
	snapshot := cloneTable(tbl)
	resolvedType := resolveIndexUsingType(stmt.IndexType, stmt.Options)
	if long {
		resolvedType = "HASH"
	}
	dropImplicitFKIndexesCoveredBy(tbl, idxCols, resolvedType, stmt.Fulltext, stmt.Spatial, stmt.Vector)
	nameExists = indexNameExists(tbl, stmt.IndexName)

	if nameExists {
		*tbl = snapshot
		if stmt.OrReplace {
			return c.replaceIndex(tbl, stmt)
		}
		return errDupKeyName(stmt.IndexName)
	}
	if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
		*tbl = snapshot
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
			*tbl = snapshot
			return err
		}
		idx.Vector = true
		idx.IndexType = "VECTOR"
	default:
		// CREATE INDEX may declare the type before ON or among the options
		// after the column list; reading only the first dropped the second.
		idx.IndexType = resolveIndexUsingType(stmt.IndexType, stmt.Options)
		idx.Unique = stmt.Unique
		if long {
			idx.IndexType = "HASH"
		}
	}

	applyIndexOptions(idx, stmt.Options)
	if err := validateNewIndexCoversPartitioning(tbl, idx); err != nil {
		*tbl = snapshot
		return err
	}
	// Hidden columns go on the table before the checks that can still fail, and a statement that
	// fails leaves the table as it was, or a retry would allocate another set beside them.
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
	if err := validateAutoIncrementKeyed(tbl); err != nil {
		*tbl = snapshot
		return err
	}
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
	if err := c.validateKeyEngineOptions(tbl, stmt.Options, stmt.Vector); err != nil {
		*tbl = snapshot
		return err
	}
	if tbl.TablespaceDiscarded && eqFoldStr(stmt.Algorithm, "COPY") {
		*tbl = snapshot
		return errTablespaceDiscarded(tbl.Name)
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
	if _, err := removeIndexFromTable(tbl, drop); err != nil {
		*tbl = snapshot
		return err
	}
	create := *stmt
	create.OrReplace = false
	if err := c.createIndex(&create); err != nil {
		*tbl = snapshot
		return err
	}
	if err := validateAutoIncrementKeyed(tbl); err != nil {
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
		if db.Views[toLower(tableName)] != nil {
			return errWrongObject(db.Name, tableName, "BASE TABLE")
		}
		return errNoSuchTable(db.Name, tableName)
	}

	// Dropping a key can unkey an AUTO_INCREMENT column (1075), promote an
	// IGNORED index to the implicit primary key (4174), or leave a foreign key
	// without a covering index (1553), all verified on 11.8.8, so keep the table
	// to put back when any fails.
	var snapshot *Table
	if hasIgnoredIndex(tbl) || needsAutoIncrementCheck(tbl, nil) || (c.foreignKeyChecks && tableHasForeignKey(tbl)) {
		s := cloneTable(tbl)
		snapshot = &s
	}

	removed, err := removeIndexFromTable(tbl, stmt)
	if err != nil || !removed {
		return err
	}

	if snapshot != nil {
		if c.foreignKeyChecks && !checkFKIndexCoverage(tbl) {
			*tbl = *snapshot
			return errDropIndexFK(stmt.Name)
		}
		if err := validateAutoIncrementKeyed(tbl); err != nil {
			*tbl = *snapshot
			return err
		}
		if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
			*tbl = *snapshot
			return err
		}
	}
	return nil
}

func removeIndexFromTable(tbl *Table, stmt *nodes.DropIndexStmt) (bool, error) {
	key := toLower(stmt.Name)
	var dropped *Index
	for i, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			dropped = idx
			tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
			break
		}
	}
	if dropped == nil {
		// IF EXISTS turns the missing index into a warning and leaves the table
		// as it was: mariadb:11.8.8 answers 1091 for DROP INDEX nope ON t and 0
		// for DROP INDEX IF EXISTS nope ON t. A missing table is still an error
		// either way (1146), which the lookup above has already raised.
		if stmt.IfExists {
			return false, nil
		}
		return false, errCantDropKey("INDEX", stmt.Name)
	}
	removeFunctionalIndexHiddenColumns(tbl, dropped)

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if constraintOfIndex(con, key) {
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			break
		}
	}
	return true, nil
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
