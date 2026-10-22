# WALK_THROUGH on every SQL Review engine: design

Status: decided 2026-10-22. §8 records the decisions. It builds on
[2026-10-19-walk-through-design.md](2026-10-19-walk-through-design.md) (the PostgreSQL design, "the PG spec" below), shipped in #493.

## 1. Goal

Bring WALK_THROUGH to every engine that has a SQL Review package in omni: MySQL (and OceanBase in MySQL mode, served by `mysql/review`), MariaDB, TiDB, MSSQL, Oracle, and Snowflake. The two requirements stay:

1. **No catalog.** The rule follows a change on the engine-neutral ledger `review/walk`, never on an engine's simulator.
2. **No false finding.** A finding is a statement that certainly fails. Misses are acceptable.

Why now:
- Bytebase's review run already sends WALK_THROUGH to omni for all eight engines (`reviewrun/review_funcs.go`). The rule is P0 and on in the default policy.
- Only `pg/review` implements it. The other six packages list it as "not implemented", so WALK_THROUGH reports nothing on them today.
- Bytebase's v1 plan check still simulates on a catalog for MySQL, MariaDB, OceanBase, and TiDB (`plugin/schema/mysql/walk_through.go`, `plugin/schema/tidb/walk_through.go`). MSSQL, Oracle, and Snowflake have no walk-through anywhere.

## 2. What carries over from PostgreSQL unchanged

- **The ledger:**
  - three-valued facts;
  - effects that are exact or May-effects over scopes;
  - an allowlist of statements, where anything else ends the walk;
  - the first certain violation is the finding.
- **What a finding claims** (PG spec §2): statement k fails, given that statements 1..k-1 succeeded. A later report after a failure the walk missed is conditional, so it is a recall gap, not a false finding.
- **Authority:** Present is trusted from the sync. Absent is trusted only in a space the sync lists completely.
- **`Target.OtherSpecs`:** UI specs on one database are unordered for every engine (`isSequentialTask`), so the flag means the same everywhere: an all-Unknown snapshot, and DML ends the walk.
- **Verification:** a curated oracle and a random sweep, both against the real engine.

## 3. Lessons from #493 that shape the plan

Codex reviewed #493 in eleven rounds. More than a dozen of its comments were real false positives, besides three the sweep found before review. Each class below gets a structural answer here, so it is not rediscovered per engine.

| What went wrong | Answer in this design |
| --- | --- |
| `IF NOT EXISTS` on a name that may exist was lowered as if the object were created (partition state, inline constraints, key indexes) | The shared builder offers `optional(prior)`: an Unknown prior turns every exact effect of the part into a May-effect (§4.1) |
| Effects within one statement were applied in written order, while the server runs `ALTER TABLE` in passes and `DROP` after all lookups | Each statement form declares an order policy, verified once per engine (§4.1) |
| Auxiliary state outlived its object (partition flags after `DROP SCHEMA` or a rename) | Engines keep no state beside the ledger unless a test pins its life cycle; prefer ledger kinds |
| `DROP TYPE` never lowered, because the parser gives a node shape the code did not expect | A coverage test: every row of an engine's statement table must lower, not stop (§4.2) |
| The oracle misjudged its own claims by parsing messages | Claims are structured: violation key, kind, and fact. The harness verifies the key, never the text (§4.2) |
| Scope grew round by round, and the PG side reached 3.5k lines against a 1.5k budget | The statement table and the "not covered" list are fixed before code. Review asks outside the table are declined by scope (§7) |

## 4. Shared work before the second engine

### 4.1 Lift the builder into `review/walk`

`pg/review/walk_through.go`'s `builder` is engine-neutral apart from its note type. It provides:
- parts, with `need`, `fail`, the effects, and `session`;
- the touched rule;
- `lookups`;
- `step()`.

Every engine needs it, and most of #493's fixes landed there. Move it to `review/walk` as `walk.Builder`:

- `Need(key, want, kinds, note any)`, `Fail(keys, got, note any)`, and the effects `Set`, `Clear`, `MayAdd`, `MayRemove`, `Forget`, `Move`, `Session(fn)`. `Last` ends the walk after the statement's checks.
- A failure's keys are its claim: `walk.Run` reports them in `Violation.Keys` with the facts it read for them (`Violation.Facts`), so an oracle confirms the claim key by key. A failure with no keys rests on the text.
- `Next()` starts the next part (subcommand or object).
- **Order policy per statement**, which replaces the PG builder's single touched rule and `lookups` flag:
  - `Lookups`: every check holds against the state before the statement, as in DROP of several objects.
  - `Passes(order)`: the engine orders the parts by its passes. A check on a key another part changes is dropped. PG's `ALTER TABLE` is drops first, then the rest.
  - `Sequential`, where each part sees the effects of the parts before it, is not built. It needs staged effects; build it when an engine's oracle shows one of its statements runs that way.
- `Optional()` marks a part the statement may skip, as with `IF NOT EXISTS` on a name that may exist:
  - its checks and failures are dropped;
  - a `Set` leaves its key Unknown, a `Clear` becomes a `MayRemove`, and a `Move` forgets both paths;
  - session effects run unchanged, so the engine adds only ones that hold either way.

  PostgreSQL keeps its own explicit handling of these cases, which its oracle pins; new engines use `Optional`.

pg/review then uses `walk.Builder`, with no change of behavior; its oracle and sweep must stay green.

### 4.2 A shared oracle harness: `review/walk/walktest`

This is test support with no dependency of its own. Each engine supplies an adapter:

```go
type Engine interface { // one change, with the target synced
	Statements(change string) []string
	Exec(t testing.TB, stmt string) Outcome // Code "" when accepted; a rejected statement leaves the database as it was
	Report(t testing.TB, change string) *Finding
	Violation(t testing.TB, change string) *walk.Violation
	Codes(message string) []string
	Holds(t testing.TB, v *walk.Violation, message string) (bool, string) // the facts of walktest.Claims(v), on the catalog
}
```

Resetting the database and syncing it stay in the engine's test, which builds an `Engine` per change.

The harness provides:
- `Judge(t, engine, change, complete)`: the accepted-prefix-plus-rejected-statement method, and the rule that a fully accepted change reports nothing.
- `Claims(v)`: the facts a violation rests on, which `Holds` confirms.
- `SweepParams`: the size, seed, and image of a sweep from `WALK_SWEEP_N`, `WALK_SWEEP_SEED`, and `WALK_SWEEP_IMAGE`. The templates stay per engine.
- `Coverage(t, forms, lowers)`: every statement form in the engine's table lowers, so none ends the walk unseen.

**Engines whose DDL commits implicitly** (MySQL family, Oracle, Snowflake) have no savepoints around DDL. The harness then relies on statement atomicity: a rejected statement changes nothing, so the live catalog right before the next statement is the state the walk judges against. Such an engine implements `Digester`, and `Judge` compares a catalog digest before and after each rejection. A server that applied part of a rejected statement (MySQL 5.7, or MariaDB before 10.6, dropping one table of `DROP TABLE a, missing`) fails the case. The pinned images are atomic.

### 4.3 Contract and plumbing

- Add `OtherSpecs` to an engine's `groupTargets` key (each is `{Schema, BackupDatabaseExists, ...}` today) in that engine's WALK_THROUGH PR, the first rule there to read it, with a test that pins it.
- **Name form rule for an unknown collation** (generalising PG spec §4.4):
  - an existing object is claimed only for an exact match of the stored spelling;
  - an absent one only when no ASCII case variant exists;
  - a name with a non-ASCII character is Unknown unless the spelling is identical.

  This holds whether the server compares case-sensitively or not.
- **Candidate `Target` fields** (decision D2):
  - `DefaultSchema` for MSSQL, which the sync does not record;
  - the current schema for Snowflake, which nothing sets at execution.

## 5. Engines

Each engine section gives:
- **scope:** what a Bytebase database is, and the current database or schema;
- **spaces:** which names share a namespace;
- **authority:** where Absent is trusted;
- **name form;**
- **resolution;**
- **stops:** what ends the walk;
- **a v1 statement table;**
- **facts to verify on the engine before code.**

The sync facts come from the Bytebase survey of 2026-10-22, from `backend/plugin/db/<engine>/sync.go` at origin/main 844a319bf0.

### 5.1 MySQL (and OceanBase, MySQL mode)

- **Scope:**
  - A Bytebase database is one MySQL database. The metadata holds a single schema named `""`.
  - The DSN selects the target database. A `USE` in the change carries over to later statements, because one connection runs them.
  - The walk tracks the current database: unqualified names are judged only while it is the target. Today's review only knows that some `USE` happened (`usesDatabase`).
- **Spaces:**
  - tables and views: one space per database;
  - columns: per table;
  - indexes, including `PRIMARY`: per table;
  - foreign key and check constraint names: per database (verify, D5);
  - triggers: per database;
  - procedures and functions: two spaces per database;
  - events: per database.
- **Authority (Absent trusted):**
  - tables (`BASE TABLE`), views, the columns and indexes of synced tables, foreign keys, triggers, routines, and events;
  - check constraints: never trusted as Absent. The sync reads them only from 8.0.16, and the server version is not an input.
  - OceanBase: no check constraints and no partitions in its sync. Its 5.7 version string makes it look older.
  - Assumed: the admin user sees every object of the database. `information_schema` lists only what it may see, and it is also the user that runs the change.
- **Name form:**
  - database and table names follow `Target.LowerCaseTableNames`: 0 compares exactly; 1 stores and compares lowercase; 2 stores as written and compares lowercase;
  - column, index, trigger, and routine names compare without case;
  - identifiers longer than 64 characters fail (1059) rather than truncate.
- **Stops:**
  - `CALL`, `PREPARE` / `EXECUTE`, and stored-program bodies;
  - DML wherever a synced or created trigger may fire;
  - assumed: events do not change the schema during the change, the analogue of PG's implicit-user-code assumption.
- **v1 statements:**
  - `CREATE`/`DROP DATABASE` (existence only);
  - `CREATE TABLE [IF NOT EXISTS]` (plain, `LIKE`, `AS SELECT`);
  - `DROP TABLE [IF EXISTS] a, b`;
  - `RENAME TABLE`;
  - `ALTER TABLE`: ADD/DROP/CHANGE/MODIFY/RENAME COLUMN; ADD/DROP INDEX, KEY, PRIMARY KEY, FOREIGN KEY, CHECK, CONSTRAINT; RENAME INDEX; RENAME TO;
  - `CREATE`/`DROP INDEX`;
  - `CREATE [OR REPLACE]`/`ALTER`/`DROP VIEW`;
  - `TRUNCATE`, DML targets, `USE`;
  - `CREATE`/`DROP TRIGGER`, `PROCEDURE`, `FUNCTION`, and `EVENT` as existence only.
  - Generated names (index named after its column, `t_ibfk_N`, `t_chk_N`) are prefix May-effects.
- **Verify:**
  - the order of `ALTER TABLE` clauses: are drops judged against the original table and adds in written order?
  - per-database FK and CHECK names, and the case of trigger and routine names;
  - atomic multi-table `DROP` and `RENAME` on the pinned 8.0 image;
  - `CREATE TABLE ... SELECT` failure atomicity.
- **Oracle environment:** `mysql:8.0` (testcontainers, already used). OceanBase has no CI container: DSN-gated (D4).

### 5.2 MariaDB (delta from MySQL)

- **Authority problem:** the sync keeps only `TABLE_TYPE = 'BASE TABLE'`, so sequences and system-versioned tables are missing. Both live in the table namespace. A system-versioned table that exists but is not synced would make `DROP TABLE t` look like "unknown table": a false finding. **Relation absence is untrusted on MariaDB until Bytebase syncs those types** (D5).
- **Dialect:**
  - `CREATE OR REPLACE` on tables, indexes, and more;
  - `IF [NOT] EXISTS` on `ALTER TABLE` subcommands.

  The omni MariaDB parser does not yet accept `ADD ... IF NOT EXISTS`, `ALTER TABLE IF EXISTS`, or `DROP INDEX IF EXISTS`. Those give Syntax findings today, and fixing that is parser work, outside this design.
- **Atomic DDL** since 10.6; pin an image at or above it.
- **Oracle environment:** `mariadb:11.8` (already used).

### 5.3 TiDB (delta from MySQL)

- Names always compare without case: `lower_case_table_names = 2` semantics, whatever the target says.
- No triggers, functions, or events (the parser rejects them). Sequences are synced, so they are trusted.
- `ALTER TABLE` with several changes exists only since 6.2 (verify). The AST cannot express `RENAME TO db.t`, so that form ends the walk.
- **Oracle environment:** a TiDB container (verify CI capacity).

### 5.4 MSSQL

- **Scope:**
  - A Bytebase database is one SQL Server database with many schemas.
  - Unqualified names resolve to the login's default schema, which Bytebase neither sets nor syncs. It is `dbo` for sysadmin logins (D2).
  - `USE` switches the database for the rest of the change, which today's review only detects.
- **Spaces:**
  - `sys.objects` names share one space per schema: tables, views, procedures, functions, PK/UQ/FK/CK/DF constraints, triggers, sequences, synonyms;
  - indexes: per table; columns: per table; types: their own space.
- **Authority:**
  - **Collision claims** ("already exists"): only on Present, since synonyms, CLR routines, and disabled triggers are not synced.
  - **Absence claims** for DDL that needs a table or view (`DROP TABLE`, `ALTER TABLE`): trusted, since `sys.tables` and views in the listed schemas are complete.
  - **Absence claims for DML and anything else that resolves through synonyms:** never trusted.
  - The database collation is not recorded, so names follow the unknown-collation rule of §4.3.
- **Statements:**
  - `ALTER VIEW`/`PROCEDURE`/`FUNCTION`/`TRIGGER` is indistinguishable from `CREATE OR ALTER` in the AST. The lowering reads the first token at `Loc.Start`: plain `ALTER` needs the object, while `CREATE OR ALTER` does not.
  - `sp_rename` with literal arguments is lowered; any other `EXEC` ends the walk.
  - `CreateTableStmt.IfNotExists` is never set by the parser.
- **Batches:**
  - A `GO` batch compiles as a unit. A compile error, such as a column a statement of the same batch adds, fails the batch before any of it runs.
  - That only adds failures, never a false finding.
  - Control flow (`IF`, `WHILE`, `TRY`) ends the walk.
- **Oracle environment:** the MSSQL testcontainer (already used by `mssql/review`).

### 5.5 Oracle

- **Scope:**
  - A Bytebase database is one schema (user).
  - `Open` sets `CURRENT_SCHEMA` to it, so unqualified names resolve there.
  - Other schemas are not synced: names in them are Unknown.
- **Spaces:**
  - one schema object space: tables, views, sequences, private synonyms, procedures, functions, packages, types, materialized views;
  - indexes, constraints, and triggers: a space each, per schema.

  The existing review already encodes this (`sharesTableNamespace`).
- **Authority:**
  - **Collision claims:** only on Present, since synonyms, types, and package bodies are not synced.
  - **Absence claims:** trusted for tables, views, columns, named constraints, and named indexes. Never trusted for `SYS_C*` constraints and `SYS_*` indexes, which the sync skips.
  - **DML:** unqualified names resolve through public synonyms, which are not synced, so DML absence is never trusted. Each DDL form is probed for whether it follows synonyms.
- **Name form:**
  - the lexer folds unquoted names to upper case, which is the stored form, so names compare exactly;
  - an identifier that is too long fails rather than truncates, and the limit (30 or 128 bytes) depends on the version, so the walk makes no claim on length.
- **Stops:**
  - PL/SQL blocks and `EXECUTE IMMEDIATE`, `CALL`;
  - `FLASHBACK TABLE ... TO BEFORE DROP`, since the recycle bin is unknown;
  - DML wherever a synced trigger may fire;
  - assumed: no schema-level DDL trigger, which the sync excludes, as PG assumes for a sibling spec's event triggers.
- **Oracle environment:** `gvenzl/oracle-free` (already used). Its container hung once in CI (#493's 4829c2fb run), so its connect must time out.

### 5.6 Snowflake

- **Scope:**
  - A Bytebase database is one Snowflake database with many schemas.
  - Bytebase sets no schema at execution, and the session default is unverified, so unqualified names are Unknown until the change runs `USE SCHEMA` (D2).
- **Name form:** unquoted names fold to upper case (`Ident.Normalize`); quoted ones keep their case.
- **Authority:**
  - The sync lists only `BASE TABLE` tables, views, columns, streams, and tasks, filtered by the default role's privileges.
  - External, event, dynamic, and temporary tables and materialized views share the table namespace (verify) but are not synced.
  - **Schemas:** trusted Absent.
  - **Relations:** absence trusted only for DDL forms that no unsynced kind satisfies (to probe).
- **Statements:** `CREATE OR REPLACE`, `UNDROP`, `SWAP WITH`, and `CLONE` are common. `AlterTableAction` is not an AST node, so `ast.Inspect` does not reach it.
- **Oracle environment:** none. There is no emulator, so precision cannot be verified in CI (D3).

## 6. Bytebase side

Required for the rule to run:
- set `Target.OtherSpecs` in `reviewrun`;
- bump the omni pin per engine.

The rule is already registered for all eight engines through `reviewFuncs`, so no registration step is needed.

Sync improvements, each turning Unknown into trusted facts. None is required to start:
- MariaDB: sequences and system-versioned tables. Required before MariaDB relation absence can be trusted.
- MSSQL: the database collation and the login's default schema.
- Oracle: public synonyms, and an explicit recycle-bin filter.
- Snowflake: the other table kinds, materialized views, and sequences.

v1 and V2: v1's plan-check walk-through keeps running until V2 covers an engine. Retiring it per engine is a Bytebase decision after V2 ships there.

## 7. Delivery

One PR per step, each with its own spec section and size budget:

1. **Shared:** `walk.Builder`, its order policies, and `Optional`; the failure facts in `Violation`; `walktest`; pg/review moved onto both, with a coverage test of its statement table. The PG oracle and sweep must stay green, with no behavior change.
2. **MySQL**, with OceanBase's authority deltas. Budget: lowering about 1.2k lines, since the builder and harness are shared.
3. **MariaDB and TiDB** as deltas of MySQL.
4. **MSSQL**, after D2.
5. **Oracle.**
6. **Snowflake**, after D3.

Review policy for every engine PR:
- A false finding is fixed.
- A miss is fixed when it is cheap and inside the engine's statement table; otherwise it is declined with the §5 "not covered" list.
- Growth past the budget is a decision for a person.

## 8. Decisions

Decided 2026-10-22, each as recommended:

- **D1. The builder and the harness move to `review/walk` first (§4)**, before the second engine.
- **D2. Default and current schema:**
  - Bytebase adds `Target.DefaultSchema`, read at sync as the admin login: `SELECT SCHEMA_NAME()` for MSSQL, `CURRENT_SCHEMA()` for Snowflake.
  - Until the field lands, unqualified names on MSSQL and Snowflake are Unknown.
- **D3. Snowflake waits for a DSN-gated test account.** Without an oracle, its precision cannot be shown. A text-only subset is the fallback only if no account can be had.
- **D4. OceanBase is verified by a DSN-gated oracle run.** Its container is too heavy for CI.
- **D5. MariaDB relation absence is untrusted** until Bytebase syncs sequences and system-versioned tables. The Bytebase fix is to be filed.
- **D6. Accepted assumptions across engines:**
  - the admin user sees every object of its database (MySQL family, Snowflake);
  - MySQL events change no schema while a change runs;
  - no Oracle schema has a schema-level DDL trigger;
  - the role running a change may use every schema of its search path. This was pending from #493 and is now confirmed for PostgreSQL too.
- **D7. Order as in §7**, one PR per step, each with its budget set in its spec section before code.
