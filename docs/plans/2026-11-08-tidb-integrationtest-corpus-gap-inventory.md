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

| | first | after OA fix | after DESC fix |
|---|---:|---:|---:|
| both accept | 26,941 | 26,941 | 27,152 |
| gap: TiDB accepts, omni rejects | 268 | 268 | 57 |
| **over-acceptance: TiDB rejects, omni accepts** | **6** | **0** | **0** |
| both reject | 31 | 37 | 37 |
| no verdict | 0 | 0 | 0 |

omni parses **27,152 of the 27,209 statements TiDB accepts — 99.8%** (was 99.0%),
with **zero over-acceptances**. The both-reject rows are statements the
server itself refuses (`load stats` with no path, `foreign key … reference`
misspelled, `require ssl ssl`, a `clustered` column attribute in the wrong place,
a `case` expression missing a space); they are not work.

## Over-acceptances — 6 rows — ✅ CLOSED

These were the half §1 of the design calls most at risk: omni accepting what the
server rejects. All six are closed by `fix(tidb): reject the integrationtest
over-acceptances`, each measured against live TiDB v8.5.5
(`over_acceptance_fixes_oracle_test.go`). The fix is tidb-only: the corpus is
TiDB's, so these are TiDB verdicts; the forks' analogues were not measured here
and are a possible follow-up.

### OA-1. `IDENTIFIED BY PASSWORD '<hash>'` — 1 row — ✅ CLOSED

```
ALTER USER USER() IDENTIFIED BY PASSWORD '*B50FBDB37F1256824274912F2A1CE648082C3F1F'
```

Only the `ALTER USER USER()` self-form is affected: its grammar is `ALTER USER
USER() IDENTIFIED BY AuthString` (a plain string, `parser.y:13745`), so `BY
PASSWORD '<hash>'` there is 1064. `CREATE USER` and named `ALTER USER` still take
`IDENTIFIED BY PASSWORD` — omni was calling the general auth-option parser for
the `USER()` form, which accepted it; the fix parses only `IDENTIFIED BY <string>`
there.

### OA-2. Function-name keyword as a bare table name — 3 rows — ✅ CLOSED

```
CREATE TABLE COUNT(a bigint)
CREATE TABLE BIT_AND(a bigint)
CREATE TABLE NOW(a bigint)
```

A function-name keyword immediately followed by `(` is a function call, so it
cannot be a table name there: `CREATE TABLE COUNT(a int)` is 1064 while `CREATE
TABLE COUNT (a int)` (a space before the paren) and a backquoted name are fine.
omni took `COUNT(` as an ordinary identifier. The fix refuses the 26 aggregate
and temporal-scalar function names (measured against v8.5.5) in the CREATE TABLE
table-name position when `(` is adjacent — not a blanket reservation, so they
stay legal column names and spaced/quoted table names.

### OA-3. Out-of-range numeric literal at parse time — 2 rows — ✅ CLOSED

```
select * from select_limit limit 18446744073709551616 offset 3
insert into tt1 values(4556414e723532)
```

`18446744073709551616` is 2^64, one past the `LIMIT` clause's unsigned range,
and `4556414e723532` is a float whose exponent overflows; TiDB rejects both
(1064 / out of range). Two narrow checks, where TiDB applies them: `parseLimitValue`
rejects a `LIMIT`/`OFFSET` integer past `uint64` (a bare integer literal that
large stays a decimal elsewhere, as TiDB keeps it), and `scanNumber` rejects a
float literal that overflows the double to infinity (underflow to zero is kept).

## The gaps, grouped — 268 rows

### A. `DESC`/`DESCRIBE` as an `EXPLAIN` synonym — 211 rows — ✅ CLOSED

```
desc select * from t
desc format='brief' select /*+ HASH_AGG() */ max(a) from t
desc analyze select * from t
desc update t set a = 1
desc with cte as (select 1) select * from cte
```

**Closed** by `fix(tidb): read DESC/DESCRIBE as full EXPLAIN synonyms`. omni read
`DESC`/`DESCRIBE` only as the describe-a-table form; TiDB takes all three of
`EXPLAIN`/`DESCRIBE`/`DESC` as the full `ExplainStmt` (`DESC [ANALYZE] [FORMAT =
fmt] explainable_stmt` over `SELECT`/`UPDATE`/`DELETE`/`WITH`/`INSERT`/`TABLE`/
`ALTER TABLE`, plus `FOR CONNECTION`). They are the one `ExplainSym` in parser.y,
and measured on v8.5.5 all three take exactly the same forms, so `parseExplainStmt`
now shares its body instead of short-circuiting `DESC`/`DESCRIBE` to the table
form. The table form (`DESC tbl [col]`) still parses. tidb-only: MySQL and the
forks take `DESC`/`DESCRIBE` only as the table form (a TiDB dialect delta).

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

The 6 over-acceptances (OA-1..3) are closed. Next, the gaps by size: A
(`DESC`/`DESCRIBE` EXPLAIN synonym, 211) buys the most for the least grammar, B
(interval expression, 51) is the one to port across the forks. Each closes as
§3.2's rows did — a small PR with §3.5's tree checks and zero new
over-acceptances — then re-record this corpus with `-update` and read the diff.
Sample 3 (Bytebase's own TiDB test files) is tracked separately.
