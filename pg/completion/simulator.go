// Package completion provides SQL auto-completion using pgparser's LALR parse tables.
package completion

import (
	"github.com/bytebase/omni/pg/yacc"
)

// parseState holds the LALR state stack after simulating a parse.
type parseState struct {
	stack []int // state numbers
}

// simulateParse feeds tokens through the LALR state machine without executing
// semantic actions. It returns the state stack at the point where tokens run out.
// This closely follows the parse loop in yacc.go (pgParserImpl.Parse).
func simulateParse(tokens []tokenInfo) parseState {
	pact := yacc.Pact()
	act := yacc.Act()
	chk := yacc.Chk()
	def := yacc.Def()
	r1 := yacc.R1()
	r2 := yacc.R2()
	pgo := yacc.Pgo()
	exca := yacc.Exca()
	flag := int32(yacc.Flag())
	last := yacc.Last()
	eofCode := yacc.EofCode()

	// State stack
	stack := make([]int, 1, 64)
	stack[0] = 0

	state := 0
	tokenIdx := 0
	token := -1 // -1 means "no token read yet"
	errCount := 0
	const maxErrors = 50

	for {
		// Push state
		if len(stack) > 10000 {
			// Safety: prevent unbounded stack growth
			return parseState{stack: stack}
		}

		// ---- pgnewstate ----
		pgn := int(pact[state])
		if pgn <= int(flag) {
			goto pgdefault
		}

		// Need a token?
		if token < 0 {
			if tokenIdx < len(tokens) {
				token = tokens[tokenIdx].parserToken
				tokenIdx++
			} else {
				// No more tokens — this is where we want to stop for completion
				return parseState{stack: stack}
			}
		}

		pgn += token
		if pgn < 0 || pgn >= last {
			goto pgdefault
		}
		pgn = int(act[pgn])
		if int(chk[pgn]) == token {
			// Valid shift
			token = -1 // consume token
			state = pgn
			// Push new state onto stack
			stack = append(stack, state)
			errCount = 0
			continue
		}

	pgdefault:
		pgn = int(def[state])
		if pgn == -2 {
			// Exception table: need a token first
			if token < 0 {
				if tokenIdx < len(tokens) {
					token = tokens[tokenIdx].parserToken
					tokenIdx++
				} else {
					return parseState{stack: stack}
				}
			}

			// Find state in exception table
			xi := 0
			for xi < len(exca)-1 {
				if exca[xi] == -1 && int(exca[xi+1]) == state {
					break
				}
				xi += 2
			}
			// Scan for matching token or default
			xi += 2
			for xi < len(exca)-1 {
				pgn = int(exca[xi])
				if pgn < 0 || pgn == token {
					break
				}
				xi += 2
			}
			pgn = int(exca[xi+1])
			if pgn < 0 {
				// Accept
				return parseState{stack: stack}
			}
		}

		if pgn == 0 {
			// Error
			errCount++
			if errCount > maxErrors {
				return parseState{stack: stack}
			}

			// Error recovery: try to find a state that accepts "error" token
			errCode := yacc.ErrCode()
			recovered := false
			for sp := len(stack) - 1; sp >= 0; sp-- {
				en := int(pact[stack[sp]]) + errCode
				if en >= 0 && en < last {
					es := int(act[en])
					if int(chk[es]) == errCode {
						// Found recovery state — truncate stack and shift error
						stack = stack[:sp+1]
						state = es
						stack = append(stack, state)
						// Discard current token
						if token != eofCode {
							token = -1
						}
						recovered = true
						break
					}
				}
			}
			if !recovered {
				// Can't recover — discard token and try again
				if token == eofCode {
					return parseState{stack: stack}
				}
				token = -1
			}
			continue
		}

		// Reduction by production pgn
		rhsLen := int(r2[pgn])
		if rhsLen > len(stack) {
			// Stack underflow — return what we have
			return parseState{stack: stack}
		}
		// Pop rhsLen states
		stack = stack[:len(stack)-rhsLen]
		if len(stack) == 0 {
			return parseState{stack: []int{0}}
		}

		// Goto
		lhs := int(r1[pgn])
		gotoBase := int(pgo[lhs])
		topState := stack[len(stack)-1]
		gotoIdx := gotoBase + topState + 1
		if gotoIdx >= last {
			state = int(act[gotoBase])
		} else {
			state = int(act[gotoIdx])
			if int(chk[state]) != -lhs {
				state = int(act[gotoBase])
			}
		}
		stack = append(stack, state)
		// Don't consume token — loop back to check new state
	}
}
