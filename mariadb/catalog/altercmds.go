package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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

	repartition := alterRepartition(stmt.Commands)
	if tbl.Partitioning == nil && repartition == nil {
		return c.alterTableCommands(db, tbl, key, stmt)
	}
	// A unique key on a partitioned table must cover the partitioning columns. MariaDB
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

// alterTableCommands runs the commands of an ALTER TABLE on tbl, found under key in db.
func (c *Catalog) alterTableCommands(db *Database, tbl *Table, key string, stmt *nodes.AlterTableStmt) error {
	if len(stmt.Commands) == 0 {
		return nil
	}

	// A column-level WITH/WITHOUT SYSTEM VERSIONING is valid only when the table
	// is already system-versioned at the start of the statement — a later
	// ADD SYSTEM VERSIONING in the same ALTER does not retroactively allow it.
	// Checked up front (against the pre-statement state) so command order does
	// not matter.
	if !tbl.SystemVersioned {
		for _, cmd := range stmt.Commands {
			if addsColumnVersioning(cmd) {
				return errNotSystemVersioned(tbl.Name)
			}
		}
	} else {
		// Ordinary column add/drop/modify/change is not allowed on a table that
		// is already system-versioned at statement start (MariaDB default; it
		// requires @@system_versioning_alter_history, which is not modeled).
		for _, cmd := range stmt.Commands {
			if isColumnMutation(cmd) {
				return errAlterColumnOnSystemVersioned(tbl.Name)
			}
		}
	}

	if err := checkAddIfNotExistsTwice(tbl, stmt.Commands); err != nil {
		return err
	}

	skip, err := checkAlterIndexIgnorability(tbl, stmt.Commands)
	if err != nil {
		return err
	}

	dropsPK := alterDropsPrimaryKey(stmt.Commands)

	if len(stmt.Commands) == 1 {
		cmd := stmt.Commands[0]
		if skip[cmd] {
			return nil
		}
		// A versioned table can be made inconsistent by any column change (a
		// dropped/modified/renamed ROW START/END or period column), so route the
		// command through post-mutation validation when the table is already
		// versioned or the command itself carries versioning metadata. So does
		// a table with an IGNORED index, which a dropped key or a column made
		// NOT NULL can promote to the implicit primary key. Plain tables keep
		// the no-clone fast path.
		if !alterNeedsFinalValidation(tbl, cmd) && !hasIgnoredIndex(tbl) &&
			cmd.Type != nodes.ATAlterIndexInvisible && !addNeedsRollback(cmd) {
			return c.execAlterCmd(db, tbl, cmd, dropsPK)
		}
		// Versioning commands need a post-mutation consistency check, so snapshot
		// to roll back an invalid result (e.g. DROP SYSTEM VERSIONING that leaves
		// orphaned period columns -> 4125).
		snapshot := cloneTable(tbl)
		if err := c.execAlterCmd(db, tbl, cmd, dropsPK); err != nil {
			*tbl = snapshot
			return err
		}
		if err := validateSystemVersioning(tbl); err != nil {
			*tbl = snapshot
			return err
		}
		if err := validateApplicationTimeFinalState(tbl); err != nil {
			*tbl = snapshot
			return err
		}
		if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
			*tbl = snapshot
			return err
		}
		return nil
	}

	// Multi-command ALTER: MySQL treats this as atomic.
	// Snapshot the table so we can rollback on any sub-command failure.
	snapshot := cloneTable(tbl)
	origKey := key

	rollback := func() { undoAlter(db, origKey, tbl, snapshot) }

	for _, cmd := range vectorKeysLast(stmt.Commands) {
		if skip[cmd] {
			continue
		}
		if err := c.execAlterCmd(db, tbl, cmd, dropsPK); err != nil {
			rollback()
			return err
		}
	}
	// Validate system-versioning consistency across the atomic ALTER as a whole
	// ("ADD PERIOD ..., ADD SYSTEM VERSIONING" is valid; either one alone is not).
	if err := validateSystemVersioning(tbl); err != nil {
		rollback()
		return err
	}
	if err := validateApplicationTimeFinalState(tbl); err != nil {
		rollback()
		return err
	}
	if err := validateIgnoredPrimaryKey(tbl, false); err != nil {
		rollback()
		return err
	}
	// Clear transient cleanup tracking after successful multi-command ALTER.
	tbl.droppedByCleanup = nil
	return nil
}

// addsVectorKey reports whether cmd adds a VECTOR key.
func addsVectorKey(cmd *nodes.AlterTableCmd) bool {
	return (cmd.Type == nodes.ATAddIndex || cmd.Type == nodes.ATAddConstraint) &&
		cmd.Constraint != nil && cmd.Constraint.Type == nodes.ConstrVector
}

// vectorKeysLast returns the commands with the ones adding a VECTOR key after the rest, each
// group in the order written. MariaDB holds a vector key to the table the statement leaves, not
// the one its clause sits in, on 11.8.8: ADD VECTOR i (v), ADD COLUMN v VECTOR(3) NOT NULL runs,
// as does ADD VECTOR i (v), MODIFY v VECTOR(3) NOT NULL and ADD VECTOR j (v), DROP INDEX i, while
// ADD VECTOR i (v), DROP COLUMN v and ADD VECTOR i (w), CHANGE w z ... are 1072.
func vectorKeysLast(cmds []*nodes.AlterTableCmd) []*nodes.AlterTableCmd {
	ordered := make([]*nodes.AlterTableCmd, 0, len(cmds))
	for _, cmd := range cmds {
		if !addsVectorKey(cmd) {
			ordered = append(ordered, cmd)
		}
	}
	for _, cmd := range cmds {
		if addsVectorKey(cmd) {
			ordered = append(ordered, cmd)
		}
	}
	return ordered
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

// execAlterCmd runs one ALTER command. dropsPK says whether another command of
// the same statement drops the primary key, which decides whether dropping one
// of its columns leaves a half-built key behind.
func (c *Catalog) execAlterCmd(db *Database, tbl *Table, cmd *nodes.AlterTableCmd, dropsPK bool) error {
	switch cmd.Type {
	case nodes.ATAddColumn:
		return c.alterAddColumn(tbl, cmd)
	case nodes.ATDropColumn:
		return c.alterDropColumn(tbl, cmd, dropsPK)
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
	case nodes.ATAddSystemVersioning:
		if tbl.SystemVersioned {
			return errAlreadySystemVersioned(tbl.Name)
		}
		tbl.SystemVersioned = true
		return nil
	case nodes.ATDropSystemVersioning:
		if !tbl.SystemVersioned {
			return errNotSystemVersioned(tbl.Name)
		}
		tbl.SystemVersioned = false
		return nil
	case nodes.ATAddPeriod:
		if !isSystemTimePeriodName(cmd.PeriodName) {
			return c.alterAddAppPeriod(tbl, cmd)
		}
		if tbl.PeriodStartCol != "" {
			return errAlreadySystemVersioned(tbl.Name)
		}
		tbl.PeriodStartCol = cmd.PeriodStartCol
		tbl.PeriodEndCol = cmd.PeriodEndCol
		return nil
	case nodes.ATDropPeriod:
		// IF EXISTS exists only in the application-time arm, DROP PERIOD
		// [IF EXISTS] FOR ident; DROP PERIOD FOR SYSTEM_TIME takes none, and
		// the parser rejects an unquoted SYSTEM_TIME after IF EXISTS. So
		// IF EXISTS FOR `SYSTEM_TIME` names an application-time period, which
		// MariaDB 11.8.8 skips when absent, even on a system-versioned table.
		if cmd.IfExists || !isSystemTimePeriodName(cmd.PeriodName) {
			return c.alterDropAppPeriod(tbl, cmd)
		}
		// Check period presence first so the result is independent of command
		// order: a DROP SYSTEM VERSIONING earlier in the same statement flips the
		// versioned flag but does not clear the period.
		if tbl.PeriodStartCol != "" {
			// Dropping an existing period orphans the ROW START/END columns;
			// MariaDB requires dropping them too (4125).
			return errMissingSystemVersioning(tbl.Name)
		}
		if !tbl.SystemVersioned {
			return errNotSystemVersioned(tbl.Name)
		}
		return errCantDropKey("PERIOD", "SYSTEM_TIME")
	default:
		// Unsupported alter command; silently ignore.
		return nil
	}
}

// isSystemTimePeriodName reports whether an ADD/DROP PERIOD command targets the
// system-versioning period rather than an application-time period.
func isSystemTimePeriodName(name string) bool {
	return name == "" || strings.EqualFold(name, "SYSTEM_TIME")
}

// alterAddAppPeriod adds an application-time period to an existing table,
// applying the same validations as CREATE (4154/1110/1060/1054/1063) and making
// the period columns NOT NULL. Grounded vs 11.8.8.
func (c *Catalog) alterAddAppPeriod(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.AppPeriodName != "" {
		return errMultipleAppPeriods() // 4154
	}
	if err := validateAppPeriodColumns(tbl, cmd.PeriodName, cmd.PeriodStartCol, cmd.PeriodEndCol); err != nil {
		return err
	}
	tbl.AppPeriodName = cmd.PeriodName
	tbl.AppPeriodStartCol = cmd.PeriodStartCol
	tbl.AppPeriodEndCol = cmd.PeriodEndCol
	// MariaDB makes application-time period columns NOT NULL.
	for _, colName := range []string{cmd.PeriodStartCol, cmd.PeriodEndCol} {
		if col := tbl.GetColumn(colName); col != nil {
			normalizeAppPeriodColumn(col)
		}
	}
	return nil
}

// alterDropAppPeriod removes the application-time period: 1091 if no such period
// exists, or nothing with IF EXISTS (MariaDB leaves note 1091). A WITHOUT
// OVERLAPS key still referencing it is caught by the final-state pass (4156),
// so a later DROP KEY in the same statement can remove the dependency first.
// The period columns keep their NOT NULL. Grounded vs 11.8.8.
func (c *Catalog) alterDropAppPeriod(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.AppPeriodName == "" || !strings.EqualFold(tbl.AppPeriodName, cmd.PeriodName) {
		if cmd.IfExists {
			return nil
		}
		return errCantDropPeriod(cmd.PeriodName) // 1091
	}
	tbl.AppPeriodName = ""
	tbl.AppPeriodStartCol = ""
	tbl.AppPeriodEndCol = ""
	return nil
}

// alterNeedsFinalValidation reports whether an ALTER command requires the
// post-mutation consistency checks (system-versioning + application-time), so the
// no-snapshot fast path is unsafe: the table is already system-versioned or has
// an application-time period (any column change can invalidate it), the command
// adds a key (a WITHOUT OVERLAPS key needs the final-state period check), or the
// command carries versioning / period metadata.
func alterNeedsFinalValidation(tbl *Table, cmd *nodes.AlterTableCmd) bool {
	if tbl.SystemVersioned || tbl.AppPeriodName != "" {
		return true
	}
	switch cmd.Type {
	case nodes.ATAddIndex, nodes.ATAddConstraint:
		return true
	}
	return cmdAffectsVersioning(cmd)
}

// cmdAffectsVersioning reports whether an ALTER command can change a table's
// system-versioning consistency (and so needs a post-ALTER validation). This is
// the four ADD/DROP SYSTEM VERSIONING / PERIOD commands, plus ADD COLUMN whose
// definition carries versioning metadata (ROW START/END or WITH/WITHOUT).
func cmdAffectsVersioning(cmd *nodes.AlterTableCmd) bool {
	switch cmd.Type {
	case nodes.ATAddSystemVersioning, nodes.ATDropSystemVersioning,
		nodes.ATAddPeriod, nodes.ATDropPeriod:
		return true
	case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
		if colDefHasVersioningMeta(cmd.Column) {
			return true
		}
		for _, c := range cmd.Columns {
			if colDefHasVersioningMeta(c) {
				return true
			}
		}
	}
	return false
}

// isColumnMutation reports whether an ALTER command adds, drops, modifies, or
// changes a column. RENAME COLUMN and index/table-option changes are excluded —
// MariaDB allows those on a system-versioned table.
func isColumnMutation(cmd *nodes.AlterTableCmd) bool {
	switch cmd.Type {
	case nodes.ATAddColumn, nodes.ATDropColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
		return true
	}
	return false
}

// addsColumnVersioning reports whether an ALTER command introduces a column
// carrying a column-level WITH/WITHOUT SYSTEM VERSIONING attribute (which
// requires the table to already be system-versioned, independent of command
// order). ADD/MODIFY/CHANGE COLUMN all carry a column definition.
func addsColumnVersioning(cmd *nodes.AlterTableCmd) bool {
	switch cmd.Type {
	case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
	default:
		return false
	}
	if cmd.Column != nil && cmd.Column.SystemVersioning != nodes.ColVersioningNone {
		return true
	}
	for _, c := range cmd.Columns {
		if c != nil && c.SystemVersioning != nodes.ColVersioningNone {
			return true
		}
	}
	return false
}

// colDefHasVersioningMeta reports whether a column definition carries
// system-versioning metadata (a ROW START/END bound or WITH/WITHOUT versioning).
func colDefHasVersioningMeta(colDef *nodes.ColumnDef) bool {
	if colDef == nil {
		return false
	}
	if colDef.SystemVersioning != nodes.ColVersioningNone {
		return true
	}
	return colDef.Generated != nil && colDef.Generated.RowBound != nodes.RowBoundNone
}

// alterAddColumn adds a new column to the table.
func (c *Catalog) alterAddColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	// The parenthesized form reads a table element list, so it carries the
	// constraints written beside the columns:
	// `ADD COLUMN (a1 INT, PRIMARY KEY (a1))` adds both. The columns come
	// first, as the CREATE TABLE path builds them, and MariaDB 11.8.8 makes a1
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
			}); err != nil {
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
	if err := validateColumnDefSemantics(colDef, indexedColumnsOf(tbl)); err != nil {
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
		case nodes.ColConstrCheck:
			// A check declared on the added column belongs to it and takes
			// the column's name, as one declared with the table does.
			// Measured on mariadb:11.8.8: `ALTER TABLE t ADD COLUMN b INT
			// CHECK (b > 0)` reports `b` int(11) DEFAULT NULL CHECK (`b` > 0)
			// and holds the check under the name b, even where the table
			// already has a check of that name — it takes no name of its own,
			// so it collides with none. See checkConstraintNameExists.
			conName := cc.Name
			if conName == "" {
				conName = colDef.Name
			}
			if err := validateColumnCheckReferences(colDef.Name, conName, cc.Expr); err != nil {
				return err
			}
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:        conName,
				Type:        ConCheck,
				Table:       tbl,
				Column:      colDef.Name,
				CheckExpr:   nodeToSQL(cc.Expr),
				NotEnforced: cc.NotEnforced,
			})
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

	// A JSON column carries json_valid only where it declares no check of its
	// own, so this is decided after the column's constraints are in: measured
	// on mariadb:11.8.8, `ALTER TABLE t ADD COLUMN j JSON CHECK (j IS NOT
	// NULL)` reports that check alone, where `ADD COLUMN j JSON` reports
	// json_valid(`j`).
	if col.IsJSON && !columnOwnsCheck(tbl, col.Name) {
		tbl.Constraints = append(tbl.Constraints, jsonValidConstraint(col))
	}

	return nil
}

// alterDropColumn removes a column from the table.
// alterDropsPrimaryKey reports whether any command of the statement drops the
// primary key. The parser builds `DROP PRIMARY KEY` as a DROP CONSTRAINT named
// PRIMARY.
func alterDropsPrimaryKey(cmds []*nodes.AlterTableCmd) bool {
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATDropConstraint && strings.EqualFold(cmd.Name, "PRIMARY") {
			return true
		}
	}
	return false
}

func (c *Catalog) alterDropColumn(tbl *Table, cmd *nodes.AlterTableCmd, droppingPrimaryKey bool) error {
	colKey := toLower(cmd.Name)
	if _, exists := tbl.colByName[colKey]; !exists {
		if cmd.IfExists {
			return nil
		}
		return errCantDropKey("COLUMN", cmd.Name)
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

	// A column a composite PRIMARY KEY is built on cannot be dropped while
	// that key survives the statement: MariaDB answers 1072, "Key column
	// '<col>' doesn't exist in table", rather than rebuilding the key over the
	// columns that are left. The same ALTER may drop the key itself, and then
	// the column goes with it — `DROP COLUMN a, DROP PRIMARY KEY` is accepted
	// whichever order the clauses come in. A single-column primary key goes
	// with its column. Measured on mariadb:11.8.8.
	if !droppingPrimaryKey {
		for _, idx := range tbl.Indexes {
			if !idx.Primary || len(idx.Columns) < 2 {
				continue
			}
			for _, ic := range idx.Columns {
				if toLower(ic.Name) == colKey {
					return &Error{
						Code:     1072,
						SQLState: "42000",
						Message:  fmt.Sprintf("Key column '%s' doesn't exist in table", cmd.Name),
					}
				}
			}
		}
	}

	// The check the column carries goes with it; leaving it behind would name
	// a column the table no longer has. Measured on mariadb:11.8.8: after
	// `ALTER TABLE t DROP COLUMN j` the table reports no check.
	dropColumnCheck(tbl, cmd.Name)

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
		// Remove the column from this index, tracking whether any ordinary
		// (non-WITHOUT OVERLAPS) key part remains.
		newCols := make([]*IndexColumn, 0, len(idx.Columns))
		ordinaryRemaining := false
		for _, ic := range idx.Columns {
			if toLower(ic.Name) == colKey {
				continue
			}
			newCols = append(newCols, ic)
			if !ic.WithoutOverlaps {
				ordinaryRemaining = true
			}
		}
		if !ordinaryRemaining {
			// No ordinary key parts left — the index is empty, or only a
			// WITHOUT OVERLAPS period part remains (not a real key column).
			// Remove the index entirely, like MariaDB.
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

	// Clean up constraints that reference removed indexes. Only a constraint
	// an index holds up goes with it; a CHECK merely carrying an index's name
	// has nothing to do with it. Measured on mariadb:11.8.8: after
	// `CREATE TABLE t (a INT, j JSON, KEY j(a)); ALTER TABLE t DROP COLUMN a`
	// the json_valid check column j carries is still there, and so is
	// `CONSTRAINT ck CHECK (b > 0)` beside a `KEY ck(a)`.
	if len(removedIndexNames) > 0 {
		newConstraints := make([]*Constraint, 0, len(tbl.Constraints))
		for _, con := range tbl.Constraints {
			if con.Type != ConCheck &&
				(removedIndexNames[toLower(con.IndexName)] || removedIndexNames[toLower(con.Name)]) {
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
	// Whether a key covers the column is a question about the column being
	// replaced, not about the name the replacement carries: a CHANGE renames
	// it and the key follows. Measured on mariadb:11.8.8, after
	// `CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))`:
	// `ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL` answers 1901,
	// where `CHANGE a b DOUBLE AS (2) VIRTUAL` is accepted and keeps
	// KEY `a` (`b`).
	indexed := indexedColumnsOf(tbl)
	if indexed[toLower(oldName)] {
		indexed[toLower(colDef.Name)] = true
	}
	if err := validateColumnDefSemantics(colDef, indexed); err != nil {
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

	// MODIFY and CHANGE replace the column's definition, so the check the old
	// definition carried goes with it and the new one's takes its place: a
	// JSON column's is json_valid, and any other definition's is whatever it
	// declares, which is nothing here since a column-level CHECK is not part
	// of what MODIFY accepts. Measured on mariadb:11.8.8:
	// `ALTER TABLE t MODIFY COLUMN a INT` against
	// `CREATE TABLE t (a INT CHECK (a > 0))` leaves no check, and
	// `ALTER TABLE t MODIFY COLUMN j JSON` against
	// `CREATE TABLE t (j JSON CHECK (j IS NOT NULL))` replaces that check
	// with json_valid(`j`).
	dropColumnCheck(tbl, oldName)
	if col.IsJSON {
		tbl.Constraints = append(tbl.Constraints, jsonValidConstraint(col))
	}

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
// checkAddIfNotExistsTwice answers for what IF NOT EXISTS does not cover. The clause looks at
// the table as it was when the statement began, so a primary key or CHECK the statement also
// adds is not there to find and the two fail, in either order, where a table that has the key
// makes the clause a note. A second index of a name goes as a note too, but only of its own kind:
// mariadb:11.8.8 answers 1061 for ADD INDEX i (a), ADD UNIQUE IF NOT EXISTS i (a) and for the two
// the other way round, and runs ADD INDEX i (a), ADD KEY IF NOT EXISTS i (a) (read on 2026-11-05).
func checkAddIfNotExistsTwice(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	var primaries int
	checks := make(map[string]bool)
	kinds := make(map[string]nodes.ConstraintType)
	see := func(con *nodes.Constraint) error {
		if con.IfNotExists && addedConstraintExists(tbl, con, extractColumnNames(con)) {
			return nil
		}
		switch con.Type {
		case nodes.ConstrPrimaryKey:
			if primaries++; primaries > 1 {
				return errMultiplePriKey()
			}
		case nodes.ConstrCheck:
			if con.Name == "" {
				return nil
			}
			name := toLower(con.Name)
			if checks[name] {
				return &Error{Code: ErrFKDupName, SQLState: "HY000", Message: fmt.Sprintf("Duplicate CHECK constraint name '%s'", con.Name)}
			}
			checks[name] = true
		case nodes.ConstrForeignKey:
		default:
			name := con.Name
			if name == "" {
				return nil
			}
			key := toLower(name)
			if kind, added := kinds[key]; added && kind != con.Type {
				return errDupKeyName(name)
			}
			kinds[key] = con.Type
		}
		return nil
	}
	for _, cmd := range cmds {
		if cmd.Constraint != nil {
			if err := see(cmd.Constraint); err != nil {
				return err
			}
		}
		for _, con := range cmd.Constraints {
			if err := see(con); err != nil {
				return err
			}
		}
		// A column's own PRIMARY KEY counts too, in ADD, MODIFY and CHANGE: ALTER TABLE t ADD b
		// INT PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (a) is 1068 in either order.
		columns := cmd.Columns
		if cmd.Column != nil {
			columns = append([]*nodes.ColumnDef{cmd.Column}, columns...)
		}
		for _, col := range columns {
			for _, cc := range col.Constraints {
				if cc.Type == nodes.ColConstrPrimaryKey {
					if primaries++; primaries > 1 {
						return errMultiplePriKey()
					}
				}
			}
		}
	}
	return nil
}

// addedConstraintExists reports whether ADD ... IF NOT EXISTS finds what it would add,
// which makes the clause a note and leaves the table as it was, ahead of every other check:
// mariadb:11.8.8 runs `ADD INDEX IF NOT EXISTS k (nosuch)` where k exists and answers 1072
// where it does not (read on 2026-11-05). What it looks for depends on the kind:
//
//   - a primary key looks for any primary key;
//   - a CHECK looks for a check constraint of the name;
//   - a foreign key looks for a foreign key on this table named after the name that follows
//     FOREIGN KEY or else the symbol. One with neither has nothing to look for and is
//     always added, and the symbol is not looked at when both are written, so `CONSTRAINT
//     fky FOREIGN KEY IF NOT EXISTS fkw` on a table with fky is 1005;
//   - any other key looks for an index of its name, which left out is its first column's:
//     `ADD UNIQUE IF NOT EXISTS (a)` is a note beside `KEY a (a)`.
func addedConstraintExists(tbl *Table, con *nodes.Constraint, cols []string) bool {
	switch con.Type {
	case nodes.ConstrPrimaryKey:
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				return true
			}
		}
		return false
	case nodes.ConstrCheck:
		// A check that belongs to a column does not occupy the name, so
		// IF NOT EXISTS over a JSON column's own name still adds the
		// statement's check. See checkConstraintNameExists.
		for _, existing := range tbl.Constraints {
			if existing.Type == ConCheck && !isColumnCheck(existing) &&
				toLower(existing.Name) == toLower(con.Name) {
				return true
			}
		}
		return false
	case nodes.ConstrForeignKey:
		name := con.IndexName
		if name == "" {
			name = con.Name
		}
		if name == "" {
			return false
		}
		// A foreign key of the name is what it looks for, on this table. An index of the
		// name is not: with KEY k on t, ADD FOREIGN KEY IF NOT EXISTS k (b) REFERENCES r (id)
		// adds it on 11.8.8.
		for _, existing := range tbl.Constraints {
			if existing.Type == ConForeignKey && toLower(existing.Name) == toLower(name) {
				return true
			}
		}
		return false
	}
	name := con.Name
	if name == "" && len(cols) > 0 {
		name = cols[0]
	}
	return name != "" && indexNameExists(tbl, name)
}

func (c *Catalog) alterAddConstraint(tbl *Table, cmd *nodes.AlterTableCmd) error {
	con := cmd.Constraint
	if con == nil {
		return nil
	}

	cols := extractColumnNames(con)
	if con.IfNotExists && addedConstraintExists(tbl, con, cols) {
		return nil
	}
	if err := validateKeyPartPrefixes(con.IndexColumns, con.Type == nodes.ConstrSpatialIndex); err != nil {
		return err
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
		// The ALTER's snapshot puts the columns back: MariaDB 11.8.8 leaves
		// them nullable.
		if !pkIdx.Visible {
			return errPKIndexCantBeIgnored()
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
		conName := fkName(con)
		if conName == "" {
			conName = fmt.Sprintf("%s_ibfk_%d", tbl.Name, nextFKGeneratedNumber(tbl, tbl.Name))
		}
		if foreignKeyConstraintNameExists(tbl.Database, nil, conName) {
			return errFKDupName(conName)
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
		ensureFKBackingIndex(tbl, fkName(con), cols, buildIndexColumns(con))

	case nodes.ConstrCheck:
		conName := con.Name
		if conName == "" {
			conName = fmt.Sprintf("CONSTRAINT_%d", nextCheckNumber(tbl))
		}
		if checkConstraintNameExists(tbl.Database, tbl, conName) {
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
		// added by ALTER lost its comment and block size.
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
		tbl.Indexes = append(tbl.Indexes, ftIdx)

	case nodes.ConstrVector:
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
		if err := validateVectorIndex(tbl, idxCols); err != nil {
			return err
		}
		vecIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			Vector:    true,
			IndexType: "VECTOR",
			Visible:   true,
		}
		applyIndexOptions(vecIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, vecIdx)

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
		return errCantDropKey("INDEX", name)
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

// dropColumnCheck removes the check a column carries, which goes with the
// column: dropped with it, and replaced along with its definition. Measured on
// mariadb:11.8.8: after `ALTER TABLE t DROP COLUMN j` a JSON column's
// json_valid check is gone, and `ALTER TABLE t MODIFY COLUMN a INT` against
// `CREATE TABLE t (a INT CHECK (a > 0))` leaves no check behind — the new
// definition declares none, so the column has none.
//
// It finds the check by the column it belongs to rather than by its name,
// which can have moved on: measured, after `ALTER TABLE t RENAME COLUMN j TO
// z` the check reads json_valid(`z`) while its name stays j.
func dropColumnCheck(tbl *Table, colName string) {
	key := toLower(colName)
	kept := tbl.Constraints[:0]
	for _, con := range tbl.Constraints {
		if isColumnCheck(con) && toLower(con.Column) == key {
			continue
		}
		kept = append(kept, con)
	}
	tbl.Constraints = kept
}

// isImplicitJSONCheck reports whether the constraint is the json_valid check a
// JSON column carries, which the catalog names after that column. It is not a
// constraint a statement can name: it comes and goes with its column.
//
// The answer comes from the constraint's own record of how it was made, not
// from its SQL, since a statement may declare a check with the same name and
// expression and that one is ordinary. See Constraint.Implicit.
func isImplicitJSONCheck(con *Constraint) bool {
	return con != nil && con.Type == ConCheck && con.Implicit
}

// renameImplicitJSONCheck follows a renamed JSON column: the check belongs to
// the column and references it, so both move with it. Measured on
// mariadb:11.8.8: after `ALTER TABLE t CHANGE COLUMN j k JSON` the check reads
// CHECK (json_valid(`k`)).
//
// The check is found by the column it belongs to, since another constraint may
// have taken its name in the meantime.
func renameImplicitJSONCheck(tbl *Table, oldName, newName string) {
	key := toLower(oldName)
	for _, con := range tbl.Constraints {
		if isImplicitJSONCheck(con) && toLower(con.Column) == key {
			con.Name = newName
			con.Column = newName
			con.CheckExpr = fmt.Sprintf("json_valid(`%s`)", newName)
		}
	}
}

// indexedColumnsOf returns the lowercased names of the columns the table's
// indexes cover.
func indexedColumnsOf(tbl *Table) map[string]bool {
	out := make(map[string]bool)
	if tbl == nil {
		return out
	}
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			out[toLower(ic.Name)] = true
		}
	}
	return out
}

// alterDropConstraint removes a constraint by name.
func (c *Catalog) alterDropConstraint(tbl *Table, cmd *nodes.AlterTableCmd) error {
	name := cmd.Name
	key := toLower(name)

	found := false
	isForeignKey := false
	for i, con := range tbl.Constraints {
		if toLower(con.Name) != key {
			continue
		}
		// A check that belongs to a column carries the column's name but is
		// not a constraint the statement can name: measured on
		// mariadb:11.8.8, the server answers 1091 both for
		// `ALTER TABLE t DROP CONSTRAINT j` on a JSON column `j` and for
		// `ALTER TABLE t DROP CONSTRAINT a` against
		// `CREATE TABLE t (a INT CHECK (a > 0))`. Such a check goes only with
		// its column.
		if isColumnCheck(con) {
			continue
		}
		isForeignKey = (con.Type == ConForeignKey)
		tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
		found = true
		break
	}

	if !found {
		if cmd.IfExists {
			return nil
		}
		// The server names the spelling it was given, so a missing name under
		// DROP FOREIGN KEY reports FOREIGN KEY where DROP CONSTRAINT reports
		// CONSTRAINT. The AST folds both to ATDropConstraint and keeps no
		// record of which was written, so this reports CONSTRAINT for both.
		return errCantDropKey("CONSTRAINT", name)
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
	// Keep PERIOD FOR SYSTEM_TIME referencing the renamed column (MariaDB
	// rewrites it).
	if strings.EqualFold(tbl.PeriodStartCol, cmd.Name) {
		tbl.PeriodStartCol = cmd.NewName
	}
	if strings.EqualFold(tbl.PeriodEndCol, cmd.Name) {
		tbl.PeriodEndCol = cmd.NewName
	}
	// Likewise keep the application-time period referencing the renamed column.
	if strings.EqualFold(tbl.AppPeriodStartCol, cmd.Name) {
		tbl.AppPeriodStartCol = cmd.NewName
	}
	if strings.EqualFold(tbl.AppPeriodEndCol, cmd.Name) {
		tbl.AppPeriodEndCol = cmd.NewName
	}
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
	// have left both e.u and e.v pointing at the table; MariaDB 11.8.8 runs
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
		tbl.Engine = opt.Value
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

// alterIndexVisibility toggles the Visible flag on an index.
// alterIndexVisibility applies ALTER {INDEX | KEY} [IF EXISTS] name [NOT]
// IGNORED. checkAlterIndexIgnorability has held the name to the table the
// statement found and dropped the IF EXISTS commands that name no index of
// it, so an index missing here went with another command of the statement,
// which MariaDB 11.8.8 answers with 1176, IF EXISTS or not.
func (c *Catalog) alterIndexVisibility(tbl *Table, cmd *nodes.AlterTableCmd, visible bool) error {
	key := toLower(cmd.Name)
	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == key {
			idx.Visible = visible
			return nil
		}
	}
	return errKeyDoesNotExist(cmd.Name, tbl.Name)
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
		// A column added or retyped JSON takes the same shape it would in a
		// CREATE TABLE: longtext with the binary collation, and the
		// json_valid check its caller attaches. Measured on mariadb:11.8.8:
		// `ALTER TABLE t ADD COLUMN j JSON` reports
		// `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
		// CHECK (json_valid(`j`)).
		if col.DataType == "json" {
			applyJSONColumn(col)
		}
		if isNationalStringType(colDef.TypeName.Name) {
			col.Charset = "utf8mb3"
		}
		if colDef.TypeName.Charset != "" {
			col.Charset = normalizeCharsetName(colDef.TypeName.Charset)
		}
		if colDef.TypeName.Collate != "" {
			col.Collation = colDef.TypeName.Collate
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
	if colDef.DefaultValue != nil {
		setColumnDefaultFromExpr(col, colDef.DefaultValue)
	}
	if colDef.OnUpdate != nil {
		setColumnOnUpdateFromExpr(col, colDef.OnUpdate)
	}
	if colDef.Generated != nil {
		col.Generated = &GeneratedColumnInfo{
			Expr:             nodeToSQLGenerated(colDef.Generated.Expr, tbl.Charset),
			Stored:           colDef.Generated.Stored,
			RowBound:         rowBoundSQL(colDef.Generated.RowBound),
			Nondeterministic: nondeterministicFuncIn(colDef.Generated.Expr),
		}
	}
	col.SystemVersioning = colVersioningSQL(colDef.SystemVersioning)

	// Process column-level constraints (non-index-producing ones).
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull:
			col.Nullable = false
		case nodes.ColConstrNull:
			col.Nullable = true
		case nodes.ColConstrDefault:
			if cc.Expr != nil {
				setColumnDefaultFromExpr(col, cc.Expr)
			}
		case nodes.ColConstrAutoIncrement:
			col.AutoIncrement = true
			col.Nullable = false
		case nodes.ColConstrInvisible:
			col.Invisible = true
		case nodes.ColConstrCollate:
			if cc.Expr != nil {
				if s, ok := cc.Expr.(*nodes.StringLit); ok {
					col.Collation = s.Value
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
	renameImplicitJSONCheck(tbl, oldName, newName)
}

// alterCheckEnforced refuses `ALTER TABLE ... ALTER CHECK <name> [NOT]
// ENFORCED`, which MariaDB has no grammar for: a CHECK constraint carries no
// enforcement attribute there, and information_schema.TABLE_CONSTRAINTS has no
// ENFORCED column. MySQL 8.0 has both. The parser still accepts the clause —
// one of the MySQL-isms mariadb/parser's divergence inventory pins — so the
// refusal lands here. Measured on mariadb:11.8.8: the server answers 1064.
func (c *Catalog) alterCheckEnforced(tbl *Table, cmd *nodes.AlterTableCmd) error {
	key := toLower(cmd.Name)
	for _, con := range tbl.Constraints {
		if toLower(con.Name) == key && con.Type == ConCheck {
			return &Error{
				Code:     1064,
				SQLState: "42000",
				Message: fmt.Sprintf("You have an error in your SQL syntax; check the manual that "+
					"corresponds to your MariaDB server version for the right syntax to use near "+
					"'CHECK %s %s' at line 1", cmd.Name, cmd.NewName),
			}
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
		return fmt.Errorf("ALTER TABLE ADD PARTITION: table '%s' is not partitioned", tbl.Name)
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
				pdi.Engine = opt.Value
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
		return fmt.Errorf("ALTER TABLE DROP PARTITION: table '%s' is not partitioned", tbl.Name)
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
		return fmt.Errorf("ALTER TABLE TRUNCATE PARTITION: table '%s' is not partitioned", tbl.Name)
	}
	// Truncate is a data operation; for DDL catalog purposes, it's a no-op.
	return nil
}

// alterCoalescePartition reduces the number of partitions for HASH/KEY partitioned tables.
func (c *Catalog) alterCoalescePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return fmt.Errorf("ALTER TABLE COALESCE PARTITION: table '%s' is not partitioned", tbl.Name)
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
		return fmt.Errorf("ALTER TABLE REORGANIZE PARTITION: table '%s' is not partitioned", tbl.Name)
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
		return fmt.Errorf("ALTER TABLE EXCHANGE PARTITION: table '%s' is not partitioned", tbl.Name)
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
// MariaDB 11.8.8 leaves the table untouched for all of them:
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
