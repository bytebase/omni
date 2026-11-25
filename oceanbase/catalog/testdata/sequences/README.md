# OceanBase sequence catalog simulation

`simulate.tsv` holds `CREATE`, `ALTER` and `DROP SEQUENCE` statements judged
against the pinned OceanBase CE MySQL-mode image (`../../../testdata/pins.tsv`,
`oceanbase/oceanbase-ce:4.3.5.6-106000012026040916`), run across five fixtures
in `../fixtures/`:

- `one_table@missing` — base table `t` (no object named `a`)
- `one_table@table` — base table `a`
- `objects@view` — base table `t` and view `a`
- `objects@sequence` — sequence `a` with the default options
- `objects@sequence_options` — descending, cyclic, ordered sequence `a` with
  every option set

`TestSequenceSimulation` (`../../simulate_sequence_test.go`) takes the 25
sequence statements OceanBase accepts in `oceanbase/parser/testdata/compat.tsv`
(one per grammar alternative, `IF [NOT] EXISTS` at every site) and a curated
matrix of 96 more: each check of OceanBase's option builder and the order it
makes them in, numbers past 64 bits and with fractions, every `ALTER SEQUENCE`
option on a sequence that exists, and names in another database, in another
case and too long.

OceanBase has no `SHOW CREATE SEQUENCE` and lists no sequence in `SHOW FULL
TABLES`, so the test reads each sequence from `oceanbase.DBA_SEQUENCES` and
writes it as one `CREATE SEQUENCE` with every option; tables and views are
compared as `../simulate.tsv` compares them. It replays from
`../tapes/TestSequenceSimulation.tape` without starting a server.
`oceanbase/review`'s `TestWalkSimulation` also checks every row offline
(`sequences` subtest) with 0 false positives: Bytebase's sync reads no
sequence, so the walk finds only the refusals that do not depend on the
database, 280 of 460.

**605 rows (121 statements × 5 fixtures): 605 agree, 0 disagree (100.00%).**

## Commands

From `<repo root>`. Commit rows a catalog fix moved to OceanBase's outcome:

```bash
go test ./oceanbase/catalog/ -run '^TestSequenceSimulation$' -update-compat
```

Record all rows against pinned `oceanbase/oceanbase-ce:4.3.5.6-106000012026040916`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./oceanbase/catalog/ -run '^TestSequenceSimulation$'
```

Verify recorded answers against pinned `oceanbase/oceanbase-ce:4.3.5.6-106000012026040916` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./oceanbase/catalog/ -run '^TestSequenceSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./oceanbase/catalog/ -run '^TestSequenceSimulation$' -args -simulate-show='ALTER SEQUENCE'
```
