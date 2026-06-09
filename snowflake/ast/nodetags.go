package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under snowflake/ast must declare a unique
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
	// returned by F4 (parser-entry).
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified 1/2/3-part name.
	T_ObjectName

	// T_TypeName is the tag for *TypeName, a data type reference.
	T_TypeName

	T_Literal
	T_ColumnRef
	T_StarExpr
	T_DollarRef
	T_BinaryExpr
	T_UnaryExpr
	T_ParenExpr
	T_CastExpr
	T_CaseExpr
	T_FuncCallExpr
	T_IffExpr
	T_CollateExpr
	T_IsExpr
	T_BetweenExpr
	T_InExpr
	T_LikeExpr
	T_AccessExpr
	T_ArrayLiteralExpr
	T_JsonLiteralExpr
	T_LambdaExpr
	T_SubqueryExpr
	T_ExistsExpr
	T_SelectStmt
	T_SetOperationStmt
	T_TableRef
	T_JoinExpr
	T_CreateTableStmt
	T_ColumnDef
	T_TableConstraint
	T_CreateDatabaseStmt
	T_AlterDatabaseStmt
	T_DropDatabaseStmt
	T_UndropDatabaseStmt
	T_CreateSchemaStmt
	T_AlterSchemaStmt
	T_DropSchemaStmt
	T_UndropSchemaStmt
	T_DropStmt
	T_UndropStmt

	// DML statement tags (T5.1)
	T_InsertStmt
	T_InsertMultiStmt
	T_UpdateStmt
	T_DeleteStmt
	T_MergeStmt

	// VIEW + MATERIALIZED VIEW tags (T2.4)
	T_CreateViewStmt
	T_CreateMaterializedViewStmt
	T_AlterViewStmt
	T_AlterMaterializedViewStmt

	// ALTER TABLE tag (T2.3)
	T_AlterTableStmt

	// DCL: GRANT / REVOKE tags (T6.1)
	T_GrantStmt
	T_RevokeStmt
	T_Grantee
	T_GrantTarget

	// Utility / introspection tags (T6.3)
	T_ShowStmt
	T_DescribeStmt
	T_UseStmt
	T_SetStmt
	T_UnsetStmt
	T_CommentStmt
	T_TruncateStmt

	// Bulk data movement tags (T5.2)
	T_StageLocation
	T_CopyOption
	T_CopyIntoTableStmt
	T_CopyIntoLocationStmt
	T_PutStmt
	T_GetStmt
	T_ListStmt
	T_RemoveStmt

	// Stage DDL tags (T4.1)
	T_CreateStageStmt
	T_AlterStageStmt

	// File-format DDL tags (T4.2)
	T_CreateFileFormatStmt
	T_AlterFileFormatStmt

	// Integration-object DDL tags (T4.7)
	T_ResourceMonitorTrigger
	T_ConnectionReplica
	T_CreateIntegrationStmt
	T_AlterIntegrationStmt

	// Replication & sharing DDL tags (T4.8)
	T_GroupOption
	T_CreateReplicationGroupStmt
	T_AlterReplicationGroupStmt
	T_CreateAccountStmt
	T_AlterAccountStmt
	T_CreateShareStmt
	T_AlterShareStmt

	// Tag / semantic-view / dataset DDL tags (T4.9)
	T_CreateTagStmt
	T_AlterTagStmt
	T_SemanticViewSection
	T_CreateSemanticViewStmt
	T_AlterSemanticViewStmt
	T_CreateDatasetStmt

	// Data-pipeline DDL tags (T4.3)
	T_CreatePipeStmt
	T_AlterPipeStmt
	T_CreateStreamStmt
	T_AlterStreamStmt
	T_CreateTaskStmt
	T_AlterTaskStmt
	T_CreateAlertStmt
	T_AlterAlertStmt

	// Routine DDL tags (T4.5)
	T_CreateRoutineStmt
	T_AlterRoutineStmt

	// TCL (transaction control) tags (T6.2)
	T_BeginStmt
	T_CommitStmt
	T_RollbackStmt

	// Table-variant + sequence DDL tags (T4.4)
	T_CreateDynamicTableStmt
	T_AlterDynamicTableStmt
	T_ExternalColumnDef
	T_CreateExternalTableStmt
	T_AlterExternalTableStmt
	T_CreateEventTableStmt
	T_CreateSequenceStmt
	T_AlterSequenceStmt

	// Warehouse DDL tags (gap-warehouse)
	T_CreateWarehouseStmt
	T_AlterWarehouseStmt

	// Network rule DDL tags (gap-network-rule)
	T_CreateNetworkRuleStmt
	T_AlterNetworkRuleStmt

	// Access-control DDL tags (T4.6)
	T_PolicyArg
	T_CreateRoleStmt
	T_CreateUserStmt
	T_CreatePolicyStmt
	T_AlterRoleStmt
	T_AlterUserStmt
	T_AlterPolicyStmt

	// CALL / EXECUTE IMMEDIATE / EXECUTE TASK / EXPLAIN tags (T5.4)
	T_CallArg
	T_CallStmt
	T_ExecuteImmediateStmt
	T_ExecuteTaskStmt
	T_ExplainStmt

	// Table-attached query clause tags (T5.3)
	T_PivotClause
	T_PivotInClause
	T_PivotValue
	T_UnpivotClause
	T_UnpivotColumn
	T_MatchRecognizeClause
	T_MatchMeasure
	T_RowsPerMatch
	T_AfterMatchSkip
	T_MatchDefine
	T_RowPattern
	T_SampleClause
	T_TimeTravelClause
	T_ChangesClause

	// Snowflake Scripting (T7.1) — statement nodes plus the structural
	// sub-nodes (declarations, IF/CASE branches, exception handlers). The
	// sub-nodes are Nodes (unlike the CaseExpr/WhenClause precedent) so the
	// generated walker descends into their nested bodies and queries, which an
	// analysis pass needs to reach.
	T_ScriptDeclaration
	T_ScriptExceptionHandler
	T_ScriptIfBranch
	T_ScriptCaseWhen
	T_ScriptBlockStmt
	T_ScriptAssignStmt
	T_ScriptLetStmt
	T_ScriptIfStmt
	T_ScriptCaseStmt
	T_ScriptForStmt
	T_ScriptWhileStmt
	T_ScriptRepeatStmt
	T_ScriptLoopStmt
	T_ScriptBreakStmt
	T_ScriptContinueStmt
	T_ScriptReturnStmt
	T_ScriptOpenStmt
	T_ScriptFetchStmt
	T_ScriptCloseStmt
	T_ScriptRaiseStmt
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_ObjectName:
		return "ObjectName"
	case T_TypeName:
		return "TypeName"
	case T_Literal:
		return "Literal"
	case T_ColumnRef:
		return "ColumnRef"
	case T_StarExpr:
		return "StarExpr"
	case T_DollarRef:
		return "DollarRef"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_ParenExpr:
		return "ParenExpr"
	case T_CastExpr:
		return "CastExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_FuncCallExpr:
		return "FuncCallExpr"
	case T_IffExpr:
		return "IffExpr"
	case T_CollateExpr:
		return "CollateExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_InExpr:
		return "InExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_AccessExpr:
		return "AccessExpr"
	case T_ArrayLiteralExpr:
		return "ArrayLiteralExpr"
	case T_JsonLiteralExpr:
		return "JsonLiteralExpr"
	case T_LambdaExpr:
		return "LambdaExpr"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
	case T_SelectStmt:
		return "SelectStmt"
	case T_SetOperationStmt:
		return "SetOperationStmt"
	case T_TableRef:
		return "TableRef"
	case T_JoinExpr:
		return "JoinExpr"
	case T_CreateTableStmt:
		return "CreateTableStmt"
	case T_ColumnDef:
		return "ColumnDef"
	case T_TableConstraint:
		return "TableConstraint"
	case T_CreateDatabaseStmt:
		return "CreateDatabaseStmt"
	case T_AlterDatabaseStmt:
		return "AlterDatabaseStmt"
	case T_DropDatabaseStmt:
		return "DropDatabaseStmt"
	case T_UndropDatabaseStmt:
		return "UndropDatabaseStmt"
	case T_CreateSchemaStmt:
		return "CreateSchemaStmt"
	case T_AlterSchemaStmt:
		return "AlterSchemaStmt"
	case T_DropSchemaStmt:
		return "DropSchemaStmt"
	case T_UndropSchemaStmt:
		return "UndropSchemaStmt"
	case T_DropStmt:
		return "DropStmt"
	case T_UndropStmt:
		return "UndropStmt"
	case T_InsertStmt:
		return "InsertStmt"
	case T_InsertMultiStmt:
		return "InsertMultiStmt"
	case T_UpdateStmt:
		return "UpdateStmt"
	case T_DeleteStmt:
		return "DeleteStmt"
	case T_MergeStmt:
		return "MergeStmt"
	case T_CreateViewStmt:
		return "CreateViewStmt"
	case T_CreateMaterializedViewStmt:
		return "CreateMaterializedViewStmt"
	case T_AlterViewStmt:
		return "AlterViewStmt"
	case T_AlterMaterializedViewStmt:
		return "AlterMaterializedViewStmt"
	case T_AlterTableStmt:
		return "AlterTableStmt"
	case T_GrantStmt:
		return "GrantStmt"
	case T_RevokeStmt:
		return "RevokeStmt"
	case T_Grantee:
		return "Grantee"
	case T_GrantTarget:
		return "GrantTarget"
	case T_ShowStmt:
		return "ShowStmt"
	case T_DescribeStmt:
		return "DescribeStmt"
	case T_UseStmt:
		return "UseStmt"
	case T_SetStmt:
		return "SetStmt"
	case T_UnsetStmt:
		return "UnsetStmt"
	case T_CommentStmt:
		return "CommentStmt"
	case T_TruncateStmt:
		return "TruncateStmt"
	case T_StageLocation:
		return "StageLocation"
	case T_CopyOption:
		return "CopyOption"
	case T_CopyIntoTableStmt:
		return "CopyIntoTableStmt"
	case T_CopyIntoLocationStmt:
		return "CopyIntoLocationStmt"
	case T_PutStmt:
		return "PutStmt"
	case T_GetStmt:
		return "GetStmt"
	case T_ListStmt:
		return "ListStmt"
	case T_RemoveStmt:
		return "RemoveStmt"
	case T_CreateStageStmt:
		return "CreateStageStmt"
	case T_AlterStageStmt:
		return "AlterStageStmt"
	case T_CreateFileFormatStmt:
		return "CreateFileFormatStmt"
	case T_AlterFileFormatStmt:
		return "AlterFileFormatStmt"
	case T_CreatePipeStmt:
		return "CreatePipeStmt"
	case T_AlterPipeStmt:
		return "AlterPipeStmt"
	case T_CreateStreamStmt:
		return "CreateStreamStmt"
	case T_AlterStreamStmt:
		return "AlterStreamStmt"
	case T_CreateTaskStmt:
		return "CreateTaskStmt"
	case T_AlterTaskStmt:
		return "AlterTaskStmt"
	case T_CreateAlertStmt:
		return "CreateAlertStmt"
	case T_AlterAlertStmt:
		return "AlterAlertStmt"
	case T_ResourceMonitorTrigger:
		return "ResourceMonitorTrigger"
	case T_ConnectionReplica:
		return "ConnectionReplica"
	case T_CreateIntegrationStmt:
		return "CreateIntegrationStmt"
	case T_AlterIntegrationStmt:
		return "AlterIntegrationStmt"
	case T_GroupOption:
		return "GroupOption"
	case T_CreateReplicationGroupStmt:
		return "CreateReplicationGroupStmt"
	case T_AlterReplicationGroupStmt:
		return "AlterReplicationGroupStmt"
	case T_CreateAccountStmt:
		return "CreateAccountStmt"
	case T_AlterAccountStmt:
		return "AlterAccountStmt"
	case T_CreateShareStmt:
		return "CreateShareStmt"
	case T_AlterShareStmt:
		return "AlterShareStmt"
	case T_CreateRoutineStmt:
		return "CreateRoutineStmt"
	case T_AlterRoutineStmt:
		return "AlterRoutineStmt"
	case T_BeginStmt:
		return "BeginStmt"
	case T_CommitStmt:
		return "CommitStmt"
	case T_RollbackStmt:
		return "RollbackStmt"
	case T_CreateDynamicTableStmt:
		return "CreateDynamicTableStmt"
	case T_AlterDynamicTableStmt:
		return "AlterDynamicTableStmt"
	case T_ExternalColumnDef:
		return "ExternalColumnDef"
	case T_CreateExternalTableStmt:
		return "CreateExternalTableStmt"
	case T_AlterExternalTableStmt:
		return "AlterExternalTableStmt"
	case T_CreateEventTableStmt:
		return "CreateEventTableStmt"
	case T_CreateSequenceStmt:
		return "CreateSequenceStmt"
	case T_AlterSequenceStmt:
		return "AlterSequenceStmt"
	case T_CreateWarehouseStmt:
		return "CreateWarehouseStmt"
	case T_AlterWarehouseStmt:
		return "AlterWarehouseStmt"
	case T_CreateNetworkRuleStmt:
		return "CreateNetworkRuleStmt"
	case T_AlterNetworkRuleStmt:
		return "AlterNetworkRuleStmt"
	case T_PolicyArg:
		return "PolicyArg"
	case T_CreateRoleStmt:
		return "CreateRoleStmt"
	case T_CreateUserStmt:
		return "CreateUserStmt"
	case T_CreatePolicyStmt:
		return "CreatePolicyStmt"
	case T_AlterRoleStmt:
		return "AlterRoleStmt"
	case T_AlterUserStmt:
		return "AlterUserStmt"
	case T_AlterPolicyStmt:
		return "AlterPolicyStmt"
	case T_CreateTagStmt:
		return "CreateTagStmt"
	case T_AlterTagStmt:
		return "AlterTagStmt"
	case T_SemanticViewSection:
		return "SemanticViewSection"
	case T_CreateSemanticViewStmt:
		return "CreateSemanticViewStmt"
	case T_AlterSemanticViewStmt:
		return "AlterSemanticViewStmt"
	case T_CreateDatasetStmt:
		return "CreateDatasetStmt"
	case T_CallArg:
		return "CallArg"
	case T_CallStmt:
		return "CallStmt"
	case T_ExecuteImmediateStmt:
		return "ExecuteImmediateStmt"
	case T_ExecuteTaskStmt:
		return "ExecuteTaskStmt"
	case T_ExplainStmt:
		return "ExplainStmt"
	case T_PivotClause:
		return "PivotClause"
	case T_PivotInClause:
		return "PivotInClause"
	case T_PivotValue:
		return "PivotValue"
	case T_UnpivotClause:
		return "UnpivotClause"
	case T_UnpivotColumn:
		return "UnpivotColumn"
	case T_MatchRecognizeClause:
		return "MatchRecognizeClause"
	case T_MatchMeasure:
		return "MatchMeasure"
	case T_RowsPerMatch:
		return "RowsPerMatch"
	case T_AfterMatchSkip:
		return "AfterMatchSkip"
	case T_MatchDefine:
		return "MatchDefine"
	case T_RowPattern:
		return "RowPattern"
	case T_SampleClause:
		return "SampleClause"
	case T_TimeTravelClause:
		return "TimeTravelClause"
	case T_ChangesClause:
		return "ChangesClause"
	case T_ScriptDeclaration:
		return "ScriptDeclaration"
	case T_ScriptExceptionHandler:
		return "ScriptExceptionHandler"
	case T_ScriptIfBranch:
		return "ScriptIfBranch"
	case T_ScriptCaseWhen:
		return "ScriptCaseWhen"
	case T_ScriptBlockStmt:
		return "ScriptBlockStmt"
	case T_ScriptAssignStmt:
		return "ScriptAssignStmt"
	case T_ScriptLetStmt:
		return "ScriptLetStmt"
	case T_ScriptIfStmt:
		return "ScriptIfStmt"
	case T_ScriptCaseStmt:
		return "ScriptCaseStmt"
	case T_ScriptForStmt:
		return "ScriptForStmt"
	case T_ScriptWhileStmt:
		return "ScriptWhileStmt"
	case T_ScriptRepeatStmt:
		return "ScriptRepeatStmt"
	case T_ScriptLoopStmt:
		return "ScriptLoopStmt"
	case T_ScriptBreakStmt:
		return "ScriptBreakStmt"
	case T_ScriptContinueStmt:
		return "ScriptContinueStmt"
	case T_ScriptReturnStmt:
		return "ScriptReturnStmt"
	case T_ScriptOpenStmt:
		return "ScriptOpenStmt"
	case T_ScriptFetchStmt:
		return "ScriptFetchStmt"
	case T_ScriptCloseStmt:
		return "ScriptCloseStmt"
	case T_ScriptRaiseStmt:
		return "ScriptRaiseStmt"
	default:
		return "Unknown"
	}
}
