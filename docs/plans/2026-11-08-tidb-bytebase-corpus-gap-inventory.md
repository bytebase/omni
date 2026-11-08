# TiDB real-SQL corpus gap inventory (§3.4 sample 3)

Date: 2026-11-08. The first measurement of §3.4 sample 3, "Bytebase's own TiDB
test files" (`tidb/parser/testdata/bytebase_tidb_corpus.tsv`). This is the SQL
Bytebase already promises to handle, so an over-acceptance here is a form
Bytebase ships against. Sibling to the sample 1 inventory
(`2026-11-08-tidb-real-sql-corpus-gap-inventory.md`).

## The extraction

The 64 Bytebase TiDB test files, at bytebase commit `5a6096dfd5`, under three
directories (paths under `backend/plugin`):

- `parser/tidb/test-data` — 7 files (query_span, query_type,
  test_statement_type, statement-ranges/standard, test_split, test_backup,
  test_restore);
- `advisor/tidb/test` — 53 files;
- `schema/tidb/testdata` — walk_through, generate_migration,
  get_database_metadata (get_database_definition holds no input SQL).

A throwaway Go extractor read the **input-side** SQL — the values at the
top-level case keys `statement`, `input`, `oldSchema`, `newSchema`, `schema` —
split each script with `tidb/parser.Split`, trimmed, and kept the non-empty
segments. Top-level keys only, so a query-span source column's nested `schema`
and a backup case's generated `result` statements are not picked up.
Generated-output SQL (backup/restore rollback, `generate_migration` `expected`,
`get_database_definition`) is Bytebase-synthesized, not promised input, and is
excluded.

From 600 scripts / 945 split segments, **527 rows** remain after dropping:

- 194 multi-line statements (embedded newline/tab — the format is one statement
  per line, as the pingcap corpora are; most are pretty-printed multi-line
  `CREATE TABLE` DDL from the schema fixtures);
- 209 internal duplicates;
- 15 already carried by `pingcap_parser_corpus.tsv` /
  `pingcap_parser_rest_corpus.tsv`.

Each row is measured against omni and TiDB v8.5.5 via the same `PREPARE`
harness the other corpora use.

## The measurement

| | first |
|---|---:|
| both accept | 526 |
| gap: TiDB accepts, omni rejects | 0 |
| **over-acceptance: TiDB rejects, omni accepts** | **0** |
| both reject | 1 |
| no verdict | 0 |

omni parses **526 of the 526 statements TiDB accepts — 100.0%**, with **zero
over-acceptances**. There is no grammar work list from this sample: it is a
clean sheet against the forms Bytebase ships against.

## What the rows are, by family

The 526 accepted rows, grouped by leading form:

| family | rows |
|---|---:|
| `ALTER TABLE` | 194 |
| `CREATE TABLE` | 142 |
| `UPDATE` | 49 |
| `SELECT` | 45 |
| `CREATE INDEX` / `CREATE UNIQUE`/`FULLTEXT INDEX` | 30 |
| `INSERT` / `REPLACE` | 19 |
| `DROP` | 12 |
| `DELETE` | 11 |
| `EXPLAIN` | 6 |
| `WITH` (CTE) | 5 |
| `RENAME` | 5 |
| `ALTER DATABASE` | 4 |
| `SET` | 3 |
| `CREATE VIEW`, conditional `/*! */` comment | 2 |

The weight is DDL — column rules, index limits, naming conventions and schema
syncs are what Bytebase's advisor and schema code exercise — which is where a
hand-written parser is most prone to over-accept. None did.

## The one both-reject

```
SELECT FROM
```

A `SELECT` with no select list before `FROM`, from a query-span error case.
Both omni and the server reject it; it is not work.

## Next

Nothing to close. This sample is clean on its first measurement, so it adds no
rows to §3.2's work list. If the bytebase pin or these fixtures change, re-run
the extraction at the new commit and re-record with
`TestTiDBBytebaseCorpusEngineVerdicts -update`; the diff is the list of forms
Bytebase added or removed. Sample 2 (TiDB integrationtest) is a separate session
and will extend the sample 1 work list.
