package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseSelectStmt parses a SELECT statement.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	select:
//	    [ with_clause ] subquery [ for_update_clause ] ;
//
//	subquery:
//	    query_block
//	    [ { UNION [ ALL ] | INTERSECT | MINUS } query_block ]...
//	    [ order_by_clause ]
//	    [ row_limiting_clause ]
//
//	query_block:
//	    SELECT [ hint ] [ DISTINCT | UNIQUE | ALL ]
//	    select_list
//	    [ FROM { table_reference | join_clause } [, { table_reference | join_clause } ]... ]
//	    [ where_clause ]
//	    [ hierarchical_query_clause ]
//	    [ group_by_clause ]
//	    [ model_clause ]
//	    [ window_clause ]
//	    [ qualify_clause ]
//
//	order_by_clause:
//	    ORDER [ SIBLINGS ] BY
//	    { expr | position | c_alias } [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ]
//	    [, { expr | position | c_alias } [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ] ]...
//
//	row_limiting_clause:
//	    [ OFFSET offset { ROW | ROWS } ]
//	    [ FETCH { FIRST | NEXT } [ { rowcount | percent PERCENT } ] { ROW | ROWS }
//	      { ONLY | WITH TIES } ]
//
//	for_update_clause:
//	    FOR UPDATE
//	    [ OF [ [ schema. ] { table | view } . ] column
//	      [, [ [ schema. ] { table | view } . ] column ]... ]
//	    [ NOWAIT | WAIT integer | SKIP LOCKED ]
//
//	window_clause:
//	    WINDOW window_name AS ( window_specification )
//	    [, window_name AS ( window_specification ) ]...
//
//	qualify_clause:
//	    QUALIFY condition
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error) {
	start := p.pos()
	sel := &nodes.SelectStmt{
		TargetList: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// WITH clause
	if p.cur.Type == kwWITH {
		var parseErr941 error
		sel.WithClause, parseErr941 = p.parseWithClause()
		if parseErr941 != nil {
			return nil, parseErr941
		}
	}

	if p.cur.Type != kwSELECT {
		sel.Loc.End = p.prev.End
		return sel, nil
	}
	p.advance() // consume SELECT

	// Hints
	if p.cur.Type == tokHINT {
		tok := p.cur
		sel.Hints = &nodes.List{}
		sel.Hints.Items = append(sel.Hints.Items, &nodes.Hint{
			Text: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: tok.End},
		})
		p.advance()
	}

	// ALL | DISTINCT | UNIQUE
	switch p.cur.Type {
	case kwALL:
		sel.All = true
		p.advance()
	case kwDISTINCT:
		sel.Distinct = true
		p.advance()
	case kwUNIQUE:
		sel.UniqueKw = true
		sel.Distinct = true
		p.advance()
	}
	var parseErr942 error

	// Select list
	sel.TargetList, parseErr942 = p.parseSelectList()
	if parseErr942 != nil {
		return nil, parseErr942
	}
	if sel.TargetList.Len() == 0 {

		return sel, p.syntaxErrorAtCur(

		// INTO (PL/SQL SELECT ... INTO variable_list FROM ...)
		)
	}

	if p.cur.Type == kwINTO {
		p.advance()
		var // consume INTO
		parseErr943 error
		sel.IntoVars, parseErr943 = p.parseExprList()
		if parseErr943 !=

			// FROM
			nil {
			return nil, parseErr943
		}
	}

	if p.cur.Type == kwFROM {
		p.advance()
		var parseErr944 error
		sel.FromClause, parseErr944 = p.parseFromClause()
		if parseErr944 !=

			// PIVOT / UNPIVOT (parsed after FROM, before WHERE)
			nil {
			return nil, parseErr944
		}
	}

	if p.cur.Type == kwPIVOT {
		var parseErr945 error
		sel.Pivot, parseErr945 = p.parsePivotClause()
		if parseErr945 != nil {
			return nil, parseErr945
		}
	} else if p.cur.Type == kwUNPIVOT {
		var parseErr946 error
		sel.Unpivot, parseErr946 = p.parseUnpivotClause()
		if parseErr946 !=

			// WHERE
			nil {
			return nil, parseErr946
		}
	}

	if p.cur.Type == kwWHERE {
		p.advance()
		var parseErr947 error
		sel.WhereClause, parseErr947 = p.parseExpr()
		if parseErr947 != nil {
			return nil, parseErr947
		}
		if sel.WhereClause == nil {

			return sel, p.syntaxErrorAtCur(

			// START WITH / CONNECT BY (either order)
			)
		}
	}

	if p.cur.Type == kwSTART || p.cur.Type == kwCONNECT {
		var parseErr948 error
		sel.Hierarchical, parseErr948 = p.parseHierarchicalClause()
		if parseErr948 !=

			// GROUP BY
			nil {
			return nil, parseErr948
		}
	}

	if p.cur.Type == kwGROUP {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		var parseErr949 error
		sel.GroupClause, parseErr949 = p.parseGroupByList()
		if parseErr949 != nil {
			return nil, parseErr949
		}
		if sel.GroupClause.Len() == 0 {

			return sel, p.syntaxErrorAtCur(

			// HAVING
			)
		}
	}

	if p.cur.Type == kwHAVING {
		p.advance()
		var parseErr950 error
		sel.HavingClause, parseErr950 = p.parseExpr()
		if parseErr950 != nil {
			return nil, parseErr950
		}
		if sel.HavingClause == nil {

			return sel, p.syntaxErrorAtCur(

			// MODEL clause
			)
		}
	}

	if p.cur.Type == kwMODEL {
		var parseErr951 error
		sel.ModelClause, parseErr951 = p.parseModelClause()
		if parseErr951 !=

			// WINDOW clause
			nil {
			return nil, parseErr951
		}
	}

	if p.isIdentLikeStr("WINDOW") {
		var parseErr952 error
		sel.WindowDefs, parseErr952 = p.parseWindowClause()
		if parseErr952 !=

			// QUALIFY clause
			nil {
			return nil, parseErr952
		}
	}

	if p.isIdentLikeStr("QUALIFY") {
		p.advance()
		var parseErr953 error
		sel.QualifyClause, parseErr953 = p.parseExpr()
		if parseErr953 !=

			// ORDER [SIBLINGS] BY
			nil {
			return nil, parseErr953
		}
	}

	if p.cur.Type == kwORDER {
		p.advance()
		if p.isIdentLikeStr("SIBLINGS") {
			sel.SiblingsOrder = true
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
		}
		var parseErr954 error
		sel.OrderBy, parseErr954 = p.parseOrderByList()
		if parseErr954 != nil {
			return nil, parseErr954
		}
		if sel.OrderBy.Len() == 0 {

			return sel, p.syntaxErrorAtCur(

			// FOR UPDATE
			)
		}
	}

	if p.cur.Type == kwFOR {
		var parseErr955 error
		sel.ForUpdate, parseErr955 = p.parseForUpdateClause()
		if parseErr955 !=

			// OFFSET / FETCH FIRST
			nil {
			return nil, parseErr955
		}
	}

	if p.cur.Type == kwOFFSET || p.cur.Type == kwFETCH {
		var parseErr956 error
		sel.FetchFirst, parseErr956 = p.parseFetchFirstClause()
		if parseErr956 !=

			// Set operations: UNION, INTERSECT, MINUS
			nil {
			return nil, parseErr956
		}
	}

	switch p.cur.Type {
	case kwUNION:
		p.advance()
		sel.Op = nodes.SETOP_UNION
		if p.cur.Type == kwALL {
			sel.SetAll = true
			p.advance()
		}
		if p.cur.Type == tokEOF {

			return sel, p.syntaxErrorAtCur()
		}
		var parseErr957 error
		sel.Rarg, parseErr957 = p.parseSelectStmt()
		if parseErr957 != nil {
			return nil, parseErr957
		}
	case kwINTERSECT:
		p.advance()
		sel.Op = nodes.SETOP_INTERSECT
		if p.cur.Type == kwALL {
			sel.SetAll = true
			p.advance()
		}
		if p.cur.Type == tokEOF {

			return sel, p.syntaxErrorAtCur()
		}
		var parseErr958 error
		sel.Rarg, parseErr958 = p.parseSelectStmt()
		if parseErr958 != nil {
			return nil, parseErr958
		}
	case kwMINUS:
		p.advance()
		sel.Op = nodes.SETOP_MINUS
		if p.cur.Type == tokEOF {

			return sel, p.syntaxErrorAtCur()
		}
		var parseErr959 error
		sel.Rarg, parseErr959 = p.parseSelectStmt()
		if parseErr959 != nil {
			return nil, parseErr959
		}
	}

	sel.Loc.End = p.prev.End
	return sel, nil
}

// parseSelectList parses the select list (target expressions).
func (p *Parser) parseSelectList() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		if p.isSelectListTerminator() {
			break
		}
		rt, parseErr960 := p.parseResTarget()
		if parseErr960 != nil {
			return nil, parseErr960
		}
		if rt != nil {
			list.Items = append(list.Items, rt)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	return list, nil
}

func (p *Parser) isSelectListTerminator() bool {
	switch p.cur.Type {
	case tokEOF, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION,
		kwINTERSECT, kwMINUS, kwFOR, kwCONNECT, kwSTART, kwFETCH, kwOFFSET:
		return true
	default:
		return false
	}
}

// parseResTarget parses a single target expression (with optional alias).
func (p *Parser) parseResTarget() (*nodes.ResTarget, error) {
	start := p.pos()
	expr, parseErr961 := p.parseExpr()
	if parseErr961 != nil {
		return nil, parseErr961
	}
	if expr == nil {
		return nil, nil
	}

	rt := &nodes.ResTarget{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// Optional alias: AS name or just name
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr962 error
		rt.Name, parseErr962 = p.parseIdentifier()
		if parseErr962 != nil {
			return nil, parseErr962
		}
	} else if p.isAliasCandidate() {
		var parseErr963 error
		rt.Name, parseErr963 = p.parseIdentifier()
		if parseErr963 != nil {
			return nil, parseErr963
		}
	}

	rt.Loc.End = p.prev.End
	return rt, nil
}

// isAliasCandidate returns true if the current token can be an implicit alias.
// We exclude keywords that start clauses to avoid consuming FROM, WHERE, etc. as aliases.
func (p *Parser) isAliasCandidate() bool {
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	// Disallow clause-starting keywords as implicit aliases
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT,
		kwMINUS, kwFOR, kwCONNECT, kwSTART, kwFETCH, kwOFFSET, kwON,
		kwLEFT, kwRIGHT, kwINNER, kwOUTER, kwCROSS, kwFULL, kwNATURAL,
		kwJOIN, kwWHEN, kwTHEN, kwELSE, kwEND, kwAND, kwOR, kwNOT,
		kwIS, kwIN, kwBETWEEN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4,
		kwINTO, kwVALUES, kwSET, kwRETURNING, kwPIVOT, kwUNPIVOT,
		kwMODEL, kwWITH, kwKEEP, kwOVER:
		return false
	}
	return false
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		expr, parseErr964 := p.parseExpr()
		if parseErr964 != nil {
			return nil, parseErr964
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()

		}
		list.Items = append(list.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()

		}
	}
	return list, nil
}

// parseGroupByList parses a comma-separated GROUP BY list, handling
// GROUPING SETS, CUBE, and ROLLUP extensions.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	group_by_clause ::=
//	    GROUP BY { expr | rollup_cube | grouping_sets } [, ...]
//	rollup_cube ::= { ROLLUP | CUBE } ( expr [, ...] )
//	grouping_sets ::= GROUPING SETS ( { rollup_cube | expr } [, ...] )
func (p *Parser) parseGroupByList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		item, parseErr965 := p.parseGroupByItem()
		if parseErr965 != nil {
			return nil, parseErr965
		}
		if item == nil {
			return nil, p.syntaxErrorAtCur()

		}
		list.Items = append(list.Items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()

		}
	}
	return list, nil
}

// parseGroupByItem parses a single GROUP BY item: expression, ROLLUP(...), CUBE(...),
// or GROUPING SETS(...).
func (p *Parser) parseGroupByItem() (nodes.Node, error) {
	start := p.pos()

	switch p.cur.Type {
	case kwROLLUP:
		p.advance() // consume ROLLUP
		rc := &nodes.RollupClause{Loc: nodes.Loc{Start: start}}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr966 error
			rc.Args, parseErr966 = p.parseExprList()
			if parseErr966 != nil {
				return nil, parseErr966
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		rc.Loc.End = p.prev.End
		return rc, nil

	case kwCUBE:
		p.advance() // consume CUBE
		cc := &nodes.CubeClause{Loc: nodes.Loc{Start: start}}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr967 error
			cc.Args, parseErr967 = p.parseExprList()
			if parseErr967 != nil {
				return nil, parseErr967
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		cc.Loc.End = p.prev.End
		return cc, nil

	case kwGROUPING:
		// GROUPING SETS(...)
		if p.peekNext().Type == kwSETS {
			p.advance() // consume GROUPING
			p.advance() // consume SETS
			gs := &nodes.GroupingSetsClause{Loc: nodes.Loc{Start: start}}
			if p.cur.Type == '(' {
				p.advance()
				gs.Sets = &nodes.List{}
				for {
					item, parseErr968 := p.parseGroupByItem()
					if parseErr968 != nil {
						return nil, parseErr968
					}
					if item == nil {
						break
					}
					gs.Sets.Items = append(gs.Sets.Items, item)
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			gs.Loc.End = p.prev.End
			return gs, nil
		}
		// GROUPING(expr) is a function call, fall through to parseExpr
		return p.parseExpr()

	default:
		return p.parseExpr()
	}
}

// parseFromClause parses a FROM clause (comma-separated table references).
func (p *Parser) parseFromClause() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		tref, parseErr969 := p.parseTableRef()
		if parseErr969 != nil {
			return nil, parseErr969

			// Check for JOINs
		}
		if tref == nil {
			break
		}
		var parseErr970 error

		tref, parseErr970 = p.parseJoinContinuation(tref)
		if parseErr970 != nil {
			return nil, parseErr970
		}

		list.Items = append(list.Items, tref)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseTableRef parses a single table reference.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	table_reference:
//	    { query_table_expression | ( join_clause ) | inline_external_table } [ t_alias ]
//
//	query_table_expression:
//	    { [ schema. ] { table | view | materialized_view | analytic_view | hierarchy }
//	      [ partition_extension_clause ] [ @ dblink ]
//	    | ( subquery [ subquery_restriction_clause ] )
//	    | table_collection_expression
//	    | inline_analytic_view
//	    }
//	    [ flashback_query_clause ] [ sample_clause ]
//	    [ pivot_clause | unpivot_clause ] [ row_pattern_clause ]
//	    [ containers_clause ] [ shards_clause ]
func (p *Parser) parseTableRef() (nodes.TableExpr, error) {
	start := p.pos()

	// LATERAL ( subquery )
	if p.cur.Type == kwLATERAL {
		return p.parseLateralRef(start)
	}

	// XMLTABLE(...)
	if p.cur.Type == kwXMLTABLE {
		return p.parseXmlTableRef(start)
	}

	// JSON_TABLE(...)
	if p.cur.Type == kwJSON_TABLE {
		return p.parseJsonTableRef(start)
	}

	// TABLE(collection_expression) — table collection expression
	if p.cur.Type == kwTABLE && p.peekNext().Type == '(' {
		return p.parseTableCollectionExpr(start)
	}

	// CONTAINERS(table) or SHARDS(table)
	if p.isIdentLikeStr("CONTAINERS") && p.peekNext().Type == '(' {
		return p.parseContainersOrShards(start, false)
	}
	if p.isIdentLikeStr("SHARDS") && p.peekNext().Type == '(' {
		return p.parseContainersOrShards(start, true)
	}

	// Subquery: ( SELECT ... )
	if p.cur.Type == '(' {
		return p.parseSubqueryRef(start)
	}

	// MATCH_RECOGNIZE as standalone (rare, usually post-table)
	if p.isIdentLikeStr("MATCH_RECOGNIZE") {
		return p.parseMatchRecognize(start)
	}

	// Table name
	if !p.isIdentLike() {
		return nil, nil
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}

	name, parseErr971 := p.parseReservedCheckedObjectName()
	if parseErr971 != nil {
		return nil, parseErr971
	}
	tr := &nodes.TableRef{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// Partition extension clause: PARTITION (name) | PARTITION FOR (key) | SUBPARTITION ...
	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		var parseErr972 error
		tr.PartitionExt, parseErr972 = p.parsePartitionExtClause()
		if parseErr972 !=

			// @ dblink
			nil {
			return nil, parseErr972
		}
	}

	if p.cur.Type == '@' {
		p.advance()
		var parseErr973 error
		tr.Dblink, parseErr973 = p.parseIdentifier()
		if parseErr973 !=

			// Flashback query: VERSIONS BETWEEN ... AND ... or AS OF SCN/TIMESTAMP
			nil {
			return nil, parseErr973
		}
	}

	if p.cur.Type == kwVERSIONS || (p.cur.Type == kwAS && p.peekNext().Type == kwOF) {
		var parseErr974 error
		tr.Flashback, parseErr974 = p.parseFlashbackClause()
		if parseErr974 !=

			// SAMPLE [BLOCK] (percent) [SEED (value)]
			nil {
			return nil, parseErr974
		}
	}

	if p.cur.Type == kwSAMPLE {
		var parseErr975 error
		tr.Sample, parseErr975 = p.parseSampleClause()
		if parseErr975 !=

			// MATCH_RECOGNIZE after table reference
			nil {
			return nil, parseErr975
		}
	}

	if p.isIdentLikeStr("MATCH_RECOGNIZE") {
		mrStart := p.pos()
		mrRef, parseErr976 := p.parseMatchRecognize(mrStart)
		if parseErr976 != nil {
			return nil, parseErr976
		}
		if mrClause, ok := mrRef.(*nodes.MatchRecognizeClause); ok {
			// Wrap: the table ref becomes the left side of a join-like construct
			// For simplicity, return the MATCH_RECOGNIZE with the table ref embedded
			_ = mrClause // MATCH_RECOGNIZE is returned directly
			return mrRef, nil
		}
	}

	// Optional alias
	if p.cur.Type == kwAS {
		// Only consume AS as alias intro if next is not OF (flashback already handled)
		p.advance()
		var parseErr977 error
		tr.Alias, parseErr977 = p.parseAlias()
		if parseErr977 != nil {
			return nil, parseErr977
		}
	} else if p.isTableAliasCandidate() {
		var parseErr978 error
		tr.Alias, parseErr978 = p.parseAlias()
		if parseErr978 != nil {
			return nil, parseErr978
		}
	}

	tr.Loc.End = p.prev.End
	return tr, nil
}

// isTableAliasCandidate checks if current token can be a table alias.
// Excludes soft keywords that start clauses in SELECT statements.
func (p *Parser) isTableAliasCandidate() bool {
	if p.cur.Type == tokQIDENT {
		return true
	}
	if p.cur.Type == tokIDENT {
		switch strings.ToUpper(p.cur.Str) {
		case "WINDOW", "QUALIFY", "MATCH_RECOGNIZE", "CONTAINERS", "SHARDS",
			"SIBLINGS", "APPLY", "SEARCH", "CYCLE", "VERSIONS", "PERIOD",
			"XML":
			return false
		}
		return true
	}
	return false
}

// parseSubqueryRef parses a subquery in FROM: ( SELECT ... ) alias.
func (p *Parser) parseSubqueryRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume '('

	subSel, parseErr979 := p.parseSelectStmt()
	if parseErr979 != nil {
		return nil, parseErr979
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	ref := &nodes.SubqueryRef{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start},
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr980 error
		ref.Alias, parseErr980 = p.parseAlias()
		if parseErr980 != nil {
			return nil, parseErr980
		}
	} else if p.isTableAliasCandidate() {
		var parseErr981 error
		ref.Alias, parseErr981 = p.parseAlias()
		if parseErr981 != nil {
			return nil, parseErr981
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseLateralRef parses a LATERAL inline view.
//
//	LATERAL ( subquery ) [ alias ]
func (p *Parser) parseLateralRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume LATERAL

	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance() // consume '('

	subSel, parseErr982 := p.parseSelectStmt()
	if parseErr982 != nil {
		return nil, parseErr982
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	ref := &nodes.LateralRef{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr983 error
		ref.Alias, parseErr983 = p.parseAlias()
		if parseErr983 != nil {
			return nil, parseErr983
		}
	} else if p.isTableAliasCandidate() {
		var parseErr984 error
		ref.Alias, parseErr984 = p.parseAlias()
		if parseErr984 != nil {
			return nil, parseErr984
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseXmlTableRef parses an XMLTABLE expression in FROM.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/XMLTABLE.html
//
//	XMLTABLE ( xpath_string PASSING xml_expr COLUMNS column_def [, ...] )
func (p *Parser) parseXmlTableRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume XMLTABLE

	ref := &nodes.XmlTableRef{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		return ref, nil
	}
	p.advance()
	var // consume '('
	parseErr985 error

	// XPath expression (usually a string literal)
	ref.XPath, parseErr985 = p.parseExpr()
	if parseErr985 !=

		// PASSING xml_expr
		nil {
		return nil, parseErr985
	}

	if p.cur.Type == kwPASSING {
		p.advance()
		var parseErr986 error
		ref.Passing, parseErr986 = p.parseExpr()
		if parseErr986 !=

			// COLUMNS column_def [, ...]
			nil {
			return nil, parseErr986
		}
	}

	if p.cur.Type == kwCOLUMNS {
		p.advance()
		var parseErr987 error
		ref.Columns, parseErr987 = p.parseXmlTableColumns()
		if parseErr987 != nil {
			return nil, parseErr987
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr988 error
		ref.Alias, parseErr988 = p.parseAlias()
		if parseErr988 != nil {
			return nil, parseErr988
		}
	} else if p.isTableAliasCandidate() {
		var parseErr989 error
		ref.Alias, parseErr989 = p.parseAlias()
		if parseErr989 != nil {
			return nil, parseErr989
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseXmlTableColumns parses column definitions in XMLTABLE.
func (p *Parser) parseXmlTableColumns() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr990 := p.parseXmlTableColumn()
		if parseErr990 != nil {
			return nil, parseErr990
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseXmlTableColumn parses a single XMLTABLE column definition.
//
//	name { datatype [PATH path] [DEFAULT default] | FOR ORDINALITY }
func (p *Parser) parseXmlTableColumn() (*nodes.XmlTableColumn, error) {
	if !p.isIdentLike() {
		return nil, nil
	}
	start := p.pos()
	col := &nodes.XmlTableColumn{Loc: nodes.Loc{Start: start}}
	var parseErr991 error
	col.Name, parseErr991 = p.parseIdentifier()
	if parseErr991 !=

		// FOR ORDINALITY
		nil {
		return nil, parseErr991
	}

	if p.cur.Type == kwFOR && p.peekNext().Type == kwORDINALITY {
		p.advance() // consume FOR
		p.advance() // consume ORDINALITY
		col.ForOrdinality = true
		col.Loc.End = p.prev.End
		return col, nil
	}
	var parseErr992 error

	// Data type
	col.TypeName, parseErr992 = p.parseTypeName()
	if parseErr992 !=

		// PATH
		nil {
		return nil, parseErr992
	}

	if p.cur.Type == kwPATH {
		p.advance()
		var parseErr993 error
		col.Path, parseErr993 = p.parseExpr()
		if parseErr993 !=

			// DEFAULT
			nil {
			return nil, parseErr993
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance()
		var parseErr994 error
		col.Default, parseErr994 = p.parseExpr()
		if parseErr994 != nil {
			return nil, parseErr994
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// parseJsonTableRef parses a JSON_TABLE expression in FROM.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_TABLE.html
//
//	JSON_TABLE ( expr, path_string COLUMNS ( column_def [, ...] ) )
func (p *Parser) parseJsonTableRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume JSON_TABLE

	ref := &nodes.JsonTableRef{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		return ref, nil
	}
	p.advance()
	var // consume '('
	parseErr995 error

	// JSON expression
	ref.Expr, parseErr995 = p.parseExpr()
	if parseErr995 !=

		// Comma separator
		nil {
		return nil, parseErr995
	}

	if p.cur.Type == ',' {
		p.advance()
	}
	var parseErr996 error

	// Path expression (string literal)
	ref.Path, parseErr996 = p.parseExpr()
	if parseErr996 !=

		// COLUMNS ( column_def [, ...] )
		nil {
		return nil, parseErr996
	}

	if p.cur.Type == kwCOLUMNS {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr997 error
			ref.Columns, parseErr997 = p.parseJsonTableColumns()
			if parseErr997 != nil {
				return nil, parseErr997
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr998 error
		ref.Alias, parseErr998 = p.parseAlias()
		if parseErr998 != nil {
			return nil, parseErr998
		}
	} else if p.isTableAliasCandidate() {
		var parseErr999 error
		ref.Alias, parseErr999 = p.parseAlias()
		if parseErr999 != nil {
			return nil, parseErr999
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseJsonTableColumns parses column definitions in JSON_TABLE.
func (p *Parser) parseJsonTableColumns() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr1000 := p.parseJsonTableColumn()
		if parseErr1000 != nil {
			return nil, parseErr1000
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseJsonTableColumn parses a single JSON_TABLE column definition.
//
//	name datatype [PATH path] | name FOR ORDINALITY | NESTED [PATH] path COLUMNS ( ... )
func (p *Parser) parseJsonTableColumn() (*nodes.JsonTableColumn, error) {
	start := p.pos()

	// NESTED [PATH] path COLUMNS (...)
	if p.cur.Type == kwNESTED {
		p.advance() // consume NESTED
		col := &nodes.JsonTableColumn{Loc: nodes.Loc{Start: start}}
		if p.cur.Type == kwPATH {
			p.advance() // consume PATH
		}
		nested := &nodes.JsonTableRef{Loc: nodes.Loc{Start: start}}
		var parseErr1001 error
		nested.Path, parseErr1001 = p.parseExpr()
		if parseErr1001 != nil {
			return nil, parseErr1001
		}
		if p.cur.Type == kwCOLUMNS {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				var parseErr1002 error
				nested.Columns, parseErr1002 = p.parseJsonTableColumns()
				if parseErr1002 != nil {
					return nil, parseErr1002
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
		nested.Loc.End = p.prev.End
		col.Nested = nested
		col.Loc.End = p.prev.End
		return col, nil
	}

	if !p.isIdentLike() {
		return nil, nil
	}

	col := &nodes.JsonTableColumn{Loc: nodes.Loc{Start: start}}
	var parseErr1003 error
	col.Name, parseErr1003 = p.parseIdentifier()
	if parseErr1003 !=

		// FOR ORDINALITY
		nil {
		return nil, parseErr1003
	}

	if p.cur.Type == kwFOR && p.peekNext().Type == kwORDINALITY {
		p.advance() // consume FOR
		p.advance() // consume ORDINALITY
		col.ForOrdinality = true
		col.Loc.End = p.prev.End
		return col, nil
	}
	var parseErr1004 error

	// Data type
	col.TypeName, parseErr1004 = p.parseTypeName()
	if parseErr1004 !=

		// EXISTS
		nil {
		return nil, parseErr1004
	}

	if p.cur.Type == kwEXISTS {
		col.Exists = true
		p.advance()
	}

	// PATH
	if p.cur.Type == kwPATH {
		p.advance()
		var parseErr1005 error
		col.Path, parseErr1005 = p.parseExpr()
		if parseErr1005 != nil {
			return nil, parseErr1005
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// parseJoinContinuation parses any JOIN clauses that follow a table reference.
func (p *Parser) parseJoinContinuation(left nodes.TableExpr) (nodes.TableExpr, error) {
	for {
		joinStart := p.pos()
		jt, ok := p.matchJoinType()
		if !ok {
			break
		}

		right, parseErr1006 := p.parseTableRef()
		if parseErr1006 != nil {
			return nil, parseErr1006
		}
		if right == nil {

			return left, p.syntaxErrorAtCur()
		}

		jc := &nodes.JoinClause{
			Type:  jt,
			Left:  left,
			Right: right,
			Loc:   nodes.Loc{Start: joinStart},
		}

		// ON condition
		if p.cur.Type == kwON {
			p.advance()
			var parseErr1007 error
			jc.On, parseErr1007 = p.parseExpr()
			if parseErr1007 != nil {
				return nil, parseErr1007
			}
			if jc.On == nil {
				return nil, p.syntaxErrorAtCur()

				jc.Loc.End = p.prev.End
				return jc, nil
			}
		}

		// USING ( col1, col2, ... )
		if p.cur.Type == kwUSING {
			p.advance()
			if p.cur.Type != '(' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			jc.Using = &nodes.List{}
			for {
				name, parseErr1008 := p.parseIdentifier()
				if parseErr1008 != nil {
					return nil, parseErr1008
				}
				if name == "" {
					return nil, p.syntaxErrorAtCur()
				}
				jc.Using.Items = append(jc.Using.Items, &nodes.String{Str: name})
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type != ')' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		}

		jc.Loc.End = p.prev.End
		left = jc
	}
	return left, nil
}

// matchJoinType tries to match a JOIN keyword sequence.
// Returns the JoinType and true if matched, or false if no JOIN found.
func (p *Parser) matchJoinType() (nodes.JoinType, bool) {
	natural := false
	if p.cur.Type == kwNATURAL {
		natural = true
		p.advance()
	}

	switch p.cur.Type {
	case kwJOIN:
		p.advance()
		if natural {
			return nodes.JOIN_NATURAL_INNER, true
		}
		return nodes.JOIN_INNER, true

	case kwINNER:
		p.advance()
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_INNER, true
		}
		return nodes.JOIN_INNER, true

	case kwLEFT:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_LEFT, true
		}
		return nodes.JOIN_LEFT, true

	case kwRIGHT:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_RIGHT, true
		}
		return nodes.JOIN_RIGHT, true

	case kwFULL:
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		if natural {
			return nodes.JOIN_NATURAL_FULL, true
		}
		return nodes.JOIN_FULL, true

	case kwCROSS:
		p.advance()
		// CROSS APPLY
		if p.isIdentLikeStr("APPLY") {
			p.advance()
			return nodes.JOIN_CROSS_APPLY, true
		}
		if p.cur.Type == kwJOIN {
			p.advance()
		}
		return nodes.JOIN_CROSS, true

	case kwOUTER:
		// OUTER APPLY
		if p.peekNext().Type == tokIDENT || p.isIdentLikeStr("APPLY") {
			p.advance() // consume OUTER
			if p.isIdentLikeStr("APPLY") {
				p.advance()
				return nodes.JOIN_OUTER_APPLY, true
			}
		}
		// OUTER JOIN handled by LEFT/RIGHT/FULL above
		return 0, false
	}

	if natural {
		// NATURAL without a recognized join keyword — treat as NATURAL INNER JOIN
		return nodes.JOIN_NATURAL_INNER, true
	}

	return 0, false
}

// parseHierarchicalClause parses START WITH / CONNECT BY clauses (either order).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Hierarchical-Queries.html
//
//	[ START WITH condition ] CONNECT BY [ NOCYCLE ] condition
//	CONNECT BY [ NOCYCLE ] condition [ START WITH condition ]
func (p *Parser) parseHierarchicalClause() (*nodes.HierarchicalClause, error) {
	start := p.pos()
	hc := &nodes.HierarchicalClause{
		Loc: nodes.Loc{Start: start},
	}

	// START WITH first, then CONNECT BY
	if p.cur.Type == kwSTART {
		p.advance() // START
		if p.cur.Type == kwWITH {
			p.advance() // WITH
		}
		var parseErr1009 error
		hc.StartWith, parseErr1009 = p.parseExpr()
		if parseErr1009 != nil {
			return nil, parseErr1009
		}
	}

	if p.cur.Type == kwCONNECT {
		p.advance() // CONNECT
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.isIdentLikeStr("NOCYCLE") {
			hc.IsNocycle = true
			p.advance()
		}
		var parseErr1010 error
		hc.ConnectBy, parseErr1010 = p.parseExpr()
		if parseErr1010 !=

			// START WITH may come after CONNECT BY
			nil {
			return nil, parseErr1010
		}
	}

	if hc.StartWith == nil && p.cur.Type == kwSTART {
		p.advance()
		if p.cur.Type == kwWITH {
			p.advance()
		}
		var parseErr1011 error
		hc.StartWith, parseErr1011 = p.parseExpr()
		if parseErr1011 != nil {
			return nil, parseErr1011
		}
	}

	hc.Loc.End = p.prev.End
	return hc, nil
}

// parseOrderByList parses a comma-separated list of ORDER BY items.
func (p *Parser) parseOrderByList() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		sb, parseErr1012 := p.parseSortBy()
		if parseErr1012 != nil {
			return nil, parseErr1012
		}
		if sb == nil {
			break
		}
		list.Items = append(list.Items, sb)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseSortBy parses a single ORDER BY item.
//
//	sort_key [ ASC | DESC ] [ NULLS { FIRST | LAST } ]
func (p *Parser) parseSortBy() (*nodes.SortBy, error) {
	start := p.pos()
	expr, parseErr1013 := p.parseExpr()
	if parseErr1013 != nil {
		return nil, parseErr1013
	}
	if expr == nil {
		return nil, nil
	}

	sb := &nodes.SortBy{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		sb.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		sb.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			sb.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			sb.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	sb.Loc.End = p.prev.End
	return sb, nil
}

// parseForUpdateClause parses FOR UPDATE [OF ...] [NOWAIT | WAIT n | SKIP LOCKED].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html#GUID-CFA006CA-6FF1-4972-821E-6996142A51C6
func (p *Parser) parseForUpdateClause() (*nodes.ForUpdateClause, error) {
	start := p.pos()
	p.advance() // consume FOR

	if p.cur.Type == kwUPDATE {
		p.advance()
	}

	fu := &nodes.ForUpdateClause{
		Loc: nodes.Loc{Start: start},
	}

	// OF table_list
	if p.cur.Type == kwOF {
		p.advance()
		var parseErr1014 error
		fu.Tables, parseErr1014 = p.parseExprList()
		if parseErr1014 !=

			// NOWAIT | WAIT n | SKIP LOCKED
			nil {
			return nil, parseErr1014
		}
	}

	switch p.cur.Type {
	case kwNOWAIT:
		fu.NoWait = true
		p.advance()
	case kwWAIT:
		p.advance()
		var parseErr1015 error
		fu.Wait, parseErr1015 = p.parseExpr()
		if parseErr1015 != nil {
			return nil, parseErr1015
		}
	case kwSKIP:
		p.advance()
		if p.cur.Type == kwLOCKED {
			fu.SkipLocked = true
			p.advance()
		}
	}

	fu.Loc.End = p.prev.End
	return fu, nil
}

// parseFetchFirstClause parses OFFSET/FETCH FIRST clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	[ OFFSET n { ROW | ROWS } ]
//	FETCH { FIRST | NEXT } [ n [ PERCENT ] ] { ROW | ROWS } { ONLY | WITH TIES }
func (p *Parser) parseFetchFirstClause() (*nodes.FetchFirstClause, error) {
	start := p.pos()
	fc := &nodes.FetchFirstClause{
		Loc: nodes.Loc{Start: start},
	}

	// OFFSET n ROWS
	if p.cur.Type == kwOFFSET {
		p.advance()
		var parseErr1016 error
		fc.Offset, parseErr1016 = p.parseExpr()
		if parseErr1016 !=
			// ROW | ROWS
			nil {
			return nil, parseErr1016
		}
		if fc.Offset == nil {
			return nil, p.syntaxErrorAtCur()
		}

		if p.cur.Type == kwROW || p.cur.Type == kwROWS {
			p.advance()
		}
	}

	// FETCH FIRST|NEXT
	if p.cur.Type == kwFETCH {
		p.advance()
		// FIRST | NEXT
		if p.cur.Type != kwFIRST && p.cur.Type != kwNEXT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

		// count expression
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			var parseErr1017 error
			fc.Count, parseErr1017 = p.parseExpr()
			if parseErr1017 !=

				// PERCENT
				nil {
				return nil, parseErr1017
			}
			if fc.Count == nil {
				return nil, p.syntaxErrorAtCur()
			}
		}

		if p.cur.Type == kwPERCENT {
			fc.Percent = true
			p.advance()
		}

		// ROW | ROWS
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

		// ONLY | WITH TIES
		if p.cur.Type == kwONLY {
			p.advance()
		} else if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type != kwTIES {
				return nil, p.syntaxErrorAtCur()
			}
			fc.WithTies = true
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseWithClause parses a WITH clause (common table expressions).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	WITH [ RECURSIVE ] cte_name [ ( col1, col2, ... ) ] AS ( subquery ) [, ...]
func (p *Parser) parseWithClause() (*nodes.WithClause, error) {
	start := p.pos()
	p.advance() // consume WITH

	wc := &nodes.WithClause{
		CTEs: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	if p.cur.Type == kwRECURSIVE {
		wc.Recursive = true
		p.advance()
	}

	for {
		cte, parseErr1018 := p.parseCTE()
		if parseErr1018 != nil {
			return nil, parseErr1018
		}
		if cte == nil {
			break
		}
		wc.CTEs.Items = append(wc.CTEs.Items, cte)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	wc.Loc.End = p.prev.End
	return wc, nil
}

// parseCTE parses a single common table expression.
func (p *Parser) parseCTE() (*nodes.CTE, error) {
	if !p.isIdentLike() {
		return nil, nil
	}

	start := p.pos()
	name, parseErr1019 := p.parseIdentifier()
	if parseErr1019 != nil {
		return nil, parseErr1019
	}

	cte := &nodes.CTE{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		cte.Columns = &nodes.List{}
		for {
			col, parseErr1020 := p.parseIdentifier()
			if parseErr1020 != nil {
				return nil, parseErr1020
			}
			if col != "" {
				cte.Columns.Items = append(cte.Columns.Items, &nodes.String{Str: col})
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// AS
	if p.cur.Type == kwAS {
		p.advance()
	}

	// ( subquery )
	if p.cur.Type == '(' {
		p.advance()
		var parseErr1021 error
		cte.Query, parseErr1021 = p.parseSelectStmt()
		if parseErr1021 != nil {
			return nil, parseErr1021
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// SEARCH { BREADTH FIRST | DEPTH FIRST } BY col [, ...] SET ordering_column
	if p.isIdentLikeStr("SEARCH") {
		var parseErr1022 error
		cte.Search, parseErr1022 = p.parseCTESearchClause()
		if parseErr1022 !=

			// CYCLE col [, ...] SET cycle_mark TO value DEFAULT no_value
			nil {
			return nil, parseErr1022
		}
	}

	if p.cur.Type == kwCYCLE {
		var parseErr1023 error
		cte.Cycle, parseErr1023 = p.parseCTECycleClause()
		if parseErr1023 != nil {
			return nil, parseErr1023
		}
	}

	cte.Loc.End = p.prev.End
	return cte, nil
}

// parseCTESearchClause parses a SEARCH clause on a CTE.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	search_clause:
//	    SEARCH { BREADTH FIRST | DEPTH FIRST }
//	    BY c_alias [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ]
//	       [, c_alias [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ] ]...
//	    SET ordering_column
func (p *Parser) parseCTESearchClause() (*nodes.CTESearchClause, error) {
	start := p.pos()
	p.advance() // consume SEARCH

	sc := &nodes.CTESearchClause{Loc: nodes.Loc{Start: start}}

	// BREADTH FIRST | DEPTH FIRST
	if p.isIdentLikeStr("BREADTH") {
		sc.BreadthFirst = true
		p.advance()
	} else if p.isIdentLikeStr("DEPTH") {
		p.advance()
	}
	if p.cur.Type == kwFIRST {
		p.advance()
	}

	// BY col [ASC|DESC] [NULLS FIRST|LAST] [, ...]
	if p.cur.Type == kwBY {
		p.advance()
	}
	var parseErr1024 error
	sc.Columns, parseErr1024 = p.parseOrderByList()
	if parseErr1024 !=

		// SET ordering_column
		nil {
		return nil, parseErr1024
	}

	if p.cur.Type == kwSET {
		p.advance()
		var parseErr1025 error
		sc.SetColumn, parseErr1025 = p.parseIdentifier()
		if parseErr1025 != nil {
			return nil, parseErr1025
		}
	}

	sc.Loc.End = p.prev.End
	return sc, nil
}

// parseCTECycleClause parses a CYCLE clause on a CTE.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	cycle_clause:
//	    CYCLE c_alias [, c_alias ]...
//	    SET cycle_mark_c_alias TO cycle_value
//	    DEFAULT no_cycle_value
func (p *Parser) parseCTECycleClause() (*nodes.CTECycleClause, error) {
	start := p.pos()
	p.advance() // consume CYCLE

	cc := &nodes.CTECycleClause{Loc: nodes.Loc{Start: start}}

	// Column list
	cc.Columns = &nodes.List{}
	for {
		col, parseErr1026 := p.parseIdentifier()
		if parseErr1026 != nil {
			return nil, parseErr1026
		}
		if col == "" {
			break
		}
		cc.Columns.Items = append(cc.Columns.Items, &nodes.String{Str: col})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// SET cycle_mark_alias
	if p.cur.Type == kwSET {
		p.advance()
		var parseErr1027 error
		cc.SetColumn, parseErr1027 = p.parseIdentifier()
		if parseErr1027 !=

			// TO cycle_value
			nil {
			return nil, parseErr1027
		}
	}

	if p.cur.Type == kwTO {
		p.advance()
		var parseErr1028 error
		cc.CycleValue, parseErr1028 = p.parseExpr()
		if parseErr1028 !=

			// DEFAULT no_cycle_value
			nil {
			return nil, parseErr1028
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance()
		var parseErr1029 error
		cc.NoCycleValue, parseErr1029 = p.parseExpr()
		if parseErr1029 != nil {
			return nil, parseErr1029
		}
	}

	cc.Loc.End = p.prev.End
	return cc, nil
}

// parsePivotClause parses a PIVOT clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	pivot_clause ::=
//	    PIVOT (
//	        aggregate_function ( expr ) [ [ AS ] c_alias ]
//	        [, aggregate_function ( expr ) [ [ AS ] c_alias ] ] ...
//	        FOR { column | ( column [, column ] ... ) }
//	        IN ( { { expr | ( expr [, expr] ... ) } [ [ AS ] c_alias ] } [, ...] )
//	    )
func (p *Parser) parsePivotClause() (*nodes.PivotClause, error) {
	start := p.pos()
	p.advance() // consume PIVOT

	pc := &nodes.PivotClause{
		Loc: nodes.Loc{Start: start},
	}

	// PIVOT [ XML ]
	if p.isIdentLikeStr("XML") {
		pc.XML = true
		p.advance()
	}

	if p.cur.Type != '(' {
		pc.Loc.End = p.prev.End
		return pc, nil
	}
	p.advance() // consume '('

	// Parse aggregate function list (before FOR keyword)
	pc.AggFuncs = &nodes.List{}
	for {
		agg, parseErr1030 := p.parseResTarget()
		if parseErr1030 != nil {
			return nil, parseErr1030
		}
		if agg != nil {
			pc.AggFuncs.Items = append(pc.AggFuncs.Items, agg)
		}
		if p.cur.Type != ',' {
			break
		}
		// Peek ahead to see if this comma separates aggregates (before FOR)
		// or if we've reached the FOR keyword
		if p.cur.Type == kwFOR {
			break
		}
		p.advance() // consume ','
		// If the next token is FOR, we went too far — this shouldn't happen
		// because FOR is not an expression start, but be safe
		if p.cur.Type == kwFOR {
			break
		}
	}

	// FOR column | ( column, column, ... )
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if p.cur.Type == '(' {
			// Multi-column: ( col1, col2, ... )
			p.advance() // consume '('
			colList := &nodes.List{}
			for {
				col, parseErr1031 := p.parseColumnRef()
				if parseErr1031 != nil {
					return nil, parseErr1031
				}
				if col != nil {
					colList.Items = append(colList.Items, col)
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
			if len(colList.Items) == 1 {
				if e, ok := colList.Items[0].(nodes.ExprNode); ok {
					pc.ForCol = e
				}
			} else {
				pc.ForCols = colList
			}
		} else {
			var parseErr1032 error
			// Single column reference (not a full expression, to avoid consuming IN)
			pc.ForCol, parseErr1032 = p.parseColumnRef()
			if parseErr1032 !=

				// IN ( ... )
				nil {
				return nil, parseErr1032
			}
		}
	}

	if p.cur.Type == kwIN {
		p.advance() // consume IN
		if p.cur.Type == '(' {
			if pc.XML && p.peekNext().Type == kwSELECT {
				// PIVOT XML allows a subquery in the IN clause
				p.advance() // consume '('
				sub, parseErr1033 := p.parseSelectStmt()
				if parseErr1033 != nil {
					return nil, parseErr1033
				}
				if p.cur.Type == ')' {
					p.advance()
				}
				pc.InList = &nodes.List{Items: []nodes.Node{sub}}
			} else if pc.XML && p.peekNext().Type == kwANY {
				// PIVOT XML allows ANY in IN clause
				p.advance() // consume '('
				anyStart := p.pos()
				p.advance() // consume ANY
				pc.InList = &nodes.List{Items: []nodes.Node{&nodes.ColumnRef{
					Column: "ANY",
					Loc:    nodes.Loc{Start: anyStart, End: p.prev.End},
				}}}
				if p.cur.Type == ')' {
					p.advance()
				}
			} else {
				p.advance()
				var // consume '('
				parseErr1034 error
				pc.InList, parseErr1034 = p.parsePivotInList()
				if parseErr1034 != nil {
					return nil, parseErr1034
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance() // consume outer ')'
	}

	pc.Loc.End = p.prev.End
	return pc, nil
}

// parsePivotInList parses the IN list of a PIVOT clause.
// Each item is: expr [ [ AS ] c_alias ] | ( expr, expr, ... ) [ [ AS ] c_alias ]
func (p *Parser) parsePivotInList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		start := p.pos()
		var expr nodes.ExprNode

		if p.cur.Type == '(' {
			// Tuple value: ( expr, expr, ... )
			p.advance() // consume '('
			tupleList, parseErr1035 := p.parseExprList()
			if parseErr1035 != nil {
				return nil, parseErr1035
			}
			if p.cur.Type == ')' {
				p.advance()
			}
			// Wrap tuple as a ParenExpr containing the first item
			// For multi-value, store as a List via a special representation
			if tupleList.Len() == 1 {
				if e, ok := tupleList.Items[0].(nodes.ExprNode); ok {
					expr = &nodes.ParenExpr{Expr: e, Loc: nodes.Loc{Start: start, End: p.prev.End}}
				}
			} else {
				// Use a FuncCallExpr with empty name to represent a row/tuple
				args := &nodes.List{Items: tupleList.Items}
				expr = &nodes.FuncCallExpr{
					FuncName: &nodes.ObjectName{Name: "", Loc: nodes.Loc{Start: start, End: p.prev.End}},
					Args:     args,
					Loc:      nodes.Loc{Start: start, End: p.prev.End},
				}
			}
		} else {
			var parseErr1036 error
			expr, parseErr1036 = p.parseExpr()
			if parseErr1036 != nil {
				return nil, parseErr1036
			}
		}

		if expr == nil {
			break
		}

		rt := &nodes.ResTarget{
			Expr: expr,
			Loc:  nodes.Loc{Start: start},
		}

		// Optional alias: [ AS ] c_alias
		if p.cur.Type == kwAS {
			p.advance()
			var parseErr1037 error
			rt.Name, parseErr1037 = p.parseIdentifier()
			if parseErr1037 != nil {
				return nil, parseErr1037
			}
		} else if p.isAliasCandidate() {
			var parseErr1038 error
			rt.Name, parseErr1038 = p.parseIdentifier()
			if parseErr1038 != nil {
				return nil, parseErr1038
			}
		}

		rt.Loc.End = p.prev.End
		list.Items = append(list.Items, rt)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseUnpivotClause parses an UNPIVOT clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	unpivot_clause ::=
//	    UNPIVOT [ { INCLUDE | EXCLUDE } NULLS ]
//	    (
//	        column
//	        FOR column
//	        IN ( column [ [ AS ] literal ] [, column [ [ AS ] literal ] ] ... )
//	    )
func (p *Parser) parseUnpivotClause() (*nodes.UnpivotClause, error) {
	start := p.pos()
	p.advance() // consume UNPIVOT

	uc := &nodes.UnpivotClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ INCLUDE | EXCLUDE ] NULLS
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type == kwNULLS {
			p.advance()
		}
		uc.IncludeNulls = true
	} else if p.isIdentLikeStr("EXCLUDE") {
		p.advance()
		if p.cur.Type == kwNULLS {
			p.advance()
		}
		// IncludeNulls stays false (EXCLUDE is the default)
	}

	if p.cur.Type != '(' {
		uc.Loc.End = p.prev.End
		return uc, nil
	}
	p.advance()
	var // consume '('
	parseErr1039 error

	// Value column(s)
	uc.ValueCol, parseErr1039 = p.parseColumnRef()
	if parseErr1039 !=

		// FOR pivot_column(s)
		nil {
		return nil, parseErr1039
	}

	if p.cur.Type == kwFOR {
		p.advance()
		var parseErr1040 error
		uc.PivotCol, parseErr1040 = p.parseColumnRef()
		if parseErr1040 !=

			// IN ( column [ AS literal ], ... )
			nil {
			return nil, parseErr1040
		}
	}

	if p.cur.Type == kwIN {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			uc.InList = &nodes.List{}
			for {
				start := p.pos()
				col, parseErr1041 := p.parseColumnRef()
				if parseErr1041 != nil {
					return nil, parseErr1041
				}
				if col == nil || col.Column == "" {
					break
				}
				rt := &nodes.ResTarget{
					Expr: col,
					Loc:  nodes.Loc{Start: start},
				}
				// Optional AS alias (can be identifier or string literal)
				if p.cur.Type == kwAS {
					p.advance()
					if p.cur.Type == tokSCONST {
						rt.Name = p.cur.Str
						p.advance()
					} else {
						var parseErr1042 error
						rt.Name, parseErr1042 = p.parseIdentifier()
						if parseErr1042 != nil {
							return nil, parseErr1042
						}
					}
				}
				rt.Loc.End = p.prev.End
				uc.InList.Items = append(uc.InList.Items, rt)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance() // consume outer ')'
	}

	uc.Loc.End = p.prev.End
	return uc, nil
}

// parseModelClause parses an Oracle MODEL clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	model_clause ::=
//	    MODEL
//	      [ cell_reference_options ]
//	      [ return_rows_clause ]
//	      [ reference_model ]...
//	      main_model
//
//	cell_reference_options ::=
//	    { IGNORE NAV | KEEP NAV }
//	    { UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE }
//
//	return_rows_clause ::=
//	    RETURN { UPDATED | ALL } ROWS
//
//	reference_model ::=
//	    REFERENCE reference_model_name ON ( subquery )
//	        model_column_clauses
//	        [ cell_reference_options ]
//
//	main_model ::=
//	    [ MAIN main_model_name ]
//	        model_column_clauses
//	        [ cell_reference_options ]
//	        model_rules_clause
//
//	model_column_clauses ::=
//	    [ PARTITION BY ( expr [ [ AS ] c_alias ] [, ...] ) ]
//	    DIMENSION BY ( expr [ [ AS ] c_alias ] [, ...] )
//	    MEASURES ( expr [ [ AS ] c_alias ] [, ...] )
//
//	model_rules_clause ::=
//	    [ RULES ]
//	    [ { UPDATE | UPSERT [ ALL ] } ]
//	    [ { AUTOMATIC | SEQUENTIAL } ORDER ]
//	    [ model_iterate_clause ]
//	    ( cell_assignment [, ...] )
//
//	model_iterate_clause ::=
//	    ITERATE ( number ) [ UNTIL ( condition ) ]
//
//	cell_assignment ::=
//	    measure_column [ dimension_subscripts ] = expr
//
//	single_column_for_loop ::=
//	    FOR dimension_column
//	      { IN ( { literal [, ...] | subquery } )
//	      | [ LIKE pattern ] FROM literal TO literal { INCREMENT | DECREMENT } literal
//	      }
//
//	multi_column_for_loop ::=
//	    FOR ( dimension_column [, ...] ) IN
//	      ( ( literal [, ...] ) [, ( literal [, ...] ) ]...
//	      | subquery
//	      )
func (p *Parser) parseModelClause() (*nodes.ModelClause, error) {
	start := p.pos()
	p.advance() // consume MODEL

	mc := &nodes.ModelClause{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr1043 error

	// cell_reference_options (optional, before RETURN/REFERENCE/MAIN/PARTITION/DIMENSION)
	mc.CellRefOptions, parseErr1043 = p.parseModelCellRefOptions()
	if parseErr1043 !=

		// RETURN { UPDATED | ALL } ROWS
		nil {
		return nil, parseErr1043
	}

	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if p.cur.Type == kwUPDATED {
			mc.ReturnRows = "UPDATED"
			p.advance()
		} else if p.cur.Type == kwALL {
			mc.ReturnRows = "ALL"
			p.advance()
		}
		if p.cur.Type == kwROWS {
			p.advance()
		}
	}

	// REFERENCE models (zero or more)
	for p.cur.Type == kwREFERENCE {
		ref, parseErr1044 := p.parseModelRefModel()
		if parseErr1044 != nil {
			return nil, parseErr1044
		}
		mc.RefModels = append(mc.RefModels, ref)
	}
	var parseErr1045 error

	// main_model
	mc.MainModel, parseErr1045 = p.parseModelMainModel()
	if parseErr1045 != nil {
		return nil, parseErr1045

		// parseModelCellRefOptions parses optional cell_reference_options.
	}

	mc.Loc.End = p.prev.End
	return mc, nil
}

func (p *Parser) parseModelCellRefOptions() (*nodes.ModelCellRefOptions, error) {
	opts := &nodes.ModelCellRefOptions{Loc: nodes.Loc{Start: p.pos()}}
	found := false

	// IGNORE NAV | KEEP NAV
	if p.isIdentLikeStr("IGNORE") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.IgnoreNav = true
		found = true
	} else if p.isIdentLikeStr("KEEP") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.KeepNav = true
		found = true
	}

	// UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE
	if p.cur.Type == kwUNIQUE {
		p.advance()
		if p.cur.Type == kwDIMENSION {
			opts.UniqueDimension = true
			p.advance()
			found = true
		} else if p.isIdentLikeStr("SINGLE") {
			p.advance()
			if p.cur.Type == kwREFERENCE {
				p.advance()
			}
			opts.UniqueSingleRef = true
			found = true
		}
	}

	if !found {
		return nil, nil
	}

	opts.Loc.End = p.prev.End
	return opts, nil
}

// parseModelRefModel parses a REFERENCE model.
func (p *Parser) parseModelRefModel() (*nodes.ModelRefModel, error) {
	start := p.pos()
	p.advance() // consume REFERENCE

	ref := &nodes.ModelRefModel{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr1046 error

	// reference_model_name
	ref.Name, parseErr1046 = p.parseIdentifier()
	if parseErr1046 !=

		// ON ( subquery )
		nil {
		return nil, parseErr1046
	}

	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1047 error
			ref.Subquery, parseErr1047 = p.parseSelectStmt()
			if parseErr1047 != nil {
				return nil, parseErr1047
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}
	var parseErr1048 error

	// model_column_clauses
	ref.ColumnClauses, parseErr1048 = p.parseModelColumnClauses()
	if parseErr1048 !=

		// optional cell_reference_options
		nil {
		return nil, parseErr1048
	}
	var parseErr1049 error

	ref.CellRefOptions, parseErr1049 = p.parseModelCellRefOptions()
	if parseErr1049 != nil {
		return nil, parseErr1049
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseModelMainModel parses the main_model.
func (p *Parser) parseModelMainModel() (*nodes.ModelMainModel, error) {
	start := p.pos()
	mm := &nodes.ModelMainModel{
		Loc: nodes.Loc{Start: start},
	}

	// [ MAIN main_model_name ]
	if p.cur.Type == kwMAIN {
		p.advance()
		var parseErr1050 error
		mm.Name, parseErr1050 = p.parseIdentifier()
		if parseErr1050 !=

			// model_column_clauses
			nil {
			return nil, parseErr1050
		}
	}
	var parseErr1051 error

	mm.ColumnClauses, parseErr1051 = p.parseModelColumnClauses()
	if parseErr1051 !=

		// optional cell_reference_options
		nil {
		return nil, parseErr1051
	}
	var parseErr1052 error

	mm.CellRefOptions, parseErr1052 = p.parseModelCellRefOptions()
	if parseErr1052 !=

		// model_rules_clause
		nil {
		return nil, parseErr1052
	}
	var parseErr1053 error

	mm.RulesClause, parseErr1053 = p.parseModelRulesClause()
	if parseErr1053 != nil {
		return nil, parseErr1053

		// parseModelColumnClauses parses PARTITION BY / DIMENSION BY / MEASURES.
	}

	mm.Loc.End = p.prev.End
	return mm, nil
}

func (p *Parser) parseModelColumnClauses() (*nodes.ModelColumnClauses, error) {
	start := p.pos()
	cc := &nodes.ModelColumnClauses{
		Loc: nodes.Loc{Start: start},
	}

	// [ PARTITION BY ( expr [AS alias], ... ) ]
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1054 error
			cc.PartitionBy, parseErr1054 = p.parseModelColumnList()
			if parseErr1054 != nil {
				return nil, parseErr1054
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// DIMENSION BY ( expr [AS alias], ... )
	if p.cur.Type == kwDIMENSION {
		p.advance() // DIMENSION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1055 error
			cc.DimensionBy, parseErr1055 = p.parseModelColumnList()
			if parseErr1055 != nil {
				return nil, parseErr1055
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// MEASURES ( expr [AS alias], ... )
	if p.cur.Type == kwMEASURES {
		p.advance() // MEASURES
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1056 error
			cc.Measures, parseErr1056 = p.parseModelColumnList()
			if parseErr1056 != nil {
				return nil, parseErr1056
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	cc.Loc.End = p.prev.End
	return cc, nil
}

// parseModelColumnList parses a comma-separated list of expr [ [AS] alias ] for MODEL clauses.
func (p *Parser) parseModelColumnList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		rt, parseErr1057 := p.parseResTarget()
		if parseErr1057 != nil {
			return nil, parseErr1057
		}
		if rt == nil {
			break
		}
		list.Items = append(list.Items, rt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseModelRulesClause parses the model_rules_clause.
func (p *Parser) parseModelRulesClause() (*nodes.ModelRulesClause, error) {
	start := p.pos()
	rc := &nodes.ModelRulesClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ RULES ]
	if p.cur.Type == kwRULES {
		p.advance()
	}

	// [ UPDATE | UPSERT [ALL] ]
	if p.cur.Type == kwUPDATE {
		rc.UpdateMode = "UPDATE"
		p.advance()
	} else if p.cur.Type == kwUPSERT {
		p.advance()
		if p.cur.Type == kwALL {
			rc.UpdateMode = "UPSERT ALL"
			p.advance()
		} else {
			rc.UpdateMode = "UPSERT"
		}
	}

	// [ AUTOMATIC | SEQUENTIAL ] ORDER
	if p.cur.Type == kwAUTOMATIC {
		rc.OrderMode = "AUTOMATIC"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	} else if p.cur.Type == kwSEQUENTIAL {
		rc.OrderMode = "SEQUENTIAL"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	}

	// [ ITERATE ( number ) [ UNTIL ( condition ) ] ]
	if p.cur.Type == kwITERATE {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1058 error
			rc.Iterate, parseErr1058 = p.parseExpr()
			if parseErr1058 != nil {
				return nil, parseErr1058
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// [ UNTIL ( condition ) ]
		if p.cur.Type == kwUNTIL {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				var parseErr1059 error
				rc.Until, parseErr1059 = p.parseExpr()
				if parseErr1059 != nil {
					return nil, parseErr1059
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	// ( cell_assignment [, ...] )
	if p.cur.Type == '(' {
		p.advance()
		rc.Rules = &nodes.List{}
		for {
			rule, parseErr1060 := p.parseModelRule()
			if parseErr1060 != nil {
				return nil, parseErr1060
			}
			if rule == nil {
				break
			}
			rc.Rules.Items = append(rc.Rules.Items, rule)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	rc.Loc.End = p.prev.End
	return rc, nil
}

// parseModelRule parses a single cell_assignment: measure_column[dim1, dim2] = expr.
func (p *Parser) parseModelRule() (*nodes.ModelRule, error) {
	start := p.pos()

	// Parse the left side: measure_column [ dim_subscripts ]
	// The left side is an identifier (possibly qualified) followed by [ subscripts ]
	cellRef, parseErr1061 := p.parseModelCellRef()
	if parseErr1061 != nil {
		return nil, parseErr1061
	}
	if cellRef == nil {
		return nil, nil
	}

	rule := &nodes.ModelRule{
		CellRef: cellRef,
		Loc:     nodes.Loc{Start: start},
	}

	// = expr
	if p.cur.Type == '=' {
		p.advance()
		var parseErr1062 error
		rule.Expr, parseErr1062 = p.parseExpr()
		if parseErr1062 != nil {
			return nil, parseErr1062
		}
	}

	rule.Loc.End = p.prev.End
	return rule, nil
}

// parseModelCellRef parses a model cell reference: measure[dim1, dim2, ...].
// Dimension subscripts can be expressions, ANY, or FOR loops (single/multi column).
func (p *Parser) parseModelCellRef() (nodes.ExprNode, error) {
	start := p.pos()

	// Parse the measure column name as a column reference
	if !p.isIdentLike() {
		return nil, nil
	}
	name, parseErr1063 := p.parseIdentifier()
	if parseErr1063 != nil {
		return nil, parseErr1063
	}
	col := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}

	// [ dim_subscripts ]
	if p.cur.Type != '[' {
		return col, nil
	}
	p.advance() // consume '['

	// Parse subscript expressions (comma-separated)
	args := &nodes.List{}
	for {
		if p.cur.Type == ']' {
			break
		}
		expr, parseErr1064 := p.parseExpr()
		if parseErr1064 != nil {
			return nil, parseErr1064
		}
		if expr == nil {
			break
		}
		args.Items = append(args.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ']' {
		p.advance()
	}

	// Represent as a FuncCallExpr with the column name and subscript args
	return &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     args,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseSampleClause parses a SAMPLE clause on a table reference.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	sample_clause ::=
//	    SAMPLE [ BLOCK ] ( sample_percent )
//	    [ SEED ( seed_value ) ]
func (p *Parser) parseSampleClause() (*nodes.SampleClause, error) {
	start := p.pos()
	p.advance() // consume SAMPLE

	sc := &nodes.SampleClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ BLOCK ]
	if p.cur.Type == kwBLOCK {
		sc.Block = true
		p.advance()
	}

	// ( percent )
	if p.cur.Type == '(' {
		p.advance()
		var parseErr1065 error
		sc.Percent, parseErr1065 = p.parseExpr()
		if parseErr1065 != nil {
			return nil, parseErr1065
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// [ SEED ( value ) ]
	if p.cur.Type == kwSEED {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1066 error
			sc.Seed, parseErr1066 = p.parseExpr()
			if parseErr1066 != nil {
				return nil, parseErr1066
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	sc.Loc.End = p.prev.End
	return sc, nil
}

// parseFlashbackClause parses a flashback query clause on a table reference.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	flashback_query_clause ::=
//	    { VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
//	    | AS OF { SCN | TIMESTAMP } expr
//	    }
func (p *Parser) parseFlashbackClause() (*nodes.FlashbackClause, error) {
	start := p.pos()
	fc := &nodes.FlashbackClause{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwVERSIONS {
		// VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
		// VERSIONS PERIOD FOR valid_time_column BETWEEN expr AND expr
		fc.IsVersions = true
		p.advance() // consume VERSIONS

		// PERIOD FOR valid_time_column
		if p.isIdentLikeStr("PERIOD") {
			fc.IsPeriodFor = true
			p.advance() // consume PERIOD
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
			}
			var parseErr1067 error
			fc.PeriodColumn, parseErr1067 = p.parseIdentifier()
			if parseErr1067 != nil {
				return nil, parseErr1067
			}
		}

		if p.cur.Type == kwBETWEEN {
			p.advance() // consume BETWEEN
		}

		// SCN | TIMESTAMP type marker
		if p.cur.Type == kwSCN {
			fc.Type = "SCN"
			p.advance()
		} else if p.cur.Type == kwTIMESTAMP {
			fc.Type = "TIMESTAMP"
			p.advance()
		}
		var parseErr1068 error

		// low expr — parse above AND precedence to not consume the BETWEEN...AND
		fc.VersionsLow, parseErr1068 = p.parseExprPrec(precNot)
		if parseErr1068 !=

			// AND
			nil {
			return nil, parseErr1068
		}

		if p.cur.Type == kwAND {
			p.advance()
		}

		// Skip repeated SCN/TIMESTAMP keyword before high expr
		if p.cur.Type == kwSCN || p.cur.Type == kwTIMESTAMP {
			p.advance()
		}
		var parseErr1069 error

		// high expr
		fc.VersionsHigh, parseErr1069 = p.parseExpr()
		if parseErr1069 != nil {
			return nil, parseErr1069

			// AS OF { SCN | TIMESTAMP } expr
			// AS OF PERIOD FOR valid_time_column expr
		}

	} else if p.cur.Type == kwAS {

		p.advance() // consume AS
		if p.cur.Type == kwOF {
			p.advance() // consume OF
		}

		// PERIOD FOR valid_time_column
		if p.isIdentLikeStr("PERIOD") {
			fc.IsPeriodFor = true
			p.advance() // consume PERIOD
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
			}
			var parseErr1070 error
			fc.PeriodColumn, parseErr1070 = p.parseIdentifier()
			if parseErr1070 !=
				// SCN | TIMESTAMP type marker
				nil {
				return nil, parseErr1070
			}

			if p.cur.Type == kwSCN {
				fc.Type = "SCN"
				p.advance()
			} else if p.cur.Type == kwTIMESTAMP {
				fc.Type = "TIMESTAMP"
				p.advance()
			}
			var parseErr1071 error
			fc.Expr, parseErr1071 = p.parseExpr()
			if parseErr1071 !=

				// SCN | TIMESTAMP
				nil {
				return nil, parseErr1071
			}
		} else {

			if p.cur.Type == kwSCN {
				fc.Type = "SCN"
				p.advance()
			} else if p.cur.Type == kwTIMESTAMP {
				fc.Type = "TIMESTAMP"
				p.advance()
			}
			var parseErr1072 error
			// expr
			fc.Expr, parseErr1072 = p.parseExpr()
			if parseErr1072 != nil {
				return nil, parseErr1072
			}
		}
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseWindowClause parses a WINDOW clause.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	window_clause:
//	    WINDOW window_name AS ( window_specification )
//	    [, window_name AS ( window_specification ) ]...
func (p *Parser) parseWindowClause() ([]*nodes.WindowDef, error) {
	p.advance() // consume WINDOW (soft keyword)

	var defs []*nodes.WindowDef
	for {
		if !p.isIdentLike() {
			break
		}
		start := p.pos()
		name, parseErr1073 := p.parseIdentifier()
		if parseErr1073 != nil {
			return nil, parseErr1073
		}
		wd := &nodes.WindowDef{
			Name: name,
			Loc:  nodes.Loc{Start: start},
		}

		// AS ( window_specification )
		if p.cur.Type == kwAS {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.advance()
			ws := &nodes.WindowSpec{Loc: nodes.Loc{Start: p.pos()}}

			// [ existing_window_name ] — only if followed by something other than BY
			if p.isIdentLike() && p.peekNext().Type != kwBY &&
				p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
				p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
				var parseErr1074 error
				ws.WindowName, parseErr1074 = p.parseIdentifier()
				if parseErr1074 !=

					// PARTITION BY
					nil {
					return nil, parseErr1074
				}
			}

			if p.cur.Type == kwPARTITION {
				p.advance()
				if p.cur.Type == kwBY {
					p.advance()
				}
				ws.PartitionBy = &nodes.List{}
				for {
					expr, parseErr1075 := p.parseExpr()
					if parseErr1075 != nil {
						return nil, parseErr1075
					}
					if expr != nil {
						ws.PartitionBy.Items = append(ws.PartitionBy.Items, expr)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
			}

			// ORDER BY
			if p.cur.Type == kwORDER {
				p.advance()
				if p.cur.Type == kwBY {
					p.advance()
				}
				var parseErr1076 error
				ws.OrderBy, parseErr1076 = p.parseOrderByList()
				if parseErr1076 !=

					// Windowing clause: ROWS | RANGE | GROUPS
					nil {
					return nil, parseErr1076
				}
			}

			if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
				var parseErr1077 error
				ws.Frame, parseErr1077 = p.parseWindowFrame()
				if parseErr1077 != nil {
					return nil, parseErr1077
				}
			}

			if p.cur.Type == ')' {
				p.advance()
			}
			ws.Loc.End = p.prev.End
			wd.Spec = ws
		}

		wd.Loc.End = p.prev.End
		defs = append(defs, wd)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return defs, nil
}

// parsePartitionExtClause parses a partition extension clause.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	partition_extension_clause:
//	    { PARTITION ( partition_name )
//	    | PARTITION FOR ( partition_key_value [, partition_key_value ]... )
//	    | SUBPARTITION ( subpartition_name )
//	    | SUBPARTITION FOR ( subpartition_key_value [, subpartition_key_value ]... )
//	    }
func (p *Parser) parsePartitionExtClause() (*nodes.PartitionExtClause, error) {
	start := p.pos()
	pe := &nodes.PartitionExtClause{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwSUBPARTITION {
		pe.IsSubpartition = true
	}
	p.advance() // consume PARTITION or SUBPARTITION

	// FOR keyword
	if p.cur.Type == kwFOR {
		pe.IsFor = true
		p.advance()
	}

	if p.cur.Type == '(' {
		p.advance()
		if pe.IsFor {
			var parseErr1078 error
			pe.Keys, parseErr1078 = p.parseExprList()
			if parseErr1078 != nil {
				return nil, parseErr1078
			}
		} else {
			var parseErr1079 error
			pe.Name, parseErr1079 = p.parseIdentifier()
			if parseErr1079 != nil {
				return nil, parseErr1079
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	pe.Loc.End = p.prev.End
	return pe, nil
}

// parseTableCollectionExpr parses TABLE(collection_expression) [(+)].
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	table_collection_expression:
//	    TABLE ( collection_expression ) [ ( + ) ]
func (p *Parser) parseTableCollectionExpr(start int) (nodes.TableExpr, error) {
	p.advance() // consume TABLE
	tc := &nodes.TableCollectionExpr{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr1080 error
		tc.Expr, parseErr1080 = p.parseExpr()
		if parseErr1080 != nil {
			return nil, parseErr1080
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional (+)
	if p.cur.Type == '(' && p.peekNext().Type == '+' {
		p.advance() // (
		p.advance() // +
		if p.cur.Type == ')' {
			p.advance()
		}
		tc.OuterJoin = true
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1081 error
		tc.Alias, parseErr1081 = p.parseAlias()
		if parseErr1081 != nil {
			return nil, parseErr1081
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1082 error
		tc.Alias, parseErr1082 = p.parseAlias()
		if parseErr1082 != nil {
			return nil, parseErr1082
		}
	}

	tc.Loc.End = p.prev.End
	return tc, nil
}

// parseContainersOrShards parses CONTAINERS(table) or SHARDS(table).
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	containers_clause:
//	    CONTAINERS ( [ schema. ] { table | view } )
//	shards_clause:
//	    SHARDS ( [ schema. ] { table | view } )
func (p *Parser) parseContainersOrShards(start int, isShards bool) (nodes.TableExpr, error) {
	p.advance() // consume CONTAINERS or SHARDS
	ce := &nodes.ContainersExpr{
		IsShards: isShards,
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr1083 error
		ce.Name, parseErr1083 = p.parseObjectName()
		if parseErr1083 != nil {
			return nil, parseErr1083
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1084 error
		ce.Alias, parseErr1084 = p.parseAlias()
		if parseErr1084 != nil {
			return nil, parseErr1084
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1085 error
		ce.Alias, parseErr1085 = p.parseAlias()
		if parseErr1085 != nil {
			return nil, parseErr1085
		}
	}

	ce.Loc.End = p.prev.End
	return ce, nil
}

// parseMatchRecognize parses a MATCH_RECOGNIZE clause.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	row_pattern_clause:
//	    MATCH_RECOGNIZE (
//	      [ PARTITION BY column [, column ]... ]
//	      [ ORDER BY column [ ASC | DESC ] [, column [ ASC | DESC ] ]... ]
//	      [ MEASURES row_pattern_measure_column [, row_pattern_measure_column ]... ]
//	      [ ONE ROW PER MATCH | ALL ROWS PER MATCH [ { SHOW | OMIT } EMPTY MATCHES ] ]
//	      [ AFTER MATCH SKIP
//	        { PAST LAST ROW | TO NEXT ROW | TO FIRST variable_name
//	        | TO LAST variable_name | TO variable_name } ]
//	      PATTERN ( row_pattern )
//	      [ SUBSET subset_item [, subset_item ]... ]
//	      DEFINE row_pattern_definition [, row_pattern_definition ]...
//	    )
func (p *Parser) parseMatchRecognize(start int) (nodes.TableExpr, error) {
	p.advance() // consume MATCH_RECOGNIZE
	mr := &nodes.MatchRecognizeClause{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		mr.Loc.End = p.prev.End
		return mr, nil
	}
	p.advance() // consume '('

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		var parseErr1086 error
		mr.PartitionBy, parseErr1086 = p.parseExprList()
		if parseErr1086 !=

			// ORDER BY
			nil {
			return nil, parseErr1086
		}
	}

	if p.cur.Type == kwORDER {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		var parseErr1087 error
		mr.OrderBy, parseErr1087 = p.parseOrderByList()
		if parseErr1087 !=

			// MEASURES
			nil {
			return nil, parseErr1087
		}
	}

	if p.cur.Type == kwMEASURES {
		p.advance()
		mr.Measures = &nodes.List{}
		for {
			rt, parseErr1088 := p.parseResTarget()
			if parseErr1088 != nil {
				return nil, parseErr1088
			}
			if rt == nil {
				break
			}
			mr.Measures.Items = append(mr.Measures.Items, rt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// ONE ROW PER MATCH | ALL ROWS PER MATCH [...]
	if p.isIdentLikeStr("ONE") {
		p.advance() // ONE
		if p.cur.Type == kwROW {
			p.advance()
		}
		if p.isIdentLikeStr("PER") {
			p.advance()
		}
		if p.isIdentLikeStr("MATCH") {
			p.advance()
		}
		mr.RowsPerMatch = "ONE ROW PER MATCH"
	} else if p.cur.Type == kwALL {
		p.advance() // ALL
		if p.cur.Type == kwROWS {
			p.advance()
		}
		if p.isIdentLikeStr("PER") {
			p.advance()
		}
		if p.isIdentLikeStr("MATCH") {
			p.advance()
		}
		mr.RowsPerMatch = "ALL ROWS PER MATCH"
		if p.isIdentLikeStr("SHOW") {
			p.advance()
			if p.isIdentLikeStr("EMPTY") {
				p.advance()
			}
			if p.isIdentLikeStr("MATCHES") {
				p.advance()
			}
			mr.RowsPerMatch = "ALL ROWS PER MATCH SHOW EMPTY MATCHES"
		} else if p.isIdentLikeStr("OMIT") {
			p.advance()
			if p.isIdentLikeStr("EMPTY") {
				p.advance()
			}
			if p.isIdentLikeStr("MATCHES") {
				p.advance()
			}
			mr.RowsPerMatch = "ALL ROWS PER MATCH OMIT EMPTY MATCHES"
		}
	}

	// AFTER MATCH SKIP ...
	if p.isIdentLikeStr("AFTER") {
		p.advance()
		if p.isIdentLikeStr("MATCH") {
			p.advance()
		}
		if p.cur.Type == kwSKIP {
			p.advance()
		}
		if p.isIdentLikeStr("PAST") {
			p.advance()
			if p.cur.Type == kwLAST {
				p.advance()
			}
			if p.cur.Type == kwROW {
				p.advance()
			}
			mr.AfterMatch = "PAST LAST ROW"
		} else if p.cur.Type == kwTO {
			p.advance()
			if p.cur.Type == kwNEXT {
				p.advance()
				if p.cur.Type == kwROW {
					p.advance()
				}
				mr.AfterMatch = "TO NEXT ROW"
			} else if p.cur.Type == kwFIRST {
				p.advance()
				name, parseErr1089 := p.parseIdentifier()
				if parseErr1089 != nil {
					return nil, parseErr1089
				}
				mr.AfterMatch = "TO FIRST " + name
			} else if p.cur.Type == kwLAST {
				p.advance()
				name, parseErr1090 := p.parseIdentifier()
				if parseErr1090 != nil {
					return nil, parseErr1090
				}
				mr.AfterMatch = "TO LAST " + name
			} else {
				name, parseErr1091 := p.parseIdentifier()
				if parseErr1091 != nil {
					return nil, parseErr1091
				}
				mr.AfterMatch = "TO " + name
			}
		}
	}

	// PATTERN ( row_pattern ) — capture raw text between parens
	if p.isIdentLikeStr("PATTERN") {
		p.advance()
		if p.cur.Type == '(' {
			mr.Pattern = p.consumeBalancedParensAsText()
		}
	}

	// SUBSET
	if p.isIdentLikeStr("SUBSET") {
		p.advance()
		mr.Subsets = &nodes.List{}
		for {
			if !p.isIdentLike() {
				break
			}
			subStart := p.pos()
			name, parseErr1092 := p.parseIdentifier()
			if parseErr1092 != nil {
				return nil, parseErr1092
			}
			rt := &nodes.ResTarget{
				Name: name,
				Loc:  nodes.Loc{Start: subStart},
			}
			// = ( var1, var2, ... )
			if p.cur.Type == '=' {
				p.advance()
				if p.cur.Type == '(' {
					p.advance()
					var parseErr1093 error
					rt.Expr, parseErr1093 = p.parseExpr()
					if parseErr1093 != nil {
						return nil, parseErr1093
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				}
			}
			rt.Loc.End = p.prev.End
			mr.Subsets.Items = append(mr.Subsets.Items, rt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// DEFINE
	if p.isIdentLikeStr("DEFINE") {
		p.advance()
		mr.Definitions = &nodes.List{}
		for {
			if !p.isIdentLike() {
				break
			}
			defStart := p.pos()
			name, parseErr1094 := p.parseIdentifier()
			if parseErr1094 != nil {
				return nil, parseErr1094
			}
			rt := &nodes.ResTarget{
				Name: name,
				Loc:  nodes.Loc{Start: defStart},
			}
			if p.cur.Type == kwAS {
				p.advance()
			}
			var parseErr1095 error
			rt.Expr, parseErr1095 = p.parseExpr()
			if parseErr1095 != nil {
				return nil, parseErr1095
			}
			rt.Loc.End = p.prev.End
			mr.Definitions.Items = append(mr.Definitions.Items, rt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1096 error
		mr.Alias, parseErr1096 = p.parseAlias()
		if parseErr1096 != nil {
			return nil, parseErr1096
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1097 error
		mr.Alias, parseErr1097 = p.parseAlias()
		if parseErr1097 != nil {
			return nil, parseErr1097
		}
	}

	mr.Loc.End = p.prev.End
	return mr, nil
}

// consumeBalancedParensAsText consumes tokens between ( and ) and returns them as a concatenated string.
func (p *Parser) consumeBalancedParensAsText() string {
	if p.cur.Type != '(' {
		return ""
	}
	p.advance() // consume '('
	depth := 1
	var parts []string
	for depth > 0 {
		if p.cur.Type == '(' {
			depth++
			parts = append(parts, "(")
		} else if p.cur.Type == ')' {
			depth--
			if depth == 0 {
				p.advance() // consume final ')'
				break
			}
			parts = append(parts, ")")
		} else if p.cur.Type == tokEOF {
			break
		} else {
			parts = append(parts, p.cur.Str)
		}
		p.advance()
	}
	result := ""
	for i, s := range parts {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}
