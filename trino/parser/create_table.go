package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's CREATE
// TABLE statements (the createTable and createTableAsSelect alternatives) plus
// the shared property-list and column-definition helpers the rest of the DDL
// node reuses (ALTER TABLE, CREATE SCHEMA/VIEW/MATERIALIZED VIEW, CREATE
// CATALOG, ANALYZE all parse a `WITH ( property = value, … )` list).
//
// The statement / clause structs live here in package parser (matching the
// Trino convention for parser-built node types — Expr in expr.go, DataType in
// datatypes.go, the DML nodes in insert.go). They satisfy ast.Node via a Tag()
// returning the placeholder ast.T_Invalid because the ast node-tag set is
// closed to the ast-core node (File / Identifier / QualifiedName) plus the
// utility/dcl-tcl statement tags; downstream consumers (analysis, completion,
// deparse) reach the concrete statement type by a Go type switch, exactly as
// they do for QueryStmt / InsertStmt.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : CREATE_ (OR_ REPLACE_)? TABLE_ (IF_ NOT_ EXISTS_)? qualifiedName
//	        columnAliases? (COMMENT_ string_)? (WITH_ properties)?
//	        AS_ (rootQuery | LPAREN_ rootQuery RPAREN_) (WITH_ (NO_)? DATA_)?  # createTableAsSelect
//	    | CREATE_ (OR_ REPLACE_)? TABLE_ (IF_ NOT_ EXISTS_)? qualifiedName
//	        LPAREN_ tableElement (COMMA_ tableElement)* RPAREN_
//	        (COMMENT_ string_)? (WITH_ properties)?                           # createTable
//	    ;
//	tableElement     : columnDefinition | likeClause ;
//	columnDefinition : identifier type (NOT_ NULL_)? (COMMENT_ string_)? (WITH_ properties)? ;
//	likeClause       : LIKE_ qualifiedName ((INCLUDING_ | EXCLUDING_) PROPERTIES_)? ;
//	properties           : LPAREN_ propertyAssignments RPAREN_ ;
//	propertyAssignments  : property (COMMA_ property)* ;
//	property             : identifier EQ_ propertyValue ;
//	propertyValue        : DEFAULT_ | expression ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar (which lags 481). Oracle-confirmed facts baked in:
//
//	D-CT1 (target ≤ catalog.schema.table). `CREATE TABLE a.b.c.d (…)` is a
//	   SYNTAX_ERROR; the target qualifiedName is bounded to 3 parts.
//	D-CT2 (column DEFAULT + FIXED clause order). Trino 481 accepts `col type
//	   DEFAULT expr` in a column definition — the legacy grammar omits DEFAULT
//	   (docs-ahead, ddl.md create-table). The optional column-constraint clauses
//	   appear in a FIXED order: DEFAULT, then NOT NULL, then COMMENT, then WITH,
//	   each at most once. 481 REJECTS out-of-order clauses (`NOT NULL DEFAULT 5`,
//	   `COMMENT 'x' NOT NULL`, `WITH (...) NOT NULL` are all SYNTAX_ERRORs), so
//	   parseColumnDefinition consumes them as a strict sequence, not a free set
//	   — verified case-by-case against the oracle.
//	D-CT6 (OR REPLACE ⊕ IF NOT EXISTS). CREATE TABLE rejects the two together at
//	   PARSE time (SYNTAX_ERROR). This is TABLE-specific: the same combination on
//	   CREATE MATERIALIZED VIEW is NOT_SUPPORTED (semantic, grammar-accepted).
//	D-CT3 (empty body rejected). `CREATE TABLE t ()` is a SYNTAX_ERROR — the
//	   parenthesized form requires at least one tableElement.
//	D-CT4 (CTAS requires AS query). The createTable form (paren column list) and
//	   the createTableAsSelect form (`AS query`) are distinct; a bare
//	   `CREATE TABLE t` with neither is a SYNTAX_ERROR.
//	D-CT5 (CTAS column aliases are bare identifiers). `CREATE TABLE t (a, b) AS …`
//	   parses (a, b) as output-column aliases, NOT column definitions: the
//	   trailing `AS query` disambiguates. A parenthesized list whose entries are
//	   `ident type …` is the column-definition form and must NOT be followed by AS.

// createObjectKeyword reports the object keyword of a CREATE statement — the
// token kind directly after CREATE, or after an optional `OR REPLACE` prefix
// (CREATE OR REPLACE { TABLE | VIEW | MATERIALIZED VIEW }). It does NOT consume
// input (it uses checkpoint/restore for the OR-REPLACE lookahead, which the
// parser's two-token window cannot reach). The CREATE dispatch uses it to route
// to the right per-object DDL parser; ROLE / FUNCTION are handled before this
// is reached. On entry cur is the CREATE keyword.
func (p *Parser) createObjectKeyword() TokenKind {
	if p.peekNext().Kind != kwOR {
		return p.peekNext().Kind
	}
	cp := p.checkpoint()
	defer p.restore(cp)
	p.advance() // consume CREATE
	p.advance() // consume OR (already confirmed)
	if p.cur.Kind != kwREPLACE {
		return tokInvalid
	}
	p.advance() // consume REPLACE
	return p.cur.Kind
}

// ---------------------------------------------------------------------------
// Property list (shared across the DDL node)
// ---------------------------------------------------------------------------

// Property is one `identifier = propertyValue` assignment inside a
// `WITH ( … )` list (or an ALTER … SET PROPERTIES list). Value is nil when the
// property uses the `= DEFAULT` reset form (IsDefault is then true).
type Property struct {
	Name      *ast.Identifier
	Value     Expr // nil iff IsDefault
	IsDefault bool // the `name = DEFAULT` reset form
	Loc       ast.Loc
}

// parsePropertyList parses a `( property (, property)* )` list. On entry cur is
// the opening '('. The list must be non-empty (Trino rejects `WITH ()`). Each
// property is `identifier = ( DEFAULT | expression )`.
func (p *Parser) parsePropertyList() ([]*Property, ast.Loc, error) {
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, ast.Loc{}, err
	}
	first, err := p.parseProperty()
	if err != nil {
		return nil, ast.Loc{}, err
	}
	props := []*Property{first}
	for {
		if _, ok := p.match(int(',')); !ok {
			break
		}
		next, err := p.parseProperty()
		if err != nil {
			return nil, ast.Loc{}, err
		}
		props = append(props, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.Loc{}, err
	}
	return props, ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}, nil
}

// parseProperty parses one `identifier = ( DEFAULT | expression )`.
func (p *Parser) parseProperty() (*Property, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	// propertyValue: DEFAULT keyword (reset) or an expression. DEFAULT is a
	// reserved keyword, so a bare `= DEFAULT` cannot be a column reference.
	if defTok, ok := p.match(kwDEFAULT); ok {
		return &Property{
			Name:      name,
			IsDefault: true,
			Loc:       ast.Loc{Start: name.Loc.Start, End: defTok.Loc.End},
		}, nil
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &Property{
		Name:  name,
		Value: value,
		Loc:   ast.Loc{Start: name.Loc.Start, End: value.Span().End},
	}, nil
}

// parseOptionalWithProperties consumes a `WITH ( property, … )` clause when the
// current token is WITH, returning the property list (nil when WITH is absent).
func (p *Parser) parseOptionalWithProperties() ([]*Property, error) {
	if _, ok := p.match(kwWITH); !ok {
		return nil, nil
	}
	props, _, err := p.parsePropertyList()
	return props, err
}

// ---------------------------------------------------------------------------
// Column definitions / table elements
// ---------------------------------------------------------------------------

// ColumnDefinition is one `identifier type [DEFAULT expr] [NOT NULL]
// [COMMENT str] [WITH (props)]` column of a CREATE TABLE body (columnDefinition
// rule). Each constraint clause is optional and at most once.
type ColumnDefinition struct {
	Name       *ast.Identifier
	Type       *DataType
	Default    Expr // DEFAULT expr (D-CT2 docs-ahead extension); nil when absent
	NotNull    bool
	Comment    *string // COMMENT string literal value (decoded); nil when absent
	Properties []*Property
	Loc        ast.Loc
}

// LikeClause is a `LIKE qualifiedName [{INCLUDING | EXCLUDING} PROPERTIES]`
// table element. Including is true for INCLUDING PROPERTIES, false for
// EXCLUDING PROPERTIES or no option (EXCLUDING is the default).
type LikeClause struct {
	Source    *ast.QualifiedName
	HasOption bool // an explicit INCLUDING/EXCLUDING PROPERTIES was given
	Including bool // INCLUDING PROPERTIES (only meaningful when HasOption)
	Loc       ast.Loc
}

// TableElement is one entry of a CREATE TABLE parenthesized body: either a
// column definition or a LIKE clause. Exactly one of Column / Like is non-nil.
type TableElement struct {
	Column *ColumnDefinition
	Like   *LikeClause
}

// parseTableElement parses one tableElement: a `LIKE …` clause or a column
// definition. The LIKE keyword unambiguously selects the likeClause form;
// anything else is a columnDefinition (which begins with an identifier).
func (p *Parser) parseTableElement() (*TableElement, error) {
	if p.cur.Kind == kwLIKE {
		like, err := p.parseLikeClause()
		if err != nil {
			return nil, err
		}
		return &TableElement{Like: like}, nil
	}
	col, err := p.parseColumnDefinition()
	if err != nil {
		return nil, err
	}
	return &TableElement{Column: col}, nil
}

// parseLikeClause parses `LIKE qualifiedName [{INCLUDING | EXCLUDING}
// PROPERTIES]`. On entry cur is the LIKE keyword.
func (p *Parser) parseLikeClause() (*LikeClause, error) {
	likeTok := p.advance() // consume LIKE
	src, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	lc := &LikeClause{Source: src, Loc: ast.Loc{Start: likeTok.Loc.Start, End: src.Loc.End}}
	if opt, ok := p.match(kwINCLUDING, kwEXCLUDING); ok {
		propsTok, err := p.expect(kwPROPERTIES)
		if err != nil {
			return nil, err
		}
		lc.HasOption = true
		lc.Including = opt.Kind == kwINCLUDING
		lc.Loc.End = propsTok.Loc.End
	}
	return lc, nil
}

// parseColumnDefinition parses `identifier type [DEFAULT expr] [NOT NULL]
// [COMMENT str] [WITH (props)]`. The optional constraint clauses appear in a
// FIXED order — DEFAULT, then NOT NULL, then COMMENT, then WITH — each at most
// once. The order is oracle-pinned (D-CT2): Trino 481 rejects out-of-order
// clauses (`NOT NULL DEFAULT 5`, `COMMENT 'x' NOT NULL`, `WITH (...) NOT NULL`
// are all SYNTAX_ERRORs), so the parser consumes them as a strict sequence, not
// an order-tolerant set. The legacy grammar omits DEFAULT entirely; 481 adds it
// as the first clause (ddl.md create-table).
func (p *Parser) parseColumnDefinition() (*ColumnDefinition, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	colType, err := p.parseType()
	if err != nil {
		return nil, err
	}
	col := &ColumnDefinition{
		Name: name,
		Type: colType,
		Loc:  ast.Loc{Start: name.Loc.Start, End: colType.Loc.End},
	}

	// 1. DEFAULT expr (D-CT2 docs-ahead extension).
	if _, ok := p.match(kwDEFAULT); ok {
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		col.Default = def
		col.Loc.End = def.Span().End
	}
	// 2. NOT NULL.
	if p.cur.Kind == kwNOT {
		p.advance() // consume NOT
		nullTok, err := p.expect(kwNULL)
		if err != nil {
			return nil, err
		}
		col.NotNull = true
		col.Loc.End = nullTok.Loc.End
	}
	// 3. COMMENT str.
	if _, ok := p.match(kwCOMMENT); ok {
		str, err := p.expectStringLiteral("column COMMENT")
		if err != nil {
			return nil, err
		}
		s := str.Str
		col.Comment = &s
		col.Loc.End = str.Loc.End
	}
	// 4. WITH (props).
	if _, ok := p.match(kwWITH); ok {
		props, loc, err := p.parsePropertyList()
		if err != nil {
			return nil, err
		}
		col.Properties = props
		col.Loc.End = loc.End
	}
	return col, nil
}

// ---------------------------------------------------------------------------
// CREATE TABLE
// ---------------------------------------------------------------------------

// CreateTableStmt is a `CREATE [OR REPLACE] TABLE [IF NOT EXISTS] name (…)
// [COMMENT str] [WITH (props)]` statement (the createTable alternative — the
// explicit-column-list form, NOT the AS-query form, which is CreateTableAsStmt).
type CreateTableStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *ast.QualifiedName
	Elements    []*TableElement // at least one (D-CT3)
	Comment     *string         // table COMMENT; nil when absent
	Properties  []*Property
	Loc         ast.Loc
}

func (n *CreateTableStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateTableStmt) Span() ast.Loc    { return n.Loc }

// CreateTableAsStmt is a `CREATE [OR REPLACE] TABLE [IF NOT EXISTS] name
// [(col_alias, …)] [COMMENT str] [WITH (props)] AS query [WITH [NO] DATA]`
// statement (the createTableAsSelect alternative). ColumnAliases renames the
// query's output columns; WithData / NoData record the trailing data clause.
type CreateTableAsStmt struct {
	OrReplace     bool
	IfNotExists   bool
	Name          *ast.QualifiedName
	ColumnAliases []*ast.Identifier // optional output-column aliases (nil when absent)
	Comment       *string
	Properties    []*Property
	Query         *Query
	HasDataClause bool // a trailing WITH [NO] DATA was given
	NoData        bool // WITH NO DATA (only meaningful when HasDataClause)
	Loc           ast.Loc
}

func (n *CreateTableAsStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateTableAsStmt) Span() ast.Loc    { return n.Loc }

var (
	_ ast.Node = (*CreateTableStmt)(nil)
	_ ast.Node = (*CreateTableAsStmt)(nil)
)

// parseCreateTableStmt parses both CREATE TABLE forms. On entry cur is the
// CREATE keyword (not yet consumed); startOffset is CREATE's byte offset.
//
// The createTable (paren column list) and createTableAsSelect (AS query) forms
// share a prefix up to and past the target name; they diverge on whether a '('
// here opens a column-definition list or an output-alias list. The
// disambiguation cannot be made locally (`CREATE TABLE t (a` could be either a
// column `a <type>` or an alias `a`), so the parenthesized list is parsed once
// in a mode that accepts BOTH shapes, and the presence of a trailing AS decides
// which statement form was meant (D-CT5).
func (p *Parser) parseCreateTableStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume CREATE
	orReplace := false
	if _, ok := p.match(kwOR); ok {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		orReplace = true
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	ifNotExistsTok := p.cur
	ifNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	// D-CT6 (oracle-pinned): CREATE TABLE rejects OR REPLACE together with
	// IF NOT EXISTS at PARSE time (SYNTAX_ERROR), matching the docs ("OR REPLACE
	// and IF NOT EXISTS cannot be used together"). Note this is TABLE-specific:
	// the same combination on CREATE MATERIALIZED VIEW is NOT_SUPPORTED
	// (semantic, accepted by the grammar), so that parser does NOT reject it.
	if orReplace && ifNotExists {
		return nil, &ParseError{
			Loc: ifNotExistsTok.Loc,
			Msg: "CREATE TABLE cannot combine OR REPLACE with IF NOT EXISTS",
		}
	}

	// D-CT1: target is at most catalog.schema.table.
	name, err := p.parseBoundedQualifiedName(3, "table name")
	if err != nil {
		return nil, err
	}

	// A '(' here opens either a column-definition list (createTable) or an
	// output-alias list (createTableAsSelect). Parse the contents once,
	// recording enough to build either node, then let the trailing AS decide.
	var (
		elements      []*TableElement
		columnAliases []*ast.Identifier
		hadParen      bool
		aliasOnly     bool // every paren entry was a bare identifier (alias-shaped)
	)
	if p.cur.Kind == int('(') {
		hadParen = true
		elements, columnAliases, aliasOnly, err = p.parseCreateTableParenBody()
		if err != nil {
			return nil, err
		}
	}

	// Optional COMMENT then optional WITH (properties), in that grammar order.
	comment, err := p.parseOptionalComment()
	if err != nil {
		return nil, err
	}
	properties, err := p.parseOptionalWithProperties()
	if err != nil {
		return nil, err
	}

	// The trailing AS decides the form. With AS → createTableAsSelect (the paren
	// list, if any, is column aliases — must be alias-shaped). Without AS →
	// createTable (the paren list is column definitions and is mandatory).
	if p.cur.Kind == kwAS {
		return p.finishCreateTableAs(startOffset, orReplace, ifNotExists, name, columnAliases, aliasOnly, hadParen, comment, properties)
	}

	// createTable: the parenthesized column-definition body is mandatory (D-CT4).
	if !hadParen {
		return nil, p.syntaxErrorAtCur()
	}
	// A list that was purely bare identifiers (no types) is only valid as CTAS
	// aliases, which require AS — without AS it is an incomplete column
	// definition (each column needs a type). Reject to match Trino.
	if aliasOnly {
		return nil, &ParseError{Loc: name.Loc, Msg: "expected column type in CREATE TABLE column definition"}
	}
	stmt := &CreateTableStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
		Elements:    elements,
		Comment:     comment,
		Properties:  properties,
		Loc:         ast.Loc{Start: startOffset, End: p.prev.Loc.End},
	}
	return stmt, nil
}

// finishCreateTableAs completes a createTableAsSelect after the shared prefix.
// On entry cur is the AS keyword. columnAliases holds the optional output-alias
// list (already parsed); aliasOnly reports whether the paren body was purely
// bare identifiers (the only legal CTAS alias shape); hadParen reports whether
// any '(' was present.
func (p *Parser) finishCreateTableAs(startOffset int, orReplace, ifNotExists bool, name *ast.QualifiedName, columnAliases []*ast.Identifier, aliasOnly, hadParen bool, comment *string, properties []*Property) (ast.Node, error) {
	// If a paren list was present it must be alias-shaped (bare identifiers);
	// a column-definition-shaped list (`a bigint`) before AS is a syntax error.
	if hadParen && !aliasOnly {
		return nil, &ParseError{Loc: name.Loc, Msg: "CREATE TABLE AS column list must be plain column aliases"}
	}
	p.advance() // consume AS

	// The query may be parenthesized: `AS (query)`. parseQuery handles the
	// SELECT / VALUES / TABLE / WITH … forms and a leading '(' subquery.
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}

	stmt := &CreateTableAsStmt{
		OrReplace:     orReplace,
		IfNotExists:   ifNotExists,
		Name:          name,
		ColumnAliases: columnAliases,
		Comment:       comment,
		Properties:    properties,
		Query:         query,
		Loc:           ast.Loc{Start: startOffset, End: query.Loc.End},
	}

	// Optional trailing WITH [NO] DATA.
	if _, ok := p.match(kwWITH); ok {
		stmt.HasDataClause = true
		if _, ok := p.match(kwNO); ok {
			stmt.NoData = true
		}
		dataTok, err := p.expect(kwDATA)
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = dataTok.Loc.End
	}
	return stmt, nil
}

// parseCreateTableParenBody parses the `( … )` body shared by the createTable
// and createTableAsSelect forms. It returns BOTH interpretations: elements is
// the column-definition / LIKE list, columnAliases is the alias list (populated
// only when every entry was a bare identifier), and aliasOnly reports the
// alias-shaped case. On entry cur is the opening '('.
//
// Each entry is parsed as a full tableElement (column definition or LIKE). When
// an entry is a column whose definition is exactly a bare identifier with no
// type (only possible as a CTAS alias), it is recorded as an alias instead.
// Because a column definition requires a type, the parser cannot know upfront
// which form it is in; it parses optimistically and tracks whether the list
// stayed alias-shaped throughout.
func (p *Parser) parseCreateTableParenBody() (elements []*TableElement, columnAliases []*ast.Identifier, aliasOnly bool, err error) {
	if _, err = p.expect(int('(')); err != nil {
		return nil, nil, false, err
	}
	aliasOnly = true
	for {
		// Peek: a bare identifier immediately followed by ',' or ')' is an alias
		// (no type), the only CTAS-alias shape. Anything else is a real
		// tableElement (column definition with a type, or LIKE).
		if isIdentifierStart(p.cur.Kind) && p.cur.Kind != kwLIKE &&
			(p.peekNext().Kind == int(',') || p.peekNext().Kind == int(')')) {
			alias := identFromToken(p.advance())
			columnAliases = append(columnAliases, alias)
			// Also record a column element so the createTable error path can
			// report a missing type rather than losing the entry.
			elements = append(elements, &TableElement{Column: &ColumnDefinition{Name: alias, Loc: alias.Loc}})
		} else {
			elem, perr := p.parseTableElement()
			if perr != nil {
				return nil, nil, false, perr
			}
			elements = append(elements, elem)
			aliasOnly = false
		}
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	if _, err = p.expect(int(')')); err != nil {
		return nil, nil, false, err
	}
	if !aliasOnly {
		columnAliases = nil
	}
	return elements, columnAliases, aliasOnly, nil
}

// ---------------------------------------------------------------------------
// Shared small clause helpers (used across the DDL node)
// ---------------------------------------------------------------------------

// parseOptionalIfNotExists consumes an `IF NOT EXISTS` clause when present,
// returning whether it was found. A bare `IF` not followed by `NOT EXISTS` is a
// syntax error (IF is otherwise not a statement-position keyword here).
func (p *Parser) parseOptionalIfNotExists() (bool, error) {
	if _, ok := p.match(kwIF); !ok {
		return false, nil
	}
	if _, err := p.expect(kwNOT); err != nil {
		return false, err
	}
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// parseOptionalIfExists consumes an `IF EXISTS` clause when present, returning
// whether it was found.
func (p *Parser) parseOptionalIfExists() (bool, error) {
	if _, ok := p.match(kwIF); !ok {
		return false, nil
	}
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// parseOptionalComment consumes a `COMMENT string_` clause when present,
// returning the decoded string value (nil when absent).
func (p *Parser) parseOptionalComment() (*string, error) {
	if _, ok := p.match(kwCOMMENT); !ok {
		return nil, nil
	}
	str, err := p.expectStringLiteral("COMMENT")
	if err != nil {
		return nil, err
	}
	s := str.Str
	return &s, nil
}
