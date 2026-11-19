# TiDB catalog test data

`simulate.tsv` holds one row per DDL statement of the simulated families that
the pinned TiDB (`../../testdata/pins.tsv`, `pingcap/tidb:v8.5.5`) takes in
`tidb/parser/testdata/compat.tsv`, run once on each fixture in `fixtures/`:
four namings of one table, `@missing`, `@table`, `@column` and `@index`. A row
holds omni's outcome and the engine's — an error number, or `ok:` and a hash of
the `SHOW CREATE` texts the database holds afterwards. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` and starts
no server. The design is `docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`
§3. 1,724 rows: 1,456 agree and 268 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./tidb/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned server again, which a new row or fixture needs:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestSimulation$'
```

Show both schemas of the rows that hold a statement's text:

```bash
go test ./tidb/catalog/ -run '^TestSimulation$' -args -simulate-show='PARTITION BY'
```

Record `TestSimulationCausesOracle`'s tape after changing what it sends. Each
cause of the file is one case there, answered by the server and by the catalog:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestSimulationCausesOracle$'
```

## What the 268 rows still disagree on

### 8200: TiDB's DDL worker refuses forms its parser takes (204 rows)

**Parked, on purpose.** 68 distinct `ALTER TABLE` forms, whose omni outcome is
`ok`, `1060` or `1068` and whose engine outcome is `8200`. They are three kinds:
unimplemented `ALTER TABLE` specifications and table options, a *value*-dependent
refusal (`ALTER TABLE a CHARSET utf8mb4` is taken where `CONVERT TO CHARSET
'binary'` is 8200), and TiDB's refusal to add a `UNIQUE`, `PRIMARY KEY` or
`AUTO_INCREMENT` column in the same `ALTER`, which would replace omni's sharper
1060/1068 with a blanket 8200. Mirroring them would encode the pin's
implementation gaps. The judgment is #1464's and nothing here moves those rows.

### The rest, measured and left (64 rows)

Each of these is a property of the deployment or of a name's registry, not a
rule about the statement, or an answer that depends on the day. The messages
are the pinned server's, read on 2026-11-19. omni answers `ok` — it takes the
statement and stores what the grammar says — except where a count is shown.

| Rows | Statement | TiDB | Why omni does not state it |
|---:|---|---|---|
| 12 | `CREATE TABLE ident SELECT *` (and `… UNION …`, and `WITH … SELECT`) | `1105 'CREATE TABLE ... SELECT' is not implemented yet` | A feature the pin has not implemented. |
| 12 | `ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident)`, `… DEPENDENCY(ident)`, `… CORRELATION(ident)`, `DROP STATS_EXTENDED ident` | `1105 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF` | Extended statistics are behind a GA flag and a system variable, so the answer is the deployment's. |
| 3 | `ALTER TABLE ident SET TIFLASH REPLICA 1` | `1105 the tiflash replica count: 1 should be less than the total tiflash server count: 0` | The count is the cluster's: the same statement is taken on a cluster with a TiFlash store. `SET TIFLASH REPLICA 0` is taken here, and the catalog's differential tests use 0. |
| 6 | `ALTER TABLE ident COMPACT PARTITION ident` (and `… TIFLASH REPLICA`) | `1105 table:ident is not a partition table, but user specify partition name list:[ident]` | TiDB's other partition actions answer 1505 on the same table (the causes oracle holds those); COMPACT answers an uncoded 1105 whose message names the partition list. The catalog records that it takes the statement rather than mirror an internal message. |
| 3 | `ALTER TABLE ident AUTO_ID_CACHE 1` | `1105 Can't Alter AUTO_ID_CACHE between 1 and non-1, the underlying implementation is different` | A value-dependent claim about a table metadata option: `AUTO_ID_CACHE 0` on the same table is taken. |
| 3 | `ALTER TABLE ident PARTITION ident COMMENT 'binary'` | `1105 unknown partition option` | The compat row samples a partition option the grammar derives and TiDB's DDL does not read. The answer is an internal "unknown option" for the sampled text. |
| 8 | `ALTER TABLE ident ENGINE 'binary'`, `… STORAGE ENGINE 'binary'` | `1286 Unknown storage engine 'binary'` | A name in TiDB's engine registry, checked while the statement is prepared (so even with no table: 2 of the 8 rows). TiDB takes `InnoDB`, `MyISAM`, `MEMORY`, `CSV`, `ARCHIVE`, `BLACKHOLE`, `FEDERATED` and `MERGE` and refuses anything else, `TiFlash`, `ROCKSDB` and `tikv` among them. |
| 4 | `ALTER TABLE ident INSERT_METHOD 'binary'` | `8233 CREATE/ALTER table with insert method option is not supported` | TiDB's CREATE/ALTER layer does not implement the option, whatever its value. |
| 4 | `ALTER TABLE ident UNION()` | `8232 CREATE/ALTER table with union option is not supported` | Same: the MERGE-table option is refused as unimplemented. |
| 3 | `ALTER TABLE ident ATTRIBUTES 'binary'` | `8237 Invalid attributes: attributes should be in format 'key=value': binary` | The sampled value `'binary'` is not a `key=value` pair; the value decides the answer, so the statement's rule is not the message's. |
| 3 | `ALTER TABLE ident AFFINITY 'binary'` | `8266 Invalid AFFINITY 'binary'` | A value-dependent table option: the sampled `'binary'` is not one of TiDB's affinities. |
| 3 | `ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE` (and `CURRENT_DATE()` and `CURDATE()`) | `1067 Invalid default value for 'ident'` | The date evaluates to a number (20251119 on the day this was read) and a TINYINT cannot hold it, while the same default on an INT column is taken. The answer changes with the day, so the catalog does not state it. #1464 kept these three rows for the same reason. |

`TiFlash`'s and the engine registry's answers above are also why no case in
`TestSimulationCausesOracle` states them: a catalog fix for them would have to
carry a list of names a release can change, which is the deployment's, not the
grammar's.
