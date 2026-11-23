# Oracle catalog test data

`simulate.tsv` holds one row per core schema DDL statement that the pinned
Oracle 23ai Free (`../../testdata/pins.tsv`, `default` role:
`gvenzl/oracle-free:23.26.3-slim-faststart`, release `23.26.3.0.0`) accepts in
`oracle/parser/testdata/compat.tsv` (grammar-generated positive rules,
excluding `suite:*` hand corpora and near-miss `!` mutations), run once on each
fixture in `fixtures/`:

- Four namings of one plain table: `one_table@missing`, `one_table@table`,
  `one_table@column`, `one_table@index`.
- Four object-kind fixtures: `objects@view`, `objects@sequence`,
  `objects@synonym`, `objects@fk`.

Simulated statement families cover core relational schema DDL on tables
(`CREATE/ALTER/DROP TABLE`), indexes (`CREATE/DROP INDEX`), views
(`CREATE/ALTER/DROP VIEW`), sequences (`CREATE/ALTER/DROP SEQUENCE`), private
and public synonyms (`CREATE/DROP SYNONYM`), table/column comments
(`COMMENT ON TABLE/COLUMN`), `TRUNCATE TABLE`, and `RENAME` (309 statements ×
8 fixtures = 2,472 rows).

A row holds omni's outcome and Oracle's — the `ORA-` error number (without
leading zeros) when refused, or `ok:` and a hash of the stored schema DDL read
back from Oracle's data dictionary afterwards. `TestSimulation`
(`../simulate_test.go`) replays Oracle's side from `tapes/TestSimulation.tape`
with zero Docker containers. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.

**2,355 / 2,472 rows (95.27%) agree** (up from 225 / 2,472 baseline before
`oracle/catalog` `Exec` error-code and schema-state fidelity fixes); **117 rows
disagree** (4 only engine takes, 22 only omni takes, 15 other schema, 76 other
code).

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to Oracle's answer:

```bash
go test ./oracle/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row live against pinned Oracle 23ai Free again (needed when adding
rows or fixtures):

```bash
OMNI_RECORD=1 go test -v -timeout 30m ./oracle/catalog/ -run '^TestSimulation$'
```

Verify that live Oracle matches the committed tape and `simulate.tsv` with zero
differences:

```bash
OMNI_RECORD=verify go test -v -timeout 30m ./oracle/catalog/ -run '^TestSimulation$'
```

Show both schemas and the catalog diff for rows matching text:

```bash
go test -v ./oracle/catalog/ -run '^TestSimulation$' -args -simulate-show='DROP PRIMARY KEY'
```

## What the 117 remaining rows disagree on

### Deployment and server-mode configuration (93 rows)

These refusals depend on the server's deployment mode or configuration rather
than schema structure, so omni does not hardcode them:

| Rows | Oracle error | Statements | Why omni leaves them open |
|---:|---|---|---|
| 40 | `ORA-43929` | `ALTER TABLE ... COLLATE a`, `CREATE VIEW ... DEFAULT COLLATION a` (5 statements × 8 fixtures) | The `default` pin runs with `MAX_STRING_SIZE=STANDARD`, which refuses column/view collation clauses before checking the table or column. On `MAX_STRING_SIZE=EXTENDED` (`TestCollationMigrationRunsOnOracle`), Oracle parses and applies them. |
| 24 | `ORA-02511` | `CREATE/ALTER SEQUENCE a SHARD / NOSHARD` (3 statements × 8 fixtures) | Shard DDL is disabled on a non-sharded database. |
| 24 | `ORA-65021` | `CREATE SEQUENCE/SYNONYM ... SHARING = DATA \| METADATA` (3 statements × 8 fixtures) | `SHARING = DATA \| METADATA` is only valid inside an application root container. |
| 5 | `ORA-28365` | `ALTER TABLE a ADD/MODIFY (... ENCRYPT ...)` on non-FK columns | Transparent Data Encryption (TDE) keystore/wallet is not open on the default container. (Encrypting a foreign-key column is still refused with `ORA-28335` on `objects@fk`.) |

### Unmodeled catalog metadata and runtime state (24 rows)

| Rows | Category | Statements | Why omni leaves them open |
|---:|---|---|---|
| 11 | `other schema` | `ALTER TABLE a ADD (a AS (1) ...)` on `one_table@table` | `oracle/catalog` intentionally keeps `ColumnMetadata.Type` empty for untyped virtual columns so `DatabaseDefinition` round-trips `AS (expr)` verbatim (pinned by `TestVirtualColumnWithoutTypeStaysUntyped`), whereas Oracle infers `NUMBER`. |
| 4 | `other schema` | `CREATE TABLE a(PERIOD FOR a)`, `ALTER TABLE a ADD(PERIOD FOR a ...)` | Oracle synthesizes hidden `TIMESTAMP WITH TIME ZONE` start/end columns and a check constraint for temporal validity (`PERIOD FOR`), which `oracle/catalog` does not model yet. |
| 4 | `only engine takes` | `ALTER TABLE a SET UNUSED COLUMN a [ONLINE]` on `one_table@column` and `objects@fk` | `SET UNUSED` is intentionally refused by `oracle/catalog` (pinned by `TestSetUnusedIsRejected`); missing-table (`ORA-00942`) and missing-column (`ORA-00904`) checks still run first. |
| 3 | `ORA-38802` | `ALTER TABLE a ADD (a AS (1) EVALUATE USING EDITION a / UNUSABLE ... EDITION a)` on `one_table@table` | Editions are not tracked in `DatabaseSchemaMetadata`. |
| 2 | `ORA-04007` | `ALTER SEQUENCE a MINVALUE 1`, `ALTER SEQUENCE a NOMINVALUE` on `objects@sequence` | `ORA-04007` ("MINVALUE cannot exceed the current value") depends on runtime sequence `CURRVAL` state before the first `NEXTVAL`. |
