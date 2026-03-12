// Package parser - alter_table.go implements T-SQL ALTER TABLE statement parsing.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseAlterTableStmt parses an ALTER TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-table-transact-sql
//
//	ALTER TABLE [ database_name . [ schema_name ] . | schema_name . ] table_name
//	{
//	    ALTER COLUMN column_name
//	    {
//	        [ type_schema_name. ] type_name
//	            [ ( { precision [ , scale ] | max | xml_schema_collection } ) ]
//	        [ COLLATE collation_name ]
//	        [ NULL | NOT NULL ] [ SPARSE ]
//	      | { ADD | DROP }
//	          { ROWGUIDCOL | PERSISTED | NOT FOR REPLICATION | SPARSE | HIDDEN }
//	      | { ADD | DROP } MASKED [ WITH ( FUNCTION = ' mask_function ') ]
//	    }
//	    [ WITH ( ONLINE = ON | OFF ) ]
//
//	    | [ WITH { CHECK | NOCHECK } ] ADD
//	    {
//	        <column_definition>
//	      | <computed_column_definition>
//	      | <table_constraint>
//	    } [ ,...n ]
//
//	    | DROP
//	     [ {
//	         [ CONSTRAINT ] [ IF EXISTS ]
//	         {
//	              constraint_name
//	         } [ ,...n ]
//	         | COLUMN [ IF EXISTS ]
//	         {
//	              column_name
//	         } [ ,...n ]
//	     } [ ,...n ] ]
//
//	    | [ WITH { CHECK | NOCHECK } ] { CHECK | NOCHECK } CONSTRAINT
//	        { ALL | constraint_name [ ,...n ] }
//
//	    | { ENABLE | DISABLE } TRIGGER
//	        { ALL | trigger_name [ ,...n ] }
//
//	    | { ENABLE | DISABLE } CHANGE_TRACKING
//	        [ WITH ( TRACK_COLUMNS_UPDATED = { ON | OFF } ) ]
//
//	    | SWITCH [ PARTITION source_partition_number_expression ]
//	        TO target_table
//	        [ PARTITION target_partition_number_expression ]
//	        [ WITH ( <low_priority_lock_wait> ) ]
//
//	    | SET
//	        (
//	            [ FILESTREAM_ON = { partition_scheme_name | filegroup | "default" | "NULL" } ]
//	            | SYSTEM_VERSIONING =
//	                  { OFF | ON [ ( HISTORY_TABLE = schema_name.history_table_name
//	                      [, DATA_CONSISTENCY_CHECK = { ON | OFF } ] ) ] }
//	            | LOCK_ESCALATION = { AUTO | TABLE | DISABLE }
//	        )
//
//	    | REBUILD
//	      [ [ PARTITION = ALL ]
//	        [ WITH ( <rebuild_option> [ ,...n ] ) ]
//	      | [ PARTITION = partition_number
//	           [ WITH ( <single_partition_rebuild_option> [ ,...n ] ) ] ]
//	      ]
//	}
func (p *Parser) parseAlterTableStmt() *nodes.AlterTableStmt {
	loc := p.pos()

	stmt := &nodes.AlterTableStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name
	stmt.Name = p.parseTableRef()

	// Parse optional WITH { CHECK | NOCHECK } prefix
	withCheck := ""
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwCHECK || next.Type == kwNOCHECK {
			p.advance() // consume WITH
			if p.cur.Type == kwCHECK {
				withCheck = "CHECK"
			} else {
				withCheck = "NOCHECK"
			}
			p.advance() // consume CHECK/NOCHECK
		}
	}

	// Parse actions
	var actions []nodes.Node
	switch {
	case p.cur.Type == kwADD:
		actions = p.parseAlterTableAdd()
	case p.cur.Type == kwDROP:
		actions = p.parseAlterTableDrop()
	case p.cur.Type == kwALTER:
		action := p.parseAlterTableAlterColumn()
		if action != nil {
			actions = append(actions, action)
		}
	case p.cur.Type == kwCHECK:
		action := p.parseAlterTableCheckConstraint(withCheck, true)
		actions = append(actions, action)
	case p.cur.Type == kwNOCHECK:
		action := p.parseAlterTableCheckConstraint(withCheck, false)
		actions = append(actions, action)
	case p.cur.Type == kwSET:
		action := p.parseAlterTableSet()
		actions = append(actions, action)
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "ENABLE"):
		action := p.parseAlterTableEnableDisable(true)
		if action != nil {
			actions = append(actions, action)
		}
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "DISABLE"):
		action := p.parseAlterTableEnableDisable(false)
		if action != nil {
			actions = append(actions, action)
		}
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "SWITCH"):
		action := p.parseAlterTableSwitch()
		actions = append(actions, action)
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "REBUILD"):
		action := p.parseAlterTableRebuild()
		actions = append(actions, action)
	default:
		// Try legacy single action for backwards compatibility
		action := p.parseAlterTableActionLegacy()
		if action != nil {
			actions = append(actions, action)
		}
	}

	stmt.Actions = &nodes.List{Items: actions}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterTableAdd parses ADD with comma-separated columns/constraints.
//
//	ADD { <column_definition> | <table_constraint> } [ ,...n ]
func (p *Parser) parseAlterTableAdd() []nodes.Node {
	p.advance() // consume ADD
	var actions []nodes.Node

	for {
		loc := p.pos()
		action := &nodes.AlterTableAction{
			Loc: nodes.Loc{Start: loc},
		}

		if p.cur.Type == kwCONSTRAINT {
			action.Type = nodes.ATAddConstraint
			action.Constraint = p.parseTableConstraint()
		} else if p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE ||
			p.cur.Type == kwFOREIGN || p.cur.Type == kwCHECK || p.cur.Type == kwDEFAULT {
			// Unnamed constraint
			action.Type = nodes.ATAddConstraint
			action.Constraint = p.parseTableConstraint()
		} else {
			action.Type = nodes.ATAddColumn
			action.Column = p.parseColumnDef()
		}

		action.Loc.End = p.pos()
		actions = append(actions, action)

		if _, ok := p.match(','); !ok {
			break
		}
	}

	return actions
}

// parseAlterTableDrop parses DROP with comma-separated columns/constraints.
//
//	DROP { [ CONSTRAINT ] [ IF EXISTS ] constraint_name [ ,...n ]
//	     | COLUMN [ IF EXISTS ] column_name [ ,...n ] }
func (p *Parser) parseAlterTableDrop() []nodes.Node {
	p.advance() // consume DROP
	var actions []nodes.Node

	if p.cur.Type == kwCOLUMN {
		p.advance() // consume COLUMN
		// IF EXISTS
		if p.cur.Type == kwIF {
			next := p.peekNext()
			if next.Type == kwEXISTS {
				p.advance() // IF
				p.advance() // EXISTS
			}
		}
		// Parse comma-separated column names
		for {
			loc := p.pos()
			action := &nodes.AlterTableAction{
				Type: nodes.ATDropColumn,
				Loc:  nodes.Loc{Start: loc},
			}
			name, _ := p.parseIdentifier()
			action.ColName = name
			action.Loc.End = p.pos()
			actions = append(actions, action)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	} else if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		// IF EXISTS
		if p.cur.Type == kwIF {
			next := p.peekNext()
			if next.Type == kwEXISTS {
				p.advance() // IF
				p.advance() // EXISTS
			}
		}
		// Parse comma-separated constraint names
		for {
			loc := p.pos()
			action := &nodes.AlterTableAction{
				Type: nodes.ATDropConstraint,
				Loc:  nodes.Loc{Start: loc},
			}
			action.Constraint = &nodes.ConstraintDef{
				Loc: nodes.Loc{Start: p.pos()},
			}
			name, _ := p.parseIdentifier()
			action.Constraint.Name = name
			action.Constraint.Loc.End = p.pos()
			action.Loc.End = p.pos()
			actions = append(actions, action)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	} else {
		// No COLUMN or CONSTRAINT keyword - treat as DROP column (legacy)
		loc := p.pos()
		action := &nodes.AlterTableAction{
			Type: nodes.ATDropColumn,
			Loc:  nodes.Loc{Start: loc},
		}
		name, _ := p.parseIdentifier()
		action.ColName = name
		action.Loc.End = p.pos()
		actions = append(actions, action)
	}

	return actions
}

// parseAlterTableAlterColumn parses ALTER COLUMN with extended options.
//
//	ALTER COLUMN column_name
//	{
//	    [ type_schema_name. ] type_name [ ( precision [ , scale ] | max ) ]
//	    [ COLLATE collation_name ]
//	    [ NULL | NOT NULL ] [ SPARSE ]
//	  | { ADD | DROP } { ROWGUIDCOL | PERSISTED | NOT FOR REPLICATION | SPARSE | HIDDEN }
//	  | { ADD | DROP } MASKED [ WITH ( FUNCTION = 'mask_function' ) ]
//	}
func (p *Parser) parseAlterTableAlterColumn() *nodes.AlterTableAction {
	loc := p.pos()
	p.advance() // consume ALTER

	if p.cur.Type == kwCOLUMN {
		p.advance() // consume COLUMN
	}

	name, _ := p.parseIdentifier()

	// Check if this is ADD/DROP form (alter column attribute)
	if p.cur.Type == kwADD || p.cur.Type == kwDROP {
		return p.parseAlterColumnAddDrop(loc, name)
	}

	// Type change form
	action := &nodes.AlterTableAction{
		Type:    nodes.ATAlterColumn,
		ColName: name,
		Loc:     nodes.Loc{Start: loc},
	}

	action.DataType = p.parseDataType()

	// COLLATE collation_name
	if p.cur.Type == kwCOLLATE {
		p.advance() // consume COLLATE
		collation, _ := p.parseIdentifier()
		action.Collation = collation
	}

	// NULL / NOT NULL
	if p.cur.Type == kwNULL {
		p.advance()
	} else if p.cur.Type == kwNOT {
		next := p.peekNext()
		if next.Type == kwNULL {
			p.advance() // NOT
			p.advance() // NULL
		}
	}

	// SPARSE (optional after NULL/NOT NULL)
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "SPARSE") {
		p.advance()
	}

	action.Loc.End = p.pos()
	return action
}

// parseAlterColumnAddDrop parses ALTER COLUMN col {ADD|DROP} {ROWGUIDCOL|PERSISTED|...}
func (p *Parser) parseAlterColumnAddDrop(loc int, colName string) *nodes.AlterTableAction {
	action := &nodes.AlterTableAction{
		Type:    nodes.ATAlterColumnAddDrop,
		ColName: colName,
		Loc:     nodes.Loc{Start: loc},
	}

	addOrDrop := "ADD"
	if p.cur.Type == kwDROP {
		addOrDrop = "DROP"
	}
	p.advance() // consume ADD/DROP

	var optItems []nodes.Node
	optItems = append(optItems, &nodes.String{Str: addOrDrop})

	// Determine what's being added/dropped
	switch {
	case p.cur.Type == kwROWGUIDCOL:
		optItems = append(optItems, &nodes.String{Str: "ROWGUIDCOL"})
		p.advance()
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "PERSISTED"):
		optItems = append(optItems, &nodes.String{Str: "PERSISTED"})
		p.advance()
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "SPARSE"):
		optItems = append(optItems, &nodes.String{Str: "SPARSE"})
		p.advance()
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "HIDDEN"):
		optItems = append(optItems, &nodes.String{Str: "HIDDEN"})
		p.advance()
	case p.isIdentLike() && strings.EqualFold(p.cur.Str, "MASKED"):
		optItems = append(optItems, &nodes.String{Str: "MASKED"})
		p.advance()
		// WITH ( FUNCTION = 'mask_function' )
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			if p.cur.Type == '(' {
				opts := p.parseKeyValueOptionList()
				if opts != nil {
					optItems = append(optItems, opts.Items...)
				}
			}
		}
	case p.cur.Type == kwNOT:
		// NOT FOR REPLICATION
		p.advance() // NOT
		if p.cur.Type == kwFOR {
			p.advance() // FOR
		}
		if p.cur.Type == kwREPLICATION {
			p.advance() // REPLICATION
		}
		optItems = append(optItems, &nodes.String{Str: "NOT FOR REPLICATION"})
	default:
		// Unknown option, consume as identifier
		if p.isIdentLike() {
			optItems = append(optItems, &nodes.String{Str: p.cur.Str})
			p.advance()
		}
	}

	action.Options = &nodes.List{Items: optItems}
	action.Loc.End = p.pos()
	return action
}

// parseAlterTableCheckConstraint parses CHECK/NOCHECK CONSTRAINT { ALL | name [,...n] }.
func (p *Parser) parseAlterTableCheckConstraint(withCheck string, isCheck bool) *nodes.AlterTableAction {
	loc := p.pos()
	action := &nodes.AlterTableAction{
		Loc:       nodes.Loc{Start: loc},
		WithCheck: withCheck,
	}

	if isCheck {
		action.Type = nodes.ATCheckConstraint
	} else {
		action.Type = nodes.ATNocheckConstraint
	}
	p.advance() // consume CHECK/NOCHECK

	// CONSTRAINT keyword
	if p.cur.Type == kwCONSTRAINT {
		p.advance()
	}

	// ALL or constraint_name list
	var names []nodes.Node
	if p.cur.Type == kwALL {
		names = append(names, &nodes.String{Str: "ALL"})
		p.advance()
	} else {
		for {
			name, _ := p.parseIdentifier()
			names = append(names, &nodes.String{Str: name})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	action.Names = &nodes.List{Items: names}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableEnableDisable parses ENABLE/DISABLE TRIGGER or CHANGE_TRACKING.
func (p *Parser) parseAlterTableEnableDisable(enable bool) *nodes.AlterTableAction {
	loc := p.pos()
	p.advance() // consume ENABLE/DISABLE

	// TRIGGER or CHANGE_TRACKING
	if p.cur.Type == kwTRIGGER {
		return p.parseAlterTableEnableDisableTrigger(loc, enable)
	}
	if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CHANGE_TRACKING") {
		return p.parseAlterTableChangeTracking(loc, enable)
	}

	return nil
}

// parseAlterTableEnableDisableTrigger parses ENABLE/DISABLE TRIGGER { ALL | name [,...n] }.
func (p *Parser) parseAlterTableEnableDisableTrigger(loc int, enable bool) *nodes.AlterTableAction {
	action := &nodes.AlterTableAction{
		Loc: nodes.Loc{Start: loc},
	}

	if enable {
		action.Type = nodes.ATEnableTrigger
	} else {
		action.Type = nodes.ATDisableTrigger
	}
	p.advance() // consume TRIGGER

	// ALL or trigger_name list
	var names []nodes.Node
	if p.cur.Type == kwALL {
		names = append(names, &nodes.String{Str: "ALL"})
		p.advance()
	} else {
		for {
			name, _ := p.parseIdentifier()
			names = append(names, &nodes.String{Str: name})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	action.Names = &nodes.List{Items: names}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableChangeTracking parses ENABLE/DISABLE CHANGE_TRACKING [WITH (...)].
func (p *Parser) parseAlterTableChangeTracking(loc int, enable bool) *nodes.AlterTableAction {
	action := &nodes.AlterTableAction{
		Loc: nodes.Loc{Start: loc},
	}

	if enable {
		action.Type = nodes.ATEnableChangeTracking
	} else {
		action.Type = nodes.ATDisableChangeTracking
	}
	p.advance() // consume CHANGE_TRACKING

	// WITH ( TRACK_COLUMNS_UPDATED = { ON | OFF } )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			action.Options = p.parseKeyValueOptionList()
		}
	}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableSwitch parses SWITCH [PARTITION n] TO target [PARTITION n].
func (p *Parser) parseAlterTableSwitch() *nodes.AlterTableAction {
	loc := p.pos()
	action := &nodes.AlterTableAction{
		Type: nodes.ATSwitchPartition,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume SWITCH

	// Optional PARTITION source_partition_number
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		action.Partition = p.parseExpr()
	}

	// TO target_table
	if p.cur.Type == kwTO {
		p.advance() // consume TO
	}
	action.TargetName = p.parseTableRef()

	// Optional PARTITION target_partition_number
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		action.TargetPart = p.parseExpr()
	}

	// Optional WITH ( <low_priority_lock_wait> )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			action.Options = p.parseOptionList()
		}
	}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableRebuild parses REBUILD [PARTITION = ALL|n] [WITH (...)].
func (p *Parser) parseAlterTableRebuild() *nodes.AlterTableAction {
	loc := p.pos()
	action := &nodes.AlterTableAction{
		Type: nodes.ATRebuild,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume REBUILD

	// Optional PARTITION = { ALL | partition_number }
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.cur.Type == kwALL {
			action.ColName = "ALL" // reuse ColName to indicate PARTITION = ALL
			p.advance()
		} else {
			action.Partition = p.parseExpr()
		}
	}

	// Optional WITH ( rebuild_option [,...n] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			action.Options = p.parseOptionList()
		}
	}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableSet parses SET ( option = value [,...] ).
//
//	SET ( LOCK_ESCALATION = { AUTO | TABLE | DISABLE }
//	    | FILESTREAM_ON = { partition_scheme_name | filegroup | "default" | "NULL" }
//	    | SYSTEM_VERSIONING = { OFF | ON [ ( HISTORY_TABLE = schema_name.table_name [...] ) ] }
//	    )
func (p *Parser) parseAlterTableSet() *nodes.AlterTableAction {
	loc := p.pos()
	action := &nodes.AlterTableAction{
		Type: nodes.ATSet,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume SET

	if p.cur.Type == '(' {
		action.Options = p.parseKeyValueOptionList()
	}

	action.Loc.End = p.pos()
	return action
}

// parseAlterTableActionLegacy is a fallback for unrecognized ALTER TABLE actions.
func (p *Parser) parseAlterTableActionLegacy() *nodes.AlterTableAction {
	return nil
}

// parseKeyValueOptionList parses ( NAME = value [, ...] ) where values can be
// keywords like ON, OFF, TABLE, AUTO, DISABLE etc. Returns a list of String
// nodes in alternating name/value pairs. Nested parentheses are handled for
// sub-options like SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.t).
func (p *Parser) parseKeyValueOptionList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Parse name (identifier-like, but also handle keywords)
		name := p.consumeAnyIdent()
		items = append(items, &nodes.String{Str: name})

		// = sign
		if p.cur.Type == '=' {
			p.advance()
		}

		// Parse value - could be keyword, identifier, number, string, or nested (...)
		if p.cur.Type == '(' {
			// Nested option list, e.g., SYSTEM_VERSIONING = ON (HISTORY_TABLE = ...)
			nested := p.parseKeyValueOptionList()
			items = append(items, nested)
		} else {
			val := p.consumeAnyIdent()
			items = append(items, &nodes.String{Str: val})

			// Check for nested parens after value, e.g., ON (HISTORY_TABLE = ...)
			if p.cur.Type == '(' {
				nested := p.parseKeyValueOptionList()
				items = append(items, nested)
			}
		}

		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}

// consumeAnyIdent consumes the current token as an identifier string,
// regardless of whether it's a keyword token or a regular identifier.
// This is used in option-list contexts where keywords like ON, OFF, TABLE, etc.
// should be treated as identifiers.
func (p *Parser) consumeAnyIdent() string {
	var s string
	if p.cur.Type == tokICONST {
		s = p.cur.Str
		if s == "" {
			s = fmt.Sprintf("%d", p.cur.Ival)
		}
		p.advance()
		return s
	}
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		s = p.cur.Str
		p.advance()
		return s
	}
	if p.cur.Str != "" {
		s = p.cur.Str
	} else {
		// Map keyword token back to string
		switch p.cur.Type {
		case kwON:
			s = "ON"
		case kwOFF:
			s = "OFF"
		case kwTABLE:
			s = "TABLE"
		case kwALL:
			s = "ALL"
		case kwDEFAULT:
			s = "DEFAULT"
		case kwNULL:
			s = "NULL"
		case kwSET:
			s = "SET"
		case kwINDEX:
			s = "INDEX"
		default:
			s = "?"
		}
	}
	p.advance()
	// Handle dotted names (e.g., dbo.history)
	for p.cur.Type == '.' {
		p.advance() // consume .
		part := ""
		if p.cur.Str != "" {
			part = p.cur.Str
		}
		p.advance()
		s = s + "." + part
	}
	return s
}
