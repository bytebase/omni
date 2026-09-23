# Legacy ANTLR4 GoogleSQL Grammar

These two grammar files are mirrored verbatim from `bytebase/parser/googlesql/`.
They are the source of truth for the keyword gates in `googlesql/parser`
(`TestKeywordCompleteness`, `TestReservedKeywordSplit` in
`keyword_completeness_test.go`):

| File | Used for |
|---|---|
| `GoogleSQLLexer.g4` | the full word-keyword set: every `*_SYMBOL: 'WORD';` rule must have a `keywordMap` entry, and nothing else may |
| `GoogleSQLParser.g4` | the reserved/non-reserved split: a word-keyword is non-reserved iff it appears in `common_keyword_as_identifier` (plus `SIMPLE` via `keyword_as_identifier`) |

The lexer grammar is a hand-port of ZetaSQL's `flex_tokenizer.l`; `tokens.go`
and `keywords.go` were generated from these files and the gates keep them from
drifting apart.

When updating: re-copy both files from `bytebase/parser/googlesql/` upstream,
bump the source commit below, and regenerate `tokens.go`/`keywords.go` if the
gates fail. Do not edit these files in place.

Source commit: `bytebase/parser@57b6ef7a2640481d8734cd63af0c7b781fa85f22` (2026-04-17).

## Licensing

`bytebase/parser` is distributed under the BSD-3-Clause license
(Copyright (c) 2025, Bytebase). That license permits redistribution with
attribution; this file is that attribution. The grammars are test fixtures only
and are not compiled into omni.
