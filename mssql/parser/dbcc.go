// Package parser - dbcc.go implements T-SQL DBCC statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseDbccStmt parses a DBCC (Database Console Command) statement.
//
// BNF: mssql/parser/bnf/dbcc-transact-sql.bnf
//
//	DBCC command_name
//	    [ ( argument [ , ...n ] ) ]
//	    [ WITH option [ , ...n ] ]
func (p *Parser) parseDbccStmt() (*nodes.DbccStmt, error) {
	loc := p.pos()
	p.advance() // consume DBCC

	stmt := &nodes.DbccStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Parse command name - it may be a keyword or identifier.
	// DBCC command names like CHECKDB, SHRINKDATABASE, etc. are always treated as identifiers here.
	if p.isIdentLike() {
		stmt.Command = strings.ToUpper(p.cur.Str)
		p.advance()
	} else {
		// Command could be a reserved keyword used as a DBCC command name.
		// Consume whatever token is there and use its string representation.
		stmt.Command = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// Optional argument list: ( arg [, ...] )
	if _, ok := p.match('('); ok {
		var args []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg, _ := p.parseExpr()
			args = append(args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}
	}

	// Optional WITH options: WITH option [, ...]
	//
	//   WITH { NO_INFOMSGS | ALL_ERRORMSGS | PHYSICAL_ONLY
	//        | EXTENDED_LOGICAL_CHECKS | DATA_PURITY | TABLOCK
	//        | ESTIMATEONLY | COUNT_ROWS | TABLERESULTS } [ , ... ]
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		var opts []nodes.Node
		for {
			if p.isIdentLike() {
				optLoc := p.pos()
				optName := strings.ToUpper(p.cur.Str)
				p.advance()
				opts = append(opts, &nodes.DbccOption{
					Name: optName,
					Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
				})
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
