// Package parser - create_statistics.go implements T-SQL CREATE/UPDATE/DROP STATISTICS parsing.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// statisticsPseudoColumns are the pseudo-columns SQL Server takes in a
// CREATE STATISTICS column list: the graph ones, $IDENTITY and $ROWGUID.
var statisticsPseudoColumns = map[string]bool{
	"$identity": true,
	"$rowguid":  true,
	"$node_id":  true,
	"$edge_id":  true,
	"$from_id":  true,
	"$to_id":    true,
}

// parseCreateStatisticsStmt parses a CREATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-statistics-transact-sql
//
//	createStatisticsStatement (TSql160.g:13230)
//	    : Create Statistics identifier On schemaObjectThreePartName identifierColumnList
//	      ( filterClause )? ( With createStatisticsStatementWithOption ( , ... )* )?
func (p *Parser) parseCreateStatisticsStmt() (*nodes.CreateStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.CreateStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	stmt.Name = name

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.indexTargetTable(3)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// ( column [ ,...n ] ): one-part names or pseudo-columns (engine-verified:
	// CREATE STATISTICS st ON Person ($node_id) executes; a dotted name is
	// Msg 4146).
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		if p.cur.Type == tokPSEUDOCOL {
			if !statisticsPseudoColumns[strings.ToLower(p.cur.Str)] {
				return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
			}
			col := &nodes.String{Str: p.cur.Str}
			p.advance()
			return col, nil
		}
		col, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		return &nodes.String{Str: col}, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Columns = &nodes.List{Items: cols}

	// WHERE filter predicate (optional)
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.Where, err = p.indexFilterExpression()
		if err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options, err = p.parseStatisticsWithOptions(false, 0)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseUpdateStatisticsStmt parses an UPDATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/update-statistics-transact-sql
//
//	updateStatisticsStatement (TSql160.g:13430)
//	    : Update Statistics schemaObjectThreePartName ( columnNameList | identifier )?
//	      ( With updateStatisticsStatementWithOption ( , ... )* )?
func (p *Parser) parseUpdateStatisticsStmt() (*nodes.UpdateStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.UpdateStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	table, err := p.indexTargetTable(3)
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// Optional statistics name or list
	nameCount := 0
	if p.cur.Type == '(' {
		p.advance()
		names, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			n, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			return &nodes.String{Str: n}, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Names = &nodes.List{Items: names}
		// Also set Name to first for backward compatibility
		stmt.Name = names[0].(*nodes.String).Str
		nameCount = len(names)
	} else if n, ok := p.parseIdentifier(); ok {
		stmt.Name = n
		nameCount = 1
	}

	if p.cur.Type == kwWITH {
		p.advance()
		stmt.Options, err = p.parseStatisticsWithOptions(true, nameCount)
		if err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropStatisticsStmt parses a DROP STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-statistics-transact-sql
//
//	dropStatisticsStatement (TSql160.g:15030) : Statistics statisticsDropObject ( , ... )*
//	statisticsDropObject (TSql160.g:15045) : [ schema . ] table . statistics_name
//
// A name without its table is SQL46038; one with a database is Msg 166.
func (p *Parser) parseDropStatisticsStmt() (*nodes.DropStatisticsStmt, error) {
	loc := p.pos()
	// STATISTICS keyword already consumed by caller

	stmt := &nodes.DropStatisticsStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var names []nodes.Node
	for {
		start := p.pos()
		first, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		parts := []string{first}
		for p.cur.Type == '.' {
			p.advance()
			part, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			parts = append(parts, part)
		}
		if len(parts) < 2 || len(parts) > 3 {
			return nil, p.newParseError(start, "DROP STATISTICS needs table.statistics_name")
		}
		names = append(names, &nodes.String{Str: strings.Join(parts, ".")})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Names = &nodes.List{Items: names}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseStatisticsWithOptions parses the options after WITH in CREATE or
// UPDATE STATISTICS. Each becomes a String holding the option's name.
//
//	createStatisticsStatementWithOption (TSql160.g:13412)
//	    : INCREMENTAL = ( On | Off ) | AUTO_DROP = ( On | Off )
//	    | PERSIST_SAMPLE_PERCENT = ( On | Off ) | SAMPLE integer ( ROWS | PERCENT )
//	    | STATS_STREAM = binary | FULLSCAN | NORECOMPUTE | MAXDOP = integer
//	updateStatisticsStatementWithOption (TSql160.g:13472) adds
//	    RESAMPLE ( On PARTITIONS ( statisticsPartitionRange , ... ) )?
//	    | PAGECOUNT = integer | ROWCOUNT = integer | ALL | COLUMNS | INDEX
//
// SQL Server refuses, while it parses, an option given twice (Msg 1039),
// two of FULLSCAN, SAMPLE, RESAMPLE and STATS_STREAM (1052), ROWCOUNT or
// PAGECOUNT beside a scan, PERSIST_SAMPLE_PERCENT beside RESAMPLE or
// STATS_STREAM (1052) or without FULLSCAN or SAMPLE (153), two of ALL,
// COLUMNS and INDEX (1052), and STATS_STREAM in UPDATE STATISTICS without
// exactly one statistics name (1092).
func (p *Parser) parseStatisticsWithOptions(update bool, nameCount int) (*nodes.List, error) {
	seen := map[string]bool{}
	scan := "" // FULLSCAN, SAMPLE, RESAMPLE or STATS_STREAM
	target := ""
	start := p.pos()
	var opts []nodes.Node
	for {
		word := p.indexOptionWord()
		optLoc := p.pos()
		if seen[word] {
			return nil, p.newParseError(optLoc, fmt.Sprintf("option %s is specified more than once", word))
		}
		switch word {
		case "FULLSCAN", "NORECOMPUTE":
			p.advance()
		case "SAMPLE":
			p.advance()
			n, err := p.indexExpectInteger()
			if err != nil {
				return nil, err
			}
			if p.cur.Type != kwPERCENT && p.cur.Type != kwROWS {
				return nil, p.unexpectedToken()
			}
			if v, convErr := strconv.Atoi(n); p.cur.Type == kwPERCENT && (convErr != nil || v > 100) {
				// Msg 1031: percent values must be between 0 and 100.
				return nil, p.newParseError(p.prev.Loc, "percent values must be between 0 and 100")
			}
			p.advance()
		case "STATS_STREAM":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if !indexTokIsBinary(p.cur) {
				return nil, p.unexpectedToken()
			}
			p.advance()
			if update && nameCount != 1 {
				return nil, p.newParseError(optLoc, "STATS_STREAM needs exactly one statistics name")
			}
		case "INCREMENTAL", "AUTO_DROP", "PERSIST_SAMPLE_PERCENT":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if _, err := p.indexExpectOnOff(); err != nil {
				return nil, err
			}
		case "MAXDOP":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if _, err := p.indexExpectInteger(); err != nil {
				return nil, err
			}
		default:
			if !update {
				return nil, p.unexpectedToken()
			}
			switch word {
			case "RESAMPLE":
				p.advance()
				if err := p.statisticsResamplePartitions(); err != nil {
					return nil, err
				}
			case "ROWCOUNT", "PAGECOUNT":
				p.advance()
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				if _, err := p.indexExpectInteger(); err != nil {
					return nil, err
				}
			case "ALL", "COLUMNS", "INDEX":
				if target != "" {
					return nil, p.newParseError(optLoc, fmt.Sprintf("conflicting options %s and %s", target, word))
				}
				target = word
				p.advance()
			default:
				return nil, p.unexpectedToken()
			}
		}
		switch word {
		case "FULLSCAN", "SAMPLE", "RESAMPLE", "STATS_STREAM":
			if scan != "" {
				return nil, p.newParseError(optLoc, fmt.Sprintf("conflicting options %s and %s", scan, word))
			}
			scan = word
		}
		seen[word] = true
		opts = append(opts, &nodes.String{Str: word})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if seen["PERSIST_SAMPLE_PERCENT"] && scan != "FULLSCAN" && scan != "SAMPLE" {
		return nil, p.newParseError(start, "PERSIST_SAMPLE_PERCENT needs FULLSCAN or SAMPLE")
	}
	if (seen["ROWCOUNT"] || seen["PAGECOUNT"]) && scan != "" && scan != "STATS_STREAM" {
		return nil, p.newParseError(start, "ROWCOUNT and PAGECOUNT conflict with "+scan)
	}
	return &nodes.List{Items: opts}, nil
}

// statisticsResamplePartitions parses the optional ON PARTITIONS ( range
// [ , ...n ] ) after RESAMPLE; statisticsPartitionRange (TSql160.g:13328)
// is integer ( TO integer )?.
func (p *Parser) statisticsResamplePartitions() error {
	if p.cur.Type != kwON {
		return nil
	}
	p.advance()
	if p.cur.Type != kwPARTITIONS {
		return p.unexpectedToken()
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return err
	}
	_, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		from, err := p.indexExpectInteger()
		if err != nil {
			return nil, err
		}
		if _, ok := p.match(kwTO); ok {
			if _, err := p.indexExpectInteger(); err != nil {
				return nil, err
			}
		}
		return &nodes.String{Str: from}, nil
	})
	if err != nil {
		return err
	}
	_, err = p.expect(')')
	return err
}
