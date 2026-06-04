package parser

import "strings"

// This file ships the GoogleSQL identifier foundation: string-level
// normalization (the byte-for-byte port of the legacy bytebase
// unquoteIdentifierByText helper that the bigquery/spanner query-span
// extractors use to canonicalize name parts) plus the token-level identifier
// predicates the statement dispatcher and later grammar nodes (expressions,
// types, SELECT, DDL, DML) consume.
//
// It deliberately does NOT build AST identifier nodes: the GoogleSQL AST
// (googlesql/ast) ships only the File root container in the foundation; the
// concrete Identifier / path-expression node types and their token→node parse
// functions are owned by the expressions node. The predicates here are the
// shared, node-independent half (which tokens may begin / continue an
// identifier), kept in the foundation so every downstream grammar node agrees
// on the rule.
//
// GoogleSQL identifier rules (GoogleSQLLexer.g4 + GoogleSQLParser.g4):
//
//   - UNQUOTED_IDENTIFIER: [A-Za-z_][A-Za-z0-9_]*  (case-insensitive,
//     normalized to lower case by Spanner/BigQuery name resolution — but the
//     parser preserves the source spelling; folding is a resolution concern).
//   - `backtick`-quoted IDENTIFIER: any text; the lexer strips the backticks
//     and emits tokIdentifier with the body in Token.Str.
//   - A NON-reserved word-keyword may stand in as an identifier
//     (GoogleSQLParser.g4 common_keyword_as_identifier / keyword_as_identifier).
//     The reserved/non-reserved split lives in keywords.go (IsReservedKeyword).

// isIdentifierStart reports whether tokenType may begin an identifier in a
// position where keyword-as-identifier is allowed (the common case in
// GoogleSQL paths and aliases):
//   - tokIdentifier: an UNQUOTED_IDENTIFIER or a `backtick`-quoted identifier
//     (the lexer collapses both to tokIdentifier).
//   - a non-reserved word-keyword (tokenType >= keywordBase and not reserved).
//
// Reserved keywords are excluded: GoogleSQL rejects a reserved word as a bare
// identifier (it must be backtick-quoted, which the lexer would have tokenized
// as tokIdentifier). The lexer never enforces the reserved split — it always
// emits the kw* token for a keyword — so this predicate is where the parser
// applies it.
func isIdentifierStart(tokenType int) bool {
	if tokenType == tokIdentifier {
		return true
	}
	return tokenType >= keywordBase && !keywordReserved[tokenType]
}

// isAnyKeywordIdentifier reports whether tokenType may stand in as an
// identifier when EVERY word-keyword — reserved or not — is permitted. A few
// GoogleSQL contexts (notably dotted path components after the first, and
// generic-entity object kinds) accept reserved keywords as name parts. Callers
// that must follow the strict identifier rule use isIdentifierStart instead;
// this predicate is the permissive variant for those name-part positions.
func isAnyKeywordIdentifier(tokenType int) bool {
	return tokenType == tokIdentifier || tokenType >= keywordBase
}

// NormalizeGoogleSQLIdentifier returns the canonical form of a single
// GoogleSQL identifier given its raw source spelling. A `backtick`-quoted
// identifier (at least two backticks with content between them) has its
// surrounding backticks stripped; any other spelling is returned unchanged.
//
// This is a byte-for-byte port of the legacy bytebase unquoteIdentifierByText
// helper (backend/plugin/parser/{bigquery,spanner}/query_span_extractor.go) so
// the cutover keeps name comparison (field lineage, data-access-control name
// resolution, masking) identical to the ANTLR path. In particular it
// intentionally matches legacy by:
//   - stripping ONLY surrounding backticks (the legacy `len >= 3 && hasPrefix
//     "`" && hasSuffix "`"` guard, which leaves a lone "`" or "“" untouched),
//   - NOT lower-casing (GoogleSQL unquoted identifiers ARE case-insensitive for
//     resolution, but the legacy extractor preserved the source case here and
//     left case-folding to the metadata layer; preserving that avoids a silent
//     behavior change), and
//   - NOT collapsing any backslash escape inside the backticks (the legacy
//     helper sliced the inter-backtick bytes verbatim).
func NormalizeGoogleSQLIdentifier(raw string) string {
	if len(raw) >= 3 && strings.HasPrefix(raw, "`") && strings.HasSuffix(raw, "`") {
		return raw[1 : len(raw)-1]
	}
	return raw
}
