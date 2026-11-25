# Oracle review test data

This directory holds recorded Oracle 23 verdicts and SQL tapes for `oracle/review`,
recorded against the pinned server in `../../testdata/pins.tsv`
(`gvenzl/oracle-free:23.26.3-slim-faststart`). All tests replay 100% offline with
zero Docker in CI.

## `walk_simulation_verdicts.tsv` and `TestWalkSimulation`

`walk_simulation_verdicts.tsv` holds every unique Oracle 23-accepted DDL and DML
statement in `../../parser/testdata/compat.tsv` (3,077 statements across tables,
views, materialized views, indexes, sequences, synonyms, procedures, functions,
packages, types, triggers, comments, renames, and DML — excluding materialized
view logs, which trigger Oracle server parse hangs) executed against 11 starting
fixtures in `fixtures/`:

- `one_table@missing` — base table `t` only (`a` absent from every namespace)
- `one_table@table` — table `a(b PRIMARY KEY, c)` with index `a_c_idx`
- `one_table@column` — table `a(b, a)` with check constraint `a` and no primary key
- `objects@partitioned` — range-partitioned table `a` with partition `a`
- `objects@index_trigger` — index `a` and trigger `a` on table `t`
- `objects@view` — view `a` on table `t`
- `objects@matview` — materialized view `a` on table `t`
- `objects@sequence` — sequence `a`
- `objects@synonym` — private synonym `a` for table `t`
- `objects@function` — function `a`
- `objects@type` — object type `a` (unsynced in Bytebase's Oracle metadata sync)

Total: **33,847 rows** (3,753 accepted `ok`, 30,094 refused).
`TestWalkSimulation` (`../walk_simulation_test.go`) enforces **zero false positives**
on all 3,753 Oracle-accepted rows and verifies every finding on refused rows against
either Oracle's error code or the fixture's real Oracle inventory synced via
`tapes/TestWalkSimulation.tape`. Currently catches 25,768 / 30,094 refusals (85.6%).

Its `queries` subtest walks `../../catalog/testdata/queries/simulate.tsv`, each
query as `CREATE VIEW v AS <query>` on the catalog's fixtures, and asks Oracle for
each verdict on the same tape. The walk asks `oracle/catalog` about a view's query
when every table it reads is synced in full, and catches 205 / 256 refusals (80.1%)
with zero false positives. It leaves the rest to Oracle, as the sync cannot show
them: a missing column or name, which a function, a public synonym or another
schema may supply (00904, 00942); a type mismatch, which an object or collection
type may cause (00932, 01790, 22905); flashback and periods (08181, 55603); database
links (02019); `GROUP BY 1` (03162, by `GROUP_BY_POSITION_ENABLED`); and a view
the query reads (00957). `FOR UPDATE` in a view (03048) is a syntax error the parser
reports, which the subtest does not count.

Its `programs` subtest walks `../../catalog/testdata/programs/simulate.tsv` the
same way: procedures, functions, packages, triggers and types with bodies, and
DDL on the objects they read, on `programs@objects`. Oracle creates a unit it
cannot compile, so most rows are accepted; the walk catches the refusals it
models (an existing name, a missing object to drop, a trigger on a missing
table, an INSTEAD OF trigger on a table, and what `catalog.TriggerFormError`
finds in a trigger's own text), 27 / 39 refusals (69.2%), and is held to zero
false positives and a recall floor of its own. It leaves the rest to Oracle,
as the sync cannot show them: types (`CREATE TYPE`, `DROP TYPE`, a column of a
missing type), a package's body (`CREATE PACKAGE BODY` of one that has a body;
Bytebase's sync appends it to the package's definition, the walk's sync does
not read it), triggers on another schema's tables (`DROP TRIGGER`, `ALTER
TRIGGER` of a missing one), whether a view is an editioning view (a `BEFORE`
trigger on one), and other users (`CREATE PROCEDURE nope.p`). `RENAME` fails
with 03001 for a user who does not own the schema, as the walk's does.

## Commands

Run from `<repo root>`.

Offline replay (default, zero Docker):

```bash
go test -v ./oracle/review/ -run '^TestWalkSimulation$'
```

Re-record `TestWalkSimulation` (`walk_simulation_verdicts.tsv` and `tapes/TestWalkSimulation.tape`) against pinned Oracle 23:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestWalkSimulation$'
```

Read-only verification against pinned Oracle 23 (writes nothing to disk):

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./oracle/review/ -run '^TestWalkSimulation$'
```

Record or verify the other Oracle-backed review tapes:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestWalkThroughAgainstOracle$'
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestWalkThroughSweepOracle$'
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestPriorBackupAgainstOracle$'
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestScanRulesAgainstOracle$'
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestRenameAgainstOracle$'
OMNI_RECORD=1 go test -v -timeout 20m ./oracle/review/ -run '^TestRequireIsNullAgainstOracle$'
```
