package parser

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Extended GoogleSQL & BigQuery statements:
//   - EXPLAIN <statement>
//   - DEFINE TABLE <path> (<options>)
//   - CREATE MODEL (shared base + BigQuery ML clauses)
//   - CREATE EXTERNAL TABLE / CREATE EXTERNAL SCHEMA (BigQuery)
//   - CREATE CONNECTION / CREATE CONSTANT (BigQuery)
//   - CREATE / ALTER / DROP PRIVILEGE RESTRICTION (BigQuery)
//   - ALTER [ALL] ROW ACCESS POLICY / POLICIES (BigQuery)
//   - ALTER SEARCH | VECTOR INDEX (BigQuery multi-action)
// ---------------------------------------------------------------------------

// parseExplainStmt parses `EXPLAIN <unterminated_sql_statement>`.
// EXPLAIN is the current token.
func (p *Parser) parseExplainStmt() (ast.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	explainTok := p.advance() // EXPLAIN
	if p.cur.Type == int('@') {
		if herr := p.skipStatementLevelHint(); herr != nil {
			return nil, herr
		}
	}
	if p.cur.Type == tokEOF {
		return nil, p.syntaxErrorAtCur()
	}
	inner, err := p.parseStmt()
	if err != nil {
		if pe, ok := err.(*ParseError); ok && strings.HasPrefix(pe.Message, "unknown or unsupported statement") {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, err
	}
	return &ast.ExplainStmt{
		Stmt: inner,
		Loc:  ast.Loc{Start: explainTok.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseDefineStmt parses `DEFINE TABLE path_expression options_list`.
// Other DEFINE forms (such as DEFINE MACRO) remain unsupported.
func (p *Parser) parseDefineStmt() (ast.Node, error) {
	if p.peekNext().Type != kwTABLE {
		return p.unsupported("DEFINE")
	}
	defineTok := p.advance() // DEFINE
	p.advance()              // TABLE

	name, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	opts, err := p.parseOptionsParens()
	if err != nil {
		return nil, err
	}
	stmt := &ast.DefineTableStmt{
		Name:    name,
		Options: opts,
		Loc:     ast.Loc{Start: defineTok.Loc.Start, End: p.prev.Loc.End},
	}
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseCreateModel parses:
//
//	CREATE [OR REPLACE] [scope] MODEL [IF NOT EXISTS] path_expression
//	  [INPUT table_element_list OUTPUT table_element_list]
//	  [TRANSFORM ( select_list )]
//	  [REMOTE [WITH CONNECTION ...] | WITH CONNECTION ...]
//	  [OPTIONS options_list]
//	  [AS query | AS ( aliased_query_list )]
//
// In DialectSpanner, only the shared base clauses (OR REPLACE, IF NOT EXISTS,
// non-empty INPUT/OUTPUT, REMOTE) are accepted. In DialectDefault and
// DialectBigQuery, all BigQuery ML clauses are accepted.
func (p *Parser) parseCreateModel(create Token, orReplace bool, scope string) (ast.Node, error) {
	if scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // MODEL

	stmt := &ast.CreateModelStmt{
		OrReplace: orReplace,
		Scope:     scope,
		Loc:       ast.Loc{Start: create.Loc.Start},
	}

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// input_output_clause?: INPUT table_element_list OUTPUT table_element_list
	// Both Spanner and BigQuery reject empty INPUT() / OUTPUT().
	if p.cur.Type == kwINPUT {
		p.advance() // INPUT
		inCols, _, err := p.parseTableElementList()
		if err != nil {
			return nil, err
		}
		if len(inCols) == 0 {
			return nil, &ParseError{Position: p.prev.Loc.Start, End: p.prev.Loc.End, Message: "syntax error: INPUT column list must not be empty"}
		}
		if _, err := p.expect(kwOUTPUT); err != nil {
			return nil, err
		}
		outCols, _, err := p.parseTableElementList()
		if err != nil {
			return nil, err
		}
		if len(outCols) == 0 {
			return nil, &ParseError{Position: p.prev.Loc.Start, End: p.prev.Loc.End, Message: "syntax error: OUTPUT column list must not be empty"}
		}
		stmt.InputCols = inCols
		stmt.OutputCols = outCols
	}

	// remote_with_connection_clause?:
	//   REMOTE [with_connection_clause] | with_connection_clause
	if p.cur.Type == kwREMOTE {
		p.advance() // REMOTE
		stmt.Remote = true
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}
