-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestSetOptionListOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_PADDING, ANSI_WARNINGS, CONCAT_NULL_YIELDS_NULL, ARITHABORT, QUOTED_IDENTIFIER, ANSI_NULLS ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NUMERIC_ROUNDABORT, ARITHIGNORE, CURSOR_CLOSE_ON_COMMIT, IMPLICIT_TRANSACTIONS, ANSI_NULL_DFLT_ON OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULL_DFLT_OFF, FORCEPLAN, REMOTE_PROC_TRANSACTIONS, NO_BROWSETABLE OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_DEFAULTS, FMTONLY, NOEXEC OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NOCOUNT, XACT_ABORT ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET PARSEONLY, NOCOUNT ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, ANSI_NULLS ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| set ansi_nulls /* c */ , quoted_identifier on
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER OFF SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| IF 1 = 1 SET ANSI_NULLS, ANSI_WARNINGS ON ELSE SET ANSI_NULLS, ANSI_WARNINGS OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE PROCEDURE dbo.p_set_list AS SET NOCOUNT, XACT_ABORT ON; SELECT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, TIME ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS PROFILE, XML, IO, TIME OFF
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, IO ON
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER
|
< error
| 102 1 15 1 Incorrect syntax near 'QUOTED_IDENTIFIER'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER 1
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS ON, QUOTED_IDENTIFIER ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER ON, ARITHABORT ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS,, QUOTED_IDENTIFIER ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET , ANSI_NULLS ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS QUOTED_IDENTIFIER ON
|
< error
| 102 4 15 1 Incorrect syntax near 'ANSI_NULLS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, FOO ON
|
< error
| 195 5 15 1 'FOO' is not a recognized SET option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, [QUOTED_IDENTIFIER] ON
|
< error
| 102 1 15 1 Incorrect syntax near 'QUOTED_IDENTIFIER'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET [ANSI_NULLS], QUOTED_IDENTIFIER ON
|
< error
| 102 1 15 1 Incorrect syntax near 'ANSI_NULLS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, SHOWPLAN_XML OFF
|
< error
| 1067 1 15 0 The SET SHOWPLAN statements must be the only statements in the batch.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET SHOWPLAN_TEXT, NOCOUNT OFF
|
< error
| 1067 1 15 0 The SET SHOWPLAN statements must be the only statements in the batch.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET DISABLE_DEF_CNST_CHK, NOCOUNT OFF
|
< error
| 195 5 15 1 'DISABLE_DEF_CNST_CHK' is not a recognized SET option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, DATEFORMAT mdy
|
< error
| 195 5 15 1 'DATEFORMAT' is not a recognized SET option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET DATEFORMAT, DATEFIRST mdy
|
< error
| 195 5 15 1 'DATEFORMAT' is not a recognized SET option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET LANGUAGE us_english, NOCOUNT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET DATEFIRST 1, NOCOUNT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ROWCOUNT 10, NOCOUNT ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, ROWCOUNT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ROWCOUNT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, TEXTSIZE ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TEXTSIZE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET ANSI_NULLS, QUOTED_IDENTIFIER OFF, DATEFORMAT mdy
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NOCOUNT, TRANSACTION ISOLATION LEVEL READ COMMITTED
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TRANSACTION'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET TRANSACTION ISOLATION LEVEL READ COMMITTED, NOCOUNT ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NOCOUNT, IDENTITY_INSERT dbo.t ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IDENTITY_INSERT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET IDENTITY_INSERT dbo.t, NOCOUNT ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NOCOUNT, OFFSETS SELECT ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'OFFSETS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET NOCOUNT, STATISTICS IO ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'STATISTICS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, STATISTICS TIME ON
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'STATISTICS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, NOCOUNT ON
|
< error
| 195 8 15 1 'NOCOUNT' is not a recognized statistics option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, FOO ON
|
< error
| 195 8 15 1 'FOO' is not a recognized statistics option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS FOO, IO ON
|
< error
| 195 8 15 1 'FOO' is not a recognized statistics option.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, TIME
|
< error
| 102 1 15 1 Incorrect syntax near 'TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS, IO ON
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SET STATISTICS IO, [TIME] ON
|
< error
| 102 1 15 1 Incorrect syntax near 'TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestSetOptionListOracle]
< ok 0 -
