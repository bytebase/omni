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
		fmt.Fprintf(sb, "{DEFAULT :loc %d}", n.Loc.Start)
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
		fmt.Fprintf(sb, "{TABLE_OPT :loc %d :name %s :val %s}", n.Loc.Start, n.Name, n.Value)
	case *DatabaseOption:
		fmt.Fprintf(sb, "{DB_OPT :loc %d :name %s :val %s}", n.Loc.Start, n.Name, n.Value)
	case *PartitionClause:
		writePartitionClause(sb, n)
	case *PartitionDef:
		writePartitionDef(sb, n)
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
	}
	if n.Having != nil {
		sb.WriteString(" :having ")
		writeNode(sb, n.Having)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString(" :order_by ")
		writeOrderByList(sb, n.OrderBy)
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
	if len(n.SetList) > 0 {
		sb.WriteString(" :set ")
		for i, a := range n.SetList {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
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
	if n.Format != "" {
		fmt.Fprintf(sb, " :format %s", n.Format)
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
	sb.WriteString("{SAVEPOINT")
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
	sb.WriteString("}")
}

func writeRevokeStmt(sb *strings.Builder, n *RevokeStmt) {
	sb.WriteString("{REVOKE")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.AllPriv {
		sb.WriteString(" :all_priv true")
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
	if len(n.Users) > 0 {
		sb.WriteString(" :users ")
		for i, u := range n.Users {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, u)
		}
	}
	sb.WriteString("}")
}

func writeCreateFunctionStmt(sb *strings.Builder, n *CreateFunctionStmt) {
	sb.WriteString("{CREATE_FUNCTION")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.IsProcedure {
		sb.WriteString(" :is_procedure true")
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
	if n.Body != "" {
		fmt.Fprintf(sb, " :body %q", n.Body)
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
	if n.Body != "" {
		fmt.Fprintf(sb, " :body %q", n.Body)
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
	if n.Body != "" {
		fmt.Fprintf(sb, " :body %q", n.Body)
	}
	sb.WriteString("}")
}

func writeLoadDataStmt(sb *strings.Builder, n *LoadDataStmt) {
	sb.WriteString("{LOAD_DATA")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
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
	if n.LinesTerminatedBy != "" {
		fmt.Fprintf(sb, " :lines_terminated %q", n.LinesTerminatedBy)
	}
	if n.FieldsTerminatedBy != "" {
		fmt.Fprintf(sb, " :fields_terminated %q", n.FieldsTerminatedBy)
	}
	if n.FieldsEnclosedBy != "" {
		fmt.Fprintf(sb, " :fields_enclosed %q", n.FieldsEnclosedBy)
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

func writeOptimizeTableStmt(sb *strings.Builder, n *OptimizeTableStmt) {
	sb.WriteString("{OPTIMIZE_TABLE")
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

func writeFlushStmt(sb *strings.Builder, n *FlushStmt) {
	sb.WriteString("{FLUSH")
	fmt.Fprintf(sb, " :loc %d", n.Loc.Start)
	if len(n.Options) > 0 {
		fmt.Fprintf(sb, " :options %s", strings.Join(n.Options, ", "))
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
	fmt.Fprintf(sb, " :op %s :left ", binaryOpStr(n.Op))
	writeNode(sb, n.Left)
	sb.WriteString(" :right ")
	writeNode(sb, n.Right)
	sb.WriteString("}")
}

func writeUnaryExpr(sb *strings.Builder, n *UnaryExpr) {
	sb.WriteString("{UNARY")
	fmt.Fprintf(sb, " :op %s :operand ", unaryOpStr(n.Op))
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
	if n.OnDelete != RefActNone {
		fmt.Fprintf(sb, " :on_delete %d", n.OnDelete)
	}
	if n.OnUpdate != RefActNone {
		fmt.Fprintf(sb, " :on_update %d", n.OnUpdate)
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
	if n.IndexType != "" {
		fmt.Fprintf(sb, " :index_type %s", n.IndexType)
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
	sb.WriteString("}")
}

func writePartitionClause(sb *strings.Builder, n *PartitionClause) {
	sb.WriteString("{PARTITION_CLAUSE")
	fmt.Fprintf(sb, " :loc %d :type %d", n.Loc.Start, n.Type)
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if len(n.Columns) > 0 {
		fmt.Fprintf(sb, " :columns %s", strings.Join(n.Columns, ", "))
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
