package v1review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/ast"
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
// V1 rule grouped tables: by its schema, "public" when the statement names
// none, and its table name. V1 read the catalog a statement names and never
// used it, so a name qualified by another database groups with the same schema
// and table of the reviewed one.
type priorBackupTable struct {
	schema        string
	table         string
	statementType priorBackupStatementType
}

// PriorBackupCheck reports what keeps a prior backup from taking the
// statements, in the V1 advisor's order: each statement is DDL; each statement
// holds a data-modifying WITH query; the database the backup writes to is
// missing, which does not end the check; and one table is changed by more than
// one kind of DML. V1 reported every DDL statement, whether or not the batch
// changes rows as well. Tables are grouped as "schema.table".
//
// The size finding is the one judgment this rule adds to the V1 advisor's, and
// it is the sibling engines' rule: V1's PostgreSQL advisor had no size cap, so
// a reviewed script above facts.MaxSize passed review and then failed the
// backup task, which refuses one ("statement size %d exceeds the limit %d,
// please disable data backup").
//
// A statement carries the index of the statement it reports and no range,
// since V1 anchored every finding at the first content line of the statement's
// text. V1 reported the mixed DML tables in its map's order, which varied
// between runs; they come here in the order the statements first name them.
// Ported from Bytebase's V1 rule BUILTIN_PRIOR_BACKUP_CHECK; Bytebase calls it
// only when the review enables prior backup.
func PriorBackupCheck(stmts []ast.Node, facts v1review.PriorBackupFacts) []v1review.Finding {
	var findings []v1review.Finding

	if facts.TotalSize > facts.MaxSize {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupScriptTooLarge, Value: strconv.Itoa(facts.MaxSize)})
	}

	for i, stmt := range stmts {
		node := statementNode(stmt)
		if isDDLStatement(node) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDDLAndDML, Statement: i})
		}
		if names := analysis.DataModifyingCTENames(node); len(names) > 0 {
			quoted := make([]string, 0, len(names))
			for _, name := range names {
				quoted = append(quoted, fmt.Sprintf("%q", name))
			}
			findings = append(findings, v1review.Finding{
				Code:      v1review.PriorBackupDataModifyingCTE,
				Value:     strings.Join(quoted, ", "),
				Statement: i,
			})
		}
	}

	if !facts.DatabaseExists {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupDatabaseMissing, Name: facts.Database})
	}

	var order []string
	types := make(map[string][]priorBackupStatementType)
	for _, stmt := range stmts {
		var table *priorBackupTable
		switch n := statementNode(stmt).(type) {
		case *ast.UpdateStmt:
			table = priorBackupTableOf(n.Relation, priorBackupStatementUpdate)
		case *ast.DeleteStmt:
			table = priorBackupTableOf(n.Relation, priorBackupStatementDelete)
		default:
		}
		if table == nil {
			continue
		}
		key := fmt.Sprintf("%s.%s", table.schema, table.table)
		if _, ok := types[key]; !ok {
			order = append(order, key)
		}
		types[key] = append(types[key], table.statementType)
	}

	for _, key := range order {
		if mixedDML(types[key]) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDMLOnTable, Table: key})
		}
	}

	return findings
}

// priorBackupTableOf returns the table a DML's target relation names, nil when
// the statement names no relation. The alias V1 kept is not read: V1 grouped
// by the relation the DML names.
func priorBackupTableOf(relation *ast.RangeVar, statementType priorBackupStatementType) *priorBackupTable {
	if relation == nil || relation.Relname == "" {
		return nil
	}
	schema := relation.Schemaname
	if schema == "" {
		schema = defaultSchema
	}
	return &priorBackupTable{
		schema:        schema,
		table:         relation.Relname,
		statementType: statementType,
	}
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

// isDDLStatement reports whether the V1 rule counted node as a DDL statement.
// The list is the one the V1 advisor kept, not every statement that changes a
// schema: a GRANT, a COMMENT ON ROLE, or a CREATE/ALTER/DROP ROLE is not in it.
func isDDLStatement(node ast.Node) bool {
	switch node.(type) {
	// CREATE statements
	case *ast.CreateStmt,
		*ast.CreateTableAsStmt,
		*ast.CreateDomainStmt,
		*ast.CreateExtensionStmt,
		*ast.CreateFdwStmt,
		*ast.CreateForeignServerStmt,
		*ast.CreateForeignTableStmt,
		*ast.CreateFunctionStmt,
		*ast.CreateRoleStmt,
		*ast.CreatePolicyStmt,
		*ast.CreatePLangStmt,
		*ast.CreateSchemaStmt,
		*ast.CreateSeqStmt,
		*ast.CreateSubscriptionStmt,
		*ast.CreateStatsStmt,
		*ast.CreateTableSpaceStmt,
		*ast.CreateTransformStmt,
		*ast.CreateTrigStmt,
		*ast.CreateEventTrigStmt,
		*ast.CreateUserMappingStmt,
		*ast.CreatedbStmt,
		*ast.CreateAmStmt,
		*ast.CreatePublicationStmt,
		*ast.CreateOpClassStmt,
		*ast.CreateOpFamilyStmt,
		*ast.CreateCastStmt,
		*ast.CreateConversionStmt,
		*ast.CreateRangeStmt,
		*ast.CompositeTypeStmt,
		*ast.CreateEnumStmt,
		*ast.IndexStmt,
		*ast.ViewStmt,
		// ALTER statements
		*ast.AlterEventTrigStmt,
		*ast.AlterCollationStmt,
		*ast.AlterDatabaseStmt,
		*ast.AlterDatabaseSetStmt,
		*ast.AlterDefaultPrivilegesStmt,
		*ast.AlterDomainStmt,
		*ast.AlterEnumStmt,
		*ast.AlterExtensionStmt,
		*ast.AlterExtensionContentsStmt,
		*ast.AlterFdwStmt,
		*ast.AlterForeignServerStmt,
		*ast.AlterFunctionStmt,
		*ast.AlterObjectDependsStmt,
		*ast.AlterObjectSchemaStmt,
		*ast.AlterOwnerStmt,
		*ast.AlterOperatorStmt,
		*ast.AlterTypeStmt,
		*ast.AlterPolicyStmt,
		*ast.AlterSeqStmt,
		*ast.AlterSystemStmt,
		*ast.AlterTableStmt,
		*ast.AlterTableSpaceOptionsStmt,
		*ast.AlterPublicationStmt,
		*ast.AlterRoleSetStmt,
		*ast.AlterRoleStmt,
		*ast.AlterSubscriptionStmt,
		*ast.AlterStatsStmt,
		*ast.AlterTSConfigurationStmt,
		*ast.AlterTSDictionaryStmt,
		*ast.AlterUserMappingStmt,
		*ast.AlterOpFamilyStmt,
		// DROP statements
		*ast.DropStmt,
		*ast.DropOwnedStmt,
		*ast.DropSubscriptionStmt,
		*ast.DropTableSpaceStmt,
		*ast.DropRoleStmt,
		*ast.DropUserMappingStmt,
		*ast.DropdbStmt,
		// Other DDL statements
		*ast.TruncateStmt,
		*ast.RenameStmt,
		*ast.CommentStmt,
		*ast.DefineStmt,
		*ast.ReindexStmt,
		*ast.ClusterStmt,
		*ast.RefreshMatViewStmt,
		*ast.RuleStmt,
		*ast.SecLabelStmt,
		*ast.ReassignOwnedStmt:
		return true
	default:
		return false
	}
}
