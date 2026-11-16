package v1review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
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
// database name, and whether the statement wrote the qualifier kept, since
// the LONG/LONG RAW lookup reads the snapshot differently for each.
type priorBackupTable struct {
	schema        string
	hasSchema     bool
	table         string
	statementType priorBackupStatementType
}

// PriorBackupCheck reports what keeps a prior backup from taking the
// statements, in the V1 advisor's order: the reviewed SQL is larger than facts
// allows; the database the backup writes to is missing, which ends the check;
// the batch holds a DDL statement; one table is changed by more than one kind
// of DML; and a table the statements change holds a LONG or LONG RAW column,
// which CREATE TABLE AS SELECT cannot copy (ORA-00997). Tables are grouped as
// "schema.table", an unqualified reference resolved to schema's database. V1
// reported the DDL once for the batch and anchored it at the script's start,
// so that finding carries nothing else. It reported the mixed DML tables and
// the LONG tables in its map's order, which varied between runs; they come
// here in the order the statements first name them. Ported from Bytebase's V1
// rule BUILTIN_PRIOR_BACKUP_CHECK; Bytebase calls it only when the review
// enables prior backup.
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
	changed := make(map[string]*priorBackupTable)
	hasDDL := false
	for _, stmt := range stmts {
		var table *priorBackupTable
		switch n := stmt.(type) {
		case *ast.UpdateStmt:
			table = priorBackupTableOf(schema.DatabaseName(), n.Table, priorBackupStatementUpdate)
		case *ast.DeleteStmt:
			table = priorBackupTableOf(schema.DatabaseName(), n.Table, priorBackupStatementDelete)
		default:
			if isDDL(stmt) {
				hasDDL = true
			}
		}
		if table == nil {
			continue
		}
		key := fmt.Sprintf("%s.%s", table.schema, table.table)
		if _, ok := types[key]; !ok {
			order = append(order, key)
		}
		types[key] = append(types[key], table.statementType)
		changed[key] = table
	}

	if hasDDL {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupDDLInBatch})
	}

	for _, key := range order {
		if mixedDML(types[key]) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDMLOnTable, Table: key})
		}
	}

	for _, key := range order {
		table := changed[key]
		columns := longColumns(schema, table)
		if len(columns) == 0 {
			continue
		}
		findings = append(findings, v1review.Finding{
			Code:   v1review.PriorBackupUncopyableColumn,
			Schema: table.schema,
			Table:  table.table,
			Value:  strings.Join(columns, ", "),
		})
	}

	return findings
}

// priorBackupTableOf returns the table name a DML changes, nil when it names
// no table. The alias is not kept: V1 grouped the table the DML names.
func priorBackupTableOf(databaseName string, name *ast.ObjectName, statementType priorBackupStatementType) *priorBackupTable {
	if name == nil {
		return nil
	}
	schemaName := name.Schema
	hasSchema := schemaName != ""
	if schemaName == "" {
		schemaName = databaseName
	}
	return &priorBackupTable{schema: schemaName, hasSchema: hasSchema, table: name.Name, statementType: statementType}
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
// CREATE USER, a GRANT, or a CREATE SEQUENCE is not in it.
func isDDL(node ast.Node) bool {
	switch node.(type) {
	case *ast.CreateTableStmt, *ast.AlterTableStmt, *ast.DropStmt, *ast.CreateIndexStmt,
		*ast.CreateViewStmt, *ast.TruncateStmt, *ast.CommentStmt:
		return true
	default:
		return false
	}
}

// nameEqual compares an AST object name with a snapshot name the way the V1
// advisor did. The parser delivers names in catalog form — an unquoted
// identifier upper-cased, a quoted one exact-case with the quotes stripped —
// and sync stores data-dictionary names verbatim, so in case-sensitive mode
// this is plain equality: folding there would conflate T_LONG with a distinct
// quoted "t_long" table and warn about the wrong table's columns.
func nameEqual(astName, snapshotName string, objectCaseSensitive bool) bool {
	if objectCaseSensitive {
		return astName == snapshotName
	}
	return strings.EqualFold(astName, snapshotName)
}

// longColumns returns the names of the LONG and LONG RAW columns of table,
// each double-quoted as V1 printed it, or nil when the snapshot does not hold
// the table. Missing metadata must not produce false alarms.
//
// The snapshot stores the connection schema's tables under a schema named "",
// while an unqualified DML is resolved to the database name, so the empty
// schema matches the current owner alone: a DML explicitly qualified with
// another owner stays silent rather than borrowing the current schema's table
// of the same name.
func longColumns(schema v1review.Schema, table *priorBackupTable) []string {
	var result []string
	for _, schemaMeta := range schema.Metadata.GetSchemas() {
		if schemaMeta.GetName() == "" {
			if table.hasSchema && !nameEqual(table.schema, schema.Metadata.GetName(), schema.ObjectCaseSensitive) {
				continue
			}
		} else if !nameEqual(table.schema, schemaMeta.GetName(), schema.ObjectCaseSensitive) {
			continue
		}
		for _, tableMeta := range schemaMeta.GetTables() {
			if !nameEqual(table.table, tableMeta.GetName(), schema.ObjectCaseSensitive) {
				continue
			}
			for _, column := range tableMeta.GetColumns() {
				switch strings.ToUpper(strings.TrimSpace(column.GetType())) {
				case "LONG", "LONG RAW":
					result = append(result, fmt.Sprintf("%q", column.GetName()))
				}
			}
		}
	}
	return result
}
