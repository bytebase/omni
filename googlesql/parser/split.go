package parser

// This file ships the GoogleSQL statement splitters. GoogleSQL is consumed by
// two bytebase engines whose legacy splitters DIVERGE (see the migration
// contract, §4 "Statement split"):
//
//   - BigQuery uses a pure lexer/token split — every top-level ';' ends a
//     statement (base.SplitSQLByLexer over the GoogleSQL lexer). BigQuery's
//     editor surface is single statements / scripts without nested control
//     flow at this layer.
//   - Spanner uses a parse-tree split that keeps procedural BEGIN/END (and
//     IF/CASE/LOOP/WHILE/REPEAT/FOR) blocks whole so a ';' separating control
//     statements inside a stored-procedure body does not split the CREATE
//     PROCEDURE into pieces (see the legacy split_begin_end_test.go).
//
// The foundation ships BOTH as pure functions over Segment so the
// bytebase-switch node can wire each engine to the matching variant:
//
//   - Split      — block-aware (Spanner semantics; also correct for BigQuery
//                  scripting, since a flat split is the degenerate case when
//                  no block keywords appear).
//   - SplitFlat  — pure lexer split on every top-level ';' (BigQuery
//                  semantics).
//
// Both reuse the hand-written Lexer, which already hides string / bytes /
// backtick-identifier / comment content, so a ';' inside a literal, a quoted
// identifier, or a comment is never mistaken for a delimiter.

// Segment represents one top-level SQL statement extracted from a source
// string.
//
// Text is the raw substring of the input from ByteStart (inclusive) to ByteEnd
// (exclusive). It is NOT trimmed — leading/trailing whitespace and comments are
// preserved verbatim. The trailing ';' delimiter (if present) is NOT part of
// Text or ByteEnd; it lives between this segment's ByteEnd and the next
// segment's ByteStart.
type Segment struct {
	Text      string // raw text of the statement (no trailing semicolon)
	ByteStart int    // inclusive start byte offset in the original source
	ByteEnd   int    // exclusive end byte offset; points AT the trailing ';' if present, else at len(input)
}

// Empty reports whether the segment contains no meaningful SQL content — only
// whitespace and comments. It re-lexes Text and checks whether the first token
// is tokEOF. Empty segments are filtered out by both splitters.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Type == tokEOF
}

// splitState is the state-machine state for the block-aware Split.
type splitState int

const (
	stateTop     splitState = iota // between statements / inside a non-block statement
	stateInBlock                   // inside a procedural BEGIN/IF/CASE/LOOP/WHILE/REPEAT/FOR block body
)

// SplitFlat splits input on every top-level ';' — the BigQuery lexer-split
// semantics. It mirrors the legacy bigquery.SplitSQL (base.SplitSQLByLexer over
// the GoogleSQL lexer): control-flow keywords are NOT treated specially, so a
// ';' inside a procedural body WOULD split (BigQuery's wired surface does not
// rely on block-aware splitting at this layer).
//
// The returned slice contains one Segment per non-empty statement; empty
// statements (lone semicolons, comment-only chunks) are filtered out. SplitFlat
// is infallible — lexing errors are suppressed here; callers who need them use
// NewLexer(input).Errors() or Parse, which promotes lex errors to diagnostics.
func SplitFlat(input string) []Segment {
	if len(input) == 0 {
		return nil
	}
	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0

	emit := func(end int) {
		seg := Segment{Text: input[stmtStart:end], ByteStart: stmtStart, ByteEnd: end}
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}

	for {
		tok := l.NextToken()
		if tok.Type == tokEOF {
			break
		}
		if tok.Type == int(';') {
			emit(tok.Loc.Start)
			stmtStart = tok.Loc.End
		}
	}
	if stmtStart < len(input) {
		emit(len(input))
	}
	return segments
}

// Split extracts top-level SQL statements from input, keeping procedural
// BEGIN/END (and IF/CASE/LOOP/WHILE/REPEAT/FOR) blocks whole — the Spanner
// parse-tree-split semantics. A ';' that separates control statements inside a
// stored-procedure body does not split the enclosing statement.
//
// The returned slice contains one Segment per non-empty statement; empty
// statements (lone semicolons, comment-only chunks) are filtered out. Split is
// infallible (lexing errors are suppressed here, as in SplitFlat).
//
// Block model (GoogleSQLParser.g4 procedural statements):
//   - Openers: BEGIN, IF, CASE, LOOP, WHILE, REPEAT, FOR.
//   - Closer: END, optionally followed by the block-kind keyword
//     (END IF / END CASE / END LOOP / END WHILE / END REPEAT / END FOR).
//   - Blocks nest; an inner block's END only closes that inner block.
//
// A top-level BEGIN is ambiguous: BEGIN [TRANSACTION|WORK] is a transaction
// start (TCL), while BEGIN ... END is a procedural block. One-token lookahead
// after BEGIN disambiguates: a TCL follower (';', EOF, TRANSACTION) keeps the
// statement flat; anything else opens a block. (CASE is doubly ambiguous — it
// is also an expression; but a CASE expression's END never carries a trailing
// ';' issue at top level because the surrounding query splits on its own ';'.
// The block-aware closer counting is harmless for the expression form: a CASE
// expression with no internal ';' splits identically either way, which the
// legacy split_begin_end_test.go relies on.)
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0 // block-nesting depth while in stateInBlock

	// prevWasEnd is set immediately after an END token so a block keyword that
	// completes a typed closer — END IF / END CASE / END LOOP / END WHILE /
	// END REPEAT / END FOR — is recognized as the closer's suffix rather than
	// miscounted as a NEW block opener (which would desync depth and swallow
	// the statement-terminating ';' after the block).
	prevWasEnd := false

	// We need one-token lookahead after a top-level BEGIN to tell a
	// transaction-start (TCL) from a procedural block opener. A one-slot
	// buffered lookahead does this without disturbing the main scan.
	var pending *Token
	nextToken := func() Token {
		if pending != nil {
			t := *pending
			pending = nil
			return t
		}
		return l.NextToken()
	}
	peekToken := func() Token {
		if pending == nil {
			t := l.NextToken()
			pending = &t
		}
		return *pending
	}

	emit := func(end int) {
		seg := Segment{Text: input[stmtStart:end], ByteStart: stmtStart, ByteEnd: end}
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}

	for {
		tok := nextToken()
		if tok.Type == tokEOF {
			break
		}

		switch state {
		case stateTop:
			switch tok.Type {
			case kwBEGIN, kwIF, kwCASE, kwLOOP, kwWHILE, kwREPEAT, kwFOR:
				switch {
				case prevWasEnd:
					// This block keyword immediately follows the END that closed
					// the OUTERMOST block (depth 1→0), so it is the suffix of a
					// typed closer (END IF / END CASE / END LOOP / …), not a new
					// opener. Consume it without reopening. Without this, a typed
					// closer at the top level would re-enter a block and swallow
					// the statement-terminating ';' that follows.
					prevWasEnd = false
				case tok.Type == kwBEGIN && isTCLBeginFollower(peekToken()):
					// BEGIN [TRANSACTION] is a transaction start — stay flat. The
					// buffered follower is returned by the next nextToken() and
					// handled normally (including the ';' emission path).
				default:
					// A top-level control-flow opener begins a procedural block
					// body whose internal ';' must be kept whole.
					state = stateInBlock
					depth = 1
				}
			case int(';'):
				prevWasEnd = false
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
			default:
				prevWasEnd = false
			}

		case stateInBlock:
			switch tok.Type {
			case kwBEGIN, kwIF, kwCASE, kwLOOP, kwWHILE, kwREPEAT, kwFOR:
				// A block keyword right after END is the suffix of a typed
				// closer (END IF / END CASE / …), not a new opener.
				if prevWasEnd {
					prevWasEnd = false
				} else {
					depth++
				}
			case kwEND:
				prevWasEnd = true
				if depth > 0 {
					depth--
					if depth == 0 {
						// Return to stateTop but KEEP prevWasEnd set: a typed
						// closer's suffix keyword (END IF / END WHILE / …) may
						// follow and must be recognized as the suffix by the
						// stateTop arm above, not mistaken for a new opener.
						state = stateTop
					}
				}
			default:
				prevWasEnd = false
			}
			// ';' and every other token are absorbed inside the block body.
		}
	}

	// Trailing segment — whatever remains after the last ';' (or the whole
	// input when there was none). ByteEnd is len(input) here.
	if stmtStart < len(input) {
		emit(len(input))
	}

	return segments
}

// isTCLBeginFollower reports whether tok, appearing immediately after a
// top-level BEGIN keyword, indicates a transaction start (TCL) rather than a
// procedural block opener.
//
// GoogleSQL begin_statement: begin_transaction_keywords transaction_mode_list?
// where begin_transaction_keywords is `BEGIN [TRANSACTION]` and a
// transaction_mode begins with READ (READ ONLY | READ WRITE) or ISOLATION
// (ISOLATION LEVEL <id>). So a BEGIN is a transaction start when it is
// immediately followed by:
//   - ';' or EOF                 — a bare `BEGIN;` / `BEGIN`
//   - TRANSACTION                — `BEGIN TRANSACTION ...`
//   - READ / ISOLATION           — `BEGIN READ ONLY|WRITE`, `BEGIN ISOLATION LEVEL ...`
//     (the optional TRANSACTION keyword may be omitted, so a mode list can
//     follow BEGIN directly — oracle-confirmed against the Spanner emulator).
//
// Anything else after BEGIN (a DECLARE, a SELECT, an IF, ...) opens a
// procedural begin_end_block, whose internal ';' must be kept whole.
func isTCLBeginFollower(tok Token) bool {
	switch tok.Type {
	case int(';'), tokEOF, kwTRANSACTION, kwREAD, kwISOLATION:
		return true
	}
	return false
}
