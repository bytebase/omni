CREATE FUNCTION myfunc();
CREATE FUNCTION myfunc()
RETURNS string LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc()
RETURNS string LANGUAGE testlang;
CREATE PRIVATE FUNCTION myfunc()
RETURNS string LANGUAGE testlang;
CREATE FUNCTION myfunc()
RETURNS string VOLATILE LANGUAGE testlang;
CREATE FUNCTION deterministic()
RETURNS immutable STABLE LANGUAGE volatile;
CREATE FUNCTION IF NOT EXISTS myfunc()
RETURNS string LANGUAGE testlang;
CREATE OR REPLACE FUNCTION myfunc()
RETURNS string LANGUAGE testlang;
CREATE TEMP FUNCTION IF NOT EXISTS myfunc()
RETURNS string LANGUAGE testlang;
CREATE TEMP FUNCTION IF NOT EXISTS myfunc()
RETURNS string LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(param_a int32)
RETURNS STRUCT< x string, y boolean > LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(param_a int32)
RETURNS a.b.c LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(a int32, b STRUCT< x string, y int32 >, c ARRAY< boolean >)
RETURNS string LANGUAGE testlang;
CREATE FUNCTION fn(s string)
RETURNS string LANGUAGE testlang AS """ return
  "presto!" + s + '\n';
""";
CREATE FUNCTION fn(s string)
RETURNS string LANGUAGE testlang OPTIONS
  (a = b, bruce = lee) AS "return 'a';";
CREATE FUNCTION fn(s string)
RETURNS string LANGUAGE testlang OPTIONS
  (a = b, bruce = lee) AS "return 'a';";
CREATE FUNCTION fn(s string)
RETURNS string;
CREATE FUNCTION fn(string)
RETURNS string LANGUAGE testlang AS "return 'a';";
CREATE FUNCTION fn(string, s string, int32, i int32)
RETURNS string LANGUAGE testlang AS "return 'a'";
CREATE FUNCTION mypackage.myfunc(int32 AS alias)
RETURNS string LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(int32 alias)
RETURNS string LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(x int32 AS alias)
RETURNS string LANGUAGE testlang;
CREATE FUNCTION mypackage.myfunc(x int32 AS alias NOT AGGREGATE)
RETURNS string LANGUAGE testlang;
CREATE FUNCTION sql_func(x int32 AS type_alias, y type_alias)
RETURNS int32 AS (
  (
    SELECT
      CAST(1 AS type_alias) + x - y
  )
);
CREATE FUNCTION myfunc()
RETURNS interval;
CREATE FUNCTION myfunc(a int64, b string DEFAULT "abc")
RETURNS double;
CREATE FUNCTION myfunc(a int64 DEFAULT -314, b string DEFAULT "abc", c double DEFAULT 3.14)
RETURNS double;
CREATE FUNCTION myfunc(a int64, b ANY TYPE DEFAULT "abc", c ANY TYPE DEFAULT 3.14)
RETURNS double;
CREATE FUNCTION myfunc()
RETURNS string REMOTE;
CREATE FUNCTION myfunc()
RETURNS string REMOTE WITH CONNECTION myconn;
CREATE FUNCTION myfunc()
RETURNS string REMOTE WITH CONNECTION DEFAULT;
CREATE FUNCTION myfunc()
RETURNS string REMOTE WITH CONNECTION `myconn-abc`;
CREATE FUNCTION myfunc()
RETURNS string REMOTE WITH CONNECTION myconn OPTIONS
  (a = b, bruce = lee);
CREATE FUNCTION myfunc()
RETURNS string VOLATILE REMOTE;
CREATE FUNCTION myfunc()
RETURNS int LANGUAGE js REMOTE;
CREATE FUNCTION myfunc()
RETURNS int LANGUAGE js REMOTE WITH CONNECTION c;
CREATE FUNCTION myfunc()
RETURNS int REMOTE AS (
  5
);
CREATE FUNCTION remote(int remote)
RETURNS int REMOTE WITH CONNECTION remote;
CREATE FUNCTION returns_remote()
RETURNS remote AS (
  NULL
);
CREATE FUNCTION returns_remote()
RETURNS remote REMOTE;
CREATE FUNCTION my_func(x FUNCTION<() -> STRUCT< INT64 > >)
AS (
  x()
);
CREATE FUNCTION my_func(x FUNCTION<(STRUCT< INT64 >) -> INT64 >)
AS (
  x(1)
);
CREATE FUNCTION my_func(x FUNCTION<(STRUCT< INT64 >) -> INT64 >)
AS (
  x(1)
);
CREATE FUNCTION my_func(x FUNCTION<(STRUCT< INT64 >, STRUCT< INT64 >) -> INT64 >)
AS (
  x(1, 2)
);
CREATE FUNCTION `function`(`function` FUNCTION<() -> INT64 >)
AS (
  `function`()
);
CREATE FUNCTION func(x FUNCTION<() -> INT64 > (0))
AS (
  1
);
