-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDumpRestores default#1
> exec
| -- The schema the dump tests read: a MergeTree table with an index, a primary
| -- key, a setting and a comment; a table the materialized views write to; a view,
| -- and a view over it under a name that sorts first; a materialized view, and one
| -- over that under a name that sorts first; a dictionary; and a materialized view
| -- with its own engine, whose inner table only the server names.
| CREATE TABLE src (
|   id UInt64,
|   name String,
|   INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| SETTINGS index_granularity = 4096
| COMMENT 'the source';
< ok 0 -
> exec
|
|
| CREATE TABLE target (
|   id UInt64,
|   name String
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| -- A name ClickHouse writes in backquotes, so the dump quotes it back.
| CREATE TABLE `odd-name` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE VIEW v_names AS SELECT id, name FROM src;
< ok 0 -
> exec
|
|
| -- A view states its columns, which ClickHouse creates without resolving its
| -- query, so this one replays before the view it reads.
| CREATE VIEW a_first AS SELECT id FROM v_names;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW mv_to TO target AS SELECT id, name FROM src;
< ok 0 -
> exec
|
|
| -- A materialized view, and one over it under a name that sorts first: a
| -- materialized view resolves its source when it is created, so the dump has to
| -- write mv_base before a_top whatever their names say.
| CREATE MATERIALIZED VIEW mv_base TO target AS SELECT id FROM src;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW a_top TO target AS SELECT id FROM mv_base;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY dict (
|   id UInt64,
|   name String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> query
|
| 		SELECT
| 			table_name,
| 			column_name,
| 			ordinal_position,
| 			column_default,
| 			is_nullable,
| 			column_type,
| 			ifNull(character_set_name, ''),
| 			ifNull(collation_name, ''),
| 			column_comment
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| omni_testdumprestores
< result
< column "table_name" "String" "string"
< column "column_name" "String" "string"
< column "ordinal_position" "UInt64" "uint64"
< column "column_default" "String" "string"
< column "is_nullable" "String" "string"
< column "column_type" "String" "string"
< column "ifNull(character_set_name, '')" "String" "string"
< column "ifNull(collation_name, '')" "String" "string"
< column "column_comment" "String" "string"
< row
: string
| a_first
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| a_top
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| mv_base
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| odd-name
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| src
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| src
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| target
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| target
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| v_names
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| v_names
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
> query
|
| 		SELECT
| 			name,
| 			engine,
| 			ifNull(total_rows, 0),
| 			ifNull(total_bytes, 0),
| 			create_table_query,
| 			comment,
| 			sorting_key,
| 			primary_key
| 		FROM system.tables
| 		WHERE database = $1
| 		ORDER BY name
: string
| omni_testdumprestores
< result
< column "name" "String" "string"
< column "engine" "String" "string"
< column "ifNull(total_rows, 0)" "UInt64" "uint64"
< column "ifNull(total_bytes, 0)" "UInt64" "uint64"
< column "create_table_query" "String" "string"
< column "comment" "String" "string"
< column "sorting_key" "String" "string"
< column "primary_key" "String" "string"
< row
: string
| a_first
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.a_first (`id` UInt64) AS SELECT id FROM omni_testdumprestores.v_names
: string
|
: string
|
: string
|
< row
: string
| a_top
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target (`id` UInt64) AS SELECT id FROM omni_testdumprestores.mv_base
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.dict (`id` UInt64, `name` String) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'src')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| mv_base
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target (`id` UInt64) AS SELECT id FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_to TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| odd-name
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`odd-name` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| src
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.src (`id` UInt64, `name` String, INDEX name_bloom name TYPE bloom_filter GRANULARITY 4) ENGINE = MergeTree PRIMARY KEY id ORDER BY (id, name) SETTINGS index_granularity = 4096 COMMENT 'the source'
: string
| the source
: string
| id, name
: string
| id
< row
: string
| target
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.target (`id` UInt64, `name` String) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| v_names
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_names (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
> query
|
| 		SELECT
| 			name,
| 			type,
| 			expr,
| 			granularity
| 		FROM system.data_skipping_indices
| 		WHERE database = $1 AND table = $2
| 		ORDER BY name
: string
| omni_testdumprestores
: string
| a_top
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| mv_base
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| mv_to
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| odd-name
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| src
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
< row
: string
| name_bloom
: string
| bloom_filter
: string
| name
: uint64 4
> query ^
: string
| omni_testdumprestores
: string
| target
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query
|
| 		SELECT
| 			name,
| 			engine
| 		FROM system.tables
| 		WHERE database = ? AND name NOT LIKE '.%'
| 		ORDER BY name
: string
| omni_testdumprestores
< result
< column "name" "String" "string"
< column "engine" "String" "string"
< row
: string
| a_first
: string
| View
< row
: string
| a_top
: string
| MaterializedView
< row
: string
| dict
: string
| Dictionary
< row
: string
| mv_base
: string
| MaterializedView
< row
: string
| mv_to
: string
| MaterializedView
< row
: string
| odd-name
: string
| MergeTree
< row
: string
| src
: string
| MergeTree
< row
: string
| target
: string
| MergeTree
< row
: string
| v_names
: string
| View
> query
| SHOW CREATE TABLE omni_testdumprestores.a_first
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdumprestores.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.v_names
> query
| SHOW CREATE TABLE omni_testdumprestores.a_top
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.mv_base
> query
| SHOW CREATE DICTIONARY omni_testdumprestores.dict
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdumprestores.dict
| (
|     `id` UInt64,
|     `name` String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE omni_testdumprestores.mv_base
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.src
> query
| SHOW CREATE TABLE omni_testdumprestores.mv_to
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_to TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.src
> query
| SHOW CREATE TABLE omni_testdumprestores.`odd-name`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.`odd-name`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE omni_testdumprestores.src
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.src
| (
|     `id` UInt64,
|     `name` String,
|     INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| PRIMARY KEY id
| ORDER BY (id, name)
| SETTINGS index_granularity = 4096
| COMMENT 'the source'
> query
| SHOW CREATE TABLE omni_testdumprestores.target
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE omni_testdumprestores.v_names
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdumprestores.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.src
> exec
| DROP DATABASE IF EXISTS omni_testdumprestores
< ok 0 -
> exec
| CREATE DATABASE omni_testdumprestores
< ok 0 -
> begin
< ok
> exec
| --
| -- Table structure for `odd-name`
| --
| CREATE TABLE omni_testdumprestores.`odd-name`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192;
< ok 0 -
> exec
|
|
| --
| -- Table structure for src
| --
| CREATE TABLE omni_testdumprestores.src
| (
|     `id` UInt64,
|     `name` String,
|     INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| PRIMARY KEY id
| ORDER BY (id, name)
| SETTINGS index_granularity = 4096
| COMMENT 'the source';
< ok 0 -
> exec
|
|
| --
| -- Table structure for target
| --
| CREATE TABLE omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192;
< ok 0 -
> exec
|
|
| --
| -- Dictionary structure for dict
| --
| CREATE DICTIONARY omni_testdumprestores.dict
| (
|     `id` UInt64,
|     `name` String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT());
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for mv_base
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.src;
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for mv_to
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_to TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.src;
< ok 0 -
> exec
|
|
| --
| -- View structure for v_names
| --
| CREATE VIEW omni_testdumprestores.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.src;
< ok 0 -
> exec
|
|
| --
| -- View structure for a_first
| --
| CREATE VIEW omni_testdumprestores.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.v_names;
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for a_top
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.mv_base;
< ok 0 -
> commit
< ok
> query
|
| 		SELECT
| 			table_name,
| 			column_name,
| 			ordinal_position,
| 			column_default,
| 			is_nullable,
| 			column_type,
| 			ifNull(character_set_name, ''),
| 			ifNull(collation_name, ''),
| 			column_comment
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| omni_testdumprestores
< result
< column "table_name" "String" "string"
< column "column_name" "String" "string"
< column "ordinal_position" "UInt64" "uint64"
< column "column_default" "String" "string"
< column "is_nullable" "String" "string"
< column "column_type" "String" "string"
< column "ifNull(character_set_name, '')" "String" "string"
< column "ifNull(collation_name, '')" "String" "string"
< column "column_comment" "String" "string"
< row
: string
| a_first
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| a_top
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| mv_base
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| odd-name
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| src
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| src
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| target
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| target
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
< row
: string
| v_names
: string
| id
: uint64 1
: string
|
: string
| 0
: string
| UInt64
: string
|
: string
|
: string
|
< row
: string
| v_names
: string
| name
: uint64 2
: string
|
: string
| 0
: string
| String
: string
|
: string
|
: string
|
> query
|
| 		SELECT
| 			name,
| 			engine,
| 			ifNull(total_rows, 0),
| 			ifNull(total_bytes, 0),
| 			create_table_query,
| 			comment,
| 			sorting_key,
| 			primary_key
| 		FROM system.tables
| 		WHERE database = $1
| 		ORDER BY name
: string
| omni_testdumprestores
< result
< column "name" "String" "string"
< column "engine" "String" "string"
< column "ifNull(total_rows, 0)" "UInt64" "uint64"
< column "ifNull(total_bytes, 0)" "UInt64" "uint64"
< column "create_table_query" "String" "string"
< column "comment" "String" "string"
< column "sorting_key" "String" "string"
< column "primary_key" "String" "string"
< row
: string
| a_first
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.a_first (`id` UInt64) AS SELECT id FROM omni_testdumprestores.v_names
: string
|
: string
|
: string
|
< row
: string
| a_top
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target (`id` UInt64) AS SELECT id FROM omni_testdumprestores.mv_base
: string
|
: string
|
: string
|
< row
: string
| dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.dict (`id` UInt64, `name` String) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'src')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| mv_base
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target (`id` UInt64) AS SELECT id FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| mv_to
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_to TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| odd-name
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`odd-name` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| src
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.src (`id` UInt64, `name` String, INDEX name_bloom name TYPE bloom_filter GRANULARITY 4) ENGINE = MergeTree PRIMARY KEY id ORDER BY (id, name) SETTINGS index_granularity = 4096 COMMENT 'the source'
: string
| the source
: string
| id, name
: string
| id
< row
: string
| target
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.target (`id` UInt64, `name` String) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| v_names
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_names (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
> query
|
| 		SELECT
| 			name,
| 			type,
| 			expr,
| 			granularity
| 		FROM system.data_skipping_indices
| 		WHERE database = $1 AND table = $2
| 		ORDER BY name
: string
| omni_testdumprestores
: string
| a_top
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| mv_base
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| mv_to
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| odd-name
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| src
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
< row
: string
| name_bloom
: string
| bloom_filter
: string
| name
: uint64 4
> query ^
: string
| omni_testdumprestores
: string
| target
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
