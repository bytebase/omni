package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func (c *Catalog) alterTable(stmt *nodes.AlterTableStmt) error {
	// A NOW()-family DEFAULT the column's type cannot hold is 1067 while TiDB
	// prepares the statement, before it resolves the database and before it
	// looks the table up: `ALTER TABLE nodb.a ADD x INT DEFAULT NOW()` is 1067
	// rather than 1049 or 1146 (measured on v8.5.5). A charset or collation a
	// column definition names is resolved as early: `ADD x TINYTEXT CHARACTER
	// SET ucs2` is 1115 even with no such table (measured on v8.5.5).
	if err := checkAlterUnsupportedTableOptions(stmt.Commands); err != nil {
		return err
	}
	if err := checkAlterNowDefaults(stmt.Commands); err != nil {
		return err
	}
	if err := checkAlterColumnCharsets(stmt.Commands); err != nil {
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
	key := toLower(tableName)
	if db.TempTables[key] != nil {
		return errLocalTemporaryDDL("ALTER TABLE")
	}
	tbl := db.Tables[key]
	if tbl == nil {
		if db.Views[key] != nil {
			// TiDB v8.5.5 runs ALTER TABLE <view> COMPACT [TIFLASH REPLICA]
			// without error, answers 1505 for ANALYZE PARTITION ... INDEX and
			// 1105 for ANALYZE PARTITION without INDEX, and answers 1347
			// ("is not BASE TABLE") for every other ALTER TABLE on a view.
			if len(stmt.Commands) == 1 && stmt.Commands[0] != nil {
				switch stmt.Commands[0].Type {
				case nodes.ATCompact:
					if len(stmt.Commands[0].PartitionNames) > 0 {
						return errCompactNonPartitionedTable(tableName, stmt.Commands[0].PartitionNames)
					}
					return nil
				case nodes.ATAnalyzePartition:
					if stmt.Commands[0].AnalyzeIndexFlag {
						return errPartitionMgmtOnNonPartitioned()
					}
					return &Error{Code: ErrUnknownError, SQLState: sqlState(ErrUnknownError),
						Message: fmt.Sprintf("analyze view %s is not supported now", tableName)}
				}
			}
			return errWrongObjectKind(db.Name, tableName, "BASE TABLE")
		}
		if db.Sequences[key] != nil {
			return errWrongObjectKind(db.Name, tableName, "BASE TABLE")
		}
		return errNoSuchTable(db.Name, tableName)
	}

	if len(stmt.Commands) == 0 {
		return nil
	}
	// A cached table takes CACHE and NOCACHE and nothing else, one command at a
	// time: v8.5.5 answers 8242 to every other ALTER, `NOCACHE, ADD COLUMN`
	// included (read on v8.5.5).
	if tbl.Cached && !singleCacheCommand(stmt.Commands) {
		return errCachedTable("Alter Table")
	}
	// The foreign keys that reference the table follow its renames once the
	// whole statement holds (pendingRetargets), so a refused statement leaves
	// them as they were.
	c.pendingRetargets = nil
	c.alterExisting = make(map[*Constraint]bool, len(tbl.Constraints))
	for _, con := range tbl.Constraints {
		c.alterExisting[con] = true
	}
	defer func() { c.pendingRetargets, c.alterExisting = nil, nil }()
	if len(stmt.Commands) == 1 && !addNeedsRollback(stmt.Commands[0]) {
		// Single command carrying one table element: no rollback needed.
		if err := c.execAlterCmd(db, tbl, stmt.Commands[0]); err != nil {
			return err
		}
		dropChecksOfMissingColumns(tbl)
		c.runPendingRetargets()
		return nil
	}

	// TiDB resolves every subcommand of a multi-command ALTER against the
	// pre-ALTER schema, so a reference to a column added earlier in the same
	// statement is rejected. Check that before applying any subcommand (omni
	// would otherwise apply them left-to-right and accept the DDL).
	if err := checkMultiCmdGenerated(tbl, stmt.Commands); err != nil {
		return err
	}
	// A key the statement drops is judged against the foreign keys the table
	// began with, one the statement also drops included (v8.5.5: DROP FOREIGN
	// KEY fk_customer, DROP INDEX idx_customer is 1553).
	if c.foreignKeyChecks {
		for _, cmd := range stmt.Commands {
			if cmd == nil || cmd.Type != nodes.ATDropIndex {
				continue
			}
			for _, idx := range tbl.Indexes {
				if toLower(idx.Name) == toLower(cmd.Name) {
					if err := c.checkDropIndexNeededByFK(tbl, idx); err != nil {
						return err
					}
				}
			}
		}
	}
	// ALGORITHM = INSTANT adds no key or constraint: v8.5.5 answers 1846 for
	// ADD INDEX and for ADD CONSTRAINT ... CHECK under it.
	if alterAlgorithmIs(stmt.Commands, "INSTANT") {
		for _, cmd := range stmt.Commands {
			if cmd != nil && (cmd.Type == nodes.ATAddIndex || cmd.Type == nodes.ATAddConstraint) {
				return &Error{Code: 1846, SQLState: "0A000", Message: "ALGORITHM=INSTANT is not supported. Reason: Cannot alter table by INSTANT. Try ALGORITHM=INPLACE."}
			}
		}
	}
	// RENAME TO runs alone (v8.5.5: ALTER TABLE entries RENAME TO
	// simulation_other.entries, ADD COLUMN note INT is 8200).
	for _, cmd := range stmt.Commands {
		if cmd != nil && cmd.Type == nodes.ATRenameTable {
			return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported multi schema change for rename table"}
		}
	}
	if err := checkMultiCmdColumnRefs(tbl, stmt.Commands); err != nil {
		return err
	}

	// Multi-command ALTER: MySQL treats this as atomic.
	// Snapshot the table so we can rollback on any sub-command failure.
	snapshot := cloneTable(tbl)
	origKey := key

	for i := 0; i < len(stmt.Commands); {
		cmd := stmt.Commands[i]
		// Combine contiguous charset/collation ATTableOption commands in the
		// same option run, matching TiDB's GetCharsetAndCollateInTableOption:
		// `CHARSET=gbk COLLATE=gbk_bin` must resolve and validate as a pair
		// rather than validating an intermediate default collation first.
		if isCharsetOrCollateOptionCmd(cmd) {
			end := i + 1
			for end < len(stmt.Commands) && stmt.Commands[end] != nil && stmt.Commands[end].Type == nodes.ATTableOption {
				end++
			}
			if err := c.execAlterTableOptionRun(tbl, stmt.Commands[i:end]); err != nil {
				undoAlter(db, origKey, tbl, snapshot)
				return err
			}
			i = end
			continue
		}
		if len(stmt.Commands) > 1 && c.checkConstraints && isDropOfCheck(tbl, cmd) {
			undoAlter(db, origKey, tbl, snapshot)
			return errMultiSchemaDropCheck()
		}
		// An added column is judged against the table the statement began
		// with: v8.5.5 answers 1060 for DROP COLUMN c, ADD COLUMN c BIGINT,
		// once the commands before it held (a DROP COLUMN a CHECK reads is
		// 3959 first).
		if cmd.Type == nodes.ATAddColumn && !cmd.IfNotExists {
			defs := cmd.Columns
			if cmd.Column != nil {
				defs = append([]*nodes.ColumnDef{cmd.Column}, defs...)
			}
			for _, def := range defs {
				if def != nil && snapshot.GetColumn(def.Name) != nil && tbl.GetColumn(def.Name) == nil {
					undoAlter(db, origKey, tbl, snapshot)
					return errDupColumn(def.Name)
				}
			}
		}
		// So is a column's new name: RENAME COLUMN b TO a, CHANGE b a INT
		// and DROP COLUMN a, RENAME COLUMN b TO a are 1060 when the table
		// began with an a (v8.5.5).
		if to := renamedColumnTo(cmd); to != "" && !strings.EqualFold(to, cmd.Name) && snapshot.GetColumn(to) != nil {
			undoAlter(db, origKey, tbl, snapshot)
			return errDupColumn(to)
		}
		if err := c.execAlterCmd(db, tbl, cmd); err != nil {
			undoAlter(db, origKey, tbl, snapshot)
			return err
		}
		i++
	}
	dropChecksOfMissingColumns(tbl)
	// Clear transient cleanup tracking after successful multi-command ALTER.
	tbl.droppedByCleanup = nil
	c.runPendingRetargets()
	return nil
}

// renamedColumnTo returns the name a RENAME COLUMN or CHANGE gives a column,
// or "" for any other command.
func renamedColumnTo(cmd *nodes.AlterTableCmd) string {
	switch cmd.Type {
	case nodes.ATRenameColumn:
		return cmd.NewName
	case nodes.ATChangeColumn:
		if cmd.Column != nil {
			return cmd.Column.Name
		}
	}
	return ""
}

// runPendingRetargets applies the foreign key retargeting an ALTER TABLE's
// renames queued, in the order they ran.
func (c *Catalog) runPendingRetargets() {
	for _, retarget := range c.pendingRetargets {
		retarget()
	}
	c.pendingRetargets = nil
}

func isCharsetOptionName(name string) bool {
	switch toLower(name) {
	case "charset", "character set", "default charset", "default character set":
		return true
	}
	return false
}

func isCollateOptionName(name string) bool {
	switch toLower(name) {
	case "collate", "default collate":
		return true
	}
	return false
}

func isCharsetOrCollateOptionCmd(cmd *nodes.AlterTableCmd) bool {
	if cmd == nil || cmd.Type != nodes.ATTableOption || cmd.Option == nil {
		return false
	}
	return isCharsetOptionName(cmd.Option.Name) || isCollateOptionName(cmd.Option.Name)
}

// execAlterTableOptionRun applies a contiguous run of ATTableOption commands,
// resolving any charset and collation options in the run together once (just
// as TiDB's GetCharsetAndCollateInTableOption does) while applying all other
// options in order.
func (c *Catalog) execAlterTableOptionRun(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	var csOpt, collOpt string
	var hasCharsetOrCollate bool
	for _, cmd := range cmds {
		if cmd == nil || cmd.Option == nil {
			continue
		}
		switch {
		case isCharsetOptionName(cmd.Option.Name):
			csOpt = cmd.Option.Value
			hasCharsetOrCollate = true
		case isCollateOptionName(cmd.Option.Name):
			collOpt = cmd.Option.Value
			hasCharsetOrCollate = true
		}
	}
	charsetApplied := false
	for _, cmd := range cmds {
		if isCharsetOrCollateOptionCmd(cmd) {
			if !charsetApplied && hasCharsetOrCollate {
				cs, coll, err := resolveCharsetCollation(csOpt, collOpt)
				if err != nil {
					return err
				}
				if err := checkAlterTableCharset(tbl.Charset, tbl.Collation, cs, coll); err != nil {
					return err
				}
				tbl.Charset, tbl.Collation = cs, coll
				charsetApplied = true
			}
			continue
		}
		if err := c.alterTableOption(tbl, cmd); err != nil {
			return err
		}
	}
	return nil
}

// singleCacheCommand reports whether the ALTER TABLE is the one form a cached
// table takes: a single CACHE or NOCACHE.
func singleCacheCommand(cmds []*nodes.AlterTableCmd) bool {
	if len(cmds) != 1 || cmds[0] == nil {
		return false
	}
	return cmds[0].Type == nodes.ATCache || cmds[0].Type == nodes.ATNoCache
}

// checkMultiCmdColumnRefs enforces TiDB v8.5.5's rule that every column a
// subcommand of a multi-command ALTER references — an ADD/MODIFY/CHANGE COLUMN's
// AFTER target, an ADD INDEX/KEY/CONSTRAINT's key columns — must exist in the
// PRE-ALTER table. TiDB resolves all subcommands against the original schema, so
// a reference to a column added earlier in the same statement (ADD y AFTER x,
// ADD INDEX (c) after ADD COLUMN c) is rejected (1054 / 1072) rather than
// resolved left-to-right as MySQL does.
func checkMultiCmdColumnRefs(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	orig := make(map[string]bool, len(tbl.Columns))
	for _, col := range tbl.Columns {
		orig[toLower(col.Name)] = true
	}
	for _, cmd := range cmds {
		switch cmd.Type {
		case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
			if cmd.After != "" && !orig[toLower(cmd.After)] {
				return errNoSuchColumn(cmd.After, tbl.Name)
			}
		case nodes.ATAddIndex, nodes.ATAddConstraint:
			// A FULLTEXT key adds nothing on v8.5.5 and names no column it
			// is held to: ADD COLUMN e VARCHAR(10), ADD FULLTEXT INDEX ft (e)
			// runs.
			if cmd.Constraint != nil && cmd.Constraint.Type != nodes.ConstrFulltextIndex {
				for _, name := range extractColumnNames(cmd.Constraint) {
					if name != "" && !orig[toLower(name)] {
						return errKeyColumnDoesNotExist(name)
					}
				}
			}
		}
	}
	return nil
}

// undoAlter puts tbl back as snapshot holds it, under key in db. A RENAME TO
// earlier in the statement may have moved the table to another database, under
// another name, so its entry there goes first: a rollback that repaired only
// db's own map left the table in two databases at once, and left db without it
// at all when the move kept the name.
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
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		tbl.Partitioning = nil
		return nil
	case nodes.ATAnalyzePartition:
		// ANALYZE PARTITION recomputes statistics and leaves the schema as it
		// is, and a table without partitions is still refused (1505, TiDB
		// v8.5.5), where it takes COMPACT PARTITION, which answers 1105.
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		for _, idxName := range cmd.AnalyzeIndexes {
			if findIndexByName(tbl, idxName) == nil {
				return errAnalyzeIndexNotFound(idxName, tbl.Name)
			}
		}
		return nil
	case nodes.ATPartitionAttributes:
		// A partition's ATTRIBUTES is a server-side setting no field of Table
		// or PartitionDefInfo holds; a table without partitions is still
		// refused (1505, TiDB v8.5.5), and the attribute format is validated (8237).
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		return checkAlterAttributes(cmd.Attributes)
	case nodes.ATAddLastPartition, nodes.ATDropFirstPartition:
		// The interval partitions FIRST and LAST name are server-side settings
		// no field of Table or PartitionDefInfo holds; a table without
		// partitions is still refused (1505, TiDB v8.5.5), and a non-constant
		// bound referencing a column is refused with 1054.
		if tbl.Partitioning == nil {
			return errPartitionMgmtOnNonPartitioned()
		}
		if name := firstColumnRef(cmd.PartitionBound); name != "" {
			return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
				Message: fmt.Sprintf("Unknown column '%s' in 'expression'", name)}
		}
		return nil
	case nodes.ATPartitionBy:
		return validateRepartition(tbl, cmd.PartitionBy)
	case nodes.ATSetTiFlashReplica:
		// SET HYPO TIFLASH REPLICA pins a hypothetical replica for the
		// optimizer and leaves the physical replica untouched, so the catalog
		// keeps its real TiFlashReplica and labels: overwriting them would
		// record one real replica and no labels for a table that has neither
		// changed.
		if cmd.Hypo {
			return nil
		}
		if c.tiflashServerCount >= 0 && cmd.TiFlashReplica > c.tiflashServerCount {
			return &Error{
				Code:     ErrUnknownError,
				SQLState: sqlState(ErrUnknownError),
				Message:  fmt.Sprintf("the tiflash replica count: %d should be less than the total tiflash server count: %d", cmd.TiFlashReplica, c.tiflashServerCount),
			}
		}
		tbl.TiFlashReplica = cmd.TiFlashReplica
		// LOCATION LABELS is optional; assign (possibly empty slice)
		// so SET TIFLASH REPLICA 0 with no clause clears prior labels.
		tbl.TiFlashLocationLabels = append([]string(nil), cmd.TiFlashLocationLabels...)
		return nil
	case nodes.ATRemoveTTL:
		tbl.TTLColumn = ""
		tbl.TTLInterval = ""
		tbl.TTLEnable = false
		tbl.TTLJobInterval = ""
		return nil

	case nodes.ATCache:
		// CACHE changes no column or key, but SHOW CREATE TABLE writes the
		// attribute back as `/* CACHED ON */`, so it is part of the schema the
		// catalog stores (read on v8.5.5). A partitioned table cannot be cached
		// at all: v8.5.5 answers 8242 "'partition mode' is unsupported on cache
		// tables".
		if tbl.Partitioning != nil {
			return errCachedTable("partition mode")
		}
		tbl.Cached = true
		return nil
	case nodes.ATNoCache:
		tbl.Cached = false
		return nil

	case nodes.ATWriteable:
		return nil
	case nodes.ATCompact:
		if len(cmd.PartitionNames) > 0 && tbl.Partitioning == nil {
			return errCompactNonPartitionedTable(tbl.Name, cmd.PartitionNames)
		}
		return nil
	case nodes.ATAttributes:
		return checkAlterAttributes(cmd.Attributes)
	case nodes.ATPartitionOptions:
		for _, opt := range cmd.PartitionOptions {
			if opt == nil {
				continue
			}
			if toLower(opt.Name) != "placement policy" {
				return &Error{
					Code:     ErrUnknownError,
					SQLState: sqlState(ErrUnknownError),
					Message:  "unknown partition option",
				}
			}
		}
		return nil
	case nodes.ATAddStatistics, nodes.ATDropStatistics:
		if !c.extendedStats {
			return &Error{
				Code:     ErrUnknownError,
				SQLState: sqlState(ErrUnknownError),
				Message:  "Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF",
			}
		}
		return nil
	case nodes.ATCheckPartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported check partition"}
	case nodes.ATOptimizePartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported optimize partition"}
	case nodes.ATRebuildPartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported rebuild partition"}
	case nodes.ATRepairPartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported repair partition"}
	case nodes.ATReorganizeLastPartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported SPLIT LAST PARTITION"}
	case nodes.ATReorganizeFirstPartition:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported MERGE FIRST PARTITION"}
	case nodes.ATDiscardPartitionTablespace, nodes.ATImportPartitionTablespace,
		nodes.ATDiscardTablespace, nodes.ATImportTablespace,
		nodes.ATForce, nodes.ATSecondaryLoad, nodes.ATSecondaryUnload,
		nodes.ATStatsOptions:
		return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported Unsupported/unknown ALTER TABLE specification"}

	default:
		// Unsupported alter command; silently ignore.
		return nil
	}
}

// alterAddColumn adds a new column to the table.
func (c *Catalog) alterAddColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	// The parenthesized form reads a table element list, so it carries the
	// constraints written beside the columns:
	// `ADD COLUMN (a1 INT, PRIMARY KEY (a1))` adds both. The columns come
	// first, as the CREATE TABLE path builds them.
	//
	// TiDB resolves a constraint's columns against the table as it stood
	// before the statement, so `ADD COLUMN (c INT, PRIMARY KEY (c))` is an
	// error (1072, "column does not exist: c") on v8.5.5 where MySQL 8.0.46
	// makes c the primary key. This catalog does not model that ordering: it
	// adds the columns first, so the key resolves against them and the
	// statement is accepted. A key part naming a column that this statement
	// does not add either is rejected (1072), as everywhere else.
	if len(cmd.Columns) > 0 || len(cmd.Constraints) > 0 {
		for _, colDef := range cmd.Columns {
			if err := c.addSingleColumnIfNotExists(tbl, colDef, cmd.IfNotExists, false, ""); err != nil {
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

	return c.addSingleColumnIfNotExists(tbl, colDef, cmd.IfNotExists, cmd.First, cmd.After)
}

// addSingleColumnIfNotExists adds one column, and does nothing rather than
// failing when the column is already there and the action was written with
// IF NOT EXISTS. TiDB v8.5.5 raises warning 1060 and leaves the table alone;
// a catalog has no warnings, so the table is what it records.
func (c *Catalog) addSingleColumnIfNotExists(tbl *Table, colDef *nodes.ColumnDef, ifNotExists, first bool, after string) error {
	if err := checkUnsupportedAddColumnConstraint(tbl, colDef); err != nil {
		return err
	}
	if ifNotExists {
		if _, exists := tbl.colByName[toLower(colDef.Name)]; exists {
			return nil
		}
	}
	return c.addSingleColumn(tbl, colDef, first, after)
}

// addSingleColumn adds one column definition to the table.
func (c *Catalog) addSingleColumn(tbl *Table, colDef *nodes.ColumnDef, first bool, after string) error {
	// TiDB v8.5.5 refuses inline AUTO_RANDOM (8216), AUTO_INCREMENT, PRIMARY
	// KEY, and UNIQUE KEY constraints on ADD COLUMN (8200) ahead of
	// duplicate-column (1060) and multiple-PK (1068) checks.
	if err := checkUnsupportedAddColumnConstraint(tbl, colDef); err != nil {
		return err
	}
	// A DEFAULT that reads a sequence resolves it first: a name that is a table
	// but not a sequence is 1347, and nothing of the name is 1146, both ahead
	// of the duplicate-column check below (read on v8.5.5).
	if err := checkColumnSequenceDefault(tbl, colDef); err != nil {
		return err
	}
	colKey := toLower(colDef.Name)
	if _, exists := tbl.colByName[colKey]; exists {
		return errDupColumn(colDef.Name)
	}
	// TiDB adds no STORED generated column, and the duplicate-column check
	// comes first (both read on v8.5.5).
	if colDef.Generated != nil && colDef.Generated.Stored {
		return errAlterStoredGeneratedColumn()
	}

	col, err := buildColumnFromDef(tbl, colDef)
	if err != nil {
		return err
	}
	if err := c.checkColumnLength(colDef, col.Charset); err != nil {
		return err
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
				// TiDB: hoist CLUSTERED/NONCLUSTERED from the inline PK
				// constraint of an ALTER TABLE ADD COLUMN ... PRIMARY KEY
				// CLUSTERED into the synthesized table-level constraint.
				// Mirrors the CREATE TABLE path in tablecmds.go.
				Clustered: cc.Clustered,
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

	// TiDB v8.5.5 refuses dropping a table's only remaining column (1090). A
	// multi-subcommand ALTER mutates the table between drops, so the last one
	// is caught here.
	if len(tbl.Columns) == 1 {
		return errCantRemoveAllFields(cmd.Name, tbl.Name)
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

	// TiDB v8.5.5 (8200) rejects dropping a column that the primary key covers,
	// or that a composite (multi-column) index covers. A column covered only by
	// single-column secondary indexes drops cleanly — those indexes drop with
	// it (see cleanupIndexesForDroppedColumn). Measured against the container.
	// A column another table's foreign key references is 1829, ahead of the
	// keys that cover it (v8.5.5: DROP COLUMN id over customers, the primary
	// key orders' fk_customer references, is 1829 and not 8200).
	if c.foreignKeyChecks {
		for _, ref := range c.referencingForeignKeys(tbl) {
			if ref.child == tbl {
				continue
			}
			for _, rc := range ref.con.RefColumns {
				if toLower(rc) == colKey {
					return &Error{Code: 1829, SQLState: "HY000", Message: fmt.Sprintf(
						"Cannot drop column '%s': needed in a foreign key constraint '%s' of table '%s'", cmd.Name, ref.con.Name, ref.child.Name)}
				}
			}
		}
	}
	if err := checkDropColumnCovered(tbl, cmd.Name); err != nil {
		return err
	}
	if err := checkDropColumnWithCheck(tbl, cmd.Name); err != nil {
		return err
	}

	// Remove column from indexes; if index becomes empty, remove it entirely.
	cleanupIndexesForDroppedColumn(tbl, cmd.Name)

	idx := tbl.colByName[colKey]
	tbl.Columns = append(tbl.Columns[:idx], tbl.Columns[idx+1:]...)
	rebuildColIndex(tbl)
	return nil
}

// checkDropColumnCovered enforces TiDB v8.5.5's 8200 rejections for DROP COLUMN.
// A column that is part of the primary key, or of any composite (multi-column)
// index, cannot be dropped — not even by dropping the covering index in the
// same statement (measured). The single integer, clustered primary key (the row
// handle, pkIsHandle) carries its own message. A column covered only by
// single-column secondary indexes is droppable and returns nil here.
func checkDropColumnCovered(tbl *Table, colName string) error {
	colKey := toLower(colName)
	for _, idx := range tbl.Indexes {
		covered := false
		for _, ic := range idx.Columns {
			if toLower(ic.Name) == colKey {
				covered = true
				break
			}
		}
		if !covered {
			continue
		}
		if idx.Primary {
			if isRowHandle(tbl, idx) && primaryKeyClustered(tbl, idx) {
				return errDropIntegerPrimaryKey()
			}
			return errDropColumnCovered(colName)
		}
		if len(idx.Columns) > 1 {
			return errDropColumnCovered(colName)
		}
	}
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
	if err := checkColumnQualifier(cmd.Column.Schema, cmd.Column.Table, tbl.Database.Name, tbl.Name); err != nil {
		return err
	}
	if skipMissingColumn(tbl, cmd.Column.Name, cmd.IfExists) {
		return nil
	}
	return c.alterReplaceColumn(tbl, cmd.Column.Name, cmd)
}

// skipMissingColumn reports whether an action written with IF EXISTS names a
// column the table does not have, which TiDB answers with a note (1054,
// verified on v8.5.5) rather than an error.
func skipMissingColumn(tbl *Table, name string, ifExists bool) bool {
	if !ifExists {
		return false
	}
	_, exists := tbl.colByName[toLower(name)]
	return !exists
}

// alterChangeColumn replaces a column (old name -> new name + new definition).
func (c *Catalog) alterChangeColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if cmd.Column == nil {
		return nil
	}
	if err := checkColumnQualifier(cmd.OldColumnSchema, cmd.OldColumnTable, tbl.Database.Name, tbl.Name); err != nil {
		return err
	}
	if err := checkColumnQualifier(cmd.Column.Schema, cmd.Column.Table, tbl.Database.Name, tbl.Name); err != nil {
		return err
	}
	if skipMissingColumn(tbl, cmd.Name, cmd.IfExists) {
		return nil
	}
	return c.alterReplaceColumn(tbl, cmd.Name, cmd)
}

// alterReplaceColumn is the shared implementation for MODIFY and CHANGE COLUMN.
// oldName is the existing column to replace; cmd.Column defines the new column.
func (c *Catalog) alterReplaceColumn(tbl *Table, oldName string, cmd *nodes.AlterTableCmd) error {
	colDef := cmd.Column
	// v8.5.5 redefines no column with a CHECK, whether CHECK constraints are
	// on or not: MODIFY age INT CHECK (age > 1) is 8200.
	for _, cc := range colDef.Constraints {
		if cc != nil && cc.Type == nodes.ColConstrCheck {
			return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported modify column: can't modify with check"}
		}
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
	// A VIRTUAL column made a plain one is 3106 'Changing the STORED status',
	// and one a key reads redefined as generated again 3106 'modifying an
	// indexed column' (v8.5.5, features@generated: MODIFY full_name
	// VARCHAR(101) and MODIFY full_name ... AS (...) VIRTUAL under idx_full).
	if oldCol.Generated != nil && !oldCol.Generated.Stored {
		switch {
		case colDef.Generated == nil:
			return errUnsupportedGeneratedStorageChange(colDef.Name, tbl.Name)
		case columnInAnyIndex(tbl, oldCol.Name):
			return &Error{Code: ErrUnsupportedGeneratedStorageChange, SQLState: sqlState(ErrUnsupportedGeneratedStorageChange),
				Message: "'modifying an indexed column' is not supported for generated columns."}
		}
	}
	// A column a generated column reads keeps its name (v8.5.5: CHANGE
	// first_name given_name ... under full_name is 3108).
	if newKey != oldKey {
		for _, gc := range tbl.Columns {
			if gc.Generated != nil && gc != oldCol && generatedExprReferencesColumn(gc.Generated.Expr, oldCol.Name) {
				return errDependentByGeneratedColumn(oldCol.Name, gc.Name, tbl.Name)
			}
		}
	}

	col, err := buildColumnFromDef(tbl, colDef)
	if err != nil {
		return err
	}
	if err := c.checkColumnLength(colDef, col.Charset); err != nil {
		return err
	}
	// A string column keeps to the character set changes TiDB makes (v8.5.5:
	// MODIFY tag VARCHAR(20) over a utf8mb4 tag in a latin1 table is 8200
	// "Unsupported modify charset from utf8mb4 to latin1").
	if isStringType(oldCol.DataType) && isStringType(col.DataType) {
		if err := checkColumnCharsetChange(columnCharsetOf(tbl, oldCol), columnCharsetOf(tbl, col)); err != nil {
			return err
		}
	}
	if c.foreignKeyChecks {
		if err := c.checkModifyColumnFKCompatibility(tbl, oldName, col); err != nil {
			return err
		}
	}
	if err := checkModifyPrimaryKeyColumnType(tbl, oldCol, col); err != nil {
		return err
	}
	col.Position = idx + 1
	tbl.Columns[idx] = col

	// Update index/constraint column references if name changed, and the
	// foreign keys that reference the column (fk_retarget.go).
	if newKey != oldKey {
		updateColumnRefsInIndexes(tbl, oldName, colDef.Name)
		db, name, to, existing := tbl.Database, tbl.Name, colDef.Name, c.alterExisting
		c.pendingRetargets = append(c.pendingRetargets, func() {
			c.retargetReferences(tbl, db, name, map[string]string{oldKey: to}, existing)
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
func (c *Catalog) alterAddConstraint(tbl *Table, cmd *nodes.AlterTableCmd) error {
	con := cmd.Constraint
	if con == nil {
		return nil
	}

	// IF NOT EXISTS covers an index the table already has: TiDB v8.5.5 raises
	// note 1061 and leaves the table alone rather than failing, so
	// `ADD INDEX IF NOT EXISTS i (a)` run twice is one index.
	//
	// A foreign key is not in that namespace and is not covered either way:
	// an index named fk does not stand in the way of a foreign key of the
	// name -- v8.5.5 runs `ADD CONSTRAINT fk FOREIGN KEY IF NOT EXISTS (a)
	// REFERENCES p (a)` against a table that has `KEY fk (a)` and keeps both
	// -- and a foreign key of the name is 1826 even with IF NOT EXISTS
	// written, which is an error and not a note (read on 2026-11-04). So the
	// guard skips foreign keys and leaves them to the duplicate check below.
	if con.IfNotExists && con.Type != nodes.ConstrForeignKey &&
		con.Name != "" && indexNameExists(tbl, con.Name) {
		return nil
	}

	cols := extractColumnNames(con)

	switch con.Type {
	case nodes.ConstrPrimaryKey:
		for _, ic := range con.IndexColumns {
			if ic != nil && ic.Functional {
				return &Error{Code: 3756, SQLState: "HY000", Message: "The primary key cannot be an expression index"}
			}
		}
		// Check for duplicate PK.
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				return errMultiplePriKey()
			}
		}
		idxCols := buildIndexColumns(con)
		// v8.5.5 words 1072 for a key added by ALTER as its secondary-index
		// path does, even for a PRIMARY KEY.
		if err := validateIndexColumnsExist(tbl, idxCols); err != nil {
			return err
		}
		if err := checkNonVectorHNSW(resolveConstraintIndexType(con)); err != nil {
			return err
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
		// The CREATE TABLE path applies the options of every index kind; this one
		// used to drop them, so an index added by ALTER lost its comment,
		// visibility and block size.
		applyIndexOptions(pkIdx, con.IndexOptions)
		// TiDB v8.5.5 refuses the statement rather than add an invisible key, and
		// refuses it whole: the columns stay nullable.
		if !pkIdx.Visible {
			return errPrimaryKeyInvisible()
		}
		pkIdx.addedByAlter = true
		// TiDB refuses an added CLUSTERED primary key outright (8200), before it
		// would look at the partitioning.
		if con.Clustered == nil || !*con.Clustered {
			if err := validateNewUniqueKey(tbl, pkIdx); err != nil {
				return err
			}
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
			// TiDB: propagate CLUSTERED/NONCLUSTERED from ALTER TABLE
			// ADD PRIMARY KEY (...) CLUSTERED / NONCLUSTERED. Mirrors
			// the CREATE TABLE path in tablecmds.go.
			Clustered: con.Clustered,
		})

	case nodes.ConstrUnique:
		if err := checkIndexName(con.Name); err != nil {
			return err
		}
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" && indexNameExists(tbl, idxName) {
			return errDupKeyName(idxName)
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumnsExist(tbl, idxCols); err != nil {
			return err
		}
		if err := checkNonVectorHNSW(resolveConstraintIndexType(con)); err != nil {
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
		if err := validateNewUniqueKey(tbl, uqIdx); err != nil {
			return err
		}
		tbl.Indexes = append(tbl.Indexes, uqIdx)
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:      idxName,
			Type:      ConUniqueKey,
			Table:     tbl,
			Columns:   cols,
			IndexName: idxName,
		})

	case nodes.ConstrForeignKey:
		// A referenced key part that is an expression names no column, and
		// TiDB's DDL fails on it (v8.5.5 raises 1105), so refuse the key rather
		// than record one with an empty referenced-column list. This mirrors the
		// CREATE TABLE path.
		for _, kp := range con.RefKeyParts {
			if kp.Functional {
				return errCannotAddForeign()
			}
		}
		// The key takes the next number only once it is added: TiDB keeps
		// MaxForeignKeyID in the table's meta, which a refused statement
		// leaves alone.
		fkID := tbl.maxForeignKeyID + 1
		conName := fkName(con)
		if conName == "" {
			conName = fmt.Sprintf("fk_%d", fkID)
		}
		// A table-level FOREIGN KEY names the referenced columns it matches, and
		// TiDB refuses one that names none (1239), whether or not the
		// referenced table exists: `ADD FOREIGN KEY (c) REFERENCES c` is 1239
		// even with no table c (read on v8.5.5). The column-level REFERENCES a
		// column may carry takes the referenced table's key instead, and is not
		// this rule's.
		if len(con.RefColumns) == 0 {
			return errFKKeyRefMismatch(conName)
		}
		// A partitioned table takes no foreign key, which v8.5.5 judges after
		// the key's own definition (1239 first): ADD CONSTRAINT fk_archive
		// FOREIGN KEY ... on a RANGE table is 1506.
		if tbl.Partitioning != nil {
			return errForeignKeyOnPartitioned()
		}
		// A name another foreign key of the table holds is 1826 (v8.5.5).
		if findForeignKey(tbl, conName) != nil {
			return &Error{Code: 1826, SQLState: "HY000", Message: fmt.Sprintf("Duplicate foreign key constraint name '%s'", conName)}
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
		// The backing index TiDB creates first takes the foreign key's name.
		if !hasIndexCoveringColumns(tbl, cols) && indexNameExists(tbl, conName) {
			return errDupKeyName(conName)
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
		if err := ensureFKBackingIndex(tbl, conName, cols, buildIndexColumns(con)); err != nil {
			return err
		}
		tbl.maxForeignKeyID = fkID

	case nodes.ConstrCheck:
		if !c.checkConstraints {
			break
		}
		conName := con.Name
		if conName == "" {
			conName = fmt.Sprintf("%s_chk_%d", tbl.Name, nextCheckNumber(tbl))
		}
		if err := validateCheckConstraint(tbl, conName, con.Expr); err != nil {
			return err
		}
		tbl.Constraints = append(tbl.Constraints, &Constraint{
			Name:        conName,
			Type:        ConCheck,
			Table:       tbl,
			CheckExpr:   checkExprSQL(con.Expr),
			NotEnforced: con.NotEnforced,
		})

	case nodes.ConstrIndex:
		if err := checkIndexName(con.Name); err != nil {
			return err
		}
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" && indexNameExists(tbl, idxName) {
			return errDupKeyName(idxName)
		}
		idxCols := buildIndexColumns(con)
		if err := validateIndexColumnsExist(tbl, idxCols); err != nil {
			return err
		}
		// An expression key part is held to what CREATE TABLE holds it to
		// (validateFunctionalIndexColumns), and names columns the table has
		// (v8.5.5: ADD INDEX idx_bad ((nope + 1)) is 1054 'expression').
		if err := validateFunctionalIndexColumns(con.IndexColumns); err != nil {
			return err
		}
		if err := checkFunctionalKeyColumns(tbl, con.IndexColumns); err != nil {
			return err
		}
		if err := checkNonVectorHNSW(resolveConstraintIndexType(con)); err != nil {
			return err
		}
		keyIdx := &Index{
			Name:      idxName,
			Table:     tbl,
			Columns:   idxCols,
			IndexType: resolveConstraintIndexType(con),
			Visible:   true,
		}
		applyIndexOptions(keyIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, keyIdx)

	case nodes.ConstrFulltextIndex:
		// ALTER TABLE ADD FULLTEXT succeeds on v8.5.5 and adds nothing, as the
		// key does in CREATE TABLE (read on 2026-11-06).

	case nodes.ConstrSpatialIndex:
		idxName := con.Name
		if idxName == "" && len(cols) > 0 {
			idxName = allocIndexName(tbl, cols[0])
		} else if idxName != "" && indexNameExists(tbl, idxName) {
			return errDupKeyName(idxName)
		}
		idxCols := buildIndexColumns(con)
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

	case nodes.ConstrVector:
		// TiDB's VECTOR INDEX, which is an index like the rest: it takes a
		// name, key parts and index options, and `ADD VECTOR INDEX` that
		// applied nothing would read as a successful no-op for a statement
		// that changes the schema. The index type is the USING name (HNSW),
		// which is where a vector index carries its algorithm.
		if con.Name != "" && indexNameExists(tbl, con.Name) {
			return errDupKeyName(con.Name)
		}
		vecIdx := &Index{
			Name:      vectorIndexName(tbl, con, cols),
			Table:     tbl,
			Columns:   buildIndexColumns(con),
			Vector:    true,
			IndexType: vectorIndexHNSW,
			Visible:   true,
		}
		applyIndexOptions(vecIdx, con.IndexOptions)
		tbl.Indexes = append(tbl.Indexes, vecIdx)
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
			// DROP INDEX `PRIMARY` hits the same clustered-PK rejection as
			// DROP PRIMARY KEY on TiDB v8.5.5 (8200). Measured against v8.5.5.
			if idx.Primary && primaryKeyClustered(tbl, idx) {
				return errDropClusteredPrimaryKey()
			}
			if err := c.checkDropIndexNeededByFK(tbl, idx); err != nil {
				return err
			}
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
		return errIndexDoesntExist(name)
	}

	// Also remove any unique/PK constraint that references this index.
	for i, con := range tbl.Constraints {
		if con.Type == ConForeignKey || con.Type == ConCheck {
			continue
		}
		if toLower(con.IndexName) == key || toLower(con.Name) == key {
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

	switch {
	case cmd.PrimaryKey:
		// DROP PRIMARY KEY: TiDB v8.5.5 (8200) rejects dropping the primary key
		// when the table clusters by it. A NONCLUSTERED primary key, or one
		// added by ALTER TABLE ADD PRIMARY KEY, drops cleanly.
		for _, idx := range tbl.Indexes {
			if idx.Primary {
				if primaryKeyClustered(tbl, idx) {
					return errDropClusteredPrimaryKey()
				}
				if err := c.checkDropIndexNeededByFK(tbl, idx); err != nil {
					return err
				}
			}
		}
		found := false
		for i, con := range tbl.Constraints {
			if con.Type == ConPrimaryKey {
				tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
				found = true
				break
			}
		}
		for i, idx := range tbl.Indexes {
			if idx.Primary {
				tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
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
		return nil

	case cmd.ForeignKey:
		// For FK constraints, TiDB keeps the backing index when dropping the FK.
		for i, con := range tbl.Constraints {
			if con.Type == ConForeignKey && toLower(con.Name) == key {
				tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
				return nil
			}
		}
		if cmd.IfExists {
			return nil
		}
		return errCantDropKey(name)

	default:
		// DROP CHECK and DROP CONSTRAINT only target CHECK constraints in TiDB
		// (never foreign keys, primary keys, or unique keys). With
		// tidb_enable_check_constraint off, TiDB treats both as a no-op with a
		// warning (v8.5.5).
		if !c.checkConstraints {
			return nil
		}
		for i, con := range tbl.Constraints {
			if con.Type == ConCheck && toLower(con.Name) == key {
				tbl.Constraints = append(tbl.Constraints[:i], tbl.Constraints[i+1:]...)
				return nil
			}
		}
		if cmd.IfExists {
			return nil
		}
		// v8.5.5 answers 3940 for a CHECK the table does not have (DROP CHECK
		// nope, and DROP CONSTRAINT naming a foreign key).
		return &Error{Code: 3940, SQLState: "HY000", Message: fmt.Sprintf("Constraint '%s' does not exist.", name)}
	}
}

// alterRenameColumn changes a column name in-place.
func (c *Catalog) alterRenameColumn(tbl *Table, cmd *nodes.AlterTableCmd) error {
	oldKey := toLower(cmd.Name)
	idx, exists := tbl.colByName[oldKey]
	if !exists {
		return errNoSuchColumn(cmd.Name, tbl.Name)
	}
	if err := checkRenameColumnWithCheck(tbl, tbl.Columns[idx].Name); err != nil {
		return err
	}

	newKey := toLower(cmd.NewName)
	if newKey != oldKey {
		if _, dup := tbl.colByName[newKey]; dup {
			return errDupColumn(cmd.NewName)
		}
		// A column the partitioning reads keeps its name (v8.5.5: RENAME
		// COLUMN created TO created_at over RANGE (YEAR(created)) is 3855).
		if partitionReadsColumn(tbl, tbl.Columns[idx].Name) {
			return errPartitionColumnDependency(tbl.Columns[idx].Name)
		}
	}

	oldName := tbl.Columns[idx].Name
	tbl.Columns[idx].Name = cmd.NewName
	updateColumnRefsInIndexes(tbl, cmd.Name, cmd.NewName)
	rebuildColIndex(tbl)
	// The foreign keys that reference the column follow it (v8.5.5: RENAME
	// COLUMN id TO cid over customers leaves orders' fk_customer referencing
	// `cid`).
	if newKey != oldKey {
		db, name, from, to, existing := tbl.Database, tbl.Name, toLower(oldName), cmd.NewName, c.alterExisting
		c.pendingRetargets = append(c.pendingRetargets, func() {
			c.retargetReferences(tbl, db, name, map[string]string{from: to}, existing)
		})
	}
	return nil
}

// alterRenameIndex changes an index name in-place.
func (c *Catalog) alterRenameIndex(tbl *Table, cmd *nodes.AlterTableCmd) error {
	oldKey := toLower(cmd.Name)
	newKey := toLower(cmd.NewName)

	if newKey != oldKey && indexNameExists(tbl, cmd.NewName) {
		return errDupKeyName(cmd.NewName)
	}

	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == oldKey {
			idx.Name = cmd.NewName
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

	return errKeyDoesNotExist(cmd.Name, tbl.Name)
}

// alterRenameTable moves a table to a new name.
func (c *Catalog) alterRenameTable(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	// The new name may carry a database qualifier, and then the rename moves
	// the table to that database, exactly as RENAME TABLE does: TiDB v8.5.5,
	// MySQL 8.0.46 and MariaDB 11.8.8 all run `ALTER TABLE t RENAME TO d.t1`.
	// NewName holds the name as written, so reading it alone would have made
	// a table called "d.t1" in the database the statement started in.
	//
	// A destination with no qualifier lands in the current database, not in
	// the renamed table's: `USE e; ALTER TABLE d.t RENAME TO u` leaves u in e
	// on TiDB v8.5.5 and on MySQL 8.0.46 (read on 2026-11-04).
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
	// entry behind, so `RENAME TO e.u, RENAME TO e.v` would have left both
	// e.u and e.v pointing at the table. TiDB v8.5.5 answers 8200
	// ("Unsupported multi schema change for rename table") to that statement,
	// but MySQL 8.0.46 and MariaDB 11.8.8 both run it and end with e.v alone.
	oldDB := db
	if tbl.Database != nil {
		oldDB = tbl.Database
	}
	newKey := toLower(newName)
	oldKey := toLower(tbl.Name)

	if (newDB != oldDB || newKey != oldKey) && (newDB.Tables[newKey] != nil || newDB.Views[newKey] != nil || newDB.Sequences[newKey] != nil) {
		return errDupTable(newName)
	}

	oldName := tbl.Name
	delete(oldDB.Tables, oldKey)
	tbl.Name = newName
	tbl.Database = newDB
	newDB.Tables[newKey] = tbl
	// The foreign keys that reference the table follow it (fk_retarget.go).
	c.pendingRetargets = append(c.pendingRetargets, func() {
		c.retargetReferences(tbl, oldDB, oldName, nil, nil)
	})
	return nil
}

// alterTableOption applies a table option (ENGINE, CHARSET, etc.).
func (c *Catalog) alterTableOption(tbl *Table, cmd *nodes.AlterTableCmd) error {
	opt := cmd.Option
	if opt == nil {
		return nil
	}

	switch toLower(opt.Name) {
	case "engine", "storage engine":
		tbl.Engine = opt.Value
	case "charset", "character set", "default charset", "default character set":
		cs, coll, err := resolveCharsetCollation(opt.Value, "")
		if err != nil {
			return err
		}
		if err := checkAlterTableCharset(tbl.Charset, tbl.Collation, cs, coll); err != nil {
			return err
		}
		tbl.Charset, tbl.Collation = cs, coll
	case "collate", "default collate":
		cs, coll, err := resolveCharsetCollation("", opt.Value)
		if err != nil {
			return err
		}
		if err := checkAlterTableCharset(tbl.Charset, tbl.Collation, cs, coll); err != nil {
			return err
		}
		tbl.Charset, tbl.Collation = cs, coll
	case "comment":
		tbl.Comment = opt.Value
	case "auto_increment":
		fmt.Sscanf(opt.Value, "%d", &tbl.AutoIncrement)
	case "row_format":
		tbl.RowFormat = opt.Value
	// TiDB-specific ALTER TABLE options.
	case "shard_row_id_bits":
		var bits int
		fmt.Sscanf(opt.Value, "%d", &bits)
		if bits > 0 && !tableHasTiDBRowID(tbl) {
			return &Error{
				Code:     ErrUnsupportedDDL,
				SQLState: sqlState(ErrUnsupportedDDL),
				Message:  "Unsupported shard_row_id_bits for table with primary key as row id",
			}
		}
		tbl.ShardRowIDBits = bits
	case "auto_id_cache":
		var newCache int64
		fmt.Sscanf(opt.Value, "%d", &newCache)
		if (newCache == 1) != (tbl.AutoIDCache == 1) {
			return &Error{
				Code:     ErrUnknownError,
				SQLState: sqlState(ErrUnknownError),
				Message:  "Can't Alter AUTO_ID_CACHE between 1 and non-1, the underlying implementation is different",
			}
		}
		tbl.AutoIDCache = newCache
	case "auto_random_base":
		// The base is a setting on an auto-random column's table, so a table
		// whose columns carry none is refused (8216, read on v8.5.5).
		if !hasAutoRandom(tbl) {
			return errAutoRandomBaseNoColumn()
		}
		fmt.Sscanf(opt.Value, "%d", &tbl.AutoRandomBase)
	case "placement policy":
		if err := c.validatePolicyRef(opt.Value); err != nil {
			return err
		}
		tbl.PlacementPolicy = resolvePolicyRef(opt.Value)
	case "ttl":
		col, interval, err := extractTTLParts(opt.Value)
		if err != nil {
			return err
		}
		// TiDB resolves the TTL column before it records the config.
		if err := checkTTLColumn(tbl, col); err != nil {
			return err
		}
		if tbl.TTLColumn == "" {
			// A table that gets its first TTL from ALTER TABLE has it
			// enabled, as one created with it does (v8.5.5 writes
			// TTL_ENABLE='ON', read on 2026-11-06).
			tbl.TTLEnable = true
		}
		tbl.TTLColumn = col
		tbl.TTLInterval = interval
	case "ttl_enable":
		tbl.TTLEnable = strings.EqualFold(opt.Value, "ON")
	case "ttl_job_interval":
		tbl.TTLJobInterval = opt.Value
	case "affinity":
		if err := checkAlterAffinity(opt.Value); err != nil {
			return err
		}
	default:
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  "This type of ALTER TABLE is currently unsupported",
		}
	}
	return nil
}

// checkTTLColumn resolves the column a TTL config names: a name the table does
// not have is 1054 ("in 'TTL config'") and a column whose type is not a
// DATETIME, DATE or TIMESTAMP one is 8148 (both read on v8.5.5).
func checkTTLColumn(tbl *Table, col string) error {
	ttlCol := tbl.GetColumn(col)
	if ttlCol == nil {
		return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
			Message: fmt.Sprintf("Unknown column '%s' in 'TTL config'", col)}
	}
	switch toLower(ttlCol.DataType) {
	case "datetime", "date", "timestamp":
		return nil
	}
	return errTTLColumnType(col)
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
		s := defaultToSQL(cmd.DefaultExpr)
		probe := *col
		probe.Default = &s
		if err := checkLiteralDefaultAllowed(&probe); err != nil {
			return err
		}
		if err := checkSetDefaultColumnRef(tbl, col, cmd.DefaultExpr); err != nil {
			return err
		}
		if err := checkExistingColumnDefault(tbl, col, cmd.DefaultExpr); err != nil {
			return err
		}
		if err := checkExistingColumnNowDefault(col, cmd.DefaultExpr); err != nil {
			return err
		}
		col.Default = &s
		col.DefaultDropped = false
	} else {
		// DROP DEFAULT — MySQL shows no default at all (not even DEFAULT NULL).
		col.Default = nil
		col.DefaultDropped = true
	}
	return nil
}

// alterIndexVisibility toggles the Visible flag on an index.
func (c *Catalog) alterIndexVisibility(tbl *Table, cmd *nodes.AlterTableCmd, visible bool) error {
	key := toLower(cmd.Name)
	for _, idx := range tbl.Indexes {
		// The primary key is no key ALTER INDEX names (v8.5.5: ALTER INDEX
		// `PRIMARY` INVISIBLE is 1176).
		if toLower(idx.Name) == key && !idx.Primary {
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
			// TiDB v8.5.5 resolves an AFTER position through its column-dependency
			// pass and names the context "generated column function" for an
			// unknown target column (measured on a live container, 2026-11-09).
			return errNoSuchColumn(after, "generated column function")
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
func buildColumnFromDef(tbl *Table, colDef *nodes.ColumnDef) (*Column, error) {
	col := &Column{
		Name:     colDef.Name,
		Nullable: true,
	}

	// Type info.
	if colDef.TypeName != nil {
		col.DataType = toLower(colDef.TypeName.Name)
		// MySQL 8.0 normalizes GEOMETRYCOLLECTION → geomcollection.
		if col.DataType == "geometrycollection" {
			col.DataType = "geomcollection"
		}
		col.ColumnType = formatColumnType(colDef.TypeName)
		cs, coll, err := resolveCharsetCollation(colDef.TypeName.Charset, colDef.TypeName.Collate)
		if err != nil {
			return nil, err
		}
		if cs != "" {
			col.Charset, col.Collation = cs, coll
		}

		// MySQL converts non-ENUM/SET string columns with CHARACTER SET binary to
		// binary types (CHAR->BINARY, VARCHAR->VARBINARY, TEXT->BLOB), like the
		// CREATE TABLE path. ENUM/SET keep the binary charset annotation.
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
				if dc, ok := defaultCollationForCharset[toLower(col.Charset)]; ok {
					col.Collation = dc
				}
			} else {
				col.Collation = tbl.Collation
			}
		}
	}

	// Apply TiDB's implicit length-driven type promotions (FLOAT(M>24) → DOUBLE,
	// an over-long VARCHAR → TEXT, TEXT(N)/BLOB(N) → the sized variant), the same
	// as the CREATE TABLE path, so an ADD/MODIFY/CHANGE column records the same
	// base type a CREATE does (keeping e.g. FK type comparisons consistent).
	// Skipped for a CHARACTER SET binary column, already mapped to its BLOB form.
	if colDef.TypeName != nil && strings.EqualFold(col.DataType, toLower(colDef.TypeName.Name)) {
		if promoted := promoteDeclaredType(colDef.TypeName, col.Charset); promoted != "" {
			col.DataType = promoted
			// FLOAT → DOUBLE is already reflected in formatColumnType's output
			// (with UNSIGNED/ZEROFILL), so keep that ColumnType; the TEXT/BLOB/
			// VARCHAR promotions yield a length-less, modifier-free type name.
			if !strings.EqualFold(toLower(colDef.TypeName.Name), "float") {
				col.ColumnType = promoted
			}
		}
	}

	// Standalone BINARY modifier (CHAR(10) BINARY) → the effective charset's
	// binary collation, matching TiDB. Resolved after the charset default above
	// so a bare CHAR BINARY uses the table/db charset. An explicit COLLATE wins
	// over BINARY, so only apply when no COLLATE was given.
	if colDef.TypeName != nil && colDef.TypeName.Binary && isStringType(col.DataType) &&
		colDef.TypeName.Collate == "" {
		col.Collation = binModifierCollation(col.Charset)
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
		s := defaultToSQL(colDef.DefaultValue)
		col.Default = &s
	}
	if colDef.OnUpdate != nil {
		col.OnUpdate = onUpdateToSQL(colDef.OnUpdate)
	}
	if colDef.Generated != nil {
		col.Generated = &GeneratedColumnInfo{
			Expr:   generationExpr(colDef.Generated.Expr),
			Stored: colDef.Generated.Stored,
		}
	}

	// Process column-level constraints (non-index-producing ones).
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull:
			col.Nullable = false
		case nodes.ColConstrNull:
			col.Nullable = true
		case nodes.ColConstrDefault:
			if cc.Expr != nil {
				s := defaultToSQL(cc.Expr)
				col.Default = &s
			}
		case nodes.ColConstrAutoIncrement:
			col.AutoIncrement = true
			col.Nullable = false
		case nodes.ColConstrCollate:
			if cc.Expr != nil {
				if s, ok := cc.Expr.(*nodes.StringLit); ok {
					if err := applyColumnCollate(col, s.Value, colDef.TypeName != nil && colDef.TypeName.Charset != ""); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if err := checkLiteralDefaultAllowed(col); err != nil {
		return nil, err
	}
	// ALTER TABLE ADD / MODIFY / CHANGE COLUMN is held to the same
	// invalid-default rules as a column written in CREATE TABLE. The statement
	// carries no table-level PRIMARY KEY of its own, so the 1171 case cannot
	// arise here and the nil constraint list is the whole story.
	if err := checkColumnDefaultValue(colDef, col.DataType, nil); err != nil {
		return nil, err
	}
	if err := checkAlterColumnDefinition(tbl, colDef, col); err != nil {
		return nil, err
	}
	return col, nil
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
	if !c.checkConstraints {
		return nil
	}
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
	// CONVERT TO CHARACTER SET DEFAULT parses and v8.5.5 refuses it with
	// 1105 "Unknown charset " (read on 2026-11-06).
	if strings.EqualFold(cmd.Name, "DEFAULT") {
		return &Error{Code: ErrUnknownError, SQLState: sqlState(ErrUnknownError), Message: "Unknown charset "}
	}
	charset, collation, err := resolveCharsetCollation(cmd.Name, cmd.NewName)
	if err != nil {
		return err
	}
	if err := checkAlterTableCharset(tbl.Charset, tbl.Collation, charset, collation); err != nil {
		return err
	}
	// Each string column's own character set is held to the same changes
	// (v8.5.5: CONVERT TO CHARACTER SET utf8mb4 over a latin1 table with an
	// ascii column is 8200 "Unsupported modify charset from ascii to
	// utf8mb4").
	for _, col := range tbl.Columns {
		if isStringType(col.DataType) && !isEnumSetType(col.DataType) {
			if err := checkColumnCharsetChange(columnCharsetOf(tbl, col), charset); err != nil {
				return err
			}
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
	if err := refuseHistoryPartitions(tbl, cmd.PartitionDefs, true); err != nil {
		return err
	}
	// RANGE and LIST tables require explicit partition definitions; bare
	// ADD PARTITION or ADD PARTITION PARTITIONS n is only for HASH/KEY (1492).
	if len(cmd.PartitionDefs) == 0 {
		switch tbl.Partitioning.Type {
		case "RANGE", "RANGE COLUMNS":
			return errPartitionsMustBeDefined("RANGE")
		case "LIST", "LIST COLUMNS":
			return errPartitionsMustBeDefined("LIST")
		}
		// ADD PARTITION PARTITIONS n names the partitions it appends the way
		// CREATE TABLE does, continuing the pN run: TiDB v8.5.5 answers p0..p4 to
		// `PARTITION BY HASH(id) PARTITIONS 3` followed by `ADD PARTITION
		// PARTITIONS 2`, as MySQL 8.0 does. `ADD PARTITION` with neither a count
		// nor definitions changes nothing, which is also what v8.5.5 does with
		// it on a HASH table.
		if cmd.Number > 0 {
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
		}
		return nil
	}

	// Validate all new partition definitions before mutating the table.
	seenNew := make(map[string]bool, len(cmd.PartitionDefs))
	var toAdd []*PartitionDefInfo
	for _, pd := range cmd.PartitionDefs {
		if err := validateAlterPartitionDef(tbl.Partitioning, pd); err != nil {
			return err
		}
		nameKey := toLower(pd.Name)
		if seenNew[nameKey] {
			return errDuplicatePartitionName(pd.Name)
		}
		if partitionExists(tbl, pd.Name) {
			// IF NOT EXISTS covers a partition the table already has: TiDB
			// v8.5.5 raises warning 1517 and adds nothing.
			if cmd.IfNotExists {
				continue
			}
			return errDuplicatePartitionName(pd.Name)
		}
		seenNew[nameKey] = true
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
		toAdd = append(toAdd, pdi)
	}
	// The partitions it leaves keep rising bounds and distinct values
	// (v8.5.5: ADD PARTITION (PARTITION p2022 VALUES LESS THAN (2023)) below
	// p2023 is 1493, VALUES IN (2, 5) beside IN (1, 2) is 1495).
	all := append(append([]*PartitionDefInfo(nil), tbl.Partitioning.Partitions...), toAdd...)
	if err := checkPartitionSequence(tbl, tbl.Partitioning, all); err != nil {
		return err
	}
	tbl.Partitioning.Partitions = all
	// A HASH or KEY table keeps counting its partitions (v8.5.5: ADD
	// PARTITION (PARTITION p4) over PARTITIONS 4 prints PARTITIONS 5).
	if tbl.Partitioning.NumParts > 0 {
		tbl.Partitioning.NumParts = len(all)
	}
	return nil
}

// validateAlterPartitionDef validates one PartitionDef in ADD PARTITION or
// REORGANIZE PARTITION against the table's partitioning type (measured on
// TiDB v8.5.5): mismatched VALUES IN / VALUES LESS THAN (1480), missing
// VALUES clause (1479), multiple values on non-COLUMNS RANGE (1657), and
// column references inside constant bound expressions (1054).
func validateAlterPartitionDef(pi *PartitionInfo, pd *nodes.PartitionDef) error {
	switch pi.Type {
	case "RANGE", "RANGE COLUMNS":
		switch pd.BoundKind {
		case nodes.PartBoundIn:
			return errOnlyListPartitioningValuesIn()
		case nodes.PartBoundNone:
			return errRangePartitionRequiresValues()
		case nodes.PartBoundLessThan:
			if pi.Type == "RANGE" {
				if list, ok := pd.Values.(*nodes.List); ok && len(list.Items) > 1 {
					return errRangeTooManyValues()
				}
			}
			if name := firstPartitionValueColumnRef(pd.Values); name != "" {
				return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
					Message: fmt.Sprintf("Unknown column '%s' in 'expression'", name)}
			}
		}
	case "LIST", "LIST COLUMNS":
		switch pd.BoundKind {
		case nodes.PartBoundLessThan:
			return errOnlyRangePartitioningValuesLessThan()
		case nodes.PartBoundNone:
			return errListPartitionRequiresValues()
		case nodes.PartBoundIn:
			if name := firstPartitionValueColumnRef(pd.Values); name != "" {
				return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
					Message: fmt.Sprintf("Unknown column '%s' in 'expression'", name)}
			}
		}
	case "HASH", "KEY":
		switch pd.BoundKind {
		case nodes.PartBoundLessThan:
			return errOnlyRangePartitioningValuesLessThan()
		case nodes.PartBoundIn:
			return errOnlyListPartitioningValuesIn()
		}
	}
	return nil
}

// firstPartitionValueColumnRef returns the first column reference inside a
// partition bound node (single expression or tuple/list of expressions), or ""
// when none is present.
func firstPartitionValueColumnRef(v nodes.Node) string {
	switch n := v.(type) {
	case *nodes.List:
		for _, item := range n.Items {
			if name := firstPartitionValueColumnRef(item); name != "" {
				return name
			}
		}
	case nodes.ExprNode:
		return firstColumnRef(n)
	}
	return ""
}

// partitionExists reports whether the table already has a partition of this
// name, matched the way the server matches partition names: case-insensitively.
func partitionExists(tbl *Table, name string) bool {
	if tbl.Partitioning == nil {
		return false
	}
	for _, pd := range tbl.Partitioning.Partitions {
		if toLower(pd.Name) == toLower(name) {
			return true
		}
	}
	return false
}

// alterDropPartition drops named partitions from a partitioned table.
func (c *Catalog) alterDropPartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		// DROP PARTITION IF EXISTS on an unpartitioned table is still an
		// error on TiDB v8.5.5 (1505); IF EXISTS covers a missing partition,
		// not a missing partitioning.
		return errPartitionMgmtOnNonPartitioned()
	}
	// v8.5.5 drops partitions of a RANGE or LIST table alone (1512), each one
	// it names (1507), and not all of them (1508).
	if !strings.HasPrefix(tbl.Partitioning.Type, "RANGE") && !strings.HasPrefix(tbl.Partitioning.Type, "LIST") {
		return errDropOnlyRangeList()
	}
	dropSet := make(map[string]bool)
	for _, name := range cmd.PartitionNames {
		if !partitionExists(tbl, name) {
			if cmd.IfExists {
				continue
			}
			return errDropPartitionList()
		}
		dropSet[toLower(name)] = true
	}
	var remaining []*PartitionDefInfo
	for _, pd := range tbl.Partitioning.Partitions {
		if !dropSet[toLower(pd.Name)] {
			remaining = append(remaining, pd)
		}
	}
	if len(remaining) == 0 {
		return errDropLastPartition()
	}
	tbl.Partitioning.Partitions = remaining
	return nil
}

// alterTruncatePartition truncates named partitions (no-op for metadata catalog).
func (c *Catalog) alterTruncatePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
	}
	// Truncate is a data operation; for DDL catalog purposes, it's a no-op. A
	// partition the table does not have is 1735 (v8.5.5).
	if name := firstUnknownPartition(tbl.Partitioning, cmd.PartitionNames); !cmd.AllPartitions && name != "" {
		return errUnknownPartition(name, tbl.Name)
	}
	return nil
}

// alterCoalescePartition reduces the number of partitions for HASH/KEY partitioned tables.
func (c *Catalog) alterCoalescePartition(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
	}
	if tbl.Partitioning.Type != "HASH" && tbl.Partitioning.Type != "KEY" {
		return errCoalesceOnlyOnHashPartition()
	}
	// Determine current partition count.
	currentCount := len(tbl.Partitioning.Partitions)
	if currentCount == 0 {
		currentCount = tbl.Partitioning.NumParts
	}
	// Coalescing none is 1515, and every partition 1508 (v8.5.5).
	if cmd.Number == 0 {
		return errCoalesceNone()
	}
	if cmd.Number >= currentCount {
		return errDropLastPartition()
	}
	newCount := currentCount - cmd.Number
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
	// v8.5.5 reorganizes RANGE and LIST partitions alone: over HASH or KEY,
	// with partitions named or not, it answers 8200.
	if tbl.Partitioning.Type == "HASH" || tbl.Partitioning.Type == "KEY" {
		return errUnsupportedReorganize()
	}
	// TiDB v8.5.5 (pkg/ddl/executor.go ReorganizePartitions) requires RANGE
	// reorganization to name a non-empty contiguous slice of existing partitions.
	if tbl.Partitioning.Type == "RANGE" || tbl.Partitioning.Type == "RANGE COLUMNS" {
		firstIdx, lastIdx, matched := -1, -1, 0
		nameSet := make(map[string]bool, len(cmd.PartitionNames))
		for _, n := range cmd.PartitionNames {
			nameSet[toLower(n)] = true
		}
		for i, pd := range tbl.Partitioning.Partitions {
			if nameSet[toLower(pd.Name)] {
				matched++
				if firstIdx == -1 {
					firstIdx = i
				}
				lastIdx = i
			}
		}
		if len(cmd.PartitionNames) == 0 || (matched == len(nameSet) && matched != (lastIdx-firstIdx+1)) {
			return &Error{
				Code:     ErrUnsupportedDDL,
				SQLState: sqlState(ErrUnsupportedDDL),
				Message:  "Unsupported REORGANIZE PARTITION of RANGE; not adjacent partitions",
			}
		}
	}
	if err := refuseHistoryPartitions(tbl, cmd.PartitionDefs, false); err != nil {
		return err
	}
	for _, pd := range cmd.PartitionDefs {
		if err := validateAlterPartitionDef(tbl.Partitioning, pd); err != nil {
			return err
		}
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
	// A new name a kept partition holds is 1517, and the partitions left keep
	// rising bounds (v8.5.5: REORGANIZE PARTITION p2024 INTO (PARTITION p2024
	// VALUES LESS THAN (2030)) before p2025's 2026 is 1493).
	seen := map[string]bool{}
	for _, pd := range remaining {
		seen[toLower(pd.Name)] = true
	}
	for _, pd := range newParts {
		if seen[toLower(pd.Name)] {
			return errDuplicatePartitionName(pd.Name)
		}
		seen[toLower(pd.Name)] = true
	}
	if err := checkPartitionSequence(tbl, tbl.Partitioning, result); err != nil {
		return err
	}
	tbl.Partitioning.Partitions = result
	return nil
}

// alterExchangePartition exchanges a partition with a non-partitioned table.
func (c *Catalog) alterExchangePartition(db *Database, tbl *Table, cmd *nodes.AlterTableCmd) error {
	if tbl.Partitioning == nil {
		return errPartitionMgmtOnNonPartitioned()
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
		if exchTbl.Partitioning != nil {
			return errExchangeTableIsPartitioned(exchTbl.Name)
		}
	}
	if name := firstUnknownPartition(tbl.Partitioning, []string{cmd.Name}); name != "" {
		return errUnknownPartition(name, tbl.Name)
	}
	return nil
}

// checkModifyColumnFKCompatibility verifies that replacing oldName with newCol
// keeps every foreign key referencing or referenced by the column type-
// compatible (3780). A NOT NULL column under a SET NULL action is taken: v8.5.5
// runs MODIFY x INT NOT NULL under ON UPDATE SET NULL and ON DELETE SET NULL.
func (c *Catalog) checkModifyColumnFKCompatibility(tbl *Table, oldName string, newCol *Column) error {
	oldKey := toLower(oldName)
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey {
			continue
		}
		for i, colName := range con.Columns {
			if toLower(colName) != oldKey {
				continue
			}
			if i >= len(con.RefColumns) {
				continue
			}
			refDB := tbl.Database
			if con.RefDatabase != "" {
				refDB = c.GetDatabase(con.RefDatabase)
			}
			var refTbl *Table
			if refDB != nil {
				if tbl.Database != nil && refDB == tbl.Database && toLower(con.RefTable) == toLower(tbl.Name) {
					refTbl = tbl
				} else {
					refTbl = refDB.GetTable(con.RefTable)
				}
			}
			if refTbl == nil {
				continue
			}
			var refCol *Column
			if refTbl == tbl && toLower(con.RefColumns[i]) == oldKey {
				refCol = newCol
			} else {
				refCol = refTbl.GetColumn(con.RefColumns[i])
			}
			if refCol != nil && !fkTypesCompatible(newCol, refCol) {
				return errFKIncompatibleColumns(newCol.Name, refCol.Name, con.Name)
			}
		}
	}
	if tbl.Database == nil {
		return nil
	}
	dbKey := toLower(tbl.Database.Name)
	tblKey := toLower(tbl.Name)
	for _, otherDB := range c.databases {
		for _, child := range otherDB.Tables {
			for _, con := range child.Constraints {
				if con.Type != ConForeignKey {
					continue
				}
				refDB := con.RefDatabase
				if refDB == "" {
					refDB = otherDB.Name
				}
				if toLower(refDB) != dbKey || toLower(con.RefTable) != tblKey {
					continue
				}
				for i, refColName := range con.RefColumns {
					if toLower(refColName) != oldKey || i >= len(con.Columns) {
						continue
					}
					var childCol *Column
					if child == tbl && toLower(con.Columns[i]) == oldKey {
						childCol = newCol
					} else {
						childCol = child.GetColumn(con.Columns[i])
					}
					if childCol != nil && !fkTypesCompatible(childCol, newCol) {
						return errFKIncompatibleColumns(childCol.Name, newCol.Name, con.Name)
					}
				}
			}
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
	// A backtick-quoted identifier outside any string literal: '`id`' is a
	// string, not the column.
	target := strings.ToLower(colName)
	for i := 0; i < len(expr); i++ {
		switch ch := expr[i]; ch {
		case '\'', '"':
			for i++; i < len(expr) && expr[i] != ch; i++ {
				if expr[i] == '\\' {
					i++
				}
			}
		case '`':
			start := i + 1
			for i++; i < len(expr); i++ {
				if expr[i] == '`' {
					if i+1 < len(expr) && expr[i+1] == '`' {
						i++
						continue
					}
					break
				}
			}
			if i <= len(expr) && strings.ToLower(strings.ReplaceAll(expr[start:min(i, len(expr))], "``", "`")) == target {
				return true
			}
		}
	}
	return false
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
// TiDB v8.5.5 leaves the table untouched for all of them:
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
