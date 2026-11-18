package parser

import "github.com/bytebase/omni/trino/ast"

// This file implements the JSON_TABLE relation primary and the plan and column
// grammar it carries (SqlBase.g4 at 482):
//
//	relationPrimary
//	    : ...
//	    | JSON_TABLE '('
//	        jsonPathInvocation
//	        COLUMNS '(' jsonTableColumn (',' jsonTableColumn)* ')'
//	        (PLAN '(' jsonTableSpecificPlan ')'
//	        | PLAN DEFAULT '(' jsonTableDefaultPlan ')'
//	        )?
//	        ((ERROR | EMPTY) ON ERROR)?
//	      ')'                                                             #jsonTable
//	    ;
//
//	jsonTableColumn
//	    : identifier FOR ORDINALITY                                     #ordinalityColumn
//	    | identifier type
//	        (PATH string)?
//	        (emptyBehavior=jsonValueBehavior ON EMPTY)?
//	        (errorBehavior=jsonValueBehavior ON ERROR)?                 #valueColumn
//	    | identifier type FORMAT jsonRepresentation
//	        (PATH string)?
//	        (jsonQueryWrapperBehavior WRAPPER)?
//	        ((KEEP | OMIT) QUOTES (ON SCALAR TEXT_STRING)?)?
//	        (emptyBehavior=jsonQueryBehavior ON EMPTY)?
//	        (errorBehavior=jsonQueryBehavior ON ERROR)?                 #queryColumn
//	    | NESTED PATH? string (AS identifier)?
//	        COLUMNS '(' jsonTableColumn (',' jsonTableColumn)* ')'      #nestedColumns
//	    ;
//
//	jsonTableSpecificPlan
//	    : jsonTablePathName                                         #leafPlan
//	    | jsonTablePathName (OUTER | INNER) planPrimary             #joinPlan
//	    | planPrimary UNION planPrimary (UNION planPrimary)*        #unionPlan
//	    | planPrimary CROSS planPrimary (CROSS planPrimary)*        #crossPlan
//	    ;
//	jsonTablePathName : identifier ;
//	planPrimary : jsonTablePathName | '(' jsonTableSpecificPlan ')' ;
//	jsonTableDefaultPlan
//	    : (OUTER | INNER) (',' (UNION | CROSS))?
//	    | (UNION | CROSS) (',' (OUTER | INNER))?
//	    ;
//
// The shared jsonPathInvocation (with its optional AS pathName) lives in
// json.go, and the value/query behavior, wrapper, format and QUOTES clauses
// reuse that file's helpers. JSON_TABLE's column clauses are positional, like
// the JSON_* functions' (J1 there): an ON ERROR clause cannot precede ON EMPTY,
// and the query column's FORMAT/PATH/WRAPPER/QUOTES order is fixed.

// JSONTableRelation is the `JSON_TABLE ( … )` relation primary (the jsonTable
// alternative). Path is its jsonPathInvocation; Columns is the COLUMNS list.
// Exactly one of Plan / DefaultPlan is non-nil when a PLAN clause is present.
// OnError is "", "ERROR", or "EMPTY" for the optional trailing ON ERROR
// behavior.
type JSONTableRelation struct {
	Path        *JSONPathInvocation
	Columns     []JSONTableColumn
	Plan        *JSONTableSpecificPlan
	DefaultPlan *JSONTableDefaultPlan
	OnError     string // "", "ERROR", or "EMPTY"
	Loc         ast.Loc
}

func (n *JSONTableRelation) Span() ast.Loc { return n.Loc }
func (*JSONTableRelation) relationNode()   {}

// JSONTableColumnKind classifies one jsonTableColumn alternative.
type JSONTableColumnKind int

const (
	// JSONTableOrdinalityColumn is `identifier FOR ORDINALITY`.
	JSONTableOrdinalityColumn JSONTableColumnKind = iota
	// JSONTableValueColumn is `identifier type [PATH string] [behaviors]`.
	JSONTableValueColumn
	// JSONTableQueryColumn is `identifier type FORMAT jsonRepresentation …`.
	JSONTableQueryColumn
	// JSONTableNestedColumn is `NESTED [PATH] string [AS identifier] COLUMNS (…)`.
	JSONTableNestedColumn
)

// JSONTableColumn is one column of a JSON_TABLE COLUMNS list. Kind selects the
// alternative that populated it:
//
//   - ordinality: Name.
//   - value: Name, Type, Path, OnEmpty, OnError.
//   - query: Name, Type, Format, Path, Wrapper, WrapperArray, Quotes,
//     QuotesOnScalar, QueryOnEmpty, QueryOnError.
//   - nested: Path, Alias, Columns.
//
// Path is nil when no PATH clause is present; an empty path string is a
// present path. The behavior fields are nil when the clause is absent.
type JSONTableColumn struct {
	Kind           JSONTableColumnKind
	Name           *ast.Identifier    // ordinality / value / query column name
	Type           *DataType          // value / query column type
	Path           *string            // PATH string, nil when absent
	Format         *JSONFormat        // query column FORMAT jsonRepresentation
	Wrapper        string             // query column wrapper, "" when absent
	WrapperArray   bool               // the wrapper's ARRAY keyword
	Quotes         string             // "", "KEEP", or "OMIT"
	QuotesOnScalar bool               // the QUOTES clause's ON SCALAR TEXT_STRING suffix
	OnEmpty        *JSONValueBehavior // value column behavior
	OnError        *JSONValueBehavior // value column behavior
	QueryOnEmpty   *JSONQueryBehavior // query column behavior
	QueryOnError   *JSONQueryBehavior // query column behavior
	Alias          *ast.Identifier    // nested column AS identifier
	Columns        []JSONTableColumn  // nested column list
	Loc            ast.Loc
}

// JSONTablePlanKind classifies one jsonTableSpecificPlan alternative.
type JSONTablePlanKind int

const (
	// JSONTableLeafPlan is `jsonTablePathName`.
	JSONTableLeafPlan JSONTablePlanKind = iota
	// JSONTableJoinPlan is `jsonTablePathName (OUTER | INNER) planPrimary`.
	JSONTableJoinPlan
	// JSONTableUnionPlan is a `planPrimary UNION planPrimary …` chain.
	JSONTableUnionPlan
	// JSONTableCrossPlan is a `planPrimary CROSS planPrimary …` chain.
	JSONTableCrossPlan
)

// JSONTableSpecificPlan is one jsonTableSpecificPlan. For a leaf plan PathName
// is the path name; for a join plan PathName is the left path name, JoinType is
// "OUTER" or "INNER", and Right is the right plan primary; for a UNION or CROSS
// plan Operands holds the two or more operands in order.
type JSONTableSpecificPlan struct {
	Kind     JSONTablePlanKind
	PathName *ast.Identifier
	JoinType string // join plan: "OUTER" or "INNER"
	Right    *JSONTablePlanPrimary
	Operands []*JSONTablePlanPrimary
	Loc      ast.Loc
}

// JSONTablePlanPrimary is one planPrimary: a path name or a parenthesized
// jsonTableSpecificPlan. Exactly one of PathName / Plan is non-nil.
type JSONTablePlanPrimary struct {
	PathName *ast.Identifier
	Plan     *JSONTableSpecificPlan
	Loc      ast.Loc
}

// JSONTableDefaultPlan is `PLAN DEFAULT ( (OUTER|INNER) [, (UNION|CROSS)] |
// (UNION|CROSS) [, (OUTER|INNER)] )`. First is the leading keyword; Second is
// the keyword after the comma, or "" when absent.
type JSONTableDefaultPlan struct {
	First  string // "OUTER", "INNER", "UNION", or "CROSS"
	Second string // "", "OUTER", "INNER", "UNION", or "CROSS"
	Loc    ast.Loc
}

// ---------------------------------------------------------------------------
// JSON_TABLE
// ---------------------------------------------------------------------------

// parseJSONTable parses the `JSON_TABLE ( … )` relation primary (JSON_TABLE is
// current).
func (p *Parser) parseJSONTable() (Relation, error) {
	startTok := p.advance() // consume JSON_TABLE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	path, err := p.parseJSONPathInvocation()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCOLUMNS); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	columns, err := p.parseJSONTableColumns()
	if err != nil {
		return nil, err
	}
	columnsClose, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	rel := &JSONTableRelation{
		Path:    path,
		Columns: columns,
		Loc:     ast.Loc{Start: startTok.Loc.Start, End: columnsClose.Loc.End},
	}

	if p.cur.Kind == kwPLAN {
		p.advance() // consume PLAN
		if p.cur.Kind == kwDEFAULT {
			plan, err := p.parseJSONTableDefaultPlan()
			if err != nil {
				return nil, err
			}
			rel.DefaultPlan = plan
			rel.Loc.End = plan.Loc.End
		} else {
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			plan, err := p.parseJSONTableSpecificPlan()
			if err != nil {
				return nil, err
			}
			planClose, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			rel.Plan = plan
			rel.Loc.End = planClose.Loc.End
		}
	}

	if p.cur.Kind == kwERROR || p.cur.Kind == kwEMPTY {
		if p.cur.Kind == kwERROR {
			rel.OnError = "ERROR"
		} else {
			rel.OnError = "EMPTY"
		}
		p.advance() // consume ERROR or EMPTY
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwERROR); err != nil {
			return nil, err
		}
		rel.Loc.End = p.prev.Loc.End
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	rel.Loc.End = closeTok.Loc.End
	return rel, nil
}

// parseJSONTableColumns parses `jsonTableColumn (',' jsonTableColumn)*`.
func (p *Parser) parseJSONTableColumns() ([]JSONTableColumn, error) {
	first, err := p.parseJSONTableColumn()
	if err != nil {
		return nil, err
	}
	columns := []JSONTableColumn{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseJSONTableColumn()
		if err != nil {
			return nil, err
		}
		columns = append(columns, next)
	}
	return columns, nil
}

// parseJSONTableColumn parses one jsonTableColumn alternative.
func (p *Parser) parseJSONTableColumn() (JSONTableColumn, error) {
	if err := p.enterNesting(); err != nil {
		return JSONTableColumn{}, err
	}
	defer p.leaveNesting()

	// NESTED is a non-reserved keyword (Trino's nonReserved rule), so it starts
	// the nestedColumns alternative only when PATH or the path string follows.
	// Otherwise ANTLR reads it as an ordinary column identifier: `NESTED zzz`
	// is the value column `NESTED` of type `zzz`, and the grammar rejects what
	// follows at its first unexpected token. Oracle-confirmed by the near-miss
	// rows for `NESTED COLUMNS(…)`, `NESTED zzz …`, and `NESTED PATTERN …`.
	if p.cur.Kind == kwNESTED {
		next := p.peekNext().Kind
		if next == kwPATH || next == tokString || next == tokUnicodeString {
			return p.parseNestedJSONTableColumn()
		}
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return JSONTableColumn{}, err
	}
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		ordTok, err := p.expect(kwORDINALITY)
		if err != nil {
			return JSONTableColumn{}, err
		}
		return JSONTableColumn{
			Kind: JSONTableOrdinalityColumn,
			Name: name,
			Loc:  ast.Loc{Start: name.Loc.Start, End: ordTok.Loc.End},
		}, nil
	}

	typ, err := p.parseType()
	if err != nil {
		return JSONTableColumn{}, err
	}
	col := JSONTableColumn{Name: name, Type: typ, Loc: ast.Loc{Start: name.Loc.Start, End: typ.Loc.End}}

	if p.cur.Kind == kwFORMAT {
		col.Kind = JSONTableQueryColumn
		format, err := p.parseJSONFormat()
		if err != nil {
			return JSONTableColumn{}, err
		}
		col.Format = format
		col.Loc.End = format.Loc.End

		if col.Path, err = p.parseJSONTablePath(); err != nil {
			return JSONTableColumn{}, err
		}
		if col.Path != nil {
			col.Loc.End = p.prev.Loc.End
		}

		if wrapper, array, ok := p.tryJSONQueryWrapper(); ok {
			col.Wrapper = wrapper
			col.WrapperArray = array
			wrapperTok, err := p.expect(kwWRAPPER)
			if err != nil {
				return JSONTableColumn{}, err
			}
			col.Loc.End = wrapperTok.Loc.End
		}

		if p.cur.Kind == kwKEEP || p.cur.Kind == kwOMIT {
			if p.cur.Kind == kwKEEP {
				col.Quotes = "KEEP"
			} else {
				col.Quotes = "OMIT"
			}
			p.advance() // consume KEEP or OMIT
			if _, err := p.expect(kwQUOTES); err != nil {
				return JSONTableColumn{}, err
			}
			col.Loc.End = p.prev.Loc.End
			// After QUOTES an ON can only begin the optional ON SCALAR STRING
			// suffix (behavior clauses start with ERROR/NULL/EMPTY), so ON
			// commits to the suffix and a bad next keyword is reported there.
			if p.cur.Kind == kwON {
				p.advance() // consume ON
				if _, err := p.expect(kwSCALAR); err != nil {
					return JSONTableColumn{}, err
				}
				if _, err := p.expect(kwSTRING); err != nil {
					return JSONTableColumn{}, err
				}
				col.QuotesOnScalar = true
				col.Loc.End = p.prev.Loc.End
			}
		}

		for p.startsJSONQueryBehavior() {
			behavior, slot, err := p.parseJSONQueryBehaviorClause()
			if err != nil {
				return JSONTableColumn{}, err
			}
			switch slot {
			case kwEMPTY:
				if col.QueryOnEmpty != nil || col.QueryOnError != nil {
					return JSONTableColumn{}, p.exprErrorAt("duplicate or misordered ON EMPTY in JSON_TABLE column")
				}
				col.QueryOnEmpty = behavior
			case kwERROR:
				if col.QueryOnError != nil {
					return JSONTableColumn{}, p.exprErrorAt("duplicate ON ERROR in JSON_TABLE column")
				}
				col.QueryOnError = behavior
			}
			col.Loc.End = behavior.Loc.End
		}
		return col, nil
	}

	col.Kind = JSONTableValueColumn
	if col.Path, err = p.parseJSONTablePath(); err != nil {
		return JSONTableColumn{}, err
	}
	if col.Path != nil {
		col.Loc.End = p.prev.Loc.End
	}
	for p.startsJSONValueBehavior() {
		behavior, slot, err := p.parseJSONValueBehaviorClause()
		if err != nil {
			return JSONTableColumn{}, err
		}
		switch slot {
		case kwEMPTY:
			if col.OnEmpty != nil || col.OnError != nil {
				return JSONTableColumn{}, p.exprErrorAt("duplicate or misordered ON EMPTY in JSON_TABLE column")
			}
			col.OnEmpty = behavior
		case kwERROR:
			if col.OnError != nil {
				return JSONTableColumn{}, p.exprErrorAt("duplicate ON ERROR in JSON_TABLE column")
			}
			col.OnError = behavior
		}
		col.Loc.End = behavior.Loc.End
	}
	return col, nil
}

// parseNestedJSONTableColumn parses `NESTED PATH? string (AS identifier)?
// COLUMNS '(' jsonTableColumn (',' jsonTableColumn)* ')'` (NESTED is current).
func (p *Parser) parseNestedJSONTableColumn() (JSONTableColumn, error) {
	nestedTok := p.advance() // consume NESTED
	path, err := p.parseJSONTablePath()
	if err != nil {
		return JSONTableColumn{}, err
	}
	if path == nil {
		// NESTED's PATH keyword is optional, but its string is not: without a
		// PATH keyword the string itself must follow NESTED.
		if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
			return JSONTableColumn{}, p.exprErrorAt("expected a string literal JSON path")
		}
		pathTok := p.advance()
		s := pathTok.Str
		path = &s
	}
	col := JSONTableColumn{
		Kind: JSONTableNestedColumn,
		Path: path,
		Loc:  ast.Loc{Start: nestedTok.Loc.Start, End: p.prev.Loc.End},
	}
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		alias, err := p.parseIdentifier()
		if err != nil {
			return JSONTableColumn{}, err
		}
		col.Alias = alias
		col.Loc.End = alias.Loc.End
	}
	if _, err := p.expect(kwCOLUMNS); err != nil {
		return JSONTableColumn{}, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return JSONTableColumn{}, err
	}
	columns, err := p.parseJSONTableColumns()
	if err != nil {
		return JSONTableColumn{}, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return JSONTableColumn{}, err
	}
	col.Columns = columns
	col.Loc.End = closeTok.Loc.End
	return col, nil
}

// parseJSONTablePath parses the optional `PATH string` clause of a value, query,
// or nested column. It returns nil (consuming nothing) when PATH is absent.
func (p *Parser) parseJSONTablePath() (*string, error) {
	if p.cur.Kind != kwPATH {
		return nil, nil
	}
	p.advance() // consume PATH
	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return nil, p.exprErrorAt("expected a string literal JSON path")
	}
	pathTok := p.advance()
	s := pathTok.Str
	return &s, nil
}

// ---------------------------------------------------------------------------
// PLAN
// ---------------------------------------------------------------------------

// parseJSONTableDefaultPlan parses `DEFAULT '(' jsonTableDefaultPlan ')'`
// (PLAN has been consumed and DEFAULT is current).
func (p *Parser) parseJSONTableDefaultPlan() (*JSONTableDefaultPlan, error) {
	defaultTok := p.advance() // consume DEFAULT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	first, err := p.parseJSONTableDefaultPlanKeywordAny()
	if err != nil {
		return nil, err
	}
	plan := &JSONTableDefaultPlan{First: first, Loc: ast.Loc{Start: defaultTok.Loc.Start}}
	if p.cur.Kind == int(',') {
		p.advance() // consume ','
		// After the comma the grammar allows the other pair: a direction
		// (OUTER/INNER) after a set operator (UNION/CROSS), or a set operator
		// after a direction.
		direction := first == "OUTER" || first == "INNER"
		second, err := p.parseJSONTableDefaultPlanKeyword(!direction)
		if err != nil {
			return nil, err
		}
		plan.Second = second
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	plan.Loc.End = closeTok.Loc.End
	return plan, nil
}

// parseJSONTableDefaultPlanKeywordAny parses the leading keyword of a
// jsonTableDefaultPlan, which may be a direction or a set operator.
func (p *Parser) parseJSONTableDefaultPlanKeywordAny() (string, error) {
	switch p.cur.Kind {
	case kwOUTER:
		p.advance()
		return "OUTER", nil
	case kwINNER:
		p.advance()
		return "INNER", nil
	case kwUNION:
		p.advance()
		return "UNION", nil
	case kwCROSS:
		p.advance()
		return "CROSS", nil
	default:
		return "", p.exprErrorAt("expected a JSON_TABLE default plan keyword")
	}
}

// parseJSONTableDefaultPlanKeyword parses the keyword after the comma of a
// jsonTableDefaultPlan. allowDirection selects the pair: OUTER/INNER when true,
// UNION/CROSS when false.
func (p *Parser) parseJSONTableDefaultPlanKeyword(allowDirection bool) (string, error) {
	if allowDirection {
		switch p.cur.Kind {
		case kwOUTER:
			p.advance()
			return "OUTER", nil
		case kwINNER:
			p.advance()
			return "INNER", nil
		}
	} else {
		switch p.cur.Kind {
		case kwUNION:
			p.advance()
			return "UNION", nil
		case kwCROSS:
			p.advance()
			return "CROSS", nil
		}
	}
	return "", p.exprErrorAt("expected a JSON_TABLE default plan keyword")
}

// parseJSONTableSpecificPlan parses one jsonTableSpecificPlan.
func (p *Parser) parseJSONTableSpecificPlan() (*JSONTableSpecificPlan, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	left, err := p.parseJSONTablePlanPrimary()
	if err != nil {
		return nil, err
	}

	switch p.cur.Kind {
	case kwUNION, kwCROSS:
		kind := JSONTableUnionPlan
		opKind := kwUNION
		if p.cur.Kind == kwCROSS {
			kind = JSONTableCrossPlan
			opKind = kwCROSS
		}
		p.advance() // consume UNION or CROSS
		right, err := p.parseJSONTablePlanPrimary()
		if err != nil {
			return nil, err
		}
		plan := &JSONTableSpecificPlan{
			Kind:     kind,
			Operands: []*JSONTablePlanPrimary{left, right},
			Loc:      ast.Loc{Start: left.Loc.Start, End: right.Loc.End},
		}
		for p.cur.Kind == opKind {
			p.advance() // consume the next UNION or CROSS
			next, err := p.parseJSONTablePlanPrimary()
			if err != nil {
				return nil, err
			}
			plan.Operands = append(plan.Operands, next)
			plan.Loc.End = next.Loc.End
		}
		return plan, nil

	case kwOUTER, kwINNER:
		if left.PathName == nil {
			return nil, p.exprErrorAt("expected a path name before OUTER or INNER")
		}
		joinType := "OUTER"
		if p.cur.Kind == kwINNER {
			joinType = "INNER"
		}
		p.advance() // consume OUTER or INNER
		right, err := p.parseJSONTablePlanPrimary()
		if err != nil {
			return nil, err
		}
		return &JSONTableSpecificPlan{
			Kind:     JSONTableJoinPlan,
			PathName: left.PathName,
			JoinType: joinType,
			Right:    right,
			Loc:      ast.Loc{Start: left.Loc.Start, End: right.Loc.End},
		}, nil

	default:
		if left.PathName == nil {
			return nil, p.exprErrorAt("expected UNION or CROSS after a parenthesized plan")
		}
		return &JSONTableSpecificPlan{
			Kind:     JSONTableLeafPlan,
			PathName: left.PathName,
			Loc:      left.Loc,
		}, nil
	}
}

// parseJSONTablePlanPrimary parses `jsonTablePathName | '('
// jsonTableSpecificPlan ')'`.
func (p *Parser) parseJSONTablePlanPrimary() (*JSONTablePlanPrimary, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	if p.cur.Kind == int('(') {
		openTok := p.advance() // consume '('
		plan, err := p.parseJSONTableSpecificPlan()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		return &JSONTablePlanPrimary{
			Plan: plan,
			Loc:  ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &JSONTablePlanPrimary{PathName: name, Loc: name.Loc}, nil
}
