// Package parser implements a recursive descent SQL parser for Oracle PL/SQL.
//
// This parser produces AST nodes from the oracle/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Parser is a recursive descent parser for Oracle SQL/PL/SQL.
//
// Oracle parser helpers follow the project-wide dual-return contract: required
// parse functions return (Node, error), optional probes use a tryParse or
// parseOptional name and may return (nil, nil), and nested parse errors must be
// returned to the caller instead of stored or discarded.
type Parser struct {
	lexer   *Lexer
	source  string // original SQL input
	cur     Token  // current token
	prev    Token  // previous token (for error reporting)
	nextBuf Token  // buffered next token for 2-token lookahead
	hasNext bool   // whether nextBuf is valid

	// caseKinds tells, for each CASE token from the first one parsePrimary
	// met to the end of the range, whether it opens a CASE expression (true)
	// or names a column (false). decideCaseTokens fills it in one pass.
	caseKinds map[int]bool

	// parenQueries tells, for each '(' from the first one a query lookup
	// met to the end of the range, how its group reads as a query.
	// decideParenQueries fills it in one pass.
	parenQueries map[int]parenQuery

	// groupFroms tells, for each parenthesized group from the first TRIM
	// call atTrimSyntax met to the end of the range, by the offset of the
	// token after its '(', whether a FROM stands at the group's own depth.
	// decideGroupFroms fills it in one pass.
	groupFroms map[int]bool

	// queryIterand is set while the iteration control of a FOR loop that
	// opens with '(' is parsed: a query there is a cursor's, which may end
	// in a plain ORDER BY (Oracle 23ai runs FOR r IN (SELECT a FROM t ORDER
	// BY a) LOOP). parseParenExpr takes it for its own parentheses alone,
	// and records the query in iterandQuery and the token after it in
	// iterandAfter, for parseIterandExpr to check it stands alone.
	queryIterand bool
	iterandQuery *nodes.SubqueryExpr
	iterandAfter Token

	// typeMods says what the parenthesized modifiers of the datatype being
	// parsed may hold; parsePLSQLDatatype sets it for PL/SQL datatypes.
	typeMods typeModMode

	// qualifiedCharsets holds the components of each item%CHARSET source
	// that is a qualified name (k.c), as opposed to one identifier that may
	// itself be a quoted "K.C".
	qualifiedCharsets map[*nodes.TypeName][]string

	// typeScopes holds, for each open PL/SQL declaration scope, the names
	// its TYPE and SUBTYPE declarations have declared so far. Oracle 23ai
	// lets one redeclare a predefined type name, NUMBER or SYS_REFCURSOR
	// alike, and the redeclaration then shadows STANDARD's type.
	typeScopes []map[string]bool

	// inPackageBody is set while a package body's items are parsed: its
	// specification may redeclare any predefined type name out of sight.
	inPackageBody bool

	// shadowedTypes records the datatypes spelled as a predefined type
	// where a redeclaration of that name is, or may be, in scope.
	shadowedTypes map[*nodes.TypeName]bool

	// addedConstraints holds the constraints of the ALTER TABLE ADD ( ... )
	// list just parsed, actions of their own after its columns'.
	addedConstraints []*nodes.AlterTableCmd

	// attributePaths records the names of each constraint column that
	// names an attribute of an object column (o.r), which a quoted column
	// name ("O.R") spells alike.
	attributePaths map[*nodes.String][]string

	// inAnonymousBlock is set while a top-level anonymous block is parsed:
	// Oracle 23ai rejects RESULT_CACHE on any subprogram in one (PLS-00999).
	inAnonymousBlock bool

	// plsqlCode is set while a PL/SQL declaration or statement is parsed,
	// and cleared in the SQL statements it embeds: the expressions of
	// PL/SQL alone take the choices of a qualified expression.
	plsqlCode bool
}

// setPLSQLCode sets plsqlCode and returns its previous value, for
// restorePLSQLCode: defer p.restorePLSQLCode(p.setPLSQLCode(true)).
func (p *Parser) setPLSQLCode(plsql bool) bool {
	saved := p.plsqlCode
	p.plsqlCode = plsql
	return saved
}

// restorePLSQLCode restores plsqlCode to what setPLSQLCode returned.
func (p *Parser) restorePLSQLCode(saved bool) {
	p.plsqlCode = saved
}

// parseAnonymousBlock parses a top-level anonymous block, noting that its
// subprograms are in one.
func (p *Parser) parseAnonymousBlock() (nodes.StmtNode, error) {
	saved := p.inAnonymousBlock
	p.inAnonymousBlock = true
	block, err := p.parsePLSQLBlock()
	p.inAnonymousBlock = saved
	if err != nil {
		return nil, err
	}
	if err := p.checkBlockCharsetScopes(block); err != nil {
		return nil, err
	}
	return block, nil
}

// Parse parses a SQL string into an AST list.
// Each statement is wrapped in a *RawStmt.
func Parse(sql string) (*nodes.List, error) {
	return ParseRange(sql, 0, len(sql))
}

// ParseRange parses source[start:end] in place. Every Loc and error Position
// is an offset into source, so a caller that split a script can parse each
// segment with absolute positions without copying or padding the text.
func ParseRange(source string, start, end int) (*nodes.List, error) {
	if err := validateBalancedDelimiters(source, start, end); err != nil {
		return nil, err
	}

	p := &Parser{
		lexer:  NewLexerRange(source, start, end),
		source: source,
	}
	p.advance()

	var stmts []nodes.Node
	needSeparator := false
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			needSeparator = false
			continue
		}
		if needSeparator {
			return nil, p.syntaxErrorAtCur()
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}

		stmtLoc := p.pos()
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		if p.cur.Type == tokEOF && (p.isIncompleteStatementEnd(stmt) || isIncompleteParsedStmt(stmt)) {
			return nil, p.syntaxErrorAtCur()
		}
		if stmt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if lock := lockedSubquery(stmt); lock != nil {
			return nil, p.syntaxErrorAtNode(lock)
		}
		if sel := sourcelessStar(stmt); sel != nil {
			// At the block's start, its SELECT or its WITH, past the
			// parentheses its span takes.
			lex := NewLexerRange(p.source, sel.Loc.Start, sel.Loc.End)
			tok := lex.NextToken()
			for tok.Type == '(' {
				tok = lex.NextToken()
			}
			return nil, p.syntaxErrorAtTok(tok)
		}

		raw := &nodes.RawStmt{
			Stmt: stmt,
			Loc:  nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}
		stmts = append(stmts, raw)
		needSeparator = true
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

func validateBalancedDelimiters(sql string, start, end int) error {
	lexer := NewLexerRange(sql, start, end)
	depth := 0
	var javaHead embeddedSourceHead
	var prev1, prev2, prev3 Token
	for {
		tok := lexer.NextToken()
		if lexer.Err != nil {
			return nil
		}
		if javaHead.observe(tok) {
			// The rest of the range is Java or JavaScript source, not SQL.
			return nil
		}
		if (prev3.Type == kwIS || prev3.Type == kwAS) &&
			prev2.Type == tokIDENT && prev2.Str == "MLE" && prev1.Type == tokIDENT && prev1.Str == "LANGUAGE" {
			// tok names the language of an inline MLE call spec; its
			// delimited code is not SQL.
			if _, codeEnd, _, ok := mleInlineCodeEnd(sql, tok.End, end); ok {
				lexer.pos = codeEnd
			}
		}
		prev3, prev2, prev1 = prev2, prev1, tok
		switch tok.Type {
		case ';':
			javaHead = embeddedSourceHead{}
		case tokEOF:
			if depth > 0 {
				return &ParseError{Message: "syntax error at end of input", Position: tok.Loc}
			}
			return nil
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				text := ")"
				return &ParseError{Message: fmt.Sprintf("syntax error at or near %q", text), Position: tok.Loc}
			}
		}
	}
}

func (p *Parser) isIncompleteStatementEnd(stmt nodes.StmtNode) bool {
	switch p.prev.Type {
	case '(', ',', '.', '=', '+', '-', '*', '/', tokASSIGN, tokASSOC, tokCONCAT:
		return true
	}

	// Keywords Oracle does not reserve (CONTENT, JOIN, USING) are left out: a
	// statement can end with a column of that name, and a clause they open
	// without its operand ("FROM t JOIN") fails in the clause's own parser.
	switch stmt.(type) {
	case *nodes.SelectStmt:
		switch p.prev.Type {
		case kwAND, kwAS, kwBY, kwFROM, kwGROUP, kwHAVING,
			kwIS, kwNOT, kwON, kwOR, kwORDER, kwTHEN, kwUNION,
			kwWHERE:
			return true
		}
	case *nodes.InsertStmt:
		return p.prev.Type == kwINTO || p.prev.Type == kwVALUES
	case *nodes.UpdateStmt:
		return p.prev.Type == kwUPDATE || p.prev.Type == kwSET || p.prev.Type == kwWHERE
	case *nodes.DeleteStmt:
		return p.prev.Type == kwFROM || p.prev.Type == kwWHERE
	case *nodes.MergeStmt:
		switch p.prev.Type {
		case kwINTO, kwON, kwSET, kwTHEN, kwWHERE:
			return true
		}
	case *nodes.CreateTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwDEFAULT
	case *nodes.CreateIndexStmt:
		return p.prev.Type == kwINDEX
	case *nodes.CreateViewStmt:
		return p.prev.Type == kwAS
	case *nodes.AlterTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwADD || p.prev.Type == kwMODIFY
	case *nodes.CreateProcedureStmt, *nodes.CreateFunctionStmt, *nodes.CreatePackageStmt,
		*nodes.CreateTriggerStmt, *nodes.PLSQLBlock:
		return p.prev.Type == kwBEGIN || p.prev.Type == kwTHEN
	case *nodes.CreateUserStmt:
		return p.prev.Type == kwUSER
	case *nodes.GrantStmt:
		return p.prev.Type == kwON || p.prev.Type == kwTO
	case *nodes.DropStmt:
		return p.prev.Type == kwDROP
	}

	return false
}

func isIncompleteParsedStmt(stmt nodes.StmtNode) bool {
	switch s := stmt.(type) {
	case *nodes.SelectStmt:
		return querySelectsSourcelessStar(s)
	case *nodes.DropStmt:
		if s.ObjectType == nodes.OBJECT_DATABASE && (s.Names == nil || s.Names.Len() == 0) {
			return false
		}
		if s.Names == nil || s.Names.Len() == 0 {
			return true
		}
		for _, item := range s.Names.Items {
			name, ok := item.(*nodes.ObjectName)
			if ok && name.Name == "" {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// sourcelessStar returns the first query block in stmt that selects *
// without a FROM clause, which Oracle 23ai rejects with ORA-54058 at the
// block's start, its SELECT or its WITH, wherever it stands: CREATE TABLE ... AS (SELECT *), a view,
// an INSERT, a subquery, an inline view, a WITH element, MERGE USING, or a
// cursor FOR loop. The subquery of EXISTS is the exception when it is one
// query block with no row-limiting clause, as Oracle does not expand its
// select list: WHERE EXISTS (SELECT *) runs, and EXISTS (SELECT * UNION
// SELECT 1 FROM dual) and EXISTS (SELECT * FETCH FIRST 1 ROWS ONLY) fail.
func sourcelessStar(stmt nodes.Node) *nodes.SelectStmt {
	exempt := map[*nodes.SelectStmt]bool{}
	var bad *nodes.SelectStmt
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		if bad != nil {
			return false
		}
		switch n := n.(type) {
		case *nodes.ExistsExpr:
			sel, ok := n.Subquery.(*nodes.SelectStmt)
			// A WITH clause before a parenthesized block wraps it, one
			// block still: Oracle 23ai runs EXISTS (WITH w AS (...)
			// (SELECT *)).
			for ok && sel.Larg != nil && sel.Op == nodes.SETOP_NONE && sel.FetchFirst == nil {
				sel = sel.Larg
			}
			if ok && sel.Larg == nil && sel.Op == nodes.SETOP_NONE && sel.FetchFirst == nil {
				exempt[sel] = true
			}
		case *nodes.SelectStmt:
			// A set operation holds its first block itself, Larg nil, or
			// wraps a parenthesized one in Larg.
			if n.Larg == nil && !exempt[n] &&
				(n.FromClause == nil || n.FromClause.Len() == 0) && selectTargetListContainsStar(n.TargetList) {
				bad = n
			}
		}
		return true
	})
	return bad
}

// lockedSubquery returns the FOR UPDATE clause that ends a subquery in
// stmt. Oracle 23ai locks rows for a top-level query and a cursor's alone:
// it rejects FOR UPDATE in an expression subquery, an inline view or a
// DML target, a WITH element, MERGE USING, CURSOR(...), INSERT ... SELECT,
// and CREATE TABLE or VIEW ... AS (ORA-00907, ORA-03048), and runs it in a
// cursor declaration, OPEN FOR, and a FOR loop's query, VALUES OF (...)
// too. A subquery it does not list keeps its FOR UPDATE.
func lockedSubquery(stmt nodes.Node) nodes.Node {
	var lock nodes.Node
	check := func(query nodes.StmtNode) {
		if sel, ok := query.(*nodes.SelectStmt); ok && lock == nil {
			lock = trailingForUpdate(sel)
		}
	}
	iterands := map[*nodes.SubqueryExpr]bool{}
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		if lock != nil {
			return false
		}
		switch n := n.(type) {
		case *nodes.PLSQLIterationControl:
			if sub, ok := n.Expr.(*nodes.SubqueryExpr); ok {
				iterands[sub] = true
			}
		case *nodes.SubqueryExpr:
			if !iterands[n] {
				check(n.Subquery)
			}
		case *nodes.ExistsExpr:
			check(n.Subquery)
		case *nodes.CursorExpr:
			check(n.Subquery)
		case *nodes.SubqueryRef:
			check(n.Subquery)
		case *nodes.CTE:
			check(n.Query)
		case *nodes.InsertStmt:
			if n.Select != nil {
				check(n.Select)
			}
			check(n.Subquery)
		case *nodes.CreateTableStmt:
			check(n.AsQuery)
		case *nodes.CreateViewStmt:
			check(n.Query)
		}
		return true
	})
	return lock
}

// querySelectsSourcelessStar reports whether a query block of query selects
// * without a FROM clause, which Oracle 23ai rejects with ORA-54058. A set
// operation's operands and a parenthesized operand's wrapper are searched:
// (SELECT *) UNION SELECT 1 FROM dual fails like SELECT * alone.
func querySelectsSourcelessStar(query *nodes.SelectStmt) bool {
	if query == nil {
		return false
	}
	if query.Larg != nil {
		if querySelectsSourcelessStar(query.Larg) {
			return true
		}
	} else if (query.FromClause == nil || query.FromClause.Len() == 0) && selectTargetListContainsStar(query.TargetList) {
		return true
	}
	return querySelectsSourcelessStar(query.Rarg)
}

// selectTargetListContainsStar reports whether list selects * or table.*:
// Oracle 23ai rejects (SELECT x.*) as it does (SELECT *) (ORA-00904), and
// runs EXISTS (SELECT x.*).
func selectTargetListContainsStar(list *nodes.List) bool {
	return list != nil && selectListHasStar(list)
}

// parseStmt dispatches to statement-specific parsers.
// Statement parsers must return syntax and lexer errors directly so Parse can
// fail at the earliest invalid token while preserving partial-node Loc data.
func (p *Parser) parseStmt() (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwSELECT, kwWITH, '(':
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwDROP:
		return p.parseDropStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVEPOINT:
		return p.parseSavepointStmt()
	case kwSET:
		next := p.peekNext()
		if next.Type == kwTRANSACTION {
			return p.parseSetTransactionAtSet()
		}
		if next.Type == kwROLE {
			return p.parseSetRoleStmt()
		}
		if next.Type == kwCONSTRAINT || next.Type == kwCONSTRAINTS {
			return p.parseSetConstraintsStmt()
		}
		return nil, p.syntaxErrorAtCur()
	case kwAUDIT:
		return p.parseAuditStmt()
	case kwNOAUDIT:
		return p.parseNoauditStmt()
	case kwASSOCIATE:
		return p.parseAssociateStatisticsStmt()
	case kwDISASSOCIATE:
		return p.parseDisassociateStatisticsStmt()
	case kwCOMMENT:
		return p.parseCommentStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()
	case kwANALYZE:
		return p.parseAnalyzeStmt()
	case kwEXPLAIN:
		return p.parseExplainPlanStmt()
	case kwFLASHBACK:
		next := p.peekNext()
		if next.Type == kwDATABASE ||
			(next.Type == tokIDENT && (next.Str == "STANDBY" || next.Str == "PLUGGABLE")) {
			return p.parseFlashbackDatabaseStmt()
		}
		return p.parseFlashbackTableStmt()
	case kwPURGE:
		return p.parsePurgeStmt()
	case kwLOCK:
		return p.parseLockTableStmt()
	case kwCALL:
		return p.parseCallStmt()
	case kwRENAME:
		return p.parseRenameStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwCREATE:
		return p.parseCreateStmt()
	case kwDECLARE, kwBEGIN, tokLABELOPEN:
		return p.parseAnonymousBlock()
	default:
		// Handle identifier-based statements
		if p.isIdentLike() {
			switch p.cur.Str {
			case "ADMINISTER":
				return p.parseAdministerKeyManagementStmt()
			}
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// peekAhead returns the n tokens after cur without consuming them; tokens
// past the end of input are tokEOF. It lexes a copy of the lexer, so the
// parser's position is untouched.
func (p *Parser) peekAhead(n int) []Token {
	toks := make([]Token, 0, n)
	if n == 0 {
		return toks
	}
	toks = append(toks, p.peekNext())
	lexer := *p.lexer
	for len(toks) < n {
		toks = append(toks, lexer.NextToken())
	}
	return toks
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// isKeyword checks whether the current token is a specific keyword.
func (p *Parser) isKeyword(kw int) bool {
	return p.cur.Type == kw
}

// matchKeyword consumes the token if it is the given keyword.
func (p *Parser) matchKeyword(kw int) (Token, bool) {
	if p.cur.Type == kw {
		return p.advance(), true
	}
	return Token{}, false
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Severity string // e.g., "ERROR", "WARNING"; defaults to "ERROR"
	Code     string // SQLSTATE code; defaults to "42601"
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	sev := e.Severity
	if sev == "" {
		sev = "ERROR"
	}
	code := e.Code
	if code == "" {
		code = "42601"
	}
	return fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)
}

// tokenText returns the original source text covered by tok.
func (p *Parser) tokenText(tok Token) string {
	if tok.Type == tokEOF {
		return ""
	}
	if tok.Loc >= 0 && tok.End >= tok.Loc && tok.End <= len(p.source) {
		return p.source[tok.Loc:tok.End]
	}
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return ""
}

// syntaxErrorAtCur returns a syntax error for the current token.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	if p.lexer.Err != nil {
		return p.lexerError()
	}
	return p.syntaxErrorAtTok(p.cur)
}

// syntaxErrorAtTok returns a PG-style syntax error for tok.
func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	text := p.tokenText(tok)
	msg := "syntax error at end of input"
	if text != "" {
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}

// lexerError returns the lexer's error with token context.
func (p *Parser) lexerError() *ParseError {
	tok := p.cur
	if tok.Type == tokEOF && p.prev.Type != 0 {
		tok = p.prev
	}
	text := p.tokenText(tok)
	msg := p.lexer.Err.Error()
	if text != "" {
		msg = fmt.Sprintf("%s at or near %q", msg, text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}
