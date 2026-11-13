-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestKeywordFunctionNamesOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT absent()
< error
| 1305 42000 FUNCTION test.absent does not exist
> exec ^
: string
| SELECT absent()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT accessible()
< error
| 1064 42000 …near 'accessible()' at line 1
> exec ^
: string
| SELECT accessible(1)
< error
| 1064 42000 …near 'accessible(1)' at line 1
> exec ^
: string
| SELECT accessible(1,2)
< error
| 1064 42000 …near 'accessible(1,2)' at line 1
> exec ^
: string
| SELECT accessible(1,2,3)
< error
| 1064 42000 …near 'accessible(1,2,3)' at line 1
> exec ^
: string
| SELECT accessible(1,2,3,4)
< error
| 1064 42000 …near 'accessible(1,2,3,4)' at line 1
> exec ^
: string
| SELECT accessible(*)
< error
| 1064 42000 …near 'accessible(*)' at line 1
> exec ^
: string
| SELECT accessible('a')
< error
| 1064 42000 …near 'accessible('a')' at line 1
> exec ^
: string
| SELECT accessible(a USING utf8mb4)
< error
| 1064 42000 …near 'accessible(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT accessible(1 AS SIGNED)
< error
| 1064 42000 …near 'accessible(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT accessible(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'accessible(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT accessible(1 FROM 2)
< error
| 1064 42000 …near 'accessible(1 FROM 2)' at line 1
> exec ^
: string
| SELECT accessible(1 IN 2)
< error
| 1064 42000 …near 'accessible(1 IN 2)' at line 1
> exec ^
: string
| SELECT accessible(DISTINCT 1)
< error
| 1064 42000 …near 'accessible(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT accessible() OVER ()
< error
| 1064 42000 …near 'accessible() OVER ()' at line 1
> exec ^
: string
| SELECT accessible(1) OVER ()
< error
| 1064 42000 …near 'accessible(1) OVER ()' at line 1
> exec ^
: string
| SELECT accessible(1,2) OVER ()
< error
| 1064 42000 …near 'accessible(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT accessible(1,2,3) OVER ()
< error
| 1064 42000 …near 'accessible(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT accessible(1) FROM t
< error
| 1064 42000 …near 'accessible(1) FROM t' at line 1
> exec ^
: string
| SELECT accessible(SELECT 1)
< error
| 1064 42000 …near 'accessible(SELECT 1)' at line 1
> exec ^
: string
| SELECT accessible(YEAR FROM 1)
< error
| 1064 42000 …near 'accessible(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT accessible(DAY, 1, 2)
< error
| 1064 42000 …near 'accessible(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT accessible(a)
< error
| 1064 42000 …near 'accessible(a)' at line 1
> exec ^
: string
| SELECT accessible(a, b)
< error
| 1064 42000 …near 'accessible(a, b)' at line 1
> exec ^
: string
| SELECT accessible(a) OVER ()
< error
| 1064 42000 …near 'accessible(a) OVER ()' at line 1
> exec ^
: string
| SELECT account()
< error
| 1630 42000 FUNCTION test.account does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT account()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT action()
< error
| 1630 42000 FUNCTION test.action does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT action()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT active()
< error
| 1305 42000 FUNCTION test.active does not exist
> exec ^
: string
| SELECT active()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT add()
< error
| 1064 42000 …near 'add()' at line 1
> exec ^
: string
| SELECT add(1)
< error
| 1064 42000 …near 'add(1)' at line 1
> exec ^
: string
| SELECT add(1,2)
< error
| 1064 42000 …near 'add(1,2)' at line 1
> exec ^
: string
| SELECT add(1,2,3)
< error
| 1064 42000 …near 'add(1,2,3)' at line 1
> exec ^
: string
| SELECT add(1,2,3,4)
< error
| 1064 42000 …near 'add(1,2,3,4)' at line 1
> exec ^
: string
| SELECT add(*)
< error
| 1064 42000 …near 'add(*)' at line 1
> exec ^
: string
| SELECT add('a')
< error
| 1064 42000 …near 'add('a')' at line 1
> exec ^
: string
| SELECT add(a USING utf8mb4)
< error
| 1064 42000 …near 'add(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT add(1 AS SIGNED)
< error
| 1064 42000 …near 'add(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT add(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'add(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT add(1 FROM 2)
< error
| 1064 42000 …near 'add(1 FROM 2)' at line 1
> exec ^
: string
| SELECT add(1 IN 2)
< error
| 1064 42000 …near 'add(1 IN 2)' at line 1
> exec ^
: string
| SELECT add(DISTINCT 1)
< error
| 1064 42000 …near 'add(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT add() OVER ()
< error
| 1064 42000 …near 'add() OVER ()' at line 1
> exec ^
: string
| SELECT add(1) OVER ()
< error
| 1064 42000 …near 'add(1) OVER ()' at line 1
> exec ^
: string
| SELECT add(1,2) OVER ()
< error
| 1064 42000 …near 'add(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT add(1,2,3) OVER ()
< error
| 1064 42000 …near 'add(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT add(1) FROM t
< error
| 1064 42000 …near 'add(1) FROM t' at line 1
> exec ^
: string
| SELECT add(SELECT 1)
< error
| 1064 42000 …near 'add(SELECT 1)' at line 1
> exec ^
: string
| SELECT add(YEAR FROM 1)
< error
| 1064 42000 …near 'add(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT add(DAY, 1, 2)
< error
| 1064 42000 …near 'add(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT add(a)
< error
| 1064 42000 …near 'add(a)' at line 1
> exec ^
: string
| SELECT add(a, b)
< error
| 1064 42000 …near 'add(a, b)' at line 1
> exec ^
: string
| SELECT add(a) OVER ()
< error
| 1064 42000 …near 'add(a) OVER ()' at line 1
> exec ^
: string
| SELECT adddate()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT adddate(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT adddate(1,2)
< ok 0 0
> exec ^
: string
| SELECT admin()
< error
| 1630 42000 FUNCTION test.admin does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT admin()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT after()
< error
| 1630 42000 FUNCTION test.after does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT after()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT against()
< error
| 1630 42000 FUNCTION test.against does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT against()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT aggregate()
< error
| 1630 42000 FUNCTION test.aggregate does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT aggregate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT algorithm()
< error
| 1630 42000 FUNCTION test.algorithm does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT algorithm()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT all()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT all(1)
< ok 0 0
> exec ^
: string
| SELECT allow_missing_files()
< error
| 1305 42000 FUNCTION test.allow_missing_files does not exist
> exec ^
: string
| SELECT allow_missing_files()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT alter()
< error
| 1064 42000 …near 'alter()' at line 1
> exec ^
: string
| SELECT alter(1)
< error
| 1064 42000 …near 'alter(1)' at line 1
> exec ^
: string
| SELECT alter(1,2)
< error
| 1064 42000 …near 'alter(1,2)' at line 1
> exec ^
: string
| SELECT alter(1,2,3)
< error
| 1064 42000 …near 'alter(1,2,3)' at line 1
> exec ^
: string
| SELECT alter(1,2,3,4)
< error
| 1064 42000 …near 'alter(1,2,3,4)' at line 1
> exec ^
: string
| SELECT alter(*)
< error
| 1064 42000 …near 'alter(*)' at line 1
> exec ^
: string
| SELECT alter('a')
< error
| 1064 42000 …near 'alter('a')' at line 1
> exec ^
: string
| SELECT alter(a USING utf8mb4)
< error
| 1064 42000 …near 'alter(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT alter(1 AS SIGNED)
< error
| 1064 42000 …near 'alter(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT alter(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'alter(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT alter(1 FROM 2)
< error
| 1064 42000 …near 'alter(1 FROM 2)' at line 1
> exec ^
: string
| SELECT alter(1 IN 2)
< error
| 1064 42000 …near 'alter(1 IN 2)' at line 1
> exec ^
: string
| SELECT alter(DISTINCT 1)
< error
| 1064 42000 …near 'alter(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT alter() OVER ()
< error
| 1064 42000 …near 'alter() OVER ()' at line 1
> exec ^
: string
| SELECT alter(1) OVER ()
< error
| 1064 42000 …near 'alter(1) OVER ()' at line 1
> exec ^
: string
| SELECT alter(1,2) OVER ()
< error
| 1064 42000 …near 'alter(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT alter(1,2,3) OVER ()
< error
| 1064 42000 …near 'alter(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT alter(1) FROM t
< error
| 1064 42000 …near 'alter(1) FROM t' at line 1
> exec ^
: string
| SELECT alter(SELECT 1)
< error
| 1064 42000 …near 'alter(SELECT 1)' at line 1
> exec ^
: string
| SELECT alter(YEAR FROM 1)
< error
| 1064 42000 …near 'alter(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT alter(DAY, 1, 2)
< error
| 1064 42000 …near 'alter(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT alter(a)
< error
| 1064 42000 …near 'alter(a)' at line 1
> exec ^
: string
| SELECT alter(a, b)
< error
| 1064 42000 …near 'alter(a, b)' at line 1
> exec ^
: string
| SELECT alter(a) OVER ()
< error
| 1064 42000 …near 'alter(a) OVER ()' at line 1
> exec ^
: string
| SELECT always()
< error
| 1630 42000 FUNCTION test.always does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT always()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT analyze()
< error
| 1064 42000 …near 'analyze()' at line 1
> exec ^
: string
| SELECT analyze(1)
< error
| 1064 42000 …near 'analyze(1)' at line 1
> exec ^
: string
| SELECT analyze(1,2)
< error
| 1064 42000 …near 'analyze(1,2)' at line 1
> exec ^
: string
| SELECT analyze(1,2,3)
< error
| 1064 42000 …near 'analyze(1,2,3)' at line 1
> exec ^
: string
| SELECT analyze(1,2,3,4)
< error
| 1064 42000 …near 'analyze(1,2,3,4)' at line 1
> exec ^
: string
| SELECT analyze(*)
< error
| 1064 42000 …near 'analyze(*)' at line 1
> exec ^
: string
| SELECT analyze('a')
< error
| 1064 42000 …near 'analyze('a')' at line 1
> exec ^
: string
| SELECT analyze(a USING utf8mb4)
< error
| 1064 42000 …near 'analyze(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT analyze(1 AS SIGNED)
< error
| 1064 42000 …near 'analyze(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT analyze(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'analyze(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT analyze(1 FROM 2)
< error
| 1064 42000 …near 'analyze(1 FROM 2)' at line 1
> exec ^
: string
| SELECT analyze(1 IN 2)
< error
| 1064 42000 …near 'analyze(1 IN 2)' at line 1
> exec ^
: string
| SELECT analyze(DISTINCT 1)
< error
| 1064 42000 …near 'analyze(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT analyze() OVER ()
< error
| 1064 42000 …near 'analyze() OVER ()' at line 1
> exec ^
: string
| SELECT analyze(1) OVER ()
< error
| 1064 42000 …near 'analyze(1) OVER ()' at line 1
> exec ^
: string
| SELECT analyze(1,2) OVER ()
< error
| 1064 42000 …near 'analyze(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT analyze(1,2,3) OVER ()
< error
| 1064 42000 …near 'analyze(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT analyze(1) FROM t
< error
| 1064 42000 …near 'analyze(1) FROM t' at line 1
> exec ^
: string
| SELECT analyze(SELECT 1)
< error
| 1064 42000 …near 'analyze(SELECT 1)' at line 1
> exec ^
: string
| SELECT analyze(YEAR FROM 1)
< error
| 1064 42000 …near 'analyze(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT analyze(DAY, 1, 2)
< error
| 1064 42000 …near 'analyze(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT analyze(a)
< error
| 1064 42000 …near 'analyze(a)' at line 1
> exec ^
: string
| SELECT analyze(a, b)
< error
| 1064 42000 …near 'analyze(a, b)' at line 1
> exec ^
: string
| SELECT analyze(a) OVER ()
< error
| 1064 42000 …near 'analyze(a) OVER ()' at line 1
> exec ^
: string
| SELECT and()
< error
| 1064 42000 …near 'and()' at line 1
> exec ^
: string
| SELECT and(1)
< error
| 1064 42000 …near 'and(1)' at line 1
> exec ^
: string
| SELECT and(1,2)
< error
| 1064 42000 …near 'and(1,2)' at line 1
> exec ^
: string
| SELECT and(1,2,3)
< error
| 1064 42000 …near 'and(1,2,3)' at line 1
> exec ^
: string
| SELECT and(1,2,3,4)
< error
| 1064 42000 …near 'and(1,2,3,4)' at line 1
> exec ^
: string
| SELECT and(*)
< error
| 1064 42000 …near 'and(*)' at line 1
> exec ^
: string
| SELECT and('a')
< error
| 1064 42000 …near 'and('a')' at line 1
> exec ^
: string
| SELECT and(a USING utf8mb4)
< error
| 1064 42000 …near 'and(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT and(1 AS SIGNED)
< error
| 1064 42000 …near 'and(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT and(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'and(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT and(1 FROM 2)
< error
| 1064 42000 …near 'and(1 FROM 2)' at line 1
> exec ^
: string
| SELECT and(1 IN 2)
< error
| 1064 42000 …near 'and(1 IN 2)' at line 1
> exec ^
: string
| SELECT and(DISTINCT 1)
< error
| 1064 42000 …near 'and(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT and() OVER ()
< error
| 1064 42000 …near 'and() OVER ()' at line 1
> exec ^
: string
| SELECT and(1) OVER ()
< error
| 1064 42000 …near 'and(1) OVER ()' at line 1
> exec ^
: string
| SELECT and(1,2) OVER ()
< error
| 1064 42000 …near 'and(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT and(1,2,3) OVER ()
< error
| 1064 42000 …near 'and(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT and(1) FROM t
< error
| 1064 42000 …near 'and(1) FROM t' at line 1
> exec ^
: string
| SELECT and(SELECT 1)
< error
| 1064 42000 …near 'and(SELECT 1)' at line 1
> exec ^
: string
| SELECT and(YEAR FROM 1)
< error
| 1064 42000 …near 'and(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT and(DAY, 1, 2)
< error
| 1064 42000 …near 'and(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT and(a)
< error
| 1064 42000 …near 'and(a)' at line 1
> exec ^
: string
| SELECT and(a, b)
< error
| 1064 42000 …near 'and(a, b)' at line 1
> exec ^
: string
| SELECT and(a) OVER ()
< error
| 1064 42000 …near 'and(a) OVER ()' at line 1
> exec ^
: string
| SELECT any()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT any(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT any(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT any(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT any(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT any(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT any('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT any(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT any(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT any(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT any(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT any(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT any(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT any() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT any(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT any(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT any(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT any(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT any(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT any(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT any(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT any(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT any(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT any(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT array()
< error
| 1305 42000 FUNCTION test.array does not exist
> exec ^
: string
| SELECT array()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT as()
< error
| 1064 42000 …near 'as()' at line 1
> exec ^
: string
| SELECT as(1)
< error
| 1064 42000 …near 'as(1)' at line 1
> exec ^
: string
| SELECT as(1,2)
< error
| 1064 42000 …near 'as(1,2)' at line 1
> exec ^
: string
| SELECT as(1,2,3)
< error
| 1064 42000 …near 'as(1,2,3)' at line 1
> exec ^
: string
| SELECT as(1,2,3,4)
< error
| 1064 42000 …near 'as(1,2,3,4)' at line 1
> exec ^
: string
| SELECT as(*)
< error
| 1064 42000 …near 'as(*)' at line 1
> exec ^
: string
| SELECT as('a')
< error
| 1064 42000 …near 'as('a')' at line 1
> exec ^
: string
| SELECT as(a USING utf8mb4)
< error
| 1064 42000 …near 'as(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT as(1 AS SIGNED)
< error
| 1064 42000 …near 'as(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT as(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'as(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT as(1 FROM 2)
< error
| 1064 42000 …near 'as(1 FROM 2)' at line 1
> exec ^
: string
| SELECT as(1 IN 2)
< error
| 1064 42000 …near 'as(1 IN 2)' at line 1
> exec ^
: string
| SELECT as(DISTINCT 1)
< error
| 1064 42000 …near 'as(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT as() OVER ()
< error
| 1064 42000 …near 'as() OVER ()' at line 1
> exec ^
: string
| SELECT as(1) OVER ()
< error
| 1064 42000 …near 'as(1) OVER ()' at line 1
> exec ^
: string
| SELECT as(1,2) OVER ()
< error
| 1064 42000 …near 'as(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT as(1,2,3) OVER ()
< error
| 1064 42000 …near 'as(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT as(1) FROM t
< error
| 1064 42000 …near 'as(1) FROM t' at line 1
> exec ^
: string
| SELECT as(SELECT 1)
< error
| 1064 42000 …near 'as(SELECT 1)' at line 1
> exec ^
: string
| SELECT as(YEAR FROM 1)
< error
| 1064 42000 …near 'as(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT as(DAY, 1, 2)
< error
| 1064 42000 …near 'as(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT as(a)
< error
| 1064 42000 …near 'as(a)' at line 1
> exec ^
: string
| SELECT as(a, b)
< error
| 1064 42000 …near 'as(a, b)' at line 1
> exec ^
: string
| SELECT as(a) OVER ()
< error
| 1064 42000 …near 'as(a) OVER ()' at line 1
> exec ^
: string
| SELECT asc()
< error
| 1064 42000 …near 'asc()' at line 1
> exec ^
: string
| SELECT asc(1)
< error
| 1064 42000 …near 'asc(1)' at line 1
> exec ^
: string
| SELECT asc(1,2)
< error
| 1064 42000 …near 'asc(1,2)' at line 1
> exec ^
: string
| SELECT asc(1,2,3)
< error
| 1064 42000 …near 'asc(1,2,3)' at line 1
> exec ^
: string
| SELECT asc(1,2,3,4)
< error
| 1064 42000 …near 'asc(1,2,3,4)' at line 1
> exec ^
: string
| SELECT asc(*)
< error
| 1064 42000 …near 'asc(*)' at line 1
> exec ^
: string
| SELECT asc('a')
< error
| 1064 42000 …near 'asc('a')' at line 1
> exec ^
: string
| SELECT asc(a USING utf8mb4)
< error
| 1064 42000 …near 'asc(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT asc(1 AS SIGNED)
< error
| 1064 42000 …near 'asc(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT asc(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'asc(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT asc(1 FROM 2)
< error
| 1064 42000 …near 'asc(1 FROM 2)' at line 1
> exec ^
: string
| SELECT asc(1 IN 2)
< error
| 1064 42000 …near 'asc(1 IN 2)' at line 1
> exec ^
: string
| SELECT asc(DISTINCT 1)
< error
| 1064 42000 …near 'asc(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT asc() OVER ()
< error
| 1064 42000 …near 'asc() OVER ()' at line 1
> exec ^
: string
| SELECT asc(1) OVER ()
< error
| 1064 42000 …near 'asc(1) OVER ()' at line 1
> exec ^
: string
| SELECT asc(1,2) OVER ()
< error
| 1064 42000 …near 'asc(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT asc(1,2,3) OVER ()
< error
| 1064 42000 …near 'asc(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT asc(1) FROM t
< error
| 1064 42000 …near 'asc(1) FROM t' at line 1
> exec ^
: string
| SELECT asc(SELECT 1)
< error
| 1064 42000 …near 'asc(SELECT 1)' at line 1
> exec ^
: string
| SELECT asc(YEAR FROM 1)
< error
| 1064 42000 …near 'asc(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT asc(DAY, 1, 2)
< error
| 1064 42000 …near 'asc(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT asc(a)
< error
| 1064 42000 …near 'asc(a)' at line 1
> exec ^
: string
| SELECT asc(a, b)
< error
| 1064 42000 …near 'asc(a, b)' at line 1
> exec ^
: string
| SELECT asc(a) OVER ()
< error
| 1064 42000 …near 'asc(a) OVER ()' at line 1
> exec ^
: string
| SELECT ascii()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT ascii(1)
< ok 0 0
> exec ^
: string
| SELECT asensitive()
< error
| 1064 42000 …near 'asensitive()' at line 1
> exec ^
: string
| SELECT asensitive(1)
< error
| 1064 42000 …near 'asensitive(1)' at line 1
> exec ^
: string
| SELECT asensitive(1,2)
< error
| 1064 42000 …near 'asensitive(1,2)' at line 1
> exec ^
: string
| SELECT asensitive(1,2,3)
< error
| 1064 42000 …near 'asensitive(1,2,3)' at line 1
> exec ^
: string
| SELECT asensitive(1,2,3,4)
< error
| 1064 42000 …near 'asensitive(1,2,3,4)' at line 1
> exec ^
: string
| SELECT asensitive(*)
< error
| 1064 42000 …near 'asensitive(*)' at line 1
> exec ^
: string
| SELECT asensitive('a')
< error
| 1064 42000 …near 'asensitive('a')' at line 1
> exec ^
: string
| SELECT asensitive(a USING utf8mb4)
< error
| 1064 42000 …near 'asensitive(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT asensitive(1 AS SIGNED)
< error
| 1064 42000 …near 'asensitive(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT asensitive(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'asensitive(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT asensitive(1 FROM 2)
< error
| 1064 42000 …near 'asensitive(1 FROM 2)' at line 1
> exec ^
: string
| SELECT asensitive(1 IN 2)
< error
| 1064 42000 …near 'asensitive(1 IN 2)' at line 1
> exec ^
: string
| SELECT asensitive(DISTINCT 1)
< error
| 1064 42000 …near 'asensitive(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT asensitive() OVER ()
< error
| 1064 42000 …near 'asensitive() OVER ()' at line 1
> exec ^
: string
| SELECT asensitive(1) OVER ()
< error
| 1064 42000 …near 'asensitive(1) OVER ()' at line 1
> exec ^
: string
| SELECT asensitive(1,2) OVER ()
< error
| 1064 42000 …near 'asensitive(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT asensitive(1,2,3) OVER ()
< error
| 1064 42000 …near 'asensitive(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT asensitive(1) FROM t
< error
| 1064 42000 …near 'asensitive(1) FROM t' at line 1
> exec ^
: string
| SELECT asensitive(SELECT 1)
< error
| 1064 42000 …near 'asensitive(SELECT 1)' at line 1
> exec ^
: string
| SELECT asensitive(YEAR FROM 1)
< error
| 1064 42000 …near 'asensitive(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT asensitive(DAY, 1, 2)
< error
| 1064 42000 …near 'asensitive(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT asensitive(a)
< error
| 1064 42000 …near 'asensitive(a)' at line 1
> exec ^
: string
| SELECT asensitive(a, b)
< error
| 1064 42000 …near 'asensitive(a, b)' at line 1
> exec ^
: string
| SELECT asensitive(a) OVER ()
< error
| 1064 42000 …near 'asensitive(a) OVER ()' at line 1
> exec ^
: string
| SELECT assign_gtids_to_anonymous_transactions()
< error
| 1305 42000 FUNCTION test.assign_gtids_to_anonymous_transactions does not exist
> exec ^
: string
| SELECT assign_gtids_to_anonymous_transactions()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT at()
< error
| 1630 42000 FUNCTION test.at does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT at()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT attribute()
< error
| 1305 42000 FUNCTION test.attribute does not exist
> exec ^
: string
| SELECT attribute()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT authentication()
< error
| 1305 42000 FUNCTION test.authentication does not exist
> exec ^
: string
| SELECT authentication()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT auto()
< error
| 1630 42000 FUNCTION test.auto does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT auto()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT auto_increment()
< error
| 1630 42000 FUNCTION test.auto_increment does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT auto_increment()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT auto_refresh()
< error
| 1305 42000 FUNCTION test.auto_refresh does not exist
> exec ^
: string
| SELECT auto_refresh()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT auto_refresh_source()
< error
| 1305 42000 FUNCTION test.auto_refresh_source does not exist
> exec ^
: string
| SELECT auto_refresh_source()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT autoextend_size()
< error
| 1630 42000 FUNCTION test.autoextend_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT autoextend_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT avg()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT avg(1)
< ok 0 0
> exec ^
: string
| SELECT avg_row_length()
< error
| 1630 42000 FUNCTION test.avg_row_length does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT avg_row_length()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT backup()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT backup(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT backup(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT backup(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT backup(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT backup(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT backup('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT backup(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT backup(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT backup(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT backup(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT backup(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT backup(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT backup() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT backup(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT backup(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT backup(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT backup(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT backup(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT backup(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT backup(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT backup(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT backup(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT backup(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT before()
< error
| 1064 42000 …near 'before()' at line 1
> exec ^
: string
| SELECT before(1)
< error
| 1064 42000 …near 'before(1)' at line 1
> exec ^
: string
| SELECT before(1,2)
< error
| 1064 42000 …near 'before(1,2)' at line 1
> exec ^
: string
| SELECT before(1,2,3)
< error
| 1064 42000 …near 'before(1,2,3)' at line 1
> exec ^
: string
| SELECT before(1,2,3,4)
< error
| 1064 42000 …near 'before(1,2,3,4)' at line 1
> exec ^
: string
| SELECT before(*)
< error
| 1064 42000 …near 'before(*)' at line 1
> exec ^
: string
| SELECT before('a')
< error
| 1064 42000 …near 'before('a')' at line 1
> exec ^
: string
| SELECT before(a USING utf8mb4)
< error
| 1064 42000 …near 'before(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT before(1 AS SIGNED)
< error
| 1064 42000 …near 'before(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT before(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'before(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT before(1 FROM 2)
< error
| 1064 42000 …near 'before(1 FROM 2)' at line 1
> exec ^
: string
| SELECT before(1 IN 2)
< error
| 1064 42000 …near 'before(1 IN 2)' at line 1
> exec ^
: string
| SELECT before(DISTINCT 1)
< error
| 1064 42000 …near 'before(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT before() OVER ()
< error
| 1064 42000 …near 'before() OVER ()' at line 1
> exec ^
: string
| SELECT before(1) OVER ()
< error
| 1064 42000 …near 'before(1) OVER ()' at line 1
> exec ^
: string
| SELECT before(1,2) OVER ()
< error
| 1064 42000 …near 'before(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT before(1,2,3) OVER ()
< error
| 1064 42000 …near 'before(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT before(1) FROM t
< error
| 1064 42000 …near 'before(1) FROM t' at line 1
> exec ^
: string
| SELECT before(SELECT 1)
< error
| 1064 42000 …near 'before(SELECT 1)' at line 1
> exec ^
: string
| SELECT before(YEAR FROM 1)
< error
| 1064 42000 …near 'before(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT before(DAY, 1, 2)
< error
| 1064 42000 …near 'before(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT before(a)
< error
| 1064 42000 …near 'before(a)' at line 1
> exec ^
: string
| SELECT before(a, b)
< error
| 1064 42000 …near 'before(a, b)' at line 1
> exec ^
: string
| SELECT before(a) OVER ()
< error
| 1064 42000 …near 'before(a) OVER ()' at line 1
> exec ^
: string
| SELECT begin()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT begin(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT begin(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT begin(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT begin(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT begin(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT begin('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT begin(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT begin(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT begin(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT begin(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT begin(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT begin(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT begin() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT begin(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT begin(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT begin(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT begin(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT begin(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT begin(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT begin(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT begin(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT begin(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT begin(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT bernoulli()
< error
| 1305 42000 FUNCTION test.bernoulli does not exist
> exec ^
: string
| SELECT bernoulli()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT between()
< error
| 1064 42000 …near 'between()' at line 1
> exec ^
: string
| SELECT between(1)
< error
| 1064 42000 …near 'between(1)' at line 1
> exec ^
: string
| SELECT between(1,2)
< error
| 1064 42000 …near 'between(1,2)' at line 1
> exec ^
: string
| SELECT between(1,2,3)
< error
| 1064 42000 …near 'between(1,2,3)' at line 1
> exec ^
: string
| SELECT between(1,2,3,4)
< error
| 1064 42000 …near 'between(1,2,3,4)' at line 1
> exec ^
: string
| SELECT between(*)
< error
| 1064 42000 …near 'between(*)' at line 1
> exec ^
: string
| SELECT between('a')
< error
| 1064 42000 …near 'between('a')' at line 1
> exec ^
: string
| SELECT between(a USING utf8mb4)
< error
| 1064 42000 …near 'between(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT between(1 AS SIGNED)
< error
| 1064 42000 …near 'between(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT between(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'between(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT between(1 FROM 2)
< error
| 1064 42000 …near 'between(1 FROM 2)' at line 1
> exec ^
: string
| SELECT between(1 IN 2)
< error
| 1064 42000 …near 'between(1 IN 2)' at line 1
> exec ^
: string
| SELECT between(DISTINCT 1)
< error
| 1064 42000 …near 'between(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT between() OVER ()
< error
| 1064 42000 …near 'between() OVER ()' at line 1
> exec ^
: string
| SELECT between(1) OVER ()
< error
| 1064 42000 …near 'between(1) OVER ()' at line 1
> exec ^
: string
| SELECT between(1,2) OVER ()
< error
| 1064 42000 …near 'between(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT between(1,2,3) OVER ()
< error
| 1064 42000 …near 'between(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT between(1) FROM t
< error
| 1064 42000 …near 'between(1) FROM t' at line 1
> exec ^
: string
| SELECT between(SELECT 1)
< error
| 1064 42000 …near 'between(SELECT 1)' at line 1
> exec ^
: string
| SELECT between(YEAR FROM 1)
< error
| 1064 42000 …near 'between(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT between(DAY, 1, 2)
< error
| 1064 42000 …near 'between(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT between(a)
< error
| 1064 42000 …near 'between(a)' at line 1
> exec ^
: string
| SELECT between(a, b)
< error
| 1064 42000 …near 'between(a, b)' at line 1
> exec ^
: string
| SELECT between(a) OVER ()
< error
| 1064 42000 …near 'between(a) OVER ()' at line 1
> exec ^
: string
| SELECT bigint()
< error
| 1064 42000 …near 'bigint()' at line 1
> exec ^
: string
| SELECT bigint(1)
< error
| 1064 42000 …near 'bigint(1)' at line 1
> exec ^
: string
| SELECT bigint(1,2)
< error
| 1064 42000 …near 'bigint(1,2)' at line 1
> exec ^
: string
| SELECT bigint(1,2,3)
< error
| 1064 42000 …near 'bigint(1,2,3)' at line 1
> exec ^
: string
| SELECT bigint(1,2,3,4)
< error
| 1064 42000 …near 'bigint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT bigint(*)
< error
| 1064 42000 …near 'bigint(*)' at line 1
> exec ^
: string
| SELECT bigint('a')
< error
| 1064 42000 …near 'bigint('a')' at line 1
> exec ^
: string
| SELECT bigint(a USING utf8mb4)
< error
| 1064 42000 …near 'bigint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT bigint(1 AS SIGNED)
< error
| 1064 42000 …near 'bigint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT bigint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'bigint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT bigint(1 FROM 2)
< error
| 1064 42000 …near 'bigint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT bigint(1 IN 2)
< error
| 1064 42000 …near 'bigint(1 IN 2)' at line 1
> exec ^
: string
| SELECT bigint(DISTINCT 1)
< error
| 1064 42000 …near 'bigint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT bigint() OVER ()
< error
| 1064 42000 …near 'bigint() OVER ()' at line 1
> exec ^
: string
| SELECT bigint(1) OVER ()
< error
| 1064 42000 …near 'bigint(1) OVER ()' at line 1
> exec ^
: string
| SELECT bigint(1,2) OVER ()
< error
| 1064 42000 …near 'bigint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT bigint(1,2,3) OVER ()
< error
| 1064 42000 …near 'bigint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT bigint(1) FROM t
< error
| 1064 42000 …near 'bigint(1) FROM t' at line 1
> exec ^
: string
| SELECT bigint(SELECT 1)
< error
| 1064 42000 …near 'bigint(SELECT 1)' at line 1
> exec ^
: string
| SELECT bigint(YEAR FROM 1)
< error
| 1064 42000 …near 'bigint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT bigint(DAY, 1, 2)
< error
| 1064 42000 …near 'bigint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT bigint(a)
< error
| 1064 42000 …near 'bigint(a)' at line 1
> exec ^
: string
| SELECT bigint(a, b)
< error
| 1064 42000 …near 'bigint(a, b)' at line 1
> exec ^
: string
| SELECT bigint(a) OVER ()
< error
| 1064 42000 …near 'bigint(a) OVER ()' at line 1
> exec ^
: string
| SELECT binary()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT binary(1)
< ok 0 0
> exec ^
: string
| SELECT binlog()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT binlog(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT binlog(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT binlog(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT binlog(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT binlog(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT binlog('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT binlog(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT binlog(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT binlog(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT binlog(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT binlog(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT binlog(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT binlog() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT binlog(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT binlog(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT binlog(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT binlog(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT binlog(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT binlog(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT binlog(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT binlog(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT binlog(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT binlog(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT bit()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT bit(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT bit(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT bit(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT bit(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT bit(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT bit('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT bit(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT bit(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT bit(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT bit(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT bit(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT bit(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT bit() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT bit(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT bit(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT bit(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT bit(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT bit(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT bit(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT bit(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT bit(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT bit(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT bit(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT bit_and()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT bit_and(1)
< ok 0 0
> exec ^
: string
| SELECT bit_or()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT bit_or(1)
< ok 0 0
> exec ^
: string
| SELECT bit_xor()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT bit_xor(1)
< ok 0 0
> exec ^
: string
| SELECT blob()
< error
| 1064 42000 …near 'blob()' at line 1
> exec ^
: string
| SELECT blob(1)
< error
| 1064 42000 …near 'blob(1)' at line 1
> exec ^
: string
| SELECT blob(1,2)
< error
| 1064 42000 …near 'blob(1,2)' at line 1
> exec ^
: string
| SELECT blob(1,2,3)
< error
| 1064 42000 …near 'blob(1,2,3)' at line 1
> exec ^
: string
| SELECT blob(1,2,3,4)
< error
| 1064 42000 …near 'blob(1,2,3,4)' at line 1
> exec ^
: string
| SELECT blob(*)
< error
| 1064 42000 …near 'blob(*)' at line 1
> exec ^
: string
| SELECT blob('a')
< error
| 1064 42000 …near 'blob('a')' at line 1
> exec ^
: string
| SELECT blob(a USING utf8mb4)
< error
| 1064 42000 …near 'blob(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT blob(1 AS SIGNED)
< error
| 1064 42000 …near 'blob(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT blob(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'blob(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT blob(1 FROM 2)
< error
| 1064 42000 …near 'blob(1 FROM 2)' at line 1
> exec ^
: string
| SELECT blob(1 IN 2)
< error
| 1064 42000 …near 'blob(1 IN 2)' at line 1
> exec ^
: string
| SELECT blob(DISTINCT 1)
< error
| 1064 42000 …near 'blob(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT blob() OVER ()
< error
| 1064 42000 …near 'blob() OVER ()' at line 1
> exec ^
: string
| SELECT blob(1) OVER ()
< error
| 1064 42000 …near 'blob(1) OVER ()' at line 1
> exec ^
: string
| SELECT blob(1,2) OVER ()
< error
| 1064 42000 …near 'blob(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT blob(1,2,3) OVER ()
< error
| 1064 42000 …near 'blob(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT blob(1) FROM t
< error
| 1064 42000 …near 'blob(1) FROM t' at line 1
> exec ^
: string
| SELECT blob(SELECT 1)
< error
| 1064 42000 …near 'blob(SELECT 1)' at line 1
> exec ^
: string
| SELECT blob(YEAR FROM 1)
< error
| 1064 42000 …near 'blob(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT blob(DAY, 1, 2)
< error
| 1064 42000 …near 'blob(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT blob(a)
< error
| 1064 42000 …near 'blob(a)' at line 1
> exec ^
: string
| SELECT blob(a, b)
< error
| 1064 42000 …near 'blob(a, b)' at line 1
> exec ^
: string
| SELECT blob(a) OVER ()
< error
| 1064 42000 …near 'blob(a) OVER ()' at line 1
> exec ^
: string
| SELECT block()
< error
| 1630 42000 FUNCTION test.block does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT block()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT bool()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT bool(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT bool(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT bool(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT bool(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT bool(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT bool('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT bool(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT bool(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT bool(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT bool(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT bool(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT bool(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT bool() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT bool(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT bool(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT bool(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT bool(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT bool(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT bool(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT bool(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT bool(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT bool(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT bool(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT boolean()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT boolean(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT boolean(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT boolean(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT boolean(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT boolean(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT boolean('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT boolean(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT boolean(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT boolean(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT boolean(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT boolean(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT boolean(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT boolean() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT boolean(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT boolean(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT boolean(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT boolean(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT boolean(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT boolean(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT boolean(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT boolean(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT boolean(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT boolean(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT both()
< error
| 1064 42000 …near 'both()' at line 1
> exec ^
: string
| SELECT both(1)
< error
| 1064 42000 …near 'both(1)' at line 1
> exec ^
: string
| SELECT both(1,2)
< error
| 1064 42000 …near 'both(1,2)' at line 1
> exec ^
: string
| SELECT both(1,2,3)
< error
| 1064 42000 …near 'both(1,2,3)' at line 1
> exec ^
: string
| SELECT both(1,2,3,4)
< error
| 1064 42000 …near 'both(1,2,3,4)' at line 1
> exec ^
: string
| SELECT both(*)
< error
| 1064 42000 …near 'both(*)' at line 1
> exec ^
: string
| SELECT both('a')
< error
| 1064 42000 …near 'both('a')' at line 1
> exec ^
: string
| SELECT both(a USING utf8mb4)
< error
| 1064 42000 …near 'both(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT both(1 AS SIGNED)
< error
| 1064 42000 …near 'both(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT both(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'both(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT both(1 FROM 2)
< error
| 1064 42000 …near 'both(1 FROM 2)' at line 1
> exec ^
: string
| SELECT both(1 IN 2)
< error
| 1064 42000 …near 'both(1 IN 2)' at line 1
> exec ^
: string
| SELECT both(DISTINCT 1)
< error
| 1064 42000 …near 'both(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT both() OVER ()
< error
| 1064 42000 …near 'both() OVER ()' at line 1
> exec ^
: string
| SELECT both(1) OVER ()
< error
| 1064 42000 …near 'both(1) OVER ()' at line 1
> exec ^
: string
| SELECT both(1,2) OVER ()
< error
| 1064 42000 …near 'both(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT both(1,2,3) OVER ()
< error
| 1064 42000 …near 'both(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT both(1) FROM t
< error
| 1064 42000 …near 'both(1) FROM t' at line 1
> exec ^
: string
| SELECT both(SELECT 1)
< error
| 1064 42000 …near 'both(SELECT 1)' at line 1
> exec ^
: string
| SELECT both(YEAR FROM 1)
< error
| 1064 42000 …near 'both(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT both(DAY, 1, 2)
< error
| 1064 42000 …near 'both(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT both(a)
< error
| 1064 42000 …near 'both(a)' at line 1
> exec ^
: string
| SELECT both(a, b)
< error
| 1064 42000 …near 'both(a, b)' at line 1
> exec ^
: string
| SELECT both(a) OVER ()
< error
| 1064 42000 …near 'both(a) OVER ()' at line 1
> exec ^
: string
| SELECT btree()
< error
| 1630 42000 FUNCTION test.btree does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT btree()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT buckets()
< error
| 1305 42000 FUNCTION test.buckets does not exist
> exec ^
: string
| SELECT buckets()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT bulk()
< error
| 1305 42000 FUNCTION test.bulk does not exist
> exec ^
: string
| SELECT bulk()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT by()
< error
| 1064 42000 …near 'by()' at line 1
> exec ^
: string
| SELECT by(1)
< error
| 1064 42000 …near 'by(1)' at line 1
> exec ^
: string
| SELECT by(1,2)
< error
| 1064 42000 …near 'by(1,2)' at line 1
> exec ^
: string
| SELECT by(1,2,3)
< error
| 1064 42000 …near 'by(1,2,3)' at line 1
> exec ^
: string
| SELECT by(1,2,3,4)
< error
| 1064 42000 …near 'by(1,2,3,4)' at line 1
> exec ^
: string
| SELECT by(*)
< error
| 1064 42000 …near 'by(*)' at line 1
> exec ^
: string
| SELECT by('a')
< error
| 1064 42000 …near 'by('a')' at line 1
> exec ^
: string
| SELECT by(a USING utf8mb4)
< error
| 1064 42000 …near 'by(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT by(1 AS SIGNED)
< error
| 1064 42000 …near 'by(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT by(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'by(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT by(1 FROM 2)
< error
| 1064 42000 …near 'by(1 FROM 2)' at line 1
> exec ^
: string
| SELECT by(1 IN 2)
< error
| 1064 42000 …near 'by(1 IN 2)' at line 1
> exec ^
: string
| SELECT by(DISTINCT 1)
< error
| 1064 42000 …near 'by(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT by() OVER ()
< error
| 1064 42000 …near 'by() OVER ()' at line 1
> exec ^
: string
| SELECT by(1) OVER ()
< error
| 1064 42000 …near 'by(1) OVER ()' at line 1
> exec ^
: string
| SELECT by(1,2) OVER ()
< error
| 1064 42000 …near 'by(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT by(1,2,3) OVER ()
< error
| 1064 42000 …near 'by(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT by(1) FROM t
< error
| 1064 42000 …near 'by(1) FROM t' at line 1
> exec ^
: string
| SELECT by(SELECT 1)
< error
| 1064 42000 …near 'by(SELECT 1)' at line 1
> exec ^
: string
| SELECT by(YEAR FROM 1)
< error
| 1064 42000 …near 'by(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT by(DAY, 1, 2)
< error
| 1064 42000 …near 'by(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT by(a)
< error
| 1064 42000 …near 'by(a)' at line 1
> exec ^
: string
| SELECT by(a, b)
< error
| 1064 42000 …near 'by(a, b)' at line 1
> exec ^
: string
| SELECT by(a) OVER ()
< error
| 1064 42000 …near 'by(a) OVER ()' at line 1
> exec ^
: string
| SELECT byte()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT byte(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT byte(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT byte(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT byte(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT byte(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT byte('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT byte(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT byte(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT byte(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT byte(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT byte(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT byte(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT byte() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT byte(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT byte(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT byte(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT byte(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT byte(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT byte(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT byte(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT byte(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT byte(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT byte(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT cache()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT cache(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT cache(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT cache(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT cache(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT cache(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT cache('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT cache(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cache(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT cache(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT cache(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT cache(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT cache(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT cache() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT cache(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT cache(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT cache(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT cache(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT cache(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT cache(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT cache(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT cache(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT cache(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT cache(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT call()
< error
| 1064 42000 …near 'call()' at line 1
> exec ^
: string
| SELECT call(1)
< error
| 1064 42000 …near 'call(1)' at line 1
> exec ^
: string
| SELECT call(1,2)
< error
| 1064 42000 …near 'call(1,2)' at line 1
> exec ^
: string
| SELECT call(1,2,3)
< error
| 1064 42000 …near 'call(1,2,3)' at line 1
> exec ^
: string
| SELECT call(1,2,3,4)
< error
| 1064 42000 …near 'call(1,2,3,4)' at line 1
> exec ^
: string
| SELECT call(*)
< error
| 1064 42000 …near 'call(*)' at line 1
> exec ^
: string
| SELECT call('a')
< error
| 1064 42000 …near 'call('a')' at line 1
> exec ^
: string
| SELECT call(a USING utf8mb4)
< error
| 1064 42000 …near 'call(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT call(1 AS SIGNED)
< error
| 1064 42000 …near 'call(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT call(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'call(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT call(1 FROM 2)
< error
| 1064 42000 …near 'call(1 FROM 2)' at line 1
> exec ^
: string
| SELECT call(1 IN 2)
< error
| 1064 42000 …near 'call(1 IN 2)' at line 1
> exec ^
: string
| SELECT call(DISTINCT 1)
< error
| 1064 42000 …near 'call(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT call() OVER ()
< error
| 1064 42000 …near 'call() OVER ()' at line 1
> exec ^
: string
| SELECT call(1) OVER ()
< error
| 1064 42000 …near 'call(1) OVER ()' at line 1
> exec ^
: string
| SELECT call(1,2) OVER ()
< error
| 1064 42000 …near 'call(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT call(1,2,3) OVER ()
< error
| 1064 42000 …near 'call(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT call(1) FROM t
< error
| 1064 42000 …near 'call(1) FROM t' at line 1
> exec ^
: string
| SELECT call(SELECT 1)
< error
| 1064 42000 …near 'call(SELECT 1)' at line 1
> exec ^
: string
| SELECT call(YEAR FROM 1)
< error
| 1064 42000 …near 'call(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT call(DAY, 1, 2)
< error
| 1064 42000 …near 'call(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT call(a)
< error
| 1064 42000 …near 'call(a)' at line 1
> exec ^
: string
| SELECT call(a, b)
< error
| 1064 42000 …near 'call(a, b)' at line 1
> exec ^
: string
| SELECT call(a) OVER ()
< error
| 1064 42000 …near 'call(a) OVER ()' at line 1
> exec ^
: string
| SELECT cascade()
< error
| 1064 42000 …near 'cascade()' at line 1
> exec ^
: string
| SELECT cascade(1)
< error
| 1064 42000 …near 'cascade(1)' at line 1
> exec ^
: string
| SELECT cascade(1,2)
< error
| 1064 42000 …near 'cascade(1,2)' at line 1
> exec ^
: string
| SELECT cascade(1,2,3)
< error
| 1064 42000 …near 'cascade(1,2,3)' at line 1
> exec ^
: string
| SELECT cascade(1,2,3,4)
< error
| 1064 42000 …near 'cascade(1,2,3,4)' at line 1
> exec ^
: string
| SELECT cascade(*)
< error
| 1064 42000 …near 'cascade(*)' at line 1
> exec ^
: string
| SELECT cascade('a')
< error
| 1064 42000 …near 'cascade('a')' at line 1
> exec ^
: string
| SELECT cascade(a USING utf8mb4)
< error
| 1064 42000 …near 'cascade(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cascade(1 AS SIGNED)
< error
| 1064 42000 …near 'cascade(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT cascade(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'cascade(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT cascade(1 FROM 2)
< error
| 1064 42000 …near 'cascade(1 FROM 2)' at line 1
> exec ^
: string
| SELECT cascade(1 IN 2)
< error
| 1064 42000 …near 'cascade(1 IN 2)' at line 1
> exec ^
: string
| SELECT cascade(DISTINCT 1)
< error
| 1064 42000 …near 'cascade(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT cascade() OVER ()
< error
| 1064 42000 …near 'cascade() OVER ()' at line 1
> exec ^
: string
| SELECT cascade(1) OVER ()
< error
| 1064 42000 …near 'cascade(1) OVER ()' at line 1
> exec ^
: string
| SELECT cascade(1,2) OVER ()
< error
| 1064 42000 …near 'cascade(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT cascade(1,2,3) OVER ()
< error
| 1064 42000 …near 'cascade(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT cascade(1) FROM t
< error
| 1064 42000 …near 'cascade(1) FROM t' at line 1
> exec ^
: string
| SELECT cascade(SELECT 1)
< error
| 1064 42000 …near 'cascade(SELECT 1)' at line 1
> exec ^
: string
| SELECT cascade(YEAR FROM 1)
< error
| 1064 42000 …near 'cascade(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT cascade(DAY, 1, 2)
< error
| 1064 42000 …near 'cascade(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT cascade(a)
< error
| 1064 42000 …near 'cascade(a)' at line 1
> exec ^
: string
| SELECT cascade(a, b)
< error
| 1064 42000 …near 'cascade(a, b)' at line 1
> exec ^
: string
| SELECT cascade(a) OVER ()
< error
| 1064 42000 …near 'cascade(a) OVER ()' at line 1
> exec ^
: string
| SELECT cascaded()
< error
| 1630 42000 FUNCTION test.cascaded does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cascaded()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT case()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT case(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(1,2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(1,2,3)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(1,2,3,4)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT case('a')
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT case(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT case(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT case(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT case(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT case(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT case() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT case(1) OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec ^
: string
| SELECT case(1,2) OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec ^
: string
| SELECT case(1,2,3) OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec ^
: string
| SELECT case(1) FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec ^
: string
| SELECT case(SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT case(DAY, 1, 2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(a)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(a, b)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT case(a) OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec ^
: string
| SELECT cast()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT cast(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT cast(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT cast(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT cast(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT cast(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT cast('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT cast(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cast(1 AS SIGNED)
< ok 0 0
> exec ^
: string
| SELECT catalog_name()
< error
| 1630 42000 FUNCTION test.catalog_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT catalog_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT chain()
< error
| 1630 42000 FUNCTION test.chain does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT chain()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT challenge_response()
< error
| 1305 42000 FUNCTION test.challenge_response does not exist
> exec ^
: string
| SELECT challenge_response()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT change()
< error
| 1064 42000 …near 'change()' at line 1
> exec ^
: string
| SELECT change(1)
< error
| 1064 42000 …near 'change(1)' at line 1
> exec ^
: string
| SELECT change(1,2)
< error
| 1064 42000 …near 'change(1,2)' at line 1
> exec ^
: string
| SELECT change(1,2,3)
< error
| 1064 42000 …near 'change(1,2,3)' at line 1
> exec ^
: string
| SELECT change(1,2,3,4)
< error
| 1064 42000 …near 'change(1,2,3,4)' at line 1
> exec ^
: string
| SELECT change(*)
< error
| 1064 42000 …near 'change(*)' at line 1
> exec ^
: string
| SELECT change('a')
< error
| 1064 42000 …near 'change('a')' at line 1
> exec ^
: string
| SELECT change(a USING utf8mb4)
< error
| 1064 42000 …near 'change(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT change(1 AS SIGNED)
< error
| 1064 42000 …near 'change(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT change(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'change(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT change(1 FROM 2)
< error
| 1064 42000 …near 'change(1 FROM 2)' at line 1
> exec ^
: string
| SELECT change(1 IN 2)
< error
| 1064 42000 …near 'change(1 IN 2)' at line 1
> exec ^
: string
| SELECT change(DISTINCT 1)
< error
| 1064 42000 …near 'change(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT change() OVER ()
< error
| 1064 42000 …near 'change() OVER ()' at line 1
> exec ^
: string
| SELECT change(1) OVER ()
< error
| 1064 42000 …near 'change(1) OVER ()' at line 1
> exec ^
: string
| SELECT change(1,2) OVER ()
< error
| 1064 42000 …near 'change(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT change(1,2,3) OVER ()
< error
| 1064 42000 …near 'change(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT change(1) FROM t
< error
| 1064 42000 …near 'change(1) FROM t' at line 1
> exec ^
: string
| SELECT change(SELECT 1)
< error
| 1064 42000 …near 'change(SELECT 1)' at line 1
> exec ^
: string
| SELECT change(YEAR FROM 1)
< error
| 1064 42000 …near 'change(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT change(DAY, 1, 2)
< error
| 1064 42000 …near 'change(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT change(a)
< error
| 1064 42000 …near 'change(a)' at line 1
> exec ^
: string
| SELECT change(a, b)
< error
| 1064 42000 …near 'change(a, b)' at line 1
> exec ^
: string
| SELECT change(a) OVER ()
< error
| 1064 42000 …near 'change(a) OVER ()' at line 1
> exec ^
: string
| SELECT changed()
< error
| 1630 42000 FUNCTION test.changed does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT changed()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT channel()
< error
| 1630 42000 FUNCTION test.channel does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT channel()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT char()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT char(1)
< ok 0 0
> exec ^
: string
| SELECT character()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT character(1)
< ok 0 0
> exec ^
: string
| SELECT charset()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT charset(1)
< ok 0 0
> exec ^
: string
| SELECT check()
< error
| 1064 42000 …near 'check()' at line 1
> exec ^
: string
| SELECT check(1)
< error
| 1064 42000 …near 'check(1)' at line 1
> exec ^
: string
| SELECT check(1,2)
< error
| 1064 42000 …near 'check(1,2)' at line 1
> exec ^
: string
| SELECT check(1,2,3)
< error
| 1064 42000 …near 'check(1,2,3)' at line 1
> exec ^
: string
| SELECT check(1,2,3,4)
< error
| 1064 42000 …near 'check(1,2,3,4)' at line 1
> exec ^
: string
| SELECT check(*)
< error
| 1064 42000 …near 'check(*)' at line 1
> exec ^
: string
| SELECT check('a')
< error
| 1064 42000 …near 'check('a')' at line 1
> exec ^
: string
| SELECT check(a USING utf8mb4)
< error
| 1064 42000 …near 'check(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT check(1 AS SIGNED)
< error
| 1064 42000 …near 'check(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT check(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'check(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT check(1 FROM 2)
< error
| 1064 42000 …near 'check(1 FROM 2)' at line 1
> exec ^
: string
| SELECT check(1 IN 2)
< error
| 1064 42000 …near 'check(1 IN 2)' at line 1
> exec ^
: string
| SELECT check(DISTINCT 1)
< error
| 1064 42000 …near 'check(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT check() OVER ()
< error
| 1064 42000 …near 'check() OVER ()' at line 1
> exec ^
: string
| SELECT check(1) OVER ()
< error
| 1064 42000 …near 'check(1) OVER ()' at line 1
> exec ^
: string
| SELECT check(1,2) OVER ()
< error
| 1064 42000 …near 'check(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT check(1,2,3) OVER ()
< error
| 1064 42000 …near 'check(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT check(1) FROM t
< error
| 1064 42000 …near 'check(1) FROM t' at line 1
> exec ^
: string
| SELECT check(SELECT 1)
< error
| 1064 42000 …near 'check(SELECT 1)' at line 1
> exec ^
: string
| SELECT check(YEAR FROM 1)
< error
| 1064 42000 …near 'check(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT check(DAY, 1, 2)
< error
| 1064 42000 …near 'check(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT check(a)
< error
| 1064 42000 …near 'check(a)' at line 1
> exec ^
: string
| SELECT check(a, b)
< error
| 1064 42000 …near 'check(a, b)' at line 1
> exec ^
: string
| SELECT check(a) OVER ()
< error
| 1064 42000 …near 'check(a) OVER ()' at line 1
> exec ^
: string
| SELECT checksum()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT checksum(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT checksum(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT checksum(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT checksum(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT checksum(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT checksum('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT checksum(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT checksum(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT checksum(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT checksum(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT checksum(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT checksum(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT checksum() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT checksum(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT checksum(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT checksum(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT checksum(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT checksum(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT checksum(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT checksum(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT checksum(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT checksum(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT checksum(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT cipher()
< error
| 1630 42000 FUNCTION test.cipher does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cipher()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT class_origin()
< error
| 1630 42000 FUNCTION test.class_origin does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT class_origin()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT client()
< error
| 1630 42000 FUNCTION test.client does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT client()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT clone()
< error
| 1305 42000 FUNCTION test.clone does not exist
> exec ^
: string
| SELECT clone()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT close()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT close(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT close(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT close(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT close(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT close(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT close('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT close(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT close(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT close(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT close(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT close(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT close(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT close() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT close(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT close(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT close(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT close(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT close(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT close(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT close(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT close(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT close(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT close(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT coalesce()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'coalesce'
> exec ^
: string
| SELECT coalesce()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'coalesce'
> exec ^
: string
| SELECT coalesce(1)
< ok 0 0
> exec ^
: string
| SELECT code()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT code(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT code(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT code(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT code(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT code(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT code('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT code(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT code(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT code(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT code(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT code(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT code(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT code() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT code(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT code(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT code(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT code(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT code(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT code(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT code(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT code(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT code(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT code(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT collate()
< error
| 1064 42000 …near 'collate()' at line 1
> exec ^
: string
| SELECT collate(1)
< error
| 1064 42000 …near 'collate(1)' at line 1
> exec ^
: string
| SELECT collate(1,2)
< error
| 1064 42000 …near 'collate(1,2)' at line 1
> exec ^
: string
| SELECT collate(1,2,3)
< error
| 1064 42000 …near 'collate(1,2,3)' at line 1
> exec ^
: string
| SELECT collate(1,2,3,4)
< error
| 1064 42000 …near 'collate(1,2,3,4)' at line 1
> exec ^
: string
| SELECT collate(*)
< error
| 1064 42000 …near 'collate(*)' at line 1
> exec ^
: string
| SELECT collate('a')
< error
| 1064 42000 …near 'collate('a')' at line 1
> exec ^
: string
| SELECT collate(a USING utf8mb4)
< error
| 1064 42000 …near 'collate(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT collate(1 AS SIGNED)
< error
| 1064 42000 …near 'collate(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT collate(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'collate(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT collate(1 FROM 2)
< error
| 1064 42000 …near 'collate(1 FROM 2)' at line 1
> exec ^
: string
| SELECT collate(1 IN 2)
< error
| 1064 42000 …near 'collate(1 IN 2)' at line 1
> exec ^
: string
| SELECT collate(DISTINCT 1)
< error
| 1064 42000 …near 'collate(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT collate() OVER ()
< error
| 1064 42000 …near 'collate() OVER ()' at line 1
> exec ^
: string
| SELECT collate(1) OVER ()
< error
| 1064 42000 …near 'collate(1) OVER ()' at line 1
> exec ^
: string
| SELECT collate(1,2) OVER ()
< error
| 1064 42000 …near 'collate(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT collate(1,2,3) OVER ()
< error
| 1064 42000 …near 'collate(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT collate(1) FROM t
< error
| 1064 42000 …near 'collate(1) FROM t' at line 1
> exec ^
: string
| SELECT collate(SELECT 1)
< error
| 1064 42000 …near 'collate(SELECT 1)' at line 1
> exec ^
: string
| SELECT collate(YEAR FROM 1)
< error
| 1064 42000 …near 'collate(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT collate(DAY, 1, 2)
< error
| 1064 42000 …near 'collate(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT collate(a)
< error
| 1064 42000 …near 'collate(a)' at line 1
> exec ^
: string
| SELECT collate(a, b)
< error
| 1064 42000 …near 'collate(a, b)' at line 1
> exec ^
: string
| SELECT collate(a) OVER ()
< error
| 1064 42000 …near 'collate(a) OVER ()' at line 1
> exec ^
: string
| SELECT collation()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'collation'
> exec ^
: string
| SELECT collation()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'collation'
> exec ^
: string
| SELECT collation(1)
< ok 0 0
> exec ^
: string
| SELECT column()
< error
| 1064 42000 …near 'column()' at line 1
> exec ^
: string
| SELECT column(1)
< error
| 1064 42000 …near 'column(1)' at line 1
> exec ^
: string
| SELECT column(1,2)
< error
| 1064 42000 …near 'column(1,2)' at line 1
> exec ^
: string
| SELECT column(1,2,3)
< error
| 1064 42000 …near 'column(1,2,3)' at line 1
> exec ^
: string
| SELECT column(1,2,3,4)
< error
| 1064 42000 …near 'column(1,2,3,4)' at line 1
> exec ^
: string
| SELECT column(*)
< error
| 1064 42000 …near 'column(*)' at line 1
> exec ^
: string
| SELECT column('a')
< error
| 1064 42000 …near 'column('a')' at line 1
> exec ^
: string
| SELECT column(a USING utf8mb4)
< error
| 1064 42000 …near 'column(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT column(1 AS SIGNED)
< error
| 1064 42000 …near 'column(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT column(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'column(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT column(1 FROM 2)
< error
| 1064 42000 …near 'column(1 FROM 2)' at line 1
> exec ^
: string
| SELECT column(1 IN 2)
< error
| 1064 42000 …near 'column(1 IN 2)' at line 1
> exec ^
: string
| SELECT column(DISTINCT 1)
< error
| 1064 42000 …near 'column(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT column() OVER ()
< error
| 1064 42000 …near 'column() OVER ()' at line 1
> exec ^
: string
| SELECT column(1) OVER ()
< error
| 1064 42000 …near 'column(1) OVER ()' at line 1
> exec ^
: string
| SELECT column(1,2) OVER ()
< error
| 1064 42000 …near 'column(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT column(1,2,3) OVER ()
< error
| 1064 42000 …near 'column(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT column(1) FROM t
< error
| 1064 42000 …near 'column(1) FROM t' at line 1
> exec ^
: string
| SELECT column(SELECT 1)
< error
| 1064 42000 …near 'column(SELECT 1)' at line 1
> exec ^
: string
| SELECT column(YEAR FROM 1)
< error
| 1064 42000 …near 'column(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT column(DAY, 1, 2)
< error
| 1064 42000 …near 'column(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT column(a)
< error
| 1064 42000 …near 'column(a)' at line 1
> exec ^
: string
| SELECT column(a, b)
< error
| 1064 42000 …near 'column(a, b)' at line 1
> exec ^
: string
| SELECT column(a) OVER ()
< error
| 1064 42000 …near 'column(a) OVER ()' at line 1
> exec ^
: string
| SELECT column_format()
< error
| 1305 42000 FUNCTION test.column_format does not exist
> exec ^
: string
| SELECT column_format()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT column_name()
< error
| 1630 42000 FUNCTION test.column_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT column_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT columns()
< error
| 1630 42000 FUNCTION test.columns does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT columns()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT comment()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT comment(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT comment(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT comment(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT comment(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT comment('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT comment(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT comment(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT comment(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT comment(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT comment(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT comment(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT comment() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT comment(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT comment(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT comment(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT comment(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT comment(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT comment(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT comment(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT comment(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT comment(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT comment(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT commit()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT commit(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT commit(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT commit(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT commit(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT commit(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT commit('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT commit(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT commit(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT commit(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT commit(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT commit(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT commit(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT commit() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT commit(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT commit(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT commit(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT commit(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT commit(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT commit(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT commit(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT commit(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT commit(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT commit(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT committed()
< error
| 1630 42000 FUNCTION test.committed does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT committed()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT compact()
< error
| 1630 42000 FUNCTION test.compact does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT compact()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT completion()
< error
| 1630 42000 FUNCTION test.completion does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT completion()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT component()
< error
| 1305 42000 FUNCTION test.component does not exist
> exec ^
: string
| SELECT component()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT compressed()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT compressed(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT compressed(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT compressed(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT compressed(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT compressed(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT compressed('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT compressed(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT compressed(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT compressed(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT compressed(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT compressed(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT compressed(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT compressed() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT compressed(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT compressed(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT compressed(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT compressed(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT compressed(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT compressed(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT compressed(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT compressed(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT compressed(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT compressed(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT compression()
< error
| 1305 42000 FUNCTION test.compression does not exist
> exec ^
: string
| SELECT compression()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT concat()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'concat'
> exec ^
: string
| SELECT concat()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'concat'
> exec ^
: string
| SELECT concat(1)
< ok 0 0
> exec ^
: string
| SELECT concurrent()
< error
| 1630 42000 FUNCTION test.concurrent does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT concurrent()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT condition()
< error
| 1064 42000 …near 'condition()' at line 1
> exec ^
: string
| SELECT condition(1)
< error
| 1064 42000 …near 'condition(1)' at line 1
> exec ^
: string
| SELECT condition(1,2)
< error
| 1064 42000 …near 'condition(1,2)' at line 1
> exec ^
: string
| SELECT condition(1,2,3)
< error
| 1064 42000 …near 'condition(1,2,3)' at line 1
> exec ^
: string
| SELECT condition(1,2,3,4)
< error
| 1064 42000 …near 'condition(1,2,3,4)' at line 1
> exec ^
: string
| SELECT condition(*)
< error
| 1064 42000 …near 'condition(*)' at line 1
> exec ^
: string
| SELECT condition('a')
< error
| 1064 42000 …near 'condition('a')' at line 1
> exec ^
: string
| SELECT condition(a USING utf8mb4)
< error
| 1064 42000 …near 'condition(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT condition(1 AS SIGNED)
< error
| 1064 42000 …near 'condition(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT condition(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'condition(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT condition(1 FROM 2)
< error
| 1064 42000 …near 'condition(1 FROM 2)' at line 1
> exec ^
: string
| SELECT condition(1 IN 2)
< error
| 1064 42000 …near 'condition(1 IN 2)' at line 1
> exec ^
: string
| SELECT condition(DISTINCT 1)
< error
| 1064 42000 …near 'condition(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT condition() OVER ()
< error
| 1064 42000 …near 'condition() OVER ()' at line 1
> exec ^
: string
| SELECT condition(1) OVER ()
< error
| 1064 42000 …near 'condition(1) OVER ()' at line 1
> exec ^
: string
| SELECT condition(1,2) OVER ()
< error
| 1064 42000 …near 'condition(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT condition(1,2,3) OVER ()
< error
| 1064 42000 …near 'condition(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT condition(1) FROM t
< error
| 1064 42000 …near 'condition(1) FROM t' at line 1
> exec ^
: string
| SELECT condition(SELECT 1)
< error
| 1064 42000 …near 'condition(SELECT 1)' at line 1
> exec ^
: string
| SELECT condition(YEAR FROM 1)
< error
| 1064 42000 …near 'condition(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT condition(DAY, 1, 2)
< error
| 1064 42000 …near 'condition(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT condition(a)
< error
| 1064 42000 …near 'condition(a)' at line 1
> exec ^
: string
| SELECT condition(a, b)
< error
| 1064 42000 …near 'condition(a, b)' at line 1
> exec ^
: string
| SELECT condition(a) OVER ()
< error
| 1064 42000 …near 'condition(a) OVER ()' at line 1
> exec ^
: string
| SELECT connection()
< error
| 1630 42000 FUNCTION test.connection does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT connection()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT consistent()
< error
| 1630 42000 FUNCTION test.consistent does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT consistent()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT constraint()
< error
| 1064 42000 …near 'constraint()' at line 1
> exec ^
: string
| SELECT constraint(1)
< error
| 1064 42000 …near 'constraint(1)' at line 1
> exec ^
: string
| SELECT constraint(1,2)
< error
| 1064 42000 …near 'constraint(1,2)' at line 1
> exec ^
: string
| SELECT constraint(1,2,3)
< error
| 1064 42000 …near 'constraint(1,2,3)' at line 1
> exec ^
: string
| SELECT constraint(1,2,3,4)
< error
| 1064 42000 …near 'constraint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT constraint(*)
< error
| 1064 42000 …near 'constraint(*)' at line 1
> exec ^
: string
| SELECT constraint('a')
< error
| 1064 42000 …near 'constraint('a')' at line 1
> exec ^
: string
| SELECT constraint(a USING utf8mb4)
< error
| 1064 42000 …near 'constraint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT constraint(1 AS SIGNED)
< error
| 1064 42000 …near 'constraint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT constraint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'constraint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT constraint(1 FROM 2)
< error
| 1064 42000 …near 'constraint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT constraint(1 IN 2)
< error
| 1064 42000 …near 'constraint(1 IN 2)' at line 1
> exec ^
: string
| SELECT constraint(DISTINCT 1)
< error
| 1064 42000 …near 'constraint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT constraint() OVER ()
< error
| 1064 42000 …near 'constraint() OVER ()' at line 1
> exec ^
: string
| SELECT constraint(1) OVER ()
< error
| 1064 42000 …near 'constraint(1) OVER ()' at line 1
> exec ^
: string
| SELECT constraint(1,2) OVER ()
< error
| 1064 42000 …near 'constraint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT constraint(1,2,3) OVER ()
< error
| 1064 42000 …near 'constraint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT constraint(1) FROM t
< error
| 1064 42000 …near 'constraint(1) FROM t' at line 1
> exec ^
: string
| SELECT constraint(SELECT 1)
< error
| 1064 42000 …near 'constraint(SELECT 1)' at line 1
> exec ^
: string
| SELECT constraint(YEAR FROM 1)
< error
| 1064 42000 …near 'constraint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT constraint(DAY, 1, 2)
< error
| 1064 42000 …near 'constraint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT constraint(a)
< error
| 1064 42000 …near 'constraint(a)' at line 1
> exec ^
: string
| SELECT constraint(a, b)
< error
| 1064 42000 …near 'constraint(a, b)' at line 1
> exec ^
: string
| SELECT constraint(a) OVER ()
< error
| 1064 42000 …near 'constraint(a) OVER ()' at line 1
> exec ^
: string
| SELECT constraint_catalog()
< error
| 1630 42000 FUNCTION test.constraint_catalog does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT constraint_catalog()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT constraint_name()
< error
| 1630 42000 FUNCTION test.constraint_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT constraint_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT constraint_schema()
< error
| 1630 42000 FUNCTION test.constraint_schema does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT constraint_schema()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT contains()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'contains()'
> exec ^
: string
| SELECT contains()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'contains()
| ; SELECT FROM'
> exec ^
: string
| SELECT contains(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'contains(1)'
> exec ^
: string
| SELECT contains(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'contains(1)
| ; SELECT FROM'
> exec ^
: string
| SELECT contains(1,2)
< error
| 4079 HY000 Illegal parameter data type int for operation 'st_contains'
> exec ^
: string
| SELECT contains(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT context()
< error
| 1630 42000 FUNCTION test.context does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT context()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT continue()
< error
| 1064 42000 …near 'continue()' at line 1
> exec ^
: string
| SELECT continue(1)
< error
| 1064 42000 …near 'continue(1)' at line 1
> exec ^
: string
| SELECT continue(1,2)
< error
| 1064 42000 …near 'continue(1,2)' at line 1
> exec ^
: string
| SELECT continue(1,2,3)
< error
| 1064 42000 …near 'continue(1,2,3)' at line 1
> exec ^
: string
| SELECT continue(1,2,3,4)
< error
| 1064 42000 …near 'continue(1,2,3,4)' at line 1
> exec ^
: string
| SELECT continue(*)
< error
| 1064 42000 …near 'continue(*)' at line 1
> exec ^
: string
| SELECT continue('a')
< error
| 1064 42000 …near 'continue('a')' at line 1
> exec ^
: string
| SELECT continue(a USING utf8mb4)
< error
| 1064 42000 …near 'continue(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT continue(1 AS SIGNED)
< error
| 1064 42000 …near 'continue(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT continue(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'continue(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT continue(1 FROM 2)
< error
| 1064 42000 …near 'continue(1 FROM 2)' at line 1
> exec ^
: string
| SELECT continue(1 IN 2)
< error
| 1064 42000 …near 'continue(1 IN 2)' at line 1
> exec ^
: string
| SELECT continue(DISTINCT 1)
< error
| 1064 42000 …near 'continue(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT continue() OVER ()
< error
| 1064 42000 …near 'continue() OVER ()' at line 1
> exec ^
: string
| SELECT continue(1) OVER ()
< error
| 1064 42000 …near 'continue(1) OVER ()' at line 1
> exec ^
: string
| SELECT continue(1,2) OVER ()
< error
| 1064 42000 …near 'continue(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT continue(1,2,3) OVER ()
< error
| 1064 42000 …near 'continue(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT continue(1) FROM t
< error
| 1064 42000 …near 'continue(1) FROM t' at line 1
> exec ^
: string
| SELECT continue(SELECT 1)
< error
| 1064 42000 …near 'continue(SELECT 1)' at line 1
> exec ^
: string
| SELECT continue(YEAR FROM 1)
< error
| 1064 42000 …near 'continue(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT continue(DAY, 1, 2)
< error
| 1064 42000 …near 'continue(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT continue(a)
< error
| 1064 42000 …near 'continue(a)' at line 1
> exec ^
: string
| SELECT continue(a, b)
< error
| 1064 42000 …near 'continue(a, b)' at line 1
> exec ^
: string
| SELECT continue(a) OVER ()
< error
| 1064 42000 …near 'continue(a) OVER ()' at line 1
> exec ^
: string
| SELECT convert()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT convert(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT convert(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT convert(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT convert(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT convert(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT convert('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT convert(a USING utf8mb4)
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec ^
: string
| SELECT convert(a USING utf8mb4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT count()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT count(1)
< ok 0 0
> exec ^
: string
| SELECT cpu()
< error
| 1630 42000 FUNCTION test.cpu does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cpu()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT create()
< error
| 1064 42000 …near 'create()' at line 1
> exec ^
: string
| SELECT create(1)
< error
| 1064 42000 …near 'create(1)' at line 1
> exec ^
: string
| SELECT create(1,2)
< error
| 1064 42000 …near 'create(1,2)' at line 1
> exec ^
: string
| SELECT create(1,2,3)
< error
| 1064 42000 …near 'create(1,2,3)' at line 1
> exec ^
: string
| SELECT create(1,2,3,4)
< error
| 1064 42000 …near 'create(1,2,3,4)' at line 1
> exec ^
: string
| SELECT create(*)
< error
| 1064 42000 …near 'create(*)' at line 1
> exec ^
: string
| SELECT create('a')
< error
| 1064 42000 …near 'create('a')' at line 1
> exec ^
: string
| SELECT create(a USING utf8mb4)
< error
| 1064 42000 …near 'create(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT create(1 AS SIGNED)
< error
| 1064 42000 …near 'create(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT create(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'create(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT create(1 FROM 2)
< error
| 1064 42000 …near 'create(1 FROM 2)' at line 1
> exec ^
: string
| SELECT create(1 IN 2)
< error
| 1064 42000 …near 'create(1 IN 2)' at line 1
> exec ^
: string
| SELECT create(DISTINCT 1)
< error
| 1064 42000 …near 'create(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT create() OVER ()
< error
| 1064 42000 …near 'create() OVER ()' at line 1
> exec ^
: string
| SELECT create(1) OVER ()
< error
| 1064 42000 …near 'create(1) OVER ()' at line 1
> exec ^
: string
| SELECT create(1,2) OVER ()
< error
| 1064 42000 …near 'create(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT create(1,2,3) OVER ()
< error
| 1064 42000 …near 'create(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT create(1) FROM t
< error
| 1064 42000 …near 'create(1) FROM t' at line 1
> exec ^
: string
| SELECT create(SELECT 1)
< error
| 1064 42000 …near 'create(SELECT 1)' at line 1
> exec ^
: string
| SELECT create(YEAR FROM 1)
< error
| 1064 42000 …near 'create(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT create(DAY, 1, 2)
< error
| 1064 42000 …near 'create(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT create(a)
< error
| 1064 42000 …near 'create(a)' at line 1
> exec ^
: string
| SELECT create(a, b)
< error
| 1064 42000 …near 'create(a, b)' at line 1
> exec ^
: string
| SELECT create(a) OVER ()
< error
| 1064 42000 …near 'create(a) OVER ()' at line 1
> exec ^
: string
| SELECT cross()
< error
| 1064 42000 …near 'cross()' at line 1
> exec ^
: string
| SELECT cross(1)
< error
| 1064 42000 …near 'cross(1)' at line 1
> exec ^
: string
| SELECT cross(1,2)
< error
| 1064 42000 …near 'cross(1,2)' at line 1
> exec ^
: string
| SELECT cross(1,2,3)
< error
| 1064 42000 …near 'cross(1,2,3)' at line 1
> exec ^
: string
| SELECT cross(1,2,3,4)
< error
| 1064 42000 …near 'cross(1,2,3,4)' at line 1
> exec ^
: string
| SELECT cross(*)
< error
| 1064 42000 …near 'cross(*)' at line 1
> exec ^
: string
| SELECT cross('a')
< error
| 1064 42000 …near 'cross('a')' at line 1
> exec ^
: string
| SELECT cross(a USING utf8mb4)
< error
| 1064 42000 …near 'cross(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cross(1 AS SIGNED)
< error
| 1064 42000 …near 'cross(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT cross(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'cross(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT cross(1 FROM 2)
< error
| 1064 42000 …near 'cross(1 FROM 2)' at line 1
> exec ^
: string
| SELECT cross(1 IN 2)
< error
| 1064 42000 …near 'cross(1 IN 2)' at line 1
> exec ^
: string
| SELECT cross(DISTINCT 1)
< error
| 1064 42000 …near 'cross(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT cross() OVER ()
< error
| 1064 42000 …near 'cross() OVER ()' at line 1
> exec ^
: string
| SELECT cross(1) OVER ()
< error
| 1064 42000 …near 'cross(1) OVER ()' at line 1
> exec ^
: string
| SELECT cross(1,2) OVER ()
< error
| 1064 42000 …near 'cross(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT cross(1,2,3) OVER ()
< error
| 1064 42000 …near 'cross(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT cross(1) FROM t
< error
| 1064 42000 …near 'cross(1) FROM t' at line 1
> exec ^
: string
| SELECT cross(SELECT 1)
< error
| 1064 42000 …near 'cross(SELECT 1)' at line 1
> exec ^
: string
| SELECT cross(YEAR FROM 1)
< error
| 1064 42000 …near 'cross(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT cross(DAY, 1, 2)
< error
| 1064 42000 …near 'cross(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT cross(a)
< error
| 1064 42000 …near 'cross(a)' at line 1
> exec ^
: string
| SELECT cross(a, b)
< error
| 1064 42000 …near 'cross(a, b)' at line 1
> exec ^
: string
| SELECT cross(a) OVER ()
< error
| 1064 42000 …near 'cross(a) OVER ()' at line 1
> exec ^
: string
| SELECT cube()
< error
| 1630 42000 FUNCTION test.cube does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cube()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT cume_dist()
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT cume_dist(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT cume_dist(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT cume_dist(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT cume_dist(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT cume_dist(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT cume_dist('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT cume_dist(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cume_dist(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT cume_dist(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT cume_dist(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT cume_dist(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT cume_dist(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT cume_dist() OVER ()
< error
| 4018 HY000 No order list in window specification for 'cume_dist'
> exec ^
: string
| SELECT cume_dist() OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT curdate()
< ok 0 0
> exec ^
: string
| SELECT current()
< error
| 1630 42000 FUNCTION test.current does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT current()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT current_date()
< ok 0 0
> exec ^
: string
| SELECT current_time()
< ok 0 0
> exec ^
: string
| SELECT current_timestamp()
< ok 0 0
> exec ^
: string
| SELECT current_user()
< ok 0 0
> exec ^
: string
| SELECT cursor()
< error
| 1064 42000 …near 'cursor()' at line 1
> exec ^
: string
| SELECT cursor(1)
< error
| 1064 42000 …near 'cursor(1)' at line 1
> exec ^
: string
| SELECT cursor(1,2)
< error
| 1064 42000 …near 'cursor(1,2)' at line 1
> exec ^
: string
| SELECT cursor(1,2,3)
< error
| 1064 42000 …near 'cursor(1,2,3)' at line 1
> exec ^
: string
| SELECT cursor(1,2,3,4)
< error
| 1064 42000 …near 'cursor(1,2,3,4)' at line 1
> exec ^
: string
| SELECT cursor(*)
< error
| 1064 42000 …near 'cursor(*)' at line 1
> exec ^
: string
| SELECT cursor('a')
< error
| 1064 42000 …near 'cursor('a')' at line 1
> exec ^
: string
| SELECT cursor(a USING utf8mb4)
< error
| 1064 42000 …near 'cursor(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT cursor(1 AS SIGNED)
< error
| 1064 42000 …near 'cursor(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT cursor(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'cursor(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT cursor(1 FROM 2)
< error
| 1064 42000 …near 'cursor(1 FROM 2)' at line 1
> exec ^
: string
| SELECT cursor(1 IN 2)
< error
| 1064 42000 …near 'cursor(1 IN 2)' at line 1
> exec ^
: string
| SELECT cursor(DISTINCT 1)
< error
| 1064 42000 …near 'cursor(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT cursor() OVER ()
< error
| 1064 42000 …near 'cursor() OVER ()' at line 1
> exec ^
: string
| SELECT cursor(1) OVER ()
< error
| 1064 42000 …near 'cursor(1) OVER ()' at line 1
> exec ^
: string
| SELECT cursor(1,2) OVER ()
< error
| 1064 42000 …near 'cursor(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT cursor(1,2,3) OVER ()
< error
| 1064 42000 …near 'cursor(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT cursor(1) FROM t
< error
| 1064 42000 …near 'cursor(1) FROM t' at line 1
> exec ^
: string
| SELECT cursor(SELECT 1)
< error
| 1064 42000 …near 'cursor(SELECT 1)' at line 1
> exec ^
: string
| SELECT cursor(YEAR FROM 1)
< error
| 1064 42000 …near 'cursor(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT cursor(DAY, 1, 2)
< error
| 1064 42000 …near 'cursor(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT cursor(a)
< error
| 1064 42000 …near 'cursor(a)' at line 1
> exec ^
: string
| SELECT cursor(a, b)
< error
| 1064 42000 …near 'cursor(a, b)' at line 1
> exec ^
: string
| SELECT cursor(a) OVER ()
< error
| 1064 42000 …near 'cursor(a) OVER ()' at line 1
> exec ^
: string
| SELECT cursor_name()
< error
| 1630 42000 FUNCTION test.cursor_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cursor_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT curtime()
< ok 0 0
> exec ^
: string
| SELECT cycle()
< error
| 1630 42000 FUNCTION test.cycle does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT cycle()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT data()
< error
| 1630 42000 FUNCTION test.data does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT data()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT database()
< ok 0 0
> exec ^
: string
| SELECT databases()
< error
| 1064 42000 …near 'databases()' at line 1
> exec ^
: string
| SELECT databases(1)
< error
| 1064 42000 …near 'databases(1)' at line 1
> exec ^
: string
| SELECT databases(1,2)
< error
| 1064 42000 …near 'databases(1,2)' at line 1
> exec ^
: string
| SELECT databases(1,2,3)
< error
| 1064 42000 …near 'databases(1,2,3)' at line 1
> exec ^
: string
| SELECT databases(1,2,3,4)
< error
| 1064 42000 …near 'databases(1,2,3,4)' at line 1
> exec ^
: string
| SELECT databases(*)
< error
| 1064 42000 …near 'databases(*)' at line 1
> exec ^
: string
| SELECT databases('a')
< error
| 1064 42000 …near 'databases('a')' at line 1
> exec ^
: string
| SELECT databases(a USING utf8mb4)
< error
| 1064 42000 …near 'databases(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT databases(1 AS SIGNED)
< error
| 1064 42000 …near 'databases(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT databases(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'databases(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT databases(1 FROM 2)
< error
| 1064 42000 …near 'databases(1 FROM 2)' at line 1
> exec ^
: string
| SELECT databases(1 IN 2)
< error
| 1064 42000 …near 'databases(1 IN 2)' at line 1
> exec ^
: string
| SELECT databases(DISTINCT 1)
< error
| 1064 42000 …near 'databases(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT databases() OVER ()
< error
| 1064 42000 …near 'databases() OVER ()' at line 1
> exec ^
: string
| SELECT databases(1) OVER ()
< error
| 1064 42000 …near 'databases(1) OVER ()' at line 1
> exec ^
: string
| SELECT databases(1,2) OVER ()
< error
| 1064 42000 …near 'databases(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT databases(1,2,3) OVER ()
< error
| 1064 42000 …near 'databases(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT databases(1) FROM t
< error
| 1064 42000 …near 'databases(1) FROM t' at line 1
> exec ^
: string
| SELECT databases(SELECT 1)
< error
| 1064 42000 …near 'databases(SELECT 1)' at line 1
> exec ^
: string
| SELECT databases(YEAR FROM 1)
< error
| 1064 42000 …near 'databases(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT databases(DAY, 1, 2)
< error
| 1064 42000 …near 'databases(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT databases(a)
< error
| 1064 42000 …near 'databases(a)' at line 1
> exec ^
: string
| SELECT databases(a, b)
< error
| 1064 42000 …near 'databases(a, b)' at line 1
> exec ^
: string
| SELECT databases(a) OVER ()
< error
| 1064 42000 …near 'databases(a) OVER ()' at line 1
> exec ^
: string
| SELECT datafile()
< error
| 1630 42000 FUNCTION test.datafile does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT datafile()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT date()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date(1)
< ok 0 0
> exec ^
: string
| SELECT date_add()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_add(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_add(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT date_add(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT date_add(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT date_add(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT date_add('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_add(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT date_add(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT date_add(1, INTERVAL 1 DAY)
< ok 0 0
> exec ^
: string
| SELECT date_sub()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_sub(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_sub(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT date_sub(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT date_sub(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT date_sub(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT date_sub('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT date_sub(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT date_sub(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT date_sub(1, INTERVAL 1 DAY)
< ok 0 0
> exec ^
: string
| SELECT datetime()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT datetime(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT datetime(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT datetime(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT datetime(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT datetime(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT datetime('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT datetime(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT datetime(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT datetime(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT datetime(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT datetime(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT datetime(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT datetime() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT datetime(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT datetime(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT datetime(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT datetime(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT datetime(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT datetime(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT datetime(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT datetime(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT datetime(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT datetime(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT day()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT day(1)
< ok 0 0
> exec ^
: string
| SELECT day_hour()
< error
| 1064 42000 …near 'day_hour()' at line 1
> exec ^
: string
| SELECT day_hour(1)
< error
| 1064 42000 …near 'day_hour(1)' at line 1
> exec ^
: string
| SELECT day_hour(1,2)
< error
| 1064 42000 …near 'day_hour(1,2)' at line 1
> exec ^
: string
| SELECT day_hour(1,2,3)
< error
| 1064 42000 …near 'day_hour(1,2,3)' at line 1
> exec ^
: string
| SELECT day_hour(1,2,3,4)
< error
| 1064 42000 …near 'day_hour(1,2,3,4)' at line 1
> exec ^
: string
| SELECT day_hour(*)
< error
| 1064 42000 …near 'day_hour(*)' at line 1
> exec ^
: string
| SELECT day_hour('a')
< error
| 1064 42000 …near 'day_hour('a')' at line 1
> exec ^
: string
| SELECT day_hour(a USING utf8mb4)
< error
| 1064 42000 …near 'day_hour(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT day_hour(1 AS SIGNED)
< error
| 1064 42000 …near 'day_hour(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT day_hour(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'day_hour(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT day_hour(1 FROM 2)
< error
| 1064 42000 …near 'day_hour(1 FROM 2)' at line 1
> exec ^
: string
| SELECT day_hour(1 IN 2)
< error
| 1064 42000 …near 'day_hour(1 IN 2)' at line 1
> exec ^
: string
| SELECT day_hour(DISTINCT 1)
< error
| 1064 42000 …near 'day_hour(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT day_hour() OVER ()
< error
| 1064 42000 …near 'day_hour() OVER ()' at line 1
> exec ^
: string
| SELECT day_hour(1) OVER ()
< error
| 1064 42000 …near 'day_hour(1) OVER ()' at line 1
> exec ^
: string
| SELECT day_hour(1,2) OVER ()
< error
| 1064 42000 …near 'day_hour(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT day_hour(1,2,3) OVER ()
< error
| 1064 42000 …near 'day_hour(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT day_hour(1) FROM t
< error
| 1064 42000 …near 'day_hour(1) FROM t' at line 1
> exec ^
: string
| SELECT day_hour(SELECT 1)
< error
| 1064 42000 …near 'day_hour(SELECT 1)' at line 1
> exec ^
: string
| SELECT day_hour(YEAR FROM 1)
< error
| 1064 42000 …near 'day_hour(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT day_hour(DAY, 1, 2)
< error
| 1064 42000 …near 'day_hour(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT day_hour(a)
< error
| 1064 42000 …near 'day_hour(a)' at line 1
> exec ^
: string
| SELECT day_hour(a, b)
< error
| 1064 42000 …near 'day_hour(a, b)' at line 1
> exec ^
: string
| SELECT day_hour(a) OVER ()
< error
| 1064 42000 …near 'day_hour(a) OVER ()' at line 1
> exec ^
: string
| SELECT day_microsecond()
< error
| 1064 42000 …near 'day_microsecond()' at line 1
> exec ^
: string
| SELECT day_microsecond(1)
< error
| 1064 42000 …near 'day_microsecond(1)' at line 1
> exec ^
: string
| SELECT day_microsecond(1,2)
< error
| 1064 42000 …near 'day_microsecond(1,2)' at line 1
> exec ^
: string
| SELECT day_microsecond(1,2,3)
< error
| 1064 42000 …near 'day_microsecond(1,2,3)' at line 1
> exec ^
: string
| SELECT day_microsecond(1,2,3,4)
< error
| 1064 42000 …near 'day_microsecond(1,2,3,4)' at line 1
> exec ^
: string
| SELECT day_microsecond(*)
< error
| 1064 42000 …near 'day_microsecond(*)' at line 1
> exec ^
: string
| SELECT day_microsecond('a')
< error
| 1064 42000 …near 'day_microsecond('a')' at line 1
> exec ^
: string
| SELECT day_microsecond(a USING utf8mb4)
< error
| 1064 42000 …near 'day_microsecond(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT day_microsecond(1 AS SIGNED)
< error
| 1064 42000 …near 'day_microsecond(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT day_microsecond(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'day_microsecond(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT day_microsecond(1 FROM 2)
< error
| 1064 42000 …near 'day_microsecond(1 FROM 2)' at line 1
> exec ^
: string
| SELECT day_microsecond(1 IN 2)
< error
| 1064 42000 …near 'day_microsecond(1 IN 2)' at line 1
> exec ^
: string
| SELECT day_microsecond(DISTINCT 1)
< error
| 1064 42000 …near 'day_microsecond(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT day_microsecond() OVER ()
< error
| 1064 42000 …near 'day_microsecond() OVER ()' at line 1
> exec ^
: string
| SELECT day_microsecond(1) OVER ()
< error
| 1064 42000 …near 'day_microsecond(1) OVER ()' at line 1
> exec ^
: string
| SELECT day_microsecond(1,2) OVER ()
< error
| 1064 42000 …near 'day_microsecond(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT day_microsecond(1,2,3) OVER ()
< error
| 1064 42000 …near 'day_microsecond(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT day_microsecond(1) FROM t
< error
| 1064 42000 …near 'day_microsecond(1) FROM t' at line 1
> exec ^
: string
| SELECT day_microsecond(SELECT 1)
< error
| 1064 42000 …near 'day_microsecond(SELECT 1)' at line 1
> exec ^
: string
| SELECT day_microsecond(YEAR FROM 1)
< error
| 1064 42000 …near 'day_microsecond(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT day_microsecond(DAY, 1, 2)
< error
| 1064 42000 …near 'day_microsecond(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT day_microsecond(a)
< error
| 1064 42000 …near 'day_microsecond(a)' at line 1
> exec ^
: string
| SELECT day_microsecond(a, b)
< error
| 1064 42000 …near 'day_microsecond(a, b)' at line 1
> exec ^
: string
| SELECT day_microsecond(a) OVER ()
< error
| 1064 42000 …near 'day_microsecond(a) OVER ()' at line 1
> exec ^
: string
| SELECT day_minute()
< error
| 1064 42000 …near 'day_minute()' at line 1
> exec ^
: string
| SELECT day_minute(1)
< error
| 1064 42000 …near 'day_minute(1)' at line 1
> exec ^
: string
| SELECT day_minute(1,2)
< error
| 1064 42000 …near 'day_minute(1,2)' at line 1
> exec ^
: string
| SELECT day_minute(1,2,3)
< error
| 1064 42000 …near 'day_minute(1,2,3)' at line 1
> exec ^
: string
| SELECT day_minute(1,2,3,4)
< error
| 1064 42000 …near 'day_minute(1,2,3,4)' at line 1
> exec ^
: string
| SELECT day_minute(*)
< error
| 1064 42000 …near 'day_minute(*)' at line 1
> exec ^
: string
| SELECT day_minute('a')
< error
| 1064 42000 …near 'day_minute('a')' at line 1
> exec ^
: string
| SELECT day_minute(a USING utf8mb4)
< error
| 1064 42000 …near 'day_minute(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT day_minute(1 AS SIGNED)
< error
| 1064 42000 …near 'day_minute(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT day_minute(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'day_minute(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT day_minute(1 FROM 2)
< error
| 1064 42000 …near 'day_minute(1 FROM 2)' at line 1
> exec ^
: string
| SELECT day_minute(1 IN 2)
< error
| 1064 42000 …near 'day_minute(1 IN 2)' at line 1
> exec ^
: string
| SELECT day_minute(DISTINCT 1)
< error
| 1064 42000 …near 'day_minute(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT day_minute() OVER ()
< error
| 1064 42000 …near 'day_minute() OVER ()' at line 1
> exec ^
: string
| SELECT day_minute(1) OVER ()
< error
| 1064 42000 …near 'day_minute(1) OVER ()' at line 1
> exec ^
: string
| SELECT day_minute(1,2) OVER ()
< error
| 1064 42000 …near 'day_minute(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT day_minute(1,2,3) OVER ()
< error
| 1064 42000 …near 'day_minute(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT day_minute(1) FROM t
< error
| 1064 42000 …near 'day_minute(1) FROM t' at line 1
> exec ^
: string
| SELECT day_minute(SELECT 1)
< error
| 1064 42000 …near 'day_minute(SELECT 1)' at line 1
> exec ^
: string
| SELECT day_minute(YEAR FROM 1)
< error
| 1064 42000 …near 'day_minute(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT day_minute(DAY, 1, 2)
< error
| 1064 42000 …near 'day_minute(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT day_minute(a)
< error
| 1064 42000 …near 'day_minute(a)' at line 1
> exec ^
: string
| SELECT day_minute(a, b)
< error
| 1064 42000 …near 'day_minute(a, b)' at line 1
> exec ^
: string
| SELECT day_minute(a) OVER ()
< error
| 1064 42000 …near 'day_minute(a) OVER ()' at line 1
> exec ^
: string
| SELECT day_second()
< error
| 1064 42000 …near 'day_second()' at line 1
> exec ^
: string
| SELECT day_second(1)
< error
| 1064 42000 …near 'day_second(1)' at line 1
> exec ^
: string
| SELECT day_second(1,2)
< error
| 1064 42000 …near 'day_second(1,2)' at line 1
> exec ^
: string
| SELECT day_second(1,2,3)
< error
| 1064 42000 …near 'day_second(1,2,3)' at line 1
> exec ^
: string
| SELECT day_second(1,2,3,4)
< error
| 1064 42000 …near 'day_second(1,2,3,4)' at line 1
> exec ^
: string
| SELECT day_second(*)
< error
| 1064 42000 …near 'day_second(*)' at line 1
> exec ^
: string
| SELECT day_second('a')
< error
| 1064 42000 …near 'day_second('a')' at line 1
> exec ^
: string
| SELECT day_second(a USING utf8mb4)
< error
| 1064 42000 …near 'day_second(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT day_second(1 AS SIGNED)
< error
| 1064 42000 …near 'day_second(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT day_second(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'day_second(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT day_second(1 FROM 2)
< error
| 1064 42000 …near 'day_second(1 FROM 2)' at line 1
> exec ^
: string
| SELECT day_second(1 IN 2)
< error
| 1064 42000 …near 'day_second(1 IN 2)' at line 1
> exec ^
: string
| SELECT day_second(DISTINCT 1)
< error
| 1064 42000 …near 'day_second(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT day_second() OVER ()
< error
| 1064 42000 …near 'day_second() OVER ()' at line 1
> exec ^
: string
| SELECT day_second(1) OVER ()
< error
| 1064 42000 …near 'day_second(1) OVER ()' at line 1
> exec ^
: string
| SELECT day_second(1,2) OVER ()
< error
| 1064 42000 …near 'day_second(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT day_second(1,2,3) OVER ()
< error
| 1064 42000 …near 'day_second(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT day_second(1) FROM t
< error
| 1064 42000 …near 'day_second(1) FROM t' at line 1
> exec ^
: string
| SELECT day_second(SELECT 1)
< error
| 1064 42000 …near 'day_second(SELECT 1)' at line 1
> exec ^
: string
| SELECT day_second(YEAR FROM 1)
< error
| 1064 42000 …near 'day_second(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT day_second(DAY, 1, 2)
< error
| 1064 42000 …near 'day_second(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT day_second(a)
< error
| 1064 42000 …near 'day_second(a)' at line 1
> exec ^
: string
| SELECT day_second(a, b)
< error
| 1064 42000 …near 'day_second(a, b)' at line 1
> exec ^
: string
| SELECT day_second(a) OVER ()
< error
| 1064 42000 …near 'day_second(a) OVER ()' at line 1
> exec ^
: string
| SELECT deallocate()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT deallocate(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT deallocate(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT deallocate(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT deallocate(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT deallocate(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT deallocate('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT deallocate(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT deallocate(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT deallocate(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT deallocate(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT deallocate(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT deallocate(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT deallocate() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT deallocate(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT deallocate(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT deallocate(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT deallocate(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT deallocate(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT deallocate(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT deallocate(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT deallocate(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT deallocate(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT deallocate(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT dec()
< error
| 1064 42000 …near 'dec()' at line 1
> exec ^
: string
| SELECT dec(1)
< error
| 1064 42000 …near 'dec(1)' at line 1
> exec ^
: string
| SELECT dec(1,2)
< error
| 1064 42000 …near 'dec(1,2)' at line 1
> exec ^
: string
| SELECT dec(1,2,3)
< error
| 1064 42000 …near 'dec(1,2,3)' at line 1
> exec ^
: string
| SELECT dec(1,2,3,4)
< error
| 1064 42000 …near 'dec(1,2,3,4)' at line 1
> exec ^
: string
| SELECT dec(*)
< error
| 1064 42000 …near 'dec(*)' at line 1
> exec ^
: string
| SELECT dec('a')
< error
| 1064 42000 …near 'dec('a')' at line 1
> exec ^
: string
| SELECT dec(a USING utf8mb4)
< error
| 1064 42000 …near 'dec(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT dec(1 AS SIGNED)
< error
| 1064 42000 …near 'dec(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT dec(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'dec(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT dec(1 FROM 2)
< error
| 1064 42000 …near 'dec(1 FROM 2)' at line 1
> exec ^
: string
| SELECT dec(1 IN 2)
< error
| 1064 42000 …near 'dec(1 IN 2)' at line 1
> exec ^
: string
| SELECT dec(DISTINCT 1)
< error
| 1064 42000 …near 'dec(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT dec() OVER ()
< error
| 1064 42000 …near 'dec() OVER ()' at line 1
> exec ^
: string
| SELECT dec(1) OVER ()
< error
| 1064 42000 …near 'dec(1) OVER ()' at line 1
> exec ^
: string
| SELECT dec(1,2) OVER ()
< error
| 1064 42000 …near 'dec(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT dec(1,2,3) OVER ()
< error
| 1064 42000 …near 'dec(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT dec(1) FROM t
< error
| 1064 42000 …near 'dec(1) FROM t' at line 1
> exec ^
: string
| SELECT dec(SELECT 1)
< error
| 1064 42000 …near 'dec(SELECT 1)' at line 1
> exec ^
: string
| SELECT dec(YEAR FROM 1)
< error
| 1064 42000 …near 'dec(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT dec(DAY, 1, 2)
< error
| 1064 42000 …near 'dec(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT dec(a)
< error
| 1064 42000 …near 'dec(a)' at line 1
> exec ^
: string
| SELECT dec(a, b)
< error
| 1064 42000 …near 'dec(a, b)' at line 1
> exec ^
: string
| SELECT dec(a) OVER ()
< error
| 1064 42000 …near 'dec(a) OVER ()' at line 1
> exec ^
: string
| SELECT decimal()
< error
| 1064 42000 …near 'decimal()' at line 1
> exec ^
: string
| SELECT decimal(1)
< error
| 1064 42000 …near 'decimal(1)' at line 1
> exec ^
: string
| SELECT decimal(1,2)
< error
| 1064 42000 …near 'decimal(1,2)' at line 1
> exec ^
: string
| SELECT decimal(1,2,3)
< error
| 1064 42000 …near 'decimal(1,2,3)' at line 1
> exec ^
: string
| SELECT decimal(1,2,3,4)
< error
| 1064 42000 …near 'decimal(1,2,3,4)' at line 1
> exec ^
: string
| SELECT decimal(*)
< error
| 1064 42000 …near 'decimal(*)' at line 1
> exec ^
: string
| SELECT decimal('a')
< error
| 1064 42000 …near 'decimal('a')' at line 1
> exec ^
: string
| SELECT decimal(a USING utf8mb4)
< error
| 1064 42000 …near 'decimal(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT decimal(1 AS SIGNED)
< error
| 1064 42000 …near 'decimal(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT decimal(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'decimal(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT decimal(1 FROM 2)
< error
| 1064 42000 …near 'decimal(1 FROM 2)' at line 1
> exec ^
: string
| SELECT decimal(1 IN 2)
< error
| 1064 42000 …near 'decimal(1 IN 2)' at line 1
> exec ^
: string
| SELECT decimal(DISTINCT 1)
< error
| 1064 42000 …near 'decimal(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT decimal() OVER ()
< error
| 1064 42000 …near 'decimal() OVER ()' at line 1
> exec ^
: string
| SELECT decimal(1) OVER ()
< error
| 1064 42000 …near 'decimal(1) OVER ()' at line 1
> exec ^
: string
| SELECT decimal(1,2) OVER ()
< error
| 1064 42000 …near 'decimal(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT decimal(1,2,3) OVER ()
< error
| 1064 42000 …near 'decimal(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT decimal(1) FROM t
< error
| 1064 42000 …near 'decimal(1) FROM t' at line 1
> exec ^
: string
| SELECT decimal(SELECT 1)
< error
| 1064 42000 …near 'decimal(SELECT 1)' at line 1
> exec ^
: string
| SELECT decimal(YEAR FROM 1)
< error
| 1064 42000 …near 'decimal(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT decimal(DAY, 1, 2)
< error
| 1064 42000 …near 'decimal(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT decimal(a)
< error
| 1064 42000 …near 'decimal(a)' at line 1
> exec ^
: string
| SELECT decimal(a, b)
< error
| 1064 42000 …near 'decimal(a, b)' at line 1
> exec ^
: string
| SELECT decimal(a) OVER ()
< error
| 1064 42000 …near 'decimal(a) OVER ()' at line 1
> exec ^
: string
| SELECT declare()
< error
| 1064 42000 …near 'declare()' at line 1
> exec ^
: string
| SELECT declare(1)
< error
| 1064 42000 …near 'declare(1)' at line 1
> exec ^
: string
| SELECT declare(1,2)
< error
| 1064 42000 …near 'declare(1,2)' at line 1
> exec ^
: string
| SELECT declare(1,2,3)
< error
| 1064 42000 …near 'declare(1,2,3)' at line 1
> exec ^
: string
| SELECT declare(1,2,3,4)
< error
| 1064 42000 …near 'declare(1,2,3,4)' at line 1
> exec ^
: string
| SELECT declare(*)
< error
| 1064 42000 …near 'declare(*)' at line 1
> exec ^
: string
| SELECT declare('a')
< error
| 1064 42000 …near 'declare('a')' at line 1
> exec ^
: string
| SELECT declare(a USING utf8mb4)
< error
| 1064 42000 …near 'declare(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT declare(1 AS SIGNED)
< error
| 1064 42000 …near 'declare(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT declare(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'declare(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT declare(1 FROM 2)
< error
| 1064 42000 …near 'declare(1 FROM 2)' at line 1
> exec ^
: string
| SELECT declare(1 IN 2)
< error
| 1064 42000 …near 'declare(1 IN 2)' at line 1
> exec ^
: string
| SELECT declare(DISTINCT 1)
< error
| 1064 42000 …near 'declare(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT declare() OVER ()
< error
| 1064 42000 …near 'declare() OVER ()' at line 1
> exec ^
: string
| SELECT declare(1) OVER ()
< error
| 1064 42000 …near 'declare(1) OVER ()' at line 1
> exec ^
: string
| SELECT declare(1,2) OVER ()
< error
| 1064 42000 …near 'declare(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT declare(1,2,3) OVER ()
< error
| 1064 42000 …near 'declare(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT declare(1) FROM t
< error
| 1064 42000 …near 'declare(1) FROM t' at line 1
> exec ^
: string
| SELECT declare(SELECT 1)
< error
| 1064 42000 …near 'declare(SELECT 1)' at line 1
> exec ^
: string
| SELECT declare(YEAR FROM 1)
< error
| 1064 42000 …near 'declare(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT declare(DAY, 1, 2)
< error
| 1064 42000 …near 'declare(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT declare(a)
< error
| 1064 42000 …near 'declare(a)' at line 1
> exec ^
: string
| SELECT declare(a, b)
< error
| 1064 42000 …near 'declare(a, b)' at line 1
> exec ^
: string
| SELECT declare(a) OVER ()
< error
| 1064 42000 …near 'declare(a) OVER ()' at line 1
> exec ^
: string
| SELECT default()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT default(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT default(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT default(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT default(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT default(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT default('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT default(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT default(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT default(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT default(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT default(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT default(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT default() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT default(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT default(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT default(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT default(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT default(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT default(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT default(DAY, 1, 2)
< error
| 1064 42000 …near ' 1, 2)' at line 1
> exec ^
: string
| SELECT default(a)
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec ^
: string
| SELECT default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT default_auth()
< error
| 1305 42000 FUNCTION test.default_auth does not exist
> exec ^
: string
| SELECT default_auth()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT definer()
< error
| 1630 42000 FUNCTION test.definer does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT definer()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT definition()
< error
| 1305 42000 FUNCTION test.definition does not exist
> exec ^
: string
| SELECT definition()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT delay_key_write()
< error
| 1630 42000 FUNCTION test.delay_key_write does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT delay_key_write()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT delayed()
< error
| 1064 42000 …near 'delayed()' at line 1
> exec ^
: string
| SELECT delayed(1)
< error
| 1064 42000 …near 'delayed(1)' at line 1
> exec ^
: string
| SELECT delayed(1,2)
< error
| 1064 42000 …near 'delayed(1,2)' at line 1
> exec ^
: string
| SELECT delayed(1,2,3)
< error
| 1064 42000 …near 'delayed(1,2,3)' at line 1
> exec ^
: string
| SELECT delayed(1,2,3,4)
< error
| 1064 42000 …near 'delayed(1,2,3,4)' at line 1
> exec ^
: string
| SELECT delayed(*)
< error
| 1064 42000 …near 'delayed(*)' at line 1
> exec ^
: string
| SELECT delayed('a')
< error
| 1064 42000 …near 'delayed('a')' at line 1
> exec ^
: string
| SELECT delayed(a USING utf8mb4)
< error
| 1064 42000 …near 'delayed(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT delayed(1 AS SIGNED)
< error
| 1064 42000 …near 'delayed(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT delayed(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'delayed(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT delayed(1 FROM 2)
< error
| 1064 42000 …near 'delayed(1 FROM 2)' at line 1
> exec ^
: string
| SELECT delayed(1 IN 2)
< error
| 1064 42000 …near 'delayed(1 IN 2)' at line 1
> exec ^
: string
| SELECT delayed(DISTINCT 1)
< error
| 1064 42000 …near 'delayed(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT delayed() OVER ()
< error
| 1064 42000 …near 'delayed() OVER ()' at line 1
> exec ^
: string
| SELECT delayed(1) OVER ()
< error
| 1064 42000 …near 'delayed(1) OVER ()' at line 1
> exec ^
: string
| SELECT delayed(1,2) OVER ()
< error
| 1064 42000 …near 'delayed(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT delayed(1,2,3) OVER ()
< error
| 1064 42000 …near 'delayed(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT delayed(1) FROM t
< error
| 1064 42000 …near 'delayed(1) FROM t' at line 1
> exec ^
: string
| SELECT delayed(SELECT 1)
< error
| 1064 42000 …near 'delayed(SELECT 1)' at line 1
> exec ^
: string
| SELECT delayed(YEAR FROM 1)
< error
| 1064 42000 …near 'delayed(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT delayed(DAY, 1, 2)
< error
| 1064 42000 …near 'delayed(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT delayed(a)
< error
| 1064 42000 …near 'delayed(a)' at line 1
> exec ^
: string
| SELECT delayed(a, b)
< error
| 1064 42000 …near 'delayed(a, b)' at line 1
> exec ^
: string
| SELECT delayed(a) OVER ()
< error
| 1064 42000 …near 'delayed(a) OVER ()' at line 1
> exec ^
: string
| SELECT delete()
< error
| 1064 42000 …near 'delete()' at line 1
> exec ^
: string
| SELECT delete(1)
< error
| 1064 42000 …near 'delete(1)' at line 1
> exec ^
: string
| SELECT delete(1,2)
< error
| 1064 42000 …near 'delete(1,2)' at line 1
> exec ^
: string
| SELECT delete(1,2,3)
< error
| 1064 42000 …near 'delete(1,2,3)' at line 1
> exec ^
: string
| SELECT delete(1,2,3,4)
< error
| 1064 42000 …near 'delete(1,2,3,4)' at line 1
> exec ^
: string
| SELECT delete(*)
< error
| 1064 42000 …near 'delete(*)' at line 1
> exec ^
: string
| SELECT delete('a')
< error
| 1064 42000 …near 'delete('a')' at line 1
> exec ^
: string
| SELECT delete(a USING utf8mb4)
< error
| 1064 42000 …near 'delete(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT delete(1 AS SIGNED)
< error
| 1064 42000 …near 'delete(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT delete(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'delete(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT delete(1 FROM 2)
< error
| 1064 42000 …near 'delete(1 FROM 2)' at line 1
> exec ^
: string
| SELECT delete(1 IN 2)
< error
| 1064 42000 …near 'delete(1 IN 2)' at line 1
> exec ^
: string
| SELECT delete(DISTINCT 1)
< error
| 1064 42000 …near 'delete(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT delete() OVER ()
< error
| 1064 42000 …near 'delete() OVER ()' at line 1
> exec ^
: string
| SELECT delete(1) OVER ()
< error
| 1064 42000 …near 'delete(1) OVER ()' at line 1
> exec ^
: string
| SELECT delete(1,2) OVER ()
< error
| 1064 42000 …near 'delete(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT delete(1,2,3) OVER ()
< error
| 1064 42000 …near 'delete(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT delete(1) FROM t
< error
| 1064 42000 …near 'delete(1) FROM t' at line 1
> exec ^
: string
| SELECT delete(SELECT 1)
< error
| 1064 42000 …near 'delete(SELECT 1)' at line 1
> exec ^
: string
| SELECT delete(YEAR FROM 1)
< error
| 1064 42000 …near 'delete(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT delete(DAY, 1, 2)
< error
| 1064 42000 …near 'delete(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT delete(a)
< error
| 1064 42000 …near 'delete(a)' at line 1
> exec ^
: string
| SELECT delete(a, b)
< error
| 1064 42000 …near 'delete(a, b)' at line 1
> exec ^
: string
| SELECT delete(a) OVER ()
< error
| 1064 42000 …near 'delete(a) OVER ()' at line 1
> exec ^
: string
| SELECT dense_rank()
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT dense_rank(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT dense_rank(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT dense_rank(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT dense_rank(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT dense_rank(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT dense_rank('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT dense_rank(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT dense_rank(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT dense_rank(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT dense_rank(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT dense_rank(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT dense_rank(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT dense_rank() OVER ()
< error
| 4018 HY000 No order list in window specification for 'dense_rank'
> exec ^
: string
| SELECT dense_rank() OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT desc()
< error
| 1064 42000 …near 'desc()' at line 1
> exec ^
: string
| SELECT desc(1)
< error
| 1064 42000 …near 'desc(1)' at line 1
> exec ^
: string
| SELECT desc(1,2)
< error
| 1064 42000 …near 'desc(1,2)' at line 1
> exec ^
: string
| SELECT desc(1,2,3)
< error
| 1064 42000 …near 'desc(1,2,3)' at line 1
> exec ^
: string
| SELECT desc(1,2,3,4)
< error
| 1064 42000 …near 'desc(1,2,3,4)' at line 1
> exec ^
: string
| SELECT desc(*)
< error
| 1064 42000 …near 'desc(*)' at line 1
> exec ^
: string
| SELECT desc('a')
< error
| 1064 42000 …near 'desc('a')' at line 1
> exec ^
: string
| SELECT desc(a USING utf8mb4)
< error
| 1064 42000 …near 'desc(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT desc(1 AS SIGNED)
< error
| 1064 42000 …near 'desc(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT desc(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'desc(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT desc(1 FROM 2)
< error
| 1064 42000 …near 'desc(1 FROM 2)' at line 1
> exec ^
: string
| SELECT desc(1 IN 2)
< error
| 1064 42000 …near 'desc(1 IN 2)' at line 1
> exec ^
: string
| SELECT desc(DISTINCT 1)
< error
| 1064 42000 …near 'desc(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT desc() OVER ()
< error
| 1064 42000 …near 'desc() OVER ()' at line 1
> exec ^
: string
| SELECT desc(1) OVER ()
< error
| 1064 42000 …near 'desc(1) OVER ()' at line 1
> exec ^
: string
| SELECT desc(1,2) OVER ()
< error
| 1064 42000 …near 'desc(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT desc(1,2,3) OVER ()
< error
| 1064 42000 …near 'desc(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT desc(1) FROM t
< error
| 1064 42000 …near 'desc(1) FROM t' at line 1
> exec ^
: string
| SELECT desc(SELECT 1)
< error
| 1064 42000 …near 'desc(SELECT 1)' at line 1
> exec ^
: string
| SELECT desc(YEAR FROM 1)
< error
| 1064 42000 …near 'desc(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT desc(DAY, 1, 2)
< error
| 1064 42000 …near 'desc(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT desc(a)
< error
| 1064 42000 …near 'desc(a)' at line 1
> exec ^
: string
| SELECT desc(a, b)
< error
| 1064 42000 …near 'desc(a, b)' at line 1
> exec ^
: string
| SELECT desc(a) OVER ()
< error
| 1064 42000 …near 'desc(a) OVER ()' at line 1
> exec ^
: string
| SELECT describe()
< error
| 1064 42000 …near 'describe()' at line 1
> exec ^
: string
| SELECT describe(1)
< error
| 1064 42000 …near 'describe(1)' at line 1
> exec ^
: string
| SELECT describe(1,2)
< error
| 1064 42000 …near 'describe(1,2)' at line 1
> exec ^
: string
| SELECT describe(1,2,3)
< error
| 1064 42000 …near 'describe(1,2,3)' at line 1
> exec ^
: string
| SELECT describe(1,2,3,4)
< error
| 1064 42000 …near 'describe(1,2,3,4)' at line 1
> exec ^
: string
| SELECT describe(*)
< error
| 1064 42000 …near 'describe(*)' at line 1
> exec ^
: string
| SELECT describe('a')
< error
| 1064 42000 …near 'describe('a')' at line 1
> exec ^
: string
| SELECT describe(a USING utf8mb4)
< error
| 1064 42000 …near 'describe(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT describe(1 AS SIGNED)
< error
| 1064 42000 …near 'describe(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT describe(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'describe(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT describe(1 FROM 2)
< error
| 1064 42000 …near 'describe(1 FROM 2)' at line 1
> exec ^
: string
| SELECT describe(1 IN 2)
< error
| 1064 42000 …near 'describe(1 IN 2)' at line 1
> exec ^
: string
| SELECT describe(DISTINCT 1)
< error
| 1064 42000 …near 'describe(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT describe() OVER ()
< error
| 1064 42000 …near 'describe() OVER ()' at line 1
> exec ^
: string
| SELECT describe(1) OVER ()
< error
| 1064 42000 …near 'describe(1) OVER ()' at line 1
> exec ^
: string
| SELECT describe(1,2) OVER ()
< error
| 1064 42000 …near 'describe(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT describe(1,2,3) OVER ()
< error
| 1064 42000 …near 'describe(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT describe(1) FROM t
< error
| 1064 42000 …near 'describe(1) FROM t' at line 1
> exec ^
: string
| SELECT describe(SELECT 1)
< error
| 1064 42000 …near 'describe(SELECT 1)' at line 1
> exec ^
: string
| SELECT describe(YEAR FROM 1)
< error
| 1064 42000 …near 'describe(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT describe(DAY, 1, 2)
< error
| 1064 42000 …near 'describe(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT describe(a)
< error
| 1064 42000 …near 'describe(a)' at line 1
> exec ^
: string
| SELECT describe(a, b)
< error
| 1064 42000 …near 'describe(a, b)' at line 1
> exec ^
: string
| SELECT describe(a) OVER ()
< error
| 1064 42000 …near 'describe(a) OVER ()' at line 1
> exec ^
: string
| SELECT description()
< error
| 1305 42000 FUNCTION test.description does not exist
> exec ^
: string
| SELECT description()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT deterministic()
< error
| 1064 42000 …near 'deterministic()' at line 1
> exec ^
: string
| SELECT deterministic(1)
< error
| 1064 42000 …near 'deterministic(1)' at line 1
> exec ^
: string
| SELECT deterministic(1,2)
< error
| 1064 42000 …near 'deterministic(1,2)' at line 1
> exec ^
: string
| SELECT deterministic(1,2,3)
< error
| 1064 42000 …near 'deterministic(1,2,3)' at line 1
> exec ^
: string
| SELECT deterministic(1,2,3,4)
< error
| 1064 42000 …near 'deterministic(1,2,3,4)' at line 1
> exec ^
: string
| SELECT deterministic(*)
< error
| 1064 42000 …near 'deterministic(*)' at line 1
> exec ^
: string
| SELECT deterministic('a')
< error
| 1064 42000 …near 'deterministic('a')' at line 1
> exec ^
: string
| SELECT deterministic(a USING utf8mb4)
< error
| 1064 42000 …near 'deterministic(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT deterministic(1 AS SIGNED)
< error
| 1064 42000 …near 'deterministic(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT deterministic(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'deterministic(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT deterministic(1 FROM 2)
< error
| 1064 42000 …near 'deterministic(1 FROM 2)' at line 1
> exec ^
: string
| SELECT deterministic(1 IN 2)
< error
| 1064 42000 …near 'deterministic(1 IN 2)' at line 1
> exec ^
: string
| SELECT deterministic(DISTINCT 1)
< error
| 1064 42000 …near 'deterministic(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT deterministic() OVER ()
< error
| 1064 42000 …near 'deterministic() OVER ()' at line 1
> exec ^
: string
| SELECT deterministic(1) OVER ()
< error
| 1064 42000 …near 'deterministic(1) OVER ()' at line 1
> exec ^
: string
| SELECT deterministic(1,2) OVER ()
< error
| 1064 42000 …near 'deterministic(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT deterministic(1,2,3) OVER ()
< error
| 1064 42000 …near 'deterministic(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT deterministic(1) FROM t
< error
| 1064 42000 …near 'deterministic(1) FROM t' at line 1
> exec ^
: string
| SELECT deterministic(SELECT 1)
< error
| 1064 42000 …near 'deterministic(SELECT 1)' at line 1
> exec ^
: string
| SELECT deterministic(YEAR FROM 1)
< error
| 1064 42000 …near 'deterministic(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT deterministic(DAY, 1, 2)
< error
| 1064 42000 …near 'deterministic(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT deterministic(a)
< error
| 1064 42000 …near 'deterministic(a)' at line 1
> exec ^
: string
| SELECT deterministic(a, b)
< error
| 1064 42000 …near 'deterministic(a, b)' at line 1
> exec ^
: string
| SELECT deterministic(a) OVER ()
< error
| 1064 42000 …near 'deterministic(a) OVER ()' at line 1
> exec ^
: string
| SELECT diagnostics()
< error
| 1630 42000 FUNCTION test.diagnostics does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT diagnostics()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT directory()
< error
| 1630 42000 FUNCTION test.directory does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT directory()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT disable()
< error
| 1630 42000 FUNCTION test.disable does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT disable()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT discard()
< error
| 1630 42000 FUNCTION test.discard does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT discard()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT disk()
< error
| 1630 42000 FUNCTION test.disk does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT disk()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT distinct()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT distinct(1)
< ok 0 0
> exec ^
: string
| SELECT distinctrow()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT distinctrow(1)
< ok 0 0
> exec ^
: string
| SELECT div()
< error
| 1064 42000 …near 'div()' at line 1
> exec ^
: string
| SELECT div(1)
< error
| 1064 42000 …near 'div(1)' at line 1
> exec ^
: string
| SELECT div(1,2)
< error
| 1064 42000 …near 'div(1,2)' at line 1
> exec ^
: string
| SELECT div(1,2,3)
< error
| 1064 42000 …near 'div(1,2,3)' at line 1
> exec ^
: string
| SELECT div(1,2,3,4)
< error
| 1064 42000 …near 'div(1,2,3,4)' at line 1
> exec ^
: string
| SELECT div(*)
< error
| 1064 42000 …near 'div(*)' at line 1
> exec ^
: string
| SELECT div('a')
< error
| 1064 42000 …near 'div('a')' at line 1
> exec ^
: string
| SELECT div(a USING utf8mb4)
< error
| 1064 42000 …near 'div(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT div(1 AS SIGNED)
< error
| 1064 42000 …near 'div(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT div(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'div(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT div(1 FROM 2)
< error
| 1064 42000 …near 'div(1 FROM 2)' at line 1
> exec ^
: string
| SELECT div(1 IN 2)
< error
| 1064 42000 …near 'div(1 IN 2)' at line 1
> exec ^
: string
| SELECT div(DISTINCT 1)
< error
| 1064 42000 …near 'div(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT div() OVER ()
< error
| 1064 42000 …near 'div() OVER ()' at line 1
> exec ^
: string
| SELECT div(1) OVER ()
< error
| 1064 42000 …near 'div(1) OVER ()' at line 1
> exec ^
: string
| SELECT div(1,2) OVER ()
< error
| 1064 42000 …near 'div(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT div(1,2,3) OVER ()
< error
| 1064 42000 …near 'div(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT div(1) FROM t
< error
| 1064 42000 …near 'div(1) FROM t' at line 1
> exec ^
: string
| SELECT div(SELECT 1)
< error
| 1064 42000 …near 'div(SELECT 1)' at line 1
> exec ^
: string
| SELECT div(YEAR FROM 1)
< error
| 1064 42000 …near 'div(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT div(DAY, 1, 2)
< error
| 1064 42000 …near 'div(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT div(a)
< error
| 1064 42000 …near 'div(a)' at line 1
> exec ^
: string
| SELECT div(a, b)
< error
| 1064 42000 …near 'div(a, b)' at line 1
> exec ^
: string
| SELECT div(a) OVER ()
< error
| 1064 42000 …near 'div(a) OVER ()' at line 1
> exec ^
: string
| SELECT do()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT do(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT do(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT do(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT do(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT do(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT do('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT do(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT do(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT do(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT do(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT do(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT do(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT do() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT do(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT do(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT do(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT do(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT do(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT do(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT do(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT do(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT do(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT do(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT double()
< error
| 1064 42000 …near 'double()' at line 1
> exec ^
: string
| SELECT double(1)
< error
| 1064 42000 …near 'double(1)' at line 1
> exec ^
: string
| SELECT double(1,2)
< error
| 1064 42000 …near 'double(1,2)' at line 1
> exec ^
: string
| SELECT double(1,2,3)
< error
| 1064 42000 …near 'double(1,2,3)' at line 1
> exec ^
: string
| SELECT double(1,2,3,4)
< error
| 1064 42000 …near 'double(1,2,3,4)' at line 1
> exec ^
: string
| SELECT double(*)
< error
| 1064 42000 …near 'double(*)' at line 1
> exec ^
: string
| SELECT double('a')
< error
| 1064 42000 …near 'double('a')' at line 1
> exec ^
: string
| SELECT double(a USING utf8mb4)
< error
| 1064 42000 …near 'double(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT double(1 AS SIGNED)
< error
| 1064 42000 …near 'double(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT double(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'double(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT double(1 FROM 2)
< error
| 1064 42000 …near 'double(1 FROM 2)' at line 1
> exec ^
: string
| SELECT double(1 IN 2)
< error
| 1064 42000 …near 'double(1 IN 2)' at line 1
> exec ^
: string
| SELECT double(DISTINCT 1)
< error
| 1064 42000 …near 'double(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT double() OVER ()
< error
| 1064 42000 …near 'double() OVER ()' at line 1
> exec ^
: string
| SELECT double(1) OVER ()
< error
| 1064 42000 …near 'double(1) OVER ()' at line 1
> exec ^
: string
| SELECT double(1,2) OVER ()
< error
| 1064 42000 …near 'double(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT double(1,2,3) OVER ()
< error
| 1064 42000 …near 'double(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT double(1) FROM t
< error
| 1064 42000 …near 'double(1) FROM t' at line 1
> exec ^
: string
| SELECT double(SELECT 1)
< error
| 1064 42000 …near 'double(SELECT 1)' at line 1
> exec ^
: string
| SELECT double(YEAR FROM 1)
< error
| 1064 42000 …near 'double(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT double(DAY, 1, 2)
< error
| 1064 42000 …near 'double(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT double(a)
< error
| 1064 42000 …near 'double(a)' at line 1
> exec ^
: string
| SELECT double(a, b)
< error
| 1064 42000 …near 'double(a, b)' at line 1
> exec ^
: string
| SELECT double(a) OVER ()
< error
| 1064 42000 …near 'double(a) OVER ()' at line 1
> exec ^
: string
| SELECT drop()
< error
| 1064 42000 …near 'drop()' at line 1
> exec ^
: string
| SELECT drop(1)
< error
| 1064 42000 …near 'drop(1)' at line 1
> exec ^
: string
| SELECT drop(1,2)
< error
| 1064 42000 …near 'drop(1,2)' at line 1
> exec ^
: string
| SELECT drop(1,2,3)
< error
| 1064 42000 …near 'drop(1,2,3)' at line 1
> exec ^
: string
| SELECT drop(1,2,3,4)
< error
| 1064 42000 …near 'drop(1,2,3,4)' at line 1
> exec ^
: string
| SELECT drop(*)
< error
| 1064 42000 …near 'drop(*)' at line 1
> exec ^
: string
| SELECT drop('a')
< error
| 1064 42000 …near 'drop('a')' at line 1
> exec ^
: string
| SELECT drop(a USING utf8mb4)
< error
| 1064 42000 …near 'drop(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT drop(1 AS SIGNED)
< error
| 1064 42000 …near 'drop(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT drop(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'drop(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT drop(1 FROM 2)
< error
| 1064 42000 …near 'drop(1 FROM 2)' at line 1
> exec ^
: string
| SELECT drop(1 IN 2)
< error
| 1064 42000 …near 'drop(1 IN 2)' at line 1
> exec ^
: string
| SELECT drop(DISTINCT 1)
< error
| 1064 42000 …near 'drop(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT drop() OVER ()
< error
| 1064 42000 …near 'drop() OVER ()' at line 1
> exec ^
: string
| SELECT drop(1) OVER ()
< error
| 1064 42000 …near 'drop(1) OVER ()' at line 1
> exec ^
: string
| SELECT drop(1,2) OVER ()
< error
| 1064 42000 …near 'drop(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT drop(1,2,3) OVER ()
< error
| 1064 42000 …near 'drop(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT drop(1) FROM t
< error
| 1064 42000 …near 'drop(1) FROM t' at line 1
> exec ^
: string
| SELECT drop(SELECT 1)
< error
| 1064 42000 …near 'drop(SELECT 1)' at line 1
> exec ^
: string
| SELECT drop(YEAR FROM 1)
< error
| 1064 42000 …near 'drop(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT drop(DAY, 1, 2)
< error
| 1064 42000 …near 'drop(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT drop(a)
< error
| 1064 42000 …near 'drop(a)' at line 1
> exec ^
: string
| SELECT drop(a, b)
< error
| 1064 42000 …near 'drop(a, b)' at line 1
> exec ^
: string
| SELECT drop(a) OVER ()
< error
| 1064 42000 …near 'drop(a) OVER ()' at line 1
> exec ^
: string
| SELECT dual()
< error
| 1064 42000 …near 'dual()' at line 1
> exec ^
: string
| SELECT dual(1)
< error
| 1064 42000 …near 'dual(1)' at line 1
> exec ^
: string
| SELECT dual(1,2)
< error
| 1064 42000 …near 'dual(1,2)' at line 1
> exec ^
: string
| SELECT dual(1,2,3)
< error
| 1064 42000 …near 'dual(1,2,3)' at line 1
> exec ^
: string
| SELECT dual(1,2,3,4)
< error
| 1064 42000 …near 'dual(1,2,3,4)' at line 1
> exec ^
: string
| SELECT dual(*)
< error
| 1064 42000 …near 'dual(*)' at line 1
> exec ^
: string
| SELECT dual('a')
< error
| 1064 42000 …near 'dual('a')' at line 1
> exec ^
: string
| SELECT dual(a USING utf8mb4)
< error
| 1064 42000 …near 'dual(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT dual(1 AS SIGNED)
< error
| 1064 42000 …near 'dual(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT dual(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'dual(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT dual(1 FROM 2)
< error
| 1064 42000 …near 'dual(1 FROM 2)' at line 1
> exec ^
: string
| SELECT dual(1 IN 2)
< error
| 1064 42000 …near 'dual(1 IN 2)' at line 1
> exec ^
: string
| SELECT dual(DISTINCT 1)
< error
| 1064 42000 …near 'dual(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT dual() OVER ()
< error
| 1064 42000 …near 'dual() OVER ()' at line 1
> exec ^
: string
| SELECT dual(1) OVER ()
< error
| 1064 42000 …near 'dual(1) OVER ()' at line 1
> exec ^
: string
| SELECT dual(1,2) OVER ()
< error
| 1064 42000 …near 'dual(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT dual(1,2,3) OVER ()
< error
| 1064 42000 …near 'dual(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT dual(1) FROM t
< error
| 1064 42000 …near 'dual(1) FROM t' at line 1
> exec ^
: string
| SELECT dual(SELECT 1)
< error
| 1064 42000 …near 'dual(SELECT 1)' at line 1
> exec ^
: string
| SELECT dual(YEAR FROM 1)
< error
| 1064 42000 …near 'dual(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT dual(DAY, 1, 2)
< error
| 1064 42000 …near 'dual(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT dual(a)
< error
| 1064 42000 …near 'dual(a)' at line 1
> exec ^
: string
| SELECT dual(a, b)
< error
| 1064 42000 …near 'dual(a, b)' at line 1
> exec ^
: string
| SELECT dual(a) OVER ()
< error
| 1064 42000 …near 'dual(a) OVER ()' at line 1
> exec ^
: string
| SELECT duality()
< error
| 1305 42000 FUNCTION test.duality does not exist
> exec ^
: string
| SELECT duality()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT dumpfile()
< error
| 1630 42000 FUNCTION test.dumpfile does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT dumpfile()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT duplicate()
< error
| 1630 42000 FUNCTION test.duplicate does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT duplicate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT dynamic()
< error
| 1630 42000 FUNCTION test.dynamic does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT dynamic()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT each()
< error
| 1064 42000 …near 'each()' at line 1
> exec ^
: string
| SELECT each(1)
< error
| 1064 42000 …near 'each(1)' at line 1
> exec ^
: string
| SELECT each(1,2)
< error
| 1064 42000 …near 'each(1,2)' at line 1
> exec ^
: string
| SELECT each(1,2,3)
< error
| 1064 42000 …near 'each(1,2,3)' at line 1
> exec ^
: string
| SELECT each(1,2,3,4)
< error
| 1064 42000 …near 'each(1,2,3,4)' at line 1
> exec ^
: string
| SELECT each(*)
< error
| 1064 42000 …near 'each(*)' at line 1
> exec ^
: string
| SELECT each('a')
< error
| 1064 42000 …near 'each('a')' at line 1
> exec ^
: string
| SELECT each(a USING utf8mb4)
< error
| 1064 42000 …near 'each(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT each(1 AS SIGNED)
< error
| 1064 42000 …near 'each(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT each(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'each(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT each(1 FROM 2)
< error
| 1064 42000 …near 'each(1 FROM 2)' at line 1
> exec ^
: string
| SELECT each(1 IN 2)
< error
| 1064 42000 …near 'each(1 IN 2)' at line 1
> exec ^
: string
| SELECT each(DISTINCT 1)
< error
| 1064 42000 …near 'each(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT each() OVER ()
< error
| 1064 42000 …near 'each() OVER ()' at line 1
> exec ^
: string
| SELECT each(1) OVER ()
< error
| 1064 42000 …near 'each(1) OVER ()' at line 1
> exec ^
: string
| SELECT each(1,2) OVER ()
< error
| 1064 42000 …near 'each(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT each(1,2,3) OVER ()
< error
| 1064 42000 …near 'each(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT each(1) FROM t
< error
| 1064 42000 …near 'each(1) FROM t' at line 1
> exec ^
: string
| SELECT each(SELECT 1)
< error
| 1064 42000 …near 'each(SELECT 1)' at line 1
> exec ^
: string
| SELECT each(YEAR FROM 1)
< error
| 1064 42000 …near 'each(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT each(DAY, 1, 2)
< error
| 1064 42000 …near 'each(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT each(a)
< error
| 1064 42000 …near 'each(a)' at line 1
> exec ^
: string
| SELECT each(a, b)
< error
| 1064 42000 …near 'each(a, b)' at line 1
> exec ^
: string
| SELECT each(a) OVER ()
< error
| 1064 42000 …near 'each(a) OVER ()' at line 1
> exec ^
: string
| SELECT else()
< error
| 1064 42000 …near 'else()' at line 1
> exec ^
: string
| SELECT else(1)
< error
| 1064 42000 …near 'else(1)' at line 1
> exec ^
: string
| SELECT else(1,2)
< error
| 1064 42000 …near 'else(1,2)' at line 1
> exec ^
: string
| SELECT else(1,2,3)
< error
| 1064 42000 …near 'else(1,2,3)' at line 1
> exec ^
: string
| SELECT else(1,2,3,4)
< error
| 1064 42000 …near 'else(1,2,3,4)' at line 1
> exec ^
: string
| SELECT else(*)
< error
| 1064 42000 …near 'else(*)' at line 1
> exec ^
: string
| SELECT else('a')
< error
| 1064 42000 …near 'else('a')' at line 1
> exec ^
: string
| SELECT else(a USING utf8mb4)
< error
| 1064 42000 …near 'else(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT else(1 AS SIGNED)
< error
| 1064 42000 …near 'else(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT else(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'else(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT else(1 FROM 2)
< error
| 1064 42000 …near 'else(1 FROM 2)' at line 1
> exec ^
: string
| SELECT else(1 IN 2)
< error
| 1064 42000 …near 'else(1 IN 2)' at line 1
> exec ^
: string
| SELECT else(DISTINCT 1)
< error
| 1064 42000 …near 'else(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT else() OVER ()
< error
| 1064 42000 …near 'else() OVER ()' at line 1
> exec ^
: string
| SELECT else(1) OVER ()
< error
| 1064 42000 …near 'else(1) OVER ()' at line 1
> exec ^
: string
| SELECT else(1,2) OVER ()
< error
| 1064 42000 …near 'else(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT else(1,2,3) OVER ()
< error
| 1064 42000 …near 'else(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT else(1) FROM t
< error
| 1064 42000 …near 'else(1) FROM t' at line 1
> exec ^
: string
| SELECT else(SELECT 1)
< error
| 1064 42000 …near 'else(SELECT 1)' at line 1
> exec ^
: string
| SELECT else(YEAR FROM 1)
< error
| 1064 42000 …near 'else(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT else(DAY, 1, 2)
< error
| 1064 42000 …near 'else(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT else(a)
< error
| 1064 42000 …near 'else(a)' at line 1
> exec ^
: string
| SELECT else(a, b)
< error
| 1064 42000 …near 'else(a, b)' at line 1
> exec ^
: string
| SELECT else(a) OVER ()
< error
| 1064 42000 …near 'else(a) OVER ()' at line 1
> exec ^
: string
| SELECT elseif()
< error
| 1064 42000 …near 'elseif()' at line 1
> exec ^
: string
| SELECT elseif(1)
< error
| 1064 42000 …near 'elseif(1)' at line 1
> exec ^
: string
| SELECT elseif(1,2)
< error
| 1064 42000 …near 'elseif(1,2)' at line 1
> exec ^
: string
| SELECT elseif(1,2,3)
< error
| 1064 42000 …near 'elseif(1,2,3)' at line 1
> exec ^
: string
| SELECT elseif(1,2,3,4)
< error
| 1064 42000 …near 'elseif(1,2,3,4)' at line 1
> exec ^
: string
| SELECT elseif(*)
< error
| 1064 42000 …near 'elseif(*)' at line 1
> exec ^
: string
| SELECT elseif('a')
< error
| 1064 42000 …near 'elseif('a')' at line 1
> exec ^
: string
| SELECT elseif(a USING utf8mb4)
< error
| 1064 42000 …near 'elseif(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT elseif(1 AS SIGNED)
< error
| 1064 42000 …near 'elseif(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT elseif(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'elseif(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT elseif(1 FROM 2)
< error
| 1064 42000 …near 'elseif(1 FROM 2)' at line 1
> exec ^
: string
| SELECT elseif(1 IN 2)
< error
| 1064 42000 …near 'elseif(1 IN 2)' at line 1
> exec ^
: string
| SELECT elseif(DISTINCT 1)
< error
| 1064 42000 …near 'elseif(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT elseif() OVER ()
< error
| 1064 42000 …near 'elseif() OVER ()' at line 1
> exec ^
: string
| SELECT elseif(1) OVER ()
< error
| 1064 42000 …near 'elseif(1) OVER ()' at line 1
> exec ^
: string
| SELECT elseif(1,2) OVER ()
< error
| 1064 42000 …near 'elseif(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT elseif(1,2,3) OVER ()
< error
| 1064 42000 …near 'elseif(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT elseif(1) FROM t
< error
| 1064 42000 …near 'elseif(1) FROM t' at line 1
> exec ^
: string
| SELECT elseif(SELECT 1)
< error
| 1064 42000 …near 'elseif(SELECT 1)' at line 1
> exec ^
: string
| SELECT elseif(YEAR FROM 1)
< error
| 1064 42000 …near 'elseif(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT elseif(DAY, 1, 2)
< error
| 1064 42000 …near 'elseif(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT elseif(a)
< error
| 1064 42000 …near 'elseif(a)' at line 1
> exec ^
: string
| SELECT elseif(a, b)
< error
| 1064 42000 …near 'elseif(a, b)' at line 1
> exec ^
: string
| SELECT elseif(a) OVER ()
< error
| 1064 42000 …near 'elseif(a) OVER ()' at line 1
> exec ^
: string
| SELECT empty()
< error
| 1630 42000 FUNCTION test.empty does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT empty()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT enable()
< error
| 1630 42000 FUNCTION test.enable does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT enable()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT enclosed()
< error
| 1064 42000 …near 'enclosed()' at line 1
> exec ^
: string
| SELECT enclosed(1)
< error
| 1064 42000 …near 'enclosed(1)' at line 1
> exec ^
: string
| SELECT enclosed(1,2)
< error
| 1064 42000 …near 'enclosed(1,2)' at line 1
> exec ^
: string
| SELECT enclosed(1,2,3)
< error
| 1064 42000 …near 'enclosed(1,2,3)' at line 1
> exec ^
: string
| SELECT enclosed(1,2,3,4)
< error
| 1064 42000 …near 'enclosed(1,2,3,4)' at line 1
> exec ^
: string
| SELECT enclosed(*)
< error
| 1064 42000 …near 'enclosed(*)' at line 1
> exec ^
: string
| SELECT enclosed('a')
< error
| 1064 42000 …near 'enclosed('a')' at line 1
> exec ^
: string
| SELECT enclosed(a USING utf8mb4)
< error
| 1064 42000 …near 'enclosed(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT enclosed(1 AS SIGNED)
< error
| 1064 42000 …near 'enclosed(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT enclosed(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'enclosed(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT enclosed(1 FROM 2)
< error
| 1064 42000 …near 'enclosed(1 FROM 2)' at line 1
> exec ^
: string
| SELECT enclosed(1 IN 2)
< error
| 1064 42000 …near 'enclosed(1 IN 2)' at line 1
> exec ^
: string
| SELECT enclosed(DISTINCT 1)
< error
| 1064 42000 …near 'enclosed(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT enclosed() OVER ()
< error
| 1064 42000 …near 'enclosed() OVER ()' at line 1
> exec ^
: string
| SELECT enclosed(1) OVER ()
< error
| 1064 42000 …near 'enclosed(1) OVER ()' at line 1
> exec ^
: string
| SELECT enclosed(1,2) OVER ()
< error
| 1064 42000 …near 'enclosed(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT enclosed(1,2,3) OVER ()
< error
| 1064 42000 …near 'enclosed(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT enclosed(1) FROM t
< error
| 1064 42000 …near 'enclosed(1) FROM t' at line 1
> exec ^
: string
| SELECT enclosed(SELECT 1)
< error
| 1064 42000 …near 'enclosed(SELECT 1)' at line 1
> exec ^
: string
| SELECT enclosed(YEAR FROM 1)
< error
| 1064 42000 …near 'enclosed(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT enclosed(DAY, 1, 2)
< error
| 1064 42000 …near 'enclosed(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT enclosed(a)
< error
| 1064 42000 …near 'enclosed(a)' at line 1
> exec ^
: string
| SELECT enclosed(a, b)
< error
| 1064 42000 …near 'enclosed(a, b)' at line 1
> exec ^
: string
| SELECT enclosed(a) OVER ()
< error
| 1064 42000 …near 'enclosed(a) OVER ()' at line 1
> exec ^
: string
| SELECT encryption()
< error
| 1305 42000 FUNCTION test.encryption does not exist
> exec ^
: string
| SELECT encryption()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT end()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT end(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT end(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT end(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT end(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT end(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT end('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT end(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT end(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT end(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT end(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT end(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT end(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT end() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT end(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT end(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT end(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT end(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT end(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT end(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT end(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT end(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT end(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT end(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT ends()
< error
| 1630 42000 FUNCTION test.ends does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT ends()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT enforced()
< error
| 1305 42000 FUNCTION test.enforced does not exist
> exec ^
: string
| SELECT enforced()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT engine()
< error
| 1630 42000 FUNCTION test.engine does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT engine()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT engine_attribute()
< error
| 1305 42000 FUNCTION test.engine_attribute does not exist
> exec ^
: string
| SELECT engine_attribute()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT engines()
< error
| 1630 42000 FUNCTION test.engines does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT engines()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT enum()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT enum(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT enum(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT enum(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT enum(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT enum(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT enum('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT enum(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT enum(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT enum(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT enum(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT enum(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT enum(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT enum() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT enum(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT enum(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT enum(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT enum(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT enum(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT enum(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT enum(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT enum(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT enum(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT enum(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT error()
< error
| 1630 42000 FUNCTION test.error does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT error()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT errors()
< error
| 1630 42000 FUNCTION test.errors does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT errors()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT escape()
< error
| 1630 42000 FUNCTION test.escape does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT escape()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT escaped()
< error
| 1064 42000 …near 'escaped()' at line 1
> exec ^
: string
| SELECT escaped(1)
< error
| 1064 42000 …near 'escaped(1)' at line 1
> exec ^
: string
| SELECT escaped(1,2)
< error
| 1064 42000 …near 'escaped(1,2)' at line 1
> exec ^
: string
| SELECT escaped(1,2,3)
< error
| 1064 42000 …near 'escaped(1,2,3)' at line 1
> exec ^
: string
| SELECT escaped(1,2,3,4)
< error
| 1064 42000 …near 'escaped(1,2,3,4)' at line 1
> exec ^
: string
| SELECT escaped(*)
< error
| 1064 42000 …near 'escaped(*)' at line 1
> exec ^
: string
| SELECT escaped('a')
< error
| 1064 42000 …near 'escaped('a')' at line 1
> exec ^
: string
| SELECT escaped(a USING utf8mb4)
< error
| 1064 42000 …near 'escaped(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT escaped(1 AS SIGNED)
< error
| 1064 42000 …near 'escaped(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT escaped(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'escaped(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT escaped(1 FROM 2)
< error
| 1064 42000 …near 'escaped(1 FROM 2)' at line 1
> exec ^
: string
| SELECT escaped(1 IN 2)
< error
| 1064 42000 …near 'escaped(1 IN 2)' at line 1
> exec ^
: string
| SELECT escaped(DISTINCT 1)
< error
| 1064 42000 …near 'escaped(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT escaped() OVER ()
< error
| 1064 42000 …near 'escaped() OVER ()' at line 1
> exec ^
: string
| SELECT escaped(1) OVER ()
< error
| 1064 42000 …near 'escaped(1) OVER ()' at line 1
> exec ^
: string
| SELECT escaped(1,2) OVER ()
< error
| 1064 42000 …near 'escaped(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT escaped(1,2,3) OVER ()
< error
| 1064 42000 …near 'escaped(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT escaped(1) FROM t
< error
| 1064 42000 …near 'escaped(1) FROM t' at line 1
> exec ^
: string
| SELECT escaped(SELECT 1)
< error
| 1064 42000 …near 'escaped(SELECT 1)' at line 1
> exec ^
: string
| SELECT escaped(YEAR FROM 1)
< error
| 1064 42000 …near 'escaped(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT escaped(DAY, 1, 2)
< error
| 1064 42000 …near 'escaped(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT escaped(a)
< error
| 1064 42000 …near 'escaped(a)' at line 1
> exec ^
: string
| SELECT escaped(a, b)
< error
| 1064 42000 …near 'escaped(a, b)' at line 1
> exec ^
: string
| SELECT escaped(a) OVER ()
< error
| 1064 42000 …near 'escaped(a) OVER ()' at line 1
> exec ^
: string
| SELECT event()
< error
| 1630 42000 FUNCTION test.event does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT event()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT events()
< error
| 1630 42000 FUNCTION test.events does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT events()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT every()
< error
| 1630 42000 FUNCTION test.every does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT every()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT except()
< error
| 1064 42000 …near 'except()' at line 1
> exec ^
: string
| SELECT except(1)
< error
| 1064 42000 …near 'except(1)' at line 1
> exec ^
: string
| SELECT except(1,2)
< error
| 1064 42000 …near 'except(1,2)' at line 1
> exec ^
: string
| SELECT except(1,2,3)
< error
| 1064 42000 …near 'except(1,2,3)' at line 1
> exec ^
: string
| SELECT except(1,2,3,4)
< error
| 1064 42000 …near 'except(1,2,3,4)' at line 1
> exec ^
: string
| SELECT except(*)
< error
| 1064 42000 …near 'except(*)' at line 1
> exec ^
: string
| SELECT except('a')
< error
| 1064 42000 …near 'except('a')' at line 1
> exec ^
: string
| SELECT except(a USING utf8mb4)
< error
| 1064 42000 …near 'except(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT except(1 AS SIGNED)
< error
| 1064 42000 …near 'except(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT except(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'except(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT except(1 FROM 2)
< error
| 1064 42000 …near 'except(1 FROM 2)' at line 1
> exec ^
: string
| SELECT except(1 IN 2)
< error
| 1064 42000 …near 'except(1 IN 2)' at line 1
> exec ^
: string
| SELECT except(DISTINCT 1)
< error
| 1064 42000 …near 'except(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT except() OVER ()
< error
| 1064 42000 …near 'except() OVER ()' at line 1
> exec ^
: string
| SELECT except(1) OVER ()
< error
| 1064 42000 …near 'except(1) OVER ()' at line 1
> exec ^
: string
| SELECT except(1,2) OVER ()
< error
| 1064 42000 …near 'except(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT except(1,2,3) OVER ()
< error
| 1064 42000 …near 'except(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT except(1) FROM t
< error
| 1064 42000 …near 'except(1) FROM t' at line 1
> exec ^
: string
| SELECT except(SELECT 1)
< error
| 1064 42000 …near 'except(SELECT 1)' at line 1
> exec ^
: string
| SELECT except(YEAR FROM 1)
< error
| 1064 42000 …near 'except(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT except(DAY, 1, 2)
< error
| 1064 42000 …near 'except(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT except(a)
< error
| 1064 42000 …near 'except(a)' at line 1
> exec ^
: string
| SELECT except(a, b)
< error
| 1064 42000 …near 'except(a, b)' at line 1
> exec ^
: string
| SELECT except(a) OVER ()
< error
| 1064 42000 …near 'except(a) OVER ()' at line 1
> exec ^
: string
| SELECT exchange()
< error
| 1630 42000 FUNCTION test.exchange does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT exchange()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT exclude()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT exclude(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT exclude(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT exclude(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT exclude(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT exclude(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT exclude('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT exclude(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT exclude(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT exclude(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT exclude(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT exclude(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT exclude(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT exclude() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT exclude(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT exclude(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT exclude(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT exclude(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT exclude(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT exclude(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT exclude(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT exclude(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT exclude(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT exclude(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT execute()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT execute(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT execute(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT execute(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT execute(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT execute(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT execute('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT execute(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT execute(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT execute(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT execute(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT execute(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT execute(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT execute() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT execute(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT execute(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT execute(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT execute(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT execute(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT execute(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT execute(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT execute(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT execute(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT execute(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT exists()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT exists(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT exists(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT exists(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT exists(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT exists(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT exists('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT exists(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT exists(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT exists(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT exists(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT exists(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT exists(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT exists() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT exists(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT exists(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT exists(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT exists(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT exists(SELECT 1)
< ok 0 0
> exec ^
: string
| SELECT exit()
< error
| 1064 42000 …near 'exit()' at line 1
> exec ^
: string
| SELECT exit(1)
< error
| 1064 42000 …near 'exit(1)' at line 1
> exec ^
: string
| SELECT exit(1,2)
< error
| 1064 42000 …near 'exit(1,2)' at line 1
> exec ^
: string
| SELECT exit(1,2,3)
< error
| 1064 42000 …near 'exit(1,2,3)' at line 1
> exec ^
: string
| SELECT exit(1,2,3,4)
< error
| 1064 42000 …near 'exit(1,2,3,4)' at line 1
> exec ^
: string
| SELECT exit(*)
< error
| 1064 42000 …near 'exit(*)' at line 1
> exec ^
: string
| SELECT exit('a')
< error
| 1064 42000 …near 'exit('a')' at line 1
> exec ^
: string
| SELECT exit(a USING utf8mb4)
< error
| 1064 42000 …near 'exit(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT exit(1 AS SIGNED)
< error
| 1064 42000 …near 'exit(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT exit(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'exit(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT exit(1 FROM 2)
< error
| 1064 42000 …near 'exit(1 FROM 2)' at line 1
> exec ^
: string
| SELECT exit(1 IN 2)
< error
| 1064 42000 …near 'exit(1 IN 2)' at line 1
> exec ^
: string
| SELECT exit(DISTINCT 1)
< error
| 1064 42000 …near 'exit(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT exit() OVER ()
< error
| 1064 42000 …near 'exit() OVER ()' at line 1
> exec ^
: string
| SELECT exit(1) OVER ()
< error
| 1064 42000 …near 'exit(1) OVER ()' at line 1
> exec ^
: string
| SELECT exit(1,2) OVER ()
< error
| 1064 42000 …near 'exit(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT exit(1,2,3) OVER ()
< error
| 1064 42000 …near 'exit(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT exit(1) FROM t
< error
| 1064 42000 …near 'exit(1) FROM t' at line 1
> exec ^
: string
| SELECT exit(SELECT 1)
< error
| 1064 42000 …near 'exit(SELECT 1)' at line 1
> exec ^
: string
| SELECT exit(YEAR FROM 1)
< error
| 1064 42000 …near 'exit(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT exit(DAY, 1, 2)
< error
| 1064 42000 …near 'exit(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT exit(a)
< error
| 1064 42000 …near 'exit(a)' at line 1
> exec ^
: string
| SELECT exit(a, b)
< error
| 1064 42000 …near 'exit(a, b)' at line 1
> exec ^
: string
| SELECT exit(a) OVER ()
< error
| 1064 42000 …near 'exit(a) OVER ()' at line 1
> exec ^
: string
| SELECT expansion()
< error
| 1630 42000 FUNCTION test.expansion does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT expansion()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT expire()
< error
| 1630 42000 FUNCTION test.expire does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT expire()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT explain()
< error
| 1064 42000 …near 'explain()' at line 1
> exec ^
: string
| SELECT explain(1)
< error
| 1064 42000 …near 'explain(1)' at line 1
> exec ^
: string
| SELECT explain(1,2)
< error
| 1064 42000 …near 'explain(1,2)' at line 1
> exec ^
: string
| SELECT explain(1,2,3)
< error
| 1064 42000 …near 'explain(1,2,3)' at line 1
> exec ^
: string
| SELECT explain(1,2,3,4)
< error
| 1064 42000 …near 'explain(1,2,3,4)' at line 1
> exec ^
: string
| SELECT explain(*)
< error
| 1064 42000 …near 'explain(*)' at line 1
> exec ^
: string
| SELECT explain('a')
< error
| 1064 42000 …near 'explain('a')' at line 1
> exec ^
: string
| SELECT explain(a USING utf8mb4)
< error
| 1064 42000 …near 'explain(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT explain(1 AS SIGNED)
< error
| 1064 42000 …near 'explain(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT explain(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'explain(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT explain(1 FROM 2)
< error
| 1064 42000 …near 'explain(1 FROM 2)' at line 1
> exec ^
: string
| SELECT explain(1 IN 2)
< error
| 1064 42000 …near 'explain(1 IN 2)' at line 1
> exec ^
: string
| SELECT explain(DISTINCT 1)
< error
| 1064 42000 …near 'explain(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT explain() OVER ()
< error
| 1064 42000 …near 'explain() OVER ()' at line 1
> exec ^
: string
| SELECT explain(1) OVER ()
< error
| 1064 42000 …near 'explain(1) OVER ()' at line 1
> exec ^
: string
| SELECT explain(1,2) OVER ()
< error
| 1064 42000 …near 'explain(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT explain(1,2,3) OVER ()
< error
| 1064 42000 …near 'explain(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT explain(1) FROM t
< error
| 1064 42000 …near 'explain(1) FROM t' at line 1
> exec ^
: string
| SELECT explain(SELECT 1)
< error
| 1064 42000 …near 'explain(SELECT 1)' at line 1
> exec ^
: string
| SELECT explain(YEAR FROM 1)
< error
| 1064 42000 …near 'explain(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT explain(DAY, 1, 2)
< error
| 1064 42000 …near 'explain(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT explain(a)
< error
| 1064 42000 …near 'explain(a)' at line 1
> exec ^
: string
| SELECT explain(a, b)
< error
| 1064 42000 …near 'explain(a, b)' at line 1
> exec ^
: string
| SELECT explain(a) OVER ()
< error
| 1064 42000 …near 'explain(a) OVER ()' at line 1
> exec ^
: string
| SELECT export()
< error
| 1630 42000 FUNCTION test.export does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT export()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT extended()
< error
| 1630 42000 FUNCTION test.extended does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT extended()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT extent_size()
< error
| 1630 42000 FUNCTION test.extent_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT extent_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT external()
< error
| 1305 42000 FUNCTION test.external does not exist
> exec ^
: string
| SELECT external()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT external_format()
< error
| 1305 42000 FUNCTION test.external_format does not exist
> exec ^
: string
| SELECT external_format()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT extract()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT extract(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT extract(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT extract(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT extract(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT extract(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT extract('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT extract(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT extract(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT extract(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT extract(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT extract(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT extract(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT extract() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT extract(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT extract(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT extract(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT extract(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT extract(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT extract(YEAR FROM 1)
< ok 0 0
> exec ^
: string
| SELECT factor()
< error
| 1305 42000 FUNCTION test.factor does not exist
> exec ^
: string
| SELECT factor()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT failed_login_attempts()
< error
| 1305 42000 FUNCTION test.failed_login_attempts does not exist
> exec ^
: string
| SELECT failed_login_attempts()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT false()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT false(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT false(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT false(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT false(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT false(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT false('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT false(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT false(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT false(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT false(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT false(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT false(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT false() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT false(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT false(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT false(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT false(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT false(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT false(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT false(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT false(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT false(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT false(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT fast()
< error
| 1630 42000 FUNCTION test.fast does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT fast()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT faults()
< error
| 1630 42000 FUNCTION test.faults does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT faults()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT fetch()
< error
| 1064 42000 …near 'fetch()' at line 1
> exec ^
: string
| SELECT fetch(1)
< error
| 1064 42000 …near 'fetch(1)' at line 1
> exec ^
: string
| SELECT fetch(1,2)
< error
| 1064 42000 …near 'fetch(1,2)' at line 1
> exec ^
: string
| SELECT fetch(1,2,3)
< error
| 1064 42000 …near 'fetch(1,2,3)' at line 1
> exec ^
: string
| SELECT fetch(1,2,3,4)
< error
| 1064 42000 …near 'fetch(1,2,3,4)' at line 1
> exec ^
: string
| SELECT fetch(*)
< error
| 1064 42000 …near 'fetch(*)' at line 1
> exec ^
: string
| SELECT fetch('a')
< error
| 1064 42000 …near 'fetch('a')' at line 1
> exec ^
: string
| SELECT fetch(a USING utf8mb4)
< error
| 1064 42000 …near 'fetch(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT fetch(1 AS SIGNED)
< error
| 1064 42000 …near 'fetch(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT fetch(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'fetch(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT fetch(1 FROM 2)
< error
| 1064 42000 …near 'fetch(1 FROM 2)' at line 1
> exec ^
: string
| SELECT fetch(1 IN 2)
< error
| 1064 42000 …near 'fetch(1 IN 2)' at line 1
> exec ^
: string
| SELECT fetch(DISTINCT 1)
< error
| 1064 42000 …near 'fetch(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT fetch() OVER ()
< error
| 1064 42000 …near 'fetch() OVER ()' at line 1
> exec ^
: string
| SELECT fetch(1) OVER ()
< error
| 1064 42000 …near 'fetch(1) OVER ()' at line 1
> exec ^
: string
| SELECT fetch(1,2) OVER ()
< error
| 1064 42000 …near 'fetch(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT fetch(1,2,3) OVER ()
< error
| 1064 42000 …near 'fetch(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT fetch(1) FROM t
< error
| 1064 42000 …near 'fetch(1) FROM t' at line 1
> exec ^
: string
| SELECT fetch(SELECT 1)
< error
| 1064 42000 …near 'fetch(SELECT 1)' at line 1
> exec ^
: string
| SELECT fetch(YEAR FROM 1)
< error
| 1064 42000 …near 'fetch(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT fetch(DAY, 1, 2)
< error
| 1064 42000 …near 'fetch(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT fetch(a)
< error
| 1064 42000 …near 'fetch(a)' at line 1
> exec ^
: string
| SELECT fetch(a, b)
< error
| 1064 42000 …near 'fetch(a, b)' at line 1
> exec ^
: string
| SELECT fetch(a) OVER ()
< error
| 1064 42000 …near 'fetch(a) OVER ()' at line 1
> exec ^
: string
| SELECT fields()
< error
| 1630 42000 FUNCTION test.fields does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT fields()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file()
< error
| 1630 42000 FUNCTION test.file does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT file()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file_block_size()
< error
| 1305 42000 FUNCTION test.file_block_size does not exist
> exec ^
: string
| SELECT file_block_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file_format()
< error
| 1305 42000 FUNCTION test.file_format does not exist
> exec ^
: string
| SELECT file_format()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file_name()
< error
| 1305 42000 FUNCTION test.file_name does not exist
> exec ^
: string
| SELECT file_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file_pattern()
< error
| 1305 42000 FUNCTION test.file_pattern does not exist
> exec ^
: string
| SELECT file_pattern()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT file_prefix()
< error
| 1305 42000 FUNCTION test.file_prefix does not exist
> exec ^
: string
| SELECT file_prefix()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT files()
< error
| 1305 42000 FUNCTION test.files does not exist
> exec ^
: string
| SELECT files()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT filter()
< error
| 1305 42000 FUNCTION test.filter does not exist
> exec ^
: string
| SELECT filter()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT finish()
< error
| 1305 42000 FUNCTION test.finish does not exist
> exec ^
: string
| SELECT finish()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT first()
< error
| 1630 42000 FUNCTION test.first does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT first()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT first_value()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT first_value(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT first_value(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT first_value(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT first_value(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT first_value(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT first_value('a')
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT first_value(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT first_value(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT first_value(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ' INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT first_value(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT first_value(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT first_value(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT first_value() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT first_value(1) OVER ()
< ok 0 0
> exec ^
: string
| SELECT fixed()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT fixed(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT fixed(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT fixed(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT fixed(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT fixed(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT fixed('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT fixed(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT fixed(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT fixed(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT fixed(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT fixed(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT fixed(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT fixed() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT fixed(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT fixed(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT fixed(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT fixed(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT fixed(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT fixed(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT fixed(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT fixed(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT fixed(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT fixed(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT float()
< error
| 1064 42000 …near 'float()' at line 1
> exec ^
: string
| SELECT float(1)
< error
| 1064 42000 …near 'float(1)' at line 1
> exec ^
: string
| SELECT float(1,2)
< error
| 1064 42000 …near 'float(1,2)' at line 1
> exec ^
: string
| SELECT float(1,2,3)
< error
| 1064 42000 …near 'float(1,2,3)' at line 1
> exec ^
: string
| SELECT float(1,2,3,4)
< error
| 1064 42000 …near 'float(1,2,3,4)' at line 1
> exec ^
: string
| SELECT float(*)
< error
| 1064 42000 …near 'float(*)' at line 1
> exec ^
: string
| SELECT float('a')
< error
| 1064 42000 …near 'float('a')' at line 1
> exec ^
: string
| SELECT float(a USING utf8mb4)
< error
| 1064 42000 …near 'float(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT float(1 AS SIGNED)
< error
| 1064 42000 …near 'float(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT float(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'float(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT float(1 FROM 2)
< error
| 1064 42000 …near 'float(1 FROM 2)' at line 1
> exec ^
: string
| SELECT float(1 IN 2)
< error
| 1064 42000 …near 'float(1 IN 2)' at line 1
> exec ^
: string
| SELECT float(DISTINCT 1)
< error
| 1064 42000 …near 'float(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT float() OVER ()
< error
| 1064 42000 …near 'float() OVER ()' at line 1
> exec ^
: string
| SELECT float(1) OVER ()
< error
| 1064 42000 …near 'float(1) OVER ()' at line 1
> exec ^
: string
| SELECT float(1,2) OVER ()
< error
| 1064 42000 …near 'float(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT float(1,2,3) OVER ()
< error
| 1064 42000 …near 'float(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT float(1) FROM t
< error
| 1064 42000 …near 'float(1) FROM t' at line 1
> exec ^
: string
| SELECT float(SELECT 1)
< error
| 1064 42000 …near 'float(SELECT 1)' at line 1
> exec ^
: string
| SELECT float(YEAR FROM 1)
< error
| 1064 42000 …near 'float(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT float(DAY, 1, 2)
< error
| 1064 42000 …near 'float(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT float(a)
< error
| 1064 42000 …near 'float(a)' at line 1
> exec ^
: string
| SELECT float(a, b)
< error
| 1064 42000 …near 'float(a, b)' at line 1
> exec ^
: string
| SELECT float(a) OVER ()
< error
| 1064 42000 …near 'float(a) OVER ()' at line 1
> exec ^
: string
| SELECT float4()
< error
| 1064 42000 …near 'float4()' at line 1
> exec ^
: string
| SELECT float4(1)
< error
| 1064 42000 …near 'float4(1)' at line 1
> exec ^
: string
| SELECT float4(1,2)
< error
| 1064 42000 …near 'float4(1,2)' at line 1
> exec ^
: string
| SELECT float4(1,2,3)
< error
| 1064 42000 …near 'float4(1,2,3)' at line 1
> exec ^
: string
| SELECT float4(1,2,3,4)
< error
| 1064 42000 …near 'float4(1,2,3,4)' at line 1
> exec ^
: string
| SELECT float4(*)
< error
| 1064 42000 …near 'float4(*)' at line 1
> exec ^
: string
| SELECT float4('a')
< error
| 1064 42000 …near 'float4('a')' at line 1
> exec ^
: string
| SELECT float4(a USING utf8mb4)
< error
| 1064 42000 …near 'float4(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT float4(1 AS SIGNED)
< error
| 1064 42000 …near 'float4(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT float4(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'float4(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT float4(1 FROM 2)
< error
| 1064 42000 …near 'float4(1 FROM 2)' at line 1
> exec ^
: string
| SELECT float4(1 IN 2)
< error
| 1064 42000 …near 'float4(1 IN 2)' at line 1
> exec ^
: string
| SELECT float4(DISTINCT 1)
< error
| 1064 42000 …near 'float4(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT float4() OVER ()
< error
| 1064 42000 …near 'float4() OVER ()' at line 1
> exec ^
: string
| SELECT float4(1) OVER ()
< error
| 1064 42000 …near 'float4(1) OVER ()' at line 1
> exec ^
: string
| SELECT float4(1,2) OVER ()
< error
| 1064 42000 …near 'float4(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT float4(1,2,3) OVER ()
< error
| 1064 42000 …near 'float4(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT float4(1) FROM t
< error
| 1064 42000 …near 'float4(1) FROM t' at line 1
> exec ^
: string
| SELECT float4(SELECT 1)
< error
| 1064 42000 …near 'float4(SELECT 1)' at line 1
> exec ^
: string
| SELECT float4(YEAR FROM 1)
< error
| 1064 42000 …near 'float4(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT float4(DAY, 1, 2)
< error
| 1064 42000 …near 'float4(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT float4(a)
< error
| 1064 42000 …near 'float4(a)' at line 1
> exec ^
: string
| SELECT float4(a, b)
< error
| 1064 42000 …near 'float4(a, b)' at line 1
> exec ^
: string
| SELECT float4(a) OVER ()
< error
| 1064 42000 …near 'float4(a) OVER ()' at line 1
> exec ^
: string
| SELECT float8()
< error
| 1064 42000 …near 'float8()' at line 1
> exec ^
: string
| SELECT float8(1)
< error
| 1064 42000 …near 'float8(1)' at line 1
> exec ^
: string
| SELECT float8(1,2)
< error
| 1064 42000 …near 'float8(1,2)' at line 1
> exec ^
: string
| SELECT float8(1,2,3)
< error
| 1064 42000 …near 'float8(1,2,3)' at line 1
> exec ^
: string
| SELECT float8(1,2,3,4)
< error
| 1064 42000 …near 'float8(1,2,3,4)' at line 1
> exec ^
: string
| SELECT float8(*)
< error
| 1064 42000 …near 'float8(*)' at line 1
> exec ^
: string
| SELECT float8('a')
< error
| 1064 42000 …near 'float8('a')' at line 1
> exec ^
: string
| SELECT float8(a USING utf8mb4)
< error
| 1064 42000 …near 'float8(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT float8(1 AS SIGNED)
< error
| 1064 42000 …near 'float8(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT float8(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'float8(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT float8(1 FROM 2)
< error
| 1064 42000 …near 'float8(1 FROM 2)' at line 1
> exec ^
: string
| SELECT float8(1 IN 2)
< error
| 1064 42000 …near 'float8(1 IN 2)' at line 1
> exec ^
: string
| SELECT float8(DISTINCT 1)
< error
| 1064 42000 …near 'float8(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT float8() OVER ()
< error
| 1064 42000 …near 'float8() OVER ()' at line 1
> exec ^
: string
| SELECT float8(1) OVER ()
< error
| 1064 42000 …near 'float8(1) OVER ()' at line 1
> exec ^
: string
| SELECT float8(1,2) OVER ()
< error
| 1064 42000 …near 'float8(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT float8(1,2,3) OVER ()
< error
| 1064 42000 …near 'float8(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT float8(1) FROM t
< error
| 1064 42000 …near 'float8(1) FROM t' at line 1
> exec ^
: string
| SELECT float8(SELECT 1)
< error
| 1064 42000 …near 'float8(SELECT 1)' at line 1
> exec ^
: string
| SELECT float8(YEAR FROM 1)
< error
| 1064 42000 …near 'float8(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT float8(DAY, 1, 2)
< error
| 1064 42000 …near 'float8(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT float8(a)
< error
| 1064 42000 …near 'float8(a)' at line 1
> exec ^
: string
| SELECT float8(a, b)
< error
| 1064 42000 …near 'float8(a, b)' at line 1
> exec ^
: string
| SELECT float8(a) OVER ()
< error
| 1064 42000 …near 'float8(a) OVER ()' at line 1
> exec ^
: string
| SELECT flush()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT flush(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT flush(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT flush(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT flush(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT flush(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT flush('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT flush(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT flush(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT flush(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT flush(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT flush(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT flush(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT flush() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT flush(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT flush(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT flush(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT flush(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT flush(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT flush(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT flush(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT flush(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT flush(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT flush(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT following()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT following(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT following(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT following(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT following(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT following(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT following('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT following(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT following(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT following(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT following(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT following(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT following(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT following() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT following(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT following(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT following(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT following(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT following(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT following(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT following(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT following(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT following(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT following(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT follows()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT follows(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT follows(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT follows(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT follows(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT follows(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT follows('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT follows(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT follows(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT follows(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT follows(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT follows(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT follows(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT follows() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT follows(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT follows(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT follows(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT follows(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT follows(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT follows(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT follows(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT follows(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT follows(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT follows(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT for()
< error
| 1064 42000 …near 'for()' at line 1
> exec ^
: string
| SELECT for(1)
< error
| 1064 42000 …near 'for(1)' at line 1
> exec ^
: string
| SELECT for(1,2)
< error
| 1064 42000 …near 'for(1,2)' at line 1
> exec ^
: string
| SELECT for(1,2,3)
< error
| 1064 42000 …near 'for(1,2,3)' at line 1
> exec ^
: string
| SELECT for(1,2,3,4)
< error
| 1064 42000 …near 'for(1,2,3,4)' at line 1
> exec ^
: string
| SELECT for(*)
< error
| 1064 42000 …near 'for(*)' at line 1
> exec ^
: string
| SELECT for('a')
< error
| 1064 42000 …near 'for('a')' at line 1
> exec ^
: string
| SELECT for(a USING utf8mb4)
< error
| 1064 42000 …near 'for(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT for(1 AS SIGNED)
< error
| 1064 42000 …near 'for(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT for(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'for(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT for(1 FROM 2)
< error
| 1064 42000 …near 'for(1 FROM 2)' at line 1
> exec ^
: string
| SELECT for(1 IN 2)
< error
| 1064 42000 …near 'for(1 IN 2)' at line 1
> exec ^
: string
| SELECT for(DISTINCT 1)
< error
| 1064 42000 …near 'for(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT for() OVER ()
< error
| 1064 42000 …near 'for() OVER ()' at line 1
> exec ^
: string
| SELECT for(1) OVER ()
< error
| 1064 42000 …near 'for(1) OVER ()' at line 1
> exec ^
: string
| SELECT for(1,2) OVER ()
< error
| 1064 42000 …near 'for(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT for(1,2,3) OVER ()
< error
| 1064 42000 …near 'for(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT for(1) FROM t
< error
| 1064 42000 …near 'for(1) FROM t' at line 1
> exec ^
: string
| SELECT for(SELECT 1)
< error
| 1064 42000 …near 'for(SELECT 1)' at line 1
> exec ^
: string
| SELECT for(YEAR FROM 1)
< error
| 1064 42000 …near 'for(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT for(DAY, 1, 2)
< error
| 1064 42000 …near 'for(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT for(a)
< error
| 1064 42000 …near 'for(a)' at line 1
> exec ^
: string
| SELECT for(a, b)
< error
| 1064 42000 …near 'for(a, b)' at line 1
> exec ^
: string
| SELECT for(a) OVER ()
< error
| 1064 42000 …near 'for(a) OVER ()' at line 1
> exec ^
: string
| SELECT force()
< error
| 1064 42000 …near 'force()' at line 1
> exec ^
: string
| SELECT force(1)
< error
| 1064 42000 …near 'force(1)' at line 1
> exec ^
: string
| SELECT force(1,2)
< error
| 1064 42000 …near 'force(1,2)' at line 1
> exec ^
: string
| SELECT force(1,2,3)
< error
| 1064 42000 …near 'force(1,2,3)' at line 1
> exec ^
: string
| SELECT force(1,2,3,4)
< error
| 1064 42000 …near 'force(1,2,3,4)' at line 1
> exec ^
: string
| SELECT force(*)
< error
| 1064 42000 …near 'force(*)' at line 1
> exec ^
: string
| SELECT force('a')
< error
| 1064 42000 …near 'force('a')' at line 1
> exec ^
: string
| SELECT force(a USING utf8mb4)
< error
| 1064 42000 …near 'force(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT force(1 AS SIGNED)
< error
| 1064 42000 …near 'force(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT force(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'force(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT force(1 FROM 2)
< error
| 1064 42000 …near 'force(1 FROM 2)' at line 1
> exec ^
: string
| SELECT force(1 IN 2)
< error
| 1064 42000 …near 'force(1 IN 2)' at line 1
> exec ^
: string
| SELECT force(DISTINCT 1)
< error
| 1064 42000 …near 'force(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT force() OVER ()
< error
| 1064 42000 …near 'force() OVER ()' at line 1
> exec ^
: string
| SELECT force(1) OVER ()
< error
| 1064 42000 …near 'force(1) OVER ()' at line 1
> exec ^
: string
| SELECT force(1,2) OVER ()
< error
| 1064 42000 …near 'force(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT force(1,2,3) OVER ()
< error
| 1064 42000 …near 'force(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT force(1) FROM t
< error
| 1064 42000 …near 'force(1) FROM t' at line 1
> exec ^
: string
| SELECT force(SELECT 1)
< error
| 1064 42000 …near 'force(SELECT 1)' at line 1
> exec ^
: string
| SELECT force(YEAR FROM 1)
< error
| 1064 42000 …near 'force(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT force(DAY, 1, 2)
< error
| 1064 42000 …near 'force(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT force(a)
< error
| 1064 42000 …near 'force(a)' at line 1
> exec ^
: string
| SELECT force(a, b)
< error
| 1064 42000 …near 'force(a, b)' at line 1
> exec ^
: string
| SELECT force(a) OVER ()
< error
| 1064 42000 …near 'force(a) OVER ()' at line 1
> exec ^
: string
| SELECT foreign()
< error
| 1064 42000 …near 'foreign()' at line 1
> exec ^
: string
| SELECT foreign(1)
< error
| 1064 42000 …near 'foreign(1)' at line 1
> exec ^
: string
| SELECT foreign(1,2)
< error
| 1064 42000 …near 'foreign(1,2)' at line 1
> exec ^
: string
| SELECT foreign(1,2,3)
< error
| 1064 42000 …near 'foreign(1,2,3)' at line 1
> exec ^
: string
| SELECT foreign(1,2,3,4)
< error
| 1064 42000 …near 'foreign(1,2,3,4)' at line 1
> exec ^
: string
| SELECT foreign(*)
< error
| 1064 42000 …near 'foreign(*)' at line 1
> exec ^
: string
| SELECT foreign('a')
< error
| 1064 42000 …near 'foreign('a')' at line 1
> exec ^
: string
| SELECT foreign(a USING utf8mb4)
< error
| 1064 42000 …near 'foreign(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT foreign(1 AS SIGNED)
< error
| 1064 42000 …near 'foreign(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT foreign(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'foreign(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT foreign(1 FROM 2)
< error
| 1064 42000 …near 'foreign(1 FROM 2)' at line 1
> exec ^
: string
| SELECT foreign(1 IN 2)
< error
| 1064 42000 …near 'foreign(1 IN 2)' at line 1
> exec ^
: string
| SELECT foreign(DISTINCT 1)
< error
| 1064 42000 …near 'foreign(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT foreign() OVER ()
< error
| 1064 42000 …near 'foreign() OVER ()' at line 1
> exec ^
: string
| SELECT foreign(1) OVER ()
< error
| 1064 42000 …near 'foreign(1) OVER ()' at line 1
> exec ^
: string
| SELECT foreign(1,2) OVER ()
< error
| 1064 42000 …near 'foreign(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT foreign(1,2,3) OVER ()
< error
| 1064 42000 …near 'foreign(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT foreign(1) FROM t
< error
| 1064 42000 …near 'foreign(1) FROM t' at line 1
> exec ^
: string
| SELECT foreign(SELECT 1)
< error
| 1064 42000 …near 'foreign(SELECT 1)' at line 1
> exec ^
: string
| SELECT foreign(YEAR FROM 1)
< error
| 1064 42000 …near 'foreign(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT foreign(DAY, 1, 2)
< error
| 1064 42000 …near 'foreign(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT foreign(a)
< error
| 1064 42000 …near 'foreign(a)' at line 1
> exec ^
: string
| SELECT foreign(a, b)
< error
| 1064 42000 …near 'foreign(a, b)' at line 1
> exec ^
: string
| SELECT foreign(a) OVER ()
< error
| 1064 42000 …near 'foreign(a) OVER ()' at line 1
> exec ^
: string
| SELECT format()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'format'
> exec ^
: string
| SELECT format()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'format'
> exec ^
: string
| SELECT format(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'format'
> exec ^
: string
| SELECT format(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'format'
> exec ^
: string
| SELECT format(1,2)
< ok 0 0
> exec ^
: string
| SELECT found()
< error
| 1630 42000 FUNCTION test.found does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT found()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT from()
< error
| 1064 42000 …near 'from()' at line 1
> exec ^
: string
| SELECT from(1)
< error
| 1064 42000 …near 'from(1)' at line 1
> exec ^
: string
| SELECT from(1,2)
< error
| 1064 42000 …near 'from(1,2)' at line 1
> exec ^
: string
| SELECT from(1,2,3)
< error
| 1064 42000 …near 'from(1,2,3)' at line 1
> exec ^
: string
| SELECT from(1,2,3,4)
< error
| 1064 42000 …near 'from(1,2,3,4)' at line 1
> exec ^
: string
| SELECT from(*)
< error
| 1064 42000 …near 'from(*)' at line 1
> exec ^
: string
| SELECT from('a')
< error
| 1064 42000 …near 'from('a')' at line 1
> exec ^
: string
| SELECT from(a USING utf8mb4)
< error
| 1064 42000 …near 'from(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT from(1 AS SIGNED)
< error
| 1064 42000 …near 'from(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT from(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'from(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT from(1 FROM 2)
< error
| 1064 42000 …near 'from(1 FROM 2)' at line 1
> exec ^
: string
| SELECT from(1 IN 2)
< error
| 1064 42000 …near 'from(1 IN 2)' at line 1
> exec ^
: string
| SELECT from(DISTINCT 1)
< error
| 1064 42000 …near 'from(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT from() OVER ()
< error
| 1064 42000 …near 'from() OVER ()' at line 1
> exec ^
: string
| SELECT from(1) OVER ()
< error
| 1064 42000 …near 'from(1) OVER ()' at line 1
> exec ^
: string
| SELECT from(1,2) OVER ()
< error
| 1064 42000 …near 'from(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT from(1,2,3) OVER ()
< error
| 1064 42000 …near 'from(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT from(1) FROM t
< error
| 1064 42000 …near 'from(1) FROM t' at line 1
> exec ^
: string
| SELECT from(SELECT 1)
< error
| 1064 42000 …near 'from(SELECT 1)' at line 1
> exec ^
: string
| SELECT from(YEAR FROM 1)
< error
| 1064 42000 …near 'from(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT from(DAY, 1, 2)
< error
| 1064 42000 …near 'from(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT from(a)
< error
| 1064 42000 …near 'from(a)' at line 1
> exec ^
: string
| SELECT from(a, b)
< error
| 1064 42000 …near 'from(a, b)' at line 1
> exec ^
: string
| SELECT from(a) OVER ()
< error
| 1064 42000 …near 'from(a) OVER ()' at line 1
> exec ^
: string
| SELECT full()
< error
| 1630 42000 FUNCTION test.full does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT full()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT fulltext()
< error
| 1064 42000 …near 'fulltext()' at line 1
> exec ^
: string
| SELECT fulltext(1)
< error
| 1064 42000 …near 'fulltext(1)' at line 1
> exec ^
: string
| SELECT fulltext(1,2)
< error
| 1064 42000 …near 'fulltext(1,2)' at line 1
> exec ^
: string
| SELECT fulltext(1,2,3)
< error
| 1064 42000 …near 'fulltext(1,2,3)' at line 1
> exec ^
: string
| SELECT fulltext(1,2,3,4)
< error
| 1064 42000 …near 'fulltext(1,2,3,4)' at line 1
> exec ^
: string
| SELECT fulltext(*)
< error
| 1064 42000 …near 'fulltext(*)' at line 1
> exec ^
: string
| SELECT fulltext('a')
< error
| 1064 42000 …near 'fulltext('a')' at line 1
> exec ^
: string
| SELECT fulltext(a USING utf8mb4)
< error
| 1064 42000 …near 'fulltext(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT fulltext(1 AS SIGNED)
< error
| 1064 42000 …near 'fulltext(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT fulltext(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'fulltext(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT fulltext(1 FROM 2)
< error
| 1064 42000 …near 'fulltext(1 FROM 2)' at line 1
> exec ^
: string
| SELECT fulltext(1 IN 2)
< error
| 1064 42000 …near 'fulltext(1 IN 2)' at line 1
> exec ^
: string
| SELECT fulltext(DISTINCT 1)
< error
| 1064 42000 …near 'fulltext(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT fulltext() OVER ()
< error
| 1064 42000 …near 'fulltext() OVER ()' at line 1
> exec ^
: string
| SELECT fulltext(1) OVER ()
< error
| 1064 42000 …near 'fulltext(1) OVER ()' at line 1
> exec ^
: string
| SELECT fulltext(1,2) OVER ()
< error
| 1064 42000 …near 'fulltext(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT fulltext(1,2,3) OVER ()
< error
| 1064 42000 …near 'fulltext(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT fulltext(1) FROM t
< error
| 1064 42000 …near 'fulltext(1) FROM t' at line 1
> exec ^
: string
| SELECT fulltext(SELECT 1)
< error
| 1064 42000 …near 'fulltext(SELECT 1)' at line 1
> exec ^
: string
| SELECT fulltext(YEAR FROM 1)
< error
| 1064 42000 …near 'fulltext(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT fulltext(DAY, 1, 2)
< error
| 1064 42000 …near 'fulltext(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT fulltext(a)
< error
| 1064 42000 …near 'fulltext(a)' at line 1
> exec ^
: string
| SELECT fulltext(a, b)
< error
| 1064 42000 …near 'fulltext(a, b)' at line 1
> exec ^
: string
| SELECT fulltext(a) OVER ()
< error
| 1064 42000 …near 'fulltext(a) OVER ()' at line 1
> exec ^
: string
| SELECT function()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT function(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT function(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT function(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT function(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT function(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT function('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT function(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT function(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT function(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT function(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT function(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT function(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT function() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT function(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT function(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT function(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT function(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT function(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT function(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT function(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT function(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT function(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT function(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT general()
< error
| 1630 42000 FUNCTION test.general does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT general()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT generate()
< error
| 1305 42000 FUNCTION test.generate does not exist
> exec ^
: string
| SELECT generate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT generated()
< error
| 1630 42000 FUNCTION test.generated does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT generated()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT geomcollection()
< error
| 1305 42000 FUNCTION test.geomcollection does not exist
> exec ^
: string
| SELECT geomcollection()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT geometry()
< error
| 1305 42000 FUNCTION test.geometry does not exist
> exec ^
: string
| SELECT geometry()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT geometrycollection()
< error
| 1305 42000 FUNCTION test.geometrycollection does not exist
> exec ^
: string
| SELECT geometrycollection()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT get()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT get(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT get(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT get(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT get(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT get(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT get('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT get(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT get(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT get(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT get(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT get(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT get(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT get() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT get(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT get(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT get(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT get(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT get(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT get(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT get(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT get(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT get(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT get(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT get_format()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT get_format(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT get_format(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT get_format(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT get_format(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT get_format(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT get_format('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT get_format(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT get_format(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT get_format(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT get_format(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT get_format(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT get_format(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT get_format() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT get_format(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT get_format(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT get_format(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT get_format(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT get_format(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT get_format(YEAR FROM 1)
< error
| 1064 42000 …near 'YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT get_format(DAY, 1, 2)
< error
| 1064 42000 …near 'DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT get_format(a)
< error
| 1064 42000 …near 'a)' at line 1
> exec ^
: string
| SELECT get_format(a, b)
< error
| 1064 42000 …near 'a, b)' at line 1
> exec ^
: string
| SELECT get_format(a) OVER ()
< error
| 1064 42000 …near 'a) OVER ()' at line 1
> exec ^
: string
| SELECT get_master_public_key()
< error
| 1305 42000 FUNCTION test.get_master_public_key does not exist
> exec ^
: string
| SELECT get_master_public_key()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT get_source_public_key()
< error
| 1305 42000 FUNCTION test.get_source_public_key does not exist
> exec ^
: string
| SELECT get_source_public_key()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT global()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT global(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT global(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT global(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT global(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT global(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT global('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT global(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT global(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT global(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT global(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT global(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT global(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT global() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT global(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT global(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT global(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT global(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT global(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT global(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT global(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT global(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT global(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT global(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT grant()
< error
| 1064 42000 …near 'grant()' at line 1
> exec ^
: string
| SELECT grant(1)
< error
| 1064 42000 …near 'grant(1)' at line 1
> exec ^
: string
| SELECT grant(1,2)
< error
| 1064 42000 …near 'grant(1,2)' at line 1
> exec ^
: string
| SELECT grant(1,2,3)
< error
| 1064 42000 …near 'grant(1,2,3)' at line 1
> exec ^
: string
| SELECT grant(1,2,3,4)
< error
| 1064 42000 …near 'grant(1,2,3,4)' at line 1
> exec ^
: string
| SELECT grant(*)
< error
| 1064 42000 …near 'grant(*)' at line 1
> exec ^
: string
| SELECT grant('a')
< error
| 1064 42000 …near 'grant('a')' at line 1
> exec ^
: string
| SELECT grant(a USING utf8mb4)
< error
| 1064 42000 …near 'grant(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT grant(1 AS SIGNED)
< error
| 1064 42000 …near 'grant(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT grant(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'grant(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT grant(1 FROM 2)
< error
| 1064 42000 …near 'grant(1 FROM 2)' at line 1
> exec ^
: string
| SELECT grant(1 IN 2)
< error
| 1064 42000 …near 'grant(1 IN 2)' at line 1
> exec ^
: string
| SELECT grant(DISTINCT 1)
< error
| 1064 42000 …near 'grant(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT grant() OVER ()
< error
| 1064 42000 …near 'grant() OVER ()' at line 1
> exec ^
: string
| SELECT grant(1) OVER ()
< error
| 1064 42000 …near 'grant(1) OVER ()' at line 1
> exec ^
: string
| SELECT grant(1,2) OVER ()
< error
| 1064 42000 …near 'grant(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT grant(1,2,3) OVER ()
< error
| 1064 42000 …near 'grant(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT grant(1) FROM t
< error
| 1064 42000 …near 'grant(1) FROM t' at line 1
> exec ^
: string
| SELECT grant(SELECT 1)
< error
| 1064 42000 …near 'grant(SELECT 1)' at line 1
> exec ^
: string
| SELECT grant(YEAR FROM 1)
< error
| 1064 42000 …near 'grant(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT grant(DAY, 1, 2)
< error
| 1064 42000 …near 'grant(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT grant(a)
< error
| 1064 42000 …near 'grant(a)' at line 1
> exec ^
: string
| SELECT grant(a, b)
< error
| 1064 42000 …near 'grant(a, b)' at line 1
> exec ^
: string
| SELECT grant(a) OVER ()
< error
| 1064 42000 …near 'grant(a) OVER ()' at line 1
> exec ^
: string
| SELECT grants()
< error
| 1630 42000 FUNCTION test.grants does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT grants()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT greatest()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'greatest'
> exec ^
: string
| SELECT greatest()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'greatest'
> exec ^
: string
| SELECT greatest(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'greatest'
> exec ^
: string
| SELECT greatest(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'greatest'
> exec ^
: string
| SELECT greatest(1,2)
< ok 0 0
> exec ^
: string
| SELECT group()
< error
| 1064 42000 …near 'group()' at line 1
> exec ^
: string
| SELECT group(1)
< error
| 1064 42000 …near 'group(1)' at line 1
> exec ^
: string
| SELECT group(1,2)
< error
| 1064 42000 …near 'group(1,2)' at line 1
> exec ^
: string
| SELECT group(1,2,3)
< error
| 1064 42000 …near 'group(1,2,3)' at line 1
> exec ^
: string
| SELECT group(1,2,3,4)
< error
| 1064 42000 …near 'group(1,2,3,4)' at line 1
> exec ^
: string
| SELECT group(*)
< error
| 1064 42000 …near 'group(*)' at line 1
> exec ^
: string
| SELECT group('a')
< error
| 1064 42000 …near 'group('a')' at line 1
> exec ^
: string
| SELECT group(a USING utf8mb4)
< error
| 1064 42000 …near 'group(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT group(1 AS SIGNED)
< error
| 1064 42000 …near 'group(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT group(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'group(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT group(1 FROM 2)
< error
| 1064 42000 …near 'group(1 FROM 2)' at line 1
> exec ^
: string
| SELECT group(1 IN 2)
< error
| 1064 42000 …near 'group(1 IN 2)' at line 1
> exec ^
: string
| SELECT group(DISTINCT 1)
< error
| 1064 42000 …near 'group(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT group() OVER ()
< error
| 1064 42000 …near 'group() OVER ()' at line 1
> exec ^
: string
| SELECT group(1) OVER ()
< error
| 1064 42000 …near 'group(1) OVER ()' at line 1
> exec ^
: string
| SELECT group(1,2) OVER ()
< error
| 1064 42000 …near 'group(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT group(1,2,3) OVER ()
< error
| 1064 42000 …near 'group(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT group(1) FROM t
< error
| 1064 42000 …near 'group(1) FROM t' at line 1
> exec ^
: string
| SELECT group(SELECT 1)
< error
| 1064 42000 …near 'group(SELECT 1)' at line 1
> exec ^
: string
| SELECT group(YEAR FROM 1)
< error
| 1064 42000 …near 'group(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT group(DAY, 1, 2)
< error
| 1064 42000 …near 'group(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT group(a)
< error
| 1064 42000 …near 'group(a)' at line 1
> exec ^
: string
| SELECT group(a, b)
< error
| 1064 42000 …near 'group(a, b)' at line 1
> exec ^
: string
| SELECT group(a) OVER ()
< error
| 1064 42000 …near 'group(a) OVER ()' at line 1
> exec ^
: string
| SELECT group_concat()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT group_concat(1)
< ok 0 0
> exec ^
: string
| SELECT group_replication()
< error
| 1305 42000 FUNCTION test.group_replication does not exist
> exec ^
: string
| SELECT group_replication()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT grouping()
< error
| 1305 42000 FUNCTION test.grouping does not exist
> exec ^
: string
| SELECT grouping()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT groups()
< error
| 1305 42000 FUNCTION test.groups does not exist
> exec ^
: string
| SELECT groups()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT gtid_only()
< error
| 1305 42000 FUNCTION test.gtid_only does not exist
> exec ^
: string
| SELECT gtid_only()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT gtids()
< error
| 1305 42000 FUNCTION test.gtids does not exist
> exec ^
: string
| SELECT gtids()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT guided()
< error
| 1305 42000 FUNCTION test.guided does not exist
> exec ^
: string
| SELECT guided()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT handler()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT handler(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT handler(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT handler(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT handler(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT handler(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT handler('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT handler(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT handler(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT handler(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT handler(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT handler(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT handler(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT handler() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT handler(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT handler(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT handler(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT handler(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT handler(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT handler(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT handler(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT handler(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT handler(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT handler(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT hash()
< error
| 1630 42000 FUNCTION test.hash does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT hash()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT having()
< error
| 1064 42000 …near 'having()' at line 1
> exec ^
: string
| SELECT having(1)
< error
| 1064 42000 …near 'having(1)' at line 1
> exec ^
: string
| SELECT having(1,2)
< error
| 1064 42000 …near 'having(1,2)' at line 1
> exec ^
: string
| SELECT having(1,2,3)
< error
| 1064 42000 …near 'having(1,2,3)' at line 1
> exec ^
: string
| SELECT having(1,2,3,4)
< error
| 1064 42000 …near 'having(1,2,3,4)' at line 1
> exec ^
: string
| SELECT having(*)
< error
| 1064 42000 …near 'having(*)' at line 1
> exec ^
: string
| SELECT having('a')
< error
| 1064 42000 …near 'having('a')' at line 1
> exec ^
: string
| SELECT having(a USING utf8mb4)
< error
| 1064 42000 …near 'having(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT having(1 AS SIGNED)
< error
| 1064 42000 …near 'having(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT having(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'having(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT having(1 FROM 2)
< error
| 1064 42000 …near 'having(1 FROM 2)' at line 1
> exec ^
: string
| SELECT having(1 IN 2)
< error
| 1064 42000 …near 'having(1 IN 2)' at line 1
> exec ^
: string
| SELECT having(DISTINCT 1)
< error
| 1064 42000 …near 'having(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT having() OVER ()
< error
| 1064 42000 …near 'having() OVER ()' at line 1
> exec ^
: string
| SELECT having(1) OVER ()
< error
| 1064 42000 …near 'having(1) OVER ()' at line 1
> exec ^
: string
| SELECT having(1,2) OVER ()
< error
| 1064 42000 …near 'having(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT having(1,2,3) OVER ()
< error
| 1064 42000 …near 'having(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT having(1) FROM t
< error
| 1064 42000 …near 'having(1) FROM t' at line 1
> exec ^
: string
| SELECT having(SELECT 1)
< error
| 1064 42000 …near 'having(SELECT 1)' at line 1
> exec ^
: string
| SELECT having(YEAR FROM 1)
< error
| 1064 42000 …near 'having(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT having(DAY, 1, 2)
< error
| 1064 42000 …near 'having(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT having(a)
< error
| 1064 42000 …near 'having(a)' at line 1
> exec ^
: string
| SELECT having(a, b)
< error
| 1064 42000 …near 'having(a, b)' at line 1
> exec ^
: string
| SELECT having(a) OVER ()
< error
| 1064 42000 …near 'having(a) OVER ()' at line 1
> exec ^
: string
| SELECT header()
< error
| 1305 42000 FUNCTION test.header does not exist
> exec ^
: string
| SELECT header()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT help()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT help(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT help(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT help(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT help(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT help(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT help('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT help(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT help(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT help(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT help(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT help(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT help(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT help() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT help(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT help(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT help(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT help(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT help(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT help(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT help(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT help(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT help(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT help(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT high_priority()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT high_priority(1)
< ok 0 0
> exec ^
: string
| SELECT histogram()
< error
| 1305 42000 FUNCTION test.histogram does not exist
> exec ^
: string
| SELECT histogram()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT history()
< error
| 1630 42000 FUNCTION test.history does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT history()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT host()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT host(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT host(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT host(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT host(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT host(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT host('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT host(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT host(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT host(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT host(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT host(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT host(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT host() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT host(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT host(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT host(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT host(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT host(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT host(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT host(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT host(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT host(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT host(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT hosts()
< error
| 1630 42000 FUNCTION test.hosts does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT hosts()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT hour()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT hour(1)
< ok 0 0
> exec ^
: string
| SELECT hour_microsecond()
< error
| 1064 42000 …near 'hour_microsecond()' at line 1
> exec ^
: string
| SELECT hour_microsecond(1)
< error
| 1064 42000 …near 'hour_microsecond(1)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1,2)
< error
| 1064 42000 …near 'hour_microsecond(1,2)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1,2,3)
< error
| 1064 42000 …near 'hour_microsecond(1,2,3)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1,2,3,4)
< error
| 1064 42000 …near 'hour_microsecond(1,2,3,4)' at line 1
> exec ^
: string
| SELECT hour_microsecond(*)
< error
| 1064 42000 …near 'hour_microsecond(*)' at line 1
> exec ^
: string
| SELECT hour_microsecond('a')
< error
| 1064 42000 …near 'hour_microsecond('a')' at line 1
> exec ^
: string
| SELECT hour_microsecond(a USING utf8mb4)
< error
| 1064 42000 …near 'hour_microsecond(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1 AS SIGNED)
< error
| 1064 42000 …near 'hour_microsecond(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'hour_microsecond(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1 FROM 2)
< error
| 1064 42000 …near 'hour_microsecond(1 FROM 2)' at line 1
> exec ^
: string
| SELECT hour_microsecond(1 IN 2)
< error
| 1064 42000 …near 'hour_microsecond(1 IN 2)' at line 1
> exec ^
: string
| SELECT hour_microsecond(DISTINCT 1)
< error
| 1064 42000 …near 'hour_microsecond(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT hour_microsecond() OVER ()
< error
| 1064 42000 …near 'hour_microsecond() OVER ()' at line 1
> exec ^
: string
| SELECT hour_microsecond(1) OVER ()
< error
| 1064 42000 …near 'hour_microsecond(1) OVER ()' at line 1
> exec ^
: string
| SELECT hour_microsecond(1,2) OVER ()
< error
| 1064 42000 …near 'hour_microsecond(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT hour_microsecond(1,2,3) OVER ()
< error
| 1064 42000 …near 'hour_microsecond(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT hour_microsecond(1) FROM t
< error
| 1064 42000 …near 'hour_microsecond(1) FROM t' at line 1
> exec ^
: string
| SELECT hour_microsecond(SELECT 1)
< error
| 1064 42000 …near 'hour_microsecond(SELECT 1)' at line 1
> exec ^
: string
| SELECT hour_microsecond(YEAR FROM 1)
< error
| 1064 42000 …near 'hour_microsecond(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT hour_microsecond(DAY, 1, 2)
< error
| 1064 42000 …near 'hour_microsecond(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT hour_microsecond(a)
< error
| 1064 42000 …near 'hour_microsecond(a)' at line 1
> exec ^
: string
| SELECT hour_microsecond(a, b)
< error
| 1064 42000 …near 'hour_microsecond(a, b)' at line 1
> exec ^
: string
| SELECT hour_microsecond(a) OVER ()
< error
| 1064 42000 …near 'hour_microsecond(a) OVER ()' at line 1
> exec ^
: string
| SELECT hour_minute()
< error
| 1064 42000 …near 'hour_minute()' at line 1
> exec ^
: string
| SELECT hour_minute(1)
< error
| 1064 42000 …near 'hour_minute(1)' at line 1
> exec ^
: string
| SELECT hour_minute(1,2)
< error
| 1064 42000 …near 'hour_minute(1,2)' at line 1
> exec ^
: string
| SELECT hour_minute(1,2,3)
< error
| 1064 42000 …near 'hour_minute(1,2,3)' at line 1
> exec ^
: string
| SELECT hour_minute(1,2,3,4)
< error
| 1064 42000 …near 'hour_minute(1,2,3,4)' at line 1
> exec ^
: string
| SELECT hour_minute(*)
< error
| 1064 42000 …near 'hour_minute(*)' at line 1
> exec ^
: string
| SELECT hour_minute('a')
< error
| 1064 42000 …near 'hour_minute('a')' at line 1
> exec ^
: string
| SELECT hour_minute(a USING utf8mb4)
< error
| 1064 42000 …near 'hour_minute(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT hour_minute(1 AS SIGNED)
< error
| 1064 42000 …near 'hour_minute(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT hour_minute(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'hour_minute(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT hour_minute(1 FROM 2)
< error
| 1064 42000 …near 'hour_minute(1 FROM 2)' at line 1
> exec ^
: string
| SELECT hour_minute(1 IN 2)
< error
| 1064 42000 …near 'hour_minute(1 IN 2)' at line 1
> exec ^
: string
| SELECT hour_minute(DISTINCT 1)
< error
| 1064 42000 …near 'hour_minute(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT hour_minute() OVER ()
< error
| 1064 42000 …near 'hour_minute() OVER ()' at line 1
> exec ^
: string
| SELECT hour_minute(1) OVER ()
< error
| 1064 42000 …near 'hour_minute(1) OVER ()' at line 1
> exec ^
: string
| SELECT hour_minute(1,2) OVER ()
< error
| 1064 42000 …near 'hour_minute(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT hour_minute(1,2,3) OVER ()
< error
| 1064 42000 …near 'hour_minute(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT hour_minute(1) FROM t
< error
| 1064 42000 …near 'hour_minute(1) FROM t' at line 1
> exec ^
: string
| SELECT hour_minute(SELECT 1)
< error
| 1064 42000 …near 'hour_minute(SELECT 1)' at line 1
> exec ^
: string
| SELECT hour_minute(YEAR FROM 1)
< error
| 1064 42000 …near 'hour_minute(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT hour_minute(DAY, 1, 2)
< error
| 1064 42000 …near 'hour_minute(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT hour_minute(a)
< error
| 1064 42000 …near 'hour_minute(a)' at line 1
> exec ^
: string
| SELECT hour_minute(a, b)
< error
| 1064 42000 …near 'hour_minute(a, b)' at line 1
> exec ^
: string
| SELECT hour_minute(a) OVER ()
< error
| 1064 42000 …near 'hour_minute(a) OVER ()' at line 1
> exec ^
: string
| SELECT hour_second()
< error
| 1064 42000 …near 'hour_second()' at line 1
> exec ^
: string
| SELECT hour_second(1)
< error
| 1064 42000 …near 'hour_second(1)' at line 1
> exec ^
: string
| SELECT hour_second(1,2)
< error
| 1064 42000 …near 'hour_second(1,2)' at line 1
> exec ^
: string
| SELECT hour_second(1,2,3)
< error
| 1064 42000 …near 'hour_second(1,2,3)' at line 1
> exec ^
: string
| SELECT hour_second(1,2,3,4)
< error
| 1064 42000 …near 'hour_second(1,2,3,4)' at line 1
> exec ^
: string
| SELECT hour_second(*)
< error
| 1064 42000 …near 'hour_second(*)' at line 1
> exec ^
: string
| SELECT hour_second('a')
< error
| 1064 42000 …near 'hour_second('a')' at line 1
> exec ^
: string
| SELECT hour_second(a USING utf8mb4)
< error
| 1064 42000 …near 'hour_second(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT hour_second(1 AS SIGNED)
< error
| 1064 42000 …near 'hour_second(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT hour_second(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'hour_second(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT hour_second(1 FROM 2)
< error
| 1064 42000 …near 'hour_second(1 FROM 2)' at line 1
> exec ^
: string
| SELECT hour_second(1 IN 2)
< error
| 1064 42000 …near 'hour_second(1 IN 2)' at line 1
> exec ^
: string
| SELECT hour_second(DISTINCT 1)
< error
| 1064 42000 …near 'hour_second(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT hour_second() OVER ()
< error
| 1064 42000 …near 'hour_second() OVER ()' at line 1
> exec ^
: string
| SELECT hour_second(1) OVER ()
< error
| 1064 42000 …near 'hour_second(1) OVER ()' at line 1
> exec ^
: string
| SELECT hour_second(1,2) OVER ()
< error
| 1064 42000 …near 'hour_second(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT hour_second(1,2,3) OVER ()
< error
| 1064 42000 …near 'hour_second(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT hour_second(1) FROM t
< error
| 1064 42000 …near 'hour_second(1) FROM t' at line 1
> exec ^
: string
| SELECT hour_second(SELECT 1)
< error
| 1064 42000 …near 'hour_second(SELECT 1)' at line 1
> exec ^
: string
| SELECT hour_second(YEAR FROM 1)
< error
| 1064 42000 …near 'hour_second(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT hour_second(DAY, 1, 2)
< error
| 1064 42000 …near 'hour_second(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT hour_second(a)
< error
| 1064 42000 …near 'hour_second(a)' at line 1
> exec ^
: string
| SELECT hour_second(a, b)
< error
| 1064 42000 …near 'hour_second(a, b)' at line 1
> exec ^
: string
| SELECT hour_second(a) OVER ()
< error
| 1064 42000 …near 'hour_second(a) OVER ()' at line 1
> exec ^
: string
| SELECT identified()
< error
| 1630 42000 FUNCTION test.identified does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT identified()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT if()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT if(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT if(1,2)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT if(1,2,3)
< ok 0 0
> exec ^
: string
| SELECT ignore()
< error
| 1064 42000 …near 'ignore()' at line 1
> exec ^
: string
| SELECT ignore(1)
< error
| 1064 42000 …near 'ignore(1)' at line 1
> exec ^
: string
| SELECT ignore(1,2)
< error
| 1064 42000 …near 'ignore(1,2)' at line 1
> exec ^
: string
| SELECT ignore(1,2,3)
< error
| 1064 42000 …near 'ignore(1,2,3)' at line 1
> exec ^
: string
| SELECT ignore(1,2,3,4)
< error
| 1064 42000 …near 'ignore(1,2,3,4)' at line 1
> exec ^
: string
| SELECT ignore(*)
< error
| 1064 42000 …near 'ignore(*)' at line 1
> exec ^
: string
| SELECT ignore('a')
< error
| 1064 42000 …near 'ignore('a')' at line 1
> exec ^
: string
| SELECT ignore(a USING utf8mb4)
< error
| 1064 42000 …near 'ignore(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT ignore(1 AS SIGNED)
< error
| 1064 42000 …near 'ignore(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT ignore(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'ignore(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT ignore(1 FROM 2)
< error
| 1064 42000 …near 'ignore(1 FROM 2)' at line 1
> exec ^
: string
| SELECT ignore(1 IN 2)
< error
| 1064 42000 …near 'ignore(1 IN 2)' at line 1
> exec ^
: string
| SELECT ignore(DISTINCT 1)
< error
| 1064 42000 …near 'ignore(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT ignore() OVER ()
< error
| 1064 42000 …near 'ignore() OVER ()' at line 1
> exec ^
: string
| SELECT ignore(1) OVER ()
< error
| 1064 42000 …near 'ignore(1) OVER ()' at line 1
> exec ^
: string
| SELECT ignore(1,2) OVER ()
< error
| 1064 42000 …near 'ignore(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT ignore(1,2,3) OVER ()
< error
| 1064 42000 …near 'ignore(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT ignore(1) FROM t
< error
| 1064 42000 …near 'ignore(1) FROM t' at line 1
> exec ^
: string
| SELECT ignore(SELECT 1)
< error
| 1064 42000 …near 'ignore(SELECT 1)' at line 1
> exec ^
: string
| SELECT ignore(YEAR FROM 1)
< error
| 1064 42000 …near 'ignore(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT ignore(DAY, 1, 2)
< error
| 1064 42000 …near 'ignore(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT ignore(a)
< error
| 1064 42000 …near 'ignore(a)' at line 1
> exec ^
: string
| SELECT ignore(a, b)
< error
| 1064 42000 …near 'ignore(a, b)' at line 1
> exec ^
: string
| SELECT ignore(a) OVER ()
< error
| 1064 42000 …near 'ignore(a) OVER ()' at line 1
> exec ^
: string
| SELECT ignore_server_ids()
< error
| 1630 42000 FUNCTION test.ignore_server_ids does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT ignore_server_ids()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT import()
< error
| 1630 42000 FUNCTION test.import does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT import()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT in()
< error
| 1064 42000 …near 'in()' at line 1
> exec ^
: string
| SELECT in(1)
< error
| 1064 42000 …near 'in(1)' at line 1
> exec ^
: string
| SELECT in(1,2)
< error
| 1064 42000 …near 'in(1,2)' at line 1
> exec ^
: string
| SELECT in(1,2,3)
< error
| 1064 42000 …near 'in(1,2,3)' at line 1
> exec ^
: string
| SELECT in(1,2,3,4)
< error
| 1064 42000 …near 'in(1,2,3,4)' at line 1
> exec ^
: string
| SELECT in(*)
< error
| 1064 42000 …near 'in(*)' at line 1
> exec ^
: string
| SELECT in('a')
< error
| 1064 42000 …near 'in('a')' at line 1
> exec ^
: string
| SELECT in(a USING utf8mb4)
< error
| 1064 42000 …near 'in(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT in(1 AS SIGNED)
< error
| 1064 42000 …near 'in(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT in(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'in(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT in(1 FROM 2)
< error
| 1064 42000 …near 'in(1 FROM 2)' at line 1
> exec ^
: string
| SELECT in(1 IN 2)
< error
| 1064 42000 …near 'in(1 IN 2)' at line 1
> exec ^
: string
| SELECT in(DISTINCT 1)
< error
| 1064 42000 …near 'in(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT in() OVER ()
< error
| 1064 42000 …near 'in() OVER ()' at line 1
> exec ^
: string
| SELECT in(1) OVER ()
< error
| 1064 42000 …near 'in(1) OVER ()' at line 1
> exec ^
: string
| SELECT in(1,2) OVER ()
< error
| 1064 42000 …near 'in(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT in(1,2,3) OVER ()
< error
| 1064 42000 …near 'in(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT in(1) FROM t
< error
| 1064 42000 …near 'in(1) FROM t' at line 1
> exec ^
: string
| SELECT in(SELECT 1)
< error
| 1064 42000 …near 'in(SELECT 1)' at line 1
> exec ^
: string
| SELECT in(YEAR FROM 1)
< error
| 1064 42000 …near 'in(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT in(DAY, 1, 2)
< error
| 1064 42000 …near 'in(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT in(a)
< error
| 1064 42000 …near 'in(a)' at line 1
> exec ^
: string
| SELECT in(a, b)
< error
| 1064 42000 …near 'in(a, b)' at line 1
> exec ^
: string
| SELECT in(a) OVER ()
< error
| 1064 42000 …near 'in(a) OVER ()' at line 1
> exec ^
: string
| SELECT inactive()
< error
| 1305 42000 FUNCTION test.inactive does not exist
> exec ^
: string
| SELECT inactive()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT increment()
< error
| 1630 42000 FUNCTION test.increment does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT increment()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT index()
< error
| 1064 42000 …near 'index()' at line 1
> exec ^
: string
| SELECT index(1)
< error
| 1064 42000 …near 'index(1)' at line 1
> exec ^
: string
| SELECT index(1,2)
< error
| 1064 42000 …near 'index(1,2)' at line 1
> exec ^
: string
| SELECT index(1,2,3)
< error
| 1064 42000 …near 'index(1,2,3)' at line 1
> exec ^
: string
| SELECT index(1,2,3,4)
< error
| 1064 42000 …near 'index(1,2,3,4)' at line 1
> exec ^
: string
| SELECT index(*)
< error
| 1064 42000 …near 'index(*)' at line 1
> exec ^
: string
| SELECT index('a')
< error
| 1064 42000 …near 'index('a')' at line 1
> exec ^
: string
| SELECT index(a USING utf8mb4)
< error
| 1064 42000 …near 'index(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT index(1 AS SIGNED)
< error
| 1064 42000 …near 'index(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT index(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'index(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT index(1 FROM 2)
< error
| 1064 42000 …near 'index(1 FROM 2)' at line 1
> exec ^
: string
| SELECT index(1 IN 2)
< error
| 1064 42000 …near 'index(1 IN 2)' at line 1
> exec ^
: string
| SELECT index(DISTINCT 1)
< error
| 1064 42000 …near 'index(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT index() OVER ()
< error
| 1064 42000 …near 'index() OVER ()' at line 1
> exec ^
: string
| SELECT index(1) OVER ()
< error
| 1064 42000 …near 'index(1) OVER ()' at line 1
> exec ^
: string
| SELECT index(1,2) OVER ()
< error
| 1064 42000 …near 'index(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT index(1,2,3) OVER ()
< error
| 1064 42000 …near 'index(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT index(1) FROM t
< error
| 1064 42000 …near 'index(1) FROM t' at line 1
> exec ^
: string
| SELECT index(SELECT 1)
< error
| 1064 42000 …near 'index(SELECT 1)' at line 1
> exec ^
: string
| SELECT index(YEAR FROM 1)
< error
| 1064 42000 …near 'index(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT index(DAY, 1, 2)
< error
| 1064 42000 …near 'index(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT index(a)
< error
| 1064 42000 …near 'index(a)' at line 1
> exec ^
: string
| SELECT index(a, b)
< error
| 1064 42000 …near 'index(a, b)' at line 1
> exec ^
: string
| SELECT index(a) OVER ()
< error
| 1064 42000 …near 'index(a) OVER ()' at line 1
> exec ^
: string
| SELECT indexes()
< error
| 1630 42000 FUNCTION test.indexes does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT indexes()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT infile()
< error
| 1064 42000 …near 'infile()' at line 1
> exec ^
: string
| SELECT infile(1)
< error
| 1064 42000 …near 'infile(1)' at line 1
> exec ^
: string
| SELECT infile(1,2)
< error
| 1064 42000 …near 'infile(1,2)' at line 1
> exec ^
: string
| SELECT infile(1,2,3)
< error
| 1064 42000 …near 'infile(1,2,3)' at line 1
> exec ^
: string
| SELECT infile(1,2,3,4)
< error
| 1064 42000 …near 'infile(1,2,3,4)' at line 1
> exec ^
: string
| SELECT infile(*)
< error
| 1064 42000 …near 'infile(*)' at line 1
> exec ^
: string
| SELECT infile('a')
< error
| 1064 42000 …near 'infile('a')' at line 1
> exec ^
: string
| SELECT infile(a USING utf8mb4)
< error
| 1064 42000 …near 'infile(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT infile(1 AS SIGNED)
< error
| 1064 42000 …near 'infile(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT infile(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'infile(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT infile(1 FROM 2)
< error
| 1064 42000 …near 'infile(1 FROM 2)' at line 1
> exec ^
: string
| SELECT infile(1 IN 2)
< error
| 1064 42000 …near 'infile(1 IN 2)' at line 1
> exec ^
: string
| SELECT infile(DISTINCT 1)
< error
| 1064 42000 …near 'infile(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT infile() OVER ()
< error
| 1064 42000 …near 'infile() OVER ()' at line 1
> exec ^
: string
| SELECT infile(1) OVER ()
< error
| 1064 42000 …near 'infile(1) OVER ()' at line 1
> exec ^
: string
| SELECT infile(1,2) OVER ()
< error
| 1064 42000 …near 'infile(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT infile(1,2,3) OVER ()
< error
| 1064 42000 …near 'infile(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT infile(1) FROM t
< error
| 1064 42000 …near 'infile(1) FROM t' at line 1
> exec ^
: string
| SELECT infile(SELECT 1)
< error
| 1064 42000 …near 'infile(SELECT 1)' at line 1
> exec ^
: string
| SELECT infile(YEAR FROM 1)
< error
| 1064 42000 …near 'infile(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT infile(DAY, 1, 2)
< error
| 1064 42000 …near 'infile(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT infile(a)
< error
| 1064 42000 …near 'infile(a)' at line 1
> exec ^
: string
| SELECT infile(a, b)
< error
| 1064 42000 …near 'infile(a, b)' at line 1
> exec ^
: string
| SELECT infile(a) OVER ()
< error
| 1064 42000 …near 'infile(a) OVER ()' at line 1
> exec ^
: string
| SELECT initial()
< error
| 1305 42000 FUNCTION test.initial does not exist
> exec ^
: string
| SELECT initial()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT initial_size()
< error
| 1630 42000 FUNCTION test.initial_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT initial_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT initiate()
< error
| 1305 42000 FUNCTION test.initiate does not exist
> exec ^
: string
| SELECT initiate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT inner()
< error
| 1064 42000 …near 'inner()' at line 1
> exec ^
: string
| SELECT inner(1)
< error
| 1064 42000 …near 'inner(1)' at line 1
> exec ^
: string
| SELECT inner(1,2)
< error
| 1064 42000 …near 'inner(1,2)' at line 1
> exec ^
: string
| SELECT inner(1,2,3)
< error
| 1064 42000 …near 'inner(1,2,3)' at line 1
> exec ^
: string
| SELECT inner(1,2,3,4)
< error
| 1064 42000 …near 'inner(1,2,3,4)' at line 1
> exec ^
: string
| SELECT inner(*)
< error
| 1064 42000 …near 'inner(*)' at line 1
> exec ^
: string
| SELECT inner('a')
< error
| 1064 42000 …near 'inner('a')' at line 1
> exec ^
: string
| SELECT inner(a USING utf8mb4)
< error
| 1064 42000 …near 'inner(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT inner(1 AS SIGNED)
< error
| 1064 42000 …near 'inner(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT inner(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'inner(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT inner(1 FROM 2)
< error
| 1064 42000 …near 'inner(1 FROM 2)' at line 1
> exec ^
: string
| SELECT inner(1 IN 2)
< error
| 1064 42000 …near 'inner(1 IN 2)' at line 1
> exec ^
: string
| SELECT inner(DISTINCT 1)
< error
| 1064 42000 …near 'inner(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT inner() OVER ()
< error
| 1064 42000 …near 'inner() OVER ()' at line 1
> exec ^
: string
| SELECT inner(1) OVER ()
< error
| 1064 42000 …near 'inner(1) OVER ()' at line 1
> exec ^
: string
| SELECT inner(1,2) OVER ()
< error
| 1064 42000 …near 'inner(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT inner(1,2,3) OVER ()
< error
| 1064 42000 …near 'inner(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT inner(1) FROM t
< error
| 1064 42000 …near 'inner(1) FROM t' at line 1
> exec ^
: string
| SELECT inner(SELECT 1)
< error
| 1064 42000 …near 'inner(SELECT 1)' at line 1
> exec ^
: string
| SELECT inner(YEAR FROM 1)
< error
| 1064 42000 …near 'inner(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT inner(DAY, 1, 2)
< error
| 1064 42000 …near 'inner(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT inner(a)
< error
| 1064 42000 …near 'inner(a)' at line 1
> exec ^
: string
| SELECT inner(a, b)
< error
| 1064 42000 …near 'inner(a, b)' at line 1
> exec ^
: string
| SELECT inner(a) OVER ()
< error
| 1064 42000 …near 'inner(a) OVER ()' at line 1
> exec ^
: string
| SELECT inout()
< error
| 1064 42000 …near 'inout()' at line 1
> exec ^
: string
| SELECT inout(1)
< error
| 1064 42000 …near 'inout(1)' at line 1
> exec ^
: string
| SELECT inout(1,2)
< error
| 1064 42000 …near 'inout(1,2)' at line 1
> exec ^
: string
| SELECT inout(1,2,3)
< error
| 1064 42000 …near 'inout(1,2,3)' at line 1
> exec ^
: string
| SELECT inout(1,2,3,4)
< error
| 1064 42000 …near 'inout(1,2,3,4)' at line 1
> exec ^
: string
| SELECT inout(*)
< error
| 1064 42000 …near 'inout(*)' at line 1
> exec ^
: string
| SELECT inout('a')
< error
| 1064 42000 …near 'inout('a')' at line 1
> exec ^
: string
| SELECT inout(a USING utf8mb4)
< error
| 1064 42000 …near 'inout(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT inout(1 AS SIGNED)
< error
| 1064 42000 …near 'inout(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT inout(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'inout(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT inout(1 FROM 2)
< error
| 1064 42000 …near 'inout(1 FROM 2)' at line 1
> exec ^
: string
| SELECT inout(1 IN 2)
< error
| 1064 42000 …near 'inout(1 IN 2)' at line 1
> exec ^
: string
| SELECT inout(DISTINCT 1)
< error
| 1064 42000 …near 'inout(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT inout() OVER ()
< error
| 1064 42000 …near 'inout() OVER ()' at line 1
> exec ^
: string
| SELECT inout(1) OVER ()
< error
| 1064 42000 …near 'inout(1) OVER ()' at line 1
> exec ^
: string
| SELECT inout(1,2) OVER ()
< error
| 1064 42000 …near 'inout(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT inout(1,2,3) OVER ()
< error
| 1064 42000 …near 'inout(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT inout(1) FROM t
< error
| 1064 42000 …near 'inout(1) FROM t' at line 1
> exec ^
: string
| SELECT inout(SELECT 1)
< error
| 1064 42000 …near 'inout(SELECT 1)' at line 1
> exec ^
: string
| SELECT inout(YEAR FROM 1)
< error
| 1064 42000 …near 'inout(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT inout(DAY, 1, 2)
< error
| 1064 42000 …near 'inout(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT inout(a)
< error
| 1064 42000 …near 'inout(a)' at line 1
> exec ^
: string
| SELECT inout(a, b)
< error
| 1064 42000 …near 'inout(a, b)' at line 1
> exec ^
: string
| SELECT inout(a) OVER ()
< error
| 1064 42000 …near 'inout(a) OVER ()' at line 1
> exec ^
: string
| SELECT insensitive()
< error
| 1064 42000 …near 'insensitive()' at line 1
> exec ^
: string
| SELECT insensitive(1)
< error
| 1064 42000 …near 'insensitive(1)' at line 1
> exec ^
: string
| SELECT insensitive(1,2)
< error
| 1064 42000 …near 'insensitive(1,2)' at line 1
> exec ^
: string
| SELECT insensitive(1,2,3)
< error
| 1064 42000 …near 'insensitive(1,2,3)' at line 1
> exec ^
: string
| SELECT insensitive(1,2,3,4)
< error
| 1064 42000 …near 'insensitive(1,2,3,4)' at line 1
> exec ^
: string
| SELECT insensitive(*)
< error
| 1064 42000 …near 'insensitive(*)' at line 1
> exec ^
: string
| SELECT insensitive('a')
< error
| 1064 42000 …near 'insensitive('a')' at line 1
> exec ^
: string
| SELECT insensitive(a USING utf8mb4)
< error
| 1064 42000 …near 'insensitive(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT insensitive(1 AS SIGNED)
< error
| 1064 42000 …near 'insensitive(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT insensitive(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'insensitive(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT insensitive(1 FROM 2)
< error
| 1064 42000 …near 'insensitive(1 FROM 2)' at line 1
> exec ^
: string
| SELECT insensitive(1 IN 2)
< error
| 1064 42000 …near 'insensitive(1 IN 2)' at line 1
> exec ^
: string
| SELECT insensitive(DISTINCT 1)
< error
| 1064 42000 …near 'insensitive(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT insensitive() OVER ()
< error
| 1064 42000 …near 'insensitive() OVER ()' at line 1
> exec ^
: string
| SELECT insensitive(1) OVER ()
< error
| 1064 42000 …near 'insensitive(1) OVER ()' at line 1
> exec ^
: string
| SELECT insensitive(1,2) OVER ()
< error
| 1064 42000 …near 'insensitive(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT insensitive(1,2,3) OVER ()
< error
| 1064 42000 …near 'insensitive(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT insensitive(1) FROM t
< error
| 1064 42000 …near 'insensitive(1) FROM t' at line 1
> exec ^
: string
| SELECT insensitive(SELECT 1)
< error
| 1064 42000 …near 'insensitive(SELECT 1)' at line 1
> exec ^
: string
| SELECT insensitive(YEAR FROM 1)
< error
| 1064 42000 …near 'insensitive(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT insensitive(DAY, 1, 2)
< error
| 1064 42000 …near 'insensitive(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT insensitive(a)
< error
| 1064 42000 …near 'insensitive(a)' at line 1
> exec ^
: string
| SELECT insensitive(a, b)
< error
| 1064 42000 …near 'insensitive(a, b)' at line 1
> exec ^
: string
| SELECT insensitive(a) OVER ()
< error
| 1064 42000 …near 'insensitive(a) OVER ()' at line 1
> exec ^
: string
| SELECT insert()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT insert(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT insert(1,2)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT insert(1,2,3)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT insert(1,2,3,4)
< ok 0 0
> exec ^
: string
| SELECT insert_method()
< error
| 1630 42000 FUNCTION test.insert_method does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT insert_method()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT install()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT install(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT install(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT install(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT install(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT install(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT install('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT install(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT install(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT install(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT install(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT install(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT install(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT install() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT install(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT install(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT install(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT install(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT install(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT install(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT install(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT install(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT install(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT install(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT instance()
< error
| 1305 42000 FUNCTION test.instance does not exist
> exec ^
: string
| SELECT instance()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT int()
< error
| 1064 42000 …near 'int()' at line 1
> exec ^
: string
| SELECT int(1)
< error
| 1064 42000 …near 'int(1)' at line 1
> exec ^
: string
| SELECT int(1,2)
< error
| 1064 42000 …near 'int(1,2)' at line 1
> exec ^
: string
| SELECT int(1,2,3)
< error
| 1064 42000 …near 'int(1,2,3)' at line 1
> exec ^
: string
| SELECT int(1,2,3,4)
< error
| 1064 42000 …near 'int(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int(*)
< error
| 1064 42000 …near 'int(*)' at line 1
> exec ^
: string
| SELECT int('a')
< error
| 1064 42000 …near 'int('a')' at line 1
> exec ^
: string
| SELECT int(a USING utf8mb4)
< error
| 1064 42000 …near 'int(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int(1 AS SIGNED)
< error
| 1064 42000 …near 'int(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int(1 FROM 2)
< error
| 1064 42000 …near 'int(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int(1 IN 2)
< error
| 1064 42000 …near 'int(1 IN 2)' at line 1
> exec ^
: string
| SELECT int(DISTINCT 1)
< error
| 1064 42000 …near 'int(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int() OVER ()
< error
| 1064 42000 …near 'int() OVER ()' at line 1
> exec ^
: string
| SELECT int(1) OVER ()
< error
| 1064 42000 …near 'int(1) OVER ()' at line 1
> exec ^
: string
| SELECT int(1,2) OVER ()
< error
| 1064 42000 …near 'int(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int(1,2,3) OVER ()
< error
| 1064 42000 …near 'int(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int(1) FROM t
< error
| 1064 42000 …near 'int(1) FROM t' at line 1
> exec ^
: string
| SELECT int(SELECT 1)
< error
| 1064 42000 …near 'int(SELECT 1)' at line 1
> exec ^
: string
| SELECT int(YEAR FROM 1)
< error
| 1064 42000 …near 'int(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int(DAY, 1, 2)
< error
| 1064 42000 …near 'int(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int(a)
< error
| 1064 42000 …near 'int(a)' at line 1
> exec ^
: string
| SELECT int(a, b)
< error
| 1064 42000 …near 'int(a, b)' at line 1
> exec ^
: string
| SELECT int(a) OVER ()
< error
| 1064 42000 …near 'int(a) OVER ()' at line 1
> exec ^
: string
| SELECT int1()
< error
| 1064 42000 …near 'int1()' at line 1
> exec ^
: string
| SELECT int1(1)
< error
| 1064 42000 …near 'int1(1)' at line 1
> exec ^
: string
| SELECT int1(1,2)
< error
| 1064 42000 …near 'int1(1,2)' at line 1
> exec ^
: string
| SELECT int1(1,2,3)
< error
| 1064 42000 …near 'int1(1,2,3)' at line 1
> exec ^
: string
| SELECT int1(1,2,3,4)
< error
| 1064 42000 …near 'int1(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int1(*)
< error
| 1064 42000 …near 'int1(*)' at line 1
> exec ^
: string
| SELECT int1('a')
< error
| 1064 42000 …near 'int1('a')' at line 1
> exec ^
: string
| SELECT int1(a USING utf8mb4)
< error
| 1064 42000 …near 'int1(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int1(1 AS SIGNED)
< error
| 1064 42000 …near 'int1(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int1(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int1(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int1(1 FROM 2)
< error
| 1064 42000 …near 'int1(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int1(1 IN 2)
< error
| 1064 42000 …near 'int1(1 IN 2)' at line 1
> exec ^
: string
| SELECT int1(DISTINCT 1)
< error
| 1064 42000 …near 'int1(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int1() OVER ()
< error
| 1064 42000 …near 'int1() OVER ()' at line 1
> exec ^
: string
| SELECT int1(1) OVER ()
< error
| 1064 42000 …near 'int1(1) OVER ()' at line 1
> exec ^
: string
| SELECT int1(1,2) OVER ()
< error
| 1064 42000 …near 'int1(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int1(1,2,3) OVER ()
< error
| 1064 42000 …near 'int1(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int1(1) FROM t
< error
| 1064 42000 …near 'int1(1) FROM t' at line 1
> exec ^
: string
| SELECT int1(SELECT 1)
< error
| 1064 42000 …near 'int1(SELECT 1)' at line 1
> exec ^
: string
| SELECT int1(YEAR FROM 1)
< error
| 1064 42000 …near 'int1(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int1(DAY, 1, 2)
< error
| 1064 42000 …near 'int1(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int1(a)
< error
| 1064 42000 …near 'int1(a)' at line 1
> exec ^
: string
| SELECT int1(a, b)
< error
| 1064 42000 …near 'int1(a, b)' at line 1
> exec ^
: string
| SELECT int1(a) OVER ()
< error
| 1064 42000 …near 'int1(a) OVER ()' at line 1
> exec ^
: string
| SELECT int2()
< error
| 1064 42000 …near 'int2()' at line 1
> exec ^
: string
| SELECT int2(1)
< error
| 1064 42000 …near 'int2(1)' at line 1
> exec ^
: string
| SELECT int2(1,2)
< error
| 1064 42000 …near 'int2(1,2)' at line 1
> exec ^
: string
| SELECT int2(1,2,3)
< error
| 1064 42000 …near 'int2(1,2,3)' at line 1
> exec ^
: string
| SELECT int2(1,2,3,4)
< error
| 1064 42000 …near 'int2(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int2(*)
< error
| 1064 42000 …near 'int2(*)' at line 1
> exec ^
: string
| SELECT int2('a')
< error
| 1064 42000 …near 'int2('a')' at line 1
> exec ^
: string
| SELECT int2(a USING utf8mb4)
< error
| 1064 42000 …near 'int2(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int2(1 AS SIGNED)
< error
| 1064 42000 …near 'int2(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int2(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int2(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int2(1 FROM 2)
< error
| 1064 42000 …near 'int2(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int2(1 IN 2)
< error
| 1064 42000 …near 'int2(1 IN 2)' at line 1
> exec ^
: string
| SELECT int2(DISTINCT 1)
< error
| 1064 42000 …near 'int2(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int2() OVER ()
< error
| 1064 42000 …near 'int2() OVER ()' at line 1
> exec ^
: string
| SELECT int2(1) OVER ()
< error
| 1064 42000 …near 'int2(1) OVER ()' at line 1
> exec ^
: string
| SELECT int2(1,2) OVER ()
< error
| 1064 42000 …near 'int2(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int2(1,2,3) OVER ()
< error
| 1064 42000 …near 'int2(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int2(1) FROM t
< error
| 1064 42000 …near 'int2(1) FROM t' at line 1
> exec ^
: string
| SELECT int2(SELECT 1)
< error
| 1064 42000 …near 'int2(SELECT 1)' at line 1
> exec ^
: string
| SELECT int2(YEAR FROM 1)
< error
| 1064 42000 …near 'int2(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int2(DAY, 1, 2)
< error
| 1064 42000 …near 'int2(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int2(a)
< error
| 1064 42000 …near 'int2(a)' at line 1
> exec ^
: string
| SELECT int2(a, b)
< error
| 1064 42000 …near 'int2(a, b)' at line 1
> exec ^
: string
| SELECT int2(a) OVER ()
< error
| 1064 42000 …near 'int2(a) OVER ()' at line 1
> exec ^
: string
| SELECT int3()
< error
| 1064 42000 …near 'int3()' at line 1
> exec ^
: string
| SELECT int3(1)
< error
| 1064 42000 …near 'int3(1)' at line 1
> exec ^
: string
| SELECT int3(1,2)
< error
| 1064 42000 …near 'int3(1,2)' at line 1
> exec ^
: string
| SELECT int3(1,2,3)
< error
| 1064 42000 …near 'int3(1,2,3)' at line 1
> exec ^
: string
| SELECT int3(1,2,3,4)
< error
| 1064 42000 …near 'int3(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int3(*)
< error
| 1064 42000 …near 'int3(*)' at line 1
> exec ^
: string
| SELECT int3('a')
< error
| 1064 42000 …near 'int3('a')' at line 1
> exec ^
: string
| SELECT int3(a USING utf8mb4)
< error
| 1064 42000 …near 'int3(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int3(1 AS SIGNED)
< error
| 1064 42000 …near 'int3(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int3(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int3(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int3(1 FROM 2)
< error
| 1064 42000 …near 'int3(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int3(1 IN 2)
< error
| 1064 42000 …near 'int3(1 IN 2)' at line 1
> exec ^
: string
| SELECT int3(DISTINCT 1)
< error
| 1064 42000 …near 'int3(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int3() OVER ()
< error
| 1064 42000 …near 'int3() OVER ()' at line 1
> exec ^
: string
| SELECT int3(1) OVER ()
< error
| 1064 42000 …near 'int3(1) OVER ()' at line 1
> exec ^
: string
| SELECT int3(1,2) OVER ()
< error
| 1064 42000 …near 'int3(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int3(1,2,3) OVER ()
< error
| 1064 42000 …near 'int3(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int3(1) FROM t
< error
| 1064 42000 …near 'int3(1) FROM t' at line 1
> exec ^
: string
| SELECT int3(SELECT 1)
< error
| 1064 42000 …near 'int3(SELECT 1)' at line 1
> exec ^
: string
| SELECT int3(YEAR FROM 1)
< error
| 1064 42000 …near 'int3(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int3(DAY, 1, 2)
< error
| 1064 42000 …near 'int3(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int3(a)
< error
| 1064 42000 …near 'int3(a)' at line 1
> exec ^
: string
| SELECT int3(a, b)
< error
| 1064 42000 …near 'int3(a, b)' at line 1
> exec ^
: string
| SELECT int3(a) OVER ()
< error
| 1064 42000 …near 'int3(a) OVER ()' at line 1
> exec ^
: string
| SELECT int4()
< error
| 1064 42000 …near 'int4()' at line 1
> exec ^
: string
| SELECT int4(1)
< error
| 1064 42000 …near 'int4(1)' at line 1
> exec ^
: string
| SELECT int4(1,2)
< error
| 1064 42000 …near 'int4(1,2)' at line 1
> exec ^
: string
| SELECT int4(1,2,3)
< error
| 1064 42000 …near 'int4(1,2,3)' at line 1
> exec ^
: string
| SELECT int4(1,2,3,4)
< error
| 1064 42000 …near 'int4(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int4(*)
< error
| 1064 42000 …near 'int4(*)' at line 1
> exec ^
: string
| SELECT int4('a')
< error
| 1064 42000 …near 'int4('a')' at line 1
> exec ^
: string
| SELECT int4(a USING utf8mb4)
< error
| 1064 42000 …near 'int4(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int4(1 AS SIGNED)
< error
| 1064 42000 …near 'int4(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int4(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int4(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int4(1 FROM 2)
< error
| 1064 42000 …near 'int4(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int4(1 IN 2)
< error
| 1064 42000 …near 'int4(1 IN 2)' at line 1
> exec ^
: string
| SELECT int4(DISTINCT 1)
< error
| 1064 42000 …near 'int4(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int4() OVER ()
< error
| 1064 42000 …near 'int4() OVER ()' at line 1
> exec ^
: string
| SELECT int4(1) OVER ()
< error
| 1064 42000 …near 'int4(1) OVER ()' at line 1
> exec ^
: string
| SELECT int4(1,2) OVER ()
< error
| 1064 42000 …near 'int4(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int4(1,2,3) OVER ()
< error
| 1064 42000 …near 'int4(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int4(1) FROM t
< error
| 1064 42000 …near 'int4(1) FROM t' at line 1
> exec ^
: string
| SELECT int4(SELECT 1)
< error
| 1064 42000 …near 'int4(SELECT 1)' at line 1
> exec ^
: string
| SELECT int4(YEAR FROM 1)
< error
| 1064 42000 …near 'int4(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int4(DAY, 1, 2)
< error
| 1064 42000 …near 'int4(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int4(a)
< error
| 1064 42000 …near 'int4(a)' at line 1
> exec ^
: string
| SELECT int4(a, b)
< error
| 1064 42000 …near 'int4(a, b)' at line 1
> exec ^
: string
| SELECT int4(a) OVER ()
< error
| 1064 42000 …near 'int4(a) OVER ()' at line 1
> exec ^
: string
| SELECT int8()
< error
| 1064 42000 …near 'int8()' at line 1
> exec ^
: string
| SELECT int8(1)
< error
| 1064 42000 …near 'int8(1)' at line 1
> exec ^
: string
| SELECT int8(1,2)
< error
| 1064 42000 …near 'int8(1,2)' at line 1
> exec ^
: string
| SELECT int8(1,2,3)
< error
| 1064 42000 …near 'int8(1,2,3)' at line 1
> exec ^
: string
| SELECT int8(1,2,3,4)
< error
| 1064 42000 …near 'int8(1,2,3,4)' at line 1
> exec ^
: string
| SELECT int8(*)
< error
| 1064 42000 …near 'int8(*)' at line 1
> exec ^
: string
| SELECT int8('a')
< error
| 1064 42000 …near 'int8('a')' at line 1
> exec ^
: string
| SELECT int8(a USING utf8mb4)
< error
| 1064 42000 …near 'int8(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT int8(1 AS SIGNED)
< error
| 1064 42000 …near 'int8(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT int8(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'int8(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT int8(1 FROM 2)
< error
| 1064 42000 …near 'int8(1 FROM 2)' at line 1
> exec ^
: string
| SELECT int8(1 IN 2)
< error
| 1064 42000 …near 'int8(1 IN 2)' at line 1
> exec ^
: string
| SELECT int8(DISTINCT 1)
< error
| 1064 42000 …near 'int8(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT int8() OVER ()
< error
| 1064 42000 …near 'int8() OVER ()' at line 1
> exec ^
: string
| SELECT int8(1) OVER ()
< error
| 1064 42000 …near 'int8(1) OVER ()' at line 1
> exec ^
: string
| SELECT int8(1,2) OVER ()
< error
| 1064 42000 …near 'int8(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT int8(1,2,3) OVER ()
< error
| 1064 42000 …near 'int8(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT int8(1) FROM t
< error
| 1064 42000 …near 'int8(1) FROM t' at line 1
> exec ^
: string
| SELECT int8(SELECT 1)
< error
| 1064 42000 …near 'int8(SELECT 1)' at line 1
> exec ^
: string
| SELECT int8(YEAR FROM 1)
< error
| 1064 42000 …near 'int8(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT int8(DAY, 1, 2)
< error
| 1064 42000 …near 'int8(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT int8(a)
< error
| 1064 42000 …near 'int8(a)' at line 1
> exec ^
: string
| SELECT int8(a, b)
< error
| 1064 42000 …near 'int8(a, b)' at line 1
> exec ^
: string
| SELECT int8(a) OVER ()
< error
| 1064 42000 …near 'int8(a) OVER ()' at line 1
> exec ^
: string
| SELECT integer()
< error
| 1064 42000 …near 'integer()' at line 1
> exec ^
: string
| SELECT integer(1)
< error
| 1064 42000 …near 'integer(1)' at line 1
> exec ^
: string
| SELECT integer(1,2)
< error
| 1064 42000 …near 'integer(1,2)' at line 1
> exec ^
: string
| SELECT integer(1,2,3)
< error
| 1064 42000 …near 'integer(1,2,3)' at line 1
> exec ^
: string
| SELECT integer(1,2,3,4)
< error
| 1064 42000 …near 'integer(1,2,3,4)' at line 1
> exec ^
: string
| SELECT integer(*)
< error
| 1064 42000 …near 'integer(*)' at line 1
> exec ^
: string
| SELECT integer('a')
< error
| 1064 42000 …near 'integer('a')' at line 1
> exec ^
: string
| SELECT integer(a USING utf8mb4)
< error
| 1064 42000 …near 'integer(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT integer(1 AS SIGNED)
< error
| 1064 42000 …near 'integer(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT integer(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'integer(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT integer(1 FROM 2)
< error
| 1064 42000 …near 'integer(1 FROM 2)' at line 1
> exec ^
: string
| SELECT integer(1 IN 2)
< error
| 1064 42000 …near 'integer(1 IN 2)' at line 1
> exec ^
: string
| SELECT integer(DISTINCT 1)
< error
| 1064 42000 …near 'integer(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT integer() OVER ()
< error
| 1064 42000 …near 'integer() OVER ()' at line 1
> exec ^
: string
| SELECT integer(1) OVER ()
< error
| 1064 42000 …near 'integer(1) OVER ()' at line 1
> exec ^
: string
| SELECT integer(1,2) OVER ()
< error
| 1064 42000 …near 'integer(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT integer(1,2,3) OVER ()
< error
| 1064 42000 …near 'integer(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT integer(1) FROM t
< error
| 1064 42000 …near 'integer(1) FROM t' at line 1
> exec ^
: string
| SELECT integer(SELECT 1)
< error
| 1064 42000 …near 'integer(SELECT 1)' at line 1
> exec ^
: string
| SELECT integer(YEAR FROM 1)
< error
| 1064 42000 …near 'integer(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT integer(DAY, 1, 2)
< error
| 1064 42000 …near 'integer(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT integer(a)
< error
| 1064 42000 …near 'integer(a)' at line 1
> exec ^
: string
| SELECT integer(a, b)
< error
| 1064 42000 …near 'integer(a, b)' at line 1
> exec ^
: string
| SELECT integer(a) OVER ()
< error
| 1064 42000 …near 'integer(a) OVER ()' at line 1
> exec ^
: string
| SELECT intersect()
< error
| 1064 42000 …near 'intersect()' at line 1
> exec ^
: string
| SELECT intersect(1)
< error
| 1064 42000 …near 'intersect(1)' at line 1
> exec ^
: string
| SELECT intersect(1,2)
< error
| 1064 42000 …near 'intersect(1,2)' at line 1
> exec ^
: string
| SELECT intersect(1,2,3)
< error
| 1064 42000 …near 'intersect(1,2,3)' at line 1
> exec ^
: string
| SELECT intersect(1,2,3,4)
< error
| 1064 42000 …near 'intersect(1,2,3,4)' at line 1
> exec ^
: string
| SELECT intersect(*)
< error
| 1064 42000 …near 'intersect(*)' at line 1
> exec ^
: string
| SELECT intersect('a')
< error
| 1064 42000 …near 'intersect('a')' at line 1
> exec ^
: string
| SELECT intersect(a USING utf8mb4)
< error
| 1064 42000 …near 'intersect(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT intersect(1 AS SIGNED)
< error
| 1064 42000 …near 'intersect(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT intersect(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'intersect(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT intersect(1 FROM 2)
< error
| 1064 42000 …near 'intersect(1 FROM 2)' at line 1
> exec ^
: string
| SELECT intersect(1 IN 2)
< error
| 1064 42000 …near 'intersect(1 IN 2)' at line 1
> exec ^
: string
| SELECT intersect(DISTINCT 1)
< error
| 1064 42000 …near 'intersect(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT intersect() OVER ()
< error
| 1064 42000 …near 'intersect() OVER ()' at line 1
> exec ^
: string
| SELECT intersect(1) OVER ()
< error
| 1064 42000 …near 'intersect(1) OVER ()' at line 1
> exec ^
: string
| SELECT intersect(1,2) OVER ()
< error
| 1064 42000 …near 'intersect(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT intersect(1,2,3) OVER ()
< error
| 1064 42000 …near 'intersect(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT intersect(1) FROM t
< error
| 1064 42000 …near 'intersect(1) FROM t' at line 1
> exec ^
: string
| SELECT intersect(SELECT 1)
< error
| 1064 42000 …near 'intersect(SELECT 1)' at line 1
> exec ^
: string
| SELECT intersect(YEAR FROM 1)
< error
| 1064 42000 …near 'intersect(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT intersect(DAY, 1, 2)
< error
| 1064 42000 …near 'intersect(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT intersect(a)
< error
| 1064 42000 …near 'intersect(a)' at line 1
> exec ^
: string
| SELECT intersect(a, b)
< error
| 1064 42000 …near 'intersect(a, b)' at line 1
> exec ^
: string
| SELECT intersect(a) OVER ()
< error
| 1064 42000 …near 'intersect(a) OVER ()' at line 1
> exec ^
: string
| SELECT interval()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT interval(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT interval(1,2)
< ok 0 0
> exec ^
: string
| SELECT into()
< error
| 1064 42000 …near 'into()' at line 1
> exec ^
: string
| SELECT into(1)
< error
| 1064 42000 …near 'into(1)' at line 1
> exec ^
: string
| SELECT into(1,2)
< error
| 1064 42000 …near 'into(1,2)' at line 1
> exec ^
: string
| SELECT into(1,2,3)
< error
| 1064 42000 …near 'into(1,2,3)' at line 1
> exec ^
: string
| SELECT into(1,2,3,4)
< error
| 1064 42000 …near 'into(1,2,3,4)' at line 1
> exec ^
: string
| SELECT into(*)
< error
| 1064 42000 …near 'into(*)' at line 1
> exec ^
: string
| SELECT into('a')
< error
| 1064 42000 …near 'into('a')' at line 1
> exec ^
: string
| SELECT into(a USING utf8mb4)
< error
| 1064 42000 …near 'into(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT into(1 AS SIGNED)
< error
| 1064 42000 …near 'into(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT into(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'into(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT into(1 FROM 2)
< error
| 1064 42000 …near 'into(1 FROM 2)' at line 1
> exec ^
: string
| SELECT into(1 IN 2)
< error
| 1064 42000 …near 'into(1 IN 2)' at line 1
> exec ^
: string
| SELECT into(DISTINCT 1)
< error
| 1064 42000 …near 'into(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT into() OVER ()
< error
| 1064 42000 …near 'into() OVER ()' at line 1
> exec ^
: string
| SELECT into(1) OVER ()
< error
| 1064 42000 …near 'into(1) OVER ()' at line 1
> exec ^
: string
| SELECT into(1,2) OVER ()
< error
| 1064 42000 …near 'into(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT into(1,2,3) OVER ()
< error
| 1064 42000 …near 'into(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT into(1) FROM t
< error
| 1064 42000 …near 'into(1) FROM t' at line 1
> exec ^
: string
| SELECT into(SELECT 1)
< error
| 1064 42000 …near 'into(SELECT 1)' at line 1
> exec ^
: string
| SELECT into(YEAR FROM 1)
< error
| 1064 42000 …near 'into(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT into(DAY, 1, 2)
< error
| 1064 42000 …near 'into(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT into(a)
< error
| 1064 42000 …near 'into(a)' at line 1
> exec ^
: string
| SELECT into(a, b)
< error
| 1064 42000 …near 'into(a, b)' at line 1
> exec ^
: string
| SELECT into(a) OVER ()
< error
| 1064 42000 …near 'into(a) OVER ()' at line 1
> exec ^
: string
| SELECT invisible()
< error
| 1630 42000 FUNCTION test.invisible does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT invisible()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT invoker()
< error
| 1630 42000 FUNCTION test.invoker does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT invoker()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT io()
< error
| 1630 42000 FUNCTION test.io does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT io()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT io_after_gtids()
< error
| 1305 42000 FUNCTION test.io_after_gtids does not exist
> exec ^
: string
| SELECT io_after_gtids()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT io_before_gtids()
< error
| 1305 42000 FUNCTION test.io_before_gtids does not exist
> exec ^
: string
| SELECT io_before_gtids()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT io_thread()
< error
| 1630 42000 FUNCTION test.io_thread does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT io_thread()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ipc()
< error
| 1630 42000 FUNCTION test.ipc does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT ipc()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT is()
< error
| 1064 42000 …near 'is()' at line 1
> exec ^
: string
| SELECT is(1)
< error
| 1064 42000 …near 'is(1)' at line 1
> exec ^
: string
| SELECT is(1,2)
< error
| 1064 42000 …near 'is(1,2)' at line 1
> exec ^
: string
| SELECT is(1,2,3)
< error
| 1064 42000 …near 'is(1,2,3)' at line 1
> exec ^
: string
| SELECT is(1,2,3,4)
< error
| 1064 42000 …near 'is(1,2,3,4)' at line 1
> exec ^
: string
| SELECT is(*)
< error
| 1064 42000 …near 'is(*)' at line 1
> exec ^
: string
| SELECT is('a')
< error
| 1064 42000 …near 'is('a')' at line 1
> exec ^
: string
| SELECT is(a USING utf8mb4)
< error
| 1064 42000 …near 'is(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT is(1 AS SIGNED)
< error
| 1064 42000 …near 'is(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT is(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'is(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT is(1 FROM 2)
< error
| 1064 42000 …near 'is(1 FROM 2)' at line 1
> exec ^
: string
| SELECT is(1 IN 2)
< error
| 1064 42000 …near 'is(1 IN 2)' at line 1
> exec ^
: string
| SELECT is(DISTINCT 1)
< error
| 1064 42000 …near 'is(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT is() OVER ()
< error
| 1064 42000 …near 'is() OVER ()' at line 1
> exec ^
: string
| SELECT is(1) OVER ()
< error
| 1064 42000 …near 'is(1) OVER ()' at line 1
> exec ^
: string
| SELECT is(1,2) OVER ()
< error
| 1064 42000 …near 'is(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT is(1,2,3) OVER ()
< error
| 1064 42000 …near 'is(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT is(1) FROM t
< error
| 1064 42000 …near 'is(1) FROM t' at line 1
> exec ^
: string
| SELECT is(SELECT 1)
< error
| 1064 42000 …near 'is(SELECT 1)' at line 1
> exec ^
: string
| SELECT is(YEAR FROM 1)
< error
| 1064 42000 …near 'is(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT is(DAY, 1, 2)
< error
| 1064 42000 …near 'is(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT is(a)
< error
| 1064 42000 …near 'is(a)' at line 1
> exec ^
: string
| SELECT is(a, b)
< error
| 1064 42000 …near 'is(a, b)' at line 1
> exec ^
: string
| SELECT is(a) OVER ()
< error
| 1064 42000 …near 'is(a) OVER ()' at line 1
> exec ^
: string
| SELECT isolation()
< error
| 1630 42000 FUNCTION test.isolation does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT isolation()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT issuer()
< error
| 1630 42000 FUNCTION test.issuer does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT issuer()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT iterate()
< error
| 1064 42000 …near 'iterate()' at line 1
> exec ^
: string
| SELECT iterate(1)
< error
| 1064 42000 …near 'iterate(1)' at line 1
> exec ^
: string
| SELECT iterate(1,2)
< error
| 1064 42000 …near 'iterate(1,2)' at line 1
> exec ^
: string
| SELECT iterate(1,2,3)
< error
| 1064 42000 …near 'iterate(1,2,3)' at line 1
> exec ^
: string
| SELECT iterate(1,2,3,4)
< error
| 1064 42000 …near 'iterate(1,2,3,4)' at line 1
> exec ^
: string
| SELECT iterate(*)
< error
| 1064 42000 …near 'iterate(*)' at line 1
> exec ^
: string
| SELECT iterate('a')
< error
| 1064 42000 …near 'iterate('a')' at line 1
> exec ^
: string
| SELECT iterate(a USING utf8mb4)
< error
| 1064 42000 …near 'iterate(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT iterate(1 AS SIGNED)
< error
| 1064 42000 …near 'iterate(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT iterate(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'iterate(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT iterate(1 FROM 2)
< error
| 1064 42000 …near 'iterate(1 FROM 2)' at line 1
> exec ^
: string
| SELECT iterate(1 IN 2)
< error
| 1064 42000 …near 'iterate(1 IN 2)' at line 1
> exec ^
: string
| SELECT iterate(DISTINCT 1)
< error
| 1064 42000 …near 'iterate(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT iterate() OVER ()
< error
| 1064 42000 …near 'iterate() OVER ()' at line 1
> exec ^
: string
| SELECT iterate(1) OVER ()
< error
| 1064 42000 …near 'iterate(1) OVER ()' at line 1
> exec ^
: string
| SELECT iterate(1,2) OVER ()
< error
| 1064 42000 …near 'iterate(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT iterate(1,2,3) OVER ()
< error
| 1064 42000 …near 'iterate(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT iterate(1) FROM t
< error
| 1064 42000 …near 'iterate(1) FROM t' at line 1
> exec ^
: string
| SELECT iterate(SELECT 1)
< error
| 1064 42000 …near 'iterate(SELECT 1)' at line 1
> exec ^
: string
| SELECT iterate(YEAR FROM 1)
< error
| 1064 42000 …near 'iterate(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT iterate(DAY, 1, 2)
< error
| 1064 42000 …near 'iterate(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT iterate(a)
< error
| 1064 42000 …near 'iterate(a)' at line 1
> exec ^
: string
| SELECT iterate(a, b)
< error
| 1064 42000 …near 'iterate(a, b)' at line 1
> exec ^
: string
| SELECT iterate(a) OVER ()
< error
| 1064 42000 …near 'iterate(a) OVER ()' at line 1
> exec ^
: string
| SELECT join()
< error
| 1064 42000 …near 'join()' at line 1
> exec ^
: string
| SELECT join(1)
< error
| 1064 42000 …near 'join(1)' at line 1
> exec ^
: string
| SELECT join(1,2)
< error
| 1064 42000 …near 'join(1,2)' at line 1
> exec ^
: string
| SELECT join(1,2,3)
< error
| 1064 42000 …near 'join(1,2,3)' at line 1
> exec ^
: string
| SELECT join(1,2,3,4)
< error
| 1064 42000 …near 'join(1,2,3,4)' at line 1
> exec ^
: string
| SELECT join(*)
< error
| 1064 42000 …near 'join(*)' at line 1
> exec ^
: string
| SELECT join('a')
< error
| 1064 42000 …near 'join('a')' at line 1
> exec ^
: string
| SELECT join(a USING utf8mb4)
< error
| 1064 42000 …near 'join(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT join(1 AS SIGNED)
< error
| 1064 42000 …near 'join(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT join(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'join(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT join(1 FROM 2)
< error
| 1064 42000 …near 'join(1 FROM 2)' at line 1
> exec ^
: string
| SELECT join(1 IN 2)
< error
| 1064 42000 …near 'join(1 IN 2)' at line 1
> exec ^
: string
| SELECT join(DISTINCT 1)
< error
| 1064 42000 …near 'join(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT join() OVER ()
< error
| 1064 42000 …near 'join() OVER ()' at line 1
> exec ^
: string
| SELECT join(1) OVER ()
< error
| 1064 42000 …near 'join(1) OVER ()' at line 1
> exec ^
: string
| SELECT join(1,2) OVER ()
< error
| 1064 42000 …near 'join(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT join(1,2,3) OVER ()
< error
| 1064 42000 …near 'join(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT join(1) FROM t
< error
| 1064 42000 …near 'join(1) FROM t' at line 1
> exec ^
: string
| SELECT join(SELECT 1)
< error
| 1064 42000 …near 'join(SELECT 1)' at line 1
> exec ^
: string
| SELECT join(YEAR FROM 1)
< error
| 1064 42000 …near 'join(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT join(DAY, 1, 2)
< error
| 1064 42000 …near 'join(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT join(a)
< error
| 1064 42000 …near 'join(a)' at line 1
> exec ^
: string
| SELECT join(a, b)
< error
| 1064 42000 …near 'join(a, b)' at line 1
> exec ^
: string
| SELECT join(a) OVER ()
< error
| 1064 42000 …near 'join(a) OVER ()' at line 1
> exec ^
: string
| SELECT json()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT json(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT json(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT json(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT json(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT json(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT json('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT json(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT json(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT json(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT json(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT json(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT json(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT json() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT json(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT json(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT json(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT json(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT json(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT json(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT json(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT json(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT json(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT json(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT json_arrayagg()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT json_arrayagg(1)
< ok 0 0
> exec ^
: string
| SELECT json_duality_object()
< error
| 1305 42000 FUNCTION test.json_duality_object does not exist
> exec ^
: string
| SELECT json_duality_object()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT json_objectagg()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT json_objectagg(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT json_objectagg(1,2)
< ok 0 0
> exec ^
: string
| SELECT json_table()
< error
| 1630 42000 FUNCTION test.json_table does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT json_table()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT json_value()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'json_value'
> exec ^
: string
| SELECT json_value()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'json_value'
> exec ^
: string
| SELECT json_value(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'json_value'
> exec ^
: string
| SELECT json_value(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'json_value'
> exec ^
: string
| SELECT json_value(1,2)
< ok 0 0
> exec ^
: string
| SELECT key()
< error
| 1064 42000 …near 'key()' at line 1
> exec ^
: string
| SELECT key(1)
< error
| 1064 42000 …near 'key(1)' at line 1
> exec ^
: string
| SELECT key(1,2)
< error
| 1064 42000 …near 'key(1,2)' at line 1
> exec ^
: string
| SELECT key(1,2,3)
< error
| 1064 42000 …near 'key(1,2,3)' at line 1
> exec ^
: string
| SELECT key(1,2,3,4)
< error
| 1064 42000 …near 'key(1,2,3,4)' at line 1
> exec ^
: string
| SELECT key(*)
< error
| 1064 42000 …near 'key(*)' at line 1
> exec ^
: string
| SELECT key('a')
< error
| 1064 42000 …near 'key('a')' at line 1
> exec ^
: string
| SELECT key(a USING utf8mb4)
< error
| 1064 42000 …near 'key(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT key(1 AS SIGNED)
< error
| 1064 42000 …near 'key(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT key(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'key(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT key(1 FROM 2)
< error
| 1064 42000 …near 'key(1 FROM 2)' at line 1
> exec ^
: string
| SELECT key(1 IN 2)
< error
| 1064 42000 …near 'key(1 IN 2)' at line 1
> exec ^
: string
| SELECT key(DISTINCT 1)
< error
| 1064 42000 …near 'key(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT key() OVER ()
< error
| 1064 42000 …near 'key() OVER ()' at line 1
> exec ^
: string
| SELECT key(1) OVER ()
< error
| 1064 42000 …near 'key(1) OVER ()' at line 1
> exec ^
: string
| SELECT key(1,2) OVER ()
< error
| 1064 42000 …near 'key(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT key(1,2,3) OVER ()
< error
| 1064 42000 …near 'key(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT key(1) FROM t
< error
| 1064 42000 …near 'key(1) FROM t' at line 1
> exec ^
: string
| SELECT key(SELECT 1)
< error
| 1064 42000 …near 'key(SELECT 1)' at line 1
> exec ^
: string
| SELECT key(YEAR FROM 1)
< error
| 1064 42000 …near 'key(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT key(DAY, 1, 2)
< error
| 1064 42000 …near 'key(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT key(a)
< error
| 1064 42000 …near 'key(a)' at line 1
> exec ^
: string
| SELECT key(a, b)
< error
| 1064 42000 …near 'key(a, b)' at line 1
> exec ^
: string
| SELECT key(a) OVER ()
< error
| 1064 42000 …near 'key(a) OVER ()' at line 1
> exec ^
: string
| SELECT key_block_size()
< error
| 1630 42000 FUNCTION test.key_block_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT key_block_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT keyring()
< error
| 1305 42000 FUNCTION test.keyring does not exist
> exec ^
: string
| SELECT keyring()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT keys()
< error
| 1064 42000 …near 'keys()' at line 1
> exec ^
: string
| SELECT keys(1)
< error
| 1064 42000 …near 'keys(1)' at line 1
> exec ^
: string
| SELECT keys(1,2)
< error
| 1064 42000 …near 'keys(1,2)' at line 1
> exec ^
: string
| SELECT keys(1,2,3)
< error
| 1064 42000 …near 'keys(1,2,3)' at line 1
> exec ^
: string
| SELECT keys(1,2,3,4)
< error
| 1064 42000 …near 'keys(1,2,3,4)' at line 1
> exec ^
: string
| SELECT keys(*)
< error
| 1064 42000 …near 'keys(*)' at line 1
> exec ^
: string
| SELECT keys('a')
< error
| 1064 42000 …near 'keys('a')' at line 1
> exec ^
: string
| SELECT keys(a USING utf8mb4)
< error
| 1064 42000 …near 'keys(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT keys(1 AS SIGNED)
< error
| 1064 42000 …near 'keys(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT keys(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'keys(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT keys(1 FROM 2)
< error
| 1064 42000 …near 'keys(1 FROM 2)' at line 1
> exec ^
: string
| SELECT keys(1 IN 2)
< error
| 1064 42000 …near 'keys(1 IN 2)' at line 1
> exec ^
: string
| SELECT keys(DISTINCT 1)
< error
| 1064 42000 …near 'keys(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT keys() OVER ()
< error
| 1064 42000 …near 'keys() OVER ()' at line 1
> exec ^
: string
| SELECT keys(1) OVER ()
< error
| 1064 42000 …near 'keys(1) OVER ()' at line 1
> exec ^
: string
| SELECT keys(1,2) OVER ()
< error
| 1064 42000 …near 'keys(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT keys(1,2,3) OVER ()
< error
| 1064 42000 …near 'keys(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT keys(1) FROM t
< error
| 1064 42000 …near 'keys(1) FROM t' at line 1
> exec ^
: string
| SELECT keys(SELECT 1)
< error
| 1064 42000 …near 'keys(SELECT 1)' at line 1
> exec ^
: string
| SELECT keys(YEAR FROM 1)
< error
| 1064 42000 …near 'keys(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT keys(DAY, 1, 2)
< error
| 1064 42000 …near 'keys(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT keys(a)
< error
| 1064 42000 …near 'keys(a)' at line 1
> exec ^
: string
| SELECT keys(a, b)
< error
| 1064 42000 …near 'keys(a, b)' at line 1
> exec ^
: string
| SELECT keys(a) OVER ()
< error
| 1064 42000 …near 'keys(a) OVER ()' at line 1
> exec ^
: string
| SELECT kill()
< error
| 1064 42000 …near 'kill()' at line 1
> exec ^
: string
| SELECT kill(1)
< error
| 1064 42000 …near 'kill(1)' at line 1
> exec ^
: string
| SELECT kill(1,2)
< error
| 1064 42000 …near 'kill(1,2)' at line 1
> exec ^
: string
| SELECT kill(1,2,3)
< error
| 1064 42000 …near 'kill(1,2,3)' at line 1
> exec ^
: string
| SELECT kill(1,2,3,4)
< error
| 1064 42000 …near 'kill(1,2,3,4)' at line 1
> exec ^
: string
| SELECT kill(*)
< error
| 1064 42000 …near 'kill(*)' at line 1
> exec ^
: string
| SELECT kill('a')
< error
| 1064 42000 …near 'kill('a')' at line 1
> exec ^
: string
| SELECT kill(a USING utf8mb4)
< error
| 1064 42000 …near 'kill(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT kill(1 AS SIGNED)
< error
| 1064 42000 …near 'kill(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT kill(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'kill(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT kill(1 FROM 2)
< error
| 1064 42000 …near 'kill(1 FROM 2)' at line 1
> exec ^
: string
| SELECT kill(1 IN 2)
< error
| 1064 42000 …near 'kill(1 IN 2)' at line 1
> exec ^
: string
| SELECT kill(DISTINCT 1)
< error
| 1064 42000 …near 'kill(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT kill() OVER ()
< error
| 1064 42000 …near 'kill() OVER ()' at line 1
> exec ^
: string
| SELECT kill(1) OVER ()
< error
| 1064 42000 …near 'kill(1) OVER ()' at line 1
> exec ^
: string
| SELECT kill(1,2) OVER ()
< error
| 1064 42000 …near 'kill(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT kill(1,2,3) OVER ()
< error
| 1064 42000 …near 'kill(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT kill(1) FROM t
< error
| 1064 42000 …near 'kill(1) FROM t' at line 1
> exec ^
: string
| SELECT kill(SELECT 1)
< error
| 1064 42000 …near 'kill(SELECT 1)' at line 1
> exec ^
: string
| SELECT kill(YEAR FROM 1)
< error
| 1064 42000 …near 'kill(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT kill(DAY, 1, 2)
< error
| 1064 42000 …near 'kill(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT kill(a)
< error
| 1064 42000 …near 'kill(a)' at line 1
> exec ^
: string
| SELECT kill(a, b)
< error
| 1064 42000 …near 'kill(a, b)' at line 1
> exec ^
: string
| SELECT kill(a) OVER ()
< error
| 1064 42000 …near 'kill(a) OVER ()' at line 1
> exec ^
: string
| SELECT lag()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT lag(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lag(1,2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lag(1,2,3)
< error
| 1064 42000 …near '3)' at line 1
> exec ^
: string
| SELECT lag(1,2,3,4)
< error
| 1064 42000 …near '3,4)' at line 1
> exec ^
: string
| SELECT lag(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT lag('a')
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lag(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT lag(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT lag(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT lag(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT lag(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT lag(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT lag() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT lag(1) OVER ()
< error
| 4018 HY000 No order list in window specification for 'lag'
> exec ^
: string
| SELECT lag(1) OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT language()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT language(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT language(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT language(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT language(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT language(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT language('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT language(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT language(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT language(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT language(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT language(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT language(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT language() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT language(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT language(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT language(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT language(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT language(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT language(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT language(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT language(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT language(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT language(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT last()
< error
| 1630 42000 FUNCTION test.last does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT last()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT last_value()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT last_value(1)
< ok 0 0
> exec ^
: string
| SELECT lateral()
< error
| 1305 42000 FUNCTION test.lateral does not exist
> exec ^
: string
| SELECT lateral()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT lead()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT lead(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lead(1,2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lead(1,2,3)
< error
| 1064 42000 …near '3)' at line 1
> exec ^
: string
| SELECT lead(1,2,3,4)
< error
| 1064 42000 …near '3,4)' at line 1
> exec ^
: string
| SELECT lead(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT lead('a')
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT lead(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT lead(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT lead(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT lead(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT lead(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT lead(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT lead() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT lead(1) OVER ()
< error
| 4018 HY000 No order list in window specification for 'lead'
> exec ^
: string
| SELECT lead(1) OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT leading()
< error
| 1064 42000 …near 'leading()' at line 1
> exec ^
: string
| SELECT leading(1)
< error
| 1064 42000 …near 'leading(1)' at line 1
> exec ^
: string
| SELECT leading(1,2)
< error
| 1064 42000 …near 'leading(1,2)' at line 1
> exec ^
: string
| SELECT leading(1,2,3)
< error
| 1064 42000 …near 'leading(1,2,3)' at line 1
> exec ^
: string
| SELECT leading(1,2,3,4)
< error
| 1064 42000 …near 'leading(1,2,3,4)' at line 1
> exec ^
: string
| SELECT leading(*)
< error
| 1064 42000 …near 'leading(*)' at line 1
> exec ^
: string
| SELECT leading('a')
< error
| 1064 42000 …near 'leading('a')' at line 1
> exec ^
: string
| SELECT leading(a USING utf8mb4)
< error
| 1064 42000 …near 'leading(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT leading(1 AS SIGNED)
< error
| 1064 42000 …near 'leading(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT leading(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'leading(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT leading(1 FROM 2)
< error
| 1064 42000 …near 'leading(1 FROM 2)' at line 1
> exec ^
: string
| SELECT leading(1 IN 2)
< error
| 1064 42000 …near 'leading(1 IN 2)' at line 1
> exec ^
: string
| SELECT leading(DISTINCT 1)
< error
| 1064 42000 …near 'leading(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT leading() OVER ()
< error
| 1064 42000 …near 'leading() OVER ()' at line 1
> exec ^
: string
| SELECT leading(1) OVER ()
< error
| 1064 42000 …near 'leading(1) OVER ()' at line 1
> exec ^
: string
| SELECT leading(1,2) OVER ()
< error
| 1064 42000 …near 'leading(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT leading(1,2,3) OVER ()
< error
| 1064 42000 …near 'leading(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT leading(1) FROM t
< error
| 1064 42000 …near 'leading(1) FROM t' at line 1
> exec ^
: string
| SELECT leading(SELECT 1)
< error
| 1064 42000 …near 'leading(SELECT 1)' at line 1
> exec ^
: string
| SELECT leading(YEAR FROM 1)
< error
| 1064 42000 …near 'leading(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT leading(DAY, 1, 2)
< error
| 1064 42000 …near 'leading(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT leading(a)
< error
| 1064 42000 …near 'leading(a)' at line 1
> exec ^
: string
| SELECT leading(a, b)
< error
| 1064 42000 …near 'leading(a, b)' at line 1
> exec ^
: string
| SELECT leading(a) OVER ()
< error
| 1064 42000 …near 'leading(a) OVER ()' at line 1
> exec ^
: string
| SELECT least()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'least'
> exec ^
: string
| SELECT least()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'least'
> exec ^
: string
| SELECT least(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'least'
> exec ^
: string
| SELECT least(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'least'
> exec ^
: string
| SELECT least(1,2)
< ok 0 0
> exec ^
: string
| SELECT leave()
< error
| 1064 42000 …near 'leave()' at line 1
> exec ^
: string
| SELECT leave(1)
< error
| 1064 42000 …near 'leave(1)' at line 1
> exec ^
: string
| SELECT leave(1,2)
< error
| 1064 42000 …near 'leave(1,2)' at line 1
> exec ^
: string
| SELECT leave(1,2,3)
< error
| 1064 42000 …near 'leave(1,2,3)' at line 1
> exec ^
: string
| SELECT leave(1,2,3,4)
< error
| 1064 42000 …near 'leave(1,2,3,4)' at line 1
> exec ^
: string
| SELECT leave(*)
< error
| 1064 42000 …near 'leave(*)' at line 1
> exec ^
: string
| SELECT leave('a')
< error
| 1064 42000 …near 'leave('a')' at line 1
> exec ^
: string
| SELECT leave(a USING utf8mb4)
< error
| 1064 42000 …near 'leave(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT leave(1 AS SIGNED)
< error
| 1064 42000 …near 'leave(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT leave(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'leave(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT leave(1 FROM 2)
< error
| 1064 42000 …near 'leave(1 FROM 2)' at line 1
> exec ^
: string
| SELECT leave(1 IN 2)
< error
| 1064 42000 …near 'leave(1 IN 2)' at line 1
> exec ^
: string
| SELECT leave(DISTINCT 1)
< error
| 1064 42000 …near 'leave(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT leave() OVER ()
< error
| 1064 42000 …near 'leave() OVER ()' at line 1
> exec ^
: string
| SELECT leave(1) OVER ()
< error
| 1064 42000 …near 'leave(1) OVER ()' at line 1
> exec ^
: string
| SELECT leave(1,2) OVER ()
< error
| 1064 42000 …near 'leave(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT leave(1,2,3) OVER ()
< error
| 1064 42000 …near 'leave(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT leave(1) FROM t
< error
| 1064 42000 …near 'leave(1) FROM t' at line 1
> exec ^
: string
| SELECT leave(SELECT 1)
< error
| 1064 42000 …near 'leave(SELECT 1)' at line 1
> exec ^
: string
| SELECT leave(YEAR FROM 1)
< error
| 1064 42000 …near 'leave(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT leave(DAY, 1, 2)
< error
| 1064 42000 …near 'leave(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT leave(a)
< error
| 1064 42000 …near 'leave(a)' at line 1
> exec ^
: string
| SELECT leave(a, b)
< error
| 1064 42000 …near 'leave(a, b)' at line 1
> exec ^
: string
| SELECT leave(a) OVER ()
< error
| 1064 42000 …near 'leave(a) OVER ()' at line 1
> exec ^
: string
| SELECT leaves()
< error
| 1630 42000 FUNCTION test.leaves does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT leaves()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT left()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT left(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT left(1,2)
< ok 0 0
> exec ^
: string
| SELECT less()
< error
| 1630 42000 FUNCTION test.less does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT less()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT level()
< error
| 1630 42000 FUNCTION test.level does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT level()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT library()
< error
| 1305 42000 FUNCTION test.library does not exist
> exec ^
: string
| SELECT library()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT like()
< error
| 1064 42000 …near 'like()' at line 1
> exec ^
: string
| SELECT like(1)
< error
| 1064 42000 …near 'like(1)' at line 1
> exec ^
: string
| SELECT like(1,2)
< error
| 1064 42000 …near 'like(1,2)' at line 1
> exec ^
: string
| SELECT like(1,2,3)
< error
| 1064 42000 …near 'like(1,2,3)' at line 1
> exec ^
: string
| SELECT like(1,2,3,4)
< error
| 1064 42000 …near 'like(1,2,3,4)' at line 1
> exec ^
: string
| SELECT like(*)
< error
| 1064 42000 …near 'like(*)' at line 1
> exec ^
: string
| SELECT like('a')
< error
| 1064 42000 …near 'like('a')' at line 1
> exec ^
: string
| SELECT like(a USING utf8mb4)
< error
| 1064 42000 …near 'like(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT like(1 AS SIGNED)
< error
| 1064 42000 …near 'like(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT like(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'like(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT like(1 FROM 2)
< error
| 1064 42000 …near 'like(1 FROM 2)' at line 1
> exec ^
: string
| SELECT like(1 IN 2)
< error
| 1064 42000 …near 'like(1 IN 2)' at line 1
> exec ^
: string
| SELECT like(DISTINCT 1)
< error
| 1064 42000 …near 'like(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT like() OVER ()
< error
| 1064 42000 …near 'like() OVER ()' at line 1
> exec ^
: string
| SELECT like(1) OVER ()
< error
| 1064 42000 …near 'like(1) OVER ()' at line 1
> exec ^
: string
| SELECT like(1,2) OVER ()
< error
| 1064 42000 …near 'like(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT like(1,2,3) OVER ()
< error
| 1064 42000 …near 'like(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT like(1) FROM t
< error
| 1064 42000 …near 'like(1) FROM t' at line 1
> exec ^
: string
| SELECT like(SELECT 1)
< error
| 1064 42000 …near 'like(SELECT 1)' at line 1
> exec ^
: string
| SELECT like(YEAR FROM 1)
< error
| 1064 42000 …near 'like(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT like(DAY, 1, 2)
< error
| 1064 42000 …near 'like(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT like(a)
< error
| 1064 42000 …near 'like(a)' at line 1
> exec ^
: string
| SELECT like(a, b)
< error
| 1064 42000 …near 'like(a, b)' at line 1
> exec ^
: string
| SELECT like(a) OVER ()
< error
| 1064 42000 …near 'like(a) OVER ()' at line 1
> exec ^
: string
| SELECT limit()
< error
| 1064 42000 …near 'limit()' at line 1
> exec ^
: string
| SELECT limit(1)
< error
| 1064 42000 …near 'limit(1)' at line 1
> exec ^
: string
| SELECT limit(1,2)
< error
| 1064 42000 …near 'limit(1,2)' at line 1
> exec ^
: string
| SELECT limit(1,2,3)
< error
| 1064 42000 …near 'limit(1,2,3)' at line 1
> exec ^
: string
| SELECT limit(1,2,3,4)
< error
| 1064 42000 …near 'limit(1,2,3,4)' at line 1
> exec ^
: string
| SELECT limit(*)
< error
| 1064 42000 …near 'limit(*)' at line 1
> exec ^
: string
| SELECT limit('a')
< error
| 1064 42000 …near 'limit('a')' at line 1
> exec ^
: string
| SELECT limit(a USING utf8mb4)
< error
| 1064 42000 …near 'limit(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT limit(1 AS SIGNED)
< error
| 1064 42000 …near 'limit(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT limit(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'limit(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT limit(1 FROM 2)
< error
| 1064 42000 …near 'limit(1 FROM 2)' at line 1
> exec ^
: string
| SELECT limit(1 IN 2)
< error
| 1064 42000 …near 'limit(1 IN 2)' at line 1
> exec ^
: string
| SELECT limit(DISTINCT 1)
< error
| 1064 42000 …near 'limit(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT limit() OVER ()
< error
| 1064 42000 …near 'limit() OVER ()' at line 1
> exec ^
: string
| SELECT limit(1) OVER ()
< error
| 1064 42000 …near 'limit(1) OVER ()' at line 1
> exec ^
: string
| SELECT limit(1,2) OVER ()
< error
| 1064 42000 …near 'limit(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT limit(1,2,3) OVER ()
< error
| 1064 42000 …near 'limit(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT limit(1) FROM t
< error
| 1064 42000 …near 'limit(1) FROM t' at line 1
> exec ^
: string
| SELECT limit(SELECT 1)
< error
| 1064 42000 …near 'limit(SELECT 1)' at line 1
> exec ^
: string
| SELECT limit(YEAR FROM 1)
< error
| 1064 42000 …near 'limit(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT limit(DAY, 1, 2)
< error
| 1064 42000 …near 'limit(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT limit(a)
< error
| 1064 42000 …near 'limit(a)' at line 1
> exec ^
: string
| SELECT limit(a, b)
< error
| 1064 42000 …near 'limit(a, b)' at line 1
> exec ^
: string
| SELECT limit(a) OVER ()
< error
| 1064 42000 …near 'limit(a) OVER ()' at line 1
> exec ^
: string
| SELECT linear()
< error
| 1064 42000 …near 'linear()' at line 1
> exec ^
: string
| SELECT linear(1)
< error
| 1064 42000 …near 'linear(1)' at line 1
> exec ^
: string
| SELECT linear(1,2)
< error
| 1064 42000 …near 'linear(1,2)' at line 1
> exec ^
: string
| SELECT linear(1,2,3)
< error
| 1064 42000 …near 'linear(1,2,3)' at line 1
> exec ^
: string
| SELECT linear(1,2,3,4)
< error
| 1064 42000 …near 'linear(1,2,3,4)' at line 1
> exec ^
: string
| SELECT linear(*)
< error
| 1064 42000 …near 'linear(*)' at line 1
> exec ^
: string
| SELECT linear('a')
< error
| 1064 42000 …near 'linear('a')' at line 1
> exec ^
: string
| SELECT linear(a USING utf8mb4)
< error
| 1064 42000 …near 'linear(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT linear(1 AS SIGNED)
< error
| 1064 42000 …near 'linear(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT linear(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'linear(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT linear(1 FROM 2)
< error
| 1064 42000 …near 'linear(1 FROM 2)' at line 1
> exec ^
: string
| SELECT linear(1 IN 2)
< error
| 1064 42000 …near 'linear(1 IN 2)' at line 1
> exec ^
: string
| SELECT linear(DISTINCT 1)
< error
| 1064 42000 …near 'linear(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT linear() OVER ()
< error
| 1064 42000 …near 'linear() OVER ()' at line 1
> exec ^
: string
| SELECT linear(1) OVER ()
< error
| 1064 42000 …near 'linear(1) OVER ()' at line 1
> exec ^
: string
| SELECT linear(1,2) OVER ()
< error
| 1064 42000 …near 'linear(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT linear(1,2,3) OVER ()
< error
| 1064 42000 …near 'linear(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT linear(1) FROM t
< error
| 1064 42000 …near 'linear(1) FROM t' at line 1
> exec ^
: string
| SELECT linear(SELECT 1)
< error
| 1064 42000 …near 'linear(SELECT 1)' at line 1
> exec ^
: string
| SELECT linear(YEAR FROM 1)
< error
| 1064 42000 …near 'linear(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT linear(DAY, 1, 2)
< error
| 1064 42000 …near 'linear(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT linear(a)
< error
| 1064 42000 …near 'linear(a)' at line 1
> exec ^
: string
| SELECT linear(a, b)
< error
| 1064 42000 …near 'linear(a, b)' at line 1
> exec ^
: string
| SELECT linear(a) OVER ()
< error
| 1064 42000 …near 'linear(a) OVER ()' at line 1
> exec ^
: string
| SELECT lines()
< error
| 1064 42000 …near 'lines()' at line 1
> exec ^
: string
| SELECT lines(1)
< error
| 1064 42000 …near 'lines(1)' at line 1
> exec ^
: string
| SELECT lines(1,2)
< error
| 1064 42000 …near 'lines(1,2)' at line 1
> exec ^
: string
| SELECT lines(1,2,3)
< error
| 1064 42000 …near 'lines(1,2,3)' at line 1
> exec ^
: string
| SELECT lines(1,2,3,4)
< error
| 1064 42000 …near 'lines(1,2,3,4)' at line 1
> exec ^
: string
| SELECT lines(*)
< error
| 1064 42000 …near 'lines(*)' at line 1
> exec ^
: string
| SELECT lines('a')
< error
| 1064 42000 …near 'lines('a')' at line 1
> exec ^
: string
| SELECT lines(a USING utf8mb4)
< error
| 1064 42000 …near 'lines(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT lines(1 AS SIGNED)
< error
| 1064 42000 …near 'lines(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT lines(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'lines(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT lines(1 FROM 2)
< error
| 1064 42000 …near 'lines(1 FROM 2)' at line 1
> exec ^
: string
| SELECT lines(1 IN 2)
< error
| 1064 42000 …near 'lines(1 IN 2)' at line 1
> exec ^
: string
| SELECT lines(DISTINCT 1)
< error
| 1064 42000 …near 'lines(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT lines() OVER ()
< error
| 1064 42000 …near 'lines() OVER ()' at line 1
> exec ^
: string
| SELECT lines(1) OVER ()
< error
| 1064 42000 …near 'lines(1) OVER ()' at line 1
> exec ^
: string
| SELECT lines(1,2) OVER ()
< error
| 1064 42000 …near 'lines(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT lines(1,2,3) OVER ()
< error
| 1064 42000 …near 'lines(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT lines(1) FROM t
< error
| 1064 42000 …near 'lines(1) FROM t' at line 1
> exec ^
: string
| SELECT lines(SELECT 1)
< error
| 1064 42000 …near 'lines(SELECT 1)' at line 1
> exec ^
: string
| SELECT lines(YEAR FROM 1)
< error
| 1064 42000 …near 'lines(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT lines(DAY, 1, 2)
< error
| 1064 42000 …near 'lines(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT lines(a)
< error
| 1064 42000 …near 'lines(a)' at line 1
> exec ^
: string
| SELECT lines(a, b)
< error
| 1064 42000 …near 'lines(a, b)' at line 1
> exec ^
: string
| SELECT lines(a) OVER ()
< error
| 1064 42000 …near 'lines(a) OVER ()' at line 1
> exec ^
: string
| SELECT linestring()
< error
| 1305 42000 FUNCTION test.linestring does not exist
> exec ^
: string
| SELECT linestring()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT list()
< error
| 1630 42000 FUNCTION test.list does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT list()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT load()
< error
| 1064 42000 …near 'load()' at line 1
> exec ^
: string
| SELECT load(1)
< error
| 1064 42000 …near 'load(1)' at line 1
> exec ^
: string
| SELECT load(1,2)
< error
| 1064 42000 …near 'load(1,2)' at line 1
> exec ^
: string
| SELECT load(1,2,3)
< error
| 1064 42000 …near 'load(1,2,3)' at line 1
> exec ^
: string
| SELECT load(1,2,3,4)
< error
| 1064 42000 …near 'load(1,2,3,4)' at line 1
> exec ^
: string
| SELECT load(*)
< error
| 1064 42000 …near 'load(*)' at line 1
> exec ^
: string
| SELECT load('a')
< error
| 1064 42000 …near 'load('a')' at line 1
> exec ^
: string
| SELECT load(a USING utf8mb4)
< error
| 1064 42000 …near 'load(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT load(1 AS SIGNED)
< error
| 1064 42000 …near 'load(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT load(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'load(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT load(1 FROM 2)
< error
| 1064 42000 …near 'load(1 FROM 2)' at line 1
> exec ^
: string
| SELECT load(1 IN 2)
< error
| 1064 42000 …near 'load(1 IN 2)' at line 1
> exec ^
: string
| SELECT load(DISTINCT 1)
< error
| 1064 42000 …near 'load(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT load() OVER ()
< error
| 1064 42000 …near 'load() OVER ()' at line 1
> exec ^
: string
| SELECT load(1) OVER ()
< error
| 1064 42000 …near 'load(1) OVER ()' at line 1
> exec ^
: string
| SELECT load(1,2) OVER ()
< error
| 1064 42000 …near 'load(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT load(1,2,3) OVER ()
< error
| 1064 42000 …near 'load(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT load(1) FROM t
< error
| 1064 42000 …near 'load(1) FROM t' at line 1
> exec ^
: string
| SELECT load(SELECT 1)
< error
| 1064 42000 …near 'load(SELECT 1)' at line 1
> exec ^
: string
| SELECT load(YEAR FROM 1)
< error
| 1064 42000 …near 'load(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT load(DAY, 1, 2)
< error
| 1064 42000 …near 'load(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT load(a)
< error
| 1064 42000 …near 'load(a)' at line 1
> exec ^
: string
| SELECT load(a, b)
< error
| 1064 42000 …near 'load(a, b)' at line 1
> exec ^
: string
| SELECT load(a) OVER ()
< error
| 1064 42000 …near 'load(a) OVER ()' at line 1
> exec ^
: string
| SELECT local()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT local(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT local(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT local(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT local(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT local(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT local('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT local(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT local(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT local(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT local(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT local(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT local(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT local() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT local(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT local(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT local(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT local(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT local(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT local(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT local(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT local(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT local(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT local(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT localtime()
< ok 0 0
> exec ^
: string
| SELECT localtimestamp()
< ok 0 0
> exec ^
: string
| SELECT lock()
< error
| 1064 42000 …near 'lock()' at line 1
> exec ^
: string
| SELECT lock(1)
< error
| 1064 42000 …near 'lock(1)' at line 1
> exec ^
: string
| SELECT lock(1,2)
< error
| 1064 42000 …near 'lock(1,2)' at line 1
> exec ^
: string
| SELECT lock(1,2,3)
< error
| 1064 42000 …near 'lock(1,2,3)' at line 1
> exec ^
: string
| SELECT lock(1,2,3,4)
< error
| 1064 42000 …near 'lock(1,2,3,4)' at line 1
> exec ^
: string
| SELECT lock(*)
< error
| 1064 42000 …near 'lock(*)' at line 1
> exec ^
: string
| SELECT lock('a')
< error
| 1064 42000 …near 'lock('a')' at line 1
> exec ^
: string
| SELECT lock(a USING utf8mb4)
< error
| 1064 42000 …near 'lock(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT lock(1 AS SIGNED)
< error
| 1064 42000 …near 'lock(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT lock(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'lock(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT lock(1 FROM 2)
< error
| 1064 42000 …near 'lock(1 FROM 2)' at line 1
> exec ^
: string
| SELECT lock(1 IN 2)
< error
| 1064 42000 …near 'lock(1 IN 2)' at line 1
> exec ^
: string
| SELECT lock(DISTINCT 1)
< error
| 1064 42000 …near 'lock(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT lock() OVER ()
< error
| 1064 42000 …near 'lock() OVER ()' at line 1
> exec ^
: string
| SELECT lock(1) OVER ()
< error
| 1064 42000 …near 'lock(1) OVER ()' at line 1
> exec ^
: string
| SELECT lock(1,2) OVER ()
< error
| 1064 42000 …near 'lock(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT lock(1,2,3) OVER ()
< error
| 1064 42000 …near 'lock(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT lock(1) FROM t
< error
| 1064 42000 …near 'lock(1) FROM t' at line 1
> exec ^
: string
| SELECT lock(SELECT 1)
< error
| 1064 42000 …near 'lock(SELECT 1)' at line 1
> exec ^
: string
| SELECT lock(YEAR FROM 1)
< error
| 1064 42000 …near 'lock(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT lock(DAY, 1, 2)
< error
| 1064 42000 …near 'lock(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT lock(a)
< error
| 1064 42000 …near 'lock(a)' at line 1
> exec ^
: string
| SELECT lock(a, b)
< error
| 1064 42000 …near 'lock(a, b)' at line 1
> exec ^
: string
| SELECT lock(a) OVER ()
< error
| 1064 42000 …near 'lock(a) OVER ()' at line 1
> exec ^
: string
| SELECT locked()
< error
| 1630 42000 FUNCTION test.locked does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT locked()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT locks()
< error
| 1630 42000 FUNCTION test.locks does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT locks()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT log()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'log'
> exec ^
: string
| SELECT log()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'log'
> exec ^
: string
| SELECT log(1)
< ok 0 0
> exec ^
: string
| SELECT logfile()
< error
| 1630 42000 FUNCTION test.logfile does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT logfile()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT logs()
< error
| 1630 42000 FUNCTION test.logs does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT logs()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT long()
< error
| 1064 42000 …near 'long()' at line 1
> exec ^
: string
| SELECT long(1)
< error
| 1064 42000 …near 'long(1)' at line 1
> exec ^
: string
| SELECT long(1,2)
< error
| 1064 42000 …near 'long(1,2)' at line 1
> exec ^
: string
| SELECT long(1,2,3)
< error
| 1064 42000 …near 'long(1,2,3)' at line 1
> exec ^
: string
| SELECT long(1,2,3,4)
< error
| 1064 42000 …near 'long(1,2,3,4)' at line 1
> exec ^
: string
| SELECT long(*)
< error
| 1064 42000 …near 'long(*)' at line 1
> exec ^
: string
| SELECT long('a')
< error
| 1064 42000 …near 'long('a')' at line 1
> exec ^
: string
| SELECT long(a USING utf8mb4)
< error
| 1064 42000 …near 'long(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT long(1 AS SIGNED)
< error
| 1064 42000 …near 'long(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT long(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'long(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT long(1 FROM 2)
< error
| 1064 42000 …near 'long(1 FROM 2)' at line 1
> exec ^
: string
| SELECT long(1 IN 2)
< error
| 1064 42000 …near 'long(1 IN 2)' at line 1
> exec ^
: string
| SELECT long(DISTINCT 1)
< error
| 1064 42000 …near 'long(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT long() OVER ()
< error
| 1064 42000 …near 'long() OVER ()' at line 1
> exec ^
: string
| SELECT long(1) OVER ()
< error
| 1064 42000 …near 'long(1) OVER ()' at line 1
> exec ^
: string
| SELECT long(1,2) OVER ()
< error
| 1064 42000 …near 'long(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT long(1,2,3) OVER ()
< error
| 1064 42000 …near 'long(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT long(1) FROM t
< error
| 1064 42000 …near 'long(1) FROM t' at line 1
> exec ^
: string
| SELECT long(SELECT 1)
< error
| 1064 42000 …near 'long(SELECT 1)' at line 1
> exec ^
: string
| SELECT long(YEAR FROM 1)
< error
| 1064 42000 …near 'long(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT long(DAY, 1, 2)
< error
| 1064 42000 …near 'long(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT long(a)
< error
| 1064 42000 …near 'long(a)' at line 1
> exec ^
: string
| SELECT long(a, b)
< error
| 1064 42000 …near 'long(a, b)' at line 1
> exec ^
: string
| SELECT long(a) OVER ()
< error
| 1064 42000 …near 'long(a) OVER ()' at line 1
> exec ^
: string
| SELECT longblob()
< error
| 1064 42000 …near 'longblob()' at line 1
> exec ^
: string
| SELECT longblob(1)
< error
| 1064 42000 …near 'longblob(1)' at line 1
> exec ^
: string
| SELECT longblob(1,2)
< error
| 1064 42000 …near 'longblob(1,2)' at line 1
> exec ^
: string
| SELECT longblob(1,2,3)
< error
| 1064 42000 …near 'longblob(1,2,3)' at line 1
> exec ^
: string
| SELECT longblob(1,2,3,4)
< error
| 1064 42000 …near 'longblob(1,2,3,4)' at line 1
> exec ^
: string
| SELECT longblob(*)
< error
| 1064 42000 …near 'longblob(*)' at line 1
> exec ^
: string
| SELECT longblob('a')
< error
| 1064 42000 …near 'longblob('a')' at line 1
> exec ^
: string
| SELECT longblob(a USING utf8mb4)
< error
| 1064 42000 …near 'longblob(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT longblob(1 AS SIGNED)
< error
| 1064 42000 …near 'longblob(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT longblob(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'longblob(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT longblob(1 FROM 2)
< error
| 1064 42000 …near 'longblob(1 FROM 2)' at line 1
> exec ^
: string
| SELECT longblob(1 IN 2)
< error
| 1064 42000 …near 'longblob(1 IN 2)' at line 1
> exec ^
: string
| SELECT longblob(DISTINCT 1)
< error
| 1064 42000 …near 'longblob(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT longblob() OVER ()
< error
| 1064 42000 …near 'longblob() OVER ()' at line 1
> exec ^
: string
| SELECT longblob(1) OVER ()
< error
| 1064 42000 …near 'longblob(1) OVER ()' at line 1
> exec ^
: string
| SELECT longblob(1,2) OVER ()
< error
| 1064 42000 …near 'longblob(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT longblob(1,2,3) OVER ()
< error
| 1064 42000 …near 'longblob(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT longblob(1) FROM t
< error
| 1064 42000 …near 'longblob(1) FROM t' at line 1
> exec ^
: string
| SELECT longblob(SELECT 1)
< error
| 1064 42000 …near 'longblob(SELECT 1)' at line 1
> exec ^
: string
| SELECT longblob(YEAR FROM 1)
< error
| 1064 42000 …near 'longblob(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT longblob(DAY, 1, 2)
< error
| 1064 42000 …near 'longblob(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT longblob(a)
< error
| 1064 42000 …near 'longblob(a)' at line 1
> exec ^
: string
| SELECT longblob(a, b)
< error
| 1064 42000 …near 'longblob(a, b)' at line 1
> exec ^
: string
| SELECT longblob(a) OVER ()
< error
| 1064 42000 …near 'longblob(a) OVER ()' at line 1
> exec ^
: string
| SELECT longtext()
< error
| 1064 42000 …near 'longtext()' at line 1
> exec ^
: string
| SELECT longtext(1)
< error
| 1064 42000 …near 'longtext(1)' at line 1
> exec ^
: string
| SELECT longtext(1,2)
< error
| 1064 42000 …near 'longtext(1,2)' at line 1
> exec ^
: string
| SELECT longtext(1,2,3)
< error
| 1064 42000 …near 'longtext(1,2,3)' at line 1
> exec ^
: string
| SELECT longtext(1,2,3,4)
< error
| 1064 42000 …near 'longtext(1,2,3,4)' at line 1
> exec ^
: string
| SELECT longtext(*)
< error
| 1064 42000 …near 'longtext(*)' at line 1
> exec ^
: string
| SELECT longtext('a')
< error
| 1064 42000 …near 'longtext('a')' at line 1
> exec ^
: string
| SELECT longtext(a USING utf8mb4)
< error
| 1064 42000 …near 'longtext(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT longtext(1 AS SIGNED)
< error
| 1064 42000 …near 'longtext(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT longtext(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'longtext(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT longtext(1 FROM 2)
< error
| 1064 42000 …near 'longtext(1 FROM 2)' at line 1
> exec ^
: string
| SELECT longtext(1 IN 2)
< error
| 1064 42000 …near 'longtext(1 IN 2)' at line 1
> exec ^
: string
| SELECT longtext(DISTINCT 1)
< error
| 1064 42000 …near 'longtext(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT longtext() OVER ()
< error
| 1064 42000 …near 'longtext() OVER ()' at line 1
> exec ^
: string
| SELECT longtext(1) OVER ()
< error
| 1064 42000 …near 'longtext(1) OVER ()' at line 1
> exec ^
: string
| SELECT longtext(1,2) OVER ()
< error
| 1064 42000 …near 'longtext(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT longtext(1,2,3) OVER ()
< error
| 1064 42000 …near 'longtext(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT longtext(1) FROM t
< error
| 1064 42000 …near 'longtext(1) FROM t' at line 1
> exec ^
: string
| SELECT longtext(SELECT 1)
< error
| 1064 42000 …near 'longtext(SELECT 1)' at line 1
> exec ^
: string
| SELECT longtext(YEAR FROM 1)
< error
| 1064 42000 …near 'longtext(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT longtext(DAY, 1, 2)
< error
| 1064 42000 …near 'longtext(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT longtext(a)
< error
| 1064 42000 …near 'longtext(a)' at line 1
> exec ^
: string
| SELECT longtext(a, b)
< error
| 1064 42000 …near 'longtext(a, b)' at line 1
> exec ^
: string
| SELECT longtext(a) OVER ()
< error
| 1064 42000 …near 'longtext(a) OVER ()' at line 1
> exec ^
: string
| SELECT loop()
< error
| 1064 42000 …near 'loop()' at line 1
> exec ^
: string
| SELECT loop(1)
< error
| 1064 42000 …near 'loop(1)' at line 1
> exec ^
: string
| SELECT loop(1,2)
< error
| 1064 42000 …near 'loop(1,2)' at line 1
> exec ^
: string
| SELECT loop(1,2,3)
< error
| 1064 42000 …near 'loop(1,2,3)' at line 1
> exec ^
: string
| SELECT loop(1,2,3,4)
< error
| 1064 42000 …near 'loop(1,2,3,4)' at line 1
> exec ^
: string
| SELECT loop(*)
< error
| 1064 42000 …near 'loop(*)' at line 1
> exec ^
: string
| SELECT loop('a')
< error
| 1064 42000 …near 'loop('a')' at line 1
> exec ^
: string
| SELECT loop(a USING utf8mb4)
< error
| 1064 42000 …near 'loop(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT loop(1 AS SIGNED)
< error
| 1064 42000 …near 'loop(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT loop(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'loop(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT loop(1 FROM 2)
< error
| 1064 42000 …near 'loop(1 FROM 2)' at line 1
> exec ^
: string
| SELECT loop(1 IN 2)
< error
| 1064 42000 …near 'loop(1 IN 2)' at line 1
> exec ^
: string
| SELECT loop(DISTINCT 1)
< error
| 1064 42000 …near 'loop(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT loop() OVER ()
< error
| 1064 42000 …near 'loop() OVER ()' at line 1
> exec ^
: string
| SELECT loop(1) OVER ()
< error
| 1064 42000 …near 'loop(1) OVER ()' at line 1
> exec ^
: string
| SELECT loop(1,2) OVER ()
< error
| 1064 42000 …near 'loop(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT loop(1,2,3) OVER ()
< error
| 1064 42000 …near 'loop(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT loop(1) FROM t
< error
| 1064 42000 …near 'loop(1) FROM t' at line 1
> exec ^
: string
| SELECT loop(SELECT 1)
< error
| 1064 42000 …near 'loop(SELECT 1)' at line 1
> exec ^
: string
| SELECT loop(YEAR FROM 1)
< error
| 1064 42000 …near 'loop(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT loop(DAY, 1, 2)
< error
| 1064 42000 …near 'loop(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT loop(a)
< error
| 1064 42000 …near 'loop(a)' at line 1
> exec ^
: string
| SELECT loop(a, b)
< error
| 1064 42000 …near 'loop(a, b)' at line 1
> exec ^
: string
| SELECT loop(a) OVER ()
< error
| 1064 42000 …near 'loop(a) OVER ()' at line 1
> exec ^
: string
| SELECT low_priority()
< error
| 1064 42000 …near 'low_priority()' at line 1
> exec ^
: string
| SELECT low_priority(1)
< error
| 1064 42000 …near 'low_priority(1)' at line 1
> exec ^
: string
| SELECT low_priority(1,2)
< error
| 1064 42000 …near 'low_priority(1,2)' at line 1
> exec ^
: string
| SELECT low_priority(1,2,3)
< error
| 1064 42000 …near 'low_priority(1,2,3)' at line 1
> exec ^
: string
| SELECT low_priority(1,2,3,4)
< error
| 1064 42000 …near 'low_priority(1,2,3,4)' at line 1
> exec ^
: string
| SELECT low_priority(*)
< error
| 1064 42000 …near 'low_priority(*)' at line 1
> exec ^
: string
| SELECT low_priority('a')
< error
| 1064 42000 …near 'low_priority('a')' at line 1
> exec ^
: string
| SELECT low_priority(a USING utf8mb4)
< error
| 1064 42000 …near 'low_priority(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT low_priority(1 AS SIGNED)
< error
| 1064 42000 …near 'low_priority(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT low_priority(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'low_priority(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT low_priority(1 FROM 2)
< error
| 1064 42000 …near 'low_priority(1 FROM 2)' at line 1
> exec ^
: string
| SELECT low_priority(1 IN 2)
< error
| 1064 42000 …near 'low_priority(1 IN 2)' at line 1
> exec ^
: string
| SELECT low_priority(DISTINCT 1)
< error
| 1064 42000 …near 'low_priority(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT low_priority() OVER ()
< error
| 1064 42000 …near 'low_priority() OVER ()' at line 1
> exec ^
: string
| SELECT low_priority(1) OVER ()
< error
| 1064 42000 …near 'low_priority(1) OVER ()' at line 1
> exec ^
: string
| SELECT low_priority(1,2) OVER ()
< error
| 1064 42000 …near 'low_priority(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT low_priority(1,2,3) OVER ()
< error
| 1064 42000 …near 'low_priority(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT low_priority(1) FROM t
< error
| 1064 42000 …near 'low_priority(1) FROM t' at line 1
> exec ^
: string
| SELECT low_priority(SELECT 1)
< error
| 1064 42000 …near 'low_priority(SELECT 1)' at line 1
> exec ^
: string
| SELECT low_priority(YEAR FROM 1)
< error
| 1064 42000 …near 'low_priority(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT low_priority(DAY, 1, 2)
< error
| 1064 42000 …near 'low_priority(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT low_priority(a)
< error
| 1064 42000 …near 'low_priority(a)' at line 1
> exec ^
: string
| SELECT low_priority(a, b)
< error
| 1064 42000 …near 'low_priority(a, b)' at line 1
> exec ^
: string
| SELECT low_priority(a) OVER ()
< error
| 1064 42000 …near 'low_priority(a) OVER ()' at line 1
> exec ^
: string
| SELECT manual()
< error
| 1305 42000 FUNCTION test.manual does not exist
> exec ^
: string
| SELECT manual()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master()
< error
| 1630 42000 FUNCTION test.master does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_auto_position()
< error
| 1305 42000 FUNCTION test.master_auto_position does not exist
> exec ^
: string
| SELECT master_auto_position()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_bind()
< error
| 1305 42000 FUNCTION test.master_bind does not exist
> exec ^
: string
| SELECT master_bind()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_compression_algorithms()
< error
| 1305 42000 FUNCTION test.master_compression_algorithms does not exist
> exec ^
: string
| SELECT master_compression_algorithms()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_connect_retry()
< error
| 1630 42000 FUNCTION test.master_connect_retry does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_connect_retry()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_delay()
< error
| 1630 42000 FUNCTION test.master_delay does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_delay()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_heartbeat_period()
< error
| 1630 42000 FUNCTION test.master_heartbeat_period does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_heartbeat_period()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_host()
< error
| 1630 42000 FUNCTION test.master_host does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_host()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_log_file()
< error
| 1630 42000 FUNCTION test.master_log_file does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_log_file()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_log_pos()
< error
| 1630 42000 FUNCTION test.master_log_pos does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_log_pos()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_password()
< error
| 1630 42000 FUNCTION test.master_password does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_password()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_port()
< error
| 1630 42000 FUNCTION test.master_port does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_port()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_public_key_path()
< error
| 1305 42000 FUNCTION test.master_public_key_path does not exist
> exec ^
: string
| SELECT master_public_key_path()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_retry_count()
< error
| 1305 42000 FUNCTION test.master_retry_count does not exist
> exec ^
: string
| SELECT master_retry_count()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl()
< error
| 1630 42000 FUNCTION test.master_ssl does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_ca()
< error
| 1630 42000 FUNCTION test.master_ssl_ca does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_ca()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_capath()
< error
| 1630 42000 FUNCTION test.master_ssl_capath does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_capath()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_cert()
< error
| 1630 42000 FUNCTION test.master_ssl_cert does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_cert()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_cipher()
< error
| 1630 42000 FUNCTION test.master_ssl_cipher does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_cipher()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_crl()
< error
| 1630 42000 FUNCTION test.master_ssl_crl does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_crl()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_crlpath()
< error
| 1630 42000 FUNCTION test.master_ssl_crlpath does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_crlpath()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_key()
< error
| 1630 42000 FUNCTION test.master_ssl_key does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_key()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_ssl_verify_server_cert()
< error
| 1630 42000 FUNCTION test.master_ssl_verify_server_cert does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_ssl_verify_server_cert()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_tls_ciphersuites()
< error
| 1305 42000 FUNCTION test.master_tls_ciphersuites does not exist
> exec ^
: string
| SELECT master_tls_ciphersuites()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_tls_version()
< error
| 1305 42000 FUNCTION test.master_tls_version does not exist
> exec ^
: string
| SELECT master_tls_version()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_user()
< error
| 1630 42000 FUNCTION test.master_user does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT master_user()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT master_zstd_compression_level()
< error
| 1305 42000 FUNCTION test.master_zstd_compression_level does not exist
> exec ^
: string
| SELECT master_zstd_compression_level()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT match()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT match(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT match(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT match(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT match(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT match(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT match('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT match(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT match(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT match(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT match(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT match(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT match(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT match() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT match(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT match(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT match(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT match(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT match(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT match(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT match(DAY, 1, 2)
< error
| 1064 42000 …near '1, 2)' at line 1
> exec ^
: string
| SELECT match(a)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT match(a, b)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT match(a) OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec ^
: string
| SELECT materialized()
< error
| 1305 42000 FUNCTION test.materialized does not exist
> exec ^
: string
| SELECT materialized()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT max(1)
< ok 0 0
> exec ^
: string
| SELECT max_connections_per_hour()
< error
| 1630 42000 FUNCTION test.max_connections_per_hour does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_connections_per_hour()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max_queries_per_hour()
< error
| 1630 42000 FUNCTION test.max_queries_per_hour does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_queries_per_hour()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max_rows()
< error
| 1630 42000 FUNCTION test.max_rows does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_rows()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max_size()
< error
| 1630 42000 FUNCTION test.max_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max_updates_per_hour()
< error
| 1630 42000 FUNCTION test.max_updates_per_hour does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_updates_per_hour()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT max_user_connections()
< error
| 1630 42000 FUNCTION test.max_user_connections does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT max_user_connections()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT maxvalue()
< error
| 1064 42000 …near 'maxvalue()' at line 1
> exec ^
: string
| SELECT maxvalue(1)
< error
| 1064 42000 …near 'maxvalue(1)' at line 1
> exec ^
: string
| SELECT maxvalue(1,2)
< error
| 1064 42000 …near 'maxvalue(1,2)' at line 1
> exec ^
: string
| SELECT maxvalue(1,2,3)
< error
| 1064 42000 …near 'maxvalue(1,2,3)' at line 1
> exec ^
: string
| SELECT maxvalue(1,2,3,4)
< error
| 1064 42000 …near 'maxvalue(1,2,3,4)' at line 1
> exec ^
: string
| SELECT maxvalue(*)
< error
| 1064 42000 …near 'maxvalue(*)' at line 1
> exec ^
: string
| SELECT maxvalue('a')
< error
| 1064 42000 …near 'maxvalue('a')' at line 1
> exec ^
: string
| SELECT maxvalue(a USING utf8mb4)
< error
| 1064 42000 …near 'maxvalue(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT maxvalue(1 AS SIGNED)
< error
| 1064 42000 …near 'maxvalue(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT maxvalue(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'maxvalue(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT maxvalue(1 FROM 2)
< error
| 1064 42000 …near 'maxvalue(1 FROM 2)' at line 1
> exec ^
: string
| SELECT maxvalue(1 IN 2)
< error
| 1064 42000 …near 'maxvalue(1 IN 2)' at line 1
> exec ^
: string
| SELECT maxvalue(DISTINCT 1)
< error
| 1064 42000 …near 'maxvalue(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT maxvalue() OVER ()
< error
| 1064 42000 …near 'maxvalue() OVER ()' at line 1
> exec ^
: string
| SELECT maxvalue(1) OVER ()
< error
| 1064 42000 …near 'maxvalue(1) OVER ()' at line 1
> exec ^
: string
| SELECT maxvalue(1,2) OVER ()
< error
| 1064 42000 …near 'maxvalue(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT maxvalue(1,2,3) OVER ()
< error
| 1064 42000 …near 'maxvalue(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT maxvalue(1) FROM t
< error
| 1064 42000 …near 'maxvalue(1) FROM t' at line 1
> exec ^
: string
| SELECT maxvalue(SELECT 1)
< error
| 1064 42000 …near 'maxvalue(SELECT 1)' at line 1
> exec ^
: string
| SELECT maxvalue(YEAR FROM 1)
< error
| 1064 42000 …near 'maxvalue(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT maxvalue(DAY, 1, 2)
< error
| 1064 42000 …near 'maxvalue(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT maxvalue(a)
< error
| 1064 42000 …near 'maxvalue(a)' at line 1
> exec ^
: string
| SELECT maxvalue(a, b)
< error
| 1064 42000 …near 'maxvalue(a, b)' at line 1
> exec ^
: string
| SELECT maxvalue(a) OVER ()
< error
| 1064 42000 …near 'maxvalue(a) OVER ()' at line 1
> exec ^
: string
| SELECT medium()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT medium(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT medium(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT medium(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT medium(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT medium(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT medium('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT medium(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT medium(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT medium(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT medium(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT medium(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT medium(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT medium() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT medium(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT medium(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT medium(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT medium(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT medium(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT medium(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT medium(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT medium(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT medium(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT medium(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT mediumblob()
< error
| 1064 42000 …near 'mediumblob()' at line 1
> exec ^
: string
| SELECT mediumblob(1)
< error
| 1064 42000 …near 'mediumblob(1)' at line 1
> exec ^
: string
| SELECT mediumblob(1,2)
< error
| 1064 42000 …near 'mediumblob(1,2)' at line 1
> exec ^
: string
| SELECT mediumblob(1,2,3)
< error
| 1064 42000 …near 'mediumblob(1,2,3)' at line 1
> exec ^
: string
| SELECT mediumblob(1,2,3,4)
< error
| 1064 42000 …near 'mediumblob(1,2,3,4)' at line 1
> exec ^
: string
| SELECT mediumblob(*)
< error
| 1064 42000 …near 'mediumblob(*)' at line 1
> exec ^
: string
| SELECT mediumblob('a')
< error
| 1064 42000 …near 'mediumblob('a')' at line 1
> exec ^
: string
| SELECT mediumblob(a USING utf8mb4)
< error
| 1064 42000 …near 'mediumblob(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT mediumblob(1 AS SIGNED)
< error
| 1064 42000 …near 'mediumblob(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT mediumblob(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'mediumblob(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT mediumblob(1 FROM 2)
< error
| 1064 42000 …near 'mediumblob(1 FROM 2)' at line 1
> exec ^
: string
| SELECT mediumblob(1 IN 2)
< error
| 1064 42000 …near 'mediumblob(1 IN 2)' at line 1
> exec ^
: string
| SELECT mediumblob(DISTINCT 1)
< error
| 1064 42000 …near 'mediumblob(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT mediumblob() OVER ()
< error
| 1064 42000 …near 'mediumblob() OVER ()' at line 1
> exec ^
: string
| SELECT mediumblob(1) OVER ()
< error
| 1064 42000 …near 'mediumblob(1) OVER ()' at line 1
> exec ^
: string
| SELECT mediumblob(1,2) OVER ()
< error
| 1064 42000 …near 'mediumblob(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT mediumblob(1,2,3) OVER ()
< error
| 1064 42000 …near 'mediumblob(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT mediumblob(1) FROM t
< error
| 1064 42000 …near 'mediumblob(1) FROM t' at line 1
> exec ^
: string
| SELECT mediumblob(SELECT 1)
< error
| 1064 42000 …near 'mediumblob(SELECT 1)' at line 1
> exec ^
: string
| SELECT mediumblob(YEAR FROM 1)
< error
| 1064 42000 …near 'mediumblob(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT mediumblob(DAY, 1, 2)
< error
| 1064 42000 …near 'mediumblob(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT mediumblob(a)
< error
| 1064 42000 …near 'mediumblob(a)' at line 1
> exec ^
: string
| SELECT mediumblob(a, b)
< error
| 1064 42000 …near 'mediumblob(a, b)' at line 1
> exec ^
: string
| SELECT mediumblob(a) OVER ()
< error
| 1064 42000 …near 'mediumblob(a) OVER ()' at line 1
> exec ^
: string
| SELECT mediumint()
< error
| 1064 42000 …near 'mediumint()' at line 1
> exec ^
: string
| SELECT mediumint(1)
< error
| 1064 42000 …near 'mediumint(1)' at line 1
> exec ^
: string
| SELECT mediumint(1,2)
< error
| 1064 42000 …near 'mediumint(1,2)' at line 1
> exec ^
: string
| SELECT mediumint(1,2,3)
< error
| 1064 42000 …near 'mediumint(1,2,3)' at line 1
> exec ^
: string
| SELECT mediumint(1,2,3,4)
< error
| 1064 42000 …near 'mediumint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT mediumint(*)
< error
| 1064 42000 …near 'mediumint(*)' at line 1
> exec ^
: string
| SELECT mediumint('a')
< error
| 1064 42000 …near 'mediumint('a')' at line 1
> exec ^
: string
| SELECT mediumint(a USING utf8mb4)
< error
| 1064 42000 …near 'mediumint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT mediumint(1 AS SIGNED)
< error
| 1064 42000 …near 'mediumint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT mediumint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'mediumint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT mediumint(1 FROM 2)
< error
| 1064 42000 …near 'mediumint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT mediumint(1 IN 2)
< error
| 1064 42000 …near 'mediumint(1 IN 2)' at line 1
> exec ^
: string
| SELECT mediumint(DISTINCT 1)
< error
| 1064 42000 …near 'mediumint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT mediumint() OVER ()
< error
| 1064 42000 …near 'mediumint() OVER ()' at line 1
> exec ^
: string
| SELECT mediumint(1) OVER ()
< error
| 1064 42000 …near 'mediumint(1) OVER ()' at line 1
> exec ^
: string
| SELECT mediumint(1,2) OVER ()
< error
| 1064 42000 …near 'mediumint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT mediumint(1,2,3) OVER ()
< error
| 1064 42000 …near 'mediumint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT mediumint(1) FROM t
< error
| 1064 42000 …near 'mediumint(1) FROM t' at line 1
> exec ^
: string
| SELECT mediumint(SELECT 1)
< error
| 1064 42000 …near 'mediumint(SELECT 1)' at line 1
> exec ^
: string
| SELECT mediumint(YEAR FROM 1)
< error
| 1064 42000 …near 'mediumint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT mediumint(DAY, 1, 2)
< error
| 1064 42000 …near 'mediumint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT mediumint(a)
< error
| 1064 42000 …near 'mediumint(a)' at line 1
> exec ^
: string
| SELECT mediumint(a, b)
< error
| 1064 42000 …near 'mediumint(a, b)' at line 1
> exec ^
: string
| SELECT mediumint(a) OVER ()
< error
| 1064 42000 …near 'mediumint(a) OVER ()' at line 1
> exec ^
: string
| SELECT mediumtext()
< error
| 1064 42000 …near 'mediumtext()' at line 1
> exec ^
: string
| SELECT mediumtext(1)
< error
| 1064 42000 …near 'mediumtext(1)' at line 1
> exec ^
: string
| SELECT mediumtext(1,2)
< error
| 1064 42000 …near 'mediumtext(1,2)' at line 1
> exec ^
: string
| SELECT mediumtext(1,2,3)
< error
| 1064 42000 …near 'mediumtext(1,2,3)' at line 1
> exec ^
: string
| SELECT mediumtext(1,2,3,4)
< error
| 1064 42000 …near 'mediumtext(1,2,3,4)' at line 1
> exec ^
: string
| SELECT mediumtext(*)
< error
| 1064 42000 …near 'mediumtext(*)' at line 1
> exec ^
: string
| SELECT mediumtext('a')
< error
| 1064 42000 …near 'mediumtext('a')' at line 1
> exec ^
: string
| SELECT mediumtext(a USING utf8mb4)
< error
| 1064 42000 …near 'mediumtext(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT mediumtext(1 AS SIGNED)
< error
| 1064 42000 …near 'mediumtext(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT mediumtext(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'mediumtext(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT mediumtext(1 FROM 2)
< error
| 1064 42000 …near 'mediumtext(1 FROM 2)' at line 1
> exec ^
: string
| SELECT mediumtext(1 IN 2)
< error
| 1064 42000 …near 'mediumtext(1 IN 2)' at line 1
> exec ^
: string
| SELECT mediumtext(DISTINCT 1)
< error
| 1064 42000 …near 'mediumtext(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT mediumtext() OVER ()
< error
| 1064 42000 …near 'mediumtext() OVER ()' at line 1
> exec ^
: string
| SELECT mediumtext(1) OVER ()
< error
| 1064 42000 …near 'mediumtext(1) OVER ()' at line 1
> exec ^
: string
| SELECT mediumtext(1,2) OVER ()
< error
| 1064 42000 …near 'mediumtext(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT mediumtext(1,2,3) OVER ()
< error
| 1064 42000 …near 'mediumtext(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT mediumtext(1) FROM t
< error
| 1064 42000 …near 'mediumtext(1) FROM t' at line 1
> exec ^
: string
| SELECT mediumtext(SELECT 1)
< error
| 1064 42000 …near 'mediumtext(SELECT 1)' at line 1
> exec ^
: string
| SELECT mediumtext(YEAR FROM 1)
< error
| 1064 42000 …near 'mediumtext(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT mediumtext(DAY, 1, 2)
< error
| 1064 42000 …near 'mediumtext(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT mediumtext(a)
< error
| 1064 42000 …near 'mediumtext(a)' at line 1
> exec ^
: string
| SELECT mediumtext(a, b)
< error
| 1064 42000 …near 'mediumtext(a, b)' at line 1
> exec ^
: string
| SELECT mediumtext(a) OVER ()
< error
| 1064 42000 …near 'mediumtext(a) OVER ()' at line 1
> exec ^
: string
| SELECT member()
< error
| 1305 42000 FUNCTION test.member does not exist
> exec ^
: string
| SELECT member()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT memory()
< error
| 1630 42000 FUNCTION test.memory does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT memory()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT merge()
< error
| 1630 42000 FUNCTION test.merge does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT merge()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT message_text()
< error
| 1630 42000 FUNCTION test.message_text does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT message_text()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT microsecond()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'microsecond'
> exec ^
: string
| SELECT microsecond()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'microsecond'
> exec ^
: string
| SELECT microsecond(1)
< ok 0 0
> exec ^
: string
| SELECT mid()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT mid(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT mid(1,2)
< ok 0 0
> exec ^
: string
| SELECT middleint()
< error
| 1064 42000 …near 'middleint()' at line 1
> exec ^
: string
| SELECT middleint(1)
< error
| 1064 42000 …near 'middleint(1)' at line 1
> exec ^
: string
| SELECT middleint(1,2)
< error
| 1064 42000 …near 'middleint(1,2)' at line 1
> exec ^
: string
| SELECT middleint(1,2,3)
< error
| 1064 42000 …near 'middleint(1,2,3)' at line 1
> exec ^
: string
| SELECT middleint(1,2,3,4)
< error
| 1064 42000 …near 'middleint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT middleint(*)
< error
| 1064 42000 …near 'middleint(*)' at line 1
> exec ^
: string
| SELECT middleint('a')
< error
| 1064 42000 …near 'middleint('a')' at line 1
> exec ^
: string
| SELECT middleint(a USING utf8mb4)
< error
| 1064 42000 …near 'middleint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT middleint(1 AS SIGNED)
< error
| 1064 42000 …near 'middleint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT middleint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'middleint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT middleint(1 FROM 2)
< error
| 1064 42000 …near 'middleint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT middleint(1 IN 2)
< error
| 1064 42000 …near 'middleint(1 IN 2)' at line 1
> exec ^
: string
| SELECT middleint(DISTINCT 1)
< error
| 1064 42000 …near 'middleint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT middleint() OVER ()
< error
| 1064 42000 …near 'middleint() OVER ()' at line 1
> exec ^
: string
| SELECT middleint(1) OVER ()
< error
| 1064 42000 …near 'middleint(1) OVER ()' at line 1
> exec ^
: string
| SELECT middleint(1,2) OVER ()
< error
| 1064 42000 …near 'middleint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT middleint(1,2,3) OVER ()
< error
| 1064 42000 …near 'middleint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT middleint(1) FROM t
< error
| 1064 42000 …near 'middleint(1) FROM t' at line 1
> exec ^
: string
| SELECT middleint(SELECT 1)
< error
| 1064 42000 …near 'middleint(SELECT 1)' at line 1
> exec ^
: string
| SELECT middleint(YEAR FROM 1)
< error
| 1064 42000 …near 'middleint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT middleint(DAY, 1, 2)
< error
| 1064 42000 …near 'middleint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT middleint(a)
< error
| 1064 42000 …near 'middleint(a)' at line 1
> exec ^
: string
| SELECT middleint(a, b)
< error
| 1064 42000 …near 'middleint(a, b)' at line 1
> exec ^
: string
| SELECT middleint(a) OVER ()
< error
| 1064 42000 …near 'middleint(a) OVER ()' at line 1
> exec ^
: string
| SELECT migrate()
< error
| 1630 42000 FUNCTION test.migrate does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT migrate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT min()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT min(1)
< ok 0 0
> exec ^
: string
| SELECT min_rows()
< error
| 1630 42000 FUNCTION test.min_rows does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT min_rows()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT minute()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT minute(1)
< ok 0 0
> exec ^
: string
| SELECT minute_microsecond()
< error
| 1064 42000 …near 'minute_microsecond()' at line 1
> exec ^
: string
| SELECT minute_microsecond(1)
< error
| 1064 42000 …near 'minute_microsecond(1)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1,2)
< error
| 1064 42000 …near 'minute_microsecond(1,2)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1,2,3)
< error
| 1064 42000 …near 'minute_microsecond(1,2,3)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1,2,3,4)
< error
| 1064 42000 …near 'minute_microsecond(1,2,3,4)' at line 1
> exec ^
: string
| SELECT minute_microsecond(*)
< error
| 1064 42000 …near 'minute_microsecond(*)' at line 1
> exec ^
: string
| SELECT minute_microsecond('a')
< error
| 1064 42000 …near 'minute_microsecond('a')' at line 1
> exec ^
: string
| SELECT minute_microsecond(a USING utf8mb4)
< error
| 1064 42000 …near 'minute_microsecond(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1 AS SIGNED)
< error
| 1064 42000 …near 'minute_microsecond(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'minute_microsecond(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1 FROM 2)
< error
| 1064 42000 …near 'minute_microsecond(1 FROM 2)' at line 1
> exec ^
: string
| SELECT minute_microsecond(1 IN 2)
< error
| 1064 42000 …near 'minute_microsecond(1 IN 2)' at line 1
> exec ^
: string
| SELECT minute_microsecond(DISTINCT 1)
< error
| 1064 42000 …near 'minute_microsecond(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT minute_microsecond() OVER ()
< error
| 1064 42000 …near 'minute_microsecond() OVER ()' at line 1
> exec ^
: string
| SELECT minute_microsecond(1) OVER ()
< error
| 1064 42000 …near 'minute_microsecond(1) OVER ()' at line 1
> exec ^
: string
| SELECT minute_microsecond(1,2) OVER ()
< error
| 1064 42000 …near 'minute_microsecond(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT minute_microsecond(1,2,3) OVER ()
< error
| 1064 42000 …near 'minute_microsecond(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT minute_microsecond(1) FROM t
< error
| 1064 42000 …near 'minute_microsecond(1) FROM t' at line 1
> exec ^
: string
| SELECT minute_microsecond(SELECT 1)
< error
| 1064 42000 …near 'minute_microsecond(SELECT 1)' at line 1
> exec ^
: string
| SELECT minute_microsecond(YEAR FROM 1)
< error
| 1064 42000 …near 'minute_microsecond(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT minute_microsecond(DAY, 1, 2)
< error
| 1064 42000 …near 'minute_microsecond(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT minute_microsecond(a)
< error
| 1064 42000 …near 'minute_microsecond(a)' at line 1
> exec ^
: string
| SELECT minute_microsecond(a, b)
< error
| 1064 42000 …near 'minute_microsecond(a, b)' at line 1
> exec ^
: string
| SELECT minute_microsecond(a) OVER ()
< error
| 1064 42000 …near 'minute_microsecond(a) OVER ()' at line 1
> exec ^
: string
| SELECT minute_second()
< error
| 1064 42000 …near 'minute_second()' at line 1
> exec ^
: string
| SELECT minute_second(1)
< error
| 1064 42000 …near 'minute_second(1)' at line 1
> exec ^
: string
| SELECT minute_second(1,2)
< error
| 1064 42000 …near 'minute_second(1,2)' at line 1
> exec ^
: string
| SELECT minute_second(1,2,3)
< error
| 1064 42000 …near 'minute_second(1,2,3)' at line 1
> exec ^
: string
| SELECT minute_second(1,2,3,4)
< error
| 1064 42000 …near 'minute_second(1,2,3,4)' at line 1
> exec ^
: string
| SELECT minute_second(*)
< error
| 1064 42000 …near 'minute_second(*)' at line 1
> exec ^
: string
| SELECT minute_second('a')
< error
| 1064 42000 …near 'minute_second('a')' at line 1
> exec ^
: string
| SELECT minute_second(a USING utf8mb4)
< error
| 1064 42000 …near 'minute_second(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT minute_second(1 AS SIGNED)
< error
| 1064 42000 …near 'minute_second(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT minute_second(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'minute_second(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT minute_second(1 FROM 2)
< error
| 1064 42000 …near 'minute_second(1 FROM 2)' at line 1
> exec ^
: string
| SELECT minute_second(1 IN 2)
< error
| 1064 42000 …near 'minute_second(1 IN 2)' at line 1
> exec ^
: string
| SELECT minute_second(DISTINCT 1)
< error
| 1064 42000 …near 'minute_second(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT minute_second() OVER ()
< error
| 1064 42000 …near 'minute_second() OVER ()' at line 1
> exec ^
: string
| SELECT minute_second(1) OVER ()
< error
| 1064 42000 …near 'minute_second(1) OVER ()' at line 1
> exec ^
: string
| SELECT minute_second(1,2) OVER ()
< error
| 1064 42000 …near 'minute_second(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT minute_second(1,2,3) OVER ()
< error
| 1064 42000 …near 'minute_second(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT minute_second(1) FROM t
< error
| 1064 42000 …near 'minute_second(1) FROM t' at line 1
> exec ^
: string
| SELECT minute_second(SELECT 1)
< error
| 1064 42000 …near 'minute_second(SELECT 1)' at line 1
> exec ^
: string
| SELECT minute_second(YEAR FROM 1)
< error
| 1064 42000 …near 'minute_second(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT minute_second(DAY, 1, 2)
< error
| 1064 42000 …near 'minute_second(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT minute_second(a)
< error
| 1064 42000 …near 'minute_second(a)' at line 1
> exec ^
: string
| SELECT minute_second(a, b)
< error
| 1064 42000 …near 'minute_second(a, b)' at line 1
> exec ^
: string
| SELECT minute_second(a) OVER ()
< error
| 1064 42000 …near 'minute_second(a) OVER ()' at line 1
> exec ^
: string
| SELECT minvalue()
< error
| 1630 42000 FUNCTION test.minvalue does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT minvalue()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT mod()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT mod(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT mod(1,2)
< ok 0 0
> exec ^
: string
| SELECT mode()
< error
| 1630 42000 FUNCTION test.mode does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT mode()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT modifies()
< error
| 1064 42000 …near 'modifies()' at line 1
> exec ^
: string
| SELECT modifies(1)
< error
| 1064 42000 …near 'modifies(1)' at line 1
> exec ^
: string
| SELECT modifies(1,2)
< error
| 1064 42000 …near 'modifies(1,2)' at line 1
> exec ^
: string
| SELECT modifies(1,2,3)
< error
| 1064 42000 …near 'modifies(1,2,3)' at line 1
> exec ^
: string
| SELECT modifies(1,2,3,4)
< error
| 1064 42000 …near 'modifies(1,2,3,4)' at line 1
> exec ^
: string
| SELECT modifies(*)
< error
| 1064 42000 …near 'modifies(*)' at line 1
> exec ^
: string
| SELECT modifies('a')
< error
| 1064 42000 …near 'modifies('a')' at line 1
> exec ^
: string
| SELECT modifies(a USING utf8mb4)
< error
| 1064 42000 …near 'modifies(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT modifies(1 AS SIGNED)
< error
| 1064 42000 …near 'modifies(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT modifies(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'modifies(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT modifies(1 FROM 2)
< error
| 1064 42000 …near 'modifies(1 FROM 2)' at line 1
> exec ^
: string
| SELECT modifies(1 IN 2)
< error
| 1064 42000 …near 'modifies(1 IN 2)' at line 1
> exec ^
: string
| SELECT modifies(DISTINCT 1)
< error
| 1064 42000 …near 'modifies(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT modifies() OVER ()
< error
| 1064 42000 …near 'modifies() OVER ()' at line 1
> exec ^
: string
| SELECT modifies(1) OVER ()
< error
| 1064 42000 …near 'modifies(1) OVER ()' at line 1
> exec ^
: string
| SELECT modifies(1,2) OVER ()
< error
| 1064 42000 …near 'modifies(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT modifies(1,2,3) OVER ()
< error
| 1064 42000 …near 'modifies(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT modifies(1) FROM t
< error
| 1064 42000 …near 'modifies(1) FROM t' at line 1
> exec ^
: string
| SELECT modifies(SELECT 1)
< error
| 1064 42000 …near 'modifies(SELECT 1)' at line 1
> exec ^
: string
| SELECT modifies(YEAR FROM 1)
< error
| 1064 42000 …near 'modifies(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT modifies(DAY, 1, 2)
< error
| 1064 42000 …near 'modifies(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT modifies(a)
< error
| 1064 42000 …near 'modifies(a)' at line 1
> exec ^
: string
| SELECT modifies(a, b)
< error
| 1064 42000 …near 'modifies(a, b)' at line 1
> exec ^
: string
| SELECT modifies(a) OVER ()
< error
| 1064 42000 …near 'modifies(a) OVER ()' at line 1
> exec ^
: string
| SELECT modify()
< error
| 1630 42000 FUNCTION test.modify does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT modify()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT month()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT month(1)
< ok 0 0
> exec ^
: string
| SELECT multilinestring()
< error
| 1305 42000 FUNCTION test.multilinestring does not exist
> exec ^
: string
| SELECT multilinestring()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT multipoint()
< error
| 1305 42000 FUNCTION test.multipoint does not exist
> exec ^
: string
| SELECT multipoint()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT multipolygon()
< error
| 1305 42000 FUNCTION test.multipolygon does not exist
> exec ^
: string
| SELECT multipolygon()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT mutex()
< error
| 1630 42000 FUNCTION test.mutex does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT mutex()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT mysql_errno()
< error
| 1630 42000 FUNCTION test.mysql_errno does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT mysql_errno()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT name()
< error
| 1630 42000 FUNCTION test.name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT names()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT names(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT names(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT names(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT names(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT names(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT names('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT names(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT names(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT names(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT names(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT names(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT names(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT names() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT names(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT names(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT names(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT names(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT names(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT names(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT names(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT names(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT names(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT names(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT national()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT national(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT national(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT national(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT national(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT national(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT national('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT national(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT national(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT national(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT national(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT national(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT national(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT national() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT national(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT national(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT national(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT national(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT national(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT national(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT national(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT national(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT national(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT national(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT natural()
< error
| 1064 42000 …near 'natural()' at line 1
> exec ^
: string
| SELECT natural(1)
< error
| 1064 42000 …near 'natural(1)' at line 1
> exec ^
: string
| SELECT natural(1,2)
< error
| 1064 42000 …near 'natural(1,2)' at line 1
> exec ^
: string
| SELECT natural(1,2,3)
< error
| 1064 42000 …near 'natural(1,2,3)' at line 1
> exec ^
: string
| SELECT natural(1,2,3,4)
< error
| 1064 42000 …near 'natural(1,2,3,4)' at line 1
> exec ^
: string
| SELECT natural(*)
< error
| 1064 42000 …near 'natural(*)' at line 1
> exec ^
: string
| SELECT natural('a')
< error
| 1064 42000 …near 'natural('a')' at line 1
> exec ^
: string
| SELECT natural(a USING utf8mb4)
< error
| 1064 42000 …near 'natural(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT natural(1 AS SIGNED)
< error
| 1064 42000 …near 'natural(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT natural(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'natural(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT natural(1 FROM 2)
< error
| 1064 42000 …near 'natural(1 FROM 2)' at line 1
> exec ^
: string
| SELECT natural(1 IN 2)
< error
| 1064 42000 …near 'natural(1 IN 2)' at line 1
> exec ^
: string
| SELECT natural(DISTINCT 1)
< error
| 1064 42000 …near 'natural(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT natural() OVER ()
< error
| 1064 42000 …near 'natural() OVER ()' at line 1
> exec ^
: string
| SELECT natural(1) OVER ()
< error
| 1064 42000 …near 'natural(1) OVER ()' at line 1
> exec ^
: string
| SELECT natural(1,2) OVER ()
< error
| 1064 42000 …near 'natural(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT natural(1,2,3) OVER ()
< error
| 1064 42000 …near 'natural(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT natural(1) FROM t
< error
| 1064 42000 …near 'natural(1) FROM t' at line 1
> exec ^
: string
| SELECT natural(SELECT 1)
< error
| 1064 42000 …near 'natural(SELECT 1)' at line 1
> exec ^
: string
| SELECT natural(YEAR FROM 1)
< error
| 1064 42000 …near 'natural(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT natural(DAY, 1, 2)
< error
| 1064 42000 …near 'natural(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT natural(a)
< error
| 1064 42000 …near 'natural(a)' at line 1
> exec ^
: string
| SELECT natural(a, b)
< error
| 1064 42000 …near 'natural(a, b)' at line 1
> exec ^
: string
| SELECT natural(a) OVER ()
< error
| 1064 42000 …near 'natural(a) OVER ()' at line 1
> exec ^
: string
| SELECT nchar()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT nchar(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT nchar(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT nchar(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT nchar(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT nchar(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT nchar('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT nchar(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT nchar(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT nchar(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT nchar(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT nchar(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT nchar(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT nchar() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT nchar(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT nchar(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT nchar(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT nchar(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT nchar(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT nchar(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT nchar(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT nchar(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT nchar(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT nchar(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT ndb()
< error
| 1305 42000 FUNCTION test.ndb does not exist
> exec ^
: string
| SELECT ndb()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ndbcluster()
< error
| 1305 42000 FUNCTION test.ndbcluster does not exist
> exec ^
: string
| SELECT ndbcluster()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nested()
< error
| 1630 42000 FUNCTION test.nested does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nested()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT network_namespace()
< error
| 1305 42000 FUNCTION test.network_namespace does not exist
> exec ^
: string
| SELECT network_namespace()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT never()
< error
| 1630 42000 FUNCTION test.never does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT never()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT new()
< error
| 1305 42000 FUNCTION test.new does not exist
> exec ^
: string
| SELECT new()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT next()
< error
| 1630 42000 FUNCTION test.next does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT next()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT no()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT no(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT no(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT no(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT no(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT no(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT no('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT no(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT no(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT no(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT no(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT no(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT no(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT no() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT no(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT no(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT no(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT no(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT no(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT no(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT no(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT no(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT no(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT no(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT no_wait()
< error
| 1630 42000 FUNCTION test.no_wait does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT no_wait()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT no_write_to_binlog()
< error
| 1064 42000 …near 'no_write_to_binlog()' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1)
< error
| 1064 42000 …near 'no_write_to_binlog(1)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1,2)
< error
| 1064 42000 …near 'no_write_to_binlog(1,2)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1,2,3)
< error
| 1064 42000 …near 'no_write_to_binlog(1,2,3)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1,2,3,4)
< error
| 1064 42000 …near 'no_write_to_binlog(1,2,3,4)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(*)
< error
| 1064 42000 …near 'no_write_to_binlog(*)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog('a')
< error
| 1064 42000 …near 'no_write_to_binlog('a')' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(a USING utf8mb4)
< error
| 1064 42000 …near 'no_write_to_binlog(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1 AS SIGNED)
< error
| 1064 42000 …near 'no_write_to_binlog(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'no_write_to_binlog(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1 FROM 2)
< error
| 1064 42000 …near 'no_write_to_binlog(1 FROM 2)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1 IN 2)
< error
| 1064 42000 …near 'no_write_to_binlog(1 IN 2)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(DISTINCT 1)
< error
| 1064 42000 …near 'no_write_to_binlog(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog() OVER ()
< error
| 1064 42000 …near 'no_write_to_binlog() OVER ()' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1) OVER ()
< error
| 1064 42000 …near 'no_write_to_binlog(1) OVER ()' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1,2) OVER ()
< error
| 1064 42000 …near 'no_write_to_binlog(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1,2,3) OVER ()
< error
| 1064 42000 …near 'no_write_to_binlog(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(1) FROM t
< error
| 1064 42000 …near 'no_write_to_binlog(1) FROM t' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(SELECT 1)
< error
| 1064 42000 …near 'no_write_to_binlog(SELECT 1)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(YEAR FROM 1)
< error
| 1064 42000 …near 'no_write_to_binlog(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(DAY, 1, 2)
< error
| 1064 42000 …near 'no_write_to_binlog(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(a)
< error
| 1064 42000 …near 'no_write_to_binlog(a)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(a, b)
< error
| 1064 42000 …near 'no_write_to_binlog(a, b)' at line 1
> exec ^
: string
| SELECT no_write_to_binlog(a) OVER ()
< error
| 1064 42000 …near 'no_write_to_binlog(a) OVER ()' at line 1
> exec ^
: string
| SELECT nocache()
< error
| 1630 42000 FUNCTION test.nocache does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nocache()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nocycle()
< error
| 1630 42000 FUNCTION test.nocycle does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nocycle()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nodegroup()
< error
| 1630 42000 FUNCTION test.nodegroup does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nodegroup()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nomaxvalue()
< error
| 1630 42000 FUNCTION test.nomaxvalue does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nomaxvalue()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nominvalue()
< error
| 1630 42000 FUNCTION test.nominvalue does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nominvalue()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT none()
< error
| 1630 42000 FUNCTION test.none does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT none()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT not()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT not(1)
< ok 0 0
> exec ^
: string
| SELECT now()
< ok 0 0
> exec ^
: string
| SELECT nowait()
< error
| 1630 42000 FUNCTION test.nowait does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT nowait()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nth_value()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT nth_value(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT nth_value(1,2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT nth_value(1,2,3)
< error
| 1064 42000 …near '3)' at line 1
> exec ^
: string
| SELECT nth_value(1,2,3,4)
< error
| 1064 42000 …near '3,4)' at line 1
> exec ^
: string
| SELECT nth_value(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT nth_value('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT nth_value(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT nth_value(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT nth_value(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT nth_value(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT nth_value(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT nth_value(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT nth_value() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT nth_value(1) OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT nth_value(1,2) OVER ()
< ok 0 0
> exec ^
: string
| SELECT ntile()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT ntile(1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT ntile(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT ntile(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT ntile(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT ntile(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT ntile('a')
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT ntile(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT ntile(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT ntile(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ' INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT ntile(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT ntile(1 IN 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT ntile(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT ntile() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT ntile(1) OVER ()
< ok 0 0
> exec ^
: string
| SELECT null()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT null(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT null(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT null(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT null(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT null(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT null('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT null(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT null(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT null(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT null(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT null(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT null(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT null() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT null(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT null(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT null(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT null(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT null(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT null(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT null(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT null(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT null(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT null(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT nullif()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nullif'
> exec ^
: string
| SELECT nullif()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nullif'
> exec ^
: string
| SELECT nullif(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nullif'
> exec ^
: string
| SELECT nullif(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nullif'
> exec ^
: string
| SELECT nullif(1,2)
< ok 0 0
> exec ^
: string
| SELECT nulls()
< error
| 1305 42000 FUNCTION test.nulls does not exist
> exec ^
: string
| SELECT nulls()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT number()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT number(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT number(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT number(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT number(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT number(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT number('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT number(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT number(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT number(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT number(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT number(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT number(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT number() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT number(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT number(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT number(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT number(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT number(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT number(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT number(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT number(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT number(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT number(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT numeric()
< error
| 1064 42000 …near 'numeric()' at line 1
> exec ^
: string
| SELECT numeric(1)
< error
| 1064 42000 …near 'numeric(1)' at line 1
> exec ^
: string
| SELECT numeric(1,2)
< error
| 1064 42000 …near 'numeric(1,2)' at line 1
> exec ^
: string
| SELECT numeric(1,2,3)
< error
| 1064 42000 …near 'numeric(1,2,3)' at line 1
> exec ^
: string
| SELECT numeric(1,2,3,4)
< error
| 1064 42000 …near 'numeric(1,2,3,4)' at line 1
> exec ^
: string
| SELECT numeric(*)
< error
| 1064 42000 …near 'numeric(*)' at line 1
> exec ^
: string
| SELECT numeric('a')
< error
| 1064 42000 …near 'numeric('a')' at line 1
> exec ^
: string
| SELECT numeric(a USING utf8mb4)
< error
| 1064 42000 …near 'numeric(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT numeric(1 AS SIGNED)
< error
| 1064 42000 …near 'numeric(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT numeric(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'numeric(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT numeric(1 FROM 2)
< error
| 1064 42000 …near 'numeric(1 FROM 2)' at line 1
> exec ^
: string
| SELECT numeric(1 IN 2)
< error
| 1064 42000 …near 'numeric(1 IN 2)' at line 1
> exec ^
: string
| SELECT numeric(DISTINCT 1)
< error
| 1064 42000 …near 'numeric(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT numeric() OVER ()
< error
| 1064 42000 …near 'numeric() OVER ()' at line 1
> exec ^
: string
| SELECT numeric(1) OVER ()
< error
| 1064 42000 …near 'numeric(1) OVER ()' at line 1
> exec ^
: string
| SELECT numeric(1,2) OVER ()
< error
| 1064 42000 …near 'numeric(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT numeric(1,2,3) OVER ()
< error
| 1064 42000 …near 'numeric(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT numeric(1) FROM t
< error
| 1064 42000 …near 'numeric(1) FROM t' at line 1
> exec ^
: string
| SELECT numeric(SELECT 1)
< error
| 1064 42000 …near 'numeric(SELECT 1)' at line 1
> exec ^
: string
| SELECT numeric(YEAR FROM 1)
< error
| 1064 42000 …near 'numeric(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT numeric(DAY, 1, 2)
< error
| 1064 42000 …near 'numeric(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT numeric(a)
< error
| 1064 42000 …near 'numeric(a)' at line 1
> exec ^
: string
| SELECT numeric(a, b)
< error
| 1064 42000 …near 'numeric(a, b)' at line 1
> exec ^
: string
| SELECT numeric(a) OVER ()
< error
| 1064 42000 …near 'numeric(a) OVER ()' at line 1
> exec ^
: string
| SELECT nvarchar()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT nvarchar(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT nvarchar(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT nvarchar(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT nvarchar(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT nvarchar(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT nvarchar('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT nvarchar(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT nvarchar(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT nvarchar(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT nvarchar(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT nvarchar(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT nvarchar(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT nvarchar() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT nvarchar(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT nvarchar(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT nvarchar(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT nvarchar(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT nvarchar(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT nvarchar(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT nvarchar(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT nvarchar(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT nvarchar(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT nvarchar(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT of()
< error
| 1630 42000 FUNCTION test.of does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT of()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT off()
< error
| 1305 42000 FUNCTION test.off does not exist
> exec ^
: string
| SELECT off()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT offset()
< error
| 1064 42000 …near 'offset()' at line 1
> exec ^
: string
| SELECT offset(1)
< error
| 1064 42000 …near 'offset(1)' at line 1
> exec ^
: string
| SELECT offset(1,2)
< error
| 1064 42000 …near 'offset(1,2)' at line 1
> exec ^
: string
| SELECT offset(1,2,3)
< error
| 1064 42000 …near 'offset(1,2,3)' at line 1
> exec ^
: string
| SELECT offset(1,2,3,4)
< error
| 1064 42000 …near 'offset(1,2,3,4)' at line 1
> exec ^
: string
| SELECT offset(*)
< error
| 1064 42000 …near 'offset(*)' at line 1
> exec ^
: string
| SELECT offset('a')
< error
| 1064 42000 …near 'offset('a')' at line 1
> exec ^
: string
| SELECT offset(a USING utf8mb4)
< error
| 1064 42000 …near 'offset(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT offset(1 AS SIGNED)
< error
| 1064 42000 …near 'offset(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT offset(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'offset(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT offset(1 FROM 2)
< error
| 1064 42000 …near 'offset(1 FROM 2)' at line 1
> exec ^
: string
| SELECT offset(1 IN 2)
< error
| 1064 42000 …near 'offset(1 IN 2)' at line 1
> exec ^
: string
| SELECT offset(DISTINCT 1)
< error
| 1064 42000 …near 'offset(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT offset() OVER ()
< error
| 1064 42000 …near 'offset() OVER ()' at line 1
> exec ^
: string
| SELECT offset(1) OVER ()
< error
| 1064 42000 …near 'offset(1) OVER ()' at line 1
> exec ^
: string
| SELECT offset(1,2) OVER ()
< error
| 1064 42000 …near 'offset(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT offset(1,2,3) OVER ()
< error
| 1064 42000 …near 'offset(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT offset(1) FROM t
< error
| 1064 42000 …near 'offset(1) FROM t' at line 1
> exec ^
: string
| SELECT offset(SELECT 1)
< error
| 1064 42000 …near 'offset(SELECT 1)' at line 1
> exec ^
: string
| SELECT offset(YEAR FROM 1)
< error
| 1064 42000 …near 'offset(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT offset(DAY, 1, 2)
< error
| 1064 42000 …near 'offset(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT offset(a)
< error
| 1064 42000 …near 'offset(a)' at line 1
> exec ^
: string
| SELECT offset(a, b)
< error
| 1064 42000 …near 'offset(a, b)' at line 1
> exec ^
: string
| SELECT offset(a) OVER ()
< error
| 1064 42000 …near 'offset(a) OVER ()' at line 1
> exec ^
: string
| SELECT oj()
< error
| 1305 42000 FUNCTION test.oj does not exist
> exec ^
: string
| SELECT oj()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT old()
< error
| 1305 42000 FUNCTION test.old does not exist
> exec ^
: string
| SELECT old()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT on()
< error
| 1064 42000 …near 'on()' at line 1
> exec ^
: string
| SELECT on(1)
< error
| 1064 42000 …near 'on(1)' at line 1
> exec ^
: string
| SELECT on(1,2)
< error
| 1064 42000 …near 'on(1,2)' at line 1
> exec ^
: string
| SELECT on(1,2,3)
< error
| 1064 42000 …near 'on(1,2,3)' at line 1
> exec ^
: string
| SELECT on(1,2,3,4)
< error
| 1064 42000 …near 'on(1,2,3,4)' at line 1
> exec ^
: string
| SELECT on(*)
< error
| 1064 42000 …near 'on(*)' at line 1
> exec ^
: string
| SELECT on('a')
< error
| 1064 42000 …near 'on('a')' at line 1
> exec ^
: string
| SELECT on(a USING utf8mb4)
< error
| 1064 42000 …near 'on(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT on(1 AS SIGNED)
< error
| 1064 42000 …near 'on(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT on(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'on(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT on(1 FROM 2)
< error
| 1064 42000 …near 'on(1 FROM 2)' at line 1
> exec ^
: string
| SELECT on(1 IN 2)
< error
| 1064 42000 …near 'on(1 IN 2)' at line 1
> exec ^
: string
| SELECT on(DISTINCT 1)
< error
| 1064 42000 …near 'on(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT on() OVER ()
< error
| 1064 42000 …near 'on() OVER ()' at line 1
> exec ^
: string
| SELECT on(1) OVER ()
< error
| 1064 42000 …near 'on(1) OVER ()' at line 1
> exec ^
: string
| SELECT on(1,2) OVER ()
< error
| 1064 42000 …near 'on(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT on(1,2,3) OVER ()
< error
| 1064 42000 …near 'on(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT on(1) FROM t
< error
| 1064 42000 …near 'on(1) FROM t' at line 1
> exec ^
: string
| SELECT on(SELECT 1)
< error
| 1064 42000 …near 'on(SELECT 1)' at line 1
> exec ^
: string
| SELECT on(YEAR FROM 1)
< error
| 1064 42000 …near 'on(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT on(DAY, 1, 2)
< error
| 1064 42000 …near 'on(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT on(a)
< error
| 1064 42000 …near 'on(a)' at line 1
> exec ^
: string
| SELECT on(a, b)
< error
| 1064 42000 …near 'on(a, b)' at line 1
> exec ^
: string
| SELECT on(a) OVER ()
< error
| 1064 42000 …near 'on(a) OVER ()' at line 1
> exec ^
: string
| SELECT one()
< error
| 1630 42000 FUNCTION test.one does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT one()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT only()
< error
| 1630 42000 FUNCTION test.only does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT only()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT open()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT open(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT open(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT open(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT open(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT open(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT open('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT open(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT open(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT open(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT open(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT open(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT open(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT open() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT open(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT open(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT open(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT open(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT open(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT open(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT open(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT open(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT open(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT open(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT optimize()
< error
| 1064 42000 …near 'optimize()' at line 1
> exec ^
: string
| SELECT optimize(1)
< error
| 1064 42000 …near 'optimize(1)' at line 1
> exec ^
: string
| SELECT optimize(1,2)
< error
| 1064 42000 …near 'optimize(1,2)' at line 1
> exec ^
: string
| SELECT optimize(1,2,3)
< error
| 1064 42000 …near 'optimize(1,2,3)' at line 1
> exec ^
: string
| SELECT optimize(1,2,3,4)
< error
| 1064 42000 …near 'optimize(1,2,3,4)' at line 1
> exec ^
: string
| SELECT optimize(*)
< error
| 1064 42000 …near 'optimize(*)' at line 1
> exec ^
: string
| SELECT optimize('a')
< error
| 1064 42000 …near 'optimize('a')' at line 1
> exec ^
: string
| SELECT optimize(a USING utf8mb4)
< error
| 1064 42000 …near 'optimize(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT optimize(1 AS SIGNED)
< error
| 1064 42000 …near 'optimize(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT optimize(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'optimize(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT optimize(1 FROM 2)
< error
| 1064 42000 …near 'optimize(1 FROM 2)' at line 1
> exec ^
: string
| SELECT optimize(1 IN 2)
< error
| 1064 42000 …near 'optimize(1 IN 2)' at line 1
> exec ^
: string
| SELECT optimize(DISTINCT 1)
< error
| 1064 42000 …near 'optimize(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT optimize() OVER ()
< error
| 1064 42000 …near 'optimize() OVER ()' at line 1
> exec ^
: string
| SELECT optimize(1) OVER ()
< error
| 1064 42000 …near 'optimize(1) OVER ()' at line 1
> exec ^
: string
| SELECT optimize(1,2) OVER ()
< error
| 1064 42000 …near 'optimize(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT optimize(1,2,3) OVER ()
< error
| 1064 42000 …near 'optimize(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT optimize(1) FROM t
< error
| 1064 42000 …near 'optimize(1) FROM t' at line 1
> exec ^
: string
| SELECT optimize(SELECT 1)
< error
| 1064 42000 …near 'optimize(SELECT 1)' at line 1
> exec ^
: string
| SELECT optimize(YEAR FROM 1)
< error
| 1064 42000 …near 'optimize(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT optimize(DAY, 1, 2)
< error
| 1064 42000 …near 'optimize(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT optimize(a)
< error
| 1064 42000 …near 'optimize(a)' at line 1
> exec ^
: string
| SELECT optimize(a, b)
< error
| 1064 42000 …near 'optimize(a, b)' at line 1
> exec ^
: string
| SELECT optimize(a) OVER ()
< error
| 1064 42000 …near 'optimize(a) OVER ()' at line 1
> exec ^
: string
| SELECT optimizer_costs()
< error
| 1305 42000 FUNCTION test.optimizer_costs does not exist
> exec ^
: string
| SELECT optimizer_costs()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT option()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT option(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT option(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT option(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT option(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT option(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT option('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT option(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT option(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT option(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT option(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT option(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT option(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT option() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT option(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT option(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT option(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT option(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT option(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT option(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT option(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT option(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT option(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT option(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT optional()
< error
| 1305 42000 FUNCTION test.optional does not exist
> exec ^
: string
| SELECT optional()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT optionally()
< error
| 1064 42000 …near 'optionally()' at line 1
> exec ^
: string
| SELECT optionally(1)
< error
| 1064 42000 …near 'optionally(1)' at line 1
> exec ^
: string
| SELECT optionally(1,2)
< error
| 1064 42000 …near 'optionally(1,2)' at line 1
> exec ^
: string
| SELECT optionally(1,2,3)
< error
| 1064 42000 …near 'optionally(1,2,3)' at line 1
> exec ^
: string
| SELECT optionally(1,2,3,4)
< error
| 1064 42000 …near 'optionally(1,2,3,4)' at line 1
> exec ^
: string
| SELECT optionally(*)
< error
| 1064 42000 …near 'optionally(*)' at line 1
> exec ^
: string
| SELECT optionally('a')
< error
| 1064 42000 …near 'optionally('a')' at line 1
> exec ^
: string
| SELECT optionally(a USING utf8mb4)
< error
| 1064 42000 …near 'optionally(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT optionally(1 AS SIGNED)
< error
| 1064 42000 …near 'optionally(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT optionally(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'optionally(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT optionally(1 FROM 2)
< error
| 1064 42000 …near 'optionally(1 FROM 2)' at line 1
> exec ^
: string
| SELECT optionally(1 IN 2)
< error
| 1064 42000 …near 'optionally(1 IN 2)' at line 1
> exec ^
: string
| SELECT optionally(DISTINCT 1)
< error
| 1064 42000 …near 'optionally(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT optionally() OVER ()
< error
| 1064 42000 …near 'optionally() OVER ()' at line 1
> exec ^
: string
| SELECT optionally(1) OVER ()
< error
| 1064 42000 …near 'optionally(1) OVER ()' at line 1
> exec ^
: string
| SELECT optionally(1,2) OVER ()
< error
| 1064 42000 …near 'optionally(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT optionally(1,2,3) OVER ()
< error
| 1064 42000 …near 'optionally(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT optionally(1) FROM t
< error
| 1064 42000 …near 'optionally(1) FROM t' at line 1
> exec ^
: string
| SELECT optionally(SELECT 1)
< error
| 1064 42000 …near 'optionally(SELECT 1)' at line 1
> exec ^
: string
| SELECT optionally(YEAR FROM 1)
< error
| 1064 42000 …near 'optionally(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT optionally(DAY, 1, 2)
< error
| 1064 42000 …near 'optionally(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT optionally(a)
< error
| 1064 42000 …near 'optionally(a)' at line 1
> exec ^
: string
| SELECT optionally(a, b)
< error
| 1064 42000 …near 'optionally(a, b)' at line 1
> exec ^
: string
| SELECT optionally(a) OVER ()
< error
| 1064 42000 …near 'optionally(a) OVER ()' at line 1
> exec ^
: string
| SELECT options()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT options(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT options(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT options(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT options(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT options(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT options('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT options(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT options(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT options(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT options(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT options(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT options(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT options() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT options(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT options(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT options(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT options(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT options(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT options(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT options(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT options(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT options(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT options(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT or()
< error
| 1064 42000 …near 'or()' at line 1
> exec ^
: string
| SELECT or(1)
< error
| 1064 42000 …near 'or(1)' at line 1
> exec ^
: string
| SELECT or(1,2)
< error
| 1064 42000 …near 'or(1,2)' at line 1
> exec ^
: string
| SELECT or(1,2,3)
< error
| 1064 42000 …near 'or(1,2,3)' at line 1
> exec ^
: string
| SELECT or(1,2,3,4)
< error
| 1064 42000 …near 'or(1,2,3,4)' at line 1
> exec ^
: string
| SELECT or(*)
< error
| 1064 42000 …near 'or(*)' at line 1
> exec ^
: string
| SELECT or('a')
< error
| 1064 42000 …near 'or('a')' at line 1
> exec ^
: string
| SELECT or(a USING utf8mb4)
< error
| 1064 42000 …near 'or(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT or(1 AS SIGNED)
< error
| 1064 42000 …near 'or(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT or(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'or(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT or(1 FROM 2)
< error
| 1064 42000 …near 'or(1 FROM 2)' at line 1
> exec ^
: string
| SELECT or(1 IN 2)
< error
| 1064 42000 …near 'or(1 IN 2)' at line 1
> exec ^
: string
| SELECT or(DISTINCT 1)
< error
| 1064 42000 …near 'or(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT or() OVER ()
< error
| 1064 42000 …near 'or() OVER ()' at line 1
> exec ^
: string
| SELECT or(1) OVER ()
< error
| 1064 42000 …near 'or(1) OVER ()' at line 1
> exec ^
: string
| SELECT or(1,2) OVER ()
< error
| 1064 42000 …near 'or(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT or(1,2,3) OVER ()
< error
| 1064 42000 …near 'or(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT or(1) FROM t
< error
| 1064 42000 …near 'or(1) FROM t' at line 1
> exec ^
: string
| SELECT or(SELECT 1)
< error
| 1064 42000 …near 'or(SELECT 1)' at line 1
> exec ^
: string
| SELECT or(YEAR FROM 1)
< error
| 1064 42000 …near 'or(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT or(DAY, 1, 2)
< error
| 1064 42000 …near 'or(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT or(a)
< error
| 1064 42000 …near 'or(a)' at line 1
> exec ^
: string
| SELECT or(a, b)
< error
| 1064 42000 …near 'or(a, b)' at line 1
> exec ^
: string
| SELECT or(a) OVER ()
< error
| 1064 42000 …near 'or(a) OVER ()' at line 1
> exec ^
: string
| SELECT order()
< error
| 1064 42000 …near 'order()' at line 1
> exec ^
: string
| SELECT order(1)
< error
| 1064 42000 …near 'order(1)' at line 1
> exec ^
: string
| SELECT order(1,2)
< error
| 1064 42000 …near 'order(1,2)' at line 1
> exec ^
: string
| SELECT order(1,2,3)
< error
| 1064 42000 …near 'order(1,2,3)' at line 1
> exec ^
: string
| SELECT order(1,2,3,4)
< error
| 1064 42000 …near 'order(1,2,3,4)' at line 1
> exec ^
: string
| SELECT order(*)
< error
| 1064 42000 …near 'order(*)' at line 1
> exec ^
: string
| SELECT order('a')
< error
| 1064 42000 …near 'order('a')' at line 1
> exec ^
: string
| SELECT order(a USING utf8mb4)
< error
| 1064 42000 …near 'order(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT order(1 AS SIGNED)
< error
| 1064 42000 …near 'order(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT order(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'order(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT order(1 FROM 2)
< error
| 1064 42000 …near 'order(1 FROM 2)' at line 1
> exec ^
: string
| SELECT order(1 IN 2)
< error
| 1064 42000 …near 'order(1 IN 2)' at line 1
> exec ^
: string
| SELECT order(DISTINCT 1)
< error
| 1064 42000 …near 'order(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT order() OVER ()
< error
| 1064 42000 …near 'order() OVER ()' at line 1
> exec ^
: string
| SELECT order(1) OVER ()
< error
| 1064 42000 …near 'order(1) OVER ()' at line 1
> exec ^
: string
| SELECT order(1,2) OVER ()
< error
| 1064 42000 …near 'order(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT order(1,2,3) OVER ()
< error
| 1064 42000 …near 'order(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT order(1) FROM t
< error
| 1064 42000 …near 'order(1) FROM t' at line 1
> exec ^
: string
| SELECT order(SELECT 1)
< error
| 1064 42000 …near 'order(SELECT 1)' at line 1
> exec ^
: string
| SELECT order(YEAR FROM 1)
< error
| 1064 42000 …near 'order(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT order(DAY, 1, 2)
< error
| 1064 42000 …near 'order(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT order(a)
< error
| 1064 42000 …near 'order(a)' at line 1
> exec ^
: string
| SELECT order(a, b)
< error
| 1064 42000 …near 'order(a, b)' at line 1
> exec ^
: string
| SELECT order(a) OVER ()
< error
| 1064 42000 …near 'order(a) OVER ()' at line 1
> exec ^
: string
| SELECT ordinality()
< error
| 1630 42000 FUNCTION test.ordinality does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT ordinality()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT organization()
< error
| 1305 42000 FUNCTION test.organization does not exist
> exec ^
: string
| SELECT organization()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT others()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT others(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT others(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT others(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT others(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT others(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT others('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT others(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT others(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT others(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT others(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT others(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT others(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT others() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT others(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT others(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT others(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT others(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT others(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT others(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT others(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT others(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT others(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT others(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT out()
< error
| 1064 42000 …near 'out()' at line 1
> exec ^
: string
| SELECT out(1)
< error
| 1064 42000 …near 'out(1)' at line 1
> exec ^
: string
| SELECT out(1,2)
< error
| 1064 42000 …near 'out(1,2)' at line 1
> exec ^
: string
| SELECT out(1,2,3)
< error
| 1064 42000 …near 'out(1,2,3)' at line 1
> exec ^
: string
| SELECT out(1,2,3,4)
< error
| 1064 42000 …near 'out(1,2,3,4)' at line 1
> exec ^
: string
| SELECT out(*)
< error
| 1064 42000 …near 'out(*)' at line 1
> exec ^
: string
| SELECT out('a')
< error
| 1064 42000 …near 'out('a')' at line 1
> exec ^
: string
| SELECT out(a USING utf8mb4)
< error
| 1064 42000 …near 'out(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT out(1 AS SIGNED)
< error
| 1064 42000 …near 'out(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT out(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'out(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT out(1 FROM 2)
< error
| 1064 42000 …near 'out(1 FROM 2)' at line 1
> exec ^
: string
| SELECT out(1 IN 2)
< error
| 1064 42000 …near 'out(1 IN 2)' at line 1
> exec ^
: string
| SELECT out(DISTINCT 1)
< error
| 1064 42000 …near 'out(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT out() OVER ()
< error
| 1064 42000 …near 'out() OVER ()' at line 1
> exec ^
: string
| SELECT out(1) OVER ()
< error
| 1064 42000 …near 'out(1) OVER ()' at line 1
> exec ^
: string
| SELECT out(1,2) OVER ()
< error
| 1064 42000 …near 'out(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT out(1,2,3) OVER ()
< error
| 1064 42000 …near 'out(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT out(1) FROM t
< error
| 1064 42000 …near 'out(1) FROM t' at line 1
> exec ^
: string
| SELECT out(SELECT 1)
< error
| 1064 42000 …near 'out(SELECT 1)' at line 1
> exec ^
: string
| SELECT out(YEAR FROM 1)
< error
| 1064 42000 …near 'out(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT out(DAY, 1, 2)
< error
| 1064 42000 …near 'out(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT out(a)
< error
| 1064 42000 …near 'out(a)' at line 1
> exec ^
: string
| SELECT out(a, b)
< error
| 1064 42000 …near 'out(a, b)' at line 1
> exec ^
: string
| SELECT out(a) OVER ()
< error
| 1064 42000 …near 'out(a) OVER ()' at line 1
> exec ^
: string
| SELECT outer()
< error
| 1064 42000 …near 'outer()' at line 1
> exec ^
: string
| SELECT outer(1)
< error
| 1064 42000 …near 'outer(1)' at line 1
> exec ^
: string
| SELECT outer(1,2)
< error
| 1064 42000 …near 'outer(1,2)' at line 1
> exec ^
: string
| SELECT outer(1,2,3)
< error
| 1064 42000 …near 'outer(1,2,3)' at line 1
> exec ^
: string
| SELECT outer(1,2,3,4)
< error
| 1064 42000 …near 'outer(1,2,3,4)' at line 1
> exec ^
: string
| SELECT outer(*)
< error
| 1064 42000 …near 'outer(*)' at line 1
> exec ^
: string
| SELECT outer('a')
< error
| 1064 42000 …near 'outer('a')' at line 1
> exec ^
: string
| SELECT outer(a USING utf8mb4)
< error
| 1064 42000 …near 'outer(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT outer(1 AS SIGNED)
< error
| 1064 42000 …near 'outer(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT outer(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'outer(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT outer(1 FROM 2)
< error
| 1064 42000 …near 'outer(1 FROM 2)' at line 1
> exec ^
: string
| SELECT outer(1 IN 2)
< error
| 1064 42000 …near 'outer(1 IN 2)' at line 1
> exec ^
: string
| SELECT outer(DISTINCT 1)
< error
| 1064 42000 …near 'outer(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT outer() OVER ()
< error
| 1064 42000 …near 'outer() OVER ()' at line 1
> exec ^
: string
| SELECT outer(1) OVER ()
< error
| 1064 42000 …near 'outer(1) OVER ()' at line 1
> exec ^
: string
| SELECT outer(1,2) OVER ()
< error
| 1064 42000 …near 'outer(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT outer(1,2,3) OVER ()
< error
| 1064 42000 …near 'outer(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT outer(1) FROM t
< error
| 1064 42000 …near 'outer(1) FROM t' at line 1
> exec ^
: string
| SELECT outer(SELECT 1)
< error
| 1064 42000 …near 'outer(SELECT 1)' at line 1
> exec ^
: string
| SELECT outer(YEAR FROM 1)
< error
| 1064 42000 …near 'outer(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT outer(DAY, 1, 2)
< error
| 1064 42000 …near 'outer(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT outer(a)
< error
| 1064 42000 …near 'outer(a)' at line 1
> exec ^
: string
| SELECT outer(a, b)
< error
| 1064 42000 …near 'outer(a, b)' at line 1
> exec ^
: string
| SELECT outer(a) OVER ()
< error
| 1064 42000 …near 'outer(a) OVER ()' at line 1
> exec ^
: string
| SELECT outfile()
< error
| 1064 42000 …near 'outfile()' at line 1
> exec ^
: string
| SELECT outfile(1)
< error
| 1064 42000 …near 'outfile(1)' at line 1
> exec ^
: string
| SELECT outfile(1,2)
< error
| 1064 42000 …near 'outfile(1,2)' at line 1
> exec ^
: string
| SELECT outfile(1,2,3)
< error
| 1064 42000 …near 'outfile(1,2,3)' at line 1
> exec ^
: string
| SELECT outfile(1,2,3,4)
< error
| 1064 42000 …near 'outfile(1,2,3,4)' at line 1
> exec ^
: string
| SELECT outfile(*)
< error
| 1064 42000 …near 'outfile(*)' at line 1
> exec ^
: string
| SELECT outfile('a')
< error
| 1064 42000 …near 'outfile('a')' at line 1
> exec ^
: string
| SELECT outfile(a USING utf8mb4)
< error
| 1064 42000 …near 'outfile(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT outfile(1 AS SIGNED)
< error
| 1064 42000 …near 'outfile(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT outfile(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'outfile(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT outfile(1 FROM 2)
< error
| 1064 42000 …near 'outfile(1 FROM 2)' at line 1
> exec ^
: string
| SELECT outfile(1 IN 2)
< error
| 1064 42000 …near 'outfile(1 IN 2)' at line 1
> exec ^
: string
| SELECT outfile(DISTINCT 1)
< error
| 1064 42000 …near 'outfile(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT outfile() OVER ()
< error
| 1064 42000 …near 'outfile() OVER ()' at line 1
> exec ^
: string
| SELECT outfile(1) OVER ()
< error
| 1064 42000 …near 'outfile(1) OVER ()' at line 1
> exec ^
: string
| SELECT outfile(1,2) OVER ()
< error
| 1064 42000 …near 'outfile(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT outfile(1,2,3) OVER ()
< error
| 1064 42000 …near 'outfile(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT outfile(1) FROM t
< error
| 1064 42000 …near 'outfile(1) FROM t' at line 1
> exec ^
: string
| SELECT outfile(SELECT 1)
< error
| 1064 42000 …near 'outfile(SELECT 1)' at line 1
> exec ^
: string
| SELECT outfile(YEAR FROM 1)
< error
| 1064 42000 …near 'outfile(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT outfile(DAY, 1, 2)
< error
| 1064 42000 …near 'outfile(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT outfile(a)
< error
| 1064 42000 …near 'outfile(a)' at line 1
> exec ^
: string
| SELECT outfile(a, b)
< error
| 1064 42000 …near 'outfile(a, b)' at line 1
> exec ^
: string
| SELECT outfile(a) OVER ()
< error
| 1064 42000 …near 'outfile(a) OVER ()' at line 1
> exec ^
: string
| SELECT over()
< error
| 1064 42000 …near 'over()' at line 1
> exec ^
: string
| SELECT over(1)
< error
| 1064 42000 …near 'over(1)' at line 1
> exec ^
: string
| SELECT over(1,2)
< error
| 1064 42000 …near 'over(1,2)' at line 1
> exec ^
: string
| SELECT over(1,2,3)
< error
| 1064 42000 …near 'over(1,2,3)' at line 1
> exec ^
: string
| SELECT over(1,2,3,4)
< error
| 1064 42000 …near 'over(1,2,3,4)' at line 1
> exec ^
: string
| SELECT over(*)
< error
| 1064 42000 …near 'over(*)' at line 1
> exec ^
: string
| SELECT over('a')
< error
| 1064 42000 …near 'over('a')' at line 1
> exec ^
: string
| SELECT over(a USING utf8mb4)
< error
| 1064 42000 …near 'over(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT over(1 AS SIGNED)
< error
| 1064 42000 …near 'over(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT over(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'over(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT over(1 FROM 2)
< error
| 1064 42000 …near 'over(1 FROM 2)' at line 1
> exec ^
: string
| SELECT over(1 IN 2)
< error
| 1064 42000 …near 'over(1 IN 2)' at line 1
> exec ^
: string
| SELECT over(DISTINCT 1)
< error
| 1064 42000 …near 'over(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT over() OVER ()
< error
| 1064 42000 …near 'over() OVER ()' at line 1
> exec ^
: string
| SELECT over(1) OVER ()
< error
| 1064 42000 …near 'over(1) OVER ()' at line 1
> exec ^
: string
| SELECT over(1,2) OVER ()
< error
| 1064 42000 …near 'over(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT over(1,2,3) OVER ()
< error
| 1064 42000 …near 'over(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT over(1) FROM t
< error
| 1064 42000 …near 'over(1) FROM t' at line 1
> exec ^
: string
| SELECT over(SELECT 1)
< error
| 1064 42000 …near 'over(SELECT 1)' at line 1
> exec ^
: string
| SELECT over(YEAR FROM 1)
< error
| 1064 42000 …near 'over(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT over(DAY, 1, 2)
< error
| 1064 42000 …near 'over(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT over(a)
< error
| 1064 42000 …near 'over(a)' at line 1
> exec ^
: string
| SELECT over(a, b)
< error
| 1064 42000 …near 'over(a, b)' at line 1
> exec ^
: string
| SELECT over(a) OVER ()
< error
| 1064 42000 …near 'over(a) OVER ()' at line 1
> exec ^
: string
| SELECT overlaps()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'overlaps()'
> exec ^
: string
| SELECT overlaps()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'overlaps()
| ; SELECT FROM'
> exec ^
: string
| SELECT overlaps(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'overlaps(1)'
> exec ^
: string
| SELECT overlaps(1)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'overlaps(1)
| ; SELECT FROM'
> exec ^
: string
| SELECT overlaps(1,2)
< error
| 4079 HY000 Illegal parameter data type int for operation 'mbroverlaps'
> exec ^
: string
| SELECT overlaps(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT overlay()
< error
| 1305 42000 FUNCTION test.overlay does not exist
> exec ^
: string
| SELECT overlay()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT owner()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT owner(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT owner(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT owner(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT owner(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT owner(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT owner('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT owner(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT owner(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT owner(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT owner(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT owner(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT owner(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT owner() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT owner(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT owner(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT owner(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT owner(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT owner(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT owner(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT owner(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT owner(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT owner(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT owner(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT pack_keys()
< error
| 1630 42000 FUNCTION test.pack_keys does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT pack_keys()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT page()
< error
| 1630 42000 FUNCTION test.page does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT page()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT parallel()
< error
| 1305 42000 FUNCTION test.parallel does not exist
> exec ^
: string
| SELECT parallel()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT parameters()
< error
| 1305 42000 FUNCTION test.parameters does not exist
> exec ^
: string
| SELECT parameters()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT parse_tree()
< error
| 1305 42000 FUNCTION test.parse_tree does not exist
> exec ^
: string
| SELECT parse_tree()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT parser()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT parser(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT parser(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT parser(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT parser(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT parser(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT parser('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT parser(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT parser(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT parser(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT parser(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT parser(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT parser(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT parser() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT parser(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT parser(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT parser(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT parser(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT parser(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT parser(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT parser(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT parser(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT parser(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT parser(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT partial()
< error
| 1630 42000 FUNCTION test.partial does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT partial()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT partition()
< error
| 1064 42000 …near 'partition()' at line 1
> exec ^
: string
| SELECT partition(1)
< error
| 1064 42000 …near 'partition(1)' at line 1
> exec ^
: string
| SELECT partition(1,2)
< error
| 1064 42000 …near 'partition(1,2)' at line 1
> exec ^
: string
| SELECT partition(1,2,3)
< error
| 1064 42000 …near 'partition(1,2,3)' at line 1
> exec ^
: string
| SELECT partition(1,2,3,4)
< error
| 1064 42000 …near 'partition(1,2,3,4)' at line 1
> exec ^
: string
| SELECT partition(*)
< error
| 1064 42000 …near 'partition(*)' at line 1
> exec ^
: string
| SELECT partition('a')
< error
| 1064 42000 …near 'partition('a')' at line 1
> exec ^
: string
| SELECT partition(a USING utf8mb4)
< error
| 1064 42000 …near 'partition(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT partition(1 AS SIGNED)
< error
| 1064 42000 …near 'partition(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT partition(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'partition(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT partition(1 FROM 2)
< error
| 1064 42000 …near 'partition(1 FROM 2)' at line 1
> exec ^
: string
| SELECT partition(1 IN 2)
< error
| 1064 42000 …near 'partition(1 IN 2)' at line 1
> exec ^
: string
| SELECT partition(DISTINCT 1)
< error
| 1064 42000 …near 'partition(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT partition() OVER ()
< error
| 1064 42000 …near 'partition() OVER ()' at line 1
> exec ^
: string
| SELECT partition(1) OVER ()
< error
| 1064 42000 …near 'partition(1) OVER ()' at line 1
> exec ^
: string
| SELECT partition(1,2) OVER ()
< error
| 1064 42000 …near 'partition(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT partition(1,2,3) OVER ()
< error
| 1064 42000 …near 'partition(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT partition(1) FROM t
< error
| 1064 42000 …near 'partition(1) FROM t' at line 1
> exec ^
: string
| SELECT partition(SELECT 1)
< error
| 1064 42000 …near 'partition(SELECT 1)' at line 1
> exec ^
: string
| SELECT partition(YEAR FROM 1)
< error
| 1064 42000 …near 'partition(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT partition(DAY, 1, 2)
< error
| 1064 42000 …near 'partition(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT partition(a)
< error
| 1064 42000 …near 'partition(a)' at line 1
> exec ^
: string
| SELECT partition(a, b)
< error
| 1064 42000 …near 'partition(a, b)' at line 1
> exec ^
: string
| SELECT partition(a) OVER ()
< error
| 1064 42000 …near 'partition(a) OVER ()' at line 1
> exec ^
: string
| SELECT partitioning()
< error
| 1630 42000 FUNCTION test.partitioning does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT partitioning()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT partitions()
< error
| 1630 42000 FUNCTION test.partitions does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT partitions()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT password()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT password(1)
< ok 0 0
> exec ^
: string
| SELECT password_lock_time()
< error
| 1305 42000 FUNCTION test.password_lock_time does not exist
> exec ^
: string
| SELECT password_lock_time()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT path()
< error
| 1630 42000 FUNCTION test.path does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT path()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT percent_rank()
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT percent_rank(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT percent_rank(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT percent_rank(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT percent_rank(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT percent_rank(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT percent_rank('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT percent_rank(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT percent_rank(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT percent_rank(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT percent_rank(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT percent_rank(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT percent_rank(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT percent_rank() OVER ()
< error
| 4018 HY000 No order list in window specification for 'percent_rank'
> exec ^
: string
| SELECT percent_rank() OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT persist()
< error
| 1305 42000 FUNCTION test.persist does not exist
> exec ^
: string
| SELECT persist()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT persist_only()
< error
| 1305 42000 FUNCTION test.persist_only does not exist
> exec ^
: string
| SELECT persist_only()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT phase()
< error
| 1630 42000 FUNCTION test.phase does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT phase()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT placing()
< error
| 1305 42000 FUNCTION test.placing does not exist
> exec ^
: string
| SELECT placing()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT plugin()
< error
| 1630 42000 FUNCTION test.plugin does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT plugin()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT plugin_dir()
< error
| 1305 42000 FUNCTION test.plugin_dir does not exist
> exec ^
: string
| SELECT plugin_dir()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT plugins()
< error
| 1630 42000 FUNCTION test.plugins does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT plugins()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT point()
< error
| 1305 42000 FUNCTION test.point does not exist
> exec ^
: string
| SELECT point()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT polygon()
< error
| 1305 42000 FUNCTION test.polygon does not exist
> exec ^
: string
| SELECT polygon()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT port()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT port(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT port(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT port(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT port(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT port(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT port('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT port(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT port(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT port(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT port(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT port(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT port(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT port() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT port(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT port(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT port(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT port(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT port(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT port(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT port(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT port(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT port(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT port(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT portion()
< error
| 1064 42000 …near 'portion()' at line 1
> exec ^
: string
| SELECT portion(1)
< error
| 1064 42000 …near 'portion(1)' at line 1
> exec ^
: string
| SELECT portion(1,2)
< error
| 1064 42000 …near 'portion(1,2)' at line 1
> exec ^
: string
| SELECT portion(1,2,3)
< error
| 1064 42000 …near 'portion(1,2,3)' at line 1
> exec ^
: string
| SELECT portion(1,2,3,4)
< error
| 1064 42000 …near 'portion(1,2,3,4)' at line 1
> exec ^
: string
| SELECT portion(*)
< error
| 1064 42000 …near 'portion(*)' at line 1
> exec ^
: string
| SELECT portion('a')
< error
| 1064 42000 …near 'portion('a')' at line 1
> exec ^
: string
| SELECT portion(a USING utf8mb4)
< error
| 1064 42000 …near 'portion(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT portion(1 AS SIGNED)
< error
| 1064 42000 …near 'portion(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT portion(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'portion(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT portion(1 FROM 2)
< error
| 1064 42000 …near 'portion(1 FROM 2)' at line 1
> exec ^
: string
| SELECT portion(1 IN 2)
< error
| 1064 42000 …near 'portion(1 IN 2)' at line 1
> exec ^
: string
| SELECT portion(DISTINCT 1)
< error
| 1064 42000 …near 'portion(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT portion() OVER ()
< error
| 1064 42000 …near 'portion() OVER ()' at line 1
> exec ^
: string
| SELECT portion(1) OVER ()
< error
| 1064 42000 …near 'portion(1) OVER ()' at line 1
> exec ^
: string
| SELECT portion(1,2) OVER ()
< error
| 1064 42000 …near 'portion(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT portion(1,2,3) OVER ()
< error
| 1064 42000 …near 'portion(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT portion(1) FROM t
< error
| 1064 42000 …near 'portion(1) FROM t' at line 1
> exec ^
: string
| SELECT portion(SELECT 1)
< error
| 1064 42000 …near 'portion(SELECT 1)' at line 1
> exec ^
: string
| SELECT portion(YEAR FROM 1)
< error
| 1064 42000 …near 'portion(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT portion(DAY, 1, 2)
< error
| 1064 42000 …near 'portion(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT portion(a)
< error
| 1064 42000 …near 'portion(a)' at line 1
> exec ^
: string
| SELECT portion(a, b)
< error
| 1064 42000 …near 'portion(a, b)' at line 1
> exec ^
: string
| SELECT portion(a) OVER ()
< error
| 1064 42000 …near 'portion(a) OVER ()' at line 1
> exec ^
: string
| SELECT position()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT position(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT position(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT position(1,2,3)
< error
| 1064 42000 …near '2,3)' at line 1
> exec ^
: string
| SELECT position(1,2,3,4)
< error
| 1064 42000 …near '2,3,4)' at line 1
> exec ^
: string
| SELECT position(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT position('a')
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT position(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT position(1 AS SIGNED)
< error
| 1064 42000 …near 'AS SIGNED)' at line 1
> exec ^
: string
| SELECT position(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near ' INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT position(1 FROM 2)
< error
| 1064 42000 …near 'FROM 2)' at line 1
> exec ^
: string
| SELECT position(1 IN 2)
< ok 0 0
> exec ^
: string
| SELECT precedes()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT precedes(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT precedes(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT precedes(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT precedes(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT precedes(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT precedes('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT precedes(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT precedes(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT precedes(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT precedes(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT precedes(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT precedes(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT precedes() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT precedes(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT precedes(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT precedes(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT precedes(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT precedes(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT precedes(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT precedes(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT precedes(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT precedes(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT precedes(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT preceding()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT preceding(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT preceding(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT preceding(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT preceding(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT preceding(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT preceding('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT preceding(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT preceding(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT preceding(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT preceding(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT preceding(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT preceding(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT preceding() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT preceding(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT preceding(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT preceding(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT preceding(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT preceding(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT preceding(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT preceding(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT preceding(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT preceding(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT preceding(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT precision()
< error
| 1064 42000 …near 'precision()' at line 1
> exec ^
: string
| SELECT precision(1)
< error
| 1064 42000 …near 'precision(1)' at line 1
> exec ^
: string
| SELECT precision(1,2)
< error
| 1064 42000 …near 'precision(1,2)' at line 1
> exec ^
: string
| SELECT precision(1,2,3)
< error
| 1064 42000 …near 'precision(1,2,3)' at line 1
> exec ^
: string
| SELECT precision(1,2,3,4)
< error
| 1064 42000 …near 'precision(1,2,3,4)' at line 1
> exec ^
: string
| SELECT precision(*)
< error
| 1064 42000 …near 'precision(*)' at line 1
> exec ^
: string
| SELECT precision('a')
< error
| 1064 42000 …near 'precision('a')' at line 1
> exec ^
: string
| SELECT precision(a USING utf8mb4)
< error
| 1064 42000 …near 'precision(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT precision(1 AS SIGNED)
< error
| 1064 42000 …near 'precision(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT precision(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'precision(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT precision(1 FROM 2)
< error
| 1064 42000 …near 'precision(1 FROM 2)' at line 1
> exec ^
: string
| SELECT precision(1 IN 2)
< error
| 1064 42000 …near 'precision(1 IN 2)' at line 1
> exec ^
: string
| SELECT precision(DISTINCT 1)
< error
| 1064 42000 …near 'precision(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT precision() OVER ()
< error
| 1064 42000 …near 'precision() OVER ()' at line 1
> exec ^
: string
| SELECT precision(1) OVER ()
< error
| 1064 42000 …near 'precision(1) OVER ()' at line 1
> exec ^
: string
| SELECT precision(1,2) OVER ()
< error
| 1064 42000 …near 'precision(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT precision(1,2,3) OVER ()
< error
| 1064 42000 …near 'precision(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT precision(1) FROM t
< error
| 1064 42000 …near 'precision(1) FROM t' at line 1
> exec ^
: string
| SELECT precision(SELECT 1)
< error
| 1064 42000 …near 'precision(SELECT 1)' at line 1
> exec ^
: string
| SELECT precision(YEAR FROM 1)
< error
| 1064 42000 …near 'precision(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT precision(DAY, 1, 2)
< error
| 1064 42000 …near 'precision(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT precision(a)
< error
| 1064 42000 …near 'precision(a)' at line 1
> exec ^
: string
| SELECT precision(a, b)
< error
| 1064 42000 …near 'precision(a, b)' at line 1
> exec ^
: string
| SELECT precision(a) OVER ()
< error
| 1064 42000 …near 'precision(a) OVER ()' at line 1
> exec ^
: string
| SELECT prepare()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT prepare(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT prepare(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT prepare(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT prepare(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT prepare(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT prepare('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT prepare(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT prepare(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT prepare(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT prepare(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT prepare(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT prepare(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT prepare() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT prepare(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT prepare(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT prepare(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT prepare(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT prepare(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT prepare(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT prepare(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT prepare(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT prepare(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT prepare(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT preserve()
< error
| 1630 42000 FUNCTION test.preserve does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT preserve()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT prev()
< error
| 1630 42000 FUNCTION test.prev does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT prev()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT previous()
< error
| 1630 42000 FUNCTION test.previous does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT previous()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT primary()
< error
| 1064 42000 …near 'primary()' at line 1
> exec ^
: string
| SELECT primary(1)
< error
| 1064 42000 …near 'primary(1)' at line 1
> exec ^
: string
| SELECT primary(1,2)
< error
| 1064 42000 …near 'primary(1,2)' at line 1
> exec ^
: string
| SELECT primary(1,2,3)
< error
| 1064 42000 …near 'primary(1,2,3)' at line 1
> exec ^
: string
| SELECT primary(1,2,3,4)
< error
| 1064 42000 …near 'primary(1,2,3,4)' at line 1
> exec ^
: string
| SELECT primary(*)
< error
| 1064 42000 …near 'primary(*)' at line 1
> exec ^
: string
| SELECT primary('a')
< error
| 1064 42000 …near 'primary('a')' at line 1
> exec ^
: string
| SELECT primary(a USING utf8mb4)
< error
| 1064 42000 …near 'primary(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT primary(1 AS SIGNED)
< error
| 1064 42000 …near 'primary(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT primary(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'primary(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT primary(1 FROM 2)
< error
| 1064 42000 …near 'primary(1 FROM 2)' at line 1
> exec ^
: string
| SELECT primary(1 IN 2)
< error
| 1064 42000 …near 'primary(1 IN 2)' at line 1
> exec ^
: string
| SELECT primary(DISTINCT 1)
< error
| 1064 42000 …near 'primary(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT primary() OVER ()
< error
| 1064 42000 …near 'primary() OVER ()' at line 1
> exec ^
: string
| SELECT primary(1) OVER ()
< error
| 1064 42000 …near 'primary(1) OVER ()' at line 1
> exec ^
: string
| SELECT primary(1,2) OVER ()
< error
| 1064 42000 …near 'primary(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT primary(1,2,3) OVER ()
< error
| 1064 42000 …near 'primary(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT primary(1) FROM t
< error
| 1064 42000 …near 'primary(1) FROM t' at line 1
> exec ^
: string
| SELECT primary(SELECT 1)
< error
| 1064 42000 …near 'primary(SELECT 1)' at line 1
> exec ^
: string
| SELECT primary(YEAR FROM 1)
< error
| 1064 42000 …near 'primary(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT primary(DAY, 1, 2)
< error
| 1064 42000 …near 'primary(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT primary(a)
< error
| 1064 42000 …near 'primary(a)' at line 1
> exec ^
: string
| SELECT primary(a, b)
< error
| 1064 42000 …near 'primary(a, b)' at line 1
> exec ^
: string
| SELECT primary(a) OVER ()
< error
| 1064 42000 …near 'primary(a) OVER ()' at line 1
> exec ^
: string
| SELECT privilege_checks_user()
< error
| 1305 42000 FUNCTION test.privilege_checks_user does not exist
> exec ^
: string
| SELECT privilege_checks_user()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT privileges()
< error
| 1630 42000 FUNCTION test.privileges does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT privileges()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT procedure()
< error
| 1064 42000 …near 'procedure()' at line 1
> exec ^
: string
| SELECT procedure(1)
< error
| 1064 42000 …near 'procedure(1)' at line 1
> exec ^
: string
| SELECT procedure(1,2)
< error
| 1064 42000 …near 'procedure(1,2)' at line 1
> exec ^
: string
| SELECT procedure(1,2,3)
< error
| 1064 42000 …near 'procedure(1,2,3)' at line 1
> exec ^
: string
| SELECT procedure(1,2,3,4)
< error
| 1064 42000 …near 'procedure(1,2,3,4)' at line 1
> exec ^
: string
| SELECT procedure(*)
< error
| 1064 42000 …near 'procedure(*)' at line 1
> exec ^
: string
| SELECT procedure('a')
< error
| 1064 42000 …near 'procedure('a')' at line 1
> exec ^
: string
| SELECT procedure(a USING utf8mb4)
< error
| 1064 42000 …near 'procedure(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT procedure(1 AS SIGNED)
< error
| 1064 42000 …near 'procedure(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT procedure(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'procedure(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT procedure(1 FROM 2)
< error
| 1064 42000 …near 'procedure(1 FROM 2)' at line 1
> exec ^
: string
| SELECT procedure(1 IN 2)
< error
| 1064 42000 …near 'procedure(1 IN 2)' at line 1
> exec ^
: string
| SELECT procedure(DISTINCT 1)
< error
| 1064 42000 …near 'procedure(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT procedure() OVER ()
< error
| 1064 42000 …near 'procedure() OVER ()' at line 1
> exec ^
: string
| SELECT procedure(1) OVER ()
< error
| 1064 42000 …near 'procedure(1) OVER ()' at line 1
> exec ^
: string
| SELECT procedure(1,2) OVER ()
< error
| 1064 42000 …near 'procedure(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT procedure(1,2,3) OVER ()
< error
| 1064 42000 …near 'procedure(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT procedure(1) FROM t
< error
| 1064 42000 …near 'procedure(1) FROM t' at line 1
> exec ^
: string
| SELECT procedure(SELECT 1)
< error
| 1064 42000 …near 'procedure(SELECT 1)' at line 1
> exec ^
: string
| SELECT procedure(YEAR FROM 1)
< error
| 1064 42000 …near 'procedure(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT procedure(DAY, 1, 2)
< error
| 1064 42000 …near 'procedure(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT procedure(a)
< error
| 1064 42000 …near 'procedure(a)' at line 1
> exec ^
: string
| SELECT procedure(a, b)
< error
| 1064 42000 …near 'procedure(a, b)' at line 1
> exec ^
: string
| SELECT procedure(a) OVER ()
< error
| 1064 42000 …near 'procedure(a) OVER ()' at line 1
> exec ^
: string
| SELECT process()
< error
| 1630 42000 FUNCTION test.process does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT process()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT processlist()
< error
| 1630 42000 FUNCTION test.processlist does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT processlist()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT profile()
< error
| 1630 42000 FUNCTION test.profile does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT profile()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT profiles()
< error
| 1630 42000 FUNCTION test.profiles does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT profiles()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT proxy()
< error
| 1630 42000 FUNCTION test.proxy does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT proxy()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT purge()
< error
| 1064 42000 …near 'purge()' at line 1
> exec ^
: string
| SELECT purge(1)
< error
| 1064 42000 …near 'purge(1)' at line 1
> exec ^
: string
| SELECT purge(1,2)
< error
| 1064 42000 …near 'purge(1,2)' at line 1
> exec ^
: string
| SELECT purge(1,2,3)
< error
| 1064 42000 …near 'purge(1,2,3)' at line 1
> exec ^
: string
| SELECT purge(1,2,3,4)
< error
| 1064 42000 …near 'purge(1,2,3,4)' at line 1
> exec ^
: string
| SELECT purge(*)
< error
| 1064 42000 …near 'purge(*)' at line 1
> exec ^
: string
| SELECT purge('a')
< error
| 1064 42000 …near 'purge('a')' at line 1
> exec ^
: string
| SELECT purge(a USING utf8mb4)
< error
| 1064 42000 …near 'purge(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT purge(1 AS SIGNED)
< error
| 1064 42000 …near 'purge(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT purge(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'purge(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT purge(1 FROM 2)
< error
| 1064 42000 …near 'purge(1 FROM 2)' at line 1
> exec ^
: string
| SELECT purge(1 IN 2)
< error
| 1064 42000 …near 'purge(1 IN 2)' at line 1
> exec ^
: string
| SELECT purge(DISTINCT 1)
< error
| 1064 42000 …near 'purge(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT purge() OVER ()
< error
| 1064 42000 …near 'purge() OVER ()' at line 1
> exec ^
: string
| SELECT purge(1) OVER ()
< error
| 1064 42000 …near 'purge(1) OVER ()' at line 1
> exec ^
: string
| SELECT purge(1,2) OVER ()
< error
| 1064 42000 …near 'purge(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT purge(1,2,3) OVER ()
< error
| 1064 42000 …near 'purge(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT purge(1) FROM t
< error
| 1064 42000 …near 'purge(1) FROM t' at line 1
> exec ^
: string
| SELECT purge(SELECT 1)
< error
| 1064 42000 …near 'purge(SELECT 1)' at line 1
> exec ^
: string
| SELECT purge(YEAR FROM 1)
< error
| 1064 42000 …near 'purge(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT purge(DAY, 1, 2)
< error
| 1064 42000 …near 'purge(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT purge(a)
< error
| 1064 42000 …near 'purge(a)' at line 1
> exec ^
: string
| SELECT purge(a, b)
< error
| 1064 42000 …near 'purge(a, b)' at line 1
> exec ^
: string
| SELECT purge(a) OVER ()
< error
| 1064 42000 …near 'purge(a) OVER ()' at line 1
> exec ^
: string
| SELECT qualify()
< error
| 1305 42000 FUNCTION test.qualify does not exist
> exec ^
: string
| SELECT qualify()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT quarter()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'quarter'
> exec ^
: string
| SELECT quarter()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'quarter'
> exec ^
: string
| SELECT quarter(1)
< ok 0 0
> exec ^
: string
| SELECT query()
< error
| 1630 42000 FUNCTION test.query does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT query()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT quick()
< error
| 1630 42000 FUNCTION test.quick does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT quick()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT random()
< error
| 1305 42000 FUNCTION test.random does not exist
> exec ^
: string
| SELECT random()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT range()
< error
| 1064 42000 …near 'range()' at line 1
> exec ^
: string
| SELECT range(1)
< error
| 1064 42000 …near 'range(1)' at line 1
> exec ^
: string
| SELECT range(1,2)
< error
| 1064 42000 …near 'range(1,2)' at line 1
> exec ^
: string
| SELECT range(1,2,3)
< error
| 1064 42000 …near 'range(1,2,3)' at line 1
> exec ^
: string
| SELECT range(1,2,3,4)
< error
| 1064 42000 …near 'range(1,2,3,4)' at line 1
> exec ^
: string
| SELECT range(*)
< error
| 1064 42000 …near 'range(*)' at line 1
> exec ^
: string
| SELECT range('a')
< error
| 1064 42000 …near 'range('a')' at line 1
> exec ^
: string
| SELECT range(a USING utf8mb4)
< error
| 1064 42000 …near 'range(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT range(1 AS SIGNED)
< error
| 1064 42000 …near 'range(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT range(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'range(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT range(1 FROM 2)
< error
| 1064 42000 …near 'range(1 FROM 2)' at line 1
> exec ^
: string
| SELECT range(1 IN 2)
< error
| 1064 42000 …near 'range(1 IN 2)' at line 1
> exec ^
: string
| SELECT range(DISTINCT 1)
< error
| 1064 42000 …near 'range(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT range() OVER ()
< error
| 1064 42000 …near 'range() OVER ()' at line 1
> exec ^
: string
| SELECT range(1) OVER ()
< error
| 1064 42000 …near 'range(1) OVER ()' at line 1
> exec ^
: string
| SELECT range(1,2) OVER ()
< error
| 1064 42000 …near 'range(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT range(1,2,3) OVER ()
< error
| 1064 42000 …near 'range(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT range(1) FROM t
< error
| 1064 42000 …near 'range(1) FROM t' at line 1
> exec ^
: string
| SELECT range(SELECT 1)
< error
| 1064 42000 …near 'range(SELECT 1)' at line 1
> exec ^
: string
| SELECT range(YEAR FROM 1)
< error
| 1064 42000 …near 'range(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT range(DAY, 1, 2)
< error
| 1064 42000 …near 'range(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT range(a)
< error
| 1064 42000 …near 'range(a)' at line 1
> exec ^
: string
| SELECT range(a, b)
< error
| 1064 42000 …near 'range(a, b)' at line 1
> exec ^
: string
| SELECT range(a) OVER ()
< error
| 1064 42000 …near 'range(a) OVER ()' at line 1
> exec ^
: string
| SELECT rank()
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT rank(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT rank(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT rank(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT rank(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT rank(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT rank('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT rank(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT rank(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT rank(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT rank(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT rank(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT rank(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT rank() OVER ()
< error
| 4018 HY000 No order list in window specification for 'rank'
> exec ^
: string
| SELECT rank() OVER ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT read()
< error
| 1064 42000 …near 'read()' at line 1
> exec ^
: string
| SELECT read(1)
< error
| 1064 42000 …near 'read(1)' at line 1
> exec ^
: string
| SELECT read(1,2)
< error
| 1064 42000 …near 'read(1,2)' at line 1
> exec ^
: string
| SELECT read(1,2,3)
< error
| 1064 42000 …near 'read(1,2,3)' at line 1
> exec ^
: string
| SELECT read(1,2,3,4)
< error
| 1064 42000 …near 'read(1,2,3,4)' at line 1
> exec ^
: string
| SELECT read(*)
< error
| 1064 42000 …near 'read(*)' at line 1
> exec ^
: string
| SELECT read('a')
< error
| 1064 42000 …near 'read('a')' at line 1
> exec ^
: string
| SELECT read(a USING utf8mb4)
< error
| 1064 42000 …near 'read(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT read(1 AS SIGNED)
< error
| 1064 42000 …near 'read(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT read(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'read(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT read(1 FROM 2)
< error
| 1064 42000 …near 'read(1 FROM 2)' at line 1
> exec ^
: string
| SELECT read(1 IN 2)
< error
| 1064 42000 …near 'read(1 IN 2)' at line 1
> exec ^
: string
| SELECT read(DISTINCT 1)
< error
| 1064 42000 …near 'read(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT read() OVER ()
< error
| 1064 42000 …near 'read() OVER ()' at line 1
> exec ^
: string
| SELECT read(1) OVER ()
< error
| 1064 42000 …near 'read(1) OVER ()' at line 1
> exec ^
: string
| SELECT read(1,2) OVER ()
< error
| 1064 42000 …near 'read(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT read(1,2,3) OVER ()
< error
| 1064 42000 …near 'read(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT read(1) FROM t
< error
| 1064 42000 …near 'read(1) FROM t' at line 1
> exec ^
: string
| SELECT read(SELECT 1)
< error
| 1064 42000 …near 'read(SELECT 1)' at line 1
> exec ^
: string
| SELECT read(YEAR FROM 1)
< error
| 1064 42000 …near 'read(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT read(DAY, 1, 2)
< error
| 1064 42000 …near 'read(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT read(a)
< error
| 1064 42000 …near 'read(a)' at line 1
> exec ^
: string
| SELECT read(a, b)
< error
| 1064 42000 …near 'read(a, b)' at line 1
> exec ^
: string
| SELECT read(a) OVER ()
< error
| 1064 42000 …near 'read(a) OVER ()' at line 1
> exec ^
: string
| SELECT read_only()
< error
| 1630 42000 FUNCTION test.read_only does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT read_only()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT read_write()
< error
| 1064 42000 …near 'read_write()' at line 1
> exec ^
: string
| SELECT read_write(1)
< error
| 1064 42000 …near 'read_write(1)' at line 1
> exec ^
: string
| SELECT read_write(1,2)
< error
| 1064 42000 …near 'read_write(1,2)' at line 1
> exec ^
: string
| SELECT read_write(1,2,3)
< error
| 1064 42000 …near 'read_write(1,2,3)' at line 1
> exec ^
: string
| SELECT read_write(1,2,3,4)
< error
| 1064 42000 …near 'read_write(1,2,3,4)' at line 1
> exec ^
: string
| SELECT read_write(*)
< error
| 1064 42000 …near 'read_write(*)' at line 1
> exec ^
: string
| SELECT read_write('a')
< error
| 1064 42000 …near 'read_write('a')' at line 1
> exec ^
: string
| SELECT read_write(a USING utf8mb4)
< error
| 1064 42000 …near 'read_write(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT read_write(1 AS SIGNED)
< error
| 1064 42000 …near 'read_write(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT read_write(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'read_write(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT read_write(1 FROM 2)
< error
| 1064 42000 …near 'read_write(1 FROM 2)' at line 1
> exec ^
: string
| SELECT read_write(1 IN 2)
< error
| 1064 42000 …near 'read_write(1 IN 2)' at line 1
> exec ^
: string
| SELECT read_write(DISTINCT 1)
< error
| 1064 42000 …near 'read_write(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT read_write() OVER ()
< error
| 1064 42000 …near 'read_write() OVER ()' at line 1
> exec ^
: string
| SELECT read_write(1) OVER ()
< error
| 1064 42000 …near 'read_write(1) OVER ()' at line 1
> exec ^
: string
| SELECT read_write(1,2) OVER ()
< error
| 1064 42000 …near 'read_write(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT read_write(1,2,3) OVER ()
< error
| 1064 42000 …near 'read_write(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT read_write(1) FROM t
< error
| 1064 42000 …near 'read_write(1) FROM t' at line 1
> exec ^
: string
| SELECT read_write(SELECT 1)
< error
| 1064 42000 …near 'read_write(SELECT 1)' at line 1
> exec ^
: string
| SELECT read_write(YEAR FROM 1)
< error
| 1064 42000 …near 'read_write(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT read_write(DAY, 1, 2)
< error
| 1064 42000 …near 'read_write(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT read_write(a)
< error
| 1064 42000 …near 'read_write(a)' at line 1
> exec ^
: string
| SELECT read_write(a, b)
< error
| 1064 42000 …near 'read_write(a, b)' at line 1
> exec ^
: string
| SELECT read_write(a) OVER ()
< error
| 1064 42000 …near 'read_write(a) OVER ()' at line 1
> exec ^
: string
| SELECT reads()
< error
| 1064 42000 …near 'reads()' at line 1
> exec ^
: string
| SELECT reads(1)
< error
| 1064 42000 …near 'reads(1)' at line 1
> exec ^
: string
| SELECT reads(1,2)
< error
| 1064 42000 …near 'reads(1,2)' at line 1
> exec ^
: string
| SELECT reads(1,2,3)
< error
| 1064 42000 …near 'reads(1,2,3)' at line 1
> exec ^
: string
| SELECT reads(1,2,3,4)
< error
| 1064 42000 …near 'reads(1,2,3,4)' at line 1
> exec ^
: string
| SELECT reads(*)
< error
| 1064 42000 …near 'reads(*)' at line 1
> exec ^
: string
| SELECT reads('a')
< error
| 1064 42000 …near 'reads('a')' at line 1
> exec ^
: string
| SELECT reads(a USING utf8mb4)
< error
| 1064 42000 …near 'reads(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT reads(1 AS SIGNED)
< error
| 1064 42000 …near 'reads(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT reads(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'reads(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT reads(1 FROM 2)
< error
| 1064 42000 …near 'reads(1 FROM 2)' at line 1
> exec ^
: string
| SELECT reads(1 IN 2)
< error
| 1064 42000 …near 'reads(1 IN 2)' at line 1
> exec ^
: string
| SELECT reads(DISTINCT 1)
< error
| 1064 42000 …near 'reads(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT reads() OVER ()
< error
| 1064 42000 …near 'reads() OVER ()' at line 1
> exec ^
: string
| SELECT reads(1) OVER ()
< error
| 1064 42000 …near 'reads(1) OVER ()' at line 1
> exec ^
: string
| SELECT reads(1,2) OVER ()
< error
| 1064 42000 …near 'reads(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT reads(1,2,3) OVER ()
< error
| 1064 42000 …near 'reads(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT reads(1) FROM t
< error
| 1064 42000 …near 'reads(1) FROM t' at line 1
> exec ^
: string
| SELECT reads(SELECT 1)
< error
| 1064 42000 …near 'reads(SELECT 1)' at line 1
> exec ^
: string
| SELECT reads(YEAR FROM 1)
< error
| 1064 42000 …near 'reads(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT reads(DAY, 1, 2)
< error
| 1064 42000 …near 'reads(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT reads(a)
< error
| 1064 42000 …near 'reads(a)' at line 1
> exec ^
: string
| SELECT reads(a, b)
< error
| 1064 42000 …near 'reads(a, b)' at line 1
> exec ^
: string
| SELECT reads(a) OVER ()
< error
| 1064 42000 …near 'reads(a) OVER ()' at line 1
> exec ^
: string
| SELECT real()
< error
| 1064 42000 …near 'real()' at line 1
> exec ^
: string
| SELECT real(1)
< error
| 1064 42000 …near 'real(1)' at line 1
> exec ^
: string
| SELECT real(1,2)
< error
| 1064 42000 …near 'real(1,2)' at line 1
> exec ^
: string
| SELECT real(1,2,3)
< error
| 1064 42000 …near 'real(1,2,3)' at line 1
> exec ^
: string
| SELECT real(1,2,3,4)
< error
| 1064 42000 …near 'real(1,2,3,4)' at line 1
> exec ^
: string
| SELECT real(*)
< error
| 1064 42000 …near 'real(*)' at line 1
> exec ^
: string
| SELECT real('a')
< error
| 1064 42000 …near 'real('a')' at line 1
> exec ^
: string
| SELECT real(a USING utf8mb4)
< error
| 1064 42000 …near 'real(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT real(1 AS SIGNED)
< error
| 1064 42000 …near 'real(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT real(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'real(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT real(1 FROM 2)
< error
| 1064 42000 …near 'real(1 FROM 2)' at line 1
> exec ^
: string
| SELECT real(1 IN 2)
< error
| 1064 42000 …near 'real(1 IN 2)' at line 1
> exec ^
: string
| SELECT real(DISTINCT 1)
< error
| 1064 42000 …near 'real(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT real() OVER ()
< error
| 1064 42000 …near 'real() OVER ()' at line 1
> exec ^
: string
| SELECT real(1) OVER ()
< error
| 1064 42000 …near 'real(1) OVER ()' at line 1
> exec ^
: string
| SELECT real(1,2) OVER ()
< error
| 1064 42000 …near 'real(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT real(1,2,3) OVER ()
< error
| 1064 42000 …near 'real(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT real(1) FROM t
< error
| 1064 42000 …near 'real(1) FROM t' at line 1
> exec ^
: string
| SELECT real(SELECT 1)
< error
| 1064 42000 …near 'real(SELECT 1)' at line 1
> exec ^
: string
| SELECT real(YEAR FROM 1)
< error
| 1064 42000 …near 'real(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT real(DAY, 1, 2)
< error
| 1064 42000 …near 'real(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT real(a)
< error
| 1064 42000 …near 'real(a)' at line 1
> exec ^
: string
| SELECT real(a, b)
< error
| 1064 42000 …near 'real(a, b)' at line 1
> exec ^
: string
| SELECT real(a) OVER ()
< error
| 1064 42000 …near 'real(a) OVER ()' at line 1
> exec ^
: string
| SELECT rebuild()
< error
| 1630 42000 FUNCTION test.rebuild does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT rebuild()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT recover()
< error
| 1630 42000 FUNCTION test.recover does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT recover()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT recursive()
< error
| 1064 42000 …near 'recursive()' at line 1
> exec ^
: string
| SELECT recursive(1)
< error
| 1064 42000 …near 'recursive(1)' at line 1
> exec ^
: string
| SELECT recursive(1,2)
< error
| 1064 42000 …near 'recursive(1,2)' at line 1
> exec ^
: string
| SELECT recursive(1,2,3)
< error
| 1064 42000 …near 'recursive(1,2,3)' at line 1
> exec ^
: string
| SELECT recursive(1,2,3,4)
< error
| 1064 42000 …near 'recursive(1,2,3,4)' at line 1
> exec ^
: string
| SELECT recursive(*)
< error
| 1064 42000 …near 'recursive(*)' at line 1
> exec ^
: string
| SELECT recursive('a')
< error
| 1064 42000 …near 'recursive('a')' at line 1
> exec ^
: string
| SELECT recursive(a USING utf8mb4)
< error
| 1064 42000 …near 'recursive(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT recursive(1 AS SIGNED)
< error
| 1064 42000 …near 'recursive(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT recursive(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'recursive(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT recursive(1 FROM 2)
< error
| 1064 42000 …near 'recursive(1 FROM 2)' at line 1
> exec ^
: string
| SELECT recursive(1 IN 2)
< error
| 1064 42000 …near 'recursive(1 IN 2)' at line 1
> exec ^
: string
| SELECT recursive(DISTINCT 1)
< error
| 1064 42000 …near 'recursive(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT recursive() OVER ()
< error
| 1064 42000 …near 'recursive() OVER ()' at line 1
> exec ^
: string
| SELECT recursive(1) OVER ()
< error
| 1064 42000 …near 'recursive(1) OVER ()' at line 1
> exec ^
: string
| SELECT recursive(1,2) OVER ()
< error
| 1064 42000 …near 'recursive(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT recursive(1,2,3) OVER ()
< error
| 1064 42000 …near 'recursive(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT recursive(1) FROM t
< error
| 1064 42000 …near 'recursive(1) FROM t' at line 1
> exec ^
: string
| SELECT recursive(SELECT 1)
< error
| 1064 42000 …near 'recursive(SELECT 1)' at line 1
> exec ^
: string
| SELECT recursive(YEAR FROM 1)
< error
| 1064 42000 …near 'recursive(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT recursive(DAY, 1, 2)
< error
| 1064 42000 …near 'recursive(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT recursive(a)
< error
| 1064 42000 …near 'recursive(a)' at line 1
> exec ^
: string
| SELECT recursive(a, b)
< error
| 1064 42000 …near 'recursive(a, b)' at line 1
> exec ^
: string
| SELECT recursive(a) OVER ()
< error
| 1064 42000 …near 'recursive(a) OVER ()' at line 1
> exec ^
: string
| SELECT redo_buffer_size()
< error
| 1630 42000 FUNCTION test.redo_buffer_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT redo_buffer_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT redundant()
< error
| 1630 42000 FUNCTION test.redundant does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT redundant()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT reference()
< error
| 1305 42000 FUNCTION test.reference does not exist
> exec ^
: string
| SELECT reference()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT references()
< error
| 1064 42000 …near 'references()' at line 1
> exec ^
: string
| SELECT references(1)
< error
| 1064 42000 …near 'references(1)' at line 1
> exec ^
: string
| SELECT references(1,2)
< error
| 1064 42000 …near 'references(1,2)' at line 1
> exec ^
: string
| SELECT references(1,2,3)
< error
| 1064 42000 …near 'references(1,2,3)' at line 1
> exec ^
: string
| SELECT references(1,2,3,4)
< error
| 1064 42000 …near 'references(1,2,3,4)' at line 1
> exec ^
: string
| SELECT references(*)
< error
| 1064 42000 …near 'references(*)' at line 1
> exec ^
: string
| SELECT references('a')
< error
| 1064 42000 …near 'references('a')' at line 1
> exec ^
: string
| SELECT references(a USING utf8mb4)
< error
| 1064 42000 …near 'references(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT references(1 AS SIGNED)
< error
| 1064 42000 …near 'references(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT references(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'references(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT references(1 FROM 2)
< error
| 1064 42000 …near 'references(1 FROM 2)' at line 1
> exec ^
: string
| SELECT references(1 IN 2)
< error
| 1064 42000 …near 'references(1 IN 2)' at line 1
> exec ^
: string
| SELECT references(DISTINCT 1)
< error
| 1064 42000 …near 'references(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT references() OVER ()
< error
| 1064 42000 …near 'references() OVER ()' at line 1
> exec ^
: string
| SELECT references(1) OVER ()
< error
| 1064 42000 …near 'references(1) OVER ()' at line 1
> exec ^
: string
| SELECT references(1,2) OVER ()
< error
| 1064 42000 …near 'references(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT references(1,2,3) OVER ()
< error
| 1064 42000 …near 'references(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT references(1) FROM t
< error
| 1064 42000 …near 'references(1) FROM t' at line 1
> exec ^
: string
| SELECT references(SELECT 1)
< error
| 1064 42000 …near 'references(SELECT 1)' at line 1
> exec ^
: string
| SELECT references(YEAR FROM 1)
< error
| 1064 42000 …near 'references(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT references(DAY, 1, 2)
< error
| 1064 42000 …near 'references(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT references(a)
< error
| 1064 42000 …near 'references(a)' at line 1
> exec ^
: string
| SELECT references(a, b)
< error
| 1064 42000 …near 'references(a, b)' at line 1
> exec ^
: string
| SELECT references(a) OVER ()
< error
| 1064 42000 …near 'references(a) OVER ()' at line 1
> exec ^
: string
| SELECT regexp()
< error
| 1064 42000 …near 'regexp()' at line 1
> exec ^
: string
| SELECT regexp(1)
< error
| 1064 42000 …near 'regexp(1)' at line 1
> exec ^
: string
| SELECT regexp(1,2)
< error
| 1064 42000 …near 'regexp(1,2)' at line 1
> exec ^
: string
| SELECT regexp(1,2,3)
< error
| 1064 42000 …near 'regexp(1,2,3)' at line 1
> exec ^
: string
| SELECT regexp(1,2,3,4)
< error
| 1064 42000 …near 'regexp(1,2,3,4)' at line 1
> exec ^
: string
| SELECT regexp(*)
< error
| 1064 42000 …near 'regexp(*)' at line 1
> exec ^
: string
| SELECT regexp('a')
< error
| 1064 42000 …near 'regexp('a')' at line 1
> exec ^
: string
| SELECT regexp(a USING utf8mb4)
< error
| 1064 42000 …near 'regexp(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT regexp(1 AS SIGNED)
< error
| 1064 42000 …near 'regexp(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT regexp(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'regexp(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT regexp(1 FROM 2)
< error
| 1064 42000 …near 'regexp(1 FROM 2)' at line 1
> exec ^
: string
| SELECT regexp(1 IN 2)
< error
| 1064 42000 …near 'regexp(1 IN 2)' at line 1
> exec ^
: string
| SELECT regexp(DISTINCT 1)
< error
| 1064 42000 …near 'regexp(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT regexp() OVER ()
< error
| 1064 42000 …near 'regexp() OVER ()' at line 1
> exec ^
: string
| SELECT regexp(1) OVER ()
< error
| 1064 42000 …near 'regexp(1) OVER ()' at line 1
> exec ^
: string
| SELECT regexp(1,2) OVER ()
< error
| 1064 42000 …near 'regexp(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT regexp(1,2,3) OVER ()
< error
| 1064 42000 …near 'regexp(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT regexp(1) FROM t
< error
| 1064 42000 …near 'regexp(1) FROM t' at line 1
> exec ^
: string
| SELECT regexp(SELECT 1)
< error
| 1064 42000 …near 'regexp(SELECT 1)' at line 1
> exec ^
: string
| SELECT regexp(YEAR FROM 1)
< error
| 1064 42000 …near 'regexp(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT regexp(DAY, 1, 2)
< error
| 1064 42000 …near 'regexp(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT regexp(a)
< error
| 1064 42000 …near 'regexp(a)' at line 1
> exec ^
: string
| SELECT regexp(a, b)
< error
| 1064 42000 …near 'regexp(a, b)' at line 1
> exec ^
: string
| SELECT regexp(a) OVER ()
< error
| 1064 42000 …near 'regexp(a) OVER ()' at line 1
> exec ^
: string
| SELECT registration()
< error
| 1305 42000 FUNCTION test.registration does not exist
> exec ^
: string
| SELECT registration()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relational()
< error
| 1305 42000 FUNCTION test.relational does not exist
> exec ^
: string
| SELECT relational()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relay()
< error
| 1630 42000 FUNCTION test.relay does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT relay()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relay_log_file()
< error
| 1630 42000 FUNCTION test.relay_log_file does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT relay_log_file()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relay_log_pos()
< error
| 1630 42000 FUNCTION test.relay_log_pos does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT relay_log_pos()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relay_thread()
< error
| 1630 42000 FUNCTION test.relay_thread does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT relay_thread()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT relaylog()
< error
| 1630 42000 FUNCTION test.relaylog does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT relaylog()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT release()
< error
| 1064 42000 …near 'release()' at line 1
> exec ^
: string
| SELECT release(1)
< error
| 1064 42000 …near 'release(1)' at line 1
> exec ^
: string
| SELECT release(1,2)
< error
| 1064 42000 …near 'release(1,2)' at line 1
> exec ^
: string
| SELECT release(1,2,3)
< error
| 1064 42000 …near 'release(1,2,3)' at line 1
> exec ^
: string
| SELECT release(1,2,3,4)
< error
| 1064 42000 …near 'release(1,2,3,4)' at line 1
> exec ^
: string
| SELECT release(*)
< error
| 1064 42000 …near 'release(*)' at line 1
> exec ^
: string
| SELECT release('a')
< error
| 1064 42000 …near 'release('a')' at line 1
> exec ^
: string
| SELECT release(a USING utf8mb4)
< error
| 1064 42000 …near 'release(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT release(1 AS SIGNED)
< error
| 1064 42000 …near 'release(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT release(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'release(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT release(1 FROM 2)
< error
| 1064 42000 …near 'release(1 FROM 2)' at line 1
> exec ^
: string
| SELECT release(1 IN 2)
< error
| 1064 42000 …near 'release(1 IN 2)' at line 1
> exec ^
: string
| SELECT release(DISTINCT 1)
< error
| 1064 42000 …near 'release(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT release() OVER ()
< error
| 1064 42000 …near 'release() OVER ()' at line 1
> exec ^
: string
| SELECT release(1) OVER ()
< error
| 1064 42000 …near 'release(1) OVER ()' at line 1
> exec ^
: string
| SELECT release(1,2) OVER ()
< error
| 1064 42000 …near 'release(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT release(1,2,3) OVER ()
< error
| 1064 42000 …near 'release(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT release(1) FROM t
< error
| 1064 42000 …near 'release(1) FROM t' at line 1
> exec ^
: string
| SELECT release(SELECT 1)
< error
| 1064 42000 …near 'release(SELECT 1)' at line 1
> exec ^
: string
| SELECT release(YEAR FROM 1)
< error
| 1064 42000 …near 'release(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT release(DAY, 1, 2)
< error
| 1064 42000 …near 'release(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT release(a)
< error
| 1064 42000 …near 'release(a)' at line 1
> exec ^
: string
| SELECT release(a, b)
< error
| 1064 42000 …near 'release(a, b)' at line 1
> exec ^
: string
| SELECT release(a) OVER ()
< error
| 1064 42000 …near 'release(a) OVER ()' at line 1
> exec ^
: string
| SELECT reload()
< error
| 1630 42000 FUNCTION test.reload does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT reload()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT remove()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT remove(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT remove(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT remove(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT remove(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT remove(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT remove('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT remove(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT remove(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT remove(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT remove(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT remove(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT remove(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT remove() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT remove(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT remove(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT remove(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT remove(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT remove(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT remove(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT remove(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT remove(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT remove(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT remove(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT rename()
< error
| 1064 42000 …near 'rename()' at line 1
> exec ^
: string
| SELECT rename(1)
< error
| 1064 42000 …near 'rename(1)' at line 1
> exec ^
: string
| SELECT rename(1,2)
< error
| 1064 42000 …near 'rename(1,2)' at line 1
> exec ^
: string
| SELECT rename(1,2,3)
< error
| 1064 42000 …near 'rename(1,2,3)' at line 1
> exec ^
: string
| SELECT rename(1,2,3,4)
< error
| 1064 42000 …near 'rename(1,2,3,4)' at line 1
> exec ^
: string
| SELECT rename(*)
< error
| 1064 42000 …near 'rename(*)' at line 1
> exec ^
: string
| SELECT rename('a')
< error
| 1064 42000 …near 'rename('a')' at line 1
> exec ^
: string
| SELECT rename(a USING utf8mb4)
< error
| 1064 42000 …near 'rename(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT rename(1 AS SIGNED)
< error
| 1064 42000 …near 'rename(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT rename(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'rename(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT rename(1 FROM 2)
< error
| 1064 42000 …near 'rename(1 FROM 2)' at line 1
> exec ^
: string
| SELECT rename(1 IN 2)
< error
| 1064 42000 …near 'rename(1 IN 2)' at line 1
> exec ^
: string
| SELECT rename(DISTINCT 1)
< error
| 1064 42000 …near 'rename(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT rename() OVER ()
< error
| 1064 42000 …near 'rename() OVER ()' at line 1
> exec ^
: string
| SELECT rename(1) OVER ()
< error
| 1064 42000 …near 'rename(1) OVER ()' at line 1
> exec ^
: string
| SELECT rename(1,2) OVER ()
< error
| 1064 42000 …near 'rename(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT rename(1,2,3) OVER ()
< error
| 1064 42000 …near 'rename(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT rename(1) FROM t
< error
| 1064 42000 …near 'rename(1) FROM t' at line 1
> exec ^
: string
| SELECT rename(SELECT 1)
< error
| 1064 42000 …near 'rename(SELECT 1)' at line 1
> exec ^
: string
| SELECT rename(YEAR FROM 1)
< error
| 1064 42000 …near 'rename(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT rename(DAY, 1, 2)
< error
| 1064 42000 …near 'rename(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT rename(a)
< error
| 1064 42000 …near 'rename(a)' at line 1
> exec ^
: string
| SELECT rename(a, b)
< error
| 1064 42000 …near 'rename(a, b)' at line 1
> exec ^
: string
| SELECT rename(a) OVER ()
< error
| 1064 42000 …near 'rename(a) OVER ()' at line 1
> exec ^
: string
| SELECT reorganize()
< error
| 1630 42000 FUNCTION test.reorganize does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT reorganize()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT repair()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT repair(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT repair(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT repair(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT repair(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT repair(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT repair('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT repair(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT repair(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT repair(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT repair(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT repair(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT repair(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT repair() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT repair(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT repair(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT repair(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT repair(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT repair(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT repair(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT repair(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT repair(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT repair(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT repair(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT repeat()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT repeat(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT repeat(1,2)
< ok 0 0
> exec ^
: string
| SELECT repeatable()
< error
| 1630 42000 FUNCTION test.repeatable does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT repeatable()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replace()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT replace(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT replace(1,2)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT replace(1,2,3)
< ok 0 0
> exec ^
: string
| SELECT replica()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT replica(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT replica(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT replica(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT replica(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT replica(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT replica('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT replica(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT replica(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT replica(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT replica(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT replica(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT replica(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT replica() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT replica(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT replica(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT replica(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT replica(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT replica(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT replica(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT replica(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT replica(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT replica(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT replica(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT replicas()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT replicas(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT replicas(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT replicas(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT replicas(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT replicas(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT replicas('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT replicas(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT replicas(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT replicas(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT replicas(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT replicas(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT replicas(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT replicas() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT replicas(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT replicas(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT replicas(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT replicas(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT replicas(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT replicas(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT replicas(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT replicas(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT replicas(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT replicas(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT replicate_do_db()
< error
| 1305 42000 FUNCTION test.replicate_do_db does not exist
> exec ^
: string
| SELECT replicate_do_db()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_do_table()
< error
| 1305 42000 FUNCTION test.replicate_do_table does not exist
> exec ^
: string
| SELECT replicate_do_table()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_ignore_db()
< error
| 1305 42000 FUNCTION test.replicate_ignore_db does not exist
> exec ^
: string
| SELECT replicate_ignore_db()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_ignore_table()
< error
| 1305 42000 FUNCTION test.replicate_ignore_table does not exist
> exec ^
: string
| SELECT replicate_ignore_table()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_rewrite_db()
< error
| 1305 42000 FUNCTION test.replicate_rewrite_db does not exist
> exec ^
: string
| SELECT replicate_rewrite_db()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_wild_do_table()
< error
| 1305 42000 FUNCTION test.replicate_wild_do_table does not exist
> exec ^
: string
| SELECT replicate_wild_do_table()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replicate_wild_ignore_table()
< error
| 1305 42000 FUNCTION test.replicate_wild_ignore_table does not exist
> exec ^
: string
| SELECT replicate_wild_ignore_table()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT replication()
< error
| 1630 42000 FUNCTION test.replication does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT replication()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT require()
< error
| 1064 42000 …near 'require()' at line 1
> exec ^
: string
| SELECT require(1)
< error
| 1064 42000 …near 'require(1)' at line 1
> exec ^
: string
| SELECT require(1,2)
< error
| 1064 42000 …near 'require(1,2)' at line 1
> exec ^
: string
| SELECT require(1,2,3)
< error
| 1064 42000 …near 'require(1,2,3)' at line 1
> exec ^
: string
| SELECT require(1,2,3,4)
< error
| 1064 42000 …near 'require(1,2,3,4)' at line 1
> exec ^
: string
| SELECT require(*)
< error
| 1064 42000 …near 'require(*)' at line 1
> exec ^
: string
| SELECT require('a')
< error
| 1064 42000 …near 'require('a')' at line 1
> exec ^
: string
| SELECT require(a USING utf8mb4)
< error
| 1064 42000 …near 'require(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT require(1 AS SIGNED)
< error
| 1064 42000 …near 'require(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT require(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'require(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT require(1 FROM 2)
< error
| 1064 42000 …near 'require(1 FROM 2)' at line 1
> exec ^
: string
| SELECT require(1 IN 2)
< error
| 1064 42000 …near 'require(1 IN 2)' at line 1
> exec ^
: string
| SELECT require(DISTINCT 1)
< error
| 1064 42000 …near 'require(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT require() OVER ()
< error
| 1064 42000 …near 'require() OVER ()' at line 1
> exec ^
: string
| SELECT require(1) OVER ()
< error
| 1064 42000 …near 'require(1) OVER ()' at line 1
> exec ^
: string
| SELECT require(1,2) OVER ()
< error
| 1064 42000 …near 'require(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT require(1,2,3) OVER ()
< error
| 1064 42000 …near 'require(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT require(1) FROM t
< error
| 1064 42000 …near 'require(1) FROM t' at line 1
> exec ^
: string
| SELECT require(SELECT 1)
< error
| 1064 42000 …near 'require(SELECT 1)' at line 1
> exec ^
: string
| SELECT require(YEAR FROM 1)
< error
| 1064 42000 …near 'require(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT require(DAY, 1, 2)
< error
| 1064 42000 …near 'require(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT require(a)
< error
| 1064 42000 …near 'require(a)' at line 1
> exec ^
: string
| SELECT require(a, b)
< error
| 1064 42000 …near 'require(a, b)' at line 1
> exec ^
: string
| SELECT require(a) OVER ()
< error
| 1064 42000 …near 'require(a) OVER ()' at line 1
> exec ^
: string
| SELECT require_row_format()
< error
| 1305 42000 FUNCTION test.require_row_format does not exist
> exec ^
: string
| SELECT require_row_format()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT require_table_primary_key_check()
< error
| 1305 42000 FUNCTION test.require_table_primary_key_check does not exist
> exec ^
: string
| SELECT require_table_primary_key_check()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT reset()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT reset(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT reset(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT reset(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT reset(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT reset(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT reset('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT reset(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT reset(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT reset(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT reset(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT reset(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT reset(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT reset() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT reset(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT reset(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT reset(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT reset(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT reset(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT reset(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT reset(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT reset(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT reset(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT reset(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT resignal()
< error
| 1064 42000 …near 'resignal()' at line 1
> exec ^
: string
| SELECT resignal(1)
< error
| 1064 42000 …near 'resignal(1)' at line 1
> exec ^
: string
| SELECT resignal(1,2)
< error
| 1064 42000 …near 'resignal(1,2)' at line 1
> exec ^
: string
| SELECT resignal(1,2,3)
< error
| 1064 42000 …near 'resignal(1,2,3)' at line 1
> exec ^
: string
| SELECT resignal(1,2,3,4)
< error
| 1064 42000 …near 'resignal(1,2,3,4)' at line 1
> exec ^
: string
| SELECT resignal(*)
< error
| 1064 42000 …near 'resignal(*)' at line 1
> exec ^
: string
| SELECT resignal('a')
< error
| 1064 42000 …near 'resignal('a')' at line 1
> exec ^
: string
| SELECT resignal(a USING utf8mb4)
< error
| 1064 42000 …near 'resignal(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT resignal(1 AS SIGNED)
< error
| 1064 42000 …near 'resignal(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT resignal(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'resignal(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT resignal(1 FROM 2)
< error
| 1064 42000 …near 'resignal(1 FROM 2)' at line 1
> exec ^
: string
| SELECT resignal(1 IN 2)
< error
| 1064 42000 …near 'resignal(1 IN 2)' at line 1
> exec ^
: string
| SELECT resignal(DISTINCT 1)
< error
| 1064 42000 …near 'resignal(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT resignal() OVER ()
< error
| 1064 42000 …near 'resignal() OVER ()' at line 1
> exec ^
: string
| SELECT resignal(1) OVER ()
< error
| 1064 42000 …near 'resignal(1) OVER ()' at line 1
> exec ^
: string
| SELECT resignal(1,2) OVER ()
< error
| 1064 42000 …near 'resignal(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT resignal(1,2,3) OVER ()
< error
| 1064 42000 …near 'resignal(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT resignal(1) FROM t
< error
| 1064 42000 …near 'resignal(1) FROM t' at line 1
> exec ^
: string
| SELECT resignal(SELECT 1)
< error
| 1064 42000 …near 'resignal(SELECT 1)' at line 1
> exec ^
: string
| SELECT resignal(YEAR FROM 1)
< error
| 1064 42000 …near 'resignal(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT resignal(DAY, 1, 2)
< error
| 1064 42000 …near 'resignal(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT resignal(a)
< error
| 1064 42000 …near 'resignal(a)' at line 1
> exec ^
: string
| SELECT resignal(a, b)
< error
| 1064 42000 …near 'resignal(a, b)' at line 1
> exec ^
: string
| SELECT resignal(a) OVER ()
< error
| 1064 42000 …near 'resignal(a) OVER ()' at line 1
> exec ^
: string
| SELECT resource()
< error
| 1305 42000 FUNCTION test.resource does not exist
> exec ^
: string
| SELECT resource()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT respect()
< error
| 1305 42000 FUNCTION test.respect does not exist
> exec ^
: string
| SELECT respect()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT restart()
< error
| 1630 42000 FUNCTION test.restart does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT restart()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT restore()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT restore(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT restore(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT restore(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT restore(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT restore(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT restore('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT restore(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT restore(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT restore(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT restore(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT restore(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT restore(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT restore() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT restore(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT restore(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT restore(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT restore(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT restore(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT restore(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT restore(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT restore(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT restore(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT restore(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT restrict()
< error
| 1064 42000 …near 'restrict()' at line 1
> exec ^
: string
| SELECT restrict(1)
< error
| 1064 42000 …near 'restrict(1)' at line 1
> exec ^
: string
| SELECT restrict(1,2)
< error
| 1064 42000 …near 'restrict(1,2)' at line 1
> exec ^
: string
| SELECT restrict(1,2,3)
< error
| 1064 42000 …near 'restrict(1,2,3)' at line 1
> exec ^
: string
| SELECT restrict(1,2,3,4)
< error
| 1064 42000 …near 'restrict(1,2,3,4)' at line 1
> exec ^
: string
| SELECT restrict(*)
< error
| 1064 42000 …near 'restrict(*)' at line 1
> exec ^
: string
| SELECT restrict('a')
< error
| 1064 42000 …near 'restrict('a')' at line 1
> exec ^
: string
| SELECT restrict(a USING utf8mb4)
< error
| 1064 42000 …near 'restrict(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT restrict(1 AS SIGNED)
< error
| 1064 42000 …near 'restrict(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT restrict(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'restrict(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT restrict(1 FROM 2)
< error
| 1064 42000 …near 'restrict(1 FROM 2)' at line 1
> exec ^
: string
| SELECT restrict(1 IN 2)
< error
| 1064 42000 …near 'restrict(1 IN 2)' at line 1
> exec ^
: string
| SELECT restrict(DISTINCT 1)
< error
| 1064 42000 …near 'restrict(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT restrict() OVER ()
< error
| 1064 42000 …near 'restrict() OVER ()' at line 1
> exec ^
: string
| SELECT restrict(1) OVER ()
< error
| 1064 42000 …near 'restrict(1) OVER ()' at line 1
> exec ^
: string
| SELECT restrict(1,2) OVER ()
< error
| 1064 42000 …near 'restrict(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT restrict(1,2,3) OVER ()
< error
| 1064 42000 …near 'restrict(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT restrict(1) FROM t
< error
| 1064 42000 …near 'restrict(1) FROM t' at line 1
> exec ^
: string
| SELECT restrict(SELECT 1)
< error
| 1064 42000 …near 'restrict(SELECT 1)' at line 1
> exec ^
: string
| SELECT restrict(YEAR FROM 1)
< error
| 1064 42000 …near 'restrict(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT restrict(DAY, 1, 2)
< error
| 1064 42000 …near 'restrict(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT restrict(a)
< error
| 1064 42000 …near 'restrict(a)' at line 1
> exec ^
: string
| SELECT restrict(a, b)
< error
| 1064 42000 …near 'restrict(a, b)' at line 1
> exec ^
: string
| SELECT restrict(a) OVER ()
< error
| 1064 42000 …near 'restrict(a) OVER ()' at line 1
> exec ^
: string
| SELECT resume()
< error
| 1630 42000 FUNCTION test.resume does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT resume()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT retain()
< error
| 1305 42000 FUNCTION test.retain does not exist
> exec ^
: string
| SELECT retain()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT return()
< error
| 1064 42000 …near 'return()' at line 1
> exec ^
: string
| SELECT return(1)
< error
| 1064 42000 …near 'return(1)' at line 1
> exec ^
: string
| SELECT return(1,2)
< error
| 1064 42000 …near 'return(1,2)' at line 1
> exec ^
: string
| SELECT return(1,2,3)
< error
| 1064 42000 …near 'return(1,2,3)' at line 1
> exec ^
: string
| SELECT return(1,2,3,4)
< error
| 1064 42000 …near 'return(1,2,3,4)' at line 1
> exec ^
: string
| SELECT return(*)
< error
| 1064 42000 …near 'return(*)' at line 1
> exec ^
: string
| SELECT return('a')
< error
| 1064 42000 …near 'return('a')' at line 1
> exec ^
: string
| SELECT return(a USING utf8mb4)
< error
| 1064 42000 …near 'return(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT return(1 AS SIGNED)
< error
| 1064 42000 …near 'return(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT return(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'return(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT return(1 FROM 2)
< error
| 1064 42000 …near 'return(1 FROM 2)' at line 1
> exec ^
: string
| SELECT return(1 IN 2)
< error
| 1064 42000 …near 'return(1 IN 2)' at line 1
> exec ^
: string
| SELECT return(DISTINCT 1)
< error
| 1064 42000 …near 'return(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT return() OVER ()
< error
| 1064 42000 …near 'return() OVER ()' at line 1
> exec ^
: string
| SELECT return(1) OVER ()
< error
| 1064 42000 …near 'return(1) OVER ()' at line 1
> exec ^
: string
| SELECT return(1,2) OVER ()
< error
| 1064 42000 …near 'return(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT return(1,2,3) OVER ()
< error
| 1064 42000 …near 'return(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT return(1) FROM t
< error
| 1064 42000 …near 'return(1) FROM t' at line 1
> exec ^
: string
| SELECT return(SELECT 1)
< error
| 1064 42000 …near 'return(SELECT 1)' at line 1
> exec ^
: string
| SELECT return(YEAR FROM 1)
< error
| 1064 42000 …near 'return(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT return(DAY, 1, 2)
< error
| 1064 42000 …near 'return(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT return(a)
< error
| 1064 42000 …near 'return(a)' at line 1
> exec ^
: string
| SELECT return(a, b)
< error
| 1064 42000 …near 'return(a, b)' at line 1
> exec ^
: string
| SELECT return(a) OVER ()
< error
| 1064 42000 …near 'return(a) OVER ()' at line 1
> exec ^
: string
| SELECT returned_sqlstate()
< error
| 1630 42000 FUNCTION test.returned_sqlstate does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT returned_sqlstate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT returning()
< error
| 1064 42000 …near 'returning()' at line 1
> exec ^
: string
| SELECT returning(1)
< error
| 1064 42000 …near 'returning(1)' at line 1
> exec ^
: string
| SELECT returning(1,2)
< error
| 1064 42000 …near 'returning(1,2)' at line 1
> exec ^
: string
| SELECT returning(1,2,3)
< error
| 1064 42000 …near 'returning(1,2,3)' at line 1
> exec ^
: string
| SELECT returning(1,2,3,4)
< error
| 1064 42000 …near 'returning(1,2,3,4)' at line 1
> exec ^
: string
| SELECT returning(*)
< error
| 1064 42000 …near 'returning(*)' at line 1
> exec ^
: string
| SELECT returning('a')
< error
| 1064 42000 …near 'returning('a')' at line 1
> exec ^
: string
| SELECT returning(a USING utf8mb4)
< error
| 1064 42000 …near 'returning(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT returning(1 AS SIGNED)
< error
| 1064 42000 …near 'returning(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT returning(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'returning(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT returning(1 FROM 2)
< error
| 1064 42000 …near 'returning(1 FROM 2)' at line 1
> exec ^
: string
| SELECT returning(1 IN 2)
< error
| 1064 42000 …near 'returning(1 IN 2)' at line 1
> exec ^
: string
| SELECT returning(DISTINCT 1)
< error
| 1064 42000 …near 'returning(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT returning() OVER ()
< error
| 1064 42000 …near 'returning() OVER ()' at line 1
> exec ^
: string
| SELECT returning(1) OVER ()
< error
| 1064 42000 …near 'returning(1) OVER ()' at line 1
> exec ^
: string
| SELECT returning(1,2) OVER ()
< error
| 1064 42000 …near 'returning(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT returning(1,2,3) OVER ()
< error
| 1064 42000 …near 'returning(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT returning(1) FROM t
< error
| 1064 42000 …near 'returning(1) FROM t' at line 1
> exec ^
: string
| SELECT returning(SELECT 1)
< error
| 1064 42000 …near 'returning(SELECT 1)' at line 1
> exec ^
: string
| SELECT returning(YEAR FROM 1)
< error
| 1064 42000 …near 'returning(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT returning(DAY, 1, 2)
< error
| 1064 42000 …near 'returning(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT returning(a)
< error
| 1064 42000 …near 'returning(a)' at line 1
> exec ^
: string
| SELECT returning(a, b)
< error
| 1064 42000 …near 'returning(a, b)' at line 1
> exec ^
: string
| SELECT returning(a) OVER ()
< error
| 1064 42000 …near 'returning(a) OVER ()' at line 1
> exec ^
: string
| SELECT returns()
< error
| 1630 42000 FUNCTION test.returns does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT returns()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT reuse()
< error
| 1630 42000 FUNCTION test.reuse does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT reuse()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT reverse()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'reverse'
> exec ^
: string
| SELECT reverse()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'reverse'
> exec ^
: string
| SELECT reverse(1)
< ok 0 0
> exec ^
: string
| SELECT revoke()
< error
| 1064 42000 …near 'revoke()' at line 1
> exec ^
: string
| SELECT revoke(1)
< error
| 1064 42000 …near 'revoke(1)' at line 1
> exec ^
: string
| SELECT revoke(1,2)
< error
| 1064 42000 …near 'revoke(1,2)' at line 1
> exec ^
: string
| SELECT revoke(1,2,3)
< error
| 1064 42000 …near 'revoke(1,2,3)' at line 1
> exec ^
: string
| SELECT revoke(1,2,3,4)
< error
| 1064 42000 …near 'revoke(1,2,3,4)' at line 1
> exec ^
: string
| SELECT revoke(*)
< error
| 1064 42000 …near 'revoke(*)' at line 1
> exec ^
: string
| SELECT revoke('a')
< error
| 1064 42000 …near 'revoke('a')' at line 1
> exec ^
: string
| SELECT revoke(a USING utf8mb4)
< error
| 1064 42000 …near 'revoke(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT revoke(1 AS SIGNED)
< error
| 1064 42000 …near 'revoke(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT revoke(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'revoke(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT revoke(1 FROM 2)
< error
| 1064 42000 …near 'revoke(1 FROM 2)' at line 1
> exec ^
: string
| SELECT revoke(1 IN 2)
< error
| 1064 42000 …near 'revoke(1 IN 2)' at line 1
> exec ^
: string
| SELECT revoke(DISTINCT 1)
< error
| 1064 42000 …near 'revoke(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT revoke() OVER ()
< error
| 1064 42000 …near 'revoke() OVER ()' at line 1
> exec ^
: string
| SELECT revoke(1) OVER ()
< error
| 1064 42000 …near 'revoke(1) OVER ()' at line 1
> exec ^
: string
| SELECT revoke(1,2) OVER ()
< error
| 1064 42000 …near 'revoke(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT revoke(1,2,3) OVER ()
< error
| 1064 42000 …near 'revoke(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT revoke(1) FROM t
< error
| 1064 42000 …near 'revoke(1) FROM t' at line 1
> exec ^
: string
| SELECT revoke(SELECT 1)
< error
| 1064 42000 …near 'revoke(SELECT 1)' at line 1
> exec ^
: string
| SELECT revoke(YEAR FROM 1)
< error
| 1064 42000 …near 'revoke(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT revoke(DAY, 1, 2)
< error
| 1064 42000 …near 'revoke(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT revoke(a)
< error
| 1064 42000 …near 'revoke(a)' at line 1
> exec ^
: string
| SELECT revoke(a, b)
< error
| 1064 42000 …near 'revoke(a, b)' at line 1
> exec ^
: string
| SELECT revoke(a) OVER ()
< error
| 1064 42000 …near 'revoke(a) OVER ()' at line 1
> exec ^
: string
| SELECT right()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT right(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT right(1,2)
< ok 0 0
> exec ^
: string
| SELECT rlike()
< error
| 1064 42000 …near 'rlike()' at line 1
> exec ^
: string
| SELECT rlike(1)
< error
| 1064 42000 …near 'rlike(1)' at line 1
> exec ^
: string
| SELECT rlike(1,2)
< error
| 1064 42000 …near 'rlike(1,2)' at line 1
> exec ^
: string
| SELECT rlike(1,2,3)
< error
| 1064 42000 …near 'rlike(1,2,3)' at line 1
> exec ^
: string
| SELECT rlike(1,2,3,4)
< error
| 1064 42000 …near 'rlike(1,2,3,4)' at line 1
> exec ^
: string
| SELECT rlike(*)
< error
| 1064 42000 …near 'rlike(*)' at line 1
> exec ^
: string
| SELECT rlike('a')
< error
| 1064 42000 …near 'rlike('a')' at line 1
> exec ^
: string
| SELECT rlike(a USING utf8mb4)
< error
| 1064 42000 …near 'rlike(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT rlike(1 AS SIGNED)
< error
| 1064 42000 …near 'rlike(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT rlike(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'rlike(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT rlike(1 FROM 2)
< error
| 1064 42000 …near 'rlike(1 FROM 2)' at line 1
> exec ^
: string
| SELECT rlike(1 IN 2)
< error
| 1064 42000 …near 'rlike(1 IN 2)' at line 1
> exec ^
: string
| SELECT rlike(DISTINCT 1)
< error
| 1064 42000 …near 'rlike(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT rlike() OVER ()
< error
| 1064 42000 …near 'rlike() OVER ()' at line 1
> exec ^
: string
| SELECT rlike(1) OVER ()
< error
| 1064 42000 …near 'rlike(1) OVER ()' at line 1
> exec ^
: string
| SELECT rlike(1,2) OVER ()
< error
| 1064 42000 …near 'rlike(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT rlike(1,2,3) OVER ()
< error
| 1064 42000 …near 'rlike(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT rlike(1) FROM t
< error
| 1064 42000 …near 'rlike(1) FROM t' at line 1
> exec ^
: string
| SELECT rlike(SELECT 1)
< error
| 1064 42000 …near 'rlike(SELECT 1)' at line 1
> exec ^
: string
| SELECT rlike(YEAR FROM 1)
< error
| 1064 42000 …near 'rlike(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT rlike(DAY, 1, 2)
< error
| 1064 42000 …near 'rlike(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT rlike(a)
< error
| 1064 42000 …near 'rlike(a)' at line 1
> exec ^
: string
| SELECT rlike(a, b)
< error
| 1064 42000 …near 'rlike(a, b)' at line 1
> exec ^
: string
| SELECT rlike(a) OVER ()
< error
| 1064 42000 …near 'rlike(a) OVER ()' at line 1
> exec ^
: string
| SELECT role()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT role(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT role(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT role(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT role(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT role(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT role('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT role(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT role(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT role(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT role(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT role(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT role(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT role() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT role(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT role(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT role(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT role(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT role(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT role(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT role(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT role(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT role(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT role(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT rollback()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT rollback(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT rollback(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT rollback(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT rollback(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT rollback(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT rollback('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT rollback(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT rollback(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT rollback(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT rollback(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT rollback(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT rollback(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT rollback() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT rollback(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT rollback(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT rollback(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT rollback(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT rollback(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT rollback(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT rollback(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT rollback(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT rollback(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT rollback(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT rollup()
< error
| 1630 42000 FUNCTION test.rollup does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT rollup()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT rotate()
< error
| 1305 42000 FUNCTION test.rotate does not exist
> exec ^
: string
| SELECT rotate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT routine()
< error
| 1630 42000 FUNCTION test.routine does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT routine()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT row()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT row(1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| SELECT row(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT row_count()
< ok 0 0
> exec ^
: string
| SELECT row_format()
< error
| 1630 42000 FUNCTION test.row_format does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT row_format()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT row_number()
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT row_number(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT row_number(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT row_number(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT row_number(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT row_number(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT row_number('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT row_number(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT row_number(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT row_number(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT row_number(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT row_number(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT row_number(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT row_number() OVER ()
< ok 0 0
> exec ^
: string
| SELECT rows()
< error
| 1064 42000 …near 'rows()' at line 1
> exec ^
: string
| SELECT rows(1)
< error
| 1064 42000 …near 'rows(1)' at line 1
> exec ^
: string
| SELECT rows(1,2)
< error
| 1064 42000 …near 'rows(1,2)' at line 1
> exec ^
: string
| SELECT rows(1,2,3)
< error
| 1064 42000 …near 'rows(1,2,3)' at line 1
> exec ^
: string
| SELECT rows(1,2,3,4)
< error
| 1064 42000 …near 'rows(1,2,3,4)' at line 1
> exec ^
: string
| SELECT rows(*)
< error
| 1064 42000 …near 'rows(*)' at line 1
> exec ^
: string
| SELECT rows('a')
< error
| 1064 42000 …near 'rows('a')' at line 1
> exec ^
: string
| SELECT rows(a USING utf8mb4)
< error
| 1064 42000 …near 'rows(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT rows(1 AS SIGNED)
< error
| 1064 42000 …near 'rows(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT rows(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'rows(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT rows(1 FROM 2)
< error
| 1064 42000 …near 'rows(1 FROM 2)' at line 1
> exec ^
: string
| SELECT rows(1 IN 2)
< error
| 1064 42000 …near 'rows(1 IN 2)' at line 1
> exec ^
: string
| SELECT rows(DISTINCT 1)
< error
| 1064 42000 …near 'rows(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT rows() OVER ()
< error
| 1064 42000 …near 'rows() OVER ()' at line 1
> exec ^
: string
| SELECT rows(1) OVER ()
< error
| 1064 42000 …near 'rows(1) OVER ()' at line 1
> exec ^
: string
| SELECT rows(1,2) OVER ()
< error
| 1064 42000 …near 'rows(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT rows(1,2,3) OVER ()
< error
| 1064 42000 …near 'rows(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT rows(1) FROM t
< error
| 1064 42000 …near 'rows(1) FROM t' at line 1
> exec ^
: string
| SELECT rows(SELECT 1)
< error
| 1064 42000 …near 'rows(SELECT 1)' at line 1
> exec ^
: string
| SELECT rows(YEAR FROM 1)
< error
| 1064 42000 …near 'rows(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT rows(DAY, 1, 2)
< error
| 1064 42000 …near 'rows(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT rows(a)
< error
| 1064 42000 …near 'rows(a)' at line 1
> exec ^
: string
| SELECT rows(a, b)
< error
| 1064 42000 …near 'rows(a, b)' at line 1
> exec ^
: string
| SELECT rows(a) OVER ()
< error
| 1064 42000 …near 'rows(a) OVER ()' at line 1
> exec ^
: string
| SELECT rtree()
< error
| 1630 42000 FUNCTION test.rtree does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT rtree()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT s3()
< error
| 1305 42000 FUNCTION test.s3 does not exist
> exec ^
: string
| SELECT s3()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT savepoint()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT savepoint(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT savepoint(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT savepoint(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT savepoint(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT savepoint(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT savepoint('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT savepoint(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT savepoint(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT savepoint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT savepoint(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT savepoint(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT savepoint(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT savepoint() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT savepoint(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT savepoint(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT savepoint(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT savepoint(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT savepoint(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT savepoint(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT savepoint(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT savepoint(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT savepoint(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT savepoint(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT schedule()
< error
| 1630 42000 FUNCTION test.schedule does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT schedule()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT schema()
< ok 0 0
> exec ^
: string
| SELECT schema_name()
< error
| 1630 42000 FUNCTION test.schema_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT schema_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT schemas()
< error
| 1064 42000 …near 'schemas()' at line 1
> exec ^
: string
| SELECT schemas(1)
< error
| 1064 42000 …near 'schemas(1)' at line 1
> exec ^
: string
| SELECT schemas(1,2)
< error
| 1064 42000 …near 'schemas(1,2)' at line 1
> exec ^
: string
| SELECT schemas(1,2,3)
< error
| 1064 42000 …near 'schemas(1,2,3)' at line 1
> exec ^
: string
| SELECT schemas(1,2,3,4)
< error
| 1064 42000 …near 'schemas(1,2,3,4)' at line 1
> exec ^
: string
| SELECT schemas(*)
< error
| 1064 42000 …near 'schemas(*)' at line 1
> exec ^
: string
| SELECT schemas('a')
< error
| 1064 42000 …near 'schemas('a')' at line 1
> exec ^
: string
| SELECT schemas(a USING utf8mb4)
< error
| 1064 42000 …near 'schemas(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT schemas(1 AS SIGNED)
< error
| 1064 42000 …near 'schemas(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT schemas(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'schemas(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT schemas(1 FROM 2)
< error
| 1064 42000 …near 'schemas(1 FROM 2)' at line 1
> exec ^
: string
| SELECT schemas(1 IN 2)
< error
| 1064 42000 …near 'schemas(1 IN 2)' at line 1
> exec ^
: string
| SELECT schemas(DISTINCT 1)
< error
| 1064 42000 …near 'schemas(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT schemas() OVER ()
< error
| 1064 42000 …near 'schemas() OVER ()' at line 1
> exec ^
: string
| SELECT schemas(1) OVER ()
< error
| 1064 42000 …near 'schemas(1) OVER ()' at line 1
> exec ^
: string
| SELECT schemas(1,2) OVER ()
< error
| 1064 42000 …near 'schemas(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT schemas(1,2,3) OVER ()
< error
| 1064 42000 …near 'schemas(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT schemas(1) FROM t
< error
| 1064 42000 …near 'schemas(1) FROM t' at line 1
> exec ^
: string
| SELECT schemas(SELECT 1)
< error
| 1064 42000 …near 'schemas(SELECT 1)' at line 1
> exec ^
: string
| SELECT schemas(YEAR FROM 1)
< error
| 1064 42000 …near 'schemas(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT schemas(DAY, 1, 2)
< error
| 1064 42000 …near 'schemas(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT schemas(a)
< error
| 1064 42000 …near 'schemas(a)' at line 1
> exec ^
: string
| SELECT schemas(a, b)
< error
| 1064 42000 …near 'schemas(a, b)' at line 1
> exec ^
: string
| SELECT schemas(a) OVER ()
< error
| 1064 42000 …near 'schemas(a) OVER ()' at line 1
> exec ^
: string
| SELECT second()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT second(1)
< ok 0 0
> exec ^
: string
| SELECT second_microsecond()
< error
| 1064 42000 …near 'second_microsecond()' at line 1
> exec ^
: string
| SELECT second_microsecond(1)
< error
| 1064 42000 …near 'second_microsecond(1)' at line 1
> exec ^
: string
| SELECT second_microsecond(1,2)
< error
| 1064 42000 …near 'second_microsecond(1,2)' at line 1
> exec ^
: string
| SELECT second_microsecond(1,2,3)
< error
| 1064 42000 …near 'second_microsecond(1,2,3)' at line 1
> exec ^
: string
| SELECT second_microsecond(1,2,3,4)
< error
| 1064 42000 …near 'second_microsecond(1,2,3,4)' at line 1
> exec ^
: string
| SELECT second_microsecond(*)
< error
| 1064 42000 …near 'second_microsecond(*)' at line 1
> exec ^
: string
| SELECT second_microsecond('a')
< error
| 1064 42000 …near 'second_microsecond('a')' at line 1
> exec ^
: string
| SELECT second_microsecond(a USING utf8mb4)
< error
| 1064 42000 …near 'second_microsecond(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT second_microsecond(1 AS SIGNED)
< error
| 1064 42000 …near 'second_microsecond(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT second_microsecond(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'second_microsecond(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT second_microsecond(1 FROM 2)
< error
| 1064 42000 …near 'second_microsecond(1 FROM 2)' at line 1
> exec ^
: string
| SELECT second_microsecond(1 IN 2)
< error
| 1064 42000 …near 'second_microsecond(1 IN 2)' at line 1
> exec ^
: string
| SELECT second_microsecond(DISTINCT 1)
< error
| 1064 42000 …near 'second_microsecond(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT second_microsecond() OVER ()
< error
| 1064 42000 …near 'second_microsecond() OVER ()' at line 1
> exec ^
: string
| SELECT second_microsecond(1) OVER ()
< error
| 1064 42000 …near 'second_microsecond(1) OVER ()' at line 1
> exec ^
: string
| SELECT second_microsecond(1,2) OVER ()
< error
| 1064 42000 …near 'second_microsecond(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT second_microsecond(1,2,3) OVER ()
< error
| 1064 42000 …near 'second_microsecond(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT second_microsecond(1) FROM t
< error
| 1064 42000 …near 'second_microsecond(1) FROM t' at line 1
> exec ^
: string
| SELECT second_microsecond(SELECT 1)
< error
| 1064 42000 …near 'second_microsecond(SELECT 1)' at line 1
> exec ^
: string
| SELECT second_microsecond(YEAR FROM 1)
< error
| 1064 42000 …near 'second_microsecond(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT second_microsecond(DAY, 1, 2)
< error
| 1064 42000 …near 'second_microsecond(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT second_microsecond(a)
< error
| 1064 42000 …near 'second_microsecond(a)' at line 1
> exec ^
: string
| SELECT second_microsecond(a, b)
< error
| 1064 42000 …near 'second_microsecond(a, b)' at line 1
> exec ^
: string
| SELECT second_microsecond(a) OVER ()
< error
| 1064 42000 …near 'second_microsecond(a) OVER ()' at line 1
> exec ^
: string
| SELECT secondary()
< error
| 1305 42000 FUNCTION test.secondary does not exist
> exec ^
: string
| SELECT secondary()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT secondary_engine()
< error
| 1305 42000 FUNCTION test.secondary_engine does not exist
> exec ^
: string
| SELECT secondary_engine()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT secondary_engine_attribute()
< error
| 1305 42000 FUNCTION test.secondary_engine_attribute does not exist
> exec ^
: string
| SELECT secondary_engine_attribute()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT secondary_load()
< error
| 1305 42000 FUNCTION test.secondary_load does not exist
> exec ^
: string
| SELECT secondary_load()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT secondary_unload()
< error
| 1305 42000 FUNCTION test.secondary_unload does not exist
> exec ^
: string
| SELECT secondary_unload()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT security()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT security(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT security(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT security(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT security(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT security(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT security('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT security(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT security(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT security(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT security(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT security(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT security(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT security() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT security(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT security(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT security(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT security(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT security(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT security(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT security(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT security(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT security(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT security(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT select()
< error
| 1064 42000 …near 'select()' at line 1
> exec ^
: string
| SELECT select(1)
< error
| 1064 42000 …near 'select(1)' at line 1
> exec ^
: string
| SELECT select(1,2)
< error
| 1064 42000 …near 'select(1,2)' at line 1
> exec ^
: string
| SELECT select(1,2,3)
< error
| 1064 42000 …near 'select(1,2,3)' at line 1
> exec ^
: string
| SELECT select(1,2,3,4)
< error
| 1064 42000 …near 'select(1,2,3,4)' at line 1
> exec ^
: string
| SELECT select(*)
< error
| 1064 42000 …near 'select(*)' at line 1
> exec ^
: string
| SELECT select('a')
< error
| 1064 42000 …near 'select('a')' at line 1
> exec ^
: string
| SELECT select(a USING utf8mb4)
< error
| 1064 42000 …near 'select(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT select(1 AS SIGNED)
< error
| 1064 42000 …near 'select(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT select(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'select(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT select(1 FROM 2)
< error
| 1064 42000 …near 'select(1 FROM 2)' at line 1
> exec ^
: string
| SELECT select(1 IN 2)
< error
| 1064 42000 …near 'select(1 IN 2)' at line 1
> exec ^
: string
| SELECT select(DISTINCT 1)
< error
| 1064 42000 …near 'select(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT select() OVER ()
< error
| 1064 42000 …near 'select() OVER ()' at line 1
> exec ^
: string
| SELECT select(1) OVER ()
< error
| 1064 42000 …near 'select(1) OVER ()' at line 1
> exec ^
: string
| SELECT select(1,2) OVER ()
< error
| 1064 42000 …near 'select(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT select(1,2,3) OVER ()
< error
| 1064 42000 …near 'select(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT select(1) FROM t
< error
| 1064 42000 …near 'select(1) FROM t' at line 1
> exec ^
: string
| SELECT select(SELECT 1)
< error
| 1064 42000 …near 'select(SELECT 1)' at line 1
> exec ^
: string
| SELECT select(YEAR FROM 1)
< error
| 1064 42000 …near 'select(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT select(DAY, 1, 2)
< error
| 1064 42000 …near 'select(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT select(a)
< error
| 1064 42000 …near 'select(a)' at line 1
> exec ^
: string
| SELECT select(a, b)
< error
| 1064 42000 …near 'select(a, b)' at line 1
> exec ^
: string
| SELECT select(a) OVER ()
< error
| 1064 42000 …near 'select(a) OVER ()' at line 1
> exec ^
: string
| SELECT sensitive()
< error
| 1064 42000 …near 'sensitive()' at line 1
> exec ^
: string
| SELECT sensitive(1)
< error
| 1064 42000 …near 'sensitive(1)' at line 1
> exec ^
: string
| SELECT sensitive(1,2)
< error
| 1064 42000 …near 'sensitive(1,2)' at line 1
> exec ^
: string
| SELECT sensitive(1,2,3)
< error
| 1064 42000 …near 'sensitive(1,2,3)' at line 1
> exec ^
: string
| SELECT sensitive(1,2,3,4)
< error
| 1064 42000 …near 'sensitive(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sensitive(*)
< error
| 1064 42000 …near 'sensitive(*)' at line 1
> exec ^
: string
| SELECT sensitive('a')
< error
| 1064 42000 …near 'sensitive('a')' at line 1
> exec ^
: string
| SELECT sensitive(a USING utf8mb4)
< error
| 1064 42000 …near 'sensitive(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sensitive(1 AS SIGNED)
< error
| 1064 42000 …near 'sensitive(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sensitive(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sensitive(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sensitive(1 FROM 2)
< error
| 1064 42000 …near 'sensitive(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sensitive(1 IN 2)
< error
| 1064 42000 …near 'sensitive(1 IN 2)' at line 1
> exec ^
: string
| SELECT sensitive(DISTINCT 1)
< error
| 1064 42000 …near 'sensitive(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sensitive() OVER ()
< error
| 1064 42000 …near 'sensitive() OVER ()' at line 1
> exec ^
: string
| SELECT sensitive(1) OVER ()
< error
| 1064 42000 …near 'sensitive(1) OVER ()' at line 1
> exec ^
: string
| SELECT sensitive(1,2) OVER ()
< error
| 1064 42000 …near 'sensitive(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sensitive(1,2,3) OVER ()
< error
| 1064 42000 …near 'sensitive(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sensitive(1) FROM t
< error
| 1064 42000 …near 'sensitive(1) FROM t' at line 1
> exec ^
: string
| SELECT sensitive(SELECT 1)
< error
| 1064 42000 …near 'sensitive(SELECT 1)' at line 1
> exec ^
: string
| SELECT sensitive(YEAR FROM 1)
< error
| 1064 42000 …near 'sensitive(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sensitive(DAY, 1, 2)
< error
| 1064 42000 …near 'sensitive(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sensitive(a)
< error
| 1064 42000 …near 'sensitive(a)' at line 1
> exec ^
: string
| SELECT sensitive(a, b)
< error
| 1064 42000 …near 'sensitive(a, b)' at line 1
> exec ^
: string
| SELECT sensitive(a) OVER ()
< error
| 1064 42000 …near 'sensitive(a) OVER ()' at line 1
> exec ^
: string
| SELECT separator()
< error
| 1064 42000 …near 'separator()' at line 1
> exec ^
: string
| SELECT separator(1)
< error
| 1064 42000 …near 'separator(1)' at line 1
> exec ^
: string
| SELECT separator(1,2)
< error
| 1064 42000 …near 'separator(1,2)' at line 1
> exec ^
: string
| SELECT separator(1,2,3)
< error
| 1064 42000 …near 'separator(1,2,3)' at line 1
> exec ^
: string
| SELECT separator(1,2,3,4)
< error
| 1064 42000 …near 'separator(1,2,3,4)' at line 1
> exec ^
: string
| SELECT separator(*)
< error
| 1064 42000 …near 'separator(*)' at line 1
> exec ^
: string
| SELECT separator('a')
< error
| 1064 42000 …near 'separator('a')' at line 1
> exec ^
: string
| SELECT separator(a USING utf8mb4)
< error
| 1064 42000 …near 'separator(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT separator(1 AS SIGNED)
< error
| 1064 42000 …near 'separator(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT separator(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'separator(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT separator(1 FROM 2)
< error
| 1064 42000 …near 'separator(1 FROM 2)' at line 1
> exec ^
: string
| SELECT separator(1 IN 2)
< error
| 1064 42000 …near 'separator(1 IN 2)' at line 1
> exec ^
: string
| SELECT separator(DISTINCT 1)
< error
| 1064 42000 …near 'separator(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT separator() OVER ()
< error
| 1064 42000 …near 'separator() OVER ()' at line 1
> exec ^
: string
| SELECT separator(1) OVER ()
< error
| 1064 42000 …near 'separator(1) OVER ()' at line 1
> exec ^
: string
| SELECT separator(1,2) OVER ()
< error
| 1064 42000 …near 'separator(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT separator(1,2,3) OVER ()
< error
| 1064 42000 …near 'separator(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT separator(1) FROM t
< error
| 1064 42000 …near 'separator(1) FROM t' at line 1
> exec ^
: string
| SELECT separator(SELECT 1)
< error
| 1064 42000 …near 'separator(SELECT 1)' at line 1
> exec ^
: string
| SELECT separator(YEAR FROM 1)
< error
| 1064 42000 …near 'separator(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT separator(DAY, 1, 2)
< error
| 1064 42000 …near 'separator(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT separator(a)
< error
| 1064 42000 …near 'separator(a)' at line 1
> exec ^
: string
| SELECT separator(a, b)
< error
| 1064 42000 …near 'separator(a, b)' at line 1
> exec ^
: string
| SELECT separator(a) OVER ()
< error
| 1064 42000 …near 'separator(a) OVER ()' at line 1
> exec ^
: string
| SELECT sequence()
< error
| 1630 42000 FUNCTION test.sequence does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT sequence()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT serial()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT serial(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT serial(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT serial(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT serial(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT serial(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT serial('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT serial(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT serial(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT serial(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT serial(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT serial(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT serial(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT serial() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT serial(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT serial(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT serial(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT serial(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT serial(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT serial(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT serial(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT serial(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT serial(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT serial(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT serializable()
< error
| 1630 42000 FUNCTION test.serializable does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT serializable()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT server()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT server(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT server(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT server(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT server(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT server(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT server('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT server(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT server(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT server(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT server(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT server(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT server(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT server() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT server(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT server(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT server(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT server(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT server(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT server(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT server(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT server(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT server(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT server(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT session()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT session(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT session(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT session(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT session(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT session(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT session('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT session(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT session(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT session(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT session(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT session(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT session(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT session() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT session(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT session(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT session(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT session(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT session(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT session(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT session(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT session(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT session(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT session(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT session_user()
< ok 0 0
> exec ^
: string
| SELECT set()
< error
| 1064 42000 …near 'set()' at line 1
> exec ^
: string
| SELECT set(1)
< error
| 1064 42000 …near 'set(1)' at line 1
> exec ^
: string
| SELECT set(1,2)
< error
| 1064 42000 …near 'set(1,2)' at line 1
> exec ^
: string
| SELECT set(1,2,3)
< error
| 1064 42000 …near 'set(1,2,3)' at line 1
> exec ^
: string
| SELECT set(1,2,3,4)
< error
| 1064 42000 …near 'set(1,2,3,4)' at line 1
> exec ^
: string
| SELECT set(*)
< error
| 1064 42000 …near 'set(*)' at line 1
> exec ^
: string
| SELECT set('a')
< error
| 1064 42000 …near 'set('a')' at line 1
> exec ^
: string
| SELECT set(a USING utf8mb4)
< error
| 1064 42000 …near 'set(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT set(1 AS SIGNED)
< error
| 1064 42000 …near 'set(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT set(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'set(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT set(1 FROM 2)
< error
| 1064 42000 …near 'set(1 FROM 2)' at line 1
> exec ^
: string
| SELECT set(1 IN 2)
< error
| 1064 42000 …near 'set(1 IN 2)' at line 1
> exec ^
: string
| SELECT set(DISTINCT 1)
< error
| 1064 42000 …near 'set(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT set() OVER ()
< error
| 1064 42000 …near 'set() OVER ()' at line 1
> exec ^
: string
| SELECT set(1) OVER ()
< error
| 1064 42000 …near 'set(1) OVER ()' at line 1
> exec ^
: string
| SELECT set(1,2) OVER ()
< error
| 1064 42000 …near 'set(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT set(1,2,3) OVER ()
< error
| 1064 42000 …near 'set(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT set(1) FROM t
< error
| 1064 42000 …near 'set(1) FROM t' at line 1
> exec ^
: string
| SELECT set(SELECT 1)
< error
| 1064 42000 …near 'set(SELECT 1)' at line 1
> exec ^
: string
| SELECT set(YEAR FROM 1)
< error
| 1064 42000 …near 'set(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT set(DAY, 1, 2)
< error
| 1064 42000 …near 'set(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT set(a)
< error
| 1064 42000 …near 'set(a)' at line 1
> exec ^
: string
| SELECT set(a, b)
< error
| 1064 42000 …near 'set(a, b)' at line 1
> exec ^
: string
| SELECT set(a) OVER ()
< error
| 1064 42000 …near 'set(a) OVER ()' at line 1
> exec ^
: string
| SELECT sets()
< error
| 1305 42000 FUNCTION test.sets does not exist
> exec ^
: string
| SELECT sets()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT share()
< error
| 1630 42000 FUNCTION test.share does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT share()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT show()
< error
| 1064 42000 …near 'show()' at line 1
> exec ^
: string
| SELECT show(1)
< error
| 1064 42000 …near 'show(1)' at line 1
> exec ^
: string
| SELECT show(1,2)
< error
| 1064 42000 …near 'show(1,2)' at line 1
> exec ^
: string
| SELECT show(1,2,3)
< error
| 1064 42000 …near 'show(1,2,3)' at line 1
> exec ^
: string
| SELECT show(1,2,3,4)
< error
| 1064 42000 …near 'show(1,2,3,4)' at line 1
> exec ^
: string
| SELECT show(*)
< error
| 1064 42000 …near 'show(*)' at line 1
> exec ^
: string
| SELECT show('a')
< error
| 1064 42000 …near 'show('a')' at line 1
> exec ^
: string
| SELECT show(a USING utf8mb4)
< error
| 1064 42000 …near 'show(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT show(1 AS SIGNED)
< error
| 1064 42000 …near 'show(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT show(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'show(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT show(1 FROM 2)
< error
| 1064 42000 …near 'show(1 FROM 2)' at line 1
> exec ^
: string
| SELECT show(1 IN 2)
< error
| 1064 42000 …near 'show(1 IN 2)' at line 1
> exec ^
: string
| SELECT show(DISTINCT 1)
< error
| 1064 42000 …near 'show(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT show() OVER ()
< error
| 1064 42000 …near 'show() OVER ()' at line 1
> exec ^
: string
| SELECT show(1) OVER ()
< error
| 1064 42000 …near 'show(1) OVER ()' at line 1
> exec ^
: string
| SELECT show(1,2) OVER ()
< error
| 1064 42000 …near 'show(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT show(1,2,3) OVER ()
< error
| 1064 42000 …near 'show(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT show(1) FROM t
< error
| 1064 42000 …near 'show(1) FROM t' at line 1
> exec ^
: string
| SELECT show(SELECT 1)
< error
| 1064 42000 …near 'show(SELECT 1)' at line 1
> exec ^
: string
| SELECT show(YEAR FROM 1)
< error
| 1064 42000 …near 'show(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT show(DAY, 1, 2)
< error
| 1064 42000 …near 'show(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT show(a)
< error
| 1064 42000 …near 'show(a)' at line 1
> exec ^
: string
| SELECT show(a, b)
< error
| 1064 42000 …near 'show(a, b)' at line 1
> exec ^
: string
| SELECT show(a) OVER ()
< error
| 1064 42000 …near 'show(a) OVER ()' at line 1
> exec ^
: string
| SELECT shutdown()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT shutdown(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT shutdown(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT shutdown(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT shutdown(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT shutdown(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT shutdown('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT shutdown(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT shutdown(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT shutdown(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT shutdown(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT shutdown(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT shutdown(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT shutdown() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT shutdown(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT shutdown(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT shutdown(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT shutdown(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT shutdown(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT shutdown(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT shutdown(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT shutdown(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT shutdown(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT shutdown(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT signal()
< error
| 1064 42000 …near 'signal()' at line 1
> exec ^
: string
| SELECT signal(1)
< error
| 1064 42000 …near 'signal(1)' at line 1
> exec ^
: string
| SELECT signal(1,2)
< error
| 1064 42000 …near 'signal(1,2)' at line 1
> exec ^
: string
| SELECT signal(1,2,3)
< error
| 1064 42000 …near 'signal(1,2,3)' at line 1
> exec ^
: string
| SELECT signal(1,2,3,4)
< error
| 1064 42000 …near 'signal(1,2,3,4)' at line 1
> exec ^
: string
| SELECT signal(*)
< error
| 1064 42000 …near 'signal(*)' at line 1
> exec ^
: string
| SELECT signal('a')
< error
| 1064 42000 …near 'signal('a')' at line 1
> exec ^
: string
| SELECT signal(a USING utf8mb4)
< error
| 1064 42000 …near 'signal(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT signal(1 AS SIGNED)
< error
| 1064 42000 …near 'signal(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT signal(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'signal(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT signal(1 FROM 2)
< error
| 1064 42000 …near 'signal(1 FROM 2)' at line 1
> exec ^
: string
| SELECT signal(1 IN 2)
< error
| 1064 42000 …near 'signal(1 IN 2)' at line 1
> exec ^
: string
| SELECT signal(DISTINCT 1)
< error
| 1064 42000 …near 'signal(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT signal() OVER ()
< error
| 1064 42000 …near 'signal() OVER ()' at line 1
> exec ^
: string
| SELECT signal(1) OVER ()
< error
| 1064 42000 …near 'signal(1) OVER ()' at line 1
> exec ^
: string
| SELECT signal(1,2) OVER ()
< error
| 1064 42000 …near 'signal(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT signal(1,2,3) OVER ()
< error
| 1064 42000 …near 'signal(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT signal(1) FROM t
< error
| 1064 42000 …near 'signal(1) FROM t' at line 1
> exec ^
: string
| SELECT signal(SELECT 1)
< error
| 1064 42000 …near 'signal(SELECT 1)' at line 1
> exec ^
: string
| SELECT signal(YEAR FROM 1)
< error
| 1064 42000 …near 'signal(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT signal(DAY, 1, 2)
< error
| 1064 42000 …near 'signal(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT signal(a)
< error
| 1064 42000 …near 'signal(a)' at line 1
> exec ^
: string
| SELECT signal(a, b)
< error
| 1064 42000 …near 'signal(a, b)' at line 1
> exec ^
: string
| SELECT signal(a) OVER ()
< error
| 1064 42000 …near 'signal(a) OVER ()' at line 1
> exec ^
: string
| SELECT signed()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT signed(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT signed(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT signed(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT signed(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT signed(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT signed('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT signed(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT signed(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT signed(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT signed(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT signed(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT signed(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT signed() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT signed(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT signed(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT signed(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT signed(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT signed(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT signed(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT signed(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT signed(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT signed(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT signed(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT simple()
< error
| 1630 42000 FUNCTION test.simple does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT simple()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT skip()
< error
| 1630 42000 FUNCTION test.skip does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT skip()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT slave()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT slave(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT slave(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT slave(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT slave(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT slave(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT slave('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT slave(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT slave(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT slave(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT slave(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT slave(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT slave(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT slave() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT slave(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT slave(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT slave(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT slave(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT slave(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT slave(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT slave(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT slave(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT slave(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT slave(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT slow()
< error
| 1630 42000 FUNCTION test.slow does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT slow()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT smallint()
< error
| 1064 42000 …near 'smallint()' at line 1
> exec ^
: string
| SELECT smallint(1)
< error
| 1064 42000 …near 'smallint(1)' at line 1
> exec ^
: string
| SELECT smallint(1,2)
< error
| 1064 42000 …near 'smallint(1,2)' at line 1
> exec ^
: string
| SELECT smallint(1,2,3)
< error
| 1064 42000 …near 'smallint(1,2,3)' at line 1
> exec ^
: string
| SELECT smallint(1,2,3,4)
< error
| 1064 42000 …near 'smallint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT smallint(*)
< error
| 1064 42000 …near 'smallint(*)' at line 1
> exec ^
: string
| SELECT smallint('a')
< error
| 1064 42000 …near 'smallint('a')' at line 1
> exec ^
: string
| SELECT smallint(a USING utf8mb4)
< error
| 1064 42000 …near 'smallint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT smallint(1 AS SIGNED)
< error
| 1064 42000 …near 'smallint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT smallint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'smallint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT smallint(1 FROM 2)
< error
| 1064 42000 …near 'smallint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT smallint(1 IN 2)
< error
| 1064 42000 …near 'smallint(1 IN 2)' at line 1
> exec ^
: string
| SELECT smallint(DISTINCT 1)
< error
| 1064 42000 …near 'smallint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT smallint() OVER ()
< error
| 1064 42000 …near 'smallint() OVER ()' at line 1
> exec ^
: string
| SELECT smallint(1) OVER ()
< error
| 1064 42000 …near 'smallint(1) OVER ()' at line 1
> exec ^
: string
| SELECT smallint(1,2) OVER ()
< error
| 1064 42000 …near 'smallint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT smallint(1,2,3) OVER ()
< error
| 1064 42000 …near 'smallint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT smallint(1) FROM t
< error
| 1064 42000 …near 'smallint(1) FROM t' at line 1
> exec ^
: string
| SELECT smallint(SELECT 1)
< error
| 1064 42000 …near 'smallint(SELECT 1)' at line 1
> exec ^
: string
| SELECT smallint(YEAR FROM 1)
< error
| 1064 42000 …near 'smallint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT smallint(DAY, 1, 2)
< error
| 1064 42000 …near 'smallint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT smallint(a)
< error
| 1064 42000 …near 'smallint(a)' at line 1
> exec ^
: string
| SELECT smallint(a, b)
< error
| 1064 42000 …near 'smallint(a, b)' at line 1
> exec ^
: string
| SELECT smallint(a) OVER ()
< error
| 1064 42000 …near 'smallint(a) OVER ()' at line 1
> exec ^
: string
| SELECT snapshot()
< error
| 1630 42000 FUNCTION test.snapshot does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT snapshot()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT socket()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT socket(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT socket(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT socket(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT socket(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT socket(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT socket('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT socket(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT socket(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT socket(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT socket(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT socket(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT socket(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT socket() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT socket(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT socket(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT socket(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT socket(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT socket(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT socket(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT socket(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT socket(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT socket(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT socket(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT some()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT some(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT some(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT some(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT some(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT some(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT some('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT some(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT some(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT some(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT some(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT some(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT some(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT some() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT some(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT some(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT some(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT some(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT some(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT some(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT some(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT some(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT some(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT some(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT soname()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT soname(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT soname(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT soname(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT soname(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT soname(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT soname('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT soname(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT soname(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT soname(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT soname(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT soname(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT soname(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT soname() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT soname(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT soname(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT soname(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT soname(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT soname(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT soname(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT soname(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT soname(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT soname(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT soname(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT sounds()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT sounds(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT sounds(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT sounds(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT sounds(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sounds(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT sounds('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT sounds(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sounds(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sounds(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sounds(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sounds(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT sounds(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sounds() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT sounds(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT sounds(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sounds(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sounds(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT sounds(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT sounds(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sounds(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sounds(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT sounds(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT sounds(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT source()
< error
| 1630 42000 FUNCTION test.source does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT source()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_auto_position()
< error
| 1305 42000 FUNCTION test.source_auto_position does not exist
> exec ^
: string
| SELECT source_auto_position()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_bind()
< error
| 1305 42000 FUNCTION test.source_bind does not exist
> exec ^
: string
| SELECT source_bind()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_compression_algorithms()
< error
| 1305 42000 FUNCTION test.source_compression_algorithms does not exist
> exec ^
: string
| SELECT source_compression_algorithms()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_connect_retry()
< error
| 1305 42000 FUNCTION test.source_connect_retry does not exist
> exec ^
: string
| SELECT source_connect_retry()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_connection_auto_failover()
< error
| 1305 42000 FUNCTION test.source_connection_auto_failover does not exist
> exec ^
: string
| SELECT source_connection_auto_failover()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_delay()
< error
| 1305 42000 FUNCTION test.source_delay does not exist
> exec ^
: string
| SELECT source_delay()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_heartbeat_period()
< error
| 1305 42000 FUNCTION test.source_heartbeat_period does not exist
> exec ^
: string
| SELECT source_heartbeat_period()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_host()
< error
| 1305 42000 FUNCTION test.source_host does not exist
> exec ^
: string
| SELECT source_host()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_log_file()
< error
| 1305 42000 FUNCTION test.source_log_file does not exist
> exec ^
: string
| SELECT source_log_file()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_log_pos()
< error
| 1305 42000 FUNCTION test.source_log_pos does not exist
> exec ^
: string
| SELECT source_log_pos()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_password()
< error
| 1305 42000 FUNCTION test.source_password does not exist
> exec ^
: string
| SELECT source_password()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_port()
< error
| 1305 42000 FUNCTION test.source_port does not exist
> exec ^
: string
| SELECT source_port()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_public_key_path()
< error
| 1305 42000 FUNCTION test.source_public_key_path does not exist
> exec ^
: string
| SELECT source_public_key_path()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_retry_count()
< error
| 1305 42000 FUNCTION test.source_retry_count does not exist
> exec ^
: string
| SELECT source_retry_count()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl()
< error
| 1305 42000 FUNCTION test.source_ssl does not exist
> exec ^
: string
| SELECT source_ssl()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_ca()
< error
| 1305 42000 FUNCTION test.source_ssl_ca does not exist
> exec ^
: string
| SELECT source_ssl_ca()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_capath()
< error
| 1305 42000 FUNCTION test.source_ssl_capath does not exist
> exec ^
: string
| SELECT source_ssl_capath()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_cert()
< error
| 1305 42000 FUNCTION test.source_ssl_cert does not exist
> exec ^
: string
| SELECT source_ssl_cert()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_cipher()
< error
| 1305 42000 FUNCTION test.source_ssl_cipher does not exist
> exec ^
: string
| SELECT source_ssl_cipher()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_crl()
< error
| 1305 42000 FUNCTION test.source_ssl_crl does not exist
> exec ^
: string
| SELECT source_ssl_crl()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_crlpath()
< error
| 1305 42000 FUNCTION test.source_ssl_crlpath does not exist
> exec ^
: string
| SELECT source_ssl_crlpath()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_key()
< error
| 1305 42000 FUNCTION test.source_ssl_key does not exist
> exec ^
: string
| SELECT source_ssl_key()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_ssl_verify_server_cert()
< error
| 1305 42000 FUNCTION test.source_ssl_verify_server_cert does not exist
> exec ^
: string
| SELECT source_ssl_verify_server_cert()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_tls_ciphersuites()
< error
| 1305 42000 FUNCTION test.source_tls_ciphersuites does not exist
> exec ^
: string
| SELECT source_tls_ciphersuites()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_tls_version()
< error
| 1305 42000 FUNCTION test.source_tls_version does not exist
> exec ^
: string
| SELECT source_tls_version()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_user()
< error
| 1305 42000 FUNCTION test.source_user does not exist
> exec ^
: string
| SELECT source_user()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT source_zstd_compression_level()
< error
| 1305 42000 FUNCTION test.source_zstd_compression_level does not exist
> exec ^
: string
| SELECT source_zstd_compression_level()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT spatial()
< error
| 1064 42000 …near 'spatial()' at line 1
> exec ^
: string
| SELECT spatial(1)
< error
| 1064 42000 …near 'spatial(1)' at line 1
> exec ^
: string
| SELECT spatial(1,2)
< error
| 1064 42000 …near 'spatial(1,2)' at line 1
> exec ^
: string
| SELECT spatial(1,2,3)
< error
| 1064 42000 …near 'spatial(1,2,3)' at line 1
> exec ^
: string
| SELECT spatial(1,2,3,4)
< error
| 1064 42000 …near 'spatial(1,2,3,4)' at line 1
> exec ^
: string
| SELECT spatial(*)
< error
| 1064 42000 …near 'spatial(*)' at line 1
> exec ^
: string
| SELECT spatial('a')
< error
| 1064 42000 …near 'spatial('a')' at line 1
> exec ^
: string
| SELECT spatial(a USING utf8mb4)
< error
| 1064 42000 …near 'spatial(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT spatial(1 AS SIGNED)
< error
| 1064 42000 …near 'spatial(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT spatial(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'spatial(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT spatial(1 FROM 2)
< error
| 1064 42000 …near 'spatial(1 FROM 2)' at line 1
> exec ^
: string
| SELECT spatial(1 IN 2)
< error
| 1064 42000 …near 'spatial(1 IN 2)' at line 1
> exec ^
: string
| SELECT spatial(DISTINCT 1)
< error
| 1064 42000 …near 'spatial(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT spatial() OVER ()
< error
| 1064 42000 …near 'spatial() OVER ()' at line 1
> exec ^
: string
| SELECT spatial(1) OVER ()
< error
| 1064 42000 …near 'spatial(1) OVER ()' at line 1
> exec ^
: string
| SELECT spatial(1,2) OVER ()
< error
| 1064 42000 …near 'spatial(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT spatial(1,2,3) OVER ()
< error
| 1064 42000 …near 'spatial(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT spatial(1) FROM t
< error
| 1064 42000 …near 'spatial(1) FROM t' at line 1
> exec ^
: string
| SELECT spatial(SELECT 1)
< error
| 1064 42000 …near 'spatial(SELECT 1)' at line 1
> exec ^
: string
| SELECT spatial(YEAR FROM 1)
< error
| 1064 42000 …near 'spatial(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT spatial(DAY, 1, 2)
< error
| 1064 42000 …near 'spatial(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT spatial(a)
< error
| 1064 42000 …near 'spatial(a)' at line 1
> exec ^
: string
| SELECT spatial(a, b)
< error
| 1064 42000 …near 'spatial(a, b)' at line 1
> exec ^
: string
| SELECT spatial(a) OVER ()
< error
| 1064 42000 …near 'spatial(a) OVER ()' at line 1
> exec ^
: string
| SELECT specific()
< error
| 1064 42000 …near 'specific()' at line 1
> exec ^
: string
| SELECT specific(1)
< error
| 1064 42000 …near 'specific(1)' at line 1
> exec ^
: string
| SELECT specific(1,2)
< error
| 1064 42000 …near 'specific(1,2)' at line 1
> exec ^
: string
| SELECT specific(1,2,3)
< error
| 1064 42000 …near 'specific(1,2,3)' at line 1
> exec ^
: string
| SELECT specific(1,2,3,4)
< error
| 1064 42000 …near 'specific(1,2,3,4)' at line 1
> exec ^
: string
| SELECT specific(*)
< error
| 1064 42000 …near 'specific(*)' at line 1
> exec ^
: string
| SELECT specific('a')
< error
| 1064 42000 …near 'specific('a')' at line 1
> exec ^
: string
| SELECT specific(a USING utf8mb4)
< error
| 1064 42000 …near 'specific(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT specific(1 AS SIGNED)
< error
| 1064 42000 …near 'specific(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT specific(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'specific(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT specific(1 FROM 2)
< error
| 1064 42000 …near 'specific(1 FROM 2)' at line 1
> exec ^
: string
| SELECT specific(1 IN 2)
< error
| 1064 42000 …near 'specific(1 IN 2)' at line 1
> exec ^
: string
| SELECT specific(DISTINCT 1)
< error
| 1064 42000 …near 'specific(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT specific() OVER ()
< error
| 1064 42000 …near 'specific() OVER ()' at line 1
> exec ^
: string
| SELECT specific(1) OVER ()
< error
| 1064 42000 …near 'specific(1) OVER ()' at line 1
> exec ^
: string
| SELECT specific(1,2) OVER ()
< error
| 1064 42000 …near 'specific(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT specific(1,2,3) OVER ()
< error
| 1064 42000 …near 'specific(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT specific(1) FROM t
< error
| 1064 42000 …near 'specific(1) FROM t' at line 1
> exec ^
: string
| SELECT specific(SELECT 1)
< error
| 1064 42000 …near 'specific(SELECT 1)' at line 1
> exec ^
: string
| SELECT specific(YEAR FROM 1)
< error
| 1064 42000 …near 'specific(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT specific(DAY, 1, 2)
< error
| 1064 42000 …near 'specific(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT specific(a)
< error
| 1064 42000 …near 'specific(a)' at line 1
> exec ^
: string
| SELECT specific(a, b)
< error
| 1064 42000 …near 'specific(a, b)' at line 1
> exec ^
: string
| SELECT specific(a) OVER ()
< error
| 1064 42000 …near 'specific(a) OVER ()' at line 1
> exec ^
: string
| SELECT sql()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT sql(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT sql(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT sql(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT sql(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sql(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT sql('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT sql(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sql(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sql(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sql(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sql(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT sql(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sql() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT sql(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT sql(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sql(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sql(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT sql(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT sql(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sql(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sql(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT sql(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT sql(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_gtids()
< error
| 1064 42000 …near 'sql_after_gtids()' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1)
< error
| 1064 42000 …near 'sql_after_gtids(1)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1,2)
< error
| 1064 42000 …near 'sql_after_gtids(1,2)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1,2,3)
< error
| 1064 42000 …near 'sql_after_gtids(1,2,3)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1,2,3,4)
< error
| 1064 42000 …near 'sql_after_gtids(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(*)
< error
| 1064 42000 …near 'sql_after_gtids(*)' at line 1
> exec ^
: string
| SELECT sql_after_gtids('a')
< error
| 1064 42000 …near 'sql_after_gtids('a')' at line 1
> exec ^
: string
| SELECT sql_after_gtids(a USING utf8mb4)
< error
| 1064 42000 …near 'sql_after_gtids(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1 AS SIGNED)
< error
| 1064 42000 …near 'sql_after_gtids(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sql_after_gtids(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1 FROM 2)
< error
| 1064 42000 …near 'sql_after_gtids(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1 IN 2)
< error
| 1064 42000 …near 'sql_after_gtids(1 IN 2)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(DISTINCT 1)
< error
| 1064 42000 …near 'sql_after_gtids(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sql_after_gtids() OVER ()
< error
| 1064 42000 …near 'sql_after_gtids() OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1) OVER ()
< error
| 1064 42000 …near 'sql_after_gtids(1) OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1,2) OVER ()
< error
| 1064 42000 …near 'sql_after_gtids(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1,2,3) OVER ()
< error
| 1064 42000 …near 'sql_after_gtids(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_gtids(1) FROM t
< error
| 1064 42000 …near 'sql_after_gtids(1) FROM t' at line 1
> exec ^
: string
| SELECT sql_after_gtids(SELECT 1)
< error
| 1064 42000 …near 'sql_after_gtids(SELECT 1)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(YEAR FROM 1)
< error
| 1064 42000 …near 'sql_after_gtids(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(DAY, 1, 2)
< error
| 1064 42000 …near 'sql_after_gtids(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(a)
< error
| 1064 42000 …near 'sql_after_gtids(a)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(a, b)
< error
| 1064 42000 …near 'sql_after_gtids(a, b)' at line 1
> exec ^
: string
| SELECT sql_after_gtids(a) OVER ()
< error
| 1064 42000 …near 'sql_after_gtids(a) OVER ()' at line 1
> exec ^
: string
| SELECT sql_after_mts_gaps()
< error
| 1305 42000 FUNCTION test.sql_after_mts_gaps does not exist
> exec ^
: string
| SELECT sql_after_mts_gaps()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT sql_before_gtids()
< error
| 1064 42000 …near 'sql_before_gtids()' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1)
< error
| 1064 42000 …near 'sql_before_gtids(1)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1,2)
< error
| 1064 42000 …near 'sql_before_gtids(1,2)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1,2,3)
< error
| 1064 42000 …near 'sql_before_gtids(1,2,3)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1,2,3,4)
< error
| 1064 42000 …near 'sql_before_gtids(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(*)
< error
| 1064 42000 …near 'sql_before_gtids(*)' at line 1
> exec ^
: string
| SELECT sql_before_gtids('a')
< error
| 1064 42000 …near 'sql_before_gtids('a')' at line 1
> exec ^
: string
| SELECT sql_before_gtids(a USING utf8mb4)
< error
| 1064 42000 …near 'sql_before_gtids(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1 AS SIGNED)
< error
| 1064 42000 …near 'sql_before_gtids(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sql_before_gtids(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1 FROM 2)
< error
| 1064 42000 …near 'sql_before_gtids(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1 IN 2)
< error
| 1064 42000 …near 'sql_before_gtids(1 IN 2)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(DISTINCT 1)
< error
| 1064 42000 …near 'sql_before_gtids(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sql_before_gtids() OVER ()
< error
| 1064 42000 …near 'sql_before_gtids() OVER ()' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1) OVER ()
< error
| 1064 42000 …near 'sql_before_gtids(1) OVER ()' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1,2) OVER ()
< error
| 1064 42000 …near 'sql_before_gtids(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1,2,3) OVER ()
< error
| 1064 42000 …near 'sql_before_gtids(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sql_before_gtids(1) FROM t
< error
| 1064 42000 …near 'sql_before_gtids(1) FROM t' at line 1
> exec ^
: string
| SELECT sql_before_gtids(SELECT 1)
< error
| 1064 42000 …near 'sql_before_gtids(SELECT 1)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(YEAR FROM 1)
< error
| 1064 42000 …near 'sql_before_gtids(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(DAY, 1, 2)
< error
| 1064 42000 …near 'sql_before_gtids(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(a)
< error
| 1064 42000 …near 'sql_before_gtids(a)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(a, b)
< error
| 1064 42000 …near 'sql_before_gtids(a, b)' at line 1
> exec ^
: string
| SELECT sql_before_gtids(a) OVER ()
< error
| 1064 42000 …near 'sql_before_gtids(a) OVER ()' at line 1
> exec ^
: string
| SELECT sql_big_result()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_big_result(1)
< ok 0 0
> exec ^
: string
| SELECT sql_buffer_result()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_buffer_result(1)
< ok 0 0
> exec ^
: string
| SELECT sql_calc_found_rows()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_calc_found_rows(1)
< ok 0 0
> exec ^
: string
| SELECT sql_no_cache()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_no_cache(1)
< ok 0 0
> exec ^
: string
| SELECT sql_small_result()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_small_result(1)
< ok 0 0
> exec ^
: string
| SELECT sql_thread()
< error
| 1630 42000 FUNCTION test.sql_thread does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT sql_thread()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT sql_tsi_day()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_day(1)
< ok 0 0
> exec ^
: string
| SELECT sql_tsi_hour()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_hour(1)
< ok 0 0
> exec ^
: string
| SELECT sql_tsi_minute()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_minute(1)
< ok 0 0
> exec ^
: string
| SELECT sql_tsi_month()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_month(1)
< ok 0 0
> exec ^
: string
| SELECT sql_tsi_quarter()
< error
| 1630 42000 FUNCTION test.sql_tsi_quarter does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT sql_tsi_quarter()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT sql_tsi_second()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_second(1)
< ok 0 0
> exec ^
: string
| SELECT sql_tsi_week()
< error
| 1630 42000 FUNCTION test.sql_tsi_week does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT sql_tsi_week()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT sql_tsi_year()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sql_tsi_year(1)
< ok 0 0
> exec ^
: string
| SELECT sqlexception()
< error
| 1064 42000 …near 'sqlexception()' at line 1
> exec ^
: string
| SELECT sqlexception(1)
< error
| 1064 42000 …near 'sqlexception(1)' at line 1
> exec ^
: string
| SELECT sqlexception(1,2)
< error
| 1064 42000 …near 'sqlexception(1,2)' at line 1
> exec ^
: string
| SELECT sqlexception(1,2,3)
< error
| 1064 42000 …near 'sqlexception(1,2,3)' at line 1
> exec ^
: string
| SELECT sqlexception(1,2,3,4)
< error
| 1064 42000 …near 'sqlexception(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sqlexception(*)
< error
| 1064 42000 …near 'sqlexception(*)' at line 1
> exec ^
: string
| SELECT sqlexception('a')
< error
| 1064 42000 …near 'sqlexception('a')' at line 1
> exec ^
: string
| SELECT sqlexception(a USING utf8mb4)
< error
| 1064 42000 …near 'sqlexception(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sqlexception(1 AS SIGNED)
< error
| 1064 42000 …near 'sqlexception(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sqlexception(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sqlexception(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sqlexception(1 FROM 2)
< error
| 1064 42000 …near 'sqlexception(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sqlexception(1 IN 2)
< error
| 1064 42000 …near 'sqlexception(1 IN 2)' at line 1
> exec ^
: string
| SELECT sqlexception(DISTINCT 1)
< error
| 1064 42000 …near 'sqlexception(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sqlexception() OVER ()
< error
| 1064 42000 …near 'sqlexception() OVER ()' at line 1
> exec ^
: string
| SELECT sqlexception(1) OVER ()
< error
| 1064 42000 …near 'sqlexception(1) OVER ()' at line 1
> exec ^
: string
| SELECT sqlexception(1,2) OVER ()
< error
| 1064 42000 …near 'sqlexception(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sqlexception(1,2,3) OVER ()
< error
| 1064 42000 …near 'sqlexception(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sqlexception(1) FROM t
< error
| 1064 42000 …near 'sqlexception(1) FROM t' at line 1
> exec ^
: string
| SELECT sqlexception(SELECT 1)
< error
| 1064 42000 …near 'sqlexception(SELECT 1)' at line 1
> exec ^
: string
| SELECT sqlexception(YEAR FROM 1)
< error
| 1064 42000 …near 'sqlexception(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sqlexception(DAY, 1, 2)
< error
| 1064 42000 …near 'sqlexception(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sqlexception(a)
< error
| 1064 42000 …near 'sqlexception(a)' at line 1
> exec ^
: string
| SELECT sqlexception(a, b)
< error
| 1064 42000 …near 'sqlexception(a, b)' at line 1
> exec ^
: string
| SELECT sqlexception(a) OVER ()
< error
| 1064 42000 …near 'sqlexception(a) OVER ()' at line 1
> exec ^
: string
| SELECT sqlstate()
< error
| 1064 42000 …near 'sqlstate()' at line 1
> exec ^
: string
| SELECT sqlstate(1)
< error
| 1064 42000 …near 'sqlstate(1)' at line 1
> exec ^
: string
| SELECT sqlstate(1,2)
< error
| 1064 42000 …near 'sqlstate(1,2)' at line 1
> exec ^
: string
| SELECT sqlstate(1,2,3)
< error
| 1064 42000 …near 'sqlstate(1,2,3)' at line 1
> exec ^
: string
| SELECT sqlstate(1,2,3,4)
< error
| 1064 42000 …near 'sqlstate(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sqlstate(*)
< error
| 1064 42000 …near 'sqlstate(*)' at line 1
> exec ^
: string
| SELECT sqlstate('a')
< error
| 1064 42000 …near 'sqlstate('a')' at line 1
> exec ^
: string
| SELECT sqlstate(a USING utf8mb4)
< error
| 1064 42000 …near 'sqlstate(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sqlstate(1 AS SIGNED)
< error
| 1064 42000 …near 'sqlstate(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sqlstate(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sqlstate(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sqlstate(1 FROM 2)
< error
| 1064 42000 …near 'sqlstate(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sqlstate(1 IN 2)
< error
| 1064 42000 …near 'sqlstate(1 IN 2)' at line 1
> exec ^
: string
| SELECT sqlstate(DISTINCT 1)
< error
| 1064 42000 …near 'sqlstate(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sqlstate() OVER ()
< error
| 1064 42000 …near 'sqlstate() OVER ()' at line 1
> exec ^
: string
| SELECT sqlstate(1) OVER ()
< error
| 1064 42000 …near 'sqlstate(1) OVER ()' at line 1
> exec ^
: string
| SELECT sqlstate(1,2) OVER ()
< error
| 1064 42000 …near 'sqlstate(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sqlstate(1,2,3) OVER ()
< error
| 1064 42000 …near 'sqlstate(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sqlstate(1) FROM t
< error
| 1064 42000 …near 'sqlstate(1) FROM t' at line 1
> exec ^
: string
| SELECT sqlstate(SELECT 1)
< error
| 1064 42000 …near 'sqlstate(SELECT 1)' at line 1
> exec ^
: string
| SELECT sqlstate(YEAR FROM 1)
< error
| 1064 42000 …near 'sqlstate(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sqlstate(DAY, 1, 2)
< error
| 1064 42000 …near 'sqlstate(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sqlstate(a)
< error
| 1064 42000 …near 'sqlstate(a)' at line 1
> exec ^
: string
| SELECT sqlstate(a, b)
< error
| 1064 42000 …near 'sqlstate(a, b)' at line 1
> exec ^
: string
| SELECT sqlstate(a) OVER ()
< error
| 1064 42000 …near 'sqlstate(a) OVER ()' at line 1
> exec ^
: string
| SELECT sqlwarning()
< error
| 1064 42000 …near 'sqlwarning()' at line 1
> exec ^
: string
| SELECT sqlwarning(1)
< error
| 1064 42000 …near 'sqlwarning(1)' at line 1
> exec ^
: string
| SELECT sqlwarning(1,2)
< error
| 1064 42000 …near 'sqlwarning(1,2)' at line 1
> exec ^
: string
| SELECT sqlwarning(1,2,3)
< error
| 1064 42000 …near 'sqlwarning(1,2,3)' at line 1
> exec ^
: string
| SELECT sqlwarning(1,2,3,4)
< error
| 1064 42000 …near 'sqlwarning(1,2,3,4)' at line 1
> exec ^
: string
| SELECT sqlwarning(*)
< error
| 1064 42000 …near 'sqlwarning(*)' at line 1
> exec ^
: string
| SELECT sqlwarning('a')
< error
| 1064 42000 …near 'sqlwarning('a')' at line 1
> exec ^
: string
| SELECT sqlwarning(a USING utf8mb4)
< error
| 1064 42000 …near 'sqlwarning(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT sqlwarning(1 AS SIGNED)
< error
| 1064 42000 …near 'sqlwarning(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT sqlwarning(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'sqlwarning(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT sqlwarning(1 FROM 2)
< error
| 1064 42000 …near 'sqlwarning(1 FROM 2)' at line 1
> exec ^
: string
| SELECT sqlwarning(1 IN 2)
< error
| 1064 42000 …near 'sqlwarning(1 IN 2)' at line 1
> exec ^
: string
| SELECT sqlwarning(DISTINCT 1)
< error
| 1064 42000 …near 'sqlwarning(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT sqlwarning() OVER ()
< error
| 1064 42000 …near 'sqlwarning() OVER ()' at line 1
> exec ^
: string
| SELECT sqlwarning(1) OVER ()
< error
| 1064 42000 …near 'sqlwarning(1) OVER ()' at line 1
> exec ^
: string
| SELECT sqlwarning(1,2) OVER ()
< error
| 1064 42000 …near 'sqlwarning(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT sqlwarning(1,2,3) OVER ()
< error
| 1064 42000 …near 'sqlwarning(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT sqlwarning(1) FROM t
< error
| 1064 42000 …near 'sqlwarning(1) FROM t' at line 1
> exec ^
: string
| SELECT sqlwarning(SELECT 1)
< error
| 1064 42000 …near 'sqlwarning(SELECT 1)' at line 1
> exec ^
: string
| SELECT sqlwarning(YEAR FROM 1)
< error
| 1064 42000 …near 'sqlwarning(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT sqlwarning(DAY, 1, 2)
< error
| 1064 42000 …near 'sqlwarning(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT sqlwarning(a)
< error
| 1064 42000 …near 'sqlwarning(a)' at line 1
> exec ^
: string
| SELECT sqlwarning(a, b)
< error
| 1064 42000 …near 'sqlwarning(a, b)' at line 1
> exec ^
: string
| SELECT sqlwarning(a) OVER ()
< error
| 1064 42000 …near 'sqlwarning(a) OVER ()' at line 1
> exec ^
: string
| SELECT srid()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'srid'
> exec ^
: string
| SELECT srid()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'srid'
> exec ^
: string
| SELECT srid(1)
< error
| 4079 HY000 Illegal parameter data type int for operation 'srid'
> exec ^
: string
| SELECT srid(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ssl()
< error
| 1064 42000 …near 'ssl()' at line 1
> exec ^
: string
| SELECT ssl(1)
< error
| 1064 42000 …near 'ssl(1)' at line 1
> exec ^
: string
| SELECT ssl(1,2)
< error
| 1064 42000 …near 'ssl(1,2)' at line 1
> exec ^
: string
| SELECT ssl(1,2,3)
< error
| 1064 42000 …near 'ssl(1,2,3)' at line 1
> exec ^
: string
| SELECT ssl(1,2,3,4)
< error
| 1064 42000 …near 'ssl(1,2,3,4)' at line 1
> exec ^
: string
| SELECT ssl(*)
< error
| 1064 42000 …near 'ssl(*)' at line 1
> exec ^
: string
| SELECT ssl('a')
< error
| 1064 42000 …near 'ssl('a')' at line 1
> exec ^
: string
| SELECT ssl(a USING utf8mb4)
< error
| 1064 42000 …near 'ssl(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT ssl(1 AS SIGNED)
< error
| 1064 42000 …near 'ssl(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT ssl(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'ssl(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT ssl(1 FROM 2)
< error
| 1064 42000 …near 'ssl(1 FROM 2)' at line 1
> exec ^
: string
| SELECT ssl(1 IN 2)
< error
| 1064 42000 …near 'ssl(1 IN 2)' at line 1
> exec ^
: string
| SELECT ssl(DISTINCT 1)
< error
| 1064 42000 …near 'ssl(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT ssl() OVER ()
< error
| 1064 42000 …near 'ssl() OVER ()' at line 1
> exec ^
: string
| SELECT ssl(1) OVER ()
< error
| 1064 42000 …near 'ssl(1) OVER ()' at line 1
> exec ^
: string
| SELECT ssl(1,2) OVER ()
< error
| 1064 42000 …near 'ssl(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT ssl(1,2,3) OVER ()
< error
| 1064 42000 …near 'ssl(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT ssl(1) FROM t
< error
| 1064 42000 …near 'ssl(1) FROM t' at line 1
> exec ^
: string
| SELECT ssl(SELECT 1)
< error
| 1064 42000 …near 'ssl(SELECT 1)' at line 1
> exec ^
: string
| SELECT ssl(YEAR FROM 1)
< error
| 1064 42000 …near 'ssl(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT ssl(DAY, 1, 2)
< error
| 1064 42000 …near 'ssl(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT ssl(a)
< error
| 1064 42000 …near 'ssl(a)' at line 1
> exec ^
: string
| SELECT ssl(a, b)
< error
| 1064 42000 …near 'ssl(a, b)' at line 1
> exec ^
: string
| SELECT ssl(a) OVER ()
< error
| 1064 42000 …near 'ssl(a) OVER ()' at line 1
> exec ^
: string
| SELECT st_collect()
< error
| 1305 42000 FUNCTION test.st_collect does not exist
> exec ^
: string
| SELECT st_collect()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT stacked()
< error
| 1305 42000 FUNCTION test.stacked does not exist
> exec ^
: string
| SELECT stacked()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT start()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT start(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT start(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT start(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT start(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT start(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT start('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT start(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT start(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT start(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT start(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT start(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT start(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT start() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT start(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT start(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT start(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT start(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT start(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT start(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT start(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT start(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT start(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT start(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT starting()
< error
| 1064 42000 …near 'starting()' at line 1
> exec ^
: string
| SELECT starting(1)
< error
| 1064 42000 …near 'starting(1)' at line 1
> exec ^
: string
| SELECT starting(1,2)
< error
| 1064 42000 …near 'starting(1,2)' at line 1
> exec ^
: string
| SELECT starting(1,2,3)
< error
| 1064 42000 …near 'starting(1,2,3)' at line 1
> exec ^
: string
| SELECT starting(1,2,3,4)
< error
| 1064 42000 …near 'starting(1,2,3,4)' at line 1
> exec ^
: string
| SELECT starting(*)
< error
| 1064 42000 …near 'starting(*)' at line 1
> exec ^
: string
| SELECT starting('a')
< error
| 1064 42000 …near 'starting('a')' at line 1
> exec ^
: string
| SELECT starting(a USING utf8mb4)
< error
| 1064 42000 …near 'starting(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT starting(1 AS SIGNED)
< error
| 1064 42000 …near 'starting(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT starting(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'starting(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT starting(1 FROM 2)
< error
| 1064 42000 …near 'starting(1 FROM 2)' at line 1
> exec ^
: string
| SELECT starting(1 IN 2)
< error
| 1064 42000 …near 'starting(1 IN 2)' at line 1
> exec ^
: string
| SELECT starting(DISTINCT 1)
< error
| 1064 42000 …near 'starting(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT starting() OVER ()
< error
| 1064 42000 …near 'starting() OVER ()' at line 1
> exec ^
: string
| SELECT starting(1) OVER ()
< error
| 1064 42000 …near 'starting(1) OVER ()' at line 1
> exec ^
: string
| SELECT starting(1,2) OVER ()
< error
| 1064 42000 …near 'starting(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT starting(1,2,3) OVER ()
< error
| 1064 42000 …near 'starting(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT starting(1) FROM t
< error
| 1064 42000 …near 'starting(1) FROM t' at line 1
> exec ^
: string
| SELECT starting(SELECT 1)
< error
| 1064 42000 …near 'starting(SELECT 1)' at line 1
> exec ^
: string
| SELECT starting(YEAR FROM 1)
< error
| 1064 42000 …near 'starting(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT starting(DAY, 1, 2)
< error
| 1064 42000 …near 'starting(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT starting(a)
< error
| 1064 42000 …near 'starting(a)' at line 1
> exec ^
: string
| SELECT starting(a, b)
< error
| 1064 42000 …near 'starting(a, b)' at line 1
> exec ^
: string
| SELECT starting(a) OVER ()
< error
| 1064 42000 …near 'starting(a) OVER ()' at line 1
> exec ^
: string
| SELECT starts()
< error
| 1630 42000 FUNCTION test.starts does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT starts()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT stats_auto_recalc()
< error
| 1064 42000 …near 'stats_auto_recalc()' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1)
< error
| 1064 42000 …near 'stats_auto_recalc(1)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1,2)
< error
| 1064 42000 …near 'stats_auto_recalc(1,2)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1,2,3)
< error
| 1064 42000 …near 'stats_auto_recalc(1,2,3)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1,2,3,4)
< error
| 1064 42000 …near 'stats_auto_recalc(1,2,3,4)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(*)
< error
| 1064 42000 …near 'stats_auto_recalc(*)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc('a')
< error
| 1064 42000 …near 'stats_auto_recalc('a')' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(a USING utf8mb4)
< error
| 1064 42000 …near 'stats_auto_recalc(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1 AS SIGNED)
< error
| 1064 42000 …near 'stats_auto_recalc(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'stats_auto_recalc(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1 FROM 2)
< error
| 1064 42000 …near 'stats_auto_recalc(1 FROM 2)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1 IN 2)
< error
| 1064 42000 …near 'stats_auto_recalc(1 IN 2)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(DISTINCT 1)
< error
| 1064 42000 …near 'stats_auto_recalc(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc() OVER ()
< error
| 1064 42000 …near 'stats_auto_recalc() OVER ()' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1) OVER ()
< error
| 1064 42000 …near 'stats_auto_recalc(1) OVER ()' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1,2) OVER ()
< error
| 1064 42000 …near 'stats_auto_recalc(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1,2,3) OVER ()
< error
| 1064 42000 …near 'stats_auto_recalc(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(1) FROM t
< error
| 1064 42000 …near 'stats_auto_recalc(1) FROM t' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(SELECT 1)
< error
| 1064 42000 …near 'stats_auto_recalc(SELECT 1)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(YEAR FROM 1)
< error
| 1064 42000 …near 'stats_auto_recalc(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(DAY, 1, 2)
< error
| 1064 42000 …near 'stats_auto_recalc(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(a)
< error
| 1064 42000 …near 'stats_auto_recalc(a)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(a, b)
< error
| 1064 42000 …near 'stats_auto_recalc(a, b)' at line 1
> exec ^
: string
| SELECT stats_auto_recalc(a) OVER ()
< error
| 1064 42000 …near 'stats_auto_recalc(a) OVER ()' at line 1
> exec ^
: string
| SELECT stats_persistent()
< error
| 1064 42000 …near 'stats_persistent()' at line 1
> exec ^
: string
| SELECT stats_persistent(1)
< error
| 1064 42000 …near 'stats_persistent(1)' at line 1
> exec ^
: string
| SELECT stats_persistent(1,2)
< error
| 1064 42000 …near 'stats_persistent(1,2)' at line 1
> exec ^
: string
| SELECT stats_persistent(1,2,3)
< error
| 1064 42000 …near 'stats_persistent(1,2,3)' at line 1
> exec ^
: string
| SELECT stats_persistent(1,2,3,4)
< error
| 1064 42000 …near 'stats_persistent(1,2,3,4)' at line 1
> exec ^
: string
| SELECT stats_persistent(*)
< error
| 1064 42000 …near 'stats_persistent(*)' at line 1
> exec ^
: string
| SELECT stats_persistent('a')
< error
| 1064 42000 …near 'stats_persistent('a')' at line 1
> exec ^
: string
| SELECT stats_persistent(a USING utf8mb4)
< error
| 1064 42000 …near 'stats_persistent(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT stats_persistent(1 AS SIGNED)
< error
| 1064 42000 …near 'stats_persistent(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT stats_persistent(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'stats_persistent(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT stats_persistent(1 FROM 2)
< error
| 1064 42000 …near 'stats_persistent(1 FROM 2)' at line 1
> exec ^
: string
| SELECT stats_persistent(1 IN 2)
< error
| 1064 42000 …near 'stats_persistent(1 IN 2)' at line 1
> exec ^
: string
| SELECT stats_persistent(DISTINCT 1)
< error
| 1064 42000 …near 'stats_persistent(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT stats_persistent() OVER ()
< error
| 1064 42000 …near 'stats_persistent() OVER ()' at line 1
> exec ^
: string
| SELECT stats_persistent(1) OVER ()
< error
| 1064 42000 …near 'stats_persistent(1) OVER ()' at line 1
> exec ^
: string
| SELECT stats_persistent(1,2) OVER ()
< error
| 1064 42000 …near 'stats_persistent(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT stats_persistent(1,2,3) OVER ()
< error
| 1064 42000 …near 'stats_persistent(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT stats_persistent(1) FROM t
< error
| 1064 42000 …near 'stats_persistent(1) FROM t' at line 1
> exec ^
: string
| SELECT stats_persistent(SELECT 1)
< error
| 1064 42000 …near 'stats_persistent(SELECT 1)' at line 1
> exec ^
: string
| SELECT stats_persistent(YEAR FROM 1)
< error
| 1064 42000 …near 'stats_persistent(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT stats_persistent(DAY, 1, 2)
< error
| 1064 42000 …near 'stats_persistent(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT stats_persistent(a)
< error
| 1064 42000 …near 'stats_persistent(a)' at line 1
> exec ^
: string
| SELECT stats_persistent(a, b)
< error
| 1064 42000 …near 'stats_persistent(a, b)' at line 1
> exec ^
: string
| SELECT stats_persistent(a) OVER ()
< error
| 1064 42000 …near 'stats_persistent(a) OVER ()' at line 1
> exec ^
: string
| SELECT stats_sample_pages()
< error
| 1064 42000 …near 'stats_sample_pages()' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1)
< error
| 1064 42000 …near 'stats_sample_pages(1)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1,2)
< error
| 1064 42000 …near 'stats_sample_pages(1,2)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1,2,3)
< error
| 1064 42000 …near 'stats_sample_pages(1,2,3)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1,2,3,4)
< error
| 1064 42000 …near 'stats_sample_pages(1,2,3,4)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(*)
< error
| 1064 42000 …near 'stats_sample_pages(*)' at line 1
> exec ^
: string
| SELECT stats_sample_pages('a')
< error
| 1064 42000 …near 'stats_sample_pages('a')' at line 1
> exec ^
: string
| SELECT stats_sample_pages(a USING utf8mb4)
< error
| 1064 42000 …near 'stats_sample_pages(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1 AS SIGNED)
< error
| 1064 42000 …near 'stats_sample_pages(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'stats_sample_pages(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1 FROM 2)
< error
| 1064 42000 …near 'stats_sample_pages(1 FROM 2)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1 IN 2)
< error
| 1064 42000 …near 'stats_sample_pages(1 IN 2)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(DISTINCT 1)
< error
| 1064 42000 …near 'stats_sample_pages(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT stats_sample_pages() OVER ()
< error
| 1064 42000 …near 'stats_sample_pages() OVER ()' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1) OVER ()
< error
| 1064 42000 …near 'stats_sample_pages(1) OVER ()' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1,2) OVER ()
< error
| 1064 42000 …near 'stats_sample_pages(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1,2,3) OVER ()
< error
| 1064 42000 …near 'stats_sample_pages(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT stats_sample_pages(1) FROM t
< error
| 1064 42000 …near 'stats_sample_pages(1) FROM t' at line 1
> exec ^
: string
| SELECT stats_sample_pages(SELECT 1)
< error
| 1064 42000 …near 'stats_sample_pages(SELECT 1)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(YEAR FROM 1)
< error
| 1064 42000 …near 'stats_sample_pages(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(DAY, 1, 2)
< error
| 1064 42000 …near 'stats_sample_pages(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(a)
< error
| 1064 42000 …near 'stats_sample_pages(a)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(a, b)
< error
| 1064 42000 …near 'stats_sample_pages(a, b)' at line 1
> exec ^
: string
| SELECT stats_sample_pages(a) OVER ()
< error
| 1064 42000 …near 'stats_sample_pages(a) OVER ()' at line 1
> exec ^
: string
| SELECT status()
< error
| 1630 42000 FUNCTION test.status does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT status()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT std()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT std(1)
< ok 0 0
> exec ^
: string
| SELECT stddev()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT stddev(1)
< ok 0 0
> exec ^
: string
| SELECT stddev_pop()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT stddev_pop(1)
< ok 0 0
> exec ^
: string
| SELECT stddev_samp()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT stddev_samp(1)
< ok 0 0
> exec ^
: string
| SELECT stop()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT stop(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT stop(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT stop(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT stop(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT stop(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT stop('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT stop(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT stop(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT stop(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT stop(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT stop(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT stop(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT stop() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT stop(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT stop(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT stop(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT stop(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT stop(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT stop(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT stop(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT stop(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT stop(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT stop(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT storage()
< error
| 1630 42000 FUNCTION test.storage does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT storage()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT stored()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT stored(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT stored(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT stored(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT stored(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT stored(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT stored('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT stored(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT stored(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT stored(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT stored(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT stored(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT stored(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT stored() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT stored(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT stored(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT stored(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT stored(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT stored(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT stored(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT stored(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT stored(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT stored(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT stored(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT straight_join()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT straight_join(1)
< ok 0 0
> exec ^
: string
| SELECT stream()
< error
| 1305 42000 FUNCTION test.stream does not exist
> exec ^
: string
| SELECT stream()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT strict_load()
< error
| 1305 42000 FUNCTION test.strict_load does not exist
> exec ^
: string
| SELECT strict_load()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT string()
< error
| 1630 42000 FUNCTION test.string does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT string()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT subclass_origin()
< error
| 1630 42000 FUNCTION test.subclass_origin does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT subclass_origin()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT subdate()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT subdate(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT subdate(1,2)
< ok 0 0
> exec ^
: string
| SELECT subject()
< error
| 1630 42000 FUNCTION test.subject does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT subject()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT subpartition()
< error
| 1630 42000 FUNCTION test.subpartition does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT subpartition()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT subpartitions()
< error
| 1630 42000 FUNCTION test.subpartitions does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT subpartitions()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT substr()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT substr(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT substr(1,2)
< ok 0 0
> exec ^
: string
| SELECT substring()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT substring(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT substring(1,2)
< ok 0 0
> exec ^
: string
| SELECT sum()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT sum(1)
< ok 0 0
> exec ^
: string
| SELECT super()
< error
| 1630 42000 FUNCTION test.super does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT super()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT suspend()
< error
| 1630 42000 FUNCTION test.suspend does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT suspend()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT swaps()
< error
| 1630 42000 FUNCTION test.swaps does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT swaps()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT switches()
< error
| 1630 42000 FUNCTION test.switches does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT switches()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT sysdate()
< ok 0 0
> exec ^
: string
| SELECT system()
< error
| 1630 42000 FUNCTION test.system does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT system()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT system_user()
< ok 0 0
> exec ^
: string
| SELECT table()
< error
| 1064 42000 …near 'table()' at line 1
> exec ^
: string
| SELECT table(1)
< error
| 1064 42000 …near 'table(1)' at line 1
> exec ^
: string
| SELECT table(1,2)
< error
| 1064 42000 …near 'table(1,2)' at line 1
> exec ^
: string
| SELECT table(1,2,3)
< error
| 1064 42000 …near 'table(1,2,3)' at line 1
> exec ^
: string
| SELECT table(1,2,3,4)
< error
| 1064 42000 …near 'table(1,2,3,4)' at line 1
> exec ^
: string
| SELECT table(*)
< error
| 1064 42000 …near 'table(*)' at line 1
> exec ^
: string
| SELECT table('a')
< error
| 1064 42000 …near 'table('a')' at line 1
> exec ^
: string
| SELECT table(a USING utf8mb4)
< error
| 1064 42000 …near 'table(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT table(1 AS SIGNED)
< error
| 1064 42000 …near 'table(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT table(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'table(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT table(1 FROM 2)
< error
| 1064 42000 …near 'table(1 FROM 2)' at line 1
> exec ^
: string
| SELECT table(1 IN 2)
< error
| 1064 42000 …near 'table(1 IN 2)' at line 1
> exec ^
: string
| SELECT table(DISTINCT 1)
< error
| 1064 42000 …near 'table(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT table() OVER ()
< error
| 1064 42000 …near 'table() OVER ()' at line 1
> exec ^
: string
| SELECT table(1) OVER ()
< error
| 1064 42000 …near 'table(1) OVER ()' at line 1
> exec ^
: string
| SELECT table(1,2) OVER ()
< error
| 1064 42000 …near 'table(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT table(1,2,3) OVER ()
< error
| 1064 42000 …near 'table(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT table(1) FROM t
< error
| 1064 42000 …near 'table(1) FROM t' at line 1
> exec ^
: string
| SELECT table(SELECT 1)
< error
| 1064 42000 …near 'table(SELECT 1)' at line 1
> exec ^
: string
| SELECT table(YEAR FROM 1)
< error
| 1064 42000 …near 'table(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT table(DAY, 1, 2)
< error
| 1064 42000 …near 'table(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT table(a)
< error
| 1064 42000 …near 'table(a)' at line 1
> exec ^
: string
| SELECT table(a, b)
< error
| 1064 42000 …near 'table(a, b)' at line 1
> exec ^
: string
| SELECT table(a) OVER ()
< error
| 1064 42000 …near 'table(a) OVER ()' at line 1
> exec ^
: string
| SELECT table_checksum()
< error
| 1630 42000 FUNCTION test.table_checksum does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT table_checksum()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT table_name()
< error
| 1630 42000 FUNCTION test.table_name does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT table_name()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT tables()
< error
| 1630 42000 FUNCTION test.tables does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT tables()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT tablesample()
< error
| 1305 42000 FUNCTION test.tablesample does not exist
> exec ^
: string
| SELECT tablesample()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT tablespace()
< error
| 1630 42000 FUNCTION test.tablespace does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT tablespace()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT temporary()
< error
| 1630 42000 FUNCTION test.temporary does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT temporary()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT temptable()
< error
| 1630 42000 FUNCTION test.temptable does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT temptable()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT terminated()
< error
| 1064 42000 …near 'terminated()' at line 1
> exec ^
: string
| SELECT terminated(1)
< error
| 1064 42000 …near 'terminated(1)' at line 1
> exec ^
: string
| SELECT terminated(1,2)
< error
| 1064 42000 …near 'terminated(1,2)' at line 1
> exec ^
: string
| SELECT terminated(1,2,3)
< error
| 1064 42000 …near 'terminated(1,2,3)' at line 1
> exec ^
: string
| SELECT terminated(1,2,3,4)
< error
| 1064 42000 …near 'terminated(1,2,3,4)' at line 1
> exec ^
: string
| SELECT terminated(*)
< error
| 1064 42000 …near 'terminated(*)' at line 1
> exec ^
: string
| SELECT terminated('a')
< error
| 1064 42000 …near 'terminated('a')' at line 1
> exec ^
: string
| SELECT terminated(a USING utf8mb4)
< error
| 1064 42000 …near 'terminated(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT terminated(1 AS SIGNED)
< error
| 1064 42000 …near 'terminated(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT terminated(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'terminated(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT terminated(1 FROM 2)
< error
| 1064 42000 …near 'terminated(1 FROM 2)' at line 1
> exec ^
: string
| SELECT terminated(1 IN 2)
< error
| 1064 42000 …near 'terminated(1 IN 2)' at line 1
> exec ^
: string
| SELECT terminated(DISTINCT 1)
< error
| 1064 42000 …near 'terminated(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT terminated() OVER ()
< error
| 1064 42000 …near 'terminated() OVER ()' at line 1
> exec ^
: string
| SELECT terminated(1) OVER ()
< error
| 1064 42000 …near 'terminated(1) OVER ()' at line 1
> exec ^
: string
| SELECT terminated(1,2) OVER ()
< error
| 1064 42000 …near 'terminated(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT terminated(1,2,3) OVER ()
< error
| 1064 42000 …near 'terminated(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT terminated(1) FROM t
< error
| 1064 42000 …near 'terminated(1) FROM t' at line 1
> exec ^
: string
| SELECT terminated(SELECT 1)
< error
| 1064 42000 …near 'terminated(SELECT 1)' at line 1
> exec ^
: string
| SELECT terminated(YEAR FROM 1)
< error
| 1064 42000 …near 'terminated(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT terminated(DAY, 1, 2)
< error
| 1064 42000 …near 'terminated(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT terminated(a)
< error
| 1064 42000 …near 'terminated(a)' at line 1
> exec ^
: string
| SELECT terminated(a, b)
< error
| 1064 42000 …near 'terminated(a, b)' at line 1
> exec ^
: string
| SELECT terminated(a) OVER ()
< error
| 1064 42000 …near 'terminated(a) OVER ()' at line 1
> exec ^
: string
| SELECT text()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT text(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT text(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT text(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT text(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT text(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT text('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT text(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT text(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT text(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT text(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT text(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT text(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT text() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT text(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT text(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT text(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT text(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT text(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT text(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT text(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT text(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT text(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT text(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT than()
< error
| 1630 42000 FUNCTION test.than does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT than()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT then()
< error
| 1064 42000 …near 'then()' at line 1
> exec ^
: string
| SELECT then(1)
< error
| 1064 42000 …near 'then(1)' at line 1
> exec ^
: string
| SELECT then(1,2)
< error
| 1064 42000 …near 'then(1,2)' at line 1
> exec ^
: string
| SELECT then(1,2,3)
< error
| 1064 42000 …near 'then(1,2,3)' at line 1
> exec ^
: string
| SELECT then(1,2,3,4)
< error
| 1064 42000 …near 'then(1,2,3,4)' at line 1
> exec ^
: string
| SELECT then(*)
< error
| 1064 42000 …near 'then(*)' at line 1
> exec ^
: string
| SELECT then('a')
< error
| 1064 42000 …near 'then('a')' at line 1
> exec ^
: string
| SELECT then(a USING utf8mb4)
< error
| 1064 42000 …near 'then(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT then(1 AS SIGNED)
< error
| 1064 42000 …near 'then(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT then(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'then(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT then(1 FROM 2)
< error
| 1064 42000 …near 'then(1 FROM 2)' at line 1
> exec ^
: string
| SELECT then(1 IN 2)
< error
| 1064 42000 …near 'then(1 IN 2)' at line 1
> exec ^
: string
| SELECT then(DISTINCT 1)
< error
| 1064 42000 …near 'then(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT then() OVER ()
< error
| 1064 42000 …near 'then() OVER ()' at line 1
> exec ^
: string
| SELECT then(1) OVER ()
< error
| 1064 42000 …near 'then(1) OVER ()' at line 1
> exec ^
: string
| SELECT then(1,2) OVER ()
< error
| 1064 42000 …near 'then(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT then(1,2,3) OVER ()
< error
| 1064 42000 …near 'then(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT then(1) FROM t
< error
| 1064 42000 …near 'then(1) FROM t' at line 1
> exec ^
: string
| SELECT then(SELECT 1)
< error
| 1064 42000 …near 'then(SELECT 1)' at line 1
> exec ^
: string
| SELECT then(YEAR FROM 1)
< error
| 1064 42000 …near 'then(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT then(DAY, 1, 2)
< error
| 1064 42000 …near 'then(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT then(a)
< error
| 1064 42000 …near 'then(a)' at line 1
> exec ^
: string
| SELECT then(a, b)
< error
| 1064 42000 …near 'then(a, b)' at line 1
> exec ^
: string
| SELECT then(a) OVER ()
< error
| 1064 42000 …near 'then(a) OVER ()' at line 1
> exec ^
: string
| SELECT thread_priority()
< error
| 1305 42000 FUNCTION test.thread_priority does not exist
> exec ^
: string
| SELECT thread_priority()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ties()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT ties(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT ties(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT ties(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT ties(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT ties(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT ties('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT ties(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT ties(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT ties(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT ties(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT ties(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT ties(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT ties() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT ties(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT ties(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT ties(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT ties(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT ties(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT ties(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT ties(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT ties(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT ties(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT ties(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT time()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT time(1)
< ok 0 0
> exec ^
: string
| SELECT timestamp()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT timestamp(1)
< ok 0 0
> exec ^
: string
| SELECT timestampadd()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT timestampadd(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT timestampadd(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT timestampadd(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT timestampadd(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT timestampadd(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT timestampadd('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT timestampadd(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT timestampadd(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT timestampadd(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT timestampadd(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT timestampadd(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT timestampadd(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT timestampadd() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT timestampadd(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT timestampadd(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT timestampadd(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT timestampadd(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT timestampadd(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT timestampadd(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT timestampadd(DAY, 1, 2)
< ok 0 0
> exec ^
: string
| SELECT timestampdiff()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT timestampdiff(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT timestampdiff(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT timestampdiff(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT timestampdiff(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT timestampdiff(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT timestampdiff('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT timestampdiff(a USING utf8mb4)
< error
| 1064 42000 …near 'a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT timestampdiff(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT timestampdiff(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT timestampdiff(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT timestampdiff(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT timestampdiff(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT timestampdiff() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT timestampdiff(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT timestampdiff(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT timestampdiff(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT timestampdiff(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT timestampdiff(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT timestampdiff(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT timestampdiff(DAY, 1, 2)
< ok 0 0
> exec ^
: string
| SELECT tinyblob()
< error
| 1064 42000 …near 'tinyblob()' at line 1
> exec ^
: string
| SELECT tinyblob(1)
< error
| 1064 42000 …near 'tinyblob(1)' at line 1
> exec ^
: string
| SELECT tinyblob(1,2)
< error
| 1064 42000 …near 'tinyblob(1,2)' at line 1
> exec ^
: string
| SELECT tinyblob(1,2,3)
< error
| 1064 42000 …near 'tinyblob(1,2,3)' at line 1
> exec ^
: string
| SELECT tinyblob(1,2,3,4)
< error
| 1064 42000 …near 'tinyblob(1,2,3,4)' at line 1
> exec ^
: string
| SELECT tinyblob(*)
< error
| 1064 42000 …near 'tinyblob(*)' at line 1
> exec ^
: string
| SELECT tinyblob('a')
< error
| 1064 42000 …near 'tinyblob('a')' at line 1
> exec ^
: string
| SELECT tinyblob(a USING utf8mb4)
< error
| 1064 42000 …near 'tinyblob(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT tinyblob(1 AS SIGNED)
< error
| 1064 42000 …near 'tinyblob(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT tinyblob(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'tinyblob(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT tinyblob(1 FROM 2)
< error
| 1064 42000 …near 'tinyblob(1 FROM 2)' at line 1
> exec ^
: string
| SELECT tinyblob(1 IN 2)
< error
| 1064 42000 …near 'tinyblob(1 IN 2)' at line 1
> exec ^
: string
| SELECT tinyblob(DISTINCT 1)
< error
| 1064 42000 …near 'tinyblob(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT tinyblob() OVER ()
< error
| 1064 42000 …near 'tinyblob() OVER ()' at line 1
> exec ^
: string
| SELECT tinyblob(1) OVER ()
< error
| 1064 42000 …near 'tinyblob(1) OVER ()' at line 1
> exec ^
: string
| SELECT tinyblob(1,2) OVER ()
< error
| 1064 42000 …near 'tinyblob(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT tinyblob(1,2,3) OVER ()
< error
| 1064 42000 …near 'tinyblob(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT tinyblob(1) FROM t
< error
| 1064 42000 …near 'tinyblob(1) FROM t' at line 1
> exec ^
: string
| SELECT tinyblob(SELECT 1)
< error
| 1064 42000 …near 'tinyblob(SELECT 1)' at line 1
> exec ^
: string
| SELECT tinyblob(YEAR FROM 1)
< error
| 1064 42000 …near 'tinyblob(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT tinyblob(DAY, 1, 2)
< error
| 1064 42000 …near 'tinyblob(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT tinyblob(a)
< error
| 1064 42000 …near 'tinyblob(a)' at line 1
> exec ^
: string
| SELECT tinyblob(a, b)
< error
| 1064 42000 …near 'tinyblob(a, b)' at line 1
> exec ^
: string
| SELECT tinyblob(a) OVER ()
< error
| 1064 42000 …near 'tinyblob(a) OVER ()' at line 1
> exec ^
: string
| SELECT tinyint()
< error
| 1064 42000 …near 'tinyint()' at line 1
> exec ^
: string
| SELECT tinyint(1)
< error
| 1064 42000 …near 'tinyint(1)' at line 1
> exec ^
: string
| SELECT tinyint(1,2)
< error
| 1064 42000 …near 'tinyint(1,2)' at line 1
> exec ^
: string
| SELECT tinyint(1,2,3)
< error
| 1064 42000 …near 'tinyint(1,2,3)' at line 1
> exec ^
: string
| SELECT tinyint(1,2,3,4)
< error
| 1064 42000 …near 'tinyint(1,2,3,4)' at line 1
> exec ^
: string
| SELECT tinyint(*)
< error
| 1064 42000 …near 'tinyint(*)' at line 1
> exec ^
: string
| SELECT tinyint('a')
< error
| 1064 42000 …near 'tinyint('a')' at line 1
> exec ^
: string
| SELECT tinyint(a USING utf8mb4)
< error
| 1064 42000 …near 'tinyint(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT tinyint(1 AS SIGNED)
< error
| 1064 42000 …near 'tinyint(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT tinyint(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'tinyint(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT tinyint(1 FROM 2)
< error
| 1064 42000 …near 'tinyint(1 FROM 2)' at line 1
> exec ^
: string
| SELECT tinyint(1 IN 2)
< error
| 1064 42000 …near 'tinyint(1 IN 2)' at line 1
> exec ^
: string
| SELECT tinyint(DISTINCT 1)
< error
| 1064 42000 …near 'tinyint(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT tinyint() OVER ()
< error
| 1064 42000 …near 'tinyint() OVER ()' at line 1
> exec ^
: string
| SELECT tinyint(1) OVER ()
< error
| 1064 42000 …near 'tinyint(1) OVER ()' at line 1
> exec ^
: string
| SELECT tinyint(1,2) OVER ()
< error
| 1064 42000 …near 'tinyint(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT tinyint(1,2,3) OVER ()
< error
| 1064 42000 …near 'tinyint(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT tinyint(1) FROM t
< error
| 1064 42000 …near 'tinyint(1) FROM t' at line 1
> exec ^
: string
| SELECT tinyint(SELECT 1)
< error
| 1064 42000 …near 'tinyint(SELECT 1)' at line 1
> exec ^
: string
| SELECT tinyint(YEAR FROM 1)
< error
| 1064 42000 …near 'tinyint(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT tinyint(DAY, 1, 2)
< error
| 1064 42000 …near 'tinyint(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT tinyint(a)
< error
| 1064 42000 …near 'tinyint(a)' at line 1
> exec ^
: string
| SELECT tinyint(a, b)
< error
| 1064 42000 …near 'tinyint(a, b)' at line 1
> exec ^
: string
| SELECT tinyint(a) OVER ()
< error
| 1064 42000 …near 'tinyint(a) OVER ()' at line 1
> exec ^
: string
| SELECT tinytext()
< error
| 1064 42000 …near 'tinytext()' at line 1
> exec ^
: string
| SELECT tinytext(1)
< error
| 1064 42000 …near 'tinytext(1)' at line 1
> exec ^
: string
| SELECT tinytext(1,2)
< error
| 1064 42000 …near 'tinytext(1,2)' at line 1
> exec ^
: string
| SELECT tinytext(1,2,3)
< error
| 1064 42000 …near 'tinytext(1,2,3)' at line 1
> exec ^
: string
| SELECT tinytext(1,2,3,4)
< error
| 1064 42000 …near 'tinytext(1,2,3,4)' at line 1
> exec ^
: string
| SELECT tinytext(*)
< error
| 1064 42000 …near 'tinytext(*)' at line 1
> exec ^
: string
| SELECT tinytext('a')
< error
| 1064 42000 …near 'tinytext('a')' at line 1
> exec ^
: string
| SELECT tinytext(a USING utf8mb4)
< error
| 1064 42000 …near 'tinytext(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT tinytext(1 AS SIGNED)
< error
| 1064 42000 …near 'tinytext(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT tinytext(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'tinytext(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT tinytext(1 FROM 2)
< error
| 1064 42000 …near 'tinytext(1 FROM 2)' at line 1
> exec ^
: string
| SELECT tinytext(1 IN 2)
< error
| 1064 42000 …near 'tinytext(1 IN 2)' at line 1
> exec ^
: string
| SELECT tinytext(DISTINCT 1)
< error
| 1064 42000 …near 'tinytext(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT tinytext() OVER ()
< error
| 1064 42000 …near 'tinytext() OVER ()' at line 1
> exec ^
: string
| SELECT tinytext(1) OVER ()
< error
| 1064 42000 …near 'tinytext(1) OVER ()' at line 1
> exec ^
: string
| SELECT tinytext(1,2) OVER ()
< error
| 1064 42000 …near 'tinytext(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT tinytext(1,2,3) OVER ()
< error
| 1064 42000 …near 'tinytext(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT tinytext(1) FROM t
< error
| 1064 42000 …near 'tinytext(1) FROM t' at line 1
> exec ^
: string
| SELECT tinytext(SELECT 1)
< error
| 1064 42000 …near 'tinytext(SELECT 1)' at line 1
> exec ^
: string
| SELECT tinytext(YEAR FROM 1)
< error
| 1064 42000 …near 'tinytext(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT tinytext(DAY, 1, 2)
< error
| 1064 42000 …near 'tinytext(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT tinytext(a)
< error
| 1064 42000 …near 'tinytext(a)' at line 1
> exec ^
: string
| SELECT tinytext(a, b)
< error
| 1064 42000 …near 'tinytext(a, b)' at line 1
> exec ^
: string
| SELECT tinytext(a) OVER ()
< error
| 1064 42000 …near 'tinytext(a) OVER ()' at line 1
> exec ^
: string
| SELECT tls()
< error
| 1305 42000 FUNCTION test.tls does not exist
> exec ^
: string
| SELECT tls()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT to()
< error
| 1064 42000 …near 'to()' at line 1
> exec ^
: string
| SELECT to(1)
< error
| 1064 42000 …near 'to(1)' at line 1
> exec ^
: string
| SELECT to(1,2)
< error
| 1064 42000 …near 'to(1,2)' at line 1
> exec ^
: string
| SELECT to(1,2,3)
< error
| 1064 42000 …near 'to(1,2,3)' at line 1
> exec ^
: string
| SELECT to(1,2,3,4)
< error
| 1064 42000 …near 'to(1,2,3,4)' at line 1
> exec ^
: string
| SELECT to(*)
< error
| 1064 42000 …near 'to(*)' at line 1
> exec ^
: string
| SELECT to('a')
< error
| 1064 42000 …near 'to('a')' at line 1
> exec ^
: string
| SELECT to(a USING utf8mb4)
< error
| 1064 42000 …near 'to(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT to(1 AS SIGNED)
< error
| 1064 42000 …near 'to(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT to(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'to(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT to(1 FROM 2)
< error
| 1064 42000 …near 'to(1 FROM 2)' at line 1
> exec ^
: string
| SELECT to(1 IN 2)
< error
| 1064 42000 …near 'to(1 IN 2)' at line 1
> exec ^
: string
| SELECT to(DISTINCT 1)
< error
| 1064 42000 …near 'to(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT to() OVER ()
< error
| 1064 42000 …near 'to() OVER ()' at line 1
> exec ^
: string
| SELECT to(1) OVER ()
< error
| 1064 42000 …near 'to(1) OVER ()' at line 1
> exec ^
: string
| SELECT to(1,2) OVER ()
< error
| 1064 42000 …near 'to(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT to(1,2,3) OVER ()
< error
| 1064 42000 …near 'to(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT to(1) FROM t
< error
| 1064 42000 …near 'to(1) FROM t' at line 1
> exec ^
: string
| SELECT to(SELECT 1)
< error
| 1064 42000 …near 'to(SELECT 1)' at line 1
> exec ^
: string
| SELECT to(YEAR FROM 1)
< error
| 1064 42000 …near 'to(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT to(DAY, 1, 2)
< error
| 1064 42000 …near 'to(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT to(a)
< error
| 1064 42000 …near 'to(a)' at line 1
> exec ^
: string
| SELECT to(a, b)
< error
| 1064 42000 …near 'to(a, b)' at line 1
> exec ^
: string
| SELECT to(a) OVER ()
< error
| 1064 42000 …near 'to(a) OVER ()' at line 1
> exec ^
: string
| SELECT trailing()
< error
| 1064 42000 …near 'trailing()' at line 1
> exec ^
: string
| SELECT trailing(1)
< error
| 1064 42000 …near 'trailing(1)' at line 1
> exec ^
: string
| SELECT trailing(1,2)
< error
| 1064 42000 …near 'trailing(1,2)' at line 1
> exec ^
: string
| SELECT trailing(1,2,3)
< error
| 1064 42000 …near 'trailing(1,2,3)' at line 1
> exec ^
: string
| SELECT trailing(1,2,3,4)
< error
| 1064 42000 …near 'trailing(1,2,3,4)' at line 1
> exec ^
: string
| SELECT trailing(*)
< error
| 1064 42000 …near 'trailing(*)' at line 1
> exec ^
: string
| SELECT trailing('a')
< error
| 1064 42000 …near 'trailing('a')' at line 1
> exec ^
: string
| SELECT trailing(a USING utf8mb4)
< error
| 1064 42000 …near 'trailing(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT trailing(1 AS SIGNED)
< error
| 1064 42000 …near 'trailing(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT trailing(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'trailing(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT trailing(1 FROM 2)
< error
| 1064 42000 …near 'trailing(1 FROM 2)' at line 1
> exec ^
: string
| SELECT trailing(1 IN 2)
< error
| 1064 42000 …near 'trailing(1 IN 2)' at line 1
> exec ^
: string
| SELECT trailing(DISTINCT 1)
< error
| 1064 42000 …near 'trailing(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT trailing() OVER ()
< error
| 1064 42000 …near 'trailing() OVER ()' at line 1
> exec ^
: string
| SELECT trailing(1) OVER ()
< error
| 1064 42000 …near 'trailing(1) OVER ()' at line 1
> exec ^
: string
| SELECT trailing(1,2) OVER ()
< error
| 1064 42000 …near 'trailing(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT trailing(1,2,3) OVER ()
< error
| 1064 42000 …near 'trailing(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT trailing(1) FROM t
< error
| 1064 42000 …near 'trailing(1) FROM t' at line 1
> exec ^
: string
| SELECT trailing(SELECT 1)
< error
| 1064 42000 …near 'trailing(SELECT 1)' at line 1
> exec ^
: string
| SELECT trailing(YEAR FROM 1)
< error
| 1064 42000 …near 'trailing(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT trailing(DAY, 1, 2)
< error
| 1064 42000 …near 'trailing(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT trailing(a)
< error
| 1064 42000 …near 'trailing(a)' at line 1
> exec ^
: string
| SELECT trailing(a, b)
< error
| 1064 42000 …near 'trailing(a, b)' at line 1
> exec ^
: string
| SELECT trailing(a) OVER ()
< error
| 1064 42000 …near 'trailing(a) OVER ()' at line 1
> exec ^
: string
| SELECT transaction()
< error
| 1630 42000 FUNCTION test.transaction does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT transaction()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT trigger()
< error
| 1064 42000 …near 'trigger()' at line 1
> exec ^
: string
| SELECT trigger(1)
< error
| 1064 42000 …near 'trigger(1)' at line 1
> exec ^
: string
| SELECT trigger(1,2)
< error
| 1064 42000 …near 'trigger(1,2)' at line 1
> exec ^
: string
| SELECT trigger(1,2,3)
< error
| 1064 42000 …near 'trigger(1,2,3)' at line 1
> exec ^
: string
| SELECT trigger(1,2,3,4)
< error
| 1064 42000 …near 'trigger(1,2,3,4)' at line 1
> exec ^
: string
| SELECT trigger(*)
< error
| 1064 42000 …near 'trigger(*)' at line 1
> exec ^
: string
| SELECT trigger('a')
< error
| 1064 42000 …near 'trigger('a')' at line 1
> exec ^
: string
| SELECT trigger(a USING utf8mb4)
< error
| 1064 42000 …near 'trigger(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT trigger(1 AS SIGNED)
< error
| 1064 42000 …near 'trigger(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT trigger(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'trigger(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT trigger(1 FROM 2)
< error
| 1064 42000 …near 'trigger(1 FROM 2)' at line 1
> exec ^
: string
| SELECT trigger(1 IN 2)
< error
| 1064 42000 …near 'trigger(1 IN 2)' at line 1
> exec ^
: string
| SELECT trigger(DISTINCT 1)
< error
| 1064 42000 …near 'trigger(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT trigger() OVER ()
< error
| 1064 42000 …near 'trigger() OVER ()' at line 1
> exec ^
: string
| SELECT trigger(1) OVER ()
< error
| 1064 42000 …near 'trigger(1) OVER ()' at line 1
> exec ^
: string
| SELECT trigger(1,2) OVER ()
< error
| 1064 42000 …near 'trigger(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT trigger(1,2,3) OVER ()
< error
| 1064 42000 …near 'trigger(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT trigger(1) FROM t
< error
| 1064 42000 …near 'trigger(1) FROM t' at line 1
> exec ^
: string
| SELECT trigger(SELECT 1)
< error
| 1064 42000 …near 'trigger(SELECT 1)' at line 1
> exec ^
: string
| SELECT trigger(YEAR FROM 1)
< error
| 1064 42000 …near 'trigger(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT trigger(DAY, 1, 2)
< error
| 1064 42000 …near 'trigger(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT trigger(a)
< error
| 1064 42000 …near 'trigger(a)' at line 1
> exec ^
: string
| SELECT trigger(a, b)
< error
| 1064 42000 …near 'trigger(a, b)' at line 1
> exec ^
: string
| SELECT trigger(a) OVER ()
< error
| 1064 42000 …near 'trigger(a) OVER ()' at line 1
> exec ^
: string
| SELECT triggers()
< error
| 1630 42000 FUNCTION test.triggers does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT triggers()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT trim()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT trim(1)
< ok 0 0
> exec ^
: string
| SELECT true()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT true(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT true(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT true(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT true(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT true(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT true('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT true(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT true(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT true(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT true(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT true(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT true(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT true() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT true(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT true(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT true(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT true(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT true(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT true(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT true(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT true(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT true(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT true(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT truncate()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT truncate(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT truncate(1,2)
< ok 0 0
> exec ^
: string
| SELECT type()
< error
| 1630 42000 FUNCTION test.type does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT type()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT types()
< error
| 1305 42000 FUNCTION test.types does not exist
> exec ^
: string
| SELECT types()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT unbounded()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT unbounded(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT unbounded(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT unbounded(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT unbounded(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT unbounded(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT unbounded('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT unbounded(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT unbounded(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT unbounded(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT unbounded(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT unbounded(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT unbounded(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT unbounded() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT unbounded(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT unbounded(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT unbounded(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT unbounded(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT unbounded(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT unbounded(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT unbounded(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT unbounded(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT unbounded(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT unbounded(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT uncommitted()
< error
| 1630 42000 FUNCTION test.uncommitted does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT uncommitted()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT undefined()
< error
| 1630 42000 FUNCTION test.undefined does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT undefined()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT undo()
< error
| 1064 42000 …near 'undo()' at line 1
> exec ^
: string
| SELECT undo(1)
< error
| 1064 42000 …near 'undo(1)' at line 1
> exec ^
: string
| SELECT undo(1,2)
< error
| 1064 42000 …near 'undo(1,2)' at line 1
> exec ^
: string
| SELECT undo(1,2,3)
< error
| 1064 42000 …near 'undo(1,2,3)' at line 1
> exec ^
: string
| SELECT undo(1,2,3,4)
< error
| 1064 42000 …near 'undo(1,2,3,4)' at line 1
> exec ^
: string
| SELECT undo(*)
< error
| 1064 42000 …near 'undo(*)' at line 1
> exec ^
: string
| SELECT undo('a')
< error
| 1064 42000 …near 'undo('a')' at line 1
> exec ^
: string
| SELECT undo(a USING utf8mb4)
< error
| 1064 42000 …near 'undo(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT undo(1 AS SIGNED)
< error
| 1064 42000 …near 'undo(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT undo(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'undo(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT undo(1 FROM 2)
< error
| 1064 42000 …near 'undo(1 FROM 2)' at line 1
> exec ^
: string
| SELECT undo(1 IN 2)
< error
| 1064 42000 …near 'undo(1 IN 2)' at line 1
> exec ^
: string
| SELECT undo(DISTINCT 1)
< error
| 1064 42000 …near 'undo(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT undo() OVER ()
< error
| 1064 42000 …near 'undo() OVER ()' at line 1
> exec ^
: string
| SELECT undo(1) OVER ()
< error
| 1064 42000 …near 'undo(1) OVER ()' at line 1
> exec ^
: string
| SELECT undo(1,2) OVER ()
< error
| 1064 42000 …near 'undo(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT undo(1,2,3) OVER ()
< error
| 1064 42000 …near 'undo(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT undo(1) FROM t
< error
| 1064 42000 …near 'undo(1) FROM t' at line 1
> exec ^
: string
| SELECT undo(SELECT 1)
< error
| 1064 42000 …near 'undo(SELECT 1)' at line 1
> exec ^
: string
| SELECT undo(YEAR FROM 1)
< error
| 1064 42000 …near 'undo(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT undo(DAY, 1, 2)
< error
| 1064 42000 …near 'undo(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT undo(a)
< error
| 1064 42000 …near 'undo(a)' at line 1
> exec ^
: string
| SELECT undo(a, b)
< error
| 1064 42000 …near 'undo(a, b)' at line 1
> exec ^
: string
| SELECT undo(a) OVER ()
< error
| 1064 42000 …near 'undo(a) OVER ()' at line 1
> exec ^
: string
| SELECT undo_buffer_size()
< error
| 1630 42000 FUNCTION test.undo_buffer_size does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT undo_buffer_size()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT undofile()
< error
| 1630 42000 FUNCTION test.undofile does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT undofile()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT unicode()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT unicode(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT unicode(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT unicode(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT unicode(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT unicode(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT unicode('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT unicode(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT unicode(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT unicode(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT unicode(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT unicode(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT unicode(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT unicode() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT unicode(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT unicode(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT unicode(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT unicode(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT unicode(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT unicode(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT unicode(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT unicode(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT unicode(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT unicode(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT uninstall()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT uninstall(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT uninstall(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT uninstall(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT uninstall(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT uninstall(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT uninstall('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT uninstall(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT uninstall(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT uninstall(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT uninstall(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT uninstall(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT uninstall(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT uninstall() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT uninstall(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT uninstall(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT uninstall(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT uninstall(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT uninstall(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT uninstall(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT uninstall(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT uninstall(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT uninstall(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT uninstall(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT union()
< error
| 1064 42000 …near 'union()' at line 1
> exec ^
: string
| SELECT union(1)
< error
| 1064 42000 …near 'union(1)' at line 1
> exec ^
: string
| SELECT union(1,2)
< error
| 1064 42000 …near 'union(1,2)' at line 1
> exec ^
: string
| SELECT union(1,2,3)
< error
| 1064 42000 …near 'union(1,2,3)' at line 1
> exec ^
: string
| SELECT union(1,2,3,4)
< error
| 1064 42000 …near 'union(1,2,3,4)' at line 1
> exec ^
: string
| SELECT union(*)
< error
| 1064 42000 …near 'union(*)' at line 1
> exec ^
: string
| SELECT union('a')
< error
| 1064 42000 …near 'union('a')' at line 1
> exec ^
: string
| SELECT union(a USING utf8mb4)
< error
| 1064 42000 …near 'union(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT union(1 AS SIGNED)
< error
| 1064 42000 …near 'union(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT union(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'union(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT union(1 FROM 2)
< error
| 1064 42000 …near 'union(1 FROM 2)' at line 1
> exec ^
: string
| SELECT union(1 IN 2)
< error
| 1064 42000 …near 'union(1 IN 2)' at line 1
> exec ^
: string
| SELECT union(DISTINCT 1)
< error
| 1064 42000 …near 'union(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT union() OVER ()
< error
| 1064 42000 …near 'union() OVER ()' at line 1
> exec ^
: string
| SELECT union(1) OVER ()
< error
| 1064 42000 …near 'union(1) OVER ()' at line 1
> exec ^
: string
| SELECT union(1,2) OVER ()
< error
| 1064 42000 …near 'union(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT union(1,2,3) OVER ()
< error
| 1064 42000 …near 'union(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT union(1) FROM t
< error
| 1064 42000 …near 'union(1) FROM t' at line 1
> exec ^
: string
| SELECT union(SELECT 1)
< error
| 1064 42000 …near 'union(SELECT 1)' at line 1
> exec ^
: string
| SELECT union(YEAR FROM 1)
< error
| 1064 42000 …near 'union(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT union(DAY, 1, 2)
< error
| 1064 42000 …near 'union(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT union(a)
< error
| 1064 42000 …near 'union(a)' at line 1
> exec ^
: string
| SELECT union(a, b)
< error
| 1064 42000 …near 'union(a, b)' at line 1
> exec ^
: string
| SELECT union(a) OVER ()
< error
| 1064 42000 …near 'union(a) OVER ()' at line 1
> exec ^
: string
| SELECT unique()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT unique(1)
< ok 0 0
> exec ^
: string
| SELECT unknown()
< error
| 1630 42000 FUNCTION test.unknown does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT unknown()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT unlock()
< error
| 1064 42000 …near 'unlock()' at line 1
> exec ^
: string
| SELECT unlock(1)
< error
| 1064 42000 …near 'unlock(1)' at line 1
> exec ^
: string
| SELECT unlock(1,2)
< error
| 1064 42000 …near 'unlock(1,2)' at line 1
> exec ^
: string
| SELECT unlock(1,2,3)
< error
| 1064 42000 …near 'unlock(1,2,3)' at line 1
> exec ^
: string
| SELECT unlock(1,2,3,4)
< error
| 1064 42000 …near 'unlock(1,2,3,4)' at line 1
> exec ^
: string
| SELECT unlock(*)
< error
| 1064 42000 …near 'unlock(*)' at line 1
> exec ^
: string
| SELECT unlock('a')
< error
| 1064 42000 …near 'unlock('a')' at line 1
> exec ^
: string
| SELECT unlock(a USING utf8mb4)
< error
| 1064 42000 …near 'unlock(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT unlock(1 AS SIGNED)
< error
| 1064 42000 …near 'unlock(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT unlock(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'unlock(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT unlock(1 FROM 2)
< error
| 1064 42000 …near 'unlock(1 FROM 2)' at line 1
> exec ^
: string
| SELECT unlock(1 IN 2)
< error
| 1064 42000 …near 'unlock(1 IN 2)' at line 1
> exec ^
: string
| SELECT unlock(DISTINCT 1)
< error
| 1064 42000 …near 'unlock(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT unlock() OVER ()
< error
| 1064 42000 …near 'unlock() OVER ()' at line 1
> exec ^
: string
| SELECT unlock(1) OVER ()
< error
| 1064 42000 …near 'unlock(1) OVER ()' at line 1
> exec ^
: string
| SELECT unlock(1,2) OVER ()
< error
| 1064 42000 …near 'unlock(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT unlock(1,2,3) OVER ()
< error
| 1064 42000 …near 'unlock(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT unlock(1) FROM t
< error
| 1064 42000 …near 'unlock(1) FROM t' at line 1
> exec ^
: string
| SELECT unlock(SELECT 1)
< error
| 1064 42000 …near 'unlock(SELECT 1)' at line 1
> exec ^
: string
| SELECT unlock(YEAR FROM 1)
< error
| 1064 42000 …near 'unlock(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT unlock(DAY, 1, 2)
< error
| 1064 42000 …near 'unlock(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT unlock(a)
< error
| 1064 42000 …near 'unlock(a)' at line 1
> exec ^
: string
| SELECT unlock(a, b)
< error
| 1064 42000 …near 'unlock(a, b)' at line 1
> exec ^
: string
| SELECT unlock(a) OVER ()
< error
| 1064 42000 …near 'unlock(a) OVER ()' at line 1
> exec ^
: string
| SELECT unregister()
< error
| 1305 42000 FUNCTION test.unregister does not exist
> exec ^
: string
| SELECT unregister()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT unsigned()
< error
| 1064 42000 …near 'unsigned()' at line 1
> exec ^
: string
| SELECT unsigned(1)
< error
| 1064 42000 …near 'unsigned(1)' at line 1
> exec ^
: string
| SELECT unsigned(1,2)
< error
| 1064 42000 …near 'unsigned(1,2)' at line 1
> exec ^
: string
| SELECT unsigned(1,2,3)
< error
| 1064 42000 …near 'unsigned(1,2,3)' at line 1
> exec ^
: string
| SELECT unsigned(1,2,3,4)
< error
| 1064 42000 …near 'unsigned(1,2,3,4)' at line 1
> exec ^
: string
| SELECT unsigned(*)
< error
| 1064 42000 …near 'unsigned(*)' at line 1
> exec ^
: string
| SELECT unsigned('a')
< error
| 1064 42000 …near 'unsigned('a')' at line 1
> exec ^
: string
| SELECT unsigned(a USING utf8mb4)
< error
| 1064 42000 …near 'unsigned(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT unsigned(1 AS SIGNED)
< error
| 1064 42000 …near 'unsigned(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT unsigned(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'unsigned(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT unsigned(1 FROM 2)
< error
| 1064 42000 …near 'unsigned(1 FROM 2)' at line 1
> exec ^
: string
| SELECT unsigned(1 IN 2)
< error
| 1064 42000 …near 'unsigned(1 IN 2)' at line 1
> exec ^
: string
| SELECT unsigned(DISTINCT 1)
< error
| 1064 42000 …near 'unsigned(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT unsigned() OVER ()
< error
| 1064 42000 …near 'unsigned() OVER ()' at line 1
> exec ^
: string
| SELECT unsigned(1) OVER ()
< error
| 1064 42000 …near 'unsigned(1) OVER ()' at line 1
> exec ^
: string
| SELECT unsigned(1,2) OVER ()
< error
| 1064 42000 …near 'unsigned(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT unsigned(1,2,3) OVER ()
< error
| 1064 42000 …near 'unsigned(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT unsigned(1) FROM t
< error
| 1064 42000 …near 'unsigned(1) FROM t' at line 1
> exec ^
: string
| SELECT unsigned(SELECT 1)
< error
| 1064 42000 …near 'unsigned(SELECT 1)' at line 1
> exec ^
: string
| SELECT unsigned(YEAR FROM 1)
< error
| 1064 42000 …near 'unsigned(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT unsigned(DAY, 1, 2)
< error
| 1064 42000 …near 'unsigned(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT unsigned(a)
< error
| 1064 42000 …near 'unsigned(a)' at line 1
> exec ^
: string
| SELECT unsigned(a, b)
< error
| 1064 42000 …near 'unsigned(a, b)' at line 1
> exec ^
: string
| SELECT unsigned(a) OVER ()
< error
| 1064 42000 …near 'unsigned(a) OVER ()' at line 1
> exec ^
: string
| SELECT until()
< error
| 1630 42000 FUNCTION test.until does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT until()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT update()
< error
| 1064 42000 …near 'update()' at line 1
> exec ^
: string
| SELECT update(1)
< error
| 1064 42000 …near 'update(1)' at line 1
> exec ^
: string
| SELECT update(1,2)
< error
| 1064 42000 …near 'update(1,2)' at line 1
> exec ^
: string
| SELECT update(1,2,3)
< error
| 1064 42000 …near 'update(1,2,3)' at line 1
> exec ^
: string
| SELECT update(1,2,3,4)
< error
| 1064 42000 …near 'update(1,2,3,4)' at line 1
> exec ^
: string
| SELECT update(*)
< error
| 1064 42000 …near 'update(*)' at line 1
> exec ^
: string
| SELECT update('a')
< error
| 1064 42000 …near 'update('a')' at line 1
> exec ^
: string
| SELECT update(a USING utf8mb4)
< error
| 1064 42000 …near 'update(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT update(1 AS SIGNED)
< error
| 1064 42000 …near 'update(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT update(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'update(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT update(1 FROM 2)
< error
| 1064 42000 …near 'update(1 FROM 2)' at line 1
> exec ^
: string
| SELECT update(1 IN 2)
< error
| 1064 42000 …near 'update(1 IN 2)' at line 1
> exec ^
: string
| SELECT update(DISTINCT 1)
< error
| 1064 42000 …near 'update(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT update() OVER ()
< error
| 1064 42000 …near 'update() OVER ()' at line 1
> exec ^
: string
| SELECT update(1) OVER ()
< error
| 1064 42000 …near 'update(1) OVER ()' at line 1
> exec ^
: string
| SELECT update(1,2) OVER ()
< error
| 1064 42000 …near 'update(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT update(1,2,3) OVER ()
< error
| 1064 42000 …near 'update(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT update(1) FROM t
< error
| 1064 42000 …near 'update(1) FROM t' at line 1
> exec ^
: string
| SELECT update(SELECT 1)
< error
| 1064 42000 …near 'update(SELECT 1)' at line 1
> exec ^
: string
| SELECT update(YEAR FROM 1)
< error
| 1064 42000 …near 'update(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT update(DAY, 1, 2)
< error
| 1064 42000 …near 'update(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT update(a)
< error
| 1064 42000 …near 'update(a)' at line 1
> exec ^
: string
| SELECT update(a, b)
< error
| 1064 42000 …near 'update(a, b)' at line 1
> exec ^
: string
| SELECT update(a) OVER ()
< error
| 1064 42000 …near 'update(a) OVER ()' at line 1
> exec ^
: string
| SELECT upgrade()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT upgrade(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT upgrade(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT upgrade(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT upgrade(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT upgrade(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT upgrade('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT upgrade(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT upgrade(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT upgrade(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT upgrade(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT upgrade(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT upgrade(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT upgrade() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT upgrade(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT upgrade(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT upgrade(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT upgrade(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT upgrade(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT upgrade(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT upgrade(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT upgrade(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT upgrade(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT upgrade(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT uri()
< error
| 1305 42000 FUNCTION test.uri does not exist
> exec ^
: string
| SELECT uri()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT url()
< error
| 1305 42000 FUNCTION test.url does not exist
> exec ^
: string
| SELECT url()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT usage()
< error
| 1064 42000 …near 'usage()' at line 1
> exec ^
: string
| SELECT usage(1)
< error
| 1064 42000 …near 'usage(1)' at line 1
> exec ^
: string
| SELECT usage(1,2)
< error
| 1064 42000 …near 'usage(1,2)' at line 1
> exec ^
: string
| SELECT usage(1,2,3)
< error
| 1064 42000 …near 'usage(1,2,3)' at line 1
> exec ^
: string
| SELECT usage(1,2,3,4)
< error
| 1064 42000 …near 'usage(1,2,3,4)' at line 1
> exec ^
: string
| SELECT usage(*)
< error
| 1064 42000 …near 'usage(*)' at line 1
> exec ^
: string
| SELECT usage('a')
< error
| 1064 42000 …near 'usage('a')' at line 1
> exec ^
: string
| SELECT usage(a USING utf8mb4)
< error
| 1064 42000 …near 'usage(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT usage(1 AS SIGNED)
< error
| 1064 42000 …near 'usage(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT usage(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'usage(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT usage(1 FROM 2)
< error
| 1064 42000 …near 'usage(1 FROM 2)' at line 1
> exec ^
: string
| SELECT usage(1 IN 2)
< error
| 1064 42000 …near 'usage(1 IN 2)' at line 1
> exec ^
: string
| SELECT usage(DISTINCT 1)
< error
| 1064 42000 …near 'usage(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT usage() OVER ()
< error
| 1064 42000 …near 'usage() OVER ()' at line 1
> exec ^
: string
| SELECT usage(1) OVER ()
< error
| 1064 42000 …near 'usage(1) OVER ()' at line 1
> exec ^
: string
| SELECT usage(1,2) OVER ()
< error
| 1064 42000 …near 'usage(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT usage(1,2,3) OVER ()
< error
| 1064 42000 …near 'usage(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT usage(1) FROM t
< error
| 1064 42000 …near 'usage(1) FROM t' at line 1
> exec ^
: string
| SELECT usage(SELECT 1)
< error
| 1064 42000 …near 'usage(SELECT 1)' at line 1
> exec ^
: string
| SELECT usage(YEAR FROM 1)
< error
| 1064 42000 …near 'usage(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT usage(DAY, 1, 2)
< error
| 1064 42000 …near 'usage(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT usage(a)
< error
| 1064 42000 …near 'usage(a)' at line 1
> exec ^
: string
| SELECT usage(a, b)
< error
| 1064 42000 …near 'usage(a, b)' at line 1
> exec ^
: string
| SELECT usage(a) OVER ()
< error
| 1064 42000 …near 'usage(a) OVER ()' at line 1
> exec ^
: string
| SELECT use()
< error
| 1064 42000 …near 'use()' at line 1
> exec ^
: string
| SELECT use(1)
< error
| 1064 42000 …near 'use(1)' at line 1
> exec ^
: string
| SELECT use(1,2)
< error
| 1064 42000 …near 'use(1,2)' at line 1
> exec ^
: string
| SELECT use(1,2,3)
< error
| 1064 42000 …near 'use(1,2,3)' at line 1
> exec ^
: string
| SELECT use(1,2,3,4)
< error
| 1064 42000 …near 'use(1,2,3,4)' at line 1
> exec ^
: string
| SELECT use(*)
< error
| 1064 42000 …near 'use(*)' at line 1
> exec ^
: string
| SELECT use('a')
< error
| 1064 42000 …near 'use('a')' at line 1
> exec ^
: string
| SELECT use(a USING utf8mb4)
< error
| 1064 42000 …near 'use(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT use(1 AS SIGNED)
< error
| 1064 42000 …near 'use(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT use(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'use(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT use(1 FROM 2)
< error
| 1064 42000 …near 'use(1 FROM 2)' at line 1
> exec ^
: string
| SELECT use(1 IN 2)
< error
| 1064 42000 …near 'use(1 IN 2)' at line 1
> exec ^
: string
| SELECT use(DISTINCT 1)
< error
| 1064 42000 …near 'use(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT use() OVER ()
< error
| 1064 42000 …near 'use() OVER ()' at line 1
> exec ^
: string
| SELECT use(1) OVER ()
< error
| 1064 42000 …near 'use(1) OVER ()' at line 1
> exec ^
: string
| SELECT use(1,2) OVER ()
< error
| 1064 42000 …near 'use(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT use(1,2,3) OVER ()
< error
| 1064 42000 …near 'use(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT use(1) FROM t
< error
| 1064 42000 …near 'use(1) FROM t' at line 1
> exec ^
: string
| SELECT use(SELECT 1)
< error
| 1064 42000 …near 'use(SELECT 1)' at line 1
> exec ^
: string
| SELECT use(YEAR FROM 1)
< error
| 1064 42000 …near 'use(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT use(DAY, 1, 2)
< error
| 1064 42000 …near 'use(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT use(a)
< error
| 1064 42000 …near 'use(a)' at line 1
> exec ^
: string
| SELECT use(a, b)
< error
| 1064 42000 …near 'use(a, b)' at line 1
> exec ^
: string
| SELECT use(a) OVER ()
< error
| 1064 42000 …near 'use(a) OVER ()' at line 1
> exec ^
: string
| SELECT use_frm()
< error
| 1630 42000 FUNCTION test.use_frm does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT use_frm()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT user()
< ok 0 0
> exec ^
: string
| SELECT user_resources()
< error
| 1630 42000 FUNCTION test.user_resources does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT user_resources()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT using()
< error
| 1064 42000 …near 'using()' at line 1
> exec ^
: string
| SELECT using(1)
< error
| 1064 42000 …near 'using(1)' at line 1
> exec ^
: string
| SELECT using(1,2)
< error
| 1064 42000 …near 'using(1,2)' at line 1
> exec ^
: string
| SELECT using(1,2,3)
< error
| 1064 42000 …near 'using(1,2,3)' at line 1
> exec ^
: string
| SELECT using(1,2,3,4)
< error
| 1064 42000 …near 'using(1,2,3,4)' at line 1
> exec ^
: string
| SELECT using(*)
< error
| 1064 42000 …near 'using(*)' at line 1
> exec ^
: string
| SELECT using('a')
< error
| 1064 42000 …near 'using('a')' at line 1
> exec ^
: string
| SELECT using(a USING utf8mb4)
< error
| 1064 42000 …near 'using(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT using(1 AS SIGNED)
< error
| 1064 42000 …near 'using(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT using(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'using(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT using(1 FROM 2)
< error
| 1064 42000 …near 'using(1 FROM 2)' at line 1
> exec ^
: string
| SELECT using(1 IN 2)
< error
| 1064 42000 …near 'using(1 IN 2)' at line 1
> exec ^
: string
| SELECT using(DISTINCT 1)
< error
| 1064 42000 …near 'using(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT using() OVER ()
< error
| 1064 42000 …near 'using() OVER ()' at line 1
> exec ^
: string
| SELECT using(1) OVER ()
< error
| 1064 42000 …near 'using(1) OVER ()' at line 1
> exec ^
: string
| SELECT using(1,2) OVER ()
< error
| 1064 42000 …near 'using(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT using(1,2,3) OVER ()
< error
| 1064 42000 …near 'using(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT using(1) FROM t
< error
| 1064 42000 …near 'using(1) FROM t' at line 1
> exec ^
: string
| SELECT using(SELECT 1)
< error
| 1064 42000 …near 'using(SELECT 1)' at line 1
> exec ^
: string
| SELECT using(YEAR FROM 1)
< error
| 1064 42000 …near 'using(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT using(DAY, 1, 2)
< error
| 1064 42000 …near 'using(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT using(a)
< error
| 1064 42000 …near 'using(a)' at line 1
> exec ^
: string
| SELECT using(a, b)
< error
| 1064 42000 …near 'using(a, b)' at line 1
> exec ^
: string
| SELECT using(a) OVER ()
< error
| 1064 42000 …near 'using(a) OVER ()' at line 1
> exec ^
: string
| SELECT utc_date()
< ok 0 0
> exec ^
: string
| SELECT utc_time()
< ok 0 0
> exec ^
: string
| SELECT utc_timestamp()
< ok 0 0
> exec ^
: string
| SELECT validate()
< error
| 1305 42000 FUNCTION test.validate does not exist
> exec ^
: string
| SELECT validate()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT validation()
< error
| 1630 42000 FUNCTION test.validation does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT validation()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT value()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT value(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT value(1,2)
< error
| 1064 42000 …near '1,2)' at line 1
> exec ^
: string
| SELECT value(1,2,3)
< error
| 1064 42000 …near '1,2,3)' at line 1
> exec ^
: string
| SELECT value(1,2,3,4)
< error
| 1064 42000 …near '1,2,3,4)' at line 1
> exec ^
: string
| SELECT value(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT value('a')
< error
| 1064 42000 …near ''a')' at line 1
> exec ^
: string
| SELECT value(a USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec ^
: string
| SELECT value(1 AS SIGNED)
< error
| 1064 42000 …near '1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT value(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT value(1 FROM 2)
< error
| 1064 42000 …near '1 FROM 2)' at line 1
> exec ^
: string
| SELECT value(1 IN 2)
< error
| 1064 42000 …near '1 IN 2)' at line 1
> exec ^
: string
| SELECT value(DISTINCT 1)
< error
| 1064 42000 …near 'DISTINCT 1)' at line 1
> exec ^
: string
| SELECT value() OVER ()
< error
| 1064 42000 …near ') OVER ()' at line 1
> exec ^
: string
| SELECT value(1) OVER ()
< error
| 1064 42000 …near '1) OVER ()' at line 1
> exec ^
: string
| SELECT value(1,2) OVER ()
< error
| 1064 42000 …near '1,2) OVER ()' at line 1
> exec ^
: string
| SELECT value(1,2,3) OVER ()
< error
| 1064 42000 …near '1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT value(1) FROM t
< error
| 1064 42000 …near '1) FROM t' at line 1
> exec ^
: string
| SELECT value(SELECT 1)
< error
| 1064 42000 …near 'SELECT 1)' at line 1
> exec ^
: string
| SELECT value(YEAR FROM 1)
< error
| 1064 42000 …near 'FROM 1)' at line 1
> exec ^
: string
| SELECT value(DAY, 1, 2)
< error
| 1064 42000 …near ' 1, 2)' at line 1
> exec ^
: string
| SELECT value(a)
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec ^
: string
| SELECT value(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT values()
< error
| 1064 42000 …near 'values()' at line 1
> exec ^
: string
| SELECT values(1)
< error
| 1064 42000 …near 'values(1)' at line 1
> exec ^
: string
| SELECT values(1,2)
< error
| 1064 42000 …near 'values(1,2)' at line 1
> exec ^
: string
| SELECT values(1,2,3)
< error
| 1064 42000 …near 'values(1,2,3)' at line 1
> exec ^
: string
| SELECT values(1,2,3,4)
< error
| 1064 42000 …near 'values(1,2,3,4)' at line 1
> exec ^
: string
| SELECT values(*)
< error
| 1064 42000 …near 'values(*)' at line 1
> exec ^
: string
| SELECT values('a')
< error
| 1064 42000 …near 'values('a')' at line 1
> exec ^
: string
| SELECT values(a USING utf8mb4)
< error
| 1064 42000 …near 'values(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT values(1 AS SIGNED)
< error
| 1064 42000 …near 'values(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT values(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'values(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT values(1 FROM 2)
< error
| 1064 42000 …near 'values(1 FROM 2)' at line 1
> exec ^
: string
| SELECT values(1 IN 2)
< error
| 1064 42000 …near 'values(1 IN 2)' at line 1
> exec ^
: string
| SELECT values(DISTINCT 1)
< error
| 1064 42000 …near 'values(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT values() OVER ()
< error
| 1064 42000 …near 'values() OVER ()' at line 1
> exec ^
: string
| SELECT values(1) OVER ()
< error
| 1064 42000 …near 'values(1) OVER ()' at line 1
> exec ^
: string
| SELECT values(1,2) OVER ()
< error
| 1064 42000 …near 'values(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT values(1,2,3) OVER ()
< error
| 1064 42000 …near 'values(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT values(1) FROM t
< error
| 1064 42000 …near 'values(1) FROM t' at line 1
> exec ^
: string
| SELECT values(SELECT 1)
< error
| 1064 42000 …near 'values(SELECT 1)' at line 1
> exec ^
: string
| SELECT values(YEAR FROM 1)
< error
| 1064 42000 …near 'values(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT values(DAY, 1, 2)
< error
| 1064 42000 …near 'values(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT values(a)
< error
| 1064 42000 …near 'values(a)' at line 1
> exec ^
: string
| SELECT values(a, b)
< error
| 1064 42000 …near 'values(a, b)' at line 1
> exec ^
: string
| SELECT values(a) OVER ()
< error
| 1064 42000 …near 'values(a) OVER ()' at line 1
> exec ^
: string
| SELECT var_pop()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT var_pop(1)
< ok 0 0
> exec ^
: string
| SELECT var_samp()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT var_samp(1)
< ok 0 0
> exec ^
: string
| SELECT varbinary()
< error
| 1064 42000 …near 'varbinary()' at line 1
> exec ^
: string
| SELECT varbinary(1)
< error
| 1064 42000 …near 'varbinary(1)' at line 1
> exec ^
: string
| SELECT varbinary(1,2)
< error
| 1064 42000 …near 'varbinary(1,2)' at line 1
> exec ^
: string
| SELECT varbinary(1,2,3)
< error
| 1064 42000 …near 'varbinary(1,2,3)' at line 1
> exec ^
: string
| SELECT varbinary(1,2,3,4)
< error
| 1064 42000 …near 'varbinary(1,2,3,4)' at line 1
> exec ^
: string
| SELECT varbinary(*)
< error
| 1064 42000 …near 'varbinary(*)' at line 1
> exec ^
: string
| SELECT varbinary('a')
< error
| 1064 42000 …near 'varbinary('a')' at line 1
> exec ^
: string
| SELECT varbinary(a USING utf8mb4)
< error
| 1064 42000 …near 'varbinary(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT varbinary(1 AS SIGNED)
< error
| 1064 42000 …near 'varbinary(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT varbinary(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'varbinary(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT varbinary(1 FROM 2)
< error
| 1064 42000 …near 'varbinary(1 FROM 2)' at line 1
> exec ^
: string
| SELECT varbinary(1 IN 2)
< error
| 1064 42000 …near 'varbinary(1 IN 2)' at line 1
> exec ^
: string
| SELECT varbinary(DISTINCT 1)
< error
| 1064 42000 …near 'varbinary(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT varbinary() OVER ()
< error
| 1064 42000 …near 'varbinary() OVER ()' at line 1
> exec ^
: string
| SELECT varbinary(1) OVER ()
< error
| 1064 42000 …near 'varbinary(1) OVER ()' at line 1
> exec ^
: string
| SELECT varbinary(1,2) OVER ()
< error
| 1064 42000 …near 'varbinary(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT varbinary(1,2,3) OVER ()
< error
| 1064 42000 …near 'varbinary(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT varbinary(1) FROM t
< error
| 1064 42000 …near 'varbinary(1) FROM t' at line 1
> exec ^
: string
| SELECT varbinary(SELECT 1)
< error
| 1064 42000 …near 'varbinary(SELECT 1)' at line 1
> exec ^
: string
| SELECT varbinary(YEAR FROM 1)
< error
| 1064 42000 …near 'varbinary(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT varbinary(DAY, 1, 2)
< error
| 1064 42000 …near 'varbinary(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT varbinary(a)
< error
| 1064 42000 …near 'varbinary(a)' at line 1
> exec ^
: string
| SELECT varbinary(a, b)
< error
| 1064 42000 …near 'varbinary(a, b)' at line 1
> exec ^
: string
| SELECT varbinary(a) OVER ()
< error
| 1064 42000 …near 'varbinary(a) OVER ()' at line 1
> exec ^
: string
| SELECT varchar()
< error
| 1064 42000 …near 'varchar()' at line 1
> exec ^
: string
| SELECT varchar(1)
< error
| 1064 42000 …near 'varchar(1)' at line 1
> exec ^
: string
| SELECT varchar(1,2)
< error
| 1064 42000 …near 'varchar(1,2)' at line 1
> exec ^
: string
| SELECT varchar(1,2,3)
< error
| 1064 42000 …near 'varchar(1,2,3)' at line 1
> exec ^
: string
| SELECT varchar(1,2,3,4)
< error
| 1064 42000 …near 'varchar(1,2,3,4)' at line 1
> exec ^
: string
| SELECT varchar(*)
< error
| 1064 42000 …near 'varchar(*)' at line 1
> exec ^
: string
| SELECT varchar('a')
< error
| 1064 42000 …near 'varchar('a')' at line 1
> exec ^
: string
| SELECT varchar(a USING utf8mb4)
< error
| 1064 42000 …near 'varchar(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT varchar(1 AS SIGNED)
< error
| 1064 42000 …near 'varchar(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT varchar(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'varchar(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT varchar(1 FROM 2)
< error
| 1064 42000 …near 'varchar(1 FROM 2)' at line 1
> exec ^
: string
| SELECT varchar(1 IN 2)
< error
| 1064 42000 …near 'varchar(1 IN 2)' at line 1
> exec ^
: string
| SELECT varchar(DISTINCT 1)
< error
| 1064 42000 …near 'varchar(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT varchar() OVER ()
< error
| 1064 42000 …near 'varchar() OVER ()' at line 1
> exec ^
: string
| SELECT varchar(1) OVER ()
< error
| 1064 42000 …near 'varchar(1) OVER ()' at line 1
> exec ^
: string
| SELECT varchar(1,2) OVER ()
< error
| 1064 42000 …near 'varchar(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT varchar(1,2,3) OVER ()
< error
| 1064 42000 …near 'varchar(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT varchar(1) FROM t
< error
| 1064 42000 …near 'varchar(1) FROM t' at line 1
> exec ^
: string
| SELECT varchar(SELECT 1)
< error
| 1064 42000 …near 'varchar(SELECT 1)' at line 1
> exec ^
: string
| SELECT varchar(YEAR FROM 1)
< error
| 1064 42000 …near 'varchar(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT varchar(DAY, 1, 2)
< error
| 1064 42000 …near 'varchar(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT varchar(a)
< error
| 1064 42000 …near 'varchar(a)' at line 1
> exec ^
: string
| SELECT varchar(a, b)
< error
| 1064 42000 …near 'varchar(a, b)' at line 1
> exec ^
: string
| SELECT varchar(a) OVER ()
< error
| 1064 42000 …near 'varchar(a) OVER ()' at line 1
> exec ^
: string
| SELECT varcharacter()
< error
| 1064 42000 …near 'varcharacter()' at line 1
> exec ^
: string
| SELECT varcharacter(1)
< error
| 1064 42000 …near 'varcharacter(1)' at line 1
> exec ^
: string
| SELECT varcharacter(1,2)
< error
| 1064 42000 …near 'varcharacter(1,2)' at line 1
> exec ^
: string
| SELECT varcharacter(1,2,3)
< error
| 1064 42000 …near 'varcharacter(1,2,3)' at line 1
> exec ^
: string
| SELECT varcharacter(1,2,3,4)
< error
| 1064 42000 …near 'varcharacter(1,2,3,4)' at line 1
> exec ^
: string
| SELECT varcharacter(*)
< error
| 1064 42000 …near 'varcharacter(*)' at line 1
> exec ^
: string
| SELECT varcharacter('a')
< error
| 1064 42000 …near 'varcharacter('a')' at line 1
> exec ^
: string
| SELECT varcharacter(a USING utf8mb4)
< error
| 1064 42000 …near 'varcharacter(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT varcharacter(1 AS SIGNED)
< error
| 1064 42000 …near 'varcharacter(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT varcharacter(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'varcharacter(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT varcharacter(1 FROM 2)
< error
| 1064 42000 …near 'varcharacter(1 FROM 2)' at line 1
> exec ^
: string
| SELECT varcharacter(1 IN 2)
< error
| 1064 42000 …near 'varcharacter(1 IN 2)' at line 1
> exec ^
: string
| SELECT varcharacter(DISTINCT 1)
< error
| 1064 42000 …near 'varcharacter(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT varcharacter() OVER ()
< error
| 1064 42000 …near 'varcharacter() OVER ()' at line 1
> exec ^
: string
| SELECT varcharacter(1) OVER ()
< error
| 1064 42000 …near 'varcharacter(1) OVER ()' at line 1
> exec ^
: string
| SELECT varcharacter(1,2) OVER ()
< error
| 1064 42000 …near 'varcharacter(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT varcharacter(1,2,3) OVER ()
< error
| 1064 42000 …near 'varcharacter(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT varcharacter(1) FROM t
< error
| 1064 42000 …near 'varcharacter(1) FROM t' at line 1
> exec ^
: string
| SELECT varcharacter(SELECT 1)
< error
| 1064 42000 …near 'varcharacter(SELECT 1)' at line 1
> exec ^
: string
| SELECT varcharacter(YEAR FROM 1)
< error
| 1064 42000 …near 'varcharacter(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT varcharacter(DAY, 1, 2)
< error
| 1064 42000 …near 'varcharacter(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT varcharacter(a)
< error
| 1064 42000 …near 'varcharacter(a)' at line 1
> exec ^
: string
| SELECT varcharacter(a, b)
< error
| 1064 42000 …near 'varcharacter(a, b)' at line 1
> exec ^
: string
| SELECT varcharacter(a) OVER ()
< error
| 1064 42000 …near 'varcharacter(a) OVER ()' at line 1
> exec ^
: string
| SELECT variables()
< error
| 1630 42000 FUNCTION test.variables does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT variables()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT variance()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT variance(1)
< ok 0 0
> exec ^
: string
| SELECT varying()
< error
| 1064 42000 …near 'varying()' at line 1
> exec ^
: string
| SELECT varying(1)
< error
| 1064 42000 …near 'varying(1)' at line 1
> exec ^
: string
| SELECT varying(1,2)
< error
| 1064 42000 …near 'varying(1,2)' at line 1
> exec ^
: string
| SELECT varying(1,2,3)
< error
| 1064 42000 …near 'varying(1,2,3)' at line 1
> exec ^
: string
| SELECT varying(1,2,3,4)
< error
| 1064 42000 …near 'varying(1,2,3,4)' at line 1
> exec ^
: string
| SELECT varying(*)
< error
| 1064 42000 …near 'varying(*)' at line 1
> exec ^
: string
| SELECT varying('a')
< error
| 1064 42000 …near 'varying('a')' at line 1
> exec ^
: string
| SELECT varying(a USING utf8mb4)
< error
| 1064 42000 …near 'varying(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT varying(1 AS SIGNED)
< error
| 1064 42000 …near 'varying(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT varying(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'varying(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT varying(1 FROM 2)
< error
| 1064 42000 …near 'varying(1 FROM 2)' at line 1
> exec ^
: string
| SELECT varying(1 IN 2)
< error
| 1064 42000 …near 'varying(1 IN 2)' at line 1
> exec ^
: string
| SELECT varying(DISTINCT 1)
< error
| 1064 42000 …near 'varying(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT varying() OVER ()
< error
| 1064 42000 …near 'varying() OVER ()' at line 1
> exec ^
: string
| SELECT varying(1) OVER ()
< error
| 1064 42000 …near 'varying(1) OVER ()' at line 1
> exec ^
: string
| SELECT varying(1,2) OVER ()
< error
| 1064 42000 …near 'varying(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT varying(1,2,3) OVER ()
< error
| 1064 42000 …near 'varying(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT varying(1) FROM t
< error
| 1064 42000 …near 'varying(1) FROM t' at line 1
> exec ^
: string
| SELECT varying(SELECT 1)
< error
| 1064 42000 …near 'varying(SELECT 1)' at line 1
> exec ^
: string
| SELECT varying(YEAR FROM 1)
< error
| 1064 42000 …near 'varying(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT varying(DAY, 1, 2)
< error
| 1064 42000 …near 'varying(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT varying(a)
< error
| 1064 42000 …near 'varying(a)' at line 1
> exec ^
: string
| SELECT varying(a, b)
< error
| 1064 42000 …near 'varying(a, b)' at line 1
> exec ^
: string
| SELECT varying(a) OVER ()
< error
| 1064 42000 …near 'varying(a) OVER ()' at line 1
> exec ^
: string
| SELECT vcpu()
< error
| 1305 42000 FUNCTION test.vcpu does not exist
> exec ^
: string
| SELECT vcpu()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT vector()
< error
| 1064 42000 …near 'vector()' at line 1
> exec ^
: string
| SELECT vector(1)
< error
| 1064 42000 …near 'vector(1)' at line 1
> exec ^
: string
| SELECT vector(1,2)
< error
| 1064 42000 …near 'vector(1,2)' at line 1
> exec ^
: string
| SELECT vector(1,2,3)
< error
| 1064 42000 …near 'vector(1,2,3)' at line 1
> exec ^
: string
| SELECT vector(1,2,3,4)
< error
| 1064 42000 …near 'vector(1,2,3,4)' at line 1
> exec ^
: string
| SELECT vector(*)
< error
| 1064 42000 …near 'vector(*)' at line 1
> exec ^
: string
| SELECT vector('a')
< error
| 1064 42000 …near 'vector('a')' at line 1
> exec ^
: string
| SELECT vector(a USING utf8mb4)
< error
| 1064 42000 …near 'vector(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT vector(1 AS SIGNED)
< error
| 1064 42000 …near 'vector(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT vector(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'vector(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT vector(1 FROM 2)
< error
| 1064 42000 …near 'vector(1 FROM 2)' at line 1
> exec ^
: string
| SELECT vector(1 IN 2)
< error
| 1064 42000 …near 'vector(1 IN 2)' at line 1
> exec ^
: string
| SELECT vector(DISTINCT 1)
< error
| 1064 42000 …near 'vector(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT vector() OVER ()
< error
| 1064 42000 …near 'vector() OVER ()' at line 1
> exec ^
: string
| SELECT vector(1) OVER ()
< error
| 1064 42000 …near 'vector(1) OVER ()' at line 1
> exec ^
: string
| SELECT vector(1,2) OVER ()
< error
| 1064 42000 …near 'vector(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT vector(1,2,3) OVER ()
< error
| 1064 42000 …near 'vector(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT vector(1) FROM t
< error
| 1064 42000 …near 'vector(1) FROM t' at line 1
> exec ^
: string
| SELECT vector(SELECT 1)
< error
| 1064 42000 …near 'vector(SELECT 1)' at line 1
> exec ^
: string
| SELECT vector(YEAR FROM 1)
< error
| 1064 42000 …near 'vector(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT vector(DAY, 1, 2)
< error
| 1064 42000 …near 'vector(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT vector(a)
< error
| 1064 42000 …near 'vector(a)' at line 1
> exec ^
: string
| SELECT vector(a, b)
< error
| 1064 42000 …near 'vector(a, b)' at line 1
> exec ^
: string
| SELECT vector(a) OVER ()
< error
| 1064 42000 …near 'vector(a) OVER ()' at line 1
> exec ^
: string
| SELECT verify_key_constraints()
< error
| 1305 42000 FUNCTION test.verify_key_constraints does not exist
> exec ^
: string
| SELECT verify_key_constraints()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT view()
< error
| 1630 42000 FUNCTION test.view does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT view()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT virtual()
< error
| 1630 42000 FUNCTION test.virtual does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT virtual()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT visible()
< error
| 1630 42000 FUNCTION test.visible does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT visible()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT wait()
< error
| 1630 42000 FUNCTION test.wait does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT wait()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT warnings()
< error
| 1630 42000 FUNCTION test.warnings does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT warnings()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT week()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'week'
> exec ^
: string
| SELECT week()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'week'
> exec ^
: string
| SELECT week(1)
< ok 0 0
> exec ^
: string
| SELECT weight_string()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT weight_string(1)
< ok 0 0
> exec ^
: string
| SELECT when()
< error
| 1064 42000 …near 'when()' at line 1
> exec ^
: string
| SELECT when(1)
< error
| 1064 42000 …near 'when(1)' at line 1
> exec ^
: string
| SELECT when(1,2)
< error
| 1064 42000 …near 'when(1,2)' at line 1
> exec ^
: string
| SELECT when(1,2,3)
< error
| 1064 42000 …near 'when(1,2,3)' at line 1
> exec ^
: string
| SELECT when(1,2,3,4)
< error
| 1064 42000 …near 'when(1,2,3,4)' at line 1
> exec ^
: string
| SELECT when(*)
< error
| 1064 42000 …near 'when(*)' at line 1
> exec ^
: string
| SELECT when('a')
< error
| 1064 42000 …near 'when('a')' at line 1
> exec ^
: string
| SELECT when(a USING utf8mb4)
< error
| 1064 42000 …near 'when(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT when(1 AS SIGNED)
< error
| 1064 42000 …near 'when(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT when(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'when(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT when(1 FROM 2)
< error
| 1064 42000 …near 'when(1 FROM 2)' at line 1
> exec ^
: string
| SELECT when(1 IN 2)
< error
| 1064 42000 …near 'when(1 IN 2)' at line 1
> exec ^
: string
| SELECT when(DISTINCT 1)
< error
| 1064 42000 …near 'when(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT when() OVER ()
< error
| 1064 42000 …near 'when() OVER ()' at line 1
> exec ^
: string
| SELECT when(1) OVER ()
< error
| 1064 42000 …near 'when(1) OVER ()' at line 1
> exec ^
: string
| SELECT when(1,2) OVER ()
< error
| 1064 42000 …near 'when(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT when(1,2,3) OVER ()
< error
| 1064 42000 …near 'when(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT when(1) FROM t
< error
| 1064 42000 …near 'when(1) FROM t' at line 1
> exec ^
: string
| SELECT when(SELECT 1)
< error
| 1064 42000 …near 'when(SELECT 1)' at line 1
> exec ^
: string
| SELECT when(YEAR FROM 1)
< error
| 1064 42000 …near 'when(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT when(DAY, 1, 2)
< error
| 1064 42000 …near 'when(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT when(a)
< error
| 1064 42000 …near 'when(a)' at line 1
> exec ^
: string
| SELECT when(a, b)
< error
| 1064 42000 …near 'when(a, b)' at line 1
> exec ^
: string
| SELECT when(a) OVER ()
< error
| 1064 42000 …near 'when(a) OVER ()' at line 1
> exec ^
: string
| SELECT where()
< error
| 1064 42000 …near 'where()' at line 1
> exec ^
: string
| SELECT where(1)
< error
| 1064 42000 …near 'where(1)' at line 1
> exec ^
: string
| SELECT where(1,2)
< error
| 1064 42000 …near 'where(1,2)' at line 1
> exec ^
: string
| SELECT where(1,2,3)
< error
| 1064 42000 …near 'where(1,2,3)' at line 1
> exec ^
: string
| SELECT where(1,2,3,4)
< error
| 1064 42000 …near 'where(1,2,3,4)' at line 1
> exec ^
: string
| SELECT where(*)
< error
| 1064 42000 …near 'where(*)' at line 1
> exec ^
: string
| SELECT where('a')
< error
| 1064 42000 …near 'where('a')' at line 1
> exec ^
: string
| SELECT where(a USING utf8mb4)
< error
| 1064 42000 …near 'where(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT where(1 AS SIGNED)
< error
| 1064 42000 …near 'where(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT where(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'where(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT where(1 FROM 2)
< error
| 1064 42000 …near 'where(1 FROM 2)' at line 1
> exec ^
: string
| SELECT where(1 IN 2)
< error
| 1064 42000 …near 'where(1 IN 2)' at line 1
> exec ^
: string
| SELECT where(DISTINCT 1)
< error
| 1064 42000 …near 'where(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT where() OVER ()
< error
| 1064 42000 …near 'where() OVER ()' at line 1
> exec ^
: string
| SELECT where(1) OVER ()
< error
| 1064 42000 …near 'where(1) OVER ()' at line 1
> exec ^
: string
| SELECT where(1,2) OVER ()
< error
| 1064 42000 …near 'where(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT where(1,2,3) OVER ()
< error
| 1064 42000 …near 'where(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT where(1) FROM t
< error
| 1064 42000 …near 'where(1) FROM t' at line 1
> exec ^
: string
| SELECT where(SELECT 1)
< error
| 1064 42000 …near 'where(SELECT 1)' at line 1
> exec ^
: string
| SELECT where(YEAR FROM 1)
< error
| 1064 42000 …near 'where(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT where(DAY, 1, 2)
< error
| 1064 42000 …near 'where(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT where(a)
< error
| 1064 42000 …near 'where(a)' at line 1
> exec ^
: string
| SELECT where(a, b)
< error
| 1064 42000 …near 'where(a, b)' at line 1
> exec ^
: string
| SELECT where(a) OVER ()
< error
| 1064 42000 …near 'where(a) OVER ()' at line 1
> exec ^
: string
| SELECT while()
< error
| 1064 42000 …near 'while()' at line 1
> exec ^
: string
| SELECT while(1)
< error
| 1064 42000 …near 'while(1)' at line 1
> exec ^
: string
| SELECT while(1,2)
< error
| 1064 42000 …near 'while(1,2)' at line 1
> exec ^
: string
| SELECT while(1,2,3)
< error
| 1064 42000 …near 'while(1,2,3)' at line 1
> exec ^
: string
| SELECT while(1,2,3,4)
< error
| 1064 42000 …near 'while(1,2,3,4)' at line 1
> exec ^
: string
| SELECT while(*)
< error
| 1064 42000 …near 'while(*)' at line 1
> exec ^
: string
| SELECT while('a')
< error
| 1064 42000 …near 'while('a')' at line 1
> exec ^
: string
| SELECT while(a USING utf8mb4)
< error
| 1064 42000 …near 'while(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT while(1 AS SIGNED)
< error
| 1064 42000 …near 'while(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT while(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'while(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT while(1 FROM 2)
< error
| 1064 42000 …near 'while(1 FROM 2)' at line 1
> exec ^
: string
| SELECT while(1 IN 2)
< error
| 1064 42000 …near 'while(1 IN 2)' at line 1
> exec ^
: string
| SELECT while(DISTINCT 1)
< error
| 1064 42000 …near 'while(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT while() OVER ()
< error
| 1064 42000 …near 'while() OVER ()' at line 1
> exec ^
: string
| SELECT while(1) OVER ()
< error
| 1064 42000 …near 'while(1) OVER ()' at line 1
> exec ^
: string
| SELECT while(1,2) OVER ()
< error
| 1064 42000 …near 'while(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT while(1,2,3) OVER ()
< error
| 1064 42000 …near 'while(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT while(1) FROM t
< error
| 1064 42000 …near 'while(1) FROM t' at line 1
> exec ^
: string
| SELECT while(SELECT 1)
< error
| 1064 42000 …near 'while(SELECT 1)' at line 1
> exec ^
: string
| SELECT while(YEAR FROM 1)
< error
| 1064 42000 …near 'while(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT while(DAY, 1, 2)
< error
| 1064 42000 …near 'while(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT while(a)
< error
| 1064 42000 …near 'while(a)' at line 1
> exec ^
: string
| SELECT while(a, b)
< error
| 1064 42000 …near 'while(a, b)' at line 1
> exec ^
: string
| SELECT while(a) OVER ()
< error
| 1064 42000 …near 'while(a) OVER ()' at line 1
> exec ^
: string
| SELECT window()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT window(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT window(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT window(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT window(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT window(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT window('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT window(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT window(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT window(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT window(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT window(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT window(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT window() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT window(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT window(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT window(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT window(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT window(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT window(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT window(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT window(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT window(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT window(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT with()
< error
| 1064 42000 …near 'with()' at line 1
> exec ^
: string
| SELECT with(1)
< error
| 1064 42000 …near 'with(1)' at line 1
> exec ^
: string
| SELECT with(1,2)
< error
| 1064 42000 …near 'with(1,2)' at line 1
> exec ^
: string
| SELECT with(1,2,3)
< error
| 1064 42000 …near 'with(1,2,3)' at line 1
> exec ^
: string
| SELECT with(1,2,3,4)
< error
| 1064 42000 …near 'with(1,2,3,4)' at line 1
> exec ^
: string
| SELECT with(*)
< error
| 1064 42000 …near 'with(*)' at line 1
> exec ^
: string
| SELECT with('a')
< error
| 1064 42000 …near 'with('a')' at line 1
> exec ^
: string
| SELECT with(a USING utf8mb4)
< error
| 1064 42000 …near 'with(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT with(1 AS SIGNED)
< error
| 1064 42000 …near 'with(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT with(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'with(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT with(1 FROM 2)
< error
| 1064 42000 …near 'with(1 FROM 2)' at line 1
> exec ^
: string
| SELECT with(1 IN 2)
< error
| 1064 42000 …near 'with(1 IN 2)' at line 1
> exec ^
: string
| SELECT with(DISTINCT 1)
< error
| 1064 42000 …near 'with(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT with() OVER ()
< error
| 1064 42000 …near 'with() OVER ()' at line 1
> exec ^
: string
| SELECT with(1) OVER ()
< error
| 1064 42000 …near 'with(1) OVER ()' at line 1
> exec ^
: string
| SELECT with(1,2) OVER ()
< error
| 1064 42000 …near 'with(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT with(1,2,3) OVER ()
< error
| 1064 42000 …near 'with(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT with(1) FROM t
< error
| 1064 42000 …near 'with(1) FROM t' at line 1
> exec ^
: string
| SELECT with(SELECT 1)
< error
| 1064 42000 …near 'with(SELECT 1)' at line 1
> exec ^
: string
| SELECT with(YEAR FROM 1)
< error
| 1064 42000 …near 'with(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT with(DAY, 1, 2)
< error
| 1064 42000 …near 'with(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT with(a)
< error
| 1064 42000 …near 'with(a)' at line 1
> exec ^
: string
| SELECT with(a, b)
< error
| 1064 42000 …near 'with(a, b)' at line 1
> exec ^
: string
| SELECT with(a) OVER ()
< error
| 1064 42000 …near 'with(a) OVER ()' at line 1
> exec ^
: string
| SELECT without()
< error
| 1630 42000 FUNCTION test.without does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT without()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT work()
< error
| 1630 42000 FUNCTION test.work does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT work()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT wrapper()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT wrapper(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT wrapper(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT wrapper(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT wrapper(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT wrapper(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT wrapper('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT wrapper(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT wrapper(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT wrapper(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT wrapper(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT wrapper(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT wrapper(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT wrapper() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT wrapper(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT wrapper(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT wrapper(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT wrapper(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT wrapper(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT wrapper(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT wrapper(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT wrapper(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT wrapper(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT wrapper(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT write()
< error
| 1064 42000 …near 'write()' at line 1
> exec ^
: string
| SELECT write(1)
< error
| 1064 42000 …near 'write(1)' at line 1
> exec ^
: string
| SELECT write(1,2)
< error
| 1064 42000 …near 'write(1,2)' at line 1
> exec ^
: string
| SELECT write(1,2,3)
< error
| 1064 42000 …near 'write(1,2,3)' at line 1
> exec ^
: string
| SELECT write(1,2,3,4)
< error
| 1064 42000 …near 'write(1,2,3,4)' at line 1
> exec ^
: string
| SELECT write(*)
< error
| 1064 42000 …near 'write(*)' at line 1
> exec ^
: string
| SELECT write('a')
< error
| 1064 42000 …near 'write('a')' at line 1
> exec ^
: string
| SELECT write(a USING utf8mb4)
< error
| 1064 42000 …near 'write(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT write(1 AS SIGNED)
< error
| 1064 42000 …near 'write(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT write(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'write(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT write(1 FROM 2)
< error
| 1064 42000 …near 'write(1 FROM 2)' at line 1
> exec ^
: string
| SELECT write(1 IN 2)
< error
| 1064 42000 …near 'write(1 IN 2)' at line 1
> exec ^
: string
| SELECT write(DISTINCT 1)
< error
| 1064 42000 …near 'write(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT write() OVER ()
< error
| 1064 42000 …near 'write() OVER ()' at line 1
> exec ^
: string
| SELECT write(1) OVER ()
< error
| 1064 42000 …near 'write(1) OVER ()' at line 1
> exec ^
: string
| SELECT write(1,2) OVER ()
< error
| 1064 42000 …near 'write(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT write(1,2,3) OVER ()
< error
| 1064 42000 …near 'write(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT write(1) FROM t
< error
| 1064 42000 …near 'write(1) FROM t' at line 1
> exec ^
: string
| SELECT write(SELECT 1)
< error
| 1064 42000 …near 'write(SELECT 1)' at line 1
> exec ^
: string
| SELECT write(YEAR FROM 1)
< error
| 1064 42000 …near 'write(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT write(DAY, 1, 2)
< error
| 1064 42000 …near 'write(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT write(a)
< error
| 1064 42000 …near 'write(a)' at line 1
> exec ^
: string
| SELECT write(a, b)
< error
| 1064 42000 …near 'write(a, b)' at line 1
> exec ^
: string
| SELECT write(a) OVER ()
< error
| 1064 42000 …near 'write(a) OVER ()' at line 1
> exec ^
: string
| SELECT x509()
< error
| 1630 42000 FUNCTION test.x509 does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT x509()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT xa()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| SELECT xa(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT xa(1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| SELECT xa(1,2,3)
< error
| 1064 42000 …near '(1,2,3)' at line 1
> exec ^
: string
| SELECT xa(1,2,3,4)
< error
| 1064 42000 …near '(1,2,3,4)' at line 1
> exec ^
: string
| SELECT xa(*)
< error
| 1064 42000 …near '(*)' at line 1
> exec ^
: string
| SELECT xa('a')
< error
| 1064 42000 …near '('a')' at line 1
> exec ^
: string
| SELECT xa(a USING utf8mb4)
< error
| 1064 42000 …near '(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT xa(1 AS SIGNED)
< error
| 1064 42000 …near '(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT xa(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near '(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT xa(1 FROM 2)
< error
| 1064 42000 …near '(1 FROM 2)' at line 1
> exec ^
: string
| SELECT xa(1 IN 2)
< error
| 1064 42000 …near '(1 IN 2)' at line 1
> exec ^
: string
| SELECT xa(DISTINCT 1)
< error
| 1064 42000 …near '(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT xa() OVER ()
< error
| 1064 42000 …near '() OVER ()' at line 1
> exec ^
: string
| SELECT xa(1) OVER ()
< error
| 1064 42000 …near '(1) OVER ()' at line 1
> exec ^
: string
| SELECT xa(1,2) OVER ()
< error
| 1064 42000 …near '(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT xa(1,2,3) OVER ()
< error
| 1064 42000 …near '(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT xa(1) FROM t
< error
| 1064 42000 …near '(1) FROM t' at line 1
> exec ^
: string
| SELECT xa(SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT xa(YEAR FROM 1)
< error
| 1064 42000 …near '(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT xa(DAY, 1, 2)
< error
| 1064 42000 …near '(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT xa(a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| SELECT xa(a, b)
< error
| 1064 42000 …near '(a, b)' at line 1
> exec ^
: string
| SELECT xa(a) OVER ()
< error
| 1064 42000 …near '(a) OVER ()' at line 1
> exec ^
: string
| SELECT xid()
< error
| 1305 42000 FUNCTION test.xid does not exist
> exec ^
: string
| SELECT xid()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT xml()
< error
| 1630 42000 FUNCTION test.xml does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT xml()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT xor()
< error
| 1064 42000 …near 'xor()' at line 1
> exec ^
: string
| SELECT xor(1)
< error
| 1064 42000 …near 'xor(1)' at line 1
> exec ^
: string
| SELECT xor(1,2)
< error
| 1064 42000 …near 'xor(1,2)' at line 1
> exec ^
: string
| SELECT xor(1,2,3)
< error
| 1064 42000 …near 'xor(1,2,3)' at line 1
> exec ^
: string
| SELECT xor(1,2,3,4)
< error
| 1064 42000 …near 'xor(1,2,3,4)' at line 1
> exec ^
: string
| SELECT xor(*)
< error
| 1064 42000 …near 'xor(*)' at line 1
> exec ^
: string
| SELECT xor('a')
< error
| 1064 42000 …near 'xor('a')' at line 1
> exec ^
: string
| SELECT xor(a USING utf8mb4)
< error
| 1064 42000 …near 'xor(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT xor(1 AS SIGNED)
< error
| 1064 42000 …near 'xor(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT xor(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'xor(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT xor(1 FROM 2)
< error
| 1064 42000 …near 'xor(1 FROM 2)' at line 1
> exec ^
: string
| SELECT xor(1 IN 2)
< error
| 1064 42000 …near 'xor(1 IN 2)' at line 1
> exec ^
: string
| SELECT xor(DISTINCT 1)
< error
| 1064 42000 …near 'xor(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT xor() OVER ()
< error
| 1064 42000 …near 'xor() OVER ()' at line 1
> exec ^
: string
| SELECT xor(1) OVER ()
< error
| 1064 42000 …near 'xor(1) OVER ()' at line 1
> exec ^
: string
| SELECT xor(1,2) OVER ()
< error
| 1064 42000 …near 'xor(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT xor(1,2,3) OVER ()
< error
| 1064 42000 …near 'xor(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT xor(1) FROM t
< error
| 1064 42000 …near 'xor(1) FROM t' at line 1
> exec ^
: string
| SELECT xor(SELECT 1)
< error
| 1064 42000 …near 'xor(SELECT 1)' at line 1
> exec ^
: string
| SELECT xor(YEAR FROM 1)
< error
| 1064 42000 …near 'xor(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT xor(DAY, 1, 2)
< error
| 1064 42000 …near 'xor(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT xor(a)
< error
| 1064 42000 …near 'xor(a)' at line 1
> exec ^
: string
| SELECT xor(a, b)
< error
| 1064 42000 …near 'xor(a, b)' at line 1
> exec ^
: string
| SELECT xor(a) OVER ()
< error
| 1064 42000 …near 'xor(a) OVER ()' at line 1
> exec ^
: string
| SELECT year()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT year(1)
< ok 0 0
> exec ^
: string
| SELECT year_month()
< error
| 1064 42000 …near 'year_month()' at line 1
> exec ^
: string
| SELECT year_month(1)
< error
| 1064 42000 …near 'year_month(1)' at line 1
> exec ^
: string
| SELECT year_month(1,2)
< error
| 1064 42000 …near 'year_month(1,2)' at line 1
> exec ^
: string
| SELECT year_month(1,2,3)
< error
| 1064 42000 …near 'year_month(1,2,3)' at line 1
> exec ^
: string
| SELECT year_month(1,2,3,4)
< error
| 1064 42000 …near 'year_month(1,2,3,4)' at line 1
> exec ^
: string
| SELECT year_month(*)
< error
| 1064 42000 …near 'year_month(*)' at line 1
> exec ^
: string
| SELECT year_month('a')
< error
| 1064 42000 …near 'year_month('a')' at line 1
> exec ^
: string
| SELECT year_month(a USING utf8mb4)
< error
| 1064 42000 …near 'year_month(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT year_month(1 AS SIGNED)
< error
| 1064 42000 …near 'year_month(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT year_month(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'year_month(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT year_month(1 FROM 2)
< error
| 1064 42000 …near 'year_month(1 FROM 2)' at line 1
> exec ^
: string
| SELECT year_month(1 IN 2)
< error
| 1064 42000 …near 'year_month(1 IN 2)' at line 1
> exec ^
: string
| SELECT year_month(DISTINCT 1)
< error
| 1064 42000 …near 'year_month(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT year_month() OVER ()
< error
| 1064 42000 …near 'year_month() OVER ()' at line 1
> exec ^
: string
| SELECT year_month(1) OVER ()
< error
| 1064 42000 …near 'year_month(1) OVER ()' at line 1
> exec ^
: string
| SELECT year_month(1,2) OVER ()
< error
| 1064 42000 …near 'year_month(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT year_month(1,2,3) OVER ()
< error
| 1064 42000 …near 'year_month(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT year_month(1) FROM t
< error
| 1064 42000 …near 'year_month(1) FROM t' at line 1
> exec ^
: string
| SELECT year_month(SELECT 1)
< error
| 1064 42000 …near 'year_month(SELECT 1)' at line 1
> exec ^
: string
| SELECT year_month(YEAR FROM 1)
< error
| 1064 42000 …near 'year_month(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT year_month(DAY, 1, 2)
< error
| 1064 42000 …near 'year_month(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT year_month(a)
< error
| 1064 42000 …near 'year_month(a)' at line 1
> exec ^
: string
| SELECT year_month(a, b)
< error
| 1064 42000 …near 'year_month(a, b)' at line 1
> exec ^
: string
| SELECT year_month(a) OVER ()
< error
| 1064 42000 …near 'year_month(a) OVER ()' at line 1
> exec ^
: string
| SELECT zerofill()
< error
| 1064 42000 …near 'zerofill()' at line 1
> exec ^
: string
| SELECT zerofill(1)
< error
| 1064 42000 …near 'zerofill(1)' at line 1
> exec ^
: string
| SELECT zerofill(1,2)
< error
| 1064 42000 …near 'zerofill(1,2)' at line 1
> exec ^
: string
| SELECT zerofill(1,2,3)
< error
| 1064 42000 …near 'zerofill(1,2,3)' at line 1
> exec ^
: string
| SELECT zerofill(1,2,3,4)
< error
| 1064 42000 …near 'zerofill(1,2,3,4)' at line 1
> exec ^
: string
| SELECT zerofill(*)
< error
| 1064 42000 …near 'zerofill(*)' at line 1
> exec ^
: string
| SELECT zerofill('a')
< error
| 1064 42000 …near 'zerofill('a')' at line 1
> exec ^
: string
| SELECT zerofill(a USING utf8mb4)
< error
| 1064 42000 …near 'zerofill(a USING utf8mb4)' at line 1
> exec ^
: string
| SELECT zerofill(1 AS SIGNED)
< error
| 1064 42000 …near 'zerofill(1 AS SIGNED)' at line 1
> exec ^
: string
| SELECT zerofill(1, INTERVAL 1 DAY)
< error
| 1064 42000 …near 'zerofill(1, INTERVAL 1 DAY)' at line 1
> exec ^
: string
| SELECT zerofill(1 FROM 2)
< error
| 1064 42000 …near 'zerofill(1 FROM 2)' at line 1
> exec ^
: string
| SELECT zerofill(1 IN 2)
< error
| 1064 42000 …near 'zerofill(1 IN 2)' at line 1
> exec ^
: string
| SELECT zerofill(DISTINCT 1)
< error
| 1064 42000 …near 'zerofill(DISTINCT 1)' at line 1
> exec ^
: string
| SELECT zerofill() OVER ()
< error
| 1064 42000 …near 'zerofill() OVER ()' at line 1
> exec ^
: string
| SELECT zerofill(1) OVER ()
< error
| 1064 42000 …near 'zerofill(1) OVER ()' at line 1
> exec ^
: string
| SELECT zerofill(1,2) OVER ()
< error
| 1064 42000 …near 'zerofill(1,2) OVER ()' at line 1
> exec ^
: string
| SELECT zerofill(1,2,3) OVER ()
< error
| 1064 42000 …near 'zerofill(1,2,3) OVER ()' at line 1
> exec ^
: string
| SELECT zerofill(1) FROM t
< error
| 1064 42000 …near 'zerofill(1) FROM t' at line 1
> exec ^
: string
| SELECT zerofill(SELECT 1)
< error
| 1064 42000 …near 'zerofill(SELECT 1)' at line 1
> exec ^
: string
| SELECT zerofill(YEAR FROM 1)
< error
| 1064 42000 …near 'zerofill(YEAR FROM 1)' at line 1
> exec ^
: string
| SELECT zerofill(DAY, 1, 2)
< error
| 1064 42000 …near 'zerofill(DAY, 1, 2)' at line 1
> exec ^
: string
| SELECT zerofill(a)
< error
| 1064 42000 …near 'zerofill(a)' at line 1
> exec ^
: string
| SELECT zerofill(a, b)
< error
| 1064 42000 …near 'zerofill(a, b)' at line 1
> exec ^
: string
| SELECT zerofill(a) OVER ()
< error
| 1064 42000 …near 'zerofill(a) OVER ()' at line 1
> exec ^
: string
| SELECT zone()
< error
| 1305 42000 FUNCTION test.zone does not exist
> exec ^
: string
| SELECT zone()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
