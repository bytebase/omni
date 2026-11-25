# Cloud Spanner GoogleSQL truth1 — DML

All forms extracted from official Cloud Spanner GoogleSQL DML reference.
Primary source: https://cloud.google.com/spanner/docs/dml-syntax
version: spanner-current

Note: **MERGE is NOT supported** in Cloud Spanner GoogleSQL DML.

---

## DML-001: INSERT — basic value list

```
[statement_hint_expr]
INSERT [INTO] [schema_name.]table_name
  ( column_name [, ...] )
  VALUES ( value_expr [, ...] ) [, ( value_expr [, ...] ) [, ...] ]
  [THEN RETURN [WITH ACTION [AS alias]] select_list]
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#insert-syntax

```sql
INSERT INTO Singers (SingerId, FirstName, LastName)
VALUES (1, 'Marc', 'Richards');

INSERT INTO Singers (SingerId, FirstName, LastName)
VALUES (1, 'Marc', 'Richards'),
       (2, 'Catalina', 'Smith');
```

---

## DML-002: INSERT OR IGNORE

```
INSERT OR IGNORE [INTO] [schema_name.]table_name
  ( column_name [, ...] )
  VALUES ( value_expr [, ...] ) [, ...]
  [THEN RETURN ...]
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#insert-or-ignore

```sql
INSERT OR IGNORE INTO Singers (SingerId, FirstName, LastName)
VALUES (1, 'Marc', 'Richards');
-- Silently skips if SingerId=1 already exists
```

---

## DML-003: INSERT OR UPDATE

```
INSERT OR UPDATE [INTO] [schema_name.]table_name
  ( column_name [, ...] )
  VALUES ( value_expr [, ...] ) [, ...]
  [THEN RETURN ...]
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#insert-or-update

```sql
INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)
VALUES (1, 'Marc', 'Richards');
-- Inserts or updates matching row by primary key
```

---

## DML-004: INSERT ... SELECT

```
[statement_hint_expr]
INSERT [INTO] [schema_name.]table_name
  ( column_name [, ...] )
  query_statement
  [THEN RETURN ...]
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#insert-select

```sql
INSERT INTO ArchivedSingers (SingerId, FirstName, LastName)
SELECT SingerId, FirstName, LastName FROM Singers WHERE Retired = TRUE;
```

---

## DML-005: THEN RETURN clause (INSERT/UPDATE/DELETE)

```
THEN RETURN [WITH ACTION [AS action_alias]]
  { * | select_expression [, ...] }
  [EXCEPT (column_name [, ...])]
  [REPLACE (expression AS column_name [, ...])]

-- WITH ACTION returns a string column containing 'INSERT', 'UPDATE', or 'DELETE'
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#then-return

```sql
INSERT INTO Singers (SingerId, FirstName, LastName)
VALUES (3, 'Alice', 'Liu')
THEN RETURN SingerId, FirstName;

INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)
VALUES (1, 'Marc', 'Richards')
THEN RETURN WITH ACTION AS op_type *;

DELETE FROM Singers WHERE SingerId = 5
THEN RETURN SingerId;
```

---

## DML-006: UPDATE

```
[statement_hint_expr]
UPDATE [schema_name.]table_name [table_hint_expr]
  [[AS] table_alias]
  SET set_clause [, ...]
  [FROM from_clause]
  WHERE bool_expression
  [THEN RETURN [WITH ACTION [AS alias]] select_list]

set_clause:
    column_name = value_expression
  | column_name = DEFAULT
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#update-syntax

```sql
UPDATE Singers
SET FirstName = 'Alice'
WHERE SingerId = 1;

UPDATE Albums
SET MarketingBudget = MarketingBudget * 1.1
WHERE SingerId = 1
  AND AlbumId IN (SELECT AlbumId FROM TopAlbums);
```

---

## DML-007: UPDATE with FROM clause

```
UPDATE [schema_name.]table_name [[AS] alias]
  SET set_clause [, ...]
  FROM from_item [, ...]
  WHERE bool_expression
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#update-from

```sql
UPDATE Singers s
SET s.LastName = u.NewLastName
FROM Updates u
WHERE s.SingerId = u.SingerId;
```

---

## DML-008: DELETE

```
[statement_hint_expr]
DELETE [FROM] [schema_name.]table_name [table_hint_expr]
  [[AS] table_alias]
  WHERE bool_expression
  [THEN RETURN [WITH ACTION [AS alias]] select_list]
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#delete-syntax

```sql
DELETE FROM Singers WHERE SingerId = 1;

DELETE FROM Albums
WHERE SingerId NOT IN (SELECT SingerId FROM Singers);

DELETE Singers AS s WHERE s.SingerId = 42
THEN RETURN s.FirstName, s.LastName;
```

---

## DML-009: Statement hint on DML

```
@{ hint_key = hint_value [, ...] }
DML_statement
```

source_url: https://cloud.google.com/spanner/docs/dml-syntax#statement-hints

```sql
@{PDML_MAX_PARALLELISM=4}
DELETE FROM OldEvents WHERE CreatedAt < '2020-01-01';

@{LOCK_SCANNED_RANGES=shared}
UPDATE Accounts SET Balance = 0 WHERE Status = 'closed';
```

---

## DML-010: MERGE — NOT SUPPORTED

MERGE is **not** supported in Cloud Spanner GoogleSQL. The dialect only supports INSERT, INSERT OR IGNORE, INSERT OR UPDATE, UPDATE, and DELETE.

source_url: https://cloud.google.com/spanner/docs/dml-syntax

```sql
-- MERGE is not available in Cloud Spanner GoogleSQL.
-- Use INSERT OR UPDATE or INSERT OR IGNORE instead.
```
