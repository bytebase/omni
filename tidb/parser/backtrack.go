package parser

// tokenStreamState captures every parser+lexer field that mutates while
// consuming tokens, so a speculative scan can be fully rewound.
//
// `input` is captured because the lexer rewrites it in place when expanding
// MySQL conditional comments (lexer.go: `l.input = l.input[:l.pos] + ...`).
// The classifier always snapshots before any splice it triggers, so position
// alone would suffice (the splice preserves the prefix `l.input[:l.pos]`); we
// capture `input` anyway — O(1), Go strings are immutable — as defensive
// insurance so the primitive is correct for any caller, not just a forward
// scan.
//
// Completion state (candidates/collecting) is intentionally NOT captured: the
// classifier scan is a pure advance() loop, and advance()/peekNext() never
// collect candidates (collection lives only in grammar-rule functions, which
// the scan does not call).
type tokenStreamState struct {
	cur, prev, nextBuf Token
	hasNext            bool
	input              string
	pos, start         int
	prevToken          int
	prevTokenEnd       int
}

// snapshotTokenStream captures the current token-stream position for later
// restoration via restoreTokenStream.
func (p *Parser) snapshotTokenStream() tokenStreamState {
	return tokenStreamState{
		cur: p.cur, prev: p.prev, nextBuf: p.nextBuf, hasNext: p.hasNext,
		input: p.lexer.input, pos: p.lexer.pos, start: p.lexer.start,
		prevToken: p.lexer.prevToken, prevTokenEnd: p.lexer.prevTokenEnd,
	}
}

// restoreTokenStream rewinds parser+lexer state to a previously captured
// snapshot. After restore, the next advance() emits the same token it would
// have at the moment snapshotTokenStream() was called.
func (p *Parser) restoreTokenStream(s tokenStreamState) {
	p.cur, p.prev, p.nextBuf, p.hasNext = s.cur, s.prev, s.nextBuf, s.hasNext
	p.lexer.input = s.input
	p.lexer.pos, p.lexer.start = s.pos, s.start
	p.lexer.prevToken, p.lexer.prevTokenEnd = s.prevToken, s.prevTokenEnd
}
