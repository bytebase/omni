# Cloud Spanner GoogleSQL truth1 — Query Syntax

All forms extracted from official Cloud Spanner GoogleSQL query syntax reference.
Primary source: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax
version: spanner-current

---

## QUERY-001: query_statement (top-level grammar)

```
query_statement:
  [statement_hint_expr]
  query_expr

query_expr:
  [WITH cte [, ...]]
  { select | ( query_expr ) | set_operation }
  [ORDER BY expression [{ ASC | DESC }] [NULLS { FIRST | LAST }] [, ...]]
  [LIMIT count [OFFSET skip_rows]]
  [FOR UPDATE]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#sql_syntax

```sql
SELECT SingerId, FirstName FROM Singers ORDER BY LastName LIMIT 10 OFFSET 5;

(SELECT 1 AS n UNION ALL SELECT 2 AS n) ORDER BY n DESC;
```

---

## QUERY-002: SELECT statement

```
SELECT
  [{ ALL | DISTINCT }]
  [AS { STRUCT | VALUE }]
  select_list

select_list:
  { select_all | select_expression } [, ...]

select_all:
  [expression.] *
  [EXCEPT ( column_name [, ...] )]
  [REPLACE ( expression AS column_name [, ...] )]

select_expression:
  expression [[AS] alias]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_list

```sql
SELECT *;
SELECT * EXCEPT (SingerInfo);
SELECT * REPLACE (UPPER(LastName) AS LastName);
SELECT DISTINCT LastName FROM Singers;
SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers;
SELECT AS VALUE SingerId FROM Singers;
```

---

## QUERY-003: FROM clause

```
FROM from_item [, ...]

from_item:
    [schema_name.]table_name [table_hint_expr] [[AS] alias]
  | join_operation
  | (join_operation)
  | ( query_expr ) [table_hint_expr] [[AS] alias]
  | field_path
  | unnest_operator
  | cte_name [table_hint_expr] [[AS] alias]
  [tablesample_operator]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#from_clause

```sql
SELECT s.FirstName FROM Singers AS s;
SELECT a.Title FROM myschema.Albums AS a;
SELECT * FROM (SELECT 1 AS x) AS sub;
```

---

## QUERY-004: JOIN operations

```
join_operation:
  { cross_join_operation | condition_join_operation }

cross_join_operation:
  from_item CROSS JOIN [join_hint_expr] from_item
  | from_item , from_item

condition_join_operation:
  from_item
    { [INNER] [[HASH] JOIN]
    | LEFT [OUTER] [[HASH] JOIN]
    | RIGHT [OUTER] [[HASH] JOIN]
    | FULL [OUTER] [[HASH] JOIN]
    }
    [join_hint_expr]
    from_item
  join_condition

join_condition:
    ON bool_expression
  | USING ( column_name [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_types

```sql
SELECT * FROM Singers s CROSS JOIN Albums a;

SELECT s.FirstName, a.Title
FROM Singers s INNER JOIN Albums a ON s.SingerId = a.SingerId;

SELECT s.FirstName, a.Title
FROM Singers s LEFT OUTER JOIN Albums a USING (SingerId);

SELECT s.FirstName, a.Title
FROM Singers s FULL OUTER JOIN Albums a ON s.SingerId = a.SingerId;

-- HASH join hint
SELECT * FROM Singers s HASH JOIN Albums a ON s.SingerId = a.SingerId;
```

---

## QUERY-005: WHERE clause

```
WHERE bool_expression
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#where_clause

```sql
SELECT * FROM Singers WHERE LastName = 'Richards';
SELECT * FROM Albums WHERE MarketingBudget > 10000 AND SingerId = 1;
```

---

## QUERY-006: GROUP BY clause

```
GROUP BY [group_hint_expr] group_by_specification

group_by_specification:
    expression [, ...]
  | ROLLUP ( expression [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#group_by_clause

```sql
SELECT LastName, COUNT(*) FROM Singers GROUP BY LastName;
SELECT SingerId, AlbumId, SUM(Revenue) FROM Sales GROUP BY ROLLUP (SingerId, AlbumId);
```

---

## QUERY-007: HAVING clause

```
HAVING bool_expression
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#having_clause

```sql
SELECT LastName, COUNT(*) AS cnt
FROM Singers
GROUP BY LastName
HAVING cnt > 1;
```

---

## QUERY-008: ORDER BY clause

```
ORDER BY expression [{ ASC | DESC }] [NULLS { FIRST | LAST }] [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#order_by_clause

```sql
SELECT * FROM Singers ORDER BY LastName ASC, FirstName DESC;
SELECT * FROM Singers ORDER BY BirthDate NULLS LAST;
```

---

## QUERY-009: LIMIT and OFFSET

```
LIMIT count [OFFSET skip_rows]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#limit_and_offset_clause

```sql
SELECT * FROM Singers LIMIT 5;
SELECT * FROM Singers LIMIT 10 OFFSET 20;
```

---

## QUERY-010: FOR UPDATE

```
FOR UPDATE
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#for_update

```sql
SELECT * FROM Singers WHERE SingerId = 1 FOR UPDATE;
```

---

## QUERY-011: Set operations — UNION

```
query_expr UNION { ALL | DISTINCT } query_expr
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators

```sql
SELECT SingerId FROM Singers
UNION ALL
SELECT SingerId FROM ArchivedSingers;

SELECT LastName FROM Singers
UNION DISTINCT
SELECT LastName FROM ArchivedSingers;
```

---

## QUERY-012: Set operations — INTERSECT

```
query_expr INTERSECT { ALL | DISTINCT } query_expr
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators

```sql
SELECT SingerId FROM Singers
INTERSECT DISTINCT
SELECT SingerId FROM Albums;
```

---

## QUERY-013: Set operations — EXCEPT

```
query_expr EXCEPT { ALL | DISTINCT } query_expr
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators

```sql
SELECT SingerId FROM Singers
EXCEPT DISTINCT
SELECT SingerId FROM BlacklistSingers;
```

---

## QUERY-014: WITH clause (CTEs)

```
WITH cte [, ...] query_expr

cte:
  cte_name [( column_name [, ...] )] AS ( query_expr )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#with_clause

```sql
WITH TopSingers AS (
  SELECT SingerId, FirstName, LastName
  FROM Singers
  ORDER BY LastName
  LIMIT 10
)
SELECT * FROM TopSingers WHERE FirstName LIKE 'A%';

-- Multiple CTEs
WITH
  Cte1 AS (SELECT 1 AS n),
  Cte2 AS (SELECT n + 1 AS n FROM Cte1)
SELECT * FROM Cte2;
```

---

## QUERY-015: WITH RECURSIVE (not supported note)

Cloud Spanner GoogleSQL does **not** support `WITH RECURSIVE`. The `RECURSIVE` keyword appears in the reserved keyword list but is not implemented for CTEs.

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#with_clause

```sql
-- WITH RECURSIVE is NOT supported in Cloud Spanner GoogleSQL.
-- RECURSIVE is a reserved keyword but recursive CTEs are not available.
```

---

## QUERY-016: Subqueries (scalar, EXISTS, IN, FROM)

```
-- Scalar subquery
( query_expr )

-- EXISTS subquery
EXISTS ( query_expr )

-- IN subquery
expression [NOT] IN ( query_expr )

-- Subquery in FROM
FROM ( query_expr ) [[AS] alias]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#subqueries

```sql
-- Scalar subquery
SELECT (SELECT MAX(MarketingBudget) FROM Albums) AS MaxBudget;

-- EXISTS
SELECT * FROM Singers WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerId = Singers.SingerId);

-- IN subquery
SELECT * FROM Singers WHERE SingerId IN (SELECT SingerId FROM Albums WHERE Title LIKE 'A%');

-- FROM subquery
SELECT * FROM (SELECT SingerId, COUNT(*) AS albums FROM Albums GROUP BY SingerId) AS sub
WHERE albums > 2;
```

---

## QUERY-017: UNNEST operator

```
unnest_operator:
    UNNEST ( array_expression ) [[AS] alias] [WITH OFFSET [[AS] alias]]
  | array_path [[AS] alias] [WITH OFFSET [[AS] alias]]

-- UNNEST expands an array into rows; WITH OFFSET adds 0-based index column
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#unnest_operator

```sql
SELECT item FROM UNNEST(['apple', 'banana', 'cherry']) AS item;

SELECT item, pos
FROM UNNEST(['apple', 'banana', 'cherry']) AS item WITH OFFSET AS pos;

-- Unnesting an array column
SELECT SingerId, album
FROM Singers, UNNEST(AlbumList) AS album;
```

---

## QUERY-018: Field path in FROM clause

```
-- A field path refers to nested struct/array fields
SELECT field_path FROM table_name

-- Can be used directly in FROM to unnest array fields
field_path
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#field_path

```sql
SELECT s.Name.FirstName, s.Name.LastName FROM StructTable AS s;
SELECT arr_elem FROM Table t, t.ArrayColumn arr_elem;
```

---

## QUERY-019: TABLESAMPLE operator

```
tablesample_operator:
  TABLESAMPLE sample_method ( sample_size { PERCENT | ROWS } )

sample_method:
    BERNOULLI
  | RESERVOIR
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#tablesample_operator

```sql
SELECT * FROM Singers TABLESAMPLE BERNOULLI (10 PERCENT);
SELECT * FROM Albums TABLESAMPLE RESERVOIR (100 ROWS);
```

---

## QUERY-020: Implicit comma join (cross join shorthand)

```
FROM table1, table2
-- equivalent to FROM table1 CROSS JOIN table2
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#implicit_alias

```sql
SELECT s.FirstName, a.Title
FROM Singers s, Albums a
WHERE s.SingerId = a.SingerId;
```

---

## QUERY-021: SELECT AS STRUCT

```
SELECT AS STRUCT expr [AS field_name] [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_as_struct

```sql
SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers LIMIT 1;
-- Returns rows of STRUCT type
```

---

## QUERY-022: SELECT AS VALUE

```
SELECT AS VALUE expr
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_as_value

```sql
SELECT AS VALUE SingerId FROM Singers;
-- Produces a value table (single-value rows)
```

---

## QUERY-023: QUALIFY clause (post-window filter)

Cloud Spanner GoogleSQL supports the QUALIFY clause for filtering results of window functions.

```
QUALIFY bool_expression
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#qualify_clause

```sql
SELECT SingerId, AlbumId,
  ROW_NUMBER() OVER (PARTITION BY SingerId ORDER BY MarketingBudget DESC) AS rn
FROM Albums
QUALIFY rn = 1;
```

---

## QUERY-024: WINDOW clause (named window definitions)

```
WINDOW window_name AS ( window_specification ) [, ...]

window_specification:
  [window_name]
  [PARTITION BY expression [, ...]]
  [ORDER BY expression [{ ASC | DESC }] [, ...]]
  [window_frame_clause]

window_frame_clause:
  { ROWS | RANGE } { frame_start | BETWEEN frame_start AND frame_end }

frame_start / frame_end:
    UNBOUNDED PRECEDING
  | integer PRECEDING
  | CURRENT ROW
  | integer FOLLOWING
  | UNBOUNDED FOLLOWING
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#window_clause

```sql
SELECT SingerId, AlbumId, MarketingBudget,
  AVG(MarketingBudget) OVER w AS AvgBudget
FROM Albums
WINDOW w AS (PARTITION BY SingerId ORDER BY AlbumId
             ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW);
```

---

## QUERY-025: Parenthesized query expression

```
( query_expr )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#parenthesized_query_expressions

```sql
(SELECT 1 AS x) UNION ALL (SELECT 2 AS x);
```

---

## QUERY-026: GraphSQL (GRAPH_TABLE operator) — skip/not in scope

The `GRAPH_TABLE` operator and property graph syntax (MATCH, GQL) are documented for Spanner but considered a distinct sub-language. The `GRAPH_TABLE` keyword appears in the reserved list.

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#graph_table_operator

```sql
-- GRAPH_TABLE operator (Spanner Graph - not detailed in this corpus)
-- SELECT * FROM GRAPH_TABLE (my_graph MATCH (n:Person)-[e]->(m) RETURN n.name, m.name);
```
