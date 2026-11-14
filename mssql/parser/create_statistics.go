// Package parser - create_statistics.go implements T-SQL CREATE/UPDATE/DROP STATISTICS parsing.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// statisticsWithOptions: UPDATE/CREATE STATISTICS WITH ... option names.
var statisticsWithOptions = newOptionSet(kwFULL, kwNOCOUNT).withIdents(
	"FULLSCAN", "SAMPLE", "RESAMPLE", "NORECOMPUTE", "INCREMENTAL",
	"MAXDOP", "ALL", "STATS_STREAM", "ROWCOUNT", "PAGECOUNT",
	"PERSIST_SAMPLE_PERCENT", "AUTO_DROP",
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

	// Column list. Graph pseudo-columns are valid here (engine-verified:
	// CREATE STATISTICS st ON Person ($node_id) executes).
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == tokPSEUDOCOL {
				if !graphPseudoColumns[strings.ToLower(p.cur.Str)] {
					return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
				}
				cols = append(cols, &nodes.String{Str: p.cur.Str})
				p.advance()
			} else if p.isIdentLike() {
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
//
// updateStatisticsStatement (TSql160.g:13430); the options are
// updateStatisticsStatementWithOption (TSql160.g:13472), checked as
// parseUpdateStatsOptions describes.
func (p *Parser) parseUpdateStatisticsStmt() (*nodes.UpdateStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.UpdateStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Table = table

	// Optional statistics name or list
	nameCount := 0
	if p.cur.Type == '(' {
		p.advance()
		var names []nodes.Node
		for {
			name, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			names = append(names, &nodes.String{Str: name})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Names = &nodes.List{Items: names}
		// Also set Name to first for backward compatibility
		stmt.Name = names[0].(*nodes.String).Str
		nameCount = len(names)
	} else if p.isIdentLike() && p.cur.Type != kwWITH {
		stmt.Name = p.cur.Str
		p.advance()
		nameCount = 1
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance()
		opts, err := p.parseUpdateStatsOptions(nameCount)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// updateStatsOptionGroup sorts UPDATE STATISTICS options into the groups SQL
// Server allows one of (Msg 1052): how the rows are read (FULLSCAN, SAMPLE,
// RESAMPLE), or statistics given outright (STATS_STREAM, ROWCOUNT,
// PAGECOUNT), which also excludes ALL, COLUMNS and INDEX.
var updateStatsOptionGroup = map[string]string{
	"FULLSCAN": "scan", "SAMPLE": "scan", "RESAMPLE": "scan",
	"STATS_STREAM": "stream", "ROWCOUNT": "stream", "PAGECOUNT": "stream",
	"ALL": "target", "COLUMNS": "target", "INDEX": "target",
}

// parseUpdateStatsOptions parses the options after UPDATE STATISTICS ...
// WITH, each given once (Msg 1039):
//
//	FULLSCAN | SAMPLE integer { PERCENT | ROWS } | RESAMPLE [ ON PARTITIONS ( integer [ TO integer ] [ ,...n ] ) ]
//	| STATS_STREAM = binary | ROWCOUNT = integer | PAGECOUNT = integer
//	| ALL | COLUMNS | INDEX | NORECOMPUTE | MAXDOP = integer
//	| { INCREMENTAL | AUTO_DROP | PERSIST_SAMPLE_PERCENT } = { ON | OFF }
//
// SQL Server refuses two options of a group (updateStatsOptionGroup), a
// stream option beside a scan or target one, PERSIST_SAMPLE_PERCENT without
// FULLSCAN or SAMPLE (Msg 153), STATS_STREAM without exactly one statistics
// name (Msg 1092), and decimal or negative numbers.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/update-statistics-transact-sql
func (p *Parser) parseUpdateStatsOptions(nameCount int) (*nodes.List, error) {
	var opts []nodes.Node
	seen := map[string]bool{}
	groups := map[string]string{}
	for {
		optLoc := p.pos()
		if !p.isKeywordOrIdent() || p.isDelimited(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		switch name {
		case "FULLSCAN", "ALL", "COLUMNS", "INDEX", "NORECOMPUTE":
		case "SAMPLE":
			if p.cur.Type != tokICONST || execIsBinary(p.cur) {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != kwPERCENT && p.cur.Type != kwROWS {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		case "RESAMPLE":
			if p.cur.Type == kwON {
				p.advance() // ON
				if err := p.execExpectWord("PARTITIONS"); err != nil {
					return nil, err
				}
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				for {
					// statisticsPartitionRange (TSql160.g:13328)
					if p.cur.Type != tokICONST || execIsBinary(p.cur) {
						return nil, p.syntaxErrorAtCur()
					}
					p.advance()
					if _, ok := p.match(kwTO); ok {
						if p.cur.Type != tokICONST || execIsBinary(p.cur) {
							return nil, p.syntaxErrorAtCur()
						}
						p.advance()
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
		case "STATS_STREAM":
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if (p.cur.Type != tokICONST && p.cur.Type != tokFCONST) || !execIsBinary(p.cur) {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if nameCount != 1 {
				return nil, p.newParseError(optLoc, "STATS_STREAM needs exactly one statistics name")
			}
		case "ROWCOUNT", "PAGECOUNT", "MAXDOP":
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if p.cur.Type != tokICONST || execIsBinary(p.cur) {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		case "INCREMENTAL", "AUTO_DROP", "PERSIST_SAMPLE_PERCENT":
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if p.cur.Type != kwON && p.cur.Type != kwOFF {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		default:
			return nil, p.newParseError(optLoc, "'"+name+"' is not a recognized UPDATE STATISTICS option")
		}
		if seen[name] {
			return nil, p.newParseError(optLoc, "option '"+name+"' is specified more than once")
		}
		seen[name] = true
		if g := updateStatsOptionGroup[name]; g != "" {
			if prev, ok := groups[g]; ok {
				return nil, p.newParseError(optLoc, "conflicting UPDATE STATISTICS options "+prev+" and "+name)
			}
			groups[g] = name
		}
		opts = append(opts, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, hasScan := groups["scan"]
	_, hasTarget := groups["target"]
	if _, ok := groups["stream"]; ok && (hasScan || hasTarget) {
		return nil, p.newParseError(p.prev.Loc, "conflicting UPDATE STATISTICS options")
	}
	if seen["PERSIST_SAMPLE_PERCENT"] && !seen["FULLSCAN"] && !seen["SAMPLE"] {
		return nil, p.newParseError(p.prev.Loc, "invalid usage of the option PERSIST_SAMPLE_PERCENT in the UPDATE STATISTICS statement")
	}
	return &nodes.List{Items: opts}, nil
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseStatisticsWithOptions parses the WITH clause of CREATE/UPDATE STATISTICS.
func (p *Parser) parseStatisticsWithOptions() *nodes.List {
	var opts []nodes.Node
	for {
		if p.isValidOption(statisticsWithOptions) {
			opt := strings.ToUpper(p.cur.Str)
			isSample := p.cur.Type == kwSAMPLE
			isResample := p.cur.Type == kwRESAMPLE
			p.advance()
			// Handle SAMPLE number { PERCENT | ROWS }.
			if isSample {
				p.parseExpr() // number
				if p.cur.Type == kwPERCENT || p.cur.Type == kwROWS {
					p.advance()
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
				if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
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
