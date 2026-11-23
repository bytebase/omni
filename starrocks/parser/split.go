package parser

// Segment represents one top-level SQL statement extracted from a source string.
//
// Text is the raw substring of the input from ByteStart (inclusive) to
// ByteEnd (exclusive). It is NOT trimmed — leading/trailing whitespace and
// comments are preserved verbatim. The trailing `;` delimiter (if present)
// is NOT part of Text or ByteEnd; it lives between this segment's ByteEnd
// and the next segment's ByteStart.
type Segment struct {
	Text      string // the raw text of the statement (no trailing semicolon)
	ByteStart int    // inclusive start byte offset in the original source
	ByteEnd   int    // exclusive end byte offset; points AT the trailing ; if present, otherwise at len(input)
}

// Empty reports whether the segment contains no meaningful SQL content —
// that is, only whitespace and comments. It works by re-lexing Text and
// checking whether the first token is tokEOF.
//
// Split returns empty segments; a caller that parses or executes statements
// skips them.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Kind == tokEOF
}

// splitState is the state machine state for Split.
type splitState int

const (
	stateTop     splitState = iota
	stateInBlock            // inside a BEGIN..END compound block
)

// Split extracts top-level SQL statements from input.
//
// The returned slice contains one Segment per delimiter-separated chunk,
// including the chunks that hold no statement — a lone semicolon, whitespace,
// a comment, or comment-only text after the last `;` — which Segment.Empty
// reports. A caller that accounts for every delimiter
// (analysis.SplitStatements) needs them; a caller that parses or executes
// statements skips them.
//
// Split is infallible — it always returns a result, even for malformed
// input. Lexing errors are suppressed internally; callers who need them
// should use NewLexer(input).Errors() directly.
//
// Split correctly handles:
//   - Single-quoted strings with doubled-quote and \ escapes
//   - Double-quoted strings with "" escapes
//   - Backtick-quoted identifiers
//   - X'...' hex literals and B'...' bit literals
//   - Line comments (--) and block comments (/* */)
//   - Doris BEGIN..END compound blocks (for stored procedures)
//
// Split does NOT handle:
//   - DELIMITER directive (MySQL-specific)
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0
	// prefix holds the first tokens of the statement being split, up to the
	// ones beginOpensBlock reads, so a BEGIN inside an expression is told
	// apart from one that opens a statement's compound block.
	var prefix []TokenKind

	// We need lookahead after kwBEGIN to disambiguate TCL from a compound
	// block: one token for most followers, two for BEGIN WITH LABEL (WITH
	// alone could open a block whose first statement is a CTE). Use a small
	// buffered queue.
	var pending []Token
	nextToken := func() Token {
		if len(pending) > 0 {
			t := pending[0]
			pending = pending[1:]
			return t
		}
		return l.NextToken()
	}
	peekToken := func() Token {
		if len(pending) == 0 {
			pending = append(pending, l.NextToken())
		}
		return pending[0]
	}
	peekToken2 := func() Token {
		for len(pending) < 2 {
			pending = append(pending, l.NextToken())
		}
		return pending[1]
	}

	// emit creates a Segment covering input[stmtStart:end] — where end is
	// the byte offset BEFORE any trailing delimiter. The caller is
	// responsible for advancing stmtStart past the delimiter.
	emit := func(end int) {
		segments = append(segments, Segment{
			Text:      input[stmtStart:end],
			ByteStart: stmtStart,
			ByteEnd:   end,
		})
		prefix = prefix[:0]
	}

	for {
		tok := nextToken()
		if tok.Kind == tokEOF {
			break
		}

		switch state {
		case stateTop:
			switch tok.Kind {
			case kwBEGIN:
				if !beginOpensBlock(prefix) {
					// BEGIN is not the statement's first word and not a
					// procedure's body opener: it is a name, as BEGIN is
					// nonReserved. `SELECT BEGIN zzz;` stays one statement.
					break
				}
				next := peekToken()
				// BEGIN WITH LABEL x is the transaction form. WITH alone does
				// not decide — a compound block's first statement can be a
				// CTE — so the LABEL after it is what disambiguates.
				isTCL := isDorisTCLBeginFollower(next) ||
					(next.Kind == kwWITH && peekToken2().Kind == kwLABEL)
				if isTCL {
					// TCL: BEGIN;, BEGIN TRANSACTION, BEGIN WORK, BEGIN EOF.
					// Stay in stateTop; the buffered peek will be returned
					// by the next nextToken() call and handled normally.
				} else {
					// Compound block BEGIN..END.
					state = stateInBlock
					depth = 1
				}
			case int(';'):
				// tok.Loc.Start is the position of the `;`, tok.Loc.End is
				// one past the `;`. Segment text stops BEFORE the `;`; the
				// next segment starts AFTER the `;`.
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
				continue
			}
			if len(prefix) < blockBeginPrefix {
				prefix = append(prefix, tok.Kind)
			}

		case stateInBlock:
			switch tok.Kind {
			case kwBEGIN:
				depth++
			case kwEND:
				if depth > 0 {
					depth--
					if depth == 0 {
						state = stateTop
					}
				}
			}
			// Semicolons and all other tokens are absorbed in stateInBlock.
		}
	}

	// Trailing segment — whatever remains after the last `;` (or the whole
	// input if there was no `;`). ByteEnd is len(input) in this case.
	if stmtStart < len(input) {
		emit(len(input))
	}

	return segments
}

// blockBeginPrefix is how many tokens of a statement blockBeginPrefix reads.
const blockBeginPrefix = 5

// beginOpensBlock reports whether a BEGIN that appears at the top level of a
// statement opens a compound block rather than naming a column. BEGIN is
// nonReserved, so `SELECT BEGIN zzz` names one; only a statement's first
// token, or the body of a CREATE [OR REPLACE] PROCEDURE, opens a block.
func beginOpensBlock(prefix []TokenKind) bool {
	if len(prefix) == 0 {
		return true
	}
	if prefix[0] != kwCREATE || len(prefix) < 2 {
		return false
	}
	if prefix[1] == kwPROCEDURE {
		return true
	}
	return len(prefix) >= 4 && prefix[1] == kwOR &&
		prefix[2] == kwREPLACE && prefix[3] == kwPROCEDURE
}

// isDorisTCLBeginFollower reports whether tok, appearing immediately after a
// top-level BEGIN keyword, indicates a transaction-start (TCL) rather than
// a compound block opener.
//
// TCL forms: BEGIN;, BEGIN TRANSACTION, BEGIN WORK, BEGIN EOF
func isDorisTCLBeginFollower(tok Token) bool {
	switch tok.Kind {
	case int(';'), tokEOF, kwTRANSACTION, kwWORK:
		return true
	}
	return false
}
