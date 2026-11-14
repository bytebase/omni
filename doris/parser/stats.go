package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// ANALYZE statement
// ---------------------------------------------------------------------------

// parseAnalyze parses ANALYZE as both of Doris's parsers take it
// (DorisParser.g4:752-755; the CUP parser's analyze_stmt,
// DorisLegacyParser.g4:888):
//
//	ANALYZE TABLE multipartIdentifier partitionSpec? identifierList?
//	    (WITH analyzeProperties)* propertyClause?
//	ANALYZE DATABASE (identifier '.')? identifier
//	    (WITH analyzeProperties)* propertyClause?
//
// with the CUP parser's 'k' = 'v' properties.
//
// The ANALYZE keyword has already been consumed; cur is the next token.
func (p *Parser) parseAnalyze() (ast.Node, error) {
	startLoc := p.prev.Loc
	stmt := &ast.AnalyzeStmt{}

	switch p.cur.Kind {
	case kwDATABASE:
		p.advance() // consume DATABASE
		stmt.TargetType = "DATABASE"
		name, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = &ast.ObjectName{Parts: []string{name}, Loc: loc}
		if p.cur.Kind == int('.') {
			p.advance() // consume '.'
			second, loc2, err := p.parseIdentifierQualified()
			if err != nil {
				return nil, err
			}
			stmt.Target.Parts = append(stmt.Target.Parts, second)
			stmt.Target.Loc.End = loc2.End
		}

	case kwTABLE:
		p.advance() // consume TABLE
		stmt.TargetType = "TABLE"
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
		if stmt.Partitions, stmt.TemporaryPartition, err = p.parsePartitionSpec(); err != nil {
			return nil, err
		}
		if p.cur.Kind == int('(') {
			if stmt.Columns, err = p.parseLoadIdentifierList(); err != nil {
				return nil, err
			}
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	for p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		prop, _, err := p.parseAnalyzeWith()
		if err != nil {
			return nil, err
		}
		stmt.Properties = append(stmt.Properties, prop)
	}
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = append(stmt.Properties, props...)
	}

	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAnalyzeWith parses one analyzeProperties after WITH
// (DorisParser.g4:781; the CUP parser's with_analysis_properties):
//
//	SYNC | INCREMENTAL | FULL | SQL | HISTOGRAM
//	| SAMPLE (ROWS | PERCENT) INTEGER_VALUE
//	| BUCKETS INTEGER_VALUE | PERIOD INTEGER_VALUE | CRON STRING_LITERAL
//
// WITH has already been consumed. It returns Property{Key: the words,
// Value: the value or ""}.
func (p *Parser) parseAnalyzeWith() (*ast.Property, ast.Loc, error) {
	startLoc := p.cur.Loc
	key := strings.ToUpper(p.cur.Str)
	valueKind := 0
	switch p.cur.Kind {
	case kwSYNC, kwINCREMENTAL, kwFULL, kwSQL, kwHISTOGRAM:
		p.advance()
	case kwSAMPLE:
		p.advance() // consume SAMPLE
		if p.cur.Kind != kwROWS && p.cur.Kind != kwPERCENT {
			return nil, ast.Loc{}, p.syntaxErrorAtCur()
		}
		key += " " + strings.ToUpper(p.advance().Str)
		valueKind = tokInt
	case kwBUCKETS, kwPERIOD:
		p.advance()
		valueKind = tokInt
	case kwCRON:
		p.advance()
		valueKind = tokString
	default:
		return nil, ast.Loc{}, p.syntaxErrorAtCur()
	}
	value := ""
	if valueKind != 0 {
		if p.cur.Kind != valueKind {
			return nil, ast.Loc{}, p.syntaxErrorAtCur()
		}
		tok := p.advance()
		value = tok.Str
		if tok.Kind == tokInt {
			value = p.sourceText(tok.Loc)
		}
	}
	loc := ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return &ast.Property{Key: key, Value: value, Loc: loc}, loc, nil
}

// ---------------------------------------------------------------------------
// SHOW variants: SHOW ANALYZE, SHOW STATS, SHOW CONSTRAINTS
// ---------------------------------------------------------------------------

// parseShowAnalyze parses SHOW [ALL|QUEUED] ANALYZE [TASK STATUS job_id | JOB? ...]
//
// The SHOW keyword has already been consumed. cur may be ALL, QUEUED, or
// ANALYZE directly.
func (p *Parser) parseShowAnalyze(startLoc ast.Loc, all, queued bool) (ast.Node, error) {
	// cur is ANALYZE — consume it.
	p.advance()

	stmt := &ast.ShowAnalyzeStmt{
		All:    all,
		Queued: queued,
	}
	endLoc := p.prev.Loc

	// SHOW ANALYZE TASK STATUS job_id
	if p.cur.Kind == kwTASK {
		p.advance() // consume TASK
		// Consume optional STATUS keyword
		if p.cur.Kind == kwSTATUS {
			p.advance()
		}
		stmt.IsTask = true
		if p.cur.Kind == tokInt {
			stmt.JobID = p.cur.Ival
			endLoc = p.cur.Loc
			p.advance()
		}
		stmt.Loc = startLoc.Merge(endLoc)
		return stmt, nil
	}

	// Optional JOB keyword (ignored, consumed)
	if p.cur.Kind == kwJOB {
		p.advance()
	}

	// Optional: job_id (integer), FOR table, bare table name, or LIKE pattern
	switch p.cur.Kind {
	case tokInt:
		stmt.JobID = p.cur.Ival
		endLoc = p.cur.Loc
		p.advance()

	case kwFOR:
		p.advance() // consume FOR
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.For = target
		endLoc = target.Loc

	case kwLIKE:
		p.advance() // consume LIKE
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Like = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()

	default:
		// Bare identifier used as table target (e.g., SHOW ANALYZE test1 WHERE ...)
		if isIdentifierToken(p.cur.Kind) {
			target, err := p.parseMultipartIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.For = target
			endLoc = target.Loc
		}
	}

	// Optional WHERE clause
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		endLoc = ast.NodeLoc(where)
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseShowStats parses SHOW [COLUMN|TABLE|INDEX|PARTITION] STATS [target] [WHERE...]
//
// The SHOW keyword has already been consumed. The optional type keyword (COLUMN,
// TABLE, etc.) may or may not be present; the STATS keyword must follow.
// statsType is the already-consumed type qualifier (empty string if absent).
func (p *Parser) parseShowStats(startLoc ast.Loc, statsType string) (ast.Node, error) {
	// cur is STATS — consume it
	p.advance()

	stmt := &ast.ShowStatsStmt{
		Type: statsType,
	}
	endLoc := p.prev.Loc

	// Optional target: a name, or SHOW TABLE STATS's table id.
	if statsType == "TABLE" && p.cur.Kind == tokInt {
		stmt.TableID = p.cur.Ival
		endLoc = p.advance().Loc
	} else if isIdentifierToken(p.cur.Kind) {
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
		endLoc = target.Loc
	}

	// Optional WHERE clause
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		endLoc = ast.NodeLoc(where)
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseShowConstraints parses SHOW CONSTRAINTS FROM table
//
// SHOW has already been consumed; cur is CONSTRAINTS.
func (p *Parser) parseShowConstraints(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume CONSTRAINTS

	stmt := &ast.ShowConstraintsStmt{}
	endLoc := p.prev.Loc

	// Expect FROM table
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM / IN
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Table = target
		endLoc = target.Loc
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseShow is the top-level SHOW dispatcher.
//
// (Note: parseShow is defined in show.go (T7.3); analyze/stats/constraints
//  routing is wired into that function. The standalone parseShow here was
//  removed during the T8.3 rebase onto main.)

// ---------------------------------------------------------------------------
// DROP STATS / DROP EXPIRED STATS / DROP CACHED STATS
// ---------------------------------------------------------------------------

// parseDropStats parses (DorisParser.g4:760-763; the CUP parser's
// drop_stmt, DorisLegacyParser.g4:983-986):
//
//	DROP STATS multipartIdentifier identifierList? partitionSpec?
//	DROP CACHED STATS multipartIdentifier
//	DROP EXPIRED STATS
//	DROP ANALYZE JOB INTEGER_VALUE
//
// DROP and the words up to STATS (or ANALYZE) have been consumed; variant
// is "", "CACHED", "EXPIRED" or "ANALYZE JOB".
func (p *Parser) parseDropStats(startLoc ast.Loc, variant string) (ast.Node, error) {
	stmt := &ast.DropStatsStmt{Variant: variant}
	var err error
	switch variant {
	case "EXPIRED":
	case "ANALYZE JOB":
		if _, err := p.expect(kwJOB); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.JobID = p.advance().Ival
	default:
		if stmt.Target, err = p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
		if variant == "" {
			if p.cur.Kind == int('(') {
				if stmt.Columns, err = p.parseLoadIdentifierList(); err != nil {
					return nil, err
				}
			}
			if stmt.Partitions, stmt.TemporaryPartition, err = p.parsePartitionSpec(); err != nil {
				return nil, err
			}
		}
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

// parseKillAnalyze parses KILL ANALYZE job_id; the id is required
// (DorisParser.g4:766 #killAnalyzeJob, DorisLegacyParser.g4:1888
// kill_analysis_job_stmt).
//
// KILL has already been consumed; cur is ANALYZE.
func (p *Parser) parseKillAnalyze(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ANALYZE

	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	stmt := &ast.KillAnalyzeStmt{JobID: p.cur.Ival, Loc: startLoc.Merge(p.cur.Loc)}
	p.advance()
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER TABLE ADD/DROP CONSTRAINT
// ---------------------------------------------------------------------------

// parseAddConstraint parses:
//
//	ALTER TABLE name ADD CONSTRAINT cname PRIMARY KEY (cols)
//	ALTER TABLE name ADD CONSTRAINT cname UNIQUE (cols)
//	ALTER TABLE name ADD CONSTRAINT cname FOREIGN KEY (cols) REFERENCES ref (cols)
//
// Each column list is an identifierList, and the referenced columns are not
// optional (DorisParser.g4:954 constraint).
//
// ALTER TABLE and the table name have been parsed; table is the table ObjectName.
// ADD has already been consumed; cur is CONSTRAINT.
func (p *Parser) parseAddConstraint(startLoc ast.Loc, table *ast.ObjectName) (ast.Node, error) {
	p.advance() // consume CONSTRAINT

	stmt := &ast.AddConstraintStmt{Table: table}
	endLoc := p.prev.Loc

	// constraint name
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc = nameLoc

	switch p.cur.Kind {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		stmt.Type = "PRIMARY KEY"
		cols, colsEnd, err := p.alterParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		endLoc = colsEnd

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		stmt.Type = "UNIQUE"
		cols, colsEnd, err := p.alterParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		endLoc = colsEnd

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		stmt.Type = "FOREIGN KEY"
		cols, _, err := p.alterParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols

		if _, err := p.expect(kwREFERENCES); err != nil {
			return nil, err
		}
		refTable, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.RefTable = refTable
		endLoc = refTable.Loc

		refCols, refColsEnd, err := p.alterParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.RefColumns = refCols
		endLoc = refColsEnd

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropConstraint parses:
//
//	ALTER TABLE name DROP CONSTRAINT cname
//
// DROP and CONSTRAINT have not yet been consumed; cur is DROP.
// table is the already-parsed table ObjectName.
func (p *Parser) parseDropConstraint(startLoc ast.Loc, table *ast.ObjectName) (ast.Node, error) {
	p.advance() // consume DROP
	p.advance() // consume CONSTRAINT

	stmt := &ast.DropConstraintStmt{Table: table}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(nameLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helper: parseParenIdentifierList — parses (id1, id2, ...) returning bare names
// ---------------------------------------------------------------------------

// parseParenIdentifierList parses a parenthesised comma-separated list of
// identifiers: (col1, col2, ...).
// Returns the names and the Loc of the closing ')'.
func (p *Parser) parseParenIdentifierList() ([]string, ast.Loc, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, ast.Loc{}, err
	}

	var names []string
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, ast.Loc{}, err
		}
		names = append(names, name)
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.Loc{}, err
	}
	return names, closeTok.Loc, nil
}
