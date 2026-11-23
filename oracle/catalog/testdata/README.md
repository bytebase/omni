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

**2,472 / 2,472 rows (100.00%) agree** (up from 225 / 2,472 baseline and
2,348 / 2,472 initial simulation suite); **0 rows disagree**.

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

## Server configuration and catalog fidelity

`TestSimulation` configures the catalog with `StandardOracle23ServerConfig()`,
matching the pinned Oracle 23ai Free PDB:

- `MAX_STRING_SIZE=STANDARD` (`ORA-43929` on column `COLLATE` and table/view `DEFAULT COLLATION`).
- Non-sharded database (`ORA-02511` on sequence `SHARD` / `NOSHARD` unless enabled via `ALTER SESSION ENABLE SHARD DDL`).
- Non-application-root container (`ORA-65021` on `SHARING = METADATA | DATA | EXTENDED DATA`).
- TDE keystore closed (`ORA-28365` on column `ENCRYPT` after table, column, and FK checks).
- Isolated single-schema external-reference validation (`ORA-00942` on `REFERENCES other.t`).
- Live `ALTER TABLE ... SET UNUSED COLUMN` execution (with `ORA-00904`, `ORA-00957`, `ORA-12983`, `ORA-12992` checks).

Catalog state tracking also models database editions (seeded with `ORA$BASE`,
updated by `CREATE`/`DROP EDITION`, validated with `ORA-38802` on virtual column
edition clauses), numeric type inference for untyped virtual columns,
temporal-validity `PERIOD FOR` check constraints, and initial sequence current
values (`ORA-04007` / `ORA-04009` on `ALTER SEQUENCE`).
