package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseCreateMTMV parses createMaterializedViewStatement (StarRocks.g4:646):
//
//	CREATE MATERIALIZED VIEW (IF NOT EXISTS)? qualifiedName
//	    ('(' columnNameWithComment (',' columnNameWithComment)* (',' indexDesc)* ')')?
//	    comment?
//	    materializedViewDesc*
//	    AS queryStatement
//	materializedViewDesc : PARTITION BY mvPartitionExprs | distributionDesc
//	    | orderByDesc | refreshSchemeDesc | properties
//
// On entry, CREATE has been consumed and cur is MATERIALIZED.
func (p *Parser) parseCreateMTMV(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	stmt := &ast.CreateMTMVStmt{}
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
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Kind == int('(') {
		if err := p.parseMTMVColumns(stmt); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.advance().Str
	}

	for {
		switch p.cur.Kind {
		case kwPARTITION:
			if err := p.parseMVPartitionExprs(stmt); err != nil {
				return nil, err
			}
			continue
		case kwDISTRIBUTED:
			if stmt.DistributedBy, err = p.parseDistributedBy(); err != nil {
				return nil, err
			}
			continue
		case kwORDER:
			// orderByDesc (:436): ORDER BY identifierList
			p.advance() // consume ORDER
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			if stmt.OrderBy, err = p.parseIdentifierList(); err != nil {
				return nil, err
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			continue
		case kwREFRESH:
			if err := p.parseMVRefreshScheme(stmt); err != nil {
				return nil, err
			}
			continue
		case kwPROPERTIES:
			if stmt.Properties, err = p.parseProperties(); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	var queryNode ast.Node
	if p.cur.Kind == kwWITH {
		wq, err := p.parseWithSelect()
		if err != nil {
			return nil, err
		}
		queryNode = wq
	} else {
		query, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		queryNode, err = p.parseSetOpTail(query)
		if err != nil {
			return nil, err
		}
	}
	stmt.Query = queryNode

	stmt.Loc = startLoc.Merge(ast.NodeLoc(stmt.Query))
	return stmt, nil
}

// parseMTMVColumns reads the column list after the view's name:
//
//	'(' columnNameWithComment (',' columnNameWithComment)* (',' indexDesc)* ')'
//	columnNameWithComment : identifier comment?
//	indexDesc : INDEX identifier identifierList (indexType propertyList?)? comment?
//	indexType : USING (BITMAP | GIN | NGRAMBF | VECTOR)
func (p *Parser) parseMTMVColumns(stmt *ast.CreateMTMVStmt) error {
	p.advance() // consume '('
	for {
		start := p.cur.Loc
		if p.cur.Kind == kwINDEX {
			if len(stmt.Columns) == 0 {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume INDEX
			if _, _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if _, err := p.expect(int('(')); err != nil {
				return err
			}
			if _, err := p.parseIdentifierList(); err != nil {
				return err
			}
			if _, err := p.expect(int(')')); err != nil {
				return err
			}
			if p.cur.Kind == kwUSING {
				p.advance() // consume USING
				word := strings.ToUpper(p.cur.Str)
				if p.cur.Kind != kwBITMAP && (p.cur.Kind != tokIdent ||
					(word != "GIN" && word != "NGRAMBF" && word != "VECTOR")) {
					return p.syntaxErrorAtCur()
				}
				p.advance()
				if p.cur.Kind == int('(') {
					if _, err := p.parsePropertyList(); err != nil {
						return err
					}
				}
			}
			if p.cur.Kind == kwCOMMENT {
				p.advance() // consume COMMENT
				if p.cur.Kind != tokString {
					return p.syntaxErrorAtCur()
				}
				p.advance()
			}
			stmt.Indexes = append(stmt.Indexes, p.input[start.Start-p.baseOffset:p.prev.Loc.End-p.baseOffset])
		} else {
			if len(stmt.Indexes) > 0 {
				return p.syntaxErrorAtCur()
			}
			name, _, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			col := &ast.ViewColumn{Name: name, Loc: start}
			if p.cur.Kind == kwCOMMENT {
				p.advance() // consume COMMENT
				if p.cur.Kind != tokString {
					return p.syntaxErrorAtCur()
				}
				col.Comment = p.advance().Str
			}
			col.Loc.End = p.prev.Loc.End
			stmt.Columns = append(stmt.Columns, col)
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(int(')'))
	return err
}

// parseMVPartitionExprs reads PARTITION BY mvPartitionExprs (StarRocks.g4:654):
//
//	primaryExpression | '(' primaryExpression (',' primaryExpression)* ')'
func (p *Parser) parseMVPartitionExprs(stmt *ast.CreateMTMVStmt) error {
	start := p.advance().Loc // consume PARTITION
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	desc := &ast.PartitionDesc{}
	// StarRocks refuses a partition expression that names no column of the
	// query: testdata/compat.tsv refuses PARTITION BY CURRENT_USER and
	// PARTITION BY (CURRENT_USER) beside AS SELECT CURRENT_USER. omni, which
	// does not resolve the query's columns here, reads a primaryExpression
	// that starts with a name: a column reference or a function call.
	read := func() error {
		exprStart := p.cur.Loc
		if !isIdentifierToken(p.cur.Kind) {
			return p.syntaxErrorAtCur()
		}
		if _, err := p.parsePrimaryExpr(); err != nil {
			return err
		}
		desc.FuncExprs = append(desc.FuncExprs, p.input[exprStart.Start-p.baseOffset:p.prev.Loc.End-p.baseOffset])
		return nil
	}
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		for {
			if err := read(); err != nil {
				return err
			}
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
	} else if err := read(); err != nil {
		return err
	}
	desc.Loc = ast.Loc{Start: start.Start, End: p.prev.Loc.End}
	stmt.PartitionBy = desc
	return nil
}

// parseMVRefreshScheme reads refreshSchemeDesc (StarRocks.g4:2687):
//
//	REFRESH (IMMEDIATE | DEFERRED)?
//	    (ASYNC | ASYNC (START '(' string ')')? EVERY '(' interval ')' | INCREMENTAL | MANUAL)
func (p *Parser) parseMVRefreshScheme(stmt *ast.CreateMTMVStmt) error {
	p.advance() // consume REFRESH
	if p.cur.Kind == kwIMMEDIATE || p.cur.Kind == kwDEFERRED {
		stmt.BuildMode = strings.ToUpper(p.advance().Str)
	}
	switch {
	case p.cur.Kind == kwINCREMENTAL, p.cur.Kind == kwMANUAL:
		stmt.RefreshMethod = strings.ToUpper(p.advance().Str)
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "ASYNC"):
		p.advance()
		stmt.RefreshMethod = "ASYNC"
		if p.cur.Kind == kwSTART || p.cur.Kind == kwEVERY {
			trigger, err := p.parseRefreshAsyncSchedule()
			if err != nil {
				return err
			}
			stmt.RefreshTrigger = trigger
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseMTMVRefreshTrigger parses the ON clause of CREATE MATERIALIZED VIEW:
//
//	ON SCHEDULE EVERY interval [STARTS timestamp]
//	ON COMMIT
//	ON MANUAL
//
// cur is ON on entry; ON is consumed here.
func (p *Parser) parseMTMVRefreshTrigger() (*ast.MTMVRefreshTrigger, error) {
	startLoc := p.cur.Loc
	p.advance() // consume ON

	trigger := &ast.MTMVRefreshTrigger{Loc: startLoc}

	switch p.cur.Kind {
	case kwSCHEDULE:
		trigger.OnSchedule = true
		p.advance() // consume SCHEDULE

		// EVERY interval_expr
		if _, err := p.expect(kwEVERY); err != nil {
			return nil, err
		}

		// Collect interval: integer + unit keyword (e.g., "1 DAY")
		// The interval is stored as raw text.
		intervalStart := p.cur.Loc.Start
		var parts []string
		// First token should be an integer
		if p.cur.Kind == tokInt {
			parts = append(parts, p.cur.Str)
			p.advance()
		}
		// Unit keyword (DAY, HOUR, MINUTE, SECOND, WEEK, MONTH, YEAR, etc.)
		if p.cur.Kind >= 700 || p.cur.Kind == tokIdent {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		}
		if len(parts) > 0 {
			trigger.Interval = strings.Join(parts, " ")
		}
		_ = intervalStart

		// Optional STARTS timestamp_string
		if p.cur.Kind == kwSTARTS {
			p.advance() // consume STARTS
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			trigger.StartsAt = p.cur.Str
			trigger.Loc.End = p.cur.Loc.End
			p.advance()
		}

	case kwCOMMIT:
		trigger.OnCommit = true
		trigger.Loc.End = p.cur.Loc.End
		p.advance()

	case kwMANUAL:
		trigger.OnManual = true
		trigger.Loc.End = p.cur.Loc.End
		p.advance()

	default:
		// Tolerate unknown ON clause variant.
		trigger.Loc.End = p.cur.Loc.End
		p.advance()
	}

	return trigger, nil
}

// parseRefreshAsyncSchedule parses the StarRocks async-MV scheduling tail that
// attaches directly to REFRESH [IMMEDIATE|DEFERRED] ASYNC (no ON SCHEDULE):
//
//	[START '(' string ')'] EVERY '(' INTERVAL n unit ')'
//
// This differs from the doris ON SCHEDULE form (bare EVERY <n unit>, STARTS
// <string>) parsed by parseMTMVRefreshTrigger. cur is START or EVERY on entry.
func (p *Parser) parseRefreshAsyncSchedule() (*ast.MTMVRefreshTrigger, error) {
	trigger := &ast.MTMVRefreshTrigger{OnSchedule: true, Loc: p.cur.Loc}

	// Optional START '(' string ')'
	if p.cur.Kind == kwSTART {
		p.advance() // consume START
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		trigger.StartsAt = p.cur.Str
		p.advance()
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		trigger.Loc.End = p.prev.Loc.End
	}

	// EVERY '(' INTERVAL n unit ')' — mandatory once the async tail is present
	// (the grammar arm is `(START('...'))? EVERY '(' interval ')'`).
	if _, err := p.expect(kwEVERY); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwINTERVAL); err != nil {
		return nil, err
	}
	if p.cur.Kind != tokInt {
		return nil, p.syntaxErrorAtCur()
	}
	parts := []string{p.cur.Str}
	p.advance()
	// Unit keyword (DAY, HOUR, MINUTE, ...) or ident (MILLISECOND/MICROSECOND).
	if p.cur.Kind < 700 && p.cur.Kind != tokIdent {
		return nil, p.syntaxErrorAtCur()
	}
	parts = append(parts, strings.ToUpper(p.cur.Str))
	p.advance()
	trigger.Interval = strings.Join(parts, " ")
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	trigger.Loc.End = closeTok.Loc.End

	return trigger, nil
}

// parseAlterMTMV parses alterMaterializedViewStatement (StarRocks.g4:675).
// On entry, ALTER has been consumed and cur is MATERIALIZED.
//
//	ALTER MATERIALIZED VIEW mvName=qualifiedName
//	    ( refreshSchemeDesc | tableRenameClause | modifyPropertiesClause | swapTableClause )
//	ALTER MATERIALIZED VIEW mvName=qualifiedName statusDesc
//
// Best-effort parsing also takes the Doris forms the fork began with
// (REFRESH COMPLETE / AUTO / NEVER, REPLACE WITH MATERIALIZED VIEW).
func (p *Parser) parseAlterMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.AlterMTMVStmt{}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.cur.Kind == kwRENAME:
		// StarRocks.g4:1012 tableRenameClause: RENAME identifier
		p.advance() // consume RENAME
		newName, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = &ast.ObjectName{Parts: []string{newName}, Loc: loc}

	case p.cur.Kind == kwREFRESH:
		if err := p.parseAlterMTMVRefresh(stmt); err != nil {
			return nil, err
		}

	case p.cur.Kind == kwSET:
		// StarRocks.g4:1020 modifyPropertiesClause: SET propertyList
		p.advance() // consume SET
		props, err := p.parseMTMVPropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props

	case p.alterAt("SWAP"):
		// StarRocks.g4:1016 swapTableClause: SWAP WITH identifier
		p.advance()
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		other, loc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.SwapWith = &ast.ObjectName{Parts: []string{other}, Loc: loc}

	case p.alterAt("ACTIVE") || p.alterAt("INACTIVE"):
		// StarRocks.g4:2694 statusDesc: ACTIVE | INACTIVE
		stmt.Status = strings.ToUpper(p.advance().Str)

	case p.cur.Kind == kwREPLACE && !p.strictTrailing:
		// REPLACE WITH MATERIALIZED VIEW other_mv (Doris; best-effort only)
		p.advance() // consume REPLACE
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwMATERIALIZED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Replace = true
		stmt.ReplaceTarget = target

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(p.prev.Loc)
	return stmt, nil
}

// parseAlterMTMVRefresh reads refreshSchemeDesc (StarRocks.g4:2687):
//
//	REFRESH (IMMEDIATE | DEFERRED)? (ASYNC
//	    | ASYNC (START '(' string ')')? EVERY '(' interval ')'
//	    | INCREMENTAL
//	    | MANUAL)
//
// The interval reads INTERVAL INTEGER_VALUE unitIdentifier: the engine
// refuses EVERY (INTERVAL CURRENT_USER YEAR). cur is REFRESH.
func (p *Parser) parseAlterMTMVRefresh(stmt *ast.AlterMTMVStmt) error {
	p.advance() // consume REFRESH
	if !p.strictTrailing {
		switch p.cur.Kind {
		case kwCOMPLETE, kwAUTO, kwNEVER:
			// Doris refresh methods; best-effort only.
			stmt.RefreshMethod = strings.ToUpper(p.advance().Str)
			return nil
		}
	}
	var method []string
	if p.alterAt("IMMEDIATE") || p.alterAt("DEFERRED") {
		method = append(method, strings.ToUpper(p.advance().Str))
	}
	switch {
	case p.alterAt("ASYNC"):
		method = append(method, "ASYNC")
		p.advance()
		trigger := &ast.MTMVRefreshTrigger{OnSchedule: true, Loc: p.cur.Loc}
		if p.alterAt("START") {
			p.advance()
			if _, err := p.expect(int('(')); err != nil {
				return err
			}
			start, err := p.alterExpectString()
			if err != nil {
				return err
			}
			trigger.StartsAt = start.Str
			if _, err := p.expect(int(')')); err != nil {
				return err
			}
			if !p.alterAt("EVERY") {
				return p.syntaxErrorAtCur()
			}
		}
		if p.alterAt("EVERY") {
			p.advance()
			if _, err := p.expect(int('(')); err != nil {
				return err
			}
			if _, err := p.expect(kwINTERVAL); err != nil {
				return err
			}
			n, err := p.alterExpectInt()
			if err != nil {
				return err
			}
			if !alterIsUnitIdentifier(p.cur) {
				return p.syntaxErrorAtCur()
			}
			trigger.Interval = n.Str + " " + strings.ToUpper(p.advance().Str)
			if _, err := p.expect(int(')')); err != nil {
				return err
			}
			trigger.Loc.End = p.prev.Loc.End
			stmt.RefreshTrigger = trigger
		}
	case p.cur.Kind == kwINCREMENTAL, p.alterAt("MANUAL"):
		method = append(method, strings.ToUpper(p.advance().Str))
	default:
		return p.syntaxErrorAtCur()
	}
	stmt.RefreshMethod = strings.Join(method, " ")
	return nil
}

// parseDropMTMV parses a DROP MATERIALIZED VIEW statement.
// On entry, DROP has been consumed and cur is MATERIALIZED.
//
// Syntax:
//
//	DROP MATERIALIZED VIEW [IF EXISTS] mv_name [ON base_table]
func (p *Parser) parseDropMTMV(startLoc ast.Loc) (ast.Node, error) {
	// Consume MATERIALIZED
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	// Consume VIEW
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.DropMTMVStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// MV name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))

	// Optional ON base_table (for sync MVs)
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		baseName, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.OnBase = baseName
		stmt.Loc = startLoc.Merge(ast.NodeLoc(baseName))
	}

	return stmt, nil
}

// parseRefreshMTMV parses refreshMaterializedViewStatement (StarRocks.g4:684):
//
//	REFRESH MATERIALIZED VIEW qualifiedName
//	    (PARTITION (partitionRangeDesc | listPartitionValues))? FORCE?
//	    (WITH (SYNC | ASYNC) MODE)? (WITH PRIORITY INTEGER_VALUE)?
//	partitionRangeDesc  : START '(' string ')' END '(' string ')'
//	listPartitionValues : '(' (NULL | string) (',' (NULL | string))* ')'
//	    | '(' ('(' (NULL | string) (',' ...)* ')') (',' ...)* ')'
//
// On entry, REFRESH has been consumed and cur is MATERIALIZED.
func (p *Parser) parseRefreshMTMV(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	stmt := &ast.RefreshMTMVStmt{}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
		start := p.cur.Loc
		if p.cur.Kind == kwSTART {
			for _, kw := range []int{kwSTART, kwEND} {
				if _, err := p.expect(kw); err != nil {
					return nil, err
				}
				if _, err := p.expect(int('(')); err != nil {
					return nil, err
				}
				if p.cur.Kind != tokString {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
			}
		} else if err := p.parseListPartitionValues(true); err != nil {
			return nil, err
		}
		stmt.Partitions = []string{p.input[start.Start-p.baseOffset : p.prev.Loc.End-p.baseOffset]}
	}
	if p.cur.Kind == kwFORCE {
		p.advance()
		stmt.Force = true
	}
	if p.cur.Kind == kwWITH && p.peekNext().Kind != tokEOF {
		next := p.peekNext()
		if next.Kind == kwSYNC || (next.Kind == tokIdent && strings.EqualFold(next.Str, "ASYNC")) {
			p.advance() // consume WITH
			stmt.Mode = strings.ToUpper(p.advance().Str)
			if p.cur.Kind != tokIdent || !strings.EqualFold(p.cur.Str, "MODE") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		}
	}
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if p.cur.Kind != tokIdent || !strings.EqualFold(p.cur.Str, "PRIORITY") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Priority = p.advance().Ival
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseListPartitionValues reads listPartitionValues (StarRocks.g4:2624): a
// parenthesized list of NULL or strings, or, when nested is set, a
// parenthesized list of such lists.
func (p *Parser) parseListPartitionValues(nested bool) error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	multi := nested && p.cur.Kind == int('(')
	for {
		if multi {
			if err := p.parseListPartitionValues(false); err != nil {
				return err
			}
		} else if p.cur.Kind == kwNULL || p.cur.Kind == tokString {
			p.advance()
		} else {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(int(')'))
	return err
}

// parseMTMVPropertyList parses a parenthesised key=value list:
//
//	( "key" = "val" [, "key" = "val"] ... )
//
// cur must be '(' on entry. Used for ALTER MATERIALIZED VIEW ... SET (...).
func (p *Parser) parseMTMVPropertyList() ([]*ast.Property, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	return p.parsePropertyItems(p.parseMTMVProperty)
}

// parseMTMVProperty reads one item of a materialized view's property list,
// whose key may be written bare as well as quoted.
func (p *Parser) parseMTMVProperty() (*ast.Property, error) {
	startLoc := p.cur.Loc

	var key string
	if p.cur.Kind == tokString {
		key = p.cur.Str
		p.advance()
	} else {
		var err error
		key, _, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}

	if p.cur.Kind != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	endLoc := p.cur.Loc
	p.advance()

	return &ast.Property{
		Key:   key,
		Value: val,
		Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}
