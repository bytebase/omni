package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Warehouse DDL — CREATE / ALTER WAREHOUSE (gap-warehouse)
// ---------------------------------------------------------------------------
//
// A warehouse carries a large, version-growing vocabulary of object properties
// (WAREHOUSE_SIZE, WAREHOUSE_TYPE, RESOURCE_CONSTRAINT, AUTO_RESUME,
// AUTO_SUSPEND, INITIALLY_SUSPENDED, GENERATION, MAX_CLUSTER_COUNT,
// MIN_CLUSTER_COUNT, SCALING_POLICY, ENABLE_QUERY_ACCELERATION, COMMENT, ...)
// alongside object-parameters and session-parameters that share the same
// `KEY = value` shape. Rather than mirror the legacy ANTLR grammar's finite,
// already-stale wh_properties / wh_common_size enumeration (it predates
// WAREHOUSE_TYPE, RESOURCE_CONSTRAINT, GENERATION, ENABLE_QUERY_ACCELERATION,
// QUERY_ACCELERATION_MAX_SCALE_FACTOR, ...), every property is parsed as an
// open-ended `KEY = <value>` pair (ast.CopyOption), reusing the merged COPY
// (T5.2) machinery (parseCopyOption / startsCopyOption) exactly as STAGE (T4.1)
// and FILE FORMAT (T4.2) do. The trailing [WITH] TAG (...) clause is the one
// structural anchor; the ALTER action keywords (SUSPEND / RESUME / ABORT /
// RENAME / SET / UNSET / ADD|REMOVE|DROP TABLES) anchor the ALTER grammar. The
// catalog/semantic layer, not the parser, validates that a property is real.

// warehouseOptionCap bounds an open-ended warehouse property run. Real
// statements carry well under a dozen properties; a far larger run signals a
// runaway parse and aborts loudly rather than spinning.
const warehouseOptionCap = 1024

// ---------------------------------------------------------------------------
// CREATE WAREHOUSE
// ---------------------------------------------------------------------------

// parseCreateWarehouseStmt parses the body of a
//
//	CREATE [ OR REPLACE | OR ALTER ] WAREHOUSE [ IF NOT EXISTS ] <name>
//	  [ WITH ] <prop> [ <prop> ... ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//
// statement. The CREATE keyword and the optional OR REPLACE / OR ALTER
// modifiers have already been consumed by parseCreateStmt; start is the Loc of
// the CREATE token, and cur is the WAREHOUSE keyword.
func (p *Parser) parseCreateWarehouseStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume WAREHOUSE

	stmt := &ast.CreateWarehouseStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Warehouse name.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional leading WITH preceding the property list (CREATE WAREHOUSE w WITH
	// WAREHOUSE_SIZE=...). It is purely cosmetic and only consumed when it does
	// NOT introduce a [WITH] TAG clause (WITH TAG anchors the tag clause, handled
	// by the property/tag loop below).
	if p.cur.Type == kwWITH && p.peekNext().Type != kwTAG {
		p.advance() // consume the cosmetic leading WITH
	}

	// Trailing clauses: an open-ended run of `KEY = value` properties interleaved
	// with the [WITH] TAG (...) clause, mirroring CREATE STAGE. The loop ends at a
	// statement boundary or any token that begins neither a property nor a tag
	// clause. The loop-guard aborts if an iteration fails to consume any token.
	for i := 0; ; i++ {
		if i >= warehouseOptionCap {
			return nil, p.syntaxErrorAtCur()
		}
		before := p.cur.Loc.Start
		if p.startsWarehouseTags() {
			tags, err := p.parseWarehouseWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
		} else if p.startsWarehouseOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
		} else {
			break
		}
		if p.cur.Loc.Start == before && p.cur.Type != tokEOF {
			// No forward progress on a non-EOF token: abort rather than spin.
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE INTERACTIVE WAREHOUSE
// ---------------------------------------------------------------------------

// interactiveWarehouseProps are the objectProperties of CREATE INTERACTIVE
// WAREHOUSE, and interactiveWarehouseParams its objectParams
// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-warehouse).
// WAREHOUSE_SIZE takes the sizes the page's parameters list, X2LARGE and
// X3LARGE among them; its examples quote the size (WAREHOUSE_SIZE = 'MEDIUM'),
// as the hyphenated sizes must be ('2X-LARGE'). AUTO_SUSPEND is
// { <num> | NULL }, which interactiveWarehouseValueOK checks.
var (
	interactiveWarehouseProps = objectProps{
		"WAREHOUSE_SIZE":      propStringOr("XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "X2LARGE", "XXXLARGE", "X3LARGE", "X4LARGE"),
		"MAX_CLUSTER_COUNT":   propNumber,
		"MIN_CLUSTER_COUNT":   propNumber,
		"AUTO_SUSPEND":        propNumber,
		"AUTO_RESUME":         propBool,
		"INITIALLY_SUSPENDED": propBool,
		"RESOURCE_MONITOR":    propName,
		"FALLBACK_WAREHOUSE":  propName,
		"COMMENT":             propString,
	}
	interactiveWarehouseParams = objectProps{
		"MAX_CONCURRENCY_LEVEL":               propNumber,
		"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS": propNumber,
		"STATEMENT_TIMEOUT_IN_SECONDS":        propNumber,
	}
)

// interactiveWarehouseValueOK reports whether o holds a value of the shape
// props documents for it.
func interactiveWarehouseValueOK(props objectProps, o *ast.CopyOption) bool {
	if o.Name == "AUTO_SUSPEND" && o.Words == "NULL" {
		return true
	}
	return props[o.Name].matches(o)
}

// parseCreateInteractiveWarehouseStmt parses
//
//	CREATE [ OR REPLACE ] INTERACTIVE WAREHOUSE [ IF NOT EXISTS ] <name>
//	  [ TABLES ( <table_name> [ , <table_name> ... ] ) ]
//	  [ [ WITH ] objectProperties ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ objectParams ]
//
// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-warehouse).
// The page writes IF NOT EXISTS only without OR REPLACE, and documents no
// OR ALTER. cur is INTERACTIVE; start is the CREATE token's Loc.
func (p *Parser) parseCreateInteractiveWarehouseStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume INTERACTIVE
	p.advance() // consume WAREHOUSE
	stmt := &ast.CreateWarehouseStmt{
		OrReplace:   orReplace,
		Interactive: true,
		Loc:         ast.Loc{Start: start.Start},
	}
	if !orReplace {
		if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
			return nil, err
		}
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if tokIsWord(p.cur, "TABLES") {
		p.advance() // consume TABLES
		tables, err := p.parseWarehouseTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables
	}

	// The properties come first, then the tag clause, then the parameters:
	// phase says which of them may still follow.
	const (
		phaseProps = iota
		phaseParams
	)
	phase := phaseProps
	seen := map[string]bool{}
	for {
		switch {
		case p.startsWarehouseTags():
			if phase != phaseProps || stmt.Tags != nil {
				return nil, p.syntaxErrorAtCur()
			}
			tags, err := p.parseWarehouseWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = tags
			phase = phaseParams
		case p.cur.Type == kwWITH:
			// WITH introduces the properties, so it comes before the first.
			if phase != phaseProps || len(stmt.Options) > 0 {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume WITH
			if _, ok := interactiveWarehouseProps[strings.ToUpper(p.cur.Str)]; !ok || !p.startsObjectProperty() {
				return nil, p.syntaxErrorAtCur()
			}
		case p.startsObjectProperty():
			opt, err := p.parseObjectProperty()
			if err != nil {
				return nil, err
			}
			props := interactiveWarehouseProps
			if _, ok := props[opt.Name]; !ok {
				props = interactiveWarehouseParams
				phase = phaseParams
			} else if phase != phaseProps {
				props = nil
			}
			if seen[opt.Name] || props == nil || !interactiveWarehouseValueOK(props, opt) {
				return nil, &ParseError{
					Position: opt.Loc.Start, End: opt.Loc.End,
					Message: "syntax error at or near " + opt.Name,
				}
			}
			seen[opt.Name] = true
			stmt.Options = append(stmt.Options, opt)
		default:
			stmt.Loc.End = p.prev.Loc.End
			return stmt, nil
		}
	}
}

// ---------------------------------------------------------------------------
// ALTER WAREHOUSE
// ---------------------------------------------------------------------------

// parseAlterWarehouseStmt parses ALTER WAREHOUSE [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the WAREHOUSE keyword.
//
//	SUSPEND
//	RESUME [ IF SUSPENDED ]
//	ABORT ALL QUERIES
//	RENAME TO <new_name>
//	SET <prop> [ <prop> ... ]               -- open-ended KEY = value params
//	UNSET <key> [ , ... ]
//	ADD TABLES ( <id> [ , ... ] )           -- Unistore/interactive variant
//	{ REMOVE | DROP } TABLES ( <id> [ , ... ] )
func (p *Parser) parseAlterWarehouseStmt() (ast.Node, error) {
	whTok := p.advance() // consume WAREHOUSE (anchors Loc.Start at the object keyword)
	stmt := &ast.AlterWarehouseStmt{Loc: ast.Loc{Start: whTok.Loc.Start}}

	// Optional IF EXISTS.
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	// Warehouse name, which SUSPEND, RESUME and ABORT ALL QUERIES may leave
	// out to act on the session's current warehouse.
	if !p.startsNamelessWarehouseAction() {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// Action branch.
	switch p.cur.Type {
	case kwENABLE, kwDISABLE:
		if stmt.Name == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterWarehouseEnable
		if p.advance().Type == kwDISABLE {
			stmt.Action = ast.AlterWarehouseDisable
		}

	case kwSUSPEND:
		p.advance() // consume SUSPEND
		stmt.Action = ast.AlterWarehouseSuspend

	case kwRESUME:
		p.advance() // consume RESUME
		stmt.Action = ast.AlterWarehouseResume
		// Optional IF SUSPENDED.
		if p.cur.Type == kwIF && p.peekNext().Type == kwSUSPENDED {
			p.advance() // consume IF
			p.advance() // consume SUSPENDED
			stmt.ResumeIfSuspended = true
		}

	case kwABORT:
		// ABORT ALL QUERIES
		p.advance() // consume ABORT
		if _, err := p.expect(kwALL); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwQUERIES); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterWarehouseAbort

	case kwRENAME:
		// RENAME TO <new_name>
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterWarehouseRename
		stmt.NewName = newName

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			// SET TAG <tag> = '<value>' [ , ... ] — unparenthesized assignment list
			// (reuses the stage SET TAG helper). Anchored before the open-ended
			// property run so a literal property named TAG cannot shadow it.
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterWarehouseSetTag
			stmt.Tags = tags
		} else {
			// SET <prop> [ <prop> ... ] — open-ended KEY = value params.
			opts, err := p.parseWarehouseOptions()
			if err != nil {
				return nil, err
			}
			if len(opts) == 0 {
				// SET with nothing settable is a syntax error.
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterWarehouseSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			// UNSET TAG <tag> [ , ... ] — unparenthesized name list (reuses the
			// stage UNSET TAG helper). Anchored before the open-ended key run so a
			// `UNSET TAG t` is not mis-read as unsetting keys "TAG" and "t".
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterWarehouseUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET <key> [ , ... ]
			keys, err := p.parseWarehouseUnsetKeys()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterWarehouseUnset
			stmt.UnsetKeys = keys
		}

	case kwADD:
		// ADD TABLES ( <id> [ , ... ] )
		p.advance() // consume ADD
		if _, err := p.expect(kwTABLES); err != nil {
			return nil, err
		}
		ids, err := p.parseWarehouseTableList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterWarehouseAddTables
		stmt.Tables = ids

	case kwREMOVE, kwDROP:
		// { REMOVE | DROP } TABLES ( <id> [ , ... ] )
		p.advance() // consume REMOVE / DROP
		if _, err := p.expect(kwTABLES); err != nil {
			return nil, err
		}
		ids, err := p.parseWarehouseTableList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterWarehouseRemoveTables
		stmt.Tables = ids

	default:
		return nil, p.syntaxErrorAtCur()
	}
	if stmt.Name == nil && stmt.Action != ast.AlterWarehouseSuspend &&
		stmt.Action != ast.AlterWarehouseResume && stmt.Action != ast.AlterWarehouseAbort {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsNamelessWarehouseAction reports whether cur begins an ALTER WAREHOUSE
// action written without a warehouse name
// (https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse):
//
//	ALTER WAREHOUSE [ IF EXISTS ] [ <name> ] { SUSPEND | RESUME [ IF SUSPENDED ] }
//	ALTER WAREHOUSE [ IF EXISTS ] [ <name> ] ABORT ALL QUERIES
//
// SUSPEND, RESUME and ABORT are not reserved, so each also names a warehouse:
// it is the action only where no action follows it (SUSPEND or RESUME at the
// statement's end, RESUME IF SUSPENDED, ABORT ALL). It consumes nothing.
func (p *Parser) startsNamelessWarehouseAction() bool {
	next := p.peekNext().Type
	end := next == tokEOF || next == ';' || next == tokFlow
	switch p.cur.Type {
	case kwSUSPEND:
		return end
	case kwRESUME:
		if end {
			return true
		}
		if next != kwIF {
			return false
		}
		snap := p.snapshot()
		defer p.restore(snap)
		p.advance() // RESUME
		p.advance() // IF
		return p.cur.Type == kwSUSPENDED
	case kwABORT:
		return next == kwALL
	}
	return false
}

// ---------------------------------------------------------------------------
// Shared warehouse helpers
// ---------------------------------------------------------------------------

// parseWarehouseOptions parses a run of zero or more warehouse properties, each
// an open-ended `KEY = <value>` pair (reusing parseCopyOption). The run
// continues while the current token begins another property and terminates at a
// statement boundary, the WITH / TAG clause, or any token that does not begin a
// property. A loop-guard aborts if an iteration consumes no tokens.
func (p *Parser) parseWarehouseOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for i := 0; p.startsWarehouseOption(); i++ {
		if i >= warehouseOptionCap {
			return nil, p.syntaxErrorAtCur()
		}
		before := p.cur.Loc.Start
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if p.cur.Loc.Start == before && p.cur.Type != tokEOF {
			return nil, p.syntaxErrorAtCur()
		}
	}
	return opts, nil
}

// startsWarehouseOption reports whether the current token can begin a warehouse
// property. It is the COPY option predicate minus WITH and TAG, which anchor the
// trailing [WITH] TAG clause and must not be swallowed as a property name.
func (p *Parser) startsWarehouseOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG:
		return false
	}
	return p.startsCopyOption()
}

// startsWarehouseTags reports whether the current position begins a [WITH] TAG
// clause: either the TAG keyword directly, or WITH immediately followed by TAG.
func (p *Parser) startsWarehouseTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
}

// parseWarehouseWithTags parses a [WITH] TAG (...) clause and returns the tag
// assignments. It accepts both the WITH TAG and the bare TAG spellings. The
// caller must have confirmed startsWarehouseTags.
func (p *Parser) parseWarehouseWithTags() ([]*ast.TagAssignment, error) {
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH (peekNext == TAG guaranteed by startsWarehouseTags)
	}
	return p.parseTagAssignments()
}

// parseWarehouseUnsetKeys parses an UNSET <key> [ , ... ] property-name list.
// Each key is a single name word (identifier or keyword), uppercased; a comma
// separates keys. Consumes at least one key.
func (p *Parser) parseWarehouseUnsetKeys() ([]string, error) {
	var keys []string
	for {
		if !p.isOptionWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		keys = append(keys, strings.ToUpper(p.advance().Str))
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return keys, nil
}

// parseWarehouseTableList parses a parenthesized comma-separated identifier list
// for ADD / REMOVE / DROP TABLES ( <id> [ , ... ] ). On entry cur is '('. At
// least one identifier is required.
func (p *Parser) parseWarehouseTableList() ([]*ast.ObjectName, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var ids []*ast.ObjectName
	for {
		id, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return ids, nil
}
