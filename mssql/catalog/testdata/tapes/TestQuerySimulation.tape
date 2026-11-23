-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestQuerySimulation default#1
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
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS WITH CHANGE_TRACKING_CONTEXT(0x00) SELECT *
|
< error
| 22108 1 16 1 proc=_v The CHANGE_TRACKING_CONTEXT WITH clause cannot be used with a SELECT statement.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS WITH XMLNAMESPACES(DEFAULT 'a') SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS WITH a AS(SELECT *) SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'a'.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS WITH a(a) AS(SELECT *) SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS WITH XMLNAMESPACES('a' AS a) SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 OFFSET 1 ROW
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FOR BROWSE
|
< error
| 176 1 15 1 proc=_v The FOR BROWSE clause is no longer supported in views.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(SELECT *) a
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(VALUES(1)) a
|
< error
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(SELECT *) FOR PATH a
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(SELECT *) a(a)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(VALUES(1),(1)) a
|
< error
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT *) zzz
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * UNION SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * EXCEPT SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * INTERSECT SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * UNION ALL SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT((SELECT *))
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * UNION SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * SELECT *
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'SELECT'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * EXCEPT SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * INTERSECT SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * UNION ALL SELECT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS (SELECT *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT ALL *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * FROM @a)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT TOP 1 *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT *, *)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
| 116 1 16 1 proc=_v Only one expression can be specified in the select list when the subquery is not introduced with EXISTS.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * WHERE 1 = 1)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * GROUP BY 1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * HAVING 1 = 1)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * WINDOW a AS())
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT * ORDER BY 1 OFFSET 1 ROW)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT *, *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * INTO a
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'INTO'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WHERE 1 = 1
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY 1
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * HAVING 1 = 1
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WINDOW a AS()
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ALL *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DISTINCT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 OFFSET 1 ROWS
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 OFFSET 1 ROW FETCH FIRST 1 ROW ONLY
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 OFFSET 1 ROW FETCH NEXT 1 ROW ONLY
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP(1) *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 PERCENT *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT @a = 1
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(SELECT 1 zzz)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 'a' = 1
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 'a' = 1
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 'a'
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 'a'
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 AS 'a'
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 AS 'a'
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(a)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a TABLESAMPLE(1)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a TABLESAMPLE SYSTEM(1)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a TABLESAMPLE(1 ROWS)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a TABLESAMPLE(1 PERCENT)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a TABLESAMPLE(1) REPEATABLE(1)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a zzz
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a, @a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a, @a zzz
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR PATH
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a FOR PATH
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a PATH
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a PATH
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR PATH zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a FOR PATH zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME ALL
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME ALL zzz
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME AS OF 'a'
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME AS OF 'a' zzz
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME FROM 'a' TO 'a'
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME FROM 'a' TO 'a' zzz
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME BETWEEN 'a' AND 'a'
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME BETWEEN 'a' AND 'a' zzz
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME CONTAINED IN('a', 'a')
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME CONTAINED IN('a', 'a') zzz
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a CROSS JOIN @a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a CROSS JOIN @a zzz
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME } zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME } zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN TRUNCATE(1, 1) }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN DATABASE() }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN USER() }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_DATE() }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIMESTAMP }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME() }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME(zzz) }
|
< error
| 207 1 16 1 proc=_v Invalid column name 'zzz'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT { FN CURRENT_TIME(1) }
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM((@a CROSS JOIN @a))
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(@a CROSS JOIN @a)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(@a CROSS JOIN @a zzz)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM((@a CROSS JOIN @a) CROSS JOIN @a)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM((@a CROSS JOIN @a) CROSS JOIN @a zzz)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(@a CROSS JOIN @a CROSS JOIN @a)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM(@a CROSS JOIN @a CROSS JOIN @a zzz)
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a UNPIVOT(a FOR a IN(a)) a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a UNPIVOT(a FOR a IN(a, a)) a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CHANGETABLE(a, 1) a
|
< error
| 208 1 16 1 proc=_v Invalid object name 'CHANGETABLE'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a()
|
< error
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a() zzz
|
< error
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a.a() a(a)
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM.a() a(a)
|
< error
| 317 1 16 1 proc=_v Table-valued function 'a' cannot have a column alias.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a')
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a') zzz
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTENT(a, a, 'a')
|
< error
| 208 1 16 1 proc=_v Invalid object name 'CONTENT'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a, a)
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a, a) zzz
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a')
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') zzz
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPEN_EXISTING(@a, 'a')
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a OUTER APPLY @a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a CROSS APPLY @a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a INNER JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a INPUT JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a LEFT JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a LEGACY_CARDINALITY_ESTIMATION JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a RIGHT JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a ROLE JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a RIGHT OUTER JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a FULL JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a FULLSCAN JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a FULL OUTER JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a LEFT OUTER JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a INNER MERGE JOIN @a ON 1 = 1
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a() a
|
< error
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(zzz)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'zzz'.
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(1)
|
< error
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(1, 1)
|
< error
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(FORCE_LAST_GOOD_PLAN)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'FORCE_LAST_GOOD_PLAN'.
| 215 1 16 1 proc=_v Parameters supplied for object 'a' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK) a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a(FORCESEEK) a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a a(FORCESEEK)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a a(FORCESEEK)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM @a AS a
|
< error
| 1087 2 15 1 proc=_v Must declare the table variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a') a
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a', 1)
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM FREETEXTTABLE(a, a, 'a')
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, PROPERTY(a, 'a'), 'a')
|
< error
| 7601 4 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on column 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a, *)
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a,(*))
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a,(a))
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a,(a, a))
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a', LANGUAGE 0x00)
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM CONTAINSTABLE(a, a, 'a', LANGUAGE 0x00, 1)
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'a' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYTABLE(a, a, 1)
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYTABLE(a, a, 1) zzz
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEND(a, a, 1)
|
< error
| 208 1 16 1 proc=_v Invalid object name 'SEND'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(a, a, 1, a, 1)
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(a, a, 1, a, 1) zzz
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a, a) a
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICKEYPHRASETABLE(a, a, 1)
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYTABLE(a, a, 1) a
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(a, a, 1, a, 1) a
|
< error
| 41202 1 16 1 proc=_v The source table 'a' specified in the SEMANTICSIMILARITYTABLE, SEMANTICKEYPHRASETABLE or SEMANTICSIMILARITYDETAILSTABLE function doesn't have a full-text index that uses the STATISTICAL_SEMANTICS option.  A full-text index using the STATISTICAL_SEMANTICS option is required to use this function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') a
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH a
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH a zzz
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH(a a)
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH(a a) zzz
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a', 1)
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH(a a, a a)
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENXML(@a, 'a') WITH(a a 'a')
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
| 319 1 15 1 proc=_v Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) WITH(a a)
|
< error
| 2715 14 16 1 proc=_v Column, parameter, or variable #0: Cannot find data type a.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) WITH(a a) zzz
|
< error
| 2715 14 16 1 proc=_v Column, parameter, or variable #0: Cannot find data type a.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1, 1)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1, 1)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) WITH(a a, a a)
|
< error
| 2715 14 16 1 proc=_v Column, parameter, or variable #0: Cannot find data type a.
| 2715 14 16 1 proc=_v Column, parameter, or variable #1: Cannot find data type a.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) WITH(a a 'a')
|
< error
| 2715 14 16 1 proc=_v Column, parameter, or variable #0: Cannot find data type a.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) WITH(a a AS JSON)
|
< error
| 2715 14 16 1 proc=_v Column, parameter, or variable #0: Cannot find data type a.
| 13618 1 16 1 proc=_v AS JSON option can be specified only for column of nvarchar(max) type in WITH clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a WITH(1)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a WITH(1)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a WITH(FORCESEEK)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a WITH(FORCESEEK)
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ALL 1
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY 1, 1
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY CUBE(1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY(1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY ROLLUP(1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS(1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING(1)
|
< error
| 144 1 15 1 proc=_v Cannot use an aggregate or a subquery in an expression used for the group by list of a GROUP BY clause.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY()
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS(())
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS((1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS((1, 1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS(CUBE(1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS(ROLLUP(1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY CUBE(1, 1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY CUBE((1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY CUBE((1, 1))
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY 1 WITH(DISTRIBUTED_AGG)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1, 1
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 ASC
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(USE PLAN 'a')
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
| 154 1 15 1 proc=_v a USE database statement is not allowed in a procedure, function or trigger.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(OPTIMIZE FOR UNKNOWN)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(TABLE HINT(a))
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(TABLE HINT(a, FORCESEEK))
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(MERGE JOIN)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(MERGE UNION)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(FORCE ORDER)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(HASH GROUP)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(ORDER GROUP)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(OPTIMIZE FOR(@a UNKNOWN))
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(OPTIMIZE FOR(@a UNKNOWN, @a UNKNOWN))
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * OPTION(OPTIMIZE FOR(@a = 1))
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1E1 *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1.5 *
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a zzz
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 + 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 + 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 + 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 - 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 - 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 - 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 & 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 & 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 & 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 | 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 | 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 | 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 ^ 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 ^ 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 ^ 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 * 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 * 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 * 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 / 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 / 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 / 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 % 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 % 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 % 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 << 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 << 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 << 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 >> 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 >> 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 >> 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT + 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT + 1 zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT + 1 zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT - 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ~ 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a a
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a zzz
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1.a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1.a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1.
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(1 AS a)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(1 AS a) zzz
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CAST(1 AS a)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CAST(1 AS a) zzz
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT PARSE(1 AS a)
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT PARSE(1 AS a) zzz
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_PARSE(1 AS a)
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function TRY_PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_PARSE(1 AS a) zzz
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function TRY_PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT IIF(1 = 1, 1, 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT IIF(1 = 1, 1, 1) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT IIF(1 = 1, 1, 1) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT()
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT() zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT JSON_OBJECT() zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY()
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY() zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT JSON_ARRAY() zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(1)
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 1 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(1) zzz
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 1 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $PARTITION.a(1)
|
< error
| 208 80 16 1 proc=_v Invalid object name '$partition.a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $PARTITION.a(1) zzz
|
< error
| 208 80 16 1 proc=_v Invalid object name '$partition.a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a zzz
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NULLIF(1, 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NULLIF(1, 1) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT NULLIF(1, 1) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COALESCE(1, 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COALESCE(1, 1) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT COALESCE(1, 1) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE 1 WHEN 1 THEN 1 END
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE 1 WHEN 1 THEN 1 END zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CASE 1 WHEN 1 THEN 1 END zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONVERT(a, 1)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONVERT(a, 1) zzz
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CONVERT(a, 1)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CONVERT(a, 1) zzz
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT USER
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT USER zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT USER zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT USER_DB
|
< error
| 207 1 16 1 proc=_v Invalid column name 'USER_DB'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT(1) zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT(1) zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a()
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a() zzz
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: zzz a
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a COLLATE a
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a :: a a
|
< error
| 243 4 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONVERT(a, 1, 1)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CONVERT(a, 1, 1)
|
< error
| 243 1 16 1 proc=_v Type a is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CURRENT_USER
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CURRENT_USER zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CURRENT_USER zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SESSION_USER
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SESSION_USER zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT SESSION_USER zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SYSTEM_USER
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SYSTEM_USER zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT SYSTEM_USER zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SYSTEM_VERSIONING
|
< error
| 207 1 16 1 proc=_v Invalid column name 'SYSTEM_VERSIONING'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP zzz
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CURRENT_TIMESTAMP zzz
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY(1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY(1, 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT(a: 1)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT(a: 1, a: 1)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'a'.
| 207 1 16 1 proc=_v Invalid column name 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT(NULL ON NULL)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT(ABSENT ON NULL)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_OBJECT(1 : 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WINDOW a AS(zzz)
|
< error
| 5362 7 15 1 proc=_v Window 'zzz' is undefined.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WINDOW a AS(a)
|
< error
| 5362 7 15 1 proc=_v Window 'a' is undefined.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WINDOW a AS(PARTITION BY 1)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * WINDOW a AS(ORDER BY 1)
|
< error
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY(zzz)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'zzz'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_ARRAY(NULL ON NULL)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(1 FROM 1)
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 1 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COALESCE(1, 1, 1)
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN 1 = 1 THEN 1 END
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE 1 WHEN 1 THEN 1 ELSE 1 END
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE 1 WHEN 1 THEN 1 WHEN 1 THEN 1 END
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN 1 = 1 THEN 1 WHEN 1 = 1 THEN 1 END
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $IDENTITY
|
< error
| 207 1 16 1 proc=_v Invalid column name '$IDENTITY'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $IDENTITY zzz
|
< error
| 207 1 16 1 proc=_v Invalid column name '$IDENTITY'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT IDENTITYCOL
|
< error
| 207 1 16 1 proc=_v Invalid column name 'identitycol'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT IDENTITYCOL zzz
|
< error
| 207 1 16 1 proc=_v Invalid column name 'identitycol'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ROWGUIDCOL
|
< error
| 207 1 16 1 proc=_v Invalid column name 'rowguidcol'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ROWGUIDCOL zzz
|
< error
| 207 1 16 1 proc=_v Invalid column name 'rowguidcol'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ROWS
|
< error
| 207 1 16 1 proc=_v Invalid column name 'ROWS'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.$PARTITION.a(1)
|
< error
| 208 80 16 1 proc=_v Invalid object name 'a.$partition.a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT.$PARTITION.a(1)
|
< error
| 208 80 16 1 proc=_v Invalid object name '$partition.a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.IDENTITYCOL
|
< error
| 107 1 15 1 proc=_v The column prefix 'a' does not match with a table name or alias name used in the query.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.IDENTITYCOL zzz
|
< error
| 107 1 15 1 proc=_v The column prefix 'a' does not match with a table name or alias name used in the query.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.*
|
< error
| 107 1 15 1 proc=_v The column prefix 'a' does not match with a table name or alias name used in the query.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * ORDER BY 1 DESC
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 COLLATE a
|
< error
| 447 0 16 1 proc=_v Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 COLLATE a zzz
|
< error
| 447 0 16 1 proc=_v Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME AS OF @a
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a FOR SYSTEM_TIME AS OF @a zzz
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@a".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM OPENJSON(1) a
|
|
< row
: string
| TABLE
: string
| dbo.a
: string
| CREATE TABLE [dbo].[a] ([b] int NOT NULL, [a] int NULL, [d] int NULL);
| ALTER TABLE [dbo].[a] ADD CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b] ASC);
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d] ASC);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * GROUP BY GROUPING SETS(1, 1)
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
| 263 1 16 1 proc=_v Must specify table to select from.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT PARSE(1 AS a USING 1)
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_PARSE(1 AS a USING 1)
|
< error
| 10761 2 15 1 proc=_v Invalid data type a in function TRY_PARSE.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(LEADING 1 FROM 1)
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(LEADING 1 FROM 1) zzz
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(TRAILING 1 FROM 1)
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(TRAILING 1 FROM 1) zzz
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH 1 FROM 1)
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(BOTH 1 FROM 1) zzz
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 2 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a OVER a
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a a
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE NONCLUSTERED INDEX [a_d_idx] ON [dbo].[a] ([d]);
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR a OVER a zzz
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.* FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT a.* FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT dbo.a.* FROM dbo.a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT dbo.a.* FROM dbo.a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b, a, d FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT b, a, d FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b AS x, a y, z = d FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT b AS x, a y, z = d FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT [b] AS [x y], 'lit' AS [lit] FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT [b] AS [x y], 'lit' AS [lit] FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b, b FROM a
|
< error
| 4506 1 16 1 proc=_v Column names in each view or function must be unique. Column name 'b' in view or function '_v' is specified more than once.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b AS x, a AS x FROM a
|
< error
| 4506 1 16 1 proc=_v Column names in each view or function must be unique. Column name 'x' in view or function '_v' is specified more than once.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b + 1 FROM a
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 AS one
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT 1 AS one
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 AS one, 2
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 2.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT a.b FROM a AS t
|
< error
| 4104 1 16 1 proc=_v The multi-part identifier "a.b" could not be bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT t.b FROM a AS t
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT t.b FROM a AS t
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT nope FROM a
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.b FROM a
|
< error
| 4104 1 16 1 proc=_v The multi-part identifier "x.b" could not be bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM nope
|
< error
| 208 1 16 1 proc=_v Invalid object name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM dbo.nope
|
< error
| 208 1 16 1 proc=_v Invalid object name 'dbo.nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM s.orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM s.orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, note FROM s.orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, note FROM s.orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM v_customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM v_customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, name FROM dbo.v_customers WHERE id > 1
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, name FROM dbo.v_customers WHERE id > 1
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DISTINCT region FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT DISTINCT region FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ALL region FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT ALL region FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $IDENTITY FROM customers
|
< error
| 207 1 16 1 proc=_v Invalid column name '$IDENTITY'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT $ROWGUID FROM customers
|
< error
| 207 1 16 1 proc=_v Invalid column name '$ROWGUID'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT @x AS v
|
< error
| 137 2 15 1 proc=_v Must declare the scalar variable "@x".
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b INTO t2 FROM a
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'INTO'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a ORDER BY b
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM #tmp
|
< error
| 4508 1 16 1 proc=_v Views or functions are not allowed on temporary tables. Table names that begin with '#' denote temporary tables.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a WITH (NOLOCK)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT b FROM a WITH (NOLOCK)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a WITH (NOEXPAND)
|
< error
| 8171 2 16 1 proc=_v Hint 'noexpand' on object 'a' is invalid.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a OPTION (MAXDOP 1)
|
< error
| 156 1 15 1 proc=_v Incorrect syntax near the keyword 'OPTION'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE Latin1_General_BIN2 AS n FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name COLLATE Latin1_General_BIN2 AS n FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customers.name FROM dbo.customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customers.name FROM dbo.customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT dbo.customers.name FROM dbo.customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT dbo.customers.name FROM dbo.customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id FROM customers c WHERE c.region = 'EU'
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id FROM customers c WHERE c.region = 'EU'
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o INNER JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o INNER JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o LEFT OUTER JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o LEFT OUTER JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o RIGHT JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o RIGHT JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o CROSS JOIN customers c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o CROSS JOIN customers c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id, c.id FROM orders o JOIN customers c ON c.id = o.customer_id
|
< error
| 4506 1 16 1 proc=_v Column names in each view or function must be unique. Column name 'id' in view or function '_v' is specified more than once.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders o JOIN customers c ON c.id = o.customer_id
|
< error
| 209 1 16 1 proc=_v Ambiguous column name 'id'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id
|
< error
| 4506 1 16 1 proc=_v Column names in each view or function must be unique. Column name 'id' in view or function '_v' is specified more than once.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN customers c ON c.nope = o.customer_id
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN customers o ON o.id = o.id
|
< error
| 1011 1 16 1 proc=_v The correlation name 'o' is specified multiple times in a FROM clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT orders.id FROM orders JOIN s.orders ON s.orders.id = dbo.orders.id
|
< error
| 1013 1 16 1 proc=_v The objects "s.orders" and "orders" in the FROM clause have the same exposed names. Use correlation names to distinguish them.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o INNER HASH JOIN customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id FROM orders o INNER HASH JOIN customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN customers c JOIN a ON a.b = c.id ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN customers c JOIN a ON a.b = c.id ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM (orders o JOIN customers c ON c.id = o.customer_id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id FROM (orders o JOIN customers c ON c.id = o.customer_id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN customers c ON o.id = x.id
|
< error
| 4104 1 16 1 proc=_v The multi-part identifier "x.id" could not be bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount BETWEEN 1 AND 10
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount BETWEEN 1 AND 10
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE status IN ('new', 'paid')
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE status IN ('new', 'paid')
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE status NOT IN ('new')
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE status NOT IN ('new')
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE status LIKE 'p%' ESCAPE '\'
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE status LIKE 'p%' ESCAPE '\'
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE status IS NULL OR amount IS NOT NULL
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE status IS NULL OR amount IS NOT NULL
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE NOT (amount > 1 AND amount < 5)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE NOT (amount > 1 AND amount < 5)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE EXISTS (SELECT 1 FROM customers c WHERE c.id = orders.customer_id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE EXISTS (SELECT 1 FROM customers c WHERE c.id = orders.customer_id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > ALL (SELECT amount FROM orders)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > ALL (SELECT amount FROM orders)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > ANY (SELECT amount FROM orders)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > ANY (SELECT amount FROM orders)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount = SOME (SELECT amount FROM orders)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount = SOME (SELECT amount FROM orders)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE customer_id IN (SELECT nope FROM customers)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE customer_id IN (SELECT id, name FROM customers)
|
< error
| 116 1 16 1 proc=_v Only one expression can be specified in the select list when the subquery is not introduced with EXISTS.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > (SELECT MAX(amount) FROM orders)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > (SELECT MAX(amount) FROM orders)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount > (SELECT amount, id FROM orders)
|
< error
| 116 1 16 1 proc=_v Only one expression can be specified in the select list when the subquery is not introduced with EXISTS.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE IIF(amount > 1, 1, 0) = 1
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE IIF(amount > 1, 1, 0) = 1
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount IS DISTINCT FROM 1
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders WHERE amount IS DISTINCT FROM 1
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE CONTAINS(status, 'x')
|
< error
| 7601 2 16 1 proc=_v Cannot use a CONTAINS or FREETEXT predicate on table or indexed view 'orders' because it is not full-text indexed.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE ROW_NUMBER() OVER (ORDER BY id) = 1
|
< error
| 4108 1 15 1 proc=_v Windowed functions can only appear in the SELECT or ORDER BY clauses.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders WHERE SUM(amount) > 1
|
< error
| 147 1 15 1 proc=_v An aggregate may not appear in the WHERE clause unless it is in a subquery contained in a HAVING clause or a select list, and the column being aggregated is an outer reference.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN amount > 100 THEN 'big' ELSE 'small' END AS size FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CASE WHEN amount > 100 THEN 'big' ELSE 'small' END AS size FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE status WHEN 'new' THEN 1 WHEN 'paid' THEN 2 END AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CASE status WHEN 'new' THEN 1 WHEN 'paid' THEN 2 END AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN amount > 100 THEN 'big' END FROM orders
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN nope > 1 THEN 1 END AS c FROM orders
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT IIF(amount > 1, 'y', 'n') AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT IIF(amount > 1, 'y', 'n') AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CHOOSE(2, 'a', 'b', 'c') AS c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CHOOSE(2, 'a', 'b', 'c') AS c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COALESCE(status, 'none') AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT COALESCE(status, 'none') AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ISNULL(status, 'none') AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT ISNULL(status, 'none') AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NULLIF(status, '') AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT NULLIF(status, '') AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COALESCE(status) AS c FROM orders
|
< error
| 102 1 15 1 proc=_v Incorrect syntax near ')'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(amount AS INT) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CAST(amount AS INT) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(amount AS VARCHAR(20)) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CAST(amount AS VARCHAR(20)) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(amount AS NOPE) AS c FROM orders
|
< error
| 243 1 16 1 proc=_v Type NOPE is not a defined system type.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONVERT(VARCHAR(10), created, 120) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CONVERT(VARCHAR(10), created, 120) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONVERT(INT, amount) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CONVERT(INT, amount) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CAST(status AS INT) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TRY_CAST(status AS INT) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_CONVERT(INT, status) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TRY_CONVERT(INT, status) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT PARSE('1' AS INT) AS c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT PARSE('1' AS INT) AS c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRY_PARSE('1' AS INT USING 'en-US') AS c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TRY_PARSE('1' AS INT USING 'en-US') AS c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(doc AS NVARCHAR(MAX)) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CAST(doc AS NVARCHAR(MAX)) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT created AT TIME ZONE 'UTC' AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT created AT TIME ZONE 'UTC' AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CAST(amount AS INT) FROM orders
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id COLLATE Latin1_General_BIN2 AS c FROM customers
|
< error
| 447 0 16 1 proc=_v Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT LEN(name) AS l, UPPER(name) AS u, LOWER(name) AS lo FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT LEN(name) AS l, UPPER(name) AS u, LOWER(name) AS lo FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(name, 1, 3) AS s, LEFT(name, 2) AS l, RIGHT(name, 2) AS r FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT SUBSTRING(name, 1, 3) AS s, LEFT(name, 2) AS l, RIGHT(name, 2) AS r FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT LTRIM(RTRIM(name)) AS t, TRIM(name) AS t2 FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT LTRIM(RTRIM(name)) AS t, TRIM(name) AS t2 FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM('x' FROM name) AS t FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TRIM('x' FROM name) AS t FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT REPLACE(name, 'a', 'b') AS r, CHARINDEX('a', name) AS c, PATINDEX('%a%', name) AS p FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT REPLACE(name, 'a', 'b') AS r, CHARINDEX('a', name) AS c, PATINDEX('%a%', name) AS p FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONCAT(name, '-', region) AS c, CONCAT_WS(',', name, region) AS w FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CONCAT(name, '-', region) AS c, CONCAT_WS(',', name, region) AS w FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT STUFF(name, 1, 1, 'x') AS s, REVERSE(name) AS r, REPLICATE('x', 3) AS p FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT STUFF(name, 1, 1, 'x') AS s, REVERSE(name) AS r, REPLICATE('x', 3) AS p FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT FORMAT(amount, 'N2') AS f, STR(amount, 10, 2) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT FORMAT(amount, 'N2') AS f, STR(amount, 10, 2) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT GETDATE() AS d, SYSDATETIME() AS s, GETUTCDATE() AS u, CURRENT_TIMESTAMP AS c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT GETDATE() AS d, SYSDATETIME() AS s, GETUTCDATE() AS u, CURRENT_TIMESTAMP AS c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DATEADD(day, 1, created) AS d, DATEDIFF(day, created, GETDATE()) AS n FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT DATEADD(day, 1, created) AS d, DATEDIFF(day, created, GETDATE()) AS n FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DATEPART(year, created) AS y, DATENAME(month, created) AS m, YEAR(created) AS y2, MONTH(created) AS m2, DAY(created) AS d2 FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT DATEPART(year, created) AS y, DATENAME(month, created) AS m, YEAR(created) AS y2, MONTH(created) AS m2, DAY(created) AS d2 FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT EOMONTH(created) AS e, DATEFROMPARTS(2020, 1, 1) AS d FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT EOMONTH(created) AS e, DATEFROMPARTS(2020, 1, 1) AS d FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DATEADD(nope, 1, created) AS d FROM orders
|
< error
| 155 1 15 1 proc=_v 'nope' is not a recognized dateadd option.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ABS(amount) AS a, ROUND(amount, 1) AS r, CEILING(amount) AS c, FLOOR(amount) AS f, POWER(amount, 2) AS p, SQRT(amount) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT ABS(amount) AS a, ROUND(amount, 1) AS r, CEILING(amount) AS c, FLOOR(amount) AS f, POWER(amount, 2) AS p, SQRT(amount) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEWID() AS g
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT NEWID() AS g
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT RAND() AS r
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT RAND() AS r
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT @@ROWCOUNT AS r
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT @@ROWCOUNT AS r
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT @@SPID AS s, @@VERSION AS v
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT @@SPID AS s, @@VERSION AS v
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DB_NAME() AS d, SCHEMA_NAME() AS s, OBJECT_ID('dbo.a') AS o, USER_NAME() AS u, SUSER_SNAME() AS n
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT DB_NAME() AS d, SCHEMA_NAME() AS s, OBJECT_ID('dbo.a') AS o, USER_NAME() AS u, SUSER_SNAME() AS n
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ISJSON(status) AS j, JSON_VALUE(status, '$.a') AS v, JSON_QUERY(status, '$.a') AS q FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT ISJSON(status) AS j, JSON_VALUE(status, '$.a') AS v, JSON_QUERY(status, '$.a') AS q FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT JSON_MODIFY(status, '$.a', 1) AS m FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT JSON_MODIFY(status, '$.a', 1) AS m FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT doc.value('(/a)[1]', 'INT') AS v FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT doc.value('(/a)[1]', 'INT') AS v FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT doc.query('/a') AS q, doc.exist('/a') AS e FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT doc.query('/a') AS q, doc.exist('/a') AS e FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT STRING_SPLIT('a,b', ',') AS s
|
< error
| 195 10 15 1 proc=_v 'STRING_SPLIT' is not a recognized built-in function name.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT LEN(name, 1) AS l FROM customers
|
< error
| 174 1 15 1 proc=_v The len function requires 1 argument(s).
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT LEN() AS l FROM customers
|
< error
| 174 1 15 1 proc=_v The len function requires 1 argument(s).
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NOPE(name) AS l FROM customers
|
< error
| 195 10 15 1 proc=_v 'NOPE' is not a recognized built-in function name.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT dbo.f_double(id) AS d FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT dbo.f_double(id) AS d FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT f_double(id) AS d FROM customers
|
< error
| 195 10 15 1 proc=_v 'f_double' is not a recognized built-in function name.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT dbo.f_double(id, 1) AS d FROM customers
|
< error
| 8144 2 16 1 proc=_v Procedure or function dbo.f_double has too many arguments specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT dbo.nope(id) AS d FROM customers
|
< error
| 4121 1 16 1 proc=_v Cannot find either column "dbo" or the user-defined function or aggregate "dbo.nope", or the name is ambiguous.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT GREATEST(1, 2) AS g, LEAST(1, 2) AS l
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT GREATEST(1, 2) AS g, LEAST(1, 2) AS l
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DATETRUNC(month, created) AS d FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT DATETRUNC(month, created) AS d FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT STRING_ESCAPE(name, 'json') AS e FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT STRING_ESCAPE(name, 'json') AS e FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRANSLATE(name, 'ab', 'cd') AS t FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TRANSLATE(name, 'ab', 'cd') AS t FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SOUNDEX(name) AS s, DIFFERENCE(name, 'x') AS d FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT SOUNDEX(name) AS s, DIFFERENCE(name, 'x') AS d FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT HASHBYTES('SHA2_256', name) AS h FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT HASHBYTES('SHA2_256', name) AS h FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CHECKSUM(*) AS c FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CHECKSUM(*) AS c FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT BINARY_CHECKSUM(name) AS c FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT BINARY_CHECKSUM(name) AS c FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT NEXT VALUE FOR nope AS n
|
< error
| 11719 1 15 1 proc=_v NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CONCAT(name) AS c FROM customers
|
< error
| 189 1 15 1 proc=_v The concat function requires 2 to 254 arguments.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUBSTRING(name, 1) AS s FROM customers
|
< error
| 174 1 15 1 proc=_v The substring function requires 3 argument(s).
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT ROUND(amount) AS r FROM orders
|
< error
| 189 1 15 1 proc=_v The round function requires 2 to 3 arguments.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DATEDIFF(nope, created, created) AS d FROM orders
|
< error
| 155 1 15 1 proc=_v 'nope' is not a recognized datediff option.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TRIM(id) AS t FROM customers
|
< error
| 8116 1 16 1 proc=_v Argument data type int is invalid for argument 1 of Trim function.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT UPPER(id) AS u FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT UPPER(id) AS u FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n, SUM(amount) AS s FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n, SUM(amount) AS s FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, COUNT(DISTINCT status) AS n FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, COUNT(DISTINCT status) AS n FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, AVG(amount) AS a, MIN(amount) AS mi, MAX(amount) AS ma FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, AVG(amount) AS a, MIN(amount) AS mi, MAX(amount) AS ma FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, COUNT_BIG(*) AS n FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, COUNT_BIG(*) AS n FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, STDEV(amount) AS s, VAR(amount) AS v FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, STDEV(amount) AS s, VAR(amount) AS v FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, STRING_AGG(status, ',') AS s FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, STRING_AGG(status, ',') AS s FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, STRING_AGG(status, ',') WITHIN GROUP (ORDER BY status) AS s FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, STRING_AGG(status, ',') WITHIN GROUP (ORDER BY status) AS s FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING SUM(amount) > 10
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING SUM(amount) > 10
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY customer_id
|
< error
| 8120 1 16 1 proc=_v Column 'orders.status' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT status FROM orders GROUP BY customer_id
|
< error
| 8120 1 16 1 proc=_v Column 'orders.status' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id FROM orders GROUP BY customer_id HAVING status = 'new'
|
< error
| 8121 1 16 1 proc=_v Column 'orders.status' is invalid in the HAVING clause because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COUNT(*) AS n FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT COUNT(*) AS n FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT COUNT(*) FROM orders
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY CUBE (customer_id, status)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY CUBE (customer_id, status)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY GROUPING SETS ((customer_id), (status), ())
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, status, COUNT(*) AS n FROM orders GROUP BY GROUPING SETS ((customer_id), (status), ())
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, GROUPING(customer_id) AS g, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, GROUPING(customer_id) AS g, COUNT(*) AS n FROM orders GROUP BY ROLLUP (customer_id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY customer_id WITH ROLLUP
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, COUNT(*) AS n FROM orders GROUP BY customer_id WITH ROLLUP
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUM(SUM(amount)) AS s FROM orders GROUP BY customer_id
|
< error
| 130 1 15 1 proc=_v Cannot perform an aggregate function on an expression containing an aggregate or a subquery.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUM(name) AS s FROM customers
|
< error
| 8117 1 16 1 proc=_v Operand data type nvarchar is invalid for sum operator.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id FROM orders GROUP BY nope
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUM(amount) AS s FROM orders GROUP BY 1
|
< error
| 164 1 15 1 proc=_v Each GROUP BY expression must contain at least one column that is not an outer reference.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, APPROX_COUNT_DISTINCT(status) AS n FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, APPROX_COUNT_DISTINCT(status) AS n FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT CASE WHEN amount > 10 THEN 1 ELSE 0 END AS big, COUNT(*) AS n FROM orders GROUP BY CASE WHEN amount > 10 THEN 1 ELSE 0 END
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT CASE WHEN amount > 10 THEN 1 ELSE 0 END AS big, COUNT(*) AS n FROM orders GROUP BY CASE WHEN amount > 10 THEN 1 ELSE 0 END
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT YEAR(created) AS y, COUNT(*) AS n FROM orders GROUP BY YEAR(created)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT YEAR(created) AS y, COUNT(*) AS n FROM orders GROUP BY YEAR(created)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT YEAR(created) AS y, COUNT(*) AS n FROM orders GROUP BY MONTH(created)
|
< error
| 8120 1 16 1 proc=_v Column 'orders.created' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.customer_id, c.name, SUM(o.amount) AS s FROM orders o JOIN customers c ON c.id = o.customer_id GROUP BY o.customer_id, c.name
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.customer_id, c.name, SUM(o.amount) AS s FROM orders o JOIN customers c ON c.id = o.customer_id GROUP BY o.customer_id, c.name
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING COUNT(*) > 1 AND MAX(amount) > 5
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id HAVING COUNT(*) > 1 AND MAX(amount) > 5
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 3 customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id ORDER BY s DESC
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 3 customer_id, SUM(amount) AS s FROM orders GROUP BY customer_id ORDER BY s DESC
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 3 customer_id FROM orders GROUP BY customer_id ORDER BY SUM(amount) DESC
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 3 customer_id FROM orders GROUP BY customer_id ORDER BY SUM(amount) DESC
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 3 customer_id FROM orders GROUP BY customer_id ORDER BY amount
|
< error
| 8127 1 16 1 proc=_v Column "orders.amount" is invalid in the ORDER BY clause because it is not contained in either an aggregate function or the GROUP BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id + 1 AS c FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id + 1 AS c FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT (SELECT COUNT(*) FROM customers c WHERE c.id = o.customer_id) AS n FROM orders o GROUP BY o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT (SELECT COUNT(*) FROM customers c WHERE c.id = o.customer_id) AS n FROM orders o GROUP BY o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT AVG(name) AS a FROM customers
|
< error
| 8117 1 16 1 proc=_v Operand data type nvarchar is invalid for avg operator.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT SUM(born) AS s FROM customers
|
< error
| 8117 1 16 1 proc=_v Operand data type date is invalid for sum operator.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT MAX(name) AS m, MIN(born) AS b FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT MAX(name) AS m, MIN(born) AS b FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS rn FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS rn FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY created DESC) AS rn FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY created DESC) AS rn FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, RANK() OVER (ORDER BY amount) AS r, DENSE_RANK() OVER (ORDER BY amount) AS d, NTILE(4) OVER (ORDER BY amount) AS n FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, RANK() OVER (ORDER BY amount) AS r, DENSE_RANK() OVER (ORDER BY amount) AS d, NTILE(4) OVER (ORDER BY amount) AS n FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, AVG(amount) OVER (ORDER BY id ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, AVG(amount) OVER (ORDER BY id ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, SUM(amount) OVER (ORDER BY id RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS s FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, LAG(amount) OVER (ORDER BY id) AS p, LEAD(amount, 1, 0) OVER (ORDER BY id) AS n FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, LAG(amount) OVER (ORDER BY id) AS p, LEAD(amount, 1, 0) OVER (ORDER BY id) AS n FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, FIRST_VALUE(amount) OVER (ORDER BY id) AS f, LAST_VALUE(amount) OVER (ORDER BY id) AS l FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, FIRST_VALUE(amount) OVER (ORDER BY id) AS f, LAST_VALUE(amount) OVER (ORDER BY id) AS l FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, PERCENT_RANK() OVER (ORDER BY amount) AS p, CUME_DIST() OVER (ORDER BY amount) AS c FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, PERCENT_RANK() OVER (ORDER BY amount) AS p, CUME_DIST() OVER (ORDER BY amount) AS c FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) OVER (PARTITION BY customer_id) AS m FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) OVER (PARTITION BY customer_id) AS m FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER () AS rn FROM orders
|
< error
| 4112 1 15 1 proc=_v The function 'ROW_NUMBER' must have an OVER clause with ORDER BY.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, LAG(amount) OVER (PARTITION BY customer_id) AS p FROM orders
|
< error
| 4112 1 15 1 proc=_v The function 'LAG' must have an OVER clause with ORDER BY.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER w AS rn FROM orders WINDOW w AS (ORDER BY id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER w AS rn FROM orders WINDOW w AS (ORDER BY id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER (ORDER BY nope) AS rn FROM orders
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, COUNT(*) OVER () AS n FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, COUNT(*) OVER () AS n FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT customer_id, SUM(SUM(amount)) OVER () AS total FROM orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT customer_id, SUM(SUM(amount)) OVER () AS total FROM orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER w AS rn, RANK() OVER w AS r FROM orders WINDOW w AS (PARTITION BY customer_id ORDER BY id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER w AS rn, RANK() OVER w AS r FROM orders WINDOW w AS (PARTITION BY customer_id ORDER BY id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER w AS rn FROM orders WINDOW w AS (PARTITION BY customer_id)
|
< error
| 5366 3 15 1 proc=_v The function 'row_number' must have an OVER clause or a WINDOW with ORDER BY.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, ROW_NUMBER() OVER nope AS rn FROM orders
|
< error
| 5362 5 15 1 proc=_v Window 'nope' is undefined.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 10 id FROM orders ORDER BY id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 10 id FROM orders ORDER BY id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP (10) id FROM orders ORDER BY created DESC
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP (10) id FROM orders ORDER BY created DESC
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 10 PERCENT id FROM orders ORDER BY id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 10 PERCENT id FROM orders ORDER BY id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 10 WITH TIES id FROM orders ORDER BY amount
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 10 WITH TIES id FROM orders ORDER BY amount
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 10 WITH TIES id FROM orders
|
< error
| 1062 1 15 1 proc=_v The TOP N WITH TIES clause is not allowed without a corresponding ORDER BY clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 100 PERCENT id FROM orders ORDER BY id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 100 PERCENT id FROM orders ORDER BY id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 5 id FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 5 id FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 0 ROWS
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 0 ROWS
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH FIRST 5 ROW ONLY
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY id OFFSET 10 ROWS FETCH FIRST 5 ROW ONLY
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 5 id FROM orders ORDER BY id OFFSET 1 ROWS
|
< error
| 10741 1 15 1 proc=_v A TOP can not be used in the same query or sub-query as a OFFSET.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders ORDER BY nope
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 id FROM orders ORDER BY 1
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 1 id FROM orders ORDER BY 1
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 id AS x FROM orders ORDER BY x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 1 id AS x FROM orders ORDER BY x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 id FROM orders ORDER BY 2
|
< error
| 108 1 16 1 proc=_v The ORDER BY position number 2 is out of range of the number of items in the select list.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT DISTINCT TOP 1 status FROM orders ORDER BY id
|
< error
| 145 1 15 1 proc=_v ORDER BY items must appear in the select list if SELECT DISTINCT is specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id FROM orders) AS x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id FROM orders) AS x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.n FROM (SELECT id AS n FROM orders) x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.n FROM (SELECT id AS n FROM orders) x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.n FROM (SELECT id FROM orders) x (n)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.n FROM (SELECT id FROM orders) x (n)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id FROM orders) x (n)
|
< error
| 207 1 16 1 proc=_v Invalid column name 'id'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id FROM orders)
|
< error
| 102 1 15 1 proc=_v Incorrect syntax near ')'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.c FROM (SELECT COUNT(*) FROM orders) x
|
< error
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'x'.
| 207 1 16 1 proc=_v Invalid column name 'c'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.c FROM (SELECT COUNT(*) FROM orders) x (c)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.c FROM (SELECT COUNT(*) FROM orders) x (c)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id, id FROM orders) x
|
< error
| 8156 1 16 1 proc=_v The column 'id' was specified multiple times for 'x'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.a FROM (VALUES (1, 2), (3, 4)) AS x (a, b)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.a FROM (VALUES (1, 2), (3, 4)) AS x (a, b)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.a FROM (VALUES (1, 2), (3, 4)) AS x (a)
|
< error
| 8158 1 16 1 proc=_v 'x' has more columns than were specified in the column list.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.a FROM (VALUES (1, 2), (3)) AS x (a, b)
|
< error
| 10709 1 16 1 proc=_v The number of columns for each row in a table value constructor must be the same.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT (SELECT MAX(amount) FROM orders) AS m
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT (SELECT MAX(amount) FROM orders) AS m
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, (SELECT COUNT(*) FROM orders o WHERE o.customer_id = c.id) AS n FROM customers c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, (SELECT COUNT(*) FROM orders o WHERE o.customer_id = c.id) AS n FROM customers c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT id FROM orders ORDER BY id) x
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT x.id FROM (SELECT TOP 1 id FROM orders ORDER BY id) x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT x.id FROM (SELECT TOP 1 id FROM orders ORDER BY id) x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT f.id, f.amount FROM dbo.f_orders(1) AS f
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT f.id, f.amount FROM dbo.f_orders(1) AS f
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT f.id FROM dbo.f_orders(1, 2) AS f
|
< error
| 8144 3 16 1 proc=_v Procedure or function dbo.f_orders has too many arguments specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM dbo.f_double(1)
|
< error
| 208 3 16 1 proc=_v Invalid object name 'dbo.f_double'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT f.nope FROM dbo.f_orders(1) AS f
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT value FROM STRING_SPLIT('a,b', ',')
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT value FROM STRING_SPLIT('a,b', ',')
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT [key], value FROM OPENJSON('{}')
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT [key], value FROM OPENJSON('{}')
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT j.a FROM OPENJSON('{}') WITH (a INT '$.a') AS j
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT j.a FROM OPENJSON('{}') WITH (a INT '$.a') AS j
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT value FROM GENERATE_SERIES(1, 10)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT value FROM GENERATE_SERIES(1, 10)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT d.c FROM (SELECT customer_id AS c FROM orders GROUP BY customer_id) d JOIN customers ON customers.id = d.c
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT d.c FROM (SELECT customer_id AS c FROM orders GROUP BY customer_id) d JOIN customers ON customers.id = d.c
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o WHERE o.amount > (SELECT AVG(i.amount) FROM orders i WHERE i.customer_id = o.customer_id)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT o.id FROM orders o WHERE o.amount > (SELECT AVG(i.amount) FROM orders i WHERE i.customer_id = o.customer_id)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT o.id FROM orders o JOIN (SELECT id FROM customers WHERE id = o.customer_id) c ON c.id = o.id
|
< error
| 4104 1 16 1 proc=_v The multi-part identifier "o.customer_id" could not be bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x (n) AS (SELECT id FROM orders) SELECT n FROM x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x (n) AS (SELECT id FROM orders) SELECT n FROM x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x (n) AS (SELECT id FROM orders) SELECT id FROM x
|
< error
| 207 1 16 1 proc=_v Invalid column name 'id'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x (n, m) AS (SELECT id FROM orders) SELECT n FROM x
|
< error
| 8159 1 16 1 proc=_v 'x' has fewer columns than were specified in the column list.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders), y AS (SELECT id FROM x) SELECT y.id FROM y JOIN x ON x.id = y.id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders), y AS (SELECT id FROM x) SELECT y.id FROM y JOIN x ON x.id = y.id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT COUNT(*) FROM orders) SELECT * FROM x
|
< error
| 8155 2 16 1 proc=_v No column name was specified for column 1 of 'x'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders), x AS (SELECT id FROM orders) SELECT id FROM x
|
< error
| 239 1 16 1 proc=_v Duplicate common table expression name 'x' was specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 10) SELECT n FROM r
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 10) SELECT n FROM r
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH r (n) AS (SELECT n + 1 FROM r) SELECT n FROM r
|
< error
| 252 1 16 1 proc=_v Recursive common table expression 'r' does not contain a top-level UNION ALL operator.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM nope) SELECT id FROM x
|
< error
| 208 1 16 1 proc=_v Invalid object name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM orders
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM orders
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders ORDER BY id) SELECT id FROM x
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT a.id FROM x a JOIN x b ON a.id = b.id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT a.id FROM x a JOIN x b ON a.id = b.id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH r (n, p) AS (SELECT id, CAST(NULL AS INT) FROM customers UNION ALL SELECT r.n, r.p FROM r JOIN customers c ON c.id = r.n WHERE 1 = 0) SELECT n, p FROM r
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH r (n, p) AS (SELECT id, CAST(NULL AS INT) FROM customers UNION ALL SELECT r.n, r.p FROM r JOIN customers c ON c.id = r.n WHERE 1 = 0) SELECT n, p FROM r
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM x WHERE id IN (SELECT id FROM x)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS WITH x AS (SELECT id FROM orders) SELECT id FROM x WHERE id IN (SELECT id FROM x)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS WITH x AS (SELECT 1 AS n), y AS (SELECT n FROM z), z AS (SELECT 1 AS n) SELECT n FROM y
|
< error
| 208 1 16 1 proc=_v Invalid object name 'z'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders UNION ALL SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders UNION ALL SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders INTERSECT SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders INTERSECT SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders EXCEPT SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders EXCEPT SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, status FROM orders UNION SELECT id FROM customers
|
< error
| 205 1 16 1 proc=_v All queries combined using a UNION, INTERSECT or EXCEPT operator must have an equal number of expressions in their target lists.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id AS x FROM orders UNION SELECT id AS y FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id AS x FROM orders UNION SELECT id AS y FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 1 UNION SELECT 2
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS (SELECT id FROM orders) UNION (SELECT id FROM customers)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS (SELECT id FROM orders) UNION (SELECT id FROM customers)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT id FROM customers ORDER BY id
|
< error
| 1033 1 15 1 proc=_v The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 1 id FROM orders UNION SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 1 id FROM orders UNION SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT nope FROM customers
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT id FROM customers UNION ALL SELECT b FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders UNION SELECT id FROM customers UNION ALL SELECT b FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TOP 2 id FROM orders UNION ALL SELECT TOP 2 id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TOP 2 id FROM orders UNION ALL SELECT TOP 2 id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM orders EXCEPT SELECT b FROM a INTERSECT SELECT id FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT id FROM orders EXCEPT SELECT b FROM a INTERSECT SELECT id FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.amount FROM customers c OUTER APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, x.amount FROM customers c OUTER APPLY (SELECT TOP 1 amount FROM orders o WHERE o.customer_id = c.id ORDER BY created DESC) x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, f.amount FROM customers c CROSS APPLY dbo.f_orders(c.id) f
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, f.amount FROM customers c CROSS APPLY dbo.f_orders(c.id) f
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, s.value FROM customers c CROSS APPLY STRING_SPLIT(c.name, ' ') s
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, s.value FROM customers c CROSS APPLY STRING_SPLIT(c.name, ' ') s
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c CROSS APPLY (SELECT c.nope AS n) x
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c CROSS APPLY (VALUES (c.id), (c.id + 1)) x (n)
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c CROSS APPLY (VALUES (c.id), (c.id + 1)) x (n)
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT c.id, x.n FROM customers c JOIN (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x ON 1 = 1
|
< error
| 4104 1 16 1 proc=_v The multi-part identifier "c.id" could not be bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT p.customer_id, p.[new], p.[paid] FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT p.customer_id, p.[new], p.[paid] FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT * FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new], [paid])) AS p
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT p.customer_id FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR nope IN ([new])) AS p
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT p.amount FROM (SELECT customer_id, status, amount FROM orders) AS src PIVOT (SUM(amount) FOR status IN ([new])) AS p
|
< error
| 207 1 16 1 proc=_v Invalid column name 'amount'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT u.id, u.k, u.v FROM (SELECT id, a, d FROM a AS t (id, a, d)) AS src UNPIVOT (v FOR k IN (a, d)) AS u
|
< error
| 321 1 15 1 proc=_v "id" is not a recognized table hints option.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT u.id, u.k, u.v FROM (SELECT b AS id, a, d FROM a) AS src UNPIVOT (v FOR k IN (a, d)) AS u
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT u.id, u.k, u.v FROM (SELECT b AS id, a, d FROM a) AS src UNPIVOT (v FOR k IN (a, d)) AS u
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT u.b, u.k, u.v FROM a UNPIVOT (v FOR k IN (a, d)) AS u
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT u.b, u.k, u.v FROM a UNPIVOT (v FOR k IN (a, d)) AS u
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT u.b, u.k, u.v FROM a UNPIVOT (v FOR k IN (a, nope)) AS u
|
< error
| 207 1 16 1 proc=_v Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, name FROM customers FOR XML PATH('c'), ROOT('cs')
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, name FROM customers FOR XML AUTO
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id, name FROM customers FOR JSON PATH
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT (SELECT id, name FROM customers FOR JSON PATH) AS j
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT (SELECT id, name FROM customers FOR JSON PATH) AS j
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT (SELECT id, name FROM customers FOR XML PATH('c'), TYPE) AS x
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT (SELECT id, name FROM customers FOR XML PATH('c'), TYPE) AS x
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM customers ORDER BY id FOR XML AUTO
|
< error
| 4511 1 16 1 proc=_v Create View or Function failed because no column name was specified for column 1.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name FROM sys.objects
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name FROM sys.objects
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name FROM sys.nope
|
< error
| 208 1 16 1 proc=_v Invalid object name 'sys.nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a TABLESAMPLE (10 PERCENT)
|
< error
| 478 0 16 1 proc=_v The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name FROM master.dbo.spt_values
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name FROM master.dbo.spt_values
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM nope.dbo.a
|
< error
| 208 1 16 1 proc=_v Invalid object name 'nope.dbo.a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM OPENROWSET('x', 'y', 'z')
|
< error
| 7222 255 16 1 proc=_v Only a SQL Server provider is allowed on this instance.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a FOR SYSTEM_TIME AS OF '2020-01-01'
|
< error
| 13544 2 16 1 proc=_v Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. 'omni_TestQuerySimulation.dbo.a' is not a system-versioned table.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a WITH CHECK OPTION
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT b FROM a WITH CHECK OPTION
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x, y) AS SELECT b, a FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v (x, y) AS SELECT b, a FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x) AS SELECT b, a FROM a
|
< error
| 8158 1 16 1 proc=_v '_v' has more columns than were specified in the column list.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x, y, z) AS SELECT b, a FROM a
|
< error
| 8159 1 16 1 proc=_v '_v' has fewer columns than were specified in the column list.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x, x) AS SELECT b, a FROM a
|
< error
| 8156 1 16 1 proc=_v The column 'x' was specified multiple times for '_v'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x, y) AS SELECT b, b FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v (x, y) AS SELECT b, b FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v (x) AS SELECT COUNT(*) FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v (x) AS SELECT COUNT(*) FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT b, a FROM dbo.a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT b, a FROM dbo.a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT b, a FROM a
|
< error
| 4512 3 16 1 proc=_v Cannot schema bind view '_v' because name 'a' is invalid for schema binding. Names must be in two-part format and an object cannot reference itself.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT * FROM dbo.a
|
< error
| 1054 6 15 1 proc=_v Syntax '*' is not allowed in schema-bound objects.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT dbo.f_double(b) AS x FROM dbo.a
|
< error
| 4513 2 16 1 proc=_v Cannot schema bind view '_v'. 'dbo.f_double' is not schema bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT id FROM dbo.v_customers
|
< error
| 4513 2 16 1 proc=_v Cannot schema bind view '_v'. 'dbo.v_customers' is not schema bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH VIEW_METADATA AS SELECT b FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v WITH VIEW_METADATA AS SELECT b FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW dbo._v AS SELECT b FROM a
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW dbo._v AS SELECT b FROM a
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW s._v AS SELECT id FROM orders
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| VIEW
: string
| s._v
: string
| CREATE VIEW s._v AS SELECT id FROM orders
|
|
< row
: string
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW nope._v AS SELECT b FROM a
|
< error
| 2760 1 16 1 proc=_v The specified schema name "nope" either does not exist or you do not have permission to use it.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW a AS SELECT b FROM a
|
< error
| 2714 3 16 1 proc=a There is already an object named 'a' in the database.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW v_customers AS SELECT b FROM a
|
< error
| 2714 3 16 1 proc=v_customers There is already an object named 'v_customers' in the database.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE OR ALTER VIEW v_customers AS SELECT b FROM a
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE   VIEW v_customers AS SELECT b FROM a
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
< row
: string
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER VIEW v_customers AS SELECT b AS id, a AS name FROM a
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| VIEW
: string
| dbo.v_customers
: string
| CREATE VIEW v_customers AS SELECT b AS id, a AS name FROM a
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
< row
: string
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER VIEW v_customers AS SELECT nope FROM a
|
< error
| 207 1 16 1 proc=v_customers Invalid column name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT b FROM a WHERE b IN (SELECT b FROM _v)
|
< error
| 208 1 16 1 proc=_v Invalid object name '_v'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT o.id, c.name FROM dbo.orders o JOIN dbo.customers c ON c.id = o.customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT o.id, c.name FROM dbo.orders o JOIN dbo.customers c ON c.id = o.customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT id FROM dbo.f_orders(1)
|
< error
| 4513 2 16 1 proc=_v Cannot schema bind view '_v'. 'dbo.f_orders' is not schema bound.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT COUNT_BIG(*) AS n, customer_id FROM dbo.orders GROUP BY customer_id
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT COUNT_BIG(*) AS n, customer_id FROM dbo.orders GROUP BY customer_id
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v WITH SCHEMABINDING AS SELECT name FROM sys.objects
|
< error
| 2720 1 16 1 proc=_v Cannot schema bind view '_v' because it references system object 'sys.objects'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE a AS n FROM customers
|
< error
| 448 1 16 1 proc=_v Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE latin1_general_bin2 AS n FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name COLLATE latin1_general_bin2 AS n FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE database_default AS n FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name COLLATE database_default AS n FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE Latin1_General_CI_AS_WS_KS AS n FROM customers
|
< error
| 448 1 16 1 proc=_v Invalid collation 'Latin1_General_CI_AS_WS_KS'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT name COLLATE Latin1_General_CI_AS_KS_WS AS n FROM customers
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
| VIEW
: string
| dbo._v
: string
| CREATE VIEW _v AS SELECT name COLLATE Latin1_General_CI_AS_KS_WS AS n FROM customers
|
|
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT 'x' COLLATE a AS c FROM nope
|
< error
| 208 1 16 1 proc=_v Invalid object name 'nope'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT id FROM customers ORDER BY name COLLATE a OFFSET 0 ROWS
|
< error
| 448 1 16 1 proc=_v Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE VIEW _v AS SELECT * FROM OPENJSON(N'[]') WITH (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 proc=_v Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10) COLLATE latin1_general_100_ci_as_sc_utf8)
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.t
: string
| CREATE TABLE [dbo].[t] ([c] nvarchar(10) NULL);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10) COLLATE database_default)
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.t
: string
| CREATE TABLE [dbo].[t] ([c] nvarchar(10) NULL);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE nope.t (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE customers (c nope COLLATE a)
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c INT, d AS (c COLLATE a))
|
< error
| 447 0 16 1 Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10), d AS ([c] COLLATE a))
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10), d AS ([c] COLLATE Latin1_General_BIN2))
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.t
: string
| CREATE TABLE [dbo].[t] ([c] nvarchar(10) NULL, [d] AS (([c]) collate Latin1_General_BIN2));
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10) DEFAULT (1 COLLATE a))
|
< error
| 447 0 16 1 Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE nope.t (c NVARCHAR(10) DEFAULT (N'x' COLLATE a))
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10), CONSTRAINT q CHECK ([c] = N'x' COLLATE latin1_general_bin))
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
|
< row
: string
| TABLE
: string
| dbo.t
: string
| CREATE TABLE [dbo].[t] ([c] nvarchar(10) NULL);
| ALTER TABLE [dbo].[t] ADD CONSTRAINT [q] CHECK ([c]=(N'x') collate latin1_general_bin);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TABLE t (c NVARCHAR(10), CHECK (c = N'x' COLLATE a))
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ALTER COLUMN name NVARCHAR(100) COLLATE sql_latin1_general_cp1_cs_as
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL);
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ALTER COLUMN name NVARCHAR(100) COLLATE a
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE nope ALTER COLUMN name NVARCHAR(100) COLLATE a
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE nope ADD c NVARCHAR(10) COLLATE a
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE nope ADD c AS (N'x' COLLATE a)
|
< error
| 4902 1 16 1 Cannot find the object "nope" because it does not exist or you do not have permissions.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD c NVARCHAR(10) DEFAULT (N'x' COLLATE a)
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD c NVARCHAR(10) DEFAULT (N'x' COLLATE latin1_general_ci_as)
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
| dbo.customers
: string
| CREATE TABLE [dbo].[customers] ([id] int NOT NULL, [name] nvarchar(100) NULL, [region] varchar(10) NULL, [born] date NULL, [c] nvarchar(10) NULL DEFAULT ((N'x') collate latin1_general_ci_as));
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
| CREATE TABLE [dbo].[orders] ([id] int NOT NULL, [customer_id] int NULL, [amount] decimal(10,2) NULL, [status] nvarchar(20) NULL, [created] datetime2(3) NULL, [doc] xml NULL);
| ALTER TABLE [dbo].[orders] ADD CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id] ASC);
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
| TABLE
: string
| s.orders
: string
| CREATE TABLE [s].[orders] ([id] int NOT NULL, [note] nvarchar(max) NULL);
|
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD CHECK (id COLLATE Latin1_General_CI_AS = N'x')
|
< error
| 447 0 16 1 Expression type int is invalid for COLLATE clause.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD CHECK (name = N'x' COLLATE a)
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD CONSTRAINT k DEFAULT (N'x' COLLATE a) FOR nope
|
< error
| 1752 0 16 1 Column 'nope' in table 'customers' is invalid for creating a default constraint.
| 1750 0 16 1 Could not create constraint or index. See previous errors.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER TABLE customers ADD CONSTRAINT k DEFAULT (N'x' COLLATE a) FOR region
|
< error
| 448 1 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE TYPE tt AS TABLE (c NVARCHAR(10) COLLATE a)
|
< error
| 448 2 16 1 Invalid collation 'a'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE PARTITION FUNCTION pf (NVARCHAR(10) COLLATE database_default) AS RANGE FOR VALUES ()
|
< error
| 448 3 16 1 Invalid collation 'database_default'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| CREATE DATABASE zz COLLATE database_default
|
< error
| 448 3 16 1 Invalid collation 'database_default'.
> exec
| EXEC [dbo].[_sim_reset]
< ok 0 -
> exec
| CREATE TABLE [dbo].[a] ([b] INT NOT NULL, [a] INT NULL, [d] INT NULL, CONSTRAINT [PK_a] PRIMARY KEY CLUSTERED ([b]));
| CREATE TABLE [dbo].[customers] ([id] INT NOT NULL, [name] NVARCHAR(100) NULL, [region] VARCHAR(10) NULL, [born] DATE NULL, CONSTRAINT [PK_customers] PRIMARY KEY CLUSTERED ([id]));
| CREATE TABLE [dbo].[orders] ([id] INT NOT NULL, [customer_id] INT NULL, [amount] DECIMAL(10,2) NULL, [status] NVARCHAR(20) NULL, [created] DATETIME2(3) NULL, [doc] XML NULL, CONSTRAINT [PK_orders] PRIMARY KEY CLUSTERED ([id]));
< ok 0 -
> exec
| CREATE SCHEMA [s];
< ok 0 -
> exec
| CREATE TABLE [s].[orders] ([id] INT NOT NULL, [note] NVARCHAR(MAX) NULL);
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_double](@x INT) RETURNS INT AS BEGIN RETURN @x * 2 END;
< ok 0 -
> exec
| CREATE FUNCTION [dbo].[f_orders](@c INT) RETURNS TABLE AS RETURN SELECT [id], [amount] FROM [dbo].[orders] WHERE [customer_id] = @c;
< ok 0 -
> exec
| CREATE VIEW [dbo].[v_customers] AS SELECT [id], [name] FROM [dbo].[customers];
< ok 0 -
> exec
| ALTER DATABASE CURRENT COLLATE a
|
< error
| 448 3 16 1 Invalid collation 'a'.
