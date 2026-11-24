# MySQL 8.4 LTS stored-program simulation

`simulate.tsv` holds `../../programs/simulate.tsv`'s rows (stored procedures,
functions, triggers and events on five fixtures), recorded against the pinned
MySQL 8.4 LTS (`../../../../testdata/pins.tsv`, role `mysql84`: `mysql:8.4.13`).
`TestProgramSimulationMySQL84` (`../../../simulate_mysql84_test.go`) replays it
from `../../tapes/TestProgramSimulationMySQL84.tape` with the catalog at `MySQL84`.

**455 rows: 452 agree, 3 disagree.**

## Known gap

The 3 rows that disagree set an event `DISABLE ON SLAVE`. MySQL 8.4 stores that
status as `DISABLE ON REPLICA`, an arm 8.4's `opt_ev_status` adds and 8.0.46's
lacks. `mysql/parser` holds to 8.0.46's grammar and has no version switch, so
the catalog cannot load what 8.4 stores and the rows read as a schema
difference. Closing them needs a version-aware parser.

## Commands

From `<repo root>`:

```bash
go test ./mysql/catalog/ -run '^TestProgramSimulationMySQL84$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestProgramSimulationMySQL84$'
```
