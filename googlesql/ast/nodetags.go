package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under googlesql/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable — do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// DCL — GRANT / REVOKE (parser-dcl node).
	//
	// These mirror the legacy ANTLR grant_statement / revoke_statement family
	// (GoogleSQLParser.g4): one statement node per verb plus the shared
	// privilege and grantee leaf nodes.
	T_GrantStmt
	T_RevokeStmt
	T_Privilege
	T_Grantee

	// Expressions (googlesql/expressions node).
	//
	// The full GoogleSQL expression tree: the precedence chain
	// (expression / expression_higher_prec_than_and / and_expression), every
	// primary form, and the operator / predicate / constructor families. These
	// mirror the legacy ANTLR expression rules (GoogleSQLParser.g4 §2.17-§2.18),
	// a hand-port of ZetaSQL.
	T_Identifier
	T_PathExpr
	T_TypeRef
	T_StarExpr
	T_Literal
	T_TypedLiteral
	T_IntervalExpr
	T_Parameter
	T_SystemVariable
	T_ParenExpr
	T_UnaryExpr
	T_BinaryExpr
	T_CompareExpr
	T_IsExpr
	T_InExpr
	T_BetweenExpr
	T_LikeExpr
	T_CaseExpr
	T_WhenClause
	T_CastExpr
	T_ExtractExpr
	T_FuncCall
	T_NamedArg
	T_LambdaExpr
	T_SequenceArg
	T_StarModifiers
	T_ArrayExpr
	T_StructExpr
	T_StructFieldExpr
	T_FieldAccess
	T_IndexAccess
	T_ExtensionAccess
	T_SubqueryExpr
	T_ExistsExpr
	T_ArraySubqueryExpr
	T_NewConstructor
	T_BracedConstructor
	T_ReplaceFieldsExpr
	T_WithExpr
	T_WindowSpec
	T_WindowFrame
	T_OrderItem
	T_HavingModifier
	T_ClampedModifier

	// Query / SELECT core (googlesql/parser-select node).
	//
	// The query stack: the QueryStmt wrapper (WITH + body + trailing ORDER BY/
	// LIMIT/FOR UPDATE), the SELECT block, set operations, CTEs, FROM sources
	// (table / subquery / TVF / UNNEST), joins, GROUP BY (incl. ROLLUP/CUBE/
	// GROUPING SETS), and named WINDOW definitions. These mirror the legacy
	// ANTLR query grammar (GoogleSQLParser.g4 §2.13-§2.16), a hand-port of
	// ZetaSQL.
	T_QueryStmt
	T_WithClause
	T_CTE
	T_RecursionDepth
	T_SelectStmt
	T_SelectItem
	T_SetOperation
	T_TableExpr
	T_UnnestExpr
	T_JoinExpr
	T_GroupByClause
	T_GroupingItem
	T_WindowDef

	// Core DDL (googlesql/parser-ddl node).
	//
	// CREATE / ALTER / DROP over TABLE / VIEW / INDEX / SCHEMA / DATABASE for the
	// BigQuery + Spanner union, plus their structural sub-nodes (column / option /
	// constraint / key-part). These mirror the legacy ANTLR DDL grammar
	// (GoogleSQLParser.g4 §2.2-§2.6), a hand-port of ZetaSQL.
	T_OptionsEntry
	T_KeyPart
	T_ForeignKeyRef
	T_ColumnDef
	T_TableConstraint
	T_InterleaveClause
	T_CreateTableStmt
	T_ViewColumn
	T_CreateViewStmt
	T_CreateIndexStmt
	T_CreateSchemaStmt
	T_CreateDatabaseStmt
	T_AlterStmt
	T_AlterAction
	T_DropStmt

	// BigQuery-specific DDL (parser-ddl-bigquery node): CREATE [AGGREGATE]
	// FUNCTION / TABLE FUNCTION / PROCEDURE, CREATE MATERIALIZED|APPROX VIEW,
	// CREATE [SEARCH|VECTOR] INDEX, CREATE SNAPSHOT TABLE, CREATE ROW ACCESS
	// POLICY, the generic-entity CREATE/ALTER/DROP path
	// (CAPACITY/RESERVATION/ASSIGNMENT), and their BigQuery-only ALTER/DROP forms.
	T_FunctionParam
	T_CreateFunctionStmt
	T_CreateProcedureStmt
	T_CreateMaterializedViewStmt
	T_CreateSnapshotStmt
	T_SearchVectorIndexStmt
	T_CreateRowAccessPolicyStmt
	T_CreateEntityStmt
	T_BQAlterStmt
	T_BQDropStmt
	T_DropAllRowAccessPoliciesStmt

	// Spanner-specific DDL (parser-ddl-spanner node): CHANGE STREAM, SEQUENCE,
	// ROLE, LOCALITY GROUP, and PROTO BUNDLE CREATE/ALTER/DROP forms. These have
	// no first-class rule in the legacy ANTLR grammar (it rides them on the
	// generic-entity hook or rejects them); the omni parser models them directly,
	// authoritatively verified against the live Cloud Spanner emulator. (Role-based
	// GRANT/REVOKE reuses T_GrantStmt/T_RevokeStmt with the GranteeRole kind.)
	T_ChangeStreamTrackedTable
	T_CreateChangeStreamStmt
	T_AlterChangeStreamStmt
	T_DropChangeStreamStmt
	T_CreateSequenceStmt
	T_AlterSequenceStmt
	T_DropSequenceStmt
	T_CreateRoleStmt
	T_DropRoleStmt
	T_CreateLocalityGroupStmt
	T_AlterLocalityGroupStmt
	T_DropLocalityGroupStmt
	T_CreateProtoBundleStmt
	T_AlterProtoBundleStmt

	// DML — INSERT / UPDATE / DELETE / MERGE / TRUNCATE (parser-dml node).
	//
	// The data-manipulation family for the BigQuery + Spanner union, plus their
	// structural sub-nodes (VALUES row, SET item, ON CONFLICT, MERGE WHEN/action,
	// THEN RETURN). These mirror the legacy ANTLR DML grammar
	// (GoogleSQLParser.g4 §2.7), a hand-port of ZetaSQL.
	T_DefaultExpr
	T_InsertStmt
	T_InsertRow
	T_InsertTable
	T_OnConflict
	T_UpdateStmt
	T_UpdateItem
	T_DeleteStmt
	T_MergeStmt
	T_MergeWhen
	T_MergeAction
	T_TruncateStmt
	T_Returning

	// Transactions / batch + utility statements (parser-utility node).
	//
	// The §2.9 transaction/batch family (BEGIN/START TRANSACTION/COMMIT/ROLLBACK,
	// START/RUN/ABORT BATCH) and the §2.10 + §2.3 utility statements (ASSERT,
	// ANALYZE, DESCRIBE, RENAME, CALL). These mirror the legacy ANTLR grammar
	// (GoogleSQLParser.g4), a hand-port of ZetaSQL.
	T_TransactionStmt
	T_TransactionMode
	T_BatchStmt
	T_AssertStmt
	T_AnalyzeStmt
	T_TableAndColumnInfo
	T_DescribeStmt
	T_RenameStmt
	T_CallArg
	T_CallStmt
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_GrantStmt:
		return "GrantStmt"
	case T_RevokeStmt:
		return "RevokeStmt"
	case T_Privilege:
		return "Privilege"
	case T_Grantee:
		return "Grantee"
	case T_Identifier:
		return "Identifier"
	case T_PathExpr:
		return "PathExpr"
	case T_TypeRef:
		return "TypeRef"
	case T_StarExpr:
		return "StarExpr"
	case T_Literal:
		return "Literal"
	case T_TypedLiteral:
		return "TypedLiteral"
	case T_IntervalExpr:
		return "IntervalExpr"
	case T_Parameter:
		return "Parameter"
	case T_SystemVariable:
		return "SystemVariable"
	case T_ParenExpr:
		return "ParenExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_CompareExpr:
		return "CompareExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_InExpr:
		return "InExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_WhenClause:
		return "WhenClause"
	case T_CastExpr:
		return "CastExpr"
	case T_ExtractExpr:
		return "ExtractExpr"
	case T_FuncCall:
		return "FuncCall"
	case T_NamedArg:
		return "NamedArg"
	case T_LambdaExpr:
		return "LambdaExpr"
	case T_SequenceArg:
		return "SequenceArg"
	case T_StarModifiers:
		return "StarModifiers"
	case T_ArrayExpr:
		return "ArrayExpr"
	case T_StructExpr:
		return "StructExpr"
	case T_StructFieldExpr:
		return "StructFieldExpr"
	case T_FieldAccess:
		return "FieldAccess"
	case T_IndexAccess:
		return "IndexAccess"
	case T_ExtensionAccess:
		return "ExtensionAccess"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
	case T_ArraySubqueryExpr:
		return "ArraySubqueryExpr"
	case T_NewConstructor:
		return "NewConstructor"
	case T_BracedConstructor:
		return "BracedConstructor"
	case T_ReplaceFieldsExpr:
		return "ReplaceFieldsExpr"
	case T_WithExpr:
		return "WithExpr"
	case T_WindowSpec:
		return "WindowSpec"
	case T_WindowFrame:
		return "WindowFrame"
	case T_OrderItem:
		return "OrderItem"
	case T_HavingModifier:
		return "HavingModifier"
	case T_ClampedModifier:
		return "ClampedModifier"
	case T_QueryStmt:
		return "QueryStmt"
	case T_WithClause:
		return "WithClause"
	case T_CTE:
		return "CTE"
	case T_RecursionDepth:
		return "RecursionDepth"
	case T_SelectStmt:
		return "SelectStmt"
	case T_SelectItem:
		return "SelectItem"
	case T_SetOperation:
		return "SetOperation"
	case T_TableExpr:
		return "TableExpr"
	case T_UnnestExpr:
		return "UnnestExpr"
	case T_JoinExpr:
		return "JoinExpr"
	case T_GroupByClause:
		return "GroupByClause"
	case T_GroupingItem:
		return "GroupingItem"
	case T_WindowDef:
		return "WindowDef"
	case T_OptionsEntry:
		return "OptionsEntry"
	case T_KeyPart:
		return "KeyPart"
	case T_ForeignKeyRef:
		return "ForeignKeyRef"
	case T_ColumnDef:
		return "ColumnDef"
	case T_TableConstraint:
		return "TableConstraint"
	case T_InterleaveClause:
		return "InterleaveClause"
	case T_CreateTableStmt:
		return "CreateTableStmt"
	case T_ViewColumn:
		return "ViewColumn"
	case T_CreateViewStmt:
		return "CreateViewStmt"
	case T_CreateIndexStmt:
		return "CreateIndexStmt"
	case T_CreateSchemaStmt:
		return "CreateSchemaStmt"
	case T_CreateDatabaseStmt:
		return "CreateDatabaseStmt"
	case T_AlterStmt:
		return "AlterStmt"
	case T_AlterAction:
		return "AlterAction"
	case T_DropStmt:
		return "DropStmt"
	case T_FunctionParam:
		return "FunctionParam"
	case T_CreateFunctionStmt:
		return "CreateFunctionStmt"
	case T_CreateProcedureStmt:
		return "CreateProcedureStmt"
	case T_CreateMaterializedViewStmt:
		return "CreateMaterializedViewStmt"
	case T_CreateSnapshotStmt:
		return "CreateSnapshotStmt"
	case T_SearchVectorIndexStmt:
		return "SearchVectorIndexStmt"
	case T_CreateRowAccessPolicyStmt:
		return "CreateRowAccessPolicyStmt"
	case T_CreateEntityStmt:
		return "CreateEntityStmt"
	case T_BQAlterStmt:
		return "BQAlterStmt"
	case T_BQDropStmt:
		return "BQDropStmt"
	case T_DropAllRowAccessPoliciesStmt:
		return "DropAllRowAccessPoliciesStmt"
	case T_ChangeStreamTrackedTable:
		return "ChangeStreamTrackedTable"
	case T_CreateChangeStreamStmt:
		return "CreateChangeStreamStmt"
	case T_AlterChangeStreamStmt:
		return "AlterChangeStreamStmt"
	case T_DropChangeStreamStmt:
		return "DropChangeStreamStmt"
	case T_CreateSequenceStmt:
		return "CreateSequenceStmt"
	case T_AlterSequenceStmt:
		return "AlterSequenceStmt"
	case T_DropSequenceStmt:
		return "DropSequenceStmt"
	case T_CreateRoleStmt:
		return "CreateRoleStmt"
	case T_DropRoleStmt:
		return "DropRoleStmt"
	case T_CreateLocalityGroupStmt:
		return "CreateLocalityGroupStmt"
	case T_AlterLocalityGroupStmt:
		return "AlterLocalityGroupStmt"
	case T_DropLocalityGroupStmt:
		return "DropLocalityGroupStmt"
	case T_CreateProtoBundleStmt:
		return "CreateProtoBundleStmt"
	case T_AlterProtoBundleStmt:
		return "AlterProtoBundleStmt"
	case T_DefaultExpr:
		return "DefaultExpr"
	case T_InsertStmt:
		return "InsertStmt"
	case T_InsertRow:
		return "InsertRow"
	case T_InsertTable:
		return "InsertTable"
	case T_OnConflict:
		return "OnConflict"
	case T_UpdateStmt:
		return "UpdateStmt"
	case T_UpdateItem:
		return "UpdateItem"
	case T_DeleteStmt:
		return "DeleteStmt"
	case T_MergeStmt:
		return "MergeStmt"
	case T_MergeWhen:
		return "MergeWhen"
	case T_MergeAction:
		return "MergeAction"
	case T_TruncateStmt:
		return "TruncateStmt"
	case T_Returning:
		return "Returning"
	case T_TransactionStmt:
		return "TransactionStmt"
	case T_TransactionMode:
		return "TransactionMode"
	case T_BatchStmt:
		return "BatchStmt"
	case T_AssertStmt:
		return "AssertStmt"
	case T_AnalyzeStmt:
		return "AnalyzeStmt"
	case T_TableAndColumnInfo:
		return "TableAndColumnInfo"
	case T_DescribeStmt:
		return "DescribeStmt"
	case T_RenameStmt:
		return "RenameStmt"
	case T_CallArg:
		return "CallArg"
	case T_CallStmt:
		return "CallStmt"
	default:
		return "Unknown"
	}
}
