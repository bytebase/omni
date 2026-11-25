# BigQuery GoogleSQL truth1 — Query Syntax

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax

---

## QUERY-001: query_statement (top-level grammar)

```
query_statement:
  query_expr

query_expr:
  [ WITH [ RECURSIVE ] { non_recursive_cte | recursive_cte }[, ...] ]
  { select | ( query_expr ) | set_operation }
  [ ORDER BY expression [{ ASC | DESC }] [, ...] ]
  [ LIMIT count [ OFFSET skip_rows ] ]

select:
  SELECT
    [ WITH differential_privacy_clause ]
    [ { ALL | DISTINCT } ]
    [ AS { STRUCT | VALUE } ]
    select_list
  [ FROM from_clause[, ...] ]
  [ WHERE bool_expression ]
  [ GROUP BY group_by_specification ]
  [ HAVING bool_expression ]
  [ QUALIFY bool_expression ]
  [ WINDOW window_clause ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#sql_syntax
version: bigquery-current

```sql
-- basic select
SELECT name, age FROM dataset.users WHERE age > 18 ORDER BY age DESC LIMIT 10;

-- parenthesized query_expr
SELECT * FROM (SELECT 1 AS n UNION ALL SELECT 2 AS n);
```

---

## QUERY-002: SELECT statement

```
SELECT
  [ WITH differential_privacy_clause ]
  [ { ALL | DISTINCT } ]
  [ AS { STRUCT | VALUE } ]
  select_list

select_list:
  { select_all | select_expression } [, ...]

select_all:
  [ expression. ]*
  [ EXCEPT ( column_name [, ...] ) ]
  [ REPLACE ( expression AS column_name [, ...] ) ]

select_expression:
  expression [ [ AS ] alias ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_statement
version: bigquery-current

```sql
SELECT * FROM Roster;
SELECT name AS n, age FROM dataset.users;
SELECT g.* FROM groceries AS g;
```

---

## QUERY-003: SELECT * EXCEPT

```
SELECT * EXCEPT ( column_name [, ...] )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_except
version: bigquery-current

```sql
WITH orders AS (SELECT 5 AS order_id, "sprocket" AS item_name, 200 AS quantity)
SELECT * EXCEPT (order_id) FROM orders;
```

---

## QUERY-004: SELECT * REPLACE

```
SELECT * REPLACE ( expression AS column_name [, ...] )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_replace
version: bigquery-current

```sql
WITH orders AS (SELECT 5 AS order_id, "sprocket" AS item_name, 200 AS quantity)
SELECT * REPLACE ("widget" AS item_name) FROM orders;

WITH orders AS (SELECT 5 AS order_id, "sprocket" AS item_name, 200 AS quantity)
SELECT * REPLACE (quantity/2 AS quantity) FROM orders;
```

---

## QUERY-005: SELECT DISTINCT

```
SELECT DISTINCT select_list ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_distinct
version: bigquery-current

```sql
SELECT DISTINCT Name FROM PlayerStats;
```

---

## QUERY-006: SELECT ALL

```
SELECT ALL select_list ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_all
version: bigquery-current

```sql
SELECT ALL name FROM dataset.users;
```

---

## QUERY-007: SELECT AS STRUCT

```
SELECT AS STRUCT select_list ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_as_struct
version: bigquery-current

```sql
SELECT AS STRUCT 1 AS a, 2 AS b;

DECLARE corpus_count, word_count INT64;
SET (corpus_count, word_count) = (
  SELECT AS STRUCT COUNT(DISTINCT corpus), SUM(word_count)
  FROM bigquery-public-data.samples.shakespeare
  WHERE LOWER(word) = 'methinks'
);
```

---

## QUERY-008: SELECT AS VALUE

```
SELECT AS VALUE select_expression ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_as_value
version: bigquery-current

```sql
SELECT AS VALUE STRUCT(1 AS a, 2 AS b);
```

---

## QUERY-009: SELECT expression.*

```
SELECT expression.* ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#select_expression_star
version: bigquery-current

```sql
WITH groceries AS (SELECT "milk" AS dairy, "eggs" AS protein, "bread" AS grain)
SELECT g.* FROM groceries AS g;

WITH locations AS (SELECT STRUCT("Seattle" AS city, "Washington" AS state) AS location UNION ALL
                   SELECT STRUCT("Phoenix" AS city, "Arizona" AS state) AS location)
SELECT l.location.* FROM locations l;
```

---

## QUERY-010: FROM clause

```
FROM from_clause[, ...]

from_clause:
  from_item
  [ { pivot_operator | unpivot_operator | match_recognize_clause } ]
  [ tablesample_operator ]

from_item:
  {
    table_name [ as_alias ] [ FOR SYSTEM_TIME AS OF timestamp_expression ]
    | { join_operation | ( join_operation ) }
    | ( query_expr ) [ as_alias ]
    | field_path
    | unnest_operator
    | cte_name [ as_alias ]
    | graph_table_operator [ as_alias ]
  }

as_alias:
  [ AS ] alias
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#from_clause
version: bigquery-current

```sql
SELECT * FROM Roster;
SELECT * FROM dataset.Roster;
SELECT * FROM project.dataset.Roster;
SELECT * FROM (SELECT "apple" AS fruit) AS t;
```

---

## QUERY-011: FOR SYSTEM_TIME AS OF (time travel)

```
table_name FOR SYSTEM_TIME AS OF timestamp_expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#for_system_time_as_of
version: bigquery-current

```sql
SELECT * FROM t FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR);
SELECT * FROM t FOR SYSTEM_TIME AS OF '2017-01-01 10:00:00-07:00';
```

---

## QUERY-012: UNNEST operator

```
unnest_operator:
  {
    UNNEST( array ) [ as_alias ]
    | array_path [ as_alias ]
  }
  [ WITH OFFSET [ as_alias ] ]

array:
  { array_expression | array_path }

as_alias:
  [AS] alias
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#unnest_operator
version: bigquery-current

```sql
SELECT * FROM UNNEST([1, 2, 3]) AS num;
SELECT num, pos FROM UNNEST([10, 20, 30]) AS num WITH OFFSET AS pos;
SELECT * FROM T1 t1, t1.array_column;
```

---

## QUERY-013: PIVOT operator

```
FROM from_item[, ...] pivot_operator

pivot_operator:
  PIVOT(
    aggregate_function_call [as_alias][, ...]
    FOR input_column
    IN ( pivot_column [as_alias][, ...] )
  ) [AS alias]

as_alias:
  [AS] alias
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#pivot_operator
version: bigquery-current

```sql
SELECT * FROM
  Produce
  PIVOT(SUM(sales) FOR quarter IN ('Q1', 'Q2', 'Q3', 'Q4'));

SELECT * FROM
  (SELECT product, sales, quarter FROM Produce)
  PIVOT(SUM(sales) AS total_sales, COUNT(*) AS num_records FOR quarter IN ('Q1', 'Q2'));
```

---

## QUERY-014: UNPIVOT operator

```
FROM from_item[, ...] unpivot_operator

unpivot_operator:
  UNPIVOT [ { INCLUDE NULLS | EXCLUDE NULLS } ] (
    { single_column_unpivot | multi_column_unpivot }
  ) [unpivot_alias]

single_column_unpivot:
  values_column
  FOR name_column
  IN (columns_to_unpivot)

multi_column_unpivot:
  values_column_set
  FOR name_column
  IN (column_sets_to_unpivot)

values_column_set:
  (values_column[, ...])

columns_to_unpivot:
  unpivot_column [row_value_alias][, ...]

column_sets_to_unpivot:
  (unpivot_column [row_value_alias][, ...])

unpivot_alias and row_value_alias:
  [AS] alias
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#unpivot_operator
version: bigquery-current

```sql
SELECT * FROM sales_table
  UNPIVOT(total_sales FOR quarter IN (Q1, Q2, Q3, Q4));

SELECT * FROM sales_table
  UNPIVOT INCLUDE NULLS (total_sales FOR quarter IN (Q1, Q2, Q3, Q4));
```

---

## QUERY-015: TABLESAMPLE operator

```
TABLESAMPLE SYSTEM ( percent PERCENT )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#tablesample_operator
version: bigquery-current

```sql
SELECT * FROM dataset.my_table TABLESAMPLE SYSTEM (10 PERCENT);
```

---

## QUERY-016: MATCH_RECOGNIZE clause

```
FROM from_item
MATCH_RECOGNIZE (
  [ PARTITION BY partition_expr [, ... ] ]
  ORDER BY order_expr [{ ASC | DESC }] [{ NULLS FIRST | NULLS LAST }] [, ...]
  MEASURES { measures_expr [AS] alias } [, ... ]
  [ AFTER MATCH SKIP { PAST LAST ROW | TO NEXT ROW } ]
  PATTERN (pattern)
  DEFINE symbol AS boolean_expr [, ... ]
  [ OPTIONS ( [ use_longest_match = { TRUE | FALSE } ] ) ]
)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#match_recognize_clause
version: bigquery-current

```sql
SELECT *
FROM stock_prices
MATCH_RECOGNIZE (
  PARTITION BY symbol
  ORDER BY date
  MEASURES FIRST(price) AS first_price, LAST(price) AS last_price
  PATTERN (RISE+ FALL+)
  DEFINE
    RISE AS price > LAG(price),
    FALL AS price < LAG(price)
);
```

---

## QUERY-017: [INNER] JOIN

```
from_item [INNER] JOIN from_item join_condition

join_condition:
  { on_clause | using_clause }

on_clause:
  ON bool_expression

using_clause:
  USING ( column_list )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#inner_join
version: bigquery-current

```sql
FROM A INNER JOIN B ON A.w = B.y
FROM A INNER JOIN B USING (x)
SELECT Roster.LastName, TeamMascot.Mascot FROM Roster JOIN TeamMascot ON Roster.SchoolID = TeamMascot.SchoolID;
```

---

## QUERY-018: CROSS JOIN

```
from_item CROSS JOIN from_item
-- or comma cross join:
from_item , from_item
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#cross_join
version: bigquery-current

```sql
FROM A CROSS JOIN B
FROM A, B
SELECT Roster.LastName, TeamMascot.Mascot FROM Roster CROSS JOIN TeamMascot;
```

---

## QUERY-019: FULL [OUTER] JOIN

```
from_item FULL [OUTER] JOIN from_item join_condition
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#full_outer_join
version: bigquery-current

```sql
FROM A FULL OUTER JOIN B ON A.w = B.y
FROM A FULL OUTER JOIN B USING (x)
SELECT Roster.LastName, TeamMascot.Mascot FROM Roster FULL JOIN TeamMascot ON Roster.SchoolID = TeamMascot.SchoolID;
```

---

## QUERY-020: LEFT [OUTER] JOIN

```
from_item LEFT [OUTER] JOIN from_item join_condition
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#left_outer_join
version: bigquery-current

```sql
FROM A LEFT OUTER JOIN B ON A.w = B.y
FROM A LEFT OUTER JOIN B USING (x)
SELECT Roster.LastName, TeamMascot.Mascot FROM Roster LEFT JOIN TeamMascot ON Roster.SchoolID = TeamMascot.SchoolID;
```

---

## QUERY-021: RIGHT [OUTER] JOIN

```
from_item RIGHT [OUTER] JOIN from_item join_condition
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#right_outer_join
version: bigquery-current

```sql
FROM A RIGHT OUTER JOIN B ON A.w = B.y
FROM A RIGHT OUTER JOIN B USING (x)
SELECT Roster.LastName, TeamMascot.Mascot FROM Roster RIGHT JOIN TeamMascot ON Roster.SchoolID = TeamMascot.SchoolID;
```

---

## QUERY-022: Join operation (full grammar)

```
join_operation:
  { cross_join_operation | condition_join_operation }

cross_join_operation:
  from_item cross_join_operator from_item

condition_join_operation:
  from_item condition_join_operator from_item join_condition

cross_join_operator:
  { CROSS JOIN | , }

condition_join_operator:
  {
    [INNER] JOIN
    | FULL [OUTER] JOIN
    | LEFT [OUTER] JOIN
    | RIGHT [OUTER] JOIN
  }

join_condition:
  { on_clause | using_clause }

on_clause:
  ON bool_expression

using_clause:
  USING ( column_list )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#join_operation
version: bigquery-current

```sql
SELECT * FROM a LEFT JOIN b ON a.id = b.id RIGHT JOIN c USING (id);
```

---

## QUERY-023: WHERE clause

```
WHERE bool_expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#where_clause
version: bigquery-current

```sql
SELECT * FROM Roster WHERE SchoolID > 52;
SELECT * FROM Roster WHERE name LIKE 'A%' AND age > 18;
```

---

## QUERY-024: GROUP BY clause (full grammar)

```
GROUP BY group_by_specification

group_by_specification:
  {
    groupable_items
    | ALL
    | grouping_sets_specification
    | rollup_specification
    | cube_specification
    | ()
  }
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#group_by_clause
version: bigquery-current

```sql
SELECT SUM(PointsScored) AS total_points, LastName FROM PlayerStats GROUP BY LastName;
SELECT SUM(PointsScored) AS total_points, LastName, FirstName FROM PlayerStats GROUP BY 2, 3;
```

---

## QUERY-025: GROUP BY ALL

```
GROUP BY ALL
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#group_by_all
version: bigquery-current

```sql
SELECT SUM(PointsScored) AS total_points, FirstName AS first_name, LastName AS last_name
FROM PlayerStats
GROUP BY ALL;
```

---

## QUERY-026: GROUP BY GROUPING SETS

```
GROUP BY GROUPING SETS ( grouping_list )

grouping_list:
  {
    rollup_specification
    | cube_specification
    | groupable_item
    | groupable_item_set
  }[, ...]

groupable_item_set:
  ( [ groupable_item[, ...] ] )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#group_by_grouping_sets
version: bigquery-current

```sql
SELECT product_type, product_name, SUM(product_count) AS product_sum
FROM Products
GROUP BY GROUPING SETS (product_type, product_name, ());
```

---

## QUERY-027: GROUP BY ROLLUP

```
GROUP BY ROLLUP ( grouping_list )

grouping_list:
  { groupable_item | groupable_item_set }[, ...]

groupable_item_set:
  ( groupable_item[, ...] )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#group_by_rollup
version: bigquery-current

```sql
SELECT product_type, product_name, SUM(product_count) AS product_sum
FROM Products
GROUP BY ROLLUP (product_type, product_name)
ORDER BY product_type, product_name;
```

---

## QUERY-028: GROUP BY CUBE

```
GROUP BY CUBE ( grouping_list )

grouping_list:
  { groupable_item | groupable_item_set }[, ...]

groupable_item_set:
  ( groupable_item[, ...] )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#group_by_cube
version: bigquery-current

```sql
SELECT product_type, product_name, SUM(product_count) AS product_sum
FROM Products
GROUP BY CUBE (product_type, product_name)
ORDER BY product_type, product_name;
```

---

## QUERY-029: HAVING clause

```
HAVING bool_expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#having_clause
version: bigquery-current

```sql
SELECT LastName FROM Roster GROUP BY LastName HAVING SUM(PointsScored) > 15;
SELECT LastName, SUM(PointsScored) AS ps FROM Roster GROUP BY LastName HAVING ps > 0;
```

---

## QUERY-030: ORDER BY clause

```
ORDER BY expression
  [{ ASC | DESC }]
  [{ NULLS FIRST | NULLS LAST }]
  [, ...]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#order_by_clause
version: bigquery-current

```sql
SELECT x, y FROM t ORDER BY x;
SELECT x, y FROM t ORDER BY x NULLS LAST;
SELECT x, y FROM t ORDER BY x DESC NULLS FIRST;
SELECT SUM(PointsScored), LastName FROM PlayerStats GROUP BY LastName ORDER BY 2;
```

---

## QUERY-031: QUALIFY clause

```
QUALIFY bool_expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#qualify_clause
version: bigquery-current

```sql
SELECT item, RANK() OVER (PARTITION BY category ORDER BY purchases DESC) AS rank
FROM Produce
WHERE Produce.category = 'vegetable'
QUALIFY rank <= 3;

SELECT item
FROM Produce
WHERE Produce.category = 'vegetable'
QUALIFY RANK() OVER (PARTITION BY category ORDER BY purchases DESC) <= 3;
```

---

## QUERY-032: WINDOW clause

```
WINDOW named_window_expression [, ...]

named_window_expression:
  named_window AS { named_window | ( [ window_specification ] ) }
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#window_clause
version: bigquery-current

```sql
SELECT item, purchases, category, LAST_VALUE(item)
  OVER (item_window) AS most_popular
FROM Produce
WINDOW item_window AS (
  PARTITION BY category
  ORDER BY purchases
  ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING);

SELECT item, purchases, category, LAST_VALUE(item) OVER (d) AS most_popular
FROM Produce
WINDOW
  a AS (PARTITION BY category),
  b AS (a ORDER BY purchases),
  c AS (b ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING),
  d AS (c);
```

---

## QUERY-033: Set operators (UNION / INTERSECT / EXCEPT)

```
  query_expr
  [ { INNER | [ { FULL | LEFT } [ OUTER ] ] } ]
  {
    UNION { ALL | DISTINCT } |
    INTERSECT DISTINCT |
    EXCEPT DISTINCT
  }
  [ { BY NAME [ ON (column_list) ] | [ STRICT ] CORRESPONDING [ BY (column_list) ] } ]
  query_expr
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#set_operators
version: bigquery-current

```sql
SELECT * FROM Roster UNION ALL SELECT * FROM TeamMascot ORDER BY SchoolID;
SELECT * FROM Roster UNION DISTINCT SELECT * FROM TeamMascot;
SELECT * FROM A INTERSECT DISTINCT SELECT * FROM B;
SELECT * FROM A EXCEPT DISTINCT SELECT * FROM B;

-- BY NAME (column-name matching)
SELECT 1 AS one_digit, 10 AS two_digit
UNION ALL BY NAME
SELECT 20 AS two_digit, 2 AS one_digit;
```

---

## QUERY-034: LIMIT and OFFSET clause

```
LIMIT count [ OFFSET skip_rows ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#limit_and_offset_clause
version: bigquery-current

```sql
SELECT * FROM UNNEST(ARRAY<STRING>['a','b','c','d','e']) AS letter
ORDER BY letter ASC LIMIT 2;

SELECT * FROM UNNEST(ARRAY<STRING>['a','b','c','d','e']) AS letter
ORDER BY letter ASC LIMIT 3 OFFSET 1;
```

---

## QUERY-035: WITH clause (non-recursive CTE)

```
WITH [ RECURSIVE ] { non_recursive_cte | recursive_cte }[, ...]

non_recursive_cte:
  cte_name AS ( query_expr )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#with_clause
version: bigquery-current

```sql
WITH subQ1 AS (SELECT SchoolID FROM Roster),
     subQ2 AS (SELECT OpponentID FROM PlayerStats)
SELECT * FROM subQ1
UNION ALL
SELECT * FROM subQ2;

WITH subQ1 AS (SELECT * FROM Roster WHERE SchoolID = 52),
     subQ2 AS (SELECT SchoolID FROM subQ1)
SELECT DISTINCT * FROM subQ2;
```

---

## QUERY-036: WITH RECURSIVE (recursive CTE)

```
WITH RECURSIVE { non_recursive_cte | recursive_cte }[, ...]

recursive_cte:
  cte_name AS ( recursive_union_operation )

recursive_union_operation:
  base_term union_operator recursive_term

base_term:
  query_expr

recursive_term:
  query_expr

union_operator:
  UNION ALL
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#recursive_keyword
version: bigquery-current

```sql
WITH RECURSIVE
  T1 AS ( (SELECT 1 AS n) UNION ALL (SELECT n + 1 AS n FROM T1 WHERE n < 3) )
SELECT n FROM T1;

WITH RECURSIVE
  T1 AS (
    (SELECT 1 AS n) UNION ALL
    (SELECT n + 2 FROM T1 WHERE n < 4))
SELECT * FROM T1 ORDER BY n;
```

---

## QUERY-037: Correlated subquery / correlated join

```
-- Correlated join: from_item references columns from outer query
SELECT * FROM T1 t1, t1.array_column;
SELECT * FROM T1 t1, t1.struct_column.array_field;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#correlated_join_operation
version: bigquery-current

```sql
SELECT (SELECT ARRAY_AGG(c) FROM t1.array_column c) FROM T1 t1;
SELECT a.struct_field1 FROM T1 t1, t1.array_of_structs a;
```

---

## QUERY-038: Field path in FROM (implicit UNNEST)

```
-- field_path is any path resolving to a field in a data type
FROM table_alias.array_field
FROM table_alias.struct_field.array_field
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#field_path
version: bigquery-current

```sql
SELECT * FROM T1 t1, t1.array_column;
SELECT * FROM T1 t1, t1.struct_column.array_field;
SELECT (SELECT STRING_AGG(a.struct_field1) FROM t1.array_of_structs a) FROM T1 t1;
```

---

## QUERY-039: AGGREGATION_THRESHOLD clause

```
SELECT WITH AGGREGATION_THRESHOLD
  OPTIONS (
    threshold = threshold_value
    [, privacy_unit_column = column_name]
  )
  select_list
FROM ...
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#aggregation_threshold_clause
version: bigquery-current

```sql
SELECT WITH AGGREGATION_THRESHOLD
  OPTIONS (threshold=50, privacy_unit_column=id)
  zip_code, COUNT(*) AS num_users
FROM dataset.users
GROUP BY zip_code;
```

---

## QUERY-040: Table function call in FROM

```
function_name ( [ argument [, ...] ] ) [ as_alias ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#table_function_calls
version: bigquery-current

```sql
SELECT * FROM my_tvf(arg1, arg2);
SELECT * FROM mydataset.my_tvf('param') AS t;
```
