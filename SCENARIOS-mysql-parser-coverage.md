# MySQL Parser Syntax Coverage Scenarios

> Goal: Close all syntax-bearing gaps between omni's MySQL parser and MySQL 8.0, with oracle-verified correctness
> Verification: `go test ./mysql/parser/...` — each scenario has a parse test (valid SQL parses, correct AST) or error test (invalid SQL rejected)
> Reference sources: MySQL 8.0 docs (https://dev.mysql.com/doc/refman/8.0/en/), real MySQL 8.0 via testcontainers

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Keyword & Type Foundation

### 1.1 Missing Type Synonyms — Numeric

These are MySQL 8.0 numeric type aliases that the parser should accept and map to canonical types.

- [x] `CREATE TABLE t (a REAL)` — parses, maps to DOUBLE
- [x] `CREATE TABLE t (a REAL(10,2))` — parses with precision, maps to DOUBLE(10,2)
- [x] `CREATE TABLE t (a REAL UNSIGNED)` — parses with UNSIGNED modifier
- [x] `CREATE TABLE t (a DEC)` — parses, maps to DECIMAL
- [x] `CREATE TABLE t (a DEC(10,2))` — parses with precision, maps to DECIMAL(10,2)
- [x] `CREATE TABLE t (a DEC UNSIGNED)` — parses with UNSIGNED modifier
- [x] `CREATE TABLE t (a FIXED)` — parses, maps to DECIMAL
- [x] `CREATE TABLE t (a FIXED(10,2))` — parses with precision, maps to DECIMAL(10,2)
- [x] `CREATE TABLE t (a FIXED UNSIGNED ZEROFILL)` — parses with all modifiers
- [x] `CREATE TABLE t (a INT1)` — parses, maps to TINYINT
- [x] `CREATE TABLE t (a INT2)` — parses, maps to SMALLINT
- [x] `CREATE TABLE t (a INT3)` — parses, maps to MEDIUMINT
- [x] `CREATE TABLE t (a INT4)` — parses, maps to INT
- [x] `CREATE TABLE t (a INT8)` — parses, maps to BIGINT
- [x] `CREATE TABLE t (a MIDDLEINT)` — parses, maps to MEDIUMINT
- [x] `CREATE TABLE t (a FLOAT4)` — parses, maps to FLOAT
- [x] `CREATE TABLE t (a FLOAT8)` — parses, maps to DOUBLE
- [x] `CREATE TABLE t (a DOUBLE PRECISION)` — existing synonym still works (no regression)

### 1.2 Missing Type Synonyms — String & Binary

- [x] `CREATE TABLE t (a LONG)` — parses, maps to MEDIUMTEXT
- [x] `CREATE TABLE t (a LONG VARCHAR)` — parses, maps to MEDIUMTEXT
- [x] `CREATE TABLE t (a LONG VARBINARY)` — parses, maps to MEDIUMBLOB
- [x] `CREATE TABLE t (a LONG CHARACTER SET utf8mb4)` — LONG with charset
- [x] `CREATE TABLE t (a LONG VARCHAR COLLATE utf8mb4_general_ci)` — LONG VARCHAR with collation
- [x] `CREATE TABLE t (a TEXT(1000))` — TEXT with optional length (existing, regression check)
- [x] `CREATE TABLE t (a NATIONAL CHAR(10))` — existing synonym (regression check)
- [x] `CREATE TABLE t (a NCHAR(10))` — existing synonym (regression check)
- [x] `CREATE TABLE t (a NVARCHAR(100))` — existing synonym (regression check)

### 1.3 SIGNED Modifier

- [x] `CREATE TABLE t (a INT SIGNED)` — parses, SIGNED accepted as modifier
- [x] `CREATE TABLE t (a BIGINT SIGNED)` — parses on BIGINT
- [x] `CREATE TABLE t (a TINYINT SIGNED)` — parses on TINYINT
- [x] `CREATE TABLE t (a SMALLINT SIGNED)` — parses on SMALLINT
- [x] `CREATE TABLE t (a MEDIUMINT SIGNED)` — parses on MEDIUMINT
- [x] `CREATE TABLE t (a TINYINT SIGNED NOT NULL)` — SIGNED followed by column constraints
- [x] `CREATE TABLE t (a DECIMAL SIGNED)` — oracle-verify: check if MySQL 8.0 accepts or rejects this
- [x] `SELECT CAST(x AS SIGNED)` — still works (no regression)
- [x] `SELECT CAST(x AS SIGNED INTEGER)` — still works (no regression)

### 1.4 Reserved Word Registration

Keywords that MySQL 8.0 reserves and that affect parsing when used in SQL.

- [x] `SELECT 1 FROM DUAL` — DUAL as keyword in FROM clause
- [x] `PARTITION p1 VALUES LESS THAN (MAXVALUE)` in CREATE TABLE — MAXVALUE recognized
- [x] `SELECT GROUPING(a) FROM t GROUP BY a WITH ROLLUP` — GROUPING function
- [x] `SELECT VALUE FROM t` — VALUE as reserved word (verify behavior with MySQL 8.0)

### 1.5 Numeric Type Shorthands in Combinations

- [x] `CREATE TABLE t (a INT1 UNSIGNED ZEROFILL)` — full modifier chain on shorthand
- [x] `CREATE TABLE t (a INT4(11) UNSIGNED)` — display width on shorthand
- [x] `ALTER TABLE t ADD COLUMN b REAL DEFAULT 0.0` — type synonym in ALTER TABLE context
- [x] `ALTER TABLE t MODIFY c DEC(10,2) NOT NULL` — type synonym in MODIFY COLUMN
- [x] `CREATE TABLE t (a FLOAT4, b FLOAT8, c INT1, d INT2, e INT3, f INT4, g INT8)` — all shorthands in one table
- [x] `CREATE TABLE t (a SERIAL)` — SERIAL type (existing, regression check)

---

## Phase 2: Statement Syntax Gaps

### 2.1 ALTER TABLE PARTITION BY

- [x] `ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4` — basic hash repartition
- [x] `ALTER TABLE t PARTITION BY KEY(id) PARTITIONS 4` — key repartition
- [x] `ALTER TABLE t PARTITION BY RANGE(id) (PARTITION p0 VALUES LESS THAN (100), PARTITION p1 VALUES LESS THAN (MAXVALUE))` — range repartition
- [x] `ALTER TABLE t PARTITION BY LIST(status) (PARTITION p0 VALUES IN (1,2), PARTITION p1 VALUES IN (3,4))` — list repartition
- [x] `ALTER TABLE t PARTITION BY RANGE COLUMNS(created_at) (PARTITION p0 VALUES LESS THAN ('2024-01-01'))` — range columns
- [x] `ALTER TABLE t PARTITION BY LINEAR HASH(id) PARTITIONS 8` — linear hash
- [x] `ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4 SUBPARTITION BY KEY(name) SUBPARTITIONS 2` — with subpartition
- [x] `ALTER TABLE t PARTITION BY KEY ALGORITHM=2 (id) PARTITIONS 4` — KEY with algorithm
- [x] `ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, ALGORITHM=INPLACE` — combined with ALTER TABLE options
- [x] `ALTER TABLE t PARTITION BY` — rejected: missing partition specification
- [x] `ALTER TABLE t REMOVE PARTITIONING` — existing operation (regression check)

### 2.2 Window Function Names as Reserved Words

MySQL 8.0 reserves these as keywords. Parser must accept them as function names in window expressions.

- [x] `SELECT RANK() OVER (ORDER BY score DESC) FROM t` — RANK
- [x] `SELECT DENSE_RANK() OVER (ORDER BY score DESC) FROM t` — DENSE_RANK
- [x] `SELECT ROW_NUMBER() OVER (ORDER BY id) FROM t` — ROW_NUMBER
- [x] `SELECT NTILE(4) OVER (ORDER BY id) FROM t` — NTILE with argument
- [x] `SELECT LAG(val, 1) OVER (ORDER BY id) FROM t` — LAG with offset
- [x] `SELECT LEAD(val, 1, 0) OVER (ORDER BY id) FROM t` — LEAD with offset and default
- [x] `SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t` — FIRST_VALUE
- [x] `SELECT LAST_VALUE(val) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) FROM t` — LAST_VALUE with frame
- [x] `SELECT NTH_VALUE(val, 2) OVER (ORDER BY id) FROM t` — NTH_VALUE
- [x] `SELECT PERCENT_RANK() OVER (ORDER BY score) FROM t` — PERCENT_RANK
- [x] `SELECT CUME_DIST() OVER (ORDER BY score) FROM t` — CUME_DIST

### 2.3 Interval Unit Keywords

MySQL reserves interval units for `INTERVAL expr unit` expressions and temporal arithmetic.

- [x] `SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM t` — DAY
- [x] `SELECT DATE_ADD(d, INTERVAL 1 HOUR) FROM t` — HOUR
- [x] `SELECT DATE_ADD(d, INTERVAL 1 MINUTE) FROM t` — MINUTE
- [x] `SELECT DATE_ADD(d, INTERVAL 1 SECOND) FROM t` — SECOND
- [x] `SELECT DATE_ADD(d, INTERVAL 1 MONTH) FROM t` — MONTH
- [x] `SELECT DATE_ADD(d, INTERVAL 1 WEEK) FROM t` — WEEK
- [x] `SELECT DATE_ADD(d, INTERVAL 1 QUARTER) FROM t` — QUARTER
- [x] `SELECT DATE_ADD(d, INTERVAL 1 YEAR) FROM t` — YEAR (already a keyword)
- [x] `SELECT DATE_ADD(d, INTERVAL 1 MICROSECOND) FROM t` — MICROSECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '1:30' HOUR_MINUTE) FROM t` — compound HOUR_MINUTE
- [x] `SELECT DATE_ADD(d, INTERVAL '1 1:30:00' DAY_SECOND) FROM t` — compound DAY_SECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '1-6' YEAR_MONTH) FROM t` — compound YEAR_MONTH
- [x] `SELECT DATE_ADD(d, INTERVAL '1:30:00.5' HOUR_MICROSECOND) FROM t` — compound HOUR_MICROSECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '1 12' DAY_HOUR) FROM t` — compound DAY_HOUR
- [x] `SELECT DATE_ADD(d, INTERVAL '1 12:30' DAY_MINUTE) FROM t` — compound DAY_MINUTE
- [x] `SELECT DATE_ADD(d, INTERVAL '30:00' MINUTE_SECOND) FROM t` — compound MINUTE_SECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '30:00.5' MINUTE_MICROSECOND) FROM t` — compound MINUTE_MICROSECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '59.999999' SECOND_MICROSECOND) FROM t` — compound SECOND_MICROSECOND
- [x] `SELECT DATE_ADD(d, INTERVAL '1:30:59' HOUR_SECOND) FROM t` — compound HOUR_SECOND
- [x] `SELECT d + INTERVAL 1 DAY FROM t` — interval in arithmetic expression
- [x] `SELECT d - INTERVAL 1 MONTH FROM t` — interval subtraction

### 2.4 UTC Temporal Functions

MySQL 8.0 reserves UTC_DATE, UTC_TIME, UTC_TIMESTAMP as no-argument function calls (no parentheses required).

- [x] `SELECT UTC_DATE` — no-parens form
- [x] `SELECT UTC_DATE()` — with-parens form
- [x] `SELECT UTC_TIME` — no-parens form
- [x] `SELECT UTC_TIME()` — with-parens form
- [x] `SELECT UTC_TIME(3)` — with fsp argument
- [x] `SELECT UTC_TIMESTAMP` — no-parens form
- [x] `SELECT UTC_TIMESTAMP()` — with-parens form
- [x] `SELECT UTC_TIMESTAMP(6)` — with fsp argument

---

## Phase 3: Oracle Corpus Verification

Oracle-verified corpus: each scenario is a concrete SQL statement sent to both MySQL 8.0 (via testcontainers) and omni parser, asserting matching accept/reject decisions.

### 3.1 Oracle Corpus — Integer Types × Modifiers

Each scenario is a CREATE TABLE with one column exercising a specific type+modifier combination.

- [x] `CREATE TABLE t (a INT)` — baseline
- [x] `CREATE TABLE t (a INT UNSIGNED)` — unsigned
- [x] `CREATE TABLE t (a INT SIGNED)` — signed (redundant but valid)
- [x] `CREATE TABLE t (a INT(11))` — display width
- [x] `CREATE TABLE t (a INT UNSIGNED ZEROFILL)` — unsigned + zerofill
- [x] `CREATE TABLE t (a TINYINT SIGNED)` — tinyint signed
- [x] `CREATE TABLE t (a TINYINT UNSIGNED ZEROFILL)` — tinyint unsigned zerofill
- [x] `CREATE TABLE t (a SMALLINT(5) UNSIGNED)` — smallint with width
- [x] `CREATE TABLE t (a MEDIUMINT SIGNED)` — mediumint signed
- [x] `CREATE TABLE t (a BIGINT(20) UNSIGNED)` — bigint with width
- [x] `CREATE TABLE t (a INT1 UNSIGNED)` — shorthand + modifier
- [x] `CREATE TABLE t (a INT8 SIGNED)` — shorthand + signed

### 3.2 Oracle Corpus — Decimal & Float Types

- [x] `CREATE TABLE t (a DECIMAL)` — bare decimal
- [x] `CREATE TABLE t (a DECIMAL(10))` — single precision
- [x] `CREATE TABLE t (a DECIMAL(10,2))` — full precision
- [x] `CREATE TABLE t (a DECIMAL(10,2) UNSIGNED)` — oracle-verify unsigned decimal
- [x] `CREATE TABLE t (a NUMERIC(10,2))` — numeric synonym
- [x] `CREATE TABLE t (a DEC(10,2))` — dec synonym
- [x] `CREATE TABLE t (a FIXED(10,2))` — fixed synonym
- [x] `CREATE TABLE t (a FLOAT)` — bare float
- [x] `CREATE TABLE t (a FLOAT(10,2))` — float with precision
- [x] `CREATE TABLE t (a FLOAT(24))` — single float precision
- [x] `CREATE TABLE t (a FLOAT(25))` — becomes double precision
- [x] `CREATE TABLE t (a DOUBLE)` — double
- [x] `CREATE TABLE t (a DOUBLE PRECISION)` — double precision synonym
- [x] `CREATE TABLE t (a REAL)` — real synonym
- [x] `CREATE TABLE t (a FLOAT4)` — float4 synonym
- [x] `CREATE TABLE t (a FLOAT8)` — float8 synonym

### 3.3 Oracle Corpus — String & Binary Types

- [x] `CREATE TABLE t (a CHAR(10))` — fixed-length string
- [x] `CREATE TABLE t (a CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)` — with charset + collation
- [x] `CREATE TABLE t (a VARCHAR(255))` — variable string
- [x] `CREATE TABLE t (a TEXT)` — text
- [x] `CREATE TABLE t (a TEXT(1000))` — text with length
- [x] `CREATE TABLE t (a TINYTEXT)` — tinytext
- [x] `CREATE TABLE t (a MEDIUMTEXT CHARACTER SET latin1)` — mediumtext with charset
- [x] `CREATE TABLE t (a LONGTEXT)` — longtext
- [x] `CREATE TABLE t (a LONG)` — long synonym
- [x] `CREATE TABLE t (a LONG VARCHAR)` — long varchar synonym
- [x] `CREATE TABLE t (a BINARY(16))` — fixed binary
- [x] `CREATE TABLE t (a VARBINARY(255))` — variable binary
- [x] `CREATE TABLE t (a BLOB)` — blob
- [x] `CREATE TABLE t (a BLOB(1000))` — blob with length
- [x] `CREATE TABLE t (a TINYBLOB)` — tinyblob
- [x] `CREATE TABLE t (a MEDIUMBLOB)` — mediumblob
- [x] `CREATE TABLE t (a LONGBLOB)` — longblob
- [x] `CREATE TABLE t (a LONG VARBINARY)` — long varbinary synonym
- [x] `CREATE TABLE t (a NATIONAL CHAR(10))` — national char
- [x] `CREATE TABLE t (a NCHAR(10))` — nchar synonym
- [x] `CREATE TABLE t (a NVARCHAR(100))` — nvarchar synonym

### 3.4 Oracle Corpus — Date/Time, JSON, Spatial & Special Types

- [x] `CREATE TABLE t (a DATE)` — date
- [x] `CREATE TABLE t (a TIME)` — time
- [x] `CREATE TABLE t (a TIME(3))` — time with fsp
- [x] `CREATE TABLE t (a DATETIME)` — datetime
- [x] `CREATE TABLE t (a DATETIME(6))` — datetime with fsp
- [x] `CREATE TABLE t (a TIMESTAMP)` — timestamp
- [x] `CREATE TABLE t (a TIMESTAMP(3))` — timestamp with fsp
- [x] `CREATE TABLE t (a YEAR)` — year
- [x] `CREATE TABLE t (a BIT(8))` — bit
- [x] `CREATE TABLE t (a BOOL)` — boolean
- [x] `CREATE TABLE t (a JSON)` — json
- [x] `CREATE TABLE t (a SERIAL)` — serial (auto_increment bigint)
- [x] `CREATE TABLE t (a ENUM('a','b','c'))` — enum
- [x] `CREATE TABLE t (a SET('x','y','z'))` — set
- [x] `CREATE TABLE t (a GEOMETRY)` — geometry
- [x] `CREATE TABLE t (a POINT)` — point
- [x] `CREATE TABLE t (a LINESTRING)` — linestring
- [x] `CREATE TABLE t (a POLYGON)` — polygon
- [x] `CREATE TABLE t (a GEOMETRYCOLLECTION)` — geometrycollection

### 3.5 Oracle Corpus — Window Functions

- [x] `SELECT RANK() OVER w FROM t WINDOW w AS (ORDER BY id)` — named window
- [x] `SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t` — partition by
- [x] `SELECT LAG(val) OVER (ORDER BY id) FROM t` — lag without offset (default 1)
- [x] `SELECT LAG(val, 2, 'N/A') OVER (ORDER BY id) FROM t` — lag with all args
- [x] `SELECT NTH_VALUE(val, 3) FROM FIRST OVER (ORDER BY id) FROM t` — with FROM FIRST
- [x] `SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t` — aggregate as window
- [x] `SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t` — multiple window functions

### 3.6 Oracle Corpus — Interval Expressions

- [x] `SELECT NOW() + INTERVAL 1 DAY` — interval in expression
- [x] `SELECT NOW() - INTERVAL 2 HOUR` — interval subtraction
- [x] `SELECT DATE_ADD('2024-01-01', INTERVAL 1 MONTH)` — with literal date
- [x] `SELECT DATE_SUB(NOW(), INTERVAL 30 MINUTE)` — date_sub
- [x] `SELECT TIMESTAMPADD(MINUTE, 30, NOW())` — timestampadd form
- [~] `SELECT EXTRACT(HOUR FROM NOW())` — extract — MISMATCH: omni rejects, MySQL accepts (parser lacks EXTRACT support)
- [x] `SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY` — in WHERE clause
- [x] `SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()` — interval in BETWEEN

### 3.7 Oracle Corpus — Rejected SQL

Verify omni matches MySQL 8.0 rejection behavior.

- [x] `CREATE TABLE t (a VARCHAR UNSIGNED)` — rejected: UNSIGNED on string type
- [x] `CREATE TABLE t (a TEXT ZEROFILL)` — rejected: ZEROFILL on text type
- [x] `CREATE TABLE t (a JSON UNSIGNED)` — rejected: UNSIGNED on JSON
- [~] `CREATE TABLE select (a INT)` — MISMATCH: omni accepts, MySQL rejects (parser too lenient with reserved words)
- [~] `CREATE TABLE t (select INT)` — MISMATCH: omni accepts, MySQL rejects (parser too lenient with reserved words)
- [x] `ALTER TABLE t PARTITION BY` — rejected: incomplete partition spec
- [x] `ALTER TABLE t PARTITION BY RANGE(id)` — rejected: missing partition definitions for RANGE
- [~] `SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t` — MISMATCH: omni accepts, MySQL rejects (parser doesn't validate interval units)
