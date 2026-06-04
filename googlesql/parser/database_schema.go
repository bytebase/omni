package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — CREATE SCHEMA / CREATE DATABASE (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR create_schema_statement and create_database_statement
// (GoogleSQLParser.g4 §2.2):
//
//	create_schema_statement:
//	  CREATE opt_or_replace? SCHEMA opt_if_not_exists? path_expression
//	    opt_default_collate_clause? opt_options_list?
//	create_database_statement:
//	  CREATE DATABASE path_expression opt_options_list?
//
// CREATE SCHEMA is shared by both dialects (Spanner named schema; BigQuery
// dataset). CREATE DATABASE is Spanner (BigQuery has no CREATE DATABASE — its
// CreateDatabase support is gated off, see contract.md). DROP SCHEMA lives in
// drop.go. CREATE EXTERNAL SCHEMA is a BigQuery dialect-node form.
//
// Verified against the live Spanner emulator (database_schema_oracle_test.go):
// CREATE SCHEMA [IF NOT EXISTS] and CREATE DATABASE accept; BigQuery's
// `DEFAULT COLLATE` / `OPTIONS` on CREATE SCHEMA are non-authoritative against
// the Spanner oracle (triangulated from the .g4 + BigQuery truth1).

// parseCreateSchema parses a CREATE SCHEMA body after the shared CREATE prefix
// (OR REPLACE) has been consumed and cur is at the SCHEMA keyword. (opt_create_scope
// is not part of create_schema_statement; a scope before SCHEMA is rejected by
// parseCreateStmt's dispatch reaching here only via the SCHEMA case, which is not
// preceded by a scope check — so a `CREATE TEMP SCHEMA` would arrive here with
// scope set and we reject it.)
func (p *Parser) parseCreateSchema(create Token, orReplace bool) (ast.Node, error) {
	p.advance() // consume SCHEMA

	stmt := &ast.CreateSchemaStmt{OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// path_expression (the schema name; the grammar uses a plain path_expression,
	// not a dashed one, but parseTablePath is a superset that also handles the
	// undashed case).
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_default_collate_clause?  — DEFAULT collate_clause (BigQuery).
	if p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		stmt.DefaultCollate = coll
	}

	// opt_options_list?  — OPTIONS ( … ) (BigQuery).
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCreateDatabase parses a CREATE DATABASE body after the shared CREATE
// prefix has been consumed and cur is at the DATABASE keyword. create_database_statement
// is `CREATE DATABASE path_expression opt_options_list?` — it takes NO OR REPLACE
// or scope, so a leading OR REPLACE / scope is a syntax error.
func (p *Parser) parseCreateDatabase(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		// CREATE DATABASE has no opt_or_replace / opt_create_scope production.
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume DATABASE

	stmt := &ast.CreateDatabaseStmt{}
	stmt.Loc.Start = create.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
