package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

func (c *Catalog) alterTable(stmt *nodes.AlterTableStmt) error {
	// The checks MySQL makes while preparing the statement answer ahead of the
	// table it names and of the database: mysql:8.0.46 answers 1439 for
	// `ALTER TABLE nosuchdb.a ADD a INT(256)`.
	if err := c.alterDefinitions(stmt.Commands); err != nil {
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
	// EXCHANGE PARTITION naming the table being altered as the table to swap with
	// is 1066, "Not unique table/alias", which MySQL answers before it looks
	// either table up: mysql:8.0.46 answers it for a table that exists and for
	// one that does not.
	if err := alterExchangeSelf(c, db, stmt); err != nil {
		return err
	}

	tableName := stmt.Table.Name
	key := toLower(tableName)
	tbl := db.Tables[key]
	if tbl == nil {
		return errNoSuchTable(db.Name, tableName)
	}
	// PERFORMANCE_SCHEMA holds the performance schema's own tables: on a table
	// that exists MySQL answers 1031 for it, after the prepare-time checks and
	// before the commands (mysql:8.0.46).
	if err := validateTableEngineReplaceable(tbl, stmt.Commands); err != nil {
		return err
	}
	// A column the statement adds is not one it can drop, rename or change:
	// MySQL resolves those against the table the statement found, so
	// `ALTER TABLE a ADD COLUMN x INT, DROP COLUMN x` answers 1091 and
	// `ADD COLUMN x INT, CHANGE x y INT` 1054 (mysql:8.0.46). Both answers come
	// before the key checks the commands make.
	if err := validateAlterColumnSources(tbl, stmt.Commands); err != nil {
		return err
	}

	repartition := alterRepartition(stmt.Commands)
	if tbl.Partitioning == nil && repartition == nil {
		return c.alterTableCommands(db, tbl, key, stmt)
	}
	// A unique key on a partitioned table must cover the partitioning columns. MySQL
	// checks the table the statement leaves, after every other error (1060, 1061 and
	// 1091 come first in ADD UNIQUE (b), ADD COLUMN a ...), and leaves the table as it
	// was: snapshot it to undo the commands when the check fails.
	before := make(map[*Index]bool, len(tbl.Indexes))
	for _, idx := range tbl.Indexes {
		before[idx] = true
	}
	snapshot := cloneTable(tbl)
	if err := c.alterTableCommands(db, tbl, key, stmt); err != nil {
		return err
	}
	if err := validateAlteredUniqueKeys(tbl, before, repartition); err != nil {
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

// validateTableEngineReplaceable refuses an ENGINE the table cannot take:
// PERFORMANCE_SCHEMA, which holds the performance schema's own tables (1031).
func validateTableEngineReplaceable(tbl *Table, commands []*nodes.AlterTableCmd) error {
	for _, cmd := range commands {
		if cmd == nil || cmd.Type != nodes.ATTableOption || cmd.Option == nil || !engineOptionName(cmd.Option.Name) {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(cmd.Option.Value), performanceSchemaEngine) {
			return errTableStorageEngineOption(tbl.Name)
		}
	}
	return nil
}

// validateAlterColumnSources refuses a command naming a column the table the
// statement found does not hold, in the statement's order: a DROP COLUMN is
// 1091 "Can't DROP 'x'"; the column a MODIFY, CHANGE or RENAME starts from is
// 1054 "Unknown column 'x' in 'tbl'". A column an earlier command of the same
// statement adds is not one of them: MySQL resolves these against the table the
// statement found.
func validateAlterColumnSources(tbl *Table, commands []*nodes.AlterTableCmd) error {
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATDropColumn:
			if cmd.IfExists {
				continue
			}
			if tbl.GetColumn(cmd.Name) == nil {
				return errCantDropKey(cmd.Name)
			}
		case nodes.ATModifyColumn, nodes.ATRenameColumn, nodes.ATChangeColumn:
			if cmd.Name != "" && tbl.GetColumn(cmd.Name) == nil {
				return errNoSuchColumn(cmd.Name, tbl.Name)
			}
		}
	}
	return nil
}

// alterFinalColumns is the set of columns the statement leaves, which MySQL
// judges the keys it adds against: `ALTER TABLE a ADD KEY(c2), ADD c2 INT` is
// taken, and a key on a column the statement drops answers 1072
// (mysql:8.0.46).
func alterFinalColumns(tbl *Table, commands []*nodes.AlterTableCmd) map[string]bool {
	cols := make(map[string]bool, len(tbl.Columns))
	for _, col := range tbl.Columns {
		cols[toLower(col.Name)] = true
	}
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATAddColumn:
			for _, colDef := range alterColumnDefs(cmd) {
				if colDef != nil {
					cols[toLower(colDef.Name)] = true
				}
			}
		case nodes.ATDropColumn:
			delete(cols, toLower(cmd.Name))
		case nodes.ATRenameColumn:
			delete(cols, toLower(cmd.Name))
			cols[toLower(cmd.NewName)] = true
		case nodes.ATChangeColumn:
			if cmd.Column != nil {
				delete(cols, toLower(cmd.Name))
				cols[toLower(cmd.Column.Name)] = true
			}
		}
	}
	return cols
}

// alterTableCommands runs the commands of an ALTER TABLE on tbl, found under key in db.
func (c *Catalog) alterTableCommands(db *Database, tbl *Table, key string, stmt *nodes.AlterTableStmt) error {
	if len(stmt.Commands) == 0 {
		return nil
	}
	leaves := alterFinalColumns(tbl, stmt.Commands)
	if len(stmt.Commands) == 1 && !addNeedsRollback(stmt.Commands[0]) {
		// Single command carrying one table element: no rollback needed.
		return c.execAlterCmd(db, tbl, stmt.Commands[0], leaves)
	}

	// Multi-command ALTER: MySQL treats this as atomic.
	// Snapshot the table so we can rollback on any sub-command failure.
	snapshot := cloneTable(tbl)
	origKey := key

	for _, cmd := range stmt.Commands {
		if err := c.execAlterCmd(db, tbl, cmd, leaves); err != nil {
			undoAlter(db, origKey, tbl, snapshot)
			return err
		}
	}
	if err := checkKeysOnFinalTable(tbl, stmt.Commands); err != nil {
		undoAlter(db, origKey, tbl, snapshot)
		return err
	}
	// Clear transient cleanup tracking after successful multi-command ALTER.
	tbl.droppedByCleanup = nil
	return nil
}

// checkKeysOnFinalTable re-checks the keys a statement added once every command
// has run. MySQL judges a key against the table the statement leaves, where a
// column a later command adds carries the definition the check needs: a TEXT
// part without a length is 1170, a prefix longer than the column 1089, and a
// SPATIAL key's part must be NOT NULL (1252). A primary key also makes the
// columns it names NOT NULL, which its own command could not do while a later
// column was still missing (measured on mysql:8.0.46).
func checkKeysOnFinalTable(tbl *Table, commands []*nodes.AlterTableCmd) error {
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		for _, con := range alterConstraints(cmd) {
			if con == nil {
				continue
			}
			switch con.Type {
			case nodes.ConstrPrimaryKey:
				cols := extractColumnNames(con)
				if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
					return err
				}
				for _, name := range cols {
					if col := tbl.GetColumn(name); col != nil {
						col.Nullable = false
					}
				}
			case nodes.ConstrUnique, nodes.ConstrIndex:
				if err := validateIndexColumns(tbl, buildIndexColumns(con), false, false); err != nil {
					return err
				}
			case nodes.ConstrFulltextIndex:
				if err := validateIndexColumns(tbl, buildIndexColumns(con), true, false); err != nil {
					return err
				}
			case nodes.ConstrSpatialIndex:
				if err := validateIndexColumns(tbl, buildIndexColumns(con), false, true); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// alterExchangeSelf refuses an EXCHANGE PARTITION that names the table being
// altered as the table to swap with: 1066 "Not unique table/alias", answered
// before either table is looked up (measured on mysql:8.0.46). A name in another
// database is a different table, and that database is looked up with the table.
func alterExchangeSelf(c *Catalog, db *Database, stmt *nodes.AlterTableStmt) error {
	if stmt.Table == nil {
		return nil
	}
	for _, cmd := range stmt.Commands {
		if cmd == nil || cmd.Type != nodes.ATExchangePartition || cmd.ExchangeTable == nil {
			continue
		}
		exchDB := db
		if cmd.ExchangeTable.Schema != "" {
			exchDB = c.GetDatabase(cmd.ExchangeTable.Schema)
			if exchDB == nil {
				continue
			}
		}
		if exchDB == db && strings.EqualFold(cmd.ExchangeTable.Name, stmt.Table.Name) {
			return errNonUniqueTable(cmd.ExchangeTable.Name)
		}
	}
	return nil
}

// alterDefinitions runs the checks MySQL makes before it opens the table an
// ALTER TABLE names, so each answers ahead of 1146: a FULLTEXT key's parser
// (1128), which its parser looks up, then, in the statement's order, a column's
// display width (1439) and an ENGINE the server does not ship (1286). Measured
// on mysql:8.0.46: `ALTER TABLE a ADD e INT(256), ENGINE = nosuch` answers 1439
// where `ALTER TABLE a ENGINE = nosuch, ADD e INT(256)` answers 1286.
func (c *Catalog) alterDefinitions(commands []*nodes.AlterTableCmd) error {
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		for _, con := range alterConstraints(cmd) {
			if err := validateFulltextParser(con); err != nil {
				return err
			}
		}
	}
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		for _, colDef := range alterColumnDefs(cmd) {
			if err := validateColumnDisplayWidth(colDef); err != nil {
				return err
			}
		}
		if cmd.Type == nodes.ATTableOption && cmd.Option != nil && engineOptionName(cmd.Option.Name) {
			if err := validateStorageEngine(cmd.Option.Value, c.noEngineSubstitution()); err != nil {
				return err
			}
		}
		// A partition definition's own errors answer ahead of the table it would
		// be added to, as ADD's do: `ALTER TABLE missing ADD PARTITION (PARTITION
		// p0 VALUES LESS THAN (b))` is 1054, not 1146 (mysql:8.0.46).
		if cmd.Type == nodes.ATAddPartition || cmd.Type == nodes.ATReorganizePartition {
			if err := validatePartitionDefinitions(cmd.PartitionDefs, c.noEngineSubstitution()); err != nil {
				return err
			}
		}
	}
	return nil
}

// validatePartitionDefinitions runs the checks MySQL makes on an ADD PARTITION's
// definitions before it judges the table they are added to: a name its bound
// reads as a column (1054), then an engine it does not ship (1286). Measured on
// mysql:8.0.46, both answer ahead of 1505.
func validatePartitionDefinitions(defs []*nodes.PartitionDef, refuseUnknownEngine bool) error {
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		if name, ok := partitionBoundColumn(pd.Values); ok {
			return errNoSuchColumn(name, "partition function")
		}
		if err := validatePartitionEngines(pd.Options, pd.SubPartitions, refuseUnknownEngine); err != nil {
			return err
		}
	}
	return nil
}

// validatePartitionEngines refuses an ENGINE a partition or subpartition names
// and the server does not ship: 1286. Every one of them carries an engine, so a
// valid parent engine does not stop its subpartitions from being judged.
func validatePartitionEngines(opts []*nodes.TableOption, subs []*nodes.SubPartitionDef, refuseUnknownEngine bool) error {
	for _, opt := range opts {
		if opt == nil || !engineOptionName(opt.Name) {
			continue
		}
		if err := validateStorageEngine(opt.Value, refuseUnknownEngine); err != nil {
			return err
		}
	}
	for _, sub := range subs {
		if sub == nil {
			continue
		}
		if err := validatePartitionEngines(sub.Options, nil, refuseUnknownEngine); err != nil {
			return err
		}
	}
	return nil
}

// alterColumnDefs returns the column definitions a command carries.
func alterColumnDefs(cmd *nodes.AlterTableCmd) []*nodes.ColumnDef {
	switch cmd.Type {
	case nodes.ATAddColumn:
		defs := make([]*nodes.ColumnDef, 0, len(cmd.Columns)+1)
		if cmd.Column != nil {
			defs = append(defs, cmd.Column)
		}
		return append(defs, cmd.Columns...)
	case nodes.ATModifyColumn, nodes.ATChangeColumn:
		if cmd.Column != nil {
			return []*nodes.ColumnDef{cmd.Column}
		}
	}
	return nil
}

// alterConstraints returns the constraints a command carries, from ADD
// CONSTRAINT and from the table element list `ADD (col, CONSTRAINT, ...)`.
func alterConstraints(cmd *nodes.AlterTableCmd) []*nodes.Constraint {
	cons := make([]*nodes.Constraint, 0, len(cmd.Constraints)+1)
	if cmd.Constraint != nil {
		cons = append(cons, cmd.Constraint)
	}
	return append(cons, cmd.Constraints...)
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

func (c *Catalog) execAlterCmd(db *Database, tbl *Table, cmd *nodes.AlterTableCmd, leaves map[string]bool) error {
	switch cmd.Type {
	case nodes.ATAddColumn:
		return c.alterAddColumn(tbl, cmd, leaves)
	case nodes.ATDropColumn:
		return c.alterDropColumn(tbl, cmd)
	case nodes.ATModifyColumn:
		return c.alterModifyColumn(tbl, cmd)
	case nodes.ATChangeColumn:
		return c.alterChangeColumn(tbl, cmd)
	case nodes.ATAddIndex, nodes.ATAddConstraint:
		return c.alterAddConstraint(tbl, cmd, leaves)
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
	case nodes.ATAlterColumnVisible:
		return c.alterColumnVisibility(tbl, cmd, false)
	case nodes.ATAlterColumnInvisible:
		return c.alterColumnVisibility(tbl, cmd, true)
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
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		tbl.Partitioning = nil
		return nil
	case nodes.ATRebuildPartition, nodes.ATDiscardPartitionTablespace, nodes.ATImportPartitionTablespace:
		// These leave the schema as it is, and a table without partitions is
		// still refused (1505, measured on mysql:8.0.46). ANALYZE, CHECK,
		// OPTIMIZE and REPAIR PARTITION on such a table are taken, so they stay
		// the default no-op.
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		return nil
	default:
		// Unsupported alter command; silently ignore.
		return nil
	}
}

// alterAddColumn adds a new column to the table.
func (c *Catalog) alterAddColumn(tbl *Table, cmd *nodes.AlterTableCmd, leaves map[string]bool) error {
	// The parenthesized form reads a table element list, so it carries the
	// constraints written beside the columns:
	// `ADD COLUMN (a1 INT, PRIMARY KEY (a1))` adds both. The columns come
	// first, as the CREATE TABLE path builds them, and MySQL 8.0.46 makes a1
	// the primary key for exactly that statement.
	if len(cmd.Columns) > 0 || len(cmd.Constraints) > 0 {
		for _, colDef := range cmd.Columns {
			if err := c.addSingleColumn(tbl, colDef, false, ""); err != nil {
				return err
			}
		}
		for _, con := range cmd.Constraints {
			if err := c.alterAddConstraint(tbl, &nodes.AlterTableCmd{
				Type:       nodes.ATAddConstraint,
				Constraint: con,
			}, leaves); err != nil {
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
	if err := validateColumnDefSemantics(colDef); err != nil {
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

		case nodes.ColConstrCheck:
			// A column's CHECK is a table constraint like any other, and the
			// CREATE TABLE path has always built one; the ALTER path dropped
			// it, so `ADD b INT CONSTRAINT positive_b CHECK (b > 0)` reported
			// success and installed only the column. MySQL 8.0.46 installs
			// the constraint -- `CONSTRAINT positive_b CHECK ((b > 0))` reads
			// back from SHOW CREATE TABLE -- and names an unnamed one
			// <table>_chk_<n>, continuing the table's run (read on
			// 2026-11-04).
			if cc.Expr == nil {
				break
			}
			conName := cc.Name
			if conName == "" {
				conName = fmt.Sprintf("%s_chk_%d", tbl.Name, nextCheckNumber(tbl))
			}
			if checkConstraintNameExists(tbl.Database, nil, conName) {
				return errCheckConstraintDupName(conName)
			}
			if err := validateColumnCheckReferences(colDef.Name, conName, cc.Expr); err != nil {
				return err
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:        conName,
				Type:        ConCheck,
				Table:       tbl,
				CheckExpr:   nodeToSQL(cc.Expr),
				NotEnforced: cc.NotEnforced,
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
func (c *Catalog) alterAddConstraint(tbl *Table, cmd *nodes.AlterTableCmd, leaves map[string]bool) error {
	con := cmd.Constraint
	if con == nil {
		return nil
	}

	cols := extractColumnNames(con)

	switch con.Type {
	case nodes.ConstrPrimaryKey:
		// Check for duplicate PK.
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				return errMultiplePriKey()
			}
		}
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
		}
		if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
			return err
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
		// Refused before the columns change: MySQL 8.0 leaves them nullable.
		if !pkIdx.Visible {
			return &Error{Code: 3522, SQLState: "HY000", Message: "A primary key index cannot be invisible"}
		}
		// Mark PK columns as NOT NULL.
		for _, colName := range cols {
			col := tbl.GetColumn(colName)
			if col != nil {
				col.Nullable = false
			}
		}
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
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
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
		coerceInnoDBHashIndex(tbl, uqIdx)
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
			conName = fmt.Sprintf("%s_ibfk_%d", tbl.Name, nextFKGeneratedNumber(tbl, tbl.Name))
		}
		if foreignKeyConstraintNameExists(tbl.Database, nil, conName) {
			return errFKDupName(conName)
		}
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
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
		// Add implicit backing index for FK if needed.
		ensureFKBackingIndex(tbl, fkIndexName(con), cols, buildIndexColumns(con))

	case nodes.ConstrCheck:
		conName := con.Name
		if conName == "" {
			conName = fmt.Sprintf("%s_chk_%d", tbl.Name, nextCheckNumber(tbl))
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
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
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
		// added by ALTER lost its comment, visibility and block size.
		applyIndexOptions(keyIdx, con.IndexOptions)
		coerceInnoDBHashIndex(tbl, keyIdx)
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
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, idxCols, true, false); err != nil {
			return err
		}
		ftIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			Fulltext:  true,
			IndexType: "FULLTEXT",
			Visible:   true,
		}
		applyIndexOptions(ftIdx, con.IndexOptions)
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
		if err := validateKeyColumnsExist(tbl, leaves, con); err != nil {
			return err
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, idxCols, false, true); err != nil {
			return err
		}
		if strings.EqualFold(resolveConstraintIndexType(con), "BTREE") {
			return &Error{Code: 3500, SQLState: "HY000", Message: "Index type not supported for spatial index"}
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

	// For FK constraints, MySQL keeps the backing index when dropping the FK.
	// For other constraints (e.g., PRIMARY KEY), also remove the corresponding index.
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
	// have left both e.u and e.v pointing at the table; MySQL 8.0.46 runs
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
	case "engine":
		// `ENGINE = 'DEFAULT'` names the server's configured default engine,
		// which is the catalog's empty engine, and MySQL switches the table to
		// it: a MyISAM table reads back InnoDB (measured on mysql:8.0.46). A
		// name the server substitutes for reads as that default too.
		tbl.Engine = storedEngineName(opt.Value)
	case "charset", "character set", "default charset", "default character set":
		tbl.Charset = normalizeCharsetName(opt.Value)
		// Update collation to the default for this charset.
		if defColl, ok := defaultCollationForCharset[toLower(tbl.Charset)]; ok {
			tbl.Collation = defColl
		}
	case "collate", "default collate":
		tbl.Collation = opt.Value
	case "comment":
		tbl.Comment = opt.Value
	case "auto_increment":
		fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
	case "row_format":
		tbl.RowFormat = opt.Value
	case "key_block_size":
		fmt.Sscanf(opt.Value, "%d", &tbl.KeyBlockSize)
	case "compression":
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
	case "checksum", "table_checksum":
		tbl.Checksum = normalizedChecksum(opt.Value)
	case "delay_key_write":
		tbl.DelayKeyWrite = opt.Value
	}
	return nil
}

// alterColumnDefault sets or drops the default on an existing column.
func (c *Catalog) alterColumnDefault(tbl *Table, cmd *nodes.AlterTableCmd) error {
	colKey := toLower(cmd.Name)
	idx, exists := tbl.colByName[colKey]
	if !exists {
		return errNoSuchColumn(cmd.Name, tbl.Name)
	}

	col := tbl.Columns[idx]
	if cmd.DefaultExpr != nil {
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

// alterColumnVisibility toggles the INVISIBLE flag on a column.
func (c *Catalog) alterColumnVisibility(tbl *Table, cmd *nodes.AlterTableCmd, invisible bool) error {
	colKey := toLower(cmd.Name)
	idx, exists := tbl.colByName[colKey]
	if !exists {
		return errNoSuchColumn(cmd.Name, tbl.Name)
	}
	tbl.Columns[idx].Invisible = invisible
	return nil
}

// alterIndexVisibility toggles the Visible flag on an index.
func (c *Catalog) alterIndexVisibility(tbl *Table, cmd *nodes.AlterTableCmd, visible bool) error {
	key := toLower(cmd.Name)
	for _, idx := range tbl.Indexes {
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
		col.DataType = normalizedColumnDataType(colDef.TypeName)
		col.ColumnType = formatColumnType(colDef.TypeName)
		if isNationalStringType(colDef.TypeName.Name) {
			col.Charset = "utf8mb3"
		}
		if colDef.TypeName.Charset != "" {
			col.Charset = normalizeCharsetName(colDef.TypeName.Charset)
			col.CharsetExplicit = true
		}
		if colDef.TypeName.Collate != "" {
			col.Collation = colDef.TypeName.Collate
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
			// If column charset differs from table charset, use the default
			// collation for the column's charset, not the table's collation.
			if !strings.EqualFold(col.Charset, tbl.Charset) {
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
		case nodes.ColConstrNotSecondary:
			col.NotSecondary = true
		case nodes.ColConstrCollate:
			if cc.Expr != nil {
				if s, ok := cc.Expr.(*nodes.StringLit); ok {
					col.Collation = s.Value
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
	if isLiteralDefaultForbiddenType(col.DataType) && isLiteralDefault(expr) {
		return &Error{Code: 1101, SQLState: "42000", Message: fmtCantHaveDefault(col.Name)}
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
		return errPartitionMgmtOnNonPartitioned()
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
	for _, pd := range cmd.PartitionDefs {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, partitionTypeFromString(tbl.Partitioning.Type))
		}
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			// STORAGE ENGINE is the other spelling of ENGINE, and names the
			// same engine.
			case "engine", "storage engine":
				pdi.Engine = storedEngineName(opt.Value)
			case "comment":
				pdi.Comment = opt.Value
			}
		}
		tbl.Partitioning.Partitions = append(tbl.Partitioning.Partitions, pdi)
	}
	if len(cmd.PartitionDefs) > 0 {
		tbl.Partitioning.NumParts = len(tbl.Partitioning.Partitions)
	}
	return nil
}

// alterDropPartition drops named partitions from a partitioned table.
func (c *Catalog) alterDropPartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
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
		return errPartitionMgmtOnNonPartitioned()
	}
	// Truncate is a data operation; for DDL catalog purposes, it's a no-op.
	return nil
}

// alterCoalescePartition reduces the number of partitions for HASH/KEY partitioned tables.
func (c *Catalog) alterCoalescePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
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
		return errPartitionMgmtOnNonPartitioned()
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
			// STORAGE ENGINE is the other spelling of ENGINE, and names the
			// same engine.
			case "engine", "storage engine":
				pdi.Engine = storedEngineName(opt.Value)
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
	// For DDL catalog purposes, exchange is primarily a data operation.
	// We just validate both tables exist. The exchange table is looked up before
	// the table being altered is judged: mysql:8.0.46 answers 1146 for
	// `ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE missing` on a table with no
	// partitions, and 1505 only when that table exists.
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
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
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

// addNeedsRollback reports whether an ADD action can fail after it has begun
// changing the table, so that it needs the rollback a multi-command ALTER
// gets rather than the single-command fast path. Two shapes can:
//
//   - `ADD [COLUMN] ( ... )` carries a table element list, so an element that
//     fails part-way through would leave the earlier ones installed;
//   - a column with constraints has them applied after the column itself is
//     inserted, and those can fail -- a second primary key, a duplicate
//     CHECK name, a CHECK naming another column.
//
// MySQL 8.0.46 leaves the table untouched for all of them:
// `ADD (b INT, KEY i (b))` against a table that already has i,
// `ADD b INT PRIMARY KEY` against a table that already has one, and
// `ADD b INT CONSTRAINT ck CHECK (b > 0)` against a table that already has ck
// (read on 2026-11-04). Any column constraint counts, rather than the three
// that fail today, so a validation added later cannot reopen this.
func addNeedsRollback(cmd *nodes.AlterTableCmd) bool {
	if cmd == nil || cmd.Type != nodes.ATAddColumn {
		return false
	}
	if len(cmd.Columns)+len(cmd.Constraints) > 1 {
		return true
	}
	if cmd.Column != nil && len(cmd.Column.Constraints) > 0 {
		return true
	}
	for _, col := range cmd.Columns {
		if col != nil && len(col.Constraints) > 0 {
			return true
		}
	}
	return false
}
