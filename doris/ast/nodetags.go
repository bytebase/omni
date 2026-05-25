package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under doris/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable -- do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified multi-part name
	// (e.g., catalog.db.table).
	T_ObjectName

	// T_TypeName is the tag for *TypeName, a Doris data type as it appears in
	// SQL source (e.g., INT, VARCHAR(255), ARRAY<INT>, MAP<STRING,INT>).
	T_TypeName

	// T_CreateIndexStmt is the tag for *CreateIndexStmt.
	T_CreateIndexStmt

	// T_DropIndexStmt is the tag for *DropIndexStmt.
	T_DropIndexStmt

	// T_BuildIndexStmt is the tag for *BuildIndexStmt.
	T_BuildIndexStmt

	// T_CreateDatabaseStmt is the tag for *CreateDatabaseStmt.
	T_CreateDatabaseStmt

	// T_AlterDatabaseStmt is the tag for *AlterDatabaseStmt.
	T_AlterDatabaseStmt

	// T_DropDatabaseStmt is the tag for *DropDatabaseStmt.
	T_DropDatabaseStmt

	// T_Property is the tag for *Property, a key=value pair in PROPERTIES clauses.
	T_Property

	// Expression nodes (T1.3).

	// T_BinaryExpr is the tag for *BinaryExpr (left op right).
	T_BinaryExpr

	// T_UnaryExpr is the tag for *UnaryExpr (op expr).
	T_UnaryExpr

	// T_IsExpr is the tag for *IsExpr (expr IS [NOT] NULL/TRUE/FALSE).
	T_IsExpr

	// T_BetweenExpr is the tag for *BetweenExpr (expr [NOT] BETWEEN low AND high).
	T_BetweenExpr

	// T_InExpr is the tag for *InExpr (expr [NOT] IN (...)).
	T_InExpr

	// T_LikeExpr is the tag for *LikeExpr (expr [NOT] LIKE pattern [ESCAPE esc]).
	T_LikeExpr

	// T_RegexpExpr is the tag for *RegexpExpr (expr [NOT] REGEXP/RLIKE pattern).
	T_RegexpExpr

	// T_FuncCallExpr is the tag for *FuncCallExpr (name(args...)).
	T_FuncCallExpr

	// T_CastExpr is the tag for *CastExpr (CAST/TRY_CAST).
	T_CastExpr

	// T_CaseExpr is the tag for *CaseExpr (CASE...END).
	T_CaseExpr

	// T_WhenClause is the tag for *WhenClause (WHEN cond THEN result).
	T_WhenClause

	// T_SubqueryExpr is the tag for *SubqueryExpr (placeholder for subqueries).
	T_SubqueryExpr

	// T_ColumnRef is the tag for *ColumnRef (qualified column reference).
	T_ColumnRef

	// T_Literal is the tag for *Literal (numeric, string, boolean, NULL).
	T_Literal

	// T_ParenExpr is the tag for *ParenExpr ((expr)).
	T_ParenExpr

	// T_ExistsExpr is the tag for *ExistsExpr (EXISTS (subquery)).
	T_ExistsExpr

	// T_IntervalExpr is the tag for *IntervalExpr (INTERVAL expr unit).
	T_IntervalExpr

	// T_OrderByItem is the tag for *OrderByItem (expr ASC/DESC NULLS FIRST/LAST).
	T_OrderByItem

	// CTE / WITH clause nodes (T1.6).

	// T_WithClause is the tag for *WithClause.
	T_WithClause

	// T_CTE is the tag for *CTE (one named CTE entry inside a WITH clause).
	T_CTE

	// SELECT statement nodes (T1.4).

	// T_SelectStmt is the tag for *SelectStmt.
	T_SelectStmt

	// T_SelectItem is the tag for *SelectItem (one item in the SELECT list).
	T_SelectItem

	// T_TableRef is the tag for *TableRef (table reference in FROM clause).
	T_TableRef

	// T_JoinClause is the tag for *JoinClause (JOIN expression in FROM).
	T_JoinClause

	// T_SetOpStmt is the tag for *SetOpStmt (UNION/INTERSECT/EXCEPT/MINUS).
	T_SetOpStmt

	// DDL — CREATE TABLE nodes (T2.1).

	// T_CreateTableStmt is the tag for *CreateTableStmt.
	T_CreateTableStmt

	// T_ColumnDef is the tag for *ColumnDef.
	T_ColumnDef

	// T_IndexDef is the tag for *IndexDef (inline index in CREATE TABLE).
	T_IndexDef

	// T_TableConstraint is the tag for *TableConstraint (PRIMARY KEY, UNIQUE).
	T_TableConstraint

	// T_KeyDesc is the tag for *KeyDesc (AGGREGATE/UNIQUE/DUPLICATE KEY).
	T_KeyDesc

	// T_PartitionDesc is the tag for *PartitionDesc (PARTITION BY RANGE/LIST).
	T_PartitionDesc

	// T_PartitionItem is the tag for *PartitionItem.
	T_PartitionItem

	// T_DistributionDesc is the tag for *DistributionDesc (DISTRIBUTED BY HASH/RANDOM).
	T_DistributionDesc

	// T_RollupDef is the tag for *RollupDef.
	T_RollupDef

	// T_RawQuery is the tag for *RawQuery (placeholder for unparsed SELECT).
	T_RawQuery

	// DDL — ALTER TABLE nodes (T2.2).

	// T_AlterTableStmt is the tag for *AlterTableStmt.
	T_AlterTableStmt

	// T_AlterTableAction is the tag for *AlterTableAction.
	T_AlterTableAction

	// DDL — VIEW nodes (T2.4).

	// T_ViewColumn is the tag for *ViewColumn (one column in a view column list).
	T_ViewColumn

	// T_CreateViewStmt is the tag for *CreateViewStmt.
	T_CreateViewStmt

	// T_AlterViewStmt is the tag for *AlterViewStmt.
	T_AlterViewStmt

	// T_DropViewStmt is the tag for *DropViewStmt.
	T_DropViewStmt

	// DML — INSERT statement (T4.1).

	// T_InsertStmt is the tag for *InsertStmt.
	T_InsertStmt

	// DML — UPDATE / DELETE nodes (T4.2).

	// T_Assignment is the tag for *Assignment (col = expr in UPDATE SET clause).
	T_Assignment

	// T_UpdateStmt is the tag for *UpdateStmt.
	T_UpdateStmt

	// T_DeleteStmt is the tag for *DeleteStmt.
	T_DeleteStmt

	// DML — MERGE INTO nodes (T4.3).

	// T_MergeStmt is the tag for *MergeStmt.
	T_MergeStmt

	// T_MergeClause is the tag for *MergeClause (one WHEN clause inside MERGE).
	T_MergeClause

	// DML — TRUNCATE TABLE / COPY INTO / LOAD / EXPORT nodes (T6.1).

	// T_TruncateTableStmt is the tag for *TruncateTableStmt.
	T_TruncateTableStmt

	// T_CopyIntoStmt is the tag for *CopyIntoStmt.
	T_CopyIntoStmt

	// T_LoadDataDesc is the tag for *LoadDataDesc.
	T_LoadDataDesc

	// T_LoadDataStmt is the tag for *LoadDataStmt.
	T_LoadDataStmt

	// T_ExportStmt is the tag for *ExportStmt.
	T_ExportStmt
	// Streaming load management nodes (T6.2).

	// T_CreateRoutineLoadStmt is the tag for *CreateRoutineLoadStmt.
	T_CreateRoutineLoadStmt

	// T_AlterRoutineLoadStmt is the tag for *AlterRoutineLoadStmt.
	T_AlterRoutineLoadStmt

	// T_PauseRoutineLoadStmt is the tag for *PauseRoutineLoadStmt.
	T_PauseRoutineLoadStmt

	// T_ResumeRoutineLoadStmt is the tag for *ResumeRoutineLoadStmt.
	T_ResumeRoutineLoadStmt

	// T_StopRoutineLoadStmt is the tag for *StopRoutineLoadStmt.
	T_StopRoutineLoadStmt

	// T_ShowRoutineLoadStmt is the tag for *ShowRoutineLoadStmt.
	T_ShowRoutineLoadStmt

	// T_ShowRoutineLoadTaskStmt is the tag for *ShowRoutineLoadTaskStmt.
	T_ShowRoutineLoadTaskStmt

	// T_SyncStmt is the tag for *SyncStmt.
	T_SyncStmt

	// DDL — Materialized View nodes (T5.1).

	// T_MTMVRefreshTrigger is the tag for *MTMVRefreshTrigger.
	T_MTMVRefreshTrigger

	// T_CreateMTMVStmt is the tag for *CreateMTMVStmt.
	T_CreateMTMVStmt

	// T_AlterMTMVStmt is the tag for *AlterMTMVStmt.
	T_AlterMTMVStmt

	// T_DropMTMVStmt is the tag for *DropMTMVStmt.
	T_DropMTMVStmt

	// T_RefreshMTMVStmt is the tag for *RefreshMTMVStmt.
	T_RefreshMTMVStmt

	// T_PauseMTMVJobStmt is the tag for *PauseMTMVJobStmt.
	T_PauseMTMVJobStmt

	// T_ResumeMTMVJobStmt is the tag for *ResumeMTMVJobStmt.
	T_ResumeMTMVJobStmt

	// T_CancelMTMVTaskStmt is the tag for *CancelMTMVTaskStmt.
	T_CancelMTMVTaskStmt

	// DDL — CATALOG nodes (T5.2).

	// T_CreateCatalogStmt is the tag for *CreateCatalogStmt.
	T_CreateCatalogStmt

	// T_AlterCatalogStmt is the tag for *AlterCatalogStmt.
	T_AlterCatalogStmt

	// T_DropCatalogStmt is the tag for *DropCatalogStmt.
	T_DropCatalogStmt

	// T_RefreshCatalogStmt is the tag for *RefreshCatalogStmt.
	T_RefreshCatalogStmt
	// DDL — STORAGE VAULT nodes (T5.3).

	// T_CreateStorageVaultStmt is the tag for *CreateStorageVaultStmt.
	T_CreateStorageVaultStmt

	// T_AlterStorageVaultStmt is the tag for *AlterStorageVaultStmt.
	T_AlterStorageVaultStmt

	// T_DropStorageVaultStmt is the tag for *DropStorageVaultStmt.
	T_DropStorageVaultStmt

	// T_SetDefaultStorageVaultStmt is the tag for *SetDefaultStorageVaultStmt.
	T_SetDefaultStorageVaultStmt

	// T_UnsetDefaultStorageVaultStmt is the tag for *UnsetDefaultStorageVaultStmt.
	T_UnsetDefaultStorageVaultStmt

	// DDL — STORAGE POLICY nodes (T5.3).

	// T_CreateStoragePolicyStmt is the tag for *CreateStoragePolicyStmt.
	T_CreateStoragePolicyStmt

	// T_AlterStoragePolicyStmt is the tag for *AlterStoragePolicyStmt.
	T_AlterStoragePolicyStmt

	// T_DropStoragePolicyStmt is the tag for *DropStoragePolicyStmt.
	T_DropStoragePolicyStmt

	// DDL — REPOSITORY nodes (T5.3).

	// T_CreateRepositoryStmt is the tag for *CreateRepositoryStmt.
	T_CreateRepositoryStmt

	// T_AlterRepositoryStmt is the tag for *AlterRepositoryStmt.
	T_AlterRepositoryStmt

	// T_DropRepositoryStmt is the tag for *DropRepositoryStmt.
	T_DropRepositoryStmt

	// DDL — STAGE nodes (T5.3).

	// T_CreateStageStmt is the tag for *CreateStageStmt.
	T_CreateStageStmt

	// T_DropStageStmt is the tag for *DropStageStmt.
	T_DropStageStmt

	// DDL — FILE nodes (T5.3).

	// T_CreateFileStmt is the tag for *CreateFileStmt.
	T_CreateFileStmt

	// T_DropFileStmt is the tag for *DropFileStmt.
	T_DropFileStmt
	// Workload management DDL nodes (T5.4).

	// T_CreateWorkloadGroupStmt is the tag for *CreateWorkloadGroupStmt.
	T_CreateWorkloadGroupStmt

	// T_AlterWorkloadGroupStmt is the tag for *AlterWorkloadGroupStmt.
	T_AlterWorkloadGroupStmt

	// T_DropWorkloadGroupStmt is the tag for *DropWorkloadGroupStmt.
	T_DropWorkloadGroupStmt

	// T_WorkloadPolicyItem is the tag for *WorkloadPolicyItem.
	T_WorkloadPolicyItem

	// T_CreateWorkloadPolicyStmt is the tag for *CreateWorkloadPolicyStmt.
	T_CreateWorkloadPolicyStmt

	// T_AlterWorkloadPolicyStmt is the tag for *AlterWorkloadPolicyStmt.
	T_AlterWorkloadPolicyStmt

	// T_DropWorkloadPolicyStmt is the tag for *DropWorkloadPolicyStmt.
	T_DropWorkloadPolicyStmt

	// T_CreateResourceStmt is the tag for *CreateResourceStmt.
	T_CreateResourceStmt

	// T_AlterResourceStmt is the tag for *AlterResourceStmt.
	T_AlterResourceStmt

	// T_DropResourceStmt is the tag for *DropResourceStmt.
	T_DropResourceStmt

	// T_CreateSQLBlockRuleStmt is the tag for *CreateSQLBlockRuleStmt.
	T_CreateSQLBlockRuleStmt

	// T_AlterSQLBlockRuleStmt is the tag for *AlterSQLBlockRuleStmt.
	T_AlterSQLBlockRuleStmt

	// T_DropSQLBlockRuleStmt is the tag for *DropSQLBlockRuleStmt.
	T_DropSQLBlockRuleStmt
	// Security DDL nodes (T5.5).

	// T_CreateRowPolicyStmt is the tag for *CreateRowPolicyStmt.
	T_CreateRowPolicyStmt

	// T_DropRowPolicyStmt is the tag for *DropRowPolicyStmt.
	T_DropRowPolicyStmt

	// T_CreateEncryptKeyStmt is the tag for *CreateEncryptKeyStmt.
	T_CreateEncryptKeyStmt

	// T_DropEncryptKeyStmt is the tag for *DropEncryptKeyStmt.
	T_DropEncryptKeyStmt

	// T_DictionaryColumn is the tag for *DictionaryColumn.
	T_DictionaryColumn

	// T_CreateDictionaryStmt is the tag for *CreateDictionaryStmt.
	T_CreateDictionaryStmt

	// T_AlterDictionaryStmt is the tag for *AlterDictionaryStmt.
	T_AlterDictionaryStmt

	// T_DropDictionaryStmt is the tag for *DropDictionaryStmt.
	T_DropDictionaryStmt

	// T_RefreshDictionaryStmt is the tag for *RefreshDictionaryStmt.
	T_RefreshDictionaryStmt

	// T_CreateRoleStmt is the tag for *CreateRoleStmt.
	T_CreateRoleStmt

	// T_AlterRoleStmt is the tag for *AlterRoleStmt.
	T_AlterRoleStmt

	// T_DropRoleStmt is the tag for *DropRoleStmt.
	T_DropRoleStmt

	// T_UserIdentity is the tag for *UserIdentity ('user'@'host').
	T_UserIdentity

	// T_CreateUserStmt is the tag for *CreateUserStmt.
	T_CreateUserStmt

	// T_AlterUserStmt is the tag for *AlterUserStmt.
	T_AlterUserStmt

	// T_DropUserStmt is the tag for *DropUserStmt.
	T_DropUserStmt

	// T_SetPasswordStmt is the tag for *SetPasswordStmt.
	T_SetPasswordStmt

	// DCL nodes (T7.1).

	// T_GrantStmt is the tag for *GrantStmt.
	T_GrantStmt

	// T_RevokeStmt is the tag for *RevokeStmt.
	T_RevokeStmt
	// TCL — Transaction Control Language nodes (T7.2).

	// T_BeginStmt is the tag for *BeginStmt.
	T_BeginStmt

	// T_CommitStmt is the tag for *CommitStmt.
	T_CommitStmt

	// T_RollbackStmt is the tag for *RollbackStmt.
	T_RollbackStmt
	// Utility statement nodes (T7.3).

	// T_ShowStmt is the tag for *ShowStmt.
	T_ShowStmt

	// T_DescribeStmt is the tag for *DescribeStmt.
	T_DescribeStmt

	// T_ExplainStmt is the tag for *ExplainStmt.
	T_ExplainStmt

	// T_UseStmt is the tag for *UseStmt.
	T_UseStmt

	// T_SetStmt is the tag for *SetStmt (generic variable assignment).
	T_SetStmt

	// T_SetItem is the tag for *SetItem (one item in a SET statement).
	T_SetItem

	// T_UnsetStmt is the tag for *UnsetStmt.
	T_UnsetStmt

	// T_HelpStmt is the tag for *HelpStmt.
	T_HelpStmt
	// Admin / System cluster-management nodes (T7.4).

	// T_AdminStmt is the tag for *AdminStmt (ADMIN ... variants).
	T_AdminStmt

	// T_SystemAlterStmt is the tag for *SystemAlterStmt (ALTER SYSTEM ...).
	T_SystemAlterStmt

	// T_CancelDecommissionStmt is the tag for *CancelDecommissionStmt.
	T_CancelDecommissionStmt
	// Utility / admin command nodes (T7.5).

	// T_BackupStmt is the tag for *BackupStmt.
	T_BackupStmt

	// T_RestoreStmt is the tag for *RestoreStmt.
	T_RestoreStmt

	// T_KillStmt is the tag for *KillStmt.
	T_KillStmt

	// T_LockTablesStmt is the tag for *LockTablesStmt.
	T_LockTablesStmt

	// T_LockItem is the tag for *LockItem.
	T_LockItem

	// T_UnlockTablesStmt is the tag for *UnlockTablesStmt.
	T_UnlockTablesStmt

	// T_InstallPluginStmt is the tag for *InstallPluginStmt.
	T_InstallPluginStmt

	// T_UninstallPluginStmt is the tag for *UninstallPluginStmt.
	T_UninstallPluginStmt

	// T_WarmUpStmt is the tag for *WarmUpStmt.
	T_WarmUpStmt

	// T_CleanStmt is the tag for *CleanStmt.
	T_CleanStmt

	// T_CancelStmt is the tag for *CancelStmt (generic, non-MTMV).
	T_CancelStmt

	// T_RecoverStmt is the tag for *RecoverStmt.
	T_RecoverStmt
	// Job DDL nodes (T8.1).

	// T_JobSchedule is the tag for *JobSchedule.
	T_JobSchedule

	// T_CreateJobStmt is the tag for *CreateJobStmt.
	T_CreateJobStmt

	// T_AlterJobStmt is the tag for *AlterJobStmt.
	T_AlterJobStmt

	// T_DropJobStmt is the tag for *DropJobStmt.
	T_DropJobStmt

	// T_PauseJobStmt is the tag for *PauseJobStmt.
	T_PauseJobStmt

	// T_ResumeJobStmt is the tag for *ResumeJobStmt.
	T_ResumeJobStmt

	// T_CancelTaskStmt is the tag for *CancelTaskStmt.
	T_CancelTaskStmt

	// T_ShowJobStmt is the tag for *ShowJobStmt.
	T_ShowJobStmt

	// T_ShowJobTaskStmt is the tag for *ShowJobTaskStmt.
	T_ShowJobTaskStmt
	// Statistics / Analyze / Constraint management nodes (T8.3).

	// T_AddConstraintStmt is the tag for *AddConstraintStmt.
	T_AddConstraintStmt

	// T_DropConstraintStmt is the tag for *DropConstraintStmt.
	T_DropConstraintStmt

	// T_ShowConstraintsStmt is the tag for *ShowConstraintsStmt.
	T_ShowConstraintsStmt

	// T_AnalyzeStmt is the tag for *AnalyzeStmt.
	T_AnalyzeStmt

	// T_ShowAnalyzeStmt is the tag for *ShowAnalyzeStmt.
	T_ShowAnalyzeStmt

	// T_ShowStatsStmt is the tag for *ShowStatsStmt.
	T_ShowStatsStmt

	// T_DropStatsStmt is the tag for *DropStatsStmt.
	T_DropStatsStmt

	// T_KillAnalyzeStmt is the tag for *KillAnalyzeStmt.
	T_KillAnalyzeStmt
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
	case T_CreateIndexStmt:
		return "CreateIndexStmt"
	case T_DropIndexStmt:
		return "DropIndexStmt"
	case T_BuildIndexStmt:
		return "BuildIndexStmt"
	case T_CreateDatabaseStmt:
		return "CreateDatabaseStmt"
	case T_AlterDatabaseStmt:
		return "AlterDatabaseStmt"
	case T_DropDatabaseStmt:
		return "DropDatabaseStmt"
	case T_Property:
		return "Property"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_InExpr:
		return "InExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_RegexpExpr:
		return "RegexpExpr"
	case T_FuncCallExpr:
		return "FuncCallExpr"
	case T_CastExpr:
		return "CastExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_WhenClause:
		return "WhenClause"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ColumnRef:
		return "ColumnRef"
	case T_Literal:
		return "Literal"
	case T_ParenExpr:
		return "ParenExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
	case T_IntervalExpr:
		return "IntervalExpr"
	case T_OrderByItem:
		return "OrderByItem"
	case T_WithClause:
		return "WithClause"
	case T_CTE:
		return "CTE"
	case T_SelectStmt:
		return "SelectStmt"
	case T_SelectItem:
		return "SelectItem"
	case T_TableRef:
		return "TableRef"
	case T_JoinClause:
		return "JoinClause"
	case T_SetOpStmt:
		return "SetOpStmt"
	case T_CreateTableStmt:
		return "CreateTableStmt"
	case T_ColumnDef:
		return "ColumnDef"
	case T_IndexDef:
		return "IndexDef"
	case T_TableConstraint:
		return "TableConstraint"
	case T_KeyDesc:
		return "KeyDesc"
	case T_PartitionDesc:
		return "PartitionDesc"
	case T_PartitionItem:
		return "PartitionItem"
	case T_DistributionDesc:
		return "DistributionDesc"
	case T_RollupDef:
		return "RollupDef"
	case T_RawQuery:
		return "RawQuery"
	case T_AlterTableStmt:
		return "AlterTableStmt"
	case T_AlterTableAction:
		return "AlterTableAction"
	case T_ViewColumn:
		return "ViewColumn"
	case T_CreateViewStmt:
		return "CreateViewStmt"
	case T_AlterViewStmt:
		return "AlterViewStmt"
	case T_DropViewStmt:
		return "DropViewStmt"
	case T_InsertStmt:
		return "InsertStmt"
	case T_Assignment:
		return "Assignment"
	case T_UpdateStmt:
		return "UpdateStmt"
	case T_DeleteStmt:
		return "DeleteStmt"
	case T_MergeStmt:
		return "MergeStmt"
	case T_MergeClause:
		return "MergeClause"
	case T_TruncateTableStmt:
		return "TruncateTableStmt"
	case T_CopyIntoStmt:
		return "CopyIntoStmt"
	case T_LoadDataDesc:
		return "LoadDataDesc"
	case T_LoadDataStmt:
		return "LoadDataStmt"
	case T_ExportStmt:
		return "ExportStmt"
	case T_CreateRoutineLoadStmt:
		return "CreateRoutineLoadStmt"
	case T_AlterRoutineLoadStmt:
		return "AlterRoutineLoadStmt"
	case T_PauseRoutineLoadStmt:
		return "PauseRoutineLoadStmt"
	case T_ResumeRoutineLoadStmt:
		return "ResumeRoutineLoadStmt"
	case T_StopRoutineLoadStmt:
		return "StopRoutineLoadStmt"
	case T_ShowRoutineLoadStmt:
		return "ShowRoutineLoadStmt"
	case T_ShowRoutineLoadTaskStmt:
		return "ShowRoutineLoadTaskStmt"
	case T_SyncStmt:
		return "SyncStmt"
	case T_MTMVRefreshTrigger:
		return "MTMVRefreshTrigger"
	case T_CreateMTMVStmt:
		return "CreateMTMVStmt"
	case T_AlterMTMVStmt:
		return "AlterMTMVStmt"
	case T_DropMTMVStmt:
		return "DropMTMVStmt"
	case T_RefreshMTMVStmt:
		return "RefreshMTMVStmt"
	case T_PauseMTMVJobStmt:
		return "PauseMTMVJobStmt"
	case T_ResumeMTMVJobStmt:
		return "ResumeMTMVJobStmt"
	case T_CancelMTMVTaskStmt:
		return "CancelMTMVTaskStmt"
	case T_CreateCatalogStmt:
		return "CreateCatalogStmt"
	case T_AlterCatalogStmt:
		return "AlterCatalogStmt"
	case T_DropCatalogStmt:
		return "DropCatalogStmt"
	case T_RefreshCatalogStmt:
		return "RefreshCatalogStmt"
	case T_CreateStorageVaultStmt:
		return "CreateStorageVaultStmt"
	case T_AlterStorageVaultStmt:
		return "AlterStorageVaultStmt"
	case T_DropStorageVaultStmt:
		return "DropStorageVaultStmt"
	case T_SetDefaultStorageVaultStmt:
		return "SetDefaultStorageVaultStmt"
	case T_UnsetDefaultStorageVaultStmt:
		return "UnsetDefaultStorageVaultStmt"
	case T_CreateStoragePolicyStmt:
		return "CreateStoragePolicyStmt"
	case T_AlterStoragePolicyStmt:
		return "AlterStoragePolicyStmt"
	case T_DropStoragePolicyStmt:
		return "DropStoragePolicyStmt"
	case T_CreateRepositoryStmt:
		return "CreateRepositoryStmt"
	case T_AlterRepositoryStmt:
		return "AlterRepositoryStmt"
	case T_DropRepositoryStmt:
		return "DropRepositoryStmt"
	case T_CreateStageStmt:
		return "CreateStageStmt"
	case T_DropStageStmt:
		return "DropStageStmt"
	case T_CreateFileStmt:
		return "CreateFileStmt"
	case T_DropFileStmt:
		return "DropFileStmt"
	case T_CreateWorkloadGroupStmt:
		return "CreateWorkloadGroupStmt"
	case T_AlterWorkloadGroupStmt:
		return "AlterWorkloadGroupStmt"
	case T_DropWorkloadGroupStmt:
		return "DropWorkloadGroupStmt"
	case T_WorkloadPolicyItem:
		return "WorkloadPolicyItem"
	case T_CreateWorkloadPolicyStmt:
		return "CreateWorkloadPolicyStmt"
	case T_AlterWorkloadPolicyStmt:
		return "AlterWorkloadPolicyStmt"
	case T_DropWorkloadPolicyStmt:
		return "DropWorkloadPolicyStmt"
	case T_CreateResourceStmt:
		return "CreateResourceStmt"
	case T_AlterResourceStmt:
		return "AlterResourceStmt"
	case T_DropResourceStmt:
		return "DropResourceStmt"
	case T_CreateSQLBlockRuleStmt:
		return "CreateSQLBlockRuleStmt"
	case T_AlterSQLBlockRuleStmt:
		return "AlterSQLBlockRuleStmt"
	case T_DropSQLBlockRuleStmt:
		return "DropSQLBlockRuleStmt"
	case T_CreateRowPolicyStmt:
		return "CreateRowPolicyStmt"
	case T_DropRowPolicyStmt:
		return "DropRowPolicyStmt"
	case T_CreateEncryptKeyStmt:
		return "CreateEncryptKeyStmt"
	case T_DropEncryptKeyStmt:
		return "DropEncryptKeyStmt"
	case T_DictionaryColumn:
		return "DictionaryColumn"
	case T_CreateDictionaryStmt:
		return "CreateDictionaryStmt"
	case T_AlterDictionaryStmt:
		return "AlterDictionaryStmt"
	case T_DropDictionaryStmt:
		return "DropDictionaryStmt"
	case T_RefreshDictionaryStmt:
		return "RefreshDictionaryStmt"
	case T_CreateRoleStmt:
		return "CreateRoleStmt"
	case T_AlterRoleStmt:
		return "AlterRoleStmt"
	case T_DropRoleStmt:
		return "DropRoleStmt"
	case T_UserIdentity:
		return "UserIdentity"
	case T_CreateUserStmt:
		return "CreateUserStmt"
	case T_AlterUserStmt:
		return "AlterUserStmt"
	case T_DropUserStmt:
		return "DropUserStmt"
	case T_SetPasswordStmt:
		return "SetPasswordStmt"
	case T_GrantStmt:
		return "GrantStmt"
	case T_RevokeStmt:
		return "RevokeStmt"
	case T_BeginStmt:
		return "BeginStmt"
	case T_CommitStmt:
		return "CommitStmt"
	case T_RollbackStmt:
		return "RollbackStmt"
	case T_ShowStmt:
		return "ShowStmt"
	case T_DescribeStmt:
		return "DescribeStmt"
	case T_ExplainStmt:
		return "ExplainStmt"
	case T_UseStmt:
		return "UseStmt"
	case T_SetStmt:
		return "SetStmt"
	case T_SetItem:
		return "SetItem"
	case T_UnsetStmt:
		return "UnsetStmt"
	case T_HelpStmt:
		return "HelpStmt"
	case T_AdminStmt:
		return "AdminStmt"
	case T_SystemAlterStmt:
		return "SystemAlterStmt"
	case T_CancelDecommissionStmt:
		return "CancelDecommissionStmt"
	case T_BackupStmt:
		return "BackupStmt"
	case T_RestoreStmt:
		return "RestoreStmt"
	case T_KillStmt:
		return "KillStmt"
	case T_LockTablesStmt:
		return "LockTablesStmt"
	case T_LockItem:
		return "LockItem"
	case T_UnlockTablesStmt:
		return "UnlockTablesStmt"
	case T_InstallPluginStmt:
		return "InstallPluginStmt"
	case T_UninstallPluginStmt:
		return "UninstallPluginStmt"
	case T_WarmUpStmt:
		return "WarmUpStmt"
	case T_CleanStmt:
		return "CleanStmt"
	case T_CancelStmt:
		return "CancelStmt"
	case T_RecoverStmt:
		return "RecoverStmt"
	case T_JobSchedule:
		return "JobSchedule"
	case T_CreateJobStmt:
		return "CreateJobStmt"
	case T_AlterJobStmt:
		return "AlterJobStmt"
	case T_DropJobStmt:
		return "DropJobStmt"
	case T_PauseJobStmt:
		return "PauseJobStmt"
	case T_ResumeJobStmt:
		return "ResumeJobStmt"
	case T_CancelTaskStmt:
		return "CancelTaskStmt"
	case T_ShowJobStmt:
		return "ShowJobStmt"
	case T_ShowJobTaskStmt:
		return "ShowJobTaskStmt"
	case T_AddConstraintStmt:
		return "AddConstraintStmt"
	case T_DropConstraintStmt:
		return "DropConstraintStmt"
	case T_ShowConstraintsStmt:
		return "ShowConstraintsStmt"
	case T_AnalyzeStmt:
		return "AnalyzeStmt"
	case T_ShowAnalyzeStmt:
		return "ShowAnalyzeStmt"
	case T_ShowStatsStmt:
		return "ShowStatsStmt"
	case T_DropStatsStmt:
		return "DropStatsStmt"
	case T_KillAnalyzeStmt:
		return "KillAnalyzeStmt"
	default:
		return "Unknown"
	}
}
