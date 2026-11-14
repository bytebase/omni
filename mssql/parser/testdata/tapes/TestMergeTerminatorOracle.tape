-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestMergeTerminatorOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE; SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE /* c */ ;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
| ;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE -- c
| ;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OUTPUT $action;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OPTION (RECOMPILE);
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT 1 AS a) MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BEGIN MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE; END
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE; ELSE SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BEGIN TRY MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE; END TRY BEGIN CATCH SELECT 1 END CATCH
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WHILE 1 = 0 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE; MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN UPDATE SET a = 1;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE p AS MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE SELECT 1
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OUTPUT $action
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OPTION (RECOMPILE)
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT 1 AS a) MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BEGIN MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE END
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE ELSE SELECT 1
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BEGIN TRY MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE END TRY BEGIN CATCH SELECT 1 END CATCH
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WHILE 1 = 0 MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN UPDATE SET a = 1
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN NOT MATCHED THEN INSERT (a) VALUES (1)
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN NOT MATCHED THEN INSERT DEFAULT VALUES
|
< error
| 10713 1 15 1 A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE p AS MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE
|
< error
| 10713 1 15 1 proc=p A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE END
|
< error
| 10713 1 15 1 proc=p A MERGE statement must be terminated by a semi-colon (;).
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1; ELSE SELECT 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1 ELSE SELECT 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1;; ELSE SELECT 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1; ELSE SELECT 2;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 BEGIN SELECT 1 END; ELSE SELECT 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 IF 2 = 2 SELECT 1; ELSE SELECT 2; ELSE SELECT 3
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1; SELECT 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SELECT 1; SELECT 3; ELSE SELECT 2
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ELSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeTerminatorOracle]
< ok 0 -
