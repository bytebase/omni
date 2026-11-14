package analysis

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
)

// classifyQueryType classifies an omni AST node into a QueryType.
// allSystems indicates whether all referenced tables are system/info_schema tables.
func classifyQueryType(node ast.Node, allSystems bool) (queryType QueryType, isExplainAnalyze bool) {
	if node == nil {
		return QueryTypeUnknown, false
	}

	switch n := node.(type) {
	// DML statements
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt, *ast.CopyStmt:
		return DML, false

	// SELECT: check for INTO clause (becomes DDL)
	case *ast.SelectStmt:
		if hasIntoClause(n) {
			return DDL, false
		}
		// A data-modifying CTE writes, so the SELECT needs DML permission.
		if containsWriteCTE(n) {
			return DML, false
		}
		if allSystems {
			return SelectInfoSchema, false
		}
		return Select, false

	// SET is treated as safe (Select)
	case *ast.VariableSetStmt:
		return Select, false

	// SHOW → info schema
	case *ast.VariableShowStmt:
		return SelectInfoSchema, false

	// EXPLAIN: check for ANALYZE option
	case *ast.ExplainStmt:
		if isExplainAnalyzeStmt(n) {
			qt := classifyExplainedQuery(n.Query)
			return qt, true
		}
		return Explain, false

	// REFRESH MATERIALIZED VIEW is DML
	case *ast.RefreshMatViewStmt:
		return DML, false

	// CALL stored procedure is DML
	case *ast.CallStmt:
		return DML, false

	// All DDL statements
	case *ast.CreateStmt, *ast.DropStmt, *ast.AlterTableStmt,
		*ast.CreateTableAsStmt, *ast.CreateSeqStmt, *ast.CreateSchemaStmt,
		*ast.CreatedbStmt, *ast.CreateFunctionStmt, *ast.CreateRoleStmt,
		*ast.IndexStmt, *ast.CreateExtensionStmt, *ast.CreateTrigStmt,
		*ast.CreateEventTrigStmt, *ast.CreateDomainStmt, *ast.CreateConversionStmt,
		*ast.CreateCastStmt, *ast.CreateOpClassStmt, *ast.CreateOpFamilyStmt,
		*ast.CreatePolicyStmt, *ast.CreateAmStmt, *ast.CreateTransformStmt,
		*ast.CreateStatsStmt, *ast.CreateTableSpaceStmt,
		*ast.CreateFdwStmt, *ast.CreateForeignServerStmt,
		*ast.CreateForeignTableStmt, *ast.CreatePLangStmt,
		*ast.CreatePublicationStmt, *ast.CreateSubscriptionStmt,
		*ast.CreateUserMappingStmt,
		*ast.ViewStmt,
		*ast.AlterSeqStmt, *ast.AlterDatabaseStmt, *ast.AlterDatabaseSetStmt,
		*ast.AlterFunctionStmt, *ast.AlterRoleStmt, *ast.AlterRoleSetStmt,
		*ast.AlterCollationStmt, *ast.AlterDomainStmt,
		*ast.AlterExtensionStmt, *ast.AlterExtensionContentsStmt,
		*ast.AlterFdwStmt, *ast.AlterForeignServerStmt,
		*ast.AlterOpFamilyStmt, *ast.AlterPolicyStmt,
		*ast.AlterEventTrigStmt, *ast.AlterObjectDependsStmt,
		*ast.AlterObjectSchemaStmt, *ast.AlterOwnerStmt,
		*ast.AlterOperatorStmt, *ast.AlterTypeStmt, *ast.AlterEnumStmt,
		*ast.AlterStatsStmt, *ast.AlterTableSpaceOptionsStmt,
		*ast.AlterSystemStmt, *ast.AlterPublicationStmt,
		*ast.AlterSubscriptionStmt, *ast.AlterUserMappingStmt,
		*ast.CompositeTypeStmt, *ast.AlterDefaultPrivilegesStmt,
		*ast.AlterTSConfigurationStmt, *ast.AlterTSDictionaryStmt,
		*ast.DropdbStmt, *ast.DropRoleStmt, *ast.DropOwnedStmt,
		*ast.DropSubscriptionStmt, *ast.DropUserMappingStmt,
		*ast.TruncateStmt, *ast.CommentStmt,
		*ast.GrantStmt, *ast.GrantRoleStmt,
		*ast.ClusterStmt, *ast.VacuumStmt, *ast.LockStmt,
		*ast.ReindexStmt, *ast.RuleStmt,
		*ast.RenameStmt, *ast.ReassignOwnedStmt,
		*ast.SecLabelStmt,
		*ast.DoStmt, *ast.DiscardStmt,
		*ast.FetchStmt, *ast.ConstraintsSetStmt, *ast.CheckPointStmt,
		*ast.CreateEnumStmt:
		return DDL, false

	default:
		return QueryTypeUnknown, false
	}
}

// hasIntoClause checks if a SelectStmt has an INTO clause (SELECT INTO).
func hasIntoClause(n *ast.SelectStmt) bool {
	return intoClause(n) != nil
}

// intoClause returns the INTO clause of n, searching set-operation arms
// (UNION/INTERSECT/EXCEPT) — the parser attaches INTO to the first arm, not the root.
func intoClause(n *ast.SelectStmt) *ast.IntoClause {
	if n == nil {
		return nil
	}
	if n.IntoClause != nil {
		return n.IntoClause
	}
	if c := intoClause(n.Larg); c != nil {
		return c
	}
	return intoClause(n.Rarg)
}

// isExplainAnalyzeStmt reports whether an ExplainStmt executes its query: its last ANALYZE option
// is not FALSE, OFF, or 0, the values PostgreSQL reads as false.
func isExplainAnalyzeStmt(n *ast.ExplainStmt) bool {
	if n.Options == nil {
		return false
	}
	analyze := false
	for _, item := range n.Options.Items {
		de, ok := item.(*ast.DefElem)
		if !ok || !strings.EqualFold(de.Defname, "analyze") {
			continue
		}
		switch arg := de.Arg.(type) {
		case *ast.String:
			analyze = !strings.EqualFold(arg.Str, "false") && !strings.EqualFold(arg.Str, "off")
		case *ast.Integer:
			analyze = arg.Ival != 0
		default:
			analyze = true
		}
	}
	return analyze
}

// classifyExplainedQuery returns the QueryType for the query inside EXPLAIN
// ANALYZE, which actually runs it — so unlike classifyQueryType's own default,
// an unrecognized node type here must fail closed rather than pass as Select.
// The cases below are exhaustive over what PostgreSQL actually accepts as an
// EXPLAIN target (CALL, TRUNCATE, COPY, DDL like DROP/ALTER do not parse
// there), matching classifyQueryType's DML/DDL classification for the same
// node types; the default exists only for a future grammar addition.
func classifyExplainedQuery(query ast.Node) QueryType {
	if query == nil {
		return Select
	}
	switch n := query.(type) {
	case *ast.SelectStmt:
		if hasIntoClause(n) {
			return DDL
		}
		if containsWriteCTE(n) {
			return DML
		}
		return Select
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		return DML
	case *ast.DeclareCursorStmt:
		return Select
	case *ast.CreateTableAsStmt:
		return DDL
	case *ast.RefreshMatViewStmt:
		return DML
	case *ast.ExecuteStmt:
		return Select
	default:
		return QueryTypeUnknown
	}
}

// containsWriteCTE reports whether any common table expression in the tree is
// data-modifying. Fail-safe by construction — a whitelist, not a denylist: a CTE
// term is read-only ONLY if it is an *ast.SelectStmt, which is how pg models
// SELECT, VALUES, TABLE, set operations and WITH RECURSIVE. Anything else —
// INSERT/UPDATE/DELETE/MERGE today, or any write node added to the grammar
// tomorrow — counts as a write, so the next forgotten write type fails closed.
func containsWriteCTE(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if cte, ok := n.(*ast.CommonTableExpr); ok {
			if _, isSelect := cte.Ctequery.(*ast.SelectStmt); !isSelect {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
