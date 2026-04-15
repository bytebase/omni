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
	case *Property:
		return v.Loc
	case *BinaryExpr:
		return v.Loc
	case *UnaryExpr:
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
	case *JoinClause:
		return v.Loc
	case *SetOpStmt:
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
	case *MergeStmt:
		return v.Loc
	case *MergeClause:
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
	default:
		return NoLoc()
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
