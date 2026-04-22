// Package parser - create_statistics.go implements T-SQL CREATE/UPDATE/DROP STATISTICS parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateStatisticsStmt parses a CREATE STATISTICS statement.
//
// BNF: mssql/parser/bnf/create-statistics-transact-sql.bnf
//
//	CREATE STATISTICS statistics_name
//	    ON { table_or_indexed_view_name } ( column [ ,...n ] )
//	    [ WHERE <filter_predicate> ]
//	    [ WITH
//	        [ FULLSCAN
//	            [ [ , ] PERSIST_SAMPLE_PERCENT = { ON | OFF } ]
//	          | SAMPLE number { PERCENT | ROWS }
//	            [ [ , ] PERSIST_SAMPLE_PERCENT = { ON | OFF } ]
//	          | <update_stats_stream_option> [ ,...n ]
//	        ]
//	        [ [ , ] NORECOMPUTE ]
//	        [ [ , ] INCREMENTAL = { ON | OFF } ]
//	        [ [ , ] MAXDOP = max_degree_of_parallelism ]
//	        [ [ , ] AUTO_DROP = { ON | OFF } ]
//	    ]
func (p *Parser) parseCreateStatisticsStmt() (*nodes.CreateStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.CreateStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Statistics name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ON table
	if _, ok := p.match(kwON); ok {
		var err error
		stmt.Table, err = p.parseTableRef()
		if err != nil {
			return nil, err
		}
	}

	// Column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isAnyKeywordIdent() {
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
		var err error
		stmt.Where, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options = p.parseStatisticsWithOptions()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseUpdateStatisticsStmt parses an UPDATE STATISTICS statement.
//
// BNF: mssql/parser/bnf/update-statistics-transact-sql.bnf
//
//	UPDATE STATISTICS table_or_indexed_view_name
//	    [
//	      { index_or_statistics_name }
//	      | ( { index_or_statistics_name } [ ,...n ] )
//	    ]
//	    [ WITH
//	        [ FULLSCAN
//	            [ [ , ] PERSIST_SAMPLE_PERCENT = { ON | OFF } ]
//	          | SAMPLE number { PERCENT | ROWS }
//	            [ [ , ] PERSIST_SAMPLE_PERCENT = { ON | OFF } ]
//	          | RESAMPLE
//	            [ ON PARTITIONS ( { <partition_number> | <range> } [ ,...n ] ) ]
//	          | <update_stats_stream_option> [ ,...n ]
//	        ]
//	        [ [ , ] [ ALL | COLUMNS | INDEX ]
//	        [ [ , ] NORECOMPUTE ]
//	        [ [ , ] INCREMENTAL = { ON | OFF } ]
//	        [ [ , ] MAXDOP = max_degree_of_parallelism ]
//	        [ [ , ] AUTO_DROP = { ON | OFF } ]
//	    ]
func (p *Parser) parseUpdateStatisticsStmt() (*nodes.UpdateStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.UpdateStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name
	var err error
	stmt.Table, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// Optional statistics name or list
	if p.cur.Type == '(' {
		p.advance()
		var names []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isAnyKeywordIdent() {
				names = append(names, &nodes.String{Str: p.cur.Str})
				p.advance()
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(names) > 0 {
			stmt.Names = &nodes.List{Items: names}
			// Also set Name to first for backward compatibility
			if s, ok := names[0].(*nodes.String); ok {
				stmt.Name = s.Str
			}
		}
	} else if p.isAnyKeywordIdent() && p.cur.Type != kwWITH {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options = p.parseStatisticsWithOptions()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropStatisticsStmt parses a DROP STATISTICS statement.
//
// BNF: mssql/parser/bnf/drop-statistics-transact-sql.bnf
//
//	DROP STATISTICS table.statistics_name | view.statistics_name [ ,...n ]
func (p *Parser) parseDropStatisticsStmt() (*nodes.DropStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.DropStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var names []nodes.Node
	for {
		// Each name is [schema.]table.stats_name (two or three-part dot-separated)
		if !p.isAnyKeywordIdent() {
			break
		}
		// Collect all dotted parts
		parts := p.cur.Str
		p.advance()
		for p.cur.Type == '.' {
			p.advance()
			if p.isAnyKeywordIdent() {
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseStatisticsWithOptions parses the WITH clause of CREATE/UPDATE STATISTICS.
func (p *Parser) parseStatisticsWithOptions() *nodes.List {
	var opts []nodes.Node
	for {
		if p.isAnyKeywordIdent() || p.cur.Type == kwFULL || p.cur.Type == kwNOCOUNT {
			opt := strings.ToUpper(p.cur.Str)
			isSample := p.cur.Type == kwSAMPLE
			isResample := p.cur.Type == kwRESAMPLE
			p.advance()
			// Handle SAMPLE number PERCENT|ROWS
			if isSample {
				p.parseExpr() // number
				if p.isAnyKeywordIdent() {
					p.advance() // PERCENT or ROWS
				}
			} else if isResample {
				// RESAMPLE [ ON PARTITIONS ( { partition_number | range } [ ,...n ] ) ]
				if p.cur.Type == kwON {
					p.advance() // ON
					if p.matchIdentCI("PARTITIONS") {
						if p.cur.Type == '(' {
							p.advance()
							for p.cur.Type != ')' && p.cur.Type != tokEOF {
								p.parseExpr() // partition number
								// TO for range
								if p.cur.Type == kwTO {
									p.advance()
									p.parseExpr()
								}
								if _, ok := p.match(','); !ok {
									break
								}
							}
							p.match(')')
						}
					}
				}
			} else if p.cur.Type == '=' {
				// key = value (e.g., INCREMENTAL = ON, MAXDOP = 4, AUTO_DROP = ON)
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == kwON || p.cur.Type == kwOFF {
					p.advance()
				} else {
					p.parseExpr()
				}
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
