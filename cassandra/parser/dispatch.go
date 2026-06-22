package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseCreate dispatches CREATE statements.
func (p *Parser) parseCreate() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // CREATE

	// Handle OR REPLACE for FUNCTION/AGGREGATE.
	orReplace := false
	if p.cur.Type == tokOR {
		p.advance() // OR
		if err := p.expectKeyword(tokREPLACE); err != nil {
			return nil, err
		}
		orReplace = true
	}

	// Handle CUSTOM INDEX.
	if p.cur.Type == tokCUSTOM {
		stmt, err := p.parseCreateIndex()
		if err != nil {
			return nil, err
		}
		stmt.IsCustom = true
		stmt.Loc.Start = start
		return stmt, nil
	}

	switch p.cur.Type {
	case tokKEYSPACE:
		stmt, err := p.parseCreateKeyspace()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTABLE:
		stmt, err := p.parseCreateTable()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokINDEX:
		stmt, err := p.parseCreateIndex()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTYPE:
		stmt, err := p.parseCreateType()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokMATERIALIZED:
		p.advance() // MATERIALIZED
		if err := p.expectKeyword(tokVIEW); err != nil {
			return nil, err
		}
		return p.parseCreateMVInline(start)
	case tokFUNCTION:
		stmt, err := p.parseCreateFunction()
		if err != nil {
			return nil, err
		}
		stmt.OrReplace = orReplace
		stmt.Loc.Start = start
		return stmt, nil
	case tokAGGREGATE:
		stmt, err := p.parseCreateAggregate()
		if err != nil {
			return nil, err
		}
		stmt.OrReplace = orReplace
		stmt.Loc.Start = start
		return stmt, nil
	case tokTRIGGER:
		stmt, err := p.parseCreateTrigger()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokROLE:
		stmt, err := p.parseCreateRole()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokUSER:
		stmt, err := p.parseCreateUser()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	default:
		return nil, p.errorf("expected object type after CREATE, got %s", p.tokenDesc())
	}
}

// parseCreateMVInline handles CREATE MATERIALIZED VIEW when we've already consumed
// CREATE MATERIALIZED VIEW tokens.
func (p *Parser) parseCreateMVInline(start int) (*ast.CreateMVStmt, error) {
	ifNotExists := p.parseIfNotExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokAS); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokSELECT); err != nil {
		return nil, err
	}

	stmt := &ast.CreateMVStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	// Select columns
	for {
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.SelectColumns = append(stmt.SelectColumns, col)
		if !p.match(tokCOMMA) {
			break
		}
	}

	if err := p.expectKeyword(tokFROM); err != nil {
		return nil, err
	}
	from, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.FromTable = from

	// WHERE
	if err := p.expectKeyword(tokWHERE); err != nil {
		return nil, err
	}

	for {
		if !p.isColumnNotNull() {
			break
		}
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		p.advance() // IS
		p.advance() // NOT
		p.advance() // NULL
		stmt.WhereNotNull = append(stmt.WhereNotNull, col)
		if p.cur.Type != tokAND {
			break
		}
		next := p.peekNext()
		if !p.isNextColumnNotNull(next) {
			break
		}
		p.advance() // AND
	}

	if p.match(tokAND) {
		rels, err := p.parseRelationElements()
		if err != nil {
			return nil, err
		}
		stmt.WhereRelations = rels
	}

	// PRIMARY KEY
	if err := p.expectKeyword(tokPRIMARY); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokKEY); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	pk := &ast.PrimaryKeyDef{Loc: ast.Loc{Start: p.curLoc()}}
	if p.cur.Type == tokLPAREN {
		p.advance()
		for {
			col, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			pk.PartitionKeys = append(pk.PartitionKeys, col)
			if !p.match(tokCOMMA) {
				break
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	} else {
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		pk.PartitionKeys = append(pk.PartitionKeys, col)
	}

	for p.match(tokCOMMA) {
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		pk.ClusteringKeys = append(pk.ClusteringKeys, col)
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	pk.Loc = p.makeLoc(pk.Loc.Start)
	stmt.PrimaryKey = pk

	if p.match(tokWITH) {
		if err := p.parseMVOptions(stmt); err != nil {
			return nil, err
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseAlter dispatches ALTER statements.
func (p *Parser) parseAlter() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // ALTER

	switch p.cur.Type {
	case tokKEYSPACE:
		stmt, err := p.parseAlterKeyspace()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTABLE:
		stmt, err := p.parseAlterTable()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTYPE:
		stmt, err := p.parseAlterType()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokMATERIALIZED:
		p.advance() // MATERIALIZED
		if err := p.expectKeyword(tokVIEW); err != nil {
			return nil, err
		}
		stmt, err := p.parseAlterMV()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokROLE:
		stmt, err := p.parseAlterRole()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokUSER:
		stmt, err := p.parseAlterUser()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	default:
		return nil, p.errorf("expected object type after ALTER, got %s", p.tokenDesc())
	}
}

// parseDrop dispatches DROP statements.
func (p *Parser) parseDrop() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // DROP

	switch p.cur.Type {
	case tokKEYSPACE:
		stmt, err := p.parseDropKeyspace()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTABLE:
		stmt, err := p.parseDropTable()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokINDEX:
		stmt, err := p.parseDropIndex()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTYPE:
		stmt, err := p.parseDropType()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokMATERIALIZED:
		p.advance() // MATERIALIZED
		if err := p.expectKeyword(tokVIEW); err != nil {
			return nil, err
		}
		stmt, err := p.parseDropMV()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokFUNCTION:
		stmt, err := p.parseDropFunction()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokAGGREGATE:
		stmt, err := p.parseDropAggregate()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokTRIGGER:
		stmt, err := p.parseDropTrigger()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokROLE:
		stmt, err := p.parseDropRole()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	case tokUSER:
		stmt, err := p.parseDropUser()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	default:
		return nil, p.errorf("expected object type after DROP, got %s", p.tokenDesc())
	}
}

// parseTruncate parses TRUNCATE [TABLE] [keyspace.]table.
func (p *Parser) parseTruncate() (*ast.TruncateStmt, error) {
	start := p.curLoc()
	p.advance() // TRUNCATE

	// Optional TABLE keyword
	p.match(tokTABLE)

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.TruncateStmt{
		Table: name,
		Loc:   p.makeLoc(start),
	}, nil
}

// parseUse parses USE keyspace.
func (p *Parser) parseUse() (*ast.UseStmt, error) {
	start := p.curLoc()
	p.advance() // USE

	ks, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.UseStmt{
		Keyspace: ks,
		Loc:      p.makeLoc(start),
	}, nil
}
