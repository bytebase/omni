package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// ALTER FUNCTION / PROCEDURE / ROUTINE
// ---------------------------------------------------------------------------

// parseAlterFunctionStmt parses ALTER FUNCTION/PROCEDURE/ROUTINE ...
// ALTER has already been consumed. Current token is FUNCTION/PROCEDURE/ROUTINE.
func (p *Parser) parseAlterFunctionStmt() nodes.Node {
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

	fwa := p.parseFunctionWithArgtypes()

	// Dispatch on the next token to determine which ALTER variant.
	switch p.cur.Type {
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: objtype,
			Object:     fwa,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: objtype,
			Object:     fwa,
			Newowner:   roleSpec,
		}
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, _ := p.parseName()
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: objtype,
				Object:     fwa,
				Newschema:  newschema,
			}
		}
		// Otherwise it's alterfunc_opt_list (e.g., SET search_path ...)
		actions := p.parseAlterfuncOptList()
		p.parseOptRestrict()
		return &nodes.AlterFunctionStmt{
			Objtype: objtype,
			Func:    fwa,
			Actions: actions,
		}
	case NO, DEPENDS:
		remove := p.parseOptNo()
		p.expect(DEPENDS)
		p.expect(ON)
		p.expect(EXTENSION)
		extname, _ := p.parseName()
		return &nodes.AlterObjectDependsStmt{
			ObjectType: objtype,
			Object:     fwa,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
		}
	default:
		// alterfunc_opt_list opt_restrict (e.g., IMMUTABLE, STABLE, etc.)
		actions := p.parseAlterfuncOptList()
		p.parseOptRestrict()
		return &nodes.AlterFunctionStmt{
			Objtype: objtype,
			Func:    fwa,
			Actions: actions,
		}
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
func (p *Parser) parseAlterTypeStmt() nodes.Node {
	p.advance() // consume TYPE
	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case ADD_P:
		// ALTER TYPE name ADD VALUE ...
		next := p.peekNext()
		if next.Type == VALUE_P {
			return p.parseAlterEnumAddValue(names)
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
			newname, _ := p.parseName()
			return &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TYPE,
				Object:     names,
				Newname:    newname,
			}
		case VALUE_P:
			// ALTER TYPE name RENAME VALUE 'old' TO 'new'
			p.advance() // consume VALUE
			oldval := p.cur.Str
			p.advance() // consume Sconst
			p.expect(TO)
			newval := p.cur.Str
			p.advance() // consume Sconst
			return &nodes.AlterEnumStmt{
				Typname: names,
				Oldval:  oldval,
				Newval:  newval,
			}
		case ATTRIBUTE:
			// ALTER TYPE name RENAME ATTRIBUTE name TO name opt_drop_behavior
			return p.parseAlterCompositeTypeRename(names)
		default:
			return nil
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TYPE,
			Object:     names,
			Newowner:   roleSpec,
		}
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, _ := p.parseName()
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TYPE,
				Object:     names,
				Newschema:  newschema,
			}
		}
		// ALTER TYPE name SET (operator_def_list) -> AlterTypeStmt
		p.advance() // consume SET
		p.expect('(')
		opts := p.parseOperatorDefList()
		p.expect(')')
		return &nodes.AlterTypeStmt{
			TypeName: names,
			Options:  opts,
		}
	default:
		return nil
	}
}

// parseAlterEnumAddValue parses ALTER TYPE name ADD VALUE ...
// Current token is ADD.
func (p *Parser) parseAlterEnumAddValue(typname *nodes.List) *nodes.AlterEnumStmt {
	p.advance() // consume ADD
	p.advance() // consume VALUE

	skipIfExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
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

	return stmt
}

// parseAlterCompositeType parses AlterCompositeTypeStmt (ALTER TYPE name alter_type_cmds).
// Current token is ADD/DROP/ALTER.
func (p *Parser) parseAlterCompositeType(names *nodes.List) *nodes.AlterTableStmt {
	cmds := p.parseAlterTypeCmds()
	rv := makeRangeVarFromCompositeType(names)
	return &nodes.AlterTableStmt{
		Relation: rv,
		Cmds:     cmds,
		ObjType:  int(nodes.OBJECT_TYPE),
	}
}

// parseAlterCompositeTypeRename parses ALTER TYPE name RENAME ATTRIBUTE name TO name opt_drop_behavior.
// RENAME has been consumed. Current token is ATTRIBUTE.
func (p *Parser) parseAlterCompositeTypeRename(names *nodes.List) *nodes.RenameStmt {
	p.advance() // consume ATTRIBUTE
	subname, _ := p.parseName()
	p.expect(TO)
	newname, _ := p.parseName()
	behavior := p.parseOptDropBehavior()
	return &nodes.RenameStmt{
		RenameType:   nodes.OBJECT_ATTRIBUTE,
		RelationType: nodes.OBJECT_TYPE,
		Relation:     makeRangeVarFromAnyName(names),
		Subname:      subname,
		Newname:      newname,
		Behavior:     nodes.DropBehavior(behavior),
	}
}

// makeRangeVarFromCompositeType creates a RangeVar from composite type names with Inh=true.
func makeRangeVarFromCompositeType(names *nodes.List) *nodes.RangeVar {
	rv := &nodes.RangeVar{
		Inh:      true,
		Location: -1,
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
func (p *Parser) parseAlterTypeCmds() *nodes.List {
	cmd := p.parseAlterTypeCmd()
	items := []nodes.Node{cmd}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseAlterTypeCmd())
	}
	return &nodes.List{Items: items}
}

// parseAlterTypeCmd parses alter_type_cmd.
func (p *Parser) parseAlterTypeCmd() *nodes.AlterTableCmd {
	switch p.cur.Type {
	case ADD_P:
		p.advance()
		p.expect(ATTRIBUTE)
		elem := p.parseTableFuncElement()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:  int(nodes.AT_AddColumn),
			Def:      elem,
			Behavior: behavior,
		}
	case DROP:
		p.advance()
		p.expect(ATTRIBUTE)
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		colname, _ := p.parseColId()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropColumn),
			Name:       colname,
			Behavior:   behavior,
			Missing_ok: missingOk,
		}
	case ALTER:
		p.advance()
		p.expect(ATTRIBUTE)
		colname, _ := p.parseColId()
		// SET DATA TYPE or just TYPE
		if p.cur.Type == SET {
			p.advance()
			p.expect(DATA_P)
		}
		p.expect(TYPE_P)
		typename, _ := p.parseTypename()
		var collClause *nodes.CollateClause
		if p.cur.Type == COLLATE {
			p.advance()
			collname, _ := p.parseAnyName()
			collClause = &nodes.CollateClause{Collname: collname, Location: -1}
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
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER DOMAIN
// ---------------------------------------------------------------------------

// parseAlterDomainOwnerOrOther parses ALTER DOMAIN ...
// ALTER has already been consumed. Current token is DOMAIN_P.
func (p *Parser) parseAlterDomainOwnerOrOther() nodes.Node {
	p.advance() // consume DOMAIN

	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance() // consume SET
			p.advance() // consume SCHEMA
			newschema, _ := p.parseName()
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_DOMAIN,
				Object:     names,
				Newschema:  newschema,
			}
		}
		if next.Type == NOT {
			// SET NOT NULL
			p.advance() // consume SET
			p.advance() // consume NOT
			p.expect(NULL_P)
			return &nodes.AlterDomainStmt{
				Subtype: 'O',
				Typname: names,
			}
		}
		if next.Type == DEFAULT {
			// SET DEFAULT a_expr
			p.advance() // consume SET
			p.advance() // consume DEFAULT
			expr := p.parseAExpr(0)
			return &nodes.AlterDomainStmt{
				Subtype: 'T',
				Typname: names,
				Def:     expr,
			}
		}
		return nil
	case DROP:
		p.advance() // consume DROP
		if p.cur.Type == NOT {
			// DROP NOT NULL
			p.advance()
			p.expect(NULL_P)
			return &nodes.AlterDomainStmt{
				Subtype: 'N',
				Typname: names,
			}
		}
		if p.cur.Type == DEFAULT {
			// DROP DEFAULT
			p.advance()
			return &nodes.AlterDomainStmt{
				Subtype: 'T',
				Typname: names,
			}
		}
		if p.cur.Type == CONSTRAINT {
			// DROP CONSTRAINT [IF EXISTS] name opt_drop_behavior
			p.advance()
			missingOk := false
			if p.cur.Type == IF_P {
				p.advance()
				p.expect(EXISTS)
				missingOk = true
			}
			cname, _ := p.parseName()
			behavior := p.parseOptDropBehavior()
			return &nodes.AlterDomainStmt{
				Subtype:   'X',
				Typname:   names,
				Name:      cname,
				Behavior:  nodes.DropBehavior(behavior),
				MissingOk: missingOk,
			}
		}
		return nil
	case ADD_P:
		// ADD [CONSTRAINT name] CHECK (expr)
		p.advance() // consume ADD
		return p.parseDomainConstraintForAlter(names)
	case VALIDATE:
		p.advance() // consume VALIDATE
		p.expect(CONSTRAINT)
		cname, _ := p.parseName()
		return &nodes.AlterDomainStmt{
			Subtype: 'V',
			Typname: names,
			Name:    cname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_DOMAIN,
			Object:     names,
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		if p.cur.Type == CONSTRAINT {
			p.advance()
			subname, _ := p.parseName()
			p.expect(TO)
			newname, _ := p.parseName()
			return &nodes.RenameStmt{
				RenameType: nodes.OBJECT_DOMCONSTRAINT,
				Object:     names,
				Subname:    subname,
				Newname:    newname,
			}
		}
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_DOMAIN,
			Object:     names,
			Newname:    newname,
		}
	default:
		return nil
	}
}

// parseDomainConstraintForAlter parses DomainConstraint for ALTER DOMAIN ... ADD.
func (p *Parser) parseDomainConstraintForAlter(typname *nodes.List) *nodes.AlterDomainStmt {
	var conname string
	if p.cur.Type == CONSTRAINT {
		p.advance()
		conname, _ = p.parseName()
	}
	// DomainConstraintElem: CHECK '(' a_expr ')' ConstraintAttributeSpec
	p.expect(CHECK)
	p.expect('(')
	expr := p.parseAExpr(0)
	p.expect(')')
	// ConstraintAttributeSpec - parse optional attributes
	constraint := &nodes.Constraint{
		Contype:  nodes.CONSTR_CHECK,
		Conname:  conname,
		RawExpr:  expr,
		Location: -1,
	}
	attrs := p.parseConstraintAttributeSpec()
	applyConstraintAttrs(constraint, attrs)
	return &nodes.AlterDomainStmt{
		Subtype: 'C',
		Typname: typname,
		Def:     constraint,
	}
}

// ---------------------------------------------------------------------------
// ALTER SCHEMA
// ---------------------------------------------------------------------------

// parseAlterSchemaOwner parses ALTER SCHEMA ...
// ALTER has already been consumed. Current token is SCHEMA.
func (p *Parser) parseAlterSchemaOwner() nodes.Node {
	p.advance() // consume SCHEMA
	name, _ := p.parseName()

	switch p.cur.Type {
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_SCHEMA,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_SCHEMA,
			Subname:    name,
			Newname:    newname,
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER COLLATION
// ---------------------------------------------------------------------------

// parseAlterCollationStmt parses ALTER COLLATION ...
// ALTER has already been consumed. Current token is COLLATION.
func (p *Parser) parseAlterCollationStmt() nodes.Node {
	p.advance() // consume COLLATION
	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case REFRESH:
		p.advance()
		p.expect(VERSION_P)
		return &nodes.AlterCollationStmt{Collname: names}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newowner:   roleSpec,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_COLLATION,
			Object:     names,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER CONVERSION
// ---------------------------------------------------------------------------

// parseAlterConversionStmt parses ALTER CONVERSION ...
// ALTER has already been consumed. Current token is CONVERSION_P.
func (p *Parser) parseAlterConversionStmt() nodes.Node {
	p.advance() // consume CONVERSION
	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newowner:   roleSpec,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_CONVERSION,
			Object:     names,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER AGGREGATE
// ---------------------------------------------------------------------------

// parseAlterAggregateStmt parses ALTER AGGREGATE ...
// ALTER has already been consumed. Current token is AGGREGATE.
func (p *Parser) parseAlterAggregateStmt() nodes.Node {
	p.advance() // consume AGGREGATE
	agg := p.parseAggregateWithArgtypesLocal()

	switch p.cur.Type {
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newowner:   roleSpec,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_AGGREGATE,
			Object:     agg,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// parseAggregateWithArgtypesLocal parses aggregate_with_argtypes: func_name aggr_args
func (p *Parser) parseAggregateWithArgtypesLocal() *nodes.ObjectWithArgs {
	funcname, _ := p.parseFuncName()
	aggrArgs := p.parseAggrArgs()
	return &nodes.ObjectWithArgs{
		Objname: funcname,
		Objargs: extractAggrArgTypesLocal(aggrArgs),
	}
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
func (p *Parser) parseAlterTextSearchStmt() nodes.Node {
	p.advance() // consume TEXT
	p.expect(SEARCH)

	switch p.cur.Type {
	case DICTIONARY:
		return p.parseAlterTSDictionary()
	case CONFIGURATION:
		return p.parseAlterTSConfiguration()
	case PARSER:
		return p.parseAlterTSParserOrTemplate(nodes.OBJECT_TSPARSER)
	case TEMPLATE:
		return p.parseAlterTSParserOrTemplate(nodes.OBJECT_TSTEMPLATE)
	default:
		return nil
	}
}

// parseAlterTSDictionary parses ALTER TEXT SEARCH DICTIONARY ...
func (p *Parser) parseAlterTSDictionary() nodes.Node {
	p.advance() // consume DICTIONARY
	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case '(':
		// ALTER TEXT SEARCH DICTIONARY name definition
		def := p.parseDefinition()
		return &nodes.AlterTSDictionaryStmt{
			Dictname: names,
			Options:  def,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newowner:   roleSpec,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_TSDICTIONARY,
			Object:     names,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// parseAlterTSConfiguration parses ALTER TEXT SEARCH CONFIGURATION ...
func (p *Parser) parseAlterTSConfiguration() nodes.Node {
	p.advance() // consume CONFIGURATION
	cfgname, _ := p.parseAnyName()

	switch p.cur.Type {
	case ADD_P:
		// ADD MAPPING FOR name_list WITH any_name_list
		p.advance()
		p.expect(MAPPING)
		p.expect(FOR)
		tokentype, _ := p.parseNameList()
		p.expect(WITH)
		dicts, _ := p.parseAnyNameList()
		return &nodes.AlterTSConfigurationStmt{
			Kind:      nodes.ALTER_TSCONFIG_ADD_MAPPING,
			Cfgname:   cfgname,
			Tokentype: tokentype,
			Dicts:     dicts,
		}
	case ALTER:
		// ALTER MAPPING ...
		p.advance()
		p.expect(MAPPING)
		if p.cur.Type == FOR {
			// ALTER MAPPING FOR name_list WITH/REPLACE ...
			p.advance()
			tokentype, _ := p.parseNameList()
			if p.cur.Type == WITH {
				p.advance()
				dicts, _ := p.parseAnyNameList()
				return &nodes.AlterTSConfigurationStmt{
					Kind:      nodes.ALTER_TSCONFIG_ALTER_MAPPING_FOR_TOKEN,
					Cfgname:   cfgname,
					Tokentype: tokentype,
					Dicts:     dicts,
					Override:  true,
				}
			}
			// ALTER MAPPING FOR name_list REPLACE any_name WITH any_name
			p.expect(REPLACE)
			oldDict, _ := p.parseAnyName()
			p.expect(WITH)
			newDict, _ := p.parseAnyName()
			return &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_REPLACE_DICT_FOR_TOKEN,
				Cfgname:   cfgname,
				Tokentype: tokentype,
				Dicts:     &nodes.List{Items: []nodes.Node{oldDict, newDict}},
				Replace:   true,
			}
		}
		// ALTER MAPPING REPLACE any_name WITH any_name
		p.expect(REPLACE)
		oldDict, _ := p.parseAnyName()
		p.expect(WITH)
		newDict, _ := p.parseAnyName()
		return &nodes.AlterTSConfigurationStmt{
			Kind:    nodes.ALTER_TSCONFIG_REPLACE_DICT,
			Cfgname: cfgname,
			Dicts:   &nodes.List{Items: []nodes.Node{oldDict, newDict}},
			Replace: true,
		}
	case DROP:
		// DROP MAPPING [IF EXISTS] FOR name_list
		p.advance()
		p.expect(MAPPING)
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		p.expect(FOR)
		tokentype, _ := p.parseNameList()
		return &nodes.AlterTSConfigurationStmt{
			Kind:      nodes.ALTER_TSCONFIG_DROP_MAPPING,
			Cfgname:   cfgname,
			Tokentype: tokentype,
			MissingOk: missingOk,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newowner:   roleSpec,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_TSCONFIGURATION,
			Object:     cfgname,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// parseAlterTSParserOrTemplate parses ALTER TEXT SEARCH PARSER/TEMPLATE ...
func (p *Parser) parseAlterTSParserOrTemplate(objtype nodes.ObjectType) nodes.Node {
	p.advance() // consume PARSER or TEMPLATE
	names, _ := p.parseAnyName()

	switch p.cur.Type {
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: objtype,
			Object:     names,
			Newname:    newname,
		}
	case SET:
		p.advance()
		p.expect(SCHEMA)
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: objtype,
			Object:     names,
			Newschema:  newschema,
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER LANGUAGE
// ---------------------------------------------------------------------------

// parseAlterLanguageStmt parses ALTER [PROCEDURAL] LANGUAGE ...
// ALTER has already been consumed. Current token is LANGUAGE or PROCEDURAL.
func (p *Parser) parseAlterLanguageStmt() nodes.Node {
	if p.cur.Type == PROCEDURAL {
		p.advance() // consume PROCEDURAL
	}
	p.advance() // consume LANGUAGE
	name, _ := p.parseName()

	switch p.cur.Type {
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_LANGUAGE,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_LANGUAGE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER LARGE OBJECT
// ---------------------------------------------------------------------------

// parseAlterLargeObjectStmt parses ALTER LARGE OBJECT ...
// ALTER has already been consumed. Current token is LARGE_P.
func (p *Parser) parseAlterLargeObjectStmt() nodes.Node {
	p.advance() // consume LARGE
	p.expect(OBJECT_P)
	numVal := p.parseNumericOnly()
	p.expect(OWNER)
	p.expect(TO)
	roleSpec := p.parseRoleSpec()
	return &nodes.AlterOwnerStmt{
		ObjectType: nodes.OBJECT_LARGEOBJECT,
		Object:     numVal,
		Newowner:   roleSpec,
	}
}

// ---------------------------------------------------------------------------
// ALTER EVENT TRIGGER
// ---------------------------------------------------------------------------

// parseAlterEventTriggerOwner parses ALTER EVENT TRIGGER ...
// ALTER has already been consumed. Current token is EVENT.
func (p *Parser) parseAlterEventTriggerOwner() nodes.Node {
	p.advance() // consume EVENT
	p.expect(TRIGGER)
	name, _ := p.parseName()

	switch p.cur.Type {
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_EVENT_TRIGGER,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_EVENT_TRIGGER,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}
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
		}
	case DISABLE_P:
		p.advance()
		return &nodes.AlterEventTrigStmt{
			Trigname:  name,
			Tgenabled: 'D',
		}
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER TABLESPACE
// ---------------------------------------------------------------------------

// parseAlterTablespaceOwner parses ALTER TABLESPACE ...
// ALTER has already been consumed. Current token is TABLESPACE.
func (p *Parser) parseAlterTablespaceOwner() nodes.Node {
	p.advance() // consume TABLESPACE
	name, _ := p.parseName()

	switch p.cur.Type {
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_TABLESPACE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TABLESPACE,
			Subname:    name,
			Newname:    newname,
		}
	case SET:
		p.advance()
		// SET reloptions or SET (...)
		opts := p.parseDefinition()
		_ = opts
		return nil // Tablespace SET (...) is handled differently in yacc, skip for now
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------
// ALTER TRIGGER ... DEPENDS ON EXTENSION
// ---------------------------------------------------------------------------

// parseAlterTriggerDependsOnExtension parses ALTER TRIGGER ...
// ALTER has already been consumed. Current token is TRIGGER.
func (p *Parser) parseAlterTriggerDependsOnExtension() nodes.Node {
	p.advance() // consume TRIGGER
	trigname, _ := p.parseName()
	p.expect(ON)
	// qualified_name
	relname, _ := p.parseAnyName()
	rel := makeRangeVarFromAnyName(relname)

	remove := p.parseOptNo()
	p.expect(DEPENDS)
	p.expect(ON)
	p.expect(EXTENSION)
	extname, _ := p.parseName()

	return &nodes.AlterObjectDependsStmt{
		ObjectType: nodes.OBJECT_TRIGGER,
		Relation:   rel,
		Object:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: trigname}}},
		Extname:    &nodes.String{Str: extname},
		Remove:     remove,
	}
}
