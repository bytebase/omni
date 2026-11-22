# MSSQL review test data

`tapes/` holds recorded SQL Server 2022 (`../testdata/pins.tsv`,
`mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04`) exchanges for the
oracle and grammar-derived simulation tests in `mssql/review/` so CI replays
them offline with zero Docker:

- `TestPriorBackupAgainstMSSQL.tape` — prior-backup oracle (`oracle_test.go`)
- `TestScanRulesAgainstMSSQL.tape` — constraint/PK scan rules oracle (`oracle_test.go`)
- `TestWalkThroughAgainstMSSQL.tape` — curated `WALK_THROUGH` oracle (`walk_through_oracle_test.go`)
- `TestWalkThroughBatchesAgainstMSSQL.tape` — batch-boundary `WALK_THROUGH` oracle (`walk_through_oracle_test.go`)
- `TestWalkThroughSweepMSSQL.tape` — randomized collation/collision sweep (`walk_through_oracle_test.go`)
- `TestWalkSimulation.tape` — grammar-derived `WALK_THROUGH` simulation (`walk_simulation_test.go`)

`fixtures/` holds the starting schemas for `TestWalkSimulation`, and
`walk_simulation_verdicts.tsv` records SQL Server 2022's verdict (`ok` or error
number) for every `(fixture, statement)` pair derived from
`mssql/parser/testdata/compat.tsv`.

## Commands

From `<repo root>`. Replay everything offline (no Docker):

```bash
go test ./mssql/review/...
```

Record `TestWalkSimulation` (updates `tapes/TestWalkSimulation.tape` and
`walk_simulation_verdicts.tsv` against the pinned SQL Server 2022 container):

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/review/ -run '^TestWalkSimulation$'
```

Verify `TestWalkSimulation` read-only against the pinned SQL Server 2022
container (writes nothing to disk, fails on any mismatch):

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mssql/review/ -run '^TestWalkSimulation$'
```

Re-record all SQL Server oracle tapes:

```bash
OMNI_RECORD=1 go test -timeout 20m ./mssql/review/ -run '^(TestPriorBackupAgainstMSSQL|TestScanRulesAgainstMSSQL|TestWalkThroughAgainstMSSQL|TestWalkThroughBatchesAgainstMSSQL|TestWalkThroughSweepMSSQL|TestWalkSimulation)$'
```
