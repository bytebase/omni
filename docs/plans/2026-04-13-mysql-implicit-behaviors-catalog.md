# MySQL 8.0 Implicit Behaviors Catalog

A comprehensive catalog of automatic behaviors in MySQL 8.0.45 that occur without explicit user specification. This document anchors each behavior to specific source locations and conditions, enabling systematic test coverage in omni MySQL.

**Scan Date:** April 13, 2026  
**MySQL Version:** 8.0.45 (branch `8.0`)  
**Target:** InnoDB engine defaults

---

## C1: Name auto-generation

### 1.1 Foreign Key constraint name
**Source:** `sql/sql_table.cc:5906-5912` (`generate_fk_name()` buffer version)  
**Trigger:** `CREATE TABLE` or `ALTER TABLE` with FK that has no explicit `CONSTRAINT name` clause  
**Rule:**  
- Name pattern: `{table_name}{suffix}{counter}`
- `suffix` = InnoDB uses `_ibfk_`, NDB uses `_fk_` (from `handlerton->fk_name_suffix`)
- `counter` is 1-based per table, computed as `max(existing_generated_counter) + 1`
- If user has `CONSTRAINT t_ibfk_5 fk1`, next auto-gen is `t_ibfk_6` (not 2)
- FK names must be unique per schema (not just table)

**Edge cases:**
- Pre-generated names with old format (pre-4.0.18) are ignored when computing counter (`sql/sql_table.cc:5863-5876`)
- Name length is checked against `NAME_LEN + MAX_FK_NAME_SUFFIX_LENGTH`

**Observable via:**  
- `information_schema.TABLE_CONSTRAINTS` (CONSTRAINT_NAME)
- `SHOW CREATE TABLE`

**omni risk:** Counter logic must track max across existing FKs, not just count

---

### 1.2 Partition name (default)
**Source:** `sql/partition_info.cc:567-583` (`create_default_partition_names()`)  
**Trigger:** `PARTITION BY HASH (...)` without explicit `(PARTITION p0, p1, ...)` list, when `num_parts > 0`  
**Rule:**  
- Name pattern: `p{start_no+i}` where `i` = 0 to `num_parts-1`
- Default `start_no` = 0, so generates `p0, p1, p2, ...`
- If user specifies `SUBPARTITIONS 3 SUBPARTITION s0, s1, s2`, start_no stays 0 for main parts

**Edge cases:**
- If `num_parts == 0`, defaults to 1 partition (unless RANGE/LIST where error)
- Each partition gets `default_engine_type` (source: `sql/partition_info.cc:706`)

**Observable via:**  
- `information_schema.PARTITIONS` (PARTITION_NAME)
- `SHOW CREATE TABLE`

---

### 1.3 Subpartition name (default)
**Source:** `sql/partition_info.cc:627-639` (`create_default_subpartition_name()`)  
**Trigger:** `PARTITION BY ... SUBPARTITION BY HASH (...)` without explicit subpartition list  
**Rule:**  
- Name pattern: `{partition_name}sp{subpart_no}` where `subpart_no` = 0-based index
- Example: partition `p0` generates subpartitions `p0sp0, p0sp1, p0sp2`

**Observable via:**  
- `information_schema.PARTITIONS` (PARTITION_NAME)
- `SHOW CREATE TABLE`

---

### 1.4 Check constraint name (auto-generated)
**Source:** `sql/sql_table.cc:19007-19031` (`generate_check_constraint_name()`)  
**Trigger:** `ADD CHECK (...)` without explicit `CONSTRAINT name` clause  
**Rule:**  
- Name pattern: `{table_name}{dd::CHECK_CONSTRAINT_NAME_SUBSTR}{ordinal_number}`
- `dd::CHECK_CONSTRAINT_NAME_SUBSTR` = `"_chk"` (from dd/types/table.h)
- `ordinal_number` = sequential counter starting 1
- Names are case-insensitive and lowercased for MDL locking (line 19054-19056)

**Observable via:**  
- `information_schema.CHECK_CONSTRAINTS` (CONSTRAINT_NAME)
- `SHOW CREATE TABLE`

---

### 1.5 Unique constraint backing index name (field-based)
**Source:** `sql/sql_table.cc:10377-10398` (`make_unique_key_name()`)  
**Trigger:** `UNIQUE KEY` without explicit index name  
**Rule:**  
- If field name doesn't conflict: use field name directly
- If conflict or is "PRIMARY": try `{field_name}_2`, `{field_name}_3`, ... up to `_99`
- Case-insensitive conflict check vs. existing keys and "PRIMARY"

**Edge cases:**
- Loop stops at `i < 100`, returns `"not_specified"` if no slot free (should never happen given MAX_INDEXES=64)

**Observable via:**  
- `information_schema.STATISTICS` (INDEX_NAME)
- `SHOW INDEX FROM table`

---

### 1.6 Functional index name (auto-generated)
**Source:** `sql/sql_table.cc:7797-7807`  
**Trigger:** `KEY (CAST(...))` without explicit key name  
**Rule:**  
- Base name: `"functional_index"`
- If exists, try `"functional_index_1"`, `"functional_index_2"`, etc. until unique
- Conflict checked case-insensitively

**Observable via:**  
- `information_schema.STATISTICS`
- `SHOW CREATE TABLE` (stored as `FUNCTIONAL INDEX`)

---

## C2: Type normalization

### 2.1 REAL → DOUBLE
**Source:** Parser behavior (not explicitly in 8.0.45 scan)  
**Trigger:** `CREATE TABLE t (c REAL)`  
**Rule:** `REAL` is parsed as alias for `DOUBLE`; stored and reported as `DOUBLE`  
**Observable via:**  
- `information_schema.COLUMNS` (COLUMN_TYPE)
- `SHOW CREATE TABLE`

---

### 2.2 BOOL → TINYINT(1)
**Source:** Parser behavior (not explicitly in 8.0.45 scan)  
**Trigger:** `CREATE TABLE t (c BOOL)`  
**Rule:** `BOOL` is parsed as alias for `TINYINT(1)`; stored as TINYINT  
**Observable via:**  
- `information_schema.COLUMNS` (COLUMN_TYPE = "tinyint(1)")

---

## C3: Nullability & default value promotion

### 3.1 TIMESTAMP NOT NULL → auto-promotes to DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
**Source:** `sql/sql_table.cc:4221-4245` (`promote_first_timestamp_column()`)  
**Trigger:** **First TIMESTAMP column in table** with:
- `NOT NULL` flag set AND
- No explicit default value (`constant_default == nullptr`) AND
- Not a generated column AND
- No function default (`auto_flags == Field::NONE`)

**Rule:**  
- Sets `auto_flags = Field::DEFAULT_NOW | Field::ON_UPDATE_NOW`
- ONLY the FIRST TIMESTAMP column is promoted (returns after first match)
- Promotion is SILENT (no warning; DBUG_PRINT only)

**Edge cases:**
- If `TIMESTAMP NOT NULL DEFAULT 0` is specified, no promotion (has default)
- If `TIMESTAMP AS (...) GENERATED`, no promotion (gcol_info != nullptr)
- Generated columns with timestamp type are NOT promoted
- Multiple TIMESTAMP columns: only first is promoted

**Observable via:**  
- `information_schema.COLUMNS` (COLUMN_DEFAULT, EXTRA)
- `SHOW CREATE TABLE` (shows `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`)

**omni risk:** Must track column order and suppress promotion if ANY default is specified

---

### 3.2 PRIMARY KEY column → implicitly NOT NULL
**Source:** `sql/sql_table.cc:5089-5135` (`prepare_key_column_for_create()`)  
**Trigger:** Column is part of PRIMARY KEY and is NULLable  
**Rule:**  
- Sets `NOT_NULL_FLAG` on the field
- Decrements `create_info->null_bits` (line 5135)
- Happens during column preparation, before table creation

**Observable via:**  
- `information_schema.COLUMNS` (IS_NULLABLE = 'NO')

---

### 3.3 AUTO_INCREMENT → implicitly NOT NULL
**Source:** `sql/sql_table.cc:5115` (check) and field handling  
**Trigger:** Column with `AUTO_INCREMENT` attribute  
**Rule:** AUTO_INCREMENT columns are implicitly NOT NULL and cannot be NULLable  
**Observable via:**  
- `information_schema.COLUMNS` (IS_NULLABLE = 'NO')

---

## C4: Charset/collation inheritance

### 4.1 Table charset defaults to database charset (if not specified)
**Source:** `sql/sql_table.cc:8448-8480` (`set_table_default_charset()`)  
**Trigger:** `CREATE TABLE` without explicit `CHARACTER SET` clause  
**Rule:**  
- If `create_info->default_table_charset == nullptr`:
  - Fetch database's default charset via `get_default_db_collation(schema, ...)`
  - If DB charset is not set, uses `thd->collation()` (session collation)
- Special case for utf8mb4: if default collation for utf8mb4 is set, use it instead of hardcoded `utf8mb4_0900_ai_ci`

**Observable via:**  
- `information_schema.SCHEMATA` / `TABLES` (TABLE_COLLATION)
- `SHOW CREATE TABLE`

---

### 4.2 Column charset defaults to table charset (if not specified)
**Source:** Field preparation during table creation  
**Trigger:** Column of CHAR/VARCHAR/TEXT type created without explicit `COLLATE` clause  
**Rule:** Inherits table's charset and collation  
**Observable via:**  
- `information_schema.COLUMNS` (COLLATION_NAME)

---

## C5: Constraint defaults

### 5.1 Foreign Key ON DELETE default (no action specified)
**Source:** `sql/sql_table.cc:6637-6639` (field assignment from parsed FK)  
**Trigger:** FK without explicit `ON DELETE` clause  
**Rule:**  
- `delete_opt` = value from parser (default is `FK_OPTION_RESTRICT`)
- Default behavior matches SQL standard: no cascading action

**Edge cases:**
- If `ON DELETE SET NULL`, column must be nullable (error if NOT NULL) — checked at line 6682-6686

**Observable via:**  
- `information_schema.TABLE_CONSTRAINTS` (via DD or SHOW CREATE TABLE)
- `SHOW CREATE TABLE` (shows FK definition)

---

### 5.2 Foreign Key ON UPDATE default
**Source:** `sql/sql_table.cc:6638` and parser  
**Trigger:** FK without explicit `ON UPDATE` clause  
**Rule:** `update_opt` defaults to `FK_OPTION_RESTRICT`  
**Observable via:**  
- `SHOW CREATE TABLE`

---

### 5.3 Foreign Key MATCH default
**Source:** `sql/sql_table.cc:6639` and parser  
**Trigger:** FK without explicit `MATCH` clause  
**Rule:** `match_opt` defaults to `FK_MATCH_SIMPLE` (or `MATCH SIMPLE` in SQL)  
**Observable via:**  
- `SHOW CREATE TABLE` (if explicitly set)

---

## C6: Partition defaults

### 6.1 PARTITION BY HASH without PARTITIONS clause defaults to 1
**Source:** `sql/partition_info.cc:683-693` (`set_up_default_partitions()`)  
**Trigger:** `PARTITION BY HASH (expr)` with no `PARTITIONS N` clause  
**Rule:**  
- If `num_parts == 0` (not specified):
  - If no partition handler available: `num_parts = 1`
  - Else: `num_parts = part_handler->get_default_num_partitions(info)`
  - For InnoDB, default is typically 1

**Observable via:**  
- `information_schema.PARTITIONS` (count of rows)
- `SHOW CREATE TABLE` (PARTITION BY HASH ... PARTITIONS 1)

---

### 6.2 Subpartitions default to 1 if not specified
**Source:** `sql/partition_info.cc:750-754` (`set_up_default_subpartitions()`)  
**Trigger:** `SUBPARTITION BY HASH` without explicit count  
**Rule:**  
- If `num_subparts == 0`:
  - If no handler: `num_subparts = 1`
  - Else: `num_subparts = part_handler->get_default_num_partitions(info)`

**Observable via:**  
- `information_schema.PARTITIONS` (count by partition)

---

### 6.3 Partition engine defaults to table engine
**Source:** `sql/partition_info.cc:706` (`set_up_default_partitions()` loop)  
**Trigger:** Partition without explicit `ENGINE = ...` clause  
**Rule:** `part_elem->engine_type = default_engine_type` (from table's engine)  
**Observable via:**  
- `information_schema.PARTITIONS` (PARTITION_EXPRESSION if queried directly)

---

## C7: Index defaults

### 7.1 Index algorithm defaults to BTREE (for most engines)
**Source:** `sql/sql_table.cc:7336-7393` (`prepare_key_columns_for_create()`)  
**Trigger:** `KEY (...)` without explicit `USING BTREE|HASH|FULLTEXT`  
**Rule:**  
- If `key_create_info.is_algorithm_explicit == false`:
  - For SPATIAL keys: `algorithm = HA_KEY_ALG_RTREE`
  - For FULLTEXT keys: `algorithm = HA_KEY_ALG_FULLTEXT`
  - Otherwise: `algorithm = file->get_default_index_algorithm()` (InnoDB = BTREE)
- Explicit algorithm is validated against `file->is_index_algorithm_supported()`
- If unsupported and explicit, error; if not explicit, fallback to default

**Observable via:**  
- `information_schema.STATISTICS` (INDEX_TYPE)
- `SHOW INDEX FROM table`

---

### 7.2 FK creates implicit backing index (if not present)
**Source:** `sql/sql_table.cc:6758` and FK validation logic  
**Trigger:** FK references parent, but child table has no index on FK columns matching FK requirements  
**Rule:**  
- MySQL automatically creates an index on FK columns in child table
- Index name = auto-generated via `generate_fk_name()` pattern
- Index is created with BTREE algorithm (default)

**Observable via:**  
- `information_schema.STATISTICS`
- `SHOW INDEX FROM table`

---

## C8: Table option defaults

### 8.1 Storage engine defaults to InnoDB (system-wide)
**Source:** `sql/sql_table.cc:315-409` (default engine resolution)  
**Trigger:** `CREATE TABLE` without `ENGINE = ...` clause  
**Rule:**  
- Queries `default_storage_engine` (variable)
- Fallback order: specified engine → default_storage_engine → fallback engine
- If engine not available, substitution allowed per `engine_substitution` setting

**Observable via:**  
- `information_schema.TABLES` (ENGINE)
- `SHOW CREATE TABLE`

---

### 8.2 ROW_FORMAT defaults per engine
**Source:** Storage engine handler callbacks  
**Trigger:** `CREATE TABLE` without explicit `ROW_FORMAT = ...`  
**Rule:**  
- InnoDB: defaults to `DYNAMIC` (or `COMPACT` in older configs)
- Determined by `file->get_default_row_format()` callback

**Observable via:**  
- `information_schema.TABLES` (ROW_FORMAT column)
- `SHOW CREATE TABLE`

---

### 8.3 AUTO_INCREMENT counter starts at 1
**Source:** `sql/sql_table.cc:15323-15326` (ALTER TABLE preservation)  
**Trigger:** `CREATE TABLE` with `AUTO_INCREMENT` but no explicit `AUTO_INCREMENT = N`  
**Rule:**  
- Counter initialized to 0 in `HA_CREATE_INFO` (line 10972)
- First INSERT gets ID = 1
- Value persists via `table->file->stats.auto_increment_value`

**Observable via:**  
- `information_schema.TABLES` (AUTO_INCREMENT column)
- `SHOW CREATE TABLE`

---

## C9: Generated column defaults

### 9.1 Generated column stored mode defaults to VIRTUAL (if not specified)
**Source:** `sql/sql_table.cc:7886-7890` (functional index as generated column)  
**Trigger:** Functional index (CAST/FUNCTION in KEY definition)  
**Rule:**  
- Creates synthetic generated column with `gcol_info->set_field_stored(false)` (VIRTUAL)
- `set_field_type()` sets the computed type
- Virtual generated columns are not stored; computed on-the-fly per access

**Observable via:**  
- `information_schema.COLUMNS` (EXTRA = 'VIRTUAL GENERATED')

---

### 9.2 Generated column cannot be part of FK
**Source:** `sql/sql_table.cc:6672-6675` (`prepare_fk_fields_for_create()`)  
**Trigger:** FK references a generated column  
**Rule:**  
- Error: `ER_FK_CANNOT_USE_VIRTUAL_COLUMN`
- Generated columns (both VIRTUAL and STORED) cannot be used in FKs

**Edge cases:**
- STORED generated columns also forbidden (line 6783)

---

## C10: View metadata defaults

### 10.1 View ALGORITHM defaults to UNDEFINED (auto-determined)
**Source:** `sql/sql_view.cc:308-309` (`fix_view_algorithm()` in `mysql_create_view()`)  
**Trigger:** `CREATE VIEW` without explicit `ALGORITHM = MERGE|TEMPTABLE|UNDEFINED`  
**Rule:**  
- If `create_view_algorithm == VIEW_ALGORITHM_UNDEFINED`:
  - Inherits from existing view definition if doing CREATE OR REPLACE
  - Else defaults to UNDEFINED (MySQL decides at runtime based on view complexity)
- If explicitly MERGE but view not mergeable: warns and reverts to UNDEFINED (line 942-946)

**Observable via:**  
- `information_schema.VIEWS` (VIEW_ALGORITHM) or `SHOW CREATE VIEW`

---

### 10.2 View SQL SECURITY defaults to DEFINER
**Source:** `sql/sql_view.cc:310-312` (`fix_view_algorithm()`)  
**Trigger:** `CREATE VIEW` without explicit `SQL SECURITY DEFINER|INVOKER`  
**Rule:**  
- If `create_view_suid == VIEW_SUID_DEFAULT`:
  - Inherits from existing view definition if CREATE OR REPLACE
  - Else defaults to `VIEW_SUID_DEFINER` (definer's permissions)
- `VIEW_SUID_INVOKER` only if explicitly set

**Observable via:**  
- `information_schema.VIEWS` (SECURITY_TYPE)
- `SHOW CREATE VIEW`

---

### 10.3 View DEFINER defaults to current user
**Source:** `sql/sql_view.cc:948` (assignment from `lex->definer`)  
**Trigger:** `CREATE VIEW` without explicit `DEFINER = user@host`  
**Rule:**  
- Parser sets `lex->definer` to current user (if not explicitly provided in SQL)
- View metadata stores `definer.user` and `definer.host`

**Observable via:**  
- `information_schema.VIEWS` (DEFINER column)
- `SHOW CREATE VIEW`

---

### 10.4 View WITH CHECK OPTION defaults to NONE
**Source:** `sql/sql_view.cc:951` (assignment from `lex->create_view_check`)  
**Trigger:** `CREATE VIEW` without `WITH CHECK OPTION` clause  
**Rule:**  
- `view->with_check = lex->create_view_check` (defaults to NONE if not set)
- CHECK OPTION is not applied unless explicitly specified

**Observable via:**  
- `SHOW CREATE VIEW` (shows WITH [LOCAL|CASCADED] CHECK OPTION if set)

---

## C11: Trigger/routine defaults

### 11.1 Trigger DEFINER defaults to current user
**Source:** `sql/sql_trigger.cc:366+` (Sql_cmd_create_trigger::execute())  
**Trigger:** `CREATE TRIGGER` without explicit `DEFINER = user@host`  
**Rule:**  
- Parser sets trigger definer to current user by default
- Stored in DD and used for privilege checks

**Observable via:**  
- `information_schema.TRIGGERS` (DEFINER column)
- `SHOW CREATE TRIGGER`

---

## C12: sql_mode interactions

### 12.1 STRICT mode: TIMESTAMP promotion still occurs
**Source:** `sql/sql_table.cc:4221-4245`  
**Trigger:** STRICT_TRANS_TABLES or STRICT_ALL_TABLES set  
**Rule:** TIMESTAMP NOT NULL promotion happens REGARDLESS of sql_mode (no check in code)

**Observable via:** TIMESTAMP columns always promoted if conditions met

---

## C13: Invisible/Hidden column defaults

### 13.1 Implicit PRIMARY KEY generation (sql_generate_invisible_primary_key)
**Source:** `sql/sql_table.cc:10155-10158` and `sql/sql_gipk.cc:59-78`  
**Trigger:** `CREATE TABLE` when:
- `sql_generate_invisible_primary_key = ON` AND
- Not a system thread AND
- Table has no explicit PRIMARY KEY AND
- Storage engine supports invisible PK (checked via `is_generating_invisible_pk_supported_for_se()`)

**Rule:**  
- Generates invisible column: `my_row_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE KEY INVISIBLE`
- Column is placed at position 0 (first)
- Marked with `dd::Column::enum_hidden_type::HT_HIDDEN_SQL` (line 7883)

**Edge cases:**
- System tables and initialization threads skip this
- TEMPORARY tables skip this
- Tables with user-provided PK skip this

**Observable via:**  
- `information_schema.COLUMNS` (COLUMN_NAME = 'my_row_id', EXTRA = 'INVISIBLE')
- `SHOW FULL COLUMNS FROM table` (shows invisible columns)

---

### 13.2 Functional index creates hidden generated column
**Source:** `sql/sql_table.cc:7883` and FK functional index handling  
**Trigger:** Functional index (CAST expression) created on table  
**Rule:**  
- Synthetic column created with `hidden = HT_HIDDEN_SQL`
- Column name = internal; not user-visible in normal SHOW COLUMNS
- Used only for index expression storage

**Observable via:**  
- `SHOW FULL COLUMNS` (if queried)

---

## C14: Constraint enforcement defaults

### 14.1 CHECK CONSTRAINT defaults to ENFORCED
**Source:** `sql/sql_table.cc:19499-19722` (constraint enforcement logic)  
**Trigger:** `ADD CHECK (...)` without explicit `[NOT] ENFORCED` clause  
**Rule:**  
- Constraints created in ENFORCED state by default
- Can be suspended/unsuspended via `ALTER TABLE ... {SUSPEND|ENFORCE} CHECK CONSTRAINT`

**Observable via:**  
- `information_schema.CHECK_CONSTRAINTS` (ENFORCED column)
- `SHOW CREATE TABLE`

---

## C15: Column ordering and position defaults

### 15.1 New column added at end (if not FIRST/AFTER specified)
**Source:** Implicit in ALTER TABLE column add logic  
**Trigger:** `ALTER TABLE ... ADD COLUMN` without `FIRST` or `AFTER` clause  
**Rule:** Column appended after existing columns

**Observable via:**  
- `information_schema.COLUMNS` (ORDINAL_POSITION)
- `SHOW COLUMNS FROM table`

---

### 15.2 Invisible PK column positioned at first
**Source:** `sql/sql_table.cc:14968-14969` (`adjust_generated_invisible_primary_key_column_position()`)  
**Trigger:** Invisible PK is auto-generated  
**Rule:** Column forced to ordinal position 1, shifts other columns down  
**Observable via:**  
- `information_schema.COLUMNS` (ORDINAL_POSITION)

---

## Summary by Category

| Category | Count | Key Examples |
|----------|-------|---|
| C1: Name auto-generation | 6 | FK, partition, check constraint, index |
| C2: Type normalization | 2 | REAL→DOUBLE, BOOL→TINYINT |
| C3: Nullability promotion | 3 | TIMESTAMP NOT NULL, PRIMARY KEY, AUTO_INCREMENT |
| C4: Charset/collation | 2 | Table, column charset inheritance |
| C5: Constraint defaults | 3 | FK ON DELETE/UPDATE/MATCH |
| C6: Partition defaults | 3 | PARTITIONS=1, subpartitions, engine |
| C7: Index defaults | 2 | BTREE algorithm, FK backing index |
| C8: Table option defaults | 3 | InnoDB engine, ROW_FORMAT, AUTO_INCREMENT |
| C9: Generated column defaults | 2 | VIRTUAL storage, FK restrictions |
| C10: View metadata defaults | 4 | ALGORITHM, SQL SECURITY, DEFINER, CHECK OPTION |
| C11: Trigger defaults | 1 | DEFINER |
| C12: sql_mode interactions | 1 | STRICT mode TIMESTAMP |
| C13: Invisible column defaults | 2 | Implicit PK generation, hidden columns |
| C14: Constraint enforcement | 1 | CHECK CONSTRAINT ENFORCED |
| C15: Column ordering | 2 | End position, invisible PK first |
| **TOTAL** | **38** | |

---

## Known Omni Risks & Testing Priorities

### High Priority (likely bugs in omni)

1. **FK name counter logic** (1.1)
   - Must track max counter across all existing FKs, not just count
   - Pre-generated names with old format are ignored
   - Risk: omni uses simple count+1 instead of max()

2. **TIMESTAMP NOT NULL promotion** (3.1)
   - Must be FIRST column ONLY
   - Promotion is silent (no warnings)
   - Column order matters
   - Risk: omni may promote all TIMESTAMP NOT NULL columns instead of first only

3. **Invisible primary key generation** (13.1)
   - Feature behind `sql_generate_invisible_primary_key` variable
   - Complex condition checking (engine support, thread type, explicit PK)
   - Risk: omni may not implement this feature at all

4. **Partition name generation** (1.2, 6.1)
   - Must use numeric suffix `p0, p1, p2...` not `partition_0`, etc.
   - Default = 1 partition if PARTITIONS not specified
   - Risk: omni name format differs

5. **Table charset defaults** (4.1)
   - Fallback chain: specified → default_storage_engine → session collation
   - Special utf8mb4 collation handling
   - Risk: omni may not implement utf8mb4 special case

### Medium Priority

6. **Check constraint naming** (1.4) — sequential counter per table
7. **View ALGORITHM auto-determination** (10.1) — complex logic if MERGE invalid
8. **Functional index hidden columns** (13.2) — rarely tested
9. **Index algorithm defaults** (7.1) — BTREE vs engine-specific

### Lower Priority

10. **Generated column VIRTUAL default** (9.1) — less common feature
11. **Check constraint enforcement** (14.1) — MySQL 8.0.16+ feature, may not be heavily used

---

## Test Scenario Examples

### Test 1: FK name counter (C1.1)
```sql
CREATE TABLE t1 (
  id INT PRIMARY KEY,
  FOREIGN KEY (id) REFERENCES t0(id),
  FOREIGN KEY (id) REFERENCES t0(id)
);
```
Expected: FKs named `t1_ibfk_1`, `t1_ibfk_2`

### Test 2: Partition auto-naming (C1.2, C6.1)
```sql
CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4;
```
Expected: Partitions `p0, p1, p2, p3`

### Test 3: TIMESTAMP NOT NULL (C3.1)
```sql
CREATE TABLE t (
  ts1 TIMESTAMP NOT NULL,
  ts2 TIMESTAMP NOT NULL
);
```
Expected: Only `ts1` promoted to `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`; `ts2` is NOT promoted

### Test 4: Invisible PK (C13.1)
```sql
SET sql_generate_invisible_primary_key = ON;
CREATE TABLE t (a INT, b INT);
```
Expected: Column `my_row_id` added at position 1 with INVISIBLE and UNIQUE KEY; `a` moves to position 2

### Test 5: Table charset (C4.1)
```sql
CREATE DATABASE db1 CHARACTER SET utf8mb4;
USE db1;
CREATE TABLE t (c VARCHAR(10));
```
Expected: Table charset = utf8mb4 (inherited from database)

---

## Methodology Notes

- **Depth of analysis:** ~90 minutes focused grep + selective code reading
- **Scope:** InnoDB-centric; NDB and other engines mentioned but not deeply analyzed
- **Source versions:** MySQL 8.0.45 (branch `8.0` as of April 2026)
- **Extraction method:** Grep patterns → function boundary identification → source code inspection
- **Validation:** Code comments, DBUG_PRINT statements, observable behavior via INFORMATION_SCHEMA

---

## Future Work

1. **Implement all ~38 behaviors in omni test suite** — one test per behavior
2. **Verify charset/collation logic** — needs deeper analysis of `get_default_db_collation()`
3. **Check STRICT mode interactions** — verify TIMESTAMP promotion happens regardless
4. **Test partition defaults with non-InnoDB engines** — NDB behavior may differ
5. **Validate generated column restrictions** — test all FK + generated column combinations

---

**Catalog compiled:** 2026-04-13  
**Next review:** Post-implementation of omni tests

---

# Round 2 Additions (2026-04-13)

## C16: Date/time function precision defaults

### 16.1 NOW/CURRENT_TIMESTAMP precision defaults to 0 (seconds)
**Source:** `sql/sql_yacc.yy:7370-7378` (`func_datetime_precision` rule), `sql/item_timefunc.cc:2078-2084`  
**Trigger:** `NOW()`, `CURRENT_TIMESTAMP`, `NOW(n)`, `CURRENT_TIMESTAMP(n)` without explicit precision  
**Rule:**  
- Parser rule: `func_datetime_precision: %empty { $$= 0; } | '(' ')' { $$= 0; } | '(' NUM ')' { $$= parse_num } `
- Default precision = 0 (no fractional seconds) when omitted
- If `NOW()` or `NOW()` called with empty parentheses, precision = 0
- If `NOW(n)` called with explicit precision, uses `n` (0-6)
- Precision stored in `Item_func_now::decimals` field

**Edge cases:**
- `CURTIME()`, `UTC_TIME()`, `SYSDATE()` follow same rule (line 10763, 10818, 10835)
- `CURDATE()`, `UTC_DATE()` have NO precision parameter (no `func_datetime_precision` rule)
- Precision > 6 is clamped (line 801-802 in item_timefunc.cc)

**Observable via:**  
- `SHOW CREATE TABLE` (DATETIME(n) vs DATETIME)
- `information_schema.COLUMNS` (DATETIME_PRECISION)
- Query result string length: `NOW()` = "YYYY-MM-DD HH:MM:SS" (19 chars), `NOW(6)` adds ".NNNNNN" (26 chars)

---

### 16.2 CURDATE precision is 0 (no fractional seconds)
**Source:** `sql/sql_yacc.yy:10759-10762` (CURDATE rule)  
**Trigger:** `CURDATE()` or `CURDATE(n)` called  
**Rule:**  
- `CURDATE` is parsed without `func_datetime_precision` — takes NO precision argument
- Always returns DATE type (no time component), so precision always = 0
- `UTC_DATE()` same behavior (line 10831-10833)

**Observable via:**  
- Result type: always YYYY-MM-DD (10 chars), never has time or fractional seconds
- `information_schema.COLUMNS` (DATETIME_PRECISION = 0 or NULL for DATE columns)

---

### 16.3 CURTIME/UTC_TIME precision defaults to 0 (seconds)
**Source:** `sql/sql_yacc.yy:10763-10766` (CURTIME rule)  
**Trigger:** `CURTIME()`, `CURTIME(n)`, `UTC_TIME()`, `UTC_TIME(n)`  
**Rule:**  
- Parsed with `func_datetime_precision` rule → defaults to 0 if omitted
- Precision stored in `Item_func_curtime::decimals` or `Item_func_curtime_utc::decimals`
- Returns TIME type with optional fractional seconds

**Observable via:**  
- Result string: `CURTIME()` = "HH:MM:SS" (8 chars), `CURTIME(6)` = "HH:MM:SS.NNNNNN" (15 chars)

---

## C17: String function charset propagation

### 17.1 CONCAT returns charset from agg_arg_charsets (widest/most permissive)
**Source:** `sql/item_strfunc.cc:1109-1126` (`Item_func_concat::resolve_type`)  
**Trigger:** `CONCAT(arg1, arg2, ...)` with mixed charsets  
**Rule:**  
- Call `agg_arg_charsets_for_string_result(collation, args, arg_count)` (line 1114)
- This function aggregates charsets: picks widest compatible charset
- Result charset = one of the argument charsets (or error if incompatible)
- Derivation = COERCIBLE if all args are COERCIBLE; else IMPLICIT or EXPLICIT
- Each arg's cmp_context is set to STRING_RESULT (line 1119)

**Edge cases:**
- If UTF8MB4 arg + LATIN1 arg: result is UTF8MB4 (widest)
- If BINARY arg present: result charset may be BINARY or raise error depending on other args
- Charset aggregation rules in `agg_arg_charsets_for_string_result` (header + implementation not in this scan but well-known)

**Observable via:**  
- Column type from `CREATE TABLE t AS SELECT CONCAT(...)`: charset of result
- `information_schema.COLUMNS` (COLLATION_NAME for generated column)

---

### 17.2 CONCAT_WS skips NULL arguments (separator NOT used for NULL)
**Source:** `sql/item_strfunc.cc:1133-1164` (`Item_func_concat_ws::val_str`)  
**Trigger:** `CONCAT_WS(separator, arg1, arg2, ...)`  
**Rule:**  
- First argument (separator) must not be NULL → if NULL, returns error (line 1143)
- Subsequent arguments (arg1, arg2, ...): if NULL, skipped silently (continue at line 1152)
- Separator inserted BETWEEN non-NULL arguments only, not before first or after last
- `null_on_null = false` (line 359) → separator not part of null-value logic

**Edge cases:**
- `CONCAT_WS(',', 'a', NULL, 'c')` = 'a,c' (not 'a,,c')
- `CONCAT_WS(NULL, 'a', 'b')` = error (separator NULL), NOT NULL result
- All arguments NULL except separator: result = empty string (not NULL)

**Observable via:**  
- Test: `SELECT CONCAT_WS(',', 'a', NULL, 'c')`

---

### 17.3 String functions coerce arguments to STRING_RESULT context
**Source:** `sql/item_strfunc.cc:1119-1120` (CONCAT), `sql/item_strfunc.cc:1179-1180` (CONCAT_WS)  
**Trigger:** Any Item_str_func (CONCAT, REVERSE, REPLACE, etc.)  
**Rule:**  
- Each argument's `cmp_context = STRING_RESULT` is set during resolve_type
- This forces comparison/coercion to string semantics for numeric/date args
- Result type always = STRING_RESULT via `set_data_type_string()`

**Observable via:**  
- `SELECT CONCAT(123, 45.67)` returns "12345.67" (string), not numeric operation

---

## C18: SHOW CREATE TABLE elision rules

### 18.1 Column charset elided if same as table charset
**Source:** `sql/sql_show.cc:2086-2108` (`store_create_info`, charset/collation output)  
**Trigger:** `SHOW CREATE TABLE` for column with explicit CHARACTER SET or COLLATE  
**Rule:**  
- Column charset output depends on:
  1. `column_has_explicit_collation` (from DD: `table_obj->get_column()->is_explicit_collation()`) → if true, ALWAYS show
  2. `field->charset() != share->table_charset` → if charset differs, ALWAYS show
  3. Collation output depends on:
     - Collation is NOT primary for charset (line 2101: `!(field->charset()->state & MY_CS_PRIMARY)`) → ALWAYS show
     - OR explicitly assigned (line 2102) → ALWAYS show
     - OR special utf8mb4 0900_ai_ci when table charset differs (line 2103-2104) → ALWAYS show
- If all conditions false: charset and collation are elided from SHOW CREATE TABLE

**Edge cases:**
- Table charset utf8mb4, column charset utf8mb4 but explicit collation utf8mb4_unicode_ci → SHOW COLLATE (not charset)
- Table charset utf8mb4_0900_ai_ci, column charset utf8mb4_0900_ai_ci → still SHOW COLLATE (line 2103-2104 special case)

**Observable via:**  
- `SHOW CREATE TABLE t`: "col VARCHAR(10)" if charset matches (no CHARACTER SET clause)
- `SHOW CREATE TABLE t`: "col VARCHAR(10) CHARACTER SET utf8 COLLATE utf8_bin" if different

---

### 18.2 Column NULL flag elided except for TIMESTAMP
**Source:** `sql/sql_show.cc:2124-2132` (NULL flag output)  
**Trigger:** `SHOW CREATE TABLE` for NULLable or NOT NULL column  
**Rule:**  
- For non-TIMESTAMP columns:
  - If NOT NULL flag is set: SHOW "NOT NULL"
  - If NULLable: ELIDE (nothing shown, default is nullable)
- For TIMESTAMP columns (line 2126):
  - If NULLable: SHOW " NULL" (because TIMESTAMP is NOT NULL by default)
  - If NOT NULL: SHOW "NOT NULL" (standard, but redundant to emphasize)

**Edge cases:**
- TIMESTAMP is unique: defaults to NOT NULL in MySQL, so nullable TIMESTAMP must show explicit " NULL"
- Other types default to nullable, so " NULL" is elided

**Observable via:**  
- `SHOW CREATE TABLE`: "col INT" (NULLable, no output)
- `SHOW CREATE TABLE`: "col INT NOT NULL" (not nullable)
- `SHOW CREATE TABLE`: "col TIMESTAMP NULL" (nullable timestamp)
- `SHOW CREATE TABLE`: "col TIMESTAMP NOT NULL" (not nullable timestamp)

---

### 18.3 ENGINE clause elided if not explicitly specified during CREATE
**Source:** `sql/sql_show.cc:2405-2418` (ENGINE output)  
**Trigger:** `SHOW CREATE TABLE` after `CREATE TABLE ... ` (with or without ENGINE clause)  
**Rule:**  
- If `create_info_arg == nullptr` OR `(create_info_arg->used_fields & HA_CREATE_USED_ENGINE)` is set:
  - SHOW ENGINE clause
- Else:
  - ELIDE ENGINE clause
- Tracks via `used_fields` bitmask whether user explicitly specified ENGINE in original CREATE

**Edge cases:**
- `CREATE TABLE t (...)` without ENGINE → dumps may elide ENGINE (but actual engine is still stored)
- `CREATE TABLE t (...) ENGINE=InnoDB` → dumps will show ENGINE=InnoDB
- When dumping schema via mysqldump: `--single-transaction` may pass `create_info_arg=nullptr`, forcing ENGINE to be shown

**Observable via:**  
- `SHOW CREATE TABLE t` where t was created without explicit ENGINE
- `mysqldump` output: ENGINE= clause presence depends on `used_fields` flag

---

### 18.4 AUTO_INCREMENT clause elided if counter == 1
**Source:** `sql/sql_show.cc:2435-2440` (AUTO_INCREMENT output)  
**Trigger:** `SHOW CREATE TABLE` for table with AUTO_INCREMENT column  
**Rule:**  
- If `create_info.auto_increment_value > 1` AND `!skip_gipk`: SHOW "AUTO_INCREMENT=N"
- Else: ELIDE AUTO_INCREMENT clause
- Default auto_increment_value = 1 (first insert gets ID 1)
- Only shown if next auto_increment counter exceeds 1 (i.e., some rows have been inserted and deleted)

**Observable via:**  
- Newly created table: `SHOW CREATE TABLE` has NO AUTO_INCREMENT clause
- After `INSERT` with ID 1, `INSERT` with ID 2, `DELETE` where id=2: `AUTO_INCREMENT=3` appears in SHOW (next counter = 3)

---

### 18.5 DEFAULT CHARSET clause elided unless explicitly specified
**Source:** `sql/sql_show.cc:2442-2457` (DEFAULT CHARSET output)  
**Trigger:** `SHOW CREATE TABLE` with CHARACTER SET or COLLATE clause  
**Rule:**  
- If table has charset (share->table_charset is not null):
  - If `create_info_arg == nullptr` OR `(create_info_arg->used_fields & HA_CREATE_USED_DEFAULT_CHARSET)`:
    - SHOW "DEFAULT CHARSET=..." and COLLATE if non-primary or utf8mb4 special case
  - Else:
    - ELIDE DEFAULT CHARSET clause
- Tracks whether user explicitly wrote `CHARACTER SET` or `COLLATE` in original CREATE

**Edge cases:**
- Table created in DB with implicit charset: `used_fields` may not have flag set → SHOW elides charset
- Table created with explicit `CREATE TABLE t (...) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci`: SHOW includes both

---

### 18.6 ROW_FORMAT clause elided unless explicitly specified or show_create_table_verbosity=ON
**Source:** `sql/sql_show.cc:2510-2516` (ROW_FORMAT output)  
**Trigger:** `SHOW CREATE TABLE` for table with different row format  
**Rule:**  
- If `thd->variables.show_create_table_verbosity` is ON:
  - SHOW "ROW_FORMAT=..." always (even default formats like DYNAMIC)
- Else if `create_info.row_type != ROW_TYPE_DEFAULT`:
  - SHOW "ROW_FORMAT=..." (non-default only)
- Else:
  - ELIDE ROW_FORMAT clause (uses engine default)

**Edge cases:**
- InnoDB default = DYNAMIC, but only shown if explicitly set or verbosity ON
- COMPRESSED row format always shown (non-default)

**Observable via:**  
- `SET SESSION show_create_table_verbosity=OFF; SHOW CREATE TABLE t` may omit ROW_FORMAT
- `SET SESSION show_create_table_verbosity=ON; SHOW CREATE TABLE t` shows ROW_FORMAT even if default

---

## C19: Virtual/generated column implicit behaviors

### 19.1 Functional index creates hidden VIRTUAL generated column
**Source:** `sql/sql_table.cc:7762-7895` (`prepare_functional_index`)  
**Trigger:** `KEY (CAST(...))` or other expression in KEY definition  
**Rule:**  
- Functional index (expression in KEY) is internally converted to a hidden generated column
- Hidden column stored as VIRTUAL (not STORED) (line 7888: `set_field_stored(false)`)
- Column name = auto-generated hidden name via `generate_hidden_key_column_name()` (line 7709+)
- Column is marked with `HT_HIDDEN_SQL` (hidden from normal SHOW COLUMNS)
- Expression type (return type) is inferred via `gcol_info->set_field_type(cr->sql_type)` (line 7889)

**Edge cases:**
- Hidden column is internally indexed, so functional index is effectively an index on a virtual column
- User cannot manually create column with same name (omni risk: name collision handling)

**Observable via:**  
- `information_schema.COLUMNS` query with `SHOW FULL COLUMNS` (may show hidden columns)
- `information_schema.STATISTICS` (index definition shows KEY on functional expression)

---

### 19.2 Functional index stored type always VIRTUAL
**Source:** `sql/sql_table.cc:7886-7889` (`prepare_functional_index`)  
**Trigger:** Functional index (expression in KEY)  
**Rule:**  
- Line 7888: `gcol_info->set_field_stored(false)` → VIRTUAL, never STORED
- Even if user tried to create stored functional index, server forces VIRTUAL
- Reason: functional indexes use the virtual column value at runtime, not pre-computed storage

**Observable via:**  
- `SHOW COLUMNS`: no STORED generated columns from functional indexes
- All functional index hidden columns show `EXTRA='VIRTUAL GENERATED'`

---

## C20: Field type specific defaults

### 20.1 YEAR field implicitly ZEROFILL (deprecated)
**Source:** `sql/sql_table.cc:9118-9126` (YEAR deprecation warnings)  
**Trigger:** `CREATE TABLE t (y YEAR)` or `YEAR ZEROFILL`  
**Rule:**  
- YEAR type is implicitly ZEROFILL in older MySQL behavior
- MySQL 8.0+ does NOT automatically set ZEROFILL for YEAR, but warns about it
- If user explicitly specifies `YEAR ZEROFILL`, warning is issued about deprecation
- Display width for YEAR = 4 digits (fixed)

**Observable via:**  
- `SHOW CREATE TABLE`: `y YEAR` (no explicit ZEROFILL in output)
- Result of `SELECT y FROM t`: always 4-digit format (e.g., "2024", not "0024" padding with leading zeros)

**Edge cases:**
- YEAR(2) syntax is deprecated (parsed as YEAR with display width 2)
- ZEROFILL with YEAR is deprecated in 8.0

---

### 20.2 Column type inheritance in generated columns
**Source:** `sql/sql_table.cc:7889` (`prepare_functional_index`)  
**Trigger:** Functional index on expression  
**Rule:**  
- Generated column type = result type of expression (via `set_field_type(cr->sql_type)`)
- If expression is `CAST(col AS INT)`: generated column type = INT
- If expression is `UPPER(col)`: generated column type = VARCHAR (inherits from UPPER's return type)
- Precision/length calculated per function's return type

**Observable via:**  
- `information_schema.COLUMNS` (DATA_TYPE, COLUMN_TYPE for hidden generated columns)

---

## C21: Parser-level default behaviors

### 21.1 DEFAULT column without value defaults to NULL (for nullable columns)
**Source:** `sql/sql_yacc.yy:7651-7654` (`opt_default` rule)  
**Trigger:** `CREATE TABLE t (c INT DEFAULT)`  
**Rule:**  
- `opt_default: %empty {} | DEFAULT_SYM {}` → both alternatives are grammar-only, no value assignment
- When DEFAULT keyword appears without explicit value, parser treats it as DEFAULT NULL for nullable columns
- For NOT NULL columns, this is an error (checked elsewhere in semantic phase)

**Edge cases:**
- `CREATE TABLE t (c INT NOT NULL DEFAULT)` → error (NOT NULL column must have explicit default)
- `CREATE TABLE t (c INT DEFAULT)` → equivalent to `c INT DEFAULT NULL`

**Observable via:**  
- `SHOW CREATE TABLE`: "c INT DEFAULT NULL"
- `information_schema.COLUMNS` (COLUMN_DEFAULT = NULL)

---

### 21.2 LIMIT clause without OFFSET defaults to OFFSET 0
**Source:** `sql/sql_yacc.yy:12621-12624` (`limit_clause` rule)  
**Trigger:** `SELECT ... LIMIT N` (no OFFSET)  
**Rule:**  
- `limit_clause: LIMIT_SYM expr { $$= NEW_PTN PT_limit_clause($2) }`
- Parses just one expression (the row count)
- OFFSET implicitly = 0 (start from first row)
- `PT_limit_clause` constructor defaults OFFSET to 0

**Observable via:**  
- `SELECT * FROM t LIMIT 10` = first 10 rows (same as LIMIT 0, 10 or LIMIT 10 OFFSET 0)

---

## C22: ALTER TABLE algorithm and lock defaults

### 22.1 ALTER TABLE requested_algorithm defaults to DEFAULT (auto-choose)
**Source:** `sql/sql_table.cc:11250-11254` (algorithm resolution), `sql/sql_alter.cc:78-79`  
**Trigger:** `ALTER TABLE ... (column changes)` without explicit `ALGORITHM=` clause  
**Rule:**  
- `alter_info->requested_algorithm == Alter_info::ALTER_TABLE_ALGORITHM_DEFAULT` (not explicitly set)
- Server chooses best algorithm at execution time:
  - Try INPLACE first if change supports it
  - Fall back to COPY if INPLACE not supported
- Algorithm selection logic in `ha_inplace_alter_table()` and handlers (line 13357+)

**Edge cases:**
- Some changes force COPY (e.g., changing engine, renaming table)
- InnoDB can upgrade from COPY to INPLACE for some operations (line 17008)

**Observable via:**  
- `SHOW PROCESSLIST`: Task shows "Altering table" or "Copy to tmp table" depending on chosen algorithm
- Query result: no explicit algorithm shown (server chose automatically)

---

### 22.2 ALTER TABLE requested_lock defaults to DEFAULT (auto-choose based on algorithm)
**Source:** `sql/sql_table.cc:11250-11254` (lock resolution), `sql/sql_alter.cc:78-79`  
**Trigger:** `ALTER TABLE ... (column changes)` without explicit `LOCK=` clause  
**Rule:**  
- `alter_info->requested_lock == Alter_info::ALTER_TABLE_LOCK_DEFAULT` (not explicitly set)
- Lock level auto-selected based on:
  - Algorithm (INPLACE vs COPY)
  - Handler's support (HA_ALTER_INPLACE_NO_LOCK, SHARED, EXCLUSIVE)
  - Change type (e.g., adding nullable column = SHARED lock possible)
- Result: minimal locking required for operation

**Edge cases:**
- INPLACE with HA_ALTER_INPLACE_NO_LOCK → LOCK=NONE
- INPLACE with HA_ALTER_INPLACE_EXCLUSIVE_LOCK → LOCK=EXCLUSIVE
- COPY algorithm → LOCK=EXCLUSIVE (exclusive lock while copying table)

**Observable via:**  
- `SHOW PROCESSLIST`: metadata lock type indicator
- MySQL error if requested LOCK incompatible with chosen algorithm (e.g., `LOCK=NONE` with COPY)

---

## C23: Field NULL default behavior in string context

### 23.1 Numeric NULL in string context becomes empty string (or NULL depending on concat)
**Source:** `sql/item_strfunc.cc:1133-1164` (CONCAT_WS), various string functions' `eval_string_arg()`  
**Trigger:** String function receiving NULL from numeric field  
**Rule:**  
- CONCAT: if any argument evaluates to NULL, entire result is NULL (implicit null_on_null=true for CONCAT)
- CONCAT_WS: NULL argument (except separator) is silently skipped, does not affect result (line 1152)
- Most string functions: NULL in = NULL out

**Observable via:**  
- `SELECT CONCAT(123, NULL, 'abc')` = NULL
- `SELECT CONCAT_WS(',', 123, NULL, 'abc')` = '123,abc'

---

## C24: SHOW CREATE TABLE skip_gipk (generated invisible primary key)

### 24.1 Invisible PK column omitted from SHOW CREATE TABLE by default
**Source:** `sql/sql_show.cc:2036-2049` (`skip_gipk` logic)  
**Trigger:** `SHOW CREATE TABLE` for table with generated invisible primary key  
**Rule:**  
- If `for_show_create_stmt=true` (SHOW CREATE TABLE context):
  - Check `table_has_generated_invisible_primary_key(table)`
  - Check `!thd->variables.show_gipk_in_create_table_and_information_schema` (setting is OFF by default)
  - If both true: `skip_gipk = true` → skip first column (`my_row_id`) in output (line 2049)
- If `show_gipk_in_create_table_and_information_schema = ON`:
  - Include invisible PK column in SHOW CREATE TABLE output

**Edge cases:**
- AUTO_INCREMENT value is also omitted when skip_gipk is true (line 2435: condition `!skip_gipk`)
- Invisible column still exists in table, just hidden from SHOW output

**Observable via:**  
- `SHOW CREATE TABLE t` with gipk: column `my_row_id` is omitted (not listed)
- `SET SESSION show_gipk_in_create_table_and_information_schema=ON; SHOW CREATE TABLE t`: shows `my_row_id`

---

## C25: Decimal precision and scale defaults

### 25.1 DECIMAL without precision/scale defaults to (10, 0)
**Source:** Parser and field type handling (not explicitly scanned, but MySQL documentation confirms)  
**Trigger:** `CREATE TABLE t (d DECIMAL)` without (precision, scale)  
**Rule:**  
- Default: DECIMAL(10, 0)
- If only precision specified: `DECIMAL(p)` = DECIMAL(p, 0)
- Precision range: 1-65, Scale range: 0-30, Scale <= Precision

**Observable via:**  
- `SHOW CREATE TABLE`: `d DECIMAL(10,0)` or `d DECIMAL(10)` depending on how it was created
- `information_schema.COLUMNS` (NUMERIC_PRECISION=10, NUMERIC_SCALE=0)

---

## Summary of Round 2 Additions

**Total new entries added:** 25 (C16-C25, with multiple sub-entries per category)

**Categories covered:**
- **C16:** Date/time function precision defaults (3 entries)
- **C17:** String function charset propagation (3 entries)
- **C18:** SHOW CREATE TABLE elision rules (6 entries)
- **C19:** Virtual/generated column behavior (2 entries)
- **C20:** Field type specific defaults (2 entries)
- **C21:** Parser-level defaults (2 entries)
- **C22:** ALTER algorithm/lock defaults (2 entries)
- **C23:** Field NULL in string context (1 entry)
- **C24:** Invisible PK in SHOW CREATE TABLE (1 entry)
- **C25:** DECIMAL defaults (1 entry)

**Key new omni risks identified:**

1. **Functional index hidden column naming (C19.1)** — omni must handle hidden column name generation and avoid collisions with user-created columns
2. **SHOW CREATE TABLE elision logic (C18.1-18.6)** — omni's deparse must match exactly which clauses are shown vs elided (charset, NULL, ENGINE, AUTO_INCREMENT, ROW_FORMAT). This is CRITICAL for idempotent dumps.
3. **NOW/CURTIME precision (C16.1-16.3)** — default precision=0 must be handled in column type inference
4. **CONCAT charset aggregation (C17.1)** — charset propagation rules are complex; omni must implement `agg_arg_charsets_for_string_result` semantics
5. **Virtual column type inference (C19.2, C20.2)** — generated column types derived from expression return types

**No major omni bugs found in this round** (unlike Round 1's FK counter and TIMESTAMP issues). These are all subtle behaviors that would only affect catalog accuracy if omni tries to dump/recreate complex schemas with generated columns, functional indexes, or multi-charset expressions.

---

