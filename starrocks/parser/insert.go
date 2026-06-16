package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// INSERT statement parser (T4.1)
// ---------------------------------------------------------------------------

// parseInsert parses an INSERT INTO or INSERT OVERWRITE TABLE statement.
// The INSERT keyword has NOT yet been consumed when this is called.
//
// Syntax:
//
//	INSERT [INTO | OVERWRITE TABLE] [TEMPORARY] [PARTITION(p1, p2, ...) | PARTITION(*)] table_name
//	    [WITH LABEL label_name]
//	    [(col1, col2, ...)]
//	    { VALUES (expr, ...) [, (...)] | SELECT ... | WITH ... SELECT ... }
func (p *Parser) parseInsert() (*ast.InsertStmt, error) {
	insertTok, err := p.expect(kwINSERT)
	if err != nil {
		return nil, err
	}

	stmt := &ast.InsertStmt{
		Loc: ast.Loc{Start: insertTok.Loc.Start},
	}

	// OVERWRITE [TABLE] or INTO. doris requires TABLE after OVERWRITE; StarRocks
	// omits it (INSERT OVERWRITE t) — accept both (additive).
	if p.cur.Kind == kwOVERWRITE {
		p.advance() // consume OVERWRITE
		stmt.Overwrite = true
		if p.cur.Kind == kwTABLE {
			p.advance() // optional TABLE
		}
	} else if p.cur.Kind == kwINTO {
		p.advance() // consume INTO
	}
	// If neither INTO nor OVERWRITE TABLE: bare INSERT table — still valid in
	// some dialects but not standard Doris; fall through and parse the table name.

	// Optional TEMPORARY keyword (before PARTITION clause)
	if p.cur.Kind == kwTEMPORARY {
		p.advance() // consume TEMPORARY
		stmt.TempPartition = true
	}

	// Optional PARTITION clause (may come before table name in some forms,
	// but in Doris the order is: INSERT INTO [TEMP] table [PARTITION(...)] [WITH LABEL] [(cols)] source)
	// According to the legacy corpus, PARTITION comes after table name.
	// Parse table name first.
	// Target: a FILES(propertyList) sink (StarRocks) or a table name. FILES is
	// not a keyword, and `t (cols)` looks similar, so speculatively parse the
	// property list and fall back to the table-name path if it isn't one.
	if isIdentText(p, "files") && p.peekNext().Kind == int('(') {
		saved := p.save()
		p.advance() // consume FILES
		props, perr := p.parsePropertyList()
		if perr == nil {
			stmt.FileTarget = props
		} else {
			p.restore(saved)
		}
	}

	if stmt.FileTarget == nil {
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target

		// Optional PARTITION(p1, p2, ...) or PARTITION(*) — table targets only.
		if p.cur.Kind == kwPARTITION {
			partitions, star, err := p.parseInsertPartition()
			if err != nil {
				return nil, err
			}
			stmt.Partition = partitions
			stmt.PartitionStar = star
		}
	}

	// Post-target modifiers in any order (grammar: insertLabelOrColumnAliases*).
	// WITH LABEL and BY NAME apply to both table and FILES targets; the column
	// list is table-only. BY NAME and the column list are mutually exclusive.
	// (WITH not followed by LABEL is left for the query source: WITH ... SELECT.)
modifiers:
	for {
		switch {
		case p.cur.Kind == kwWITH && p.peekNext().Kind == kwLABEL:
			p.advance() // consume WITH
			p.advance() // consume LABEL
			label, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Label = label
		case p.cur.Kind == kwBY && p.peekNext().Kind == kwNAME:
			if stmt.Columns != nil {
				return nil, p.syntaxErrorAtCur() // BY NAME and a column list are mutually exclusive
			}
			p.advance() // consume BY
			p.advance() // consume NAME
			stmt.ByName = true
		case p.cur.Kind == int('(') && stmt.FileTarget == nil && !stmt.ByName:
			cols, isColList, err := p.tryParseColumnList()
			if err != nil {
				return nil, err
			}
			if !isColList {
				break modifiers
			}
			stmt.Columns = cols
		default:
			break modifiers
		}
	}

	// Source: VALUES (...), (...) | SELECT ... | WITH ... SELECT ...
	switch p.cur.Kind {
	case kwVALUES:
		rows, err := p.parseInsertValues()
		if err != nil {
			return nil, err
		}
		stmt.Values = rows

	case kwSELECT:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		result, err := p.parseSetOpTail(sel)
		if err != nil {
			return nil, err
		}
		stmt.Query = result

	case kwWITH:
		withStmt, err := p.parseWithSelect()
		if err != nil {
			return nil, err
		}
		stmt.Query = withStmt

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseInsertPartition parses PARTITION(p1, p2, ...) or PARTITION(*).
// The PARTITION keyword has NOT yet been consumed.
// Returns (partitionNames, isStar, error).
func (p *Parser) parseInsertPartition() ([]string, bool, error) {
	p.advance() // consume PARTITION
	if _, err := p.expect(int('(')); err != nil {
		return nil, false, err
	}

	// PARTITION(*) — star means all partitions
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		if _, err := p.expect(int(')')); err != nil {
			return nil, false, err
		}
		return nil, true, nil
	}

	var parts []string
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, false, err
	}
	parts = append(parts, name)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		name, _, err = p.parseIdentifier()
		if err != nil {
			return nil, false, err
		}
		parts = append(parts, name)
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}

	return parts, false, nil
}

// tryParseColumnList attempts to parse a parenthesized column list. It succeeds
// and returns (cols, true, nil) only when the token sequence is:
//
//	'(' ident [, ident]* ')' followed by VALUES / SELECT / WITH
//
// Otherwise it returns (nil, false, nil) and the parser state is unchanged
// (no tokens consumed). This is important: if it is not a column list,
// the '(' must remain for the VALUES parser to consume.
//
// Implementation: we use the two-token lookahead available in the parser to
// detect the pattern. For deeper lookahead we parse tentatively and restore
// state on failure by returning isColList=false — but that is complex. Instead
// we use a structural heuristic: peek inside via peekNext() to see if the first
// token after '(' is an identifier-like token. If the content looks like
// identifiers, we parse the list. If the source type keyword follows the ')',
// it was indeed a column list. If not (e.g., we consumed a VALUES row), we
// report an error because we're committed.
//
// In practice Doris INSERT always puts VALUES/SELECT/WITH right after the
// column list or right after the table name, so the heuristic is reliable:
//   - '(' followed immediately by VALUES/SELECT/WITH inside → not a column list
//     (impossible: VALUES starts outside parens)
//   - '(' ident ',' ... ')' VALUES / SELECT / WITH → column list
func (p *Parser) tryParseColumnList() ([]string, bool, error) {
	// p.cur.Kind == '(' at entry.

	// Peek at the token after '(' to decide.
	next := p.peekNext()

	// If the token after '(' is '*' it's likely VALUES(*) form — not a col list.
	if next.Kind == int('*') {
		return nil, false, nil
	}

	// If the token after '(' is not identifier-like, it can't be a column list.
	if !isIdentifierToken(next.Kind) {
		return nil, false, nil
	}

	// Looks like a column list — parse it optimistically.
	// We are now committed; any syntax error here is a real error.
	p.advance() // consume '('

	var cols []string
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, false, err
	}
	cols = append(cols, name)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		name, _, err = p.parseIdentifier()
		if err != nil {
			return nil, false, err
		}
		cols = append(cols, name)
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}

	// After a column list we expect VALUES, SELECT, or WITH.
	if p.cur.Kind == kwVALUES || p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
		return cols, true, nil
	}

	// Not a valid column list position — this shouldn't happen with well-formed
	// Doris SQL, but return an error rather than silently misparse.
	return nil, false, p.syntaxErrorAtCur()
}

// parseInsertValues parses VALUES (expr, ...) [, (expr, ...)] ...
// The VALUES keyword has NOT yet been consumed.
func (p *Parser) parseInsertValues() ([][]ast.Node, error) {
	p.advance() // consume VALUES

	var rows [][]ast.Node

	row, err := p.parseInsertValueRow()
	if err != nil {
		return nil, err
	}
	rows = append(rows, row)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		row, err = p.parseInsertValueRow()
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

// parseInsertValueRow parses one VALUES row: '(' expr [, expr]* ')'.
// Also handles DEFAULT as a special expression in insert value context.
func (p *Parser) parseInsertValueRow() ([]ast.Node, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var exprs []ast.Node

	expr, err := p.parseInsertValueExpr()
	if err != nil {
		return nil, err
	}
	exprs = append(exprs, expr)

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		expr, err = p.parseInsertValueExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return exprs, nil
}

// parseInsertValueExpr parses one expression in a VALUES row.
// In INSERT context, DEFAULT is a valid standalone "expression".
func (p *Parser) parseInsertValueExpr() (ast.Node, error) {
	if p.cur.Kind == kwDEFAULT {
		tok := p.advance() // consume DEFAULT
		return &ast.Literal{
			Kind:  ast.LitKeyword,
			Value: "DEFAULT",
			Loc:   tok.Loc,
		}, nil
	}
	return p.parseExpr()
}
