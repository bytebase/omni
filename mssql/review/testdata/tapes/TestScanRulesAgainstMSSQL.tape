-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestScanRulesAgainstMSSQL default#1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE t DROP CONSTRAINT fk_p, ck_n
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE t DROP CONSTRAINT uk_code
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE plain DROP CONSTRAINT pk_plain
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE nokey ADD CONSTRAINT ck_a CHECK (a > 0)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| nokey
: string
| C
: string
| ck_a
: string
| ([a]>(0))
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| nokey
: string
| C
: string
| ck_a
: string
| ([a]>(0))
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE nokey DROP CONSTRAINT ck_a
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE nokey ADD CONSTRAINT df_a DEFAULT 0 FOR a
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE nokey DROP CONSTRAINT df_a
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| CREATE TABLE a (id INT)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| CREATE TABLE b (id INT PRIMARY KEY)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| b
: string
| PK
: string
| PK__b__standin1
: string
|
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| b
: int64 1
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| b
: int64 1
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| CREATE TABLE a (id INT NOT NULL)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 0
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE a ADD CONSTRAINT pk_a PRIMARY KEY (id)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| a
: string
| PK
: string
| pk_a
: string
|
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 1
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| a
: int64 1
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE plain DROP CONSTRAINT pk_plain
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
> exec
| ALTER TABLE plain ADD CONSTRAINT pk_plain2 PRIMARY KEY (id)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain2
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| BEGIN TRAN
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| ALTER TABLE plain DROP CONSTRAINT pk_plain
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 0
< row
: string
| t
: int64 1
> exec
| ROLLBACK
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| IF 1 = 1 ALTER TABLE t DROP CONSTRAINT ck_n
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
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
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| CREATE DATABASE d
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| CREATE TABLE d.dbo.e (id INT)
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> exec
| DROP DATABASE d
< ok 0 -
> query
|
| 		SELECT OBJECT_NAME(k.parent_object_id), k.type, k.name, '' FROM sys.key_constraints k
| 		UNION ALL SELECT OBJECT_NAME(f.parent_object_id), 'F', f.name, '' FROM sys.foreign_keys f
| 		UNION ALL SELECT OBJECT_NAME(c.parent_object_id), 'C', c.name, c.definition FROM sys.check_constraints c
| 		ORDER BY 1, 2, 3
< result
< column "" "NVARCHAR"
< column "type" "VARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| p
: string
| PK
: string
| pk_p
: string
|
< row
: string
| plain
: string
| PK
: string
| pk_plain
: string
|
< row
: string
| t
: string
| C
: string
| ck_n
: string
| ([n]>(0))
< row
: string
| t
: string
| F
: string
| fk_p
: string
|
< row
: string
| t
: string
| PK
: string
| pk_t
: string
|
< row
: string
| t
: string
| UQ
: string
| uk_code
: string
|
> query
| SELECT o.name, CASE WHEN EXISTS (SELECT 1 FROM sys.key_constraints k WHERE k.parent_object_id = o.object_id AND k.type = 'PK') THEN 1 ELSE 0 END FROM sys.objects o WHERE o.type = 'U' ORDER BY o.name
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
> query ^
< result
< column "name" "NVARCHAR"
< column "" "INT"
< row
: string
| nokey
: int64 0
< row
: string
| p
: int64 1
< row
: string
| plain
: int64 1
< row
: string
| t
: int64 1
