package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseCreateTable parses createTableStatement (Parser.g:779) and
// tableDefinition (Parser.g:786):
//
//	CREATE TABLE [IF NOT EXISTS] t '(' tableColumns (',' tableColumns?)* ')'
//	[WITH tableProperty (AND tableProperty)*]
//
// A trailing ',' is allowed, and PRIMARY KEY (...) may stand among the
// columns. The actions refuse a column declared twice and a second primary
// key (CreateTableStatement.Raw.addColumn, setPartitionKeyColumns).
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

	columns := map[string]bool{}
	hasPK := false
	for first := true; ; first = false {
		if !first && (p.cur.Type == tokCOMMA || p.cur.Type == tokRPAREN) {
			// ',' tableColumns? : an empty element.
		} else if p.cur.Type == tokPRIMARY {
			pkTok := p.cur
			pk, err := p.parsePrimaryKeyElement()
			if err != nil {
				return nil, err
			}
			if hasPK {
				return nil, p.errorAt(pkTok, "Multiple PRIMARY KEY specified for table '%s' (exactly one required)", name.Parts[len(name.Parts)-1].Name)
			}
			hasPK = true
			stmt.PrimaryKey = pk
		} else {
			colTok := p.cur
			col, err := p.parseColumnDefinition(true)
			if err != nil {
				return nil, err
			}
			if columns[identKey(col.Name)] {
				return nil, p.errorAt(colTok, "Duplicate column '%s' declaration for table '%s'", identKey(col.Name), name.Parts[len(name.Parts)-1].Name)
			}
			columns[identKey(col.Name)] = true
			if col.PrimaryKey {
				if hasPK {
					return nil, p.errorAt(colTok, "Multiple PRIMARY KEY specified for table '%s' (exactly one required)", name.Parts[len(name.Parts)-1].Name)
				}
				hasPK = true
			}
			stmt.Columns = append(stmt.Columns, col)
		}
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

// parseColumnDefinition parses ident comparatorType [STATIC] [columnMask],
// then [PRIMARY KEY] where a CREATE TABLE column takes one (Parser.g:791).
func (p *Parser) parseColumnDefinition(primaryKey bool) (*ast.ColumnDef, error) {
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
	if p.curIs("masked") {
		mask, err := p.parseColumnMask()
		if err != nil {
			return nil, err
		}
		col.Mask = mask
	}
	if primaryKey && p.cur.Type == tokPRIMARY && p.peekNext().Type == tokKEY {
		col.PrimaryKey = true
		p.advance() // PRIMARY
		p.advance() // KEY
	}

	col.Loc = p.makeLoc(start)
	return col, nil
}

// parseColumnMask parses columnMask (Parser.g:798):
//
//	MASKED WITH DEFAULT | MASKED WITH functionName '(' [term (',' term)*] ')'
func (p *Parser) parseColumnMask() (*ast.ColumnMask, error) {
	start := p.curLoc()
	p.advance() // MASKED
	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	mask := &ast.ColumnMask{}
	// DEFAULT is an unreserved keyword: DEFAULT(...) calls a function of
	// that name.
	if p.cur.Type == tokDEFAULT && p.peekNext().Type != tokLPAREN && p.peekNext().Type != tokDOT {
		p.advance()
		mask.Default = true
		mask.Loc = p.makeLoc(start)
		return mask, nil
	}
	fn, err := p.parseFunctionQualifiedName()
	if err != nil {
		return nil, err
	}
	mask.Function = fn
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}
	if p.cur.Type != tokRPAREN {
		for {
			arg, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			mask.Args = append(mask.Args, arg)
			if !p.match(tokCOMMA) {
				break
			}
		}
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	mask.Loc = p.makeLoc(start)
	return mask, nil
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

// parseTableOptions parses tableProperty (AND tableProperty)* (Parser.g:815):
// a property, COMPACT STORAGE, or CLUSTERING ORDER BY (...). CLUSTERING and
// COMPACT are unreserved keywords, so each is a property name unless ORDER
// or STORAGE follows it.
func (p *Parser) parseTableOptions(stmt *ast.CreateTableStmt) error {
	props := map[string]bool{}
	orders := map[string]bool{}
	for {
		switch {
		case p.cur.Type == tokCOMPACT && p.peekNext().Type == tokSTORAGE:
			p.advance() // COMPACT
			p.advance() // STORAGE
			stmt.CompactStorage = true
		case p.cur.Type == tokCLUSTERING && p.peekNext().Type == tokORDER:
			co, err := p.parseClusteringOrderClause(orders)
			if err != nil {
				return err
			}
			stmt.ClusteringOrders = append(stmt.ClusteringOrders, co...)
		default:
			opt, err := p.parseProperty(props)
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

// parseClusteringOrderClause parses CLUSTERING ORDER BY '(' ident (ASC |
// DESC) (',' ident (ASC | DESC))* ')' (Parser.g:818, Parser.g:821). seen
// holds the columns ordered so far in the same WITH; one given twice is
// refused (CreateTableStatement.Raw.extendClusteringOrder).
func (p *Parser) parseClusteringOrderClause(seen map[string]bool) ([]*ast.ClusteringOrder, error) {
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
		colTok := p.cur
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		var dir string
		switch p.cur.Type {
		case tokASC:
			dir = "ASC"
		case tokDESC:
			dir = "DESC"
		default:
			return nil, p.errorf("expected ASC or DESC, got %s", p.tokenDesc())
		}
		p.advance()
		if seen[identKey(col)] {
			return nil, p.errorAt(colTok, "Duplicate column '%s' in CLUSTERING ORDER BY clause", identKey(col))
		}
		seen[identKey(col)] = true
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

// parseAlterTable parses alterTableStatement (Parser.g:957).
func (p *Parser) parseAlterTable() (*ast.AlterTableStmt, error) {
	start := p.curLoc()
	p.advance() // TABLE

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterTableStmt{IfExists: ifExists, Name: name}

	switch p.cur.Type {
	case tokALTER:
		// ALTER col TYPE type | ALTER [IF EXISTS] col (columnMask | DROP MASKED)
		p.advance()
		stmt.MaskIfExists = p.parseIfExists()
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.AlterColumn = col
		switch {
		case p.cur.Type == tokTYPE && !stmt.MaskIfExists:
			p.advance()
			stmt.Op = ast.AlterTableAlterType
			dt, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.AlterType = dt
		case p.curIs("masked"):
			stmt.Op = ast.AlterTableMask
			mask, err := p.parseColumnMask()
			if err != nil {
				return nil, err
			}
			stmt.Mask = mask
		case p.cur.Type == tokDROP:
			p.advance()
			if !p.curIs("masked") {
				return nil, p.errorf("expected MASKED, got %s", p.tokenDesc())
			}
			p.advance()
			stmt.Op = ast.AlterTableMask
			stmt.DropMask = true
		default:
			return nil, p.errorf("expected TYPE, MASKED or DROP MASKED, got %s", p.tokenDesc())
		}
	case tokADD:
		// ADD [IF NOT EXISTS] (column | '(' column (',' column)* ')')
		stmt.Op = ast.AlterTableAdd
		p.advance()
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		stmt.AddIfNotExists = ine
		paren := p.match(tokLPAREN)
		for {
			col, err := p.parseColumnDefinition(false)
			if err != nil {
				return nil, err
			}
			stmt.AddColumns = append(stmt.AddColumns, col)
			if !paren || !p.match(tokCOMMA) {
				break
			}
		}
		if paren {
			if _, err := p.expect(tokRPAREN); err != nil {
				return nil, err
			}
		}
	case tokDROP:
		p.advance()
		if p.cur.Type == tokCOMPACT && p.peekNext().Type == tokSTORAGE {
			stmt.Op = ast.AlterTableDropCompactStorage
			p.advance() // COMPACT
			p.advance() // STORAGE
			break
		}
		// DROP [IF EXISTS] (col | '(' col (',' col)* ')') [USING TIMESTAMP n]
		stmt.Op = ast.AlterTableDrop
		stmt.DropIfExists = p.parseIfExists()
		paren := p.match(tokLPAREN)
		for {
			col, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.DropColumns = append(stmt.DropColumns, col)
			if !paren || !p.match(tokCOMMA) {
				break
			}
		}
		if paren {
			if _, err := p.expect(tokRPAREN); err != nil {
				return nil, err
			}
		}
		if p.match(tokUSING) {
			if err := p.expectKeyword(tokTIMESTAMP); err != nil {
				return nil, err
			}
			// An INTEGER takes a sign written against it (Lexer.g:337).
			start, sign := p.curLoc(), ""
			if p.cur.Type == tokMINUS && p.peekNext().Type == tokINTEGER && p.peekNext().Loc == p.cur.End {
				sign = "-"
				p.advance()
			}
			if p.cur.Type != tokINTEGER {
				return nil, p.errorf("expected integer, got %s", p.tokenDesc())
			}
			stmt.DropTimestamp = &ast.IntegerLit{Val: sign + p.cur.Str, Loc: ast.Loc{Start: start, End: p.cur.End}}
			p.advance()
		}
	case tokRENAME:
		stmt.Op = ast.AlterTableRename
		p.advance()
		stmt.RenameIfExists = p.parseIfExists()
		for {
			rStart := p.curLoc()
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
			stmt.RenameItems = append(stmt.RenameItems, &ast.AlterTableRenameItem{
				From: from,
				To:   to,
				Loc:  p.makeLoc(rStart),
			})
			if !p.match(tokAND) {
				break
			}
		}
	case tokWITH:
		stmt.Op = ast.AlterTableWith
		p.advance()
		opts, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	default:
		return nil, p.errorf("expected ALTER, ADD, DROP, RENAME, or WITH after ALTER TABLE, got %s", p.tokenDesc())
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
