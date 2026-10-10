package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// Clauses of subprogram and package headings.
//
// Each clause stands at most once in a heading: Oracle 23ai rejects a second
// AUTHID, ACCESSIBLE BY, DETERMINISTIC, PIPELINED, PARALLEL_ENABLE,
// RESULT_CACHE, SQL_MACRO, ORDER BY, or CLUSTER BY with PLS-00371.

// clauseSeen records the clauses a heading has used, for the PLS-00371 check.
type clauseSeen map[string]bool

// subprogramLevel says where a subprogram heading stands, which decides the
// clauses it takes. On Oracle 23ai AUTHID stands only on a schema-level unit
// (PLS-00157 on a packaged, nested, or object type method), and ACCESSIBLE
// BY only on a schema-level unit or a package subprogram (PLS-00262 on a
// type method, PLS-00262 or PLS-00263 on a nested subprogram). DEFAULT
// COLLATION is documented for schema-level units only.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/DEFAULT-COLLATION-clause.html
type subprogramLevel int

const (
	subprogramSchema   subprogramLevel = iota // CREATE PROCEDURE or FUNCTION
	subprogramPackaged                        // in a package specification or body
	subprogramNested                          // declared in a block or another subprogram
	subprogramMethod                          // an object type method
)

// takesAccessibleBy reports whether a heading at level takes ACCESSIBLE BY.
func (l subprogramLevel) takesAccessibleBy() bool {
	return l == subprogramSchema || l == subprogramPackaged
}

// firstClause reports a syntax error at the current token, which starts clause,
// when the heading already has that clause.
func (p *Parser) firstClause(seen clauseSeen, clause string) error {
	if seen[clause] {
		return p.syntaxErrorAtCur()
	}
	seen[clause] = true
	return nil
}

// parseAccessibleByClause parses ACCESSIBLE BY ( accessor [, accessor ]... ),
// the current token at ACCESSIBLE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/ACCESSIBLE-BY-clause.html
//
//	accessor  ::= [ unit_kind ] [ schema. ] unit_name
//	unit_kind ::= { FUNCTION | PROCEDURE | PACKAGE | TRIGGER | TYPE }
//
// Oracle 23ai reads an unquoted unit kind word as the kind, so ACCESSIBLE BY
// (PACKAGE) is PLS-00103, as are an empty list, a missing parenthesis, and an
// accessor with a database link.
func (p *Parser) parseAccessibleByClause() (*nodes.List, error) {
	p.advance() // consume ACCESSIBLE
	if p.cur.Type != kwBY {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for {
		acc := &nodes.Accessor{Loc: nodes.Loc{Start: p.pos()}}
		switch p.cur.Type {
		case kwFUNCTION:
			acc.UnitKind = "FUNCTION"
		case kwPROCEDURE:
			acc.UnitKind = "PROCEDURE"
		case kwPACKAGE:
			acc.UnitKind = "PACKAGE"
		case kwTRIGGER:
			acc.UnitKind = "TRIGGER"
		case kwTYPE:
			acc.UnitKind = "TYPE"
		}
		if acc.UnitKind != "" {
			p.advance()
		}
		name, err := p.parseUnitName()
		if err != nil {
			return nil, err
		}
		acc.Name = name
		acc.Loc.End = p.prev.End
		list.Items = append(list.Items, acc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseUnitName parses [ schema. ] name, without a database link. Neither
// part may be an unquoted PL/SQL reserved word: Oracle 23ai rejects
// ACCESSIBLE BY (SELECT), (BEGIN), and (sys.select), and RELIES_ON (BEGIN),
// with PLS-00103, and accepts (loop).
func (p *Parser) parseUnitName() (*nodes.ObjectName, error) {
	start := p.pos()
	if !p.isPLSQLIdentifier() {
		return nil, p.syntaxErrorAtCur()
	}
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	name := &nodes.ObjectName{Name: first, Loc: nodes.Loc{Start: start}}
	if p.cur.Type == '.' {
		p.advance()
		if !p.isPLSQLIdentifier() {
			return nil, p.syntaxErrorAtCur()
		}
		second, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		name.Schema, name.Name = first, second
	}
	name.Loc.End = p.prev.End
	return name, nil
}

// parseDefaultCollationClause parses DEFAULT COLLATION USING_NLS_COMP, the
// current token at DEFAULT, and returns the option. USING_NLS_COMP, unquoted,
// is the one option a PL/SQL unit takes.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/DEFAULT-COLLATION-clause.html
func (p *Parser) parseDefaultCollationClause() (string, error) {
	p.advance() // consume DEFAULT
	p.advance() // consume COLLATION
	if !p.isKeywordStr("USING_NLS_COMP") {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return "USING_NLS_COMP", nil
}

// atDefaultCollation reports whether the current tokens are DEFAULT COLLATION.
func (p *Parser) atDefaultCollation() bool {
	return p.cur.Type == kwDEFAULT && p.isIdentLikeStrAt(p.peekNext(), "COLLATION")
}

// parseProcedureProperties parses the clauses between a procedure's
// parameters and its IS | AS or ';', in any order; level limits them as
// subprogramLevel says.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-PROCEDURE-statement.html
//
//	[ { invoker_rights_clause | accessible_by_clause | default_collation_clause }... ]
func (p *Parser) parseProcedureProperties(stmt *nodes.CreateProcedureStmt, level subprogramLevel) error {
	seen := clauseSeen{}
	for {
		switch {
		case p.isKeywordStr("AUTHID") && level == subprogramSchema:
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE") && level.takesAccessibleBy():
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation() && level == subprogramSchema:
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			collation, err := p.parseDefaultCollationClause()
			if err != nil {
				return err
			}
			stmt.DefaultCollation = collation
		default:
			return nil
		}
	}
}

// parsePackageProperties parses the clauses between a package
// specification's name and its IS | AS, in any order.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-PACKAGE-statement.html
//
//	[ { default_collation_clause | invoker_rights_clause | accessible_by_clause }... ]
func (p *Parser) parsePackageProperties(stmt *nodes.CreatePackageStmt) error {
	seen := clauseSeen{}
	for {
		switch {
		case p.isKeywordStr("AUTHID"):
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE"):
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation():
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			collation, err := p.parseDefaultCollationClause()
			if err != nil {
				return err
			}
			stmt.DefaultCollation = collation
		default:
			return nil
		}
	}
}

// parseParallelEnableSpec parses the partitioning of PARALLEL_ENABLE, the
// current token at '('.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/PARALLEL_ENABLE-clause.html
//
//	( PARTITION argument BY { ANY | { HASH | RANGE | VALUE } ( expr [, expr ]... ) } )
//
// Oracle 23ai rejects an empty partitioning, PARTITION without BY, other
// methods (LIST), and an empty column list with PLS-00103. That VALUE takes a
// single column is checked when the unit compiles.
func (p *Parser) parseParallelEnableSpec() (*nodes.ParallelEnableClause, Token, error) {
	spec := &nodes.ParallelEnableClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // consume '('
	if p.cur.Type != kwPARTITION {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isIdentLike() {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	argTok := p.cur
	arg, err := p.parseIdentifier()
	if err != nil {
		return nil, Token{}, err
	}
	spec.Argument = arg
	if p.cur.Type != kwBY {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.cur.Type == kwANY:
		spec.PartitionBy = "ANY"
		p.advance()
	case p.cur.Type == kwHASH, p.cur.Type == kwRANGE, p.isKeywordStr("VALUE"):
		spec.PartitionBy = p.cur.Str
		if p.cur.Type == kwHASH {
			spec.PartitionBy = "HASH"
		} else if p.cur.Type == kwRANGE {
			spec.PartitionBy = "RANGE"
		}
		p.advance()
		cols, err := p.parseParenExprList()
		if err != nil {
			return nil, Token{}, err
		}
		// VALUE partitioning takes one column (PLS-00757); HASH and RANGE
		// take several.
		if spec.PartitionBy == "VALUE" && cols.Len() > 1 {
			return nil, Token{}, p.syntaxErrorAtNode(cols.Items[1])
		}
		if err := p.checkColumnNames(cols); err != nil {
			return nil, Token{}, err
		}
		spec.Columns = cols
	default:
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	p.advance()
	spec.Loc.End = p.prev.End
	return spec, argTok, nil
}

// parseStreamingClause parses { ORDER | CLUSTER } argument BY ( expr [, expr ]... ),
// the current token at ORDER or CLUSTER. Oracle 23ai takes it anywhere in a
// function's clause list, with or without PARALLEL_ENABLE, and requires a
// plain parameter name (ORDER c.x BY is PLS-00103).
func (p *Parser) parseStreamingClause() (*nodes.StreamingClause, Token, error) {
	clause := &nodes.StreamingClause{Kind: "ORDER", Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == kwCLUSTER {
		clause.Kind = "CLUSTER"
	}
	p.advance()
	if !p.isIdentLike() {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	argTok := p.cur
	arg, err := p.parseIdentifier()
	if err != nil {
		return nil, Token{}, err
	}
	clause.Argument = arg
	if p.cur.Type != kwBY {
		return nil, Token{}, p.syntaxErrorAtCur()
	}
	p.advance()
	cols, err := p.parseParenExprList()
	if err != nil {
		return nil, Token{}, err
	}
	if err := p.checkColumnNames(cols); err != nil {
		return nil, Token{}, err
	}
	clause.Columns = cols
	clause.Loc.End = p.prev.End
	return clause, argTok, nil
}

// checkColumnNames checks that each item of a PARALLEL_ENABLE partitioning
// or a streaming clause is a simple column name: Oracle 23ai rejects an
// expression, a qualified name (c.a), or a literal with PLS-00670, and
// accepts a quoted name.
func (p *Parser) checkColumnNames(list *nodes.List) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		ref, ok := item.(*nodes.ColumnRef)
		quoted := ok && ref.Loc.Start >= 0 && ref.Loc.Start < len(p.source) && p.source[ref.Loc.Start] == '"'
		// A quoted "*" names a column (Oracle 23ai compiles HASH ("*"));
		// only the unquoted * is the wildcard.
		if !ok || ref.Table != "" || ref.Schema != "" || (ref.Column == "*" && !quoted) || ref.OuterJoin {
			return p.syntaxErrorAtNode(item)
		}
	}
	return nil
}

// syntaxErrorAtNode returns a syntax error at node n.
func (p *Parser) syntaxErrorAtNode(n nodes.Node) *ParseError {
	loc := nodes.NodeLoc(n)
	return p.syntaxErrorAtTok(Token{Type: tokIDENT, Loc: loc.Start, End: loc.End})
}

// checkParallelArguments checks the arguments PARALLEL_ENABLE partitions and
// the streaming clauses order or cluster, each from the function's own text
// as Oracle 23ai does: the partitioned argument is one of the function's IN
// parameters (PLS-00626, PLS-00625); a streaming clause needs a PARTITION BY
// (PLS-00654, PLS-00665) on the same argument (PLS-00631). partTok and
// streamToks are the argument tokens, for the error position.
func (p *Parser) checkParallelArguments(stmt *nodes.CreateFunctionStmt, partTok Token, streamToks []Token) error {
	if stmt.ParallelSpec == nil {
		if len(streamToks) > 0 {
			return p.syntaxErrorAtTok(streamToks[0])
		}
		return nil
	}
	var param *nodes.Parameter
	if stmt.Parameters != nil {
		for _, item := range stmt.Parameters.Items {
			if pr, ok := item.(*nodes.Parameter); ok && pr.Name == stmt.ParallelSpec.Argument {
				param = pr
				break
			}
		}
	}
	if param == nil || (param.Mode != "" && param.Mode != "IN") {
		return p.syntaxErrorAtTok(partTok)
	}
	// The partitioned argument is a REF CURSOR, strongly typed unless
	// partitioned BY ANY without a streaming clause (PLS-00627): a
	// predefined non-cursor type never qualifies, nor SYS_REFCURSOR with
	// HASH, RANGE, VALUE, ORDER BY, or CLUSTER BY. A named type is left to
	// the engine.
	if name, ok := p.predefinedTypeName(param.TypeName); ok {
		weakAllowed := name == "SYS_REFCURSOR" && stmt.ParallelSpec.PartitionBy == "ANY" && len(streamToks) == 0
		if !weakAllowed {
			return p.syntaxErrorAtTok(partTok)
		}
	}
	if stmt.Streaming != nil {
		for i, item := range stmt.Streaming.Items {
			if item.(*nodes.StreamingClause).Argument != stmt.ParallelSpec.Argument {
				return p.syntaxErrorAtTok(streamToks[i])
			}
		}
	}
	return nil
}

// predefinedTypeName returns the name of tn when it is a predefined datatype
// that is not a REF CURSOR of the program's own: SYS_REFCURSOR, or a scalar,
// character, LOB, or %ROWTYPE record type.
func (p *Parser) predefinedTypeName(tn *nodes.TypeName) (string, bool) {
	if tn == nil {
		return "", false
	}
	if tn.IsPercRowtype {
		return "%ROWTYPE", true
	}
	switch lead := p.predefinedTypeLead(tn); lead {
	case "CHAR", "CHARACTER", "VARCHAR2", "VARCHAR", "STRING", "CLOB":
		return lead, true
	default:
		return lead, plsqlNonCharacterTypes[lead]
	}
}

// checkPolymorphicSignature checks a polymorphic table function's signature
// as Oracle 23ai does from the text: it returns TABLE (PLS-00767) and takes
// exactly one TABLE parameter (PLS-00773, PLS-00766), which has no default
// (PLS-00768), while a COLUMNS parameter defaults to NULL if at all
// (PLS-00769).
func (p *Parser) checkPolymorphicSignature(stmt *nodes.CreateFunctionStmt) error {
	if p.pseudoType(stmt.ReturnType) != "TABLE" {
		return p.syntaxErrorAtType(stmt.ReturnType)
	}
	tables := 0
	if stmt.Parameters != nil {
		for _, item := range stmt.Parameters.Items {
			pr, ok := item.(*nodes.Parameter)
			if !ok {
				continue
			}
			switch p.pseudoType(pr.TypeName) {
			case "TABLE":
				tables++
				if tables > 1 {
					return p.syntaxErrorAtType(pr.TypeName)
				}
				if pr.Default != nil {
					return p.syntaxErrorAtNode(pr.Default)
				}
			case "COLUMNS":
				if _, isNull := pr.Default.(*nodes.NullLiteral); pr.Default != nil && !isNull {
					return p.syntaxErrorAtNode(pr.Default)
				}
			}
		}
	}
	if tables == 0 {
		return p.syntaxErrorAtType(stmt.ReturnType)
	}
	return nil
}

// checkCharsetSources checks each item%CHARSET in a subprogram heading that
// names one of its own parameters: Oracle 23ai rejects one whose parameter
// has a non-character type, such as NUMBER or DATE (PLS-00550). A CHAR,
// VARCHAR2, CLOB, or national character parameter qualifies with or
// without ANY_CS; an item outside the heading is left to the engine.
func (p *Parser) checkCharsetSources(params *nodes.List, result *nodes.TypeName) error {
	if params == nil {
		return nil
	}
	// params holds subprogram parameters (*Parameter) or cursor parameters
	// (*PLSQLVarDecl); Oracle 23ai applies the rule to both.
	formal := func(item nodes.Node) (string, *nodes.TypeName, bool) {
		switch f := item.(type) {
		case *nodes.Parameter:
			return f.Name, f.TypeName, true
		case *nodes.PLSQLVarDecl:
			return f.Name, f.TypeName, true
		}
		return "", nil, false
	}
	byName := map[string]*nodes.TypeName{}
	for _, item := range params.Items {
		if name, tn, ok := formal(item); ok {
			byName[name] = tn
		}
	}
	check := func(tn *nodes.TypeName) error {
		// A qualified source (k.c%CHARSET) names an item outside the
		// heading, even when a quoted formal is spelled "K.C" (Oracle 23ai
		// compiles that function); a single quoted "K.C" names the formal.
		if tn == nil || !tn.IsPercCharset || p.qualifiedCharsets[tn] {
			return nil
		}
		if src, ok := byName[tn.CharacterSet]; ok && !p.mayBeCharacterType(src) {
			return p.syntaxErrorAtType(tn)
		}
		return nil
	}
	for _, item := range params.Items {
		if _, tn, ok := formal(item); ok {
			if err := check(tn); err != nil {
				return err
			}
		}
	}
	return check(result)
}

// charsetScopes checks the item%CHARSET source of each declaration in a
// unit against the declarations in scope, innermost first, as Oracle 23ai
// resolves it: a source of a non-character type is PLS-00550, from a
// package item, an outer block, or an enclosing subprogram alike, while an
// inner declaration of the same name shadows the outer one. A unit is
// walked once, each scope's map filled as its declarations are read.
type charsetScopes struct {
	p      *Parser
	scopes []map[string]*nodes.TypeName
}

// checkCharsetScopes walks a subprogram body (or an anonymous block, with
// no parameters): the formals and the body's declarations share a scope.
func (p *Parser) checkCharsetScopes(params *nodes.List, body nodes.StmtNode) error {
	c := &charsetScopes{p: p}
	return c.subprogram(params, body)
}

// checkPackageCharsetScopes walks the items of a package specification or
// body, and the subprograms among them.
func (p *Parser) checkPackageCharsetScopes(items *nodes.List) error {
	if items == nil {
		return nil
	}
	c := &charsetScopes{p: p}
	c.push()
	return c.declarations(items)
}

func (c *charsetScopes) push() {
	c.scopes = append(c.scopes, map[string]*nodes.TypeName{})
}

func (c *charsetScopes) pop() {
	c.scopes = c.scopes[:len(c.scopes)-1]
}

// declare enters name in the innermost scope; tn is nil for an item that
// is not a typed variable (a cursor, a type, a subprogram), which shadows
// an outer name without being a %CHARSET source the text can judge.
func (c *charsetScopes) declare(name string, tn *nodes.TypeName) {
	c.scopes[len(c.scopes)-1][name] = tn
}

func (c *charsetScopes) check(tn *nodes.TypeName) error {
	if tn == nil || !tn.IsPercCharset || c.p.qualifiedCharsets[tn] {
		return nil
	}
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if src, ok := c.scopes[i][tn.CharacterSet]; ok {
			if src != nil && !c.p.mayBeCharacterType(src) {
				return c.p.syntaxErrorAtType(tn)
			}
			return nil
		}
	}
	return nil
}

func (c *charsetScopes) subprogram(params *nodes.List, body nodes.StmtNode) error {
	block, ok := body.(*nodes.PLSQLBlock)
	if !ok || block == nil {
		return nil
	}
	c.push()
	defer c.pop()
	if params != nil {
		for _, item := range params.Items {
			if pr, ok := item.(*nodes.Parameter); ok {
				c.declare(pr.Name, pr.TypeName)
			}
		}
	}
	return c.blockContents(block)
}

func (c *charsetScopes) block(b *nodes.PLSQLBlock) error {
	c.push()
	defer c.pop()
	return c.blockContents(b)
}

// blockContents walks a block's declarations, then each block nested in its
// statements and handlers, without descending into those blocks twice.
func (c *charsetScopes) blockContents(b *nodes.PLSQLBlock) error {
	if b.Declarations != nil {
		if err := c.declarations(b.Declarations); err != nil {
			return err
		}
	}
	var err error
	visit := func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		if nested, ok := n.(*nodes.PLSQLBlock); ok {
			err = c.block(nested)
			return false
		}
		return true
	}
	if b.Statements != nil {
		nodes.Inspect(b.Statements, visit)
	}
	if b.Exceptions != nil {
		nodes.Inspect(b.Exceptions, visit)
	}
	return err
}

func (c *charsetScopes) declarations(items *nodes.List) error {
	for _, item := range items.Items {
		switch d := item.(type) {
		case *nodes.PLSQLVarDecl:
			if err := c.check(d.TypeName); err != nil {
				return err
			}
			c.declare(d.Name, d.TypeName)
		case *nodes.PLSQLSubtypeDecl:
			if err := c.check(d.BaseType); err != nil {
				return err
			}
			c.declare(d.Name, nil)
		case *nodes.PLSQLCursorDecl:
			c.declare(d.Name, nil)
		case *nodes.PLSQLTypeDecl:
			c.declare(d.Name, nil)
		case *nodes.CreateProcedureStmt:
			if d.Name != nil {
				c.declare(d.Name.Name, nil)
			}
			if err := c.subprogram(d.Parameters, d.Body); err != nil {
				return err
			}
		case *nodes.CreateFunctionStmt:
			if d.Name != nil {
				c.declare(d.Name.Name, nil)
			}
			if err := c.subprogram(d.Parameters, d.Body); err != nil {
				return err
			}
		}
	}
	return nil
}

// tablePseudoType returns the first TABLE or COLUMNS pseudo-type among a
// subprogram's parameter types and its result type, or nil. Both stand only
// in a polymorphic table function: Oracle 23ai rejects them elsewhere with
// PLS-00765, in a procedure, an ordinary function, or a SQL macro alike.
func (p *Parser) tablePseudoType(params *nodes.List, result *nodes.TypeName) *nodes.TypeName {
	if params != nil {
		for _, item := range params.Items {
			if pr, ok := item.(*nodes.Parameter); ok && p.pseudoType(pr.TypeName) != "" {
				return pr.TypeName
			}
		}
	}
	if p.pseudoType(result) != "" {
		return result
	}
	return nil
}

// pseudoType returns TABLE or COLUMNS when tn is that polymorphic table
// function pseudo-type, "" otherwise. Only the unquoted word is the
// pseudo-type: a quoted "TABLE" or "COLUMNS" names a user type, which a
// procedure or function parameter may have, while no type can be named
// COLUMNS unquoted (CREATE TYPE columns is PLS-00103 on Oracle 23ai).
func (p *Parser) pseudoType(tn *nodes.TypeName) string {
	if tn == nil || tn.IsPercType || tn.IsPercRowtype || tn.Names.Len() != 1 {
		return ""
	}
	if tn.Loc.Start >= 0 && tn.Loc.Start < len(p.source) && p.source[tn.Loc.Start] == '"' {
		return ""
	}
	s, ok := tn.Names.Items[0].(*nodes.String)
	if !ok || (s.Str != "TABLE" && s.Str != "COLUMNS") {
		return ""
	}
	return s.Str
}

// syntaxErrorAtType returns a syntax error at datatype tn.
func (p *Parser) syntaxErrorAtType(tn *nodes.TypeName) *ParseError {
	return p.syntaxErrorAtTok(Token{Type: tokIDENT, Loc: tn.Loc.Start, End: tn.Loc.End})
}

// parseParenExprList parses ( expr [, expr ]... ) with at least one
// expression, the current token at '('.
func (p *Parser) parseParenExprList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		list.Items = append(list.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseResultCacheReliesOn parses RELIES_ON ( [ data_source [, data_source ]... ] ),
// the current token at RELIES_ON. The list may be empty on Oracle 23ai.
func (p *Parser) parseResultCacheReliesOn() (*nodes.List, error) {
	p.advance() // consume RELIES_ON
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	list := &nodes.List{}
	for p.cur.Type != ')' {
		// A ',' needs a data source after it: RELIES_ON (t1,) is PLS-00103.
		name, err := p.parseUnitName()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == ')' {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseSqlMacroType parses ( [ TYPE => ] { SCALAR | TABLE } ) after SQL_MACRO,
// the current token at '('.
func (p *Parser) parseSqlMacroType() (string, error) {
	p.advance() // consume '('
	if p.cur.Type == kwTYPE && p.peekNext().Type == tokASSOC {
		p.advance() // consume TYPE
		p.advance() // consume =>
	}
	var kind string
	switch {
	case p.cur.Type == kwTABLE:
		kind = "TABLE"
	case p.isKeywordStr("SCALAR"):
		kind = "SCALAR"
	default:
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return kind, nil
}
