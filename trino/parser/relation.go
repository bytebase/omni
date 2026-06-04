package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `parser-select` DAG node (with select.go, setops.go,
// ctes.go and window.go): it implements Trino's FROM-clause relation grammar —
// the `relation`, `joinType`, `joinCriteria`, `sampledRelation`,
// `aliasedRelation`, `relationPrimary`, and `tableFunctionCall` rules.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	relation
//	    : left (CROSS JOIN right
//	           | joinType JOIN rightRelation joinCriteria
//	           | NATURAL joinType JOIN right)        # joinRelation
//	    | sampledRelation                            # relationDefault ;
//	joinType        : INNER? | (LEFT|RIGHT|FULL) OUTER? ;
//	joinCriteria    : ON booleanExpression | USING ( identifier (, …)* ) ;
//	sampledRelation : patternRecognition (TABLESAMPLE sampleType ( percentage ))? ;
//	aliasedRelation : relationPrimary (AS? identifier columnAliases?)? ;
//	relationPrimary
//	    : qualifiedName queryPeriod?                 # tableName
//	    | ( query )                                  # subqueryRelation
//	    | UNNEST ( expression (, …)* ) (WITH ORDINALITY)?  # unnest
//	    | LATERAL ( query )                          # lateral
//	    | TABLE ( tableFunctionCall )                # tableFunctionInvocation
//	    | ( relation )                               # parenthesizedRelation ;
//
// The implementation is adjudicated against the live Trino 481 oracle.
//
// JOIN associativity & the CROSS/NATURAL operand asymmetry (oracle-confirmed via
// the grammar shape): a JOIN chain is LEFT-associative. The qualified-join
// alternative's right operand is `relation` (recursively), but CROSS JOIN and
// NATURAL JOIN take a `sampledRelation` right operand — so `a JOIN b JOIN c` and
// `a CROSS JOIN b CROSS JOIN c` both group left as `((a J b) J c)`. parseRelation
// parses a sampledRelation then folds a left-associative loop of join suffixes.
//
// Node-scope boundary D4 (select.go): sampledRelation's patternRecognition adds
// the MATCH_RECOGNIZE row-pattern subsystem, which is the parser-match-recognize
// DAG node (row_pattern.go). Here a relation primary is its aliasedRelation with
// the optional TABLESAMPLE only; a MATCH_RECOGNIZE following an aliasedRelation
// is left for that node (its absence here means omni rejects MATCH_RECOGNIZE
// until that node lands — tracked separately, not a parser-select divergence).
//
// Node-scope boundary B1 (expr.go): a LATERAL / subquery / TABLE(query) inner
// `query` IS parsed into a real Query tree here (this node owns the query rule),
// unlike expression-embedded subqueries which stay raw-text placeholders.

// ---------------------------------------------------------------------------
// Relation node hierarchy (parser-package; not ast.Node — see select.go header)
// ---------------------------------------------------------------------------

// Relation is the interface implemented by every FROM-clause relation node.
// Span returns the source byte range; concrete fields are reached by a Go type
// switch. It is the result type of the `relation` / `sampledRelation` /
// `relationPrimary` rules.
type Relation interface {
	// Span returns the source byte range covered by the relation.
	Span() ast.Loc
	// relationNode is a marker preventing unrelated types from satisfying Relation.
	relationNode()
}

// JoinType enumerates the join flavors (joinType + the CROSS/NATURAL forms).
type JoinType int

const (
	// JoinInner is `[INNER] JOIN` (the default qualified join).
	JoinInner JoinType = iota
	// JoinLeft is `LEFT [OUTER] JOIN`.
	JoinLeft
	// JoinRight is `RIGHT [OUTER] JOIN`.
	JoinRight
	// JoinFull is `FULL [OUTER] JOIN`.
	JoinFull
	// JoinCross is `CROSS JOIN` (no criteria).
	JoinCross
)

// String returns the canonical spelling of the join type (without OUTER).
func (j JoinType) String() string {
	switch j {
	case JoinLeft:
		return "LEFT"
	case JoinRight:
		return "RIGHT"
	case JoinFull:
		return "FULL"
	case JoinCross:
		return "CROSS"
	default:
		return "INNER"
	}
}

// Join is a binary join (the joinRelation rule): `left <jointype> JOIN right
// [criteria]`. Type is the join flavor; Natural marks a NATURAL join; Outer
// records the explicit OUTER keyword. On is the `ON booleanExpression` criteria
// (nil otherwise); Using is the `USING ( col, … )` criteria (nil otherwise).
// CROSS and NATURAL joins carry no criteria.
type Join struct {
	Type    JoinType
	Natural bool
	Outer   bool
	Left    Relation
	Right   Relation
	On      Expr              // ON condition, nil when absent
	Using   []*ast.Identifier // USING column list, nil when absent
	Loc     ast.Loc
}

func (n *Join) Span() ast.Loc { return n.Loc }
func (*Join) relationNode()   {}

// AliasedRelation is a relation primary with an optional alias and column-alias
// list (the aliasedRelation rule applied to a relationPrimary). It also carries
// the optional TABLESAMPLE (sampledRelation) wrapping the primary. Inner is the
// underlying relationPrimary node. Alias / ColumnAliases are nil when absent.
// Sample is nil when there is no TABLESAMPLE.
type AliasedRelation struct {
	Inner         Relation
	Alias         *ast.Identifier
	ColumnAliases []*ast.Identifier
	Sample        *TableSample // TABLESAMPLE, nil when absent
	Loc           ast.Loc
}

func (n *AliasedRelation) Span() ast.Loc { return n.Loc }
func (*AliasedRelation) relationNode()   {}

// TableSample is the `TABLESAMPLE (BERNOULLI|SYSTEM) ( percentage )` modifier
// (the sampledRelation TABLESAMPLE clause). Method is "BERNOULLI" or "SYSTEM";
// Percentage is the sample-rate expression.
type TableSample struct {
	Method     string
	Percentage Expr
	Loc        ast.Loc
}

// TableRelation is the `qualifiedName queryPeriod?` relation primary (the
// tableName rule): a base table or view reference, optionally with a time-travel
// queryPeriod. Period is nil when absent.
type TableRelation struct {
	Name   *ast.QualifiedName
	Period *QueryPeriod // FOR (TIMESTAMP|VERSION) AS OF …, nil when absent
	Loc    ast.Loc
}

func (n *TableRelation) Span() ast.Loc { return n.Loc }
func (*TableRelation) relationNode()   {}

// QueryPeriod is the `FOR (TIMESTAMP|VERSION) AS OF valueExpression` time-travel
// clause (the queryPeriod rule). RangeType is "TIMESTAMP" or "VERSION"; End is
// the as-of value expression.
type QueryPeriod struct {
	RangeType string // "TIMESTAMP" or "VERSION"
	End       Expr
	Loc       ast.Loc
}

// SubqueryRelation is the `( query )` relation primary (the subqueryRelation
// rule): a derived table. Query is the parsed inner query (B1 boundary does not
// apply — this node owns the query rule).
type SubqueryRelation struct {
	Query *Query
	Loc   ast.Loc
}

func (n *SubqueryRelation) Span() ast.Loc { return n.Loc }
func (*SubqueryRelation) relationNode()   {}

// UnnestRelation is the `UNNEST ( expression (, …)* ) (WITH ORDINALITY)?`
// relation primary (the unnest rule): expands array/map expressions into rows.
// WithOrdinality marks the trailing WITH ORDINALITY.
type UnnestRelation struct {
	Exprs          []Expr
	WithOrdinality bool
	Loc            ast.Loc
}

func (n *UnnestRelation) Span() ast.Loc { return n.Loc }
func (*UnnestRelation) relationNode()   {}

// LateralRelation is the `LATERAL ( query )` relation primary (the lateral
// rule): a correlated subquery that may reference earlier FROM items. Query is
// the parsed inner query.
type LateralRelation struct {
	Query *Query
	Loc   ast.Loc
}

func (n *LateralRelation) Span() ast.Loc { return n.Loc }
func (*LateralRelation) relationNode()   {}

// ParenRelation is the `( relation )` relation primary (the parenthesizedRelation
// rule): a parenthesized relation, used to group joins explicitly. Inner is the
// wrapped relation.
type ParenRelation struct {
	Inner Relation
	Loc   ast.Loc
}

func (n *ParenRelation) Span() ast.Loc { return n.Loc }
func (*ParenRelation) relationNode()   {}

// TableFunctionRelation is the `TABLE ( tableFunctionCall )` relation primary
// (the tableFunctionInvocation rule): a polymorphic table function invocation.
// Call is the parsed table-function call.
type TableFunctionRelation struct {
	Call *TableFunctionCall
	Loc  ast.Loc
}

func (n *TableFunctionRelation) Span() ast.Loc { return n.Loc }
func (*TableFunctionRelation) relationNode()   {}

// TableFunctionCall is `qualifiedName ( tableFunctionArgument (, …)* )
// (COPARTITION …)?` (the tableFunctionCall rule). Name is the function name;
// Args is the argument list; Copartition holds each COPARTITION group's tables.
type TableFunctionCall struct {
	Name        *ast.QualifiedName
	Args        []TableFunctionArgument
	Copartition [][]*ast.QualifiedName // each inner slice is one COPARTITION group
	Loc         ast.Loc
}

// TableFunctionArgKind classifies a TableFunctionArgument's value.
type TableFunctionArgKind int

const (
	// TFArgExpr is a scalar `expression` argument.
	TFArgExpr TableFunctionArgKind = iota
	// TFArgTable is a `TABLE ( … )` table argument (tableArgument).
	TFArgTable
	// TFArgDescriptor is a `DESCRIPTOR ( field, … )` or `CAST(NULL AS DESCRIPTOR)`
	// argument (descriptorArgument).
	TFArgDescriptor
)

// TableFunctionArgument is one argument of a table-function call (the
// tableFunctionArgument rule): an optional `name =>` followed by a scalar
// expression, a table argument, or a descriptor argument. Name is the optional
// named-argument identifier (nil for a positional argument).
type TableFunctionArgument struct {
	Name       *ast.Identifier // name => …, nil for positional
	Kind       TableFunctionArgKind
	Expr       Expr           // for TFArgExpr
	Table      *TableArgument // for TFArgTable
	Descriptor *DescriptorArg // for TFArgDescriptor
	Loc        ast.Loc
}

// TableArgument is a table-function `TABLE ( name | query )` argument with the
// optional PARTITION BY / PRUNE|KEEP WHEN EMPTY / ORDER BY modifiers (the
// tableArgument + tableArgumentRelation rules). Exactly one of Name / Query is
// non-nil. Alias / ColumnAliases carry the optional `AS? identifier (cols)?`.
type TableArgument struct {
	Name          *ast.QualifiedName // TABLE ( qualifiedName )
	Query         *Query             // TABLE ( query )
	Alias         *ast.Identifier
	ColumnAliases []*ast.Identifier
	PartitionBy   []Expr     // PARTITION BY …, nil when absent
	Pruning       string     // "", "PRUNE WHEN EMPTY", or "KEEP WHEN EMPTY"
	OrderBy       []SortItem // ORDER BY …, nil when absent
	Loc           ast.Loc
}

// DescriptorArg is a `DESCRIPTOR ( field, … )` or `CAST(NULL AS DESCRIPTOR)`
// table-function argument (the descriptorArgument rule). NullCast marks the
// CAST(NULL AS DESCRIPTOR) form (Fields is nil then). Otherwise Fields lists the
// descriptor fields.
type DescriptorArg struct {
	NullCast bool
	Fields   []DescriptorField
	Loc      ast.Loc
}

// DescriptorField is one `identifier type?` of a DESCRIPTOR argument (the
// descriptorField rule). Type is nil when the field has no declared type.
type DescriptorField struct {
	Name *ast.Identifier
	Type *DataType // nil when absent
	Loc  ast.Loc
}

// ---------------------------------------------------------------------------
// relation list / relation (joins)
// ---------------------------------------------------------------------------

// parseRelationList parses `relation (, relation)*` — the comma-separated FROM
// item list. A comma between relations is an implicit cross join; each item is a
// full `relation` (so `FROM a JOIN b, c` is `(a JOIN b)` and `c`).
func (p *Parser) parseRelationList() ([]Relation, error) {
	first, err := p.parseRelation()
	if err != nil {
		return nil, err
	}
	rels := []Relation{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseRelation()
		if err != nil {
			return nil, err
		}
		rels = append(rels, next)
	}
	return rels, nil
}

// parseRelation parses one `relation`: a sampledRelation followed by a
// left-associative chain of join suffixes (CROSS JOIN, qualified JOIN with
// criteria, or NATURAL JOIN). The left-associative loop mirrors ANTLR's
// left-recursive `relation` rule.
func (p *Parser) parseRelation() (Relation, error) {
	left, err := p.parseSampledRelation()
	if err != nil {
		return nil, err
	}
	for {
		joined, ok, err := p.tryParseJoinSuffix(left)
		if err != nil {
			return nil, err
		}
		if !ok {
			return left, nil
		}
		left = joined
	}
}

// tryParseJoinSuffix parses one join suffix on `left` if the current token
// begins a join (CROSS / NATURAL / a joinType keyword / JOIN). It returns
// ok=false (no join) when no join keyword leads. The three join shapes:
//
//   - CROSS JOIN sampledRelation         (no criteria)
//   - joinType JOIN relation joinCriteria
//   - NATURAL joinType JOIN sampledRelation (no criteria)
func (p *Parser) tryParseJoinSuffix(left Relation) (Relation, bool, error) {
	switch p.cur.Kind {
	case kwCROSS:
		return p.parseCrossJoin(left)
	case kwNATURAL:
		return p.parseNaturalJoin(left)
	case kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL:
		j, err := p.parseQualifiedJoin(left)
		if err != nil {
			return nil, false, err
		}
		return j, true, nil
	default:
		return nil, false, nil
	}
}

// parseCrossJoin parses `CROSS JOIN sampledRelation` on `left` (CROSS is
// current). A cross join carries no criteria and takes a sampledRelation (not a
// full relation) right operand, keeping a chain left-associative.
func (p *Parser) parseCrossJoin(left Relation) (Relation, bool, error) {
	p.advance() // consume CROSS
	if _, err := p.expect(kwJOIN); err != nil {
		return nil, false, err
	}
	right, err := p.parseSampledRelation()
	if err != nil {
		return nil, false, err
	}
	return &Join{
		Type:  JoinCross,
		Left:  left,
		Right: right,
		Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
	}, true, nil
}

// parseNaturalJoin parses `NATURAL joinType JOIN sampledRelation` on `left`
// (NATURAL is current). A natural join carries no explicit criteria and takes a
// sampledRelation right operand.
func (p *Parser) parseNaturalJoin(left Relation) (Relation, bool, error) {
	p.advance() // consume NATURAL
	jt, outer := p.parseJoinType()
	if _, err := p.expect(kwJOIN); err != nil {
		return nil, false, err
	}
	right, err := p.parseSampledRelation()
	if err != nil {
		return nil, false, err
	}
	return &Join{
		Type:    jt,
		Natural: true,
		Outer:   outer,
		Left:    left,
		Right:   right,
		Loc:     ast.Loc{Start: left.Span().Start, End: right.Span().End},
	}, true, nil
}

// parseQualifiedJoin parses `joinType JOIN relation joinCriteria` on `left`
// (the leading joinType keyword or JOIN is current). The right operand is a full
// relation; the criteria (ON / USING) is mandatory for a qualified join.
func (p *Parser) parseQualifiedJoin(left Relation) (Relation, error) {
	jt, outer := p.parseJoinType()
	if _, err := p.expect(kwJOIN); err != nil {
		return nil, err
	}
	right, err := p.parseRelation()
	if err != nil {
		return nil, err
	}
	join := &Join{
		Type:  jt,
		Outer: outer,
		Left:  left,
		Right: right,
		Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
	}
	if err := p.parseJoinCriteria(join); err != nil {
		return nil, err
	}
	return join, nil
}

// parseJoinType parses the optional `joinType` keyword sequence (INNER? |
// (LEFT|RIGHT|FULL) OUTER?) and returns the flavor plus whether an explicit
// OUTER was present. A bare JOIN (no leading keyword) is INNER. INNER never
// takes OUTER.
func (p *Parser) parseJoinType() (JoinType, bool) {
	switch p.cur.Kind {
	case kwINNER:
		p.advance() // consume INNER
		return JoinInner, false
	case kwLEFT:
		p.advance()
		return JoinLeft, p.matchOptionalOuter()
	case kwRIGHT:
		p.advance()
		return JoinRight, p.matchOptionalOuter()
	case kwFULL:
		p.advance()
		return JoinFull, p.matchOptionalOuter()
	default:
		// Bare JOIN — INNER.
		return JoinInner, false
	}
}

// matchOptionalOuter consumes an optional OUTER keyword, returning whether it
// was present.
func (p *Parser) matchOptionalOuter() bool {
	_, ok := p.match(kwOUTER)
	return ok
}

// parseJoinCriteria parses the mandatory `joinCriteria` of a qualified join
// (ON booleanExpression | USING ( identifier (, …)* )) and stores it on the
// join. ON / USING is the current token.
func (p *Parser) parseJoinCriteria(join *Join) error {
	switch p.cur.Kind {
	case kwON:
		p.advance() // consume ON
		cond, err := p.parseExpr()
		if err != nil {
			return err
		}
		join.On = cond
		join.Loc.End = cond.Span().End
		return nil
	case kwUSING:
		p.advance() // consume USING
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		first, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		cols := []*ast.Identifier{first}
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			cols = append(cols, next)
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return err
		}
		join.Using = cols
		join.Loc.End = closeTok.Loc.End
		return nil
	default:
		return p.exprErrorAt("expected ON or USING join criteria")
	}
}

// ---------------------------------------------------------------------------
// sampledRelation / aliasedRelation
// ---------------------------------------------------------------------------

// parseSampledRelation parses a `sampledRelation`: an aliasedRelation with an
// optional trailing `TABLESAMPLE (BERNOULLI|SYSTEM) ( percentage )`. (The
// MATCH_RECOGNIZE patternRecognition layer between them is the
// parser-match-recognize node, D4.)
func (p *Parser) parseSampledRelation() (Relation, error) {
	rel, err := p.parseAliasedRelation()
	if err != nil {
		return nil, err
	}
	if p.cur.Kind != kwTABLESAMPLE {
		return rel, nil
	}
	sample, err := p.parseTableSample()
	if err != nil {
		return nil, err
	}
	// Attach the sample to the aliased relation; if rel is not an AliasedRelation
	// (it always is, parseAliasedRelation wraps every primary) fall back to a
	// fresh wrapper so the sample is never dropped.
	if ar, ok := rel.(*AliasedRelation); ok {
		ar.Sample = sample
		ar.Loc.End = sample.Loc.End
		return ar, nil
	}
	return &AliasedRelation{
		Inner:  rel,
		Sample: sample,
		Loc:    ast.Loc{Start: rel.Span().Start, End: sample.Loc.End},
	}, nil
}

// parseTableSample parses `TABLESAMPLE (BERNOULLI|SYSTEM) ( percentage )`
// (TABLESAMPLE is current).
func (p *Parser) parseTableSample() (*TableSample, error) {
	sampleTok := p.advance() // consume TABLESAMPLE
	methodTok, ok := p.match(kwBERNOULLI, kwSYSTEM)
	if !ok {
		return nil, p.exprErrorAt("expected BERNOULLI or SYSTEM after TABLESAMPLE")
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	pct, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &TableSample{
		Method:     methodTok.Str,
		Percentage: pct,
		Loc:        ast.Loc{Start: sampleTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseAliasedRelation parses an `aliasedRelation`: a relationPrimary with an
// optional `AS? identifier columnAliases?` alias. The alias identifier is
// optional and (per the grammar) follows an optional AS; a column-alias list may
// follow the alias name.
func (p *Parser) parseAliasedRelation() (Relation, error) {
	primary, err := p.parseRelationPrimary()
	if err != nil {
		return nil, err
	}

	alias, cols, end, ok, err := p.tryParseRelationAlias()
	if err != nil {
		return nil, err
	}
	if !ok {
		// No alias: still wrap in AliasedRelation so a later TABLESAMPLE can attach
		// uniformly, and so consumers see a single relation-primary node shape.
		return &AliasedRelation{Inner: primary, Loc: primary.Span()}, nil
	}
	return &AliasedRelation{
		Inner:         primary,
		Alias:         alias,
		ColumnAliases: cols,
		Loc:           ast.Loc{Start: primary.Span().Start, End: end},
	}, nil
}

// tryParseRelationAlias parses the optional `AS? identifier columnAliases?` after
// a relation primary. It returns ok=false when no alias is present. An explicit
// AS makes the identifier mandatory; without AS, an identifier-start token is
// taken as the alias UNLESS it is a non-reserved clause keyword sitting in
// clause-introducing position (see atRelationClauseStart) — those introduce the
// surrounding querySpecification's WINDOW / LIMIT / OFFSET / FETCH clause or the
// sampledRelation's TABLESAMPLE, not an alias. (Most clause/join keywords —
// WHERE/GROUP/HAVING/ORDER/JOIN/ON/… — are reserved, so isIdentifierStart already
// excludes them; only WINDOW/LIMIT/OFFSET/FETCH/TABLESAMPLE are non-reserved and
// need the extra guard. A bare non-reserved keyword — `FROM t window`,
// `FROM t limit` — is still a valid alias and is taken here.)
func (p *Parser) tryParseRelationAlias() (alias *ast.Identifier, cols []*ast.Identifier, end int, ok bool, err error) {
	if p.cur.Kind == kwAS {
		p.advance() // consume AS — an identifier MUST follow (even a clause keyword)
		id, e := p.parseIdentifier()
		if e != nil {
			return nil, nil, 0, false, e
		}
		alias = id
		end = id.Loc.End
	} else if isIdentifierStart(p.cur.Kind) && !p.atRelationClauseStart() {
		alias = identFromToken(p.advance())
		end = alias.Loc.End
	} else {
		return nil, nil, 0, false, nil
	}

	// Optional column-alias list `( col, … )`.
	if p.cur.Kind == int('(') {
		c, e, cerr := p.parseColumnAliases()
		if cerr != nil {
			return nil, nil, 0, false, cerr
		}
		cols = c
		end = e
	}
	return alias, cols, end, true, nil
}

// atRelationClauseStart reports whether the cursor sits at a NON-RESERVED clause
// keyword that, here, introduces a clause rather than naming a relation alias —
// disambiguated by the token that follows it (matching Trino 481, probed):
//
//	WINDOW      + identifier      → the querySpecification WINDOW clause
//	LIMIT       + (ALL|int|?)     → the queryNoWith LIMIT clause
//	OFFSET      + (int|?)         → the queryNoWith OFFSET clause
//	FETCH       + (FIRST|NEXT)    → the queryNoWith FETCH clause
//	TABLESAMPLE + (BERNOULLI|SYSTEM) → the sampledRelation TABLESAMPLE clause
//
// A bare occurrence (e.g. `FROM t window`, `FROM t limit`) does NOT match and is
// taken as an alias by the caller. The reserved clause/join keywords need no
// guard (isIdentifierStart already rejects them).
func (p *Parser) atRelationClauseStart() bool {
	switch p.cur.Kind {
	case kwWINDOW:
		return isIdentifierStart(p.peekNext().Kind)
	case kwLIMIT:
		nk := p.peekNext().Kind
		return nk == kwALL || nk == tokInteger || nk == tokQuestion || nk == int('?')
	case kwOFFSET:
		nk := p.peekNext().Kind
		return nk == tokInteger || nk == tokQuestion || nk == int('?')
	case kwFETCH:
		nk := p.peekNext().Kind
		return nk == kwFIRST || nk == kwNEXT
	case kwTABLESAMPLE:
		nk := p.peekNext().Kind
		return nk == kwBERNOULLI || nk == kwSYSTEM
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// relationPrimary
// ---------------------------------------------------------------------------

// parseRelationPrimary parses one `relationPrimary`:
//
//   - UNNEST ( … ) [WITH ORDINALITY]      (unnest)
//   - LATERAL ( query )                   (lateral)
//   - TABLE ( tableFunctionCall )         (tableFunctionInvocation)
//   - ( query ) | ( relation )            (subqueryRelation / parenthesizedRelation)
//   - qualifiedName [queryPeriod]         (tableName)
//
// UNNEST and TABLE are reserved keywords, so they unambiguously begin their forms
// (and require the following '('). LATERAL is NON-RESERVED, so it begins the
// lateral form only when directly followed by '('; a bare LATERAL (or LATERAL not
// followed by '(') is an ordinary table-name reference (`FROM lateral`,
// `FROM lateral x`), oracle-confirmed in Trino 481.
//
// JSON_TABLE (a 481 relationPrimary absent from the legacy grammar) is the
// deferred P1 extension D1 (select.go); it is not handled here.
func (p *Parser) parseRelationPrimary() (Relation, error) {
	switch p.cur.Kind {
	case kwUNNEST:
		return p.parseUnnest()
	case kwLATERAL:
		// LATERAL is non-reserved: `LATERAL ( query )` is the lateral form, but a
		// LATERAL not followed by '(' is a table named "lateral".
		if p.peekNext().Kind == int('(') {
			return p.parseLateral()
		}
		return p.parseTableRelation()
	case kwTABLE:
		return p.parseTableFunctionInvocation()
	case int('('):
		return p.parseParenRelationOrSubquery()
	default:
		return p.parseTableRelation()
	}
}

// parseTableRelation parses `qualifiedName queryPeriod?` (the tableName rule):
// a base table/view reference with an optional time-travel period.
func (p *Parser) parseTableRelation() (Relation, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rel := &TableRelation{Name: name, Loc: name.Loc}
	if p.cur.Kind == kwFOR {
		period, err := p.parseQueryPeriod()
		if err != nil {
			return nil, err
		}
		rel.Period = period
		rel.Loc.End = period.Loc.End
	}
	return rel, nil
}

// parseQueryPeriod parses `FOR (TIMESTAMP|VERSION) AS OF valueExpression`
// (the queryPeriod rule; FOR is current). The as-of value is a valueExpression
// (not a full booleanExpression).
func (p *Parser) parseQueryPeriod() (*QueryPeriod, error) {
	forTok := p.advance() // consume FOR
	rangeTok, ok := p.match(kwTIMESTAMP, kwVERSION)
	if !ok {
		return nil, p.exprErrorAt("expected TIMESTAMP or VERSION after FOR")
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}
	end, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	return &QueryPeriod{
		RangeType: rangeTok.Str,
		End:       end,
		Loc:       ast.Loc{Start: forTok.Loc.Start, End: end.Span().End},
	}, nil
}

// parseUnnest parses `UNNEST ( expression (, …)* ) (WITH ORDINALITY)?` (UNNEST
// is current).
func (p *Parser) parseUnnest() (Relation, error) {
	unnestTok := p.advance() // consume UNNEST
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprs := []Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	rel := &UnnestRelation{
		Exprs: exprs,
		Loc:   ast.Loc{Start: unnestTok.Loc.Start, End: closeTok.Loc.End},
	}
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		ordTok, err := p.expect(kwORDINALITY)
		if err != nil {
			return nil, err
		}
		rel.WithOrdinality = true
		rel.Loc.End = ordTok.Loc.End
	}
	return rel, nil
}

// parseLateral parses `LATERAL ( query )` (LATERAL is current). The inner query
// is parsed into a real Query tree (this node owns the query rule).
func (p *Parser) parseLateral() (Relation, error) {
	lateralTok := p.advance() // consume LATERAL
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &LateralRelation{
		Query: q,
		Loc:   ast.Loc{Start: lateralTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseParenRelationOrSubquery disambiguates a leading '(' between a
// parenthesized query `( query )` (subqueryRelation) and a parenthesized
// relation `( relation )` (parenthesizedRelation; e.g. an explicitly grouped
// join). The token after '(' decides:
//
//   - SELECT / WITH / VALUES / TABLE-name → unambiguously a query; parse one.
//   - a nested '(' → AMBIGUOUS: `((SELECT 1))` is a nested subquery, while
//     `((a JOIN b ON …) JOIN c ON …)` is a nested parenthesized relation. Both
//     start `((`, and the deciding token lies arbitrarily deep, so a checkpoint
//     speculates: try a query first; if that fails (or does not land on the
//     matching ')'), rewind and parse a relation.
//   - anything else (an identifier table name, UNNEST, LATERAL, TABLE(func)) →
//     a relation.
func (p *Parser) parseParenRelationOrSubquery() (Relation, error) {
	openTok := p.advance() // consume '('

	switch {
	case p.startsQueryRelation():
		return p.finishSubqueryRelation(openTok.Loc.Start)
	case p.cur.Kind == int('('):
		// Ambiguous `((…`. Speculatively read a query; fall back to a relation.
		cp := p.checkpoint()
		if q, err := p.parseQuery(); err == nil && p.cur.Kind == int(')') {
			closeTok := p.advance() // consume ')'
			return &SubqueryRelation{
				Query: q,
				Loc:   ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
			}, nil
		}
		p.restore(cp)
		return p.finishParenRelation(openTok.Loc.Start)
	default:
		return p.finishParenRelation(openTok.Loc.Start)
	}
}

// finishSubqueryRelation parses `query )` after the opening '(' of a relation
// subquery (the '(' already consumed; startOffset is its byte offset).
func (p *Parser) finishSubqueryRelation(startOffset int) (Relation, error) {
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &SubqueryRelation{
		Query: q,
		Loc:   ast.Loc{Start: startOffset, End: closeTok.Loc.End},
	}, nil
}

// finishParenRelation parses `relation )` after the opening '(' of a
// parenthesized relation (the '(' already consumed; startOffset is its byte
// offset).
func (p *Parser) finishParenRelation(startOffset int) (Relation, error) {
	inner, err := p.parseRelation()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ParenRelation{
		Inner: inner,
		Loc:   ast.Loc{Start: startOffset, End: closeTok.Loc.End},
	}, nil
}

// startsQueryRelation reports whether the current token UNAMBIGUOUSLY begins a
// `query` in a relation-subquery position: SELECT / WITH / VALUES, or TABLE not
// followed by '(' (the `TABLE qualifiedName` query primary). A nested '(' is
// deliberately EXCLUDED — it is ambiguous between a nested subquery and a nested
// parenthesized relation, and is resolved by speculation in the caller. A
// relation, by contrast, starts with an identifier (table name), UNNEST,
// LATERAL, or TABLE ( … ) (the table-function relation primary).
func (p *Parser) startsQueryRelation() bool {
	switch p.cur.Kind {
	case kwSELECT, kwWITH, kwVALUES:
		return true
	case kwTABLE:
		// TABLE qualifiedName is a query primary; TABLE ( … ) is the table-function
		// relation primary. Decide by the token after TABLE.
		return p.peekNext().Kind != int('(')
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// table function invocation
// ---------------------------------------------------------------------------

// parseTableFunctionInvocation parses `TABLE ( tableFunctionCall )`
// (the tableFunctionInvocation relation primary; TABLE is current). Note this is
// distinct from the `TABLE qualifiedName` query primary — here TABLE is
// immediately followed by '(' and wraps a tableFunctionCall.
func (p *Parser) parseTableFunctionInvocation() (Relation, error) {
	tableTok := p.advance() // consume TABLE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	call, err := p.parseTableFunctionCall()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &TableFunctionRelation{
		Call: call,
		Loc:  ast.Loc{Start: tableTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// atCopartitionClause reports whether the cursor sits at the COPARTITION clause
// of a tableFunctionCall, as opposed to a `copartition => …` named argument.
// COPARTITION is a non-reserved keyword, so a leading COPARTITION followed by
// `=>` is an argument name, not the clause. Oracle-confirmed: Trino 481 accepts
// `TABLE(f(copartition => 1))` as a named-argument call.
func (p *Parser) atCopartitionClause() bool {
	return p.cur.Kind == kwCOPARTITION && p.peekNext().Kind != tokDoubleArrow
}

// parseTableFunctionCall parses `qualifiedName ( tableFunctionArgument (, …)* )
// (COPARTITION copartitionTables (, …)*)?` (the tableFunctionCall rule).
func (p *Parser) parseTableFunctionCall() (*TableFunctionCall, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	call := &TableFunctionCall{Name: name, Loc: ast.Loc{Start: name.Loc.Start}}

	// Argument list (possibly empty). COPARTITION is non-reserved, so a leading
	// COPARTITION is the COPARTITION clause only when it is NOT a `copartition =>`
	// named argument (atCopartitionClause); otherwise it begins an argument.
	if p.cur.Kind != int(')') && !p.atCopartitionClause() {
		first, err := p.parseTableFunctionArgument()
		if err != nil {
			return nil, err
		}
		call.Args = []TableFunctionArgument{first}
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseTableFunctionArgument()
			if err != nil {
				return nil, err
			}
			call.Args = append(call.Args, next)
		}
	}

	// Optional COPARTITION groups (the clause, not a `copartition =>` argument).
	for p.atCopartitionClause() {
		p.advance() // consume COPARTITION
		group, err := p.parseCopartitionTables()
		if err != nil {
			return nil, err
		}
		call.Copartition = append(call.Copartition, group)
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	call.Loc.End = closeTok.Loc.End
	return call, nil
}

// parseTableFunctionArgument parses one `tableFunctionArgument`: an optional
// `name =>` followed by a table argument, a descriptor argument, or a scalar
// expression. The descriptor alternative is tried before expression (its DESCRIPTOR
// keyword leads), and TABLE leads the table argument.
func (p *Parser) parseTableFunctionArgument() (TableFunctionArgument, error) {
	arg := TableFunctionArgument{}
	start := p.cur.Loc.Start

	// Optional `name =>` prefix.
	if isIdentifierStart(p.cur.Kind) && p.peekNext().Kind == tokDoubleArrow {
		arg.Name = identFromToken(p.advance())
		p.advance() // consume =>
	}

	switch {
	case p.cur.Kind == kwTABLE:
		ta, err := p.parseTableArgument()
		if err != nil {
			return TableFunctionArgument{}, err
		}
		arg.Kind = TFArgTable
		arg.Table = ta
		arg.Loc = ast.Loc{Start: start, End: ta.Loc.End}
	case (p.cur.Kind == kwDESCRIPTOR && p.peekNext().Kind == int('(')) || (p.cur.Kind == kwCAST && p.peekNext().Kind == int('(')):
		// DESCRIPTOR(...) or CAST(NULL AS DESCRIPTOR). DESCRIPTOR is non-reserved,
		// so a bare `descriptor` (not followed by '(') is an ordinary column
		// reference handled by the default expression branch — only DESCRIPTOR '('
		// enters here. For CAST, only the CAST(NULL AS DESCRIPTOR) shape is a
		// descriptor argument; an ordinary CAST is a scalar expression. Try the
		// descriptor form, fall back to a scalar expression.
		da, ok, err := p.tryParseDescriptorArg(start)
		if err != nil {
			return TableFunctionArgument{}, err
		}
		if ok {
			arg.Kind = TFArgDescriptor
			arg.Descriptor = da
			arg.Loc = ast.Loc{Start: start, End: da.Loc.End}
		} else {
			e, err := p.parseExpr()
			if err != nil {
				return TableFunctionArgument{}, err
			}
			arg.Kind = TFArgExpr
			arg.Expr = e
			arg.Loc = ast.Loc{Start: start, End: e.Span().End}
		}
	default:
		e, err := p.parseExpr()
		if err != nil {
			return TableFunctionArgument{}, err
		}
		arg.Kind = TFArgExpr
		arg.Expr = e
		arg.Loc = ast.Loc{Start: start, End: e.Span().End}
	}
	return arg, nil
}

// atTableArgModifier reports whether the cursor sits at a table-argument
// modifier in position — PARTITION/ORDER followed by BY, or PRUNE/KEEP followed
// by WHEN. These four words are non-reserved keywords (usable as relation
// aliases), so the table-argument parser must not consume them as an alias when
// they actually introduce a modifier clause.
func (p *Parser) atTableArgModifier() bool {
	switch p.cur.Kind {
	case kwPARTITION, kwORDER:
		return p.peekNext().Kind == kwBY
	case kwPRUNE, kwKEEP:
		return p.peekNext().Kind == kwWHEN
	default:
		return false
	}
}

// parseTableArgument parses a `tableArgument`: a `TABLE ( name | query )`
// relation with the optional PARTITION BY / PRUNE|KEEP WHEN EMPTY / ORDER BY
// modifiers (TABLE is current).
func (p *Parser) parseTableArgument() (*TableArgument, error) {
	tableTok := p.advance() // consume TABLE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	ta := &TableArgument{Loc: ast.Loc{Start: tableTok.Loc.Start}}

	// TABLE ( query ) vs TABLE ( qualifiedName ). A query starts with
	// SELECT/WITH/VALUES/TABLE-name; a bare table name is a qualifiedName. A
	// leading '(' here is unambiguously a (parenthesized) query — a qualifiedName
	// cannot start with '(' — so it is parsed as a query, not via the relation
	// speculation parseParenRelationOrSubquery uses (a tableArgument's content is
	// `query | qualifiedName`, never a parenthesized relation).
	if p.startsQueryRelation() || p.cur.Kind == int('(') {
		q, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		ta.Query = q
	} else {
		name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		ta.Name = name
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ta.Loc.End = closeTok.Loc.End

	// Optional `AS? identifier columnAliases?`. The alias is NOT taken when the
	// current token is a table-argument MODIFIER in position — PARTITION/ORDER
	// followed by BY, or PRUNE/KEEP followed by WHEN — because those non-reserved
	// keywords would otherwise be mis-consumed as the alias (e.g. the alias parser
	// would grab PARTITION, leaving a dangling BY). An explicit AS still forces an
	// alias. Probed against Trino 481 (`TABLE(orders) PARTITION BY …`).
	if !p.atTableArgModifier() {
		if alias, cols, end, ok, err := p.tryParseRelationAlias(); err != nil {
			return nil, err
		} else if ok {
			ta.Alias = alias
			ta.ColumnAliases = cols
			ta.Loc.End = end
		}
	}

	// Optional PARTITION BY ( … ) | expression.
	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, end, err := p.parseParenOrBareExprList()
		if err != nil {
			return nil, err
		}
		ta.PartitionBy = exprs
		ta.Loc.End = end
	}

	// Optional PRUNE WHEN EMPTY | KEEP WHEN EMPTY.
	if p.cur.Kind == kwPRUNE || p.cur.Kind == kwKEEP {
		kwTok := p.advance() // consume PRUNE / KEEP
		if _, err := p.expect(kwWHEN); err != nil {
			return nil, err
		}
		emptyTok, err := p.expect(kwEMPTY)
		if err != nil {
			return nil, err
		}
		ta.Pruning = kwTok.Str + " WHEN EMPTY"
		ta.Loc.End = emptyTok.Loc.End
	}

	// Optional ORDER BY ( sortItem, … ) | sortItem.
	if p.cur.Kind == kwORDER {
		items, end, err := p.parseParenOrBareSortItems()
		if err != nil {
			return nil, err
		}
		ta.OrderBy = items
		ta.Loc.End = end
	}

	return ta, nil
}

// parseParenOrBareExprList parses `( expression (, …)* )` (possibly empty) or a
// single bare expression, returning the expressions and the end offset. Used by
// a table argument's PARTITION BY, which allows either spelling.
func (p *Parser) parseParenOrBareExprList() ([]Expr, int, error) {
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		exprs, closeTok, err := p.parseBracketedExprList(int(')'))
		if err != nil {
			return nil, 0, err
		}
		return exprs, closeTok.Loc.End, nil
	}
	e, err := p.parseExpr()
	if err != nil {
		return nil, 0, err
	}
	return []Expr{e}, e.Span().End, nil
}

// parseParenOrBareSortItems parses `ORDER BY ( sortItem (, …)* )` or `ORDER BY
// sortItem` (ORDER is current), returning the sort items and the end offset.
// Used by a table argument's ORDER BY, which allows either spelling.
func (p *Parser) parseParenOrBareSortItems() ([]SortItem, int, error) {
	p.advance() // consume ORDER
	if _, err := p.expect(kwBY); err != nil {
		return nil, 0, err
	}
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		first, err := p.parseSortItem()
		if err != nil {
			return nil, 0, err
		}
		items := []SortItem{first}
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseSortItem()
			if err != nil {
				return nil, 0, err
			}
			items = append(items, next)
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, 0, err
		}
		return items, closeTok.Loc.End, nil
	}
	first, err := p.parseSortItem()
	if err != nil {
		return nil, 0, err
	}
	items := []SortItem{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSortItem()
		if err != nil {
			return nil, 0, err
		}
		items = append(items, next)
	}
	return items, items[len(items)-1].Loc.End, nil
}

// tryParseDescriptorArg parses a `descriptorArgument`:
// `DESCRIPTOR ( field, … )` or `CAST ( NULL AS DESCRIPTOR )`. start is the
// argument's first-token offset. It returns ok=false (parser rewound) when a
// leading CAST is an ordinary CAST expression rather than the
// CAST(NULL AS DESCRIPTOR) form.
func (p *Parser) tryParseDescriptorArg(start int) (*DescriptorArg, bool, error) {
	if p.cur.Kind == kwDESCRIPTOR {
		p.advance() // consume DESCRIPTOR
		if _, err := p.expect(int('(')); err != nil {
			return nil, false, err
		}
		first, err := p.parseDescriptorField()
		if err != nil {
			return nil, false, err
		}
		fields := []DescriptorField{first}
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseDescriptorField()
			if err != nil {
				return nil, false, err
			}
			fields = append(fields, next)
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, false, err
		}
		return &DescriptorArg{Fields: fields, Loc: ast.Loc{Start: start, End: closeTok.Loc.End}}, true, nil
	}

	// CAST ( NULL AS DESCRIPTOR ) — speculatively match the exact shape; on any
	// mismatch rewind so the caller parses an ordinary CAST expression.
	cp := p.checkpoint()
	p.advance() // consume CAST
	if p.cur.Kind != int('(') {
		p.restore(cp)
		return nil, false, nil
	}
	p.advance() // consume '('
	if p.cur.Kind != kwNULL {
		p.restore(cp)
		return nil, false, nil
	}
	p.advance() // consume NULL
	if p.cur.Kind != kwAS {
		p.restore(cp)
		return nil, false, nil
	}
	p.advance() // consume AS
	if p.cur.Kind != kwDESCRIPTOR {
		p.restore(cp)
		return nil, false, nil
	}
	p.advance() // consume DESCRIPTOR
	if p.cur.Kind != int(')') {
		p.restore(cp)
		return nil, false, nil
	}
	closeTok := p.advance() // consume ')'
	return &DescriptorArg{NullCast: true, Loc: ast.Loc{Start: start, End: closeTok.Loc.End}}, true, nil
}

// parseDescriptorField parses one `descriptorField`: `identifier type?`. The
// type is parsed only when a type-start token follows the name; a bare
// identifier (followed by ',' or ')') is a type-less field.
func (p *Parser) parseDescriptorField() (DescriptorField, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return DescriptorField{}, err
	}
	field := DescriptorField{Name: name, Loc: name.Loc}
	// A type follows unless the field ends here (at ',' or ')').
	if p.cur.Kind != int(',') && p.cur.Kind != int(')') {
		ty, err := p.parseType()
		if err != nil {
			return DescriptorField{}, err
		}
		field.Type = ty
		field.Loc.End = ty.Loc.End
	}
	return field, nil
}

// parseCopartitionTables parses one `copartitionTables` group:
// `( qualifiedName , qualifiedName (, …)* )` (the COPARTITION keyword already
// consumed). Per the grammar a group has AT LEAST TWO tables — the first comma
// is mandatory — so a single-table group `COPARTITION (a)` is a syntax error
// (oracle-confirmed: Trino 481 rejects it).
func (p *Parser) parseCopartitionTables() ([]*ast.QualifiedName, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	first, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	// The second table (and its comma) is required.
	if _, err := p.expect(int(',')); err != nil {
		return nil, err
	}
	second, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	tables := []*ast.QualifiedName{first, second}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		tables = append(tables, next)
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return tables, nil
}
