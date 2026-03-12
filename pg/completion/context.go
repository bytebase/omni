package completion

import (
	"github.com/bytebase/omni/pg/yacc"
)

// GrammarContext represents the grammar rule context at the cursor position.
type GrammarContext int

const (
	CtxNone        GrammarContext = iota
	CtxColumnRef                  // cursor is in a column reference position
	CtxRelationRef                // cursor is in a table/view/sequence reference position
	CtxFuncName                   // cursor is in a function name position
	CtxTypeName                   // cursor is in a type name position
)

// inferContexts determines what grammar contexts are active at the cursor position.
// It explores reachable states (via reduce chains from the stack top), tries to
// shift IDENT from each, then traces the post-IDENT reduce chain to classify nonterminals.
func inferContexts(stack []int) []GrammarContext {
	if len(stack) == 0 {
		return nil
	}

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
	toknames := yacc.TokNames()

	identTok := findTokenID(toknames, "IDENT")
	if identTok == 0 {
		return nil
	}

	seen := make(map[GrammarContext]bool)
	var contexts []GrammarContext
	addCtx := func(ctx GrammarContext) {
		if ctx != CtxNone && !seen[ctx] {
			seen[ctx] = true
			contexts = append(contexts, ctx)
		}
	}

	// Explore reachable stacks (same logic as collector's reduce chain exploration)
	// and try shifting IDENT from each.
	visited := make(map[int]bool)
	type stackEntry struct {
		stack []int
		depth int
	}
	queue := []stackEntry{{stack: stack, depth: 0}}

	for len(queue) > 0 {
		entry := queue[0]
		queue = queue[1:]

		if entry.depth > maxReduceDepth || len(entry.stack) == 0 {
			continue
		}
		state := entry.stack[len(entry.stack)-1]
		if visited[state] {
			continue
		}
		visited[state] = true

		// Try shifting IDENT from this state
		base := int(pact[state])
		if int32(base) > flag {
			n := base + identTok
			if n >= 0 && n < last {
				nn := int(act[n])
				if int(chk[nn]) == identTok {
					// IDENT can be shifted here → trace reduce chain
					postStack := make([]int, len(entry.stack)+1)
					copy(postStack, entry.stack)
					postStack[len(entry.stack)] = nn

					prods := collectPossibleProductions(nn, def, exca)
					for _, prod := range prods {
						traceReduceChain(postStack, prod, r1, r2, pgo, act, chk, def, exca, last, addCtx)
					}
				}
			}
		}

		// Follow reductions to discover more reachable states
		prod := int(def[state])
		if prod == -2 {
			prod = lookupExcaDefault(exca, state)
		}
		if prod > 0 {
			lhs := int(r1[prod])
			rhsLen := int(r2[prod])
			if rhsLen <= len(entry.stack)-1 {
				newStack := make([]int, len(entry.stack)-rhsLen)
				copy(newStack, entry.stack[:len(entry.stack)-rhsLen])
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
				queue = append(queue, stackEntry{stack: newStack, depth: entry.depth + 1})
			}
		}
	}

	return contexts
}

// collectPossibleProductions collects all possible reduction rule numbers from a state.
// For states with def > 0, returns just the default production.
// For states with def == -2 (exception table), returns all productions from exception entries.
func collectPossibleProductions(state int, def []int16, exca []int16) []int {
	d := int(def[state])
	if d > 0 {
		return []int{d}
	}
	if d != -2 {
		return nil
	}

	// Scan exception table for all reductions in this state.
	// Format: (-1, state), (tok, action)..., (-2, default_action), (-1, next_state)...
	prodSet := make(map[int]bool)
	i := 0
	for i < len(exca)-1 {
		if exca[i] == -1 && int(exca[i+1]) == state {
			break
		}
		i += 2
	}
	// Skip past the (-1, state) header
	i += 2
	for i < len(exca)-1 {
		tok := exca[i]
		action := int(exca[i+1])
		if tok == -1 {
			break // next state
		}
		if tok == -2 {
			if action > 0 {
				prodSet[action] = true
			}
			i += 2
			continue
		}
		if action > 0 {
			prodSet[action] = true
		}
		i += 2
	}

	prods := make([]int, 0, len(prodSet))
	for p := range prodSet {
		prods = append(prods, p)
	}
	return prods
}

// traceReduceChain follows a reduce chain from the given stack and initial production,
// classifying each nonterminal encountered.
func traceReduceChain(
	stack []int, startProd int,
	r1 []int16, r2 []int8, pgo []int16,
	act []int16, chk []int16, def []int16, exca []int16,
	last int, addCtx func(GrammarContext),
) {
	testStack := make([]int, len(stack))
	copy(testStack, stack)
	prod := startProd

	for depth := 0; depth < maxReduceDepth; depth++ {
		if prod <= 0 || prod >= len(ruleLHS) {
			break
		}

		name := ruleLHS[prod]
		addCtx(classifyNonterminal(name))

		lhs := int(r1[prod])
		rhsLen := int(r2[prod])
		if rhsLen > len(testStack)-1 {
			break
		}
		testStack = testStack[:len(testStack)-rhsLen]

		gotoBase := int(pgo[lhs])
		topState := testStack[len(testStack)-1]
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
		testStack = append(testStack, newState)

		curState := testStack[len(testStack)-1]
		d := int(def[curState])
		if d > 0 {
			prod = d
		} else if d == -2 {
			prod = lookupExcaDefault(exca, curState)
		} else {
			break
		}
	}
}

// lookupExcaDefault finds the default reduction in the exception table for a state.
// The default entry uses token=-2: (-2, default_action).
func lookupExcaDefault(exca []int16, state int) int {
	i := 0
	for i < len(exca)-1 {
		if exca[i] == -1 && int(exca[i+1]) == state {
			break
		}
		i += 2
	}
	i += 2
	for i < len(exca)-1 {
		tok := exca[i]
		if tok == -1 {
			break
		}
		if tok == -2 {
			return int(exca[i+1])
		}
		i += 2
	}
	return 0
}

// classifyNonterminal maps a nonterminal name to a GrammarContext.
func classifyNonterminal(name string) GrammarContext {
	switch name {
	case "columnref", "ColId":
		return CtxColumnRef
	case "relation_expr", "qualified_name", "table_ref",
		"relation_expr_opt_alias", "insert_target":
		return CtxRelationRef
	case "func_name", "func_application", "func_expr_common_subexpr":
		return CtxFuncName
	case "Typename", "SimpleTypename", "GenericType":
		return CtxTypeName
	}
	return CtxNone
}

// findTokenID finds the internal token ID for a given token name.
func findTokenID(toknames []string, name string) int {
	for i, n := range toknames {
		if n == name {
			return i + 1
		}
	}
	return 0
}
