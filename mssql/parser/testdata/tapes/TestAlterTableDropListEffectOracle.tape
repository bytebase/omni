-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestAlterTableDropListEffectOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP COLUMN a, b
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP ck_probe, COLUMN a
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP CONSTRAINT ck_probe, df_probe, COLUMN c, d
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP COLUMN a, CONSTRAINT ck_probe, df_probe
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP COLUMN a, ck_probe
< error
| 4924 1 16 1 ALTER TABLE DROP COLUMN failed because column 'ck_probe' does not exist in table 'drop_probe'.
> rollback
< error
| 3903 2 16 1 The ROLLBACK TRANSACTION request has no corresponding BEGIN TRANSACTION.
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP COLUMN IF EXISTS nope, b
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP CONSTRAINT IF EXISTS nope1, CONSTRAINT IF EXISTS ck_probe, COLUMN IF EXISTS nope2
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP COLUMN IF EXISTS a, nope
< error
| 4924 1 16 1 ALTER TABLE DROP COLUMN failed because column 'nope' does not exist in table 'drop_probe'.
> rollback
< error
| 3903 2 16 1 The ROLLBACK TRANSACTION request has no corresponding BEGIN TRANSACTION.
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP CONSTRAINT IF EXISTS ck_probe, nope
< error
| 3728 1 16 1 'nope' is not a constraint.
| 3727 0 16 1 Could not drop constraint. See previous errors.
> rollback
< error
| 3903 2 16 1 The ROLLBACK TRANSACTION request has no corresponding BEGIN TRANSACTION.
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP PERIOD FOR SYSTEM_TIME, COLUMN vs, ve
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
> rollback
< ok
> begin
< ok
> exec
| CREATE TABLE dbo.drop_probe (
| 	id INT CONSTRAINT pk_probe PRIMARY KEY,
| 	a INT,
| 	b INT,
| 	c INT CONSTRAINT ck_probe CHECK (c > 0),
| 	d INT CONSTRAINT df_probe DEFAULT 1,
| 	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
| 	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
| 	PERIOD FOR SYSTEM_TIME (vs, ve)
| )
< ok 0 -
> query
| SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
| UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
| UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')
| ORDER BY 1
< result
< column "" "NVARCHAR"
< row
: string
| column a
< row
: string
| column b
< row
: string
| column c
< row
: string
| column d
< row
: string
| column id
< row
: string
| column ve
< row
: string
| column vs
< row
: string
| constraint ck_probe
< row
: string
| constraint df_probe
< row
: string
| constraint pk_probe
< row
: string
| period
> exec
| ALTER TABLE dbo.drop_probe DROP PERIOD FOR SYSTEM_TIME, ck_probe
< error
| 13594 1 16 1 Table 'omni_TestAlterTableDropListEffectOracle.dbo.drop_probe' does not contain SYSTEM_TIME period 'ck_probe'.
> rollback
< error
| 3903 2 16 1 The ROLLBACK TRANSACTION request has no corresponding BEGIN TRANSACTION.
