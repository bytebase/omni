# MySQL 8.4 LTS query simulation

`simulate.tsv` holds `../../queries/simulate.tsv`'s rows (every query and `DO`
expression MySQL accepts in `mysql/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on `one_table@column`), recorded against the pinned
MySQL 8.4 LTS (`../../../../testdata/pins.tsv`, role `mysql84`: `mysql:8.4.13`).
`TestQuerySimulationMySQL84` (`../../../simulate_mysql84_test.go`) replays it
from `../../tapes/TestQuerySimulationMySQL84.tape` with the catalog at `MySQL84`.

**698 rows: 698 agree, 0 disagree.**

## Commands

From `<repo root>`:

```bash
go test ./mysql/catalog/ -run '^TestQuerySimulationMySQL84$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestQuerySimulationMySQL84$'
```
