-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestSyncDatabase default#1
> exec
| CREATE SCHEMA app;
< ok 0 -
> exec
| CREATE TABLE app.accounts (
| 	id BIGINT IDENTITY(1, 2) NOT NULL,
| 	email NVARCHAR(190) NOT NULL,
| 	bio NVARCHAR(MAX) NULL,
| 	status VARCHAR(10) NOT NULL CONSTRAINT df_accounts_status DEFAULT 'open',
| 	balance DECIMAL(12, 4) NOT NULL CONSTRAINT df_accounts_balance DEFAULT 0.0000,
| 	flags VARBINARY(16) NULL,
| 	token UNIQUEIDENTIFIER NULL,
| 	created_at DATETIME2(3) NOT NULL CONSTRAINT df_accounts_created DEFAULT SYSUTCDATETIME(),
| 	offset_at DATETIMEOFFSET(0) NULL,
| 	xml_doc XML NULL,
| 	updated_at DATETIME2(3) NULL,
| 	email_lower AS LOWER(email) PERSISTED,
| 	PRIMARY KEY (id),
| 	CONSTRAINT uk_accounts_email UNIQUE (email),
| 	CONSTRAINT chk_accounts_balance CHECK (balance >= 0)
| );
< ok 0 -
> exec
| EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'customer accounts', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'TABLE', @level1name = N'accounts';
< ok 2 -
> exec
| EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'owner''s note', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'TABLE', @level1name = N'accounts', @level2type = N'COLUMN', @level2name = N'bio';
< ok 2 -
> exec
| CREATE INDEX idx_accounts_status ON app.accounts (status) INCLUDE (email) WHERE status = 'open';
< ok 0 -
> exec
| CREATE PRIMARY XML INDEX pxml_accounts_doc ON app.accounts (xml_doc);
< ok 0 -
> exec
| CREATE TABLE app.orders (
| 	id BIGINT IDENTITY NOT NULL PRIMARY KEY,
| 	account_id BIGINT NOT NULL CONSTRAINT fk_orders_account REFERENCES app.accounts (id) ON DELETE CASCADE,
| 	total INT NOT NULL,
| 	note NVARCHAR(100) NULL
| );
< ok 0 -
> exec
| CREATE TABLE app.metrics (
| 	ts INT NOT NULL,
| 	val BIGINT NOT NULL,
| 	CONSTRAINT pk_metrics PRIMARY KEY (ts)
| );
< ok 0 -
> exec
| CREATE NONCLUSTERED COLUMNSTORE INDEX ncc_metrics ON app.metrics (ts, val);
< ok 0 -
> exec
| CREATE SEQUENCE app.seq_orders AS BIGINT START WITH 10 INCREMENT BY 2 MINVALUE 1 MAXVALUE 10000 CYCLE CACHE 50;
< ok 0 -
> exec
| EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'order sequence', @level0type = N'SCHEMA', @level0name = N'app', @level1type = N'SEQUENCE', @level1name = N'seq_orders';
< ok 2 -
> exec
| CREATE VIEW app.v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM app.accounts WHERE status = 'open';
< ok 0 -
> exec
| CREATE FUNCTION app.fn_format_tag (@prefix NVARCHAR(16), @id BIGINT) RETURNS NVARCHAR(64) AS
| BEGIN
| 	RETURN @prefix + N'-' + CAST(@id AS NVARCHAR(20));
| END;
< ok 0 -
> exec
| CREATE FUNCTION app.fn_orders_of (@account BIGINT) RETURNS TABLE AS
| 	RETURN (SELECT id, total FROM app.orders WHERE account_id = @account);
< ok 0 -
> exec
| CREATE PROCEDURE app.sp_order_total @id BIGINT, @total INT OUTPUT AS
| BEGIN
| 	SELECT @total = total FROM app.orders WHERE id = @id;
| END;
< ok 0 -
> exec
| CREATE PROCEDURE app.sp_account_orders @account BIGINT AS
| BEGIN
| 	SELECT id, total FROM app.orders WHERE account_id = @account;
| 	SELECT COUNT(*) AS n FROM app.orders WHERE account_id = @account;
| END;
< ok 0 -
> exec
| CREATE TRIGGER app.trg_open_accounts_io ON app.v_open_accounts INSTEAD OF INSERT AS
| BEGIN
| 	SET NOCOUNT ON;
| END;
< ok 0 -
> exec
| CREATE TYPE app.big_id FROM BIGINT NOT NULL;
< ok 0 -
> exec
| CREATE SEQUENCE app.seq_alias AS app.big_id START WITH 1 INCREMENT BY 1;
< ok 0 -
> exec
| CREATE TRIGGER app.trg_orders_bi ON app.orders AFTER INSERT AS
| BEGIN
| 	UPDATE app.orders SET total = 0 WHERE total < 0 AND id IN (SELECT id FROM inserted);
| END;
< ok 0 -
> exec
| CREATE SYNONYM app.syn_accounts FOR app.accounts;
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
| app
< row
: string
| dbo
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
< row
: int64 901578250
: int64 3
: string
| owner's note
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
< row
: int64 949578421
: string
| ('open')
: string
| df_accounts_status
< row
: int64 965578478
: string
| ((0.0000))
: string
| df_accounts_balance
< row
: int64 981578535
: string
| (sysutcdatetime())
: string
| df_accounts_created
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
| 	WHERE s.name in (N'app',N'dbo')
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
| app
: string
| accounts
: int64 901578250
: int64 1
: string
| id
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool true
: int64 0
: int64 1
: int64 2
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 2
: string
| email
: string
| nvarchar
: bool false
: null
: null
: int64 380
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 3
: string
| bio
: string
| nvarchar
: bool false
: null
: null
: int64 -1
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 4
: string
| status
: string
| varchar
: bool false
: null
: null
: int64 10
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool false
: bool false
: int64 949578421
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 5
: string
| balance
: string
| decimal
: bool false
: null
: null
: int64 9
: int64 12
: int64 4
: null
: bool false
: bool false
: int64 965578478
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 6
: string
| flags
: string
| varbinary
: bool false
: null
: null
: int64 16
: int64 0
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 7
: string
| token
: string
| uniqueidentifier
: bool false
: null
: null
: int64 16
: int64 0
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 8
: string
| created_at
: string
| datetime2
: bool false
: null
: null
: int64 7
: int64 23
: int64 3
: null
: bool false
: bool false
: int64 981578535
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 9
: string
| offset_at
: string
| datetimeoffset
: bool false
: null
: null
: int64 8
: int64 26
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 10
: string
| xml_doc
: string
| xml
: bool false
: null
: null
: int64 -1
: int64 0
: int64 0
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 11
: string
| updated_at
: string
| datetime2
: bool false
: null
: null
: int64 7
: int64 23
: int64 3
: null
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 12
: string
| email_lower
: string
| nvarchar
: bool true
: string
| (lower([email]))
: bool true
: int64 380
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| orders
: int64 1029578706
: int64 1
: string
| id
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool true
: int64 0
: int64 1
: int64 1
< row
: string
| app
: string
| orders
: int64 1029578706
: int64 2
: string
| account_id
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| orders
: int64 1029578706
: int64 3
: string
| total
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
| app
: string
| orders
: int64 1029578706
: int64 4
: string
| note
: string
| nvarchar
: bool false
: null
: null
: int64 200
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool true
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| metrics
: int64 1077578877
: int64 1
: string
| ts
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
| app
: string
| metrics
: int64 1077578877
: int64 2
: string
| val
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| v_open_accounts
: int64 1125579048
: int64 1
: string
| acct_id
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool true
: int64 0
: null
: null
< row
: string
| app
: string
| v_open_accounts
: int64 1125579048
: int64 2
: string
| acct_email
: string
| nvarchar
: bool false
: null
: null
: int64 380
: int64 0
: int64 0
: string
| SQL_Latin1_General_CP1_CI_AS
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| fn_orders_of
: int64 1157579162
: int64 1
: string
| id
: string
| bigint
: bool false
: null
: null
: int64 8
: int64 19
: int64 0
: null
: bool false
: bool false
: int64 0
: null
: null
< row
: string
| app
: string
| fn_orders_of
: int64 1157579162
: int64 2
: string
| total
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
| 	WHERE i.index_id > 0 AND (i.is_primary_key = 1 OR i.is_unique_constraint = 1) AND o.type IN ('U', 'V') AND s.name in (N'app',N'dbo')
| 	ORDER BY s.name ASC, i.name ASC, ic.key_ordinal ASC
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
| app
: string
| accounts
: string
| PK__accounts__standin1
: string
| id
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
< row
: string
| app
: string
| orders
: string
| PK__orders__standin2
: string
| id
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
< row
: string
| app
: string
| metrics
: string
| pk_metrics
: string
| ts
: bool false
: bool true
: bool false
: string
| CLUSTERED
: null
< row
: string
| app
: string
| accounts
: string
| uk_accounts_email
: string
| email
: bool false
: bool false
: bool true
: string
| NONCLUSTERED
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
| 		xi.secondary_type_desc,
| 		pi.name AS using_index_name,
| 		col.name AS column_name,
| 		ic.is_descending_key
| 	FROM sys.indexes i
| 	INNER JOIN sys.all_objects o ON o.object_id = i.object_id
| 	INNER JOIN sys.schemas s ON s.schema_id = o.schema_id
| 	LEFT JOIN sys.xml_indexes xi ON xi.object_id = i.object_id AND xi.index_id = i.index_id
| 	LEFT JOIN sys.indexes pi ON pi.object_id = xi.object_id AND pi.index_id = xi.using_xml_index_id
| 	INNER JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
| 	INNER JOIN sys.all_columns col ON ic.column_id = col.column_id AND ic.object_id = col.object_id
| 	WHERE i.index_id > 0 AND i.is_primary_key = 0 AND i.is_unique_constraint = 0
| 		AND s.name in (N'app',N'dbo') AND o.type IN ('U', 'S', 'V')
| 		-- A rowstore index's INCLUDE columns carry key_ordinal 0, and belong to no
| 		-- key; a columnstore or XML index states every one of its columns that way.
| 		AND (ic.key_ordinal > 0 OR i.type_desc IN ('CLUSTERED COLUMNSTORE', 'NONCLUSTERED COLUMNSTORE', 'XML'))
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
< column "secondary_type_desc" "NVARCHAR"
< column "using_index_name" "NVARCHAR"
< column "column_name" "NVARCHAR"
< column "is_descending_key" "BIT"
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 3
: string
| idx_accounts_status
: string
| NONCLUSTERED
: bool false
: null
: null
: string
| status
: bool false
< row
: string
| app
: string
| accounts
: int64 901578250
: int64 256000
: string
| pxml_accounts_doc
: string
| XML
: bool false
: null
: null
: string
| xml_doc
: bool false
< row
: string
| app
: string
| metrics
: int64 1077578877
: int64 2
: string
| ncc_metrics
: string
| NONCLUSTERED COLUMNSTORE
: bool false
: null
: null
: string
| ts
: bool false
< row
: string
| app
: string
| metrics
: int64 1077578877
: int64 2
: string
| ncc_metrics
: string
| NONCLUSTERED COLUMNSTORE
: bool false
: null
: null
: string
| val
: bool false
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
| 	WHERE s.name IN (N'app',N'dbo') AND o.type IN ('U', 'S', 'V')
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
| 	WHERE s.name in (N'app',N'dbo')
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
< row
: string
| app
: string
| orders
: int64 1061578820
: string
| fk_orders_account
: string
| app
: string
| accounts
: int64 1
: int64 0
: string
| account_id
: string
| id
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
| 	WHERE s.name in (N'app',N'dbo') AND o.type IN ('U')
| 	
< result
< column "object_id" "INT"
< column "schema_name" "NVARCHAR"
< column "table_name" "NVARCHAR"
< row
: int64 901578250
: string
| app
: string
| accounts
< row
: int64 1029578706
: string
| app
: string
| orders
< row
: int64 1077578877
: string
| app
: string
| metrics
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
< row
: int64 901578250
: string
| chk_accounts_balance
: int64 997578592
: string
| ([balance]>=(0))
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
| app
: string
| accounts
< row
: int64 1029578706
: string
| app
: string
| orders
< row
: int64 1077578877
: string
| app
: string
| metrics
> query
|
| 		SELECT major_id, value
| 		FROM sys.extended_properties
| 		WHERE name = 'MS_Description' AND minor_id = 0 AND class = 1
< result
< column "major_id" "INT"
< column "value" "SQL_VARIANT"
< row
: int64 901578250
: string
| customer accounts
< row
: int64 1109578991
: string
| order sequence
> query
|
| 		SELECT object_id, SUM(row_count) AS total_rows
| 		FROM sys.dm_db_partition_stats
| 		WHERE index_id < 2
| 		GROUP BY object_id
< result
< column "object_id" "INT"
< column "total_rows" "BIGINT"
< row
: int64 3
: int64 1530
< row
: int64 5
: int64 196
< row
: int64 6
: int64 0
< row
: int64 7
: int64 234
< row
: int64 8
: int64 2
< row
: int64 9
: int64 3
< row
: int64 17
: int64 0
< row
: int64 18
: int64 0
< row
: int64 19
: int64 0
< row
: int64 20
: int64 2
< row
: int64 21
: int64 0
< row
: int64 22
: int64 0
< row
: int64 23
: int64 1
< row
: int64 24
: int64 2
< row
: int64 25
: int64 0
< row
: int64 27
: int64 14
< row
: int64 28
: int64 0
< row
: int64 29
: int64 229
< row
: int64 34
: int64 2800
< row
: int64 35
: int64 0
< row
: int64 36
: int64 0
< row
: int64 37
: int64 0
< row
: int64 38
: int64 0
< row
: int64 40
: int64 0
< row
: int64 41
: int64 1264
< row
: int64 42
: int64 0
< row
: int64 43
: int64 0
< row
: int64 44
: int64 1
< row
: int64 45
: int64 0
< row
: int64 46
: int64 0
< row
: int64 47
: int64 0
< row
: int64 48
: int64 0
< row
: int64 49
: int64 3
< row
: int64 50
: int64 35
< row
: int64 51
: int64 0
< row
: int64 54
: int64 229
< row
: int64 55
: int64 458
< row
: int64 56
: int64 0
< row
: int64 57
: int64 0
< row
: int64 58
: int64 23
< row
: int64 59
: int64 0
< row
: int64 60
: int64 288
< row
: int64 62
: int64 0
< row
: int64 63
: int64 0
< row
: int64 64
: int64 27
< row
: int64 65
: int64 0
< row
: int64 67
: int64 0
< row
: int64 68
: int64 0
< row
: int64 69
: int64 1
< row
: int64 71
: int64 0
< row
: int64 72
: int64 0
< row
: int64 73
: int64 0
< row
: int64 74
: int64 210
< row
: int64 75
: int64 123
< row
: int64 78
: int64 0
< row
: int64 79
: int64 0
< row
: int64 80
: int64 0
< row
: int64 82
: int64 0
< row
: int64 84
: int64 0
< row
: int64 85
: int64 0
< row
: int64 86
: int64 0
< row
: int64 87
: int64 0
< row
: int64 89
: int64 0
< row
: int64 90
: int64 98
< row
: int64 91
: int64 100
< row
: int64 92
: int64 112
< row
: int64 93
: int64 19
< row
: int64 94
: int64 0
< row
: int64 95
: int64 0
< row
: int64 96
: int64 0
< row
: int64 97
: int64 3
< row
: int64 98
: int64 6
< row
: int64 101575400
: int64 0
< row
: int64 117575457
: int64 0
< row
: int64 133575514
: int64 0
< row
: int64 165575628
: int64 0
< row
: int64 181575685
: int64 0
< row
: int64 197575742
: int64 0
< row
: int64 213575799
: int64 0
< row
: int64 229575856
: int64 1
< row
: int64 245575913
: int64 3
< row
: int64 261575970
: int64 0
< row
: int64 277576027
: int64 0
< row
: int64 293576084
: int64 0
< row
: int64 341576255
: int64 0
< row
: int64 357576312
: int64 0
< row
: int64 373576369
: int64 0
< row
: int64 389576426
: int64 0
< row
: int64 405576483
: int64 0
< row
: int64 565577053
: int64 0
< row
: int64 581577110
: int64 0
< row
: int64 597577167
: int64 0
< row
: int64 613577224
: int64 0
< row
: int64 629577281
: int64 0
< row
: int64 645577338
: int64 0
< row
: int64 661577395
: int64 0
< row
: int64 677577452
: int64 0
< row
: int64 693577509
: int64 0
< row
: int64 885578193
: int64 0
< row
: int64 901578250
: int64 0
< row
: int64 1013578649
: int64 0
< row
: int64 1029578706
: int64 0
< row
: int64 1077578877
: int64 0
< row
: int64 1993058136
: int64 0
< row
: int64 2025058250
: int64 0
< row
: int64 2057058364
: int64 0
< row
: int64 2073058421
: int64 0
< row
: int64 2089058478
: int64 0
< row
: int64 2105058535
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
| app
: string
| accounts
< row
: int64 1077578877
: string
| app
: string
| metrics
< row
: int64 1029578706
: string
| app
: string
| orders
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
< row
: string
| app
: string
| v_open_accounts
: int64 1125579048
: string
| CREATE VIEW app.v_open_accounts AS
| 	SELECT id AS acct_id, email AS acct_email FROM app.accounts WHERE status = 'open';
> query
|
| 	SELECT
| 		s.name,
| 		seq.name,
| 		tp.name,
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
< column "comment" "NVARCHAR"
< row
: string
| app
: string
| seq_alias
: string
| big_id
: null
< row
: string
| app
: string
| seq_orders
: string
| bigint
: string
| order sequence
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
< row
: string
| app
: string
| fn_format_tag
: string
| CREATE FUNCTION app.fn_format_tag (@prefix NVARCHAR(16), @id BIGINT) RETURNS NVARCHAR(64) AS
| BEGIN
| 	RETURN @prefix + N'-' + CAST(@id AS NVARCHAR(20));
| END;
: null
< row
: string
| app
: string
| fn_orders_of
: string
| CREATE FUNCTION app.fn_orders_of (@account BIGINT) RETURNS TABLE AS
| 	RETURN (SELECT id, total FROM app.orders WHERE account_id = @account);
: null
> query
|
| 	SELECT
| 		SCHEMA_NAME(ao.schema_id) AS schema_name,
| 		ao.name AS procedure_name,
| 		asm.definition
| 	FROM sys.all_objects ao
|         INNER JOIN sys.all_sql_modules asm ON asm.object_id = ao.object_id
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
< row
: string
| app
: string
| sp_account_orders
: string
| CREATE PROCEDURE app.sp_account_orders @account BIGINT AS
| BEGIN
| 	SELECT id, total FROM app.orders WHERE account_id = @account;
| 	SELECT COUNT(*) AS n FROM app.orders WHERE account_id = @account;
| END;
< row
: string
| app
: string
| sp_order_total
: string
| CREATE PROCEDURE app.sp_order_total @id BIGINT, @total INT OUTPUT AS
| BEGIN
| 	SELECT @total = total FROM app.orders WHERE id = @id;
| END;
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
< row
: int64 1205579333
: string
| trg_open_accounts_io
: string
| INSERT
: string
| INSTEAD OF
: string
| CREATE TRIGGER app.trg_open_accounts_io ON app.v_open_accounts INSTEAD OF INSERT AS
| BEGIN
| 	SET NOCOUNT ON;
| END;
: string
| v_open_accounts
: string
| app
: string
| V 
< row
: int64 1237579447
: string
| trg_orders_bi
: string
| INSERT
: string
| AFTER
: string
| CREATE TRIGGER app.trg_orders_bi ON app.orders AFTER INSERT AS
| BEGIN
| 	UPDATE app.orders SET total = 0 WHERE total < 0 AND id IN (SELECT id FROM inserted);
| END;
: string
| orders
: string
| app
: string
| U 
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
< row
: string
| app
: string
| syn_accounts
: null
: null
: string
| app
: string
| accounts
> query
| SELECT SCHEMA_NAME()
< result
< column "" "NVARCHAR"
< row
: string
| dbo
> commit
< ok
