package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

func (c *Catalog) alterTable(stmt *nodes.AlterTableStmt) error {
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
	key := toLower(tableName)
	// An ENGINE option naming an engine the pinned server does not ship is
	// 1286, answered before the table is looked up (measured on
	// mariadb:11.8.8: `ALTER TABLE a ENGINE zzz` on no table answers 1286,
	// not 1146).
	if err := validateAlteredStorageEngines(stmt.Commands); err != nil {
		return err
	}
	tbl := db.Tables[key]
	if tbl == nil {
		// IF EXISTS makes a missing table a note (1146 as a warning).
		if stmt.IfExists {
			return nil
		}
		// A partition maintenance command is not an error on a table that is not
		// there: OPTIMIZE, ANALYZE, CHECK and REPAIR PARTITION run and report
		// nothing (measured on mariadb:11.8.8), where REBUILD PARTITION and
		// REMOVE [FORCE] PARTITIONING answer 1146.
		if onlyPartitionMaintenance(stmt.Commands) {
			return nil
		}
		return errNoSuchTable(db.Name, tableName)
	}

	return c.alterTableKeepingOptions(db, tbl, key, stmt)
}

// alterTableKeepingOptions runs the ALTER and then drops PAGE_CHECKSUM unless the
// statement set it or only renamed the table, or the table is Aria's: MariaDB
// 11.8.8 keeps it through a RENAME and through nothing else that rewrites an
// InnoDB or MyISAM table, and Aria's handler carries it over.
func (c *Catalog) alterTableKeepingOptions(db *Database, tbl *Table, key string, stmt *nodes.AlterTableStmt) error {
	// A ROW_FORMAT FIXED or PAGE the statement writes is the engine's to take, and
	// InnoDB refuses it with 1478 — after the statement's other commands, which
	// answer first: `ROW_FORMAT FIXED, ADD KEY(a)` is 1072. MariaDB leaves the
	// table as it found it, so hold a copy to put back when the row format is what
	// the statement fails on.
	// A discarded tablespace leaves nothing to copy the table into, so the ALTERs
	// that need a copy of it are refused (1814) where the in-place ones run:
	// measured on mariadb:11.8.8 over a table whose tablespace was DISCARDed,
	// `ROW_FORMAT=COMPACT`, `ENGINE=MyISAM`, `MODIFY b BIGINT` and
	// `ALGORITHM=COPY, ADD d INT` answer 1814, and `ADD c INT`, `ADD INDEX i (b)`
	// and `RENAME TO u` run.
	if tbl.TablespaceDiscarded {
		if err := refuseCopyAlterOnDiscarded(tbl, stmt.Commands); err != nil {
			return err
		}
	}
	refusesRowFormat := writesFixedRowFormat(tbl, stmt.Commands)
	// MariaDB holds every table an ALTER leaves to the AUTO_INCREMENT rule
	// (1075), so the statement that can unkey the column needs a snapshot to put
	// the table back with when it fails. See validateAutoIncrementKeyed.
	checksAutoIncrement := needsAutoIncrementCheck(tbl, stmt.Commands)
	var snapshot Table
	if refusesRowFormat || checksAutoIncrement {
		snapshot = cloneTable(tbl)
	}
	if err := c.alterTableRun(db, tbl, key, stmt); err != nil {
		return err
	}
	if checksAutoIncrement {
		if err := validateAutoIncrementKeyed(tbl); err != nil {
			undoAlter(db, key, tbl, snapshot)
			return err
		}
	}
	if refusesRowFormat {
		if err := validateAlteredRowFormat(tbl); err != nil {
			undoAlter(db, key, tbl, snapshot)
			return err
		}
	}
	rebuilt, setsPageChecksum := false, false
	for _, cmd := range stmt.Commands {
		if cmd.Type != nodes.ATRenameTable {
			rebuilt = true
		}
		if cmd.Type == nodes.ATTableOption && cmd.Option != nil && eqFoldStr(cmd.Option.Name, "PAGE_CHECKSUM") {
			setsPageChecksum = true
		}
	}
	// Aria recovers the checksum it has when the ALTER names none.
	if rebuilt && !setsPageChecksum && !eqFoldStr(tbl.Engine, "Aria") {
		tbl.PageChecksum = ""
	}
	return nil
}

// refuseCopyAlterOnDiscarded refuses an ALTER that needs a copy of a table whose
// tablespace the statement DISCARDed: 1814. The catalog knows the commands that
// rewrite the table's storage — a ROW_FORMAT or an ENGINE option, and a column
// whose type changes — and not the rest (MariaDB also refuses FORCE and an added
// CHECK there).
func refuseCopyAlterOnDiscarded(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	for _, cmd := range cmds {
		switch cmd.Type {
		case nodes.ATForce, nodes.ATOrderBy:
			return errTablespaceDiscarded(tbl.Name)
		case nodes.ATAlgorithm:
			// ALGORITHM=COPY asks for the copy by name, where INPLACE and DEFAULT
			// leave the choice to the server: measured on mariadb:11.8.8 over a
			// discarded table, COPY is 1814 and the other two run.
			if eqFoldStr(cmd.Name, "COPY") {
				return errTablespaceDiscarded(tbl.Name)
			}
		case nodes.ATTableOption:
			if cmd.Option == nil {
				continue
			}
			switch toLower(cmd.Option.Name) {
			case "engine", "row_format":
				return errTablespaceDiscarded(tbl.Name)
			}
		case nodes.ATAddIndex, nodes.ATAddConstraint:
			// A UNIQUE, PRIMARY, FOREIGN, FULLTEXT or CHECK key rebuilds the
			// table where an ordinary or SPATIAL one does not: measured on
			// mariadb:11.8.8 over a discarded table, `ADD INDEX i (b)` and
			// `ADD SPATIAL KEY s (g)` run and the other five answer 1814.
			if cmd.Constraint == nil {
				continue
			}
			switch cmd.Constraint.Type {
			case nodes.ConstrPrimaryKey, nodes.ConstrUnique, nodes.ConstrForeignKey,
				nodes.ConstrCheck, nodes.ConstrFulltextIndex:
				return errTablespaceDiscarded(tbl.Name)
			}
		case nodes.ATModifyColumn, nodes.ATChangeColumn:
			for _, colDef := range append([]*nodes.ColumnDef{cmd.Column}, cmd.Columns...) {
				if colDef == nil || colDef.TypeName == nil {
					continue
				}
				col := tbl.GetColumn(cmd.Name)
				if col == nil {
					col = tbl.GetColumn(colDef.Name)
				}
				if col == nil {
					continue
				}
				if !eqFoldStr(col.ColumnType, formatColumnType(colDef.TypeName)) {
					return errTablespaceDiscarded(tbl.Name)
				}
			}
		}
	}
	return nil
}

// hasSecondaryIndex reports whether the table holds an index that is not its
// primary key, which is what IMPORT TABLESPACE asks to be dropped.
func hasSecondaryIndex(tbl *Table) bool {
	for _, idx := range tbl.Indexes {
		if !idx.Primary {
			return true
		}
	}
	return false
}

// needsAutoIncrementCheck reports whether the ALTER can leave the table with an
// AUTO_INCREMENT column unkeyed or with two: the table has one, a column
// definition may add one, or a command adds or drops a key or a column. Such a
// statement is held to validateAutoIncrementKeyed once its commands have run.
func needsAutoIncrementCheck(tbl *Table, cmds []*nodes.AlterTableCmd) bool {
	for _, col := range tbl.Columns {
		if col.AutoIncrement {
			return true
		}
	}
	for _, cmd := range cmds {
		switch cmd.Type {
		case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn,
			nodes.ATDropColumn, nodes.ATAddIndex, nodes.ATAddConstraint,
			nodes.ATDropIndex, nodes.ATDropConstraint:
			return true
		}
	}
	return false
}

// validateAutoIncrementKeyed holds the table an ALTER leaves to MariaDB's
// AUTO_INCREMENT rule: at most one such column, and it the first part of a key
// (1075). Measured on mariadb:11.8.8: `ALTER TABLE a ADD c INT AUTO_INCREMENT`
// answers it over a table whose keys do not lead with c, and `DROP KEY kk` does
// over a table whose only key over the auto column is kk, where the statement
// that drops kk and adds a key over c runs.
func validateAutoIncrementKeyed(tbl *Table) error {
	auto := ""
	for _, col := range tbl.Columns {
		if !col.AutoIncrement {
			continue
		}
		if auto != "" {
			return errIncorrectAutoIncrementDefinition()
		}
		auto = col.Name
	}
	if auto == "" || autoIncrementKeyed(tbl, auto) {
		return nil
	}
	return errIncorrectAutoIncrementDefinition()
}

// autoIncrementKeyed reports whether the table's keys take an AUTO_INCREMENT
// column where its engine needs it: the first part of a key on InnoDB, and
// anywhere in a key on MyISAM and Aria. Measured on mariadb:11.8.8:
// `CREATE TABLE t (g INT, id INT AUTO_INCREMENT, KEY k (g, id))` runs on MyISAM
// and Aria and answers 1075 on InnoDB, and `DROP KEY k` over the MyISAM table
// answers 1075 once no key holds the column at all.
func autoIncrementKeyed(tbl *Table, name string) bool {
	key := toLower(name)
	first, anywhere := false, false
	for _, idx := range tbl.Indexes {
		for i, ic := range idx.Columns {
			if ic == nil || ic.Name == "" || toLower(ic.Name) != key {
				continue
			}
			anywhere = true
			if i == 0 {
				first = true
			}
		}
	}
	// InnoDB and MEMORY want the column first, MyISAM and Aria take it anywhere,
	// measured on mariadb:11.8.8: `CREATE TABLE t (g INT, id INT AUTO_INCREMENT,
	// KEY k (g, id))` answers 1075 on InnoDB and MEMORY and runs on MyISAM and
	// Aria. An engine the catalog does not model keeps the stricter rule.
	switch engine := toLower(tbl.Engine); engine {
	case "myisam", "aria":
		return anywhere
	}
	return first
}

// alterTableRun runs the commands of an ALTER TABLE, checking a partitioned
// table's unique keys after them.
func (c *Catalog) alterTableRun(db *Database, tbl *Table, key string, stmt *nodes.AlterTableStmt) error {
	repartition := alterRepartition(stmt.Commands)
	// SYSTEM_TIME partitioning is a versioned table's: measured on
	// mariadb:11.8.8, `ALTER TABLE a PARTITION BY SYSTEM_TIME` over a table
	// that is not system-versioned answers 4124.
	// The table the statement leaves decides it: `ALTER TABLE t ADD SYSTEM
	// VERSIONING PARTITION BY SYSTEM_TIME` runs, measured on mariadb:11.8.8, the
	// other orders being 1064.
	if repartition != nil && repartition.Type == nodes.PartitionSystemTime && !tbl.SystemVersioned &&
		!alterAddsSystemVersioning(stmt.Commands) {
		return errNotSystemVersioned(tbl.Name)
	}
	systemTimeRepartition := repartition != nil && repartition.Type == nodes.PartitionSystemTime
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
	// The table the statement leaves is the one SYSTEM_TIME partitioning needs
	// versioned: `ALTER TABLE t DROP SYSTEM VERSIONING PARTITION BY SYSTEM_TIME`
	// answers 4124, where the ADD form runs (measured on mariadb:11.8.8).
	if systemTimeRepartition && !tbl.SystemVersioned {
		undoAlter(db, key, tbl, snapshot)
		return errNotSystemVersioned(tbl.Name)
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
		// So is a ROW START or ROW END column, unless the statement makes the
		// table system-versioned too. Measured on mariadb:11.8.8, which answers
		// 4124 for `ALTER TABLE t ADD s TIMESTAMP(6) AS ROW START` and for MODIFY,
		// before it looks at the column's name or keys, and runs the ADD of both
		// columns, the period and SYSTEM VERSIONING in one statement.
		if !alterAddsSystemVersioning(stmt.Commands) {
			for _, cmd := range stmt.Commands {
				if addsRowBoundColumn(cmd) {
					return errNotSystemVersioned(tbl.Name)
				}
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

	// The keys IF NOT EXISTS finds are notes whose definitions are never read,
	// and which is found is the table the statement began with.
	skippedKeys := map[*nodes.Constraint]bool{}
	for _, cmd := range stmt.Commands {
		cons := cmd.Constraints
		if cmd.Constraint != nil {
			cons = append([]*nodes.Constraint{cmd.Constraint}, cons...)
		}
		for _, con := range cons {
			if con != nil && con.IfNotExists && addedConstraintExists(tbl, con, extractColumnNames(con)) {
				skippedKeys[con] = true
			}
		}
	}

	// A key part a statement adds WITHOUT OVERLAPS names a period, and MariaDB
	// judges it before the key's own columns: measured on mariadb:11.8.8,
	// `ALTER TABLE a ADD UNIQUE(a, a WITHOUT OVERLAPS)` over a table with neither
	// a column a nor a period a answers 4156 (the missing period), not 1072.
	// The statement may add the period after the key — `ADD UNIQUE uk (id,
	// app_time WITHOUT OVERLAPS), ADD PERIOD FOR app_time(s,e)` runs — so the
	// periods it adds count too.
	if err := validateAddedWithoutOverlapsParts(tbl, stmt.Commands, skippedKeys); err != nil {
		return err
	}

	skip, err := checkAlterIndexIgnorability(tbl, stmt.Commands)
	if err != nil {
		return err
	}

	// The columns the statement adds are read before any of them is built, so a
	// generated one may not read a generated column the same statement adds after
	// it. That is a check over the whole statement, and it runs before any command
	// changes the table.
	if err := c.validateAddedColumnOrder(tbl, stmt.Commands); err != nil {
		return err
	}

	dropsPK := alterDropsPrimaryKey(stmt.Commands)
	droppedColumns := alterDroppedColumns(stmt.Commands)
	declaredChecks := alterDeclaredCheckNames(stmt.Commands)

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
			return c.execAlterCmd(db, tbl, cmd, dropsPK, droppedColumns, declaredChecks)
		}
		// Versioning commands need a post-mutation consistency check, so snapshot
		// to roll back an invalid result (e.g. DROP SYSTEM VERSIONING that leaves
		// orphaned period columns -> 4125).
		snapshot := cloneTable(tbl)
		if err := c.execAlterCmd(db, tbl, cmd, dropsPK, droppedColumns, declaredChecks); err != nil {
			*tbl = snapshot
			return err
		}
		dropIgnoredSystemTimePeriod(tbl)
		if err := validateSystemVersioning(tbl); err != nil {
			*tbl = snapshot
			return err
		}
		if err := validateApplicationTimeFinalState(tbl); err != nil {
			*tbl = snapshot
			return err
		}
		if err := validateAlteredColumnExprColumns(tbl, stmt.Commands); err != nil {
			*tbl = snapshot
			return err
		}
		if err := c.validateAddedKeyEngineOptions(tbl, stmt.Commands, skippedKeys); err != nil {
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
		if err := c.execAlterCmd(db, tbl, cmd, dropsPK, droppedColumns, declaredChecks); err != nil {
			rollback()
			return err
		}
	}
	// Validate system-versioning consistency across the atomic ALTER as a whole
	// ("ADD PERIOD ..., ADD SYSTEM VERSIONING" is valid; either one alone is not).
	dropIgnoredSystemTimePeriod(tbl)
	if err := validateSystemVersioning(tbl); err != nil {
		rollback()
		return err
	}
	if err := validateApplicationTimeFinalState(tbl); err != nil {
		rollback()
		return err
	}
	if err := validateAlteredColumnExprColumns(tbl, stmt.Commands); err != nil {
		rollback()
		return err
	}
	if err := c.validateAddedKeyEngineOptions(tbl, stmt.Commands, skippedKeys); err != nil {
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

// dropIgnoredSystemTimePeriod drops a PERIOD FOR SYSTEM_TIME that sits over
// ordinary columns of a table that is not system-versioned. MariaDB's ALTER
// accepts one and stores nothing: measured on mariadb:11.8.8,
// `ALTER TABLE a ADD PERIOD FOR SYSTEM_TIME(a, a)` over a plain table a runs and
// SHOW CREATE TABLE a shows the table as it was, where the same PERIOD written
// in a CREATE TABLE is 4125. A period over ROW START/ROW END columns stays for
// validateSystemVersioning to refuse — those columns are already refused up
// front (4124) on a table that is not versioned.
func dropIgnoredSystemTimePeriod(tbl *Table) {
	if tbl.SystemVersioned || tbl.PeriodStartCol == "" {
		return
	}
	start, end := tbl.GetColumn(tbl.PeriodStartCol), tbl.GetColumn(tbl.PeriodEndCol)
	if start != nil && end != nil &&
		start.Generated != nil && start.Generated.RowBound == "ROW START" &&
		end.Generated != nil && end.Generated.RowBound == "ROW END" {
		return
	}
	tbl.PeriodStartCol, tbl.PeriodEndCol = "", ""
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
// the same statement drops the primary key and droppedColumns holds every
// column the statement drops, which together decide whether dropping one of a
// composite key's columns would leave a half-built key behind; declaredChecks
// holds the CHECK names the statement gives explicitly, which the unnamed ones
// are named around.
func (c *Catalog) execAlterCmd(db *Database, tbl *Table, cmd *nodes.AlterTableCmd, dropsPK bool, droppedColumns map[string]bool, declaredChecks map[string]bool) error {
	switch cmd.Type {
	case nodes.ATAddColumn:
		return c.alterAddColumn(tbl, cmd)
	case nodes.ATDropColumn:
		return c.alterDropColumn(tbl, cmd, dropsPK, droppedColumns)
	case nodes.ATModifyColumn:
		return c.alterModifyColumn(tbl, cmd)
	case nodes.ATChangeColumn:
		return c.alterChangeColumn(tbl, cmd)
	case nodes.ATAddIndex, nodes.ATAddConstraint:
		return c.alterAddConstraint(tbl, cmd, declaredChecks)
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
	case nodes.ATConvertPartitionToTable:
		return c.alterConvertPartitionToTable(db, tbl, cmd)
	case nodes.ATConvertTableToPartition:
		return c.alterConvertTableToPartition(db, tbl, cmd)
	case nodes.ATImportTablespace:
		// A tablespace is InnoDB's: measured on mariadb:11.8.8, MyISAM, Aria and
		// MEMORY answer 1031 for both commands. Over an InnoDB table the server
		// gave itself a tablespace, so IMPORT answers 1813 until the statement
		// DISCARDs it, and 1030 from then on (there is no file to import).
		if err := refuseTablespaceEngine(tbl); err != nil {
			return err
		}
		if tbl.TablespaceDiscarded {
			// The import reads the .cfg file DISCARD removed: a table with a
			// secondary index asks for it to be dropped first (1815), and one
			// without reports the tablespace missing (1030). Measured on
			// mariadb:11.8.8.
			if hasSecondaryIndex(tbl) {
				return errImportNeedsIndexesDropped(db.Name, tbl.Name)
			}
			return errTablespaceMissing()
		}
		return errTablespaceExists(db.Name, tbl.Name)
	case nodes.ATDiscardTablespace:
		if err := refuseTablespaceEngine(tbl); err != nil {
			return err
		}
		tbl.TablespaceDiscarded = true
		return nil
	case nodes.ATRemovePartitioning:
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonpartitioned()
		}
		tbl.Partitioning = nil
		return nil
	case nodes.ATRebuildPartition, nodes.ATOptimizePartition, nodes.ATAnalyzePartition,
		nodes.ATCheckPartition, nodes.ATRepairPartition:
		// These leave the schema as it is; a table without partitions is
		// still refused (1505, measured on mariadb:11.8.8).
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonpartitioned()
		}
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

// alterAddsSystemVersioning reports whether a command of the statement is ADD
// SYSTEM VERSIONING.
func alterAddsSystemVersioning(cmds []*nodes.AlterTableCmd) bool {
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATAddSystemVersioning {
			return true
		}
	}
	return false
}

// addsRowBoundColumn reports whether an ADD, MODIFY or CHANGE gives a column
// GENERATED ALWAYS AS ROW START or ROW END.
func addsRowBoundColumn(cmd *nodes.AlterTableCmd) bool {
	switch cmd.Type {
	case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
	default:
		return false
	}
	for _, c := range append([]*nodes.ColumnDef{cmd.Column}, cmd.Columns...) {
		if c != nil && c.Generated != nil && c.Generated.RowBound != nodes.RowBoundNone {
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
			if cmd.IfNotExists && tableHasColumn(tbl, colDef.Name) {
				continue
			}
			if err := c.addSingleColumn(tbl, colDef, false, ""); err != nil {
				return err
			}
		}
		// The names these constraints give explicitly are reserved for each
		// other, as a CREATE TABLE's are.
		declared := declaredCheckNames(cmd.Constraints)
		for _, con := range cmd.Constraints {
			if err := c.alterAddConstraint(tbl, &nodes.AlterTableCmd{
				Type:       nodes.ATAddConstraint,
				Constraint: con,
			}, declared); err != nil {
				return err
			}
		}
		return nil
	}

	colDef := cmd.Column
	if colDef == nil {
		return nil
	}
	// ADD IF NOT EXISTS leaves a column already there as it is, with a note
	// (1060), and so does it for each column of the paren form. A plain ADD
	// of the name after it still fails: MariaDB 11.8.8 answers 1060 for
	// `ADD IF NOT EXISTS c INT, ADD c INT` and a note for the other order.
	if cmd.IfNotExists && tableHasColumn(tbl, colDef.Name) {
		return nil
	}

	return c.addSingleColumn(tbl, colDef, cmd.First, cmd.After)
}

// tableHasColumn reports whether the table has a column of the name.
func tableHasColumn(tbl *Table, name string) bool {
	_, ok := tbl.colByName[toLower(name)]
	return ok
}

// addSingleColumn adds one column definition to the table.
func (c *Catalog) addSingleColumn(tbl *Table, colDef *nodes.ColumnDef, first bool, after string) error {
	// The name is looked up before the definition is read, as MariaDB does:
	// measured on mariadb:11.8.8, `ALTER TABLE a ADD a ENUM('s', 's')` over a
	// table that already has a answers 1060 (the duplicate column), not 1291
	// (the duplicated ENUM value).
	colKey := toLower(colDef.Name)
	if _, exists := tbl.colByName[colKey]; exists {
		return errDupColumn(colDef.Name)
	}
	if err := validateColumnDefSemantics(colDef, indexedColumnsOf(tbl), c.effectiveColumnCollation(colDef, tbl.Collation)); err != nil {
		return err
	}
	// 4029, after the name is known to be new: the column's own expression may
	// not read it.
	if err := validateColumnSelfReference(colDef); err != nil {
		return err
	}

	col := buildColumnFromDef(tbl, colDef)
	isSerial := colDef.TypeName != nil && strings.EqualFold(colDef.TypeName.Name, "serial")
	if err := c.applyColumnEngineOptions(tbl, col, colDef); err != nil {
		return err
	}
	if col.AutoIncrement && tableHasOtherAutoIncrementColumn(tbl, "") {
		return errIncorrectAutoIncrementDefinition()
	}
	// A materialized generated column cannot read a non-deterministic one,
	// whichever statement adds it. See nondeterministicDependency.
	if col.Generated != nil && (col.Generated.Stored || columnDefHasKey(colDef)) {
		if dep := nondeterministicDependency(tbl, col, nil); dep != "" {
			return errUnusableInGeneratedColumn(dep, col.Name)
		}
	}
	if err := insertColumn(tbl, col, first, after); err != nil {
		return err
	}
	// SERIAL implies a unique key, as it does on the CREATE TABLE path.
	if isSerial {
		idxName := allocIndexName(tbl, col.Name)
		tbl.Indexes = append(tbl.Indexes, &Index{
			Name:    idxName,
			Table:   tbl,
			Columns: []*IndexColumn{{Name: col.Name}},
			Unique:  true,
			Visible: true,
		})
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      idxName,
			Type:      ConUniqueKey,
			Table:     tbl,
			Columns:   []string{col.Name},
			IndexName: idxName,
		})
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
				IndexType: longUniqueIndexType(tbl.GetColumn(colDef.Name)),
				Visible:   true,
			})
			tbl.Constraints = append(tbl.Constraints, &Constraint{
				Name:      idxName,
				Type:      ConUniqueKey,
				Table:     tbl,
				Columns:   []string{colDef.Name},
				IndexName: idxName,
			})
		case nodes.ColConstrReferences:
			// MariaDB makes a real foreign key of a column-level REFERENCES
			// the ALTER adds, as it does of the CREATE TABLE path's, and
			// judges it against the table the statement leaves: measured on
			// mariadb:11.8.8, `ALTER TABLE a ADD a INT REFERENCES a` over a
			// table whose keys do not lead with a answers 1005.
			conName := cc.Name
			if conName == "" {
				conName = fmt.Sprintf("%s_ibfk_%d", tbl.Name, nextFKGeneratedNumber(tbl, tbl.Name))
			}
			if foreignKeyConstraintNameExists(tbl.Database, nil, conName) {
				return errFKDupName(conName)
			}
			refDB, refTable := "", ""
			if cc.RefTable != nil {
				refDB = cc.RefTable.Schema
				refTable = cc.RefTable.Name
			}
			fkCols := []string{colDef.Name}
			fkCon := &Constraint{
				Name:        conName,
				Type:        ConForeignKey,
				Table:       tbl,
				Columns:     fkCols,
				RefDatabase: refDB,
				RefTable:    refTable,
				RefColumns:  refColumnsOrDefault(cc.RefColumns, fkCols),
				OnDelete:    refActionToString(cc.OnDelete),
				OnUpdate:    refActionToString(cc.OnUpdate),
			}
			if c.foreignKeyChecks {
				if err := c.validateSingleFK(tbl.Database, tbl, fkCon); err != nil {
					return err
				}
			}
			tbl.Constraints = append(tbl.Constraints, fkCon)
			// The backing index is named after the constraint only where the
			// statement named one, as the CREATE TABLE path's is, and a name
			// the table's indexes already hold is 1061, measured on
			// mariadb:11.8.8: `ADD a INT CONSTRAINT k REFERENCES p(id)` over
			// a table holding KEY k(b) answers it rather than adding a second
			// index called k.
			if !hasIndexCoveringColumns(tbl, fkCols) && cc.Name != "" && indexNameExists(tbl, cc.Name) {
				return errDupKeyName(cc.Name)
			}
			insertFKBackingIndex(tbl, cc.Name, fkCols, []*IndexColumn{{Name: colDef.Name}}, len(tbl.Indexes))
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
		// The primary key has more than one spelling, and the parser files
		// them under two commands: DROP PRIMARY KEY and the two DROP
		// CONSTRAINT forms are a dropped constraint named PRIMARY, while
		// DROP INDEX `PRIMARY` and DROP KEY `PRIMARY` are a dropped index of
		// that name. Measured on mariadb:11.8.8, every one of them lets the
		// same statement drop the key's columns, whichever order the clauses
		// come in.
		switch cmd.Type {
		case nodes.ATDropConstraint, nodes.ATDropIndex:
			if strings.EqualFold(cmd.Name, "PRIMARY") {
				return true
			}
		}
	}
	return false
}

// alterDroppedColumns returns the columns the statement drops, which is what
// decides whether a composite primary key survives it.
// alterDeclaredCheckNames collects the names this statement's added CHECK
// constraints give explicitly, lowered, which allocCheckName reserves while it
// names the unnamed ones: measured on mariadb:11.8.8,
// `ALTER TABLE t ADD CHECK (a > 0), ADD CONSTRAINT CONSTRAINT_1 CHECK (a < 9)`
// holds CONSTRAINT_2 and CONSTRAINT_1, where naming them in the order the
// commands come takes CONSTRAINT_1 twice.
func alterDeclaredCheckNames(cmds []*nodes.AlterTableCmd) map[string]bool {
	var declared map[string]bool
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATAddConstraint || cmd.Constraint == nil {
			continue
		}
		con := cmd.Constraint
		if con.Type != nodes.ConstrCheck || con.Name == "" {
			continue
		}
		if declared == nil {
			declared = map[string]bool{}
		}
		declared[toLower(con.Name)] = true
	}
	return declared
}

func alterDroppedColumns(cmds []*nodes.AlterTableCmd) map[string]bool {
	out := make(map[string]bool)
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATDropColumn {
			out[toLower(cmd.Name)] = true
		}
	}
	return out
}

func (c *Catalog) alterDropColumn(tbl *Table, cmd *nodes.AlterTableCmd, droppingPrimaryKey bool, droppedColumns map[string]bool) error {
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
	// part of that key survives the statement: MariaDB answers 1072, "Key
	// column '<col>' doesn't exist in table", rather than rebuilding the key
	// over the columns that are left. Where the statement drops every part the
	// key goes with them and the drop is accepted — measured on
	// mariadb:11.8.8, `PRIMARY KEY (a,b)` with `DROP COLUMN a, DROP COLUMN b`
	// in either order, and `PRIMARY KEY (a,b,c)` with all three, while
	// dropping two of those three answers 1072 for the first of them. The
	// statement may also drop the key itself, and then the column goes with
	// it: `DROP COLUMN a, DROP PRIMARY KEY` is accepted whichever order the
	// clauses come in. A single-column primary key goes with its column, and a
	// key that is not the primary one is rebuilt over what is left.
	if !droppingPrimaryKey {
		for _, idx := range tbl.Indexes {
			if !idx.Primary || len(idx.Columns) < 2 {
				continue
			}
			covers, survives, last := false, false, ""
			for _, ic := range idx.Columns {
				key := toLower(ic.Name)
				if key == colKey {
					covers = true
				}
				if droppedColumns[key] {
					// The column named is the last part of the key the
					// statement drops, in the key's own order and not the
					// statement's: measured, `PRIMARY KEY (a,b,c,d)` reports
					// 'b' for `DROP COLUMN a, DROP COLUMN b` and 'c' for
					// `DROP COLUMN c, DROP COLUMN a`.
					last = ic.Name
					continue
				}
				survives = true
			}
			if covers && survives {
				return &Error{
					Code:     1072,
					SQLState: "42000",
					Message:  fmt.Sprintf("Key column '%s' doesn't exist in table", last),
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
	if err := validateColumnDefSemantics(colDef, indexed, c.effectiveColumnCollation(colDef, tbl.Collation)); err != nil {
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
	// A generated column whose replacement expression reads it is 4029, for the
	// column that was generated already: turning a plain column into a generated
	// one is 1907 first, which the catalog does not model.
	if oldCol.Generated != nil {
		if err := validateColumnSelfReference(colDef); err != nil {
			return err
		}
	}

	col := buildColumnFromDef(tbl, colDef)
	if err := c.applyColumnEngineOptions(tbl, col, colDef); err != nil {
		return err
	}
	if col.AutoIncrement && tableHasOtherAutoIncrementColumn(tbl, oldName) {
		return errIncorrectAutoIncrementDefinition()
	}
	// The replacement is held to the same rule as a column added outright: a
	// materialized generated column cannot read a non-deterministic one.
	if col.Generated != nil && (col.Generated.Stored || indexed[toLower(col.Name)]) {
		if dep := nondeterministicDependency(tbl, col, nil); dep != "" {
			return errUnusableInGeneratedColumn(dep, col.Name)
		}
	}
	// So are the columns that read this one: replacing an expression can make
	// a materialized dependent invalid, and the table would then hold a schema
	// the server refuses. Measured on mariadb:11.8.8 with
	// `b DOUBLE AS (a) STORED` in place, `ALTER TABLE t MODIFY a DOUBLE AS
	// (RAND()) VIRTUAL` answers 1901 naming a, and so does the same table with
	// b VIRTUAL under a key, while an unmaterialized b takes it.
	//
	// The question is about the table as it would stand, so the column goes in
	// for the length of the check and comes back out. Nothing else has been
	// touched yet, so a refusal leaves the table as it was.
	if col.Generated != nil {
		prev := tbl.Columns[idx]
		tbl.Columns[idx] = col
		depErr := validateMaterializedGeneratedColumns(tbl, indexed)
		tbl.Columns[idx] = prev
		if depErr != nil {
			return depErr
		}
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

	// SERIAL implies a unique key here too: a replacement keeps the keys the
	// old definition had and adds the new one's, so MODIFY of a SERIAL column
	// leaves two unique keys. Measured on mariadb:11.8.8.
	if colDef.TypeName != nil && strings.EqualFold(colDef.TypeName.Name, "serial") {
		idxName := allocIndexName(tbl, col.Name)
		tbl.Indexes = append(tbl.Indexes, &Index{
			Name:    idxName,
			Table:   tbl,
			Columns: []*IndexColumn{{Name: col.Name}},
			Unique:  true,
			Visible: true,
		})
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      idxName,
			Type:      ConUniqueKey,
			Table:     tbl,
			Columns:   []string{col.Name},
			IndexName: idxName,
		})
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
		// A column the table already holds is not read for its keys at all: the
		// duplicate name answers 1060 first, measured on mariadb:11.8.8, where
		// `ALTER TABLE a ADD a INT KEY KEY` over a table that has a answers 1060
		// and the same statement over one that does not answers 1068.
		columns := cmd.Columns
		if cmd.Column != nil {
			columns = append([]*nodes.ColumnDef{cmd.Column}, columns...)
		}
		for _, col := range columns {
			// A plain ADD of a column the table already holds fails 1060 before
			// its keys are read (measured on mariadb:11.8.8), where MODIFY,
			// CHANGE and ADD IF NOT EXISTS — which leaves the column as it is —
			// still count them.
			if col == nil || (cmd.Type == nodes.ATAddColumn && !cmd.IfNotExists && tableHasColumn(tbl, col.Name)) {
				continue
			}
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

// validateAddedKeyColumns refuses a key part that names a column the table
// does not have: mariadb:11.8.8 answers 1072 for ADD KEY, FULLTEXT, SPATIAL,
// UNIQUE and FOREIGN KEY alike. A functional part is an expression, and a
// WITHOUT OVERLAPS part names the application-time period; neither is a
// column, and the checks that judge them run later.
func validateAddedKeyColumns(tbl *Table, cols []*nodes.IndexColumn) error {
	for _, ic := range cols {
		if ic == nil || ic.Functional || ic.WithoutOverlaps {
			continue
		}
		ref, ok := ic.Expr.(*nodes.ColumnRef)
		if !ok || tbl.GetColumn(ref.Column) != nil {
			continue
		}
		// A part with no position did not come from SQL: LoadMetadata keeps a
		// key expression it could not parse as a name, so that the index still
		// installs, and that name is no column the statement failed to find.
		if ref.Loc.End == 0 {
			continue
		}
		return errKeyColumnDoesNotExist(ref.Column)
	}
	return nil
}

func (c *Catalog) alterAddConstraint(tbl *Table, cmd *nodes.AlterTableCmd, declaredChecks map[string]bool) error {
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
	// A FULLTEXT key's parser is looked up before its columns: mariadb:11.8.8
	// answers 1128 for `ADD FULLTEXT(a) WITH PARSER a` even where a is not a
	// column, and 1072 only for a key without a parser.
	if con.Type == nodes.ConstrFulltextIndex {
		if err := validateFulltextParser(con.IndexOptions); err != nil {
			return err
		}
	}
	if err := validateAddedKeyColumns(tbl, con.IndexColumns); err != nil {
		return err
	}
	// A type the server answers 1210 for is refused before the key is built,
	// after the columns are looked up: measured on mariadb:11.8.8,
	// `ADD KEY USING RTREE(a)` answers 1210 where `ADD KEY USING RTREE(nosuch)`
	// answers 1072.
	if err := validateKeyTypeArguments(resolveConstraintIndexType(con)); err != nil {
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
		idxCols := buildIndexColumns(con)
		if err := validatePrimaryKeyColumns(tbl, cols); err != nil {
			return err
		}
		// The key's parts are held to their types before a repeated name is 1060:
		// measured on mariadb:11.8.8, `PRIMARY KEY (a(1), a)` over an INT column
		// answers 1089 and `PRIMARY KEY (a)` over a TEXT one 1170.
		if err := validateIndexColumns(tbl, idxCols, false, false); err != nil {
			return err
		}
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
			return err
		}
		// Mark PK columns as NOT NULL.
		for _, colName := range cols {
			col := tbl.GetColumn(colName)
			if col != nil {
				col.Nullable = false
			}
		}
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
		long, err := validateUniqueIndexColumns(tbl, idxCols)
		if err != nil {
			return err
		}
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
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
		if long {
			uqIdx.IndexType = "HASH"
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
			RefColumns:  refColumnsOrDefault(con.RefColumns, cols),
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
		// The backing index a foreign key brings with it materializes a
		// VIRTUAL generated column as any other key does, so the column is
		// held to the same rule: measured on mariadb:11.8.8,
		// `ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)`
		// over `a DOUBLE AS (RAND()) VIRTUAL` answers 1901, and over a column
		// reading one it names that column. Checked before the constraint
		// goes in, so a refusal leaves the table as it was.
		fkIdxCols := buildIndexColumns(con)
		if err := validateIndexColumns(tbl, fkIdxCols, false, false); err != nil {
			return err
		}
		tbl.Constraints = append(tbl.Constraints, fkCon)
		// Add implicit backing index for FK if needed.
		ensureFKBackingIndex(tbl, fkName(con), cols, fkIdxCols)

	case nodes.ConstrCheck:
		conName := con.Name
		if conName == "" {
			conName = allocCheckName(tbl.Database, tbl, declaredChecks)
		}
		if checkConstraintNameExists(tbl.Database, tbl, conName) {
			return errCheckConstraintDupName(conName)
		}
		if err := validateCheckExpr(conName, con.Expr); err != nil {
			return err
		}
		// The check names columns of the table: measured on mariadb:11.8.8,
		// `ALTER TABLE a ADD CHECK(a)` over a table without column a answers
		// 1054 "Unknown column 'a' in 'CHECK'".
		if err := validateExprColumns(tbl, con.Expr, "CHECK"); err != nil {
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
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
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
		if err := validateFulltextParser(con.IndexOptions); err != nil {
			return err
		}
		if err := validateFulltextColumns(tbl, idxCols); err != nil {
			return err
		}
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
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
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
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
		if err := validateNoRepeatedKeyColumns(idxCols); err != nil {
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

// moveColumnChecks follows a renamed column: the check it carries belongs to
// it, so the ownership moves to the new name. Every check the column carries
// moves, not only the json_valid one the catalog wrote — ownership left behind
// would make the check the next column of that name's, and go with that one
// when it is dropped.
//
// The name stays as it was: measured on mariadb:11.8.8, after
// `CREATE TABLE t (a INT CHECK (a > 0)); ALTER TABLE t RENAME COLUMN a TO z`
// the check is written on z and reads CHECK (`z` > 0), while
// information_schema.CHECK_CONSTRAINTS still calls it a — and still a after
// a second rename to w. A JSON column's check behaves the same way.
//
// The expression is rewritten only for the json_valid check, which the catalog
// wrote itself. Rewriting the column references inside a check a statement
// declared is the general stored-expression rewrite the catalog does not do:
// measured, the server turns CHECK (`a` > 0) into CHECK (`z` > 0) there too.
//
// Checks are found by the column they belong to, since another constraint may
// carry the old name.
func moveColumnChecks(tbl *Table, oldName, newName string) {
	key := toLower(oldName)
	for _, con := range tbl.Constraints {
		if !isColumnCheck(con) || toLower(con.Column) != key {
			continue
		}
		con.Column = newName
		if con.Implicit {
			con.CheckExpr = "json_valid(" + backtickIDQ(newName) + ")"
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

	// The key the statement renames is looked up before the name it gives is
	// judged: measured on mariadb:11.8.8, `ALTER TABLE a RENAME KEY a TO idx_d`
	// over a table holding idx_d alone answers 1176, not 1061, so a missing key
	// is the answer whatever the new name holds.
	var found *Index
	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == oldKey {
			found = idx
			break
		}
	}
	if found == nil {
		return errKeyDoesNotExist(cmd.Name, tbl.Name)
	}

	if newKey != oldKey {
		if err := validateNonPrimaryIndexName(cmd.NewName); err != nil {
			return err
		}
		if indexNameExists(tbl, cmd.NewName) {
			return errDupKeyName(cmd.NewName)
		}
	}

	found.Name = cmd.NewName
	if newKey != oldKey {
		renameFunctionalIndexHiddenColumns(tbl, found, cmd.NewName)
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

	// Tables, views and sequences share one name space.
	if (newDB != oldDB || newKey != oldKey) && (newDB.Tables[newKey] != nil || newDB.Views[newKey] != nil || newDB.Sequences[newKey] != nil) {
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

	// An engine-defined option is the engine's by its form, whatever its name.
	if opt.EngineDefined {
		return tbl.setEngineOption(opt.Name, opt.Value, opt.StringValue)
	}
	switch toLower(opt.Name) {
	case "engine":
		tbl.Engine = canonicalEngineName(opt.Value)
	case "charset", "character set", "default charset", "default character set":
		tbl.Charset = normalizeCharsetName(opt.Value)
		// Update collation to the default for this charset.
		if defColl := c.charsetCollation(tbl.Charset); defColl != "" {
			tbl.Collation = defColl
		}
	case "collate", "default collate":
		tbl.Collation = opt.Value
	case "comment":
		tbl.Comment = opt.Value
	case "auto_increment":
		fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
	case "row_format":
		// DEFAULT clears the option: MariaDB prints no ROW_FORMAT for
		// `ALTER TABLE a ROW_FORMAT DEFAULT`, measured on mariadb:11.8.8.
		tbl.RowFormat = canonicalRowFormat(opt.Value)
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
		tbl.MinRows = canonicalRows(opt.Value)
	case "max_rows":
		tbl.MaxRows = canonicalRows(opt.Value)
	case "avg_row_length":
		tbl.AvgRowLength = opt.Value
	case "tablespace":
		// MariaDB 11.8.8 accepts the option and prints no TABLESPACE for it,
		// so the catalog keeps none. Measured on `ALTER TABLE a TABLESPACE a`
		// and on a CREATE TABLE with it.
	case "pack_keys":
		tbl.PackKeys = opt.Value
	case "checksum", "table_checksum":
		tbl.Checksum = normalizedChecksum(opt.Value)
	case "delay_key_write":
		tbl.DelayKeyWrite = opt.Value
	case "page_checksum":
		tbl.PageChecksum = choiceValue(opt.Value)
	case "transactional":
		tbl.Transactional = choiceValue(opt.Value)
	case "sequence":
		if opt.Value == "1" {
			return errSequenceInvalidTableStructure(tbl.Database.Name, tbl.Name)
		}
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
		// A default expression that reads the column it defaults is 4029, after
		// the column is known to exist (1054).
		key := toLower(col.Name)
		for _, ref := range generatedColumnRefs(cmd.DefaultExpr) {
			if toLower(ref) == key {
				return errUninitializedField(col.Name, col.Name)
			}
		}
		if err := validateDefaultExprForExistingColumn(col, cmd.DefaultExpr); err != nil {
			return err
		}
		setColumnDefaultFromExpr(col, qualifySequenceCalls(cmd.DefaultExpr, tbl))
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
		// SERIAL is BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, as it is on the
		// CREATE TABLE path; the caller adds its unique key.
		if strings.EqualFold(colDef.TypeName.Name, "serial") {
			col.DataType = "bigint"
			col.ColumnType = "bigint(20) unsigned"
			col.AutoIncrement = true
			col.Nullable = false
		}
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
		// A string type with CHARACTER SET binary is the binary type: MariaDB
		// stores `ADD a CHAR BYTE` and `ADD a CHAR COLLATE BINARY` as
		// `a` binary(1), measured on mariadb:11.8.8.
		if strings.EqualFold(col.Charset, "binary") && isStringType(col.DataType) && !isEnumSetType(col.DataType) {
			col = convertToBinaryType(col, colDef.TypeName)
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
				if dc := versionCollation(tbl.Database, col.Charset); dc != "" {
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
		setColumnDefaultFromExpr(col, qualifySequenceCalls(colDef.DefaultValue, tbl))
	}
	if colDef.OnUpdate != nil {
		setColumnOnUpdateFromExpr(col, colDef.OnUpdate)
	}
	if colDef.Generated != nil {
		col.Generated = &GeneratedColumnInfo{
			Expr:             nodeToSQL(colDef.Generated.Expr),
			Stored:           colDef.Generated.Stored,
			RowBound:         rowBoundSQL(colDef.Generated.RowBound),
			Nondeterministic: nondeterministicFuncIn(colDef.Generated.Expr),
			Refs:             generatedColumnRefs(colDef.Generated.Expr),
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
				setColumnDefaultFromExpr(col, qualifySequenceCalls(cc.Expr, tbl))
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
	moveColumnChecks(tbl, oldName, newName)
	// What a generated expression reads follows the column too. The server
	// rewrites the expression's own text, which this does not do yet, but the
	// recorded names have to keep pointing at the column or the dependency
	// goes unseen: measured on mariadb:11.8.8, after
	// `CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED);
	// ALTER TABLE t RENAME COLUMN a TO z`, `ALTER TABLE t MODIFY z DOUBLE AS
	// (RAND()) VIRTUAL` answers 1901 naming z.
	for _, col := range tbl.Columns {
		if col.Generated == nil {
			continue
		}
		// A fresh slice rather than a write through this one, which a
		// snapshot of the table may still be sharing.
		refs := make([]string, len(col.Generated.Refs))
		for i, ref := range col.Generated.Refs {
			if toLower(ref) == oldKey {
				ref = newName
			}
			refs[i] = ref
		}
		col.Generated.Refs = refs
	}
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
		if defColl := c.charsetCollation(charset); defColl != "" {
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
		return errPartitionMgmtOnNonpartitioned()
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
	if err := checkAddedPartitionKinds(tbl.Partitioning, cmd.PartitionDefs); err != nil {
		return err
	}
	if err := checkAddedRangeBounds(tbl, cmd.PartitionDefs); err != nil {
		return err
	}
	for _, pd := range cmd.PartitionDefs {
		pdi := &PartitionDefInfo{
			Name: pd.Name,
		}
		if pd.Values != nil {
			pdi.ValueExpr = partitionValueToString(pd.Values, partitionTypeFromString(tbl.Partitioning.Type))
			pdi.bound = pd.Values
		}
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			// STORAGE ENGINE is the other spelling of ENGINE, and names the
			// same engine.
			case "engine", "storage engine":
				pdi.Engine = canonicalEngineName(opt.Value)
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
		return errPartitionMgmtOnNonpartitioned()
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
		return errPartitionMgmtOnNonpartitioned()
	}
	// Truncate is a data operation; for DDL catalog purposes, it's a no-op.
	return nil
}

// alterCoalescePartition reduces the number of partitions for HASH/KEY partitioned tables.
func (c *Catalog) alterCoalescePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonpartitioned()
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
		return errPartitionMgmtOnNonpartitioned()
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
			pdi.bound = pd.Values
		}
		for _, opt := range pd.Options {
			switch toLower(opt.Name) {
			// STORAGE ENGINE is the other spelling of ENGINE, and names the
			// same engine.
			case "engine", "storage engine":
				pdi.Engine = canonicalEngineName(opt.Value)
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

// validateAddedColumnOrder refuses, before any of them is added, a generated column
// the statement adds whose expression reads a generated column the same statement adds
// after it: 4029 naming the field read, as CREATE TABLE answers it. Measured on
// mariadb:11.8.8: `ALTER TABLE t ADD c INT AS(d), ADD d INT AS(c)` is 4029, and
// `ADD c INT AS(d), ADD d INT` runs.
//
// MariaDB reads every column first, so a name the table already holds is 1060 before
// this — `ADD b INT AS(a), ADD a INT AS(b)` over a table that has a and b is 1060 — and
// so is a definition's own semantics, as `ADD a INT AS(a), ADD c INT NOT NULL DEFAULT
// 'x'` answers 1067 for c. The statement's own commands then judge nothing again: what
// a single ADD answers is checked where it is added (see addSingleColumn).
func (c *Catalog) validateAddedColumnOrder(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	type addedColumn struct {
		def         *nodes.ColumnDef
		ifNotExists bool
	}
	var added []addedColumn
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATAddColumn {
			continue
		}
		defs := cmd.Columns
		if cmd.Column != nil {
			defs = append([]*nodes.ColumnDef{cmd.Column}, defs...)
		}
		for _, def := range defs {
			if def != nil {
				added = append(added, addedColumn{def: def, ifNotExists: cmd.IfNotExists})
			}
		}
	}
	generated := false
	for _, a := range added {
		if a.def.Generated != nil {
			generated = true
			break
		}
	}
	if !generated {
		return nil
	}
	for _, a := range added {
		if !a.ifNotExists && tableHasColumn(tbl, a.def.Name) {
			return errDupColumn(a.def.Name)
		}
	}
	indexed := indexedColumnsOf(tbl)
	for _, a := range added {
		if err := validateColumnDefSemantics(a.def, indexed, c.effectiveColumnCollation(a.def, tbl.Collation)); err != nil {
			return err
		}
	}
	for i, a := range added {
		later := make([]*nodes.ColumnDef, 0, len(added)-i-1)
		for _, b := range added[i+1:] {
			later = append(later, b.def)
		}
		if err := validateColumnGeneratedOrder(a.def, later); err != nil {
			return err
		}
	}
	return nil
}

// onlyPartitionMaintenance reports whether every command of an ALTER TABLE is one of
// the partition maintenance commands, which a missing table takes silently.
func onlyPartitionMaintenance(cmds []*nodes.AlterTableCmd) bool {
	if len(cmds) == 0 {
		return false
	}
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATOptimizePartition, nodes.ATAnalyzePartition, nodes.ATCheckPartition, nodes.ATRepairPartition:
		default:
			return false
		}
	}
	return true
}

// writesFixedRowFormat reports whether the statement writes ROW_FORMAT FIXED or PAGE
// and leaves the table on an engine that does not take it, which is 1478. The engine
// is the one the statement leaves, so `ROW_FORMAT FIXED ENGINE=MyISAM` is taken.
func writesFixedRowFormat(tbl *Table, cmds []*nodes.AlterTableCmd) bool {
	engine, written := tbl.Engine, false
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATTableOption || cmd.Option == nil {
			continue
		}
		switch {
		case eqFoldStr(cmd.Option.Name, "ENGINE"):
			engine = canonicalEngineName(cmd.Option.Value)
		case eqFoldStr(cmd.Option.Name, "ROW_FORMAT"):
			written = writesFixedRowFormatValue(canonicalRowFormat(cmd.Option.Value))
		}
	}
	return written && !engineTakesFixedRowFormat(engine)
}

func writesFixedRowFormatValue(value string) bool {
	switch strings.ToLower(value) {
	case "fixed", "page":
		return true
	}
	return false
}

// engineTakesFixedRowFormat reports whether an engine takes ROW_FORMAT FIXED and PAGE,
// which MyISAM, Aria, MEMORY, CSV and MRG_MyISAM do on mariadb:11.8.8. A table that
// names no engine is InnoDB's.
func engineTakesFixedRowFormat(engine string) bool {
	if engine == "" {
		return false
	}
	switch strings.ToLower(engine) {
	case "myisam", "aria", "memory", "csv", "mrg_myisam":
		return true
	}
	return false
}

// validateAlteredRowFormat refuses the ROW_FORMAT the table was left with when its engine
// does not take it: InnoDB answers 1478, "Table storage engine 'InnoDB' does not support
// the create option 'ROW_TYPE'" (measured on mariadb:11.8.8), where MyISAM, Aria,
// MEMORY, CSV and MRG_MyISAM take FIXED and PAGE. writesFixedRowFormat decides whether
// the statement wrote one; a row format the table had already and an engine change
// alone are left out, as MariaDB answers 1005 for those and not 1478.
func validateAlteredRowFormat(tbl *Table) error {
	engine := tbl.Engine
	if engine == "" {
		engine = "InnoDB"
	}
	return &Error{Code: 1478, SQLState: "HY000",
		Message: fmt.Sprintf("Table storage engine '%s' does not support the create option 'ROW_TYPE'", engine)}
}

// alterExchangePartition exchanges a partition with a non-partitioned table.
func (c *Catalog) alterExchangePartition(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonpartitioned()
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

// alterConvertPartitionToTable runs CONVERT PARTITION p TO TABLE t: the partition
// leaves the table and becomes t, a table of the same definition with no
// partitioning. Measured on mariadb:11.8.8, which answers 1505 on a table with no
// partitions, 1512 on one partitioned other than by RANGE or LIST, 4184 on one
// with subpartitions, 1507 for a
// partition it does not have, 1508 for its last partition and 1050 for a table t
// that exists.
func (c *Catalog) alterConvertPartitionToTable(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonpartitioned()
	}
	if !strings.HasPrefix(tbl.Partitioning.Type, "RANGE") && !strings.HasPrefix(tbl.Partitioning.Type, "LIST") {
		return errOnlyOnRangeListPartitioning("CONVERT PARTITION")
	}
	if tbl.Partitioning.SubType != "" {
		return &Error{Code: 4184, SQLState: "HY000", Message: "Convert partition is not supported for subpartitioned table."}
	}
	at := -1
	for i, pd := range tbl.Partitioning.Partitions {
		if strings.EqualFold(pd.Name, cmd.Name) {
			at = i
		}
	}
	if at < 0 {
		return errDropPartitionNonExistent()
	}
	if len(tbl.Partitioning.Partitions) == 1 {
		return errDropLastPartition()
	}
	target, err := c.resolveDatabase(cmd.ConvertTable.Schema)
	if err != nil {
		return err
	}
	key := toLower(cmd.ConvertTable.Name)
	switch {
	case target.Tables[key] == tbl:
		// The table itself: MariaDB 11.8.8 runs it, and SHOW CREATE TABLE then
		// shows the table without the partition and no other table. The
		// server's files are not right after it (a later DROP TABLE leaves a
		// table behind), so no oracle test pins it.
	case target.Tables[key] != nil || target.Views[key] != nil:
		return errDupTable(cmd.ConvertTable.Name)
	default:
		if err := c.createTableLike(target, cmd.ConvertTable.Name, key, tbl, &nodes.CreateTableStmt{}); err != nil {
			return err
		}
	}
	parts := tbl.Partitioning.Partitions
	tbl.Partitioning.Partitions = append(parts[:at:at], parts[at+1:]...)
	tbl.Partitioning.NumParts = len(tbl.Partitioning.Partitions)
	return nil
}

// alterConvertTableToPartition runs CONVERT TABLE t TO PARTITION part_definition:
// t becomes a new partition of the table and is gone. Measured on mariadb:11.8.8,
// which answers 1505 on a table with no partitions, 1734 on one with
// subpartitions, 1146 for a table t that does
// not exist, 1736 for one whose definition is not the table's or that has
// partitions of its own, and 1517 for a partition name the table has.
func (c *Catalog) alterConvertTableToPartition(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonpartitioned()
	}
	if tbl.Partitioning.SubType != "" {
		return &Error{Code: 1734, SQLState: "HY000", Message: "Subpartitioned table, use subpartition instead of partition"}
	}
	source, err := c.resolveDatabase(cmd.ConvertTable.Schema)
	if err != nil {
		return err
	}
	src := source.GetTable(cmd.ConvertTable.Name)
	if src == nil {
		return errNoSuchTable(source.Name, cmd.ConvertTable.Name)
	}
	// The source is a table of no partitions: a partitioned one is 1736 whatever
	// its columns, and so is a definition other than the table's.
	if src.Partitioning != nil || c.tableDefinitionText(src) != c.tableDefinitionText(tbl) {
		return errTablesDifferentMetadata()
	}
	// The new partition's name is one the table does not have (1517), and its
	// bound comes after the table's last (1493, 1481), checked before either
	// table changes.
	for _, pd := range cmd.PartitionDefs {
		for _, have := range tbl.Partitioning.Partitions {
			if strings.EqualFold(have.Name, pd.Name) {
				return errSameNamePartition(pd.Name)
			}
		}
	}
	// alterAddPartition orders the bound before it adds anything.
	if err := c.alterAddPartition(tbl, &nodes.AlterTableCmd{PartitionDefs: cmd.PartitionDefs}); err != nil {
		return err
	}
	c.dropReplacedTable(source, src)
	return nil
}

// tableDefinitionText is SHOW CREATE TABLE of the table without its name and its
// partitioning, which is what CONVERT TABLE holds the two tables to. The table is
// printed with no partitioning rather than cut at the clause, which a comment can
// spell too, and loses the exact head that names it, which a quoted name can make
// look like anything.
func (c *Catalog) tableDefinitionText(tbl *Table) string {
	partitioning := tbl.Partitioning
	tbl.Partitioning = nil
	text := c.ShowCreateTable(tbl.Database.Name, tbl.Name)
	tbl.Partitioning = partitioning
	text = strings.TrimPrefix(text, "CREATE TABLE "+backtickIDQ(tbl.Name))
	return strings.TrimPrefix(text, "CREATE TEMPORARY TABLE "+backtickIDQ(tbl.Name))
}

// checkAddedPartitionKinds holds each partition ADD PARTITION and CONVERT TABLE
// add to the bound its table's partitioning takes: VALUES LESS THAN on RANGE,
// VALUES IN or DEFAULT on LIST, HISTORY or CURRENT on SYSTEM_TIME, and none on
// HASH and KEY. Measured on mariadb:11.8.8: VALUES IN or DEFAULT on another is
// 1480, LESS THAN on another is 1480, no bound on RANGE or LIST is 1479, and
// HISTORY or CURRENT on another is 4113.
func checkAddedPartitionKinds(pi *PartitionInfo, defs []*nodes.PartitionDef) error {
	isRange := strings.HasPrefix(pi.Type, "RANGE")
	isList := strings.HasPrefix(pi.Type, "LIST")
	for _, pd := range defs {
		def, _ := pd.Values.(*nodes.String)
		isDefault := def != nil && strings.EqualFold(def.Str, "DEFAULT")
		switch {
		case pd.SystemTime != "":
			if pi.Type != "SYSTEM_TIME" {
				return &Error{Code: 4113, SQLState: "HY000", Message: "Only SYSTEM_TIME PARTITIONING can use HISTORY or CURRENT in partition definition"}
			}
		case pd.ValuesIn || isDefault:
			if !isList {
				return &Error{Code: 1480, SQLState: "HY000", Message: "Only LIST PARTITIONING can use VALUES IN in partition definition"}
			}
		case pd.Values != nil:
			if !isRange {
				return &Error{Code: 1480, SQLState: "HY000", Message: "Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition"}
			}
		case isRange:
			return &Error{Code: 1479, SQLState: "HY000", Message: "Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition"}
		case isList:
			return &Error{Code: 1479, SQLState: "HY000", Message: "Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition"}
		}
		if !isDefault && pd.Values != nil && !boundFitsColumns(pi, pd) {
			return &Error{Code: 1653, SQLState: "HY000", Message: "Inconsistency in usage of column lists for partitioning"}
		}
	}
	return nil
}

// boundFitsColumns reports whether a partition's bound has the width its table's
// partitioning takes: one value per COLUMNS column, one value with no COLUMNS, and
// for LIST each value of that width. MariaDB 11.8.8 answers 1653 for LESS THAN (20)
// on RANGE COLUMNS (a, b), for LESS THAN (20, 30) on RANGE (a), for IN (3) on LIST
// COLUMNS (a, b) and for IN ((3, 4)) on LIST (a).
func boundFitsColumns(pi *PartitionInfo, pd *nodes.PartitionDef) bool {
	width := len(pi.Columns)
	if width == 0 {
		width = 1
	}
	if !pd.ValuesIn {
		return len(partitionBoundItems(pd.Values)) == width
	}
	for _, item := range partitionBoundItems(pd.Values) {
		if len(partitionBoundItems(item)) != width {
			return false
		}
	}
	return true
}

// checkAddedRangeBounds holds the partitions ADD PARTITION and CONVERT TABLE add
// to a LIST table to values of their own (checkAddedListValues), and to a RANGE or
// RANGE COLUMNS table to bounds that go on rising from the table's
// last: past a MAXVALUE partition is 1481 on RANGE, and a bound not above the one
// before is 1493. Measured on mariadb:11.8.8, which answers 1493 for
// `ADD PARTITION (PARTITION p1 VALUES LESS THAN (5))` after LESS THAN (10), and
// for a bound equal to it. A bound the catalog cannot order passes.
func checkAddedRangeBounds(tbl *Table, defs []*nodes.PartitionDef) error {
	pi := tbl.Partitioning
	if pi != nil && (pi.Type == "LIST" || pi.Type == "LIST COLUMNS") {
		return checkAddedListValues(pi, defs)
	}
	if pi == nil || (pi.Type != "RANGE" && pi.Type != "RANGE COLUMNS") {
		return nil
	}
	var prev nodes.Node
	if n := len(pi.Partitions); n > 0 {
		prev = pi.Partitions[n-1].bound
	}
	kinds := rangeColumnsBoundKinds(tbl, pi.Columns)
	for _, pd := range defs {
		cur := pd.Values
		if prev != nil && cur != nil {
			if pi.Type == "RANGE COLUMNS" {
				if rangeColumnsNotIncreasing(prev, cur, kinds) {
					return errRangeNotIncreasing()
				}
			} else {
				p, c := partitionBoundItems(prev), partitionBoundItems(cur)
				if len(p) == 1 && isPartitionMaxvalue(p[0]) {
					return &Error{Code: 1481, SQLState: "HY000", Message: "MAXVALUE can only be used in last partition definition"}
				}
				if len(p) == 1 && len(c) == 1 {
					pv, ok1 := partitionBoundInt(p[0])
					cv, ok2 := partitionBoundInt(c[0])
					if ok1 && ok2 && cv <= pv {
						return errRangeNotIncreasing()
					}
				}
			}
		}
		prev = cur
	}
	return nil
}

// checkAddedListValues holds the partitions added to a LIST or LIST COLUMNS table
// to values no partition holds yet, theirs or each other's: MariaDB 11.8.8 answers
// 1495 for `ADD PARTITION (PARTITION p1 VALUES IN (1))` beside IN (1, 2), for two
// new partitions of one value, for a tuple again and for NULL again. A partition
// whose values the catalog does not hold is passed over.
func checkAddedListValues(pi *PartitionInfo, defs []*nodes.PartitionDef) error {
	seen := make(map[string]bool)
	add := func(bound nodes.Node) bool {
		for _, item := range partitionBoundItems(bound) {
			key := listValueKey(item)
			if key == "" {
				continue
			}
			if seen[key] {
				return false
			}
			seen[key] = true
		}
		return true
	}
	for _, have := range pi.Partitions {
		if have.bound != nil {
			add(have.bound)
		}
	}
	for _, pd := range defs {
		if pd.Values != nil && !add(pd.Values) {
			return &Error{Code: 1495, SQLState: "HY000", Message: "Multiple definition of same constant in list partitioning"}
		}
	}
	return nil
}

// listValueKey is a LIST value as checkAddedListValues compares it: an integer by
// its value, a tuple by its members, a tuple of one as its member, and anything
// else by its text as written; "" for a value it cannot compare. Text is compared
// exactly, since whether 'a' and 'A' are one value is the column's collation's to
// say: a utf8mb4_bin column holds both, and the catalog passes what it cannot
// settle.
func listValueKey(n nodes.Node) string {
	if l, ok := n.(*nodes.List); ok {
		if len(l.Items) == 1 {
			return listValueKey(l.Items[0])
		}
		keys := make([]string, len(l.Items))
		for i, item := range l.Items {
			if keys[i] = listValueKey(item); keys[i] == "" {
				return ""
			}
		}
		return "(" + strings.Join(keys, ", ") + ")"
	}
	if v, ok := partitionBoundInt(n); ok {
		return strconv.FormatInt(v, 10)
	}
	// An integer partitionBoundInt declines is one the lexer clamped, so two of
	// them written apart can read alike: not comparable.
	if _, ok := n.(*nodes.IntLit); ok {
		return ""
	}
	switch v := n.(type) {
	case nodes.ExprNode:
		return nodeToSQL(v)
	case *nodes.String:
		return v.Str
	}
	return ""
}

// errRangeNotIncreasing is 1493, a RANGE bound not above the one before it.
func errRangeNotIncreasing() error {
	return &Error{Code: 1493, SQLState: "HY000", Message: "VALUES LESS THAN value must be strictly increasing for each partition"}
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
	if cmd == nil {
		return false
	}
	if cmd.Type == nodes.ATModifyColumn || cmd.Type == nodes.ATChangeColumn {
		return cmd.Column != nil && columnNeedsRollback(cmd.Column)
	}
	if cmd.Type != nodes.ATAddColumn {
		return false
	}
	if len(cmd.Columns)+len(cmd.Constraints) > 1 {
		return true
	}
	if cmd.Column != nil && columnNeedsRollback(cmd.Column) {
		return true
	}
	for _, col := range cmd.Columns {
		if col != nil && columnNeedsRollback(col) {
			return true
		}
	}
	return false
}

// columnNeedsRollback reports whether a column the statement adds carries
// something judged after it is in the table, so the ADD has to be undone when
// that judgement fails: a constraint, or a DEFAULT expression whose columns are
// read against the table with the column in it.
func columnNeedsRollback(col *nodes.ColumnDef) bool {
	return len(col.Constraints) > 0 || col.DefaultValue != nil
}
