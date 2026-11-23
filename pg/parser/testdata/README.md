# pg parser test data

`compat.tsv` holds a statement per row, with omni's verdict and the pinned
PostgreSQL's (`../../testdata/pins.tsv`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes, and starts no server.

## Rows

| Rule | Rows | From |
|---|---|---|
| Generated (`Rule/symbols`, near-misses `!drop-2`, `!swap`) | 13,971 | `../../testdata/gram.y` and `../../testdata/tokens.tsv` |
| `suite:lexer` | 250 | the token forms `tokens.tsv` does not spell (below) |
| `suite:psql` | 7 | psql-adjacent non-metacommand and variable forms rejected by both server and omni |
| `suite:paren-fuzz` | 2,827 | the paren dispatch fuzz (`../paren_oracle_fuzz_test.go`) |

**Generated rows.** For each alternative of each rule in `gram.y`, the
shortest statement that reaches it, then its near-misses: the statement with
one token dropped, the end cut off, a stray word inserted, or a keyword
swapped. 13 alternatives get no row: the five `MODE_*` entry points
`raw_parser` uses for a type name and for PL/pgSQL, and the rules only they
reach. A rule `tokens.tsv` marks `(every site)` (`opt_drop_behavior`,
`opt_or_replace`, `OptTemp`, `opt_if_not_exists` and the other optional
clauses listed there) also gets a row at each alternative that names it, with
the clause written out (`Rule/symbols+opt_drop_behavior`), whose near-misses
drop the clause's tokens. `TestCompatibility` regenerates these rows and fails
if the file holds others.

**`suite:lexer`** is every form `scan.l` reads that a generated row does not
write, each in a `SELECT`: strings with `''` and with a newline between
literals, `E''` and its escapes (octal, hex, `\u`, `\U`, surrogate pairs, an
escape that makes a zero byte or an invalid character), `U&''` and `U&""`
with and without `UESCAPE`, `B''` and `X''`, dollar quoting (tags, nesting,
case, a tag that is not one), quoted and unquoted identifiers, numbers
(decimal, `0x`, `0o`, `0b`, `_` separators, exponents, trailing junk),
parameters, operators (the `+-` and comment rules of `{operator}`, the
two-character tokens, `OPERATOR()`), nested and unterminated comments, and the
whitespace characters.

**`suite:psql`** holds 7 psql-adjacent forms (`SELECT 1 \;`, `\`, `\1`, `\\`,
`SELECT :a`, `SELECT :'a'`, `SELECT :"a"`) that neither PostgreSQL nor omni
accepts (both reject at the same byte). The 8 client-side `psql` backslash
metacommands (`\d`, `\dt`, `\copy`, `\restrict`, `\echo`, `\gset`, `\g`,
`\gexec`) are intentional client-mode extensions that `Parse` reads as comments
(`../lexer.go`, `../../internal/metacmd`) so scripts written for `psql` and
`pg_dump` parse cleanly; those 8 forms are pinned directly in
`../metacmd_test.go` rather than `compat.tsv`.

**`suite:paren-fuzz`** is the statements `TestParenFuzzRows` generates: 10,000
draws from a fixed seed over the template families, 2,827 distinct. The test
fails when the rows are not what it generates, and checks the one thing a
verdict does not say: where PostgreSQL accepts a statement whose first FROM
item opens with a paren, omni routes it to a subquery or a joined table. It
replaces the paren oracle's own corpus (`paren-fuzz-corpus/`,
`paren-fuzz-defer/`) and the `PAREN_FUZZ_SIZE` it read.

## Where the verdict comes from

The judge sends `SELECT 1/0;` and the statement on the next line, in one
simple query (`compatJudge`, read by `compat.RawParseVerdict`). PostgreSQL
parses the whole text before it runs any of it, so:

- **Division by zero** means the statement parsed: `accept`. Nothing of the
  statement ran, so the judge needs no transaction and leaves nothing behind.
- **Any other error** is the raw parser's, whatever its SQLSTATE: 42601 from
  `scan.l` and `gram.y`, and what a grammar action chooses (0A000 for
  `CREATE ASSERTION`, which `gram.y` refuses). The verdict is `reject@<byte>`,
  the byte the server points at, and `omniVerdict` writes `ParseError.Position`
  the same way. So a row agrees only when omni rejects at the same byte.
- **No position:** 34 rows, where an action reports none: an unknown
  partitioning strategy, a constraint attribute on a domain's `NOT NULL`
  (`processCASbits`), `WITH TIES` without `ORDER BY`. The verdict is a bare
  `reject`.

What this cannot tell:

- **Analysis is not parsing.** A statement `gram.y` takes and parse analysis
  refuses (`SELECT a FROM` a missing table) is `accept`, as it should be.
- **A zero byte** cannot be sent: the protocol carries a query as a C string.
  No row holds one.

## Counts at `REL_17_12`

PostgreSQL 17.12 server compatibility is **17,055 / 17,055 (100.00%)** (`13,971`
PG 17 grammar rows + `250` `suite:lexer` + `7` `suite:psql` server-matching
forms + `2,827` `suite:paren-fuzz`).

| Rows | Total | Same verdict | Same byte | Only PostgreSQL accepts | Only omni accepts |
|---|---|---|---|---|---|
| Alternatives | 3,591 | 3,591 | 3,591 | 0 | 0 |
| Near-misses | 10,380 | 10,380 | 10,380 | 0 | 0 |
| `suite:lexer` | 250 | 250 | 250 | 0 | 0 |
| `suite:psql` | 7 | 7 | 7 | 0 | 0 |
| `suite:paren-fuzz` | 2,827 | 2,827 | 2,827 | 0 | 0 |
| **All** | **17,055** | **17,055** | **17,055** | **0** | **0** |

PostgreSQL accepts 3,499 of the 3,591 alternatives. A grammar action refuses
the other 92 as the generator writes them: `CREATE ASSERTION`, an operator
written with one argument type (`DROP OPERATOR ~(a)`), a publication object
with no `TABLE` before it, a `JSON_TABLE` path with no name, a window frame
that ends before it starts, `OVERLAPS` with rows of the wrong length,
`CREATE OR REPLACE CONSTRAINT TRIGGER`.

## Commands

From `<repo root>`. Judge new rows, or every row after a pin bump:

```bash
OMNI_RECORD=1 go test ./pg/parser/ -run '^TestCompatibility$'
```

It starts the pinned image and judges all 17,055 rows in under ten seconds.
Commit rows a parser fix moved:

```bash
go test ./pg/parser/ -run '^TestCompatibility$' -update-compat
```

After a change to the paren fuzz's generator, rewrite its rows, then judge
the new ones:

```bash
go test ./pg/parser/ -run '^TestParenFuzzRows$' -update-compat
```

## `tapes/`: the oracle tests' recordings

Every other test here that asks PostgreSQL something reaches it through
`testcontainer.Open` (`pg/internal/testcontainer`) and replays
`tapes/<Test>.tape`: each request the test sends, and the pinned server's
answer. `enginetest/tape` describes the format and what fails a replay. After
changing what a test sends, record it:

```bash
OMNI_RECORD=1 go test ./pg/parser/ -run '^TestName$'
```

The whole package records in about 15 seconds: 29 tapes, 26,923 requests.
