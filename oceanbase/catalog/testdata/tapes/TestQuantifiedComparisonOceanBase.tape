-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQuantifiedComparisonOceanBase default#1
> exec
| CREATE TABLE qt (a INT)
< ok 0 0
> exec
| CREATE TABLE qn (a INT)
< ok 0 0
> exec
| CREATE TABLE qe (a INT)
< ok 0 0
> exec
| INSERT INTO qt VALUES (0), (2)
< ok 2 0
> exec
| INSERT INTO qn VALUES (1), (NULL)
< ok 2 0
> exec
| SELECT 1 <=> ANY (SELECT a FROM qt)
< error
| 1064 42000 …near '(SELECT a FROM qt)' at line 1
> query
| SELECT 1 = ANY (SELECT a FROM qt) AS r1, 1 = ANY (SELECT a FROM qn) AS r2, 1 = ANY (SELECT a FROM qe) AS r3, NULL = ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| select (1 = any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 = any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 = any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL = any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| SELECT 1 = SOME (SELECT a FROM qt) AS r1, 1 = SOME (SELECT a FROM qn) AS r2, 1 = SOME (SELECT a FROM qe) AS r3, NULL = SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| select (1 = any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 = any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 = any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL = any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 1
: int64 0
: null
> query
| SELECT 1 = ALL (SELECT a FROM qt) AS r1, 1 = ALL (SELECT a FROM qn) AS r2, 1 = ALL (SELECT a FROM qe) AS r3, NULL = ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 = all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 = all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 = all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL = all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 <> ANY (SELECT a FROM qt) AS r1, 1 <> ANY (SELECT a FROM qn) AS r2, 1 <> ANY (SELECT a FROM qe) AS r3, NULL <> ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 <> SOME (SELECT a FROM qt) AS r1, 1 <> SOME (SELECT a FROM qn) AS r2, 1 <> SOME (SELECT a FROM qe) AS r3, NULL <> SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 <> ALL (SELECT a FROM qt) AS r1, 1 <> ALL (SELECT a FROM qn) AS r2, 1 <> ALL (SELECT a FROM qe) AS r3, NULL <> ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| select (1 != all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| SELECT 1 != ANY (SELECT a FROM qt) AS r1, 1 != ANY (SELECT a FROM qn) AS r2, 1 != ANY (SELECT a FROM qe) AS r3, NULL != ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 != SOME (SELECT a FROM qt) AS r1, 1 != SOME (SELECT a FROM qn) AS r2, 1 != SOME (SELECT a FROM qe) AS r3, NULL != SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 != ALL (SELECT a FROM qt) AS r1, 1 != ALL (SELECT a FROM qn) AS r2, 1 != ALL (SELECT a FROM qe) AS r3, NULL != ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| select (1 != all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 != all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 != all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL != all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 0
: int64 1
: null
> query
| SELECT 1 < ANY (SELECT a FROM qt) AS r1, 1 < ANY (SELECT a FROM qn) AS r2, 1 < ANY (SELECT a FROM qe) AS r3, NULL < ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 < any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 < any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 < any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL < any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 < SOME (SELECT a FROM qt) AS r1, 1 < SOME (SELECT a FROM qn) AS r2, 1 < SOME (SELECT a FROM qe) AS r3, NULL < SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 < any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 < any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 < any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL < any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 < ALL (SELECT a FROM qt) AS r1, 1 < ALL (SELECT a FROM qn) AS r2, 1 < ALL (SELECT a FROM qe) AS r3, NULL < ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| select (1 < all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 < all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 < all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL < all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| SELECT 1 <= ANY (SELECT a FROM qt) AS r1, 1 <= ANY (SELECT a FROM qn) AS r2, 1 <= ANY (SELECT a FROM qe) AS r3, NULL <= ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 <= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 <= any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 <= any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL <= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 <= SOME (SELECT a FROM qt) AS r1, 1 <= SOME (SELECT a FROM qn) AS r2, 1 <= SOME (SELECT a FROM qe) AS r3, NULL <= SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 <= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 <= any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 <= any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL <= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 <= ALL (SELECT a FROM qt) AS r1, 1 <= ALL (SELECT a FROM qn) AS r2, 1 <= ALL (SELECT a FROM qe) AS r3, NULL <= ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 <= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 <= all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 <= all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL <= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 > ANY (SELECT a FROM qt) AS r1, 1 > ANY (SELECT a FROM qn) AS r2, 1 > ANY (SELECT a FROM qe) AS r3, NULL > ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 > any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 > any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 > any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL > any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 > SOME (SELECT a FROM qt) AS r1, 1 > SOME (SELECT a FROM qn) AS r2, 1 > SOME (SELECT a FROM qe) AS r3, NULL > SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| select (1 > any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 > any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 > any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL > any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: null
: int64 0
: null
> query
| SELECT 1 > ALL (SELECT a FROM qt) AS r1, 1 > ALL (SELECT a FROM qn) AS r2, 1 > ALL (SELECT a FROM qe) AS r3, NULL > ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| select (1 > all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 > all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 > all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL > all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: int64 0
: int64 1
: null
> query
| SELECT 1 >= ANY (SELECT a FROM qt) AS r1, 1 >= ANY (SELECT a FROM qn) AS r2, 1 >= ANY (SELECT a FROM qe) AS r3, NULL >= ANY (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 >= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 >= any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 >= any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL >= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 >= SOME (SELECT a FROM qt) AS r1, 1 >= SOME (SELECT a FROM qn) AS r2, 1 >= SOME (SELECT a FROM qe) AS r3, NULL >= SOME (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| select (1 >= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 >= any (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 >= any (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL >= any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 1
: int64 1
: int64 0
: null
> query
| SELECT 1 >= ALL (SELECT a FROM qt) AS r1, 1 >= ALL (SELECT a FROM qn) AS r2, 1 >= ALL (SELECT a FROM qe) AS r3, NULL >= ALL (SELECT a FROM qt) AS r4
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| select (1 >= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r1`, (1 >= all (select `qn`.`a` from `omni_t_oceanbase_catalog_quantified`.`qn`)) as `r2`, (1 >= all (select `qe`.`a` from `omni_t_oceanbase_catalog_quantified`.`qe`)) as `r3`, (NULL >= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r4`
< result
< column "r1" "INT"
< column "r2" "INT"
< column "r3" "INT"
< column "r4" "INT"
< row
: int64 0
: null
: int64 1
: null
> query
| SELECT 1 > ALL (SELECT 0 UNION ALL SELECT 2) AS r1, 1 > ANY (SELECT 0 UNION ALL SELECT 2) AS r2
< result
< column "r1" "INT"
< column "r2" "INT"
< row
: int64 0
: int64 1
> query
| select (1 > all (select 0 as `0` union all select 2 as `2`)) as `r1`, (1 > any (select 0 as `0` union all select 2 as `2`)) as `r2`
< result
< column "r1" "INT"
< column "r2" "INT"
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
| select `q`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt` as `q` where (`q`.`a` >= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) order by `q`.`a`
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
| select `q`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt` as `q` where (`q`.`a` < any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) order by `q`.`a`
< result
< column "a" "INT"
< row
: int64 0
> query
| SELECT 1 = (2 > ANY (SELECT a FROM qt)) AS r
< result
< column "r" "INT"
< row
: int64 1
> query
| select (1 = (2 > any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`))) as `r`
< result
< column "r" "INT"
< row
: int64 1
> query
| SELECT NOT (1 = ANY (SELECT a FROM qt)) AS r
< result
< column "r" "INT"
< row
: int64 1
> query
| select (not (1 = any (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`))) as `r`
< result
< column "r" "INT"
< row
: int64 1
> query
| SELECT (1 + 1) <= ALL (SELECT a + 2 FROM qt) AS r
< result
< column "r" "INT"
< row
: int64 1
> query
| select ((1 + 1) <= all (select (`qt`.`a` + 2) as `?` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r`
< result
< column "r" "INT"
< row
: int64 1
> query
| SELECT MAX(q.a) >= ALL (SELECT a FROM qt) AS r FROM qt AS q
< result
< column "r" "INT"
< row
: int64 1
> query
| select (max(`q`.`a`) >= all (select `qt`.`a` from `omni_t_oceanbase_catalog_quantified`.`qt`)) as `r` from `omni_t_oceanbase_catalog_quantified`.`qt` as `q`
< result
< column "r" "INT"
< row
: int64 1
