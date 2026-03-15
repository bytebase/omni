package ast

import (
	"fmt"
	"strings"
)

// NodeToString converts a Node to its string representation.
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

	// Literal nodes
	case *NullLiteral:
		sb.WriteString("{NULL")
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
		sb.WriteString("}")
	case *StringLiteral:
		sb.WriteString("{STRLIT")
		sb.WriteString(fmt.Sprintf(" :val %q", n.Val))
		if n.IsNChar {
			sb.WriteString(" :isnchar true")
		}
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
		sb.WriteString("}")
	case *NumberLiteral:
		sb.WriteString("{NUMLIT")
		sb.WriteString(fmt.Sprintf(" :val %q", n.Val))
		if n.IsFloat {
			sb.WriteString(" :isfloat true")
		}
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
		sb.WriteString("}")

	// Name/reference nodes
	case *ColumnRef:
		writeColumnRef(sb, n)
	case *ObjectName:
		writeObjectName(sb, n)
	case *Alias:
		writeAlias(sb, n)
	case *TypeName:
		writeTypeName(sb, n)
	case *BindVariable:
		sb.WriteString(fmt.Sprintf("{BINDVAR :name %q :loc_start %d :loc_end %d}", n.Name, n.Loc.Start, n.Loc.End))
	case *PseudoColumn:
		sb.WriteString(fmt.Sprintf("{PSEUDOCOL :type %d :loc_start %d :loc_end %d}", n.Type, n.Loc.Start, n.Loc.End))
	case *Hint:
		writeHint(sb, n)
	case *Star:
		sb.WriteString("{STAR")
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
		sb.WriteString("}")

	// Expression nodes
	case *BinaryExpr:
		writeBinaryExpr(sb, n)
	case *UnaryExpr:
		writeUnaryExpr(sb, n)
	case *BoolExpr:
		writeBoolExpr(sb, n)
	case *FuncCallExpr:
		writeFuncCallExpr(sb, n)
	case *CaseExpr:
		writeCaseExpr(sb, n)
	case *CaseWhen:
		writeCaseWhen(sb, n)
	case *DecodeExpr:
		writeDecodeExpr(sb, n)
	case *DecodePair:
		writeDecodePair(sb, n)
	case *CastExpr:
		writeCastExpr(sb, n)
	case *MultisetExpr:
		writeMultisetExpr(sb, n)
	case *CursorExpr:
		writeCursorExpr(sb, n)
	case *TreatExpr:
		writeTreatExpr(sb, n)
	case *IntervalExpr:
		writeIntervalExpr(sb, n)
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
	case *SubqueryExpr:
		writeSubqueryExpr(sb, n)
	case *ParenExpr:
		writeParenExpr(sb, n)

	// Window nodes
	case *WindowSpec:
		writeWindowSpec(sb, n)
	case *WindowFrame:
		writeWindowFrame(sb, n)
	case *WindowBound:
		writeWindowBound(sb, n)
	case *KeepClause:
		writeKeepClause(sb, n)

	// Grouping extension nodes
	case *GroupingSetsClause:
		writeGroupingSetsClause(sb, n)
	case *CubeClause:
		writeCubeClause(sb, n)
	case *RollupClause:
		writeRollupClause(sb, n)

	// Table nodes
	case *JoinClause:
		writeJoinClause(sb, n)
	case *TableRef:
		writeTableRef(sb, n)
	case *SubqueryRef:
		writeSubqueryRef(sb, n)
	case *LateralRef:
		writeLateralRef(sb, n)
	case *XmlTableRef:
		writeXmlTableRef(sb, n)
	case *XmlTableColumn:
		writeXmlTableColumn(sb, n)
	case *JsonTableRef:
		writeJsonTableRef(sb, n)
	case *JsonTableColumn:
		writeJsonTableColumn(sb, n)
	case *SampleClause:
		writeSampleClause(sb, n)
	case *PivotClause:
		writePivotClause(sb, n)
	case *UnpivotClause:
		writeUnpivotClause(sb, n)
	case *TableCollectionExpr:
		writeTableCollectionExpr(sb, n)
	case *MatchRecognizeClause:
		writeMatchRecognizeClause(sb, n)
	case *ContainersExpr:
		writeContainersExpr(sb, n)
	case *InlineExternalTable:
		writeInlineExternalTable(sb, n)

	// Clause nodes
	case *HierarchicalClause:
		writeHierarchicalClause(sb, n)
	case *WithClause:
		writeWithClause(sb, n)
	case *CTE:
		writeCTE(sb, n)
	case *CTESearchClause:
		writeCTESearchClause(sb, n)
	case *CTECycleClause:
		writeCTECycleClause(sb, n)
	case *ForUpdateClause:
		writeForUpdateClause(sb, n)
	case *FetchFirstClause:
		writeFetchFirstClause(sb, n)
	case *ModelClause:
		writeModelClause(sb, n)
	case *ModelCellRefOptions:
		writeModelCellRefOptions(sb, n)
	case *ModelRefModel:
		writeModelRefModel(sb, n)
	case *ModelMainModel:
		writeModelMainModel(sb, n)
	case *ModelColumnClauses:
		writeModelColumnClauses(sb, n)
	case *ModelRulesClause:
		writeModelRulesClause(sb, n)
	case *ModelRule:
		writeModelRule(sb, n)
	case *ModelForLoop:
		writeModelForLoop(sb, n)
	case *FlashbackClause:
		writeFlashbackClause(sb, n)
	case *WindowDef:
		writeWindowDef(sb, n)
	case *PartitionExtClause:
		writePartitionExtClause(sb, n)
	case *SortBy:
		writeSortBy(sb, n)
	case *ResTarget:
		writeResTarget(sb, n)

	// DML helpers
	case *InsertIntoClause:
		writeInsertIntoClause(sb, n)
	case *ErrorLogClause:
		writeErrorLogClause(sb, n)
	case *SetClause:
		writeSetClause(sb, n)
	case *MergeClause:
		writeMergeClause(sb, n)

	// DDL nodes
	case *ColumnDef:
		writeColumnDef(sb, n)
	case *ColumnConstraint:
		writeColumnConstraint(sb, n)
	case *TableConstraint:
		writeTableConstraint(sb, n)
	case *IdentityClause:
		writeIdentityClause(sb, n)
	case *StorageClause:
		writeStorageClause(sb, n)
	case *PartitionClause:
		writePartitionClause(sb, n)
	case *PartitionDef:
		writePartitionDef(sb, n)
	case *AlterTableCmd:
		writeAlterTableCmd(sb, n)
	case *IndexColumn:
		writeIndexColumn(sb, n)

	// Statements
	case *RawStmt:
		writeRawStmt(sb, n)
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
	case *CreateSequenceStmt:
		writeCreateSequenceStmt(sb, n)
	case *CreateSynonymStmt:
		writeCreateSynonymStmt(sb, n)
	case *CreateDatabaseLinkStmt:
		writeCreateDatabaseLinkStmt(sb, n)
	case *CreateTypeStmt:
		writeCreateTypeStmt(sb, n)
	case *TypeBodyMember:
		writeTypeBodyMember(sb, n)
	case *CreatePackageStmt:
		writeCreatePackageStmt(sb, n)
	case *CreateProcedureStmt:
		writeCreateProcedureStmt(sb, n)
	case *CreateFunctionStmt:
		writeCreateFunctionStmt(sb, n)
	case *Parameter:
		writeParameter(sb, n)
	case *CreateTriggerStmt:
		writeCreateTriggerStmt(sb, n)
	case *TruncateStmt:
		writeTruncateStmt(sb, n)
	case *GrantStmt:
		writeGrantStmt(sb, n)
	case *RevokeStmt:
		writeRevokeStmt(sb, n)
	case *CommentStmt:
		writeCommentStmt(sb, n)
	case *AlterSessionStmt:
		writeAlterSessionStmt(sb, n)
	case *AlterSystemStmt:
		writeAlterSystemStmt(sb, n)
	case *SetParam:
		writeSetParam(sb, n)
	case *CommitStmt:
		writeCommitStmt(sb, n)
	case *RollbackStmt:
		writeRollbackStmt(sb, n)
	case *SavepointStmt:
		writeSavepointStmt(sb, n)
	case *SetTransactionStmt:
		writeSetTransactionStmt(sb, n)
	case *AnalyzeStmt:
		writeAnalyzeStmt(sb, n)
	case *ExplainPlanStmt:
		writeExplainPlanStmt(sb, n)
	case *FlashbackTableStmt:
		writeFlashbackTableStmt(sb, n)
	case *FlashbackDatabaseStmt:
		writeFlashbackDatabaseStmt(sb, n)
	case *PurgeStmt:
		writePurgeStmt(sb, n)
	case *LockTableStmt:
		writeLockTableStmt(sb, n)
	case *CallStmt:
		writeCallStmt(sb, n)
	case *RenameStmt:
		writeRenameStmt(sb, n)
	case *IdentifiedClause:
		writeIdentifiedClause(sb, n)
	case *UserQuotaClause:
		writeUserQuotaClause(sb, n)
	case *DefaultRoleClause:
		writeDefaultRoleClause(sb, n)
	case *ProfileLimit:
		writeProfileLimit(sb, n)
	case *CreateUserStmt:
		writeCreateUserStmt(sb, n)
	case *AlterUserStmt:
		writeAlterUserStmt(sb, n)
	case *CreateRoleStmt:
		writeCreateRoleStmt(sb, n)
	case *CreateProfileStmt:
		writeCreateProfileStmt(sb, n)
	case *AdminDDLStmt:
		writeAdminDDLStmt(sb, n)
	case *DDLOption:
		writeDDLOption(sb, n)
	case *CreateSchemaStmt:
		writeCreateSchemaStmt(sb, n)
	case *AlterDatabaseLinkStmt:
		writeAlterDatabaseLinkStmt(sb, n)
	case *AlterSynonymStmt:
		writeAlterSynonymStmt(sb, n)
	case *AlterMaterializedViewStmt:
		writeAlterMaterializedViewStmt(sb, n)
	case *AlterIndexStmt:
		writeAlterIndexStmt(sb, n)
	case *AlterViewStmt:
		writeAlterViewStmt(sb, n)
	case *AlterSequenceStmt:
		writeAlterSequenceStmt(sb, n)
	case *AlterProcedureStmt:
		writeAlterProcedureStmt(sb, n)
	case *AlterFunctionStmt:
		writeAlterFunctionStmt(sb, n)
	case *AlterPackageStmt:
		writeAlterPackageStmt(sb, n)
	case *AlterTriggerStmt:
		writeAlterTriggerStmt(sb, n)
	case *AlterTypeStmt:
		writeAlterTypeStmt(sb, n)
	case *TypeAttribute:
		writeTypeAttribute(sb, n)
	case *SetRoleStmt:
		writeSetRoleStmt(sb, n)
	case *SetConstraintsStmt:
		writeSetConstraintsStmt(sb, n)
	case *AuditStmt:
		writeAuditStmt(sb, n)
	case *NoauditStmt:
		writeNoauditStmt(sb, n)
	case *AssociateStatisticsStmt:
		writeAssociateStatisticsStmt(sb, n)
	case *DisassociateStatisticsStmt:
		writeDisassociateStatisticsStmt(sb, n)
	case *CreateTablespaceStmt:
		writeCreateTablespaceStmt(sb, n)
	case *CreateClusterStmt:
		writeCreateClusterStmt(sb, n)
	case *CreateDimensionStmt:
		writeCreateDimensionStmt(sb, n)
	case *DatafileClause:
		writeDatafileClause(sb, n)
	case *AutoextendClause:
		writeAutoextendClause(sb, n)
	case *ClusterColumn:
		writeClusterColumn(sb, n)
	case *DimensionLevel:
		writeDimensionLevel(sb, n)
	case *DimensionHierarchy:
		writeDimensionHierarchy(sb, n)
	case *DimensionJoinKey:
		writeDimensionJoinKey(sb, n)
	case *DimensionAttribute:
		writeDimensionAttribute(sb, n)

	// PL/SQL nodes
	case *PLSQLBlock:
		writePLSQLBlock(sb, n)
	case *ExceptionHandler:
		writeExceptionHandler(sb, n)
	case *PLSQLIf:
		writePLSQLIf(sb, n)
	case *PLSQLElsIf:
		writePLSQLElsIf(sb, n)
	case *PLSQLLoop:
		writePLSQLLoop(sb, n)
	case *PLSQLReturn:
		writePLSQLReturn(sb, n)
	case *PLSQLGoto:
		writePLSQLGoto(sb, n)
	case *PLSQLAssign:
		writePLSQLAssign(sb, n)
	case *PLSQLRaise:
		writePLSQLRaise(sb, n)
	case *PLSQLNull:
		writePLSQLNull(sb, n)
	case *PLSQLVarDecl:
		writePLSQLVarDecl(sb, n)
	case *PLSQLCursorDecl:
		writePLSQLCursorDecl(sb, n)
	case *PLSQLExecImmediate:
		writePLSQLExecImmediate(sb, n)
	case *PLSQLOpen:
		writePLSQLOpen(sb, n)
	case *PLSQLFetch:
		writePLSQLFetch(sb, n)
	case *PLSQLClose:
		writePLSQLClose(sb, n)
	case *PLSQLForall:
		writePLSQLForall(sb, n)
	case *PLSQLExit:
		writePLSQLExit(sb, n)
	case *PLSQLContinue:
		writePLSQLContinue(sb, n)
	case *PLSQLPipeRow:
		writePLSQLPipeRow(sb, n)
	case *PLSQLPragma:
		writePLSQLPragma(sb, n)
	case *PLSQLCase:
		writePLSQLCase(sb, n)
	case *PLSQLWhen:
		writePLSQLWhen(sb, n)
	case *PLSQLTypeDecl:
		writePLSQLTypeDecl(sb, n)
	case *PLSQLCall:
		writePLSQLCall(sb, n)

	default:
		sb.WriteString("{UNKNOWN_NODE}")
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

func writeColumnRef(sb *strings.Builder, n *ColumnRef) {
	sb.WriteString("{COLUMNREF")
	if n.Schema != "" {
		sb.WriteString(fmt.Sprintf(" :schema %q", n.Schema))
	}
	if n.Table != "" {
		sb.WriteString(fmt.Sprintf(" :table %q", n.Table))
	}
	sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
	if n.OuterJoin {
		sb.WriteString(" :outerJoin true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeObjectName(sb *strings.Builder, n *ObjectName) {
	sb.WriteString("{OBJNAME")
	if n.Schema != "" {
		sb.WriteString(fmt.Sprintf(" :schema %q", n.Schema))
	}
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.DBLink != "" {
		sb.WriteString(fmt.Sprintf(" :dblink %q", n.DBLink))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlias(sb *strings.Builder, n *Alias) {
	sb.WriteString("{ALIAS")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Cols != nil {
		sb.WriteString(" :cols ")
		writeNode(sb, n.Cols)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTypeName(sb *strings.Builder, n *TypeName) {
	sb.WriteString("{TYPENAME")
	if n.Names != nil {
		sb.WriteString(" :names ")
		writeNode(sb, n.Names)
	}
	if n.TypeMods != nil {
		sb.WriteString(" :typeMods ")
		writeNode(sb, n.TypeMods)
	}
	if n.IsPercType {
		sb.WriteString(" :pctType true")
	}
	if n.IsPercRowtype {
		sb.WriteString(" :pctRowtype true")
	}
	if n.ArrayBounds != nil {
		sb.WriteString(" :arrayBounds ")
		writeNode(sb, n.ArrayBounds)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeHint(sb *strings.Builder, n *Hint) {
	sb.WriteString("{HINT")
	sb.WriteString(fmt.Sprintf(" :text %q", n.Text))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

func writeBinaryExpr(sb *strings.Builder, n *BinaryExpr) {
	sb.WriteString("{BINEXPR")
	sb.WriteString(fmt.Sprintf(" :op %q", n.Op))
	sb.WriteString(" :left ")
	writeNode(sb, n.Left)
	sb.WriteString(" :right ")
	writeNode(sb, n.Right)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUnaryExpr(sb *strings.Builder, n *UnaryExpr) {
	sb.WriteString("{UNARYEXPR")
	sb.WriteString(fmt.Sprintf(" :op %q", n.Op))
	sb.WriteString(" :operand ")
	writeNode(sb, n.Operand)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeBoolExpr(sb *strings.Builder, n *BoolExpr) {
	sb.WriteString("{BOOLEXPR")
	sb.WriteString(fmt.Sprintf(" :boolop %d", n.Boolop))
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFuncCallExpr(sb *strings.Builder, n *FuncCallExpr) {
	sb.WriteString("{FUNCCALL")
	if n.FuncName != nil {
		sb.WriteString(" :funcname ")
		writeNode(sb, n.FuncName)
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
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	if n.KeepClause != nil {
		sb.WriteString(" :keepClause ")
		writeNode(sb, n.KeepClause)
	}
	if n.Over != nil {
		sb.WriteString(" :over ")
		writeNode(sb, n.Over)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCaseExpr(sb *strings.Builder, n *CaseExpr) {
	sb.WriteString("{CASE")
	if n.Arg != nil {
		sb.WriteString(" :arg ")
		writeNode(sb, n.Arg)
	}
	if n.Whens != nil {
		sb.WriteString(" :whens ")
		writeNode(sb, n.Whens)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCaseWhen(sb *strings.Builder, n *CaseWhen) {
	sb.WriteString("{CASEWHEN")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Result != nil {
		sb.WriteString(" :result ")
		writeNode(sb, n.Result)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDecodeExpr(sb *strings.Builder, n *DecodeExpr) {
	sb.WriteString("{DECODE")
	if n.Arg != nil {
		sb.WriteString(" :arg ")
		writeNode(sb, n.Arg)
	}
	if n.Pairs != nil {
		sb.WriteString(" :pairs ")
		writeNode(sb, n.Pairs)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDecodePair(sb *strings.Builder, n *DecodePair) {
	sb.WriteString("{DECODEPAIR")
	if n.Search != nil {
		sb.WriteString(" :search ")
		writeNode(sb, n.Search)
	}
	if n.Result != nil {
		sb.WriteString(" :result ")
		writeNode(sb, n.Result)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCastExpr(sb *strings.Builder, n *CastExpr) {
	sb.WriteString("{CAST")
	sb.WriteString(" :arg ")
	writeNode(sb, n.Arg)
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMultisetExpr(sb *strings.Builder, n *MultisetExpr) {
	sb.WriteString("{MULTISET")
	sb.WriteString(fmt.Sprintf(" :op %q", n.Op))
	sb.WriteString(" :left ")
	writeNode(sb, n.Left)
	sb.WriteString(" :right ")
	writeNode(sb, n.Right)
	if n.All {
		sb.WriteString(" :all true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCursorExpr(sb *strings.Builder, n *CursorExpr) {
	sb.WriteString("{CURSOR")
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTreatExpr(sb *strings.Builder, n *TreatExpr) {
	sb.WriteString("{TREAT")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIntervalExpr(sb *strings.Builder, n *IntervalExpr) {
	sb.WriteString("{INTERVAL")
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	if n.From != "" {
		sb.WriteString(fmt.Sprintf(" :from %q", n.From))
	}
	if n.To != "" {
		sb.WriteString(fmt.Sprintf(" :to %q", n.To))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeBetweenExpr(sb *strings.Builder, n *BetweenExpr) {
	sb.WriteString("{BETWEEN")
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(" :low ")
	writeNode(sb, n.Low)
	sb.WriteString(" :high ")
	writeNode(sb, n.High)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeInExpr(sb *strings.Builder, n *InExpr) {
	sb.WriteString("{IN")
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	if n.List != nil {
		sb.WriteString(" :list ")
		writeNode(sb, n.List)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLikeExpr(sb *strings.Builder, n *LikeExpr) {
	sb.WriteString("{LIKE")
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(" :pattern ")
	writeNode(sb, n.Pattern)
	if n.Escape != nil {
		sb.WriteString(" :escape ")
		writeNode(sb, n.Escape)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIsExpr(sb *strings.Builder, n *IsExpr) {
	sb.WriteString("{IS")
	if n.Not {
		sb.WriteString(" :not true")
	}
	sb.WriteString(fmt.Sprintf(" :test %q", n.Test))
	sb.WriteString(" :expr ")
	writeNode(sb, n.Expr)
	if n.TypeList != nil {
		sb.WriteString(" :typeList ")
		writeNode(sb, n.TypeList)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExistsExpr(sb *strings.Builder, n *ExistsExpr) {
	sb.WriteString("{EXISTS :subquery ")
	writeNode(sb, n.Subquery)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSubqueryExpr(sb *strings.Builder, n *SubqueryExpr) {
	sb.WriteString("{SUBQUERY :query ")
	writeNode(sb, n.Subquery)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeParenExpr(sb *strings.Builder, n *ParenExpr) {
	sb.WriteString("{PAREN :expr ")
	writeNode(sb, n.Expr)
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Window nodes
// ---------------------------------------------------------------------------

func writeWindowSpec(sb *strings.Builder, n *WindowSpec) {
	sb.WriteString("{WINDOWSPEC")
	if n.WindowName != "" {
		sb.WriteString(fmt.Sprintf(" :windowName %q", n.WindowName))
	}
	if n.PartitionBy != nil {
		sb.WriteString(" :partitionBy ")
		writeNode(sb, n.PartitionBy)
	}
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	if n.Frame != nil {
		sb.WriteString(" :frame ")
		writeNode(sb, n.Frame)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
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
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWindowBound(sb *strings.Builder, n *WindowBound) {
	sb.WriteString("{WINDOWBOUND")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeKeepClause(sb *strings.Builder, n *KeepClause) {
	sb.WriteString("{KEEP")
	if n.IsFirst {
		sb.WriteString(" :isFirst true")
	}
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Grouping extension nodes
// ---------------------------------------------------------------------------

func writeGroupingSetsClause(sb *strings.Builder, n *GroupingSetsClause) {
	sb.WriteString("{GROUPING_SETS")
	if n.Sets != nil {
		sb.WriteString(" :sets ")
		writeNode(sb, n.Sets)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCubeClause(sb *strings.Builder, n *CubeClause) {
	sb.WriteString("{CUBE")
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRollupClause(sb *strings.Builder, n *RollupClause) {
	sb.WriteString("{ROLLUP")
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Table nodes
// ---------------------------------------------------------------------------

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
	if n.On != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.On)
	}
	if n.Using != nil {
		sb.WriteString(" :using ")
		writeNode(sb, n.Using)
	}
	if n.OracleJoin {
		sb.WriteString(" :oracleJoin true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTableRef(sb *strings.Builder, n *TableRef) {
	sb.WriteString("{TABLEREF")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	if n.Sample != nil {
		sb.WriteString(" :sample ")
		writeNode(sb, n.Sample)
	}
	if n.Flashback != nil {
		sb.WriteString(" :flashback ")
		writeNode(sb, n.Flashback)
	}
	if n.PartitionExt != nil {
		sb.WriteString(" :partitionExt ")
		writeNode(sb, n.PartitionExt)
	}
	if n.Dblink != "" {
		sb.WriteString(fmt.Sprintf(" :dblink %q", n.Dblink))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSubqueryRef(sb *strings.Builder, n *SubqueryRef) {
	sb.WriteString("{SUBQUERYREF")
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLateralRef(sb *strings.Builder, n *LateralRef) {
	sb.WriteString("{LATERAL")
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeXmlTableRef(sb *strings.Builder, n *XmlTableRef) {
	sb.WriteString("{XMLTABLE")
	if n.XPath != nil {
		sb.WriteString(" :xpath ")
		writeNode(sb, n.XPath)
	}
	if n.Passing != nil {
		sb.WriteString(" :passing ")
		writeNode(sb, n.Passing)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeXmlTableColumn(sb *strings.Builder, n *XmlTableColumn) {
	sb.WriteString("{XMLTABLECOL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.ForOrdinality {
		sb.WriteString(" :forOrdinality true")
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	if n.Path != nil {
		sb.WriteString(" :path ")
		writeNode(sb, n.Path)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeJsonTableRef(sb *strings.Builder, n *JsonTableRef) {
	sb.WriteString("{JSON_TABLE")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.Path != nil {
		sb.WriteString(" :path ")
		writeNode(sb, n.Path)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeJsonTableColumn(sb *strings.Builder, n *JsonTableColumn) {
	sb.WriteString("{JSON_TABLE_COL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.ForOrdinality {
		sb.WriteString(" :forOrdinality true")
	}
	if n.Exists {
		sb.WriteString(" :exists true")
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	if n.Path != nil {
		sb.WriteString(" :path ")
		writeNode(sb, n.Path)
	}
	if n.Nested != nil {
		sb.WriteString(" :nested ")
		writeNode(sb, n.Nested)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSampleClause(sb *strings.Builder, n *SampleClause) {
	sb.WriteString("{SAMPLE")
	if n.Percent != nil {
		sb.WriteString(" :percent ")
		writeNode(sb, n.Percent)
	}
	if n.Seed != nil {
		sb.WriteString(" :seed ")
		writeNode(sb, n.Seed)
	}
	if n.Block {
		sb.WriteString(" :block true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePivotClause(sb *strings.Builder, n *PivotClause) {
	sb.WriteString("{PIVOT")
	if n.XML {
		sb.WriteString(" :xml true")
	}
	if n.AggFuncs != nil {
		sb.WriteString(" :aggFuncs ")
		writeNode(sb, n.AggFuncs)
	}
	if n.ForCol != nil {
		sb.WriteString(" :forCol ")
		writeNode(sb, n.ForCol)
	}
	if n.ForCols != nil {
		sb.WriteString(" :forCols ")
		writeNode(sb, n.ForCols)
	}
	if n.InList != nil {
		sb.WriteString(" :inList ")
		writeNode(sb, n.InList)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUnpivotClause(sb *strings.Builder, n *UnpivotClause) {
	sb.WriteString("{UNPIVOT")
	if n.ValueCol != nil {
		sb.WriteString(" :valueCol ")
		writeNode(sb, n.ValueCol)
	}
	if n.PivotCol != nil {
		sb.WriteString(" :pivotCol ")
		writeNode(sb, n.PivotCol)
	}
	if n.InList != nil {
		sb.WriteString(" :inList ")
		writeNode(sb, n.InList)
	}
	if n.IncludeNulls {
		sb.WriteString(" :includeNulls true")
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Clause nodes
// ---------------------------------------------------------------------------

func writeHierarchicalClause(sb *strings.Builder, n *HierarchicalClause) {
	sb.WriteString("{HIERARCHICAL")
	if n.ConnectBy != nil {
		sb.WriteString(" :connectBy ")
		writeNode(sb, n.ConnectBy)
	}
	if n.StartWith != nil {
		sb.WriteString(" :startWith ")
		writeNode(sb, n.StartWith)
	}
	if n.IsNocycle {
		sb.WriteString(" :nocycle true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWithClause(sb *strings.Builder, n *WithClause) {
	sb.WriteString("{WITH")
	if n.Recursive {
		sb.WriteString(" :recursive true")
	}
	if n.CTEs != nil {
		sb.WriteString(" :ctes ")
		writeNode(sb, n.CTEs)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCTE(sb *strings.Builder, n *CTE) {
	sb.WriteString("{CTE")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	if n.Search != nil {
		sb.WriteString(" :search ")
		writeNode(sb, n.Search)
	}
	if n.Cycle != nil {
		sb.WriteString(" :cycle ")
		writeNode(sb, n.Cycle)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeForUpdateClause(sb *strings.Builder, n *ForUpdateClause) {
	sb.WriteString("{FORUPDATE")
	if n.Tables != nil {
		sb.WriteString(" :tables ")
		writeNode(sb, n.Tables)
	}
	if n.NoWait {
		sb.WriteString(" :nowait true")
	}
	if n.Wait != nil {
		sb.WriteString(" :wait ")
		writeNode(sb, n.Wait)
	}
	if n.SkipLocked {
		sb.WriteString(" :skipLocked true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFetchFirstClause(sb *strings.Builder, n *FetchFirstClause) {
	sb.WriteString("{FETCHFIRST")
	if n.Offset != nil {
		sb.WriteString(" :offset ")
		writeNode(sb, n.Offset)
	}
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
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelClause(sb *strings.Builder, n *ModelClause) {
	sb.WriteString("{MODEL")
	if n.CellRefOptions != nil {
		sb.WriteString(" :cellRefOptions ")
		writeNode(sb, n.CellRefOptions)
	}
	if n.ReturnRows != "" {
		sb.WriteString(fmt.Sprintf(" :returnRows %q", n.ReturnRows))
	}
	for _, ref := range n.RefModels {
		sb.WriteString(" :refModel ")
		writeNode(sb, ref)
	}
	if n.MainModel != nil {
		sb.WriteString(" :mainModel ")
		writeNode(sb, n.MainModel)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelCellRefOptions(sb *strings.Builder, n *ModelCellRefOptions) {
	sb.WriteString("{MODEL_CELL_REF_OPTIONS")
	if n.IgnoreNav {
		sb.WriteString(" :ignoreNav true")
	}
	if n.KeepNav {
		sb.WriteString(" :keepNav true")
	}
	if n.UniqueDimension {
		sb.WriteString(" :uniqueDimension true")
	}
	if n.UniqueSingleRef {
		sb.WriteString(" :uniqueSingleRef true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelRefModel(sb *strings.Builder, n *ModelRefModel) {
	sb.WriteString("{MODEL_REF")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	if n.ColumnClauses != nil {
		sb.WriteString(" :columnClauses ")
		writeNode(sb, n.ColumnClauses)
	}
	if n.CellRefOptions != nil {
		sb.WriteString(" :cellRefOptions ")
		writeNode(sb, n.CellRefOptions)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelMainModel(sb *strings.Builder, n *ModelMainModel) {
	sb.WriteString("{MODEL_MAIN")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.ColumnClauses != nil {
		sb.WriteString(" :columnClauses ")
		writeNode(sb, n.ColumnClauses)
	}
	if n.CellRefOptions != nil {
		sb.WriteString(" :cellRefOptions ")
		writeNode(sb, n.CellRefOptions)
	}
	if n.RulesClause != nil {
		sb.WriteString(" :rulesClause ")
		writeNode(sb, n.RulesClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelColumnClauses(sb *strings.Builder, n *ModelColumnClauses) {
	sb.WriteString("{MODEL_COLUMNS")
	if n.PartitionBy != nil {
		sb.WriteString(" :partitionBy ")
		writeNode(sb, n.PartitionBy)
	}
	if n.DimensionBy != nil {
		sb.WriteString(" :dimensionBy ")
		writeNode(sb, n.DimensionBy)
	}
	if n.Measures != nil {
		sb.WriteString(" :measures ")
		writeNode(sb, n.Measures)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelRulesClause(sb *strings.Builder, n *ModelRulesClause) {
	sb.WriteString("{MODEL_RULES")
	if n.UpdateMode != "" {
		sb.WriteString(fmt.Sprintf(" :updateMode %q", n.UpdateMode))
	}
	if n.OrderMode != "" {
		sb.WriteString(fmt.Sprintf(" :orderMode %q", n.OrderMode))
	}
	if n.Iterate != nil {
		sb.WriteString(" :iterate ")
		writeNode(sb, n.Iterate)
	}
	if n.Until != nil {
		sb.WriteString(" :until ")
		writeNode(sb, n.Until)
	}
	if n.Rules != nil {
		sb.WriteString(" :rules ")
		writeNode(sb, n.Rules)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelRule(sb *strings.Builder, n *ModelRule) {
	sb.WriteString("{MODEL_RULE")
	if n.CellRef != nil {
		sb.WriteString(" :cellRef ")
		writeNode(sb, n.CellRef)
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeModelForLoop(sb *strings.Builder, n *ModelForLoop) {
	sb.WriteString("{MODEL_FOR_LOOP")
	if n.Column != "" {
		sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
	}
	if n.InList != nil {
		sb.WriteString(" :inList ")
		writeNode(sb, n.InList)
	}
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	if n.LikePattern != nil {
		sb.WriteString(" :likePattern ")
		writeNode(sb, n.LikePattern)
	}
	if n.FromExpr != nil {
		sb.WriteString(" :from ")
		writeNode(sb, n.FromExpr)
	}
	if n.ToExpr != nil {
		sb.WriteString(" :to ")
		writeNode(sb, n.ToExpr)
	}
	if n.Increment {
		sb.WriteString(" :increment true")
	}
	if n.IncrExpr != nil {
		sb.WriteString(" :incrExpr ")
		writeNode(sb, n.IncrExpr)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFlashbackClause(sb *strings.Builder, n *FlashbackClause) {
	sb.WriteString("{FLASHBACK")
	if n.Type != "" {
		sb.WriteString(fmt.Sprintf(" :type %q", n.Type))
	}
	if n.IsVersions {
		sb.WriteString(" :isVersions true")
	}
	if n.IsPeriodFor {
		sb.WriteString(" :isPeriodFor true")
	}
	if n.PeriodColumn != "" {
		sb.WriteString(fmt.Sprintf(" :periodColumn %q", n.PeriodColumn))
	}
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.VersionsLow != nil {
		sb.WriteString(" :versionsLow ")
		writeNode(sb, n.VersionsLow)
	}
	if n.VersionsHigh != nil {
		sb.WriteString(" :versionsHigh ")
		writeNode(sb, n.VersionsHigh)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSortBy(sb *strings.Builder, n *SortBy) {
	sb.WriteString("{SORTBY")
	if n.Expr != nil {
		sb.WriteString(" :node ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :dir %d", n.Dir))
	sb.WriteString(fmt.Sprintf(" :nulls %d", n.NullOrder))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeWindowDef(sb *strings.Builder, n *WindowDef) {
	sb.WriteString("{WINDOWDEF")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Spec != nil {
		sb.WriteString(" :spec ")
		writeNode(sb, n.Spec)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCTESearchClause(sb *strings.Builder, n *CTESearchClause) {
	sb.WriteString("{CTESEARCH")
	if n.BreadthFirst {
		sb.WriteString(" :breadthFirst true")
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.SetColumn != "" {
		sb.WriteString(fmt.Sprintf(" :setColumn %q", n.SetColumn))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCTECycleClause(sb *strings.Builder, n *CTECycleClause) {
	sb.WriteString("{CTECYCLE")
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.SetColumn != "" {
		sb.WriteString(fmt.Sprintf(" :setColumn %q", n.SetColumn))
	}
	if n.CycleValue != nil {
		sb.WriteString(" :cycleValue ")
		writeNode(sb, n.CycleValue)
	}
	if n.NoCycleValue != nil {
		sb.WriteString(" :noCycleValue ")
		writeNode(sb, n.NoCycleValue)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePartitionExtClause(sb *strings.Builder, n *PartitionExtClause) {
	sb.WriteString("{PARTITIONEXT")
	if n.IsSubpartition {
		sb.WriteString(" :isSubpartition true")
	}
	if n.IsFor {
		sb.WriteString(" :isFor true")
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Keys != nil {
		sb.WriteString(" :keys ")
		writeNode(sb, n.Keys)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTableCollectionExpr(sb *strings.Builder, n *TableCollectionExpr) {
	sb.WriteString("{TABLECOLLECTION")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	if n.OuterJoin {
		sb.WriteString(" :outerJoin true")
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMatchRecognizeClause(sb *strings.Builder, n *MatchRecognizeClause) {
	sb.WriteString("{MATCHRECOGNIZE")
	if n.PartitionBy != nil {
		sb.WriteString(" :partitionBy ")
		writeNode(sb, n.PartitionBy)
	}
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	if n.Measures != nil {
		sb.WriteString(" :measures ")
		writeNode(sb, n.Measures)
	}
	if n.RowsPerMatch != "" {
		sb.WriteString(fmt.Sprintf(" :rowsPerMatch %q", n.RowsPerMatch))
	}
	if n.AfterMatch != "" {
		sb.WriteString(fmt.Sprintf(" :afterMatch %q", n.AfterMatch))
	}
	if n.Pattern != "" {
		sb.WriteString(fmt.Sprintf(" :pattern %q", n.Pattern))
	}
	if n.Subsets != nil {
		sb.WriteString(" :subsets ")
		writeNode(sb, n.Subsets)
	}
	if n.Definitions != nil {
		sb.WriteString(" :definitions ")
		writeNode(sb, n.Definitions)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeContainersExpr(sb *strings.Builder, n *ContainersExpr) {
	sb.WriteString("{CONTAINERS")
	if n.IsShards {
		sb.WriteString(" :isShards true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeInlineExternalTable(sb *strings.Builder, n *InlineExternalTable) {
	sb.WriteString("{INLINEEXTERNALTABLE")
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Type != "" {
		sb.WriteString(fmt.Sprintf(" :type %q", n.Type))
	}
	if n.Directory != "" {
		sb.WriteString(fmt.Sprintf(" :directory %q", n.Directory))
	}
	if n.Location != "" {
		sb.WriteString(fmt.Sprintf(" :location %q", n.Location))
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeResTarget(sb *strings.Builder, n *ResTarget) {
	sb.WriteString("{RESTARGET")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Expr != nil {
		sb.WriteString(" :val ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// DML helper nodes
// ---------------------------------------------------------------------------

func writeInsertIntoClause(sb *strings.Builder, n *InsertIntoClause) {
	sb.WriteString("{INSERTINTO")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	if n.When != nil {
		sb.WriteString(" :when ")
		writeNode(sb, n.When)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeErrorLogClause(sb *strings.Builder, n *ErrorLogClause) {
	sb.WriteString("{ERRORLOG")
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.Tag != nil {
		sb.WriteString(" :tag ")
		writeNode(sb, n.Tag)
	}
	if n.Reject != nil {
		sb.WriteString(" :reject ")
		writeNode(sb, n.Reject)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetClause(sb *strings.Builder, n *SetClause) {
	sb.WriteString("{SETCLAUSE")
	if n.Column != nil {
		sb.WriteString(" :column ")
		writeNode(sb, n.Column)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeClause(sb *strings.Builder, n *MergeClause) {
	sb.WriteString("{MERGECLAUSE")
	if n.Matched {
		sb.WriteString(" :matched true")
	}
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.UpdateSet != nil {
		sb.WriteString(" :updateSet ")
		writeNode(sb, n.UpdateSet)
	}
	if n.UpdateWhere != nil {
		sb.WriteString(" :updateWhere ")
		writeNode(sb, n.UpdateWhere)
	}
	if n.DeleteWhere != nil {
		sb.WriteString(" :deleteWhere ")
		writeNode(sb, n.DeleteWhere)
	}
	if n.InsertCols != nil {
		sb.WriteString(" :insertCols ")
		writeNode(sb, n.InsertCols)
	}
	if n.InsertVals != nil {
		sb.WriteString(" :insertVals ")
		writeNode(sb, n.InsertVals)
	}
	if n.InsertWhere != nil {
		sb.WriteString(" :insertWhere ")
		writeNode(sb, n.InsertWhere)
	}
	if n.IsDelete {
		sb.WriteString(" :isDelete true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// DDL nodes
// ---------------------------------------------------------------------------

func writeColumnDef(sb *strings.Builder, n *ColumnDef) {
	sb.WriteString("{COLUMNDEF")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	if n.Identity != nil {
		sb.WriteString(" :identity ")
		writeNode(sb, n.Identity)
	}
	if n.Virtual != nil {
		sb.WriteString(" :virtual ")
		writeNode(sb, n.Virtual)
	}
	if n.Invisible {
		sb.WriteString(" :invisible true")
	}
	if n.NotNull {
		sb.WriteString(" :notNull true")
	}
	if n.Null {
		sb.WriteString(" :null true")
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	if n.Collation != "" {
		sb.WriteString(fmt.Sprintf(" :collation %q", n.Collation))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeColumnConstraint(sb *strings.Builder, n *ColumnConstraint) {
	sb.WriteString("{COLCONSTRAINT")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
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
	if n.OnDelete != "" {
		sb.WriteString(fmt.Sprintf(" :onDelete %q", n.OnDelete))
	}
	if n.Deferrable {
		sb.WriteString(" :deferrable true")
	}
	if n.Initially != "" {
		sb.WriteString(fmt.Sprintf(" :initially %q", n.Initially))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTableConstraint(sb *strings.Builder, n *TableConstraint) {
	sb.WriteString("{TBLCONSTRAINT")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
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
	if n.OnDelete != "" {
		sb.WriteString(fmt.Sprintf(" :onDelete %q", n.OnDelete))
	}
	if n.Deferrable {
		sb.WriteString(" :deferrable true")
	}
	if n.Initially != "" {
		sb.WriteString(fmt.Sprintf(" :initially %q", n.Initially))
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIdentityClause(sb *strings.Builder, n *IdentityClause) {
	sb.WriteString("{IDENTITY")
	if n.Always {
		sb.WriteString(" :always true")
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeStorageClause(sb *strings.Builder, n *StorageClause) {
	sb.WriteString("{STORAGE")
	if n.Initial != "" {
		sb.WriteString(fmt.Sprintf(" :initial %q", n.Initial))
	}
	if n.Next != "" {
		sb.WriteString(fmt.Sprintf(" :next %q", n.Next))
	}
	if n.PctIncrease != "" {
		sb.WriteString(fmt.Sprintf(" :pctIncrease %q", n.PctIncrease))
	}
	if n.MinExtents != "" {
		sb.WriteString(fmt.Sprintf(" :minExtents %q", n.MinExtents))
	}
	if n.MaxExtents != "" {
		sb.WriteString(fmt.Sprintf(" :maxExtents %q", n.MaxExtents))
	}
	if n.PctFree != "" {
		sb.WriteString(fmt.Sprintf(" :pctFree %q", n.PctFree))
	}
	if n.PctUsed != "" {
		sb.WriteString(fmt.Sprintf(" :pctUsed %q", n.PctUsed))
	}
	if n.Logging {
		sb.WriteString(" :logging true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePartitionClause(sb *strings.Builder, n *PartitionClause) {
	sb.WriteString("{PARTITION")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Partitions != nil {
		sb.WriteString(" :partitions ")
		writeNode(sb, n.Partitions)
	}
	if n.Subpartition != nil {
		sb.WriteString(" :subpartition ")
		writeNode(sb, n.Subpartition)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePartitionDef(sb *strings.Builder, n *PartitionDef) {
	sb.WriteString("{PARTITIONDEF")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTableCmd(sb *strings.Builder, n *AlterTableCmd) {
	sb.WriteString("{ALTERTABLECMD")
	sb.WriteString(fmt.Sprintf(" :action %d", n.Action))
	if n.ColumnDef != nil {
		sb.WriteString(" :columnDef ")
		writeNode(sb, n.ColumnDef)
	}
	if n.ColumnName != "" {
		sb.WriteString(fmt.Sprintf(" :columnName %q", n.ColumnName))
	}
	if n.NewName != "" {
		sb.WriteString(fmt.Sprintf(" :newName %q", n.NewName))
	}
	if n.Constraint != nil {
		sb.WriteString(" :constraint ")
		writeNode(sb, n.Constraint)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIndexColumn(sb *strings.Builder, n *IndexColumn) {
	sb.WriteString("{INDEXCOL")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :dir %d", n.Dir))
	sb.WriteString(fmt.Sprintf(" :nulls %d", n.NullOrder))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// Statement nodes
// ---------------------------------------------------------------------------

func writeRawStmt(sb *strings.Builder, n *RawStmt) {
	sb.WriteString("{RAWSTMT")
	if n.Stmt != nil {
		sb.WriteString(" :stmt ")
		writeNode(sb, n.Stmt)
	}
	sb.WriteString(fmt.Sprintf(" :stmt_location %d", n.StmtLocation))
	sb.WriteString(fmt.Sprintf(" :stmt_len %d", n.StmtLen))
	sb.WriteString("}")
}

func writeSelectStmt(sb *strings.Builder, n *SelectStmt) {
	sb.WriteString("{SELECT")
	if n.WithClause != nil {
		sb.WriteString(" :withClause ")
		writeNode(sb, n.WithClause)
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.UniqueKw {
		sb.WriteString(" :unique true")
	}
	if n.All {
		sb.WriteString(" :all true")
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	if n.TargetList != nil {
		sb.WriteString(" :targetList ")
		writeNode(sb, n.TargetList)
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.FromClause != nil {
		sb.WriteString(" :fromClause ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.Hierarchical != nil {
		sb.WriteString(" :hierarchical ")
		writeNode(sb, n.Hierarchical)
	}
	if n.GroupClause != nil {
		sb.WriteString(" :groupClause ")
		writeNode(sb, n.GroupClause)
	}
	if n.HavingClause != nil {
		sb.WriteString(" :havingClause ")
		writeNode(sb, n.HavingClause)
	}
	if n.ModelClause != nil {
		sb.WriteString(" :modelClause ")
		writeNode(sb, n.ModelClause)
	}
	if len(n.WindowDefs) > 0 {
		sb.WriteString(" :windowDefs (")
		for i, wd := range n.WindowDefs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, wd)
		}
		sb.WriteString(")")
	}
	if n.QualifyClause != nil {
		sb.WriteString(" :qualifyClause ")
		writeNode(sb, n.QualifyClause)
	}
	if n.SiblingsOrder {
		sb.WriteString(" :siblingsOrder true")
	}
	if n.OrderBy != nil {
		sb.WriteString(" :orderBy ")
		writeNode(sb, n.OrderBy)
	}
	if n.ForUpdate != nil {
		sb.WriteString(" :forUpdate ")
		writeNode(sb, n.ForUpdate)
	}
	if n.FetchFirst != nil {
		sb.WriteString(" :fetchFirst ")
		writeNode(sb, n.FetchFirst)
	}
	if n.Pivot != nil {
		sb.WriteString(" :pivot ")
		writeNode(sb, n.Pivot)
	}
	if n.Unpivot != nil {
		sb.WriteString(" :unpivot ")
		writeNode(sb, n.Unpivot)
	}
	if n.Op != SETOP_NONE {
		sb.WriteString(fmt.Sprintf(" :op %d", n.Op))
	}
	if n.SetAll {
		sb.WriteString(" :setAll true")
	}
	if n.Larg != nil {
		sb.WriteString(" :larg ")
		writeNode(sb, n.Larg)
	}
	if n.Rarg != nil {
		sb.WriteString(" :rarg ")
		writeNode(sb, n.Rarg)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeInsertStmt(sb *strings.Builder, n *InsertStmt) {
	sb.WriteString("{INSERT")
	sb.WriteString(fmt.Sprintf(" :insertType %d", n.InsertType))
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.PartitionExt != nil {
		sb.WriteString(" :partitionExt ")
		writeNode(sb, n.PartitionExt)
	}
	if n.Dblink != "" {
		sb.WriteString(fmt.Sprintf(" :dblink %q", n.Dblink))
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Values != nil {
		sb.WriteString(" :values ")
		writeNode(sb, n.Values)
	}
	if n.SetClauses != nil {
		sb.WriteString(" :setClauses ")
		writeNode(sb, n.SetClauses)
	}
	if n.ByName {
		sb.WriteString(" :byName true")
	}
	if n.ByPosition {
		sb.WriteString(" :byPosition true")
	}
	if n.Select != nil {
		sb.WriteString(" :select ")
		writeNode(sb, n.Select)
	}
	if n.MultiTable != nil {
		sb.WriteString(" :multiTable ")
		writeNode(sb, n.MultiTable)
	}
	if n.Subquery != nil {
		sb.WriteString(" :subquery ")
		writeNode(sb, n.Subquery)
	}
	if n.Returning != nil {
		sb.WriteString(" :returning ")
		writeNode(sb, n.Returning)
	}
	if n.ErrorLog != nil {
		sb.WriteString(" :errorLog ")
		writeNode(sb, n.ErrorLog)
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeUpdateStmt(sb *strings.Builder, n *UpdateStmt) {
	sb.WriteString("{UPDATE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.PartitionExt != nil {
		sb.WriteString(" :partitionExt ")
		writeNode(sb, n.PartitionExt)
	}
	if n.Dblink != "" {
		sb.WriteString(fmt.Sprintf(" :dblink %q", n.Dblink))
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	if n.SetClauses != nil {
		sb.WriteString(" :setClauses ")
		writeNode(sb, n.SetClauses)
	}
	if n.FromClause != nil {
		sb.WriteString(" :fromClause ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.Returning != nil {
		sb.WriteString(" :returning ")
		writeNode(sb, n.Returning)
	}
	if n.ErrorLog != nil {
		sb.WriteString(" :errorLog ")
		writeNode(sb, n.ErrorLog)
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDeleteStmt(sb *strings.Builder, n *DeleteStmt) {
	sb.WriteString("{DELETE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.PartitionExt != nil {
		sb.WriteString(" :partitionExt ")
		writeNode(sb, n.PartitionExt)
	}
	if n.Dblink != "" {
		sb.WriteString(fmt.Sprintf(" :dblink %q", n.Dblink))
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.Returning != nil {
		sb.WriteString(" :returning ")
		writeNode(sb, n.Returning)
	}
	if n.ErrorLog != nil {
		sb.WriteString(" :errorLog ")
		writeNode(sb, n.ErrorLog)
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeMergeStmt(sb *strings.Builder, n *MergeStmt) {
	sb.WriteString("{MERGE")
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	if n.TargetAlias != nil {
		sb.WriteString(" :targetAlias ")
		writeNode(sb, n.TargetAlias)
	}
	if n.Source != nil {
		sb.WriteString(" :source ")
		writeNode(sb, n.Source)
	}
	if n.SourceAlias != nil {
		sb.WriteString(" :sourceAlias ")
		writeNode(sb, n.SourceAlias)
	}
	if n.On != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.On)
	}
	if n.Clauses != nil {
		sb.WriteString(" :clauses ")
		writeNode(sb, n.Clauses)
	}
	if n.ErrorLog != nil {
		sb.WriteString(" :errorLog ")
		writeNode(sb, n.ErrorLog)
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTableStmt(sb *strings.Builder, n *CreateTableStmt) {
	sb.WriteString("{CREATETABLE")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Global {
		sb.WriteString(" :global true")
	}
	if n.Private {
		sb.WriteString(" :private true")
	}
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
	if n.AsQuery != nil {
		sb.WriteString(" :asQuery ")
		writeNode(sb, n.AsQuery)
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	if n.Storage != nil {
		sb.WriteString(" :storage ")
		writeNode(sb, n.Storage)
	}
	if n.Partition != nil {
		sb.WriteString(" :partition ")
		writeNode(sb, n.Partition)
	}
	if n.OnCommit != "" {
		sb.WriteString(fmt.Sprintf(" :onCommit %q", n.OnCommit))
	}
	if n.Parallel != "" {
		sb.WriteString(fmt.Sprintf(" :parallel %q", n.Parallel))
	}
	if n.Compress != "" {
		sb.WriteString(fmt.Sprintf(" :compress %q", n.Compress))
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Hints != nil {
		sb.WriteString(" :hints ")
		writeNode(sb, n.Hints)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
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
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDropStmt(sb *strings.Builder, n *DropStmt) {
	sb.WriteString("{DROP")
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Names != nil {
		sb.WriteString(" :names ")
		writeNode(sb, n.Names)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.Purge {
		sb.WriteString(" :purge true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateIndexStmt(sb *strings.Builder, n *CreateIndexStmt) {
	sb.WriteString("{CREATEINDEX")
	if n.Unique {
		sb.WriteString(" :unique true")
	}
	if n.Bitmap {
		sb.WriteString(" :bitmap true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.FunctionBased {
		sb.WriteString(" :functionBased true")
	}
	if n.Reverse {
		sb.WriteString(" :reverse true")
	}
	if n.Local {
		sb.WriteString(" :local true")
	}
	if n.Global {
		sb.WriteString(" :global true")
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	if n.Parallel != "" {
		sb.WriteString(fmt.Sprintf(" :parallel %q", n.Parallel))
	}
	if n.Compress != "" {
		sb.WriteString(fmt.Sprintf(" :compress %q", n.Compress))
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateViewStmt(sb *strings.Builder, n *CreateViewStmt) {
	sb.WriteString("{CREATEVIEW")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.NoForce {
		sb.WriteString(" :noForce true")
	}
	if n.Materialized {
		sb.WriteString(" :materialized true")
	}
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
	if n.WithCheckOpt {
		sb.WriteString(" :withCheckOpt true")
	}
	if n.WithReadOnly {
		sb.WriteString(" :withReadOnly true")
	}
	if n.BuildMode != "" {
		sb.WriteString(fmt.Sprintf(" :buildMode %q", n.BuildMode))
	}
	if n.RefreshMode != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMode %q", n.RefreshMode))
	}
	if n.EnableQuery {
		sb.WriteString(" :enableQuery true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateSequenceStmt(sb *strings.Builder, n *CreateSequenceStmt) {
	sb.WriteString("{CREATESEQUENCE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IncrementBy != nil {
		sb.WriteString(" :incrementBy ")
		writeNode(sb, n.IncrementBy)
	}
	if n.StartWith != nil {
		sb.WriteString(" :startWith ")
		writeNode(sb, n.StartWith)
	}
	if n.MaxValue != nil {
		sb.WriteString(" :maxValue ")
		writeNode(sb, n.MaxValue)
	}
	if n.MinValue != nil {
		sb.WriteString(" :minValue ")
		writeNode(sb, n.MinValue)
	}
	if n.NoMaxValue {
		sb.WriteString(" :noMaxValue true")
	}
	if n.NoMinValue {
		sb.WriteString(" :noMinValue true")
	}
	if n.Cycle {
		sb.WriteString(" :cycle true")
	}
	if n.NoCycle {
		sb.WriteString(" :noCycle true")
	}
	if n.Cache != nil {
		sb.WriteString(" :cache ")
		writeNode(sb, n.Cache)
	}
	if n.NoCache {
		sb.WriteString(" :noCache true")
	}
	if n.Order {
		sb.WriteString(" :order true")
	}
	if n.NoOrder {
		sb.WriteString(" :noOrder true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateSynonymStmt(sb *strings.Builder, n *CreateSynonymStmt) {
	sb.WriteString("{CREATESYNONYM")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateDatabaseLinkStmt(sb *strings.Builder, n *CreateDatabaseLinkStmt) {
	sb.WriteString("{CREATEDBLINK")
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.ConnectTo != "" {
		sb.WriteString(fmt.Sprintf(" :connectTo %q", n.ConnectTo))
	}
	if n.Identified != "" {
		sb.WriteString(fmt.Sprintf(" :identified %q", n.Identified))
	}
	if n.Using != "" {
		sb.WriteString(fmt.Sprintf(" :using %q", n.Using))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTypeStmt(sb *strings.Builder, n *CreateTypeStmt) {
	sb.WriteString("{CREATETYPE")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Attributes != nil {
		sb.WriteString(" :attributes ")
		writeNode(sb, n.Attributes)
	}
	if n.AsTable != nil {
		sb.WriteString(" :asTable ")
		writeNode(sb, n.AsTable)
	}
	if n.AsVarray != nil {
		sb.WriteString(" :asVarray ")
		writeNode(sb, n.AsVarray)
	}
	if n.VarraySize != nil {
		sb.WriteString(" :varraySize ")
		writeNode(sb, n.VarraySize)
	}
	if n.IsBody {
		sb.WriteString(" :isBody true")
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTypeBodyMember(sb *strings.Builder, n *TypeBodyMember) {
	sb.WriteString("{TYPEBODYMEMBER")
	kindStr := "MEMBER"
	switch n.Kind {
	case TYPE_BODY_STATIC:
		kindStr = "STATIC"
	case TYPE_BODY_MAP:
		kindStr = "MAP"
	case TYPE_BODY_ORDER:
		kindStr = "ORDER"
	case TYPE_BODY_CONSTRUCTOR:
		kindStr = "CONSTRUCTOR"
	}
	sb.WriteString(fmt.Sprintf(" :kind %q", kindStr))
	if n.Subprog != nil {
		sb.WriteString(" :subprog ")
		writeNode(sb, n.Subprog)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreatePackageStmt(sb *strings.Builder, n *CreatePackageStmt) {
	sb.WriteString("{CREATEPACKAGE")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IsBody {
		sb.WriteString(" :isBody true")
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateProcedureStmt(sb *strings.Builder, n *CreateProcedureStmt) {
	sb.WriteString("{CREATEPROCEDURE")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Parameters != nil {
		sb.WriteString(" :parameters ")
		writeNode(sb, n.Parameters)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateFunctionStmt(sb *strings.Builder, n *CreateFunctionStmt) {
	sb.WriteString("{CREATEFUNCTION")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Parameters != nil {
		sb.WriteString(" :parameters ")
		writeNode(sb, n.Parameters)
	}
	if n.ReturnType != nil {
		sb.WriteString(" :returnType ")
		writeNode(sb, n.ReturnType)
	}
	if n.Deterministic {
		sb.WriteString(" :deterministic true")
	}
	if n.Pipelined {
		sb.WriteString(" :pipelined true")
	}
	if n.Parallel {
		sb.WriteString(" :parallel true")
	}
	if n.ResultCache {
		sb.WriteString(" :resultCache true")
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeParameter(sb *strings.Builder, n *Parameter) {
	sb.WriteString("{PARAM")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Mode != "" {
		sb.WriteString(fmt.Sprintf(" :mode %q", n.Mode))
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTriggerStmt(sb *strings.Builder, n *CreateTriggerStmt) {
	sb.WriteString("{CREATETRIGGER")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString(fmt.Sprintf(" :timing %d", n.Timing))
	if n.Events != nil {
		sb.WriteString(" :events ")
		writeNode(sb, n.Events)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.ForEachRow {
		sb.WriteString(" :forEachRow true")
	}
	if n.When != nil {
		sb.WriteString(" :when ")
		writeNode(sb, n.When)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	if n.Compound {
		sb.WriteString(" :compound true")
	}
	if n.Enable {
		sb.WriteString(" :enable true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTruncateStmt(sb *strings.Builder, n *TruncateStmt) {
	sb.WriteString("{TRUNCATE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Cluster {
		sb.WriteString(" :cluster true")
	}
	if n.PurgeMVLog {
		sb.WriteString(" :purgeMVLog true")
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeGrantStmt(sb *strings.Builder, n *GrantStmt) {
	sb.WriteString("{GRANT")
	if n.AllPriv {
		sb.WriteString(" :allPriv true")
	}
	if n.Privileges != nil {
		sb.WriteString(" :privileges ")
		writeNode(sb, n.Privileges)
	}
	if n.OnObject != nil {
		sb.WriteString(" :onObject ")
		writeNode(sb, n.OnObject)
	}
	sb.WriteString(fmt.Sprintf(" :onType %d", n.OnType))
	if n.Grantees != nil {
		sb.WriteString(" :grantees ")
		writeNode(sb, n.Grantees)
	}
	if n.WithAdmin {
		sb.WriteString(" :withAdmin true")
	}
	if n.WithGrant {
		sb.WriteString(" :withGrant true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRevokeStmt(sb *strings.Builder, n *RevokeStmt) {
	sb.WriteString("{REVOKE")
	if n.AllPriv {
		sb.WriteString(" :allPriv true")
	}
	if n.Privileges != nil {
		sb.WriteString(" :privileges ")
		writeNode(sb, n.Privileges)
	}
	if n.OnObject != nil {
		sb.WriteString(" :onObject ")
		writeNode(sb, n.OnObject)
	}
	sb.WriteString(fmt.Sprintf(" :onType %d", n.OnType))
	if n.Grantees != nil {
		sb.WriteString(" :grantees ")
		writeNode(sb, n.Grantees)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCommentStmt(sb *strings.Builder, n *CommentStmt) {
	sb.WriteString("{COMMENT")
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Object != nil {
		sb.WriteString(" :object ")
		writeNode(sb, n.Object)
	}
	if n.Column != "" {
		sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
	}
	sb.WriteString(fmt.Sprintf(" :comment %q", n.Comment))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterSessionStmt(sb *strings.Builder, n *AlterSessionStmt) {
	sb.WriteString("{ALTERSESSION")
	if n.SetParams != nil {
		sb.WriteString(" :setParams ")
		writeNode(sb, n.SetParams)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterSystemStmt(sb *strings.Builder, n *AlterSystemStmt) {
	sb.WriteString("{ALTERSYSTEM")
	if n.SetParams != nil {
		sb.WriteString(" :setParams ")
		writeNode(sb, n.SetParams)
	}
	if n.Kill != "" {
		sb.WriteString(fmt.Sprintf(" :kill %q", n.Kill))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetParam(sb *strings.Builder, n *SetParam) {
	sb.WriteString("{SETPARAM")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCommitStmt(sb *strings.Builder, n *CommitStmt) {
	sb.WriteString("{COMMIT")
	if n.Work {
		sb.WriteString(" :work true")
	}
	if n.Comment != "" {
		sb.WriteString(fmt.Sprintf(" :comment %q", n.Comment))
	}
	if n.Force != "" {
		sb.WriteString(fmt.Sprintf(" :force %q", n.Force))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRollbackStmt(sb *strings.Builder, n *RollbackStmt) {
	sb.WriteString("{ROLLBACK")
	if n.Work {
		sb.WriteString(" :work true")
	}
	if n.ToSavepoint != "" {
		sb.WriteString(fmt.Sprintf(" :toSavepoint %q", n.ToSavepoint))
	}
	if n.Force != "" {
		sb.WriteString(fmt.Sprintf(" :force %q", n.Force))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSavepointStmt(sb *strings.Builder, n *SavepointStmt) {
	sb.WriteString("{SAVEPOINT")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetTransactionStmt(sb *strings.Builder, n *SetTransactionStmt) {
	sb.WriteString("{SETTRANSACTION")
	if n.ReadOnly {
		sb.WriteString(" :readOnly true")
	}
	if n.ReadWrite {
		sb.WriteString(" :readWrite true")
	}
	if n.IsolLevel != "" {
		sb.WriteString(fmt.Sprintf(" :isolLevel %q", n.IsolLevel))
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAnalyzeStmt(sb *strings.Builder, n *AnalyzeStmt) {
	sb.WriteString("{ANALYZE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExplainPlanStmt(sb *strings.Builder, n *ExplainPlanStmt) {
	sb.WriteString("{EXPLAINPLAN")
	if n.StatementID != "" {
		sb.WriteString(fmt.Sprintf(" :statementID %q", n.StatementID))
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.Statement != nil {
		sb.WriteString(" :statement ")
		writeNode(sb, n.Statement)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFlashbackTableStmt(sb *strings.Builder, n *FlashbackTableStmt) {
	sb.WriteString("{FLASHBACKTABLE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.ToSCN != nil {
		sb.WriteString(" :toSCN ")
		writeNode(sb, n.ToSCN)
	}
	if n.ToTimestamp != nil {
		sb.WriteString(" :toTimestamp ")
		writeNode(sb, n.ToTimestamp)
	}
	if n.ToBeforeDrop {
		sb.WriteString(" :toBeforeDrop true")
	}
	if n.Rename != "" {
		sb.WriteString(fmt.Sprintf(" :rename %q", n.Rename))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFlashbackDatabaseStmt(sb *strings.Builder, n *FlashbackDatabaseStmt) {
	sb.WriteString("{FLASHBACKDATABASE")
	if n.ToSCN != nil {
		sb.WriteString(" :toSCN ")
		writeNode(sb, n.ToSCN)
	}
	if n.ToTimestamp != nil {
		sb.WriteString(" :toTimestamp ")
		writeNode(sb, n.ToTimestamp)
	}
	if n.ToRestorePoint != "" {
		sb.WriteString(fmt.Sprintf(" :toRestorePoint %q", n.ToRestorePoint))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePurgeStmt(sb *strings.Builder, n *PurgeStmt) {
	sb.WriteString("{PURGE")
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLockTableStmt(sb *strings.Builder, n *LockTableStmt) {
	sb.WriteString("{LOCK_TABLE")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.LockMode != "" {
		sb.WriteString(fmt.Sprintf(" :lockMode %q", n.LockMode))
	}
	if n.Nowait {
		sb.WriteString(" :nowait true")
	}
	if n.Wait != nil {
		sb.WriteString(" :wait ")
		writeNode(sb, n.Wait)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCallStmt(sb *strings.Builder, n *CallStmt) {
	sb.WriteString("{CALL")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeRenameStmt(sb *strings.Builder, n *RenameStmt) {
	sb.WriteString("{RENAME")
	if n.OldName != nil {
		sb.WriteString(" :oldName ")
		writeNode(sb, n.OldName)
	}
	if n.NewName != nil {
		sb.WriteString(" :newName ")
		writeNode(sb, n.NewName)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIdentifiedClause(sb *strings.Builder, n *IdentifiedClause) {
	sb.WriteString("{IDENTIFIED")
	switch n.Type {
	case IDENTIFIED_BY:
		sb.WriteString(" :type BY")
		if n.Password != "" {
			sb.WriteString(fmt.Sprintf(" :password %q", n.Password))
		}
		if n.OldPass != "" {
			sb.WriteString(fmt.Sprintf(" :replace %q", n.OldPass))
		}
	case IDENTIFIED_EXTERNALLY:
		sb.WriteString(" :type EXTERNALLY")
		if n.ExternalAs != "" {
			sb.WriteString(fmt.Sprintf(" :as %q", n.ExternalAs))
		}
	case IDENTIFIED_GLOBALLY:
		sb.WriteString(" :type GLOBALLY")
		if n.ExternalAs != "" {
			sb.WriteString(fmt.Sprintf(" :as %q", n.ExternalAs))
		}
	case IDENTIFIED_NO_AUTH:
		sb.WriteString(" :type NO_AUTHENTICATION")
	}
	sb.WriteString("}")
}

func writeUserQuotaClause(sb *strings.Builder, n *UserQuotaClause) {
	sb.WriteString("{QUOTA")
	sb.WriteString(fmt.Sprintf(" :size %q", n.Size))
	if n.Tablespace != nil {
		sb.WriteString(" :on ")
		writeNode(sb, n.Tablespace)
	}
	sb.WriteString("}")
}

func writeDefaultRoleClause(sb *strings.Builder, n *DefaultRoleClause) {
	sb.WriteString("{DEFAULT_ROLE")
	if n.AllRoles {
		if n.ExceptAll {
			sb.WriteString(" :allExcept true")
		} else {
			sb.WriteString(" :all true")
		}
	}
	if n.NoneRole {
		sb.WriteString(" :none true")
	}
	for _, r := range n.Roles {
		sb.WriteString(" :role ")
		writeNode(sb, r)
	}
	sb.WriteString("}")
}

func writeProfileLimit(sb *strings.Builder, n *ProfileLimit) {
	sb.WriteString("{LIMIT")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	sb.WriteString(fmt.Sprintf(" :value %q", n.Value))
	sb.WriteString("}")
}

func writeCreateUserStmt(sb *strings.Builder, n *CreateUserStmt) {
	sb.WriteString("{CREATE_USER")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Identified != nil {
		sb.WriteString(" :identified ")
		writeIdentifiedClause(sb, n.Identified)
	}
	if n.DefaultCollation != "" {
		sb.WriteString(fmt.Sprintf(" :defaultCollation %q", n.DefaultCollation))
	}
	if n.DefaultTablespace != "" {
		sb.WriteString(fmt.Sprintf(" :defaultTablespace %q", n.DefaultTablespace))
	}
	if n.TempTablespace != "" {
		if n.LocalTemp {
			sb.WriteString(" :localTemp true")
		}
		sb.WriteString(fmt.Sprintf(" :tempTablespace %q", n.TempTablespace))
	}
	for _, q := range n.Quotas {
		sb.WriteString(" :quota ")
		writeUserQuotaClause(sb, q)
	}
	if n.Profile != "" {
		sb.WriteString(fmt.Sprintf(" :profile %q", n.Profile))
	}
	if n.PasswordExpire {
		sb.WriteString(" :passwordExpire true")
	}
	if n.AccountLock != nil {
		if *n.AccountLock {
			sb.WriteString(" :accountLock true")
		} else {
			sb.WriteString(" :accountUnlock true")
		}
	}
	if n.EnableEditions {
		sb.WriteString(" :enableEditions true")
	}
	if n.ContainerAll != nil {
		if *n.ContainerAll {
			sb.WriteString(" :container ALL")
		} else {
			sb.WriteString(" :container CURRENT")
		}
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterUserStmt(sb *strings.Builder, n *AlterUserStmt) {
	sb.WriteString("{ALTER_USER")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Identified != nil {
		sb.WriteString(" :identified ")
		writeIdentifiedClause(sb, n.Identified)
	}
	if n.DefaultCollation != "" {
		sb.WriteString(fmt.Sprintf(" :defaultCollation %q", n.DefaultCollation))
	}
	if n.DefaultTablespace != "" {
		sb.WriteString(fmt.Sprintf(" :defaultTablespace %q", n.DefaultTablespace))
	}
	if n.TempTablespace != "" {
		if n.LocalTemp {
			sb.WriteString(" :localTemp true")
		}
		sb.WriteString(fmt.Sprintf(" :tempTablespace %q", n.TempTablespace))
	}
	for _, q := range n.Quotas {
		sb.WriteString(" :quota ")
		writeUserQuotaClause(sb, q)
	}
	if n.Profile != "" {
		sb.WriteString(fmt.Sprintf(" :profile %q", n.Profile))
	}
	if n.DefaultRole != nil {
		sb.WriteString(" :defaultRole ")
		writeDefaultRoleClause(sb, n.DefaultRole)
	}
	if n.PasswordExpire {
		sb.WriteString(" :passwordExpire true")
	}
	if n.AccountLock != nil {
		if *n.AccountLock {
			sb.WriteString(" :accountLock true")
		} else {
			sb.WriteString(" :accountUnlock true")
		}
	}
	if n.EnableEditions {
		sb.WriteString(" :enableEditions true")
	}
	if n.ContainerAll != nil {
		if *n.ContainerAll {
			sb.WriteString(" :container ALL")
		} else {
			sb.WriteString(" :container CURRENT")
		}
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateRoleStmt(sb *strings.Builder, n *CreateRoleStmt) {
	sb.WriteString("{CREATE_ROLE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IdentifyBy != "" {
		sb.WriteString(fmt.Sprintf(" :identifiedBy %q", n.IdentifyBy))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateProfileStmt(sb *strings.Builder, n *CreateProfileStmt) {
	sb.WriteString("{CREATE_PROFILE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Mandatory {
		sb.WriteString(" :mandatory true")
	}
	for _, lim := range n.Limits {
		sb.WriteString(" :limit ")
		writeProfileLimit(sb, lim)
	}
	if n.ContainerAll != nil {
		if *n.ContainerAll {
			sb.WriteString(" :container ALL")
		} else {
			sb.WriteString(" :container CURRENT")
		}
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAdminDDLStmt(sb *strings.Builder, n *AdminDDLStmt) {
	sb.WriteString("{ADMIN_DDL")
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	sb.WriteString(fmt.Sprintf(" :objectType %d", n.ObjectType))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Options != nil && len(n.Options.Items) > 0 {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDDLOption(sb *strings.Builder, n *DDLOption) {
	sb.WriteString("{DDL_OPTION")
	sb.WriteString(fmt.Sprintf(" :key %q", n.Key))
	if n.Value != "" {
		sb.WriteString(fmt.Sprintf(" :value %q", n.Value))
	}
	if n.Items != nil && len(n.Items.Items) > 0 {
		sb.WriteString(" :items ")
		writeNode(sb, n.Items)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateSchemaStmt(sb *strings.Builder, n *CreateSchemaStmt) {
	sb.WriteString("{CREATE_SCHEMA")
	sb.WriteString(fmt.Sprintf(" :schemaName %q", n.SchemaName))
	if n.Stmts != nil && len(n.Stmts.Items) > 0 {
		sb.WriteString(" :stmts ")
		writeNode(sb, n.Stmts)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterDatabaseLinkStmt(sb *strings.Builder, n *AlterDatabaseLinkStmt) {
	sb.WriteString("{ALTER_DATABASE_LINK")
	if n.Shared {
		sb.WriteString(" :shared true")
	}
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.ConnectUser != "" {
		sb.WriteString(fmt.Sprintf(" :connectUser %q", n.ConnectUser))
	}
	if n.ConnectPassword != "" {
		sb.WriteString(fmt.Sprintf(" :connectPassword %q", n.ConnectPassword))
	}
	if n.AuthenticatedUser != "" {
		sb.WriteString(fmt.Sprintf(" :authenticatedUser %q", n.AuthenticatedUser))
	}
	if n.AuthenticatedPass != "" {
		sb.WriteString(fmt.Sprintf(" :authenticatedPass %q", n.AuthenticatedPass))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterSynonymStmt(sb *strings.Builder, n *AlterSynonymStmt) {
	sb.WriteString("{ALTER_SYNONYM")
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterMaterializedViewStmt(sb *strings.Builder, n *AlterMaterializedViewStmt) {
	sb.WriteString("{ALTER_MATERIALIZED_VIEW")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.RefreshMethod != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMethod %q", n.RefreshMethod))
	}
	if n.RefreshMode != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMode %q", n.RefreshMode))
	}
	if n.StartWith != nil {
		sb.WriteString(" :startWith ")
		writeNode(sb, n.StartWith)
	}
	if n.Next != nil {
		sb.WriteString(" :next ")
		writeNode(sb, n.Next)
	}
	if n.WithPrimaryKey {
		sb.WriteString(" :withPrimaryKey true")
	}
	if n.UsingRollbackSegment != "" {
		sb.WriteString(fmt.Sprintf(" :usingRollbackSegment %q", n.UsingRollbackSegment))
	}
	if n.UsingConstraints != "" {
		sb.WriteString(fmt.Sprintf(" :usingConstraints %q", n.UsingConstraints))
	}
	if n.EnableOnQueryComputation {
		sb.WriteString(" :enableOnQueryComputation true")
	}
	if n.DisableOnQueryComputation {
		sb.WriteString(" :disableOnQueryComputation true")
	}
	if n.Compact {
		sb.WriteString(" :compact true")
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.ParallelDegree != "" {
		sb.WriteString(fmt.Sprintf(" :parallelDegree %q", n.ParallelDegree))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetRoleStmt(sb *strings.Builder, n *SetRoleStmt) {
	sb.WriteString("{SET_ROLE")
	if n.All {
		sb.WriteString(" :all true")
	}
	if n.None {
		sb.WriteString(" :none true")
	}
	if len(n.Roles) > 0 {
		sb.WriteString(" :roles (")
		for i, r := range n.Roles {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, r)
		}
		sb.WriteString(")")
	}
	if len(n.Except) > 0 {
		sb.WriteString(" :except (")
		for i, r := range n.Except {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, r)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeSetConstraintsStmt(sb *strings.Builder, n *SetConstraintsStmt) {
	sb.WriteString("{SET_CONSTRAINTS")
	if n.All {
		sb.WriteString(" :all true")
	}
	if len(n.Constraints) > 0 {
		sb.WriteString(" :constraints (")
		for i, c := range n.Constraints {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, c)
		}
		sb.WriteString(")")
	}
	if n.Deferred {
		sb.WriteString(" :deferred true")
	} else {
		sb.WriteString(" :immediate true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAuditStmt(sb *strings.Builder, n *AuditStmt) {
	sb.WriteString("{AUDIT")
	if len(n.Actions) > 0 {
		sb.WriteString(" :actions (")
		for i, a := range n.Actions {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
	if n.Object != nil {
		sb.WriteString(" :object ")
		writeNode(sb, n.Object)
	}
	if n.By != "" {
		sb.WriteString(fmt.Sprintf(" :by %q", n.By))
	}
	if n.When != "" {
		sb.WriteString(fmt.Sprintf(" :when %q", n.When))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeNoauditStmt(sb *strings.Builder, n *NoauditStmt) {
	sb.WriteString("{NOAUDIT")
	if len(n.Actions) > 0 {
		sb.WriteString(" :actions (")
		for i, a := range n.Actions {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
	if n.Object != nil {
		sb.WriteString(" :object ")
		writeNode(sb, n.Object)
	}
	if n.When != "" {
		sb.WriteString(fmt.Sprintf(" :when %q", n.When))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAssociateStatisticsStmt(sb *strings.Builder, n *AssociateStatisticsStmt) {
	sb.WriteString("{ASSOCIATE_STATISTICS")
	if n.ObjectType != "" {
		sb.WriteString(fmt.Sprintf(" :objectType %q", n.ObjectType))
	}
	if len(n.Objects) > 0 {
		sb.WriteString(" :objects (")
		for i, o := range n.Objects {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString(")")
	}
	if n.Using != nil {
		sb.WriteString(" :using ")
		writeNode(sb, n.Using)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDisassociateStatisticsStmt(sb *strings.Builder, n *DisassociateStatisticsStmt) {
	sb.WriteString("{DISASSOCIATE_STATISTICS")
	if n.ObjectType != "" {
		sb.WriteString(fmt.Sprintf(" :objectType %q", n.ObjectType))
	}
	if len(n.Objects) > 0 {
		sb.WriteString(" :objects (")
		for i, o := range n.Objects {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString(")")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// ---------------------------------------------------------------------------
// PL/SQL nodes
// ---------------------------------------------------------------------------

func writePLSQLBlock(sb *strings.Builder, n *PLSQLBlock) {
	sb.WriteString("{PLSQLBLOCK")
	if n.Label != "" {
		sb.WriteString(fmt.Sprintf(" :label %q", n.Label))
	}
	if n.Declarations != nil {
		sb.WriteString(" :declarations ")
		writeNode(sb, n.Declarations)
	}
	if n.Statements != nil {
		sb.WriteString(" :statements ")
		writeNode(sb, n.Statements)
	}
	if n.Exceptions != nil {
		sb.WriteString(" :exceptions ")
		writeNode(sb, n.Exceptions)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeExceptionHandler(sb *strings.Builder, n *ExceptionHandler) {
	sb.WriteString("{EXCEPTIONHANDLER")
	if n.Exceptions != nil {
		sb.WriteString(" :exceptions ")
		writeNode(sb, n.Exceptions)
	}
	if n.Statements != nil {
		sb.WriteString(" :statements ")
		writeNode(sb, n.Statements)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLIf(sb *strings.Builder, n *PLSQLIf) {
	sb.WriteString("{PLSQLIF")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Then != nil {
		sb.WriteString(" :then ")
		writeNode(sb, n.Then)
	}
	if n.ElsIfs != nil {
		sb.WriteString(" :elsifs ")
		writeNode(sb, n.ElsIfs)
	}
	if n.Else != nil {
		sb.WriteString(" :else ")
		writeNode(sb, n.Else)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLElsIf(sb *strings.Builder, n *PLSQLElsIf) {
	sb.WriteString("{PLSQLELSIF")
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Then != nil {
		sb.WriteString(" :then ")
		writeNode(sb, n.Then)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLLoop(sb *strings.Builder, n *PLSQLLoop) {
	sb.WriteString("{PLSQLLOOP")
	sb.WriteString(fmt.Sprintf(" :type %d", n.Type))
	if n.Label != "" {
		sb.WriteString(fmt.Sprintf(" :label %q", n.Label))
	}
	if n.Condition != nil {
		sb.WriteString(" :condition ")
		writeNode(sb, n.Condition)
	}
	if n.Iterator != "" {
		sb.WriteString(fmt.Sprintf(" :iterator %q", n.Iterator))
	}
	if n.LowerBound != nil {
		sb.WriteString(" :lowerBound ")
		writeNode(sb, n.LowerBound)
	}
	if n.UpperBound != nil {
		sb.WriteString(" :upperBound ")
		writeNode(sb, n.UpperBound)
	}
	if n.Reverse {
		sb.WriteString(" :reverse true")
	}
	if n.CursorName != "" {
		sb.WriteString(fmt.Sprintf(" :cursorName %q", n.CursorName))
	}
	if n.CursorArgs != nil {
		sb.WriteString(" :cursorArgs ")
		writeNode(sb, n.CursorArgs)
	}
	if n.Statements != nil {
		sb.WriteString(" :statements ")
		writeNode(sb, n.Statements)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLReturn(sb *strings.Builder, n *PLSQLReturn) {
	sb.WriteString("{PLSQLRETURN")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLGoto(sb *strings.Builder, n *PLSQLGoto) {
	sb.WriteString("{PLSQLGOTO")
	sb.WriteString(fmt.Sprintf(" :label %q", n.Label))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLAssign(sb *strings.Builder, n *PLSQLAssign) {
	sb.WriteString("{PLSQLASSIGN")
	if n.Target != nil {
		sb.WriteString(" :target ")
		writeNode(sb, n.Target)
	}
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLRaise(sb *strings.Builder, n *PLSQLRaise) {
	sb.WriteString("{PLSQLRAISE")
	if n.Exception != "" {
		sb.WriteString(fmt.Sprintf(" :exception %q", n.Exception))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLNull(sb *strings.Builder, n *PLSQLNull) {
	sb.WriteString("{PLSQLNULL")
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLVarDecl(sb *strings.Builder, n *PLSQLVarDecl) {
	sb.WriteString("{PLSQLVARDECL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	if n.Constant {
		sb.WriteString(" :constant true")
	}
	if n.NotNull {
		sb.WriteString(" :notNull true")
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLCursorDecl(sb *strings.Builder, n *PLSQLCursorDecl) {
	sb.WriteString("{PLSQLCURSORDECL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Parameters != nil {
		sb.WriteString(" :parameters ")
		writeNode(sb, n.Parameters)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLExecImmediate(sb *strings.Builder, n *PLSQLExecImmediate) {
	sb.WriteString("{PLSQLEXECIMMEDIATE")
	if n.SQL != nil {
		sb.WriteString(" :sql ")
		writeNode(sb, n.SQL)
	}
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.Using != nil {
		sb.WriteString(" :using ")
		writeNode(sb, n.Using)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLOpen(sb *strings.Builder, n *PLSQLOpen) {
	sb.WriteString("{PLSQLOPEN")
	sb.WriteString(fmt.Sprintf(" :cursor %q", n.Cursor))
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.ForQuery != nil {
		sb.WriteString(" :forQuery ")
		writeNode(sb, n.ForQuery)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLFetch(sb *strings.Builder, n *PLSQLFetch) {
	sb.WriteString("{PLSQLFETCH")
	sb.WriteString(fmt.Sprintf(" :cursor %q", n.Cursor))
	if n.Into != nil {
		sb.WriteString(" :into ")
		writeNode(sb, n.Into)
	}
	if n.Bulk {
		sb.WriteString(" :bulk true")
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLClose(sb *strings.Builder, n *PLSQLClose) {
	sb.WriteString("{PLSQLCLOSE")
	sb.WriteString(fmt.Sprintf(" :cursor %q", n.Cursor))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLForall(sb *strings.Builder, n *PLSQLForall) {
	sb.WriteString("{PLSQLFORALL")
	sb.WriteString(fmt.Sprintf(" :index %q", n.Index))
	if n.Lower != nil {
		sb.WriteString(" :lower ")
		writeNode(sb, n.Lower)
	}
	if n.Upper != nil {
		sb.WriteString(" :upper ")
		writeNode(sb, n.Upper)
	}
	if n.Body != nil {
		sb.WriteString(" :body ")
		writeNode(sb, n.Body)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLExit(sb *strings.Builder, n *PLSQLExit) {
	sb.WriteString("{PLSQLEXIT")
	if n.Label != "" {
		sb.WriteString(fmt.Sprintf(" :label %q", n.Label))
	}
	if n.Condition != nil {
		sb.WriteString(" :when ")
		writeNode(sb, n.Condition)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLContinue(sb *strings.Builder, n *PLSQLContinue) {
	sb.WriteString("{PLSQLCONTINUE")
	if n.Label != "" {
		sb.WriteString(fmt.Sprintf(" :label %q", n.Label))
	}
	if n.Condition != nil {
		sb.WriteString(" :when ")
		writeNode(sb, n.Condition)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLPipeRow(sb *strings.Builder, n *PLSQLPipeRow) {
	sb.WriteString("{PLSQLPIPEROW")
	if n.Row != nil {
		sb.WriteString(" :row ")
		writeNode(sb, n.Row)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLPragma(sb *strings.Builder, n *PLSQLPragma) {
	sb.WriteString("{PLSQLPRAGMA")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Args != nil && len(n.Args.Items) > 0 {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLCase(sb *strings.Builder, n *PLSQLCase) {
	sb.WriteString("{PLSQLCASE")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
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
	if len(n.Else) > 0 {
		sb.WriteString(" :else (")
		for i, s := range n.Else {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLWhen(sb *strings.Builder, n *PLSQLWhen) {
	sb.WriteString("{PLSQLWHEN")
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
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
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLTypeDecl(sb *strings.Builder, n *PLSQLTypeDecl) {
	sb.WriteString("{PLSQLTYPEDECL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	sb.WriteString(fmt.Sprintf(" :kind %q", n.Kind))
	if n.ElementType != nil {
		sb.WriteString(" :elementType ")
		writeNode(sb, n.ElementType)
	}
	if n.IndexBy != nil {
		sb.WriteString(" :indexBy ")
		writeNode(sb, n.IndexBy)
	}
	if n.Limit != nil {
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit)
	}
	if n.Fields != nil && len(n.Fields.Items) > 0 {
		sb.WriteString(" :fields ")
		writeNode(sb, n.Fields)
	}
	if n.ReturnType != nil {
		sb.WriteString(" :returnType ")
		writeNode(sb, n.ReturnType)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writePLSQLCall(sb *strings.Builder, n *PLSQLCall) {
	sb.WriteString("{PLSQLCALL")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTablespaceStmt(sb *strings.Builder, n *CreateTablespaceStmt) {
	sb.WriteString("{CREATE_TABLESPACE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Bigfile {
		sb.WriteString(" :bigfile true")
	}
	if n.Smallfile {
		sb.WriteString(" :smallfile true")
	}
	if n.Temporary {
		sb.WriteString(" :temporary true")
	}
	if n.Undo {
		sb.WriteString(" :undo true")
	}
	if len(n.Datafiles) > 0 {
		sb.WriteString(" :datafiles (")
		for i, d := range n.Datafiles {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, d)
		}
		sb.WriteString(")")
	}
	if n.Size != "" {
		sb.WriteString(fmt.Sprintf(" :size %q", n.Size))
	}
	if n.Autoextend != nil {
		sb.WriteString(" :autoextend ")
		writeNode(sb, n.Autoextend)
	}
	if n.Logging != "" {
		sb.WriteString(fmt.Sprintf(" :logging %q", n.Logging))
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Offline {
		sb.WriteString(" :offline true")
	}
	if n.Extent != "" {
		sb.WriteString(fmt.Sprintf(" :extent %q", n.Extent))
	}
	if n.Segment != "" {
		sb.WriteString(fmt.Sprintf(" :segment %q", n.Segment))
	}
	if n.Blocksize != "" {
		sb.WriteString(fmt.Sprintf(" :blocksize %q", n.Blocksize))
	}
	if n.Retention != "" {
		sb.WriteString(fmt.Sprintf(" :retention %q", n.Retention))
	}
	if n.Encryption != "" {
		sb.WriteString(fmt.Sprintf(" :encryption %q", n.Encryption))
	}
	if n.Compress != "" {
		sb.WriteString(fmt.Sprintf(" :compress %q", n.Compress))
	}
	if n.MaxSize != "" {
		sb.WriteString(fmt.Sprintf(" :maxsize %q", n.MaxSize))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDatafileClause(sb *strings.Builder, n *DatafileClause) {
	sb.WriteString("{DATAFILE")
	if n.Filename != "" {
		sb.WriteString(fmt.Sprintf(" :filename %q", n.Filename))
	}
	if n.Size != "" {
		sb.WriteString(fmt.Sprintf(" :size %q", n.Size))
	}
	if n.Reuse {
		sb.WriteString(" :reuse true")
	}
	if n.Autoextend != nil {
		sb.WriteString(" :autoextend ")
		writeNode(sb, n.Autoextend)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAutoextendClause(sb *strings.Builder, n *AutoextendClause) {
	sb.WriteString("{AUTOEXTEND")
	if n.On {
		sb.WriteString(" :on true")
	}
	if n.Next != "" {
		sb.WriteString(fmt.Sprintf(" :next %q", n.Next))
	}
	if n.MaxSize != "" {
		sb.WriteString(fmt.Sprintf(" :maxsize %q", n.MaxSize))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateClusterStmt(sb *strings.Builder, n *CreateClusterStmt) {
	sb.WriteString("{CREATE_CLUSTER")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
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
	if n.PctFree != nil {
		sb.WriteString(fmt.Sprintf(" :pctfree %d", *n.PctFree))
	}
	if n.PctUsed != nil {
		sb.WriteString(fmt.Sprintf(" :pctused %d", *n.PctUsed))
	}
	if n.InitTrans != nil {
		sb.WriteString(fmt.Sprintf(" :initrans %d", *n.InitTrans))
	}
	if n.Size != "" {
		sb.WriteString(fmt.Sprintf(" :size %q", n.Size))
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	if n.IsIndex {
		sb.WriteString(" :isIndex true")
	}
	if n.IsHash {
		sb.WriteString(" :isHash true")
	}
	if n.HashKeys != "" {
		sb.WriteString(fmt.Sprintf(" :hashkeys %q", n.HashKeys))
	}
	if n.SingleTable {
		sb.WriteString(" :singleTable true")
	}
	if n.HashExpr != nil {
		sb.WriteString(" :hashExpr ")
		writeNode(sb, n.HashExpr)
	}
	if n.Cache {
		sb.WriteString(" :cache true")
	}
	if n.NoCache {
		sb.WriteString(" :nocache true")
	}
	if n.Parallel != "" {
		sb.WriteString(fmt.Sprintf(" :parallel %q", n.Parallel))
	}
	if n.RowDep {
		sb.WriteString(" :rowdep true")
	}
	if n.NoRowDep {
		sb.WriteString(" :norowdep true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeClusterColumn(sb *strings.Builder, n *ClusterColumn) {
	sb.WriteString("{CLUSTER_COL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.DataType != nil {
		sb.WriteString(" :datatype ")
		writeNode(sb, n.DataType)
	}
	if n.Sort {
		sb.WriteString(" :sort true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateDimensionStmt(sb *strings.Builder, n *CreateDimensionStmt) {
	sb.WriteString("{CREATE_DIMENSION")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Levels) > 0 {
		sb.WriteString(" :levels (")
		for i, l := range n.Levels {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, l)
		}
		sb.WriteString(")")
	}
	if len(n.Hierarchies) > 0 {
		sb.WriteString(" :hierarchies (")
		for i, h := range n.Hierarchies {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, h)
		}
		sb.WriteString(")")
	}
	if len(n.Attributes) > 0 {
		sb.WriteString(" :attributes (")
		for i, a := range n.Attributes {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDimensionLevel(sb *strings.Builder, n *DimensionLevel) {
	sb.WriteString("{DIM_LEVEL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
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
	if n.SkipWhenNull {
		sb.WriteString(" :skipWhenNull true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDimensionHierarchy(sb *strings.Builder, n *DimensionHierarchy) {
	sb.WriteString("{DIM_HIERARCHY")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if len(n.Levels) > 0 {
		sb.WriteString(" :levels (")
		for i, l := range n.Levels {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", l))
		}
		sb.WriteString(")")
	}
	if len(n.JoinKeys) > 0 {
		sb.WriteString(" :joinKeys (")
		for i, j := range n.JoinKeys {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, j)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDimensionJoinKey(sb *strings.Builder, n *DimensionJoinKey) {
	sb.WriteString("{DIM_JOIN_KEY")
	if len(n.ChildKeys) > 0 {
		sb.WriteString(" :childKeys (")
		for i, k := range n.ChildKeys {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, k)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :parentLevel %q", n.ParentLevel))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDimensionAttribute(sb *strings.Builder, n *DimensionAttribute) {
	sb.WriteString("{DIM_ATTR")
	sb.WriteString(fmt.Sprintf(" :attrName %q", n.AttrName))
	if n.LevelName != "" {
		sb.WriteString(fmt.Sprintf(" :levelName %q", n.LevelName))
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
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterIndexStmt(sb *strings.Builder, n *AlterIndexStmt) {
	sb.WriteString("{ALTER_INDEX")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	if n.Partition != "" {
		sb.WriteString(fmt.Sprintf(" :partition %q", n.Partition))
	}
	if n.Subpartition != "" {
		sb.WriteString(fmt.Sprintf(" :subpartition %q", n.Subpartition))
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Reverse {
		sb.WriteString(" :reverse true")
	}
	if n.NoReverse {
		sb.WriteString(" :noreverse true")
	}
	if n.NewName != "" {
		sb.WriteString(fmt.Sprintf(" :newName %q", n.NewName))
	}
	if n.Cleanup {
		sb.WriteString(" :cleanup true")
	}
	if n.CleanupOnly {
		sb.WriteString(" :cleanupOnly true")
	}
	if n.Compact {
		sb.WriteString(" :compact true")
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.Parallel != "" {
		sb.WriteString(fmt.Sprintf(" :parallel %q", n.Parallel))
	}
	if n.NoParallel {
		sb.WriteString(" :noparallel true")
	}
	if n.Logging {
		sb.WriteString(" :logging true")
	}
	if n.NoLogging {
		sb.WriteString(" :nologging true")
	}
	if n.Compress != "" {
		sb.WriteString(fmt.Sprintf(" :compress %q", n.Compress))
	}
	if n.NoCompress {
		sb.WriteString(" :nocompress true")
	}
	if n.IndexingFull {
		sb.WriteString(" :indexingFull true")
	}
	if n.IndexingPartial {
		sb.WriteString(" :indexingPartial true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterViewStmt(sb *strings.Builder, n *AlterViewStmt) {
	sb.WriteString("{ALTER_VIEW")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	if n.Constraint != nil {
		sb.WriteString(" :constraint ")
		writeNode(sb, n.Constraint)
	}
	if n.ConstraintName != "" {
		sb.WriteString(fmt.Sprintf(" :constraintName %q", n.ConstraintName))
	}
	if n.Rely {
		sb.WriteString(" :rely true")
	}
	if n.NoRely {
		sb.WriteString(" :norely true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterSequenceStmt(sb *strings.Builder, n *AlterSequenceStmt) {
	sb.WriteString("{ALTER_SEQUENCE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.IncrementBy != nil {
		sb.WriteString(" :incrementBy ")
		writeNode(sb, n.IncrementBy)
	}
	if n.MaxValue != nil {
		sb.WriteString(" :maxValue ")
		writeNode(sb, n.MaxValue)
	}
	if n.MinValue != nil {
		sb.WriteString(" :minValue ")
		writeNode(sb, n.MinValue)
	}
	if n.NoMaxValue {
		sb.WriteString(" :noMaxValue true")
	}
	if n.NoMinValue {
		sb.WriteString(" :noMinValue true")
	}
	if n.Cycle {
		sb.WriteString(" :cycle true")
	}
	if n.NoCycle {
		sb.WriteString(" :noCycle true")
	}
	if n.Cache != nil {
		sb.WriteString(" :cache ")
		writeNode(sb, n.Cache)
	}
	if n.NoCache {
		sb.WriteString(" :noCache true")
	}
	if n.Order {
		sb.WriteString(" :order true")
	}
	if n.NoOrder {
		sb.WriteString(" :noOrder true")
	}
	if n.Keep {
		sb.WriteString(" :keep true")
	}
	if n.NoKeep {
		sb.WriteString(" :noKeep true")
	}
	if n.Restart {
		sb.WriteString(" :restart true")
	}
	if n.RestartWith != nil {
		sb.WriteString(" :restartWith ")
		writeNode(sb, n.RestartWith)
	}
	if n.Scale {
		sb.WriteString(" :scale true")
	}
	if n.NoScale {
		sb.WriteString(" :noScale true")
	}
	if n.ScaleExtend {
		sb.WriteString(" :scaleExtend true")
	}
	if n.ScaleNoExtend {
		sb.WriteString(" :scaleNoExtend true")
	}
	if n.Shard {
		sb.WriteString(" :shard true")
	}
	if n.NoShard {
		sb.WriteString(" :noShard true")
	}
	if n.ShardExtend {
		sb.WriteString(" :shardExtend true")
	}
	if n.ShardNoExtend {
		sb.WriteString(" :shardNoExtend true")
	}
	if n.Global {
		sb.WriteString(" :global true")
	}
	if n.Session {
		sb.WriteString(" :session true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCompileFields(sb *strings.Builder, compile bool, debug bool, reuseSettings bool, compilerParams []*SetParam) {
	if compile {
		sb.WriteString(" :compile true")
	}
	if debug {
		sb.WriteString(" :debug true")
	}
	if reuseSettings {
		sb.WriteString(" :reuseSettings true")
	}
	if len(compilerParams) > 0 {
		sb.WriteString(" :compilerParams (")
		for i, p := range compilerParams {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString(")")
	}
}

func writeAlterProcedureStmt(sb *strings.Builder, n *AlterProcedureStmt) {
	sb.WriteString("{ALTER_PROCEDURE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	writeCompileFields(sb, n.Compile, n.Debug, n.ReuseSettings, n.CompilerParams)
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterFunctionStmt(sb *strings.Builder, n *AlterFunctionStmt) {
	sb.WriteString("{ALTER_FUNCTION")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	writeCompileFields(sb, n.Compile, n.Debug, n.ReuseSettings, n.CompilerParams)
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterPackageStmt(sb *strings.Builder, n *AlterPackageStmt) {
	sb.WriteString("{ALTER_PACKAGE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	writeCompileFields(sb, n.Compile, n.Debug, n.ReuseSettings, n.CompilerParams)
	if n.CompileTarget != "" {
		sb.WriteString(fmt.Sprintf(" :compileTarget %q", n.CompileTarget))
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTriggerStmt(sb *strings.Builder, n *AlterTriggerStmt) {
	sb.WriteString("{ALTER_TRIGGER")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	writeCompileFields(sb, n.Action == "COMPILE", n.Debug, n.ReuseSettings, n.CompilerParams)
	if n.NewName != "" {
		sb.WriteString(fmt.Sprintf(" :newName %q", n.NewName))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTypeStmt(sb *strings.Builder, n *AlterTypeStmt) {
	sb.WriteString("{ALTER_TYPE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.CompileTarget != "" {
		sb.WriteString(fmt.Sprintf(" :compileTarget %q", n.CompileTarget))
	}
	writeCompileFields(sb, n.Action == "COMPILE", n.Debug, n.ReuseSettings, n.CompilerParams)
	if len(n.Attributes) > 0 {
		sb.WriteString(" :attributes (")
		for i, a := range n.Attributes {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString(")")
	}
	if n.MethodKind != "" {
		sb.WriteString(fmt.Sprintf(" :methodKind %q", n.MethodKind))
	}
	if n.MethodType != "" {
		sb.WriteString(fmt.Sprintf(" :methodType %q", n.MethodType))
	}
	if n.MethodName != "" {
		sb.WriteString(fmt.Sprintf(" :methodName %q", n.MethodName))
	}
	if len(n.MethodParams) > 0 {
		sb.WriteString(" :methodParams (")
		for i, p := range n.MethodParams {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString(")")
	}
	if n.MethodReturn != nil {
		sb.WriteString(" :methodReturn ")
		writeNode(sb, n.MethodReturn)
	}
	if n.LimitValue != nil {
		sb.WriteString(" :limitValue ")
		writeNode(sb, n.LimitValue)
	}
	if n.ElementType != nil {
		sb.WriteString(" :elementType ")
		writeNode(sb, n.ElementType)
	}
	if n.Invalidate {
		sb.WriteString(" :invalidate true")
	}
	if n.Cascade {
		sb.WriteString(" :cascade true")
	}
	if n.IncludeData != nil {
		if *n.IncludeData {
			sb.WriteString(" :includeData true")
		} else {
			sb.WriteString(" :includeData false")
		}
	}
	if n.ConvertToSubst {
		sb.WriteString(" :convertToSubst true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeTypeAttribute(sb *strings.Builder, n *TypeAttribute) {
	sb.WriteString("{TYPE_ATTR")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

// escapeString escapes special characters in a string for output.
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
