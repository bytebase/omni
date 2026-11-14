-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestKeywordOracleOptionPositions default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE INDEX ix ON dbo.t(a) WITH (FILLFACTOR = 80)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE INDEX ix ON dbo.t(a) WITH (SELECT = 1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE INDEX ix ON dbo.t(a) WITH (FROM = 1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE INDEX ix ON dbo.t(a) WITH (DROP = ON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'DROP'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET QUOTED_IDENTIFIER OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ARITHABORT ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET SELECT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SET'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET FROM OFF
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb SET RECOVERY SIMPLE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb SET READ_ONLY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb SET ANSI_NULLS ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb SET SELECT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER DATABASE testdb SET FROM OFF
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t WITH (ROWLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t WITH (HOLDLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t WITH (SELECT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t WITH (FROM)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t OPTION (RECOMPILE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t OPTION (MAXDOP 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t OPTION (SELECT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t OPTION (FROM)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR FAST_FORWARD FOR SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR SCROLL FOR SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR STATIC FOR SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR SELECT FOR SELECT 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE c CURSOR FROM FOR SELECT 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DBCC CHECKDB WITH (NO_INFOMSGS)
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DBCC CHECKDB WITH (ALL_ERRORMSGS)
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DBCC CHECKDB WITH (SELECT)
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DBCC CHECKDB WITH (FROM)
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH COMPRESSION
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH INIT
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH SELECT
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BACKUP DATABASE testdb TO DISK = '/tmp/test.bak' WITH FROM
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH NORECOVERY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH REPLACE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH SELECT
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| RESTORE DATABASE testdb FROM DISK = '/tmp/test.bak' WITH FROM
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR XML RAW, ELEMENTS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR XML RAW, ROOT('r')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR XML RAW, SELECT
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR XML RAW, FROM
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR JSON PATH, ROOT('r')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR JSON PATH, WITHOUT_ARRAY_WRAPPER
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR JSON PATH, SELECT
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 AS val FOR JSON PATH, FROM
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BULK INSERT dbo.t FROM '/tmp/data.csv' WITH (SELECT = 1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BULK INSERT dbo.t FROM '/tmp/data.csv' WITH (FROM = 1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_bad1 (a INT) WITH (SELECT = ON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_bad2 (a INT) WITH (FROM = ON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE dbo.p_recomp WITH RECOMPILE AS SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE dbo.p_enc WITH ENCRYPTION AS SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE dbo.p_bad1 WITH SELECT AS SELECT 1
|
< error
| 156 1 15 1 proc=p_bad1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE dbo.p_bad2 WITH FROM AS SELECT 1
|
< error
| 156 1 15 1 proc=p_bad2 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TRIGGER dbo.tr_enc ON dbo.t WITH ENCRYPTION FOR INSERT AS SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TRIGGER dbo.tr_bad1 ON dbo.t WITH SELECT FOR INSERT AS SELECT 1
|
< error
| 156 1 15 1 proc=tr_bad1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TRIGGER dbo.tr_bad2 ON dbo.t WITH FROM FOR INSERT AS SELECT 1
|
< error
| 156 1 15 1 proc=tr_bad2 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE VIEW dbo.v_sb WITH SCHEMABINDING AS SELECT 1 AS val
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE VIEW dbo.v_enc WITH ENCRYPTION AS SELECT 1 AS val
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE VIEW dbo.v_bad1 WITH SELECT AS SELECT 1 AS val
|
< error
| 156 1 15 1 proc=v_bad1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE VIEW dbo.v_bad2 WITH FROM AS SELECT 1 AS val
|
< error
| 156 1 15 1 proc=v_bad2 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FULLTEXT INDEX ON dbo.t (a) KEY INDEX ix WITH (SELECT = ON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE SERVICE svc ON QUEUE dbo.q (SELECT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE AVAILABILITY GROUP ag WITH (SELECT = ON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE ENDPOINT ep STATE = STARTED AS TCP (LISTENER_PORT = 5022) FOR DATABASE_MIRRORING (ROLE = PARTNER)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE ENDPOINT ep SELECT = STARTED AS TCP (LISTENER_PORT = 5022) FOR DATABASE_MIRRORING (ROLE = PARTNER)
|
< error
| 7801 4 16 1 The required parameter AS was not specified.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleOptionPositions]
< ok 0 -
