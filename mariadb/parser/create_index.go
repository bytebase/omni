package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
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
func (p *Parser) parseCreateIndexStmt(unique bool, fulltext bool, spatial bool) (*nodes.CreateIndexStmt, error) {
	start := p.pos()
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

	// Optional index_type: USING {BTREE | HASH}
	if p.cur.Type == kwUSING {
		p.advance()
		typeName, _, err := p.parseKeywordOrIdent()
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

	// (key_part, ...) — WITHOUT OVERLAPS allowed only on a UNIQUE index.
	cols, err := p.parseParenIndexKeyParts(stmt.Unique)
	if err != nil {
		return nil, err
	}
	stmt.Columns = cols

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
	for {
		if p.cur.Type == kwALGORITHM {
			p.advance()
			p.match('=') // optional =
			val, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return nil, err
			}
			stmt.Algorithm = val
		} else if p.cur.Type == kwLOCK {
			p.advance()
			p.match('=') // optional =
			val, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return nil, err
			}
			stmt.Lock = val
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
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
			Loc:    nodes.Loc{Start: start, End: p.pos()},
			Column: colName,
		}

		// Optional (length) — when present, the length is required (no empty "()").
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			col.Length = int(p.cur.Ival)
			col.HasPrefix = true
			p.advance()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}

	// Optional ASC | DESC
	hasOrdering := false
	if _, ok := p.match(kwASC); ok {
		hasOrdering = true // ASC is the default, but the token was given
	} else if _, ok := p.match(kwDESC); ok {
		col.Desc = true
		hasOrdering = true
	}

	// Optional WITHOUT OVERLAPS — valid only on a bare column key part: no
	// functional expression, no prefix length, no ordering token (else 1064).
	// OVERLAPS is non-reserved (matched by text).
	if p.cur.Type == kwWITHOUT && p.peekNext().Type == kwOVERLAPS {
		if col.Functional || col.Length > 0 || hasOrdering {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // WITHOUT
		p.advance() // OVERLAPS
		col.WithoutOverlaps = true
	}

	col.Loc.End = p.pos()
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

// constrAllowsOverlaps reports whether a table constraint type permits a
// WITHOUT OVERLAPS key part (UNIQUE / PRIMARY KEY only).
func constrAllowsOverlaps(t nodes.ConstraintType) bool {
	return t == nodes.ConstrPrimaryKey || t == nodes.ConstrUnique
}

// parseParenIndexKeyParts parses a parenthesized list of index key parts:
//
//	(key_part [, key_part] ...)
//
// allowOverlaps is true only for UNIQUE / PRIMARY KEY definitions, where a
// WITHOUT OVERLAPS application-time period key part is valid.
func (p *Parser) parseParenIndexKeyParts(allowOverlaps bool) ([]*nodes.IndexColumn, error) {
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
	if err := p.validateOverlapsParts(cols, allowOverlaps); err != nil {
		return nil, err
	}
	return cols, nil
}

// validateOverlapsParts enforces the WITHOUT OVERLAPS grammar rules (1064): it
// is valid only on a UNIQUE/PRIMARY KEY, appearing exactly once as the last key
// part, alongside at least one ordinary key part.
func (p *Parser) validateOverlapsParts(cols []*nodes.IndexColumn, allowOverlaps bool) error {
	overlaps, ordinary := 0, 0
	for _, c := range cols {
		if c.WithoutOverlaps {
			overlaps++
		} else {
			ordinary++
		}
	}
	if overlaps == 0 {
		return nil
	}
	if !allowOverlaps || overlaps > 1 || ordinary == 0 || !cols[len(cols)-1].WithoutOverlaps {
		return p.syntaxErrorAtCur()
	}
	return nil
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
		p.match('=') // optional =
		if p.cur.Type == tokICONST {
			val := p.cur.Ival
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "KEY_BLOCK_SIZE",
				Value: &nodes.IntLit{Loc: nodes.Loc{Start: start}, Value: val},
			}, true, nil
		}
		// Could be an identifier value
		v, _, err := p.parseIdentifier()
		if err != nil {
			return nil, false, err
		}
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.pos()},
			Name:  "KEY_BLOCK_SIZE",
			Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: v},
		}, true, nil

	case p.cur.Type == kwUSING:
		p.advance()
		typeName, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, false, err
		}
		return &nodes.IndexOption{
			Loc:   nodes.Loc{Start: start, End: p.pos()},
			Name:  "USING",
			Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: typeName},
		}, true, nil

	case p.cur.Type == kwWITH:
		next := p.peekNext()
		if next.Type == kwPARSER {
			p.advance() // consume WITH
			p.advance() // consume PARSER
			parserName, _, err := p.parseIdent()
			if err != nil {
				return nil, false, err
			}
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "PARSER",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: parserName},
			}, true, nil
		}
		return nil, false, nil

	case p.cur.Type == kwCOMMENT:
		p.advance()
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "COMMENT",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: val},
			}, true, nil
		}
		return nil, false, &ParseError{Message: "expected string after COMMENT", Position: p.cur.Loc}

	case p.cur.Type == kwVISIBLE:
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Name: "VISIBLE",
		}, true, nil

	case p.cur.Type == kwINVISIBLE:
		p.advance()
		return &nodes.IndexOption{
			Loc:  nodes.Loc{Start: start, End: p.pos()},
			Name: "INVISIBLE",
		}, true, nil

	case p.cur.Type == kwENGINE_ATTRIBUTE:
		p.advance()
		p.match('=')
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.IndexOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "ENGINE_ATTRIBUTE",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: val},
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
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "SECONDARY_ENGINE_ATTRIBUTE",
				Value: &nodes.StringLit{Loc: nodes.Loc{Start: start}, Value: val},
			}, true, nil
		}
		return nil, false, &ParseError{Message: "expected string after SECONDARY_ENGINE_ATTRIBUTE", Position: p.cur.Loc}
	}

	return nil, false, nil
}
