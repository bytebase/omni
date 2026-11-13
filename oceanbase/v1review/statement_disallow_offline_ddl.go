package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowOfflineDDLLookups returns, in order and with repeats, the
// columns the V1 rule STATEMENT_DISALLOW_OFFLINE_DDL looked up in the
// database: each column a top-level ALTER TABLE ... DROP COLUMN drops, in the
// table's database or else currentDatabase, when the database, table and
// column are all named. Bytebase reads each column's
// INFORMATION_SCHEMA.COLUMNS.EXTRA and passes the values to
// StatementDisallowOfflineDDL.
func StatementDisallowOfflineDDLLookups(stmts []ast.Node, currentDatabase string) []v1review.ColumnRef {
	var refs []v1review.ColumnRef
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		database, table := offlineDDLTable(alter, currentDatabase)
		for _, cmd := range alter.Commands {
			if cmd == nil || cmd.Type != ast.ATDropColumn {
				continue
			}
			if database != "" && table != "" && cmd.Name != "" {
				refs = append(refs, v1review.ColumnRef{Database: database, Table: table, Column: cmd.Name})
			}
		}
	}
	return refs
}

// StatementDisallowOfflineDDL reports each operation top-level statements run
// that OceanBase performs offline, Name the operation as the V1 advisor
// printed it: DROP TABLE ("Dropping tables") and TRUNCATE ("Truncating
// tables"), anchoring the statement; and per ALTER TABLE command, anchoring
// the command: adding a column that is stored generated, auto-increment or a
// primary key, checked in that order ("Adding <kind> columns", once for the
// command's column and once for each of its column list); adding a primary
// key; modifying or changing a column to such a column ("Modifying columns
// to <kind>", "Changing columns to <kind>"); dropping a column whose EXTRA,
// looked up by the ref StatementDisallowOfflineDDLLookups gave, contains
// "stored" in any case; dropping the constraint named PRIMARY in any case;
// dropping or truncating partitions; every other partition operation
// ("Modifying partitions"); converting the charset; and a table option whose
// name contains charset, character set, collate or collation in any case
// ("Changing charset or collate"). A column with no EXTRA in extras is not
// stored. OfflineDDL is the only code. Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_OFFLINE_DDL.
func StatementDisallowOfflineDDL(stmts []ast.Node, currentDatabase string, extras map[v1review.ColumnRef]string) []v1review.Finding {
	var findings []v1review.Finding
	add := func(loc ast.Loc, operation string) {
		start, end := anchor(loc)
		findings = append(findings, v1review.Finding{Code: v1review.OfflineDDL, Name: operation, Start: start, End: end})
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.AlterTableStmt:
			database, table := offlineDDLTable(n, currentDatabase)
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				offlineDDLCommand(cmd, database, table, extras, func(operation string) { add(cmd.Loc, operation) })
			}
		case *ast.DropTableStmt:
			add(n.Loc, "Dropping tables")
		case *ast.TruncateStmt:
			add(n.Loc, "Truncating tables")
		}
	}
	return findings
}

// offlineDDLTable returns the database and table an ALTER TABLE names, the
// database defaulting to currentDatabase.
func offlineDDLTable(alter *ast.AlterTableStmt, currentDatabase string) (string, string) {
	var database, table string
	if alter.Table != nil {
		database, table = alter.Table.Schema, alter.Table.Name
	}
	if database == "" {
		database = currentDatabase
	}
	return database, table
}

func offlineDDLCommand(cmd *ast.AlterTableCmd, database, table string, extras map[v1review.ColumnRef]string, add func(string)) {
	switch cmd.Type {
	case ast.ATAddColumn:
		if kind := offlineColumnKind(cmd.Column); kind != "" {
			add("Adding " + kind + " columns")
		}
		for _, column := range cmd.Columns {
			if kind := offlineColumnKind(column); kind != "" {
				add("Adding " + kind + " columns")
			}
		}
	case ast.ATAddConstraint:
		if cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrPrimaryKey {
			add("Adding primary keys")
		}
	case ast.ATModifyColumn:
		if kind := offlineColumnKind(cmd.Column); kind != "" {
			add("Modifying columns to " + kind)
		}
	case ast.ATChangeColumn:
		if kind := offlineColumnKind(cmd.Column); kind != "" {
			add("Changing columns to " + kind)
		}
	case ast.ATDropColumn:
		if database != "" && table != "" && cmd.Name != "" {
			if extra, ok := extras[v1review.ColumnRef{Database: database, Table: table, Column: cmd.Name}]; ok && strings.Contains(strings.ToLower(extra), "stored") {
				add("Dropping stored generated columns")
			}
		}
	case ast.ATDropConstraint:
		if strings.EqualFold(cmd.Name, "PRIMARY") {
			add("Dropping primary keys")
		}
	case ast.ATDropPartition:
		add("Dropping partitions")
	case ast.ATTruncatePartition:
		add("Truncating partitions")
	case ast.ATAddPartition, ast.ATCoalescePartition, ast.ATReorganizePartition, ast.ATExchangePartition, ast.ATAnalyzePartition, ast.ATCheckPartition, ast.ATOptimizePartition, ast.ATRebuildPartition, ast.ATRepairPartition, ast.ATDiscardPartitionTablespace, ast.ATImportPartitionTablespace, ast.ATRemovePartitioning, ast.ATPartitionBy:
		add("Modifying partitions")
	case ast.ATConvertCharset:
		add("Changing charset or collate")
	case ast.ATTableOption:
		if cmd.Option != nil {
			name := strings.ToLower(cmd.Option.Name)
			if strings.Contains(name, "charset") || strings.Contains(name, "character set") || strings.Contains(name, "collate") || strings.Contains(name, "collation") {
				add("Changing charset or collate")
			}
		}
	}
}

// offlineColumnKind names what makes adding a column offline, empty for
// nothing: stored generated, then auto increment, then primary key.
func offlineColumnKind(column *ast.ColumnDef) string {
	if column == nil {
		return ""
	}
	switch {
	case column.Generated != nil && column.Generated.Stored:
		return "stored generated"
	case column.AutoIncrement:
		return "auto increment"
	}
	for _, constraint := range column.Constraints {
		if constraint.Type == ast.ColConstrPrimaryKey {
			return "primary key"
		}
	}
	return ""
}
