# TiDB query & expression view simulation

`simulate.tsv` implements Spec §10 Decision D6
(`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):
every unique generated query (`SELECT`, `WITH`, `TABLE`, `VALUES`) and `DO`
expression accepted by the pinned TiDB (`../../../testdata/pins.tsv`,
`pingcap/tidb:v8.5.5`) in `tidb/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on the `one_table@column` fixture. TiDB's stored
`SHOW CREATE VIEW` text is the engine's own reading of each query and
expression. `TestQuerySimulation` (`../../simulate_query_test.go`) replays it
from `../tapes/TestQuerySimulation.tape` with no server.

382 rows: 357 agree (187 stored views matching byte-for-byte, 170 matching
refusals) and 25 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./tidb/catalog/ -run '^TestQuerySimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestQuerySimulation$'
```

Verify committed answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test ./tidb/catalog/ -run '^TestQuerySimulation$'
```

## What the 25 rows still disagree on

| Rows | Category | Omni vs TiDB | Why parked |
|---:|---|---|---|
| 5 | `VALUES ROW(...)` views | `ok` vs `1051` | TiDB's grammar implements MySQL's `VALUES` query primary, but v8.5.5 cannot execute it; omni intentionally keeps `VALUES` view support (`table_values_view_test.go`). |
| 5 | `MATCH(...) AGAINST(...)` | `ok` vs `1105` | TiDB v8.5.5 does not implement full-text `MATCH ... AGAINST` (`UnknownType: *ast.MatchAgainst`). |
| 3 | `GROUP BY ... ASC\|DESC`, `SQL_CALC_FOUND_ROWS` | `1051` vs `1235` | Deployment/session-variable dependent (`tidb_enable_noop_functions`); same discipline as `tidb_enable_extended_stats` in `../README.md`. |
| 2 | `CAST(... AS ... ARRAY)`, `GROUPS` window frame | `ok` vs `1235` | Unimplemented feature refusals in the pinned TiDB release (`This version of TiDB doesn't yet support ...`). |
| 3 | `ident()`, `ident(ident)`, `POINT()` | `ok` vs `1305` | Unqualified function existence depends on the engine's full builtin registry, which changes across releases; schema-qualified calls (`ident.ident()`) are rejected with `1305` since TiDB has no stored functions. |
| 3 | `DATE 'binary'`, `TIME 'binary'`, `TIMESTAMP 'binary'` | `ok` vs `1292` | Value-dependent temporal literal validation. |
| 3 | `ident MEMBER OF(ident)`, `ident -> 'binary'`, `ident ->> 'binary'` | `ok` vs `3146` | Expression type-checking rejecting non-JSON/non-string arguments. |
| 1 | `APPROX_PERCENTILE(ident)` | `ok` vs `1105` | Uncoded `1105` message (`APPROX_PERCENTILE should take 2 arguments`). |
