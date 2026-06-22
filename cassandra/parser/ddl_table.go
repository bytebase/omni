package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateTable() (*ast.CreateTableStmt, error) {
	start := p.curLoc()
	p.advance() // TABLE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	stmt := &ast.CreateTableStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	// Parse column definitions and optional trailing PRIMARY KEY.
	for {
		if p.cur.Type == tokPRIMARY {
			pk, err := p.parsePrimaryKeyElement()
			if err != nil {
				return nil, err
			}
			stmt.PrimaryKey = pk
			break
		}
		col, err := p.parseColumnDefinition()
		if err != nil {
			return nil, err
		}
		stmt.Columns = append(stmt.Columns, col)
		if !p.match(tokCOMMA) {
			break
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	// Parse optional WITH clause.
	if p.match(tokWITH) {
		if err := p.parseTableOptions(stmt); err != nil {
			return nil, err
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseColumnDefinition() (*ast.ColumnDef, error) {
	start := p.curLoc()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col := &ast.ColumnDef{Name: name, Type: dt}

	if p.cur.Type == tokSTATIC {
		col.Static = true
		p.advance()
	}
	if p.cur.Type == tokPRIMARY && p.peekNext().Type == tokKEY {
		col.PrimaryKey = true
		p.advance() // PRIMARY
		p.advance() // KEY
	}

	col.Loc = p.makeLoc(start)
	return col, nil
}

func (p *Parser) parsePrimaryKeyElement() (*ast.PrimaryKeyDef, error) {
	start := p.curLoc()
	p.advance() // PRIMARY
	if err := p.expectKeyword(tokKEY); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	pk := &ast.PrimaryKeyDef{}

	// Composite key: ((pk1, pk2), ck1, ck2)
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
		// Single or compound key: pk, ck1, ck2
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		pk.PartitionKeys = append(pk.PartitionKeys, col)
	}

	// Clustering keys
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

	pk.Loc = p.makeLoc(start)
	return pk, nil
}

func (p *Parser) parseTableOptions(stmt *ast.CreateTableStmt) error {
	for {
		// Check for COMPACT STORAGE
		if p.cur.Type == tokCOMPACT {
			p.advance()
			if err := p.expectKeyword(tokSTORAGE); err != nil {
				return err
			}
			stmt.CompactStorage = true
		} else if p.cur.Type == tokCLUSTERING {
			co, err := p.parseClusteringOrderClause()
			if err != nil {
				return err
			}
			stmt.ClusteringOrders = co
		} else {
			opt, err := p.parseTableOption()
			if err != nil {
				return err
			}
			stmt.Options = append(stmt.Options, opt)
		}
		if !p.match(tokAND) {
			break
		}
	}
	return nil
}

func (p *Parser) parseTableOption() (*ast.TableOption, error) {
	start := p.curLoc()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}
	var value ast.ExprNode
	if p.cur.Type == tokLBRACE {
		hash, err := p.parseOptionHash()
		if err != nil {
			return nil, err
		}
		value = hash
	} else {
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		value = val
	}
	return &ast.TableOption{Name: name, Value: value, Loc: p.makeLoc(start)}, nil
}

func (p *Parser) parseClusteringOrderClause() ([]*ast.ClusteringOrder, error) {
	p.advance() // CLUSTERING
	if err := p.expectKeyword(tokORDER); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokBY); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	var orders []*ast.ClusteringOrder
	for {
		start := p.curLoc()
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		dir := ""
		if p.cur.Type == tokASC {
			dir = "ASC"
			p.advance()
		} else if p.cur.Type == tokDESC {
			dir = "DESC"
			p.advance()
		}
		orders = append(orders, &ast.ClusteringOrder{Column: col, Direction: dir, Loc: p.makeLoc(start)})
		if !p.match(tokCOMMA) {
			break
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return orders, nil
}

func (p *Parser) parseAlterTable() (*ast.AlterTableStmt, error) {
	start := p.curLoc()
	p.advance() // TABLE

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterTableStmt{Name: name}

	switch p.cur.Type {
	case tokADD:
		stmt.Op = ast.AlterTableAdd
		p.advance()
		for {
			col, err := p.parseColumnDefinition()
			if err != nil {
				return nil, err
			}
			stmt.AddColumns = append(stmt.AddColumns, col)
			if !p.match(tokCOMMA) {
				break
			}
		}
	case tokDROP:
		p.advance()
		if p.cur.Type == tokCOMPACT {
			stmt.Op = ast.AlterTableDropCompactStorage
			p.advance() // COMPACT
			if err := p.expectKeyword(tokSTORAGE); err != nil {
				return nil, err
			}
		} else {
			stmt.Op = ast.AlterTableDrop
			for {
				col, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.DropColumns = append(stmt.DropColumns, col)
				if !p.match(tokCOMMA) {
					break
				}
			}
		}
	case tokRENAME:
		stmt.Op = ast.AlterTableRename
		p.advance()
		from, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokTO); err != nil {
			return nil, err
		}
		to, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.RenameFrom = from
		stmt.RenameTo = to
	case tokWITH:
		stmt.Op = ast.AlterTableWith
		p.advance()
		for {
			opt, err := p.parseTableOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			if !p.match(tokAND) {
				break
			}
		}
	default:
		return nil, p.errorf("expected ADD, DROP, RENAME, or WITH after ALTER TABLE, got %s", p.tokenDesc())
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropTable() (*ast.DropTableStmt, error) {
	start := p.curLoc()
	p.advance() // TABLE

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropTableStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}
