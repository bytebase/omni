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
	default:
		return "Unknown"
	}
}
