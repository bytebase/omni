package ast

// ---------------------------------------------------------------------------
// Statement nodes and supporting types.
//
// This file is built in two task groups:
//   1. Top-level statements (this section): SELECT, set ops, EXPLAIN,
//      DML, DDL, EXEC.
//   2. Clause helpers: TargetEntry, GroupBy, OrderBy, OnConflict, Returning,
//      etc. (Task 9 — appended below).
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   root                   (lines 17–18)
//   statement              (lines 20–25)
//   dql                    (lines 55–56)
//   dml                    (lines 94–100)
//   dmlBaseCommand         (lines 102–108)
//   ddl                    (lines 73–76)
//   createCommand          (lines 78–81)
//   dropCommand            (lines 83–86)
//   execCommand            (lines 64–65)
//   selectClause           (lines 216–221)
//   exprBagOp              (lines 449–454)
//   exprSelect             (lines 456–467)
//   insertCommand          (lines 133–137)
//   insertCommandReturning (lines 130–131)
//   updateClause           (lines 182–183)
//   deleteCommand          (lines 191–192)
//   upsertCommand          (lines 124–125)
//   replaceCommand         (lines 120–121)
//   removeCommand          (lines 127–128)
//   setCommand             (lines 185–186)
// Each type below cites its specific rule#Label.
// ---------------------------------------------------------------------------

// ===========================================================================
// Statement enums
// ===========================================================================

// SetOpKind identifies the set operation in SetOpStmt.
type SetOpKind int

const (
	SetOpInvalid SetOpKind = iota
	SetOpUnion
	SetOpIntersect
	SetOpExcept
)

func (k SetOpKind) String() string {
	switch k {
	case SetOpUnion:
		return "UNION"
	case SetOpIntersect:
		return "INTERSECT"
	case SetOpExcept:
		return "EXCEPT"
	default:
		return "INVALID"
	}
}

// OnConflictAction discriminates the body of an ON CONFLICT clause.
// PartiQL's legacy ANTLR grammar leaves DO REPLACE/UPDATE as stubs that
// only accept the EXCLUDED keyword (see doReplace/doUpdate rules in
// PartiQLParser.g4 lines 168–180); the omni AST matches that scope.
type OnConflictAction int

const (
	OnConflictInvalid OnConflictAction = iota
	OnConflictDoNothing
	OnConflictDoReplaceExcluded
	OnConflictDoUpdateExcluded
)

func (a OnConflictAction) String() string {
	switch a {
	case OnConflictDoNothing:
		return "DO_NOTHING"
	case OnConflictDoReplaceExcluded:
		return "DO_REPLACE_EXCLUDED"
	case OnConflictDoUpdateExcluded:
		return "DO_UPDATE_EXCLUDED"
	default:
		return "INVALID"
	}
}

// ReturningStatus is the (MODIFIED|ALL) modifier on each RETURNING item.
type ReturningStatus int

const (
	ReturningStatusInvalid ReturningStatus = iota
	ReturningStatusModified
	ReturningStatusAll
)

func (s ReturningStatus) String() string {
	switch s {
	case ReturningStatusModified:
		return "MODIFIED"
	case ReturningStatusAll:
		return "ALL"
	default:
		return "INVALID"
	}
}

// ReturningMapping is the (OLD|NEW) modifier on each RETURNING item.
type ReturningMapping int

const (
	ReturningMappingInvalid ReturningMapping = iota
	ReturningMappingOld
	ReturningMappingNew
)

func (m ReturningMapping) String() string {
	switch m {
	case ReturningMappingOld:
		return "OLD"
	case ReturningMappingNew:
		return "NEW"
	default:
		return "INVALID"
	}
}

// ===========================================================================
// Top-level statements — all implement StmtNode
// ===========================================================================

// SelectStmt is a single SELECT query (not a set-op combination — see
// SetOpStmt for that). Carries the full clause set:
//
//	SELECT [DISTINCT|ALL] (* | items | VALUE expr | PIVOT v AT k)
//	FROM ...
//	[LET ...]
//	[WHERE ...]
//	[GROUP [PARTIAL] BY ... [GROUP AS alias]]
//	[HAVING ...]
//	[ORDER BY ...]
//	[LIMIT ...]
//	[OFFSET ...]
//
// Quantifier holds the optional DISTINCT/ALL modifier (NONE if absent).
// Star is true for `SELECT *`. Value holds the expression for
// `SELECT VALUE expr`. Pivot holds the PIVOT projection. At most one of
// (Star, Value != nil, Pivot != nil, len(Targets) > 0) is set per instance.
//
// Grammar: exprSelect#SfwQuery (with selectClause#SelectAll / selectClause#SelectItems /
//
//	selectClause#SelectValue / selectClause#SelectPivot)
type SelectStmt struct {
	Quantifier QuantifierKind   // NONE / DISTINCT / ALL
	Star       bool             // SELECT *
	Value      ExprNode         // SELECT VALUE expr (PartiQL-unique)
	Pivot      *PivotProjection // PIVOT v AT k (PartiQL-unique)
	Targets    []*TargetEntry   // SELECT items
	From       TableExpr        // FROM clause
	Let        []*LetBinding    // LET bindings (PartiQL-unique)
	Where      ExprNode
	GroupBy    *GroupByClause
	Having     ExprNode
	OrderBy    []*OrderByItem
	Limit      ExprNode
	Offset     ExprNode
	Loc        Loc
}

func (*SelectStmt) nodeTag()      {}
func (n *SelectStmt) GetLoc() Loc { return n.Loc }
func (*SelectStmt) stmtNode()     {}

// SetOpStmt combines two queries with UNION/INTERSECT/EXCEPT.
// Quantifier is the DISTINCT/ALL modifier; Outer is true for the
// PartiQL-specific OUTER variant. Both Left and Right may themselves
// be SetOpStmt nodes for nested set operations.
//
// Grammar: exprBagOp#Union, exprBagOp#Intersect, exprBagOp#Except
type SetOpStmt struct {
	Op         SetOpKind
	Quantifier QuantifierKind
	Outer      bool
	Left       StmtNode
	Right      StmtNode
	Loc        Loc
}

func (*SetOpStmt) nodeTag()      {}
func (n *SetOpStmt) GetLoc() Loc { return n.Loc }
func (*SetOpStmt) stmtNode()     {}

// ExplainStmt wraps any other StmtNode with an EXPLAIN prefix.
//
// Grammar: root (EXPLAIN prefix on the root rule)
type ExplainStmt struct {
	Inner StmtNode
	Loc   Loc
}

func (*ExplainStmt) nodeTag()      {}
func (n *ExplainStmt) GetLoc() Loc { return n.Loc }
func (*ExplainStmt) stmtNode()     {}

// InsertStmt represents `INSERT INTO target [AS alias] VALUE expr
// [AT pos] [ON CONFLICT ...] [RETURNING ...]`. Covers both legacy
// (INSERT INTO p VALUE … [AT pos]) and RFC 0011 (INSERT INTO c AS a VALUE …)
// forms.
//
// Mutual exclusion between the two forms:
//   - Legacy form (insertCommand#InsertLegacy, insertCommandReturning):
//     AsAlias is nil; Pos may be set.
//   - RFC 0011 form (insertCommand#Insert): AsAlias may be set; Pos is nil.
//   - The grammar's `insertCommandReturning` (line 130–131) is the only
//     alternative that allows `RETURNING`; on `insertCommand#InsertLegacy`
//     (line 134) and `insertCommand#Insert` (line 136), Returning is nil.
//
// Grammar: insertCommand#InsertLegacy, insertCommand#Insert, insertCommandReturning
type InsertStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	Pos        ExprNode // legacy `AT pos` clause; nil for RFC 0011 form
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*InsertStmt) nodeTag()      {}
func (n *InsertStmt) GetLoc() Loc { return n.Loc }
func (*InsertStmt) stmtNode()     {}

// UpdateStmt represents `UPDATE source SET ... [WHERE ...] [RETURNING ...]`
// and the equivalent `FROM source SET ...` form.
//
// Grammar: updateClause, dml#DmlBaseWrapper (with dmlBaseCommand containing setCommand)
type UpdateStmt struct {
	Source    TableExpr
	Sets      []*SetAssignment
	Where     ExprNode
	Returning *ReturningClause
	Loc       Loc
}

func (*UpdateStmt) nodeTag()      {}
func (n *UpdateStmt) GetLoc() Loc { return n.Loc }
func (*UpdateStmt) stmtNode()     {}

// DeleteStmt represents `DELETE FROM source [WHERE ...] [RETURNING ...]`
// and the equivalent `FROM source DELETE ...` form.
//
// Grammar: deleteCommand
type DeleteStmt struct {
	Source    TableExpr
	Where     ExprNode
	Returning *ReturningClause
	Loc       Loc
}

func (*DeleteStmt) nodeTag()      {}
func (n *DeleteStmt) GetLoc() Loc { return n.Loc }
func (*DeleteStmt) stmtNode()     {}

// UpsertStmt represents `UPSERT INTO target [AS alias] VALUE expr ...`.
// Same shape as InsertStmt; semantically the conflict-merge is implicit.
//
// Grammar: upsertCommand
type UpsertStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*UpsertStmt) nodeTag()      {}
func (n *UpsertStmt) GetLoc() Loc { return n.Loc }
func (*UpsertStmt) stmtNode()     {}

// ReplaceStmt represents `REPLACE INTO target [AS alias] VALUE expr ...`.
// Same shape as InsertStmt.
//
// Grammar: replaceCommand
type ReplaceStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*ReplaceStmt) nodeTag()      {}
func (n *ReplaceStmt) GetLoc() Loc { return n.Loc }
func (*ReplaceStmt) stmtNode()     {}

// RemoveStmt represents `REMOVE path` — removes an element from a
// collection. PartiQL-unique.
//
// Grammar: removeCommand
type RemoveStmt struct {
	Path *PathExpr
	Loc  Loc
}

func (*RemoveStmt) nodeTag()      {}
func (n *RemoveStmt) GetLoc() Loc { return n.Loc }
func (*RemoveStmt) stmtNode()     {}

// CreateTableStmt represents `CREATE TABLE name`. PartiQL DDL has no
// column definitions or constraints — just the table name.
//
// Grammar: createCommand#CreateTable
type CreateTableStmt struct {
	Name *VarRef
	Loc  Loc
}

func (*CreateTableStmt) nodeTag()      {}
func (n *CreateTableStmt) GetLoc() Loc { return n.Loc }
func (*CreateTableStmt) stmtNode()     {}

// CreateIndexStmt represents `CREATE INDEX ON table (path, path, ...)`.
//
// Grammar: createCommand#CreateIndex
type CreateIndexStmt struct {
	Table *VarRef
	Paths []*PathExpr
	Loc   Loc
}

func (*CreateIndexStmt) nodeTag()      {}
func (n *CreateIndexStmt) GetLoc() Loc { return n.Loc }
func (*CreateIndexStmt) stmtNode()     {}

// DropTableStmt represents `DROP TABLE name`.
//
// Grammar: dropCommand#DropTable
type DropTableStmt struct {
	Name *VarRef
	Loc  Loc
}

func (*DropTableStmt) nodeTag()      {}
func (n *DropTableStmt) GetLoc() Loc { return n.Loc }
func (*DropTableStmt) stmtNode()     {}

// DropIndexStmt represents `DROP INDEX index ON table`.
//
// Grammar: dropCommand#DropIndex
type DropIndexStmt struct {
	Index *VarRef
	Table *VarRef
	Loc   Loc
}

func (*DropIndexStmt) nodeTag()      {}
func (n *DropIndexStmt) GetLoc() Loc { return n.Loc }
func (*DropIndexStmt) stmtNode()     {}

// ExecStmt represents `EXEC name [arg, arg, ...]`. Per the grammar
// (PartiQLParser.g4 line 65: `EXEC name=expr ...`), the procedure name
// is itself an expression — typically a VarRef but may be any ExprNode
// (e.g., a parameter `?`) so the AST keeps the full breadth.
//
// Grammar: execCommand
type ExecStmt struct {
	Name ExprNode
	Args []ExprNode
	Loc  Loc
}

func (*ExecStmt) nodeTag()      {}
func (n *ExecStmt) GetLoc() Loc { return n.Loc }
func (*ExecStmt) stmtNode()     {}

// ===========================================================================
// SELECT clause helpers — bare Node (no sub-interface marker).
//
// These types appear only as fields/elements inside SelectStmt and never
// stand alone in scalar/statement/table-expr position, so they don't need
// a sub-interface marker.
// ===========================================================================

// TargetEntry represents one item in a SELECT projection list:
// `expr [AS alias]`.
//
// Grammar: projectionItem
type TargetEntry struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*TargetEntry) nodeTag()      {}
func (n *TargetEntry) GetLoc() Loc { return n.Loc }

// PivotProjection represents the body of a `PIVOT v AT k` projection.
// Used as SelectStmt.Pivot. PartiQL-unique.
//
// Grammar: selectClause#SelectPivot
type PivotProjection struct {
	Value ExprNode
	At    ExprNode
	Loc   Loc
}

func (*PivotProjection) nodeTag()      {}
func (n *PivotProjection) GetLoc() Loc { return n.Loc }

// LetBinding represents one `expr AS alias` binding inside a LET clause.
// PartiQL-unique.
//
// Grammar: letBinding
type LetBinding struct {
	Expr  ExprNode
	Alias string
	Loc   Loc
}

func (*LetBinding) nodeTag()      {}
func (n *LetBinding) GetLoc() Loc { return n.Loc }

// GroupByClause represents `GROUP [PARTIAL] BY items [GROUP AS alias]`.
//
// Grammar: groupClause
type GroupByClause struct {
	Partial bool
	Items   []*GroupByItem
	GroupAs *string
	Loc     Loc
}

func (*GroupByClause) nodeTag()      {}
func (n *GroupByClause) GetLoc() Loc { return n.Loc }

// GroupByItem represents one `expr [AS alias]` item in a GROUP BY list.
//
// Grammar: groupKey
type GroupByItem struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*GroupByItem) nodeTag()      {}
func (n *GroupByItem) GetLoc() Loc { return n.Loc }

// OrderByItem represents one `expr [ASC|DESC] [NULLS FIRST|LAST]` item
// in an ORDER BY list. NullsExplicit is true when the source text included
// a NULLS clause; NullsFirst is the resulting setting when it did.
//
// Grammar: orderSortSpec
type OrderByItem struct {
	Expr          ExprNode
	Desc          bool
	NullsFirst    bool
	NullsExplicit bool
	Loc           Loc
}

func (*OrderByItem) nodeTag()      {}
func (n *OrderByItem) GetLoc() Loc { return n.Loc }

// ===========================================================================
// DML helpers — bare Node.
// ===========================================================================

// SetAssignment represents one `target = value` assignment in an UPDATE
// SET clause (or in a chained-DML SET op).
//
// Grammar: setAssignment
type SetAssignment struct {
	Target *PathExpr
	Value  ExprNode
	Loc    Loc
}

func (*SetAssignment) nodeTag()      {}
func (n *SetAssignment) GetLoc() Loc { return n.Loc }

// OnConflict represents the body of an ON CONFLICT clause:
// `ON CONFLICT [target] [WHERE expr] action`.
//
// Grammar: onConflictClause#OnConflict, onConflictClause#OnConflictLegacy
type OnConflict struct {
	Target *OnConflictTarget
	Action OnConflictAction
	Where  ExprNode
	Loc    Loc
}

func (*OnConflict) nodeTag()      {}
func (n *OnConflict) GetLoc() Loc { return n.Loc }

// OnConflictTarget represents the target of an ON CONFLICT clause —
// either a list of column references `(col, col, ...)` or
// `ON CONSTRAINT name`. Exactly one of Cols/ConstraintName is set.
//
// Grammar: conflictTarget
type OnConflictTarget struct {
	Cols           []*VarRef
	ConstraintName string
	Loc            Loc
}

func (*OnConflictTarget) nodeTag()      {}
func (n *OnConflictTarget) GetLoc() Loc { return n.Loc }

// ReturningClause represents `RETURNING item, item, …`.
//
// Grammar: returningClause
type ReturningClause struct {
	Items []*ReturningItem
	Loc   Loc
}

func (*ReturningClause) nodeTag()      {}
func (n *ReturningClause) GetLoc() Loc { return n.Loc }

// ReturningItem represents one `(MODIFIED|ALL) (OLD|NEW) (* | expr)` entry
// in a RETURNING clause. Star is true for the `*` form; otherwise Expr
// holds the projection expression.
//
// Grammar: returningColumn
type ReturningItem struct {
	Status  ReturningStatus
	Mapping ReturningMapping
	Star    bool
	Expr    ExprNode
	Loc     Loc
}

func (*ReturningItem) nodeTag()      {}
func (n *ReturningItem) GetLoc() Loc { return n.Loc }
