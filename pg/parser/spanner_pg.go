package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Dialect selects dialect-specific syntax rules within PostgreSQL.
type Dialect int

const (
	// DialectDefault is standard PostgreSQL 17 used by Parse and ParseRange.
	DialectDefault Dialect = iota
	// DialectSpannerPG applies Cloud Spanner PostgreSQL dialect rules:
	// matches Cloud Spanner's PostgreSQL 16-based parser restrictions (no
	// XML expressions, no PG 17 SQL/JSON query functions or PG 17 grammar
	// additions, non-empty statement inputs, and Spanner-PG grammar
	// restrictions) and enables Cloud Spanner PostgreSQL DDL extensions
	// (table/index interleaving, TTL clauses, change streams, SQL SECURITY
	// views, bit-reversed sequences, search indexes, locality groups, and
	// HIDDEN / VIRTUAL generated columns).
	DialectSpannerPG
)

// ParseWithDialect parses a SQL string using the specified PostgreSQL dialect.
func ParseWithDialect(sql string, dialect Dialect) (*nodes.List, error) {
	return ParseRangeWithDialect(sql, 0, len(sql), dialect)
}

func (p *Parser) isSpannerPG() bool {
	return p.dialect == DialectSpannerPG
}

func isSpannerIdentTok(tok Token, kw string) bool {
	return tok.Type == IDENT && strings.EqualFold(tok.Str, kw)
}

func (p *Parser) isSpannerIdent(kw string) bool {
	return isSpannerIdentTok(p.cur, kw)
}

// isPG17SQLJsonKeyword reports whether tok is one of the SQL/JSON query
// keywords added in PostgreSQL 17 that are ordinary function identifiers in
// Spanner's PostgreSQL 16-based parser.
func isPG17SQLJsonKeyword(tok int) bool {
	switch tok {
	case JSON, JSON_SCALAR, JSON_SERIALIZE, JSON_QUERY, JSON_EXISTS, JSON_VALUE, JSON_TABLE:
		return true
	}
	return false
}

func (p *Parser) spannerXMLUnsupportedError() *ParseError {
	return &ParseError{
		Code:     "0A000",
		Message:  "XML expressions are not supported",
		Position: p.pos(),
	}
}

// parseSpannerTableInterleave parses Spanner PG's table interleaving clause:
//
//	INTERLEAVE IN [PARENT] qualified_name [ON DELETE <action>]
//
// When allowFullKeyAction is true (CREATE TABLE), any PostgreSQL key_action is
// accepted after ON DELETE (matching Spanner PG's grammar, which rejects
// non-CASCADE/NO ACTION actions semantically). When false (ALTER TABLE SET
// INTERLEAVE), only CASCADE and NO ACTION are accepted.
func (p *Parser) parseSpannerTableInterleave(allowFullKeyAction bool) (*nodes.SpannerInterleaveClause, error) {
	loc := p.pos()
	p.advance() // consume INTERLEAVE
	if _, err := p.expect(IN_P); err != nil {
		return nil, err
	}

	inParent := false
	// PARENT is an unreserved identifier; if it is followed by another ColId
	// that is not ON and not TTL, treat it as the optional PARENT keyword.
	// Otherwise PARENT itself is the parent table name.
	if p.isSpannerIdent("parent") {
		next := p.peekNext()
		if isColIdToken(next) && next.Type != ON && !isSpannerIdentTok(next, "ttl") {
			p.advance() // consume PARENT
			inParent = true
		}
	}

	parentStart := p.pos()
	parentNames, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	parentRv := makeRangeVarFromNames(parentNames, nodes.Loc{Start: parentStart, End: p.prev.End})

	var onDelete byte
	if p.cur.Type == ON {
		p.advance() // consume ON
		if _, err := p.expect(DELETE_P); err != nil {
			return nil, err
		}
		if allowFullKeyAction {
			var cols *nodes.List
			onDelete, err = p.parseKeyActionType(&cols)
			if err != nil {
				return nil, err
			}
		} else {
			onDelete, err = p.parseSpannerOnDeleteAction()
			if err != nil {
				return nil, err
			}
		}
	}

	return &nodes.SpannerInterleaveClause{
		Parent:   parentRv,
		InParent: inParent,
		OnDelete: onDelete,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseSpannerOnDeleteAction parses CASCADE | NO ACTION for ALTER TABLE ... SET ON DELETE.
func (p *Parser) parseSpannerOnDeleteAction() (byte, error) {
	switch p.cur.Type {
	case CASCADE:
		p.advance()
		return 'c', nil
	case NO:
		p.advance()
		if _, err := p.expect(ACTION); err != nil {
			return 0, err
		}
		return 'a', nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
}

// parseSpannerTTLClause parses Spanner PG's TTL clause:
//
//	TTL INTERVAL 'literal' [qualifier] ON ColId
func (p *Parser) parseSpannerTTLClause() (*nodes.SpannerTTLClause, error) {
	loc := p.pos()
	p.advance() // consume TTL
	if p.cur.Type != INTERVAL {
		return nil, p.syntaxErrorAtCur()
	}
	intervalExpr, err := p.parseIntervalConst()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	colName, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.SpannerTTLClause{
		Interval: intervalExpr,
		Column:   colName,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseCreateChangeStreamStmt parses Spanner PG's CREATE CHANGE STREAM statement:
//
//	CREATE CHANGE STREAM [IF NOT EXISTS] qualified_name
//	    [FOR ALL | FOR table_name [([col, ...])] [, ...]]
//	    [WITH (option = value [, ...])]
//
// CREATE has already been consumed; cur is IDENT("change").
func (p *Parser) parseCreateChangeStreamStmt(startLoc int) (*nodes.CreateChangeStreamStmt, error) {
	p.advance() // consume CHANGE
	if !p.isSpannerIdent("stream") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume STREAM

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}

	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	var forAll bool
	var watched *nodes.List
	if p.cur.Type == FOR {
		p.advance() // consume FOR
		forAll, watched, err = p.parseChangeStreamForTargets()
		if err != nil {
			return nil, err
		}
	}

	options, err := p.parseOptReloptions()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateChangeStreamStmt{
		Name:           rv,
		IfNotExists:    ifNotExists,
		ForAll:         forAll,
		WatchedObjects: watched,
		Options:        options,
		Loc:            nodes.Loc{Start: startLoc, End: p.prev.End},
	}, nil
}

// parseAlterChangeStreamStmt parses Spanner PG's ALTER CHANGE STREAM statement:
//
//	ALTER CHANGE STREAM qualified_name SET FOR {ALL | table_list}
//	ALTER CHANGE STREAM qualified_name SET (option = value [, ...])
//	ALTER CHANGE STREAM qualified_name RESET (option [, ...])
//	ALTER CHANGE STREAM qualified_name DROP FOR {ALL | table_list}
//
// ALTER has already been consumed; cur is IDENT("change").
func (p *Parser) parseAlterChangeStreamStmt(startLoc int) (*nodes.AlterChangeStreamStmt, error) {
	p.advance() // consume CHANGE
	if !p.isSpannerIdent("stream") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume STREAM

	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	stmt := &nodes.AlterChangeStreamStmt{Name: rv}
	switch p.cur.Type {
	case SET:
		p.advance() // consume SET
		if p.cur.Type == FOR {
			p.advance() // consume FOR
			forAll, watched, err := p.parseChangeStreamForTargets()
			if err != nil {
				return nil, err
			}
			stmt.Action = "set_for"
			stmt.ForAll = forAll
			stmt.WatchedObjects = watched
		} else if p.cur.Type == '(' {
			opts, err := p.parseReloptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = "set_options"
			stmt.Options = opts
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	case RESET:
		p.advance() // consume RESET
		opts, err := p.parseReloptions()
		if err != nil {
			return nil, err
		}
		stmt.Action = "reset_options"
		stmt.ResetOptions = opts
	case DROP:
		p.advance() // consume DROP
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		forAll, watched, err := p.parseChangeStreamForTargets()
		if err != nil {
			return nil, err
		}
		stmt.Action = "drop_for"
		stmt.ForAll = forAll
		stmt.WatchedObjects = watched
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = nodes.Loc{Start: startLoc, End: p.prev.End}
	return stmt, nil
}

// parseDropChangeStreamStmt parses Spanner PG's DROP CHANGE STREAM statement:
//
//	DROP CHANGE STREAM [IF EXISTS] qualified_name
//
// DROP has already been consumed; cur is IDENT("change").
func (p *Parser) parseDropChangeStreamStmt() (*nodes.DropStmt, error) {
	p.advance() // consume CHANGE
	if !p.isSpannerIdent("stream") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume STREAM

	missingOk := p.parseOptIfExists()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{names}},
		RemoveType: int(nodes.OBJECT_CHANGE_STREAM),
		Missing_ok: missingOk,
	}, nil
}

// parseChangeStreamForTargets parses ALL or a comma-separated list of
// qualified_name [ '(' [columnList] ')' ].
func (p *Parser) parseChangeStreamForTargets() (bool, *nodes.List, error) {
	if p.cur.Type == ALL {
		p.advance()
		return true, nil, nil
	}
	var items []nodes.Node
	for {
		itemLoc := p.pos()
		names, err := p.parseQualifiedName()
		if err != nil {
			return false, nil, err
		}
		rv := makeRangeVarFromNames(names, nodes.Loc{Start: itemLoc, End: p.prev.End})
		item := &nodes.ChangeStreamForItem{
			Table:      rv,
			AllColumns: true,
		}
		if p.cur.Type == '(' {
			p.advance() // consume '('
			item.AllColumns = false
			if p.cur.Type != ')' {
				cols, err := p.parseColumnList()
				if err != nil {
					return false, nil, err
				}
				item.Columns = cols
			} else {
				item.Columns = &nodes.List{}
			}
			if _, err := p.expect(')'); err != nil {
				return false, nil, err
			}
		}
		item.Loc = nodes.Loc{Start: itemLoc, End: p.prev.End}
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return false, &nodes.List{Items: items}, nil
}

// parseCreateSearchIndexStmt parses Spanner PG's CREATE SEARCH INDEX statement:
//
//	CREATE SEARCH INDEX name ON relation_expr '(' index_params ')'
//	    [INCLUDE '(' index_params ')']
//	    [PARTITION BY col1 [, ...]]
//	    [ORDER BY sortby_list]
//	    [INTERLEAVE IN qualified_name]
//	    [WHERE a_expr]
//	    [WITH (reloption_list)]
//
// CREATE has already been consumed; cur is SEARCH.
func (p *Parser) parseCreateSearchIndexStmt(startLoc int) (*nodes.IndexStmt, error) {
	p.advance() // consume SEARCH
	if _, err := p.expect(INDEX); err != nil {
		return nil, err
	}

	idxname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	indexParams, err := p.parseIndexParams()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Optional INCLUDE '(' index_params ')'
	var includeParams *nodes.List
	if p.cur.Type == INCLUDE {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		includeParams, err = p.parseIndexParams()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// Optional PARTITION BY col1 [, ...]
	var partitionBy *nodes.List
	if p.cur.Type == PARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		partitionBy, err = p.parseColumnList()
		if err != nil {
			return nil, err
		}
	}

	// Optional ORDER BY sortby_list
	var orderBy *nodes.List
	if p.cur.Type == ORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		orderBy, err = p.parseSortByList()
		if err != nil {
			return nil, err
		}
	}

	// Optional INTERLEAVE IN qualified_name
	var interleaveParent *nodes.RangeVar
	if p.isSpannerIdent("interleave") {
		p.advance() // consume INTERLEAVE
		if _, err := p.expect(IN_P); err != nil {
			return nil, err
		}
		parentStart := p.pos()
		parentNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		interleaveParent = makeRangeVarFromNames(parentNames, nodes.Loc{Start: parentStart, End: p.prev.End})
	}

	// Optional WHERE predicate
	whereClause, err := p.parseWhereClause()
	if err != nil {
		return nil, err
	}

	// Optional WITH (reloption_list)
	options, err := p.parseOptReloptions()
	if err != nil {
		return nil, err
	}

	return &nodes.IndexStmt{
		Idxname:              idxname,
		Relation:             rel,
		IndexParams:          indexParams,
		IndexIncludingParams: includeParams,
		PartitionBy:          partitionBy,
		OrderBy:              orderBy,
		InterleaveParent:     interleaveParent,
		WhereClause:          whereClause,
		Options:              options,
		IsSearch:             true,
		Loc:                  nodes.Loc{Start: startLoc, End: p.prev.End},
	}, nil
}

// parseAlterSearchIndexStmt parses Spanner PG's ALTER SEARCH INDEX statement:
//
//	ALTER SEARCH INDEX qualified_name {ADD | DROP} COLUMN ColId
//
// ALTER has already been consumed; cur is SEARCH.
func (p *Parser) parseAlterSearchIndexStmt(alterLoc int) (*nodes.AlterTableStmt, error) {
	p.advance() // consume SEARCH
	if _, err := p.expect(INDEX); err != nil {
		return nil, err
	}

	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	cmdLoc := p.pos()
	var subtype nodes.AlterTableType
	switch p.cur.Type {
	case ADD_P:
		p.advance()
		subtype = nodes.AT_AddColumn
	case DROP:
		p.advance()
		subtype = nodes.AT_DropColumn
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(COLUMN); err != nil {
		return nil, err
	}
	colName, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	cmd := &nodes.AlterTableCmd{
		Subtype: int(subtype),
		Name:    colName,
		Loc:     nodes.Loc{Start: cmdLoc, End: p.prev.End},
	}
	return &nodes.AlterTableStmt{
		Relation: rv,
		Cmds:     &nodes.List{Items: []nodes.Node{cmd}},
		ObjType:  int(nodes.OBJECT_SEARCH_INDEX),
		Loc:      nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseDropSearchIndexStmt parses Spanner PG's DROP SEARCH INDEX statement:
//
//	DROP SEARCH INDEX [IF EXISTS] qualified_name
//
// DROP has already been consumed; cur is SEARCH.
func (p *Parser) parseDropSearchIndexStmt() (*nodes.DropStmt, error) {
	p.advance() // consume SEARCH
	if _, err := p.expect(INDEX); err != nil {
		return nil, err
	}
	missingOk := p.parseOptIfExists()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{names}},
		RemoveType: int(nodes.OBJECT_SEARCH_INDEX),
		Missing_ok: missingOk,
	}, nil
}

// parseCreateLocalityGroupStmt parses Spanner PG's CREATE LOCALITY GROUP statement:
//
//	CREATE LOCALITY GROUP [IF NOT EXISTS] qualified_name
//
// CREATE has already been consumed; cur is IDENT("locality").
func (p *Parser) parseCreateLocalityGroupStmt(startLoc int) (*nodes.CreateLocalityGroupStmt, error) {
	p.advance() // consume LOCALITY
	if _, err := p.expect(GROUP_P); err != nil {
		return nil, err
	}
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	return &nodes.CreateLocalityGroupStmt{
		Name:        rv,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: startLoc, End: p.prev.End},
	}, nil
}

// parseDropLocalityGroupStmt parses Spanner PG's DROP LOCALITY GROUP statement:
//
//	DROP LOCALITY GROUP [IF EXISTS] qualified_name
//
// DROP has already been consumed; cur is IDENT("locality").
func (p *Parser) parseDropLocalityGroupStmt() (*nodes.DropStmt, error) {
	p.advance() // consume LOCALITY
	if _, err := p.expect(GROUP_P); err != nil {
		return nil, err
	}
	missingOk := p.parseOptIfExists()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{names}},
		RemoveType: int(nodes.OBJECT_LOCALITY_GROUP),
		Missing_ok: missingOk,
	}, nil
}

// isSpannerSeqOptElemLead reports whether cur starts a Spanner-PG-specific
// sequence option (BIT_REVERSED_POSITIVE or SKIP RANGE).
func (p *Parser) isSpannerSeqOptElemLead() bool {
	if !p.isSpannerPG() {
		return false
	}
	return p.isSpannerIdent("bit_reversed_positive") || p.cur.Type == SKIP
}
