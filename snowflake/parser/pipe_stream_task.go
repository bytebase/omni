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
func (p *Parser) parseCreatePipeStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume PIPE

	stmt := &ast.CreatePipeStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
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
	for !p.curIs(kwAS) && p.startsCopyOption() {
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
	if p.curIs('(') {
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

	switch {
	case p.curIs(kwSET):
		p.advance() // consume SET
		if p.curIs(kwTAG) {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPipeSetTag
			stmt.Tags = tags
		} else {
			// The properties may be separated by commas (SET,
			// https://docs.snowflake.com/en/sql-reference/sql/alter-pipe).
			opts, err := p.parseCopyOptionsSep(true)
			if err != nil {
				return nil, err
			}
			if len(opts) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterPipeSet
			stmt.Options = opts
		}

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
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

	case p.curIs(kwREFRESH):
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
//	  ON { TABLE | VIEW | STAGE | EXTERNAL TABLE | EVENT TABLE | DYNAMIC TABLE } <object_name>
//	  [ { AT | BEFORE } ( <key> => <value> ) ]   -- not ON EVENT or DYNAMIC TABLE
//	  [ <config options> ]
//
//	CREATE [ OR REPLACE ] STREAM <name> CLONE <source> [ COPY GRANTS ] [ COPY = ( GRANTS ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-stream
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the STREAM
// keyword. Per the docs the [WITH] TAG clause precedes COPY GRANTS, which
// precedes ON. The trailing config options (APPEND_ONLY / INSERT_ONLY /
// SHOW_INITIAL_ROWS / COMMENT / …) are captured open-ended.
func (p *Parser) parseCreateStreamStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume STREAM

	stmt := &ast.CreateStreamStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
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

	// CLONE form: CREATE STREAM <name> CLONE <source> [ COPY GRANTS ]
	// [ COPY = ( GRANTS ) ].
	if p.curIs(kwCLONE) {
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
		if p.curIs(kwCOPY) && p.nextIs('=') {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			if len(opt.Group) != 1 || !opt.Group[0].Bare || opt.Group[0].Name != "GRANTS" {
				return nil, &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error: COPY takes ( GRANTS )"}
			}
			stmt.Options = []*ast.CopyOption{opt}
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
	switch {
	case p.curIs(kwTABLE):
		p.advance() // consume TABLE
		stmt.SourceKind = ast.StreamOnTable
	case p.curIs(kwVIEW):
		p.advance() // consume VIEW
		stmt.SourceKind = ast.StreamOnView
	case p.curIs(kwSTAGE):
		p.advance() // consume STAGE
		stmt.SourceKind = ast.StreamOnStage
	case p.curIs(kwEXTERNAL), p.curIs(kwEVENT), p.curIs(kwDYNAMIC):
		stmt.SourceKind = ast.StreamOnExternalTable
		if p.curIs(kwEVENT) {
			stmt.SourceKind = ast.StreamOnEventTable
		} else if p.curIs(kwDYNAMIC) {
			stmt.SourceKind = ast.StreamOnDynamicTable
		}
		p.advance() // consume EXTERNAL / EVENT / DYNAMIC
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
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
	// A stream on an event or dynamic table takes none.
	if (p.curIs(kwAT) || p.curIs(kwBEFORE)) &&
		stmt.SourceKind != ast.StreamOnEventTable && stmt.SourceKind != ast.StreamOnDynamicTable {
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
	if p.curIs(kwTAG) {
		return true
	}
	return p.curIs(kwWITH) && p.nextIs(kwTAG)
}

// parseStreamWithTags parses a [WITH] TAG (...) clause, accepting both the WITH
// TAG and bare TAG spellings. The caller must have confirmed startsStreamWithTags.
func (p *Parser) parseStreamWithTags() ([]*ast.TagAssignment, error) {
	if p.curIs(kwWITH) {
		p.advance() // consume WITH (peekNext == TAG guaranteed by caller)
	}
	return p.parseTagAssignments()
}

// startsCopyGrants reports whether cur begins a COPY GRANTS clause (COPY
// immediately followed by GRANTS).
func (p *Parser) startsCopyGrants() bool {
	return p.curIs(kwCOPY) && p.nextIs(kwGRANTS)
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

	switch {
	case p.curIs(kwSET):
		p.advance() // consume SET
		if p.curIs(kwTAG) {
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

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
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
//	  [ WITH CONTACT ( <purpose> = <contact> [ , ... ] ) ]
//	  [ <config options> ]                     (KEY = value; session parameters may be comma-separated)
//	  [ EXECUTE AS USER <user> ]
//	  [ AFTER <predecessor> [ , ... ] ]
//	  [ WHEN <boolean_expr> ]
//	  AS <sql>
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the TASK keyword.
// The config options (WAREHOUSE / SCHEDULE / CONFIG / OVERLAP_POLICY /
// session params / …) are captured open-ended, with AFTER / WHEN / AS as the
// structural anchors. AS's body is captured by parseEmbeddedBody (reuses
// parseStmt, with verbatim fallback for scripting / unsupported bodies).
func (p *Parser) parseCreateTaskStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume TASK

	stmt := &ast.CreateTaskStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
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
	if p.curIs(kwCLONE) {
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

	// WITH CONTACT ( <purpose> = <contact> [ , ... ] ) follows the tags.
	if p.curIs(kwWITH) && p.nextWord("CONTACT") {
		p.advance() // consume WITH
		p.advance() // consume CONTACT
		group, err := p.parseTaskGraphGroup(false, false)
		if err != nil {
			return nil, err
		}
		if err := checkTaskGraphContacts(group); err != nil {
			return nil, err
		}
		stmt.Contacts = group.Group
	}

	// Config options, AFTER, and WHEN may appear before AS. WHEN and AFTER are
	// the structural anchors; AS terminates the option run. AFTER and WHEN are
	// keywords startsCopyOption would treat as option names, so they are matched
	// before the option branch.
	for {
		switch {
		case p.curIs(kwAFTER):
			after, err := p.parseTaskAfterList()
			if err != nil {
				return nil, err
			}
			stmt.After = after
			continue
		case p.curIs(kwWHEN):
			p.advance() // consume WHEN
			when, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.When = when
			continue
		case p.curIs(kwAS):
			// fallthrough out of the loop below
		default:
			if start := p.cur.Loc.Start; p.acceptWords("EXECUTE", "AS", "USER") {
				user, err := p.parseObjectSimpleName()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, &ast.CopyOption{Name: "EXECUTE AS USER", Words: user.String(), Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
				continue
			}
			if p.startsTaskOption() {
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, opt)
				if err := p.skipTaskOptionComma(opt); err != nil {
					return nil, err
				}
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

// startsTaskOption reports whether cur begins a task's `KEY = value` option.
// A task has no value-less options.
func (p *Parser) startsTaskOption() bool {
	return p.startsCopyOption() && p.nextIs('=')
}

// skipTaskOptionComma consumes the ',' that separates session parameters
// (<session_parameter> = <value> [ , ... ]) after option opt, which must then
// be followed by another `KEY = value`.
// https://docs.snowflake.com/en/sql-reference/sql/create-task
func (p *Parser) skipTaskOptionComma(opt *ast.CopyOption) error {
	if !p.curIs(',') {
		return nil
	}
	p.advance() // consume ','
	if !p.startsTaskOption() {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseTaskSetOptions reads ALTER TASK ... SET's options: open-ended
// `KEY = value` pairs, session parameters separated by commas,
// CONTACT <purpose> = <contact> [ , ... ] (the option CONTACT, its Group the
// entries) and EXECUTE AS USER <user> (the option "EXECUTE AS USER").
// https://docs.snowflake.com/en/sql-reference/sql/alter-task
func (p *Parser) parseTaskSetOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for {
		start := p.cur.Loc.Start
		switch {
		case p.curWord("CONTACT") && isNameWordTok(p.peekNext()):
			p.advance() // consume CONTACT
			group := &ast.CopyOption{}
			for {
				e, err := p.parseTaskGraphEntry(false, false)
				if err != nil {
					return nil, err
				}
				group.Group = append(group.Group, e)
				if !p.curIs(',') || !p.commaStartsContactEntry() {
					break
				}
				p.advance() // consume ','
			}
			if err := checkTaskGraphContacts(group); err != nil {
				return nil, err
			}
			opts = append(opts, &ast.CopyOption{Name: "CONTACT", Group: group.Group, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		case p.acceptWords("EXECUTE", "AS", "USER"):
			user, err := p.parseObjectSimpleName()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &ast.CopyOption{Name: "EXECUTE AS USER", Words: user.String(), Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		case p.startsTaskOption():
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
			if err := p.skipTaskOptionComma(opt); err != nil {
				return nil, err
			}
		default:
			if len(opts) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			return opts, nil
		}
	}
}

// parseTaskUnsetList reads ALTER TASK ... UNSET's <property> [ , ... ], where
// a property is a name word, CONTACT <purpose> or EXECUTE AS USER.
// https://docs.snowflake.com/en/sql-reference/sql/alter-task
func (p *Parser) parseTaskUnsetList() ([]string, error) {
	var props []string
	for {
		switch {
		case p.curWord("CONTACT") && isNameWordTok(p.peekNext()):
			p.advance() // consume CONTACT
			props = append(props, "CONTACT "+strings.ToUpper(p.advance().Str))
		case p.acceptWords("EXECUTE", "AS", "USER"):
			props = append(props, "EXECUTE AS USER")
		case p.isOptionWord(p.cur.Type):
			props = append(props, strings.ToUpper(p.advance().Str))
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if !p.curIs(',') {
			return props, nil
		}
		p.advance() // consume ','
	}
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
		if p.curIs(',') {
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
//	REMOVE WHEN
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-task
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

	switch {
	case p.curIs(kwRESUME):
		p.advance()
		stmt.Action = ast.AlterTaskResume

	case p.curIs(kwSUSPEND):
		p.advance()
		stmt.Action = ast.AlterTaskSuspend

	case p.curIs(kwADD):
		p.advance() // consume ADD
		after, err := p.parseTaskAfterList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterTaskAddAfter
		stmt.After = after

	case p.curIs(kwREMOVE):
		p.advance() // consume REMOVE
		if p.curIs(kwWHEN) {
			// REMOVE WHEN
			p.advance() // consume WHEN
			stmt.Action = ast.AlterTaskRemoveWhen
			break
		}
		after, err := p.parseTaskAfterList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterTaskRemoveAfter
		stmt.After = after

	case p.curIs(kwSET):
		p.advance() // consume SET
		if p.curIs(kwTAG) {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskSetTag
			stmt.Tags = tags
		} else {
			opts, err := p.parseTaskSetOptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskSet
			stmt.Options = opts
		}

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseTaskUnsetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskUnset
			stmt.UnsetProps = props
		}

	case p.curIs(kwMODIFY):
		p.advance() // consume MODIFY
		switch {
		case p.curIs(kwAS):
			p.advance() // consume AS
			body, raw, err := p.parseEmbeddedBody()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTaskModifyAs
			stmt.Body = body
			stmt.BodyRaw = raw
		case p.curIs(kwWHEN):
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
//
// The page also gives CREATE [ OR REPLACE ] ALERT <name> CLONE <source> and
// CREATE ... ALERT <name> FROM TEMPLATE <template> [ [ WITH ] TAG ( ... ) ]
// [ <option> ... ], which take no condition or action.
// https://docs.snowflake.com/en/sql-reference/sql/create-alert
func (p *Parser) parseCreateAlertStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume ALERT

	stmt := &ast.CreateAlertStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
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

	// CREATE [ OR REPLACE ] ALERT <name> CLONE <source_alert>
	if p.curIs(kwCLONE) && !stmt.IfNotExists && !orAlter {
		p.advance() // consume CLONE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Clone = src
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// CREATE ... ALERT <name> FROM TEMPLATE <template> [ [ WITH ] TAG ( ... ) ]
	//   [ <option> ... ]
	// CREATE OR ALTER ALERT <name> FROM TEMPLATE ... creates or re-renders a
	// templated alert; the page gives it no WITH TAG and no IF NOT EXISTS.
	if p.curIs(kwFROM) && !(orAlter && stmt.IfNotExists) {
		p.advance() // consume FROM
		if err := p.expectWords("TEMPLATE"); err != nil {
			return nil, err
		}
		tmpl, err := p.parseObjectSimpleName()
		if err != nil {
			return nil, err
		}
		stmt.Template = tmpl
		if !orAlter && p.startsStreamWithTags() {
			tags, err := p.parseStreamWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = tags
		}
		opts, err := p.parseAlertTemplateOptions(false)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
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

	// Config options precede IF. IF is a keyword startsCopyOption would treat as
	// an option name, so the loop stops explicitly at IF.
	for !p.curIs(kwIF) && p.startsCopyOption() {
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

// alertTemplateProps are the options an alert made from a template takes, as
// CREATE ALERT ... FROM TEMPLATE and ALTER ALERT ... FROM TEMPLATE ... SET
// document them.
// https://docs.snowflake.com/en/sql-reference/sql/create-alert
// https://docs.snowflake.com/en/sql-reference/sql/alter-alert
var alertTemplateProps = objectProps{
	"WAREHOUSE":                        propName,
	"SCHEDULE":                         propString,
	"COMMENT":                          propString,
	"RUNBOOK":                          propString,
	"SUSPEND_ALERT_AFTER_NUM_FAILURES": propNumber,
	"TEMPLATE_PARAMS":                  propString,
}

// parseAlertTemplateOptions reads the `KEY = value` options of an alert made
// from a template, each one of alertTemplateProps, once. need says whether
// the run may be empty.
func (p *Parser) parseAlertTemplateOptions(need bool) ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	seen := map[string]bool{}
	for p.startsObjectProperty() {
		o, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if v, ok := alertTemplateProps[o.Name]; !ok || !v.matches(o) || seen[o.Name] {
			return nil, &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error at or near " + o.Name}
		}
		seen[o.Name] = true
		opts = append(opts, o)
	}
	if need && len(opts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return opts, nil
}

// parseAlertCondition parses the query inside IF(EXISTS(<condition>)). Per the
// docs the condition is a SELECT, SHOW, or CALL. SELECT/WITH reuse the shared
// query-expression parser; SHOW reuses parseShowStmt; any other condition
// (e.g. CALL, unsupported by parseStmt) is captured verbatim by consuming
// balanced tokens up to the closing ')' so the alert still parses. Returns the
// parsed node (nil for the verbatim case) and the verbatim condition text.
func (p *Parser) parseAlertCondition() (ast.Node, string, error) {
	condStartAbs := p.cur.Loc.Start

	switch {
	case p.curIs(kwSELECT), p.curIs('('):
		node, err := p.parseQueryExpr()
		if err != nil {
			return nil, "", err
		}
		return node, p.srcSlice(condStartAbs, p.prev.Loc.End), nil
	case p.curIs(kwWITH):
		node, err := p.parseWithQueryExpr()
		if err != nil {
			return nil, "", err
		}
		return node, p.srcSlice(condStartAbs, p.prev.Loc.End), nil
	case p.curIs(kwSHOW):
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
	for !p.curIs(tokEOF) {
		if p.curIs(')') && depth == 0 {
			break
		}
		switch {
		case p.curIs('('):
			depth++
		case p.curIs(')'):
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
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	MODIFY CONDITION EXISTS ( <query> )
//	MODIFY ACTION <action>
//	FROM TEMPLATE <template> SET <option> [ ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-alert
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

	switch {
	case p.curIs(kwRESUME):
		p.advance()
		stmt.Action = ast.AlterAlertResume

	case p.curIs(kwSUSPEND):
		p.advance()
		stmt.Action = ast.AlterAlertSuspend

	case p.curIs(kwFROM):
		// FROM TEMPLATE <template> SET <option> [ ... ]
		p.advance() // consume FROM
		if err := p.expectWords("TEMPLATE"); err != nil {
			return nil, err
		}
		tmpl, err := p.parseObjectSimpleName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		opts, err := p.parseAlertTemplateOptions(true)
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAlertSet
		stmt.Template = tmpl
		stmt.Options = opts

	case p.curIs(kwSET):
		p.advance() // consume SET
		if p.curIs(kwTAG) {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterAlertSetTag
			stmt.Tags = tags
			break
		}
		opts, err := p.parseRequiredOptions()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAlertSet
		stmt.Options = opts

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterAlertUnsetTag
			stmt.UnsetTags = names
			break
		}
		props, err := p.parseUnsetPropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAlertUnset
		stmt.UnsetProps = props

	case p.curIs(kwMODIFY):
		p.advance() // consume MODIFY
		switch {
		case p.curIs(kwCONDITION):
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
		case p.curIs(kwACTION):
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
	if p.curIs(tokEOF) {
		// AS / THEN with no body is a syntax error.
		return nil, "", p.syntaxErrorAtCur()
	}
	bodyStartAbs := p.cur.Loc.Start

	// A BEGIN / DECLARE opener is a scripting block parseStmt cannot parse; go
	// straight to verbatim capture without attempting (and erroring) parseStmt.
	if !p.curIs(kwBEGIN) && !p.curIs(kwDECLARE) {
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
	for !p.curIs(tokEOF) {
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
		if p.curIs(',') {
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
	if p.curIs(kwIF) && p.nextIs(kwNOT) {
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
	if p.curIs(kwIF) && p.nextIs(kwEXISTS) {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		*dst = true
	}
	return nil
}
