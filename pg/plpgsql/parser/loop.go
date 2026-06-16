package parser

import (
	"strings"

	"github.com/bytebase/omni/pg/plpgsql/ast"

	pgparser "github.com/bytebase/omni/pg/parser"
)

// parseEndLoop parses END LOOP [label] ; and validates the label.
func (p *Parser) parseEndLoop(label string) error {
	if !p.isKeyword("END") {
		return p.errorf("syntax error at or near %q, expected END LOOP", p.tokenText(p.cur))
	}
	p.advance()

	if !p.isKeyword("LOOP") {
		return p.errorf("syntax error at or near %q, expected LOOP after END", p.tokenText(p.cur))
	}
	p.advance()

	// Check for optional trailing label
	if (p.isIdent() || p.isAnyKeywordAsIdent()) &&
		!p.isKeyword("IF") && !p.isKeyword("LOOP") && !p.isKeyword("CASE") {
		endLabel := p.identText()
		if label == "" || !strings.EqualFold(endLabel, label) {
			if label == "" {
				return p.errorf("end label %q specified for unlabeled loop", endLabel)
			}
			return p.errorf("end label %q does not match loop label %q", endLabel, label)
		}
		p.advance()
	}

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}
	return nil
}

// parseLoop parses: LOOP stmts END LOOP [label] ;
func (p *Parser) parseLoop(label string) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLLoop{
		Label: label,
		Body:  body,
		Loc:   ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseWhile parses: WHILE condition LOOP stmts END LOOP [label] ;
func (p *Parser) parseWhile(label string) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume WHILE

	condition, err := p.collectUntil("LOOP")
	if err != nil {
		return nil, err
	}
	if !p.isKeyword("LOOP") {
		return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
	}
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLWhile{
		Label:     label,
		Condition: condition,
		Body:      body,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseFor parses all FOR loop variants:
//   - Integer FOR: FOR i IN [REVERSE] lower..upper [BY step] LOOP
//   - Query FOR: FOR rec IN SELECT/WITH ... LOOP
//   - Cursor FOR: FOR rec IN cursor_var [(args)] LOOP
//   - Dynamic FOR: FOR rec IN EXECUTE expr [USING params] LOOP
func (p *Parser) parseFor(label string) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume FOR

	varName, err := p.parseForTargetList()
	if err != nil {
		return nil, err
	}

	// Expect IN
	if err := p.expectKeyword("IN"); err != nil {
		return nil, err
	}

	// Check for REVERSE (integer FOR)
	reverse := false
	if p.isKeyword("REVERSE") {
		reverse = true
		p.advance()
	}

	// Determine which FOR variant:
	// 1. EXECUTE -> dynamic FOR
	// 2. SELECT / WITH -> query FOR
	// 3. Otherwise, scan ahead to check for DOT_DOT (integer range) vs cursor

	if p.isKeyword("EXECUTE") {
		return p.parseForDynamic(label, varName, startPos)
	}

	if p.isKeyword("SELECT") || p.isKeyword("WITH") || p.isParenthesizedQueryStart() {
		return p.parseForQuery(label, varName, startPos)
	}

	// Collect tokens until LOOP, looking for DOT_DOT to disambiguate integer vs cursor
	return p.parseForIntOrCursor(label, varName, reverse, startPos)
}

func (p *Parser) isParenthesizedQueryStart() bool {
	if p.cur.Type != '(' {
		return false
	}
	next := p.peekNext()
	return next.Type == pgparser.SELECT || next.Type == pgparser.WITH ||
		next.Type == pgparser.VALUES || next.Type == pgparser.TABLE
}

// parseForTargetList parses the target before IN in a FOR statement. Query FOR
// accepts a record/row target or a comma-separated scalar-variable list.
func (p *Parser) parseForTargetList() (string, error) {
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return "", p.errorf("syntax error at or near %q, expected loop variable name", p.tokenText(p.cur))
	}
	start := p.pos()
	depth := 0
	for !p.isEOF() {
		if p.cur.Type == '(' || p.cur.Type == '[' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' || p.cur.Type == ']' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}
		if depth == 0 && p.isKeyword("IN") {
			break
		}
		p.advance()
	}
	target := strings.TrimSpace(p.source[start:p.prev.End])
	if target == "" {
		return "", p.errorf("syntax error at or near %q, expected loop variable name", p.tokenText(p.cur))
	}
	return target, nil
}

// parseForDynamic parses: FOR rec IN EXECUTE expr [USING params] LOOP stmts END LOOP [label] ;
func (p *Parser) parseForDynamic(label, varName string, startPos int) (ast.Node, error) {
	p.advance() // consume EXECUTE

	// Collect the dynamic query expression until LOOP or USING
	query, err := p.collectUntil("LOOP", "USING")
	if err != nil {
		return nil, err
	}

	// Parse optional USING clause
	var params []string
	if p.isKeyword("USING") {
		p.advance()
		params, err = p.parseForUsingParams()
		if err != nil {
			return nil, err
		}
	}

	if !p.isKeyword("LOOP") {
		return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
	}
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLForDynS{
		Label:  label,
		Var:    varName,
		Query:  query,
		Params: params,
		Body:   body,
		Loc:    ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseUsingParams parses a comma-separated list of expressions until LOOP keyword.
func (p *Parser) parseForUsingParams() ([]string, error) {
	var params []string
	for {
		expr, err := p.collectUntil("LOOP")
		if err != nil {
			return nil, err
		}
		// The expression might end with a comma - if so, split
		if strings.HasSuffix(expr, ",") {
			expr = strings.TrimSpace(expr[:len(expr)-1])
			params = append(params, expr)
			continue
		}
		params = append(params, expr)
		break
	}
	return params, nil
}

// parseForQuery parses: FOR rec IN SELECT/WITH ... LOOP stmts END LOOP [label] ;
func (p *Parser) parseForQuery(label, varName string, startPos int) (ast.Node, error) {
	// Collect SQL query text until LOOP keyword
	query, err := p.collectUntil("LOOP")
	if err != nil {
		return nil, err
	}

	if !p.isKeyword("LOOP") {
		return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
	}
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLForS{
		Label: label,
		Var:   varName,
		Query: query,
		Body:  body,
		Loc:   ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseForIntOrCursor disambiguates integer range FOR (with ..) from cursor FOR.
// It scans tokens looking for DOT_DOT at paren depth 0.
func (p *Parser) parseForIntOrCursor(label, varName string, reverse bool, startPos int) (ast.Node, error) {
	// Save the start position for collecting text
	exprStart := p.pos()
	depth := 0
	hasDotDot := false

	// Scan to find DOT_DOT or LOOP at depth 0
	// We need to find DOT_DOT to know it's integer, otherwise it's cursor
	// Save token positions as we scan
	type savedToken struct {
		tok pgparser.Token
	}
	var tokens []savedToken

	scanStart := p.pos()
	_ = scanStart

	// We'll manually scan, tracking what we see
	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			tokens = append(tokens, savedToken{p.cur})
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			tokens = append(tokens, savedToken{p.cur})
			p.advance()
			continue
		}
		if depth == 0 {
			if p.cur.Type == pgparser.DOT_DOT {
				hasDotDot = true
				break
			}
			if p.isKeyword("LOOP") {
				break
			}
		}
		tokens = append(tokens, savedToken{p.cur})
		p.advance()
	}

	if hasDotDot {
		if hasTopLevelComma(varName) {
			return nil, p.errorf("syntax error at or near %q, expected single loop variable", varName)
		}

		// Integer FOR: lower..upper [BY step]
		lower := strings.TrimSpace(p.source[exprStart:p.cur.Loc])
		p.advance() // consume DOT_DOT

		// Collect upper bound until BY or LOOP
		upper, err := p.collectUntil("BY", "LOOP")
		if err != nil {
			return nil, err
		}

		// Parse optional BY step
		var step string
		if p.isKeyword("BY") {
			p.advance()
			step, err = p.collectUntil("LOOP")
			if err != nil {
				return nil, err
			}
		}

		if !p.isKeyword("LOOP") {
			return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
		}
		p.advance() // consume LOOP

		body, err := p.parseStmtList()
		if err != nil {
			return nil, err
		}

		if err := p.parseEndLoop(label); err != nil {
			return nil, err
		}

		return &ast.PLForI{
			Label:   label,
			Var:     varName,
			Lower:   lower,
			Upper:   upper,
			Step:    step,
			Reverse: reverse,
			Body:    body,
			Loc:     ast.Loc{Start: startPos, End: p.prev.End},
		}, nil
	}

	// Cursor FOR: the tokens we collected form the cursor variable name,
	// possibly followed by (args)
	if hasTopLevelComma(varName) {
		return nil, p.errorf("syntax error at or near %q, expected single loop variable", varName)
	}
	cursorText := strings.TrimSpace(p.source[exprStart:p.cur.Loc])

	// Parse cursor variable name and optional args from the collected text
	// The cursor var is a single identifier, optionally followed by (args)
	cursorVar := cursorText
	var argQuery string

	// Check if there are parens in cursorText indicating arguments
	parenIdx := strings.Index(cursorText, "(")
	if parenIdx >= 0 {
		cursorVar = strings.TrimSpace(cursorText[:parenIdx])
		// Extract args between outer parens
		lastParen := strings.LastIndex(cursorText, ")")
		if lastParen > parenIdx {
			argQuery = strings.TrimSpace(cursorText[parenIdx+1 : lastParen])
		}
	}

	if !p.isKeyword("LOOP") {
		return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
	}
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLForC{
		Label:     label,
		Var:       varName,
		CursorVar: cursorVar,
		ArgQuery:  argQuery,
		Body:      body,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

func hasTopLevelComma(text string) bool {
	depth := 0
	for _, r := range text {
		switch r {
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

// parseForEach parses: FOREACH var [SLICE n] IN ARRAY expr LOOP stmts END LOOP [label] ;
func (p *Parser) parseForEach(label string) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume FOREACH

	// Parse loop variable name
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return nil, p.errorf("syntax error at or near %q, expected loop variable name", p.tokenText(p.cur))
	}
	varName := p.identText()
	p.advance()

	// Parse optional SLICE n
	sliceDim := 0
	if p.isKeyword("SLICE") {
		p.advance()
		// Expect integer constant
		if p.cur.Type != pgparser.ICONST {
			return nil, p.errorf("syntax error at or near %q, expected integer after SLICE", p.tokenText(p.cur))
		}
		// Parse integer from token text
		sliceText := p.source[p.cur.Loc:p.cur.End]
		for _, ch := range sliceText {
			sliceDim = sliceDim*10 + int(ch-'0')
		}
		p.advance()
	}

	// Expect IN
	if err := p.expectKeyword("IN"); err != nil {
		return nil, err
	}

	// Expect ARRAY
	if err := p.expectKeyword("ARRAY"); err != nil {
		return nil, err
	}

	// Collect array expression until LOOP
	arrayExpr, err := p.collectUntil("LOOP")
	if err != nil {
		return nil, err
	}

	if !p.isKeyword("LOOP") {
		return nil, p.errorf("syntax error at or near %q, expected LOOP", p.tokenText(p.cur))
	}
	p.advance() // consume LOOP

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	if err := p.parseEndLoop(label); err != nil {
		return nil, err
	}

	return &ast.PLForEachA{
		Label:     label,
		Var:       varName,
		SliceDim:  sliceDim,
		ArrayExpr: arrayExpr,
		Body:      body,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// Ensure imports are used.
var _ = pgparser.DOT_DOT
