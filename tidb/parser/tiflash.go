package parser

// parseLocationLabelList parses the optional `LOCATION LABELS 'l1', 'l2',
// ...` clause shared by `SET TIFLASH REPLICA` on tables (parser.y:2193),
// databases (parser.y:4482), and the HYPO variant (2204, unsupported).
//
// Returns an empty slice when the clause is absent — upstream's
// LocationLabelList rule (parser.y:2176-2183) has an epsilon arm, so
// omitting the clause is valid.
//
// When present, the label list is non-empty (StringList at parser.y:13570
// is one-or-more, with mandatory comma separators between items, no
// parentheses, no trailing comma).
func (p *Parser) parseLocationLabelList() ([]string, error) {
	// Completion: after the REPLICA count, offer LOCATION as an
	// optional continuation. The labels themselves are user strings,
	// not candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwLOCATION)
		return nil, &ParseError{Message: "collecting"}
	}
	if p.cur.Type != kwLOCATION {
		return nil, nil
	}
	p.advance() // consume LOCATION
	if _, err := p.expect(kwLABELS); err != nil {
		return nil, err
	}
	// Non-empty string list: at least one stringLit, then zero-or-more
	// `, stringLit` pairs.
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	labels := []string{p.cur.Str}
	p.advance()
	for p.cur.Type == ',' {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		labels = append(labels, p.cur.Str)
		p.advance()
	}
	return labels, nil
}
