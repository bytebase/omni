package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — CREATE VIEW (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Ports the plain-view alternative of the legacy ANTLR create_view_statement
// (GoogleSQLParser.g4 §2.2):
//
//	CREATE opt_or_replace? opt_create_scope? RECURSIVE? VIEW opt_if_not_exists?
//	  maybe_dashed_path_expression column_with_options_list? opt_sql_security_clause?
//	  opt_options_list? as_query
//
// The MATERIALIZED and APPROX view alternatives (which add PARTITION/CLUSTER and
// the REPLICA OF source) are owned by parser-ddl-bigquery and route to the
// unsupported stub from parseCreateStmt before reaching here.
//
// ORACLE NOTE — the Spanner emulator (the harness) is a SUBSET oracle for VIEW:
//   - it accepts only `SQL SECURITY INVOKER`, rejecting `DEFINER` with a parse-
//     prefixed error, even though the legacy .g4 (sql_security_clause_kind:
//     INVOKER | DEFINER) and BigQuery accept DEFINER. So a DEFINER reject from the
//     emulator is NON-AUTHORITATIVE; the parser accepts both (divergence).
//   - it requires the SQL SECURITY clause and named view columns (semantic), which
//     the grammar leaves optional — those are accepted-but-semantically-rejected
//     by the emulator (classified ACCEPT) and harmless.
// Verified against the live emulator in view_oracle_test.go.

// parseCreateView parses the plain CREATE VIEW body after the shared CREATE
// prefix (OR REPLACE / scope) and an optional RECURSIVE have been consumed and
// cur is at the VIEW keyword.
func (p *Parser) parseCreateView(create Token, orReplace bool, scope string, recursive bool) (ast.Node, error) {
	p.advance() // consume VIEW

	stmt := &ast.CreateViewStmt{OrReplace: orReplace, Scope: scope, Recursive: recursive}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// column_with_options_list?  — `( column_with_options (, …)* )`, where each is
	// `identifier opt_options_list?`. Distinguished from anything else by the
	// leading '('.
	if p.cur.Type == int('(') {
		cols, err := p.parseViewColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// opt_sql_security_clause?  — SQL SECURITY {INVOKER | DEFINER}.
	if p.cur.Type == kwSQL {
		sec, err := p.parseSQLSecurityClause()
		if err != nil {
			return nil, err
		}
		stmt.SQLSecurity = sec
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// as_query  — AS query (required).
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	stmt.AsQuery = q

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseViewColumnList parses a column_with_options_list: `( column_with_options
// (, column_with_options)* )` where cur is the opening '(' and each entry is
// `identifier opt_options_list?`. Requires at least one column (the grammar's
// list is non-empty).
func (p *Parser) parseViewColumnList() ([]*ast.ViewColumn, error) {
	p.advance() // consume '('
	var cols []*ast.ViewColumn
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		col := &ast.ViewColumn{Name: name, Loc: tok.Loc}
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			col.Options = opts
		}
		col.Loc.End = p.prev.Loc.End
		cols = append(cols, col)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseSQLSecurityClause parses an opt_sql_security_clause: `SQL SECURITY
// {INVOKER | DEFINER}`. cur is at SQL. Both kinds are accepted (the legacy .g4
// and BigQuery permit DEFINER; the Spanner emulator's INVOKER-only restriction is
// a non-authoritative semantic reject — divergence).
func (p *Parser) parseSQLSecurityClause() (string, error) {
	p.advance() // SQL
	if _, err := p.expect(kwSECURITY); err != nil {
		return "", err
	}
	switch p.cur.Type {
	case kwINVOKER:
		p.advance()
		return "INVOKER", nil
	case kwDEFINER:
		p.advance()
		return "DEFINER", nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}
