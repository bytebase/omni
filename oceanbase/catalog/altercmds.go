package catalog

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
)

func (c *Catalog) alterTable(stmt *nodes.AlterTableStmt) error {
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
	key := toLower(tableName)
	tbl := db.Tables[key]
	if tbl == nil {
		return errNoSuchTable(db.Name, tableName)
	}
	if err := checkKeyOptions(alterKeys(stmt.Commands)); err != nil {
		return err
	}

	repartition := alterRepartition(stmt.Commands)
	if repartition != nil {
		if err := checkPartitionNames(tbl, repartition); err != nil {
			return err
		}
	}
	if tbl.Partitioning == nil && repartition == nil {
		return c.alterTableCommands(db, tbl, key, stmt)
	}
	// A unique key on a partitioned table must cover the partitioning columns.
	// OceanBase checks the table the statement leaves, after a duplicate column (1060)
	// or a missing index to drop (1091) in ADD UNIQUE (b), ADD COLUMN a ..., and leaves
	// the table as it was: snapshot it to undo the commands when the check fails.
	before := make(map[*Index]bool, len(tbl.Indexes))
	for _, idx := range tbl.Indexes {
		before[idx] = true
	}
	snapshot := cloneTable(tbl)
	if err := c.alterTableCommands(db, tbl, key, stmt); err != nil {
		return err
	}
	if err := validateAlteredUniqueKeys(tbl, before, stmt.Commands, repartition); err != nil {
		undoAlter(db, key, tbl, snapshot)
		return err
	}
	if repartition != nil && tbl.Partitioning != nil {
		// The catalog does not apply a repartition, so the partitioning it keeps no
		// longer says which columns a unique key needs.
		tbl.Partitioning.keyColumns = nil
	}
	return nil
}

// alterTableCommands runs the commands of an ALTER TABLE on tbl, found under key in db.
func (c *Catalog) alterTableCommands(db *Database, tbl *Table, key string, stmt *nodes.AlterTableStmt) error {
	if len(stmt.Commands) <= 1 {
		// Single command: no rollback needed.
		if len(stmt.Commands) == 1 {
			return c.execAlterCmd(db, tbl, stmt.Commands[0])
		}
		return nil
	}

	// Multi-command ALTER: MySQL treats this as atomic.
	// Snapshot the table so we can rollback on any sub-command failure.
	snapshot := cloneTable(tbl)
	origKey := key

	for _, cmd := range stmt.Commands {
		if err := c.execAlterCmd(db, tbl, cmd); err != nil {
			undoAlter(db, origKey, tbl, snapshot)
			return err
		}
	}
	// Clear transient cleanup tracking after successful multi-command ALTER.
	tbl.droppedByCleanup = nil
	return nil
}

// undoAlter puts tbl back as snapshot holds it, under key in db. A RENAME TO earlier
// in the statement may have moved the table to another database, under another name,
// so its entry there goes first.
func undoAlter(db *Database, key string, tbl *Table, snapshot Table) {
	if cur := tbl.Database; cur != nil && cur.Tables[toLower(tbl.Name)] == tbl {
		delete(cur.Tables, toLower(tbl.Name))
	}
	*tbl = snapshot
	db.Tables[key] = tbl
}

func (c *Catalog) execAlterCmd(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	switch cmd.Type {
	case nodes.ATAddColumn:
		return c.alterAddColumn(tbl, cmd)
	case nodes.ATDropColumn:
		return c.alterDropColumn(tbl, cmd)
	case nodes.ATModifyColumn:
		return c.alterModifyColumn(tbl, cmd)
	case nodes.ATChangeColumn:
		return c.alterChangeColumn(tbl, cmd)
	case nodes.ATAddIndex, nodes.ATAddConstraint:
		return c.alterAddConstraint(tbl, cmd)
	case nodes.ATDropIndex:
		return c.alterDropIndex(tbl, cmd)
	case nodes.ATDropConstraint:
		return c.alterDropConstraint(tbl, cmd)
	case nodes.ATRenameColumn:
		return c.alterRenameColumn(tbl, cmd)
	case nodes.ATRenameIndex:
		return c.alterRenameIndex(tbl, cmd)
	case nodes.ATRenameTable:
		return c.alterRenameTable(db, tbl, cmd)
	case nodes.ATTableOption:
		return c.alterTableOption(tbl, cmd)
	case nodes.ATAlterColumnDefault:
		return c.alterColumnDefault(tbl, cmd)
	case nodes.ATAlterIndexVisible:
		return c.alterIndexVisibility(tbl, cmd, true)
	case nodes.ATAlterIndexInvisible:
		return c.alterIndexVisibility(tbl, cmd, false)
	case nodes.ATAlterCheckEnforced:
		return c.alterCheckEnforced(tbl, cmd)
	case nodes.ATConvertCharset:
		return c.alterConvertCharset(tbl, cmd)
	case nodes.ATAddPartition:
		return c.alterAddPartition(tbl, cmd)
	case nodes.ATDropPartition:
		return c.alterDropPartition(tbl, cmd)
	case nodes.ATTruncatePartition:
		return c.alterTruncatePartition(tbl, cmd)
	case nodes.ATCoalescePartition:
		return c.alterCoalescePartition(tbl, cmd)
	case nodes.ATReorganizePartition:
		return c.alterReorganizePartition(tbl, cmd)
	case nodes.ATExchangePartition:
		return c.alterExchangePartition(db, tbl, cmd)
	case nodes.ATRemovePartitioning:
		tbl.Partitioning = nil
		return nil
	default:
		// Unsupported alter command; silently ignore.
		return nil
	}
}

// alterAddColumn adds a new column to the table.
func (c *Catalog) alterAddColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	// Handle multi-column parenthesized form: ADD (col1 INT, col2 INT, ...)
	if len(cmd.Columns) > 0 {
		for _, colDef := range cmd.Columns {
			if err := c.addSingleColumn(tbl, colDef, false, ""); err != nil {
				return err
			}
		}
		return nil
	}

	colDef := cmd.Column
	if colDef == nil {
		return nil
	}

	return c.addSingleColumn(tbl, colDef, cmd.First, cmd.After)
}

// addSingleColumn adds one column definition to the table.
func (c *Catalog) addSingleColumn(tbl *Table, colDef *nodes.ColumnDef, first bool, after string) error {
	if err := checkAddedGeneratedCheck(colDef); err != nil {
		return err
	}
	if err := validateColumnDefSemantics(colDef); err != nil {
		return err
	}
	if err := checkAddedColumn(tbl, colDef); err != nil {
		return err
	}
	colKey := toLower(colDef.Name)
	if _, exists := tbl.colByName[colKey]; exists {
		return errDupColumn(colDef.Name)
	}

	col := buildColumnFromDef(tbl, colDef)
	if col.AutoIncrement && tableHasOtherAutoIncrementColumn(tbl, "") {
		return errIncorrectAutoIncrementDefinition()
	}
	if err := insertColumn(tbl, col, first, after); err != nil {
		return err
	}

	// Process column-level constraints that produce indexes/constraints.
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrPrimaryKey:
			// Check for duplicate PK.
			for _, idx := range tbl.Indexes {
				if idx.Primary {
					return errMultiplePriKey()
				}
			}
			if err := validatePrimaryKeyColumns(tbl, []string{colDef.Name}); err != nil {
				return err
			}
			col.Nullable = false
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      "PRIMARY",
				Table:     tbl,
				Columns:   []*IndexColumn{{Name: colDef.Name}},
				Unique:    true,
				Primary:   true,
				IndexType: "",
				Visible:   true,
			})
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      "PRIMARY",
				Type:      ConPrimaryKey,
				Table:     tbl,
				Columns:   []string{colDef.Name},
				IndexName: "PRIMARY",
			})
		case nodes.ColConstrUnique:
			idxName := allocIndexName(tbl, colDef.Name)
			tbl.Indexes = append(tbl.Indexes, &Index{
				Name:      idxName,
				Table:     tbl,
				Columns:   []*IndexColumn{{Name: colDef.Name}},
				Unique:    true,
				IndexType: "",
				Visible:   true,
			})
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      idxName,
				Type:      ConUniqueKey,
				Table:     tbl,
				Columns:   []string{colDef.Name},
				IndexName: idxName,
			})
		}
	}

	return nil
}

// alterDropColumn removes a column from the table.
func (c *Catalog) alterDropColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	colKey := toLower(cmd.Name)
	if _, exists := tbl.colByName[colKey]; !exists {
		if cmd.IfExists {
			return nil
		}
		// MySQL 8.0 returns error 1091 for DROP COLUMN on nonexistent column,
		// same as DROP INDEX: "Can't DROP 'x'; check that column/key exists".
		return errCantDropKey(cmd.Name)
	}
	idx := tbl.colByName[colKey]
	if tbl.Columns[idx].Hidden == ColumnHiddenSystem {
		return errCannotDropColumnFunctionalIndex(cmd.Name)
	}
	// OceanBase keeps a table's rows in primary-key order and drops no column of the key
	// (TestOracle_OceanBaseRefusals), where MySQL rebuilds the key without it.
	for _, index := range tbl.Indexes {
		if !index.Primary {
			continue
		}
		for _, ic := range index.Columns {
			if toLower(ic.Name) == colKey {
				return &Error{Code: 1235, SQLState: "0A000", Message: "drop rowkey column is not supported"}
			}
		}
	}
	if len(tbl.Columns) == 1 {
		return errCantRemoveAllFields()
	}

	// Check if column is referenced by a generated column expression.
	for _, col := range tbl.Columns {
		if col.Generated != nil && generatedExprReferencesColumn(col.Generated.Expr, cmd.Name) {
			return errDependentByGeneratedColumn(cmd.Name, col.Name, tbl.Name)
		}
	}

	// Check if column is referenced by a foreign key constraint.
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			for _, col := range con.Columns {
				if toLower(col) == colKey {
					return &Error{
						Code:     1828,
						SQLState: "HY000",
						Message:  fmt.Sprintf("Cannot drop column '%s': needed in a foreign key constraint '%s'", cmd.Name, con.Name),
					}
				}
			}
		}
	}

	// Remove column from indexes; if index becomes empty, remove it entirely.
	cleanupIndexesForDroppedColumn(tbl, cmd.Name)

	tbl.Columns = append(tbl.Columns[:idx], tbl.Columns[idx+1:]...)
	rebuildColIndex(tbl)
	return nil
}

// cleanupIndexesForDroppedColumn removes references to a dropped column from
// all indexes. If an index loses all columns, it is removed entirely.
// Associated constraints are also cleaned up.
func cleanupIndexesForDroppedColumn(tbl *Table, colName string) {
	colKey := toLower(colName)

	// Clean up indexes.
	newIndexes := make([]*Index, 0, len(tbl.Indexes))
	removedIndexNames := make(map[string]bool)
	for _, idx := range tbl.Indexes {
		// A column dropped from a STORING goes without taking the index along: 4.3.5.6
		// reads KEY k (a) STORING (b) back as KEY k (a) once b is gone.
		idx.Storing = slices.DeleteFunc(idx.Storing, func(name string) bool { return toLower(name) == colKey })

		// Remove the column from this index.
		newCols := make([]*IndexColumn, 0, len(idx.Columns))
		for _, ic := range idx.Columns {
			if toLower(ic.Name) != colKey {
				newCols = append(newCols, ic)
			}
		}
		if len(newCols) == 0 {
			// Index has no columns left — remove it.
			nameKey := toLower(idx.Name)
			removedIndexNames[nameKey] = true
			// Track for multi-command ALTER so explicit DROP INDEX succeeds.
			if tbl.droppedByCleanup == nil {
				tbl.droppedByCleanup = make(map[string]bool)
			}
			tbl.droppedByCleanup[nameKey] = true
			continue
		}
		idx.Columns = newCols
		newIndexes = append(newIndexes, idx)
	}
	tbl.Indexes = newIndexes

	// Clean up constraints that reference removed indexes.
	if len(removedIndexNames) > 0 {
		newConstraints := make([]*Constraint, 0, len(tbl.Constraints))
		for _, con := range tbl.Constraints {
			if removedIndexNames[toLower(con.IndexName)] || removedIndexNames[toLower(con.Name)] {
				continue
			}
			newConstraints = append(newConstraints, con)
		}
		tbl.Constraints = newConstraints
	}

	// Also update constraint column lists for remaining constraints.
	for _, con := range tbl.Constraints {
		newCols := make([]string, 0, len(con.Columns))
		for _, col := range con.Columns {
			if toLower(col) != colKey {
				newCols = append(newCols, col)
			}
		}
		con.Columns = newCols
	}
}

// alterModifyColumn replaces a column definition in-place (same name).
func (c *Catalog) alterModifyColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if cmd.Column == nil {
		return nil
	}
	return c.alterReplaceColumn(tbl, cmd.Column.Name, cmd)
}

// alterChangeColumn replaces a column (old name -> new name + new definition).
func (c *Catalog) alterChangeColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if cmd.Column == nil {
		return nil
	}
	return c.alterReplaceColumn(tbl, cmd.Name, cmd)
}

// alterReplaceColumn is the shared implementation for MODIFY and CHANGE COLUMN.
// oldName is the existing column to replace; cmd.Column defines the new column.
func (c *Catalog) alterReplaceColumn(tbl *Table, oldName string, cmd *nodes.AlterTableCmd) error {
	colDef := cmd.Column
	if err := validateColumnDefSemantics(colDef); err != nil {
		return err
	}
	oldKey := toLower(oldName)
	idx, exists := tbl.colByName[oldKey]
	if !exists {
		return errNoSuchColumn(oldName, tbl.Name)
	}

	// Check if new name conflicts with existing column (unless same).
	newKey := toLower(colDef.Name)
	if newKey != oldKey {
		if _, dup := tbl.colByName[newKey]; dup {
			return errDupColumn(colDef.Name)
		}
	}

	// Check for VIRTUAL<->STORED storage type change (MySQL 8.0 error 3106).
	oldCol := tbl.Columns[idx]
	if oldCol.Generated != nil && colDef.Generated != nil {
		if oldCol.Generated.Stored != colDef.Generated.Stored {
			return errUnsupportedGeneratedStorageChange(colDef.Name, tbl.Name)
		}
	}

	col := buildColumnFromDef(tbl, colDef)
	if col.AutoIncrement && tableHasOtherAutoIncrementColumn(tbl, oldName) {
		return errIncorrectAutoIncrementDefinition()
	}
	col.Position = idx + 1
	tbl.Columns[idx] = col

	// Update index/constraint column references if name changed.
	if newKey != oldKey {
		updateColumnRefsInIndexes(tbl, oldName, colDef.Name)
	}

	// Handle repositioning.
	if cmd.First || cmd.After != "" {
		tbl.Columns = append(tbl.Columns[:idx], tbl.Columns[idx+1:]...)
		rebuildColIndex(tbl)
		if err := insertColumn(tbl, col, cmd.First, cmd.After); err != nil {
			return err
		}
	} else {
		rebuildColIndex(tbl)
	}

	return nil
}

// alterAddConstraint adds a constraint or index to the table.
func (c *Catalog) alterAddConstraint(tbl *Table, cmd *nodes.AlterTableCmd) error {
	con := cmd.Constraint
	if con == nil {
		return nil
	}

	cols := extractColumnNames(con)
	// A primary key's STORING is taken and dropped unread: on 4.3.5.6, ALTER TABLE t
	// ADD PRIMARY KEY (a) STORING (nosuch) succeeds and the key reads back without it.
	if con.Type != nodes.ConstrPrimaryKey {
		if err := validateStoringColumns(tbl, con.IndexOptions); err != nil {
			return err
		}
	}

	switch con.Type {
	case nodes.ConstrPrimaryKey:
		// Check for duplicate PK.
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				return errMultiplePriKey()
			}
		}
		if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
			return err
		}
		// Mark PK columns as NOT NULL.
		for _, colName := range cols {
			col := tbl.GetColumn(colName)
			if col != nil {
				col.Nullable = false
			}
		}
		idxCols := buildIndexColumns(con)
		pkIdx := &Index{
			Name:      "PRIMARY",
			Table:     tbl,
			Columns:   idxCols,
			Unique:    true,
			Primary:   true,
			IndexType: "",
			Visible:   true,
		}
		applyIndexOptions(pkIdx, con.IndexOptions)
		// MySQL refuses an invisible primary key with 3522. OceanBase takes
		// INVISIBLE here and ignores it: on 4.3.5.6, ALTER TABLE t ADD PRIMARY KEY
		// (id) INVISIBLE succeeds, SHOW CREATE TABLE prints PRIMARY KEY (`id`), and
		// information_schema.STATISTICS reports the key visible. It drops a COMMENT
		// given here the same way, though it keeps one given in CREATE TABLE
		// (TestOracle_OceanBaseAlterAddPrimaryKey). It drops STORING too. CREATE TABLE
		// refuses INVISIBLE and STORING at parse time.
		pkIdx.Visible = true
		pkIdx.Comment = ""
		pkIdx.Storing = nil
		tbl.Indexes = append(tbl.Indexes, pkIdx)
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      "PRIMARY",
			Type:      ConPrimaryKey,
			Table:     tbl,
			Columns:   cols,
			IndexName: "PRIMARY",
		})

	case nodes.ConstrUnique:
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" {
			if err := validateNonPrimaryIndexName(idxName); err != nil {
				return err
			}
			if indexNameExists(tbl, idxName) {
				return errDupKeyName(idxName)
			}
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, idxCols, false, false); err != nil {
			return err
		}
		uqIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			Unique:    true,
			IndexType: resolveConstraintIndexType(con),
			Visible:   true,
		}
		applyIndexOptions(uqIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, uqIdx)
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      idxName,
			Type:      ConUniqueKey,
			Table:     tbl,
			Columns:   cols,
			IndexName: idxName,
		})

	case nodes.ConstrForeignKey:
		conName := con.Name
		if conName == "" {
			conName = c.madeUpConstraintName(tbl.Name, madeUpForeignKeyInfix, func(name string) bool {
				return foreignKeyConstraintNameExists(tbl.Database, nil, name)
			})
		}
		if foreignKeyConstraintNameExists(tbl.Database, nil, conName) {
			return errFKDupName(tbl.Name)
		}
		refDBName := ""
		refTable := ""
		if con.RefTable != nil {
			refDBName = con.RefTable.Schema
			refTable = con.RefTable.Name
		}
		fkCon := &Constraint{
			Name:        conName,
			Type:        ConForeignKey,
			Table:       tbl,
			Columns:     cols,
			RefDatabase: refDBName,
			RefTable:    refTable,
			RefColumns:  con.RefColumns,
			OnDelete:    refActionToString(con.OnDelete),
			OnUpdate:    refActionToString(con.OnUpdate),
		}
		// Validate FK before adding (unless foreign_key_checks=0).
		db := tbl.Database
		if c.foreignKeyChecks {
			if err := c.validateSingleFK(db, tbl, fkCon); err != nil {
				return err
			}
		}
		tbl.Constraints = append(tbl.Constraints, fkCon)
		// No index is made for it: see the foreign keys of CREATE TABLE.

	case nodes.ConstrCheck:
		conName := con.Name
		if conName == "" {
			conName = c.madeUpConstraintName(tbl.Name, madeUpCheckInfix, func(name string) bool {
				return checkConstraintNameExists(tbl.Database, nil, name)
			})
		}
		if err := checkTableCheckExpr(tbl, conName, con.Expr); err != nil {
			return err
		}
		if checkConstraintNameExists(tbl.Database, nil, conName) {
			return errCheckConstraintDupName(conName)
		}
		if err := validateCheckExpr(conName, con.Expr); err != nil {
			return err
		}
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:        conName,
			Type:        ConCheck,
			Table:       tbl,
			CheckExpr:   nodeToSQL(con.Expr),
			NotEnforced: con.NotEnforced,
		})

	case nodes.ConstrIndex:
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" {
			if err := validateNonPrimaryIndexName(idxName); err != nil {
				return err
			}
			if indexNameExists(tbl, idxName) {
				return errDupKeyName(idxName)
			}
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, idxCols, false, false); err != nil {
			return err
		}
		keyIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			IndexType: resolveConstraintIndexType(con),
			Visible:   true,
		}
		// The CREATE TABLE path reads these; this one used to drop them, so an index
		// added by ALTER lost its comment, visibility, block size and scope.
		applyIndexOptions(keyIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, keyIdx)

	case nodes.ConstrFulltextIndex:
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" {
			if err := validateNonPrimaryIndexName(idxName); err != nil {
				return err
			}
			if indexNameExists(tbl, idxName) {
				return errDupKeyName(idxName)
			}
		}
		idxCols := buildIndexColumns(con)
		ftIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			Fulltext:  true,
			IndexType: "FULLTEXT",
			Visible:   true,
		}
		applyIndexOptions(ftIdx, con.IndexOptions)
		if err := checkIndexOptions(con.IndexOptions, false); err != nil {
			return err
		}
		tbl.Indexes = append(tbl.Indexes, ftIdx)

	case nodes.ConstrSpatialIndex:
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" {
			if err := validateNonPrimaryIndexName(idxName); err != nil {
				return err
			}
			if indexNameExists(tbl, idxName) {
				return errDupKeyName(idxName)
			}
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, idxCols, false, true); err != nil {
			return err
		}
		spIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			Spatial:   true,
			IndexType: "SPATIAL",
			Visible:   true,
		}
		applyIndexOptions(spIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, spIdx)
	}

	return nil
}

// alterDropIndex removes an index (and any associated constraint) by name.
func (c *Catalog) alterDropIndex(tbl *Table, cmd *nodes.AlterTableCmd) error {
	name := cmd.Name
	key := toLower(name)

	found := false
	for i, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		if cmd.IfExists {
			return nil
		}
		// If the index was auto-removed by DROP COLUMN cleanup in this
		// multi-command ALTER, treat as success (matches MySQL 8.0 behavior).
		if tbl.droppedByCleanup[key] {
			return nil
		}
		return errCantDropKey(name)
	}

	// Also remove any constraint that references this index.
	for i, con := range tbl.Constraints {
		if constraintOfIndex(con, key) {
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			break
		}
	}

	return nil
}

// alterDropConstraint removes a constraint by name.
func (c *Catalog) alterDropConstraint(tbl *Table, cmd *nodes.AlterTableCmd) error {
	name := cmd.Name
	key := toLower(name)

	found := false
	isForeignKey := false
	for i, con := range tbl.Constraints {
		if toLower(con.Name) == key {
			isForeignKey = (con.Type == ConForeignKey)
			tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		if cmd.IfExists {
			return nil
		}
		return errCantDropKey(name)
	}

	// A foreign key has no index of its own, so an index of its name is another object and
	// stays. For other constraints (e.g., PRIMARY KEY), also remove the corresponding index.
	if !isForeignKey {
		for i, idx := range tbl.Indexes {
			if toLower(idx.Name) == key {
				tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
				break
			}
		}
	}

	return nil
}

// alterRenameColumn changes a column name in-place.
func (c *Catalog) alterRenameColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	oldKey := toLower(cmd.Name)
	idx, exists := tbl.colByName[oldKey]
	if !exists {
		return errNoSuchColumn(cmd.Name, tbl.Name)
	}

	newKey := toLower(cmd.NewName)
	if newKey != oldKey {
		if _, dup := tbl.colByName[newKey]; dup {
			return errDupColumn(cmd.NewName)
		}
	}

	tbl.Columns[idx].Name = cmd.NewName
	updateColumnRefsInIndexes(tbl, cmd.Name, cmd.NewName)
	rebuildColIndex(tbl)
	return nil
}

// alterRenameIndex changes an index name in-place.
func (c *Catalog) alterRenameIndex(tbl *Table, cmd *nodes.AlterTableCmd) error {
	oldKey := toLower(cmd.Name)
	newKey := toLower(cmd.NewName)

	if newKey != oldKey {
		if err := validateNonPrimaryIndexName(cmd.NewName); err != nil {
			return err
		}
	}
	if newKey != oldKey && indexNameExists(tbl, cmd.NewName) {
		return errDupKeyName(cmd.NewName)
	}

	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == oldKey {
			idx.Name = cmd.NewName
			if newKey != oldKey {
				renameFunctionalIndexHiddenColumns(tbl, idx, cmd.NewName)
			}
			// Also update any constraint that references this index.
			for _, con := range tbl.Constraints {
				if toLower(con.IndexName) == oldKey {
					con.IndexName = cmd.NewName
					con.Name = cmd.NewName
				}
			}
			return nil
		}
	}

	return &Error{
		Code:     ErrDupKeyName,
		SQLState: sqlState(ErrDupKeyName),
		Message:  fmt.Sprintf("Key '%s' doesn't exist in table '%s'", cmd.Name, tbl.Name),
	}
}

// alterRenameTable moves a table to a new name.
func (c *Catalog) alterRenameTable(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	newTable := cmd.NewTable
	if newTable == nil {
		newTable = &nodes.TableRef{Name: cmd.NewName}
	}
	newDB, err := c.resolveDatabase(newTable.Schema)
	if err != nil {
		return err
	}
	newName := newTable.Name
	// The table leaves the database it is in now, which an earlier rename in
	// the same statement may already have changed: db is where the statement
	// started. Deleting from db instead left the intermediate destination's
	// entry behind, so `ALTER TABLE d.t RENAME TO e.u, RENAME TO e.v` would
	// have left both e.u and e.v pointing at the table; OceanBase CE 4.3.5.6 runs
	// that statement and ends with e.v alone (read on 2026-11-04).
	oldDB := db
	if tbl.Database != nil {
		oldDB = tbl.Database
	}
	newKey := toLower(newName)
	oldKey := toLower(tbl.Name)

	if (newDB != oldDB || newKey != oldKey) && newDB.Tables[newKey] != nil {
		return errDupTable(newName)
	}

	delete(oldDB.Tables, oldKey)
	tbl.Name = newName
	tbl.Database = newDB
	newDB.Tables[newKey] = tbl
	return nil
}

// alterTableOption applies a table option (ENGINE, CHARSET, etc.).
func (c *Catalog) alterTableOption(tbl *Table, cmd *nodes.AlterTableCmd) error {
	opt := cmd.Option
	if opt == nil {
		return nil
	}

	switch toLower(opt.Name) {
	case "table_id":
		// OceanBase parses ALTER TABLE t TABLE_ID = n and refuses to run it: the
		// pinned server answers 1064 on a table that exists and 1146 on one that
		// does not (testdata/simulate.tsv, rule table_option/TABLE_ID). A table's
		// id is given when it is created and never altered.
		return &Error{Code: 1064, SQLState: "42000", Message: "TABLE_ID cannot be altered"}
	case "engine":
		tbl.Engine = opt.Value
	case "charset", "character set", "default charset", "default character set":
		tbl.Charset = normalizeCharsetName(opt.Value)
		// Update collation to the default for this charset.
		if defColl, ok := defaultCollationForCharset[toLower(tbl.Charset)]; ok {
			tbl.Collation = defColl
		}
	case "collate", "default collate":
		tbl.Collation = normalizeCollationName(opt.Value)
	case "comment":
		tbl.Comment = opt.Value
	case "auto_increment":
		fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
	case "row_format":
		tbl.RowFormat = opt.Value
	case "key_block_size":
		fmt.Sscanf(opt.Value, "%d", &tbl.KeyBlockSize)
	case "compression":
		if err := validateCompressor(opt.Value); err != nil {
			return err
		}
		tbl.Compression = opt.Value
	case "encryption":
		tbl.Encryption = opt.Value
	case "stats_persistent":
		tbl.StatsPersistent = opt.Value
	case "stats_auto_recalc":
		tbl.StatsAutoRecalc = opt.Value
	case "stats_sample_pages":
		tbl.StatsSamplePages = opt.Value
	case "min_rows":
		tbl.MinRows = opt.Value
	case "max_rows":
		tbl.MaxRows = opt.Value
	case "avg_row_length":
		tbl.AvgRowLength = opt.Value
	case "tablespace":
		tbl.Tablespace = opt.Value
	case "pack_keys":
		tbl.PackKeys = opt.Value
	case "checksum":
		tbl.Checksum = opt.Value
	case "delay_key_write":
		tbl.DelayKeyWrite = opt.Value
	}
	return nil
}

// alterColumnDefault sets or drops the default on an existing column.
func (c *Catalog) alterColumnDefault(tbl *Table, cmd *nodes.AlterTableCmd) error {
	// SET DEFAULT takes an expression only in parentheses, which reach here unwrapped:
	// one that reads a column is 1235 before the column is looked up.
	if cmd.DefaultExpr != nil && hasColumnRef(cmd.DefaultExpr) {
		return errNotSupportedDefault()
	}
	colKey := toLower(cmd.Name)
	idx, exists := tbl.colByName[colKey]
	if !exists {
		return errNoSuchColumn(cmd.Name, tbl.Name)
	}

	col := tbl.Columns[idx]
	if cmd.DefaultExpr != nil {
		if err := checkDefault(col.Name, col.typeNode, cmd.DefaultExpr); err != nil {
			return err
		}
		if err := validateDefaultExprForExistingColumn(col, cmd.DefaultExpr); err != nil {
			return err
		}
		setColumnDefaultFromExpr(col, cmd.DefaultExpr)
	} else {
		// DROP DEFAULT — MySQL shows no default at all (not even DEFAULT NULL).
		col.Default = nil
		col.DefaultKind = ColumnDefaultNone
		col.DefaultDropped = true
	}
	return nil
}

// alterIndexVisibility toggles the Visible flag on an index.
func (c *Catalog) alterIndexVisibility(tbl *Table, cmd *nodes.AlterTableCmd, visible bool) error {
	key := toLower(cmd.Name)
	for _, idx := range tbl.Indexes {
		// OceanBase does not reach the primary key through ALTER INDEX: on 4.3.5.6,
		// ALTER TABLE t ALTER INDEX `PRIMARY` {VISIBLE | INVISIBLE} answers 1176, Key
		// 'PRIMARY' doesn't exist, where MySQL answers 3522 for INVISIBLE. Setting it
		// invisible here made a key the engine never hides.
		if idx.Primary {
			continue
		}
		if toLower(idx.Name) == key {
			idx.Visible = visible
			return nil
		}
	}
	return &Error{
		Code:     ErrDupKeyName,
		SQLState: sqlState(ErrDupKeyName),
		Message:  fmt.Sprintf("Key '%s' doesn't exist in table '%s'", cmd.Name, tbl.Name),
	}
}

// insertColumn inserts col into tbl at the position specified by first/after.
// If neither first nor after is set, appends at end. Always rebuilds the column index.
func insertColumn(tbl *Table, col *Column, first bool, after string) error {
	if first {
		tbl.Columns = append([]*Column{col}, tbl.Columns...)
	} else if after != "" {
		afterIdx, ok := tbl.colByName[toLower(after)]
		if !ok {
			return errNoSuchColumn(after, tbl.Name)
		}
		pos := afterIdx + 1
		tbl.Columns = append(tbl.Columns, nil)
		copy(tbl.Columns[pos+1:], tbl.Columns[pos:])
		tbl.Columns[pos] = col
	} else {
		tbl.Columns = append(tbl.Columns, col)
	}
	rebuildColIndex(tbl)
	return nil
}

// rebuildColIndex rebuilds tbl.colByName and updates Position fields.
func rebuildColIndex(tbl *Table) {
	tbl.colByName = make(map[string]int, len(tbl.Columns))
	for i, col := range tbl.Columns {
		col.Position = i + 1
		tbl.colByName[toLower(col.Name)] = i
	}
}

// buildColumnFromDef builds a catalog Column from an AST ColumnDef.
func buildColumnFromDef(tbl *Table, colDef *nodes.ColumnDef) *Column {
	col := &Column{
		Name:     colDef.Name,
		Nullable: true,
	}

	// Type info.
	if colDef.TypeName != nil {
		col.typeNode = colDef.TypeName
		col.DataType = normalizedColumnDataType(colDef.TypeName)
		col.ColumnType = formatColumnType(colDef.TypeName)
		if isNationalStringType(colDef.TypeName.Name) {
			// As in CREATE TABLE: the national character set, named as CHARACTER SET is.
			col.Charset = "utf8mb4"
			col.CharsetExplicit = true
		}
		if colDef.TypeName.Charset != "" {
			col.Charset = normalizeCharsetName(colDef.TypeName.Charset)
			col.CharsetExplicit = true
		}
		if colDef.TypeName.Collate != "" {
			col.Collation = normalizeCollationName(colDef.TypeName.Collate)
			col.CollationExplicit = true
			if col.Charset == "" {
				col.Charset = normalizeCharsetName(charsetForCollation(col.Collation))
			}
		}
	}

	// Default charset/collation for string types.
	if isStringType(col.DataType) {
		if col.Charset == "" {
			col.Charset = tbl.Charset
		}
		if col.Collation == "" {
			// A column that names its character set, or has another than its table's,
			// takes that character set's default collation, not the table's collation.
			if col.CharsetExplicit || !strings.EqualFold(col.Charset, tbl.Charset) {
				if dc, ok := defaultCollationForCharset[toLower(col.Charset)]; ok {
					col.Collation = dc
				}
			} else {
				col.Collation = tbl.Collation
			}
		}
		applyBinaryModifierCollation(col, colDef.TypeName)
	}

	// Top-level column properties.
	if colDef.AutoIncrement {
		col.AutoIncrement = true
		col.Nullable = false
	}
	if colDef.Comment != "" {
		col.Comment = colDef.Comment
	}
	// SRID, copied by PRESENCE (HasSRID) so an explicit `SRID 0` on an ALTER ADD/MODIFY
	// COLUMN is preserved and distinct from a no-SRID column (mirrors tablecmds.go).
	if colDef.TypeName != nil && colDef.TypeName.HasSRID {
		col.SRID = colDef.TypeName.SRID
		col.HasSRID = true
	}
	if colDef.DefaultValue != nil {
		setColumnDefaultFromExpr(col, colDef.DefaultValue)
	}
	if colDef.OnUpdate != nil {
		setColumnOnUpdateFromExpr(col, colDef.OnUpdate)
	}
	if colDef.Generated != nil {
		// Generated expressions are stored by the engine in resolver-rewritten
		// form — TIMESTAMPADD/DATE_ADD fold to interval arithmetic, NOT folds
		// to inverted comparisons — identical to view bodies (oracle 8.0.32:
		// AS (TIMESTAMPADD(HOUR,1,a)) stores as ((`a` + interval 1 hour)),
		// AS (NOT a) as ((0 = `a`))). Apply the same rewrite pass before
		// rendering so user-declared SDL canonicalizes to the stored form.
		col.Generated = &GeneratedColumnInfo{
			Expr:   nodeToSQLGenerated(deparse.RewriteExpr(colDef.Generated.Expr), tbl.Charset),
			Stored: colDef.Generated.Stored,
		}
	}

	// Process column-level constraints (non-index-producing ones).
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull:
			col.Nullable = false
			col.NullExplicit = true
		case nodes.ColConstrNull:
			col.Nullable = true
			col.NullExplicit = true
		case nodes.ColConstrDefault:
			if cc.Expr != nil {
				setColumnDefaultFromExpr(col, cc.Expr)
			}
		case nodes.ColConstrAutoIncrement:
			col.AutoIncrement = true
			col.Nullable = false
		case nodes.ColConstrVisible:
			col.Invisible = false
		case nodes.ColConstrInvisible:
			col.Invisible = true
		case nodes.ColConstrCollate:
			if cc.Expr != nil {
				if s, ok := cc.Expr.(*nodes.StringLit); ok {
					col.Collation = normalizeCollationName(s.Value)
					col.CollationExplicit = true
				}
			}
		}
	}

	return col
}

func tableHasOtherAutoIncrementColumn(tbl *Table, exceptName string) bool {
	exceptKey := toLower(exceptName)
	for _, col := range tbl.Columns {
		if col.AutoIncrement && (exceptKey == "" || toLower(col.Name) != exceptKey) {
			return true
		}
	}
	return false
}

func validateDefaultExprForExistingColumn(col *Column, expr nodes.ExprNode) error {
	if isNullDefault(expr) && !col.Nullable {
		return errInvalidDefault(col.Name)
	}
	if err := validateDefaultOfType(col.Name, col.DataType, expr); err != nil {
		return err
	}
	return validateTemporalExprForColumn(col.Name, col.DataType, columnTemporalFSP(col), expr, true)
}

func columnTemporalFSP(col *Column) int {
	if col == nil || (col.DataType != "datetime" && col.DataType != "timestamp" && col.DataType != "time") {
		return 0
	}
	open := strings.IndexByte(col.ColumnType, '(')
	close := strings.IndexByte(col.ColumnType, ')')
	if open < 0 || close <= open+1 {
		return 0
	}
	n := 0
	for _, ch := range col.ColumnType[open+1 : close] {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

// updateColumnRefsInIndexes updates index and constraint column references
// when a column is renamed.
func updateColumnRefsInIndexes(tbl *Table, oldName, newName string) {
	oldKey := toLower(oldName)
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if toLower(ic.Name) == oldKey {
				ic.Name = newName
			}
		}
		for i, name := range idx.Storing {
			if toLower(name) == oldKey {
				idx.Storing[i] = newName
			}
		}
	}
	for _, con := range tbl.Constraints {
		for i, col := range con.Columns {
			if toLower(col) == oldKey {
				con.Columns[i] = newName
			}
		}
	}
	renamePartitionKeyColumn(tbl, oldName, newName)
}

// alterCheckEnforced toggles the ENFORCED / NOT ENFORCED flag on a CHECK constraint.
func (c *Catalog) alterCheckEnforced(tbl *Table, cmd *nodes.AlterTableCmd) error {
	key := toLower(cmd.Name)
	for _, con := range tbl.Constraints {
		if toLower(con.Name) == key && con.Type == ConCheck {
			con.NotEnforced = (cmd.NewName == "NOT ENFORCED")
			return nil
		}
	}
	return &Error{
		Code:     3940,
		SQLState: "HY000",
		Message:  fmt.Sprintf("Constraint '%s' does not exist.", cmd.Name),
	}
}

// alterConvertCharset handles CONVERT TO CHARACTER SET charset [COLLATE collation].
// This changes the table's default charset/collation AND converts all existing
// string columns to the new charset/collation.
func (c *Catalog) alterConvertCharset(tbl *Table, cmd *nodes.AlterTableCmd) error {
	charset := cmd.Name
	collation := cmd.NewName

	// If no collation specified, use the default collation for the charset.
	if collation == "" {
		if defColl, ok := defaultCollationForCharset[toLower(charset)]; ok {
			collation = defColl
		}
	}

	tbl.Charset = charset
	tbl.Collation = collation

	// Convert all string-type columns to the new charset/collation.
	for _, col := range tbl.Columns {
		if isStringType(col.DataType) {
			col.Charset = charset
			col.Collation = collation
		}
	}

	return nil
}

// Ensure strings import is used (for toLower references via strings package).
var _ = strings.ToLower

// alterAddPartition adds partition definitions to a partitioned table.
func (c *Catalog) alterAddPartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	if cmd.Number > 0 && len(cmd.PartitionDefs) == 0 {
		start := len(tbl.Partitioning.Partitions)
		if start == 0 {
			start = tbl.Partitioning.NumParts
		}
		for i := 0; i < cmd.Number; i++ {
			tbl.Partitioning.Partitions = append(tbl.Partitioning.Partitions, &PartitionDefInfo{
				Name: fmt.Sprintf("p%d", start+i),
			})
		}
		tbl.Partitioning.NumParts = len(tbl.Partitioning.Partitions)
		return nil
	}
	if tbl.Partitioning.SubType != "" {
		for _, pd := range cmd.PartitionDefs {
			if err := checkSubpartitionBoundKinds(partitionTypeFromString(tbl.Partitioning.SubType), pd.SubPartitions); err != nil {
				return err
			}
		}
	}
	for _, pd := range cmd.PartitionDefs {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, partitionTypeFromString(tbl.Partitioning.Type))
		}
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			case "engine":
				pdi.Engine = opt.Value
			case "comment":
				pdi.Comment = opt.Value
			}
		}
		addedPartitionSubpartitions(tbl.Partitioning, pdi, pd)
		tbl.Partitioning.Partitions = append(tbl.Partitioning.Partitions, pdi)
	}
	if len(cmd.PartitionDefs) > 0 {
		tbl.Partitioning.NumParts = len(tbl.Partitioning.Partitions)
	}
	return nil
}

// addedPartitionSubpartitions gives a partition that ALTER TABLE adds to a
// subpartitioned table its subpartitions, as OceanBase 4.3.5.6 does
// (addPartitionScenarios):
//
//   - the ones the partition defines. SHOW CREATE TABLE then prints every
//     partition's subpartitions rather than a template or a count, so the table
//     stops rendering as SUBPARTITIONS N; the template itself survives, and a
//     partition added after that without subpartitions still takes it;
//   - else the template's, named as templateSubpartitions names them, which
//     covers SUBPARTITIONS N;
//   - else one, <partition>sp0, which is what the engine gives a table whose
//     partitions each listed their own.
func addedPartitionSubpartitions(pi *PartitionInfo, pdi *PartitionDefInfo, pd *nodes.PartitionDef) {
	if pi.SubType == "" {
		return
	}
	switch {
	case len(pd.SubPartitions) > 0:
		for _, spd := range pd.SubPartitions {
			spdi := &SubPartitionDefInfo{Name: spd.Name, Engine: pdi.Engine, Comment: pdi.Comment}
			if spd.Values != nil {
				spdi.ValueExpr = partitionValueToString(spd.Values, partitionTypeFromString(pi.SubType))
			}
			for _, opt := range spd.Options {
				switch toLower(opt.Name) {
				case "engine":
					spdi.Engine = opt.Value
				case "comment":
					spdi.Comment = opt.Value
				}
			}
			pdi.SubPartitions = append(pdi.SubPartitions, spdi)
		}
		pi.UseDefaultSubpartitions = false
		pi.UseDefaultNumSubpartitions = false
	case len(pi.SubPartTemplate) > 0:
		pdi.SubPartitions = templateSubpartitions(pi.SubPartTemplate, pdi)
	default:
		pdi.SubPartitions = []*SubPartitionDefInfo{{Name: pdi.Name + "sp0", ValueExpr: defaultSubpartitionValue(pi.SubType)}}
	}
}

// alterDropPartition drops named partitions from a partitioned table.
func (c *Catalog) alterDropPartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	dropSet := make(map[string]bool)
	for _, name := range cmd.PartitionNames {
		dropSet[toLower(name)] = true
	}
	var remaining []*PartitionDefInfo
	for _, pd := range tbl.Partitioning.Partitions {
		if !dropSet[toLower(pd.Name)] {
			remaining = append(remaining, pd)
		}
	}
	tbl.Partitioning.Partitions = remaining
	return nil
}

// alterTruncatePartition truncates named partitions (no-op for metadata catalog).
func (c *Catalog) alterTruncatePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	// Truncate is a data operation; for DDL catalog purposes, it's a no-op.
	return nil
}

// alterCoalescePartition reduces the number of partitions for HASH/KEY partitioned tables.
func (c *Catalog) alterCoalescePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	// Determine current partition count.
	currentCount := len(tbl.Partitioning.Partitions)
	if currentCount == 0 {
		currentCount = tbl.Partitioning.NumParts
	}
	newCount := currentCount - cmd.Number
	if newCount < 1 {
		newCount = 1
	}
	if len(tbl.Partitioning.Partitions) > 0 {
		tbl.Partitioning.Partitions = tbl.Partitioning.Partitions[:newCount]
	}
	tbl.Partitioning.NumParts = newCount
	return nil
}

// alterReorganizePartition reorganizes partitions into new definitions.
func (c *Catalog) alterReorganizePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	// OceanBase refuses REORGANIZE PARTITION on a subpartitioned table, splitting or
	// merging, templated or not: 4.3.5.6 answers 1235 "split subpartition is not
	// supported". Reorganizing it here would invent partitions the engine never makes,
	// and used to drop their subpartitions besides.
	if tbl.Partitioning.SubType != "" {
		return &Error{Code: 1235, SQLState: "42000", Message: "split subpartition is not supported"}
	}
	// It reorganizes a table only by its primary key: the table has one, and the
	// partitioning key is the key's leading columns, in the key's order and with no
	// expression over them (TestOracle_OceanBaseRefusals).
	if err := reorganizeNeedsPrimaryKeyPrefix(tbl); err != nil {
		return err
	}
	// Remove the old partitions.
	dropSet := make(map[string]bool)
	for _, name := range cmd.PartitionNames {
		dropSet[toLower(name)] = true
	}
	var remaining []*PartitionDefInfo
	insertPos := -1
	for i, pd := range tbl.Partitioning.Partitions {
		if dropSet[toLower(pd.Name)] {
			if insertPos < 0 {
				insertPos = i
			}
			continue
		}
		remaining = append(remaining, pd)
	}
	if insertPos < 0 {
		insertPos = len(remaining)
	}

	// Build new partitions.
	var newParts []*PartitionDefInfo
	for _, pd := range cmd.PartitionDefs {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, partitionTypeFromString(tbl.Partitioning.Type))
		}
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			case "engine":
				pdi.Engine = opt.Value
			case "comment":
				pdi.Comment = opt.Value
			}
		}
		newParts = append(newParts, pdi)
	}

	// Insert new partitions at the position of the first removed partition.
	result := make([]*PartitionDefInfo, 0, len(remaining)+len(newParts))
	for i, pd := range remaining {
		if i == insertPos {
			result = append(result, newParts...)
		}
		result = append(result, pd)
	}
	if insertPos >= len(remaining) {
		result = append(result, newParts...)
	}
	tbl.Partitioning.Partitions = result
	return nil
}

// alterExchangePartition exchanges a partition with a non-partitioned table.
func (c *Catalog) alterExchangePartition(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errNotPartitioned()
	}
	// For DDL catalog purposes, exchange is primarily a data operation.
	// We just validate both tables exist.
	if cmd.ExchangeTable != nil {
		exchDB := db
		if cmd.ExchangeTable.Schema != "" {
			exchDB = c.GetDatabase(cmd.ExchangeTable.Schema)
			if exchDB == nil {
				return errNoSuchTable(cmd.ExchangeTable.Schema, cmd.ExchangeTable.Name)
			}
		}
		exchTbl := exchDB.GetTable(cmd.ExchangeTable.Name)
		if exchTbl == nil {
			return errNoSuchTable(exchDB.Name, cmd.ExchangeTable.Name)
		}
	}
	return nil
}

// partitionTypeFromString converts a string partition type to AST PartitionType.
func partitionTypeFromString(t string) nodes.PartitionType {
	switch t {
	case "RANGE", "RANGE COLUMNS":
		return nodes.PartitionRange
	case "LIST", "LIST COLUMNS":
		return nodes.PartitionList
	case "HASH":
		return nodes.PartitionHash
	case "KEY":
		return nodes.PartitionKey
	default:
		return nodes.PartitionRange
	}
}

// generatedExprReferencesColumn checks if a generated column expression
// references a column by name. The expression uses backtick-quoted identifiers
// (e.g., `col_name`), so we search for the backtick-quoted form.
func generatedExprReferencesColumn(expr, colName string) bool {
	target := "`" + strings.ToLower(colName) + "`"
	return strings.Contains(strings.ToLower(expr), target)
}

// reorganizeNeedsPrimaryKeyPrefix refuses REORGANIZE PARTITION as OceanBase does on a table
// whose primary key does not start with the partitioning key.
func reorganizeNeedsPrimaryKeyPrefix(tbl *Table) error {
	var pk *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			pk = idx
		}
	}
	if pk == nil {
		return &Error{Code: 1235, SQLState: "0A000", Message: "reorganizing table without primary key(s) is not supported"}
	}
	key := tbl.Partitioning.Columns
	if len(key) == 0 {
		if column, ok := partitionKeyColumn(tbl.Partitioning.Expr); ok {
			key = []string{column}
		}
	}
	matches := len(key) > 0 && len(key) <= len(pk.Columns)
	for i := 0; matches && i < len(key); i++ {
		matches = toLower(pk.Columns[i].Name) == toLower(key[i])
	}
	if !matches {
		return &Error{Code: 1235, SQLState: "0A000", Message: "mismatching between primary key prefix and partition key is not supported"}
	}
	return nil
}

// partitionKeyColumn returns the column a partitioning expression is, when it is one column
// and nothing else: a bare name, or a name in backticks, which may hold anything and writes
// a backtick of its own twice.
func partitionKeyColumn(expr string) (string, bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "", false
	}
	if expr[0] != '`' {
		for _, r := range expr {
			if r != '_' && r != '$' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return "", false
			}
		}
		return expr, true
	}
	var name strings.Builder
	for i := 1; i < len(expr); i++ {
		switch {
		case expr[i] != '`':
			name.WriteByte(expr[i])
		case i+1 < len(expr) && expr[i+1] == '`':
			name.WriteByte('`')
			i++
		default:
			// The closing backtick: the expression is a column only if it ends here.
			if i == len(expr)-1 && name.Len() > 0 {
				return name.String(), true
			}
			return "", false
		}
	}
	return "", false
}
