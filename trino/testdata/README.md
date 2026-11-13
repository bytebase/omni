# trino test data

What `trino`'s committed engine answers are judged by and generated from. The
design is `docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`.

- **`pins.tsv`** names the Trino server by image digest and the release it was
  built from. `default` (482) judges `../parser/testdata/compat.tsv`.
- **`SqlBase.g4`** is `core/trino-grammar/src/main/antlr4/io/trino/grammar/sql/SqlBase.g4`
  of [trinodb/trino](https://github.com/trinodb/trino) at tag `482` (commit
  `f04d222fbeedaf888ac3c907748209c7e716a4c2`), unchanged: SHA-256
  `77d09cc5afdb14d900af342d8ca780cc4917f27309230e177c29f365a2380cdb`. The
  generated rows of `compat.tsv` are its alternatives and their near-misses.
  The five `standalone*` entry points are not statements and get no row. A
  pin bump replaces it with the new tag's.
- **`tokens.tsv`** gives the text each token class of `SqlBase.g4` is written
  as, citing its lexer rule, and the keyword written for `methodName`'s
  wildcard.

`SqlBase.g4` is distributed under the Apache License 2.0, in `Trino.LICENSE`.
It is a test fixture read at test time and is not compiled into omni.
