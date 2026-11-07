package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateJob parses a CREATE JOB statement.
// On entry, CREATE has been consumed and cur is JOB.
//
// Syntax:
//
//	CREATE JOB [IF NOT EXISTS] job_name
//	    ON SCHEDULE { EVERY interval [STARTS ts] [ENDS ts] | AT ts }
//	    [COMMENT 'text']
//	    DO dml_statement
//
// or:
//
//	CREATE JOB [IF NOT EXISTS] job_name AS STREAMING
//	    [COMMENT 'text']
//	    DO dml_statement
func (p *Parser) parseCreateJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	stmt := &ast.CreateJobStmt{}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Determine job type: ON SCHEDULE or AS STREAMING
	if p.cur.Kind == kwAS {
		p.advance() // consume AS
		if _, err := p.expect(kwSTREAMING); err != nil {
			return nil, err
		}
		stmt.JobType = "STREAMING"
	} else if p.cur.Kind == kwON {
		p.advance() // consume ON
		if _, err := p.expect(kwSCHEDULE); err != nil {
			return nil, err
		}
		stmt.JobType = "SCHEDULE"

		sched, err := p.parseJobSchedule()
		if err != nil {
			return nil, err
		}
		stmt.Schedule = sched
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	// Optional COMMENT 'text'
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	// DO statement
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}

	doStmt, err := p.parseJobDoBody()
	if err != nil {
		return nil, err
	}
	stmt.DoStmt = doStmt

	stmt.Loc = startLoc.Merge(ast.NodeLoc(doStmt))
	return stmt, nil
}

// Doris 4.0 supportedDmlStatement includes MERGE and TRUNCATE.
// https://github.com/apache/doris/blob/branch-4.0/fe/fe-core/src/main/antlr4/org/apache/doris/nereids/DorisParser.g4
// parseJobDoBody parses the statement after a job's DO, which the grammar
// limits to a DML statement:
//
//	CREATE JOB ... DO supportedDmlStatement
//	supportedDmlStatement : explain? cte? INSERT ... | explain? cte? UPDATE ...
//	                      | explain? cte? DELETE ... | LOAD LABEL ...
//	                      | EXPORT TABLE ... | replayCommand
//
// (DorisParser.g4). It used to take any statement, so DO SELECT 1 INTO
// OUTFILE 'p', DO CREATE TABLE ... and DO SHOW TABLES parsed. The engine
// rejects each of those at the body's first token and names DELETE, DESC,
// DESCRIBE, EXPLAIN, EXPORT, INSERT, LOAD, PLAN, UPDATE and WITH as what it
// expects there; it also rejects DO EXPLAIN SELECT 1 and DO WITH ... SELECT,
// where the query after the prefix is not DML (container-verified).
//
// ALTER JOB ... DO takes the same body. The pinned engine has no ALTER JOB at
// all, so that statement form is a divergence of its own; holding its body to
// the same rule keeps it from carrying what CREATE JOB may not.
func (p *Parser) parseJobDoBody() (ast.Node, error) {
	start := p.cur.Loc
	switch p.cur.Kind {
	case kwINSERT, kwUPDATE, kwDELETE, kwLOAD, kwEXPORT,
		kwWITH, kwEXPLAIN, kwDESC, kwDESCRIBE, kwPLAN, kwMERGE, kwTRUNCATE:
	default:
		return nil, p.syntaxErrorAtCur()
	}
	body, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if !isJobDML(body) {
		// The prefix took a query that is not DML: EXPLAIN SELECT ...,
		// WITH ... SELECT ..., DESC t. Point at the query when there is one.
		at := start
		if explain, ok := body.(*ast.ExplainStmt); ok && explain.Query != nil {
			at = ast.NodeLoc(explain.Query)
		}
		return nil, &ParseError{
			Position: at.Start, End: at.End,
			Message: "a job's DO takes INSERT, UPDATE, DELETE, LOAD or EXPORT",
		}
	}
	return body, nil
}

// isJobDML reports whether n is a statement a job's DO may run: one of the
// grammar's supportedDmlStatement forms that this parser builds.
func isJobDML(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.LoadDataStmt, *ast.ExportStmt, *ast.MergeStmt, *ast.TruncateTableStmt:
		return true
	case *ast.ExplainStmt:
		// explain? applies to INSERT, UPDATE and DELETE only.
		switch n.Query.(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			return true
		}
	}
	return false
}

// parseJobSchedule parses the schedule portion after ON SCHEDULE:
//
//	EVERY interval [STARTS timestamp] [ENDS timestamp]
//	AT timestamp
//
// cur is at EVERY or AT on entry.
func (p *Parser) parseJobSchedule() (*ast.JobSchedule, error) {
	startLoc := p.cur.Loc
	sched := &ast.JobSchedule{Loc: startLoc}

	switch p.cur.Kind {
	case kwEVERY:
		p.advance() // consume EVERY

		// Collect interval: integer + unit keyword (e.g., "1 DAY")
		var parts []string
		if p.cur.Kind == tokInt {
			parts = append(parts, p.cur.Str)
			p.advance()
		}
		// Unit keyword (DAY, HOUR, MINUTE, SECOND, WEEK, MONTH, YEAR, etc.)
		if p.cur.Kind >= 700 || p.cur.Kind == tokIdent {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		}
		sched.Every = strings.Join(parts, " ")

		// Optional STARTS timestamp
		if p.cur.Kind == kwSTARTS {
			p.advance() // consume STARTS
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			sched.Starts = p.cur.Str
			sched.Loc.End = p.cur.Loc.End
			p.advance()
		}

		// Optional ENDS timestamp
		if p.cur.Kind == kwENDS {
			p.advance() // consume ENDS
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			sched.Ends = p.cur.Str
			sched.Loc.End = p.cur.Loc.End
			p.advance()
		}

	case kwAT:
		p.advance() // consume AT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		sched.At = p.cur.Str
		sched.Loc.End = p.cur.Loc.End
		p.advance()

	default:
		return nil, p.syntaxErrorAtCur()
	}

	return sched, nil
}

// parseAlterJob parses an ALTER JOB statement.
// On entry, ALTER has been consumed and cur is JOB.
//
// Syntax:
//
//	ALTER JOB job_name { PROPERTIES(...) | DO dml_statement }
func (p *Parser) parseAlterJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	stmt := &ast.AlterJobStmt{}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Kind {
	case kwPROPERTIES:
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		stmt.Loc = startLoc.Merge(p.prev.Loc)

	case kwDO:
		p.advance() // consume DO
		doStmt, err := p.parseJobDoBody()
		if err != nil {
			return nil, err
		}
		stmt.NewStmt = doStmt
		stmt.Loc = startLoc.Merge(ast.NodeLoc(doStmt))

	default:
		stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	}

	return stmt, nil
}

// parseDropJob parses a DROP JOB statement.
// On entry, DROP has been consumed and cur is JOB.
//
// Syntax:
//
//	DROP JOB [IF EXISTS] job_name
//	DROP JOB WHERE expr
func (p *Parser) parseDropJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	stmt := &ast.DropJobStmt{}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc = startLoc.Merge(ast.NodeLoc(where))
		return stmt, nil
	}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parsePauseJob parses a PAUSE JOB statement.
// On entry, PAUSE has been consumed and cur is JOB.
//
// Syntax:
//
//	PAUSE JOB job_name
//	PAUSE JOB WHERE expr
func (p *Parser) parsePauseJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	stmt := &ast.PauseJobStmt{}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc = startLoc.Merge(ast.NodeLoc(where))
		return stmt, nil
	}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parseResumeJob parses a RESUME JOB statement.
// On entry, RESUME has been consumed and cur is JOB.
//
// Syntax:
//
//	RESUME JOB job_name
//	RESUME JOB WHERE expr
func (p *Parser) parseResumeJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	stmt := &ast.ResumeJobStmt{}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc = startLoc.Merge(ast.NodeLoc(where))
		return stmt, nil
	}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parseCancelTask parses a CANCEL TASK FOR job_name task_id statement.
// On entry, CANCEL has been consumed and cur is TASK.
//
// Syntax:
//
//	CANCEL TASK FOR job_name task_id
func (p *Parser) parseCancelTask(startLoc ast.Loc) (ast.Node, error) {
	// Consume TASK
	if _, err := p.expect(kwTASK); err != nil {
		return nil, err
	}
	// Consume FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	stmt := &ast.CancelTaskStmt{}

	// job_name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.For = name

	// task_id — integer
	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.TaskID = p.cur.Ival
	stmt.Loc = startLoc.Merge(p.cur.Loc)
	p.advance()

	return stmt, nil
}

// parseShowJob parses a SHOW JOB or SHOW JOB TASK statement.
// On entry, SHOW has been consumed and cur is JOB.
//
// Syntax:
//
//	SHOW JOB [LIKE 'pat' | WHERE expr]
//	SHOW JOB TASK FOR job_name
func (p *Parser) parseShowJob(startLoc ast.Loc) (ast.Node, error) {
	// Consume JOB
	if _, err := p.expect(kwJOB); err != nil {
		return nil, err
	}

	// Check for SHOW JOB TASK FOR job_name
	if p.cur.Kind == kwTASK {
		p.advance() // consume TASK
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt := &ast.ShowJobTaskStmt{
			For: name,
			Loc: startLoc.Merge(ast.NodeLoc(name)),
		}
		return stmt, nil
	}

	stmt := &ast.ShowJobStmt{}

	switch p.cur.Kind {
	case kwLIKE:
		p.advance() // consume LIKE
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Like = p.cur.Str
		stmt.Loc = startLoc.Merge(p.cur.Loc)
		p.advance()

	case kwWHERE:
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
		stmt.Loc = startLoc.Merge(ast.NodeLoc(where))

	default:
		stmt.Loc = startLoc.Merge(p.prev.Loc)
	}

	return stmt, nil
}
