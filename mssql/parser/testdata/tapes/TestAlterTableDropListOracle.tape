-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestAlterTableDropListOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT fk_p, CONSTRAINT ck_n
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT uk_code, COLUMN n
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT ck_n
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT [fk p], COLUMN [n]
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, b
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP fk_p, ck_n
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP ck_n, COLUMN a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, b, CONSTRAINT ck_n, df_n, COLUMN c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN IF EXISTS a, b, COLUMN IF EXISTS c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT IF EXISTS ck_n WITH (ONLINE = ON)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP IF EXISTS ck_n
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, IF EXISTS b
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS a, IF EXISTS b
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP IF EXISTS COLUMN a
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH (ONLINE = ON), COLUMN b
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH (ONLINE = ON), uk_code WITH (MAXDOP = 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, b WITH (ONLINE = ON)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a WITH (ONLINE = ON)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN IF EXISTS a WITH (ONLINE = ON), b
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH
|
< error
| 102 1 15 1 Incorrect syntax near 'WITH'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME, COLUMN vs, ve
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT ck_n, PERIOD FOR SYSTEM_TIME
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME, ck_n WITH (ONLINE = ON)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, PERIOD
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME WITH (ONLINE = ON)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT ck_n, PERIOD FOR
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP PERIOD SYSTEM_TIME
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSTEM_TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP
|
< error
| 102 1 15 1 Incorrect syntax near 'DROP'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT
|
< error
| 102 1 15 1 Incorrect syntax near 'CONSTRAINT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN
|
< error
| 102 1 15 1 Incorrect syntax near 'COLUMN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a,
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT ck_n, COLUMN
|
< error
| 102 1 15 1 Incorrect syntax near 'COLUMN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT
|
< error
| 102 1 15 1 Incorrect syntax near 'CONSTRAINT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN a, , b
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP COLUMN COLUMN a
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLUMN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER TABLE dbo.t DROP CONSTRAINT COLUMN a
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLUMN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAlterTableDropListOracle]
< ok 0 -
