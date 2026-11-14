package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ANALYZE statement
// ---------------------------------------------------------------------------

// parseAnalyze parses the ANALYZE statements of StarRocks.g4:
//
//	ANALYZE (FULL | SAMPLE)? TABLE qualifiedName ('(' qualifiedName (',' qualifiedName)* ')')?
//	    (WITH (SYNC | ASYNC) MODE)? properties?                         (:1297)
//	ANALYZE TABLE qualifiedName UPDATE HISTOGRAM ON qualifiedName (',' qualifiedName)*
//	    (WITH INTEGER_VALUE BUCKETS)? properties? (WITH (SYNC | ASYNC) MODE)?   (:1307)
//	ANALYZE TABLE qualifiedName DROP HISTOGRAM ON qualifiedName (',' qualifiedName)*   (:1318)
//	ANALYZE PROFILE FROM string (',' INTEGER_VALUE)*                   (:1351)
//
// The ANALYZE keyword has already been consumed; cur is the next token.
func (p *Parser) parseAnalyze() (ast.Node, error) {
	startLoc := p.prev.Loc
	stmt := &ast.AnalyzeStmt{}

	if p.cur.Kind == kwPROFILE {
		p.advance() // consume PROFILE
		stmt.TargetType = "PROFILE"
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.ProfileID = p.advance().Str
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			if p.cur.Kind != tokInt {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.PlanNodeIDs = append(stmt.PlanNodeIDs, p.advance().Ival)
		}
		stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
		return stmt, nil
	}

	if err := p.parseAnalyzeTable(stmt, true); err != nil {
		return nil, err
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAnalyzeTable reads (FULL | SAMPLE)? TABLE qualifiedName and what
// follows it in an ANALYZE or a CREATE ANALYZE (StarRocks.g4:1297-1325):
// a column list, or UPDATE HISTOGRAM ON (or, unless job, DROP HISTOGRAM
// ON) a column list; then the clauses each form takes. A job (CREATE
// ANALYZE) has no WITH ... MODE outside a histogram.
func (p *Parser) parseAnalyzeTable(stmt *ast.AnalyzeStmt, analyze bool) error {
	if p.cur.Kind == kwFULL || p.cur.Kind == kwSAMPLE {
		stmt.Method = strings.ToUpper(p.advance().Str)
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return err
	}
	stmt.TargetType = "TABLE"
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	stmt.Target = target

	if stmt.Method == "" && (p.cur.Kind == kwUPDATE || (analyze && p.cur.Kind == kwDROP)) {
		// histogramStatement / dropHistogramStatement
		stmt.Histogram = strings.ToUpper(p.advance().Str)
		if _, err := p.expect(kwHISTOGRAM); err != nil {
			return err
		}
		if _, err := p.expect(kwON); err != nil {
			return err
		}
		if stmt.Columns, err = p.parseAnalyzeColumns(); err != nil {
			return err
		}
		if stmt.Histogram == "DROP" {
			return nil
		}
		if p.cur.Kind == kwWITH && p.peekNext().Kind == tokInt {
			p.advance() // consume WITH
			stmt.Buckets = p.advance().Ival
			if _, err := p.expect(kwBUCKETS); err != nil {
				return err
			}
		}
		if err := p.parseAnalyzeProperties(stmt); err != nil {
			return err
		}
		if analyze {
			return p.parseAnalyzeMode(stmt)
		}
		return nil
	}

	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		if stmt.Columns, err = p.parseAnalyzeColumns(); err != nil {
			return err
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
	}
	if analyze {
		if err := p.parseAnalyzeMode(stmt); err != nil {
			return err
		}
	}
	return p.parseAnalyzeProperties(stmt)
}

// parseAnalyzeColumns reads qualifiedName (',' qualifiedName)*.
func (p *Parser) parseAnalyzeColumns() ([]string, error) {
	var cols []string
	for {
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		cols = append(cols, strings.Join(name.Parts, "."))
		if p.cur.Kind != int(',') {
			return cols, nil
		}
		p.advance() // consume ','
	}
}

// parseAnalyzeMode reads an optional WITH (SYNC | ASYNC) MODE.
func (p *Parser) parseAnalyzeMode(stmt *ast.AnalyzeStmt) error {
	if p.cur.Kind != kwWITH {
		return nil
	}
	p.advance() // consume WITH
	switch {
	case p.cur.Kind == kwSYNC:
		stmt.Mode = "SYNC"
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "ASYNC"):
		stmt.Mode = "ASYNC"
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Kind != tokIdent || !strings.EqualFold(p.cur.Str, "MODE") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseAnalyzeProperties reads an optional properties clause.
func (p *Parser) parseAnalyzeProperties(stmt *ast.AnalyzeStmt) error {
	if p.cur.Kind != kwPROPERTIES {
		return nil
	}
	props, err := p.parseProperties()
	stmt.Properties = props
	return err
}

// parseCreateAnalyze parses createAnalyzeStatement (StarRocks.g4:1322):
//
//	CREATE ANALYZE (FULL | SAMPLE)? ALL properties?
//	| CREATE ANALYZE (FULL | SAMPLE)? DATABASE identifier properties?
//	| CREATE ANALYZE (FULL | SAMPLE)? TABLE qualifiedName ('(' qualifiedName (',' qualifiedName)* ')')? properties?
//	| CREATE histogramStatement
//
// CREATE has been consumed; cur is ANALYZE.
func (p *Parser) parseCreateAnalyze(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ANALYZE
	stmt := &ast.AnalyzeStmt{Create: true}
	method := ""
	if p.cur.Kind == kwFULL || p.cur.Kind == kwSAMPLE {
		method = strings.ToUpper(p.cur.Str)
	}
	switch k := p.peekPast(method != ""); k {
	case kwALL, kwDATABASE:
		if method != "" {
			p.advance()
		}
		stmt.Method = method
		p.advance() // consume ALL / DATABASE
		stmt.TargetType = "ALL"
		if k == kwDATABASE {
			stmt.TargetType = "DATABASE"
			name, loc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Target = &ast.ObjectName{Parts: []string{name}, Loc: loc}
		}
		if err := p.parseAnalyzeProperties(stmt); err != nil {
			return nil, err
		}
	default:
		if err := p.parseAnalyzeTable(stmt, false); err != nil {
			return nil, err
		}
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// peekPast returns the kind of cur, or of the token after it when skip is
// set.
func (p *Parser) peekPast(skip bool) TokenKind {
	if skip {
		return p.peekNext().Kind
	}
	return p.cur.Kind
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

	// Optional target name
	if isIdentifierToken(p.cur.Kind) {
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

// parseDropStats parses dropStatsStatement and dropAnalyzeJobStatement
// (StarRocks.g4:1303, :1329):
//
//	DROP STATS qualifiedName
//	DROP ANALYZE INTEGER_VALUE
//
// DROP and STATS (or ANALYZE) have been consumed; variant is "" or "ANALYZE".
func (p *Parser) parseDropStats(startLoc ast.Loc, variant string) (ast.Node, error) {
	stmt := &ast.DropStatsStmt{Variant: variant}
	if variant == "ANALYZE" {
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.JobID = p.advance().Ival
	} else {
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

// parseKillAnalyze parses KILL ANALYZE job_id; the id is required
// (StarRocks.g4:1345 killAnalyzeStatement).
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
		cols, colsEnd, err := p.parseParenIdentifierList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		endLoc = colsEnd

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		stmt.Type = "UNIQUE"
		cols, colsEnd, err := p.parseParenIdentifierList()
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
		cols, _, err := p.parseParenIdentifierList()
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

		if p.cur.Kind == int('(') {
			refCols, refColsEnd, err := p.parseParenIdentifierList()
			if err != nil {
				return nil, err
			}
			stmt.RefColumns = refCols
			endLoc = refColsEnd
		}

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
