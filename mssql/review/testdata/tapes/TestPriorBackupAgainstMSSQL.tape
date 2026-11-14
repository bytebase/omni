-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestPriorBackupAgainstMSSQL default#1
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_0] FROM (SELECT [shop].[dbo].[t].* FROM t WHERE id = 1) AS backup_table
< ok 1 -
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[nokey].* FROM nokey WHERE a = 1) AS backup_table
< ok 1 -
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_0] FROM (SELECT [shop].[dbo].[p].* TOP (1) FROM p WHERE id > 0) AS backup_table
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TOP'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[n2].* FROM n2 WHERE id = 1) AS backup_table
< error
| 208 1 16 1 Invalid object name 'n2'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [x].* FROM t AS x JOIN n2 ON x.id = n2.id WHERE x.n = 1) AS backup_table
< error
| 208 1 16 1 Invalid object name 'n2'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[t].* FROM t WHERE id IN (SELECT id FROM n2)) AS backup_table
< error
| 208 1 16 1 Invalid object name 'n2'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[t].* FROM t WHERE gone = 1) AS backup_table
< error
| 207 1 16 1 Invalid column name 'gone'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[nokey].* FROM nokey WHERE nokey.gone = 1) AS backup_table
< error
| 207 1 16 1 Invalid column name 'gone'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [x].* FROM t AS x JOIN p ON p.memo = x.code WHERE x.id = 3) AS backup_table
< error
| 207 1 16 1 Invalid column name 'memo'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_2] FROM (SELECT [shop].[dbo].[nokey].* FROM nokey WHERE b = 1) AS backup_table
< error
| 207 1 16 1 Invalid column name 'b'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[plain].* FROM plain WHERE id IN (SELECT gone FROM t)) AS backup_table
< error
| 207 1 16 1 Invalid column name 'gone'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[t].* FROM t WHERE id = 1) AS backup_table
< ok 1 -
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[nokey].* FROM dbo.nokey WHERE b = 1) AS backup_table
< error
| 207 1 16 1 Invalid column name 'b'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_0] FROM (SELECT [shop].[dbo].[t].* FROM dbo.p AS p) AS backup_table
< error
| 107 1 15 1 The column prefix 'shop.dbo.t' does not match with a table name or alias name used in the query.
| 8155 2 16 1 No column name was specified for column 1 of 'backup_table'.
| 1038 5 15 1 An object or column name is missing or empty. For SELECT INTO statements, verify each column has a name. For other statements, look for empty alias names. Aliases defined as "" or [] are not allowed. Change the alias to a valid name.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_0] FROM (SELECT [shop].[dbo].[t].* FROM dbo.p AS t) AS backup_table
< error
| 107 1 15 1 The column prefix 'shop.dbo.t' does not match with a table name or alias name used in the query.
| 8155 2 16 1 No column name was specified for column 1 of 'backup_table'.
| 1038 5 15 1 An object or column name is missing or empty. For SELECT INTO statements, verify each column has a name. For other statements, look for empty alias names. Aliases defined as "" or [] are not allowed. Change the alias to a valid name.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_0] FROM (SELECT [shop].[dbo].[t].* FROM dbo.p AS t2 JOIN s.t2 AS s ON s.id = t2.id) AS backup_table
< error
| 208 1 16 1 Invalid object name 's.t2'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[n8].* FROM dbo.n8 WHERE a = 1) AS backup_table
< error
| 208 1 16 1 Invalid object name 'dbo.n8'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[n6].* FROM s3.n6 WHERE id = 1) AS backup_table
< error
| 208 1 16 1 Invalid object name 's3.n6'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[n4].* FROM n4 WHERE id = 1) AS backup_table
< error
| 208 1 16 1 Invalid object name 'n4'.
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop
< ok 0 -
> exec
| IF DB_ID('bbdataarchive') IS NOT NULL BEGIN ALTER DATABASE bbdataarchive SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE bbdataarchive; END
< ok 0 -
> exec
| CREATE DATABASE bbdataarchive
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.p (id INT NOT NULL CONSTRAINT pk_p PRIMARY KEY, note VARCHAR(20))
< ok 0 -
> exec
| CREATE TABLE dbo.t (id INT NOT NULL CONSTRAINT pk_t PRIMARY KEY, code VARCHAR(20) CONSTRAINT uk_code UNIQUE, p_id INT CONSTRAINT fk_p REFERENCES dbo.p (id), n INT CONSTRAINT ck_n CHECK (n > 0))
< ok 0 -
> exec
| CREATE INDEX ix_n ON dbo.t (n)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a INT)
< ok 0 -
> exec
| CREATE TABLE dbo.plain (id INT NOT NULL CONSTRAINT pk_plain PRIMARY KEY, v INT)
< ok 0 -
> exec
| INSERT INTO dbo.p VALUES (1, 'a'), (2, 'b')
< ok 2 -
> exec
| INSERT INTO dbo.t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2)
< ok 2 -
> exec
| INSERT INTO dbo.nokey VALUES (1), (2)
< ok 2 -
> exec
| INSERT INTO dbo.plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT s.name, o.name, o.type FROM sys.objects o JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type IN ('U', 'V') AND o.is_ms_shipped = 0 ORDER BY s.name, o.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| nokey
: string
| U 
< row
: string
| dbo
: string
| p
: string
| U 
< row
: string
| dbo
: string
| plain
: string
| U 
< row
: string
| dbo
: string
| t
: string
| U 
> query
| SELECT s.name, o.name, c.name FROM sys.columns c JOIN sys.objects o ON o.object_id = c.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
: string
| a
< row
: string
| dbo
: string
| p
: string
| id
< row
: string
| dbo
: string
| p
: string
| note
< row
: string
| dbo
: string
| plain
: string
| id
< row
: string
| dbo
: string
| plain
: string
| v
< row
: string
| dbo
: string
| t
: string
| id
< row
: string
| dbo
: string
| t
: string
| code
< row
: string
| dbo
: string
| t
: string
| p_id
< row
: string
| dbo
: string
| t
: string
| n
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint FROM sys.indexes i JOIN sys.objects o ON o.object_id = i.object_id JOIN sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' AND i.name IS NOT NULL ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< row
: string
| dbo
: string
| p
: string
| pk_p
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| plain
: string
| pk_plain
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| ix_n
: bool false
: bool false
: bool false
< row
: string
| dbo
: string
| t
: string
| pk_t
: bool true
: bool true
: bool false
< row
: string
| dbo
: string
| t
: string
| uk_code
: bool false
: bool true
: bool true
> query
| SELECT s.name, o.name, f.name FROM sys.foreign_keys f JOIN sys.objects o ON o.object_id = f.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, f.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| fk_p
> query
| SELECT s.name, o.name, c.name FROM sys.check_constraints c JOIN sys.objects o ON o.object_id = c.parent_object_id JOIN sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, c.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| t
: string
| ck_n
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_1] FROM (SELECT [shop].[dbo].[nokey].* FROM nokey WHERE gone = 1) AS backup_table
< error
| 207 1 16 1 Invalid column name 'gone'.
> exec
| SELECT * INTO [bbdataarchive].[dbo].[bk_3] FROM (SELECT [shop].[dbo].[nokey].* FROM shop.dbo.nokey WHERE a = 2) AS backup_table
< ok 1 -
