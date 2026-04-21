// Package parser - create_view.go implements T-SQL CREATE VIEW statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateViewStmt parses a CREATE [OR ALTER] VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-view-transact-sql
//
//	CREATE [ OR ALTER ] VIEW [ schema_name . ] view_name [ ( column [ ,...n ] ) ]
//	[ WITH <view_attribute> [ ,...n ] ]
//	AS select_statement
//	[ WITH CHECK OPTION ]
//
//	<view_attribute> ::=
//	{
//	    [ ENCRYPTION ]
//	    [ SCHEMABINDING ]
//	    [ VIEW_METADATA ]
//	}
func (p *Parser) parseCreateViewStmt(orAlter bool) (*nodes.CreateViewStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER VIEW → identifier or view name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("view_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateViewStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			colName, ok := p.parseIdentifier()
			if !ok {
				return nil, p.unexpectedToken()
			}
			return &nodes.String{Str: colName}, nil
		})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Columns = &nodes.List{Items: cols}
	}

	// WITH <view_attribute> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isViewOption(next) {
			p.advance() // consume WITH
			opts, err := p.parseViewOptionList()
			if err != nil {
				return nil, err
			}
			stmt.Options = opts
			// Set SchemaBinding flag for backward compat
			if stmt.Options != nil {
				for _, item := range stmt.Options.Items {
					if s, ok := item.(*nodes.String); ok && s.Str == "SCHEMABINDING" {
						stmt.SchemaBinding = true
					}
				}
			}
		}
	}

	// AS (required in CREATE VIEW syntax)
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Completion: after CREATE/ALTER VIEW v AS → SELECT
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return nil, errCollecting
	}

	// SELECT query
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	// WITH CHECK OPTION
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwCHECK {
			p.advance() // WITH
			p.advance() // CHECK
			// OPTION
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithCheck = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateMaterializedViewStmt parses a CREATE MATERIALIZED VIEW AS SELECT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-materialized-view-as-select-transact-sql
//
//	CREATE MATERIALIZED VIEW [ schema_name. ] materialized_view_name
//	    WITH (
//	      <distribution_option>
//	      [, FOR_APPEND ]
//	    )
//	    AS <select_statement>
//
//	<distribution_option> ::=
//	    {
//	        DISTRIBUTION = HASH ( distribution_column_name [, ...n] )
//	      | DISTRIBUTION = ROUND_ROBIN
//	    }
func (p *Parser) parseCreateMaterializedViewStmt() (*nodes.CreateMaterializedViewStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// WITH ( <distribution_option> [, FOR_APPEND] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			switch p.cur.Type {
			case kwDISTRIBUTION:
				p.advance() // consume DISTRIBUTION
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				switch p.cur.Type {
				case kwHASH:
					p.advance() // consume HASH
					stmt.Distribution = "HASH"
					if _, err := p.expect('('); err != nil {
						return nil, err
					}
					var cols []nodes.Node
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						colName, ok := p.parseIdentifier()
						if !ok {
							break
						}
						cols = append(cols, &nodes.String{Str: colName})
						if _, ok := p.match(','); !ok {
							break
						}
					}
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
					stmt.HashColumns = &nodes.List{Items: cols}
				case kwROUND_ROBIN:
					p.advance() // consume ROUND_ROBIN
					stmt.Distribution = "ROUND_ROBIN"
				default:
					p.parseIdentifier() // consume unknown distribution type
				}
			case kwFOR_APPEND:
				p.advance() // consume FOR_APPEND
				stmt.ForAppend = true
			default:
				p.parseIdentifier() // consume unknown option
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// AS
	p.match(kwAS)

	// SELECT query
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// isViewOption checks if a token is a valid CREATE VIEW option.
// Valid view options per SqlScriptDOM ViewOptionHelper:
//   - ENCRYPTION
//   - SCHEMABINDING
//   - VIEW_METADATA
func (p *Parser) isViewOption(tok Token) bool {
	switch tok.Type {
	case kwENCRYPTION, kwSCHEMABINDING:
		return true
	}
	if tok.Str != "" {
		switch strings.ToUpper(tok.Str) {
		case "ENCRYPTION", "SCHEMABINDING", "VIEW_METADATA":
			return true
		}
	}
	return false
}

// parseViewOptionList parses a comma-separated list of view options.
// Only ENCRYPTION, SCHEMABINDING, and VIEW_METADATA are valid.
func (p *Parser) parseViewOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt := p.parseViewOption()
		if opt == nil {
			break
		}
		items = append(items, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseViewOption parses a single view option.
func (p *Parser) parseViewOption() *nodes.String {
	// SCHEMABINDING as keyword token
	if p.cur.Type == kwSCHEMABINDING {
		p.advance()
		return &nodes.String{Str: "SCHEMABINDING"}
	}
	// ENCRYPTION as keyword token
	if p.cur.Type == kwENCRYPTION {
		p.advance()
		return &nodes.String{Str: "ENCRYPTION"}
	}
	// VIEW_METADATA as identifier
	if p.cur.Type == tokIDENT {
		s := strings.ToUpper(p.cur.Str)
		switch s {
		case "VIEW_METADATA":
			p.advance()
			return &nodes.String{Str: s}
		}
	}
	return nil
}

// parseAlterMaterializedViewStmt parses an ALTER MATERIALIZED VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-materialized-view-transact-sql
//
//	ALTER MATERIALIZED VIEW [ schema_name. ] view_name
//	{
//	    REBUILD | DISABLE
//	}
func (p *Parser) parseAlterMaterializedViewStmt() (*nodes.AlterMaterializedViewStmt, error) {
	loc := p.pos()

	stmt := &nodes.AlterMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// REBUILD | DISABLE
	action, ok := p.parseIdentifier()
	if ok {
		stmt.Action = action
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
