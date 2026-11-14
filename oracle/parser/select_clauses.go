package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the clauses of a query's FROM items by the grammar of
// the SELECT page (bnf/SELECT.bnf, from
// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/SELECT.html),
// where Oracle 23ai reads them as documented, and as the engine reads them
// where it does not.

// takeSelectStar takes the * at the current token.
func (p *Parser) takeSelectStar() *nodes.Star {
	start := p.pos()
	p.advance()
	return &nodes.Star{Loc: nodes.Loc{Start: start, End: p.prev.End}}
}

// parseSelectWindowSpec parses the parenthesized window_specification of a
// WINDOW clause.
//
//	window_specification:
//	    [ existing_window_name ]
//	    [ query_partition_clause ]
//	    [ order_by_clause ]
//	    [ windowing_clause ]
//
// The existing window is one named earlier in the clause, a key of
// ordered (ORA-32784 at its name), and a window built on it takes no
// PARTITION BY of its own (ORA-32783 at the window's name, nameTok). A
// windowing clause needs an ORDER BY, its own or an existing window's,
// which ordered records for each window (ORA-30485 at ROWS, RANGE or
// GROUPS).
func (p *Parser) parseSelectWindowSpec(nameTok Token, ordered map[string]bool) (*nodes.WindowSpec, error) {
	p.advance() // consume '('
	ws := &nodes.WindowSpec{Loc: nodes.Loc{Start: p.pos()}}
	var err error

	based, baseOrdered := false, false
	if p.isIdentLike() && !p.isSelectWindowSpecKeyword() {
		if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
			return nil, err
		}
		refTok := p.cur
		if ws.WindowName, err = p.parseIdentifier(); err != nil {
			return nil, err
		}
		if baseOrdered, based = ordered[ws.WindowName]; !based {
			return nil, p.syntaxErrorAtTok(refTok)
		}
	}

	if p.cur.Type == kwPARTITION {
		if based {
			return nil, p.syntaxErrorAtTok(nameTok)
		}
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		ws.PartitionBy = &nodes.List{}
		for {
			if p.cur.Type == '(' && p.peekNext().Type == ')' {
				p.advance()
				return nil, p.syntaxErrorAtCur()
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			ws.PartitionBy.Items = append(ws.PartitionBy.Items, expr)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == kwORDER {
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if ws.OrderBy, err = p.parseOrderByList(); err != nil {
			return nil, err
		}
		if ws.OrderBy.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		if ws.OrderBy == nil && !baseOrdered {
			return nil, p.syntaxErrorAtCur()
		}
		if ws.Frame, err = p.parseWindowFrame(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	ws.Loc.End = p.prev.End
	return ws, nil
}

// isSelectWindowSpecKeyword reports whether the current token opens a
// clause of a window specification rather than naming an existing window.
func (p *Parser) isSelectWindowSpecKeyword() bool {
	switch p.cur.Type {
	case kwPARTITION, kwORDER, kwROWS, kwRANGE, kwGROUPS:
		return true
	}
	return false
}

// parseSelectTableFunctionRef parses a table function called in FROM
// without TABLE, name ( args ) [ t_alias ], which Oracle 23ai takes as it
// takes TABLE ( name ( args ) ).
func (p *Parser) parseSelectTableFunctionRef(name *nodes.ObjectName, start int) (nodes.TableExpr, error) {
	call, err := p.parseFuncCall(name.Name, name.Schema, start)
	if err != nil {
		return nil, err
	}
	tc := &nodes.TableCollectionExpr{Expr: call, Loc: nodes.Loc{Start: start}}
	if fn, ok := call.(*nodes.FuncCallExpr); ok {
		tc.FunctionCall = fn
	} else {
		tc.CollectionExpr = call
	}
	if p.isTableAliasCandidate() {
		if tc.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	}
	tc.Loc.End = p.prev.End
	return p.parsePivotSuffix(tc)
}

// parseSelectQueryPartition parses a query_partition_clause, at PARTITION:
//
//	PARTITION BY { expr [, expr ]... | ( expr [, expr ]... ) }
//
// Oracle 23ai takes only the parenthesized form before an outer join
// (parenthesized, ORA-00906 at the first expression otherwise).
func (p *Parser) parseSelectQueryPartition(parenthesized bool) (*nodes.List, error) {
	p.advance() // consume PARTITION
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if parenthesized {
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		list, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if list == nil || list.Len() == 0 || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return list, nil
	}
	list, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	if list == nil || list.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return list, nil
}

// isSelectQueryPartitionStart reports whether the current token opens a
// query_partition_clause: PARTITION BY.
func (p *Parser) isSelectQueryPartitionStart() bool {
	return p.cur.Type == kwPARTITION && p.peekNext().Type == kwBY
}

// parseSelectWithDeclarations parses the plsql_declarations of a WITH
// clause: function and procedure definitions, each ending in its ';'.
func (p *Parser) parseSelectWithDeclarations() (*nodes.List, error) {
	defer p.restorePLSQLCode(p.setPLSQLCode(true))
	decls := &nodes.List{}
	for p.cur.Type == kwFUNCTION || p.cur.Type == kwPROCEDURE {
		var decl nodes.Node
		var err error
		if p.cur.Type == kwFUNCTION {
			decl, err = p.parsePackageFuncDecl(subprogramNested)
		} else {
			decl, err = p.parsePackageProcDecl(subprogramNested)
		}
		if err != nil {
			return nil, err
		}
		if p.prev.Type != ';' {
			return nil, p.syntaxErrorAtCur()
		}
		decls.Items = append(decls.Items, decl)
	}
	return decls, nil
}

// parseInlineExternalTable parses Oracle's inline external table source.
//
//	inline_external_table:
//	    EXTERNAL ( ( column_definition [, column_definition ]... )
//	      inline_external_table_properties ) [ t_alias ]
//
//	inline_external_table_properties:
//	    TYPE external_table_type
//	    DEFAULT DIRECTORY directory
//	    [ ACCESS PARAMETERS { ( opaque_format_spec ) | USING CLOB subquery } ]
//	    LOCATION ( 'file_uri_list' )
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
//
// Oracle 23ai reads the properties in that order but each of them as
// optional: EXTERNAL ((a NUMBER)) and EXTERNAL ((a NUMBER) LOCATION ('f'))
// parse. It takes ACCESS PARAMETERS only in parentheses (ORA-00922 at
// USING CLOB), LOCATION as a list of strings with no directory: prefix
// (ORA-00931 at a name), and any expression as REJECT LIMIT ('x', 1+1, :b).
// A column definition is a column, a built-in datatype, and NOT NULL alone
// (selectExternalColumn).
func (p *Parser) parseInlineExternalTable(start int) (nodes.TableExpr, error) {
	p.advance() // consume EXTERNAL
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume outer '('

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume column-definition '('

	ref := &nodes.InlineExternalTable{
		Columns: &nodes.List{},
		Loc:     nodes.Loc{Start: start},
	}
	for {
		col, err := p.parseSelectExternalColumn()
		if err != nil {
			return nil, err
		}
		ref.Columns.Items = append(ref.Columns.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume column-definition ')'

	var err error
	if p.cur.Type == kwTYPE {
		p.advance()
		if ref.Type, err = p.parseSelectExternalName(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if p.cur.Type != kwDIRECTORY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if ref.Directory, err = p.parseSelectExternalName(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwACCESS {
		p.advance()
		if !p.isKeywordStr("PARAMETERS") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if ref.AccessParams, err = p.collectBalancedRaw(); err != nil {
			return nil, err
		}
	}
	if p.isKeywordStr("LOCATION") {
		p.advance()
		if ref.Location, err = p.parseSelectExternalLocation(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwREJECT {
		p.advance()
		if p.cur.Type != kwLIMIT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if ref.RejectLimit, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if ref.RejectLimit == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume outer ')'

	if p.cur.Type == kwAS {
		p.advance()
		if ref.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	} else if p.isTableAliasCandidate() {
		if ref.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	}

	ref.Loc.End = p.prev.End
	return p.parsePivotSuffix(ref)
}

// parseSelectExternalName parses the name after TYPE or DEFAULT DIRECTORY.
func (p *Parser) parseSelectExternalName() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return "", err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// parseSelectExternalColumn parses a column definition of an inline
// external table. Oracle 23ai takes a column, a built-in datatype, and NOT
// NULL alone there: any other column property is ORA-00907 where it starts
// (DEFAULT, ENCRYPT, an identity clause, a constraint, NULL, VISIBLE, SORT,
// COLLATE, NOT NULL ENABLE), and a datatype no external table holds is
// ORA-00902 where the type starts: a user-defined type, an XMLTYPE, a REF, a
// BFILE, a LONG, a domain (selectExternalDatatype).
func (p *Parser) parseSelectExternalColumn() (*nodes.ColumnDef, error) {
	start := p.pos()
	if err := p.syntaxErrorIfReservedColumnIdentifier(); err != nil {
		return nil, err
	}
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	col := &nodes.ColumnDef{
		Name:        name,
		Constraints: &nodes.List{},
		Loc:         nodes.Loc{Start: start},
	}
	if !p.isTypeName() || p.isKeywordStr("DOMAIN") {
		return nil, p.syntaxErrorAtCur()
	}
	if col.TypeName, err = p.parseTypeName(); err != nil {
		return nil, err
	}
	if !selectExternalDatatype[p.predefinedTypeSpelling(col.TypeName)] {
		return nil, p.syntaxErrorAtNode(col.TypeName)
	}
	// VARCHAR2, NVARCHAR2 and RAW take a size (ORA-00906).
	switch p.predefinedTypeSpelling(col.TypeName) {
	case "VARCHAR2", "VARCHAR", "NVARCHAR2", "RAW":
		if col.TypeName.TypeMods.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type != kwNULL {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		col.NotNull = true
	}
	col.Loc.End = p.prev.End
	return col, nil
}

// selectExternalDatatype holds the datatypes, by their first word, an
// inline external table column takes on Oracle 23ai.
var selectExternalDatatype = map[string]bool{
	"NUMBER": true, "INTEGER": true, "INT": true, "SMALLINT": true, "DECIMAL": true,
	"DEC": true, "NUMERIC": true, "FLOAT": true, "REAL": true, "DOUBLE": true,
	"BINARY_FLOAT": true, "BINARY_DOUBLE": true, "BOOLEAN": true, "BOOL": true,
	"CHAR": true, "CHARACTER": true, "VARCHAR2": true, "VARCHAR": true,
	"NCHAR": true, "NVARCHAR2": true, "CLOB": true, "NCLOB": true, "BLOB": true,
	"JSON": true, "VECTOR": true, "DATE": true, "TIMESTAMP": true, "INTERVAL": true,
	"RAW": true, "ROWID": true, "UROWID": true,
}

// parseSelectExternalLocation parses LOCATION's parenthesized list of
// strings. Oracle 23ai takes no directory: prefix and no empty list there
// (ORA-00931).
func (p *Parser) parseSelectExternalLocation() (string, error) {
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	start := p.pos()
	var files []string
	for {
		if p.cur.Type != tokSCONST && p.cur.Type != tokNCHARLIT {
			return "", p.syntaxErrorAtCur()
		}
		files = append(files, p.cur.Str)
		p.advance()
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	raw := strings.TrimSpace(p.source[start:p.cur.Loc])
	p.advance()
	if len(files) == 1 {
		return files[0], nil
	}
	return raw, nil
}

// collectBalancedRaw consumes a parenthesized group and returns the text
// inside it.
func (p *Parser) collectBalancedRaw() (string, error) {
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	start := p.pos()
	depth := 1
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
			if depth == 0 {
				raw := strings.TrimSpace(p.source[start:p.cur.Loc])
				p.advance()
				return raw, nil
			}
		}
		p.advance()
	}
	return "", p.syntaxErrorAtCur()
}

// parseJoinContinuation parses the joins that follow a table reference.
//
//	join_clause:
//	    table_reference
//	    { inner_cross_join_clause | outer_join_clause | cross_outer_apply_clause }...
//
//	inner_cross_join_clause:
//	    { [ INNER ] JOIN table_reference { ON condition | USING ( column [, column ]... ) }
//	    | CROSS JOIN table_reference
//	    | NATURAL [ INNER ] JOIN table_reference
//	    }
//
//	outer_join_clause:
//	    { outer_join_type JOIN table_reference
//	          { ON condition | USING ( column [, column ]... ) }
//	    | NATURAL outer_join_type JOIN table_reference
//	    }
//
//	outer_join_type:
//	    { LEFT | RIGHT | FULL } [ OUTER ]
//
//	cross_outer_apply_clause:
//	    { CROSS | OUTER } APPLY { table_reference | collection_expression }
//
// A join that takes a condition takes exactly one, ON or USING; a natural
// join, a cross join and an apply take none. When the table after JOIN is
// itself followed by a join rather than the condition, Oracle 23ai nests
// that join on the right and reads the condition after it: a JOIN b CROSS
// JOIN c ON x joins a to (b CROSS JOIN c), and a JOIN b JOIN c ON x lacks
// the outer join's condition (ORA-02000 at the end).
func (p *Parser) parseJoinContinuation(left nodes.TableExpr) (nodes.TableExpr, error) {
	for p.isSelectJoinStart() || p.isSelectQueryPartitionStart() {
		joinStart := p.pos()
		// A partitioned outer join partitions its left table:
		// t PARTITION BY (...) LEFT JOIN u ON ... (an outer join follows).
		var leftPartition *nodes.List
		if p.isSelectQueryPartitionStart() {
			var err error
			if leftPartition, err = p.parseSelectQueryPartition(true); err != nil {
				return nil, err
			}
			switch p.cur.Type {
			case kwLEFT, kwRIGHT, kwFULL:
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}
		jt, err := p.matchJoinType()
		if err != nil {
			return nil, err
		}

		right, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if right == nil {
			return nil, p.syntaxErrorAtCur()
		}

		// A join starts where its left table does.
		jc := &nodes.JoinClause{
			Type:            jt,
			Left:            left,
			Right:           right,
			LeftPartitionBy: leftPartition,
			Loc:             nodes.Loc{Start: exprLocStart(left, joinStart)},
		}
		// ... or its right one: LEFT JOIN u PARTITION BY (...) ON ...
		if selectOuterJoin(jt) && p.isSelectQueryPartitionStart() {
			if jc.RightPartitionBy, err = p.parseSelectQueryPartition(false); err != nil {
				return nil, err
			}
		}

		if selectJoinTakesCondition(jt) {
			if p.isSelectJoinStart() {
				if jc.Right, err = p.parseSelectNestedJoin(right); err != nil {
					return nil, err
				}
			}
			switch p.cur.Type {
			case kwON:
				p.advance()
				if jc.On, err = p.parseExpr(); err != nil {
					return nil, err
				}
				if jc.On == nil {
					return nil, p.syntaxErrorAtCur()
				}
			case kwUSING:
				if jc.Using, err = p.parseSelectJoinUsing(); err != nil {
					return nil, err
				}
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}

		jc.Loc.End = p.prev.End
		left = jc
	}
	return left, nil
}

// parseSelectNestedJoin parses the joins nested on the right of a join that
// still owes its condition.
func (p *Parser) parseSelectNestedJoin(right nodes.TableExpr) (nodes.TableExpr, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	return p.parseJoinContinuation(right)
}

// parseSelectJoinUsing parses USING ( column [, column ]... ).
func (p *Parser) parseSelectJoinUsing() (*nodes.List, error) {
	p.advance() // consume USING
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	using := &nodes.List{}
	for {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		using.Items = append(using.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return using, nil
}

// isSelectJoinStart reports whether the current token starts a join after
// a table reference.
func (p *Parser) isSelectJoinStart() bool {
	switch p.cur.Type {
	case kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwNATURAL, kwOUTER:
		return true
	}
	return false
}

// selectOuterJoin reports whether a join of type jt is an outer join, which
// a query_partition_clause may partition.
func selectOuterJoin(jt nodes.JoinType) bool {
	switch jt {
	case nodes.JOIN_LEFT, nodes.JOIN_RIGHT, nodes.JOIN_FULL,
		nodes.JOIN_NATURAL_LEFT, nodes.JOIN_NATURAL_RIGHT, nodes.JOIN_NATURAL_FULL:
		return true
	}
	return false
}

// selectJoinTakesCondition reports whether a join of type jt takes ON or
// USING.
func selectJoinTakesCondition(jt nodes.JoinType) bool {
	switch jt {
	case nodes.JOIN_INNER, nodes.JOIN_LEFT, nodes.JOIN_RIGHT, nodes.JOIN_FULL:
		return true
	}
	return false
}

// matchJoinType reads the words of a join up to and including JOIN or
// APPLY, at a token isSelectJoinStart accepts. Oracle 23ai fails at the
// first word out of place: NATURAL a (ORA-03049), LEFT a, INNER a, CROSS a,
// OUTER JOIN (ORA-02000).
func (p *Parser) matchJoinType() (nodes.JoinType, error) {
	natural := p.cur.Type == kwNATURAL
	if natural {
		p.advance()
	}
	expectJoin := func() error {
		if p.cur.Type != kwJOIN {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	}

	var jt nodes.JoinType
	switch p.cur.Type {
	case kwJOIN, kwINNER:
		if p.cur.Type == kwINNER {
			p.advance()
		}
		if err := expectJoin(); err != nil {
			return 0, err
		}
		jt = nodes.JOIN_INNER
		if natural {
			jt = nodes.JOIN_NATURAL_INNER
		}
	case kwLEFT, kwRIGHT, kwFULL:
		side := p.cur.Type
		p.advance()
		if p.cur.Type == kwOUTER {
			p.advance()
		}
		if err := expectJoin(); err != nil {
			return 0, err
		}
		switch {
		case side == kwLEFT && natural:
			jt = nodes.JOIN_NATURAL_LEFT
		case side == kwLEFT:
			jt = nodes.JOIN_LEFT
		case side == kwRIGHT && natural:
			jt = nodes.JOIN_NATURAL_RIGHT
		case side == kwRIGHT:
			jt = nodes.JOIN_RIGHT
		case natural:
			jt = nodes.JOIN_NATURAL_FULL
		default:
			jt = nodes.JOIN_FULL
		}
	case kwCROSS, kwOUTER:
		if natural {
			return 0, p.syntaxErrorAtCur()
		}
		cross := p.cur.Type == kwCROSS
		p.advance()
		switch {
		case p.isKeywordStr("APPLY"):
			p.advance()
			if cross {
				return nodes.JOIN_CROSS_APPLY, nil
			}
			return nodes.JOIN_OUTER_APPLY, nil
		case cross && p.cur.Type == kwJOIN:
			p.advance()
			return nodes.JOIN_CROSS, nil
		}
		return 0, p.syntaxErrorAtCur()
	default:
		return 0, p.syntaxErrorAtCur()
	}
	return jt, nil
}

// parseFlashbackClause parses a flashback query clause on a table reference.
//
//	flashback_query_clause:
//	    { VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
//	    | VERSIONS PERIOD FOR valid_time_column BETWEEN expr AND expr
//	    | AS OF { SCN | TIMESTAMP } expr
//	    | AS OF PERIOD FOR valid_time_column expr
//	    }
//
// The valid-time column is a name, neither reserved (ORA-03050 for
// BETWEEN) nor a number (ORA-03054).
func (p *Parser) parseFlashbackClause() (*nodes.FlashbackClause, error) {
	start := p.pos()
	fc := &nodes.FlashbackClause{Loc: nodes.Loc{Start: start}}

	periodFor := func() error {
		if !p.isKeywordStr("PERIOD") {
			return nil
		}
		fc.IsPeriodFor = true
		p.advance()
		if p.cur.Type != kwFOR {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
			return err
		}
		var err error
		fc.PeriodColumn, err = p.parseIdentifier()
		return err
	}
	scnOrTimestamp := func() error {
		switch p.cur.Type {
		case kwSCN:
			fc.Type = "SCN"
		case kwTIMESTAMP:
			fc.Type = "TIMESTAMP"
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	}
	expr := func(prec int) (nodes.ExprNode, error) {
		// A bound is missing before AND (ORA-00936), which no expression
		// starts.
		if p.cur.Type == kwAND {
			return nil, p.syntaxErrorAtCur()
		}
		e, err := p.parseExprPrec(prec)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return e, nil
	}

	var err error
	if p.cur.Type == kwVERSIONS {
		fc.IsVersions = true
		p.advance()
		if err := periodFor(); err != nil {
			return nil, err
		}
		if p.cur.Type != kwBETWEEN {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if !fc.IsPeriodFor {
			if err := scnOrTimestamp(); err != nil {
				return nil, err
			}
		}
		// The low bound binds tighter than AND, which ends it.
		if fc.VersionsLow, err = expr(precNot); err != nil {
			return nil, err
		}
		if p.cur.Type != kwAND {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if fc.VersionsHigh, err = expr(precOr); err != nil {
			return nil, err
		}
	} else {
		p.advance() // consume AS
		if p.cur.Type != kwOF {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if err := periodFor(); err != nil {
			return nil, err
		}
		if !fc.IsPeriodFor {
			if err := scnOrTimestamp(); err != nil {
				return nil, err
			}
		}
		if fc.Expr, err = expr(precOr); err != nil {
			return nil, err
		}
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseSampleClause parses a SAMPLE clause on a table reference.
//
//	sample_clause:
//	    SAMPLE [ BLOCK ] ( sample_percent ) [ SEED ( seed_value ) ]
//
// Oracle 23ai takes an unsigned numeric literal for each value, not an
// expression, a bind or a name (ORA-03048 at the '(': SAMPLE (1+1),
// SAMPLE (:b), SAMPLE ()), and a second literal after the percent,
// SAMPLE (1, 2).
func (p *Parser) parseSampleClause() (*nodes.SampleClause, error) {
	start := p.pos()
	p.advance() // consume SAMPLE

	sc := &nodes.SampleClause{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwBLOCK {
		sc.Block = true
		p.advance()
	}
	var err error
	if sc.Percent, err = p.parseSelectSampleValue(true); err != nil {
		return nil, err
	}
	if p.cur.Type == kwSEED {
		p.advance()
		if sc.Seed, err = p.parseSelectSampleValue(false); err != nil {
			return nil, err
		}
	}
	sc.Loc.End = p.prev.End
	return sc, nil
}

// isSelectSampleClauseAhead reports whether the SAMPLE at the current token
// opens a sample_clause: SAMPLE ( or SAMPLE BLOCK.
func (p *Parser) isSelectSampleClauseAhead() bool {
	next := p.peekNext().Type
	return next == '(' || next == kwBLOCK
}

// parseSelectSampleValue parses ( number ), or ( number [, number ] ) when
// pair, failing at the '(' when what it holds is not that.
func (p *Parser) parseSelectSampleValue(pair bool) (nodes.ExprNode, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	open := p.cur
	p.advance()
	if !p.isSelectSampleNumber() {
		return nil, p.syntaxErrorAtTok(open)
	}
	value, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	if pair && p.cur.Type == ',' {
		p.advance()
		if !p.isSelectSampleNumber() {
			return nil, p.syntaxErrorAtTok(open)
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtTok(open)
	}
	p.advance()
	return value, nil
}

// isSelectSampleNumber reports whether the current token is an unsigned
// NUMBER literal, not a BINARY_FLOAT or BINARY_DOUBLE one (1d, 1f).
func (p *Parser) isSelectSampleNumber() bool {
	return p.cur.Type == tokICONST || p.cur.Type == tokFCONST && !strings.ContainsAny(p.cur.Str, "fFdD")
}

// parsePivotClause parses a PIVOT clause.
//
//	pivot_clause:
//	    PIVOT [ XML ]
//	    ( aggregate_function ( expr ) [ [ AS ] c_alias ]
//	      [, aggregate_function ( expr ) [ [ AS ] c_alias ] ]...
//	      pivot_for_clause
//	      pivot_in_clause
//	    )
//
//	pivot_for_clause:
//	    FOR { column | ( column [, column ]... ) }
//
//	pivot_in_clause:
//	    IN ( { { expr | ( expr [, expr ]... ) } [ [ AS ] c_alias ]
//	         [, { expr | ( expr [, expr ]... ) } [ [ AS ] c_alias ] ]...
//	         | subquery
//	         | ANY [, ANY ]...
//	         } )
//
// The subquery and ANY forms are PIVOT XML's, and its only ones (ORA-02000
// at a value); without XML they are ORA-00936. A FOR column is a simple name
// (ORA-01748 at a '.'). With one FOR column an IN value is an expression;
// with n of them it is a parenthesized list of n (ORA-00906, ORA-00907), and
// PIVOT XML takes n ANYs. Oracle 23ai parses an empty IN list.
func (p *Parser) parsePivotClause() (*nodes.PivotClause, error) {
	start := p.pos()
	p.advance() // consume PIVOT

	pc := &nodes.PivotClause{Loc: nodes.Loc{Start: start}}
	if p.isKeywordStr("XML") {
		pc.XML = true
		p.advance()
	}
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	pc.AggFuncs = &nodes.List{}
	pc.Aggregates = &nodes.List{}
	for {
		agg, err := p.parseResTarget()
		if err != nil {
			return nil, err
		}
		if agg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		pc.AggFuncs.Items = append(pc.AggFuncs.Items, agg)
		pc.Aggregates.Items = append(pc.Aggregates.Items, &nodes.PivotAggregate{
			Expr:  agg.Expr,
			Alias: agg.Name,
			Loc:   agg.Loc,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type != kwFOR {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	parenthesized := p.cur.Type == '('
	cols, err := p.parseSelectPivotColumns()
	if err != nil {
		return nil, err
	}
	pc.ForColumns = cols
	if parenthesized && cols.Len() > 1 {
		pc.ForCols = &nodes.List{Items: cols.Items}
	} else {
		pc.ForCol = cols.Items[0].(nodes.ExprNode)
	}

	if p.cur.Type != kwIN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case !pc.XML:
		if pc.InList, pc.InItems, err = p.parsePivotInList(cols.Len()); err != nil {
			return nil, err
		}
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH:
		sub, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		pc.InList = &nodes.List{Items: []nodes.Node{sub}}
	case p.cur.Type == kwANY:
		pc.InList = &nodes.List{}
		for i := 0; i < cols.Len(); i++ {
			if i > 0 {
				if p.cur.Type != ',' {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
			}
			if p.cur.Type != kwANY {
				return nil, p.syntaxErrorAtCur()
			}
			pc.InList.Items = append(pc.InList.Items, &nodes.ColumnRef{
				Column: "ANY",
				Loc:    nodes.Loc{Start: p.pos(), End: p.cur.End},
			})
			p.advance()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	pc.Loc.End = p.prev.End
	return pc, nil
}

// parseSelectPivotColumns parses column | ( column [, column ]... ), each a
// simple name, for PIVOT's and UNPIVOT's FOR and UNPIVOT's value columns.
func (p *Parser) parseSelectPivotColumns() (*nodes.List, error) {
	cols := &nodes.List{}
	if p.cur.Type != '(' {
		col, err := p.parseSelectMatchColumn()
		if err != nil {
			return nil, err
		}
		cols.Items = append(cols.Items, col)
		return cols, nil
	}
	p.advance()
	for {
		col, err := p.parseSelectMatchColumn()
		if err != nil {
			return nil, err
		}
		cols.Items = append(cols.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return cols, nil
}

// parseSelectExprTuple parses ( expr [, expr ]... ) of exactly n
// expressions, failing at the token where the count goes wrong.
func (p *Parser) parseSelectExprTuple(n int) (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for i := 0; i < n; i++ {
		if i > 0 {
			if p.cur.Type != ',' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, p.syntaxErrorAtCur()
		}
		list.Items = append(list.Items, e)
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parsePivotInList parses the values of a PIVOT IN list for n FOR columns,
// up to its closing parenthesis.
func (p *Parser) parsePivotInList(n int) (*nodes.List, *nodes.List, error) {
	list := &nodes.List{}
	items := &nodes.List{}
	if p.cur.Type == ')' {
		return list, items, nil
	}
	for {
		start := p.pos()
		var expr nodes.ExprNode
		values := &nodes.List{}
		if n > 1 {
			tuple, err := p.parseSelectExprTuple(n)
			if err != nil {
				return nil, nil, err
			}
			// A tuple is kept as a call with no name.
			expr = &nodes.FuncCallExpr{
				FuncName: &nodes.ObjectName{Loc: nodes.Loc{Start: start, End: p.prev.End}},
				Args:     tuple,
				Loc:      nodes.Loc{Start: start, End: p.prev.End},
			}
			values.Items = append(values.Items, tuple.Items...)
		} else {
			// ANY, a subquery and * are PIVOT XML's (ORA-00936).
			if p.isSelectPivotXMLValueStart() {
				return nil, nil, p.syntaxErrorAtCur()
			}
			p.scalarParen = p.cur.Type == '('
			var err error
			expr, err = p.parseExpr()
			p.scalarParen = false
			if err != nil {
				return nil, nil, err
			}
			if expr == nil {
				return nil, nil, p.syntaxErrorAtCur()
			}
			if paren, ok := expr.(*nodes.ParenExpr); ok {
				values.Items = append(values.Items, paren.Expr)
			} else {
				values.Items = append(values.Items, expr)
			}
		}

		rt := &nodes.ResTarget{Expr: expr, Loc: nodes.Loc{Start: start}}
		if p.cur.Type == kwAS {
			p.advance()
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, nil, err
			}
			if name == "" {
				return nil, nil, p.syntaxErrorAtCur()
			}
			rt.Name = name
		} else if p.isAliasCandidate() {
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, nil, err
			}
			rt.Name = name
		}
		rt.Loc.End = p.prev.End
		list.Items = append(list.Items, rt)
		items.Items = append(items.Items, &nodes.PivotInItem{Values: values, Alias: rt.Name, Loc: rt.Loc})

		if p.cur.Type != ',' {
			return list, items, nil
		}
		p.advance()
	}
}

// isSelectPivotXMLValueStart reports whether the current token starts what
// only a PIVOT XML IN list holds, ANY or a subquery, or the * no IN list
// holds.
func (p *Parser) isSelectPivotXMLValueStart() bool {
	switch p.cur.Type {
	case kwANY, kwSELECT, kwWITH, '*':
		return true
	}
	return false
}

// parseUnpivotClause parses an UNPIVOT clause.
//
//	unpivot_clause:
//	    UNPIVOT [ { INCLUDE | EXCLUDE } NULLS ]
//	    ( { column | ( column [, column ]... ) }
//	      pivot_for_clause
//	      unpivot_in_clause
//	    )
//
//	unpivot_in_clause:
//	    IN ( { column | ( column [, column ]... ) }
//	         [ AS { literal | ( literal [, literal ]... ) } ]
//	         [, { column | ( column [, column ]... ) }
//	            [ AS { literal | ( literal [, literal ]... ) } ] ]...
//	       )
//
// Every column is a simple name (ORA-02000, ORA-01738, ORA-01748 at a '.').
// With n value columns an IN item is a parenthesized list of n columns, and
// with one it is a column, parenthesized or not; with m FOR columns a label
// is a parenthesized list of m values, and with one it is a value. Oracle
// 23ai parses any expression as a label value and refuses one that is not
// constant after the parse (ORA-56901).
func (p *Parser) parseUnpivotClause() (*nodes.UnpivotClause, error) {
	start := p.pos()
	p.advance() // consume UNPIVOT

	uc := &nodes.UnpivotClause{Loc: nodes.Loc{Start: start}}
	switch {
	case p.cur.Type == kwINCLUDE, p.isKeywordStr("EXCLUDE"):
		uc.IncludeNulls = p.cur.Type == kwINCLUDE
		p.advance()
		if p.cur.Type != kwNULLS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	valueCols, err := p.parseSelectPivotColumns()
	if err != nil {
		return nil, err
	}
	uc.ValueColumns = valueCols
	uc.ValueCol = valueCols.Items[0].(nodes.ExprNode)

	if p.cur.Type != kwFOR {
		return nil, p.syntaxErrorAtCur()
	}
	forStart := p.pos()
	p.advance()
	parenthesized := p.cur.Type == '('
	forCols, err := p.parseSelectPivotColumns()
	if err != nil {
		return nil, err
	}
	if parenthesized && forCols.Len() > 1 {
		// Several label columns are kept as a call with no name.
		uc.PivotColumn = &nodes.FuncCallExpr{
			FuncName: &nodes.ObjectName{Loc: nodes.Loc{Start: forStart, End: p.prev.End}},
			Args:     forCols,
			Loc:      nodes.Loc{Start: forStart, End: p.prev.End},
		}
	} else {
		uc.PivotColumn = forCols.Items[0].(nodes.ExprNode)
	}
	uc.PivotCol = uc.PivotColumn

	if p.cur.Type != kwIN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	uc.InList = &nodes.List{}
	uc.InputMappings = &nodes.List{}
	for {
		itemStart := p.pos()
		inputCols := &nodes.List{}
		if valueCols.Len() > 1 || p.cur.Type == '(' {
			if p.cur.Type != '(' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			for i := 0; i < valueCols.Len(); i++ {
				if i > 0 {
					if p.cur.Type != ',' {
						return nil, p.syntaxErrorAtCur()
					}
					p.advance()
				}
				col, err := p.parseSelectMatchColumn()
				if err != nil {
					return nil, err
				}
				inputCols.Items = append(inputCols.Items, col)
			}
			if p.cur.Type != ')' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		} else {
			col, err := p.parseSelectMatchColumn()
			if err != nil {
				return nil, err
			}
			inputCols.Items = append(inputCols.Items, col)
		}
		rt := &nodes.ResTarget{Expr: inputCols.Items[0].(nodes.ExprNode), Loc: nodes.Loc{Start: itemStart}}
		item := &nodes.UnpivotInItem{InputColumns: inputCols, Loc: nodes.Loc{Start: itemStart}}
		if p.cur.Type == kwAS {
			p.advance()
			labelStart := p.pos()
			if forCols.Len() > 1 {
				tuple, err := p.parseSelectExprTuple(forCols.Len())
				if err != nil {
					return nil, err
				}
				item.Label = &nodes.FuncCallExpr{
					FuncName: &nodes.ObjectName{Loc: nodes.Loc{Start: labelStart, End: p.prev.End}},
					Args:     tuple,
					Loc:      nodes.Loc{Start: labelStart, End: p.prev.End},
				}
			} else {
				p.scalarParen = p.cur.Type == '('
				label, err := p.parseExpr()
				p.scalarParen = false
				if err != nil {
					return nil, err
				}
				if label == nil {
					return nil, p.syntaxErrorAtCur()
				}
				item.Label = label
			}
			switch label := item.Label.(type) {
			case *nodes.StringLiteral:
				item.Alias = label.Val
			case *nodes.ColumnRef:
				item.Alias = label.Column
			default:
				item.Alias = strings.TrimSpace(p.source[labelStart:p.prev.End])
			}
			rt.Name = item.Alias
		}
		rt.Loc.End = p.prev.End
		item.Loc.End = p.prev.End
		uc.InList.Items = append(uc.InList.Items, rt)
		uc.InputMappings.Items = append(uc.InputMappings.Items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	uc.Loc.End = p.prev.End
	return uc, nil
}

// parseMatchRecognize parses a MATCH_RECOGNIZE clause.
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
//
//	row_pattern_measure_column:
//	    expr AS c_alias
//
//	subset_item:
//	    variable_name = ( variable_name [, variable_name ]... )
//
//	row_pattern_definition:
//	    variable_name AS condition
//
// Oracle 23ai reads the clauses in that order. A partitioning or ordering
// column is a bare column name (ORA-00907 at a '.', an operator or a
// number), and an ordering column takes NULLS FIRST or LAST too. A measure
// takes its alias with or without AS, and ALL ROWS PER MATCH takes WITH
// UNMATCHED ROWS too, as the Data Warehousing Guide documents. A pattern
// variable is a name that is not reserved (ORA-03050 for ROW).
func (p *Parser) parseMatchRecognize(start int) (nodes.TableExpr, error) {
	p.advance() // consume MATCH_RECOGNIZE
	mr := &nodes.MatchRecognizeClause{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		mr.Loc.End = p.prev.End
		return mr, nil
	}
	p.advance() // consume '('

	var err error
	if p.cur.Type == kwPARTITION {
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		mr.PartitionBy = &nodes.List{}
		for {
			col, err := p.parseSelectMatchColumn()
			if err != nil {
				return nil, err
			}
			mr.PartitionBy.Items = append(mr.PartitionBy.Items, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == kwORDER {
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		mr.OrderBy = &nodes.List{}
		for {
			sb, err := p.parseSelectMatchSortBy()
			if err != nil {
				return nil, err
			}
			mr.OrderBy.Items = append(mr.OrderBy.Items, sb)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == kwMEASURES {
		p.advance()
		mr.Measures = &nodes.List{}
		for {
			rt, err := p.parseSelectMatchMeasure()
			if err != nil {
				return nil, err
			}
			mr.Measures.Items = append(mr.Measures.Items, rt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if mr.RowsPerMatch, err = p.parseSelectMatchRowsPerMatch(); err != nil {
		return nil, err
	}
	if mr.AfterMatch, err = p.parseSelectMatchAfterMatch(); err != nil {
		return nil, err
	}

	if !p.isKeywordStr("PATTERN") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parts []string
	if err := p.parseSelectRowPattern(&parts); err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	mr.Pattern = strings.Join(parts, " ")

	if p.isKeywordStr("SUBSET") {
		p.advance()
		mr.Subsets = &nodes.List{}
		for {
			rt, err := p.parseSelectMatchSubsetItem()
			if err != nil {
				return nil, err
			}
			mr.Subsets.Items = append(mr.Subsets.Items, rt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if !p.isKeywordStr("DEFINE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	mr.Definitions = &nodes.List{}
	for {
		defStart := p.pos()
		name, err := p.parseSelectPatternVariable()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != kwAS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		rt := &nodes.ResTarget{Name: name, Loc: nodes.Loc{Start: defStart}}
		if rt.Expr, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if rt.Expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		rt.Loc.End = p.prev.End
		mr.Definitions.Items = append(mr.Definitions.Items, rt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if p.cur.Type == kwAS {
		p.advance()
		if mr.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	} else if p.isTableAliasCandidate() {
		if mr.Alias, err = p.parseAlias(); err != nil {
			return nil, err
		}
	}

	mr.Loc.End = p.prev.End
	return mr, nil
}

// parseSelectMatchColumn parses a MATCH_RECOGNIZE partitioning or ordering
// column: a name, not qualified.
func (p *Parser) parseSelectMatchColumn() (*nodes.ColumnRef, error) {
	start := p.pos()
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &nodes.ColumnRef{Column: name, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
}

// parseSelectMatchSortBy parses column [ ASC | DESC ] [ NULLS { FIRST | LAST } ].
func (p *Parser) parseSelectMatchSortBy() (*nodes.SortBy, error) {
	start := p.pos()
	col, err := p.parseSelectMatchColumn()
	if err != nil {
		return nil, err
	}
	sb := &nodes.SortBy{Expr: col, Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwASC:
		sb.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		sb.Dir = nodes.SORTBY_DESC
		p.advance()
	}
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			sb.NullOrder = nodes.SORTBY_NULLS_FIRST
		case kwLAST:
			sb.NullOrder = nodes.SORTBY_NULLS_LAST
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}
	sb.Loc.End = p.prev.End
	return sb, nil
}

// parseSelectMatchMeasure parses expr [ AS ] c_alias. The alias is
// required: Oracle 23ai reads MEASURES 1 PATTERN (a) as a measure named
// PATTERN and fails at the '('.
func (p *Parser) parseSelectMatchMeasure() (*nodes.ResTarget, error) {
	start := p.pos()
	// A measure's expression comes before AS (ORA-00936).
	if p.cur.Type == kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwAS {
		p.advance()
	}
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	rt := &nodes.ResTarget{Expr: expr, Loc: nodes.Loc{Start: start}}
	if rt.Name, err = p.parseIdentifier(); err != nil {
		return nil, err
	}
	rt.Loc.End = p.prev.End
	return rt, nil
}

// parseSelectMatchRowsPerMatch parses
//
//	[ ONE ROW PER MATCH
//	| ALL ROWS PER MATCH [ { SHOW | OMIT } EMPTY MATCHES | WITH UNMATCHED ROWS ] ]
func (p *Parser) parseSelectMatchRowsPerMatch() (string, error) {
	var words []string
	switch {
	case p.isKeywordStr("ONE"):
		words = []string{"ONE", "ROW", "PER", "MATCH"}
	case p.cur.Type == kwALL:
		words = []string{"ALL", "ROWS", "PER", "MATCH"}
	default:
		return "", nil
	}
	if err := p.expectSelectWords(words...); err != nil {
		return "", err
	}
	if words[0] == "ONE" {
		return "ONE ROW PER MATCH", nil
	}
	switch {
	case p.isKeywordStr("SHOW"), p.isKeywordStr("OMIT"):
		mode := p.cur.Str
		p.advance()
		if err := p.expectSelectWords("EMPTY", "MATCHES"); err != nil {
			return "", err
		}
		return "ALL ROWS PER MATCH " + mode + " EMPTY MATCHES", nil
	case p.cur.Type == kwWITH:
		p.advance()
		if err := p.expectSelectWords("UNMATCHED", "ROWS"); err != nil {
			return "", err
		}
		return "ALL ROWS PER MATCH WITH UNMATCHED ROWS", nil
	}
	return "ALL ROWS PER MATCH", nil
}

// parseSelectMatchAfterMatch parses
//
//	[ AFTER MATCH SKIP { PAST LAST ROW | TO NEXT ROW | TO FIRST variable_name
//	                   | TO LAST variable_name | TO variable_name } ]
func (p *Parser) parseSelectMatchAfterMatch() (string, error) {
	if !p.isKeywordStr("AFTER") {
		return "", nil
	}
	if err := p.expectSelectWords("AFTER", "MATCH", "SKIP"); err != nil {
		return "", err
	}
	if p.isKeywordStr("PAST") {
		if err := p.expectSelectWords("PAST", "LAST", "ROW"); err != nil {
			return "", err
		}
		return "PAST LAST ROW", nil
	}
	if p.cur.Type != kwTO {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	prefix := "TO "
	switch p.cur.Type {
	case kwNEXT:
		if err := p.expectSelectWords("NEXT", "ROW"); err != nil {
			return "", err
		}
		return "TO NEXT ROW", nil
	case kwFIRST, kwLAST:
		prefix += p.cur.Str + " "
		p.advance()
	}
	name, err := p.parseSelectPatternVariable()
	if err != nil {
		return "", err
	}
	return prefix + name, nil
}

// expectSelectWords consumes the unquoted words in order, failing at the
// first token that is not the word expected there.
func (p *Parser) expectSelectWords(words ...string) error {
	for _, w := range words {
		if !p.isKeywordStr(w) {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return nil
}

// parseSelectPatternVariable parses a pattern variable name: an identifier
// that is not reserved.
func (p *Parser) parseSelectPatternVariable() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return "", err
	}
	return p.parseIdentifier()
}

// parseSelectMatchSubsetItem parses variable_name = ( variable_name [, ...] ).
// A list of more than one variable is kept as a tuple, as the PIVOT IN list
// keeps one: a call with no name.
func (p *Parser) parseSelectMatchSubsetItem() (*nodes.ResTarget, error) {
	start := p.pos()
	name, err := p.parseSelectPatternVariable()
	if err != nil {
		return nil, err
	}
	rt := &nodes.ResTarget{Name: name, Loc: nodes.Loc{Start: start}}
	if p.cur.Type != '=' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	listStart := p.pos()
	p.advance()
	vars := &nodes.List{}
	for {
		varStart := p.pos()
		v, err := p.parseSelectPatternVariable()
		if err != nil {
			return nil, err
		}
		vars.Items = append(vars.Items, &nodes.ColumnRef{Column: v, Loc: nodes.Loc{Start: varStart, End: p.prev.End}})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if vars.Len() == 1 {
		rt.Expr = vars.Items[0].(*nodes.ColumnRef)
	} else {
		rt.Expr = &nodes.FuncCallExpr{
			FuncName: &nodes.ObjectName{Loc: nodes.Loc{Start: listStart, End: p.prev.End}},
			Args:     vars,
			Loc:      nodes.Loc{Start: listStart, End: p.prev.End},
		}
	}
	rt.Loc.End = p.prev.End
	return rt, nil
}

// parseSelectRowPattern parses a row pattern, appending the text of each
// token it reads to parts.
//
//	row_pattern:
//	    row_pattern_term [ '|' row_pattern_term ]...
//
//	row_pattern_term:
//	    row_pattern_factor [ row_pattern_factor ]...
//
//	row_pattern_factor:
//	    row_pattern_primary [ row_pattern_quantifier ]
func (p *Parser) parseSelectRowPattern(parts *[]string) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	for {
		if !p.isSelectRowPatternPrimaryStart() {
			return p.syntaxErrorAtCur()
		}
		for p.isSelectRowPatternPrimaryStart() {
			if err := p.parseSelectRowPatternFactor(parts); err != nil {
				return err
			}
		}
		if p.cur.Type != '|' {
			return nil
		}
		*parts = append(*parts, p.cur.Str)
		p.advance()
	}
}

// isSelectRowPatternPrimaryStart reports whether the current token starts
// a row_pattern_primary: a variable, '$', '^', '(', '{-', or PERMUTE (.
func (p *Parser) isSelectRowPatternPrimaryStart() bool {
	switch p.cur.Type {
	case '$', '^', '(':
		return true
	case '{':
		return p.peekNext().Type == '-'
	}
	return p.isIdentLike()
}

// parseSelectRowPatternFactor parses a row_pattern_primary and its
// quantifier.
//
//	row_pattern_primary:
//	    variable_name | '$' | '^' | ( [ row_pattern ] )
//	    | '{-' row_pattern '-}' | PERMUTE ( row_pattern [, row_pattern ]... )
//
// An anchor takes no quantifier (ORA-00931 at the '*' of ^*).
func (p *Parser) parseSelectRowPatternFactor(parts *[]string) error {
	take := func() {
		*parts = append(*parts, p.cur.Str)
		p.advance()
	}
	switch {
	case p.cur.Type == '$' || p.cur.Type == '^':
		take()
		return nil
	case p.cur.Type == '(':
		take()
		if p.cur.Type != ')' {
			if err := p.parseSelectRowPattern(parts); err != nil {
				return err
			}
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		take()
	case p.cur.Type == '{':
		take() // {
		take() // -
		if err := p.parseSelectRowPattern(parts); err != nil {
			return err
		}
		if p.cur.Type != '-' || p.peekNext().Type != '}' {
			return p.syntaxErrorAtCur()
		}
		take()
		take()
	case p.isKeywordStr("PERMUTE"):
		take()
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		take()
		for {
			if err := p.parseSelectRowPattern(parts); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			take()
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		take()
	default:
		if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
			return err
		}
		take()
	}
	return p.parseSelectRowPatternQuantifier(parts)
}

// parseSelectRowPatternQuantifier parses an optional row_pattern_quantifier.
//
//	row_pattern_quantifier:
//	    { '*' | '+' | '?' } [ '?' ]
//	    | '{' integer [ ',' [ integer ] ] '}' [ '?' ]
//	    | '{' ',' integer '}' [ '?' ]
//
// Oracle 23ai also takes {,} (any number of times).
func (p *Parser) parseSelectRowPatternQuantifier(parts *[]string) error {
	take := func() {
		*parts = append(*parts, p.cur.Str)
		p.advance()
	}
	switch p.cur.Type {
	case '*', '+', '?':
		take()
	case '{':
		if p.peekNext().Type == '-' {
			return nil // the next factor's exclusion
		}
		take()
		if p.isSelectPatternBound() {
			take()
			if p.cur.Type == ',' {
				take()
				if p.isSelectPatternBound() {
					take()
				}
			}
		} else if p.cur.Type == ',' {
			take()
			if p.isSelectPatternBound() {
				take()
			}
		} else {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != '}' {
			return p.syntaxErrorAtCur()
		}
		take()
	default:
		return nil
	}
	if p.cur.Type == '?' {
		take()
	}
	return nil
}

// isSelectPatternBound reports whether the current token is a quantifier
// bound: an unsigned integer, which Oracle 23ai also takes in exponent form
// (a{1e1}) but not with a fraction (ORA-00931 for a{1.5}).
func (p *Parser) isSelectPatternBound() bool {
	return p.cur.Type == tokICONST || p.cur.Type == tokFCONST && !strings.ContainsAny(p.cur.Str, ".fFdD")
}
