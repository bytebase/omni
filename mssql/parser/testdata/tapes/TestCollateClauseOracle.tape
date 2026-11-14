-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5
-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestCollateClauseOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE database_default FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS, y COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (x COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT ((x COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN) COLLATE SQL_Latin1_General_CP1_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT -(x COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT COALESCE(x COLLATE Latin1_General_CI_AS, y) COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CASE WHEN 1 = 1 THEN 'a' COLLATE Latin1_General_CI_AS END COLLATE Latin1_General_BIN
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS AT TIME ZONE 'UTC' COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS AT TIME ZONE 'UTC' COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x COLLATE Latin1_General_CI_AS = y COLLATE Latin1_General_BIN
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x FROM t ORDER BY (x COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b' COLLATE Latin1_General_BIN
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS * 2 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS - y COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS & 1 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT -x COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(x AS VARCHAR(10)) COLLATE Latin1_General_CI_AS + 'b' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (SELECT 'a') COLLATE Latin1_General_CI_AS + 'b'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS + 'b' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS AT TIME ZONE 'UTC' + 'x' COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x COLLATE Latin1_General_CI_AS + 'a' = y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x COLLATE Latin1_General_CI_AS LIKE 'a' COLLATE Latin1_General_BIN + '%'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t SET x = y COLLATE Latin1_General_CI_AS + 'z'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE Latin1_General_CI_AS + 'y')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN + 'b'
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1, COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE COLLATE Latin1_General_CI_AS = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_CI_AS FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN COLLATE SQL_Latin1_General_CP1_CI_AS FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS -- comment
|   COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT -x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' + 'b' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (x + y) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(x AS VARCHAR(10)) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CASE WHEN 1 = 1 THEN 'a' END COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (SELECT 'a') COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN = y
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x FROM t ORDER BY x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t SET x = y COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @v VARCHAR(10) = 'x' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS SPARSE NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS DEFAULT 'x' COLLATE Latin1_General_BIN)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) DEFAULT ('x' COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (x VARCHAR(10), a AS (x COLLATE Latin1_General_CI_AS) COLLATE Latin1_General_BIN)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @t TABLE (a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE DATABASE d2 COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL COLLATE Latin1_General_BIN)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) PRIMARY KEY COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) CHECK (a <> '') COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (x VARCHAR(10), a AS x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE t ADD c VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @t TABLE (a VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TYPE tt AS TABLE (a VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE [Latin1_General_CI_AS] FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE "Latin1_General_CI_AS" FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE [database_default] FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'database_default'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE t3 (a VARCHAR(10) COLLATE [Latin1_General_CI_AS])
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE [Latin1_General_CI_AS]
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE DATABASE d2 COLLATE [Latin1_General_CI_AS]
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb COLLATE [Latin1_General_CI_AS]
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollateClauseOracle]
< ok 0 -

== TestCollateClauseOracle scriptdom#1
> query
| SELECT 'a' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN + 'b'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":41,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1, COLLATE Latin1_General_CI_AS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":11,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 FROM t WHERE COLLATE Latin1_General_CI_AS = 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":39,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":39,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN COLLATE SQL_Latin1_General_CP1_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":39,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE Latin1_General_CI_AS -- comment
|   COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":2,"column":3,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT -x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":40,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 'a' + 'b' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":47,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT (x + y) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":45,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CAST(x AS VARCHAR(10)) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":60,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT CASE WHEN 1 = 1 THEN 'a' END COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":66,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT (SELECT 'a') COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":50,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":58,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 FROM t WHERE x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN = y
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":54,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x FROM t ORDER BY x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":57,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| UPDATE t SET x = y COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":49,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| DECLARE @v VARCHAR(10) = 'x' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":59,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":null,"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":61,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":73,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":41,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL COLLATE Latin1_General_BIN)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":70,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) PRIMARY KEY COLLATE Latin1_General_CI_AS)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":44,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) CHECK (a <> '') COLLATE Latin1_General_CI_AS)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":48,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (x VARCHAR(10), a AS x COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":69,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER TABLE t ADD c VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":null,"exprs":null}
> query
| DECLARE @t TABLE (a VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TYPE tt AS TABLE (a VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_BIN)
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":69,"message":"Incorrect syntax near \u0027COLLATE\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE [Latin1_General_CI_AS] FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":18,"message":"Incorrect syntax near \u0027[Latin1_General_CI_AS]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE "Latin1_General_CI_AS" FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":18,"message":"Incorrect syntax near \u0027\u0022Latin1_General_CI_AS\u0022\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE [database_default] FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":18,"message":"Incorrect syntax near \u0027[database_default]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE TABLE t3 (a VARCHAR(10) COLLATE [Latin1_General_CI_AS])
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":40,"message":"Incorrect syntax near \u0027[Latin1_General_CI_AS]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE [Latin1_General_CI_AS]
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":50,"message":"Incorrect syntax near \u0027[Latin1_General_CI_AS]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| CREATE DATABASE d2 COLLATE [Latin1_General_CI_AS]
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":28,"message":"Incorrect syntax near \u0027[Latin1_General_CI_AS]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| ALTER DATABASE testdb COLLATE [Latin1_General_CI_AS]
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":31,"message":"Incorrect syntax near \u0027[Latin1_General_CI_AS]\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
