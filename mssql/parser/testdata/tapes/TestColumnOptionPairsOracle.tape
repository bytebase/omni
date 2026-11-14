-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestColumnOptionPairsOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS SPARSE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS MASKED WITH (FUNCTION = 'default()'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS GENERATED ALWAYS AS ROW START)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) COLLATE Latin1_General_CI_AS REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE MASKED WITH (FUNCTION = 'default()'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) SPARSE REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) MASKED WITH (FUNCTION = 'default()') REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START HIDDEN)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) GENERATED ALWAYS AS ROW START REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN COLLATE Latin1_General_CI_AS)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN NOT NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN ROWGUIDCOL)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN IDENTITY(1, 1))
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN DEFAULT 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN CONSTRAINT df DEFAULT 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN CONSTRAINT nn NOT NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN PRIMARY KEY)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN UNIQUE)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN CHECK (a > 0))
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) HIDDEN REFERENCES dbo.t (a))
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT NULL REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NULL REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) ROWGUIDCOL REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) NOT FOR REPLICATION)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) IDENTITY(1, 1) REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION COLLATE Latin1_General_CI_AS)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION NOT NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION ROWGUIDCOL)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION IDENTITY(1, 1))
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION DEFAULT 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION CONSTRAINT df DEFAULT 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION CONSTRAINT nn NOT NULL)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION PRIMARY KEY)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION UNIQUE)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION CHECK (a > 0))
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) NOT FOR REPLICATION REFERENCES dbo.t (a))
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 COLLATE Latin1_General_CI_AS)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
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
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) DEFAULT 1 REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 COLLATE Latin1_General_CI_AS)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT df DEFAULT 1 REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CONSTRAINT nn NOT NULL REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) PRIMARY KEY REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) UNIQUE REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) NOT FOR REPLICATION)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) CHECK (a > 0) REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) COLLATE Latin1_General_CI_AS)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) SPARSE)
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) MASKED WITH (FUNCTION = 'default()'))
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'))
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) GENERATED ALWAYS AS ROW START)
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) HIDDEN)
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) ROWGUIDCOL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) IDENTITY(1, 1))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) NOT FOR REPLICATION)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) CONSTRAINT df DEFAULT 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) CONSTRAINT nn NOT NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) PRIMARY KEY)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) UNIQUE)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) CHECK (a > 0))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t9 (a VARCHAR(10) REFERENCES dbo.t (a) REFERENCES dbo.t (a))
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS SPARSE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS MASKED WITH (FUNCTION = 'default()')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS GENERATED ALWAYS AS ROW START
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) COLLATE Latin1_General_CI_AS REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE MASKED WITH (FUNCTION = 'default()')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) SPARSE REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) MASKED WITH (FUNCTION = 'default()') REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256') REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START HIDDEN
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) GENERATED ALWAYS AS ROW START REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN COLLATE Latin1_General_CI_AS
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN NOT NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN ROWGUIDCOL
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN IDENTITY(1, 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN DEFAULT 1
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN CONSTRAINT df DEFAULT 1
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN CONSTRAINT nn NOT NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN PRIMARY KEY
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN UNIQUE
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN CHECK (a > 0)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) HIDDEN REFERENCES dbo.t (a)
|
< error
| 102 1 15 1 Incorrect syntax near 'hidden'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT NULL REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NULL REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) ROWGUIDCOL REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) NOT FOR REPLICATION
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) IDENTITY(1, 1) REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION COLLATE Latin1_General_CI_AS
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION NOT NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION ROWGUIDCOL
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION IDENTITY(1, 1)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION DEFAULT 1
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION CONSTRAINT df DEFAULT 1
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION CONSTRAINT nn NOT NULL
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION PRIMARY KEY
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION UNIQUE
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION CHECK (a > 0)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) NOT FOR REPLICATION REFERENCES dbo.t (a)
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) DEFAULT 1 REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT df DEFAULT 1 REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CONSTRAINT nn NOT NULL REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) PRIMARY KEY REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) UNIQUE REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) NOT FOR REPLICATION
|
< error
| 102 1 15 1 Incorrect syntax near 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) CHECK (a > 0) REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) SPARSE
|
< error
| 102 1 15 1 Incorrect syntax near 'SPARSE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) MASKED WITH (FUNCTION = 'default()')
|
< error
| 102 1 15 1 Incorrect syntax near 'MASKED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = k, ENCRYPTION_TYPE = DETERMINISTIC, ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256')
|
< error
| 102 1 15 1 Incorrect syntax near 'ENCRYPTED'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) GENERATED ALWAYS AS ROW START
|
< error
| 102 1 15 1 Incorrect syntax near 'GENERATED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) HIDDEN
|
< error
| 102 1 15 1 Incorrect syntax near 'HIDDEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) ROWGUIDCOL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) IDENTITY(1, 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) NOT FOR REPLICATION
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) CONSTRAINT df DEFAULT 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) CONSTRAINT nn NOT NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) PRIMARY KEY
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) UNIQUE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) CHECK (a > 0)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t ADD c VARCHAR(10) REFERENCES dbo.t (a) REFERENCES dbo.t (a)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestColumnOptionPairsOracle]
< ok 0 -
