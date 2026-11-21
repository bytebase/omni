package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// viewAsTable builds a read-only proxy Table representing a view's resolved
// columns, with no indexes, constraints, or partitioning. OceanBase's resolver
// validates many ALTER TABLE and CREATE INDEX clauses against the view's table
// schema before RootService rejects the non-base-table target with 1347.
func (c *Catalog) viewAsTable(db *Database, v *View) *Table {
	vt := &Table{
		Name:      v.Name,
		Database:  db,
		Charset:   db.Charset,
		Collation: db.Collation,
		colByName: make(map[string]int),
	}
	for i, cm := range c.viewColumnsToProto(db, v) {
		colType := cm.GetType()
		col := &Column{
			Position:   i + 1,
			Name:       cm.GetName(),
			DataType:   baseTypeFromColumnType(colType),
			ColumnType: colType,
			Nullable:   cm.GetNullable(),
			Charset:    cm.GetCharacterSet(),
			Collation:  cm.GetCollation(),
		}
		vt.Columns = append(vt.Columns, col)
		vt.colByName[toLower(col.Name)] = i
	}
	return vt
}

// baseTypeFromColumnType extracts the leading SQL type word from a formatted
// ColumnType string such as "int(11) unsigned" -> "int".
func baseTypeFromColumnType(colType string) string {
	s := strings.TrimSpace(colType)
	if idx := strings.IndexAny(s, " ("); idx >= 0 {
		s = s[:idx]
	}
	return strings.ToLower(s)
}

// alterViewAsTable runs OceanBase 4.3.5.6's ALTER TABLE resolver checks when
// the target relation is a view. Resolver-stage errors (table options, column
// definition semantics, key part types, check expressions, partition clauses,
// column groups, unknown charset/collation) are returned ahead of RootService's
// 1347 ('<db>.<view>' is not BASE TABLE). A statement consisting solely of
// COMMENT table options succeeds as a no-op on the view, matching 4.3.5.6.
func (c *Catalog) alterViewAsTable(db *Database, v *View, stmt *nodes.AlterTableStmt) error {
	vt := c.viewAsTable(db, v)

	// Phase 1: resolver-stage table option validation.
	// STORAGE_FORMAT_VERSION, TABLEGROUP, and DYNAMIC_PARTITION_POLICY are
	// checked in RootService after the 1347 base-table check.
	for _, cmd := range stmt.Commands {
		if cmd.Type != nodes.ATTableOption || cmd.Option == nil {
			continue
		}
		upper := strings.ToUpper(cmd.Option.Name)
		switch upper {
		case "STORAGE_FORMAT_VERSION", "TABLEGROUP", "DYNAMIC_PARTITION_POLICY":
			continue
		case "TABLE_ID":
			return &Error{Code: 1064, SQLState: "42000", Message: "TABLE_ID cannot be altered"}
		case "CHARSET", "CHARACTER SET", "DEFAULT CHARSET", "DEFAULT CHARACTER SET":
			if err := checkCharsetName(normalizeCharsetName(cmd.Option.Value)); err != nil {
				return err
			}
			continue
		case "COLLATE", "DEFAULT COLLATE":
			if err := checkCollationName(cmd.Option.Value); err != nil {
				return err
			}
			continue
		case "COMPRESSION":
			if err := validateCompressor(cmd.Option.Value); err != nil {
				return err
			}
			continue
		case "TTL":
			if err := checkTTL(vt, cmd.Option); err != nil {
				return err
			}
			continue
		}
		if err := checkAlterTableOption(vt, cmd.Option); err != nil {
			return err
		}
	}

	if err := checkKeyOptions(alterKeys(stmt.Commands)); err != nil {
		return err
	}

	if repartition := alterRepartition(stmt.Commands); repartition != nil {
		if repartition.Type == nodes.PartitionKey && len(repartition.Columns) == 0 {
			return &Error{Code: 1488, SQLState: "HY000", Message: "Field in list of fields for partition function not found in table"}
		}
		if repartition.Type == nodes.PartitionRange && ((repartition.Expr == nil && len(repartition.Columns) == 0) || len(repartition.Partitions) == 0 || repartition.SplitSize != "") {
			return errNotSupported("using auto-partition clause for non-user-table is")
		}
		if err := checkPartitionStructure(vt, repartition, true); err != nil {
			return err
		}
		return errUniqueKeyNeedAllFieldsInPF("PRIMARY KEY")
	}

	onlyComment := len(stmt.Commands) > 0
	addedCols := make(map[string]bool)
	for _, cmd := range stmt.Commands {
		if cmd.Type != nodes.ATTableOption || cmd.Option == nil || !strings.EqualFold(cmd.Option.Name, "COMMENT") {
			onlyComment = false
		}
		if err := c.checkAlterViewCmd(vt, cmd, addedCols); err != nil {
			return err
		}
	}
	if onlyComment {
		return nil
	}
	return errNotBaseTable(db.Name, stmt.Table.Name)
}

func (c *Catalog) checkAlterViewCmd(vt *Table, cmd *nodes.AlterTableCmd, addedCols map[string]bool) error {
	switch cmd.Type {
	case nodes.ATAddColumn:
		if alterCarriesColumnID(cmd) {
			return &Error{Code: 1064, SQLState: "42000", Message: "a column id cannot be altered"}
		}
		if cmd.LobStorage != nil {
			return &Error{Code: 1210, SQLState: "HY000", Message: "Invalid argument"}
		}
		defs := cmd.Columns
		if len(defs) == 0 && cmd.Column != nil {
			defs = []*nodes.ColumnDef{cmd.Column}
		}
		for _, colDef := range defs {
			if err := checkAddedGeneratedCheck(colDef); err != nil {
				return err
			}
			if err := validateColumnDefSemanticsMode(colDef, true); err != nil {
				return err
			}
			if err := checkAddedColumn(vt, colDef); err != nil {
				return err
			}
			key := toLower(colDef.Name)
			if addedCols[key] {
				return errDupColumn(colDef.Name)
			}
			addedCols[key] = true
		}
	case nodes.ATDropColumn:
		if !cmd.IfExists && vt.GetColumn(cmd.Name) == nil {
			return errCantDropKey(cmd.Name)
		}
	case nodes.ATAddIndex, nodes.ATAddConstraint:
		con := cmd.Constraint
		if con == nil {
			return nil
		}
		if kind, ok := constraintKeyKind(con); ok {
			if err := checkKeyPartsMode(vt, buildIndexColumns(con), kind, con.IndexOptions, true); err != nil {
				return err
			}
		}
		switch con.Type {
		case nodes.ConstrPrimaryKey:
			return errMultiplePriKey()
		case nodes.ConstrCheck:
			conName := con.Name
			if conName == "" {
				conName = fmt.Sprintf("%s%s0", vt.Name, madeUpCheckInfix)
			}
			if err := checkTableCheckExpr(vt, conName, con.Expr); err != nil {
				return err
			}
			if err := validateCheckExpr(conName, con.Expr); err != nil {
				return err
			}
		case nodes.ConstrForeignKey:
			if err := rejectRemoteRef(con.RefTable); err != nil {
				return err
			}
			return errCannotAddForeign()
		}
	case nodes.ATDropConstraint:
		if len(cmd.ConstraintNames) > 0 {
			if !cmd.IfExists {
				return errCheckConstraintNotFound(cmd.ConstraintNames[0])
			}
			return nil
		}
		if cmd.IfExists {
			return nil
		}
		switch cmd.ConstraintKeyword {
		case nodes.ConstraintKeywordCheck:
			return errCheckConstraintNotFound(cmd.Name)
		case nodes.ConstraintKeywordConstraint:
			return errConstraintDoesNotExist(cmd.Name)
		}
	case nodes.ATAlterCheckEnforced:
		if cmd.ConstraintKeyword == nodes.ConstraintKeywordCheck {
			return errCheckConstraintNotFound(cmd.Name)
		}
		return errConstraintDoesNotExist(cmd.Name)
	case nodes.ATConvertCharset:
		if err := checkCharsetName(normalizeCharsetName(cmd.Name)); err != nil {
			return err
		}
		if cmd.NewName != "" {
			if err := checkCollationName(cmd.NewName); err != nil {
				return err
			}
			if err := checkCharsetCollationPair(cmd.Name, cmd.NewName); err != nil {
				return err
			}
		}
	case nodes.ATAlterColumnDefault:
		if cmd.DefaultExpr != nil && hasColumnRef(cmd.DefaultExpr) {
			return errNotSupportedDefault()
		}
	case nodes.ATAlterIndexStorageCachePolicy:
		return errNotSupported("storage cache policy")
	case nodes.ATAddPartition, nodes.ATDropPartition, nodes.ATTruncatePartition,
		nodes.ATDropSubpartition, nodes.ATTruncateSubpartition, nodes.ATAlterPartitionStorageCachePolicy,
		nodes.ATAddSubpartition, nodes.ATCoalescePartition, nodes.ATReorganizePartition, nodes.ATExchangePartition:
		return errNotPartitioned()
	case nodes.ATRefresh:
		return &Error{Code: 1210, SQLState: "HY000", Message: "Incorrect arguments to REFRESH"}
	case nodes.ATAddColumnGroup, nodes.ATDropColumnGroup:
		return errColumnGroupNotSupported()
	}
	return nil
}
