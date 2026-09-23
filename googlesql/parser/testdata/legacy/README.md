# Legacy ANTLR4 GoogleSQL Parser Examples

These SQL files are mirrored from `bytebase/parser/googlesql/examples/`. They form
**Corpus A** of the omni googlesql parser test suite: the legacy-corpus closure
gate (`TestLegacyCorpusTokenizes`, `TestLegacyCorpusParses`,
`TestDML_LegacyCorpusAccepts`, `TestSelect_TPCHCorpus` in `googlesql/parser`).
The layout is preserved verbatim because the skip-list in
`legacy_corpus_test.go` is keyed by path relative to this directory.

| Path | Contents |
|---|---|
| `easy.sql` | hand-written smoke file |
| `zetasql/parser/testdata/*.sql` | statements lifted from the ZetaSQL parser testdata (49 files) |
| `zetasql/examples/tpch/*.sql` | the 22 TPC-H reference queries in GoogleSQL |

When updating: re-copy from `bytebase/parser/googlesql/examples/` upstream and
re-baseline `legacyCorpusParseSkips`. Do not edit these files in place.

Source commit: `bytebase/parser@57b6ef7a2640481d8734cd63af0c7b781fa85f22` (2026-04-17).

## Licensing

- `bytebase/parser` is distributed under the BSD-3-Clause license
  (Copyright (c) 2025, Bytebase).
- The `zetasql/` subtree is derived from Google's ZetaSQL project
  (<https://github.com/google/zetasql>, `zetasql/parser/testdata/` and
  `zetasql/examples/tpch/`), distributed under the Apache License 2.0
  (Copyright 2019 Google LLC).

Both licenses permit redistribution with attribution; this file is that
attribution. The files are test fixtures only and are not compiled into omni.
