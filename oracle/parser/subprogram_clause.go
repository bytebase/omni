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
	clause.Columns = cols
	clause.Loc.End = p.prev.End
	return clause, argTok, nil
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
	if name, ok := predefinedTypeName(param.TypeName); ok {
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
func predefinedTypeName(tn *nodes.TypeName) (string, bool) {
	if tn == nil {
		return "", false
	}
	if tn.IsPercRowtype {
		return "%ROWTYPE", true
	}
	if tn.IsPercType || tn.Names.Len() != 1 {
		return "", false
	}
	s, ok := tn.Names.Items[0].(*nodes.String)
	if !ok {
		return "", false
	}
	switch s.Str {
	case "CHAR", "CHARACTER", "VARCHAR2", "VARCHAR", "STRING", "CLOB":
		return s.Str, true
	}
	return s.Str, plsqlNonCharacterTypes[s.Str]
}

// checkPolymorphicSignature checks a polymorphic table function's signature
// as Oracle 23ai does from the text: it returns TABLE (PLS-00767) and takes
// exactly one TABLE parameter (PLS-00773, PLS-00766).
func (p *Parser) checkPolymorphicSignature(stmt *nodes.CreateFunctionStmt) error {
	if !isTablePseudoType(stmt.ReturnType) {
		return p.syntaxErrorAtType(stmt.ReturnType)
	}
	tables := 0
	if stmt.Parameters != nil {
		for _, item := range stmt.Parameters.Items {
			if pr, ok := item.(*nodes.Parameter); ok && isTablePseudoType(pr.TypeName) {
				tables++
				if tables > 1 {
					return p.syntaxErrorAtType(pr.TypeName)
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
	byName := map[string]*nodes.Parameter{}
	for _, item := range params.Items {
		if pr, ok := item.(*nodes.Parameter); ok {
			byName[pr.Name] = pr
		}
	}
	check := func(tn *nodes.TypeName) error {
		if tn == nil || !tn.IsPercCharset {
			return nil
		}
		if src, ok := byName[tn.CharacterSet]; ok && !mayBeCharacterType(src.TypeName) {
			return p.syntaxErrorAtType(tn)
		}
		return nil
	}
	for _, item := range params.Items {
		if pr, ok := item.(*nodes.Parameter); ok {
			if err := check(pr.TypeName); err != nil {
				return err
			}
		}
	}
	return check(result)
}

// tablePseudoType returns the TABLE pseudo-type among a subprogram's
// parameter types and its result type, or nil. TABLE stands only in a
// polymorphic table function: Oracle 23ai rejects it elsewhere with
// PLS-00765, in a procedure, an ordinary function, or a SQL macro alike.
func tablePseudoType(params *nodes.List, result *nodes.TypeName) *nodes.TypeName {
	if params != nil {
		for _, item := range params.Items {
			if pr, ok := item.(*nodes.Parameter); ok && isTablePseudoType(pr.TypeName) {
				return pr.TypeName
			}
		}
	}
	if isTablePseudoType(result) {
		return result
	}
	return nil
}

// isTablePseudoType reports whether tn is the TABLE pseudo-type.
func isTablePseudoType(tn *nodes.TypeName) bool {
	if tn == nil || tn.IsPercType || tn.IsPercRowtype || tn.Names.Len() != 1 {
		return false
	}
	s, ok := tn.Names.Items[0].(*nodes.String)
	return ok && s.Str == "TABLE"
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
