package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's
// COMMENT ON and ANALYZE statements.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	| COMMENT_ ON_ TABLE_ qualifiedName IS_ (string_ | NULL_)  # commentTable
//	| COMMENT_ ON_ VIEW_ qualifiedName IS_ (string_ | NULL_)   # commentView
//	| COMMENT_ ON_ COLUMN_ qualifiedName IS_ (string_ | NULL_) # commentColumn
//	| ANALYZE_ qualifiedName (WITH_ properties)?               # analyze
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts:
//
//	D-CM1 (COMMENT value is a string OR NULL). `COMMENT ON TABLE t IS 'x'` and
//	   `COMMENT ON TABLE t IS NULL` (clear) are both accepted; anything else
//	   after IS is a SYNTAX_ERROR.
//	D-CM2 (object name part-count). TABLE/VIEW names are ≤ catalog.schema.object
//	   (3 parts); a COLUMN name is catalog.schema.table.column (≤ 4 parts), since
//	   it includes the column component (ddl.md comment-on-column).
//	D-AN1 (ANALYZE target ≤ 3 parts). The ANALYZE table name is bounded to
//	   catalog.schema.table.

// ---------------------------------------------------------------------------
// COMMENT ON
// ---------------------------------------------------------------------------

// CommentObjectKind classifies the COMMENT ON target.
type CommentObjectKind int

const (
	// CommentOnTable is COMMENT ON TABLE.
	CommentOnTable CommentObjectKind = iota
	// CommentOnView is COMMENT ON VIEW.
	CommentOnView
	// CommentOnColumn is COMMENT ON COLUMN.
	CommentOnColumn
)

// CommentStmt is `COMMENT ON {TABLE | VIEW | COLUMN} name IS (string | NULL)`.
// IsNull is true for the `IS NULL` clear form; Comment holds the decoded string
// otherwise.
type CommentStmt struct {
	Object  CommentObjectKind
	Name    *ast.QualifiedName
	IsNull  bool
	Comment string // decoded string value (empty and meaningless when IsNull)
	Loc     ast.Loc
}

func (n *CommentStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CommentStmt) Span() ast.Loc    { return n.Loc }

// parseCommentStmt parses COMMENT ON {TABLE|VIEW|COLUMN} name IS (string|NULL).
// On entry cur is COMMENT.
func (p *Parser) parseCommentStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume COMMENT
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	var (
		object   CommentObjectKind
		maxParts int
	)
	switch p.cur.Kind {
	case kwTABLE:
		p.advance()
		object, maxParts = CommentOnTable, 3
	case kwVIEW:
		p.advance()
		object, maxParts = CommentOnView, 3
	case kwCOLUMN:
		p.advance()
		// D-CM2: a column reference carries its column component too.
		object, maxParts = CommentOnColumn, 4
	default:
		return nil, p.syntaxErrorAtCur()
	}

	name, err := p.parseBoundedQualifiedName(maxParts, "comment target")
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIS); err != nil {
		return nil, err
	}

	stmt := &CommentStmt{Object: object, Name: name, Loc: ast.Loc{Start: startOffset}}

	// D-CM1: the comment value is a string literal or NULL.
	if nullTok, ok := p.match(kwNULL); ok {
		stmt.IsNull = true
		stmt.Loc.End = nullTok.Loc.End
		return stmt, nil
	}
	str, err := p.expectStringLiteral("COMMENT value")
	if err != nil {
		return nil, err
	}
	stmt.Comment = str.Str
	stmt.Loc.End = str.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ANALYZE
// ---------------------------------------------------------------------------

// AnalyzeStmt is `ANALYZE table_name [WITH (props)]`.
type AnalyzeStmt struct {
	Name       *ast.QualifiedName
	Properties []*Property
	Loc        ast.Loc
}

func (n *AnalyzeStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *AnalyzeStmt) Span() ast.Loc    { return n.Loc }

// parseAnalyzeStmt parses ANALYZE name [WITH (props)]. On entry cur is ANALYZE.
func (p *Parser) parseAnalyzeStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume ANALYZE
	// D-AN1: target is at most catalog.schema.table.
	name, err := p.parseBoundedQualifiedName(3, "analyze target")
	if err != nil {
		return nil, err
	}
	stmt := &AnalyzeStmt{Name: name, Loc: ast.Loc{Start: startOffset, End: name.Loc.End}}
	props, err := p.parseOptionalWithProperties()
	if err != nil {
		return nil, err
	}
	if props != nil {
		stmt.Properties = props
		stmt.Loc.End = p.prev.Loc.End
	}
	return stmt, nil
}

var (
	_ ast.Node = (*CommentStmt)(nil)
	_ ast.Node = (*AnalyzeStmt)(nil)
)
