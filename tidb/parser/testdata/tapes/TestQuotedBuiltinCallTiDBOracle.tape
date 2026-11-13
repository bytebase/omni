-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQuotedBuiltinCallTiDBOracle default#1
> exec
| USE test
< ok 0 0
> exec
| DROP TABLE IF EXISTS tq
< ok 0 0
> exec
| CREATE TABLE tq (a INT, b VARCHAR(20), d DATE)
< ok 0 0
> exec
| SELECT `count`(a) FROM tq
< error
| 1305 42000 FUNCTION test.count does not exist
> exec
| SELECT `count`()
< error
| 1305 42000 FUNCTION test.count does not exist
> exec
| SELECT `sum`(a) FROM tq
< error
| 1305 42000 FUNCTION test.sum does not exist
> exec
| SELECT `trim`(b) FROM tq
< ok 0 0
> exec
| SELECT `trim`(b, 'x') FROM tq
< ok 0 0
> exec
| SELECT `trim`() FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'trim'
> exec
| SELECT `substring`(b, 1, 2) FROM tq
< ok 0 0
> exec
| SELECT `substring`(b) FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'substring'
> exec
| SELECT `substring`(b, 1, 2, 3) FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'substring'
> exec
| SELECT `substr`(b, 1) FROM tq
< ok 0 0
> exec
| SELECT `group_concat`(b) FROM tq
< error
| 1305 42000 FUNCTION test.group_concat does not exist
> exec
| SELECT `group_concat`(a, b) FROM tq
< error
| 1305 42000 FUNCTION test.group_concat does not exist
> exec
| SELECT `rank`() FROM tq
< error
| 1305 42000 FUNCTION test.rank does not exist
> exec
| SELECT `row_number`() FROM tq
< error
| 1305 42000 FUNCTION test.row_number does not exist
> exec
| SELECT `cast`(a) FROM tq
< ok 0 0
> exec
| SELECT `extract`(d) FROM tq
< error
| 1105 HY000 runtime error: index out of range [1] with length 1
> exec
| SELECT `position`(b) FROM tq
< error
| 1105 HY000 runtime error: index out of range [1] with length 1
> exec
| SELECT `char`(65)
< error
| 1305 42000 FUNCTION test.char does not exist
> exec
| SELECT `date_add`(d, 1) FROM tq
< error
| 1105 HY000 runtime error: index out of range [2] with length 2
> exec
| SELECT `f`(d + INTERVAL 1 DAY) FROM tq
< error
| 1305 42000 FUNCTION test.f does not exist
> exec
| SELECT test.count(a) FROM tq
< error
| 1305 42000 FUNCTION test.count does not exist
> exec
| SELECT test.trim(b) FROM tq
< error
| 1305 42000 FUNCTION test.trim does not exist
> exec
| SELECT test.trim(b, 'x') FROM tq
< error
| 1305 42000 FUNCTION test.trim does not exist
> exec
| SELECT test.substring(b, 1, 2) FROM tq
< error
| 1305 42000 FUNCTION test.substring does not exist
> exec
| SELECT test.substr(b, 1) FROM tq
< error
| 1305 42000 FUNCTION test.substr does not exist
> exec
| SELECT `count`(*) FROM tq
< error
| 1064 42000 …line 1 column 16 near "*) FROM tq" 
> exec
| SELECT `COUNT`(*) FROM tq
< error
| 1064 42000 …line 1 column 16 near "*) FROM tq" 
> exec
| SELECT `count`(*, a) FROM tq
< error
| 1064 42000 …line 1 column 16 near "*, a) FROM tq" 
> exec
| SELECT `count`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 23 near "DISTINCT a) FROM tq" 
> exec
| SELECT `count`(DISTINCT a, b) FROM tq
< error
| 1064 42000 …line 1 column 23 near "DISTINCT a, b) FROM tq" 
> exec
| SELECT `sum`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 21 near "DISTINCT a) FROM tq" 
> exec
| SELECT `avg`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 21 near "DISTINCT a) FROM tq" 
> exec
| SELECT `min`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 21 near "DISTINCT a) FROM tq" 
> exec
| SELECT `max`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 21 near "DISTINCT a) FROM tq" 
> exec
| SELECT `json_arrayagg`(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 31 near "DISTINCT a) FROM tq" 
> exec
| SELECT `group_concat`(DISTINCT b) FROM tq
< error
| 1064 42000 …line 1 column 30 near "DISTINCT b) FROM tq" 
> exec
| SELECT `group_concat`(b SEPARATOR ',') FROM tq
< error
| 1064 42000 …line 1 column 33 near "SEPARATOR ',') FROM tq" 
> exec
| SELECT `group_concat`(b ORDER BY b) FROM tq
< error
| 1064 42000 …line 1 column 29 near "ORDER BY b) FROM tq" 
> exec
| SELECT `group_concat`(b ORDER BY b SEPARATOR ',') FROM tq
< error
| 1064 42000 …line 1 column 29 near "ORDER BY b SEPARATOR ',') FROM tq" 
> exec
| SELECT `trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 21 near "LEADING 'x' FROM b) FROM tq" 
> exec
| SELECT `trim`(BOTH FROM b) FROM tq
< error
| 1064 42000 …line 1 column 18 near "BOTH FROM b) FROM tq" 
> exec
| SELECT `trim`('x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 22 near "FROM b) FROM tq" 
> exec
| SELECT `substring`(b FROM 1) FROM tq
< error
| 1064 42000 …line 1 column 25 near "FROM 1) FROM tq" 
> exec
| SELECT `substring`(b FROM 1 FOR 2) FROM tq
< error
| 1064 42000 …line 1 column 25 near "FROM 1 FOR 2) FROM tq" 
> exec
| SELECT `substr`(b FROM 1) FROM tq
< error
| 1064 42000 …line 1 column 22 near "FROM 1) FROM tq" 
> exec
| SELECT `date_add`(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …line 1 column 36 near ") FROM tq" 
> exec
| SELECT `adddate`(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …line 1 column 35 near ") FROM tq" 
> exec
| SELECT `sum`(a) OVER () FROM tq
< error
| 1064 42000 …line 1 column 20 near "OVER () FROM tq" 
> exec
| SELECT `rank`() OVER () FROM tq
< error
| 1064 42000 …line 1 column 20 near "OVER () FROM tq" 
> exec
| SELECT `row_number`() OVER () FROM tq
< error
| 1064 42000 …line 1 column 26 near "OVER () FROM tq" 
> exec
| SELECT `lag`(a) OVER () FROM tq
< error
| 1064 42000 …line 1 column 20 near "OVER () FROM tq" 
> exec
| SELECT `json_arrayagg`(a) OVER () FROM tq
< error
| 1064 42000 …line 1 column 30 near "OVER () FROM tq" 
> exec
| SELECT `json_objectagg`(a, b) OVER () FROM tq
< error
| 1064 42000 …line 1 column 34 near "OVER () FROM tq" 
> exec
| SELECT `count`(a) OVER w FROM tq WINDOW w AS ()
< error
| 1064 42000 …line 1 column 22 near "OVER w FROM tq WINDOW w AS ()" 
> exec
| SELECT `count`(ALL a) FROM tq
< error
| 1064 42000 …line 1 column 18 near "ALL a) FROM tq" 
> exec
| SELECT `cast`(a AS CHAR) FROM tq
< error
| 1064 42000 …line 1 column 18 near "AS CHAR) FROM tq" 
> exec
| SELECT `convert`(b USING utf8mb4) FROM tq
< error
| 1064 42000 …line 1 column 24 near "USING utf8mb4) FROM tq" 
> exec
| SELECT `convert`(b, CHAR) FROM tq
< error
| 1064 42000 …line 1 column 25 near ") FROM tq" 
> exec
| SELECT `extract`(YEAR FROM d) FROM tq
< error
| 1064 42000 …line 1 column 26 near "FROM d) FROM tq" 
> exec
| SELECT `position`('a' IN b) FROM tq
< error
| 1064 42000 …line 1 column 26 near "b) FROM tq" 
> exec
| SELECT `char`(65 USING utf8mb4)
< error
| 1064 42000 …line 1 column 22 near "USING utf8mb4)" 
> exec
| SELECT * FROM tq WHERE `match`(b) AGAINST ('x')
< error
| 1064 42000 …line 1 column 41 near "AGAINST ('x')" 
> exec
| SELECT * FROM `json_table`('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt
< error
| 1064 42000 …line 1 column 27 near "('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt" 
> exec
| SELECT test.count(*) FROM tq
< error
| 1064 42000 …line 1 column 19 near "*) FROM tq" 
> exec
| SELECT test.count(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 26 near "DISTINCT a) FROM tq" 
> exec
| SELECT test.sum(DISTINCT a) FROM tq
< error
| 1064 42000 …line 1 column 24 near "DISTINCT a) FROM tq" 
> exec
| SELECT test.trim(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 24 near "LEADING 'x' FROM b) FROM tq" 
> exec
| SELECT test.trim(TRAILING 'x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 25 near "TRAILING 'x' FROM b) FROM tq" 
> exec
| SELECT test.trim('x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 25 near "FROM b) FROM tq" 
> exec
| SELECT `test`.trim(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 26 near "LEADING 'x' FROM b) FROM tq" 
> exec
| SELECT test.`trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …line 1 column 26 near "LEADING 'x' FROM b) FROM tq" 
> exec
| SELECT test.substring(b FROM 1) FROM tq
< error
| 1064 42000 …line 1 column 28 near "FROM 1) FROM tq" 
> exec
| SELECT test.substring(b FROM 1 FOR 2) FROM tq
< error
| 1064 42000 …line 1 column 28 near "FROM 1 FOR 2) FROM tq" 
> exec
| SELECT test.substr(b FROM 1) FROM tq
< error
| 1064 42000 …line 1 column 25 near "FROM 1) FROM tq" 
> exec
| SELECT test.group_concat(b SEPARATOR ',') FROM tq
< error
| 1064 42000 …line 1 column 36 near "SEPARATOR ',') FROM tq" 
> exec
| SELECT test.group_concat(b ORDER BY b) FROM tq
< error
| 1064 42000 …line 1 column 32 near "ORDER BY b) FROM tq" 
> exec
| SELECT test.group_concat(DISTINCT b) FROM tq
< error
| 1064 42000 …line 1 column 33 near "DISTINCT b) FROM tq" 
> exec
| SELECT test.date_add(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …line 1 column 39 near ") FROM tq" 
> exec
| SELECT test.sum(a) OVER () FROM tq
< error
| 1064 42000 …line 1 column 23 near "OVER () FROM tq" 
> exec
| SELECT test.count(a) OVER () FROM tq
< error
| 1064 42000 …line 1 column 25 near "OVER () FROM tq" 
> exec
| SELECT COUNT(*) FROM tq
< ok 0 0
> exec
| SELECT COUNT(DISTINCT a) FROM tq
< ok 0 0
> exec
| SELECT TRIM(LEADING 'x' FROM b) FROM tq
< ok 0 0
> exec
| SELECT SUBSTRING(b FROM 1 FOR 2) FROM tq
< ok 0 0
> exec
| SELECT GROUP_CONCAT(b ORDER BY b SEPARATOR ',') FROM tq
< ok 0 0
> exec
| SELECT SUM(a) OVER () FROM tq
< ok 0 0
> exec
| SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM tq
< ok 0 0
> exec
| SELECT JSON_ARRAYAGG(a) OVER () FROM tq
< ok 0 0
> exec
| DROP TABLE IF EXISTS tq
< ok 0 0
