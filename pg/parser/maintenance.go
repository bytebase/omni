package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseVacuumStmt parses a VACUUM statement.
//
//	VacuumStmt:
//	    VACUUM opt_full opt_freeze opt_verbose opt_analyze opt_vacuum_relation_list
//	    | VACUUM '(' utility_option_list ')' opt_vacuum_relation_list
func (p *Parser) parseVacuumStmt() (nodes.Node, error) {
	p.advance() // consume VACUUM

	// VACUUM '(' utility_option_list ')' opt_vacuum_relation_list
	if p.cur.Type == '(' {
		p.advance() // consume '('
		opts := p.parseUtilityOptionList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		rels := p.parseOptVacuumRelationList()
		return &nodes.VacuumStmt{
			Options:     opts,
			Rels:        rels,
			IsVacuumCmd: true,
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

	n.Rels = p.parseOptVacuumRelationList()
	return n, nil
}

// parseAnalyzeStmt parses an ANALYZE/ANALYSE statement.
//
//	AnalyzeStmt:
//	    analyze_keyword opt_verbose opt_vacuum_relation_list
//	    | analyze_keyword '(' utility_option_list ')' opt_vacuum_relation_list
func (p *Parser) parseAnalyzeStmt() (nodes.Node, error) {
	p.advance() // consume ANALYZE/ANALYSE

	// analyze_keyword '(' utility_option_list ')' opt_vacuum_relation_list
	if p.cur.Type == '(' {
		p.advance() // consume '('
		opts := p.parseUtilityOptionList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		rels := p.parseOptVacuumRelationList()
		return &nodes.VacuumStmt{
			Options:     opts,
			Rels:        rels,
			IsVacuumCmd: false,
		}, nil
	}

	// analyze_keyword opt_verbose opt_vacuum_relation_list
	n := &nodes.VacuumStmt{
		IsVacuumCmd: false,
	}
	if p.cur.Type == VERBOSE {
		p.advance()
		n.Options = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
	}

	n.Rels = p.parseOptVacuumRelationList()
	return n, nil
}

// parseOptVacuumRelationList parses opt_vacuum_relation_list.
//
//	opt_vacuum_relation_list:
//	    vacuum_relation_list
//	    | /* EMPTY */
func (p *Parser) parseOptVacuumRelationList() *nodes.List {
	// Check if current token could start a vacuum_relation (i.e. a qualified_name).
	if !p.isColId() {
		return nil
	}
	return p.parseVacuumRelationList()
}

// parseVacuumRelationList parses vacuum_relation_list.
//
//	vacuum_relation_list:
//	    vacuum_relation
//	    | vacuum_relation_list ',' vacuum_relation
func (p *Parser) parseVacuumRelationList() *nodes.List {
	rel := p.parseVacuumRelation()
	if rel == nil {
		return nil
	}
	items := []nodes.Node{rel}
	for p.cur.Type == ',' {
		p.advance()
		r := p.parseVacuumRelation()
		if r != nil {
			items = append(items, r)
		}
	}
	return &nodes.List{Items: items}
}

// parseVacuumRelation parses vacuum_relation.
//
//	vacuum_relation:
//	    qualified_name opt_column_list
func (p *Parser) parseVacuumRelation() *nodes.VacuumRelation {
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil
	}
	rv := makeRangeVarFromNames(names)
	cols, _ := p.parseOptColumnList()
	return &nodes.VacuumRelation{
		Relation: rv,
		VaCols:   cols,
	}
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
	p.advance() // consume CLUSTER

	// CLUSTER '(' utility_option_list ')' [qualified_name cluster_index_specification]
	if p.cur.Type == '(' {
		p.advance() // consume '('
		params := p.parseUtilityOptionList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		// Check if a qualified_name follows
		if p.isColId() {
			names, err := p.parseQualifiedName()
			if err != nil {
				return &nodes.ClusterStmt{Params: params}, nil
			}
			rv := makeRangeVarFromNames(names)
			idxName := p.parseClusterIndexSpecification()
			return &nodes.ClusterStmt{
				Relation:  rv,
				Indexname: idxName,
				Params:    params,
			}, nil
		}
		return &nodes.ClusterStmt{Params: params}, nil
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
		n := &nodes.ClusterStmt{}
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

	names, err := p.parseQualifiedName()
	if err != nil {
		n := &nodes.ClusterStmt{}
		if verbose {
			n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
		}
		return n, nil
	}

	if p.cur.Type == ON {
		// Backwards-compatible: CLUSTER [VERBOSE] name ON qualified_name
		// The "name" we parsed is actually the index name (single name, not qualified).
		p.advance() // consume ON
		indexName := ""
		if names != nil && len(names.Items) == 1 {
			indexName = names.Items[0].(*nodes.String).Str
		}
		tableNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		rv := makeRangeVarFromNames(tableNames)
		n := &nodes.ClusterStmt{
			Relation:  rv,
			Indexname: indexName,
		}
		if verbose {
			n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
		}
		return n, nil
	}

	// CLUSTER [VERBOSE] qualified_name [USING name]
	rv := makeRangeVarFromNames(names)
	idxName := p.parseClusterIndexSpecification()
	n := &nodes.ClusterStmt{
		Relation:  rv,
		Indexname: idxName,
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
func (p *Parser) parseClusterIndexSpecification() string {
	if p.cur.Type == USING {
		p.advance()
		name, _ := p.parseName()
		return name
	}
	return ""
}

// parseReindexStmt parses a REINDEX statement.
//
//	ReindexStmt:
//	    REINDEX opt_reindex_option_list reindex_target_type opt_concurrently qualified_name
//	    | REINDEX opt_reindex_option_list SCHEMA opt_concurrently name
//	    | REINDEX opt_reindex_option_list reindex_target_multitable opt_concurrently name
func (p *Parser) parseReindexStmt() (nodes.Node, error) {
	p.advance() // consume REINDEX

	// opt_reindex_option_list: '(' utility_option_list ')' | /* EMPTY */
	var params *nodes.List
	if p.cur.Type == '(' {
		p.advance()
		params = p.parseUtilityOptionList()
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

		names, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		rv := makeRangeVarFromNames(names)
		return &nodes.ReindexStmt{
			Kind:     kind,
			Relation: rv,
			Params:   params,
		}, nil

	case SCHEMA:
		p.advance()

		// opt_concurrently
		if p.cur.Type == CONCURRENTLY {
			p.advance()
			params = appendParamsList(params, makeDefElem("concurrently", nil))
		}

		name, _ := p.parseName()
		return &nodes.ReindexStmt{
			Kind:   nodes.REINDEX_OBJECT_SCHEMA,
			Name:   name,
			Params: params,
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

		name, _ := p.parseName()
		return &nodes.ReindexStmt{
			Kind:   kind,
			Name:   name,
			Params: params,
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

// parseUtilityOptionList parses utility_option_list.
func (p *Parser) parseUtilityOptionList() *nodes.List {
	elem := p.parseUtilityOptionElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseUtilityOptionElem())
	}
	return &nodes.List{Items: items}
}

// parseUtilityOptionElem parses utility_option_elem.
func (p *Parser) parseUtilityOptionElem() *nodes.DefElem {
	name := p.parseUtilityOptionName()
	arg := p.parseUtilityOptionArg()
	return &nodes.DefElem{Defname: name, Arg: arg, Loc: nodes.NoLoc()}
}

// parseUtilityOptionName parses utility_option_name.
func (p *Parser) parseUtilityOptionName() string {
	switch p.cur.Type {
	case ANALYZE, ANALYSE:
		p.advance()
		return "analyze"
	case FORMAT:
		p.advance()
		return "format"
	case DEFAULT:
		p.advance()
		return "default"
	default:
		s := p.cur.Str
		p.advance()
		return s
	}
}

// parseUtilityOptionArg parses utility_option_arg.
func (p *Parser) parseUtilityOptionArg() nodes.Node {
	switch p.cur.Type {
	case '*':
		starLoc := p.pos()
		p.advance()
		return &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}
	case DEFAULT:
		p.advance()
		return &nodes.String{Str: "default"}
	case '(':
		p.advance()
		list := p.parseUtilityOptionArgList()
		p.expect(')')
		return list
	case ICONST, FCONST:
		return p.parseNumericOnly()
	case '+', '-':
		return p.parseNumericOnly()
	case TRUE_P:
		p.advance()
		return &nodes.String{Str: "true"}
	case FALSE_P:
		p.advance()
		return &nodes.String{Str: "false"}
	case ON:
		p.advance()
		return &nodes.String{Str: "on"}
	case SCONST:
		s := p.cur.Str
		p.advance()
		return &nodes.String{Str: s}
	default:
		if p.isNonReservedWord() {
			s := p.cur.Str
			p.advance()
			return &nodes.String{Str: s}
		}
		return nil
	}
}

// parseUtilityOptionArgList parses utility_option_arg_list.
func (p *Parser) parseUtilityOptionArgList() *nodes.List {
	s := p.parseOptBooleanOrString()
	items := []nodes.Node{&nodes.String{Str: s}}
	for p.cur.Type == ',' {
		p.advance()
		s = p.parseOptBooleanOrString()
		items = append(items, &nodes.String{Str: s})
	}
	return &nodes.List{Items: items}
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
