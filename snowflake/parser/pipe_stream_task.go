package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Data-pipeline DDL — CREATE / ALTER PIPE / STREAM / TASK / ALERT (T4.3)
// ---------------------------------------------------------------------------
//
// PIPE / STREAM / TASK / ALERT are Snowflake's data-pipeline objects. Several
// embed other statements:
//
//   - PIPE   AS <copy_into_table>          — reuses parseCopyStmt (T5.2).
//   - TASK   AS <sql>                       — reuses parseStmt (any statement);
//                                             falls back to verbatim capture for a
//                                             Snowflake Scripting BEGIN…END /
//                                             DECLARE…BEGIN…END body or any
//                                             statement parseStmt cannot yet parse.
//   - ALERT  IF(EXISTS(<query>)) THEN <act> — condition reuses parseQueryExpr; the
//                                             THEN action reuses parseStmt with the
//                                             same verbatim fallback.
//
// Each object's growable config vocabulary (PIPE's AUTO_INGEST / INTEGRATION /
// AWS_SNS_TOPIC / ERROR_INTEGRATION; TASK's WAREHOUSE / SCHEDULE / CONFIG /
// OVERLAP_POLICY / SUCCESS_INTEGRATION / LOG_LEVEL / session params / …; ALERT's
// WAREHOUSE / SCHEDULE / CONFIG / RUNBOOK / …) is parsed as open-ended
// `KEY = <value>` pairs (ast.CopyOption), reusing the merged STAGE (T4.1) /
// FILE FORMAT (T4.2) / COPY (T5.2) machinery rather than mirroring the legacy
// ANTLR grammar's stale enumerations. Only the structural anchors are modeled
// explicitly: the embedded bodies, STREAM's ON <source> + AT/BEFORE clause,
// TASK's AFTER list and WHEN predicate, and the [WITH] TAG / COPY GRANTS clauses.

// ---------------------------------------------------------------------------
// CREATE PIPE
// ---------------------------------------------------------------------------

// parseCreatePipeStmt parses
//
//	CREATE [ OR REPLACE ] PIPE [ IF NOT EXISTS ] <name>
//	  [ <config options> ] AS <copy_into_table>
//
// The CREATE keyword and the optional OR REPLACE modifier have already been
// consumed by parseCreateStmt; start is the Loc of the CREATE token and cur is
// the PIPE keyword. Every option before AS (AUTO_INGEST / ERROR_INTEGRATION /
// AWS_SNS_TOPIC / INTEGRATION / COMMENT / …) is captured open-ended. AS is the
// structural anchor; its body is a COPY INTO, optionally wrapped in parens.
func (p *Parser) parseCreatePipeStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume PIPE

	stmt := &ast.CreatePipeStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Config options precede AS. AS is a keyword that startsCopyOption would
	// otherwise treat as an option name, so the loop stops explicitly at AS.
	for p.cur.Type != kwAS && p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// AS <copy_into_table>, optionally parenthesized: AS (COPY INTO …). The
	// parens are consumed and not retained. parseCopyStmt expects cur at COPY.
	copyNode, err := p.parseParenthesizedCopy()
	if err != nil {
		return nil, err
	}
	stmt.Copy = copyNode

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseParenthesizedCopy parses a COPY statement that may be wrapped in a single
// layer of parentheses (the pipe `AS (COPY INTO …)` form seen in the docs
// corpus). A leading '(' is consumed and its matching ')' is required after the
// COPY body; otherwise the COPY is parsed bare. cur must be '(' or COPY.
func (p *Parser) parseParenthesizedCopy() (ast.Node, error) {
	if p.cur.Type == '(' {
		p.advance() // consume '('
		copyNode, err := p.parseCopyStmt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return copyNode, nil
	}
	return p.parseCopyStmt()
}

// ---------------------------------------------------------------------------
// ALTER PIPE
// ---------------------------------------------------------------------------

// parseAlterPipeStmt parses ALTER PIPE [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the PIPE keyword.
//
//	SET <options>
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]                 (e.g. PIPE_EXECUTION_PAUSED, COMMENT)
//	REFRESH [ PREFIX = '<path>' ] [ MODIFIED_AFTER = '<timestamp>' ]
func (p *Parser) parseAlterPipeStmt() (ast.Node, error) {
	altTok := p.advance() // consume PIPE
	stmt := &ast.AlterPipeStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPipeSetTag
			stmt.Tags = tags
		} else {
			opts, err := p.parseRequiredOptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPipeSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPipeUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseUnsetPropertyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPipeUnset
			stmt.UnsetProps = props
		}

	case kwREFRESH:
		// REFRESH [ PREFIX = '<path>' ] [ MODIFIED_AFTER = '<timestamp>' ] — both
		// optional, captured open-ended.
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterPipeRefresh
		for p.startsCopyOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE STREAM
// ---------------------------------------------------------------------------

// parseCreateStreamStmt parses
//
//	CREATE [ OR REPLACE ] STREAM [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( … ) ] [ COPY GRANTS ]
//	  ON { TABLE | VIEW | STAGE | EXTERNAL TABLE } <object_name>
//	  [ { AT | BEFORE } ( <key> => <value> ) ]
//	  [ <config options> ]
//
//	CREATE [ OR REPLACE ] STREAM <name> CLONE <source> [ COPY GRANTS ]
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the STREAM
// keyword. Per the docs the [WITH] TAG clause precedes COPY GRANTS, which
// precedes ON. The trailing config options (APPEND_ONLY / INSERT_ONLY /
// SHOW_INITIAL_ROWS / COMMENT / …) are captured open-ended.
func (p *Parser) parseCreateStreamStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume STREAM

	stmt := &ast.CreateStreamStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CLONE form: CREATE STREAM <name> CLONE <source> [COPY GRANTS].
	if p.cur.Type == kwCLONE {
		p.advance() // consume CLONE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Clone = src
		if p.startsCopyGrants() {
			p.advance() // consume COPY
			p.advance() // consume GRANTS
			stmt.CopyGrants = true
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// [WITH] TAG (...) and COPY GRANTS both precede ON, in either order
	// defensively (docs list TAG then COPY GRANTS; the legacy grammar agrees).
	for {
		if p.startsStreamWithTags() {
			tags, err := p.parseStreamWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		if p.startsCopyGrants() {
			p.advance() // consume COPY
			p.advance() // consume GRANTS
			stmt.CopyGrants = true
			continue
		}
		break
	}

	// ON { TABLE | VIEW | STAGE | EXTERNAL TABLE } <object_name> — mandatory.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt.SourceKind = ast.StreamOnTable
	case kwVIEW:
		p.advance() // consume VIEW
		stmt.SourceKind = ast.StreamOnView
	case kwSTAGE:
		p.advance() // consume STAGE
		stmt.SourceKind = ast.StreamOnStage
	case kwEXTERNAL:
		p.advance() // consume EXTERNAL
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		stmt.SourceKind = ast.StreamOnExternalTable
	default:
		return nil, p.syntaxErrorAtCur()
	}
	src, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Source = src

	// Optional { AT | BEFORE } ( <key> => <value> ) time-travel clause. Parsed
	// inline/open-ended (not via the T5.3 query-clauses node).
	if p.cur.Type == kwAT || p.cur.Type == kwBEFORE {
		tt, err := p.parseStreamTimeTravel()
		if err != nil {
			return nil, err
		}
		stmt.TimeTravel = tt
	}

	// Trailing config options: APPEND_ONLY / INSERT_ONLY / SHOW_INITIAL_ROWS /
	// COMMENT / … — open-ended KEY = value pairs.
	for p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseStreamTimeTravel parses a stream's { AT | BEFORE } ( <key> => <value> )
// clause. cur is AT or BEFORE. key ∈ { TIMESTAMP | OFFSET | STATEMENT | STREAM }
// (captured verbatim/uppercased, not enumerated, to tolerate doc growth); the
// value after => is a general expression (the corpus shows TO_TIMESTAMP(...),
// arithmetic like -60*5, and string literals), parsed by the shared expression
// parser.
func (p *Parser) parseStreamTimeTravel() (*ast.StreamTimeTravel, error) {
	atTok := p.advance() // consume AT or BEFORE
	tt := &ast.StreamTimeTravel{
		AtBefore: strings.ToUpper(atTok.Str),
		Loc:      ast.Loc{Start: atTok.Loc.Start},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if !p.isOptionWord(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	tt.Key = strings.ToUpper(p.advance().Str)

	if _, err := p.expect(tokAssoc); err != nil {
		return nil, err
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	tt.Value = val

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	tt.Loc.End = p.prev.Loc.End
	return tt, nil
}

// startsStreamWithTags reports whether cur begins a [WITH] TAG clause: the TAG
// keyword directly, or WITH immediately followed by TAG.
func (p *Parser) startsStreamWithTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
}

// parseStreamWithTags parses a [WITH] TAG (...) clause, accepting both the WITH
// TAG and bare TAG spellings. The caller must have confirmed startsStreamWithTags.
func (p *Parser) parseStreamWithTags() ([]*ast.TagAssignment, error) {
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH (peekNext == TAG guaranteed by caller)
	}
	return p.parseTagAssignments()
}

// startsCopyGrants reports whether cur begins a COPY GRANTS clause (COPY
// immediately followed by GRANTS).
func (p *Parser) startsCopyGrants() bool {
	return p.cur.Type == kwCOPY && p.peekNext().Type == kwGRANTS
}

// ---------------------------------------------------------------------------
// ALTER STREAM
// ---------------------------------------------------------------------------

// parseAlterStreamStmt parses ALTER STREAM [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the STREAM keyword.
//
//	SET <options>                              (COMMENT = '…', …)
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]                  (e.g. COMMENT)
func (p *Parser) parseAlterStreamStmt() (ast.Node, error) {
	altTok := p.advance() // consume STREAM
	stmt := &ast.AlterStreamStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStreamSetTag
			stmt.Tags = tags
		} else {
			opts, err := p.parseRequiredOptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStreamSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStreamUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseUnsetPropertyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStreamUnset
			stmt.UnsetProps = props
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE TASK
// ---------------------------------------------------------------------------

// parseCreateTaskStmt parses
//
//	CREATE [ OR REPLACE ] TASK [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( … ) ]
//	  [ <config options> ]
//	  [ AFTER <predecessor> [ , ... ] ]
//	  [ WHEN <boolean_expr> ]
//	  AS <sql>
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the TASK keyword.
// The config options (WAREHOUSE / SCHEDULE / CONFIG / OVERLAP_POLICY /
// session params / …) are captured open-ended, with AFTER / WHEN / AS as the
// structural anchors. AS's body is captured by parseEmbeddedBody (reuses
// parseStmt, with verbatim fallback for scripting / unsupported bodies).
func (p *Parser) parseCreateTaskStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume TASK

	stmt := &ast.CreateTaskStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CLONE form: CREATE [OR REPLACE] TASK <name> CLONE <source_task>.
	if p.cur.Type == kwCLONE {
		p.advance() // consume CLONE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Clone = src
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// [WITH] TAG (...) precedes the config options.
	for p.startsStreamWithTags() {
		tags, err := p.parseStreamWithTags()
		if err != nil {
			return nil, err
		}
		stmt.Tags = append(stmt.Tags, tags...)
	}

	// Config options, AFTER, and WHEN may appear before AS. WHEN and AFTER are
	// the structural anchors; AS terminates the option run. AFTER and WHEN are
	// keywords startsCopyOption would treat as option names, so they are matched
	// before the option branch.
	for {
		switch p.cur.Type {
		case kwAFTER:
			after, err := p.parseTaskAfterList()
			if err != nil {
				return nil, err
			}
			stmt.After = after
			continue
		case kwWHEN:
			p.advance() // consume WHEN
			when, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.When = when
			continue
		case kwAS:
			// fallthrough out of the loop below
		default:
			if p.startsCopyOption() {
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, opt)
				continue
			}
		}
		break
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	body, raw, err := p.parseEmbeddedBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body
	stmt.BodyRaw = raw

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTaskAfterList parses AFTER <predecessor> [ , <predecessor> , ... ]. The
// docs spell predecessors as strings while the corpus uses bare identifiers;
// parseObjectName covers both (a bare/quoted name). The AFTER keyword is
// consumed here.
func (p *Parser) parseTaskAfterList() ([]*ast.ObjectName, error) {
	if _, err := p.expect(kwAFTER); err != nil {
		return nil, err
	}
	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return names, nil
}

// ---------------------------------------------------------------------------
// ALTER TASK
// ---------------------------------------------------------------------------

// parseAlterTaskStmt parses ALTER TASK [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the TASK keyword.
//
//	RESUME | SUSPEND
//	{ ADD | REMOVE } AFTER <task> [ , ... ]
//	SET <options>
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]
//	MODIFY AS <sql>
//	MODIFY WHEN <boolean_expr>
func (p *Parser) parseAlterTaskStmt() (ast.Node, error) {
	altTok := p.advance() // consume TASK
	stmt := &ast.AlterTaskStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwRESUME:
		p.advance()
		stmt.Action = ast.AlterTaskResume

	case kwSUSPEND:
		p.advance()
		stmt.Action = ast.AlterTaskSuspend

	case kwADD:
		p.advance() // consume ADD
		after, err := p.parseTaskAfterList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterTaskAddAfter
		stmt.After = after

	case kwREMOVE:
		p.advance() // consume REMOVE
		after, err := p.parseTaskAfterList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterTaskRemoveAfter
		stmt.After = after

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskSetTag
			stmt.Tags = tags
		} else {
			opts, err := p.parseRequiredOptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseUnsetPropertyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskUnset
			stmt.UnsetProps = props
		}

	case kwMODIFY:
		p.advance() // consume MODIFY
		switch p.cur.Type {
		case kwAS:
			p.advance() // consume AS
			body, raw, err := p.parseEmbeddedBody()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskModifyAs
			stmt.Body = body
			stmt.BodyRaw = raw
		case kwWHEN:
			p.advance() // consume WHEN
			when, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskModifyWhen
			stmt.When = when
		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE ALERT
// ---------------------------------------------------------------------------

// parseCreateAlertStmt parses
//
//	CREATE [ OR REPLACE ] ALERT [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( … ) ]
//	  [ <config options> ]                     (WAREHOUSE / SCHEDULE / COMMENT / …)
//	  IF ( EXISTS ( <condition> ) )
//	  THEN <action>
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the ALERT
// keyword. WAREHOUSE and SCHEDULE are NOT lifted out — WAREHOUSE is optional for
// serverless alerts and SCHEDULE's value grows — so all config is open-ended.
// IF / THEN are the structural anchors. The condition reuses parseQueryExpr; the
// action reuses parseEmbeddedBody (parseStmt + verbatim fallback).
func (p *Parser) parseCreateAlertStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume ALERT

	stmt := &ast.CreateAlertStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// [WITH] TAG (...) precedes the config options.
	for p.startsStreamWithTags() {
		tags, err := p.parseStreamWithTags()
		if err != nil {
			return nil, err
		}
		stmt.Tags = append(stmt.Tags, tags...)
	}

	// Config options precede IF. IF is a keyword startsCopyOption would treat as
	// an option name, so the loop stops explicitly at IF.
	for p.cur.Type != kwIF && p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
	}

	// IF ( EXISTS ( <condition> ) ) THEN <action>.
	if _, err := p.expect(kwIF); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwEXISTS); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cond, condRaw, err := p.parseAlertCondition()
	if err != nil {
		return nil, err
	}
	stmt.Condition = cond
	stmt.ConditionRaw = condRaw
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}

	action, actionRaw, err := p.parseEmbeddedBody()
	if err != nil {
		return nil, err
	}
	stmt.Action = action
	stmt.ActionRaw = actionRaw

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlertCondition parses the query inside IF(EXISTS(<condition>)). Per the
// docs the condition is a SELECT, SHOW, or CALL. SELECT/WITH reuse the shared
// query-expression parser; SHOW reuses parseShowStmt; any other condition
// (e.g. CALL, unsupported by parseStmt) is captured verbatim by consuming
// balanced tokens up to the closing ')' so the alert still parses. Returns the
// parsed node (nil for the verbatim case) and the verbatim condition text.
func (p *Parser) parseAlertCondition() (ast.Node, string, error) {
	condStartAbs := p.cur.Loc.Start

	switch p.cur.Type {
	case kwSELECT, '(':
		node, err := p.parseQueryExpr()
		if err != nil {
			return nil, "", err
		}
		return node, p.srcSlice(condStartAbs, p.prev.Loc.End), nil
	case kwWITH:
		node, err := p.parseWithQueryExpr()
		if err != nil {
			return nil, "", err
		}
		return node, p.srcSlice(condStartAbs, p.prev.Loc.End), nil
	case kwSHOW:
		node, err := p.parseShowStmt()
		if err != nil {
			return nil, "", err
		}
		return node, p.srcSlice(condStartAbs, p.prev.Loc.End), nil
	default:
		// CALL or any other condition: capture verbatim up to the matching ')'
		// that closes EXISTS( … ), tracking nested parens.
		raw := p.captureBalancedRaw(condStartAbs)
		return nil, raw, nil
	}
}

// captureBalancedRaw consumes tokens from the current position until the parser
// reaches the ')' that closes the paren level the caller is already inside (i.e.
// until nesting depth would go negative), without consuming that ')'. It returns
// the verbatim source text from startAbs to the end of the last consumed token.
// Used to capture an embedded condition (e.g. a CALL) that no available
// sub-parser handles, so the enclosing statement still parses cleanly.
func (p *Parser) captureBalancedRaw(startAbs int) string {
	depth := 0
	endAbs := startAbs
	for p.cur.Type != tokEOF {
		if p.cur.Type == ')' && depth == 0 {
			break
		}
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			depth--
		}
		endAbs = p.cur.Loc.End
		p.advance()
	}
	return p.srcSlice(startAbs, endAbs)
}

// ---------------------------------------------------------------------------
// ALTER ALERT
// ---------------------------------------------------------------------------

// parseAlterAlertStmt parses ALTER ALERT [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the ALERT keyword.
//
//	RESUME | SUSPEND
//	SET <options>                              (WAREHOUSE = … | SCHEDULE = … | COMMENT = …)
//	UNSET <property> [ , ... ]                  (WAREHOUSE | SCHEDULE | COMMENT)
//	MODIFY CONDITION EXISTS ( <query> )
//	MODIFY ACTION <action>
func (p *Parser) parseAlterAlertStmt() (ast.Node, error) {
	altTok := p.advance() // consume ALERT
	stmt := &ast.AlterAlertStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwRESUME:
		p.advance()
		stmt.Action = ast.AlterAlertResume

	case kwSUSPEND:
		p.advance()
		stmt.Action = ast.AlterAlertSuspend

	case kwSET:
		p.advance() // consume SET
		opts, err := p.parseRequiredOptions()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAlertSet
		stmt.Options = opts

	case kwUNSET:
		p.advance() // consume UNSET
		props, err := p.parseUnsetPropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAlertUnset
		stmt.UnsetProps = props

	case kwMODIFY:
		p.advance() // consume MODIFY
		switch p.cur.Type {
		case kwCONDITION:
			// MODIFY CONDITION EXISTS ( <query> )
			p.advance() // consume CONDITION
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			cond, condRaw, err := p.parseAlertCondition()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterAlertModifyCondition
			stmt.Condition = cond
			stmt.ConditionRaw = condRaw
		case kwACTION:
			// MODIFY ACTION <action>
			p.advance() // consume ACTION
			action, raw, err := p.parseEmbeddedBody()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterAlertModifyAction
			stmt.ActionBody = action
			stmt.ActionRaw = raw
		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// parseEmbeddedBody parses an embedded statement body (TASK's AS <sql>, ALERT's
// THEN <action>), reusing the top-level parseStmt where it can and falling back
// to a verbatim capture otherwise. It returns the parsed node (nil on fallback)
// and the verbatim body source text (always non-empty for a well-formed body).
//
// A Snowflake Scripting block (BEGIN…END / DECLARE…BEGIN…END) and any statement
// parseStmt does not yet support (e.g. CALL) cannot be parsed structurally; for
// those the parser snapshot is restored and the remaining tokens of the segment
// are consumed verbatim, so the enclosing statement is never rejected on account
// of its body. (F3's Split keeps a `… AS BEGIN … END;` body as one segment, so
// the remaining tokens are exactly the body.)
func (p *Parser) parseEmbeddedBody() (ast.Node, string, error) {
	if p.cur.Type == tokEOF {
		// AS / THEN with no body is a syntax error.
		return nil, "", p.syntaxErrorAtCur()
	}
	bodyStartAbs := p.cur.Loc.Start

	// A BEGIN / DECLARE opener is a scripting block parseStmt cannot parse; go
	// straight to verbatim capture without attempting (and erroring) parseStmt.
	if p.cur.Type != kwBEGIN && p.cur.Type != kwDECLARE {
		snap := p.snapshot()
		node, err := p.parseStmt()
		if err == nil && node != nil {
			return node, p.srcSlice(bodyStartAbs, p.prev.Loc.End), nil
		}
		// Unsupported / unparsable body: restore and fall through to verbatim.
		p.restore(snap)
	}

	raw := p.captureRawToEOF(bodyStartAbs)
	return nil, raw, nil
}

// captureRawToEOF consumes every remaining token of the current segment and
// returns the verbatim source text from startAbs to the end of the last token.
func (p *Parser) captureRawToEOF(startAbs int) string {
	endAbs := startAbs
	for p.cur.Type != tokEOF {
		endAbs = p.cur.Loc.End
		p.advance()
	}
	return p.srcSlice(startAbs, endAbs)
}

// parserSnapshot captures enough Parser state to roll back a speculative parse
// (used by parseEmbeddedBody before attempting parseStmt on a body that may turn
// out to be unsupported). The Lexer is a value-copyable struct (string + int
// cursors + an error slice), so a shallow copy plus the buffered-token and
// error-length fields fully restores the position.
type parserSnapshot struct {
	lexer     Lexer
	cur       Token
	prev      Token
	nextBuf   Token
	hasNext   bool
	numErrors int
}

// snapshot records the current Parser state for a later restore.
func (p *Parser) snapshot() parserSnapshot {
	return parserSnapshot{
		lexer:     *p.lexer,
		cur:       p.cur,
		prev:      p.prev,
		nextBuf:   p.nextBuf,
		hasNext:   p.hasNext,
		numErrors: len(p.errors),
	}
}

// restore rolls the Parser back to a previously captured snapshot, discarding
// any tokens consumed and any errors appended since.
func (p *Parser) restore(s parserSnapshot) {
	*p.lexer = s.lexer
	p.cur = s.cur
	p.prev = s.prev
	p.nextBuf = s.nextBuf
	p.hasNext = s.hasNext
	p.errors = p.errors[:s.numErrors]
}

// parseRequiredOptions parses a run of one or more open-ended KEY = value
// options (reusing parseCopyOption). An empty run — a SET with nothing settable
// — is a syntax error. Shared by the ALTER … SET branches.
func (p *Parser) parseRequiredOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	if len(opts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return opts, nil
}

// parseUnsetPropertyList parses an UNSET <property> [ , ... ] list, where each
// property is a single name word (identifier or keyword), uppercased. Property
// names are not enumerated. Consumes at least one property.
func (p *Parser) parseUnsetPropertyList() ([]string, error) {
	var props []string
	for {
		if !p.isOptionWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		props = append(props, strings.ToUpper(p.advance().Str))
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return props, nil
}

// parseIfNotExistsInto consumes an optional IF NOT EXISTS prefix, setting *dst
// when present. cur should be just past the object-type keyword.
func (p *Parser) parseIfNotExistsInto(dst *bool) error {
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if _, err := p.expect(kwEXISTS); err != nil {
			return err
		}
		*dst = true
	}
	return nil
}

// parseIfExistsInto consumes an optional IF EXISTS prefix, setting *dst when
// present. cur should be just past the object-type keyword.
func (p *Parser) parseIfExistsInto(dst *bool) error {
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		*dst = true
	}
	return nil
}
