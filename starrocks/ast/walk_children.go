package ast

// walkChildren walks the child nodes of node, calling Walk(v, child)
// for each child. Maintained manually until the node count warrants a
// code generator.
func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *File:
		walkNodes(v, n.Stmts)
	case *ObjectName:
		// leaf node, no children
	case *TypeName:
		if n.ElementType != nil {
			Walk(v, n.ElementType)
		}
		if n.ValueType != nil {
			Walk(v, n.ValueType)
		}
		for _, f := range n.Fields {
			if f.Type != nil {
				Walk(v, f.Type)
			}
		}
	case *CreateIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *DropIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *BuildIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *CreateDatabaseStmt:
		Walk(v, n.Name)
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterDatabaseStmt:
		Walk(v, n.Name)
		if n.NewName != nil {
			Walk(v, n.NewName)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropDatabaseStmt:
		Walk(v, n.Name)
	case *DropTableStmt:
		Walk(v, n.Name)
	case *CreateFunctionStmt:
		// leaf-ish node, signature stored as raw text
	case *DropFunctionStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *Property:
		// leaf node, no children

	// Expression nodes (T1.3).
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *UnaryExpr:
		Walk(v, n.Expr)
	case *IsExpr:
		Walk(v, n.Expr)
	case *BetweenExpr:
		Walk(v, n.Expr)
		Walk(v, n.Low)
		Walk(v, n.High)
	case *InExpr:
		Walk(v, n.Expr)
		walkNodes(v, n.Values)
	case *LikeExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
		if n.Escape != nil {
			Walk(v, n.Escape)
		}
	case *RegexpExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
	case *FuncCallExpr:
		Walk(v, n.Name)
		walkNodes(v, n.Args)
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
		if n.Over != nil {
			walkNodes(v, n.Over.PartitionBy)
			for _, item := range n.Over.OrderBy {
				Walk(v, item)
			}
			if n.Over.Frame != nil {
				if n.Over.Frame.StartExpr != nil {
					Walk(v, n.Over.Frame.StartExpr)
				}
				if n.Over.Frame.EndExpr != nil {
					Walk(v, n.Over.Frame.EndExpr)
				}
			}
		}
	case *CastExpr:
		Walk(v, n.Expr)
		Walk(v, n.TypeName)
	case *CaseExpr:
		if n.Operand != nil {
			Walk(v, n.Operand)
		}
		for _, w := range n.Whens {
			Walk(v, w)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *WhenClause:
		Walk(v, n.Cond)
		Walk(v, n.Result)
	case *SubqueryExpr:
		// placeholder leaf — no parsed children yet
	case *ColumnRef:
		Walk(v, n.Name)
	case *Literal:
		// leaf node, no children
	case *ParenExpr:
		Walk(v, n.Expr)
	case *ExistsExpr:
		Walk(v, n.Subquery)
	case *IntervalExpr:
		Walk(v, n.Value)
	case *OrderByItem:
		Walk(v, n.Expr)

	// CTE / WITH clause nodes (T1.6).
	case *WithClause:
		for _, cte := range n.CTEs {
			Walk(v, cte)
		}
	case *CTE:
		if n.Query != nil {
			Walk(v, n.Query)
		}

	// SELECT statement nodes (T1.4).
	case *SelectStmt:
		if n.With != nil {
			Walk(v, n.With)
		}
		for _, item := range n.Items {
			Walk(v, item)
		}
		walkNodes(v, n.From)
		if n.Where != nil {
			Walk(v, n.Where)
		}
		walkNodes(v, n.GroupBy)
		if n.Having != nil {
			Walk(v, n.Having)
		}
		if n.Qualify != nil {
			Walk(v, n.Qualify)
		}
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
		if n.Limit != nil {
			Walk(v, n.Limit)
		}
		if n.Offset != nil {
			Walk(v, n.Offset)
		}
	case *SelectItem:
		if n.Expr != nil {
			Walk(v, n.Expr)
		}
		if n.TableName != nil {
			Walk(v, n.TableName)
		}
	case *TableRef:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *JoinClause:
		Walk(v, n.Left)
		Walk(v, n.Right)
		if n.On != nil {
			Walk(v, n.On)
		}
	case *SetOpStmt:
		Walk(v, n.Left)
		Walk(v, n.Right)

	// DDL — CREATE TABLE nodes (T2.1).
	case *CreateTableStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		for _, idx := range n.Indexes {
			Walk(v, idx)
		}
		for _, c := range n.Constraints {
			Walk(v, c)
		}
		if n.KeyDesc != nil {
			Walk(v, n.KeyDesc)
		}
		if n.PartitionBy != nil {
			Walk(v, n.PartitionBy)
		}
		if n.DistributedBy != nil {
			Walk(v, n.DistributedBy)
		}
		for _, r := range n.Rollup {
			Walk(v, r)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.Like != nil {
			Walk(v, n.Like)
		}
		if n.AsSelect != nil {
			Walk(v, n.AsSelect)
		}
	case *ColumnDef:
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Default != nil {
			Walk(v, n.Default)
		}
		if n.Generated != nil {
			Walk(v, n.Generated)
		}
	case *IndexDef:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *TableConstraint:
		// leaf-ish node, no Node children
	case *KeyDesc:
		// leaf-ish node, no Node children
	case *PartitionDesc:
		for _, p := range n.Partitions {
			Walk(v, p)
		}
	case *PartitionItem:
		// leaf-ish node, values stored as strings
	case *DistributionDesc:
		// leaf-ish node, no Node children
	case *RollupDef:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *RawQuery:
		// leaf node, no parsed children

	// DDL — ALTER TABLE nodes (T2.2).
	case *AlterTableStmt:
		Walk(v, n.Name)
		for _, action := range n.Actions {
			Walk(v, action)
		}
	case *AlterTableAction:
		if n.Column != nil {
			Walk(v, n.Column)
		}
		if n.FieldType != nil {
			Walk(v, n.FieldType)
		}
		if n.NewTableName != nil {
			Walk(v, n.NewTableName)
		}
		if n.Partition != nil {
			Walk(v, n.Partition)
		}
		if n.PartitionDist != nil {
			Walk(v, n.PartitionDist)
		}
		for _, prop := range n.PartitionProps {
			Walk(v, prop)
		}
		if n.Rollup != nil {
			Walk(v, n.Rollup)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.Distribution != nil {
			Walk(v, n.Distribution)
		}

	// DDL — VIEW nodes (T2.4).
	case *ViewColumn:
		// leaf-ish node, no Node children
	case *CreateViewStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}
	case *AlterViewStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}
	case *DropViewStmt:
		Walk(v, n.Name)

	// DML — INSERT statement (T4.1).
	case *InsertStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		for _, row := range n.Values {
			walkNodes(v, row)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}

	// DML — UPDATE / DELETE nodes (T4.2).
	case *Assignment:
		Walk(v, n.Column)
		Walk(v, n.Value)
	case *UpdateStmt:
		Walk(v, n.Target)
		for _, a := range n.Assignments {
			Walk(v, a)
		}
		walkNodes(v, n.From)
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *DeleteStmt:
		Walk(v, n.Target)
		walkNodes(v, n.Using)
		if n.Where != nil {
			Walk(v, n.Where)
		}

	// DML — TRUNCATE TABLE / COPY INTO / LOAD / EXPORT nodes (T6.1).
	case *TruncateTableStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
	case *CopyIntoStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *LoadDataDesc:
		if n.Target != nil {
			Walk(v, n.Target)
		}
	case *LoadDataStmt:
		for _, desc := range n.DataDescs {
			Walk(v, desc)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *ExportStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}

	// Streaming load management nodes (T6.2).
	case *CreateRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.OnTable != nil {
			Walk(v, n.OnTable)
		}
		for _, prop := range n.JobProperties {
			Walk(v, prop)
		}
		for _, prop := range n.DataSourceProperties {
			Walk(v, prop)
		}
	case *AlterRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		for _, prop := range n.DataSourceProperties {
			Walk(v, prop)
		}
	case *PauseRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *ResumeRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *StopRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *ShowRoutineLoadStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *ShowRoutineLoadTaskStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *SyncStmt:
		// leaf node, no children

	// DML — MERGE INTO nodes (T4.3).
	case *MergeStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		if n.Source != nil {
			Walk(v, n.Source)
		}
		if n.On != nil {
			Walk(v, n.On)
		}
		for _, clause := range n.Clauses {
			Walk(v, clause)
		}
	case *MergeClause:
		if n.And != nil {
			Walk(v, n.And)
		}
		for _, a := range n.Assignments {
			if a.Value != nil {
				Walk(v, a.Value)
			}
		}
		for _, val := range n.Values {
			if val != nil {
				Walk(v, val)
			}
		}

	// DDL — Materialized View nodes (T5.1).
	case *MTMVRefreshTrigger:
		// leaf-ish node, no Node children

	case *CreateMTMVStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		if n.RefreshTrigger != nil {
			Walk(v, n.RefreshTrigger)
		}
		if n.PartitionBy != nil {
			Walk(v, n.PartitionBy)
		}
		if n.DistributedBy != nil {
			Walk(v, n.DistributedBy)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}

	case *AlterMTMVStmt:
		Walk(v, n.Name)
		if n.NewName != nil {
			Walk(v, n.NewName)
		}
		if n.ReplaceTarget != nil {
			Walk(v, n.ReplaceTarget)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}

	case *DropMTMVStmt:
		Walk(v, n.Name)
		if n.OnBase != nil {
			Walk(v, n.OnBase)
		}

	case *RefreshMTMVStmt:
		Walk(v, n.Name)

	case *PauseMTMVJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	case *ResumeMTMVJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	case *CancelMTMVTaskStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	// DDL — CATALOG nodes (T5.2).
	case *CreateCatalogStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterCatalogStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropCatalogStmt:
		// leaf-ish node, name stored as string
	case *RefreshCatalogStmt:
	// DDL — STORAGE VAULT nodes (T5.3).
	case *CreateStorageVaultStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterStorageVaultStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStorageVaultStmt:
		// leaf-ish node, no Node children
	case *SetDefaultStorageVaultStmt:
		// leaf-ish node, no Node children
	case *UnsetDefaultStorageVaultStmt:
		// leaf-ish node, no Node children

	// DDL — STORAGE POLICY nodes (T5.3).
	case *CreateStoragePolicyStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterStoragePolicyStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStoragePolicyStmt:
		// leaf-ish node, no Node children

	// DDL — REPOSITORY nodes (T5.3).
	case *CreateRepositoryStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterRepositoryStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropRepositoryStmt:
		// leaf-ish node, no Node children

	// DDL — STAGE nodes (T5.3).
	case *CreateStageStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStageStmt:
		// leaf-ish node, no Node children

	// DDL — FILE nodes (T5.3).
	case *CreateFileStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropFileStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	// Workload management DDL nodes (T5.4).
	case *CreateWorkloadGroupStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterWorkloadGroupStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropWorkloadGroupStmt:
		// leaf-ish node, no Node children
	case *WorkloadPolicyItem:
		// leaf node, raw text only
	case *CreateWorkloadPolicyStmt:
		for _, c := range n.Conditions {
			Walk(v, c)
		}
		for _, a := range n.Actions {
			Walk(v, a)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterWorkloadPolicyStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropWorkloadPolicyStmt:
		// leaf-ish node, no Node children
	case *CreateResourceStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterResourceStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropResourceStmt:
		// leaf-ish node, no Node children
	case *CreateSQLBlockRuleStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterSQLBlockRuleStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropSQLBlockRuleStmt:
		// leaf-ish node, no Node children
	// Security DDL nodes (T5.5).
	case *CreateRowPolicyStmt:
		if n.On != nil {
			Walk(v, n.On)
		}
	case *DropRowPolicyStmt:
		if n.On != nil {
			Walk(v, n.On)
		}
	case *CreateEncryptKeyStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *DropEncryptKeyStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *DictionaryColumn:
		// leaf-ish node, no Node children
	case *CreateDictionaryStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.UsingTable != nil {
			Walk(v, n.UsingTable)
		}
		for _, col := range n.Columns {
			Walk(v, col)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterDictionaryStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropDictionaryStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *RefreshDictionaryStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *CreateRoleStmt:
		// leaf-ish node, no Node children
	case *AlterRoleStmt:
		// leaf-ish node, no Node children
	case *DropRoleStmt:
		// leaf-ish node, no Node children
	case *UserIdentity:
		// leaf node, no Node children
	case *CreateUserStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *AlterUserStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *DropUserStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *SetPasswordStmt:
		if n.For != nil {
			Walk(v, n.For)
		}
	// DCL nodes (T7.1).
	case *GrantStmt:
		if n.Object != nil {
			Walk(v, n.Object)
		}
	case *RevokeStmt:
		if n.Object != nil {
			Walk(v, n.Object)
		}

	// TCL — Transaction Control Language nodes (T7.2).
	case *BeginStmt:
		// leaf node, no Node children
	case *CommitStmt:
		// leaf node, no Node children
	case *RollbackStmt:
		// leaf node, no Node children

	// Utility statement nodes (T7.3).
	case *ShowStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}

	// Job DDL nodes (T8.1).
	case *JobSchedule:
		// leaf-ish node, no Node children

	case *CreateJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Schedule != nil {
			Walk(v, n.Schedule)
		}
		if n.DoStmt != nil {
			Walk(v, n.DoStmt)
		}

	case *AlterJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.NewStmt != nil {
			Walk(v, n.NewStmt)
		}

	case *DropJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *DescribeStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
	case *ExplainStmt:
		if n.Query != nil {
			Walk(v, n.Query)
		}
	case *UseStmt:
		// leaf-ish node, names stored as strings
	case *SetStmt:
		for _, item := range n.Items {
			Walk(v, item)
		}
	case *SetItem:
		if n.Value != nil {
			Walk(v, n.Value)
		}
	case *UnsetStmt:
		// leaf-ish node, names stored as strings
	case *HelpStmt:
		// leaf node
	// Admin / System cluster-management nodes (T7.4).
	case *AdminStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *SystemAlterStmt:
		for _, prop := range n.SetClause {
			Walk(v, prop)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *CancelDecommissionStmt:
		// leaf-ish node, hosts stored as strings

	// Utility / admin command nodes (T7.5).
	case *BackupStmt:
		for _, tbl := range n.Tables {
			Walk(v, tbl)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *RestoreStmt:
		for _, tbl := range n.Tables {
			Walk(v, tbl)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *KillStmt:
		// leaf-ish node, no Node children
	case *LockTablesStmt:
		for _, item := range n.Items {
			Walk(v, item)
		}
	case *LockItem:
		if n.Table != nil {
			Walk(v, n.Table)
		}
	case *UnlockTablesStmt:
		// leaf node, no children
	case *InstallPluginStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *UninstallPluginStmt:
		// leaf node, no children
	case *WarmUpStmt:
		// leaf-ish node, target stored as raw text
	case *CleanStmt:
		// leaf-ish node, target stored as raw text
	case *CancelStmt:
		// leaf-ish node, target stored as raw text
	case *RecoverStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.FromTable != nil {
			Walk(v, n.FromTable)
		}

	case *PauseJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}

	case *ResumeJobStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}

	case *CancelTaskStmt:
		if n.For != nil {
			Walk(v, n.For)
		}

	case *ShowJobStmt:
		if n.Where != nil {
			Walk(v, n.Where)
		}

	case *ShowJobTaskStmt:
		if n.For != nil {
			Walk(v, n.For)
		}

	// Statistics / Analyze / Constraint management nodes (T8.3).
	case *AddConstraintStmt:
		if n.Table != nil {
			Walk(v, n.Table)
		}
		if n.RefTable != nil {
			Walk(v, n.RefTable)
		}
	case *DropConstraintStmt:
		if n.Table != nil {
			Walk(v, n.Table)
		}
	case *ShowConstraintsStmt:
		if n.Table != nil {
			Walk(v, n.Table)
		}
	case *AnalyzeStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *ShowAnalyzeStmt:
		if n.For != nil {
			Walk(v, n.For)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *ShowStatsStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *DropStatsStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
	case *KillAnalyzeStmt:
		// leaf node, no children

	// DDL/DML — Stored Procedure nodes (T8.2).
	case *ProcedureParam:
		if n.Type != nil {
			Walk(v, n.Type)
		}
	case *CreateProcedureStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Parameters {
			Walk(v, param)
		}
	case *CallProcedureStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkNodes(v, n.Args)
	case *DropProcedureStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	}
}
