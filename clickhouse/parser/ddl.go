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
	if action == "CREATE" && p.cur().Type == KwOr && p.peek(1).Type == KwReplace {
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
		return p.parseCreateAccessStmt(start, action, orReplace)
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

	// Optional engineClause
	if p.cur().Type == KwEngine {
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
	stmt.End = sub.End
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
	stmt.End = sub.End
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
	for {
		if err := p.parseDictionaryAttrDfnt(); err != nil {
			return nil, err
		}
		if !p.match(TokComma) {
			break
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return nil, err
	}

	// dictionaryEngineClause
	if p.cur().Type == KwPrimary && p.peek(1).Type == KwKey {
		p.advance()
		p.advance()
		for {
			if p.cur().Type == TokStar {
				if _, err := p.parseStarExpr(nil, p.cur().Pos); err != nil {
					return nil, err
				}
			} else {
				if _, err := p.parseExpr(); err != nil {
					return nil, err
				}
			}
			if !p.match(TokComma) {
				break
			}
		}
	}
	seen := map[string]bool{}
	for {
		switch p.cur().Type {
		case KwSource:
			if seen["source"] {
				return nil, p.errAt(p.cur(), "duplicate SOURCE clause")
			}
			seen["source"] = true
			if err := p.parseDictSourceOrLayout(); err != nil {
				return nil, err
			}
		case KwLifetime:
			if seen["lifetime"] {
				return nil, p.errAt(p.cur(), "duplicate LIFETIME clause")
			}
			seen["lifetime"] = true
			if err := p.parseDictLifetime(); err != nil {
				return nil, err
			}
		case KwLayout:
			if seen["layout"] {
				return nil, p.errAt(p.cur(), "duplicate LAYOUT clause")
			}
			seen["layout"] = true
			if err := p.parseDictSourceOrLayout(); err != nil {
				return nil, err
			}
		case KwRange:
			if seen["range"] {
				return nil, p.errAt(p.cur(), "duplicate RANGE clause")
			}
			seen["range"] = true
			if err := p.parseDictRange(); err != nil {
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
			if _, err := p.parseSettingExprList(); err != nil {
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
	}, nil
}

func (p *parser) parseDictionaryAttrDfnt() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.parseDataType(); err != nil {
		return err
	}
	seen := map[string]bool{}
	for {
		switch p.cur().Type {
		case KwDefault:
			if seen["default"] {
				return p.errAt(p.cur(), "duplicate DEFAULT attribute")
			}
			seen["default"] = true
			p.advance()
			if _, err := p.parseLiteral(); err != nil {
				return err
			}
		case KwExpression:
			if seen["expression"] {
				return p.errAt(p.cur(), "duplicate EXPRESSION attribute")
			}
			seen["expression"] = true
			p.advance()
			if _, err := p.parseExpr(); err != nil {
				return err
			}
		case KwHierarchical:
			if seen["hierarchical"] {
				return p.errAt(p.cur(), "duplicate HIERARCHICAL attribute")
			}
			seen["hierarchical"] = true
			p.advance()
		case KwInjective:
			if seen["injective"] {
				return p.errAt(p.cur(), "duplicate INJECTIVE attribute")
			}
			seen["injective"] = true
			p.advance()
		case KwIsObjectId:
			if seen["is_object_id"] {
				return p.errAt(p.cur(), "duplicate IS_OBJECT_ID attribute")
			}
			seen["is_object_id"] = true
			p.advance()
		default:
			return nil
		}
	}
}

func (p *parser) parseDictSourceOrLayout() error {
	p.advance() // SOURCE or LAYOUT
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	for p.cur().Type != TokRParen && p.cur().Type != TokEOF {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if isLiteralStart(p.cur(), p.peek(1)) {
			if _, err := p.parseLiteral(); err != nil {
				return err
			}
		} else {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if p.match(TokLParen) {
				if _, err := p.expect(TokRParen); err != nil {
					return err
				}
			}
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseDictLifetime() error {
	p.advance() // LIFETIME
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.match(KwMin) {
		if _, err := p.expect(TokInt); err != nil {
			return err
		}
		if _, err := p.expect(KwMax); err != nil {
			return err
		}
		if _, err := p.expect(TokInt); err != nil {
			return err
		}
	} else if p.match(KwMax) {
		if _, err := p.expect(TokInt); err != nil {
			return err
		}
		if _, err := p.expect(KwMin); err != nil {
			return err
		}
		if _, err := p.expect(TokInt); err != nil {
			return err
		}
	} else {
		if _, err := p.expect(TokInt); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseDictRange() error {
	p.advance() // RANGE
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.match(KwMin) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(KwMax); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else if p.match(KwMax) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(KwMin); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else {
		return p.errAt(p.cur(), "expected MIN or MAX in RANGE clause")
	}
	_, err := p.expect(TokRParen)
	return err
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

func (p *parser) parseCreateAccessStmt(start int, action string, orReplace bool) (*AccessControlStmt, error) {
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

	if action == "ALTER" {
		p.matchIfExists()
	} else {
		p.matchIfNotExists()
	}

	// Entity name(s) (comma-separated)
	if err := p.parseIdentList(); err != nil {
		return nil, err
	}
	// ROW POLICY requires `ON table1 [, table2 ...]`
	if entity == "ROW POLICY" || entity == "POLICY" {
		if _, err := p.expect(KwOn); err != nil {
			return nil, err
		}
		for {
			if _, _, err := p.parseTableIdentifier(); err != nil {
				return nil, err
			}
			if !p.match(TokComma) {
				break
			}
		}
	}
	if _, err := p.parseClusterClause(); err != nil {
		return nil, err
	}

	// Structured entity-specific clauses
	switch entity {
	case "USER":
		for !p.isStatementEndOrOutput() {
			switch {
			case p.match(KwIdentified):
				if p.match(KwWith) {
					if _, err := p.parseIdentifier(); err != nil {
						return nil, err
					}
				}
				if p.match(KwBy) {
					if _, err := p.parseStringLiteral(); err != nil {
						return nil, err
					}
				}
			case p.cur().Type == KwDefault && p.peek(1).Type == KwRole:
				p.advance()
				p.advance()
				if p.isWord("NONE") || p.match(KwAll) {
					if p.isWord("NONE") {
						p.advance()
					}
				} else if err := p.parseIdentList(); err != nil {
					return nil, err
				}
			case p.cur().Type == KwSettings:
				if _, err := p.parseSettingsClause(); err != nil {
					return nil, err
				}
			case p.isWord("HOST"):
				p.advance()
				if p.cur().Type == KwAny || p.isWord("NONE") || p.isWord("LOCAL") {
					p.advance()
				} else {
					for {
						if _, err := p.parseIdentifier(); err != nil {
							return nil, err
						}
						if _, err := p.parseStringLiteral(); err != nil {
							return nil, err
						}
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
			if _, err := p.parseSettingsClause(); err != nil {
				return nil, err
			}
		}
		if p.match(KwTo) {
			if err := p.parseGranteeList(); err != nil {
				return nil, err
			}
		}
	case "ROW POLICY", "POLICY":
		for !p.isStatementEndOrOutput() {
			switch {
			case p.match(KwFor):
				if _, err := p.expect(KwSelect); err != nil {
					return nil, err
				}
			case p.match(KwUsing):
				if _, err := p.parseExpr(); err != nil {
					return nil, err
				}
			case p.match(KwTo):
				if err := p.parseGranteeList(); err != nil {
					return nil, err
				}
			default:
				return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q in %s", p.cur().Raw, kind))
			}
		}
	case "QUOTA":
		if p.match(KwTo) {
			if err := p.parseGranteeList(); err != nil {
				return nil, err
			}
		}
	}

	return &AccessControlStmt{
		span: span{Pos: start, End: p.peek(-1).End},
		Kind: kind,
	}, nil
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
		if isIdentifier(p.peek(1)) && p.peek(2).Type == TokLParen {
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
	if _, err := p.expect(KwGranularity); err != nil {
		return nil, err
	}
	gran, err := p.expect(TokInt)
	if err != nil {
		return nil, err
	}
	return &TableElement{
		span:        span{Pos: start, End: gran.End},
		Kind:        "INDEX",
		Name:        name,
		Type:        dt,
		Expr:        expr,
		Granularity: gran.Raw,
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
		gb, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		proj.GroupBy = gb
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

func (p *parser) parseEngineClause() (*EngineClause, error) {
	eng, err := p.parseEngineExpr()
	if err != nil {
		return nil, err
	}
	ec := &EngineClause{span: eng.span, Engine: eng}
	seen := map[string]bool{}
	for {
		switch {
		case p.cur().Type == KwOrder && p.peek(1).Type == KwBy:
			if seen["order"] {
				return nil, p.errAt(p.cur(), "duplicate ORDER BY clause")
			}
			seen["order"] = true
			p.advance()
			p.advance()
			items, err := p.parseOrderExprList()
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
			ttls, err := p.parseTTLClause()
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

func (p *parser) parseTTLClause() ([]Expr, error) {
	if _, err := p.expect(KwTtl); err != nil {
		return nil, err
	}
	var exprs []Expr
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		switch {
		case p.match(KwDelete):
			if p.match(KwWhere) {
				if _, err := p.parseExpr(); err != nil {
					return nil, err
				}
			}
		case p.cur().Type == KwTo && (p.peek(1).Type == KwDisk || p.peek(1).Type == KwVolume):
			p.advance()
			p.advance()
			if _, err := p.parseStringLiteral(); err != nil {
				return nil, err
			}
		case p.cur().Type == KwGroup && p.peek(1).Type == KwBy:
			p.advance() // GROUP
			p.advance() // BY
			if err := p.parseTTLGroupByKeys(); err != nil {
				return nil, err
			}
			if _, err := p.expect(KwSet); err != nil {
				return nil, err
			}
			for {
				// ClickHouse requires a column identifier on the LHS of TTL SET assignments.
				if _, err := p.parseNestedIdentifier(); err != nil {
					return nil, err
				}
				if _, err := p.expect(TokEq); err != nil {
					return nil, err
				}
				if _, err := p.parseExpr(); err != nil {
					return nil, err
				}
				if !p.match(TokComma) {
					break
				}
			}
		}
		exprs = append(exprs, expr)
		if p.cur().Type == TokComma && !isAlterClauseStart(p.peek(1).Type) {
			p.advance()
			continue
		}
		break
	}
	return exprs, nil
}

func (p *parser) parseTTLGroupByKeys() error {
	if p.cur().Type == KwAll && p.peek(1).Type != TokLParen {
		p.advance()
		return nil
	}
	if (p.cur().Type == KwCube || p.cur().Type == KwRollup) && p.peek(1).Type == TokLParen {
		p.advance()
		p.advance()
		for {
			if _, err := p.parseExpr(); err != nil {
				return err
			}
			if !p.match(TokComma) {
				break
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if p.cur().Type == KwGrouping && p.peek(1).Type == KwSets && p.peek(2).Type == TokLParen {
		p.advance()
		p.advance()
		p.advance()
		for {
			if _, err := p.parseExpr(); err != nil {
				return err
			}
			if !p.match(TokComma) {
				break
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}
	for {
		if p.cur().Type == TokStar {
			if _, err := p.parseStarExpr(nil, p.cur().Pos); err != nil {
				return err
			}
		} else if _, err := p.parseExpr(); err != nil {
			return err
		}
		if !p.match(TokComma) {
			break
		}
	}
	return nil
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
		return p.parseCreateAccessStmt(start, "ALTER", false)
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
			if p.cur().Type == KwIf && p.peek(1).Type == KwNot {
				if p.peek(2).Type != KwExists {
					return nil, p.errAt(p.peek(2), "expected EXISTS after IF NOT")
				}
				p.advance()
				p.advance()
				p.advance()
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
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD COLUMN", Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwIndex:
			p.advance()
			p.matchIfNotExists()
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
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD INDEX", Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwProjection:
			p.advance()
			p.matchIfNotExists()
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
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD PROJECTION", Name: elem.Name, Element: elem, After: after, First: first}, nil
		case KwConstraint:
			p.advance()
			p.matchIfNotExists()
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
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "ADD CONSTRAINT", Name: name, Element: elem}, nil
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
		p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		var part Expr
		if p.match(KwIn) {
			part, err = p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, Name: name, Partition: part}, nil

	case KwComment:
		p.advance()
		if _, err := p.expect(KwColumn); err != nil {
			return nil, err
		}
		p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "COMMENT COLUMN", Name: name, Comment: lit.Val}, nil

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
			p.matchIfExists()
			name, err := p.parseNestedIdentifier()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, Name: name}, nil
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
		var part Expr
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
		p.matchIfExists()
		name, err := p.parseNestedIdentifier()
		if err != nil {
			return nil, err
		}
		var part Expr
		if p.match(KwIn) {
			part, err = p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
		}
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: kind, Name: name, Partition: part}, nil

	case KwModify:
		p.advance()
		switch p.cur().Type {
		case KwColumn:
			p.advance()
			p.matchIfExists()
			// Could be:
			//   nestedIdentifier codecExpr
			//   nestedIdentifier COMMENT stringLiteral
			//   nestedIdentifier REMOVE tableColumnPropertyType
			//   tableColumnDfnt
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
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN CODEC", Name: name, Codec: codec}, nil
			case KwComment:
				p.advance()
				lit, err := p.parseStringLiteral()
				if err != nil {
					return nil, err
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN COMMENT", Name: name, Comment: lit.Val}, nil
			case KwRemove:
				p.advance()
				var prop string
				switch p.cur().Type {
				case KwAlias, KwCodec, KwComment, KwDefault, KwMaterialized, KwTtl:
					prop = strings.ToUpper(p.advance().Val)
				default:
					return nil, p.errAt(p.cur(), "expected column property type after REMOVE")
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN REMOVE", Name: name, RemoveProp: prop}, nil
			default:
				p.pos = save
				elem, err := p.parseTableColumnDfnt()
				if err != nil {
					return nil, err
				}
				return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY COLUMN", Name: elem.Name, Element: elem}, nil
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
			ttls, err := p.parseTTLClause()
			if err != nil {
				return nil, err
			}
			return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "MODIFY TTL", TTL: ttls}, nil
		case KwSetting:
			p.advance()
			settings, err := p.parseSettingExprList()
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
		p.matchIfExists()
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
		return &AlterClause{span: span{Pos: start, End: p.peek(-1).End}, Action: "RENAME COLUMN", Name: oldName, NewName: newName}, nil

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

func (p *parser) parsePartitionClause() (Expr, error) {
	if _, err := p.expect(KwPartition); err != nil {
		return nil, err
	}
	if p.cur().Type == KwId {
		p.advance()
		return p.parseStringLiteral()
	}
	return p.parseExpr()
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
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		for p.match(TokComma) {
			if _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
		if (stmt.ObjectType == "ROW POLICY" || stmt.ObjectType == "POLICY") && p.match(KwOn) {
			for {
				if _, _, err := p.parseTableIdentifier(); err != nil {
					return nil, err
				}
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
