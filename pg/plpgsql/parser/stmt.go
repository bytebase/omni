package parser

import (
	"strings"

	"github.com/bytebase/omni/pg/plpgsql/ast"

	pgparser "github.com/bytebase/omni/pg/parser"
)

// --------------------------------------------------------------------------
// Section 3.1: Variable Assignment
// --------------------------------------------------------------------------

// parseAssignOrCall parses an assignment statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-statements.html#PLPGSQL-STATEMENTS-ASSIGNMENT
//
//	variable { := | = } expression ;
//
// The target can be a simple variable, a record field (rec.field),
// a multi-level field (rec.nested.field), an array element (arr[1]),
// or an array slice (arr[1:3]).
func (p *Parser) parseAssignOrCall() (ast.Node, error) {
	startPos := p.pos()

	// Collect the target (variable name, possibly with field/subscript)
	target := p.identText()
	p.advance()

	// Handle dotted names and subscripts
	for {
		if p.cur.Type == '.' {
			p.advance()
			if p.isIdent() || p.isAnyKeywordAsIdent() {
				target += "." + p.identText()
				p.advance()
			}
		} else if p.cur.Type == '[' {
			// Collect subscript including brackets
			depth := 1
			target += "["
			p.advance()
			for !p.isEOF() && depth > 0 {
				if p.cur.Type == '[' {
					depth++
				} else if p.cur.Type == ']' {
					depth--
				}
				target += p.source[p.cur.Loc:p.cur.End]
				p.advance()
			}
		} else {
			break
		}
	}

	// Expect := or =
	if p.cur.Type == pgparser.COLON_EQUALS {
		p.advance()
	} else if p.cur.Type == '=' {
		p.advance()
	} else {
		return nil, p.errorf("syntax error at or near %q, expected := or =", p.tokenText(p.cur))
	}

	// Collect expression until semicolon
	expr, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLAssign{
		Target: target,
		Expr:   expr,
		Loc:    ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// --------------------------------------------------------------------------
// Section 3.2: RETURN Variants
// --------------------------------------------------------------------------

// parseReturn parses all RETURN variants.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-control-structures.html#PLPGSQL-STATEMENTS-RETURNING
//
//	RETURN [ expression ] ;
//	RETURN NEXT [ expression ] ;
//	RETURN QUERY query ;
//	RETURN QUERY EXECUTE command-string [ USING expression [, ...] ] ;
func (p *Parser) parseReturn() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume RETURN

	// RETURN NEXT
	if p.isKeyword("NEXT") {
		p.advance() // consume NEXT
		return p.parseReturnNext(startPos)
	}

	// RETURN QUERY
	if p.isKeyword("QUERY") {
		p.advance() // consume QUERY
		return p.parseReturnQuery(startPos)
	}

	// Simple RETURN [expr] ;
	// If immediately followed by semicolon, it's a bare RETURN.
	if p.cur.Type == ';' {
		p.advance()
		return &ast.PLReturn{
			Expr: "",
			Loc:  ast.Loc{Start: startPos, End: p.prev.End},
		}, nil
	}

	// RETURN with expression
	expr, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLReturn{
		Expr: expr,
		Loc:  ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseReturnNext parses: RETURN NEXT [expression] ;
// The RETURN and NEXT keywords have already been consumed.
func (p *Parser) parseReturnNext(startPos int) (ast.Node, error) {
	// Bare RETURN NEXT (for OUT params)
	if p.cur.Type == ';' {
		p.advance()
		return &ast.PLReturnNext{
			Expr: "",
			Loc:  ast.Loc{Start: startPos, End: p.prev.End},
		}, nil
	}

	expr, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLReturnNext{
		Expr: expr,
		Loc:  ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseReturnQuery parses: RETURN QUERY {query | EXECUTE expr [USING ...]} ;
// The RETURN and QUERY keywords have already been consumed.
func (p *Parser) parseReturnQuery(startPos int) (ast.Node, error) {
	// RETURN QUERY EXECUTE ...
	if p.isKeyword("EXECUTE") {
		p.advance() // consume EXECUTE
		return p.parseReturnQueryExecute(startPos)
	}

	// RETURN QUERY static-query ;
	query, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLReturnQuery{
		Query: query,
		Loc:   ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseReturnQueryExecute parses: RETURN QUERY EXECUTE expr [USING params] ;
// The RETURN, QUERY, and EXECUTE keywords have already been consumed.
func (p *Parser) parseReturnQueryExecute(startPos int) (ast.Node, error) {
	// Collect expression until USING or semicolon
	dynQuery, err := p.collectUntil("USING")
	if err != nil {
		return nil, err
	}

	var params []string
	if p.isKeyword("USING") {
		p.advance() // consume USING
		params, err = p.parseReturnQueryUsingParams()
		if err != nil {
			return nil, err
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLReturnQuery{
		DynQuery: dynQuery,
		Params:   params,
		Loc:      ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseReturnQueryUsingParams parses comma-separated USING parameter expressions
// until a semicolon at depth 0.
func (p *Parser) parseReturnQueryUsingParams() ([]string, error) {
	var params []string
	for {
		start := p.pos()
		depth := 0
		for !p.isEOF() {
			if p.cur.Type == '(' {
				depth++
				p.advance()
				continue
			}
			if p.cur.Type == ')' {
				if depth > 0 {
					depth--
				}
				p.advance()
				continue
			}
			if depth == 0 {
				if p.cur.Type == ',' {
					break
				}
				if p.cur.Type == ';' {
					break
				}
			}
			p.advance()
		}
		expr := strings.TrimSpace(p.source[start:p.cur.Loc])
		if expr == "" {
			return nil, p.errorf("syntax error at or near %q, expected expression in USING clause", p.tokenText(p.cur))
		}
		params = append(params, expr)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume comma
	}
	return params, nil
}

// --------------------------------------------------------------------------
// Section 3.3: PERFORM and Bare SQL
// --------------------------------------------------------------------------

// parsePerform parses a PERFORM statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-statements.html#PLPGSQL-STATEMENTS-SQL-NORESULT
//
//	PERFORM query ;
//
// PERFORM runs a SELECT query but discards the result.
func (p *Parser) parsePerform() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume PERFORM

	expr, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLPerform{
		Expr: expr,
		Loc:  ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseExecSQL parses an inline SQL statement (INSERT, UPDATE, DELETE, SELECT, MERGE, WITH, IMPORT).
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-statements.html#PLPGSQL-STATEMENTS-SQL-ONEROW
//
//	sql_statement [ INTO [STRICT] target [, ...] ] ;
//
// For SELECT statements, INTO [STRICT] can appear after the select list.
// For INSERT/UPDATE/DELETE/MERGE, INTO [STRICT] applies after RETURNING.
func (p *Parser) parseExecSQL() (ast.Node, error) {
	startPos := p.pos()

	// Collect the entire SQL statement while tracking INTO within it.
	// PL/pgSQL INTO is not the same as SQL's INSERT INTO, so detection is
	// enabled only after a top-level SELECT or top-level RETURNING clause.
	sqlText, into, strict, err := p.collectSQLWithInto()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLExecSQL{
		SQLText: sqlText,
		Into:    into,
		Strict:  strict,
		Loc:     ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// collectSQLWithInto collects an SQL statement text until semicolon,
// while detecting PL/pgSQL INTO [STRICT] target clauses.
// SELECT can place INTO after the select list. Row-returning DML places INTO
// after RETURNING. WITH is handled by observing the top-level command that
// follows the CTE list, while nested CTE query text is ignored by depth.
// The INTO clause is extracted from the SQL text and the remaining SQL is reassembled.
func (p *Parser) collectSQLWithInto() (string, []string, bool, error) {
	start := p.pos()
	depth := 0
	var into []string
	strict := false
	intoStart := -1
	intoEnd := -1

	intoAllowed := false
	foundInto := false

	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}
		if depth == 0 && p.cur.Type == ';' {
			break
		}

		if depth == 0 {
			switch p.cur.Type {
			case pgparser.SELECT, pgparser.RETURNING:
				intoAllowed = true
			case pgparser.INSERT, pgparser.UPDATE, pgparser.DELETE_P, pgparser.MERGE:
				intoAllowed = false
			}
		}

		// Detect PL/pgSQL INTO at depth 0 after SELECT or RETURNING.
		if intoAllowed && depth == 0 && !foundInto && p.isKeyword("INTO") {
			intoStart = p.cur.Loc
			p.advance() // consume INTO

			// Check for STRICT
			if p.isKeyword("STRICT") {
				strict = true
				p.advance()
			}

			// Parse comma-separated target identifiers
			targetStart := len(into)
			for {
				if !p.isIdent() && !p.isAnyKeywordAsIdent() {
					break
				}
				into = append(into, p.identText())
				p.advance()
				if p.cur.Type == ',' {
					// Check if next token is an ident (part of INTO targets)
					// vs. a comma in the SELECT list
					next := p.peekNext()
					nextIsIdent := (next.Type == pgparser.IDENT || (next.Type > 256 && next.Str != "" && next.Type != pgparser.Op))
					if !nextIsIdent {
						break
					}
					p.advance() // consume comma
				} else {
					break
				}
			}
			if len(into) == targetStart {
				return "", nil, false, p.errorf("syntax error at or near %q, expected identifier", p.tokenText(p.cur))
			}
			intoEnd = p.cur.Loc
			foundInto = true
			continue
		}

		p.advance()
	}

	// Build the SQL text
	if foundInto && intoStart >= 0 && intoEnd >= 0 {
		// Remove the INTO ... clause from the SQL text
		before := p.source[start:intoStart]
		after := p.source[intoEnd:p.cur.Loc]
		sqlText := strings.TrimSpace(before + after)
		return sqlText, into, strict, nil
	}

	sqlText := strings.TrimSpace(p.source[start:p.cur.Loc])
	return sqlText, nil, false, nil
}

// Ensure imports are used.
var _ = pgparser.COLON_EQUALS
