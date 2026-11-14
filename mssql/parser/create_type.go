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
//
// The grammar is createTypeStatement (TSql160.g:10727): FROM takes a type,
// EXTERNAL NAME an assembly with an optional class (assemblyName,
// TSql160.g:33181), and AS TABLE a table definition that is not empty, with
// no named or FOREIGN KEY constraint (Msg 156), and WITH
// ( MEMORY_OPTIMIZED = { ON | OFF } ) (withTypeTableOptions,
// TSql160.g:10815).
func (p *Parser) parseCreateTypeStmt() (*nodes.CreateTypeStmt, error) {
	stmt := &nodes.CreateTypeStmt{}

	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil {
		return nil, p.syntaxErrorAtCur()
	}

	switch {
	case p.cur.Type == kwFROM:
		p.advance()
		stmt.BaseType, err = p.parseDataType()
		if err != nil {
			return nil, err
		}
		if stmt.BaseType == nil {
			return nil, p.syntaxErrorAtCur()
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
		if _, err := p.expect(kwNAME); err != nil {
			return nil, err
		}
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		extName := name
		if p.cur.Type == '.' {
			p.advance()
			className, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			extName += "." + className
		}
		stmt.ExternalName = extName

	case p.cur.Type == kwAS:
		p.advance()
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var elements []nodes.Node
		for {
			switch p.cur.Type {
			case kwPRIMARY, kwUNIQUE, kwCHECK:
				constraint, err := p.parseTableConstraint()
				if err != nil {
					return nil, err
				}
				if constraint == nil {
					return nil, p.syntaxErrorAtCur()
				}
				elements = append(elements, constraint)
			case kwINDEX:
				idx, err := p.parseTableTypeIndex()
				if err != nil {
					return nil, err
				}
				elements = append(elements, idx)
			default:
				col, err := p.parseTypedColumnDef()
				if err != nil {
					return nil, err
				}
				if col == nil {
					return nil, p.syntaxErrorAtCur()
				}
				elements = append(elements, col)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.TableDef = &nodes.List{Items: elements}
		// [ WITH ( MEMORY_OPTIMIZED = { ON | OFF } ) ]
		if p.cur.Type == kwWITH {
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if err := p.ddlExpectWord("MEMORY_OPTIMIZED"); err != nil {
				return nil, err
			}
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			switch p.cur.Type {
			case kwON:
				stmt.MemoryOptimized = true
			case kwOFF:
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseTableTypeIndex parses an INDEX clause within CREATE TYPE AS TABLE.
//
// BNF:
//
//	INDEX index_name [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ] [ HASH ]
//	    [ ( column_name [ ASC | DESC ] [ ,...n ] ) ]
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
//	    [ WITH ( BUCKET_COUNT = count ) ]
//
// The grammar is inlineIndexTableDefinition (TSql160.g, the rule under
// tableElement): the column list is optional and comes before WITH; SQL
// Server refuses WITH before it (Msg 102).
func (p *Parser) parseTableTypeIndex() (*nodes.TableTypeIndex, error) {
	idx := &nodes.TableTypeIndex{Loc: nodes.Loc{Start: p.pos(), End: -1}}

	p.match(kwINDEX) // consume INDEX keyword

	// index_name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	idx.Name = name

	// [ UNIQUE ]
	if p.cur.Type == kwUNIQUE {
		idx.Unique = true
		p.advance()
	}

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
	if p.ddlIsWord("HASH") {
		p.advance()
		idx.Hash = true
	}

	// [ ( column_name [ ASC | DESC ] [ ,...n ] ) ]
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for {
			col := &nodes.IndexColumn{Loc: nodes.Loc{Start: p.pos(), End: -1}}
			var ok bool
			col.Name, ok = p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
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
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		idx.Columns = &nodes.List{Items: cols}
	}

	// [ INCLUDE ( column_name [ ,...n ] ) ]
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var incCols []nodes.Node
		for {
			name, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			incCols = append(incCols, &nodes.String{Str: name})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		idx.IncludeCols = &nodes.List{Items: incCols}
	}

	// [ WITH ( BUCKET_COUNT = count ) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if err := p.ddlExpectWord("BUCKET_COUNT"); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		count, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		idx.BucketCount = count
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	idx.Loc.End = p.prevEnd()
	return idx, nil
}
