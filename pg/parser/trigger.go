package parser

import (
	"strconv"
	nodes "github.com/bytebase/omni/pg/ast"
)

func (p *Parser) parseCreateTrigStmt(replace bool) nodes.Node {
	if p.cur.Type == CONSTRAINT {
		return p.parseCreateConstraintTrigger(replace)
	}
	p.expect(TRIGGER)
	trigname, _ := p.parseName()
	timing := p.parseTriggerActionTime()
	events, columns := p.parseTriggerEvents()
	p.expect(ON)
	relNames, _ := p.parseQualifiedName()
	relation := makeRangeVarFromAnyName(relNames)
	transitionRels := p.parseTriggerReferencing()
	row := p.parseTriggerForSpec()
	whenClause := p.parseTriggerWhen()
	p.expect(EXECUTE)
	if p.cur.Type == FUNCTION || p.cur.Type == PROCEDURE { p.advance() }
	funcname, _ := p.parseFuncName()
	p.expect('(')
	args := p.parseTriggerFuncArgs()
	p.expect(')')
	return &nodes.CreateTrigStmt{
		Replace: replace, IsConstraint: false, Trigname: trigname,
		Relation: relation, Funcname: funcname, Args: args, Row: row,
		Timing: int16(timing), Events: int16(events), Columns: columns,
		WhenClause: whenClause, TransitionRels: transitionRels,
		Deferrable: false, Initdeferred: false,
	}
}

func (p *Parser) parseCreateConstraintTrigger(replace bool) nodes.Node {
	p.expect(CONSTRAINT)
	p.expect(TRIGGER)
	trigname, _ := p.parseName()
	p.expect(AFTER)
	events, columns := p.parseTriggerEvents()
	p.expect(ON)
	relNames, _ := p.parseQualifiedName()
	relation := makeRangeVarFromAnyName(relNames)
	var constrrel *nodes.RangeVar
	if p.cur.Type == FROM {
		p.advance()
		fromNames, _ := p.parseQualifiedName()
		constrrel = makeRangeVarFromAnyName(fromNames)
	}
	casBits := p.parseConstraintAttributeSpec()
	deferrable := (casBits & int64(nodes.CAS_DEFERRABLE)) != 0
	initdeferred := (casBits & int64(nodes.CAS_INITIALLY_DEFERRED)) != 0
	p.expect(FOR)
	if p.cur.Type == EACH { p.advance() }
	p.expect(ROW)
	whenClause := p.parseTriggerWhen()
	p.expect(EXECUTE)
	if p.cur.Type == FUNCTION || p.cur.Type == PROCEDURE { p.advance() }
	funcname, _ := p.parseFuncName()
	p.expect('(')
	args := p.parseTriggerFuncArgs()
	p.expect(')')
	return &nodes.CreateTrigStmt{
		Replace: replace, IsConstraint: true, Trigname: trigname,
		Relation: relation, Funcname: funcname, Args: args, Row: true,
		Timing: int16(nodes.TRIGGER_TYPE_AFTER), Events: int16(events),
		Columns: columns, WhenClause: whenClause,
		Deferrable: deferrable, Initdeferred: initdeferred, Constrrel: constrrel,
	}
}

func (p *Parser) parseTriggerActionTime() int64 {
	switch p.cur.Type {
	case BEFORE: p.advance(); return int64(nodes.TRIGGER_TYPE_BEFORE)
	case AFTER: p.advance(); return int64(nodes.TRIGGER_TYPE_AFTER)
	case INSTEAD: p.advance(); p.expect(OF); return int64(nodes.TRIGGER_TYPE_INSTEAD)
	}
	return int64(nodes.TRIGGER_TYPE_AFTER)
}

func (p *Parser) parseTriggerEvents() (int64, *nodes.List) {
	events, cols := p.parseTriggerOneEvent()
	for p.cur.Type == OR {
		p.advance()
		ev2, cols2 := p.parseTriggerOneEvent()
		events |= ev2
		cols = concatColumnLists(cols, cols2)
	}
	return events, cols
}

func (p *Parser) parseTriggerOneEvent() (int64, *nodes.List) {
	switch p.cur.Type {
	case INSERT: p.advance(); return int64(nodes.TRIGGER_TYPE_INSERT), nil
	case DELETE_P: p.advance(); return int64(nodes.TRIGGER_TYPE_DELETE), nil
	case UPDATE:
		p.advance()
		if p.cur.Type == OF { p.advance(); return int64(nodes.TRIGGER_TYPE_UPDATE), p.parseColumnList() }
		return int64(nodes.TRIGGER_TYPE_UPDATE), nil
	case TRUNCATE: p.advance(); return int64(nodes.TRIGGER_TYPE_TRUNCATE), nil
	}
	return 0, nil
}

func concatColumnLists(a, b *nodes.List) *nodes.List {
	if a == nil && b == nil { return nil }
	var items []nodes.Node
	if a != nil { items = append(items, a.Items...) }
	if b != nil { items = append(items, b.Items...) }
	if len(items) == 0 { return nil }
	return &nodes.List{Items: items}
}

func (p *Parser) parseTriggerReferencing() *nodes.List {
	if p.cur.Type != REFERENCING { return nil }
	p.advance()
	var items []nodes.Node
	for {
		tt := p.parseTriggerTransition()
		if tt == nil { break }
		items = append(items, tt)
	}
	if len(items) == 0 { return nil }
	return &nodes.List{Items: items}
}

func (p *Parser) parseTriggerTransition() *nodes.TriggerTransition {
	var isNew bool
	switch p.cur.Type {
	case NEW: isNew = true; p.advance()
	case OLD: isNew = false; p.advance()
	default: return nil
	}
	isTable := false
	if p.cur.Type == TABLE { isTable = true; p.advance() } else if p.cur.Type == ROW { p.advance() }
	if p.cur.Type == AS { p.advance() }
	name, _ := p.parseColId()
	return &nodes.TriggerTransition{Name: name, IsNew: isNew, IsTable: isTable}
}

func (p *Parser) parseTriggerForSpec() bool {
	if p.cur.Type != FOR { return false }
	p.advance()
	if p.cur.Type == EACH { p.advance() }
	if p.cur.Type == ROW { p.advance(); return true }
	if p.cur.Type == STATEMENT { p.advance(); return false }
	return false
}

func (p *Parser) parseTriggerWhen() nodes.Node {
	if p.cur.Type != WHEN { return nil }
	p.advance()
	p.expect('(')
	expr := p.parseAExpr(0)
	p.expect(')')
	return expr
}

func (p *Parser) parseTriggerFuncArgs() *nodes.List {
	if p.cur.Type == ')' { return nil }
	var items []nodes.Node
	items = append(items, p.parseTriggerFuncArg())
	for p.cur.Type == ',' { p.advance(); items = append(items, p.parseTriggerFuncArg()) }
	return &nodes.List{Items: items}
}

func (p *Parser) parseTriggerFuncArg() nodes.Node {
	switch p.cur.Type {
	case ICONST: tok := p.advance(); return &nodes.String{Str: strconv.FormatInt(tok.Ival, 10)}
	case FCONST: tok := p.advance(); return &nodes.String{Str: tok.Str}
	case SCONST: tok := p.advance(); return &nodes.String{Str: tok.Str}
	default: label, _ := p.parseColLabel(); return &nodes.String{Str: label}
	}
}

func (p *Parser) parseCreateEventTrigStmt() nodes.Node {
	p.expect(TRIGGER)
	trigname, _ := p.parseName()
	p.expect(ON)
	eventname, _ := p.parseColLabel()
	var whenclause *nodes.List
	if p.cur.Type == WHEN {
		p.advance()
		var items []nodes.Node
		items = append(items, p.parseEventTriggerWhenItem())
		for p.cur.Type == AND { p.advance(); items = append(items, p.parseEventTriggerWhenItem()) }
		whenclause = &nodes.List{Items: items}
	}
	p.expect(EXECUTE)
	if p.cur.Type == FUNCTION || p.cur.Type == PROCEDURE { p.advance() }
	funcname, _ := p.parseFuncName()
	p.expect('(')
	p.expect(')')
	return &nodes.CreateEventTrigStmt{Trigname: trigname, Eventname: eventname, Whenclause: whenclause, Funcname: funcname}
}

func (p *Parser) parseEventTriggerWhenItem() nodes.Node {
	name, _ := p.parseColId()
	p.expect(IN_P)
	p.expect('(')
	tok := p.advance()
	items := []nodes.Node{&nodes.String{Str: tok.Str}}
	for p.cur.Type == ',' { p.advance(); tok = p.advance(); items = append(items, &nodes.String{Str: tok.Str}) }
	p.expect(')')
	return &nodes.DefElem{Defname: name, Arg: &nodes.List{Items: items}}
}

func (p *Parser) parseAlterEventTrigStmt() nodes.Node {
	p.expect(TRIGGER)
	trigname, _ := p.parseName()
	tgenabled := p.parseEnableTrigger()
	return &nodes.AlterEventTrigStmt{Trigname: trigname, Tgenabled: byte(tgenabled)}
}

func (p *Parser) parseEnableTrigger() int64 {
	switch p.cur.Type {
	case ENABLE_P:
		p.advance()
		switch p.cur.Type {
		case REPLICA: p.advance(); return int64(nodes.TRIGGER_FIRES_ON_REPLICA)
		case ALWAYS: p.advance(); return int64(nodes.TRIGGER_FIRES_ALWAYS)
		default: return int64(nodes.TRIGGER_FIRES_ON_ORIGIN)
		}
	case DISABLE_P: p.advance(); return int64(nodes.TRIGGER_DISABLED)
	}
	return int64(nodes.TRIGGER_FIRES_ON_ORIGIN)
}
