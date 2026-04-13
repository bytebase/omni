# Section PS — CREATE vs ALTER path-split bugs

Bugs discovered while implementing `TestScenario_PS` in
`mysql/catalog/scenarios_ps_test.go`. These are expected failures; the
scenario tests still run the assertion so the bugs surface in CI as
regressions-to-fix.

### PS.2 CHECK constraint counter — ALTER path uses fresh counter (regression)
- **Discovered**: 2026-04-13
- **Section**: PS
- **Scenario**:
  ```sql
  CREATE TABLE t (a INT, b INT, CONSTRAINT t_chk_20 CHECK (a>0));
  ALTER TABLE t ADD CHECK (b>0);
  ```
- **Expected (MySQL 8.0.45)**: `[t_chk_20, t_chk_21]`. The ALTER path
  seeds the CHECK counter from the max existing generated number.
- **Actual (omni)**: `[t_chk_1, t_chk_20]`. The ALTER path generates
  `t_chk_1` instead of `t_chk_21`, indicating `nextCheckNumber` for
  ALTER is NOT scanning existing CHECK names for the max suffix.
  Confirms the status note in the dispatch prompt was optimistic;
  omni's ALTER check-name counter needs the same max+1 logic applied to
  FK counter in `altercmds.go` (`nextFKGeneratedNumber`).
- **Severity**: HIGH
- **Fix hint**: in `mysql/catalog/altercmds.go`, the CHECK constraint
  ADD path for ALTER TABLE should call a `nextCheckGeneratedNumber`
  helper that mirrors `nextFKGeneratedNumber` — iterate existing
  CHECK constraints matching `<table>_chk_<N>` and return max+1.
  Reference `sql/sql_table.cc:~19280` (`prepare_check_constraints_for_alter`).

### PS.5 DATETIME(6) DEFAULT NOW() fsp mismatch not rejected
- **Discovered**: 2026-04-13
- **Section**: PS
- **Scenario**: `CREATE TABLE t (a DATETIME(6) DEFAULT NOW())`
- **Expected (MySQL 8.0.45)**: `ER_INVALID_DEFAULT` (1067) — the implicit
  `NOW()` has fsp=0 but column has fsp=6, so the default cannot be
  represented losslessly.
- **Actual (omni)**: accepts the DDL silently.
- **Severity**: HIGH
- **Fix hint**: strictness / analyze gap. The check belongs near
  column default validation in `mysql/catalog/tablecmds.go` or the
  analyze phase that builds column specs. Reference `sql/sql_parse.cc:5521`
  (`Alter_info::add_field`).

### PS.7 FK name collision between user-named and auto-generated not detected
- **Discovered**: 2026-04-13
- **Section**: PS
- **Scenario**:
  ```sql
  CREATE TABLE p (id INT PRIMARY KEY);
  CREATE TABLE c (
      a INT,
      CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id),
      b INT,
      FOREIGN KEY (b) REFERENCES p(id)
  );
  ```
- **Expected (MySQL 8.0.45)**: `ER_FK_DUP_NAME` (1826) — the counter
  generates `c_ibfk_1` for the second (unnamed) FK, which collides with
  the user-named `c_ibfk_1`.
- **Actual (omni)**: silently succeeds. The CREATE path generates a
  non-colliding auto name because the user-named entry is NOT seeded,
  but MySQL still rejects the collision by name.
- **Severity**: HIGH
- **Fix hint**: add a pre-insert collision check in the CREATE path of
  `mysql/catalog/tablecmds.go` that compares each assigned FK constraint
  name against all other FK names in the same CREATE TABLE. Reference
  `sql/sql_table.cc:6614`.

### PS.8 CHECK schema-scoped name collision not detected
- **Discovered**: 2026-04-13
- **Section**: PS
- **Scenario**:
  ```sql
  CREATE TABLE t1 (a INT, CONSTRAINT my_rule CHECK (a > 0));
  CREATE TABLE t2 (b INT, CONSTRAINT my_rule CHECK (b > 0));
  ```
- **Expected (MySQL 8.0.45)**: second CREATE fails with
  `ER_CHECK_CONSTRAINT_DUP_NAME` (3822). CHECK constraint names are
  schema-scoped — they must be unique across ALL tables in the schema.
- **Actual (omni)**: both CREATE statements succeed.
- **Severity**: MED
- **Fix hint**: during CREATE/ALTER, scan all tables in the target
  database for existing CHECK constraints with the same name and reject
  with `ER_CHECK_CONSTRAINT_DUP_NAME` if found. Reference
  `sql/sql_table.cc:19594-19601`.

### PS.6 HASH partition ADD not supported (placeholder)
- **Discovered**: 2026-04-13
- **Section**: PS
- **Scenario**:
  ```sql
  CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3;
  ALTER TABLE t ADD PARTITION PARTITIONS 2;
  ```
- **Expected (MySQL 8.0.45)**: partition names `[p0, p1, p2, p3, p4]`.
- **Actual (omni)**: omni has no `ALTER TABLE ... ADD PARTITION` support;
  ALTER either errors or leaves partition list unchanged.
- **Severity**: LOW (omni-wide partition ALTER gap, tracked separately)
- **Fix hint**: add HASH partition ADD support in
  `mysql/catalog/altercmds.go`. Counter seeds from existing partition
  count. Reference `sql/sql_partition.cc:4506`.
