package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// View check option constants matching PostgreSQL's VIEW_CHECK_OPTION_* values.
const (
	VIEW_CHECK_OPTION_NONE     = 0
	VIEW_CHECK_OPTION_LOCAL    = 1
	VIEW_CHECK_OPTION_CASCADED = 2
)

// parseViewStmt parses a CREATE [OR REPLACE] [TEMP|TEMPORARY] [RECURSIVE] VIEW statement.
// The CREATE keyword has already been consumed by the caller (parseCreateDispatch).
//
//	ViewStmt:
//	    CREATE OptTemp VIEW qualified_name opt_column_list opt_reloptions AS SelectStmt opt_check_option
//	    | CREATE OR REPLACE OptTemp VIEW qualified_name opt_column_list opt_reloptions AS SelectStmt opt_check_option
//	    | CREATE OptTemp RECURSIVE VIEW qualified_name '(' columnList ')' opt_reloptions AS SelectStmt opt_check_option
//	    | CREATE OR REPLACE OptTemp RECURSIVE VIEW qualified_name '(' columnList ')' opt_reloptions AS SelectStmt opt_check_option
func (p *Parser) parseViewStmt(replace bool) (*nodes.ViewStmt, error) {
	// OptTemp
	relpersistence := p.parseOptTemp()

	// RECURSIVE?
	recursive := false
	if p.cur.Type == RECURSIVE {
		p.advance()
		recursive = true
	}

	// VIEW
	if _, err := p.expect(VIEW); err != nil {
		return nil, err
	}

	// qualified_name
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = relpersistence

	var aliases *nodes.List
	if recursive {
		// '(' columnList ')'
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		aliases = p.parseColumnList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		// opt_column_list
		aliases, _ = p.parseOptColumnList()
	}

	// opt_reloptions
	options := p.parseOptReloptions()

	// AS
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	// SelectStmt
	var query nodes.Node

	if recursive {
		// For RECURSIVE VIEW, we need to wrap the select in a CTE.
		selectStmt, _ := p.parseSelectNoParens()

		cte := &nodes.CommonTableExpr{
			Ctename:       rv.Relname,
			Ctequery:      selectStmt,
			Aliascolnames: aliases,
		}
		wc := &nodes.WithClause{
			Ctes:      &nodes.List{Items: []nodes.Node{cte}},
			Recursive: true,
		}
		cr := &nodes.ColumnRef{
			Fields:   &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
			Loc: nodes.NoLoc(),
		}
		rt := &nodes.ResTarget{Val: cr, Loc: nodes.NoLoc()}
		fromRv := &nodes.RangeVar{Relname: rv.Relname, Inh: true, Relpersistence: 'p'}
		sel := &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{rt}},
			FromClause: &nodes.List{Items: []nodes.Node{fromRv}},
			WithClause: wc,
		}
		query = sel
	} else {
		query, _ = p.parseSelectNoParens()
	}

	// opt_check_option
	checkOption := p.parseOptCheckOption()

	return &nodes.ViewStmt{
		View:            rv,
		Aliases:         aliases,
		Options:         options,
		Query:           query,
		Replace:         replace,
		WithCheckOption: checkOption,
	}, nil
}

// parseOptCheckOption parses opt_check_option.
//
//	opt_check_option:
//	    WITH CHECK OPTION
//	    | WITH CASCADED CHECK OPTION
//	    | WITH LOCAL CHECK OPTION
//	    | /* EMPTY */
func (p *Parser) parseOptCheckOption() int {
	if p.cur.Type != WITH {
		return VIEW_CHECK_OPTION_NONE
	}

	next := p.peekNext()
	switch next.Type {
	case CHECK:
		// WITH CHECK OPTION
		p.advance() // WITH
		p.advance() // CHECK
		p.expect(OPTION)
		return VIEW_CHECK_OPTION_LOCAL
	case CASCADED:
		// WITH CASCADED CHECK OPTION
		p.advance() // WITH
		p.advance() // CASCADED
		p.expect(CHECK)
		p.expect(OPTION)
		return VIEW_CHECK_OPTION_CASCADED
	case LOCAL:
		// WITH LOCAL CHECK OPTION
		p.advance() // WITH
		p.advance() // LOCAL
		p.expect(CHECK)
		p.expect(OPTION)
		return VIEW_CHECK_OPTION_LOCAL
	default:
		return VIEW_CHECK_OPTION_NONE
	}
}

// parseOptReloptions parses opt_reloptions.
//
//	opt_reloptions:
//	    WITH reloptions
//	    | /* EMPTY */
func (p *Parser) parseOptReloptions() *nodes.List {
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == '(' {
			p.advance() // WITH
			return p.parseReloptions()
		}
	}
	return nil
}

// parseCreateTableAsStmt parses a CREATE [TEMP] TABLE ... AS statement.
// The CREATE keyword has already been consumed. OptTemp has been parsed.
// TABLE has been consumed. IF NOT EXISTS has been parsed.
//
//	CreateAsStmt:
//	    CREATE OptTemp TABLE create_as_target AS SelectStmt opt_with_data
//	    | CREATE OptTemp TABLE IF NOT EXISTS create_as_target AS SelectStmt opt_with_data
//	    | CREATE OptTemp TABLE create_as_target AS EXECUTE name execute_param_clause opt_with_data
//	    | CREATE OptTemp TABLE IF NOT EXISTS create_as_target AS EXECUTE name execute_param_clause opt_with_data
func (p *Parser) parseCreateTableAsStmt(relpersistence byte, ifNotExists bool) (*nodes.CreateTableAsStmt, error) {
	// create_as_target: qualified_name opt_column_list OptAccessMethod OptWith OnCommitOption OptTableSpace
	into, err := p.parseCreateAsTarget()
	if err != nil {
		return nil, err
	}
	into.Rel.Relpersistence = relpersistence

	// AS
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	var query nodes.Node

	if p.cur.Type == EXECUTE {
		// AS EXECUTE name execute_param_clause
		p.advance() // EXECUTE
		name, _ := p.parseName()
		params, err := p.parseExecuteParamClause()
		if err != nil {
			return nil, err
		}
		query = &nodes.ExecuteStmt{
			Name:   name,
			Params: params,
		}
	} else {
		// AS SelectStmt
		query, _ = p.parseSelectNoParens()
	}

	// opt_with_data
	withData := p.parseOptWithData()
	into.SkipData = !withData

	return &nodes.CreateTableAsStmt{
		Query:       query,
		Into:        into,
		Objtype:     nodes.OBJECT_TABLE,
		IfNotExists: ifNotExists,
	}, nil
}

// parseCreateAsTarget parses create_as_target.
//
//	create_as_target:
//	    qualified_name opt_column_list OptAccessMethod OptWith OnCommitOption OptTableSpace
func (p *Parser) parseCreateAsTarget() (*nodes.IntoClause, error) {
	loc := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromAnyName(names)

	colNames, _ := p.parseOptColumnList()
	accessMethod := p.parseOptAccessMethod()
	options := p.parseOptWith()
	onCommit := p.parseOnCommitOption()
	tableSpace := p.parseOptTableSpace()

	return &nodes.IntoClause{
		Rel:            rv,
		ColNames:       colNames,
		AccessMethod:   accessMethod,
		Options:        options,
		OnCommit:       onCommit,
		TableSpaceName: tableSpace,
		Loc:            nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseCreateMatViewStmt parses a CREATE [UNLOGGED] MATERIALIZED VIEW statement.
// The CREATE keyword has already been consumed by the caller.
// The MATERIALIZED keyword has NOT been consumed yet.
//
//	CreateMatViewStmt:
//	    CREATE OptNoLog MATERIALIZED VIEW create_mv_target AS SelectStmt opt_with_data
//	    | CREATE OptNoLog MATERIALIZED VIEW IF NOT EXISTS create_mv_target AS SelectStmt opt_with_data
func (p *Parser) parseCreateMatViewStmt(relpersistence byte) (*nodes.CreateTableAsStmt, error) {
	// MATERIALIZED has already been consumed by the caller.
	// VIEW
	if _, err := p.expect(VIEW); err != nil {
		return nil, err
	}

	// IF NOT EXISTS
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}

	// create_mv_target: qualified_name opt_column_list OptAccessMethod opt_reloptions OptTableSpace
	into, err := p.parseCreateMvTarget()
	if err != nil {
		return nil, err
	}
	into.Rel.Relpersistence = relpersistence

	// AS
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	// SelectStmt
	query, _ := p.parseSelectNoParens()

	// opt_with_data
	withData := p.parseOptWithData()
	into.SkipData = !withData

	return &nodes.CreateTableAsStmt{
		Query:       query,
		Into:        into,
		Objtype:     nodes.OBJECT_MATVIEW,
		IfNotExists: ifNotExists,
	}, nil
}

// parseCreateMvTarget parses create_mv_target.
//
//	create_mv_target:
//	    qualified_name opt_column_list OptAccessMethod opt_reloptions OptTableSpace
func (p *Parser) parseCreateMvTarget() (*nodes.IntoClause, error) {
	loc := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromAnyName(names)

	colNames, _ := p.parseOptColumnList()
	accessMethod := p.parseOptAccessMethod()
	options := p.parseOptReloptions()
	tableSpace := p.parseOptTableSpace()

	return &nodes.IntoClause{
		Rel:            rv,
		ColNames:       colNames,
		AccessMethod:   accessMethod,
		Options:        options,
		TableSpaceName: tableSpace,
		Loc:            nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseOptWithData parses opt_with_data.
//
//	opt_with_data:
//	    WITH DATA
//	    | WITH NO DATA
//	    | /* EMPTY */
func (p *Parser) parseOptWithData() bool {
	if p.cur.Type != WITH {
		return true // default: WITH DATA
	}

	next := p.peekNext()
	switch next.Type {
	case DATA_P:
		p.advance() // WITH
		p.advance() // DATA
		return true
	case NO:
		p.advance() // WITH
		p.advance() // NO
		p.expect(DATA_P)
		return false
	default:
		return true
	}
}

// parseRefreshMatViewStmt parses a REFRESH MATERIALIZED VIEW statement.
// The REFRESH keyword has already been consumed by the caller.
//
//	RefreshMatViewStmt:
//	    REFRESH MATERIALIZED VIEW opt_concurrently qualified_name opt_with_data
func (p *Parser) parseRefreshMatViewStmt() (*nodes.RefreshMatViewStmt, error) {
	// MATERIALIZED
	if _, err := p.expect(MATERIALIZED); err != nil {
		return nil, err
	}
	// VIEW
	if _, err := p.expect(VIEW); err != nil {
		return nil, err
	}

	// opt_concurrently
	concurrent := false
	if p.cur.Type == CONCURRENTLY {
		p.advance()
		concurrent = true
	}

	// qualified_name
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromAnyName(names)

	// opt_with_data
	withData := p.parseOptWithData()

	return &nodes.RefreshMatViewStmt{
		Concurrent: concurrent,
		SkipData:   !withData,
		Relation:   rv,
	}, nil
}

// parseExecuteParamClause parses execute_param_clause.
//
//	execute_param_clause:
//	    '(' expr_list ')'
//	    | /* EMPTY */
func (p *Parser) parseExecuteParamClause() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance() // '('
	list := p.parseExprList()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}
