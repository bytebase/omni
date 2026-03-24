# MySQL Parser Error Quality & Recovery Scenarios

> Goal: (1) Fix all sites where parseIdentifier() errors are silently discarded via `_, _, _ :=` pattern, (2) Improve error message quality to include "at end of input" / "at or near" context, (3) Add position info to ParseError output
> Verification: `go build ./mysql/parser/ && go test ./mysql/parser/ -count=1`; error_test.go table-driven tests per scenario
> Reference sources: PG parser error patterns, MySQL 8.0 native error format

Status: [ ] pending, [x] passing, [~] partial

Note: The MySQL parser already returns errors (not nil,nil) for most truncated inputs.
The work is improving message quality and fixing the ~32 sites where errors are silently ignored.

---

## Phase 1: Error Infrastructure & Ignored-Error Fixes

Priority: fix real bugs (silently ignored errors) first, then improve message quality.

### 1.1 Error Infrastructure

```
[x] syntaxErrorAtCur() helper — returns "syntax error at end of input" for EOF, "syntax error at or near \"X\"" for other tokens
[x] syntaxErrorAtTok(tok) helper — same logic for arbitrary token
[x] ParseError.Error() includes position in output: "syntax error at end of input (line 1, column 10)"
[x] error_test.go test harness with parseExpectError(sql, expectedSubstring) helper
```

### 1.2 Ignored Errors — stmt.go (10 sites)

These dispatch functions call parseXxxStmt() after advance() and return the result directly without checking for nil. If the sub-parser silently fails, the error is lost.

```
[x] parseChangeMasterStmtInner — propagate error after advance past MASTER
[x] parseStartReplicaStmt — propagate error after advance past REPLICA/SLAVE
[x] parseStartGroupReplicationStmt — propagate error after advance past GROUP_REPLICATION
[x] parseStopReplicaStmt — propagate error after advance past REPLICA/SLAVE
[x] parseStopGroupReplicationStmt — propagate error after advance past GROUP_REPLICATION
[x] parseLoadIndexIntoCacheStmt — propagate error after advance past INDEX
[x] parseCreateResourceGroupStmt — propagate error after advance past RESOURCE
[x] parseAlterResourceGroupStmt — propagate error after advance past RESOURCE
[x] parseDropSpatialRefSysStmt — propagate error after advance past SPATIAL
[x] parseDropResourceGroupStmt — propagate error after advance past RESOURCE
```

### 1.3 Ignored Errors — DDL Files (9 sites)

parseIdentifier() returns (string, int, error) but error is discarded via `_, _, _ :=`.

```
[x] create_function.go: LANGUAGE identifier — check error
[x] create_function.go: SQL SECURITY identifier — check error
[~] create_table.go: table option identifier (2 sites) — 1 of 2 fixed; consumeOptionValue deferred (20+ callers)
[x] create_database.go: charset/collation identifier (3 sites) — check error
[x] alter_table.go: AFTER column identifier — check error
[x] create_index.go: index type identifier — already passing
```

### 1.4 Ignored Errors — DML & Other Files (13 sites)

```
[x] select.go: table alias parseIdentifier (2 sites) — check error
[x] set_show.go: parseSetDefaultRoleStmt — already passing
[x] set_show.go: parseSetRoleStmt — already passing
[x] grant.go: privilege/object identifier (4 sites) — check error
[x] load_data.go: charset/file identifier (2 sites) — check error
[x] replication.go: channel identifier — check error
[x] utility.go: identifier (2 sites) — check error
[x] expr.go: window spec identifier — check error
```

## Phase 2: Expression Error Messages

Improve error messages from generic "expected expression" to contextual "syntax error at end of input" / "at or near".

### 2.1 Arithmetic Operators

Each test verifies: parse returns error AND message contains "at end of input".

```
[x] `SELECT 1 +` → error contains "at end of input"
[x] `SELECT 1 -` → error contains "at end of input"
[x] `SELECT 1 *` → error contains "at end of input"
[x] `SELECT 1 /` → error contains "at end of input"
[x] `SELECT 1 %` → error contains "at end of input"
[x] `SELECT 1 DIV` → error contains "at end of input"
[x] `SELECT 1 MOD` → error contains "at end of input"
```

### 2.2 Comparison Operators

```
[x] `SELECT 1 =` → error contains "at end of input"
[x] `SELECT 1 <` → error contains "at end of input"
[x] `SELECT 1 >` → error contains "at end of input"
[x] `SELECT 1 <=` → error contains "at end of input"
[x] `SELECT 1 >=` → error contains "at end of input"
[x] `SELECT 1 <>` → error contains "at end of input"
[x] `SELECT 1 !=` → error contains "at end of input"
[x] `SELECT 1 <=>` → error contains "at end of input"
```

### 2.3 Logical & Bitwise Operators

```
[x] `SELECT 1 AND` → error contains "at end of input"
[x] `SELECT 1 OR` → error contains "at end of input"
[x] `SELECT 1 XOR` → error contains "at end of input"
[x] `SELECT 1 |` → error contains "at end of input"
[x] `SELECT 1 &` → error contains "at end of input"
[x] `SELECT 1 ^` → error contains "at end of input"
[x] `SELECT 1 <<` → error contains "at end of input"
[x] `SELECT 1 >>` → error contains "at end of input"
```

### 2.4 Unary & Prefix Operators

```
[x] `SELECT NOT` → error contains "at end of input"
[x] `SELECT -` → error contains "at end of input"
[x] `SELECT ~` → error contains "at end of input"
[x] `SELECT !` → error contains "at end of input"
```

## Phase 3: Predicate & Special Expression Error Messages

### 3.1 BETWEEN, LIKE, IN, REGEXP

```
[ ] `SELECT 1 BETWEEN` → error contains "at end of input"
[ ] `SELECT 1 BETWEEN 0 AND` → error contains "at end of input"
[ ] `SELECT 1 NOT BETWEEN` → error contains "at end of input"
[ ] `SELECT 'a' LIKE` → error contains "at end of input"
[ ] `SELECT 'a' LIKE 'b' ESCAPE` → error contains "at end of input"
[ ] `SELECT 'a' NOT LIKE` → error contains "at end of input"
[ ] `SELECT 1 IN (` → error contains "at end of input"
[ ] `SELECT 1 IN (1,` → error contains "at end of input"
[ ] `SELECT 1 REGEXP` → error contains "at end of input"
[ ] `SELECT 1 NOT REGEXP` → error contains "at end of input"
```

### 3.2 CASE, CAST, IS, JSON, COLLATE

```
[ ] `SELECT CASE WHEN` → error contains "at end of input"
[ ] `SELECT CASE WHEN 1 THEN` → error contains "at end of input"
[ ] `SELECT CAST(1 AS` → error contains "at end of input"
[ ] `SELECT CONVERT(1 USING` → error contains "at end of input"
[ ] `SELECT 1 IS` → error contains "at end of input"
[ ] `SELECT 1 IS NOT` → error contains "at end of input"
[ ] `SELECT a COLLATE` → error contains "at end of input"
[ ] `SELECT a SOUNDS` → error contains "at end of input"
```

### 3.3 Mid-token Errors ("at or near")

When truncation happens at a non-EOF token, error should include the problematic token.

```
[ ] `SELECT CAST(1 + AS int)` → error contains "at or near" with "AS"
[ ] `SELECT 1 AND FROM t` → error contains "at or near" with "FROM"
[ ] `SELECT 1 + FROM t` → error contains "at or near" with "FROM"
```

## Phase 4: Statement-Level Error Messages

### 4.1 SELECT Clauses

```
[ ] `SELECT * FROM t JOIN` → error contains "at end of input"
[ ] `SELECT * FROM t LEFT JOIN` → error contains "at end of input"
[ ] `SELECT * FROM t JOIN t2 ON` → error contains "at end of input"
[ ] `SELECT * FROM t WHERE` → error contains "at end of input"
[ ] `SELECT * FROM t GROUP BY` → error contains "at end of input"
[ ] `SELECT * FROM t GROUP BY a,` → error contains "at end of input"
[ ] `SELECT * FROM t HAVING` → error contains "at end of input"
[ ] `SELECT * FROM t ORDER BY` → error contains "at end of input"
[ ] `SELECT * FROM t ORDER BY a,` → error contains "at end of input"
[ ] `SELECT * FROM t LIMIT` → error contains "at end of input"
[ ] `SELECT 1 UNION` → error contains "at end of input"
[ ] `SELECT 1 UNION ALL` → error contains "at end of input"
```

### 4.2 CTE & Subqueries

```
[ ] `WITH` → error contains "at end of input"
[ ] `WITH cte AS (` → error contains "at end of input"
[ ] `SELECT EXISTS (` → error contains "at end of input"
```

### 4.3 DDL Truncation

```
[ ] `CREATE TABLE t (a INT DEFAULT` → error contains "at end of input"
[ ] `CREATE TABLE t (a INT REFERENCES` → error contains "at end of input"
[ ] `CREATE TABLE t (a INT CHECK (` → error contains "at end of input"
[ ] `ALTER TABLE t ADD COLUMN` → error contains "at end of input"
[ ] `ALTER TABLE t DROP COLUMN` → error contains "at end of input"
[ ] `ALTER TABLE t RENAME TO` → error contains "at end of input"
[ ] `CREATE INDEX idx ON` → error contains "at end of input"
[ ] `CREATE VIEW v AS` → error contains "at end of input"
```

### 4.4 DML Truncation

```
[ ] `INSERT INTO t VALUES (` → error contains "at end of input"
[ ] `INSERT INTO t (` → error contains "at end of input"
[ ] `INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a =` → error contains "at end of input"
[ ] `UPDATE t SET` → error contains "at end of input"
[ ] `UPDATE t SET a =` → error contains "at end of input"
[ ] `DELETE FROM t WHERE` → error contains "at end of input"
[ ] `SET @x =` → error contains "at end of input"
[ ] `USE` → error contains "at end of input"
[ ] `DROP TABLE` → error contains "at end of input"
```
