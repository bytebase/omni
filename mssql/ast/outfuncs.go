package ast

import (
	"fmt"
	"strings"
)

// NodeToString converts a Node to its string representation for deterministic AST comparison.
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
	case *Integer:
		sb.WriteString(fmt.Sprintf("%d", n.Ival))
	case *Float:
		sb.WriteString(n.Fval)
	case *Boolean:
		if n.Boolval {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	case *String:
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Str))
		sb.WriteString("\"")
	case *SelectStmt:
		writeSelectStmt(sb, n)
	case *InsertStmt:
		writeInsertStmt(sb, n)
	case *UpdateStmt:
		writeUpdateStmt(sb, n)
	case *DeleteStmt:
		writeDeleteStmt(sb, n)
	case *MergeStmt:
		writeMergeStmt(sb, n)
	case *CreateTableStmt:
		writeCreateTableStmt(sb, n)
	case *AlterTableStmt:
		writeAlterTableStmt(sb, n)
	case *DropStmt:
		writeDropStmt(sb, n)
	case *CreateIndexStmt:
		writeCreateIndexStmt(sb, n)
	case *CreateViewStmt:
		writeCreateViewStmt(sb, n)
	case *CreateTriggerStmt:
		writeCreateTriggerStmt(sb, n)
	case *CreateFunctionStmt:
		writeCreateFunctionStmt(sb, n)
	case *CreateProcedureStmt:
		writeCreateProcedureStmt(sb, n)
	case *CreateDatabaseStmt:
		writeCreateDatabaseStmt(sb, n)
	case *AlterDatabaseStmt:
		writeAlterDatabaseStmt(sb, n)
	case *AlterIndexStmt:
		writeAlterIndexStmt(sb, n)
	case *TruncateStmt:
		writeTruncateStmt(sb, n)
	case *DeclareStmt:
		writeDeclareStmt(sb, n)
	case *SetStmt:
		writeSetStmt(sb, n)
	case *IfStmt:
		writeIfStmt(sb, n)
	case *WhileStmt:
		writeWhileStmt(sb, n)
	case *BeginEndStmt:
		writeBeginEndStmt(sb, n)
	case *TryCatchStmt:
		writeTryCatchStmt(sb, n)
	case *ExecStmt:
		writeExecStmt(sb, n)
	case *PrintStmt:
		writePrintStmt(sb, n)
	case *RaiseErrorStmt:
		writeRaiseErrorStmt(sb, n)
	case *ThrowStmt:
		writeThrowStmt(sb, n)
	case *UseStmt:
		writeUseStmt(sb, n)
	case *GoStmt:
		writeGoStmt(sb, n)
	case *ReturnStmt:
		writeReturnStmt(sb, n)
	case *BreakStmt:
		sb.WriteString(fmt.Sprintf("{BREAK :loc %d %d}", n.Loc.Start, n.Loc.End))
	case *ContinueStmt:
		sb.WriteString(fmt.Sprintf("{CONTINUE :loc %d %d}", n.Loc.Start, n.Loc.End))
	case *GotoStmt:
		sb.WriteString(fmt.Sprintf("{GOTO :label \"%s\" :loc %d %d}", escapeString(n.Label), n.Loc.Start, n.Loc.End))
	case *LabelStmt:
		sb.WriteString(fmt.Sprintf("{LABEL :label \"%s\" :loc %d %d}", escapeString(n.Label), n.Loc.Start, n.Loc.End))
	case *WaitForStmt:
		writeWaitForStmt(sb, n)
	case *BeginTransStmt:
		sb.WriteString(fmt.Sprintf("{BEGINTRANS :name \"%s\" :loc %d %d}", escapeString(n.Name), n.Loc.Start, n.Loc.End))
	case *CommitTransStmt:
		sb.WriteString(fmt.Sprintf("{COMMITTRANS :name \"%s\" :loc %d %d}", escapeString(n.Name), n.Loc.Start, n.Loc.End))
	case *RollbackTransStmt:
		writeRollbackTransStmt(sb, n)
	case *SaveTransStmt:
		sb.WriteString(fmt.Sprintf("{SAVETRANS :name \"%s\" :loc %d %d}", escapeString(n.Name), n.Loc.Start, n.Loc.End))
	case *SecurityStmt:
		writeSecurityStmt(sb, n)
	case *CreateSchemaStmt:
		writeCreateSchemaStmt(sb, n)
	case *AlterSchemaStmt:
		writeAlterSchemaStmt(sb, n)
	case *CreateTypeStmt:
		writeCreateTypeStmt(sb, n)
	case *CreateSequenceStmt:
		writeCreateSequenceStmt(sb, n)
	case *AlterSequenceStmt:
		writeAlterSequenceStmt(sb, n)
	case *CreateSynonymStmt:
		writeCreateSynonymStmt(sb, n)
	case *GrantStmt:
		writeGrantStmt(sb, n)
	case *BinaryExpr:
		writeBinaryExpr(sb, n)
	case *UnaryExpr:
		writeUnaryExpr(sb, n)
	case *FuncCallExpr:
		writeFuncCallExpr(sb, n)
	case *CaseExpr:
		writeCaseExpr(sb, n)
	case *CaseWhen:
		writeCaseWhen(sb, n)
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
	case *ConvertExpr:
		writeConvertExpr(sb, n)
	case *TryCastExpr:
		writeTryCastExpr(sb, n)
	case *TryConvertExpr:
		writeTryConvertExpr(sb, n)
	case *CoalesceExpr:
		writeCoalesceExpr(sb, n)
	case *NullifExpr:
		writeNullifExpr(sb, n)
	case *IifExpr:
		writeIifExpr(sb, n)
	case *ColumnRef:
		writeColumnRef(sb, n)
	case *VariableRef:
		sb.WriteString(fmt.Sprintf("{VARREF :name \"%s\" :loc %d %d}", escapeString(n.Name), n.Loc.Start, n.Loc.End))
	case *StarExpr:
		writeStarExpr(sb, n)
	case *Literal:
		writeLiteral(sb, n)
	case *SubqueryExpr:
		writeSubqueryExpr(sb, n)
	case *ParenExpr:
		writeParenExpr(sb, n)
	case *TableRef:
		writeTableRef(sb, n)
	case *DataType:
		writeDataType(sb, n)
	case *JoinClause:
		writeJoinClause(sb, n)
	case *AliasedTableRef:
		writeAliasedTableRef(sb, n)
	case *OverClause:
		writeOverClause(sb, n)
	case *OrderByItem:
		writeOrderByItem(sb, n)
	case *ResTarget:
		writeResTarget(sb, n)
	case *SetExpr:
		writeSetExpr(sb, n)
	case *ValuesClause:
		writeValuesClause(sb, n)
	case *TopClause:
		writeTopClause(sb, n)
	case *WithClause:
		writeWithClause(sb, n)
	case *CommonTableExpr:
		writeCommonTableExpr(sb, n)
	case *OutputClause:
		writeOutputClause(sb, n)
	case *ColumnDef:
		writeColumnDef(sb, n)
	case *ConstraintDef:
		writeConstraintDef(sb, n)
	case *MergeWhenClause:
		writeMergeWhenClause(sb, n)
	case *IndexColumn:
		writeIndexColumn(sb, n)
	case *ParamDef:
		writeParamDef(sb, n)
	case *VariableDecl:
		writeVariableDecl(sb, n)
	case *SelectAssign:
		writeSelectAssign(sb, n)
	case *MethodCallExpr:
		writeMethodCallExpr(sb, n)
	case *FetchClause:
		writeFetchClause(sb, n)
	case *ForClause:
		writeForClause(sb, n)
	case *NullableSpec:
		writeNullableSpec(sb, n)
	case *IdentitySpec:
		writeIdentitySpec(sb, n)
	case *ComputedColumnDef:
		writeComputedColumnDef(sb, n)
	case *AlterTableAction:
		writeAlterTableAction(sb, n)
	case *ReturnsTableDef:
		writeReturnsTableDef(sb, n)
	case *ExecArg:
		writeExecArg(sb, n)
	case *MergeUpdateAction:
		writeMergeUpdateAction(sb, n)
	case *MergeDeleteAction:
		writeMergeDeleteAction(sb, n)
	case *MergeInsertAction:
		writeMergeInsertAction(sb, n)
	case *WindowFrame:
		writeWindowFrame(sb, n)
	case *WindowBound:
		writeWindowBound(sb, n)
	case *DeclareCursorStmt:
		writeDeclareCursorStmt(sb, n)
	case *OpenCursorStmt:
		writeOpenCursorStmt(sb, n)
	case *FetchCursorStmt:
		writeFetchCursorStmt(sb, n)
	case *CloseCursorStmt:
		writeCloseCursorStmt(sb, n)
	case *DeallocateCursorStmt:
		writeDeallocateCursorStmt(sb, n)
	case *BulkInsertStmt:
		writeBulkInsertStmt(sb, n)
	case *DbccStmt:
		writeDbccStmt(sb, n)
	case *BackupStmt:
		writeBackupStmt(sb, n)
	case *RestoreStmt:
		writeRestoreStmt(sb, n)
	case *SecurityKeyStmt:
		writeSecurityKeyStmt(sb, n)
	case *BeginDistributedTransStmt:
		fmt.Fprintf(sb, "{BEGINDISTRIBUTEDTRANS :name \"%s\" :loc %d %d}", escapeString(n.Name), n.Loc.Start, n.Loc.End)
	case *CreateStatisticsStmt:
		writeCreateStatisticsStmt(sb, n)
	case *UpdateStatisticsStmt:
		writeUpdateStatisticsStmt(sb, n)
	case *DropStatisticsStmt:
		writeDropStatisticsStmt(sb, n)
	case *SetOptionStmt:
		writeSetOptionStmt(sb, n)
	case *CreatePartitionFunctionStmt:
		writeCreatePartitionFunctionStmt(sb, n)
	case *AlterPartitionFunctionStmt:
		writeAlterPartitionFunctionStmt(sb, n)
	case *CreatePartitionSchemeStmt:
		writeCreatePartitionSchemeStmt(sb, n)
	case *AlterPartitionSchemeStmt:
		fmt.Fprintf(sb, "{ALTERPARTITIONSCHEME :name \"%s\" :filegroup \"%s\" :loc %d %d}", escapeString(n.Name), escapeString(n.FileGroup), n.Loc.Start, n.Loc.End)
	case *CreateFulltextIndexStmt:
		writeCreateFulltextIndexStmt(sb, n)
	case *AlterFulltextIndexStmt:
		writeAlterFulltextIndexStmt(sb, n)
	case *CreateFulltextCatalogStmt:
		writeCreateFulltextCatalogStmt(sb, n)
	case *AlterFulltextCatalogStmt:
		fmt.Fprintf(sb, "{ALTERFULLTEXTCATALOG :name \"%s\" :action \"%s\" :loc %d %d}", escapeString(n.Name), escapeString(n.Action), n.Loc.Start, n.Loc.End)
	case *CreateXmlSchemaCollectionStmt:
		writeCreateXmlSchemaCollectionStmt(sb, n)
	case *AlterXmlSchemaCollectionStmt:
		writeAlterXmlSchemaCollectionStmt(sb, n)
	case *CreateAssemblyStmt:
		writeCreateAssemblyStmt(sb, n)
	case *AlterAssemblyStmt:
		writeAlterAssemblyStmt(sb, n)
	case *ServiceBrokerStmt:
		writeServiceBrokerStmt(sb, n)
	case *CheckpointStmt:
		writeCheckpointStmt(sb, n)
	case *ReconfigureStmt:
		fmt.Fprintf(sb, "{RECONFIGURE :withOverride %v :loc %d %d}", n.WithOverride, n.Loc.Start, n.Loc.End)
	case *ShutdownStmt:
		fmt.Fprintf(sb, "{SHUTDOWN :withNoWait %v :loc %d %d}", n.WithNoWait, n.Loc.Start, n.Loc.End)
	case *KillStmt:
		writeKillStmt(sb, n)
	case *ReadtextStmt:
		writeReadtextStmt(sb, n)
	case *WritetextStmt:
		writeWritetextStmt(sb, n)
	case *UpdatetextStmt:
		writeUpdatetextStmt(sb, n)
	case *PivotExpr:
		writePivotExpr(sb, n)
	case *UnpivotExpr:
		writeUnpivotExpr(sb, n)
	case *TableSampleClause:
		writeTableSampleClause(sb, n)
	case *GroupingSetsExpr:
		writeGroupingSetsExpr(sb, n)
	case *RollupExpr:
		writeRollupExpr(sb, n)
	case *CubeExpr:
		writeCubeExpr(sb, n)
	default:
		sb.WriteString("{UNKNOWN}")
	}
}

func writeList(sb *strings.Builder, n *List) {
	sb.WriteString("(")
	for i, item := range n.Items {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, item)
	}
	sb.WriteString(")")
}

func writeSelectStmt(sb *strings.Builder, n *SelectStmt) {
	sb.WriteString("{SELECTSTMT")
	if n.WithClause != nil {
		sb.WriteString(" :with ")
		writeNode(sb, n.WithClause)
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.All {
		sb.WriteString(" :all true")
	}
	if n.Top != nil {
		sb.WriteString(" :top ")
		writeNode(sb, n.Top)
	}
	if n.TargetList != nil {
		sb.WriteString(" :targetList ")
		writeNode(sb, n.TargetList)
	}
	if n.IntoTable != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.IntoTable)
	}
	if n.FromClause != nil {
		sb.WriteString(" :from ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.WhereClause)
	}
	if n.GroupByClause != nil {
		sb.WriteString(" :groupBy ")
		writeNode(sb, n.GroupByClause)
	}
	if n.HavingClause != nil {
		sb.WriteString(" :having ")
		writeNode(sb, n.HavingClause)
	}
	if n.OrderByClause != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderByClause)
	}
	if n.OffsetClause != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.OffsetClause)
	}
	if n.FetchClause != nil {
		sb.WriteString(" :fetch ")
		writeNode(sb, n.FetchClause)
	}
	if n.ForClause != nil {
		sb.WriteString(" :for ")
		writeNode(sb, n.ForClause)
	}
	if n.OptionClause != nil {
		sb.WriteString(" :option ")
		writeNode(sb, n.OptionClause)
	}
	if n.Op != SetOpNone {
		sb.WriteString(fmt.Sprintf(" :op %d", n.Op))
	}
	if n.Larg != nil {
		sb.WriteString(" :larg ")
		writeNode(sb, n.Larg)
	}
	if n.Rarg != nil {
		sb.WriteString(" :rarg ")
		writeNode(sb, n.Rarg)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeInsertStmt(sb *strings.Builder, n *InsertStmt) {
	sb.WriteString("{INSERTSTMT")
	if n.WithClause != nil {
		sb.WriteString(" :with ")
		writeNode(sb, n.WithClause)
	}
	if n.Top != nil {
		sb.WriteString(" :top ")
		writeNode(sb, n.Top)
	}
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.Cols != nil {
		sb.WriteString(" :cols ")
		writeNode(sb, n.Cols)
	}
	if n.Source != nil {
		sb.WriteString(" :source ")
		writeNode(sb, n.Source)
	}
	if n.OutputClause != nil {
		sb.WriteString(" :output ")
		writeNode(sb, n.OutputClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUpdateStmt(sb *strings.Builder, n *UpdateStmt) {
	sb.WriteString("{UPDATESTMT")
	if n.WithClause != nil {
		sb.WriteString(" :with ")
		writeNode(sb, n.WithClause)
	}
	if n.Top != nil {
		sb.WriteString(" :top ")
		writeNode(sb, n.Top)
	}
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.SetClause != nil {
		sb.WriteString(" :set ")
		writeNode(sb, n.SetClause)
	}
	if n.FromClause != nil {
		sb.WriteString(" :from ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.WhereClause)
	}
	if n.OutputClause != nil {
		sb.WriteString(" :output ")
		writeNode(sb, n.OutputClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDeleteStmt(sb *strings.Builder, n *DeleteStmt) {
	sb.WriteString("{DELETESTMT")
	if n.WithClause != nil {
		sb.WriteString(" :with ")
		writeNode(sb, n.WithClause)
	}
	if n.Top != nil {
		sb.WriteString(" :top ")
		writeNode(sb, n.Top)
	}
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.FromClause != nil {
		sb.WriteString(" :from ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.WhereClause)
	}
	if n.OutputClause != nil {
		sb.WriteString(" :output ")
		writeNode(sb, n.OutputClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeStmt(sb *strings.Builder, n *MergeStmt) {
	sb.WriteString("{MERGESTMT")
	if n.WithClause != nil {
		sb.WriteString(" :with ")
		writeNode(sb, n.WithClause)
	}
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	if n.Source != nil {
		sb.WriteString(" :source ")
		writeNode(sb, n.Source)
	}
	if n.SourceAlias != "" {
		sb.WriteString(fmt.Sprintf(" :sourceAlias \"%s\"", escapeString(n.SourceAlias)))
	}
	if n.OnCondition != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.OnCondition)
	}
	if n.WhenClauses != nil {
		sb.WriteString(" :whenClauses ")
		writeNode(sb, n.WhenClauses)
	}
	if n.OutputClause != nil {
		sb.WriteString(" :output ")
		writeNode(sb, n.OutputClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTableStmt(sb *strings.Builder, n *CreateTableStmt) {
	sb.WriteString("{CREATETABLE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	sb.WriteString(fmt.Sprintf(" :ifNotExists %t", n.IfNotExists))
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTableStmt(sb *strings.Builder, n *AlterTableStmt) {
	sb.WriteString("{ALTERTABLE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Actions != nil {
		sb.WriteString(" :actions ")
		writeNode(sb, n.Actions)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDropStmt(sb *strings.Builder, n *DropStmt) {
	sb.WriteString("{DROPSTMT")
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Names != nil {
		sb.WriteString(" :names ")
		writeNode(sb, n.Names)
	}
	sb.WriteString(fmt.Sprintf(" :ifExists %t", n.IfExists))
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateIndexStmt(sb *strings.Builder, n *CreateIndexStmt) {
	sb.WriteString("{CREATEINDEX")
	sb.WriteString(fmt.Sprintf(" :unique %t", n.Unique))
	if n.Clustered != nil {
		sb.WriteString(fmt.Sprintf(" :clustered %t", *n.Clustered))
	}
	if n.Columnstore {
		sb.WriteString(" :columnstore true")
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.IncludeCols != nil {
		sb.WriteString(" :include ")
		writeNode(sb, n.IncludeCols)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.WhereClause)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	if n.OnFileGroup != "" {
		sb.WriteString(fmt.Sprintf(" :onFileGroup \"%s\"", escapeString(n.OnFileGroup)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateViewStmt(sb *strings.Builder, n *CreateViewStmt) {
	sb.WriteString("{CREATEVIEW")
	sb.WriteString(fmt.Sprintf(" :orAlter %t", n.OrAlter))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	if n.WithCheck {
		sb.WriteString(" :withCheck true")
	}
	if n.SchemaBinding {
		sb.WriteString(" :schemaBinding true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTriggerStmt(sb *strings.Builder, n *CreateTriggerStmt) {
	sb.WriteString("{CREATETRIGGER")
	sb.WriteString(fmt.Sprintf(" :orAlter %t", n.OrAlter))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.OnDatabase {
		sb.WriteString(" :onDatabase true")
	}
	if n.OnAllServer {
		sb.WriteString(" :onAllServer true")
	}
	if n.TriggerType != "" {
		sb.WriteString(fmt.Sprintf(" :triggerType \"%s\"", n.TriggerType))
	}
	if n.Events != nil {
		sb.WriteString(" :events ")
		writeNode(sb, n.Events)
	}
	if n.WithAppend {
		sb.WriteString(" :withAppend true")
	}
	if n.NotForReplication {
		sb.WriteString(" :notForReplication true")
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateFunctionStmt(sb *strings.Builder, n *CreateFunctionStmt) {
	sb.WriteString("{CREATEFUNCTION")
	sb.WriteString(fmt.Sprintf(" :orAlter %t", n.OrAlter))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Params != nil {
		sb.WriteString(" :params ")
		writeNode(sb, n.Params)
	}
	if n.ReturnType != nil {
		sb.WriteString(" :returnType ")
		writeNode(sb, n.ReturnType)
	}
	if n.ReturnsTable != nil {
		sb.WriteString(" :returnsTable ")
		writeNode(sb, n.ReturnsTable)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateProcedureStmt(sb *strings.Builder, n *CreateProcedureStmt) {
	sb.WriteString("{CREATEPROCEDURE")
	sb.WriteString(fmt.Sprintf(" :orAlter %t", n.OrAlter))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Params != nil {
		sb.WriteString(" :params ")
		writeNode(sb, n.Params)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateDatabaseStmt(sb *strings.Builder, n *CreateDatabaseStmt) {
	sb.WriteString("{CREATEDATABASE")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTruncateStmt(sb *strings.Builder, n *TruncateStmt) {
	sb.WriteString("{TRUNCATE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDeclareStmt(sb *strings.Builder, n *DeclareStmt) {
	sb.WriteString("{DECLARE")
	if n.Variables != nil {
		sb.WriteString(" :variables ")
		writeNode(sb, n.Variables)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetStmt(sb *strings.Builder, n *SetStmt) {
	sb.WriteString("{SET")
	sb.WriteString(fmt.Sprintf(" :variable \"%s\"", escapeString(n.Variable)))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIfStmt(sb *strings.Builder, n *IfStmt) {
	sb.WriteString("{IF")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Then != nil {
		sb.WriteString(" :then ")
		writeNode(sb, n.Then)
	}
	if n.Else != nil {
		sb.WriteString(" :else ")
		writeNode(sb, n.Else)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWhileStmt(sb *strings.Builder, n *WhileStmt) {
	sb.WriteString("{WHILE")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeBeginEndStmt(sb *strings.Builder, n *BeginEndStmt) {
	sb.WriteString("{BEGINEND")
	if n.Stmts != nil {
		sb.WriteString(" :stmts ")
		writeNode(sb, n.Stmts)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTryCatchStmt(sb *strings.Builder, n *TryCatchStmt) {
	sb.WriteString("{TRYCATCH")
	if n.TryBlock != nil {
		sb.WriteString(" :try ")
		writeNode(sb, n.TryBlock)
	}
	if n.CatchBlock != nil {
		sb.WriteString(" :catch ")
		writeNode(sb, n.CatchBlock)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExecStmt(sb *strings.Builder, n *ExecStmt) {
	sb.WriteString("{EXEC")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.ReturnVar != "" {
		sb.WriteString(fmt.Sprintf(" :returnVar \"%s\"", escapeString(n.ReturnVar)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePrintStmt(sb *strings.Builder, n *PrintStmt) {
	sb.WriteString("{PRINT")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRaiseErrorStmt(sb *strings.Builder, n *RaiseErrorStmt) {
	sb.WriteString("{RAISERROR")
	if n.Message != nil {
		sb.WriteString(" :message ")
		writeNode(sb, n.Message)
	}
	if n.Severity != nil {
		sb.WriteString(" :severity ")
		writeNode(sb, n.Severity)
	}
	if n.State != nil {
		sb.WriteString(" :state ")
		writeNode(sb, n.State)
	}
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeThrowStmt(sb *strings.Builder, n *ThrowStmt) {
	sb.WriteString("{THROW")
	if n.ErrorNumber != nil {
		sb.WriteString(" :errorNumber ")
		writeNode(sb, n.ErrorNumber)
	}
	if n.Message != nil {
		sb.WriteString(" :message ")
		writeNode(sb, n.Message)
	}
	if n.State != nil {
		sb.WriteString(" :state ")
		writeNode(sb, n.State)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUseStmt(sb *strings.Builder, n *UseStmt) {
	sb.WriteString(fmt.Sprintf("{USE :database \"%s\" :loc %d %d}", escapeString(n.Database), n.Loc.Start, n.Loc.End))
}

func writeGoStmt(sb *strings.Builder, n *GoStmt) {
	sb.WriteString(fmt.Sprintf("{GO :count %d :loc %d %d}", n.Count, n.Loc.Start, n.Loc.End))
}

func writeReturnStmt(sb *strings.Builder, n *ReturnStmt) {
	sb.WriteString("{RETURN")
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWaitForStmt(sb *strings.Builder, n *WaitForStmt) {
	sb.WriteString("{WAITFOR")
	sb.WriteString(fmt.Sprintf(" :waitType \"%s\"", escapeString(n.WaitType)))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRollbackTransStmt(sb *strings.Builder, n *RollbackTransStmt) {
	sb.WriteString("{ROLLBACKTRANS")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Savepoint != "" {
		sb.WriteString(fmt.Sprintf(" :savepoint \"%s\"", escapeString(n.Savepoint)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSecurityStmt(sb *strings.Builder, n *SecurityStmt) {
	sb.WriteString("{SECURITY")
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	fmt.Fprintf(sb, " :objectType \"%s\"", escapeString(n.ObjectType))
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Options != nil && len(n.Options.Items) > 0 {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeCreateSchemaStmt(sb *strings.Builder, n *CreateSchemaStmt) {
	sb.WriteString("{CREATESCHEMA")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Authorization != "" {
		fmt.Fprintf(sb, " :authorization \"%s\"", escapeString(n.Authorization))
	}
	if n.Elements != nil && len(n.Elements.Items) > 0 {
		sb.WriteString(" :elements ")
		writeNode(sb, n.Elements)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeAlterSchemaStmt(sb *strings.Builder, n *AlterSchemaStmt) {
	sb.WriteString("{ALTERSCHEMA")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.TransferType != "" {
		fmt.Fprintf(sb, " :transferType \"%s\"", escapeString(n.TransferType))
	}
	if n.TransferEntity != "" {
		fmt.Fprintf(sb, " :transferEntity \"%s\"", escapeString(n.TransferEntity))
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeCreateTypeStmt(sb *strings.Builder, n *CreateTypeStmt) {
	sb.WriteString("{CREATETYPE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.BaseType != nil {
		sb.WriteString(" :baseType ")
		writeNode(sb, n.BaseType)
	}
	if n.Nullable != nil {
		fmt.Fprintf(sb, " :nullable %v", *n.Nullable)
	}
	if n.ExternalName != "" {
		fmt.Fprintf(sb, " :externalName \"%s\"", escapeString(n.ExternalName))
	}
	if n.TableDef != nil && len(n.TableDef.Items) > 0 {
		sb.WriteString(" :tableDef ")
		writeNode(sb, n.TableDef)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeCreateSequenceStmt(sb *strings.Builder, n *CreateSequenceStmt) {
	sb.WriteString("{CREATESEQUENCE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Start != nil {
		sb.WriteString(" :start ")
		writeNode(sb, n.Start)
	}
	if n.Increment != nil {
		sb.WriteString(" :increment ")
		writeNode(sb, n.Increment)
	}
	if n.MinValue != nil {
		sb.WriteString(" :minValue ")
		writeNode(sb, n.MinValue)
	}
	if n.NoMinVal {
		sb.WriteString(" :noMinVal true")
	}
	if n.MaxValue != nil {
		sb.WriteString(" :maxValue ")
		writeNode(sb, n.MaxValue)
	}
	if n.NoMaxVal {
		sb.WriteString(" :noMaxVal true")
	}
	if n.Cycle != nil {
		fmt.Fprintf(sb, " :cycle %t", *n.Cycle)
	}
	if n.Cache != nil {
		sb.WriteString(" :cache ")
		writeNode(sb, n.Cache)
	}
	if n.NoCache {
		sb.WriteString(" :noCache true")
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeAlterSequenceStmt(sb *strings.Builder, n *AlterSequenceStmt) {
	sb.WriteString("{ALTERSEQUENCE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Restart {
		sb.WriteString(" :restart true")
	}
	if n.RestartWith != nil {
		sb.WriteString(" :restartWith ")
		writeNode(sb, n.RestartWith)
	}
	if n.Increment != nil {
		sb.WriteString(" :increment ")
		writeNode(sb, n.Increment)
	}
	if n.MinValue != nil {
		sb.WriteString(" :minValue ")
		writeNode(sb, n.MinValue)
	}
	if n.NoMinVal {
		sb.WriteString(" :noMinVal true")
	}
	if n.MaxValue != nil {
		sb.WriteString(" :maxValue ")
		writeNode(sb, n.MaxValue)
	}
	if n.NoMaxVal {
		sb.WriteString(" :noMaxVal true")
	}
	if n.Cycle != nil {
		fmt.Fprintf(sb, " :cycle %t", *n.Cycle)
	}
	if n.Cache != nil {
		sb.WriteString(" :cache ")
		writeNode(sb, n.Cache)
	}
	if n.NoCache {
		sb.WriteString(" :noCache true")
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeCreateSynonymStmt(sb *strings.Builder, n *CreateSynonymStmt) {
	sb.WriteString("{CREATESYNONYM")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeGrantStmt(sb *strings.Builder, n *GrantStmt) {
	sb.WriteString("{GRANT")
	sb.WriteString(fmt.Sprintf(" :stmtType %d", n.StmtType))
	if n.Privileges != nil {
		sb.WriteString(" :privileges ")
		writeNode(sb, n.Privileges)
	}
	if n.OnType != "" {
		sb.WriteString(fmt.Sprintf(" :onType \"%s\"", escapeString(n.OnType)))
	}
	if n.OnName != nil {
		sb.WriteString(" :onName ")
		writeNode(sb, n.OnName)
	}
	if n.Principals != nil {
		sb.WriteString(" :principals ")
		writeNode(sb, n.Principals)
	}
	if n.WithGrant {
		sb.WriteString(" :withGrant true")
	}
	if n.CascadeOpt {
		sb.WriteString(" :cascade true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeBinaryExpr(sb *strings.Builder, n *BinaryExpr) {
	sb.WriteString("{BINEXPR")
	sb.WriteString(fmt.Sprintf(" :op %d", n.Op))
	if n.Left != nil {
		sb.WriteString(" :left ")
		writeNode(sb, n.Left)
	}
	if n.Right != nil {
		sb.WriteString(" :right ")
		writeNode(sb, n.Right)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUnaryExpr(sb *strings.Builder, n *UnaryExpr) {
	sb.WriteString("{UNARYEXPR")
	sb.WriteString(fmt.Sprintf(" :op %d", n.Op))
	if n.Operand != nil {
		sb.WriteString(" :operand ")
		writeNode(sb, n.Operand)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFuncCallExpr(sb *strings.Builder, n *FuncCallExpr) {
	sb.WriteString("{FUNCCALL")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.Star {
		sb.WriteString(" :star true")
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.Over != nil {
		sb.WriteString(" :over ")
		writeNode(sb, n.Over)
	}
	if n.Within != nil {
		sb.WriteString(" :within ")
		writeNode(sb, n.Within)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCaseExpr(sb *strings.Builder, n *CaseExpr) {
	sb.WriteString("{CASE")
	if n.Arg != nil {
		sb.WriteString(" :arg ")
		writeNode(sb, n.Arg)
	}
	if n.WhenList != nil {
		sb.WriteString(" :whenList ")
		writeNode(sb, n.WhenList)
	}
	if n.ElseExpr != nil {
		sb.WriteString(" :else ")
		writeNode(sb, n.ElseExpr)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCaseWhen(sb *strings.Builder, n *CaseWhen) {
	sb.WriteString("{WHEN")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Result != nil {
		sb.WriteString(" :result ")
		writeNode(sb, n.Result)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeBetweenExpr(sb *strings.Builder, n *BetweenExpr) {
	sb.WriteString("{BETWEEN")
	sb.WriteString(fmt.Sprintf(" :not %t", n.Not))
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Low != nil {
		sb.WriteString(" :low ")
		writeNode(sb, n.Low)
	}
	if n.High != nil {
		sb.WriteString(" :high ")
		writeNode(sb, n.High)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeInExpr(sb *strings.Builder, n *InExpr) {
	sb.WriteString("{IN")
	sb.WriteString(fmt.Sprintf(" :not %t", n.Not))
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.List != nil {
		sb.WriteString(" :list ")
		writeNode(sb, n.List)
	}
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLikeExpr(sb *strings.Builder, n *LikeExpr) {
	sb.WriteString("{LIKE")
	sb.WriteString(fmt.Sprintf(" :not %t", n.Not))
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Pattern != nil {
		sb.WriteString(" :pattern ")
		writeNode(sb, n.Pattern)
	}
	if n.Escape != nil {
		sb.WriteString(" :escape ")
		writeNode(sb, n.Escape)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIsExpr(sb *strings.Builder, n *IsExpr) {
	sb.WriteString("{IS")
	sb.WriteString(fmt.Sprintf(" :testType %d", n.TestType))
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExistsExpr(sb *strings.Builder, n *ExistsExpr) {
	sb.WriteString("{EXISTS")
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCastExpr(sb *strings.Builder, n *CastExpr) {
	sb.WriteString("{CAST")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeConvertExpr(sb *strings.Builder, n *ConvertExpr) {
	sb.WriteString("{CONVERT")
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Style != nil {
		sb.WriteString(" :style ")
		writeNode(sb, n.Style)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTryCastExpr(sb *strings.Builder, n *TryCastExpr) {
	sb.WriteString("{TRY_CAST")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTryConvertExpr(sb *strings.Builder, n *TryConvertExpr) {
	sb.WriteString("{TRY_CONVERT")
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Style != nil {
		sb.WriteString(" :style ")
		writeNode(sb, n.Style)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCoalesceExpr(sb *strings.Builder, n *CoalesceExpr) {
	sb.WriteString("{COALESCE")
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeNullifExpr(sb *strings.Builder, n *NullifExpr) {
	sb.WriteString("{NULLIF")
	if n.Left != nil {
		sb.WriteString(" :left ")
		writeNode(sb, n.Left)
	}
	if n.Right != nil {
		sb.WriteString(" :right ")
		writeNode(sb, n.Right)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIifExpr(sb *strings.Builder, n *IifExpr) {
	sb.WriteString("{IIF")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.TrueVal != nil {
		sb.WriteString(" :trueVal ")
		writeNode(sb, n.TrueVal)
	}
	if n.FalseVal != nil {
		sb.WriteString(" :falseVal ")
		writeNode(sb, n.FalseVal)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeColumnRef(sb *strings.Builder, n *ColumnRef) {
	sb.WriteString("{COLREF")
	if n.Server != "" {
		sb.WriteString(fmt.Sprintf(" :server \"%s\"", escapeString(n.Server)))
	}
	if n.Database != "" {
		sb.WriteString(fmt.Sprintf(" :database \"%s\"", escapeString(n.Database)))
	}
	if n.Schema != "" {
		sb.WriteString(fmt.Sprintf(" :schema \"%s\"", escapeString(n.Schema)))
	}
	if n.Table != "" {
		sb.WriteString(fmt.Sprintf(" :table \"%s\"", escapeString(n.Table)))
	}
	if n.Column != "" {
		sb.WriteString(fmt.Sprintf(" :column \"%s\"", escapeString(n.Column)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeStarExpr(sb *strings.Builder, n *StarExpr) {
	sb.WriteString("{STAR")
	if n.Qualifier != "" {
		sb.WriteString(fmt.Sprintf(" :qualifier \"%s\"", escapeString(n.Qualifier)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLiteral(sb *strings.Builder, n *Literal) {
	sb.WriteString("{LITERAL")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	switch n.Type {
	case LitString:
		sb.WriteString(fmt.Sprintf(" :str \"%s\"", escapeString(n.Str)))
		if n.IsNChar {
			sb.WriteString(" :nchar true")
		}
	case LitInteger:
		sb.WriteString(fmt.Sprintf(" :ival %d", n.Ival))
	case LitFloat:
		sb.WriteString(fmt.Sprintf(" :str \"%s\"", n.Str))
	case LitNull:
		sb.WriteString(" :null true")
	case LitDefault:
		sb.WriteString(" :default true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSubqueryExpr(sb *strings.Builder, n *SubqueryExpr) {
	sb.WriteString("{SUBQUERY")
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeParenExpr(sb *strings.Builder, n *ParenExpr) {
	sb.WriteString("{PAREN")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTableRef(sb *strings.Builder, n *TableRef) {
	sb.WriteString("{TABLEREF")
	if n.Server != "" {
		sb.WriteString(fmt.Sprintf(" :server \"%s\"", escapeString(n.Server)))
	}
	if n.Database != "" {
		sb.WriteString(fmt.Sprintf(" :database \"%s\"", escapeString(n.Database)))
	}
	if n.Schema != "" {
		sb.WriteString(fmt.Sprintf(" :schema \"%s\"", escapeString(n.Schema)))
	}
	if n.Object != "" {
		sb.WriteString(fmt.Sprintf(" :object \"%s\"", escapeString(n.Object)))
	}
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias \"%s\"", escapeString(n.Alias)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDataType(sb *strings.Builder, n *DataType) {
	sb.WriteString("{DATATYPE")
	if n.Schema != "" {
		sb.WriteString(fmt.Sprintf(" :schema \"%s\"", escapeString(n.Schema)))
	}
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Length != nil {
		sb.WriteString(" :length ")
		writeNode(sb, n.Length)
	}
	if n.MaxLength {
		sb.WriteString(" :max true")
	}
	if n.Precision != nil {
		sb.WriteString(" :precision ")
		writeNode(sb, n.Precision)
	}
	if n.Scale != nil {
		sb.WriteString(" :scale ")
		writeNode(sb, n.Scale)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeJoinClause(sb *strings.Builder, n *JoinClause) {
	sb.WriteString("{JOIN")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Left != nil {
		sb.WriteString(" :left ")
		writeNode(sb, n.Left)
	}
	if n.Right != nil {
		sb.WriteString(" :right ")
		writeNode(sb, n.Right)
	}
	if n.Condition != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.Condition)
	}
	if n.Using != nil {
		sb.WriteString(" :using ")
		writeNode(sb, n.Using)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAliasedTableRef(sb *strings.Builder, n *AliasedTableRef) {
	sb.WriteString("{ALIASEDTABLE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias \"%s\"", escapeString(n.Alias)))
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.TableSample != nil {
		sb.WriteString(" :tablesample ")
		writeNode(sb, n.TableSample)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeOverClause(sb *strings.Builder, n *OverClause) {
	sb.WriteString("{OVER")
	if n.PartitionBy != nil {
		sb.WriteString(" :partitionBy ")
		writeNode(sb, n.PartitionBy)
	}
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	if n.WindowFrame != nil {
		sb.WriteString(" :windowFrame ")
		writeNode(sb, n.WindowFrame)
	}
	if n.WindowName != "" {
		sb.WriteString(fmt.Sprintf(" :windowName \"%s\"", escapeString(n.WindowName)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeOrderByItem(sb *strings.Builder, n *OrderByItem) {
	sb.WriteString("{ORDERBY")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.SortDir != SortDefault {
		sb.WriteString(fmt.Sprintf(" :dir %d", n.SortDir))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeResTarget(sb *strings.Builder, n *ResTarget) {
	sb.WriteString("{RESTARGET")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Val != nil {
		sb.WriteString(" :val ")
		writeNode(sb, n.Val)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetExpr(sb *strings.Builder, n *SetExpr) {
	sb.WriteString("{SETEXPR")
	if n.Column != nil {
		sb.WriteString(" :column ")
		writeNode(sb, n.Column)
	}
	if n.Variable != "" {
		sb.WriteString(fmt.Sprintf(" :variable \"%s\"", escapeString(n.Variable)))
	}
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeValuesClause(sb *strings.Builder, n *ValuesClause) {
	sb.WriteString("{VALUES")
	if n.Rows != nil {
		sb.WriteString(" :rows ")
		writeNode(sb, n.Rows)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTopClause(sb *strings.Builder, n *TopClause) {
	sb.WriteString("{TOP")
	if n.Count != nil {
		sb.WriteString(" :count ")
		writeNode(sb, n.Count)
	}
	if n.Percent {
		sb.WriteString(" :percent true")
	}
	if n.WithTies {
		sb.WriteString(" :withTies true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWithClause(sb *strings.Builder, n *WithClause) {
	sb.WriteString("{WITH")
	if n.CTEs != nil {
		sb.WriteString(" :ctes ")
		writeNode(sb, n.CTEs)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCommonTableExpr(sb *strings.Builder, n *CommonTableExpr) {
	sb.WriteString("{CTE")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeOutputClause(sb *strings.Builder, n *OutputClause) {
	sb.WriteString("{OUTPUT")
	if n.Targets != nil {
		sb.WriteString(" :targets ")
		writeNode(sb, n.Targets)
	}
	if n.IntoTable != nil {
		sb.WriteString(" :intoTable ")
		writeNode(sb, n.IntoTable)
	}
	if n.IntoCols != nil {
		sb.WriteString(" :intoCols ")
		writeNode(sb, n.IntoCols)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeColumnDef(sb *strings.Builder, n *ColumnDef) {
	sb.WriteString("{COLUMNDEF")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Identity != nil {
		sb.WriteString(" :identity ")
		writeNode(sb, n.Identity)
	}
	if n.Computed != nil {
		sb.WriteString(" :computed ")
		writeNode(sb, n.Computed)
	}
	if n.DefaultExpr != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.DefaultExpr)
	}
	if n.Collation != "" {
		sb.WriteString(fmt.Sprintf(" :collation \"%s\"", escapeString(n.Collation)))
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	if n.Nullable != nil {
		sb.WriteString(" :nullable ")
		writeNode(sb, n.Nullable)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeConstraintDef(sb *strings.Builder, n *ConstraintDef) {
	sb.WriteString("{CONSTRAINT")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.RefTable != nil {
		sb.WriteString(" :refTable ")
		writeNode(sb, n.RefTable)
	}
	if n.RefColumns != nil {
		sb.WriteString(" :refColumns ")
		writeNode(sb, n.RefColumns)
	}
	if n.OnDelete != RefActNone {
		sb.WriteString(fmt.Sprintf(" :onDelete %d", n.OnDelete))
	}
	if n.OnUpdate != RefActNone {
		sb.WriteString(fmt.Sprintf(" :onUpdate %d", n.OnUpdate))
	}
	if n.Clustered != nil {
		sb.WriteString(fmt.Sprintf(" :clustered %t", *n.Clustered))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeWhenClause(sb *strings.Builder, n *MergeWhenClause) {
	sb.WriteString("{MERGEWHEN")
	sb.WriteString(fmt.Sprintf(" :matched %t", n.Matched))
	sb.WriteString(fmt.Sprintf(" :byTarget %t", n.ByTarget))
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Action != nil {
		sb.WriteString(" :action ")
		writeNode(sb, n.Action)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIndexColumn(sb *strings.Builder, n *IndexColumn) {
	sb.WriteString("{INDEXCOL")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.SortDir != SortDefault {
		sb.WriteString(fmt.Sprintf(" :dir %d", n.SortDir))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeParamDef(sb *strings.Builder, n *ParamDef) {
	sb.WriteString("{PARAM")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	if n.Output {
		sb.WriteString(" :output true")
	}
	if n.ReadOnly {
		sb.WriteString(" :readOnly true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeVariableDecl(sb *strings.Builder, n *VariableDecl) {
	sb.WriteString("{VARDECL")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	if n.IsTable {
		sb.WriteString(" :isTable true")
	}
	if n.TableDef != nil {
		sb.WriteString(" :tableDef ")
		writeNode(sb, n.TableDef)
	}
	if n.IsCursor {
		sb.WriteString(" :isCursor true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSelectAssign(sb *strings.Builder, n *SelectAssign) {
	sb.WriteString("{SELECTASSIGN")
	sb.WriteString(fmt.Sprintf(" :variable \"%s\"", escapeString(n.Variable)))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMethodCallExpr(sb *strings.Builder, n *MethodCallExpr) {
	sb.WriteString("{METHODCALL")
	if n.Type != nil {
		sb.WriteString(" :type ")
		writeNode(sb, n.Type)
	}
	sb.WriteString(fmt.Sprintf(" :method \"%s\"", escapeString(n.Method)))
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFetchClause(sb *strings.Builder, n *FetchClause) {
	sb.WriteString("{FETCH")
	if n.Count != nil {
		sb.WriteString(" :count ")
		writeNode(sb, n.Count)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeForClause(sb *strings.Builder, n *ForClause) {
	sb.WriteString("{FOR")
	sb.WriteString(fmt.Sprintf(" :mode %d", n.Mode))
	if n.SubMode != "" {
		sb.WriteString(fmt.Sprintf(" :subMode \"%s\"", escapeString(n.SubMode)))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeNullableSpec(sb *strings.Builder, n *NullableSpec) {
	sb.WriteString("{NULLABLE")
	sb.WriteString(fmt.Sprintf(" :notNull %t", n.NotNull))
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIdentitySpec(sb *strings.Builder, n *IdentitySpec) {
	sb.WriteString("{IDENTITY")
	sb.WriteString(fmt.Sprintf(" :seed %d", n.Seed))
	sb.WriteString(fmt.Sprintf(" :increment %d", n.Increment))
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeComputedColumnDef(sb *strings.Builder, n *ComputedColumnDef) {
	sb.WriteString("{COMPUTED")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Persisted {
		sb.WriteString(" :persisted true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTableAction(sb *strings.Builder, n *AlterTableAction) {
	sb.WriteString("{ALTERACTION")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Column != nil {
		sb.WriteString(" :column ")
		writeNode(sb, n.Column)
	}
	if n.ColName != "" {
		sb.WriteString(fmt.Sprintf(" :colName \"%s\"", escapeString(n.ColName)))
	}
	if n.Constraint != nil {
		sb.WriteString(" :constraint ")
		writeNode(sb, n.Constraint)
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeReturnsTableDef(sb *strings.Builder, n *ReturnsTableDef) {
	sb.WriteString("{RETURNSTABLE")
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Variable != "" {
		sb.WriteString(fmt.Sprintf(" :variable \"%s\"", escapeString(n.Variable)))
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExecArg(sb *strings.Builder, n *ExecArg) {
	sb.WriteString("{EXECARG")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	}
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	if n.Output {
		sb.WriteString(" :output true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeUpdateAction(sb *strings.Builder, n *MergeUpdateAction) {
	sb.WriteString("{MERGEUPDATE")
	if n.SetClause != nil {
		sb.WriteString(" :set ")
		writeNode(sb, n.SetClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeDeleteAction(sb *strings.Builder, n *MergeDeleteAction) {
	sb.WriteString(fmt.Sprintf("{MERGEDELETE :loc %d %d}", n.Loc.Start, n.Loc.End))
}

func writeMergeInsertAction(sb *strings.Builder, n *MergeInsertAction) {
	sb.WriteString("{MERGEINSERT")
	if n.Cols != nil {
		sb.WriteString(" :cols ")
		writeNode(sb, n.Cols)
	}
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWindowFrame(sb *strings.Builder, n *WindowFrame) {
	sb.WriteString("{WINDOWFRAME")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Start != nil {
		sb.WriteString(" :start ")
		writeNode(sb, n.Start)
	}
	if n.End != nil {
		sb.WriteString(" :end ")
		writeNode(sb, n.End)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWindowBound(sb *strings.Builder, n *WindowBound) {
	sb.WriteString("{WINDOWBOUND")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Offset != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.Offset)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDeclareCursorStmt(sb *strings.Builder, n *DeclareCursorStmt) {
	sb.WriteString("{DECLARECURSOR")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Insensitive {
		sb.WriteString(" :insensitive true")
	}
	if n.Scroll {
		sb.WriteString(" :scroll true")
	}
	if n.Scope != "" {
		sb.WriteString(fmt.Sprintf(" :scope \"%s\"", n.Scope))
	}
	if n.ForwardOnly {
		sb.WriteString(" :forwardOnly true")
	}
	if n.CursorType != "" {
		sb.WriteString(fmt.Sprintf(" :cursorType \"%s\"", n.CursorType))
	}
	if n.Concurrency != "" {
		sb.WriteString(fmt.Sprintf(" :concurrency \"%s\"", n.Concurrency))
	}
	if n.TypeWarning {
		sb.WriteString(" :typeWarning true")
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	if n.ForUpdate {
		sb.WriteString(" :forUpdate true")
	}
	if n.UpdateCols != nil {
		sb.WriteString(" :updateCols ")
		writeNode(sb, n.UpdateCols)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeOpenCursorStmt(sb *strings.Builder, n *OpenCursorStmt) {
	sb.WriteString("{OPENCURSOR")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Global {
		sb.WriteString(" :global true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFetchCursorStmt(sb *strings.Builder, n *FetchCursorStmt) {
	sb.WriteString("{FETCHCURSOR")
	if n.Orientation != "" {
		sb.WriteString(fmt.Sprintf(" :orientation \"%s\"", n.Orientation))
	}
	if n.FetchOffset != nil {
		sb.WriteString(" :fetchOffset ")
		writeNode(sb, n.FetchOffset)
	}
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Global {
		sb.WriteString(" :global true")
	}
	if n.IntoVars != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.IntoVars)
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCloseCursorStmt(sb *strings.Builder, n *CloseCursorStmt) {
	sb.WriteString("{CLOSECURSOR")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Global {
		sb.WriteString(" :global true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDeallocateCursorStmt(sb *strings.Builder, n *DeallocateCursorStmt) {
	sb.WriteString("{DEALLOCATECURSOR")
	sb.WriteString(fmt.Sprintf(" :name \"%s\"", escapeString(n.Name)))
	if n.Global {
		sb.WriteString(" :global true")
	}
	sb.WriteString(fmt.Sprintf(" :loc %d %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterDatabaseStmt(sb *strings.Builder, n *AlterDatabaseStmt) {
	sb.WriteString("{ALTERDATABASE")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeAlterIndexStmt(sb *strings.Builder, n *AlterIndexStmt) {
	sb.WriteString("{ALTERINDEX")
	fmt.Fprintf(sb, " :indexName \"%s\"", escapeString(n.IndexName))
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeBulkInsertStmt(sb *strings.Builder, n *BulkInsertStmt) {
	sb.WriteString("{BULKINSERT")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	fmt.Fprintf(sb, " :dataFile \"%s\"", escapeString(n.DataFile))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeDbccStmt(sb *strings.Builder, n *DbccStmt) {
	sb.WriteString("{DBCC")
	fmt.Fprintf(sb, " :command \"%s\"", escapeString(n.Command))
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeBackupStmt(sb *strings.Builder, n *BackupStmt) {
	sb.WriteString("{BACKUP")
	fmt.Fprintf(sb, " :type \"%s\"", escapeString(n.Type))
	fmt.Fprintf(sb, " :database \"%s\"", escapeString(n.Database))
	fmt.Fprintf(sb, " :target \"%s\"", escapeString(n.Target))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeRestoreStmt(sb *strings.Builder, n *RestoreStmt) {
	sb.WriteString("{RESTORE")
	fmt.Fprintf(sb, " :type \"%s\"", escapeString(n.Type))
	fmt.Fprintf(sb, " :database \"%s\"", escapeString(n.Database))
	fmt.Fprintf(sb, " :source \"%s\"", escapeString(n.Source))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

func writeSecurityKeyStmt(sb *strings.Builder, n *SecurityKeyStmt) {
	sb.WriteString("{SECURITYKEY")
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	fmt.Fprintf(sb, " :objectType \"%s\"", escapeString(n.ObjectType))
	if n.Name != "" {
		fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d", n.Loc.Start, n.Loc.End)
	sb.WriteString("}")
}

// ---------- Batch 39-48 write functions ----------

func writeCreateStatisticsStmt(sb *strings.Builder, n *CreateStatisticsStmt) {
	sb.WriteString("{CREATESTATISTICS")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeUpdateStatisticsStmt(sb *strings.Builder, n *UpdateStatisticsStmt) {
	sb.WriteString("{UPDATESTATISTICS")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Name != "" {
		fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeDropStatisticsStmt(sb *strings.Builder, n *DropStatisticsStmt) {
	sb.WriteString("{DROPSTATISTICS")
	if n.Names != nil {
		sb.WriteString(" :names ")
		writeNode(sb, n.Names)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeSetOptionStmt(sb *strings.Builder, n *SetOptionStmt) {
	sb.WriteString("{SETOPTION")
	fmt.Fprintf(sb, " :option \"%s\"", escapeString(n.Option))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreatePartitionFunctionStmt(sb *strings.Builder, n *CreatePartitionFunctionStmt) {
	sb.WriteString("{CREATEPARTITIONFUNCTION")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.InputType != nil {
		sb.WriteString(" :inputType ")
		writeNode(sb, n.InputType)
	}
	fmt.Fprintf(sb, " :range \"%s\"", escapeString(n.Range))
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeAlterPartitionFunctionStmt(sb *strings.Builder, n *AlterPartitionFunctionStmt) {
	sb.WriteString("{ALTERPARTITIONFUNCTION")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	if n.BoundaryValue != nil {
		sb.WriteString(" :boundaryValue ")
		writeNode(sb, n.BoundaryValue)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreatePartitionSchemeStmt(sb *strings.Builder, n *CreatePartitionSchemeStmt) {
	sb.WriteString("{CREATEPARTITIONSCHEME")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	fmt.Fprintf(sb, " :functionName \"%s\"", escapeString(n.FunctionName))
	if n.AllToFileGroup != "" {
		fmt.Fprintf(sb, " :allTo \"%s\"", escapeString(n.AllToFileGroup))
	}
	if n.FileGroups != nil {
		sb.WriteString(" :fileGroups ")
		writeNode(sb, n.FileGroups)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreateFulltextIndexStmt(sb *strings.Builder, n *CreateFulltextIndexStmt) {
	sb.WriteString("{CREATEFULLTEXTINDEX")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	fmt.Fprintf(sb, " :keyIndex \"%s\"", escapeString(n.KeyIndex))
	if n.CatalogName != "" {
		fmt.Fprintf(sb, " :catalog \"%s\"", escapeString(n.CatalogName))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeAlterFulltextIndexStmt(sb *strings.Builder, n *AlterFulltextIndexStmt) {
	sb.WriteString("{ALTERFULLTEXTINDEX")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreateFulltextCatalogStmt(sb *strings.Builder, n *CreateFulltextCatalogStmt) {
	sb.WriteString("{CREATEFULLTEXTCATALOG")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreateXmlSchemaCollectionStmt(sb *strings.Builder, n *CreateXmlSchemaCollectionStmt) {
	sb.WriteString("{CREATEXMLSCHEMACOLLECTION")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.XmlSchemaNamespaces != nil {
		sb.WriteString(" :schema ")
		writeNode(sb, n.XmlSchemaNamespaces)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeAlterXmlSchemaCollectionStmt(sb *strings.Builder, n *AlterXmlSchemaCollectionStmt) {
	sb.WriteString("{ALTERXMLSCHEMACOLLECTION")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.XmlSchemaNamespaces != nil {
		sb.WriteString(" :schema ")
		writeNode(sb, n.XmlSchemaNamespaces)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCreateAssemblyStmt(sb *strings.Builder, n *CreateAssemblyStmt) {
	sb.WriteString("{CREATEASSEMBLY")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Authorization != "" {
		fmt.Fprintf(sb, " :authorization \"%s\"", escapeString(n.Authorization))
	}
	if n.FromFiles != nil {
		sb.WriteString(" :fromFiles ")
		writeNode(sb, n.FromFiles)
	}
	if n.PermissionSet != "" {
		fmt.Fprintf(sb, " :permissionSet \"%s\"", escapeString(n.PermissionSet))
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeAlterAssemblyStmt(sb *strings.Builder, n *AlterAssemblyStmt) {
	sb.WriteString("{ALTERASSEMBLY")
	fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	if n.Actions != nil {
		sb.WriteString(" :actions ")
		writeNode(sb, n.Actions)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeServiceBrokerStmt(sb *strings.Builder, n *ServiceBrokerStmt) {
	sb.WriteString("{SERVICEBROKER")
	fmt.Fprintf(sb, " :action \"%s\"", escapeString(n.Action))
	fmt.Fprintf(sb, " :objectType \"%s\"", escapeString(n.ObjectType))
	if n.Name != "" {
		fmt.Fprintf(sb, " :name \"%s\"", escapeString(n.Name))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCheckpointStmt(sb *strings.Builder, n *CheckpointStmt) {
	sb.WriteString("{CHECKPOINT")
	if n.Duration != nil {
		sb.WriteString(" :duration ")
		writeNode(sb, n.Duration)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeKillStmt(sb *strings.Builder, n *KillStmt) {
	sb.WriteString("{KILL")
	if n.SessionID != nil {
		sb.WriteString(" :sessionId ")
		writeNode(sb, n.SessionID)
	}
	fmt.Fprintf(sb, " :statusOnly %v :loc %d %d}", n.StatusOnly, n.Loc.Start, n.Loc.End)
}

func writeReadtextStmt(sb *strings.Builder, n *ReadtextStmt) {
	sb.WriteString("{READTEXT")
	if n.Column != nil {
		sb.WriteString(" :column ")
		writeNode(sb, n.Column)
	}
	if n.TextPtr != nil {
		sb.WriteString(" :textPtr ")
		writeNode(sb, n.TextPtr)
	}
	if n.Offset != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.Offset)
	}
	if n.Size != nil {
		sb.WriteString(" :size ")
		writeNode(sb, n.Size)
	}
	fmt.Fprintf(sb, " :holdLock %v :loc %d %d}", n.HoldLock, n.Loc.Start, n.Loc.End)
}

func writeWritetextStmt(sb *strings.Builder, n *WritetextStmt) {
	sb.WriteString("{WRITETEXT")
	if n.Column != nil {
		sb.WriteString(" :column ")
		writeNode(sb, n.Column)
	}
	if n.TextPtr != nil {
		sb.WriteString(" :textPtr ")
		writeNode(sb, n.TextPtr)
	}
	if n.Data != nil {
		sb.WriteString(" :data ")
		writeNode(sb, n.Data)
	}
	fmt.Fprintf(sb, " :withLog %v :loc %d %d}", n.WithLog, n.Loc.Start, n.Loc.End)
}

func writeUpdatetextStmt(sb *strings.Builder, n *UpdatetextStmt) {
	sb.WriteString("{UPDATETEXT")
	if n.DestColumn != nil {
		sb.WriteString(" :destColumn ")
		writeNode(sb, n.DestColumn)
	}
	if n.DestTextPtr != nil {
		sb.WriteString(" :destTextPtr ")
		writeNode(sb, n.DestTextPtr)
	}
	if n.InsertOffset != nil {
		sb.WriteString(" :insertOffset ")
		writeNode(sb, n.InsertOffset)
	}
	if n.DeleteLength != nil {
		sb.WriteString(" :deleteLength ")
		writeNode(sb, n.DeleteLength)
	}
	if n.InsertedData != nil {
		sb.WriteString(" :insertedData ")
		writeNode(sb, n.InsertedData)
	}
	fmt.Fprintf(sb, " :withLog %v :loc %d %d}", n.WithLog, n.Loc.Start, n.Loc.End)
}

func writePivotExpr(sb *strings.Builder, n *PivotExpr) {
	sb.WriteString("{PIVOT")
	if n.Source != nil {
		sb.WriteString(" :source ")
		writeNode(sb, n.Source)
	}
	if n.AggFunc != nil {
		sb.WriteString(" :aggFunc ")
		writeNode(sb, n.AggFunc)
	}
	if n.ForCol != "" {
		fmt.Fprintf(sb, " :forCol \"%s\"", escapeString(n.ForCol))
	}
	if n.InValues != nil {
		sb.WriteString(" :inValues ")
		writeNode(sb, n.InValues)
	}
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias \"%s\"", escapeString(n.Alias))
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeUnpivotExpr(sb *strings.Builder, n *UnpivotExpr) {
	sb.WriteString("{UNPIVOT")
	if n.Source != nil {
		sb.WriteString(" :source ")
		writeNode(sb, n.Source)
	}
	if n.ValueCol != "" {
		fmt.Fprintf(sb, " :valueCol \"%s\"", escapeString(n.ValueCol))
	}
	if n.ForCol != "" {
		fmt.Fprintf(sb, " :forCol \"%s\"", escapeString(n.ForCol))
	}
	if n.InCols != nil {
		sb.WriteString(" :inCols ")
		writeNode(sb, n.InCols)
	}
	if n.Alias != "" {
		fmt.Fprintf(sb, " :alias \"%s\"", escapeString(n.Alias))
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeTableSampleClause(sb *strings.Builder, n *TableSampleClause) {
	sb.WriteString("{TABLESAMPLE")
	if n.Size != nil {
		sb.WriteString(" :size ")
		writeNode(sb, n.Size)
	}
	if n.Unit != "" {
		fmt.Fprintf(sb, " :unit \"%s\"", escapeString(n.Unit))
	}
	if n.Repeatable != nil {
		sb.WriteString(" :repeatable ")
		writeNode(sb, n.Repeatable)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeGroupingSetsExpr(sb *strings.Builder, n *GroupingSetsExpr) {
	sb.WriteString("{GROUPINGSETS")
	if n.Sets != nil {
		sb.WriteString(" :sets ")
		writeNode(sb, n.Sets)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeRollupExpr(sb *strings.Builder, n *RollupExpr) {
	sb.WriteString("{ROLLUP")
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func writeCubeExpr(sb *strings.Builder, n *CubeExpr) {
	sb.WriteString("{CUBE")
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	fmt.Fprintf(sb, " :loc %d %d}", n.Loc.Start, n.Loc.End)
}

func escapeString(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
