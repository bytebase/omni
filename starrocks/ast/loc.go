package ast

// NodeLoc returns the source location of n, or NoLoc() if n is nil or its
// concrete type carries no Loc field. Every concrete node type added under
// doris/ast must add a case here.
//
// The pattern matches snowflake/ast.NodeLoc.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	case *File:
		return v.Loc
	case *ObjectName:
		return v.Loc
	case *TypeName:
		return v.Loc
	case *CreateIndexStmt:
		return v.Loc
	case *DropIndexStmt:
		return v.Loc
	case *BuildIndexStmt:
		return v.Loc
	case *CreateDatabaseStmt:
		return v.Loc
	case *AlterDatabaseStmt:
		return v.Loc
	case *DropDatabaseStmt:
		return v.Loc
	case *DropTableStmt:
		return v.Loc
	case *CreateFunctionStmt:
		return v.Loc
	case *DropFunctionStmt:
		return v.Loc
	case *Property:
		return v.Loc
	case *BinaryExpr:
		return v.Loc
	case *UnaryExpr:
		return v.Loc
	case *CollateExpr:
		return v.Loc
	case *MatchExpr:
		return v.Loc
	case *DereferenceExpr:
		return v.Loc
	case *RowExpr:
		return v.Loc
	case *IsExpr:
		return v.Loc
	case *BetweenExpr:
		return v.Loc
	case *InExpr:
		return v.Loc
	case *LikeExpr:
		return v.Loc
	case *RegexpExpr:
		return v.Loc
	case *FuncCallExpr:
		return v.Loc
	case *CastExpr:
		return v.Loc
	case *ExtractExpr:
		return v.Loc
	case *VariableRef:
		return v.Loc
	case *LambdaExpr:
		return v.Loc
	case *CaseExpr:
		return v.Loc
	case *WhenClause:
		return v.Loc
	case *SubqueryExpr:
		return v.Loc
	case *ColumnRef:
		return v.Loc
	case *Literal:
		return v.Loc
	case *MapLiteral:
		return v.Loc
	case *MapEntry:
		return v.Loc
	case *ElementAtExpr:
		return v.Loc
	case *GroupingSetsExpr:
		return v.Loc
	case *ArrayLiteral:
		return v.Loc
	case *ParenExpr:
		return v.Loc
	case *ExistsExpr:
		return v.Loc
	case *IntervalExpr:
		return v.Loc
	case *OrderByItem:
		return v.Loc
	case *WithClause:
		return v.Loc
	case *CTE:
		return v.Loc
	case *SelectStmt:
		return v.Loc
	case *SelectItem:
		return v.Loc
	case *TableRef:
		return v.Loc
	case *InlineTable:
		return v.Loc
	case *TableFunctionRef:
		return v.Loc
	case *JoinClause:
		return v.Loc
	case *SetOpStmt:
		return v.Loc
	case *ParenSelect:
		return v.Loc
	case *IntoOutfileClause:
		return v.Loc
	case *PivotClause:
		return v.Loc
	case *PivotAggregation:
		return v.Loc
	case *PivotValue:
		return v.Loc
	case *CreateTableStmt:
		return v.Loc
	case *ColumnDef:
		return v.Loc
	case *IndexDef:
		return v.Loc
	case *TableConstraint:
		return v.Loc
	case *KeyDesc:
		return v.Loc
	case *PartitionDesc:
		return v.Loc
	case *PartitionItem:
		return v.Loc
	case *DistributionDesc:
		return v.Loc
	case *RollupDef:
		return v.Loc
	case *RawQuery:
		return v.Loc
	case *AlterTableStmt:
		return v.Loc
	case *AlterTableAction:
		return v.Loc
	case *ViewColumn:
		return v.Loc
	case *CreateViewStmt:
		return v.Loc
	case *AlterViewStmt:
		return v.Loc
	case *DropViewStmt:
		return v.Loc
	case *InsertStmt:
		return v.Loc
	case *Assignment:
		return v.Loc
	case *UpdateStmt:
		return v.Loc
	case *DeleteStmt:
		return v.Loc
	case *DataCacheSelectStmt:
		return v.Loc
	case *TruncateTableStmt:
		return v.Loc
	case *LoadDataDesc:
		return v.Loc
	case *LoadDataStmt:
		return v.Loc
	case *AlterLoadStmt:
		return v.Loc
	case *ExportStmt:
		return v.Loc
	case *CreateRoutineLoadStmt:
		return v.Loc
	case *AlterRoutineLoadStmt:
		return v.Loc
	case *PauseRoutineLoadStmt:
		return v.Loc
	case *ResumeRoutineLoadStmt:
		return v.Loc
	case *StopRoutineLoadStmt:
		return v.Loc
	case *ShowRoutineLoadStmt:
		return v.Loc
	case *ShowRoutineLoadTaskStmt:
		return v.Loc
	case *SyncStmt:
		return v.Loc
	case *MTMVRefreshTrigger:
		return v.Loc
	case *CreateMTMVStmt:
		return v.Loc
	case *AlterMTMVStmt:
		return v.Loc
	case *DropMTMVStmt:
		return v.Loc
	case *RefreshMTMVStmt:
		return v.Loc
	case *PauseMTMVJobStmt:
		return v.Loc
	case *ResumeMTMVJobStmt:
		return v.Loc
	case *CancelMTMVTaskStmt:
		return v.Loc
	case *CreateCatalogStmt:
		return v.Loc
	case *AlterCatalogStmt:
		return v.Loc
	case *DropCatalogStmt:
		return v.Loc
	case *RefreshCatalogStmt:
		return v.Loc
	case *CreateStorageVolumeStmt:
		return v.Loc
	case *AlterStorageVolumeStmt:
		return v.Loc
	case *AlterStorageVolumeClause:
		return v.Loc
	case *DropStorageVolumeStmt:
		return v.Loc
	case *DescStorageVolumeStmt:
		return v.Loc
	case *CreateWarehouseStmt:
		return v.Loc
	case *DropWarehouseStmt:
		return v.Loc
	case *SuspendWarehouseStmt:
		return v.Loc
	case *ResumeWarehouseStmt:
		return v.Loc
	case *CreateStoragePolicyStmt:
		return v.Loc
	case *AlterStoragePolicyStmt:
		return v.Loc
	case *DropStoragePolicyStmt:
		return v.Loc
	case *CreateRepositoryStmt:
		return v.Loc
	case *AlterRepositoryStmt:
		return v.Loc
	case *DropRepositoryStmt:
		return v.Loc
	case *CreateFileStmt:
		return v.Loc
	case *DropFileStmt:
		return v.Loc
	case *CreateResourceStmt:
		return v.Loc
	case *AlterResourceStmt:
		return v.Loc
	case *DropResourceStmt:
		return v.Loc
	case *CreateRowPolicyStmt:
		return v.Loc
	case *DropRowPolicyStmt:
		return v.Loc
	case *CreateEncryptKeyStmt:
		return v.Loc
	case *DropEncryptKeyStmt:
		return v.Loc
	case *DictionaryColumn:
		return v.Loc
	case *CreateDictionaryStmt:
		return v.Loc
	case *AlterDictionaryStmt:
		return v.Loc
	case *DropDictionaryStmt:
		return v.Loc
	case *RefreshDictionaryStmt:
		return v.Loc
	case *CreateRoleStmt:
		return v.Loc
	case *AlterRoleStmt:
		return v.Loc
	case *DropRoleStmt:
		return v.Loc
	case *UserIdentity:
		return v.Loc
	case *CreateUserStmt:
		return v.Loc
	case *AlterUserStmt:
		return v.Loc
	case *DropUserStmt:
		return v.Loc
	case *SetPasswordStmt:
		return v.Loc
	case *GrantStmt:
		return v.Loc
	case *RevokeStmt:
		return v.Loc
	case *BeginStmt:
		return v.Loc
	case *CommitStmt:
		return v.Loc
	case *RollbackStmt:
		return v.Loc
	case *ShowStmt:
		return v.Loc
	case *DescribeStmt:
		return v.Loc
	case *ExplainStmt:
		return v.Loc
	case *OptimizerTraceStmt:
		return v.Loc
	case *UseStmt:
		return v.Loc
	case *SetStmt:
		return v.Loc
	case *SetItem:
		return v.Loc
	case *UnsetStmt:
		return v.Loc
	case *HelpStmt:
		return v.Loc
	case *AdminStmt:
		return v.Loc
	case *SystemAlterStmt:
		return v.Loc
	case *CancelDecommissionStmt:
		return v.Loc
	case *BackupStmt:
		return v.Loc
	case *RestoreStmt:
		return v.Loc
	case *KillStmt:
		return v.Loc
	case *LockTablesStmt:
		return v.Loc
	case *LockItem:
		return v.Loc
	case *UnlockTablesStmt:
		return v.Loc
	case *InstallPluginStmt:
		return v.Loc
	case *UninstallPluginStmt:
		return v.Loc
	case *WarmUpStmt:
		return v.Loc
	case *CleanStmt:
		return v.Loc
	case *CancelStmt:
		return v.Loc
	case *RecoverStmt:
		return v.Loc

	// Job DDL nodes (T8.1).
	case *JobSchedule:
		return v.Loc
	case *CreateJobStmt:
		return v.Loc
	case *AlterJobStmt:
		return v.Loc
	case *DropJobStmt:
		return v.Loc
	case *PauseJobStmt:
		return v.Loc
	case *ResumeJobStmt:
		return v.Loc
	case *CancelTaskStmt:
		return v.Loc
	case *ShowJobStmt:
		return v.Loc
	case *ShowJobTaskStmt:
		return v.Loc
	case *TaskSchedule:
		return v.Loc
	case *SubmitTaskStmt:
		return v.Loc
	case *DropTaskStmt:
		return v.Loc
	case *ShowConstraintsStmt:
		return v.Loc
	case *AnalyzeStmt:
		return v.Loc
	case *ShowAnalyzeStmt:
		return v.Loc
	case *ShowStatsStmt:
		return v.Loc
	case *DropStatsStmt:
		return v.Loc
	case *KillAnalyzeStmt:
		return v.Loc
	case *ProcedureParam:
		return v.Loc
	case *CreateProcedureStmt:
		return v.Loc
	case *DropProcedureStmt:
		return v.Loc
	case *RefreshExternalTableStmt:
		return v.Loc
	case *ExternalCatalogRef:
		return v.Loc
	case *CreateResourceGroupStmt:
		return v.Loc
	case *DropResourceGroupStmt:
		return v.Loc
	case *AlterResourceGroupStmt:
		return v.Loc
	case *CreateDataCacheRuleStmt:
		return v.Loc
	case *DropDataCacheRuleStmt:
		return v.Loc
	case *ClearDataCacheRulesStmt:
		return v.Loc
	case *CreatePipeStmt:
		return v.Loc
	case *DropPipeStmt:
		return v.Loc
	case *AlterPipeStmt:
		return v.Loc
	case *DescPipeStmt:
		return v.Loc
	case *AlterPlanAdvisorStmt:
		return v.Loc
	case *TruncatePlanAdvisorStmt:
		return v.Loc
	case *TranslateStmt:
		return v.Loc
	case *PrepareStmt:
		return v.Loc
	case *ExecuteStmt:
		return v.Loc
	case *DeallocateStmt:
		return v.Loc
	case *ExecuteAsStmt:
		return v.Loc
	case *BlacklistStmt:
		return v.Loc
	default:
		return NoLoc()
	}
}

// SetNodeLoc sets n's Loc, the counterpart of NodeLoc. The parser uses it for a
// statement it read at other bytes and now reports at its own: a PREPARE
// string's body, whose decoded text does not stand in the statement byte for
// byte. Like NodeLoc it knows every node that carries one.
func SetNodeLoc(n Node, loc Loc) {
	switch v := n.(type) {
	case *File:
		v.Loc = loc
	case *ObjectName:
		v.Loc = loc
	case *TypeName:
		v.Loc = loc
	case *CreateIndexStmt:
		v.Loc = loc
	case *DropIndexStmt:
		v.Loc = loc
	case *BuildIndexStmt:
		v.Loc = loc
	case *CreateDatabaseStmt:
		v.Loc = loc
	case *AlterDatabaseStmt:
		v.Loc = loc
	case *DropDatabaseStmt:
		v.Loc = loc
	case *DropTableStmt:
		v.Loc = loc
	case *CreateFunctionStmt:
		v.Loc = loc
	case *DropFunctionStmt:
		v.Loc = loc
	case *Property:
		v.Loc = loc
	case *BinaryExpr:
		v.Loc = loc
	case *UnaryExpr:
		v.Loc = loc
	case *CollateExpr:
		v.Loc = loc
	case *MatchExpr:
		v.Loc = loc
	case *DereferenceExpr:
		v.Loc = loc
	case *RowExpr:
		v.Loc = loc
	case *IsExpr:
		v.Loc = loc
	case *BetweenExpr:
		v.Loc = loc
	case *InExpr:
		v.Loc = loc
	case *LikeExpr:
		v.Loc = loc
	case *RegexpExpr:
		v.Loc = loc
	case *FuncCallExpr:
		v.Loc = loc
	case *CastExpr:
		v.Loc = loc
	case *ExtractExpr:
		v.Loc = loc
	case *VariableRef:
		v.Loc = loc
	case *LambdaExpr:
		v.Loc = loc
	case *CaseExpr:
		v.Loc = loc
	case *WhenClause:
		v.Loc = loc
	case *SubqueryExpr:
		v.Loc = loc
	case *ColumnRef:
		v.Loc = loc
	case *Literal:
		v.Loc = loc
	case *MapLiteral:
		v.Loc = loc
	case *MapEntry:
		v.Loc = loc
	case *ElementAtExpr:
		v.Loc = loc
	case *GroupingSetsExpr:
		v.Loc = loc
	case *ArrayLiteral:
		v.Loc = loc
	case *ParenExpr:
		v.Loc = loc
	case *ExistsExpr:
		v.Loc = loc
	case *IntervalExpr:
		v.Loc = loc
	case *OrderByItem:
		v.Loc = loc
	case *WithClause:
		v.Loc = loc
	case *CTE:
		v.Loc = loc
	case *SelectStmt:
		v.Loc = loc
	case *SelectItem:
		v.Loc = loc
	case *TableRef:
		v.Loc = loc
	case *InlineTable:
		v.Loc = loc
	case *TableFunctionRef:
		v.Loc = loc
	case *JoinClause:
		v.Loc = loc
	case *SetOpStmt:
		v.Loc = loc
	case *ParenSelect:
		v.Loc = loc
	case *IntoOutfileClause:
		v.Loc = loc
	case *PivotClause:
		v.Loc = loc
	case *PivotAggregation:
		v.Loc = loc
	case *PivotValue:
		v.Loc = loc
	case *CreateTableStmt:
		v.Loc = loc
	case *ColumnDef:
		v.Loc = loc
	case *IndexDef:
		v.Loc = loc
	case *TableConstraint:
		v.Loc = loc
	case *KeyDesc:
		v.Loc = loc
	case *PartitionDesc:
		v.Loc = loc
	case *PartitionItem:
		v.Loc = loc
	case *DistributionDesc:
		v.Loc = loc
	case *RollupDef:
		v.Loc = loc
	case *RawQuery:
		v.Loc = loc
	case *AlterTableStmt:
		v.Loc = loc
	case *AlterTableAction:
		v.Loc = loc
	case *ViewColumn:
		v.Loc = loc
	case *CreateViewStmt:
		v.Loc = loc
	case *AlterViewStmt:
		v.Loc = loc
	case *DropViewStmt:
		v.Loc = loc
	case *InsertStmt:
		v.Loc = loc
	case *Assignment:
		v.Loc = loc
	case *UpdateStmt:
		v.Loc = loc
	case *DeleteStmt:
		v.Loc = loc
	case *DataCacheSelectStmt:
		v.Loc = loc
	case *TruncateTableStmt:
		v.Loc = loc
	case *LoadDataDesc:
		v.Loc = loc
	case *LoadDataStmt:
		v.Loc = loc
	case *AlterLoadStmt:
		v.Loc = loc
	case *ExportStmt:
		v.Loc = loc
	case *CreateRoutineLoadStmt:
		v.Loc = loc
	case *AlterRoutineLoadStmt:
		v.Loc = loc
	case *PauseRoutineLoadStmt:
		v.Loc = loc
	case *ResumeRoutineLoadStmt:
		v.Loc = loc
	case *StopRoutineLoadStmt:
		v.Loc = loc
	case *ShowRoutineLoadStmt:
		v.Loc = loc
	case *ShowRoutineLoadTaskStmt:
		v.Loc = loc
	case *SyncStmt:
		v.Loc = loc
	case *MTMVRefreshTrigger:
		v.Loc = loc
	case *CreateMTMVStmt:
		v.Loc = loc
	case *AlterMTMVStmt:
		v.Loc = loc
	case *DropMTMVStmt:
		v.Loc = loc
	case *RefreshMTMVStmt:
		v.Loc = loc
	case *PauseMTMVJobStmt:
		v.Loc = loc
	case *ResumeMTMVJobStmt:
		v.Loc = loc
	case *CancelMTMVTaskStmt:
		v.Loc = loc
	case *CreateCatalogStmt:
		v.Loc = loc
	case *AlterCatalogStmt:
		v.Loc = loc
	case *DropCatalogStmt:
		v.Loc = loc
	case *RefreshCatalogStmt:
		v.Loc = loc
	case *CreateStorageVolumeStmt:
		v.Loc = loc
	case *AlterStorageVolumeStmt:
		v.Loc = loc
	case *AlterStorageVolumeClause:
		v.Loc = loc
	case *DropStorageVolumeStmt:
		v.Loc = loc
	case *DescStorageVolumeStmt:
		v.Loc = loc
	case *CreateWarehouseStmt:
		v.Loc = loc
	case *DropWarehouseStmt:
		v.Loc = loc
	case *SuspendWarehouseStmt:
		v.Loc = loc
	case *ResumeWarehouseStmt:
		v.Loc = loc
	case *CreateStoragePolicyStmt:
		v.Loc = loc
	case *AlterStoragePolicyStmt:
		v.Loc = loc
	case *DropStoragePolicyStmt:
		v.Loc = loc
	case *CreateRepositoryStmt:
		v.Loc = loc
	case *AlterRepositoryStmt:
		v.Loc = loc
	case *DropRepositoryStmt:
		v.Loc = loc
	case *CreateFileStmt:
		v.Loc = loc
	case *DropFileStmt:
		v.Loc = loc
	case *CreateResourceStmt:
		v.Loc = loc
	case *AlterResourceStmt:
		v.Loc = loc
	case *DropResourceStmt:
		v.Loc = loc
	case *CreateRowPolicyStmt:
		v.Loc = loc
	case *DropRowPolicyStmt:
		v.Loc = loc
	case *CreateEncryptKeyStmt:
		v.Loc = loc
	case *DropEncryptKeyStmt:
		v.Loc = loc
	case *DictionaryColumn:
		v.Loc = loc
	case *CreateDictionaryStmt:
		v.Loc = loc
	case *AlterDictionaryStmt:
		v.Loc = loc
	case *DropDictionaryStmt:
		v.Loc = loc
	case *RefreshDictionaryStmt:
		v.Loc = loc
	case *CreateRoleStmt:
		v.Loc = loc
	case *AlterRoleStmt:
		v.Loc = loc
	case *DropRoleStmt:
		v.Loc = loc
	case *UserIdentity:
		v.Loc = loc
	case *CreateUserStmt:
		v.Loc = loc
	case *AlterUserStmt:
		v.Loc = loc
	case *DropUserStmt:
		v.Loc = loc
	case *SetPasswordStmt:
		v.Loc = loc
	case *GrantStmt:
		v.Loc = loc
	case *RevokeStmt:
		v.Loc = loc
	case *BeginStmt:
		v.Loc = loc
	case *CommitStmt:
		v.Loc = loc
	case *RollbackStmt:
		v.Loc = loc
	case *ShowStmt:
		v.Loc = loc
	case *DescribeStmt:
		v.Loc = loc
	case *ExplainStmt:
		v.Loc = loc
	case *OptimizerTraceStmt:
		v.Loc = loc
	case *UseStmt:
		v.Loc = loc
	case *SetStmt:
		v.Loc = loc
	case *SetItem:
		v.Loc = loc
	case *UnsetStmt:
		v.Loc = loc
	case *HelpStmt:
		v.Loc = loc
	case *AdminStmt:
		v.Loc = loc
	case *SystemAlterStmt:
		v.Loc = loc
	case *CancelDecommissionStmt:
		v.Loc = loc
	case *BackupStmt:
		v.Loc = loc
	case *RestoreStmt:
		v.Loc = loc
	case *KillStmt:
		v.Loc = loc
	case *LockTablesStmt:
		v.Loc = loc
	case *LockItem:
		v.Loc = loc
	case *UnlockTablesStmt:
		v.Loc = loc
	case *InstallPluginStmt:
		v.Loc = loc
	case *UninstallPluginStmt:
		v.Loc = loc
	case *WarmUpStmt:
		v.Loc = loc
	case *CleanStmt:
		v.Loc = loc
	case *CancelStmt:
		v.Loc = loc
	case *RecoverStmt:
		v.Loc = loc
	case *JobSchedule:
		v.Loc = loc
	case *CreateJobStmt:
		v.Loc = loc
	case *AlterJobStmt:
		v.Loc = loc
	case *DropJobStmt:
		v.Loc = loc
	case *PauseJobStmt:
		v.Loc = loc
	case *ResumeJobStmt:
		v.Loc = loc
	case *CancelTaskStmt:
		v.Loc = loc
	case *ShowJobStmt:
		v.Loc = loc
	case *ShowJobTaskStmt:
		v.Loc = loc
	case *TaskSchedule:
		v.Loc = loc
	case *SubmitTaskStmt:
		v.Loc = loc
	case *DropTaskStmt:
		v.Loc = loc
	case *ShowConstraintsStmt:
		v.Loc = loc
	case *AnalyzeStmt:
		v.Loc = loc
	case *ShowAnalyzeStmt:
		v.Loc = loc
	case *ShowStatsStmt:
		v.Loc = loc
	case *DropStatsStmt:
		v.Loc = loc
	case *KillAnalyzeStmt:
		v.Loc = loc
	case *ProcedureParam:
		v.Loc = loc
	case *CreateProcedureStmt:
		v.Loc = loc
	case *DropProcedureStmt:
		v.Loc = loc
	case *RefreshExternalTableStmt:
		v.Loc = loc
	case *ExternalCatalogRef:
		v.Loc = loc
	case *CreateResourceGroupStmt:
		v.Loc = loc
	case *DropResourceGroupStmt:
		v.Loc = loc
	case *AlterResourceGroupStmt:
		v.Loc = loc
	case *CreateDataCacheRuleStmt:
		v.Loc = loc
	case *DropDataCacheRuleStmt:
		v.Loc = loc
	case *ClearDataCacheRulesStmt:
		v.Loc = loc
	case *CreatePipeStmt:
		v.Loc = loc
	case *DropPipeStmt:
		v.Loc = loc
	case *AlterPipeStmt:
		v.Loc = loc
	case *DescPipeStmt:
		v.Loc = loc
	case *AlterPlanAdvisorStmt:
		v.Loc = loc
	case *TruncatePlanAdvisorStmt:
		v.Loc = loc
	case *TranslateStmt:
		v.Loc = loc
	case *PrepareStmt:
		v.Loc = loc
	case *ExecuteStmt:
		v.Loc = loc
	case *DeallocateStmt:
		v.Loc = loc
	case *ExecuteAsStmt:
		v.Loc = loc
	case *BlacklistStmt:
		v.Loc = loc
	}
}

// SpanNodes returns the smallest Loc that covers every node in nodes.
// Nil entries and nodes whose Loc is invalid are skipped. Returns NoLoc()
// when no node has a valid Loc (including the empty-args case).
func SpanNodes(nodes ...Node) Loc {
	out := NoLoc()
	for _, n := range nodes {
		if n == nil {
			continue
		}
		out = out.Merge(NodeLoc(n))
	}
	return out
}
