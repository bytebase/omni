package catalog

import "strings"

type hintTokKind int

const (
	hintTokWord hintTokKind = iota
	hintTokInt
	hintTokString
	hintTokLParen
	hintTokRParen
	hintTokComma
	hintTokOther
)

type hintToken struct {
	kind hintTokKind
	val  string
}

// validateOptimizerHint validates OceanBase 4.3.5.6 semantic rules for
// optimizer hints attached to SELECT (`/*+ ... */`):
//   - LOAD_BATCH_SIZE(int) / LOAD_BATCH_SIZE(int, 'str') is a LOAD DATA hint
//     that OceanBase rejects in SELECT/view queries with 5063 (Unknown hint).
//   - DIRECT(bool, int, 'load_method') requires load_method to be one of
//     'full', 'inc', or 'inc_replace'; any other string fails with 1210
//     (Incorrect arguments to load method in direct hint).
//
// Syntactically malformed hint comments are ignored without error, matching
// OceanBase's hint parser which discards invalid hint blocks with a warning.
func validateOptimizerHint(body string) error {
	toks, ok := tokenizeHintBody(body)
	if !ok || len(toks) == 0 {
		return nil
	}

	var pendingErr error
	for i := 0; i < len(toks); {
		switch toks[i].kind {
		case hintTokComma:
			i++
		case hintTokWord:
			upper := strings.ToUpper(toks[i].val)
			if i+1 < len(toks) && toks[i+1].kind == hintTokLParen {
				closeIdx := findMatchingHintRParen(toks, i+1)
				if closeIdx < 0 {
					return nil
				}
				inner := toks[i+2 : closeIdx]
				switch upper {
				case "LOAD_BATCH_SIZE":
					if !isValidLoadBatchSizeHintArgs(inner) {
						return nil
					}
					if pendingErr == nil {
						pendingErr = errOBUnknownHint()
					}
				case "DIRECT":
					err, validSyntax := validateDirectHintArgs(inner)
					if !validSyntax {
						return nil
					}
					if err != nil && pendingErr == nil {
						pendingErr = err
					}
				}
				i = closeIdx + 1
			} else {
				i++
			}
		default:
			// Unexpected top-level token (such as a bare `(...)` without a hint
			// name) makes the whole hint block syntactically invalid in
			// OceanBase's hint parser, which discards the block with a warning.
			return nil
		}
	}
	return pendingErr
}

// isValidLoadBatchSizeHintArgs reports whether inner matches OceanBase's
// global_hint grammar for LOAD_BATCH_SIZE:
//
//	LOAD_BATCH_SIZE '(' INTNUM ')'
//	LOAD_BATCH_SIZE '(' INTNUM ',' STRING_VALUE ')'
func isValidLoadBatchSizeHintArgs(inner []hintToken) bool {
	if len(inner) == 1 {
		return inner[0].kind == hintTokInt
	}
	if len(inner) == 3 {
		return inner[0].kind == hintTokInt &&
			inner[1].kind == hintTokComma &&
			inner[2].kind == hintTokString
	}
	return false
}

// validateDirectHintArgs checks OceanBase's global_hint grammar for DIRECT:
//
//	DIRECT '(' BOOL_VALUE ',' INTNUM ')'
//	DIRECT '(' BOOL_VALUE ',' INTNUM ',' STRING_VALUE ')'
//
// Returns (semanticErr, validSyntax). When validSyntax is false, the hint
// comment is syntactically malformed and discarded by OceanBase.
func validateDirectHintArgs(inner []hintToken) (error, bool) {
	if len(inner) == 3 {
		if isHintBoolToken(inner[0]) && inner[1].kind == hintTokComma && inner[2].kind == hintTokInt {
			return nil, true
		}
		return nil, false
	}
	if len(inner) == 5 {
		if isHintBoolToken(inner[0]) &&
			inner[1].kind == hintTokComma &&
			inner[2].kind == hintTokInt &&
			inner[3].kind == hintTokComma &&
			inner[4].kind == hintTokString {
			switch strings.ToLower(inner[4].val) {
			case "full", "inc", "inc_replace":
				return nil, true
			default:
				return &Error{
					Code:     ErrWrongArguments,
					SQLState: sqlState(ErrWrongArguments),
					Message:  "Incorrect arguments to load method in direct hint",
				}, true
			}
		}
		return nil, false
	}
	return nil, false
}

func isHintBoolToken(tok hintToken) bool {
	return tok.kind == hintTokWord &&
		(strings.EqualFold(tok.val, "TRUE") || strings.EqualFold(tok.val, "FALSE"))
}

func findMatchingHintRParen(toks []hintToken, openIdx int) int {
	depth := 0
	for i := openIdx; i < len(toks); i++ {
		switch toks[i].kind {
		case hintTokLParen:
			depth++
		case hintTokRParen:
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// tokenizeHintBody lexes the body of a `/*+ ... */` hint comment and checks
// that parentheses and string literals are balanced. Returns (nil, false) if
// the hint body is lexically malformed.
func tokenizeHintBody(s string) ([]hintToken, bool) {
	var toks []hintToken
	depth := 0
	for i := 0; i < len(s); {
		ch := s[i]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' || ch == '\f' || ch == '\v' {
			i++
			continue
		}
		switch {
		case ch == '(':
			depth++
			toks = append(toks, hintToken{kind: hintTokLParen, val: "("})
			i++
		case ch == ')':
			depth--
			if depth < 0 {
				return nil, false
			}
			toks = append(toks, hintToken{kind: hintTokRParen, val: ")"})
			i++
		case ch == ',':
			toks = append(toks, hintToken{kind: hintTokComma, val: ","})
			i++
		case ch == '\'' || ch == '"':
			quote := ch
			i++
			var sb strings.Builder
			closed := false
			for i < len(s) {
				c := s[i]
				if c == '\\' && i+1 < len(s) {
					sb.WriteByte(s[i+1])
					i += 2
					continue
				}
				if c == quote {
					if i+1 < len(s) && s[i+1] == quote {
						sb.WriteByte(quote)
						i += 2
						continue
					}
					i++
					closed = true
					break
				}
				sb.WriteByte(c)
				i++
			}
			if !closed {
				return nil, false
			}
			toks = append(toks, hintToken{kind: hintTokString, val: sb.String()})
		case ch >= '0' && ch <= '9':
			start := i
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				i++
			}
			toks = append(toks, hintToken{kind: hintTokInt, val: s[start:i]})
		case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_':
			start := i
			for i < len(s) && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= '0' && s[i] <= '9') || s[i] == '_' || s[i] == '$') {
				i++
			}
			toks = append(toks, hintToken{kind: hintTokWord, val: s[start:i]})
		default:
			toks = append(toks, hintToken{kind: hintTokOther, val: s[i : i+1]})
			i++
		}
	}
	if depth != 0 {
		return nil, false
	}
	return toks, true
}
