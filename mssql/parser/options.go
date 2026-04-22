// Package parser - options.go implements option name validation for T-SQL WITH clauses.
//
// Many T-SQL statements accept WITH (<option_name> = <value>) clauses where
// <option_name> must be one of a finite set of valid keywords. The parser was
// previously too permissive, accepting any keyword/identifier as an option name.
//
// This file provides a lightweight framework for declaring valid option sets
// and checking/consuming option tokens against them.
package parser

import "strings"

// optionSet declares a set of valid option names for a particular parsing position.
// Keys are keyword token types (kw* constants). For option names that are not
// registered keywords, use the special key 0 (tokEOF) mapped to a set of
// uppercase identifier strings via identOptions.
type optionSet struct {
	// tokens contains the set of keyword token types that are valid options.
	tokens map[int]bool
	// idents contains uppercase identifier strings that are valid options
	// but are not registered as keyword tokens.
	idents map[string]bool
}

// newOptionSet creates an optionSet from a list of keyword token types.
func newOptionSet(keywords ...int) optionSet {
	s := optionSet{tokens: make(map[int]bool, len(keywords))}
	for _, kw := range keywords {
		s.tokens[kw] = true
	}
	return s
}

// withIdents returns a copy of the optionSet with additional valid identifier strings.
// Use this for option names that are not registered as keyword tokens.
func (s optionSet) withIdents(names ...string) optionSet {
	out := optionSet{
		tokens: s.tokens,
		idents: make(map[string]bool, len(names)),
	}
	for _, n := range names {
		out.idents[strings.ToUpper(n)] = true
	}
	return out
}

// isValidOption returns true if the current token is a valid option name
// according to the given optionSet.
//
// It checks:
//  1. Whether the token type is a keyword in the set (handles Core and Context keywords).
//  2. Whether the token is an identifier whose uppercase text matches a registered
//     keyword in the set (for context keywords scanned as identifiers).
//  3. Whether the token is an identifier whose uppercase text is in the idents set
//     (for option names not registered as keywords).
func (p *Parser) isValidOption(opts optionSet) bool {
	// Direct keyword match.
	if opts.tokens[p.cur.Type] {
		return true
	}
	// Identifier (including bracketed) whose text matches.
	if p.cur.Type == tokIDENT {
		kw := lookupKeyword(p.cur.Str)
		if kw != tokIDENT && opts.tokens[kw] {
			return true
		}
		if len(opts.idents) > 0 && opts.idents[strings.ToUpper(p.cur.Str)] {
			return true
		}
		return false
	}
	// Registered keyword whose uppercase name matches a declared ident. This
	// lets `.withIdents("FULL", "SIMPLE", ...)` accept both the registered
	// kwFULL token and the unregistered "SIMPLE" identifier in a single enum
	// declaration, without requiring callers to track which values are
	// lexer-registered.
	if len(opts.idents) > 0 && p.cur.Str != "" && opts.idents[strings.ToUpper(p.cur.Str)] {
		return true
	}
	return false
}

// isValidOptionToken checks whether the given token (not necessarily the current
// token) is a valid option name according to the given optionSet. This is useful
// for peeking ahead without advancing the parser.
func (p *Parser) isValidOptionToken(opts optionSet, tok Token) bool {
	if opts.tokens[tok.Type] {
		return true
	}
	if tok.Type == tokIDENT {
		kw := lookupKeyword(tok.Str)
		if kw != tokIDENT && opts.tokens[kw] {
			return true
		}
		if len(opts.idents) > 0 && opts.idents[strings.ToUpper(tok.Str)] {
			return true
		}
	}
	return false
}

// expectOption consumes the current token if it is a valid option name
// according to the given optionSet, returning the uppercase option name.
// Returns a syntax error if the current token is not a valid option.
func (p *Parser) expectOption(opts optionSet) (string, error) {
	if !p.isValidOption(opts) {
		return "", p.syntaxErrorAtCur()
	}
	name := strings.ToUpper(p.cur.Str)
	p.advance()
	return name, nil
}
