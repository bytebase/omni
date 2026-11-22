# MariaDB query & expression view simulation

`simulate.tsv` implements Spec §10 Decision D6
(`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):
every unique generated query (`SELECT`, `WITH`, `TABLE`, `VALUES`) and `DO`
expression accepted by the pinned MariaDB (`../../../testdata/pins.tsv`,
`mariadb:11.8.8`) in `mariadb/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on the `one_table@column` fixture. MariaDB's stored
`SHOW CREATE VIEW` text is the engine's own reading of each query and
expression. `TestQuerySimulation` (`../../simulate_query_test.go`) replays it
from `../tapes/TestQuerySimulation.tape` with no server.

457 rows: 415 agree (100% of stored views matching byte-for-byte with 0 schema
mismatches) and 42 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./mariadb/catalog/ -run '^TestQuerySimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test ./mariadb/catalog/ -run '^TestQuerySimulation$'
```

Verify committed answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test ./mariadb/catalog/ -run '^TestQuerySimulation$'
```

## What the 42 rows still disagree on

| Rows | Category | Omni vs MariaDB | Why parked |
|---:|---|---|---|
| 33 | Unknown keyword-named functions (`AVG_ROW_LENGTH(a)`, `CPU(a)`, `ACTION()`, `a.SUM()`, …) | `ok` vs `1630` | Function existence checking requires the pinned server's complete builtin and UDF registry (and keyword-collision vs routine-lookup split), which varies across releases and loaded plugins. |
| 7 | Unknown identifier functions (`a()`, `a(a)`, `a(zzz)`, `a.a()`, `a.a.a()`, `` `a`() ``) | `ok`/`1054` vs `1305` | Same as above: omni's view catalog does not maintain a closed-world builtin/UDF registry or require external schemas' stored routines to be preloaded. |
| 2 | `CURRENT_TIME(2147483648)`, `CURRENT_TIME(9223372036854775808)` | `1426` vs `ok` | Upstream MariaDB 11.8.8 C++ signed integer overflow wraps out-of-range `opt_fsp` values into `curtime()` (0 precision) rather than rejecting with `1426` (`Too big precision specified`). Omni intentionally keeps rejecting out-of-range precision (`> 6`). |
