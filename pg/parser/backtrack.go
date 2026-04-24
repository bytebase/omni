package parser

// tokenStreamState captures the parser + lexer state needed to rewind the
// token stream to a previous token boundary. This is sufficient for the
// "speculative parse, then rollback if it doesn't match" pattern at every
// site in pg/parser as of this writing.
//
// SCOPE: this is a TOKEN-STREAM snapshot, not a complete parser/lexer
// snapshot. It does NOT cover mid-token-content lexer state (literalbuf,
// dolqstart, utf16FirstPart, xcdepth, stateBeforeStrStop, warning flags)
// or completion-mode state (candidates, collecting). Lexer internals are
// reset at token boundaries. Completion-mode speculative callers that can
// scan past the cursor should use snapshotTokenStreamAndCompletion.
//
// If a future caller needs to roll back from INSIDE a token (e.g., from
// inside a string literal or dollar-quoted block), this struct is
// insufficient — extend it carefully and update the doc comment.
//
// Why this exists: see commit history for two prior incidents where
// hand-rolled rollback machinery captured only a subset of the necessary
// state and corrupted the parser:
//
//   - create_function.go pushBack: rewrote token type to IDENT
//     unconditionally on rollback, breaking CREATE FUNCTION arguments
//     ending in 'double precision' (DOUBLE_P is the only UnreservedKeyword
//     in the type lead set, so only that one keyword surfaced the bug).
//
//   - type.go parseFuncType speculative branch: saved cur/prev/nextBuf/
//     hasNext/lexer.Err but NOT lexer.pos/start/state, so any qualified
//     type at a parseFuncType call site (function param, function return,
//     RETURNS TABLE column, CREATE OPERATOR LEFTARG/RIGHTARG) corrupted
//     the lexer position on rollback.
//
// Both bugs would have been caught by using this helper from day one.
type tokenStreamState struct {
	cur, prev, nextBuf Token
	hasNext            bool
	lexerErr           error
	lexerPos           int
	lexerStart         int
	lexerState         LexerState
}

type tokenStreamAndCompletionState struct {
	tokenStream  tokenStreamState
	collecting   bool
	collectDepth int
	candidates   *CandidateSet
}

// snapshotTokenStream captures the current token-stream position for
// later restoration via restoreTokenStream. See tokenStreamState for
// scope and limitations.
func (p *Parser) snapshotTokenStream() tokenStreamState {
	return tokenStreamState{
		cur:        p.cur,
		prev:       p.prev,
		nextBuf:    p.nextBuf,
		hasNext:    p.hasNext,
		lexerErr:   p.lexer.Err,
		lexerPos:   p.lexer.pos,
		lexerStart: p.lexer.start,
		lexerState: p.lexer.state,
	}
}

// snapshotTokenStreamAndCompletion captures token-stream state plus the
// completion state that advance() can mutate when a speculative walk crosses
// the cursor. Use this for completion-mode lookahead that can scan arbitrary
// user input before rolling back.
func (p *Parser) snapshotTokenStreamAndCompletion() tokenStreamAndCompletionState {
	s := tokenStreamAndCompletionState{
		tokenStream:  p.snapshotTokenStream(),
		collecting:   p.collecting,
		collectDepth: p.collectDepth,
	}
	if p.candidates != nil {
		s.candidates = p.candidates.snapshot()
	}
	return s
}

// restoreTokenStream rewinds parser + lexer state to a previously
// captured snapshot. After restore, the next advance() will emit the
// same token as it would have at the moment snapshotTokenStream() was
// called.
//
// Caller responsibility: do not interleave restore with completion-mode
// queries or with any operation that mutates lexer state outside the token
// stream. Use restoreTokenStreamAndCompletion for lookahead that may cross
// the completion cursor.
func (p *Parser) restoreTokenStream(s tokenStreamState) {
	p.cur = s.cur
	p.prev = s.prev
	p.nextBuf = s.nextBuf
	p.hasNext = s.hasNext
	p.lexer.Err = s.lexerErr
	p.lexer.pos = s.lexerPos
	p.lexer.start = s.lexerStart
	p.lexer.state = s.lexerState
}

// restoreTokenStreamAndCompletion rewinds token-stream and completion state
// captured by snapshotTokenStreamAndCompletion.
func (p *Parser) restoreTokenStreamAndCompletion(s tokenStreamAndCompletionState) {
	p.restoreTokenStream(s.tokenStream)
	p.collecting = s.collecting
	p.collectDepth = s.collectDepth
	if p.candidates != nil && s.candidates != nil {
		p.candidates.restore(s.candidates)
	}
}
