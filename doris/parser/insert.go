package parser

import (
	"github.com/bytebase/omni/doris/ast"
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

	// OVERWRITE TABLE or INTO
	if p.cur.Kind == kwOVERWRITE {
		p.advance() // consume OVERWRITE
		stmt.Overwrite = true
		// TABLE keyword is required after OVERWRITE
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
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
	target, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Optional PARTITION(p1, p2, ...) or PARTITION(*)
	if p.cur.Kind == kwPARTITION {
		partitions, star, err := p.parseInsertPartition()
		if err != nil {
			return nil, err
		}
		stmt.Partition = partitions
		stmt.PartitionStar = star
	}

	// Optional WITH LABEL label_name
	if p.cur.Kind == kwWITH {
		// Peek ahead: WITH LABEL means insert label; WITH ident AS means CTE.
		// We need to disambiguate: if WITH is followed by LABEL, it's a label clause.
		// Otherwise leave it for the query source parser.
		if p.peekNext().Kind == kwLABEL {
			p.advance() // consume WITH
			p.advance() // consume LABEL
			label, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Label = label
		}
		// If WITH is NOT followed by LABEL, fall through — it will be parsed as
		// part of the query source (WITH ... SELECT ...).
	}

	// Optional column list: (col1, col2, ...)
	// Appears only when the current token is '(' and it's NOT the start of VALUES
	// or a SELECT. We detect this by checking what follows the closing ')'.
	// Simpler heuristic: if cur='(' and peekNext() is an identifier (not SELECT/
	// WITH/VALUES keyword), it's a column list.
	if p.cur.Kind == int('(') {
		// Peek inside to determine if this is a column list or a VALUES row.
		// A column list always has identifiers; a VALUES row may have expressions
		// but the distinction here is that column lists appear before the source.
		// Heuristic: if what we see is '(' <ident/keyword_as_ident> [, ...] ')'
		// followed by VALUES or SELECT/WITH, it's a column list.
		cols, isColList, err := p.tryParseColumnList()
		if err != nil {
			return nil, err
		}
		if isColList {
			stmt.Columns = cols
		}
	}

	// Source: VALUES (...), (...) | SELECT ... | WITH ... SELECT ...
	switch p.cur.Kind {
	case kwVALUES:
		// The source is a query (DorisParser.g4:144), and a row list is its
		// inlineTable primary (:1486), so it takes the set-op and trailing
		// clauses a query does: INSERT INTO t VALUES(1) UNION SELECT 2.
		values, err := p.parseValuesStmt()
		if err != nil {
			return nil, err
		}
		source, err := p.parseQueryTail(values)
		if err != nil {
			return nil, err
		}
		// A bare row list stays in Values, which is what consumers of an
		// INSERT count rows from; a query that is more than the rows — a set
		// operation, ORDER BY, LIMIT, INTO OUTFILE, a WITH or an alias — is
		// kept as the query it is.
		if plain, ok := source.(*ast.ValuesStmt); ok && plainValues(plain) {
			stmt.Values = valuesRows(plain)
		} else {
			stmt.Query = source
		}

	case kwSELECT:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		result, err := p.parseQueryTail(sel)
		if err != nil {
			return nil, err
		}
		stmt.Query = result

	case kwWITH:
		withStmt, err := p.parseWithSelect()
		if err != nil {
			return nil, err
		}
		result, err := p.parseQueryTail(withStmt)
		if err != nil {
			return nil, err
		}
		stmt.Query = result

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

// plainValues reports whether a VALUES query is just its rows: no WITH clause,
// no trailing clause and no aliased item, all of which InsertStmt.Values cannot
// carry.
func plainValues(v *ast.ValuesStmt) bool {
	if v.With != nil || v.Into != nil || len(v.OrderBy) > 0 || v.Limit != nil || v.Offset != nil {
		return false
	}
	for _, row := range v.Rows {
		for _, item := range row {
			if item.Aliased {
				return false
			}
		}
	}
	return true
}

// valuesRows flattens a plain VALUES query's rows into InsertStmt.Values, which
// carries each item's expression.
func valuesRows(v *ast.ValuesStmt) [][]ast.Node {
	rows := make([][]ast.Node, len(v.Rows))
	for i, row := range v.Rows {
		rows[i] = make([]ast.Node, len(row))
		for j, item := range row {
			rows[i][j] = item.Expr
		}
	}
	return rows
}
