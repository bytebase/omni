# Cloud Spanner GoogleSQL truth1 — Query Hints

All forms extracted from official Cloud Spanner GoogleSQL hints reference.
Primary sources:
- https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#hints
- https://cloud.google.com/spanner/docs/query-execution-hints
version: spanner-current

---

## HINT-001: Statement hint syntax

```
@{ hint_key = hint_value [, hint_key = hint_value ...] }
statement

-- Hint block precedes the SQL statement
-- Multiple hints are comma-separated within a single @{...} block
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#statement_hints

```sql
@{USE_ADDITIONAL_PARALLELISM=TRUE}
SELECT * FROM Singers;

@{OPTIMIZER_VERSION=2, OPTIMIZER_STATISTICS_PACKAGE='auto_20240101'}
SELECT * FROM Albums;
```

---

## HINT-002: Table hint syntax

```
FROM table_name @{ hint_key = hint_value [, ...] } [[AS] alias]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#table_hints

```sql
SELECT * FROM Singers @{FORCE_INDEX=SingersByName};
SELECT * FROM Albums @{FORCE_INDEX=AlbumsByTitle, GROUPBY_SCAN_OPTIMIZATION=TRUE} AS a;
```

---

## HINT-003: Join hint syntax

```
from_item JOIN @{ hint_key = hint_value [, ...] } from_item ON condition
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_hints

```sql
SELECT *
FROM Singers s
JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON s.SingerId = a.SingerId;

SELECT *
FROM Singers s
JOIN @{JOIN_METHOD=APPLY_JOIN} Albums a ON s.SingerId = a.SingerId;
```

---

## HINT-004: FORCE_INDEX table hint

```
@{ FORCE_INDEX = index_name }
-- Forces use of specified index; use _BASE_TABLE to force base table scan
@{ FORCE_INDEX = _BASE_TABLE }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#force_index_hint

```sql
SELECT * FROM Singers @{FORCE_INDEX=SingersByLastName};
SELECT * FROM Singers @{FORCE_INDEX=_BASE_TABLE};
```

---

## HINT-005: GROUPBY_SCAN_OPTIMIZATION table hint

```
@{ GROUPBY_SCAN_OPTIMIZATION = { TRUE | FALSE } }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#groupby_scan_optimization_hint

```sql
SELECT SingerId, COUNT(*) FROM Albums @{GROUPBY_SCAN_OPTIMIZATION=TRUE} GROUP BY SingerId;
```

---

## HINT-006: OPTIMIZER_VERSION statement hint

```
@{ OPTIMIZER_VERSION = { integer | 'latest' } }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#optimizer_version_hint

```sql
@{OPTIMIZER_VERSION=4}
SELECT * FROM Singers;

@{OPTIMIZER_VERSION=latest}
SELECT * FROM Singers;
```

---

## HINT-007: OPTIMIZER_STATISTICS_PACKAGE statement hint

```
@{ OPTIMIZER_STATISTICS_PACKAGE = 'package_name' }
```

source_url: https://cloud.google.com/spanner/docs/query-execution-hints#optimizer_statistics_package_hint

```sql
@{OPTIMIZER_STATISTICS_PACKAGE='auto_20240601_14_00_00'}
SELECT * FROM Albums;
```

---

## HINT-008: USE_ADDITIONAL_PARALLELISM statement hint

```
@{ USE_ADDITIONAL_PARALLELISM = { TRUE | FALSE } }
```

source_url: https://cloud.google.com/spanner/docs/query-execution-hints#use_additional_parallelism

```sql
@{USE_ADDITIONAL_PARALLELISM=TRUE}
SELECT COUNT(*) FROM LargeTable;
```

---

## HINT-009: PDML_MAX_PARALLELISM DML hint

```
@{ PDML_MAX_PARALLELISM = integer }
-- Used with Partitioned DML (PDML) for UPDATE/DELETE
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#statement-hints

```sql
@{PDML_MAX_PARALLELISM=4}
DELETE FROM OldLogs WHERE CreatedAt < '2023-01-01';
```

---

## HINT-010: JOIN_METHOD join hint

```
JOIN @{ JOIN_METHOD = { HASH_JOIN | APPLY_JOIN | LOOP_JOIN | MERGE_JOIN } }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_hints

```sql
SELECT *
FROM Singers s
INNER JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON s.SingerId = a.SingerId;

SELECT *
FROM Singers s
LEFT JOIN @{JOIN_METHOD=APPLY_JOIN} Albums a ON s.SingerId = a.SingerId;
```

---

## HINT-011: LOCK_SCANNED_RANGES statement hint

```
@{ LOCK_SCANNED_RANGES = { exclusive | shared } }
```

source_url: https://cloud.google.com/spanner/docs/query-execution-hints#lock_scanned_ranges

```sql
@{LOCK_SCANNED_RANGES=exclusive}
SELECT * FROM Accounts WHERE AccountId = 42 FOR UPDATE;

@{LOCK_SCANNED_RANGES=shared}
UPDATE Counters SET Value = Value + 1 WHERE CounterId = 1;
```

---

## HINT-012: SCAN_METHOD table hint

```
@{ SCAN_METHOD = { BATCH | ROW } }
```

source_url: https://cloud.google.com/spanner/docs/query-execution-hints#scan_method

```sql
SELECT * FROM LargeTable @{SCAN_METHOD=BATCH};
```

---

## HINT-013: Multiple hints combined

```
@{ hint1 = val1, hint2 = val2 }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#hints

```sql
@{OPTIMIZER_VERSION=4, USE_ADDITIONAL_PARALLELISM=TRUE}
SELECT * FROM Singers @{FORCE_INDEX=SingersByName}
JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON Singers.SingerId = a.SingerId;
```
