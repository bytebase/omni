package parser

import (
	"fmt"
	"strings"
)

func (p *parser) parseAttachOrCreateStmt() (Statement, error) {
	// In ClickHouse, ATTACH DICTIONARY only attaches an existing dictionary by
	// name (no schema/SOURCE/LAYOUT/LIFETIME clause).
	if p.peek(1).Type == KwDictionary {
		start := p.advance().Pos // ATTACH
		p.advance()              // DICTIONARY
		p.matchIfNotExists()
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.parseUUIDClause(); err != nil {
			return nil, err
		}
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		return &AttachStmt{
			span:       span{Pos: start, End: p.peek(-1).End},
			ObjectType: "DICTIONARY",
			Database:   db,
			Name:       name,
			Cluster:    cluster,
		}, nil
	}
	return p.parseCreateStmt()
}

func (p *parser) isStatementEndOrOutput() bool {
	t := p.cur().Type
	return t == TokEOF || t == TokSemicolon || t == KwFormat || (t == KwInto && p.peek(1).Type == KwOutfile)
}

func (p *parser) parseCreateStmt() (Statement, error) {
	start := p.cur().Pos
	action := strings.ToUpper(p.advance().Val) // CREATE, ATTACH, or REPLACE

	orReplace := false
	// ATTACH takes OR REPLACE before VIEW only: ClickHouseParser.g4:103
	// writes (ATTACH | CREATE) (OR REPLACE)? VIEW, and ParserCreateViewQuery
	// reads OR REPLACE after either keyword.
	if (action == "CREATE" || action == "ATTACH" && p.peek(2).Type == KwView) &&
		p.cur().Type == KwOr && p.peek(1).Type == KwReplace {
		p.advance()
		p.advance()
		orReplace = true
	}

	switch p.cur().Type {
	case KwDatabase:
		if action == "REPLACE" || orReplace {
			return nil, p.errAt(p.cur(), "OR REPLACE / REPLACE is not allowed for DATABASE")
		}
		return p.parseCreateDatabaseStmt(start, action == "ATTACH")
	case KwDictionary:
		return p.parseCreateDictionaryStmt(start, action, orReplace)
	case KwMaterialized:
		if action == "REPLACE" || orReplace {
			return nil, p.errAt(p.cur(), "OR REPLACE / REPLACE is not allowed for MATERIALIZED VIEW")
		}
		return p.parseCreateMaterializedViewStmt(start, action)
	case KwView:
		if action == "REPLACE" {
			return nil, p.errAt(p.cur(), "REPLACE VIEW is not supported; use CREATE OR REPLACE VIEW")
		}
		return p.parseCreateViewStmt(start, action, orReplace)
	case KwTemporary, KwTable:
		return p.parseCreateTableStmt(start, action, orReplace)
	case KwFunction:
		return p.parseCreateFunctionStmt(start, orReplace)
	case KwUser, KwRole, KwQuota, KwSettings, KwProfile, KwRow, KwPolicy:
		// An access entity takes OR REPLACE after its keyword, as
		// ParserCreateUserQuery and its siblings read it: the server answers
		// CREATE OR REPLACE ROLE with a syntax error at ROLE.
		if orReplace {
			return nil, p.errAt(p.cur(), "OR REPLACE follows the entity: CREATE "+p.cur().Raw+" OR REPLACE")
		}
		return p.parseCreateAccessStmt(start, action)
	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after %s", p.cur().Raw, action))
	}
}

func (p *parser) parseCreateDatabaseStmt(start int, attach bool) (*CreateDatabaseStmt, error) {
	p.advance() // DATABASE
	ifNotExists := p.matchIfNotExists()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	var eng *EngineExpr
	if p.cur().Type == KwEngine {
		eng, err = p.parseEngineExpr()
		if err != nil {
			return nil, err
		}
	}
	var comment string
	if p.match(KwComment) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		comment = lit.Val
	}
	return &CreateDatabaseStmt{
		span:        span{Pos: start, End: p.peek(-1).End},
		Attach:      attach,
		IfNotExists: ifNotExists,
		Name:        name,
		Cluster:     cluster,
		Engine:      eng,
		Comment:     comment,
	}, nil
}

func (p *parser) parseCreateTableStmt(start int, action string, orReplace bool) (*CreateTableStmt, error) {
	temp := p.match(KwTemporary)
	if _, err := p.expect(KwTable); err != nil {
		return nil, err
	}
	ifNotExists := p.matchIfNotExists()
	db, name, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	uuid, err := p.parseUUIDClause()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt := &CreateTableStmt{
		span:        span{Pos: start},
		Action:      action,
		OrReplace:   orReplace,
		Temporary:   temp,
		IfNotExists: ifNotExists,
		Database:    db,
		Name:        name,
		UUID:        uuid,
		Cluster:     cluster,
	}

	// Optional tableSchemaClause: `(elements)` or (for CREATE/REPLACE only) `AS table` / `AS func`
	// Careful: `AS SELECT ...` is subqueryClause, NOT tableSchemaClause!
	if p.cur().Type == TokLParen {
		elems, err := p.parseTableSchemaParens()
		if err != nil {
			return nil, err
		}
		stmt.Columns = elems
	} else if action != "ATTACH" && p.cur().Type == KwAs && !p.looksLikeSelectAfter(1) {
		p.advance() // AS
		if isIdentifier(p.cur()) && p.peek(1).Type == TokLParen {
			fn, err := p.parseTableFunctionExpr()
			if err != nil {
				return nil, err
			}
			stmt.AsFunc = fn
		} else {
			asDB, asTbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			if asDB != "" {
				stmt.AsTable = asDB + "." + asTbl
			} else {
				stmt.AsTable = asTbl
			}
		}
	}

	// Optional engineClause (ENGINE itself is optional when relying on default_table_engine)
	if p.isEngineClauseStart() {
		ec, err := p.parseEngineClause()
		if err != nil {
			return nil, err
		}
		stmt.Engine = ec
	}

	// Optional subqueryClause: AS selectUnionStmt (CREATE/REPLACE only)
	if action != "ATTACH" && p.match(KwAs) {
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		stmt.Subquery = sub
	}

	if p.match(KwComment) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Comment = lit.Val
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) looksLikeSelectAfter(offset int) bool {
	save := p.pos
	p.pos += offset
	ok := p.looksLikeSelect()
	p.pos = save
	return ok
}

func (p *parser) parseCreateMaterializedViewStmt(start int, action string) (*CreateViewStmt, error) {
	p.advance() // MATERIALIZED
	if _, err := p.expect(KwView); err != nil {
		return nil, err
	}
	ifNotExists := p.matchIfNotExists()
	db, name, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	uuid, err := p.parseUUIDClause()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt := &CreateViewStmt{
		span:         span{Pos: start},
		Action:       action,
		Materialized: true,
		IfNotExists:  ifNotExists,
		Database:     db,
		Name:         name,
		UUID:         uuid,
		Cluster:      cluster,
	}
	if p.cur().Type == TokLParen {
		elems, err := p.parseTableSchemaParens()
		if err != nil {
			return nil, err
		}
		stmt.Columns = elems
	}
	if p.match(KwTo) {
		toDB, toTbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		if toDB != "" {
			stmt.ToTable = toDB + "." + toTbl
		} else {
			stmt.ToTable = toTbl
		}
	} else if p.cur().Type == KwEngine {
		ec, err := p.parseEngineClause()
		if err != nil {
			return nil, err
		}
		stmt.Engine = ec
		stmt.Populate = p.match(KwPopulate)
	} else {
		return nil, p.errAt(p.cur(), "expected TO or ENGINE in CREATE MATERIALIZED VIEW")
	}
	// Optional schema after TO table (ClickHouse permits `TO dst (col Type) AS SELECT`)
	if len(stmt.Columns) == 0 && p.cur().Type == TokLParen {
		elems, err := p.parseTableSchemaParens()
		if err != nil {
			return nil, err
		}
		stmt.Columns = elems
	}
	if _, err := p.expect(KwAs); err != nil {
		return nil, err
	}
	sub, err := p.parseSelectUnionStmt()
	if err != nil {
		return nil, err
	}
	stmt.Subquery = sub
	if p.match(KwComment) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Comment = lit.Val
	}
	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseCreateViewStmt(start int, action string, orReplace bool) (*CreateViewStmt, error) {
	p.advance() // VIEW
	ifNotExists := p.matchIfNotExists()
	db, name, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	uuid, err := p.parseUUIDClause()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt := &CreateViewStmt{
		span:        span{Pos: start},
		Action:      action,
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Database:    db,
		Name:        name,
		UUID:        uuid,
		Cluster:     cluster,
	}
	if p.cur().Type == TokLParen {
		elems, err := p.parseTableSchemaParens()
		if err != nil {
			return nil, err
		}
		stmt.Columns = elems
	}
	if _, err := p.expect(KwAs); err != nil {
		return nil, err
	}
	sub, err := p.parseSelectUnionStmt()
	if err != nil {
		return nil, err
	}
	stmt.Subquery = sub
	if p.match(KwComment) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Comment = lit.Val
	}
	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseCreateDictionaryStmt(start int, action string, orReplace bool) (*CreateDictionaryStmt, error) {
	p.advance() // DICTIONARY
	ifNotExists := p.matchIfNotExists()
	db, name, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	uuid, err := p.parseUUIDClause()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	// dictionarySchemaClause: '(' dictionaryAttrDfnt (',' dictionaryAttrDfnt)* ')'
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	var attrs []*DictionaryAttr
	for {
		attr, err := p.parseDictionaryAttrDfnt()
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, attr)
		if !p.match(TokComma) {
			break
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return nil, err
	}

	// dictionaryEngineClause
	var primaryKey []Expr
	if p.cur().Type == KwPrimary && p.peek(1).Type == KwKey {
		p.advance()
		p.advance()
		for {
			var expr Expr
			var err error
			if p.cur().Type == TokStar {
				expr, err = p.parseStarExpr(nil, p.cur().Pos)
			} else {
				expr, err = p.parseExpr()
			}
			if err != nil {
				return nil, err
			}
			primaryKey = append(primaryKey, expr)
			if !p.match(TokComma) {
				break
			}
		}
	}
	var (
		source   *DictionaryConfig
		layout   *DictionaryConfig
		lifetime *DictionaryLifetime
		rng      *DictionaryRange
		settings []*SettingExpr
	)
	seen := map[string]bool{}
	for {
		switch p.cur().Type {
		case KwSource:
			if seen["source"] {
				return nil, p.errAt(p.cur(), "duplicate SOURCE clause")
			}
			seen["source"] = true
			var err error
			source, err = p.parseDictSourceOrLayout()
			if err != nil {
				return nil, err
			}
		case KwLifetime:
			if seen["lifetime"] {
				return nil, p.errAt(p.cur(), "duplicate LIFETIME clause")
			}
			seen["lifetime"] = true
			var err error
			lifetime, err = p.parseDictLifetime()
			if err != nil {
				return nil, err
			}
		case KwLayout:
			if seen["layout"] {
				return nil, p.errAt(p.cur(), "duplicate LAYOUT clause")
			}
			seen["layout"] = true
			var err error
			layout, err = p.parseDictSourceOrLayout()
			if err != nil {
				return nil, err
			}
		case KwRange:
			if seen["range"] {
				return nil, p.errAt(p.cur(), "duplicate RANGE clause")
			}
			seen["range"] = true
			var err error
			rng, err = p.parseDictRange()
			if err != nil {
				return nil, err
			}
		case KwSettings:
			if seen["settings"] {
				return nil, p.errAt(p.cur(), "duplicate SETTINGS clause")
			}
			seen["settings"] = true
			p.advance()
			if _, err := p.expect(TokLParen); err != nil {
				return nil, err
			}
			var err error
			settings, err = p.parseSettingExprList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokRParen); err != nil {
				return nil, err
			}
		default:
			goto doneDict
		}
	}
doneDict:
	return &CreateDictionaryStmt{
		span:        span{Pos: start, End: p.peek(-1).End},
		Action:      action,
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Database:    db,
		Name:        name,
		UUID:        uuid,
		Cluster:     cluster,
		Attributes:  attrs,
		PrimaryKey:  primaryKey,
		Source:      source,
		Layout:      layout,
		Lifetime:    lifetime,
		Range:       rng,
		Settings:    settings,
	}, nil
}

func (p *parser) parseDictionaryAttrDfnt() (*DictionaryAttr, error) {
	start := p.cur().Pos
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	attr := &DictionaryAttr{
		span: span{Pos: start, End: dt.End},
		Name: name,
		Type: dt,
	}
	seen := map[string]bool{}
	for {
		switch p.cur().Type {
		case KwDefault:
			if seen["default"] {
				return nil, p.errAt(p.cur(), "duplicate DEFAULT attribute")
			}
			seen["default"] = true
			p.advance()
			lit, err := p.parseLiteral()
			if err != nil {
				return nil, err
			}
			attr.Default = lit
		case KwExpression:
			if seen["expression"] {
				return nil, p.errAt(p.cur(), "duplicate EXPRESSION attribute")
			}
			seen["expression"] = true
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			attr.Expression = expr
		case KwHierarchical:
			if seen["hierarchical"] {
				return nil, p.errAt(p.cur(), "duplicate HIERARCHICAL attribute")
			}
			seen["hierarchical"] = true
			p.advance()
			attr.Hierarchical = true
		case KwInjective:
			if seen["injective"] {
				return nil, p.errAt(p.cur(), "duplicate INJECTIVE attribute")
			}
			seen["injective"] = true
			p.advance()
			attr.Injective = true
		case KwIsObjectId:
			if seen["is_object_id"] {
				return nil, p.errAt(p.cur(), "duplicate IS_OBJECT_ID attribute")
			}
			seen["is_object_id"] = true
			p.advance()
			attr.IsObjectID = true
		default:
			attr.End = p.peek(-1).End
			return attr, nil
		}
	}
}

func (p *parser) parseDictSourceOrLayout() (*DictionaryConfig, error) {
	start := p.advance().Pos // SOURCE or LAYOUT
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cfg, err := p.parseDictConfigBody(start, name)
	if err != nil {
		return nil, err
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	cfg.End = rp.End
	return cfg, nil
}

func (p *parser) parseDictConfigBody(start int, name string) (*DictionaryConfig, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	var params []*SettingExpr
	var children []*DictionaryConfig
	for p.cur().Type != TokRParen && p.cur().Type != TokEOF {
		kStart := p.cur().Pos
		key, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if p.cur().Type == TokLParen {
			child, err := p.parseDictConfigBody(kStart, key)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
			continue
		}
		var val Expr
		if isLiteralStart(p.cur(), p.peek(1)) {
			val, err = p.parseLiteral()
			if err != nil {
				return nil, err
			}
		} else {
			vStart := p.cur().Pos
			vName, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if p.match(TokLParen) {
				rp, err := p.expect(TokRParen)
				if err != nil {
					return nil, err
				}
				val = &FuncCallExpr{span: span{Pos: vStart, End: rp.End}, Name: vName}
			} else {
				val = &IdentExpr{span: span{Pos: vStart, End: p.peek(-1).End}, Parts: []string{vName}}
			}
		}
		params = append(params, &SettingExpr{
			span:  span{Pos: kStart, End: val.EndPos()},
			Name:  key,
			Value: val,
		})
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	return &DictionaryConfig{
		span:     span{Pos: start, End: rp.End},
		Name:     name,
		Params:   params,
		Children: children,
	}, nil
}

func (p *parser) parseDictLifetime() (*DictionaryLifetime, error) {
	start := p.advance().Pos // LIFETIME
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	lt := &DictionaryLifetime{span: span{Pos: start}}
	if p.match(KwMin) {
		minTok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwMax); err != nil {
			return nil, err
		}
		maxTok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		lt.Min, lt.Max = minTok.Raw, maxTok.Raw
	} else if p.match(KwMax) {
		maxTok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwMin); err != nil {
			return nil, err
		}
		minTok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		lt.Min, lt.Max = minTok.Raw, maxTok.Raw
	} else {
		maxTok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		lt.Max = maxTok.Raw
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	lt.End = rp.End
	return lt, nil
}

func (p *parser) parseDictRange() (*DictionaryRange, error) {
	start := p.advance().Pos // RANGE
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	rng := &DictionaryRange{span: span{Pos: start}}
	if p.match(KwMin) {
		minCol, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwMax); err != nil {
			return nil, err
		}
		maxCol, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		rng.Min, rng.Max = minCol, maxCol
	} else if p.match(KwMax) {
		maxCol, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwMin); err != nil {
			return nil, err
		}
		minCol, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		rng.Min, rng.Max = minCol, maxCol
	} else {
		return nil, p.errAt(p.cur(), "expected MIN or MAX in RANGE clause")
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	rng.End = rp.End
	return rng, nil
}

func (p *parser) parseCreateFunctionStmt(start int, orReplace bool) (*CreateFunctionStmt, error) {
	p.advance() // FUNCTION
	ifNotExists := p.matchIfNotExists()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(KwAs); err != nil {
		return nil, err
	}
	body, err := p.parseFuncArg()
	if err != nil {
		return nil, err
	}
	return &CreateFunctionStmt{
		span:        span{Pos: start, End: body.EndPos()},
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
		Cluster:     cluster,
		Body:        body,
	}, nil
}

func (p *parser) parseCreateAccessStmt(start int, action string) (*AccessControlStmt, error) {
	var entity string
	switch p.cur().Type {
	case KwUser:
		p.advance()
		entity = "USER"
	case KwRole:
		p.advance()
		entity = "ROLE"
	case KwQuota:
		p.advance()
		entity = "QUOTA"
	case KwSettings:
		p.advance()
		if _, err := p.expect(KwProfile); err != nil {
			return nil, err
		}
		entity = "SETTINGS PROFILE"
	case KwProfile:
		p.advance()
		entity = "PROFILE"
	case KwRow:
		p.advance()
		if _, err := p.expect(KwPolicy); err != nil {
			return nil, err
		}
		entity = "ROW POLICY"
	case KwPolicy:
		p.advance()
		entity = "POLICY"
	}
	kind := action + " " + entity
	stmt := &AccessControlStmt{
		span: span{Pos: start},
		Kind: kind,
	}

	// CREATE takes IF NOT EXISTS or OR REPLACE after the entity, not both.
	if action == "ALTER" {
		stmt.IfExists = p.matchIfExists()
	} else if stmt.IfNotExists = p.matchIfNotExists(); !stmt.IfNotExists && action == "CREATE" &&
		p.cur().Type == KwOr && p.peek(1).Type == KwReplace {
		p.advance()
		p.advance()
		stmt.OrReplace = true
	}

	// Entity name(s) (comma-separated)
	names, err := p.parseIdentList()
	if err != nil {
		return nil, err
	}
	stmt.Names = names
	// ROW POLICY requires `ON table1 [, table2 ...]`
	if entity == "ROW POLICY" || entity == "POLICY" {
		if _, err := p.expect(KwOn); err != nil {
			return nil, err
		}
		for {
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Tables = append(stmt.Tables, joinDB(db, tbl))
			if !p.match(TokComma) {
				break
			}
		}
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt.Cluster = cluster

	// Structured entity-specific clauses
	switch entity {
	case "USER":
		for !p.isStatementEndOrOutput() {
			switch {
			case p.match(KwIdentified):
				if p.match(KwWith) {
					withType, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					stmt.IdentifiedWith = withType
				}
				if p.match(KwBy) {
					lit, err := p.parseStringLiteral()
					if err != nil {
						return nil, err
					}
					stmt.IdentifiedBy = lit.Val
				}
			case p.cur().Type == KwDefault && p.peek(1).Type == KwRole:
				p.advance()
				p.advance()
				if p.isWord("NONE") {
					stmt.DefaultRoles = append(stmt.DefaultRoles, strings.ToUpper(p.advance().Val))
				} else if p.match(KwAll) {
					stmt.DefaultRoles = append(stmt.DefaultRoles, "ALL")
				} else {
					roles, err := p.parseIdentList()
					if err != nil {
						return nil, err
					}
					stmt.DefaultRoles = append(stmt.DefaultRoles, roles...)
				}
			case p.cur().Type == KwSettings:
				settings, err := p.parseSettingsClause()
				if err != nil {
					return nil, err
				}
				stmt.Settings = append(stmt.Settings, settings...)
			case p.isWord("HOST"):
				p.advance()
				if p.cur().Type == KwAny || p.isWord("NONE") || p.isWord("LOCAL") {
					stmt.Hosts = append(stmt.Hosts, strings.ToUpper(p.advance().Val))
				} else {
					for {
						hType, err := p.parseIdentifier()
						if err != nil {
							return nil, err
						}
						hVal, err := p.parseStringLiteral()
						if err != nil {
							return nil, err
						}
						stmt.Hosts = append(stmt.Hosts, strings.ToUpper(hType)+" "+hVal.Val)
						if !p.match(TokComma) {
							break
						}
					}
				}
			default:
				return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q in %s", p.cur().Raw, kind))
			}
		}
	case "ROLE", "SETTINGS PROFILE", "PROFILE":
		if p.cur().Type == KwSettings {
			settings, err := p.parseSettingsClause()
			if err != nil {
				return nil, err
			}
			stmt.Settings = settings
		}
		if p.match(KwTo) {
			grantees, except, err := p.parseGranteeList()
			if err != nil {
				return nil, err
			}
			stmt.Grantees = grantees
			stmt.ExceptGrantees = except
		}
	case "ROW POLICY", "POLICY":
		for !p.isStatementEndOrOutput() {
			switch {
			case p.match(KwFor):
				if _, err := p.expect(KwSelect); err != nil {
					return nil, err
				}
				stmt.ForSelect = true
			case p.match(KwUsing):
				using, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				stmt.Using = using
			case p.match(KwTo):
				grantees, except, err := p.parseGranteeList()
				if err != nil {
					return nil, err
				}
				stmt.Grantees = grantees
				stmt.ExceptGrantees = except
			default:
				return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q in %s", p.cur().Raw, kind))
			}
		}
	case "QUOTA":
		if p.match(KwTo) {
			grantees, except, err := p.parseGranteeList()
			if err != nil {
				return nil, err
			}
			stmt.Grantees = grantees
			stmt.ExceptGrantees = except
		}
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) matchIfNotExists() bool {
	if p.cur().Type == KwIf && p.peek(1).Type == KwNot && p.peek(2).Type == KwExists {
		p.advance()
		p.advance()
		p.advance()
		return true
	}
	return false
}

func (p *parser) matchIfExists() bool {
	if p.cur().Type == KwIf && p.peek(1).Type == KwExists {
		p.advance()
		p.advance()
		return true
	}
	return false
}

func (p *parser) parseUUIDClause() (string, error) {
	if !p.match(KwUuid) {
		return "", nil
	}
	lit, err := p.parseStringLiteral()
	if err != nil {
		return "", err
	}
	// ClickHouse validates UUID format inside ParserCreateQuery (Code 376 CANNOT_PARSE_UUID).
	if !isValidClickHouseUUID(lit.Val) {
		return "", p.errAt(p.peek(-1), fmt.Sprintf("cannot parse UUID %q", lit.Val))
	}
	return lit.Val, nil
}

func isValidClickHouseUUID(s string) bool {
	if len(s) == 36 {
		for i := 0; i < 36; i++ {
			if i == 8 || i == 13 || i == 18 || i == 23 {
				if s[i] != '-' {
					return false
				}
			} else if !isHexDigit(s[i]) {
				return false
			}
		}
		return true
	}
	if len(s) == 32 {
		for i := 0; i < 32; i++ {
			if !isHexDigit(s[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (p *parser) parseTableSchemaParens() ([]*TableElement, error) {
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	var elems []*TableElement
	for {
		elem, err := p.parseTableElementExpr()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
		if !p.match(TokComma) {
			break
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return nil, err
	}
	return elems, nil
}

func (p *parser) parseTableElementExpr() (*TableElement, error) {
	start := p.cur().Pos
	switch p.cur().Type {
	case KwConstraint:
		if isIdentifier(p.peek(1)) && p.peek(2).Type == KwCheck {
			p.advance() // CONSTRAINT
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			p.advance() // CHECK
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &TableElement{span: span{Pos: start, End: expr.EndPos()}, Kind: "CONSTRAINT", Name: name, Expr: expr}, nil
		}
	case KwIndex:
		// Disambiguate INDEX definition vs column named `index`
		if isIdentifier(p.peek(1)) {
			save := p.pos
			p.advance() // INDEX
			if elem, err := p.parseTableIndexDfnt(start); err == nil {
				return elem, nil
			}
			p.pos = save
		}
	case KwProjection:
		// Disambiguate PROJECTION definition (`PROJECTION name ([WITH ...] SELECT ...)`)
		// from a column named `projection` with a parameterized type (`projection Array(String)`).
		if isIdentifier(p.peek(1)) && p.peek(2).Type == TokLParen &&
			(p.peek(3).Type == KwSelect || p.peek(3).Type == KwWith) {
			p.advance() // PROJECTION
			return p.parseTableProjectionDfnt(start)
		}
	}
	return p.parseTableColumnDfnt()
}

func (p *parser) parseTableColumnDfnt() (*TableElement, error) {
	start := p.cur().Pos
	name, err := p.parseNestedIdentifier()
	if err != nil {
		return nil, err
	}
	elem := &TableElement{span: span{Pos: start}, Kind: "COLUMN", Name: name}
	hasType := false
	hasProp := false

	// Type is optional if a column property (DEFAULT / MATERIALIZED / ALIAS) is present.
	if p.cur().Type != KwDefault && p.cur().Type != KwMaterialized && p.cur().Type != KwAlias {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		elem.Type = dt
		hasType = true
		// Optional explicit nullability modifier: NULL or NOT NULL.
		if p.match(KwNull) {
			elem.Nullability = "NULL"
		} else if p.cur().Type == KwNot && p.peek(1).Type == KwNull {
			p.advance()
			p.advance()
			elem.Nullability = "NOT NULL"
		}
	}
	if p.cur().Type == KwDefault || p.cur().Type == KwMaterialized || p.cur().Type == KwAlias {
		elem.DefaultKind = strings.ToUpper(p.advance().Val)
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		elem.Expr = expr
		hasProp = true
	}
	if !hasType && !hasProp {
		return nil, p.errAt(p.cur(), "expected column type or DEFAULT/MATERIALIZED/ALIAS expression")
	}
	// Optional column modifiers (COMMENT, CODEC, TTL) in any order
	for {
		switch p.cur().Type {
		case KwComment:
			p.advance()
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			elem.Comment = lit.Val
		case KwCodec:
			codec, err := p.parseCodecExpr()
			if err != nil {
				return nil, err
			}
			elem.Codec = codec
		case KwTtl:
			p.advance()
			ttlExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			elem.TTLExpr = ttlExpr
		default:
			elem.End = p.peek(-1).End
			return elem, nil
		}
	}
}

func (p *parser) parseTableIndexDfnt(start int) (*TableElement, error) {
	name, err := p.parseNestedIdentifier()
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(KwType); err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	end := dt.End
	granVal := ""
	if p.match(KwGranularity) {
		gran, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		granVal = gran.Raw
		end = gran.End
	}
	return &TableElement{
		span:        span{Pos: start, End: end},
		Kind:        "INDEX",
		Name:        name,
		Type:        dt,
		Expr:        expr,
		Granularity: granVal,
	}, nil
}

func (p *parser) parseTableProjectionDfnt(start int) (*TableElement, error) {
	name, err := p.parseNestedIdentifier()
	if err != nil {
		return nil, err
	}
	lp, err := p.expect(TokLParen)
	if err != nil {
		return nil, err
	}
	proj := &SelectStmt{span: span{Pos: lp.End}}
	if p.match(KwWith) {
		withs, err := p.parseWithExprList()
		if err != nil {
			return nil, err
		}
		proj.With = withs
	}
	selTok, err := p.expect(KwSelect)
	if err != nil {
		return nil, err
	}
	if len(proj.With) == 0 {
		proj.Pos = selTok.Pos
	}
	cols, err := p.parseColumnExprList()
	if err != nil {
		return nil, err
	}
	proj.Columns = cols
	if p.cur().Type == KwGroup && p.peek(1).Type == KwBy {
		if p.peek(2).Type == TokRParen {
			// ClickHouse's ParserProjectionSelectQuery permits an empty GROUP BY
			// immediately before the closing ')'.
			gbStart := p.advance().Pos // GROUP
			byTok := p.advance()       // BY
			proj.GroupBy = &GroupByClause{span: span{Pos: gbStart, End: byTok.End}}
		} else {
			gb, err := p.parseGroupByClause()
			if err != nil {
				return nil, err
			}
			proj.GroupBy = gb
		}
	}
	if p.cur().Type == KwOrder && p.peek(1).Type == KwBy {
		p.advance()
		p.advance()
		orderExprs, err := p.parseColumnExprList()
		if err != nil {
			return nil, err
		}
		for _, oe := range orderExprs {
			proj.OrderBy = append(proj.OrderBy, &OrderByItem{
				span: span{Pos: oe.StartPos(), End: oe.EndPos()},
				Expr: oe,
			})
		}
	}
	proj.End = p.peek(-1).End
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	return &TableElement{
		span:       span{Pos: start, End: rp.End},
		Kind:       "PROJECTION",
		Name:       name,
		Projection: proj,
	}, nil
}

func (p *parser) parseCodecExpr() (*CodecExpr, error) {
	tok, err := p.expect(KwCodec)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	ce := &CodecExpr{span: span{Pos: tok.Pos}}
	for {
		argStart := p.cur().Pos
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		arg := &CodecArg{span: span{Pos: argStart, End: p.peek(-1).End}, Name: name}
		if p.match(TokLParen) {
			if p.cur().Type != TokRParen {
				args, err := p.parseColumnExprList()
				if err != nil {
					return nil, err
				}
				arg.Args = args
			}
			rpArg, err := p.expect(TokRParen)
			if err != nil {
				return nil, err
			}
			arg.End = rpArg.End
		}
		ce.Codecs = append(ce.Codecs, arg)
		if !p.match(TokComma) {
			break
		}
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	ce.End = rp.End
	return ce, nil
}

func (p *parser) parseEngineExpr() (*EngineExpr, error) {
	start := p.advance().Pos // ENGINE
	p.match(TokEq)
	name, err := p.parseIdentifierOrNull()
	if err != nil {
		return nil, err
	}
	eng := &EngineExpr{span: span{Pos: start, End: p.peek(-1).End}, Name: name}
	if p.match(TokLParen) {
		if p.cur().Type != TokRParen {
			args, err := p.parseColumnExprList()
			if err != nil {
				return nil, err
			}
			eng.Args = args
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		eng.End = rp.End
	}
	return eng, nil
}

func (p *parser) isEngineClauseStart() bool {
	switch p.cur().Type {
	case KwEngine, KwTtl, KwSettings:
		return true
	case KwOrder, KwPartition, KwSample:
		return p.peek(1).Type == KwBy
	case KwPrimary:
		return p.peek(1).Type == KwKey
	default:
		return false
	}
}

func (p *parser) parseEngineClause() (*EngineClause, error) {
	start := p.cur().Pos
	ec := &EngineClause{span: span{Pos: start}}
	seen := map[string]bool{}
	if p.cur().Type == KwEngine {
		eng, err := p.parseEngineExpr()
		if err != nil {
			return nil, err
		}
		ec.Engine = eng
		seen["engine"] = true
	}
	for {
		switch {
		case p.cur().Type == KwEngine && !seen["engine"]:
			seen["engine"] = true
			eng, err := p.parseEngineExpr()
			if err != nil {
				return nil, err
			}
			ec.Engine = eng
		case p.cur().Type == KwOrder && p.peek(1).Type == KwBy:
			if seen["order"] {
				return nil, p.errAt(p.cur(), "duplicate ORDER BY clause")
			}
			seen["order"] = true
			p.advance()
			p.advance()
			items, err := p.parseEngineOrderExprList()
			if err != nil {
				return nil, err
			}
			ec.OrderBy = items
		case p.cur().Type == KwPartition && p.peek(1).Type == KwBy:
			if seen["partition"] {
				return nil, p.errAt(p.cur(), "duplicate PARTITION BY clause")
			}
			seen["partition"] = true
			p.advance()
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ec.PartitionBy = expr
		case p.cur().Type == KwPrimary && p.peek(1).Type == KwKey:
			if seen["primary"] {
				return nil, p.errAt(p.cur(), "duplicate PRIMARY KEY clause")
			}
			seen["primary"] = true
			p.advance()
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ec.PrimaryKey = expr
		case p.cur().Type == KwSample && p.peek(1).Type == KwBy:
			if seen["sample"] {
				return nil, p.errAt(p.cur(), "duplicate SAMPLE BY clause")
			}
			seen["sample"] = true
			p.advance()
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ec.SampleBy = expr
		case p.cur().Type == KwTtl:
			if seen["ttl"] {
				return nil, p.errAt(p.cur(), "duplicate TTL clause")
			}
			seen["ttl"] = true
			ttls, err := p.parseTTLClause(false)
			if err != nil {
				return nil, err
			}
			ec.TTL = ttls
		case p.cur().Type == KwSettings:
			if seen["settings"] {
				return nil, p.errAt(p.cur(), "duplicate SETTINGS clause")
			}
			seen["settings"] = true
			settings, err := p.parseSettingsClause()
			if err != nil {
				return nil, err
			}
			ec.Settings = settings
		default:
			ec.End = p.peek(-1).End
			return ec, nil
		}
	}
}

func (p *parser) parseTTLClause(inAlter bool) ([]*TTLExpr, error) {
	if _, err := p.expect(KwTtl); err != nil {
		return nil, err
	}
	var ttls []*TTLExpr
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ttl := &TTLExpr{
			span: span{Pos: expr.StartPos(), End: expr.EndPos()},
			Expr: expr,
		}
		switch {
		case p.match(KwDelete):
			ttl.Action = "DELETE"
			if p.match(KwWhere) {
				where, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				ttl.Where = where
			}
		case p.cur().Type == KwTo && (p.peek(1).Type == KwDisk || p.peek(1).Type == KwVolume):
			p.advance() // TO
			destKind := strings.ToUpper(p.advance().Val)
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			ttl.Action = "TO " + destKind
			ttl.Destination = lit.Val
		case p.cur().Type == KwGroup && p.peek(1).Type == KwBy:
			ttl.Action = "GROUP BY"
			gbStart := p.advance().Pos // GROUP
			p.advance()                // BY
			gb, err := p.parseTTLGroupByKeys(gbStart)
			if err != nil {
				return nil, err
			}
			ttl.GroupBy = gb
			if _, err := p.expect(KwSet); err != nil {
				return nil, err
			}
			for {
				aStart := p.cur().Pos
				// ClickHouse requires a column identifier on the LHS of TTL SET assignments.
				col, err := p.parseNestedIdentifier()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(TokEq); err != nil {
					return nil, err
				}
				val, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				ttl.Set = append(ttl.Set, &AssignmentExpr{
					span:   span{Pos: aStart, End: val.EndPos()},
					Column: col,
					Value:  val,
				})
				if p.cur().Type == TokComma && p.looksLikeTTLSetAssignment(1) {
					p.advance()
					continue
				}
				break
			}
		}
		ttl.End = p.peek(-1).End
		ttls = append(ttls, ttl)
		if p.cur().Type == TokComma && (!inAlter || !isAlterClauseStart(p.peek(1).Type)) {
			p.advance()
			continue
		}
		break
	}
	return ttls, nil
}

// looksLikeTTLSetAssignment reports whether tokens starting at p.peek(offset)
// match the `<nestedIdentifier> =` prefix of a TTL SET assignment.
func (p *parser) looksLikeTTLSetAssignment(offset int) bool {
	if !isIdentifier(p.peek(offset)) {
		return false
	}
	i := offset + 1
	for p.peek(i).Type == TokDot {
		if !isIdentifier(p.peek(i + 1)) {
			return false
		}
		i += 2
	}
	return p.peek(i).Type == TokEq
}

func (p *parser) parseTTLGroupByKeys(start int) (*GroupByClause, error) {
	gb := &GroupByClause{span: span{Pos: start}}
	if p.cur().Type == KwAll && p.peek(1).Type != TokLParen {
		gb.All = true
		gb.End = p.advance().End
		return gb, nil
	}
	if (p.cur().Type == KwCube || p.cur().Type == KwRollup) && p.peek(1).Type == TokLParen {
		gb.Modifier = strings.ToUpper(p.advance().Val)
		p.advance()
		for {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			gb.Exprs = append(gb.Exprs, expr)
			if !p.match(TokComma) {
				break
			}
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		gb.End = rp.End
		return gb, nil
	}
	if p.cur().Type == KwGrouping && p.peek(1).Type == KwSets && p.peek(2).Type == TokLParen {
		gb.Modifier = "GROUPING SETS"
		p.advance()
		p.advance()
		p.advance()
		for {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			gb.Exprs = append(gb.Exprs, expr)
			if !p.match(TokComma) {
				break
			}
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		gb.End = rp.End
		return gb, nil
	}
	for {
		var expr Expr
		var err error
		if p.cur().Type == TokStar {
			expr, err = p.parseStarExpr(nil, p.cur().Pos)
		} else {
			expr, err = p.parseExpr()
		}
		if err != nil {
			return nil, err
		}
		gb.Exprs = append(gb.Exprs, expr)
		if !p.match(TokComma) {
			break
		}
	}
	gb.End = p.peek(-1).End
	return gb, nil
}

func isAlterClauseStart(tt TokenType) bool {
	switch tt {
	case KwAdd, KwAttach, KwClear, KwComment, KwDelete, KwDetach, KwDrop,
		KwFreeze, KwMaterialize, KwModify, KwMove, KwRemove, KwRename,
		KwReplace, KwUpdate:
		return true
	default:
		return false
	}
}

// --- ALTER ---

func (p *parser) parseAlterStmt() (Statement, error) {
	start := p.advance().Pos // ALTER
	switch p.cur().Type {
	case KwTable:
		p.advance()
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		var clauses []*AlterClause
		for {
			cl, err := p.parseAlterTableClause()
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, cl)
			if !p.match(TokComma) {
				break
			}
		}
		return &AlterTableStmt{
			span:     span{Pos: start, End: p.peek(-1).End},
			Database: db,
			Table:    tbl,
			Cluster:  cluster,
			Clauses:  clauses,
		}, nil
	case KwUser, KwRole, KwQuota, KwSettings, KwProfile, KwRow, KwPolicy:
		return p.parseCreateAccessStmt(start, "ALTER")
	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("expected TABLE after ALTER, got %q", p.cur().Raw))
	}
}

func (p *parser) parseAlterTableClause() (*AlterClause, error) {
	start := p.cur().Pos
	switch p.cur().Type {
	case KwAdd:
		p.advance()
		switch p.cur().Type {
		case KwColumn:
			p.advance()
			ifNotExists := false
			if p.cur().Type == KwIf && p.peek(1).Type == KwNot {
				if p.peek(2).Type != KwExists {
					return nil, p.errAt(p.peek(2), "expected EXISTS after IF NOT")
				}
				p.advance()
				p.advance()
				p.advance()
				ifNotExists = true
			}
			elem, err := p.parseTableColumnDfnt()
			if err != nil {
				return nil, err
			}
			var after string
			var first bool
			if p.match(KwAfter) {
				after, err = p.parseNestedIdentifier()
				if err != nil {
					return nil, err
				}
			} else {
				first = p.match(KwFirst)
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD COLUMN", IfNotExists: ifNotExists, Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwIndex:
			p.advance()
			ifNotExists := p.matchIfNotExists()
			elem, err := p.parseTableIndexDfnt(start)
			if err != nil {
				return nil, err
			}
			var after string
			var first bool
			if p.match(KwAfter) {
				after, err = p.parseNestedIdentifier()
				if err != nil {
					return nil, err
				}
			} else {
				first = p.match(KwFirst)
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD INDEX", IfNotExists: ifNotExists, Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwProjection:
			p.advance()
			ifNotExists := p.matchIfNotExists()
			elem, err := p.parseTableProjectionDfnt(start)
			if err != nil {
				return nil, err
			}
			var after string
			var first bool
			if p.match(KwAfter) {
				after, err = p.parseNestedIdentifier()
				if err != nil {
					return nil, err
				}
			} else {
				first = p.match(KwFirst)
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD PROJECTION", IfNotExists: ifNotExists, Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwConstraint:
			p.advance()
			ifNotExists := p.matchIfNotExists()
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(KwCheck); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			elem := &TableElement{span: span{Pos: start, End: expr.EndPos()}, Kind: "CONSTRAINT", Name: name, Expr: expr}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD CONSTRAINT", IfNotExists: ifNotExists, Name: name, Element: elem}, nil
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after ADD", p.cur().Raw))
		}

	case KwAttach:
		p.advance()
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		var target string
		if p.match(KwFrom) {
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			target = joinDB(db, tbl)
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ATTACH PARTITION", Partition: part, Target: target}, nil

	case KwClear:
		p.advance()
		var kind string
		switch p.cur().Type {
		case KwColumn, KwIndex, KwProjection:
			kind = "CLEAR " + strings.ToUpper(p.advance().Val)
		default:
			return nil, p.errAt(p.cur(), "expected COLUMN, INDEX, or PROJECTION after CLEAR")
		}
		ifExists := p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		var part *PartitionClause
		if p.match(KwIn) {
			part, err = p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, IfExists: ifExists, Name: name, Partition: part}, nil

	case KwComment:
		p.advance()
		if _, err := p.expect(KwColumn); err != nil {
			return nil, err
		}
		ifExists := p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "COMMENT COLUMN", IfExists: ifExists, Name: name, Comment: lit.Val}, nil

	case KwDelete:
		p.advance()
		if _, err := p.expect(KwWhere); err != nil {
			return nil, err
		}
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "DELETE", Where: where}, nil

	case KwDetach:
		p.advance()
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "DETACH PARTITION", Partition: part}, nil

	case KwDrop:
		p.advance()
		switch p.cur().Type {
		case KwColumn, KwIndex, KwProjection, KwConstraint:
			kind := "DROP " + strings.ToUpper(p.advance().Val)
			ifExists := p.matchIfExists()
			name, err := p.parseNestedIdentifier()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, IfExists: ifExists, Name: name}, nil
		case KwPartition:
			part, err := p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "DROP PARTITION", Partition: part}, nil
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after DROP", p.cur().Raw))
		}

	case KwFreeze:
		p.advance()
		var part *PartitionClause
		if p.cur().Type == KwPartition {
			var err error
			part, err = p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "FREEZE", Partition: part}, nil

	case KwMaterialize:
		p.advance()
		var kind string
		switch p.cur().Type {
		case KwIndex, KwProjection, KwColumn:
			kind = "MATERIALIZE " + strings.ToUpper(p.advance().Val)
		default:
			return nil, p.errAt(p.cur(), "expected INDEX or PROJECTION after MATERIALIZE")
		}
		ifExists := p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		var part *PartitionClause
		if p.match(KwIn) {
			part, err = p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, IfExists: ifExists, Name: name, Partition: part}, nil

	case KwModify:
		p.advance()
		switch p.cur().Type {
		case KwColumn:
			p.advance()
			ifExists := p.matchIfExists()
			// Could be:
			//   nestedIdentifier codecExpr
			//   nestedIdentifier COMMENT stringLiteral
			//   nestedIdentifier REMOVE tableColumnPropertyType
			//   tableColumnDfnt [FIRST | AFTER col]
			save := p.pos
			name, err := p.parseNestedIdentifier()
			if err != nil {
				return nil, err
			}
			switch p.cur().Type {
			case KwCodec:
				codec, err := p.parseCodecExpr()
				if err != nil {
					return nil, err
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN CODEC", IfExists: ifExists, Name: name, Codec: codec}, nil
			case KwComment:
				p.advance()
				lit, err := p.parseStringLiteral()
				if err != nil {
					return nil, err
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN COMMENT", IfExists: ifExists, Name: name, Comment: lit.Val}, nil
			case KwRemove:
				p.advance()
				var prop string
				switch p.cur().Type {
				case KwAlias, KwCodec, KwComment, KwDefault, KwMaterialized, KwTtl:
					prop = strings.ToUpper(p.advance().Val)
				default:
					return nil, p.errAt(p.cur(), "expected column property type after REMOVE")
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN REMOVE", IfExists: ifExists, Name: name, RemoveProp: prop}, nil
			default:
				p.pos = save
				elem, err := p.parseTableColumnDfnt()
				if err != nil {
					return nil, err
				}
				var after string
				var first bool
				if p.match(KwAfter) {
					after, err = p.parseNestedIdentifier()
					if err != nil {
						return nil, err
					}
				} else {
					first = p.match(KwFirst)
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN", IfExists: ifExists, Name: elem.Name, Element: elem, After: after, First: first}, nil
			}
		case KwOrder:
			p.advance()
			if _, err := p.expect(KwBy); err != nil {
				return nil, err
			}
			orderBy, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY ORDER BY", OrderBy: orderBy}, nil
		case KwTtl:
			ttls, err := p.parseTTLClause(true)
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY TTL", TTL: ttls}, nil
		case KwSetting:
			p.advance()
			settings, err := p.parseAlterSettingExprList()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY SETTING", Settings: settings}, nil
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after MODIFY", p.cur().Raw))
		}

	case KwMove:
		p.advance()
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwTo); err != nil {
			return nil, err
		}
		var target string
		switch p.cur().Type {
		case KwDisk, KwVolume:
			destKind := strings.ToUpper(p.advance().Val)
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			target = destKind + " " + lit.Val
		case KwTable:
			p.advance()
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			target = "TABLE " + joinDB(db, tbl)
		default:
			return nil, p.errAt(p.cur(), "expected DISK, VOLUME, or TABLE after MOVE PARTITION ... TO")
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MOVE PARTITION", Partition: part, Target: target}, nil

	case KwRemove:
		p.advance()
		if _, err := p.expect(KwTtl); err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "REMOVE TTL"}, nil

	case KwRename:
		p.advance()
		if _, err := p.expect(KwColumn); err != nil {
			return nil, err
		}
		ifExists := p.matchIfExists()
		oldName, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwTo); err != nil {
			return nil, err
		}
		newName, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "RENAME COLUMN", IfExists: ifExists, Name: oldName, NewName: newName}, nil

	case KwReplace:
		p.advance()
		part, err := p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwFrom); err != nil {
			return nil, err
		}
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "REPLACE PARTITION", Partition: part, Target: joinDB(db, tbl)}, nil

	case KwUpdate:
		p.advance()
		assigns, err := p.parseAssignmentExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwWhere); err != nil {
			return nil, err
		}
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "UPDATE", Assignments: assigns, Where: where}, nil

	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q in ALTER TABLE", p.cur().Raw))
	}
}

func (p *parser) parsePartitionClause() (*PartitionClause, error) {
	partTok, err := p.expect(KwPartition)
	if err != nil {
		return nil, err
	}
	if p.cur().Type == KwId {
		p.advance()
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		return &PartitionClause{
			span: span{Pos: partTok.Pos, End: lit.EndPos()},
			ByID: true,
			ID:   lit.Val,
			Expr: lit,
		}, nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &PartitionClause{
		span: span{Pos: partTok.Pos, End: expr.EndPos()},
		Expr: expr,
	}, nil
}

// --- DROP / DETACH / TRUNCATE / RENAME / EXCHANGE / UNDROP ---

func (p *parser) parseDropStmt() (*DropStmt, error) {
	start := p.cur().Pos
	action := strings.ToUpper(p.advance().Val) // DROP or DETACH
	stmt := &DropStmt{span: span{Pos: start}, Action: action}

	switch p.cur().Type {
	case KwDatabase:
		p.advance()
		stmt.ObjectType = "DATABASE"
		stmt.IfExists = p.matchIfExists()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
	case KwDictionary, KwView, KwTemporary, KwTable:
		var obj string
		if p.match(KwTemporary) {
			if _, err := p.expect(KwTable); err != nil {
				return nil, err
			}
			obj = "TEMPORARY TABLE"
		} else {
			obj = strings.ToUpper(p.advance().Val)
		}
		stmt.ObjectType = obj
		stmt.IfExists = p.matchIfExists()
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = db
		stmt.Name = name
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
		if p.match(KwPermanently) {
			stmt.Permanently = true
		}
		if p.cur().Type == KwNo && p.peek(1).Type == KwDelay {
			p.advance()
			p.advance()
			stmt.NoDelay = true
		} else if p.match(KwSync) {
			stmt.Sync = true
		}
	case KwFunction:
		p.advance()
		stmt.ObjectType = "FUNCTION"
		stmt.IfExists = p.matchIfExists()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
	case KwUser, KwRole, KwQuota, KwSettings, KwProfile, KwRow, KwPolicy:
		if p.match(KwSettings) {
			if _, err := p.expect(KwProfile); err != nil {
				return nil, err
			}
			stmt.ObjectType = "SETTINGS PROFILE"
		} else if p.match(KwRow) {
			if _, err := p.expect(KwPolicy); err != nil {
				return nil, err
			}
			stmt.ObjectType = "ROW POLICY"
		} else {
			stmt.ObjectType = strings.ToUpper(p.advance().Val)
		}
		stmt.IfExists = p.matchIfExists()
		names, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Names = names
		stmt.Name = names[0]
		if (stmt.ObjectType == "ROW POLICY" || stmt.ObjectType == "POLICY") && p.match(KwOn) {
			for {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Tables = append(stmt.Tables, joinDB(db, tbl))
				if !p.match(TokComma) {
					break
				}
			}
		}
		cluster, err := p.parseClusterClause()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after %s", p.cur().Raw, action))
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseTruncateStmt() (*TruncateStmt, error) {
	start := p.advance().Pos // TRUNCATE
	temp := p.match(KwTemporary)
	p.match(KwTable)
	ifExists := p.matchIfExists()
	db, tbl, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	return &TruncateStmt{
		span:      span{Pos: start, End: p.peek(-1).End},
		Temporary: temp,
		IfExists:  ifExists,
		Database:  db,
		Table:     tbl,
		Cluster:   cluster,
	}, nil
}

func (p *parser) parseRenameStmt() (*RenameStmt, error) {
	start := p.advance().Pos // RENAME
	stmt := &RenameStmt{span: span{Pos: start}, ObjectType: "TABLE"}
	switch p.cur().Type {
	case KwDatabase:
		p.advance()
		stmt.ObjectType = "DATABASE"
		from, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwTo); err != nil {
			return nil, err
		}
		to, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Pairs = [][2]string{{from, to}}
	case KwDictionary:
		p.advance()
		stmt.ObjectType = "DICTIONARY"
		for {
			fromDB, fromName, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(KwTo); err != nil {
				return nil, err
			}
			toDB, toName, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Pairs = append(stmt.Pairs, [2]string{joinDB(fromDB, fromName), joinDB(toDB, toName)})
			if !p.match(TokComma) {
				break
			}
		}
	default:
		p.match(KwTable)
		for {
			fromDB, fromName, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(KwTo); err != nil {
				return nil, err
			}
			toDB, toName, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Pairs = append(stmt.Pairs, [2]string{joinDB(fromDB, fromName), joinDB(toDB, toName)})
			if !p.match(TokComma) {
				break
			}
		}
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt.Cluster = cluster
	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseExchangeStmt() (*ExchangeStmt, error) {
	start := p.advance().Pos // EXCHANGE
	var obj string
	if p.match(KwTables) {
		obj = "TABLES"
	} else if p.match(KwDictionaries) {
		obj = "DICTIONARIES"
	} else {
		return nil, p.errAt(p.cur(), "expected TABLES or DICTIONARIES after EXCHANGE")
	}
	db1, t1, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(KwAnd); err != nil {
		return nil, err
	}
	db2, t2, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	return &ExchangeStmt{
		span:       span{Pos: start, End: p.peek(-1).End},
		ObjectType: obj,
		First:      joinDB(db1, t1),
		Second:     joinDB(db2, t2),
		Cluster:    cluster,
	}, nil
}

func (p *parser) parseUndropStmt() (*UndropStmt, error) {
	start := p.advance().Pos // UNDROP
	if _, err := p.expect(KwTable); err != nil {
		return nil, err
	}
	db, tbl, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	uuid, err := p.parseUUIDClause()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	return &UndropStmt{
		span:     span{Pos: start, End: p.peek(-1).End},
		Database: db,
		Table:    tbl,
		UUID:     uuid,
		Cluster:  cluster,
	}, nil
}

func joinDB(db, name string) string {
	if db == "" {
		return name
	}
	return db + "." + name
}
