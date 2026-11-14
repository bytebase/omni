# doris test data

What `doris`'s committed engine answers are judged by and generated from. The
design is `docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`.

- **`pins.tsv`** names the Doris server by image digest and the source tag it
  was built from. `default` (3.1.4, tag `3.1.4-rc02`) judges
  `../parser/testdata/compat.tsv`.
- **`DorisParser.g4`** and **`DorisLexer.g4`** are Doris's ANTLR parser,
  `fe/fe-core/src/main/antlr4/org/apache/doris/nereids/` of
  [apache/doris](https://github.com/apache/doris) at tag `3.1.4-rc02`
  (commit `7f5ba43de6c7f148e1dacdb7cb838394d2ec1537`), unchanged. Its rows are
  `compat.tsv`'s first. `multiStatements` is the start rule; the other entry
  points (`singleStatement`, `expressionWithEof`) and `selectHint`, which the
  parser reads apart from a hint comment, get no row.
- **`sql_parser.cup`** and **`sql_scanner.flex`** are Doris's legacy parser,
  `fe/fe-core/src/main/cup/` and `fe/fe-core/src/main/jflex/` at the same tag,
  unchanged. The ANTLR parser reads every statement first, and hands the ones
  its `unsupportedStatement` rule derives (most `SHOW`, `SET`, `ADMIN`,
  `ALTER SYSTEM` and account statements), and a `LOAD LABEL` or
  `CREATE ROW POLICY`, to the CUP parser, which reads the text again: Doris
  takes such a statement only if both do.
- **`DorisLegacyParser.g4`** is the CUP grammar converted to ANTLR by
  `../parser/cmd/cup2g4`, each rule under the line of `sql_parser.cup` it came
  from. Its rows follow `DorisParser.g4`'s. It is the whole grammar the CUP
  parser's start symbol reaches, not only the statements handed to it: the
  ANTLR parser decides which those are, and the engine judges every row.
  `TestDorisLegacyParserIsTheConversion` holds it to the conversion.
- **`tokens.tsv`** gives the text each token class of `DorisLexer.g4` is written
  as, citing its lexer rule, and the tokens of `sql_scanner.flex` that
  `DorisLexer.g4` has no match for.

| File | SHA-256 |
| --- | --- |
| `DorisParser.g4` | `b1917d98bf14a66647ab43eaa8ef988a283c2e633efde00e7cc3ec875c11d2ed` |
| `DorisLexer.g4` | `e27fd7351090c7c1735d6d131c8b16f86d11b8c269408fac8fdc9ef2c1e1a2a4` |
| `sql_parser.cup` | `70a73096ae8e0d6ed87c16a069ea5ad1d446a8ccb607ed54dbf90b68caa18a5e` |
| `sql_scanner.flex` | `cea14f6c0e4a7f4b598041c54d1834787a9ecf6bf5ffffe14823bba433908e7c` |

A pin bump replaces the four files with the new tag's and reruns the
converter. They are distributed under the Apache License 2.0, in
`Doris.LICENSE` (the repository's `LICENSE.txt` at the tag). They are test
fixtures read at test time and are not compiled into omni.
