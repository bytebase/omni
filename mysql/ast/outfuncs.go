package ast

import (
	"fmt"
	"strings"
)

// NodeToString converts a Node to its deterministic string representation.
// Used by comparison tests to verify parser output.
func NodeToString(node Node) string {
	if node == nil {
		return "<>"
	}
	var sb strings.Builder
	writeNode(&sb, node)
	return sb.String()
}

func writeNode(sb *strings.Builder, node Node) {
	if node == nil {
		sb.WriteString("<>")
		return
	}

	switch n := node.(type) {
	case *List:
		writeList(sb, n)
	case *String:
		fmt.Fprintf(sb, "{STRING %q}", n.Str)
	case *Integer:
		fmt.Fprintf(sb, "{INTEGER %d}", n.Ival)

	// Statements
	case *SelectStmt:
		writeSelectStmt(sb, n)
	case *InsertStmt:
		writeInsertStmt(sb, n)
	case *UpdateStmt:
		writeUpdateStmt(sb, n)
	case *DeleteStmt:
		writeDeleteStmt(sb, n)
	case *CreateTableStmt:
		writeCreateTableStmt(sb, n)
	case *AlterTableStmt:
		writeAlterTableStmt(sb, n)
	case *DropTableStmt:
		writeDropTableStmt(sb, n)
	case *CreateIndexStmt:
		writeCreateIndexStmt(sb, n)
	case *CreateViewStmt:
		writeCreateViewStmt(sb, n)
	case *CreateDatabaseStmt:
		writeCreateDatabaseStmt(sb, n)
	case *TruncateStmt:
		writeTruncateStmt(sb, n)
	case *RenameTableStmt:
		writeRenameTableStmt(sb, n)
	case *SetStmt:
		writeSetStmt(sb, n)
	case *SetPasswordStmt:
		writeSetPasswordStmt(sb, n)
	case *ShowStmt:
		writeShowStmt(sb, n)
	case *UseStmt:
		writeUseStmt(sb, n)
	case *ExplainStmt:
		writeExplainStmt(sb, n)
	case *BeginStmt:
		writeBeginStmt(sb, n)
	case *CommitStmt:
		writeCommitStmt(sb, n)
	case *RollbackStmt:
		writeRollbackStmt(sb, n)
	case *SavepointStmt:
		writeSavepointStmt(sb, n)
	case *LockTablesStmt:
		writeLockTablesStmt(sb, n)
	case *UnlockTablesStmt:
		writeUnlockTablesStmt(sb, n)
	case *GrantStmt:
		writeGrantStmt(sb, n)
	case *RevokeStmt:
		writeRevokeStmt(sb, n)
	case *CreateUserStmt:
		writeCreateUserStmt(sb, n)
	case *DropUserStmt:
		writeDropUserStmt(sb, n)
	case *AlterUserStmt:
		writeAlterUserStmt(sb, n)
	case *CreateRoleStmt:
		writeCreateRoleStmt(sb, n)
	case *DropRoleStmt:
		writeDropRoleStmt(sb, n)
	case *SetDefaultRoleStmt:
		writeSetDefaultRoleStmt(sb, n)
	case *SetRoleStmt:
		writeSetRoleStmt(sb, n)
	case *GrantRoleStmt:
		writeGrantRoleStmt(sb, n)
	case *RevokeRoleStmt:
		writeRevokeRoleStmt(sb, n)
	case *CreateFunctionStmt:
		writeCreateFunctionStmt(sb, n)
	case *CreateTriggerStmt:
		writeCreateTriggerStmt(sb, n)
	case *CreateEventStmt:
		writeCreateEventStmt(sb, n)
	case *LoadDataStmt:
		writeLoadDataStmt(sb, n)
	case *PrepareStmt:
		writePrepareStmt(sb, n)
	case *ExecuteStmt:
		writeExecuteStmt(sb, n)
	case *DeallocateStmt:
		writeDeallocateStmt(sb, n)
	case *AnalyzeTableStmt:
		writeAnalyzeTableStmt(sb, n)
	case *OptimizeTableStmt:
		writeOptimizeTableStmt(sb, n)
	case *CheckTableStmt:
		writeCheckTableStmt(sb, n)
	case *RepairTableStmt:
		writeRepairTableStmt(sb, n)
	case *FlushStmt:
		writeFlushStmt(sb, n)
	case *KillStmt:
		writeKillStmt(sb, n)
	case *DoStmt:
		writeDoStmt(sb, n)
	case *ChecksumTableStmt:
		writeChecksumTableStmt(sb, n)
	case *ShutdownStmt:
		writeShutdownStmt(sb, n)
	case *RestartStmt:
		writeRestartStmt(sb, n)
	case *AlterDatabaseStmt:
		writeAlterDatabaseStmt(sb, n)
	case *DropDatabaseStmt:
		writeDropDatabaseStmt(sb, n)
	case *DropIndexStmt:
		writeDropIndexStmt(sb, n)
	case *DropViewStmt:
		writeDropViewStmt(sb, n)
	case *RawStmt:
		writeRawStmt(sb, n)
	case *TableStmt:
		writeTableStmt(sb, n)
	case *ValuesStmt:
		writeValuesStmt(sb, n)

	// Expressions
	case *BinaryExpr:
		writeBinaryExpr(sb, n)
	case *UnaryExpr:
		writeUnaryExpr(sb, n)
	case *FuncCallExpr:
		writeFuncCallExpr(sb, n)
	case *SubqueryExpr:
		writeSubqueryExpr(sb, n)
	case *CaseExpr:
		writeCaseExpr(sb, n)
	case *BetweenExpr:
		writeBetweenExpr(sb, n)
	case *InExpr:
		writeInExpr(sb, n)
	case *LikeExpr:
		writeLikeExpr(sb, n)
	case *IsExpr:
		writeIsExpr(sb, n)
	case *ExistsExpr:
		writeExistsExpr(sb, n)
	case *CastExpr:
		writeCastExpr(sb, n)
	case *ExtractExpr:
		writeExtractExpr(sb, n)
	case *ParenExpr:
		writeParenExpr(sb, n)
	case *CollateExpr:
		writeCollateExpr(sb, n)
	case *IntervalExpr:
		writeIntervalExpr(sb, n)
	case *MatchExpr:
		writeMatchExpr(sb, n)
	case *ConvertExpr:
		writeConvertExpr(sb, n)
	case *DefaultExpr:
		writeDefaultExpr(sb, n)
	case *RowExpr:
		writeRowExpr(sb, n)
	case *StarExpr:
		fmt.Fprintf(sb, "{STAR :loc %d}", n.Loc.Start)

	// References and literals
	case *ColumnRef:
		writeColumnRef(sb, n)
	case *TableRef:
		writeTableRef(sb, n)
	case *IntLit:
		fmt.Fprintf(sb, "{INT_LIT :val %d :loc %d}", n.Value, n.Loc.Start)
	case *FloatLit:
		fmt.Fprintf(sb, "{FLOAT_LIT :val %s :loc %d}", n.Value, n.Loc.Start)
	case *StringLit:
		writeStringLit(sb, n)
	case *BoolLit:
		fmt.Fprintf(sb, "{BOOL_LIT :val %t :loc %d}", n.Value, n.Loc.Start)
	case *NullLit:
		fmt.Fprintf(sb, "{NULL_LIT :loc %d}", n.Loc.Start)
	case *HexLit:
		fmt.Fprintf(sb, "{HEX_LIT :val %s :loc %d}", n.Value, n.Loc.Start)
	case *BitLit:
		fmt.Fprintf(sb, "{BIT_LIT :val %s :loc %d}", n.Value, n.Loc.Start)
	case *TemporalLit:
		fmt.Fprintf(sb, "{TEMPORAL_LIT :type %s :val %s :loc %d}", n.Type, n.Value, n.Loc.Start)
	case *VariableRef:
		writeVariableRef(sb, n)

	// Clause helpers
	case *Assignment:
		writeAssignment(sb, n)
	case *OrderByItem:
		writeOrderByItem(sb, n)
	case *Limit:
		writeLimit(sb, n)
	case *JoinClause:
		writeJoinClause(sb, n)
	case *OnCondition:
		writeOnCondition(sb, n)
	case *UsingCondition:
		writeUsingCondition(sb, n)
	case *ResTarget:
		writeResTarget(sb, n)
	case *WindowDef:
		writeWindowDef(sb, n)
	case *WindowFrame:
		writeWindowFrame(sb, n)
	case *WindowFrameBound:
		writeWindowFrameBound(sb, n)
	case *ColumnDef:
		writeColumnDef(sb, n)
	case *GeneratedColumn:
		writeGeneratedColumn(sb, n)
	case *ColumnConstraint:
		writeColumnConstraint(sb, n)
	case *Constraint:
		writeConstraint(sb, n)
	case *DataType:
		writeDataType(sb, n)
	case *AlterTableCmd:
		writeAlterTableCmd(sb, n)
	case *TableOption:
		if n.Storage != "" {
			fmt.Fprintf(sb, "{TABLE_OPT :loc %d :name %s :val %s :storage %s}", n.Loc.Start, n.Name, n.Value, n.Storage)
		} else {
			fmt.Fprintf(sb, "{TABLE_OPT :loc %d :name %s :val %s}", n.Loc.Start, n.Name, n.Value)
		}
	case *DatabaseOption:
		fmt.Fprintf(sb, "{DB_OPT :loc %d :name %s :val %s}", n.Loc.Start, n.Name, n.Value)
	case *PartitionClause:
		writePartitionClause(sb, n)
	case *PartitionDef:
		writePartitionDef(sb, n)
	case *SubPartitionDef:
		writeSubPartitionDef(sb, n)
	case *IndexColumn:
		writeIndexColumn(sb, n)
	case *IndexOption:
		writeIndexOption(sb, n)
	case *IndexHint:
		writeIndexHint(sb, n)
	case *CaseWhen:
		writeCaseWhen(sb, n)
	case *ForUpdate:
		writeForUpdate(sb, n)
	case *IntoClause:
		writeIntoClause(sb, n)
	case *RenameTablePair:
		writeRenameTablePair(sb, n)
	case *LockTable:
		writeLockTable(sb, n)
	case *GrantTarget:
		writeGrantTarget(sb, n)
	case *RequireClause:
		writeRequireClause(sb, n)
	case *ResourceOption:
		writeResourceOption(sb, n)
	case *UserSpec:
		writeUserSpec(sb, n)
	case *FuncParam:
		writeFuncParam(sb, n)
	case *RoutineCharacteristic:
		writeRoutineCharacteristic(sb, n)
	case *TriggerOrder:
		writeTriggerOrder(sb, n)
	case *EventSchedule:
		writeEventSchedule(sb, n)
	case *CommonTableExpr:
		writeCommonTableExpr(sb, n)
	case *SetTransactionStmt:
		writeSetTransactionStmt(sb, n)
	case *XAStmt:
		writeXAStmt(sb, n)
	case *MemberOfExpr:
		writeMemberOfExpr(sb, n)
	case *JsonTableExpr:
		writeJsonTableExpr(sb, n)
	case *JsonTableColumn:
		writeJsonTableColumn(sb, n)
	case *CallStmt:
		writeCallStmt(sb, n)
	case *HandlerOpenStmt:
		writeHandlerOpenStmt(sb, n)
	case *HandlerReadStmt:
		writeHandlerReadStmt(sb, n)
	case *HandlerCloseStmt:
		writeHandlerCloseStmt(sb, n)
	case *SignalStmt:
		writeSignalStmt(sb, n)
	case *ResignalStmt:
		writeResignalStmt(sb, n)
	case *SignalInfoItem:
		writeSignalInfoItem(sb, n)
	case *GetDiagnosticsStmt:
		writeGetDiagnosticsStmt(sb, n)
	case *DiagnosticsItem:
		writeDiagnosticsItem(sb, n)
	case *BeginEndBlock:
		writeBeginEndBlock(sb, n)
	case *DeclareVarStmt:
		writeDeclareVarStmt(sb, n)
	case *DeclareConditionStmt:
		writeDeclareConditionStmt(sb, n)
	case *DeclareHandlerStmt:
		writeDeclareHandlerStmt(sb, n)
	case *DeclareCursorStmt:
		writeDeclareCursorStmt(sb, n)
	case *IfStmt:
		writeIfStmt(sb, n)
	case *ElseIf:
		writeElseIf(sb, n)
	case *CaseStmtNode:
		writeCaseStmtNode(sb, n)
	case *CaseStmtWhen:
		writeCaseStmtWhen(sb, n)
	case *WhileStmt:
		writeWhileStmt(sb, n)
	case *RepeatStmt:
		writeRepeatStmt(sb, n)
	case *LoopStmt:
		writeLoopStmt(sb, n)
	case *LeaveStmt:
		writeLeaveStmt(sb, n)
	case *IterateStmt:
		writeIterateStmt(sb, n)
	case *ReturnStmt:
		writeReturnStmt(sb, n)
	case *OpenCursorStmt:
		writeOpenCursorStmt(sb, n)
	case *FetchCursorStmt:
		writeFetchCursorStmt(sb, n)
	case *CloseCursorStmt:
		writeCloseCursorStmt(sb, n)
	case *CloneStmt:
		writeCloneStmt(sb, n)
	case *InstallPluginStmt:
		writeInstallPluginStmt(sb, n)
	case *UninstallPluginStmt:
		writeUninstallPluginStmt(sb, n)
	case *InstallComponentStmt:
		writeInstallComponentStmt(sb, n)
	case *UninstallComponentStmt:
		writeUninstallComponentStmt(sb, n)
	case *CreateTablespaceStmt:
		writeCreateTablespaceStmt(sb, n)
	case *AlterTablespaceStmt:
		writeAlterTablespaceStmt(sb, n)
	case *DropTablespaceStmt:
		writeDropTablespaceStmt(sb, n)
	case *CreateServerStmt:
		writeCreateServerStmt(sb, n)
	case *AlterServerStmt:
		writeAlterServerStmt(sb, n)
	case *DropServerStmt:
		writeDropServerStmt(sb, n)
	case *CreateLogfileGroupStmt:
		writeCreateLogfileGroupStmt(sb, n)
	case *AlterLogfileGroupStmt:
		writeAlterLogfileGroupStmt(sb, n)
	case *DropLogfileGroupStmt:
		writeDropLogfileGroupStmt(sb, n)
	case *CreateSpatialRefSysStmt:
		writeCreateSpatialRefSysStmt(sb, n)
	case *DropSpatialRefSysStmt:
		writeDropSpatialRefSysStmt(sb, n)
	case *CreateResourceGroupStmt:
		writeCreateResourceGroupStmt(sb, n)
	case *AlterResourceGroupStmt:
		writeAlterResourceGroupStmt(sb, n)
	case *DropResourceGroupStmt:
		writeDropResourceGroupStmt(sb, n)
	case *AlterViewStmt:
		writeAlterViewStmt(sb, n)
	case *AlterEventStmt:
		writeAlterEventStmt(sb, n)
	case *AlterRoutineStmt:
		writeAlterRoutineStmt(sb, n)
	case *DropRoutineStmt:
		writeDropRoutineStmt(sb, n)
	case *DropTriggerStmt:
		writeDropTriggerStmt(sb, n)
	case *DropEventStmt:
		writeDropEventStmt(sb, n)
	case *ChangeReplicationSourceStmt:
		writeChangeReplicationSourceStmt(sb, n)
	case *ReplicationOption:
		writeReplicationOption(sb, n)
	case *ChangeReplicationFilterStmt:
		writeChangeReplicationFilterStmt(sb, n)
	case *ReplicationFilter:
		writeReplicationFilter(sb, n)
	case *StartReplicaStmt:
		writeStartReplicaStmt(sb, n)
	case *StopReplicaStmt:
		writeStopReplicaStmt(sb, n)
	case *ResetReplicaStmt:
		writeResetReplicaStmt(sb, n)
	case *PurgeBinaryLogsStmt:
		writePurgeBinaryLogsStmt(sb, n)
	case *ResetMasterStmt:
		writeResetMasterStmt(sb, n)
	case *StartGroupReplicationStmt:
		writeStartGroupReplicationStmt(sb, n)
	case *StopGroupReplicationStmt:
		writeStopGroupReplicationStmt(sb, n)
	case *AlterInstanceStmt:
		writeAlterInstanceStmt(sb, n)
	case *LockInstanceStmt:
		fmt.Fprintf(sb, "{LOCK_INSTANCE :loc %d}", n.Loc.Start)
	case *UnlockInstanceStmt:
		fmt.Fprintf(sb, "{UNLOCK_INSTANCE :loc %d}", n.Loc.Start)
	case *ImportTableStmt:
		writeImportTableStmt(sb, n)
	case *BinlogStmt:
		fmt.Fprintf(sb, "{BINLOG :loc %d :str %q}", n.Loc.Start, n.Str)
	case *CacheIndexStmt:
		writeCacheIndexStmt(sb, n)
	case *LoadIndexIntoCacheStmt:
		writeLoadIndexIntoCacheStmt(sb, n)
	case *LoadIndexTable:
		writeLoadIndexTable(sb, n)
	case *ResetPersistStmt:
		writeResetPersistStmt(sb, n)
	case *RenameUserStmt:
		writeRenameUserStmt(sb, n)
	case *RenameUserPair:
		writeRenameUserPair(sb, n)
	case *SetResourceGroupStmt:
		writeSetResourceGroupStmt(sb, n)
	case *HelpStmt:
		fmt.Fprintf(sb, "{HELP :loc %d :topic %q}", n.Loc.Start, n.Topic)
	case *VCPUSpec:
		writeVCPUSpec(sb, n)

	default:
		fmt.Fprintf(sb, "{UNKNOWN %T}", node)
	}
}

func writeList(sb *strings.Builder, l *List) {
	if l == nil || len(l.Items) == 0 {
		sb.WriteString("()")
		return
	}
	sb.WriteString("(")
	for i, item := range l.Items {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, item)
	}
	sb.WriteString(")")
}

// -----------------------------------------------------------------------
// Statement writers
// -----------------------------------------------------------------------

func writeSelectStmt(sb *strings.Builder, n *SelectStmt) {
	sb.WriteString("{SELECT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.CTEs) > 0 {
		sb.WriteString(" :ctes (")
		for i, cte := range n.CTEs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, cte)
		}
		sb.WriteString(")")
	}
	if n.DistinctKind == DistinctOn {
		sb.WriteString(" :distinct true")
	}
	if n.CalcFoundRows {
		sb.WriteString(" :calc_found_rows true")
	}
	if n.HighPriority {
		sb.WriteString(" :high_priority true")
	}
	if n.StraightJoin {
		sb.WriteString(" :straight_join true")
	}
	if n.SmallResult {
		sb.WriteString(" :small_result true")
	}
	if n.BigResult {
		sb.WriteString(" :big_result true")
	}
	if n.BufferResult {
		sb.WriteString(" :buffer_result true")
	}
	if n.NoCache {
		sb.WriteString(" :no_cache true")
	}
	if len(n.TargetList) > 0 {
		sb.WriteString(" :targets ")
		writeExprNodeList(sb, n.TargetList)
	}
	if len(n.From) > 0 {
		sb.WriteString(" :from ")
		writeTableExprList(sb, n.From)
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if len(n.GroupBy) > 0 {
		sb.WriteString(" :group_by ")
		writeExprNodeList(sb, n.GroupBy)
		if n.WithRollup {
			sb.WriteString(" :with_rollup true")
		}
	}
	if n.Having != nil {
		sb.WriteString(" :having ")
		writeNode(sb, n.Having)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
		if n.OrderByWithRollup {
			sb.WriteString(" :order_by_with_rollup true")
		}
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	if n.ForUpdate != nil {
		sb.WriteString(" :for_update ")
		writeNode(sb, n.ForUpdate)
	}
	if len(n.WindowClause) > 0 {
		sb.WriteString(" :window ")
		for i, w := range n.WindowClause {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, w)
		}
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.SetOp != SetOpNone {
		sb.WriteString(" :set_op ")
		switch n.SetOp {
		case SetOpUnion:
			sb.WriteString("UNION")
		case SetOpIntersect:
			sb.WriteString("INTERSECT")
		case SetOpExcept:
			sb.WriteString("EXCEPT")
		}
		if n.SetAll {
			sb.WriteString(" ALL")
		}
		if n.Left != nil {
			sb.WriteString(" :left ")
			writeNode(sb, n.Left)
		}
		if n.Right != nil {
			sb.WriteString(" :right ")
			writeNode(sb, n.Right)
		}
	}
	sb.WriteString("}")
}

func writeInsertStmt(sb *strings.Builder, n *InsertStmt) {
	sb.WriteString("{INSERT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IsReplace {
		sb.WriteString(" :replace true")
	}
	if n.Priority != InsertPriorityNone {
		fmt.Fprintf(sb, " :priority %d", n.Priority)
	}
	if n.Ignore {
		sb.WriteString(" :ignore true")
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.Partitions) > 0 {
		sb.WriteString(" :partitions (")
		for i, p := range n.Partitions {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns ")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	if len(n.Values) > 0 {
		sb.WriteString(" :values ")
		for i, row := range n.Values {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString("(")
			writeExprNodeList(sb, row)
			sb.WriteString(")")
		}
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	if n.TableSource != nil {
		sb.WriteString(" :table_source ")
		writeNode(sb, n.TableSource)
	}
	if len(n.SetList) > 0 {
		sb.WriteString(" :set ")
		for i, a := range n.SetList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
	}
	if n.RowAlias != "" {
		fmt.Fprintf(sb, " :row_alias %s", n.RowAlias)
	}
	if len(n.ColAliases) > 0 {
		sb.WriteString(" :col_aliases (")
		for i, c := range n.ColAliases {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(c)
		}
		sb.WriteString(")")
	}
	if len(n.OnDuplicateKey) > 0 {
		sb.WriteString(" :on_dup ")
		for i, a := range n.OnDuplicateKey {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
	}
	if len(n.Returning) > 0 {
		sb.WriteString(" :returning ")
		writeNodeList(sb, n.Returning)
	}
	sb.WriteString("}")
}

func writeUpdateStmt(sb *strings.Builder, n *UpdateStmt) {
	sb.WriteString("{UPDATE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.LowPriority {
		sb.WriteString(" :low_priority true")
	}
	if n.Ignore {
		sb.WriteString(" :ignore true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		writeTableExprList(sb, n.Tables)
	}
	if len(n.SetList) > 0 {
		sb.WriteString(" :set ")
		for i, a := range n.SetList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	sb.WriteString("}")
}

func writeDeleteStmt(sb *strings.Builder, n *DeleteStmt) {
	sb.WriteString("{DELETE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.LowPriority {
		sb.WriteString(" :low_priority true")
	}
	if n.Quick {
		sb.WriteString(" :quick true")
	}
	if n.Ignore {
		sb.WriteString(" :ignore true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		writeTableExprList(sb, n.Tables)
	}
	if len(n.Using) > 0 {
		sb.WriteString(" :using ")
		writeTableExprList(sb, n.Using)
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	sb.WriteString("}")
}

func writeCreateTableStmt(sb *strings.Builder, n *CreateTableStmt) {
	sb.WriteString("{CREATE_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.Temporary {
		sb.WriteString(" :temporary true")
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns ")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	if len(n.Constraints) > 0 {
		sb.WriteString(" :constraints ")
		for i, c := range n.Constraints {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	if n.Partitions != nil {
		sb.WriteString(" :partitions ")
		writeNode(sb, n.Partitions)
	}
	if n.Like != nil {
		sb.WriteString(" :like ")
		writeNode(sb, n.Like)
	}
	if n.Ignore {
		sb.WriteString(" :ignore true")
	}
	if n.Replace {
		sb.WriteString(" :replace true")
	}
	if n.Select != nil {
		sb.WriteString(" :as_select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeAlterTableStmt(sb *strings.Builder, n *AlterTableStmt) {
	sb.WriteString("{ALTER_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.Commands) > 0 {
		sb.WriteString(" :cmds ")
		for i, c := range n.Commands {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	sb.WriteString("}")
}

func writeDropTableStmt(sb *strings.Builder, n *DropTableStmt) {
	sb.WriteString("{DROP_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.Temporary {
		sb.WriteString(" :temporary true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.Restrict {
		sb.WriteString(" :restrict true")
	}
	sb.WriteString("}")
}

func writeCreateIndexStmt(sb *strings.Builder, n *CreateIndexStmt) {
	sb.WriteString("{CREATE_INDEX")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Unique {
		sb.WriteString(" :unique true")
	}
	if n.Fulltext {
		sb.WriteString(" :fulltext true")
	}
	if n.Spatial {
		sb.WriteString(" :spatial true")
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.IndexName != "" {
		fmt.Fprintf(sb, " :name %s", n.IndexName)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.IndexType != "" {
		fmt.Fprintf(sb, " :type %s", n.IndexType)
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns ")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	if n.Algorithm != "" {
		fmt.Fprintf(sb, " :algorithm %s", n.Algorithm)
	}
	if n.Lock != "" {
		fmt.Fprintf(sb, " :lock %s", n.Lock)
	}
	sb.WriteString("}")
}

func writeCreateViewStmt(sb *strings.Builder, n *CreateViewStmt) {
	sb.WriteString("{CREATE_VIEW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.Algorithm != "" {
		fmt.Fprintf(sb, " :algorithm %s", n.Algorithm)
	}
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.SqlSecurity != "" {
		fmt.Fprintf(sb, " :sql_security %s", n.SqlSecurity)
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Columns) > 0 {
		fmt.Fprintf(sb, " :columns %s", strings.Join(n.Columns, ", "))
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	if n.CheckOption != "" {
		fmt.Fprintf(sb, " :check_option %s", n.CheckOption)
	}
	sb.WriteString("}")
}

func writeCreateDatabaseStmt(sb *strings.Builder, n *CreateDatabaseStmt) {
	sb.WriteString("{CREATE_DATABASE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	sb.WriteString("}")
}

func writeTruncateStmt(sb *strings.Builder, n *TruncateStmt) {
	sb.WriteString("{TRUNCATE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	sb.WriteString("}")
}

func writeRenameTableStmt(sb *strings.Builder, n *RenameTableStmt) {
	sb.WriteString("{RENAME_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Pairs) > 0 {
		sb.WriteString(" :pairs ")
		for i, p := range n.Pairs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
	}
	sb.WriteString("}")
}

func writeSetStmt(sb *strings.Builder, n *SetStmt) {
	sb.WriteString("{SET")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Scope != "" {
		fmt.Fprintf(sb, " :scope %s", n.Scope)
	}
	if len(n.Assignments) > 0 {
		sb.WriteString(" :assignments ")
		for i, a := range n.Assignments {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
	}
	sb.WriteString("}")
}

func writeSetPasswordStmt(sb *strings.Builder, n *SetPasswordStmt) {
	sb.WriteString("{SET_PASSWORD")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.User != nil {
		sb.WriteString(" :user ")
		writeNode(sb, n.User)
	}
	if n.Password != "" {
		fmt.Fprintf(sb, " :password %s", n.Password)
	}
	if n.ToRandom {
		sb.WriteString(" :to_random true")
	}
	if n.Replace != "" {
		fmt.Fprintf(sb, " :replace %s", n.Replace)
	}
	if n.RetainCurrentPassword {
		sb.WriteString(" :retain_current_password true")
	}
	sb.WriteString("}")
}

func writeShowStmt(sb *strings.Builder, n *ShowStmt) {
	sb.WriteString("{SHOW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Type != "" {
		fmt.Fprintf(sb, " :type %s", n.Type)
	}
	if n.From != nil {
		sb.WriteString(" :from ")
		writeNode(sb, n.From)
	}
	if n.Like != nil {
		sb.WriteString(" :like ")
		writeNode(sb, n.Like)
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if len(n.ProfileTypes) > 0 {
		sb.WriteString(" :profile_types (")
		for i, pt := range n.ProfileTypes {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(pt)
		}
		sb.WriteString(")")
	}
	if n.ForQuery != nil {
		sb.WriteString(" :for_query ")
		writeNode(sb, n.ForQuery)
	}
	if n.ForUser != nil {
		sb.WriteString(" :for_user ")
		writeNode(sb, n.ForUser)
	}
	if len(n.Using) > 0 {
		sb.WriteString(" :using (")
		for i, u := range n.Using {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, u)
		}
		sb.WriteString(")")
	}
	if n.FromPos != nil {
		sb.WriteString(" :from_pos ")
		writeNode(sb, n.FromPos)
	}
	if n.LimitCount != nil {
		sb.WriteString(" :limit_count ")
		writeNode(sb, n.LimitCount)
	}
	if n.LimitOffset != nil {
		sb.WriteString(" :limit_offset ")
		writeNode(sb, n.LimitOffset)
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writeUseStmt(sb *strings.Builder, n *UseStmt) {
	sb.WriteString("{USE")
	fmt.Fprintf(sb, " :loc %d :database %s", n.Loc.Start, n.Database)
	sb.WriteString("}")
}

func writeExplainStmt(sb *strings.Builder, n *ExplainStmt) {
	sb.WriteString("{EXPLAIN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Analyze {
		sb.WriteString(" :analyze true")
	}
	if n.Extended {
		sb.WriteString(" :extended true")
	}
	if n.Partitions {
		sb.WriteString(" :partitions true")
	}
	if n.Format != "" {
		fmt.Fprintf(sb, " :format %s", n.Format)
	}
	if n.ForConnection != 0 {
		fmt.Fprintf(sb, " :for_connection %d", n.ForConnection)
	}
	if n.Stmt != nil {
		sb.WriteString(" :stmt ")
		writeNode(sb, n.Stmt)
	}
	sb.WriteString("}")
}

func writeBeginStmt(sb *strings.Builder, n *BeginStmt) {
	sb.WriteString("{BEGIN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.ReadOnly {
		sb.WriteString(" :read_only true")
	}
	if n.ReadWrite {
		sb.WriteString(" :read_write true")
	}
	if n.WithConsistentSnapshot {
		sb.WriteString(" :consistent_snapshot true")
	}
	sb.WriteString("}")
}

func writeCommitStmt(sb *strings.Builder, n *CommitStmt) {
	sb.WriteString("{COMMIT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Chain {
		sb.WriteString(" :chain true")
	}
	if n.Release {
		sb.WriteString(" :release true")
	}
	sb.WriteString("}")
}

func writeRollbackStmt(sb *strings.Builder, n *RollbackStmt) {
	sb.WriteString("{ROLLBACK")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Savepoint != "" {
		fmt.Fprintf(sb, " :savepoint %s", n.Savepoint)
	}
	if n.Chain {
		sb.WriteString(" :chain true")
	}
	if n.Release {
		sb.WriteString(" :release true")
	}
	sb.WriteString("}")
}

func writeSavepointStmt(sb *strings.Builder, n *SavepointStmt) {
	if n.Release {
		sb.WriteString("{RELEASE_SAVEPOINT")
	} else {
		sb.WriteString("{SAVEPOINT")
	}
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	sb.WriteString("}")
}

func writeLockTablesStmt(sb *strings.Builder, n *LockTablesStmt) {
	sb.WriteString("{LOCK_TABLES")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	sb.WriteString("}")
}

func writeUnlockTablesStmt(sb *strings.Builder, n *UnlockTablesStmt) {
	sb.WriteString("{UNLOCK_TABLES")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	sb.WriteString("}")
}

func writeGrantStmt(sb *strings.Builder, n *GrantStmt) {
	sb.WriteString("{GRANT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.AllPriv {
		sb.WriteString(" :all_priv true")
	}
	if len(n.Privileges) > 0 {
		fmt.Fprintf(sb, " :privileges %s", strings.Join(n.Privileges, ", "))
	}
	if n.ProxyUser != "" {
		fmt.Fprintf(sb, " :proxy_user %s", n.ProxyUser)
	}
	if n.On != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.On)
	}
	if len(n.To) > 0 {
		fmt.Fprintf(sb, " :to %s", strings.Join(n.To, ", "))
	}
	if n.WithGrant {
		sb.WriteString(" :with_grant true")
	}
	if n.AsUser != "" {
		fmt.Fprintf(sb, " :as %s", n.AsUser)
	}
	if n.WithRoleType != "" {
		fmt.Fprintf(sb, " :with_role_type %s", n.WithRoleType)
	}
	if len(n.WithRoles) > 0 {
		fmt.Fprintf(sb, " :with_roles %s", strings.Join(n.WithRoles, ", "))
	}
	if n.Require != nil {
		sb.WriteString(" :require ")
		writeNode(sb, n.Require)
	}
	if n.Resource != nil {
		sb.WriteString(" :resource ")
		writeNode(sb, n.Resource)
	}
	sb.WriteString("}")
}

func writeRevokeStmt(sb *strings.Builder, n *RevokeStmt) {
	sb.WriteString("{REVOKE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.AllPriv {
		sb.WriteString(" :all_priv true")
	}
	if n.GrantOption {
		sb.WriteString(" :grant_option true")
	}
	if len(n.Privileges) > 0 {
		fmt.Fprintf(sb, " :privileges %s", strings.Join(n.Privileges, ", "))
	}
	if n.On != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.On)
	}
	if len(n.From) > 0 {
		fmt.Fprintf(sb, " :from %s", strings.Join(n.From, ", "))
	}
	if n.IgnoreUnknownUser {
		sb.WriteString(" :ignore_unknown_user true")
	}
	sb.WriteString("}")
}

func writeCreateUserStmt(sb *strings.Builder, n *CreateUserStmt) {
	sb.WriteString("{CREATE_USER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if len(n.Users) > 0 {
		sb.WriteString(" :users ")
		for i, u := range n.Users {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, u)
		}
	}
	if len(n.DefaultRoles) > 0 {
		fmt.Fprintf(sb, " :default_roles %s", strings.Join(n.DefaultRoles, ", "))
	}
	if n.Require != nil {
		sb.WriteString(" :require ")
		writeNode(sb, n.Require)
	}
	if n.Resource != nil {
		sb.WriteString(" :resource ")
		writeNode(sb, n.Resource)
	}
	if n.PasswordExpire != "" {
		fmt.Fprintf(sb, " :password_expire %s", n.PasswordExpire)
	}
	if n.PasswordHistory != "" {
		fmt.Fprintf(sb, " :password_history %s", n.PasswordHistory)
	}
	if n.PasswordReuseInterval != "" {
		fmt.Fprintf(sb, " :password_reuse_interval %s", n.PasswordReuseInterval)
	}
	if n.PasswordRequireCurrent != "" {
		fmt.Fprintf(sb, " :password_require_current %s", n.PasswordRequireCurrent)
	}
	if n.HasFailedLogin {
		fmt.Fprintf(sb, " :failed_login_attempts %d", n.FailedLoginAttempts)
	}
	if n.PasswordLockTime != "" {
		fmt.Fprintf(sb, " :password_lock_time %s", n.PasswordLockTime)
	}
	if n.AccountLock {
		sb.WriteString(" :account_lock true")
	}
	if n.AccountUnlock {
		sb.WriteString(" :account_unlock true")
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.Attribute != "" {
		fmt.Fprintf(sb, " :attribute %q", n.Attribute)
	}
	sb.WriteString("}")
}

func writeDropUserStmt(sb *strings.Builder, n *DropUserStmt) {
	sb.WriteString("{DROP_USER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if len(n.Users) > 0 {
		fmt.Fprintf(sb, " :users %s", strings.Join(n.Users, ", "))
	}
	sb.WriteString("}")
}

func writeAlterUserStmt(sb *strings.Builder, n *AlterUserStmt) {
	sb.WriteString("{ALTER_USER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.DefaultRoleUser != "" {
		fmt.Fprintf(sb, " :default_role_user %s", n.DefaultRoleUser)
	}
	if n.DefaultRoleType != "" {
		fmt.Fprintf(sb, " :default_role_type %s", n.DefaultRoleType)
	}
	if len(n.DefaultRoles) > 0 {
		fmt.Fprintf(sb, " :default_roles %s", strings.Join(n.DefaultRoles, ", "))
	}
	if len(n.Users) > 0 {
		sb.WriteString(" :users ")
		for i, u := range n.Users {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, u)
		}
	}
	if n.Require != nil {
		sb.WriteString(" :require ")
		writeNode(sb, n.Require)
	}
	if n.Resource != nil {
		sb.WriteString(" :resource ")
		writeNode(sb, n.Resource)
	}
	if n.PasswordExpire != "" {
		fmt.Fprintf(sb, " :password_expire %s", n.PasswordExpire)
	}
	if n.PasswordHistory != "" {
		fmt.Fprintf(sb, " :password_history %s", n.PasswordHistory)
	}
	if n.PasswordReuseInterval != "" {
		fmt.Fprintf(sb, " :password_reuse_interval %s", n.PasswordReuseInterval)
	}
	if n.PasswordRequireCurrent != "" {
		fmt.Fprintf(sb, " :password_require_current %s", n.PasswordRequireCurrent)
	}
	if n.HasFailedLogin {
		fmt.Fprintf(sb, " :failed_login_attempts %d", n.FailedLoginAttempts)
	}
	if n.PasswordLockTime != "" {
		fmt.Fprintf(sb, " :password_lock_time %s", n.PasswordLockTime)
	}
	if n.AccountLock {
		sb.WriteString(" :account_lock true")
	}
	if n.AccountUnlock {
		sb.WriteString(" :account_unlock true")
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.Attribute != "" {
		fmt.Fprintf(sb, " :attribute %q", n.Attribute)
	}
	if n.IsUserFunc {
		sb.WriteString(" :is_user_func true")
	}
	if len(n.FactorOps) > 0 {
		sb.WriteString(" :factor_ops (")
		for i, op := range n.FactorOps {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "{%s %dFACTOR", op.Action, op.Factor)
			if op.AuthPlugin != "" {
				fmt.Fprintf(sb, " :plugin %s", op.AuthPlugin)
			}
			if op.Password != "" {
				sb.WriteString(" :password ***")
			}
			if op.AuthHash != "" {
				sb.WriteString(" :hash ***")
			}
			if op.PasswordRandom {
				sb.WriteString(" :random true")
			}
			sb.WriteString("}")
		}
		sb.WriteString(")")
	}
	if n.RegistrationOp != nil {
		op := n.RegistrationOp
		fmt.Fprintf(sb, " :registration {%dFACTOR %s", op.Factor, op.Action)
		if op.ChallengeResponse != "" {
			sb.WriteString(" :challenge ***")
		}
		sb.WriteString("}")
	}
	sb.WriteString("}")
}

func writeCreateRoleStmt(sb *strings.Builder, n *CreateRoleStmt) {
	sb.WriteString("{CREATE_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	sb.WriteString("}")
}

func writeDropRoleStmt(sb *strings.Builder, n *DropRoleStmt) {
	sb.WriteString("{DROP_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	sb.WriteString("}")
}

func writeSetDefaultRoleStmt(sb *strings.Builder, n *SetDefaultRoleStmt) {
	sb.WriteString("{SET_DEFAULT_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.None {
		sb.WriteString(" :none true")
	}
	if n.All {
		sb.WriteString(" :all true")
	}
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	if len(n.To) > 0 {
		fmt.Fprintf(sb, " :to %s", strings.Join(n.To, ", "))
	}
	sb.WriteString("}")
}

func writeSetRoleStmt(sb *strings.Builder, n *SetRoleStmt) {
	sb.WriteString("{SET_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Default {
		sb.WriteString(" :default true")
	}
	if n.None {
		sb.WriteString(" :none true")
	}
	if n.All {
		sb.WriteString(" :all true")
	}
	if len(n.AllExcept) > 0 {
		fmt.Fprintf(sb, " :all_except %s", strings.Join(n.AllExcept, ", "))
	}
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	sb.WriteString("}")
}

func writeGrantRoleStmt(sb *strings.Builder, n *GrantRoleStmt) {
	sb.WriteString("{GRANT_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	if len(n.To) > 0 {
		fmt.Fprintf(sb, " :to %s", strings.Join(n.To, ", "))
	}
	if n.WithAdmin {
		sb.WriteString(" :with_admin true")
	}
	sb.WriteString("}")
}

func writeRevokeRoleStmt(sb *strings.Builder, n *RevokeRoleStmt) {
	sb.WriteString("{REVOKE_ROLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if len(n.Roles) > 0 {
		fmt.Fprintf(sb, " :roles %s", strings.Join(n.Roles, ", "))
	}
	if len(n.From) > 0 {
		fmt.Fprintf(sb, " :from %s", strings.Join(n.From, ", "))
	}
	if n.IgnoreUnknownUser {
		sb.WriteString(" :ignore_unknown_user true")
	}
	sb.WriteString("}")
}

func writeCreateFunctionStmt(sb *strings.Builder, n *CreateFunctionStmt) {
	sb.WriteString("{CREATE_FUNCTION")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.IsProcedure {
		sb.WriteString(" :is_procedure true")
	}
	if n.IsAggregate {
		sb.WriteString(" :is_aggregate true")
	}
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Params) > 0 {
		sb.WriteString(" :params ")
		for i, p := range n.Params {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
	}
	if n.Returns != nil {
		sb.WriteString(" :returns ")
		writeNode(sb, n.Returns)
	}
	if n.Soname != "" {
		fmt.Fprintf(sb, " :soname %q", n.Soname)
	}
	if n.BodyText != "" {
		fmt.Fprintf(sb, " :body %q", n.BodyText)
	}
	if n.Body != nil {
		sb.WriteString(" :body-stmt ")
		writeNode(sb, n.Body)
	}
	if len(n.Characteristics) > 0 {
		sb.WriteString(" :characteristics ")
		for i, c := range n.Characteristics {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	sb.WriteString("}")
}

func writeCreateTriggerStmt(sb *strings.Builder, n *CreateTriggerStmt) {
	sb.WriteString("{CREATE_TRIGGER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Timing != "" {
		fmt.Fprintf(sb, " :timing %s", n.Timing)
	}
	if n.Event != "" {
		fmt.Fprintf(sb, " :event %s", n.Event)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Order != nil {
		sb.WriteString(" :order ")
		writeNode(sb, n.Order)
	}
	if n.BodyText != "" {
		fmt.Fprintf(sb, " :body %q", n.BodyText)
	}
	if n.Body != nil {
		sb.WriteString(" :body-stmt ")
		writeNode(sb, n.Body)
	}
	sb.WriteString("}")
}

func writeCreateEventStmt(sb *strings.Builder, n *CreateEventStmt) {
	sb.WriteString("{CREATE_EVENT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Schedule != nil {
		sb.WriteString(" :schedule ")
		writeNode(sb, n.Schedule)
	}
	if n.OnCompletion != "" {
		fmt.Fprintf(sb, " :on_completion %s", n.OnCompletion)
	}
	if n.Enable != "" {
		fmt.Fprintf(sb, " :enable %s", n.Enable)
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.BodyText != "" {
		fmt.Fprintf(sb, " :body %q", n.BodyText)
	}
	if n.Body != nil {
		sb.WriteString(" :body-stmt ")
		writeNode(sb, n.Body)
	}
	sb.WriteString("}")
}

func writeLoadDataStmt(sb *strings.Builder, n *LoadDataStmt) {
	if n.IsXML {
		sb.WriteString("{LOAD_XML")
	} else {
		sb.WriteString("{LOAD_DATA")
	}
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.LowPriority {
		sb.WriteString(" :low_priority true")
	}
	if n.Concurrent {
		sb.WriteString(" :concurrent true")
	}
	if n.Local {
		sb.WriteString(" :local true")
	}
	if n.Infile != "" {
		fmt.Fprintf(sb, " :infile %q", n.Infile)
	}
	if n.Replace {
		sb.WriteString(" :replace true")
	}
	if n.Ignore {
		sb.WriteString(" :ignore true")
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.Partitions) > 0 {
		fmt.Fprintf(sb, " :partitions (%s)", strings.Join(n.Partitions, ", "))
	}
	if n.CharacterSet != "" {
		fmt.Fprintf(sb, " :charset %s", n.CharacterSet)
	}
	if n.RowsIdentifiedBy != "" {
		fmt.Fprintf(sb, " :rows_identified_by %q", n.RowsIdentifiedBy)
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns ")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	if len(n.SetList) > 0 {
		sb.WriteString(" :set ")
		for i, a := range n.SetList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
	}
	if n.LinesStartingBy != "" {
		fmt.Fprintf(sb, " :lines_starting %q", n.LinesStartingBy)
	}
	if n.LinesTerminatedBy != "" {
		fmt.Fprintf(sb, " :lines_terminated %q", n.LinesTerminatedBy)
	}
	if n.FieldsTerminatedBy != "" {
		fmt.Fprintf(sb, " :fields_terminated %q", n.FieldsTerminatedBy)
	}
	if n.FieldsEnclosedBy != "" {
		fmt.Fprintf(sb, " :fields_enclosed %q", n.FieldsEnclosedBy)
	}
	if n.FieldsOptionalEncl {
		sb.WriteString(" :fields_optionally_enclosed true")
	}
	if n.FieldsEscapedBy != "" {
		fmt.Fprintf(sb, " :fields_escaped %q", n.FieldsEscapedBy)
	}
	if n.IgnoreRows > 0 {
		fmt.Fprintf(sb, " :ignore_rows %d", n.IgnoreRows)
	}
	sb.WriteString("}")
}

func writePrepareStmt(sb *strings.Builder, n *PrepareStmt) {
	sb.WriteString("{PREPARE")
	fmt.Fprintf(sb, " :loc %d :name %s :stmt %q", n.Loc.Start, n.Name, n.Stmt)
	sb.WriteString("}")
}

func writeExecuteStmt(sb *strings.Builder, n *ExecuteStmt) {
	sb.WriteString("{EXECUTE")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if len(n.Params) > 0 {
		sb.WriteString(" :params ")
		writeExprNodeList(sb, n.Params)
	}
	sb.WriteString("}")
}

func writeDeallocateStmt(sb *strings.Builder, n *DeallocateStmt) {
	sb.WriteString("{DEALLOCATE")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	sb.WriteString("}")
}

func writeAnalyzeTableStmt(sb *strings.Builder, n *AnalyzeTableStmt) {
	sb.WriteString("{ANALYZE_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.NoWriteToBinlog {
		sb.WriteString(" :no_write_to_binlog true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.HistogramOp != "" {
		fmt.Fprintf(sb, " :histogram_op %s", n.HistogramOp)
	}
	if len(n.HistogramColumns) > 0 {
		fmt.Fprintf(sb, " :histogram_columns %s", strings.Join(n.HistogramColumns, ", "))
	}
	if n.Buckets > 0 {
		fmt.Fprintf(sb, " :buckets %d", n.Buckets)
	}
	if n.UsingData != "" {
		fmt.Fprintf(sb, " :using_data \"%s\"", n.UsingData)
	}
	sb.WriteString("}")
}

func writeOptimizeTableStmt(sb *strings.Builder, n *OptimizeTableStmt) {
	sb.WriteString("{OPTIMIZE_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.NoWriteToBinlog {
		sb.WriteString(" :no_write_to_binlog true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	sb.WriteString("}")
}

func writeCheckTableStmt(sb *strings.Builder, n *CheckTableStmt) {
	sb.WriteString("{CHECK_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if len(n.Options) > 0 {
		fmt.Fprintf(sb, " :options %s", strings.Join(n.Options, ", "))
	}
	sb.WriteString("}")
}

func writeRepairTableStmt(sb *strings.Builder, n *RepairTableStmt) {
	sb.WriteString("{REPAIR_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.NoWriteToBinlog {
		sb.WriteString(" :no_write_to_binlog true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.Quick {
		sb.WriteString(" :quick true")
	}
	if n.Extended {
		sb.WriteString(" :extended true")
	}
	if n.UseFrm {
		sb.WriteString(" :use_frm true")
	}
	sb.WriteString("}")
}

func writeFlushStmt(sb *strings.Builder, n *FlushStmt) {
	sb.WriteString("{FLUSH")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.NoWriteToBinlog {
		sb.WriteString(" :no_write_to_binlog true")
	}
	if len(n.Options) > 0 {
		fmt.Fprintf(sb, " :options %s", strings.Join(n.Options, ", "))
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.WithReadLock {
		sb.WriteString(" :with_read_lock true")
	}
	if n.ForExport {
		sb.WriteString(" :for_export true")
	}
	if n.RelayChannel != "" {
		fmt.Fprintf(sb, " :relay_channel \"%s\"", n.RelayChannel)
	}
	sb.WriteString("}")
}

func writeKillStmt(sb *strings.Builder, n *KillStmt) {
	sb.WriteString("{KILL")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Query {
		sb.WriteString(" :query true")
	}
	if n.ConnectionID != nil {
		sb.WriteString(" :connection_id ")
		writeNode(sb, n.ConnectionID)
	}
	sb.WriteString("}")
}

func writeDoStmt(sb *strings.Builder, n *DoStmt) {
	sb.WriteString("{DO")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Exprs) > 0 {
		sb.WriteString(" :exprs ")
		writeExprNodeList(sb, n.Exprs)
	}
	sb.WriteString("}")
}

func writeChecksumTableStmt(sb *strings.Builder, n *ChecksumTableStmt) {
	sb.WriteString("{CHECKSUM_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.Quick {
		sb.WriteString(" :quick true")
	}
	if n.Extended {
		sb.WriteString(" :extended true")
	}
	sb.WriteString("}")
}

func writeShutdownStmt(sb *strings.Builder, n *ShutdownStmt) {
	sb.WriteString("{SHUTDOWN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	sb.WriteString("}")
}

func writeRestartStmt(sb *strings.Builder, n *RestartStmt) {
	sb.WriteString("{RESTART")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	sb.WriteString("}")
}

func writeAlterDatabaseStmt(sb *strings.Builder, n *AlterDatabaseStmt) {
	sb.WriteString("{ALTER_DATABASE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	sb.WriteString("}")
}

func writeDropDatabaseStmt(sb *strings.Builder, n *DropDatabaseStmt) {
	sb.WriteString("{DROP_DATABASE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	sb.WriteString("}")
}

func writeDropIndexStmt(sb *strings.Builder, n *DropIndexStmt) {
	sb.WriteString("{DROP_INDEX")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Algorithm != "" {
		fmt.Fprintf(sb, " :algorithm %s", n.Algorithm)
	}
	if n.Lock != "" {
		fmt.Fprintf(sb, " :lock %s", n.Lock)
	}
	sb.WriteString("}")
}

func writeDropViewStmt(sb *strings.Builder, n *DropViewStmt) {
	sb.WriteString("{DROP_VIEW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if len(n.Views) > 0 {
		sb.WriteString(" :views ")
		for i, v := range n.Views {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, v)
		}
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.Restrict {
		sb.WriteString(" :restrict true")
	}
	sb.WriteString("}")
}

func writeRawStmt(sb *strings.Builder, n *RawStmt) {
	sb.WriteString("{RAW_STMT")
	fmt.Fprintf(sb, " :loc %d :len %d :stmt ", n.Loc.Start, n.StmtLen)
	writeNode(sb, n.Stmt)
	sb.WriteString("}")
}

// -----------------------------------------------------------------------
// Expression writers
// -----------------------------------------------------------------------

func binaryOpStr(op BinaryOp) string {
	switch op {
	case BinOpAdd:
		return "+"
	case BinOpSub:
		return "-"
	case BinOpMul:
		return "*"
	case BinOpDiv:
		return "/"
	case BinOpMod:
		return "%%"
	case BinOpEq:
		return "="
	case BinOpNe:
		return "<>"
	case BinOpLt:
		return "<"
	case BinOpGt:
		return ">"
	case BinOpLe:
		return "<="
	case BinOpGe:
		return ">="
	case BinOpAnd:
		return "AND"
	case BinOpOr:
		return "OR"
	case BinOpXor:
		return "XOR"
	case BinOpBitAnd:
		return "&"
	case BinOpBitOr:
		return "|"
	case BinOpBitXor:
		return "^"
	case BinOpShiftLeft:
		return "<<"
	case BinOpShiftRight:
		return ">>"
	case BinOpDivInt:
		return "DIV"
	case BinOpRegexp:
		return "REGEXP"
	case BinOpLikeEscape:
		return "LIKE_ESCAPE"
	case BinOpNullSafeEq:
		return "<=>"
	case BinOpAssign:
		return ":="
	case BinOpJsonExtract:
		return "->"
	case BinOpJsonUnquote:
		return "->>"
	case BinOpSoundsLike:
		return "SOUNDS LIKE"
	default:
		return fmt.Sprintf("?%d", op)
	}
}

func unaryOpStr(op UnaryOp) string {
	switch op {
	case UnaryMinus:
		return "-"
	case UnaryPlus:
		return "+"
	case UnaryNot:
		return "NOT"
	case UnaryBitNot:
		return "~"
	case UnaryBinary:
		return "BINARY"
	default:
		return fmt.Sprintf("?%d", op)
	}
}

func isTestStr(t IsTestType) string {
	switch t {
	case IsNull:
		return "NULL"
	case IsTrue:
		return "TRUE"
	case IsFalse:
		return "FALSE"
	case IsUnknown:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("?%d", t)
	}
}

func writeBinaryExpr(sb *strings.Builder, n *BinaryExpr) {
	sb.WriteString("{BINEXPR")
	fmt.Fprintf(sb, " :op %s :loc %d :left ", binaryOpStr(n.Op), n.Loc.Start)
	writeNode(sb, n.Left)
	sb.WriteString(" :right ")
	writeNode(sb, n.Right)
	sb.WriteString("}")
}

func writeUnaryExpr(sb *strings.Builder, n *UnaryExpr) {
	sb.WriteString("{UNARY")
	fmt.Fprintf(sb, " :op %s :loc %d :operand ", unaryOpStr(n.Op), n.Loc.Start)
	writeNode(sb, n.Operand)
	sb.WriteString("}")
}

func writeFuncCallExpr(sb *strings.Builder, n *FuncCallExpr) {
	sb.WriteString("{FUNCCALL")
	fmt.Fprintf(sb, " :name %s :loc %d", n.Name, n.Loc.Start)
	if n.Schema != "" {
		fmt.Fprintf(sb, " :schema %s", n.Schema)
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.Star {
		sb.WriteString(" :star true")
	}
	if len(n.Args) > 0 {
		sb.WriteString(" :args ")
		writeExprNodeList(sb, n.Args)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
	}
	if n.Over != nil {
		sb.WriteString(" :over ")
		writeNode(sb, n.Over)
	}
	sb.WriteString("}")
}

func writeSubqueryExpr(sb *strings.Builder, n *SubqueryExpr) {
	sb.WriteString("{SUBQUERY")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Exists {
		sb.WriteString(" :exists true")
	}
	if n.Lateral {
		sb.WriteString(" :lateral true")
	}
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias %s", n.Alias)
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeCaseExpr(sb *strings.Builder, n *CaseExpr) {
	sb.WriteString("{CASE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Operand != nil {
		sb.WriteString(" :operand ")
		writeNode(sb, n.Operand)
	}
	if len(n.Whens) > 0 {
		sb.WriteString(" :whens ")
		for i, w := range n.Whens {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, w)
		}
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString("}")
}

func writeCaseWhen(sb *strings.Builder, n *CaseWhen) {
	sb.WriteString("{WHEN")
	fmt.Fprintf(sb, " :loc %d :cond ", n.Loc.Start)
	writeNode(sb, n.Cond)
	sb.WriteString(" :result ")
	writeNode(sb, n.Result)
	sb.WriteString("}")
}

func writeBetweenExpr(sb *strings.Builder, n *BetweenExpr) {
	sb.WriteString("{BETWEEN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(" :low ")
	writeNode(sb, n.Low)
	sb.WriteString(" :high ")
	writeNode(sb, n.High)
	sb.WriteString("}")
}

func writeInExpr(sb *strings.Builder, n *InExpr) {
	sb.WriteString("{IN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	if len(n.List) > 0 {
		sb.WriteString(" :list ")
		writeExprNodeList(sb, n.List)
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeLikeExpr(sb *strings.Builder, n *LikeExpr) {
	sb.WriteString("{LIKE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(" :pattern ")
	writeNode(sb, n.Pattern)
	if n.Escape != nil {
		sb.WriteString(" :escape ")
		writeNode(sb, n.Escape)
	}
	sb.WriteString("}")
}

func writeIsExpr(sb *strings.Builder, n *IsExpr) {
	sb.WriteString("{IS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Not {
		sb.WriteString(" :not true")
	}
	fmt.Fprintf(sb, " :test %s", isTestStr(n.Test))
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString("}")
}

func writeExistsExpr(sb *strings.Builder, n *ExistsExpr) {
	sb.WriteString("{EXISTS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeCastExpr(sb *strings.Builder, n *CastExpr) {
	sb.WriteString("{CAST")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	sb.WriteString(" :type ")
	writeNode(sb, n.TypeName)
	sb.WriteString("}")
}

func writeExtractExpr(sb *strings.Builder, n *ExtractExpr) {
	sb.WriteString("{EXTRACT")
	fmt.Fprintf(sb, " :loc %d :unit %s :expr ", n.Loc.Start, n.Unit)
	writeNode(sb, n.Expr)
	sb.WriteString("}")
}

func writeParenExpr(sb *strings.Builder, n *ParenExpr) {
	sb.WriteString("{PAREN")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	sb.WriteString("}")
}

func writeCollateExpr(sb *strings.Builder, n *CollateExpr) {
	sb.WriteString("{COLLATE")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	fmt.Fprintf(sb, " :collation %s", n.Collation)
	sb.WriteString("}")
}

func writeIntervalExpr(sb *strings.Builder, n *IntervalExpr) {
	sb.WriteString("{INTERVAL")
	fmt.Fprintf(sb, " :loc %d :val ", n.Loc.Start)
	writeNode(sb, n.Value)
	fmt.Fprintf(sb, " :unit %s", n.Unit)
	sb.WriteString("}")
}

func writeMatchExpr(sb *strings.Builder, n *MatchExpr) {
	sb.WriteString("{MATCH")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns ")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	sb.WriteString(" :against ")
	writeNode(sb, n.Against)
	if n.Modifier != "" {
		fmt.Fprintf(sb, " :modifier %s", n.Modifier)
	}
	sb.WriteString("}")
}

func writeConvertExpr(sb *strings.Builder, n *ConvertExpr) {
	sb.WriteString("{CONVERT")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	if n.TypeName != nil {
		sb.WriteString(" :type ")
		writeNode(sb, n.TypeName)
	}
	if n.Charset != "" {
		fmt.Fprintf(sb, " :charset %s", n.Charset)
	}
	sb.WriteString("}")
}

func writeDefaultExpr(sb *strings.Builder, n *DefaultExpr) {
	sb.WriteString("{DEFAULT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Column != "" {
		fmt.Fprintf(sb, " :col %s", n.Column)
	}
	sb.WriteString("}")
}

func writeRowExpr(sb *strings.Builder, n *RowExpr) {
	sb.WriteString("{ROW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Items) > 0 {
		sb.WriteString(" :items (")
		for i, item := range n.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

// -----------------------------------------------------------------------
// Reference and literal writers
// -----------------------------------------------------------------------

func writeColumnRef(sb *strings.Builder, n *ColumnRef) {
	sb.WriteString("{COLREF")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Schema != "" {
		fmt.Fprintf(sb, " :schema %s", n.Schema)
	}
	if n.Table != "" {
		fmt.Fprintf(sb, " :table %s", n.Table)
	}
	if n.Column != "" {
		fmt.Fprintf(sb, " :col %s", n.Column)
	}
	if n.Star {
		sb.WriteString(" :star true")
	}
	sb.WriteString("}")
}

func writeTableRef(sb *strings.Builder, n *TableRef) {
	sb.WriteString("{TABLEREF")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Schema != "" {
		fmt.Fprintf(sb, " :schema %s", n.Schema)
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias %s", n.Alias)
	}
	if len(n.Partitions) > 0 {
		sb.WriteString(" :partitions (")
		for i, p := range n.Partitions {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	if len(n.IndexHints) > 0 {
		sb.WriteString(" :index_hints ")
		for i, h := range n.IndexHints {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, h)
		}
	}
	sb.WriteString("}")
}

func writeStringLit(sb *strings.Builder, n *StringLit) {
	sb.WriteString("{STRING_LIT")
	fmt.Fprintf(sb, " :val %q :loc %d", n.Value, n.Loc.Start)
	if n.Charset != "" {
		fmt.Fprintf(sb, " :charset %s", n.Charset)
	}
	sb.WriteString("}")
}

func writeVariableRef(sb *strings.Builder, n *VariableRef) {
	sb.WriteString("{VAR")
	fmt.Fprintf(sb, " :name %s :loc %d", n.Name, n.Loc.Start)
	if n.System {
		sb.WriteString(" :system true")
	}
	if n.Scope != "" {
		fmt.Fprintf(sb, " :scope %s", n.Scope)
	}
	sb.WriteString("}")
}

// -----------------------------------------------------------------------
// Clause helper writers
// -----------------------------------------------------------------------

func writeAssignment(sb *strings.Builder, n *Assignment) {
	sb.WriteString("{ASSIGN")
	fmt.Fprintf(sb, " :loc %d :col ", n.Loc.Start)
	writeNode(sb, n.Column)
	sb.WriteString(" :val ")
	writeNode(sb, n.Value)
	sb.WriteString("}")
}

func writeOrderByItem(sb *strings.Builder, n *OrderByItem) {
	sb.WriteString("{ORDER_BY")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	if n.Desc {
		sb.WriteString(" :desc true")
	}
	if n.NullsFirst != nil {
		if *n.NullsFirst {
			sb.WriteString(" :nulls_first true")
		} else {
			sb.WriteString(" :nulls_last true")
		}
	}
	sb.WriteString("}")
}

func writeLimit(sb *strings.Builder, n *Limit) {
	sb.WriteString("{LIMIT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Count != nil {
		sb.WriteString(" :count ")
		writeNode(sb, n.Count)
	}
	if n.Offset != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.Offset)
	}
	sb.WriteString("}")
}

func writeJoinClause(sb *strings.Builder, n *JoinClause) {
	sb.WriteString("{JOIN")
	fmt.Fprintf(sb, " :type %d :loc %d :left ", n.Type, n.Loc.Start)
	writeNode(sb, n.Left)
	sb.WriteString(" :right ")
	writeNode(sb, n.Right)
	if n.Condition != nil {
		sb.WriteString(" :cond ")
		writeNode(sb, n.Condition)
	}
	sb.WriteString("}")
}

func writeOnCondition(sb *strings.Builder, n *OnCondition) {
	sb.WriteString("{ON")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	sb.WriteString("}")
}

func writeUsingCondition(sb *strings.Builder, n *UsingCondition) {
	sb.WriteString("{USING")
	fmt.Fprintf(sb, " :loc %d :columns %s", n.Loc.Start, strings.Join(n.Columns, ", "))
	sb.WriteString("}")
}

func writeResTarget(sb *strings.Builder, n *ResTarget) {
	sb.WriteString("{RES_TARGET")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	sb.WriteString(" :val ")
	writeNode(sb, n.Val)
	sb.WriteString("}")
}

func writeWindowDef(sb *strings.Builder, n *WindowDef) {
	sb.WriteString("{WINDOW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.RefName != "" {
		fmt.Fprintf(sb, " :ref %s", n.RefName)
	}
	if len(n.PartitionBy) > 0 {
		sb.WriteString(" :partition_by ")
		writeExprNodeList(sb, n.PartitionBy)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
	}
	if n.Frame != nil {
		sb.WriteString(" :frame ")
		writeNode(sb, n.Frame)
	}
	sb.WriteString("}")
}

func writeWindowFrame(sb *strings.Builder, n *WindowFrame) {
	sb.WriteString("{FRAME")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Start != nil {
		sb.WriteString(" :start ")
		writeNode(sb, n.Start)
	}
	if n.End != nil {
		sb.WriteString(" :end ")
		writeNode(sb, n.End)
	}
	sb.WriteString("}")
}

func writeWindowFrameBound(sb *strings.Builder, n *WindowFrameBound) {
	sb.WriteString("{FRAME_BOUND")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Offset != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.Offset)
	}
	sb.WriteString("}")
}

func writeColumnDef(sb *strings.Builder, n *ColumnDef) {
	sb.WriteString("{COLDEF")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if n.TypeName != nil {
		sb.WriteString(" :type ")
		writeNode(sb, n.TypeName)
	}
	if n.AutoIncrement {
		sb.WriteString(" :auto_increment true")
	}
	if n.DefaultValue != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.DefaultValue)
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.OnUpdate != nil {
		sb.WriteString(" :on_update ")
		writeNode(sb, n.OnUpdate)
	}
	if n.Generated != nil {
		sb.WriteString(" :generated ")
		writeNode(sb, n.Generated)
	}
	if n.ColumnFormat != "" {
		fmt.Fprintf(sb, " :column_format %s", n.ColumnFormat)
	}
	if n.Storage != "" {
		fmt.Fprintf(sb, " :storage %s", n.Storage)
	}
	if n.EngineAttribute != "" {
		fmt.Fprintf(sb, " :engine_attribute %q", n.EngineAttribute)
	}
	if n.SecondaryEngineAttribute != "" {
		fmt.Fprintf(sb, " :secondary_engine_attribute %q", n.SecondaryEngineAttribute)
	}
	if len(n.Constraints) > 0 {
		sb.WriteString(" :constraints ")
		for i, c := range n.Constraints {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	sb.WriteString("}")
}

func writeGeneratedColumn(sb *strings.Builder, n *GeneratedColumn) {
	sb.WriteString("{GENERATED")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Stored {
		sb.WriteString(" :stored true")
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString("}")
}

func writeColumnConstraint(sb *strings.Builder, n *ColumnConstraint) {
	sb.WriteString("{COL_CONSTRAINT")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.RefTable != nil {
		sb.WriteString(" :ref_table ")
		writeNode(sb, n.RefTable)
	}
	if len(n.RefColumns) > 0 {
		fmt.Fprintf(sb, " :ref_columns %s", strings.Join(n.RefColumns, ", "))
	}
	if n.Match != "" {
		fmt.Fprintf(sb, " :match %s", n.Match)
	}
	if n.OnDelete != RefActNone {
		fmt.Fprintf(sb, " :on_delete %d", n.OnDelete)
	}
	if n.OnUpdate != RefActNone {
		fmt.Fprintf(sb, " :on_update %d", n.OnUpdate)
	}
	if n.NotEnforced {
		sb.WriteString(" :not_enforced true")
	}
	sb.WriteString("}")
}

func writeConstraint(sb *strings.Builder, n *Constraint) {
	sb.WriteString("{CONSTRAINT")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if len(n.Columns) > 0 {
		fmt.Fprintf(sb, " :columns %s", strings.Join(n.Columns, ", "))
	}
	if len(n.IndexColumns) > 0 {
		sb.WriteString(" :index_columns (")
		for i, ic := range n.IndexColumns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, ic)
		}
		sb.WriteString(")")
	}
	if n.IndexType != "" {
		fmt.Fprintf(sb, " :index_type %s", n.IndexType)
	}
	if len(n.IndexOptions) > 0 {
		sb.WriteString(" :index_options (")
		for i, opt := range n.IndexOptions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, opt)
		}
		sb.WriteString(")")
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.RefTable != nil {
		sb.WriteString(" :ref_table ")
		writeNode(sb, n.RefTable)
	}
	if len(n.RefColumns) > 0 {
		fmt.Fprintf(sb, " :ref_columns %s", strings.Join(n.RefColumns, ", "))
	}
	if n.OnDelete != RefActNone {
		fmt.Fprintf(sb, " :on_delete %d", n.OnDelete)
	}
	if n.OnUpdate != RefActNone {
		fmt.Fprintf(sb, " :on_update %d", n.OnUpdate)
	}
	if n.Match != "" {
		fmt.Fprintf(sb, " :match %s", n.Match)
	}
	if n.NotEnforced {
		sb.WriteString(" :not_enforced true")
	}
	sb.WriteString("}")
}

func writeDataType(sb *strings.Builder, n *DataType) {
	sb.WriteString("{DATATYPE")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if n.Length > 0 {
		fmt.Fprintf(sb, " :len %d", n.Length)
	}
	if n.Scale > 0 {
		fmt.Fprintf(sb, " :scale %d", n.Scale)
	}
	if n.Unsigned {
		sb.WriteString(" :unsigned true")
	}
	if n.Zerofill {
		sb.WriteString(" :zerofill true")
	}
	if n.Charset != "" {
		fmt.Fprintf(sb, " :charset %s", n.Charset)
	}
	if n.Collate != "" {
		fmt.Fprintf(sb, " :collate %s", n.Collate)
	}
	if len(n.EnumValues) > 0 {
		fmt.Fprintf(sb, " :enum_values %s", strings.Join(n.EnumValues, ", "))
	}
	if n.ArrayDim > 0 {
		fmt.Fprintf(sb, " :array_dim %d", n.ArrayDim)
	}
	sb.WriteString("}")
}

func writeAlterTableCmd(sb *strings.Builder, n *AlterTableCmd) {
	sb.WriteString("{AT_CMD")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.NewName != "" {
		fmt.Fprintf(sb, " :new_name %s", n.NewName)
	}
	if n.Column != nil {
		sb.WriteString(" :col ")
		writeNode(sb, n.Column)
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns (")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
		sb.WriteString(")")
	}
	if n.DefaultExpr != nil {
		sb.WriteString(" :default_expr ")
		writeNode(sb, n.DefaultExpr)
	}
	if n.Constraint != nil {
		sb.WriteString(" :constraint ")
		writeNode(sb, n.Constraint)
	}
	if n.Option != nil {
		sb.WriteString(" :option ")
		writeNode(sb, n.Option)
	}
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.First {
		sb.WriteString(" :first true")
	}
	if n.After != "" {
		fmt.Fprintf(sb, " :after %s", n.After)
	}
	if len(n.PartitionNames) > 0 {
		sb.WriteString(" :partition_names (")
		for i, name := range n.PartitionNames {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(name)
		}
		sb.WriteString(")")
	}
	if n.AllPartitions {
		sb.WriteString(" :all_partitions true")
	}
	if len(n.PartitionDefs) > 0 {
		sb.WriteString(" :partition_defs (")
		for i, pd := range n.PartitionDefs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, pd)
		}
		sb.WriteString(")")
	}
	if n.Number > 0 {
		fmt.Fprintf(sb, " :number %d", n.Number)
	}
	if n.ExchangeTable != nil {
		sb.WriteString(" :exchange_table ")
		writeNode(sb, n.ExchangeTable)
	}
	if n.WithValidation != nil {
		if *n.WithValidation {
			sb.WriteString(" :validation with")
		} else {
			sb.WriteString(" :validation without")
		}
	}
	if len(n.OrderByItems) > 0 {
		sb.WriteString(" :order_by (")
		for i, item := range n.OrderByItems {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString(")")
	}
	if n.PartitionBy != nil {
		sb.WriteString(" :partition_by ")
		writeNode(sb, n.PartitionBy)
	}
	sb.WriteString("}")
}

func writePartitionClause(sb *strings.Builder, n *PartitionClause) {
	sb.WriteString("{PARTITION_CLAUSE")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Linear {
		sb.WriteString(" :linear true")
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if len(n.Columns) > 0 {
		fmt.Fprintf(sb, " :columns %s", strings.Join(n.Columns, ", "))
	}
	if n.Algorithm > 0 {
		fmt.Fprintf(sb, " :algorithm %d", n.Algorithm)
	}
	if n.NumParts > 0 {
		fmt.Fprintf(sb, " :num_parts %d", n.NumParts)
	}
	if len(n.Partitions) > 0 {
		sb.WriteString(" :partitions ")
		for i, p := range n.Partitions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
	}
	if n.SubPartType > 0 {
		fmt.Fprintf(sb, " :sub_type %d", n.SubPartType)
	}
	if n.SubPartExpr != nil {
		sb.WriteString(" :sub_expr ")
		writeNode(sb, n.SubPartExpr)
	}
	if len(n.SubPartColumns) > 0 {
		fmt.Fprintf(sb, " :sub_columns %s", strings.Join(n.SubPartColumns, ", "))
	}
	if n.SubPartAlgo > 0 {
		fmt.Fprintf(sb, " :sub_algorithm %d", n.SubPartAlgo)
	}
	if n.NumSubParts > 0 {
		fmt.Fprintf(sb, " :num_sub_parts %d", n.NumSubParts)
	}
	sb.WriteString("}")
}

func writePartitionDef(sb *strings.Builder, n *PartitionDef) {
	sb.WriteString("{PARTITION_DEF")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	if len(n.SubPartitions) > 0 {
		sb.WriteString(" :sub_partitions (")
		for i, sp := range n.SubPartitions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, sp)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeSubPartitionDef(sb *strings.Builder, n *SubPartitionDef) {
	sb.WriteString("{SUBPARTITION_DEF")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if len(n.Options) > 0 {
		sb.WriteString(" :options ")
		for i, o := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
	}
	sb.WriteString("}")
}

func writeIndexColumn(sb *strings.Builder, n *IndexColumn) {
	sb.WriteString("{IDX_COL")
	fmt.Fprintf(sb, " :loc %d :expr ", n.Loc.Start)
	writeNode(sb, n.Expr)
	if n.Length > 0 {
		fmt.Fprintf(sb, " :len %d", n.Length)
	}
	if n.Desc {
		sb.WriteString(" :desc true")
	}
	sb.WriteString("}")
}

func writeIndexOption(sb *strings.Builder, n *IndexOption) {
	sb.WriteString("{IDX_OPT")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if n.Value != nil {
		sb.WriteString(" :val ")
		writeNode(sb, n.Value)
	}
	sb.WriteString("}")
}

func writeIndexHint(sb *strings.Builder, n *IndexHint) {
	sb.WriteString("{IDX_HINT")
	fmt.Fprintf(sb, " :loc %d :type %d :scope %d", n.Loc.Start, n.Type, n.Scope)
	if len(n.Indexes) > 0 {
		fmt.Fprintf(sb, " :indexes %s", strings.Join(n.Indexes, ", "))
	}
	sb.WriteString("}")
}

func writeForUpdate(sb *strings.Builder, n *ForUpdate) {
	sb.WriteString("{FOR_UPDATE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Share {
		sb.WriteString(" :share true")
	}
	if n.LockInShareMode {
		sb.WriteString(" :lock_in_share_mode true")
	}
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables ")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
	}
	if n.NoWait {
		sb.WriteString(" :nowait true")
	}
	if n.SkipLocked {
		sb.WriteString(" :skip_locked true")
	}
	sb.WriteString("}")
}

func writeIntoClause(sb *strings.Builder, n *IntoClause) {
	sb.WriteString("{INTO")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Outfile != "" {
		fmt.Fprintf(sb, " :outfile %q", n.Outfile)
	}
	if n.Dumpfile != "" {
		fmt.Fprintf(sb, " :dumpfile %q", n.Dumpfile)
	}
	if n.Charset != "" {
		fmt.Fprintf(sb, " :charset %q", n.Charset)
	}
	if n.HasFieldsClause {
		sb.WriteString(" :fields true")
		if n.FieldsTerminatedBy != "" {
			fmt.Fprintf(sb, " :fields_terminated_by %q", n.FieldsTerminatedBy)
		}
		if n.FieldsEnclosedBy != "" {
			if n.FieldsOptionalEncl {
				sb.WriteString(" :optionally true")
			}
			fmt.Fprintf(sb, " :fields_enclosed_by %q", n.FieldsEnclosedBy)
		}
		if n.FieldsEscapedBy != "" {
			fmt.Fprintf(sb, " :fields_escaped_by %q", n.FieldsEscapedBy)
		}
	}
	if n.HasLinesClause {
		sb.WriteString(" :lines true")
		if n.LinesStartingBy != "" {
			fmt.Fprintf(sb, " :lines_starting_by %q", n.LinesStartingBy)
		}
		if n.LinesTerminatedBy != "" {
			fmt.Fprintf(sb, " :lines_terminated_by %q", n.LinesTerminatedBy)
		}
	}
	if len(n.Vars) > 0 {
		sb.WriteString(" :vars ")
		for i, v := range n.Vars {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, v)
		}
	}
	sb.WriteString("}")
}

func writeRenameTablePair(sb *strings.Builder, n *RenameTablePair) {
	sb.WriteString("{RENAME_PAIR")
	fmt.Fprintf(sb, " :loc %d :old ", n.Loc.Start)
	writeNode(sb, n.Old)
	sb.WriteString(" :new ")
	writeNode(sb, n.New)
	sb.WriteString("}")
}

func writeLockTable(sb *strings.Builder, n *LockTable) {
	sb.WriteString("{LOCK_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.LockType != "" {
		fmt.Fprintf(sb, " :lock_type %s", n.LockType)
	}
	sb.WriteString("}")
}

func writeGrantTarget(sb *strings.Builder, n *GrantTarget) {
	sb.WriteString("{GRANT_TARGET")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Type != "" {
		fmt.Fprintf(sb, " :type %s", n.Type)
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString("}")
}

func writeRequireClause(sb *strings.Builder, n *RequireClause) {
	sb.WriteString("{REQUIRE")
	if n.None {
		sb.WriteString(" NONE")
	}
	if n.SSL {
		sb.WriteString(" SSL")
	}
	if n.X509 {
		sb.WriteString(" X509")
	}
	if n.Cipher != "" {
		fmt.Fprintf(sb, " :cipher %q", n.Cipher)
	}
	if n.Issuer != "" {
		fmt.Fprintf(sb, " :issuer %q", n.Issuer)
	}
	if n.Subject != "" {
		fmt.Fprintf(sb, " :subject %q", n.Subject)
	}
	sb.WriteString("}")
}

func writeResourceOption(sb *strings.Builder, n *ResourceOption) {
	sb.WriteString("{RESOURCE")
	if n.HasMaxQueries {
		fmt.Fprintf(sb, " :max_queries_per_hour %d", n.MaxQueriesPerHour)
	}
	if n.HasMaxUpdates {
		fmt.Fprintf(sb, " :max_updates_per_hour %d", n.MaxUpdatesPerHour)
	}
	if n.HasMaxConnections {
		fmt.Fprintf(sb, " :max_connections_per_hour %d", n.MaxConnectionsPerHour)
	}
	if n.HasMaxUserConn {
		fmt.Fprintf(sb, " :max_user_connections %d", n.MaxUserConnections)
	}
	sb.WriteString("}")
}

func writeUserSpec(sb *strings.Builder, n *UserSpec) {
	sb.WriteString("{USER_SPEC")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Host != "" {
		fmt.Fprintf(sb, " :host %s", n.Host)
	}
	if n.AuthPlugin != "" {
		fmt.Fprintf(sb, " :auth_plugin %s", n.AuthPlugin)
	}
	if n.Password != "" {
		fmt.Fprintf(sb, " :password %s", n.Password)
	}
	if n.AuthHash != "" {
		fmt.Fprintf(sb, " :auth_hash %s", n.AuthHash)
	}
	if n.PasswordRandom {
		sb.WriteString(" :password_random true")
	}
	if n.RetainCurrentPassword {
		sb.WriteString(" :retain_current_password true")
	}
	if n.DiscardOldPassword {
		sb.WriteString(" :discard_old_password true")
	}
	if n.Replace != "" {
		fmt.Fprintf(sb, " :replace %s", n.Replace)
	}
	if len(n.AuthFactors) > 0 {
		sb.WriteString(" :auth_factors (")
		for i, af := range n.AuthFactors {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, af)
		}
		sb.WriteString(")")
	}
	if n.InitialAuthPlugin != "" {
		fmt.Fprintf(sb, " :initial_auth_plugin %s", n.InitialAuthPlugin)
	}
	if n.InitialAuthString != "" {
		fmt.Fprintf(sb, " :initial_auth_string %s", n.InitialAuthString)
	}
	if n.InitialAuthRandom {
		sb.WriteString(" :initial_auth_random true")
	}
	sb.WriteString("}")
}

func writeFuncParam(sb *strings.Builder, n *FuncParam) {
	sb.WriteString("{FUNC_PARAM")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Direction != "" {
		fmt.Fprintf(sb, " :direction %s", n.Direction)
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.TypeName != nil {
		sb.WriteString(" :type ")
		writeNode(sb, n.TypeName)
	}
	sb.WriteString("}")
}

func writeRoutineCharacteristic(sb *strings.Builder, n *RoutineCharacteristic) {
	sb.WriteString("{ROUTINE_CHAR")
	fmt.Fprintf(sb, " :loc %d :name %s :val %s", n.Loc.Start, n.Name, n.Value)
	sb.WriteString("}")
}

func writeTriggerOrder(sb *strings.Builder, n *TriggerOrder) {
	sb.WriteString("{TRIGGER_ORDER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Follows {
		sb.WriteString(" :follows true")
	}
	if n.TriggerName != "" {
		fmt.Fprintf(sb, " :trigger_name %s", n.TriggerName)
	}
	sb.WriteString("}")
}

func writeEventSchedule(sb *strings.Builder, n *EventSchedule) {
	sb.WriteString("{EVENT_SCHEDULE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.At != nil {
		sb.WriteString(" :at ")
		writeNode(sb, n.At)
	}
	if n.Every != nil {
		sb.WriteString(" :every ")
		writeNode(sb, n.Every)
	}
	if n.Starts != nil {
		sb.WriteString(" :starts ")
		writeNode(sb, n.Starts)
	}
	if n.Ends != nil {
		sb.WriteString(" :ends ")
		writeNode(sb, n.Ends)
	}
	sb.WriteString("}")
}

func writeCommonTableExpr(sb *strings.Builder, n *CommonTableExpr) {
	sb.WriteString("{CTE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Recursive {
		sb.WriteString(" :recursive true")
	}
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns (")
		for i, col := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%s", col)
		}
		sb.WriteString(")")
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeSetTransactionStmt(sb *strings.Builder, n *SetTransactionStmt) {
	sb.WriteString("{SET_TRANSACTION")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Scope != "" {
		fmt.Fprintf(sb, " :scope %s", n.Scope)
	}
	if n.IsolationLevel != "" {
		fmt.Fprintf(sb, " :isolation %s", n.IsolationLevel)
	}
	if n.AccessMode != "" {
		fmt.Fprintf(sb, " :access %s", n.AccessMode)
	}
	sb.WriteString("}")
}

func writeXAStmt(sb *strings.Builder, n *XAStmt) {
	sb.WriteString("{XA")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	switch n.Type {
	case XAStart:
		sb.WriteString(" :type START")
	case XAEnd:
		sb.WriteString(" :type END")
	case XAPrepare:
		sb.WriteString(" :type PREPARE")
	case XACommit:
		sb.WriteString(" :type COMMIT")
	case XARollback:
		sb.WriteString(" :type ROLLBACK")
	case XARecover:
		sb.WriteString(" :type RECOVER")
	}
	if len(n.Xid) > 0 {
		sb.WriteString(" :xid ")
		writeExprNodeList(sb, n.Xid)
	}
	if n.Join {
		sb.WriteString(" :join true")
	}
	if n.Resume {
		sb.WriteString(" :resume true")
	}
	if n.Suspend {
		sb.WriteString(" :suspend true")
	}
	if n.Migrate {
		sb.WriteString(" :migrate true")
	}
	if n.OnePhase {
		sb.WriteString(" :one_phase true")
	}
	if n.Convert {
		sb.WriteString(" :convert true")
	}
	sb.WriteString("}")
}

func writeMemberOfExpr(sb *strings.Builder, n *MemberOfExpr) {
	sb.WriteString("{MEMBER_OF")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	sb.WriteString(" :value ")
	writeNode(sb, n.Value)
	sb.WriteString(" :array ")
	writeNode(sb, n.Array)
	sb.WriteString("}")
}

func writeJsonTableExpr(sb *strings.Builder, n *JsonTableExpr) {
	sb.WriteString("{JSON_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(" :path ")
	writeNode(sb, n.Path)
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns (")
		for i, col := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, col)
		}
		sb.WriteString(")")
	}
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias %s", n.Alias)
	}
	sb.WriteString("}")
}

func writeJsonTableColumn(sb *strings.Builder, n *JsonTableColumn) {
	sb.WriteString("{JT_COL")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Ordinality {
		fmt.Fprintf(sb, " :name %s :ordinality true", n.Name)
	} else if n.Nested {
		fmt.Fprintf(sb, " :nested_path %q", n.NestedPath)
		if len(n.NestedCols) > 0 {
			sb.WriteString(" :columns (")
			for i, col := range n.NestedCols {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, col)
			}
			sb.WriteString(")")
		}
	} else {
		fmt.Fprintf(sb, " :name %s", n.Name)
		if n.TypeName != nil {
			sb.WriteString(" :type ")
			writeNode(sb, n.TypeName)
		}
		if n.Path != "" {
			fmt.Fprintf(sb, " :path %q", n.Path)
		}
		if n.Exists {
			sb.WriteString(" :exists true")
		}
	}
	sb.WriteString("}")
}

// -----------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------

func writeNodeList(sb *strings.Builder, nodes []Node) {
	sb.WriteString("(")
	for i, n := range nodes {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, n)
	}
	sb.WriteString(")")
}

func writeExprNodeList(sb *strings.Builder, nodes []ExprNode) {
	sb.WriteString("(")
	for i, n := range nodes {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, n)
	}
	sb.WriteString(")")
}

func writeTableExprList(sb *strings.Builder, nodes []TableExpr) {
	sb.WriteString("(")
	for i, n := range nodes {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, n)
	}
	sb.WriteString(")")
}

func writeStmtNodeList(sb *strings.Builder, nodes []StmtNode) {
	sb.WriteString("(")
	for i, n := range nodes {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, n)
	}
	sb.WriteString(")")
}

func writeOrderByList(sb *strings.Builder, items []*OrderByItem) {
	sb.WriteString("(")
	for i, item := range items {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, item)
	}
	sb.WriteString(")")
}

func writeCallStmt(sb *strings.Builder, n *CallStmt) {
	sb.WriteString("{CALL")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Args) > 0 {
		sb.WriteString(" :args ")
		writeExprNodeList(sb, n.Args)
	}
	sb.WriteString("}")
}

func writeHandlerOpenStmt(sb *strings.Builder, n *HandlerOpenStmt) {
	sb.WriteString("{HANDLER_OPEN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias %s", n.Alias)
	}
	sb.WriteString("}")
}

func writeHandlerReadStmt(sb *strings.Builder, n *HandlerReadStmt) {
	sb.WriteString("{HANDLER_READ")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Direction != "" {
		fmt.Fprintf(sb, " :direction %s", n.Direction)
	}
	if n.Index != "" {
		fmt.Fprintf(sb, " :index %s", n.Index)
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	sb.WriteString("}")
}

func writeHandlerCloseStmt(sb *strings.Builder, n *HandlerCloseStmt) {
	sb.WriteString("{HANDLER_CLOSE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	sb.WriteString("}")
}

func writeSignalStmt(sb *strings.Builder, n *SignalStmt) {
	sb.WriteString("{SIGNAL")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.ConditionValue != "" {
		fmt.Fprintf(sb, " :condition %s", n.ConditionValue)
	}
	if len(n.SetItems) > 0 {
		sb.WriteString(" :set ")
		for i, item := range n.SetItems {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
	}
	sb.WriteString("}")
}

func writeResignalStmt(sb *strings.Builder, n *ResignalStmt) {
	sb.WriteString("{RESIGNAL")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.ConditionValue != "" {
		fmt.Fprintf(sb, " :condition %s", n.ConditionValue)
	}
	if len(n.SetItems) > 0 {
		sb.WriteString(" :set ")
		for i, item := range n.SetItems {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
	}
	sb.WriteString("}")
}

func writeSignalInfoItem(sb *strings.Builder, n *SignalInfoItem) {
	sb.WriteString("{SIGNAL_INFO")
	fmt.Fprintf(sb, " :loc %d :name %s :val ", n.Loc.Start, n.Name)
	writeNode(sb, n.Value)
	sb.WriteString("}")
}

func writeGetDiagnosticsStmt(sb *strings.Builder, n *GetDiagnosticsStmt) {
	sb.WriteString("{GET_DIAGNOSTICS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Stacked {
		sb.WriteString(" :stacked true")
	}
	if n.StatementInfo {
		sb.WriteString(" :stmt_info true")
	}
	if n.ConditionNumber != nil {
		sb.WriteString(" :condition_number ")
		writeNode(sb, n.ConditionNumber)
	}
	if len(n.Items) > 0 {
		sb.WriteString(" :items ")
		for i, item := range n.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
	}
	sb.WriteString("}")
}

func writeDiagnosticsItem(sb *strings.Builder, n *DiagnosticsItem) {
	sb.WriteString("{DIAG_ITEM")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	sb.WriteString("}")
}

func writeBeginEndBlock(sb *strings.Builder, n *BeginEndBlock) {
	sb.WriteString("{BEGIN_END")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Label != "" {
		fmt.Fprintf(sb, " :label %s", n.Label)
	}
	if n.EndLabel != "" {
		fmt.Fprintf(sb, " :end_label %s", n.EndLabel)
	}
	if len(n.Stmts) > 0 {
		sb.WriteString(" :stmts")
		for _, s := range n.Stmts {
			sb.WriteString(" ")
			writeNode(sb, s)
		}
	}
	sb.WriteString("}")
}

func writeDeclareVarStmt(sb *strings.Builder, n *DeclareVarStmt) {
	sb.WriteString("{DECLARE_VAR")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Names) > 0 {
		sb.WriteString(" :names")
		for _, name := range n.Names {
			fmt.Fprintf(sb, " %s", name)
		}
	}
	if n.TypeName != nil {
		sb.WriteString(" :type ")
		writeNode(sb, n.TypeName)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString("}")
}

func writeDeclareConditionStmt(sb *strings.Builder, n *DeclareConditionStmt) {
	sb.WriteString("{DECLARE_CONDITION")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	fmt.Fprintf(sb, " :value %s", n.ConditionValue)
	sb.WriteString("}")
}

func writeDeclareHandlerStmt(sb *strings.Builder, n *DeclareHandlerStmt) {
	sb.WriteString("{DECLARE_HANDLER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :action %s", n.Action)
	if len(n.Conditions) > 0 {
		sb.WriteString(" :conditions")
		for _, c := range n.Conditions {
			fmt.Fprintf(sb, " %s", c)
		}
	}
	if n.Stmt != nil {
		sb.WriteString(" :stmt ")
		writeNode(sb, n.Stmt)
	}
	sb.WriteString("}")
}

func writeDeclareCursorStmt(sb *strings.Builder, n *DeclareCursorStmt) {
	sb.WriteString("{DECLARE_CURSOR")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	sb.WriteString("}")
}

func writeIfStmt(sb *strings.Builder, n *IfStmt) {
	sb.WriteString("{IF_STMT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Cond != nil {
		sb.WriteString(" :cond ")
		writeNode(sb, n.Cond)
	}
	if len(n.ThenList) > 0 {
		sb.WriteString(" :then (")
		for i, s := range n.ThenList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	if len(n.ElseIfs) > 0 {
		sb.WriteString(" :elseifs (")
		for i, ei := range n.ElseIfs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, ei)
		}
		sb.WriteString(")")
	}
	if len(n.ElseList) > 0 {
		sb.WriteString(" :else (")
		for i, s := range n.ElseList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeElseIf(sb *strings.Builder, n *ElseIf) {
	sb.WriteString("{ELSEIF")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Cond != nil {
		sb.WriteString(" :cond ")
		writeNode(sb, n.Cond)
	}
	if len(n.ThenList) > 0 {
		sb.WriteString(" :then (")
		for i, s := range n.ThenList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeCaseStmtNode(sb *strings.Builder, n *CaseStmtNode) {
	sb.WriteString("{CASE_STMT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Operand != nil {
		sb.WriteString(" :operand ")
		writeNode(sb, n.Operand)
	}
	if len(n.Whens) > 0 {
		sb.WriteString(" :whens (")
		for i, w := range n.Whens {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, w)
		}
		sb.WriteString(")")
	}
	if len(n.ElseList) > 0 {
		sb.WriteString(" :else (")
		for i, s := range n.ElseList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeCaseStmtWhen(sb *strings.Builder, n *CaseStmtWhen) {
	sb.WriteString("{WHEN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Cond != nil {
		sb.WriteString(" :cond ")
		writeNode(sb, n.Cond)
	}
	if len(n.ThenList) > 0 {
		sb.WriteString(" :then (")
		for i, s := range n.ThenList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeWhileStmt(sb *strings.Builder, n *WhileStmt) {
	sb.WriteString("{WHILE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Label != "" {
		fmt.Fprintf(sb, " :label %s", n.Label)
	}
	if n.EndLabel != "" {
		fmt.Fprintf(sb, " :end_label %s", n.EndLabel)
	}
	if n.Cond != nil {
		sb.WriteString(" :cond ")
		writeNode(sb, n.Cond)
	}
	if len(n.Stmts) > 0 {
		sb.WriteString(" :stmts (")
		for i, s := range n.Stmts {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeRepeatStmt(sb *strings.Builder, n *RepeatStmt) {
	sb.WriteString("{REPEAT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Label != "" {
		fmt.Fprintf(sb, " :label %s", n.Label)
	}
	if n.EndLabel != "" {
		fmt.Fprintf(sb, " :end_label %s", n.EndLabel)
	}
	if len(n.Stmts) > 0 {
		sb.WriteString(" :stmts (")
		for i, s := range n.Stmts {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	if n.Cond != nil {
		sb.WriteString(" :until ")
		writeNode(sb, n.Cond)
	}
	sb.WriteString("}")
}

func writeLoopStmt(sb *strings.Builder, n *LoopStmt) {
	sb.WriteString("{LOOP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Label != "" {
		fmt.Fprintf(sb, " :label %s", n.Label)
	}
	if n.EndLabel != "" {
		fmt.Fprintf(sb, " :end_label %s", n.EndLabel)
	}
	if len(n.Stmts) > 0 {
		sb.WriteString(" :stmts (")
		for i, s := range n.Stmts {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeLeaveStmt(sb *strings.Builder, n *LeaveStmt) {
	sb.WriteString("{LEAVE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :label %s", n.Label)
	sb.WriteString("}")
}

func writeIterateStmt(sb *strings.Builder, n *IterateStmt) {
	sb.WriteString("{ITERATE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :label %s", n.Label)
	sb.WriteString("}")
}

func writeReturnStmt(sb *strings.Builder, n *ReturnStmt) {
	sb.WriteString("{RETURN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString("}")
}

func writeOpenCursorStmt(sb *strings.Builder, n *OpenCursorStmt) {
	sb.WriteString("{OPEN_CURSOR")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	sb.WriteString("}")
}

func writeFetchCursorStmt(sb *strings.Builder, n *FetchCursorStmt) {
	sb.WriteString("{FETCH_CURSOR")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if len(n.Into) > 0 {
		sb.WriteString(" :into")
		for _, v := range n.Into {
			fmt.Fprintf(sb, " %s", v)
		}
	}
	sb.WriteString("}")
}

func writeCloseCursorStmt(sb *strings.Builder, n *CloseCursorStmt) {
	sb.WriteString("{CLOSE_CURSOR")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	sb.WriteString("}")
}

func writeInstallPluginStmt(sb *strings.Builder, n *InstallPluginStmt) {
	sb.WriteString("{INSTALL_PLUGIN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.PluginName)
	if n.Soname != "" {
		fmt.Fprintf(sb, " :soname %q", n.Soname)
	}
	sb.WriteString("}")
}

func writeUninstallPluginStmt(sb *strings.Builder, n *UninstallPluginStmt) {
	sb.WriteString("{UNINSTALL_PLUGIN")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.PluginName)
	sb.WriteString("}")
}

func writeInstallComponentStmt(sb *strings.Builder, n *InstallComponentStmt) {
	sb.WriteString("{INSTALL_COMPONENT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Components) > 0 {
		sb.WriteString(" :components")
		for _, c := range n.Components {
			fmt.Fprintf(sb, " %q", c)
		}
	}
	if len(n.SetVars) > 0 {
		sb.WriteString(" :set_vars")
		for _, v := range n.SetVars {
			sb.WriteString(" ")
			writeNode(sb, v)
		}
	}
	sb.WriteString("}")
}

func writeUninstallComponentStmt(sb *strings.Builder, n *UninstallComponentStmt) {
	sb.WriteString("{UNINSTALL_COMPONENT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Components) > 0 {
		sb.WriteString(" :components")
		for _, c := range n.Components {
			fmt.Fprintf(sb, " %q", c)
		}
	}
	sb.WriteString("}")
}

func writeTableStmt(sb *strings.Builder, n *TableStmt) {
	sb.WriteString("{TABLE_STMT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by (")
		for i, item := range n.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString(")")
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	sb.WriteString("}")
}

func writeValuesStmt(sb *strings.Builder, n *ValuesStmt) {
	sb.WriteString("{VALUES_STMT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Rows) > 0 {
		sb.WriteString(" :rows (")
		for i, row := range n.Rows {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString("(")
			for j, expr := range row {
				if j > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, expr)
			}
			sb.WriteString(")")
		}
		sb.WriteString(")")
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by (")
		for i, item := range n.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString(")")
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	sb.WriteString("}")
}

func writeCloneStmt(sb *strings.Builder, n *CloneStmt) {
	sb.WriteString("{CLONE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Local {
		sb.WriteString(" :local true")
	}
	if n.Directory != "" {
		fmt.Fprintf(sb, " :directory %q", n.Directory)
	}
	if n.User != "" {
		fmt.Fprintf(sb, " :user %q", n.User)
	}
	if n.Host != "" {
		fmt.Fprintf(sb, " :host %q", n.Host)
	}
	if n.Port != 0 {
		fmt.Fprintf(sb, " :port %d", n.Port)
	}
	if n.Password != "" {
		fmt.Fprintf(sb, " :password %q", n.Password)
	}
	if n.RequireSSL != nil {
		if *n.RequireSSL {
			sb.WriteString(" :require_ssl true")
		} else {
			sb.WriteString(" :require_ssl false")
		}
	}
	sb.WriteString("}")
}

func writeCreateTablespaceStmt(sb *strings.Builder, n *CreateTablespaceStmt) {
	sb.WriteString("{CREATE_TABLESPACE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Undo {
		sb.WriteString(" :undo true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.DataFile != "" {
		fmt.Fprintf(sb, " :datafile %q", n.DataFile)
	}
	if n.AutoextendSize != "" {
		fmt.Fprintf(sb, " :autoextend_size %s", n.AutoextendSize)
	}
	if n.FileBlockSize != "" {
		fmt.Fprintf(sb, " :file_block_size %s", n.FileBlockSize)
	}
	if n.Encryption != "" {
		fmt.Fprintf(sb, " :encryption %q", n.Encryption)
	}
	if n.UseLogfileGroup != "" {
		fmt.Fprintf(sb, " :use_logfile_group %s", n.UseLogfileGroup)
	}
	if n.ExtentSize != "" {
		fmt.Fprintf(sb, " :extent_size %s", n.ExtentSize)
	}
	if n.InitialSize != "" {
		fmt.Fprintf(sb, " :initial_size %s", n.InitialSize)
	}
	if n.MaxSize != "" {
		fmt.Fprintf(sb, " :max_size %s", n.MaxSize)
	}
	if n.NodeGroup != "" {
		fmt.Fprintf(sb, " :nodegroup %s", n.NodeGroup)
	}
	if n.Wait {
		sb.WriteString(" :wait true")
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	if n.EngineAttribute != "" {
		fmt.Fprintf(sb, " :engine_attribute %q", n.EngineAttribute)
	}
	sb.WriteString("}")
}

func writeAlterTablespaceStmt(sb *strings.Builder, n *AlterTablespaceStmt) {
	sb.WriteString("{ALTER_TABLESPACE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Undo {
		sb.WriteString(" :undo true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.AddDataFile != "" {
		fmt.Fprintf(sb, " :add_datafile %q", n.AddDataFile)
	}
	if n.DropDataFile != "" {
		fmt.Fprintf(sb, " :drop_datafile %q", n.DropDataFile)
	}
	if n.InitialSize != "" {
		fmt.Fprintf(sb, " :initial_size %s", n.InitialSize)
	}
	if n.Wait {
		sb.WriteString(" :wait true")
	}
	if n.RenameTo != "" {
		fmt.Fprintf(sb, " :rename_to %s", n.RenameTo)
	}
	if n.AutoextendSize != "" {
		fmt.Fprintf(sb, " :autoextend_size %s", n.AutoextendSize)
	}
	if n.Encryption != "" {
		fmt.Fprintf(sb, " :encryption %q", n.Encryption)
	}
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	if n.EngineAttribute != "" {
		fmt.Fprintf(sb, " :engine_attribute %q", n.EngineAttribute)
	}
	if n.SetActive {
		sb.WriteString(" :set_active true")
	}
	if n.SetInactive {
		sb.WriteString(" :set_inactive true")
	}
	sb.WriteString("}")
}

func writeDropTablespaceStmt(sb *strings.Builder, n *DropTablespaceStmt) {
	sb.WriteString("{DROP_TABLESPACE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Undo {
		sb.WriteString(" :undo true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	sb.WriteString("}")
}

func writeCreateServerStmt(sb *strings.Builder, n *CreateServerStmt) {
	sb.WriteString("{CREATE_SERVER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.WrapperName != "" {
		fmt.Fprintf(sb, " :wrapper %s", n.WrapperName)
	}
	if len(n.Options) > 0 {
		sb.WriteString(" :options (")
		for i, opt := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%q", opt)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeAlterServerStmt(sb *strings.Builder, n *AlterServerStmt) {
	sb.WriteString("{ALTER_SERVER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if len(n.Options) > 0 {
		sb.WriteString(" :options (")
		for i, opt := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%q", opt)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeDropServerStmt(sb *strings.Builder, n *DropServerStmt) {
	sb.WriteString("{DROP_SERVER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	fmt.Fprintf(sb, " :name %s", n.Name)
	sb.WriteString("}")
}

func writeCreateLogfileGroupStmt(sb *strings.Builder, n *CreateLogfileGroupStmt) {
	sb.WriteString("{CREATE_LOGFILE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.UndoFile != "" {
		fmt.Fprintf(sb, " :undo_file %q", n.UndoFile)
	}
	if n.InitialSize != "" {
		fmt.Fprintf(sb, " :initial_size %s", n.InitialSize)
	}
	if n.UndoBufferSize != "" {
		fmt.Fprintf(sb, " :undo_buffer_size %s", n.UndoBufferSize)
	}
	if n.RedoBufferSize != "" {
		fmt.Fprintf(sb, " :redo_buffer_size %s", n.RedoBufferSize)
	}
	if n.NodeGroup != "" {
		fmt.Fprintf(sb, " :nodegroup %s", n.NodeGroup)
	}
	if n.Wait {
		sb.WriteString(" :wait true")
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	sb.WriteString("}")
}

func writeAlterLogfileGroupStmt(sb *strings.Builder, n *AlterLogfileGroupStmt) {
	sb.WriteString("{ALTER_LOGFILE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.UndoFile != "" {
		fmt.Fprintf(sb, " :undo_file %q", n.UndoFile)
	}
	if n.InitialSize != "" {
		fmt.Fprintf(sb, " :initial_size %s", n.InitialSize)
	}
	if n.Wait {
		sb.WriteString(" :wait true")
	}
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	sb.WriteString("}")
}

func writeDropLogfileGroupStmt(sb *strings.Builder, n *DropLogfileGroupStmt) {
	sb.WriteString("{DROP_LOGFILE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Engine != "" {
		fmt.Fprintf(sb, " :engine %s", n.Engine)
	}
	sb.WriteString("}")
}

func writeCreateSpatialRefSysStmt(sb *strings.Builder, n *CreateSpatialRefSysStmt) {
	sb.WriteString("{CREATE_SPATIAL_REF_SYS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	fmt.Fprintf(sb, " :srid %d", n.SRID)
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %q", n.Name)
	}
	if n.Definition != "" {
		fmt.Fprintf(sb, " :definition %q", n.Definition)
	}
	if n.Organization != "" {
		fmt.Fprintf(sb, " :organization %q", n.Organization)
		fmt.Fprintf(sb, " :org_srid %d", n.OrgSRID)
	}
	if n.Description != "" {
		fmt.Fprintf(sb, " :description %q", n.Description)
	}
	sb.WriteString("}")
}

func writeDropSpatialRefSysStmt(sb *strings.Builder, n *DropSpatialRefSysStmt) {
	sb.WriteString("{DROP_SPATIAL_REF_SYS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	fmt.Fprintf(sb, " :srid %d", n.SRID)
	sb.WriteString("}")
}

func writeCreateResourceGroupStmt(sb *strings.Builder, n *CreateResourceGroupStmt) {
	sb.WriteString("{CREATE_RESOURCE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Type != "" {
		fmt.Fprintf(sb, " :type %s", n.Type)
	}
	if len(n.VCPUs) > 0 {
		sb.WriteString(" :vcpus (")
		for i, v := range n.VCPUs {
			if i > 0 {
				sb.WriteString(" ")
			}
			if v.End == -1 {
				fmt.Fprintf(sb, "%d", v.Start)
			} else {
				fmt.Fprintf(sb, "%d-%d", v.Start, v.End)
			}
		}
		sb.WriteString(")")
	}
	if n.ThreadPriority != nil {
		fmt.Fprintf(sb, " :thread_priority %d", *n.ThreadPriority)
	}
	if n.Enable != nil {
		if *n.Enable {
			sb.WriteString(" :enable true")
		} else {
			sb.WriteString(" :enable false")
		}
	}
	sb.WriteString("}")
}

func writeAlterResourceGroupStmt(sb *strings.Builder, n *AlterResourceGroupStmt) {
	sb.WriteString("{ALTER_RESOURCE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if len(n.VCPUs) > 0 {
		sb.WriteString(" :vcpus (")
		for i, v := range n.VCPUs {
			if i > 0 {
				sb.WriteString(" ")
			}
			if v.End == -1 {
				fmt.Fprintf(sb, "%d", v.Start)
			} else {
				fmt.Fprintf(sb, "%d-%d", v.Start, v.End)
			}
		}
		sb.WriteString(")")
	}
	if n.ThreadPriority != nil {
		fmt.Fprintf(sb, " :thread_priority %d", *n.ThreadPriority)
	}
	if n.Enable != nil {
		if *n.Enable {
			sb.WriteString(" :enable true")
		} else {
			sb.WriteString(" :enable false")
		}
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	sb.WriteString("}")
}

func writeDropResourceGroupStmt(sb *strings.Builder, n *DropResourceGroupStmt) {
	sb.WriteString("{DROP_RESOURCE_GROUP")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	fmt.Fprintf(sb, " :name %s", n.Name)
	if n.Force {
		sb.WriteString(" :force true")
	}
	sb.WriteString("}")
}

func writeAlterViewStmt(sb *strings.Builder, n *AlterViewStmt) {
	sb.WriteString("{ALTER_VIEW")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Algorithm != "" {
		fmt.Fprintf(sb, " :algorithm %s", n.Algorithm)
	}
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.SqlSecurity != "" {
		fmt.Fprintf(sb, " :sql_security %s", n.SqlSecurity)
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Columns) > 0 {
		fmt.Fprintf(sb, " :columns %s", strings.Join(n.Columns, ", "))
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	if n.CheckOption != "" {
		fmt.Fprintf(sb, " :check_option %s", n.CheckOption)
	}
	sb.WriteString("}")
}

func writeAlterEventStmt(sb *strings.Builder, n *AlterEventStmt) {
	sb.WriteString("{ALTER_EVENT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Definer != "" {
		fmt.Fprintf(sb, " :definer %s", n.Definer)
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	if n.Schedule != nil {
		sb.WriteString(" :schedule ")
		writeNode(sb, n.Schedule)
	}
	if n.OnCompletion != "" {
		fmt.Fprintf(sb, " :on_completion %s", n.OnCompletion)
	}
	if n.RenameTo != "" {
		fmt.Fprintf(sb, " :rename_to %s", n.RenameTo)
	}
	if n.Enable != "" {
		fmt.Fprintf(sb, " :enable %s", n.Enable)
	}
	if n.Comment != "" {
		fmt.Fprintf(sb, " :comment %q", n.Comment)
	}
	if n.BodyText != "" {
		fmt.Fprintf(sb, " :body %q", n.BodyText)
	}
	if n.Body != nil {
		sb.WriteString(" :body-stmt ")
		writeNode(sb, n.Body)
	}
	sb.WriteString("}")
}

func writeAlterRoutineStmt(sb *strings.Builder, n *AlterRoutineStmt) {
	sb.WriteString("{ALTER_ROUTINE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IsProcedure {
		sb.WriteString(" :is_procedure true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Characteristics) > 0 {
		sb.WriteString(" :characteristics ")
		for i, c := range n.Characteristics {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
	}
	sb.WriteString("}")
}

func writeDropRoutineStmt(sb *strings.Builder, n *DropRoutineStmt) {
	sb.WriteString("{DROP_ROUTINE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IsProcedure {
		sb.WriteString(" :is_procedure true")
	}
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString("}")
}

func writeDropTriggerStmt(sb *strings.Builder, n *DropTriggerStmt) {
	sb.WriteString("{DROP_TRIGGER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString("}")
}

func writeDropEventStmt(sb *strings.Builder, n *DropEventStmt) {
	sb.WriteString("{DROP_EVENT")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name %s", n.Name)
	}
	sb.WriteString("}")
}

func writeChangeReplicationSourceStmt(sb *strings.Builder, n *ChangeReplicationSourceStmt) {
	if n.Legacy {
		sb.WriteString("{CHANGE_MASTER")
	} else {
		sb.WriteString("{CHANGE_REPLICATION_SOURCE")
	}
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Options) > 0 {
		sb.WriteString(" :options (")
		for i, opt := range n.Options {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, opt)
		}
		sb.WriteString(")")
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writeReplicationOption(sb *strings.Builder, n *ReplicationOption) {
	sb.WriteString("{REPL_OPT")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if n.Value != "" {
		fmt.Fprintf(sb, " :val %s", n.Value)
	}
	if len(n.IDs) > 0 {
		sb.WriteString(" :ids (")
		for i, id := range n.IDs {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%d", id)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeChangeReplicationFilterStmt(sb *strings.Builder, n *ChangeReplicationFilterStmt) {
	sb.WriteString("{CHANGE_REPLICATION_FILTER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Filters) > 0 {
		sb.WriteString(" :filters (")
		for i, f := range n.Filters {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, f)
		}
		sb.WriteString(")")
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writeReplicationFilter(sb *strings.Builder, n *ReplicationFilter) {
	sb.WriteString("{REPL_FILTER")
	fmt.Fprintf(sb, " :loc %d :type %s", n.Loc.Start, n.Type)
	if len(n.Values) > 0 {
		sb.WriteString(" :values (")
		for i, v := range n.Values {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%q", v)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeStartReplicaStmt(sb *strings.Builder, n *StartReplicaStmt) {
	sb.WriteString("{START_REPLICA")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IOThread {
		sb.WriteString(" :io_thread true")
	}
	if n.SQLThread {
		sb.WriteString(" :sql_thread true")
	}
	if n.UntilType != "" {
		fmt.Fprintf(sb, " :until_type %s", n.UntilType)
		if n.UntilValue != "" {
			fmt.Fprintf(sb, " :until_value %s", n.UntilValue)
		}
		if n.UntilPos != 0 {
			fmt.Fprintf(sb, " :until_pos %d", n.UntilPos)
		}
	}
	if n.User != "" {
		fmt.Fprintf(sb, " :user %s", n.User)
	}
	if n.Password != "" {
		sb.WriteString(" :password ***")
	}
	if n.DefaultAuth != "" {
		fmt.Fprintf(sb, " :default_auth %s", n.DefaultAuth)
	}
	if n.PluginDir != "" {
		fmt.Fprintf(sb, " :plugin_dir %s", n.PluginDir)
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writeStopReplicaStmt(sb *strings.Builder, n *StopReplicaStmt) {
	sb.WriteString("{STOP_REPLICA")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IOThread {
		sb.WriteString(" :io_thread true")
	}
	if n.SQLThread {
		sb.WriteString(" :sql_thread true")
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writeResetReplicaStmt(sb *strings.Builder, n *ResetReplicaStmt) {
	sb.WriteString("{RESET_REPLICA")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.All {
		sb.WriteString(" :all true")
	}
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	sb.WriteString("}")
}

func writePurgeBinaryLogsStmt(sb *strings.Builder, n *PurgeBinaryLogsStmt) {
	sb.WriteString("{PURGE_BINARY_LOGS")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.To != "" {
		fmt.Fprintf(sb, " :to %s", n.To)
	}
	if n.BeforeExpr != nil {
		sb.WriteString(" :before ")
		writeNode(sb, n.BeforeExpr)
	}
	sb.WriteString("}")
}

func writeResetMasterStmt(sb *strings.Builder, n *ResetMasterStmt) {
	sb.WriteString("{RESET_MASTER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.To != 0 {
		fmt.Fprintf(sb, " :to %d", n.To)
	}
	sb.WriteString("}")
}

func writeStartGroupReplicationStmt(sb *strings.Builder, n *StartGroupReplicationStmt) {
	sb.WriteString("{START_GROUP_REPLICATION")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.User != "" {
		fmt.Fprintf(sb, " :user %s", n.User)
	}
	if n.Password != "" {
		sb.WriteString(" :password ***")
	}
	if n.DefaultAuth != "" {
		fmt.Fprintf(sb, " :default_auth %s", n.DefaultAuth)
	}
	sb.WriteString("}")
}

func writeStopGroupReplicationStmt(sb *strings.Builder, n *StopGroupReplicationStmt) {
	fmt.Fprintf(sb, "{STOP_GROUP_REPLICATION :loc %d}", n.Loc.Start)
}

func writeAlterInstanceStmt(sb *strings.Builder, n *AlterInstanceStmt) {
	sb.WriteString("{ALTER_INSTANCE")
	fmt.Fprintf(sb, " :loc %d :action %s", n.Loc.Start, n.Action)
	if n.Channel != "" {
		fmt.Fprintf(sb, " :channel %s", n.Channel)
	}
	if n.NoRollbackOnError {
		sb.WriteString(" :no_rollback_on_error true")
	}
	sb.WriteString("}")
}

func writeImportTableStmt(sb *strings.Builder, n *ImportTableStmt) {
	sb.WriteString("{IMPORT_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Files) > 0 {
		sb.WriteString(" :files (")
		for i, f := range n.Files {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%q", f)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeCacheIndexStmt(sb *strings.Builder, n *CacheIndexStmt) {
	sb.WriteString("{CACHE_INDEX")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables (")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
		sb.WriteString(")")
	}
	if len(n.Partitions) > 0 {
		sb.WriteString(" :partitions (")
		for i, p := range n.Partitions {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(p)
		}
		sb.WriteString(")")
	}
	if len(n.Indexes) > 0 {
		sb.WriteString(" :indexes (")
		for i, idx := range n.Indexes {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(idx)
		}
		sb.WriteString(")")
	}
	if n.CacheName != "" {
		fmt.Fprintf(sb, " :cache %s", n.CacheName)
	}
	sb.WriteString("}")
}

func writeLoadIndexIntoCacheStmt(sb *strings.Builder, n *LoadIndexIntoCacheStmt) {
	sb.WriteString("{LOAD_INDEX_INTO_CACHE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables (")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeLoadIndexTable(sb *strings.Builder, n *LoadIndexTable) {
	sb.WriteString("{LOAD_INDEX_TABLE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if len(n.Partitions) > 0 {
		fmt.Fprintf(sb, " :partitions (%s)", strings.Join(n.Partitions, ", "))
	}
	if len(n.Indexes) > 0 {
		fmt.Fprintf(sb, " :indexes (%s)", strings.Join(n.Indexes, ", "))
	}
	if n.IgnoreLeaves {
		sb.WriteString(" :ignore_leaves true")
	}
	sb.WriteString("}")
}

func writeResetPersistStmt(sb *strings.Builder, n *ResetPersistStmt) {
	sb.WriteString("{RESET_PERSIST")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.Variable != "" {
		fmt.Fprintf(sb, " :variable %s", n.Variable)
	}
	sb.WriteString("}")
}

func writeRenameUserStmt(sb *strings.Builder, n *RenameUserStmt) {
	sb.WriteString("{RENAME_USER")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Pairs) > 0 {
		sb.WriteString(" :pairs (")
		for i, p := range n.Pairs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeRenameUserPair(sb *strings.Builder, n *RenameUserPair) {
	sb.WriteString("{RENAME_USER_PAIR")
	fmt.Fprintf(sb, " :loc %d :old %s", n.Loc.Start, n.OldUser)
	if n.OldHost != "" {
		fmt.Fprintf(sb, "@%s", n.OldHost)
	}
	fmt.Fprintf(sb, " :new %s", n.NewUser)
	if n.NewHost != "" {
		fmt.Fprintf(sb, "@%s", n.NewHost)
	}
	sb.WriteString("}")
}

func writeSetResourceGroupStmt(sb *strings.Builder, n *SetResourceGroupStmt) {
	sb.WriteString("{SET_RESOURCE_GROUP")
	fmt.Fprintf(sb, " :loc %d :name %s", n.Loc.Start, n.Name)
	if len(n.ThreadIDs) > 0 {
		sb.WriteString(" :threads (")
		for i, id := range n.ThreadIDs {
			if i > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(sb, "%d", id)
		}
		sb.WriteString(")")
	}
	sb.WriteString("}")
}

func writeVCPUSpec(sb *strings.Builder, n *VCPUSpec) {
	if n.End < 0 {
		fmt.Fprintf(sb, "{VCPU %d}", n.Start)
	} else {
		fmt.Fprintf(sb, "{VCPU %d-%d}", n.Start, n.End)
	}
}
