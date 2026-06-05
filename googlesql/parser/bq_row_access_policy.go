package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — CREATE ROW ACCESS POLICY + DROP [ALL] ROW ACCESS POLICY/POLICIES
// (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR create_row_access_policy_statement (GoogleSQLParser.g4
// §2.2), the `DROP ROW ACCESS POLICY … ON …` alternative of drop_statement, and
// drop_all_row_access_policies_statement.
//
//	create_row_access_policy_statement:
//	  CREATE opt_or_replace? ROW ACCESS? POLICY opt_if_not_exists? identifier?
//	    ON path_expression create_row_access_policy_grant_to_clause? filter_using_clause
//	  create_row_access_policy_grant_to_clause: grant_to_clause | TO grantee_list
//	  grant_to_clause: GRANT TO ( grantee_list )
//	  filter_using_clause: FILTER? USING ( expression )
//	drop_statement: DROP ROW ACCESS POLICY opt_if_exists? identifier on_path_expression
//	drop_all_row_access_policies_statement:
//	  DROP ALL ROW ACCESS? POLICIES ON path_expression
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md). Probed 2026-06-05:
// `CREATE ROW ACCESS POLICY …` REJECTS on the Spanner emulator ("Encountered
// 'ROW' while parsing: create_statement"). Triangulated against the legacy .g4 +
// BigQuery truth1 (DDL-021/052); proven by the unit tests in
// bq_row_access_policy_test.go.

// parseCreateRowAccessPolicy parses a CREATE ROW ACCESS POLICY statement. The
// shared CREATE prefix (OR REPLACE) has been consumed by parseCreateStmt; cur is
// at ROW.
func (p *Parser) parseCreateRowAccessPolicy(create Token, orReplace bool) (ast.Node, error) {
	p.advance() // ROW
	// ACCESS? — the grammar marks ACCESS optional (ROW [ACCESS] POLICY); accept both.
	if p.cur.Type == kwACCESS {
		p.advance()
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.CreateRowAccessPolicyStmt{OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// identifier? — the policy name, optional in the grammar. It is present unless
	// the next token is ON (the required clause that follows). An identifier-start
	// that is not ON is the policy name.
	if p.cur.Type != kwON && isIdentifierStart(p.cur.Type) {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// ON path_expression — required.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// create_row_access_policy_grant_to_clause?:
	//   grant_to_clause          = GRANT TO ( grantee_list )
	//   | TO grantee_list
	switch p.cur.Type {
	case kwGRANT:
		p.advance() // GRANT
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		grantees, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		stmt.Grantees = grantees
		stmt.HasGrantTo = true
	case kwTO:
		p.advance() // TO
		grantees, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		stmt.Grantees = grantees
		stmt.HasGrantTo = true
	}

	// filter_using_clause: FILTER? USING ( expression ) — required.
	if p.cur.Type == kwFILTER {
		p.advance() // FILTER
	}
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	filter, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	stmt.Filter = filter

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseDropRowAccessPolicy parses either `DROP ROW ACCESS POLICY [IF EXISTS]
// identifier ON path` or `DROP ALL ROW ACCESS POLICIES ON path`. The DROP keyword
// has been consumed by parseDropStmt; cur is at ROW or ALL.
func (p *Parser) parseDropRowAccessPolicy(drop Token) (ast.Node, error) {
	if p.cur.Type == kwALL {
		return p.parseDropAllRowAccessPolicies(drop)
	}

	// DROP ROW ACCESS POLICY [IF EXISTS] identifier ON path.
	p.advance() // ROW
	if _, err := p.expect(kwACCESS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &ast.BQDropStmt{Object: ast.BQDropRowAccessPolicy}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	// identifier — the policy name (required for the single-policy drop).
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// on_path_expression: ON path_expression — required.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	onTable, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.OnTable = onTable

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropAllRowAccessPolicies parses `DROP ALL ROW ACCESS? POLICIES ON
// path_expression`. cur is at ALL (the DROP token is consumed).
func (p *Parser) parseDropAllRowAccessPolicies(drop Token) (ast.Node, error) {
	p.advance() // ALL
	if _, err := p.expect(kwROW); err != nil {
		return nil, err
	}
	// ACCESS? — optional in drop_all_row_access_policies_statement (ROW ACCESS?
	// POLICIES).
	if p.cur.Type == kwACCESS {
		p.advance()
	}
	if _, err := p.expect(kwPOLICIES); err != nil {
		return nil, err
	}

	stmt := &ast.DropAllRowAccessPoliciesStmt{}
	stmt.Loc.Start = drop.Loc.Start

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
