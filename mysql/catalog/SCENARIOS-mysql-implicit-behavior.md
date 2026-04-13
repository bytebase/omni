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
| 2.3 | INTEGER → INT | pending | LOW | Wave 1 C2 worker |
| 2.4 | BOOLEAN → TINYINT(1) | pending | LOW | Wave 1 C2 worker |
| 2.5 | INT1/INT2/INT3/INT4/INT8 aliases | pending | LOW | Wave 1 C2 worker |
| 2.6 | MIDDLEINT → MEDIUMINT | pending | LOW | Wave 1 C2 worker |
| 2.7 | INT(11) display width stripped | pending | HIGH | Wave 1 C2 worker |
| 2.8 | INT(N) ZEROFILL preserves width | pending | HIGH | Wave 1 C2 worker |
| 2.9 | SERIAL composite alias | pending | HIGH | Wave 1 C2 worker |
| 2.10 | NUMERIC → DECIMAL | pending | LOW | Wave 1 C2 worker |
| 2.11 | DEC / FIXED → DECIMAL | pending | LOW | Wave 1 C2 worker |
| 2.12 | DOUBLE PRECISION → DOUBLE | pending | LOW | Wave 1 C2 worker |
| 2.13 | FLOAT4 → FLOAT, FLOAT8 → DOUBLE | pending | LOW | Wave 1 C2 worker |
| 2.14 | FLOAT(p) precision split | pending-verify | HIGH | Wave 1 C2 worker |
| 2.15 | FLOAT(M,D) deprecated non-standard | pending-verify | MED | Wave 1 C2 worker |
| 2.16 | CHARACTER / CHARACTER VARYING | pending | LOW | Wave 1 C2 worker |
| 2.17 | NATIONAL CHAR / NCHAR utf8mb3 | pending-verify | MED | Wave 1 C2 worker |
| 2.18 | NVARCHAR / NCHAR VARYING utf8mb3 | pending-verify | MED | Wave 1 C2 worker |
| 2.19 | LONG / LONG VARCHAR → MEDIUMTEXT | pending | LOW | Wave 1 C2 worker |
| 2.20 | CHAR / BINARY default length 1 | pending | MED | Wave 1 C2 worker |
| 2.21 | VARCHAR without length syntax error | pending-verify | MED | Wave 1 C2 worker |
| 2.22 | TIMESTAMP/DATETIME/TIME default fsp 0 | pending | MED | Wave 1 C2 worker |
| 2.23 | YEAR(4) deprecated → bare YEAR | pending | MED | Wave 1 C2 worker |
| 2.24 | BIT without length defaults to BIT(1) | pending-verify | HIGH | Wave 1 C2 worker |
| 2.25 | VARCHAR > 65535 → TEXT family | pending-verify | MED | Wave 1 C2 worker |
| 2.26 | TEXT(N)/BLOB(N) byte-count promotion | pending-verify | HIGH | Wave 1 C2 worker |
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
| 6.4 | KEY partitioning ALGORITHM defaults to 2 | pending | HIGH | Wave 1 C6 worker |
| 6.5 | KEY() empty list defaults to PRIMARY KEY cols | pending | HIGH | Wave 1 C6 worker |
| 6.6 | LINEAR HASH / LINEAR KEY preserved | pending | LOW | Wave 1 C6 worker |
| 6.7 | RANGE/LIST require explicit definitions | pending | HIGH | Wave 1 C6 worker |
| 6.8 | MAXVALUE only in last RANGE partition | pending | MED | Wave 1 C6 worker |
| 6.9 | LIST equality vs RANGE strict less-than | pending | LOW | Wave 1 C6 worker |
| 6.10 | LIST DEFAULT partition catch-all | pending | MED | Wave 1 C6 worker |
| 6.11 | Partition function must return INTEGER | pending | HIGH | Wave 1 C6 worker |
| 6.12 | TIMESTAMP requires UNIX_TIMESTAMP() wrap | pending | MED | Wave 1 C6 worker |
| 6.13 | UNIQUE/PK must contain partition cols | pending | HIGH | Wave 1 C6 worker |
| 6.14 | Per-partition options not inherited | pending | MED | Wave 1 C6 worker |
| 6.15 | Subpartition options inherit from parent | pending | LOW | Wave 1 C6 worker |
| 6.16 | ADD PARTITION auto-naming seeds from count | pending | LOW | Wave 1 C6 worker |
| 6.17 | COALESCE/REORGANIZE partition behavior | pending | LOW | Wave 1 C6 worker |
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
| 16.2 | NOW(N) explicit precision 0..6 | pending | HIGH | Wave 1 C16 worker |
| 16.3 | CURDATE / UTC_DATE take no precision | pending | MED | Wave 1 C16 worker |
| 16.4 | CURTIME / UTC_TIME precision defaults to 0 | pending | MED | Wave 1 C16 worker |
| 16.5 | SYSDATE precision defaults to 0 | pending | LOW | Wave 1 C16 worker |
| 16.6 | UTC_TIMESTAMP precision defaults to 0 | pending | LOW | Wave 1 C16 worker |
| 16.7 | UNIX_TIMESTAMP return type by arg fsp | pending | MED | Wave 1 C16 worker |
| 16.8 | DATETIME(N) DEFAULT NOW() fsp must match | pending | HIGH | Wave 1 C16 worker |
| 16.9 | ON UPDATE NOW(N) must match column fsp | pending | HIGH | Wave 1 C16 worker |
| 16.10 | DATETIME/TIMESTAMP storage bytes by fsp | pending | MED | Wave 1 C16 worker |
| 16.11 | YEAR(N) deprecated — only YEAR(4) accepted | pending | LOW | Wave 1 C16 worker |
| 16.12 | TIMESTAMP promotion inherits fsp | pending | MED | Wave 1 C16 worker |
| 18.1 | Column charset elision | verified | HIGH | `C4_2_and_C18_1_and_C18_5_charset_inheritance_and_elision` |
| 18.2 | NOT NULL elision (TIMESTAMP) | verified | HIGH | `C18_2_NotNull_rendering` |
| 18.3 | ENGINE always rendered | pending | HIGH | |
| 18.4 | AUTO_INCREMENT elided at 1 | verified | HIGH | `C18_4_auto_increment_elision` |
| 18.5 | DEFAULT CHARSET always rendered | verified | HIGH | `C18_5_DefaultCharset_implicit` |
| 18.6 | ROW_FORMAT elided unless explicit | pending | HIGH | |
| 21.1 | DEFAULT NULL literal | verified | LOW | `C21_1_Default_NULL` |
| 22.1 | ALGORITHM=DEFAULT picks fastest supported | pending | HIGH | Wave 1 C22 worker |
| 22.2 | LOCK=DEFAULT picks least restrictive | pending | HIGH | Wave 1 C22 worker |
| 22.3 | ADD COLUMN nullable trailing is INSTANT | pending | HIGH | Wave 1 C22 worker |
| 22.4 | DROP COLUMN INSTANT (8.0.29+) else INPLACE | pending | MED | Wave 1 C22 worker |
| 22.5 | RENAME COLUMN INPLACE/INSTANT | pending | MED | Wave 1 C22 worker |
| 22.6 | CHANGE COLUMN type change forces COPY | pending | HIGH | Wave 1 C22 worker |
| 22.7 | Explicit ALGORITHM=INSTANT unsupported → error | pending | HIGH | Wave 1 C22 worker |
| 22.8 | LOCK=NONE on COPY-only → error | pending | HIGH | Wave 1 C22 worker |
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
