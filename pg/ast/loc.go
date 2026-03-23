package ast

// NodeLoc extracts the Loc from any AST node that carries one.
// Returns NoLoc() if the node is nil or its type has no Loc field.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	// Statement nodes
	case *RawStmt:
		return v.Loc
	case *SelectStmt:
		return v.Loc
	case *InsertStmt:
		return v.Loc
	case *UpdateStmt:
		return v.Loc
	case *DeleteStmt:
		return v.Loc
	case *MergeStmt:
		return v.Loc
	case *TransactionStmt:
		return v.Loc
	case *DeallocateStmt:
		return v.Loc

	// --- Section 4.1: Database & schema nodes ---
	case *CreatedbStmt:
		return v.Loc
	case *DropdbStmt:
		return v.Loc
	case *AlterDatabaseStmt:
		return v.Loc
	case *AlterDatabaseSetStmt:
		return v.Loc
	case *CreateSchemaStmt:
		return v.Loc
	case *CreateTableSpaceStmt:
		return v.Loc
	case *CommentStmt:
		return v.Loc
	case *SecLabelStmt:

		return v.Loc

	// Expression nodes
	case *A_Expr:
		return v.Loc
	case *A_Const:
		return v.Loc
	case *A_ArrayExpr:
		return v.Loc
	case *BoolExpr:
		return v.Loc
	case *NullTest:
		return v.Loc
	case *BooleanTest:
		return v.Loc
	case *ColumnRef:
		return v.Loc
	case *FuncCall:
		return v.Loc
	case *TypeCast:
		return v.Loc
	case *SubLink:
		return v.Loc
	case *ParamRef:
		return v.Loc
	case *NamedArgExpr:
		return v.Loc
	case *CollateClause:
		return v.Loc
	case *CaseExpr:
		return v.Loc
	case *CaseWhen:
		return v.Loc
	case *CoalesceExpr:
		return v.Loc
	case *MinMaxExpr:
		return v.Loc
	case *NullIfExpr:
		return v.Loc
	case *RowExpr:
		return v.Loc
	case *ArrayExpr:
		return v.Loc
	case *GroupingFunc:
		return v.Loc
	case *GroupingSet:
		return v.Loc
	case *SQLValueFunction:
		return v.Loc
	case *SetToDefault:
		return v.Loc
	case *XmlExpr:
		return v.Loc
	case *XmlSerialize:
		return v.Loc

	// Clause/definition nodes
	case *RangeVar:
		return v.Loc
	case *ResTarget:
		return v.Loc
	case *SortBy:
		return v.Loc
	case *TypeName:
		return v.Loc
	case *ColumnDef:
		return v.Loc
	case *Constraint:
		return v.Loc
	case *DefElem:
		return v.Loc
	case *WindowDef:
		return v.Loc
	case *WithClause:
		return v.Loc
	case *CommonTableExpr:
		return v.Loc
	case *CTESearchClause:
		return v.Loc
	case *CTECycleClause:
		return v.Loc
	case *RoleSpec:
		return v.Loc
	case *OnConflictClause:
		return v.Loc
	case *InferClause:
		return v.Loc
	case *PartitionSpec:
		return v.Loc
	case *PartitionElem:
		return v.Loc
	case *PartitionBoundSpec:
		return v.Loc
	case *RangeTableSample:
		return v.Loc
	case *RangeTableFunc:
		return v.Loc
	case *RangeTableFuncCol:
		return v.Loc

	// JSON nodes
	case *JsonFormat:
		return v.Loc
	case *JsonBehavior:
		return v.Loc
	case *JsonFuncExpr:
		return v.Loc
	case *JsonTablePathSpec:
		return v.Loc
	case *JsonTableColumn:
		return v.Loc
	case *JsonTable:
		return v.Loc
	case *JsonParseExpr:
		return v.Loc
	case *JsonScalarExpr:
		return v.Loc
	case *JsonSerializeExpr:
		return v.Loc
	case *JsonObjectConstructor:
		return v.Loc
	case *JsonArrayConstructor:
		return v.Loc
	case *JsonArrayQueryConstructor:
		return v.Loc
	case *JsonAggConstructor:
		return v.Loc
	case *JsonIsPredicate:
		return v.Loc

	// Publication
	case *PublicationObjSpec:
		return v.Loc

	// --- Section 1.1: FROM/Join nodes ---
	case *JoinExpr:
		return v.Loc
	case *RangeSubselect:
		return v.Loc
	case *RangeFunction:
		return v.Loc
	case *CurrentOfExpr:
		return v.Loc
	case *LockingClause:
		return v.Loc
	case *IntoClause:
		return v.Loc

	// --- Section 1.2: Expression helper nodes ---
	case *A_Indices:
		return v.Loc
	case *A_Indirection:
		return v.Loc
	case *A_Star:
		return v.Loc
	case *MergeWhenClause:
		return v.Loc
	case *MultiAssignRef:
		return v.Loc
	case *TableLikeClause:
		return v.Loc

	// --- Section 1.3: JSON nodes ---
	case *JsonValueExpr:
		return v.Loc
	case *JsonOutput:
		return v.Loc
	case *JsonArgument:
		return v.Loc
	case *JsonKeyValue:
		return v.Loc
	case *JsonObjectAgg:
		return v.Loc
	case *JsonArrayAgg:
		return v.Loc

	// --- Section 2.1: Type & operator definitions ---
	case *DefineStmt:
		return v.Loc
	case *CompositeTypeStmt:
		return v.Loc
	case *CreateEnumStmt:
		return v.Loc
	case *CreateRangeStmt:
		return v.Loc
	case *CreateOpClassStmt:
		return v.Loc
	case *CreateOpFamilyStmt:
		return v.Loc
	case *CreateOpClassItem:
		return v.Loc
	case *CreateConversionStmt:
		return v.Loc
	case *CreateStatsStmt:
		return v.Loc
	case *AlterDefaultPrivilegesStmt:
		return v.Loc
	case *AlterOpFamilyStmt:
		return v.Loc
	case *AlterOperatorStmt:
		return v.Loc
	case *AlterStatsStmt:
		return v.Loc
	case *StatsElem:
		return v.Loc

	// --- Section 2.2: Extension nodes ---
	case *CreateExtensionStmt:
		return v.Loc
	case *AlterExtensionStmt:
		return v.Loc
	case *AlterExtensionContentsStmt:
		return v.Loc
	case *CreateAmStmt:
		return v.Loc
	case *CreateCastStmt:
		return v.Loc
	case *CreateTransformStmt:
		return v.Loc

	// --- Section 2.3: Foreign data wrapper nodes ---
	case *CreateFdwStmt:
		return v.Loc
	case *CreateForeignServerStmt:
		return v.Loc
	case *CreateForeignTableStmt:
		return v.Loc
	case *CreatePLangStmt:
		return v.Loc
	case *CreateUserMappingStmt:
		return v.Loc
	case *AlterFdwStmt:
		return v.Loc
	case *AlterForeignServerStmt:
		return v.Loc
	case *AlterUserMappingStmt:
		return v.Loc
	case *DropUserMappingStmt:
		return v.Loc
	case *ImportForeignSchemaStmt:
		return v.Loc

	// --- Section 3.1: General ALTER nodes ---
	case *AlterCollationStmt:
		return v.Loc
	case *AlterDomainStmt:
		return v.Loc
	case *AlterEnumStmt:
		return v.Loc
	case *AlterEventTrigStmt:
		return v.Loc
	case *AlterFunctionStmt:
		return v.Loc
	case *AlterObjectDependsStmt:
		return v.Loc
	case *AlterObjectSchemaStmt:
		return v.Loc
	case *AlterOwnerStmt:
		return v.Loc
	case *AlterTableSpaceOptionsStmt:
		return v.Loc
	case *AlterTSConfigurationStmt:
		return v.Loc
	case *AlterTSDictionaryStmt:
		return v.Loc
	case *AlterTypeStmt:
		return v.Loc

	// --- Section 3.2: ALTER TABLE nodes ---
	case *AlterTableStmt:
		return v.Loc
	case *AlterTableCmd:
		return v.Loc
	case *AlterSeqStmt:
		return v.Loc
	case *AlterTableMoveAllStmt:
		return v.Loc
	case *PartitionCmd:
		return v.Loc
	case *RenameStmt:
		return v.Loc

	// --- Section 3.3: Grant & role nodes ---
	case *GrantStmt:
		return v.Loc
	case *GrantRoleStmt:
		return v.Loc
	case *CreateRoleStmt:
		return v.Loc
	case *AlterRoleStmt:
		return v.Loc
	case *AlterRoleSetStmt:
		return v.Loc
	case *AlterPolicyStmt:
		return v.Loc
	case *CreatePolicyStmt:
		return v.Loc
	case *DropRoleStmt:
		return v.Loc
	case *AccessPriv:
		return v.Loc

	// --- Section 3.4: Publication & subscription nodes ---
	case *CreatePublicationStmt:
		return v.Loc
	case *AlterPublicationStmt:
		return v.Loc
	case *CreateSubscriptionStmt:
		return v.Loc
	case *AlterSubscriptionStmt:
		return v.Loc
	case *PublicationTable:
		return v.Loc
	case *RuleStmt:
		return v.Loc

	// --- Section 4.1: Database & schema nodes ---

	// --- Section 4.2: Sequence, function & domain nodes ---
	case *CreateSeqStmt:
		return v.Loc
	case *CreateDomainStmt:
		return v.Loc
	case *CreateFunctionStmt:
		return v.Loc
	case *FunctionParameter:
		return v.Loc
	case *ReturnStmt:
		return v.Loc

	// --- Section 4.3: Trigger, index & view nodes ---
	case *CreateTrigStmt:
		return v.Loc
	case *CreateEventTrigStmt:
		return v.Loc
	case *TriggerTransition:
		return v.Loc
	case *IndexStmt:
		return v.Loc
	case *IndexElem:
		return v.Loc
	case *ViewStmt:
		return v.Loc
	case *CreateTableAsStmt:
		return v.Loc
	case *RefreshMatViewStmt:
		return v.Loc
	case *CreateStmt:
		return v.Loc

	// --- Section 4.4: Utility statement nodes ---
	case *ExplainStmt:
		return v.Loc
	case *CallStmt:
		return v.Loc
	case *DoStmt:
		return v.Loc
	case *CheckPointStmt:
		return v.Loc
	case *DiscardStmt:
		return v.Loc
	case *ListenStmt:
		return v.Loc
	case *NotifyStmt:
		return v.Loc
	case *UnlistenStmt:
		return v.Loc
	case *LoadStmt:
		return v.Loc
	case *ReassignOwnedStmt:
		return v.Loc

	// --- Section 4.5: Cursor, prepare & IO nodes ---
	case *DeclareCursorStmt:
		return v.Loc
	case *FetchStmt:
		return v.Loc
	case *ClosePortalStmt:
		return v.Loc
	case *PrepareStmt:
		return v.Loc
	case *ExecuteStmt:
		return v.Loc
	case *CopyStmt:
		return v.Loc
	case *LockStmt:
		return v.Loc

	// --- Section 4.6: Maintenance & SET nodes ---
	case *VacuumStmt:
		return v.Loc
	case *VacuumRelation:
		return v.Loc
	case *ClusterStmt:
		return v.Loc
	case *ReindexStmt:
		return v.Loc
	case *VariableSetStmt:
		return v.Loc
	case *VariableShowStmt:
		return v.Loc
	case *AlterSystemStmt:
		return v.Loc
	case *ConstraintsSetStmt:
		return v.Loc
	case *DropStmt:
		return v.Loc
	case *DropOwnedStmt:
		return v.Loc
	case *DropSubscriptionStmt:
		return v.Loc
	case *DropTableSpaceStmt:
		return v.Loc
	case *TruncateStmt:
		return v.Loc

	// --- Section 5.1: Investigation nodes ---




	default:
		return NoLoc()
	}
}

// ListSpan returns the byte range spanning all items in a List.
// It uses NodeLoc on the first and last items to compute the range.
// Returns NoLoc() if the list is nil, empty, or items have no Loc.
func ListSpan(list *List) Loc {
	if list == nil || len(list.Items) == 0 {
		return NoLoc()
	}
	first := NodeLoc(list.Items[0])
	last := NodeLoc(list.Items[len(list.Items)-1])
	if first.Start == -1 || last.End == -1 {
		return NoLoc()
	}
	return Loc{Start: first.Start, End: last.End}
}
