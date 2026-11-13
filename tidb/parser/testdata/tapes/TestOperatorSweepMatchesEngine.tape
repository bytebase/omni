-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestOperatorSweepMatchesEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_sweep
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 OR 1 OR NULL IN (+ NULL >= 0, 0 >> 3 < 2 <=> 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 OR 1 OR NULL IN (+ NULL >= 0, 0 >> 3 < 2 <=> 3)) <=> ((((NOT 3) OR 1) OR (NULL IN (((+ NULL) >= 0), (((0 >> 3) < 2) <=> 3)))))
< result
< column "(! 3 OR 1 OR NULL IN (+ NULL >= 0, 0 >> 3 < 2 <=> 3)) <=> ((((NOT 3) OR 1) OR (NULL IN (((+ NULL) >= 0), (((0 >> 3) < 2) <=> 3)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < 2) <=> ((1 < 2))
< result
< column "(1 < 2) <=> ((1 < 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 || 3 << - 2 NOT REGEXP NULL AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0 || 3 << - 2 NOT REGEXP NULL AND 2) <=> (((~ 0) OR (((3 << (- 2)) NOT REGEXP NULL) AND 2)))
< result
< column "(~ 0 || 3 << - 2 NOT REGEXP NULL AND 2) <=> (((~ 0) OR (((3 << (- 2)) NOT REGEXP NULL) AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (3 IN (0, 1), 0 <> 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 IN (3 IN (0, 1), 0 <> 1)) <=> ((0 IN ((3 IN (0, 1)), (0 <> 1))))
< result
< column "(0 IN (3 IN (0, 1), 0 <> 1)) <=> ((0 IN ((3 IN (0, 1)), (0 <> 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 || 2 NOT IN (NULL XOR 3, (0))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3 || 2 NOT IN (NULL XOR 3, (0))) <=> (((~ 3) OR (2 NOT IN ((NULL XOR 3), 0))))
< result
< column "(~ 3 || 2 NOT IN (NULL XOR 3, (0))) <=> (((~ 3) OR (2 NOT IN ((NULL XOR 3), 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3) <=> ((~ 3))
< result
< column "(~ 3) <=> ((~ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 NOT LIKE 1 DIV 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "DIV 0 AND 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 OR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 OR 1) <=> ((3 OR 1))
< result
< column "(3 OR 1) <=> ((3 OR 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0) <=> ((~ 0))
< result
< column "(~ 0) <=> ((~ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT LIKE 3 IN (NULL BETWEEN 2 AND 2, NULL IS NOT NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 22 near "IN (NULL BETWEEN 2 AND 2, NULL IS NOT NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 3 AND 1 BETWEEN 2 BETWEEN 0 AND 0 AND NULL NOT IN (3, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 46 near "BETWEEN 0 AND 0 AND NULL NOT IN (3, NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <=> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <=> 2) <=> ((2 <=> 2))
< result
< column "(2 <=> 2) <=> ((2 <=> 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (0, 1) * NOT 3 || 3 + 1 NOT LIKE ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "* NOT 3 || 3 + 1 NOT LIKE ~ 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 - 1) <=> ((3 - 1))
< result
< column "(3 - 1) <=> ((3 - 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && 2) <=> ((0 AND 2))
< result
< column "(0 && 2) <=> ((0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NULL | NULL / 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "| NULL / 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT IN (0, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT IN (0, 2)) <=> ((2 NOT IN (0, 2)))
< result
< column "(2 NOT IN (0, 2)) <=> ((2 NOT IN (0, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 OR + 2 <> 0 IN (2, 2) NOT BETWEEN 1 IN (3, 3) AND + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "NOT BETWEEN 1 IN (3, 3) AND + 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 BETWEEN 2 AND 3) <=> ((NOT (3 BETWEEN 2 AND 3)))
< result
< column "(NOT 3 BETWEEN 2 AND 3) <=> ((NOT (3 BETWEEN 2 AND 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2 NOT IN (0, 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((2 NOT IN (0, 2))) <=> ((2 NOT IN (0, 2)))
< result
< column "((2 NOT IN (0, 2))) <=> ((2 NOT IN (0, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0) <=> ((NOT 0))
< result
< column "(NOT 0) <=> ((NOT 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + ! 2 NOT BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ ! 2 NOT BETWEEN 0 AND NULL) <=> (((+ (NOT 2)) NOT BETWEEN 0 AND NULL))
< result
< column "(+ ! 2 NOT BETWEEN 0 AND NULL) <=> (((+ (NOT 2)) NOT BETWEEN 0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 > NULL & 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0 > NULL & 0) <=> (((+ 0) > (NULL & 0)))
< result
< column "(+ 0 > NULL & 0) <=> (((+ 0) > (NULL & 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IN (3, 1) << NULL IS NULL >= NULL IS NULL ^ NULL - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "<< NULL IS NULL >= NULL IS NULL ^ NULL - 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 IS NOT NULL << NULL NOT BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "<< NULL NOT BETWEEN NULL AND 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! ! NULL + 3 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! ! NULL + 3 AND 3) <=> ((((NOT (NOT NULL)) + 3) AND 3))
< result
< column "(! ! NULL + 3 AND 3) <=> ((((NOT (NOT NULL)) + 3) AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 2 AND 1 IN (3 > 1, NULL <=> 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 2 AND 1 IN (3 > 1, NULL <=> 2)) <=> ((NULL NOT BETWEEN 2 AND (1 IN ((3 > 1), (NULL <=> 2)))))
< result
< column "(NULL NOT BETWEEN 2 AND 1 IN (3 > 1, NULL <=> 2)) <=> ((NULL NOT BETWEEN 2 AND (1 IN ((3 > 1), (NULL <=> 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < 3 >> 2 <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < 3 >> 2 <=> 1) <=> (((1 < (3 >> 2)) <=> 1))
< result
< column "(1 < 3 >> 2 <=> 1) <=> (((1 < (3 >> 2)) <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 <=> NULL AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 <=> NULL AND NULL) <=> ((((+ 3) <=> NULL) AND NULL))
< result
< column "(+ 3 <=> NULL AND NULL) <=> ((((+ 3) <=> NULL) AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 / 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 / 1) <=> ((3 / 1))
< result
< column "(3 / 1) <=> ((3 / 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 = 1 > 1 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 = 1 > 1 IS NULL) <=> ((((1 = 1) > 1) IS NULL))
< result
< column "(1 = 1 > 1 IS NULL) <=> ((((1 = 1) > 1) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 0 AND 3 = 1 << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 0 AND 3 = 1 << 3) <=> (((2 NOT BETWEEN 0 AND 3) = (1 << 3)))
< result
< column "(2 NOT BETWEEN 0 AND 3 = 1 << 3) <=> (((2 NOT BETWEEN 0 AND 3) = (1 << 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 MOD 2 AND ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 MOD 2 AND ~ 1) <=> (((2 % 2) AND (~ 1)))
< result
< column "(2 MOD 2 AND ~ 1) <=> (((2 % 2) AND (~ 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NULL <> 2 XOR NULL < 2 NOT IN (NULL, 1) - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 53 near "- 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 * 1) <=> ((0 * 1))
< result
< column "(0 * 1) <=> ((0 * 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 * NULL) <=> ((0 * NULL))
< result
< column "(0 * NULL) <=> ((0 * NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL OR NULL & 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL OR NULL & 0) <=> (((- NULL) OR (NULL & 0)))
< result
< column "(- NULL OR NULL & 0) <=> (((- NULL) OR (NULL & 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN 0 AND NULL) <=> ((1 BETWEEN 0 AND NULL))
< result
< column "(1 BETWEEN 0 AND NULL) <=> ((1 BETWEEN 0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 2 LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 2 LIKE 2) <=> ((@v := (2 LIKE 2)))
< result
< column "(@v := 2 LIKE 2) <=> ((@v := (2 LIKE 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS UNKNOWN) <=> ((2 IS UNKNOWN))
< result
< column "(2 IS UNKNOWN) <=> ((2 IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL * 3 * 3 ^ 0 NOT BETWEEN @v := 3 OR NULL AND 1 < 3 NOT BETWEEN NULL BETWEEN 2 AND 1 AND 1 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 41 near ":= 3 OR NULL AND 1 < 3 NOT BETWEEN NULL BETWEEN 2 AND 1 AND 1 + 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (1, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (1, 0)) <=> ((0 NOT IN (1, 0)))
< result
< column "(0 NOT IN (1, 0)) <=> ((0 NOT IN (1, 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1)) <=> (1)
< result
< column "((1)) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 | NULL) <=> ((2 | NULL))
< result
< column "(2 | NULL) <=> ((2 | NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL | 1 <=> 3 IN (1, NULL) OR 1 OR 1 > 1 | 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL | 1 <=> 3 IN (1, NULL) OR 1 OR 1 > 1 | 0) <=> (((((NULL | 1) <=> (3 IN (1, NULL))) OR 1) OR (1 > (1 | 0))))
< result
< column "(NULL | 1 <=> 3 IN (1, NULL) OR 1 OR 1 > 1 | 0) <=> (((((NULL | 1) <=> (3 IN (1, NULL))) OR 1) OR (1 > (1 | 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL DIV NULL || NOT 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL DIV NULL || NOT 3) <=> (((NULL DIV NULL) OR (NOT 3)))
< result
< column "(NULL DIV NULL || NOT 3) <=> (((NULL DIV NULL) OR (NOT 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 >> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2 >> 3) <=> (((NOT 2) >> 3))
< result
< column "(! 2 >> 3) <=> (((NOT 2) >> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 XOR 3 + ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 XOR 3 + ~ NULL) <=> ((0 XOR (3 + (~ NULL))))
< result
< column "(0 XOR 3 + ~ NULL) <=> ((0 XOR (3 + (~ NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 > 2 AND 2 * 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 > 2 AND 2 * 2) <=> (((2 > 2) AND (2 * 2)))
< result
< column "(2 > 2 AND 2 * 2) <=> (((2 > 2) AND (2 * 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <> NULL) <=> ((NULL <> NULL))
< result
< column "(NULL <> NULL) <=> ((NULL <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 * NOT NULL <= 1 AND 2 > - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 16 near "NOT NULL <= 1 AND 2 > - 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 + 2) <=> ((0 + 2))
< result
< column "(0 + 2) <=> ((0 + 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> 2) <=> ((0 >> 2))
< result
< column "(0 >> 2) <=> ((0 >> 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL MOD 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL MOD 0) <=> (((+ NULL) % 0))
< result
< column "(+ NULL MOD 0) <=> (((+ NULL) % 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 2 AND 3) <=> ((NULL NOT BETWEEN 2 AND 3))
< result
< column "(NULL NOT BETWEEN 2 AND 3) <=> ((NULL NOT BETWEEN 2 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0) - - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0) - - 1) <=> ((0 - (- 1)))
< result
< column "((0) - - 1) <=> ((0 - (- 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 / 2 >= 0 MOD 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 / 2 >= 0 MOD 0) <=> (((3 / 2) >= (0 % 0)))
< result
< column "(3 / 2 >= 0 MOD 0) <=> (((3 / 2) >= (0 % 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 0 <> 3 NOT LIKE NULL ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 31 near "^ 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL - 3 LIKE 0 << 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "<< 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NULL) <=> ((NULL IS NULL))
< result
< column "(NULL IS NULL) <=> ((NULL IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 <= NULL | - 3 << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2 <= NULL | - 3 << 3) <=> (((NOT 2) <= (NULL | ((- 3) << 3))))
< result
< column "(! 2 <= NULL | - 3 << 3) <=> (((NOT 2) <= (NULL | ((- 3) << 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 DIV 2 + 2 > 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 DIV 2 + 2 > 0) <=> ((((3 DIV 2) + 2) > 0))
< result
< column "(3 DIV 2 + 2 > 0) <=> ((((3 DIV 2) + 2) > 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = NULL <=> NOT NULL != ! 1 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "NOT NULL != ! 1 IS NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 ^ 3) <=> ((2 ^ 3))
< result
< column "(2 ^ 3) <=> ((2 ^ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 <= NULL BETWEEN ~ 3 AND 1 BETWEEN 0 AND 0 >= 3 >= 0 <= 1 < 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 <= NULL BETWEEN ~ 3 AND 1 BETWEEN 0 AND 0 >= 3 >= 0 <= 1 < 3) <=> ((((((0 <= (NULL BETWEEN (~ 3) AND (1 BETWEEN 0 AND 0))) >= 3) >= 0) <= 1) < 3))
< result
< column "(0 <= NULL BETWEEN ~ 3 AND 1 BETWEEN 0 AND 0 >= 3 >= 0 <= 1 < 3) <=> ((((((0 <= (NULL BETWEEN (~ 3) AND (1 BETWEEN 0 AND 0))) >= 3) >= 0) <= 1) < 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 NOT BETWEEN 3 OR 0 AND 1 BETWEEN 2 AND 3 != + (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "OR 0 AND 1 BETWEEN 2 AND 3 != + (2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 0 % 3 < 2 IN (~ 1 | 2, 1 BETWEEN 1 AND NULL IN (3 > NULL, 3 MOD NULL))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 | 0 % 3 < 2 IN (~ 1 | 2, 1 BETWEEN 1 AND NULL IN (3 > NULL, 3 MOD NULL))) <=> (((2 | (0 % 3)) < (2 IN (((~ 1) | 2), (1 BETWEEN 1 AND (NULL IN ((3 > NULL), (3 % NULL))))))))
< result
< column "(2 | 0 % 3 < 2 IN (~ 1 | 2, 1 BETWEEN 1 AND NULL IN (3 > NULL, 3 MOD NULL))) <=> (((2 | (0 % 3)) < (2 IN (((~ 1) | 2), (1 BETWEEN 1 AND (NULL IN ((3 > NULL), (3 % NULL))))))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 NOT BETWEEN 2 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 NOT BETWEEN 2 AND 0) <=> ((NOT (2 NOT BETWEEN 2 AND 0)))
< result
< column "(NOT 2 NOT BETWEEN 2 AND 0) <=> ((NOT (2 NOT BETWEEN 2 AND 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 1 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 1 AND 0) <=> ((NULL NOT BETWEEN 1 AND 0))
< result
< column "(NULL NOT BETWEEN 1 AND 0) <=> ((NULL NOT BETWEEN 1 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL OR 0 IN (0 IN (3, 3), 3 || 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL OR 0 IN (0 IN (3, 3), 3 || 0)) <=> (((NOT NULL) OR (0 IN ((0 IN (3, 3)), (3 OR 0)))))
< result
< column "(NOT NULL OR 0 IN (0 IN (3, 3), 3 || 0)) <=> (((NOT NULL) OR (0 IN ((0 IN (3, 3)), (3 OR 0)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ + 3) <=> ((~ (+ 3)))
< result
< column "(~ + 3) <=> ((~ (+ 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT IN (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT IN (1, 2)) <=> ((NULL NOT IN (1, 2)))
< result
< column "(NULL NOT IN (1, 2)) <=> ((NULL NOT IN (1, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 DIV 3 || NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 DIV 3 || NULL IS NULL) <=> (((3 DIV 3) OR (NULL IS NULL)))
< result
< column "(3 DIV 3 || NULL IS NULL) <=> (((3 DIV 3) OR (NULL IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | 3) <=> ((0 | 3))
< result
< column "(0 | 3) <=> ((0 | 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2) <=> ((NOT 2))
< result
< column "(NOT 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> 3) <=> ((1 <=> 3))
< result
< column "(1 <=> 3) <=> ((1 <=> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS UNKNOWN || 1 DIV 1 IN (2 ^ 2, 1 >= NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS UNKNOWN || 1 DIV 1 IN (2 ^ 2, 1 >= NULL)) <=> (((3 IS UNKNOWN) OR ((1 DIV 1) IN ((2 ^ 2), (1 >= NULL)))))
< result
< column "(3 IS UNKNOWN || 1 DIV 1 IN (2 ^ 2, 1 >= NULL)) <=> (((3 IS UNKNOWN) OR ((1 DIV 1) IN ((2 ^ 2), (1 >= NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < NULL) <=> ((1 < NULL))
< result
< column "(1 < NULL) <=> ((1 < NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2) <=> ((NOT 2))
< result
< column "(! 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL MOD 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL MOD 2) <=> ((NULL % 2))
< result
< column "(NULL MOD 2) <=> ((NULL % 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN NULL AND 1) <=> ((3 BETWEEN NULL AND 1))
< result
< column "(3 BETWEEN NULL AND 1) <=> ((3 BETWEEN NULL AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 = NULL >> + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 = NULL >> + NULL) <=> ((0 = (NULL >> (+ NULL))))
< result
< column "(0 = NULL >> + NULL) <=> ((0 = (NULL >> (+ NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = 2 != NULL IS NOT FALSE IN (2 - NULL && NULL - 2, NOT 3 BETWEEN 0 AND 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 36 near "IN (2 - NULL && NULL - 2, NOT 3 BETWEEN 0 AND 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 > 2 NOT BETWEEN 1 AND 3 < 1 IN (NULL, 1) MOD 1 && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 55 near "MOD 1 && NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL) <=> ((NOT NULL))
< result
< column "(! NULL) <=> ((NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL)) <=> (NULL)
< result
< column "((NULL)) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 AND 0) <=> ((3 AND 0))
< result
< column "(3 AND 0) <=> ((3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <> 2 AND 0 NOT BETWEEN 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <> 2 AND 0 NOT BETWEEN 0 AND 3) <=> (((2 <> 2) AND (0 NOT BETWEEN 0 AND 3)))
< result
< column "(2 <> 2 AND 0 NOT BETWEEN 0 AND 3) <=> (((2 <> 2) AND (0 NOT BETWEEN 0 AND 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV 0) <=> ((1 DIV 0))
< result
< column "(1 DIV 0) <=> ((1 DIV 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN NULL AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN NULL AND NULL) <=> ((NULL BETWEEN NULL AND NULL))
< result
< column "(NULL BETWEEN NULL AND NULL) <=> ((NULL BETWEEN NULL AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 * 1 LIKE 1 < 2 != 0 / 2 <= 3 % 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 * 1 LIKE 1 < 2 != 0 / 2 <= 3 % 1) <=> ((((((1 * 1) LIKE 1) < 2) <> (0 / 2)) <= (3 % 1)))
< result
< column "(1 * 1 LIKE 1 < 2 != 0 / 2 <= 3 % 1) <=> ((((((1 * 1) LIKE 1) < 2) <> (0 / 2)) <= (3 % 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2) <=> ((~ 2))
< result
< column "(~ 2) <=> ((~ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <= NULL IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <= NULL IS UNKNOWN) <=> (((NULL <= NULL) IS UNKNOWN))
< result
< column "(NULL <= NULL IS UNKNOWN) <=> (((NULL <= NULL) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0) % 0 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0) % 0 IS NOT FALSE) <=> (((0 % 0) IS NOT FALSE))
< result
< column "((0) % 0 IS NOT FALSE) <=> (((0 % 0) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL + 0 OR ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL + 0 OR ! 2) <=> (((NULL + 0) OR (NOT 2)))
< result
< column "(NULL + 0 OR ! 2) <=> (((NULL + 0) OR (NOT 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 3 + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 3 + 0) <=> ((@v := (3 + 0)))
< result
< column "(@v := 3 + 0) <=> ((@v := (3 + 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + NULL) <=> ((3 + NULL))
< result
< column "(3 + NULL) <=> ((3 + NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 AND 0 AND 1 AND 2 << 2 ^ 3 <> 3 NOT IN (1, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 AND 0 AND 1 AND 2 << 2 ^ 3 <> 3 NOT IN (1, NULL)) <=> ((((2 AND 0) AND 1) AND ((2 << (2 ^ 3)) <> (3 NOT IN (1, NULL)))))
< result
< column "(2 AND 0 AND 1 AND 2 << 2 ^ 3 <> 3 NOT IN (1, NULL)) <=> ((((2 AND 0) AND 1) AND ((2 << (2 ^ 3)) <> (3 NOT IN (1, NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >> 1 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >> 1 IS TRUE) <=> (((1 >> 1) IS TRUE))
< result
< column "(1 >> 1 IS TRUE) <=> (((1 >> 1) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (0, 0) | 0 OR 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 20 near "| 0 OR 3 IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT IN (3, 0) % NULL IS NOT NULL > 0 REGEXP 3 ^ 2 IN (1, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "% NULL IS NOT NULL > 0 REGEXP 3 ^ 2 IN (1, 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 0 AND 2 <=> 0 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 0 AND 2 <=> 0 IS NOT NULL) <=> ((((2 NOT BETWEEN 0 AND 2) <=> 0) IS NOT NULL))
< result
< column "(2 NOT BETWEEN 0 AND 2 <=> 0 IS NOT NULL) <=> ((((2 NOT BETWEEN 0 AND 2) <=> 0) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN 0 AND 2) <=> ((NULL BETWEEN 0 AND 2))
< result
< column "(NULL BETWEEN 0 AND 2) <=> ((NULL BETWEEN 0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 <=> 2 != 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 <=> 2 != 0) <=> ((((~ 2) <=> 2) <> 0))
< result
< column "(~ 2 <=> 2 != 0) <=> ((((~ 2) <=> 2) <> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >> 1 IS NOT FALSE) <=> (((2 >> 1) IS NOT FALSE))
< result
< column "(2 >> 1 IS NOT FALSE) <=> (((2 >> 1) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 MOD 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 MOD 2) <=> ((0 % 2))
< result
< column "(0 MOD 2) <=> ((0 % 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 LIKE NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 LIKE NULL) <=> ((1 LIKE NULL))
< result
< column "(1 LIKE NULL) <=> ((1 LIKE NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> 1) <=> ((1 <=> 1))
< result
< column "(1 <=> 1) <=> ((1 <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL > 1 DIV 0 || 1 = 1 <=> 2 BETWEEN 1 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL > 1 DIV 0 || 1 = 1 <=> 2 BETWEEN 1 AND NULL) <=> ((((NOT NULL) > (1 DIV 0)) OR ((1 = 1) <=> (2 BETWEEN 1 AND NULL))))
< result
< column "(! NULL > 1 DIV 0 || 1 = 1 <=> 2 BETWEEN 1 AND NULL) <=> ((((NOT NULL) > (1 DIV 0)) OR ((1 = 1) <=> (2 BETWEEN 1 AND NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IN (3, NULL) <=> 0 + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IN (3, NULL) <=> 0 + NULL) <=> (((3 IN (3, NULL)) <=> (0 + NULL)))
< result
< column "(3 IN (3, NULL) <=> 0 + NULL) <=> (((3 IN (3, NULL)) <=> (0 + NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT ! NULL <=> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT ! NULL <=> 3) <=> ((NOT ((NOT NULL) <=> 3)))
< result
< column "(NOT ! NULL <=> 3) <=> ((NOT ((NOT NULL) <=> 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 NOT IN (3 NOT IN (2, NULL), (NULL)) NOT IN (NULL IN (1, NULL) >> (NULL), ~ NULL + 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 50 near "NOT IN (NULL IN (1, NULL) >> (NULL), ~ NULL + 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 0 AND 0 BETWEEN NULL || 2 AND 2 NOT IN (1, 0) + NULL IS NOT NULL > 2 NOT BETWEEN NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "|| 2 AND 2 NOT IN (1, 0) + NULL IS NOT NULL > 2 NOT BETWEEN NULL AND 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 0 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN 0 AND 0) <=> ((NULL BETWEEN 0 AND 0))
< result
< column "(NULL BETWEEN 0 AND 0) <=> ((NULL BETWEEN 0 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (NULL, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (NULL, 2)) <=> ((2 IN (NULL, 2)))
< result
< column "(2 IN (NULL, 2)) <=> ((2 IN (NULL, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 <=> 0 BETWEEN 1 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 <=> 0 BETWEEN 1 AND 2) <=> (((+ 3) <=> (0 BETWEEN 1 AND 2)))
< result
< column "(+ 3 <=> 0 BETWEEN 1 AND 2) <=> (((+ 3) <=> (0 BETWEEN 1 AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL % 1 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL % 1 IS NULL) <=> (((NULL % 1) IS NULL))
< result
< column "(NULL % 1 IS NULL) <=> (((NULL % 1) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL ^ 1) & 2 < 2 IN (2 MOD 1, 2 * 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL ^ 1) & 2 < 2 IN (2 MOD 1, 2 * 0)) <=> ((((NULL ^ 1) & 2) < (2 IN ((2 % 1), (2 * 0)))))
< result
< column "((NULL ^ 1) & 2 < 2 IN (2 MOD 1, 2 * 0)) <=> ((((NULL ^ 1) & 2) < (2 IN ((2 % 1), (2 * 0)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 0) <=> ((@v := 0))
< result
< column "(@v := 0) <=> ((@v := 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 < 3 < ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 < 3 < ~ 3) <=> ((((NOT 0) < 3) < (~ 3)))
< result
< column "(! 0 < 3 < ~ 3) <=> ((((NOT 0) < 3) < (~ 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >> 1) <=> ((1 >> 1))
< result
< column "(1 >> 1) <=> ((1 >> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1) <=> ((+ 1))
< result
< column "(+ 1) <=> ((+ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 BETWEEN NULL AND 2 NOT LIKE 2 >= NULL >> NULL >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2 BETWEEN NULL AND 2 NOT LIKE 2 >= NULL >> NULL >= 2) <=> (((((NOT 2) BETWEEN NULL AND (2 NOT LIKE 2)) >= (NULL >> NULL)) >= 2))
< result
< column "(! 2 BETWEEN NULL AND 2 NOT LIKE 2 >= NULL >> NULL >= 2) <=> (((((NOT 2) BETWEEN NULL AND (2 NOT LIKE 2)) >= (NULL >> NULL)) >= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 <> 2 | (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1 <> 2 | (1)) <=> (((+ 1) <> (2 | 1)))
< result
< column "(+ 1 <> 2 | (1)) <=> (((+ 1) <> (2 | 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <> 1) <=> ((1 <> 1))
< result
< column "(1 <> 1) <=> ((1 <> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS NOT FALSE & ~ NULL MOD 2 ^ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "& ~ NULL MOD 2 ^ NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL)) <=> (NULL)
< result
< column "((NULL)) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 - 2 << NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 - 2 << NULL IS NULL) <=> ((((1 - 2) << NULL) IS NULL))
< result
< column "(1 - 2 << NULL IS NULL) <=> ((((1 - 2) << NULL) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS UNKNOWN) <=> ((1 IS UNKNOWN))
< result
< column "(1 IS UNKNOWN) <=> ((1 IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT NULL) <=> ((3 IS NOT NULL))
< result
< column "(3 IS NOT NULL) <=> ((3 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 IS UNKNOWN) <=> (((NOT 0) IS UNKNOWN))
< result
< column "(! 0 IS UNKNOWN) <=> (((NOT 0) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 & 2 IS UNKNOWN && 1 BETWEEN 3 AND NULL DIV 1 NOT BETWEEN 3 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 & 2 IS UNKNOWN && 1 BETWEEN 3 AND NULL DIV 1 NOT BETWEEN 3 AND NULL) <=> ((((0 & 2) IS UNKNOWN) AND (1 BETWEEN 3 AND ((NULL DIV 1) NOT BETWEEN 3 AND NULL))))
< result
< column "(0 & 2 IS UNKNOWN && 1 BETWEEN 3 AND NULL DIV 1 NOT BETWEEN 3 AND NULL) <=> ((((0 & 2) IS UNKNOWN) AND (1 BETWEEN 3 AND ((NULL DIV 1) NOT BETWEEN 3 AND NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL <=> NULL REGEXP 0 * 3 < 1 BETWEEN ~ 3 AND 3 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "* 3 < 1 BETWEEN ~ 3 AND 3 IS NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND 3 IN (NOT 0, 3 <= 0) IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 AND 3 IN (NOT 0, 3 <= 0) IS TRUE) <=> ((3 AND ((3 IN ((NOT 0), (3 <= 0))) IS TRUE)))
< result
< column "(3 AND 3 IN (NOT 0, 3 <= 0) IS TRUE) <=> ((3 AND ((3 IN ((NOT 0), (3 <= 0))) IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL & 0 <=> NULL * 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL & 0 <=> NULL * 0) <=> (((NULL & 0) <=> (NULL * 0)))
< result
< column "(NULL & 0 <=> NULL * 0) <=> (((NULL & 0) <=> (NULL * 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0 && 3) > NOT 3 NOT BETWEEN NULL IS UNKNOWN AND 2 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "NOT 3 NOT BETWEEN NULL IS UNKNOWN AND 2 OR 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 3 || 0 IN (0, 3) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= 3 || 0 IN (0, 3) IS NOT FALSE) <=> (((2 >= 3) OR ((0 IN (0, 3)) IS NOT FALSE)))
< result
< column "(2 >= 3 || 0 IN (0, 3) IS NOT FALSE) <=> (((2 >= 3) OR ((0 IN (0, 3)) IS NOT FALSE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 1 AND 2 * NULL LIKE 3 NOT IN (1 BETWEEN NULL AND 0 MOD 0 IS UNKNOWN, 0 AND NULL / 1 / 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 45 near "NOT IN (1 BETWEEN NULL AND 0 MOD 0 IS UNKNOWN, 0 AND NULL / 1 / 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 || 0 / 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 || 0 / 1) <=> ((3 OR (0 / 1)))
< result
< column "(3 || 0 / 1) <=> ((3 OR (0 / 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 2) <=> ((1 XOR 2))
< result
< column "(1 XOR 2) <=> ((1 XOR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 != 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 != 0) <=> (((- 1) <> 0))
< result
< column "(- 1 != 0) <=> (((- 1) <> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 || 2 DIV 3 NOT BETWEEN 0 + 3 < 1 IS UNKNOWN AND 1 != 1 NOT IN (2 MOD 1, 0 BETWEEN 0 AND 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 41 near "< 1 IS UNKNOWN AND 1 != 1 NOT IN (2 MOD 1, 0 BETWEEN 0 AND 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 = NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 = NULL) <=> ((2 = NULL))
< result
< column "(2 = NULL) <=> ((2 = NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > NULL > 2 IN (2, 1) IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL > NULL > 2 IN (2, 1) IS TRUE) <=> ((((NULL > NULL) > (2 IN (2, 1))) IS TRUE))
< result
< column "(NULL > NULL > 2 IN (2, 1) IS TRUE) <=> ((((NULL > NULL) > (2 IN (2, 1))) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1 > 1 < NULL NOT IN (2, 1) XOR - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1 > 1 < NULL NOT IN (2, 1) XOR - 2) <=> (((((~ 1) > 1) < (NULL NOT IN (2, 1))) XOR (- 2)))
< result
< column "(~ 1 > 1 < NULL NOT IN (2, 1) XOR - 2) <=> (((((~ 1) > 1) < (NULL NOT IN (2, 1))) XOR (- 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 0 AND 3) <=> ((2 BETWEEN 0 AND 3))
< result
< column "(2 BETWEEN 0 AND 3) <=> ((2 BETWEEN 0 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 ^ 3 / 1 >> 3 BETWEEN NULL << 1 <> (0) AND 3 = NULL <=> 1 != 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 42 near "<> (0) AND 3 = NULL <=> 1 != 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 0 AND 1 <=> 3 IN (0, 0) BETWEEN NULL * 2 AND 2 >> 1 AND 2 >> 2 NOT LIKE 0 XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 52 near "BETWEEN NULL * 2 AND 2 >> 1 AND 2 >> 2 NOT LIKE 0 XOR 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 1 AND 2 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 1 AND 2 IS NULL) <=> (((3 BETWEEN 1 AND 2) IS NULL))
< result
< column "(3 BETWEEN 1 AND 2 IS NULL) <=> (((3 BETWEEN 1 AND 2) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 << 3 <= 1 IN (NULL, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 << 3 <= 1 IN (NULL, 0)) <=> (((0 << 3) <= (1 IN (NULL, 0))))
< result
< column "(0 << 3 <= 1 IN (NULL, 0)) <=> (((0 << 3) <= (1 IN (NULL, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL <=> 0 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL <=> 0 <> NULL) <=> ((NOT ((NULL <=> 0) <> NULL)))
< result
< column "(NOT NULL <=> 0 <> NULL) <=> ((NOT ((NULL <=> 0) <> NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL / 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL / 0) <=> ((NULL / 0))
< result
< column "(NULL / 0) <=> ((NULL / 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 < 2 <> 0 >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 < 2 <> 0 >> 1) <=> (((0 < 2) <> (0 >> 1)))
< result
< column "(0 < 2 <> 0 >> 1) <=> (((0 < 2) <> (0 >> 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (3 BETWEEN 2 AND 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((3 BETWEEN 2 AND 2)) <=> ((3 BETWEEN 2 AND 2))
< result
< column "((3 BETWEEN 2 AND 2)) <=> ((3 BETWEEN 2 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL IS TRUE NOT REGEXP 3 % 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "NOT REGEXP 3 % 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND 3 || NULL < NULL | 2 >> NULL IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND 3 || NULL < NULL | 2 >> NULL IS TRUE) <=> (((NULL AND 3) OR ((NULL < (NULL | (2 >> NULL))) IS TRUE)))
< result
< column "(NULL AND 3 || NULL < NULL | 2 >> NULL IS TRUE) <=> (((NULL AND 3) OR ((NULL < (NULL | (2 >> NULL))) IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 = 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 = 0) <=> ((1 = 0))
< result
< column "(1 = 0) <=> ((1 = 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 <=> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 <=> NULL) <=> ((0 <=> NULL))
< result
< column "(0 <=> NULL) <=> ((0 <=> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= 3 DIV 0 & NULL IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= 3 DIV 0 & NULL IS UNKNOWN) <=> (((3 <= ((3 DIV 0) & NULL)) IS UNKNOWN))
< result
< column "(3 <= 3 DIV 0 & NULL IS UNKNOWN) <=> (((3 <= ((3 DIV 0) & NULL)) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 > NULL | 1 AND 1 IN (NULL XOR 1 OR 1 BETWEEN 3 AND NULL, ! 0 + 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 > NULL | 1 AND 1 IN (NULL XOR 1 OR 1 BETWEEN 3 AND NULL, ! 0 + 0)) <=> (((1 > (NULL | 1)) AND (1 IN (((NULL XOR 1) OR (1 BETWEEN 3 AND NULL)), ((NOT 0) + 0)))))
< result
< column "(1 > NULL | 1 AND 1 IN (NULL XOR 1 OR 1 BETWEEN 3 AND NULL, ! 0 + 0)) <=> (((1 > (NULL | 1)) AND (1 IN (((NULL XOR 1) OR (1 BETWEEN 3 AND NULL)), ((NOT 0) + 0)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS UNKNOWN / 0 OR 1 LIKE NOT NULL IN (- NULL, 1 NOT REGEXP 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "/ 0 OR 1 LIKE NOT NULL IN (- NULL, 1 NOT REGEXP 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL MOD 1 * ! NULL NOT BETWEEN 3 ^ 1 * 1 XOR 2 AND NOT NULL BETWEEN 1 NOT LIKE 0 AND 0 < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 52 near "XOR 2 AND NOT NULL BETWEEN 1 NOT LIKE 0 AND 0 < 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 3 / ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 3 / ! 1) <=> ((NULL = (3 / (NOT 1))))
< result
< column "(NULL = 3 / ! 1) <=> ((NULL = (3 / (NOT 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >> 2) <=> ((2 >> 2))
< result
< column "(2 >> 2) <=> ((2 >> 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 0) <=> ((- 0))
< result
< column "(- 0) <=> ((- 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 + 0 XOR NULL IN (2 IS UNKNOWN, NULL MOD 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 + 0 XOR NULL IN (2 IS UNKNOWN, NULL MOD 1)) <=> (((1 + 0) XOR (NULL IN ((2 IS UNKNOWN), (NULL % 1)))))
< result
< column "(1 + 0 XOR NULL IN (2 IS UNKNOWN, NULL MOD 1)) <=> (((1 + 0) XOR (NULL IN ((2 IS UNKNOWN), (NULL % 1)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 0 AND 1 OR NULL <> 1 DIV NULL AND NULL BETWEEN 2 XOR NULL AND 3 < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 69 near "XOR NULL AND 3 < 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN NULL AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN NULL AND 2) <=> ((2 BETWEEN NULL AND 2))
< result
< column "(2 BETWEEN NULL AND 2) <=> ((2 BETWEEN NULL AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 ^ NULL / 3 || 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 ^ NULL / 3 || 3) <=> ((((1 ^ NULL) / 3) OR 3))
< result
< column "(1 ^ NULL / 3 || 3) <=> ((((1 ^ NULL) / 3) OR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 & 0 LIKE 1 - 3 IN (NULL <= NULL ^ 2 MOD 0, 0 BETWEEN 3 AND 0 XOR 3 | 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "- 3 IN (NULL <= NULL ^ 2 MOD 0, 0 BETWEEN 3 AND 0 XOR 3 | 1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL >> (NULL) / 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL >> (NULL) / 0) <=> ((NULL >> (NULL / 0)))
< result
< column "(NULL >> (NULL) / 0) <=> ((NULL >> (NULL / 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (3, 3) NOT IN (NOT 1, @v := 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "NOT IN (NOT 1, @v := 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 0) <=> ((- 0))
< result
< column "(- 0) <=> ((- 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 1) <=> ((NULL = 1))
< result
< column "(NULL = 1) <=> ((NULL = 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT IN (NULL, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT IN (NULL, NULL)) <=> ((2 NOT IN (NULL, NULL)))
< result
< column "(2 NOT IN (NULL, NULL)) <=> ((2 NOT IN (NULL, NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2) <=> ((~ 2))
< result
< column "(~ 2) <=> ((~ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = NULL NOT LIKE 2 NOT IN (0, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "NOT IN (0, 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL OR 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL OR 3) <=> ((NULL OR 3))
< result
< column "(NULL OR 3) <=> ((NULL OR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 2) <=> ((@v := 2))
< result
< column "(@v := 2) <=> ((@v := 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <> 0 IS NOT NULL OR 1 <> 0 | NULL ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <> 0 IS NOT NULL OR 1 <> 0 | NULL ^ 2) <=> ((((NULL <> 0) IS NOT NULL) OR (1 <> (0 | (NULL ^ 2)))))
< result
< column "(NULL <> 0 IS NOT NULL OR 1 <> 0 | NULL ^ 2) <=> ((((NULL <> 0) IS NOT NULL) OR (1 <> (0 | (NULL ^ 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <> NULL) <=> ((3 <> NULL))
< result
< column "(3 <> NULL) <=> ((3 <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <= 2 / 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <= 2 / 0) <=> ((NULL <= (2 / 0)))
< result
< column "(NULL <= 2 / 0) <=> ((NULL <= (2 / 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 REGEXP 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 REGEXP 2) <=> ((0 REGEXP 2))
< result
< column "(0 REGEXP 2) <=> ((0 REGEXP 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL - 0 <> NULL * 1 <= - 1 = 0 <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL - 0 <> NULL * 1 <= - 1 = 0 <= NULL) <=> ((((((NULL - 0) <> (NULL * 1)) <= (- 1)) = 0) <= NULL))
< result
< column "(NULL - 0 <> NULL * 1 <= - 1 = 0 <= NULL) <=> ((((((NULL - 0) <> (NULL * 1)) <= (- 1)) = 0) <= NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS UNKNOWN ^ @v := 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "^ @v := 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 LIKE 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 LIKE 1) <=> ((1 LIKE 1))
< result
< column "(1 LIKE 1) <=> ((1 LIKE 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 % 0) <=> ((3 % 0))
< result
< column "(3 % 0) <=> ((3 % 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 0 AND 1 ^ NULL & 3 NOT IN (3 AND 3 % NULL BETWEEN 2 AND 2, 3 != 1 OR 3 = 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN 0 AND 1 ^ NULL & 3 NOT IN (3 AND 3 % NULL BETWEEN 2 AND 2, 3 != 1 OR 3 = 2)) <=> ((3 NOT BETWEEN 0 AND (((1 ^ NULL) & 3) NOT IN ((3 AND ((3 % NULL) BETWEEN 2 AND 2)), ((3 <> 1) OR (3 = 2))))))
< result
< column "(3 NOT BETWEEN 0 AND 1 ^ NULL & 3 NOT IN (3 AND 3 % NULL BETWEEN 2 AND 2, 3 != 1 OR 3 = 2)) <=> ((3 NOT BETWEEN 0 AND (((1 ^ NULL) & 3) NOT IN ((3 AND ((3 % NULL) BETWEEN 2 AND 2)), ((3 <> 1) OR (3 = 2))))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE & NULL && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "& NULL && NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 OR 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 OR 3) <=> ((3 OR 3))
< result
< column "(3 OR 3) <=> ((3 OR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 && 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 && 0) <=> ((3 AND 0))
< result
< column "(3 && 0) <=> ((3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 - 0 && 3 && 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 - 0 && 3 && 3 IS TRUE) <=> ((((2 - 0) AND 3) AND (3 IS TRUE)))
< result
< column "(2 - 0 && 3 && 3 IS TRUE) <=> ((((2 - 0) AND 3) AND (3 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1) <=> ((~ 1))
< result
< column "(~ 1) <=> ((~ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 << 1 DIV 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 << 1 DIV 1) <=> ((NOT (1 << (1 DIV 1))))
< result
< column "(NOT 1 << 1 DIV 1) <=> ((NOT (1 << (1 DIV 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 << 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 << 0) <=> ((3 << 0))
< result
< column "(3 << 0) <=> ((3 << 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- - 0) <=> ((- (- 0)))
< result
< column "(- - 0) <=> ((- (- 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 3 AND 1 <> NULL NOT LIKE 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN 3 AND 1 <> NULL NOT LIKE 3) <=> (((3 NOT BETWEEN 3 AND 1) <> (NULL NOT LIKE 3)))
< result
< column "(3 NOT BETWEEN 3 AND 1 <> NULL NOT LIKE 3) <=> (((3 NOT BETWEEN 3 AND 1) <> (NULL NOT LIKE 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE NOT BETWEEN 2 MOD NULL AND 1 NOT BETWEEN 2 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "NOT BETWEEN 2 MOD NULL AND 1 NOT BETWEEN 2 AND NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL << 2 - 3 LIKE 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL << 2 - 3 LIKE 0) <=> ((((- NULL) << (2 - 3)) LIKE 0))
< result
< column "(- NULL << 2 - 3 LIKE 0) <=> ((((- NULL) << (2 - 3)) LIKE 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN 3 AND 0) <=> ((NULL BETWEEN 3 AND 0))
< result
< column "(NULL BETWEEN 3 AND 0) <=> ((NULL BETWEEN 3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT FALSE + 3 BETWEEN 2 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "+ 3 BETWEEN 2 AND 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 3 MOD 2 BETWEEN 1 AND 2 & 3 && 1 % 3 DIV 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 3 MOD 2 BETWEEN 1 AND 2 & 3 && 1 % 3 DIV 2) <=> ((1 XOR (((3 % 2) BETWEEN 1 AND (2 & 3)) AND ((1 % 3) DIV 2))))
< result
< column "(1 XOR 3 MOD 2 BETWEEN 1 AND 2 & 3 && 1 % 3 DIV 2) <=> ((1 XOR (((3 % 2) BETWEEN 1 AND (2 & 3)) AND ((1 % 3) DIV 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 != 1 >= NULL IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3 != 1 >= NULL IS NOT FALSE) <=> (((((~ 3) <> 1) >= NULL) IS NOT FALSE))
< result
< column "(~ 3 != 1 >= NULL IS NOT FALSE) <=> (((((~ 3) <> 1) >= NULL) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 3 IS NOT NULL / 2 IN (2, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "/ 2 IN (2, 1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && NULL) <=> ((0 AND NULL))
< result
< column "(0 && NULL) <=> ((0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> NULL < 1 NOT BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> NULL < 1 NOT BETWEEN NULL AND 1) <=> (((1 <=> NULL) < (1 NOT BETWEEN NULL AND 1)))
< result
< column "(1 <=> NULL < 1 NOT BETWEEN NULL AND 1) <=> (((1 <=> NULL) < (1 NOT BETWEEN NULL AND 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && 1 IS UNKNOWN) <=> ((0 AND (1 IS UNKNOWN)))
< result
< column "(0 && 1 IS UNKNOWN) <=> ((0 AND (1 IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 || 2) <=> ((3 OR 2))
< result
< column "(3 || 2) <=> ((3 OR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 / 2 >> 1 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 / 2 >> 1 <> NULL) <=> ((((2 / 2) >> 1) <> NULL))
< result
< column "(2 / 2 >> 1 <> NULL) <=> ((((2 / 2) >> 1) <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 ^ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 ^ 0) <=> ((0 ^ 0))
< result
< column "(0 ^ 0) <=> ((0 ^ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0) <=> ((NOT 0))
< result
< column "(NOT 0) <=> ((NOT 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 <=> 1 != 1 MOD 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 <=> 1 != 1 MOD 3 IS NOT NULL) <=> ((((0 <=> 1) <> (1 % 3)) IS NOT NULL))
< result
< column "(0 <=> 1 != 1 MOD 3 IS NOT NULL) <=> ((((0 <=> 1) <> (1 % 3)) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 AND NULL <= NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "NOT 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 * 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 * 0) <=> ((2 * 0))
< result
< column "(2 * 0) <=> ((2 * 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NOT NULL) <=> ((NULL IS NOT NULL))
< result
< column "(NULL IS NOT NULL) <=> ((NULL IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT IN (2, 3) > NULL IN (0, 2) NOT BETWEEN (2 && 0) AND - NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 46 near "NOT BETWEEN (2 && 0) AND - NULL IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * 1 < NULL <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 * 1 < NULL <=> 1) <=> ((((0 * 1) < NULL) <=> 1))
< result
< column "(0 * 1 < NULL <=> 1) <=> ((((0 * 1) < NULL) <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV 0 + 1 <=> 1 NOT IN (1 >= NULL << 3 < 3, 0 >= NULL DIV NULL = 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV 0 + 1 <=> 1 NOT IN (1 >= NULL << 3 < 3, 0 >= NULL DIV NULL = 2)) <=> ((((2 DIV 0) + 1) <=> (1 NOT IN (((1 >= (NULL << 3)) < 3), ((0 >= (NULL DIV NULL)) = 2)))))
< result
< column "(2 DIV 0 + 1 <=> 1 NOT IN (1 >= NULL << 3 < 3, 0 >= NULL DIV NULL = 2)) <=> ((((2 DIV 0) + 1) <=> (1 NOT IN (((1 >= (NULL << 3)) < 3), ((0 >= (NULL DIV NULL)) = 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 || 0 < ! 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 || 0 < ! 3) <=> ((3 OR (0 < (NOT 3))))
< result
< column "(3 || 0 < ! 3) <=> ((3 OR (0 < (NOT 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + 2) <=> ((3 + 2))
< result
< column "(3 + 2) <=> ((3 + 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 IN (+ 0, ~ 0) % + 2 IN (0 << 3, 1 BETWEEN 1 AND 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "% + 2 IN (0 << 3, 1 BETWEEN 1 AND 1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 && 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 && 3) <=> ((3 AND 3))
< result
< column "(3 && 3) <=> ((3 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 << 0 >= 2 MOD 0 NOT IN (0 + 1 && 2 BETWEEN 0 AND 3, + 1 & 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 << 0 >= 2 MOD 0 NOT IN (0 + 1 && 2 BETWEEN 0 AND 3, + 1 & 0)) <=> (((0 << 0) >= ((2 % 0) NOT IN (((0 + 1) AND (2 BETWEEN 0 AND 3)), ((+ 1) & 0)))))
< result
< column "(0 << 0 >= 2 MOD 0 NOT IN (0 + 1 && 2 BETWEEN 0 AND 3, + 1 & 0)) <=> (((0 << 0) >= ((2 % 0) NOT IN (((0 + 1) AND (2 BETWEEN 0 AND 3)), ((+ 1) & 0)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / 3 NOT BETWEEN @v := 2 AND (3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near ":= 2 AND (3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 / 1 BETWEEN 0 DIV 2 AND 2 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 / 1 BETWEEN 0 DIV 2 AND 2 <> 3) <=> ((((3 / 1) BETWEEN (0 DIV 2) AND 2) <> 3))
< result
< column "(3 / 1 BETWEEN 0 DIV 2 AND 2 <> 3) <=> ((((3 / 1) BETWEEN (0 DIV 2) AND 2) <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND 3) <=> ((NULL AND 3))
< result
< column "(NULL AND 3) <=> ((NULL AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL <> 0 BETWEEN 1 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT NULL <> 0 BETWEEN 1 AND 3) <=> (((3 IS NOT NULL) <> (0 BETWEEN 1 AND 3)))
< result
< column "(3 IS NOT NULL <> 0 BETWEEN 1 AND 3) <=> (((3 IS NOT NULL) <> (0 BETWEEN 1 AND 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3) <=> ((~ 3))
< result
< column "(~ 3) <=> ((~ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 >> NULL >> NULL + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 >> NULL >> NULL + 3) <=> ((((~ 2) >> NULL) >> (NULL + 3)))
< result
< column "(~ 2 >> NULL >> NULL + 3) <=> ((((~ 2) >> NULL) >> (NULL + 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 << 1 > NULL MOD 3 << 0 DIV 0 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 << 1 > NULL MOD 3 << 0 DIV 0 IS NOT FALSE) <=> ((((2 << 1) > ((NULL % 3) << (0 DIV 0))) IS NOT FALSE))
< result
< column "(2 << 1 > NULL MOD 3 << 0 DIV 0 IS NOT FALSE) <=> ((((2 << 1) > ((NULL % 3) << (0 DIV 0))) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT NULL) <=> ((1 IS NOT NULL))
< result
< column "(1 IS NOT NULL) <=> ((1 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (NULL, 2) - 3 || 1 / NOT 2 NOT REGEXP (NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "- 3 || 1 / NOT 2 NOT REGEXP (NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT BETWEEN 1 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT BETWEEN 1 AND 3) <=> ((0 NOT BETWEEN 1 AND 3))
< result
< column "(0 NOT BETWEEN 1 AND 3) <=> ((0 NOT BETWEEN 1 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 3 XOR 3 AND 3 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "XOR 3 AND 3 - 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 0 != (0) MOD - 0 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 0 != (0) MOD - 0 IS NOT NULL) <=> ((1 XOR ((0 <> (0 % (- 0))) IS NOT NULL)))
< result
< column "(1 XOR 0 != (0) MOD - 0 IS NOT NULL) <=> ((1 XOR ((0 <> (0 % (- 0))) IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! ~ 1 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! ~ 1 - 1) <=> (((NOT (~ 1)) - 1))
< result
< column "(! ~ 1 - 1) <=> (((NOT (~ 1)) - 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS NOT FALSE) <=> ((2 IS NOT FALSE))
< result
< column "(2 IS NOT FALSE) <=> ((2 IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 = 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 = 3 IS NOT NULL) <=> ((((+ 3) = 3) IS NOT NULL))
< result
< column "(+ 3 = 3 IS NOT NULL) <=> ((((+ 3) = 3) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 IS NOT FALSE & NULL = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "& NULL = 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 IN (3, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 IN (3, 2)) <=> ((NOT (1 IN (3, 2))))
< result
< column "(NOT 1 IN (3, 2)) <=> ((NOT (1 IN (3, 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 LIKE 0 + 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 17 near "+ 0 AND 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IN (0, 1) | 0 >= 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 20 near "| 0 >= 3 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 <=> 1) <=> (((NOT 3) <=> 1))
< result
< column "(! 3 <=> 1) <=> (((NOT 3) <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IN (1, 1) && 2 <= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IN (1, 1) && 2 <= 3) <=> (((NULL IN (1, 1)) AND (2 <= 3)))
< result
< column "(NULL IN (1, 1) && 2 <= 3) <=> (((NULL IN (1, 1)) AND (2 <= 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (0, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (0, 2)) <=> ((2 IN (0, 2)))
< result
< column "(2 IN (0, 2)) <=> ((2 IN (0, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS UNKNOWN) <=> ((3 IS UNKNOWN))
< result
< column "(3 IS UNKNOWN) <=> ((3 IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= 2) <=> ((3 >= 2))
< result
< column "(3 >= 2) <=> ((3 >= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <=> 1 % 3 IN (NULL, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <=> 1 % 3 IN (NULL, 0)) <=> ((3 <=> ((1 % 3) IN (NULL, 0))))
< result
< column "(3 <=> 1 % 3 IN (NULL, 0)) <=> ((3 <=> ((1 % 3) IN (NULL, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 != 0 * NULL LIKE NULL && ~ 0 - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 != 0 * NULL LIKE NULL && ~ 0 - 0) <=> (((3 <> ((0 * NULL) LIKE NULL)) AND ((~ 0) - 0)))
< result
< column "(3 != 0 * NULL LIKE NULL && ~ 0 - 0) <=> (((3 <> ((0 * NULL) LIKE NULL)) AND ((~ 0) - 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NOT 1 != 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 12 near "NOT 1 != 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 + NULL) <=> ((NOT (1 + NULL)))
< result
< column "(NOT 1 + NULL) <=> ((NOT (1 + NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 IN (2 NOT BETWEEN 0 AND 0, NOT 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 IN (2 NOT BETWEEN 0 AND 0, NOT 1)) <=> (((+ 3) IN ((2 NOT BETWEEN 0 AND 0), (NOT 1))))
< result
< column "(+ 3 IN (2 NOT BETWEEN 0 AND 0, NOT 1)) <=> (((+ 3) IN ((2 NOT BETWEEN 0 AND 0), (NOT 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 >= (1) BETWEEN 1 IN (NULL, 2) < 0 * 3 AND 2 IS NULL / 3 / 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "IN (NULL, 2) < 0 * 3 AND 2 IS NULL / 3 / 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 1 AND NULL IS TRUE IN (+ 2 BETWEEN NULL NOT BETWEEN 3 AND 3 AND 0 IS TRUE, (1) IS UNKNOWN)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "IN (+ 2 BETWEEN NULL NOT BETWEEN 3 AND 3 AND 0 IS TRUE, (1) IS UNKNOWN)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (1, NULL) IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (1, NULL) IS NULL) <=> (((2 IN (1, NULL)) IS NULL))
< result
< column "(2 IN (1, NULL) IS NULL) <=> (((2 IN (1, NULL)) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 REGEXP 1 >= 2 != 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 REGEXP 1 >= 2 != 3) <=> ((((0 REGEXP 1) >= 2) <> 3))
< result
< column "(0 REGEXP 1 >= 2 != 3) <=> ((((0 REGEXP 1) >= 2) <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 && NULL) <=> ((3 AND NULL))
< result
< column "(3 && NULL) <=> ((3 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (+ 0) MOD ! 1 NOT REGEXP 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((+ 0) MOD ! 1 NOT REGEXP 2) <=> ((((+ 0) % (NOT 1)) NOT REGEXP 2))
< result
< column "((+ 0) MOD ! 1 NOT REGEXP 2) <=> ((((+ 0) % (NOT 1)) NOT REGEXP 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2) <=> ((NOT 2))
< result
< column "(NOT 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 || 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 || 3 IS TRUE) <=> (((~ 2) OR (3 IS TRUE)))
< result
< column "(~ 2 || 3 IS TRUE) <=> (((~ 2) OR (3 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 OR 2) <=> ((2 OR 2))
< result
< column "(2 OR 2) <=> ((2 OR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV 2) <=> ((1 DIV 2))
< result
< column "(1 DIV 2) <=> ((1 DIV 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 > NULL | 0 != 1 <= + - 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 > NULL | 0 != 1 <= + - 3) <=> ((((1 > (NULL | 0)) <> 1) <= (+ (- 3))))
< result
< column "(1 > NULL | 0 != 1 <= + - 3) <=> ((((1 > (NULL | 0)) <> 1) <= (+ (- 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NULL >= 3 <> 1 && 3 NOT IN (3, NULL) < NULL IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS NULL >= 3 <> 1 && 3 NOT IN (3, NULL) < NULL IS TRUE) <=> (((((2 IS NULL) >= 3) <> 1) AND (((3 NOT IN (3, NULL)) < NULL) IS TRUE)))
< result
< column "(2 IS NULL >= 3 <> 1 && 3 NOT IN (3, NULL) < NULL IS TRUE) <=> (((((2 IS NULL) >= 3) <> 1) AND (((3 NOT IN (3, NULL)) < NULL) IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (3) ^ 1 BETWEEN 1 AND 1 + 3 / 0 <= - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((3) ^ 1 BETWEEN 1 AND 1 + 3 / 0 <= - NULL) <=> ((((3 ^ 1) BETWEEN 1 AND (1 + (3 / 0))) <= (- NULL)))
< result
< column "((3) ^ 1 BETWEEN 1 AND 1 + 3 / 0 <= - NULL) <=> ((((3 ^ 1) BETWEEN 1 AND (1 + (3 / 0))) <= (- NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL ^ NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 19 near "NOT 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 1 AND NULL >= 2 = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN 1 AND NULL >= 2 = 2) <=> ((((3 NOT BETWEEN 1 AND NULL) >= 2) = 2))
< result
< column "(3 NOT BETWEEN 1 AND NULL >= 2 = 2) <=> ((((3 NOT BETWEEN 1 AND NULL) >= 2) = 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 OR 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 OR 0 IS TRUE) <=> ((2 OR (0 IS TRUE)))
< result
< column "(2 OR 0 IS TRUE) <=> ((2 OR (0 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 0 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 0 AND 0) <=> ((NULL NOT BETWEEN 0 AND 0))
< result
< column "(NULL NOT BETWEEN 0 AND 0) <=> ((NULL NOT BETWEEN 0 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 REGEXP 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 REGEXP 1) <=> ((0 REGEXP 1))
< result
< column "(0 REGEXP 1) <=> ((0 REGEXP 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! + NULL ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! + NULL ^ 3) <=> (((NOT (+ NULL)) ^ 3))
< result
< column "(! + NULL ^ 3) <=> (((NOT (+ NULL)) ^ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 / 2 * 2 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 / 2 * 2 IS UNKNOWN) <=> ((((0 / 2) * 2) IS UNKNOWN))
< result
< column "(0 / 2 * 2 IS UNKNOWN) <=> ((((0 / 2) * 2) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1)) <=> (1)
< result
< column "((1)) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL / 3 NOT LIKE 0 / 0 && 1 + 0 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "/ 0 && 1 + 0 IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 || 2 NOT IN (NULL XOR 3, 0 && 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 || 2 NOT IN (NULL XOR 3, 0 && 2)) <=> ((1 OR (2 NOT IN ((NULL XOR 3), (0 AND 2)))))
< result
< column "(1 || 2 NOT IN (NULL XOR 3, 0 && 2)) <=> ((1 OR (2 NOT IN ((NULL XOR 3), (0 AND 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 > 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 > 1) <=> ((1 > 1))
< result
< column "(1 > 1) <=> ((1 > 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 ^ 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 ^ 1 IS NOT FALSE) <=> (((2 ^ 1) IS NOT FALSE))
< result
< column "(2 ^ 1 IS NOT FALSE) <=> (((2 ^ 1) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL >> 2) LIKE - NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL >> 2) LIKE - NULL IS NULL) <=> ((((NULL >> 2) LIKE (- NULL)) IS NULL))
< result
< column "((NULL >> 2) LIKE - NULL IS NULL) <=> ((((NULL >> 2) LIKE (- NULL)) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 < 3 NOT BETWEEN 2 <=> 3 AND 2 IS NOT NULL NOT BETWEEN NULL + 1 NOT LIKE (NULL) AND NOT NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "<=> 3 AND 2 IS NOT NULL NOT BETWEEN NULL + 1 NOT LIKE (NULL) AND NOT NOT 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL % NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL % NULL) <=> ((NULL % NULL))
< result
< column "(NULL % NULL) <=> ((NULL % NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 && 2 >> 3 BETWEEN 1 AND 1 NOT IN (! ! 0, ~ 2 IN (1 | NULL, 3 - 0))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 && 2 >> 3 BETWEEN 1 AND 1 NOT IN (! ! 0, ~ 2 IN (1 | NULL, 3 - 0))) <=> ((2 AND ((2 >> 3) BETWEEN 1 AND (1 NOT IN ((NOT (NOT 0)), ((~ 2) IN ((1 | NULL), (3 - 0))))))))
< result
< column "(2 && 2 >> 3 BETWEEN 1 AND 1 NOT IN (! ! 0, ~ 2 IN (1 | NULL, 3 - 0))) <=> ((2 AND ((2 >> 3) BETWEEN 1 AND (1 NOT IN ((NOT (NOT 0)), ((~ 2) IN ((1 | NULL), (3 - 0))))))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <> 0 NOT IN (3, NULL) IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <> 0 NOT IN (3, NULL) IS TRUE) <=> (((3 <> (0 NOT IN (3, NULL))) IS TRUE))
< result
< column "(3 <> 0 NOT IN (3, NULL) IS TRUE) <=> (((3 <> (0 NOT IN (3, NULL))) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3) <=> ((+ 3))
< result
< column "(+ 3) <=> ((+ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (2, 0) LIKE 0 <= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "LIKE 0 <= 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 < 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 < 3) <=> (((NOT 3) < 3))
< result
< column "(! 3 < 3) <=> (((NOT 3) < 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS NULL AND 2 IN (NULL, NULL) / ! 2 <= 0 < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "/ ! 2 <= 0 < 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 MOD 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 MOD 3) <=> (((- 1) % 3))
< result
< column "(- 1 MOD 3) <=> (((- 1) % 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 2) <=> ((1 XOR 2))
< result
< column "(1 XOR 2) <=> ((1 XOR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! (3) XOR 1 * NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! (3) XOR 1 * NULL) <=> (((NOT 3) XOR (1 * NULL)))
< result
< column "(! (3) XOR 1 * NULL) <=> (((NOT 3) XOR (1 * NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NOT NULL >> + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near ">> + 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL) <=> ((+ NULL))
< result
< column "(+ NULL) <=> ((+ NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 BETWEEN 2 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 BETWEEN 2 AND 1) <=> (((~ 2) BETWEEN 2 AND 1))
< result
< column "(~ 2 BETWEEN 2 AND 1) <=> (((~ 2) BETWEEN 2 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 = 3 < 2 >= 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 = 3 < 2 >= 1) <=> ((((1 = 3) < 2) >= 1))
< result
< column "(1 = 3 < 2 >= 1) <=> ((((1 = 3) < 2) >= 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (2, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 IN (2, 1)) <=> ((0 IN (2, 1)))
< result
< column "(0 IN (2, 1)) <=> ((0 IN (2, 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 NOT REGEXP NULL || 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 NOT REGEXP NULL || 1) <=> ((((NOT 0) NOT REGEXP NULL) OR 1))
< result
< column "(! 0 NOT REGEXP NULL || 1) <=> ((((NOT 0) NOT REGEXP NULL) OR 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 < 2 * NULL <> NULL << 2 << 0 IN (3 IS UNKNOWN, 2 < 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 < 2 * NULL <> NULL << 2 << 0 IN (3 IS UNKNOWN, 2 < 1)) <=> (((0 < (2 * NULL)) <> (((NULL << 2) << 0) IN ((3 IS UNKNOWN), (2 < 1)))))
< result
< column "(0 < 2 * NULL <> NULL << 2 << 0 IN (3 IS UNKNOWN, 2 < 1)) <=> (((0 < (2 * NULL)) <> (((NULL << 2) << 0) IN ((3 IS UNKNOWN), (2 < 1)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 BETWEEN 1 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 BETWEEN 1 AND NULL) <=> ((0 BETWEEN 1 AND NULL))
< result
< column "(0 BETWEEN 1 AND NULL) <=> ((0 BETWEEN 1 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL NOT BETWEEN 2 << 3 AND ! NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL NOT BETWEEN 2 << 3 AND ! NULL) <=> ((NOT (NULL NOT BETWEEN (2 << 3) AND (NOT NULL))))
< result
< column "(NOT NULL NOT BETWEEN 2 << 3 AND ! NULL) <=> ((NOT (NULL NOT BETWEEN (2 << 3) AND (NOT NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL ^ ~ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL ^ ~ 0) <=> (((+ NULL) ^ (~ 0)))
< result
< column "(+ NULL ^ ~ 0) <=> (((+ NULL) ^ (~ 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 = 1) <=> ((0 = 1))
< result
< column "(0 = 1) <=> ((0 = 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL | 0 BETWEEN NULL NOT BETWEEN 1 AND 3 AND 2 NOT BETWEEN 3 AND NULL <=> NOT 2 OR NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "NOT BETWEEN 1 AND 3 AND 2 NOT BETWEEN 3 AND NULL <=> NOT 2 OR NOT 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL MOD 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL MOD 3) <=> ((NULL % 3))
< result
< column "(NULL MOD 3) <=> ((NULL % 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3) <=> ((+ 3))
< result
< column "(+ 3) <=> ((+ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <= 1 LIKE 2 % 3 NOT IN (@v := NULL / 3, NULL && 0 IS NOT FALSE)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 22 near "% 3 NOT IN (@v := NULL / 3, NULL && 0 IS NOT FALSE)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <=> NULL NOT IN (1 DIV 2, NOT NULL) MOD - 3 = NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 48 near "MOD - 3 = NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN NULL AND 1 AND 3 IN (1, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN NULL AND 1 AND 3 IN (1, 0)) <=> (((NULL NOT BETWEEN NULL AND 1) AND (3 IN (1, 0))))
< result
< column "(NULL NOT BETWEEN NULL AND 1 AND 3 IN (1, 0)) <=> (((NULL NOT BETWEEN NULL AND 1) AND (3 IN (1, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL < 1) <=> ((NULL < 1))
< result
< column "(NULL < 1) <=> ((NULL < 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT LIKE 2 IN (2 IN (0, NULL), 2 NOT REGEXP 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 22 near "IN (2 IN (0, NULL), 2 NOT REGEXP 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 3 - 0 NOT BETWEEN 1 NOT BETWEEN 1 AND 2 AND 1 OR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 36 near "NOT BETWEEN 1 AND 2 AND 1 OR 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 | 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 | 3) <=> ((1 | 3))
< result
< column "(1 | 3) <=> ((1 | 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL >> 1) <=> ((NULL >> 1))
< result
< column "(NULL >> 1) <=> ((NULL >> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 && 3 IS NOT FALSE >= 2 = 2 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 33 near ">= 2 = 2 IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> 0 NOT BETWEEN NULL + 2 AND 2 ^ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> 0 NOT BETWEEN NULL + 2 AND 2 ^ NULL) <=> ((1 <=> (0 NOT BETWEEN (NULL + 2) AND (2 ^ NULL))))
< result
< column "(1 <=> 0 NOT BETWEEN NULL + 2 AND 2 ^ NULL) <=> ((1 <=> (0 NOT BETWEEN (NULL + 2) AND (2 ^ NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * 3 << 3 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 * 3 << 3 <> NULL) <=> ((((0 * 3) << 3) <> NULL))
< result
< column "(0 * 3 << 3 <> NULL) <=> ((((0 * 3) << 3) <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN NULL AND 2 <> 3 NOT BETWEEN NULL AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN NULL AND 2 <> 3 NOT BETWEEN NULL AND 2) <=> (((NULL NOT BETWEEN NULL AND 2) <> (3 NOT BETWEEN NULL AND 2)))
< result
< column "(NULL NOT BETWEEN NULL AND 2 <> 3 NOT BETWEEN NULL AND 2) <=> (((NULL NOT BETWEEN NULL AND 2) <> (3 NOT BETWEEN NULL AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 - NULL NOT IN (NULL, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 - NULL NOT IN (NULL, NULL)) <=> ((((- 1) - NULL) NOT IN (NULL, NULL)))
< result
< column "(- 1 - NULL NOT IN (NULL, NULL)) <=> ((((- 1) - NULL) NOT IN (NULL, NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 * 0 / 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3 * 0 / 0 IS TRUE) <=> (((((~ 3) * 0) / 0) IS TRUE))
< result
< column "(~ 3 * 0 / 0 IS TRUE) <=> (((((~ 3) * 0) / 0) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 % NULL NOT IN (2 / 0, 1 NOT BETWEEN 2 AND NULL) ^ 1 <= 3 | 2 <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 58 near "^ 1 <= 3 | 2 <= NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (NULL, 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 IN (NULL, 3)) <=> ((0 IN (NULL, 3)))
< result
< column "(0 IN (NULL, 3)) <=> ((0 IN (NULL, 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((2)) <=> (2)
< result
< column "((2)) <=> (2)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < 0) <=> ((1 < 0))
< result
< column "(1 < 0) <=> ((1 < 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 3 NOT IN (3, 0) / 3 IS UNKNOWN <=> 1 IN (2, 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "/ 3 IS UNKNOWN <=> 1 IN (2, 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 > 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 > 0) <=> ((2 > 0))
< result
< column "(2 > 0) <=> ((2 > 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 - 3) <=> ((3 - 3))
< result
< column "(3 - 3) <=> ((3 - 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS TRUE & NULL LIKE 2 NOT BETWEEN 1 = NULL >= 3 AND (2 NOT IN (2, 3))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 18 near "& NULL LIKE 2 NOT BETWEEN 1 = NULL >= 3 AND (2 NOT IN (2, 3))" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | ~ 2) <=> ((0 | (~ 2)))
< result
< column "(0 | ~ 2) <=> ((0 | (~ 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= 1) <=> ((3 <= 1))
< result
< column "(3 <= 1) <=> ((3 <= 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 3 NOT REGEXP NOT 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "NOT 3 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | NULL < NULL | 1 OR 3 <=> 3 ^ 1 < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | NULL < NULL | 1 OR 3 <=> 3 ^ 1 < NULL) <=> ((((0 | NULL) < (NULL | 1)) OR ((3 <=> (3 ^ 1)) < NULL)))
< result
< column "(0 | NULL < NULL | 1 OR 3 <=> 3 ^ 1 < NULL) <=> ((((0 | NULL) < (NULL | 1)) OR ((3 <=> (3 ^ 1)) < NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 1 AND 3 <=> NULL % 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 1 AND 3 <=> NULL % 2) <=> (((2 NOT BETWEEN 1 AND 3) <=> (NULL % 2)))
< result
< column "(2 NOT BETWEEN 1 AND 3 <=> NULL % 2) <=> (((2 NOT BETWEEN 1 AND 3) <=> (NULL % 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0) <=> (0)
< result
< column "(0) <=> (0)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 REGEXP 0 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 REGEXP 0 IS UNKNOWN) <=> (((3 REGEXP 0) IS UNKNOWN))
< result
< column "(3 REGEXP 0 IS UNKNOWN) <=> (((3 REGEXP 0) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NULL & NULL IS NOT NULL && NULL DIV 3 IN (0 << 1, 2 BETWEEN 1 AND NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 18 near "& NULL IS NOT NULL && NULL DIV 3 IN (0 << 1, 2 BETWEEN 1 AND NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= 2) <=> ((2 >= 2))
< result
< column "(2 >= 2) <=> ((2 >= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT NULL) <=> ((1 IS NOT NULL))
< result
< column "(1 IS NOT NULL) <=> ((1 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 3 & 3 IN (2 && 1, 2 NOT BETWEEN 2 AND NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 3 & 3 IN (2 && 1, 2 NOT BETWEEN 2 AND NULL)) <=> ((((- 3) & 3) IN ((2 AND 1), (2 NOT BETWEEN 2 AND NULL))))
< result
< column "(- 3 & 3 IN (2 && 1, 2 NOT BETWEEN 2 AND NULL)) <=> ((((- 3) & 3) IN ((2 AND 1), (2 NOT BETWEEN 2 AND NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + + 0 + - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ + 0 + - NULL) <=> (((+ (+ 0)) + (- NULL)))
< result
< column "(+ + 0 + - NULL) <=> (((+ (+ 0)) + (- NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 3 AND 1 & ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN 3 AND 1 & ! 2) <=> ((3 NOT BETWEEN 3 AND (1 & (NOT 2))))
< result
< column "(3 NOT BETWEEN 3 AND 1 & ! 2) <=> ((3 NOT BETWEEN 3 AND (1 & (NOT 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT LIKE 0 && 2 NOT IN (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT LIKE 0 && 2 NOT IN (1, 2)) <=> (((3 NOT LIKE 0) AND (2 NOT IN (1, 2))))
< result
< column "(3 NOT LIKE 0 && 2 NOT IN (1, 2)) <=> (((3 NOT LIKE 0) AND (2 NOT IN (1, 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 ^ 2 & 3 >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 ^ 2 & 3 >> 1) <=> (((2 ^ 2) & (3 >> 1)))
< result
< column "(2 ^ 2 & 3 >> 1) <=> (((2 ^ 2) & (3 >> 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 2 AND 3) <=> ((3 BETWEEN 2 AND 3))
< result
< column "(3 BETWEEN 2 AND 3) <=> ((3 BETWEEN 2 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 * 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 * 2) <=> ((2 * 2))
< result
< column "(2 * 2) <=> ((2 * 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 - 3 < ~ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2 - 3 < ~ 0) <=> ((((NOT 2) - 3) < (~ 0)))
< result
< column "(! 2 - 3 < ~ 0) <=> ((((NOT 2) - 3) < (~ 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 / NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 / NULL) <=> ((2 / NULL))
< result
< column "(2 / NULL) <=> ((2 / NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 3 != 3 NOT IN (0 | NULL, NOT 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 3 != 3 NOT IN (0 | NULL, NOT 3)) <=> (((- 3) <> (3 NOT IN ((0 | NULL), (NOT 3)))))
< result
< column "(- 3 != 3 NOT IN (0 | NULL, NOT 3)) <=> (((- 3) <> (3 NOT IN ((0 | NULL), (NOT 3)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL < 0 MOD 1 > 2 IS UNKNOWN MOD NULL && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 42 near "MOD NULL && NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 << 3) <=> ((1 << 3))
< result
< column "(1 << 3) <=> ((1 << 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 >= 1 >> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 >= 1 >> NULL) <=> ((NOT (3 >= (1 >> NULL))))
< result
< column "(NOT 3 >= 1 >> NULL) <=> ((NOT (3 >= (1 >> NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3) <=> ((NOT 3))
< result
< column "(NOT 3) <=> ((NOT 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS TRUE REGEXP NOT 1 >= 1 > 1 AND - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "REGEXP NOT 1 >= 1 > 1 AND - 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 XOR 2 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 XOR 2 IS NULL) <=> (((+ 3) XOR (2 IS NULL)))
< result
< column "(+ 3 XOR 2 IS NULL) <=> (((+ 3) XOR (2 IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN NULL AND 3 - 0 <=> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN NULL AND 3 - 0 <=> NULL) <=> (((NULL NOT BETWEEN NULL AND (3 - 0)) <=> NULL))
< result
< column "(NULL NOT BETWEEN NULL AND 3 - 0 <=> NULL) <=> (((NULL NOT BETWEEN NULL AND (3 - 0)) <=> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1) <=> ((- 1))
< result
< column "(- 1) <=> ((- 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT FALSE) <=> ((3 IS NOT FALSE))
< result
< column "(3 IS NOT FALSE) <=> ((3 IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE IN (NULL IS NULL, NULL NOT IN (1, 3))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "IN (NULL IS NULL, NULL NOT IN (1, 3))" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 REGEXP 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 REGEXP 0 IS TRUE) <=> (((0 REGEXP 0) IS TRUE))
< result
< column "(0 REGEXP 0 IS TRUE) <=> (((0 REGEXP 0) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 | 2) <=> ((2 | 2))
< result
< column "(2 | 2) <=> ((2 | 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> NULL * 3 && 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <> NULL * 3 && 1 IS UNKNOWN) <=> (((1 <> (NULL * 3)) AND (1 IS UNKNOWN)))
< result
< column "(1 <> NULL * 3 && 1 IS UNKNOWN) <=> (((1 <> (NULL * 3)) AND (1 IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS TRUE NOT LIKE 3 BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 20 near "NOT LIKE 3 BETWEEN 0 AND NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (3) - 2 NOT LIKE 1 << NULL IS NULL BETWEEN NULL XOR 0 AND 0 << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "<< NULL IS NULL BETWEEN NULL XOR 0 AND 0 << 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 ^ 1 MOD 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0 ^ 1 MOD 2) <=> ((((~ 0) ^ 1) % 2))
< result
< column "(~ 0 ^ 1 MOD 2) <=> ((((~ 0) ^ 1) % 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 & 1 IN (0 && NULL, 1 NOT LIKE 2) || NOT 0 = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 & 1 IN (0 && NULL, 1 NOT LIKE 2) || NOT 0 = 2) <=> ((((2 & 1) IN ((0 AND NULL), (1 NOT LIKE 2))) OR (NOT (0 = 2))))
< result
< column "(2 & 1 IN (0 && NULL, 1 NOT LIKE 2) || NOT 0 = 2) <=> ((((2 & 1) IN ((0 AND NULL), (1 NOT LIKE 2))) OR (NOT (0 = 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 IS NOT NULL BETWEEN NULL % 2 AND 2 >> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "BETWEEN NULL % 2 AND 2 >> 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / 3 REGEXP 1 OR NULL XOR 0 >= 2 DIV @v := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 50 near ":= 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 DIV NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 DIV NULL IS NOT NULL) <=> (((0 DIV NULL) IS NOT NULL))
< result
< column "(0 DIV NULL IS NOT NULL) <=> (((0 DIV NULL) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (3) % 1 NOT LIKE 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((3) % 1 NOT LIKE 0) <=> (((3 % 1) NOT LIKE 0))
< result
< column "((3) % 1 NOT LIKE 0) <=> (((3 % 1) NOT LIKE 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 1 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 1 AND 1) <=> ((2 NOT BETWEEN 1 AND 1))
< result
< column "(2 NOT BETWEEN 1 AND 1) <=> ((2 NOT BETWEEN 1 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL IN (2 | 1, NULL REGEXP 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "IN (2 | 1, NULL REGEXP 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 + 0) <=> ((1 + 0))
< result
< column "(1 + 0) <=> ((1 + 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 IS UNKNOWN) <=> (((+ 3) IS UNKNOWN))
< result
< column "(+ 3 IS UNKNOWN) <=> (((+ 3) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT FALSE) <=> ((3 IS NOT FALSE))
< result
< column "(3 IS NOT FALSE) <=> ((3 IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 BETWEEN 0 AND 2 - 2 = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1 BETWEEN 0 AND 2 - 2 = 2) <=> ((((+ 1) BETWEEN 0 AND (2 - 2)) = 2))
< result
< column "(+ 1 BETWEEN 0 AND 2 - 2 = 2) <=> ((((+ 1) BETWEEN 0 AND (2 - 2)) = 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0) <=> ((+ 0))
< result
< column "(+ 0) <=> ((+ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 > 1 NOT BETWEEN NULL / 1 IS UNKNOWN AND ~ 2 MOD 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "IS UNKNOWN AND ~ 2 MOD 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS TRUE) <=> ((1 IS TRUE))
< result
< column "(1 IS TRUE) <=> ((1 IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <= 0 & 1 IS NOT FALSE > NOT 1 ^ 0 << NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "> NOT 1 ^ 0 << NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 || 2) <=> ((1 OR 2))
< result
< column "(1 || 2) <=> ((1 OR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0 LIKE 3 % 2 NOT IN (2, 2) << NULL OR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "% 2 NOT IN (2, 2) << NULL OR 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 NOT BETWEEN 3 AND 2 || ! 0 < 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 NOT BETWEEN 3 AND 2 || ! 0 < 0) <=> (((NOT (2 NOT BETWEEN 3 AND 2)) OR ((NOT 0) < 0)))
< result
< column "(NOT 2 NOT BETWEEN 3 AND 2 || ! 0 < 0) <=> (((NOT (2 NOT BETWEEN 3 AND 2)) OR ((NOT 0) < 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 3 AND NULL ^ NULL <=> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 3 AND NULL ^ NULL <=> 0) <=> (((NULL NOT BETWEEN 3 AND (NULL ^ NULL)) <=> 0))
< result
< column "(NULL NOT BETWEEN 3 AND NULL ^ NULL <=> 0) <=> (((NULL NOT BETWEEN 3 AND (NULL ^ NULL)) <=> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 << NULL % 1 IN (3, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 << NULL % 1 IN (3, 1)) <=> ((NOT ((3 << (NULL % 1)) IN (3, 1))))
< result
< column "(NOT 3 << NULL % 1 IN (3, 1)) <=> ((NOT ((3 << (NULL % 1)) IN (3, 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 < NULL IS NOT FALSE NOT BETWEEN 0 IN (NULL, 1) / 2 NOT LIKE 3 AND 3 IS NULL LIKE NULL + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 34 near "NOT BETWEEN 0 IN (NULL, 1) / 2 NOT LIKE 3 AND 3 IS NULL LIKE NULL + 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV 0) <=> ((1 DIV 0))
< result
< column "(1 DIV 0) <=> ((1 DIV 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 && 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 && 3) <=> ((2 AND 3))
< result
< column "(2 && 3) <=> ((2 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL MOD 1 MOD NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL MOD 1 MOD NULL) <=> ((((+ NULL) % 1) % NULL))
< result
< column "(+ NULL MOD 1 MOD NULL) <=> ((((+ NULL) % 1) % NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV NULL IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV NULL IS TRUE) <=> (((1 DIV NULL) IS TRUE))
< result
< column "(1 DIV NULL IS TRUE) <=> (((1 DIV NULL) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 | 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 | 1) <=> ((1 | 1))
< result
< column "(1 | 1) <=> ((1 | 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT IN (2, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT IN (2, 2)) <=> ((NULL NOT IN (2, 2)))
< result
< column "(NULL NOT IN (2, 2)) <=> ((NULL NOT IN (2, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= 1 MOD 1 IN (3, 0) IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= 1 MOD 1 IN (3, 0) IS UNKNOWN) <=> (((3 <= ((1 % 1) IN (3, 0))) IS UNKNOWN))
< result
< column "(3 <= 1 MOD 1 IN (3, 0) IS UNKNOWN) <=> (((3 <= ((1 % 1) IN (3, 0))) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >= 0 BETWEEN NULL <=> 1 AND NULL OR 2 >= 1 || 2 + 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "<=> 1 AND NULL OR 2 >= 1 || 2 + 0 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 3 AND 2 > NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "NOT 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 - 3 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1 - 3 IS NULL) <=> ((((+ 1) - 3) IS NULL))
< result
< column "(+ 1 - 3 IS NULL) <=> ((((+ 1) - 3) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 % 2 NOT IN ((NULL), 2 + NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 % 2 NOT IN ((NULL), 2 + NULL)) <=> (((2 % 2) NOT IN (NULL, (2 + NULL))))
< result
< column "(2 % 2 NOT IN ((NULL), 2 + NULL)) <=> (((2 % 2) NOT IN (NULL, (2 + NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <= 0 NOT REGEXP 0 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <= 0 NOT REGEXP 0 IS UNKNOWN) <=> (((2 <= (0 NOT REGEXP 0)) IS UNKNOWN))
< result
< column "(2 <= 0 NOT REGEXP 0 IS UNKNOWN) <=> (((2 <= (0 NOT REGEXP 0)) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS TRUE AND 3 NOT LIKE 0 NOT IN (NULL OR 1 LIKE 0 DIV 3, 3 BETWEEN NULL - 2 AND 2 <=> NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "NOT IN (NULL OR 1 LIKE 0 DIV 3, 3 BETWEEN NULL - 2 AND 2 <=> NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2) <=> ((+ 2))
< result
< column "(+ 2) <=> ((+ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 2 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 2 AND 0) <=> ((2 BETWEEN 2 AND 0))
< result
< column "(2 BETWEEN 2 AND 0) <=> ((2 BETWEEN 2 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS NULL) <=> ((2 IS NULL))
< result
< column "(2 IS NULL) <=> ((2 IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 <= NULL) <=> ((0 <= NULL))
< result
< column "(0 <= NULL) <=> ((0 <= NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 > 3 NOT BETWEEN 1 % 0 AND 1 <> 2 - 0 NOT IN (NULL, 0) NOT BETWEEN 2 BETWEEN 0 AND 3 AND 1 NOT IN (0, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 66 near "NOT BETWEEN 2 BETWEEN 0 AND 3 AND 1 NOT IN (0, 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 & 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3 & 1) <=> (((~ 3) & 1))
< result
< column "(~ 3 & 1) <=> (((~ 3) & 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && 2) <=> ((0 AND 2))
< result
< column "(0 && 2) <=> ((0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 2 AND 3 <=> 1 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN 2 AND 3 <=> 1 OR 2) <=> ((((1 BETWEEN 2 AND 3) <=> 1) OR 2))
< result
< column "(1 BETWEEN 2 AND 3 <=> 1 OR 2) <=> ((((1 BETWEEN 2 AND 3) <=> 1) OR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL % 0 IN ((2), 2 IS TRUE) ^ 0 AND 2 ^ 2 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "^ 0 AND 2 ^ 2 || 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NOT NULL) <=> ((NULL IS NOT NULL))
< result
< column "(NULL IS NOT NULL) <=> ((NULL IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0 IN (0, NULL) = - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 0 IN (0, NULL) = - 1) <=> ((@v := ((0 IN (0, NULL)) = (- 1))))
< result
< column "(@v := 0 IN (0, NULL) = - 1) <=> ((@v := ((0 IN (0, NULL)) = (- 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT REGEXP 1 BETWEEN 1 IS TRUE AND NOT 1 != 1 IS TRUE <> - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "BETWEEN 1 IS TRUE AND NOT 1 != 1 IS TRUE <> - 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 + NULL <> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 + NULL <> 2) <=> ((NOT ((2 + NULL) <> 2)))
< result
< column "(NOT 2 + NULL <> 2) <=> ((NOT ((2 + NULL) <> 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (NULL, 2) + ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "+ ~ 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL) <=> (NULL)
< result
< column "(NULL) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL AND 0 >> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL AND 0 >> 0) <=> (((NOT NULL) AND (0 >> 0)))
< result
< column "(! NULL AND 0 >> 0) <=> (((NOT NULL) AND (0 >> 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0 | 0 DIV 3 - 0 / 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 0 | 0 DIV 3 - 0 / 1 IS NOT NULL) <=> ((@v := ((0 | ((0 DIV 3) - (0 / 1))) IS NOT NULL)))
< result
< column "(@v := 0 | 0 DIV 3 - 0 / 1 IS NOT NULL) <=> ((@v := ((0 | ((0 DIV 3) - (0 / 1))) IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 1 AND 3 BETWEEN 3 NOT BETWEEN 1 AND 3 AND 1 REGEXP NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "NOT BETWEEN 1 AND 3 AND 1 REGEXP NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ NULL % 3 && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ NULL % 3 && NULL) <=> ((((~ NULL) % 3) AND NULL))
< result
< column "(~ NULL % 3 && NULL) <=> ((((~ NULL) % 3) AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 IS NOT FALSE NOT LIKE NULL NOT REGEXP 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "NOT LIKE NULL NOT REGEXP 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 2 < - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= 2 < - 1) <=> (((3 >= 2) < (- 1)))
< result
< column "(3 >= 2 < - 1) <=> (((3 >= 2) < (- 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 3 IS TRUE) <=> ((1 XOR (3 IS TRUE)))
< result
< column "(1 XOR 3 IS TRUE) <=> ((1 XOR (3 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ ~ NULL <=> 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ ~ NULL <=> 3 IS TRUE) <=> ((((~ (~ NULL)) <=> 3) IS TRUE))
< result
< column "(~ ~ NULL <=> 3 IS TRUE) <=> ((((~ (~ NULL)) <=> 3) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (2, 2) BETWEEN NOT 1 AND 1 BETWEEN 3 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "BETWEEN NOT 1 AND 1 BETWEEN 3 AND 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (2, 1) && NULL IN (NULL, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (2, 1) && NULL IN (NULL, 1)) <=> (((2 IN (2, 1)) AND (NULL IN (NULL, 1))))
< result
< column "(2 IN (2, 1) && NULL IN (NULL, 1)) <=> (((2 IN (2, 1)) AND (NULL IN (NULL, 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 ^ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1 ^ 1) <=> (((NOT 1) ^ 1))
< result
< column "(! 1 ^ 1) <=> (((NOT 1) ^ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 2) <=> ((- 2))
< result
< column "(- 2) <=> ((- 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 XOR 3 = - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 XOR 3 = - NULL) <=> ((0 XOR (3 = (- NULL))))
< result
< column "(0 XOR 3 = - NULL) <=> ((0 XOR (3 = (- NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <> NULL && 1 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <> NULL && 1 IS NULL) <=> (((2 <> NULL) AND (1 IS NULL)))
< result
< column "(2 <> NULL && 1 IS NULL) <=> (((2 <> NULL) AND (1 IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL OR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL OR 0) <=> ((NULL OR 0))
< result
< column "(NULL OR 0) <=> ((NULL OR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (1, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (1, NULL)) <=> ((0 NOT IN (1, NULL)))
< result
< column "(0 NOT IN (1, NULL)) <=> ((0 NOT IN (1, NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 && 3 IS NULL IN (0 IS NOT FALSE >= 2 + 2, + 0 BETWEEN @v := 2 AND NULL XOR NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "IN (0 IS NOT FALSE >= 2 + 2, + 0 BETWEEN @v := 2 AND NULL XOR NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3 <=> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3 <=> 3) <=> (((+ 3) <=> 3))
< result
< column "(+ 3 <=> 3) <=> (((+ 3) <=> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV 0 * 1 && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV 0 * 1 && 2) <=> ((((2 DIV 0) * 1) AND 2))
< result
< column "(2 DIV 0 * 1 && 2) <=> ((((2 DIV 0) * 1) AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT FALSE) <=> ((1 IS NOT FALSE))
< result
< column "(1 IS NOT FALSE) <=> ((1 IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 0 % NULL <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 0 % NULL <= NULL) <=> (((NULL = (0 % NULL)) <= NULL))
< result
< column "(NULL = 0 % NULL <= NULL) <=> (((NULL = (0 % NULL)) <= NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND NULL) <=> ((NULL AND NULL))
< result
< column "(NULL AND NULL) <=> ((NULL AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 NOT BETWEEN 1 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 NOT BETWEEN 1 AND 2) <=> ((NOT (0 NOT BETWEEN 1 AND 2)))
< result
< column "(NOT 0 NOT BETWEEN 1 AND 2) <=> ((NOT (0 NOT BETWEEN 1 AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 3 OR 3 >> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= 3 OR 3 >> 2) <=> (((2 >= 3) OR (3 >> 2)))
< result
< column "(2 >= 3 OR 3 >> 2) <=> (((2 >= 3) OR (3 >> 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2) <=> ((NOT 2))
< result
< column "(NOT 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 && 2 / NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 && 2 / NULL AND 0) <=> (((3 AND (2 / NULL)) AND 0))
< result
< column "(3 && 2 / NULL AND 0) <=> (((3 AND (2 / NULL)) AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 XOR 3 < ! 3 >= 2 IS NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 XOR 3 < ! 3 >= 2 IS NULL IS NULL) <=> ((2 XOR ((((3 < (NOT 3)) >= 2) IS NULL) IS NULL)))
< result
< column "(2 XOR 3 < ! 3 >= 2 IS NULL IS NULL) <=> ((2 XOR ((((3 < (NOT 3)) >= 2) IS NULL) IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= 0) <=> ((3 <= 0))
< result
< column "(3 <= 0) <=> ((3 <= 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2) <=> (2)
< result
< column "(2) <=> (2)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IN (0, 2) >> + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near ">> + 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >> NULL) <=> ((2 >> NULL))
< result
< column "(2 >> NULL) <=> ((2 >> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 AND 0) <=> ((3 AND 0))
< result
< column "(3 AND 0) <=> ((3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 DIV 0 OR NULL IN (+ 3 OR 3, 0 REGEXP NULL > 3 IS NOT FALSE)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0 DIV 0 OR NULL IN (+ 3 OR 3, 0 REGEXP NULL > 3 IS NOT FALSE)) <=> ((((+ 0) DIV 0) OR (NULL IN (((+ 3) OR 3), (((0 REGEXP NULL) > 3) IS NOT FALSE)))))
< result
< column "(+ 0 DIV 0 OR NULL IN (+ 3 OR 3, 0 REGEXP NULL > 3 IS NOT FALSE)) <=> ((((+ 0) DIV 0) OR (NULL IN (((+ 3) OR 3), (((0 REGEXP NULL) > 3) IS NOT FALSE)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 = NULL MOD 2 BETWEEN 1 AND NULL <> 2 NOT IN (NULL, 1) - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 64 near "- NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1) <=> ((NOT 1))
< result
< column "(! 1) <=> ((NOT 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <= 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <= 1) <=> ((1 <= 1))
< result
< column "(1 <= 1) <=> ((1 <= 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 REGEXP NULL IN (3 REGEXP 2, NULL >= 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "IN (3 REGEXP 2, NULL >= 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT IN (3, NULL) NOT BETWEEN 1 NOT IN (3, 0) AND + 0 IN (3 && 3 < 3 || 1, 0 IN (1, 0) >= 0 REGEXP 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "NOT BETWEEN 1 NOT IN (3, 0) AND + 0 IN (3 && 3 < 3 || 1, 0 IN (1, 0) >= 0 REGEXP 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL NOT REGEXP NULL BETWEEN 1 AND 0 BETWEEN 2 - 0 >= 3 AND @v := 2 IN (1 IS NULL, 0 <=> 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "BETWEEN 1 AND 0 BETWEEN 2 - 0 >= 3 AND @v := 2 IN (1 IS NULL, 0 <=> 1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 0 OR 2) <=> ((@v := (0 OR 2)))
< result
< column "(@v := 0 OR 2) <=> ((@v := (0 OR 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS NOT FALSE ^ NULL ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "^ NULL ^ 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 + 1 || 1 > 0 NOT BETWEEN 2 NOT BETWEEN 2 AND 0 <= 2 DIV 0 AND 0 || NULL << NULL < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 39 near "NOT BETWEEN 2 AND 0 <= 2 DIV 0 AND 0 || NULL << NULL < 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL ^ 2) <=> ((NULL ^ 2))
< result
< column "(NULL ^ 2) <=> ((NULL ^ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 NOT BETWEEN NULL AND NULL IN (0 NOT IN (3, NULL) | - 2, + 0 BETWEEN 3 AND 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 61 near "| - 2, + 0 BETWEEN 3 AND 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= NULL) <=> ((3 <= NULL))
< result
< column "(3 <= NULL) <=> ((3 <= NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 AND 3 >> NULL DIV 3 != (+ 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 AND 3 >> NULL DIV 3 != (+ 3)) <=> ((2 AND ((3 >> (NULL DIV 3)) <> (+ 3))))
< result
< column "(2 AND 3 >> NULL DIV 3 != (+ 3)) <=> ((2 AND ((3 >> (NULL DIV 3)) <> (+ 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + 3) <=> ((3 + 3))
< result
< column "(3 + 3) <=> ((3 + 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT BETWEEN 3 AND 0) <=> ((1 NOT BETWEEN 3 AND 0))
< result
< column "(1 NOT BETWEEN 3 AND 0) <=> ((1 NOT BETWEEN 3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ ! 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ ! 0 IS TRUE) <=> (((~ (NOT 0)) IS TRUE))
< result
< column "(~ ! 0 IS TRUE) <=> (((~ (NOT 0)) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (NULL, 1) != 2 < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (NULL, 1) != 2 < 1) <=> ((((0 NOT IN (NULL, 1)) <> 2) < 1))
< result
< column "(0 NOT IN (NULL, 1) != 2 < 1) <=> ((((0 NOT IN (NULL, 1)) <> 2) < 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 0 != 0 <= NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 0 != 0 <= NULL) <=> ((((- 0) <> 0) <= NULL))
< result
< column "(- 0 != 0 <= NULL) <=> ((((- 0) <> 0) <= NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL || 0 NOT IN (2, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL || 0 NOT IN (2, 0)) <=> (((NOT NULL) OR (0 NOT IN (2, 0))))
< result
< column "(! NULL || 0 NOT IN (2, 0)) <=> (((NOT NULL) OR (0 NOT IN (2, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <= NULL && + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <= NULL && + 1) <=> (((NULL <= NULL) AND (+ 1)))
< result
< column "(NULL <= NULL && + 1) <=> (((NULL <= NULL) AND (+ 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 LIKE 2 = 1 = NULL && NULL / 2 && NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 LIKE 2 = 1 = NULL && NULL / 2 && NULL IS NULL) <=> ((((((3 LIKE 2) = 1) = NULL) AND (NULL / 2)) AND (NULL IS NULL)))
< result
< column "(3 LIKE 2 = 1 = NULL && NULL / 2 && NULL IS NULL) <=> ((((((3 LIKE 2) = 1) = NULL) AND (NULL / 2)) AND (NULL IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 1 != 1 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= 1 != 1 % 0) <=> (((3 >= 1) <> (1 % 0)))
< result
< column "(3 >= 1 != 1 % 0) <=> (((3 >= 1) <> (1 % 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NOT FALSE) <=> ((NULL IS NOT FALSE))
< result
< column "(NULL IS NOT FALSE) <=> ((NULL IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL MOD 1 && 0 <= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL MOD 1 && 0 <= 3) <=> (((NULL % 1) AND (0 <= 3)))
< result
< column "(NULL MOD 1 && 0 <= 3) <=> (((NULL % 1) AND (0 <= 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 = 1 || (0) DIV 0 AND NULL / (0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 = 1 || (0) DIV 0 AND NULL / (0)) <=> (((1 = 1) OR ((0 DIV 0) AND (NULL / 0))))
< result
< column "(1 = 1 || (0) DIV 0 AND NULL / (0)) <=> (((1 = 1) OR ((0 DIV 0) AND (NULL / 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (NULL, 3) IS NOT NULL & 3 <= 3 || 2 BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "& 3 <= 3 || 2 BETWEEN NULL AND 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 - 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 - 3) <=> ((2 - 3))
< result
< column "(2 - 3) <=> ((2 - 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 >= 3 LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 >= 3 LIKE 2) <=> ((NOT (0 >= (3 LIKE 2))))
< result
< column "(NOT 0 >= 3 LIKE 2) <=> ((NOT (0 >= (3 LIKE 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 > 0 DIV 1 BETWEEN 2 AND 0 IN (! NULL IS NOT FALSE, (NULL NOT IN (3, 0)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 > 0 DIV 1 BETWEEN 2 AND 0 IN (! NULL IS NOT FALSE, (NULL NOT IN (3, 0)))) <=> ((2 > ((0 DIV 1) BETWEEN 2 AND (0 IN (((NOT NULL) IS NOT FALSE), (NULL NOT IN (3, 0)))))))
< result
< column "(2 > 0 DIV 1 BETWEEN 2 AND 0 IN (! NULL IS NOT FALSE, (NULL NOT IN (3, 0)))) <=> ((2 > ((0 DIV 1) BETWEEN 2 AND (0 IN (((NOT NULL) IS NOT FALSE), (NULL NOT IN (3, 0)))))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 2 BETWEEN 1 XOR 2 AND NULL ^ 0 <= 1 | 0 AND NULL NOT LIKE 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "XOR 2 AND NULL ^ 0 <= 1 | 0 AND NULL NOT LIKE 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 NOT IN (1 REGEXP NULL, 0 ^ 0) IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 NOT IN (1 REGEXP NULL, 0 ^ 0) IS UNKNOWN) <=> ((((NOT 0) NOT IN ((1 REGEXP NULL), (0 ^ 0))) IS UNKNOWN))
< result
< column "(! 0 NOT IN (1 REGEXP NULL, 0 ^ 0) IS UNKNOWN) <=> ((((NOT 0) NOT IN ((1 REGEXP NULL), (0 ^ 0))) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 MOD 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 MOD 3) <=> ((3 % 3))
< result
< column "(3 MOD 3) <=> ((3 % 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 IS NULL) <=> ((0 IS NULL))
< result
< column "(0 IS NULL) <=> ((0 IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 || NULL % 2 & 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 || NULL % 2 & 3) <=> ((2 OR ((NULL % 2) & 3)))
< result
< column "(2 || NULL % 2 & 3) <=> ((2 OR ((NULL % 2) & 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 NOT IN (+ NULL, 1 && NULL) IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2 NOT IN (+ NULL, 1 && NULL) IS NULL) <=> ((((+ 2) NOT IN ((+ NULL), (1 AND NULL))) IS NULL))
< result
< column "(+ 2 NOT IN (+ NULL, 1 && NULL) IS NULL) <=> ((((+ 2) NOT IN ((+ NULL), (1 AND NULL))) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 <> NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 19 near "NOT 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >> 0 BETWEEN 0 NOT LIKE NULL AND NULL LIKE 3 - 1 <=> 2 NOT IN (NULL | 0, 0 DIV 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "NOT LIKE NULL AND NULL LIKE 3 - 1 <=> 2 NOT IN (NULL | 0, 0 DIV 1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 != 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 != 0) <=> ((1 <> 0))
< result
< column "(1 != 0) <=> ((1 <> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := NULL IS NOT NULL) <=> ((@v := (NULL IS NOT NULL)))
< result
< column "(@v := NULL IS NOT NULL) <=> ((@v := (NULL IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 3 = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 | 3 = 3) <=> (((2 | 3) = 3))
< result
< column "(2 | 3 = 3) <=> (((2 | 3) = 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT IN (NULL, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT IN (NULL, 1)) <=> ((2 NOT IN (NULL, 1)))
< result
< column "(2 NOT IN (NULL, 1)) <=> ((2 NOT IN (NULL, 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 3 AND 0) <=> ((2 NOT BETWEEN 3 AND 0))
< result
< column "(2 NOT BETWEEN 3 AND 0) <=> ((2 NOT BETWEEN 3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL) <=> (NULL)
< result
< column "(NULL) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 >> 1 <=> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 >> 1 <=> NULL) <=> ((NOT ((1 >> 1) <=> NULL)))
< result
< column "(NOT 1 >> 1 <=> NULL) <=> ((NOT ((1 >> 1) <=> NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2) <=> ((NOT 2))
< result
< column "(! 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0) <=> ((~ 0))
< result
< column "(~ 0) <=> ((~ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := 0 XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := 0 XOR 0) <=> ((@v := (0 XOR 0)))
< result
< column "(@v := 0 XOR 0) <=> ((@v := (0 XOR 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 BETWEEN 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 BETWEEN 0 AND 2) <=> ((0 BETWEEN 0 AND 2))
< result
< column "(0 BETWEEN 0 AND 2) <=> ((0 BETWEEN 0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 <=> 3 != @v := 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1 <=> 3 != @v := 0) <=> ((((NOT 1) <=> 3) <> (@v := 0)))
< result
< column "(! 1 <=> 3 != @v := 0) <=> ((((NOT 1) <=> 3) <> (@v := 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NULL <= ~ NULL IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NULL <= ~ NULL IS TRUE) <=> ((((1 IS NULL) <= (~ NULL)) IS TRUE))
< result
< column "(1 IS NULL <= ~ NULL IS TRUE) <=> ((((1 IS NULL) <= (~ NULL)) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 1 AND 0 BETWEEN 2 <= 1 AND 1 % 0 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "<= 1 AND 1 % 0 IS NOT FALSE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 MOD 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 MOD 1) <=> ((1 % 1))
< result
< column "(1 MOD 1) <=> ((1 % 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 <=> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 <=> 0) <=> ((NOT (2 <=> 0)))
< result
< column "(NOT 2 <=> 0) <=> ((NOT (2 <=> 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (3, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (3, 2)) <=> ((0 NOT IN (3, 2)))
< result
< column "(0 NOT IN (3, 2)) <=> ((0 NOT IN (3, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 ^ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 ^ NULL) <=> ((2 ^ NULL))
< result
< column "(2 ^ NULL) <=> ((2 ^ NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 AND 3) <=> ((2 AND 3))
< result
< column "(2 AND 3) <=> ((2 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL NOT BETWEEN NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL NOT BETWEEN NULL AND 0) <=> (((- NULL) NOT BETWEEN NULL AND 0))
< result
< column "(- NULL NOT BETWEEN NULL AND 0) <=> (((- NULL) NOT BETWEEN NULL AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + (0) REGEXP ! NULL BETWEEN NULL AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 34 near "BETWEEN NULL AND NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 && 2) <=> ((2 AND 2))
< result
< column "(2 && 2) <=> ((2 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1 MOD 1 NOT IN (NULL, 2) + 2 | 3 > NULL % 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "+ 2 | 3 > NULL % 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 = 3) <=> ((3 = 3))
< result
< column "(3 = 3) <=> ((3 = 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | NULL) <=> ((0 | NULL))
< result
< column "(0 | NULL) <=> ((0 | NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (3, 0) BETWEEN 3 >= 0 AND (3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "BETWEEN 3 >= 0 AND (3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 < 2 IS TRUE >> ! NULL DIV 1 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near ">> ! NULL DIV 1 <> 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL BETWEEN 2 AND 3 IN (- NOT 2, 2 NOT LIKE 2 LIKE ~ 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 39 near "NOT 2, 2 NOT LIKE 2 LIKE ~ 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 0 AND 3) <=> ((NULL NOT BETWEEN 0 AND 3))
< result
< column "(NULL NOT BETWEEN 0 AND 3) <=> ((NULL NOT BETWEEN 0 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 0 IS UNKNOWN * 2 && 1 NOT BETWEEN @v := 0 AND 2 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "* 2 && 1 NOT BETWEEN @v := 0 AND 2 IS UNKNOWN" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 - 0) <=> ((2 - 0))
< result
< column "(2 - 0) <=> ((2 - 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 OR 0 MOD 2 & 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 OR 0 MOD 2 & 1) <=> ((2 OR ((0 % 2) & 1)))
< result
< column "(2 OR 0 MOD 2 & 1) <=> ((2 OR ((0 % 2) & 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 MOD 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 MOD 2) <=> ((3 % 2))
< result
< column "(3 MOD 2) <=> ((3 % 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> NULL < ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> NULL < ! 0) <=> (((0 >> NULL) < (NOT 0)))
< result
< column "(0 >> NULL < ! 0) <=> (((0 >> NULL) < (NOT 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 REGEXP 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 REGEXP 0) <=> (((- 1) REGEXP 0))
< result
< column "(- 1 REGEXP 0) <=> (((- 1) REGEXP 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NOT 3) IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NOT 3) IS NOT NULL) <=> (((NOT 3) IS NOT NULL))
< result
< column "((NOT 3) IS NOT NULL) <=> (((NOT 3) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (3, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (3, 0)) <=> ((2 IN (3, 0)))
< result
< column "(2 IN (3, 0)) <=> ((2 IN (3, 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 0 AND 3) <=> ((3 BETWEEN 0 AND 3))
< result
< column "(3 BETWEEN 0 AND 3) <=> ((3 BETWEEN 0 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := ~ 1 = 0 > 3 <> NULL <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := ~ 1 = 0 > 3 <> NULL <> 3) <=> ((@v := (((((~ 1) = 0) > 3) <> NULL) <> 3)))
< result
< column "(@v := ~ 1 = 0 > 3 <> NULL <> 3) <=> ((@v := (((((~ 1) = 0) > 3) <> NULL) <> 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL << 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL << 1) <=> ((NULL << 1))
< result
< column "(NULL << 1) <=> ((NULL << 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE IN (- 2, 1 = 2) <> ! 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "IN (- 2, 1 = 2) <> ! 0 AND 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS TRUE IN (@v := 2, (2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 19 near "IN (@v := 2, (2))" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL) & NULL BETWEEN 0 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL) & NULL BETWEEN 0 AND 3) <=> (((NULL & NULL) BETWEEN 0 AND 3))
< result
< column "((NULL) & NULL BETWEEN 0 AND 3) <=> (((NULL & NULL) BETWEEN 0 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 && 0 IN (1 > 2 IN (3 || 2, NULL MOD NULL), 1 / 0 << 1 << 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3 && 0 IN (1 > 2 IN (3 || 2, NULL MOD NULL), 1 / 0 << 1 << 2)) <=> (((~ 3) AND (0 IN ((1 > (2 IN ((3 OR 2), (NULL % NULL)))), (((1 / 0) << 1) << 2)))))
< result
< column "(~ 3 && 0 IN (1 > 2 IN (3 || 2, NULL MOD NULL), 1 / 0 << 1 << 2)) <=> (((~ 3) AND (0 IN ((1 > (2 IN ((3 OR 2), (NULL % NULL)))), (((1 / 0) << 1) << 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 LIKE 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 LIKE 3) <=> ((1 LIKE 3))
< result
< column "(1 LIKE 3) <=> ((1 LIKE 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 / 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 / 1) <=> ((3 / 1))
< result
< column "(3 / 1) <=> ((3 / 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NULL) <=> ((NULL IS NULL))
< result
< column "(NULL IS NULL) <=> ((NULL IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL BETWEEN 3 AND 3 | 3 OR NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL BETWEEN 3 AND 3 | 3 OR NULL)) <=> (((NULL BETWEEN 3 AND (3 | 3)) OR NULL))
< result
< column "((NULL BETWEEN 3 AND 3 | 3 OR NULL)) <=> (((NULL BETWEEN 3 AND (3 | 3)) OR NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> NULL MOD 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> NULL MOD 2) <=> ((0 >> (NULL % 2)))
< result
< column "(0 >> NULL MOD 2) <=> ((0 >> (NULL % 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0)) <=> (0)
< result
< column "((0)) <=> (0)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2 ^ 3) DIV 0 XOR NULL | 1 NOT BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((2 ^ 3) DIV 0 XOR NULL | 1 NOT BETWEEN NULL AND 1) <=> ((((2 ^ 3) DIV 0) XOR ((NULL | 1) NOT BETWEEN NULL AND 1)))
< result
< column "((2 ^ 3) DIV 0 XOR NULL | 1 NOT BETWEEN NULL AND 1) <=> ((((2 ^ 3) DIV 0) XOR ((NULL | 1) NOT BETWEEN NULL AND 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <=> 0 NOT BETWEEN ! 1 AND 0 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <=> 0 NOT BETWEEN ! 1 AND 0 IS UNKNOWN) <=> (((3 <=> (0 NOT BETWEEN (NOT 1) AND 0)) IS UNKNOWN))
< result
< column "(3 <=> 0 NOT BETWEEN ! 1 AND 0 IS UNKNOWN) <=> (((3 <=> (0 NOT BETWEEN (NOT 1) AND 0)) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <= 3 && 0 IN (2, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <= 3 && 0 IN (2, 1)) <=> (((2 <= 3) AND (0 IN (2, 1))))
< result
< column "(2 <= 3 && 0 IN (2, 1)) <=> (((2 <= 3) AND (0 IN (2, 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 < 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 < 3) <=> ((3 < 3))
< result
< column "(3 < 3) <=> ((3 < 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 - NULL) <=> ((1 - NULL))
< result
< column "(1 - NULL) <=> ((1 - NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 BETWEEN 0 - 3 AND 2 IN (1, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 BETWEEN 0 - 3 AND 2 IN (1, 1)) <=> ((NOT (3 BETWEEN (0 - 3) AND (2 IN (1, 1)))))
< result
< column "(NOT 3 BETWEEN 0 - 3 AND 2 IN (1, 1)) <=> ((NOT (3 BETWEEN (0 - 3) AND (2 IN (1, 1)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 2 AND 0 IS NULL + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 41 near "+ 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 / 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 / 3) <=> ((3 / 3))
< result
< column "(3 / 3) <=> ((3 / 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 IS NOT FALSE) <=> (((- 1) IS NOT FALSE))
< result
< column "(- 1 IS NOT FALSE) <=> (((- 1) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - ~ NULL ^ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- ~ NULL ^ 1) <=> (((- (~ NULL)) ^ 1))
< result
< column "(- ~ NULL ^ 1) <=> (((- (~ NULL)) ^ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL / 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL / 3) <=> ((NULL / 3))
< result
< column "(NULL / 3) <=> ((NULL / 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT REGEXP 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT REGEXP 2) <=> ((NULL NOT REGEXP 2))
< result
< column "(NULL NOT REGEXP 2) <=> ((NULL NOT REGEXP 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 1 BETWEEN 3 NOT LIKE 0 AND 0 NOT LIKE 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "NOT LIKE 0 AND 0 NOT LIKE 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV 1 = 0 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV 1 = 0 <> NULL) <=> ((((1 DIV 1) = 0) <> NULL))
< result
< column "(1 DIV 1 = 0 <> NULL) <=> ((((1 DIV 1) = 0) <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS TRUE) <=> ((NULL IS TRUE))
< result
< column "(NULL IS TRUE) <=> ((NULL IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 || 0 REGEXP 2 BETWEEN 2 DIV 0 NOT BETWEEN 1 XOR 1 AND 2 NOT REGEXP 0 AND 0 MOD 0 < NULL DIV NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 34 near "BETWEEN 2 DIV 0 NOT BETWEEN 1 XOR 1 AND 2 NOT REGEXP 0 AND 0 MOD 0 < NULL DIV NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1) <=> ((NOT 1))
< result
< column "(! 1) <=> ((NOT 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 DIV 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 DIV 2) <=> ((3 DIV 2))
< result
< column "(3 DIV 2) <=> ((3 DIV 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 & - NULL - 2 MOD 0 && NULL BETWEEN 3 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 & - NULL - 2 MOD 0 && NULL BETWEEN 3 AND 3) <=> (((NOT (1 & ((- NULL) - (2 % 0)))) AND (NULL BETWEEN 3 AND 3)))
< result
< column "(NOT 1 & - NULL - 2 MOD 0 && NULL BETWEEN 3 AND 3) <=> (((NOT (1 & ((- NULL) - (2 % 0)))) AND (NULL BETWEEN 3 AND 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 NOT IN (1, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 NOT IN (1, NULL)) <=> (((NOT 3) NOT IN (1, NULL)))
< result
< column "(! 3 NOT IN (1, NULL)) <=> (((NOT 3) NOT IN (1, NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2) <=> ((NOT 2))
< result
< column "(! 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2) <=> ((~ 2))
< result
< column "(~ 2) <=> ((~ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0 > 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0 > 3)) <=> ((0 > 3))
< result
< column "((0 > 3)) <=> ((0 > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 XOR NULL MOD 0 * 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 XOR NULL MOD 0 * 2) <=> ((2 XOR ((NULL % 0) * 2)))
< result
< column "(2 XOR NULL MOD 0 * 2) <=> ((2 XOR ((NULL % 0) * 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 MOD 0 XOR 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 2 MOD 0 XOR 1 IS UNKNOWN) <=> ((((- 2) % 0) XOR (1 IS UNKNOWN)))
< result
< column "(- 2 MOD 0 XOR 1 IS UNKNOWN) <=> ((((- 2) % 0) XOR (1 IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 & NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 & NULL) <=> ((1 & NULL))
< result
< column "(1 & NULL) <=> ((1 & NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2) <=> ((NOT 2))
< result
< column "(! 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT LIKE NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT LIKE NULL) <=> ((0 NOT LIKE NULL))
< result
< column "(0 NOT LIKE NULL) <=> ((0 NOT LIKE NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 1 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 1 AND 1) <=> ((3 BETWEEN 1 AND 1))
< result
< column "(3 BETWEEN 1 AND 1) <=> ((3 BETWEEN 1 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 3 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN 3 AND NULL) <=> ((1 BETWEEN 3 AND NULL))
< result
< column "(1 BETWEEN 3 AND NULL) <=> ((1 BETWEEN 3 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 3 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN 3 AND NULL) <=> ((1 BETWEEN 3 AND NULL))
< result
< column "(1 BETWEEN 3 AND NULL) <=> ((1 BETWEEN 3 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | 2) <=> ((0 | 2))
< result
< column "(0 | 2) <=> ((0 | 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND 2 BETWEEN 2 << 1 AND 3 * 2 < - 0 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND 2 BETWEEN 2 << 1 AND 3 * 2 < - 0 IS NULL) <=> ((NULL AND (((2 BETWEEN (2 << 1) AND (3 * 2)) < (- 0)) IS NULL)))
< result
< column "(NULL AND 2 BETWEEN 2 << 1 AND 3 * 2 < - 0 IS NULL) <=> ((NULL AND (((2 BETWEEN (2 << 1) AND (3 * 2)) < (- 0)) IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 IS NOT NULL) <=> ((NOT (2 IS NOT NULL)))
< result
< column "(NOT 2 IS NOT NULL) <=> ((NOT (2 IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := NULL > 2 ^ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := NULL > 2 ^ 1) <=> ((@v := (NULL > (2 ^ 1))))
< result
< column "(@v := NULL > 2 ^ 1) <=> ((@v := (NULL > (2 ^ 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV NULL - + 3 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV NULL - + 3 IS NOT FALSE) <=> ((((2 DIV NULL) - (+ 3)) IS NOT FALSE))
< result
< column "(2 DIV NULL - + 3 IS NOT FALSE) <=> ((((2 DIV NULL) - (+ 3)) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 AND 0) <=> ((2 AND 0))
< result
< column "(2 AND 0) <=> ((2 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 - 2 >> 3 | NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 - 2 >> 3 | NULL) <=> (((((NOT 3) - 2) >> 3) | NULL))
< result
< column "(! 3 - 2 >> 3 | NULL) <=> (((((NOT 3) - 2) >> 3) | NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT NULL) <=> ((1 IS NOT NULL))
< result
< column "(1 IS NOT NULL) <=> ((1 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (0, 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 IN (0, 3)) <=> ((0 IN (0, 3)))
< result
< column "(0 IN (0, 3)) <=> ((0 IN (0, 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 ^ 2) <=> ((1 ^ 2))
< result
< column "(1 ^ 2) <=> ((1 ^ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1) <=> ((NOT 1))
< result
< column "(! 1) <=> ((NOT 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2) <=> ((NOT 2))
< result
< column "(! 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = NULL) <=> ((NULL = NULL))
< result
< column "(NULL = NULL) <=> ((NULL = NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= NULL > 2 DIV 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= NULL > 2 DIV 1) <=> (((3 >= NULL) > (2 DIV 1)))
< result
< column "(3 >= NULL > 2 DIV 1) <=> (((3 >= NULL) > (2 DIV 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT BETWEEN 0 AND NULL) <=> ((0 NOT BETWEEN 0 AND NULL))
< result
< column "(0 NOT BETWEEN 0 AND NULL) <=> ((0 NOT BETWEEN 0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - 1 NOT BETWEEN (NULL) AND (NULL) LIKE 3 <= 3 NOT BETWEEN 3 NOT REGEXP 1 AND 2 < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 72 near "NOT REGEXP 1 AND 2 < 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <= 1 <= 1 NOT IN (0, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <= 1 <= 1 NOT IN (0, 0)) <=> (((1 <= 1) <= (1 NOT IN (0, 0))))
< result
< column "(1 <= 1 <= 1 NOT IN (0, 0)) <=> (((1 <= 1) <= (1 NOT IN (0, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 2 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN 2 AND 1) <=> ((NULL BETWEEN 2 AND 1))
< result
< column "(NULL BETWEEN 2 AND 1) <=> ((NULL BETWEEN 2 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IN (2, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IN (2, 1)) <=> ((3 IN (2, 1)))
< result
< column "(3 IN (2, 1)) <=> ((3 IN (2, 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 & 0 % 0 NOT IN (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 & 0 % 0 NOT IN (1, 2)) <=> (((0 & (0 % 0)) NOT IN (1, 2)))
< result
< column "(0 & 0 % 0 NOT IN (1, 2)) <=> (((0 & (0 % 0)) NOT IN (1, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (3 IS UNKNOWN)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((3 IS UNKNOWN)) <=> ((3 IS UNKNOWN))
< result
< column "((3 IS UNKNOWN)) <=> ((3 IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3) <=> (3)
< result
< column "(3) <=> (3)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT LIKE 1 < 0 NOT LIKE 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT LIKE 1 < 0 NOT LIKE 0) <=> (((1 NOT LIKE 1) < (0 NOT LIKE 0)))
< result
< column "(1 NOT LIKE 1 < 0 NOT LIKE 0) <=> (((1 NOT LIKE 1) < (0 NOT LIKE 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL XOR 0) <=> ((NULL XOR 0))
< result
< column "(NULL XOR 0) <=> ((NULL XOR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 2 AND 0 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 2 AND 0 IS NOT NULL) <=> (((2 NOT BETWEEN 2 AND 0) IS NOT NULL))
< result
< column "(2 NOT BETWEEN 2 AND 0 IS NOT NULL) <=> (((2 NOT BETWEEN 2 AND 0) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - NULL % 0 NOT IN (1, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 - NULL % 0 NOT IN (1, 1)) <=> (((3 - (NULL % 0)) NOT IN (1, 1)))
< result
< column "(3 - NULL % 0 NOT IN (1, 1)) <=> (((3 - (NULL % 0)) NOT IN (1, 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 OR 0 IS UNKNOWN <= 2 NOT LIKE 0 NOT IN (0 NOT IN (1, 1), 0 BETWEEN 0 AND NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "<= 2 NOT LIKE 0 NOT IN (0 NOT IN (1, 1), 0 BETWEEN 0 AND NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 = 2) <=> ((3 = 2))
< result
< column "(3 = 2) <=> ((3 = 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL >= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL >= 3) <=> ((NULL >= 3))
< result
< column "(NULL >= 3) <=> ((NULL >= 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <=> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <=> NULL) <=> ((3 <=> NULL))
< result
< column "(3 <=> NULL) <=> ((3 <=> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 & NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 & NULL) <=> ((NOT (3 & NULL)))
< result
< column "(NOT 3 & NULL) <=> ((NOT (3 & NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 > 3) <=> ((0 > 3))
< result
< column "(0 > 3) <=> ((0 > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (3, NULL) IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (3, NULL) IS NULL) <=> (((2 IN (3, NULL)) IS NULL))
< result
< column "(2 IN (3, NULL) IS NULL) <=> (((2 IN (3, NULL)) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (3, NULL) IS NOT NULL && 2 DIV 0 ^ ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (3, NULL) IS NOT NULL && 2 DIV 0 ^ ~ NULL) <=> ((((0 NOT IN (3, NULL)) IS NOT NULL) AND (2 DIV (0 ^ (~ NULL)))))
< result
< column "(0 NOT IN (3, NULL) IS NOT NULL && 2 DIV 0 ^ ~ NULL) <=> ((((0 NOT IN (3, NULL)) IS NOT NULL) AND (2 DIV (0 ^ (~ NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 % 3 / 1 NOT IN (0, 0) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 % 3 / 1 NOT IN (0, 0) IS NOT FALSE) <=> (((((1 % 3) / 1) NOT IN (0, 0)) IS NOT FALSE))
< result
< column "(1 % 3 / 1 NOT IN (0, 0) IS NOT FALSE) <=> (((((1 % 3) / 1) NOT IN (0, 0)) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN NULL AND 0) <=> ((NULL BETWEEN NULL AND 0))
< result
< column "(NULL BETWEEN NULL AND 0) <=> ((NULL BETWEEN NULL AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS TRUE / 1 MOD NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 18 near "/ 1 MOD NULL IS NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 3 AND 3 MOD NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "NOT 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV NULL) <=> ((1 DIV NULL))
< result
< column "(1 DIV NULL) <=> ((1 DIV NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / NULL IS NULL | + 3 / 3 REGEXP 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "| + 3 / 3 REGEXP 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 2 AND 0 OR NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN 2 AND 0 OR NOT 2) <=> (((1 BETWEEN 2 AND 0) OR (NOT 2)))
< result
< column "(1 BETWEEN 2 AND 0 OR NOT 2) <=> (((1 BETWEEN 2 AND 0) OR (NOT 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (1, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (1, 0)) <=> ((2 IN (1, 0)))
< result
< column "(2 IN (1, 0)) <=> ((2 IN (1, 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 XOR 1 + NULL + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 XOR 1 + NULL + 0) <=> ((3 XOR ((1 + NULL) + 0)))
< result
< column "(3 XOR 1 + NULL + 0) <=> ((3 XOR ((1 + NULL) + 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1) <=> ((+ 1))
< result
< column "(+ 1) <=> ((+ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 2 AND 0 MOD 2 NOT IN (2, 1) > 1 NOT BETWEEN NULL AND 0 NOT BETWEEN - 2 AND 0 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 2 AND 0 MOD 2 NOT IN (2, 1) > 1 NOT BETWEEN NULL AND 0 NOT BETWEEN - 2 AND 0 <> 3) <=> ((((3 BETWEEN 2 AND ((0 % 2) NOT IN (2, 1))) > (1 NOT BETWEEN NULL AND (0 NOT BETWEEN (- 2) AND 0))) <> 3))
< result
< column "(3 BETWEEN 2 AND 0 MOD 2 NOT IN (2, 1) > 1 NOT BETWEEN NULL AND 0 NOT BETWEEN - 2 AND 0 <> 3) <=> ((((3 BETWEEN 2 AND ((0 % 2) NOT IN (2, 1))) > (1 NOT BETWEEN NULL AND (0 NOT BETWEEN (- 2) AND 0))) <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1 MOD 1) > 1 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 0 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1 MOD 1) > 1 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 0 AND 1) <=> ((((1 % 1) > (1 NOT BETWEEN 0 AND 1)) AND (1 NOT BETWEEN 0 AND 1)))
< result
< column "((1 MOD 1) > 1 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 0 AND 1) <=> ((((1 % 1) > (1 NOT BETWEEN 0 AND 1)) AND (1 NOT BETWEEN 0 AND 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 AND 3 = 2 <=> 0 XOR 2 >> 2 BETWEEN 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 AND 3 = 2 <=> 0 XOR 2 >> 2 BETWEEN 3 AND 0) <=> ((((NOT 3) AND ((3 = 2) <=> 0)) XOR ((2 >> 2) BETWEEN 3 AND 0)))
< result
< column "(NOT 3 AND 3 = 2 <=> 0 XOR 2 >> 2 BETWEEN 3 AND 0) <=> ((((NOT 3) AND ((3 = 2) <=> 0)) XOR ((2 >> 2) BETWEEN 3 AND 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 & 3 NOT BETWEEN 1 NOT IN (0, 3) AND 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "NOT IN (0, 3) AND 0 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (1, 2) IS TRUE & NULL != 2 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "& NULL != 2 IS UNKNOWN" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 << 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 << 2) <=> ((1 << 2))
< result
< column "(1 << 2) <=> ((1 << 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >= 0 <> 3 <=> NULL && 0 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >= 0 <> 3 <=> NULL && 0 IS NOT NULL) <=> (((((1 >= 0) <> 3) <=> NULL) AND (0 IS NOT NULL)))
< result
< column "(1 >= 0 <> 3 <=> NULL && 0 IS NOT NULL) <=> (((((1 >= 0) <> 3) <=> NULL) AND (0 IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2) <=> ((~ 2))
< result
< column "(~ 2) <=> ((~ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 != 3 DIV 3 >= 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 != 3 DIV 3 >= 0) <=> (((0 <> (3 DIV 3)) >= 0))
< result
< column "(0 != 3 DIV 3 >= 0) <=> (((0 <> (3 DIV 3)) >= 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >= 2) <=> ((1 >= 2))
< result
< column "(1 >= 2) <=> ((1 >= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <=> 1) <=> ((NULL <=> 1))
< result
< column "(NULL <=> 1) <=> ((NULL <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IN (0, 3) IS TRUE IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "IS NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NULL) <=> ((NULL IS NULL))
< result
< column "(NULL IS NULL) <=> ((NULL IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <= 2) <=> ((3 <= 2))
< result
< column "(3 <= 2) <=> ((3 <= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IN (0, NULL) + 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "+ 0 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >> 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >> 2) <=> ((1 >> 2))
< result
< column "(1 >> 2) <=> ((1 >> 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2) <=> (2)
< result
< column "(2) <=> (2)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 3) <=> ((~ 3))
< result
< column "(~ 3) <=> ((~ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NOT NULL) <=> ((NULL IS NOT NULL))
< result
< column "(NULL IS NOT NULL) <=> ((NULL IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 DIV NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 DIV NULL) <=> ((3 DIV NULL))
< result
< column "(3 DIV NULL) <=> ((3 DIV NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 || 1 XOR 2 IS NOT FALSE IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "IS UNKNOWN" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 3 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT BETWEEN 3 AND 1) <=> ((1 NOT BETWEEN 3 AND 1))
< result
< column "(1 NOT BETWEEN 3 AND 1) <=> ((1 NOT BETWEEN 3 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT NULL) <=> ((3 IS NOT NULL))
< result
< column "(3 IS NOT NULL) <=> ((3 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IS NOT NULL) <=> ((3 IS NOT NULL))
< result
< column "(3 IS NOT NULL) <=> ((3 IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 * NULL BETWEEN 1 = 3 AND 2 = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "= 3 AND 2 = 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL < 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL < 0) <=> ((NULL < 0))
< result
< column "(NULL < 0) <=> ((NULL < 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + 0 IN (NULL | NULL, NOT 3) >= ! 0 ^ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + 0 IN (NULL | NULL, NOT 3) >= ! 0 ^ NULL) <=> ((((3 + 0) IN ((NULL | NULL), (NOT 3))) >= ((NOT 0) ^ NULL)))
< result
< column "(3 + 0 IN (NULL | NULL, NOT 3) >= ! 0 ^ NULL) <=> ((((3 + 0) IN ((NULL | NULL), (NOT 3))) >= ((NOT 0) ^ NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 ^ NULL <> + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 ^ NULL <> + 3) <=> (((3 ^ NULL) <> (+ 3)))
< result
< column "(3 ^ NULL <> + 3) <=> (((3 ^ NULL) <> (+ 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT LIKE NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT LIKE NULL) <=> ((3 NOT LIKE NULL))
< result
< column "(3 NOT LIKE NULL) <=> ((3 NOT LIKE NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 NOT BETWEEN - 1 AND 2 XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 NOT BETWEEN - 1 AND 2 XOR 0) <=> (((NOT (1 NOT BETWEEN (- 1) AND 2)) XOR 0))
< result
< column "(NOT 1 NOT BETWEEN - 1 AND 2 XOR 0) <=> (((NOT (1 NOT BETWEEN (- 1) AND 2)) XOR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1 + 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1 + 0)) <=> ((1 + 0))
< result
< column "((1 + 0)) <=> ((1 + 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS TRUE ^ 1 / 1 && NULL + 2 BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "^ 1 / 1 && NULL + 2 BETWEEN 2 AND 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (0) <= 0 != NULL XOR 3 NOT BETWEEN NULL AND 0 LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0) <= 0 != NULL XOR 3 NOT BETWEEN NULL AND 0 LIKE 2) <=> ((((0 <= 0) <> NULL) XOR (3 NOT BETWEEN NULL AND (0 LIKE 2))))
< result
< column "((0) <= 0 != NULL XOR 3 NOT BETWEEN NULL AND 0 LIKE 2) <=> ((((0 <= 0) <> NULL) XOR (3 NOT BETWEEN NULL AND (0 LIKE 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0 < 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0 < 1)) <=> ((0 < 1))
< result
< column "((0 < 1)) <=> ((0 < 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL) <=> ((NOT NULL))
< result
< column "(NOT NULL) <=> ((NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> NULL << 3 NOT LIKE 2 & NOT 3 <=> 1 OR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 34 near "& NOT 3 <=> 1 OR 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= NULL <= 2 / NULL % 0 + 3 <= 0 MOD 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= NULL <= 2 / NULL % 0 + 3 <= 0 MOD 1) <=> ((((2 >= NULL) <= (((2 / NULL) % 0) + 3)) <= (0 % 1)))
< result
< column "(2 >= NULL <= 2 / NULL % 0 + 3 <= 0 MOD 1) <=> ((((2 >= NULL) <= (((2 / NULL) % 0) + 3)) <= (0 % 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL IS NULL IS NOT NULL) <=> (((NULL IS NULL) IS NOT NULL))
< result
< column "(NULL IS NULL IS NOT NULL) <=> (((NULL IS NULL) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 % 2 BETWEEN 1 AND 1 DIV 3 BETWEEN 0 AND 1 XOR 3 ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0 % 2 BETWEEN 1 AND 1 DIV 3 BETWEEN 0 AND 1 XOR 3 ^ 2) <=> (((((+ 0) % 2) BETWEEN 1 AND ((1 DIV 3) BETWEEN 0 AND 1)) XOR (3 ^ 2)))
< result
< column "(+ 0 % 2 BETWEEN 1 AND 1 DIV 3 BETWEEN 0 AND 1 XOR 3 ^ 2) <=> (((((+ 0) % 2) BETWEEN 1 AND ((1 DIV 3) BETWEEN 0 AND 1)) XOR (3 ^ 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 BETWEEN 2 AND 0 <= - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 BETWEEN 2 AND 0 <= - 0) <=> (((0 BETWEEN 2 AND 0) <= (- 0)))
< result
< column "(0 BETWEEN 2 AND 0 <= - 0) <=> (((0 BETWEEN 2 AND 0) <= (- 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT REGEXP NULL >= NULL != 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT REGEXP NULL >= NULL != 1) <=> ((((NULL NOT REGEXP NULL) >= NULL) <> 1))
< result
< column "(NULL NOT REGEXP NULL >= NULL != 1) <=> ((((NULL NOT REGEXP NULL) >= NULL) <> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL OR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL OR 1) <=> ((NULL OR 1))
< result
< column "(NULL OR 1) <=> ((NULL OR 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 || 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 || 0) <=> ((1 OR 0))
< result
< column "(1 || 0) <=> ((1 OR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 3 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 3 AND 2) <=> ((2 NOT BETWEEN 3 AND 2))
< result
< column "(2 NOT BETWEEN 3 AND 2) <=> ((2 NOT BETWEEN 3 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 NOT IN (1 NOT BETWEEN NULL AND 3, 2 >= 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 NOT IN (1 NOT BETWEEN NULL AND 3, 2 >= 2)) <=> ((NOT (3 NOT IN ((1 NOT BETWEEN NULL AND 3), (2 >= 2)))))
< result
< column "(NOT 3 NOT IN (1 NOT BETWEEN NULL AND 3, 2 >= 2)) <=> ((NOT (3 NOT IN ((1 NOT BETWEEN NULL AND 3), (2 >= 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL OR 0 >> 3 NOT IN (2, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL OR 0 >> 3 NOT IN (2, 2)) <=> (((+ NULL) OR ((0 >> 3) NOT IN (2, 2))))
< result
< column "(+ NULL OR 0 >> 3 NOT IN (2, 2)) <=> (((+ NULL) OR ((0 >> 3) NOT IN (2, 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 1 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT BETWEEN 1 AND 3) <=> ((1 NOT BETWEEN 1 AND 3))
< result
< column "(1 NOT BETWEEN 1 AND 3) <=> ((1 NOT BETWEEN 1 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS UNKNOWN NOT LIKE 2 > 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "NOT LIKE 2 > 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 0) <=> ((NULL = 0))
< result
< column "(NULL = 0) <=> ((NULL = 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 = 1) <=> ((0 = 1))
< result
< column "(0 = 1) <=> ((0 = 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 OR 3 IS NOT FALSE < - 0 * NULL OR 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "< - 0 * NULL OR 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0) <=> (0)
< result
< column "(0) <=> (0)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 OR NULL >> 0 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 OR NULL >> 0 % 0) <=> ((2 OR (NULL >> (0 % 0))))
< result
< column "(2 OR NULL >> 0 % 0) <=> ((2 OR (NULL >> (0 % 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 3 AND 0 | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 3 AND 0 | 2) <=> ((2 BETWEEN 3 AND (0 | 2)))
< result
< column "(2 BETWEEN 3 AND 0 | 2) <=> ((2 BETWEEN 3 AND (0 | 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 IS NOT FALSE) <=> ((NOT (1 IS NOT FALSE)))
< result
< column "(NOT 1 IS NOT FALSE) <=> ((NOT (1 IS NOT FALSE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1 NOT BETWEEN NULL AND 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1 NOT BETWEEN NULL AND 2)) <=> ((1 NOT BETWEEN NULL AND 2))
< result
< column "((1 NOT BETWEEN NULL AND 2)) <=> ((1 NOT BETWEEN NULL AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NOT NULL != 0 > 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS NOT NULL != 0 > 1) <=> ((((2 IS NOT NULL) <> 0) > 1))
< result
< column "(2 IS NOT NULL != 0 > 1) <=> ((((2 IS NOT NULL) <> 0) > 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 3) <=> ((NULL = 3))
< result
< column "(NULL = 3) <=> ((NULL = 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0) <=> ((NOT 0))
< result
< column "(! 0) <=> ((NOT 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 3 AND 0 IN (1 NOT BETWEEN 2 AND NULL, NULL IS NOT NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 3 AND 0 IN (1 NOT BETWEEN 2 AND NULL, NULL IS NOT NULL)) <=> ((2 BETWEEN 3 AND (0 IN ((1 NOT BETWEEN 2 AND NULL), (NULL IS NOT NULL)))))
< result
< column "(2 BETWEEN 3 AND 0 IN (1 NOT BETWEEN 2 AND NULL, NULL IS NOT NULL)) <=> ((2 BETWEEN 3 AND (0 IN ((1 NOT BETWEEN 2 AND NULL), (NULL IS NOT NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 / 1 << 3 << ~ NULL <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0 / 1 << 3 << ~ NULL <> 3) <=> ((((((+ 0) / 1) << 3) << (~ NULL)) <> 3))
< result
< column "(+ 0 / 1 << 3 << ~ NULL <> 3) <=> ((((((+ 0) / 1) << 3) << (~ NULL)) <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 IS TRUE) <=> (((~ 2) IS TRUE))
< result
< column "(~ 2 IS TRUE) <=> (((~ 2) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL << NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL << NULL) <=> ((NULL << NULL))
< result
< column "(NULL << NULL) <=> ((NULL << NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (0 = 2, (0)) < 1 + 2 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (0 = 2, (0)) < 1 + 2 IS NOT NULL) <=> ((((0 NOT IN ((0 = 2), 0)) < (1 + 2)) IS NOT NULL))
< result
< column "(0 NOT IN (0 = 2, (0)) < 1 + 2 IS NOT NULL) <=> ((((0 NOT IN ((0 = 2), 0)) < (1 + 2)) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 DIV 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 DIV 2) <=> ((0 DIV 2))
< result
< column "(0 DIV 2) <=> ((0 DIV 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 MOD NULL ^ NULL AND 0 << 3 LIKE 1 >= 2 < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 MOD NULL ^ NULL AND 0 << 3 LIKE 1 >= 2 < NULL) <=> (((3 % (NULL ^ NULL)) AND ((((0 << 3) LIKE 1) >= 2) < NULL)))
< result
< column "(3 MOD NULL ^ NULL AND 0 << 3 LIKE 1 >= 2 < NULL) <=> (((3 % (NULL ^ NULL)) AND ((((0 << 3) LIKE 1) >= 2) < NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (0, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT IN (0, 2)) <=> ((0 NOT IN (0, 2)))
< result
< column "(0 NOT IN (0, 2)) <=> ((0 NOT IN (0, 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 0 IN (1 <= NULL, NULL > 0) > 1 & 2 XOR + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 0 IN (1 <= NULL, NULL > 0) > 1 & 2 XOR + 0) <=> (((((- 0) IN ((1 <= NULL), (NULL > 0))) > (1 & 2)) XOR (+ 0)))
< result
< column "(- 0 IN (1 <= NULL, NULL > 0) > 1 & 2 XOR + 0) <=> (((((- 0) IN ((1 <= NULL), (NULL > 0))) > (1 & 2)) XOR (+ 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 = 3 <> NULL IS UNKNOWN >> 3 | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near ">> 3 | 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 2 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 2 AND 0) <=> (((3 NOT BETWEEN 0 AND 1) AND (1 NOT BETWEEN 2 AND 0)))
< result
< column "(3 NOT BETWEEN 0 AND 1 AND 1 NOT BETWEEN 2 AND 0) <=> (((3 NOT BETWEEN 0 AND 1) AND (1 NOT BETWEEN 2 AND 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IN (2, NULL) MOD NULL BETWEEN NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "MOD NULL BETWEEN NULL AND 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * 2 | NULL IS NOT NULL DIV 2 LIKE NULL << NULL || NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "DIV 2 LIKE NULL << NULL || NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 | 2) <=> ((3 | 2))
< result
< column "(3 | 2) <=> ((3 | 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 BETWEEN 1 AND 1 BETWEEN 2 NOT REGEXP 3 AND 3 != 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "NOT REGEXP 3 AND 3 != 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 NOT IN (0, 3) LIKE 3 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "LIKE 3 OR 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL XOR 1 = 3 IS NOT FALSE >> NOT 0 NOT LIKE 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near ">> NOT 0 NOT LIKE 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 % 1 >= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2 % 1 >= 3) <=> ((((~ 2) % 1) >= 3))
< result
< column "(~ 2 % 1 >= 3) <=> ((((~ 2) % 1) >= 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IN (2, 0) NOT BETWEEN 2 - 1 AND 1 NOT LIKE 3 || 2 OR 0 DIV ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "NOT BETWEEN 2 - 1 AND 1 NOT LIKE 3 || 2 OR 0 DIV ~ NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL > 3) <=> ((NULL > 3))
< result
< column "(NULL > 3) <=> ((NULL > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT ~ 2) <=> ((NOT (~ 2)))
< result
< column "(NOT ~ 2) <=> ((NOT (~ 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 2) <=> ((~ 2))
< result
< column "(~ 2) <=> ((~ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0 <=> 0) XOR 1 XOR 1 <= 0 > NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0 <=> 0) XOR 1 XOR 1 <= 0 > NULL) <=> ((((0 <=> 0) XOR 1) XOR ((1 <= 0) > NULL)))
< result
< column "((0 <=> 0) XOR 1 XOR 1 <= 0 > NULL) <=> ((((0 <=> 0) XOR 1) XOR ((1 <= 0) > NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL IN (3, 3) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL IN (3, 3) IS NOT FALSE) <=> ((NOT ((NULL IN (3, 3)) IS NOT FALSE)))
< result
< column "(NOT NULL IN (3, 3) IS NOT FALSE) <=> ((NOT ((NULL IN (3, 3)) IS NOT FALSE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT IN (1, 2) NOT IN (- 1, 0 < 3) < NULL IN (3 BETWEEN 3 AND 3, 1 NOT BETWEEN 2 AND 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "NOT IN (- 1, 0 < 3) < NULL IN (3 BETWEEN 3 AND 3, 1 NOT BETWEEN 2 AND 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IN (2, 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 IN (2, 3)) <=> ((3 IN (2, 3)))
< result
< column "(3 IN (2, 3)) <=> ((3 IN (2, 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 & 0 / 2 && 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 & 0 / 2 && 1) <=> (((3 & (0 / 2)) AND 1))
< result
< column "(3 & 0 / 2 && 1) <=> (((3 & (0 / 2)) AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 0 AND NULL) <=> ((2 BETWEEN 0 AND NULL))
< result
< column "(2 BETWEEN 0 AND NULL) <=> ((2 BETWEEN 0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 = 3 NOT IN (NULL, 0) MOD ! ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "MOD ! ! 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 & NULL <=> (0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 & NULL <=> (0)) <=> (((3 & NULL) <=> 0))
< result
< column "(3 & NULL <=> (0)) <=> (((3 & NULL) <=> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 ^ 3 XOR - 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 ^ 3 XOR - 3) <=> (((NOT (1 ^ 3)) XOR (- 3)))
< result
< column "(NOT 1 ^ 3 XOR - 3) <=> (((NOT (1 ^ 3)) XOR (- 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 IS NOT FALSE NOT LIKE NULL < 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "NOT LIKE NULL < 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 REGEXP 0 <> 2 >= 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 REGEXP 0 <> 2 >= 0) <=> ((((2 REGEXP 0) <> 2) >= 0))
< result
< column "(2 REGEXP 0 <> 2 >= 0) <=> ((((2 REGEXP 0) <> 2) >= 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 REGEXP 1 && 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 REGEXP 1 && 0) <=> (((1 REGEXP 1) AND 0))
< result
< column "(1 REGEXP 1 && 0) <=> (((1 REGEXP 1) AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 REGEXP 1 <> NULL <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 REGEXP 1 <> NULL <> 3) <=> ((((2 REGEXP 1) <> NULL) <> 3))
< result
< column "(2 REGEXP 1 <> NULL <> 3) <=> ((((2 REGEXP 1) <> NULL) <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (- 0) | 2 < 1 AND NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((- 0) | 2 < 1 AND NOT 1) <=> (((((- 0) | 2) < 1) AND (NOT 1)))
< result
< column "((- 0) | 2 < 1 AND NOT 1) <=> (((((- 0) | 2) < 1) AND (NOT 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 || 0 NOT BETWEEN 2 BETWEEN 1 AND NULL AND (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "BETWEEN 1 AND NULL AND (1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT REGEXP NULL NOT IN (@v := 2, 0 > 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "NOT IN (@v := 2, 0 > 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > NULL NOT LIKE @v := 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 33 near ":= 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN NULL AND NULL || 2 BETWEEN 2 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 BETWEEN NULL AND NULL || 2 BETWEEN 2 AND 0) <=> (((1 BETWEEN NULL AND NULL) OR (2 BETWEEN 2 AND 0)))
< result
< column "(1 BETWEEN NULL AND NULL || 2 BETWEEN 2 AND 0) <=> (((1 BETWEEN NULL AND NULL) OR (2 BETWEEN 2 AND 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT BETWEEN 0 AND 2) <=> ((1 NOT BETWEEN 0 AND 2))
< result
< column "(1 NOT BETWEEN 0 AND 2) <=> ((1 NOT BETWEEN 0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> 0 NOT IN (NULL != 2, 2 & 3) NOT LIKE 1 * 2 NOT LIKE 0 LIKE 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 43 near "NOT LIKE 1 * 2 NOT LIKE 0 LIKE 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL && 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL && 1) <=> ((NULL AND 1))
< result
< column "(NULL && 1) <=> ((NULL AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 / 2) <=> ((1 / 2))
< result
< column "(1 / 2) <=> ((1 / 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 > 0 > @v := 3 AND NULL > NULL REGEXP 3 NOT IN (1, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 51 near "NOT IN (1, NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 & 3 NOT BETWEEN 3 XOR NULL AND 1 & 3 <= 3 XOR NULL - 1 XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "XOR NULL AND 1 & 3 <= 3 XOR NULL - 1 XOR 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 - 1 << 0 & 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 - 1 << 0 & 3) <=> ((((2 - 1) << 0) & 3))
< result
< column "(2 - 1 << 0 & 3) <=> ((((2 - 1) << 0) & 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 / 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 / 0) <=> ((2 / 0))
< result
< column "(2 / 0) <=> ((2 / 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 1 * 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 1 * 2) <=> ((NULL = (1 * 2)))
< result
< column "(NULL = 1 * 2) <=> ((NULL = (1 * 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL || NULL ^ NULL MOD 2 << NOT NULL < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 39 near "NOT NULL < NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0) <=> (0)
< result
< column "(0) <=> (0)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 & 0 || 1 << 1 NOT IN (2 NOT IN (1, NULL) AND 3, (NULL) NOT BETWEEN 1 * 2 AND NOT 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 89 near "NOT 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 | NULL IN (NULL << 1, 1 IS NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 | NULL IN (NULL << 1, 1 IS NULL)) <=> (((1 | NULL) IN ((NULL << 1), (1 IS NULL))))
< result
< column "(1 | NULL IN (NULL << 1, 1 IS NULL)) <=> (((1 | NULL) IN ((NULL << 1), (1 IS NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! @v := NULL / + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 14 near ":= NULL / + 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2) <=> ((NOT 2))
< result
< column "(NOT 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS UNKNOWN * (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "* (1)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 != NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 != NULL) <=> ((0 <> NULL))
< result
< column "(0 != NULL) <=> ((0 <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 NOT IN (3 - 3, 1 >= 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 2 NOT IN (3 - 3, 1 >= 2)) <=> (((NOT 2) NOT IN ((3 - 3), (1 >= 2))))
< result
< column "(! 2 NOT IN (3 - 3, 1 >= 2)) <=> (((NOT 2) NOT IN ((3 - 3), (1 >= 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NULL BETWEEN NULL NOT BETWEEN 0 AND 2 AND NULL <=> NULL NOT IN (+ 0 || + 3, NULL | NULL < 1 >= 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "BETWEEN NULL NOT BETWEEN 0 AND 2 AND NULL <=> NULL NOT IN (+ 0 || + 3, NULL | NULL < 1 >= 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 IS TRUE IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND 1 && NULL - 0 << 3 << 3 > 2 <> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND 1 && NULL - 0 << 3 << 3 > 2 <> 0) <=> (((NULL AND 1) AND (((((NULL - 0) << 3) << 3) > 2) <> 0)))
< result
< column "(NULL AND 1 && NULL - 0 << 3 << 3 > 2 <> 0) <=> (((NULL AND 1) AND (((((NULL - 0) << 3) << 3) > 2) <> 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 AND 0 REGEXP 3 IS TRUE / NULL ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "/ NULL ^ 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> NULL) <=> ((0 >> NULL))
< result
< column "(0 >> NULL) <=> ((0 >> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 2 DIV 0 && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 XOR 2 DIV 0 && 2) <=> ((1 XOR ((2 DIV 0) AND 2)))
< result
< column "(1 XOR 2 DIV 0 && 2) <=> ((1 XOR ((2 DIV 0) AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 + NULL) <=> ((1 + NULL))
< result
< column "(1 + NULL) <=> ((1 + NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 REGEXP 3 NOT REGEXP NULL >= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "NOT REGEXP NULL >= 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 IS NULL IN (1 IN (0, 1) && 0 IS TRUE, 3 IS NULL XOR - 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "IN (1 IN (0, 1) && 0 IS TRUE, 3 IS NULL XOR - 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT BETWEEN NULL AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT BETWEEN NULL AND 2) <=> ((3 NOT BETWEEN NULL AND 2))
< result
< column "(3 NOT BETWEEN NULL AND 2) <=> ((3 NOT BETWEEN NULL AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ ! ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ ! ! 1) <=> ((~ (NOT (NOT 1))))
< result
< column "(~ ! ! 1) <=> ((~ (NOT (NOT 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL < 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL < 0) <=> ((NULL < 0))
< result
< column "(NULL < 0) <=> ((NULL < 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 IS NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0 IS NULL IS NULL) <=> ((((~ 0) IS NULL) IS NULL))
< result
< column "(~ 0 IS NULL IS NULL) <=> ((((~ 0) IS NULL) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 != 3 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 != 3 IS UNKNOWN) <=> (((2 <> 3) IS UNKNOWN))
< result
< column "(2 != 3 IS UNKNOWN) <=> (((2 <> 3) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1 OR 2 IN (- NULL, 1 IN (0, 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1 OR 2 IN (- NULL, 1 IN (0, 2))) <=> (((~ 1) OR (2 IN ((- NULL), (1 IN (0, 2))))))
< result
< column "(~ 1 OR 2 IN (- NULL, 1 IN (0, 2))) <=> (((~ 1) OR (2 IN ((- NULL), (1 IN (0, 2))))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ + 0 - 2 DIV NULL << NULL & 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ + 0 - 2 DIV NULL << NULL & 0) <=> (((((~ (+ 0)) - (2 DIV NULL)) << NULL) & 0))
< result
< column "(~ + 0 - 2 DIV NULL << NULL & 0) <=> (((((~ (+ 0)) - (2 DIV NULL)) << NULL) & 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 BETWEEN 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 BETWEEN 3 AND 0) <=> ((3 BETWEEN 3 AND 0))
< result
< column "(3 BETWEEN 3 AND 0) <=> ((3 BETWEEN 3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT FALSE DIV NULL OR 1 DIV NOT NULL = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "DIV NULL OR 1 DIV NOT NULL = 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 | 2 * 0 NOT BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 | 2 * 0 NOT BETWEEN 2 AND 3) <=> (((0 | (2 * 0)) NOT BETWEEN 2 AND 3))
< result
< column "(0 | 2 * 0 NOT BETWEEN 2 AND 3) <=> (((0 | (2 * 0)) NOT BETWEEN 2 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT NULL NOT BETWEEN 1 <= NULL AND 3 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "NOT BETWEEN 1 <= NULL AND 3 IS UNKNOWN" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 % 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 % 1) <=> ((1 % 1))
< result
< column "(1 % 1) <=> ((1 % 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ (0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ (0)) <=> ((~ 0))
< result
< column "(~ (0)) <=> ((~ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 IS UNKNOWN NOT LIKE 1 IS NULL IN (1 <= 2 * NULL > 2, 1 NOT BETWEEN 3 AND 1 XOR NULL XOR 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "NOT LIKE 1 IS NULL IN (1 <= 2 * NULL > 2, 1 NOT BETWEEN 3 AND 1 XOR NULL XOR 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = 0 IN (NULL > 0, 1 <> 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 = 0 IN (NULL > 0, 1 <> 3)) <=> ((3 = (0 IN ((NULL > 0), (1 <> 3)))))
< result
< column "(3 = 0 IN (NULL > 0, 1 <> 3)) <=> ((3 = (0 IN ((NULL > 0), (1 <> 3)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (NULL, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IN (NULL, 0)) <=> ((2 IN (NULL, 0)))
< result
< column "(2 IN (NULL, 0)) <=> ((2 IN (NULL, 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1) <=> ((- 1))
< result
< column "(- 1) <=> ((- 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL % 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL % 2) <=> ((NULL % 2))
< result
< column "(NULL % 2) <=> ((NULL % 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <= NULL NOT LIKE 3 MOD NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 31 near "MOD NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 != 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 != 1) <=> ((0 <> 1))
< result
< column "(0 != 1) <=> ((0 <> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 3) <=> ((+ 3))
< result
< column "(+ 3) <=> ((+ 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 * 3 NOT LIKE NULL MOD 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 30 near "MOD 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ NULL LIKE 1 IS NOT NULL BETWEEN 2 NOT IN (NULL, 2) NOT IN (0 IN (NULL, 1), 2 > 2) AND NOT NULL = 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "BETWEEN 2 NOT IN (NULL, 2) NOT IN (0 IN (NULL, 1), 2 > 2) AND NOT NULL = 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 << 1 < + 0 < 1 - 2 DIV 3 IN (3, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 << 1 < + 0 < 1 - 2 DIV 3 IN (3, NULL)) <=> ((((1 << 1) < (+ 0)) < ((1 - (2 DIV 3)) IN (3, NULL))))
< result
< column "(1 << 1 < + 0 < 1 - 2 DIV 3 IN (3, NULL)) <=> ((((1 << 1) < (+ 0)) < ((1 - (2 DIV 3)) IN (3, NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <> NULL) <=> ((1 <> NULL))
< result
< column "(1 <> NULL) <=> ((1 <> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IN (NULL, 2) IS NOT FALSE BETWEEN 3 AND NULL - - 2 AND @v := NULL >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 45 near "BETWEEN 3 AND NULL - - 2 AND @v := NULL >= 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV 3) <=> ((2 DIV 3))
< result
< column "(2 DIV 3) <=> ((2 DIV 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 < NULL || 1 DIV 0 >> - 2 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 < NULL || 1 DIV 0 >> - 2 IS NOT NULL) <=> (((2 < NULL) OR (((1 DIV 0) >> (- 2)) IS NOT NULL)))
< result
< column "(2 < NULL || 1 DIV 0 >> - 2 IS NOT NULL) <=> (((2 < NULL) OR (((1 DIV 0) >> (- 2)) IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2 NOT BETWEEN 1 AND 0) IN (2 - NULL NOT BETWEEN 3 BETWEEN 3 AND NULL AND NULL AND NULL, 1 >= 3 - 2 DIV 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 65 near "BETWEEN 3 AND NULL AND NULL AND NULL, 1 >= 3 - 2 DIV 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL LIKE NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL LIKE NULL) <=> ((NOT (NULL LIKE NULL)))
< result
< column "(NOT NULL LIKE NULL) <=> ((NOT (NULL LIKE NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 3 IS NULL NOT REGEXP 0 DIV 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 22 near "NOT REGEXP 0 DIV 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ ~ ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ ~ ~ NULL) <=> ((~ (~ (~ NULL))))
< result
< column "(~ ~ ~ NULL) <=> ((~ (~ (~ NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL) <=> ((NOT NULL))
< result
< column "(NOT NULL) <=> ((NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 <= 3 + 2 / 3 = 0 XOR 1 > 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 <= 3 + 2 / 3 = 0 XOR 1 > 1) <=> ((((2 <= (3 + (2 / 3))) = 0) XOR (1 > 1)))
< result
< column "(2 <= 3 + 2 / 3 = 0 XOR 1 > 1) <=> ((((2 <= (3 + (2 / 3))) = 0) XOR (1 > 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 NOT BETWEEN NULL AND 3 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 NOT BETWEEN NULL AND 3 IS UNKNOWN) <=> ((((NOT 0) NOT BETWEEN NULL AND 3) IS UNKNOWN))
< result
< column "(! 0 NOT BETWEEN NULL AND 3 IS UNKNOWN) <=> ((((NOT 0) NOT BETWEEN NULL AND 3) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 0 AND 2 NOT BETWEEN 3 / NULL AND 1 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 0 AND 2 NOT BETWEEN 3 / NULL AND 1 AND 1) <=> (((NULL NOT BETWEEN 0 AND (2 NOT BETWEEN (3 / NULL) AND 1)) AND 1))
< result
< column "(NULL NOT BETWEEN 0 AND 2 NOT BETWEEN 3 / NULL AND 1 AND 1) <=> (((NULL NOT BETWEEN 0 AND (2 NOT BETWEEN (3 / NULL) AND 1)) AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 - NULL NOT BETWEEN 0 BETWEEN 2 AND 0 AND 0 > NULL >> ! ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "BETWEEN 2 AND 0 AND 0 > NULL >> ! ! 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 >> 0 / 0 < 0 MOD NULL MOD 2 <> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 >> 0 / 0 < 0 MOD NULL MOD 2 <> NULL) <=> ((NOT (((3 >> (0 / 0)) < ((0 % NULL) % 2)) <> NULL)))
< result
< column "(NOT 3 >> 0 / 0 < 0 MOD NULL MOD 2 <> NULL) <=> ((NOT (((3 >> (0 / 0)) < ((0 % NULL) % 2)) <> NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IN (3, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IN (3, 0)) <=> ((1 IN (3, 0)))
< result
< column "(1 IN (3, 0)) <=> ((1 IN (3, 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL * 1 | 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL * 1 | 0) <=> ((((- NULL) * 1) | 0))
< result
< column "(- NULL * 1 | 0) <=> ((((- NULL) * 1) | 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= 2) <=> ((3 >= 2))
< result
< column "(3 >= 2) <=> ((3 >= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= 3) <=> ((2 >= 3))
< result
< column "(2 >= 3) <=> ((2 >= 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0) <=> ((~ 0))
< result
< column "(~ 0) <=> ((~ 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 OR NULL AND 0 >> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 OR NULL AND 0 >> 3) <=> ((1 OR (NULL AND (0 >> 3))))
< result
< column "(1 OR NULL AND 0 >> 3) <=> ((1 OR (NULL AND (0 >> 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL LIKE 2 XOR 0 LIKE 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL LIKE 2 XOR 0 LIKE 3) <=> (((NULL LIKE 2) XOR (0 LIKE 3)))
< result
< column "(NULL LIKE 2 XOR 0 LIKE 3) <=> (((NULL LIKE 2) XOR (0 LIKE 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 = 2 NOT BETWEEN 3 <= NULL AND NULL | 1 OR 3 IS NULL OR NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "<= NULL AND NULL | 1 OR 3 IS NULL OR NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 & 1 IN (! 3, NULL IN (3, 3)) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 & 1 IN (! 3, NULL IN (3, 3)) IS NOT FALSE) <=> ((((1 & 1) IN ((NOT 3), (NULL IN (3, 3)))) IS NOT FALSE))
< result
< column "(1 & 1 IN (! 3, NULL IN (3, 3)) IS NOT FALSE) <=> ((((1 & 1) IN ((NOT 3), (NULL IN (3, 3)))) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NULL XOR 2 * 3 IN (2 BETWEEN 3 AND NULL << ~ 1, 1 - NULL > 2 <=> 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 IS NULL XOR 2 * 3 IN (2 BETWEEN 3 AND NULL << ~ 1, 1 - NULL > 2 <=> 1)) <=> (((2 IS NULL) XOR ((2 * 3) IN ((2 BETWEEN 3 AND (NULL << (~ 1))), (((1 - NULL) > 2) <=> 1)))))
< result
< column "(2 IS NULL XOR 2 * 3 IN (2 BETWEEN 3 AND NULL << ~ 1, 1 - NULL > 2 <=> 1)) <=> (((2 IS NULL) XOR ((2 * 3) IN ((2 BETWEEN 3 AND (NULL << (~ 1))), (((1 - NULL) > 2) <=> 1)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL <=> ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL <=> ! 0) <=> (((- NULL) <=> (NOT 0)))
< result
< column "(- NULL <=> ! 0) <=> (((- NULL) <=> (NOT 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 - 1) <=> ((3 - 1))
< result
< column "(3 - 1) <=> ((3 - 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <> 3) <=> ((3 <> 3))
< result
< column "(3 <> 3) <=> ((3 <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> 3 && 0 IS TRUE - + 1 NOT IN (NULL, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "- + 1 NOT IN (NULL, NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 - 1) <=> ((0 - 1))
< result
< column "(0 - 1) <=> ((0 - 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT IN (NULL, 0) NOT BETWEEN NULL NOT BETWEEN 2 AND 0 AND NULL * 2 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "NOT BETWEEN NULL NOT BETWEEN 2 AND 0 AND NULL * 2 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >= 1 IS NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >= 1 IS NULL IS NULL) <=> ((((0 >= 1) IS NULL) IS NULL))
< result
< column "(0 >= 1 IS NULL IS NULL) <=> ((((0 >= 1) IS NULL) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NULL & 3 && NULL * NULL | NULL OR - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 18 near "& 3 && NULL * NULL | NULL OR - 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! + 1 <> 3 >= 2 >= 1 NOT BETWEEN NULL AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! + 1 <> 3 >= 2 >= 1 NOT BETWEEN NULL AND NULL) <=> (((((NOT (+ 1)) <> 3) >= 2) >= (1 NOT BETWEEN NULL AND NULL)))
< result
< column "(! + 1 <> 3 >= 2 >= 1 NOT BETWEEN NULL AND NULL) <=> (((((NOT (+ 1)) <> 3) >= 2) >= (1 NOT BETWEEN NULL AND NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && 0) <=> ((0 AND 0))
< result
< column "(0 && 0) <=> ((0 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 <=> 0 MOD ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2 <=> 0 MOD ~ 2) <=> (((+ 2) <=> (0 % (~ 2))))
< result
< column "(+ 2 <=> 0 MOD ~ 2) <=> (((+ 2) <=> (0 % (~ 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <> 3 <> 3 XOR 2 IN ((NULL >= NULL), 1 IS NOT NULL REGEXP 0 <= 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 68 near "REGEXP 0 <= 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL || 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL || 1) <=> ((NULL OR 1))
< result
< column "(NULL || 1) <=> ((NULL OR 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < NULL) <=> ((1 < NULL))
< result
< column "(1 < NULL) <=> ((1 < NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL & 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL & 1) <=> ((NULL & 1))
< result
< column "(NULL & 1) <=> ((NULL & 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL ^ NULL IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL ^ NULL IS NULL) <=> ((((NOT NULL) ^ NULL) IS NULL))
< result
< column "(! NULL ^ NULL IS NULL) <=> ((((NOT NULL) ^ NULL) IS NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT BETWEEN NULL AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT BETWEEN NULL AND 1) <=> ((1 NOT BETWEEN NULL AND 1))
< result
< column "(1 NOT BETWEEN NULL AND 1) <=> ((1 NOT BETWEEN NULL AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + 1) <=> ((3 + 1))
< result
< column "(3 + 1) <=> ((3 + 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL XOR 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL XOR 0) <=> ((NULL XOR 0))
< result
< column "(NULL XOR 0) <=> ((NULL XOR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> 3) <=> ((0 >> 3))
< result
< column "(0 >> 3) <=> ((0 >> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 = 2 && ! NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 = 2 && ! NULL) <=> (((1 = 2) AND (NOT NULL)))
< result
< column "(1 = 2 && ! NULL) <=> (((1 = 2) AND (NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 + 3 IN (2 >> 0, 1 % 0) IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 + 3 IN (2 >> 0, 1 % 0) IS NOT NULL) <=> ((((0 + 3) IN ((2 >> 0), (1 % 0))) IS NOT NULL))
< result
< column "(0 + 3 IN (2 >> 0, 1 % 0) IS NOT NULL) <=> ((((0 + 3) IN ((2 >> 0), (1 % 0))) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 / 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 / 1) <=> ((2 / 1))
< result
< column "(2 / 1) <=> ((2 / 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (3 OR NULL) OR 3 NOT IN (3, 2) << ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "<< ~ 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT NULL > 2 = NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT NULL > 2 = NULL) <=> ((((1 IS NOT NULL) > 2) = NULL))
< result
< column "(1 IS NOT NULL > 2 = NULL) <=> ((((1 IS NOT NULL) > 2) = NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >> 1) <=> ((1 >> 1))
< result
< column "(1 >> 1) <=> ((1 >> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV 0 IN (NULL NOT BETWEEN 2 AND 0, + 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV 0 IN (NULL NOT BETWEEN 2 AND 0, + 1)) <=> (((2 DIV 0) IN ((NULL NOT BETWEEN 2 AND 0), (+ 1))))
< result
< column "(2 DIV 0 IN (NULL NOT BETWEEN 2 AND 0, + 1)) <=> (((2 DIV 0) IN ((NULL NOT BETWEEN 2 AND 0), (+ 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0) <=> ((NOT 0))
< result
< column "(! 0) <=> ((NOT 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 << 0 IN (2 AND 2, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2 << 0 IN (2 AND 2, 0)) <=> ((((+ 2) << 0) IN ((2 AND 2), 0)))
< result
< column "(+ 2 << 0 IN (2 AND 2, 0)) <=> ((((+ 2) << 0) IN ((2 AND 2), 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 + NULL) <=> ((0 + NULL))
< result
< column "(0 + NULL) <=> ((0 + NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL) / ~ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL) / ~ 3) <=> ((NULL / (~ 3)))
< result
< column "((NULL) / ~ 3) <=> ((NULL / (~ 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 != 3 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 != 3 IS TRUE) <=> (((1 <> 3) IS TRUE))
< result
< column "(1 != 3 IS TRUE) <=> (((1 <> 3) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 % 0 NOT BETWEEN NULL NOT LIKE 0 IN (1 NOT IN (NULL, 1), NULL IS NULL) AND 1 * 1 <= 2 % 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 35 near "NOT LIKE 0 IN (1 NOT IN (NULL, 1), NULL IS NULL) AND 1 * 1 <= 2 % 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 12 near "NOT 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0 ^ 1 AND 0 NOT IN (NULL, 0) IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0 ^ 1 AND 0 NOT IN (NULL, 0) IS NOT NULL) <=> ((((NOT 0) ^ 1) AND ((0 NOT IN (NULL, 0)) IS NOT NULL)))
< result
< column "(! 0 ^ 1 AND 0 NOT IN (NULL, 0) IS NOT NULL) <=> ((((NOT 0) ^ 1) AND ((0 NOT IN (NULL, 0)) IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT @v := NULL <> NULL <> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (@v := NULL <> NULL <> 1) <=> ((@v := ((NULL <> NULL) <> 1)))
< result
< column "(@v := NULL <> NULL <> 1) <=> ((@v := ((NULL <> NULL) <> 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IN (3, 1) * 0 BETWEEN 0 AND NULL NOT BETWEEN 3 BETWEEN NULL AND 3 & NULL AND 2 AND 2 | 3 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 20 near "* 0 BETWEEN 0 AND NULL NOT BETWEEN 3 BETWEEN NULL AND 3 & NULL AND 2 AND 2 | 3 IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 < 1) <=> ((1 < 1))
< result
< column "(1 < 1) <=> ((1 < 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL / NULL * 1 * 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL / NULL * 1 * 0) <=> ((((NULL / NULL) * 1) * 0))
< result
< column "(NULL / NULL * 1 * 0) <=> ((((NULL / NULL) * 1) * 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL) <=> ((+ NULL))
< result
< column "(+ NULL) <=> ((+ NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN NULL AND 2 % (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN NULL AND 2 % (2)) <=> ((NULL BETWEEN NULL AND (2 % 2)))
< result
< column "(NULL BETWEEN NULL AND 2 % (2)) <=> ((NULL BETWEEN NULL AND (2 % 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 DIV 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 DIV 0) <=> ((2 DIV 0))
< result
< column "(2 DIV 0) <=> ((2 DIV 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 != 0 NOT IN (1 IN (1, 1), ~ 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 != 0 NOT IN (1 IN (1, 1), ~ 3)) <=> ((3 <> (0 NOT IN ((1 IN (1, 1)), (~ 3)))))
< result
< column "(3 != 0 NOT IN (1 IN (1, 1), ~ 3)) <=> ((3 <> (0 NOT IN ((1 IN (1, 1)), (~ 3)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT FALSE >> 3 < 3 * NULL ^ NULL < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near ">> 3 < 3 * NULL ^ NULL < 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + NULL || 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ NULL || 0 IS TRUE) <=> (((+ NULL) OR (0 IS TRUE)))
< result
< column "(+ NULL || 0 IS TRUE) <=> (((+ NULL) OR (0 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 2 IN (3, NULL) = 1 NOT IN (3, 1) / 3 OR 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 43 near "/ 3 OR 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <= 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL <= 2) <=> ((NULL <= 2))
< result
< column "(NULL <= 2) <=> ((NULL <= 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 << 3 << NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 << 3 << NULL) <=> (((0 << 3) << NULL))
< result
< column "(0 << 3 << NULL) <=> (((0 << 3) << NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3 NOT IN (1, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3 NOT IN (1, 1)) <=> ((NOT (3 NOT IN (1, 1))))
< result
< column "(NOT 3 NOT IN (1, 1)) <=> ((NOT (3 NOT IN (1, 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 & 1 = 1 / 0 <> 0 || NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0 & 1 = 1 / 0 <> 0 || NULL) <=> ((((((~ 0) & 1) = (1 / 0)) <> 0) OR NULL))
< result
< column "(~ 0 & 1 = 1 / 0 <> 0 || NULL) <=> ((((((~ 0) & 1) = (1 / 0)) <> 0) OR NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 || 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 || 1 IS NOT NULL) <=> ((2 OR (1 IS NOT NULL)))
< result
< column "(2 || 1 IS NOT NULL) <=> ((2 OR (1 IS NOT NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1) <=> ((NOT 1))
< result
< column "(! 1) <=> ((NOT 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> 1) <=> ((1 <=> 1))
< result
< column "(1 <=> 1) <=> ((1 <=> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 << NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 << NULL) <=> ((1 << NULL))
< result
< column "(1 << NULL) <=> ((1 << NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 << 2 <> NULL NOT LIKE 0 BETWEEN NULL / 3 XOR 0 || 1 AND (0 <=> 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "BETWEEN NULL / 3 XOR 0 || 1 AND (0 <=> 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 2 AND 1 && 0 IS NOT FALSE OR ! 3 NOT LIKE 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 2 AND 1 && 0 IS NOT FALSE OR ! 3 NOT LIKE 1) <=> ((((NULL NOT BETWEEN 2 AND 1) AND (0 IS NOT FALSE)) OR ((NOT 3) NOT LIKE 1)))
< result
< column "(NULL NOT BETWEEN 2 AND 1 && 0 IS NOT FALSE OR ! 3 NOT LIKE 1) <=> ((((NULL NOT BETWEEN 2 AND 1) AND (0 IS NOT FALSE)) OR ((NOT 3) NOT LIKE 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + (NULL) ^ NULL << 0 IN (3 = 1, 3 && 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ (NULL) ^ NULL << 0 IN (3 = 1, 3 && 0)) <=> (((((+ NULL) ^ NULL) << 0) IN ((3 = 1), (3 AND 0))))
< result
< column "(+ (NULL) ^ NULL << 0 IN (3 = 1, 3 && 0)) <=> (((((+ NULL) ^ NULL) << 0) IN ((3 = 1), (3 AND 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 & 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 & 2) <=> ((1 & 2))
< result
< column "(1 & 2) <=> ((1 & 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 < 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 < 2) <=> ((2 < 2))
< result
< column "(2 < 2) <=> ((2 < 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > 3 IN (1 NOT IN (0, 3), 2 << NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL > 3 IN (1 NOT IN (0, 3), 2 << NULL)) <=> ((NULL > (3 IN ((1 NOT IN (0, 3)), (2 << NULL)))))
< result
< column "(NULL > 3 IN (1 NOT IN (0, 3), 2 << NULL)) <=> ((NULL > (3 IN ((1 NOT IN (0, 3)), (2 << NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 LIKE 1 IS UNKNOWN < NULL <> NULL OR NULL > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "< NULL <> NULL OR NULL > 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND 1 OR 0 & 2 >= NULL IN (2, 3) NOT IN (1 >> NULL, 0 LIKE 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 45 near "NOT IN (1 >> NULL, 0 LIKE 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS UNKNOWN | (3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "| (3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 MOD 1 ^ 1 NOT IN (0, 1) NOT BETWEEN 1 IN (3, 1) <> ~ NULL AND 3 IS TRUE <> 3 ^ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 36 near "NOT BETWEEN 1 IN (3, 1) <> ~ NULL AND 3 IS TRUE <> 3 ^ 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >= 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >= 3) <=> ((3 >= 3))
< result
< column "(3 >= 3) <=> ((3 >= 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT LIKE 2) <=> ((0 NOT LIKE 2))
< result
< column "(0 NOT LIKE 2) <=> ((0 NOT LIKE 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (1 NOT BETWEEN 0 AND 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((1 NOT BETWEEN 0 AND 3)) <=> ((1 NOT BETWEEN 0 AND 3))
< result
< column "((1 NOT BETWEEN 0 AND 3)) <=> ((1 NOT BETWEEN 0 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 2 > 3) <=> (((- 2) > 3))
< result
< column "(- 2 > 3) <=> (((- 2) > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL OR 1 || NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL OR 1 || NULL) <=> (((NULL OR 1) OR NULL))
< result
< column "(NULL OR 1 || NULL) <=> (((NULL OR 1) OR NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 LIKE 3 DIV 1 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 19 near "DIV 1 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 3 AND 2 MOD 2 IS NOT FALSE NOT IN (1 XOR 0 IN (3 OR 2, ~ 3), - 1 IS TRUE)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 51 near "NOT IN (1 XOR 0 IN (3 OR 2, ~ 3), - 1 IS TRUE)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL NOT BETWEEN NULL AND 3) >> 3 <> 3 <=> - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL NOT BETWEEN NULL AND 3) >> 3 <> 3 <=> - 2) <=> (((((NULL NOT BETWEEN NULL AND 3) >> 3) <> 3) <=> (- 2)))
< result
< column "((NULL NOT BETWEEN NULL AND 3) >> 3 <> 3 <=> - 2) <=> (((((NULL NOT BETWEEN NULL AND 3) >> 3) <> 3) <=> (- 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (+ 1 MOD - 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((+ 1 MOD - 0)) <=> (((+ 1) % (- 0)))
< result
< column "((+ 1 MOD - 0)) <=> (((+ 1) % (- 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL XOR 2 < NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL XOR 2 < NULL) <=> (((- NULL) XOR (2 < NULL)))
< result
< column "(- NULL XOR 2 < NULL) <=> (((- NULL) XOR (2 < NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (0, NULL) NOT IN (3 BETWEEN 3 AND 3, 2 | 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "NOT IN (3 BETWEEN 3 AND 3, 2 | 3)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 % 0) <=> ((0 % 0))
< result
< column "(0 % 0) <=> ((0 % 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 & 1 OR 0 NOT BETWEEN 2 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 & 1 OR 0 NOT BETWEEN 2 AND 2) <=> (((1 & 1) OR (0 NOT BETWEEN 2 AND 2)))
< result
< column "(1 & 1 OR 0 NOT BETWEEN 2 AND 2) <=> (((1 & 1) OR (0 NOT BETWEEN 2 AND 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 BETWEEN 2 AND 3 != NULL && NULL NOT IN (NULL <= 0, 2 DIV 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 BETWEEN 2 AND 3 != NULL && NULL NOT IN (NULL <= 0, 2 DIV 2)) <=> (((NOT ((0 BETWEEN 2 AND 3) <> NULL)) AND (NULL NOT IN ((NULL <= 0), (2 DIV 2)))))
< result
< column "(NOT 0 BETWEEN 2 AND 3 != NULL && NULL NOT IN (NULL <= 0, 2 DIV 2)) <=> (((NOT ((0 BETWEEN 2 AND 3) <> NULL)) AND (NULL NOT IN ((NULL <= 0), (2 DIV 2)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <> 3) <=> ((1 <> 3))
< result
< column "(1 <> 3) <=> ((1 <> 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 - 0) <=> ((1 - 0))
< result
< column "(1 - 0) <=> ((1 - 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL * 0 <=> 3 IS NULL IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL * 0 <=> 3 IS NULL IS NOT NULL) <=> (((((NULL * 0) <=> 3) IS NULL) IS NOT NULL))
< result
< column "(NULL * 0 <=> 3 IS NULL IS NOT NULL) <=> (((((NULL * 0) <=> 3) IS NULL) IS NOT NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 != NULL % + 2 = 0 ^ 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 != NULL % + 2 = 0 ^ 0) <=> ((NOT ((0 <> (NULL % (+ 2))) = (0 ^ 0))))
< result
< column "(NOT 0 != NULL % + 2 = 0 ^ 0) <=> ((NOT ((0 <> (NULL % (+ 2))) = (0 ^ 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 2) <=> ((- 2))
< result
< column "(- 2) <=> ((- 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 <= 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 <= 0) <=> ((0 <= 0))
< result
< column "(0 <= 0) <=> ((0 <= 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 XOR 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1 XOR 3) <=> (((+ 1) XOR 3))
< result
< column "(+ 1 XOR 3) <=> (((+ 1) XOR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 2 NOT REGEXP NULL BETWEEN NULL AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 34 near "BETWEEN NULL AND 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 IS NULL <= 2 <=> 2 MOD 2 << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 1 IS NULL <= 2 <=> 2 MOD 2 << 3) <=> (((((- 1) IS NULL) <= 2) <=> ((2 % 2) << 3)))
< result
< column "(- 1 IS NULL <= 2 <=> 2 MOD 2 << 3) <=> (((((- 1) IS NULL) <= 2) <=> ((2 % 2) << 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1) <=> ((~ 1))
< result
< column "(~ 1) <=> ((~ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL <> 1 IS TRUE >= 1 NOT IN (3, 1) ^ 2 BETWEEN 2 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near ">= 1 NOT IN (3, 1) ^ 2 BETWEEN 2 AND 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 XOR 1 NOT REGEXP NOT 0 <= 0 LIKE 1 = 2 >> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "NOT 0 <= 0 LIKE 1 = 2 >> 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL LIKE 3 && 1 - 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL LIKE 3 && 1 - 0) <=> (((NULL LIKE 3) AND (1 - 0)))
< result
< column "(NULL LIKE 3 && 1 - 0) <=> (((NULL LIKE 3) AND (1 - 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 > 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 > 1) <=> ((0 > 1))
< result
< column "(0 > 1) <=> ((0 > 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL) <=> (NULL)
< result
< column "(NULL) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ (1) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ (1) IS NOT FALSE) <=> (((~ 1) IS NOT FALSE))
< result
< column "(~ (1) IS NOT FALSE) <=> (((~ 1) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1) <=> ((~ 1))
< result
< column "(~ 1) <=> ((~ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 NOT BETWEEN 2 AND 0 XOR 1 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 NOT BETWEEN 2 AND 0 XOR 1 <> 3) <=> (((NOT (0 NOT BETWEEN 2 AND 0)) XOR (1 <> 3)))
< result
< column "(NOT 0 NOT BETWEEN 2 AND 0 XOR 1 <> 3) <=> (((NOT (0 NOT BETWEEN 2 AND 0)) XOR (1 <> 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 > 2 <=> (NULL) NOT BETWEEN NULL AND 3 IS NOT FALSE <= 1 < 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 62 near "<= 1 < 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL ^ 1 IN (@v := 1, 0 >> NULL) NOT LIKE 0 NOT LIKE 3 < 2 || 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 43 near "NOT LIKE 0 NOT LIKE 3 < 2 || 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 <= 0 < 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2 <= 0 < 1 IS UNKNOWN) <=> (((((+ 2) <= 0) < 1) IS UNKNOWN))
< result
< column "(+ 2 <= 0 < 1 IS UNKNOWN) <=> (((((+ 2) <= 0) < 1) IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2) <=> (2)
< result
< column "(2) <=> (2)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 MOD 3 NOT REGEXP 2 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 MOD 3 NOT REGEXP 2 IS UNKNOWN) <=> ((NOT (((1 % 3) NOT REGEXP 2) IS UNKNOWN)))
< result
< column "(NOT 1 MOD 3 NOT REGEXP 2 IS UNKNOWN) <=> ((NOT (((1 % 3) NOT REGEXP 2) IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 <> 3 >> 2 NOT BETWEEN 0 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1 <> 3 >> 2 NOT BETWEEN 0 AND 1) <=> (((NOT 1) <> ((3 >> 2) NOT BETWEEN 0 AND 1)))
< result
< column "(! 1 <> 3 >> 2 NOT BETWEEN 0 AND 1) <=> (((NOT 1) <> ((3 >> 2) NOT BETWEEN 0 AND 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 * 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1 * 3) <=> (((NOT 1) * 3))
< result
< column "(! 1 * 3) <=> (((NOT 1) * 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS UNKNOWN) <=> ((1 IS UNKNOWN))
< result
< column "(1 IS UNKNOWN) <=> ((1 IS UNKNOWN))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT NULL = (1) + 0 NOT BETWEEN 3 AND NULL ^ 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT NULL = (1) + 0 NOT BETWEEN 3 AND NULL ^ 1 IS UNKNOWN) <=> ((NOT ((NULL = ((1 + 0) NOT BETWEEN 3 AND (NULL ^ 1))) IS UNKNOWN)))
< result
< column "(NOT NULL = (1) + 0 NOT BETWEEN 3 AND NULL ^ 1 IS UNKNOWN) <=> ((NOT ((NULL = ((1 + 0) NOT BETWEEN 3 AND (NULL ^ 1))) IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / NULL <> 1 || 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 / NULL <> 1 || 3) <=> ((((1 / NULL) <> 1) OR 3))
< result
< column "(1 / NULL <> 1 || 3) <=> ((((1 / NULL) <> 1) OR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 0 >> 2 REGEXP 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 0 >> 2 REGEXP 1) <=> (((0 >> 2) NOT REGEXP 1))
< result
< column "(NOT 0 >> 2 REGEXP 1) <=> (((0 >> 2) NOT REGEXP 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 << 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 << 1) <=> ((1 << 1))
< result
< column "(1 << 1) <=> ((1 << 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >> 0 <= 3 | 1 BETWEEN 3 BETWEEN 2 AND 1 >= 3 * 1 AND ! 1 = 0 NOT IN (0, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "BETWEEN 2 AND 1 >= 3 * 1 AND ! 1 = 0 NOT IN (0, NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND NULL NOT IN (0, 1 NOT BETWEEN 3 AND 2) XOR 2 BETWEEN 1 AND 0 = 1 LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 AND NULL NOT IN (0, 1 NOT BETWEEN 3 AND 2) XOR 2 BETWEEN 1 AND 0 = 1 LIKE 2) <=> (((3 AND (NULL NOT IN (0, (1 NOT BETWEEN 3 AND 2)))) XOR ((2 BETWEEN 1 AND 0) = (1 LIKE 2))))
< result
< column "(3 AND NULL NOT IN (0, 1 NOT BETWEEN 3 AND 2) XOR 2 BETWEEN 1 AND 0 = 1 LIKE 2) <=> (((3 AND (NULL NOT IN (0, (1 NOT BETWEEN 3 AND 2)))) XOR ((2 BETWEEN 1 AND 0) = (1 LIKE 2))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1 IS TRUE IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 21 near "IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 << 0 >> 0 NOT BETWEEN 2 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 << 0 >> 0 NOT BETWEEN 2 AND NULL) <=> ((((2 << 0) >> 0) NOT BETWEEN 2 AND NULL))
< result
< column "(2 << 0 >> 0 NOT BETWEEN 2 AND NULL) <=> ((((2 << 0) >> 0) NOT BETWEEN 2 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 % 2 & NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 % 2 & NULL) <=> ((((NOT 3) % 2) & NULL))
< result
< column "(! 3 % 2 & NULL) <=> ((((NOT 3) % 2) & NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL XOR 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL XOR 1) <=> ((NULL XOR 1))
< result
< column "(NULL XOR 1) <=> ((NULL XOR 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 * NULL BETWEEN 1 && 3 AND 2 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "&& 3 AND 2 IS TRUE" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL << 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL << 3) <=> ((NULL << 3))
< result
< column "(NULL << 3) <=> ((NULL << 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL = 1) <=> ((NULL = 1))
< result
< column "(NULL = 1) <=> ((NULL = 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 1 AND 0 && 2 IS NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 1 AND 0 && 2 IS NULL) <=> (((2 BETWEEN 1 AND 0) AND (2 IS NULL)))
< result
< column "(2 BETWEEN 1 AND 0 && 2 IS NULL) <=> (((2 BETWEEN 1 AND 0) AND (2 IS NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 - 3 <= ! 1 BETWEEN NOT 3 IS UNKNOWN AND 3 <= NULL % 1 LIKE 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 31 near "NOT 3 IS UNKNOWN AND 3 <= NULL % 1 LIKE 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (3, 1) IN (3 BETWEEN 1 AND 2, - 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 25 near "IN (3 BETWEEN 1 AND 2, - 2)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 3) <=> ((NOT 3))
< result
< column "(NOT 3) <=> ((NOT 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 != 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 != 1) <=> ((0 <> 1))
< result
< column "(0 != 1) <=> ((0 <> 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 * 0 + 1 ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 * 0 + 1 ^ 3) <=> (((2 * 0) + (1 ^ 3)))
< result
< column "(2 * 0 + 1 ^ 3) <=> (((2 * 0) + (1 ^ 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL && NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL && NULL) <=> ((NULL AND NULL))
< result
< column "(NULL && NULL) <=> ((NULL AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ + 1) <=> ((+ (+ 1)))
< result
< column "(+ + 1) <=> ((+ (+ 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 <> NULL IN (3 NOT BETWEEN 2 AND 3, 0 MOD NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 <> NULL IN (3 NOT BETWEEN 2 AND 3, 0 MOD NULL)) <=> ((3 <> (NULL IN ((3 NOT BETWEEN 2 AND 3), (0 % NULL)))))
< result
< column "(3 <> NULL IN (3 NOT BETWEEN 2 AND 3, 0 MOD NULL)) <=> ((3 <> (NULL IN ((3 NOT BETWEEN 2 AND 3), (0 % NULL)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2) <=> ((+ 2))
< result
< column "(+ 2) <=> ((+ 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0 NOT BETWEEN 1 AND 1 >> 3 * 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0 NOT BETWEEN 1 AND 1 >> 3 * 3)) <=> ((0 NOT BETWEEN 1 AND (1 >> (3 * 3))))
< result
< column "((0 NOT BETWEEN 1 AND 1 >> 3 * 3)) <=> ((0 NOT BETWEEN 1 AND (1 >> (3 * 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - ! 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- ! 3) <=> ((- (NOT 3)))
< result
< column "(- ! 3) <=> ((- (NOT 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 OR 2 != NULL REGEXP 0 REGEXP 2 NOT BETWEEN 2 AND 0 + NULL && 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 37 near "REGEXP 2 NOT BETWEEN 2 AND 0 + NULL && 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + ~ 3 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ ~ 3 IS NOT FALSE) <=> (((+ (~ 3)) IS NOT FALSE))
< result
< column "(+ ~ 3 IS NOT FALSE) <=> (((+ (~ 3)) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / 2 / 2 - 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 / 2 / 2 - 2) <=> ((((1 / 2) / 2) - 2))
< result
< column "(1 / 2 / 2 - 2) <=> ((((1 / 2) / 2) - 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL XOR 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL XOR 3) <=> ((NULL XOR 3))
< result
< column "(NULL XOR 3) <=> ((NULL XOR 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 1 AND 0 DIV + 3 BETWEEN 2 BETWEEN 1 AND 2 & NULL OR 3 AND - 1 IN (0, NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 54 near "BETWEEN 1 AND 2 & NULL OR 3 AND - 1 IN (0, NULL)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 NOT BETWEEN 1 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 NOT BETWEEN 1 AND 1) <=> ((0 NOT BETWEEN 1 AND 1))
< result
< column "(0 NOT BETWEEN 1 AND 1) <=> ((0 NOT BETWEEN 1 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 % 0) <=> ((1 % 0))
< result
< column "(1 % 0) <=> ((1 % 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 IN (0, 1) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 2 IN (0, 1) IS NOT FALSE) <=> ((((- 2) IN (0, 1)) IS NOT FALSE))
< result
< column "(- 2 IN (0, 1) IS NOT FALSE) <=> ((((- 2) IN (0, 1)) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT + 1 DIV 0 && 3 IN (NULL, 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 1 DIV 0 && 3 IN (NULL, 0)) <=> ((((+ 1) DIV 0) AND (3 IN (NULL, 0))))
< result
< column "(+ 1 DIV 0 && 3 IN (NULL, 0)) <=> ((((+ 1) DIV 0) AND (3 IN (NULL, 0))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2) <=> ((NOT 2))
< result
< column "(NOT 2) <=> ((NOT 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 + 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 + 2) <=> ((3 + 2))
< result
< column "(3 + 2) <=> ((3 + 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 >> 3 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 >> 3 IS NOT FALSE) <=> (((3 >> 3) IS NOT FALSE))
< result
< column "(3 >> 3 IS NOT FALSE) <=> (((3 >> 3) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT REGEXP 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT REGEXP 2) <=> ((1 NOT REGEXP 2))
< result
< column "(1 NOT REGEXP 2) <=> ((1 NOT REGEXP 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1) <=> ((NOT 1))
< result
< column "(NOT 1) <=> ((NOT 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1) <=> (1)
< result
< column "(1) <=> (1)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 IS NOT NULL || 3 NOT BETWEEN 3 AND 0 > 0 IS NULL ^ NULL REGEXP 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 59 near "^ NULL REGEXP 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 0) <=> ((NOT 0))
< result
< column "(! 0) <=> ((NOT 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 >> 1 / 2 NOT LIKE 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 >> 1 / 2 NOT LIKE 2) <=> (((0 >> (1 / 2)) NOT LIKE 2))
< result
< column "(0 >> 1 / 2 NOT LIKE 2) <=> (((0 >> (1 / 2)) NOT LIKE 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 / 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 / 3) <=> ((0 / 3))
< result
< column "(0 / 3) <=> ((0 / 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 > 3) <=> ((2 > 3))
< result
< column "(2 > 3) <=> ((2 > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (NULL NOT IN (0, 2) NOT BETWEEN 1 NOT BETWEEN NULL AND 3 AND (NULL), 1 | 2 IN (1 || 3, 2 NOT LIKE 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 39 near "NOT BETWEEN 1 NOT BETWEEN NULL AND 3 AND (NULL), 1 | 2 IN (1 || 3, 2 NOT LIKE 2))" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IN (NULL, NULL) IS NOT NULL ^ NOT 1 ^ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "^ NOT 1 ^ 2" 
> exec
| SET @omni_sweep = ?
: string
| SELECT (3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((3)) <=> (3)
< result
< column "((3)) <=> (3)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 1 IS UNKNOWN * NOT 1 IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near "* NOT 1 IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 % 1 DIV 1 > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 % 1 DIV 1 > 3) <=> ((((0 % 1) DIV 1) > 3))
< result
< column "(0 % 1 DIV 1 > 3) <=> ((((0 % 1) DIV 1) > 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && 0 ^ 0 NOT LIKE 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && 0 ^ 0 NOT LIKE 0) <=> ((0 AND ((0 ^ 0) NOT LIKE 0)))
< result
< column "(0 && 0 ^ 0 NOT LIKE 0) <=> ((0 AND ((0 ^ 0) NOT LIKE 0)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1 || 1 NOT IN (0 BETWEEN 0 AND NULL DIV (2), 2 NOT BETWEEN 1 AND 0 & 1 <= 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1 || 1 NOT IN (0 BETWEEN 0 AND NULL DIV (2), 2 NOT BETWEEN 1 AND 0 & 1 <= 0)) <=> (((~ 1) OR (1 NOT IN ((0 BETWEEN 0 AND (NULL DIV 2)), ((2 NOT BETWEEN 1 AND (0 & 1)) <= 0)))))
< result
< column "(~ 1 || 1 NOT IN (0 BETWEEN 0 AND NULL DIV (2), 2 NOT BETWEEN 1 AND 0 & 1 <= 0)) <=> (((~ 1) OR (1 NOT IN ((0 BETWEEN 0 AND (NULL DIV 2)), ((2 NOT BETWEEN 1 AND (0 & 1)) <= 0)))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 & 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 & 0) <=> ((3 & 0))
< result
< column "(3 & 0) <=> ((3 & 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 MOD 3 BETWEEN 1 IS UNKNOWN AND ~ NULL NOT LIKE ! 1 <> NULL && 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 27 near "IS UNKNOWN AND ~ NULL NOT LIKE ! 1 <> NULL && 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 BETWEEN 2 AND 2 > 0 NOT BETWEEN NULL AND NULL % @v := NULL / 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 62 near ":= NULL / 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <=> NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 <=> NULL) <=> ((1 <=> NULL))
< result
< column "(1 <=> NULL) <=> ((1 <=> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 AND 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 AND 0) <=> ((3 AND 0))
< result
< column "(3 AND 0) <=> ((3 AND 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 1 - 2 | 2 >= ~ NULL << 2 NOT BETWEEN 3 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 1 - 2 | 2 >= ~ NULL << 2 NOT BETWEEN 3 AND 1) <=> ((NOT (((1 - 2) | 2) >= (((~ NULL) << 2) NOT BETWEEN 3 AND 1))))
< result
< column "(NOT 1 - 2 | 2 >= ~ NULL << 2 NOT BETWEEN 3 AND 1) <=> ((NOT (((1 - 2) | 2) >= (((~ NULL) << 2) NOT BETWEEN 3 AND 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 XOR 2 IS UNKNOWN MOD 1 > NULL || 2 <= 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "MOD 1 > NULL || 2 <= 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 > 1 = (3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 > 1 = (3)) <=> (((3 > 1) = 3))
< result
< column "(3 > 1 = (3)) <=> (((3 > 1) = 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 BETWEEN 2 AND 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 BETWEEN 2 AND 1) <=> ((2 BETWEEN 2 AND 1))
< result
< column "(2 BETWEEN 2 AND 1) <=> ((2 BETWEEN 2 AND 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL) <=> (NULL)
< result
< column "(NULL) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL NOT IN (0, 0) IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! NULL NOT IN (0, 0) IS NOT FALSE) <=> ((((NOT NULL) NOT IN (0, 0)) IS NOT FALSE))
< result
< column "(! NULL NOT IN (0, 0) IS NOT FALSE) <=> ((((NOT NULL) NOT IN (0, 0)) IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL BETWEEN 1 AND 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL BETWEEN 1 AND 3) <=> ((NULL BETWEEN 1 AND 3))
< result
< column "(NULL BETWEEN 1 AND 3) <=> ((NULL BETWEEN 1 AND 3))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NOT 2 > 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NOT 2 > 3) <=> ((NOT (2 > 3)))
< result
< column "(NOT 2 > 3) <=> ((NOT (2 > 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((0)) <=> (0)
< result
< column "((0)) <=> (0)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 || 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 || 0) <=> ((0 OR 0))
< result
< column "(0 || 0) <=> ((0 OR 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 OR 1 <= 3 <> 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 OR 1 <= 3 <> 3) <=> ((3 OR ((1 <= 3) <> 3)))
< result
< column "(3 OR 1 <= 3 <> 3) <=> ((3 OR ((1 <= 3) <> 3)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3) <=> (3)
< result
< column "(3) <=> (3)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((NULL)) <=> (NULL)
< result
< column "((NULL)) <=> (NULL)" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 1 || 1 IS UNKNOWN
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 1 || 1 IS UNKNOWN) <=> (((NOT 1) OR (1 IS UNKNOWN)))
< result
< column "(! 1 || 1 IS UNKNOWN) <=> (((NOT 1) OR (1 IS UNKNOWN)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (+ 3) > NULL IS NULL | 0 <=> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 29 near "| 0 <=> 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 * 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 * 0) <=> ((2 * 0))
< result
< column "(2 * 0) <=> ((2 * 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 1) <=> ((~ 1))
< result
< column "(~ 1) <=> ((~ 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! 3 NOT IN (2 AND 3, - NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (! 3 NOT IN (2 AND 3, - NULL)) <=> (((NOT 3) NOT IN ((2 AND 3), (- NULL))))
< result
< column "(! 3 NOT IN (2 AND 3, - NULL)) <=> (((NOT 3) NOT IN ((2 AND 3), (- NULL))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 IS NOT NULL >> 2 / 2 / NULL <= 2 <= NOT 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 23 near ">> 2 / 2 / NULL <= 2 <= NOT 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 2 MOD 2 <> 3 <> 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 2 MOD 2 <> 3 <> 0) <=> (((((+ 2) % 2) <> 3) <> 0))
< result
< column "(+ 2 MOD 2 <> 3 <> 0) <=> (((((+ 2) % 2) <> 3) <> 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > NULL && 0 | 0 BETWEEN ! 2 BETWEEN NULL AND 0 AND NULL NOT BETWEEN 0 AND 1 != - NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 47 near "BETWEEN NULL AND 0 AND NULL NOT BETWEEN 0 AND 1 != - NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 + 1) <=> ((0 + 1))
< result
< column "(0 + 1) <=> ((0 + 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 DIV 3 >> - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 DIV 3 >> - 1) <=> (((1 DIV 3) >> (- 1)))
< result
< column "(1 DIV 3 >> - 1) <=> (((1 DIV 3) >> (- 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT LIKE 1 ^ 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 24 near "^ 3" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 * 1 MOD ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 * 1 MOD ~ NULL) <=> (((0 * 1) % (~ NULL)))
< result
< column "(0 * 1 MOD ~ NULL) <=> (((0 * 1) % (~ NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 IS NOT FALSE) <=> ((1 IS NOT FALSE))
< result
< column "(1 IS NOT FALSE) <=> ((1 IS NOT FALSE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 XOR 2 + 0 IN (3, 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 XOR 2 + 0 IN (3, 3)) <=> ((3 XOR ((2 + 0) IN (3, 3))))
< result
< column "(3 XOR 2 + 0 IN (3, 3)) <=> ((3 XOR ((2 + 0) IN (3, 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT FALSE NOT IN (NULL IS NOT NULL, NULL != NULL) MOD NOT 3 - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 28 near "NOT IN (NULL IS NOT NULL, NULL != NULL) MOD NOT 3 - 1" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 % 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 % 0) <=> ((0 % 0))
< result
< column "(0 % 0) <=> ((0 % 0))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 3 <> NULL & 1 DIV 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- 3 <> NULL & 1 DIV 3) <=> (((- 3) <> (NULL & (1 DIV 3))))
< result
< column "(- 3 <> NULL & 1 DIV 3) <=> (((- 3) <> (NULL & (1 DIV 3))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 MOD 1 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 MOD 1 IS TRUE) <=> (((2 % 1) IS TRUE))
< result
< column "(2 MOD 1 IS TRUE) <=> (((2 % 1) IS TRUE))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - NULL OR 2 IS NOT FALSE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (- NULL OR 2 IS NOT FALSE) <=> (((- NULL) OR (2 IS NOT FALSE)))
< result
< column "(- NULL OR 2 IS NOT FALSE) <=> (((- NULL) OR (2 IS NOT FALSE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ NULL) <=> ((~ NULL))
< result
< column "(~ NULL) <=> ((~ NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ~ 0 NOT IN (0 < 3, 1 <=> 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (~ 0 NOT IN (0 < 3, 1 <=> 1)) <=> (((~ 0) NOT IN ((0 < 3), (1 <=> 1))))
< result
< column "(~ 0 NOT IN (0 < 3, 1 <=> 1)) <=> (((~ 0) NOT IN ((0 < 3), (1 <=> 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL NOT BETWEEN 0 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL NOT BETWEEN 0 AND NULL) <=> ((NULL NOT BETWEEN 0 AND NULL))
< result
< column "(NULL NOT BETWEEN 0 AND NULL) <=> ((NULL NOT BETWEEN 0 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT ! NULL & NULL IN (NULL, NULL) NOT IN (NULL, NOT 1 - 0)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 40 near "NOT IN (NULL, NOT 1 - 0)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT + 0 != 1 AND - 2 MOD 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (+ 0 != 1 AND - 2 MOD 1) <=> ((((+ 0) <> 1) AND ((- 2) % 1)))
< result
< column "(+ 0 != 1 AND - 2 MOD 1) <=> ((((+ 0) <> 1) AND ((- 2) % 1)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 << NULL || ~ 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 << NULL || ~ 2) <=> (((3 << NULL) OR (~ 2)))
< result
< column "(3 << NULL || ~ 2) <=> (((3 << NULL) OR (~ 2)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL && 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL && 2) <=> ((NULL AND 2))
< result
< column "(NULL && 2) <=> ((NULL AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 NOT IN (0, 0) <=> ~ 2 LIKE NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 NOT IN (0, 0) <=> ~ 2 LIKE NULL) <=> (((1 NOT IN (0, 0)) <=> ((~ 2) LIKE NULL)))
< result
< column "(1 NOT IN (0, 0) <=> ~ 2 LIKE NULL) <=> (((1 NOT IN (0, 0)) <=> ((~ 2) LIKE NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 <> 2 << + 3 BETWEEN + NULL <= 3 BETWEEN 0 AND 2 AND 2 >> 0 NOT IN (3 REGEXP 0, 1 IS NOT FALSE)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 38 near "<= 3 BETWEEN 0 AND 2 AND 2 >> 0 NOT IN (3 REGEXP 0, 1 IS NOT FALSE)" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 3 NOT IN (0, 3) <= NULL MOD 1 | ~ 0 / 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (3 NOT IN (0, 3) <= NULL MOD 1 | ~ 0 / 1) <=> (((3 NOT IN (0, 3)) <= ((NULL % 1) | ((~ 0) / 1))))
< result
< column "(3 NOT IN (0, 3) <= NULL MOD 1 | ~ 0 / 1) <=> (((3 NOT IN (0, 3)) <= ((NULL % 1) | ((~ 0) / 1))))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2 <=> NULL)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((2 <=> NULL)) <=> ((2 <=> NULL))
< result
< column "((2 <=> NULL)) <=> ((2 <=> NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 >= 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 >= 1) <=> ((1 >= 1))
< result
< column "(1 >= 1) <=> ((1 >= 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 0 AND 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 0 AND 2) <=> ((2 NOT BETWEEN 0 AND 2))
< result
< column "(2 NOT BETWEEN 0 AND 2) <=> ((2 NOT BETWEEN 0 AND 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL > 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL > 1) <=> ((NULL > 1))
< result
< column "(NULL > 1) <=> ((NULL > 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL IS NOT FALSE - NULL > 1 ^ NULL IS UNKNOWN IS NOT NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 26 near "- NULL > 1 ^ NULL IS UNKNOWN IS NOT NULL" 
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 | 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 | 1) <=> ((2 | 1))
< result
< column "(2 | 1) <=> ((2 | 1))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT - 2 MOD NULL REGEXP 2 NOT BETWEEN 2 AND NULL >= + 0
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …line 1 column 32 near "NOT BETWEEN 2 AND NULL >= + 0" 
> exec
| SET @omni_sweep = ?
: string
| SELECT NULL AND 0 IS TRUE
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (NULL AND 0 IS TRUE) <=> ((NULL AND (0 IS TRUE)))
< result
< column "(NULL AND 0 IS TRUE) <=> ((NULL AND (0 IS TRUE)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 0 && NULL >= 0 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (0 && NULL >= 0 || 2) <=> (((0 AND (NULL >= 0)) OR 2))
< result
< column "(0 && NULL >= 0 || 2) <=> (((0 AND (NULL >= 0)) OR 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 >= 0 DIV NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 >= 0 DIV NULL) <=> ((2 >= (0 DIV NULL)))
< result
< column "(2 >= 0 DIV NULL) <=> ((2 >= (0 DIV NULL)))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 1 / 0 * 3 REGEXP 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (1 / 0 * 3 REGEXP 2) <=> ((((1 / 0) * 3) REGEXP 2))
< result
< column "(1 / 0 * 3 REGEXP 2) <=> ((((1 / 0) * 3) REGEXP 2))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT 2 NOT BETWEEN 3 AND NULL
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT (2 NOT BETWEEN 3 AND NULL) <=> ((2 NOT BETWEEN 3 AND NULL))
< result
< column "(2 NOT BETWEEN 3 AND NULL) <=> ((2 NOT BETWEEN 3 AND NULL))" "BIGINT"
< row
: int64 1
> exec
| SET @omni_sweep = ?
: string
| SELECT (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> query
| SELECT ((2)) <=> (2)
< result
< column "((2)) <=> (2)" "BIGINT"
< row
: int64 1
