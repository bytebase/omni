# MySQL Implicit Behaviors — Test Scenarios

> Goal: Verify omni's MySQL catalog correctly simulates the implicit / automatic behaviors
> that MySQL 8.0.45 performs at CREATE / ALTER time, against a real MySQL 8.0 container oracle.
>
> Source catalog: [`docs/plans/2026-04-13-mysql-implicit-behaviors-catalog.md`](../../docs/plans/2026-04-13-mysql-implicit-behaviors-catalog.md)
>
> Verification: Each scenario runs identical DDL on (1) a MySQL 8.0 testcontainer and (2) omni's in-memory catalog,
> then compares the observable state via `information_schema` / `SHOW CREATE TABLE` on MySQL and the equivalent
> `*Catalog` / `*Table` / `*Column` state in omni.
>
> Status legend: [ ] pending, [x] passing, [~] partial / known-limitation, [!] known omni bug
>
> Many scenarios were already spot-checked via `catalog_spotcheck_test.go` (2026-04-13) — those have
> known-correct expected values and are pre-marked `verified` in the progress tracker.

---

## Shared catalog state

Unless otherwise noted, scenarios assume:
- A fresh `catalog.New()` with `testdb` as the active database, created via
  `CREATE DATABASE testdb; USE testdb;`
- MySQL container default session (`sql_mode = ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,...`,
  `default_storage_engine = InnoDB`, `character_set_server = utf8mb4`,
  `collation_server = utf8mb4_0900_ai_ci`).
- `foreign_key_checks = 1` unless the scenario sets otherwise.

---

## Section C1: Name auto-generation

### 1.1 Foreign Key name — CREATE path (fresh counter)

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE child (
    a INT, CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES p(id),
    b INT, FOREIGN KEY (b) REFERENCES p(id)
);
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child' AND CONSTRAINT_TYPE='FOREIGN KEY'
ORDER BY CONSTRAINT_NAME;
```
Expected: `[child_ibfk_1, child_ibfk_5]` — CREATE counter starts at 0; user-named `_5` does NOT seed counter.

**omni assertion:**
```go
tbl := c.GetDatabase("testdb").GetTable("child")
// Collect FK constraint names, assert [child_ibfk_1, child_ibfk_5].
```

**See catalog:** C1.1a (`sql/sql_table.cc:9252`, `5912`). omni status: FIXED in commit `3202dab`.

---

### 1.2 Foreign Key name — ALTER path (max+1 counter)

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE child (
    a INT, b INT,
    CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES p(id)
);
ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES p(id);
```

**Oracle verification:** Same query as 1.1.
Expected: `[child_ibfk_20, child_ibfk_21]` — ALTER seeds from `max(existing generated number) + 1`.

**omni assertion:** FK list on `child` has names `{child_ibfk_20, child_ibfk_21}`.

**See catalog:** C1.1b (`sql/sql_table.cc:14345`, `5843-5877`).

---

### 1.3 Partition default naming `p0..p{n-1}`

**Setup:**
```sql
CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4;
```

**Oracle verification:**
```sql
SELECT PARTITION_NAME FROM information_schema.PARTITIONS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY PARTITION_ORDINAL_POSITION;
```
Expected: `[p0, p1, p2, p3]`.

**omni assertion:** `tbl.Partition.Partitions` has names `p0..p3`.

**See catalog:** C1.2 (`sql/partition_info.cc:567-583`).

---

### 1.4 CHECK constraint auto-name (table_chk_N)

**Setup:**
```sql
CREATE TABLE t (a INT, CHECK (a > 0));
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA='testdb';
```
Expected: `[t_chk_1]`.

**omni assertion:** `tbl` has exactly one CHECK constraint named `t_chk_1`.

**See catalog:** C1.4 (`sql/sql_table.cc:19007-19031`).

---

### 1.5 UNIQUE KEY auto-name uses field name

**Setup:**
```sql
CREATE TABLE t (a INT, UNIQUE KEY (a));
```

**Oracle verification:**
```sql
SELECT INDEX_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0;
```
Expected: `[a]`.

**omni assertion:** index on `t` with name `a`, non-unique=false.

**See catalog:** C1.5 (`sql/sql_table.cc:10377-10398`).

---

### 1.6 UNIQUE KEY name collision appends `_2`

**Setup:**
```sql
CREATE TABLE t (a INT, UNIQUE KEY a (a), UNIQUE KEY (a));
```

**Oracle verification:** Same STATISTICS query.
Expected: `[a, a_2]`.

**omni assertion:** two unique indexes on `t` with names `a` and `a_2`.

**See catalog:** C1.5 (`make_unique_key_name` loop).

---

## Section C2: Type normalization

### 2.1 REAL → DOUBLE

**Setup:**
```sql
CREATE TABLE t (c REAL);
```

**Oracle verification:**
```sql
SELECT COLUMN_TYPE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c';
```
Expected: `double`.

**omni assertion:** `tbl.GetColumn("c").Type` renders as `double` (or equivalent type enum).

**See catalog:** C2.1.

---

### 2.2 BOOL → TINYINT(1)

**Setup:**
```sql
CREATE TABLE t (c BOOL);
```

**Oracle verification:** Same query.
Expected: `tinyint(1)`.

**omni assertion:** column type renders as `tinyint(1)`.

**See catalog:** C2.2.

---

## Section C3: Nullability & default promotion

### 3.1 TIMESTAMP NOT NULL — FIRST column ONLY auto-promotes (omni risk)

**Setup:**
```sql
CREATE TABLE t (
    ts1 TIMESTAMP NOT NULL,
    ts2 TIMESTAMP NOT NULL
);
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected:
- `ts1` → `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`
- `ts2` → `TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00'` (NOT promoted)

**omni assertion:**
```go
ts1 := tbl.GetColumn("ts1")
// ts1.Default is CURRENT_TIMESTAMP, ts1.OnUpdate is CURRENT_TIMESTAMP
ts2 := tbl.GetColumn("ts2")
// ts2.Default is zero literal, ts2.OnUpdate is empty
```

**See catalog:** C3.1 (`sql/sql_table.cc:4221-4245`). **omni risk** flagged.

---

### 3.2 PRIMARY KEY implies NOT NULL

**Setup:**
```sql
CREATE TABLE t (a INT PRIMARY KEY);
```

**Oracle verification:**
```sql
SELECT IS_NULLABLE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a';
```
Expected: `NO`.

**omni assertion:** `tbl.GetColumn("a").Nullable == false`.

**See catalog:** C3.2.

---

### 3.3 AUTO_INCREMENT implies NOT NULL

**Setup:**
```sql
CREATE TABLE t (a INT AUTO_INCREMENT, KEY (a));
```

**Oracle verification:** Same query as 3.2.
Expected: `NO`.

**omni assertion:** `tbl.GetColumn("a").Nullable == false`.

**See catalog:** C3.3.

---

## Section C4: Charset / collation inheritance

### 4.1 Table charset inherits from database

**Setup:**
```sql
CREATE DATABASE db1 CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE db1;
CREATE TABLE t (c VARCHAR(10));
```

**Oracle verification:**
```sql
SELECT TABLE_COLLATION FROM information_schema.TABLES
WHERE TABLE_SCHEMA='db1' AND TABLE_NAME='t';
```
Expected: `utf8mb4_0900_ai_ci`.

**omni assertion:** `tbl.Charset == "utf8mb4"` and `tbl.Collation == "utf8mb4_0900_ai_ci"`.

**See catalog:** C4.1.

---

### 4.2 Column charset inherits from table (then elided in SHOW)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10)) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**Oracle verification:**
```sql
SELECT COLLATION_NAME FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c';
-- AND
SHOW CREATE TABLE t;
```
Expected: column collation = `utf8mb4_0900_ai_ci`; `SHOW CREATE TABLE` elides the column-level CHARACTER SET / COLLATE (matches table default), per C18.1.

**omni assertion:** column charset == table charset, and deparse output does not print a column-level CHARACTER SET clause.

**See catalog:** C4.2 + C18.1.

---

## Section C5: Constraint defaults

### 5.1 FK ON DELETE default → stored as RESTRICT, reported as NO ACTION, elided in SHOW

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id));
```

**Oracle verification:**
```sql
SELECT DELETE_RULE, UPDATE_RULE, MATCH_OPTION
FROM information_schema.REFERENTIAL_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c';
SHOW CREATE TABLE c;
```
Expected:
- `DELETE_RULE = 'NO ACTION'`
- `UPDATE_RULE = 'NO ACTION'`
- `MATCH_OPTION = 'NONE'`
- `SHOW CREATE TABLE` omits ON DELETE / ON UPDATE / MATCH clauses entirely.

**omni assertion:**
```go
fk := tbl.GetFKConstraint("...")
// fk.OnDelete, fk.OnUpdate, fk.Match are all the default (RESTRICT / SIMPLE)
// Deparse output for c has no ON DELETE / ON UPDATE / MATCH.
```

**See catalog:** C5.1, C5.2, C5.3. **Reporting discrepancy** — omni must map internal RESTRICT to `NO ACTION` when answering information_schema-style queries.

---

### 5.2 FK ON DELETE SET NULL requires nullable column

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (
    a INT NOT NULL,
    FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL
);
```

**Oracle verification:** MySQL errors (`ER_FK_COLUMN_NOT_NULL`).

**omni assertion:** `catalog.Exec(...)` returns an error for this DDL.

**See catalog:** C5.1 edge case (`sql/sql_table.cc:6682-6686`).

---

### 5.3 FK MATCH default rendered as NONE in information_schema

**Setup:** Same as 5.1.

**Oracle verification:** `REFERENTIAL_CONSTRAINTS.MATCH_OPTION = 'NONE'`.

**omni assertion:** If omni exposes an i_s-shaped view, it returns `NONE` for this FK.

**See catalog:** C5.3.

---

## Section C6: Partition defaults

### 6.1 HASH partition count defaults to 1 when PARTITIONS clause omitted

**Setup:**
```sql
CREATE TABLE t (id INT) PARTITION BY HASH(id);
```

**Oracle verification:**
```sql
SELECT COUNT(*) FROM information_schema.PARTITIONS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: 1.

**omni assertion:** `len(tbl.Partition.Partitions) == 1`, partition named `p0`.

**See catalog:** C6.1.

---

### 6.2 Subpartition count defaults to 1

**Setup:**
```sql
CREATE TABLE t (id INT, created DATE)
    PARTITION BY RANGE (YEAR(created))
    SUBPARTITION BY HASH (id)
    (PARTITION p0 VALUES LESS THAN (2020),
     PARTITION p1 VALUES LESS THAN MAXVALUE);
```

**Oracle verification:**
```sql
SELECT PARTITION_NAME, SUBPARTITION_NAME FROM information_schema.PARTITIONS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION;
```
Expected: each main partition has exactly one subpartition named `p0sp0` / `p1sp0`.

**omni assertion:** each `tbl.Partition.Partitions[i]` has 1 subpartition; names `{pN}sp0`.

**See catalog:** C6.2, C1.3.

---

### 6.3 Partition engine defaults to table engine

**Setup:**
```sql
CREATE TABLE t (id INT) ENGINE=InnoDB PARTITION BY HASH(id) PARTITIONS 2;
```

**Oracle verification:**
```sql
SELECT PARTITION_NAME, ENGINE FROM information_schema.PARTITIONS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: Both partitions report `ENGINE=InnoDB`.

**omni assertion:** `tbl.Partition.Partitions[i].Engine` either empty or equals table engine.

**See catalog:** C6.3.

---

## Section C7: Index defaults

### 7.1 Index algorithm defaults to BTREE

**Setup:**
```sql
CREATE TABLE t (a INT, KEY (a));
```

**Oracle verification:**
```sql
SELECT INDEX_TYPE FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='a';
```
Expected: `BTREE`.

**omni assertion:** `idx.Type == "BTREE"` (or default sentinel that deparses to BTREE).

**See catalog:** C7.1.

---

### 7.2 FK creates implicit backing index

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id));
```

**Oracle verification:**
```sql
SELECT INDEX_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c';
```
Expected: an index on `(a)` named `c_ibfk_1` (implicit, uses FK name).

**omni assertion:** `c` has an index on column `a`; name matches FK constraint name.

**See catalog:** C7.2.

---

## Section C8: Table option defaults

### 8.1 Storage engine defaults to InnoDB

**Setup:**
```sql
CREATE TABLE t (a INT);
```

**Oracle verification:**
```sql
SELECT ENGINE FROM information_schema.TABLES
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: `InnoDB`.

**omni assertion:** `tbl.Engine == "InnoDB"` (or empty that deparses to InnoDB).

**See catalog:** C8.1.

---

### 8.2 ROW_FORMAT defaults to DYNAMIC (InnoDB)

**Setup:** Same as 8.1.

**Oracle verification:**
```sql
SELECT ROW_FORMAT FROM information_schema.TABLES
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: `Dynamic`.

**omni assertion:** `tbl.RowFormat` == "DYNAMIC" (or default sentinel).

**See catalog:** C8.2.

---

### 8.3 AUTO_INCREMENT counter starts at 1 (new table)

**Setup:**
```sql
CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY);
```

**Oracle verification:**
```sql
SELECT AUTO_INCREMENT FROM information_schema.TABLES
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
SHOW CREATE TABLE t;
```
Expected: `AUTO_INCREMENT = 1`; `SHOW CREATE TABLE` ELIDES `AUTO_INCREMENT=` clause (C18.4).

**omni assertion:** `tbl.AutoIncrement` either 0 or 1; deparse does not emit `AUTO_INCREMENT=` clause.

**See catalog:** C8.3 + C18.4.

---

## Section C9: Generated column defaults

### 9.1 Generated column stored mode defaults to VIRTUAL

**Setup:**
```sql
CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a+1));
```

**Oracle verification:**
```sql
SELECT EXTRA FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b';
```
Expected: `VIRTUAL GENERATED`.

**omni assertion:** `tbl.GetColumn("b").Generated == GenVirtual` (or equivalent enum).

**See catalog:** C9.1.

---

### 9.2 FK on generated column is rejected

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (
    a INT,
    b INT GENERATED ALWAYS AS (a+1) STORED,
    FOREIGN KEY (b) REFERENCES p(id)
);
```

**Oracle verification:** MySQL errors (`ER_FK_CANNOT_USE_VIRTUAL_COLUMN`).

**omni assertion:** `Exec` returns an error.

**See catalog:** C9.2.

---

## Section C10: View metadata defaults

### 10.1 View ALGORITHM defaults to UNDEFINED

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE VIEW v AS SELECT a FROM t;
```

**Oracle verification:**
```sql
SELECT VIEW_DEFINITION, CHECK_OPTION, IS_UPDATABLE, SECURITY_TYPE
FROM information_schema.VIEWS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v';
SHOW CREATE VIEW v;
```
Expected:
- `CHECK_OPTION = 'NONE'`
- `SECURITY_TYPE = 'DEFINER'`
- `SHOW CREATE VIEW` contains `ALGORITHM=UNDEFINED`

**omni assertion:** view object reports `Algorithm == UNDEFINED`, `SqlSecurity == DEFINER`, `CheckOption == NONE`.

**See catalog:** C10.1, C10.2, C10.4.

---

### 10.2 View DEFINER defaults to current user

**Setup:** Same as 10.1.

**Oracle verification:**
```sql
SELECT DEFINER FROM information_schema.VIEWS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v';
```
Expected: `root@%` (or the connecting user).

**omni assertion:** view definer is set to the session user (or a sentinel that deparses to CURRENT_USER).

**See catalog:** C10.3.

---

## Section C11: Trigger defaults

### 11.1 Trigger DEFINER defaults to current user

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a;
```

**Oracle verification:**
```sql
SELECT DEFINER FROM information_schema.TRIGGERS
WHERE TRIGGER_SCHEMA='testdb' AND TRIGGER_NAME='trg';
```
Expected: the session user.

**omni assertion:** trigger definer is the session user.

**See catalog:** C11.1.

---

## Section C14: Constraint enforcement defaults

### 14.1 CHECK constraint defaults to ENFORCED

**Setup:**
```sql
CREATE TABLE t (a INT, CHECK (a > 0));
```

**Oracle verification:**
```sql
SELECT ENFORCED FROM information_schema.CHECK_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA='testdb';
```
Expected: `YES`.

**omni assertion:** `tbl.GetCheck("t_chk_1").Enforced == true`.

**See catalog:** C14.1.

---

## Section C15: Column positioning defaults

### 15.1 ALTER TABLE ADD COLUMN appends to end

**Setup:**
```sql
CREATE TABLE t (a INT, b INT);
ALTER TABLE t ADD COLUMN c INT;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY ORDINAL_POSITION;
```
Expected: `a=1, b=2, c=3`.

**omni assertion:** column ordering `[a, b, c]`.

**See catalog:** C15.1.

---

## Section C16: Date/time function precision defaults

### 16.1 NOW() precision defaults to 0

**Setup:**
```sql
CREATE TABLE t (ts DATETIME DEFAULT NOW());
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected: column rendered `ts datetime DEFAULT CURRENT_TIMESTAMP` (no `(n)` precision).

**omni assertion:** column default is `CURRENT_TIMESTAMP` with precision 0.

**See catalog:** C16.1.

---

## Section C18: SHOW CREATE TABLE elision rules

### 18.1 Column charset elided when equal to table default

**Setup:**
```sql
CREATE TABLE t (a VARCHAR(10), b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci)
  CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**Oracle verification:** `SHOW CREATE TABLE t` shows:
- `a` — no column-level CHARACTER SET / COLLATE
- `b` — `CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci` (explicit, non-primary)

**omni assertion:** deparse output matches the above structure — column-level charset appears only when it differs from table charset or is marked explicit.

**See catalog:** C18.1 (`sql/sql_show.cc:2086-2108`).

---

### 18.2 NOT NULL elision: TIMESTAMP shows NULL, others hide it

**Setup:**
```sql
CREATE TABLE t (
    i INT,
    i_nn INT NOT NULL,
    ts TIMESTAMP NULL DEFAULT NULL,
    ts_nn TIMESTAMP NOT NULL DEFAULT '2020-01-01'
);
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected:
- `i INT DEFAULT NULL` (NULL elided)
- `i_nn INT NOT NULL`
- `ts TIMESTAMP NULL DEFAULT NULL` (explicit NULL required)
- `ts_nn TIMESTAMP NOT NULL DEFAULT '2020-01-01 00:00:00'`

**omni assertion:** deparse produces the same 4 column defs.

**See catalog:** C18.2.

---

### 18.3 ENGINE clause always rendered (implicit default shown)

**Setup:**
```sql
CREATE TABLE t (a INT);
```

**Oracle verification:** `SHOW CREATE TABLE t` contains `ENGINE=InnoDB`.

**omni assertion:** deparse emits `ENGINE=InnoDB`.

**See catalog:** C18.3 (in practice always rendered for 8.0.45).

---

### 18.4 AUTO_INCREMENT clause elided when counter == 1

**Setup:**
```sql
CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY);
```

**Oracle verification:** `SHOW CREATE TABLE t` does NOT contain `AUTO_INCREMENT=`.

**omni assertion:** deparse output contains no `AUTO_INCREMENT=` clause.

**See catalog:** C18.4.

---

### 18.5 DEFAULT CHARSET clause always rendered

**Setup:**
```sql
CREATE TABLE tnocs (x INT);
```

**Oracle verification:** `SHOW CREATE TABLE tnocs` contains
`DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci`.

**omni assertion:** deparse emits both `DEFAULT CHARSET=utf8mb4` and `COLLATE=utf8mb4_0900_ai_ci`.

**See catalog:** C18.5 (spot-check 2026-04-13 overrides the "elided" reading — omni must ALWAYS render).

---

### 18.6 ROW_FORMAT clause elided when not explicitly specified

**Setup:**
```sql
CREATE TABLE t (a INT);  -- no ROW_FORMAT
CREATE TABLE t2 (a INT) ROW_FORMAT=DYNAMIC;  -- explicit
```

**Oracle verification:** `SHOW CREATE TABLE t` does NOT emit `ROW_FORMAT=`; `SHOW CREATE TABLE t2` emits `ROW_FORMAT=DYNAMIC`.

**omni assertion:** deparse emits `ROW_FORMAT=` only if `tbl.RowFormatExplicit == true`.

**See catalog:** C18.6.

---

## Section C21: Parser-level defaults

### 21.1 `DEFAULT` without value on nullable column → DEFAULT NULL

**Setup:**
```sql
CREATE TABLE t (a INT DEFAULT NULL);
```

**Oracle verification:**
```sql
SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a';
```
Expected: `COLUMN_DEFAULT = NULL`, `IS_NULLABLE = YES`.

**omni assertion:** `tbl.GetColumn("a").Default` is the NULL literal; nullable true.

**See catalog:** C21.1.

---

## Section C24: SHOW CREATE TABLE skip_gipk

### 24.1 Generated invisible primary key omitted from SHOW CREATE TABLE

**Setup:**
```sql
SET SESSION sql_generate_invisible_primary_key = ON;
CREATE TABLE t (a INT, b INT);
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;                 -- no my_row_id
SELECT COLUMN_NAME FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
SET SESSION show_gipk_in_create_table_and_information_schema = ON;
SHOW CREATE TABLE t;                 -- my_row_id visible
```

**omni assertion:** deparse under default session omits `my_row_id`; toggling the visibility flag shows it.

**See catalog:** C13.1 + C24.1. **omni risk** — likely not implemented.

---

## Section C25: DECIMAL defaults

### 25.1 DECIMAL with no precision/scale → DECIMAL(10,0)

**Setup:**
```sql
CREATE TABLE t (d DECIMAL);
```

**Oracle verification:**
```sql
SELECT COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d';
```
Expected: `decimal(10,0)`, precision=10, scale=0.

**omni assertion:** `tbl.GetColumn("d").Type` renders as `decimal(10,0)`.

**See catalog:** C25.1.

---

## Section PS: Path-split behaviors (CREATE vs ALTER)

These eight scenarios verify that omni's CREATE path and ALTER path do NOT get conflated
when they differ in MySQL. These are the highest-priority scenarios; several are already
covered by bug-fix tests.

### PS.1 CHECK constraint counter — CREATE path (fresh counter)

**Category:** Name auto-generation, refines C1.4
**MySQL source:** `sql/sql_table.cc:19068` (`prepare_check_constraints_for_create`)
**omni status:** FIXED (commit `3202dab`); verified by `TestBugFix_CheckCounterCreateTable`.

**Setup:**
```sql
CREATE TABLE t (
    a INT,
    CONSTRAINT t_chk_5 CHECK (a > 0),
    b INT,
    CHECK (b < 100)
);
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA='testdb'
ORDER BY CONSTRAINT_NAME;
```
Expected: `[t_chk_1, t_chk_5]`. Counter starts at 0; user-named `_5` is ignored.

**omni assertion:** `tbl` CHECK constraints named `{t_chk_1, t_chk_5}`.

---

### PS.2 CHECK constraint counter — ALTER path (max+1)

**MySQL source:** `sql/sql_table.cc:~19280` (`prepare_check_constraints_for_alter`)
**omni status:** VERIFIED CORRECT via `PS1_CheckCounter_ALTER_open`.

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, CONSTRAINT t_chk_20 CHECK (a>0));
ALTER TABLE t ADD CHECK (b>0);
```

**Oracle verification:** Same query as PS.1.
Expected: `[t_chk_20, t_chk_21]`.

**omni assertion:** names `{t_chk_20, t_chk_21}`.

---

### PS.3 FK counter — CREATE path (fresh, user-named NOT seeded)

**MySQL source:** `sql/sql_table.cc:9252`, `5912`.
**omni status:** FIXED (commit `3202dab`).

Identical to scenario 1.1 above. Listed again in PS section to make CREATE/ALTER split explicit.

---

### PS.4 FK counter — ALTER path (max+1 over existing generated numbers)

**MySQL source:** `sql/sql_table.cc:14345`, `5843-5877`.
**omni status:** FIXED — uses `nextFKGeneratedNumber` in `altercmds.go`.

Identical to scenario 1.2 above.

---

### PS.5 DEFAULT NOW() / fsp precision mismatch must error (parser-level, symmetric)

**MySQL source:** `sql/sql_parse.cc:5521` (`Alter_info::add_field`).
**omni status:** strictness gap to verify — spot-checked via `PS5_DatetimeFspMismatch`.

**Setup:**
```sql
CREATE TABLE t (a DATETIME(6) DEFAULT NOW());
```

**Oracle verification:** MySQL errors (`ER_INVALID_DEFAULT`).

**omni assertion:** `Exec` returns a parse / analyze error. If omni currently accepts, mark as strictness gap.

**See catalog:** PS5 (refines C16.x).

---

### PS.6 HASH partition ADD — seeded from count

**MySQL source:** `sql/sql_partition.cc:4506`.
**omni status:** N/A — omni has no ALTER TABLE ... ADD PARTITION support. Placeholder scenario, keep pending.

**Setup (once implemented):**
```sql
CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3;
ALTER TABLE t ADD PARTITION PARTITIONS 2;
```

**Oracle verification:** Partition names `[p0, p1, p2, p3, p4]`.

**omni assertion:** Partition names `[p0..p4]`.

---

### PS.7 FK name collision between user-named and auto-generated — must error

**MySQL source:** `sql/sql_table.cc:6614`.
**omni status:** **BUG** — omni silently succeeds. Spot-checked via `PS7_FKNameCollision`.

**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (
    a INT,
    CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id),
    b INT,
    FOREIGN KEY (b) REFERENCES p(id)
);
```

**Oracle verification:** MySQL errors (`ER_FK_DUP_NAME`, 1826).

**omni assertion:** `Exec` MUST return an error. Fix is a pre-insert collision check in CREATE path.

---

### PS.8 CHECK constraint duplicate name in schema — must error

**MySQL source:** `sql/sql_table.cc:19594-19601`.
**omni status:** BUG — no collision check today.

**Setup:**
```sql
CREATE TABLE t1 (a INT, CONSTRAINT my_rule CHECK (a > 0));
CREATE TABLE t2 (b INT, CONSTRAINT my_rule CHECK (b > 0));
```

**Oracle verification:** Second CREATE errors with `ER_CHECK_CONSTRAINT_DUP_NAME` (check names are schema-scoped).

**omni assertion:** second `Exec` returns an error.

---

## Progress tracker

Status values: `pending`, `verified` (spot-check done), `passing`, `bug` (omni bug known), `limitation`.

| # | Scenario | Status | Priority | Notes |
|---|----------|--------|----------|-------|
| 1.1 | FK name CREATE fresh counter | verified | HIGH | `C1_1_FK_counter_max_plus_one` + bugfix test |
| 1.2 | FK name ALTER max+1 | verified | HIGH | spot-check |
| 1.3 | Partition default naming p0..pN | verified | MED | `C1_2_partition_naming` |
| 1.4 | CHECK constraint auto-name | verified | MED | `C1_3_CheckConstraintName` |
| 1.5 | UNIQUE KEY field-name | pending | LOW | trivial |
| 1.6 | UNIQUE KEY name collision _2 | pending | LOW | |
| 2.1 | REAL → DOUBLE | verified | LOW | `C2_1_REAL_to_DOUBLE` |
| 2.2 | BOOL → TINYINT(1) | verified | LOW | `C2_2_BOOL_to_TINYINT1` |
| 3.1 | TIMESTAMP first-only promotion | verified | HIGH | `C3_1_timestamp_first_only_promotion` — omni risk |
| 3.2 | PRIMARY KEY → NOT NULL | verified | MED | `C3_2_primary_key_implies_not_null` |
| 3.3 | AUTO_INCREMENT → NOT NULL | verified | MED | `C3_3_AutoIncrement_implies_NOT_NULL` |
| 4.1 | Table charset from database | verified | HIGH | `C4_1_table_charset_from_database` |
| 4.2 | Column charset inherits + elided | verified | MED | `C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision` |
| 5.1 | FK ON DELETE default NO ACTION | verified | HIGH | `C5_1_fk_on_delete_default` — reporting discrepancy |
| 5.2 | FK SET NULL on NOT NULL errors | pending | MED | |
| 5.3 | FK MATCH default NONE | verified | MED | `C5_3_FK_MATCH_default` |
| 6.1 | HASH PARTITIONS defaults to 1 | verified | MED | `C6_1_Partition_default_count` |
| 6.2 | Subpartition count default 1 | verified | MED | `C6_2_Subpartition_count` |
| 6.3 | Partition engine inherits | pending | LOW | |
| 7.1 | Default index algorithm BTREE | verified | MED | `C7_1_Default_index_algorithm_BTREE` |
| 7.2 | FK backing index implicit | verified | MED | `C7_2_FK_backing_index` |
| 8.1 | Engine default InnoDB | verified | MED | `C8_1_Default_engine_InnoDB` |
| 8.2 | ROW_FORMAT default DYNAMIC | verified | MED | `C8_2_Default_row_format` |
| 8.3 | AUTO_INCREMENT starts at 1 | verified | MED | covered by 18.4 spot-check |
| 9.1 | Generated column VIRTUAL default | verified | LOW | `C9_1_GeneratedColumn_default_VIRTUAL` |
| 9.2 | FK on generated column errors | pending | LOW | |
| 10.1 | View ALGORITHM/CHECK/SECURITY defaults | verified | MED | `C10_1_3_4_View_defaults` |
| 10.2 | View DEFINER default | verified | MED | `C10_2_view_security_definer` |
| 11.1 | Trigger DEFINER default | verified | LOW | `C11_1_Trigger_definer_default` |
| 14.1 | CHECK enforced default | verified | LOW | `C14_1_Check_enforced_default` |
| 15.1 | ADD COLUMN appends to end | verified | LOW | `C15_1_Column_positioning_end` |
| 16.1 | NOW() precision default 0 | verified | MED | `C16_1_now_precision_default_zero` |
| 18.1 | Column charset elision | verified | HIGH | `C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision` |
| 18.2 | NOT NULL elision (TIMESTAMP) | verified | HIGH | `C18_2_NotNull_rendering` |
| 18.3 | ENGINE always rendered | pending | HIGH | |
| 18.4 | AUTO_INCREMENT elided at 1 | verified | HIGH | `C18_4_auto_increment_elision` |
| 18.5 | DEFAULT CHARSET always rendered | verified | HIGH | `C18_5_DefaultCharset_implicit` |
| 18.6 | ROW_FORMAT elided unless explicit | pending | HIGH | |
| 21.1 | DEFAULT NULL literal | verified | LOW | `C21_1_Default_NULL` |
| 24.1 | Invisible PK skip_gipk | pending | MED | omni risk — may not be implemented |
| 25.1 | DECIMAL default (10,0) | verified | LOW | `C25_1_decimal_default_10_0` |
| PS.1 | CHECK counter CREATE fresh | verified | HIGH | `PS1_CheckCounter_CREATE_fresh` + bugfix test |
| PS.2 | CHECK counter ALTER max+1 | verified | HIGH | `PS1_CheckCounter_ALTER_open` |
| PS.3 | FK counter CREATE fresh | verified | HIGH | same as 1.1 |
| PS.4 | FK counter ALTER max+1 | verified | HIGH | same as 1.2 |
| PS.5 | DEFAULT NOW() fsp mismatch errors | verified | MED | `PS5_DatetimeFspMismatch` |
| PS.6 | HASH partition ADD count-seed | pending | LOW | N/A — no partition ALTER support |
| PS.7 | FK name collision error | verified | HIGH | `PS7_FKNameCollision` — **omni bug** |
| PS.8 | CHECK schema-scoped dup name error | pending | MED | **potential omni bug** |

**Total scenarios:** 50 (13 pending, 35 verified via spot-check, 2 flagged as omni bugs).
