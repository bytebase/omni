package parser

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
// is tokEOF. Empty segments are filtered out by Split.
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Kind == tokEOF
}

// splitState is the state-machine state for Split.
type splitState int

const (
	stateTop    splitState = iota // between statements / inside a non-routine statement
	stateInBody                   // inside a CREATE/WITH FUNCTION routine body
)

// Split extracts top-level SQL statements from input.
//
// The returned slice contains one Segment per non-empty statement. Empty
// statements (lone semicolons, comment-only chunks) are filtered out. Split is
// infallible — it always returns a result even for malformed input; lexing
// errors are suppressed here (callers who need them use NewLexer(input).Errors()
// or Parse, which promotes lex errors to diagnostics).
//
// The lexer already hides string/identifier/comment content, so a ';' inside a
// '...' / U&'...' / X'...' literal, a "..." / `...` quoted identifier, or a
// line/bracketed comment is never seen as a delimiter here.
//
// The one Trino construct whose body legitimately contains top-level ';' is an
// inline SQL routine: CREATE [OR REPLACE] FUNCTION ... <controlStatement> and
// WITH FUNCTION ... <controlStatement>. A routine's controlStatement may be a
// bare `RETURN expr` (no ';', no block) or an END-terminated block
// (BEGIN..END, IF..END IF, CASE..END CASE, LOOP..END LOOP, WHILE..END WHILE,
// REPEAT..END REPEAT), and blocks nest. Split tracks routine-body context so
// the ';' separating control statements inside such a body does not split the
// CREATE/WITH FUNCTION into pieces.
//
// A top-level BEGIN is NOT treated as a block opener: in Trino a statement-
// leading BEGIN is a transaction start (BEGIN [WORK|TRANSACTION]); the routine
// BEGIN only ever appears after `FUNCTION ... RETURNS <type>`, i.e. inside a
// function-definition statement, which is exactly the context stateInBody
// tracks.
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0 // END-block nesting depth while in stateInBody

	// isFunctionDef reports whether the current statement is an inline-routine
	// definition whose body we must keep whole. It is decided from the
	// statement's leading keywords and cleared at each statement boundary.
	isFunctionDef := false
	// stmtTokenIdx counts meaningful tokens seen in the current statement, and
	// firstKind records the statement's first token kind; together they drive
	// the position-precise recognition of the `CREATE [OR REPLACE] FUNCTION` /
	// `WITH FUNCTION` prefix, so a FUNCTION keyword appearing LATER in a
	// statement (e.g. a column named `function` in a CTE) or after a different
	// leading keyword (e.g. `DROP FUNCTION`) is never mistaken for a routine
	// body whose internal ';' must be preserved.
	stmtTokenIdx := 0
	firstKind := tokInvalid
	// prevWasEnd is set inside a routine body immediately after an END token so
	// the keyword that completes a typed block closer — END IF / END CASE /
	// END LOOP / END WHILE / END REPEAT — is recognized as the closer's suffix
	// rather than miscounted as a NEW block opener (which would desync depth and
	// swallow the statement-terminating ';' after the routine).
	prevWasEnd := false

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

	// markRoutineIfPrefix sets isFunctionDef when tok at position idx completes
	// an inline-routine prefix that can carry a control-flow body:
	//   WITH FUNCTION ...                 (FUNCTION at idx 1, first kind WITH)
	//   CREATE FUNCTION ...               (FUNCTION at idx 1, first kind CREATE)
	//   CREATE OR REPLACE FUNCTION ...    (FUNCTION at idx 3, first kind CREATE)
	// Only CREATE and WITH can introduce a routine body, so other leading
	// keywords (DROP FUNCTION, SHOW FUNCTIONS) are excluded.
	markRoutineIfPrefix := func(idx int, kind TokenKind) {
		if kind != kwFUNCTION {
			return
		}
		switch {
		case idx == 1 && (firstKind == kwWITH || firstKind == kwCREATE):
			isFunctionDef = true
		case idx == 3 && firstKind == kwCREATE:
			isFunctionDef = true
		}
	}

	for {
		tok := l.NextToken()
		if tok.Kind == tokEOF {
			break
		}

		switch state {
		case stateTop:
			if stmtTokenIdx == 0 {
				firstKind = tok.Kind
			}
			markRoutineIfPrefix(stmtTokenIdx, tok.Kind)

			switch tok.Kind {
			case kwBEGIN, kwIF, kwCASE, kwLOOP, kwWHILE, kwREPEAT:
				// Block openers count only inside a function definition; a
				// top-level BEGIN/IF/CASE outside a routine is either a
				// transaction start (BEGIN) or not a statement opener at all.
				if isFunctionDef {
					state = stateInBody
					depth = 1
				}
			case int(';'):
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
				isFunctionDef = false
				firstKind = tokInvalid
				stmtTokenIdx = -1 // becomes 0 after the increment below
			}

		case stateInBody:
			switch tok.Kind {
			case kwBEGIN, kwIF, kwCASE, kwLOOP, kwWHILE, kwREPEAT:
				// A block keyword right after END is the suffix of a typed
				// closer (END IF / END CASE / …), not a new block opener.
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
						state = stateTop
						// The routine body block is closed. Clear the
						// function-definition flag so any control keyword
						// (e.g. a CASE expression) in the trailing query of a
						// `WITH FUNCTION ... <body> <query>` statement is NOT
						// re-treated as a routine block opener.
						isFunctionDef = false
						prevWasEnd = false
					}
				}
			default:
				prevWasEnd = false
			}
			// ';' and every other token are absorbed inside the routine body.
		}

		stmtTokenIdx++
	}

	// Trailing segment — whatever remains after the last ';' (or the whole
	// input when there was none). ByteEnd is len(input) here.
	if stmtStart < len(input) {
		emit(len(input))
	}

	return segments
}
