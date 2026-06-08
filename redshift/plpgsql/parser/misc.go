package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"

	pgparser "github.com/bytebase/omni/redshift/parser"
)

// Ensure pgparser import is used.
var _ = pgparser.CLOSE

// parseClose parses: CLOSE cursor_var ;
func (p *Parser) parseClose() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume CLOSE

	// Expect cursor variable name
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return nil, p.errorf("syntax error at or near %q, expected cursor variable name", p.tokenText(p.cur))
	}
	cursorVar := p.identText()
	p.advance()

	// Consume trailing semicolon
	if _, err := p.expect(';'); err != nil {
		return nil, err
	}

	return &ast.PLClose{
		CursorVar: cursorVar,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseGetDiag parses: GET [CURRENT|STACKED] DIAGNOSTICS target = ITEM [, target = ITEM ...] ;
func (p *Parser) parseGetDiag() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume GET

	node := &ast.PLGetDiag{}

	// Check for CURRENT or STACKED
	if p.isKeyword("CURRENT") {
		p.advance()
		// node.IsStacked remains false (CURRENT is default)
	} else if p.isKeyword("STACKED") {
		p.advance()
		node.IsStacked = true
	}

	// Expect DIAGNOSTICS
	if err := p.expectKeyword("DIAGNOSTICS"); err != nil {
		return nil, err
	}

	// Parse diagnostic items: target = ITEM [, target = ITEM ...]
	for {
		item, err := p.parseGetDiagItem()
		if err != nil {
			return nil, err
		}
		node.Items = append(node.Items, item)

		// Check for comma (more items)
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		break
	}

	// Consume trailing semicolon
	if _, err := p.expect(';'); err != nil {
		return nil, err
	}

	node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
	return node, nil
}

// parseGetDiagItem parses: target = ITEM  or  target := ITEM
func (p *Parser) parseGetDiagItem() (*ast.PLGetDiagItem, error) {
	startPos := p.pos()

	// Expect target variable name
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return nil, p.errorf("syntax error at or near %q, expected diagnostic target variable", p.tokenText(p.cur))
	}
	target := p.identText()
	p.advance()

	// Expect = or :=
	if p.cur.Type == pgparser.COLON_EQUALS {
		p.advance()
	} else if p.cur.Type == '=' {
		p.advance()
	} else {
		return nil, p.errorf("syntax error at or near %q, expected = or :=", p.tokenText(p.cur))
	}

	// Expect diagnostic item keyword
	kind, err := p.parseDiagItemKind()
	if err != nil {
		return nil, err
	}

	return &ast.PLGetDiagItem{
		Target: target,
		Kind:   kind,
		Loc:    ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseDiagItemKind parses a diagnostic item keyword and returns its canonical uppercase form.
func (p *Parser) parseDiagItemKind() (string, error) {
	// Current diagnostic items
	currentItems := []string{
		"ROW_COUNT", "PG_CONTEXT", "PG_ROUTINE_OID",
	}
	// Stacked diagnostic items
	stackedItems := []string{
		"RETURNED_SQLSTATE", "MESSAGE_TEXT",
		"PG_EXCEPTION_DETAIL", "PG_EXCEPTION_HINT", "PG_EXCEPTION_CONTEXT",
		"COLUMN_NAME", "CONSTRAINT_NAME", "PG_DATATYPE_NAME",
		"TABLE_NAME", "SCHEMA_NAME",
	}

	allItems := append(currentItems, stackedItems...)
	for _, item := range allItems {
		if p.isKeyword(item) {
			p.advance()
			return strings.ToUpper(item), nil
		}
	}

	return "", p.errorf("syntax error at or near %q, expected diagnostic item name", p.tokenText(p.cur))
}

// --------------------------------------------------------------------------
// Section 5.2: ASSERT, CALL, DO, COMMIT, ROLLBACK
// --------------------------------------------------------------------------

// parseAssert parses: ASSERT condition [, message] ;
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-errors-and-messages.html#PLPGSQL-STATEMENTS-ASSERT
//
//	ASSERT condition [ , message ] ;
func (p *Parser) parseAssert() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume ASSERT

	// Collect condition expression until comma or semicolon
	cond, err := p.collectUntilCommaOrSemicolon()
	if err != nil {
		return nil, err
	}

	var message string
	if p.cur.Type == ',' {
		p.advance() // consume comma
		message, err = p.collectUntilSemicolon()
		if err != nil {
			return nil, err
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLAssert{
		Condition: cond,
		Message:   message,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// collectUntilCommaOrSemicolon collects tokens until a comma or semicolon at depth 0.
// Neither the comma nor semicolon is consumed.
func (p *Parser) collectUntilCommaOrSemicolon() (string, error) {
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
			if p.cur.Type == ',' || p.cur.Type == ';' {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				return text, nil
			}
		}
		p.advance()
	}

	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}

// parseCall parses: CALL proc_name(args) ;
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-transactions.html
//
//	CALL proc_name ( [ argument [, ...] ] ) ;
//
// The entire statement text (from CALL to before the semicolon) is collected as SQL.
func (p *Parser) parseCall() (ast.Node, error) {
	startPos := p.pos()

	// Collect everything from CALL to semicolon as SQL text
	sqlText, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLCall{
		SQLText: sqlText,
		Loc:     ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseDo parses: DO $$ body $$ ;
//
// Ref: https://www.postgresql.org/docs/17/sql-do.html
//
//	DO [ LANGUAGE lang_name ] code ;
//
// Collected as SQL text like CALL.
func (p *Parser) parseDo() (ast.Node, error) {
	startPos := p.pos()

	// Collect everything from DO to semicolon as SQL text
	sqlText, err := p.collectUntilSemicolon()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLCall{
		SQLText: sqlText,
		Loc:     ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseCommit parses: COMMIT [AND [NO] CHAIN] ;
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-transactions.html
//
//	COMMIT [ AND [ NO ] CHAIN ] ;
func (p *Parser) parseCommit() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume COMMIT

	chain := p.parseOptionalChain()

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLCommit{
		Chain: chain,
		Loc:   ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseRollback parses: ROLLBACK [AND [NO] CHAIN] ;
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-transactions.html
//
//	ROLLBACK [ AND [ NO ] CHAIN ] ;
func (p *Parser) parseRollback() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume ROLLBACK

	chain := p.parseOptionalChain()

	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLRollback{
		Chain: chain,
		Loc:   ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseOptionalChain parses the optional AND [NO] CHAIN suffix.
// Returns true if AND CHAIN, false otherwise (including AND NO CHAIN).
func (p *Parser) parseOptionalChain() bool {
	if !p.isKeyword("AND") {
		return false
	}
	p.advance() // consume AND

	if p.isKeyword("NO") {
		p.advance() // consume NO
		// Expect CHAIN
		if p.isKeyword("CHAIN") {
			p.advance()
		}
		return false // AND NO CHAIN => chain = false
	}

	// Expect CHAIN
	if p.isKeyword("CHAIN") {
		p.advance()
	}
	return true // AND CHAIN => chain = true
}
