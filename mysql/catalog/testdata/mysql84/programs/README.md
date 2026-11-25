# MySQL 8.4 LTS stored-program simulation

`simulate.tsv` holds `../../programs/simulate.tsv`'s rows (stored procedures,
functions, triggers and events on five fixtures), recorded against the pinned
MySQL 8.4 LTS (`../../../../testdata/pins.tsv`, role `mysql84`: `mysql:8.4.13`).
`TestProgramSimulationMySQL84` (`../../../simulate_mysql84_test.go`) replays it
from `../../tapes/TestProgramSimulationMySQL84.tape` with the catalog at `MySQL84`.

**457 rows, all agree.**

MySQL 8.4 stores an event's `DISABLE ON SLAVE` as `DISABLE ON REPLICA`, an arm
8.4's `opt_ev_status` adds and 8.0.46's lacks. `mysql/parser` reads both, and
the catalog stores 8.4's spelling at `MySQL84` and refuses `DISABLE ON REPLICA`
with 8.0.46's 1064 at `MySQL80` (`../../programs/simulate.tsv`).

## Commands

From `<repo root>`:

```bash
go test ./mysql/catalog/ -run '^TestProgramSimulationMySQL84$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestProgramSimulationMySQL84$'
```
