// Package parser - name.go implements identifier and qualified name parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// isIdentLike returns true if the current token can be used as an identifier.
// In T-SQL, context keywords can be used as identifiers but Core keywords cannot
// (unless bracket-quoted). Use isAnyKeywordIdent() for positions that accept all keywords.
// The tokIDENT type already includes [bracketed] and "quoted" identifiers.
func (p *Parser) isIdentLike() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	return isContextKeyword(p.cur.Type)
}

// isAnyKeywordIdent returns true if the current token can be used as a label or alias.
// Like PG's isColLabel(), some positions (explicit aliases after AS, column labels)
// should accept ALL keywords including Core keywords.
// This is the old isIdentLike behavior — needed for positions where SQL Server allows
// reserved words as aliases (e.g., SELECT 1 AS select).
func (p *Parser) isAnyKeywordIdent() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	// Any keyword token can be used as a label/alias
	return p.cur.Type >= kwABSENT && p.cur.Str != ""
}

// parseIdentifier consumes and returns the current token as an identifier string.
// It accepts tokIDENT tokens and keywords used as identifiers.
// Returns ("", false) if the current token is not identifier-like.
//
//	identifier = regular_identifier | bracketed_identifier | quoted_identifier | keyword_as_identifier
func (p *Parser) parseIdentifier() (string, bool) {
	if p.cur.Type == tokIDENT {
		name := p.cur.Str
		p.advance()
		return name, true
	}
	if isContextKeyword(p.cur.Type) {
		name := p.cur.Str
		p.advance()
		return name, true
	}
	return "", false
}

// parseTableRef parses a qualified object name: [server.][database.][schema.]object
// Used for table names in DDL/DML contexts (FROM, CREATE TABLE, INSERT INTO, etc.).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/databases/database-identifiers
//
//	qualified_name = [ server_name . [ database_name ] . [ schema_name ] . ]
//	                 | [ database_name . [ schema_name ] . ]
//	                 | [ schema_name . ]
//	                 object_name
func (p *Parser) parseTableRef() (*nodes.TableRef, error) {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, nil
	}

	ref := &nodes.TableRef{
		Object: name,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// Collect dot-separated parts
	parts := []string{name}
	for p.cur.Type == '.' {
		p.advance() // consume .
		// Completion: after dot in qualified name → table_ref (schema-qualified)
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		part, ok := p.parseIdentifier()
		if !ok {
			// Handle trailing dot (e.g., "db..object" means db.dbo.object with empty schema)
			parts = append(parts, "")
			continue
		}
		parts = append(parts, part)
	}

	// Assign parts based on count: object, schema.object, db.schema.object, server.db.schema.object
	switch len(parts) {
	case 1:
		ref.Object = parts[0]
	case 2:
		ref.Schema = parts[0]
		ref.Object = parts[1]
	case 3:
		ref.Database = parts[0]
		ref.Schema = parts[1]
		ref.Object = parts[2]
	default: // 4+
		ref.Server = parts[0]
		ref.Database = parts[1]
		ref.Schema = parts[2]
		ref.Object = parts[3]
	}

	ref.Loc.End = p.prevEnd()
	return ref, nil
}

// parseVariableTableSource parses a T-SQL table variable used as a table source
// in FROM / JOIN positions. Handles three shapes:
//
//	@t                            -> *nodes.TableVarRef
//	@t [AS] alias                 -> *nodes.TableVarRef with Alias
//	@v.Method(args) [alias [(c)]] -> *nodes.TableVarMethodCallRef
//
// Mirrors SqlScriptDOM variableTableReference + variableMethodCallTableReference.
// Caller must ensure p.cur.Type == tokVARIABLE.
func (p *Parser) parseVariableTableSource() (nodes.TableExpr, error) {
	loc := p.pos()
	name := p.cur.Str // includes leading '@'
	p.advance()

	// @v.Method(args) — variable method-call table reference (e.g. XML .nodes()).
	// Any '.' after @name is necessarily a method call (table variables cannot be
	// schema-qualified), so we can commit to this branch on sight of the dot.
	if p.cur.Type == '.' {
		p.advance() // consume .
		method, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var args []nodes.ExprNode
		if p.cur.Type != ')' {
			for {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				if _, ok := p.match(','); !ok {
					break
				}
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		mc := &nodes.TableVarMethodCallRef{
			Var:    name,
			Method: method,
			Args:   args,
			Loc:    nodes.Loc{Start: loc, End: p.prevEnd()},
		}
		mc.Alias = p.parseOptionalAlias()
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				col, ok := p.parseIdentifier()
				if !ok {
					break
				}
				mc.Columns = append(mc.Columns, col)
				if _, ok := p.match(','); !ok {
					break
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		mc.Loc.End = p.prevEnd()
		return mc, nil
	}

	tv := &nodes.TableVarRef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}
	tv.Alias = p.parseOptionalAlias()
	tv.Loc.End = p.prevEnd()
	return tv, nil
}

// parseVariableDmlTarget parses a bare table variable as a DML target
// (INSERT/UPDATE/DELETE/MERGE). Does not accept alias or table hints —
// mirrors SqlScriptDOM variableDmlTarget.
// Caller must ensure p.cur.Type == tokVARIABLE.
func (p *Parser) parseVariableDmlTarget() *nodes.TableVarRef {
	loc := p.pos()
	name := p.cur.Str
	p.advance()
	return &nodes.TableVarRef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
	}
}

// parseIdentExpr parses an identifier expression (column ref, function call, or qualified name).
// This handles both simple identifiers and dot-qualified references.
//
//	ident_expr = identifier [ '(' args ')' ]
//	           | identifier '.' identifier [ '.' identifier [ '.' identifier ] ]
//	           | identifier '.' '*'
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
	loc := p.pos()
	name := p.cur.Str
	p.advance()

	// Function call: ident(...)
	if p.cur.Type == '(' {
		return p.parseFuncCall(name, loc)
	}

	// Static method call: type::Method(args)
	if p.cur.Type == tokCOLONCOLON {
		p.advance() // consume ::
		method := ""
		if p.isIdentLike() {
			method = p.cur.Str
			p.advance()
		}
		mc := &nodes.MethodCallExpr{
			Type:   &nodes.DataType{Name: name, Loc: nodes.Loc{Start: loc, End: -1}},
			Method: method,
			Loc:    nodes.Loc{Start: loc, End: -1},
		}
		if p.cur.Type == '(' {
			p.advance() // consume (
			var args []nodes.Node
			if p.cur.Type != ')' {
				for {
					arg, _ := p.parseExpr()
					args = append(args, arg)
					if _, ok := p.match(','); !ok {
						break
					}
				}
			}
			mc.Args = &nodes.List{Items: args}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		mc.Loc.End = p.prevEnd()
		return mc, nil
	}

	// Qualified name: ident.ident[.ident[.ident]] or ident.*
	if p.cur.Type == '.' {
		return p.parseQualifiedRef(name, loc)
	}

	// Simple column reference
	return &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}, nil
}

// parseQualifiedRef parses a dot-qualified column reference or star expression.
// The first part has already been consumed.
//
//	qualified_ref = first '.' ( '*' | ident [ '.' ( '*' | ident [ '.' ( '*' | ident ) ] ) ] )
func (p *Parser) parseQualifiedRef(first string, loc int) (nodes.ExprNode, error) {
	parts := []string{first}
	for p.cur.Type == '.' {
		p.advance() // consume .

		// Check for table.* or schema.table.*
		if p.cur.Type == '*' {
			p.advance()
			// Build qualifier from collected parts
			qualifier := first
			if len(parts) > 1 {
				qualifier = parts[len(parts)-1]
			}
			return &nodes.StarExpr{
				Qualifier: qualifier,
				Loc:       nodes.Loc{Start: loc, End: -1},
			}, nil
		}

		// Accept identifier or keyword-as-identifier after dot
		if p.isIdentLike() {
			partName := p.cur.Str
			p.advance()

			// Check if this part is followed by '(' -- meaning it's a function call
			// e.g., schema.function(args)
			if p.cur.Type == '(' {
				schema := first
				if len(parts) > 1 {
					schema = parts[0]
				}
				return p.parseFuncCallWithSchema(schema, partName, loc)
			}

			// Check for :: static method call: schema.type::Method(args)
			if p.cur.Type == tokCOLONCOLON {
				p.advance() // consume ::
				method := ""
				if p.isIdentLike() {
					method = p.cur.Str
					p.advance()
				}
				dt := &nodes.DataType{Name: partName, Loc: nodes.Loc{Start: loc, End: -1}}
				if len(parts) > 0 {
					dt.Schema = parts[0]
				}
				mc := &nodes.MethodCallExpr{
					Type:   dt,
					Method: method,
					Loc:    nodes.Loc{Start: loc, End: -1},
				}
				if p.cur.Type == '(' {
					p.advance() // consume (
					var args []nodes.Node
					if p.cur.Type != ')' {
						for {
							arg, _ := p.parseExpr()
							args = append(args, arg)
							if _, ok := p.match(','); !ok {
								break
							}
						}
					}
					mc.Args = &nodes.List{Items: args}
					if _, err := p.expect(')'); err != nil {
						return nil, err
					}
				}
				mc.Loc.End = p.prevEnd()
				return mc, nil
			}

			parts = append(parts, partName)
		} else {
			break
		}
	}

	ref := &nodes.ColumnRef{Loc: nodes.Loc{Start: loc, End: -1}}
	switch len(parts) {
	case 1:
		ref.Column = parts[0]
	case 2:
		ref.Table = parts[0]
		ref.Column = parts[1]
	case 3:
		ref.Schema = parts[0]
		ref.Table = parts[1]
		ref.Column = parts[2]
	case 4:
		ref.Database = parts[0]
		ref.Schema = parts[1]
		ref.Table = parts[2]
		ref.Column = parts[3]
	default: // 5 parts: server.database.schema.table.column
		ref.Server = parts[0]
		ref.Database = parts[1]
		ref.Schema = parts[2]
		ref.Table = parts[3]
		ref.Column = parts[4]
	}
	return ref, nil
}

// parseFuncCallWithSchema parses a schema-qualified function call.
// schema.func(args)
func (p *Parser) parseFuncCallWithSchema(schema, funcName string, loc int) (nodes.ExprNode, error) {
	p.advance() // consume (

	fc := &nodes.FuncCallExpr{
		Name: &nodes.TableRef{Schema: schema, Object: funcName, Loc: nodes.Loc{Start: loc, End: -1}},
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// COUNT(*) special case
	if p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if p.cur.Type == kwOVER {
			fc.Over, _ = p.parseOverClause()
		}
		return fc, nil
	}

	if p.cur.Type == ')' {
		p.advance()
		if p.cur.Type == kwOVER {
			fc.Over, _ = p.parseOverClause()
		}
		return fc, nil
	}

	// Check for DISTINCT
	if _, ok := p.match(kwDISTINCT); ok {
		fc.Distinct = true
	}

	var args []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		arg, _ := p.parseExpr()
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	fc.Args = &nodes.List{Items: args}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Check for OVER clause
	if p.cur.Type == kwOVER {
		fc.Over, _ = p.parseOverClause()
	}

	return fc, nil
}
