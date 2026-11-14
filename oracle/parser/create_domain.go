package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// domainBuiltinTypes are the built-in datatypes a domain column takes.
// Oracle 23ai refuses a user-defined type there, at the word after it
// (CREATE DOMAIN d AS t NULL).
var domainBuiltinTypes = map[string]bool{
	"NUMBER": true, "INTEGER": true, "INT": true, "SMALLINT": true, "DECIMAL": true, "DEC": true,
	"NUMERIC": true, "FLOAT": true, "REAL": true, "DOUBLE": true, "BINARY_FLOAT": true,
	"BINARY_DOUBLE": true, "CHAR": true, "CHARACTER": true, "NCHAR": true, "NATIONAL": true,
	"VARCHAR": true, "VARCHAR2": true, "NVARCHAR2": true, "CLOB": true, "NCLOB": true, "BLOB": true,
	"BFILE": true, "JSON": true, "DATE": true, "TIMESTAMP": true, "INTERVAL": true, "RAW": true,
	"LONG": true, "ROWID": true, "UROWID": true, "BOOLEAN": true, "BOOL": true, "VECTOR": true,
}

// parseDomainDataType reads a domain column's datatype, a built-in one.
func (p *Parser) parseDomainDataType() (*nodes.TypeName, error) {
	if !p.isIdentLike() || p.cur.Type == tokQIDENT {
		return nil, p.syntaxErrorAtCur()
	}
	switch p.cur.Type {
	case kwNULL, kwNOT, kwCHECK, kwDEFAULT, kwORDER:
		return nil, p.syntaxErrorAtCur()
	}
	tn, err := p.parseTypeName()
	if err != nil {
		return nil, err
	}
	if tn == nil || tn.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if s, ok := tn.Names.Items[0].(*nodes.String); !ok || !domainBuiltinTypes[s.Str] {
		return nil, p.syntaxErrorAtCur()
	}
	return tn, nil
}

// parseCreateDomainStmt parses a CREATE DOMAIN statement after DOMAIN.
//
// BNF: oracle/parser/bnf/create-domain.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-domain.html
//
//	CREATE [ OR REPLACE ] [ USECASE ] DOMAIN [ IF NOT EXISTS ] [ schema. ] domain_name
//	    AS { create_single_column_domain | create_multi_column_domain } ;
//
//	create_single_column_domain:
//	    { datatype | ENUM ( enum_list ) } [ STRICT ]
//	    [ DEFAULT [ ON NULL ] default_expression ]
//	    [ { constraint_clause }... ]
//	    [ COLLATE collation_name ]
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//	    [ annotations_clause ]
//	create_multi_column_domain:
//	    ( column_name AS datatype [ annotations_clause ] [, ...] )
//	    [ CHECK ( condition ) [ constraint_state ] ]...
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//	    [ annotations_clause ]
//
// The page's create_flexible_domain is CREATE FLEXIBLE DOMAIN, a statement
// of its own: Oracle 23ai reads AS FLEXIBLE as a datatype (ORA-00922 at
// the DOMAIN after it). A multi-column domain takes only CHECK constraints
// and no COLLATE (ORA-00922 at NULL or COLLATE), and Oracle takes its
// columns without the commas between them.
func (p *Parser) parseCreateDomainStmt(start int, orReplace, usecase bool) (*nodes.CreateDomainStmt, error) {
	stmt := &nodes.CreateDomainStmt{
		OrReplace: orReplace,
		Usecase:   usecase,
		Loc:       nodes.Loc{Start: start},
	}

	// [ IF NOT EXISTS ]: an IF not followed by NOT EXISTS is refused at
	// the IF.
	if p.cur.Type == kwIF {
		ifTok := p.cur
		p.advance() // IF
		if p.cur.Type != kwNOT || p.peekNext().Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance() // NOT
		p.advance() // EXISTS
		stmt.IfNotExists = true
	}

	var err error
	if stmt.Name, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return nil, err
	}

	switch {
	case p.isKeywordStr("ENUM"):
		stmt.DomainType = "ENUM"
		p.advance() // ENUM
		if err := p.parseDomainEnumList(stmt); err != nil {
			return nil, err
		}
		if p.isKeywordStr("STRICT") {
			stmt.Strict = true
			p.advance()
		}
		if err := p.parseDomainProperties(stmt, false); err != nil {
			return nil, err
		}

	case p.cur.Type == '(':
		stmt.DomainType = "MULTI"
		p.advance() // (
		stmt.Columns = &nodes.List{}
		for {
			col := &nodes.DomainColumn{Loc: nodes.Loc{Start: p.pos()}}
			col.Name = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
			if err := p.expectCreateTableToken(kwAS); err != nil {
				return nil, err
			}
			if col.DataType, err = p.parseDomainDataType(); err != nil {
				return nil, err
			}
			if p.isKeywordStr("ANNOTATIONS") {
				if col.Annotations, err = p.parseDomainAnnotations(false); err != nil {
					return nil, err
				}
			}
			col.Loc.End = p.prev.End
			stmt.Columns.Items = append(stmt.Columns.Items, col)
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			if p.cur.Type == ')' {
				break
			}
		}
		p.advance() // )
		if err := p.parseDomainProperties(stmt, true); err != nil {
			return nil, err
		}

	default:
		stmt.DomainType = "SINGLE"
		if stmt.DataType, err = p.parseDomainDataType(); err != nil {
			return nil, err
		}
		if p.isKeywordStr("STRICT") {
			stmt.Strict = true
			p.advance()
		}
		if err := p.parseDomainProperties(stmt, false); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDomainEnumList reads ( enum_list ).
//
//	enum_item_list: name [ = alias [ = alias ]... ] [ = value ]
//
// Oracle 23ai takes an empty list and a trailing comma (ENUM(), ENUM(a,)).
func (p *Parser) parseDomainEnumList(stmt *nodes.CreateDomainStmt) error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	stmt.EnumItems = &nodes.List{}
	for p.cur.Type != ')' {
		item := &nodes.DomainEnumItem{Loc: nodes.Loc{Start: p.pos()}}
		item.Name = p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		for p.cur.Type == '=' {
			p.advance() // =
			if item.Value != nil {
				return p.syntaxErrorAtTok(p.prev)
			}
			if p.isCreateTableNameStart() {
				item.Aliases = append(item.Aliases, p.cur.Str)
				p.advance()
				continue
			}
			value, err := p.parseExpr()
			if err != nil {
				return err
			}
			if value == nil {
				return p.syntaxErrorAtCur()
			}
			item.Value = value
		}
		item.Loc.End = p.prev.End
		stmt.EnumItems.Items = append(stmt.EnumItems.Items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ,
	}
	return p.expectCreateTableToken(')')
}

// parseDomainProperties reads what follows a domain's datatype, or its
// column list when multi.
func (p *Parser) parseDomainProperties(stmt *nodes.CreateDomainStmt, multi bool) error {
	var err error
	if !multi && p.cur.Type == kwDEFAULT {
		p.advance()
		if p.cur.Type == kwON {
			p.advance()
			if err := p.expectCreateTableToken(kwNULL); err != nil {
				return err
			}
			stmt.DefaultOnNull = true
		}
		if stmt.Default, err = p.parseExpr(); err != nil {
			return err
		}
		if stmt.Default == nil {
			return p.syntaxErrorAtCur()
		}
	}
	if stmt.Constraints, err = p.parseDomainConstraints(multi); err != nil {
		return err
	}
	if !multi && p.isKeywordStr("COLLATE") {
		p.advance()
		stmt.Collation = p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	if p.isKeywordStr("DISPLAY") {
		p.advance()
		if stmt.Display, err = p.parseRequiredExpr(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwORDER {
		p.advance()
		if stmt.Order, err = p.parseRequiredExpr(); err != nil {
			return err
		}
	}
	if p.isKeywordStr("ANNOTATIONS") {
		if stmt.Annotations, err = p.parseDomainAnnotations(false); err != nil {
			return err
		}
	}
	return nil
}

// parseDomainConstraints reads a domain's constraint_clauses.
//
//	[ CONSTRAINT constraint_name ]
//	{ NOT NULL | NULL | CHECK ( condition ) }
//	[ constraint_state ]
func (p *Parser) parseDomainConstraints(multi bool) (*nodes.List, error) {
	var items []nodes.Node
	for {
		start := p.pos()
		c := &nodes.DomainConstraint{Loc: nodes.Loc{Start: start}}
		named := false
		if p.cur.Type == kwCONSTRAINT {
			p.advance()
			c.Name = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
			named = true
		}
		switch {
		case !multi && p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
			c.Type = "NOT_NULL"
			p.advance()
			p.advance()
		case !multi && p.cur.Type == kwNULL:
			c.Type = "NULL"
			p.advance()
		case p.cur.Type == kwCHECK:
			c.Type = "CHECK"
			p.advance()
			if err := p.expectCreateTableToken('('); err != nil {
				return nil, err
			}
			var err error
			if c.CheckExpr, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return nil, err
			}
		default:
			if named {
				return nil, p.syntaxErrorAtCur()
			}
			if len(items) == 0 {
				return nil, nil
			}
			return &nodes.List{Items: items}, nil
		}
		var cs constraintState
		if err := p.parseConstraintState(&cs); err != nil {
			return nil, err
		}
		c.Loc.End = p.prev.End
		items = append(items, c)
	}
}

// parseDomainAnnotations reads a domain's annotations_clause.
//
//	ANNOTATIONS ( [ ADD | DROP ] annotation_name [ 'annotation_value' ] [, ...] )
//
// CREATE DOMAIN adds annotations; ALTER DOMAIN adds or drops them.
func (p *Parser) parseDomainAnnotations(alter bool) (*nodes.List, error) {
	p.advance() // ANNOTATIONS
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	result := &nodes.List{}
	for {
		start := p.pos()
		drop := false
		switch {
		case p.cur.Type == kwADD:
			p.advance()
		case alter && p.cur.Type == kwDROP:
			p.advance()
			drop = true
		}
		key := p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return nil, err
		}
		var val string
		if !drop && p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		result.Items = append(result.Items, &nodes.DDLOption{Key: key, Value: val, Loc: nodes.Loc{Start: start, End: p.prev.End}})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ,
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	return result, nil
}

// parseAlterDomainStmt parses an ALTER DOMAIN statement after DOMAIN.
//
// BNF: oracle/parser/bnf/alter-domain.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/alter-domain.html
//
//	ALTER [ USECASE ] DOMAIN [ IF EXISTS ] [ schema. ] domain_name
//	    { ADD DISPLAY display_expression
//	    | MODIFY DISPLAY display_expression
//	    | DROP DISPLAY
//	    | ADD ORDER order_expression
//	    | MODIFY ORDER order_expression
//	    | DROP ORDER
//	    | annotations_clause
//	    } ;
func (p *Parser) parseAlterDomainStmt(start int, usecase bool) (*nodes.AlterDomainStmt, error) {
	stmt := &nodes.AlterDomainStmt{
		Usecase: usecase,
		Loc:     nodes.Loc{Start: start},
	}

	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var err error
	if stmt.Name, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}

	switch {
	case p.cur.Type == kwADD, p.cur.Type == kwMODIFY:
		verb := "ADD"
		if p.cur.Type == kwMODIFY {
			verb = "MODIFY"
		}
		p.advance()
		switch {
		case p.isKeywordStr("DISPLAY"):
			stmt.Action = verb + "_DISPLAY"
			p.advance()
			if stmt.Display, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwORDER:
			stmt.Action = verb + "_ORDER"
			p.advance()
			if stmt.Order, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case p.cur.Type == kwDROP:
		p.advance()
		switch {
		case p.isKeywordStr("DISPLAY"):
			stmt.Action = "DROP_DISPLAY"
		case p.cur.Type == kwORDER:
			stmt.Action = "DROP_ORDER"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	case p.isKeywordStr("ANNOTATIONS"):
		stmt.Action = "ANNOTATIONS"
		if stmt.Annotations, err = p.parseDomainAnnotations(true); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
