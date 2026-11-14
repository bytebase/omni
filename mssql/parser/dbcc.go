// Package parser - dbcc.go implements T-SQL DBCC statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// dbccOptionNames are the options DBCC takes after WITH: SQL Server refuses
// any other word (Msg 195), whatever the command.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/database-console-commands/dbcc-transact-sql
var dbccOptionNames = map[string]bool{
	"ALL_ERRORMSGS": true, "ALL_LEVELS": true, "COUNT_ROWS": true, "DATA_PURITY": true,
	"DENSITY_VECTOR": true, "ESTIMATEONLY": true, "EXTENDED_LOGICAL_CHECKS": true, "FAST": true,
	"HISTOGRAM": true, "NO_INFOMSGS": true, "PHYSICAL_ONLY": true, "STAT_HEADER": true,
	"STATS_STREAM": true, "TABLERESULTS": true, "TABLOCK": true, "MAXDOP": true,
}

// parseDbccStmt parses a DBCC statement.
//
// BNF: mssql/parser/bnf/dbcc-transact-sql.bnf
//
//	DBCC command_name [ ( [ [ name = ] value [ ,...n ] ] ) ]
//	    [ WITH option [ ,...n ] ]
//
// dbccStatement (TSql160.g:6381): the command is required; each value is a
// constant, a negative number or a name (dbccNamedLiteral, TSql160.g:6502);
// each option is an undelimited word of dbccOptionNames, MAXDOP with
// = integer. SQL Server refuses the JOIN form of the options (Msg 195).
func (p *Parser) parseDbccStmt() (*nodes.DbccStmt, error) {
	loc := p.pos()
	p.advance() // consume DBCC

	stmt := &nodes.DbccStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if !p.isKeywordOrIdent() {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Command = strings.ToUpper(p.cur.Str)
	p.advance()

	if _, ok := p.match('('); ok {
		var args []nodes.Node
		for p.cur.Type != ')' {
			// [ name = ] value
			if p.isIdentLike() && p.peekNext().Type == '=' {
				p.advance() // name
				p.advance() // =
			}
			arg, err := p.parseProcConstantOrIdentifier()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
			if p.cur.Type == ')' {
				return nil, p.syntaxErrorAtCur()
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}
	}

	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		var opts []nodes.Node
		for {
			if !p.isKeywordOrIdent() || p.isDelimited(p.cur) || !dbccOptionNames[strings.ToUpper(p.cur.Str)] {
				return nil, p.syntaxErrorAtCur()
			}
			optLoc := p.pos()
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if optName == "MAXDOP" {
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				if !procIsInteger(p.cur) {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
			}
			opts = append(opts, &nodes.DbccOption{
				Name: optName,
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
