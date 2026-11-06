package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// Called when p.cur is at the INDEX keyword (after consuming CREATE [UNIQUE|FULLTEXT|SPATIAL]).
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

// parseIndexOption parses a single index_option.
//
//	index_option:
//	    KEY_BLOCK_SIZE [=] value
//	    | USING {BTREE | HASH}
//	    | WITH PARSER parser_name
//	    | COMMENT 'string'
//	    | {VISIBLE | INVISIBLE}
func (p *Parser) parseIndexOption() (*nodes.IndexOption, bool, error) {
	start := p.pos()

	switch {
	case p.cur.Type == kwKEY_BLOCK_SIZE:
		p.advance()
		p.matchEqual() // optional =
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
		p.matchEqual()
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
		p.matchEqual()
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
// index_type_name is a keyword list, not an identifier, so a name that is no
// index type and a backquoted `BTREE` are both syntax errors: mysql:8.0
// answers 1064 for USING nosuchtype and for USING `BTREE`.
//
// Ref: index_type_name in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseIndexTypeName() (string, error) {
	if tok, ok := p.match(kwBTREE, kwHASH, kwRTREE); ok {
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseIndexLockAndAlgorithm reads the ALGORITHM and LOCK clauses that close a
// CREATE INDEX or a DROP INDEX statement.
//
// The grammar takes at most one of each, in either order, so a second
// ALGORITHM or LOCK is a syntax error rather than a value that wins: mysql:8.0
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
			p.matchEqual() // optional =
			if algorithm, err = p.parseAlgorithmValue(); err != nil {
				return "", "", err
			}
			haveAlgorithm = true
		case kwLOCK:
			if haveLock {
				return "", "", p.syntaxErrorAtCur()
			}
			p.advance()
			p.matchEqual() // optional =
			if lock, err = p.parseLockValue(); err != nil {
				return "", "", err
			}
			haveLock = true
		default:
			return algorithm, lock, nil
		}
	}
}

// parseAlgorithmValue reads the value of ALGORITHM: the DEFAULT keyword, or an
// identifier naming COPY, INPLACE or INSTANT in any case, backquoted or not. Any other
// identifier is 1800, ErrUnknownAlterAlgorithm, from the action, a backquoted `default`
// and NOCOPY among them. The value is not optional, and a string, a number and a
// reserved word are no identifier: mysql:8.0 answers 1064 for ALGORITHM = 'inplace', for
// ALGORITHM = 1 and for ALGORITHM = select.
//
// Ref: alter_algorithm_option_value in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseAlgorithmValue() (string, error) {
	if tok, ok := p.match(kwDEFAULT); ok {
		return tok.Str, nil
	}
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	switch strings.ToUpper(name) {
	case "COPY", "INPLACE", "INSTANT":
		return name, nil
	}
	return "", &ParseError{Message: "Unknown ALGORITHM '" + name + "'", Position: start}
}

// parseLockValue reads the value of LOCK the same way: the DEFAULT keyword, or an
// identifier naming NONE, SHARED or EXCLUSIVE, where any other identifier is 1801,
// ErrUnknownAlterLock.
//
// Ref: alter_lock_option_value in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseLockValue() (string, error) {
	if tok, ok := p.match(kwDEFAULT); ok {
		return tok.Str, nil
	}
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	switch strings.ToUpper(name) {
	case "NONE", "SHARED", "EXCLUSIVE":
		return name, nil
	}
	return "", &ParseError{Message: "Unknown LOCK TYPE '" + name + "'", Position: start}
}
