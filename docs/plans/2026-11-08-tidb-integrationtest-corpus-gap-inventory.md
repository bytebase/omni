# TiDB integrationtest corpus gap inventory (§3.4 sample 2)

Date: 2026-11-08. The first measurement of §3.4 sample 2, "TiDB's integration
tests" (`tidb/parser/testdata/tidb_integrationtest_corpus.tsv`). Sample 1 (the
rest of pingcap's parser tests, `2026-11-08-tidb-real-sql-corpus-gap-inventory.md`)
is down to its last gap; this is the larger real-SQL list — the statements people
wrote to exercise the TiDB server.

## The measurement

27,246 statements, lifted from `tests/integrationtest/t/**/*.test` at the v8.5.5
tag (239 files): mysqltest directive/control lines dropped, `tidb/parser.Split`
run on the rest, single-line only, deduplicated within this corpus and against
the two existing corpora. Each row was measured against omni and against TiDB
v8.5.5 via `PREPARE`.

| | count |
|---|---:|
| both accept | 26,941 |
| gap: TiDB accepts, omni rejects | 268 |
| **over-acceptance: TiDB rejects, omni accepts** | **6** |
| both reject | 31 |
| no verdict | 0 |

omni parses **26,941 of the 27,209 statements TiDB accepts — 99.0%**, and
**over-accepts 6 of the 26,947 it parses — 0.02%**. The 31 both-reject rows are
statements the server itself refuses (`load stats` with no path, `foreign key …
reference` misspelled, `require ssl ssl`, a `clustered` column attribute in the
wrong place, a `case` expression missing a space); they are not work.

## Over-acceptances — 6 rows — bugs, fix first

These are the half §1 of the design calls most at risk: omni accepts what the
server rejects. Each closes like a gap, checked against mysql, mariadb and
oceanbase in the same PR.

### OA-1. `IDENTIFIED BY PASSWORD '<hash>'` — 1 row

```
ALTER USER USER() IDENTIFIED BY PASSWORD '*B50FBDB37F1256824274912F2A1CE648082C3F1F'
```

TiDB v8.5.5 dropped the MySQL-5 `IDENTIFIED BY PASSWORD '<hash>'` auth-option
form (it answers 1064). omni still reads it.

### OA-2. Function-name keyword as a bare table name — 3 rows

```
CREATE TABLE COUNT(a bigint)
CREATE TABLE BIT_AND(a bigint)
CREATE TABLE NOW(a bigint)
```

`COUNT`, `BIT_AND`, `NOW` and the other built-in function names TiDB treats as
reserved cannot be a table name unquoted; the server answers 1064. omni takes
them as an ordinary identifier. (The fix is a table-name position that refuses
the function-keyword set, not a blanket reservation — these words are still
legal column names.)

### OA-3. Out-of-range numeric literal at parse time — 2 rows

```
select * from select_limit limit 18446744073709551616 offset 3
insert into tt1 values(4556414e723532)
```

`18446744073709551616` is 2^64, one past the `LIMIT` clause's unsigned range,
and `4556414e723532` is a float whose exponent overflows; TiDB rejects both
while lexing (1064 / out of range). omni accepts the token. This is a lexer
range check, narrow to the contexts where TiDB applies it.

## The gaps, grouped — 268 rows

### A. `DESC`/`DESCRIBE` as an `EXPLAIN` synonym — 211 rows — the big one

```
desc select * from t
desc format='brief' select /*+ HASH_AGG() */ max(a) from t
desc analyze select * from t
desc update t set a = 1
desc with cte as (select 1) select * from cte
```

omni reads `DESC`/`DESCRIBE` only as the describe-a-table form (`DESC tbl
[col]`); the code says so at `tidb/parser/set_show.go:1850`. TiDB takes all three
of `EXPLAIN`/`DESCRIBE`/`DESC` as the full `ExplainStmt`: `DESC [ANALYZE]
[FORMAT = fmt] explainable_stmt`, where the statement is a `SELECT`, `UPDATE`,
`DELETE`, `WITH`, `INSERT`, etc. omni already parses that whole grammar after
`EXPLAIN`; the gap is reaching it from `DESC`/`DESCRIBE` too, after the table
form fails to apply. (MySQL and the forks take `DESC`/`DESCRIBE` only as the
table form, so this is a TiDB dialect delta — check before porting.)

### B. `INTERVAL expr unit` with a signed, decimal or operator quantity — 51 rows

```
SELECT "1900-01-01 00:00:00" + INTERVAL 1 << 37 SECOND
SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2.2 MINUTE_MICROSECOND)
SELECT * FROM t WHERE date_add(ts, interval -20 minute) = '2025-03-30 02:40:00'
insert into t2 values('-9223372036854775809', "1000-01-01 00:00:00" + INTERVAL -9223372036854775809 YEAR)
```

TiDB's interval quantity is a full `Expression` (`INTERVAL Expression TimeUnit`),
so a shift (`1 << 37`), a signed or decimal literal (`-2.2`, `-9223372036854775809`)
and any other expression are all legal, with every compound unit
(`MINUTE_MICROSECOND`, `DAY_SECOND`, `YEAR_MONTH`, …). omni reads a narrower
expression in the interval position. This is the one family with sibling reach:
the same construct appears in mysql/mariadb/oceanbase — port it across the forks
in one PR where they agree.

### C. `\N` as the `NULL` literal — 3 rows

```
select \N
select \N from test
select (\N) from test
```

MySQL and TiDB lex `\N` as a synonym for `NULL`. omni's lexer does not.

### D. `DEFAULT(qualified.column)` — 1 row

```
select default(t1.a) from t1, t2
```

The `DEFAULT(col)` function takes a qualified column reference (`t1.a`), not only
a bare column name.

### E. CREATE/ALTER table-option combinations — 2 rows

```
ALTER TABLE t TTL_ENABLE = 'OFF' REMOVE TTL
create table ttt4(a varchar(123) default null collate utf8mb4_unicode_ci)engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci
```

Two option combinations omni rejects: `TTL_ENABLE = …` followed by `REMOVE TTL`
in one `ALTER`, and a column `DEFAULT NULL COLLATE …` next to table-level
`ENGINE` / `DEFAULT CHARSET` / `COLLATE` with no space before `)engine`. Each
needs its own look before it is counted a single family.

## Next

Close the 6 over-acceptances first (OA-1..3), then the gaps by size: A
(`DESC`/`DESCRIBE` EXPLAIN synonym, 211) buys the most for the least grammar, B
(interval expression, 51) is the one to port across the forks. Each closes as
§3.2's rows did — a small PR with §3.5's tree checks and zero new
over-acceptances — then re-record this corpus with `-update` and read the diff.
Sample 3 (Bytebase's own TiDB test files) is tracked separately.
