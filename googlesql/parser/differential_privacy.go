package parser

// This file is part of the `parser-query-clauses` DAG node. It implements
// GoogleSQL's SELECT-level differential-privacy / aggregation-threshold clause
// (GoogleSQLParser.g4 §2.13 opt_select_with), a hand-port of Google's
// open-source ZetaSQL reference.
//
// The clause is the `WITH <name> [OPTIONS(...)]` that may follow SELECT before
// the set quantifier, e.g.
//
//	SELECT WITH DIFFERENTIAL_PRIVACY OPTIONS(epsilon=1, ...) col FROM t
//	SELECT WITH ANONYMIZATION OPTIONS(...) col FROM t
//	SELECT WITH AGGREGATION_THRESHOLD OPTIONS(...) col FROM t
//
// The grammar reads the privacy mechanism as a bare `identifier`
// (DIFFERENTIAL_PRIVACY / ANONYMIZATION / AGGREGATION_THRESHOLD are accepted as
// names, not dedicated keywords), so any name is accepted. The clause name is
// the load-bearing fact for the downstream consumers (query-span / diagnose);
// the OPTIONS body is validated semantically, not syntactically, so it is parsed
// and skipped (its key/value vocabulary is open in the grammar). The presence +
// name is recorded on SelectStmt.SelectWith.
//
// DIALECT NOTE. This is a BigQuery-only form. The Spanner emulator does not
// support `SELECT WITH ...` and reports "Unexpected keyword WITH" WITHOUT the
// canonical `Syntax error:` prefix, so the differential harness misclassifies it
// as a (semantic) ACCEPT — i.e. the oracle is NON-authoritative here (oracle.md
// dialect caveat). It is therefore proven by the hand-written unit tests +
// triangulation against truth1/bigquery + the legacy .g4, NOT by the differential
// (mirroring the parser-select node's exclusion of differential-privacy fixtures
// from its oracle test).

// atSelectWith reports whether the current position begins an opt_select_with
// clause: `WITH <identifier>` (NOT `WITH (`, which would be an inline WITH-expr
// column). SELECT's hint has already been consumed by the caller.
func (p *Parser) atSelectWith() bool {
	return p.cur.Type == kwWITH && p.peekNext().Type != int('(')
}

// parseSelectWith parses the opt_select_with clause and returns its mechanism
// name (the identifier after WITH). The current token is WITH. The optional
// trailing `OPTIONS(...)` is parsed and skipped (its body is semantic, not
// grammatical). The caller (parseSelectStmt) must have already checked
// atSelectWith.
//
// The mechanism name is the legacy grammar's `identifier`. The two name-shaped
// mechanisms ANONYMIZATION / AGGREGATION_THRESHOLD lex as plain identifiers, but
// DIFFERENTIAL_PRIVACY is lexed by the omni lexer as a RESERVED keyword
// (kwDIFFERENTIAL_PRIVACY) — an over-reservation versus the legacy grammar,
// which has no such token and treats the word as an ordinary identifier (the
// emulator likewise accepts `differential_privacy` as a bare identifier). So
// this name slot explicitly admits that one over-reserved keyword in addition to
// the normal identifier set, to keep the documented
// `SELECT WITH DIFFERENTIAL_PRIVACY OPTIONS(...)` form parseable. (The broader
// lexer over-reservation — `differential_privacy` cannot be used as a column /
// alias name — is a separate lexer-node concern outside this node's files.)
func (p *Parser) parseSelectWith() (string, error) {
	p.advance() // WITH
	if !isIdentifierStart(p.cur.Type) && p.cur.Type != kwDIFFERENTIAL_PRIVACY {
		return "", p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return "", err
	}
	if p.cur.Type == kwOPTIONS {
		p.advance() // OPTIONS
		if err := p.skipOptionsList(); err != nil {
			return "", err
		}
	}
	return name, nil
}
