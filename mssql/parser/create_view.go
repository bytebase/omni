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
	if name == nil {
		return nil, p.syntaxErrorAtCur()
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
	query, err := p.parseSelectStmtWithCTE()
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
			// OPTION (viewStatementBody, TSql160.g: With Check Option)
			if _, err := p.expect(kwOPTION); err != nil {
				return nil, err
			}
			stmt.WithCheck = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateMaterializedViewStmt refuses CREATE MATERIALIZED VIEW AS SELECT,
// with CREATE MATERIALIZED VIEW consumed: it is Azure Synapse Analytics
// syntax, which ScriptDOM reads (createMaterializedViewStatement, TSql160.g)
// and SQL Server 2022 refuses (Msg 102 "Incorrect syntax near
// 'CREATE MATERIALIZED VIEW'").
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-materialized-view-as-select-transact-sql
func (p *Parser) parseCreateMaterializedViewStmt() (*nodes.CreateMaterializedViewStmt, error) {
	return nil, p.newParseError(p.prev.Loc, "CREATE MATERIALIZED VIEW is Azure Synapse Analytics syntax, which SQL Server refuses")
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
// Only ENCRYPTION, SCHEMABINDING, and VIEW_METADATA are valid, each at most
// once (Msg 1039), and a comma is followed by another option.
func (p *Parser) parseViewOptionList() (*nodes.List, error) {
	var items []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		opt := p.parseViewOption()
		if opt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if seen[opt.Str] {
			return nil, p.newParseError(optLoc, "option "+opt.Str+" is specified more than once")
		}
		seen[opt.Str] = true
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
	// VIEW_METADATA as an identifier that is not delimited
	if p.ddlIsWord("VIEW_METADATA") {
		p.advance()
		return &nodes.String{Str: "VIEW_METADATA"}
	}
	return nil
}

// parseAlterMaterializedViewStmt refuses ALTER MATERIALIZED VIEW, with ALTER
// MATERIALIZED VIEW consumed: it is Azure Synapse Analytics syntax, which SQL
// Server 2022 refuses (Msg 102).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-materialized-view-transact-sql
func (p *Parser) parseAlterMaterializedViewStmt() (*nodes.AlterMaterializedViewStmt, error) {
	return nil, p.newParseError(p.prev.Loc, "ALTER MATERIALIZED VIEW is Azure Synapse Analytics syntax, which SQL Server refuses")
}
