package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the analytic view family: CREATE and ALTER ANALYTIC VIEW,
// ATTRIBUTE DIMENSION and HIERARCHY, and the older CREATE and ALTER
// DIMENSION.
//
// BNF: oracle/parser/bnf/CREATE-ANALYTIC-VIEW.bnf, ALTER-ANALYTIC-VIEW.bnf,
// CREATE-ATTRIBUTE-DIMENSION.bnf, ALTER-ATTRIBUTE-DIMENSION.bnf,
// CREATE-HIERARCHY.bnf, ALTER-HIERARCHY.bnf, CREATE-DIMENSION.bnf,
// ALTER-DIMENSION.bnf
// (https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ANALYTIC-VIEW.html
// and the pages beside it).
//
// Where the pages and the pinned Oracle 23.26.3 disagree, the parser follows
// Oracle, as its compat.tsv rows and probes on the server show:
//
//   - classification_clause is [ CAPTION 'string' ] [ DESCRIPTION 'string' ]
//     [ CLASSIFICATION name [ VALUE 'string' ] [ LANGUAGE 'string' ] ]...,
//     in that order. It comes after SHARING and, in an attribute
//     dimension, after DIMENSION TYPE; a hierarchy level takes none.
//   - source_clause has no REMOTE: USING REMOTE a reads the table REMOTE
//     with the alias a, and USING REMOTE a.b fails at the dot.
//   - An analytic view's dim_key may carry classifications after its alias;
//     KEY and REFERENCES take one name each or the same number in
//     parentheses, and HIERARCHIES is required.
//   - A base measure is [ FACT expression ] [ AGGREGATE BY aggr_function ];
//     aggr_function also takes COUNT DISTINCT.
//   - cache_clause is CACHE cache_specification..., with no parentheses or
//     commas: MEASURE GROUP { ALL | ( measure [, measure ]... ) } followed by
//     one or more LEVELS ( [ level_specification [, ...] ] ) MATERIALIZED
//     [ USING [ schema. ] table ]. ALTER ANALYTIC VIEW ADD CACHE takes one
//     LEVELS ... MATERIALIZED [ USING table ]; DROP CACHE takes LEVELS
//     without MATERIALIZED.
//   - An attribute dimension's join paths follow its sources: USING a, b
//     JOIN PATH p ON a.x = b.y. Its key takes no NOT NULL or SKIP WHEN
//     NULL, ALTERNATE KEY may repeat, dim_order_clause is ORDER BY [ MIN |
//     MAX ] attribute [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...]
//     with no parentheses, and the ALL clause needs at least one MEMBER.
//   - CREATE DIMENSION takes no OR REPLACE, and its level clauses first, then hierarchy and
//     attribute clauses; a dimension_join_clause stands inside a hierarchy
//     alone. A level's first column is qualified by its table; a column
//     list may also be a single column with no parentheses.
//   - ALTER DIMENSION takes ADD clauses or DROP clauses, any number of
//     them, and COMPILE among them; ADD with DROP is ORA-30348. DROP takes
//     LEVEL level [ RESTRICT | CASCADE ], HIERARCHY hierarchy and ATTRIBUTE
//     attribute [ LEVEL level [ COLUMN column ] ]; an extended attribute is
//     ATTRIBUTE a LEVEL l, as in CREATE, not OF l; SKIP WHEN takes NULL.
//
// The expression of a calculated measure, AS ( analytic_view_expression ),
// has a grammar of its own (hierarchy navigation, SHARE_OF, QUALIFY, ...)
// that no page here holds: its parentheses are read whole, and must not be
// empty.

// isAnalyticNameToken reports whether the current token can name an object,
// a level, an attribute or a measure: an identifier, quoted or not, that is
// not a reserved word.
func (p *Parser) isAnalyticNameToken() bool {
	switch {
	case p.cur.Type == tokQIDENT:
		return true
	case p.cur.Type == tokIDENT || p.cur.Type >= 2000:
		return !isOracleSQLReservedKeyword(p.cur)
	}
	return false
}

// anParseName reads a required name (isAnalyticNameToken).
func (p *Parser) anParseName() (string, error) {
	if !p.isAnalyticNameToken() {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// anParseObjectName reads [ schema. ] name.
func (p *Parser) anParseObjectName() (*nodes.ObjectName, error) {
	start := p.pos()
	name, err := p.anParseName()
	if err != nil {
		return nil, err
	}
	obj := &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start}}
	if p.cur.Type == '.' {
		p.advance()
		second, err := p.anParseName()
		if err != nil {
			return nil, err
		}
		obj.Schema, obj.Name = name, second
	}
	obj.Loc.End = p.prev.End
	return obj, nil
}

// anParseDottedName reads name [ . name ]... of at most max parts, joined by
// dots. A dot past max parts is left for the caller, which fails at it.
func (p *Parser) anParseDottedName(max int) (string, error) {
	name, err := p.anParseName()
	if err != nil {
		return "", err
	}
	for parts := 1; parts < max && p.cur.Type == '.'; parts++ {
		p.advance()
		part, err := p.anParseName()
		if err != nil {
			return "", err
		}
		name += "." + part
	}
	return name, nil
}

// anExpect consumes a token of type t or fails at the current token.
func (p *Parser) anExpect(t int) error {
	if p.cur.Type != t {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// anExpectWord consumes the unquoted word w or fails at the current token.
func (p *Parser) anExpectWord(w string) error {
	if !p.isKeywordStr(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// anNextIsWord reports whether the token after the current one is the
// unquoted word w.
func (p *Parser) anNextIsWord(w string) bool {
	next := p.peekNext()
	return next.Type != tokQIDENT && (next.Type == tokIDENT || next.Type >= 2000) && next.Str == w
}

// anParseString reads a required string literal: 'text', N'text' or
// q'[text]' (ORA-01780 otherwise).
func (p *Parser) anParseString() (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNCHARLIT {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// anText is the source text from start to end, as the AST's summary fields
// keep it.
func (p *Parser) anText(start, end int) string {
	if start < 0 || end > len(p.source) || start > end {
		return ""
	}
	return p.source[start:end]
}

// anOption is a DDLOption for the text from start to the previous token.
func (p *Parser) anOption(key string, start int) *nodes.DDLOption {
	return &nodes.DDLOption{
		Key:   key,
		Value: p.anText(start, p.prev.End),
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	}
}

// parseAnalyticAlias reads [ AS ] alias. Without AS a name is the alias
// unless stop reports that it opens what follows.
func (p *Parser) parseAnalyticAlias(stop func() bool) (string, error) {
	if p.cur.Type == kwAS {
		p.advance()
		return p.anParseName()
	}
	if p.isAnalyticNameToken() && (stop == nil || !stop()) {
		return p.advance().Str, nil
	}
	return "", nil
}

// parseAnalyticIfNotExists reads IF NOT EXISTS before a name. IF followed
// by NOT or EXISTS opens the clause, and IF NOT without EXISTS or IF EXISTS
// is ORA-11543 at IF; IF before anything else is a name.
func (p *Parser) parseAnalyticIfNotExists() (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	next := p.peekNext().Type
	if next != kwNOT && next != kwEXISTS {
		return false, nil
	}
	ifTok := p.cur
	p.advance() // IF
	if next == kwEXISTS {
		return false, p.syntaxErrorAtTok(ifTok)
	}
	p.advance() // NOT
	if p.cur.Type != kwEXISTS {
		return false, p.syntaxErrorAtTok(ifTok)
	}
	p.advance()
	return true, nil
}

// analyticIfExists reads IF EXISTS before a name, if there. IF not
// followed by EXISTS is the name IF.
func (p *Parser) analyticIfExists() bool {
	if p.cur.Type != kwIF || p.peekNext().Type != kwEXISTS {
		return false
	}
	p.advance() // IF
	p.advance() // EXISTS
	return true
}

// parseAnalyticSharing reads SHARING = { METADATA | NONE }, if there.
func (p *Parser) parseAnalyticSharing() (string, error) {
	if !p.isKeywordStr("SHARING") {
		return "", nil
	}
	p.advance()
	if err := p.anExpect('='); err != nil {
		return "", err
	}
	if p.isKeywordStr("METADATA") || p.isKeywordStr("NONE") {
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// isAnalyticClassificationStart reports whether the current token opens a
// classification_clause.
func (p *Parser) isAnalyticClassificationStart() bool {
	return p.isKeywordStr("CAPTION") || p.isKeywordStr("DESCRIPTION") || p.isKeywordStr("CLASSIFICATION")
}

// parseAnalyticClassifications reads classification_clause:
//
//	[ CAPTION 'caption' ] [ DESCRIPTION 'description' ]
//	[ CLASSIFICATION name [ VALUE 'value' ] [ LANGUAGE 'language' ] ]...
//
// It returns nil when there is none.
func (p *Parser) parseAnalyticClassifications() (*nodes.List, error) {
	var items []nodes.Node
	for _, key := range []string{"CAPTION", "DESCRIPTION"} {
		if !p.isKeywordStr(key) {
			continue
		}
		start := p.pos()
		p.advance()
		value, err := p.anParseString()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: start, End: p.prev.End}})
	}
	for p.isKeywordStr("CLASSIFICATION") {
		start := p.pos()
		p.advance()
		value, err := p.anParseName()
		if err != nil {
			return nil, err
		}
		for _, key := range []string{"VALUE", "LANGUAGE"} {
			if !p.isKeywordStr(key) {
				continue
			}
			p.advance()
			s, err := p.anParseString()
			if err != nil {
				return nil, err
			}
			value += " " + key + " " + s
		}
		items = append(items, &nodes.DDLOption{Key: "CLASSIFICATION", Value: value, Loc: nodes.Loc{Start: start, End: p.prev.End}})
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// skipAnalyticAlias reads [ [ AS ] alias ] the AST does not keep.
func (p *Parser) skipAnalyticAlias(stop func() bool) error {
	alias, err := p.parseAnalyticAlias(stop)
	if err != nil || alias == "" {
		return err
	}
	return nil
}

// skipAnalyticClassifications reads a classification_clause the AST does
// not keep.
func (p *Parser) skipAnalyticClassifications() error {
	list, err := p.parseAnalyticClassifications()
	if err != nil || list == nil {
		return err
	}
	return nil
}

// skipAnalyticList reads a parenthesized list the AST does not keep.
func (p *Parser) skipAnalyticList(item func() error) error {
	list, err := p.parseAnalyticList(item)
	if err != nil || list == nil {
		return err
	}
	return nil
}

// skipAnalyticAggrFunction reads an aggr_function the AST does not keep.
func (p *Parser) skipAnalyticAggrFunction() error {
	fn, err := p.parseAnalyticAggrFunction()
	if err != nil || fn == "" {
		return err
	}
	return nil
}

// anAppendItems appends the items of list to dst.
func anAppendItems(dst *nodes.List, list *nodes.List) *nodes.List {
	if list == nil || len(list.Items) == 0 {
		return dst
	}
	if dst == nil {
		dst = &nodes.List{}
	}
	dst.Items = append(dst.Items, list.Items...)
	return dst
}

// isAnalyticObjectWord reports whether tok, unquoted, opens ANALYTIC VIEW,
// ATTRIBUTE DIMENSION or HIERARCHY, the objects CREATE takes FORCE or
// NOFORCE before.
func isAnalyticObjectWord(tok Token) bool {
	if tok.Type == tokQIDENT || tok.Type != tokIDENT && tok.Type < 2000 {
		return false
	}
	switch tok.Str {
	case "ANALYTIC", "ATTRIBUTE", "HIERARCHY":
		return true
	}
	return false
}

// parseCreateForcedAnalyticObject parses CREATE [ OR REPLACE ] { FORCE |
// NOFORCE } followed by ANALYTIC VIEW, ATTRIBUTE DIMENSION or HIERARCHY.
// Called at FORCE or NOFORCE.
func (p *Parser) parseCreateForcedAnalyticObject(start int, orReplace bool) (nodes.StmtNode, error) {
	force := p.cur.Type == kwFORCE
	p.advance() // FORCE or NOFORCE
	switch {
	case p.isKeywordStr("ANALYTIC"):
		p.advance()
		if p.cur.Type != kwVIEW {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return p.parseCreateAnalyticViewStmt(start, orReplace, force, !force)
	case p.isKeywordStr("ATTRIBUTE"):
		p.advance()
		if err := p.anExpectWord("DIMENSION"); err != nil {
			return nil, err
		}
		return p.parseCreateAttributeDimensionStmt(start, orReplace, force, !force)
	default: // HIERARCHY
		p.advance()
		return p.parseCreateHierarchyStmt(start, orReplace, force, !force)
	}
}

// ---------------------------------------------------------------------------
// CREATE / ALTER ANALYTIC VIEW
// ---------------------------------------------------------------------------

// parseCreateAnalyticViewStmt parses a CREATE ANALYTIC VIEW statement.
// Called after CREATE [ OR REPLACE ] [ FORCE | NOFORCE ] ANALYTIC VIEW.
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] ANALYTIC VIEW
//	    [ IF NOT EXISTS ] [ schema. ] analytic_view_name
//	    [ SHARING = { METADATA | NONE } ]
//	    [ classification_clause ]
//	    USING [ schema. ] { table | view } [ [ AS ] alias ]
//	    DIMENSION BY ( dim_key [, dim_key ]... )
//	    MEASURES ( av_measure [, av_measure ]... )
//	    [ DEFAULT MEASURE measure_name ]
//	    [ DEFAULT AGGREGATE BY aggr_function ]
//	    [ cache_clause ]
//	    [ FACT COLUMNS ( column [ [ AS ] alias ] [, ...] ) ]
//	    [ ENABLE QUERY TRANSFORM [ RELY | NORELY ] ]
func (p *Parser) parseCreateAnalyticViewStmt(start int, orReplace, force, noForce bool) (*nodes.CreateAnalyticViewStmt, error) {
	stmt := &nodes.CreateAnalyticViewStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}
	var err error
	if stmt.IfNotExists, err = p.parseAnalyticIfNotExists(); err != nil {
		return nil, err
	}
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.Sharing, err = p.parseAnalyticSharing(); err != nil {
		return nil, err
	}
	classifications, err := p.parseAnalyticClassifications()
	if err != nil {
		return nil, err
	}
	stmt.Options = anAppendItems(stmt.Options, classifications)

	// using_clause
	if err := p.anExpect(kwUSING); err != nil {
		return nil, err
	}
	if stmt.UsingTable, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.UsingAlias, err = p.parseAnalyticAlias(func() bool { return p.isKeywordStr("DIMENSION") }); err != nil {
		return nil, err
	}

	// dim_by_clause
	if err := p.anExpectWord("DIMENSION"); err != nil {
		return nil, err
	}
	if err := p.anExpect(kwBY); err != nil {
		return nil, err
	}
	if stmt.DimBy, err = p.parseAnalyticList(p.parseAnalyticDimKey); err != nil {
		return nil, err
	}

	// measures_clause
	if err := p.anExpectWord("MEASURES"); err != nil {
		return nil, err
	}
	if stmt.Measures, err = p.parseAnalyticList(p.parseAnalyticMeasure); err != nil {
		return nil, err
	}

	// default_measure_clause, default_aggregate_clause
	if p.cur.Type == kwDEFAULT && !p.anNextIsWord("AGGREGATE") {
		p.advance()
		if err := p.anExpectWord("MEASURE"); err != nil {
			return nil, err
		}
		if stmt.DefaultMeasure, err = p.anParseName(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if err := p.anExpectWord("AGGREGATE"); err != nil {
			return nil, err
		}
		if err := p.anExpect(kwBY); err != nil {
			return nil, err
		}
		if stmt.DefaultAggregate, err = p.parseAnalyticAggrFunction(); err != nil {
			return nil, err
		}
	}

	// cache_clause
	if p.isKeywordStr("CACHE") {
		optStart := p.pos()
		p.advance()
		for first := true; first || p.isKeywordStr("MEASURE"); first = false {
			if err := p.parseAnalyticMeasureGroup(); err != nil {
				return nil, err
			}
			for firstLevels := true; firstLevels || p.isKeywordStr("LEVELS"); firstLevels = false {
				if err := p.parseAnalyticCacheLevels(true, true); err != nil {
					return nil, err
				}
			}
		}
		stmt.Options = anAppendItems(stmt.Options, &nodes.List{Items: []nodes.Node{p.anOption("CACHE", optStart)}})
	}

	// fact_columns_clause
	if p.isKeywordStr("FACT") {
		optStart := p.pos()
		p.advance()
		if err := p.anExpectWord("COLUMNS"); err != nil {
			return nil, err
		}
		if err := p.skipAnalyticList(p.parseAnalyticFactColumn); err != nil {
			return nil, err
		}
		stmt.Options = anAppendItems(stmt.Options, &nodes.List{Items: []nodes.Node{p.anOption("FACT COLUMNS", optStart)}})
	}

	// qry_transform_clause
	if p.cur.Type == kwENABLE {
		optStart := p.pos()
		p.advance()
		if err := p.anExpectWord("QUERY"); err != nil {
			return nil, err
		}
		if err := p.anExpectWord("TRANSFORM"); err != nil {
			return nil, err
		}
		if p.isKeywordStr("RELY") || p.isKeywordStr("NORELY") {
			p.advance()
		}
		stmt.Options = anAppendItems(stmt.Options, &nodes.List{Items: []nodes.Node{p.anOption("QUERY TRANSFORM", optStart)}})
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAnalyticList reads ( item [, item ]... ), each item read by item, and
// returns the text of each.
func (p *Parser) parseAnalyticList(item func() error) (*nodes.List, error) {
	if err := p.anExpect('('); err != nil {
		return nil, err
	}
	list := &nodes.List{}
	for {
		start := p.pos()
		if err := item(); err != nil {
			return nil, err
		}
		list.Items = append(list.Items, &nodes.String{Str: p.anText(start, p.prev.End)})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseAnalyticDimKey reads a dim_key:
//
//	[ schema. ] attribute_dimension [ [ AS ] alias ] [ classification_clause ]
//	KEY { column | ( column [, column ]... ) }
//	REFERENCES [ DISTINCT ] { attribute | ( attribute [, attribute ]... ) }
//	HIERARCHIES ( hier_ref [, hier_ref ]... )
//
// REFERENCES names as many attributes as KEY names columns.
func (p *Parser) parseAnalyticDimKey() error {
	if _, err := p.anParseObjectName(); err != nil {
		return err
	}
	stop := func() bool { return p.isKeywordStr("KEY") || p.isAnalyticClassificationStart() }
	if err := p.skipAnalyticAlias(stop); err != nil {
		return err
	}
	if err := p.skipAnalyticClassifications(); err != nil {
		return err
	}
	if err := p.anExpectWord("KEY"); err != nil {
		return err
	}
	keys := 1
	if p.cur.Type == '(' {
		p.advance()
		for keys = 0; ; {
			if _, err := p.anParseDottedName(2); err != nil {
				return err
			}
			keys++
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.anExpect(')'); err != nil {
			return err
		}
	} else if _, err := p.anParseDottedName(2); err != nil {
		return err
	}
	if err := p.anExpect(kwREFERENCES); err != nil {
		return err
	}
	if p.cur.Type == kwDISTINCT {
		p.advance()
	}
	if keys == 1 && p.cur.Type != '(' {
		if _, err := p.anParseName(); err != nil {
			return err
		}
	} else {
		if err := p.anExpect('('); err != nil {
			return err
		}
		for i := 0; i < keys; i++ {
			if i > 0 {
				if err := p.anExpect(','); err != nil {
					return err
				}
			}
			if _, err := p.anParseName(); err != nil {
				return err
			}
		}
		if err := p.anExpect(')'); err != nil {
			return err
		}
	}
	if err := p.anExpectWord("HIERARCHIES"); err != nil {
		return err
	}
	return p.skipAnalyticList(p.parseAnalyticHierRef)
}

// parseAnalyticHierRef reads a hier_ref:
//
//	[ schema. ] hierarchy [ [ AS ] alias ] [ DEFAULT ]
func (p *Parser) parseAnalyticHierRef() error {
	if _, err := p.anParseObjectName(); err != nil {
		return err
	}
	if err := p.skipAnalyticAlias(nil); err != nil {
		return err
	}
	if p.cur.Type == kwDEFAULT {
		p.advance()
	}
	return nil
}

// parseAnalyticMeasure reads an av_measure:
//
//	measure_name
//	    { [ FACT expression ] [ AGGREGATE BY aggr_function ]
//	    | AS ( analytic_view_expression ) }
//	    [ classification_clause ]
func (p *Parser) parseAnalyticMeasure() error {
	if _, err := p.anParseName(); err != nil {
		return err
	}
	if p.cur.Type == kwAS {
		p.advance()
		if err := p.skipAnalyticViewExpression(); err != nil {
			return err
		}
	} else {
		if p.isKeywordStr("FACT") {
			p.advance()
			if expr, err := p.parseRequiredExpr(); err != nil || expr == nil {
				return err
			}
		}
		if p.isKeywordStr("AGGREGATE") {
			p.advance()
			if err := p.anExpect(kwBY); err != nil {
				return err
			}
			if err := p.skipAnalyticAggrFunction(); err != nil {
				return err
			}
		}
	}
	return p.skipAnalyticClassifications()
}

// skipAnalyticViewExpression reads ( analytic_view_expression ) as a whole:
// a parenthesized group that is not empty (AS () is ORA-00936).
func (p *Parser) skipAnalyticViewExpression() error {
	if err := p.anExpect('('); err != nil {
		return err
	}
	if p.cur.Type == ')' {
		return p.syntaxErrorAtCur()
	}
	for depth := 1; depth > 0; {
		switch p.cur.Type {
		case tokEOF, ';':
			return p.syntaxErrorAtCur()
		case '(':
			depth++
		case ')':
			depth--
		}
		p.advance()
	}
	return nil
}

// analyticAggrFunctions are the aggr_function words of an analytic view
// (ORA-00919 for any other).
var analyticAggrFunctions = map[string]bool{
	"APPROX_COUNT_DISTINCT": true, "APPROX_COUNT_DISTINCT_AGG": true, "AVG": true,
	"COUNT": true, "MAX": true, "MIN": true, "STDDEV": true, "STDDEV_POP": true,
	"STDDEV_SAMP": true, "SUM": true, "VAR_POP": true, "VAR_SAMP": true, "VARIANCE": true,
}

// parseAnalyticAggrFunction reads aggr_function: one of
// analyticAggrFunctions, unquoted, and COUNT may take DISTINCT.
func (p *Parser) parseAnalyticAggrFunction() (string, error) {
	if p.cur.Type == tokQIDENT || !p.isIdentLike() || !analyticAggrFunctions[p.cur.Str] {
		return "", p.syntaxErrorAtCur()
	}
	fn := p.advance().Str
	if fn == "COUNT" && p.cur.Type == kwDISTINCT {
		p.advance()
		fn += " DISTINCT"
	}
	return fn, nil
}

// parseAnalyticMeasureGroup reads MEASURE GROUP { ALL | ( measure_name [,
// measure_name ]... ) }.
func (p *Parser) parseAnalyticMeasureGroup() error {
	if err := p.anExpectWord("MEASURE"); err != nil {
		return err
	}
	if err := p.anExpect(kwGROUP); err != nil {
		return err
	}
	if p.cur.Type == kwALL {
		p.advance()
		return nil
	}
	return p.skipAnalyticList(func() error {
		_, err := p.anParseName()
		return err
	})
}

// parseAnalyticCacheLevels reads LEVELS ( [ level_specification [, ...] ] ),
// then MATERIALIZED [ USING [ schema. ] table ] when materialized is set.
// A level_specification is [ [ schema. ] hierarchy. ] level. The table
// after USING is read only where using is set.
func (p *Parser) parseAnalyticCacheLevels(materialized, using bool) error {
	if err := p.anExpectWord("LEVELS"); err != nil {
		return err
	}
	if err := p.anExpect('('); err != nil {
		return err
	}
	if p.cur.Type != ')' {
		for {
			if _, err := p.anParseDottedName(3); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if err := p.anExpect(')'); err != nil {
		return err
	}
	if !materialized {
		return nil
	}
	if err := p.anExpectWord("MATERIALIZED"); err != nil {
		return err
	}
	if using && p.cur.Type == kwUSING {
		p.advance()
		if _, err := p.anParseObjectName(); err != nil {
			return err
		}
	}
	return nil
}

// parseAnalyticFactColumn reads a fact column: [ alias. ] column
// [ [ AS ] alias ].
func (p *Parser) parseAnalyticFactColumn() error {
	if _, err := p.anParseDottedName(2); err != nil {
		return err
	}
	return p.skipAnalyticAlias(nil)
}

// parseAlterAnalyticViewStmt parses an ALTER ANALYTIC VIEW statement.
// Called after ALTER ANALYTIC VIEW.
//
//	ALTER ANALYTIC VIEW [ IF EXISTS ] [ schema. ] analytic_view_name
//	    { RENAME TO new_av_name
//	    | COMPILE
//	    | ADD CACHE MEASURE GROUP { ALL | ( measure [, measure ]... ) }
//	        LEVELS ( level_specification [, ...] ) MATERIALIZED [ USING table ]
//	    | DROP CACHE MEASURE GROUP { ALL | ( measure [, measure ]... ) }
//	        LEVELS ( level_specification [, ...] ) }
func (p *Parser) parseAlterAnalyticViewStmt(start int) (*nodes.AlterAnalyticViewStmt, error) {
	stmt := &nodes.AlterAnalyticViewStmt{Loc: nodes.Loc{Start: start}}
	stmt.IfExists = p.analyticIfExists()
	var err error
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	switch {
	case p.isKeywordStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwRENAME:
		stmt.Action = "RENAME"
		if stmt.NewName, err = p.parseAnalyticRenameTo(); err != nil {
			return nil, err
		}
	case p.cur.Type == kwADD || p.cur.Type == kwDROP:
		add := p.cur.Type == kwADD
		stmt.Action = "DROP_CACHE"
		if add {
			stmt.Action = "ADD_CACHE"
		}
		optStart := p.pos()
		p.advance()
		if err := p.anExpectWord("CACHE"); err != nil {
			return nil, err
		}
		if err := p.parseAnalyticMeasureGroup(); err != nil {
			return nil, err
		}
		if err := p.parseAnalyticCacheLevels(add, add); err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: []nodes.Node{p.anOption("CACHE", optStart)}}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAnalyticRenameTo reads RENAME TO new_name, a name with no schema.
func (p *Parser) parseAnalyticRenameTo() (*nodes.ObjectName, error) {
	p.advance() // RENAME
	if err := p.anExpect(kwTO); err != nil {
		return nil, err
	}
	start := p.pos()
	name, err := p.anParseName()
	if err != nil {
		return nil, err
	}
	return &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
}

// parseAnalyticRenameOrCompile reads the action of ALTER ATTRIBUTE
// DIMENSION and ALTER HIERARCHY: { RENAME TO new_name | COMPILE }.
func (p *Parser) parseAnalyticRenameOrCompile() (string, *nodes.ObjectName, error) {
	switch {
	case p.isKeywordStr("COMPILE"):
		p.advance()
		return "COMPILE", nil, nil
	case p.cur.Type == kwRENAME:
		name, err := p.parseAnalyticRenameTo()
		if err != nil {
			return "", nil, err
		}
		return "RENAME", name, nil
	}
	return "", nil, p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// CREATE / ALTER ATTRIBUTE DIMENSION
// ---------------------------------------------------------------------------

// attrDimLevelTypes are the words of LEVEL TYPE.
var attrDimLevelTypes = map[string]bool{
	"STANDARD": true, "YEARS": true, "HALF_YEARS": true, "QUARTERS": true, "MONTHS": true,
	"WEEKS": true, "DAYS": true, "HOURS": true, "MINUTES": true, "SECONDS": true,
}

// isAttrDimLevelTypeWord reports whether tok is a LEVEL TYPE word, unquoted.
func isAttrDimLevelTypeWord(tok Token) bool {
	return tok.Type != tokQIDENT && (tok.Type == tokIDENT || tok.Type >= 2000) && attrDimLevelTypes[tok.Str]
}

// parseCreateAttributeDimensionStmt parses a CREATE ATTRIBUTE DIMENSION
// statement. Called after CREATE [ OR REPLACE ] [ FORCE | NOFORCE ]
// ATTRIBUTE DIMENSION.
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] ATTRIBUTE DIMENSION
//	    [ IF NOT EXISTS ] [ schema. ] attr_dimension
//	    [ SHARING = { METADATA | NONE } ]
//	    [ DIMENSION TYPE { STANDARD | TIME } ]
//	    [ classification_clause ]
//	    USING source_clause [, source_clause ]... [ join_path_clause ]...
//	    ATTRIBUTES ( attr_dim_attribute_clause [, ...] )
//	    attr_dim_level_clause [ attr_dim_level_clause ]...
//	    [ all_clause ]
//
//	source_clause: [ schema. ] { table | view } [ [ AS ] alias ]
//	join_path_clause: JOIN PATH join_path_name ON join_condition
//	join_condition: join_condition_elem [ AND join_condition_elem ]...
//	join_condition_elem: [ alias. ] column = [ alias. ] column
func (p *Parser) parseCreateAttributeDimensionStmt(start int, orReplace, force, noForce bool) (*nodes.CreateAttributeDimensionStmt, error) {
	stmt := &nodes.CreateAttributeDimensionStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}
	var err error
	if stmt.IfNotExists, err = p.parseAnalyticIfNotExists(); err != nil {
		return nil, err
	}
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.Sharing, err = p.parseAnalyticSharing(); err != nil {
		return nil, err
	}
	if p.isKeywordStr("DIMENSION") {
		p.advance()
		if err := p.anExpectWord("TYPE"); err != nil {
			return nil, err
		}
		if !p.isKeywordStr("STANDARD") && !p.isKeywordStr("TIME") {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.DimensionType = p.advance().Str
	}
	stmt.Classifications, err = p.parseAnalyticClassifications()
	if err != nil {
		return nil, err
	}

	// attr_dim_using_clause
	if err := p.anExpect(kwUSING); err != nil {
		return nil, err
	}
	stmt.Sources = &nodes.List{}
	for {
		src, err := p.parseAttrDimSource()
		if err != nil {
			return nil, err
		}
		stmt.Sources.Items = append(stmt.Sources.Items, src)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	for p.cur.Type == kwJOIN && p.anNextIsWord("PATH") {
		src, err := p.parseAttrDimJoinPath()
		if err != nil {
			return nil, err
		}
		stmt.Sources.Items = append(stmt.Sources.Items, src)
	}

	// attributes_clause
	if err := p.anExpectWord("ATTRIBUTES"); err != nil {
		return nil, err
	}
	if err := p.anExpect('('); err != nil {
		return nil, err
	}
	stmt.Attributes = &nodes.List{}
	for {
		attr, err := p.parseAttrDimAttribute()
		if err != nil {
			return nil, err
		}
		stmt.Attributes.Items = append(stmt.Attributes.Items, attr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}

	// attr_dim_level_clause...
	if p.cur.Type != kwLEVEL {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Levels = &nodes.List{}
	for p.cur.Type == kwLEVEL {
		lvl, err := p.parseAttrDimLevel()
		if err != nil {
			return nil, err
		}
		stmt.Levels.Items = append(stmt.Levels.Items, lvl)
	}

	// all_clause
	if p.cur.Type == kwALL {
		allc := &nodes.AttrDimAllClause{Loc: nodes.Loc{Start: p.pos()}}
		p.advance()
		if !p.isKeywordStr("MEMBER") {
			return nil, p.syntaxErrorAtCur()
		}
		if allc.MemberName, allc.MemberCaption, allc.MemberDesc, err = p.parseAttrDimMembers(); err != nil {
			return nil, err
		}
		allc.Loc.End = p.prev.End
		stmt.AllClause = allc
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAttrDimSource reads a source_clause: [ schema. ] { table | view }
// [ [ AS ] alias ].
func (p *Parser) parseAttrDimSource() (*nodes.AttrDimSourceClause, error) {
	src := &nodes.AttrDimSourceClause{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	if src.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	stop := func() bool {
		return p.isKeywordStr("ATTRIBUTES") || p.cur.Type == kwJOIN && p.anNextIsWord("PATH")
	}
	if src.Alias, err = p.parseAnalyticAlias(stop); err != nil {
		return nil, err
	}
	src.Loc.End = p.prev.End
	return src, nil
}

// parseAttrDimJoinPath reads a join_path_clause.
func (p *Parser) parseAttrDimJoinPath() (*nodes.AttrDimSourceClause, error) {
	src := &nodes.AttrDimSourceClause{IsJoinPath: true, Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // JOIN
	p.advance() // PATH
	var err error
	if src.JoinPathName, err = p.anParseName(); err != nil {
		return nil, err
	}
	if err := p.anExpect(kwON); err != nil {
		return nil, err
	}
	src.JoinCondition = &nodes.List{}
	for {
		elem := &nodes.AttrDimJoinCondElem{Loc: nodes.Loc{Start: p.pos()}}
		if elem.LeftTable, elem.LeftCol, err = p.parseAttrDimJoinColumn(); err != nil {
			return nil, err
		}
		if err := p.anExpect('='); err != nil {
			return nil, err
		}
		if elem.RightTable, elem.RightCol, err = p.parseAttrDimJoinColumn(); err != nil {
			return nil, err
		}
		elem.Loc.End = p.prev.End
		src.JoinCondition.Items = append(src.JoinCondition.Items, elem)
		if p.cur.Type != kwAND {
			break
		}
		p.advance()
	}
	src.Loc.End = p.prev.End
	return src, nil
}

// parseAttrDimJoinColumn reads [ alias. ] column.
func (p *Parser) parseAttrDimJoinColumn() (table, column string, err error) {
	if column, err = p.anParseName(); err != nil {
		return "", "", err
	}
	if p.cur.Type == '.' {
		p.advance()
		table = column
		if column, err = p.anParseName(); err != nil {
			return "", "", err
		}
	}
	return table, column, nil
}

// parseAttrDimAttribute reads an attr_dim_attribute_clause:
//
//	[ alias. ] column [ [ AS ] alias ] [ classification_clause ]
func (p *Parser) parseAttrDimAttribute() (*nodes.AttrDimAttribute, error) {
	attr := &nodes.AttrDimAttribute{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	if attr.Column, err = p.anParseDottedName(2); err != nil {
		return nil, err
	}
	if attr.Alias, err = p.parseAnalyticAlias(p.isAnalyticClassificationStart); err != nil {
		return nil, err
	}
	if attr.Classifications, err = p.parseAnalyticClassifications(); err != nil {
		return nil, err
	}
	attr.Loc.End = p.prev.End
	return attr, nil
}

// parseAttrDimLevel reads an attr_dim_level_clause:
//
//	LEVEL level_name
//	[ LEVEL TYPE { STANDARD | YEARS | HALF_YEARS | QUARTERS | MONTHS
//	             | WEEKS | DAYS | HOURS | MINUTES | SECONDS } ]
//	[ classification_clause ]
//	KEY { attribute | ( attribute [, attribute ]... ) }
//	[ ALTERNATE KEY { attribute | ( attribute [, attribute ]... ) } ]...
//	[ MEMBER NAME expression ] [ MEMBER CAPTION expression ]
//	[ MEMBER DESCRIPTION expression ]
//	[ ORDER BY dim_order_item [, dim_order_item ]... ]
//	[ DETERMINES ( attribute [, attribute ]... ) ]
//
// LEVEL after the level name opens LEVEL TYPE. Oracle reads TYPE after
// LEVEL as LEVEL TYPE when a level type word follows, so LEVEL TYPE YEARS
// cannot open a level named TYPE.
func (p *Parser) parseAttrDimLevel() (*nodes.AttrDimLevel, error) {
	lvl := &nodes.AttrDimLevel{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // LEVEL
	if p.isKeywordStr("TYPE") && isAttrDimLevelTypeWord(p.peekNext()) {
		return nil, p.syntaxErrorAtCur()
	}
	var err error
	if lvl.Name, err = p.anParseName(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwLEVEL {
		p.advance() // LEVEL
		if err := p.anExpectWord("TYPE"); err != nil {
			return nil, err
		}
		if !isAttrDimLevelTypeWord(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		lvl.LevelType = p.advance().Str
	}
	if lvl.Classifications, err = p.parseAnalyticClassifications(); err != nil {
		return nil, err
	}
	if err := p.anExpectWord("KEY"); err != nil {
		return nil, err
	}
	if lvl.KeyAttrs, err = p.parseAttrDimNames(true); err != nil {
		return nil, err
	}
	for p.isKeywordStr("ALTERNATE") {
		p.advance()
		if err := p.anExpectWord("KEY"); err != nil {
			return nil, err
		}
		alt, err := p.parseAttrDimNames(true)
		if err != nil {
			return nil, err
		}
		lvl.AltKeyAttrs = anAppendItems(lvl.AltKeyAttrs, alt)
	}
	if lvl.MemberName, lvl.MemberCaption, lvl.MemberDesc, err = p.parseAttrDimMembers(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwORDER {
		p.advance()
		if err := p.anExpect(kwBY); err != nil {
			return nil, err
		}
		lvl.OrderByAttrs = &nodes.List{}
		for {
			item, err := p.parseAttrDimOrderItem()
			if err != nil {
				return nil, err
			}
			lvl.OrderByAttrs.Items = append(lvl.OrderByAttrs.Items, item)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if p.isKeywordStr("DETERMINES") {
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		if lvl.Determines, err = p.parseAttrDimNames(false); err != nil {
			return nil, err
		}
	}
	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// parseAttrDimNames reads ( attribute [, attribute ]... ), or a single
// attribute with no parentheses where single is set.
func (p *Parser) parseAttrDimNames(single bool) (*nodes.List, error) {
	list := &nodes.List{}
	if single && p.cur.Type != '(' {
		name, err := p.anParseName()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		return list, nil
	}
	if err := p.anExpect('('); err != nil {
		return nil, err
	}
	for {
		name, err := p.anParseName()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseAttrDimMembers reads [ MEMBER NAME expression ] [ MEMBER CAPTION
// expression ] [ MEMBER DESCRIPTION expression ], in that order. MEMBER
// opens one of them: MEMBER 1 fails at 1. A MEMBER whose word comes out
// of order is left for the caller, which fails at it.
func (p *Parser) parseAttrDimMembers() (name, caption, desc nodes.ExprNode, err error) {
	exprs := []*nodes.ExprNode{&name, &caption, &desc}
	words := []string{"NAME", "CAPTION", "DESCRIPTION"}
	for next := 0; p.isKeywordStr("MEMBER"); {
		i := next
		for i < len(words) && !p.anNextIsWord(words[i]) {
			i++
		}
		if i == len(words) {
			for _, w := range words[:next] {
				if p.anNextIsWord(w) {
					return name, caption, desc, nil
				}
			}
			p.advance() // MEMBER
			return nil, nil, nil, p.syntaxErrorAtCur()
		}
		p.advance() // MEMBER
		p.advance() // word
		if *exprs[i], err = p.parseRequiredExpr(); err != nil {
			return nil, nil, nil, err
		}
		next = i + 1
	}
	return name, caption, desc, nil
}

// parseAttrDimOrderItem reads a dim_order_clause item:
//
//	[ MIN | MAX ] attribute [ ASC | DESC ] [ NULLS { FIRST | LAST } ]
func (p *Parser) parseAttrDimOrderItem() (*nodes.AttrDimOrderByItem, error) {
	item := &nodes.AttrDimOrderByItem{Loc: nodes.Loc{Start: p.pos()}}
	if p.isKeywordStr("MIN") || p.isKeywordStr("MAX") {
		p.advance()
	}
	var err error
	if item.Attribute, err = p.anParseName(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwASC || p.cur.Type == kwDESC {
		item.Direction = p.advance().Str
	}
	if p.isKeywordStr("NULLS") {
		p.advance()
		if !p.isKeywordStr("FIRST") && !p.isKeywordStr("LAST") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}
	item.Loc.End = p.prev.End
	return item, nil
}

// parseAlterAttributeDimensionStmt parses an ALTER ATTRIBUTE DIMENSION
// statement. Called after ALTER ATTRIBUTE DIMENSION.
//
//	ALTER ATTRIBUTE DIMENSION [ IF EXISTS ] [ schema. ] attr_dim_name
//	    { RENAME TO new_attr_dim_name | COMPILE }
func (p *Parser) parseAlterAttributeDimensionStmt(start int) (*nodes.AlterAttributeDimensionStmt, error) {
	stmt := &nodes.AlterAttributeDimensionStmt{Loc: nodes.Loc{Start: start}}
	stmt.IfExists = p.analyticIfExists()
	var err error
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.Action, stmt.NewName, err = p.parseAnalyticRenameOrCompile(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE / ALTER HIERARCHY
// ---------------------------------------------------------------------------

// parseCreateHierarchyStmt parses a CREATE HIERARCHY statement. Called after
// CREATE [ OR REPLACE ] [ FORCE | NOFORCE ] HIERARCHY.
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] HIERARCHY
//	    [ IF NOT EXISTS ] [ schema. ] hierarchy
//	    [ SHARING = { METADATA | NONE } ]
//	    [ classification_clause ]
//	    USING [ schema. ] attr_dimension
//	    ( level_hier_clause )
//	    [ HIERARCHICAL ATTRIBUTES ( hier_attr_clause [, ...] ) ]
//
//	level_hier_clause: level_name [ CHILD OF level_hier_clause ]
//	hier_attr_clause: hier_attr_name [ classification_clause ]
func (p *Parser) parseCreateHierarchyStmt(start int, orReplace, force, noForce bool) (*nodes.CreateHierarchyStmt, error) {
	stmt := &nodes.CreateHierarchyStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}
	var err error
	if stmt.IfNotExists, err = p.parseAnalyticIfNotExists(); err != nil {
		return nil, err
	}
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.Sharing, err = p.parseAnalyticSharing(); err != nil {
		return nil, err
	}
	stmt.Classifications, err = p.parseAnalyticClassifications()
	if err != nil {
		return nil, err
	}
	if err := p.anExpect(kwUSING); err != nil {
		return nil, err
	}
	if stmt.UsingAttrDim, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if err := p.anExpect('('); err != nil {
		return nil, err
	}
	if stmt.LevelHier, err = p.parseHierarchyLevelChain(); err != nil {
		return nil, err
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}
	if p.isKeywordStr("HIERARCHICAL") {
		p.advance()
		if err := p.anExpectWord("ATTRIBUTES"); err != nil {
			return nil, err
		}
		if err := p.anExpect('('); err != nil {
			return nil, err
		}
		stmt.HierAttrs = &nodes.List{}
		for {
			ha := &nodes.HierAttr{Loc: nodes.Loc{Start: p.pos()}}
			if ha.Name, err = p.anParseName(); err != nil {
				return nil, err
			}
			if ha.Classifications, err = p.parseAnalyticClassifications(); err != nil {
				return nil, err
			}
			ha.Loc.End = p.prev.End
			stmt.HierAttrs.Items = append(stmt.HierAttrs.Items, ha)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.anExpect(')'); err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseHierarchyLevelChain reads level_hier_clause: level_name [ CHILD OF
// level_hier_clause ].
func (p *Parser) parseHierarchyLevelChain() (*nodes.HierLevelClause, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	lvl := &nodes.HierLevelClause{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	if lvl.Name, err = p.anParseName(); err != nil {
		return nil, err
	}
	if p.isKeywordStr("CHILD") {
		p.advance()
		if err := p.anExpect(kwOF); err != nil {
			return nil, err
		}
		if lvl.ChildOf, err = p.parseHierarchyLevelChain(); err != nil {
			return nil, err
		}
	}
	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// parseAlterHierarchyStmt parses an ALTER HIERARCHY statement. Called after
// ALTER HIERARCHY.
//
//	ALTER HIERARCHY [ IF EXISTS ] [ schema. ] hierarchy_name
//	    { RENAME TO new_hier_name | COMPILE }
func (p *Parser) parseAlterHierarchyStmt(start int) (*nodes.AlterHierarchyStmt, error) {
	stmt := &nodes.AlterHierarchyStmt{Loc: nodes.Loc{Start: start}}
	stmt.IfExists = p.analyticIfExists()
	var err error
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if stmt.Action, stmt.NewName, err = p.parseAnalyticRenameOrCompile(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE / ALTER DIMENSION
// ---------------------------------------------------------------------------

// parseCreateDimensionStmt parses a CREATE DIMENSION statement. Called after
// CREATE DIMENSION.
//
//	CREATE DIMENSION [ schema. ] dimension
//	    level_clause [ level_clause ]...
//	    [ { hierarchy_clause | attribute_clause | extended_attribute_clause } ]...
//
//	level_clause:
//	    LEVEL level IS { level_column | ( level_column [, column ]... ) }
//	    [ SKIP WHEN NULL ]
//	hierarchy_clause:
//	    HIERARCHY hierarchy ( child_level CHILD OF parent_level
//	        [ CHILD OF parent_level ]... [ dimension_join_clause ]... )
//	dimension_join_clause:
//	    JOIN KEY { column | ( column [, column ]... ) } REFERENCES parent_level
//	attribute_clause, extended_attribute_clause:
//	    ATTRIBUTE attribute [ LEVEL level ]
//	    DETERMINES { column | ( column [, column ]... ) }
//
// A level_column is [ schema. ] table. column; the other columns may leave
// out their table.
func (p *Parser) parseCreateDimensionStmt(start int) (*nodes.CreateDimensionStmt, error) {
	stmt := &nodes.CreateDimensionStmt{Loc: nodes.Loc{Start: start}}
	var err error
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	if p.cur.Type != kwLEVEL {
		return nil, p.syntaxErrorAtCur()
	}
	for p.cur.Type == kwLEVEL {
		lvl, err := p.parseDimensionLevelClause()
		if err != nil {
			return nil, err
		}
		stmt.Levels = append(stmt.Levels, lvl)
	}
	for {
		switch {
		case p.isKeywordStr("HIERARCHY"):
			hier, err := p.parseDimensionHierarchyClause()
			if err != nil {
				return nil, err
			}
			stmt.Hierarchies = append(stmt.Hierarchies, hier)
		case p.isKeywordStr("ATTRIBUTE"):
			attr, err := p.parseDimensionAttributeClause()
			if err != nil {
				return nil, err
			}
			stmt.Attributes = append(stmt.Attributes, attr)
		default:
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
	}
}

// parseDimensionColumn reads a column of a dimension clause: name [ . name
// ]..., with the number of names.
func (p *Parser) parseDimensionColumn() (*nodes.ObjectName, int, error) {
	start := p.pos()
	first, err := p.anParseName()
	if err != nil {
		return nil, 0, err
	}
	parts := []string{first}
	for p.cur.Type == '.' {
		p.advance()
		part, err := p.anParseName()
		if err != nil {
			return nil, 0, err
		}
		parts = append(parts, part)
	}
	obj := &nodes.ObjectName{Name: parts[len(parts)-1], Loc: nodes.Loc{Start: start, End: p.prev.End}}
	if len(parts) > 1 {
		obj.Schema = strings.Join(parts[:len(parts)-1], ".")
	}
	return obj, len(parts), nil
}

// parseDimensionColumns reads { column | ( column [, column ]... ) }. Where
// qualified is set, the first column names its table (ORA-30347).
func (p *Parser) parseDimensionColumns(qualified bool) ([]*nodes.ObjectName, error) {
	column := func(first bool) (*nodes.ObjectName, error) {
		col, parts, err := p.parseDimensionColumn()
		if err != nil {
			return nil, err
		}
		if first && qualified && parts < 2 {
			return nil, p.syntaxErrorAtCur()
		}
		return col, nil
	}
	if p.cur.Type != '(' {
		col, err := column(true)
		if err != nil {
			return nil, err
		}
		return []*nodes.ObjectName{col}, nil
	}
	p.advance()
	var cols []*nodes.ObjectName
	for {
		col, err := column(len(cols) == 0)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseDimensionLevelClause reads a level_clause.
func (p *Parser) parseDimensionLevelClause() (*nodes.DimensionLevel, error) {
	lvl := &nodes.DimensionLevel{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // LEVEL
	var err error
	if lvl.Name, err = p.anParseName(); err != nil {
		return nil, err
	}
	if err := p.anExpect(kwIS); err != nil {
		return nil, err
	}
	if lvl.Columns, err = p.parseDimensionColumns(true); err != nil {
		return nil, err
	}
	if p.isKeywordStr("SKIP") {
		p.advance()
		if err := p.anExpect(kwWHEN); err != nil {
			return nil, err
		}
		if err := p.anExpect(kwNULL); err != nil {
			return nil, err
		}
		lvl.SkipWhenNull = true
	}
	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// parseDimensionHierarchyClause reads a hierarchy_clause.
func (p *Parser) parseDimensionHierarchyClause() (*nodes.DimensionHierarchy, error) {
	hier := &nodes.DimensionHierarchy{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // HIERARCHY
	var err error
	if hier.Name, err = p.anParseName(); err != nil {
		return nil, err
	}
	if err := p.anExpect('('); err != nil {
		return nil, err
	}
	child, err := p.anParseName()
	if err != nil {
		return nil, err
	}
	hier.Levels = append(hier.Levels, child)
	for first := true; first || p.isKeywordStr("CHILD"); first = false {
		if err := p.anExpectWord("CHILD"); err != nil {
			return nil, err
		}
		if err := p.anExpect(kwOF); err != nil {
			return nil, err
		}
		parent, err := p.anParseName()
		if err != nil {
			return nil, err
		}
		hier.Levels = append(hier.Levels, parent)
	}
	for p.cur.Type == kwJOIN {
		jk := &nodes.DimensionJoinKey{Loc: nodes.Loc{Start: p.pos()}}
		p.advance()
		if err := p.anExpectWord("KEY"); err != nil {
			return nil, err
		}
		if jk.ChildKeys, err = p.parseDimensionColumns(false); err != nil {
			return nil, err
		}
		if err := p.anExpect(kwREFERENCES); err != nil {
			return nil, err
		}
		if jk.ParentLevel, err = p.anParseName(); err != nil {
			return nil, err
		}
		jk.Loc.End = p.prev.End
		hier.JoinKeys = append(hier.JoinKeys, jk)
	}
	if err := p.anExpect(')'); err != nil {
		return nil, err
	}
	hier.Loc.End = p.prev.End
	return hier, nil
}

// parseDimensionAttributeClause reads an attribute_clause or an
// extended_attribute_clause.
func (p *Parser) parseDimensionAttributeClause() (*nodes.DimensionAttribute, error) {
	attr := &nodes.DimensionAttribute{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // ATTRIBUTE
	var err error
	if attr.AttrName, err = p.anParseName(); err != nil {
		return nil, err
	}
	if p.cur.Type == kwLEVEL {
		p.advance()
		if attr.LevelName, err = p.anParseName(); err != nil {
			return nil, err
		}
	}
	if err := p.anExpectWord("DETERMINES"); err != nil {
		return nil, err
	}
	if attr.Columns, err = p.parseDimensionColumns(false); err != nil {
		return nil, err
	}
	attr.Loc.End = p.prev.End
	return attr, nil
}

// parseAlterDimensionStmt parses an ALTER DIMENSION statement. Called after
// ALTER DIMENSION.
//
//	ALTER DIMENSION [ schema. ] dimension
//	    { ADD { level_clause | hierarchy_clause | attribute_clause
//	          | extended_attribute_clause } }...
//	  | { DROP { LEVEL level [ RESTRICT | CASCADE ]
//	           | HIERARCHY hierarchy
//	           | ATTRIBUTE attribute [ LEVEL level [ COLUMN column ] ] } }...
//	  | COMPILE
//
// COMPILE may stand among the ADD or DROP clauses; ADD and DROP together
// are ORA-30348.
func (p *Parser) parseAlterDimensionStmt(start int) (*nodes.AlterDimensionStmt, error) {
	stmt := &nodes.AlterDimensionStmt{Loc: nodes.Loc{Start: start}}
	var err error
	if stmt.Name, err = p.anParseObjectName(); err != nil {
		return nil, err
	}
	added, dropped := false, false
	for {
		switch {
		case p.isKeywordStr("COMPILE"):
			p.advance()
			stmt.Compile = true
		case p.cur.Type == kwADD:
			if dropped {
				return nil, p.syntaxErrorAtCur()
			}
			added = true
			p.advance()
			if err := p.parseAlterDimensionAdd(stmt); err != nil {
				return nil, err
			}
		case p.cur.Type == kwDROP:
			if added {
				return nil, p.syntaxErrorAtCur()
			}
			dropped = true
			p.advance()
			if err := p.parseAlterDimensionDrop(stmt); err != nil {
				return nil, err
			}
		default:
			if !added && !dropped && !stmt.Compile {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
	}
}

// parseAlterDimensionAdd reads the clause after ADD.
func (p *Parser) parseAlterDimensionAdd(stmt *nodes.AlterDimensionStmt) error {
	switch {
	case p.cur.Type == kwLEVEL:
		lvl, err := p.parseDimensionLevelClause()
		if err != nil {
			return err
		}
		stmt.AddLevels = append(stmt.AddLevels, lvl)
	case p.isKeywordStr("HIERARCHY"):
		hier, err := p.parseDimensionHierarchyClause()
		if err != nil {
			return err
		}
		stmt.AddHierarchies = append(stmt.AddHierarchies, hier)
	case p.isKeywordStr("ATTRIBUTE"):
		attr, err := p.parseDimensionAttributeClause()
		if err != nil {
			return err
		}
		stmt.AddAttributes = append(stmt.AddAttributes, attr)
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterDimensionDrop reads the clause after DROP.
func (p *Parser) parseAlterDimensionDrop(stmt *nodes.AlterDimensionStmt) error {
	switch {
	case p.cur.Type == kwLEVEL:
		p.advance()
		name, err := p.anParseName()
		if err != nil {
			return err
		}
		if p.isKeywordStr("RESTRICT") || p.isKeywordStr("CASCADE") {
			p.advance()
		}
		stmt.DropLevels = append(stmt.DropLevels, name)
	case p.isKeywordStr("HIERARCHY"):
		p.advance()
		name, err := p.anParseName()
		if err != nil {
			return err
		}
		stmt.DropHierarchies = append(stmt.DropHierarchies, name)
	case p.isKeywordStr("ATTRIBUTE"):
		p.advance()
		name, err := p.anParseName()
		if err != nil {
			return err
		}
		if p.cur.Type == kwLEVEL {
			p.advance()
			if _, err := p.anParseName(); err != nil {
				return err
			}
			if p.cur.Type == kwCOLUMN {
				p.advance()
				if _, err := p.anParseName(); err != nil {
					return err
				}
			}
		}
		stmt.DropAttributes = append(stmt.DropAttributes, name)
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}
