-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestProgramSimulation default#1
> exec
| CREATE OR ALTER PROCEDURE [dbo].[_sim_reset] AS
| BEGIN
| 	SET NOCOUNT ON;
| 	DECLARE @sql NVARCHAR(MAX) = N'';
| 	SELECT @sql += N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N' SET (SYSTEM_VERSIONING = OFF);' + CHAR(10)
| 	FROM sys.tables WHERE temporal_type = 2;
| 	IF @sql <> N'' EXEC sp_executesql @sql;
| 	SET @sql = N'';
| 	SELECT @sql += N'DROP TRIGGER ' + QUOTENAME(name) + N' ON DATABASE;' + CHAR(10)
| 	FROM sys.triggers WHERE parent_class = 0;
| 	SELECT @sql += N'DROP VIEW ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.views WHERE is_ms_shipped = 0;
| 	SELECT @sql += N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(parent.schema_id)) + N'.' + QUOTENAME(parent.name) + N' DROP CONSTRAINT ' + QUOTENAME(fk.name) + N';' + CHAR(10)
| 	FROM sys.foreign_keys fk JOIN sys.objects parent ON parent.object_id = fk.parent_object_id;
| 	IF @sql <> N'' EXEC sp_executesql @sql;
| 	SET @sql = N'';
| 	SELECT @sql += N'DROP TABLE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.tables WHERE is_ms_shipped = 0;
| 	IF @sql <> N'' EXEC sp_executesql @sql;
| 	SET @sql = N'';
| 	SELECT @sql += N'DROP PROCEDURE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.procedures WHERE is_ms_shipped = 0 AND NOT (schema_id = SCHEMA_ID(N'dbo') AND name IN (N'_sim_reset', N'_sim_forms'));
| 	SELECT @sql += N'DROP FUNCTION ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.objects WHERE type IN ('FN', 'IF', 'TF') AND is_ms_shipped = 0;
| 	SELECT @sql += N'DROP SEQUENCE ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.sequences;
| 	SELECT @sql += N'DROP SYNONYM ' + QUOTENAME(SCHEMA_NAME(schema_id)) + N'.' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.synonyms;
| 	IF @sql <> N'' EXEC sp_executesql @sql;
| 	SET @sql = N'';
| 	SELECT @sql += N'DROP SCHEMA ' + QUOTENAME(name) + N';' + CHAR(10)
| 	FROM sys.schemas WHERE schema_id > 4 AND schema_id < 16384;
| 	IF @sql <> N'' EXEC sp_executesql @sql;
| END;
< ok 0 -
> exec
| CREATE OR ALTER PROCEDURE [dbo].[_sim_forms] AS
| BEGIN
| 	SET NOCOUNT ON;
| WITH key_cols AS (
| 	SELECT i.object_id, i.index_id,
| 		STRING_AGG(CAST(QUOTENAME(c2.name) + CASE WHEN ic2.is_descending_key = 1 THEN N' DESC' ELSE N' ASC' END AS NVARCHAR(MAX)), N', ')
| 			WITHIN GROUP (ORDER BY ic2.key_ordinal) AS cols
| 	FROM sys.indexes i
| 	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id AND ic2.is_included_column = 0
| 	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
| 	WHERE i.type > 0
| 	GROUP BY i.object_id, i.index_id
| ),
| inc_cols AS (
| 	SELECT i.object_id, i.index_id,
| 		STRING_AGG(CAST(QUOTENAME(c2.name) AS NVARCHAR(MAX)), N', ')
| 			WITHIN GROUP (ORDER BY ic2.index_column_id) AS cols
| 	FROM sys.indexes i
| 	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id AND ic2.is_included_column = 1
| 	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
| 	WHERE i.type > 0
| 	GROUP BY i.object_id, i.index_id
| ),
| all_idx_cols AS (
| 	SELECT i.object_id, i.index_id,
| 		STRING_AGG(CAST(QUOTENAME(c2.name) AS NVARCHAR(MAX)), N', ')
| 			WITHIN GROUP (ORDER BY ic2.index_column_id) AS cols
| 	FROM sys.indexes i
| 	JOIN sys.index_columns ic2 ON ic2.object_id = i.object_id AND ic2.index_id = i.index_id
| 	JOIN sys.columns c2 ON c2.object_id = ic2.object_id AND c2.column_id = ic2.column_id
| 	WHERE i.type = 6
| 	GROUP BY i.object_id, i.index_id
| ),
| fk_parent_cols AS (
| 	SELECT fkc.constraint_object_id,
| 		STRING_AGG(CAST(QUOTENAME(pc.name) AS NVARCHAR(MAX)), N', ')
| 			WITHIN GROUP (ORDER BY fkc.constraint_column_id) AS cols
| 	FROM sys.foreign_key_columns fkc
| 	JOIN sys.columns pc ON pc.object_id = fkc.parent_object_id AND pc.column_id = fkc.parent_column_id
| 	GROUP BY fkc.constraint_object_id
| ),
| fk_ref_cols AS (
| 	SELECT fkc.constraint_object_id,
| 		STRING_AGG(CAST(QUOTENAME(rc.name) AS NVARCHAR(MAX)), N', ')
| 			WITHIN GROUP (ORDER BY fkc.constraint_column_id) AS cols
| 	FROM sys.foreign_key_columns fkc
| 	JOIN sys.columns rc ON rc.object_id = fkc.referenced_object_id AND rc.column_id = fkc.referenced_column_id
| 	GROUP BY fkc.constraint_object_id
| ),
| tbl_cols AS (
| 	SELECT c.object_id,
| 		STRING_AGG(
| 			CAST(
| 				QUOTENAME(c.name) + N' ' +
| 				CASE
| 					WHEN cc.object_id IS NOT NULL THEN
| 						N'AS ' + cc.definition
| 						+ CASE WHEN cc.is_persisted = 1 THEN N' PERSISTED' ELSE N'' END
| 						+ CASE WHEN cc.is_persisted = 1 AND c.is_nullable = 0 THEN N' NOT NULL' ELSE N'' END
| 					ELSE
| 						CASE
| 							WHEN ty.name IN (N'decimal', N'numeric') THEN
| 								ty.name + N'(' + CAST(c.precision AS NVARCHAR(10)) + N',' + CAST(c.scale AS NVARCHAR(10)) + N')'
| 							WHEN ty.name IN (N'datetime2', N'datetimeoffset', N'time') THEN
| 								ty.name + N'(' + CAST(c.scale AS NVARCHAR(10)) + N')'
| 							WHEN ty.name IN (N'nvarchar', N'nchar') THEN
| 								ty.name + N'(' + CASE WHEN c.max_length = -1 THEN N'max' ELSE CAST(c.max_length / 2 AS NVARCHAR(10)) END + N')'
| 							WHEN ty.name IN (N'varchar', N'char', N'varbinary', N'binary') THEN
| 								ty.name + N'(' + CASE WHEN c.max_length = -1 THEN N'max' ELSE CAST(c.max_length AS NVARCHAR(10)) END + N')'
| 							ELSE ty.name
| 						END
| 						+ CASE WHEN ic.object_id IS NOT NULL THEN
| 							N' IDENTITY(' + CAST(ic.seed_value AS NVARCHAR(40)) + N',' + CAST(ic.increment_value AS NVARCHAR(40)) + N')'
| 							ELSE N'' END
| 						+ CASE WHEN c.is_nullable = 0 THEN N' NOT NULL' ELSE N' NULL' END
| 						+ CASE WHEN dc.object_id IS NOT NULL THEN
| 							CASE WHEN dc.is_system_named = 0 THEN N' CONSTRAINT ' + QUOTENAME(dc.name) ELSE N'' END
| 							+ N' DEFAULT ' + dc.definition
| 							ELSE N'' END
| 				END
| 			AS NVARCHAR(MAX)),
| 			N', '
| 		) WITHIN GROUP (ORDER BY c.column_id) AS text
| 	FROM sys.columns c
| 	JOIN sys.tables t ON t.object_id = c.object_id AND t.is_ms_shipped = 0
| 	JOIN sys.types ty ON ty.user_type_id = c.user_type_id
| 	LEFT JOIN sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 	LEFT JOIN sys.identity_columns ic ON ic.object_id = c.object_id AND ic.column_id = c.column_id
| 	LEFT JOIN sys.default_constraints dc ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
| 	WHERE c.is_hidden = 0 AND c.graph_type IS NULL
| 	GROUP BY c.object_id
| ),
| tbl_kcs AS (
| 	SELECT kc.parent_object_id AS object_id,
| 		STRING_AGG(
| 			CAST(
| 				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
| 				+ CASE WHEN kc.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(kc.name) + N' ' ELSE N'' END
| 				+ CASE kc.type WHEN 'PK' THEN N'PRIMARY KEY ' ELSE N'UNIQUE ' END
| 				+ CASE i.type WHEN 1 THEN N'CLUSTERED' ELSE N'NONCLUSTERED' END
| 				+ N' (' + ISNULL(kcols.cols, N'') + N');' + CHAR(10)
| 			AS NVARCHAR(MAX)),
| 			N''
| 		) WITHIN GROUP (ORDER BY kc.is_system_named, CASE WHEN kc.is_system_named = 0 THEN kc.name ELSE N'' END, kc.type, i.index_id) AS text
| 	FROM sys.key_constraints kc
| 	JOIN sys.tables t ON t.object_id = kc.parent_object_id AND t.is_ms_shipped = 0
| 	JOIN sys.indexes i ON i.object_id = kc.parent_object_id AND i.index_id = kc.unique_index_id
| 	LEFT JOIN key_cols kcols ON kcols.object_id = i.object_id AND kcols.index_id = i.index_id
| 	GROUP BY kc.parent_object_id
| ),
| tbl_cks AS (
| 	SELECT ck.parent_object_id AS object_id,
| 		STRING_AGG(
| 			CAST(
| 				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
| 				+ CASE WHEN ck.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(ck.name) + N' ' ELSE N'' END
| 				+ N'CHECK ' + ck.definition + N';' + CHAR(10)
| 			AS NVARCHAR(MAX)),
| 			N''
| 		) WITHIN GROUP (ORDER BY ck.is_system_named, CASE WHEN ck.is_system_named = 0 THEN ck.name ELSE N'' END, ck.definition) AS text
| 	FROM sys.check_constraints ck
| 	JOIN sys.tables t ON t.object_id = ck.parent_object_id AND t.is_ms_shipped = 0
| 	GROUP BY ck.parent_object_id
| ),
| tbl_fks AS (
| 	SELECT fk.parent_object_id AS object_id,
| 		STRING_AGG(
| 			CAST(
| 				N'ALTER TABLE ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name) + N' ADD '
| 				+ CASE WHEN fk.is_system_named = 0 THEN N'CONSTRAINT ' + QUOTENAME(fk.name) + N' ' ELSE N'' END
| 				+ N'FOREIGN KEY (' + ISNULL(fpc.cols, N'') + N') REFERENCES '
| 				+ QUOTENAME(OBJECT_SCHEMA_NAME(fk.referenced_object_id)) + N'.' + QUOTENAME(OBJECT_NAME(fk.referenced_object_id))
| 				+ N' (' + ISNULL(frc.cols, N'') + N')'
| 				+ CASE fk.delete_referential_action WHEN 1 THEN N' ON DELETE CASCADE' WHEN 2 THEN N' ON DELETE SET NULL' WHEN 3 THEN N' ON DELETE SET DEFAULT' ELSE N'' END
| 				+ CASE fk.update_referential_action WHEN 1 THEN N' ON UPDATE CASCADE' WHEN 2 THEN N' ON UPDATE SET NULL' WHEN 3 THEN N' ON UPDATE SET DEFAULT' ELSE N'' END
| 				+ N';' + CHAR(10)
| 			AS NVARCHAR(MAX)),
| 			N''
| 		) WITHIN GROUP (ORDER BY fk.is_system_named, CASE WHEN fk.is_system_named = 0 THEN fk.name ELSE N'' END, fk.object_id) AS text
| 	FROM sys.foreign_keys fk
| 	JOIN sys.tables t ON t.object_id = fk.parent_object_id AND t.is_ms_shipped = 0
| 	LEFT JOIN fk_parent_cols fpc ON fpc.constraint_object_id = fk.object_id
| 	LEFT JOIN fk_ref_cols frc ON frc.constraint_object_id = fk.object_id
| 	GROUP BY fk.parent_object_id
| ),
| tbl_idxs AS (
| 	SELECT i.object_id,
| 		STRING_AGG(
| 			CAST(
| 				N'CREATE '
| 				+ CASE WHEN i.is_unique = 1 THEN N'UNIQUE ' ELSE N'' END
| 				+ CASE i.type
| 					WHEN 1 THEN N'CLUSTERED '
| 					WHEN 2 THEN N'NONCLUSTERED '
| 					WHEN 5 THEN N'CLUSTERED COLUMNSTORE '
| 					WHEN 6 THEN N'NONCLUSTERED COLUMNSTORE '
| 					ELSE N''
| 				  END
| 				+ N'INDEX ' + QUOTENAME(i.name) + N' ON ' + QUOTENAME(SCHEMA_NAME(t.schema_id)) + N'.' + QUOTENAME(t.name)
| 				+ CASE
| 					WHEN i.type = 5 THEN N''
| 					WHEN i.type = 6 THEN N' (' + ISNULL(acols.cols, N'') + N')'
| 					ELSE N' (' + ISNULL(kcols.cols, N'') + N')' + ISNULL(N' INCLUDE (' + icols.cols + N')', N'')
| 				  END
| 				+ CASE WHEN i.has_filter = 1 AND i.filter_definition IS NOT NULL THEN N' WHERE ' + i.filter_definition ELSE N'' END
| 				+ N';' + CHAR(10)
| 			AS NVARCHAR(MAX)),
| 			N''
| 		) WITHIN GROUP (ORDER BY i.name) AS text
| 	FROM sys.indexes i
| 	JOIN sys.tables t ON t.object_id = i.object_id AND t.is_ms_shipped = 0
| 	LEFT JOIN key_cols kcols ON kcols.object_id = i.object_id AND kcols.index_id = i.index_id
| 	LEFT JOIN inc_cols icols ON icols.object_id = i.object_id AND icols.index_id = i.index_id
| 	LEFT JOIN all_idx_cols acols ON acols.object_id = i.object_id AND acols.index_id = i.index_id
| 	WHERE i.type > 0
| 	  AND i.is_primary_key = 0
| 	  AND i.is_unique_constraint = 0
| 	  AND NOT EXISTS (
| 		  SELECT 1 FROM sys.index_columns ic3
| 		  JOIN sys.columns c3 ON c3.object_id = ic3.object_id AND c3.column_id = ic3.column_id
| 		  WHERE ic3.object_id = i.object_id AND ic3.index_id = i.index_id AND (c3.is_hidden = 1 OR c3.graph_type IS NOT NULL)
| 	  )
| 	GROUP BY i.object_id
| ),
| obj_trgs AS (
| 	SELECT tr.parent_id AS object_id,
| 		STRING_AGG(CAST(N'GO' + CHAR(10) + LTRIM(RTRIM(OBJECT_DEFINITION(tr.object_id))) + CHAR(10) AS NVARCHAR(MAX)), N'')
| 			WITHIN GROUP (ORDER BY tr.name) AS text
| 	FROM sys.triggers tr
| 	WHERE tr.parent_class = 1
| 	GROUP BY tr.parent_id
| )
| SELECT kind, name, text FROM (
| 	SELECT CAST(N'SCHEMA' AS NVARCHAR(32)) AS kind,
| 		CAST(s.name AS NVARCHAR(260)) AS name,
| 		CAST(N'CREATE SCHEMA ' + QUOTENAME(s.name) + N';' + CHAR(10) AS NVARCHAR(MAX)) AS text
| 	FROM sys.schemas s
| 	WHERE s.schema_id > 4 AND s.schema_id < 16384
| 	UNION ALL
| 	SELECT CAST(N'TABLE' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + t.name AS NVARCHAR(260)),
| 		CAST(
| 			CASE
| 				WHEN t.is_edge = 1 AND tc.text IS NULL
| 					THEN N'CREATE TABLE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(t.name) + N' AS EDGE;' + CHAR(10)
| 				ELSE
| 					N'CREATE TABLE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(t.name) + N' (' + ISNULL(tc.text, N'') + N')'
| 					+ CASE WHEN t.is_edge = 1 THEN N' AS EDGE' ELSE N'' END + N';' + CHAR(10)
| 			END
| 			+ ISNULL(tk.text, N'')
| 			+ ISNULL(tck.text, N'')
| 			+ ISNULL(tfk.text, N'')
| 			+ ISNULL(ti.text, N'')
| 			+ ISNULL(ttr.text, N'')
| 		AS NVARCHAR(MAX))
| 	FROM sys.tables t
| 	JOIN sys.schemas s ON s.schema_id = t.schema_id
| 	LEFT JOIN tbl_cols tc ON tc.object_id = t.object_id
| 	LEFT JOIN tbl_kcs tk ON tk.object_id = t.object_id
| 	LEFT JOIN tbl_cks tck ON tck.object_id = t.object_id
| 	LEFT JOIN tbl_fks tfk ON tfk.object_id = t.object_id
| 	LEFT JOIN tbl_idxs ti ON ti.object_id = t.object_id
| 	LEFT JOIN obj_trgs ttr ON ttr.object_id = t.object_id
| 	WHERE t.is_ms_shipped = 0
| 	UNION ALL
| 	SELECT CAST(N'VIEW' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + v.name AS NVARCHAR(260)),
| 		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(v.object_id))) + CHAR(10) + ISNULL(vtr.text, N'') AS NVARCHAR(MAX))
| 	FROM sys.views v
| 	JOIN sys.schemas s ON s.schema_id = v.schema_id
| 	LEFT JOIN obj_trgs vtr ON vtr.object_id = v.object_id
| 	WHERE v.is_ms_shipped = 0
| 	UNION ALL
| 	SELECT CAST(N'SEQUENCE' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + sq.name AS NVARCHAR(260)),
| 		CAST(
| 			N'CREATE SEQUENCE ' + QUOTENAME(s.name) + N'.' + QUOTENAME(sq.name)
| 			+ N' AS ' + ty.name
| 			+ N' START WITH ' + CAST(sq.start_value AS NVARCHAR(40))
| 			+ N' INCREMENT BY ' + CAST(sq.increment AS NVARCHAR(40))
| 			+ N' MINVALUE ' + CAST(sq.minimum_value AS NVARCHAR(40))
| 			+ N' MAXVALUE ' + CAST(sq.maximum_value AS NVARCHAR(40))
| 			+ CASE WHEN sq.is_cycling = 1 THEN N' CYCLE' ELSE N' NO CYCLE' END
| 			+ CASE WHEN sq.is_cached = 0 THEN N' NO CACHE' WHEN sq.cache_size IS NOT NULL THEN N' CACHE ' + CAST(sq.cache_size AS NVARCHAR(40)) ELSE N'' END
| 			+ N';' + CHAR(10)
| 		AS NVARCHAR(MAX))
| 	FROM sys.sequences sq
| 	JOIN sys.schemas s ON s.schema_id = sq.schema_id
| 	JOIN sys.types ty ON ty.user_type_id = sq.user_type_id
| 	UNION ALL
| 	SELECT CAST(N'SYNONYM' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + sy.name AS NVARCHAR(260)),
| 		CAST(N'CREATE SYNONYM ' + QUOTENAME(s.name) + N'.' + QUOTENAME(sy.name) + N' FOR ' + sy.base_object_name + N';' + CHAR(10) AS NVARCHAR(MAX))
| 	FROM sys.synonyms sy
| 	JOIN sys.schemas s ON s.schema_id = sy.schema_id
| 	UNION ALL
| 	SELECT CAST(N'PROCEDURE' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + p.name AS NVARCHAR(260)),
| 		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(p.object_id))) + CHAR(10) AS NVARCHAR(MAX))
| 	FROM sys.procedures p
| 	JOIN sys.schemas s ON s.schema_id = p.schema_id
| 	WHERE p.is_ms_shipped = 0 AND NOT (s.name = N'dbo' AND p.name IN (N'_sim_reset', N'_sim_forms'))
| 	UNION ALL
| 	SELECT CAST(N'FUNCTION' AS NVARCHAR(32)),
| 		CAST(s.name + N'.' + o.name AS NVARCHAR(260)),
| 		CAST(LTRIM(RTRIM(OBJECT_DEFINITION(o.object_id))) + CHAR(10) AS NVARCHAR(MAX))
| 	FROM sys.objects o
| 	JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	WHERE o.type IN ('FN', 'IF', 'TF') AND o.is_ms_shipped = 0
| ) o ORDER BY name COLLATE Latin1_General_BIN2, kind;
| END;
< ok 0 -
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id, name FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT id, name FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @id INT AS SELECT name FROM customers WHERE id = @id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p @id INT AS SELECT name FROM customers WHERE id = @id
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @id INT = NULL, @name NVARCHAR(100) OUTPUT AS SELECT @name = name FROM customers WHERE id = @id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p @id INT = NULL, @name NVARCHAR(100) OUTPUT AS SELECT @name = name FROM customers WHERE id = @id
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT nope FROM customers
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT nope FROM missing_t
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT nope FROM missing_t
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT c.nope FROM customers c JOIN missing_t m ON m.id = c.id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT c.nope FROM customers c JOIN missing_t m ON m.id = c.id
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN SELECT nope FROM missing_t; SELECT nope FROM customers; END
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN SET NOCOUNT ON; INSERT INTO orders (customer_id, amount) VALUES (1, 2.5); SELECT SCOPE_IDENTITY() AS id; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN SET NOCOUNT ON; INSERT INTO orders (customer_id, amount) VALUES (1, 2.5); SELECT SCOPE_IDENTITY() AS id; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS INSERT INTO orders (nope) VALUES (1)
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS INSERT INTO orders (customer_id) VALUES (1, 2)
|
< error
| 110 1 15 1 proc=p There are fewer columns in the INSERT statement than values specified in the VALUES clause. The number of values in the VALUES clause must match the number of columns specified in the INSERT statement.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS INSERT INTO orders (customer_id, amount) SELECT id FROM customers
|
< error
| 120 1 15 1 proc=p The select list for the INSERT statement contains fewer items than the insert list. The number of SELECT values must match the number of INSERT columns.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS INSERT INTO missing_t (x) VALUES (1)
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS INSERT INTO missing_t (x) VALUES (1)
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS UPDATE orders SET status = 'x' WHERE id = 1
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS UPDATE orders SET status = 'x' WHERE id = 1
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS UPDATE orders SET nope = 1
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS UPDATE o SET status = c.region FROM orders o JOIN customers c ON c.id = o.customer_id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS UPDATE o SET status = c.region FROM orders o JOIN customers c ON c.id = o.customer_id
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS DELETE FROM orders WHERE nope = 1
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS DELETE o FROM orders o JOIN customers c ON c.id = o.customer_id WHERE c.region = 'EU'
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS DELETE o FROM orders o JOIN customers c ON c.id = o.customer_id WHERE c.region = 'EU'
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS MERGE orders AS t USING customers AS s ON t.customer_id = s.id WHEN MATCHED THEN UPDATE SET t.status = 'x' WHEN NOT MATCHED THEN INSERT (customer_id) VALUES (s.id);
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS MERGE orders AS t USING customers AS s ON t.customer_id = s.id WHEN MATCHED THEN UPDATE SET t.status = 'x' WHEN NOT MATCHED THEN INSERT (customer_id) VALUES (s.id);
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS MERGE orders AS t USING customers AS s ON t.customer_id = s.nope WHEN MATCHED THEN DELETE;
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE @n INT; SET @n = (SELECT COUNT(*) FROM orders); IF @n > 10 SELECT 'many' AS r ELSE SELECT 'few' AS r; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN DECLARE @n INT; SET @n = (SELECT COUNT(*) FROM orders); IF @n > 10 SELECT 'many' AS r ELSE SELECT 'few' AS r; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE @i INT = 0; WHILE @i < 10 BEGIN SET @i += 1; END; RETURN @i; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN DECLARE @i INT = 0; WHILE @i < 10 BEGIN SET @i += 1; END; RETURN @i; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT @undeclared AS v
|
< error
| 137 2 15 1 proc=p Must declare the scalar variable "@undeclared".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN SET @x = 1; END
|
< error
| 137 1 15 1 proc=p Must declare the scalar variable "@x".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN SELECT @n AS v; DECLARE @n INT; END
|
< error
| 137 2 15 1 proc=p Must declare the scalar variable "@n".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN TRY SELECT 1 / 0 AS x; END TRY BEGIN CATCH SELECT ERROR_MESSAGE() AS m; THROW; END CATCH
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN TRY SELECT 1 / 0 AS x; END TRY BEGIN CATCH SELECT ERROR_MESSAGE() AS m; THROW; END CATCH
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); INSERT INTO #t VALUES (1); SELECT x FROM #t; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); INSERT INTO #t VALUES (1); SELECT x FROM #t; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); SELECT y FROM #t; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN CREATE TABLE #t (x INT); SELECT y FROM #t; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); SELECT x FROM @t; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); SELECT x FROM @t; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE @t TABLE (x INT); SELECT y FROM @t; END
|
< error
| 207 1 16 1 proc=p Invalid column name 'y'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT x FROM @t
|
< error
| 1087 2 15 1 proc=p Must declare the table variable "@t".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS EXEC p_existing 1
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS EXEC p_existing 1
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS EXEC nope_proc
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS EXEC nope_proc
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT dbo.f_double(id) AS d FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT dbo.f_double(id) AS d FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT dbo.nope(id) AS d FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT dbo.nope(id) AS d FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT dbo.f_double(id, 1) AS d FROM customers
|
< error
| 8144 2 16 1 proc=p Procedure or function dbo.f_double has too many arguments specified.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id FROM dbo.f_orders(1)
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT id FROM dbo.f_orders(1)
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id FROM dbo.f_orders(1, 2)
|
< error
| 8144 3 16 1 proc=p Procedure or function dbo.f_orders has too many arguments specified.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT * FROM v_customers WHERE nope = 1
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @id INT AS SELECT name FROM customers WHERE id = @id ORDER BY nope
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id INTO #copy FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT id INTO #copy FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY customer_id
|
< error
| 8120 1 16 1 proc=p Column 'orders.status' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id FROM customers UNION SELECT id, status FROM orders
|
< error
| 205 1 16 1 proc=p All queries combined using a UNION, INTERSECT or EXCEPT operator must have an equal number of expressions in their target lists.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id FROM customers c JOIN orders o ON o.customer_id = c.id
|
< error
| 209 1 16 1 proc=p Ambiguous column name 'id'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @c CURSOR VARYING OUTPUT AS BEGIN SET @c = CURSOR FOR SELECT id FROM customers; OPEN @c; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p @c CURSOR VARYING OUTPUT AS BEGIN SET @c = CURSOR FOR SELECT id FROM customers; OPEN @c; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE c CURSOR FOR SELECT nope FROM customers; OPEN c; END
|
< error
| 207 1 16 1 proc=p Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS WITH x AS (SELECT id FROM customers) SELECT id FROM x
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS WITH x AS (SELECT id FROM customers) SELECT id FROM x
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS TRUNCATE TABLE audit
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS TRUNCATE TABLE audit
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN BEGIN TRANSACTION; UPDATE orders SET amount = amount * 2; COMMIT; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN BEGIN TRANSACTION; UPDATE orders SET amount = amount * 2; COMMIT; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @n INT AS SELECT TOP (@n) id FROM orders ORDER BY id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p @n INT AS SELECT TOP (@n) id FROM orders ORDER BY id
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT id FROM customers WHERE id = @@ROWCOUNT
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT id FROM customers WHERE id = @@ROWCOUNT
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS RAISERROR('x', 16, 1)
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS RAISERROR('x', 16, 1)
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS BEGIN DECLARE @sql NVARCHAR(MAX) = N'SELECT nope FROM nowhere'; EXEC sp_executesql @sql; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS BEGIN DECLARE @sql NVARCHAR(MAX) = N'SELECT nope FROM nowhere'; EXEC sp_executesql @sql; END
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @id INT AS IF EXISTS (SELECT 1 FROM customers WHERE id = @id) RETURN 1 ELSE RETURN 0
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p @id INT AS IF EXISTS (SELECT 1 FROM customers WHERE id = @id) RETURN 1 ELSE RETURN 0
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p WITH RECOMPILE AS SELECT id FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p WITH RECOMPILE AS SELECT id FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p WITH EXECUTE AS OWNER AS SELECT id FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p WITH EXECUTE AS OWNER AS SELECT id FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE s.p AS SELECT id FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
< row
: string
| PROCEDURE
: string
| s.p
: string
| CREATE PROCEDURE s.p AS SELECT id FROM customers
|
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE nope.p AS SELECT 1
|
< error
| 2760 1 16 1 proc=p The specified schema name "nope" either does not exist or you do not have permission to use it.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p_existing AS SELECT 1
|
< error
| 2714 3 16 1 proc=p_existing There is already an object named 'p_existing' in the database.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @id NOPE AS SELECT 1
|
< error
| 2715 3 16 1 proc=p Column, parameter, or variable #1: Cannot find data type NOPE.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER PROCEDURE p_existing @id INT, @extra INT = 0 AS SELECT name FROM customers WHERE id = @id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE p_existing @id INT, @extra INT = 0 AS SELECT name FROM customers WHERE id = @id
|
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER PROCEDURE p_existing AS SELECT nope FROM customers
|
< error
| 207 1 16 1 proc=p_existing Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER PROCEDURE nope AS SELECT 1
|
< error
| 208 6 16 1 proc=nope Invalid object name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE OR ALTER PROCEDURE p_existing AS SELECT id FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE   PROCEDURE p_existing AS SELECT id FROM customers
|
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE OR ALTER PROCEDURE p AS SELECT id FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE   PROCEDURE p AS SELECT id FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP PROCEDURE p_existing
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP PROCEDURE IF EXISTS p_existing, nope
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@a INT, @b INT) RETURNS INT AS BEGIN RETURN @a + @b END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@a INT, @b INT) RETURNS INT AS BEGIN RETURN @a + @b END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS DECIMAL(10,2) AS BEGIN DECLARE @s DECIMAL(10,2); SELECT @s = SUM(amount) FROM orders WHERE customer_id = @c; RETURN @s; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS DECIMAL(10,2) AS BEGIN DECLARE @s DECIMAL(10,2); SELECT @s = SUM(amount) FROM orders WHERE customer_id = @c; RETURN @s; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS INT AS BEGIN RETURN (SELECT COUNT(*) FROM orders WHERE nope = @c) END
|
< error
| 207 1 16 1 proc=f Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT COUNT(*) FROM missing_t) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT COUNT(*) FROM missing_t) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN INSERT INTO orders (customer_id) VALUES (1); RETURN 1 END
|
< error
| 443 15 16 1 proc=f Invalid use of a side-effecting operator 'INSERT' within a function.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS UNIQUEIDENTIFIER AS BEGIN RETURN NEWID() END
|
< error
| 443 1 16 1 proc=f Invalid use of a side-effecting operator 'newid' within a function.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN EXEC p_existing 1; RETURN 1 END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN EXEC p_existing 1; RETURN 1 END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN @x END
|
< error
| 137 2 15 1 proc=f Must declare the scalar variable "@x".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN SELECT id FROM customers; RETURN 1 END
|
< error
| 444 2 16 1 proc=f Select statements included within a function cannot return data to a client.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); RETURN (SELECT MAX(x) FROM @t) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN DECLARE @t TABLE (x INT); INSERT INTO @t VALUES (1); RETURN (SELECT MAX(x) FROM @t) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM dbo.orders) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM dbo.orders) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 2 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT WITH SCHEMABINDING AS BEGIN RETURN (SELECT COUNT(*) FROM orders) END
|
< error
| 4512 3 16 1 proc=f Cannot schema bind function 'f' because name 'orders' is invalid for schema binding. Names must be in two-part format and an object cannot reference itself.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN IF @x > 0 RETURN 1; RETURN 0; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN IF @x > 0 RETURN 1; RETURN 0; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN IF @x > 0 RETURN 1; END
|
< error
| 455 2 16 1 proc=f The last statement included within a function must be a return statement.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.f_double(2) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.f_double(2) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.nope(2) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN dbo.nope(2) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT TOP 1 id FROM dbo.f_orders(1)) END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS INT AS BEGIN RETURN (SELECT TOP 1 id FROM dbo.f_orders(1)) END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS NOPE AS BEGIN RETURN 1 END
|
< error
| 2715 3 16 1 proc=f Column, parameter, or variable #0: Cannot find data type NOPE.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN DECLARE @y INT; SET @y = @x * 2; RETURN @y END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@x INT) RETURNS INT AS BEGIN DECLARE @y INT; SET @y = @x * 2; RETURN @y END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@x INT) RETURNS INT WITH RETURNS NULL ON NULL INPUT AS BEGIN RETURN @x END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@x INT) RETURNS INT WITH RETURNS NULL ON NULL INPUT AS BEGIN RETURN @x END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x END
|
< error
| 2714 3 16 1 proc=f_double There is already an object named 'f_double' in the database.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x * 3 END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x * 3 END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER FUNCTION f_double(@x INT) RETURNS TABLE AS RETURN SELECT @x AS x
|
< error
| 2010 1 16 1 proc=f_double Cannot perform alter on 'f_double' because it is an incompatible object type.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE OR ALTER FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x + 1 END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE   FUNCTION f_double(@x INT) RETURNS INT AS BEGIN RETURN @x + 1 END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP FUNCTION f_double
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT id, amount FROM orders WHERE customer_id = @c
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT id, amount FROM orders WHERE customer_id = @c
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN (SELECT id, amount FROM orders WHERE customer_id = @c)
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN (SELECT id, amount FROM orders WHERE customer_id = @c)
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT nope FROM orders
|
< error
| 207 1 16 1 proc=f Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM missing_t
|
< error
| 208 1 16 1 proc=f Invalid object name 'missing_t'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT COUNT(*) FROM orders
|
< error
| 4514 1 16 1 proc=f CREATE FUNCTION failed because a column name is not specified for column 1.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id, id FROM orders
|
< error
| 4506 1 16 1 proc=f Column names in each view or function must be unique. Column name 'id' in view or function 'f' is specified more than once.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM orders ORDER BY id
|
< error
| 1033 1 15 1 proc=f The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT TOP 5 id FROM orders ORDER BY id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT TOP 5 id FROM orders ORDER BY id
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN WITH x AS (SELECT id FROM orders) SELECT id FROM x
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS TABLE AS RETURN WITH x AS (SELECT id FROM orders) SELECT id FROM x
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT @x AS v
|
< error
| 137 2 15 1 proc=f Must declare the scalar variable "@x".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS TABLE WITH SCHEMABINDING AS RETURN SELECT id FROM dbo.orders WHERE customer_id = @c
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS TABLE WITH SCHEMABINDING AS RETURN SELECT id FROM dbo.orders WHERE customer_id = @c
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 2 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS TABLE WITH SCHEMABINDING AS RETURN SELECT id FROM orders WHERE customer_id = @c
|
< error
| 4512 3 16 1 proc=f Cannot schema bind table valued function 'f' because name 'orders' is invalid for schema binding. Names must be in two-part format and an object cannot reference itself.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT o.id, x.id AS xid FROM orders o CROSS APPLY dbo.f_orders(@c) x
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS TABLE AS RETURN SELECT o.id, x.id AS xid FROM orders o CROSS APPLY dbo.f_orders(@c) x
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM v_customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT id FROM v_customers
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER FUNCTION f_orders(@c INT) RETURNS TABLE AS RETURN SELECT id FROM orders WHERE customer_id = @c
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION f_orders(@c INT) RETURNS TABLE AS RETURN SELECT id FROM orders WHERE customer_id = @c
|
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER FUNCTION f_orders(@c INT) RETURNS TABLE AS RETURN SELECT nope FROM orders
|
< error
| 207 1 16 1 proc=f_orders Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP FUNCTION f_orders, f_double
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f(@c INT) RETURNS @r TABLE (id INT, amount DECIMAL(10,2)) AS BEGIN INSERT INTO @r SELECT id, amount FROM orders WHERE customer_id = @c; RETURN; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f(@c INT) RETURNS @r TABLE (id INT, amount DECIMAL(10,2)) AS BEGIN INSERT INTO @r SELECT id, amount FROM orders WHERE customer_id = @c; RETURN; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r SELECT nope FROM orders; RETURN; END
|
< error
| 207 1 16 1 proc=f Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r (nope) VALUES (1); RETURN; END
|
< error
| 207 1 16 1 proc=f Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT, id INT) AS BEGIN RETURN; END
|
< error
| 2705 3 16 1 proc=f Column names in each table must be unique. Column name 'id' in table '@r' is specified more than once.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT PRIMARY KEY, v NVARCHAR(10) DEFAULT 'x') AS BEGIN INSERT INTO @r (id) VALUES (1); RETURN; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS @r TABLE (id INT PRIMARY KEY, v NVARCHAR(10) DEFAULT 'x') AS BEGIN INSERT INTO @r (id) VALUES (1); RETURN; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r VALUES (1); RETURN 1; END
|
< error
| 178 1 15 1 proc=f A RETURN statement with a return value cannot be used in this context.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r SELECT id FROM missing_t; RETURN; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN INSERT INTO @r SELECT id FROM missing_t; RETURN; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id NOPE) AS BEGIN RETURN; END
|
< error
| 2715 6 16 1 proc=f Column, parameter, or variable #1: Cannot find data type NOPE.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE @r SET id = id + 1; DELETE FROM @r WHERE id > 5; RETURN; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f
: string
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE @r SET id = id + 1; DELETE FROM @r WHERE id > 5; RETURN; END
|
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (id INT) AS BEGIN UPDATE orders SET amount = 0; RETURN; END
|
< error
| 443 15 16 1 proc=f Invalid use of a side-effecting operator 'UPDATE' within a function.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS INSERT INTO audit (what) SELECT 'order ' + CAST(id AS NVARCHAR(10)) FROM inserted
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT AS INSERT INTO audit (what) SELECT 'order ' + CAST(id AS NVARCHAR(10)) FROM inserted
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER UPDATE AS IF UPDATE(amount) INSERT INTO audit (what) SELECT 'amount' FROM inserted i JOIN deleted d ON d.id = i.id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER UPDATE AS IF UPDATE(amount) INSERT INTO audit (what) SELECT 'amount' FROM inserted i JOIN deleted d ON d.id = i.id
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER DELETE AS SELECT nope FROM deleted
|
< error
| 207 1 16 1 proc=t Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT nope FROM inserted
|
< error
| 207 1 16 1 proc=t Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS IF UPDATE(nope) RETURN
|
< error
| 207 1 16 1 proc=t Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders INSTEAD OF DELETE AS UPDATE o SET status = 'deleted' FROM orders o JOIN deleted d ON d.id = o.id
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders INSTEAD OF DELETE AS UPDATE o SET status = 'deleted' FROM orders o JOIN deleted d ON d.id = o.id
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON v_customers INSTEAD OF INSERT AS INSERT INTO customers (id, name) SELECT id, name FROM inserted
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
| GO
| CREATE TRIGGER t ON v_customers INSTEAD OF INSERT AS INSERT INTO customers (id, name) SELECT id, name FROM inserted
|
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON v_customers INSTEAD OF INSERT AS INSERT INTO customers (id, name) SELECT id, region FROM inserted
|
< error
| 207 1 16 1 proc=t Invalid column name 'region'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON v_customers AFTER INSERT AS SET NOCOUNT ON
|
< error
| 8197 6 16 1 proc=t The object 'v_customers' does not exist or is invalid for this operation.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders FOR INSERT, UPDATE AS SET NOCOUNT ON
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders FOR INSERT, UPDATE AS SET NOCOUNT ON
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS BEGIN DECLARE @n INT = (SELECT COUNT(*) FROM inserted); IF @n > 1 THROW 50000, 'one at a time', 1; END
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT AS BEGIN DECLARE @n INT = (SELECT COUNT(*) FROM inserted); IF @n > 1 THROW 50000, 'one at a time', 1; END
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT COLUMNS_UPDATED() AS c
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT COLUMNS_UPDATED() AS c
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT id FROM missing_t
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT id FROM missing_t
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS SELECT @v AS v
|
< error
| 137 2 15 1 proc=t Must declare the scalar variable "@v".
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders WITH EXECUTE AS OWNER AFTER INSERT AS SET NOCOUNT ON
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders WITH EXECUTE AS OWNER AFTER INSERT AS SET NOCOUNT ON
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT NOT FOR REPLICATION AS SET NOCOUNT ON
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT NOT FOR REPLICATION AS SET NOCOUNT ON
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON orders AFTER INSERT AS ROLLBACK TRANSACTION
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER t ON orders AFTER INSERT AS ROLLBACK TRANSACTION
|
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON missing_t AFTER INSERT AS SET NOCOUNT ON
|
< error
| 8197 4 16 1 proc=t The object 'missing_t' does not exist or is invalid for this operation.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER s.t ON orders AFTER INSERT AS SET NOCOUNT ON
|
< error
| 2103 1 15 1 proc=t Cannot create trigger 's.t' because its schema is different from the schema of the target table or view.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER trg_existing ON orders AFTER INSERT AS SET NOCOUNT ON
|
< error
| 2714 2 16 1 proc=trg_existing There is already an object named 'trg_existing' in the database.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER trg_existing ON customers AFTER INSERT AS SET NOCOUNT ON
|
< error
| 2714 2 16 1 proc=trg_existing There is already an object named 'trg_existing' in the database.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER TRIGGER trg_existing ON orders AFTER UPDATE AS SET NOCOUNT ON
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER trg_existing ON orders AFTER UPDATE AS SET NOCOUNT ON
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| ALTER TRIGGER trg_existing ON customers AFTER UPDATE AS SET NOCOUNT ON
|
< error
| 2110 1 15 1 proc=trg_existing Cannot alter trigger 'trg_existing' on 'customers' because this trigger does not belong to this object. Specify the correct trigger name or the correct target object name.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE OR ALTER TRIGGER trg_existing ON orders AFTER DELETE AS SET NOCOUNT ON
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE   TRIGGER trg_existing ON orders AFTER DELETE AS SET NOCOUNT ON
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP TRIGGER trg_existing
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP TRIGGER IF EXISTS nope
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS PRINT 'created'
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE FOR DDL_TABLE_EVENTS AS SELECT EVENTDATA().value('(/EVENT_INSTANCE/ObjectName)[1]', 'nvarchar(128)') AS n
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE FOR NOPE_EVENT AS PRINT 'x'
|
< error
| 1084 1 15 1 proc=t 'NOPE_EVENT' is an invalid event type.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE AFTER DROP_TABLE AS ROLLBACK
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS SELECT nope FROM customers
|
< error
| 207 1 16 1 proc=t Invalid column name 'nope'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER t ON DATABASE FOR CREATE_TABLE AS INSERT INTO audit (what) VALUES (EVENTDATA().value('(/EVENT_INSTANCE/ObjectName)[1]', 'nvarchar(100)'))
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| DROP TRIGGER nope ON DATABASE
|
< error
| 3701 5 11 1 Cannot drop the trigger 'nope', because it does not exist or you do not have permission.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS CREATE TABLE #t (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 proc=p Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS DECLARE @t TABLE (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 proc=p Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS ALTER TABLE missing_t ADD c NVARCHAR(10) COLLATE a
|
< error
| 448 2 16 1 proc=p Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT name COLLATE a FROM customers
|
< error
| 448 1 16 1 proc=p Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT name COLLATE latin1_general_ci_as FROM customers
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT name COLLATE latin1_general_ci_as FROM customers
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p AS SELECT 'x' COLLATE a FROM missing_t
|
< ok 0 -
> query
| EXEC [dbo].[_sim_forms]
< result
< column "kind" "NVARCHAR"
< column "name" "NVARCHAR"
< column "text" "NVARCHAR"
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
|
< row
: string
| TABLE
: string
| dbo.audit
: string
| CREATE TABLE [dbo].[audit] ([id] int IDENTITY(1,1) NOT NULL, [what] nvarchar(100) NULL, [at] datetime2(3) NULL);
| ALTER TABLE [dbo].[audit] ADD CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL);
| ALTER TABLE [dbo].[customers] ADD CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| FUNCTION
: string
| dbo.f_double
: string
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
|
< row
: string
| FUNCTION
: string
| dbo.f_orders
: string
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
|
< row
: string
| TABLE
: string
| dbo.orders
: string
| CREATE TABLE [dbo].[orders] ([id] int IDENTITY(1,1) NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
| GO
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
|
< row
: string
| PROCEDURE
: string
| dbo.p
: string
| CREATE PROCEDURE p AS SELECT 'x' COLLATE a FROM missing_t
|
|
< row
: string
| PROCEDURE
: string
| dbo.p_existing
: string
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
|
< row
: string
| SCHEMA
: string
| s
: string
| CREATE SCHEMA [s];
|
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE PROCEDURE p @x NVARCHAR(10) AS SELECT @x COLLATE a
|
< error
| 448 1 16 1 proc=p Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS NVARCHAR(10) AS BEGIN RETURN N'x' COLLATE a END
|
< error
| 448 1 16 1 proc=f Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT name COLLATE a AS n FROM customers
|
< error
| 448 1 16 1 proc=f Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS @r TABLE (c NVARCHAR(10) COLLATE a) AS BEGIN RETURN END
|
< error
| 448 2 16 1 proc=f Invalid collation 'a'.
> exec
| DECLARE @sql NVARCHAR(MAX) = N'';
| SELECT @sql += N'DROP ' + CASE o.type WHEN 'V' THEN N'VIEW ' ELSE N'FUNCTION ' END + QUOTENAME(SCHEMA_NAME(o.schema_id)) + N'.' + QUOTENAME(o.name) + N';' + CHAR(10)
| FROM sys.objects o JOIN sys.sql_modules m ON m.object_id = o.object_id
| WHERE m.is_schema_bound = 1 AND o.is_ms_shipped = 0 AND o.type IN ('FN', 'IF', 'TF', 'V')
| ORDER BY CASE o.type WHEN 'V' THEN 0 ELSE 1 END, o.name;
| IF @sql <> N'' EXEC sp_executesql @sql;
< ok 1 -
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[audit] ([id] INT IDENTITY(1,1) NOT NULL, [what] NVARCHAR(100) NULL, [at] DATETIME2(3) NULL, CONSTRAINT [PK_audit] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT IDENTITY(1,1) NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE PROCEDURE [dbo].[p_existing] @id INT AS SELECT [name] FROM [dbo].[customers] WHERE [id] = @id;
< ok 0 -
> exec
| CREATE TRIGGER [dbo].[trg_existing] ON [dbo].[orders] AFTER INSERT AS SET NOCOUNT ON;
< ok 0 -
> exec
| CREATE TRIGGER tr ON customers AFTER INSERT AS SELECT name COLLATE a FROM inserted
|
< error
| 448 1 16 1 proc=tr Invalid collation 'a'.
