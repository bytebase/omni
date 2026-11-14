-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestColumnOptionOrderOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) SPARSE FILESTREAM NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM SPARSE NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) COLLATE Latin1_General_CI_AS SPARSE FILESTREAM MASKED WITH (FUNCTION = 'default()') NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 COLLATE Latin1_General_CI_AS GENERATED ALWAYS AS ROW END HIDDEN NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT IDENTITY NOT FOR REPLICATION PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT IDENTITY(1, 1) NOT FOR REPLICATION IDENTITY NOT FOR REPLICATION)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT NOT NULL REFERENCES dbo.t (a) NOT FOR REPLICATION DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT NULL NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a UNIQUEIDENTIFIER ROWGUIDCOL DEFAULT NEWID() ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 'x' COLLATE Latin1_General_BIN NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT CHECK NOT FOR REPLICATION (a > 0) CONSTRAINT nn NOT NULL UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c INT SPARSE NULL CONSTRAINT ck CHECK (c > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM FILESTREAM)
|
< error
| 102 1 15 1 Incorrect syntax near 'FILESTREAM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) NULL FILESTREAM)
|
< error
| 102 1 15 1 Incorrect syntax near 'FILESTREAM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) MASKED WITH (FUNCTION = 'default()') FILESTREAM)
|
< error
| 102 1 15 1 Incorrect syntax near 'FILESTREAM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) FILESTREAM ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'FILESTREAM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARBINARY(MAX) SPARSE FILESTREAM ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 SPARSE GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 NOT NULL GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a DATETIME2 GENERATED ALWAYS AS ROW START HIDDEN HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT NOT FOR REPLICATION IDENTITY)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT IDENTITY NULL NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT IDENTITY NOT FOR REPLICATION NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a INT CHECK (a > 0) NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c INT NOT NULL SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c INT NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @t TABLE (a INT NOT NULL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TYPE dbo.tt AS TABLE (a INT NULL SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FUNCTION dbo.f() RETURNS @t TABLE (a INT NOT FOR REPLICATION) AS BEGIN RETURN END
|
< error
| 102 1 15 1 proc=f Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionOrderOracle]
< ok 0 -
