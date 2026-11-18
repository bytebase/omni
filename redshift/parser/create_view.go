package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
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
	relpersistence, err := p.parseOptTemp()
	if err != nil {
		return nil, err
	}

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
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	rv.Relpersistence = relpersistence

	var aliases *nodes.List
	if recursive {
		// '(' columnList ')'
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var err error
		if aliases, err = p.parseColumnList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		// opt_column_list
		if aliases, err = p.parseOptColumnList(); err != nil {
			return nil, err
		}
	}

	// opt_reloptions
	options, err := p.parseOptReloptions()
	if err != nil {
		return nil, err
	}

	// AS
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	// SelectStmt
	var query nodes.Node

	if recursive {
		// For RECURSIVE VIEW, we need to wrap the select in a CTE.
		selectStmt, err := p.parseRequiredSelectNoParens()
		if err != nil {
			return nil, err
		}

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
			Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{Loc: nodes.NoLoc()}}},
			Loc:    nodes.NoLoc(),
		}
		rt := &nodes.ResTarget{Val: cr, Loc: nodes.NoLoc()}
		fromRv := &nodes.RangeVar{Relname: rv.Relname, Inh: true, Relpersistence: 'p', Loc: nodes.NoLoc()}
		sel := &nodes.SelectStmt{
			TargetList: &nodes.List{Items: []nodes.Node{rt}},
			FromClause: &nodes.List{Items: []nodes.Node{fromRv}},
			WithClause: wc,
		}
		query = sel
	} else {
		sel, err := p.parseRequiredSelectNoParens()
		if err != nil {
			return nil, err
		}
		query = sel
	}

	options, err = p.parseRedshiftViewTailOptions(options)
	if err != nil {
		return nil, err
	}

	// opt_check_option
	checkOption, err := p.parseOptCheckOption()
	if err != nil {
		return nil, err
	}

	return &nodes.ViewStmt{
		View:            rv,
		Aliases:         aliases,
		Options:         options,
		Query:           query,
		Replace:         replace,
		WithCheckOption: checkOption,
	}, nil
}

func (p *Parser) parseRedshiftViewTailOptions(options *nodes.List) (*nodes.List, error) {
	if p.cur.Type != WITH || p.peekNext().Type != NO {
		return options, nil
	}
	p.advance() // WITH
	p.advance() // NO
	if _, err := p.expect(SCHEMA); err != nil {
		return nil, err
	}
	if !p.curIdentEquals("binding") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // BINDING
	return appendToList(options, &nodes.DefElem{
		Defname: "no_schema_binding",
		Arg:     &nodes.Boolean{Boolval: true},
		Loc:     nodes.NoLoc(),
	}), nil
}

// parseOptCheckOption parses opt_check_option. It ends the view, after
// its SelectStmt, so WITH here can only lead it.
//
//	opt_check_option:
//	    WITH CHECK OPTION
//	    | WITH CASCADED CHECK OPTION
//	    | WITH LOCAL CHECK OPTION
//	    | /* EMPTY */
func (p *Parser) parseOptCheckOption() (int, error) {
	if p.cur.Type != WITH {
		return VIEW_CHECK_OPTION_NONE, nil
	}
	p.advance() // WITH

	option := VIEW_CHECK_OPTION_LOCAL
	switch p.cur.Type {
	case CASCADED:
		p.advance()
		option = VIEW_CHECK_OPTION_CASCADED
	case LOCAL:
		p.advance()
	}
	if _, err := p.expect(CHECK); err != nil {
		return 0, err
	}
	if _, err := p.expect(OPTION); err != nil {
		return 0, err
	}
	return option, nil
}

// parseOptReloptions parses opt_reloptions.
//
//	opt_reloptions:
//	    WITH reloptions
//	    | /* EMPTY */
func (p *Parser) parseOptReloptions() (*nodes.List, error) {
	if p.cur.Type == WITH && p.peekNext().Type == '(' {
		p.advance() // WITH
		return p.parseReloptions()
	}
	return nil, nil
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
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
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
		sel, err := p.parseRequiredSelectNoParens()
		if err != nil {
			return nil, err
		}
		query = sel
	}

	// opt_with_data
	withData, err := p.parseOptWithData()
	if err != nil {
		return nil, err
	}
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
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: loc, End: p.prev.End})

	colNames, err := p.parseOptColumnList()
	if err != nil {
		return nil, err
	}
	accessMethod, err := p.parseOptAccessMethod()
	if err != nil {
		return nil, err
	}
	options, err := p.parseOptWith()
	if err != nil {
		return nil, err
	}
	onCommit, err := p.parseOnCommitOption()
	if err != nil {
		return nil, err
	}
	tableSpace, err := p.parseOptTableSpace()
	if err != nil {
		return nil, err
	}

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
	query, err := p.parseRequiredSelectNoParens()
	if err != nil {
		return nil, err
	}

	// opt_with_data
	withData, err := p.parseOptWithData()
	if err != nil {
		return nil, err
	}
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
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: loc, End: p.prev.End})

	colNames, err := p.parseOptColumnList()
	if err != nil {
		return nil, err
	}
	accessMethod, err := p.parseOptAccessMethod()
	if err != nil {
		return nil, err
	}
	options, err := p.parseOptReloptions()
	if err != nil {
		return nil, err
	}
	options, err = p.parseRedshiftMaterializedViewOptionsInto(options, accessMethod)
	if err != nil {
		return nil, err
	}
	tableSpace, err := p.parseOptTableSpace()
	if err != nil {
		return nil, err
	}

	return &nodes.IntoClause{
		Rel:            rv,
		ColNames:       colNames,
		AccessMethod:   accessMethod,
		Options:        options,
		TableSpaceName: tableSpace,
		Loc:            nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRedshiftMaterializedViewOptionsInto parses the clauses
// r_CREATE_MATERIALIZED_VIEW puts after the view's name, its column list and
// its access method: BACKUP, DISTSTYLE, DISTKEY, SORTKEY and AUTO REFRESH,
// plus the LOCATION, PARTITIONED BY and TABLE PROPERTIES tail that only the
// USING ICEBERG form takes.
//
// Grammar reference: materializedViewCreateSqlCommand in RedshiftDocs.g4.
func (p *Parser) parseRedshiftMaterializedViewOptionsInto(options *nodes.List, accessMethod string) (*nodes.List, error) {
	iceberg := strings.EqualFold(accessMethod, "iceberg")
	for {
		consumed := false
		attrs, err := p.parseRedshiftCreateTableAttributes()
		if err != nil {
			return nil, err
		}
		if attrs != nil {
			options = appendNodeLists(options, attrs)
			consumed = true
		}
		if p.redshiftWordEqual("auto") {
			p.advance()
			if _, err := p.expect(REFRESH); err != nil {
				return nil, err
			}
			enabled, interval, err := p.parseRedshiftAutoRefreshSettings()
			if err != nil {
				return nil, err
			}
			options = appendToList(options, makeDefElem("auto_refresh", &nodes.Boolean{Boolval: enabled}))
			if interval != nil {
				options = appendToList(options, makeDefElem("auto_refresh_interval", interval))
			}
			consumed = true
		}
		if iceberg {
			var icebergConsumed bool
			options, icebergConsumed, err = p.parseRedshiftIcebergMaterializedViewOptions(options)
			if err != nil {
				return nil, err
			}
			consumed = consumed || icebergConsumed
		}
		if !consumed {
			return options, nil
		}
	}
}

// parseRedshiftIcebergMaterializedViewOptions parses the clauses the docs put
// after USING ICEBERG on r_CREATE_MATERIALIZED_VIEW, each optional and in this
// order:
//
//	[ LOCATION 's3://bucket/folder/' ]
//	[ PARTITIONED BY ( partition_transform [, ...] ) ]
//	[ TABLE PROPERTIES ( 'name' = 'value' [, ...] ) ]
//
// It reports whether it read one, so its caller's loop knows to go on. The
// table properties keep the shape ALTER TABLE ... SET TABLE PROPERTIES gives
// them: a list of one DefElem per 'name' = 'value' pair.
func (p *Parser) parseRedshiftIcebergMaterializedViewOptions(options *nodes.List) (*nodes.List, bool, error) {
	consumed := false
	for {
		switch {
		case p.cur.Type == LOCATION:
			p.advance()
			location, err := p.parseRedshiftStringLiteral()
			if err != nil {
				return nil, false, err
			}
			options = appendToList(options, makeDefElem("location", &nodes.String{Str: location}))
		case p.redshiftKeywordIs("partitioned"):
			// PARTITIONED BY is not a PostgreSQL clause; both words are names
			// there, so neither is a keyword token in this package.
			p.advance()
			if _, err := p.expect(BY); err != nil {
				return nil, false, err
			}
			transforms, err := p.parseRedshiftPartitionTransforms()
			if err != nil {
				return nil, false, err
			}
			options = appendToList(options, makeDefElem("partitioned_by", transforms))
		case p.cur.Type == TABLE:
			p.advance()
			if !p.consumeRedshiftWord("properties") {
				return nil, false, p.syntaxErrorAtCur()
			}
			properties, err := p.parseRedshiftStringOptionList()
			if err != nil {
				return nil, false, err
			}
			options = appendToList(options, makeDefElem("table_properties", properties))
		default:
			return options, consumed, nil
		}
		consumed = true
	}
}

// parseRedshiftPartitionTransforms parses the partition_transform list after
// PARTITIONED BY:
//
//	( partition_transform [, ...] )
func (p *Parser) parseRedshiftPartitionTransforms() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		transform, err := p.parseRedshiftPartitionTransform()
		if err != nil {
			return nil, err
		}
		items = append(items, transform)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseRedshiftPartitionTransform parses one partition_transform. The
// Parameters section of materialized-view-create-sql-command lists it, and
// docs-fixes.tsv has the converter read that section as its definition:
//
//	column_name | YEAR( column_name ) | MONTH( column_name ) |
//	DAY( column_name ) | HOUR( column_name ) |
//	BUCKET( number, column_name ) | TRUNCATE( number, column_name )
//
// Only those calls are transforms, so only those parse; a column is a ColumnRef
// and a call is a FuncCall of that name, which is what the transform means.
func (p *Parser) parseRedshiftPartitionTransform() (nodes.Node, error) {
	loc := p.pos()
	if p.isRedshiftPartitionTransformCallStart() {
		name := strings.ToLower(p.advance().Str)
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var args []nodes.Node
		if name == "bucket" || name == "truncate" {
			// BUCKET ( number, column_name )
			if !p.isRedshiftNumberStart() {
				return nil, p.syntaxErrorAtCur()
			}
			args = append(args, p.parseNumericOnly())
			if _, err := p.expect(','); err != nil {
				return nil, err
			}
		}
		column, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		args = append(args, &nodes.ColumnRef{
			Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: column}}},
			Loc:    nodes.Loc{Start: p.prev.Loc, End: p.prev.End},
		})
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.FuncCall{
			Funcname: makeFuncName(name),
			Args:     &nodes.List{Items: args},
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	column, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.ColumnRef{
		Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: column}}},
		Loc:    nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseOptWithData parses opt_with_data. It ends its statement, so WITH
// here can only lead it.
//
//	opt_with_data:
//	    WITH DATA
//	    | WITH NO DATA
//	    | /* EMPTY */
func (p *Parser) parseOptWithData() (bool, error) {
	if p.cur.Type != WITH {
		return true, nil // default: WITH DATA
	}
	p.advance() // WITH
	withData := true
	if p.cur.Type == NO {
		p.advance() // NO
		withData = false
	}
	if _, err := p.expect(DATA_P); err != nil {
		return false, err
	}
	return withData, nil
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
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	// opt_with_data
	withData, err := p.parseOptWithData()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()

	return &nodes.RefreshMatViewStmt{
		Concurrent: concurrent,
		SkipData:   !withData,
		Relation:   rv,
		Behavior:   behavior,
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
	list, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}
