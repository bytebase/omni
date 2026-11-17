# Legacy Snowflake ANTLR4 Grammar Files

These two files are mirrored verbatim from `bytebase/parser/snowflake/`. They
are the source of truth for the keyword gates in `snowflake/parser`
(`TestKeywordCompleteness`, `TestReservedKeywordsCompleteness` in
`keyword_completeness_test.go`):

| File | Used for |
|---|---|
| `SnowflakeLexer.g4` | the word-keyword set: every `NAME: 'WORD';` rule must have a `keywordMap` entry |
| `build_id_contains_non_reserved_keywords.py` | the reserved/non-reserved split: its `snowflake_reserved_keyword` dict is the reserved keyword set, minus the phantoms the lexer never emits as a token (listed in the test) |

`keywords.go` was generated from these files, and the gates keep them from
drifting apart.

When updating: re-copy both files from `bytebase/parser/snowflake/` upstream,
bump the source commit below, and regenerate `keywords.go` if the gates fail.
Do not edit these files in place.

Source commit: `bytebase/parser@36e79b3acb7e2161beecdced086a5b3eaae06605` (2026-02-27).

## Licensing

`SnowflakeLexer.g4` is MIT-licensed (Copyright (c) 2022, Michał Lorek); its
license text is in the file's header. `bytebase/parser` is distributed under the
BSD-3-Clause license (Copyright (c) 2025, Bytebase). Both licenses permit
redistribution with attribution; this file is that attribution. The files are
test fixtures only and are not compiled into omni.
