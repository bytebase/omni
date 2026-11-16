package v1review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// priorBackupStatementType is the kind of change a top-level DML makes. The
// V1 rule kept INSERT in its enum but never collected an INSERT, so no
// statement is grouped as one.
type priorBackupStatementType int

const (
	priorBackupStatementUnknown priorBackupStatementType = iota
	priorBackupStatementUpdate
	priorBackupStatementDelete
)

// priorBackupTable is a table a top-level DML changes, qualified the way the
// V1 rule grouped tables: an unqualified name resolved to the snapshot's
// database and to dbo, and an alias in the FROM clause followed to the
// physical table it names.
type priorBackupTable struct {
	database      string
	schema        string
	table         string
	statementType priorBackupStatementType
}

// PriorBackupCheck reports what keeps a prior backup from taking the
// statements, in the V1 advisor's order: the reviewed SQL is larger than facts
// allows; the database the backup writes to is missing, which ends the check;
// the batch holds a DDL statement; and one table is changed by more than one
// kind of DML. Tables are grouped as "database.schema.table", an unqualified
// reference resolved to schema's database and dbo, and a temp table skipped
// because the backup transformer never copies one. V1 reported the DDL once
// for the batch and anchored every finding at the script's first line, so no
// finding carries a statement index or a range; the tables of the mixed DML
// findings come in the order the statements first name them rather than in
// the map order V1 read them in. Ported from Bytebase's V1 rule
// BUILTIN_PRIOR_BACKUP_CHECK; Bytebase calls it only when the review enables
// prior backup.
func PriorBackupCheck(stmts []ast.Node, schema v1review.Schema, facts v1review.PriorBackupFacts) []v1review.Finding {
	var findings []v1review.Finding

	if facts.TotalSize > facts.MaxSize {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupScriptTooLarge, Value: strconv.Itoa(facts.MaxSize)})
	}

	// The missing database ended V1's check: it reported neither the DDL nor
	// the tables a backup would have to copy.
	if !facts.DatabaseExists {
		return append(findings, v1review.Finding{Code: v1review.PriorBackupDatabaseMissing, Name: facts.Database})
	}

	var order []string
	types := make(map[string][]priorBackupStatementType)
	hasDDL := false
	for _, stmt := range stmts {
		var table *priorBackupTable
		switch n := stmt.(type) {
		case *ast.UpdateStmt:
			table = priorBackupTableOf(n.Relation, n.FromClause, schema.DatabaseName(), priorBackupStatementUpdate)
		case *ast.DeleteStmt:
			table = priorBackupTableOf(n.Relation, n.FromClause, schema.DatabaseName(), priorBackupStatementDelete)
		default:
			if isDDL(stmt) {
				hasDDL = true
			}
		}
		if table == nil || table.table == "" {
			continue
		}
		key := fmt.Sprintf("%s.%s.%s", table.database, table.schema, table.table)
		if _, ok := types[key]; !ok {
			order = append(order, key)
		}
		types[key] = append(types[key], table.statementType)
	}

	if hasDDL {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupDDLInBatch})
	}

	for _, key := range order {
		if mixedDML(types[key]) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDMLOnTable, Table: key})
		}
	}

	return findings
}

// priorBackupTableOf returns the table a DML changes, nil when the target is
// not a table reference V1 resolved.
func priorBackupTableOf(relation ast.TableExpr, fromClause *ast.List, databaseName string, statementType priorBackupStatementType) *priorBackupTable {
	table := tableReferenceOf(relation, databaseName, defaultSchema)
	if table == nil {
		return nil
	}
	// An UPDATE or DELETE can name an alias the FROM clause declares; V1
	// followed it to the physical table, but only for a target it had
	// resolved to the snapshot's own database and dbo.
	if fromClause != nil && table.database == databaseName && table.schema == defaultSchema {
		if physical := physicalTableOfAlias(fromClause, table); physical != nil {
			table = physical
		}
	}
	// A temp table is the session's own and prior backup never copies it, so
	// statements on one cannot conflict with each other.
	if strings.HasPrefix(table.table, "#") {
		return nil
	}
	table.statementType = statementType
	return table
}

// tableReferenceOf returns the table a table expression names, nil when it is
// not a table reference.
func tableReferenceOf(expr ast.TableExpr, defaultDatabase, schema string) *priorBackupTable {
	ref, ok := expr.(*ast.TableRef)
	if !ok || ref == nil {
		return nil
	}
	schemaName := schema
	if ref.Schema != "" {
		schemaName = ref.Schema
	}
	databaseName := defaultDatabase
	if ref.Database != "" {
		databaseName = ref.Database
	}
	return &priorBackupTable{database: databaseName, schema: schemaName, table: ref.Object}
}

// physicalTableOfAlias returns the table the FROM list declares under the
// name table holds, nil when no item declares it.
func physicalTableOfAlias(list *ast.List, table *priorBackupTable) *priorBackupTable {
	for _, item := range list.Items {
		if result := physicalTableOfAliasIn(item, table); result != nil {
			return result
		}
	}
	return nil
}

func physicalTableOfAliasIn(node ast.Node, table *priorBackupTable) *priorBackupTable {
	switch n := node.(type) {
	case *ast.TableRef:
		if n.Alias != "" && n.Alias == table.table {
			return tableReferenceOf(n, table.database, table.schema)
		}
	case *ast.AliasedTableRef:
		if ref, ok := n.Table.(*ast.TableRef); ok && n.Alias == table.table {
			return tableReferenceOf(ref, table.database, table.schema)
		}
	case *ast.JoinClause:
		if result := physicalTableOfAliasIn(n.Left, table); result != nil {
			return result
		}
		return physicalTableOfAliasIn(n.Right, table)
	default:
	}
	return nil
}

// mixedDML reports whether types holds more than one kind of change, reading
// them as the V1 advisor did: a comparison against the first kind it read that
// it knew, so a leading unknown kind compares as nothing.
func mixedDML(types []priorBackupStatementType) bool {
	first := priorBackupStatementUnknown
	for _, t := range types {
		if first == priorBackupStatementUnknown {
			first = t
		}
		if first != t {
			return true
		}
	}
	return false
}

// isDDL reports whether the V1 rule counted node as a DDL statement. The list
// is the one the V1 advisor kept, not every statement that changes a schema: a
// CREATE USER or a GRANT is not in it.
func isDDL(node ast.Node) bool {
	switch node.(type) {
	case *ast.CreateTableStmt, *ast.AlterTableStmt, *ast.DropStmt, *ast.TruncateStmt,
		*ast.CreateIndexStmt, *ast.CreateViewStmt, *ast.CreateFunctionStmt,
		*ast.CreateProcedureStmt, *ast.CreateSchemaStmt, *ast.CreateDatabaseStmt,
		*ast.CreateTriggerStmt, *ast.CreateTypeStmt, *ast.CreateSequenceStmt,
		*ast.AlterIndexStmt, *ast.AlterSchemaStmt, *ast.AlterDatabaseStmt,
		*ast.AlterSequenceStmt, *ast.RenameStmt:
		return true
	default:
		return false
	}
}
