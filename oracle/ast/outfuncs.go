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
		sb.WriteString(fmt.Sprintf("{BINDVAR :name %q", n.Name))
		if n.Member != "" {
			sb.WriteString(fmt.Sprintf(" :member %q", n.Member))
		}
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d}", n.Loc.Start, n.Loc.End))
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
	case *ExtractExpr:
		writeExtractExpr(sb, n)
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
	case *DateTimeLiteral:
		sb.WriteString("{DATETIMELIT")
		sb.WriteString(fmt.Sprintf(" :typeName %q :val %q", n.TypeName, n.Val))
		sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
		sb.WriteString("}")
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
	case *LockTableItem:
		writeLockTableItem(sb, n)
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
	case *AlterRoleStmt:
		writeAlterRoleStmt(sb, n)
	case *CreateProfileStmt:
		writeCreateProfileStmt(sb, n)
	case *AlterProfileStmt:
		writeAlterProfileStmt(sb, n)
	case *AlterResourceCostStmt:
		writeAlterResourceCostStmt(sb, n)
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
	case *CreateAuditPolicyStmt:
		writeCreateAuditPolicyStmt(sb, n)
	case *AlterAuditPolicyStmt:
		writeAlterAuditPolicyStmt(sb, n)
	case *DropAuditPolicyStmt:
		writeDropAuditPolicyStmt(sb, n)
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
	case *AlterTablespaceStmt:
		writeAlterTablespaceStmt(sb, n)
	case *CreateTablespaceSetStmt:
		writeCreateTablespaceSetStmt(sb, n)
	case *DropTablespaceStmt:
		writeDropTablespaceStmt(sb, n)
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
	case *AlterClusterStmt:
		writeAlterClusterStmt(sb, n)
	case *AlterDimensionStmt:
		writeAlterDimensionStmt(sb, n)
	case *CreateMaterializedZonemapStmt:
		writeCreateMaterializedZonemapStmt(sb, n)
	case *AlterMaterializedZonemapStmt:
		writeAlterMaterializedZonemapStmt(sb, n)
	case *CreateInmemoryJoinGroupStmt:
		writeCreateInmemoryJoinGroupStmt(sb, n)
	case *AlterInmemoryJoinGroupStmt:
		writeAlterInmemoryJoinGroupStmt(sb, n)
	case *JoinGroupMember:
		writeJoinGroupMember(sb, n)
	case *CreateIndextypeStmt:
		writeCreateIndextypeStmt(sb, n)
	case *AlterIndextypeStmt:
		writeAlterIndextypeStmt(sb, n)
	case *IndextypeOp:
		writeIndextypeOp(sb, n)
	case *IndextypeModOp:
		writeIndextypeModOp(sb, n)
	case *CreateOperatorStmt:
		writeCreateOperatorStmt(sb, n)
	case *AlterOperatorStmt:
		writeAlterOperatorStmt(sb, n)
	case *OperatorBinding:
		writeOperatorBinding(sb, n)
	case *CreateMviewLogStmt:
		writeCreateMviewLogStmt(sb, n)
	case *AlterMviewLogStmt:
		writeAlterMviewLogStmt(sb, n)
	case *CreateAnalyticViewStmt:
		writeCreateAnalyticViewStmt(sb, n)
	case *AlterAnalyticViewStmt:
		writeAlterAnalyticViewStmt(sb, n)
	case *CreateJsonDualityViewStmt:
		writeCreateJsonDualityViewStmt(sb, n)
	case *AlterJsonDualityViewStmt:
		writeAlterJsonDualityViewStmt(sb, n)
	case *CreateAttributeDimensionStmt:
		writeCreateAttributeDimensionStmt(sb, n)
	case *AlterAttributeDimensionStmt:
		writeAlterAttributeDimensionStmt(sb, n)
	case *AttrDimSourceClause:
		writeAttrDimSourceClause(sb, n)
	case *AttrDimJoinCondElem:
		writeAttrDimJoinCondElem(sb, n)
	case *AttrDimAttribute:
		writeAttrDimAttribute(sb, n)
	case *AttrDimLevel:
		writeAttrDimLevel(sb, n)
	case *AttrDimOrderByItem:
		writeAttrDimOrderByItem(sb, n)
	case *AttrDimAllClause:
		writeAttrDimAllClause(sb, n)
	case *CreateHierarchyStmt:
		writeCreateHierarchyStmt(sb, n)
	case *AlterHierarchyStmt:
		writeAlterHierarchyStmt(sb, n)
	case *HierLevelClause:
		writeHierLevelClause(sb, n)
	case *HierAttr:
		writeHierAttr(sb, n)
	case *CreateDomainStmt:
		writeCreateDomainStmt(sb, n)
	case *AlterDomainStmt:
		writeAlterDomainStmt(sb, n)
	case *DomainColumn:
		writeDomainColumn(sb, n)
	case *DomainEnumItem:
		writeDomainEnumItem(sb, n)
	case *DomainConstraint:
		writeDomainConstraint(sb, n)
	case *CreatePropertyGraphStmt:
		writeCreatePropertyGraphStmt(sb, n)
	case *GraphTableDef:
		writeGraphTableDef(sb, n)
	case *GraphEdgeDef:
		writeGraphEdgeDef(sb, n)
	case *GraphOptions:
		writeGraphOptions(sb, n)
	case *CreateVectorIndexStmt:
		writeCreateVectorIndexStmt(sb, n)
	case *CreateLockdownProfileStmt:
		writeCreateLockdownProfileStmt(sb, n)
	case *AlterLockdownProfileStmt:
		writeAlterLockdownProfileStmt(sb, n)
	case *CreateOutlineStmt:
		writeCreateOutlineStmt(sb, n)
	case *AlterOutlineStmt:
		writeAlterOutlineStmt(sb, n)

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

func writeExtractExpr(sb *strings.Builder, n *ExtractExpr) {
	sb.WriteString("{EXTRACT")
	sb.WriteString(fmt.Sprintf(" :field %q", n.Field))
	if n.Expr != nil {
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr)
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
	if n.Domain != nil {
		sb.WriteString(" :domain ")
		writeNode(sb, n.Domain)
	}
	if n.Sort {
		sb.WriteString(" :sort true")
	}
	if n.Visible {
		sb.WriteString(" :visible true")
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	if n.DefaultOnNull {
		sb.WriteString(" :defaultOnNull true")
	}
	if n.DefaultOnNullInsertOnly {
		sb.WriteString(" :defaultOnNullInsertOnly true")
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
	if n.Encrypt != "" {
		sb.WriteString(fmt.Sprintf(" :encrypt %q", n.Encrypt))
	}
	if n.DropIdentity {
		sb.WriteString(" :dropIdentity true")
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
	if n.Interval != nil {
		sb.WriteString(" :interval ")
		writeNode(sb, n.Interval)
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
	if n.ColumnDefs != nil {
		sb.WriteString(" :columnDefs ")
		writeNode(sb, n.ColumnDefs)
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
	if n.Subtype != "" {
		sb.WriteString(fmt.Sprintf(" :subtype %q", n.Subtype))
	}
	if n.Options != nil && len(n.Options.Items) > 0 {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
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
	sb.WriteString(fmt.Sprintf(" :stmt_location %d", n.Loc.Start))
	sb.WriteString(fmt.Sprintf(" :stmt_len %d", n.Loc.End-n.Loc.Start))
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
	if n.IntoVars != nil {
		sb.WriteString(" :intoVars ")
		writeNode(sb, n.IntoVars)
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
	if n.Sharded {
		sb.WriteString(" :sharded true")
	}
	if n.Duplicated {
		sb.WriteString(" :duplicated true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
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
	if n.Collation != "" {
		sb.WriteString(fmt.Sprintf(" :collation %q", n.Collation))
	}
	if n.Organization != "" {
		sb.WriteString(fmt.Sprintf(" :organization %q", n.Organization))
	}
	if n.SegmentCreation != "" {
		sb.WriteString(fmt.Sprintf(" :segmentCreation %q", n.SegmentCreation))
	}
	if n.Logging != "" {
		sb.WriteString(fmt.Sprintf(" :logging %q", n.Logging))
	}
	if n.Cache != "" {
		sb.WriteString(fmt.Sprintf(" :cache %q", n.Cache))
	}
	if n.RowMovement != "" {
		sb.WriteString(fmt.Sprintf(" :rowMovement %q", n.RowMovement))
	}
	if n.ReadOnly != "" {
		sb.WriteString(fmt.Sprintf(" :readOnly %q", n.ReadOnly))
	}
	if n.Indexing != "" {
		sb.WriteString(fmt.Sprintf(" :indexing %q", n.Indexing))
	}
	if n.ResultCache != "" {
		sb.WriteString(fmt.Sprintf(" :resultCache %q", n.ResultCache))
	}
	if n.RowDependencies != "" {
		sb.WriteString(fmt.Sprintf(" :rowDependencies %q", n.RowDependencies))
	}
	if n.FlashbackArchive != "" {
		sb.WriteString(fmt.Sprintf(" :flashbackArchive %q", n.FlashbackArchive))
	}
	if n.ImmutableNoDrop != "" {
		sb.WriteString(fmt.Sprintf(" :immutableNoDrop %q", n.ImmutableNoDrop))
	}
	if n.ImmutableNoDel != "" {
		sb.WriteString(fmt.Sprintf(" :immutableNoDel %q", n.ImmutableNoDel))
	}
	if n.BlockchainHash != "" {
		sb.WriteString(fmt.Sprintf(" :blockchainHash %q", n.BlockchainHash))
	}
	if n.BlockchainVer != "" {
		sb.WriteString(fmt.Sprintf(" :blockchainVer %q", n.BlockchainVer))
	}
	if n.MemoptimizeRead {
		sb.WriteString(" :memoptimizeRead true")
	}
	if n.MemoptimizeWrite {
		sb.WriteString(" :memoptimizeWrite true")
	}
	if n.Parent != nil {
		sb.WriteString(" :parent ")
		writeNode(sb, n.Parent)
	}
	if n.Options != nil && len(n.Options.Items) > 0 {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
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
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.Validate {
		sb.WriteString(" :validate true")
	}
	if n.Invalidation != "" {
		sb.WriteString(fmt.Sprintf(" :invalidation %q", n.Invalidation))
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
	if n.Multivalue {
		sb.WriteString(" :multivalue true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Cluster != nil {
		sb.WriteString(" :cluster ")
		writeNode(sb, n.Cluster)
	}
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
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
	if n.Sort {
		sb.WriteString(" :sort true")
	}
	if n.NoSort {
		sb.WriteString(" :nosort true")
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
	if n.Logging {
		sb.WriteString(" :logging true")
	}
	if n.NoLogging {
		sb.WriteString(" :nologging true")
	}
	if n.Visible {
		sb.WriteString(" :visible true")
	}
	if n.Invisible {
		sb.WriteString(" :invisible true")
	}
	if n.IndexType != nil {
		sb.WriteString(" :indexType ")
		writeNode(sb, n.IndexType)
	}
	if n.Parameters != "" {
		sb.WriteString(fmt.Sprintf(" :parameters %q", n.Parameters))
	}
	if n.Invalidation != "" {
		sb.WriteString(fmt.Sprintf(" :invalidation %q", n.Invalidation))
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.FromTables != nil {
		sb.WriteString(" :fromTables ")
		writeNode(sb, n.FromTables)
	}
	if n.Where != nil {
		sb.WriteString(" :where ")
		writeNode(sb, n.Where)
	}
	if n.PctFree != "" {
		sb.WriteString(fmt.Sprintf(" :pctfree %q", n.PctFree))
	}
	if n.InitTrans != "" {
		sb.WriteString(fmt.Sprintf(" :initrans %q", n.InitTrans))
	}
	if n.MaxTrans != "" {
		sb.WriteString(fmt.Sprintf(" :maxtrans %q", n.MaxTrans))
	}
	if n.IndexingFull {
		sb.WriteString(" :indexingFull true")
	}
	if n.IndexingPartial {
		sb.WriteString(" :indexingPartial true")
	}
	if n.Options != nil && len(n.Options.Items) > 0 {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateViewStmt(sb *strings.Builder, n *CreateViewStmt) {
	sb.WriteString("{CREATEVIEW")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.NoForce {
		sb.WriteString(" :noForce true")
	}
	if n.Editioning {
		sb.WriteString(" :editioning true")
	}
	if n.Editionable != "" {
		sb.WriteString(fmt.Sprintf(" :editionable %q", n.Editionable))
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
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
	if n.ConstraintName != "" {
		sb.WriteString(fmt.Sprintf(" :constraintName %q", n.ConstraintName))
	}
	if n.Bequeath != "" {
		sb.WriteString(fmt.Sprintf(" :bequeath %q", n.Bequeath))
	}
	if n.DefaultCollation != "" {
		sb.WriteString(fmt.Sprintf(" :defaultCollation %q", n.DefaultCollation))
	}
	if n.ContainerMap {
		sb.WriteString(" :containerMap true")
	}
	if n.ContainersDefault {
		sb.WriteString(" :containersDefault true")
	}
	if n.BuildMode != "" {
		sb.WriteString(fmt.Sprintf(" :buildMode %q", n.BuildMode))
	}
	if n.RefreshMode != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMode %q", n.RefreshMode))
	}
	if n.RefreshMethod != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMethod %q", n.RefreshMethod))
	}
	if n.OnPrebuilt {
		sb.WriteString(" :onPrebuilt true")
	}
	if n.NeverRefresh {
		sb.WriteString(" :neverRefresh true")
	}
	if n.WithPK {
		sb.WriteString(" :withPK true")
	}
	if n.WithRowID {
		sb.WriteString(" :withRowID true")
	}
	if n.RefreshOnStmt {
		sb.WriteString(" :refreshOnStmt true")
	}
	if n.StartWith != nil {
		sb.WriteString(" :startWith ")
		writeNode(sb, n.StartWith)
	}
	if n.Next != nil {
		sb.WriteString(" :next ")
		writeNode(sb, n.Next)
	}
	if n.EnableQuery {
		sb.WriteString(" :enableQuery true")
	}
	if n.DisableQuery {
		sb.WriteString(" :disableQuery true")
	}
	if n.EnableOnQueryComputation {
		sb.WriteString(" :enableOnQueryComputation true")
	}
	if n.DisableOnQueryComputation {
		sb.WriteString(" :disableOnQueryComputation true")
	}
	if n.CacheMode != "" {
		sb.WriteString(fmt.Sprintf(" :cacheMode %q", n.CacheMode))
	}
	if n.ParallelMode != "" {
		sb.WriteString(fmt.Sprintf(" :parallelMode %q", n.ParallelMode))
	}
	if n.ParallelDegree != "" {
		sb.WriteString(fmt.Sprintf(" :parallelDegree %q", n.ParallelDegree))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
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
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
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
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
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
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
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
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
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
	if n.Aggregate {
		sb.WriteString(" :aggregate true")
	}
	if n.SqlMacro {
		sb.WriteString(" :sqlMacro true")
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
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Editionable {
		sb.WriteString(" :editionable true")
	}
	if n.NonEditionable {
		sb.WriteString(" :nonEditionable true")
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
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.SetParams != nil {
		sb.WriteString(" :setParams ")
		writeNode(sb, n.SetParams)
	}
	if n.AdviseAction != "" {
		sb.WriteString(fmt.Sprintf(" :adviseAction %q", n.AdviseAction))
	}
	if n.DBLink != "" {
		sb.WriteString(fmt.Sprintf(" :dbLink %q", n.DBLink))
	}
	if n.Feature != "" {
		sb.WriteString(fmt.Sprintf(" :feature %q", n.Feature))
	}
	if n.ParallelDegree != 0 {
		sb.WriteString(fmt.Sprintf(" :parallelDegree %d", n.ParallelDegree))
	}
	if n.Timeout != 0 {
		sb.WriteString(fmt.Sprintf(" :timeout %d", n.Timeout))
	}
	if n.ResumableName != "" {
		sb.WriteString(fmt.Sprintf(" :resumableName %q", n.ResumableName))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterSystemStmt(sb *strings.Builder, n *AlterSystemStmt) {
	sb.WriteString("{ALTERSYSTEM")
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.SetParams != nil {
		sb.WriteString(" :setParams ")
		writeNode(sb, n.SetParams)
	}
	if n.ResetParam != "" {
		sb.WriteString(fmt.Sprintf(" :resetParam %q", n.ResetParam))
	}
	if n.SessionID != "" {
		sb.WriteString(fmt.Sprintf(" :sessionId %q", n.SessionID))
	}
	if n.InstanceID != "" {
		sb.WriteString(fmt.Sprintf(" :instanceId %q", n.InstanceID))
	}
	if n.Immediate {
		sb.WriteString(" :immediate true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.PostTransaction {
		sb.WriteString(" :postTransaction true")
	}
	if n.NoReplay {
		sb.WriteString(" :noReplay true")
	}
	if n.Timeout != 0 {
		sb.WriteString(fmt.Sprintf(" :timeout %d", n.Timeout))
	}
	if n.FlushTarget != "" {
		sb.WriteString(fmt.Sprintf(" :flushTarget %q", n.FlushTarget))
	}
	if n.FlushScope != "" {
		sb.WriteString(fmt.Sprintf(" :flushScope %q", n.FlushScope))
	}
	if n.FlushRedoDB != "" {
		sb.WriteString(fmt.Sprintf(" :flushRedoDb %q", n.FlushRedoDB))
	}
	if n.FlushRedoConfirm != "" {
		sb.WriteString(fmt.Sprintf(" :flushRedoConfirm %q", n.FlushRedoConfirm))
	}
	if n.ArchiveLogSpec != "" {
		sb.WriteString(fmt.Sprintf(" :archiveLogSpec %q", n.ArchiveLogSpec))
	}
	if n.ArchiveLogValue != "" {
		sb.WriteString(fmt.Sprintf(" :archiveLogValue %q", n.ArchiveLogValue))
	}
	if n.ArchiveInstance != "" {
		sb.WriteString(fmt.Sprintf(" :archiveInstance %q", n.ArchiveInstance))
	}
	if n.ArchiveThread != 0 {
		sb.WriteString(fmt.Sprintf(" :archiveThread %d", n.ArchiveThread))
	}
	if n.ArchiveNoSwitch {
		sb.WriteString(" :archiveNoSwitch true")
	}
	if n.ArchiveBackupCF {
		sb.WriteString(" :archiveBackupCF true")
	}
	if n.ArchiveTo != "" {
		sb.WriteString(fmt.Sprintf(" :archiveTo %q", n.ArchiveTo))
	}
	if n.CheckScope != "" {
		sb.WriteString(fmt.Sprintf(" :checkScope %q", n.CheckScope))
	}
	if n.Feature != "" {
		sb.WriteString(fmt.Sprintf(" :feature %q", n.Feature))
	}
	if n.ShutdownDisp != "" {
		sb.WriteString(fmt.Sprintf(" :shutdownDisp %q", n.ShutdownDisp))
	}
	if n.RelocateClient != "" {
		sb.WriteString(fmt.Sprintf(" :relocateClient %q", n.RelocateClient))
	}
	if n.RollingVersion != "" {
		sb.WriteString(fmt.Sprintf(" :rollingVersion %q", n.RollingVersion))
	}
	if n.EncryptionAction != "" {
		sb.WriteString(fmt.Sprintf(" :encryptionAction %q", n.EncryptionAction))
	}
	if n.Comment != "" {
		sb.WriteString(fmt.Sprintf(" :comment %q", n.Comment))
	}
	if n.Deferred {
		sb.WriteString(" :deferred true")
	}
	if n.Scope != "" {
		sb.WriteString(fmt.Sprintf(" :scope %q", n.Scope))
	}
	if n.SID != "" {
		sb.WriteString(fmt.Sprintf(" :sid %q", n.SID))
	}
	if n.Container != "" {
		sb.WriteString(fmt.Sprintf(" :container %q", n.Container))
	}
	if n.SqlID != "" {
		sb.WriteString(fmt.Sprintf(" :sqlId %q", n.SqlID))
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
	if n.WriteWait != "" {
		sb.WriteString(fmt.Sprintf(" :writeWait %q", n.WriteWait))
	}
	if n.WriteBatch != "" {
		sb.WriteString(fmt.Sprintf(" :writeBatch %q", n.WriteBatch))
	}
	if n.Force != "" {
		sb.WriteString(fmt.Sprintf(" :force %q", n.Force))
	}
	if n.HasForceInt {
		sb.WriteString(fmt.Sprintf(" :forceInteger %d", n.ForceInteger))
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
	if n.UseRollbackSegment != "" {
		sb.WriteString(fmt.Sprintf(" :useRollbackSegment %q", n.UseRollbackSegment))
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
	if n.DeleteSystem {
		sb.WriteString(" :deleteSystem true")
	}
	if n.SampleValue != 0 {
		sb.WriteString(fmt.Sprintf(" :sampleValue %d", n.SampleValue))
	}
	if n.SampleUnit != "" {
		sb.WriteString(fmt.Sprintf(" :sampleUnit %q", n.SampleUnit))
	}
	if n.CascadeFast {
		sb.WriteString(" :cascadeFast true")
	}
	if n.IntoTable != nil {
		sb.WriteString(" :intoTable ")
		writeNode(sb, n.IntoTable)
	}
	if n.SetDanglingNull {
		sb.WriteString(" :setDanglingNull true")
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Offline {
		sb.WriteString(" :offline true")
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
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables (")
		for i, t := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
		sb.WriteString(")")
	} else if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Before {
		sb.WriteString(" :before true")
	}
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
	if n.ToBeforeDrop {
		sb.WriteString(" :toBeforeDrop true")
	}
	if n.Rename != "" {
		sb.WriteString(fmt.Sprintf(" :rename %q", n.Rename))
	}
	if n.EnableTriggers != nil {
		if *n.EnableTriggers {
			sb.WriteString(" :triggers \"ENABLE\"")
		} else {
			sb.WriteString(" :triggers \"DISABLE\"")
		}
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeFlashbackDatabaseStmt(sb *strings.Builder, n *FlashbackDatabaseStmt) {
	sb.WriteString("{FLASHBACKDATABASE")
	if n.Modifier != "" {
		sb.WriteString(fmt.Sprintf(" :modifier %q", n.Modifier))
	}
	if n.DatabaseName != nil {
		sb.WriteString(" :databaseName ")
		writeNode(sb, n.DatabaseName)
	}
	if n.Before {
		sb.WriteString(" :before true")
	}
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
	if n.ToResetlogs {
		sb.WriteString(" :toResetlogs true")
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

func writeLockTableItem(sb *strings.Builder, n *LockTableItem) {
	sb.WriteString("{LOCK_TABLE_ITEM")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.PartitionType != "" {
		sb.WriteString(fmt.Sprintf(" :partitionType %q", n.PartitionType))
		if n.PartitionFor {
			sb.WriteString(" :partitionFor true")
		}
		if n.PartitionName != "" {
			sb.WriteString(fmt.Sprintf(" :partitionName %q", n.PartitionName))
		}
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeLockTableStmt(sb *strings.Builder, n *LockTableStmt) {
	sb.WriteString("{LOCK_TABLE")
	if len(n.Tables) > 0 {
		sb.WriteString(" :tables (")
		for i, item := range n.Tables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeLockTableItem(sb, item)
		}
		sb.WriteString(")")
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

func writeRoleIdentified(sb *strings.Builder, idType RoleIdentifiedType, identifyBy, identifySchema string, hasIdentified bool) {
	if !hasIdentified && idType == ROLE_NOT_IDENTIFIED && identifyBy == "" {
		return
	}
	switch idType {
	case ROLE_NOT_IDENTIFIED:
		if hasIdentified {
			sb.WriteString(" :identified NOT_IDENTIFIED")
		}
	case ROLE_IDENTIFIED_BY:
		sb.WriteString(fmt.Sprintf(" :identified BY :password %q", identifyBy))
	case ROLE_IDENTIFIED_USING:
		if identifySchema != "" {
			sb.WriteString(fmt.Sprintf(" :identified USING :schema %q :package %q", identifySchema, identifyBy))
		} else {
			sb.WriteString(fmt.Sprintf(" :identified USING :package %q", identifyBy))
		}
	case ROLE_IDENTIFIED_EXTERNALLY:
		sb.WriteString(" :identified EXTERNALLY")
	case ROLE_IDENTIFIED_GLOBALLY:
		sb.WriteString(" :identified GLOBALLY")
		if identifyBy != "" {
			sb.WriteString(fmt.Sprintf(" :as %q", identifyBy))
		}
	}
}

func writeCreateRoleStmt(sb *strings.Builder, n *CreateRoleStmt) {
	sb.WriteString("{CREATE_ROLE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	writeRoleIdentified(sb, n.IdentifiedType, n.IdentifyBy, n.IdentifySchema, n.HasIdentified)
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

func writeAlterRoleStmt(sb *strings.Builder, n *AlterRoleStmt) {
	sb.WriteString("{ALTER_ROLE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	writeRoleIdentified(sb, n.IdentifiedType, n.IdentifyBy, n.IdentifySchema, true)
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

func writeAlterProfileStmt(sb *strings.Builder, n *AlterProfileStmt) {
	sb.WriteString("{ALTER_PROFILE")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
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

func writeAlterResourceCostStmt(sb *strings.Builder, n *AlterResourceCostStmt) {
	sb.WriteString("{ALTER_RESOURCE_COST")
	for _, c := range n.Costs {
		sb.WriteString(fmt.Sprintf(" :cost %q %q", c.Name, c.Value))
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
	if n.PctFree != 0 {
		sb.WriteString(fmt.Sprintf(" :pctfree %d", n.PctFree))
	}
	if n.PctUsed != 0 {
		sb.WriteString(fmt.Sprintf(" :pctused %d", n.PctUsed))
	}
	if n.IniTrans != 0 {
		sb.WriteString(fmt.Sprintf(" :initrans %d", n.IniTrans))
	}
	if n.AllocateSize != "" {
		sb.WriteString(fmt.Sprintf(" :allocateSize %q", n.AllocateSize))
	}
	if n.AllocateDatafile != "" {
		sb.WriteString(fmt.Sprintf(" :allocateDatafile %q", n.AllocateDatafile))
	}
	if n.AllocateInstance != 0 {
		sb.WriteString(fmt.Sprintf(" :allocateInstance %d", n.AllocateInstance))
	}
	if n.DeallocateKeep != "" {
		sb.WriteString(fmt.Sprintf(" :deallocateKeep %q", n.DeallocateKeep))
	}
	if n.EditionName != "" {
		sb.WriteString(fmt.Sprintf(" :editionName %q", n.EditionName))
	}
	if n.ScopeColumn != "" {
		sb.WriteString(fmt.Sprintf(" :scopeColumn %q", n.ScopeColumn))
	}
	if n.ScopeTable != nil {
		sb.WriteString(" :scopeTable ")
		writeNode(sb, n.ScopeTable)
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

func writeCreateAuditPolicyStmt(sb *strings.Builder, n *CreateAuditPolicyStmt) {
	sb.WriteString("{CREATE_AUDIT_POLICY")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if len(n.Privileges) > 0 {
		sb.WriteString(" :privileges (")
		for i, p := range n.Privileges {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", p))
		}
		sb.WriteString(")")
	}
	if len(n.Actions) > 0 {
		sb.WriteString(" :actions (")
		for i, a := range n.Actions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeAuditActionEntry(sb, a)
		}
		sb.WriteString(")")
	}
	if len(n.ComponentActions) > 0 {
		sb.WriteString(" :component_actions (")
		for i, c := range n.ComponentActions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeAuditComponentAction(sb, c)
		}
		sb.WriteString(")")
	}
	if len(n.Roles) > 0 {
		sb.WriteString(" :roles (")
		for i, r := range n.Roles {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", r))
		}
		sb.WriteString(")")
	}
	if n.WhenCondition != "" {
		sb.WriteString(fmt.Sprintf(" :when %q", n.WhenCondition))
	}
	if n.EvaluatePer != "" {
		sb.WriteString(fmt.Sprintf(" :evaluate_per %q", n.EvaluatePer))
	}
	if n.OnlyToplevel {
		sb.WriteString(" :only_toplevel true")
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

func writeAuditActionEntry(sb *strings.Builder, n *AuditActionEntry) {
	sb.WriteString("{ACTION")
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	if len(n.Columns) > 0 {
		sb.WriteString(" :columns (")
		for i, c := range n.Columns {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", c))
		}
		sb.WriteString(")")
	}
	if n.Object != nil {
		sb.WriteString(" :object ")
		writeNode(sb, n.Object)
	}
	if n.Directory != "" {
		sb.WriteString(fmt.Sprintf(" :directory %q", n.Directory))
	}
	if n.MiningModel {
		sb.WriteString(" :mining_model true")
	}
	sb.WriteString("}")
}

func writeAuditComponentAction(sb *strings.Builder, n *AuditComponentAction) {
	sb.WriteString("{COMPONENT")
	sb.WriteString(fmt.Sprintf(" :component %q", n.Component))
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
	if n.Object != "" {
		sb.WriteString(fmt.Sprintf(" :object %q", n.Object))
	}
	sb.WriteString("}")
}

func writeAlterAuditPolicyStmt(sb *strings.Builder, n *AlterAuditPolicyStmt) {
	sb.WriteString("{ALTER_AUDIT_POLICY")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.AddDrop != "" {
		sb.WriteString(fmt.Sprintf(" :add_drop %q", n.AddDrop))
	}
	if len(n.Privileges) > 0 {
		sb.WriteString(" :privileges (")
		for i, p := range n.Privileges {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", p))
		}
		sb.WriteString(")")
	}
	if len(n.Actions) > 0 {
		sb.WriteString(" :actions (")
		for i, a := range n.Actions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeAuditActionEntry(sb, a)
		}
		sb.WriteString(")")
	}
	if len(n.ComponentActions) > 0 {
		sb.WriteString(" :component_actions (")
		for i, c := range n.ComponentActions {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeAuditComponentAction(sb, c)
		}
		sb.WriteString(")")
	}
	if len(n.Roles) > 0 {
		sb.WriteString(" :roles (")
		for i, r := range n.Roles {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", r))
		}
		sb.WriteString(")")
	}
	if n.ConditionDrop {
		sb.WriteString(" :condition_drop true")
	}
	if n.Condition != "" {
		sb.WriteString(fmt.Sprintf(" :condition %q", n.Condition))
	}
	if n.EvaluatePer != "" {
		sb.WriteString(fmt.Sprintf(" :evaluate_per %q", n.EvaluatePer))
	}
	if n.AddToplevel {
		sb.WriteString(" :add_toplevel true")
	}
	if n.DropToplevel {
		sb.WriteString(" :drop_toplevel true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDropAuditPolicyStmt(sb *strings.Builder, n *DropAuditPolicyStmt) {
	sb.WriteString("{DROP_AUDIT_POLICY")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAuditStmt(sb *strings.Builder, n *AuditStmt) {
	sb.WriteString("{AUDIT")
	if n.Policy != "" {
		sb.WriteString(fmt.Sprintf(" :policy %q", n.Policy))
	}
	if len(n.ByUsers) > 0 {
		sb.WriteString(" :by_users (")
		for i, u := range n.ByUsers {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", u))
		}
		sb.WriteString(")")
	}
	if len(n.ExceptUsers) > 0 {
		sb.WriteString(" :except_users (")
		for i, u := range n.ExceptUsers {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", u))
		}
		sb.WriteString(")")
	}
	if len(n.WithRoles) > 0 {
		if n.WithRoleExcept {
			sb.WriteString(" :except_with_roles (")
		} else {
			sb.WriteString(" :by_with_roles (")
		}
		for i, r := range n.WithRoles {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", r))
		}
		sb.WriteString(")")
	}
	if n.ContextNS != "" {
		sb.WriteString(fmt.Sprintf(" :context_ns %q", n.ContextNS))
	}
	if len(n.ContextAttrs) > 0 {
		sb.WriteString(" :context_attrs (")
		for i, a := range n.ContextAttrs {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
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
	if n.OnDefault {
		sb.WriteString(" :on_default true")
	}
	if n.OnDirectory != "" {
		sb.WriteString(fmt.Sprintf(" :on_directory %q", n.OnDirectory))
	}
	if n.OnNetwork {
		sb.WriteString(" :on_network true")
	}
	if n.OnDirectPath {
		sb.WriteString(" :on_direct_path true")
	}
	if n.By != "" {
		sb.WriteString(fmt.Sprintf(" :by %q", n.By))
	}
	if len(n.ByUsers2) > 0 {
		sb.WriteString(" :by_users2 (")
		for i, u := range n.ByUsers2 {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", u))
		}
		sb.WriteString(")")
	}
	if n.When != "" {
		sb.WriteString(fmt.Sprintf(" :when %q", n.When))
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

func writeNoauditStmt(sb *strings.Builder, n *NoauditStmt) {
	sb.WriteString("{NOAUDIT")
	if n.Policy != "" {
		sb.WriteString(fmt.Sprintf(" :policy %q", n.Policy))
	}
	if len(n.ByUsers) > 0 {
		sb.WriteString(" :by_users (")
		for i, u := range n.ByUsers {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", u))
		}
		sb.WriteString(")")
	}
	if len(n.WithRoles) > 0 {
		sb.WriteString(" :with_roles (")
		for i, r := range n.WithRoles {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", r))
		}
		sb.WriteString(")")
	}
	if n.ContextNS != "" {
		sb.WriteString(fmt.Sprintf(" :context_ns %q", n.ContextNS))
	}
	if len(n.ContextAttrs) > 0 {
		sb.WriteString(" :context_attrs (")
		for i, a := range n.ContextAttrs {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
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
	if n.OnDefault {
		sb.WriteString(" :on_default true")
	}
	if n.OnDirectory != "" {
		sb.WriteString(fmt.Sprintf(" :on_directory %q", n.OnDirectory))
	}
	if n.OnNetwork {
		sb.WriteString(" :on_network true")
	}
	if n.OnDirectPath {
		sb.WriteString(" :on_direct_path true")
	}
	if len(n.ByUsers2) > 0 {
		sb.WriteString(" :by_users2 (")
		for i, u := range n.ByUsers2 {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", u))
		}
		sb.WriteString(")")
	}
	if n.When != "" {
		sb.WriteString(fmt.Sprintf(" :when %q", n.When))
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
	if n.Local {
		sb.WriteString(" :local true")
	}
	if n.Undo {
		sb.WriteString(" :undo true")
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
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
	if n.EncryptionAlgorithm != "" {
		sb.WriteString(fmt.Sprintf(" :encryption_algorithm %q", n.EncryptionAlgorithm))
	}
	if n.Compress != "" {
		sb.WriteString(fmt.Sprintf(" :compress %q", n.Compress))
	}
	if n.MaxSize != "" {
		sb.WriteString(fmt.Sprintf(" :maxsize %q", n.MaxSize))
	}
	if n.MinimumExtent != "" {
		sb.WriteString(fmt.Sprintf(" :minimum_extent %q", n.MinimumExtent))
	}
	if n.TablespaceGroup != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace_group %q", n.TablespaceGroup))
	}
	if n.Flashback != "" {
		sb.WriteString(fmt.Sprintf(" :flashback %q", n.Flashback))
	}
	if n.LostWriteProtection != "" {
		sb.WriteString(fmt.Sprintf(" :lost_write_protection %q", n.LostWriteProtection))
	}
	if n.DefaultParams != "" {
		sb.WriteString(fmt.Sprintf(" :default_params %q", n.DefaultParams))
	}
	if n.ForLeaf != "" {
		sb.WriteString(fmt.Sprintf(" :for_leaf %q", n.ForLeaf))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterTablespaceStmt(sb *strings.Builder, n *AlterTablespaceStmt) {
	sb.WriteString("{ALTER_TABLESPACE")
	if n.IsSet {
		sb.WriteString(" :is_set true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.DefaultParams != "" {
		sb.WriteString(fmt.Sprintf(" :default_params %q", n.DefaultParams))
	}
	if n.MinimumExtent != "" {
		sb.WriteString(fmt.Sprintf(" :minimum_extent %q", n.MinimumExtent))
	}
	if n.Resize != "" {
		sb.WriteString(fmt.Sprintf(" :resize %q", n.Resize))
	}
	if n.Coalesce {
		sb.WriteString(" :coalesce true")
	}
	if n.ShrinkSpace {
		sb.WriteString(" :shrink_space true")
	}
	if n.ShrinkKeep != "" {
		sb.WriteString(fmt.Sprintf(" :shrink_keep %q", n.ShrinkKeep))
	}
	if n.ShrinkTempfile != "" {
		sb.WriteString(fmt.Sprintf(" :shrink_tempfile %q", n.ShrinkTempfile))
	}
	if n.ShrinkTempfileKeep != "" {
		sb.WriteString(fmt.Sprintf(" :shrink_tempfile_keep %q", n.ShrinkTempfileKeep))
	}
	if n.RenameTo != "" {
		sb.WriteString(fmt.Sprintf(" :rename_to %q", n.RenameTo))
	}
	if n.BeginBackup {
		sb.WriteString(" :begin_backup true")
	}
	if n.EndBackup {
		sb.WriteString(" :end_backup true")
	}
	if n.AddDatafile || n.AddTempfile {
		if n.AddDatafile {
			sb.WriteString(" :add_datafile true")
		}
		if n.AddTempfile {
			sb.WriteString(" :add_tempfile true")
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
	}
	if n.DropDatafile {
		sb.WriteString(fmt.Sprintf(" :drop_datafile %q", n.DropFileRef))
	}
	if n.DropTempfile {
		sb.WriteString(fmt.Sprintf(" :drop_tempfile %q", n.DropFileRef))
	}
	if n.RenameDatafile {
		sb.WriteString(" :rename_datafile true")
	}
	if n.DatafileOnline {
		sb.WriteString(" :datafile_online true")
	}
	if n.DatafileOffline {
		sb.WriteString(" :datafile_offline true")
	}
	if n.TempfileOnline {
		sb.WriteString(" :tempfile_online true")
	}
	if n.TempfileOffline {
		sb.WriteString(" :tempfile_offline true")
	}
	if n.Logging != "" {
		sb.WriteString(fmt.Sprintf(" :logging %q", n.Logging))
	}
	if n.ForceLogging != "" {
		sb.WriteString(fmt.Sprintf(" :force_logging %q", n.ForceLogging))
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Offline {
		sb.WriteString(" :offline true")
	}
	if n.OfflineMode != "" {
		sb.WriteString(fmt.Sprintf(" :offline_mode %q", n.OfflineMode))
	}
	if n.ReadOnly {
		sb.WriteString(" :read_only true")
	}
	if n.ReadWrite {
		sb.WriteString(" :read_write true")
	}
	if n.Permanent {
		sb.WriteString(" :permanent true")
	}
	if n.TempState {
		sb.WriteString(" :temporary true")
	}
	if n.Autoextend != nil {
		sb.WriteString(" :autoextend ")
		writeNode(sb, n.Autoextend)
	}
	if n.Flashback != "" {
		sb.WriteString(fmt.Sprintf(" :flashback %q", n.Flashback))
	}
	if n.Retention != "" {
		sb.WriteString(fmt.Sprintf(" :retention %q", n.Retention))
	}
	if n.Encryption != "" {
		sb.WriteString(fmt.Sprintf(" :encryption %q", n.Encryption))
	}
	if n.LostWriteProtection != "" {
		sb.WriteString(fmt.Sprintf(" :lost_write_protection %q", n.LostWriteProtection))
	}
	if n.TablespaceGroup != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace_group %q", n.TablespaceGroup))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateTablespaceSetStmt(sb *strings.Builder, n *CreateTablespaceSetStmt) {
	sb.WriteString("{CREATE_TABLESPACE_SET")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Shardspace != "" {
		sb.WriteString(fmt.Sprintf(" :shardspace %q", n.Shardspace))
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
	if n.Logging != "" {
		sb.WriteString(fmt.Sprintf(" :logging %q", n.Logging))
	}
	if n.Encryption != "" {
		sb.WriteString(fmt.Sprintf(" :encryption %q", n.Encryption))
	}
	if n.DefaultParams != "" {
		sb.WriteString(fmt.Sprintf(" :default_params %q", n.DefaultParams))
	}
	if n.Extent != "" {
		sb.WriteString(fmt.Sprintf(" :extent %q", n.Extent))
	}
	if n.Segment != "" {
		sb.WriteString(fmt.Sprintf(" :segment %q", n.Segment))
	}
	if n.Flashback != "" {
		sb.WriteString(fmt.Sprintf(" :flashback %q", n.Flashback))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDropTablespaceStmt(sb *strings.Builder, n *DropTablespaceStmt) {
	sb.WriteString("{DROP_TABLESPACE")
	if n.IsSet {
		sb.WriteString(" :is_set true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.IfExists {
		sb.WriteString(" :if_exists true")
	}
	if n.DropQuota {
		sb.WriteString(" :drop_quota true")
	}
	if n.KeepQuota {
		sb.WriteString(" :keep_quota true")
	}
	if n.IncludingContents {
		sb.WriteString(" :including_contents true")
	}
	if n.AndDatafiles {
		sb.WriteString(" :and_datafiles true")
	}
	if n.KeepDatafiles {
		sb.WriteString(" :keep_datafiles true")
	}
	if n.CascadeConstraints {
		sb.WriteString(" :cascade_constraints true")
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

func writeAlterClusterStmt(sb *strings.Builder, n *AlterClusterStmt) {
	sb.WriteString("{ALTER_CLUSTER")
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
	if n.ModifyPartition != "" {
		sb.WriteString(fmt.Sprintf(" :modifyPartition %q", n.ModifyPartition))
	}
	if n.Parallel != "" {
		sb.WriteString(fmt.Sprintf(" :parallel %q", n.Parallel))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterDimensionStmt(sb *strings.Builder, n *AlterDimensionStmt) {
	sb.WriteString("{ALTER_DIMENSION")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Compile {
		sb.WriteString(" :compile true")
	}
	if len(n.AddLevels) > 0 {
		sb.WriteString(" :addLevels (")
		for i, l := range n.AddLevels {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, l)
		}
		sb.WriteString(")")
	}
	if len(n.AddHierarchies) > 0 {
		sb.WriteString(" :addHierarchies (")
		for i, h := range n.AddHierarchies {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, h)
		}
		sb.WriteString(")")
	}
	if len(n.AddAttributes) > 0 {
		sb.WriteString(" :addAttributes (")
		for i, a := range n.AddAttributes {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString(")")
	}
	if len(n.DropLevels) > 0 {
		sb.WriteString(" :dropLevels (")
		for i, l := range n.DropLevels {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", l))
		}
		sb.WriteString(")")
	}
	if len(n.DropHierarchies) > 0 {
		sb.WriteString(" :dropHierarchies (")
		for i, h := range n.DropHierarchies {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", h))
		}
		sb.WriteString(")")
	}
	if len(n.DropAttributes) > 0 {
		sb.WriteString(" :dropAttributes (")
		for i, a := range n.DropAttributes {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateMaterializedZonemapStmt(sb *strings.Builder, n *CreateMaterializedZonemapStmt) {
	sb.WriteString("{CREATE_MATERIALIZED_ZONEMAP")
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Tablespace != "" {
		sb.WriteString(fmt.Sprintf(" :tablespace %q", n.Tablespace))
	}
	if n.Scale != nil {
		sb.WriteString(fmt.Sprintf(" :scale %d", *n.Scale))
	}
	if n.PctFree != nil {
		sb.WriteString(fmt.Sprintf(" :pctfree %d", *n.PctFree))
	}
	if n.PctUsed != nil {
		sb.WriteString(fmt.Sprintf(" :pctused %d", *n.PctUsed))
	}
	if n.Cache {
		sb.WriteString(" :cache true")
	}
	if n.NoCache {
		sb.WriteString(" :nocache true")
	}
	if n.RefreshMethod != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMethod %q", n.RefreshMethod))
	}
	if n.RefreshOn != "" {
		sb.WriteString(fmt.Sprintf(" :refreshOn %q", n.RefreshOn))
	}
	if n.EnablePruning {
		sb.WriteString(" :enablePruning true")
	}
	if n.DisablePruning {
		sb.WriteString(" :disablePruning true")
	}
	if n.OnTable != nil {
		sb.WriteString(" :onTable ")
		writeNode(sb, n.OnTable)
	}
	if len(n.OnColumns) > 0 {
		sb.WriteString(" :onColumns (")
		for i, c := range n.OnColumns {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", c))
		}
		sb.WriteString(")")
	}
	if len(n.ColumnAliases) > 0 {
		sb.WriteString(" :columnAliases (")
		for i, a := range n.ColumnAliases {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
	if n.AsQuery != nil {
		sb.WriteString(" :asQuery ")
		writeNode(sb, n.AsQuery)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterMaterializedZonemapStmt(sb *strings.Builder, n *AlterMaterializedZonemapStmt) {
	sb.WriteString("{ALTER_MATERIALIZED_ZONEMAP")
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
	if n.PctFree != nil {
		sb.WriteString(fmt.Sprintf(" :pctfree %d", *n.PctFree))
	}
	if n.PctUsed != nil {
		sb.WriteString(fmt.Sprintf(" :pctused %d", *n.PctUsed))
	}
	if n.Cache {
		sb.WriteString(" :cache true")
	}
	if n.NoCache {
		sb.WriteString(" :nocache true")
	}
	if n.RefreshMethod != "" {
		sb.WriteString(fmt.Sprintf(" :refreshMethod %q", n.RefreshMethod))
	}
	if n.RefreshOn != "" {
		sb.WriteString(fmt.Sprintf(" :refreshOn %q", n.RefreshOn))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateInmemoryJoinGroupStmt(sb *strings.Builder, n *CreateInmemoryJoinGroupStmt) {
	sb.WriteString("{CREATE_INMEMORY_JOIN_GROUP")
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Members) > 0 {
		sb.WriteString(" :members (")
		for i, m := range n.Members {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, m)
		}
		sb.WriteString(")")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterInmemoryJoinGroupStmt(sb *strings.Builder, n *AlterInmemoryJoinGroupStmt) {
	sb.WriteString("{ALTER_INMEMORY_JOIN_GROUP")
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
	if n.Member != nil {
		sb.WriteString(" :member ")
		writeNode(sb, n.Member)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeJoinGroupMember(sb *strings.Builder, n *JoinGroupMember) {
	sb.WriteString("{JOIN_GROUP_MEMBER")
	if n.Table != nil {
		sb.WriteString(" :table ")
		writeNode(sb, n.Table)
	}
	if n.Column != "" {
		sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
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
	if n.PctFree != "" {
		sb.WriteString(fmt.Sprintf(" :pctfree %q", n.PctFree))
	}
	if n.PctUsed != "" {
		sb.WriteString(fmt.Sprintf(" :pctused %q", n.PctUsed))
	}
	if n.InitTrans != "" {
		sb.WriteString(fmt.Sprintf(" :initrans %q", n.InitTrans))
	}
	if n.MaxTrans != "" {
		sb.WriteString(fmt.Sprintf(" :maxtrans %q", n.MaxTrans))
	}
	if n.Parameters != "" {
		sb.WriteString(fmt.Sprintf(" :parameters %q", n.Parameters))
	}
	if n.Invalidation != "" {
		sb.WriteString(fmt.Sprintf(" :invalidation %q", n.Invalidation))
	}
	if n.DeallocateKeep != "" {
		sb.WriteString(fmt.Sprintf(" :deallocateKeep %q", n.DeallocateKeep))
	}
	if n.ModifyDefaultFor != "" {
		sb.WriteString(fmt.Sprintf(" :modifyDefaultFor %q", n.ModifyDefaultFor))
	}
	if n.AddPartitionName != "" {
		sb.WriteString(fmt.Sprintf(" :addPartitionName %q", n.AddPartitionName))
	}
	if n.ModifyPartAction != "" {
		sb.WriteString(fmt.Sprintf(" :modifyPartAction %q", n.ModifyPartAction))
	}
	if n.SplitPartition != "" {
		sb.WriteString(fmt.Sprintf(" :splitPartition %q", n.SplitPartition))
	}
	if n.SplitValues != nil {
		sb.WriteString(" :splitValues ")
		writeNode(sb, n.SplitValues)
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
	if n.Annotations != nil {
		sb.WriteString(" :annotations ")
		writeNode(sb, n.Annotations)
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
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
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

func writeCreateIndextypeStmt(sb *strings.Builder, n *CreateIndextypeStmt) {
	sb.WriteString("{CREATE_INDEXTYPE")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	if len(n.Operators) > 0 {
		sb.WriteString(" :operators (")
		for i, op := range n.Operators {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, op)
		}
		sb.WriteString(")")
	}
	if n.UsingType != nil {
		sb.WriteString(" :usingType ")
		writeNode(sb, n.UsingType)
	}
	if n.WithLocal {
		sb.WriteString(" :withLocal true")
	}
	if n.WithRange {
		sb.WriteString(" :withRange true")
	}
	if n.StorageTable != "" {
		sb.WriteString(fmt.Sprintf(" :storageTable %q", n.StorageTable))
	}
	if n.ArrayDML {
		sb.WriteString(" :arrayDML true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIndextypeOp(sb *strings.Builder, n *IndextypeOp) {
	sb.WriteString("{INDEXTYPE_OP")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.ParamTypes) > 0 {
		sb.WriteString(fmt.Sprintf(" :paramTypes %q", strings.Join(n.ParamTypes, ",")))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterIndextypeStmt(sb *strings.Builder, n *AlterIndextypeStmt) {
	sb.WriteString("{ALTER_INDEXTYPE")
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	if len(n.Modifications) > 0 {
		sb.WriteString(" :modifications (")
		for i, m := range n.Modifications {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, m)
		}
		sb.WriteString(")")
	}
	if n.UsingType != nil {
		sb.WriteString(" :usingType ")
		writeNode(sb, n.UsingType)
	}
	if n.ArrayDML {
		sb.WriteString(" :arrayDML true")
	}
	if n.WithLocal {
		sb.WriteString(" :withLocal true")
	}
	if n.StorageTable != "" {
		sb.WriteString(fmt.Sprintf(" :storageTable %q", n.StorageTable))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeIndextypeModOp(sb *strings.Builder, n *IndextypeModOp) {
	sb.WriteString("{INDEXTYPE_MOD_OP")
	if n.Add {
		sb.WriteString(" :add true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.ParamTypes) > 0 {
		sb.WriteString(fmt.Sprintf(" :paramTypes %q", strings.Join(n.ParamTypes, ",")))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateOperatorStmt(sb *strings.Builder, n *CreateOperatorStmt) {
	sb.WriteString("{CREATE_OPERATOR")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.Bindings) > 0 {
		sb.WriteString(" :bindings (")
		for i, b := range n.Bindings {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, b)
		}
		sb.WriteString(")")
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeOperatorBinding(sb *strings.Builder, n *OperatorBinding) {
	sb.WriteString("{OP_BINDING")
	if len(n.ParamTypes) > 0 {
		sb.WriteString(fmt.Sprintf(" :paramTypes %q", strings.Join(n.ParamTypes, ",")))
	}
	if n.ReturnType != "" {
		sb.WriteString(fmt.Sprintf(" :returnType %q", n.ReturnType))
	}
	if n.UsingFunc != nil {
		sb.WriteString(" :usingFunc ")
		writeNode(sb, n.UsingFunc)
	}
	if n.AncillaryTo != nil {
		sb.WriteString(" :ancillaryTo ")
		writeNode(sb, n.AncillaryTo)
	}
	if len(n.AncillaryParams) > 0 {
		sb.WriteString(fmt.Sprintf(" :ancillaryParams %q", strings.Join(n.AncillaryParams, ",")))
	}
	if n.WithIndexCtx {
		sb.WriteString(" :withIndexCtx true")
	}
	if n.ScanCtxType != "" {
		sb.WriteString(fmt.Sprintf(" :scanCtxType %q", n.ScanCtxType))
	}
	if n.ComputeAnc {
		sb.WriteString(" :computeAnc true")
	}
	if n.WithColumnCtx {
		sb.WriteString(" :withColumnCtx true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterOperatorStmt(sb *strings.Builder, n *AlterOperatorStmt) {
	sb.WriteString("{ALTER_OPERATOR")
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	if n.Binding != nil {
		sb.WriteString(" :binding ")
		writeNode(sb, n.Binding)
	}
	if len(n.DropTypes) > 0 {
		sb.WriteString(fmt.Sprintf(" :dropTypes %q", strings.Join(n.DropTypes, ",")))
	}
	if n.DropForce {
		sb.WriteString(" :dropForce true")
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

func writeCreateMviewLogStmt(sb *strings.Builder, n *CreateMviewLogStmt) {
	sb.WriteString("{CREATE_MVIEW_LOG")
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.OnTable != nil {
		sb.WriteString(" :onTable ")
		writeNode(sb, n.OnTable)
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	if n.WithPK {
		sb.WriteString(" :withPK true")
	}
	if n.WithRowID {
		sb.WriteString(" :withRowID true")
	}
	if n.WithOID {
		sb.WriteString(" :withOID true")
	}
	if n.WithSeq {
		sb.WriteString(" :withSeq true")
	}
	if n.CommitSCN {
		sb.WriteString(" :commitSCN true")
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Including {
		sb.WriteString(" :including true")
	}
	if n.Excluding {
		sb.WriteString(" :excluding true")
	}
	if n.PurgeMode != "" {
		sb.WriteString(fmt.Sprintf(" :purgeMode %q", n.PurgeMode))
	}
	if n.PurgeStart != nil {
		sb.WriteString(" :purgeStart ")
		writeNode(sb, n.PurgeStart)
	}
	if n.PurgeNext != nil {
		sb.WriteString(" :purgeNext ")
		writeNode(sb, n.PurgeNext)
	}
	if n.ForRefresh != "" {
		sb.WriteString(fmt.Sprintf(" :forRefresh %q", n.ForRefresh))
	}
	if n.StagingLog != "" {
		sb.WriteString(fmt.Sprintf(" :stagingLog %q", n.StagingLog))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterMviewLogStmt(sb *strings.Builder, n *AlterMviewLogStmt) {
	sb.WriteString("{ALTER_MVIEW_LOG")
	if n.IfExists {
		sb.WriteString(" :ifExists true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.OnTable != nil {
		sb.WriteString(" :onTable ")
		writeNode(sb, n.OnTable)
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateAnalyticViewStmt(sb *strings.Builder, n *CreateAnalyticViewStmt) {
	sb.WriteString("{CREATE_ANALYTIC_VIEW")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.NoForce {
		sb.WriteString(" :noForce true")
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.UsingTable != nil {
		sb.WriteString(" :usingTable ")
		writeNode(sb, n.UsingTable)
	}
	if n.UsingAlias != "" {
		sb.WriteString(fmt.Sprintf(" :usingAlias %q", n.UsingAlias))
	}
	if n.DimBy != nil {
		sb.WriteString(" :dimBy ")
		writeNode(sb, n.DimBy)
	}
	if n.Measures != nil {
		sb.WriteString(" :measures ")
		writeNode(sb, n.Measures)
	}
	if n.DefaultMeasure != "" {
		sb.WriteString(fmt.Sprintf(" :defaultMeasure %q", n.DefaultMeasure))
	}
	if n.DefaultAggregate != "" {
		sb.WriteString(fmt.Sprintf(" :defaultAggregate %q", n.DefaultAggregate))
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterAnalyticViewStmt(sb *strings.Builder, n *AlterAnalyticViewStmt) {
	sb.WriteString("{ALTER_ANALYTIC_VIEW")
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
	if n.NewName != nil {
		sb.WriteString(" :newName ")
		writeNode(sb, n.NewName)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateJsonDualityViewStmt(sb *strings.Builder, n *CreateJsonDualityViewStmt) {
	sb.WriteString("{CREATE_JSON_DUALITY_VIEW")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.EnableLogicalReplication {
		sb.WriteString(" :enableLogicalReplication true")
	}
	if n.DisableLogicalReplication {
		sb.WriteString(" :disableLogicalReplication true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterJsonDualityViewStmt(sb *strings.Builder, n *AlterJsonDualityViewStmt) {
	sb.WriteString("{ALTER_JSON_DUALITY_VIEW")
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

func writeCreateAttributeDimensionStmt(sb *strings.Builder, n *CreateAttributeDimensionStmt) {
	sb.WriteString("{CREATE_ATTRIBUTE_DIMENSION")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.NoForce {
		sb.WriteString(" :noForce true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	if n.DimensionType != "" {
		sb.WriteString(fmt.Sprintf(" :dimensionType %q", n.DimensionType))
	}
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	if n.Sources != nil {
		sb.WriteString(" :sources ")
		writeNode(sb, n.Sources)
	}
	if n.Attributes != nil {
		sb.WriteString(" :attributes ")
		writeNode(sb, n.Attributes)
	}
	if n.Levels != nil {
		sb.WriteString(" :levels ")
		writeNode(sb, n.Levels)
	}
	if n.AllClause != nil {
		sb.WriteString(" :allClause ")
		writeNode(sb, n.AllClause)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterAttributeDimensionStmt(sb *strings.Builder, n *AlterAttributeDimensionStmt) {
	sb.WriteString("{ALTER_ATTRIBUTE_DIMENSION")
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
	if n.NewName != nil {
		sb.WriteString(" :newName ")
		writeNode(sb, n.NewName)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAttrDimSourceClause(sb *strings.Builder, n *AttrDimSourceClause) {
	sb.WriteString("{ATTR_DIM_SOURCE")
	if n.Remote {
		sb.WriteString(" :remote true")
	}
	if n.IsJoinPath {
		sb.WriteString(" :isJoinPath true")
		if n.JoinPathName != "" {
			sb.WriteString(fmt.Sprintf(" :joinPathName %q", n.JoinPathName))
		}
		if n.JoinCondition != nil {
			sb.WriteString(" :joinCondition ")
			writeNode(sb, n.JoinCondition)
		}
	} else {
		if n.Name != nil {
			sb.WriteString(" :name ")
			writeNode(sb, n.Name)
		}
		if n.Alias != "" {
			sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
		}
	}
	sb.WriteString("}")
}

func writeAttrDimJoinCondElem(sb *strings.Builder, n *AttrDimJoinCondElem) {
	sb.WriteString("{JOIN_COND_ELEM")
	if n.LeftTable != "" {
		sb.WriteString(fmt.Sprintf(" :leftTable %q", n.LeftTable))
	}
	sb.WriteString(fmt.Sprintf(" :leftCol %q", n.LeftCol))
	if n.RightTable != "" {
		sb.WriteString(fmt.Sprintf(" :rightTable %q", n.RightTable))
	}
	sb.WriteString(fmt.Sprintf(" :rightCol %q", n.RightCol))
	sb.WriteString("}")
}

func writeAttrDimAttribute(sb *strings.Builder, n *AttrDimAttribute) {
	sb.WriteString("{ATTR_DIM_ATTR")
	sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
	}
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	sb.WriteString("}")
}

func writeAttrDimLevel(sb *strings.Builder, n *AttrDimLevel) {
	sb.WriteString("{ATTR_DIM_LEVEL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.LevelType != "" {
		sb.WriteString(fmt.Sprintf(" :levelType %q", n.LevelType))
	}
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	if n.KeyAttrs != nil {
		sb.WriteString(" :keyAttrs ")
		writeNode(sb, n.KeyAttrs)
	}
	if n.KeyNotNull {
		sb.WriteString(" :keyNotNull true")
	}
	if n.KeySkipWhenNull {
		sb.WriteString(" :keySkipWhenNull true")
	}
	if n.AltKeyAttrs != nil {
		sb.WriteString(" :altKeyAttrs ")
		writeNode(sb, n.AltKeyAttrs)
	}
	if n.MemberName != nil {
		sb.WriteString(" :memberName ")
		writeNode(sb, n.MemberName)
	}
	if n.MemberCaption != nil {
		sb.WriteString(" :memberCaption ")
		writeNode(sb, n.MemberCaption)
	}
	if n.MemberDesc != nil {
		sb.WriteString(" :memberDesc ")
		writeNode(sb, n.MemberDesc)
	}
	if n.OrderByAttrs != nil {
		sb.WriteString(" :orderByAttrs ")
		writeNode(sb, n.OrderByAttrs)
	}
	if n.Determines != nil {
		sb.WriteString(" :determines ")
		writeNode(sb, n.Determines)
	}
	sb.WriteString("}")
}

func writeAttrDimOrderByItem(sb *strings.Builder, n *AttrDimOrderByItem) {
	sb.WriteString("{ORDER_BY_ITEM")
	sb.WriteString(fmt.Sprintf(" :attribute %q", n.Attribute))
	if n.Direction != "" {
		sb.WriteString(fmt.Sprintf(" :direction %q", n.Direction))
	}
	sb.WriteString("}")
}

func writeAttrDimAllClause(sb *strings.Builder, n *AttrDimAllClause) {
	sb.WriteString("{ATTR_DIM_ALL")
	if n.MemberName != nil {
		sb.WriteString(" :memberName ")
		writeNode(sb, n.MemberName)
	}
	if n.MemberCaption != nil {
		sb.WriteString(" :memberCaption ")
		writeNode(sb, n.MemberCaption)
	}
	if n.MemberDesc != nil {
		sb.WriteString(" :memberDesc ")
		writeNode(sb, n.MemberDesc)
	}
	sb.WriteString("}")
}

func writeCreateHierarchyStmt(sb *strings.Builder, n *CreateHierarchyStmt) {
	sb.WriteString("{CREATE_HIERARCHY")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Force {
		sb.WriteString(" :force true")
	}
	if n.NoForce {
		sb.WriteString(" :noForce true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Sharing != "" {
		sb.WriteString(fmt.Sprintf(" :sharing %q", n.Sharing))
	}
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	if n.UsingAttrDim != nil {
		sb.WriteString(" :usingAttrDim ")
		writeNode(sb, n.UsingAttrDim)
	}
	if n.LevelHier != nil {
		sb.WriteString(" :levelHier ")
		writeNode(sb, n.LevelHier)
	}
	if n.HierAttrs != nil {
		sb.WriteString(" :hierAttrs ")
		writeNode(sb, n.HierAttrs)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterHierarchyStmt(sb *strings.Builder, n *AlterHierarchyStmt) {
	sb.WriteString("{ALTER_HIERARCHY")
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
	if n.NewName != nil {
		sb.WriteString(" :newName ")
		writeNode(sb, n.NewName)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeHierLevelClause(sb *strings.Builder, n *HierLevelClause) {
	sb.WriteString("{HIER_LEVEL")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	if n.ChildOf != nil {
		sb.WriteString(" :childOf ")
		writeNode(sb, n.ChildOf)
	}
	sb.WriteString("}")
}

func writeHierAttr(sb *strings.Builder, n *HierAttr) {
	sb.WriteString("{HIER_ATTR")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.Classifications != nil {
		sb.WriteString(" :classifications ")
		writeNode(sb, n.Classifications)
	}
	sb.WriteString("}")
}

func writeCreateDomainStmt(sb *strings.Builder, n *CreateDomainStmt) {
	sb.WriteString("{CREATE_DOMAIN")
	if n.OrReplace {
		sb.WriteString(" :orReplace true")
	}
	if n.Usecase {
		sb.WriteString(" :usecase true")
	}
	if n.IfNotExists {
		sb.WriteString(" :ifNotExists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.DomainType != "" {
		sb.WriteString(fmt.Sprintf(" :domainType %q", n.DomainType))
	}
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Strict {
		sb.WriteString(" :strict true")
	}
	if n.Default != nil {
		sb.WriteString(" :default ")
		writeNode(sb, n.Default)
	}
	if n.DefaultOnNull {
		sb.WriteString(" :defaultOnNull true")
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	if n.Collation != "" {
		sb.WriteString(fmt.Sprintf(" :collation %q", n.Collation))
	}
	if n.Display != nil {
		sb.WriteString(" :display ")
		writeNode(sb, n.Display)
	}
	if n.Order != nil {
		sb.WriteString(" :order ")
		writeNode(sb, n.Order)
	}
	if n.Annotations != nil {
		sb.WriteString(" :annotations ")
		writeNode(sb, n.Annotations)
	}
	if n.EnumItems != nil {
		sb.WriteString(" :enumItems ")
		writeNode(sb, n.EnumItems)
	}
	if n.Columns != nil {
		sb.WriteString(" :columns ")
		writeNode(sb, n.Columns)
	}
	if n.FlexDomainName != nil {
		sb.WriteString(" :flexDomainName ")
		writeNode(sb, n.FlexDomainName)
	}
	if n.FlexColumns != nil {
		sb.WriteString(" :flexColumns ")
		writeNode(sb, n.FlexColumns)
	}
	if n.ChooseUsing != nil {
		sb.WriteString(" :chooseUsing ")
		writeNode(sb, n.ChooseUsing)
	}
	if n.ChooseExpr != nil {
		sb.WriteString(" :chooseExpr ")
		writeNode(sb, n.ChooseExpr)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterDomainStmt(sb *strings.Builder, n *AlterDomainStmt) {
	sb.WriteString("{ALTER_DOMAIN")
	if n.Usecase {
		sb.WriteString(" :usecase true")
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
	if n.Display != nil {
		sb.WriteString(" :display ")
		writeNode(sb, n.Display)
	}
	if n.Order != nil {
		sb.WriteString(" :order ")
		writeNode(sb, n.Order)
	}
	if n.Annotations != nil {
		sb.WriteString(" :annotations ")
		writeNode(sb, n.Annotations)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeDomainColumn(sb *strings.Builder, n *DomainColumn) {
	sb.WriteString("{DOMAIN_COLUMN")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if n.DataType != nil {
		sb.WriteString(" :dataType ")
		writeNode(sb, n.DataType)
	}
	if n.Annotations != nil {
		sb.WriteString(" :annotations ")
		writeNode(sb, n.Annotations)
	}
	sb.WriteString("}")
}

func writeDomainEnumItem(sb *strings.Builder, n *DomainEnumItem) {
	sb.WriteString("{DOMAIN_ENUM_ITEM")
	sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	if len(n.Aliases) > 0 {
		sb.WriteString(" :aliases (")
		for i, a := range n.Aliases {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", a))
		}
		sb.WriteString(")")
	}
	if n.Value != nil {
		sb.WriteString(" :value ")
		writeNode(sb, n.Value)
	}
	sb.WriteString("}")
}

func writeDomainConstraint(sb *strings.Builder, n *DomainConstraint) {
	sb.WriteString("{DOMAIN_CONSTRAINT")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	sb.WriteString(fmt.Sprintf(" :type %q", n.Type))
	if n.CheckExpr != nil {
		sb.WriteString(" :checkExpr ")
		writeNode(sb, n.CheckExpr)
	}
	if n.State != nil {
		sb.WriteString(" :state ")
		writeNode(sb, n.State)
	}
	sb.WriteString("}")
}

func writeCreatePropertyGraphStmt(sb *strings.Builder, n *CreatePropertyGraphStmt) {
	sb.WriteString("{CREATE_PROPERTY_GRAPH")
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if len(n.VertexTables) > 0 {
		sb.WriteString(" :vertex_tables (")
		for i, v := range n.VertexTables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, v)
		}
		sb.WriteString(")")
	}
	if len(n.EdgeTables) > 0 {
		sb.WriteString(" :edge_tables (")
		for i, e := range n.EdgeTables {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString(")")
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeGraphTableDef(sb *strings.Builder, n *GraphTableDef) {
	sb.WriteString("{GRAPH_TABLE_DEF")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
	}
	if len(n.KeyColumns) > 0 {
		sb.WriteString(" :key_columns (")
		for i, c := range n.KeyColumns {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", c))
		}
		sb.WriteString(")")
	}
	if len(n.Labels) > 0 {
		sb.WriteString(" :labels (")
		for i, l := range n.Labels {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", l))
		}
		sb.WriteString(")")
	}
	if n.Properties != "" {
		sb.WriteString(fmt.Sprintf(" :properties %q", n.Properties))
	}
	if n.DefaultLabel {
		sb.WriteString(" :default_label true")
	}
	sb.WriteString("}")
}

func writeGraphEdgeDef(sb *strings.Builder, n *GraphEdgeDef) {
	sb.WriteString("{GRAPH_EDGE_DEF")
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Alias != "" {
		sb.WriteString(fmt.Sprintf(" :alias %q", n.Alias))
	}
	if n.SourceRef != "" {
		sb.WriteString(fmt.Sprintf(" :source_ref %q", n.SourceRef))
	}
	if n.DestRef != "" {
		sb.WriteString(fmt.Sprintf(" :dest_ref %q", n.DestRef))
	}
	sb.WriteString("}")
}

func writeGraphOptions(sb *strings.Builder, n *GraphOptions) {
	sb.WriteString("{GRAPH_OPTIONS")
	if n.Mode != "" {
		sb.WriteString(fmt.Sprintf(" :mode %q", n.Mode))
	}
	if n.MixedPropTypes != "" {
		sb.WriteString(fmt.Sprintf(" :mixed_prop_types %q", n.MixedPropTypes))
	}
	sb.WriteString("}")
}

func writeCreateVectorIndexStmt(sb *strings.Builder, n *CreateVectorIndexStmt) {
	sb.WriteString("{CREATE_VECTOR_INDEX")
	if n.IfNotExists {
		sb.WriteString(" :if_not_exists true")
	}
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.TableName != nil {
		sb.WriteString(" :table_name ")
		writeNode(sb, n.TableName)
	}
	if n.Column != "" {
		sb.WriteString(fmt.Sprintf(" :column %q", n.Column))
	}
	if n.Organization != "" {
		sb.WriteString(fmt.Sprintf(" :organization %q", n.Organization))
	}
	if n.Distance != "" {
		sb.WriteString(fmt.Sprintf(" :distance %q", n.Distance))
	}
	if n.TargetAccuracy > 0 {
		sb.WriteString(fmt.Sprintf(" :target_accuracy %d", n.TargetAccuracy))
	}
	if n.ParameterType != "" {
		sb.WriteString(fmt.Sprintf(" :parameter_type %q", n.ParameterType))
	}
	if n.Online {
		sb.WriteString(" :online true")
	}
	if n.Parallel > 0 {
		sb.WriteString(fmt.Sprintf(" :parallel %d", n.Parallel))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateLockdownProfileStmt(sb *strings.Builder, n *CreateLockdownProfileStmt) {
	sb.WriteString("{CREATE_LOCKDOWN_PROFILE")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Using != "" {
		sb.WriteString(fmt.Sprintf(" :using %q", n.Using))
	}
	if n.Including != "" {
		sb.WriteString(fmt.Sprintf(" :including %q", n.Including))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterLockdownProfileStmt(sb *strings.Builder, n *AlterLockdownProfileStmt) {
	sb.WriteString("{ALTER_LOCKDOWN_PROFILE")
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.RuleType != "" {
		sb.WriteString(fmt.Sprintf(" :rule_type %q", n.RuleType))
	}
	if len(n.RuleItems) > 0 {
		sb.WriteString(" :rule_items (")
		for i, item := range n.RuleItems {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%q", item))
		}
		sb.WriteString(")")
	}
	if n.AllItems {
		sb.WriteString(" :all_items true")
	}
	if n.Users != "" {
		sb.WriteString(fmt.Sprintf(" :users %q", n.Users))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeCreateOutlineStmt(sb *strings.Builder, n *CreateOutlineStmt) {
	sb.WriteString("{CREATE_OUTLINE")
	if n.OrReplace {
		sb.WriteString(" :or_replace true")
	}
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.Private {
		sb.WriteString(" :private true")
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Category != "" {
		sb.WriteString(fmt.Sprintf(" :category %q", n.Category))
	}
	if n.FromSource != "" {
		sb.WriteString(fmt.Sprintf(" :from_source %q", n.FromSource))
	}
	if n.FromPrivate {
		sb.WriteString(" :from_private true")
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}

func writeAlterOutlineStmt(sb *strings.Builder, n *AlterOutlineStmt) {
	sb.WriteString("{ALTER_OUTLINE")
	if n.Public {
		sb.WriteString(" :public true")
	}
	if n.Private {
		sb.WriteString(" :private true")
	}
	if n.Name != "" {
		sb.WriteString(fmt.Sprintf(" :name %q", n.Name))
	}
	if n.Action != "" {
		sb.WriteString(fmt.Sprintf(" :action %q", n.Action))
	}
	if n.NewName != "" {
		sb.WriteString(fmt.Sprintf(" :new_name %q", n.NewName))
	}
	if n.Category != "" {
		sb.WriteString(fmt.Sprintf(" :category %q", n.Category))
	}
	sb.WriteString(fmt.Sprintf(" :loc_start %d :loc_end %d", n.Loc.Start, n.Loc.End))
	sb.WriteString("}")
}
