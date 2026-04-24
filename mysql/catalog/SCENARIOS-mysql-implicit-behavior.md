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

> **Expansion note (Wave 2):** grew from 6 to 13 scenarios via systematic
> walk of MySQL 8.0 create-table-foreign-keys.html,
> create-table-check-constraints.html, create-index.html +
> targeted re-read of `sql/sql_table.cc` covering
> `prepare_key` (index name generation from first column +
> PRIMARY reservation, ~L7229-L7265),
> `make_unique_key_name` (L10377-L10398),
> `make_functional_index_column_name` (L7710-L7743),
> `add_functional_index_to_create_list` (L7783-L7808),
> `generate_fk_name` (L5906-L5948),
> `prepare_check_constraints_for_create` +
> `generate_check_constraint_name` (L19007-L19031, L19583-L19602),
> and a spot-check of `mysql/catalog/tablecmds.go` (`allocIndexName`,
> `nextFKGeneratedNumber`, `nextCheckNumber`, `ensureFKBackingIndex`).

---

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

### 1.7 PRIMARY KEY always named "PRIMARY" (user name ignored)

> **new in expansion (Wave 2)** — priority: HIGH — status: pending-verify

**Setup:**
```sql
CREATE TABLE t (
    a INT,
    CONSTRAINT my_pk PRIMARY KEY (a)
);
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME, CONSTRAINT_TYPE FROM information_schema.TABLE_CONSTRAINTS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';

SELECT INDEX_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: CONSTRAINT_NAME = `PRIMARY` (NOT `my_pk`). The `CONSTRAINT my_pk` clause is silently discarded — MySQL hard-codes `key_info->name = primary_key_name` regardless of the user-supplied `CONSTRAINT symbol`. STATISTICS INDEX_NAME = `PRIMARY`.

**omni assertion:** `tbl.PrimaryKey.Name == "PRIMARY"`, no constraint named `my_pk` exists on `tbl`.

**Anchors:**
- `sql/sql_table.cc:210` — `const char *primary_key_name = "PRIMARY";`
- `sql/sql_table.cc:7236-7237` — `if (key->type == KEYTYPE_PRIMARY) key_info->name = primary_key_name;`
- doc: [create-table.html](https://dev.mysql.com/doc/refman/8.0/en/create-table.html) — "The name of a PRIMARY KEY is always PRIMARY, which thus cannot be used as the name for any other kind of index."

**omni gap:** need to confirm omni discards `CONSTRAINT my_pk` on PK. Grep of `tablecmds.go` does not show explicit "override user PK name" logic — candidate gap.

---

### 1.8 Non-PK index cannot be named "PRIMARY" (ER_WRONG_NAME_FOR_INDEX)

> **new in expansion (Wave 2)** — priority: MED — status: verified

**Setup (error case):**
```sql
CREATE TABLE t (a INT, UNIQUE KEY PRIMARY (a));
-- or
CREATE TABLE t (a INT, INDEX primary (a));
-- or via ALTER
CREATE TABLE t (a INT);
ALTER TABLE t ADD INDEX PRIMARY (a);
```

**Oracle verification:** Expect `ERROR 1280 (42000): Incorrect index name 'PRIMARY'` on each of the three forms. Case-insensitive — `primary`, `Primary`, `PRIMARY` all error.

**omni assertion:** `LoadSDL` / `Apply(ALTER)` should return an error for these; currently likely silent or a different error.

**Anchors:**
- `sql/sql_table.cc:7229-7232` — `if (key->name.str && (key->type != KEYTYPE_PRIMARY) && !my_strcasecmp(system_charset_info, key->name.str, primary_key_name)) { my_error(ER_WRONG_NAME_FOR_INDEX, ...); return true; }`
- `sql/sql_table.cc:15087-15092` — ALTER RENAME INDEX from/to `PRIMARY` check
- doc: [create-table.html](https://dev.mysql.com/doc/refman/8.0/en/create-table.html)

**omni gap:** `indexNameExists` in `tablecmds.go` is case-insensitive but omni does not special-case the string `"PRIMARY"` as reserved. Candidate bug.

---

### 1.9 Implicit index name from first key column

> **new in expansion (Wave 2)** — priority: HIGH — status: pending-verify

**Setup:**
```sql
CREATE TABLE t (
    a INT,
    b INT,
    c INT,
    KEY (b, c)
);
```

**Oracle verification:**
```sql
SELECT INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY INDEX_NAME, SEQ_IN_INDEX;
```
Expected: `INDEX_NAME = 'b'` — the generated name is taken from the *first* key-part column (not all columns concatenated). The second column `c` contributes nothing to the name.

**omni assertion:** `tbl.Indexes[0].Name == "b"`, columns `[b, c]`.

**Anchors:**
- `sql/sql_table.cc:7240-7255` — `else { const Key_part_spec *first_col = key->columns[0]; ... key_info->name = make_unique_key_name(sql_field->field_name, ...); }`
- `mysql/catalog/tablecmds.go:915` — `allocIndexName(tbl, cols[0])` (matches rule)

**Covers omni helper:** matches `allocIndexName`.

---

### 1.10 UNIQUE index backing KEY name fallback when first column equals "PRIMARY" (string literal)

> **new in expansion (Wave 2)** — priority: LOW — status: pending-verify

**Setup:**
```sql
CREATE TABLE t (`PRIMARY` INT);  -- literal column named PRIMARY (backtick-quoted)
-- implicitly, any UNIQUE KEY (`PRIMARY`) would want name "PRIMARY",
-- but make_unique_key_name() explicitly rejects that:
ALTER TABLE t ADD UNIQUE KEY (`PRIMARY`);
```

**Oracle verification:**
```sql
SELECT INDEX_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: INDEX_NAME = `PRIMARY_2` — the first candidate `PRIMARY` is rejected because `make_unique_key_name` does `my_strcasecmp(field_name, primary_key_name)` and falls into the `_2..._99` loop, even though no conflicting index exists.

**omni assertion:** unique index name should be `PRIMARY_2` not `PRIMARY`.

**Anchors:**
- `sql/sql_table.cc:10382-10384` — `if (!check_if_keyname_exists(...) && my_strcasecmp(..., field_name, primary_key_name)) return field_name;` — the `PRIMARY` field_name fails the `my_strcasecmp != 0` guard, skipping the early return.
- `mysql/catalog/tablecmds.go:915-923` — `allocIndexName` does NOT exclude `PRIMARY` — omni gap.

**omni gap:** omni `allocIndexName` would return `"PRIMARY"` here, creating an illegal index name per 1.8.

---

### 1.11 Functional index auto-name "functional_index" with collision suffix

> **new in expansion (Wave 2)** — priority: MED — status: pending-verify

**Setup:**
```sql
CREATE TABLE t (
    a INT,
    INDEX ((a + 1)),            -- no name → functional_index
    INDEX ((a * 2))             -- no name → functional_index_2 (counter starts at 2)
);
```

**Oracle verification:**
```sql
SELECT INDEX_NAME FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY INDEX_NAME;
```
Expected: `[functional_index, functional_index_2]` — note the collision suffix starts at **2**, not 1, per source.

**omni assertion:** indexes named `functional_index` and `functional_index_2`.

**Anchors:**
- `sql/sql_table.cc:7797-7803` — `key_name.assign("functional_index"); ... while (key_name_exists(...)) { key_name.assign("functional_index_"); key_name.append(std::to_string(count++)); }` with `int count = 2;`
- doc: [create-index.html#create-index-functional-key-parts](https://dev.mysql.com/doc/refman/8.0/en/create-index.html)

**omni status 2026-04-24:** verified by `TestScenario_C1/1_11`. Unnamed
functional indexes are assigned `functional_index`, then
`functional_index_2`, matching MySQL's suffix sequence.

---

### 1.12 Functional index hidden generated column name `!hidden!{idx}!{part}!{count}`

> **new in expansion (Wave 2)** — priority: MED — status: verified

**Setup:**
```sql
CREATE TABLE t (a INT, INDEX fx ((a + 1), (a * 2)));
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, EXTRA, GENERATION_EXPRESSION
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
-- Hidden functional columns are not shown by information_schema.COLUMNS by default
-- (they are hidden from SQL layer). Observe via:
SHOW CREATE TABLE t;
-- AND internally via dd schema dumps if needed.
```
Expected internal column names (visible only in the data dictionary, hidden from user queries):
- key part 0 → `!hidden!fx!0!0`
- key part 1 → `!hidden!fx!1!0`

If a user column with a colliding name exists, the trailing `count` increments: `!hidden!fx!0!1`, etc. If the `!hidden!{key}` prefix + `!{part}!{count}` would exceed `NAME_CHAR_LEN` (64), the key_name portion is truncated, NOT the counter (otherwise the loop could never terminate).

**omni assertion:** omni's representation of functional-index hidden columns should match this name format for round-trip deparse.

**Anchors:**
- `sql/sql_table.cc:7710-7743` — `make_functional_index_column_name`:
  ```cpp
  string name("!hidden!"); name += key_name;
  string suffix("!"); suffix += to_string(key_part_number); suffix += '!'; suffix += to_string(count);
  name.resize(min(name.size(), NAME_CHAR_LEN - suffix.size()));
  name.append(suffix);
  ```
- doc: [create-index.html](https://dev.mysql.com/doc/refman/8.0/en/create-index.html) — "Functional indexes are implemented as hidden virtual generated columns"

**omni status 2026-04-24:** verified by `TestScenario_C1/1_12`. Functional
index key parts synthesize system-hidden virtual generated columns named
`!hidden!{idx}!{part}!{count}`.

---

### 1.13 CHECK constraint name is schema-scoped, not table-scoped

> **new in expansion (Wave 2)** — priority: HIGH — status: pending-verify

**Setup (error case):**
```sql
CREATE TABLE t1 (a INT, CONSTRAINT mychk CHECK (a > 0));
CREATE TABLE t2 (a INT, CONSTRAINT mychk CHECK (a < 100));
-- second CREATE TABLE fails
```

**Oracle verification:** Expect `ERROR 3822 (HY000): Duplicate check constraint name 'mychk'.` on the second CREATE. This holds even though the two constraints live on different tables.

Variant — auto-generated name collision across tables:
```sql
CREATE TABLE t1 (a INT, CHECK (a > 0));  -- auto → t1_chk_1
CREATE TABLE t1_chk_1 (a INT, CHECK (a > 0));  -- auto → t1_chk_1_chk_1, no collision
-- But:
CREATE TABLE t_chk_1 (a INT);           -- no check
CREATE TABLE t (a INT, CHECK (a > 0));  -- auto → t_chk_1 — does NOT collide because t_chk_1 is a table name, not a check constraint name
CREATE TABLE u (a INT, CONSTRAINT t_chk_1 CHECK (a > 0));  -- collides with the auto-generated t_chk_1 from table `t` → ERROR 3822
```

**omni assertion:** `LoadSDL` / `Apply(CREATE)` should error on duplicate check constraint names even across tables in the same schema. omni `nextCheckNumber` is per-table today; it does NOT detect cross-table collisions.

**Anchors:**
- `sql/sql_table.cc:19593-19602` — `for (auto cc_name : new_cc_names) { if (thd->dd_client()->check_constraint_exists(*new_schema, cc_name, &exists)) ... if (exists) { my_error(ER_CHECK_CONSTRAINT_DUP_NAME, ...); ... } }` — uses **schema** not table for lookup.
- `sql/sql_table.cc:19046-19066` — `push_check_constraint_mdl_request_to_list` acquires schema-level MDL on lowercased cc_name.
- doc: [create-table-check-constraints.html](https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html) — "CHECK constraint names must be unique per schema; no two tables in the same schema can share a CHECK constraint name."

**omni gap:** `mysql/catalog/tablecmds.go:253-254` and `altercmds.go:467` assign `{tableName}_chk_{N}` without a schema-level uniqueness check. A cross-table duplicate should raise an error — currently it will silently succeed. HIGH priority because this can corrupt SDL diffs when an unrelated table reserves a `_chk_N` name that happens to match.

---

## Section C2: Type normalization

> **Expansion note (Wave 1):** grew from 2 to 24 scenarios via systematic
> walk of MySQL 8.0 type documentation (numeric-type-syntax,
> other-vendor-data-types, date-and-time-type-syntax, string-type-syntax,
> blob, char, data-type-defaults) + targeted grep of
> `sql/sql_yacc.yy` (rules `type`, `int_type`, `real_type`,
> `numeric_type`, `nchar`, `varchar`, `nvarchar`),
> `sql/lex.h` (keyword → token aliases),
> `sql/parse_tree_column_attrs.h` (PT_numeric_type, PT_serial_type,
> PT_bit_type, PT_year_type),
> `sql/create_field.cc` (FLOAT precision split at ~L438),
> `sql/sql_table.cc` (prepare_blob_field at ~L8495, YEAR defaults at
> ~L9118, float/zerofill deprecations at ~L9080-9127).

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

### 2.3 INTEGER → INT

**Setup:**
```sql
CREATE TABLE t (c INTEGER);
```

**Oracle verification:** `SELECT COLUMN_TYPE FROM information_schema.COLUMNS ...`
Expected: `int`.

**omni assertion:** column type renders as `int`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/lex.h:342` —
`{SYM("INTEGER", INT_SYM)}` — INTEGER is a lexer-level alias for INT.
omni handles this in `mysql/parser/type.go:59` (`case kwINT, kwINTEGER`).

**Priority:** LOW — trivial synonym.
**Status:** pending

---

### 2.4 BOOLEAN → TINYINT(1)

**Setup:**
```sql
CREATE TABLE t (c BOOLEAN);
```

**Oracle verification:** Same query.
Expected: `tinyint(1)`.

**omni assertion:** column type renders as `tinyint(1)`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7092-7099`
— both `BOOL_SYM` and `BOOLEAN_SYM` produce `PT_boolean_type`, which maps
to `MYSQL_TYPE_TINY` with length 1.

**Priority:** LOW — spelled-out synonym of BOOL.
**Status:** pending

---

### 2.5 INT1/INT2/INT3/INT4/INT8 → TINYINT/SMALLINT/MEDIUMINT/INT/BIGINT

**Setup:**
```sql
CREATE TABLE t (
  a INT1,
  b INT2,
  c INT3,
  d INT4,
  e INT8
);
```

**Oracle verification:** Query `information_schema.COLUMNS` for each column.
Expected:
- `a` → `tinyint`
- `b` → `smallint`
- `c` → `mediumint`
- `d` → `int`
- `e` → `bigint`

**omni assertion:** each column renders to the normalized form above.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/lex.h:337-341` —
```
{SYM("INT1", TINYINT_SYM)},
{SYM("INT2", SMALLINT_SYM)},
{SYM("INT3", MEDIUMINT_SYM)},
{SYM("INT4", INT_SYM)},
{SYM("INT8", BIGINT_SYM)},
```
omni handles these at `mysql/parser/type.go:319-343`.

**Priority:** LOW — lexer aliases, already handled by omni.
**Status:** pending

---

### 2.6 MIDDLEINT → MEDIUMINT

**Setup:**
```sql
CREATE TABLE t (c MIDDLEINT);
```

**Oracle verification:** Same query.
Expected: `mediumint`.

**omni assertion:** column type renders as `mediumint`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/lex.h:442` —
`{SYM("MIDDLEINT", MEDIUMINT_SYM)}, /* For powerbuilder */`.
omni handles this at `mysql/parser/type.go:344-348`.

**Priority:** LOW — PowerBuilder vendor alias.
**Status:** pending

---

### 2.7 INT(11) display width deprecated → stripped from output

**Setup:**
```sql
CREATE TABLE t (c INT(11));
```

**Oracle verification:** Same query.
Expected: `int` (display width stripped in 8.0.17+; warning ER_WARN_DEPRECATED_INTEGER_DISPLAY_WIDTH emitted).

**omni assertion:** `formatColumnType` emits `int`, not `int(11)`, when
no ZEROFILL. Verified by spot-read of
`/Users/rebeliceyang/Github/omni/mysql/catalog/tablecmds.go:1293-1306`
(isIntType branch strips display width unless zerofill).

**Anchor:**
`/Users/rebeliceyang/Github/mysql-server/sql/parse_tree_column_attrs.h:668-672`
— `PT_numeric_type` pushes `ER_WARN_DEPRECATED_INTEGER_DISPLAY_WIDTH`
warning when length is explicit for an integer type.
Behavior since MySQL 8.0.17 per
https://dev.mysql.com/doc/refman/8.0/en/numeric-type-attributes.html
("As of MySQL 8.0.17, the display width attribute is deprecated for
integer data types; you should expect support for it to be removed in a
future version of MySQL.").

**Priority:** HIGH — schema dump round-trip depends on this.
**Status:** pending

---

### 2.8 INT(N) ZEROFILL → preserves display width

**Setup:**
```sql
CREATE TABLE t (c INT(5) ZEROFILL);
```

**Oracle verification:** Same query.
Expected: `int(5) unsigned zerofill` (ZEROFILL implies UNSIGNED, and
width is preserved when ZEROFILL).

**omni assertion:** `formatColumnType` emits `int(5) unsigned zerofill`.
See `mysql/catalog/tablecmds.go:1297-1305`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7396-7401`
— `ZEROFILL_SYM` sets `ZEROFILL_FLAG` and pushes deprecation warning.
`PT_numeric_type::get_type_flags` at parse_tree_column_attrs.h:675-677
adds `UNSIGNED_FLAG` whenever `ZEROFILL_FLAG` is set. MySQL 8.0
continues to render the display width while ZEROFILL is in effect (even
though it's deprecated).

**Priority:** HIGH — the one case where display width is NOT stripped.
**Status:** pending

---

### 2.9 SERIAL → BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE

**Setup:**
```sql
CREATE TABLE t (c SERIAL);
```

**Oracle verification:**
```sql
SELECT COLUMN_TYPE, IS_NULLABLE, EXTRA
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c';
SHOW INDEX FROM t;
```
Expected:
- `COLUMN_TYPE = 'bigint unsigned'`
- `IS_NULLABLE = 'NO'`
- `EXTRA = 'auto_increment'`
- a UNIQUE index is created on `c`.

**omni assertion:** column renders as `bigint unsigned`, `NOT NULL`,
`AUTO_INCREMENT`, and a UNIQUE index is auto-added. omni handles this at
`mysql/catalog/tablecmds.go:150` (SERIAL expansion) and :303 (implicit
UNIQUE). `formatColumnType` returns `"bigint unsigned"` at :1284.

**Anchor:**
`/Users/rebeliceyang/Github/mysql-server/sql/parse_tree_column_attrs.h:925-932`
— `PT_serial_type` sets
`AUTO_INCREMENT_FLAG | NOT_NULL_FLAG | UNSIGNED_FLAG | UNIQUE_FLAG`
over `MYSQL_TYPE_LONGLONG`. Docs confirm at
https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html
("SERIAL is an alias for BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE.").

**Priority:** HIGH — composite alias that touches four subsystems.
**Status:** pending

---

### 2.10 NUMERIC → DECIMAL

**Setup:**
```sql
CREATE TABLE t (c NUMERIC(10,2));
```

**Oracle verification:** Same query.
Expected: `decimal(10,2)`.

**omni assertion:** column renders as `decimal(10,2)`. Note omni parser
keeps `"NUMERIC"` as the Name (`type.go:114-119`) but `formatColumnType`
lowercases and rewrites to `decimal` at `tablecmds.go:1282-1283`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7325`
— `NUMERIC_SYM { $$= Numeric_type::DECIMAL; }`. Docs:
https://dev.mysql.com/doc/refman/8.0/en/other-vendor-data-types.html
("NUMERIC ... : Synonym for DECIMAL.").

**Priority:** LOW — well-known synonym.
**Status:** pending

---

### 2.11 DEC and FIXED → DECIMAL

**Setup:**
```sql
CREATE TABLE t (a DEC(6,2), b FIXED(6,2));
```

**Oracle verification:** Same query.
Expected: both `decimal(6,2)`.

**omni assertion:** both render as `decimal(6,2)`. omni parser rewrites
to `"DECIMAL"` at `mysql/parser/type.go:124-129`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7324-7326`
— `DECIMAL_SYM`, `NUMERIC_SYM`, and `FIXED_SYM` all map to
`Numeric_type::DECIMAL`. `DEC` is a lexer alias for `DECIMAL_SYM`.

**Priority:** LOW.
**Status:** pending

---

### 2.12 DOUBLE PRECISION → DOUBLE

**Setup:**
```sql
CREATE TABLE t (c DOUBLE PRECISION);
```

**Oracle verification:** Same query.
Expected: `double`.

**omni assertion:** column renders as `double`. omni parser consumes the
optional PRECISION keyword at `mysql/parser/type.go:99-102`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7313-7314`
and `7317-7320` — `DOUBLE_SYM opt_PRECISION` where `opt_PRECISION` is
`%empty | PRECISION`.

**Priority:** LOW — ANSI standard spelling.
**Status:** pending

---

### 2.13 FLOAT4 → FLOAT and FLOAT8 → DOUBLE

**Setup:**
```sql
CREATE TABLE t (a FLOAT4, b FLOAT8);
```

**Oracle verification:** Same query per column.
Expected: `a` → `float`, `b` → `double`.

**omni assertion:** as above. omni handles this at
`mysql/parser/type.go:349-358`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/lex.h:272-273` —
```
{SYM("FLOAT4", FLOAT_SYM)},
{SYM("FLOAT8", DOUBLE_SYM)},
```

**Priority:** LOW — C-language-influenced alias.
**Status:** pending

---

### 2.14 FLOAT(p) precision split — FLOAT(0..24) stays FLOAT, FLOAT(25..53) becomes DOUBLE

**Setup:**
```sql
CREATE TABLE t (
  a FLOAT(10),   -- single-arg precision: "bits of precision"
  b FLOAT(25)
);
```

**Oracle verification:** Same query.
Expected: `a` → `float`, `b` → `double`.

**omni assertion:** `formatColumnType` / parser must rewrite
`FLOAT(25..53)` to `double` and `FLOAT(0..24)` to `float`.
**omni gap:** `mysql/parser/type.go:90-94` keeps Name="FLOAT" and
captures length only; `tablecmds.go:1274+` does not implement the
precision split. Likely mis-normalizes `FLOAT(25)`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/create_field.cc:433-448`
— `Create_field::init` rewrites `MYSQL_TYPE_FLOAT` to `MYSQL_TYPE_DOUBLE`
when `tmp_length > PRECISION_FOR_FLOAT` (24). Docs:
https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html
("MySQL permits a nonstandard syntax: FLOAT(M,D) or REAL(M,D) ... For
FLOAT(p), MySQL treats p as specifying the precision in bits ... MySQL
uses p to determine whether to use FLOAT or DOUBLE for the resulting
data type. If p is from 0 to 24, the data type becomes FLOAT with no M
or D values. If p is from 25 to 53, the data type becomes DOUBLE with
no M or D values.").

**Priority:** HIGH — real omni gap; behavior differs from MySQL.
**Status:** pending-verify

---

### 2.15 FLOAT(M,D) and DOUBLE(M,D) are deprecated non-standard syntax

**Setup:**
```sql
CREATE TABLE t (c FLOAT(7,4));
```

**Oracle verification:** Same query.
Expected: `float(7,4)` (column type still renders the M,D; MySQL emits
warning `ER_WARN_DEPRECATED_FLOAT_DIGITS`).

**omni assertion:** column renders as `float(7,4)`. omni currently
outputs the `(length,scale)` form via `parseOptionalPrecision` +
formatColumnType tail at `tablecmds.go:1319-1320`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/parse_tree_column_attrs.h:649-653`
emits `ER_WARN_DEPRECATED_FLOAT_DIGITS` whenever a float-family type
receives a `dec` (scale). Docs call this a nonstandard MySQL extension.

**Priority:** MED — warning-only, but the `(M,D)` form is still stored
and round-tripped. Wave 5 should verify omni preserves it.
**Status:** pending-verify

---

### 2.16 CHARACTER → CHAR and CHARACTER VARYING → VARCHAR

**Setup:**
```sql
CREATE TABLE t (a CHARACTER(10), b CHARACTER VARYING(20));
```

**Oracle verification:** Same query.
Expected: `a` → `char(10)`, `b` → `varchar(20)`.

**omni assertion:** renders as above. CHARACTER is a lexer alias for
CHAR_SYM; `varchar:` rule in yacc accepts `CHAR_SYM VARYING`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7286-7289`
— `varchar: CHAR_SYM VARYING {} | VARCHAR_SYM {}`. Lexer:
`CHARACTER` is a reserved keyword entry for `CHAR_SYM`. Docs:
https://dev.mysql.com/doc/refman/8.0/en/other-vendor-data-types.html
("CHARACTER: Synonym for CHAR. CHARACTER VARYING: Synonym for VARCHAR.").

**Priority:** LOW.
**Status:** pending

---

### 2.17 NATIONAL CHAR / NCHAR → CHAR with utf8mb3 charset

**Setup:**
```sql
CREATE TABLE t (a NATIONAL CHAR(10), b NCHAR(10));
```

**Oracle verification:**
```sql
SELECT COLUMN_TYPE, CHARACTER_SET_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: both columns `char(10)` with `CHARACTER_SET_NAME = 'utf8mb3'`.

**omni assertion:** both columns render as `char(10)` with charset
`utf8mb3`. omni handles this at `mysql/parser/type.go:282-310`
(hardcodes `"utf8mb3"`).

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7110-7127,7281-7284`
— `nchar` rule with `NCHAR_SYM` or `NATIONAL_SYM CHAR_SYM` maps to
`PT_char_type(Char_type::CHAR, ..., national_charset_info)` plus
`warn_about_deprecated_national`. `national_charset_info` defaults to
`utf8mb3` in 8.0. Docs:
https://dev.mysql.com/doc/refman/8.0/en/charset-national.html.

**omni gap:** hardcoding `utf8mb3` is correct for 8.0 defaults but the
charset is technically configurable. Wave 5 should confirm against the
test container (which uses 8.0).

**Priority:** MED.
**Status:** pending-verify

---

### 2.18 NVARCHAR / NATIONAL VARCHAR / NCHAR VARYING / NCHAR VARCHAR → VARCHAR utf8mb3

**Setup:**
```sql
CREATE TABLE t (
  a NVARCHAR(10),
  b NATIONAL VARCHAR(10),
  c NCHAR VARCHAR(10),
  d NATIONAL CHAR VARYING(10),
  e NCHAR VARYING(10)
);
```

**Oracle verification:** Same.
Expected: all five columns `varchar(10)` with `CHARACTER_SET_NAME='utf8mb3'`.

**omni assertion:** all render as `varchar(10)` charset utf8mb3. omni
handles `NVARCHAR` and `NATIONAL VARCHAR` (type.go:290-317) but
**does NOT appear to cover NCHAR VARCHAR / NCHAR VARYING / NATIONAL
CHAR VARYING**. Confirm by grep — all five forms in the yacc rule must
round-trip.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7291-7297`
— `nvarchar: NATIONAL_SYM VARCHAR_SYM | NVARCHAR_SYM |
NCHAR_SYM VARCHAR_SYM | NATIONAL_SYM CHAR_SYM VARYING |
NCHAR_SYM VARYING`.

**Priority:** MED — obscure but is in the grammar.
**Status:** pending-verify

---

### 2.19 LONG / LONG VARCHAR → MEDIUMTEXT (and LONG VARBINARY → MEDIUMBLOB)

**Setup:**
```sql
CREATE TABLE t (
  a LONG,
  b LONG VARCHAR,
  c LONG VARBINARY
);
```

**Oracle verification:** Same query.
Expected: `a` and `b` → `mediumtext`; `c` → `mediumblob`.

**omni assertion:** rendered as above. omni handles this at
`mysql/parser/type.go:359-372`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7210-7218,7247-7251`
— `LONG_SYM VARBINARY_SYM` → `PT_blob_type(MEDIUM, my_charset_bin)`,
`LONG_SYM varchar ...` → MEDIUM text with given charset, bare
`LONG_SYM opt_charset_with_opt_binary` → MEDIUMTEXT (or MEDIUMBLOB if
binary-forced).

**Priority:** LOW — vendor alias; omni already handles.
**Status:** pending

---

### 2.20 CHAR and BINARY default to length 1 when no length specified

**Setup:**
```sql
CREATE TABLE t (a CHAR, b BINARY);
```

**Oracle verification:** Same.
Expected: `a` → `char(1)`, `b` → `binary(1)`.

**omni assertion:** `formatColumnType` emits `char(1)` / `binary(1)`
when `dt.Length == 0`. See `mysql/catalog/tablecmds.go:1316-1318`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7105-7109,7132-7134`
— the bare `CHAR_SYM opt_charset_with_opt_binary` and `BINARY_SYM`
alternatives construct `PT_char_type` without a length; the length
defaults to 1 in `Create_field::init`. Docs:
https://dev.mysql.com/doc/refman/8.0/en/char.html
("When CHAR values are stored, they are right-padded with spaces to the
specified length. ... If M is omitted, the length is 1.").

**Priority:** MED — omni handles it, Wave 5 should still verify.
**Status:** pending

---

### 2.21 VARCHAR without length is a syntax error

**Setup:**
```sql
CREATE TABLE t (c VARCHAR);
```

**Oracle verification:** Expected: `ERROR 1064 (42000)` near the
column definition (syntax error — VARCHAR requires a length).

**omni assertion:** parse error. Confirm omni rejects this with a
parser error; the bare `VARCHAR_SYM` alone is NOT in the `type` rule.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7136-7140`
— `varchar field_length opt_charset_with_opt_binary` is the only
production; `varchar` has no production without `field_length`.

**Priority:** MED — parser strictness check.
**Status:** pending-verify

---

### 2.22 TIMESTAMP/DATETIME/TIME default fractional-seconds precision is 0

**Setup:**
```sql
CREATE TABLE t (
  a TIMESTAMP,       -- fsp=0
  b DATETIME,        -- fsp=0
  c TIME,            -- fsp=0
  d TIMESTAMP(6),
  e DATETIME(6),
  f TIME(3)
);
```

**Oracle verification:** Query `DATETIME_PRECISION` from
`information_schema.COLUMNS`.
Expected:
- `a`/`b`/`c` → `DATETIME_PRECISION = 0` (stored as `timestamp`,
  `datetime`, `time` — no `(N)` suffix)
- `d`/`e` → `timestamp(6)`, `datetime(6)`
- `f` → `time(3)`

**omni assertion:** same. omni captures the precision via
`parseOptionalLength` at `mysql/parser/type.go:171-193`; `formatColumnType`
renders the `(N)` only when Length > 0.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7181-7192,7365-7368`
— `type_datetime_precision: %empty { $$= NULL; } | '(' NUM ')'`. When
fsp is NULL, `PT_time_type` / `PT_timestamp_type` default to 0. Docs:
https://dev.mysql.com/doc/refman/8.0/en/date-and-time-type-syntax.html
("An fsp value of 0 signifies that there is no fractional part. If
omitted, the default precision is 0.").

**Priority:** MED — round-trip asymmetry concern (`datetime` vs
`datetime(0)`).
**Status:** pending

---

### 2.23 YEAR(4) syntax deprecated → stored as bare YEAR

**Setup:**
```sql
CREATE TABLE t (c YEAR(4));
```

**Oracle verification:** Same query.
Expected: `year` (no display width in 8.0+). A deprecation warning
`YEAR(4)` is emitted at parse time.

**omni assertion:** `formatColumnType` elides the length for YEAR —
see `mysql/catalog/tablecmds.go:1314-1315`.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy:7154-7176`
— `YEAR_SYM opt_field_length field_options` rule: if length is
non-null it must be 4 (else `ER_INVALID_YEAR_COLUMN_LENGTH`), then
`push_deprecated_warn("YEAR(4)", "YEAR")`. Length is then ignored
(`PT_year_type` takes no length). See catalog C20.1.

**Priority:** MED.
**Status:** pending

---

### 2.24 BIT without length defaults to BIT(1)

**Setup:**
```sql
CREATE TABLE t (c BIT);
```

**Oracle verification:** Same.
Expected: `bit(1)`.

**omni assertion:** column renders as `bit(1)`. **Possible omni gap:**
`mysql/parser/type.go:229-232` captures length via `parseOptionalLength`
but does NOT set a default of 1 when absent; `tablecmds.go:1274+` does
not special-case `bit` the way it does `char`/`binary`. Likely renders
as `bit` (missing `(1)`).

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/parse_tree_column_attrs.h:687-695`
— `PT_bit_type() : PT_type(MYSQL_TYPE_BIT), length("1") {}` (default
ctor sets length to "1"). Grammar at sql_yacc.yy:7084-7091.

**Priority:** HIGH — round-trip omni gap.
**Status:** pending-verify

---

### 2.25 VARCHAR(N > 65535 bytes) auto-converts to TEXT family (non-strict mode)

**Setup:** (session in non-strict SQL mode, no `constant_default`)
```sql
SET sql_mode='';
CREATE TABLE t (c VARCHAR(65536));
```

**Oracle verification:** Same query.
Expected: `c` becomes `mediumtext` (or `text`, depending on exact
byte count), with a `Note 1246: Converting column 'c' from VARCHAR to
TEXT` warning.

In strict mode the same DDL fails with
`ER_TOO_BIG_FIELDLENGTH`.

**omni assertion:** omni should either error (strict) or auto-convert
(non-strict). **omni gap:** no byte-length → TEXT-family conversion in
`mysql/catalog/tablecmds.go` — verify the current behavior.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_table.cc:8495-8521`
— `prepare_blob_field` rewrites `sql_type` via `get_blob_type_from_length`
when `max_display_width_in_bytes() > MAX_FIELD_VARCHARLENGTH` and the
field has no BLOB_FLAG; pushes `ER_AUTO_CONVERT` note. Errors in
strict mode.

**Priority:** MED — edge case but a real silent-rewrite.
**Status:** pending-verify

---

### 2.26 TEXT(N) / BLOB(N) → TINYTEXT / TEXT / MEDIUMTEXT / LONGTEXT by byte count

**Setup:**
```sql
CREATE TABLE t (
  a TEXT(100),     -- ≤255 bytes → tinytext
  b TEXT(1000),    -- ≤65535 → text
  c TEXT(70000),   -- ≤16777215 → mediumtext
  d TEXT(20000000) -- ≤4294967295 → longtext
);
```

**Oracle verification:** Same query.
Expected:
- `a` → `tinytext`
- `b` → `text`
- `c` → `mediumtext`
- `d` → `longtext`

**omni assertion:** renders as above. **omni gap:** `formatColumnType`
at `mysql/catalog/tablecmds.go:1310-1313` simply *strips* the length
from `text` / `blob`, leaving the column as `text` regardless of N.
Does not promote to tinytext/mediumtext/longtext.

**Anchor:** `/Users/rebeliceyang/Github/mysql-server/sql/sql_table.cc:8546-8565`
— `prepare_blob_field` calls `get_blob_type_from_length` to pick
TINY/MEDIUM/LONG per explicit display width. Docs:
https://dev.mysql.com/doc/refman/8.0/en/blob.html
("If you specify a length for TEXT(M), MySQL chooses the smallest TEXT
type large enough to hold values M characters long. The same applies
to BLOB(M).").

**Priority:** HIGH — real omni gap that affects schema round-trip.
**Status:** pending-verify

---

## Section C3: Nullability & default promotion

> **Expansion note (Wave 4 batch A):** grew from 3 to 8 scenarios via walk of
> `sql/sql_table.cc` (`promote_first_timestamp_column` L4221-L4245,
> `mysql_prepare_create_table` NOT_NULL_FLAG propagation for PK/AUTO_INC,
> gcol nullability at L7886-L7895) and the CREATE TABLE / TIMESTAMP doc pages.

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
- `ts1` -> `TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`
- `ts2` -> `TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00'` (NOT promoted)

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

### 3.4 Explicit NULL on PRIMARY KEY column is a hard error
**Priority:** HIGH
**Status:** verified
**Setup:**
```sql
CREATE TABLE t (a INT NULL PRIMARY KEY);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "A PRIMARY KEY column ... must not contain any NULL values."
**Source:** `sql/sql_table.cc` — `mysql_prepare_create_table` sets `NOT_NULL_FLAG` when column participates in PK; if `EXPLICIT_NULL_FLAG` also set, raises `ER_PRIMARY_CANT_HAVE_NULL` (`sql/share/messages_to_clients.txt`).
**Oracle:** MySQL returns error 1171 (`ER_PRIMARY_CANT_HAVE_NULL`).
**omni pointer:** `mysql/catalog/tablecmds.go` create-path must reject explicit NULL + PK combination; `mysql/catalog/altercmds.go` `alterColumnSetNotNull`. Check if omni currently silently coerces.
**omni assertion:** `Exec` returns error; catalog unchanged.
**omni status 2026-04-24:** verified by `TestScenario_C3/3_4`.

---

### 3.5 UNIQUE KEY does NOT imply NOT NULL
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT UNIQUE);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "a UNIQUE index permits multiple NULL values for columns that can contain NULL."
**Source:** `sql/sql_table.cc:3904-3932` — UNIQUE NOT NULL only ordered before other UNIQUE keys when NOT NULL flag comes from the column definition, not from the UNIQUE index itself.
**Oracle:** `IS_NULLABLE = 'YES'` for column `a`; `information_schema.TABLE_CONSTRAINTS` shows UNIQUE; row `(NULL), (NULL)` both succeed.
**omni pointer:** `mysql/catalog/tablecmds.go` resolveColumnNullable — verify that only PK / AUTO_INC promotion flips `col.Nullable`, not UNIQUE.
**omni assertion:** `tbl.GetColumn("a").Nullable == true` and UNIQUE index present.

---

### 3.6 Generated column nullability is derived from expression, not declared NOT NULL
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (
    a INT NULL,
    b INT GENERATED ALWAYS AS (a+1) VIRTUAL
);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "A generated column is implicitly nullable if the expression can be NULL; an explicit NOT NULL makes it an error to store NULL."
**Source:** `sql/sql_table.cc` `prepare_create_field` — for a gcol, NOT_NULL_FLAG follows the expression's `maybe_null`; writing a user NOT NULL on top does not change the inferred nullability until insert time.
**Oracle:** `IS_NULLABLE = 'YES'` for `b`; insert `(NULL)` into `a` yields `b = NULL`.
**omni pointer:** `mysql/catalog/table.go:64` Column.Nullable — gcol path in `tablecmds.go` should not silently force NOT NULL from the presence of an expression; must honor the underlying column nullability.
**omni assertion:** `tbl.GetColumn("b").Nullable == true`.

---

### 3.7 AUTO_INCREMENT on an explicitly NULL column silently promotes to NOT NULL
**Priority:** MED
**Status:** verified
**Setup:**
```sql
CREATE TABLE t (id INT NULL AUTO_INCREMENT, KEY(id));
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/example-auto-increment.html — "an AUTO_INCREMENT column must not contain NULL values."
**Oracle:** MySQL 8.0.45 accepts the DDL and promotes `id` to `NOT NULL`;
the earlier hard-error expectation was incorrect.
**omni status 2026-04-24:** verified by `TestScenario_C3/3_7`; omni matches
the silent-promotion behavior.

---

### 3.8 Second TIMESTAMP column without DEFAULT is implicitly NOT NULL DEFAULT '0000-00-00 00:00:00' under explicit_defaults_for_timestamp=OFF
**Priority:** LOW
**Status:** pending
**Setup:**
```sql
SET SESSION explicit_defaults_for_timestamp=OFF;
CREATE TABLE t (ts1 TIMESTAMP, ts2 TIMESTAMP);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/timestamp-initialization.html — deprecated legacy path when `explicit_defaults_for_timestamp=OFF`.
**Source:** `sql/sql_table.cc:4221-4245` `promote_first_timestamp_column`: first TIMESTAMP gets DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP; subsequent TIMESTAMPs are forced NOT NULL DEFAULT zero even without explicit NOT NULL.
**Oracle (legacy mode):** `ts1` NOT NULL CURRENT_TIMESTAMP; `ts2` NOT NULL DEFAULT '0000-00-00 00:00:00'. Note: default 8.0 mode is ON — scenario is LOW because omni targets default session settings.
**omni pointer:** `mysql/catalog/catalog.go` session var tracking; `tablecmds.go` timestamp defaults. Mostly a known-limitation to document, not fix.
**omni assertion:** Skip in default mode; if omni adds session var support, match the legacy transform.

---


## Section C4: Charset / collation inheritance

Scope: this section covers MySQL 8.0's implicit charset/collation resolution
chain (server → database → table → column → expression), the BINARY /
NCHAR / NATIONAL special forms, literal introducers, index-prefix byte
accounting, and the precedence rules DTCollation uses to negotiate charsets
between operands. Every scenario below was chosen because it produces a
materially different `SHOW CREATE TABLE` / `information_schema.COLUMNS` /
`information_schema.TABLES` shape that omni must reproduce on round-trip.

Primary sources:
- `sql/sql_table.cc:8448` `set_table_default_charset()` — table-level default
- `sql/sql_table.cc:4188` `get_sql_field_charset()` — column-level default
- `sql/sql_table.cc:4544-4560` `prepare_create_field` BINCMP_FLAG handling
- `sql/sql_yacc.yy:489`, `sql/mysqld.cc:10208` `national_charset_info`
- `sql/item.h:160-241` `DTCollation` derivation levels
- Refman: charset-database / charset-table / charset-column /
  charset-national / charset-binary-op / charset-collate / charset-literal

---

### 4.1 Table charset inherits from database

**Priority:** HIGH
**Status:** existing (Wave 0)

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

**Priority:** HIGH
**Status:** existing (Wave 0)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10)) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**Oracle verification:**
```sql
SELECT COLLATION_NAME FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c';
SHOW CREATE TABLE t;
```
Expected: column collation = `utf8mb4_0900_ai_ci`; `SHOW CREATE TABLE` elides
the column-level CHARACTER SET / COLLATE (matches table default), per C18.1.

**omni assertion:** column charset == table charset, and deparse output does
not print a column-level CHARACTER SET clause.

**See catalog:** C4.2 + C18.1.

---

### 4.3 Column `COLLATE` alone derives `CHARACTER SET` from collation name

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10) COLLATE utf8mb4_unicode_ci)
  DEFAULT CHARSET=latin1;
```

The column was not given a `CHARACTER SET`, only a `COLLATE`. MySQL derives
the charset from the collation (`utf8mb4_unicode_ci` → `utf8mb4`), so the
column is `utf8mb4` even though the table default is `latin1`.

**Oracle verification:**
```sql
SELECT CHARACTER_SET_NAME, COLLATION_NAME
FROM information_schema.COLUMNS
WHERE TABLE_NAME='t' AND COLUMN_NAME='c';
SHOW CREATE TABLE t;
```
Expected: `CHARACTER_SET_NAME='utf8mb4'`, `COLLATION_NAME='utf8mb4_unicode_ci'`.
`SHOW CREATE TABLE` emits `CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci`
on the column (both because charset differs from table default).

**omni assertion:** `col.Charset == "utf8mb4"`, `col.Collation ==
"utf8mb4_unicode_ci"`, round-trip deparse preserves the column clause.

**Source anchor:** parser attaches `COLLATE` to `Create_field`; charset is
looked up from the collation row during `get_sql_field_charset` because
`sql_field->charset` was populated from the collation side of the parse node
(`sql/sql_yacc.yy` charset/collation attribute rules; follows the rule in
refman charset-collate.html "if only COLLATE, charset is that of collation").

**See catalog:** C4 (new subsection), C18.1.

---

### 4.4 Column `CHARACTER SET` alone derives default collation

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10) CHARACTER SET latin1)
  DEFAULT CHARSET=utf8mb4;
```

No column COLLATE given. MySQL picks the charset's **default collation**
(`latin1` → `latin1_swedish_ci`). This is the mirror of 4.3.

**Oracle verification:**
```sql
SELECT CHARACTER_SET_NAME, COLLATION_NAME
FROM information_schema.COLUMNS
WHERE TABLE_NAME='t' AND COLUMN_NAME='c';
SHOW CREATE TABLE t;
```
Expected: `COLLATION_NAME='latin1_swedish_ci'`. `SHOW CREATE TABLE` emits
`CHARACTER SET latin1 COLLATE latin1_swedish_ci` on the column (column charset
differs from table default → cannot be elided, C18.1).

**omni assertion:** `col.Charset == "latin1"`, `col.Collation ==
"latin1_swedish_ci"`; deparse reproduces the clause.

**Source anchor:** `CHARSET_INFO::default_collation` lookup in
`get_charset_by_csname(...)`; parser fills `sql_field->charset` with the
charset's primary collation row when only `CHARACTER SET` is supplied.

**See catalog:** C4 (new subsection).

---

### 4.5 Table `COLLATE` must be compatible with table `CHARACTER SET`

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup (error case):**
```sql
CREATE TABLE t (c VARCHAR(10))
  CHARACTER SET latin1 COLLATE utf8mb4_0900_ai_ci;
```

**Setup (success case):**
```sql
CREATE TABLE t (c VARCHAR(10))
  CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

If the COLLATE does not belong to the named CHARACTER SET, MySQL raises
`ER_COLLATION_CHARSET_MISMATCH` ("COLLATION 'utf8mb4_0900_ai_ci' is not valid
for CHARACTER SET 'latin1'"). Validation runs before `set_table_default_charset`
and is performed by the parser's `merge_charset_and_collation` helper /
`check_charset_collation_pair`.

**Oracle verification:**
- Mismatch case: statement fails with errno 1253 (SQLSTATE 42000).
- Compatible case: table created; `TABLE_COLLATION='utf8mb4_0900_ai_ci'`.

**omni assertion:** mismatch is rejected at parse/analyse time with a mirrored
error, not silently coerced. Compatible combination resolves both fields.

**Source anchor:** `sql/sql_parse.cc` / `sql/parse_tree_nodes.cc`
`merge_charset_and_collation()` (pre-flight compatibility check before
`set_table_default_charset` is called in `create_table_impl`).

**See catalog:** C4 (new subsection).

---

### 4.6 `BINARY` modifier on CHAR/VARCHAR/TEXT resolves to `{charset}_bin` collation

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (
  a CHAR(10) BINARY,
  b VARCHAR(10) CHARACTER SET latin1 BINARY
) DEFAULT CHARSET=utf8mb4;
```

`BINARY` here is an **attribute**, not a column type. `prepare_create_field`
(sql_table.cc:4546-4560) sees `BINCMP_FLAG` and rewrites the charset to the
same charset's `_bin` collation via
`get_charset_by_csname(cs->csname, MY_CS_BINSORT, MYF(0))`. Column `a` becomes
`utf8mb4 / utf8mb4_bin`; column `b` becomes `latin1 / latin1_bin`.

MySQL 8.0 also deprecates this form (`warn_about_deprecated_binary`, sql_yacc.yy:
"BINARY as attribute of a type … use a CHARACTER SET clause with _bin
collation"), so round-trip through `SHOW CREATE TABLE` re-emits the **canonical
form**: `CHARACTER SET xxx COLLATE xxx_bin`. The `BINARY` keyword does **not**
come back.

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
SELECT COLUMN_NAME, COLLATION_NAME FROM information_schema.COLUMNS
WHERE TABLE_NAME='t';
```
Expected:
- `a` → `utf8mb4_bin` (rendered as `COLLATE utf8mb4_bin` only — charset matches
  table default so charset clause is elided).
- `b` → `latin1_bin` (rendered as `CHARACTER SET latin1 COLLATE latin1_bin` —
  charset differs from table default).
- Neither column has the `BINARY` keyword in the reprinted DDL.

**omni assertion:** ingesting the BINARY-attribute form produces the same
resolved `col.Collation` values as ingesting the canonical `COLLATE xxx_bin`
form; deparse always emits the canonical form (matches MySQL round-trip).

**Source anchor:** `sql/sql_table.cc:4546` BINCMP_FLAG branch; `sql/sql_yacc.yy`
deprecation warning (`warn_about_deprecated_binary`).

**See catalog:** C4 (new subsection), C18.1.

---

### 4.7 `CHARACTER SET binary` is not the same as column type `BINARY`

**Priority:** MEDIUM
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (
  a BINARY(10),                    -- byte type, sql_type = MYSQL_TYPE_STRING + BINARY_FLAG
  b CHAR(10) CHARACTER SET binary, -- text type stored as bytes
  c VARBINARY(10)                  -- MYSQL_TYPE_VARCHAR + binary charset
);
```

All three end up with `COLLATION_NAME='binary'`, but they are different column
types and different `DATA_TYPE` values:
- `a` → `DATA_TYPE='binary'`
- `b` → `DATA_TYPE='char'` with `CHARACTER SET binary`
- `c` → `DATA_TYPE='varbinary'`

`get_sql_field_charset` hard-codes an early return for `&my_charset_bin` and
for array columns (sql_table.cc:4203) so that `ALTER TABLE … CONVERT TO
CHARACTER SET …` does not touch them. This is why a converted table keeps
`BINARY`/`VARBINARY` columns as bytes.

**Oracle verification:**
```sql
SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME
FROM information_schema.COLUMNS WHERE TABLE_NAME='t';
SHOW CREATE TABLE t;
```
Expected: three distinct `DATA_TYPE` / `COLUMN_TYPE` values; all three have
`COLLATION_NAME='binary'`. `SHOW CREATE TABLE` emits `binary(10)` for `a`,
`char(10) CHARACTER SET 'binary'` (or equivalent) for `b`, `varbinary(10)` for
`c`. `ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4` leaves `a` and `c`
untouched but rewrites `b` to `char(10) CHARACTER SET utf8mb4`.

**omni assertion:** the three columns parse to three different
`DataType` values (not conflated); `Collation == "binary"` for all; CONVERT TO
CHARACTER SET semantics preserve `BINARY`/`VARBINARY` on round-trip.

**Source anchor:** `sql/sql_table.cc:4203` (`cs == &my_charset_bin` early
return); refman charset-binary-op.html.

**See catalog:** C4 (new subsection).

---

### 4.8 `utf8` is an alias for `utf8mb3` (not the server default `utf8mb4`)

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10) CHARACTER SET utf8);
```

Even though the MySQL 8.0 server default is `utf8mb4`, the **literal name**
`utf8` in DDL is an alias for `utf8mb3`. The parser resolves it that way and
the resulting column is `utf8mb3 / utf8mb3_general_ci` (or whatever the
session's default collation for utf8mb3 is).

Additionally, `SHOW CREATE TABLE` in 8.0.30+ normalises the stored clause back
to `utf8mb3` — the `utf8` spelling is **not** preserved on round-trip. omni's
catalog must normalise ingested `utf8` to `utf8mb3` so that a subsequent
deparse matches MySQL.

**Oracle verification:**
```sql
SELECT CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
WHERE TABLE_NAME='t' AND COLUMN_NAME='c';
SHOW CREATE TABLE t;
```
Expected: `CHARACTER_SET_NAME='utf8mb3'`, `COLLATION_NAME='utf8mb3_general_ci'`.
`SHOW CREATE TABLE` prints `utf8mb3`, never `utf8`.

**omni assertion:** after ingestion, `col.Charset == "utf8mb3"` (not `utf8`);
deparse emits `utf8mb3`. Also applies to table-level and database-level
charset declarations.

**Source anchor:** refman
charset-collation-implementations.html; `mysys/charset.cc` aliases;
`sql/mysqld.cc:10208` national_charset_info = `my_charset_utf8mb3_general_ci`
(proves utf8 = utf8mb3 internally).

**See catalog:** C4 (new subsection).

---

### 4.9 `NATIONAL CHAR` / `NCHAR` is hard-wired to `utf8mb3` (deprecated in 8.0)

**Priority:** MEDIUM
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (
  a NCHAR(10),
  b NATIONAL CHARACTER(10),
  c NATIONAL VARCHAR(10),
  d NCHAR VARYING(10)
);
```

These column forms ignore the table/server default and use
`national_charset_info`, which `sql/mysqld.cc:10208` hard-codes to
`&my_charset_utf8mb3_general_ci`. This is SQL-2003 legacy behaviour: NATIONAL
means "the server's designated Unicode charset", and MySQL still uses
utf8mb3 for that slot. Because utf8mb3 is deprecated, MySQL 8.0 issues
`ER_DEPRECATED_NATIONAL` (sql_yacc.yy:489) and the documentation recommends
switching to explicit `CHARACTER SET utf8mb4`.

**Oracle verification:**
```sql
SELECT COLUMN_NAME, CHARACTER_SET_NAME, COLLATION_NAME
FROM information_schema.COLUMNS WHERE TABLE_NAME='t';
SHOW CREATE TABLE t;
```
Expected: all four columns `CHARACTER_SET_NAME='utf8mb3'`,
`COLLATION_NAME='utf8mb3_general_ci'`. `SHOW CREATE TABLE` rewrites every form
to the same canonical shape (typically `char(10) CHARACTER SET utf8mb3` /
`varchar(10) CHARACTER SET utf8mb3`), the NATIONAL / NCHAR keyword is **not**
preserved.

**omni assertion:** parser accepts all four forms, resolves each to
`utf8mb3 / utf8mb3_general_ci`, and deparse emits the canonical
`CHARACTER SET utf8mb3` shape (matches MySQL round-trip). Ingesting the
deparser output and re-deparsing yields an identical string.

**Source anchor:** `sql/sql_yacc.yy:7113-7144` (`national_charset_info` used for
NCHAR/NATIONAL productions); `sql/sql_yacc.yy:487-492`
`warn_about_deprecated_national`; `sql/mysqld.cc:10208`.

**See catalog:** C4 (new subsection), C17.

---

### 4.10 `ENUM` / `SET` charset inheritance follows the same column rule

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (
  a ENUM('x','y'),
  b ENUM('x','y') CHARACTER SET latin1,
  c SET('p','q') COLLATE utf8mb4_unicode_ci
) DEFAULT CHARSET=utf8mb4;
```

ENUM and SET are string columns, so they go through `get_sql_field_charset`
exactly like CHAR/VARCHAR. The critical side effect is that MySQL **converts
the enum literals** from the client charset to the column charset during
`prepare_create_field` (sql_table.cc:4567-4591,
`constant_default->safe_charset_converter`). After creation, the stored
definition lists the literals in the column charset's encoding, which is what
`information_schema.COLUMNS.COLUMN_TYPE` reflects.

**Oracle verification:**
```sql
SELECT COLUMN_NAME, COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME
FROM information_schema.COLUMNS WHERE TABLE_NAME='t';
SHOW CREATE TABLE t;
```
Expected:
- `a` → charset utf8mb4 (table default), `COLUMN_TYPE="enum('x','y')"`.
- `b` → charset latin1, default collation latin1_swedish_ci.
- `c` → charset utf8mb4 (derived from COLLATE, see 4.3),
  collation utf8mb4_unicode_ci.

**omni assertion:** enum/set elements stored on the column carry the resolved
charset; catalog reports charset on ENUM/SET columns identically to
VARCHAR; deparse round-trips values without mojibake.

**Source anchor:** `sql/sql_table.cc:4593-4596` (`prepare_set_field`,
`prepare_enum_field`) which use `sql_field->charset` set two lines earlier.

**See catalog:** C4 (new subsection).

---

### 4.11 Index prefix length is measured in **bytes**, so charset change multiplies it

**Priority:** HIGH
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t1 (c VARCHAR(10) CHARACTER SET latin1, KEY k (c(10)));
CREATE TABLE t2 (c VARCHAR(10) CHARACTER SET utf8mb4, KEY k (c(10)));
CREATE TABLE t3 (c VARCHAR(200) CHARACTER SET utf8mb4, KEY k (c(768)));
```

The `c(N)` prefix length is **characters** in the DDL but is stored and
reported in **bytes internally**. MySQL multiplies by
`charset->mbmaxlen` (sql_table.cc:5022:
`column->get_prefix_length() * sql_field->charset->mbmaxlen`), so:
- t1 `c(10)` → 10 × 1 = 10 bytes on disk.
- t2 `c(10)` → 10 × 4 = 40 bytes on disk.
- t3 `c(768)` → 768 × 4 = 3072 bytes, **exceeds** InnoDB's default 3072-byte
  per-column key limit → rejected with `ER_TOO_LONG_KEY` unless the caller
  has `innodb_large_prefix` / DYNAMIC row format.

Additionally, when a key part length is rounded, the rounding is done on byte
boundaries: `key_part_length -= key_part_length % sql_field->charset->mbmaxlen`
(sql_table.cc:5183, 5241), ensuring the stored length is a multiple of
`mbmaxlen`.

**Oracle verification:**
```sql
SELECT INDEX_NAME, COLUMN_NAME, SUB_PART
FROM information_schema.STATISTICS
WHERE TABLE_NAME IN ('t1','t2');
SHOW CREATE TABLE t1; SHOW CREATE TABLE t2;
```
Expected: `SUB_PART` reported as 10 for both in characters; but
`innodb_sys_indexes` / internal max length differs by 4×. The t3 CREATE fails
with error 1071 "Specified key was too long".

**omni assertion:** index prefix validation uses the resolved column charset's
mbmaxlen; omni must reject the t3 form with the same error. Deparse preserves
the user-written prefix length in characters (not bytes).

**Source anchor:** `sql/sql_table.cc:5022` (prefix × mbmaxlen), 5183, 5241
(rounding).

**See catalog:** C4 (new subsection), cross-link to C6/C16 (index
defaults).

---

### 4.12 Expression collation negotiation: column COLLATE vs `COLLATE` expression vs literal introducer

**Priority:** MEDIUM
**Status:** new (Wave 2)

**Setup:**
```sql
CREATE TABLE t (c VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci);
-- query time:
SELECT c = 'abc' FROM t;                              -- implicit vs coercible
SELECT c = _utf8mb4'abc' COLLATE utf8mb4_bin FROM t;  -- coercible vs explicit
SELECT c COLLATE utf8mb4_bin = _latin1'abc' FROM t;   -- explicit vs coercible
SELECT CAST('abc' AS CHAR CHARACTER SET latin1) = c FROM t; -- explicit vs implicit
```

MySQL's `DTCollation` class (sql/item.h:160-241) classifies every expression
with a derivation level:

| Level                 | Value | Source                                       |
|-----------------------|-------|----------------------------------------------|
| DERIVATION_EXPLICIT   | 0     | `COLLATE x` clause in the expression         |
| DERIVATION_NONE       | 1     | result of failed aggregation                 |
| DERIVATION_SYSCONST   | 2     | system constants (USER(), VERSION())         |
| DERIVATION_IMPLICIT   | 3     | column reference                             |
| DERIVATION_COERCIBLE  | 4     | string literal                               |
| DERIVATION_NUMERIC    | 5     | numeric / binary                             |

When aggregating two operands (`agg_item_charsets_for_string_result`), the
one with the **lower (stronger)** derivation wins. If both are EXPLICIT and
their collations differ → `ER_CANT_AGGREGATE_2COLLATIONS`. Column vs literal:
column wins (IMPLICIT < COERCIBLE). Two columns with different collations:
error unless one is EXPLICIT. Literal introducers (`_utf8mb4'abc'`) do **not**
raise derivation — the literal stays COERCIBLE; to force it, add a
`COLLATE` clause, which raises it to EXPLICIT.

**Oracle verification:** each query above executes with the documented
outcome:
- `c = 'abc'` uses the column's collation (implicit wins over coercible).
- `c = _utf8mb4'abc' COLLATE utf8mb4_bin` uses utf8mb4_bin (explicit wins).
- `c COLLATE utf8mb4_bin = _latin1'abc'` coerces latin1 to utf8mb4_bin
  (explicit wins; cross-charset conversion allowed because derivation differs).
- `CAST('abc' AS CHAR CHARACTER SET latin1) = c` → `ER_CANT_AGGREGATE_2COLLATIONS`
  because CAST produces an implicit latin1 vs implicit utf8mb4 column, and
  neither is explicit.

**omni assertion:** omni's expression-type resolver must track derivation on
every string-producing node (column ref, literal, CAST, COLLATE) and replay
the same aggregation outcome. In particular, the catalog must **not** reject
cross-charset comparisons that MySQL allows, and must reject the ones MySQL
rejects — both are needed for correct CHECK-constraint / generated-column
validation and for SDL-diff stability on computed-column definitions.

**Source anchor:** `sql/item.h:160-241` DTCollation class and derivation enum;
`sql/item.cc` `DTCollation::aggregate`; refman charset-collation-coercibility.html.

**See catalog:** C4 (new subsection); cross-link to C18 (generated columns).

---

## Section C5: Constraint defaults

> **Expansion note (Wave 4 batch A):** grew from 3 to 10 scenarios via walk of
> `sql/sql_table.cc` FK preparation (`prepare_fk_parent_key` L6340-L6500,
> FK action validation L6672-L6710, gcol FK rejection L6783/L6883/L9464) and
> CHECK constraint preparation (`prepare_check_constraints_for_create` L19007-L19031,
> `is_enforced` handling L19180/L19379/L19443/L19512/L19727).
> Does not duplicate Wave 1 C22 (ALTER algorithm) or Wave 2 C1 (name auto-gen).

### 5.1 FK ON DELETE default -> stored as RESTRICT, reported as NO ACTION, elided in SHOW

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

### 5.4 FK ON UPDATE default is independent of ON DELETE
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id) ON DELETE CASCADE);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html — "If ON DELETE or ON UPDATE are not specified, the default action is NO ACTION."
**Source:** `sql/sql_yacc.yy` `opt_on_update_delete` — DELETE and UPDATE options are parsed into independent fields; absence of one does not inherit from the other. `sql/sql_table.cc:6682-6701`.
**Oracle:** `REFERENTIAL_CONSTRAINTS.UPDATE_RULE = 'NO ACTION'` while `DELETE_RULE = 'CASCADE'`. `SHOW CREATE TABLE` renders `ON DELETE CASCADE` only, no `ON UPDATE` clause.
**omni pointer:** `mysql/catalog/constraint.go` FK struct — verify OnUpdate is not defaulted from OnDelete in `tablecmds.go` `buildFK`. `mysql/catalog/show.go:443` rendering path.
**omni assertion:** `fk.OnDelete == "CASCADE"`, `fk.OnUpdate == "" / RESTRICT`; deparse omits ON UPDATE.

---

### 5.5 FK SET DEFAULT is parsed but rejected at runtime (InnoDB limitation)
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (
    a INT DEFAULT 0,
    FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT
);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html — "InnoDB parses but rejects SET DEFAULT; it is not supported."
**Source:** `sql/sql_table.cc:6697-6701` — `FK_OPTION_DEFAULT` is explicitly listed among referential actions that conflict with CHECK constraints. InnoDB handler emits `ER_FK_NO_INDEX_CHILD` / not-supported at handler level.
**Oracle:** CREATE TABLE succeeds in MySQL 8.0.16+ parser but InnoDB returns an error like `ER_NOT_SUPPORTED_YET` ("ON DELETE SET DEFAULT") or, depending on version, accepts the syntax silently. Exact expected-behavior capture is in spot-check phase.
**omni pointer:** `mysql/catalog/tablecmds.go` FK options handler — must either reject SET DEFAULT or record it with a documented limitation flag to match InnoDB.
**omni assertion:** omni should reject SET DEFAULT on CREATE to match InnoDB behavior.

---

### 5.6 FK column type must match referenced column type (size/sign)
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE p (id BIGINT PRIMARY KEY);
CREATE TABLE c (a INT, FOREIGN KEY (a) REFERENCES p(id));
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-foreign-keys.html — "Corresponding columns in the foreign key and the referenced key must have similar data types. The size and sign of integer types must be the same."
**Source:** `sql/sql_table.cc` `prepare_fk_parent_key` and InnoDB handler `dict_foreign_find_col` — size/sign mismatch raises `ER_FK_INCOMPATIBLE_COLUMNS`.
**Oracle:** MySQL error 3780 (`ER_FK_INCOMPATIBLE_COLUMNS`).
**omni pointer:** `mysql/catalog/tablecmds.go` FK resolution — must compare column type including size and unsigned flag.
**omni assertion:** `Exec` returns error for mismatched FK types.

---

### 5.7 FK on a VIRTUAL generated column is rejected (CREATE)
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE p (id INT PRIMARY KEY);
CREATE TABLE c (
    a INT,
    b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
    FOREIGN KEY (b) REFERENCES p(id)
);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "Foreign keys referencing a virtual generated column are not permitted."
**Source:** `sql/sql_table.cc:6672` / `6783` / `6883` / `9464` all raise `ER_FK_CANNOT_USE_VIRTUAL_COLUMN`.
**Oracle:** MySQL error 3104. (Note: STORED gcol is also rejected as FK **child**; allowed only as FK **parent** in limited cases — capture separately as spot-check.)
**omni pointer:** `mysql/catalog/tablecmds.go` FK build — must check the referencing column is not a virtual gcol. Overlaps with C9.2 (kept here for the "which FK actions also conflict" enumeration).
**omni assertion:** `Exec` returns error; cross-ref to 9.4.

---

### 5.8 CHECK NOT ENFORCED flag defaults to ENFORCED (is_enforced = true)
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (
    a INT,
    CONSTRAINT chk_pos CHECK (a > 0)
);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html — "If ENFORCED/NOT ENFORCED is omitted, the default is ENFORCED."
**Source:** `sql/sql_table.cc:19180` — `cc_spec->is_enforced = table_cc.is_enforced()`; parser default is true.
**Oracle:** `information_schema.CHECK_CONSTRAINTS` row for `chk_pos` exists; `information_schema.TABLE_CONSTRAINTS.ENFORCED = 'YES'`. `SHOW CREATE TABLE` does not render `NOT ENFORCED`.
**omni pointer:** `mysql/catalog/constraint.go:25` `NotEnforced bool` — zero value (false) matches MySQL default. `mysql/catalog/show.go:443` — only renders `NOT ENFORCED` when flag is true. Already looks correct; scenario is a regression fence.
**omni assertion:** `con.NotEnforced == false`; deparse omits the clause.

---

### 5.9 CHECK constraint at column level vs table level is equivalent in information_schema
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t1 (a INT CHECK (a > 0));
CREATE TABLE t2 (a INT, CHECK (a > 0));
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html — "A column-level CHECK constraint can reference only that column; otherwise it is stored identically to a table-level CHECK."
**Source:** `sql/sql_table.cc:19007-19031` `prepare_check_constraints_for_create` — merges column-level and table-level CHECK lists and assigns auto names uniformly.
**Oracle:** Both tables expose CHECK constraints in `information_schema.CHECK_CONSTRAINTS`; auto-names follow same `tblname_chk_N` scheme. Expression stored post-normalization; both render as `CHECK ((`a` > 0))` in `SHOW CREATE TABLE`.
**omni pointer:** `mysql/catalog/tablecmds.go` `buildCheckConstraints` / `collectColumnChecks` — ensure column-level CHECKs are merged into the same constraint list before numbering.
**omni assertion:** Both tables expose a single CHECK with identical serialized expression.

---

### 5.10 Column-level CHECK referencing another column is rejected (column-scope constraint)
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT CHECK (a > b), b INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html — "A column-level CHECK constraint must refer only to the column on which it is defined. Referring to other columns produces ER_CHECK_CONSTRAINT_REFERS."
**Source:** `sql/sql_table.cc` `prepare_check_constraints_for_create` — column-scoped CHECK has its referenced columns validated; mismatch raises `ER_CHECK_CONSTRAINT_REFERS` (3823).
**Oracle:** Error 3823.
**omni pointer:** `mysql/catalog/tablecmds.go` `buildCheckConstraints` must distinguish column-level vs table-level scope when validating referenced identifiers.
**omni assertion:** `Exec` returns error. (Note: MySQL allows table-level CHECK to reference any table column — confirm this path still succeeds.)

---


## Section C6: Partition defaults

Wave 1 worker output for category **C6 — Partition defaults**. Existing scenarios 6.1–6.3 are preserved (verify-only notes). New scenarios start at 6.4.

Primary docs: https://dev.mysql.com/doc/refman/8.0/en/partitioning.html and sub-pages.
Primary source: `sql/partition_info.cc`, `sql/sql_partition.cc`.
Scope note: NDB-specific partition behaviors are deliberately excluded — omni targets InnoDB.

---

### 6.1 PARTITION BY HASH without PARTITIONS clause defaults to 1
**Priority:** HIGH
**Status:** pending-verify (already in SCENARIOS; re-run oracle to confirm InnoDB returns 1, not engine-specific default)
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-hash.html — "If PARTITIONS num is not specified, the number of partitions defaults to 1."
**Source:** `sql/partition_info.cc:683-693` (`set_up_default_partitions`).
**omni pointer:** `mysql/catalog/tablecmds.go:705` — `buildPartitionInfo` fills `pi.Partitions` only when `NumParts > 0`; need to verify default-to-1 fallback when `PARTITIONS` clause absent.

---

### 6.2 Subpartitions default to 1 if not specified
**Priority:** MED
**Status:** pending-verify
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-subpartitions.html
**Source:** `sql/partition_info.cc:750-754` (`set_up_default_subpartitions`).
**omni pointer:** `mysql/catalog/tablecmds.go:717-720`.

---

### 6.3 Partition ENGINE defaults to table ENGINE
**Priority:** MED
**Status:** pending-verify
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html (partition_options)
**Source:** `sql/partition_info.cc:706` — `part_elem->engine_type = default_engine_type` in `set_up_default_partitions`; `sql/partition_info.cc:1494-1516` propagates table engine to per-partition elements during `check_partition_info`.
**omni pointer:** `mysql/catalog/table.go:47` — `PartitionDefInfo.Engine`; `mysql/catalog/show.go:603` render path.

---

### 6.4 KEY partitioning ALGORITHM defaults to 2 on CREATE/ALTER
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-key.html — "The default algorithm used by PARTITION BY KEY is 2 (default in MySQL 5.6 and later)."
**Source:** `sql/partition_info.cc:2427-2451` (`fix_parser_data`): for `PARTITION BY KEY` (i.e. `part_type == HASH && list_of_part_fields`) with `key_algorithm == KEY_ALGORITHM_NONE`, CREATE/ALTER forces `key_algorithm = KEY_ALGORITHM_55` (= 2). Same block lines 2441-2451 applies to KEY subpartitioning.
**Trigger:** `CREATE TABLE t (id INT) PARTITION BY KEY(id) PARTITIONS 4;`
**Rule:** SHOW CREATE TABLE must emit `PARTITION BY KEY (id) /*!50100 PARTITIONS 4 */` — no explicit `ALGORITHM=` — because algorithm 2 is the default and not printed; but when the user wrote `ALGORITHM=1`, SHOW CREATE TABLE must print `KEY ALGORITHM=1 (id)`.
**Observable via:** `SHOW CREATE TABLE t`.
**omni pointer:** `mysql/catalog/table.go:32` — `Algorithm int`; `mysql/catalog/show.go:567` unconditionally prints `KEY ALGORITHM = %d` when `Algorithm != 0`. **GAP:** omni never applies a default of 2, and because zero-valued `Algorithm` skips the branch, omni outputs `KEY (col)` instead of `KEY (col)`; for users who explicitly wrote `ALGORITHM=2` the value round-trips. Need to confirm omni parser does not pre-fill `Algorithm=2` for unset cases (would be wrong for SHOW CREATE) and that omni catalog's internal default for query-planning still matches `2`.

---

### 6.5 KEY partitioning with empty column list defaults to the PRIMARY KEY columns
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-key.html — "It is possible to use PARTITION BY KEY() (with an empty list) ... MySQL uses the table's primary key, if there is one, as the partitioning key."
**Source:** `sql/sql_partition.cc:763-793` (`set_up_field_array`): when `is_list_empty && part_type == HASH` (true for `PARTITION BY KEY()`), iterates `table->key_info[primary_key].key_part` and flags each PK column as a partition field. If no PK exists and engine lacks `HA_USE_AUTO_PARTITION`, fails with `ER_FIELD_NOT_FOUND_PART_ERROR` (InnoDB does not set the flag, so error fires).
**Trigger:** `CREATE TABLE t (id INT PRIMARY KEY, v INT) PARTITION BY KEY() PARTITIONS 4;`
**Rule:**
- If table has a PK → partition columns = PK columns.
- If no PK → error (for InnoDB).
- BLOB/TEXT columns in PK → `ER_BLOB_FIELD_IN_PART_FUNC_ERROR` at partition-field setup time.
**Observable via:** `SHOW CREATE TABLE` (renders `PARTITION BY KEY (id)` with the PK column name).
**omni pointer:** `mysql/catalog/tablecmds.go:645-665` — builds `pi.Columns` directly from AST column list; does not inject PK columns when list empty. **GAP:** omni likely renders `PARTITION BY KEY ()` (or an empty col list), and does not error when no PK is present.

---

### 6.6 LINEAR HASH / LINEAR KEY only changes the placement algorithm, not the storage layout
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-linear-hash.html — "LINEAR HASH uses a power-of-two algorithm, whereas HASH uses modulus. This has no effect on SHOW CREATE TABLE storage defaults — `LINEAR` is retained verbatim."
**Source:** partition expression parser + `partition_info::linear_hash_ind` field; SHOW CREATE uses `PARTITION BY LINEAR HASH`/`LINEAR KEY` when `linear_hash_ind` set. No default fallback — user-specified only.
**Rule:** `LINEAR` token is preserved, defaults to off; power-of-two placement does not change partition count defaults.
**omni pointer:** `mysql/catalog/table.go:29` — `Linear bool`; `show.go` renders. Spot-check: omni preserves `LINEAR` on round-trip.

---

### 6.7 RANGE / LIST partitioning require explicit partition definitions (no `PARTITIONS n` shortcut)
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-range.html and .../partitioning-list.html
**Source:** `sql/partition_info.cc:673-679` — `set_up_default_partitions` errors with `ER_PARTITIONS_MUST_BE_DEFINED_ERROR` for any `part_type` other than HASH. So `PARTITION BY RANGE (id) PARTITIONS 4;` (no `(PARTITION p0 VALUES ..., ...)`) is rejected.
**Trigger:** `CREATE TABLE t (id INT) PARTITION BY RANGE (id) PARTITIONS 4;`
**Rule:** Error `1708 (HY000): For RANGE partitions each partition must be defined`. Same for LIST. Only HASH/KEY may rely on implicit definition-generation from `PARTITIONS n`.
**omni pointer:** `mysql/catalog/tablecmds.go:705-715` — omni's fallback that materializes `p0..p{n-1}` when `Partitions == nil && NumParts > 0` will do so unconditionally; **GAP:** should reject for RANGE/LIST.

---

### 6.8 MAXVALUE must appear in the last RANGE partition only
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-range.html
**Source:** `sql/partition_info.cc:2245-2259` (`fix_partition_values`): `col_val->max_value` is only accepted when `part_id == num_parts - 1`; duplicate MAXVALUE or misplaced MAXVALUE → `ER_PARTITION_MAXVALUE_ERROR`. Also `NULL` in `VALUES LESS THAN` → `ER_NULL_IN_VALUES_LESS_THAN` (line 2275, 2468).
**Rule:**
- `VALUES LESS THAN MAXVALUE` is only legal on the final range partition.
- `VALUES LESS THAN (NULL)` is illegal.
**omni pointer:** `mysql/catalog/tablecmds.go:620-699` — check if `buildPartitionInfo` tracks `MAXVALUE` position / rejects multiple. **Potential GAP.**

---

### 6.9 LIST partitioning comparison is equality (`=`), RANGE is strict less-than (`<`)
**Priority:** LOW (docs-only; affects semantic analysis, not DDL shape)
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-list.html — "VALUES IN (list) matches column value by equality" vs RANGE "VALUES LESS THAN (value) is strict less-than".
**Source:** partition pruning in `sql/partition_info.cc:2268` (`INT_RESULT` check) and later pruning code; not a CREATE default, but the **semantic** difference is often missed by catalogs.
**Rule:** LIST rows whose partition expression does not match any `IN (...)` value → `ER_NO_PARTITION_FOR_GIVEN_VALUE` at INSERT time (unless a `DEFAULT` partition is defined — see 6.10). RANGE values equal to `LESS THAN (v)` go to the **next** partition.
**omni pointer:** n/a for DDL catalog, but flag so row-routing semantics aren't confused in any future partition pruning work.

---

### 6.10 LIST DEFAULT partition acts as catch-all (MySQL 8.0+)
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-list.html — "From MySQL 8.0.3, it is possible to specify a DEFAULT partition for a LIST or LIST COLUMNS partition. Rows that are not matched ... are stored in the DEFAULT partition."
**Source:** parser grammar `sql_yacc.yy` (`opt_part_values`, `DEFAULT_SYM`), `partition_element::has_default_value`; populated at parse time, printed by `sql_partition.cc` in SHOW CREATE.
**Trigger:** `PARTITION BY LIST (c) (PARTITION p0 VALUES IN (1,2), PARTITION pd DEFAULT)`
**Rule:** At most one `DEFAULT` partition per LIST/LIST COLUMNS table. RANGE does not support `DEFAULT` (use `MAXVALUE`). SHOW CREATE preserves the keyword verbatim.
**omni pointer:** `mysql/catalog/tablecmds.go:732` — `partitionValueToString`; check if `DEFAULT` token is represented and round-trips. **GAP suspected** — existing `PartitionDefInfo.ValueExpr` is a string, may not encode `DEFAULT`.

---

### 6.11 Partition function result must be INTEGER (non-COLUMNS variants)
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-types.html — "The partitioning function must return an integer value" (except for RANGE COLUMNS / LIST COLUMNS which accept typed column lists directly).
**Source:** `sql/partition_info.cc:2268-2272` — `item_expr->result_type() != INT_RESULT` → `ER_VALUES_IS_NOT_INT_TYPE_ERROR`. Upstream: `check_partition_func_processor` (`Item::check_partition_func_processor`) at 1452/1459 rejects non-deterministic and non-integer-producing functions.
**Rule:** RANGE/LIST/HASH expression must evaluate to INT. Use `UNIX_TIMESTAMP(ts_col)` to partition by TIMESTAMP; use `TO_DAYS(d)`/`YEAR(d)` for DATE.
**omni pointer:** `mysql/catalog/tablecmds.go:620-645` — omni stores `pi.Expr` as a raw string; does not validate integer-result. **GAP:** omni silently accepts `PARTITION BY RANGE (CONCAT(a,b))` or `PARTITION BY RANGE (ts_col)` where MySQL errors.

---

### 6.12 TIMESTAMP columns require UNIX_TIMESTAMP() wrapping in RANGE/LIST partitioning
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-limitations.html — explicitly: "The only legal function that can be used with a TIMESTAMP column is UNIX_TIMESTAMP()."
**Source:** `Item_func_unix_timestamp::check_partition_func_processor` + companion DATE/TIME function whitelist in `sql/item_timefunc.cc`. Any other wrapping → `ER_PARTITION_FUNCTION_IS_NOT_ALLOWED`.
**Rule:** `PARTITION BY RANGE (ts_col)` → error; `PARTITION BY RANGE (UNIX_TIMESTAMP(ts_col))` → ok.
**omni pointer:** same as 6.11; analyzer does not enforce.

---

### 6.13 Every UNIQUE KEY (and PRIMARY KEY) must contain all partition-expression columns
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-limitations-partitioning-keys-unique-keys.html — "every unique key on the table must use every column in the table's partitioning expression."
**Source:** `sql/sql_partition.cc:1006-1042` (`check_primary_key`) and `sql/sql_partition.cc:1044+` (`check_unique_keys`) — both emit `ER_UNIQUE_KEY_NEED_ALL_FIELDS_IN_PF`.
**Trigger:** `CREATE TABLE t (a INT, b INT, UNIQUE KEY (a)) PARTITION BY HASH (b) PARTITIONS 4;` → errors.
**Rule:** Applies to implicit uniqueness (PRIMARY KEY, UNIQUE). Non-unique indexes are unaffected. A unique key may be a **superset** of partition columns.
**omni pointer:** `mysql/catalog/tablecmds.go` partition build path — omni does NOT run this validation on the constraint set. **GAP — likely bug**; matches the "silently accepts schemas MySQL rejects" pattern called out in catalog PS7.

---

### 6.14 Per-partition options (DATA/INDEX DIRECTORY, MAX_ROWS, MIN_ROWS, TABLESPACE, COMMENT) are preserved verbatim and NOT inherited from the table
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — partition_options grammar.
**Source:** `sql/partition_info.cc:1402-1410` — `part_elem->data_file_name` / `index_file_name` stored per partition; `sql/partition_info.cc:2599` `partition_name, tablespace_name, data_file_name, index_file_name` list; `sql/partition_info.cc:2671-2688` copy semantics. No silent inheritance — an unset per-partition option is `nullptr`, not the table-level value.
**Rule:** `SHOW CREATE TABLE` only emits the per-partition option when it was user-specified; and the table-level equivalent is independent.
**omni pointer:** `mysql/catalog/table.go:44-50` — `PartitionDefInfo` only tracks `Name`, `ValueExpr`, `Engine`, `Comment`. **GAP:** omni drops `DATA DIRECTORY`, `INDEX DIRECTORY`, `MAX_ROWS`, `MIN_ROWS`, `TABLESPACE`, `NODEGROUP` on round-trip.

---

### 6.15 Subpartition options inherit from parent partition, then from table
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-subpartitions.html
**Source:** `sql/partition_info.cc:767` — `new (*THR_MALLOC) partition_element(part_elem)` — subpartition element **copy-constructs from parent partition**, so engine_type, tablespace, and data_file_name are inherited unless overridden in the SUBPARTITION clause. Line 773 then explicitly re-applies `default_engine_type` to engine.
**Rule:** Subpartition engine = subpartition override ?? table engine (NOT parent partition's explicit engine). Data/index dir and max_rows inherit from parent partition at copy time, then overridden.
**omni pointer:** `mysql/catalog/table.go:53` — `SubPartitionDefInfo` has only `Name, Engine, Comment`. Same coverage gap as 6.14.

---

### 6.16 ALTER TABLE ADD PARTITION auto-naming seeds from current partition count (not max-suffix+1)
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table-partition-operations.html
**Source:** `sql/sql_partition.cc:4506` — `alt_part_info->set_up_defaults_for_partitioning(part_handler, nullptr, tab_part_info->num_parts)`; seed is the count. Also documented in catalog section PS6.
**Rule:** Only legal for HASH/KEY (RANGE/LIST require explicit names). Auto-generated names are `p{num_parts}..p{num_parts + new - 1}`.
**omni pointer:** omni currently has no ALTER TABLE ADD PARTITION support per `mysql/catalog/altercmds.go` grep in catalog PS6. Flag for future.

---

### 6.17 COALESCE PARTITION removes the last N HASH partitions; REORGANIZE renames freely
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table-partition-operations.html
**Source:** `sql/sql_partition.cc:4738-4756` — `ALTER_COALESCE_PARTITION` only valid for HASH (`ER_COALESCE_ONLY_ON_HASH_PARTITION`); always drops from the **tail** of the partition list, reducing `num_parts` by the coalesce count. Lines 4827+ handle REORGANIZE which rebuilds partition definitions from the user-provided list.
**Rule:** `ALTER TABLE t COALESCE PARTITION 2` on a HASH table with `p0..p5` leaves `p0..p3`. Data in `p4, p5` is redistributed across survivors.
**omni pointer:** not implemented; track as planned ALTER-partition work. TRUNCATE PARTITION (not covered here) preserves all definitions and only clears rows.

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

### 7.3 BTREE is the default algorithm for InnoDB (HASH requires explicit USING)
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-index.html, https://dev.mysql.com/doc/refman/8.0/en/innodb-index-types.html
**Source:** `sql/sql_table.cc:7341` — `key_info->is_algorithm_explicit = false` is the default; line 7358-7393 — when `is_algorithm_explicit` is false the server assigns `file->get_default_index_algorithm()`, which for InnoDB returns `HA_KEY_ALG_BTREE`. If the user wrote nothing, `key->key_create_info.algorithm == HA_KEY_ALG_SE_SPECIFIC` (asserted line 7392). HASH on InnoDB/MyISAM is silently coerced to BTREE (`ER_UNSUPPORTED_INDEX_ALGORITHM` warning, line 7377-7381).
**Rule:**
```sql
CREATE TABLE t (a INT, KEY (a) USING HASH) ENGINE=InnoDB;
```
InnoDB does not support HASH secondary indexes; the engine coerces to BTREE and emits `Note 3502: The storage engine for the table doesn't support HASH`. `SHOW CREATE TABLE` shows `KEY a (a)` with no `USING` clause — the explicit-flag is discarded.
**Oracle verification:** `SELECT INDEX_TYPE FROM information_schema.STATISTICS WHERE TABLE_NAME='t'` returns `BTREE`; `SHOW WARNINGS` shows note 3502.
**omni assertion:** parsing `USING HASH` on InnoDB should leave the catalog index with `IndexType="BTREE"` (post-engine-normalization) and record a warning; currently omni keeps the user string verbatim.
**omni pointer:** `mysql/catalog/indexcmds.go:62` stores `stmt.IndexType` raw. Needs engine-aware normalization hook.

---

### 7.4 USING BTREE explicit vs implicit differs in SHOW CREATE output
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-index.html
**Source:** `sql/sql_table.cc:7341,7367` — `is_algorithm_explicit` is a first-class bit on `KEY_INFO`. `sql/sql_show.cc` (store_create_info) emits `USING BTREE` only when `key_info->is_algorithm_explicit` is true. `sql/sql_table.cc:12340` — ALTER path diffs this flag and triggers `CHANGE_INDEX_OPTION` even with no physical change.
**Rule:**
```sql
CREATE TABLE t1 (a INT, KEY (a));              -- SHOW CREATE: KEY `a` (`a`)
CREATE TABLE t2 (a INT, KEY (a) USING BTREE);  -- SHOW CREATE: KEY `a` (`a`) USING BTREE
```
Both index the same way; only the rendering differs. Round-trip deparsing must preserve the explicit-flag to stay byte-identical with the oracle.
**Oracle verification:** `SHOW CREATE TABLE t1\G` vs `SHOW CREATE TABLE t2\G` — compare the `KEY` line.
**omni assertion:** catalog needs a `IndexTypeExplicit bool` companion field (or nil vs `"BTREE"` distinction). Deparser must only print `USING BTREE` when explicit.
**omni pointer:** `mysql/catalog/indexcmds.go:62` — currently conflates default and explicit BTREE. Deparser in `mysql/deparse/` drops or emits unconditionally — verify.

---

### 7.5 UNIQUE index on a nullable column permits multiple NULLs
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-index.html (“A UNIQUE index permits multiple NULL values for columns that can contain NULL.”)
**Source:** `sql/sql_table.cc:5134-5143` — when a key part’s column is nullable, the server sets `HA_NULL_PART_KEY` on `key_info->flags`. The UNIQUE uniqueness check in InnoDB (`storage/innobase/row/row0ins.cc`) short-circuits on NULL key parts, treating each NULL tuple as distinct. Contrast with PRIMARY KEY path (5130-5137) which forcibly sets `NOT_NULL_FLAG` and decrements `null_bits`.
**Rule:**
```sql
CREATE TABLE t (a INT, UNIQUE KEY (a));
INSERT INTO t VALUES (NULL), (NULL), (NULL);   -- all succeed
INSERT INTO t VALUES (1), (1);                 -- second fails: duplicate key
```
**Oracle verification:** both inserts succeed; `SELECT COUNT(*) FROM t` returns 3 after the NULL inserts; second `(1)` insert returns `ER_DUP_ENTRY 1062`.
**omni assertion:** catalog must retain column nullability after UNIQUE index creation (no forced NOT NULL, unlike PK). Advisor/diff should not treat UNIQUE as implying NOT NULL.
**omni pointer:** `mysql/catalog/tablecmds.go` PK path promotes to NOT NULL; verify UNIQUE path does NOT.

---

### 7.6 Indexes default to VISIBLE; INVISIBLE must be explicit (and PK cannot be invisible)
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/invisible-indexes.html
**Source:** `sql/sql_table.cc:7342` — `key_info->is_visible = key->key_create_info.is_visible` where the parser default for `KEY_CREATE_INFO` is `is_visible=true`. Lines 7473-7474 — `if (!k->is_visible) my_error(ER_PK_INDEX_CANT_BE_INVISIBLE)`; line 15153 — ALTER INDEX VISIBLE/INVISIBLE repeats the same PK guard. Column hidden-ness (`HT_VISIBLE`, line 8357) is a separate DD concept from index visibility.
**Rule:**
```sql
CREATE TABLE t (a INT, KEY ix (a));                 -- VISIBLE (default)
CREATE TABLE t (a INT, KEY ix (a) INVISIBLE);       -- invisible, optimizer ignores
CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE);  -- ER_PK_INDEX_CANT_BE_INVISIBLE 3522
```
**Oracle verification:** `SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE INDEX_NAME='ix'` returns `YES`/`NO`.
**omni assertion:** `idx.Visible` defaults to true; parser must record the explicit `INVISIBLE` keyword; PK with INVISIBLE must be rejected at analyze time.
**omni pointer:** `mysql/catalog/indexcmds.go:51` sets `Visible: true`. Verify PRIMARY KEY path errors when INVISIBLE is present.

---

### 7.7 BLOB/TEXT columns require an explicit prefix length in an index
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-index.html (Column Prefix Key Parts), https://dev.mysql.com/doc/refman/8.0/en/innodb-limits.html
**Source:** `sql/sql_table.cc:5076-5082` — `if (is_blob(sql_field->sql_type)) { if (!column_length) { my_error(ER_BLOB_KEY_WITHOUT_LENGTH); return true; } }`. Line 5157-5194 — the prefix length is clamped to the engine max key part length (InnoDB: 767 bytes for COMPACT/REDUNDANT, 3072 bytes for DYNAMIC/COMPRESSED — row-format dependent). FULLTEXT indexes (line 5151 branch) are exempt: the whole column is tokenized.
**Rule:**
```sql
CREATE TABLE t (a TEXT, KEY (a));             -- ER_BLOB_KEY_WITHOUT_LENGTH 1170
CREATE TABLE t (a TEXT, KEY (a(100)));         -- OK
CREATE TABLE t (a TEXT, FULLTEXT KEY (a));     -- OK, no prefix required
```
**Oracle verification:** first form returns 1170; second yields `SUB_PART=100` in `information_schema.STATISTICS`; third yields `INDEX_TYPE=FULLTEXT` and `SUB_PART=NULL`.
**omni assertion:** analyze-time error on non-FULLTEXT BLOB/TEXT index without length; carry `SubPart` into catalog.
**omni pointer:** check `mysql/catalog/indexcmds.go` key-part handling for prefix length propagation.

---

### 7.8 FULLTEXT index uses the built-in parser when WITH PARSER is omitted
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/innodb-fulltext-index.html, https://dev.mysql.com/doc/refman/8.0/en/fulltext-natural-language.html
**Source:** `sql/sql_table.cc:7307-7308` — `if (key->key_create_info.parser_name.str) key_info->parser_name = key->key_create_info.parser_name;` — only set when the user wrote `WITH PARSER`. Otherwise `key_info->parser_name` is `{nullptr, 0}` and InnoDB falls back to its built-in whitespace/stopword parser (not ngram; ngram must be explicit `WITH PARSER ngram`). Line 15136 — SHOW CREATE re-emits `WITH PARSER <name>` only when `key_info->parser` is non-null.
**Rule:**
```sql
CREATE TABLE t (a TEXT, FULLTEXT KEY (a));                     -- default parser
CREATE TABLE t (a TEXT, FULLTEXT KEY (a) WITH PARSER ngram);   -- ngram explicit
```
`SHOW CREATE TABLE` for the first form does **not** emit any `WITH PARSER` clause.
**Oracle verification:** compare `SHOW CREATE TABLE` outputs; query `information_schema.INNODB_FT_INDEX_TABLE` after inserting multi-byte text to confirm tokenization differs.
**omni assertion:** catalog `ParserName` is nil (not `"default"`) when omitted; deparser suppresses the clause in that case.
**omni pointer:** `mysql/catalog/indexcmds.go` should parse optional `WITH PARSER` and leave unset when absent.

---

### 7.9 SPATIAL index requires all key columns to be NOT NULL
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/creating-spatial-indexes.html, https://dev.mysql.com/doc/refman/8.0/en/create-index.html
**Source:** `sql/sql_table.cc:5144-5147` — `if (key->type == KEYTYPE_SPATIAL || sql_field->sql_type == MYSQL_TYPE_GEOMETRY) my_error(ER_SPATIAL_CANT_HAVE_NULL)`. Lines 4998-5007 — an explicit `USING BTREE/HASH` on a SPATIAL key raises `ER_INDEX_TYPE_NOT_SUPPORTED_FOR_SPATIAL_INDEX`; line 7352-7354 — the algorithm is force-set to `HA_KEY_ALG_RTREE` and `is_algorithm_explicit` must be false. Related: SRID on the column is required for the optimizer to use the spatial index (not an error, but a usability rule).
**Rule:**
```sql
CREATE TABLE t (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g));  -- OK
CREATE TABLE t (g GEOMETRY, SPATIAL KEY (g));                     -- ER_SPATIAL_CANT_HAVE_NULL 1252
CREATE TABLE t (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE);-- ER_INDEX_TYPE_NOT_SUPPORTED_FOR_SPATIAL_INDEX 3500
```
**Oracle verification:** error codes 1252 and 3500 respectively.
**omni assertion:** analyze-time checks for (a) non-null geometry columns under SPATIAL, (b) rejection of `USING BTREE|HASH` on SPATIAL.
**omni pointer:** `mysql/catalog/indexcmds.go:58` sets `idx.Spatial=true` and `IndexType="SPATIAL"` but does not validate column nullability or algorithm conflicts.

---

### 7.10 PRIMARY KEY and UNIQUE KEY on the same columns both persist
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html (Keys section), https://dev.mysql.com/doc/refman/8.0/en/create-index.html
**Source:** `sql/sql_table.cc` — key-list processing in `mysql_prepare_create_table` enumerates each `Key_spec` in `alter_info->key_list` and appends a `KEY_INFO` per spec; there is no dedup step that collapses a UNIQUE onto an existing PRIMARY (only the opposite is forbidden — a second PRIMARY raises `ER_MULTIPLE_PRI_KEY`). The two indexes share tuple storage under InnoDB clustering rules but appear as distinct rows in `information_schema.STATISTICS`.
**Rule:**
```sql
CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY uk (a));
```
`SHOW CREATE TABLE` preserves both:
```
PRIMARY KEY (`a`),
UNIQUE KEY `uk` (`a`)
```
**Oracle verification:** `SELECT INDEX_NAME FROM information_schema.STATISTICS WHERE TABLE_NAME='t'` returns two rows: `PRIMARY` and `uk`.
**omni assertion:** catalog keeps both indexes distinct; diff engine must not merge/drop the redundant UNIQUE as a “same-columns” dedup.
**omni pointer:** verify `mysql/catalog/tablecmds.go` appends both keys and diff code does not collapse.

---

## Section C8: Table option defaults

> **Expansion note (Wave 4 batch A):** grew from 3 to 10 scenarios via walk of
> `sql/sql_table.cc` (create_info option processing), `sql/handler.cc` default
> populators, and the CREATE TABLE / table_options doc section.
> Category focuses on **effective values** (what omni stores); see C18 for
> `SHOW CREATE TABLE` elision rules — a single rule can appear in both.

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

### 8.4 CHARSET defaults to database default when table omits CHARACTER SET
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE DATABASE d1 CHARACTER SET latin1;
USE d1;
CREATE TABLE t (a VARCHAR(10));
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/charset-table.html — "If the CHARACTER SET option is not given, the database character set is used."
**Source:** `sql/sql_table.cc` `mysql_prepare_create_table` — `create_info->default_table_charset` filled from database default when `HA_CREATE_USED_DEFAULT_CHARSET` is not set.
**Oracle:** `information_schema.TABLES.TABLE_COLLATION = 'latin1_swedish_ci'`; column `a` is `varchar(10) CHARACTER SET latin1`.
**omni pointer:** `mysql/catalog/tablecmds.go` `resolveTableCharset` — verify fallback to `db.Charset`; cross-ref C4.1 which covers the **inheritance chain** from server, this focuses on the **effective table-level value**.
**omni assertion:** `tbl.Charset == "latin1"`.

---

### 8.5 COLLATE independence from CHARSET — specifying COLLATE alone fills CHARSET
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT) COLLATE=latin1_german2_ci;
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/charset-table.html — "If COLLATE is given but CHARACTER SET is not, the character set is derived from the collation's base character set."
**Source:** `sql/sql_table.cc` table option collapse — `create_info->table_charset` derived from `create_info->default_table_charset->csname` after resolving the named collation.
**Oracle:** `TABLE_COLLATION = 'latin1_german2_ci'`; derived CHARSET = latin1.
**omni pointer:** `mysql/catalog/tablecmds.go` — table option path should invoke charset-from-collation resolution (mirrors C4.3 for columns).
**omni assertion:** `tbl.Charset == "latin1"`, `tbl.Collation == "latin1_german2_ci"`.

---

### 8.6 KEY_BLOCK_SIZE defaults to 0 (engine picks) and is elided in SHOW
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html table_options: "KEY_BLOCK_SIZE = value ... A value of 0 means use the default." https://dev.mysql.com/doc/refman/8.0/en/innodb-file-format.html — default for COMPRESSED format comes from `innodb_page_size/2`.
**Source:** `sql/handler.cc` HA_CREATE_INFO default — `key_block_size` initialized to 0.
**Oracle:** `information_schema.TABLES.CREATE_OPTIONS` has no `key_block_size`. `SHOW CREATE TABLE` omits clause. Cross-ref C18.
**omni pointer:** `mysql/catalog/table.go:17` `KeyBlockSize int` — zero-value is the default; verify `mysql/catalog/show.go` does not render when 0.
**omni assertion:** `tbl.KeyBlockSize == 0`; deparse omits clause.

---

### 8.7 COMPRESSION defaults to "" (None) and is elided in SHOW
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "COMPRESSION [=] {'ZLIB'|'LZ4'|'NONE'}. The default is None."
**Source:** `sql/handler.cc` HA_CREATE_INFO — `compress` is LEX_CSTRING empty. `sql/sql_table.cc:10990` mentions "Do not keep ENCRYPTION clause for unencrypted table" — parallel elision rule.
**Oracle:** `SHOW CREATE TABLE` omits `COMPRESSION=`; `information_schema.TABLES.CREATE_OPTIONS` is empty string.
**omni pointer:** omni Table struct (`mysql/catalog/table.go:3-24`) has **no** `Compression` field — this is an **omni gap**. Either add field or document as unsupported.
**omni assertion:** CREATE succeeds; compression option parsed and stored (or explicitly documented as dropped).

---

### 8.8 ENCRYPTION defaults depend on server var `default_table_encryption`
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "If ENCRYPTION is not specified, the default depends on the default_table_encryption system variable."
**Source:** `sql/sql_table.cc:9032-9220` ENCRYPTION validation + default application; `sql/sql_table.cc:10990` elides ENCRYPTION from SHOW when unencrypted.
**Oracle:** With default `default_table_encryption=OFF`, `CREATE_OPTIONS = ''`, `SHOW CREATE TABLE` omits ENCRYPTION.
**omni pointer:** omni Table has no `Encryption` field — **omni gap**. Document that encryption is not modeled (out of scope for file-format simulation).
**omni assertion:** CREATE succeeds; option silently dropped or flagged.

---

### 8.9 STATS_PERSISTENT defaults to DEFAULT (engine-level innodb_stats_persistent)
**Priority:** LOW
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "STATS_PERSISTENT = {DEFAULT|0|1}; DEFAULT means use innodb_stats_persistent."
**Source:** `sql/sql_table.cc:15361-15364` — `HA_OPTION_STATS_PERSISTENT` / `HA_OPTION_NO_STATS_PERSISTENT` bits cleared when not explicitly set.
**Oracle:** `CREATE_OPTIONS` empty; `SHOW CREATE TABLE` omits the clause.
**omni pointer:** omni Table has no stats fields — **omni gap**. Acceptable since omni does not model row statistics.
**omni assertion:** CREATE succeeds; option discarded.

---

### 8.10 TABLESPACE defaults to the innodb_file_per_table single-table space (rendered as elided)
**Priority:** LOW
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table.html — "TABLESPACE tablespace_name [STORAGE {DISK|MEMORY}] ... If omitted, the file_per_table tablespace is used when innodb_file_per_table=ON."
**Source:** `sql/sql_table.cc:9183-9220` — tablespace option rejection for temporary tables and default application.
**Oracle:** `information_schema.INNODB_TABLES.SPACE` per-table; `SHOW CREATE TABLE` omits `TABLESPACE=` clause.
**omni pointer:** omni does not model tablespaces — **omni gap**. Document as explicitly out of scope; scenario is LOW because no observable catalog difference.
**omni assertion:** CREATE succeeds; no tablespace field stored.

---


## Section C9: Generated column defaults

> **Expansion note (Wave 4 batch A):** grew from 2 to 8 scenarios via walk of
> `sql/sql_table.cc` generated-column handling (`gcol_info` storage L4232/L4739/L7886-L7890),
> FK-vs-gcol interactions L6672/L6783/L6883/L9464, ALTER GCOL expression tracking
> L12108-L12213, and the create-table-generated-columns doc page.

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

### 9.3 VIRTUAL gcol cannot be part of an InnoDB index unless secondary; STORED has no restriction
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, KEY (b));
CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL PRIMARY KEY);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "You cannot define a PRIMARY KEY on a virtual generated column."
**Source:** `sql/sql_table.cc` gcol + key preparation — flags `ER_KEY_BASED_ON_GENERATED_COLUMN` or `ER_WRONG_KEY_COLUMN` when PK contains VIRTUAL gcol. Secondary index is allowed (InnoDB supports indexes on virtual gcols since 5.7).
**Oracle:** `t1` succeeds with secondary KEY on `b`; `t2` errors.
**omni pointer:** `mysql/catalog/tablecmds.go` key resolution must reject VIRTUAL gcol in PK but allow in secondary index. Check `allocIndexName` path.
**omni assertion:** `t1` builds with index; `t2` returns error.

---

### 9.4 Generated column expression must be deterministic — forbids NOW(), UUID(), RAND()
**Priority:** HIGH
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (a INT, b TIMESTAMP GENERATED ALWAYS AS (NOW()) VIRTUAL);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "The generated column expression must be deterministic. Non-deterministic functions such as NOW(), CURRENT_TIMESTAMP, UUID(), RAND() are not permitted."
**Source:** `sql/sql_table.cc` `prepare_create_field` — calls `Item::walk` to check for `FUNC_NONDETERMINISTIC`; raises `ER_GENERATED_COLUMN_NON_DETERMINISTIC` / `ER_GENERATED_COLUMN_FUNCTION_IS_NOT_ALLOWED`.
**Oracle:** Error 3103/3102.
**omni pointer:** `mysql/catalog/tablecmds.go` gcol expression validation — omni's analyzer (Phase 3 `GeneratedAnalyzed`) must flag non-deterministic functions.
**omni assertion:** `Exec` returns error.

---

### 9.5 UNIQUE on STORED gcol is allowed; UNIQUE on VIRTUAL gcol is allowed but limited
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) STORED UNIQUE);
CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL UNIQUE);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "UNIQUE constraints and foreign keys referencing generated columns have restrictions depending on storage mode."
**Source:** `sql/sql_table.cc` — UNIQUE+gcol path accepted for both VIRTUAL and STORED in InnoDB since 5.7; earlier restriction lifted. MyISAM still restricts.
**Oracle:** Both tables succeed under InnoDB; UNIQUE index present for both.
**omni pointer:** `mysql/catalog/tablecmds.go` — verify UNIQUE with gcol does not double-reject. Distinguishes from scenario 9.2/9.7 (FK-specific).
**omni assertion:** Both CREATE succeed; UNIQUE indexes present.

---

### 9.6 Generated column NOT NULL interaction — declared NOT NULL requires expression to never yield NULL
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (
    a INT NULL,
    b INT GENERATED ALWAYS AS (a+1) VIRTUAL NOT NULL
);
INSERT INTO t (a) VALUES (NULL);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "Declaring a generated column NOT NULL does not change the expression's nullability; it only causes an error if the computed value is NULL."
**Source:** `sql/sql_table.cc` `prepare_create_field` — sets column NOT_NULL_FLAG from declaration but defers actual NULL check to insert time.
**Oracle:** CREATE succeeds; INSERT with NULL `a` errors at row time (`ER_BAD_NULL_ERROR`).
**omni pointer:** `mysql/catalog/tablecmds.go` — gcol NOT NULL must be accepted at CREATE. (Runtime insert check is out of scope; mark as spot-check.)
**omni assertion:** CREATE succeeds; `tbl.GetColumn("b").Nullable == false`.

---

### 9.7 FK referencing a STORED gcol from another table is allowed (parent side only)
**Priority:** MED
**Status:** pending
**Setup:**
```sql
CREATE TABLE p (
    a INT,
    b INT GENERATED ALWAYS AS (a+1) STORED,
    UNIQUE KEY (b)
);
CREATE TABLE c (x INT, FOREIGN KEY (x) REFERENCES p(b));
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "A stored generated column can be used as the parent side of a foreign key only if no ON CASCADE / SET NULL / SET DEFAULT action would modify the parent value."
**Source:** `sql/sql_table.cc:9464` `ER_FK_CANNOT_USE_VIRTUAL_COLUMN` triggers only when **child** is virtual; parent STORED is permitted.
**Oracle:** CREATE both succeeds. Adding `ON UPDATE CASCADE` would error because the generated column is not user-writable.
**omni pointer:** `mysql/catalog/tablecmds.go` — FK parent resolution must allow STORED gcol as target; ensure FK action validation rejects cascading updates to gcol parent.
**omni assertion:** CREATE succeeds; FK registered.

---

### 9.8 Generated column expression charset is derived from expression inputs, not column declaration
**Priority:** LOW
**Status:** pending
**Setup:**
```sql
CREATE TABLE t (
    a VARCHAR(10) CHARACTER SET latin1,
    b VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a, 'x')) VIRTUAL
);
```
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-generated-columns.html — "The character set and collation of a generated column are derived from the expression's result type, honoring collation coercibility."
**Source:** `sql/sql_table.cc:7886-7890` Value_generator initialization; `sql/field.cc` `Field::sql_type` picks up the result charset via Item metadata.
**Oracle:** `information_schema.COLUMNS.CHARACTER_SET_NAME` for `b` is `latin1` (inherited from `a`), not the table default.
**omni pointer:** `mysql/catalog/tablecmds.go` gcol column build — `b.Charset` should be resolved from the analyzed expression, not from the table/column declaration path. Check `GeneratedAnalyzed` type propagation.
**omni assertion:** `tbl.GetColumn("b").Charset == "latin1"`.

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

### 10.3 View CHECK OPTION default is CASCADED when WITH CHECK OPTION lacks a qualifier
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/view-check-option.html
> "The default is CASCADED. WITH CHECK OPTION is equivalent to WITH CASCADED CHECK OPTION."
**Source:** `sql/sql_yacc.yy` (rule `view_check_option`: `WITH CHECK_SYM OPTION { $$ = VIEW_CHECK_CASCADED; } | WITH LOCAL_SYM CHECK_SYM OPTION { $$ = VIEW_CHECK_LOCAL; } | WITH CASCADED_SYM CHECK_SYM OPTION { $$ = VIEW_CHECK_CASCADED; }`); `sql/sql_view.cc` `mysql_register_view()` persists `view->with_check`.
**Trigger:** `CREATE VIEW v AS SELECT … WITH CHECK OPTION;` (no LOCAL/CASCADED qualifier).
**Rule:**
- Bare `WITH CHECK OPTION` normalizes to `VIEW_CHECK_CASCADED` (value 1), identical to `WITH CASCADED CHECK OPTION`.
- `WITH LOCAL CHECK OPTION` persists as `VIEW_CHECK_LOCAL` (value 2).
- Omitting the clause entirely stores `VIEW_CHECK_NONE` (value 0) — distinct from CASCADED; rendered as `CHECK_OPTION='NONE'` in I_S.
- Observable via `information_schema.VIEWS.CHECK_OPTION` (`NONE` / `LOCAL` / `CASCADED`).

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION;
CREATE VIEW v2 AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION;
```

**Oracle verification:**
```sql
SELECT TABLE_NAME, CHECK_OPTION FROM information_schema.VIEWS
WHERE TABLE_SCHEMA='testdb' ORDER BY TABLE_NAME;
```
Expected: `v1 → CASCADED`, `v2 → LOCAL`.

**omni pointer:** `mysql/catalog/view.go` must distinguish the three states (`NONE`, `LOCAL`, `CASCADED`) and normalize bare `WITH CHECK OPTION` to CASCADED — do not collapse it to a boolean.

---

### 10.4 ALGORITHM=UNDEFINED resolves to MERGE or TEMPTABLE based on SELECT shape
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/view-algorithms.html
> "For UNDEFINED, MySQL chooses which algorithm to use. It prefers MERGE over TEMPTABLE if possible, because MERGE is usually more efficient and because a view cannot be updatable if a temporary table is used. … If the view contains any of … aggregate functions (SUM(), MIN(), MAX(), COUNT(), and so forth), DISTINCT, GROUP BY, HAVING, LIMIT, UNION, subquery in the select list, assignment to user variables, or refers only to literal values — the view becomes a TEMPTABLE."
**Source:** `sql/sql_view.cc` `mysql_register_view()` → `check_view_mergeable()`; `sql/sql_lex.cc` sets `can_be_merged` flag; resolution stored in view's `algorithm` field on first open.
**Trigger:** `CREATE ALGORITHM=UNDEFINED VIEW v AS <SELECT with GROUP BY>;` — at CREATE time stays UNDEFINED in `.frm`, but at open/resolve time picks TEMPTABLE because SELECT is non-mergeable.
**Rule:**
- The stored `.frm` value remains `UNDEFINED` (or is absent for default) — this is a catalog-layer fact.
- MERGE eligibility is decided at VIEW open, not CREATE — but the catalog representation must preserve the user-declared algorithm verbatim.
- Explicit `ALGORITHM=MERGE` on a non-mergeable SELECT → warning `ER_WARN_VIEW_MERGE` and downgrade to UNDEFINED silently in the `.frm`.
- Explicit `ALGORITHM=TEMPTABLE` forces non-updatable regardless of SELECT shape.

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE ALGORITHM=UNDEFINED VIEW v_agg AS SELECT COUNT(*) FROM t;
CREATE ALGORITHM=MERGE VIEW v_merge_bad AS SELECT DISTINCT a FROM t;
```

**Oracle verification:**
```sql
SHOW CREATE VIEW v_agg;
SHOW CREATE VIEW v_merge_bad;
SHOW WARNINGS;
```
Expected:
- `v_agg` SHOW output contains `ALGORITHM=UNDEFINED`.
- `v_merge_bad` SHOW output contains `ALGORITHM=UNDEFINED` (downgraded); warning 1354 emitted at CREATE time.

**omni pointer:** omni must (a) record the user's declared algorithm before any downgrade and (b) still accept MERGE declarations on DISTINCT/aggregate SELECTs to match server lenience (warning, not error).

---

### 10.5 View SQL SECURITY defaults to DEFINER
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-view.html
> "The default SQL SECURITY characteristic is DEFINER."
**Source:** `sql/sql_yacc.yy` `view_suid` rule (`%empty { $$ = VIEW_SUID_DEFAULT; }`); `sql/sql_view.cc` `mysql_register_view()` — when `view->view_suid == VIEW_SUID_DEFAULT`, persisted as `SQL SECURITY DEFINER` in `.frm`.
**Trigger:** `CREATE VIEW v AS SELECT …` (no `SQL SECURITY` clause).
**Rule:**
- Absent clause → `SECURITY_TYPE = 'DEFINER'` in I_S.
- Persisted verbatim in `SHOW CREATE VIEW` output.
- Only DEFINER and INVOKER are legal (no NONE).

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE VIEW v AS SELECT a FROM t;
```

**Oracle verification:**
```sql
SELECT SECURITY_TYPE FROM information_schema.VIEWS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v';
```
Expected: `DEFINER`.

**omni pointer:** `mysql/catalog/view.go` must default `SqlSecurity` to `DEFINER` on CREATE, not leave empty — matters for deparse round-trip of SHOW CREATE VIEW.

---

### 10.6 View column names default to SELECT expression spelling; expressions get positional aliases
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-view.html
> "If a view column name is not specified, the name is derived from the SELECT statement. Column names are taken from the select-list names if they are valid MySQL identifiers, or the column aliases if AS clauses are used."
**Source:** `sql/sql_view.cc` `mysql_register_view()` — walks `THD::lex->query_block->fields`, uses each item's `item_name.ptr()` as column name; `sql/table.cc` `open_and_read_view()` fills `TABLE_SHARE::field` names.
**Trigger:** View SELECT contains column references vs expressions vs constants without explicit column list.
**Rule:**
- Plain column ref `SELECT a FROM t` → view column name = `a`.
- Expression `SELECT a+1 FROM t` → view column name = `a+1` (literal expression text, not sanitized).
- Aliased `SELECT a+1 AS s FROM t` → view column name = `s`.
- Function call `SELECT COUNT(*) FROM t` → view column name = `COUNT(*)`.
- Explicit list `CREATE VIEW v(x,y) AS SELECT …` wins over any derivation.
- Name length cap: 64 bytes (`NAME_LEN`). Longer derivations are truncated and warnings issued (`ER_NAME_BECOMES_EMPTY` / `ER_REMOVED_SPACES`).

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE VIEW v_auto AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a;
CREATE VIEW v_list (x,y,z) AS SELECT a, a+1, COUNT(*) FROM t GROUP BY a;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_auto' ORDER BY ORDINAL_POSITION;
```
Expected: `a`, `a+1`, `COUNT(*)` (exact spelling).

**omni pointer:** `mysql/catalog/view.go` column-name derivation must use the raw expression text (post-parser, pre-resolver), not a mangled form. A placeholder like `Name_exp_2` is wrong.

---

### 10.7 View updatability is derived from SELECT shape, persisted as IS_UPDATABLE
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/view-updatability.html
> "Some views are updatable and some are not. For a view to be updatable, there must be a one-to-one relationship between the rows in the view and the rows in the underlying table. There are also certain other constructs that make a view nonupdatable: Aggregate functions, DISTINCT, GROUP BY, HAVING, UNION or UNION ALL, … Subquery in the select list, …"
**Source:** `sql/sql_view.cc` `mysql_register_view()` sets `view->updatable_view` based on `can_be_merged` + `is_updatable()` on the query block; persisted in `.frm` header.
**Trigger:** Varies by SELECT shape.
**Rule:**
- `SELECT a FROM t` → updatable.
- `SELECT a FROM t1 JOIN t2 USING(a)` → updatable but columns from only one underlying table are settable per DML.
- `SELECT DISTINCT a FROM t`, `GROUP BY`, any aggregate → non-updatable.
- `ALGORITHM=TEMPTABLE` forces non-updatable regardless of shape.
- Persisted as boolean in `information_schema.VIEWS.IS_UPDATABLE` (`YES`/`NO`).

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE VIEW v_ok AS SELECT a FROM t;
CREATE VIEW v_distinct AS SELECT DISTINCT a FROM t;
CREATE ALGORITHM=TEMPTABLE VIEW v_temp AS SELECT a FROM t;
```

**Oracle verification:**
```sql
SELECT TABLE_NAME, IS_UPDATABLE FROM information_schema.VIEWS
WHERE TABLE_SCHEMA='testdb' ORDER BY TABLE_NAME;
```
Expected: `v_distinct=NO`, `v_ok=YES`, `v_temp=NO`.

**omni pointer:** catalog view representation needs an `IsUpdatable` field; cannot be computed lazily because it depends on parser+resolver analysis that omni may not run.

---

### 10.8 View column nullability is widened vs base columns
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-view.html (nullability mapping implicit).
**Source:** `sql/sql_view.cc` `mysql_register_view()` → `Create_field::init_for_tmp_table()` inherits base column `maybe_null`; outer-join/UNION branches force `maybe_null = true`.
**Trigger:** View built over LEFT JOIN, UNION, or a column reached via an outer-joined table.
**Rule:**
- Even if `t.a` is `NOT NULL`, a view `SELECT t2.a FROM t1 LEFT JOIN t2 ON …` produces a nullable view column.
- `UNION` of two NOT NULL columns stays NOT NULL only if both sides are NOT NULL; UNION with a literal `NULL` column flips to nullable.
- Aggregates over a NOT NULL column: `COUNT(*)` NOT NULL, `SUM(a)` nullable (returns NULL on empty group), `MIN/MAX` nullable.
- Observable via `information_schema.COLUMNS.IS_NULLABLE` on the view name.

**Setup:**
```sql
CREATE TABLE t1 (id INT NOT NULL, a INT NOT NULL);
CREATE TABLE t2 (id INT NOT NULL, b INT NOT NULL);
CREATE VIEW v AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.id = t2.id;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, IS_NULLABLE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v' ORDER BY ORDINAL_POSITION;
```
Expected: `a → YES` (outer-join side propagates), `b → YES` (optional side).

**omni pointer:** omni's view column resolver must track outer-join nullability, not just copy `NotNull` from the underlying column. A simple "view column inherits base column flags" rule is wrong.

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

### 11.2 Trigger SQL SECURITY is always DEFINER (no INVOKER option)
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
> "Triggers currently are not activated by foreign key actions. … Within the trigger body, CURRENT_USER returns the account used to check privileges at trigger activation time. This is the DEFINER user, not the user whose actions caused the trigger to be activated."
**Source:** `sql/sql_trigger.cc` `Trigger::create_from_parser()`; unlike views, the CREATE TRIGGER grammar has no `SQL SECURITY` clause — there is no `view_suid` equivalent.
**Trigger:** `CREATE TRIGGER …` (no way to specify security mode).
**Rule:**
- Grammar rejects `SQL SECURITY INVOKER` on triggers — parser error.
- Trigger always runs with DEFINER privileges; `information_schema.TRIGGERS` has no `SECURITY_TYPE` column for triggers.
- Catalog representation should not expose a trigger-level `SqlSecurity` field; if present, it must be locked to DEFINER.

**Setup:**
```sql
CREATE TABLE t (a INT);
-- valid
CREATE DEFINER='root'@'%' TRIGGER trg1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a=1;
-- invalid grammar:
-- CREATE TRIGGER trg2 SQL SECURITY INVOKER BEFORE INSERT ON t FOR EACH ROW ...;
```

**Oracle verification:** the second form fails with `ER_PARSE_ERROR`.

**omni pointer:** omni trigger catalog must NOT carry a nullable SqlSecurity; the value is implicit and uniform.

---

### 11.3 Trigger CHARACTER_SET_CLIENT / COLLATION_CONNECTION / DATABASE_COLLATION snapshot session state at CREATE time
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
> "MySQL takes the character set and collation of character columns into consideration when performing operations. … These attributes are stored with the trigger and used again each time the trigger is activated."
**Source:** `sql/sql_trigger.cc` — `Trigger::create_from_parser()` records `thd->variables.character_set_client`, `thd->variables.collation_connection`, and the schema's default collation at the moment of CREATE; persisted in `.TRN`/data-dictionary and surfaced as columns in `information_schema.TRIGGERS`.
**Trigger:** `CREATE TRIGGER …` under a specific `SET NAMES` session.
**Rule:**
- Three separate columns persisted per trigger: `CHARACTER_SET_CLIENT`, `COLLATION_CONNECTION`, `DATABASE_COLLATION`.
- Default values depend on server config (usually `utf8mb4` / `utf8mb4_0900_ai_ci` / schema default).
- These matter for SDL diff: if a dump preserves `SET NAMES` before the CREATE TRIGGER, the snapshot is reproducible; if not, the catalog must record them independently.

**Setup:**
```sql
SET NAMES utf8mb4;
CREATE TABLE t (a INT);
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a;
```

**Oracle verification:**
```sql
SELECT CHARACTER_SET_CLIENT, COLLATION_CONNECTION, DATABASE_COLLATION
FROM information_schema.TRIGGERS WHERE TRIGGER_NAME='trg';
```
Expected: three non-null values reflecting session state at CREATE.

**omni pointer:** `mysql/catalog/trigger.go` must carry these three fields; without them, deparse→reparse can shift character interpretation inside trigger bodies.

---

### 11.4 Trigger ACTION_ORDER defaults to 1, incremented per (table, timing, event) group
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
> "It is possible to have multiple triggers for a given table that have the same trigger event and action time. … By default, triggers that have the same trigger event and action time activate in the order they were created. To affect trigger order, specify FOLLOWS or PRECEDES."
**Source:** `sql/sql_trigger.cc` `Trigger::create_from_parser()` — when FOLLOWS/PRECEDES omitted, new trigger gets `action_order = max(existing on same table/timing/event) + 1`; stored in `information_schema.TRIGGERS.ACTION_ORDER`.
**Trigger:** Multiple `BEFORE INSERT` triggers on the same table without explicit ordering.
**Rule:**
- First such trigger: `ACTION_ORDER = 1`.
- Each subsequent: `max + 1` within the same `(table, timing, event)` tuple.
- Ordering is per-group: BEFORE INSERT and AFTER INSERT have independent sequences.
- `FOLLOWS other_trg` / `PRECEDES other_trg` splices into the existing order and shifts later triggers (catalog renumbers `action_order`).
- `ACTION_ORDER = 0` reserved for "pre-migration" triggers from older MySQL versions (no ordering info).

**Setup:**
```sql
CREATE TABLE t (a INT);
CREATE TRIGGER trg_a BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1;
CREATE TRIGGER trg_b BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 2;
CREATE TRIGGER trg_c BEFORE INSERT ON t FOR EACH ROW PRECEDES trg_a SET NEW.a = NEW.a + 10;
```

**Oracle verification:**
```sql
SELECT TRIGGER_NAME, ACTION_ORDER FROM information_schema.TRIGGERS
WHERE TRIGGER_SCHEMA='testdb' ORDER BY ACTION_ORDER;
```
Expected: `trg_c=1, trg_a=2, trg_b=3` (PRECEDES shifts trg_a and trg_b down).

**omni pointer:** trigger catalog must carry `ActionOrder` and support renumbering on FOLLOWS/PRECEDES. A name-only store cannot round-trip.

---

### 11.5 NEW and OLD pseudo-rows availability depends on event (INSERT/UPDATE/DELETE)
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/trigger-syntax.html
> "Use the OLD and NEW keywords to access columns in the rows affected by a trigger (OLD and NEW are not case-sensitive). … In an INSERT trigger, only NEW.col_name can be used; there is no old row. In a DELETE trigger, only OLD.col_name can be used; there is no new row. In an UPDATE trigger, you can use OLD.col_name to refer to the columns of a row before it is updated and NEW.col_name to refer to the columns of the row after it is updated."
**Source:** `sql/sp.cc` trigger body compile phase; `sql/item.cc` `Item_trigger_field::fix_fields()` validates NEW/OLD access against trigger event.
**Trigger:** Trigger body references `OLD.col` in an INSERT trigger or `NEW.col` in a DELETE trigger.
**Rule:**
- `INSERT` trigger: only `NEW.*` legal; `OLD.*` → `ER_TRG_NO_SUCH_ROW_IN_TRG`.
- `DELETE` trigger: only `OLD.*` legal; `NEW.*` → `ER_TRG_NO_SUCH_ROW_IN_TRG`.
- `UPDATE` trigger: both `OLD.*` and `NEW.*` legal.
- BEFORE vs AFTER further restricts writability: `SET NEW.col = …` legal only in BEFORE INSERT / BEFORE UPDATE; AFTER triggers cannot assign `NEW.*` (read-only).
- `OLD.*` is always read-only.

**Setup:**
```sql
CREATE TABLE t (a INT);
-- legal:
CREATE TRIGGER t_ins BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NEW.a + 1;
CREATE TRIGGER t_del BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a;
-- illegal (parse-accepted, rejected at CREATE TRIGGER):
-- CREATE TRIGGER bad BEFORE INSERT ON t FOR EACH ROW SET @x = OLD.a;
```

**Oracle verification:** the commented-out CREATE fails with `ER_TRG_NO_SUCH_ROW_IN_TRG` (1363).

**omni pointer:** trigger body validator in omni must check NEW/OLD references against trigger event+timing. Without this check, diffing trigger bodies can accept semantically-impossible forms.

---

### 11.6 Trigger on partitioned table is allowed; one trigger fires per affected row regardless of partition
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/partitioning-limitations.html
> "Triggers on partitioned tables are supported. A trigger defined on a partitioned table applies to all partitions."
**Source:** `sql/sql_trigger.cc` — no partition filter; trigger metadata is stored at table level, not partition level.
**Trigger:** `CREATE TRIGGER` on a `PARTITION BY HASH` / RANGE / LIST table.
**Rule:**
- Trigger definition is accepted and stored once at table level.
- Catalog must not carry any per-partition trigger metadata — the trigger→table relationship is N:1 where N is any count.
- `ALTER TABLE … DROP PARTITION` does NOT drop triggers; trigger survives partition churn.
- `ALTER TABLE … REORGANIZE PARTITION` preserves triggers.

**Setup:**
```sql
CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 4;
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.b = NEW.a * 2;
ALTER TABLE t COALESCE PARTITION 2;
```

**Oracle verification:**
```sql
SELECT TRIGGER_NAME FROM information_schema.TRIGGERS
WHERE TRIGGER_SCHEMA='testdb' AND EVENT_OBJECT_TABLE='t';
```
Expected: `trg` still present after COALESCE.

**omni pointer:** trigger storage must live on the table node, not hang off partition definitions. Confirm that partition mutation paths do not touch trigger lists.

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

### 14.2 ALTER TABLE … ALTER CHECK … NOT ENFORCED / ENFORCED toggles the flag
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html
> "An ALTER TABLE statement that modifies a constraint's enforcement does not rebuild the table. ALTER TABLE tbl ALTER CHECK chk [NOT] ENFORCED;"
**Source:** `sql/sql_table.cc` `Sql_cmd_alter_table::execute()` — `Alter_column::ALTER_CHECK_CONSTRAINT_ENFORCEMENT` path; INPLACE/INSTANT capable because only metadata flips.
**Trigger:** `ALTER TABLE t ALTER CHECK c NOT ENFORCED;` and `ALTER TABLE t ALTER CHECK c ENFORCED;`.
**Rule:**
- Toggle is metadata-only — no table rebuild, no row re-evaluation of existing data.
- Flipping to ENFORCED does NOT validate existing rows against the predicate (existing rows can violate).
- Catalog state changes atomically; both directions supported.
- `ALTER TABLE t ALTER CONSTRAINT c [NOT] ENFORCED` (new 8.0.16 syntax) is equivalent for CHECK constraints and additionally applies to FK (but FK enforcement is a separate concept).

**Setup:**
```sql
CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0));
ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED;
```

**Oracle verification:**
```sql
SELECT ENFORCED FROM information_schema.CHECK_CONSTRAINTS
WHERE CONSTRAINT_NAME='c_pos';
```
Expected: `NO` after the ALTER, `YES` after re-enabling.

**omni pointer:** catalog must expose a mutable `Enforced` flag on CheckConstraint and have an ALTER path that doesn't rebuild the check expression.

---

### 14.3 STORED generated column + CHECK constraint: predicate evaluated at INSERT/UPDATE time on the stored value
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html
> "A constraint expression can refer to any column in the table, including generated columns."
**Source:** `sql/sql_table.cc` `prepare_check_constraints_for_create()` validates that CHECK expressions referring to generated columns compile; `sql/field.cc` evaluation order guarantees STORED generated columns are computed before CHECK evaluation.
**Trigger:** `CHECK (gcol > 0)` where `gcol INT AS (a * 2) STORED`.
**Rule:**
- Legal: CHECK may reference STORED generated columns; evaluation order is (regular defaults) → (generated columns) → (CHECK).
- Legal: CHECK may reference VIRTUAL generated columns — MySQL re-evaluates the expression at check time.
- CHECK expression itself is not allowed to be "generated" — but can reference generated columns.
- ALTER that toggles a generated column's STORED/VIRTUAL mode is forbidden while a CHECK references it (locks the representation).

**Setup:**
```sql
CREATE TABLE t (
  a INT,
  g INT AS (a * 2) STORED,
  CONSTRAINT c_g_nonneg CHECK (g >= 0)
);
INSERT INTO t (a) VALUES (5);   -- OK, g=10
INSERT INTO t (a) VALUES (-1);  -- CHECK violation (g=-2)
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS
WHERE CONSTRAINT_NAME='c_g_nonneg';
```
Expected: clause reads `(`g` >= 0)`; second INSERT fails with `ER_CHECK_CONSTRAINT_VIOLATED` (3819).

**omni pointer:** CHECK expression resolver in omni must walk into generated-column references; simple "is this identifier a plain column?" analysis misses the generated case.

---

### 14.4 CHECK constraint cannot contain a subquery, stored function, or non-deterministic built-in
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/create-table-check-constraints.html
> "Subqueries are not permitted. Stored functions and loadable functions are not permitted. Stored procedures and function parameters are not permitted. Variables (system, user, and stored program local variables) are not permitted. Nondeterministic functions (such as NOW() or CONNECTION_ID()) are not permitted."
**Source:** `sql/sql_table.cc` `pre_validate_check_constraint()` walks the parsed expression tree and rejects `Item_subselect`, `Item_func_sp`, `Item_func_udf`, variable refs, and functions whose `used_tables() & RAND_TABLE_BIT` is set.
**Trigger:** `CHECK` clause containing any of: subquery, stored function, user variable, `NOW()`, `RAND()`, `CONNECTION_ID()`.
**Rule:**
- Parse-time errors (not CREATE-time errors — yacc-level rejection in `check_constraint_spec`):
  - Subquery → `ER_CHECK_CONSTRAINT_NOT_ALLOWED_CONTEXT` (3812) or `ER_CHECK_CONSTRAINT_CONTAINS_SUBQUERY`.
  - Stored/UDF → `ER_CHECK_CONSTRAINT_FUNCTION_IS_NOT_ALLOWED` (3814).
  - `NOW()` / `CURRENT_TIMESTAMP` / `RAND()` → `ER_CHECK_CONSTRAINT_NAMED_FUNCTION_IS_NOT_ALLOWED` (3815).
  - User variable `@v` / system variable `@@v` → `ER_CHECK_CONSTRAINT_VARIABLES` (3813).
- Deterministic built-ins (`LENGTH`, `CHAR_LENGTH`, arithmetic, string compare, `JSON_VALID`, `REGEXP_LIKE`) are allowed.
- omni SDL diff must preserve CHECK expression verbatim — cannot rewrite or normalize away deterministic functions.

**Setup:**
```sql
CREATE TABLE t1 (a INT, CHECK (CHAR_LENGTH(CAST(a AS CHAR)) < 10));  -- OK
-- illegal:
-- CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)));    -- 3812
-- CREATE TABLE t3 (a INT, CHECK (a < NOW()));                      -- 3815
-- CREATE TABLE t4 (a INT, CHECK (a < @max));                       -- 3813
```

**Oracle verification:** the four commented-out CREATE statements all fail at parse time with the codes listed above.

**omni pointer:** catalog CHECK validator must classify the expression tree; a string-level accept allows forbidden constructs that bytebase would round-trip and the server would then reject.

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

### 15.2 ADD COLUMN … FIRST inserts before column 1, shifting all others down
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
> "To add a column at a specific position within a table row, use FIRST or AFTER col_name. The default is to add the column last. You can also use FIRST and AFTER in CHANGE or MODIFY operations to reorder columns within a table."
**Source:** `sql/sql_table.cc` `mysql_prepare_alter_table()` — `Alter_column::ALTER_COLUMN_ORDER` with `after_field == nullptr` and `first` flag set; rebuilds `create_list` with the new column spliced at index 0.
**Trigger:** `ALTER TABLE t ADD COLUMN x INT FIRST;`.
**Rule:**
- New column's `ORDINAL_POSITION = 1`.
- All existing columns' ordinal positions shift +1.
- Operation was COPY-only pre-8.0.29; since 8.0.29 it is INSTANT when column count fits into row header reserved bits, INPLACE otherwise.
- FIRST is rejected with `ER_BAD_FIELD_ERROR` if a column with the target name already exists.

**Setup:**
```sql
CREATE TABLE t (a INT, b INT);
ALTER TABLE t ADD COLUMN c INT FIRST;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
```
Expected: `c=1, a=2, b=3`.

**omni pointer:** SDL diff for column order must emit `FIRST` rather than a chain of `AFTER` moves when the target is position 1 — otherwise round-trip produces different `.frm` ordering ancestry.

---

### 15.3 ADD COLUMN … AFTER col inserts at col+1
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table.html (same paragraph as 15.2).
**Source:** `sql/sql_table.cc` `mysql_prepare_alter_table()` `Alter_column::ALTER_COLUMN_ORDER` branch where `after_field` names an existing column; uses `List_iterator<Create_field>` to splice after the named field.
**Trigger:** `ALTER TABLE t ADD COLUMN x INT AFTER a;`.
**Rule:**
- New column's ORDINAL_POSITION = (position of `a`) + 1.
- Columns after `a` shift +1.
- `AFTER unknown_col` → `ER_BAD_FIELD_ERROR` (1054).
- `AFTER` a column added in the same ALTER (multi-ADD) is legal — the splice is resolved left-to-right.

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, c INT);
ALTER TABLE t ADD COLUMN x INT AFTER a;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
```
Expected: `a=1, x=2, b=3, c=4`.

**omni pointer:** SDL diff must treat "insert after named col" as the canonical form when the new column is not at position 1.

---

### 15.4 MODIFY col and CHANGE col both retain existing position unless FIRST/AFTER given
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
> "To change a column's type or attributes, use CHANGE or MODIFY. … If you do not use FIRST or AFTER, the column remains in its original position."
**Source:** `sql/sql_table.cc` `mysql_prepare_alter_table()` — `Alter_column::ALTER_COLUMN_CHANGE` / `ALTER_COLUMN_MODIFY` default keeps the original ordinal; only if `after_field` or `first` is set does the column move.
**Trigger:** `ALTER TABLE t MODIFY a BIGINT;` (no position clause).
**Rule:**
- Position unchanged — distinguishes MySQL's MODIFY/CHANGE from other DBs where type change can reposition.
- `CHANGE old new BIGINT FIRST` moves to position 1.
- `MODIFY a BIGINT AFTER c` moves to position after c.
- If FIRST/AFTER is supplied, the position move is a separate INSTANT-capable metadata op combined with the type change (which may itself be COPY).

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, c INT);
ALTER TABLE t MODIFY b BIGINT;
```

**Oracle verification:** ordinal positions remain `a=1, b=2, c=3`.

**omni pointer:** SDL diff must NOT emit a trailing `AFTER a` when the user only changed the type — emitting spurious position hints causes round-trip drift. Check the source column's pre-existing position before adding FIRST/AFTER.

---

### 15.5 Multiple ADD COLUMN in one ALTER: positions resolve left-to-right against the evolving column list
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
> "You can perform multiple add, alter, drop, and change operations with a single ALTER TABLE statement, separated by commas. This is a MySQL extension to standard SQL, which permits only one of each clause per ALTER TABLE statement."
**Source:** `sql/sql_table.cc` `mysql_prepare_alter_table()` iterates `alter_info->create_list` in declaration order; each ADD splices into the current working list, so a later ADD can reference a name added earlier in the same statement.
**Trigger:** `ALTER TABLE t ADD COLUMN x INT AFTER a, ADD COLUMN y INT AFTER x;`.
**Rule:**
- Order of evaluation = lexical order of ADD clauses.
- Later clauses see earlier-added columns by name.
- `ADD COLUMN x AFTER x` (self-reference) → `ER_BAD_FIELD_ERROR`.
- `ADD COLUMN x AFTER y, ADD COLUMN y AFTER a` (forward reference) → `ER_BAD_FIELD_ERROR` (left-to-right, not order-independent).
- Two `FIRST` clauses in the same ALTER: the second one wins at position 1, the first shifts to position 2.

**Setup:**
```sql
CREATE TABLE t (a INT, b INT);
ALTER TABLE t ADD COLUMN x INT AFTER a, ADD COLUMN y INT AFTER x;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
```
Expected: `a=1, x=2, y=3, b=4`.

**omni pointer:** omni ALTER batch processor must stream column list mutations left-to-right when computing resulting positions; a naive "apply all ADDs against original list" is wrong.

---


## Section C16: Date/time function precision defaults

### 16.1 NOW/CURRENT_TIMESTAMP precision defaults to 0 (seconds)
**Priority:** HIGH
**Status:** pending-verify
**Source:** `sql/sql_yacc.yy:7370-7378` (`func_datetime_precision` rule), `sql/item_timefunc.cc:2078-2084`
**Trigger:** `NOW()`, `CURRENT_TIMESTAMP`, `NOW(n)`, `CURRENT_TIMESTAMP(n)` without explicit precision
**Rule:**
- Parser rule: `func_datetime_precision: %empty { $$= 0; } | '(' ')' { $$= 0; } | '(' NUM ')' { $$= parse_num } `
- Default precision = 0 (no fractional seconds) when omitted
- If `NOW()` or `NOW` called with empty parentheses, precision = 0
- If `NOW(n)` called with explicit precision, uses `n` (0-6)
- Precision stored in `Item_func_now::decimals` field

**Edge cases:**
- `CURTIME()`, `UTC_TIME()`, `SYSDATE()` follow same rule (yacc 10763, 10818, 10835)
- `CURDATE()`, `UTC_DATE()` have NO precision parameter (no `func_datetime_precision` rule)
- Precision > 6 is clamped/rejected (`check_precision()` — sql/item_timefunc.cc:800-807, `ER_TOO_BIG_PRECISION`)

**Observable via:**
- `SHOW CREATE TABLE` on generated/virtual columns (`DATETIME(n)` vs `DATETIME`)
- `information_schema.COLUMNS.DATETIME_PRECISION`
- Query result string length: `NOW()` = "YYYY-MM-DD HH:MM:SS" (19 chars), `NOW(6)` = "YYYY-MM-DD HH:MM:SS.NNNNNN" (26 chars)

**omni pointer:** `mysql/catalog/function_types.go:32` returns `BaseTypeDateTime` for `now`/`current_timestamp` without any `decimals`/`Fsp` field — precision is lost.

---

### 16.2 NOW(N) explicit precision 0..6 (and > 6 rejected)
**Priority:** HIGH
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html
> "Functions that take temporal arguments accept values with fractional seconds. Return values from temporal functions include fractional seconds as appropriate. For example, `NOW()` returns the current date and time with no fractional part, but `NOW(6)` returns the current date and time with fractional seconds to microseconds precision."
**Source:**
- `sql/item_timefunc.cc:800-807` (`Item_temporal_func::check_precision`):
  ```cpp
  if (decimals > DATETIME_MAX_DECIMALS) {
    my_error(ER_TOO_BIG_PRECISION, MYF(0), (int)decimals, func_name(),
             DATETIME_MAX_DECIMALS);
    return true;
  }
  ```
- `DATETIME_MAX_DECIMALS` = 6 (header constant)
**Trigger:** `NOW(0)`..`NOW(6)` accepted; `NOW(7)` → `ER_TOO_BIG_PRECISION`.
**Rule:**
- Valid `N` range: 0 ≤ N ≤ 6.
- `NOW(7)`, `CURTIME(7)`, `SYSDATE(7)` all fail parse-time (precision checked during `resolve_type`).
- Negative N is blocked by lexer (`NUM` unsigned).

**Observable via:**
- `SELECT NOW(7)` → error 1426 "Too-big precision 7 specified for 'now'. Maximum is 6."
- `SELECT LENGTH(NOW(6))` = 26 (19 + 1 dot + 6 digits).

**omni pointer:** omni does not validate fsp range in `function_types.go` — `NOW(7)` would silently type-check as DATETIME.

---

### 16.3 CURDATE / CURRENT_DATE / UTC_DATE take no precision argument
**Priority:** MED
**Status:** pending
**Source:** `sql/sql_yacc.yy:10759-10762` (CURDATE rule uses `optional_braces`, not `func_datetime_precision`); `sql/sql_yacc.yy:10831-10833` (UTC_DATE_SYM optional_braces)
**Trigger:** `CURDATE()`, `CURRENT_DATE`, `CURRENT_DATE()`, `UTC_DATE()`
**Rule:**
- Parser does NOT consume a precision argument — `CURDATE(6)` is a **parse error**, not a silent drop.
- Result type is DATE (10-char `YYYY-MM-DD`), never fractional.
- `CURRENT_DATE` is a synonym for `CURDATE` (defined in same yacc rule block).
- `DATETIME_PRECISION` for DATE columns/expressions is NULL in `information_schema`.

**Observable via:**
- `SELECT CURDATE(6)` → syntax error near `(6)`.
- `CREATE TABLE t AS SELECT CURDATE()` → column type `date`, no precision.

**omni pointer:** `mysql/catalog/function_types.go:34` returns `BaseTypeDate` for `curdate`/`current_date`; omni parser should reject `CURDATE(6)` — needs verification.

---

### 16.4 CURTIME / CURRENT_TIME / UTC_TIME precision defaults to 0
**Priority:** MED
**Status:** pending
**Source:** `sql/sql_yacc.yy:10763-10766` (CURTIME `func_datetime_precision`), 10835-10838 (UTC_TIME_SYM `func_datetime_precision`)
**Trigger:** `CURTIME()`, `CURTIME(n)`, `CURRENT_TIME`, `CURRENT_TIME(n)`, `UTC_TIME()`, `UTC_TIME(n)`
**Rule:**
- Parsed with `func_datetime_precision` → defaults to 0 if omitted.
- Precision stored in `Item_func_curtime::decimals` / `Item_func_curtime_utc::decimals`.
- Returns TIME type with optional fractional seconds (0..6).
- `CURRENT_TIME` is a synonym for `CURTIME`.

**Observable via:**
- `LENGTH(CURTIME())` = 8 (`HH:MM:SS`), `LENGTH(CURTIME(6))` = 15 (`HH:MM:SS.NNNNNN`).
- `SELECT CURTIME(7)` → `ER_TOO_BIG_PRECISION`.

**omni pointer:** `mysql/catalog/function_types.go:36` returns `BaseTypeTime` for `curtime`/`current_time`; no fsp propagation.

---

### 16.5 SYSDATE precision defaults to 0 (distinct from NOW in evaluation semantics)
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_sysdate
> "SYSDATE() returns the time at which it executes. This differs from the behavior for NOW(), which returns a constant time that indicates the time at which the statement began to execute."
**Source:** `sql/sql_yacc.yy:10818-10822` (SYSDATE `func_datetime_precision` → `PTI_function_call_nonkeyword_sysdate`)
**Trigger:** `SYSDATE()`, `SYSDATE(n)`
**Rule:**
- Same fsp parsing rule as NOW — default 0, range 0..6.
- **Semantic difference from NOW:** SYSDATE is evaluated per row (non-constant), while NOW is a per-statement constant. This affects replication (`sysdate_is_now` flag) but not schema/type.
- Cannot be used as a column `DEFAULT` (not `real_type_with_now_as_default`), so PS5 fsp-match rule does not apply.

**Observable via:**
- Column spec: `DATETIME DEFAULT SYSDATE` → `ER_INVALID_DEFAULT` (not an `Item_func_now`, fails `add_field` check at `sql/sql_parse.cc:5521`).
- Result length identical to `NOW(n)` for same `n`.

**omni pointer:** `mysql/catalog/function_types.go:32` buckets `sysdate` with `now`/`current_timestamp` — loses the DEFAULT-eligibility distinction.

---

### 16.6 UTC_TIMESTAMP precision defaults to 0
**Priority:** LOW
**Status:** pending
**Source:** `sql/sql_yacc.yy:10839-10842` (UTC_TIMESTAMP_SYM `func_datetime_precision` → `Item_func_now_utc`)
**Trigger:** `UTC_TIMESTAMP()`, `UTC_TIMESTAMP(n)`
**Rule:**
- Built on `Item_func_now_utc` (a subclass of `Item_func_now`). Default fsp=0, range 0..6.
- Returns DATETIME (not TIMESTAMP) in UTC zone — see `sql/item_timefunc.cc:2078-2084` (`Item_func_now::resolve_type`).
- `UTC_TIMESTAMP()` is NOT allowed as a column `DEFAULT` — only `CURRENT_TIMESTAMP`/`NOW()` pass the `real_type_with_now_as_default` check.

**Observable via:**
- `CREATE TABLE t (a DATETIME DEFAULT UTC_TIMESTAMP)` → `ER_INVALID_DEFAULT`.
- `SELECT LENGTH(UTC_TIMESTAMP())` = 19; `LENGTH(UTC_TIMESTAMP(3))` = 23.

**omni pointer:** `mysql/catalog/function_types.go` has no `utc_timestamp` / `utc_time` / `utc_date` entries — unknown function fallback.

---

### 16.7 UNIX_TIMESTAMP return type depends on argument type and fsp
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_unix-timestamp
> "If UNIX_TIMESTAMP() is called with no date argument, it returns an unsigned integer. If UNIX_TIMESTAMP() is called with a date argument, it returns the value of the argument as seconds since '1970-01-01 00:00:00' UTC. The return value is a decimal with a scale of 0 if the argument is an integer type. Otherwise, the return value is a decimal whose scale is the same as the scale of the argument."
**Source:** `sql/item_timefunc.h:401` (`Item_func_unix_timestamp final : public Item_timeval_func`); `sql/item_create.cc:1669` (`{"UNIX_TIMESTAMP", SQL_FN_V(Item_func_unix_timestamp, 0, 1)}`)
**Trigger:** `UNIX_TIMESTAMP()`, `UNIX_TIMESTAMP(expr)`
**Rule:**
- Zero-arg: returns BIGINT UNSIGNED with fsp=0.
- One-arg: return type is `DECIMAL(M, args[0]->decimals)`. If the argument is a DATETIME(6) expression, result is `DECIMAL(..., 6)`; if DATETIME (fsp=0), result is integer-like `DECIMAL(..., 0)`.
- String argument: parses as DATETIME, fsp derived from the literal's fractional part.

**Observable via:**
- `SELECT UNIX_TIMESTAMP(NOW(6))` → value like `1712912345.123456` with DECIMAL(20,6).
- `SELECT UNIX_TIMESTAMP()` → integer.

**omni pointer:** omni's `function_types.go` has no `unix_timestamp` entry — unknown return type, no fsp inference.

---

### 16.8 DATETIME(N) column with `DEFAULT NOW()` where N > 0 → ER_INVALID_DEFAULT (refines PS5)
**Priority:** HIGH
**Status:** pending
**Cross-ref:** **PS5** in catalog lines 1300-1321.
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/fractional-seconds.html
> "To define a column that includes a fractional seconds part, use the syntax `type_name(fsp)`... The fsp value, if given, must be in the range 0 to 6."
**Source:** `sql/sql_parse.cc:5521` (`Alter_info::add_field`):
```cpp
uint8 datetime_precision = decimals ? atoi(decimals) : 0;
...
if (func->functype() != Item_func::NOW_FUNC ||
    !real_type_with_now_as_default(type) ||
    default_value->decimals != datetime_precision) {
    my_error(ER_INVALID_DEFAULT, MYF(0), field_name->str);
    return true;
}
```
**Trigger:** `CREATE TABLE t (a DATETIME(6) DEFAULT NOW())` — column fsp 6, default fsp 0.
**Rule:** The `fsp` of `DEFAULT NOW(fsp)` **must equal** the column's declared precision. MySQL does NOT silently upconvert. Valid forms:
- `DATETIME DEFAULT NOW()` — both fsp 0 ✓
- `DATETIME(6) DEFAULT NOW(6)` — both fsp 6 ✓
- `DATETIME(3) DEFAULT NOW()` — ✗ `ER_INVALID_DEFAULT` (3 ≠ 0)
- `DATETIME(3) DEFAULT NOW(6)` — ✗ `ER_INVALID_DEFAULT` (3 ≠ 6)

**Symmetry:** CREATE and ALTER (ADD/MODIFY/CHANGE COLUMN) both hit `add_field`.

**omni pointer:** `mysql/catalog/altercmds.go:828-832` stores `col.Default = &s` via `nodeToSQL(colDef.DefaultValue)` — no fsp comparison, no rejection. **Strictness gap** (catalog will accept DDL that MySQL rejects).

---

### 16.9 ON UPDATE NOW(N) must match column fsp (mirrors PS5)
**Priority:** HIGH
**Status:** pending
**Source:** `sql/sql_parse.cc:5603` — same `add_field` check applied to `on_update_value` (see PS5 note: "The same check also applies to on_update_value at line 5603").
**Trigger:** `CREATE TABLE t (a DATETIME(6) ON UPDATE NOW())` or `... ON UPDATE NOW(3)` on DATETIME(6).
**Rule:**
- `ON UPDATE` value must be `NOW_FUNC` AND column must be `real_type_with_now_as_default` (i.e. DATETIME or TIMESTAMP) AND `on_update->decimals == column_fsp`.
- `ON UPDATE` is permitted only on TIMESTAMP and DATETIME (not DATE, TIME, etc.).
- Mismatched fsp → `ER_INVALID_ON_UPDATE` (different errcode from DEFAULT, but same structural rule).

**Observable via:**
- `CREATE TABLE t (a DATETIME(6) DEFAULT NOW(6) ON UPDATE NOW())` → error (ON UPDATE fsp 0 vs col fsp 6).
- `CREATE TABLE t (a DATE ON UPDATE NOW())` → `ER_INVALID_ON_UPDATE` (DATE not allowed).

**omni pointer:** `mysql/catalog/tablecmds.go:214` writes `col.OnUpdate = nodeToSQL(colDef.OnUpdate)` verbatim; no fsp or type check.

---

### 16.10 DATETIME/TIMESTAMP storage bytes scale with fsp
**Priority:** MED
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/storage-requirements.html (Date and Time Type Storage Requirements)
> "A DATETIME value requires 5 bytes for the non-fractional part. Fractional seconds storage requires additional bytes... fsp 1-2 = 1 byte, 3-4 = 2 bytes, 5-6 = 3 bytes."
**Source:** `libbinlogevents/src/binary_log_funcs.cpp:70-88`:
```cpp
unsigned int my_datetime_binary_length(unsigned int dec) {
  return 5 + (dec + 1) / 2;
}
unsigned int my_timestamp_binary_length(unsigned int dec) {
  return 4 + (dec + 1) / 2;
}
```
Also referenced by `sql/field.h:3430` (`Field_datetimef::pack_length`).
**Trigger:** Column type inference with explicit fsp.
**Rule:**
- `DATETIME` (no fsp) → fsp=0, 5 bytes.
- `DATETIME(1)` / `DATETIME(2)` → 5 + 1 = 6 bytes.
- `DATETIME(3)` / `DATETIME(4)` → 5 + 2 = 7 bytes.
- `DATETIME(5)` / `DATETIME(6)` → 5 + 3 = 8 bytes.
- `TIMESTAMP` same formula but base = 4. `TIMESTAMP(6)` = 7 bytes.
- `TIME(fsp)` uses `my_time_binary_length` (base 3).

**Observable via:**
- `information_schema.COLUMNS.DATETIME_PRECISION` reflects the fsp.
- `SHOW CREATE TABLE` emits the explicit `(fsp)` suffix when fsp > 0; MySQL omits the parenthesized suffix when fsp = 0 (i.e. `DATETIME`, not `DATETIME(0)`).

**omni pointer:** `mysql/catalog/table.go:62` stores `DataType string` — omni keeps the raw normalized form (e.g. `"datetime(6)"`). Deparse in `mysql/catalog/show.go` needs to round-trip `DATETIME(0)` → `DATETIME` correctly. Spot-check needed.

---

### 16.11 YEAR(N) deprecated — only YEAR(4) accepted, normalized to YEAR
**Priority:** LOW
**Status:** pending
**Doc:** https://dev.mysql.com/doc/refman/8.0/en/year.html
> "As of MySQL 8.0.19, specifying the number of digits for the YEAR data type is deprecated. In MySQL 8.0, YEAR(4) and YEAR have the same meaning; YEAR(2) is not supported."
**Source:** `sql/sql_yacc.yy:7155-7175`:
```
if (errno != 0 || length != 4)
{
  my_error(ER_INVALID_YEAR_COLUMN_LENGTH, MYF(0), "YEAR");
  MYSQL_YYABORT;
}
push_deprecated_warn(YYTHD, "YEAR(4)", "YEAR");
```
**Trigger:** `CREATE TABLE t (y YEAR(2))`, `CREATE TABLE t (y YEAR(4))`.
**Rule:**
- `YEAR(2)`, `YEAR(3)`, `YEAR(5)` → `ER_INVALID_YEAR_COLUMN_LENGTH` (parse-time abort).
- `YEAR(4)` → accepted with deprecation warning; normalized to `PT_year_type` (bare YEAR).
- `YEAR UNSIGNED` also deprecated (`ER_WARN_DEPRECATED_YEAR_UNSIGNED`).
- `SHOW CREATE TABLE` always emits `year`, never `year(4)` or `year(2)`.

**Observable via:**
- `CREATE TABLE t (y YEAR(2))` → error 1818.
- Round-trip of `YEAR(4)` through SDL → `YEAR`.

**omni pointer:** Needs verification that omni parser rejects `YEAR(2)` and normalizes `YEAR(4)` → `YEAR`.

---

### 16.12 TIMESTAMP first-column promotion inherits column's declared fsp
**Priority:** MED
**Status:** pending
**Cross-ref:** C3.1 (TIMESTAMP NOT NULL first-only promotion).
**Source:** `sql/sql_table.cc` (`promote_first_timestamp_column`); default clause generated as `CURRENT_TIMESTAMP(fsp)` matching the column's fsp.
**Trigger:** `CREATE TABLE t (ts TIMESTAMP(6) NOT NULL)`.
**Rule:**
- When first-TIMESTAMP promotion fires, the synthesized `DEFAULT CURRENT_TIMESTAMP` and `ON UPDATE CURRENT_TIMESTAMP` must carry the column's fsp to satisfy the PS5 fsp-match check automatically.
- Result: `TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)`.
- If `explicit_defaults_for_timestamp=ON`, no promotion happens and the column has no default (potential error on INSERT).

**Observable via:**
- `SHOW CREATE TABLE` after `CREATE TABLE t (ts TIMESTAMP(3) NOT NULL)` → shows `CURRENT_TIMESTAMP(3)` in both clauses.

**omni pointer:** `mysql/catalog/wt_3_3_test.go:538-543` tests C3.1 promotion but only for `TIMESTAMP` (fsp 0). A `TIMESTAMP(3)`/`TIMESTAMP(6)` variant is not covered. `mysql/catalog/show.go:289-292` already normalizes `CURRENT_TIMESTAMP(N)` on output — that path is preserved, but the promoter in `tablecmds.go` needs to construct the suffixed form.

---

## Section C17: String function charset / collation propagation

> **New section (Wave 3).** Every string-valued expression in MySQL carries
> both a charset and a "derivation coercibility" level. When a string
> function combines several string args (CONCAT, CONCAT_WS, REPLACE, REPEAT,
> LPAD, RPAD, SUBSTR, INSERT, LOWER, UPPER, TRIM, LEFT, RIGHT, ELT, FIELD,
> IF/IFNULL on strings, CASE on strings, UNION column merging, comparison,
> etc.) it must aggregate the input `DTCollation`s into a single result
> `DTCollation` using the rules encoded in `DTCollation::aggregate()` and
> the helper `agg_arg_charsets_for_string_result`.
>
> **Source anchors:**
> - Docs: https://dev.mysql.com/doc/refman/8.0/en/charset-collation-coercibility.html
> - Docs: https://dev.mysql.com/doc/refman/8.0/en/string-functions.html
> - MySQL 8.0: `sql/item.h` (`DTCollation`, `Derivation` enum, l.150-247)
> - MySQL 8.0: `sql/item.cc` (`DTCollation::aggregate`, l.2400-2510; `my_coll_agg_error`, l.2513)
> - MySQL 8.0: `sql/item_strfunc.cc` (`Item_func_concat::resolve_type` l.1109, `Item_func_concat_ws::resolve_type` l.1166, `Item_func_replace::resolve_type` l.1297, `Item_func_repeat::resolve_type` l.2581, `Item_func_rpad::resolve_type` l.2724, `Item_func_lpad::resolve_type` l.2822)
> - omni: `/Users/rebeliceyang/Github/omni/mysql/catalog/function_types.go` (no per-call collation derivation)
> - omni: `/Users/rebeliceyang/Github/omni/mysql/catalog/analyze_expr.go` (no collation tracking at all)
>
> **Derivation levels** (from `enum Derivation` in `sql/item.h` — lower
> number = "stronger", wins aggregation):
>
> | Level                 | Numeric | Source example                                     |
> | --------------------- | ------- | -------------------------------------------------- |
> | `DERIVATION_EXPLICIT` | 0       | `x COLLATE utf8mb4_bin`, `_utf8mb4'foo' COLLATE …` |
> | `DERIVATION_NONE`     | 1       | result of a prior failed aggregation               |
> | `DERIVATION_IMPLICIT` | 2       | table columns, view columns, SP vars               |
> | `DERIVATION_SYSCONST` | 3       | `USER()`, `VERSION()`, `DATABASE()`                |
> | `DERIVATION_COERCIBLE`| 4       | string literals, `_latin1'x'` without COLLATE      |
> | `DERIVATION_NUMERIC`  | 5       | numeric-to-string conversions                      |
> | `DERIVATION_IGNORABLE`| 6       | `NULL`, plain numeric NULL                         |
>
> **omni gap summary (applies to every scenario in this section):**
> `mysql/catalog/analyze_expr.go` and `mysql/catalog/function_types.go` do
> not track charset, collation, or coercibility on an expression. For every
> scenario below, omni will currently: parse the statement, assign the
> function call a `VARCHAR`/`TEXT` type, silently succeed on every
> illegal-mix case that real MySQL rejects, and record a view column type
> with no charset metadata, so Phase 3 view-column type derivation cannot
> reconcile with `information_schema.COLUMNS`. All scenarios in C17 are
> therefore **omni gaps** — the fix lands in Phase 3's `function_types.go`
> (new `charsetRule` / `collationRule` side-table) and in `analyze_expr.go`
> (new `exprCollation` helper that mirrors `DTCollation::aggregate`).

### 17.1 CONCAT of two columns with identical charset/collation

**Priority:** HIGH  **Status:** pending  **Anchor:** `{#c17-1}`

**MySQL source:** `Item_func_concat::resolve_type` (`sql/item_strfunc.cc`
l.1109) → `agg_arg_charsets_for_string_result(collation, args, arg_count)`.

**Rule:** Two IMPLICIT-level columns with the same charset and collation
aggregate to that same (charset, collation). Result derivation is IMPLICIT
(because the result is a function-of-columns, not a literal). Max char
length is the sum of arg char lengths.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
);
CREATE VIEW v1 AS SELECT CONCAT(a, b) AS c FROM t;

SELECT CHARACTER_SET_NAME, COLLATION_NAME, CHARACTER_MAXIMUM_LENGTH
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v1' AND COLUMN_NAME='c';
```

**Expected:** `utf8mb4`, `utf8mb4_0900_ai_ci`, `CHARACTER_MAXIMUM_LENGTH = 20`.

**omni status:** gap. analyze_expr.go produces a bare VARCHAR type for `c`
with no charset/collation, so a later SDL-diff against a real DB round-trip
cannot match.

---

### 17.2 CONCAT mixing latin1 + utf8mb4 columns (superset conversion)

**Priority:** HIGH  **Status:** pending  **Anchor:** `{#c17-2}`

**MySQL source:** `DTCollation::aggregate` `MY_COLL_ALLOW_SUPERSET_CONV`
branch (`sql/item.cc` l.2452). `agg_arg_charsets_for_string_result` always
passes `MY_COLL_ALLOW_SUPERSET_CONV | MY_COLL_ALLOW_COERCIBLE_CONV`.

**Rule:** ascii-repertoire latin1 values are a subset of utf8mb4. When two
IMPLICIT-derivation args have different charsets where one is a strict
superset, the superset charset wins with its default collation of the
superset side. The other arg is auto-converted at runtime via
`eval_string_arg`. Result derivation stays IMPLICIT.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  b VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci
);
CREATE VIEW v2 AS SELECT CONCAT(a, b) AS c FROM t;

SELECT CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v2' AND COLUMN_NAME='c';
```

**Expected:** `utf8mb4`, `utf8mb4_0900_ai_ci` (utf8mb4 wins as superset).
Note: this specific pair works because utf8mb4 is treated as the superset
of latin1 for `agg_arg_charsets_for_string_result`. Verify the exact
collation the oracle reports.

**omni status:** gap. Same as 17.1 plus the superset computation is absent
entirely.

---

### 17.3 CONCAT with incompatible collations (ER_CANT_AGGREGATE_2COLLATIONS)

**Priority:** HIGH  **Status:** pending  **Anchor:** `{#c17-3}`

**MySQL source:** `DTCollation::aggregate` fall-through at `sql/item.cc`
l.2466-2470 → set `DERIVATION_NONE`, then `agg_item_charsets` calls
`my_coll_agg_error` which raises `ER_CANT_AGGREGATE_2COLLATIONS`. Error
literal in `sql/item.cc` l.2515.

**Rule:** Two IMPLICIT columns with the same charset (`utf8mb4`) but *two
different non-binary collations* (e.g. `utf8mb4_0900_ai_ci` and
`utf8mb4_0900_as_cs`) have no aggregation path: neither is binary, neither
is the "explicit" level, superset check does not apply (same charset), and
the coercible-conversion branch is gated on `>= DERIVATION_SYSCONST` — so
aggregation fails. Result: view creation itself fails at resolve_type time.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
);
-- Expect: ERROR 1267 (HY000): Illegal mix of collations
--   (utf8mb4_0900_ai_ci,IMPLICIT) and (utf8mb4_0900_as_cs,IMPLICIT)
--   for operation 'concat'
CREATE VIEW v3 AS SELECT CONCAT(a, b) AS c FROM t;
```

**Expected:** `ERROR 1267` at `CREATE VIEW` time. View is NOT created.

**omni status:** **silent-accept gap**. omni has no derivation/aggregation
logic, so the view parses and analyzes cleanly and the analyzer happily
assigns `c` a VARCHAR type. This is the most user-visible bug in C17 —
bytebase query span, SDL diff, and Phase 3 view types all diverge from
reality on statements that MySQL outright rejects.

---

### 17.4 CONCAT_WS separator charset + NULL skipping

**Priority:** MEDIUM  **Status:** pending  **Anchor:** `{#c17-4}`

**MySQL source:** `Item_func_concat_ws::val_str` (`sql/item_strfunc.cc`
l.1133) and `Item_func_concat_ws::resolve_type` (l.1166). `resolve_type`
uses `agg_arg_charsets_for_string_result(collation, args, arg_count)` —
**the separator (arg 0) participates in aggregation**. `val_str` skips
NULL args (continue, not return NULL) — unlike CONCAT which returns NULL
if any arg is NULL.

**Rule:** Result charset is aggregated over **all** args including the
separator. If the separator is `_utf8mb4','` (COERCIBLE) and the payload
cols are IMPLICIT utf8mb4, the IMPLICIT cols win and the result is
utf8mb4/IMPLICIT. Result length is `(arg_count - 2) * sep_len + Σ payload`.
NULL among payload args does NOT make the output NULL.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
);
CREATE VIEW v4 AS SELECT CONCAT_WS(',', a, b, NULL) AS c FROM t;

SELECT CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v4' AND COLUMN_NAME='c';
```

**Expected:** `utf8mb4`, `utf8mb4_0900_ai_ci`. Also verify runtime: for
`INSERT INTO t VALUES (NULL, 'x')`, `SELECT CONCAT(a,b)` returns NULL but
`SELECT CONCAT_WS(',',a,b)` returns `'x'` (not NULL, not `',x'`).

**omni status:** gap. analyze_expr.go does not differentiate CONCAT vs
CONCAT_WS for NULL handling in its nullability inference either — a second
sub-gap worth recording.

---

### 17.5 String literal coercibility: `_utf8mb4'x'` vs bare `'x'`

**Priority:** MEDIUM  **Status:** pending  **Anchor:** `{#c17-5}`

**MySQL source:** `Item_string` constructors in `sql/item.h` l.5317-5490,
default `Derivation dv = DERIVATION_COERCIBLE`. An introducer `_utf8mb4'x'`
still constructs at COERCIBLE — the introducer only **pins the charset**,
it does NOT escalate to EXPLICIT. Only `COLLATE` escalates to EXPLICIT
(see 17.8).

**Rule:** A bare literal `'x'` is COERCIBLE with the session's
`character_set_connection`. `_utf8mb4'x'` is COERCIBLE with utf8mb4 no
matter the session charset. When aggregated against an IMPLICIT table
column of any other charset, the IMPLICIT column wins (its derivation is
stronger), so the literal is re-coded into the column's charset. This is
exactly why `WHERE col = 'x'` almost never fails with illegal mix — the
literal is coerced.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
SET NAMES utf8mb4;
CREATE TABLE t (a VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci);
-- Case A: bare literal in session utf8mb4 is still coerced to latin1
CREATE VIEW v5a AS SELECT CONCAT(a, 'x') AS c FROM t;
-- Case B: _utf8mb4 introducer still coerces down to latin1 (because
-- IMPLICIT beats COERCIBLE even across charsets when the coercible side
-- can be represented).
CREATE VIEW v5b AS SELECT CONCAT(a, _utf8mb4'x') AS c FROM t;

SELECT TABLE_NAME, COLUMN_NAME, CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME IN ('v5a','v5b') AND COLUMN_NAME='c';
```

**Expected:** Both views report `latin1` / `latin1_swedish_ci`. (The
IMPLICIT column wins over the COERCIBLE literal in both cases.)

**omni status:** gap. omni parses `_utf8mb4'x'` but stores no coercibility
metadata for literals.

---

### 17.6 REPEAT / LPAD / RPAD pass-through of first-arg charset

**Priority:** MEDIUM  **Status:** pending  **Anchor:** `{#c17-6}`

**MySQL source:**
- `Item_func_repeat::resolve_type` (`sql/item_strfunc.cc` l.2581) →
  `agg_arg_charsets_for_string_result(collation, args, 1)` — **only arg[0]**.
- `Item_func_lpad::resolve_type` (l.2822) and `Item_func_rpad::resolve_type`
  (l.2724) — same: aggregate over `args[0]` only, then
  `simplify_string_args(thd, collation, args+2, 1)` converts the pad string
  into the first arg's charset.

**Rule:** Result charset/collation of `REPEAT(s,n)`, `LPAD(s,n,pad)`,
`RPAD(s,n,pad)` is **always** taken from `s` (`args[0]`). The pad argument
can be any compatible charset; if not compatible it gets implicitly
converted, and if conversion fails at resolve-time you get
ER_CANT_AGGREGATE_2COLLATIONS from `simplify_string_args` — but the error
is reported against the pair `(args[0], args[2])`.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET latin1  COLLATE latin1_swedish_ci,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
);
CREATE VIEW v6a AS SELECT REPEAT(a, 3) AS c FROM t;              -- latin1
CREATE VIEW v6b AS SELECT LPAD(a, 20, b) AS c FROM t;            -- latin1 (arg0)
CREATE VIEW v6c AS SELECT RPAD(b, 20, a) AS c FROM t;            -- utf8mb4 (arg0)

SELECT TABLE_NAME, CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME IN ('v6a','v6b','v6c') AND COLUMN_NAME='c'
 ORDER BY TABLE_NAME;
```

**Expected:** `v6a`→latin1/latin1_swedish_ci,
`v6b`→latin1/latin1_swedish_ci, `v6c`→utf8mb4/utf8mb4_0900_ai_ci. Directly
demonstrates the "first-arg pins charset" asymmetry vs CONCAT (17.1/17.2).

**omni status:** gap. This is the closest-to-ready omni fix — add a
per-function rule "result.collation = args[0].collation" in
`function_types.go` for the first-arg-only family. Roughly a dozen builtins
belong to this family (LPAD, RPAD, REPEAT, REVERSE, TRIM, LTRIM, RTRIM,
LEFT, RIGHT, SUBSTR, SUBSTRING, SUBSTRING_INDEX, INSERT).

---

### 17.7 `CONVERT(x USING charset)` forces the charset

**Priority:** HIGH  **Status:** pending  **Anchor:** `{#c17-7}`

**MySQL source:** `Item_func_conv_charset` (`sql/item_strfunc.cc`, search
for `Item_func_conv_charset::resolve_type`). It pins `collation` to the
target charset's default collation at `DERIVATION_IMPLICIT` — **not**
EXPLICIT, and **not** COERCIBLE. That is, `CONVERT(col USING utf8mb4)` is
treated like another IMPLICIT column of utf8mb4 for downstream aggregation.

**Rule:** Use `CONVERT(... USING cs)` to sidestep illegal-mix by promoting
a foreign-charset value to a shared charset. Result derivation is IMPLICIT,
charset is `cs`, collation is `cs`'s default.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
);
-- 17.3 says this would fail:
--   CREATE VIEW bad AS SELECT CONCAT(a, b) FROM t; -- ER 1267
-- But CONVERT rescues it:
CREATE VIEW v7 AS SELECT CONCAT(a, CONVERT(b USING utf8mb4)) AS c FROM t;

SELECT CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v7' AND COLUMN_NAME='c';
```

**Expected:** `utf8mb4` / `utf8mb4_0900_ai_ci` (the default collation for
the `utf8mb4` charset on MySQL 8.0). VERIFY AT ORACLE: if MySQL actually
rejects this combination too, update the scenario to wrap BOTH sides in
CONVERT.

**omni status:** gap. omni parses `CONVERT ... USING ...` but the function
return type is a plain VARCHAR without charset.

---

### 17.8 `COLLATE` clause is EXPLICIT — highest precedence

**Priority:** HIGH  **Status:** pending  **Anchor:** `{#c17-8}`

**MySQL source:** `Item_func_set_collation` (in `sql/item_strfunc.cc`, grep
`Item_func_set_collation::resolve_type`). Pins `collation.derivation =
DERIVATION_EXPLICIT`. In `DTCollation::aggregate`, `DERIVATION_EXPLICIT` is
numeric 0 — it beats IMPLICIT, COERCIBLE, etc. If both sides are EXPLICIT
with different collations, the aggregation fails (see l.2479-2481).

**Rule:** `x COLLATE utf8mb4_bin` pins that expression to
EXPLICIT/utf8mb4_bin. In any aggregation with non-EXPLICIT args, the
EXPLICIT side's collation wins. Two EXPLICIT args with different collations
is a hard error.

**Oracle fixture:**
```sql
DROP DATABASE IF EXISTS testdb; CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  b VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs
);
-- Case A: one side EXPLICIT — wins.
CREATE VIEW v8a AS SELECT CONCAT(a, b COLLATE utf8mb4_bin) AS c FROM t;
-- Case B: both sides EXPLICIT with different collations — must fail.
--   Expect: ERROR 1270 (HY000): Illegal mix of collations
--   (utf8mb4_bin,EXPLICIT) and (utf8mb4_0900_ai_ci,EXPLICIT) for operation 'concat'
CREATE VIEW v8b AS SELECT CONCAT(a COLLATE utf8mb4_0900_ai_ci,
                                 b COLLATE utf8mb4_bin) AS c FROM t;

SELECT CHARACTER_SET_NAME, COLLATION_NAME
  FROM information_schema.COLUMNS
 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v8a' AND COLUMN_NAME='c';
```

**Expected:**
- `v8a` created; column `c` has `utf8mb4` / `utf8mb4_bin`.
- `v8b` creation fails with ER_CANT_AGGREGATE_2COLLATIONS (error 1270 in
  some code paths, 1267 in others — oracle-verify to pin the exact error).

**omni status:** gap. omni parses `COLLATE` (accepted as a suffix in
expressions) but does not track EXPLICIT derivation. `v8b` analyzes cleanly
today — a soft-silent failure mirroring 17.3.

---

## Section C18: SHOW CREATE TABLE elision rules

> **Expansion note (Wave 2):** grew from 6 to 15 scenarios via systematic
> walk of `sql/sql_show.cc::store_create_info` (approx L1800-2600),
> `sql/sql_show.cc::store_key_options` (approx L2620-2680), and the
> `HA_CREATE_USED_*` flag table in `sql/handler.h` (L710-805). Every
> `if (create_info_arg->used_fields & HA_CREATE_USED_*)` gate and every
> `if (share->...)` renderer in that function is a potential elision
> rule; each one where MySQL hides a technically-present value is
> lifted to a scenario below. This section is the deparse contract —
> omni's `mysql/catalog/show.go` must match byte-for-byte on every
> elision rule documented here or SDL round-trip breaks.

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

**Priority:** HIGH  **Status:** verified

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

**Priority:** HIGH  **Status:** verified

**See catalog:** C18.2.

---

### 18.3 ENGINE clause always rendered (implicit default shown)

**Setup:**
```sql
CREATE TABLE t (a INT);
```

**Oracle verification:** `SHOW CREATE TABLE t` contains `ENGINE=InnoDB`.

**omni assertion:** deparse emits `ENGINE=InnoDB`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.3 (in practice always rendered for 8.0.45).

---

### 18.4 AUTO_INCREMENT clause elided when counter == 1

**Setup:**
```sql
CREATE TABLE t (id INT AUTO_INCREMENT PRIMARY KEY);
```

**Oracle verification:** `SHOW CREATE TABLE t` does NOT contain `AUTO_INCREMENT=`.

**omni assertion:** deparse output contains no `AUTO_INCREMENT=` clause.

**Priority:** HIGH  **Status:** verified

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

**Priority:** HIGH  **Status:** verified

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

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.6.

---

### 18.7 Table-level COLLATE clause rendered only when non-primary or utf8mb4_0900_ai_ci

**Setup:**
```sql
-- Primary collation for latin1 is latin1_swedish_ci.
CREATE TABLE t_prim   (x INT) CHARACTER SET latin1;
-- Non-primary collation for latin1.
CREATE TABLE t_nonprim (x INT) CHARACTER SET latin1 COLLATE latin1_bin;
-- utf8mb4_0900_ai_ci is "primary-ish" but always rendered (special-cased).
CREATE TABLE t_0900    (x INT) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_prim;     -- contains DEFAULT CHARSET=latin1; NO COLLATE= clause
SHOW CREATE TABLE t_nonprim;  -- contains DEFAULT CHARSET=latin1 COLLATE=latin1_bin
SHOW CREATE TABLE t_0900;     -- contains DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci (special-case)
```

Expected substrings:
- `t_prim`: contains `DEFAULT CHARSET=latin1`; does NOT contain `COLLATE=`.
- `t_nonprim`: contains `COLLATE=latin1_bin`.
- `t_0900`: contains `COLLATE=utf8mb4_0900_ai_ci`.

**omni assertion:** deparse emits `COLLATE=` after `DEFAULT CHARSET=` iff the collation is non-primary for the charset OR equals `utf8mb4_0900_ai_ci`. Source: `sql_show.cc:2450-2456`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.7 (`sql/sql_show.cc:2450-2456`).

---

### 18.8 Table-level KEY_BLOCK_SIZE elided unless explicitly set

**Setup:**
```sql
CREATE TABLE t_nokbs  (a INT);
CREATE TABLE t_kbs    (a INT) KEY_BLOCK_SIZE=4;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_nokbs;  -- does NOT contain KEY_BLOCK_SIZE=
SHOW CREATE TABLE t_kbs;    -- contains KEY_BLOCK_SIZE=4
```

**omni assertion:** deparse emits `KEY_BLOCK_SIZE=N` only when `share->key_block_size != 0`. Source: `sql_show.cc:2516-2520` (renders iff `table->s->key_block_size` truthy; unset = 0 = elided).

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.8 (`sql/sql_show.cc:2516-2520`; `HA_CREATE_USED_KEY_BLOCK_SIZE`).

---

### 18.9 COMPRESSION clause elided unless explicitly set

**Setup:**
```sql
CREATE TABLE t_nocomp (a INT);
CREATE TABLE t_comp   (a INT) COMPRESSION='ZLIB';
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_nocomp;  -- does NOT contain COMPRESSION=
SHOW CREATE TABLE t_comp;    -- contains COMPRESSION='ZLIB'
```

**omni assertion:** deparse emits `COMPRESSION='...'` only when `share->compress.length > 0`. Omni must elide it even though InnoDB's effective default is "no compression" (storage default; not shown). Source: `sql_show.cc:2522-2525`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.9 (`sql/sql_show.cc:2522-2525`; `HA_CREATE_USED_COMPRESS`).

---

### 18.10 STATS_PERSISTENT / STATS_AUTO_RECALC / STATS_SAMPLE_PAGES elision

**Setup:**
```sql
CREATE TABLE t_nostats (a INT);
CREATE TABLE t_stats   (a INT)
  STATS_PERSISTENT=1
  STATS_AUTO_RECALC=0
  STATS_SAMPLE_PAGES=32;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_nostats;
-- Must NOT contain STATS_PERSISTENT=, STATS_AUTO_RECALC=, or STATS_SAMPLE_PAGES=.
SHOW CREATE TABLE t_stats;
-- Must contain STATS_PERSISTENT=1, STATS_AUTO_RECALC=0, STATS_SAMPLE_PAGES=32.
```

**omni assertion:** deparse emits each of these three clauses only when the corresponding field has a non-default value:
- `STATS_PERSISTENT=1` iff `HA_OPTION_STATS_PERSISTENT` set; `=0` iff `HA_OPTION_NO_STATS_PERSISTENT` set; otherwise elided (DEFAULT).
- `STATS_AUTO_RECALC=1|0` iff `share->stats_auto_recalc` is ON/OFF; elided when DEFAULT.
- `STATS_SAMPLE_PAGES=N` iff `share->stats_sample_pages != 0`.

Source: `sql_show.cc:2481-2497`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.10 (`sql/sql_show.cc:2481-2497`; `HA_CREATE_USED_STATS_PERSISTENT` / `_STATS_AUTO_RECALC` / `_STATS_SAMPLE_PAGES`).

---

### 18.11 AVG_ROW_LENGTH / MAX_ROWS / MIN_ROWS elided when zero

**Setup:**
```sql
CREATE TABLE t_nominmax (a INT);
CREATE TABLE t_minmax   (a INT) MIN_ROWS=10 MAX_ROWS=1000 AVG_ROW_LENGTH=256;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_nominmax;
-- Must NOT contain MIN_ROWS=, MAX_ROWS=, or AVG_ROW_LENGTH=.
SHOW CREATE TABLE t_minmax;
-- Must contain MIN_ROWS=10, MAX_ROWS=1000, AVG_ROW_LENGTH=256.
```

**omni assertion:** deparse emits each clause only when the `share->min_rows` / `share->max_rows` / `share->avg_row_length` field is truthy (non-zero). Source: `sql_show.cc:2461-2479`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.11 (`sql/sql_show.cc:2461-2479`; `HA_CREATE_USED_MIN_ROWS` / `_MAX_ROWS` / `_AVG_ROW_LENGTH`).

---

### 18.12 TABLESPACE clause elision rules

**Setup:**
```sql
CREATE TABLE t_default (a INT);  -- no explicit tablespace
CREATE TABLE t_gts     (a INT) TABLESPACE=innodb_system;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_default;  -- does NOT contain TABLESPACE=
SHOW CREATE TABLE t_gts;      -- contains /*!50100 TABLESPACE `innodb_system` */
```

**omni assertion:** deparse emits a `TABLESPACE=` clause only when `HA_CREATE_USED_TABLESPACE` was set at create time — i.e. when the table was explicitly attached to a named (non-implicit) tablespace. The implicit per-file tablespace InnoDB creates for every file-per-table is NEVER rendered. Wrapped in `/*!50100 ... */` versioned comment on output.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.12 (`sql/sql_show.cc` tablespace path via `handler::append_create_info`; `HA_CREATE_USED_TABLESPACE`).

---

### 18.13 PACK_KEYS / CHECKSUM / DELAY_KEY_WRITE elision

**Setup:**
```sql
CREATE TABLE t_none (a INT);
CREATE TABLE t_opts (a INT) PACK_KEYS=1 CHECKSUM=1 DELAY_KEY_WRITE=1;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t_none;
-- Must NOT contain PACK_KEYS=, CHECKSUM=, or DELAY_KEY_WRITE=.
SHOW CREATE TABLE t_opts;
-- Must contain PACK_KEYS=1, CHECKSUM=1, DELAY_KEY_WRITE=1.
-- (InnoDB may reject or normalize some values; the elision contract is
--  what matters: none of these is ever shown unless the caller asked for it.)
```

**omni assertion:** deparse emits each clause only when the corresponding `HA_OPTION_*` bit is set on `share->db_create_options`:
- `PACK_KEYS=1` iff `HA_OPTION_PACK_KEYS`; `=0` iff `HA_OPTION_NO_PACK_KEYS`; else elided.
- `CHECKSUM=1` iff `HA_OPTION_CHECKSUM`.
- `DELAY_KEY_WRITE=1` iff `HA_OPTION_DELAY_KEY_WRITE`.

Source: `sql_show.cc:2481-2504`.

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.13 (`sql/sql_show.cc:2481-2504`; `HA_CREATE_USED_PACK_KEYS` / `_CHECKSUM` / `_DELAY_KEY_WRITE`).

---

### 18.14 Per-index COMMENT and KEY_BLOCK_SIZE rendering inside index clauses

**Setup:**
```sql
CREATE TABLE t (
    id INT PRIMARY KEY,
    a  INT,
    b  INT,
    KEY ix_plain (a),
    KEY ix_cmt   (b) COMMENT 'hello'
) KEY_BLOCK_SIZE=4;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected substrings:
- `KEY \`ix_plain\` (\`a\`)` — no `COMMENT`, no per-index `KEY_BLOCK_SIZE`.
- `KEY \`ix_cmt\` (\`b\`) COMMENT 'hello'` — per-index COMMENT appears iff `HA_USES_COMMENT` and length > 0.
- Table-level `KEY_BLOCK_SIZE=4` present; neither index prints its own `KEY_BLOCK_SIZE` (both inherit the table-level value — per-index KBS only renders when it DIFFERS from table-level).

**omni assertion:** deparse emits per-index `COMMENT '...'` iff the index has a non-empty comment, and per-index `KEY_BLOCK_SIZE=N` iff `key_info->block_size != table->s->key_block_size`. Source: `sql_show.cc:2646-2665` (`store_key_options`).

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.14 (`sql/sql_show.cc:2646-2665`).

---

### 18.15 USING BTREE/HASH clause rendered only when algorithm explicit

**Setup:**
```sql
CREATE TABLE t (
    id INT,
    a  INT,
    b  INT,
    KEY ix_default (a),
    KEY ix_btree   (a) USING BTREE,
    KEY ix_hash    (b) USING HASH
) ENGINE=InnoDB;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected substrings:
- `KEY \`ix_default\` (\`a\`)` — no `USING` clause (InnoDB silently stores BTREE, but MySQL elides the USING because the user didn't specify it).
- `KEY \`ix_btree\` (\`a\`) USING BTREE` — explicit flag preserved and rendered.
- InnoDB rewrites `USING HASH` to `USING BTREE` under the hood; the rendered form depends on what `key_info->algorithm` holds after the handler fixes up HASH → BTREE. Test asserts only that when `is_algorithm_explicit == true`, a `USING ...` clause appears (and does not appear when it is false).

**omni assertion:** deparse emits `USING BTREE|HASH|RTREE` after index-key-list only when the index was created with an explicit algorithm flag (omni must track an `AlgorithmExplicit` bit on index metadata). Default (implicit) indexes must NEVER produce `USING`. Source: `sql_show.cc:2624-2642` (`if (key_info->is_algorithm_explicit)`).

**Priority:** HIGH  **Status:** pending

**See catalog:** C18.15 (`sql/sql_show.cc:2624-2642`).

---

## Section C19: Virtual / functional indexes

> **New section (Wave 3).** MySQL 8.0.13 introduced **functional key parts**
> — indexes whose key parts are expressions rather than plain columns, e.g.
> `INDEX ((LOWER(name)))`. MySQL does not have a native "expression index"
> storage model: at DDL time, the server synthesizes a **hidden VIRTUAL
> generated column** over the expression, gives it a generated name (see
> §1.11/1.12), marks it `HT_HIDDEN_SQL`, and builds an ordinary index over
> that column. All downstream semantics flow from this rewrite: type
> inference, `SELECT *` visibility, disallowed function classes, BLOB
> restrictions, drop/rename cleanup, and replication ordering.
>
> omni now has catalog-level support for functional indexes: the catalog
> synthesizes system-hidden virtual generated columns, keeps functional key
> part expressions for `SHOW CREATE TABLE`, validates covered invalid
> expression classes, and handles hidden-column drop/rename lifecycle. Paired
> with §1.11 (auto-name `functional_index[_N]`) and §1.12 (hidden column name
> `!hidden!{key}!{part}!{count}`), C19 captures the remaining semantic
> obligations for schema sync and query-span resolution to match MySQL 8.0.
>
> **Source anchors (MySQL 8.0 `sql/sql_table.cc`):**
> `add_functional_index_to_create_list` (L7783-L7900),
> `make_functional_index_column_name` (L7710-L7743),
> `Replace_field_processor_arg::replace_field_processor` (L7516-L7550),
> `handle_drop_functional_index` (L16158-L16195),
> `handle_rename_functional_index` (L16211+).
> **Docs:** https://dev.mysql.com/doc/refman/8.0/en/create-index.html#create-index-functional-key-parts

### 19.1 Functional index creates a hidden VIRTUAL generated column

**Priority:** P0  **Status:** verified  **Anchor:** `{#c19-1}`

```sql
CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));
CREATE INDEX idx_lower ON t ((LOWER(name)));

-- User-visible SHOW COLUMNS does NOT reveal the hidden column:
SHOW COLUMNS FROM t;
-- +-------+-------------+------+-----+---------+-------+
-- | id    | int         | NO   | PRI | NULL    |       |
-- | name  | varchar(64) | YES  |     | NULL    |       |
-- +-------+-------------+------+-----+---------+-------+

-- But information_schema.COLUMNS shows nothing extra either, because
-- HT_HIDDEN_SQL columns are filtered out of I_S.COLUMNS for regular users.
-- Evidence of the hidden column is instead visible via SHOW CREATE TABLE
-- (which renders the functional key part as the expression) and via
-- information_schema.STATISTICS referencing a generated EXPRESSION:
SELECT INDEX_NAME, COLUMN_NAME, EXPRESSION
  FROM information_schema.STATISTICS
 WHERE TABLE_NAME='t' AND INDEX_NAME='idx_lower';
-- +-----------+-------------+--------------+
-- | idx_lower | NULL        | lower(`name`)|
-- +-----------+-------------+--------------+

SHOW CREATE TABLE t;
-- ... KEY `idx_lower` ((lower(`name`))) ...
```

**Expected:** The catalog holds an extra column with `hidden =
HT_HIDDEN_SQL`, `stored_in_db = false` (VIRTUAL), `gcol_info.expr_item =
LOWER(name)`, auto-named `!hidden!idx_lower!0!0`. `SHOW COLUMNS` and
user-scoped `I_S.COLUMNS` suppress it. `information_schema.STATISTICS`
surfaces it as `EXPRESSION` on the key part with `COLUMN_NAME = NULL`.

**omni assertion:** after loading the table, the catalog must expose both
the hidden column (queryable via an internal API, e.g. `tbl.HiddenColumns()`)
and an index with exactly one key part that points to that hidden column OR
carries an expression payload. `SHOW CREATE TABLE` deparse must render the
functional key part form `((expr))`, not `(hidden_col_name)`.

**MySQL source:**
- `sql/sql_table.cc:7883` — `cr->hidden = dd::Column::enum_hidden_type::HT_HIDDEN_SQL;`
- `sql/sql_table.cc:7884` — `cr->stored_in_db = false;`
- `sql/sql_table.cc:7887-7891` — `Value_generator` built with `set_field_stored(false)`.

**omni status 2026-04-24:** verified by `TestScenario_C19/19_1`. The catalog
uses `ColumnHiddenSystem` for synthesized virtual generated columns and
`SHOW CREATE TABLE` renders the functional key part expression.

---

### 19.2 Hidden column type is inferred from the expression return type

**Priority:** P0  **Status:** verified  **Anchor:** `{#c19-2}`

```sql
CREATE TABLE t (
  a INT, b INT,
  name VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  payload JSON,
  INDEX k_sum ((a + b)),                       -- hidden col: BIGINT
  INDEX k_low ((LOWER(name))),                 -- hidden col: VARCHAR(64) utf8mb4_0900_ai_ci
  INDEX k_cast ((CAST(payload->'$.age' AS UNSIGNED)))   -- hidden col: BIGINT UNSIGNED
);
```

**Expected:** MySQL resolves the expression (`kp->resolve_expression(thd)`)
and calls `generate_create_field(thd, item, &tmp_table)` to materialize a
`Create_field` whose `sql_type`, `length`, `charset`, `collation`, `flags`
(nullability, UNSIGNED, …) come from the `Item*`. So the hidden column’s
type mirrors what you’d get from `SELECT expr FROM t` type-inference.

`information_schema.STATISTICS.COLLATION` and `SHOW INDEX FROM t` will
reflect those inferred properties, and server-side query planning uses them
for index matching.

**omni assertion:** the synthesized hidden column’s `DataType` must be the
result of evaluating the expression type against the table’s other columns.
Importantly, the hidden column inherits the **expression collation**, not
the table default. (E.g. `((name))` is rejected by
ER_FUNCTIONAL_INDEX_ON_FIELD, but `((CONVERT(name USING utf8mb4) COLLATE
utf8mb4_bin))` produces a column with `utf8mb4_bin`, even if `name` is
`utf8mb4_0900_ai_ci`.)

**MySQL source:**
- `sql/sql_table.cc:7860` — `Create_field *cr = generate_create_field(thd, item, &tmp_table);`
- `sql/sql_table.cc:7864-7868` — `if (is_blob(cr->sql_type)) ER_FUNCTIONAL_INDEX_ON_LOB`.
- `sql/sql_table.cc:7889` — `gcol_info->set_field_type(cr->sql_type);`

**omni status 2026-04-24:** verified by `TestScenario_C19/19_2` for the
covered resolver surface: integer arithmetic, string functions that preserve
column string metadata, explicit CAST targets, and JSON operator return types
needed by validation.

---

### 19.3 Hidden functional column is invisible to `SELECT *` and user I_S.COLUMNS

**Priority:** P0  **Status:** verified  **Anchor:** `{#c19-3}`

```sql
CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));
CREATE INDEX idx_lower ON t ((LOWER(name)));
INSERT INTO t VALUES (1,'Alice'),(2,'BOB');

-- (a) SELECT * skips the hidden column entirely.
SELECT * FROM t;
-- +----+-------+
-- | id | name  |
-- | 1  | Alice |
-- | 2  | BOB   |
-- +----+-------+
-- (no third column)

-- (b) Hidden column cannot be referenced by name either.
SELECT `!hidden!idx_lower!0!0` FROM t;
-- ERROR 1054 (42S22): Unknown column '!hidden!idx_lower!0!0' in 'field list'

-- (c) information_schema.COLUMNS as a non-DBA user does NOT list it.
SELECT COLUMN_NAME FROM information_schema.COLUMNS
 WHERE TABLE_NAME='t';
-- id, name    (no hidden column row)

-- (d) But the key-part expression IS visible via I_S.STATISTICS:
SELECT INDEX_NAME, SEQ_IN_INDEX, EXPRESSION
  FROM information_schema.STATISTICS
 WHERE TABLE_NAME='t';
```

**Expected:** `HT_HIDDEN_SQL` hides the column from the SQL namespace
entirely — not just from `SELECT *` expansion, but also from direct name
reference, `INSERT ... VALUES` positional binding, and the user-scoped
`I_S.COLUMNS` view. Only server internals and `I_S.STATISTICS.EXPRESSION`
acknowledge its existence.

**omni assertion:** query-span and completion must treat functional-index
hidden columns as non-resolvable identifiers. `SELECT *` expansion in the
advisor must not include them. Schema diffing must not treat the hidden
column as a user-drift column that needs to be added/dropped explicitly —
it is an artifact of the index.

**MySQL source:**
- `sql/sql_base.cc` — `setup_wild` skips `HT_HIDDEN_SQL` fields.
- `sql/sql_show.cc` — `I_S.COLUMNS` DD view joins on `dd.columns.is_hidden=0` for the user view.
- `sql/item.cc` — `find_field_in_table` rejects hidden-by-system fields unless `thd->lex->allow_sum_func` / DD access path.

**omni status 2026-04-24:** verified by `TestScenario_C19/19_3`. Hidden
functional columns are retained internally but filtered from user-visible
catalog surfaces and analyzer scope.

---

### 19.4 Functional expression must be deterministic and pure

**Priority:** P1  **Status:** verified  **Anchor:** `{#c19-4}`

```sql
-- Rejected: uses non-deterministic function
CREATE TABLE t (a INT, INDEX ((a + RAND())));
-- ERROR 3757 (HY000): Expression of functional index 'functional_index'
--   contains a disallowed function.

-- Rejected: references subquery, stored function, or variable
CREATE TABLE t (a INT, INDEX ((a + @@global.long_query_time)));
-- ERROR 3757 (HY000): ... disallowed function.

-- Rejected: DEFAULT() / VALUES() — replace_field_processor path
CREATE TABLE t (a INT, b INT DEFAULT 5, INDEX ((DEFAULT(b) + a)));
-- ERROR 3757 (HY000): Expression of functional index 'functional_index'
--   contains a disallowed function.

-- Rejected: bare field reference — not actually an expression
CREATE TABLE t (a INT, INDEX ((a)));
-- ERROR 3756 (HY000): Functional index on a column is not supported.
--   Consider using a regular index instead.

-- Rejected: result type is BLOB/TEXT
CREATE TABLE t (payload JSON, INDEX ((payload->'$.name')));
-- ERROR 3754 (HY000): Cannot create a functional index on an expression
--   that returns a BLOB or TEXT. Please consider using CAST.
```

**Expected:** The set of rejections mirrors generated-column rules
(`pre_validate_value_generator_expr` with `VGS_GENERATED_COLUMN`) plus three
functional-index-specific errors:
- `ER_FUNCTIONAL_INDEX_ON_FIELD` (3756) — bare column, not an expression.
- `ER_FUNCTIONAL_INDEX_ON_LOB` (3754) — inferred type is BLOB/TEXT/JSON.
- `ER_FUNCTIONAL_INDEX_PRIMARY_KEY` (3752) — declared as PRIMARY KEY.
- `ER_FUNCTIONAL_INDEX_FUNCTION_IS_NOT_ALLOWED` (3757) — disallowed Item
  (non-deterministic, session state, subquery, DEFAULT(), VALUES(), etc.).

**omni assertion:** the omni catalog-builder (and eventually schema-diff /
advisor) must reproduce each of these rejections at load/validation time so
that invalid SDL is caught before applying. Error code & message should
match MySQL where practical.

**MySQL source:**
- `sql/sql_table.cc:7788` — `ER_FUNCTIONAL_INDEX_PRIMARY_KEY`
- `sql/sql_table.cc:7828` — `ER_FUNCTIONAL_INDEX_ON_FIELD`
- `sql/sql_table.cc:7871` — `ER_FUNCTIONAL_INDEX_ON_LOB`
- `sql/sql_table.cc:7529` — `ER_FUNCTIONAL_INDEX_FUNCTION_IS_NOT_ALLOWED`
- `sql/sql_table.cc:7830-7832` — `pre_validate_value_generator_expr(..., VGS_GENERATED_COLUMN)`

**omni status 2026-04-24:** verified by `TestScenario_C19/19_4` for covered
invalid classes: bare field functional key parts, disallowed functions and
stateful constructs, and LOB/JSON/TEXT-returning expressions without CAST.

---

### 19.5 Functional index on JSON path via `(col->>'$.path')`

**Priority:** P0  **Status:** verified  **Anchor:** `{#c19-5}`

```sql
CREATE TABLE t (
  id INT PRIMARY KEY,
  doc JSON,
  -- ->> returns TEXT, which is BLOB-family → must CAST down to a
  -- non-LOB type, otherwise ER_FUNCTIONAL_INDEX_ON_LOB.
  INDEX idx_name ((CAST(doc->>'$.name' AS CHAR(64))))
);

-- Equivalent via JSON_UNQUOTE(JSON_EXTRACT(...)):
CREATE INDEX idx_age ON t ((CAST(doc->'$.age' AS UNSIGNED)));

INSERT INTO t VALUES
  (1, '{"name":"Alice","age":30}'),
  (2, '{"name":"Bob","age":25}');

-- The optimizer uses the functional index for matching predicates:
EXPLAIN SELECT * FROM t WHERE CAST(doc->>'$.name' AS CHAR(64)) = 'Alice';
-- key: idx_name, type: ref

-- Common mistakes:
CREATE INDEX bad ON t ((doc->>'$.name'));
-- ERROR 3754 (HY000): ... BLOB or TEXT. Please consider using CAST.

CREATE INDEX bad2 ON t ((doc->'$.name'));
-- ERROR 3754 (HY000): ... JSON value → treated as LOB.
```

**Expected:** The single most common real-world use of functional indexes
is indexing a JSON path. MySQL forces the user to `CAST` because `->>`
returns `longtext` and `->` returns `json`, both of which are LOB-family.
Only after the cast does `generate_create_field` produce an acceptable
hidden column type.

For the optimizer to use the index, the query predicate must use **the
exact same expression** as the index definition (modulo commutativity);
MySQL does not normalize `doc->>'$.name'` and
`JSON_UNQUOTE(JSON_EXTRACT(doc,'$.name'))` for index-matching purposes
even though they are semantically equivalent.

**omni assertion:** the parser must accept the JSON operator forms inside
functional key parts. The catalog must recognize `CAST(... AS CHAR(N))` and
`CAST(... AS UNSIGNED/SIGNED/DECIMAL/DATE/DATETIME)` as valid cast targets,
and the advisor/query-span engine should match functional indexes against
predicates that use the identical expression tree.

**MySQL source:**
- `sql/sql_table.cc:7864` — `is_blob(cr->sql_type)` check.
- `sql/item_json_func.cc` — `Item_func_json_extract_oneline` returns JSON.
- `sql/sql_optimizer.cc` — functional-index matching in `substitute_gc()`
  / `find_func_index_on_expr()`.

**omni status 2026-04-24:** verified by `TestScenario_C19/19_5`. The catalog
rejects uncast JSON/LOB expressions and normalizes JSON path literals with
the `_utf8mb4` introducer in functional-index `SHOW CREATE TABLE` output for
the covered case.

---

### 19.6 DROP INDEX cascades to the hidden generated column

**Priority:** P0  **Status:** verified  **Anchor:** `{#c19-6}`

```sql
CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64));
CREATE INDEX idx_lower ON t ((LOWER(name)));

-- Dropping the index also drops the hidden generated column.
DROP INDEX idx_lower ON t;
SHOW CREATE TABLE t;
-- CREATE TABLE `t` (
--   `id` int NOT NULL,
--   `name` varchar(64) DEFAULT NULL,
--   PRIMARY KEY (`id`)
-- )   -- no trace of the hidden column

-- You CANNOT drop the hidden column directly.
CREATE INDEX idx_lower ON t ((LOWER(name)));
ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`;
-- ERROR 3108 (HY000): Cannot drop column '!hidden!idx_lower!0!0' because
--   it is used by a functional index. In order to drop the column, the
--   functional index must be removed first.

-- RENAME INDEX also cascades: the hidden column is renamed to match the
-- new index name because the name is derived from the key name.
ALTER TABLE t RENAME INDEX idx_lower TO idx_lc;
-- Internally: !hidden!idx_lower!0!0 → !hidden!idx_lc!0!0
```

**Expected:** `handle_drop_functional_index` walks the `drop_list`; for each
`Alter_drop::KEY` whose key parts reference `is_field_for_functional_index()`
columns, it appends synthetic `Alter_drop::COLUMN` entries. The inverse is
blocked: dropping the hidden column by name yields
`ER_CANNOT_DROP_COLUMN_FUNCTIONAL_INDEX` (3108). `ALTER TABLE ... RENAME
INDEX` cascades to a hidden-column rename because the hidden column name is
derived from the key name (see §1.12).

**omni assertion:** omni’s catalog model and schema-diff migration generator must:
1. Emit `DROP INDEX idx_name` *without* an accompanying `DROP COLUMN` for
   the hidden column (server handles cascade).
2. Reject any attempt to generate `ALTER TABLE ... DROP COLUMN
   '!hidden!...'` as a planned migration — this must be caught as a
   pre-flight validation error.
3. When diffing a rename (old SDL `idx_a` → new SDL `idx_b`, same
   expression), prefer `ALTER TABLE ... RENAME INDEX idx_a TO idx_b` over
   drop+recreate so the hidden column is renamed in place.

**omni status 2026-04-24:** catalog behavior is verified by
`TestScenario_C19/19_6`. `DROP INDEX` removes the attached system-hidden
generated column, direct `DROP COLUMN` on that column returns a
3108-equivalent functional-index error, and `RENAME INDEX` renames the
attached hidden column. The migration-generator diff preference remains a
separate integration concern.

**MySQL source:**
- `sql/sql_table.cc:16158-16195` — `handle_drop_functional_index`
- `sql/sql_table.cc:16187` — `ER_CANNOT_DROP_COLUMN_FUNCTIONAL_INDEX`
- `sql/sql_table.cc:16211+` — `handle_rename_functional_index`

**omni gap:** the catalog lifecycle is implemented. The migration generator
still needs a pre-diff normalization pass that attaches hidden columns to
their parent index so it does not plan invalid hidden-column DDL.

---

## Section C20: Field-type-specific implicit defaults

> **New section (Wave 3).** C2 (type normalization) is about how the parser
> rewrites type names (REAL → DOUBLE, SERIAL → BIGINT UNSIGNED NOT NULL
> AUTO_INCREMENT UNIQUE). C20 is about the **runtime implicit default value**
> applied when a column has no explicit `DEFAULT` clause and an INSERT either
> omits the column or writes `DEFAULT` to it. Two distinct layers, both
> load-bearing for catalog and deparse.
>
> **Primary doc:** https://dev.mysql.com/doc/refman/8.0/en/data-type-defaults.html
> **Primary source:** `sql/field.cc` (per-Field subclass `store_default` /
> `reset`), `sql/sql_data_change.cc`, `sql/sql_insert.cc` (implicit default
> fill-in), `sql/sql_table.cc` (BLOB/JSON/GEOMETRY expression-default rules),
> `sql/field.h` (`m_default_val_expr`).
>
> **Why omni cares:**
> 1. `SHOW CREATE TABLE` round-trip: omni's deparser must render `DEFAULT <x>`
>    iff the user wrote an explicit default. A column with no explicit default
>    must NOT render as `DEFAULT 0` / `DEFAULT ''` even though the runtime
>    would materialize those values.
> 2. SDL diff semantics: `col INT NOT NULL` vs `col INT NOT NULL DEFAULT 0`
>    are **different** at the catalog level even if the runtime fallback is
>    identical.
> 3. Expression defaults (8.0.13+) are a separate AST category — the parser
>    stores a parenthesized `Expr` rather than a literal.

### 20.1 Integer column, NOT NULL, no DEFAULT — implicit 0 on INSERT

**Priority:** HIGH  **Status:** pending

```sql
CREATE TABLE c20_1 (id INT NOT NULL);
INSERT INTO c20_1 () VALUES ();           -- non-strict: inserts 0
INSERT INTO c20_1 () VALUES ();           -- strict: ER_NO_DEFAULT_FOR_FIELD (1364)
SHOW CREATE TABLE c20_1;
```

**Expected catalog state:**
- `Column{Name:"id", Type:INT, NotNull:true, Default:nil, DefaultAnalyzed:nil}`
- `Default` pointer **remains nil** — there is no stored default expression.
- The implicit `0` is a runtime property of `Field_long::reset()` (field.cc),
  not a catalog-level default.

**Expected SHOW CREATE TABLE:**
```
`id` int NOT NULL
```
(No `DEFAULT 0` rendered. Oracle MySQL 8.0 behaves identically.)

**Expected INSERT behavior (oracle):**
- `sql_mode` without `STRICT_ALL_TABLES`/`STRICT_TRANS_TABLES`: row inserted
  with `id=0`, warning `1364 Field 'id' doesn't have a default value`.
- `sql_mode` with strict: error `1364`, no row inserted.

**Oracle verify:** `SHOW CREATE TABLE c20_1` + both strict and non-strict
INSERT; compare column rendering and post-INSERT `SELECT *`.

---

### 20.2 Integer column, nullable, no DEFAULT — implicit NULL

**Priority:** HIGH  **Status:** pending

```sql
CREATE TABLE c20_2 (id INT);
INSERT INTO c20_2 () VALUES ();    -- always succeeds: id=NULL
SHOW CREATE TABLE c20_2;
```

**Expected catalog state:**
- `Column{Name:"id", Type:INT, NotNull:false, Default:nil}`.
- Note: MySQL's own `SHOW CREATE TABLE` for a nullable int with no default
  renders `` `id` int DEFAULT NULL ``. This is a **deparse hint** — the
  catalog did not parse an explicit `DEFAULT NULL`, but MySQL's renderer
  always appends `DEFAULT NULL` for nullable columns that lack an explicit
  default AND lack `NOT NULL`.
- omni deparse MUST match: append `DEFAULT NULL` even when `col.Default == nil`
  and the column is nullable. This is the one case where omni **adds** a
  token absent from the catalog.

**Expected SHOW CREATE TABLE:**
```
`id` int DEFAULT NULL
```

**Deparse rule (for /mysql/deparse):**
```
if col.NotNull == false && col.Default == nil && !isAutoIncrement(col) {
    emit(" DEFAULT NULL")
}
```
(Do not emit for TIMESTAMP — that type has its own implicit-default quirk
handled in C16 / `explicit_defaults_for_timestamp`.)

**Oracle verify:** `SHOW CREATE TABLE` output must contain exactly the
`DEFAULT NULL` fragment.

---

### 20.3 String column (VARCHAR/CHAR), NOT NULL, no DEFAULT — implicit ''

**Priority:** MEDIUM  **Status:** pending

```sql
CREATE TABLE c20_3 (name VARCHAR(64) NOT NULL, code CHAR(4) NOT NULL);
INSERT INTO c20_3 () VALUES ();          -- non-strict: name='', code=''
SHOW CREATE TABLE c20_3;
```

**Expected catalog state:**
- Both columns: `Default:nil, NotNull:true`.
- No empty-string literal is stored in the catalog.

**Expected SHOW CREATE TABLE:**
```
`name` varchar(64) NOT NULL,
`code` char(4) NOT NULL
```
(No `DEFAULT ''` rendered.)

**Runtime implicit default source:** `Field_varstring::reset()` /
`Field_string::reset()` in `sql/field.cc` zero the value buffer; for char
columns, `reset()` pads with the fill char (' ').

**Oracle verify:** SHOW CREATE comparison + non-strict INSERT then
`SELECT HEX(name), HEX(code)` to confirm empty-string materialization.

---

### 20.4 ENUM NOT NULL, no DEFAULT — implicit default = first enum value

**Priority:** HIGH  **Status:** pending

```sql
CREATE TABLE c20_4 (
  status ENUM('active','archived','deleted') NOT NULL,
  kind   ENUM('a','b','c')
);
INSERT INTO c20_4 () VALUES ();          -- non-strict: status='active', kind=NULL
SHOW CREATE TABLE c20_4;
```

**Expected catalog state:**
- `status`: `Default:nil, NotNull:true`.
- `kind` (nullable): `Default:nil, NotNull:false` — implicit default is
  NULL, **not** 'a'. The "first value" rule only applies when the column is
  NOT NULL.
- The first-value behavior lives in `Field_enum::reset()` — it stores the
  ordinal `1`, which decodes to the first literal. It is **not** a catalog
  property.

**Expected SHOW CREATE TABLE:**
```
`status` enum('active','archived','deleted') NOT NULL,
`kind` enum('a','b','c') DEFAULT NULL
```
(No `DEFAULT 'active'` rendered for `status`. Do render `DEFAULT NULL` for
the nullable `kind` per rule 20.2.)

**Oracle verify:** Non-strict INSERT + `SELECT status+0, kind FROM c20_4`
confirms ordinal 1 ('active') and NULL.

---

### 20.5 DATETIME/DATE NOT NULL, no DEFAULT — zero-date vs strict error

**Priority:** HIGH  **Status:** pending

```sql
CREATE TABLE c20_5 (
  created_at DATETIME NOT NULL,
  birthday   DATE     NOT NULL
);
-- non-strict, NO_ZERO_DATE off: '0000-00-00 00:00:00' / '0000-00-00'
-- non-strict, NO_ZERO_DATE on:   warning + zero value
-- strict + NO_ZERO_DATE on:      ER_NO_DEFAULT_FOR_FIELD (1364)
INSERT INTO c20_5 () VALUES ();
SHOW CREATE TABLE c20_5;
```

**Expected catalog state:**
- Both columns: `Default:nil, NotNull:true`.
- omni MUST NOT invent a `DEFAULT '0000-00-00'` — even MySQL's own deparser
  does not. The zero value is an execution fallback only.

**Expected SHOW CREATE TABLE:**
```
`created_at` datetime NOT NULL,
`birthday` date NOT NULL
```

**sql_mode interaction:**
- Default 8.0 `sql_mode` includes `STRICT_TRANS_TABLES,NO_ZERO_DATE` so the
  naive `INSERT () VALUES ()` errors with 1364 for InnoDB tables.
- With `sql_mode=''`, the insert succeeds with `'0000-00-00'`.
- This scenario verifies omni **does not pre-apply** the zero-date — leave
  fallback to the server at INSERT time.

**Oracle verify:** Both `sql_mode='' ` and default `sql_mode` INSERT paths,
then SHOW CREATE comparison. The interesting assert is that rendering is
mode-independent.

---

### 20.6 BLOB/TEXT literal DEFAULT — parser error pre-8.0.13, still illegal

**Priority:** HIGH  **Status:** pending

```sql
-- All of these must be rejected:
CREATE TABLE c20_6a (b BLOB DEFAULT 'abc');        -- ER_BLOB_CANT_HAVE_DEFAULT (1101)
CREATE TABLE c20_6b (t TEXT DEFAULT 'hello');       -- 1101
CREATE TABLE c20_6c (g GEOMETRY DEFAULT 'x');       -- 1101
CREATE TABLE c20_6d (j JSON DEFAULT '[]');          -- 1101
```

**Expected omni behavior:**
- Parser accepts the grammar (MySQL accepts it at parse time too).
- `tablecmds.go:CREATE TABLE` application validates: for column types
  `BLOB/TEXT/LONGBLOB/LONGTEXT/MEDIUMBLOB/MEDIUMTEXT/TINYBLOB/TINYTEXT/
  GEOMETRY/POINT/LINESTRING/POLYGON/MULTI*/GEOMETRYCOLLECTION/JSON`, if
  `DefaultValue` is a plain literal (non-parenthesized `Expr`), return
  `ER_BLOB_CANT_HAVE_DEFAULT` (1101).
- The error must mention the offending column name, matching MySQL's
  `"BLOB, TEXT, GEOMETRY or JSON column '%s' can't have a default value"`.

**Expected catalog state:** table creation fails; no row added to catalog.

**Source reference:** `sql/sql_table.cc` `prepare_create_field()` — search
for `ER_BLOB_CANT_HAVE_DEFAULT`. The check is gated on
`sql_field->constant_default != nullptr` (literal) and the type being one of
the blob-family types. Expression defaults (`m_default_val_expr != nullptr`)
are allowed and take a different code path.

**Oracle verify:** Run all 4 statements on MySQL 8.0; assert error code
`1101` and compare message text (substring match on column name).

---

### 20.7 JSON / BLOB expression DEFAULT — parenthesized expression OK (8.0.13+)

**Priority:** HIGH  **Status:** pending

```sql
CREATE TABLE c20_7 (
  id INT PRIMARY KEY,
  tags  JSON      DEFAULT (JSON_ARRAY()),
  meta  JSON      DEFAULT (JSON_OBJECT('v', 1)),
  blob1 BLOB      DEFAULT (SUBSTRING('abcdef', 1, 3)),
  pt    POINT     DEFAULT (POINT(0, 0)),
  uuid  BINARY(16) DEFAULT (UUID_TO_BIN(UUID()))
);
SHOW CREATE TABLE c20_7;
INSERT INTO c20_7 (id) VALUES (1);
SELECT tags, meta FROM c20_7;
```

**Parser requirement:**
- The `(` after `DEFAULT` must open a grammar path that parses a full `Expr`
  (not a `Literal`). When the parser sees `DEFAULT` followed by `(` at
  column-option position, it must consume an expression in parens.
- An unparenthesized call `DEFAULT JSON_ARRAY()` is a **parser error** in
  MySQL 8.0 — omni must reject it with a syntax error or `ER_INVALID_DEFAULT`
  (1067). Only the parenthesized form is legal for non-constant defaults on
  these types.

**Expected catalog state:**
- `tags.DefaultAnalyzed` holds an analyzed `FuncCall{JSON_ARRAY}` node.
- `tags.Default` (raw) must preserve that the user wrote a parenthesized
  expression. Suggested representation: store the inner expression plus a
  flag `IsExprDefault bool`, or always wrap in parens on deparse when
  `DefaultAnalyzed.Kind != Literal`.
- `pt`, `uuid`, `blob1` same pattern.

**Expected SHOW CREATE TABLE rendering:**
```
`tags` json DEFAULT (json_array()),
`meta` json DEFAULT (json_object(_utf8mb4'v',1)),
...
```
(MySQL renders the inner expression back lowercased and fully qualified with
charset introducers for string literals. Match at the **AST-equivalent**
level in oracle tests — do not string-compare.)

**Forward-reference constraint:** Expression defaults cannot reference later
columns with expression defaults or generated columns. Tablecmds application
must walk `base_columns_map` (mysql-server concept) or its omni equivalent
and error with `ER_DEFAULT_VAL_GENERATED_REF_AUTO_INC` / `ER_DEFAULT_VAL_
GENERATED_...` family codes if violated.

**Oracle verify:** SHOW CREATE round-trip + `INSERT (id) VALUES (1)` then
`SELECT JSON_LENGTH(tags), JSON_EXTRACT(meta,'$.v')` to confirm the
expression actually materialized its value.

---

### 20.8 Generated column — DEFAULT clause is a grammar error

**Priority:** MEDIUM  **Status:** pending

```sql
-- ALL invalid:
CREATE TABLE c20_8a (
  a INT,
  b INT AS (a + 1) DEFAULT 0                 -- stored or virtual: illegal
);

CREATE TABLE c20_8b (
  a INT,
  b INT GENERATED ALWAYS AS (a + 1) VIRTUAL DEFAULT 0
);

CREATE TABLE c20_8c (
  a INT,
  b INT GENERATED ALWAYS AS (a + 1) STORED DEFAULT (a * 2)
);
```

**Expected omni behavior:**
- Parser: reject at grammar level. In MySQL's Bison grammar, column options
  after `GENERATED ALWAYS AS (...) [STORED|VIRTUAL]` are a restricted subset
  — `DEFAULT` is not in that subset. omni's parser should either produce
  `ER_PARSE_ERROR` or `ER_WRONG_USAGE` mentioning that generated columns
  cannot have a DEFAULT clause.
- If the parser accepts for error-recovery reasons, `tablecmds.go` must
  reject: if `col.Generated != nil && col.Default != nil`, error with
  `"A generated column cannot have a default value"`.

**Error code:** MySQL emits `ER_WRONG_USAGE (1221)` with message
`"Incorrect usage of DEFAULT and generated column"` (or parse error 1064
depending on position — test both orderings).

**Expected catalog state:** no table created.

**Note on `NOT NULL`:** `NOT NULL` IS allowed on generated columns (it's a
constraint, not a default), so the test must isolate `DEFAULT` as the
offending clause.

**Oracle verify:** All three statements against MySQL 8.0; assert error code
(1064 or 1221) and capture the exact message for omni parity.

---

## Section C21: Parser-level implicit defaults

> **Expansion note (Wave 3):** grew from 1 to 10 scenarios via systematic
> walk of `/Users/rebeliceyang/Github/mysql-server/sql/sql_yacc.yy` grammar
> action blocks. These are "invisible" grammar-level defaults — the rules
> fire action code on the empty production and don't appear in reference
> manual syntax diagrams. Category anchors: `sql_yacc.yy`, `sql_lex.cc`,
> `sql_cmd_ddl_table.cc`, `sql_view.cc`.

### 21.1 `DEFAULT` without value on nullable column → DEFAULT NULL

**Priority:** HIGH  **Status:** verified

**Setup:**
```sql
CREATE TABLE t (a INT DEFAULT NULL);
-- and the inverse (column with no DEFAULT) implicitly yields DEFAULT NULL
CREATE TABLE t2 (c INT);
```

**Grammar rule:** `sql_yacc.yy:7651 opt_default`
```
opt_default:
          %empty {}
        | DEFAULT_SYM {}
        ;
```

`column_attribute` (sql_yacc.yy:7467) has no clause that attaches an implicit
`DEFAULT NULL` marker — when the column definition omits both `NOT NULL` and
`DEFAULT <v>`, the column is nullable and the server later materializes
`DEFAULT NULL` for `INFORMATION_SCHEMA.COLUMNS.COLUMN_DEFAULT`.

**Oracle verification:**
```sql
SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a';
```
Expected: `COLUMN_DEFAULT = NULL`, `IS_NULLABLE = YES`.

**omni assertion:** `tbl.GetColumn("a").Default` is the NULL literal; nullable
true. `mysql/parser` should treat missing DEFAULT on a nullable column as
equivalent to `DEFAULT NULL` when computing view/catalog default expressions.

**See catalog:** C21.1.

---

### 21.2 Bare JOIN keyword → INNER JOIN (JTT_INNER)

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
SELECT * FROM t1 JOIN t2 ON t1.a = t2.a;       -- same as INNER JOIN
SELECT * FROM t1 INNER JOIN t2 ON t1.a = t2.a;
SELECT * FROM t1 CROSS JOIN t2;                -- also yields JTT_INNER
```

**Grammar rule:** `sql_yacc.yy:11975 inner_join_type`
```
inner_join_type:
          JOIN_SYM                         { $$= JTT_INNER; }
        | INNER_SYM JOIN_SYM               { $$= JTT_INNER; }
        | CROSS JOIN_SYM                   { $$= JTT_INNER; }
        | STRAIGHT_JOIN                    { $$= JTT_STRAIGHT_INNER; }
```

Three textually-different forms (`JOIN`, `INNER JOIN`, `CROSS JOIN`) all
collapse to the same enum value `JTT_INNER`. Additionally `opt_inner:`
(sql_yacc.yy:11986) and `opt_outer:` (sql_yacc.yy:11991) are empty-acceptance
rules that make the `INNER` / `OUTER` keywords purely syntactic noise.

**omni check:** `mysql/parser/select.go` — the AST node for a JOIN should use
a single canonical `INNER` value for all three forms, or record them
distinctly only for deparse round-tripping. The semantic IR must normalize.

---

### 21.3 ORDER BY column without direction → ORDER_NOT_RELEVANT (NOT ORDER_ASC)

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
SELECT * FROM t ORDER BY a;        -- emits ORDER_NOT_RELEVANT
SELECT * FROM t ORDER BY a ASC;    -- emits ORDER_ASC
```

**Grammar rule:** `sql_yacc.yy:12606 opt_ordering_direction`
```
opt_ordering_direction:
          %empty { $$= ORDER_NOT_RELEVANT; }
        | ordering_direction
        ;

ordering_direction:
          ASC         { $$= ORDER_ASC; }
        | DESC        { $$= ORDER_DESC; }
        ;
```

**Important correction to starmap description:** Contrary to the plan's claim
that "ORDER BY without direction → ASC", the **grammar emits a distinct third
value `ORDER_NOT_RELEVANT`**. The server later treats `ORDER_NOT_RELEVANT` the
same as `ORDER_ASC` in execution, but the AST preserves the distinction. This
matters for deparse (`SHOW CREATE VIEW` should reproduce the user's original
form without injecting `ASC`).

**omni check:** AST ordering-direction field should carry a tri-state, not
bool. Currently omni may collapse to ASC — verify `mysql/parser/select.go`.

---

### 21.4 LIMIT N without OFFSET → opt_offset = NULL (NOT OFFSET 0)

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
SELECT * FROM t LIMIT 10;            -- opt_offset = NULL
SELECT * FROM t LIMIT 10 OFFSET 0;   -- opt_offset = Item_uint(0)
SELECT * FROM t LIMIT 0, 10;         -- opt_offset = Item_uint(0), is_offset_first=true
```

**Grammar rule:** `sql_yacc.yy:12628 limit_options`
```
limit_options:
          limit_option
          {
            $$.limit= $1;
            $$.opt_offset= NULL;
            $$.is_offset_first= false;
          }
        | limit_option ',' limit_option
          {
            $$.limit= $3;
            $$.opt_offset= $1;
            $$.is_offset_first= true;
          }
        | limit_option OFFSET_SYM limit_option
          { ... }
        ;
```

**Correction:** The plan says `LIMIT N` → `OFFSET 0`. The grammar emits
`opt_offset = NULL`. Execution semantics treat NULL as 0, but deparse and
SHOW CREATE VIEW must preserve the absence. Additionally, `LIMIT a, b` and
`LIMIT b OFFSET a` are grammatically distinct (tracked by the
`is_offset_first` bool).

**omni check:** `mysql/parser` LIMIT clause should record offset as optional
pointer (nullable), not default to 0. Also must track `is_offset_first`.

---

### 21.5 FK ON DELETE omitted → FK_OPTION_UNDEF (NOT FK_OPTION_RESTRICT)

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
CREATE TABLE child (
  p INT,
  FOREIGN KEY (p) REFERENCES parent(id)    -- no ON DELETE, no ON UPDATE
);
```

**Grammar rule:** `sql_yacc.yy:7814 opt_on_update_delete`
```
opt_on_update_delete:
          %empty
          {
            $$.fk_update_opt= FK_OPTION_UNDEF;
            $$.fk_delete_opt= FK_OPTION_UNDEF;
          }
        | ON_SYM UPDATE_SYM delete_option
          {
            $$.fk_update_opt= $3;
            $$.fk_delete_opt= FK_OPTION_UNDEF;
          }
        ...
```

**Important correction to starmap description:** The plan says the parser
emits `FK_OPTION_RESTRICT` for omitted clauses. This is **wrong** — the
grammar emits `FK_OPTION_UNDEF`. The mapping to RESTRICT / NO ACTION happens
later in InnoDB (`dd::Foreign_key::delete_rule` defaults). The `FK_OPTION_*`
enum has all five explicit values (sql_yacc.yy:7844) plus `FK_OPTION_UNDEF`
which is grammar-only. Also note `RESTRICT` and `NO ACTION` parse to
**distinct enum values**, even though InnoDB treats them identically at
runtime — another implicit normalization.

**omni check:** FK parser must use a 6-value enum (including UNDEF). Treating
omitted as RESTRICT will mis-deparse.

---

### 21.6 FK ON DELETE present but ON UPDATE omitted → UPDATE gets FK_OPTION_UNDEF

**Priority:** MEDIUM  **Status:** pending

**Trigger SQL:**
```sql
CREATE TABLE child (
  p INT,
  FOREIGN KEY (p) REFERENCES parent(id) ON DELETE CASCADE
);
-- fk_update_opt = FK_OPTION_UNDEF (NOT inherited from delete rule)
```

**Grammar rule:** Same `sql_yacc.yy:7814 opt_on_update_delete`, specifically
the branch:
```
        | ON_SYM DELETE_SYM delete_option
          {
            $$.fk_update_opt= FK_OPTION_UNDEF;
            $$.fk_delete_opt= $3;
          }
```

The asymmetric branches (four total) mean specifying one action never
implicitly fills in the other — each clause defaults independently to
`FK_OPTION_UNDEF`.

**omni check:** Parser must not copy the specified action into the
unspecified slot.

---

### 21.7 CREATE INDEX without USING clause → nullptr (engine picks default)

**Priority:** MEDIUM  **Status:** pending

**Trigger SQL:**
```sql
CREATE TABLE t (a INT, KEY k (a));                 -- no USING, no TYPE
CREATE TABLE t (a INT, KEY k (a) USING BTREE);     -- explicit
```

**Grammar rule:** `sql_yacc.yy:8006 opt_index_type_clause`
```
opt_index_type_clause:
          %empty { $$ = nullptr; }
        | index_type_clause
        ;

index_type_clause:
          USING index_type    { $$= NEW_PTN PT_index_type($2); }
        | TYPE_SYM index_type { $$= NEW_PTN PT_index_type($2); }
        ;

index_type:                                        (sql_yacc.yy:8021)
          BTREE_SYM { $$= HA_KEY_ALG_BTREE; }
        | RTREE_SYM { $$= HA_KEY_ALG_RTREE; }
        | HASH_SYM  { $$= HA_KEY_ALG_HASH; }
        ;
```

The default `nullptr` bubbles up to handler; InnoDB converts
`HA_KEY_ALG_SE_SPECIFIC` (from `KEY_ALGORITHM_NONE`) to BTREE; MEMORY converts
to HASH. So "default" is engine-dependent, not constant BTREE.

**omni check:** Parser index node needs a nullable/tri-state field; the
catalog layer must resolve the engine default, not the parser.

---

### 21.8 INSERT without column list → empty PT_item_list (resolved later)

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
INSERT INTO t VALUES (1, 2, 3);            -- no column list
INSERT INTO t () VALUES ();                -- empty column list (MySQL ext)
INSERT INTO t (a, b, c) VALUES (1, 2, 3);  -- explicit
```

**Grammar rule:** `sql_yacc.yy:13241 insert_from_constructor`
```
insert_from_constructor:
          insert_values
          {
            $$.column_list= NEW_PTN PT_item_list;   // empty list
            $$.row_value_list= $1;
          }
        | '(' ')' insert_values
          {
            $$.column_list= NEW_PTN PT_item_list;   // also empty list
            $$.row_value_list= $3;
          }
        | '(' insert_columns ')' insert_values
          { ... }
        ;
```

Notice the first two branches are **grammatically distinct but semantically
identical** — both produce an empty `PT_item_list`. The parser has no way to
distinguish `INSERT INTO t VALUES(...)` from `INSERT INTO t () VALUES(...)`
downstream. The column list is resolved to all columns in their
table-declared order at name-resolution time. Same pattern mirrored for
`insert_query_expression` at sql_yacc.yy:13259.

**omni check:** Query-span analysis of INSERT must expand empty column list
to the full column order from the catalog. Mis-ordering breaks lineage.

---

### 21.9 CREATE VIEW without ALGORITHM → VIEW_ALGORITHM_UNDEFINED

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
CREATE VIEW v AS SELECT 1;                        -- default
CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT 1;    -- explicit, same result
CREATE ALGORITHM=MERGE VIEW v AS SELECT 1;
```

**Grammar rule + LEX default:** `sql_lex.cc:421`
```
create_view_algorithm = VIEW_ALGORITHM_UNDEFINED;
```

In `sql_yacc.yy:17439 view_algorithm` only explicit forms set a value; the
`view_replace_or_algorithm` rule (sql_yacc.yy:17425) is itself optional
because the CREATE VIEW grammar can start without it. The LEX struct is
zero-initialized and `lex_start()` sets `VIEW_ALGORITHM_UNDEFINED` as the
default before parsing, meaning a missing ALGORITHM clause is
indistinguishable from an explicit `ALGORITHM=UNDEFINED`.

Also note `ALTER VIEW` without an explicit algorithm re-asserts the default
at sql_yacc.yy:8226:
```
            lex->create_view_algorithm= VIEW_ALGORITHM_UNDEFINED;
```

Further post-processing in sql_view.cc:942-945 can silently downgrade
`VIEW_ALGORITHM_MERGE` to `VIEW_ALGORITHM_UNDEFINED` when the view body is
not mergeable — another implicit behavior the parser cannot see.

**omni check:** View parser should record algorithm as nullable or tri-state;
deparse should omit the ALGORITHM clause when it's the default.

---

### 21.10 CREATE TABLE without ENGINE → post-parse fill from @@default_storage_engine

**Priority:** HIGH  **Status:** pending

**Trigger SQL:**
```sql
CREATE TABLE t (a INT);                  -- no ENGINE=...
CREATE TABLE t (a INT) ENGINE=InnoDB;    -- explicit
```

**Grammar:** `opt_create_table_options_etc` (sql_yacc.yy:6237) accepts an
empty body. The ENGINE clause appears only under `create_table_option` at
sql_yacc.yy:6705:
```
          ENGINE_SYM opt_equal ident_or_text
```
No empty alternative sets a default at parse time.

**Post-parse fill:** `sql_cmd_ddl_table.cc:170-178`
```cpp
/*
  If no engine type was given, work out the default now
  rather than at parse-time.
*/
if (!(create_info.used_fields & HA_CREATE_USED_ENGINE))
  create_info.db_type = create_info.options & HA_LEX_CREATE_TMP_TABLE
                            ? ha_default_temp_handlerton(thd)
                            : ha_default_handlerton(thd);
```

Key implicit behaviors:
1. `HA_CREATE_USED_ENGINE` bit tells whether ENGINE was explicit.
2. The default is **not constant InnoDB** — it's `@@default_storage_engine`
   (session variable).
3. Temporary tables use a separate `@@default_tmp_storage_engine`.
4. Fill happens in the *command executor*, not the parser — so the AST
   produced by yacc has a NULL engine, and only the post-parse layer
   resolves it.

**omni check:** Parser must leave engine NULL; catalog/advisor logic must
respect that "no engine" ≠ "ENGINE=InnoDB". Tests for engine defaulting
must set the session variable.

---

## Section C22: ALTER TABLE algorithm / lock defaults

> **New section (Wave 1, 2026-04-13).** `ALTER TABLE`'s `ALGORITHM=` and
> `LOCK=` clauses default to a value that MySQL picks per operation, and an
> explicit request that the chosen operation can't satisfy is a hard error —
> MySQL does not silently downgrade. omni's catalog does not track algorithm
> or lock (they are execution-time concerns, not persisted schema state), but
> the SDL diff engine must know which ALTER operations are INSTANT / INPLACE /
> COPY so it can decide whether to emit explicit clauses, split a multi-clause
> ALTER into several statements, or warn the user that a generated migration
> will rebuild the table under an EXCLUSIVE lock.
>
> **MySQL source anchors (all relative to `/Users/rebeliceyang/Github/mysql-server`):**
> - `sql/sql_alter.h:354-365` — `enum enum_alter_table_algorithm { ALTER_TABLE_ALGORITHM_DEFAULT, INPLACE, INSTANT, COPY }`.
> - `sql/sql_alter.h:374-383` — `enum enum_alter_table_lock { DEFAULT, NONE, SHARED, EXCLUSIVE }`.
> - `sql/sql_alter.h:467-468` — constructor defaults both fields to `..._DEFAULT`.
> - `sql/sql_alter.cc:78-79` — copy-constructor propagates `requested_algorithm` / `requested_lock`.
> - `sql/sql_table.cc` — `mysql_alter_table()` / `ha_inplace_alter_table()` do the fallback chain (INSTANT → INPLACE → COPY) and error if the explicit request can't be honored.
>
> **Documentation anchor:** https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
> ("If the ALGORITHM clause is omitted, MySQL uses `ALGORITHM=INSTANT` for
> storage engines and `ALTER TABLE` clauses that support it. Otherwise,
> `ALGORITHM=INPLACE` is used. If `ALGORITHM=INPLACE` is not supported,
> `ALGORITHM=COPY` is used.")
>
> **omni parser state (for reference):**
> - `mysql/parser/alter_table.go:36` — grammar comment lists
>   `ALGORITHM [=] {DEFAULT | INSTANT | INPLACE | COPY}`.
> - `mysql/parser/alter_table.go:127-138` — parser recognizes `ALGORITHM` and
>   emits `cmd.Type = nodes.ATAlgorithm` (there is a parallel `ATLock`).
> - `mysql/ast/parsenodes.go:184` — `ATAlgorithm` enum member (and `ATLock`).
> - `mysql/catalog/altercmds.go` — **no** handling for `ATAlgorithm` / `ATLock`.
>   The catalog layer intentionally ignores these clauses (no persisted state
>   to mutate). That is the correct behavior: algorithm/lock are runtime-only.

---

### 22.1 `ALGORITHM=DEFAULT` picks fastest supported (INSTANT > INPLACE > COPY)

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB;
ALTER TABLE t1 ADD COLUMN b INT;   -- no ALGORITHM= clause
```
**Expected (MySQL 8.0.12+):**
- `requested_algorithm = ALTER_TABLE_ALGORITHM_DEFAULT` (the parser default per `sql/sql_alter.h:467`).
- Server walks INSTANT → INPLACE → COPY in `mysql_alter_table()` and picks `INSTANT` for this `ADD COLUMN` case.
- No error, no warning, and `SHOW CREATE TABLE` is unchanged apart from the new column.

**Oracle verification:** execute against MySQL 8.0 container; confirm via
`information_schema.innodb_metrics` (`ddl_instant_alter_table`) that the
INSTANT counter advanced. (Not performed in this scenario file — recorded so
Wave 5 can automate.)

**omni assertion:**
- Parser accepts the bare `ALTER TABLE ... ADD COLUMN` with no ALGORITHM clause.
- `mysql/catalog/altercmds.go` applies the `ADD COLUMN` to `tbl.Columns`; no
  algorithm bookkeeping is expected. SDL diff callers that care about the
  chosen algorithm must consult a helper (to be written) keyed on the
  `AlterCmd.Type` — **this scenario is the justification for that helper.**

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_alter.h:354-365,467`; docs alter-table.html "ALGORITHM" section.

---

### 22.2 `LOCK=DEFAULT` picks least restrictive supported lock

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB;
ALTER TABLE t1 ADD INDEX ix_a (a);   -- no LOCK= clause
```
**Expected:**
- `requested_lock = ALTER_TABLE_LOCK_DEFAULT` (`sql/sql_alter.h:374-383`).
- InnoDB reports `HA_ALTER_INPLACE_NO_LOCK` for secondary-index add → effective `LOCK=NONE`, concurrent reads and writes permitted.
- `CHANGE COLUMN type` (scenario 22.6) would instead map DEFAULT to `EXCLUSIVE`.

**omni assertion:**
- Parser accepts without LOCK clause. Catalog applies the index add.
- SDL diff generator's lock-classification helper must return `NONE` for
  `ATAddConstraint{Index}` over InnoDB, so generated migration scripts can
  annotate themselves as online-safe.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_alter.h:374-383`; docs innodb-online-ddl-operations.html (ADD INDEX row: LOCK=NONE).

---

### 22.3 `ADD COLUMN` (trailing, nullable) is INSTANT in 8.0.12+

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB ROW_FORMAT=DYNAMIC;
ALTER TABLE t1 ADD COLUMN b VARCHAR(32) NULL;                    -- 22.3a
ALTER TABLE t1 ADD COLUMN c INT NULL, ALGORITHM=INSTANT;         -- 22.3b
```
**Expected:**
- Both succeed with the INSTANT algorithm.
- **Restrictions (from docs innodb-online-ddl-operations.html):**
  - Not allowed for `ROW_FORMAT=COMPRESSED` tables.
  - Not allowed if the table has `FULLTEXT` indexes.
  - Max 64 row versions accumulated before a real rebuild is required.
  - After an INSTANT add, `ALTER TABLE ... EXCHANGE PARTITION` is blocked on
    the table forever (alter-table.html "Key Limitation").
- `ADD COLUMN ... FIRST` or `AFTER` still qualifies for INSTANT in 8.0.29+ but
  not in 8.0.12 — the SDL diff generator should treat positional adds as
  INPLACE to stay portable across the 8.0 line.

**omni assertion:**
- Parser accepts both. Catalog appends the column in both cases.
- SDL diff helper classifies `ATAddColumn` with nullable, default-less,
  non-positional, non-compressed, non-FULLTEXT table as INSTANT-eligible.

**Priority:** HIGH
**Status:** pending
**Source anchors:** docs innodb-online-ddl-operations.html; `sql/sql_table.cc` INSTANT eligibility checks.

---

### 22.4 `DROP COLUMN` is INSTANT-capable in 8.0.29+ but historically INPLACE rebuild

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT, b INT) ENGINE=InnoDB;
ALTER TABLE t1 DROP COLUMN b;                        -- DEFAULT → INSTANT (8.0.29+) else INPLACE
ALTER TABLE t1 DROP COLUMN b, ALGORITHM=INSTANT;     -- hard error pre-8.0.29
```
**Expected:**
- On MySQL 8.0.29+, `DROP COLUMN` is eligible for INSTANT.
- On earlier 8.0.x, INSTANT is rejected with
  `ER_ALTER_OPERATION_NOT_SUPPORTED` ("ALGORITHM=INSTANT is not supported for
  this operation. Try ALGORITHM=COPY/INPLACE.") — MySQL errors rather than
  downgrades when the explicit clause can't be honored (see 22.7).
- DEFAULT on 8.0.12–8.0.28 picks INPLACE (rebuild); 8.0.29+ picks INSTANT.

**omni assertion:**
- Parser accepts both forms. Catalog removes the column from
  `tbl.Columns` regardless of algorithm.
- The SDL diff generator should, by default, emit `DROP COLUMN` **without** an
  explicit `ALGORITHM=` clause so the server picks the best option available
  at the target version. Only emit `ALGORITHM=INSTANT` when the user has
  opted in to "instant-only migrations".

**Priority:** MED
**Status:** pending
**Source anchors:** docs innodb-online-ddl-operations.html (DROP COLUMN row); `sql/sql_table.cc` INSTANT eligibility; error class `ER_ALTER_OPERATION_NOT_SUPPORTED_REASON`.

---

### 22.5 `RENAME COLUMN` (metadata-only) is INPLACE, upgraded to INSTANT in 8.0.29+

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, old_name INT) ENGINE=InnoDB;
ALTER TABLE t1 RENAME COLUMN old_name TO new_name;                -- 22.5a, no clause
ALTER TABLE t1 CHANGE COLUMN old_name new_name INT;               -- 22.5b, CHANGE form
```
**Expected:**
- Rename without type change → INPLACE (or INSTANT on 8.0.29+), `LOCK=NONE`.
- **Restriction:** column referenced by a FOREIGN KEY cannot rename under
  INSTANT; the server drops back to INPLACE.
- CHANGE COLUMN form (22.5b) with an unchanged type is treated the same; if
  the type differs even trivially (e.g. attribute change), the server may
  fall back to COPY — see 22.6.

**omni assertion:**
- Parser already carries rename info through `ATCmd`. The catalog layer
  rewrites `Column.Name`.
- SDL diff classifies a pure rename (no type/default change, not an FK column)
  as INPLACE-`NONE`. If the column is on either side of a `ForeignKey`
  in the catalog, the classifier must downgrade to INPLACE with `LOCK=SHARED`
  on MySQL < 8.0.29.

**Priority:** MED
**Status:** pending
**Source anchors:** docs innodb-online-ddl-operations.html (RENAME COLUMN row); alter-table.html "INSTANT Algorithm Operations".

---

### 22.6 `CHANGE COLUMN` type change (e.g. INT → BIGINT) forces COPY

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB;
ALTER TABLE t1 CHANGE COLUMN a a BIGINT;                         -- 22.6a
ALTER TABLE t1 MODIFY COLUMN a BIGINT, ALGORITHM=INPLACE;        -- 22.6b → hard error
ALTER TABLE t1 MODIFY COLUMN a BIGINT, LOCK=NONE;                -- 22.6c → hard error
```
**Expected:**
- 22.6a: DEFAULT resolves to `COPY`, `LOCK=DEFAULT` resolves to `SHARED`/
  `EXCLUSIVE` (COPY always blocks concurrent DML).
- 22.6b: `ER_ALTER_OPERATION_NOT_SUPPORTED` — the server will not silently
  downgrade an explicit `INPLACE` request to COPY.
- 22.6c: `ER_ALTER_OPERATION_NOT_SUPPORTED` — `LOCK=NONE` is incompatible
  with the COPY algorithm that this operation requires.
- Type widening (INT → BIGINT), narrowing (BIGINT → INT), sign flip
  (INT → INT UNSIGNED), CHARSET change, and `ENUM` reordering all force COPY.

**omni assertion:**
- Parser accepts all three statements.
- Catalog re-types the column in all three (catalog is shape-only; it has no
  idea the second two will fail at execution). The catalog does **not** need
  to reject these at apply time.
- SDL diff's classifier must flag these as COPY and never emit
  `ALGORITHM=INPLACE` or `LOCK=NONE` for them; otherwise the generated
  migration will hit `ER_ALTER_OPERATION_NOT_SUPPORTED` against the real
  server.

**Priority:** HIGH
**Status:** pending
**Source anchors:** docs innodb-online-ddl-operations.html (CHANGE COLUMN TYPE row); alter-table.html "Behavior When Requested Algorithm Not Supported".

---

### 22.7 Explicit `ALGORITHM=INSTANT` on unsupported operation is an error (no downgrade)

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB;
ALTER TABLE t1 DROP PRIMARY KEY, ADD PRIMARY KEY (a), ALGORITHM=INSTANT;
```
**Expected:**
- `ER_ALTER_OPERATION_NOT_SUPPORTED_REASON` with reason "Cannot change
  column type INPLACE" / "PRIMARY KEY rebuild requires table copy".
- Omitting the clause would have auto-selected COPY and succeeded.
- The server **never** downgrades an explicit `INSTANT`/`INPLACE` request;
  per docs alter-table.html, "Specifying an `ALGORITHM` clause requires the
  operation to use the specified algorithm ... or fail with an error
  otherwise."

**omni assertion:**
- Parser accepts the statement; `ATAlgorithm` carries the `INSTANT` value.
- Catalog applies the PK swap (`mysql/catalog/altercmds.go` is
  algorithm-oblivious, which is correct).
- SDL diff must **not** emit `ALGORITHM=INSTANT` for PK changes. A
  classification helper keyed on `(AlterCmd.Type, op-specific details)`
  should return `AlgorithmCopy` for `DropConstraint(Primary)` followed by
  `AddConstraint(Primary)`.

**Priority:** HIGH
**Status:** pending
**Source anchors:** alter-table.html "Behavior When Requested Algorithm Not Supported"; `sql/sql_table.cc` algorithm-check branches.

---

### 22.8 `LOCK=NONE` on COPY-only operation is a hard error; `LOCK` is ignored with `ALGORITHM=INSTANT`

**Setup:**
```sql
CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB;
-- 22.8a: incompatible lock for COPY
ALTER TABLE t1 DROP PRIMARY KEY, ADD PRIMARY KEY (a), LOCK=NONE;
-- 22.8b: LOCK specified alongside INSTANT → only LOCK=DEFAULT allowed
ALTER TABLE t1 ADD COLUMN b INT, ALGORITHM=INSTANT, LOCK=NONE;
-- 22.8c: DEFAULT lock is always legal
ALTER TABLE t1 ADD COLUMN c INT, ALGORITHM=INSTANT;
```
**Expected:**
- 22.8a: `ER_ALTER_OPERATION_NOT_SUPPORTED_REASON` — PK rebuild is COPY-only
  and COPY cannot honor `LOCK=NONE`.
- 22.8b: `ER_ALTER_OPERATION_NOT_SUPPORTED_REASON` — per alter-table.html,
  "Only `LOCK = DEFAULT` is permitted for operations using
  `ALGORITHM=INSTANT`." Even if the operation would be lock-free, any
  explicit `LOCK=` other than `DEFAULT` is rejected.
- 22.8c: succeeds (INSTANT with implicit DEFAULT lock).

**omni assertion:**
- Parser accepts all three.
- Catalog applies the column add / PK swap obliviously.
- SDL diff classification helper must enforce two rules when choosing
  whether to emit an explicit `LOCK=` clause:
  1. If op is COPY-only → never emit `LOCK=NONE` or `LOCK=SHARED`.
  2. If op is INSTANT → never emit `LOCK=` at all (only `LOCK=DEFAULT` is
     legal with INSTANT, and that's identical to omitting the clause).

**Priority:** HIGH
**Status:** pending
**Source anchors:** alter-table.html "Only LOCK = DEFAULT is permitted for operations using ALGORITHM=INSTANT"; `sql/sql_alter.h:374-383`.

---

## Section C23: NULL in string context

### 23.1 CONCAT with any NULL argument → NULL result

**Setup:**
```sql
CREATE TABLE t (a VARCHAR(10), b VARCHAR(10));
INSERT INTO t VALUES ('foo', NULL), ('foo', 'bar');
```

**Oracle verification:**
```sql
SELECT CONCAT(a, b) FROM t;
SELECT CONCAT('x', NULL, 'y');
SELECT CONCAT('x', '', 'y');             -- empty string is NOT NULL → 'xy'
```
Expected: first query → `NULL`, `'foobar'`; literal `CONCAT('x', NULL, 'y')` → `NULL`; empty string → `'xy'`.

**omni assertion:** query-span / nullability analyzer infers `CONCAT(...) IS NULL` when ANY input may be NULL. An empty string argument must NOT be treated as NULL. The resulting column's nullability = OR of all input nullabilities.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/item_strfunc.cc` `Item_func_concat::val_str()` — loops arguments, returns `null_value = true` on first NULL arg; string-functions.html "CONCAT() returns NULL if any argument is NULL."

---

### 23.2 CONCAT_WS skips NULL arguments (separator non-null)

**Setup:**
```sql
CREATE TABLE t (a VARCHAR(10), b VARCHAR(10), c VARCHAR(10));
INSERT INTO t VALUES ('x', NULL, 'z'), (NULL, NULL, NULL);
```

**Oracle verification:**
```sql
SELECT CONCAT_WS(',', a, b, c) FROM t;
SELECT CONCAT_WS(',', 'x', NULL, 'z');        -- 'x,z' (NULL skipped, NO trailing/double sep)
SELECT CONCAT_WS(',', NULL, NULL, NULL);      -- '' (empty, NOT NULL)
SELECT CONCAT_WS(NULL, 'x', 'y');             -- NULL (separator NULL → NULL result)
```
Expected: row1 → `'x,z'`, row2 → `''`; literal cases per comments.

**omni assertion:** nullability analyzer must special-case CONCAT_WS: result is NULL iff the separator (arg 0) is NULL. Data args being NULL does NOT make the result NULL. When all data args are NULL, result is empty string.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/item_strfunc.cc` `Item_func_concat_ws::val_str()` — `if (i == 0) { /* sep NULL → null_value */ } else continue;`; string-functions.html "CONCAT_WS() does not skip empty strings. However, it does skip any NULL values after the separator argument." + "If the separator is NULL, the result is NULL."

---

### 23.3 IFNULL / COALESCE as rescue for CONCAT NULL-propagation

**Setup:**
```sql
CREATE TABLE t (first_name VARCHAR(20), middle_name VARCHAR(20), last_name VARCHAR(20));
INSERT INTO t VALUES ('Ada', NULL, 'Lovelace');
```

**Oracle verification:**
```sql
SELECT CONCAT(first_name, ' ', middle_name, ' ', last_name) FROM t;                 -- NULL
SELECT CONCAT(first_name, ' ', IFNULL(middle_name, ''), ' ', last_name) FROM t;     -- 'Ada  Lovelace'
SELECT CONCAT(first_name, ' ', COALESCE(middle_name, ''), ' ', last_name) FROM t;   -- 'Ada  Lovelace'
SELECT CONCAT_WS(' ', first_name, middle_name, last_name) FROM t;                   -- 'Ada Lovelace'
```

**omni assertion:** nullability analyzer treats `IFNULL(x, literal_not_null)` and `COALESCE(x, literal_not_null, ...)` where at least one non-null literal is present as non-nullable, allowing the enclosing CONCAT to be reported as non-nullable too. CONCAT_WS provides the idiomatic alternative.

**Priority:** LOW
**Status:** pending
**Source anchors:** `sql/item_cmpfunc.cc` `Item_func_ifnull::fix_length_and_dec()` sets `maybe_null` based on second arg; `Item_func_coalesce` likewise propagates the intersection of input null-ness.

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

**Priority:** MED
**Status:** verified
**Source anchors:** C13.1 + C24.1. Verified by `TestScenario_C24`.

---

### 24.2 GIPK column spec: name, type, attributes

**Setup:**
```sql
SET SESSION sql_generate_invisible_primary_key = ON;
CREATE TABLE t (a INT);
SET SESSION show_gipk_in_create_table_and_information_schema = ON;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, EXTRA, COLUMN_KEY
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
SHOW INDEX FROM t;
```
Expected: GIPK column is named exactly `my_row_id`, type `bigint unsigned`, `NOT NULL`, `EXTRA='auto_increment INVISIBLE'`, `COLUMN_KEY='PRI'`. It is the first column of the table. A PRIMARY KEY index on `my_row_id` exists.

**omni assertion:** when the catalog materializes the GIPK, the generated column must be: name `my_row_id`, type `BIGINT UNSIGNED`, NOT NULL, AUTO_INCREMENT, INVISIBLE, and inserted at position 0. A PRIMARY KEY constraint over `(my_row_id)` is added.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `validate_and_generate_invisible_primary_key()` (builds `Create_field` with `PRIMARY_KEY_NAME_STRING`, `AUTO_INCREMENT_FLAG`, `NOT_NULL_FLAG`, `FIELD_IS_INVISIBLE`); create-table-gipks.html "The generated key column has the name `my_row_id`, which means that a table cannot have any existing column named `my_row_id`."

---

### 24.3 GIPK NOT added when table already has a user-defined PK or PK-equivalent

**Setup:**
```sql
SET SESSION sql_generate_invisible_primary_key = ON;
CREATE TABLE t1 (id INT PRIMARY KEY, a INT);
CREATE TABLE t2 (id INT NOT NULL UNIQUE, a INT);         -- NOT-NULL UNIQUE is not promoted to PK here
CREATE TABLE t3 (id INT AUTO_INCREMENT, a INT, PRIMARY KEY (id));
```

**Oracle verification:**
```sql
SET SESSION show_gipk_in_create_table_and_information_schema = ON;
SHOW CREATE TABLE t1;        -- no my_row_id (user PK exists)
SHOW CREATE TABLE t2;        -- my_row_id PRESENT (UNIQUE NOT NULL does not suppress GIPK)
SHOW CREATE TABLE t3;        -- no my_row_id (user PK exists)
```

**omni assertion:** catalog suppresses GIPK only when the CREATE TABLE already declares a `PRIMARY KEY` (column-level or table-level). A `UNIQUE NOT NULL` column does NOT suppress GIPK — unlike the legacy "first UNIQUE NOT NULL becomes the clustered key" behavior used by InnoDB internally, the GIPK generation path tests strictly for a declared PRIMARY KEY.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `is_generated_invisible_primary_key_mode_active()` + `check_generated_invisible_primary_key()` — condition is `!create_info->primary_key` (declared PK), does not inspect UNIQUE keys; create-table-gipks.html "A GIPK is added if the `sql_generate_invisible_primary_key` system variable is enabled and the table being created does not have an explicit primary key."

---

### 24.4 GIPK interaction with column name collision

**Setup:**
```sql
SET SESSION sql_generate_invisible_primary_key = ON;
CREATE TABLE t (my_row_id INT, a INT);
```

**Oracle verification:**
```sql
-- fails at CREATE time
-- ERROR 4108 (HY000): Failed to generate invisible primary key. Column 'my_row_id' already exists.
```

**omni assertion:** catalog surface-loads CREATE TABLE with both `sql_generate_invisible_primary_key=ON` and a user-declared `my_row_id` column → hard error (`ER_GIPK_FAILED_AUTOINC_COLUMN_NAME_RESERVED` class). Under `sql_generate_invisible_primary_key=OFF`, the table is created normally with the user's `my_row_id` column. The reservation is only active when GIPK generation is on.

**Priority:** LOW
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `check_generated_invisible_primary_key()` — scans `alter_info->create_list` for a column named `my_row_id` and rejects; create-table-gipks.html "A table cannot have any existing column named my_row_id [when GIPK is enabled]."

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

**Priority:** LOW
**Status:** verified (`C25_1_decimal_default_10_0`)
**Source anchors:** `sql/field.cc` `Field_new_decimal`, default `precision=10`, `dec=0`.

---

### 25.2 DECIMAL precision-only → scale defaults to 0

**Setup:**
```sql
CREATE TABLE t (
  d5     DECIMAL(5),
  d15    DECIMAL(15),
  d65    DECIMAL(65)
);
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY ORDINAL_POSITION;
```
Expected: `decimal(5,0)`, `decimal(15,0)`, `decimal(65,0)`. SHOW CREATE TABLE renders `DECIMAL(5,0)`, never `DECIMAL(5)`. NUMERIC_SCALE = 0.

**omni assertion:** parser/catalog stores an explicit `scale=0` when scale is omitted; deparser always emits the `,0` form so round-trip through SDL matches MySQL's canonical SHOW CREATE TABLE output.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/field.cc` Field_new_decimal construction path; fixed-point-types.html "DECIMAL(M): This is a synonym for DECIMAL(M,0)."; precision-math-decimal-characteristics.html "If D is omitted, the default is 0."

---

### 25.3 DECIMAL max precision 65 / max scale 30 / scale > precision

**Setup:**
```sql
-- each subcase parsed independently
CREATE TABLE ok_max_p (d DECIMAL(65, 30));                 -- OK
CREATE TABLE err_p_gt_65 (d DECIMAL(66, 0));               -- error: precision out of range
CREATE TABLE err_s_gt_30 (d DECIMAL(40, 31));              -- error: scale out of range
CREATE TABLE err_s_gt_p  (d DECIMAL(5, 6));                -- error: scale larger than precision
CREATE TABLE err_neg     (d DECIMAL(-1, 0));               -- parse error
```

**Oracle verification:**
```sql
-- Expected MySQL errors:
-- ER_TOO_BIG_PRECISION       (1426) for DECIMAL(66,0)    "Too big precision 66 specified for column 'd'. Maximum is 65."
-- ER_TOO_BIG_SCALE           (1425) for DECIMAL(40,31)   "Too big scale 31 specified for column 'd'. Maximum is 30."
-- ER_M_BIGGER_THAN_D         (1427) for DECIMAL(5,6)     "For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'd')."
```

**omni assertion:** parser / catalog validate DECIMAL bounds with the three distinct error codes. The max-allowed `DECIMAL(65,30)` must be accepted as-is (this is a common real-world schema anti-pattern). `DECIMAL(65,30)` stores up to 35 integer digits and 30 fractional digits.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/field.cc` Field_new_decimal constants `DECIMAL_MAX_PRECISION=65`, `DECIMAL_MAX_SCALE=30`; `sql/sql_yacc.yy` field_options validation; precision-math-decimal-characteristics.html "The maximum number of digits (M) is 65. The maximum number of decimals (D) is 30. If D is omitted, the default is 0. If M is omitted, the default is 10."

---

### 25.4 UNSIGNED DECIMAL keeps precision/scale, just removes sign bit

**Setup:**
```sql
CREATE TABLE t (
  d1 DECIMAL(10,2) UNSIGNED,
  d2 DECIMAL(10,2) UNSIGNED ZEROFILL,
  d3 NUMERIC(10,2) UNSIGNED
);
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, COLUMN_TYPE, NUMERIC_PRECISION, NUMERIC_SCALE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
ORDER BY ORDINAL_POSITION;
SHOW WARNINGS;                -- ZEROFILL and DISPLAY WIDTH deprecation warnings
```
Expected: `decimal(10,2) unsigned`, `decimal(10,2) unsigned zerofill`, `decimal(10,2) unsigned`. NUMERIC_PRECISION=10, NUMERIC_SCALE=2 for all three. Deprecation warning `ER_WARN_DEPRECATED_ZEROFILL_ETC` for `d2`. NUMERIC is a stored-as-DECIMAL synonym and displays as `decimal` in information_schema.

**omni assertion:** catalog stores `Unsigned` flag separately from precision/scale; deparser emits `UNSIGNED` (and `ZEROFILL` if set) in that order AFTER the `(M,D)` spec. Round-trip preserves the NUMERIC vs DECIMAL choice at the token level IF the catalog records the spelled form, otherwise it canonicalizes to DECIMAL (MySQL canonicalizes — the information_schema always reports `decimal`).

**Priority:** LOW
**Status:** pending
**Source anchors:** `sql/field.cc` Field_new_decimal `unsigned_flag`; `sql/sql_yacc.yy` `NUMERIC_SYM` reduces to DECIMAL type; fixed-point-types.html "NUMERIC is implemented as DECIMAL, so the following remarks about DECIMAL apply equally to NUMERIC."

---

### 25.5 Zero-scale rendering: DECIMAL(5,0) vs DECIMAL(5) vs DECIMAL

**Setup:**
```sql
CREATE TABLE t (
  a DECIMAL,
  b DECIMAL(5),
  c DECIMAL(5,0),
  d DECIMAL(10,0)
);
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
```
Expected: SHOW CREATE TABLE renders `decimal(10,0)`, `decimal(5,0)`, `decimal(5,0)`, `decimal(10,0)`. All three input forms for a zero-scale column collapse to the fully-qualified `DECIMAL(M,0)` rendering in both SHOW CREATE TABLE and information_schema.

**omni assertion:** deparser always renders scale explicitly, even zero, even when the original SQL omitted it. SDL diff must therefore treat `DECIMAL`, `DECIMAL(10)`, and `DECIMAL(10,0)` as equivalent — do NOT emit a type-change ALTER when the user source had the short form but the loaded catalog reports the expanded form.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_show.cc` `show_create_table` column type printer; `sql/field.cc` `Field_new_decimal::sql_type()` emits `(precision,dec)` unconditionally.

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

## Section AX: ALTER TABLE sub-commands implicit behaviors

### AX.1 ADD COLUMN without FIRST/AFTER → appended at end

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, c INT);
ALTER TABLE t ADD COLUMN d INT;
ALTER TABLE t ADD COLUMN e INT FIRST;
ALTER TABLE t ADD COLUMN f INT AFTER a;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
```
Expected order after all three ALTERs: `e, a, f, b, c, d`. A bare ADD COLUMN with no position clause always appends.

**omni assertion:** catalog's ALTER-apply helper for `ADD COLUMN` must insert at the end of the current column list when no FIRST/AFTER clause is present. FIRST inserts at index 0; AFTER x inserts at `index(x)+1`. The inserted column's ordinal is the new count-1 for the plain append path.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_alter_table` — iterates new_create_list and places Alter_column entries with `after` = NULL → appended; alter-table.html "To add a column at a specific position within a table row, use FIRST or AFTER col_name. The default is to add the column last."

---

### AX.2 DROP COLUMN cascades: removes indexes containing the dropped column

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, c INT, INDEX idx_a(a), INDEX idx_ab(a, b), INDEX idx_bc(b, c));
ALTER TABLE t DROP COLUMN a;
```

**Oracle verification:**
```sql
SHOW INDEX FROM t;
SELECT INDEX_NAME, COLUMN_NAME, SEQ_IN_INDEX FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY INDEX_NAME, SEQ_IN_INDEX;
```
Expected after drop: `idx_a` is gone (single-column index on dropped column). `idx_ab` is also gone — MySQL drops the whole composite index when any of its columns is dropped. `idx_bc` remains intact.

**omni assertion:** catalog's `DROP COLUMN` helper iterates the table's indexes; any index listing the dropped column among its keyparts is removed in entirety (no auto-trim of the composite index to the surviving columns).

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_alter_table` — `drop_list` handling for indexes: when `fill_alter_inplace_info` sees a key with a dropped column, the key is added to `dropped_keys`; alter-table.html "If columns are dropped from a table, the columns are also removed from any index of which they are a part. If all columns that make up an index are dropped, the index is dropped as well."

---

### AX.3 DROP COLUMN fails if it would leave the table empty

**Setup:**
```sql
CREATE TABLE t (a INT);
ALTER TABLE t DROP COLUMN a;
```

**Oracle verification:**
```sql
-- ERROR 1090 (42000): ER_CANT_REMOVE_ALL_FIELDS
-- "You can't delete all columns with ALTER TABLE; use DROP TABLE instead"
```

**omni assertion:** catalog ALTER-apply rejects a DROP COLUMN that would remove the only remaining column with `ER_CANT_REMOVE_ALL_FIELDS`. Check happens after resolving the drop list, before applying any other sub-commands in the same ALTER (so `ALTER TABLE t DROP COLUMN a, ADD COLUMN b INT` ALSO fails — the final shape has one column but MySQL rejects at prepare time).

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_alter_table` — "if (!new_create_list.elements) { my_error(ER_CANT_REMOVE_ALL_FIELDS, MYF(0)); goto err; }"; alter-table.html "You cannot drop all columns from a table."

---

### AX.4 DROP COLUMN fails if referenced by CHECK or GENERATED column

**Setup:**
```sql
CREATE TABLE t1 (a INT, b INT, CHECK (a > 0));
ALTER TABLE t1 DROP COLUMN a;

CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a + 1));
ALTER TABLE t2 DROP COLUMN a;
```

**Oracle verification:**
```sql
-- t1:  ER_CHECK_CONSTRAINT_REFERS (3942)   "Check constraint '...' uses column 'a', hence column cannot be dropped or renamed."
-- t2:  ER_DEPENDENT_BY_GENERATED_COLUMN (3108)   "Column 'a' has a generated column dependency."
```

**omni assertion:** catalog's DROP COLUMN helper must scan (a) CHECK constraint expressions and (b) generated column expressions for references to the dropped column, and raise the appropriate error before applying the drop. The check also applies to CHANGE COLUMN that would rename the column — see AX.6.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `check_if_field_used_by_generated_column_or_default()` / check-constraint validation in `prepare_check_constraints_for_alter()`; alter-table.html "You cannot drop or rename a column that is referenced by a generated column or CHECK constraint."

---

### AX.5 MODIFY COLUMN preserves unspecified attributes? → NO, it rewrites from the given spec

**Setup:**
```sql
CREATE TABLE t (
  a INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  b INT NOT NULL DEFAULT 5 COMMENT 'x'
);
ALTER TABLE t MODIFY COLUMN b BIGINT;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b';
```
Expected: column `b` becomes `bigint DEFAULT NULL` with NO comment, NO NOT NULL, NO default. MODIFY COLUMN treats its spec as the NEW full column definition; attributes absent from the ALTER statement are NOT inherited from the old column.

**omni assertion:** catalog `MODIFY COLUMN` helper replaces the old column spec wholesale with the new one. Attribute inheritance from the prior column is NOT performed. Exception: column position is preserved unless FIRST/AFTER is specified.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_alter_table` — `Alter_info::ALTER_CHANGE_COLUMN` branch constructs a fresh `Create_field` from the new spec; alter-table.html "When you use CHANGE or MODIFY, the column definition must include the data type and all attributes that should apply to the new column, other than index attributes such as PRIMARY KEY or UNIQUE. Attributes present in the original definition but not specified for the new definition are not carried over."

---

### AX.6 CHANGE COLUMN renames + retypes atomically; MODIFY cannot rename

**Setup:**
```sql
CREATE TABLE t (a INT);
ALTER TABLE t CHANGE COLUMN a b BIGINT NOT NULL;    -- rename + retype
ALTER TABLE t MODIFY COLUMN b b INT;                -- new_name == old_name is fine (retype only)
-- ALTER TABLE t MODIFY COLUMN b c INT;             -- syntax error: MODIFY takes only one name
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
```
Expected: column renames work only via CHANGE. MODIFY parses as `MODIFY [COLUMN] col_name column_definition` — exactly one name, no rename.

**omni assertion:** parser accepts `CHANGE old new def` (two names) and `MODIFY col def` (one name). Catalog's CHANGE helper matches on the OLD name and then applies the NEW name + NEW type atomically. If the new name collides with another existing column (not the one being renamed), raise `ER_DUP_FIELDNAME`.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_yacc.yy` `alter_list_item: CHANGE ... field_ident field_ident ...` vs `MODIFY ... field_ident ...`; alter-table.html "CHANGE requires two column names and can rename a column. MODIFY does not rename."

---

### AX.7 ADD INDEX without name → auto-generated name (ALTER context)

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, INDEX (a));
ALTER TABLE t ADD INDEX (b);
ALTER TABLE t ADD INDEX (a, b);
```

**Oracle verification:**
```sql
SHOW INDEX FROM t;
```
Expected index names: `a`, `b`, `a_2`. MySQL names a nameless single-column index after the first column. When that name collides it appends `_2`, `_3`, etc. The counter is per-table, per-first-column.

**omni assertion:** catalog's ALTER ADD INDEX helper applies the same naming rule as CREATE TABLE (C1 scenario) but the collision search must include (a) pre-existing indexes on the table AND (b) indexes being added earlier in the same ALTER statement. `_2` is chosen, not `_1`.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `prepare_key_column()` / `make_unique_key_name()` — loops appending numeric suffixes; alter-table.html links to CREATE TABLE index-name rules.

---

### AX.8 ADD UNIQUE / ADD KEY / ADD FULLTEXT implicit index types

**Setup:**
```sql
CREATE TABLE t (a INT, b VARCHAR(100)) ENGINE=InnoDB;
ALTER TABLE t ADD UNIQUE (a);
ALTER TABLE t ADD KEY (b);
ALTER TABLE t ADD FULLTEXT (b);
ALTER TABLE t ADD SPATIAL KEY (b);   -- error: SPATIAL requires NOT NULL geometry
```

**Oracle verification:**
```sql
SELECT INDEX_NAME, NON_UNIQUE, INDEX_TYPE FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY INDEX_NAME;
```
Expected: UNIQUE index → `NON_UNIQUE=0`, `INDEX_TYPE='BTREE'`. Plain KEY → `NON_UNIQUE=1`, `INDEX_TYPE='BTREE'` (InnoDB). FULLTEXT → `INDEX_TYPE='FULLTEXT'`, `NON_UNIQUE=1`. SPATIAL is rejected on nullable column (`ER_SPATIAL_CANT_HAVE_NULL`).

**omni assertion:** catalog translates each ADD sub-command to its constraint kind: `UNIQUE`, `INDEX`/`KEY`, `FULLTEXT INDEX`, `SPATIAL INDEX`. Default `USING` algorithm is `BTREE` for UNIQUE/KEY on InnoDB, regardless of the `USING` clause optional. FULLTEXT and SPATIAL ignore `USING`.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_create_table` key-type mapping; alter-table.html "add_constraint : ADD [CONSTRAINT [symbol]] {PRIMARY KEY | UNIQUE | FOREIGN KEY} ...".

---

### AX.9 ADD FOREIGN KEY: column-level shorthand vs table-level form

**Setup:**
```sql
CREATE TABLE parent (id INT PRIMARY KEY);
CREATE TABLE t1 (a INT);
ALTER TABLE t1 ADD FOREIGN KEY (a) REFERENCES parent(id);     -- table-level
-- ALTER TABLE cannot use column-level shorthand — only table-level ADD is accepted
CREATE TABLE t2 (a INT REFERENCES parent(id));                -- column-level shorthand in CREATE
```

**Oracle verification:**
```sql
SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME IN ('t1','t2');
```
Expected: `t1` has a FK constraint named `t1_ibfk_1` on `(a) → parent(id)`. `t2` has NO FK — MySQL silently ignores column-level `REFERENCES` clauses in CREATE TABLE for non-FEDERATED engines. This is a notorious pitfall.

**omni assertion:** catalog's CREATE/ALTER helpers must drop the column-level `REFERENCES` clause for InnoDB (and MyISAM) without raising an error, matching MySQL's lenient parse-but-ignore semantics. Only a table-level `FOREIGN KEY (col) REFERENCES ...` actually creates an FK. ALTER TABLE does not accept column-level shorthand; a FK must be added via table-level `ADD CONSTRAINT ... FOREIGN KEY`.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `add_fk_to_list` and `mysql_prepare_create_table`; create-table-foreign-keys.html "MySQL parses but ignores 'inline REFERENCES specifications' (as defined in the SQL standard) where the references are defined as part of the column specification. MySQL accepts REFERENCES clauses only when specified as part of a separate FOREIGN KEY specification."

---

### AX.10 RENAME COLUMN syntax

**Setup:**
```sql
CREATE TABLE t (a INT NOT NULL DEFAULT 5 COMMENT 'hi', INDEX (a));
ALTER TABLE t RENAME COLUMN a TO aa;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;
SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
SHOW INDEX FROM t;
```
Expected: column is renamed to `aa`, all attributes (NOT NULL, DEFAULT 5, COMMENT) PRESERVED (unlike CHANGE/MODIFY, RENAME COLUMN does not rewrite the spec). Index name stays `a` (the index auto-name does NOT track the column rename).

**omni assertion:** catalog `RENAME COLUMN old TO new` helper mutates only the column name; all other attributes are untouched. Indexes referencing the column are updated to use the new name as a keypart but the INDEX NAME itself (if auto-generated from the old column name) is NOT renamed. Check-constraint expressions and generated column expressions are updated to refer to the new column name.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_rename_column` path via Alter_info `flags & ALTER_RENAME_COLUMN`; alter-table.html "RENAME COLUMN old_col_name TO new_col_name. RENAME COLUMN does not otherwise change the column definition."

---

### AX.11 RENAME INDEX syntax

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, INDEX old_name (a, b));
ALTER TABLE t RENAME INDEX old_name TO new_name;
```

**Oracle verification:**
```sql
SHOW INDEX FROM t;
```
Expected: index name is `new_name`, columns and order unchanged.

**omni assertion:** catalog supports `RENAME INDEX` as a pure metadata rename. PRIMARY cannot be renamed (error `ER_WRONG_NAME_FOR_INDEX` or similar if user writes `RENAME INDEX PRIMARY TO foo`). Duplicate target name → `ER_DUP_KEYNAME`.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_alter.h` `Alter_info::ALTER_RENAME_INDEX` flag; `sql/sql_table.cc` `rename_index_in_dd_table`; alter-table.html "RENAME {INDEX|KEY} old_index_name TO new_index_name. Renames an index."

---

### AX.12 RENAME TO (table rename via ALTER)

**Setup:**
```sql
CREATE TABLE t (a INT);
ALTER TABLE t RENAME TO t2;
ALTER TABLE t2 RENAME TO other_db.t3;       -- cross-db move if same filesystem
```

**Oracle verification:**
```sql
SHOW TABLES;
SHOW TABLES FROM other_db;
```

**omni assertion:** catalog recognizes `ALTER TABLE ... RENAME TO [db.]new_name` as equivalent to `RENAME TABLE`. FKs pointing at the old name are updated atomically. A cross-database rename is supported when both DBs are on the same filesystem; views, triggers, stored routines referencing the old name are NOT automatically updated.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_rename_table` invoked from ALTER path when `Alter_info::ALTER_RENAME` flag set; alter-table.html "RENAME [TO|AS] new_tbl_name. Renames the table."

---

### AX.13 ALTER COLUMN SET DEFAULT / DROP DEFAULT / SET VISIBLE / SET INVISIBLE

**Setup:**
```sql
CREATE TABLE t (a INT NOT NULL, b INT);
ALTER TABLE t ALTER COLUMN a SET DEFAULT 5;
ALTER TABLE t ALTER COLUMN a DROP DEFAULT;
ALTER TABLE t ALTER COLUMN b SET INVISIBLE;
ALTER TABLE t ALTER COLUMN b SET VISIBLE;
```

**Oracle verification:**
```sql
SHOW CREATE TABLE t;                                   -- after each step
SELECT COLUMN_NAME, COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
```
Expected: SET DEFAULT mutates only the default literal; DROP DEFAULT removes it and (for integer/non-timestamp non-nullable columns) leaves no default — inserting without a value becomes an error in strict mode. SET INVISIBLE / SET VISIBLE toggles the `INVISIBLE` flag in `EXTRA`, no effect on the stored value.

**omni assertion:** catalog's `ALTER COLUMN ... {SET DEFAULT expr | DROP DEFAULT | SET INVISIBLE | SET VISIBLE}` helper performs an in-place metadata patch: no type/nullability changes, no index rebuild, no data rewrite. SET DEFAULT accepts both literals and (as of 8.0.13) expressions wrapped in parentheses.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_alter.h` `Alter_column` class with `m_default_val`, `m_is_set_default`, `m_is_visible`; alter-table.html "ALTER COLUMN col_name {SET DEFAULT literal | (expr) | DROP DEFAULT | SET {VISIBLE | INVISIBLE}}."

---

### AX.14 ALTER INDEX ... VISIBLE / INVISIBLE

**Setup:**
```sql
CREATE TABLE t (a INT, b INT, INDEX idx_a (a));
ALTER TABLE t ALTER INDEX idx_a INVISIBLE;
```

**Oracle verification:**
```sql
SHOW INDEX FROM t;                                      -- Visible column = 'NO'
SELECT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t';
-- PRIMARY cannot be made invisible:
-- ALTER TABLE t ALTER INDEX `PRIMARY` INVISIBLE;   -- ER_PK_INDEX_CANT_BE_INVISIBLE
```

**omni assertion:** catalog toggles the `is_visible` flag on the index. Query planner consumers of the catalog must honor visibility (invisible = ignored unless `optimizer_switch='use_invisible_indexes=on'`). PRIMARY KEY's underlying index cannot be made invisible — reject with `ER_PK_INDEX_CANT_BE_INVISIBLE`.

**Priority:** MED
**Status:** pending
**Source anchors:** `sql/sql_alter.h` `Alter_index_visibility`; invisible-indexes.html "The PRIMARY (explicit or implicit) cannot be made invisible."; alter-table.html "ALTER {INDEX|KEY} index_name {VISIBLE | INVISIBLE}."

---

### AX.15 Multi-sub-command ALTER: single-pass semantics and ordering

**Setup:**
```sql
CREATE TABLE t (a INT, b INT);
ALTER TABLE t
  ADD COLUMN c INT AFTER a,
  DROP COLUMN b,
  ADD INDEX (c),
  RENAME COLUMN a TO aa;
```

**Oracle verification:**
```sql
SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY ORDINAL_POSITION;
SHOW INDEX FROM t;
```
Expected: columns `aa, c`; index `c` exists. All sub-commands apply against the PRE-ALTER schema — so `ADD COLUMN c AFTER a` uses `a` even though a later clause renames it, and `DROP COLUMN b` does not interfere with `ADD COLUMN c`. MySQL evaluates the full sub-command list in one logical pass: drops and adds are resolved against the old schema; the new column list is then computed; the rename applies after the add/drop.

**omni assertion:** catalog's ALTER helper does NOT execute sub-commands one at a time in surface order. Instead it (a) collects drops, adds, changes, renames, (b) builds the new column list by iterating the OLD list, filtering drops, applying renames/changes, and injecting adds with their FIRST/AFTER targets resolved against the OLD-name space, (c) applies index-level sub-commands against the resulting column list. Contradictory combinations (e.g. `ADD COLUMN x, DROP COLUMN x`) are rejected with `ER_CANT_DROP_FIELD_OR_KEY` on the drop side. This single-pass semantics is critical for SDL diff emission — diff must produce ONE ALTER containing all the needed sub-commands, not a sequence.

**Priority:** HIGH
**Status:** pending
**Source anchors:** `sql/sql_table.cc` `mysql_prepare_alter_table` — builds `new_create_list`, `new_key_list` in a single pass from the old definition + Alter_info lists; alter-table.html "You can issue multiple ADD, ALTER, DROP, and CHANGE clauses in a single ALTER TABLE statement, separated by commas. This is a MySQL extension to standard SQL, which permits only one of each clause per ALTER TABLE statement."

---

## Progress tracker

Status values: `pending`, `verified` (spot-check done), `passing`, `bug` (omni bug known), `limitation`.

| # | Scenario | Status | Priority | Notes |
|---|----------|--------|----------|-------|
| 1.1 | FK name CREATE fresh counter | verified | HIGH | `C1_1_FK_counter_max_plus_one` + bugfix test |
| 1.2 | FK name ALTER max+1 | verified | HIGH | spot-check |
| 1.3 | Partition default naming p0..pN | verified | MED | `C1_2_partition_naming` |
| 1.4 | CHECK constraint auto-name | verified | MED | `C1_3_CheckConstraintName` |
| 1.5 | UNIQUE KEY field-name | verified | LOW | `TestScenario_C1` reconciliation |
| 1.6 | UNIQUE KEY name collision _2 | verified | LOW | `TestScenario_C1` reconciliation |
| 1.7 | PRIMARY KEY name forced "PRIMARY" | verified | HIGH | `TestScenario_C1/1_7` reconciliation |
| 1.8 | Non-PK index cannot be named PRIMARY | verified | MED | `TestScenario_C1` reconciliation |
| 1.9 | Implicit index name from first column | verified | HIGH | `TestScenario_C1/1_9` reconciliation |
| 1.10 | UNIQUE fallback when column literally "PRIMARY" | verified | LOW | `TestScenario_C1` reconciliation |
| 1.11 | Functional index auto-name collision suffix | verified | MED | `TestScenario_C1/1_11` |
| 1.12 | Functional index hidden column name format | verified | MED | `TestScenario_C1/1_12` |
| 1.13 | CHECK constraint name is schema-scoped | verified | HIGH | `TestScenario_C1/1_13` reconciliation |
| 2.1 | REAL → DOUBLE | verified | LOW | `C2_1_REAL_to_DOUBLE` |
| 2.2 | BOOL → TINYINT(1) | verified | LOW | `C2_2_BOOL_to_TINYINT1` |
| 2.3 | INTEGER → INT | verified | LOW | `TestScenario_C2` reconciliation |
| 2.4 | BOOLEAN → TINYINT(1) | verified | LOW | `TestScenario_C2` reconciliation |
| 2.5 | INT1/INT2/INT3/INT4/INT8 aliases | verified | LOW | `TestScenario_C2` reconciliation |
| 2.6 | MIDDLEINT → MEDIUMINT | verified | LOW | `TestScenario_C2` reconciliation |
| 2.7 | INT(11) display width stripped | verified | HIGH | `TestScenario_C2` reconciliation |
| 2.8 | INT(N) ZEROFILL preserves width | verified | HIGH | `TestScenario_C2/2_8` reconciliation |
| 2.9 | SERIAL composite alias | verified | HIGH | `TestScenario_C2` reconciliation |
| 2.10 | NUMERIC → DECIMAL | verified | LOW | `TestScenario_C2` reconciliation |
| 2.11 | DEC / FIXED → DECIMAL | verified | LOW | `TestScenario_C2` reconciliation |
| 2.12 | DOUBLE PRECISION → DOUBLE | verified | LOW | `TestScenario_C2` reconciliation |
| 2.13 | FLOAT4 → FLOAT, FLOAT8 → DOUBLE | verified | LOW | `TestScenario_C2` reconciliation |
| 2.14 | FLOAT(p) precision split | verified | HIGH | `TestScenario_C2/2_14` reconciliation |
| 2.15 | FLOAT(M,D) deprecated non-standard | verified | MED | `TestScenario_C2` reconciliation |
| 2.16 | CHARACTER / CHARACTER VARYING | verified | LOW | `TestScenario_C2` reconciliation |
| 2.17 | NATIONAL CHAR / NCHAR utf8mb3 | verified | MED | `TestScenario_C2` reconciliation |
| 2.18 | NVARCHAR / NCHAR VARYING utf8mb3 | verified | MED | `TestScenario_C2` reconciliation |
| 2.19 | LONG / LONG VARCHAR → MEDIUMTEXT | verified | LOW | `TestScenario_C2` reconciliation |
| 2.20 | CHAR / BINARY default length 1 | verified | MED | `TestScenario_C2` reconciliation |
| 2.21 | VARCHAR without length syntax error | verified | MED | `TestScenario_C2` reconciliation |
| 2.22 | TIMESTAMP/DATETIME/TIME default fsp 0 | verified | MED | `TestScenario_C2` reconciliation |
| 2.23 | YEAR(4) deprecated → bare YEAR | verified | MED | `TestScenario_C2` reconciliation |
| 2.24 | BIT without length defaults to BIT(1) | verified | HIGH | `TestScenario_C2/2_24` reconciliation |
| 2.25 | VARCHAR > 65535 → TEXT family | verified | MED | `TestScenario_C2` reconciliation |
| 2.26 | TEXT(N)/BLOB(N) byte-count promotion | verified | HIGH | `TestScenario_C2/2_26` reconciliation |
| 3.1 | TIMESTAMP first-only promotion | verified | HIGH | `C3_1_timestamp_first_only_promotion` — omni risk |
| 3.2 | PRIMARY KEY → NOT NULL | verified | MED | `C3_2_primary_key_implies_not_null` |
| 3.3 | AUTO_INCREMENT → NOT NULL | verified | MED | `C3_3_AutoIncrement_implies_NOT_NULL` |
| 4.1 | Table charset from database | verified | HIGH | `C4_1_table_charset_from_database` |
| 4.2 | Column charset inherits + elided | verified | MED | `C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision` |
| 4.3 | Column COLLATE alone derives CHARACTER SET | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.4 | Column CHARACTER SET alone derives default collation | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.5 | Table COLLATE/CHARSET mismatch rejected | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.6 | BINARY attribute → {charset}_bin collation | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.7 | CHARACTER SET binary vs BINARY type | verified | MED | `TestScenario_C4` reconciliation |
| 4.8 | utf8 alias normalized to utf8mb3 | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.9 | NATIONAL / NCHAR hard-wired to utf8mb3 | verified | MED | `TestScenario_C4` reconciliation |
| 4.10 | ENUM/SET charset inheritance | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.11 | Index prefix length × mbmaxlen | verified | HIGH | `TestScenario_C4` reconciliation |
| 4.12 | Collation derivation aggregation | verified | MED | `TestScenario_C4` reconciliation |
| 5.1 | FK ON DELETE default NO ACTION | verified | HIGH | `C5_1_fk_on_delete_default` — reporting discrepancy |
| 5.2 | FK SET NULL on NOT NULL errors | verified | MED | `TestScenario_C5` reconciliation |
| 5.3 | FK MATCH default NONE | verified | MED | `C5_3_FK_MATCH_default` |
| 6.1 | HASH PARTITIONS defaults to 1 | verified | MED | `C6_1_Partition_default_count` |
| 6.2 | Subpartition count default 1 | verified | MED | `C6_2_Subpartition_count` |
| 6.3 | Partition engine inherits | verified | LOW | `TestScenario_C6` reconciliation |
| 6.4 | KEY partitioning ALGORITHM defaults to 2 | verified | HIGH | `TestScenario_C6/6_4` reconciliation |
| 6.5 | KEY() empty list defaults to PRIMARY KEY cols | verified | HIGH | `TestScenario_C6/6_5` reconciliation |
| 6.6 | LINEAR HASH / LINEAR KEY preserved | verified | LOW | `TestScenario_C6` reconciliation |
| 6.7 | RANGE/LIST require explicit definitions | verified | HIGH | `TestScenario_C6/6_7` reconciliation |
| 6.8 | MAXVALUE only in last RANGE partition | verified | MED | `TestScenario_C6` reconciliation |
| 6.9 | LIST equality vs RANGE strict less-than | verified | LOW | `TestScenario_C6` reconciliation |
| 6.10 | LIST DEFAULT partition catch-all | pending | MED | Wave 1 C6 worker |
| 6.11 | Partition function must return INTEGER | verified | HIGH | `TestScenario_C6/6_11` reconciliation |
| 6.12 | TIMESTAMP requires UNIX_TIMESTAMP() wrap | verified | MED | `TestScenario_C6` reconciliation |
| 6.13 | UNIQUE/PK must contain partition cols | verified | HIGH | `TestScenario_C6/6_13` reconciliation |
| 6.14 | Per-partition options not inherited | verified | MED | `TestScenario_C6` reconciliation |
| 6.15 | Subpartition options inherit from parent | verified | LOW | `TestScenario_C6` reconciliation |
| 6.16 | ADD PARTITION auto-naming seeds from count | verified | LOW | `TestScenario_C6` reconciliation |
| 6.17 | COALESCE/REORGANIZE partition behavior | verified | LOW | `TestScenario_C6` reconciliation |
| 7.1 | Default index algorithm BTREE | verified | MED | `C7_1_Default_index_algorithm_BTREE` |
| 7.2 | FK backing index implicit | verified | MED | `C7_2_FK_backing_index` |
| 7.3 | HASH on InnoDB silently coerced to BTREE | verified | MED | `TestScenario_C7` reconciliation |
| 7.4 | USING BTREE explicit vs implicit rendering | verified | MED | `TestScenario_C7` reconciliation |
| 7.5 | UNIQUE on nullable allows multiple NULLs | verified | HIGH | `TestScenario_C7` reconciliation |
| 7.6 | VISIBLE default; PK cannot be INVISIBLE | verified | MED | `TestScenario_C7` reconciliation |
| 7.7 | BLOB/TEXT index requires prefix length | verified | HIGH | `TestScenario_C7` reconciliation |
| 7.8 | FULLTEXT default parser when WITH PARSER omitted | verified | MED | `TestScenario_C7` reconciliation |
| 7.9 | SPATIAL requires NOT NULL; no USING BTREE/HASH | verified | HIGH | `TestScenario_C7` reconciliation |
| 7.10 | PRIMARY and UNIQUE on same columns both persist | verified | MED | `TestScenario_C7` reconciliation |
| 8.1 | Engine default InnoDB | verified | MED | `C8_1_Default_engine_InnoDB` |
| 8.2 | ROW_FORMAT default DYNAMIC | verified | MED | `C8_2_Default_row_format` |
| 8.3 | AUTO_INCREMENT starts at 1 | verified | MED | covered by 18.4 spot-check |
| 9.1 | Generated column VIRTUAL default | verified | LOW | `C9_1_GeneratedColumn_default_VIRTUAL` |
| 9.2 | FK on generated column errors | verified | LOW | `TestScenario_C9` reconciliation |
| 10.1 | View ALGORITHM/CHECK/SECURITY defaults | verified | MED | `C10_1_3_4_View_defaults` |
| 10.2 | View DEFINER default | verified | MED | `C10_2_view_security_definer` |
| 11.1 | Trigger DEFINER default | verified | LOW | `C11_1_Trigger_definer_default` |
| 14.1 | CHECK enforced default | verified | LOW | `C14_1_Check_enforced_default` |
| 15.1 | ADD COLUMN appends to end | verified | LOW | `C15_1_Column_positioning_end` |
| 16.1 | NOW() precision default 0 | verified | MED | `C16_1_now_precision_default_zero` |
| 16.2 | NOW(N) explicit precision 0..6 | verified | HIGH | `TestScenario_C16` reconciliation |
| 16.3 | CURDATE / UTC_DATE take no precision | verified | MED | `TestScenario_C16` reconciliation |
| 16.4 | CURTIME / UTC_TIME precision defaults to 0 | verified | MED | `TestScenario_C16` reconciliation |
| 16.5 | SYSDATE precision defaults to 0 | verified | LOW | `TestScenario_C16` reconciliation |
| 16.6 | UTC_TIMESTAMP precision defaults to 0 | verified | LOW | `TestScenario_C16` reconciliation |
| 16.7 | UNIX_TIMESTAMP return type by arg fsp | verified | MED | `TestScenario_C16` reconciliation |
| 16.8 | DATETIME(N) DEFAULT NOW() fsp must match | verified | HIGH | `TestScenario_C16` reconciliation |
| 16.9 | ON UPDATE NOW(N) must match column fsp | verified | HIGH | `TestScenario_C16` reconciliation |
| 16.10 | DATETIME/TIMESTAMP storage bytes by fsp | verified | MED | `TestScenario_C16` reconciliation |
| 16.11 | YEAR(N) deprecated — only YEAR(4) accepted | verified | LOW | `TestScenario_C16` reconciliation |
| 16.12 | TIMESTAMP promotion inherits fsp | pending | MED | Wave 1 C16 worker |
| 17.1 | CONCAT two cols identical charset/collation | verified | HIGH | `TestScenario_C17` reconciliation |
| 17.2 | CONCAT latin1 + utf8mb4 superset conv | verified | HIGH | `TestScenario_C17` reconciliation |
| 17.3 | CONCAT incompatible collations (ER 1267) | verified | HIGH | `TestScenario_C17` reconciliation |
| 17.4 | CONCAT_WS separator charset + NULL skip | verified | MED | `TestScenario_C17` reconciliation |
| 17.5 | String literal coercibility `_utf8mb4'x'` | verified | MED | `TestScenario_C17` reconciliation |
| 17.6 | REPEAT/LPAD/RPAD first-arg pins charset | verified | MED | `TestScenario_C17` reconciliation |
| 17.7 | CONVERT(x USING cs) forces IMPLICIT cs | verified | HIGH | `TestScenario_C17/17_7` reconciliation |
| 17.8 | COLLATE clause is EXPLICIT | verified | HIGH | `TestScenario_C17` reconciliation |
| 18.1 | Column charset elision | verified | HIGH | `C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision` |
| 18.2 | NOT NULL elision (TIMESTAMP) | verified | HIGH | `C18_2_NotNull_rendering` |
| 18.3 | ENGINE always rendered | verified | HIGH | `TestScenario_C18` reconciliation |
| 18.4 | AUTO_INCREMENT elided at 1 | verified | HIGH | `C18_4_auto_increment_elision` |
| 18.5 | DEFAULT CHARSET always rendered | verified | HIGH | `C18_5_DefaultCharset_implicit` |
| 18.6 | ROW_FORMAT elided unless explicit | verified | HIGH | `TestScenario_C18` reconciliation |
| 18.7 | Table COLLATE rendered only when non-primary | verified | HIGH | `TestScenario_C18` reconciliation |
| 18.8 | KEY_BLOCK_SIZE elided unless explicit | verified | HIGH | `TestScenario_C18` reconciliation |
| 18.9 | COMPRESSION elided unless explicit | verified | HIGH | `TestScenario_C18/18_9` reconciliation |
| 18.10 | STATS_* elision | verified | HIGH | `TestScenario_C18/18_10` reconciliation |
| 18.11 | AVG_ROW_LENGTH / MIN_ROWS / MAX_ROWS elided at zero | verified | HIGH | `TestScenario_C18/18_11` reconciliation |
| 18.12 | TABLESPACE elision | verified | HIGH | `TestScenario_C18/18_12` reconciliation |
| 18.13 | PACK_KEYS / CHECKSUM / DELAY_KEY_WRITE elision | verified | HIGH | `TestScenario_C18/18_13` reconciliation |
| 18.14 | Per-index COMMENT / KEY_BLOCK_SIZE rendering | verified | HIGH | `TestScenario_C18` reconciliation |
| 18.15 | USING BTREE/HASH emitted only when explicit | verified | HIGH | `TestScenario_C18` reconciliation |
| 19.1 | Functional index hidden VIRTUAL gen col | verified | P0 | `TestScenario_C19/19_1` |
| 19.2 | Hidden col type inferred from expression | verified | P0 | `TestScenario_C19/19_2` |
| 19.3 | Hidden col invisible to SELECT * / I_S | verified | P0 | `TestScenario_C19/19_3` |
| 19.4 | Functional expr must be deterministic/pure | verified | P1 | `TestScenario_C19/19_4` |
| 19.5 | Functional index on JSON path via CAST | verified | P0 | `TestScenario_C19/19_5` |
| 19.6 | DROP INDEX cascades to hidden gen col | verified | P0 | `TestScenario_C19/19_6` |
| 20.1 | INT NOT NULL no DEFAULT → implicit 0 | verified | HIGH | `TestScenario_C20` reconciliation |
| 20.2 | INT nullable no DEFAULT → implicit NULL | verified | HIGH | `TestScenario_C20` reconciliation |
| 20.3 | VARCHAR/CHAR NOT NULL → implicit '' | verified | MED | `TestScenario_C20` reconciliation |
| 20.4 | ENUM NOT NULL → first value default | verified | HIGH | `TestScenario_C20` reconciliation |
| 20.5 | DATETIME/DATE NOT NULL → zero-date | verified | HIGH | `TestScenario_C20` reconciliation |
| 20.6 | BLOB/TEXT literal DEFAULT error 1101 | verified | HIGH | `TestScenario_C20/20_6` reconciliation |
| 20.7 | JSON/BLOB expression DEFAULT (8.0.13+) | verified | HIGH | `TestScenario_C20` reconciliation |
| 20.8 | Generated col DEFAULT → grammar error | verified | MED | `TestScenario_C20/20_8` reconciliation |
| 21.1 | DEFAULT NULL literal | verified | HIGH | `C21_1_Default_NULL` — Wave 3 C21 worker updated |
| 21.2 | Bare JOIN → INNER (JTT_INNER) | verified | HIGH | `TestScenario_C21` reconciliation |
| 21.3 | ORDER BY no direction → ORDER_NOT_RELEVANT | verified | HIGH | `TestScenario_C21/21_3` reconciliation |
| 21.4 | LIMIT N no OFFSET → opt_offset NULL | verified | HIGH | `TestScenario_C21` reconciliation |
| 21.5 | FK clause omitted → FK_OPTION_UNDEF | verified | HIGH | `TestScenario_C21` reconciliation |
| 21.6 | FK asymmetric fill → UPDATE gets UNDEF | verified | MED | `TestScenario_C21` reconciliation |
| 21.7 | CREATE INDEX no USING → nullptr engine-default | verified | MED | `TestScenario_C21` reconciliation |
| 21.8 | INSERT no column list → empty PT_item_list | verified | HIGH | `TestScenario_C21` reconciliation |
| 21.9 | CREATE VIEW no ALGORITHM → UNDEFINED | verified | HIGH | `TestScenario_C21` reconciliation |
| 21.10 | CREATE TABLE no ENGINE → @@default_storage_engine | verified | HIGH | `TestScenario_C21/21_10` reconciliation |
| 22.1 | ALGORITHM=DEFAULT picks fastest supported | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 22.2 | LOCK=DEFAULT picks least restrictive | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 22.3 | ADD COLUMN nullable trailing is INSTANT | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 22.4 | DROP COLUMN INSTANT (8.0.29+) else INPLACE | verified | MED | `TestScenario_C22` catalog-level reconciliation |
| 22.5 | RENAME COLUMN INPLACE/INSTANT | verified | MED | `TestScenario_C22` catalog-level reconciliation |
| 22.6 | CHANGE COLUMN type change forces COPY | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 22.7 | Explicit ALGORITHM=INSTANT unsupported → error | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 22.8 | LOCK=NONE on COPY-only → error | verified | HIGH | `TestScenario_C22` catalog-level reconciliation |
| 24.1 | Invisible PK skip_gipk | verified | MED | `TestScenario_C24` |
| 25.1 | DECIMAL default (10,0) | verified | LOW | `C25_1_decimal_default_10_0` |
| PS.1 | CHECK counter CREATE fresh | verified | HIGH | `PS1_CheckCounter_CREATE_fresh` + bugfix test |
| PS.2 | CHECK counter ALTER max+1 | verified | HIGH | `PS1_CheckCounter_ALTER_open` |
| PS.3 | FK counter CREATE fresh | verified | HIGH | same as 1.1 |
| PS.4 | FK counter ALTER max+1 | verified | HIGH | same as 1.2 |
| PS.5 | DEFAULT NOW() fsp mismatch errors | verified | MED | `PS5_DatetimeFspMismatch` |
| PS.6 | HASH partition ADD count-seed | pending | LOW | N/A — no partition ALTER support |
| PS.7 | FK name collision error | verified | HIGH | `PS7_FKNameCollision` — **omni bug** |
| PS.8 | CHECK schema-scoped dup name error | pending | MED | **potential omni bug** |

| 3.4 | Explicit NULL on PK hard error | verified | HIGH | `TestScenario_C3/3_4` reconciliation |
| 3.5 | UNIQUE does NOT imply NOT NULL | verified | HIGH | `TestScenario_C3` reconciliation |
| 3.6 | Gcol nullability from expression | verified | MED | `TestScenario_C3` reconciliation |
| 3.7 | AUTO_INCREMENT + explicit NULL silent promote | verified | MED | `TestScenario_C3` reconciliation |
| 3.8 | 2nd TIMESTAMP implicit zero default (legacy) | pending | LOW | Wave 4 C3 worker |
| 5.4 | FK ON UPDATE independent of ON DELETE | verified | HIGH | `TestScenario_C5` reconciliation |
| 5.5 | FK SET DEFAULT rejected by InnoDB | verified | HIGH | `TestScenario_C5` reconciliation |
| 5.6 | FK column type must match (size/sign) | verified | HIGH | `TestScenario_C5` reconciliation |
| 5.7 | FK on VIRTUAL gcol rejected | verified | HIGH | `TestScenario_C5` reconciliation |
| 5.8 | CHECK NOT ENFORCED defaults ENFORCED | verified | HIGH | `TestScenario_C5` reconciliation |
| 5.9 | Column vs table CHECK equivalent | verified | MED | `TestScenario_C5` reconciliation |
| 5.10 | Column-CHECK cross-col reject | verified | MED | `TestScenario_C5` reconciliation |
| 8.4 | CHARSET defaults from database | verified | HIGH | `TestScenario_C8` reconciliation |
| 8.5 | COLLATE alone fills CHARSET | verified | HIGH | `TestScenario_C8` reconciliation |
| 8.6 | KEY_BLOCK_SIZE default 0 elided | verified | MED | `TestScenario_C8` reconciliation |
| 8.7 | COMPRESSION default None elided | verified | MED | `TestScenario_C8` + `TestScenario_C18/18_9` reconciliation |
| 8.8 | ENCRYPTION default from server var | verified | MED | `TestScenario_C8/8_8` reconciliation |
| 8.9 | STATS_PERSISTENT default | verified | LOW | `TestScenario_C8` + `TestScenario_C18/18_10` reconciliation |
| 8.10 | TABLESPACE default | verified | LOW | `TestScenario_C8` + `TestScenario_C18/18_12` reconciliation |
| 9.3 | VIRTUAL gcol PK rejection | verified | HIGH | `TestScenario_C9` reconciliation |
| 9.4 | Gcol expression must be deterministic | verified | HIGH | `TestScenario_C9` reconciliation |
| 9.5 | UNIQUE on STORED/VIRTUAL gcol | verified | MED | `TestScenario_C9` reconciliation |
| 9.6 | Gcol NOT NULL interaction | verified | MED | `TestScenario_C9` reconciliation |
| 9.7 | FK parent STORED gcol allowed | verified | MED | `TestScenario_C9` reconciliation |
| 9.8 | Gcol expression charset derivation | verified | LOW | `TestScenario_C9` reconciliation |
| 10.3 | View CHECK OPTION default CASCADED | verified | MED | `TestScenario_C10` reconciliation |
| 10.4 | ALGORITHM=UNDEFINED resolves MERGE/TEMPTABLE | verified | MED | `TestScenario_C10` reconciliation |
| 10.5 | View SQL SECURITY defaults DEFINER | verified | HIGH | `TestScenario_C10` reconciliation |
| 10.6 | View column names from SELECT | verified | HIGH | `TestScenario_C10` reconciliation |
| 10.7 | View updatability derivation | verified | MED | `TestScenario_C10` reconciliation |
| 10.8 | View column nullability widened | verified | MED | `TestScenario_C10` reconciliation |
| 11.2 | Trigger SQL SECURITY always DEFINER | verified | MED | `TestScenario_C11` reconciliation |
| 11.3 | Trigger charset/collation snapshot | verified | HIGH | `TestScenario_C11` reconciliation |
| 11.4 | Trigger ACTION_ORDER default | verified | HIGH | `TestScenario_C11` reconciliation |
| 11.5 | NEW/OLD pseudo-rows by event | verified | HIGH | `TestScenario_C11` reconciliation |
| 11.6 | Trigger on partitioned table | verified | LOW | `TestScenario_C11` reconciliation |
| 14.2 | ALTER CHECK ENFORCED toggle | verified | HIGH | `TestScenario_C14` reconciliation |
| 14.3 | STORED gcol + CHECK evaluation | verified | MED | `TestScenario_C14` reconciliation |
| 14.4 | CHECK forbidden constructs | verified | HIGH | `TestScenario_C14` reconciliation |
| 15.2 | ADD COLUMN FIRST | verified | HIGH | `TestScenario_C15` reconciliation |
| 15.3 | ADD COLUMN AFTER col | verified | HIGH | `TestScenario_C15` reconciliation |
| 15.4 | MODIFY/CHANGE retains position | verified | HIGH | `TestScenario_C15` reconciliation |
| 15.5 | Multiple ADD COLUMN left-to-right | verified | MED | `TestScenario_C15` reconciliation |
| 23.1 | CONCAT NULL propagation | verified | MED | `TestScenario_C23` reconciliation |
| 23.2 | CONCAT_WS skips NULL args | verified | MED | `TestScenario_C23` reconciliation |
| 23.3 | IFNULL/COALESCE rescue for CONCAT | verified | LOW | `TestScenario_C23` reconciliation |
| 24.2 | GIPK column spec | verified | MED | `TestScenario_C24` |
| 24.3 | GIPK suppressed by user PK | verified | MED | `TestScenario_C24` |
| 24.4 | GIPK my_row_id collision error | verified | LOW | `TestScenario_C24` |
| 25.2 | DECIMAL precision-only scale 0 | verified | MED | `TestScenario_C25` reconciliation |
| 25.3 | DECIMAL max precision/scale | verified | HIGH | `TestScenario_C25` reconciliation |
| 25.4 | UNSIGNED DECIMAL | verified | LOW | `TestScenario_C25` reconciliation |
| 25.5 | DECIMAL zero-scale rendering | verified | HIGH | `TestScenario_C25` reconciliation |
| AX.1 | ADD COLUMN append default | pending | HIGH | Wave 4 AX worker |
| AX.2 | DROP COLUMN cascades index | pending | HIGH | Wave 4 AX worker |
| AX.3 | DROP COLUMN last col error | pending | MED | Wave 4 AX worker |
| AX.4 | DROP COLUMN CHECK/gcol ref error | pending | HIGH | Wave 4 AX worker |
| AX.5 | MODIFY COLUMN rewrites spec | pending | HIGH | Wave 4 AX worker |
| AX.6 | CHANGE rename vs MODIFY | pending | HIGH | Wave 4 AX worker |
| AX.7 | ADD INDEX auto-name in ALTER | pending | HIGH | Wave 4 AX worker |
| AX.8 | ADD UNIQUE/KEY/FULLTEXT types | pending | MED | Wave 4 AX worker |
| AX.9 | ADD FOREIGN KEY forms | pending | HIGH | Wave 4 AX worker |
| AX.10 | RENAME COLUMN syntax | pending | HIGH | Wave 4 AX worker |
| AX.11 | RENAME INDEX syntax | pending | MED | Wave 4 AX worker |
| AX.12 | RENAME TO (table rename) | pending | MED | Wave 4 AX worker |
| AX.13 | ALTER COLUMN default/visibility | pending | MED | Wave 4 AX worker |
| AX.14 | ALTER INDEX VISIBLE/INVISIBLE | pending | MED | Wave 4 AX worker |
| AX.15 | Multi-subcommand single-pass | pending | HIGH | Wave 4 AX worker |
**Total scenarios:** 118 (81 pending, 35 verified via spot-check, 2 flagged as omni bugs).
