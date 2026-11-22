package parser

import (
	"fmt"
	"strings"
)

func (p *parser) parseDescribeStmt() (*DescribeStmt, error) {
	start := p.advance().Pos // DESCRIBE or DESC
	p.match(KwTable)
	te, err := p.parseDescribeTableExpr()
	if err != nil {
		return nil, err
	}
	return &DescribeStmt{
		span:   span{Pos: start, End: te.End},
		Target: te,
	}, nil
}

func (p *parser) parseDescribeTableExpr() (*TableExpr, error) {
	start := p.cur().Pos
	te := &TableExpr{span: span{Pos: start}}
	switch {
	case p.cur().Type == TokLParen:
		p.advance()
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		te.Subquery = sub
		te.End = rp.End
	case isIdentifier(p.cur()) && p.peek(1).Type == TokLParen:
		fn, err := p.parseTableFunctionExpr()
		if err != nil {
			return nil, err
		}
		te.Func = fn
		te.End = fn.End
	case p.cur().Type == KwFrom && !p.isEndAfter(1):
		// `DESC FROM tbl WHERE ...` is ClickHouse's DESCRIBE of a FROM-first SELECT
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		te.Subquery = sub
		te.End = sub.End
		return te, nil
	case p.cur().Type == TokLBrace:
		firstParam, err := p.parseQueryParamNode()
		if err != nil {
			return nil, err
		}
		if p.match(TokDot) {
			te.DatabaseParam = firstParam
			if p.cur().Type == TokLBrace {
				tblParam, err := p.parseQueryParamNode()
				if err != nil {
					return nil, err
				}
				te.TableParam = tblParam
			} else {
				second, err := p.parseAnyWordIdent()
				if err != nil {
					return nil, err
				}
				te.Table = second
			}
		} else {
			te.TableParam = firstParam
		}
		te.End = p.peek(-1).End
	default:
		// In DESCRIBE, ClickHouse allows INF/NAN/NULL as table identifiers too.
		first, err := p.parseAnyWordIdent()
		if err != nil {
			return nil, err
		}
		if p.match(TokDot) {
			te.Database = first
			if p.cur().Type == TokLBrace {
				tblParam, err := p.parseQueryParamNode()
				if err != nil {
					return nil, err
				}
				te.TableParam = tblParam
			} else {
				second, err := p.parseAnyWordIdent()
				if err != nil {
					return nil, err
				}
				te.Table = second
			}
		} else {
			te.Table = first
		}
		te.End = p.peek(-1).End
	}

	// Optional alias: AS ident | bare word alias (ClickHouse allows almost any
	// keyword as a table alias in DESCRIBE except clause terminators and UUID).
	if p.match(KwAs) {
		alias, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		te.Alias = alias
		te.End = p.peek(-1).End
	} else if isDescribeTableAlias(p.cur()) {
		tok := p.advance()
		te.Alias = tok.Val
		te.End = tok.End
	}
	return te, nil
}

func (p *parser) parseAnyWordIdent() (string, error) {
	tok := p.cur()
	if tok.Type == TokIdent || tok.Type == TokQuotedIdent || tok.IsKeyword() {
		p.advance()
		return tok.Val, nil
	}
	return "", p.errAt(tok, fmt.Sprintf("expected identifier, got %q", tok.Raw))
}

func (p *parser) isEndAfter(offset int) bool {
	t := p.peek(offset).Type
	return t == TokEOF || t == TokSemicolon || t == KwFormat || (t == KwInto && p.peek(offset+1).Type == KwOutfile)
}

func isDescribeTableAlias(tok Token) bool {
	if tok.Type == TokIdent || tok.Type == TokQuotedIdent {
		return true
	}
	if !tok.IsKeyword() {
		return false
	}
	switch tok.Type {
	case KwAll, KwAnti, KwAny, KwArray, KwAs, KwAsof, KwBetween, KwCross,
		KwElse, KwEnd, KwExcept, KwFinal, KwFormat, KwFrom, KwFull, KwGlobal,
		KwGroup, KwHaving, KwIlike, KwIn, KwInner, KwIntersect, KwInto, KwIs,
		KwJoin, KwLeft, KwLike, KwLimit, KwNot, KwOffset, KwOn, KwOrder,
		KwOuter, KwPrewhere, KwQualify, KwRight, KwSample, KwSemi, KwSettings,
		KwUnion, KwUsing, KwUuid, KwWhen, KwWhere, KwWindow, KwWith:
		return false
	default:
		return true
	}
}

func (p *parser) parseExistsStmt() (*ExistsStmt, error) {
	start := p.advance().Pos // EXISTS
	stmt := &ExistsStmt{span: span{Pos: start}, ObjectType: "TABLE"}
	switch p.cur().Type {
	case KwDatabase:
		p.advance()
		stmt.ObjectType = "DATABASE"
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	case KwDictionary, KwView:
		stmt.ObjectType = strings.ToUpper(p.advance().Val)
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = db
		stmt.Name = name
	case KwTemporary:
		p.advance()
		if _, err := p.expect(KwTable); err != nil {
			return nil, err
		}
		stmt.ObjectType = "TEMPORARY TABLE"
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = db
		stmt.Name = name
	default:
		p.match(KwTable)
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Database = db
		stmt.Name = name
	}
	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseCheckStmt() (*CheckStmt, error) {
	start := p.advance().Pos // CHECK
	if _, err := p.expect(KwTable); err != nil {
		return nil, err
	}
	db, tbl, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	var part *PartitionClause
	if p.cur().Type == KwPartition {
		part, err = p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
	}
	return &CheckStmt{
		span:      span{Pos: start, End: p.peek(-1).End},
		Database:  db,
		Table:     tbl,
		Partition: part,
	}, nil
}

func (p *parser) parseOptimizeStmt() (*OptimizeStmt, error) {
	start := p.advance().Pos // OPTIMIZE
	if _, err := p.expect(KwTable); err != nil {
		return nil, err
	}
	db, tbl, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	var part *PartitionClause
	if p.cur().Type == KwPartition {
		part, err = p.parsePartitionClause()
		if err != nil {
			return nil, err
		}
	}
	final := p.match(KwFinal)
	dedup := p.match(KwDeduplicate)
	var dedupBy []Expr
	if dedup && p.match(KwBy) {
		dedupBy, err = p.parseColumnExprList()
		if err != nil {
			return nil, err
		}
	}
	return &OptimizeStmt{
		span:          span{Pos: start, End: p.peek(-1).End},
		Database:      db,
		Table:         tbl,
		Cluster:       cluster,
		Partition:     part,
		Final:         final,
		Deduplicate:   dedup,
		DeduplicateBy: dedupBy,
	}, nil
}

func (p *parser) parseKillStmt() (*KillStmt, error) {
	start := p.advance().Pos // KILL
	var kind string
	if p.match(KwMutation) {
		kind = "MUTATION"
	} else if p.match(KwQuery) {
		kind = "QUERY"
	} else {
		return nil, p.errAt(p.cur(), "expected MUTATION or QUERY after KILL")
	}
	cluster, err := p.parseClusterClause()
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
	mode := ""
	switch p.cur().Type {
	case KwSync, KwAsync, KwTest:
		mode = strings.ToUpper(p.advance().Val)
	}
	return &KillStmt{
		span:    span{Pos: start, End: p.peek(-1).End},
		Kind:    kind,
		Cluster: cluster,
		Where:   where,
		Mode:    mode,
	}, nil
}

func (p *parser) parseSetStmt() (*SetStmt, error) {
	start := p.advance().Pos // SET
	settings, err := p.parseSettingExprList()
	if err != nil {
		return nil, err
	}
	return &SetStmt{
		span:     span{Pos: start, End: p.peek(-1).End},
		Settings: settings,
	}, nil
}

func (p *parser) parseUseStmt() (*UseStmt, error) {
	start := p.advance().Pos // USE
	db, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &UseStmt{
		span:     span{Pos: start, End: p.peek(-1).End},
		Database: db,
	}, nil
}

func (p *parser) parseWatchStmt() (*WatchStmt, error) {
	start := p.advance().Pos // WATCH
	db, tbl, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	events := p.match(KwEvents)
	limit := ""
	if p.match(KwLimit) {
		tok, err := p.expect(TokInt)
		if err != nil {
			return nil, err
		}
		limit = tok.Raw
	}
	return &WatchStmt{
		span:     span{Pos: start, End: p.peek(-1).End},
		Database: db,
		Table:    tbl,
		Events:   events,
		Limit:    limit,
	}, nil
}

func (p *parser) parseExplainStmt() (*ExplainStmt, error) {
	start := p.advance().Pos // EXPLAIN
	stmt := &ExplainStmt{span: span{Pos: start}}

	switch p.cur().Type {
	case KwAst, KwSyntax, KwPlan, KwPipeline, KwEstimate:
		stmt.Kind = strings.ToUpper(p.advance().Val)
	case KwQuery:
		if p.peek(1).Type == KwTree {
			p.advance()
			p.advance()
			stmt.Kind = "QUERY TREE"
		}
	case KwTable:
		if p.peek(1).Type == KwOverride {
			p.advance()
			p.advance()
			stmt.Kind = "TABLE OVERRIDE"
		}
	case KwCurrent:
		if p.peek(1).Type == TokIdent && strings.EqualFold(p.peek(1).Val, "TRANSACTION") {
			p.advance()
			p.advance()
			stmt.Kind = "CURRENT TRANSACTION"
			stmt.End = p.peek(-1).End
			return stmt, nil
		}
	default:
		if p.cur().Type == TokIdent && strings.EqualFold(p.cur().Val, "ANALYZE") && p.peek(1).Type != TokEq {
			p.advance()
			stmt.Kind = "ANALYZE"
		}
	}

	// Optional settings: ident = literal (, ident = literal)*
	if isIdentifier(p.cur()) && p.peek(1).Type == TokEq {
		settings, err := p.parseSettingExprList()
		if err != nil {
			return nil, err
		}
		stmt.Settings = settings
	}

	// Inner query
	var inner Statement
	var err error
	if stmt.Kind == "AST" {
		inner, err = p.parseTopLevel()
	} else if stmt.Kind == "TABLE OVERRIDE" {
		// EXPLAIN TABLE OVERRIDE tableIdentifier (optional override clauses)
		var db, tbl string
		db, tbl, err = p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Table = joinDB(db, tbl)
		if !p.isStatementEndOrOutput() {
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q in EXPLAIN TABLE OVERRIDE", p.cur().Raw))
		}
	} else {
		inner, err = p.parseSelectUnionStmt()
	}
	if err != nil {
		return nil, err
	}
	stmt.Inner = inner
	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseSystemStmt() (*SystemStmt, error) {
	start := p.advance().Pos // SYSTEM
	stmt := &SystemStmt{span: span{Pos: start}}
	switch p.cur().Type {
	case KwFlush:
		p.advance()
		if p.match(KwDistributed) {
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Action = "FLUSH DISTRIBUTED"
			stmt.Target = joinDB(db, tbl)
		} else if p.match(KwLogs) {
			stmt.Action = "FLUSH LOGS"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				name, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = name
			}
		} else {
			return nil, p.errAt(p.cur(), "expected DISTRIBUTED or LOGS after SYSTEM FLUSH")
		}
	case KwReload:
		p.advance()
		if p.match(KwDictionaries) {
			stmt.Action = "RELOAD DICTIONARIES"
		} else if p.match(KwDictionary) {
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Action = "RELOAD DICTIONARY"
			stmt.Target = joinDB(db, tbl)
		} else {
			return nil, p.errAt(p.cur(), "expected DICTIONARIES or DICTIONARY after SYSTEM RELOAD")
		}
	case KwStart, KwStop:
		verb := strings.ToUpper(p.advance().Val)
		switch {
		case p.cur().Type == KwDistributed && p.peek(1).Type == KwSends:
			p.advance()
			p.advance()
			stmt.Action = verb + " DISTRIBUTED SENDS"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = joinDB(db, tbl)
			}
		case p.cur().Type == KwFetches:
			p.advance()
			stmt.Action = verb + " FETCHES"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = joinDB(db, tbl)
			}
		case p.cur().Type == KwTtl && p.peek(1).Type == KwMerges:
			p.advance()
			p.advance()
			stmt.Action = verb + " TTL MERGES"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = joinDB(db, tbl)
			}
		case p.cur().Type == KwMerges:
			p.advance()
			stmt.Action = verb + " MERGES"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = joinDB(db, tbl)
			}
		case p.cur().Type == KwReplicated && p.peek(1).Type == KwSends:
			p.advance()
			p.advance()
			stmt.Action = verb + " REPLICATED SENDS"
			if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = joinDB(db, tbl)
			}
		case isIdentifier(p.cur()):
			// ClickHouse's ParserSystemQuery falls back to SYSTEM START|STOP [db.]table
			// (e.g. starting/stopping a view or table target by identifier).
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Action = verb
			stmt.Target = joinDB(db, tbl)
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after SYSTEM %s", p.cur().Raw, verb))
		}
	case KwSync:
		p.advance()
		if _, err := p.expect(KwReplica); err != nil {
			return nil, err
		}
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Action = "SYNC REPLICA"
		stmt.Target = joinDB(db, tbl)
	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after SYSTEM", p.cur().Raw))
	}
	stmt.End = p.peek(-1).End
	return stmt, nil
}

// --- SHOW statements ---

func (p *parser) parseShowStmt() (*ShowStmt, error) {
	start := p.advance().Pos // SHOW
	stmt := &ShowStmt{span: span{Pos: start}}

	switch p.cur().Type {
	case KwCreate:
		p.advance()
		switch p.cur().Type {
		case KwDatabase:
			p.advance()
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Kind, stmt.Target = "CREATE DATABASE", name
		case KwDictionary:
			p.advance()
			db, name, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Kind, stmt.Target = "CREATE DICTIONARY", joinDB(db, name)
		case KwView:
			p.advance()
			db, name, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Kind, stmt.Target = "CREATE VIEW", joinDB(db, name)
		case KwUser:
			p.advance()
			stmt.Kind = "CREATE USER"
			if p.match(KwCurrentUser) {
				stmt.Target = "CURRENT_USER"
			} else if isIdentifier(p.cur()) && !p.isStatementEndOrOutput() {
				names, err := p.parseIdentList()
				if err != nil {
					return nil, err
				}
				stmt.Targets = names
			}
		case KwRole:
			p.advance()
			if p.isStatementEndOrOutput() {
				// `SHOW CREATE ROLE` is parsed as `SHOW CREATE <table=ROLE>` in ClickHouse.
				stmt.Kind, stmt.Target = "CREATE TABLE", "ROLE"
			} else {
				stmt.Kind = "CREATE ROLE"
				names, err := p.parseIdentList()
				if err != nil {
					return nil, err
				}
				stmt.Targets = names
			}
		case KwRow, KwPolicy:
			p.match(KwRow)
			if _, err := p.expect(KwPolicy); err != nil {
				return nil, err
			}
			stmt.Kind = "CREATE ROW POLICY"
			if p.cur().Type != KwOn {
				name, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Target = name
			}
			if _, err := p.expect(KwOn); err != nil {
				return nil, err
			}
			for {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.Targets = append(stmt.Targets, joinDB(db, tbl))
				if !p.match(TokComma) {
					break
				}
			}
		case KwQuota:
			p.advance()
			if p.isStatementEndOrOutput() {
				stmt.Kind, stmt.Target = "CREATE TABLE", "QUOTA"
			} else {
				stmt.Kind = "CREATE QUOTA"
				if p.match(KwCurrent) {
					stmt.Target = "CURRENT"
				} else {
					names, err := p.parseIdentList()
					if err != nil {
						return nil, err
					}
					stmt.Targets = names
				}
			}
		case KwSettings, KwProfile:
			hasSettings := p.match(KwSettings)
			if _, err := p.expect(KwProfile); err != nil {
				return nil, err
			}
			if !hasSettings && p.isStatementEndOrOutput() {
				stmt.Kind, stmt.Target = "CREATE TABLE", "PROFILE"
			} else {
				stmt.Kind = "CREATE SETTINGS PROFILE"
				names, err := p.parseIdentList()
				if err != nil {
					return nil, err
				}
				stmt.Targets = names
			}
		default:
			stmt.Temporary = p.match(KwTemporary)
			p.match(KwTable)
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Kind, stmt.Target = "CREATE TABLE", joinDB(db, tbl)
		}

	case KwDatabase:
		p.advance()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "DATABASE", name

	case KwDictionary:
		p.advance()
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "DICTIONARY", joinDB(db, name)

	case KwView:
		p.advance()
		db, name, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "VIEW", joinDB(db, name)

	case KwDatabases:
		p.advance()
		stmt.Kind = "DATABASES"
		if err := p.parseShowLikeAndLimit(stmt); err != nil {
			return nil, err
		}

	case KwDictionaries:
		p.advance()
		stmt.Kind = "DICTIONARIES"
		if p.match(KwFrom) {
			from, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.From = from
		}
		if err := p.parseShowFilterAndLimit(stmt); err != nil {
			return nil, err
		}

	case KwFull:
		p.advance()
		stmt.Full = true
		switch p.cur().Type {
		case KwTemporary, KwTables:
			stmt.Temporary = p.match(KwTemporary)
			if _, err := p.expect(KwTables); err != nil {
				return nil, err
			}
			stmt.Kind = "TABLES"
			if err := p.parseShowTablesTail(stmt); err != nil {
				return nil, err
			}
		case KwColumns, KwFields:
			p.advance()
			stmt.Kind = "COLUMNS"
			if err := p.parseShowColumnsTail(stmt); err != nil {
				return nil, err
			}
		default:
			return nil, p.errAt(p.cur(), "expected TABLES, COLUMNS, or FIELDS after SHOW FULL")
		}

	case KwTemporary, KwTables:
		stmt.Temporary = p.match(KwTemporary)
		if _, err := p.expect(KwTables); err != nil {
			return nil, err
		}
		stmt.Kind = "TABLES"
		if err := p.parseShowTablesTail(stmt); err != nil {
			return nil, err
		}

	case KwExtended:
		p.advance()
		stmt.Extended = true
		if p.match(KwFull) {
			stmt.Full = true
			if p.cur().Type != KwColumns && p.cur().Type != KwFields {
				return nil, p.errAt(p.cur(), "expected COLUMNS or FIELDS")
			}
			p.advance()
			stmt.Kind = "COLUMNS"
			if err := p.parseShowColumnsTail(stmt); err != nil {
				return nil, err
			}
		} else if p.cur().Type == KwColumns || p.cur().Type == KwFields {
			p.advance()
			stmt.Kind = "COLUMNS"
			if err := p.parseShowColumnsTail(stmt); err != nil {
				return nil, err
			}
		} else if p.cur().Type == KwIndex || p.cur().Type == KwIndexes || p.cur().Type == KwIndices || p.cur().Type == KwKeys {
			p.advance()
			stmt.Kind = "INDEXES"
			if err := p.parseShowIndexTail(stmt); err != nil {
				return nil, err
			}
		} else {
			return nil, p.errAt(p.cur(), "expected COLUMNS, FIELDS, INDEX, INDEXES, INDICES, or KEYS after SHOW EXTENDED")
		}

	case KwColumns, KwFields:
		p.advance()
		stmt.Kind = "COLUMNS"
		if err := p.parseShowColumnsTail(stmt); err != nil {
			return nil, err
		}

	case KwIndex, KwIndexes, KwIndices, KwKeys:
		p.advance()
		stmt.Kind = "INDEXES"
		if err := p.parseShowIndexTail(stmt); err != nil {
			return nil, err
		}

	case KwProcesslist:
		p.advance()
		stmt.Kind = "PROCESSLIST"

	case KwGrants:
		p.advance()
		stmt.Kind = "GRANTS"
		if p.match(KwFor) {
			names, err := p.parseIdentList()
			if err != nil {
				return nil, err
			}
			stmt.Targets = names
		}
		if p.cur().Type == KwWith && p.peek(1).Type == KwImplicit {
			p.advance()
			p.advance()
			stmt.WithImplicit = true
		}
		stmt.Final = p.match(KwFinal)

	case KwUsers:
		p.advance()
		stmt.Kind = "USERS"

	case KwCurrent:
		p.advance()
		if p.match(KwRoles) {
			stmt.Kind = "CURRENT ROLES"
		} else if p.match(KwQuota) {
			stmt.Kind = "CURRENT QUOTA"
		} else {
			return nil, p.errAt(p.cur(), "expected ROLES or QUOTA after SHOW CURRENT")
		}

	case KwEnabled:
		p.advance()
		if _, err := p.expect(KwRoles); err != nil {
			return nil, err
		}
		stmt.Kind = "ENABLED ROLES"

	case KwRoles:
		p.advance()
		stmt.Kind = "ROLES"

	case KwProfiles:
		p.advance()
		stmt.Kind = "PROFILES"

	case KwRow:
		p.advance()
		if _, err := p.expect(KwPolicies); err != nil {
			return nil, err
		}
		stmt.Kind = "ROW POLICIES"
		if err := p.parseShowPoliciesTarget(stmt); err != nil {
			return nil, err
		}

	case KwPolicies:
		p.advance()
		stmt.Kind = "POLICIES"
		if err := p.parseShowPoliciesTarget(stmt); err != nil {
			return nil, err
		}

	case KwQuotas:
		p.advance()
		stmt.Kind = "QUOTAS"

	case KwQuota:
		p.advance()
		stmt.Kind = "QUOTA"

	case KwAccess:
		p.advance()
		stmt.Kind = "ACCESS"

	case KwCluster:
		p.advance()
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "CLUSTER", lit.Val

	case KwClusters:
		p.advance()
		stmt.Kind = "CLUSTERS"
		if err := p.parseShowLikeAndLimit(stmt); err != nil {
			return nil, err
		}

	case KwChanged:
		p.advance()
		if _, err := p.expect(KwSettings); err != nil {
			return nil, err
		}
		stmt.Kind = "CHANGED SETTINGS"
		if p.cur().Type != KwLike && p.cur().Type != KwIlike {
			return nil, p.errAt(p.cur(), "expected LIKE or ILIKE after SHOW CHANGED SETTINGS")
		}
		stmt.ILike = p.advance().Type == KwIlike
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Like = lit.Val

	case KwSettings:
		p.advance()
		if p.match(KwProfiles) {
			stmt.Kind = "SETTINGS PROFILES"
		} else {
			stmt.Kind = "SETTINGS"
			if p.cur().Type != KwLike && p.cur().Type != KwIlike {
				return nil, p.errAt(p.cur(), "expected PROFILES, LIKE, or ILIKE after SHOW SETTINGS")
			}
			stmt.ILike = p.advance().Type == KwIlike
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			stmt.Like = lit.Val
		}

	case KwSetting:
		p.advance()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "SETTING", name

	case KwFilesystem:
		p.advance()
		if _, err := p.expect(KwCaches); err != nil {
			return nil, err
		}
		stmt.Kind = "FILESYSTEM CACHES"

	case KwEngines:
		p.advance()
		stmt.Kind = "ENGINES"

	case KwFunctions:
		p.advance()
		stmt.Kind = "FUNCTIONS"
		if p.cur().Type == KwLike || p.cur().Type == KwIlike {
			stmt.ILike = p.advance().Type == KwIlike
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			stmt.Like = lit.Val
		}

	case KwMerges:
		p.advance()
		stmt.Kind = "MERGES"
		if err := p.parseShowLikeAndLimit(stmt); err != nil {
			return nil, err
		}

	case KwPrivileges:
		p.advance()
		stmt.Kind = "PRIVILEGES"

	case KwTable:
		p.advance()
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind, stmt.Target = "TABLE", joinDB(db, tbl)

	default:
		return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after SHOW", p.cur().Raw))
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

// parseShowPoliciesTarget parses the optional target of `SHOW [ROW] POLICIES`:
//   - `ON [db.]table` (or `ON *`, `ON *.*`, `ON db.*`)
//   - a single identifier or string literal policy short name (ClickHouse's
//     ParserShowAccessEntitiesQuery accepts `SHOW [ROW] POLICIES <short_name>`
//     as well as `ON <target>`, treating a bare trailing `ON` as `<short_name>`).
func (p *parser) parseShowPoliciesTarget(stmt *ShowStmt) error {
	if p.isStatementEndOrOutput() {
		return nil
	}
	if p.cur().Type == KwOn {
		p.advance() // ON
		if p.cur().Type == TokEOF || p.cur().Type == TokSemicolon {
			stmt.Target = p.peek(-1).Val
			return nil
		}
		target, err := p.parsePrivilegeTarget()
		if err != nil {
			return err
		}
		stmt.Target = target
		return nil
	}
	if isStringTok(p.cur().Type) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return err
		}
		stmt.Target = lit.Val
		return nil
	}
	if p.cur().Type == TokIdent || p.cur().Type == TokQuotedIdent || p.cur().IsKeyword() {
		stmt.Target = p.advance().Val
		return nil
	}
	return nil
}

func (p *parser) parseIdentList() ([]string, error) {
	var names []string
	for {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if !p.match(TokComma) {
			break
		}
	}
	return names, nil
}

func (p *parser) parseShowTablesTail(stmt *ShowStmt) error {
	if p.match(KwFrom) || p.match(KwIn) {
		from, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		stmt.From = from
	}
	return p.parseShowFilterAndLimit(stmt)
}

func (p *parser) parseShowColumnsTail(stmt *ShowStmt) error {
	if !p.match(KwFrom) && !p.match(KwIn) {
		return p.errAt(p.cur(), "expected FROM or IN")
	}
	target, err := p.parseShowTableOrInDB()
	if err != nil {
		return err
	}
	stmt.From = target
	return p.parseShowFilterAndLimit(stmt)
}

func (p *parser) parseShowIndexTail(stmt *ShowStmt) error {
	if !p.match(KwFrom) && !p.match(KwIn) {
		return p.errAt(p.cur(), "expected FROM or IN")
	}
	target, err := p.parseShowTableOrInDB()
	if err != nil {
		return err
	}
	stmt.From = target
	if p.match(KwWhere) {
		where, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Where = where
	}
	return nil
}

func (p *parser) parseShowTableOrInDB() (string, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.match(TokDot) {
		second, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		return first + "." + second, nil
	}
	if p.match(KwFrom) || p.match(KwIn) {
		db, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		return db + "." + first, nil
	}
	return first, nil
}

func (p *parser) parseShowFilterAndLimit(stmt *ShowStmt) error {
	if p.match(KwWhere) {
		where, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Where = where
	} else {
		not := p.match(KwNot)
		if p.cur().Type == KwLike || p.cur().Type == KwIlike {
			stmt.NotLike = not
			stmt.ILike = p.advance().Type == KwIlike
			lit, err := p.parseStringLiteral()
			if err != nil {
				return err
			}
			stmt.Like = lit.Val
		} else if not {
			return p.errAt(p.cur(), "expected LIKE or ILIKE after NOT")
		}
	}
	if p.match(KwLimit) {
		limit, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Limit = limit
	}
	return nil
}

func (p *parser) parseShowLikeAndLimit(stmt *ShowStmt) error {
	not := p.match(KwNot)
	if p.cur().Type == KwLike || p.cur().Type == KwIlike {
		stmt.NotLike = not
		stmt.ILike = p.advance().Type == KwIlike
		lit, err := p.parseStringLiteral()
		if err != nil {
			return err
		}
		stmt.Like = lit.Val
	} else if not {
		return p.errAt(p.cur(), "expected LIKE or ILIKE after NOT")
	}
	if p.match(KwLimit) {
		limit, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Limit = limit
	}
	return nil
}

func (p *parser) parseGrantRevokeStmt() (*AccessControlStmt, error) {
	start := p.cur().Pos
	isRevoke := p.cur().Type == KwRevoke
	kind := strings.ToUpper(p.advance().Val) // GRANT or REVOKE
	stmt := &AccessControlStmt{span: span{Pos: start}, Kind: kind}

	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt.Cluster = cluster

	if isRevoke {
		// Optional `ADMIN OPTION FOR` or `GRANT OPTION FOR`
		if (p.isWord("ADMIN") || p.cur().Type == KwGrant) && p.isWordAt(1, "OPTION") && p.peek(2).Type == KwFor {
			optKind := strings.ToUpper(p.advance().Val)
			p.advance() // OPTION
			p.advance() // FOR
			stmt.RevokeOption = optKind + " OPTION"
		}
	}

	// Privilege list (e.g. `SELECT(a, b), INSERT ON db.*`) OR role list (e.g. `role1, role2 TO/FROM u`)
	privs, err := p.parsePrivilegeOrRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Privileges = privs

	if p.match(KwOn) {
		target, err := p.parsePrivilegeTarget()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	if isRevoke {
		if _, err := p.expect(KwFrom); err != nil {
			return nil, err
		}
	} else {
		if _, err := p.expect(KwTo); err != nil {
			return nil, err
		}
	}

	grantees, except, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	stmt.Grantees = grantees
	stmt.ExceptGrantees = except

	// Optional trailing WITH GRANT OPTION / WITH ADMIN OPTION / WITH REPLACE OPTION
	for p.cur().Type == KwWith {
		p.advance()
		switch {
		case p.cur().Type == KwGrant && p.isWordAt(1, "OPTION"):
			p.advance()
			p.advance()
			stmt.Options = append(stmt.Options, "GRANT OPTION")
		case p.isWord("ADMIN") && p.isWordAt(1, "OPTION"):
			p.advance()
			p.advance()
			stmt.Options = append(stmt.Options, "ADMIN OPTION")
		case p.cur().Type == KwReplace && p.isWordAt(1, "OPTION"):
			p.advance()
			p.advance()
			stmt.Options = append(stmt.Options, "REPLACE OPTION")
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("unexpected token %q after WITH in %s", p.cur().Raw, kind))
		}
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parsePrivilegeOrRoleList() ([]*PrivilegeSpec, error) {
	var list []*PrivilegeSpec
	for {
		if !isIdentifier(p.cur()) {
			return nil, p.errAt(p.cur(), fmt.Sprintf("expected privilege or role name, got %q", p.cur().Raw))
		}
		start := p.cur().Pos
		// Privilege names may be multi-word (e.g. ALL PRIVILEGES, CREATE TABLE, ALTER COLUMN, SHOW TABLES)
		parts := []string{p.advance().Val}
		for isIdentifier(p.cur()) &&
			p.cur().Type != KwOn && p.cur().Type != KwTo && p.cur().Type != KwFrom &&
			p.cur().Type != KwWith {
			parts = append(parts, p.advance().Val)
		}
		spec := &PrivilegeSpec{
			span: span{Pos: start},
			Name: strings.Join(parts, " "),
		}
		// Optional column list: SELECT(col1, col2)
		if p.match(TokLParen) {
			cols, err := p.parseIdentList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokRParen); err != nil {
				return nil, err
			}
			spec.Columns = cols
		}
		spec.End = p.peek(-1).End
		list = append(list, spec)
		if !p.match(TokComma) {
			break
		}
	}
	return list, nil
}

func (p *parser) parsePrivilegeTarget() (string, error) {
	// Targets: *.*, db.*, db.tbl, tbl
	if p.match(TokStar) {
		if p.match(TokDot) {
			if _, err := p.expect(TokStar); err != nil {
				return "", err
			}
			return "*.*", nil
		}
		return "*", nil
	}
	first, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.match(TokDot) {
		if p.match(TokStar) {
			return first + ".*", nil
		}
		second, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		return first + "." + second, nil
	}
	return first, nil
}

func (p *parser) parseGranteeList() (grantees []string, except []string, err error) {
	for {
		if p.match(KwCurrentUser) {
			grantees = append(grantees, "CURRENT_USER")
		} else if p.match(KwAll) {
			grantees = append(grantees, "ALL")
		} else {
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, nil, err
			}
			grantees = append(grantees, name)
		}
		if !p.match(TokComma) {
			break
		}
	}
	if p.match(KwExcept) {
		for {
			if p.match(KwCurrentUser) {
				except = append(except, "CURRENT_USER")
			} else {
				name, err := p.parseIdentifier()
				if err != nil {
					return nil, nil, err
				}
				except = append(except, name)
			}
			if !p.match(TokComma) {
				break
			}
		}
	}
	return grantees, except, nil
}

func (p *parser) parseBackupRestoreStmt() (*BackupRestoreStmt, error) {
	start := p.cur().Pos
	isRestore := p.cur().Type == KwRestore
	kind := strings.ToUpper(p.advance().Val) // BACKUP or RESTORE
	stmt := &BackupRestoreStmt{span: span{Pos: start}, Kind: kind}

	// Comma-separated list of backup/restore elements:
	//   TABLE [db.]t [AS [db.]t2] | TEMPORARY TABLE t | DICTIONARY [db.]d | VIEW [db.]v |
	//   DATABASE db [EXCEPT TABLES ...] | ALL [EXCEPT DATABASES/TABLES ...]
	for {
		eStart := p.cur().Pos
		elem := &BackupElement{span: span{Pos: eStart}}
		switch p.cur().Type {
		case KwAll:
			p.advance()
			elem.Kind = "ALL"
		case KwTemporary:
			p.advance()
			if _, err := p.expect(KwTable); err != nil {
				return nil, err
			}
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			elem.Kind = "TEMPORARY TABLE"
			elem.Name = joinDB(db, tbl)
		case KwTable, KwDictionary, KwView:
			elem.Kind = strings.ToUpper(p.advance().Val)
			db, tbl, err := p.parseTableIdentifier()
			if err != nil {
				return nil, err
			}
			elem.Name = joinDB(db, tbl)
			if p.match(KwAs) {
				adb, atbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				elem.Alias = joinDB(adb, atbl)
			}
		case KwDatabase:
			p.advance()
			elem.Kind = "DATABASE"
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			elem.Name = name
			if p.match(KwAs) {
				alias, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				elem.Alias = alias
			}
		default:
			return nil, p.errAt(p.cur(), fmt.Sprintf("expected TABLE, DATABASE, DICTIONARY, VIEW, or ALL in %s, got %q", kind, p.cur().Raw))
		}
		if p.match(KwExcept) {
			switch p.cur().Type {
			case KwTables, KwDatabases, KwTable, KwDatabase:
				elem.ExceptKind = strings.ToUpper(p.advance().Val)
			default:
				return nil, p.errAt(p.cur(), "expected TABLE(S) or DATABASE(S) after EXCEPT")
			}
			for {
				db, tbl, err := p.parseTableIdentifier()
				if err != nil {
					return nil, err
				}
				elem.Except = append(elem.Except, joinDB(db, tbl))
				if !p.match(TokComma) {
					break
				}
			}
		}
		elem.End = p.peek(-1).End
		stmt.Items = append(stmt.Items, elem)
		if !p.match(TokComma) {
			break
		}
	}

	cluster, err := p.parseClusterClause()
	if err != nil {
		return nil, err
	}
	stmt.Cluster = cluster

	if isRestore {
		if _, err := p.expect(KwFrom); err != nil {
			return nil, err
		}
	} else {
		if _, err := p.expect(KwTo); err != nil {
			return nil, err
		}
	}

	// Destination/source expression: e.g. Disk('backups', '1.zip'), S3(...), File(...), Null
	if p.match(KwNull) {
		stmt.Destination = "Null"
	} else if isIdentifier(p.cur()) && p.peek(1).Type == TokLParen {
		fn, err := p.parseTableFunctionExpr()
		if err != nil {
			return nil, err
		}
		stmt.DestinationFunc = fn
	} else {
		dest, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Destination = dest
	}

	if p.cur().Type == KwSettings {
		settings, err := p.parseSettingsClause()
		if err != nil {
			return nil, err
		}
		stmt.Settings = settings
	}
	stmt.Async = p.match(KwAsync)

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) isWord(upper string) bool {
	return p.isWordAt(0, upper)
}

func (p *parser) isWordAt(offset int, upper string) bool {
	tok := p.peek(offset)
	return (tok.Type == TokIdent || tok.IsKeyword()) && strings.EqualFold(tok.Val, upper)
}
