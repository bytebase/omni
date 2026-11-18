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
	//
	// Hints (/*+ ... */) are read past, as the parser reads past them
	// (parser.nextToken): they hold no syntax, so they decide no boundary.
	nextRaw := func() Token {
		tok := l.NextToken()
		for tok.Kind == tokHintStart {
			for tok.Kind != tokHintEnd && tok.Kind != tokEOF {
				tok = l.NextToken()
			}
			if tok.Kind == tokEOF {
				return tok
			}
			tok = l.NextToken()
		}
		return tok
	}
	var pending []Token
	nextToken := func() Token {
		if len(pending) > 0 {
			t := pending[0]
			pending = pending[1:]
			return t
		}
		return nextRaw()
	}
	peekToken := func() Token {
		if len(pending) == 0 {
			pending = append(pending, nextRaw())
		}
		return pending[0]
	}
	peekToken2 := func() Token {
		for len(pending) < 2 {
			pending = append(pending, nextRaw())
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
				if proc.inParens {
					// A BEGIN inside a procedure's parentheses is text the
					// rule's `.*?` takes, not a compound block.
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
					proc.reset()
				}
			case int(';'):
				if proc.inParens {
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
// Split keeps it that way. The `.*?` closes at the first ')', so nothing inside
// the parentheses nests, and a BEGIN there is text, not a compound block.
type procedureHead struct {
	state    procedureHeadState
	inParens bool // the procedure's '(' has been read and not closed yet
}

type procedureHeadState int

const (
	procHeadStart           procedureHeadState = iota // the statement's first token
	procHeadCreate                                    // CREATE
	procHeadCreateOr                                  // CREATE OR
	procHeadCreateOrReplace                           // CREATE OR REPLACE
	procHeadAlter                                     // ALTER
	procHeadReplace                                   // REPLACE
	procHeadVerb                                      // PROCEDURE or PROC
	procHeadNamePart                                  // a part of the name
	procHeadNameDone                                  // a part read: '.' or '(' comes next
	procHeadNot                                       // not a procedure statement
)

// step advances the reading with one top-level token.
func (h *procedureHead) step(tok Token) {
	if h.inParens {
		if tok.Kind == int(')') {
			h.inParens = false
			h.state = procHeadNot
		}
		return
	}
	switch h.state {
	case procHeadStart:
		switch tok.Kind {
		case kwCREATE:
			h.state = procHeadCreate
		case kwALTER:
			h.state = procHeadAlter
		case kwREPLACE:
			h.state = procHeadReplace
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
	case procHeadCreateOrReplace, procHeadAlter, procHeadReplace:
		if tok.Kind == kwPROCEDURE || tok.Kind == kwPROC {
			h.state = procHeadVerb
		} else {
			h.state = procHeadNot
		}
	case procHeadVerb:
		// The rule names the procedure with multipartIdentifier, whose parts
		// are identifiers and identifiers written as keywords.
		if tok.Kind == tokIdent || tok.Kind == tokQuotedIdent || tok.Kind >= 700 {
			h.state = procHeadNameDone
		} else {
			h.state = procHeadNot
		}
	case procHeadNamePart:
		if tok.Kind == tokIdent || tok.Kind == tokQuotedIdent || tok.Kind >= 700 {
			h.state = procHeadNameDone
		} else {
			h.state = procHeadNot
		}
	case procHeadNameDone:
		switch tok.Kind {
		case int('('):
			h.inParens = true
		case int('.'):
			h.state = procHeadNamePart
		default:
			h.state = procHeadNot
		}
	}
}

func (h *procedureHead) reset() { *h = procedureHead{} }
