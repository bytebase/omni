-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestWalkThroughBatchesAgainstMSSQL default#1
> exec
| USE master
< ok 0 -
> exec
| IF DB_ID('shop') IS NOT NULL BEGIN ALTER DATABASE shop SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE shop; END
< ok 0 -
> exec
| CREATE DATABASE shop COLLATE SQL_Latin1_General_CP1_CI_AS
< ok 0 -
> exec
| USE shop
< ok 0 -
> exec
| CREATE TABLE dbo.r (id int NOT NULL CONSTRAINT pk_r PRIMARY KEY)
< ok 0 -
> exec
| CREATE TABLE dbo.t (id int NOT NULL CONSTRAINT pk_t PRIMARY KEY, a int CONSTRAINT df_a DEFAULT 0, b int CONSTRAINT ck_b CHECK (b > 0), c int CONSTRAINT uq_c UNIQUE, d int CONSTRAINT fk_r REFERENCES dbo.r (id))
< ok 0 -
> exec
| CREATE INDEX ix_a ON dbo.t (a)
< ok 0 -
> exec
| CREATE TABLE dbo.nokey (a int)
< ok 0 -
> exec
| CREATE VIEW dbo.v AS SELECT id FROM dbo.r
< ok 0 -
> exec
| CREATE PROCEDURE dbo.pr AS SELECT 1
< ok 0 -
> exec
| CREATE FUNCTION dbo.fn() RETURNS int AS BEGIN RETURN 1 END
< ok 0 -
> exec
| CREATE SEQUENCE dbo.sq
< ok 0 -
> exec
| CREATE SYNONYM dbo.syn FOR dbo.t
< ok 0 -
> exec
| CREATE SCHEMA s
< ok 0 -
> exec
| CREATE TABLE s.st (a int)
< ok 0 -
> exec
| CREATE TABLE s.t (id int)
< ok 0 -
> exec
| CREATE USER u_s WITHOUT LOGIN WITH DEFAULT_SCHEMA = s
< ok 0 -
> exec
| GRANT CONTROL ON DATABASE::shop TO u_s
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE TABLE dbo.x1 (a int);
| CREATE VIEW dbo.v1 AS SELECT 1 AS a;
< error
| 111 1 15 2 'CREATE VIEW' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SET NOCOUNT ON;
| CREATE VIEW dbo.v1 AS SELECT 1 AS a;
< error
| 111 1 15 2 'CREATE VIEW' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1;
| CREATE PROCEDURE dbo.p1 AS SELECT 1;
< error
| 111 1 15 2 proc=p1 'CREATE/ALTER PROCEDURE' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1;
| ALTER PROCEDURE dbo.pr AS SELECT 2;
< error
| 111 1 15 2 proc=pr 'CREATE/ALTER PROCEDURE' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1;
| CREATE FUNCTION dbo.f1() RETURNS int AS BEGIN RETURN 1 END;
< error
| 111 1 15 2 'CREATE FUNCTION' must be the first statement in a query batch.
| 178 1 15 2 A RETURN statement with a return value cannot be used in this context.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1;
| CREATE TRIGGER dbo.tr1 ON dbo.nokey AFTER INSERT AS SELECT 1;
< error
| 111 1 15 2 'CREATE TRIGGER' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1;
| CREATE SCHEMA s9;
< error
| 111 1 15 2 'CREATE SCHEMA' must be the first statement in a query batch.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE TABLE dbo.x1 (a int);
< ok 0 -
> exec
| CREATE VIEW dbo.v1 AS SELECT 1 AS a;
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| -- note
| /* note */
| CREATE VIEW dbo.v1 AS SELECT 1 AS a;
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE PROCEDURE dbo.p1 AS SELECT 1;
| CREATE TABLE dbo.t (a int);
| DROP TABLE dbo.nosuch;
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE PROCEDURE dbo.p1 AS BEGIN SELECT 1 END;
| CREATE TABLE dbo.t (a int);
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE TRIGGER dbo.tr1 ON dbo.nokey AFTER UPDATE AS SELECT 1;
| CREATE TABLE dbo.t (a int);
< ok 0 -
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE PROCEDURE dbo.p1 AS SELECT 1;
< ok 0 -
> exec
| CREATE TABLE dbo.t (a int);
< error
| 2714 6 16 1 There is already an object named 't' in the database.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE VIEW dbo.v1 AS SELECT 1 AS a;
| CREATE TABLE dbo.t (a int);
< error
| 156 1 15 2 proc=v1 Incorrect syntax near the keyword 'CREATE'.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| CREATE TABLE dbo.g1 (a int)
< ok 0 -
> exec ^
< error
| 2714 6 16 1 There is already an object named 'g1' in the database.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
> exec ^
< ok 0 -
> exec
| REVERT
< error
| 2812 62 16 1 proc=REVERT Could not find stored procedure 'REVERT'.
> exec
| USE shop
< ok 0 -
> exec
| BEGIN TRANSACTION
< ok 0 -
> query
| SELECT schema_name FROM shop.INFORMATION_SCHEMA.SCHEMATA WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys') ORDER BY schema_name
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| s
> query
| SELECT s.name, t.name FROM shop.sys.tables t JOIN shop.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| nokey
< row
: string
| dbo
: string
| r
< row
: string
| dbo
: string
| t
< row
: string
| s
: string
| st
< row
: string
| s
: string
| t
> query
| SELECT s.name, v.name FROM shop.sys.views v JOIN shop.sys.sql_modules m ON m.object_id = v.object_id JOIN shop.sys.schemas s ON s.schema_id = v.schema_id ORDER BY s.name, v.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| v
> query
| SELECT s.name, o.name, c.name, ISNULL(dc.name, ''), CASE WHEN c.is_computed = 1 THEN 'AS ' + ISNULL(cc.definition, '') ELSE ISNULL(ty.name, '') END, c.is_nullable FROM shop.sys.columns c
| 		JOIN shop.sys.objects o ON o.object_id = c.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN shop.sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN shop.sys.types ty ON ty.user_type_id = c.user_type_id
| 		LEFT JOIN shop.sys.default_constraints dc ON dc.object_id = c.default_object_id
| 		WHERE o.type = 'U'
| 		ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "is_nullable" "BIT"
< row
: string
| dbo
: string
| nokey
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| r
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| id
: string
|
: string
| int
: bool false
< row
: string
| dbo
: string
| t
: string
| a
: string
| df_a
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| b
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| c
: string
|
: string
| int
: bool true
< row
: string
| dbo
: string
| t
: string
| d
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| st
: string
| a
: string
|
: string
| int
: bool true
< row
: string
| s
: string
| t
: string
| id
: string
|
: string
| int
: bool true
> query
| SELECT s.name, o.name, i.name, i.is_primary_key, i.is_unique, i.is_unique_constraint, i.type_desc FROM shop.sys.indexes i
| 		JOIN shop.sys.objects o ON o.object_id = i.object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE i.index_id > 0 AND i.name IS NOT NULL AND o.type = 'U'
| 		AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1 OR EXISTS (SELECT 1 FROM shop.sys.index_columns ic JOIN shop.sys.all_columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id WHERE ic.object_id = i.object_id AND ic.index_id = i.index_id))
| 		ORDER BY s.name, o.name, i.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "is_primary_key" "BIT"
< column "is_unique" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< row
: string
| dbo
: string
| r
: string
| pk_r
: bool true
: bool true
: bool false
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| ix_a
: bool false
: bool false
: bool false
: string
| NONCLUSTERED
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
: string
| CLUSTERED
< row
: string
| dbo
: string
| t
: string
| uq_c
: bool false
: bool true
: bool true
: string
| NONCLUSTERED
> query
| SELECT s.name, o.name, f.name FROM shop.sys.foreign_keys f JOIN shop.sys.objects o ON o.object_id = f.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE EXISTS (SELECT 1 FROM shop.sys.foreign_key_columns fkc WHERE fkc.constraint_object_id = f.object_id)
| 		ORDER BY s.name, o.name, f.name
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
| fk_r
> query
| SELECT s.name, o.name, c.name FROM shop.sys.check_constraints c JOIN shop.sys.objects o ON o.object_id = c.parent_object_id JOIN shop.sys.schemas s ON s.schema_id = o.schema_id WHERE o.type = 'U' ORDER BY s.name, o.name, c.name
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
| ck_b
> query
| SELECT ss.name, so.name, so.type, st.name FROM shop.sys.triggers st JOIN shop.sys.sql_modules m ON m.object_id = st.object_id
| 		JOIN shop.sys.objects so ON so.object_id = st.parent_id JOIN shop.sys.schemas ss ON ss.schema_id = so.schema_id
| 		WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| 		ORDER BY ss.name, so.name, st.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
> query
| SELECT s.name, q.name FROM shop.sys.sequences q JOIN shop.sys.schemas s ON s.schema_id = q.schema_id ORDER BY s.name, q.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| sq
> query
| SELECT s.name, ao.name, ao.type FROM shop.sys.all_objects ao JOIN shop.sys.all_sql_modules m ON m.object_id = ao.object_id JOIN shop.sys.schemas s ON s.schema_id = ao.schema_id
| 		WHERE ao.type IN ('FN', 'IF', 'TF', 'P', 'RF') AND ao.is_ms_shipped = 0
| 		AND NOT EXISTS (SELECT 1 FROM shop.sys.extended_properties ep WHERE ep.major_id = ao.object_id AND ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'microsoft_database_tools_support')
| 		ORDER BY s.name, ao.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< row
: string
| dbo
: string
| fn
: string
| FN
< row
: string
| dbo
: string
| pr
: string
| P 
> query
| SELECT s.name, sy.name FROM shop.sys.synonyms sy JOIN shop.sys.schemas s ON s.schema_id = sy.schema_id ORDER BY s.name, sy.name
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: string
| dbo
: string
| syn
> exec
| SELECT 1
| GO -- note
| CREATE TABLE dbo.t (a int);
< error
| 2714 6 16 3 There is already an object named 't' in the database.
> exec
| IF @@TRANCOUNT > 0 ROLLBACK
< ok 0 -
