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
// CRITICAL — a block opener is recognized only at a STATEMENT-LEADING position.
// IF/CASE/LOOP/WHILE/REPEAT/FOR are ordinary syntax mid-statement (a function
// call `IF(x,1,2)`, a query clause `FOR UPDATE`, a guard `IF NOT EXISTS`, a
// keyword-as-alias, or a CASE expression), so opening a block on ANY occurrence
// would flip a plain `SELECT IF(x,1,2); SELECT 1` into block mode and absorb the
// next statement's ';'. These keywords therefore open a block only when they
// begin a statement (input start, or right after a top-level ';'); everywhere
// else they stay flat. This also keeps a stray top-level closer suffix (the IF
// in `END IF`) from reopening a block, since it is never statement-leading.
//
// BEGIN is the exception: it opens a block position-INDEPENDENTLY because the
// canonical `CREATE PROCEDURE … BEGIN` is mid-statement. A top-level BEGIN is
// still ambiguous (BEGIN [TRANSACTION|WORK] is a TCL transaction start, BEGIN …
// END is a block), so one-token lookahead after BEGIN disambiguates: a TCL
// follower (';', EOF, TRANSACTION, READ, ISOLATION) keeps it flat; anything else
// opens a block. (CASE, which is also an expression, no longer needs the
// "harmless closer-counting" argument the legacy split_begin_end_test.go relied
// on — a mid-statement CASE expression simply never opens a block now.)
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	l := NewLexer(input)
	var segments []Segment
	stmtStart := 0
	state := stateTop
	depth := 0 // block-nesting depth while in stateInBlock

	// atStmtStart reports whether the next TOP-LEVEL token begins a new
	// statement — true at input start and immediately after a top-level ';'. It
	// is the gate for the control-flow block OPENERS at the top level: IF / CASE
	// / LOOP / WHILE / REPEAT / FOR open a procedural block ONLY when they are
	// statement-leading. Mid-statement the SAME keywords are ordinary syntax — a
	// function call (IF(x,1,2)), a query clause (FOR UPDATE), a guard (IF NOT
	// EXISTS), a keyword-as-alias, or a CASE/IF expression — and must NOT flip
	// the splitter into block mode (which would absorb the next statement's
	// terminating ';'). It also makes a stray top-level closer suffix (the IF in
	// `END IF`) a no-op, since that suffix is never statement-leading.
	//
	// NOTE the scope: atStmtStart gates only the TOP-LEVEL openers. Inside an
	// already-open block body (stateInBlock) every BEGIN/IF/CASE/LOOP/WHILE/
	// REPEAT/FOR still counts as a nesting opener and every END as a closer (the
	// classic balanced count). That counting is correct for well-formed nesting
	// and self-balances for a CASE expression (CASE … END), which is all the
	// Spanner consumer of block-aware Split needs — Spanner has no stored
	// procedures, so a procedure body with a bare mid-statement IF(...)/FOR
	// clause is invalid Spanner input and out of scope here. Applying the
	// statement-position rule INSIDE a block is unsafe without full parsing,
	// because THEN/ELSE introduce a statement list in a procedural IF/CASE but a
	// value expression in a CASE *expression* (e.g. the IF(...) in
	// `CASE WHEN x THEN IF(a,b,c) ELSE 0 END` must NOT open a block).
	atStmtStart := true

	// prevWasEnd is set immediately after an END token INSIDE a block so a block
	// keyword that completes a typed closer — END IF / END CASE / END LOOP /
	// END WHILE / END REPEAT / END FOR — is recognized as the closer's suffix
	// rather than miscounted as a NEW (nested) block opener, which would desync
	// depth and swallow the statement-terminating ';' after the block. (At the
	// TOP level the statement-position gate above already prevents a closer
	// suffix from reopening, so prevWasEnd is purely a stateInBlock concern.)
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
			// A statement-leading `label:` prefix introduces a labeled procedural
			// block (`label: BEGIN|LOOP|WHILE|REPEAT|FOR …`, GoogleSQLParser.g4
			// unterminated_script_statement's labeled alternative). The label name
			// and its ':' must stay TRANSPARENT to statement-leading position so the
			// block opener that follows still opens a block — otherwise the block's
			// internal ';' would split the labeled statement. Consume the ':' here
			// and keep atStmtStart true. (Only a non-reserved identifier can be a
			// bare label; a reserved word would be backtick-quoted → tokIdentifier.)
			if atStmtStart && isIdentifierStart(tok.Type) && peekToken().Type == int(':') {
				nextToken() // consume ':'
				// atStmtStart stays true: the next token (the block opener) is still
				// statement-leading.
				continue
			}
			switch tok.Type {
			case kwBEGIN:
				// BEGIN opens a procedural block whenever the following token is
				// not a TCL follower — this is position-INDEPENDENT, because
				// `CREATE PROCEDURE … BEGIN` is mid-statement yet must open a
				// block, while a top-level `BEGIN [TRANSACTION]` is a flat
				// transaction start. (A bare kwBEGIN never appears mid-statement
				// except as a block opener; an identifier spelled `begin` lexes
				// to tokIdentifier, not kwBEGIN.) The buffered follower is
				// returned by the next nextToken() and handled normally.
				if !isTCLBeginFollower(peekToken()) {
					state = stateInBlock
					depth = 1
					prevWasEnd = false
				}
				atStmtStart = false
			case kwIF, kwCASE, kwLOOP, kwWHILE, kwREPEAT, kwFOR:
				// A control-flow keyword opens a procedural block ONLY when it is
				// statement-leading. Mid-statement the same keyword is ordinary
				// syntax (IF(...) call, FOR UPDATE clause, IF NOT EXISTS guard,
				// keyword-as-alias, CASE expression) and must stay flat so the
				// next ';' still splits. A stray top-level closer suffix (the IF
				// in `END IF`) is likewise not statement-leading, so it cannot
				// reopen a block.
				if atStmtStart {
					state = stateInBlock
					depth = 1
					prevWasEnd = false
				}
				atStmtStart = false
			case int(';'):
				emit(tok.Loc.Start)
				stmtStart = tok.Loc.End
				atStmtStart = true
			default:
				atStmtStart = false
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
						// Return to stateTop. atStmtStart is still false (it was
						// cleared when this block opened and stateInBlock never
						// sets it), so a typed closer's suffix keyword that
						// follows (END IF / END WHILE / …) is correctly seen as
						// non-statement-leading by the stateTop arm above and
						// cannot reopen a block. The next top-level ';' is the
						// only thing that re-arms atStmtStart.
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
