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
// Empty segments are filtered out by Split and kept by SplitAll.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Kind == tokEOF
}

// splitState is the state machine state for Split and SplitAll.
type splitState int

const (
	stateTop     splitState = iota
	stateInBlock            // inside a BEGIN..END compound block
)

// Split extracts top-level SQL statements from input.
//
// The returned slice contains one Segment per non-empty statement. Empty
// statements (lone semicolons, comment-only chunks) are filtered out;
// SplitAll returns them marked Empty.
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
	segments := SplitAll(input)
	kept := segments[:0]
	for _, segment := range segments {
		if !segment.Empty() {
			kept = append(kept, segment)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	return kept
}

// SplitAll is Split without the filtering: it returns one Segment per
// delimiter-separated chunk, including the chunks that hold no statement — a
// lone semicolon, whitespace, a comment — which Segment.Empty reports. A
// caller that must account for every delimiter (analysis.SplitStatements)
// needs them; a caller that executes statements wants Split.
func SplitAll(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0
	var proc procedureHead

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
	}

	for {
		tok := nextToken()
		if tok.Kind == tokEOF {
			break
		}

		switch state {
		case stateTop:
			proc.step(tok)
			switch tok.Kind {
			case kwBEGIN:
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
					proc.reset()
				}
			case int(';'):
				if proc.parens > 0 {
					// A ';' inside a procedure's parentheses separates
					// nothing: it is part of the statement
					// (DorisParser.g4:43 LEFT_PAREN .*? RIGHT_PAREN).
					continue
				}
				// tok.Loc.Start is the position of the `;`, tok.Loc.End is
				// one past the `;`. Segment text stops BEFORE the `;`; the
				// next segment starts AFTER the `;`.
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
				proc.reset()
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

// procedureHead follows a statement's tokens to a procedure's parentheses
// (DorisParser.g4:43 (PROCEDURE | PROC) name LEFT_PAREN .*? RIGHT_PAREN).
// DorisLexer.g4's SEMICOLON is otherwise only a statement delimiter
// (DorisParser.g4:29-33), so the rule's `.*?` is the one place a ';' stands
// inside a statement: CREATE PROCEDURE a(x;y) is one statement for Doris, and
// Split keeps it that way.
type procedureHead struct {
	state  procedureHeadState
	parens int // > 0 inside the procedure's parentheses
}

type procedureHeadState int

const (
	procHeadStart           procedureHeadState = iota // the statement's first token
	procHeadCreate                                    // CREATE
	procHeadCreateOr                                  // CREATE OR
	procHeadCreateOrReplace                           // CREATE OR REPLACE
	procHeadVerb                                      // PROCEDURE or PROC
	procHeadName                                      // in the name
	procHeadNot                                       // not a procedure statement
)

// step advances the reading with one top-level token.
func (h *procedureHead) step(tok Token) {
	if h.parens > 0 {
		switch tok.Kind {
		case int('('):
			h.parens++
		case int(')'):
			h.parens--
		}
		return
	}
	switch h.state {
	case procHeadStart:
		switch tok.Kind {
		case kwCREATE:
			h.state = procHeadCreate
		case kwALTER, kwREPLACE:
			h.state = procHeadVerb
		default:
			h.state = procHeadNot
		}
	case procHeadCreate:
		switch tok.Kind {
		case kwOR:
			h.state = procHeadCreateOr
		case kwPROCEDURE, kwPROC:
			h.state = procHeadVerb
		default:
			h.state = procHeadNot
		}
	case procHeadCreateOr:
		if tok.Kind == kwREPLACE {
			h.state = procHeadCreateOrReplace
		} else {
			h.state = procHeadNot
		}
	case procHeadCreateOrReplace:
		if tok.Kind == kwPROCEDURE || tok.Kind == kwPROC {
			h.state = procHeadVerb
		} else {
			h.state = procHeadNot
		}
	case procHeadVerb:
		// The rule names the procedure with multipartIdentifier, whose parts
		// are identifiers and identifiers written as keywords.
		if tok.Kind == tokIdent || tok.Kind == tokQuotedIdent || tok.Kind >= 700 {
			h.state = procHeadName
		} else {
			h.state = procHeadNot
		}
	case procHeadName:
		switch {
		case tok.Kind == int('('):
			h.parens = 1
		case tok.Kind == int('.'), tok.Kind == tokIdent, tok.Kind == tokQuotedIdent, tok.Kind >= 700:
			// Another part of the name.
		default:
			h.state = procHeadNot
		}
	}
}

func (h *procedureHead) reset() { *h = procedureHead{} }
