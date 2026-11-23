-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestComputedNullabilityFidelity default#1
> exec
| CREATE SCHEMA [fid_computed];
< ok 0 -
> exec
| CREATE TABLE [fid_computed].[t] (
|     [id] INT NOT NULL,
|     [a] INT NULL,
|     [c] AS ([a] + 1) PERSISTED,
|     CONSTRAINT [PK_t] PRIMARY KEY ([id])
| );
< ok 0 -
> begin
< ok
> query
|
| 		SELECT schema_name FROM INFORMATION_SCHEMA.SCHEMATA
| 		WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys');
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| fid_computed
> query
|
| 	SELECT
| 		p.major_id AS object_id,
| 		p.minor_id AS column_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.class = 1 AND p.name = 'MS_Description' AND p.minor_id > 0
< result
< column "object_id" "INT"
< column "column_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		so.object_id,
| 		dc.definition,
| 		so.name AS default_name
| 	FROM sys.objects so
| 	INNER JOIN sys.default_constraints dc ON dc.object_id = so.object_id
| 	WHERE so.type = 'D'
< result
< column "object_id" "INT"
< column "definition" "NVARCHAR"
< column "default_name" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		OBJECT_NAME(c.object_id) AS table_name,
| 		c.object_id,
| 		c.column_id,
| 		c.name AS column_name,
| 		t.name AS type_name,
| 		c.is_computed,
| 		cc.definition,
| 		cc.is_persisted,
| 		c.max_length,
| 		c.precision AS precision,
| 		c.scale,
| 		c.collation_name,
| 		c.is_nullable,
| 		c.is_identity,
| 		c.default_object_id,
| 		id.seed_value AS seed_value,
| 		id.increment_value AS increment_value
| 	FROM sys.columns c
| 		LEFT JOIN sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN sys.types t ON c.user_type_id = t.user_type_id
| 		LEFT JOIN sys.objects o ON o.object_id = c.object_id
| 		LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN sys.identity_columns id ON c.object_id = id.object_id AND c.column_id = id.column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, c.object_id ASC, c.column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "column_id" "INT"
< column "column_name" "NVARCHAR"
< column "type_name" "NVARCHAR"
< column "is_computed" "BIT"
< column "definition" "NVARCHAR"
< column "is_persisted" "BIT"
< column "max_length" "SMALLINT"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "collation_name" "NVARCHAR"
< column "is_nullable" "BIT"
< column "is_identity" "BIT"
< column "default_object_id" "INT"
< column "seed_value" "SQL_VARIANT"
< column "increment_value" "SQL_VARIANT"
< row
: string
| fid_computed
: string
| t
: int64 901578250
: int64 1
: string
| id
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 901578250
: int64 2
: string
| a
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 901578250
: int64 3
: string
| c
: string
| int
: bool true
: string
| ([a]+(1))
: bool true
: int64 4
: int64 10
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 	    o.name AS table_name,
| 	    i.name,
| 	    c.name AS column_name,
| 	    ic.is_descending_key,
| 	    i.is_primary_key,
| 	    i.is_unique_constraint,
| 	    i.type_desc,
| 	    CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM
| 	    sys.indexes i
| 	        LEFT JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	        LEFT JOIN sys.columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id
| 	        LEFT JOIN sys.objects co ON co.parent_object_id = i.object_id AND co.name = i.name LEFT JOIN sys.objects o ON o.object_id = i.object_id
| 	        LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	        LEFT JOIN sys.extended_properties p ON p.major_id = co.object_id AND p.class = 1 AND p.name = 'MS_Description'
| 	WHERE i.index_id > 0 AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1) AND o.type IN ('U', 'V') AND s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.name ASC, i.name ASC, ic.key_ordinal ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "is_descending_key" "BIT"
< column "is_primary_key" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< column "comment" "NVARCHAR"
< row
: string
| fid_computed
: string
| t
: string
| PK_t
: string
| id
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		i.object_id,
| 		i.index_id,
| 		i.name,
| 		i.type_desc,
| 		i.is_unique,
| 		i.filter_definition,
| 		xi.secondary_type_desc,
| 		pi.name AS using_index_name,
| 		col.name AS column_name,
| 		ic.key_ordinal,
| 		ic.is_included_column,
| 		ic.is_descending_key
| 	FROM sys.indexes i
| 	INNER JOIN sys.all_objects o ON o.object_id = i.object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	LEFT JOIN sys.xml_indexes xi ON xi.object_id = i.object_id AND xi.index_id = i.index_id
| 	LEFT JOIN sys.indexes pi ON pi.object_id = xi.object_id AND pi.index_id = xi.using_xml_index_id
| 	INNER JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	INNER JOIN sys.all_columns col ON ic.column_id = col.column_id AND ic.object_id = col.object_id
| 	WHERE i.index_id > 0 AND i.is_primary_key = 0 AND i.is_unique_constraint = 0
| 		AND s.name in (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.index_id, ic.key_ordinal, ic.index_column_id
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "name" "NVARCHAR"
< column "type_desc" "NVARCHAR"
< column "is_unique" "BIT"
< column "filter_definition" "NVARCHAR"
< column "secondary_type_desc" "NVARCHAR"
< column "using_index_name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "is_included_column" "BIT"
< column "is_descending_key" "BIT"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
| 		AND ep.major_id IN (SELECT object_id FROM sys.spatial_indexes)
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		object_id,
| 		index_id,
| 		COALESCE(bounding_box_xmin, 0),
| 		COALESCE(bounding_box_ymin, 0),
| 		COALESCE(bounding_box_xmax, 0),
| 		COALESCE(bounding_box_ymax, 0),
| 		COALESCE(level_1_grid_desc, ''),
| 		COALESCE(level_2_grid_desc, ''),
| 		COALESCE(level_3_grid_desc, ''),
| 		COALESCE(level_4_grid_desc, ''),
| 		COALESCE(cells_per_object, 0)
| 	FROM sys.spatial_index_tessellations
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "INT"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		si.object_id,
| 		si.index_id,
| 		i.name AS index_name,
| 		ic.key_ordinal,
| 		col.name AS column_name,
| 		si.spatial_index_type,
| 		si.spatial_index_type_desc,
| 		si.tessellation_scheme,
| 		col.system_type_id,
| 		i.filter_definition,
| 		i.fill_factor,
| 		CASE WHEN i.is_padded = 1 THEN 1 ELSE 0 END AS is_padded,
| 		CASE WHEN i.allow_row_locks = 1 THEN 1 ELSE 0 END AS allow_row_locks,
| 		CASE WHEN i.allow_page_locks = 1 THEN 1 ELSE 0 END AS allow_page_locks
| 	FROM sys.spatial_indexes si
| 	INNER JOIN sys.indexes i ON si.object_id = i.object_id AND si.index_id = i.index_id
| 	INNER JOIN sys.objects o ON i.object_id = o.object_id
| 	INNER JOIN sys.schemas s ON o.schema_id = s.schema_id
| 	INNER JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
| 	INNER JOIN sys.columns col ON ic.object_id = col.object_id AND ic.column_id = col.column_id
| 	WHERE s.name IN (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.name, ic.key_ordinal
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "index_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "column_name" "NVARCHAR"
< column "spatial_index_type" "INT"
< column "spatial_index_type_desc" "NVARCHAR"
< column "tessellation_scheme" "NVARCHAR"
< column "system_type_id" "TINYINT"
< column "filter_definition" "NVARCHAR"
< column "fill_factor" "TINYINT"
< column "is_padded" "INT"
< column "allow_row_locks" "INT"
< column "allow_page_locks" "INT"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.foreign_keys)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		fk.object_id,
| 		fk.name,
| 		OBJECT_SCHEMA_NAME(fk.referenced_object_id) AS referenced_schema,
| 		OBJECT_NAME(fk.referenced_object_id) AS referenced_table,
| 		fk.delete_referential_action,
| 		fk.update_referential_action,
| 		pc.name AS parent_column,
| 		rc.name AS referenced_column
| 	FROM sys.foreign_keys fk
| 	INNER JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
| 	INNER JOIN sys.objects o ON o.object_id = fk.parent_object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	INNER JOIN sys.all_columns pc ON pc.object_id = fkc.parent_object_id AND pc.column_id = fkc.parent_column_id
| 	INNER JOIN sys.all_columns rc ON rc.object_id = fkc.referenced_object_id AND rc.column_id = fkc.referenced_column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.object_id ASC, fk.object_id ASC, fkc.constraint_column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "referenced_schema" "NVARCHAR"
< column "referenced_table" "NVARCHAR"
< column "delete_referential_action" "TINYINT"
< column "update_referential_action" "TINYINT"
< column "parent_column" "NVARCHAR"
< column "referenced_column" "NVARCHAR"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.check_constraints)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 	FROM sys.objects o
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	WHERE s.name in (N'dbo',N'fid_computed') AND o.type IN ('U')
| 	
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 901578250
: string
| fid_computed
: string
| t
> query
|
| 	SELECT
| 		ch.parent_object_id,
| 		ch.name,
| 		ch.object_id,
| 		ch.definition
| 	FROM sys.check_constraints ch
| 	ORDER BY ch.parent_object_id, ch.object_id
< result
< column "parent_object_id" "INT"
< column "name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 		SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 		FROM sys.all_objects o
| 		INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE o.type = 'U'
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 -463397375
: string
| sys
: string
| trace_xe_action_map
< row
: int64 -319884821
: string
| sys
: string
| trace_xe_event_map
< row
: int64 901578250
: string
| fid_computed
: string
| t
> query
|
| 		SELECT major_id, value
| 		FROM sys.extended_properties
| 		WHERE name = 'MS_Description' AND minor_id = 0 AND class = 1
< result
< column "major_id" "INT"
< column "value" "SQL_VARIANT"
> query
|
| 		SELECT ps.object_id, SUM(ps.row_count) AS total_rows
| 		FROM sys.dm_db_partition_stats ps
| 		INNER JOIN sys.tables t ON t.object_id = ps.object_id
| 		INNER JOIN sys.schemas s ON s.schema_id = t.schema_id
| 		WHERE ps.index_id < 2 AND s.name IN (N'dbo',N'fid_computed')
| 		GROUP BY ps.object_id
< result
< column "object_id" "INT"
< column "total_rows" "BIGINT"
< row
: int64 901578250
: int64 0
> query
|
| 		SELECT
| 			t.object_id,
| 			SCHEMA_NAME(t.schema_id),
| 			t.name
| 		FROM sys.tables t
| 		ORDER BY SCHEMA_NAME(t.schema_id), t.name ASC
< result
< column "object_id" "INT"
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: int64 901578250
: string
| fid_computed
: string
| t
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.views)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(v.schema_id) AS schema_name,
| 			v.name AS view_name,
| 			v.object_id,
| 			m.definition
| 		FROM sys.views v
| 		INNER JOIN sys.sql_modules m ON v.object_id = m.object_id
| 		ORDER BY schema_name, view_name;
< result
< column "schema_name" "NVARCHAR"
< column "view_name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name,
| 		seq.name,
| 		tp.name,
| 		seq.precision,
| 		seq.scale,
| 		CAST(seq.start_value AS NVARCHAR(MAX)) AS start_value,
| 		CAST(seq.increment AS NVARCHAR(MAX)) AS increment,
| 		CAST(seq.minimum_value AS NVARCHAR(MAX)) AS minimum_value,
| 		CAST(seq.maximum_value AS NVARCHAR(MAX)) AS maximum_value,
| 		seq.is_cycling,
| 		CASE WHEN seq.is_cached = 0 THEN N'NO CACHE' ELSE CAST(seq.cache_size AS NVARCHAR(MAX)) END AS cache_size,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM
| 		sys.sequences seq
| 	INNER JOIN
| 		sys.schemas s ON s.schema_id = seq.schema_id
| 	INNER JOIN
| 		sys.types tp ON tp.user_type_id = seq.user_type_id
| 	LEFT JOIN sys.extended_properties ep ON ep.major_id = seq.object_id 
| 		AND ep.minor_id = 0 
| 		AND ep.class = 1 
| 		AND ep.name = 'MS_Description'
| 	ORDER BY s.name, seq.name;
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "start_value" "NVARCHAR"
< column "increment" "NVARCHAR"
< column "minimum_value" "NVARCHAR"
< column "maximum_value" "NVARCHAR"
< column "is_cycling" "BIT"
< column "cache_size" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS func_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id 
| 			AND ep.minor_id = 0 
| 			AND ep.class = 1 
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('FN', 'IF', 'TF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, func_name;
< result
< column "schema_name" "NVARCHAR"
< column "func_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS procedure_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id
| 			AND ep.minor_id = 0
| 			AND ep.class = 1
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('P', 'RF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, procedure_name;
< result
< column "schema_name" "NVARCHAR"
< column "procedure_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.triggers WHERE is_disabled = 0 AND is_ms_shipped = 0)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| SELECT
|     st.object_id,
|     st.name,
|     STUFF((
|         SELECT ',' + te.type_desc
|         FROM sys.trigger_events AS te
|         WHERE te.object_id = st.object_id
|         FOR XML PATH('')
|     ), 1, 1, '') AS events,
| CASE
|         WHEN st.is_instead_of_trigger = 1 THEN 'INSTEAD OF'
|         WHEN st.type IN ('TR', 'TA') THEN 'AFTER' -- DML triggers created with FOR or AFTER, and DDL triggers
|         ELSE 'UNKNOWN'
| END AS timing,
| ssm.definition AS body,
| so.name AS parentName,
| ss.name AS schemaName,
| so.type AS objectType
| FROM sys.triggers AS st
| JOIN sys.sql_modules AS ssm ON st.object_id = ssm.object_id
| JOIN sys.objects AS so ON st.parent_id = so.object_id
| JOIN sys.schemas AS ss ON so.schema_id = ss.schema_id
| WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| ORDER BY st.name;
|
< result
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "events" "NVARCHAR"
< column "timing" "VARCHAR"
< column "body" "NVARCHAR"
< column "parentName" "NVARCHAR"
< column "schemaName" "NVARCHAR"
< column "objectType" "CHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(schema_id),
| 			name,
| 			PARSENAME(base_object_name, 4),
| 			PARSENAME(base_object_name, 3),
| 			PARSENAME(base_object_name, 2),
| 			PARSENAME(base_object_name, 1)
| 		FROM sys.synonyms
| 		ORDER BY SCHEMA_NAME(schema_id), name
< result
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> query
| SELECT SCHEMA_NAME(), CONVERT(nvarchar(128), DATABASEPROPERTYEX(DB_NAME(), 'Collation'))
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| dbo
: string
| SQL_Latin1_General_CP1_CI_AS
> commit
< ok
> exec
| DROP TABLE IF EXISTS [fid_computed].[t];
| DROP SCHEMA IF EXISTS [fid_computed];
< ok 0 -
> exec
| CREATE SCHEMA [fid_computed];
< ok 0 -
> exec
| CREATE TABLE [fid_computed].[t] (
|     [id] INT NOT NULL,
|     [a] INT NULL,
|     [c] AS ([a] + 1) PERSISTED NOT NULL,
|     CONSTRAINT [PK_t] PRIMARY KEY ([id])
| );
< ok 0 -
> begin
< ok
> query
|
| 		SELECT schema_name FROM INFORMATION_SCHEMA.SCHEMATA
| 		WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys');
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| fid_computed
> query
|
| 	SELECT
| 		p.major_id AS object_id,
| 		p.minor_id AS column_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.class = 1 AND p.name = 'MS_Description' AND p.minor_id > 0
< result
< column "object_id" "INT"
< column "column_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		so.object_id,
| 		dc.definition,
| 		so.name AS default_name
| 	FROM sys.objects so
| 	INNER JOIN sys.default_constraints dc ON dc.object_id = so.object_id
| 	WHERE so.type = 'D'
< result
< column "object_id" "INT"
< column "definition" "NVARCHAR"
< column "default_name" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		OBJECT_NAME(c.object_id) AS table_name,
| 		c.object_id,
| 		c.column_id,
| 		c.name AS column_name,
| 		t.name AS type_name,
| 		c.is_computed,
| 		cc.definition,
| 		cc.is_persisted,
| 		c.max_length,
| 		c.precision AS precision,
| 		c.scale,
| 		c.collation_name,
| 		c.is_nullable,
| 		c.is_identity,
| 		c.default_object_id,
| 		id.seed_value AS seed_value,
| 		id.increment_value AS increment_value
| 	FROM sys.columns c
| 		LEFT JOIN sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN sys.types t ON c.user_type_id = t.user_type_id
| 		LEFT JOIN sys.objects o ON o.object_id = c.object_id
| 		LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN sys.identity_columns id ON c.object_id = id.object_id AND c.column_id = id.column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, c.object_id ASC, c.column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "column_id" "INT"
< column "column_name" "NVARCHAR"
< column "type_name" "NVARCHAR"
< column "is_computed" "BIT"
< column "definition" "NVARCHAR"
< column "is_persisted" "BIT"
< column "max_length" "SMALLINT"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "collation_name" "NVARCHAR"
< column "is_nullable" "BIT"
< column "is_identity" "BIT"
< column "default_object_id" "INT"
< column "seed_value" "SQL_VARIANT"
< column "increment_value" "SQL_VARIANT"
< row
: string
| fid_computed
: string
| t
: int64 933578364
: int64 1
: string
| id
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 933578364
: int64 2
: string
| a
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 933578364
: int64 3
: string
| c
: string
| int
: bool true
: string
| ([a]+(1))
: bool true
: int64 4
: int64 10
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 	    o.name AS table_name,
| 	    i.name,
| 	    c.name AS column_name,
| 	    ic.is_descending_key,
| 	    i.is_primary_key,
| 	    i.is_unique_constraint,
| 	    i.type_desc,
| 	    CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM
| 	    sys.indexes i
| 	        LEFT JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	        LEFT JOIN sys.columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id
| 	        LEFT JOIN sys.objects co ON co.parent_object_id = i.object_id AND co.name = i.name LEFT JOIN sys.objects o ON o.object_id = i.object_id
| 	        LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	        LEFT JOIN sys.extended_properties p ON p.major_id = co.object_id AND p.class = 1 AND p.name = 'MS_Description'
| 	WHERE i.index_id > 0 AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1) AND o.type IN ('U', 'V') AND s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.name ASC, i.name ASC, ic.key_ordinal ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "is_descending_key" "BIT"
< column "is_primary_key" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< column "comment" "NVARCHAR"
< row
: string
| fid_computed
: string
| t
: string
| PK_t
: string
| id
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		i.object_id,
| 		i.index_id,
| 		i.name,
| 		i.type_desc,
| 		i.is_unique,
| 		i.filter_definition,
| 		xi.secondary_type_desc,
| 		pi.name AS using_index_name,
| 		col.name AS column_name,
| 		ic.key_ordinal,
| 		ic.is_included_column,
| 		ic.is_descending_key
| 	FROM sys.indexes i
| 	INNER JOIN sys.all_objects o ON o.object_id = i.object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	LEFT JOIN sys.xml_indexes xi ON xi.object_id = i.object_id AND xi.index_id = i.index_id
| 	LEFT JOIN sys.indexes pi ON pi.object_id = xi.object_id AND pi.index_id = xi.using_xml_index_id
| 	INNER JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	INNER JOIN sys.all_columns col ON ic.column_id = col.column_id AND ic.object_id = col.object_id
| 	WHERE i.index_id > 0 AND i.is_primary_key = 0 AND i.is_unique_constraint = 0
| 		AND s.name in (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.index_id, ic.key_ordinal, ic.index_column_id
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "name" "NVARCHAR"
< column "type_desc" "NVARCHAR"
< column "is_unique" "BIT"
< column "filter_definition" "NVARCHAR"
< column "secondary_type_desc" "NVARCHAR"
< column "using_index_name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "is_included_column" "BIT"
< column "is_descending_key" "BIT"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
| 		AND ep.major_id IN (SELECT object_id FROM sys.spatial_indexes)
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		object_id,
| 		index_id,
| 		COALESCE(bounding_box_xmin, 0),
| 		COALESCE(bounding_box_ymin, 0),
| 		COALESCE(bounding_box_xmax, 0),
| 		COALESCE(bounding_box_ymax, 0),
| 		COALESCE(level_1_grid_desc, ''),
| 		COALESCE(level_2_grid_desc, ''),
| 		COALESCE(level_3_grid_desc, ''),
| 		COALESCE(level_4_grid_desc, ''),
| 		COALESCE(cells_per_object, 0)
| 	FROM sys.spatial_index_tessellations
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "INT"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		si.object_id,
| 		si.index_id,
| 		i.name AS index_name,
| 		ic.key_ordinal,
| 		col.name AS column_name,
| 		si.spatial_index_type,
| 		si.spatial_index_type_desc,
| 		si.tessellation_scheme,
| 		col.system_type_id,
| 		i.filter_definition,
| 		i.fill_factor,
| 		CASE WHEN i.is_padded = 1 THEN 1 ELSE 0 END AS is_padded,
| 		CASE WHEN i.allow_row_locks = 1 THEN 1 ELSE 0 END AS allow_row_locks,
| 		CASE WHEN i.allow_page_locks = 1 THEN 1 ELSE 0 END AS allow_page_locks
| 	FROM sys.spatial_indexes si
| 	INNER JOIN sys.indexes i ON si.object_id = i.object_id AND si.index_id = i.index_id
| 	INNER JOIN sys.objects o ON i.object_id = o.object_id
| 	INNER JOIN sys.schemas s ON o.schema_id = s.schema_id
| 	INNER JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
| 	INNER JOIN sys.columns col ON ic.object_id = col.object_id AND ic.column_id = col.column_id
| 	WHERE s.name IN (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.name, ic.key_ordinal
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "index_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "column_name" "NVARCHAR"
< column "spatial_index_type" "INT"
< column "spatial_index_type_desc" "NVARCHAR"
< column "tessellation_scheme" "NVARCHAR"
< column "system_type_id" "TINYINT"
< column "filter_definition" "NVARCHAR"
< column "fill_factor" "TINYINT"
< column "is_padded" "INT"
< column "allow_row_locks" "INT"
< column "allow_page_locks" "INT"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.foreign_keys)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		fk.object_id,
| 		fk.name,
| 		OBJECT_SCHEMA_NAME(fk.referenced_object_id) AS referenced_schema,
| 		OBJECT_NAME(fk.referenced_object_id) AS referenced_table,
| 		fk.delete_referential_action,
| 		fk.update_referential_action,
| 		pc.name AS parent_column,
| 		rc.name AS referenced_column
| 	FROM sys.foreign_keys fk
| 	INNER JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
| 	INNER JOIN sys.objects o ON o.object_id = fk.parent_object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	INNER JOIN sys.all_columns pc ON pc.object_id = fkc.parent_object_id AND pc.column_id = fkc.parent_column_id
| 	INNER JOIN sys.all_columns rc ON rc.object_id = fkc.referenced_object_id AND rc.column_id = fkc.referenced_column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.object_id ASC, fk.object_id ASC, fkc.constraint_column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "referenced_schema" "NVARCHAR"
< column "referenced_table" "NVARCHAR"
< column "delete_referential_action" "TINYINT"
< column "update_referential_action" "TINYINT"
< column "parent_column" "NVARCHAR"
< column "referenced_column" "NVARCHAR"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.check_constraints)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 	FROM sys.objects o
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	WHERE s.name in (N'dbo',N'fid_computed') AND o.type IN ('U')
| 	
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 933578364
: string
| fid_computed
: string
| t
> query
|
| 	SELECT
| 		ch.parent_object_id,
| 		ch.name,
| 		ch.object_id,
| 		ch.definition
| 	FROM sys.check_constraints ch
| 	ORDER BY ch.parent_object_id, ch.object_id
< result
< column "parent_object_id" "INT"
< column "name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 		SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 		FROM sys.all_objects o
| 		INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE o.type = 'U'
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 -463397375
: string
| sys
: string
| trace_xe_action_map
< row
: int64 -319884821
: string
| sys
: string
| trace_xe_event_map
< row
: int64 933578364
: string
| fid_computed
: string
| t
> query
|
| 		SELECT major_id, value
| 		FROM sys.extended_properties
| 		WHERE name = 'MS_Description' AND minor_id = 0 AND class = 1
< result
< column "major_id" "INT"
< column "value" "SQL_VARIANT"
> query
|
| 		SELECT ps.object_id, SUM(ps.row_count) AS total_rows
| 		FROM sys.dm_db_partition_stats ps
| 		INNER JOIN sys.tables t ON t.object_id = ps.object_id
| 		INNER JOIN sys.schemas s ON s.schema_id = t.schema_id
| 		WHERE ps.index_id < 2 AND s.name IN (N'dbo',N'fid_computed')
| 		GROUP BY ps.object_id
< result
< column "object_id" "INT"
< column "total_rows" "BIGINT"
< row
: int64 933578364
: int64 0
> query
|
| 		SELECT
| 			t.object_id,
| 			SCHEMA_NAME(t.schema_id),
| 			t.name
| 		FROM sys.tables t
| 		ORDER BY SCHEMA_NAME(t.schema_id), t.name ASC
< result
< column "object_id" "INT"
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: int64 933578364
: string
| fid_computed
: string
| t
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.views)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(v.schema_id) AS schema_name,
| 			v.name AS view_name,
| 			v.object_id,
| 			m.definition
| 		FROM sys.views v
| 		INNER JOIN sys.sql_modules m ON v.object_id = m.object_id
| 		ORDER BY schema_name, view_name;
< result
< column "schema_name" "NVARCHAR"
< column "view_name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name,
| 		seq.name,
| 		tp.name,
| 		seq.precision,
| 		seq.scale,
| 		CAST(seq.start_value AS NVARCHAR(MAX)) AS start_value,
| 		CAST(seq.increment AS NVARCHAR(MAX)) AS increment,
| 		CAST(seq.minimum_value AS NVARCHAR(MAX)) AS minimum_value,
| 		CAST(seq.maximum_value AS NVARCHAR(MAX)) AS maximum_value,
| 		seq.is_cycling,
| 		CASE WHEN seq.is_cached = 0 THEN N'NO CACHE' ELSE CAST(seq.cache_size AS NVARCHAR(MAX)) END AS cache_size,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM
| 		sys.sequences seq
| 	INNER JOIN
| 		sys.schemas s ON s.schema_id = seq.schema_id
| 	INNER JOIN
| 		sys.types tp ON tp.user_type_id = seq.user_type_id
| 	LEFT JOIN sys.extended_properties ep ON ep.major_id = seq.object_id 
| 		AND ep.minor_id = 0 
| 		AND ep.class = 1 
| 		AND ep.name = 'MS_Description'
| 	ORDER BY s.name, seq.name;
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "start_value" "NVARCHAR"
< column "increment" "NVARCHAR"
< column "minimum_value" "NVARCHAR"
< column "maximum_value" "NVARCHAR"
< column "is_cycling" "BIT"
< column "cache_size" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS func_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id 
| 			AND ep.minor_id = 0 
| 			AND ep.class = 1 
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('FN', 'IF', 'TF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, func_name;
< result
< column "schema_name" "NVARCHAR"
< column "func_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS procedure_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id
| 			AND ep.minor_id = 0
| 			AND ep.class = 1
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('P', 'RF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, procedure_name;
< result
< column "schema_name" "NVARCHAR"
< column "procedure_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.triggers WHERE is_disabled = 0 AND is_ms_shipped = 0)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| SELECT
|     st.object_id,
|     st.name,
|     STUFF((
|         SELECT ',' + te.type_desc
|         FROM sys.trigger_events AS te
|         WHERE te.object_id = st.object_id
|         FOR XML PATH('')
|     ), 1, 1, '') AS events,
| CASE
|         WHEN st.is_instead_of_trigger = 1 THEN 'INSTEAD OF'
|         WHEN st.type IN ('TR', 'TA') THEN 'AFTER' -- DML triggers created with FOR or AFTER, and DDL triggers
|         ELSE 'UNKNOWN'
| END AS timing,
| ssm.definition AS body,
| so.name AS parentName,
| ss.name AS schemaName,
| so.type AS objectType
| FROM sys.triggers AS st
| JOIN sys.sql_modules AS ssm ON st.object_id = ssm.object_id
| JOIN sys.objects AS so ON st.parent_id = so.object_id
| JOIN sys.schemas AS ss ON so.schema_id = ss.schema_id
| WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| ORDER BY st.name;
|
< result
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "events" "NVARCHAR"
< column "timing" "VARCHAR"
< column "body" "NVARCHAR"
< column "parentName" "NVARCHAR"
< column "schemaName" "NVARCHAR"
< column "objectType" "CHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(schema_id),
| 			name,
| 			PARSENAME(base_object_name, 4),
| 			PARSENAME(base_object_name, 3),
| 			PARSENAME(base_object_name, 2),
| 			PARSENAME(base_object_name, 1)
| 		FROM sys.synonyms
| 		ORDER BY SCHEMA_NAME(schema_id), name
< result
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> query
| SELECT SCHEMA_NAME(), CONVERT(nvarchar(128), DATABASEPROPERTYEX(DB_NAME(), 'Collation'))
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| dbo
: string
| SQL_Latin1_General_CP1_CI_AS
> commit
< ok
> exec
| DROP TABLE IF EXISTS [fid_computed].[t];
| DROP SCHEMA IF EXISTS [fid_computed];
< ok 0 -
> exec
| CREATE SCHEMA [fid_computed];
< ok 0 -
> exec
| CREATE TABLE [fid_computed].[t] (
|     [id] INT NOT NULL,
|     [a] INT NULL,
|     [c] AS ([a] + 1) PERSISTED,
|     CONSTRAINT [PK_t] PRIMARY KEY ([id])
| );
< ok 0 -
> exec
| ALTER TABLE [fid_computed].[t] DROP COLUMN [c];
| ALTER TABLE [fid_computed].[t] ADD [c] AS ([a]+(1)) PERSISTED NOT NULL;
< ok 0 -
> begin
< ok
> query
|
| 		SELECT schema_name FROM INFORMATION_SCHEMA.SCHEMATA
| 		WHERE schema_name NOT IN ('db_owner', 'db_accessadmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_ddladmin', 'db_denydatareader', 'db_denydatawriter', 'db_securityadmin', 'guest', 'INFORMATION_SCHEMA', 'sys');
< result
< column "schema_name" "NVARCHAR"
< row
: string
| dbo
< row
: string
| fid_computed
> query
|
| 	SELECT
| 		p.major_id AS object_id,
| 		p.minor_id AS column_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.class = 1 AND p.name = 'MS_Description' AND p.minor_id > 0
< result
< column "object_id" "INT"
< column "column_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		so.object_id,
| 		dc.definition,
| 		so.name AS default_name
| 	FROM sys.objects so
| 	INNER JOIN sys.default_constraints dc ON dc.object_id = so.object_id
| 	WHERE so.type = 'D'
< result
< column "object_id" "INT"
< column "definition" "NVARCHAR"
< column "default_name" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		OBJECT_NAME(c.object_id) AS table_name,
| 		c.object_id,
| 		c.column_id,
| 		c.name AS column_name,
| 		t.name AS type_name,
| 		c.is_computed,
| 		cc.definition,
| 		cc.is_persisted,
| 		c.max_length,
| 		c.precision AS precision,
| 		c.scale,
| 		c.collation_name,
| 		c.is_nullable,
| 		c.is_identity,
| 		c.default_object_id,
| 		id.seed_value AS seed_value,
| 		id.increment_value AS increment_value
| 	FROM sys.columns c
| 		LEFT JOIN sys.computed_columns cc ON cc.object_id = c.object_id AND cc.column_id = c.column_id
| 		LEFT JOIN sys.types t ON c.user_type_id = t.user_type_id
| 		LEFT JOIN sys.objects o ON o.object_id = c.object_id
| 		LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		LEFT JOIN sys.identity_columns id ON c.object_id = id.object_id AND c.column_id = id.column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, c.object_id ASC, c.column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "column_id" "INT"
< column "column_name" "NVARCHAR"
< column "type_name" "NVARCHAR"
< column "is_computed" "BIT"
< column "definition" "NVARCHAR"
< column "is_persisted" "BIT"
< column "max_length" "SMALLINT"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "collation_name" "NVARCHAR"
< column "is_nullable" "BIT"
< column "is_identity" "BIT"
< column "default_object_id" "INT"
< column "seed_value" "SQL_VARIANT"
< column "increment_value" "SQL_VARIANT"
< row
: string
| fid_computed
: string
| t
: int64 965578478
: int64 1
: string
| id
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 965578478
: int64 2
: string
| a
: string
| int
: bool false
: null
: null
: int64 4
: int64 10
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| fid_computed
: string
| t
: int64 965578478
: int64 4
: string
| c
: string
| int
: bool true
: string
| ([a]+(1))
: bool true
: int64 4
: int64 10
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 	    o.name AS table_name,
| 	    i.name,
| 	    c.name AS column_name,
| 	    ic.is_descending_key,
| 	    i.is_primary_key,
| 	    i.is_unique_constraint,
| 	    i.type_desc,
| 	    CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM
| 	    sys.indexes i
| 	        LEFT JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	        LEFT JOIN sys.columns c ON c.object_id = ic.object_id AND c.column_id = ic.column_id
| 	        LEFT JOIN sys.objects co ON co.parent_object_id = i.object_id AND co.name = i.name LEFT JOIN sys.objects o ON o.object_id = i.object_id
| 	        LEFT JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	        LEFT JOIN sys.extended_properties p ON p.major_id = co.object_id AND p.class = 1 AND p.name = 'MS_Description'
| 	WHERE i.index_id > 0 AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1) AND o.type IN ('U', 'V') AND s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.name ASC, i.name ASC, ic.key_ordinal ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "is_descending_key" "BIT"
< column "is_primary_key" "BIT"
< column "is_unique_constraint" "BIT"
< column "type_desc" "NVARCHAR"
< column "comment" "NVARCHAR"
< row
: string
| fid_computed
: string
| t
: string
| PK_t
: string
| id
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		i.object_id,
| 		i.index_id,
| 		i.name,
| 		i.type_desc,
| 		i.is_unique,
| 		i.filter_definition,
| 		xi.secondary_type_desc,
| 		pi.name AS using_index_name,
| 		col.name AS column_name,
| 		ic.key_ordinal,
| 		ic.is_included_column,
| 		ic.is_descending_key
| 	FROM sys.indexes i
| 	INNER JOIN sys.all_objects o ON o.object_id = i.object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	LEFT JOIN sys.xml_indexes xi ON xi.object_id = i.object_id AND xi.index_id = i.index_id
| 	LEFT JOIN sys.indexes pi ON pi.object_id = xi.object_id AND pi.index_id = xi.using_xml_index_id
| 	INNER JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	INNER JOIN sys.all_columns col ON ic.column_id = col.column_id AND ic.object_id = col.object_id
| 	WHERE i.index_id > 0 AND i.is_primary_key = 0 AND i.is_unique_constraint = 0
| 		AND s.name in (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.index_id, ic.key_ordinal, ic.index_column_id
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "name" "NVARCHAR"
< column "type_desc" "NVARCHAR"
< column "is_unique" "BIT"
< column "filter_definition" "NVARCHAR"
< column "secondary_type_desc" "NVARCHAR"
< column "using_index_name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "is_included_column" "BIT"
< column "is_descending_key" "BIT"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		ep.minor_id AS index_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.class = 7 AND ep.name = 'MS_Description' AND ep.minor_id > 0
| 		AND ep.major_id IN (SELECT object_id FROM sys.spatial_indexes)
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		object_id,
| 		index_id,
| 		COALESCE(bounding_box_xmin, 0),
| 		COALESCE(bounding_box_ymin, 0),
| 		COALESCE(bounding_box_xmax, 0),
| 		COALESCE(bounding_box_ymax, 0),
| 		COALESCE(level_1_grid_desc, ''),
| 		COALESCE(level_2_grid_desc, ''),
| 		COALESCE(level_3_grid_desc, ''),
| 		COALESCE(level_4_grid_desc, ''),
| 		COALESCE(cells_per_object, 0)
| 	FROM sys.spatial_index_tessellations
< result
< column "object_id" "INT"
< column "index_id" "INT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "FLOAT"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "INT"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		si.object_id,
| 		si.index_id,
| 		i.name AS index_name,
| 		ic.key_ordinal,
| 		col.name AS column_name,
| 		si.spatial_index_type,
| 		si.spatial_index_type_desc,
| 		si.tessellation_scheme,
| 		col.system_type_id,
| 		i.filter_definition,
| 		i.fill_factor,
| 		CASE WHEN i.is_padded = 1 THEN 1 ELSE 0 END AS is_padded,
| 		CASE WHEN i.allow_row_locks = 1 THEN 1 ELSE 0 END AS allow_row_locks,
| 		CASE WHEN i.allow_page_locks = 1 THEN 1 ELSE 0 END AS allow_page_locks
| 	FROM sys.spatial_indexes si
| 	INNER JOIN sys.indexes i ON si.object_id = i.object_id AND si.index_id = i.index_id
| 	INNER JOIN sys.objects o ON i.object_id = o.object_id
| 	INNER JOIN sys.schemas s ON o.schema_id = s.schema_id
| 	INNER JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
| 	INNER JOIN sys.columns col ON ic.object_id = col.object_id AND ic.column_id = col.column_id
| 	WHERE s.name IN (N'dbo',N'fid_computed') AND o.type IN ('U', 'S', 'V')
| 	ORDER BY s.name, o.name, i.name, ic.key_ordinal
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "index_id" "INT"
< column "index_name" "NVARCHAR"
< column "key_ordinal" "TINYINT"
< column "column_name" "NVARCHAR"
< column "spatial_index_type" "INT"
< column "spatial_index_type_desc" "NVARCHAR"
< column "tessellation_scheme" "NVARCHAR"
< column "system_type_id" "TINYINT"
< column "filter_definition" "NVARCHAR"
< column "fill_factor" "TINYINT"
< column "is_padded" "INT"
< column "allow_row_locks" "INT"
< column "allow_page_locks" "INT"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.foreign_keys)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name AS schema_name,
| 		o.name AS table_name,
| 		fk.object_id,
| 		fk.name,
| 		OBJECT_SCHEMA_NAME(fk.referenced_object_id) AS referenced_schema,
| 		OBJECT_NAME(fk.referenced_object_id) AS referenced_table,
| 		fk.delete_referential_action,
| 		fk.update_referential_action,
| 		pc.name AS parent_column,
| 		rc.name AS referenced_column
| 	FROM sys.foreign_keys fk
| 	INNER JOIN sys.foreign_key_columns fkc ON fkc.constraint_object_id = fk.object_id
| 	INNER JOIN sys.objects o ON o.object_id = fk.parent_object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	INNER JOIN sys.all_columns pc ON pc.object_id = fkc.parent_object_id AND pc.column_id = fkc.parent_column_id
| 	INNER JOIN sys.all_columns rc ON rc.object_id = fkc.referenced_object_id AND rc.column_id = fkc.referenced_column_id
| 	WHERE s.name in (N'dbo',N'fid_computed')
| 	ORDER BY s.name ASC, o.object_id ASC, fk.object_id ASC, fkc.constraint_column_id ASC
| 	
< result
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "referenced_schema" "NVARCHAR"
< column "referenced_table" "NVARCHAR"
< column "delete_referential_action" "TINYINT"
< column "update_referential_action" "TINYINT"
< column "parent_column" "NVARCHAR"
< column "referenced_column" "NVARCHAR"
> query
|
| 	SELECT 
| 		p.major_id AS object_id,
| 		CAST(p.[value] AS nvarchar(4000)) AS comment
| 	FROM sys.extended_properties p
| 	WHERE p.minor_id = 0 AND p.name = 'MS_Description'
| 		AND p.major_id IN (SELECT object_id FROM sys.check_constraints)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 	FROM sys.objects o
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	WHERE s.name in (N'dbo',N'fid_computed') AND o.type IN ('U')
| 	
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 965578478
: string
| fid_computed
: string
| t
> query
|
| 	SELECT
| 		ch.parent_object_id,
| 		ch.name,
| 		ch.object_id,
| 		ch.definition
| 	FROM sys.check_constraints ch
| 	ORDER BY ch.parent_object_id, ch.object_id
< result
< column "parent_object_id" "INT"
< column "name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 		SELECT o.object_id, s.name AS schema_name, o.name AS table_name
| 		FROM sys.all_objects o
| 		INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 		WHERE o.type = 'U'
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 -463397375
: string
| sys
: string
| trace_xe_action_map
< row
: int64 -319884821
: string
| sys
: string
| trace_xe_event_map
< row
: int64 965578478
: string
| fid_computed
: string
| t
> query
|
| 		SELECT major_id, value
| 		FROM sys.extended_properties
| 		WHERE name = 'MS_Description' AND minor_id = 0 AND class = 1
< result
< column "major_id" "INT"
< column "value" "SQL_VARIANT"
> query
|
| 		SELECT ps.object_id, SUM(ps.row_count) AS total_rows
| 		FROM sys.dm_db_partition_stats ps
| 		INNER JOIN sys.tables t ON t.object_id = ps.object_id
| 		INNER JOIN sys.schemas s ON s.schema_id = t.schema_id
| 		WHERE ps.index_id < 2 AND s.name IN (N'dbo',N'fid_computed')
| 		GROUP BY ps.object_id
< result
< column "object_id" "INT"
< column "total_rows" "BIGINT"
< row
: int64 965578478
: int64 0
> query
|
| 		SELECT
| 			t.object_id,
| 			SCHEMA_NAME(t.schema_id),
| 			t.name
| 		FROM sys.tables t
| 		ORDER BY SCHEMA_NAME(t.schema_id), t.name ASC
< result
< column "object_id" "INT"
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< row
: int64 965578478
: string
| fid_computed
: string
| t
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.views)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(v.schema_id) AS schema_name,
| 			v.name AS view_name,
| 			v.object_id,
| 			m.definition
| 		FROM sys.views v
| 		INNER JOIN sys.sql_modules m ON v.object_id = m.object_id
| 		ORDER BY schema_name, view_name;
< result
< column "schema_name" "NVARCHAR"
< column "view_name" "NVARCHAR"
< column "object_id" "INT"
< column "definition" "NVARCHAR"
> query
|
| 	SELECT
| 		s.name,
| 		seq.name,
| 		tp.name,
| 		seq.precision,
| 		seq.scale,
| 		CAST(seq.start_value AS NVARCHAR(MAX)) AS start_value,
| 		CAST(seq.increment AS NVARCHAR(MAX)) AS increment,
| 		CAST(seq.minimum_value AS NVARCHAR(MAX)) AS minimum_value,
| 		CAST(seq.maximum_value AS NVARCHAR(MAX)) AS maximum_value,
| 		seq.is_cycling,
| 		CASE WHEN seq.is_cached = 0 THEN N'NO CACHE' ELSE CAST(seq.cache_size AS NVARCHAR(MAX)) END AS cache_size,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM
| 		sys.sequences seq
| 	INNER JOIN
| 		sys.schemas s ON s.schema_id = seq.schema_id
| 	INNER JOIN
| 		sys.types tp ON tp.user_type_id = seq.user_type_id
| 	LEFT JOIN sys.extended_properties ep ON ep.major_id = seq.object_id 
| 		AND ep.minor_id = 0 
| 		AND ep.class = 1 
| 		AND ep.name = 'MS_Description'
| 	ORDER BY s.name, seq.name;
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "precision" "TINYINT"
< column "scale" "TINYINT"
< column "start_value" "NVARCHAR"
< column "increment" "NVARCHAR"
< column "minimum_value" "NVARCHAR"
< column "maximum_value" "NVARCHAR"
< column "is_cycling" "BIT"
< column "cache_size" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS func_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id 
| 			AND ep.minor_id = 0 
| 			AND ep.class = 1 
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('FN', 'IF', 'TF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, func_name;
< result
< column "schema_name" "NVARCHAR"
< column "func_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS procedure_name,
| 		asm.definition,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
| 		LEFT JOIN sys.extended_properties ep ON ep.major_id = ao.object_id
| 			AND ep.minor_id = 0
| 			AND ep.class = 1
| 			AND ep.name = 'MS_Description'
| 	WHERE ao.type IN ('P', 'RF')
| 		AND ao.is_ms_shipped = 0 AND
| 		(
| 			SELECT major_id
| 			FROM sys.extended_properties
| 			WHERE major_id = ao.object_id
| 				AND minor_id = 0
| 				AND class = 1
| 				AND name = 'microsoft_database_tools_support'
| 		) IS NULL
| 	ORDER BY schema_name, procedure_name;
< result
< column "schema_name" "NVARCHAR"
< column "procedure_name" "NVARCHAR"
< column "definition" "NVARCHAR"
< column "comment" "NVARCHAR"
> query
|
| 	SELECT 
| 		ep.major_id AS object_id,
| 		CAST(ep.value AS NVARCHAR(MAX)) AS comment
| 	FROM sys.extended_properties ep
| 	WHERE ep.minor_id = 0 AND ep.class = 1 AND ep.name = 'MS_Description'
| 		AND ep.major_id IN (SELECT object_id FROM sys.triggers WHERE is_disabled = 0 AND is_ms_shipped = 0)
< result
< column "object_id" "INT"
< column "comment" "NVARCHAR"
> query
|
| SELECT
|     st.object_id,
|     st.name,
|     STUFF((
|         SELECT ',' + te.type_desc
|         FROM sys.trigger_events AS te
|         WHERE te.object_id = st.object_id
|         FOR XML PATH('')
|     ), 1, 1, '') AS events,
| CASE
|         WHEN st.is_instead_of_trigger = 1 THEN 'INSTEAD OF'
|         WHEN st.type IN ('TR', 'TA') THEN 'AFTER' -- DML triggers created with FOR or AFTER, and DDL triggers
|         ELSE 'UNKNOWN'
| END AS timing,
| ssm.definition AS body,
| so.name AS parentName,
| ss.name AS schemaName,
| so.type AS objectType
| FROM sys.triggers AS st
| JOIN sys.sql_modules AS ssm ON st.object_id = ssm.object_id
| JOIN sys.objects AS so ON st.parent_id = so.object_id
| JOIN sys.schemas AS ss ON so.schema_id = ss.schema_id
| WHERE st.is_disabled = 0 AND st.is_ms_shipped = 0 AND st.parent_id <> 0 AND so.type IN ('U', 'V')
| ORDER BY st.name;
|
< result
< column "object_id" "INT"
< column "name" "NVARCHAR"
< column "events" "NVARCHAR"
< column "timing" "VARCHAR"
< column "body" "NVARCHAR"
< column "parentName" "NVARCHAR"
< column "schemaName" "NVARCHAR"
< column "objectType" "CHAR"
> query
|
| 		SELECT
| 			SCHEMA_NAME(schema_id),
| 			name,
| 			PARSENAME(base_object_name, 4),
| 			PARSENAME(base_object_name, 3),
| 			PARSENAME(base_object_name, 2),
| 			PARSENAME(base_object_name, 1)
| 		FROM sys.synonyms
| 		ORDER BY SCHEMA_NAME(schema_id), name
< result
< column "" "NVARCHAR"
< column "name" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> query
| SELECT SCHEMA_NAME(), CONVERT(nvarchar(128), DATABASEPROPERTYEX(DB_NAME(), 'Collation'))
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| dbo
: string
| SQL_Latin1_General_CP1_CI_AS
> commit
< ok
