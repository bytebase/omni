// Package parser - alter_objects.go implements T-SQL ALTER DATABASE and ALTER INDEX parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseAlterDatabaseStmt parses an ALTER DATABASE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql
//
//	ALTER DATABASE { database_name | CURRENT }
//	    { SET ... | MODIFY FILE ... | ADD FILE ... | COLLATE ... | ... }
//
// We capture the database name and the action keyword, then skip the rest.
func (p *Parser) parseAlterDatabaseStmt() *nodes.AlterDatabaseStmt {
	loc := p.pos()

	stmt := &nodes.AlterDatabaseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Database name or CURRENT
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Action keyword: SET, MODIFY, ADD, REMOVE, COLLATE, etc.
	if p.isIdentLike() || p.cur.Type == kwSET {
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// Skip the rest of the statement tokens (up to ; or EOF or statement-start keyword)
	for p.cur.Type != tokEOF && p.cur.Type != ';' && !p.isStatementStart() {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterIndexStmt parses an ALTER INDEX statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-transact-sql
//
//	ALTER INDEX { index_name | ALL } ON <object>
//	    { REBUILD [ PARTITION = ... ] | REORGANIZE [ PARTITION = ... ] | DISABLE | SET ( ... ) }
func (p *Parser) parseAlterIndexStmt() *nodes.AlterIndexStmt {
	loc := p.pos()

	stmt := &nodes.AlterIndexStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Index name or ALL
	if p.isIdentLike() {
		stmt.IndexName = p.cur.Str
		p.advance()
	}

	// ON table_name
	if p.cur.Type == kwON {
		p.advance() // consume ON
		stmt.Table = p.parseTableRef()
	}

	// Action: REBUILD, REORGANIZE, DISABLE, SET
	if p.cur.Type == kwSET {
		stmt.Action = "SET"
		p.advance()
	} else if p.isIdentLike() {
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// Skip any remaining tokens (e.g. PARTITION = ..., WITH (...), etc.)
	// Handle nested parentheses so that WITH (FILLFACTOR = 80) is fully consumed.
	// We treat WITH here as part of the index action (e.g. REBUILD WITH (...))
	// rather than a new statement start.
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' && depth == 0 {
			break
		}
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
				p.advance()
				continue
			}
			break
		}
		// Stop at statement-start keywords except WITH (which can be REBUILD WITH (...))
		if depth == 0 && p.isStatementStart() && p.cur.Type != kwWITH {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
