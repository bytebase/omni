package parser

// tokenStreamState captures every parser+lexer field that mutates while
// consuming tokens, so a speculative scan can be fully rewound.
//
// `execStack` and `execMarkers` are captured because the lexer pushes and pops
// executable comments (/*! ... */, /*T! ... */) as it steps through them. A
// push never writes into the saved stack's backing array, and markers are only
// appended, so keeping the slice headers is enough.
//
// Completion state (candidates/collecting) is intentionally NOT captured: the
// classifier scan is a pure advance() loop, and advance()/peekNext() never
// collect candidates (collection lives only in grammar-rule functions, which
// the scan does not call).
type tokenStreamState struct {
	cur, prev, nextBuf Token
	hasNext            bool
	pos, start         int
	prevToken          int
	prevTokenEnd       int
	execStack          []execComment
	execMarkers        []execMarker
}

// snapshotTokenStream captures the current token-stream position for later
// restoration via restoreTokenStream.
func (p *Parser) snapshotTokenStream() tokenStreamState {
	return tokenStreamState{
		cur: p.cur, prev: p.prev, nextBuf: p.nextBuf, hasNext: p.hasNext,
		pos: p.lexer.pos, start: p.lexer.start,
		prevToken: p.lexer.prevToken, prevTokenEnd: p.lexer.prevTokenEnd,
		execStack: p.lexer.execStack, execMarkers: p.lexer.execMarkers,
	}
}

// restoreTokenStream rewinds parser+lexer state to a previously captured
// snapshot. After restore, the next advance() emits the same token it would
// have at the moment snapshotTokenStream() was called.
func (p *Parser) restoreTokenStream(s tokenStreamState) {
	p.cur, p.prev, p.nextBuf, p.hasNext = s.cur, s.prev, s.nextBuf, s.hasNext
	p.lexer.pos, p.lexer.start = s.pos, s.start
	p.lexer.prevToken, p.lexer.prevTokenEnd = s.prevToken, s.prevTokenEnd
	p.lexer.execStack, p.lexer.execMarkers = s.execStack, s.execMarkers
}
