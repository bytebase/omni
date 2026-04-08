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
// [ON CONFLICT ...] [RETURNING ...]`. Covers both legacy
// (INSERT INTO p VALUE …) and RFC 0011 (INSERT INTO c AS a VALUE …) forms.
//
// Grammar: insertCommand#InsertLegacy, insertCommand#Insert, insertCommandReturning
type InsertStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
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

// ExecStmt represents `EXEC name [arg, arg, ...]`.
//
// Grammar: execCommand
type ExecStmt struct {
	Name string
	Args []ExprNode
	Loc  Loc
}

func (*ExecStmt) nodeTag()      {}
func (n *ExecStmt) GetLoc() Loc { return n.Loc }
func (*ExecStmt) stmtNode()     {}

// ---------------------------------------------------------------------------
// PLACEHOLDERS — replaced in Task 9.
//
// These are forward-declared minimal stubs so the package builds at the
// end of Task 8 even though the real types live in the next task. Each
// stub is a bare struct with the Loc field plus the methods Node requires.
// Task 9 deletes this entire block and replaces it with the full type
// definitions of the clause helpers and DML helpers.
//
// IMPORTANT: This block also INCLUDES the OrderByItem placeholder
// originally added in Task 4 — DO NOT remove it during the wholesale
// stmts.go rewrite, or exprs.go (WindowSpec.OrderBy) will fail to build.
// ---------------------------------------------------------------------------

type OrderByItem struct {
	Expr          ExprNode
	Desc          bool
	NullsFirst    bool
	NullsExplicit bool
	Loc           Loc
}

func (*OrderByItem) nodeTag()      {}
func (n *OrderByItem) GetLoc() Loc { return n.Loc }

type TargetEntry struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*TargetEntry) nodeTag()      {}
func (n *TargetEntry) GetLoc() Loc { return n.Loc }

type PivotProjection struct {
	Value ExprNode
	At    ExprNode
	Loc   Loc
}

func (*PivotProjection) nodeTag()      {}
func (n *PivotProjection) GetLoc() Loc { return n.Loc }

type LetBinding struct {
	Expr  ExprNode
	Alias string
	Loc   Loc
}

func (*LetBinding) nodeTag()      {}
func (n *LetBinding) GetLoc() Loc { return n.Loc }

type GroupByClause struct {
	Partial bool
	Items   []*GroupByItem
	GroupAs *string
	Loc     Loc
}

func (*GroupByClause) nodeTag()      {}
func (n *GroupByClause) GetLoc() Loc { return n.Loc }

type GroupByItem struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*GroupByItem) nodeTag()      {}
func (n *GroupByItem) GetLoc() Loc { return n.Loc }

type SetAssignment struct {
	Target *PathExpr
	Value  ExprNode
	Loc    Loc
}

func (*SetAssignment) nodeTag()      {}
func (n *SetAssignment) GetLoc() Loc { return n.Loc }

type OnConflict struct {
	Target *OnConflictTarget
	Action OnConflictAction
	Where  ExprNode
	Loc    Loc
}

func (*OnConflict) nodeTag()      {}
func (n *OnConflict) GetLoc() Loc { return n.Loc }

type OnConflictTarget struct {
	Cols           []*VarRef
	ConstraintName string
	Loc            Loc
}

func (*OnConflictTarget) nodeTag()      {}
func (n *OnConflictTarget) GetLoc() Loc { return n.Loc }

type ReturningClause struct {
	Items []*ReturningItem
	Loc   Loc
}

func (*ReturningClause) nodeTag()      {}
func (n *ReturningClause) GetLoc() Loc { return n.Loc }

type ReturningItem struct {
	Status  ReturningStatus
	Mapping ReturningMapping
	Star    bool
	Expr    ExprNode
	Loc     Loc
}

func (*ReturningItem) nodeTag()      {}
func (n *ReturningItem) GetLoc() Loc { return n.Loc }
