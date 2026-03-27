// Package parser - create_type.go implements CREATE TYPE parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateTypeStmt parses a CREATE TYPE statement.
//
// BNF: mssql/parser/bnf/create-type-transact-sql.bnf
//
//	CREATE TYPE [ schema_name. ] type_name
//	{
//	    FROM base_type [ ( precision [ , scale ] ) ] [ NULL | NOT NULL ]
//	  | EXTERNAL NAME assembly_name [ .class_name ]
//	  | AS TABLE ( { <column_definition> | <computed_column_definition>
//	      | <table_constraint> | <table_index> } [ ,...n ] )
//	    [ WITH ( MEMORY_OPTIMIZED = ON ) ]
//	}
//
//	<table_index> ::=
//	    INDEX index_name [ CLUSTERED | NONCLUSTERED ] [ HASH ]
//	        [ WITH ( BUCKET_COUNT = count ) ]
//	        ( column_name [ ASC | DESC ] [ ,...n ] )
//	        [ INCLUDE ( column_name [ ,...n ] ) ]
func (p *Parser) parseCreateTypeStmt() (*nodes.CreateTypeStmt, error) {
	stmt := &nodes.CreateTypeStmt{}

	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}

	switch {
	case p.cur.Type == kwFROM:
		p.advance()
		stmt.BaseType, err = p.parseDataType()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == kwNULL {
			b := true
			stmt.Nullable = &b
			p.advance()
		} else if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				b := false
				stmt.Nullable = &b
				p.advance()
				p.advance()
			}
		}

	case p.cur.Type == kwEXTERNAL:
		p.advance()
		p.matchIdentCI("NAME")
		name, _ := p.parseIdentifier()
		extName := name
		if p.cur.Type == '.' {
			p.advance()
			className, _ := p.parseIdentifier()
			extName += "." + className
		}
		stmt.ExternalName = extName

	case p.cur.Type == kwAS:
		p.advance()
		p.match(kwTABLE)
		if _, err := p.expect('('); err == nil {
			var elements []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
					p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK ||
					p.cur.Type == kwFOREIGN {
					constraint, err := p.parseTableConstraint()
					if err != nil {
						return nil, err
					}
					if constraint != nil {
						elements = append(elements, constraint)
					}
				} else if p.cur.Type == kwINDEX {
					idx, err := p.parseTableTypeIndex()
					if err != nil {
						return nil, err
					}
					if idx != nil {
						elements = append(elements, idx)
					}
				} else {
					col, err := p.parseColumnDef()
					if err != nil {
						return nil, err
					}
					if col != nil {
						elements = append(elements, col)
					}
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(elements) > 0 {
				stmt.TableDef = &nodes.List{Items: elements}
			}
		}
		// [ WITH ( MEMORY_OPTIMIZED = ON ) ]
		if p.cur.Type == kwWITH {
			p.advance()
			if _, err := p.expect('('); err == nil {
				if p.matchIdentCI("MEMORY_OPTIMIZED") {
					p.match('=')
					if p.cur.Type == kwON {
						stmt.MemoryOptimized = true
						p.advance()
					} else {
						p.advance() // OFF or other
					}
				}
				p.match(')')
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseTableTypeIndex parses an INDEX clause within CREATE TYPE AS TABLE.
//
// BNF:
//
//	INDEX index_name [ CLUSTERED | NONCLUSTERED ] [ HASH ]
//	    [ WITH ( BUCKET_COUNT = count ) ]
//	    ( column_name [ ASC | DESC ] [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
func (p *Parser) parseTableTypeIndex() (*nodes.TableTypeIndex, error) {
	idx := &nodes.TableTypeIndex{Loc: nodes.Loc{Start: p.pos()}}

	p.match(kwINDEX) // consume INDEX keyword

	// index_name
	idx.Name, _ = p.parseIdentifier()

	// [ CLUSTERED | NONCLUSTERED ]
	if p.cur.Type == kwCLUSTERED {
		b := true
		idx.Clustered = &b
		p.advance()
	} else if p.cur.Type == kwNONCLUSTERED {
		b := false
		idx.Clustered = &b
		p.advance()
	}

	// [ HASH ]
	if p.matchIdentCI("HASH") {
		idx.Hash = true
	}

	// [ WITH ( BUCKET_COUNT = count ) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect('('); err == nil {
			if p.matchIdentCI("BUCKET_COUNT") {
				p.match('=')
				idx.BucketCount, err = p.parseExpr()
				if err != nil {
					return nil, err
				}
			}
			p.match(')')
		}
	}

	// ( column_name [ ASC | DESC ] [ ,...n ] )
	if _, err := p.expect('('); err == nil {
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col := &nodes.IndexColumn{Loc: nodes.Loc{Start: p.pos()}}
			col.Name, _ = p.parseIdentifier()
			if p.cur.Type == kwASC {
				col.SortDir = nodes.SortAsc
				p.advance()
			} else if p.cur.Type == kwDESC {
				col.SortDir = nodes.SortDesc
				p.advance()
			}
			col.Loc.End = p.prevEnd()
			cols = append(cols, col)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(cols) > 0 {
			idx.Columns = &nodes.List{Items: cols}
		}
	}

	// [ INCLUDE ( column_name [ ,...n ] ) ]
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if _, err := p.expect('('); err == nil {
			var incCols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				name, _ := p.parseIdentifier()
				incCols = append(incCols, &nodes.String{Str: name})
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(incCols) > 0 {
				idx.IncludeCols = &nodes.List{Items: incCols}
			}
		}
	}

	idx.Loc.End = p.prevEnd()
	return idx, nil
}
