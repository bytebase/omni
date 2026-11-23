package parser

// Segment is one segment of a Databricks/Spark SQL script delimited by
// top-level semicolons (outside compound BEGIN ... END blocks). Text is
// script[ByteStart:ByteEnd], preserving leading whitespace and comments, the
// terminating semicolon when present, and any same-line trailing comment after
// the semicolon.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty reports whether the segment holds no Databricks SQL statement (only
// whitespace, comments, and semicolons). Malformed spans such as unclosed
// block comments, strings, or backtick identifiers are non-empty so callers
// such as Parse or ValidateQuery can report the lexical error.
func (s Segment) Empty() bool {
	l := NewLexer(s.Text)
	for {
		tok := l.NextToken()
		if l.Err() != nil {
			return false
		}
		switch tok.Kind {
		case tokEOF:
			return true
		case tokSemicolon:
			continue
		default:
			return false
		}
	}
}

type compoundKind uint8

const (
	compoundBegin compoundKind = iota + 1
	compoundIf
	compoundCase
	compoundWhile
	compoundRepeat
	compoundLoop
	compoundFor
)

type compoundFrame struct {
	kind  compoundKind
	label string // normalized begin label, or empty
}

type casePhase uint8

const (
	caseExpectWhen casePhase = iota
	caseInWhen
	caseAfterThen
)

type caseFrame struct {
	parenDepth int
	phase      casePhase
}

// Split splits a Databricks/Spark SQL script into segments at top-level
// semicolons using the Spark 4.0.5 lexer. Semicolons inside single-quoted
// strings, double-quoted strings, raw strings (R'...' / R"..."),
// backtick-quoted identifiers, line comments (`--` including Spark's `\` + LF
// continuation), nested block comments (`/* /* ... */ */`), optimizer hints
// (`/*+ ... */`), and compound `BEGIN [NOT ATOMIC] ... END` scripting blocks
// (including nested `BEGIN...END`, labeled blocks whose labels may be
// non-reserved keywords such as `loop` or `if`, `IF...END IF`,
// `CASE...END CASE`, `WHILE...END WHILE`, `REPEAT...END REPEAT`,
// `LOOP...END LOOP`, `FOR...END FOR`, and `DECLARE EXIT HANDLER FOR ...
// BEGIN...END`) do not split statements.
//
// Byte spans are contiguous across [0, len(script)) with
// Text == script[ByteStart:ByteEnd]. Empty segments (lone semicolons,
// whitespace, or comment-only chunks) are retained with Empty() == true so
// callers can account for every byte and semicolon. Lexical errors do not abort
// splitting; malformed spans stay in their segment as non-empty so Parse or
// analysis callers can report them at their byte offset.
func Split(script string) []Segment {
	if len(script) == 0 {
		return nil
	}

	l := NewLexer(script)
	var buf []Token
	nextTok := func() Token {
		if len(buf) > 0 {
			t := buf[0]
			buf = buf[1:]
			return t
		}
		return l.NextToken()
	}
	peekTok := func(offset int) Token {
		for len(buf) <= offset {
			buf = append(buf, l.NextToken())
		}
		return buf[offset]
	}

	var segments []Segment
	start := 0
	stmtStart := true
	var stack []compoundFrame
	pendingLabel := ""
	atBodyStmtStart := false
	stmtHead := tokEOF
	prevTokKind := tokEOF
	parenDepth := 0
	var caseStack []caseFrame

	canPrecedeCaseEnd := func(prevKind TokenKind) bool {
		switch prevKind {
		case tokEOF, tokSemicolon, tokComma, tokDot, tokColon, tokDoubleColon, tokLParen, tokLBracket,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat, tokTilde, tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokArrow, tokFatArrow,
			kwWHEN, kwTHEN, kwELSE, kwAND, kwOR, kwNOT, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE,
			kwBETWEEN, kwAS, kwDISTINCT, kwALL, kwSELECT, kwFROM, kwWHERE, kwHAVING, kwON, kwUSING, kwINTO:
			return false
		}
		return true
	}

	isCaseExprEnd := func(prevKind TokenKind) bool {
		if len(caseStack) == 0 {
			return false
		}
		top := caseStack[len(caseStack)-1]
		if parenDepth != top.parenDepth || top.phase != caseAfterThen {
			return false
		}
		return canPrecedeCaseEnd(prevKind)
	}

	isCaseExprStart := func(prevKind TokenKind) bool {
		if prevKind == tokDot {
			return false
		}
		first := peekTok(0).Kind
		if first == kwWHEN {
			if prevKind == kwCASE {
				// Disambiguate `CASE case WHEN ... END` (identifier operand named `case`)
				// from `CASE CASE WHEN ... END WHEN ... END` (nested searched CASE).
				relDepth := 0
				afterThen := false
				for off := 1; ; off++ {
					k := peekTok(off).Kind
					switch k {
					case tokEOF, tokSemicolon:
						return false
					case tokLParen, tokLBracket:
						relDepth++
					case tokRParen, tokRBracket:
						if relDepth == 0 {
							return false
						}
						relDepth--
					case kwWHEN:
						if relDepth == 0 {
							afterThen = false
						}
					case kwTHEN, kwELSE:
						if relDepth == 0 {
							afterThen = true
						}
					case kwEND:
						if relDepth == 0 && afterThen && canPrecedeCaseEnd(peekTok(off-1).Kind) {
							return peekTok(off+1).Kind == kwWHEN
						}
					}
				}
			}
			return true
		}
		switch first {
		case tokEOF, tokSemicolon, tokComma, tokRParen, tokRBracket, tokDot, tokLBracket, tokDoubleColon, tokColon,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat, tokShiftLeft, tokShiftRight, tokShiftRightU,
			kwTHEN, kwDO, kwELSE, kwELSEIF, kwUNTIL, kwFROM, kwWHERE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN, kwAS:
			return false
		}
		// Simple CASE expression (`CASE <operand> WHEN ...`): must reach WHEN at
		// relative paren depth 0 before any enclosing clause boundary.
		relDepth := 0
		for off := 0; ; off++ {
			k := peekTok(off).Kind
			switch k {
			case tokEOF, tokSemicolon:
				return false
			case tokLParen, tokLBracket:
				relDepth++
			case tokRParen, tokRBracket:
				if relDepth == 0 {
					return false
				}
				relDepth--
			case kwWHEN:
				if relDepth == 0 {
					return true
				}
			case kwCASE:
				if relDepth == 0 && off > 0 {
					return false
				}
			case kwEND:
				if relDepth == 0 && off > 0 && canPrecedeCaseEnd(peekTok(off-1).Kind) {
					return false
				}
			case kwTHEN, kwDO, kwELSE, kwELSEIF, kwUNTIL:
				if relDepth == 0 {
					return false
				}
			}
		}
	}

	enterBeginBlock := func(label string) {
		stack = append(stack, compoundFrame{kind: compoundBegin, label: label})
		if peekTok(0).Kind == kwNOT && peekTok(1).Kind == kwATOMIC {
			nextTok()
			nextTok()
		}
		pendingLabel = ""
		atBodyStmtStart = true
		stmtHead = tokEOF
		prevTokKind = tokEOF
		parenDepth = 0
		caseStack = caseStack[:0]
	}

	consumeMatchingEndLabel := func(label string) {
		if label == "" {
			return
		}
		next := peekTok(0)
		if next.Kind.isIdentifier() && normalizeIdentName(next.Text) == label {
			nextTok()
		}
	}

	popControlEnd := func(wantKind compoundKind, wantKw TokenKind) bool {
		if len(stack) == 0 {
			return false
		}
		top := stack[len(stack)-1]
		if top.kind != wantKind || peekTok(0).Kind != wantKw {
			return false
		}
		nextTok() // consume IF/CASE/WHILE/REPEAT/LOOP/FOR
		stack = stack[:len(stack)-1]
		consumeMatchingEndLabel(top.label)
		atBodyStmtStart = false
		stmtHead = tokEOF
		return true
	}

	for {
		tok := nextTok()
		if tok.Kind == tokEOF {
			break
		}
		prevKind := prevTokKind
		prevTokKind = tok.Kind

		// Top-level semicolon outside any compound block.
		if tok.Kind == tokSemicolon && len(stack) == 0 {
			segEnd := extendSameLineTrailingComment(script, tok.End)
			if segEnd > tok.End {
				buf = buf[:0]
				l.pos = segEnd
				l.inHint = false
			}
			segments = append(segments, Segment{
				Text:      script[start:segEnd],
				ByteStart: start,
				ByteEnd:   segEnd,
			})
			start = segEnd
			stmtStart = true
			continue
		}

		if len(stack) == 0 {
			if stmtStart {
				if tok.Kind == kwBEGIN {
					stmtStart = false
					enterBeginBlock("")
					continue
				}
				// Labeled top-level compound block: `label: BEGIN ... END label;`
				if tok.Kind.isIdentifier() && peekTok(0).Kind == tokColon && peekTok(1).Kind == kwBEGIN {
					label := normalizeIdentName(tok.Text)
					nextTok() // consume ':'
					nextTok() // consume 'BEGIN'
					stmtStart = false
					enterBeginBlock(label)
					continue
				}
				stmtStart = false
			}
			continue
		}

		// Inside a compound BEGIN ... END block (len(stack) > 0).
		if tok.Kind == tokSemicolon {
			atBodyStmtStart = true
			pendingLabel = ""
			stmtHead = tokEOF
			prevTokKind = tokEOF
			parenDepth = 0
			caseStack = caseStack[:0]
			continue
		}

		if atBodyStmtStart {
			// Labeled compound block or loop: `label: BEGIN/WHILE/REPEAT/LOOP/FOR`
			if tok.Kind.isIdentifier() && peekTok(0).Kind == tokColon {
				pendingLabel = normalizeIdentName(tok.Text)
				nextTok() // consume ':'
				continue
			}
			switch tok.Kind {
			case kwBEGIN:
				enterBeginBlock(pendingLabel)
				continue
			case kwIF:
				stack = append(stack, compoundFrame{kind: compoundIf})
				pendingLabel = ""
				stmtHead = kwIF
				atBodyStmtStart = false
				continue
			case kwCASE:
				stack = append(stack, compoundFrame{kind: compoundCase})
				pendingLabel = ""
				stmtHead = kwCASE
				atBodyStmtStart = false
				continue
			case kwWHILE:
				stack = append(stack, compoundFrame{kind: compoundWhile, label: pendingLabel})
				pendingLabel = ""
				stmtHead = kwWHILE
				atBodyStmtStart = false
				continue
			case kwREPEAT:
				stack = append(stack, compoundFrame{kind: compoundRepeat, label: pendingLabel})
				pendingLabel = ""
				// Body follows REPEAT immediately without a semicolon.
				continue
			case kwLOOP:
				stack = append(stack, compoundFrame{kind: compoundLoop, label: pendingLabel})
				pendingLabel = ""
				// Body follows LOOP immediately without a semicolon.
				continue
			case kwFOR:
				stack = append(stack, compoundFrame{kind: compoundFor, label: pendingLabel})
				pendingLabel = ""
				stmtHead = kwFOR
				atBodyStmtStart = false
				continue
			case kwELSE:
				// Body follows ELSE immediately without a semicolon.
				pendingLabel = ""
				continue
			case kwDECLARE:
				pendingLabel = ""
				// `DECLARE (EXIT|CONTINUE) HANDLER FOR <conditionValues> <action>`
				// Consume the handler header and condition list so BEGIN is only
				// recognized when it actually begins the handler's action.
				if (peekTok(0).Kind == kwEXIT || peekTok(0).Kind == kwCONTINUE) &&
					peekTok(1).Kind == kwHANDLER && peekTok(2).Kind == kwFOR {
					nextTok() // EXIT / CONTINUE
					nextTok() // HANDLER
					nextTok() // FOR
					for {
						k := peekTok(0).Kind
						if k == kwSQLSTATE {
							if peekTok(1).Kind == kwVALUE && (peekTok(2).Kind == tokString || peekTok(2).Kind == tokDoubleQuoted) {
								nextTok()
								nextTok()
								nextTok()
							} else if peekTok(1).Kind == tokString || peekTok(1).Kind == tokDoubleQuoted {
								nextTok()
								nextTok()
							} else {
								nextTok()
								for peekTok(0).Kind == tokDot && peekTok(1).Kind.isIdentifier() {
									nextTok()
									nextTok()
								}
							}
						} else if k == kwSQLEXCEPTION {
							nextTok()
						} else if k == kwNOT && peekTok(1).Kind == kwFOUND {
							nextTok()
							nextTok()
						} else if k.isIdentifier() {
							nextTok()
							for peekTok(0).Kind == tokDot && peekTok(1).Kind.isIdentifier() {
								nextTok()
								nextTok()
							}
						} else {
							break
						}
						if peekTok(0).Kind != tokComma {
							break
						}
						nextTok() // ','
					}
					if peekTok(0).Kind == kwBEGIN {
						nextTok()
						enterBeginBlock("")
						continue
					}
					if peekTok(0).Kind.isIdentifier() && peekTok(1).Kind == tokColon && peekTok(2).Kind == kwBEGIN {
						label := normalizeIdentName(nextTok().Text)
						nextTok() // ':'
						nextTok() // BEGIN
						enterBeginBlock(label)
						continue
					}
				}
				stmtHead = kwDECLARE
				atBodyStmtStart = false
				parenDepth = 0
				caseStack = caseStack[:0]
				continue
			case kwEND:
				pendingLabel = ""
				top := stack[len(stack)-1]
				switch top.kind {
				case compoundBegin:
					stack = stack[:len(stack)-1]
					consumeMatchingEndLabel(top.label)
					atBodyStmtStart = false
					stmtHead = tokEOF
				case compoundIf:
					popControlEnd(compoundIf, kwIF)
				case compoundCase:
					popControlEnd(compoundCase, kwCASE)
				case compoundWhile:
					popControlEnd(compoundWhile, kwWHILE)
				case compoundRepeat:
					popControlEnd(compoundRepeat, kwREPEAT)
				case compoundLoop:
					popControlEnd(compoundLoop, kwLOOP)
				case compoundFor:
					popControlEnd(compoundFor, kwFOR)
				}
				continue
			default:
				pendingLabel = ""
				stmtHead = tok.Kind
				atBodyStmtStart = false
				parenDepth = 0
				caseStack = caseStack[:0]
				continue
			}
		}

		// Mid-statement tokens inside a compound block.
		switch tok.Kind {
		case tokLParen, tokLBracket:
			parenDepth++
		case tokRParen, tokRBracket:
			if parenDepth > 0 {
				parenDepth--
			}
		case kwCASE:
			// Track nested SQL CASE expressions inside IF/ELSEIF/WHEN/UNTIL headers so
			// `IF CASE WHEN ... THEN ... END THEN` doesn't treat the expression's
			// THEN as the IF body start, while preserving non-reserved identifiers
			// named `case` (e.g. `IF case = 1 THEN ...`).
			if (stmtHead == kwIF || stmtHead == kwELSEIF || stmtHead == kwWHEN || stmtHead == kwCASE || stmtHead == kwUNTIL) &&
				isCaseExprStart(prevKind) {
				caseStack = append(caseStack, caseFrame{parenDepth: parenDepth, phase: caseExpectWhen})
			}
		case kwEND:
			if len(caseStack) > 0 {
				if isCaseExprEnd(prevKind) {
					caseStack = caseStack[:len(caseStack)-1]
				}
			} else if stmtHead == kwUNTIL && parenDepth == 0 {
				// `REPEAT ... UNTIL <cond> END REPEAT [label];` has no semicolon before END REPEAT.
				popControlEnd(compoundRepeat, kwREPEAT)
			}
		case kwWHEN:
			if len(caseStack) > 0 {
				top := &caseStack[len(caseStack)-1]
				if parenDepth == top.parenDepth {
					top.phase = caseInWhen
				}
			} else if stmtHead == kwCASE && parenDepth == 0 {
				stmtHead = kwWHEN
			}
		case kwTHEN:
			if len(caseStack) > 0 {
				top := &caseStack[len(caseStack)-1]
				if parenDepth == top.parenDepth {
					top.phase = caseAfterThen
				}
			} else if (stmtHead == kwIF || stmtHead == kwELSEIF || stmtHead == kwWHEN) && parenDepth == 0 {
				atBodyStmtStart = true
				stmtHead = tokEOF
			}
		case kwELSE:
			if len(caseStack) > 0 {
				top := &caseStack[len(caseStack)-1]
				if parenDepth == top.parenDepth {
					top.phase = caseAfterThen
				}
			}
		case kwDO:
			if (stmtHead == kwWHILE || stmtHead == kwFOR) && parenDepth == 0 && len(caseStack) == 0 {
				atBodyStmtStart = true
				stmtHead = tokEOF
			}
		}
	}

	if start < len(script) {
		segments = append(segments, Segment{
			Text:      script[start:],
			ByteStart: start,
			ByteEnd:   len(script),
		})
	}
	return segments
}

// extendSameLineTrailingComment advances past optional horizontal spaces/tabs
// and a trailing same-line comment after a terminating semicolon:
//   - a `-- ...` line comment (including Spark's `\` + LF continuation), up to
//     (not including) `\r`, `\n`, or EOF; or
//   - a `/* ... */` block comment (including nested block comments and
//     multiline block comments that start on the semicolon's line) PROVIDED
//     only horizontal whitespace follows its closing `*/` before `\r`, `\n`,
//     or EOF (so `SELECT 1; /* lead */ SELECT 2;` keeps `/* lead */` attached
//     to `SELECT 2;`).
//
// If no such trailing comment follows on the same line, it returns pos unchanged.
func extendSameLineTrailingComment(script string, pos int) int {
	n := len(script)
	i := pos
	for i < n && (script[i] == ' ' || script[i] == '\t') {
		i++
	}
	if i+1 < n && script[i] == '-' && script[i+1] == '-' {
		i += 2
		for i < n {
			if script[i] == '\\' && i+1 < n && script[i+1] == '\n' {
				i += 2
				continue
			}
			if script[i] == '\n' || script[i] == '\r' {
				break
			}
			i++
		}
		return i
	}
	if i+1 < n && script[i] == '/' && script[i+1] == '*' && !(i+2 < n && script[i+2] == '+') {
		j := i + 2
		depth := 1
		for j < n && depth > 0 {
			if script[j] == '/' && j+1 < n && script[j+1] == '*' {
				depth++
				j += 2
			} else if script[j] == '*' && j+1 < n && script[j+1] == '/' {
				depth--
				j += 2
			} else {
				j++
			}
		}
		if depth != 0 {
			return pos
		}
		// Only absorb the block comment into the previous statement if nothing
		// else follows on the closing line.
		k := j
		for k < n && (script[k] == ' ' || script[k] == '\t') {
			k++
		}
		if k == n || script[k] == '\n' || script[k] == '\r' {
			return j
		}
	}
	return pos
}
