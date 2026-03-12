package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseSignalStmt parses a SIGNAL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/signal.html
//
//	SIGNAL condition_value
//	    [SET signal_information_item
//	    [, signal_information_item] ...]
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
func (p *Parser) parseSignalStmt() (*nodes.SignalStmt, error) {
	start := p.pos()
	p.advance() // consume SIGNAL

	stmt := &nodes.SignalStmt{Loc: nodes.Loc{Start: start}}

	// Parse condition_value (required for SIGNAL)
	condVal, err := p.parseSignalConditionValue()
	if err != nil {
		return nil, err
	}
	stmt.ConditionValue = condVal

	// Optional SET clause
	if p.cur.Type == kwSET {
		items, err := p.parseSignalSetClause()
		if err != nil {
			return nil, err
		}
		stmt.SetItems = items
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseResignalStmt parses a RESIGNAL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/resignal.html
//
//	RESIGNAL [condition_value]
//	    [SET signal_information_item
//	    [, signal_information_item] ...]
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
func (p *Parser) parseResignalStmt() (*nodes.ResignalStmt, error) {
	start := p.pos()
	p.advance() // consume RESIGNAL

	stmt := &nodes.ResignalStmt{Loc: nodes.Loc{Start: start}}

	// Optional condition_value
	// condition_value starts with SQLSTATE or an identifier (condition_name)
	// but NOT SET (which starts the SET clause)
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "SQLSTATE") {
		condVal, err := p.parseSignalConditionValue()
		if err != nil {
			return nil, err
		}
		stmt.ConditionValue = condVal
	} else if p.cur.Type == tokIDENT && !eqFold(p.cur.Str, "SET") {
		// condition_name — a plain identifier that is not SET
		condVal, err := p.parseSignalConditionValue()
		if err != nil {
			return nil, err
		}
		stmt.ConditionValue = condVal
	}

	// Optional SET clause
	if p.cur.Type == kwSET {
		items, err := p.parseSignalSetClause()
		if err != nil {
			return nil, err
		}
		stmt.SetItems = items
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseGetDiagnosticsStmt parses a GET DIAGNOSTICS statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/get-diagnostics.html
//
//	GET [CURRENT | STACKED] DIAGNOSTICS {
//	    statement_information_item
//	    [, statement_information_item] ...
//	  | CONDITION condition_number
//	    condition_information_item
//	    [, condition_information_item] ...
//	}
//
//	statement_information_item:
//	    target = statement_information_item_name
//
//	statement_information_item_name: {
//	    NUMBER
//	  | ROW_COUNT
//	}
//
//	condition_information_item:
//	    target = condition_information_item_name
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | RETURNED_SQLSTATE
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
//
//	target:
//	    user_variable | local_variable | param
func (p *Parser) parseGetDiagnosticsStmt() (*nodes.GetDiagnosticsStmt, error) {
	start := p.pos()
	p.advance() // consume GET

	stmt := &nodes.GetDiagnosticsStmt{Loc: nodes.Loc{Start: start}}

	// Optional CURRENT | STACKED
	if p.cur.Type == kwCURRENT {
		p.advance()
	} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "STACKED") {
		p.advance()
		stmt.Stacked = true
	}

	// DIAGNOSTICS
	if _, err := p.expect(kwDIAGNOSTICS); err != nil {
		return nil, err
	}

	// CONDITION or statement_information_item
	if p.cur.Type == kwCONDITION {
		p.advance() // consume CONDITION
		stmt.StatementInfo = false

		// condition_number
		condNum, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.ConditionNumber = condNum

		// condition_information_item [, condition_information_item] ...
		items, err := p.parseDiagnosticsItems()
		if err != nil {
			return nil, err
		}
		stmt.Items = items
	} else {
		// statement_information_item [, statement_information_item] ...
		stmt.StatementInfo = true
		items, err := p.parseDiagnosticsItems()
		if err != nil {
			return nil, err
		}
		stmt.Items = items
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSignalConditionValue parses a signal condition value.
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
func (p *Parser) parseSignalConditionValue() (string, error) {
	// SQLSTATE [VALUE] 'sqlstate_value'
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "SQLSTATE") {
		p.advance() // consume SQLSTATE

		// Optional VALUE keyword
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "VALUE") {
			p.advance()
		}

		// sqlstate_value (string literal)
		if p.cur.Type != tokSCONST {
			return "", &ParseError{
				Message:  "expected SQLSTATE value string",
				Position: p.cur.Loc,
			}
		}
		val := p.cur.Str
		p.advance()
		return val, nil
	}

	// condition_name (identifier)
	name, _, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	return name, nil
}

// parseSignalSetClause parses SET signal_information_item [, signal_information_item] ...
func (p *Parser) parseSignalSetClause() ([]*nodes.SignalInfoItem, error) {
	p.advance() // consume SET

	var items []*nodes.SignalInfoItem
	for {
		item, err := p.parseSignalInfoItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseSignalInfoItem parses a single signal information item.
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
func (p *Parser) parseSignalInfoItem() (*nodes.SignalInfoItem, error) {
	start := p.pos()

	// condition_information_item_name is a keyword-like identifier
	// (MESSAGE_TEXT, MYSQL_ERRNO, CLASS_ORIGIN, etc.)
	var name string
	if p.cur.Type == tokIDENT {
		name = p.cur.Str
		p.advance()
	} else {
		// Some of these may be recognized as keywords
		name = p.cur.Str
		p.advance()
	}

	// '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// simple_value_specification (expression)
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &nodes.SignalInfoItem{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Name:  name,
		Value: val,
	}, nil
}

// parseDiagnosticsItems parses target = item_name [, target = item_name] ...
func (p *Parser) parseDiagnosticsItems() ([]*nodes.DiagnosticsItem, error) {
	var items []*nodes.DiagnosticsItem
	for {
		item, err := p.parseDiagnosticsItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseDiagnosticsItem parses a single diagnostics item (target = item_name).
func (p *Parser) parseDiagnosticsItem() (*nodes.DiagnosticsItem, error) {
	start := p.pos()

	// target must be a variable reference (@var or @@var)
	if !p.isVariableRef() {
		return nil, &ParseError{
			Message:  "expected variable reference as diagnostics target",
			Position: p.cur.Loc,
		}
	}
	target, err := p.parseVariableRef()
	if err != nil {
		return nil, err
	}

	// '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// item_name (identifier-like keyword: NUMBER, ROW_COUNT, MESSAGE_TEXT, etc.)
	// These are typically plain identifiers; some may overlap with keywords.
	var name string
	if p.cur.Type == tokIDENT {
		name = p.cur.Str
		p.advance()
	} else {
		// Accept any keyword token as name (e.g. NUMBER might not be a keyword)
		name = p.cur.Str
		p.advance()
	}

	return &nodes.DiagnosticsItem{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Target: target,
		Name:   name,
	}, nil
}
