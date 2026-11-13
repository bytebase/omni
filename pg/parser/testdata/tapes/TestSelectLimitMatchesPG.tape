-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSelectLimitMatchesPG default#1
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 5 LIMIT 10
< error
| 42601 syntax error at or near "LIMIT"
| position: 30
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 OFFSET 2
< error
| 42601 syntax error at or near "OFFSET"
| position: 31
> exec
| SELECT 1/0;
| SELECT 1 LIMIT ALL LIMIT 1
< error
| 42601 syntax error at or near "LIMIT"
| position: 32
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 FETCH FIRST 1 ROW ONLY
< error
| 42601 syntax error at or near "FETCH"
| position: 30
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ROW ONLY LIMIT 1
< error
| 42601 syntax error at or near "LIMIT"
| position: 45
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ROW ONLY FETCH FIRST 2 ROWS ONLY
< error
| 42601 syntax error at or near "FETCH"
| position: 45
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 LIMIT 1 OFFSET 2
< error
| 42601 syntax error at or near "OFFSET"
| position: 39
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 OFFSET 2 LIMIT 3
< error
| 42601 syntax error at or near "LIMIT"
| position: 39
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 ROWS OFFSET 2
< error
| 42601 syntax error at or near "OFFSET"
| position: 36
> exec
| SELECT 1/0;
| SELECT 1 LIMIT
< error
| 42601 syntax error at end of input
| position: 27
> exec
| SELECT 1/0;
| SELECT 1 OFFSET
< error
| 42601 syntax error at end of input
| position: 28
> exec
| SELECT 1/0;
| SELECT * FROM (SELECT 1 LIMIT 1 LIMIT 2) s
< error
| 42601 syntax error at or near "LIMIT"
| position: 45
> exec
| SELECT 1/0;
| SELECT 1 WHERE EXISTS (SELECT 1 LIMIT 1 LIMIT 2)
< error
| 42601 syntax error at or near "LIMIT"
| position: 53
> exec
| SELECT 1/0;
| SELECT ARRAY(SELECT 1 LIMIT 1 LIMIT 2)
< error
| 42601 syntax error at or near "LIMIT"
| position: 43
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) SELECT 1 LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 51
> exec
| SELECT 1/0;
| INSERT INTO t SELECT 1 LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 44
> exec
| SELECT 1/0;
| CREATE TABLE x AS SELECT 1 LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 48
> exec
| SELECT 1/0;
| CREATE VIEW v AS SELECT 1 OFFSET 1 OFFSET 2
< error
| 42601 syntax error at or near "OFFSET"
| position: 48
> exec
| SELECT 1/0;
| CREATE RECURSIVE VIEW v (a) AS SELECT 1 LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 61
> exec
| SELECT 1/0;
| DECLARE c CURSOR FOR SELECT 1 LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 51
> exec
| SELECT 1/0;
| TABLE t LIMIT 1 LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 29
> exec
| SELECT 1/0;
| SELECT 1 LIMIT ALL
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT NULL
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 OFFSET 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 2 LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 FETCH FIRST 1 ROW ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ROW ONLY OFFSET 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| VALUES (1) LIMIT 1 OFFSET 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT * FROM (SELECT 1 LIMIT 1) s LIMIT 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 UNION SELECT 2 LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) UNION (SELECT 2 LIMIT 1) LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) UNION SELECT 2 LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 UNION SELECT 2 ORDER BY 1 LIMIT 1 OFFSET 2 FOR UPDATE
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) SELECT * FROM a ORDER BY 1 LIMIT 1 OFFSET 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 FOR UPDATE LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 41
> exec
| SELECT 1/0;
| SELECT 1 FOR UPDATE LIMIT 2 FOR SHARE
< error
| 42601 syntax error at or near "FOR"
| position: 41
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 FOR UPDATE OFFSET 2
< error
| 42601 syntax error at or near "OFFSET"
| position: 41
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 FOR UPDATE LIMIT 2
< error
| 42601 syntax error at or near "LIMIT"
| position: 42
> exec
| SELECT 1/0;
| SELECT 1 FOR READ ONLY FOR UPDATE
< error
| 42601 syntax error at or near "FOR"
| position: 36
> exec
| SELECT 1/0;
| SELECT 1 FOR UPDATE LIMIT 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FOR UPDATE LIMIT 2 OFFSET 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 2 FOR UPDATE FOR SHARE
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FOR UPDATE FOR SHARE LIMIT 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FOR READ ONLY LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1 FOR READ ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1, 2
< error
| 42601 LIMIT #,# syntax is not supported
| position: 22
| hint: Use separate LIMIT and OFFSET clauses.
> exec
| SELECT 1/0;
| SELECT 1 LIMIT ALL, 2
< error
| 42601 LIMIT #,# syntax is not supported
| position: 22
| hint: Use separate LIMIT and OFFSET clauses.
> exec
| SELECT 1/0;
| SELECT 1 LIMIT 1,
< error
| 42601 syntax error at end of input
| position: 30
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 + 1 ROWS
< error
| 42601 syntax error at or near "ROWS"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1::int ROWS
< error
| 42601 syntax error at or near "ROWS"
| position: 36
> exec
| SELECT 1/0;
| SELECT 1 OFFSET -a ROWS
< error
| 42601 syntax error at or near "ROWS"
| position: 32
> exec
| SELECT 1/0;
| SELECT 1 OFFSET -1::int ROWS
< error
| 42601 syntax error at or near "ROWS"
| position: 37
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 ROW
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET (1 + 1) ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET -1 ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET +1 ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET - 1.5 ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET $1 ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET a[1] ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET (ARRAY[1])[1] ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET f(1) ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 + 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET -1 * 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET NOT true
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET 1 IS NULL
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ONLY
< error
| 42601 syntax error at or near "ONLY"
| position: 36
> exec
| SELECT 1/0;
| SELECT 1 FETCH 1 ROWS ONLY
< error
| 42601 syntax error at or near "1"
| position: 28
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ROWS
< error
| 42601 syntax error at end of input
| position: 40
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST +a ROWS ONLY
< error
| 42601 syntax error at or near "a"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST -a ROWS ONLY
< error
| 42601 syntax error at or near "a"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1+1 ROWS ONLY
< error
| 42601 syntax error at or near "+"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 FETCH NEXT ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST ROW ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST -1 ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST +1 ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1.5 ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST -1.5 ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST (1+1) ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST $1 ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST a[1] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST ARRAY[1][1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 42
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST f(1)[1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 38
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1[1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 'a'[1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 37
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST CASE WHEN true THEN ARRAY[1] END[1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 66
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST ARRAY(SELECT 1)[1] ROWS ONLY
< error
| 42601 syntax error at or near "["
| position: 49
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST (ARRAY[1])[1] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST $1[1] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST a.b[1] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST a[1][2] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST (SELECT ARRAY[1])[1] ROWS ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET $1[1] ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 OFFSET (SELECT ARRAY[1])[1] ROWS
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 ORDER BY 1 FETCH NEXT 2 ROW WITH TIES
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT ALL) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT ALL
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 FETCH FIRST 1 ROW ONLY) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 53
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) FETCH FIRST 2 ROWS ONLY
< error
| 42601 multiple LIMIT clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) FETCH FIRST ROW ONLY
< error
| 42601 multiple LIMIT clauses not allowed
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) FETCH FIRST -2 ROWS ONLY
< error
| 42601 multiple LIMIT clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) FETCH FIRST +2 ROWS ONLY
< error
| 42601 multiple LIMIT clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) FETCH FIRST (2) ROWS ONLY
< error
| 42601 multiple LIMIT clauses not allowed
| position: 45
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) OFFSET 2
< error
| 42601 multiple OFFSET clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) OFFSET 2 ROWS
< error
| 42601 multiple OFFSET clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) OFFSET -2 ROWS
< error
| 42601 multiple OFFSET clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1 OFFSET 1) OFFSET 2
< error
| 42601 multiple OFFSET clauses not allowed
| position: 48
> exec
| SELECT 1/0;
| ((SELECT 1 LIMIT 1)) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) ORDER BY 1
< error
| 42601 multiple ORDER BY clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) ORDER BY 2, 3
< error
| 42601 multiple ORDER BY clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) ORDER BY a.b
< error
| 42601 multiple ORDER BY clauses not allowed
| position: 44
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES) OFFSET 1
< error
| 42601 multiple limit options not allowed
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES) LIMIT 1
< error
| 42601 multiple LIMIT clauses not allowed
| position: 69
> exec
| SELECT 1/0;
| SELECT 1 FETCH FIRST 1 ROW WITH TIES
< error
| 42601 WITH TIES cannot be specified without ORDER BY clause
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) SELECT * FROM a FETCH FIRST 1 ROW WITH TIES
< error
| 42601 WITH TIES cannot be specified without ORDER BY clause
> exec
| SELECT 1/0;
| SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES FOR UPDATE SKIP LOCKED
< error
| 42601 SKIP LOCKED and WITH TIES options cannot be used together
> exec
| SELECT 1/0;
| SELECT 1 ORDER BY 1 FOR UPDATE SKIP LOCKED FETCH FIRST 1 ROW WITH TIES
< error
| 42601 SKIP LOCKED and WITH TIES options cannot be used together
> exec
| SELECT 1/0;
| (SELECT 1 FOR UPDATE SKIP LOCKED) ORDER BY 1 FETCH FIRST 1 ROW WITH TIES
< error
| 42601 SKIP LOCKED and WITH TIES options cannot be used together
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) (WITH b AS (SELECT 2) SELECT 1)
< error
| 42601 multiple WITH clauses not allowed
| position: 13
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) (WITH b AS (SELECT 2) SELECT 1) ORDER BY 1
< error
| 42601 multiple WITH clauses not allowed
| position: 13
> exec
| SELECT 1/0;
| WITH a AS (SELECT 1) (SELECT 1 LIMIT 1) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 59
> exec
| SELECT 1/0;
| CREATE VIEW v AS (SELECT 1 LIMIT 1) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 55
> exec
| SELECT 1/0;
| CREATE MATERIALIZED VIEW v AS (SELECT 1 OFFSET 1) OFFSET 2
< error
| 42601 multiple OFFSET clauses not allowed
| position: 70
> exec
| SELECT 1/0;
| CREATE TABLE x AS (SELECT 1 LIMIT 1) LIMIT 2
< error
| 42601 multiple LIMIT clauses not allowed
| position: 56
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) OFFSET 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) LIMIT 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1) ORDER BY 1 LIMIT 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 FOR UPDATE) FOR SHARE
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW ONLY) OFFSET 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) FETCH FIRST 1 ROW WITH TIES
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT (2)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 39
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT ((2))
< error
| 42601 multiple LIMIT clauses not allowed
| position: 40
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT ( /* c */ 2)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 48
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT (2)::int
< error
| 42601 multiple LIMIT clauses not allowed
| position: 39
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT (SELECT 2)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT ((SELECT 2) + 1)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 39
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT (-2)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 39
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT 1 + 1
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT -1
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT x.y
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT f(1)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 LIMIT 1) LIMIT CAST(1 AS int)
< error
| 42601 multiple LIMIT clauses not allowed
| position: 38
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) OFFSET (2)
< error
| 42601 multiple OFFSET clauses not allowed
| position: 41
> exec
| SELECT 1/0;
| (SELECT 1 OFFSET 1) OFFSET (ARRAY[1])[1] ROWS
< error
| 42601 multiple OFFSET clauses not allowed
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) ORDER BY (1) DESC
< error
| 42601 multiple ORDER BY clauses not allowed
| position: 45
> exec
| SELECT 1/0;
| (SELECT 1 ORDER BY 1) ORDER BY (a).b
< error
| 42601 multiple ORDER BY clauses not allowed
