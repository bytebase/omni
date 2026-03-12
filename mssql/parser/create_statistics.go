// Package parser - create_statistics.go implements T-SQL CREATE/UPDATE/DROP STATISTICS parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateStatisticsStmt parses a CREATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-statistics-transact-sql
//
//	CREATE STATISTICS statistics_name
//	    ON { table_or_indexed_view_name } ( column [ ,...n ] )
//	    [ WHERE <filter_predicate> ]
//	    [ WITH
//	        [ FULLSCAN [ , PERSIST_SAMPLE_PERCENT = { ON | OFF } ] ]
//	        | SAMPLE number { PERCENT | ROWS }
//	        | STATS_STREAM = stats_stream
//	        [ , NORECOMPUTE ]
//	        [ , INCREMENTAL = { ON | OFF } ]
//	        [ , MAXDOP = max_degree_of_parallelism ]
//	        [ , AUTO_DROP = { ON | OFF } ]
//	    ]
func (p *Parser) parseCreateStatisticsStmt() *nodes.CreateStatisticsStmt {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.CreateStatisticsStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Statistics name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// Column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() {
				cols = append(cols, &nodes.String{Str: p.cur.Str})
				p.advance()
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		stmt.Columns = &nodes.List{Items: cols}
	}

	// WHERE filter predicate (optional)
	if p.cur.Type == kwWHERE {
		p.advance()
		p.parseExpr() // consume but don't store
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options = p.parseStatisticsWithOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseUpdateStatisticsStmt parses an UPDATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/update-statistics-transact-sql
//
//	UPDATE STATISTICS table_or_indexed_view_name
//	    [
//	      { index_or_statistics_name }
//	      | ( { index_or_statistics_name } [ ,...n ] )
//	    ]
//	    [ WITH
//	        [ FULLSCAN [ , PERSIST_SAMPLE_PERCENT = { ON | OFF } ] ]
//	        | SAMPLE number { PERCENT | ROWS }
//	        | RESAMPLE [ ON PARTITIONS (...) ]
//	        | <update_stats_stream_option> [ ,...n ]
//	    ]
func (p *Parser) parseUpdateStatisticsStmt() *nodes.UpdateStatisticsStmt {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.UpdateStatisticsStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name
	stmt.Table = p.parseTableRef()

	// Optional statistics name or list
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() {
				if stmt.Name == "" {
					stmt.Name = p.cur.Str
				}
				p.advance()
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
	} else if p.isIdentLike() && p.cur.Type != kwWITH {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options = p.parseStatisticsWithOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropStatisticsStmt parses a DROP STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-statistics-transact-sql
//
//	DROP STATISTICS { table | view }.statistics_name [ ,...n ]
//	(fully-qualified: [schema.]table.statistics_name)
func (p *Parser) parseDropStatisticsStmt() *nodes.DropStatisticsStmt {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.DropStatisticsStmt{
		Loc: nodes.Loc{Start: loc},
	}

	var names []nodes.Node
	for {
		// Each name is [schema.]table.stats_name (two or three-part dot-separated)
		if !p.isIdentLike() {
			break
		}
		// Collect all dotted parts
		parts := p.cur.Str
		p.advance()
		for p.cur.Type == '.' {
			p.advance()
			if p.isIdentLike() {
				parts = parts + "." + p.cur.Str
				p.advance()
			}
		}
		names = append(names, &nodes.String{Str: parts})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Names = &nodes.List{Items: names}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseStatisticsWithOptions parses the WITH clause of CREATE/UPDATE STATISTICS.
func (p *Parser) parseStatisticsWithOptions() *nodes.List {
	var opts []nodes.Node
	for {
		if p.isIdentLike() || p.cur.Type == kwFULL || p.cur.Type == kwNOCOUNT {
			opt := strings.ToUpper(p.cur.Str)
			p.advance()
			// Handle SAMPLE number PERCENT|ROWS
			if strings.EqualFold(opt, "SAMPLE") {
				p.parseExpr() // number
				if p.isIdentLike() {
					p.advance() // PERCENT or ROWS
				}
			} else if p.cur.Type == '=' {
				// key = value
				p.advance()
				p.parseExpr()
			} else if p.cur.Type == kwON || p.cur.Type == kwOFF {
				p.advance()
			}
			opts = append(opts, &nodes.String{Str: opt})
		} else {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}
