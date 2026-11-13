# Redshift parser test data

## `compat.tsv`: the docs' rows

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/RedshiftDocs.g4`, then its near-misses: the statement with one
token dropped, the end cut off, a stray word inserted, or a keyword swapped.
Each row has omni's verdict and the judge's. `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes. It starts no server, regenerates the rows from the
grammar, and asks the judge again on every run.

| Rows | Total | Agree | Only the docs accept | Only omni accepts | Undocumented |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 2,241 | 1,935 | 306 | 0 | 0 |
| Near-misses | 5,839 | 105 | 40 | 0 | 5,694 |

- **The grammar** is AWS's SQL command reference
  (https://docs.aws.amazon.com/redshift/latest/dg/c_SQL_commands.html): the
  Syntax sections of the 125 pages its Topics list, converted by
  `../cmd/docsgrammar` with `enginetest/grammar/docs`, the converter
  Snowflake's grammar comes from. Each rule sits under its page's URL.
  `../../testdata/pins.tsv` pins the date the pages were read and a SHA-256 of
  the blocks.
- **AWS writes the notation its own way.** A value is in italics or in lower
  case, and a fragment is defined by a "where *name* is:" line. The
  converter's front end rewrites both into the shared notation (`<name>`,
  `name ::=`) before converting.
- **Hand fixes** are in `../../testdata/docs-fixes.tsv`, 73 of them, each with
  its reason: brackets left open or closed twice, choices not braced, prose and
  examples inside a Syntax section. The converter fails on a fix whose text the
  page no longer holds.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the names,
  literals and reserved-words pages.

To regenerate the grammar from the live docs (a pin bump), then the rows:

```bash
go run ./redshift/parser/cmd/docsgrammar -cache /tmp/redshift-docs
OMNI_RECORD=1 go test ./redshift/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./redshift/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

Redshift has no server CI can start, so every row is judged by the docs.

- **A statement the grammar derives is accepted.** A near-miss it does not
  derive is `undocumented`, never `reject`: docs say what is valid, not what
  Redshift refuses. It passes whatever omni says, so omni is never pushed to
  reject a form users write.
- **346 rows disagree, all statements the docs accept and omni rejects.**
  `GRANT` and `REVOKE` hold 102 (`ASSUMEROLE`, `DEBUG`, templates, datashares),
  `SHOW` 81 (`RLS POLICIES`, `COLUMN GRANTS`, `CONSTRAINTS`, `PARAMETERS`,
  `TEMPLATES`), templates' own statements 50, `CREATE` and `ALTER DATABASE` 41
  (`FROM INTEGRATION`, `FROM ARN`), `COPY` 18, `CREATE MATERIALIZED VIEW` 15,
  `SELECT` 12 (`EXCLUDE`, `CONNECT BY`), `UNLOAD` 9.

## What the judge cannot decide

- **What the docs leave out.** The command pages give no expression or query
  grammar: an expression is sampled as `1`, a query as `SELECT 1`, and the
  pages a command lists under Topics, such as SELECT's clauses, are not read.
- **What the docs get wrong.** The grammar is the docs as written, less the
  hand fixes: `ALTER USER`'s options with commas between them, an `UNLOAD`
  option list with commas, `REVOKE [ GRANT OPTION ]` without `FOR`. A row the
  engine would refuse is accepted until the engine judges it.
- **The engine.** With `REDSHIFT_COMPAT_DSN` set to a dedicated, scratch-only
  cluster, a recording could ask Redshift instead, through the raw-parse probe
  `compat.RawParseVerdict` reads for pg. That judge is not built: it would add
  a pin row for the cluster and give the suite `compat.SourceEngine` and a
  `Judge`.
