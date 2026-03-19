package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseIndexStmt parses a CREATE INDEX statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-createindex.html
//
//	CREATE [ UNIQUE ] INDEX [ CONCURRENTLY ] [ [ IF NOT EXISTS ] name ] ON [ ONLY ] table_name [ USING method ]
//	    ( { column_name | ( expression ) } [ COLLATE collation ] [ opclass [ ( opclass_parameter = value [, ... ] ) ] ] [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...] )
//	    [ INCLUDE ( column_name [, ...] ) ]
//	    [ NULLS [ NOT ] DISTINCT ]
//	    [ WITH ( storage_parameter [= value] [, ... ] ) ]
//	    [ TABLESPACE tablespace_name ]
//	    [ WHERE predicate ]
func (p *Parser) parseIndexStmt() *nodes.IndexStmt {
	// CREATE already consumed by parseCreateDispatch

	// opt_unique
	unique := false
	if p.cur.Type == UNIQUE {
		p.advance()
		unique = true
	}

	p.expect(INDEX)

	// opt_concurrently
	concurrent := false
	if p.cur.Type == CONCURRENTLY {
		p.advance()
		concurrent = true
	}

	// IF NOT EXISTS name | opt_single_name
	ifNotExists := false
	idxname := ""
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
		idxname, _ = p.parseName()
	} else if p.cur.Type == ON {
		// unnamed index — skip to ON
	} else {
		// opt_single_name: could be a name or empty (ON follows)
		idxname, _ = p.parseName()
	}

	// ON
	p.expect(ON)

	// relation_expr
	rel, _ := p.parseRelationExpr()

	// access_method_clause
	accessMethod := ""
	if p.cur.Type == USING {
		p.advance()
		accessMethod, _ = p.parseName()
	}

	// '(' index_params ')'
	p.expect('(')
	indexParams := p.parseIndexParams()
	p.expect(')')

	// opt_unique_null_treatment
	nullsNotDistinct := p.parseOptUniqueNullTreatment()

	// opt_include
	var includeParams *nodes.List
	if p.cur.Type == INCLUDE {
		p.advance()
		p.expect('(')
		includeParams = p.parseIndexParams()
		p.expect(')')
	}

	// opt_reloptions
	var options *nodes.List
	if p.cur.Type == WITH {
		// Need to distinguish WITH (...) from WITH_LA (which is WITH TIME/ORDINALITY)
		// The advance() reclassification handles WITH -> WITH_LA for TIME/ORDINALITY,
		// so if cur is still WITH, it's reloptions.
		p.advance()
		options = p.parseReloptions()
	}

	// OptTableSpace
	tableSpace := p.parseOptTableSpace()

	// where_clause
	whereClause, _ := p.parseWhereClause()

	return &nodes.IndexStmt{
		Idxname:              idxname,
		Relation:             rel,
		AccessMethod:         accessMethod,
		IndexParams:          indexParams,
		IndexIncludingParams: includeParams,
		Options:              options,
		TableSpace:           tableSpace,
		WhereClause:          whereClause,
		Unique:               unique,
		Nulls_not_distinct:   nullsNotDistinct,
		Concurrent:           concurrent,
		IfNotExists:          ifNotExists,
	}
}

// parseIndexParams parses a comma-separated list of index elements.
//
//	index_params:
//	    index_elem
//	    | index_params ',' index_elem
func (p *Parser) parseIndexParams() *nodes.List {
	elem := p.parseIndexElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseIndexElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseIndexElem parses a single index element.
//
//	index_elem:
//	    ColId opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
//	    | ColId opt_collate any_name reloptions opt_asc_desc opt_nulls_order
//	    | func_expr_windowless opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
//	    | func_expr_windowless opt_collate any_name reloptions opt_asc_desc opt_nulls_order
//	    | '(' a_expr ')' opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
//	    | '(' a_expr ')' opt_collate any_name reloptions opt_asc_desc opt_nulls_order
func (p *Parser) parseIndexElem() *nodes.IndexElem {
	elem := &nodes.IndexElem{}

	if p.cur.Type == '(' {
		// Expression index: '(' a_expr ')'
		p.advance()
		elem.Expr, _ = p.parseAExpr(0)
		p.expect(')')
		elem.Collation = p.parseOptCollate()
		elem.Opclass, elem.Opclassopts = p.parseIndexElemOpclass()
		elem.Ordering = nodes.SortByDir(p.parseOptAscDesc())
		elem.NullsOrdering = nodes.SortByNulls(p.parseOptNullsOrder())
	} else if p.isColId() {
		// Could be ColId or func_expr_windowless (function call starting with a name)
		// We need to check if it's a function call: name '(' ...
		// But ColId is also valid as a simple column name.
		// Strategy: parse as ColId first, then check if '(' follows for function call
		name, _ := p.parseColId()

		if p.cur.Type == '(' {
			// This is a function call - reparse as func_expr_windowless
			// We need to reconstruct the function call
			elem.Expr = p.parseIndexElemFuncCall(name)
		} else {
			elem.Name = name
		}

		elem.Collation = p.parseOptCollate()
		elem.Opclass, elem.Opclassopts = p.parseIndexElemOpclass()
		elem.Ordering = nodes.SortByDir(p.parseOptAscDesc())
		elem.NullsOrdering = nodes.SortByNulls(p.parseOptNullsOrder())
	} else {
		// func_expr_windowless (non-identifier start, e.g. special function syntax)
		elem.Expr, _ = p.parseFuncExprWindowless()
		elem.Collation = p.parseOptCollate()
		elem.Opclass, elem.Opclassopts = p.parseIndexElemOpclass()
		elem.Ordering = nodes.SortByDir(p.parseOptAscDesc())
		elem.NullsOrdering = nodes.SortByNulls(p.parseOptNullsOrder())
	}

	return elem
}

// parseIndexElemFuncCall parses a function call for an index element,
// given that the function name has already been consumed.
func (p *Parser) parseIndexElemFuncCall(name string) nodes.Node {
	p.advance() // consume '('

	var args []nodes.Node
	if p.cur.Type != ')' {
		arg, _ := p.parseAExpr(0)
		args = append(args, arg)
		for p.cur.Type == ',' {
			p.advance()
			arg, _ = p.parseAExpr(0)
			args = append(args, arg)
		}
	}
	p.expect(')')

	funcCall := &nodes.FuncCall{
		Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
		Loc: nodes.NoLoc(),
	}
	if len(args) > 0 {
		funcCall.Args = &nodes.List{Items: args}
	}
	return funcCall
}

// parseIndexElemOpclass parses the optional operator class and its options.
// It handles both:
//   - opt_qualified_name (opclass without reloptions)
//   - any_name reloptions (opclass with reloptions)
//
// Returns (opclass, opclassopts).
func (p *Parser) parseIndexElemOpclass() (*nodes.List, *nodes.List) {
	// Check if what follows looks like an opclass name.
	// An opclass can't start with ASC, DESC, NULLS, ',', ')', or end-of-input.
	// It also can't be WITH, INCLUDE, WHERE, TABLESPACE which are subsequent clauses.
	if !p.isIndexElemOpclassStart() {
		return nil, nil
	}

	opclass := p.parseOptQualifiedName()
	if opclass == nil {
		return nil, nil
	}

	// Check for reloptions: '(' ... ')'
	var opclassopts *nodes.List
	if p.cur.Type == '(' {
		opclassopts = p.parseReloptions()
	}

	return opclass, opclassopts
}

// isIndexElemOpclassStart returns true if the current token could start an opclass name.
func (p *Parser) isIndexElemOpclassStart() bool {
	switch p.cur.Type {
	case ASC, DESC, NULLS_LA, ',', ')', 0,
		WITH, INCLUDE, WHERE, TABLESPACE,
		USING, DO:
		return false
	}
	return p.isColId() || p.isTypeFunctionName()
}
