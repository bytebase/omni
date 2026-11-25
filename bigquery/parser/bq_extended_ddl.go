package parser

import (
	"strings"

	"github.com/bytebase/omni/bigquery/ast"
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

	// transform_clause?: TRANSFORM ( select_list ) (BigQuery)
	if p.cur.Type == kwTRANSFORM {
		p.advance() // TRANSFORM
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		items, err := p.parseSelectList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		stmt.Transform = items
	}

	// remote_with_connection_clause?:
	//   REMOTE [with_connection_clause] | with_connection_clause
	if p.cur.Type == kwREMOTE {
		p.advance() // REMOTE
		stmt.Remote = true
		if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
			conn, err := p.parseWithConnectionPath()
			if err != nil {
				return nil, err
			}
			stmt.Connection = conn
		}
	} else if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		conn, err := p.parseWithConnectionPath()
		if err != nil {
			return nil, err
		}
		stmt.Connection = conn
	}

	// options_?: OPTIONS options_list (BigQuery)
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// as_query_or_aliased_query_list?: AS query | AS ( aliased_query_list ) (BigQuery)
	if p.cur.Type == kwAS {
		p.advance() // AS
		if p.cur.Type == int('(') && p.isAliasedQueryListAhead() {
			p.advance() // '('
			for {
				cte, err := p.parseCTE()
				if err != nil {
					return nil, err
				}
				stmt.AliasedQueries = append(stmt.AliasedQueries, cte)
				if p.cur.Type != int(',') {
					break
				}
				p.advance()
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
		} else {
			q, err := p.parseQuery()
			if err != nil {
				return nil, err
			}
			stmt.AsQuery = q
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// isAliasedQueryListAhead reports whether cur ('(') opens a parenthesized
// aliased_query_list `( id AS ( query ), ... )` rather than a parenthesized query.
func (p *Parser) isAliasedQueryListAhead() bool {
	if p.cur.Type != int('(') {
		return false
	}
	next := p.peekNext()
	if (!isIdentifierStart(next.Type) && !p.isTmNonReservedWord(next.Type)) || isExtraReservedIdent(next) {
		return false
	}
	return p.peekSecond().Type == kwAS
}

// parseWithConnectionPath parses `WITH CONNECTION (path_expression | DEFAULT)`.
func (p *Parser) parseWithConnectionPath() (*ast.PathExpr, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCONNECTION); err != nil {
		return nil, err
	}
	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		return &ast.PathExpr{Parts: []string{"DEFAULT"}, Loc: tok.Loc}, nil
	}
	return p.parseTablePath()
}

// parseCreateExternal parses BigQuery `CREATE [OR REPLACE] [scope] EXTERNAL
// TABLE ...` and `CREATE [OR REPLACE] [scope] EXTERNAL SCHEMA ...`.
func (p *Parser) parseCreateExternal(create Token, orReplace bool, scope string) (ast.Node, error) {
	p.advance() // EXTERNAL
	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateExternalTable(create, orReplace, scope)
	case kwSCHEMA:
		return p.parseCreateExternalSchema(create, orReplace, scope)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseCreateExternalTable parses:
//
//	CREATE [OR REPLACE] [scope] EXTERNAL TABLE [IF NOT EXISTS] maybe_dashed_path_expression
//	  [table_element_list] [LIKE path] [DEFAULT COLLATE ...]
//	  [WITH PARTITION COLUMNS [table_element_list]]
//	  [WITH CONNECTION ...] OPTIONS (...)
func (p *Parser) parseCreateExternalTable(create Token, orReplace bool, scope string) (ast.Node, error) {
	p.advance() // TABLE
	stmt := &ast.CreateTableStmt{
		OrReplace: orReplace,
		Scope:     scope,
		External:  true,
		Loc:       ast.Loc{Start: create.Loc.Start},
	}

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

	if p.cur.Type == int('(') {
		cols, cons, err := p.parseTableElementList()
		if err != nil {
			return nil, err
		}
		if len(cols) == 0 && len(cons) == 0 {
			return nil, &ParseError{
				Position: stmt.Loc.Start, End: p.prev.Loc.End,
				Message: "syntax error: a table must define at least one column",
			}
		}
		stmt.Columns = cols
		stmt.Constraints = cons
	}

	if p.cur.Type == kwLIKE {
		p.advance()
		likePath, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.Like = likePath
	}

	if p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		stmt.DefaultCollate = coll
	}

	// with_partition_columns_clause?: WITH PARTITION COLUMNS [table_element_list]
	if p.cur.Type == kwWITH && p.peekNext().Type == kwPARTITION {
		p.advance() // WITH
		p.advance() // PARTITION
		if _, err := p.expect(kwCOLUMNS); err != nil {
			return nil, err
		}
		if p.cur.Type == int('(') {
			pcols, _, err := p.parseTableElementList()
			if err != nil {
				return nil, err
			}
			stmt.PartitionColumns = pcols
		}
	}

	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		if err := p.skipWithConnection(); err != nil {
			return nil, err
		}
	}

	// options_ is required for CREATE EXTERNAL TABLE.
	if p.cur.Type != kwOPTIONS {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseOptionsList()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseCreateExternalSchema parses:
//
//	CREATE [OR REPLACE] [scope] EXTERNAL SCHEMA [IF NOT EXISTS] path_expression
//	  [WITH CONNECTION ...] OPTIONS (...)
func (p *Parser) parseCreateExternalSchema(create Token, orReplace bool, scope string) (ast.Node, error) {
	p.advance() // SCHEMA
	stmt := &ast.CreateSchemaStmt{
		OrReplace: orReplace,
		Scope:     scope,
		External:  true,
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

	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		conn, err := p.parseWithConnectionPath()
		if err != nil {
			return nil, err
		}
		stmt.Connection = conn
	}

	// options_ is required for CREATE EXTERNAL SCHEMA.
	if p.cur.Type != kwOPTIONS {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseOptionsList()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseCreateConnection parses:
//
//	CREATE [OR REPLACE] CONNECTION [IF NOT EXISTS] path_expression [OPTIONS(...)]
func (p *Parser) parseCreateConnection(create Token, orReplace bool, scope string) (ast.Node, error) {
	if scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // CONNECTION
	stmt := &ast.CreateConnectionStmt{
		OrReplace: orReplace,
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

	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseCreateConstant parses:
//
//	CREATE [OR REPLACE] [scope] CONSTANT [IF NOT EXISTS] path_expression = expression
func (p *Parser) parseCreateConstant(create Token, orReplace bool, scope string) (ast.Node, error) {
	p.advance() // CONSTANT
	stmt := &ast.CreateConstantStmt{
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

	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Value = val

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parsePossiblyEmptyGranteeList parses '(' [grantee_list] ')'.
func (p *Parser) parsePossiblyEmptyGranteeList() ([]*ast.Grantee, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if p.cur.Type == int(')') {
		p.advance()
		return nil, nil
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return grantees, nil
}

// parsePrivilegeRestrictionTarget parses `ON privilege_list ON identifier path_expression`.
func (p *Parser) parsePrivilegeRestrictionTarget(stmt *ast.PrivilegeRestrictionStmt) error {
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	privs, err := p.parsePrivilegeList()
	if err != nil {
		return err
	}
	stmt.Privileges = privs

	if _, err := p.expect(kwON); err != nil {
		return err
	}
	objTok, err := p.expectIdentifier()
	if err != nil {
		return err
	}
	objType, err := p.identifierText(objTok)
	if err != nil {
		return err
	}
	stmt.ObjectType = objType

	path, err := p.parsePathExpr()
	if err != nil {
		return err
	}
	stmt.Path = path
	return nil
}

// parseCreatePrivilegeRestriction parses:
//
//	CREATE [OR REPLACE] PRIVILEGE RESTRICTION [IF NOT EXISTS]
//	  ON privilege_list ON identifier path_expression [RESTRICT TO ( [grantee_list] )]
func (p *Parser) parseCreatePrivilegeRestriction(create Token, orReplace bool, scope string) (ast.Node, error) {
	if scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // PRIVILEGE
	if _, err := p.expect(kwRESTRICTION); err != nil {
		return nil, err
	}
	stmt := &ast.PrivilegeRestrictionStmt{
		Action:    "CREATE",
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: create.Loc.Start},
	}

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	if err := p.parsePrivilegeRestrictionTarget(stmt); err != nil {
		return nil, err
	}

	if p.cur.Type == kwRESTRICT {
		p.advance() // RESTRICT
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		grantees, err := p.parsePossiblyEmptyGranteeList()
		if err != nil {
			return nil, err
		}
		stmt.RestrictTo = grantees
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterPrivilegeRestriction parses:
//
//	ALTER PRIVILEGE RESTRICTION [IF EXISTS] ON privilege_list ON identifier path_expression
//	  privilege_restriction_alter_action (, privilege_restriction_alter_action)*
func (p *Parser) parseAlterPrivilegeRestriction(alter Token) (ast.Node, error) {
	p.advance() // PRIVILEGE
	if _, err := p.expect(kwRESTRICTION); err != nil {
		return nil, err
	}
	stmt := &ast.PrivilegeRestrictionStmt{
		Action: "ALTER",
		Loc:    ast.Loc{Start: alter.Loc.Start},
	}

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	if err := p.parsePrivilegeRestrictionTarget(stmt); err != nil {
		return nil, err
	}

	for {
		actStart := p.cur.Loc.Start
		act := &ast.PrivilegeRestrictionAction{Loc: ast.Loc{Start: actStart}}
		switch p.cur.Type {
		case kwRESTRICT:
			p.advance() // RESTRICT
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			grantees, err := p.parsePossiblyEmptyGranteeList()
			if err != nil {
				return nil, err
			}
			act.Kind = "RESTRICT TO"
			act.Grantees = grantees
		case kwADD:
			p.advance() // ADD
			ine, err := p.parseIfNotExists()
			if err != nil {
				return nil, err
			}
			grantees, err := p.parsePossiblyEmptyGranteeList()
			if err != nil {
				return nil, err
			}
			act.Kind = "ADD"
			act.IfNotExists = ine
			act.Grantees = grantees
		case kwREMOVE:
			p.advance() // REMOVE
			ie, err := p.parseIfExists()
			if err != nil {
				return nil, err
			}
			grantees, err := p.parsePossiblyEmptyGranteeList()
			if err != nil {
				return nil, err
			}
			act.Kind = "REMOVE"
			act.IfExists = ie
			act.Grantees = grantees
		default:
			return nil, p.syntaxErrorAtCur()
		}
		act.Loc.End = p.prev.Loc.End
		stmt.Actions = append(stmt.Actions, act)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropPrivilegeRestriction parses:
//
//	DROP PRIVILEGE RESTRICTION [IF EXISTS] ON privilege_list ON identifier path_expression
func (p *Parser) parseDropPrivilegeRestriction(drop Token) (ast.Node, error) {
	p.advance() // PRIVILEGE
	if _, err := p.expect(kwRESTRICTION); err != nil {
		return nil, err
	}
	stmt := &ast.PrivilegeRestrictionStmt{
		Action: "DROP",
		Loc:    ast.Loc{Start: drop.Loc.Start},
	}

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	if err := p.parsePrivilegeRestrictionTarget(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterRowAccessPolicy parses:
//
//	ALTER ROW ACCESS POLICY [IF EXISTS] identifier ON path_expression row_access_policy_alter_action_list
//	ALTER ALL ROW ACCESS POLICIES ON path_expression row_access_policy_alter_action
func (p *Parser) parseAlterRowAccessPolicy(alter Token) (ast.Node, error) {
	stmt := &ast.AlterRowAccessPolicyStmt{Loc: ast.Loc{Start: alter.Loc.Start}}

	if p.cur.Type == kwALL {
		p.advance() // ALL
		stmt.All = true
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICIES); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		table, err := p.parsePathExpr()
		if err != nil {
			return nil, err
		}
		stmt.Table = table

		act, err := p.parseRowAccessPolicyAlterAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, act)
	} else {
		p.advance() // ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		stmt.IfExists = ifExists

		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		stmt.Name = name

		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		table, err := p.parsePathExpr()
		if err != nil {
			return nil, err
		}
		stmt.Table = table

		for {
			act, err := p.parseRowAccessPolicyAlterAction()
			if err != nil {
				return nil, err
			}
			stmt.Actions = append(stmt.Actions, act)
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

func (p *Parser) parseRowAccessPolicyAlterAction() (*ast.RowAccessPolicyAction, error) {
	act := &ast.RowAccessPolicyAction{Loc: ast.Loc{Start: p.cur.Loc.Start}}
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
		act.Kind = "GRANT TO"
		act.Grantees = grantees
	case kwFILTER:
		p.advance() // FILTER
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		act.Kind = "FILTER USING"
		act.Filter = expr
	case kwREVOKE:
		p.advance() // REVOKE
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		if p.cur.Type == kwALL {
			p.advance()
			act.Kind = "REVOKE FROM ALL"
		} else {
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
			act.Kind = "REVOKE FROM"
			act.Grantees = grantees
		}
	case kwRENAME:
		p.advance() // RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		act.Kind = "RENAME TO"
		act.RenameTo = name
	default:
		return nil, p.syntaxErrorAtCur()
	}
	act.Loc.End = p.prev.Loc.End
	return act, nil
}

// parseBQAlterSearchVectorIndex parses BigQuery:
//
//	ALTER (SEARCH | VECTOR) INDEX [IF EXISTS] path_expression [ON path_expression]
//	  alter_index_action (, alter_index_action)*
func (p *Parser) parseBQAlterSearchVectorIndex(alter Token) (ast.Node, error) {
	objKind := ast.BQAlterSearchIndex
	if p.cur.Type == kwVECTOR {
		objKind = ast.BQAlterVectorIndex
	}
	p.advance() // SEARCH | VECTOR
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	stmt := &ast.BQAlterStmt{
		Object: objKind,
		Loc:    ast.Loc{Start: alter.Loc.Start},
	}

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Type == kwON {
		p.advance() // ON
		table, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.OnTable = table
	}

	for {
		act, err := p.parseAlterIndexAction()
		if err != nil {
			return nil, err
		}
		if act.Kind == "REBUILD" {
			stmt.Rebuild = true
			if len(act.Options) > 0 {
				stmt.SetOptions = act.Options
			}
		} else if act.Kind == "SET OPTIONS" && len(act.Options) > 0 {
			stmt.SetOptions = act.Options
		}
		stmt.IndexActions = append(stmt.IndexActions, act)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

func (p *Parser) parseAlterIndexAction() (*ast.AlterIndexAction, error) {
	act := &ast.AlterIndexAction{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	switch {
	case p.cur.Type == kwSET:
		p.advance() // SET
		if p.cur.Type != kwOPTIONS {
			return nil, p.syntaxErrorAtCur()
		}
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		act.Kind = "SET OPTIONS"
		act.Options = opts
	case p.cur.Type == kwADD:
		p.advance() // ADD
		if _, err := p.expect(kwCOLUMN); err != nil {
			return nil, err
		}
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		colTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		col, err := p.identifierText(colTok)
		if err != nil {
			return nil, err
		}
		act.Kind = "ADD COLUMN"
		act.IfNotExists = ine
		act.Column = col
		// Optional column type (e.g. BigQuery `ADD COLUMN a <type>`) before OPTIONS / ',' / EOF.
		if p.cur.Type != kwOPTIONS && p.cur.Type != int(',') && p.cur.Type != tokEOF && p.cur.Type != int(';') {
			dt, err := p.parseType()
			if err != nil {
				return nil, err
			}
			act.Type = &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
		}
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			act.Options = opts
		}
	case p.cur.Type == kwDROP:
		p.advance() // DROP
		if _, err := p.expect(kwCOLUMN); err != nil {
			return nil, err
		}
		ie, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		colTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		col, err := p.identifierText(colTok)
		if err != nil {
			return nil, err
		}
		act.Kind = "DROP COLUMN"
		act.IfExists = ie
		act.Column = col
	case p.cur.Type == kwALTER:
		p.advance() // ALTER
		if _, err := p.expect(kwCOLUMN); err != nil {
			return nil, err
		}
		ie, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		colTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		col, err := p.identifierText(colTok)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		if p.cur.Type != kwOPTIONS {
			return nil, p.syntaxErrorAtCur()
		}
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		act.Kind = "ALTER COLUMN"
		act.IfExists = ie
		act.Column = col
		act.Options = opts
	case p.curIsWord("REBUILD"):
		p.advance() // REBUILD
		act.Kind = "REBUILD"
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			act.Options = opts
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	act.Loc.End = p.prev.Loc.End
	return act, nil
}

// isKnownBigQueryEntityType reports whether tok is a recognized generic entity
// type (BigQuery CAPACITY, RESERVATION, ASSIGNMENT, PROJECT, ORGANIZATION,
// FOLDER; Spanner PLACEMENT). Arbitrary identifiers like `a` or `LOCALITY` are
// rejected, matching both BigQuery (`<word> is not a supported object type`)
// and Spanner.
func isKnownBigQueryEntityType(tok Token) bool {
	if tok.Type == kwPROJECT {
		return true
	}
	if tok.Type != tokIdentifier || tok.Loc.End-tok.Loc.Start != len(tok.Str) {
		return false
	}
	switch strings.ToUpper(tok.Str) {
	case "CAPACITY", "RESERVATION", "ASSIGNMENT", "ORGANIZATION", "FOLDER", "PLACEMENT":
		return true
	}
	return false
}
