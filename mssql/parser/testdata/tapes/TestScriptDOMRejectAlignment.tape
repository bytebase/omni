-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestScriptDOMRejectAlignment scriptdom#1
> query
| SELECT * FROM (SELECT 1 AS a) AS t()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":35,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"QueryDerivedTable","alias":"t"}]}}]},"exprs":["*"]}
> query
| WITH c() AS (SELECT 1) SELECT * FROM c
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":8,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE VIEW v() AS SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":15,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE a IN ()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":28,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| INSERT INTO t VALUES ()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t ()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":16,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t PIVOT (SUM(x) FOR y IN ()) p
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":41,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t UNPIVOT (v FOR c IN ()) u
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) FROM t GROUP BY ROLLUP()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":40,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) FROM t GROUP BY CUBE()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT *, ROW_NUMBER() OVER (PARTITION BY) FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE INDEX ix ON t()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":22,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM (SELECT 1 AS a) AS t(a,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WITH c(a,) AS (SELECT 1) SELECT * FROM c
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":10,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WITH c AS (SELECT 1), SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE VIEW v(a,) AS SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| INSERT INTO t (a,) VALUES (1)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":18,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| INSERT INTO t VALUES (1,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE a IN (1,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT f(1,) FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":12,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM dbo.fn(1,) AS t(a)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t ORDER BY a,
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":28,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) FROM t GROUP BY ROLLUP(a,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT *, ROW_NUMBER() OVER (PARTITION BY a,) FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":45,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t PIVOT (SUM(x) FOR y IN ([a],)) p
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":45,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t UNPIVOT (v FOR c IN (a,)) u
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":40,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE INDEX ix ON t(a,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a,b,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":50,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (a INT, b INT, CONSTRAINT pk PRIMARY KEY (a,b,))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":62,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WITH (NOLOCK,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":30,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) FROM t GROUP BY GROUPING SETS (, a)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":48,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WITH XMLNAMESPACES ('foo' AS p, , 'bar' AS q) SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":33,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WITH (INDEX(123,,456))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":33,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (a INT, 123, b INT)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027123\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON(,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON(@j,)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":27,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON()
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t OPTION (TABLE HINT(t, NOLOCK,))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":46,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t OPTION (TABLE HINT(t,))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Incorrect syntax near \u0027)\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t OPTION (TABLE HINT(t,,NOLOCK))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":8,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 AS SELECT
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 AS WHERE
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027WHERE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 AS FROM
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE DATABASE SELECT
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE DATABASE FROM
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM FROM
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":15,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (SELECT INT)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":16,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET RECOVERY OFF
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027RECOVERY\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET RECOVERY BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":33,"message":"Incorrect syntax near \u0027BOGUS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET RECOVERY SELECT
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027RECOVERY\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET PAGE_VERIFY OFF
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027PAGE_VERIFY\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET PARAMETERIZATION SELECT
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027PARAMETERIZATION\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET CURSOR_DEFAULT OFF
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027CURSOR_DEFAULT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE mydb FAILOVER
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":21,"message":"Incorrect syntax near \u0027FAILOVER\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE mydb UNKNOWN_ACTION
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":21,"message":"Incorrect syntax near \u0027UNKNOWN_ACTION\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER PARTITION FUNCTION pf() SELECT RANGE (100)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET ENCRYPTION SELECT
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027ENCRYPTION\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET ENCRYPTION BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027ENCRYPTION\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET AUTO_CLOSE BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027AUTO_CLOSE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE foo SET QUOTED_IDENTIFIER BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Incorrect syntax near \u0027QUOTED_IDENTIFIER\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE ASSEMBLY a FROM 'x' WITH PERMISSION_SET = BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":50,"message":"Incorrect syntax near \u0027BOGUS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t FOR XML BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Incorrect syntax near \u0027BOGUS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t FOR JSON BOGUS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":26,"message":"Incorrect syntax near \u0027BOGUS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CONTAINS(b, 'foo')
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":1,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT FREETEXT(b, 'foo')
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":1,"message":"Incorrect syntax near \u0027SELECT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT c FROM t WHERE b > (SELECT CONTAINS(col, 'x') FROM u)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":35,"message":"Incorrect syntax near \u0027CONTAINS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COALESCE(CONTAINS(b, 'foo'), 0)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Incorrect syntax near \u0027CONTAINS\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
