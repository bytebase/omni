package completion

import (
	"github.com/bytebase/omni/pg/yacc"
)

const maxReduceDepth = 20

// collectValidTokens returns the set of token IDs that are valid at the given state stack.
// It explores shift actions directly available and follows reduce chains to discover more.
func collectValidTokens(stack []int) map[int]bool {
	tokens := make(map[int]bool)
	visited := make(map[stateStackKey]bool)
	collectFromStack(stack, tokens, visited, 0)
	return tokens
}

// stateStackKey is used to avoid revisiting the same state stack configuration.
type stateStackKey struct {
	topState int
	depth    int
}

func collectFromStack(stack []int, tokens map[int]bool, visited map[stateStackKey]bool, depth int) {
	if depth > maxReduceDepth || len(stack) == 0 {
		return
	}

	state := stack[len(stack)-1]
	key := stateStackKey{topState: state, depth: len(stack)}
	if visited[key] {
		return
	}
	visited[key] = true

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
	ntokens := yacc.NTokens()

	const TOKSTART = 4 // skip $end, error, $unk, and the empty token name

	// 1. Collect direct shifts from this state
	base := int(pact[state])
	if int32(base) > flag {
		for tok := TOKSTART; tok-1 < ntokens; tok++ {
			n := base + tok
			if n >= 0 && n < last {
				if int(chk[int(act[n])]) == tok {
					tokens[tok] = true
				}
			}
		}
	}

	// 2. Also check exception table
	if def[state] == -2 {
		i := 0
		for i < len(exca)-1 {
			if exca[i] == -1 && int(exca[i+1]) == state {
				break
			}
			i += 2
		}
		for i += 2; i < len(exca)-1 && exca[i] >= 0; i += 2 {
			tok := int(exca[i])
			if tok >= TOKSTART && exca[i+1] != 0 {
				tokens[tok] = true
			}
		}
	}

	// 3. Follow default reduction to discover more tokens
	prod := int(def[state])
	if prod > 0 {
		// Simulate the reduce and recurse
		rhsLen := int(r2[prod])
		if rhsLen > len(stack)-1 {
			return
		}
		newStack := make([]int, len(stack)-rhsLen)
		copy(newStack, stack[:len(stack)-rhsLen])

		lhs := int(r1[prod])
		gotoBase := int(pgo[lhs])
		topState := newStack[len(newStack)-1]
		gotoIdx := gotoBase + topState + 1
		var newState int
		if gotoIdx >= last {
			newState = int(act[gotoBase])
		} else {
			newState = int(act[gotoIdx])
			if int(chk[newState]) != -lhs {
				newState = int(act[gotoBase])
			}
		}
		newStack = append(newStack, newState)
		collectFromStack(newStack, tokens, visited, depth+1)
	}

	// 4. Also explore exception table reductions
	if def[state] == -2 {
		i := 0
		for i < len(exca)-1 {
			if exca[i] == -1 && int(exca[i+1]) == state {
				break
			}
			i += 2
		}
		// Check the default action at the end of the exception entries
		for i += 2; i < len(exca)-1 && exca[i] >= 0; i += 2 {
			// skip individual token entries
		}
		if i < len(exca) {
			excaDef := int(exca[i+1])
			if excaDef > 0 {
				rhsLen := int(r2[excaDef])
				if rhsLen <= len(stack)-1 {
					newStack := make([]int, len(stack)-rhsLen)
					copy(newStack, stack[:len(stack)-rhsLen])

					lhs := int(r1[excaDef])
					gotoBase := int(pgo[lhs])
					topState := newStack[len(newStack)-1]
					gotoIdx := gotoBase + topState + 1
					var newState int
					if gotoIdx >= last {
						newState = int(act[gotoBase])
					} else {
						newState = int(act[gotoIdx])
						if int(chk[newState]) != -lhs {
							newState = int(act[gotoBase])
						}
					}
					newStack = append(newStack, newState)
					collectFromStack(newStack, tokens, visited, depth+1)
				}
			}
		}
	}
}
