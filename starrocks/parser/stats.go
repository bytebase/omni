package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ANALYZE statement
// ---------------------------------------------------------------------------

// parseAnalyze parses:
//
//	ANALYZE DATABASE name
//	ANALYZE TABLE name [(col1, col2, ...)]
//	ANALYZE PROFILE
//
// followed by optional modifiers:
//
//	WITH SAMPLE PERCENT n
//	WITH SAMPLE ROWS n
//	WITH SYNC
//	WITH INCREMENTAL
//	PROPERTIES(...)
//
// The ANALYZE keyword has already been consumed; cur is the next token.
func (p *Parser) parseAnalyze() (ast.Node, error) {
	startLoc := p.prev.Loc

	stmt := &ast.AnalyzeStmt{}
	endLoc := startLoc

	switch p.cur.Kind {
	case kwDATABASE, kwSCHEMA:
		p.advance() // consume DATABASE/SCHEMA
		stmt.TargetType = "DATABASE"
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
		endLoc = target.Loc

	case kwTABLE:
		p.advance() // consume TABLE
		stmt.TargetType = "TABLE"
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
		endLoc = target.Loc

		// Optional column list: (col1, col2, ...)
		if p.cur.Kind == int('(') {
			cols, colsEnd, err := p.parseParenIdentifierList()
			if err != nil {
				return nil, err
			}
			stmt.Columns = cols
			endLoc = colsEnd
		}

	case kwPROFILE:
		p.advance() // consume PROFILE
		stmt.TargetType = "PROFILE"
		endLoc = p.prev.Loc

	default:
		// Bare ANALYZE with no target type — treat as ANALYZE TABLE (Doris also
		// allows ANALYZE name directly for table-level analyze).
		if isIdentifierToken(p.cur.Kind) {
			stmt.TargetType = "TABLE"
			target, err := p.parseMultipartIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Target = target
			endLoc = target.Loc

			// Optional column list
			if p.cur.Kind == int('(') {
				cols, colsEnd, err := p.parseParenIdentifierList()
				if err != nil {
					return nil, err
				}
				stmt.Columns = cols
				endLoc = colsEnd
			}
		}
	}

	// Optional WITH modifiers and PROPERTIES
	for {
		if p.cur.Kind == kwWITH {
			p.advance() // consume WITH
			prop, propEnd, err := p.parseAnalyzeWith()
			if err != nil {
				return nil, err
			}
			stmt.Properties = append(stmt.Properties, prop)
			endLoc = propEnd
		} else if p.cur.Kind == kwPROPERTIES {
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			stmt.Properties = append(stmt.Properties, props...)
			if len(props) > 0 {
				endLoc = ast.NodeLoc(props[len(props)-1])
			}
		} else {
			break
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseAnalyzeWith parses a single WITH modifier after ANALYZE:
//
//	WITH SAMPLE PERCENT n
//	WITH SAMPLE ROWS n
//	WITH SYNC
//	WITH INCREMENTAL
//
// WITH has already been consumed; cur is the modifier keyword.
// Returns a synthetic Property{Key: modifier, Value: n} for value-bearing
// modifiers or Property{Key: modifier, Value: ""} for flags.
func (p *Parser) parseAnalyzeWith() (*ast.Property, ast.Loc, error) {
	startLoc := p.cur.Loc

	switch p.cur.Kind {
	case kwSAMPLE:
		p.advance() // consume SAMPLE
		switch p.cur.Kind {
		case kwPERCENT:
			p.advance() // consume PERCENT
			if p.cur.Kind != tokInt && p.cur.Kind != tokFloat {
				return nil, ast.Loc{}, p.syntaxErrorAtCur()
			}
			val := p.cur.Str
			endLoc := p.cur.Loc
			p.advance()
			return &ast.Property{Key: "SAMPLE PERCENT", Value: val, Loc: startLoc.Merge(endLoc)}, endLoc, nil
		case kwROWS:
			p.advance() // consume ROWS
			if p.cur.Kind != tokInt {
				return nil, ast.Loc{}, p.syntaxErrorAtCur()
			}
			val := p.cur.Str
			endLoc := p.cur.Loc
			p.advance()
			return &ast.Property{Key: "SAMPLE ROWS", Value: val, Loc: startLoc.Merge(endLoc)}, endLoc, nil
		default:
			return nil, ast.Loc{}, p.syntaxErrorAtCur()
		}

	case kwSYNC:
		endLoc := p.cur.Loc
		p.advance()
		return &ast.Property{Key: "SYNC", Value: "", Loc: startLoc.Merge(endLoc)}, endLoc, nil

	case kwINCREMENTAL:
		endLoc := p.cur.Loc
		p.advance()
		return &ast.Property{Key: "INCREMENTAL", Value: "", Loc: startLoc.Merge(endLoc)}, endLoc, nil

	default:
		// Unknown WITH modifier — consume as raw identifier to avoid hard failure
		if isIdentifierToken(p.cur.Kind) {
			key := strings.ToUpper(p.cur.Str)
			endLoc := p.cur.Loc
			p.advance()
			return &ast.Property{Key: key, Value: "", Loc: startLoc.Merge(endLoc)}, endLoc, nil
		}
		return nil, ast.Loc{}, p.syntaxErrorAtCur()
	}
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

// parseDropStats parses:
//
//	DROP STATS target [(col1, col2)]
//	DROP EXPIRED STATS target
//	DROP CACHED STATS target
//
// DROP has already been consumed; variant is "" | "EXPIRED" | "CACHED".
// The STATS keyword has also been consumed by the caller for EXPIRED/CACHED.
// For plain DROP STATS, the caller consumed STATS; target is next.
func (p *Parser) parseDropStats(startLoc ast.Loc, variant string) (ast.Node, error) {
	stmt := &ast.DropStatsStmt{Variant: variant}
	endLoc := startLoc

	// target table
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target
	endLoc = target.Loc

	// Optional column list for plain DROP STATS
	if variant == "" && p.cur.Kind == int('(') {
		cols, colsEnd, err := p.parseParenIdentifierList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		endLoc = colsEnd
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

// parseKillAnalyze parses KILL ANALYZE job_id.
//
// KILL has already been consumed; cur is ANALYZE.
func (p *Parser) parseKillAnalyze(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume ANALYZE

	stmt := &ast.KillAnalyzeStmt{}
	endLoc := p.prev.Loc

	if p.cur.Kind == tokInt {
		stmt.JobID = p.cur.Ival
		endLoc = p.cur.Loc
		p.advance()
	}

	stmt.Loc = startLoc.Merge(endLoc)
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
