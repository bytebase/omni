package v1review

import (
	"fmt"
	"strconv"

	"github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// PriorBackupStatement is one statement of a review BUILTIN_PRIOR_BACKUP_CHECK
// judges: its node, and the statement's own text, the text the node's
// locations are relative to.
type PriorBackupStatement struct {
	Node ast.Node
	Text string
}

// PriorBackupCheck reports what keeps a prior backup from taking the
// statements, in the V1 advisor's order: the reviewed SQL is larger than facts
// allows, a statement is DDL, the database the backup writes to is missing, or
// one table is changed by more than one kind of DML. V1 reported every DDL
// statement, whether or not the batch changes rows as well. Tables are grouped
// as "database.table", an unqualified reference resolved to schema's database.
//
// V1 reported the mixed DML tables in its map's order, which varied between
// runs; they come here in the order the statements first name them. The
// findings are in V1's order otherwise: size, then each mixed DDL statement in
// statement order, then the missing database, then the mixed DML tables. A
// mixed DDL finding carries the index of its statement and no range, since V1
// anchored it at the statement's start; the others anchor nowhere. Ported from
// Bytebase's V1 rule BUILTIN_PRIOR_BACKUP_CHECK; Bytebase calls it only when
// the review enables prior backup.
func PriorBackupCheck(stmts []PriorBackupStatement, schema v1review.Schema, facts v1review.PriorBackupFacts) ([]v1review.Finding, error) {
	var findings []v1review.Finding

	if facts.TotalSize > facts.MaxSize {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupScriptTooLarge, Value: strconv.Itoa(facts.MaxSize)})
	}

	for i, stmt := range stmts {
		if isDDL(stmt.Node) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDDLAndDML, Statement: i})
		}
	}

	if !facts.DatabaseExists {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupDatabaseMissing, Name: facts.Database})
	}

	var order []string
	types := make(map[string][]analysis.DMLType)
	for i, stmt := range stmts {
		targets, err := analysis.ExtractTables(schema.DatabaseName(), stmt.Node, stmt.Text, i, analysis.Catalog{})
		if err != nil {
			return nil, err
		}
		for _, target := range targets {
			if target.Table == nil {
				continue
			}
			key := fmt.Sprintf("%s.%s", target.Table.Database, target.Table.Table)
			if _, ok := types[key]; !ok {
				order = append(order, key)
			}
			types[key] = append(types[key], target.Table.StatementType)
		}
	}
	for _, key := range order {
		if mixedDML(types[key]) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDMLOnTable, Table: key})
		}
	}

	return findings, nil
}

// mixedDML reports whether types holds more than one kind of change, reading
// them as the V1 advisor did: a comparison against the first kind it read that
// it knew, so a leading unknown kind compares as nothing.
func mixedDML(types []analysis.DMLType) bool {
	first := analysis.DMLTypeUnknown
	for _, t := range types {
		if first == analysis.DMLTypeUnknown {
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
	case *ast.CreateTableStmt, *ast.CreateDatabaseStmt, *ast.CreateIndexStmt,
		*ast.CreateViewStmt, *ast.CreateTriggerStmt, *ast.CreateEventStmt,
		*ast.CreateFunctionStmt,
		*ast.AlterTableStmt, *ast.AlterDatabaseStmt, *ast.AlterViewStmt,
		*ast.AlterEventStmt,
		*ast.DropTableStmt, *ast.DropDatabaseStmt, *ast.DropIndexStmt,
		*ast.DropViewStmt, *ast.DropTriggerStmt, *ast.DropEventStmt,
		*ast.DropRoutineStmt,
		*ast.RenameTableStmt, *ast.TruncateStmt:
		return true
	default:
		return false
	}
}
