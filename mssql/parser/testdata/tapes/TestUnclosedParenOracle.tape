-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5
-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestUnclosedParenOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT ((((1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT ((((1))))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 + (2
|
< error
| 102 1 15 1 Incorrect syntax near '2'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 + (2)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE (a = 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE (a = 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t SET a = (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t SET a = (1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (SELECT 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (SELECT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE a IN (SELECT a FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 't'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE a IN (SELECT a FROM t)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE a > ANY (SELECT a FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 't'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE a > ANY (SELECT a FROM t)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE EXISTS (SELECT 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE EXISTS (SELECT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(1 AS int
|
< error
| 102 1 15 1 Incorrect syntax near 'int'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(1 AS int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT TRY_CAST(1 AS int
|
< error
| 102 1 15 1 Incorrect syntax near 'int'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT TRY_CAST(1 AS int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CONVERT(int, 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CONVERT(int, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CONVERT(varchar(10), 1, 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CONVERT(varchar(10), 1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT TRY_CONVERT(int, 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT TRY_CONVERT(int, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COALESCE(1, 2
|
< error
| 102 1 15 1 Incorrect syntax near '2'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COALESCE(1, 2)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT NULLIF(1, 2
|
< error
| 102 1 15 1 Incorrect syntax near '2'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT NULLIF(1, 2)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT IIF(1 = 1, 1, 0
|
< error
| 102 1 15 1 Incorrect syntax near '0'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT IIF(1 = 1, 1, 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(*
|
< error
| 102 1 15 1 Incorrect syntax near '*'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(*)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(*) OVER (PARTITION BY a
|
< error
| 102 1 15 1 Incorrect syntax near 'a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(*) OVER (PARTITION BY a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT SUM(a) OVER (ORDER BY a ROWS UNBOUNDED PRECEDING
|
< error
| 102 1 15 1 Incorrect syntax near 'PRECEDING'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT SUM(a) OVER (ORDER BY a ROWS UNBOUNDED PRECEDING)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT STRING_AGG(a, ',') WITHIN GROUP (ORDER BY a
|
< error
| 102 1 15 1 Incorrect syntax near 'a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT STRING_AGG(a, ',') WITHIN GROUP (ORDER BY a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a int
|
< error
| 102 1 15 1 Incorrect syntax near 'int'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10
|
< error
| 102 1 15 1 Incorrect syntax near '10'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR XML RAW, ROOT('r'
|
< error
| 102 1 15 1 Incorrect syntax near 'r'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR XML RAW, ROOT('r')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR XML AUTO, XMLSCHEMA('urn:x'
|
< error
| 102 1 15 1 Incorrect syntax near 'urn:x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR XML AUTO, XMLSCHEMA('urn:x')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR JSON PATH, ROOT('r'
|
< error
| 102 1 15 1 Incorrect syntax near 'r'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t FOR JSON PATH, ROOT('r')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @x int = (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @x int = (1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @t TABLE (a int
|
< error
| 102 1 15 1 Incorrect syntax near 'int'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @t TABLE (a int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR FOR SELECT (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR FOR SELECT (1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE p AS RETURN (1
|
< error
| 102 1 15 1 proc=p Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE p AS RETURN (1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR (RECEIVE * FROM q
|
< error
| 102 1 15 1 Incorrect syntax near 'q'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR (RECEIVE * FROM q)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RECEIVE TOP (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RECEIVE TOP (1) * FROM q
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| EXEC ('SELECT 1'
|
< error
| 102 1 15 1 Incorrect syntax near 'SELECT 1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| EXEC ('SELECT 1')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO t EXEC ('SELECT 1'
|
< error
| 102 1 15 1 Incorrect syntax near 'SELECT 1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO t EXEC ('SELECT 1')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RAISERROR('x', 16, 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RAISERROR('x', 16, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| COMMIT TRAN WITH (DELAYED_DURABILITY = ON
|
< error
| 102 1 15 1 Incorrect syntax near 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| COMMIT TRAN WITH (DELAYED_DURABILITY = ON)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (VALUES (1, 2, (3, 4)) AS t(a, b)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t GROUP BY GROUPING SETS ((a, (a, b), ())
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM t GROUP BY GROUPING SETS ((a), (a, b), ())
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(* FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COUNT(*) FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(1 AS decimal((10, 2))
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(1 AS decimal(10, 2))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT geography::Point((1, 2, 4326)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT geography::Point(1, 2, 4326)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10 '$.a')
|
< error
| 102 1 15 1 Incorrect syntax near '$.a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10) '$.a')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR (RECEIVE TOP (1 * FROM q)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR (RECEIVE TOP (1) * FROM q)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int IDENTITY(1, 1, b int)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int IDENTITY(1, 1), b int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int CHECK (a > 0, b int)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int CHECK (a > 0), b int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int, CHECK (a > 0, b int)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t (a int, CHECK (a > 0), b int)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET TEXTSIZE (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| FETCH ABSOLUTE (1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR DELAY ('00:00:01'
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WAITFOR (RECEIVE * FROM q), TIMEOUT (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| EXEC p @a = (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO t EXEC p @a = (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RAISERROR('x', 16, (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| THROW 50000, 'x', (1
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT dbo.f(1) OVER (ORDER BY a
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ORDER'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestUnclosedParenOracle]
< ok 0 -

== TestUnclosedParenOracle scriptdom#1
> query
| SELECT ((((1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT ((((1))))
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ParenthesisExpression"}]}}]},"exprs":["(paren (paren (paren (paren #))))"]}
> query
| SELECT (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":10,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT (1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ParenthesisExpression"}]}}]},"exprs":["(paren #)"]}
> query
| SELECT 1 + (2
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":14,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 + (2)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}]}}]},"exprs":["(\u002B # (paren #))"]}
> query
| SELECT * FROM t WHERE (a = 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":29,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE (a = 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["*"]}
> query
| UPDATE t SET a = (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":20,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| UPDATE t SET a = (1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"UpdateStatement"}]},"exprs":null}
> query
| SELECT (SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT (SELECT 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ScalarSubquery"}]}}]},"exprs":["#"]}
> query
| SELECT * FROM t WHERE a IN (SELECT a FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":44,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE a IN (SELECT a FROM t)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM t WHERE a > ANY (SELECT a FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":47,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE a > ANY (SELECT a FROM t)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM t WHERE EXISTS (SELECT 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":39,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM t WHERE EXISTS (SELECT 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["*"]}
> query
| SELECT CAST(1 AS int
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":21,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CAST(1 AS int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"CastCall"}]}}]},"exprs":["#"]}
> query
| SELECT TRY_CAST(1 AS int
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT TRY_CAST(1 AS int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"TryCastCall"}]}}]},"exprs":["#"]}
> query
| SELECT CONVERT(int, 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":22,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CONVERT(int, 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ConvertCall"}]}}]},"exprs":["#"]}
> query
| SELECT CONVERT(varchar(10), 1, 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":33,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CONVERT(varchar(10), 1, 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ConvertCall"}]}}]},"exprs":["#"]}
> query
| SELECT TRY_CONVERT(int, 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":26,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT TRY_CONVERT(int, 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"TryConvertCall"}]}}]},"exprs":["#"]}
> query
| SELECT COALESCE(1, 2
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":21,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COALESCE(1, 2)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"CoalesceExpression"}]}}]},"exprs":["#"]}
> query
| SELECT NULLIF(1, 2
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":19,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT NULLIF(1, 2)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"NullIfExpression"}]}}]},"exprs":["#"]}
> query
| SELECT IIF(1 = 1, 1, 0
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT IIF(1 = 1, 1, 0)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"IIfCall"}]}}]},"exprs":["#"]}
> query
| SELECT COUNT(*
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":15,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}]}}]},"exprs":["#"]}
> query
| SELECT COUNT(*) OVER (PARTITION BY a
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":37,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) OVER (PARTITION BY a)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}]}}]},"exprs":["#"]}
> query
| SELECT SUM(a) OVER (ORDER BY a ROWS UNBOUNDED PRECEDING
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":56,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT SUM(a) OVER (ORDER BY a ROWS UNBOUNDED PRECEDING)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}]}}]},"exprs":["#"]}
> query
| SELECT STRING_AGG(a, ',') WITHIN GROUP (ORDER BY a
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":51,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT STRING_AGG(a, ',') WITHIN GROUP (ORDER BY a)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}]}}]},"exprs":["#"]}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a int
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"OpenJsonTableReference","with_cols":[{"name":"a","type":"int"}]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":49,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10))
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"OpenJsonTableReference","with_cols":[{"name":"a","type":"varchar"}]}]}}]},"exprs":["*"]}
> query
| SELECT a FROM t FOR XML RAW, ROOT('r'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":38,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT a FROM t FOR XML RAW, ROOT('r')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT a FROM t FOR XML AUTO, XMLSCHEMA('urn:x'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":48,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT a FROM t FOR XML AUTO, XMLSCHEMA('urn:x')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT a FROM t FOR JSON PATH, ROOT('r'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":40,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT a FROM t FOR JSON PATH, ROOT('r')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| DECLARE @x int = (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":20,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| DECLARE @x int = (1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"DeclareVariableStatement"}]},"exprs":null}
> query
| DECLARE @t TABLE (a int
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":24,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| DECLARE @t TABLE (a int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"DeclareTableVariableStatement"}]},"exprs":null}
> query
| DECLARE c CURSOR FOR SELECT (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| DECLARE c CURSOR FOR SELECT (1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"DeclareCursorStatement"}]},"exprs":null}
> query
| CREATE PROCEDURE p AS RETURN (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Incorrect syntax near \u00271\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE PROCEDURE p AS RETURN (1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"CreateProcedureStatement"}]},"exprs":null}
> query
| WAITFOR (RECEIVE * FROM q
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":26,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WAITFOR (RECEIVE * FROM q)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"WaitForStatement"}]},"exprs":null}
> query
| RECEIVE TOP (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":15,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| RECEIVE TOP (1) * FROM q
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"ReceiveStatement"}]},"exprs":null}
> query
| EXEC ('SELECT 1'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| EXEC ('SELECT 1')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"ExecuteStatement"}]},"exprs":null}
> query
| INSERT INTO t EXEC ('SELECT 1'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| INSERT INTO t EXEC ('SELECT 1')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Insert"}]},"exprs":null}
> query
| RAISERROR('x', 16, 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":21,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| RAISERROR('x', 16, 1)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"RaiseErrorStatement"}]},"exprs":null}
> query
| COMMIT TRAN WITH (DELAYED_DURABILITY = ON
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Unexpected end of file occurred."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| COMMIT TRAN WITH (DELAYED_DURABILITY = ON)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"CommitTransactionStatement"}]},"exprs":null}
> query
| SELECT * FROM (VALUES (1, 2, (3, 4)) AS t(a, b)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":32,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"InlineDerivedTable","alias":"t","columns":["a","b"]}]}}]},"exprs":["*"]}
> query
| SELECT a FROM t GROUP BY GROUPING SETS ((a, (a, b), ())
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":47,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT a FROM t GROUP BY GROUPING SETS ((a), (a, b), ())
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT COUNT(* FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":16,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT COUNT(*) FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT CAST(1 AS decimal((10, 2))
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CAST(1 AS decimal(10, 2))
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"CastCall"}]}}]},"exprs":["#"]}
> query
| SELECT geography::Point((1, 2, 4326)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":27,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT geography::Point(1, 2, 4326)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"FunctionCall"}]}}]},"exprs":["#"]}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10 '$.a')
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":50,"message":"Incorrect syntax near \u0027\u0027$.a\u0027\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT * FROM OPENJSON(N'[]') WITH (a varchar(10) '$.a')
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"OpenJsonTableReference","with_cols":[{"name":"a","type":"varchar"}]}]}}]},"exprs":["*"]}
> query
| WAITFOR (RECEIVE TOP (1 * FROM q)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":27,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WAITFOR (RECEIVE TOP (1) * FROM q)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"WaitForStatement"}]},"exprs":null}
> query
| CREATE TABLE t (a int IDENTITY(1, 1, b int)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":36,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (a int IDENTITY(1, 1), b int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"CreateTableStatement"}]},"exprs":null}
> query
| CREATE TABLE t (a int CHECK (a > 0, b int)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":35,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (a int CHECK (a > 0), b int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"CreateTableStatement"}]},"exprs":null}
> query
| CREATE TABLE t (a int, CHECK (a > 0, b int)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":36,"message":"Incorrect syntax near \u0027,\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t (a int, CHECK (a > 0), b int)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"CreateTableStatement"}]},"exprs":null}
> query
| SET TEXTSIZE (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":14,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| FETCH ABSOLUTE (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":17,"message":"Incorrect syntax near \u00271\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"FetchCursorStatement"}]},"exprs":null}
> query
| WAITFOR DELAY ('00:00:01'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":15,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| WAITFOR (RECEIVE * FROM q), TIMEOUT (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":37,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| EXEC p @a = (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| INSERT INTO t EXEC p @a = (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":27,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| RAISERROR('x', 16, (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":20,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| THROW 50000, 'x', (1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":19,"message":"Incorrect syntax near \u0027(\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT dbo.f(1) OVER (ORDER BY a
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Incorrect syntax near \u0027ORDER\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
