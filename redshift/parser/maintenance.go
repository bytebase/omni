package parser

import (
	nodes "github.com/bytebase/omni/redshift/ast"
)

// parseVacuumStmt parses a VACUUM statement.
//
//	VacuumStmt:
//	    VACUUM opt_full opt_freeze opt_verbose opt_analyze opt_vacuum_relation_list
//	    | VACUUM '(' utility_option_list ')' opt_vacuum_relation_list
func (p *Parser) parseVacuumStmt() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume VACUUM

	// VACUUM '(' utility_option_list ')' opt_vacuum_relation_list
	if p.cur.Type == '(' {
		p.advance() // consume '('
		opts, err := p.parseUtilityOptionList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		rels, err := p.parseOptVacuumRelationList()
		if err != nil {
			return nil, err
		}
		return &nodes.VacuumStmt{
			Options:     opts,
			Rels:        rels,
			IsVacuumCmd: true,
			Loc:         nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}

	// VACUUM opt_full opt_freeze opt_verbose opt_analyze opt_vacuum_relation_list
	n := &nodes.VacuumStmt{
		IsVacuumCmd: true,
	}
	var opts []nodes.Node

	// opt_full
	if p.cur.Type == FULL {
		p.advance()
		opts = append(opts, makeDefElem("full", nil))
		opts = append(opts, makeDefElem("mode", &nodes.String{Str: "full"}))
	} else if p.redshiftWordEqual("sort") {
		p.advance()
		if _, err := p.expect(ONLY); err != nil {
			return nil, err
		}
		opts = append(opts, makeDefElem("mode", &nodes.String{Str: "sort_only"}))
	} else if p.cur.Type == DELETE_P {
		p.advance()
		if _, err := p.expect(ONLY); err != nil {
			return nil, err
		}
		opts = append(opts, makeDefElem("mode", &nodes.String{Str: "delete_only"}))
	} else if p.cur.Type == REINDEX {
		p.advance()
		opts = append(opts, makeDefElem("mode", &nodes.String{Str: "reindex"}))
	} else if p.redshiftWordEqual("recluster") {
		p.advance()
		opts = append(opts, makeDefElem("mode", &nodes.String{Str: "recluster"}))
	}
	// opt_freeze
	if p.cur.Type == FREEZE {
		p.advance()
		opts = append(opts, makeDefElem("freeze", nil))
	}
	// opt_verbose
	if p.cur.Type == VERBOSE {
		p.advance()
		opts = append(opts, makeDefElem("verbose", nil))
	}
	// opt_analyze
	if p.cur.Type == ANALYZE || p.cur.Type == ANALYSE {
		p.advance()
		opts = append(opts, makeDefElem("analyze", nil))
	}

	if len(opts) > 0 {
		n.Options = &nodes.List{Items: opts}
	}

	if !p.redshiftVacuumOptionStart() {
		rels, err := p.parseOptVacuumRelationList()
		if err != nil {
			return nil, err
		}
		n.Rels = rels
	}
	n.Options = p.parseRedshiftVacuumOptions(n.Options)
	n.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return n, nil
}

func (p *Parser) redshiftVacuumOptionStart() bool {
	return p.cur.Type == TO || p.redshiftWordEqual("boost")
}

func (p *Parser) parseRedshiftVacuumOptions(options *nodes.List) *nodes.List {
	for {
		switch {
		case p.cur.Type == TO:
			p.advance()
			threshold := p.parseSignedIconst()
			if !p.consumeRedshiftWord("percent") {
				return options
			}
			options = appendToList(options, makeDefElem("threshold_percent", &nodes.Integer{Ival: threshold}))
		case p.redshiftWordEqual("boost"):
			p.advance()
			options = appendToList(options, makeDefElem("boost", nil))
		default:
			return options
		}
	}
}

// parseAnalyzeStmt parses an ANALYZE/ANALYSE statement.
//
//	AnalyzeStmt:
//	    analyze_keyword opt_verbose opt_vacuum_relation_list
//	    | analyze_keyword '(' utility_option_list ')' opt_vacuum_relation_list
func (p *Parser) parseAnalyzeStmt() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume ANALYZE/ANALYSE

	// analyze_keyword '(' utility_option_list ')' opt_vacuum_relation_list
	if p.cur.Type == '(' {
		p.advance() // consume '('
		opts, err := p.parseUtilityOptionList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		rels, err := p.parseOptVacuumRelationList()
		if err != nil {
			return nil, err
		}
		return &nodes.VacuumStmt{
			Options:     opts,
			Rels:        rels,
			IsVacuumCmd: false,
			Loc:         nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}

	// analyze_keyword opt_verbose opt_vacuum_relation_list
	n := &nodes.VacuumStmt{
		IsVacuumCmd: false,
	}
	if p.cur.Type == COMPRESSION {
		p.advance()
		n.Options = appendToList(n.Options, makeDefElem("compression", nil))
	}
	if p.cur.Type == VERBOSE {
		p.advance()
		n.Options = appendToList(n.Options, makeDefElem("verbose", nil))
	}

	if !p.redshiftAnalyzeOptionStart() {
		rels, err := p.parseOptVacuumRelationList()
		if err != nil {
			return nil, err
		}
		n.Rels = rels
	}
	n.Options = p.parseRedshiftAnalyzeOptions(n.Options)
	n.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return n, nil
}

func (p *Parser) redshiftAnalyzeOptionStart() bool {
	return p.cur.Type == ALL || p.redshiftWordEqual("predicate") || p.redshiftWordEqual("comprows")
}

func (p *Parser) parseRedshiftAnalyzeOptions(options *nodes.List) *nodes.List {
	for {
		switch {
		case p.redshiftWordEqual("predicate"):
			p.advance()
			p.expect(COLUMNS)
			options = appendToList(options, makeDefElem("predicate_columns", nil))
		case p.cur.Type == ALL:
			p.advance()
			p.expect(COLUMNS)
			options = appendToList(options, makeDefElem("all_columns", nil))
		case p.redshiftWordEqual("comprows"):
			p.advance()
			rows := p.parseSignedIconst()
			options = appendToList(options, makeDefElem("comprows", &nodes.Integer{Ival: rows}))
		default:
			return options
		}
	}
}

// parseOptVacuumRelationList parses opt_vacuum_relation_list.
//
//	opt_vacuum_relation_list:
//	    vacuum_relation_list
//	    | /* EMPTY */
func (p *Parser) parseOptVacuumRelationList() (*nodes.List, error) {
	// Check if current token could start a vacuum_relation (i.e. a qualified_name).
	if !p.isColId() {
		return nil, nil
	}
	return p.parseVacuumRelationList()
}

// parseVacuumRelationList parses vacuum_relation_list.
//
//	vacuum_relation_list:
//	    vacuum_relation
//	    | vacuum_relation_list ',' vacuum_relation
func (p *Parser) parseVacuumRelationList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		rel, err := p.parseVacuumRelation()
		if err != nil {
			return nil, err
		}
		items = append(items, rel)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseVacuumRelation parses vacuum_relation.
//
//	vacuum_relation:
//	    qualified_name opt_name_list
func (p *Parser) parseVacuumRelation() (*nodes.VacuumRelation, error) {
	loc := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: loc, End: p.prev.End})
	cols, err := p.parseOptColumnList()
	if err != nil {
		return nil, err
	}
	return &nodes.VacuumRelation{
		Relation: rv,
		VaCols:   cols,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseClusterStmt parses a CLUSTER statement.
//
//	ClusterStmt:
//	    CLUSTER '(' utility_option_list ')' qualified_name cluster_index_specification
//	    | CLUSTER '(' utility_option_list ')'
//	    | CLUSTER opt_verbose qualified_name cluster_index_specification
//	    | CLUSTER opt_verbose
//	    | CLUSTER opt_verbose name ON qualified_name
func (p *Parser) parseClusterStmt() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume CLUSTER

	// CLUSTER '(' utility_option_list ')' [qualified_name cluster_index_specification]
	if p.cur.Type == '(' {
		p.advance() // consume '('
		params, err := p.parseUtilityOptionList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		// Check if a qualified_name follows
		if p.isColId() {
			nameStart := p.pos()
			names, err := p.parseQualifiedName()
			if err != nil {
				// shape-III-partial fix: propagate the error instead of silently
				// returning a ClusterStmt with nil Relation.
				return nil, err
			}
			rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
			idxName, err := p.parseClusterIndexSpecification()
			if err != nil {
				return nil, err
			}
			return &nodes.ClusterStmt{
				Relation:  rv,
				Indexname: idxName,
				Params:    params,
				Loc:       nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		return &nodes.ClusterStmt{Params: params, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}

	// opt_verbose
	verbose := false
	if p.cur.Type == VERBOSE {
		p.advance()
		verbose = true
	}

	// Check if we have a qualified_name
	if !p.isColId() {
		// CLUSTER opt_verbose (no table)
		n := &nodes.ClusterStmt{Loc: nodes.Loc{Start: loc, End: p.prev.End}}
		if verbose {
			n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
		}
		return n, nil
	}

	// We have a name. Could be:
	//   CLUSTER opt_verbose qualified_name cluster_index_specification
	//   CLUSTER opt_verbose name ON qualified_name
	// We need to distinguish: after parsing the name, if ON follows, it's the
	// backwards-compatible syntax "CLUSTER [VERBOSE] name ON qualified_name".
	// Otherwise it's "CLUSTER [VERBOSE] qualified_name [USING name]".

	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		// shape-III-partial fix: propagate the error instead of silently
		// returning a ClusterStmt with nil Relation.
		return nil, err
	}
	nameLoc := nodes.Loc{Start: nameStart, End: p.prev.End}

	if p.cur.Type == ON {
		// Backwards-compatible: CLUSTER [VERBOSE] name ON qualified_name
		// The "name" we parsed is actually the index name (single name, not qualified).
		p.advance() // consume ON
		indexName := ""
		if names != nil && len(names.Items) == 1 {
			indexName = names.Items[0].(*nodes.String).Str
		}
		tableStart := p.pos()
		tableNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		rv := makeRangeVarFromNames(tableNames, nodes.Loc{Start: tableStart, End: p.prev.End})
		n := &nodes.ClusterStmt{
			Relation:  rv,
			Indexname: indexName,
			Loc:       nodes.Loc{Start: loc, End: p.prev.End},
		}
		if verbose {
			n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
		}
		return n, nil
	}

	// CLUSTER [VERBOSE] qualified_name [USING name]
	rv := makeRangeVarFromNames(names, nameLoc)
	idxName, err := p.parseClusterIndexSpecification()
	if err != nil {
		return nil, err
	}
	n := &nodes.ClusterStmt{
		Relation:  rv,
		Indexname: idxName,
		Loc:       nodes.Loc{Start: loc, End: p.prev.End},
	}
	if verbose {
		n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
	}
	return n, nil
}

// parseClusterIndexSpecification parses cluster_index_specification.
//
//	cluster_index_specification:
//	    USING name
//	    | /* EMPTY */
func (p *Parser) parseClusterIndexSpecification() (string, error) {
	if p.cur.Type != USING {
		return "", nil
	}
	p.advance()
	return p.parseName()
}

// parseReindexStmt parses a REINDEX statement.
//
//	ReindexStmt:
//	    REINDEX opt_reindex_option_list reindex_target_type opt_concurrently qualified_name
//	    | REINDEX opt_reindex_option_list SCHEMA opt_concurrently name
//	    | REINDEX opt_reindex_option_list reindex_target_all opt_concurrently opt_single_name
func (p *Parser) parseReindexStmt() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume REINDEX

	// opt_reindex_option_list: '(' utility_option_list ')' | /* EMPTY */
	var params *nodes.List
	if p.cur.Type == '(' {
		p.advance()
		var err error
		if params, err = p.parseUtilityOptionList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	switch p.cur.Type {
	case INDEX, TABLE:
		// reindex_target_type
		kind := nodes.REINDEX_OBJECT_INDEX
		if p.cur.Type == TABLE {
			kind = nodes.REINDEX_OBJECT_TABLE
		}
		p.advance()

		// opt_concurrently
		if p.cur.Type == CONCURRENTLY {
			p.advance()
			params = appendParamsList(params, makeDefElem("concurrently", nil))
		}

		nameStart := p.pos()
		names, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
		return &nodes.ReindexStmt{
			Kind:     kind,
			Relation: rv,
			Params:   params,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case SCHEMA:
		p.advance()

		// opt_concurrently
		if p.cur.Type == CONCURRENTLY {
			p.advance()
			params = appendParamsList(params, makeDefElem("concurrently", nil))
		}

		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.ReindexStmt{
			Kind:   nodes.REINDEX_OBJECT_SCHEMA,
			Name:   name,
			Params: params,
			Loc:    nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case SYSTEM_P, DATABASE:
		// reindex_target_multitable
		kind := nodes.REINDEX_OBJECT_SYSTEM
		if p.cur.Type == DATABASE {
			kind = nodes.REINDEX_OBJECT_DATABASE
		}
		p.advance()

		// opt_concurrently
		if p.cur.Type == CONCURRENTLY {
			p.advance()
			params = appendParamsList(params, makeDefElem("concurrently", nil))
		}

		// opt_single_name: ColId | /* EMPTY */
		name := ""
		if p.isColId() {
			var err error
			if name, err = p.parseColId(); err != nil {
				return nil, err
			}
		}
		return &nodes.ReindexStmt{
			Kind:   kind,
			Name:   name,
			Params: params,
			Loc:    nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// appendParamsList appends a node to a params list, creating the list if nil.
func appendParamsList(list *nodes.List, item nodes.Node) *nodes.List {
	if list == nil {
		return &nodes.List{Items: []nodes.Node{item}}
	}
	list.Items = append(list.Items, item)
	return list
}

// parseUtilityOptionList parses utility_option_list, the parenthesized
// options of EXPLAIN, VACUUM, ANALYZE, CLUSTER and REINDEX.
//
//	utility_option_list: utility_option_elem | utility_option_list ',' utility_option_elem
//	utility_option_elem: utility_option_name utility_option_arg
func (p *Parser) parseUtilityOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		name, err := p.parseUtilityOptionName()
		if err != nil {
			return nil, err
		}
		arg, err := p.parseUtilityOptionArg()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.DefElem{Defname: name, Arg: arg, Loc: nodes.NoLoc()})
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseUtilityOptionName parses utility_option_name.
//
//	utility_option_name: NonReservedWord | analyze_keyword | FORMAT_LA
//
// FORMAT is an unreserved keyword, so NonReservedWord reads it.
func (p *Parser) parseUtilityOptionName() (string, error) {
	if p.cur.Type == ANALYZE || p.cur.Type == ANALYSE {
		p.advance()
		return "analyze", nil
	}
	return p.parseNonReservedWord()
}

// parseUtilityOptionArg parses utility_option_arg, the argument of an
// option, which may be missing.
//
//	utility_option_arg: opt_boolean_or_string | NumericOnly | /* EMPTY */
func (p *Parser) parseUtilityOptionArg() (nodes.Node, error) {
	switch p.cur.Type {
	case ICONST, FCONST, '+', '-':
		return p.parseNumericOnlyRequired()
	case TRUE_P, FALSE_P, ON, SCONST:
		s, err := p.parseBooleanOrString()
		if err != nil {
			return nil, err
		}
		return &nodes.String{Str: s}, nil
	default:
		if p.isNonReservedWord() {
			s := p.cur.Str
			p.advance()
			return &nodes.String{Str: s}, nil
		}
		return nil, nil
	}
}

// parseBooleanOrString parses opt_boolean_or_string and fails where none
// is.
//
//	opt_boolean_or_string: TRUE_P | FALSE_P | ON | NonReservedWord_or_Sconst
func (p *Parser) parseBooleanOrString() (string, error) {
	switch p.cur.Type {
	case TRUE_P:
		p.advance()
		return "true", nil
	case FALSE_P:
		p.advance()
		return "false", nil
	case ON:
		p.advance()
		return "on", nil
	case SCONST:
		s := p.cur.Str
		p.advance()
		return s, nil
	}
	return p.parseNonReservedWord()
}

// parseOptBooleanOrString parses opt_boolean_or_string.
func (p *Parser) parseOptBooleanOrString() string {
	switch p.cur.Type {
	case TRUE_P:
		p.advance()
		return "true"
	case FALSE_P:
		p.advance()
		return "false"
	case ON:
		p.advance()
		return "on"
	default:
		return p.parseNonReservedWordOrSconst()
	}
}
