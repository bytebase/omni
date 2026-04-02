package completion

import "github.com/bytebase/omni/mongo/parser"

// completionContext identifies the kind of completion expected.
type completionContext int

const (
	contextTopLevel      completionContext = iota // start of input or after semicolon
	contextAfterDbDot                             // db.|
	contextAfterCollDot                           // db.users.|
	contextAfterBracket                           // db[|
	contextInsideArgs                             // db.users.find(|
	contextDocumentKey                            // {| or {age: 1, |
	contextQueryOperator                          // {age: {$|
	contextAggStage                               // [{$|
	contextCursorChain                            // db.users.find().|
	contextShowTarget                             // show |
	contextAfterRsDot                             // rs.|
	contextAfterShDot                             // sh.|
)

// detectContext analyzes the token sequence to determine the completion context.
func detectContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	if n == 0 {
		return contextTopLevel
	}

	last := tokens[n-1]

	// Ends with semicolon → top level.
	if last.Str == ";" {
		return contextTopLevel
	}

	// Ends with "." → classify dot context.
	if last.Str == "." {
		return classifyDotContext(tokens[:n-1])
	}

	// Ends with "show" keyword → show target.
	if last.Str == "show" {
		return contextShowTarget
	}

	// Ends with "[" → check if preceded by "db".
	if last.Str == "[" {
		if n >= 2 && tokens[n-2].Str == "db" {
			return contextAfterBracket
		}
		return contextTopLevel
	}

	// Ends with a string token after "db[" → bracket access with partial/full collection name.
	// Handles: db["us (unterminated string) and db["users" (complete string).
	if last.Type == parser.TokString {
		if n >= 3 && tokens[n-2].Str == "[" && tokens[n-3].Str == "db" {
			return contextAfterBracket
		}
	}

	// Ends with "(" → inside args.
	if last.Str == "(" {
		return contextInsideArgs
	}

	// Ends with "{" → classify brace context.
	if last.Str == "{" {
		return classifyBraceContext(tokens[:n-1])
	}

	// Ends with "," or ":" → check if inside unclosed brace.
	if last.Str == "," || last.Str == ":" {
		if insideBrace(tokens[:n-1]) {
			return contextDocumentKey
		}
		return contextTopLevel
	}

	// Ends with ")" followed by nothing — this shouldn't produce completions
	// in most cases, but we handle it as top level.
	if last.Str == ")" {
		return contextTopLevel
	}

	return contextTopLevel
}

// classifyDotContext determines the context when the last token is ".".
// tokens is the slice WITHOUT the trailing ".".
func classifyDotContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	if n == 0 {
		return contextTopLevel
	}

	last := tokens[n-1]

	// "db." → afterDbDot
	if last.Str == "db" && n == 1 {
		return contextAfterDbDot
	}
	if last.Str == "db" {
		// Check that what's before "db" isn't a dot (which would mean db is a property).
		prev := tokens[n-2]
		if prev.Str != "." {
			return contextAfterDbDot
		}
	}

	// "rs." → afterRsDot
	if last.Str == "rs" {
		if n == 1 || tokens[n-2].Str != "." {
			return contextAfterRsDot
		}
	}

	// "sh." → afterShDot
	if last.Str == "sh" {
		if n == 1 || tokens[n-2].Str != "." {
			return contextAfterShDot
		}
	}

	// Ends with ")" → could be cursor chain or getCollection(...).
	if last.Str == ")" {
		// Try to find the matching "(" and check if it's db.getCollection("x").
		openIdx := findMatchingOpen(tokens, n-1, "(", ")")
		if openIdx >= 0 {
			// Check for db.getCollection("x").
			if openIdx >= 2 &&
				tokens[openIdx-1].Str == "getCollection" &&
				tokens[openIdx-2].Str == "." {
				// Check if preceded by "db".
				if openIdx >= 3 && tokens[openIdx-3].Str == "db" {
					return contextAfterCollDot
				}
			}
		}
		// General case: method().
		return contextCursorChain
	}

	// Ends with "]" → check for db["coll"].
	if last.Str == "]" {
		openIdx := findMatchingOpen(tokens, n-1, "[", "]")
		if openIdx >= 1 && tokens[openIdx-1].Str == "db" {
			return contextAfterCollDot
		}
		return contextCursorChain
	}

	// Ends with a word → check if it's a collection name after "db.".
	if last.IsWord() {
		// db . <word> . → afterCollDot
		if n >= 3 && tokens[n-2].Str == "." && tokens[n-3].Str == "db" {
			// Make sure "db" isn't itself after a dot.
			if n == 3 || tokens[n-4].Str != "." {
				return contextAfterCollDot
			}
		}
	}

	return contextCursorChain
}

// classifyBraceContext determines the context when the last token is "{".
// tokens is the slice WITHOUT the trailing "{".
func classifyBraceContext(tokens []parser.Token) completionContext {
	n := len(tokens)
	if n == 0 {
		return contextDocumentKey
	}

	last := tokens[n-1]

	// Preceded by ":" → query operator context (nested document for operator).
	if last.Str == ":" {
		return contextQueryOperator
	}

	// Preceded by "[" → agg stage context (pipeline array).
	if last.Str == "[" {
		return contextAggStage
	}

	// Preceded by "," → need to check if we're inside an array (agg pipeline).
	if last.Str == "," {
		if insideArray(tokens[:n-1]) {
			return contextAggStage
		}
		return contextDocumentKey
	}

	return contextDocumentKey
}

// findMatchingOpen walks backward from pos to find the matching open delimiter.
// pos should point to the close delimiter. Returns the index of the matching open,
// or -1 if not found.
func findMatchingOpen(tokens []parser.Token, pos int, open, close string) int {
	depth := 0
	for i := pos; i >= 0; i-- {
		if tokens[i].Str == close {
			depth++
		} else if tokens[i].Str == open {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// insideBrace checks if there is an unclosed "{" in tokens.
func insideBrace(tokens []parser.Token) bool {
	depth := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].Str {
		case "}":
			depth++
		case "{":
			if depth == 0 {
				return true
			}
			depth--
		}
	}
	return false
}

// insideArray checks if there is an unclosed "[" in tokens.
func insideArray(tokens []parser.Token) bool {
	depth := 0
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].Str {
		case "]":
			depth++
		case "[":
			if depth == 0 {
				return true
			}
			depth--
		}
	}
	return false
}
