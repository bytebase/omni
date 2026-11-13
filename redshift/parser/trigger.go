package parser

import (
	nodes "github.com/bytebase/omni/redshift/ast"
	"strconv"
)

func (p *Parser) parseCreateTrigStmt(replace bool) (*nodes.CreateTrigStmt, error) {
	if p.cur.Type == CONSTRAINT {
		return p.parseCreateConstraintTrigger(replace)
	}
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}
	trigname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	timing, err := p.parseTriggerActionTime()
	if err != nil {
		return nil, err
	}
	events, columns, err := p.parseTriggerEvents()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	relStart := p.pos()
	relNames, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	relation := makeRangeVarFromAnyName(relNames, nodes.Loc{Start: relStart, End: p.prev.End})
	transitionRels, err := p.parseTriggerReferencing()
	if err != nil {
		return nil, err
	}
	row, err := p.parseTriggerForSpec()
	if err != nil {
		return nil, err
	}
	whenClause, err := p.parseTriggerWhen()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(EXECUTE); err != nil {
		return nil, err
	}
	if err := p.parseFunctionOrProcedure(); err != nil {
		return nil, err
	}
	funcname, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	args, err := p.parseTriggerFuncArgs()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateTrigStmt{
		Replace: replace, IsConstraint: false, Trigname: trigname,
		Relation: relation, Funcname: funcname, Args: args, Row: row,
		Timing: int16(timing), Events: int16(events), Columns: columns,
		WhenClause: whenClause, TransitionRels: transitionRels,
		Deferrable: false, Initdeferred: false,
	}, nil
}

func (p *Parser) parseCreateConstraintTrigger(replace bool) (*nodes.CreateTrigStmt, error) {
	if _, err := p.expect(CONSTRAINT); err != nil {
		return nil, err
	}
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}
	trigname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AFTER); err != nil {
		return nil, err
	}
	events, columns, err := p.parseTriggerEvents()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	relStart := p.pos()
	relNames, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	relation := makeRangeVarFromAnyName(relNames, nodes.Loc{Start: relStart, End: p.prev.End})
	var constrrel *nodes.RangeVar
	if p.cur.Type == FROM {
		p.advance()
		fromStart := p.pos()
		fromNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		constrrel = makeRangeVarFromAnyName(fromNames, nodes.Loc{Start: fromStart, End: p.prev.End})
	}
	casBits, err := p.parseKindAttributes("TRIGGER", true, false, false)
	if err != nil {
		return nil, err
	}
	deferrable := (casBits & int64(nodes.CAS_DEFERRABLE)) != 0
	initdeferred := (casBits & int64(nodes.CAS_INITIALLY_DEFERRED)) != 0
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if _, err := p.expect(EACH); err != nil {
		return nil, err
	}
	if _, err := p.expect(ROW); err != nil {
		return nil, err
	}
	whenClause, err := p.parseTriggerWhen()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(EXECUTE); err != nil {
		return nil, err
	}
	if err := p.parseFunctionOrProcedure(); err != nil {
		return nil, err
	}
	funcname, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	args, err := p.parseTriggerFuncArgs()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateTrigStmt{
		Replace: replace, IsConstraint: true, Trigname: trigname,
		Relation: relation, Funcname: funcname, Args: args, Row: true,
		Timing: int16(nodes.TRIGGER_TYPE_AFTER), Events: int16(events),
		Columns: columns, WhenClause: whenClause,
		Deferrable: deferrable, Initdeferred: initdeferred, Constrrel: constrrel,
	}, nil
}

// parseTriggerActionTime parses TriggerActionTime: BEFORE | AFTER | INSTEAD OF.
func (p *Parser) parseTriggerActionTime() (int64, error) {
	if p.collectMode() {
		for _, t := range []int{BEFORE, AFTER, INSTEAD} {
			p.addTokenCandidate(t)
		}
		return 0, errCollecting
	}
	switch p.cur.Type {
	case BEFORE:
		p.advance()
		return int64(nodes.TRIGGER_TYPE_BEFORE), nil
	case AFTER:
		p.advance()
		return int64(nodes.TRIGGER_TYPE_AFTER), nil
	case INSTEAD:
		p.advance()
		if _, err := p.expect(OF); err != nil {
			return 0, err
		}
		return int64(nodes.TRIGGER_TYPE_INSTEAD), nil
	}
	return 0, p.syntaxErrorAtCur()
}

// parseFunctionOrProcedure consumes FUNCTION_or_PROCEDURE, which
// CreateTrigStmt and CreateEventTrigStmt require after EXECUTE.
func (p *Parser) parseFunctionOrProcedure() error {
	if p.collectMode() {
		p.addTokenCandidate(FUNCTION)
		p.addTokenCandidate(PROCEDURE)
		return errCollecting
	}
	if p.cur.Type != FUNCTION && p.cur.Type != PROCEDURE {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

func (p *Parser) parseTriggerEvents() (int64, *nodes.List, error) {
	events, cols, err := p.parseTriggerOneEvent()
	if err != nil {
		return 0, nil, err
	}
	for p.cur.Type == OR {
		p.advance()
		ev2, cols2, err := p.parseTriggerOneEvent()
		if err != nil {
			return 0, nil, err
		}
		if events&ev2 != 0 {
			// gram.y's TriggerEvents action raises parser_yyerror, which
			// reports the last token the parser read: the event keyword
			// itself for INSERT, DELETE and TRUNCATE, and the token after
			// an UPDATE, which bison reads to look for OF or a column.
			pos := p.cur.Loc
			if ev2 != int64(nodes.TRIGGER_TYPE_UPDATE) {
				pos = p.prev.Loc
			}
			return 0, nil, &ParseError{Message: "duplicate trigger events specified", Position: pos}
		}
		events |= ev2
		cols = concatColumnLists(cols, cols2)
	}
	return events, cols, nil
}

func (p *Parser) parseTriggerOneEvent() (int64, *nodes.List, error) {
	switch p.cur.Type {
	case INSERT:
		p.advance()
		return int64(nodes.TRIGGER_TYPE_INSERT), nil, nil
	case DELETE_P:
		p.advance()
		return int64(nodes.TRIGGER_TYPE_DELETE), nil, nil
	case UPDATE:
		p.advance()
		if p.cur.Type == OF {
			p.advance()
			cols, err := p.parseColumnList()
			if err != nil {
				return 0, nil, err
			}
			return int64(nodes.TRIGGER_TYPE_UPDATE), cols, nil
		}
		return int64(nodes.TRIGGER_TYPE_UPDATE), nil, nil
	case TRUNCATE:
		p.advance()
		return int64(nodes.TRIGGER_TYPE_TRUNCATE), nil, nil
	}
	if p.collectMode() {
		for _, t := range []int{INSERT, DELETE_P, UPDATE, TRUNCATE} {
			p.addTokenCandidate(t)
		}
		return 0, nil, errCollecting
	}
	return 0, nil, p.syntaxErrorAtCur()
}

func concatColumnLists(a, b *nodes.List) *nodes.List {
	if a == nil && b == nil {
		return nil
	}
	var items []nodes.Node
	if a != nil {
		items = append(items, a.Items...)
	}
	if b != nil {
		items = append(items, b.Items...)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseTriggerReferencing parses TriggerReferencing: REFERENCING
// TriggerTransitions, a list of one or more TriggerTransition.
func (p *Parser) parseTriggerReferencing() (*nodes.List, error) {
	if p.cur.Type != REFERENCING {
		return nil, nil
	}
	p.advance()
	var items []nodes.Node
	for {
		tt, err := p.parseTriggerTransition()
		if err != nil {
			return nil, err
		}
		if tt == nil {
			break
		}
		items = append(items, tt)
	}
	if len(items) == 0 {
		if p.collectMode() {
			p.addTokenCandidate(NEW)
			p.addTokenCandidate(OLD)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseTriggerTransition() (*nodes.TriggerTransition, error) {
	loc := p.pos()
	var isNew bool
	switch p.cur.Type {
	case NEW:
		isNew = true
		p.advance()
	case OLD:
		isNew = false
		p.advance()
	// optional-probe: trigger transition list terminator — referencing loop breaks on nil.
	default:
		return nil, nil
	}
	// TransitionRowOrTable: TABLE | ROW, required (gram.y TransitionRowOrTable).
	isTable := false
	switch p.cur.Type {
	case TABLE:
		isTable = true
		p.advance()
	case ROW:
		p.advance()
	default:
		if p.collectMode() {
			p.addTokenCandidate(TABLE)
			p.addTokenCandidate(ROW)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == AS {
		p.advance()
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.TriggerTransition{Name: name, IsNew: isNew, IsTable: isTable, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseTriggerForSpec parses TriggerForSpec: FOR [EACH] ROW | STATEMENT,
// or nothing.
func (p *Parser) parseTriggerForSpec() (bool, error) {
	if p.cur.Type != FOR {
		return false, nil
	}
	p.advance()
	if p.cur.Type == EACH {
		p.advance()
	}
	switch p.cur.Type {
	case ROW:
		p.advance()
		return true, nil
	case STATEMENT:
		p.advance()
		return false, nil
	}
	if p.collectMode() {
		if p.prev.Type == FOR {
			p.addTokenCandidate(EACH)
		}
		p.addTokenCandidate(ROW)
		p.addTokenCandidate(STATEMENT)
		return false, errCollecting
	}
	return false, p.syntaxErrorAtCur()
}

func (p *Parser) parseTriggerWhen() (nodes.Node, error) {
	if p.cur.Type != WHEN {
		return nil, nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseTriggerFuncArgs parses TriggerFuncArgs:
//
//	TriggerFuncArg | TriggerFuncArgs ',' TriggerFuncArg | /*EMPTY*/
//
// The empty alternative heads the left recursion, so a list may start with
// a comma: a(, 1) passes one argument, while a(,) is rejected at ')'.
func (p *Parser) parseTriggerFuncArgs() (*nodes.List, error) {
	if p.cur.Type == ')' {
		return nil, nil
	}
	var items []nodes.Node
	if p.cur.Type != ',' {
		arg, err := p.parseTriggerFuncArg()
		if err != nil {
			return nil, err
		}
		items = append(items, arg)
	}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseTriggerFuncArg()
		if err != nil {
			return nil, err
		}
		items = append(items, arg)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseTriggerFuncArg() (nodes.Node, error) {
	switch p.cur.Type {
	case ICONST:
		tok := p.advance()
		return &nodes.String{Str: strconv.FormatInt(tok.Ival, 10)}, nil
	case FCONST:
		tok := p.advance()
		return &nodes.String{Str: tok.Str}, nil
	case SCONST:
		tok := p.advance()
		return &nodes.String{Str: tok.Str}, nil
	default:
		label, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		return &nodes.String{Str: label}, nil
	}
}

func (p *Parser) parseCreateEventTrigStmt() (*nodes.CreateEventTrigStmt, error) {
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}
	trigname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	eventname, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	var whenclause *nodes.List
	if p.cur.Type == WHEN {
		p.advance()
		var items []nodes.Node
		item, err := p.parseEventTriggerWhenItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		for p.cur.Type == AND {
			p.advance()
			item, err := p.parseEventTriggerWhenItem()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		whenclause = &nodes.List{Items: items}
	}
	if _, err := p.expect(EXECUTE); err != nil {
		return nil, err
	}
	if err := p.parseFunctionOrProcedure(); err != nil {
		return nil, err
	}
	funcname, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateEventTrigStmt{Trigname: trigname, Eventname: eventname, Whenclause: whenclause, Funcname: funcname}, nil
}

func (p *Parser) parseEventTriggerWhenItem() (*nodes.DefElem, error) {
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(IN_P); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// event_trigger_value_list: SCONST [, SCONST ...]
	tok, err := p.expect(SCONST)
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: tok.Str}}
	for p.cur.Type == ',' {
		p.advance()
		tok, err = p.expect(SCONST)
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: tok.Str})
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.DefElem{Defname: name, Arg: &nodes.List{Items: items}, Loc: nodes.NoLoc()}, nil
}

func (p *Parser) parseAlterEventTrigStmt() (*nodes.AlterEventTrigStmt, error) {
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}
	trigname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	tgenabled, err := p.parseEnableTrigger()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterEventTrigStmt{Trigname: trigname, Tgenabled: byte(tgenabled)}, nil
}

// parseEnableTrigger parses the enable_trigger production:
//
//	ENABLE_P                 -> TRIGGER_FIRES_ON_ORIGIN
//	ENABLE_P REPLICA         -> TRIGGER_FIRES_ON_REPLICA
//	ENABLE_P ALWAYS          -> TRIGGER_FIRES_ALWAYS
//	DISABLE_P                -> TRIGGER_DISABLED
//
// Anything else is a syntax error; previously the fallthrough silently
// returned TRIGGER_FIRES_ON_ORIGIN for unknown keywords (shape-III-partial).
func (p *Parser) parseEnableTrigger() (int64, error) {
	switch p.cur.Type {
	case ENABLE_P:
		p.advance()
		switch p.cur.Type {
		case REPLICA:
			p.advance()
			return int64(nodes.TRIGGER_FIRES_ON_REPLICA), nil
		case ALWAYS:
			p.advance()
			return int64(nodes.TRIGGER_FIRES_ALWAYS), nil
		default:
			return int64(nodes.TRIGGER_FIRES_ON_ORIGIN), nil
		}
	case DISABLE_P:
		p.advance()
		return int64(nodes.TRIGGER_DISABLED), nil
	}
	return 0, p.syntaxErrorAtCur()
}
