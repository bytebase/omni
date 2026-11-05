package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// Called when p.cur is at the INDEX keyword (after consuming CREATE [UNIQUE|FULLTEXT|SPATIAL|VECTOR]).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-index.html
//
//	CREATE [UNIQUE | FULLTEXT | SPATIAL | VECTOR] INDEX index_name
//	    [index_type]
//	    ON tbl_name (key_part,...)
//	    [index_option] ...
//	    [algorithm_option | lock_option] ...
func (p *Parser) parseCreateIndexStmt(start int, unique, fulltext, spatial, vector bool) (*nodes.CreateIndexStmt, error) {
	p.advance() // consume INDEX

	stmt := &nodes.CreateIndexStmt{
		Loc:      nodes.Loc{Start: start},
		Unique:   unique,
		Fulltext: fulltext,
		Spatial:  spatial,
		Vector:   vector,
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
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.IndexName = name

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
	tbl, err := p.parseTableRef()
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
		opt, ok, err := p.parseIndexOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
	}

	// [algorithm_option | lock_option] ...
	algorithm, lock, err := p.parseIndexLockAndAlgorithm()
	if err != nil {
		return nil, err
	}
	stmt.Algorithm, stmt.Lock = algorithm, lock

	if createIndexRefused(stmt) {
		return nil, refuseVectorIndexShape(start)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
		// ColumnName. The server looks the name up in the table it is indexing, whatever
		// the qualifier says, so `KEY (x.a)` indexes a.
		schema, table, colName, err := p.parseColumnName()
		if err != nil {
			return nil, err
		}
		col.Expr = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Schema: schema,
			Table:  table,
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
		if cr, ok := c.Expr.(*nodes.ColumnRef); ok {
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

// parseIndexOption parses a single index_option.
//
//	index_option:
//	    KEY_BLOCK_SIZE [=] value
//	    | USING {BTREE | HASH}
//	    | WITH PARSER parser_name
//	    | COMMENT 'string'
//	    | {VISIBLE | INVISIBLE}
//	    | {CLUSTERED | NONCLUSTERED}
//	    | {GLOBAL | LOCAL}
//
// CLUSTERED, NONCLUSTERED, GLOBAL and LOCAL are TiDB's own, and are index
// options rather than a suffix of one key kind: IndexOption in v8.5.5
// pkg/parser/parser.y reads `IndexInvisible | WithClustered | "GLOBAL" |
// "LOCAL"` beside COMMENT and the rest, so they appear in any order and more
// than once. v8.5.5 runs `ADD INDEX (a) GLOBAL COMMENT 'c'`,
// `ADD INDEX (a) COMMENT 'c' GLOBAL`, `ADD INDEX (a) GLOBAL LOCAL`,
// `ADD INDEX (a) NONCLUSTERED`, `ADD PRIMARY KEY (a) CLUSTERED GLOBAL`,
// `ADD PRIMARY KEY (a) GLOBAL CLUSTERED`, `CREATE TABLE t (a int, INDEX (a)
// GLOBAL)` and `CREATE INDEX i ON t (a) GLOBAL`; MySQL 8.0.46 and MariaDB
// 11.8.8 answer 1064 to every one of them (read on 2026-11-04).
func (p *Parser) parseIndexOption() (*nodes.IndexOption, bool, error) {
	// Completion: offer the index options, TiDB's own among them. Record the
	// candidates without failing, because this runs inside an option loop
	// whose caller expects "no option here" rather than an error.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwKEY_BLOCK_SIZE, kwUSING, kwWITH, kwCOMMENT, kwVISIBLE, kwINVISIBLE,
			kwCLUSTERED, kwNONCLUSTERED, kwGLOBAL, kwLOCAL,
		} {
			p.addTokenCandidate(t)
		}
	}

	start := p.pos()

	switch {
	case p.cur.Type == kwGLOBAL, p.cur.Type == kwLOCAL,
		p.cur.Type == kwCLUSTERED, p.cur.Type == kwNONCLUSTERED:
		name := strings.ToUpper(p.cur.Str)
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: name,
		}, true, nil

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
		return nil, false, nil

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

// isIndexTypeIntro reports whether the current token opens an index_type
// clause. MySQL spells it either way, USING or TYPE.
//
// Ref: index_type_clause in MySQL 8.0 sql/sql_yacc.yy.
// https://dev.mysql.com/doc/refman/8.0/en/create-index.html
func (p *Parser) isIndexTypeIntro() bool {
	return p.cur.Type == kwUSING || p.cur.Type == kwTYPE
}

// parseIndexTypeName reads the name an index_type clause takes after USING or
// TYPE.
//
// IndexTypeName is a keyword list, not an identifier, so a name that is no
// index type and a backquoted `BTREE` are both syntax errors: pingcap/tidb:v8.5.5
// answers 1064 for USING nosuchtype and for USING `BTREE`.
//
// TiDB adds HNSW and HYPO, the two index types MySQL has no token for, and
// takes each unquoted only: v8.5.5 parses USING HNSW, USING hnsw and
// USING HYPO, and answers 1064 for USING `HNSW` and USING `HYPO` (read on
// 2026-11-04). Neither is in this parser's keyword table, so both are read as
// bare words rather than as tokens.
//
// Ref: IndexTypeName in pingcap/tidb pkg/parser/parser.y; index_type_name in
// MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseIndexTypeName() (string, error) {
	if tok, ok := p.match(kwBTREE, kwHASH, kwRTREE); ok {
		return tok.Str, nil
	}
	if p.isBareWord("HNSW") || p.isBareWord("HYPO") {
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseIndexLockAndAlgorithm reads the ALGORITHM and LOCK clauses that close a
// CREATE INDEX or a DROP INDEX statement.
//
// The grammar takes at most one of each, in either order, so a second
// ALGORITHM or LOCK is a syntax error rather than a value that wins: pingcap/tidb:v8.5.5
// answers 1064 for ALGORITHM = INPLACE ALGORITHM = COPY.
//
// Ref: opt_index_lock_and_algorithm in MySQL 8.0 sql/sql_yacc.yy.
// https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
func (p *Parser) parseIndexLockAndAlgorithm() (algorithm string, lock string, err error) {
	var haveAlgorithm, haveLock bool
	for {
		switch p.cur.Type {
		case kwALGORITHM:
			if haveAlgorithm {
				return "", "", p.syntaxErrorAtCur()
			}
			p.advance()
			p.match('=') // optional =
			if algorithm, err = p.parseAlgorithmValue(); err != nil {
				return "", "", err
			}
			haveAlgorithm = true
		case kwLOCK:
			if haveLock {
				return "", "", p.syntaxErrorAtCur()
			}
			p.advance()
			p.match('=') // optional =
			if lock, err = p.parseLockValue(); err != nil {
				return "", "", err
			}
			haveLock = true
		default:
			return algorithm, lock, nil
		}
	}
}

// parseAlgorithmValue reads the value of an AlgorithmClause: DEFAULT, or one of the
// words COPY, INPLACE and INSTANT, each a token of its own. An identifier is any other
// name and the arm for it answers 1800, ErrUnknownAlterAlgorithm, so a backquoted
// `inplace` is refused as well as a name that is no algorithm. A reserved word is not
// an identifier, and answers 1064.
//
// Ref: AlgorithmClause in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseAlgorithmValue() (string, error) {
	if tok, ok := p.match(kwDEFAULT); ok {
		return tok.Str, nil
	}
	if p.isBareWord("COPY") || p.isBareWord("INPLACE") || p.isBareWord("INSTANT") {
		return p.advance().Str, nil
	}
	start := p.cur.Loc
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	return "", &ParseError{Message: "unknown ALGORITHM '" + name + "'", Position: start}
}

// parseLockValue reads the value of a LockClause: DEFAULT, or an Identifier naming
// NONE, SHARED or EXCLUSIVE in any case. Unlike ALGORITHM's words these are read from a
// name, so a backquoted `none` runs, and the action answers 1801, ErrUnknownAlterLock,
// for any other name, a backquoted `default` included.
//
// Ref: LockClause in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseLockValue() (string, error) {
	if tok, ok := p.match(kwDEFAULT); ok {
		return tok.Str, nil
	}
	start := p.cur.Loc
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	switch strings.ToUpper(name) {
	case "NONE", "SHARED", "EXCLUSIVE":
		return name, nil
	}
	return "", &ParseError{Message: "unknown LOCK '" + name + "'", Position: start}
}
