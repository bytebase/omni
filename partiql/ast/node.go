// Package ast defines parse-tree node types for the omni PartiQL parser.
//
// PartiQL is the SQL++-flavored query language used by AWS DynamoDB and
// Azure Cosmos DB. This package mirrors the legacy bytebase/parser/partiql
// ANTLR grammar's full coverage scope (see docs/migration/partiql/analysis.md
// and docs/migration/partiql/dag.md).
//
// AST style: sealed sub-interfaces. Every node implements Node; most also
// implement one or more of StmtNode, ExprNode, TableExpr, PathStep, TypeName,
// or PatternNode for compile-time position discipline. A handful of small
// clause helpers (e.g. TargetEntry, CaseWhen) are bare Node.
//
// See docs/superpowers/specs/2026-04-08-partiql-ast-core-design.md for the
// design rationale.
package ast

// Loc is a half-open byte range describing where a node appears in the
// original source text. Loc{-1, -1} means the position is unknown
// (synthetic nodes constructed post-parse).
type Loc struct {
	Start int // inclusive byte offset
	End   int // exclusive byte offset
}

// Node is the root interface for all PartiQL parse-tree nodes.
//
// Implementations carry a Loc field and expose it via GetLoc. The
// unexported nodeTag method seals the interface so only types in this
// package can implement it.
type Node interface {
	nodeTag()
	GetLoc() Loc
}

// StmtNode marks top-level statement nodes — SELECT, INSERT, CREATE TABLE,
// EXEC, EXPLAIN, etc. Anything that can appear at script-statement position.
type StmtNode interface {
	Node
	stmtNode()
}

// ExprNode marks scalar-position expression nodes — operators, predicates,
// function calls, literals, paths, variables, subqueries, etc.
type ExprNode interface {
	Node
	exprNode()
}

// TableExpr marks nodes that appear in FROM position — table references,
// aliased sources, joins, unpivot.
//
// PathExpr, VarRef, and SubLink (defined in exprs.go) deliberately also
// implement TableExpr because PartiQL's grammar lets the same productions
// (path navigation, identifiers, parenthesized SELECTs) appear in both
// scalar and FROM position. See tableexprs.go for details.
type TableExpr interface {
	Node
	tableExpr()
}

// PathStep marks a single step in a PartiQL path expression:
// .field, .*, [expr], [*]. Path steps are chained inside a PathExpr.
type PathStep interface {
	Node
	pathStep()
}

// TypeName marks a PartiQL type reference, used by CAST and DDL.
// PartiQL's type system is small enough that a single concrete TypeRef
// implementation in types.go covers everything.
type TypeName interface {
	Node
	typeName()
}

// PatternNode marks GPML graph-match pattern nodes — node patterns,
// edge patterns, quantifiers, selectors.
type PatternNode interface {
	Node
	patternNode()
}

// List is a generic ordered collection of nodes used as a building block
// for cases where the elements are heterogeneous. Most parser-side lists
// use a typed slice (e.g. []ExprNode, []*TargetEntry) instead; List exists
// primarily so NodeToString has a uniform way to dump heterogenous groups.
type List struct {
	Items []Node
	Loc   Loc
}

func (*List) nodeTag()      {}
func (l *List) GetLoc() Loc { return l.Loc }

// Len returns the number of items in the list, treating a nil receiver as 0.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}
