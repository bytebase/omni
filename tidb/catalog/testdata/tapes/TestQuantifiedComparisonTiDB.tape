-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQuantifiedComparisonTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_quantified
< ok 0 0
> exec
| CREATE DATABASE omni_quantified
< ok 0 0
> exec
| USE omni_quantified
< ok 0 0
> exec
| CREATE TABLE qt (a INT);
| CREATE TABLE qn (a INT);
| CREATE TABLE qe (a INT);
| INSERT INTO qt VALUES (0), (2);
| INSERT INTO qn VALUES (1), (NULL)
< ok 2 0
> query
| SELECT 1 = ANY (SELECT a FROM qt) AS r1, 1 = ANY (SELECT a FROM qn) AS r2, 1 = ANY (SELECT a FROM qe) AS r3, NULL = ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| select (1 = any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 = any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 = any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL = any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| SELECT 1 = SOME (SELECT a FROM qt) AS r1, 1 = SOME (SELECT a FROM qn) AS r2, 1 = SOME (SELECT a FROM qe) AS r3, NULL = SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| select (1 = any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 = any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 = any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL = any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| SELECT 1 = ALL (SELECT a FROM qt) AS r1, 1 = ALL (SELECT a FROM qn) AS r2, 1 = ALL (SELECT a FROM qe) AS r3, NULL = ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 = all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 = all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 = all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL = all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 <> ANY (SELECT a FROM qt) AS r1, 1 <> ANY (SELECT a FROM qn) AS r2, 1 <> ANY (SELECT a FROM qe) AS r3, NULL <> ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 <> SOME (SELECT a FROM qt) AS r1, 1 <> SOME (SELECT a FROM qn) AS r2, 1 <> SOME (SELECT a FROM qe) AS r3, NULL <> SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 <> ALL (SELECT a FROM qt) AS r1, 1 <> ALL (SELECT a FROM qn) AS r2, 1 <> ALL (SELECT a FROM qe) AS r3, NULL <> ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| select (1 != all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| SELECT 1 != ANY (SELECT a FROM qt) AS r1, 1 != ANY (SELECT a FROM qn) AS r2, 1 != ANY (SELECT a FROM qe) AS r3, NULL != ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 != SOME (SELECT a FROM qt) AS r1, 1 != SOME (SELECT a FROM qn) AS r2, 1 != SOME (SELECT a FROM qe) AS r3, NULL != SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 != ALL (SELECT a FROM qt) AS r1, 1 != ALL (SELECT a FROM qn) AS r2, 1 != ALL (SELECT a FROM qe) AS r3, NULL != ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| select (1 != all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 != all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 != all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL != all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "TINYINT"
< column "r2" "TINYINT"
< column "r3" "TINYINT"
< column "r4" "TINYINT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| SELECT 1 < ANY (SELECT a FROM qt) AS r1, 1 < ANY (SELECT a FROM qn) AS r2, 1 < ANY (SELECT a FROM qe) AS r3, NULL < ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 < any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 < any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 < any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL < any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 < SOME (SELECT a FROM qt) AS r1, 1 < SOME (SELECT a FROM qn) AS r2, 1 < SOME (SELECT a FROM qe) AS r3, NULL < SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 < any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 < any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 < any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL < any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 < ALL (SELECT a FROM qt) AS r1, 1 < ALL (SELECT a FROM qn) AS r2, 1 < ALL (SELECT a FROM qe) AS r3, NULL < ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| select (1 < all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 < all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 < all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL < all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| SELECT 1 <= ANY (SELECT a FROM qt) AS r1, 1 <= ANY (SELECT a FROM qn) AS r2, 1 <= ANY (SELECT a FROM qe) AS r3, NULL <= ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 <= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <= any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <= any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 <= SOME (SELECT a FROM qt) AS r1, 1 <= SOME (SELECT a FROM qn) AS r2, 1 <= SOME (SELECT a FROM qe) AS r3, NULL <= SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 <= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <= any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <= any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 <= ALL (SELECT a FROM qt) AS r1, 1 <= ALL (SELECT a FROM qn) AS r2, 1 <= ALL (SELECT a FROM qe) AS r3, NULL <= ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 <= all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <= all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <= all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <= all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 > ANY (SELECT a FROM qt) AS r1, 1 > ANY (SELECT a FROM qn) AS r2, 1 > ANY (SELECT a FROM qe) AS r3, NULL > ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 > any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 > any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 > any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL > any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 > SOME (SELECT a FROM qt) AS r1, 1 > SOME (SELECT a FROM qn) AS r2, 1 > SOME (SELECT a FROM qe) AS r3, NULL > SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 > any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 > any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 > any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL > any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 > ALL (SELECT a FROM qt) AS r1, 1 > ALL (SELECT a FROM qn) AS r2, 1 > ALL (SELECT a FROM qe) AS r3, NULL > ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| select (1 > all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 > all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 > all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL > all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| SELECT 1 >= ANY (SELECT a FROM qt) AS r1, 1 >= ANY (SELECT a FROM qn) AS r2, 1 >= ANY (SELECT a FROM qe) AS r3, NULL >= ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 >= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 >= any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 >= any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL >= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 >= SOME (SELECT a FROM qt) AS r1, 1 >= SOME (SELECT a FROM qn) AS r2, 1 >= SOME (SELECT a FROM qe) AS r3, NULL >= SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 >= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 >= any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 >= any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL >= any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 >= ALL (SELECT a FROM qt) AS r1, 1 >= ALL (SELECT a FROM qn) AS r2, 1 >= ALL (SELECT a FROM qe) AS r3, NULL >= ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 >= all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 >= all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 >= all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL >= all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< column "r3" "BIGINT"
< column "r4" "BIGINT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 <=> ANY (SELECT a FROM qt) AS r1, 1 <=> ANY (SELECT a FROM qn) AS r2, 1 <=> ANY (SELECT a FROM qe) AS r3, NULL <=> ANY (SELECT a FROM qt) AS r4
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| select (1 <=> any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <=> any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <=> any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <=> any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| SELECT 1 <=> SOME (SELECT a FROM qt) AS r1, 1 <=> SOME (SELECT a FROM qn) AS r2, 1 <=> SOME (SELECT a FROM qe) AS r3, NULL <=> SOME (SELECT a FROM qt) AS r4
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| select (1 <=> any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <=> any (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <=> any (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <=> any (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| SELECT 1 <=> ALL (SELECT a FROM qt) AS r1, 1 <=> ALL (SELECT a FROM qn) AS r2, 1 <=> ALL (SELECT a FROM qe) AS r3, NULL <=> ALL (SELECT a FROM qt) AS r4
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| select (1 <=> all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r1`, (1 <=> all (select `qn`.`a` from `omni_quantified`.`qn`)) as `r2`, (1 <=> all (select `qe`.`a` from `omni_quantified`.`qe`)) as `r3`, (NULL <=> all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r4`
< error
| 1105 HY000 We don't support <=> all or <=> any now
> query
| SELECT 1 > ALL (SELECT 0 UNION ALL SELECT 2) AS r1, 1 > ANY (SELECT 0 UNION ALL SELECT 2) AS r2
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< row
: int64 0
: int64 1
> query
| select (1 > all (select 0 as `0` union all select 2 as `2`)) as `r1`, (1 > any (select 0 as `0` union all select 2 as `2`)) as `r2`
< result
< column "r1" "BIGINT"
< column "r2" "BIGINT"
< row
: int64 0
: int64 1
> query
| SELECT q.a FROM qt AS q WHERE q.a >= ALL (SELECT a FROM qt) ORDER BY q.a
< result
< column "a" "INT"
< row
: int64 2
> query
| select `q`.`a` from `omni_quantified`.`qt` as `q` where (`q`.`a` >= all (select `qt`.`a` from `omni_quantified`.`qt`)) order by `q`.`a`
< result
< column "a" "INT"
< row
: int64 2
> query
| SELECT q.a FROM qt AS q WHERE q.a < SOME (SELECT a FROM qt) ORDER BY q.a
< result
< column "a" "INT"
< row
: int64 0
> query
| select `q`.`a` from `omni_quantified`.`qt` as `q` where (`q`.`a` < any (select `qt`.`a` from `omni_quantified`.`qt`)) order by `q`.`a`
< result
< column "a" "INT"
< row
: int64 0
> query
| SELECT 1 = (2 > ANY (SELECT a FROM qt)) AS r
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| select (1 = (2 > any (select `qt`.`a` from `omni_quantified`.`qt`))) as `r`
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| SELECT NOT (1 = ANY (SELECT a FROM qt)) AS r
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| select (not (1 = any (select `qt`.`a` from `omni_quantified`.`qt`))) as `r`
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| SELECT (1 + 1) <= ALL (SELECT a + 2 FROM qt) AS r
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| select ((1 + 1) <= all (select (`qt`.`a` + 2) as `?` from `omni_quantified`.`qt`)) as `r`
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| SELECT MAX(q.a) >= ALL (SELECT a FROM qt) AS r FROM qt AS q
< result
< column "r" "BIGINT"
< row
: int64 1
> query
| select (max(`q`.`a`) >= all (select `qt`.`a` from `omni_quantified`.`qt`)) as `r` from `omni_quantified`.`qt` as `q`
< result
< column "r" "BIGINT"
< row
: int64 1
> exec
| DROP DATABASE IF EXISTS omni_quantified
< ok 0 0
