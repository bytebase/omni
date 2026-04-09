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
// Empty segments are filtered out by Split.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Type == tokEOF
}

// splitState is the state machine state for Split.
type splitState int

const (
	stateTop splitState = iota
	stateInDeclare
	stateInBlock
)

// Split extracts top-level SQL statements from input.
//
// The returned slice contains one Segment per non-empty statement. Empty
// statements (lone semicolons, comment-only chunks) are filtered out.
//
// Split is infallible — it always returns a result, even for malformed
// input. Lexing errors are suppressed internally; callers who need them
// should use NewLexer(input).Errors() directly.
//
// Split correctly handles:
//   - Single-quoted strings with ” and \ escapes
//   - Double-quoted identifiers with "" escape
//   - $$...$$ dollar strings
//   - X'...' hex literals
//   - Line comments (-- and //) and block comments (/* */ including nested)
//   - Snowflake Scripting BEGIN..END blocks (including nested and DECLARE..BEGIN..END)
//   - Inline procedure bodies (CREATE TASK/PROCEDURE/FUNCTION ... AS BEGIN ... END;)
//
// Split does NOT handle:
//   - IF/FOR/WHILE/REPEAT/CASE at top level (these are only valid inside a BEGIN..END body)
//   - DECLARE CURSOR at top level without a matching BEGIN (unusual; best-effort)
//   - DELIMITER directive (MySQL-specific, not used in Snowflake)
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0

	// We need one-token lookahead after kwBEGIN to disambiguate TCL from
	// scripting. Use a one-slot buffered lookahead.
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

	// emit creates a Segment covering input[stmtStart:end] — where end is
	// the byte offset BEFORE any trailing delimiter. The caller is
	// responsible for advancing stmtStart past the delimiter.
	emit := func(end int) {
		seg := Segment{
			Text:      input[stmtStart:end],
			ByteStart: stmtStart,
			ByteEnd:   end,
		}
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
			case kwDECLARE:
				state = stateInDeclare
			case kwBEGIN:
				next := peekToken()
				if isTCLBeginFollower(next) {
					// Stay in stateTop; the buffered peek will be returned
					// by the next nextToken() call and handled normally
					// (including the `;` emission path if applicable).
				} else {
					state = stateInBlock
					depth = 1
				}
			case int(';'):
				// tok.Loc.Start is the position of the `;`, tok.Loc.End is
				// one past the `;`. Segment text stops BEFORE the `;`; the
				// next segment starts AFTER the `;`.
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
			}

		case stateInDeclare:
			switch tok.Type {
			case kwBEGIN:
				state = stateInBlock
				depth = 1
			}
			// Semicolons and all other tokens are absorbed in stateInDeclare.

		case stateInBlock:
			switch tok.Type {
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

// isTCLBeginFollower reports whether tok, appearing immediately after a
// top-level BEGIN keyword, indicates a transaction-start (TCL) rather than
// a Snowflake Scripting block opener.
//
// TCL forms: BEGIN;, BEGIN TRANSACTION, BEGIN WORK, BEGIN NAME <id>, BEGIN EOF
func isTCLBeginFollower(tok Token) bool {
	switch tok.Type {
	case int(';'), tokEOF, kwTRANSACTION, kwWORK, kwNAME:
		return true
	}
	return false
}
