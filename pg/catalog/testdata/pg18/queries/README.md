# PostgreSQL 18 query simulation

`simulate.tsv` holds `../../queries/simulate.tsv`'s rows (every generated
`SELECT` PostgreSQL accepts in `pg/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on `one_table@column`), recorded against the pinned
PostgreSQL 18 (`../../../../testdata/pins.tsv`, role `pg18`:
`postgres:18.7-alpine`). `TestQuerySimulationPG18`
(`../../../simulate_query_test.go`) replays it from
`../../tapes/TestQuerySimulationPG18.tape` with the catalog at `PG18`.

**1,151 rows: 1,151 agree, 0 disagree.** PostgreSQL 18 answers every row as 17 does.

## Commands

From `<repo root>`:

```bash
go test ./pg/catalog/ -run '^TestQuerySimulationPG18$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./pg/catalog/ -run '^TestQuerySimulationPG18$'
```
