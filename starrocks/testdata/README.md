# starrocks test data

What `starrocks`'s committed engine answers are judged by and generated from.
The design is `docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`.

- **`pins.tsv`** names the StarRocks server by image digest and the release
  it was built from. `default` (3.4.10) judges `../parser/testdata/compat.tsv`.
- **`StarRocks.g4`** and **`StarRocksLex.g4`** are
  `fe/fe-core/src/main/java/com/starrocks/sql/parser/` of
  [StarRocks/starrocks](https://github.com/StarRocks/starrocks) at tag
  `3.4.10` (commit `965be759154f297a1b9d66bef0ae5ba07da603a9`), unchanged:
  SHA-256 `6b15ef6d603292a8daae803ee8e61277f87a514ada2b5463f43ad5cc23a56745`
  and `adb275d4a4727d09704e3e79348f2b5855e38411c2b672a3517aff1ceac1ac2e`. The
  generated rows of `compat.tsv` are their alternatives and near-misses.
  `sqlStatements` is the start rule; `expressionSingleton`, another entry
  point, and `hintMap`, which the parser reads apart from a hint comment, get
  no row. A pin bump replaces them with the new tag's.
- **`tokens.tsv`** gives the text each token class of `StarRocksLex.g4` is
  written as, citing its lexer rule, and the statement written for
  `translateSQL`'s wildcard.

The two grammars are distributed under the Apache License 2.0, in
`StarRocks.LICENSE` (the repository's `LICENSE.txt` at the tag). They are test
fixtures read at test time and are not compiled into omni.
