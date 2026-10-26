// Package parser - rowset_functions.go implements T-SQL rowset function parsing.
// Rowset functions: OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
// These appear as table sources in FROM clauses.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseRowsetFunction parses OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
// as table sources. These are treated as FuncCallExpr nodes with keyword names.
func (p *Parser) parseRowsetFunction() (nodes.TableExpr, error) {
	loc := p.pos()
	funcName := p.cur.Str
	p.advance() // consume the keyword
	nameEnd := p.prevEnd()

	// OPENDATASOURCE('provider', 'connstr').db.schema.table - special case
	if p.cur.Type == '(' {
		fc := &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc, End: nameEnd}},
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		p.advance() // consume (

		argIndex := 0
		args, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
			if isRowsetObjectArg(funcName, argIndex) && p.isIdentLike() {
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				if ref == nil {
					return nil, p.unexpectedToken()
				}
				argIndex++
				return ref, nil
			}
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.unexpectedToken()
			}
			argIndex++
			return arg, nil
		})
		if err != nil {
			return nil, err
		}
		if len(args) > 0 {
			fc.Args = &nodes.List{Items: args}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Loc.End = p.prevEnd()

		// WITH schema: OPENJSON and OPENXML only. SQL Server rejects WITH
		// after OPENQUERY and OPENROWSET.
		var withClause *nodes.List
		isJSON := strings.EqualFold(funcName, "OPENJSON")
		if (isJSON || strings.EqualFold(funcName, "OPENXML")) && p.cur.Type == kwWITH && p.peekNext().Type == '(' {
			p.advance() // consume WITH
			p.advance() // consume (
			withClause, err = p.parseRowsetWithClause(isJSON)
			if err != nil {
				return nil, err
			}
		}

		alias := p.parseOptionalAlias()
		aliasEnd := p.prevEnd()

		return &nodes.AliasedTableRef{
			Table:   fc,
			Alias:   alias,
			Columns: withClause,
			Loc:     nodes.Loc{Start: loc, End: aliasEnd},
		}, nil
	}

	// Shouldn't happen but handle gracefully
	return &nodes.AliasedTableRef{
		Table: &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc, End: nameEnd}},
			Loc:  nodes.Loc{Start: loc, End: nameEnd},
		},
		Loc: nodes.Loc{Start: loc, End: nameEnd},
	}, nil
}

func isRowsetObjectArg(funcName string, argIndex int) bool {
	switch {
	case strings.EqualFold(funcName, "OPENQUERY"):
		return argIndex == 0
	case strings.EqualFold(funcName, "OPENROWSET"):
		return argIndex == 2
	default:
		return false
	}
}

// parseRowsetWithClause parses the WITH schema of OPENJSON (isJSON) or
// OPENXML. The opening '(' has already been consumed.
//
//	column_name data_type [ COLLATE collation_name ] [ column_path ] [ AS JSON ]
//
// SQL Server requires the data type, and takes the collation only
// undelimited, the path only as a string literal, and AS JSON only in
// OPENJSON, in this order.
func (p *Parser) parseRowsetWithClause(isJSON bool) (*nodes.List, error) {
	cols, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		loc := p.pos()
		colName, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if p.collectMode() {
			p.addRuleCandidate("type_name")
			return nil, errCollecting
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		if dt == nil {
			return nil, p.unexpectedToken()
		}
		col := &nodes.RowsetColumnDef{
			Name:     colName,
			DataType: dt,
			Loc:      nodes.Loc{Start: loc, End: -1},
		}
		if _, ok := p.match(kwCOLLATE); ok {
			if !p.isIdentLike() || p.source[p.cur.Loc] == '[' || p.source[p.cur.Loc] == '"' {
				return nil, p.unexpectedToken()
			}
			col.Collation = p.cur.Str
			p.advance()
		}
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			tok := p.advance()
			col.Path = &nodes.Literal{
				Type:    nodes.LitString,
				Str:     tok.Str,
				IsNChar: tok.Type == tokNSCONST,
				Loc:     nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
			}
		}
		if isJSON && p.cur.Type == kwAS && p.peekNext().Type == kwJSON {
			p.advance() // consume AS
			p.advance() // consume JSON
			col.AsJSON = true
		}
		col.Loc.End = p.prevEnd()
		return col, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: cols}, nil
}
