package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// ALTER FUNCTION / PROCEDURE / ROUTINE
// ---------------------------------------------------------------------------

// parseAlterFunctionStmt parses ALTER FUNCTION/PROCEDURE/ROUTINE ...
// ALTER has already been consumed. Current token is FUNCTION/PROCEDURE/ROUTINE.
func (p *Parser) parseAlterFunctionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	var objtype nodes.ObjectType
	switch p.cur.Type {
	case FUNCTION:
		objtype = nodes.OBJECT_FUNCTION
	case PROCEDURE:
		objtype = nodes.OBJECT_PROCEDURE
	case ROUTINE:
		objtype = nodes.OBJECT_ROUTINE
	}
	p.advance() // consume FUNCTION/PROCEDURE/ROUTINE

	fwa, err := p.parseFunctionWithArgtypes()
	if err != nil {
		return nil, err
	}

	// Dispatch on the next token to determine which ALTER variant.
	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: objtype,
			Object:     fwa,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: objtype,
			Object:     fwa,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: objtype,
				Object:     fwa,
				Newschema:  newschema,
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		// Otherwise it's alterfunc_opt_list (e.g., SET search_path ...)
		actions := p.parseAlterfuncOptList()
		p.parseOptRestrict()
		return &nodes.AlterFunctionStmt{
			Objtype: objtype,
			Func:    fwa,
			Actions: actions,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case NO, DEPENDS:
		remove := p.parseOptNo()
		if _, err := p.expect(DEPENDS); err != nil {
			return nil, err
		}
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXTENSION); err != nil {
			return nil, err
		}
		extname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectDependsStmt{
			ObjectType: objtype,
			Object:     fwa,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		// alterfunc_opt_list opt_restrict (e.g., IMMUTABLE, STABLE, etc.)
		actions := p.parseAlterfuncOptList()
		p.parseOptRestrict()
		return &nodes.AlterFunctionStmt{
			Objtype: objtype,
			Func:    fwa,
			Actions: actions,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
}

// parseAlterfuncOptList parses alterfunc_opt_list: one or more common_func_opt_item.
func (p *Parser) parseAlterfuncOptList() *nodes.List {
	item := p.parseCommonFuncOptItem()
	if item == nil {
		return nil
	}
	items := []nodes.Node{item}
	for {
		item = p.parseCommonFuncOptItem()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}
}

// parseOptRestrict consumes an optional RESTRICT keyword (ignored, for SQL compliance).
func (p *Parser) parseOptRestrict() {
	if p.cur.Type == RESTRICT {
		p.advance()
	}
}

// parseOptNo parses opt_no: NO | /* EMPTY */. Returns true if NO was consumed.
func (p *Parser) parseOptNo() bool {
	if p.cur.Type == NO {
		p.advance()
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// ALTER TYPE
// ---------------------------------------------------------------------------

// parseAlterTypeStmt parses ALTER TYPE ...
// ALTER has already been consumed. Current token is TYPE_P.
func (p *Parser) parseAlterTypeStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume TYPE
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case ADD_P:
		// ALTER TYPE name ADD VALUE ...
		next := p.peekNext()
		if next.Type == VALUE_P {
			stmt, err := p.parseAlterEnumAddValue(names)
			if stmt != nil {
				stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
			}
			return stmt, err
		}
		// ALTER TYPE name ADD ATTRIBUTE ... (AlterCompositeTypeStmt)
		return p.parseAlterCompositeType(names)
	case DROP:
		// ALTER TYPE name DROP ATTRIBUTE ... (AlterCompositeTypeStmt)
		return p.parseAlterCompositeType(names)
	case ALTER:
		// ALTER TYPE name ALTER ATTRIBUTE ... (AlterCompositeTypeStmt)
		return p.parseAlterCompositeType(names)
	case RENAME:
		p.advance() // consume RENAME
		switch p.cur.Type {
		case TO:
			p.advance()
			newname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TYPE,
				Object:     names,
				Newname:    newname,
			}, nil
		case VALUE_P:
			// ALTER TYPE name RENAME VALUE 'old' TO 'new'
			p.advance() // consume VALUE
			oldval := p.cur.Str
			p.advance() // consume Sconst
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			newval := p.cur.Str
			p.advance() // consume Sconst
			return &nodes.AlterEnumStmt{
				Typname: names,
				Oldval:  oldval,
				Newval:  newval,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		case ATTRIBUTE:
			// ALTER TYPE name RENAME ATTRIBUTE name TO name opt_drop_behavior
			return p.parseAlterCompositeTypeRename(names)
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TYPE,
			Object:     names,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TYPE,
				Object:     names,
				Newschema:  newschema,
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		// ALTER TYPE name SET (operator_def_list) -> AlterTypeStmt
		p.advance() // consume SET
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		opts, err := p.parseOperatorDefList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.AlterTypeStmt{
			TypeName: names,
			Options:  opts,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// parseAlterEnumAddValue parses ALTER TYPE name ADD VALUE ...
// Current token is ADD.
func (p *Parser) parseAlterEnumAddValue(typname *nodes.List) (*nodes.AlterEnumStmt, error) {
	p.advance() // consume ADD
	p.advance() // consume VALUE

	skipIfExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		skipIfExists = true
	}

	newval := p.cur.Str
	p.advance() // consume Sconst

	stmt := &nodes.AlterEnumStmt{
		Typname:            typname,
		Newval:             newval,
		SkipIfNewvalExists: skipIfExists,
	}

	if p.cur.Type == BEFORE {
		p.advance()
		stmt.NewvalNeighbor = p.cur.Str
		p.advance()
		stmt.NewvalIsAfter = false
	} else if p.cur.Type == AFTER {
		p.advance()
		stmt.NewvalNeighbor = p.cur.Str
		p.advance()
		stmt.NewvalIsAfter = true
	}

	return stmt, nil
}

// parseAlterCompositeType parses AlterCompositeTypeStmt (ALTER TYPE name alter_type_cmds).
// Current token is ADD/DROP/ALTER.
func (p *Parser) parseAlterCompositeType(names *nodes.List) (*nodes.AlterTableStmt, error) {
	cmds, err := p.parseAlterTypeCmds()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromCompositeType(names)
	return &nodes.AlterTableStmt{
		Relation: rv,
		Cmds:     cmds,
		ObjType:  int(nodes.OBJECT_TYPE),
	}, nil
}

// parseAlterCompositeTypeRename parses ALTER TYPE name RENAME ATTRIBUTE name TO name opt_drop_behavior.
// RENAME has been consumed. Current token is ATTRIBUTE.
func (p *Parser) parseAlterCompositeTypeRename(names *nodes.List) (*nodes.RenameStmt, error) {
	p.advance() // consume ATTRIBUTE
	subname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	newname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	return &nodes.RenameStmt{
		RenameType:   nodes.OBJECT_ATTRIBUTE,
		RelationType: nodes.OBJECT_TYPE,
		Relation:     makeRangeVarFromAnyName(names),
		Subname:      subname,
		Newname:      newname,
		Behavior:     nodes.DropBehavior(behavior),
	}, nil
}

// makeRangeVarFromCompositeType creates a RangeVar from composite type names with Inh=true.
func makeRangeVarFromCompositeType(names *nodes.List) *nodes.RangeVar {
	rv := &nodes.RangeVar{
		Inh:      true,
		Loc: nodes.NoLoc(),
	}
	if names != nil && len(names.Items) > 0 {
		switch len(names.Items) {
		case 1:
			rv.Relname = names.Items[0].(*nodes.String).Str
		case 2:
			rv.Schemaname = names.Items[0].(*nodes.String).Str
			rv.Relname = names.Items[1].(*nodes.String).Str
		case 3:
			rv.Catalogname = names.Items[0].(*nodes.String).Str
			rv.Schemaname = names.Items[1].(*nodes.String).Str
			rv.Relname = names.Items[2].(*nodes.String).Str
		}
	}
	return rv
}

// parseAlterTypeCmds parses alter_type_cmds: alter_type_cmd (',' alter_type_cmd)*
func (p *Parser) parseAlterTypeCmds() (*nodes.List, error) {
	cmd, err := p.parseAlterTypeCmd()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{cmd}
	for p.cur.Type == ',' {
		p.advance()
		cmd, err := p.parseAlterTypeCmd()
		if err != nil {
			return nil, err
		}
		items = append(items, cmd)
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTypeCmd parses alter_type_cmd.
func (p *Parser) parseAlterTypeCmd() (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case ADD_P:
		p.advance()
		if _, err := p.expect(ATTRIBUTE); err != nil {
			return nil, err
		}
		elem, err := p.parseTableFuncElement()
		if err != nil {
			return nil, err
		}
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:  int(nodes.AT_AddColumn),
			Def:      elem,
			Behavior: behavior,
		}, nil
	case DROP:
		p.advance()
		if _, err := p.expect(ATTRIBUTE); err != nil {
			return nil, err
		}
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			if _, err := p.expect(EXISTS); err != nil {
				return nil, err
			}
			missingOk = true
		}
		colname, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropColumn),
			Name:       colname,
			Behavior:   behavior,
			Missing_ok: missingOk,
		}, nil
	case ALTER:
		p.advance()
		if _, err := p.expect(ATTRIBUTE); err != nil {
			return nil, err
		}
		colname, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		// SET DATA TYPE or just TYPE
		if p.cur.Type == SET {
			p.advance()
			if _, err := p.expect(DATA_P); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(TYPE_P); err != nil {
			return nil, err
		}
		typename, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		var collClause *nodes.CollateClause
		if p.cur.Type == COLLATE {
			p.advance()
			collname, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			collClause = &nodes.CollateClause{Collname: collname, Loc: nodes.NoLoc()}
		}
		behavior := p.parseOptDropBehavior()
		coldef := &nodes.ColumnDef{
			Colname:    colname,
			TypeName:   typename,
			CollClause: collClause,
		}
		return &nodes.AlterTableCmd{
			Subtype:  int(nodes.AT_AlterColumnType),
			Name:     colname,
			Def:      coldef,
			Behavior: behavior,
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER DOMAIN
// ---------------------------------------------------------------------------

// parseAlterDomainOwnerOrOther parses ALTER DOMAIN ...
// ALTER has already been consumed. Current token is DOMAIN_P.
func (p *Parser) parseAlterDomainOwnerOrOther() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume DOMAIN

	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_DOMAIN,
				Object:     names,
				Newschema:  newschema,
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if next.Type == NOT {
			// SET NOT NULL
			p.advance() // consume SET
			p.advance() // consume NOT
			if _, err := p.expect(NULL_P); err != nil {
				return nil, err
			}
			return &nodes.AlterDomainStmt{
				Subtype: 'O',
				Typname: names,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if next.Type == DEFAULT {
			// SET DEFAULT a_expr
			p.advance() // consume SET
			p.advance() // consume DEFAULT
			expr, err := p.parseAExpr(0)
			if err != nil {
				return nil, err
			}
			return &nodes.AlterDomainStmt{
				Subtype: 'T',
				Typname: names,
				Def:     expr,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		return nil, p.syntaxErrorAtCur()
	case DROP:
		p.advance() // consume DROP
		if p.cur.Type == NOT {
			// DROP NOT NULL
			p.advance()
			if _, err := p.expect(NULL_P); err != nil {
				return nil, err
			}
			return &nodes.AlterDomainStmt{
				Subtype: 'N',
				Typname: names,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if p.cur.Type == DEFAULT {
			// DROP DEFAULT
			p.advance()
			return &nodes.AlterDomainStmt{
				Subtype: 'T',
				Typname: names,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if p.cur.Type == CONSTRAINT {
			// DROP CONSTRAINT [IF EXISTS] name opt_drop_behavior
			p.advance()
			missingOk := false
			if p.cur.Type == IF_P {
				p.advance()
				if _, err := p.expect(EXISTS); err != nil {
					return nil, err
				}
				missingOk = true
			}
			cname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			behavior := p.parseOptDropBehavior()
			return &nodes.AlterDomainStmt{
				Subtype:   'X',
				Typname:   names,
				Name:      cname,
				Behavior:  nodes.DropBehavior(behavior),
				MissingOk: missingOk,
				Loc:       nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		return nil, p.syntaxErrorAtCur()
	case ADD_P:
		// ADD [CONSTRAINT name] CHECK (expr)
		p.advance() // consume ADD
		result, err := p.parseDomainConstraintForAlter(names)
		if err != nil {
			return nil, err
		}
		if result != nil {
			result.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		}
		return result, nil
	case VALIDATE:
		p.advance() // consume VALIDATE
		if _, err := p.expect(CONSTRAINT); err != nil {
			return nil, err
		}
		cname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDomainStmt{
			Subtype: 'V',
			Typname: names,
			Name:    cname,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_DOMAIN,
			Object:     names,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if p.cur.Type == CONSTRAINT {
			p.advance()
			subname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			newname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.RenameStmt{
				RenameType: nodes.OBJECT_DOMCONSTRAINT,
				Object:     names,
				Subname:    subname,
				Newname:    newname,
			}, nil
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_DOMAIN,
			Object:     names,
			Newname:    newname,
		}, nil
	default:
		return nil, nil
	}
}

// parseDomainConstraintForAlter parses DomainConstraint for ALTER DOMAIN ... ADD.
func (p *Parser) parseDomainConstraintForAlter(typname *nodes.List) (*nodes.AlterDomainStmt, error) {
	var conname string
	if p.cur.Type == CONSTRAINT {
		p.advance()
		var err error
		conname, err = p.parseName()
		if err != nil {
			return nil, err
		}
	}
	// DomainConstraintElem: CHECK '(' a_expr ')' ConstraintAttributeSpec
	if _, err := p.expect(CHECK); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	// ConstraintAttributeSpec - parse optional attributes
	constraint := &nodes.Constraint{
		Contype:  nodes.CONSTR_CHECK,
		Conname:  conname,
		RawExpr:  expr,
		Loc: nodes.NoLoc(),
	}
	attrs := p.parseConstraintAttributeSpec()
	applyConstraintAttrs(constraint, attrs)
	return &nodes.AlterDomainStmt{
		Subtype: 'C',
		Typname: typname,
		Def:     constraint,
	}, nil
}

// ---------------------------------------------------------------------------
// ALTER SCHEMA
// ---------------------------------------------------------------------------

// parseAlterSchemaOwner parses ALTER SCHEMA ...
// ALTER has already been consumed. Current token is SCHEMA.
func (p *Parser) parseAlterSchemaOwner() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume SCHEMA
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_SCHEMA,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_SCHEMA,
			Subname:    name,
			Newname:    newname,
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER COLLATION
// ---------------------------------------------------------------------------

// parseAlterCollationStmt parses ALTER COLLATION ...
// ALTER has already been consumed. Current token is COLLATION.
func (p *Parser) parseAlterCollationStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume COLLATION
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case REFRESH:
		p.advance()
		if _, err := p.expect(VERSION_P); err != nil {
			return nil, err
		}
		return &nodes.AlterCollationStmt{Collname: names, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER CONVERSION
// ---------------------------------------------------------------------------

// parseAlterConversionStmt parses ALTER CONVERSION ...
// ALTER has already been consumed. Current token is CONVERSION_P.
func (p *Parser) parseAlterConversionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume CONVERSION
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER AGGREGATE
// ---------------------------------------------------------------------------

// parseAlterAggregateStmt parses ALTER AGGREGATE ...
// ALTER has already been consumed. Current token is AGGREGATE.
func (p *Parser) parseAlterAggregateStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume AGGREGATE
	agg, err := p.parseAggregateWithArgtypesLocal()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// parseAggregateWithArgtypesLocal parses aggregate_with_argtypes: func_name aggr_args
func (p *Parser) parseAggregateWithArgtypesLocal() (*nodes.ObjectWithArgs, error) {
	funcname, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	aggrArgs, err := p.parseAggrArgs()
	if err != nil {
		return nil, err
	}
	return &nodes.ObjectWithArgs{
		Objname: funcname,
		Objargs: extractAggrArgTypesLocal(aggrArgs),
	}, nil
}

// extractAggrArgTypesLocal extracts types from aggregate args.
func extractAggrArgTypesLocal(args *nodes.List) *nodes.List {
	if args == nil || len(args.Items) < 1 {
		return nil
	}
	argsList, ok := args.Items[0].(*nodes.List)
	if !ok || argsList == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range argsList.Items {
		if fp, ok := item.(*nodes.FunctionParameter); ok {
			result.Items = append(result.Items, fp.ArgType)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// ALTER TEXT SEARCH
// ---------------------------------------------------------------------------

// parseAlterTextSearchStmt parses ALTER TEXT SEARCH ...
// ALTER has already been consumed. Current token is TEXT_P.
func (p *Parser) parseAlterTextSearchStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume TEXT
	if _, err := p.expect(SEARCH); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case DICTIONARY:
		return p.parseAlterTSDictionary(loc)
	case CONFIGURATION:
		return p.parseAlterTSConfiguration(loc)
	case PARSER:
		return p.parseAlterTSParserOrTemplate(nodes.OBJECT_TSPARSER, loc)
	case TEMPLATE:
		return p.parseAlterTSParserOrTemplate(nodes.OBJECT_TSTEMPLATE, loc)
	default:
		return nil, nil
	}
}

// parseAlterTSDictionary parses ALTER TEXT SEARCH DICTIONARY ...
func (p *Parser) parseAlterTSDictionary(loc int) (nodes.Node, error) {
	p.advance() // consume DICTIONARY
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case '(':
		// ALTER TEXT SEARCH DICTIONARY name definition
		def, err := p.parseDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTSDictionaryStmt{
			Dictname: names,
			Options:  def,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// parseAlterTSConfiguration parses ALTER TEXT SEARCH CONFIGURATION ...
func (p *Parser) parseAlterTSConfiguration(loc int) (nodes.Node, error) {
	p.advance() // consume CONFIGURATION
	cfgname, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case ADD_P:
		// ADD MAPPING FOR name_list WITH any_name_list
		p.advance()
		if _, err := p.expect(MAPPING); err != nil {
			return nil, err
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		tokentype, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(WITH); err != nil {
			return nil, err
		}
		dicts, err := p.parseAnyNameList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTSConfigurationStmt{
			Kind:      nodes.ALTER_TSCONFIG_ADD_MAPPING,
			Cfgname:   cfgname,
			Tokentype: tokentype,
			Dicts:     dicts,
			Loc:       nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case ALTER:
		// ALTER MAPPING ...
		p.advance()
		if _, err := p.expect(MAPPING); err != nil {
			return nil, err
		}
		if p.cur.Type == FOR {
			// ALTER MAPPING FOR name_list WITH/REPLACE ...
			p.advance()
			tokentype, err := p.parseNameList()
			if err != nil {
				return nil, err
			}
			if p.cur.Type == WITH {
				p.advance()
				dicts, err := p.parseAnyNameList()
				if err != nil {
					return nil, err
				}
				return &nodes.AlterTSConfigurationStmt{
					Kind:      nodes.ALTER_TSCONFIG_ALTER_MAPPING_FOR_TOKEN,
					Cfgname:   cfgname,
					Tokentype: tokentype,
					Dicts:     dicts,
					Override:  true,
					Loc:       nodes.Loc{Start: loc, End: p.prev.End},
				}, nil
			}
			// ALTER MAPPING FOR name_list REPLACE any_name WITH any_name
			if _, err := p.expect(REPLACE); err != nil {
				return nil, err
			}
			oldDict, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(WITH); err != nil {
				return nil, err
			}
			newDict, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_REPLACE_DICT_FOR_TOKEN,
				Cfgname:   cfgname,
				Tokentype: tokentype,
				Dicts:     &nodes.List{Items: []nodes.Node{oldDict, newDict}},
				Replace:   true,
				Loc:       nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		// ALTER MAPPING REPLACE any_name WITH any_name
		if _, err := p.expect(REPLACE); err != nil {
			return nil, err
		}
		oldDict, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(WITH); err != nil {
			return nil, err
		}
		newDict, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTSConfigurationStmt{
			Kind:    nodes.ALTER_TSCONFIG_REPLACE_DICT,
			Cfgname: cfgname,
			Dicts:   &nodes.List{Items: []nodes.Node{oldDict, newDict}},
			Replace: true,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case DROP:
		// DROP MAPPING [IF EXISTS] FOR name_list
		p.advance()
		if _, err := p.expect(MAPPING); err != nil {
			return nil, err
		}
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			if _, err := p.expect(EXISTS); err != nil {
				return nil, err
			}
			missingOk = true
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		tokentype, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTSConfigurationStmt{
			Kind:      nodes.ALTER_TSCONFIG_DROP_MAPPING,
			Cfgname:   cfgname,
			Tokentype: tokentype,
			MissingOk: missingOk,
			Loc:       nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// parseAlterTSParserOrTemplate parses ALTER TEXT SEARCH PARSER/TEMPLATE ...
func (p *Parser) parseAlterTSParserOrTemplate(objtype nodes.ObjectType, loc int) (nodes.Node, error) {
	p.advance() // consume PARSER or TEMPLATE
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: objtype,
			Object:     names,
			Newname:    newname,
		}, nil
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: objtype,
			Object:     names,
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER LANGUAGE
// ---------------------------------------------------------------------------

// parseAlterLanguageStmt parses ALTER [PROCEDURAL] LANGUAGE ...
// ALTER has already been consumed. Current token is LANGUAGE or PROCEDURAL.
func (p *Parser) parseAlterLanguageStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	if p.cur.Type == PROCEDURAL {
		p.advance() // consume PROCEDURAL
	}
	p.advance() // consume LANGUAGE
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_LANGUAGE,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_LANGUAGE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER LARGE OBJECT
// ---------------------------------------------------------------------------

// parseAlterLargeObjectStmt parses ALTER LARGE OBJECT ...
// ALTER has already been consumed. Current token is LARGE_P.
func (p *Parser) parseAlterLargeObjectStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume LARGE
	if _, err := p.expect(OBJECT_P); err != nil {
		return nil, err
	}
	numVal := p.parseNumericOnly()
	if _, err := p.expect(OWNER); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	roleSpec := p.parseRoleSpec()
	return &nodes.AlterOwnerStmt{
		ObjectType: nodes.OBJECT_LARGEOBJECT,
		Object:     numVal,
		Newowner:   roleSpec,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// ---------------------------------------------------------------------------
// ALTER EVENT TRIGGER
// ---------------------------------------------------------------------------

// parseAlterEventTriggerOwner parses ALTER EVENT TRIGGER ...
// ALTER has already been consumed. Current token is EVENT.
func (p *Parser) parseAlterEventTriggerOwner() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume EVENT
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_EVENT_TRIGGER,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_EVENT_TRIGGER,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	case ENABLE_P:
		p.advance()
		tgenabled := byte(nodes.TRIGGER_FIRES_ON_ORIGIN)
		if p.cur.Type == REPLICA {
			p.advance()
			tgenabled = nodes.TRIGGER_FIRES_ON_REPLICA
		} else if p.cur.Type == ALWAYS {
			p.advance()
			tgenabled = nodes.TRIGGER_FIRES_ALWAYS
		}
		return &nodes.AlterEventTrigStmt{
			Trigname:  name,
			Tgenabled: tgenabled,
			Loc:       nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case DISABLE_P:
		p.advance()
		return &nodes.AlterEventTrigStmt{
			Trigname:  name,
			Tgenabled: 'D',
			Loc:       nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER TABLESPACE
// ---------------------------------------------------------------------------

// parseAlterTablespaceOwner parses ALTER TABLESPACE ...
// ALTER has already been consumed. Current token is TABLESPACE.
func (p *Parser) parseAlterTablespaceOwner() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume TABLESPACE
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TABLESPACE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TABLESPACE,
			Subname:    name,
			Newname:    newname,
		}, nil
	case SET:
		p.advance()
		opts := p.parseReloptions()
		return &nodes.AlterTableSpaceOptionsStmt{
			Tablespacename: name,
			Options:        opts,
			IsReset:        false,
			Loc:            nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case RESET:
		p.advance()
		opts := p.parseReloptions()
		return &nodes.AlterTableSpaceOptionsStmt{
			Tablespacename: name,
			Options:        opts,
			IsReset:        true,
			Loc:            nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// ALTER TRIGGER ... DEPENDS ON EXTENSION
// ---------------------------------------------------------------------------

// parseAlterTriggerDependsOnExtension parses ALTER TRIGGER ...
// ALTER has already been consumed. Current token is TRIGGER.
//
// Ref: https://www.postgresql.org/docs/17/sql-altertrigger.html
//
//	ALTER TRIGGER name ON table_name RENAME TO new_name
//	ALTER TRIGGER name ON table_name [ NO ] DEPENDS ON EXTENSION extension_name
func (p *Parser) parseAlterTriggerDependsOnExtension() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance() // consume TRIGGER
	trigname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	// qualified_name
	relname, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	rel := makeRangeVarFromAnyName(relname)

	// Dispatch: RENAME TO vs [NO] DEPENDS ON EXTENSION
	if p.cur.Type == RENAME {
		p.advance() // consume RENAME
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TRIGGER,
			Relation:   rel,
			Subname:    trigname,
			Newname:    newname,
		}, nil
	}

	remove := p.parseOptNo()
	if _, err := p.expect(DEPENDS); err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	if _, err := p.expect(EXTENSION); err != nil {
		return nil, err
	}
	extname, err := p.parseName()
	if err != nil {
		return nil, err
	}

	return &nodes.AlterObjectDependsStmt{
		ObjectType: nodes.OBJECT_TRIGGER,
		Relation:   rel,
		Object:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: trigname}}},
		Extname:    &nodes.String{Str: extname},
		Remove:     remove,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// ---------------------------------------------------------------------------
// ALTER RULE
// ---------------------------------------------------------------------------

// parseAlterRuleStmt parses ALTER RULE name ON qualified_name RENAME TO name.
// ALTER has already been consumed. Current token is RULE.
//
// Ref: https://www.postgresql.org/docs/17/sql-alterrule.html
//
//	ALTER RULE name ON table_name RENAME TO new_name
func (p *Parser) parseAlterRuleStmt() (nodes.Node, error) {
	p.advance() // consume RULE
	subname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rel := makeRangeVarFromNames(names)
	if _, err := p.expect(RENAME); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	newname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.RenameStmt{
		RenameType: nodes.OBJECT_RULE,
		Relation:   rel,
		Subname:    subname,
		Newname:    newname,
	}, nil
}
