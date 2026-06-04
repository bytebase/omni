package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the `parser-utility` DAG node (with session.go and
// explain_call.go): it implements Trino's SHOW family and the DESCRIBE / DESC
// aliases — the legacy TrinoParser.g4 `statement` alternatives labelled
// showCreateTable / showCreateSchema / showCreateView /
// showCreateMaterializedView / showTables / showSchemas / showCatalogs /
// showColumns / showStats / showStatsForQuery / showRoles / showRoleGrants /
// showFunctions / showSession / showGrants, plus DESCRIBE/DESC (which the
// legacy grammar folds into the showColumns label).
//
// Like expr.go's Expr and datatypes.go's DataType, the statement nodes are
// PARSER-PACKAGE types carrying their own ast.Loc; they satisfy ast.Node by
// returning a NodeTag declared in trino/ast/nodetags.go (the ast tag set is
// closed to the ast-core node — File/Identifier/QualifiedName — so a statement
// that embeds a parser-local child cannot itself live in package ast without
// an import cycle). parseStmt (parser.go) dispatches SHOW / DESCRIBE / DESC
// here.
//
// The legacy SHOW/DESCRIBE alternatives (TrinoParser.g4), folded into ShowStmt:
//
//	SHOW GRANTS (ON TABLE? qualifiedName)?                                  # showGrants
//	SHOW CREATE TABLE qualifiedName                                         # showCreateTable
//	SHOW CREATE SCHEMA qualifiedName                                        # showCreateSchema
//	SHOW CREATE VIEW qualifiedName                                          # showCreateView
//	SHOW CREATE MATERIALIZED VIEW qualifiedName                            # showCreateMaterializedView
//	SHOW TABLES ((FROM|IN) qualifiedName)? (LIKE pattern (ESCAPE escape)?)? # showTables
//	SHOW SCHEMAS ((FROM|IN) identifier)? (LIKE pattern (ESCAPE escape)?)?   # showSchemas
//	SHOW CATALOGS (LIKE pattern (ESCAPE escape)?)?                          # showCatalogs
//	SHOW COLUMNS (FROM|IN) qualifiedName? (LIKE pattern (ESCAPE escape)?)?  # showColumns
//	SHOW STATS FOR qualifiedName                                            # showStats
//	SHOW STATS FOR ( query )                                                # showStatsForQuery
//	SHOW CURRENT? ROLES ((FROM|IN) identifier)?                            # showRoles
//	SHOW ROLE GRANTS ((FROM|IN) identifier)?                              # showRoleGrants
//	DESCRIBE qualifiedName | DESC qualifiedName                            # showColumns
//	SHOW FUNCTIONS ((FROM|IN) qualifiedName)? (LIKE pattern (ESCAPE escape)?)? # showFunctions
//	SHOW SESSION (LIKE pattern (ESCAPE escape)?)?                          # showSession
//
// Adjudicated against the live Trino 481 oracle, not the literal legacy
// grammar. Oracle-confirmed divergences from a naive reading of the legacy
// rule, recorded in the migration divergence ledger:
//
//	U1 (SHOW CREATE FUNCTION). The legacy grammar has no SHOW CREATE FUNCTION
//	   alternative, but Trino 481 accepts `SHOW CREATE FUNCTION name` (oracle:
//	   NOT_SUPPORTED, a semantic verdict — so syntactically ACCEPTED). It is a
//	   documented form (trino.io/docs/current/sql/show-create-function). Added
//	   as a ShowCreateFunction kind; flagged docs-ahead-of-legacy.
//
//	U2 (SHOW COLUMNS requires a name). The legacy `SHOW COLUMNS (FROM|IN)
//	   qualifiedName?` marks the name optional, but Trino 481 rejects
//	   `SHOW COLUMNS FROM` (SYNTAX_ERROR): the name is mandatory after FROM/IN.
//	   parseShowColumns therefore requires it; flagged (legacy grammar wrong).
//
//	U3 (DESCRIBE INPUT/OUTPUT belong to the prepared-statement node). DESCRIBE
//	   INPUT name and DESCRIBE OUTPUT name (describeInput/describeOutput) are
//	   prepared-statement introspection, implemented by the parser-dcl-tcl node
//	   (prepared.go). parseShowOrDescribe dispatches them to that node's
//	   parseDescribeInputOutput when present, else routes to the unsupported
//	   stub, so this node never mis-claims them as DESCRIBE-table.

// ---------------------------------------------------------------------------
// ShowStmt
// ---------------------------------------------------------------------------

// ShowKind classifies a ShowStmt — one per legacy SHOW/DESCRIBE alternative.
type ShowKind int

const (
	// ShowInvalid is the zero value; never produced by the parser.
	ShowInvalid ShowKind = iota
	// ShowGrants is SHOW GRANTS [ON [TABLE] name].
	ShowGrants
	// ShowCreateTable is SHOW CREATE TABLE name.
	ShowCreateTable
	// ShowCreateSchema is SHOW CREATE SCHEMA name.
	ShowCreateSchema
	// ShowCreateView is SHOW CREATE VIEW name.
	ShowCreateView
	// ShowCreateMaterializedView is SHOW CREATE MATERIALIZED VIEW name.
	ShowCreateMaterializedView
	// ShowCreateFunction is SHOW CREATE FUNCTION name (U1: docs-ahead-of-legacy).
	ShowCreateFunction
	// ShowTables is SHOW TABLES [(FROM|IN) name] [LIKE pattern [ESCAPE esc]].
	ShowTables
	// ShowSchemas is SHOW SCHEMAS [(FROM|IN) catalog] [LIKE pattern [ESCAPE esc]].
	ShowSchemas
	// ShowCatalogs is SHOW CATALOGS [LIKE pattern [ESCAPE esc]].
	ShowCatalogs
	// ShowColumns is SHOW COLUMNS (FROM|IN) name [LIKE pattern [ESCAPE esc]]
	// — also the target of DESCRIBE/DESC name.
	ShowColumns
	// ShowStats is SHOW STATS FOR name.
	ShowStats
	// ShowStatsForQuery is SHOW STATS FOR ( query ).
	ShowStatsForQuery
	// ShowRoles is SHOW [CURRENT] ROLES [(FROM|IN) catalog].
	ShowRoles
	// ShowRoleGrants is SHOW ROLE GRANTS [(FROM|IN) catalog].
	ShowRoleGrants
	// ShowFunctions is SHOW FUNCTIONS [(FROM|IN) name] [LIKE pattern [ESCAPE esc]].
	ShowFunctions
	// ShowSession is SHOW SESSION [LIKE pattern [ESCAPE esc]].
	ShowSession
)

// ShowStmt is the parse node for the SHOW family and DESCRIBE/DESC. The fields
// populated depend on Kind; unused fields are zero. It carries enough structure
// for downstream consumers (query-type classification treats SHOW … as a
// read / SelectInfoSchema-style statement; completion offers object names) while
// keeping a single node type for the whole family (mirroring snowflake's
// open-ended SHOW node).
type ShowStmt struct {
	Kind ShowKind

	// Name is the object/target qualifiedName, where the form has one
	// (SHOW CREATE …, SHOW STATS FOR name, SHOW COLUMNS (FROM|IN) name,
	// DESCRIBE/DESC name). Nil otherwise.
	Name *ast.QualifiedName
	// In is the (FROM|IN) scope. For SHOW SCHEMAS / SHOW ROLES /
	// SHOW ROLE GRANTS it is a single-identifier catalog (legacy: identifier);
	// for SHOW TABLES / SHOW FUNCTIONS it is a qualifiedName. Stored as a
	// QualifiedName in both cases (single-part for the identifier forms). Nil
	// when absent.
	In *ast.QualifiedName
	// InKeyword records whether the scope was introduced by IN (true) or FROM
	// (false); meaningless when In is nil.
	InKeyword bool

	// Current is true for SHOW CURRENT ROLES.
	Current bool
	// OnTable is true for SHOW GRANTS ON TABLE name (vs. SHOW GRANTS ON name).
	OnTable bool

	// Like is the LIKE pattern string literal text (decoded, quotes stripped),
	// with HasLike marking its presence (the pattern may legitimately be "").
	Like    string
	HasLike bool
	// Escape is the ESCAPE string literal text; HasEscape marks its presence.
	Escape    string
	HasEscape bool

	// QueryText is the raw inner SQL of SHOW STATS FOR ( query ), trimmed.
	// Only set for ShowStatsForQuery. The query is captured as raw text (a
	// placeholder) rather than a parsed tree — the `query` rule belongs to the
	// parser-select node — mirroring expr.go's SubqueryExpr (B1).
	QueryText string

	Loc ast.Loc
}

// Tag implements ast.Node.
func (s *ShowStmt) Tag() ast.NodeTag { return ast.T_ShowStmt }

// Span returns the source byte range covered by the statement.
func (s *ShowStmt) Span() ast.Loc { return s.Loc }

// Compile-time assertion that *ShowStmt satisfies ast.Node.
var _ ast.Node = (*ShowStmt)(nil)

// ---------------------------------------------------------------------------
// parsing
// ---------------------------------------------------------------------------

// parseShowStmt parses a statement whose leading keyword is SHOW, dispatching
// on the second keyword. The SHOW token has NOT been consumed yet.
func (p *Parser) parseShowStmt() (ast.Node, error) {
	showTok := p.advance() // consume SHOW
	start := showTok.Loc.Start

	switch p.cur.Kind {
	case kwGRANTS:
		return p.parseShowGrants(start)
	case kwCREATE:
		return p.parseShowCreate(start)
	case kwTABLES:
		return p.parseShowTablesLike(start, ShowTables)
	case kwSCHEMAS:
		return p.parseShowSchemas(start)
	case kwCATALOGS:
		return p.parseShowCatalogs(start)
	case kwCOLUMNS:
		return p.parseShowColumns(start)
	case kwSTATS:
		return p.parseShowStats(start)
	case kwCURRENT, kwROLES:
		return p.parseShowRoles(start)
	case kwROLE:
		return p.parseShowRoleGrants(start)
	case kwFUNCTIONS:
		return p.parseShowTablesLike(start, ShowFunctions)
	case kwSESSION:
		return p.parseShowSession(start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseShowGrants parses SHOW GRANTS (ON TABLE? qualifiedName)?. GRANTS is the
// current token.
func (p *Parser) parseShowGrants(start int) (ast.Node, error) {
	grantsTok := p.advance() // consume GRANTS
	s := &ShowStmt{Kind: ShowGrants, Loc: ast.Loc{Start: start, End: grantsTok.Loc.End}}
	if _, ok := p.match(kwON); ok {
		if _, ok := p.match(kwTABLE); ok {
			s.OnTable = true
		}
		name, err := p.parseBoundedQualifiedName(3, "table name")
		if err != nil {
			return nil, err
		}
		s.Name = name
		s.Loc.End = name.Loc.End
	}
	return s, nil
}

// parseShowCreate parses SHOW CREATE (TABLE|SCHEMA|VIEW|MATERIALIZED VIEW|
// FUNCTION) qualifiedName. CREATE is the current token.
func (p *Parser) parseShowCreate(start int) (ast.Node, error) {
	p.advance() // consume CREATE
	var kind ShowKind
	switch p.cur.Kind {
	case kwTABLE:
		p.advance()
		kind = ShowCreateTable
	case kwSCHEMA:
		p.advance()
		kind = ShowCreateSchema
	case kwVIEW:
		p.advance()
		kind = ShowCreateView
	case kwMATERIALIZED:
		p.advance()
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		kind = ShowCreateMaterializedView
	case kwFUNCTION:
		// U1: docs-ahead-of-legacy; Trino 481 accepts SHOW CREATE FUNCTION.
		p.advance()
		kind = ShowCreateFunction
	default:
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.parseBoundedQualifiedName(3, "object name")
	if err != nil {
		return nil, err
	}
	return &ShowStmt{Kind: kind, Name: name, Loc: ast.Loc{Start: start, End: name.Loc.End}}, nil
}

// parseShowTablesLike parses the shape shared by SHOW TABLES and SHOW FUNCTIONS:
// `((FROM|IN) qualifiedName)? (LIKE pattern (ESCAPE escape)?)?`. The
// TABLES/FUNCTIONS keyword is the current token.
func (p *Parser) parseShowTablesLike(start int, kind ShowKind) (ast.Node, error) {
	kwTok := p.advance() // consume TABLES / FUNCTIONS
	s := &ShowStmt{Kind: kind, Loc: ast.Loc{Start: start, End: kwTok.Loc.End}}
	if scopeTok, ok := p.match(kwFROM, kwIN); ok {
		s.InKeyword = scopeTok.Kind == kwIN
		// The FROM/IN scope is a schema (catalog.schema), capped at 2 parts;
		// `SHOW TABLES FROM a.b.c` is a SYNTAX_ERROR in Trino 481.
		name, err := p.parseBoundedQualifiedName(2, "schema name")
		if err != nil {
			return nil, err
		}
		s.In = name
		s.Loc.End = name.Loc.End
	}
	if err := p.parseLikeEscape(s); err != nil {
		return nil, err
	}
	return s, nil
}

// parseShowSchemas parses SHOW SCHEMAS ((FROM|IN) identifier)?
// (LIKE pattern (ESCAPE escape)?)?. The legacy grammar takes a single
// identifier (catalog) after FROM/IN — NOT a qualifiedName — and Trino 481
// rejects a dotted scope (`SHOW SCHEMAS FROM c.x` is a SYNTAX_ERROR).
// SCHEMAS is the current token.
func (p *Parser) parseShowSchemas(start int) (ast.Node, error) {
	schemasTok := p.advance() // consume SCHEMAS
	s := &ShowStmt{Kind: ShowSchemas, Loc: ast.Loc{Start: start, End: schemasTok.Loc.End}}
	if scopeTok, ok := p.match(kwFROM, kwIN); ok {
		s.InKeyword = scopeTok.Kind == kwIN
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.In = &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc}
		s.Loc.End = ident.Loc.End
	}
	if err := p.parseLikeEscape(s); err != nil {
		return nil, err
	}
	return s, nil
}

// parseShowCatalogs parses SHOW CATALOGS (LIKE pattern (ESCAPE escape)?)?.
// CATALOGS is the current token.
func (p *Parser) parseShowCatalogs(start int) (ast.Node, error) {
	catalogsTok := p.advance() // consume CATALOGS
	s := &ShowStmt{Kind: ShowCatalogs, Loc: ast.Loc{Start: start, End: catalogsTok.Loc.End}}
	if err := p.parseLikeEscape(s); err != nil {
		return nil, err
	}
	return s, nil
}

// parseShowColumns parses SHOW COLUMNS (FROM|IN) qualifiedName
// (LIKE pattern (ESCAPE escape)?)?. COLUMNS is the current token.
//
// U2: the (FROM|IN) keyword AND the name are both mandatory — Trino 481 rejects
// both `SHOW COLUMNS` and `SHOW COLUMNS FROM` despite the legacy grammar marking
// the name optional.
func (p *Parser) parseShowColumns(start int) (ast.Node, error) {
	p.advance() // consume COLUMNS
	scopeTok, ok := p.match(kwFROM, kwIN)
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	name, err := p.parseBoundedQualifiedName(3, "table name")
	if err != nil {
		return nil, err
	}
	s := &ShowStmt{
		Kind:      ShowColumns,
		Name:      name,
		InKeyword: scopeTok.Kind == kwIN,
		Loc:       ast.Loc{Start: start, End: name.Loc.End},
	}
	if err := p.parseLikeEscape(s); err != nil {
		return nil, err
	}
	return s, nil
}

// parseShowStats parses SHOW STATS FOR (qualifiedName | ( query )). STATS is the
// current token.
func (p *Parser) parseShowStats(start int) (ast.Node, error) {
	p.advance() // consume STATS
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	// SHOW STATS FOR ( query ) — a parenthesised query. The inner must START a
	// query (SELECT/WITH/TABLE/VALUES); Trino 481 rejects `SHOW STATS FOR
	// (nation)`, `SHOW STATS FOR ()`, and other non-query parens with a
	// SYNTAX_ERROR. We gate on startsQuery() (the same check parser-select's
	// query rule begins with) and then capture the inner query as raw text — the
	// `query` grammar belongs to the parser-select node, so the inner is a
	// placeholder here, mirroring expr.go's SubqueryExpr (B1). Full inner-query
	// validation arrives with parser-select; the gate already rejects the common
	// non-query cases.
	if p.cur.Kind == int('(') {
		openTok := p.advance() // consume '('
		if !p.startsQuery() {
			return nil, p.syntaxErrorAtCur()
		}
		raw, closeTok, err := p.captureBalancedParen(openTok)
		if err != nil {
			return nil, err
		}
		return &ShowStmt{
			Kind:      ShowStatsForQuery,
			QueryText: raw,
			Loc:       ast.Loc{Start: start, End: closeTok.Loc.End},
		}, nil
	}
	name, err := p.parseBoundedQualifiedName(3, "table name")
	if err != nil {
		return nil, err
	}
	return &ShowStmt{Kind: ShowStats, Name: name, Loc: ast.Loc{Start: start, End: name.Loc.End}}, nil
}

// parseShowRoles parses SHOW CURRENT? ROLES ((FROM|IN) identifier)?. The current
// token is CURRENT or ROLES.
func (p *Parser) parseShowRoles(start int) (ast.Node, error) {
	s := &ShowStmt{Kind: ShowRoles}
	if _, ok := p.match(kwCURRENT); ok {
		s.Current = true
	}
	rolesTok, err := p.expect(kwROLES)
	if err != nil {
		return nil, err
	}
	s.Loc = ast.Loc{Start: start, End: rolesTok.Loc.End}
	if scopeTok, ok := p.match(kwFROM, kwIN); ok {
		s.InKeyword = scopeTok.Kind == kwIN
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.In = &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc}
		s.Loc.End = ident.Loc.End
	}
	return s, nil
}

// parseShowRoleGrants parses SHOW ROLE GRANTS ((FROM|IN) identifier)?. ROLE is
// the current token.
func (p *Parser) parseShowRoleGrants(start int) (ast.Node, error) {
	p.advance() // consume ROLE
	grantsTok, err := p.expect(kwGRANTS)
	if err != nil {
		return nil, err
	}
	s := &ShowStmt{Kind: ShowRoleGrants, Loc: ast.Loc{Start: start, End: grantsTok.Loc.End}}
	if scopeTok, ok := p.match(kwFROM, kwIN); ok {
		s.InKeyword = scopeTok.Kind == kwIN
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.In = &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc}
		s.Loc.End = ident.Loc.End
	}
	return s, nil
}

// parseShowSession parses SHOW SESSION (LIKE pattern (ESCAPE escape)?)?. SESSION
// is the current token.
func (p *Parser) parseShowSession(start int) (ast.Node, error) {
	sessionTok := p.advance() // consume SESSION
	s := &ShowStmt{Kind: ShowSession, Loc: ast.Loc{Start: start, End: sessionTok.Loc.End}}
	if err := p.parseLikeEscape(s); err != nil {
		return nil, err
	}
	return s, nil
}

// parseShowOrDescribe is the entry point for DESCRIBE / DESC. The DESCRIBE/DESC
// token is the current token. DESCRIBE qualifiedName and DESC qualifiedName are
// aliases for SHOW COLUMNS (ShowColumns kind). DESCRIBE INPUT name and
// DESCRIBE OUTPUT name are prepared-statement introspection owned by the
// parser-dcl-tcl node (prepared.go); they are routed there when present, else
// to the unsupported stub (U3).
func (p *Parser) parseDescribeStmt() (ast.Node, error) {
	descTok := p.advance() // consume DESCRIBE / DESC

	// DESCRIBE INPUT name / DESCRIBE OUTPUT name / DESCRIBE OUTPUT (query):
	// prepared-statement introspection (legacy describeInput/describeOutput),
	// owned by the parser-dcl-tcl node (prepared.go). The disambiguation is
	// exactly Trino 481's (oracle-confirmed):
	//
	//   - only the DESCRIBE spelling has the prepared form; `DESC INPUT a` is a
	//     SYNTAX_ERROR (DESC takes a single qualifiedName, so two bare names is
	//     invalid).
	//   - INPUT/OUTPUT followed by an identifier-start token is the prepared
	//     `DESCRIBE (INPUT|OUTPUT) statement_name` form. `DESCRIBE INPUT` (INPUT
	//     then EOF) and `DESCRIBE input.t` (INPUT then '.') are instead
	//     DESCRIBE-of-table-named-INPUT (ShowColumns), since INPUT/OUTPUT are
	//     non-reserved keywords usable as identifiers.
	//   - OUTPUT (but NOT INPUT) also has the direct `DESCRIBE OUTPUT (query)`
	//     form; `DESCRIBE INPUT (SELECT 1)` and `DESC OUTPUT (SELECT 1)` are
	//     SYNTAX_ERRORs, so only `DESCRIBE OUTPUT (` routes here on the '('.
	//
	// The prepared forms route to the "not yet supported" stub until
	// parser-dcl-tcl lands (U3); Diagnose then distinguishes valid-Trino-not-
	// yet-implemented from invalid syntax. They are out of this node's scope and
	// excluded from its oracle differential.
	if descTok.Kind == kwDESCRIBE {
		switch {
		case (p.cur.Kind == kwINPUT || p.cur.Kind == kwOUTPUT) && isIdentifierStart(p.peekNext().Kind):
			return p.unsupported("DESCRIBE INPUT/OUTPUT")
		case p.cur.Kind == kwOUTPUT && p.peekNext().Kind == int('('):
			return p.unsupported("DESCRIBE OUTPUT (query)")
		}
	}

	start := descTok.Loc.Start
	name, err := p.parseBoundedQualifiedName(3, "table name")
	if err != nil {
		return nil, err
	}
	return &ShowStmt{Kind: ShowColumns, Name: name, Loc: ast.Loc{Start: start, End: name.Loc.End}}, nil
}

// parseLikeEscape parses an optional `LIKE pattern (ESCAPE escape)?` tail into s.
// Both pattern and escape are string literals. Sets HasLike / HasEscape so an
// empty pattern is distinguishable from no LIKE clause.
func (p *Parser) parseLikeEscape(s *ShowStmt) error {
	if _, ok := p.match(kwLIKE); !ok {
		return nil
	}
	pat, err := p.expectStringLiteral("LIKE pattern")
	if err != nil {
		return err
	}
	s.Like = pat.Str
	s.HasLike = true
	s.Loc.End = pat.Loc.End
	if _, ok := p.match(kwESCAPE); ok {
		esc, err := p.expectStringLiteral("ESCAPE")
		if err != nil {
			return err
		}
		s.Escape = esc.Str
		s.HasEscape = true
		s.Loc.End = esc.Loc.End
	}
	return nil
}

// expectStringLiteral consumes the current token if it is a string literal
// ('...' or U&'...'), returning it. Otherwise returns a *ParseError naming the
// context (e.g. "LIKE pattern").
func (p *Parser) expectStringLiteral(context string) (Token, error) {
	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return Token{}, &ParseError{Loc: p.cur.Loc, Msg: "expected string literal for " + context}
	}
	return p.advance(), nil
}

// parseBoundedQualifiedName parses a qualifiedName and rejects it with a
// *ParseError when it has more than maxParts dotted components. Trino's legacy
// ANTLR grammar uses the unbounded `qualifiedName` for every name position, but
// Trino 481 enforces a part-count limit per position at PARSE time (the verdict
// is SYNTAX_ERROR), so the omni parser must too — otherwise it over-accepts
// names the oracle rejects. Object names (table/view/schema/function/procedure
// targets) cap at 3 (catalog.schema.object); a schema scope (SHOW TABLES/
// FUNCTIONS FROM) caps at 2 (catalog.schema). This is a confirmed
// docs/oracle-vs-legacy divergence (legacy grammar too permissive).
func (p *Parser) parseBoundedQualifiedName(maxParts int, context string) (*ast.QualifiedName, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if len(name.Parts) > maxParts {
		return nil, &ParseError{
			Loc: name.Loc,
			Msg: "too many dotted parts in " + context + " (expected at most " + strconv.Itoa(maxParts) + ")",
		}
	}
	return name, nil
}

// captureBalancedParen consumes tokens up to and including the ')' matching an
// already-consumed '(' (openTok), returning the inner text (trimmed) and the
// closing token. Nested parentheses are balanced. Mirrors
// expr.go's parseSubqueryPlaceholder — used to capture a SHOW STATS FOR ( query )
// inner query as raw text without parsing the query grammar (a parser-select
// concern).
func (p *Parser) captureBalancedParen(openTok Token) (string, Token, error) {
	depth := 1
	innerStart := p.cur.Loc.Start
	innerEnd := p.cur.Loc.Start
	for depth > 0 && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
			if depth == 0 {
				innerEnd = p.cur.Loc.Start
			}
		}
		if depth > 0 {
			p.advance()
		}
	}
	if depth != 0 {
		return "", Token{}, &ParseError{Loc: p.cur.Loc, Msg: "unterminated parenthesised query"}
	}
	raw := p.sourceSlice(innerStart, innerEnd)
	closeTok := p.advance() // consume ')'
	return strings.TrimSpace(raw), closeTok, nil
}
