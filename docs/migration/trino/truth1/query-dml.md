# Trino Query & DML Syntax Corpus — source: Trino 481 (current), https://trino.io/docs/current/

---

## select-full-syntax

- **syntax:**
```
[ WITH SESSION [ name = expression [, ...] ]
[ WITH [ FUNCTION udf ] [, ...] ]
[ WITH [ RECURSIVE ] with_query [, ...] ]
SELECT [ ALL | DISTINCT ] select_expression [, ...]
[ FROM from_item [, ...] ]
[ WHERE condition ]
[ GROUP BY [ ALL | DISTINCT ] grouping_element [, ...] ]
[ HAVING condition ]
[ WINDOW window_definition_list ]
[ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
[ ORDER BY expression [ ASC | DESC ] [, ...] ]
[ OFFSET count [ ROW | ROWS ] ]
[ LIMIT { count | ALL } ]
[ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } { ONLY | WITH TIES } ]
```

  where `from_item` is one of:
```
table_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]

from_item join_type from_item
  [ ON join_condition | USING ( join_column [, ...] ) ]

table_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
  MATCH_RECOGNIZE pattern_recognition_specification
    [ [ AS ] alias [ ( column_alias [, ...] ) ] ]

TABLE (table_function_invocation) [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
```

  where `join_type` is one of:
```
[ INNER ] JOIN
LEFT [ OUTER ] JOIN
RIGHT [ OUTER ] JOIN
FULL [ OUTER ] JOIN
CROSS JOIN
```

  where `grouping_element` is one of:
```
()
expression
AUTO
GROUPING SETS ( ( column [, ...] ) [, ...] )
CUBE ( column [, ...] )
ROLLUP ( column [, ...] )
```

  where `select_expression` is one of:
```
expression [ [ AS ] column_alias ]
row_expression.* [ AS ( column_alias [, ...] ) ]
relation.*
*
```

- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT * FROM orders WHERE totalprice > 1000 ORDER BY orderdate DESC LIMIT 10;
```
- **notes:** The WITH SESSION, WITH FUNCTION, and WITH RECURSIVE clauses precede the SELECT keyword. Multiple from_items separated by commas produce a cross join. The WINDOW clause defines named windows; set operations (UNION/INTERSECT/EXCEPT) appear between the WINDOW clause and ORDER BY.

---

## select-clause

- **syntax:**
```
SELECT [ ALL | DISTINCT ] select_expression [, ...]

select_expression is one of:
    expression [ [ AS ] column_alias ]
    row_expression.* [ AS ( column_alias [, ...] ) ]
    relation.*
    *
```
- **source_url:** https://trino.io/docs/current/sql/select.html#select-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- ROW type wildcard expansion with aliases
SELECT (CAST(ROW(1, true) AS ROW(field1 bigint, field2 boolean))).* AS (alias1, alias2);
```
- **notes:** `SELECT ALL` is the default (returns all rows including duplicates). `SELECT DISTINCT` removes duplicates. The `row_expression.*` form expands a ROW-typed expression into individual columns and supports renaming via `AS (col_alias [, ...])`. `relation.*` expands all columns of a named relation.

---

## from-clause

- **syntax:**
```
FROM from_item [, ...]

from_item is one of:
    table_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    from_item join_type from_item [ ON join_condition | USING ( join_column [, ...] ) ]
    table_name MATCH_RECOGNIZE ( ... ) [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    TABLE (table_function_invocation) [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    UNNEST ( array_or_map [, ...] ) [ WITH ORDINALITY ] [ AS alias ( column_alias [, ...] ) ]
    LATERAL ( subquery ) [ AS alias ]
    table_name TABLESAMPLE { BERNOULLI | SYSTEM } ( percentage )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#from-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT n.name, r.name FROM nation AS n CROSS JOIN region AS r;
```
- **notes:** Multiple from_items separated by commas are equivalent to CROSS JOIN. UNNEST, LATERAL, TABLESAMPLE, JSON_TABLE, and NEAREST are additional FROM constructs described in dedicated sections.

---

## join-cross

- **syntax:**
```
from_item CROSS JOIN from_item
```
- **source_url:** https://trino.io/docs/current/sql/select.html#cross-join
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT *
FROM nation
CROSS JOIN region;

SELECT n.name AS nation, r.name AS region
FROM nation AS n
CROSS JOIN region AS r
ORDER BY 1, 2;

-- Equivalent comma syntax
SELECT *
FROM nation, region;
```
- **notes:** Produces the Cartesian product of the two relations. Comma-separated FROM items are syntactic sugar for CROSS JOIN.

---

## join-inner

- **syntax:**
```
from_item [ INNER ] JOIN from_item { ON join_condition | USING ( join_column [, ...] ) }
```
- **source_url:** https://trino.io/docs/current/sql/select.html#join-types
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT o.*, l.*
FROM orders o
JOIN lineitem l ON o.orderkey = l.orderkey;

SELECT *
FROM orders
JOIN lineitem USING (orderkey);
```
- **notes:** The INNER keyword is optional. Both ON and USING conditions are supported.

---

## join-left-right-full-outer

- **syntax:**
```
from_item LEFT  [ OUTER ] JOIN from_item { ON join_condition | USING ( join_column [, ...] ) }
from_item RIGHT [ OUTER ] JOIN from_item { ON join_condition | USING ( join_column [, ...] ) }
from_item FULL  [ OUTER ] JOIN from_item { ON join_condition | USING ( join_column [, ...] ) }
```
- **source_url:** https://trino.io/docs/current/sql/select.html#join-types
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT o.orderkey, l.partkey
FROM orders o
LEFT JOIN lineitem l ON o.orderkey = l.orderkey;
```
- **notes:** The OUTER keyword is optional for all three forms. LEFT/RIGHT retain all rows from the respective side, padding with NULLs on the other side. FULL OUTER retains all rows from both sides.

---

## join-lateral

- **syntax:**
```
from_item CROSS JOIN LATERAL ( subquery ) [ [ AS ] alias ]
from_item LEFT  JOIN LATERAL ( subquery ) [ [ AS ] alias ] ON TRUE
from_item FULL  JOIN LATERAL ( subquery ) [ [ AS ] alias ] ON TRUE
```
- **source_url:** https://trino.io/docs/current/sql/select.html#lateral
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name, x, y
FROM nation
CROSS JOIN LATERAL (SELECT name || ' :-' AS x)
CROSS JOIN LATERAL (SELECT x || ')' AS y);
```
- **notes:** LATERAL allows a subquery in the FROM clause to reference columns from FROM items that precede it. When LATERAL appears on the right side of a FULL JOIN, the only condition supported is ON TRUE.

---

## join-nearest

- **syntax:**
```
from_item CROSS JOIN NEAREST (
    FROM relation
    [ WHERE condition ]
    MATCH comparison
)

from_item INNER JOIN NEAREST (
    FROM relation
    [ WHERE condition ]
    MATCH comparison
) ON TRUE

from_item LEFT JOIN NEAREST (
    FROM relation
    [ WHERE condition ]
    MATCH comparison
) ON TRUE
```

  where `comparison` uses one of the operators: `<`, `<=`, `>`, `>=`

- **source_url:** https://trino.io/docs/current/sql/select.html#nearest
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT trades.symbol, trades.ts, quotes.price
FROM trades
CROSS JOIN NEAREST (
    FROM quotes
    WHERE quotes.symbol = trades.symbol
    MATCH quotes.ts <= trades.ts
);

SELECT trades.symbol, quotes.price
FROM trades
LEFT JOIN NEAREST (
    FROM quotes
    WHERE quotes.symbol = trades.symbol
    MATCH quotes.ts >= trades.ts
) ON TRUE;
```
- **notes:** NEAREST selects at most one row from the inner relation per left-side row. The MATCH clause is required and uses a directional comparison: `<`/`<=` select the closest row whose match key is smaller (or equal); `>`/`>=` select the closest row whose match key is larger (or equal). The matching direction also determines ordering of candidate rows.

---

## unnest

- **syntax:**
```
UNNEST ( array_or_map_expression [, ...] ) [ WITH ORDINALITY ]
    [ AS table_alias ( column_alias [, ...] ) ]
```
- **source_url:** https://trino.io/docs/current/sql/select.html#unnest
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Expand a single array
SELECT * FROM UNNEST(ARRAY[1,2]) AS t(number);

-- Expand a map into key-value columns
SELECT * FROM UNNEST(
    map_from_entries(ARRAY[('SQL',1974),('Java', 1995)])
) AS t(language, first_appeared_year);

-- WITH ORDINALITY adds a row-number column
SELECT a, b, rownumber
FROM UNNEST(ARRAY[2, 5], ARRAY[7, 8, 9]) WITH ORDINALITY AS t(a, b, rownumber);

-- CROSS JOIN UNNEST to expand a column
SELECT student, score
FROM (
   VALUES
      ('John', ARRAY[7, 10, 9]),
      ('Mary', ARRAY[4, 8, 9])
) AS tests (student, scores)
CROSS JOIN UNNEST(scores) AS t(score);

-- Multiple arrays: null-padded to longest
SELECT numbers, animals, n, a
FROM (
  VALUES
    (ARRAY[2, 5], ARRAY['dog', 'cat', 'bird']),
    (ARRAY[7, 8, 9], ARRAY['cow', 'pig'])
) AS x (numbers, animals)
CROSS JOIN UNNEST(numbers, animals) AS t (n, a);

-- LEFT JOIN to preserve rows with empty/null arrays
SELECT runner, checkpoint
FROM (
   VALUES
      ('Joe', ARRAY[10, 20, 30, 42]),
      ('Roger', ARRAY[10]),
      ('Dave', ARRAY[]),
      ('Levi', NULL)
) AS marathon (runner, checkpoints)
LEFT JOIN UNNEST(checkpoints) AS t(checkpoint) ON TRUE;

-- ROW type expansion
SELECT *
FROM UNNEST(
    ARRAY[ROW('Java', 1995), ROW('SQL', 1974)],
    ARRAY[ROW(false), ROW(true)]
) AS t(language, first_appeared_year, declarative);
```
- **notes:** UNNEST can expand ARRAY (one output column) or MAP (two output columns: key, value). Multiple array arguments are expanded side-by-side with NULL padding for shorter arrays. WITH ORDINALITY appends a row-number column. LEFT JOIN UNNEST ON TRUE preserves source rows with empty arrays or NULL arrays. The only condition supported with LEFT JOIN UNNEST is ON TRUE.

---

## tablesample

- **syntax:**
```
table_name [ [ AS ] alias ] TABLESAMPLE { BERNOULLI | SYSTEM } ( percentage )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#tablesample
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT *
FROM users TABLESAMPLE BERNOULLI (50);

SELECT *
FROM users TABLESAMPLE SYSTEM (75);

SELECT o.*, i.*
FROM orders o TABLESAMPLE SYSTEM (10)
JOIN lineitem i TABLESAMPLE BERNOULLI (40)
  ON o.orderkey = i.orderkey;
```
- **notes:** BERNOULLI samples each row independently with the given probability percentage (0–100). SYSTEM divides the table into logical segments and samples at segment granularity. Neither method guarantees an exact row count. Percentage is a numeric literal from 0 to 100. Both methods can be combined in a single query across multiple tables.

---

## table-function-invocation

- **syntax:**
```
TABLE ( table_function_invocation ) [ [ AS ] alias [ ( column_alias [, ...] ) ] ]

table_function_invocation:
    [ catalog_name. ] [ schema_name. ] function_name ( [ argument [, ...] ] )

argument:
    scalar_argument
  | descriptor_argument
  | table_argument

scalar_argument:
    [ argument_name => ] expression

descriptor_argument:
    [ argument_name => ] DESCRIPTOR ( field_name [ data_type ] [, ...] )
  | [ argument_name => ] CAST ( NULL AS DESCRIPTOR )

table_argument:
    [ argument_name => ] TABLE ( { table_name | query } )
        [ PARTITION BY { column [, ...] | { HASH | RANGE } ON ( column [, ...] ) } ]
        [ { KEEP | PRUNE } WHEN EMPTY ]
        [ ORDER BY column [, ...] ]
```
- **source_url:** https://trino.io/docs/current/functions/table.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Built-in: exclude_columns
SELECT *
FROM TABLE(exclude_columns(
    input => TABLE(orders),
    columns => DESCRIPTOR(clerk, comment)));

-- Built-in: sequence
SELECT *
FROM TABLE(sequence(
    start => 1000000,
    stop => -2000000,
    step => -3));

-- Named arguments (any order)
SELECT * FROM TABLE(my_function(row_count => 100, column_count => 1));

-- Positional arguments
SELECT * FROM TABLE(my_function(1, 100));

-- Schema-qualified
SELECT * FROM TABLE(schema_name.my_function(1, 100));
SELECT * FROM TABLE(catalog_name.schema_name.my_function(1, 100));

-- With parameters
PREPARE stmt FROM SELECT * FROM TABLE(my_function(row_count => ? + 1, column_count => ?));
EXECUTE stmt USING 100, 1;

-- Table argument with partitioning and ordering
SELECT * FROM TABLE(my_function(
    input => TABLE(orders)
        PARTITION BY orderstatus
        KEEP WHEN EMPTY
        ORDER BY orderdate
));
```
- **notes:** Arguments may be passed by name (named => value) or positionally. DESCRIPTOR arguments describe schemas with optional type information; CAST(NULL AS DESCRIPTOR) passes an empty descriptor. Table arguments support PARTITION BY (with optional HASH or RANGE strategy), KEEP/PRUNE WHEN EMPTY (controls whether the function executes when the table argument is empty), and ORDER BY. The COPARTITION clause (for co-partitioning multiple table arguments) is referenced in the developer API but is not explicitly documented in the user-facing SQL reference for Trino 481.

---

## where-clause

- **syntax:**
```
WHERE condition
```
- **source_url:** https://trino.io/docs/current/sql/select.html#where-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name FROM nation WHERE regionkey = 1;

SELECT name FROM nation
WHERE EXISTS (
     SELECT *
     FROM region
     WHERE region.regionkey = nation.regionkey
);

SELECT name FROM nation
WHERE regionkey IN (
     SELECT regionkey
     FROM region
     WHERE name = 'AMERICA' OR name = 'AFRICA'
);

SELECT name FROM nation
WHERE regionkey = (SELECT max(regionkey) FROM region);
```
- **notes:** The WHERE clause filters rows before grouping. It supports all standard predicates including EXISTS, IN, ANY/ALL quantified comparisons, and scalar subqueries.

---

## subquery-exists

- **syntax:**
```
EXISTS ( subquery )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#exists
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name
FROM nation
WHERE EXISTS (
     SELECT *
     FROM region
     WHERE region.regionkey = nation.regionkey
);
```
- **notes:** Returns TRUE if the subquery produces at least one row. Correlated EXISTS subqueries referencing outer query columns are supported.

---

## subquery-in

- **syntax:**
```
expression IN ( subquery )
expression NOT IN ( subquery )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#in-predicate
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name
FROM nation
WHERE regionkey IN (
     SELECT regionkey
     FROM region
     WHERE name = 'AMERICA' OR name = 'AFRICA'
);
```
- **notes:** The subquery must produce exactly one column. Tests whether the expression equals any value returned by the subquery.

---

## subquery-scalar

- **syntax:**
```
( subquery )   -- used as a scalar expression
```
- **source_url:** https://trino.io/docs/current/sql/select.html#scalar-subquery
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name
FROM nation
WHERE regionkey = (SELECT max(regionkey) FROM region);
```
- **notes:** A scalar subquery returns zero or one row. If it returns zero rows, the result is NULL. If it returns more than one row, a runtime error is raised. Currently only a single column may be returned from a scalar subquery.

---

## group-by

- **syntax:**
```
GROUP BY [ ALL | DISTINCT ] grouping_element [, ...]

grouping_element is one of:
    ()
    expression
    AUTO
    GROUPING SETS ( ( column [, ...] ) [, ...] )
    CUBE ( column [, ...] )
    ROLLUP ( column [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#group-by-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT count(*), mktsegment, nationkey,
       CAST(sum(acctbal) AS bigint) AS totalbal
FROM customer
GROUP BY mktsegment, nationkey
HAVING sum(acctbal) > 5700000
ORDER BY totalbal DESC;
```
- **notes:** The ALL and DISTINCT quantifiers control deduplication when multiple complex grouping elements (CUBE, ROLLUP, GROUPING SETS) are combined. ALL (default) produces all combinations including duplicates. DISTINCT eliminates redundant grouping sets. Grouping expressions may reference column names or expressions from the SELECT list.

---

## group-by-auto

- **syntax:**
```
GROUP BY AUTO
```
- **source_url:** https://trino.io/docs/current/sql/select.html#group-by-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT mktsegment, sum(acctbal) FROM shipping GROUP BY AUTO;
```
- **notes:** The AUTO keyword instructs Trino to automatically determine grouping columns. Any column in the SELECT list that is not part of an aggregate function is implicitly treated as a grouping column. Equivalent to listing all non-aggregate SELECT columns explicitly in GROUP BY.

---

## group-by-grouping-sets

- **syntax:**
```
GROUP BY GROUPING SETS ( ( column [, ...] ) [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#grouping-sets
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT origin_state, origin_zip, destination_state, sum(package_weight)
FROM shipping
GROUP BY GROUPING SETS (
    (origin_state),
    (origin_state, origin_zip),
    (destination_state));
```
- **notes:** Equivalent to a UNION ALL of multiple GROUP BY queries, one per listed grouping set. Columns not in a given grouping set are set to NULL for that set's rows.

---

## group-by-cube

- **syntax:**
```
GROUP BY CUBE ( column [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#cube
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT origin_state, destination_state, sum(package_weight)
FROM shipping
GROUP BY CUBE (origin_state, destination_state);
```
- **notes:** CUBE generates all possible grouping sets (power set) for the given columns. For N columns, generates 2^N grouping sets. Equivalent to GROUPING SETS of all subsets including the empty set ().

---

## group-by-rollup

- **syntax:**
```
GROUP BY ROLLUP ( column [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#rollup
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT origin_state, origin_zip, sum(package_weight)
FROM shipping
GROUP BY ROLLUP (origin_state, origin_zip);
```
- **notes:** ROLLUP generates hierarchical subtotals. For columns (c1, c2, ..., cN) it generates N+1 grouping sets: (c1, c2, ..., cN), (c1, c2, ..., c(N-1)), ..., (c1), (). Useful for computing totals and subtotals in a hierarchy.

---

## group-by-all-distinct

- **syntax:**
```
GROUP BY ALL  grouping_element [, ...]
GROUP BY DISTINCT grouping_element [, ...]
```
- **source_url:** https://trino.io/docs/current/sql/select.html#grouping-sets
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT origin_state, destination_state, origin_zip, sum(package_weight)
FROM shipping
GROUP BY ALL
    CUBE (origin_state, destination_state),
    ROLLUP (origin_state, origin_zip);

SELECT origin_state, destination_state, origin_zip, sum(package_weight)
FROM shipping
GROUP BY DISTINCT
    CUBE (origin_state, destination_state),
    ROLLUP (origin_state, origin_zip);
```
- **notes:** When multiple complex grouping elements are combined, ALL (default) takes the cross-product and retains all combinations including duplicates; DISTINCT deduplicates the resulting grouping sets. These quantifiers apply to the overall GROUP BY, not to individual grouping elements.

---

## grouping-function

- **syntax:**
```
grouping ( column [, ...] ) -> bigint
```
- **source_url:** https://trino.io/docs/current/sql/select.html#grouping-operation
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT origin_state, origin_zip, destination_state, sum(package_weight),
       grouping(origin_state, origin_zip, destination_state)
FROM shipping
GROUP BY GROUPING SETS (
    (origin_state),
    (origin_state, origin_zip),
    (destination_state)
);
```
- **notes:** Returns a bitmask as a bigint indicating which columns are present in the current grouping. Bits are assigned right-to-left (rightmost argument = least-significant bit). Bit value 0 means the column IS in the grouping; 1 means it is NOT. Must be used with GROUPING SETS, ROLLUP, CUBE, or GROUP BY.

---

## having-clause

- **syntax:**
```
HAVING condition
```
- **source_url:** https://trino.io/docs/current/sql/select.html#having-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT count(*), mktsegment, nationkey,
       CAST(sum(acctbal) AS bigint) AS totalbal
FROM customer
GROUP BY mktsegment, nationkey
HAVING sum(acctbal) > 5700000
ORDER BY totalbal DESC;
```
- **notes:** HAVING filters groups after aggregation, in conjunction with GROUP BY and aggregate functions. The HAVING condition is evaluated after grouping but before projection in the SELECT list.

---

## window-clause

- **syntax:**
```
WINDOW window_definition_list

window_definition_list:
    window_name AS ( window_specification ) [, ...]

window_specification:
    [ existing_window_name ]
    [ PARTITION BY expression [, ...] ]
    [ ORDER BY expression [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...] ]
    [ { ROWS | RANGE | GROUPS }
        { UNBOUNDED PRECEDING
        | expression PRECEDING
        | CURRENT ROW
        | BETWEEN frame_bound AND frame_bound }
    ]

frame_bound:
    UNBOUNDED PRECEDING
  | expression PRECEDING
  | CURRENT ROW
  | expression FOLLOWING
  | UNBOUNDED FOLLOWING
```
- **source_url:** https://trino.io/docs/current/sql/select.html#window-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT orderkey, clerk, totalprice,
      rank() OVER w AS rnk
FROM orders
WINDOW w AS (PARTITION BY clerk ORDER BY totalprice DESC)
ORDER BY count() OVER w, clerk, rnk;

-- Aggregate as window function (rolling sum)
-- sum(totalprice) OVER (PARTITION BY clerk ORDER BY orderdate) AS rolling_sum
```
- **notes:** Named windows defined in WINDOW can be referenced in the SELECT list and ORDER BY via `OVER window_name`. Window specifications may inherit from another named window (by including the existing window name as the first element). All components (partition, ordering, frame) are optional. Default frame when ORDER BY is present: `RANGE UNBOUNDED PRECEDING` (equivalent to `RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW`). Without ORDER BY, default frame is all rows. Row pattern recognition clauses are supported within window specifications. Aggregate functions can be used as window functions by adding an OVER clause.

---

## union

- **syntax:**
```
query UNION [ ALL | DISTINCT ] [ CORRESPONDING ] query
```
- **source_url:** https://trino.io/docs/current/sql/select.html#union-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT 13
UNION
SELECT 42;

-- UNION (deduplicates): returns 13, 42
SELECT 13
UNION
SELECT * FROM (VALUES 42, 13);

-- UNION ALL (preserves duplicates): returns 13, 42, 13
SELECT 13
UNION ALL
SELECT * FROM (VALUES 42, 13);

-- CORRESPONDING: match columns by name instead of position
SELECT * FROM (VALUES (1, 'alice')) AS t(id, name)
UNION ALL CORRESPONDING
SELECT * FROM (VALUES ('bob', 2)) AS t(name, id);
```
- **notes:** Default (without ALL or DISTINCT) is DISTINCT, which deduplicates. CORRESPONDING matches output columns by name rather than ordinal position, allowing queries with different column orderings to be combined. Multiple set operations are processed left to right; INTERSECT binds more tightly than UNION and EXCEPT.

---

## intersect

- **syntax:**
```
query INTERSECT [ ALL | DISTINCT ] [ CORRESPONDING ] query
```
- **source_url:** https://trino.io/docs/current/sql/select.html#intersect-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT * FROM (VALUES 13, 42)
INTERSECT
SELECT 13;

-- CORRESPONDING variant
SELECT * FROM (VALUES (1, 'alice')) AS t(id, name)
INTERSECT CORRESPONDING
SELECT * FROM (VALUES ('alice', 1)) AS t(name, id);
```
- **notes:** Returns only rows present in both result sets. Default is DISTINCT. INTERSECT binds more tightly than UNION and EXCEPT when used without parentheses.

---

## except

- **syntax:**
```
query EXCEPT [ ALL | DISTINCT ] [ CORRESPONDING ] query
```
- **source_url:** https://trino.io/docs/current/sql/select.html#except-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT * FROM (VALUES 13, 42)
EXCEPT
SELECT 13;

-- CORRESPONDING variant
SELECT * FROM (VALUES (1, 'alice'), (2, 'bob')) AS t(id, name)
EXCEPT CORRESPONDING
SELECT * FROM (VALUES ('alice', 1)) AS t(name, id);
```
- **notes:** Returns rows from the first query that do not appear in the second query. Default is DISTINCT. CORRESPONDING matches columns by name.

---

## order-by

- **syntax:**
```
ORDER BY expression [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...]
```
- **source_url:** https://trino.io/docs/current/sql/select.html#order-by-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name FROM nation ORDER BY name OFFSET 22;

SELECT count(*), mktsegment, nationkey,
       CAST(sum(acctbal) AS bigint) AS totalbal
FROM customer
GROUP BY mktsegment, nationkey
HAVING sum(acctbal) > 5700000
ORDER BY totalbal DESC;
```
- **notes:** Default sort direction is ASC. Default null ordering is NULLS LAST regardless of direction. Expressions may be output column aliases, ordinal positions (starting at 1), or arbitrary expressions. Per SQL specification, an ORDER BY clause only affects ordering for the immediately containing query; Trino drops redundant ORDER BY in subqueries or INSERT to avoid performance overhead.

---

## offset

- **syntax:**
```
OFFSET count [ ROW | ROWS ]
```
- **source_url:** https://trino.io/docs/current/sql/select.html#offset-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT name FROM nation ORDER BY name OFFSET 22;

SELECT * FROM (VALUES 5, 2, 4, 1, 3) t(x) ORDER BY x OFFSET 2 LIMIT 2;
```
- **notes:** Discards the first `count` rows from the result. If ORDER BY is present, OFFSET applies to the sorted result and the remaining set stays sorted. ROW and ROWS are interchangeable noise words.

---

## limit

- **syntax:**
```
LIMIT { count | ALL }
```
- **source_url:** https://trino.io/docs/current/sql/select.html#limit-or-fetch-first-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
SELECT orderdate FROM orders LIMIT 5;

SELECT * FROM (VALUES 5, 2, 4, 1, 3) t(x) ORDER BY x OFFSET 2 LIMIT 2;
```
- **notes:** `LIMIT ALL` returns all rows (equivalent to omitting LIMIT). LIMIT is processed after ORDER BY and OFFSET.

---

## fetch-first

- **syntax:**
```
FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } { ONLY | WITH TIES }
```
- **source_url:** https://trino.io/docs/current/sql/select.html#limit-or-fetch-first-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Default count of 1 (FETCH FIRST ROW ONLY)
SELECT orderdate FROM orders FETCH FIRST ROW ONLY;

-- WITH TIES includes all rows tied with the last included row
SELECT name, regionkey
FROM nation
ORDER BY regionkey FETCH FIRST ROW WITH TIES;
```
- **notes:** FIRST and NEXT are interchangeable. ROW and ROWS are interchangeable. If count is omitted, it defaults to 1. WITH TIES requires ORDER BY to be present; it returns the specified number of leading rows plus any additional rows that tie with the last row according to the ORDER BY. ONLY and WITH TIES are mutually exclusive.

---

## with-cte

- **syntax:**
```
WITH [ RECURSIVE ] with_query [, ...]
SELECT ...

with_query:
    query_name [ ( column_name [, ...] ) ] AS ( query )
```
- **source_url:** https://trino.io/docs/current/sql/select.html#with-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Single CTE
WITH x AS (SELECT a, MAX(b) AS b FROM t GROUP BY a)
SELECT a, b FROM x;

-- Multiple CTEs
WITH
  t1 AS (SELECT a, MAX(b) AS b FROM x GROUP BY a),
  t2 AS (SELECT a, AVG(d) AS d FROM y GROUP BY a)
SELECT t1.*, t2.*
FROM t1
JOIN t2 ON t1.a = t2.a;

-- Chained CTEs (each can reference prior CTEs)
WITH
  x AS (SELECT a FROM t),
  y AS (SELECT a AS b FROM x),
  z AS (SELECT b AS c FROM y)
SELECT c FROM z;
```
- **notes:** CTEs define named relations available within the query scope. Currently, the SQL for a CTE is inlined wherever the named relation is used; if referenced multiple times in a non-deterministic query, results may differ across uses. Forward references (a CTE referencing a later CTE) are not allowed.

---

## with-recursive

- **syntax:**
```
WITH RECURSIVE with_query [, ...]
SELECT ...

-- Recursive with_query shape (required):
query_name ( column_name [, ...] ) AS (
    non_recursive_base_query
    UNION ALL
    recursive_step_query_referencing_query_name
)
```
- **source_url:** https://trino.io/docs/current/sql/select.html#with-recursive-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
WITH RECURSIVE t(n) AS (
    VALUES (1)
    UNION ALL
    SELECT n + 1 FROM t WHERE n < 4
)
SELECT sum(n) FROM t;
-- Result: 10 (1 + 2 + 3 + 4)
```
- **notes:** Experimental feature. The recursive query must be shaped as UNION of exactly two parts: a non-recursive base and a recursive step. Column aliases are mandatory for all recursive WITH queries. Types of columns returned by the step relation must be coercible to the base relation types. Outer joins, non-UNION set operations, LIMIT, and certain other clauses are restricted in the step relation. Default recursion depth is 10, controlled by session property `max_recursion_depth`. Query plan size grows quadratically with recursion depth.

---

## with-session

- **syntax:**
```
WITH SESSION name = expression [, ...]
SELECT ...
```
- **source_url:** https://trino.io/docs/current/sql/select.html#with-session-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
WITH
  SESSION
    query_max_execution_time='2h',
    example.query_partition_filter_required=true
SELECT *
FROM example.default.thetable
LIMIT 100;
```
- **notes:** WITH SESSION sets session and catalog session properties for the duration of the current SELECT statement only. Values override any other configuration and session property settings. Catalog-qualified properties use the form `catalog.property_name = value`.

---

## with-function

- **syntax:**
```
WITH FUNCTION function_name ( [ parameter_name data_type [, ...] ] )
    RETURNS return_type
    RETURN expression
[, FUNCTION ... ]
SELECT ...
```
- **source_url:** https://trino.io/docs/current/sql/select.html#with-function-clause
- **version:** Trino 481 (current)
- **example_sql:**
```sql
WITH 
  FUNCTION hello(name varchar)
    RETURNS varchar
    RETURN format('Hello %s!', name),
  FUNCTION bye(name varchar)
    RETURNS varchar
    RETURN format('Bye %s!', name)
SELECT hello('Finn') || ' and ' || bye('Joe');
-- Result: Hello Finn! and Bye Joe!
```
- **notes:** Inline user-defined functions defined with WITH FUNCTION are available only within the current SELECT statement. Multiple functions are separated by commas after the FUNCTION keyword list.

---

## values

- **syntax:**
```
VALUES row [, ...]

row:
    expression
  | ( column_expression [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/values.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Single column, three rows
VALUES 1, 2, 3

-- Two columns, three rows
VALUES (1, 'a'), (2, 'b'), (3, 'c')

-- With column aliases via AS
SELECT * FROM (VALUES (1, 'a'), (2, 'b'), (3, 'c')) AS t (id, name)

-- In CREATE TABLE AS
CREATE TABLE example AS
SELECT * FROM (VALUES (1, 'a'), (2, 'b'), (3, 'c')) AS t (id, name)
```
- **notes:** VALUES defines a literal inline table. It can be used anywhere a query can appear: in FROM clause, as the source for INSERT, as a standalone statement, or nested inside CTEs. Creates an anonymous table; columns can be named with an AS clause and column alias list.

---

## match-recognize

- **syntax:**
```
table_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
MATCH_RECOGNIZE (
  [ PARTITION BY column [, ...] ]
  [ ORDER BY column [, ...] ]
  [ MEASURES measure_definition [, ...] ]
  [ rows_per_match ]
  [ AFTER MATCH skip_to ]
  PATTERN ( row_pattern )
  [ SUBSET subset_definition [, ...] ]
  DEFINE variable_definition [, ...]
)
[ [ AS ] alias [ ( column_alias [, ...] ) ] ]

rows_per_match:
    ONE ROW PER MATCH
  | ALL ROWS PER MATCH
  | ALL ROWS PER MATCH SHOW EMPTY MATCHES
  | ALL ROWS PER MATCH OMIT EMPTY MATCHES
  | ALL ROWS PER MATCH WITH UNMATCHED ROWS

skip_to:
    AFTER MATCH SKIP PAST LAST ROW
  | AFTER MATCH SKIP TO NEXT ROW
  | AFTER MATCH SKIP TO [ FIRST | LAST ] pattern_variable

row_pattern (elements):
    A B+ C+ D+                    -- concatenation
  | A | B | C                     -- alternation (one component matches)
  | PERMUTE(A, B, C)              -- permutation (all in any order)
  | ( row_pattern )               -- grouping
  | ^                             -- anchor to partition start
  | $                             -- anchor to partition end
  | ()                            -- empty pattern
  | {- row_pattern -}             -- exclude from output

quantifiers (append to pattern element):
    *       -- zero or more (greedy)
    +       -- one or more (greedy)
    ?       -- zero or one (greedy)
    {n}     -- exactly n
    {m, n}  -- between m and n
    {, n}   -- up to n (0 to n)
    {m, }   -- at least m
    {,}     -- zero or more (same as *)
    -- append ? to any quantifier for reluctant (non-greedy) matching

measure_definition:
    expression AS measure_name

subset_definition:
    subset_variable = ( primary_variable [, ...] )

variable_definition:
    primary_variable AS boolean_expression
```
- **source_url:** https://trino.io/docs/current/sql/match-recognize.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- V-shape price pattern detection
SELECT * FROM orders MATCH_RECOGNIZE(
  PARTITION BY custkey
  ORDER BY orderdate
  MEASURES
    A.totalprice AS starting_price,
    LAST(B.totalprice) AS bottom_price,
    LAST(U.totalprice) AS top_price
  ONE ROW PER MATCH
  AFTER MATCH SKIP PAST LAST ROW
  PATTERN (A B+ C+ D+)
  SUBSET U = (C, D)
  DEFINE
    B AS totalprice < PREV(totalprice),
    C AS totalprice > PREV(totalprice) AND totalprice <= A.totalprice,
    D AS totalprice > PREV(totalprice)
)
```
- **notes:** PARTITION BY divides the input into independent sections processed separately (analogous to window function partitioning). ORDER BY determines matching order (e.g. by date for temporal patterns). MEASURES defines scalar output columns; use FINAL before navigation functions to evaluate with complete match visibility (MEASURES only; RUNNING is default for both DEFINE and MEASURES). ONE ROW PER MATCH (default) outputs one row per match; ALL ROWS PER MATCH outputs one row per matched row with optional handling of empty matches and unmatched rows. SUBSET defines union variables grouping multiple primary variables for convenience. DEFINE assigns boolean conditions to primary pattern variables; variables without explicit definitions implicitly match all rows. Navigation functions: FIRST(expr, offset)/LAST(expr, offset) navigate logically from first/last occurrence of a variable; PREV(expr, offset)/NEXT(expr, offset) navigate physically (default offset=1 for PREV/NEXT, 0 for FIRST/LAST). CLASSIFIER() returns the matched variable name as varchar; MATCH_NUMBER() returns the sequential match number as bigint. Standard aggregate functions can be applied over match rows using variable prefix syntax (e.g. count(A.*), avg(B.price)).

---

## insert

- **syntax:**
```
INSERT INTO table_name [ @ branch_name ] [ ( column [, ...] ) ] query
```

  where `query` is a SELECT statement or VALUES clause.

- **source_url:** https://trino.io/docs/current/sql/insert.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Insert from another table
INSERT INTO orders
SELECT * FROM new_orders;

-- Insert a single row with VALUES
INSERT INTO cities VALUES (1, 'San Francisco');

-- Insert multiple rows
INSERT INTO cities VALUES (2, 'San Jose'), (3, 'Oakland');

-- With explicit column list
INSERT INTO nation (nationkey, name, regionkey, comment)
VALUES (26, 'POLAND', 3, 'no comment');

-- Omitting a column (receives NULL)
INSERT INTO nation (nationkey, name, regionkey)
VALUES (26, 'POLAND', 3);

-- Insert into a branched table
INSERT INTO cities @ audit VALUES (1, 'San Francisco');
```
- **notes:** If a column list is specified, it must exactly match the columns produced by the query. Unlisted columns in the target table receive NULL. If no column list is provided, the query columns must exactly match the table columns in number and type.

---

## delete

- **syntax:**
```
DELETE FROM table_name [ @ branch_name ] [ WHERE condition ]
```
- **source_url:** https://trino.io/docs/current/sql/delete.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Delete with a simple condition
DELETE FROM lineitem WHERE shipmode = 'AIR';

-- Delete with a subquery condition
DELETE FROM lineitem
WHERE orderkey IN (SELECT orderkey FROM orders WHERE priority = 'LOW');

-- Delete all rows (no WHERE)
DELETE FROM orders;

-- Delete from a specific branch
DELETE FROM orders @ audit;
```
- **notes:** Without WHERE, all rows are deleted. Connector support for DELETE varies; see connector-specific documentation.

---

## update

- **syntax:**
```
UPDATE table_name [ @ branch_name ]
SET column = expression [, ...]
[ WHERE condition ]
```
- **source_url:** https://trino.io/docs/current/sql/update.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Basic conditional update
UPDATE purchases
SET status = 'OVERDUE'
WHERE ship_date IS NULL;

-- Multiple column assignment
UPDATE customers
SET account_manager = 'John Henry', assign_date = now();

-- Subquery-based update
UPDATE new_hires
SET manager = (
  SELECT e.name
  FROM employees e
  WHERE e.employee_id = new_hires.manager_id
);

-- Branch-specific update
UPDATE purchases @ audit
SET status = 'OVERDUE'
WHERE ship_date IS NULL;
```
- **notes:** All column update expressions for a matching row are evaluated before any column value is changed. Type coercions (e.g., numeric widening) are applied automatically when expression and column types differ. Connector support for UPDATE varies; see connector-specific documentation.

---

## merge

- **syntax:**
```
MERGE INTO target_table [ @ branch_name ] [ [ AS ] target_alias ]
USING { source_table | query } [ [ AS ] source_alias ]
ON search_condition
when_clause [...]

when_clause is one of:
    WHEN MATCHED [ AND condition ]
        THEN DELETE

    WHEN MATCHED [ AND condition ]
        THEN UPDATE SET ( column = expression [, ...] )

    WHEN NOT MATCHED [ AND condition ]
        THEN INSERT [ ( column [, ...] ) ] VALUES ( expression [, ...] )
```
- **source_url:** https://trino.io/docs/current/sql/merge.html
- **version:** Trino 481 (current)
- **example_sql:**
```sql
-- Simple matched delete
MERGE INTO accounts t USING monthly_accounts_update s
    ON t.customer = s.customer
    WHEN MATCHED
        THEN DELETE

-- Update matched rows, insert unmatched
MERGE INTO accounts t USING monthly_accounts_update s
    ON (t.customer = s.customer)
    WHEN MATCHED
        THEN UPDATE SET purchases = s.purchases + t.purchases
    WHEN NOT MATCHED
        THEN INSERT (customer, purchases, address)
              VALUES(s.customer, s.purchases, s.address)

-- Multiple WHEN clauses with conditions
MERGE INTO accounts t USING monthly_accounts_update s
    ON (t.customer = s.customer)
    WHEN MATCHED AND s.address = 'Centreville'
        THEN DELETE
    WHEN MATCHED
        THEN UPDATE
            SET purchases = s.purchases + t.purchases, address = s.address
    WHEN NOT MATCHED
        THEN INSERT (customer, purchases, address)
              VALUES(s.customer, s.purchases, s.address)

-- Branch-qualified target
MERGE INTO accounts @ audit t USING monthly_accounts_update s
    ON t.customer = s.customer
    WHEN MATCHED
        THEN DELETE
```
- **notes:** WHEN clauses are processed in order; only the first matching clause executes for each source row. The query fails if a single target table row matches more than one source row. UPDATE expressions may reference both target and source columns; INSERT (NOT MATCHED) expressions reference only source columns. Only connectors that support MERGE can serve as the target; any connector supporting SELECT may serve as the source. The optional column list in INSERT specifies which target columns to populate; unlisted columns receive NULL.
