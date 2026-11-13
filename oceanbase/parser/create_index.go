package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// Called when p.cur is at the INDEX keyword (after consuming CREATE [UNIQUE|FULLTEXT|SPATIAL|VECTOR]).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-index.html
//
//	CREATE [UNIQUE | FULLTEXT | SPATIAL] INDEX index_name
//	    [index_type]
//	    ON tbl_name (key_part,...)
//	    [index_option] ...
//	    [algorithm_option | lock_option] ...
func (p *Parser) parseCreateIndexStmt(start int, unique bool, fulltext bool, spatial bool) (*nodes.CreateIndexStmt, error) {
	p.advance() // consume INDEX

	stmt := &nodes.CreateIndexStmt{
		Loc:      nodes.Loc{Start: start},
		Unique:   unique,
		Fulltext: fulltext,
		Spatial:  spatial,
	}

	// IF NOT EXISTS (MySQL 8.0.27+)
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// index_name
	if err := p.parseIndexRelationName(stmt); err != nil {
		return nil, err
	}

	// Optional index_type. CREATE INDEX reads its name first and cannot leave
	// it out, so both spellings of the clause stand here, where a key
	// definition's position takes a bare TYPE as a name instead.
	if p.isIndexTypeIntro() {
		p.advance()
		typeName, err := p.parseIndexTypeName()
		if err != nil {
			return nil, err
		}
		stmt.IndexType = typeName
	}

	// ON tbl_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	// Completion: after ON, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}
	tbl, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt.Table = tbl

	// (key_part, ...)
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		col, err := p.parseIndexKeyPart()
		if err != nil {
			return nil, err
		}
		stmt.Columns = append(stmt.Columns, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// [index_option] ...
	for {
		opt, ok, err := p.parseIndexOption(stmt.Fulltext)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
	}

	// No [algorithm_option | lock_option]: MySQL closes CREATE INDEX with them
	// and OceanBase does not. 5.7.25-OceanBase_CE-v4.3.5.6 answers 1064 for
	// ALGORITHM = DEFAULT, ALGORITHM = INPLACE, LOCK = DEFAULT, LOCK = NONE and
	// LOCK = EXCLUSIVE after a column list it otherwise parses, so stmt.Algorithm
	// and stmt.Lock stay empty here (read on 2026-11-04).

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseIndexRelationName reads the name CREATE INDEX gives its index, which unlike a key
// definition's is a normal_relation_factor: i, d.i or c.d.i, and any of them with a dblink
// after @. A reserved word stands on either side of a dot: 4.3.5.6 parses `d.select` and
// `select.i`, and answers 1064 for `w.x.y.i` and `.i`. It parses `CREATE INDEX q@l ON t (a)`
// and makes the index q.
//
// Ref: normal_relation_factor in sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
func (p *Parser) parseIndexRelationName(stmt *nodes.CreateIndexStmt) error {
	var parts []string
	for {
		var part string
		switch {
		case len(parts) > 0 && p.cur.Type >= 700, p.cur.Type >= 700 && p.peekNext().Type == '.':
			part = p.advance().Str
		default:
			name, _, err := p.parseIdent()
			if err != nil {
				return err
			}
			part = name
		}
		parts = append(parts, part)
		if p.cur.Type != '.' || len(parts) == 3 {
			break
		}
		p.advance() // consume '.'
	}
	stmt.IndexName = parts[len(parts)-1]
	if len(parts) >= 2 {
		stmt.IndexSchema = parts[len(parts)-2]
	}
	if len(parts) == 3 {
		stmt.IndexCatalog = parts[0]
	}
	if p.isVariableRef() {
		stmt.IndexDBLink = strings.TrimPrefix(p.advance().Str, "@")
	}
	return nil
}

// parseIndexKeyPart parses a single key_part in an index column list.
//
//	key_part:
//	    col_name [(length)] [ASC | DESC]
//	  | (expr) [ASC | DESC]
func (p *Parser) parseIndexKeyPart() (*nodes.IndexColumn, error) {
	// Completion: offer columnref for index column position.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	start := p.pos()
	col := &nodes.IndexColumn{
		Loc: nodes.Loc{Start: start},
	}

	// Functional index: (expr)
	if p.cur.Type == '(' {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		col.Expr = expr
		col.Functional = true
	} else {
		// Column name
		colName, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		col.Expr = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Column: colName,
		}

		// Optional (length)
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type == tokICONST {
				col.Length = int(p.cur.Ival)
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}

	// Optional ASC | DESC
	if _, ok := p.match(kwASC); ok {
		// ASC is default, nothing to set
	} else if _, ok := p.match(kwDESC); ok {
		col.Desc = true
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// indexColumnsToNames extracts simple column names from index columns.
// For functional indexes (expression-based), the column name is empty.
func indexColumnsToNames(cols []*nodes.IndexColumn) []string {
	names := make([]string, 0, len(cols))
	for _, c := range cols {
		if cr, ok := c.Expr.(*nodes.ColumnRef); ok && !c.Functional {
			names = append(names, cr.Column)
		}
	}
	return names
}

// parseParenIndexKeyParts parses a parenthesized list of index key parts.
//
//	(key_part [, key_part] ...)
func (p *Parser) parseParenIndexKeyParts() ([]*nodes.IndexColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var cols []*nodes.IndexColumn
	for {
		col, err := p.parseIndexKeyPart()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseStoringOption reads STORING (col, ...): bare column names, so 4.3.5.6 answers 1064
// for (t.b), (b(3)), (b DESC), for an empty list and a trailing comma. The column is checked
// later (1054), and a second STORING adds to the first.
//
// Ref: `STORING '(' column_name_list ')'` in sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
func (p *Parser) parseStoringOption() (*nodes.IndexOption, bool, error) {
	start := p.pos()
	p.advance() // consume STORING
	if _, err := p.expect('('); err != nil {
		return nil, false, err
	}
	opt := &nodes.IndexOption{Name: "STORING"}
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		opt.Columns = append(opt.Columns, name)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, false, err
	}
	opt.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return opt, true, nil
}

// parseStorageCachePolicy reads STORAGE_CACHE_POLICY [=] (...): GLOBAL = 'string' alone, or
// a list of BOUNDARY_COLUMN = col, BOUNDARY_COLUMN_UNIT = 'string', GRANULARITY = 'string'
// and HOT_RETENTION = n unit. 4.3.5.6 parses them and refuses them later with 1235, outside
// shared storage, and answers 1064 for GLOBAL = hot, an empty list and a bare string.
//
// Ref: storage_cache_policy_attribute_list and retention_time_unit in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
func (p *Parser) parseStorageCachePolicy() (*nodes.IndexOption, bool, error) {
	start := p.pos()
	p.advance() // consume STORAGE_CACHE_POLICY
	p.match('=')
	if _, err := p.expect('('); err != nil {
		return nil, false, err
	}
	opt := &nodes.IndexOption{Name: "STORAGE_CACHE_POLICY"}
	stringValue := func(prop *nodes.ParserProperty) error {
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}, Value: p.cur.Str}
		p.advance()
		return nil
	}
	for {
		prop := &nodes.ParserProperty{Loc: nodes.Loc{Start: p.pos()}, Name: strings.ToLower(p.cur.Str)}
		global := p.cur.Type == kwGLOBAL
		if !global && p.cur.Type != tokIDENT && p.cur.Type < 700 {
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance() // consume the name
		if _, err := p.expect('='); err != nil {
			return nil, false, err
		}
		var err error
		switch prop.Name {
		case "global", "boundary_column_unit", "granularity":
			err = stringValue(prop)
		case "boundary_column":
			var col string
			valStart := p.pos()
			if col, _, err = p.parseIdent(); err == nil {
				prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: col}
			}
		case "hot_retention":
			valStart := p.pos()
			if p.cur.Type != tokICONST {
				return nil, false, p.syntaxErrorAtCur()
			}
			count := p.advance().Str
			if !isStorageRetentionUnit(p.cur) {
				return nil, false, p.syntaxErrorAtCur()
			}
			unit := p.advance().Str
			prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: count + " " + strings.ToUpper(unit)}
		default:
			return nil, false, &ParseError{Message: "unknown storage cache policy attribute", Position: prop.Loc.Start}
		}
		if err != nil {
			return nil, false, err
		}
		prop.Loc.End = p.prev.End
		opt.Properties = append(opt.Properties, prop)
		// GLOBAL stands alone.
		if global {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, false, err
	}
	opt.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return opt, true, nil
}

// isStorageRetentionUnit reports whether tok is a unit HOT_RETENTION takes.
func isStorageRetentionUnit(tok Token) bool {
	switch strings.ToUpper(tok.Str) {
	case "YEAR", "MONTH", "WEEK", "DAY", "HOUR", "MINUTE", "SECOND":
		return true
	}
	return false
}

// parseIndexParams reads the WITH (name = value, ...) a vector index takes its
// parameters from: WITH (distance = l2, type = hnsw, lib = vsag).
//
// It is an index option of every index and not only a vector one: 4.3.5.6 parses
// CREATE INDEX i ON t (a) WITH (foo=bar) and ADD INDEX i (a) WITH (foo=bar) and refuses
// them later with 1235, set index params in not vector index is not supported. The
// names and the words are not checked until then either: WITH (foo=bar) parses on a
// vector index.
//
// A name is an identifier. A value is an identifier, an unsigned number or a boolean:
// WITH (m=16), WITH (ef_search=1.5) and WITH (a=TRUE) parse, and 1064 is the answer for
// a string, DEFAULT, NULL, -1, +1, a parenthesized value, an expression and b.c, for an
// empty list, a missing = and a trailing comma.
func (p *Parser) parseIndexParams() (*nodes.IndexOption, bool, error) {
	start := p.pos()
	p.advance() // consume WITH
	p.advance() // consume (
	opt := &nodes.IndexOption{Name: "WITH"}
	for {
		prop := &nodes.ParserProperty{Loc: nodes.Loc{Start: p.pos()}}
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		prop.Name = strings.ToLower(name)
		if _, err := p.expect('='); err != nil {
			return nil, false, err
		}
		valStart := p.pos()
		switch p.cur.Type {
		case tokICONST:
			prop.Value = &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.cur.End}, Value: p.cur.Ival}
			p.advance()
		case tokFCONST, kwTRUE, kwFALSE:
			prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.cur.End}, Value: p.cur.Str}
			p.advance()
		default:
			word, _, err := p.parseIdent()
			if err != nil {
				return nil, false, err
			}
			prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: word}
		}
		prop.Loc.End = p.prev.End
		opt.Properties = append(opt.Properties, prop)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, false, err
	}
	opt.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return opt, true, nil
}

// parseIndexOption parses a single index_option.
//
//	index_option:
//	    KEY_BLOCK_SIZE [=] value
//	    | USING {BTREE | HASH}
//	    | WITH (name = value, ...)
//	    | WITH ROWID
//	    | PARALLEL [=] n | NOPARALLEL
//	    | STORING (col, ...)
//	    | DATA_TABLE_ID [=] n | INDEX_TABLE_ID [=] n | MAX_USED_PART_ID [=] n
//	    | STORAGE_CACHE_POLICY [=] (...)
//	    | WITH PARSER parser_name          -- fulltext only
//	    | PARSER_PROPERTIES [=] (...)      -- fulltext only
//	    | COMMENT 'string'
//	    | {VISIBLE | INVISIBLE}
//
// fulltext is set for a FULLTEXT index, whose options OceanBase reads as
// fulltext_index_option (sql_parser_mysql_mode.y, v4.3.5_CE_BP6): index_option, WITH
// PARSER and PARSER_PROPERTIES. Every other index reads index_option alone, so
// 4.3.5.6 answers 1064 for WITH PARSER there, at PARSER, and for PARSER_PROPERTIES
// (fulltextParserCases).
func (p *Parser) parseIndexOption(fulltext bool) (*nodes.IndexOption, bool, error) {
	start := p.pos()

	switch {
	case p.cur.Type == kwKEY_BLOCK_SIZE:
		p.advance()
		p.match('=') // optional =
		valStart := p.pos()
		if p.cur.Type == tokICONST {
			val := p.cur.Ival
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "KEY_BLOCK_SIZE",
				Value: &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val},
			}, true, nil
		}
		// Could be an identifier value
		v, _, err := p.parseIdentifier()
		if err != nil {
			return nil, false, err
		}
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "KEY_BLOCK_SIZE",
			Value: &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: v},
		}, true, nil

	case p.isIndexTypeIntro():
		p.advance()
		valStart := p.pos()
		typeName, err := p.parseIndexTypeName()
		if err != nil {
			return nil, false, err
		}
		// Either spelling is recorded as USING: TYPE is a synonym for it, and a
		// consumer reading the index type should not have to know which was written.
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "USING",
			Value: &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: typeName},
		}, true, nil

	case p.cur.Type == kwWITH:
		next := p.peekNext()
		if next.Type == kwPARSER {
			if !fulltext {
				return nil, false, p.syntaxErrorAtTok(next)
			}
			p.advance() // consume WITH
			p.advance() // consume PARSER
			valStart := p.pos()
			parserName, _, err := p.parseIdent()
			if err != nil {
				return nil, false, err
			}
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "PARSER",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: parserName},
			}, true, nil
		}
		if next.Type == '(' {
			return p.parseIndexParams()
		}
		// WITH ROWID. It changes nothing SHOW CREATE TABLE prints.
		if next.Type == tokIDENT && strings.EqualFold(next.Str, "ROWID") {
			p.advance() // consume WITH
			p.advance() // consume ROWID
			return &nodes.IndexOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "WITH ROWID"}, true, nil
		}
		return nil, false, nil

	// PARALLEL [=] n, n at least 1, and NOPARALLEL set the degree the index is built with.
	// 4.3.5.6 answers 1064 for PARALLEL 0 and takes either any number of times, and neither
	// changes what SHOW CREATE TABLE prints.
	case p.cur.Type == kwPARALLEL:
		p.advance()
		p.match('=')
		valStart := p.pos()
		if p.cur.Type != tokICONST || p.cur.Ival < 1 {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Ival
		p.advance()
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "PARALLEL",
			Value: &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val},
		}, true, nil

	case p.isBareWord("NOPARALLEL"):
		p.advance()
		return &nodes.IndexOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "NOPARALLEL"}, true, nil

	case p.isBareWord("STORING"):
		return p.parseStoringOption()

	// DATA_TABLE_ID, INDEX_TABLE_ID and MAX_USED_PART_ID [=] n. The first two parse and
	// are refused later with 4179, and VIRTUAL_COLUMN_ID, which the grammar lists beside
	// them, is a 1064 on 4.3.5.6 whatever follows it.
	case p.isBareWord("DATA_TABLE_ID"), p.isBareWord("INDEX_TABLE_ID"), p.isBareWord("MAX_USED_PART_ID"):
		name := strings.ToUpper(p.advance().Str)
		p.match('=')
		valStart := p.pos()
		if p.cur.Type != tokICONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Ival
		p.advance()
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  name,
			Value: &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val},
		}, true, nil

	case p.isBareWord("STORAGE_CACHE_POLICY"):
		return p.parseStorageCachePolicy()

	case fulltext && p.isBareWord("PARSER_PROPERTIES"):
		return p.parseParserProperties()

	case p.cur.Type == kwCOMMENT:
		p.advance()
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "COMMENT",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: p.prev.Loc, End: p.prev.End}, Value: val},
			}, true, nil
		}
		return nil, false, &ParseError{Message: "expected string after COMMENT", Position: p.cur.Loc}

	// LOCAL and GLOBAL give an index its scope: LOCAL is partitioned with the table,
	// GLOBAL spans it. 4.3.5 prints a scope on every secondary index, partitioned table
	// or not, and takes either in CREATE TABLE and in ALTER TABLE ... ADD INDEX.
	case p.cur.Type == kwLOCAL, p.cur.Type == kwGLOBAL:
		name := "LOCAL"
		if p.cur.Type == kwGLOBAL {
			name = "GLOBAL"
		}
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: name,
		}, true, nil

	// BLOCK_SIZE is OceanBase's per-index block size. 4.3.5 takes it with no '=',
	// which MySQL's KEY_BLOCK_SIZE requires, so the '=' is optional here.
	case p.cur.Type == kwBLOCK_SIZE:
		p.advance()
		p.match('=') // optional =
		valStart := p.pos()
		if p.cur.Type != tokICONST {
			return nil, false, &ParseError{Message: "expected an integer for BLOCK_SIZE", Position: p.cur.Loc}
		}
		val := p.cur.Ival
		p.advance()
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "BLOCK_SIZE",
			Value: &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val},
		}, true, nil

	case p.cur.Type == kwVISIBLE:
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: "VISIBLE",
		}, true, nil

	case p.cur.Type == kwINVISIBLE:
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: "INVISIBLE",
		}, true, nil

	case p.cur.Type == kwENGINE_ATTRIBUTE:
		p.advance()
		p.match('=')
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "ENGINE_ATTRIBUTE",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: p.prev.Loc, End: p.prev.End}, Value: val},
			}, true, nil
		}
		return nil, false, &ParseError{Message: "expected string after ENGINE_ATTRIBUTE", Position: p.cur.Loc}

	case p.cur.Type == kwSECONDARY_ENGINE_ATTRIBUTE:
		p.advance()
		p.match('=')
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "SECONDARY_ENGINE_ATTRIBUTE",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: p.prev.Loc, End: p.prev.End}, Value: val},
			}, true, nil
		}
		return nil, false, &ParseError{Message: "expected string after SECONDARY_ENGINE_ATTRIBUTE", Position: p.cur.Loc}
	}

	return nil, false, nil
}

// parserPropertyKinds maps each key OceanBase's fulltext_parser_properties takes to
// whether its value is an integer (true) or a string (false).
var parserPropertyKinds = map[string]bool{
	"min_token_size":   true,
	"max_token_size":   true,
	"ngram_token_size": true,
	"min_ngram_size":   true,
	"max_ngram_size":   true,
	"stopword_table":   false,
	"dict_table":       false,
	"quantifier_table": false,
	"ik_mode":          false,
}

// parseParserProperties parses OceanBase's PARSER_PROPERTIES option on a fulltext
// index, from PARSER_PROPERTIES on:
//
//	PARSER_PROPERTIES [=] '(' property [[,] property] ... ')'
//	property: key '=' value
//
// as fulltext_parser_properties_list reads it (sql_parser_mysql_mode.y,
// v4.3.5_CE_BP6). SHOW CREATE TABLE on 4.3.5.6 prints the option on every fulltext
// index, as PARSER_PROPERTIES=(min_token_size=3,max_token_size=84). The key is one of
// parserPropertyKinds, bare and in any case; the '=' is required; a size takes an
// unsigned integer and the rest a string literal, adjacent literals included. The
// engine answers 1064 for anything else, an empty list or a trailing comma included
// (fulltextParserCases). Which keys a parser takes, and their ranges, it checks after
// parsing.
func (p *Parser) parseParserProperties() (*nodes.IndexOption, bool, error) {
	start := p.pos()
	p.advance() // consume PARSER_PROPERTIES
	p.match('=')
	if _, err := p.expect('('); err != nil {
		return nil, false, err
	}
	opt := &nodes.IndexOption{Name: "PARSER_PROPERTIES"}
	for {
		if p.cur.Type != tokIDENT && p.cur.Type < 700 {
			return nil, false, p.syntaxErrorAtCur()
		}
		key := strings.ToLower(p.cur.Str)
		isInt, known := parserPropertyKinds[key]
		if !known || !p.isBareWord(key) {
			return nil, false, p.syntaxErrorAtCur()
		}
		propStart := p.pos()
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, false, err
		}
		prop := &nodes.ParserProperty{Name: key}
		valStart := p.pos()
		switch {
		case isInt && p.cur.Type == tokICONST:
			val := p.cur.Ival
			p.advance()
			prop.Value = &nodes.IntLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val}
		case !isInt && p.cur.Type == tokSCONST:
			var val strings.Builder
			for p.cur.Type == tokSCONST {
				val.WriteString(p.cur.Str)
				p.advance()
			}
			prop.Value = &nodes.StringLit{Loc: nodes.Loc{Start: valStart, End: p.prev.End}, Value: val.String()}
		default:
			return nil, false, p.syntaxErrorAtCur()
		}
		prop.Loc = nodes.Loc{Start: propStart, End: p.prev.End}
		opt.Properties = append(opt.Properties, prop)
		if p.cur.Type == ')' {
			break
		}
		p.match(',')
	}
	p.advance() // consume ')'
	opt.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return opt, true, nil
}

// isIndexTypeIntro reports whether the current token opens an index_type
// clause. OceanBase spells it USING and only USING: MySQL also takes TYPE, and
// 5.7.25-OceanBase_CE-v4.3.5.6 answers 1064 for every position it can stand in
// -- CREATE INDEX i ON t (a) TYPE BTREE, CREATE INDEX i TYPE BTREE ON t (a),
// KEY k (a) TYPE BTREE and KEY k TYPE BTREE (a) -- so taking it here would only
// widen an over-acceptance (read on 2026-11-04).
//
// Ref: index_using_algorithm in OceanBase sql_parser_mysql_mode.y.
// https://en.oceanbase.com/docs/common-oceanbase-database-10000000002017621
func (p *Parser) isIndexTypeIntro() bool {
	return p.cur.Type == kwUSING
}

// parseIndexTypeName reads the name an index_type clause takes after USING or
// TYPE.
//
// index_type_name is a keyword list, not an identifier, so a name that is no
// index type and a backquoted `BTREE` are both syntax errors: oceanbase/oceanbase-ce:4.3.5-lts
// answers 1064 for USING nosuchtype and for USING `BTREE`.
//
// Ref: index_type_name in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseIndexTypeName() (string, error) {
	if tok, ok := p.match(kwBTREE, kwHASH); ok {
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}
