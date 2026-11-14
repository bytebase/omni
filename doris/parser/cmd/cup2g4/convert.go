package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// A ctok is a token of a CUP file, with the line it starts on.
type ctok struct {
	text string
	line int
}

// lexCUP splits a CUP file into names and punctuation, dropping comments
// and the {: ... :} code blocks of actions and parser code.
func lexCUP(src string) ([]ctok, error) {
	var out []ctok
	line := 1
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r' || c == '\f':
			i++
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i+2:], "*/")
			if end < 0 {
				return nil, fmt.Errorf("line %d: unterminated comment", line)
			}
			line += strings.Count(src[i:i+2+end+2], "\n")
			i += 2 + end + 2
		case strings.HasPrefix(src[i:], "//"):
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				end = len(src) - i
			}
			i += end
		case strings.HasPrefix(src[i:], "{:"):
			end := strings.Index(src[i+2:], ":}")
			if end < 0 {
				return nil, fmt.Errorf("line %d: unterminated {: block", line)
			}
			line += strings.Count(src[i:i+2+end+2], "\n")
			i += 2 + end + 2
		case strings.HasPrefix(src[i:], "::="):
			out = append(out, ctok{"::=", line})
			i += 3
		case c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i
			for j < len(src) && (src[j] == '_' || src[j] >= 'a' && src[j] <= 'z' || src[j] >= 'A' && src[j] <= 'Z' || src[j] >= '0' && src[j] <= '9') {
				j++
			}
			out = append(out, ctok{src[i:j], line})
			i = j
		default:
			out = append(out, ctok{string(c), line})
			i++
		}
	}
	return out, nil
}

// A rule is one CUP production list: its name, the line it starts on, and
// its alternatives, each a list of symbols.
type rule struct {
	name string
	line int
	alts [][]string
}

// readRules returns the rules of a CUP file in order, and its start symbol.
// Everything else (imports, code blocks, declarations, precedence) is
// skipped: a rule is the only statement with ::= second.
func readRules(toks []ctok) ([]rule, string, map[string]bool, error) {
	var rules []rule
	start := ""
	declared := map[string]bool{}
	for i := 0; i < len(toks); {
		end := i
		for end < len(toks) && toks[end].text != ";" {
			end++
		}
		stmt := toks[i:end]
		i = end + 1
		switch {
		case len(stmt) == 3 && stmt[0].text == "start" && stmt[1].text == "with":
			start = stmt[2].text
		case len(stmt) > 0 && stmt[0].text == "terminal":
			// terminal [Type] A, B, ...: each name before a comma or the end.
			for k := 1; k < len(stmt); k++ {
				if isName(stmt[k].text) && (k+1 == len(stmt) || stmt[k+1].text == ",") {
					declared[stmt[k].text] = true
				}
			}
		case len(stmt) >= 2 && stmt[1].text == "::=":
			r := rule{name: stmt[0].text, line: stmt[0].line}
			var alt []string
			for k := 2; k < len(stmt); k++ {
				switch t := stmt[k].text; {
				case t == "|":
					r.alts = append(r.alts, alt)
					alt = nil
				case t == ":":
					k++ // a label: name:label
				case t == "%" && k+1 < len(stmt) && stmt[k+1].text == "prec":
					k += 2 // %prec TOKEN
				case isName(t):
					alt = append(alt, t)
				default:
					return nil, "", nil, fmt.Errorf("line %d: %q in rule %s", stmt[k].line, t, r.name)
				}
			}
			r.alts = append(r.alts, alt)
			rules = append(rules, r)
		}
	}
	if start == "" {
		return nil, "", nil, fmt.Errorf("no start with")
	}
	return rules, start, declared, nil
}

func isName(t string) bool {
	return t != "" && (t[0] == '_' || t[0] >= 'a' && t[0] <= 'z' || t[0] >= 'A' && t[0] <= 'Z')
}

var (
	// keywordMap.put("set", new Integer(SqlParserSymbols.KW_SET));
	flexKeyword = regexp.MustCompile(`keywordMap\.put\("([^"]+)",\s*new Integer\(SqlParserSymbols\.(\w+)\)\)`)
	// "," { return newToken(SqlParserSymbols.COMMA, null); }
	flexPunct = regexp.MustCompile(`(?m)^"((?:[^"\\]|\\.)+)"\s*\{\s*return newToken\(SqlParserSymbols\.(\w+),`)
)

// classTokens are the terminals sql_scanner.flex makes from a class of
// texts, each with the DorisLexer.g4 token that matches the same texts, or a
// token of its own that tokens.tsv samples.
var classTokens = map[string]string{
	"IDENT":                     "IDENTIFIER",                // {IdentifierOrKw}, {QuotedIdentifier}
	"STRING_LITERAL":            "STRING_LITERAL",            // {SingleQuoteStringLiteral}, {DoubleQuoteStringLiteral}
	"INTEGER_LITERAL":           "INTEGER_VALUE",             // {IntegerLiteral} up to 2^63-1
	"LARGE_INTEGER_LITERAL":     "LARGE_INTEGER_LITERAL",     // {IntegerLiteral} above it
	"DECIMAL_LITERAL":           "DECIMAL_VALUE",             // {DoubleLiteral}, {ExponentLiteral}
	"COMMENTED_PLAN_HINT_START": "COMMENTED_PLAN_HINT_START", // {CommentedHintBegin}, {EolHintBegin}
	"COMMENTED_PLAN_HINT_END":   "COMMENTED_PLAN_HINT_END",   // {CommentedHintEnd}
	"UNMATCHED_STRING_LITERAL":  "UNMATCHED_STRING_LITERAL",  // a lone quote, which no statement holds
	"NUMERIC_OVERFLOW":          "NUMERIC_OVERFLOW",          // a number too large for its literal
}

// terminals maps each terminal of the CUP grammar to how the .g4 writes it.
func terminals(flex string) (map[string]string, error) {
	spelled := map[string][]string{}
	for _, m := range flexKeyword.FindAllStringSubmatch(flex, -1) {
		spelled[m[2]] = append(spelled[m[2]], m[1])
	}
	out := map[string]string{}
	for sym, texts := range spelled {
		// A keyword with more than one spelling (KW_AND is "and" and "&&")
		// is written as its name says.
		text := texts[0]
		if want := strings.ToLower(strings.TrimPrefix(sym, "KW_")); slices.Contains(texts, want) {
			text = want
		}
		out[sym] = literal(strings.ToUpper(text))
	}
	for _, m := range flexPunct.FindAllStringSubmatch(flex, -1) {
		if _, class := classTokens[m[2]]; class {
			continue // a lone quote, an error token
		}
		if _, ok := out[m[2]]; !ok {
			out[m[2]] = literal(strings.ReplaceAll(m[1], `\"`, `"`))
		}
	}
	for sym, tok := range classTokens {
		if _, ok := out[sym]; ok {
			return nil, fmt.Errorf("%s is both spelled and a class of texts", sym)
		}
		out[sym] = tok
	}
	return out, nil
}

// literal writes text as an ANTLR literal.
func literal(text string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(text) + "'"
}

// convert writes the .g4 for a CUP grammar and its scanner: the rules its
// start symbol reaches, start first and then in the file's order. An
// alternative that holds error, CUP's recovery token, matches no statement
// and is dropped.
func convert(header, cup, flex string) (string, error) {
	toks, err := lexCUP(cup)
	if err != nil {
		return "", err
	}
	rules, start, declared, err := readRules(toks)
	if err != nil {
		return "", err
	}
	terms, err := terminals(flex)
	if err != nil {
		return "", err
	}
	byName := map[string]*rule{}
	for i := range rules {
		if byName[rules[i].name] != nil {
			return "", fmt.Errorf("line %d: rule %s again", rules[i].line, rules[i].name)
		}
		byName[rules[i].name] = &rules[i]
	}
	if byName[start] == nil {
		return "", fmt.Errorf("no rule for the start symbol %s", start)
	}

	reached := map[string]bool{}
	var unscanned []string
	var reach func(name string) error
	reach = func(name string) error {
		if reached[name] {
			return nil
		}
		reached[name] = true
		for _, alt := range byName[name].alts {
			for _, sym := range alt {
				if byName[sym] != nil {
					if err := reach(sym); err != nil {
						return err
					}
				} else if _, ok := terms[sym]; !ok && sym != "error" {
					if !declared[sym] {
						return fmt.Errorf("rule %s: %s is neither a rule nor a terminal", name, sym)
					}
					// Declared, and never made by the scanner: a token of
					// its own, which tokens.tsv marks unreachable.
					terms[sym] = sym
					unscanned = append(unscanned, sym)
				}
			}
		}
		return nil
	}
	if err := reach(start); err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString(header)
	if len(unscanned) > 0 {
		slices.Sort(unscanned)
		fmt.Fprintf(&b, "// sql_parser.cup declares terminals sql_scanner.flex never makes: %s.\n\n", strings.Join(unscanned, ", "))
	}
	b.WriteString("parser grammar DorisLegacyParser;\n")
	write := func(r *rule) {
		fmt.Fprintf(&b, "\n// %s:%d\n%s\n", cupFile, r.line, r.name)
		sep := ":"
		for _, alt := range r.alts {
			if slices.Contains(alt, "error") {
				continue
			}
			parts := make([]string, len(alt))
			for k, sym := range alt {
				if byName[sym] != nil {
					parts[k] = sym
				} else {
					parts[k] = terms[sym]
				}
			}
			if len(parts) == 0 {
				fmt.Fprintf(&b, "    %s /* empty */\n", sep)
			} else {
				fmt.Fprintf(&b, "    %s %s\n", sep, strings.Join(parts, " "))
			}
			sep = "|"
		}
		b.WriteString("    ;\n")
	}
	write(byName[start])
	for i := range rules {
		if rules[i].name != start && reached[rules[i].name] {
			write(&rules[i])
		}
	}
	return b.String(), nil
}
