-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDumpRestores default#1
> exec
| -- The schema the dump tests read: a MergeTree table with an index, a primary
| -- key, a setting and a comment; a table the materialized views write to; a name
| -- that needs backquotes; a user table whose name starts with the dot ClickHouse
| -- uses for a materialized view's storage; a name that holds a line break; a
| -- dictionary, and a table that reads it; a view, and a view over it under a name
| -- that sorts first; two materialized views, the first holding the other's name
| -- in a literal; and a materialized view with its own engine, whose inner table
| -- only the server names.
| --
| -- omni_testdumprestores stands for the database the test created, which the dictionary's
| -- source and the materialized view's literal name.
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
| CREATE TABLE `odd-name` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE TABLE `.events` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| -- A name the internal-table pattern must not take for one: its underscores are
| -- what a LIKE pattern would read as wildcards.
| CREATE TABLE `.innerXid.events` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE TABLE `line
| break` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
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
> exec
|
|
| -- A dictionary whose source is a view, and a table that loads it: creating the
| -- table queries the view, so both follow the views.
| CREATE VIEW v_source AS SELECT id FROM src;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY view_dict (
|   id UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'v_source'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| CREATE TABLE view_dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.view_dict');
< ok 0 -
> exec
|
|
| CREATE TABLE dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.dict');
< ok 0 -
> exec
|
|
| -- A dictionary whose source is a materialized view, and a table that loads it
| -- under a name that sorts first: creating the table queries the view, so it
| -- follows the view.
| -- Its table alias is a reference with the spaces around it, in backquotes: a name
| -- inside backquotes is part of that name, not a reference to a_mv_dict_table,
| -- which loads the dictionary this view sources.
| CREATE MATERIALIZED VIEW z_src_mv TO target AS SELECT ` omni_testdumprestores.a_mv_dict_table `.id, 'm' AS name FROM src AS ` omni_testdumprestores.a_mv_dict_table `;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY mv_dict (
|   id UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| CREATE TABLE a_mv_dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.mv_dict');
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
| CREATE MATERIALIZED VIEW mv_base TO target AS SELECT id, 'omni_testdumprestores.a_top' AS name FROM src;
< ok 0 -
> exec
|
|
|
| -- An ordinary view reads mv_base, and a materialized view reads that view: a
| -- materialized view resolves its source when it is created, through the view, so
| -- the dump has to write mv_base first whatever the names say. The literal above
| -- names a_top and is not a read.
| CREATE VIEW v_bridge AS SELECT id, name FROM mv_base;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW a_top TO target AS SELECT id, name FROM v_bridge;
< ok 0 -
> exec
|
|
| -- A materialized view whose name holds a line break, and one that reads it under
| -- a name that sorts first: the reference the dump looks for is the rendering
| -- ClickHouse writes, with the line break escaped.
| CREATE MATERIALIZED VIEW `nl
| mv` TO target AS SELECT id, 'x' AS name FROM src;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW a_after TO target AS SELECT id, 'y' AS name FROM `nl
| mv`;
< ok 0 -
> exec
|
|
| -- A materialized view reading a table on the Dictionary engine, under a name that
| -- sorts first: creating it needs the dictionary table, which loads the
| -- dictionary.
| CREATE MATERIALIZED VIEW a_dict TO target AS SELECT id, 'z' AS name FROM dict_table;
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
| .events
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
| .innerXid.events
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
| a_after
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
| a_after
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
| a_dict
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
| a_dict
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
| a_mv_dict_table
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
| a_top
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
| dict_table
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
| line
| break
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
| mv_base
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
| mv_dict
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
| nl
| mv
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
| nl
| mv
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
| v_bridge
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
| v_bridge
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
< row
: string
| v_source
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
| view_dict
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
| view_dict_table
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
| z_src_mv
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
| z_src_mv
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
| .events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`.events` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| .innerXid.events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`.innerXid.events` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| a_after
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_after TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'y' AS name FROM omni_testdumprestores.`nl\nmv`
: string
|
: string
|
: string
|
< row
: string
| a_dict
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_dict TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'z' AS name FROM omni_testdumprestores.dict_table
: string
|
: string
|
: string
|
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
| a_mv_dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.a_mv_dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.mv_dict')
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
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.v_bridge
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
| dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.dict')
: string
|
: string
|
: string
|
< row
: string
| line
| break
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`line\nbreak` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| mv_base
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'omni_testdumprestores.a_top' AS name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| mv_dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.mv_dict (`id` UInt64) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'z_src_mv')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| nl
| mv
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.`nl\nmv` TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'x' AS name FROM omni_testdumprestores.src
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
| v_bridge
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_bridge (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.mv_base
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
< row
: string
| v_source
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_source (`id` UInt64) AS SELECT id FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| view_dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.view_dict (`id` UInt64) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'v_source')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| view_dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.view_dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.view_dict')
: string
|
: string
|
: string
|
< row
: string
| z_src_mv
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.z_src_mv TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT ` omni_testdumprestores.a_mv_dict_table `.id, 'm' AS name FROM omni_testdumprestores.src AS ` omni_testdumprestores.a_mv_dict_table `
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
| .events
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| .innerXid.events
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_after
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_mv_dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
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
| dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| line
| break
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
| mv_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| nl
| mv
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
> query ^
: string
| omni_testdumprestores
: string
| view_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| view_dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| z_src_mv
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> exec
| SET format_display_secrets_in_show_and_select = 1
< ok 0 -
> query
|
| 		SELECT
| 			name
| 		FROM system.tables
| 		WHERE database = ? AND NOT match(name, '^[.]inner_id[.][0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$')
| 		ORDER BY name
: string
| omni_testdumprestores
< result
< column "name" "String" "string"
< row
: string
| .events
< row
: string
| .innerXid.events
< row
: string
| a_after
< row
: string
| a_dict
< row
: string
| a_first
< row
: string
| a_mv_dict_table
< row
: string
| a_top
< row
: string
| dict
< row
: string
| dict_table
< row
: string
| line
| break
< row
: string
| mv_base
< row
: string
| mv_dict
< row
: string
| nl
| mv
< row
: string
| odd-name
< row
: string
| src
< row
: string
| target
< row
: string
| v_bridge
< row
: string
| v_names
< row
: string
| v_source
< row
: string
| view_dict
< row
: string
| view_dict_table
< row
: string
| z_src_mv
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`.events`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.`.events`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`.innerXid.events`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.`.innerXid.events`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`a_after`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_after TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'y' AS name
| FROM omni_testdumprestores.`nl\nmv`
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`a_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_dict TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'z' AS name
| FROM omni_testdumprestores.dict_table
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`a_first`
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
| SHOW CREATE TABLE `omni_testdumprestores`.`a_mv_dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.a_mv_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.mv_dict')
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`a_top`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.v_bridge
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`dict`
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
| SHOW CREATE TABLE `omni_testdumprestores`.`dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.dict')
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`line\nbreak`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.`line\nbreak`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`mv_base`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'omni_testdumprestores.a_top' AS name
| FROM omni_testdumprestores.src
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`mv_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdumprestores.mv_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`nl\nmv`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.`nl\nmv` TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'x' AS name
| FROM omni_testdumprestores.src
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`odd-name`
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
| SHOW CREATE TABLE `omni_testdumprestores`.`src`
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
| SHOW CREATE TABLE `omni_testdumprestores`.`target`
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
| SHOW CREATE TABLE `omni_testdumprestores`.`v_bridge`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdumprestores.v_bridge
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.mv_base
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`v_names`
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
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`v_source`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdumprestores.v_source
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.src
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`view_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdumprestores.view_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'v_source'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`view_dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprestores.view_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.view_dict')
> query
| SHOW CREATE TABLE `omni_testdumprestores`.`z_src_mv`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.z_src_mv TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     ` omni_testdumprestores.a_mv_dict_table `.id,
|     'm' AS name
| FROM omni_testdumprestores.src AS ` omni_testdumprestores.a_mv_dict_table `
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.a_after TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'y' AS name
| FROM omni_testdumprestores.`nl\nmv`
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores a_after (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier a_after
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'y' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.nl
< row
: string
| mv
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.a_dict TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'z' AS name
| FROM omni_testdumprestores.dict_table
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores a_dict (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier a_dict
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'z' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.dict_table
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE VIEW omni_testdumprestores.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.v_names
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores a_first (children 4)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier a_first
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 1)
< row
: string
|      Identifier id
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.v_names
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.v_bridge
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores a_top (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier a_top
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.v_bridge
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'omni_testdumprestores.a_top' AS name
| FROM omni_testdumprestores.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores mv_base (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier mv_base
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'omni_testdumprestores.a_top' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.src
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.`nl\nmv` TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'x' AS name
| FROM omni_testdumprestores.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores nl
< row
: string
| mv (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier nl
< row
: string
| mv
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'x' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.src
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE VIEW omni_testdumprestores.v_bridge
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.mv_base
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores v_bridge (children 4)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier v_bridge
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.mv_base
> query
| EXPLAIN AST CREATE VIEW omni_testdumprestores.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores v_names (children 4)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier v_names
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.src
> query
| EXPLAIN AST CREATE VIEW omni_testdumprestores.v_source
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdumprestores.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores v_source (children 4)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier v_source
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 1)
< row
: string
|      Identifier id
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.src
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdumprestores.z_src_mv TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     ` omni_testdumprestores.a_mv_dict_table `.id,
|     'm' AS name
| FROM omni_testdumprestores.src AS ` omni_testdumprestores.a_mv_dict_table `
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdumprestores z_src_mv (children 5)
< row
: string
|  Identifier omni_testdumprestores
< row
: string
|  Identifier z_src_mv
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier  omni_testdumprestores.a_mv_dict_table .id
< row
: string
|      Literal 'm' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdumprestores.src (alias  omni_testdumprestores.a_mv_dict_table )
< row
: string
|  ViewTargets
> exec
| DROP DATABASE IF EXISTS `omni_testdumprestores`
< ok 0 -
> exec
| CREATE DATABASE `omni_testdumprestores`
< ok 0 -
> begin
< ok
> exec
| --
| -- Table structure for `.events`
| --
| CREATE TABLE omni_testdumprestores.`.events`
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
| -- Table structure for `.innerXid.events`
| --
| CREATE TABLE omni_testdumprestores.`.innerXid.events`
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
| -- Table structure for `line\nbreak`
| --
| CREATE TABLE omni_testdumprestores.`line\nbreak`
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
| -- Table structure for `src`
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
| -- Table structure for `target`
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
| -- Dictionary structure for `dict`
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
| -- Dictionary structure for `mv_dict`
| --
| CREATE DICTIONARY omni_testdumprestores.mv_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT());
< ok 0 -
> exec
|
|
| --
| -- Dictionary structure for `view_dict`
| --
| CREATE DICTIONARY omni_testdumprestores.view_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'v_source'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT());
< ok 0 -
> exec
|
|
| --
| -- View structure for `a_first`
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
| -- View structure for `v_bridge`
| --
| CREATE VIEW omni_testdumprestores.v_bridge
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.mv_base;
< ok 0 -
> exec
|
|
| --
| -- View structure for `v_names`
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
| -- View structure for `v_source`
| --
| CREATE VIEW omni_testdumprestores.v_source
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
| -- Table structure for `dict_table`
| --
| CREATE TABLE omni_testdumprestores.dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.dict');
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `mv_base`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'omni_testdumprestores.a_top' AS name
| FROM omni_testdumprestores.src;
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `nl\nmv`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.`nl\nmv` TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'x' AS name
| FROM omni_testdumprestores.src;
< ok 0 -
> exec
|
|
| --
| -- Table structure for `view_dict_table`
| --
| CREATE TABLE omni_testdumprestores.view_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.view_dict');
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `z_src_mv`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.z_src_mv TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     ` omni_testdumprestores.a_mv_dict_table `.id,
|     'm' AS name
| FROM omni_testdumprestores.src AS ` omni_testdumprestores.a_mv_dict_table `;
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `a_after`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_after TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'y' AS name
| FROM omni_testdumprestores.`nl\nmv`;
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `a_dict`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_dict TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'z' AS name
| FROM omni_testdumprestores.dict_table;
< ok 0 -
> exec
|
|
| --
| -- Table structure for `a_mv_dict_table`
| --
| CREATE TABLE omni_testdumprestores.a_mv_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdumprestores.mv_dict');
< ok 0 -
> exec
|
|
| --
| -- Materialized view structure for `a_top`
| --
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdumprestores.v_bridge;
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
| .events
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
| .innerXid.events
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
| a_after
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
| a_after
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
| a_dict
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
| a_dict
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
| a_mv_dict_table
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
| a_top
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
| dict_table
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
| line
| break
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
| mv_base
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
| mv_dict
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
| nl
| mv
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
| nl
| mv
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
| v_bridge
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
| v_bridge
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
< row
: string
| v_source
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
| view_dict
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
| view_dict_table
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
| z_src_mv
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
| z_src_mv
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
| .events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`.events` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| .innerXid.events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`.innerXid.events` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| a_after
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_after TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'y' AS name FROM omni_testdumprestores.`nl\nmv`
: string
|
: string
|
: string
|
< row
: string
| a_dict
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_dict TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'z' AS name FROM omni_testdumprestores.dict_table
: string
|
: string
|
: string
|
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
| a_mv_dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.a_mv_dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.mv_dict')
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
| CREATE MATERIALIZED VIEW omni_testdumprestores.a_top TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.v_bridge
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
| dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.dict')
: string
|
: string
|
: string
|
< row
: string
| line
| break
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.`line\nbreak` (`id` UInt64) ENGINE = MergeTree ORDER BY id SETTINGS index_granularity = 8192
: string
|
: string
| id
: string
| id
< row
: string
| mv_base
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.mv_base TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'omni_testdumprestores.a_top' AS name FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| mv_dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.mv_dict (`id` UInt64) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'z_src_mv')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| nl
| mv
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.`nl\nmv` TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT id, 'x' AS name FROM omni_testdumprestores.src
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
| v_bridge
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_bridge (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdumprestores.mv_base
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
< row
: string
| v_source
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdumprestores.v_source (`id` UInt64) AS SELECT id FROM omni_testdumprestores.src
: string
|
: string
|
: string
|
< row
: string
| view_dict
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE DICTIONARY omni_testdumprestores.view_dict (`id` UInt64) PRIMARY KEY id SOURCE(CLICKHOUSE(TABLE 'v_source')) LIFETIME(MIN 0 MAX 300) LAYOUT(FLAT())
: string
|
: string
|
: string
|
< row
: string
| view_dict_table
: string
| Dictionary
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdumprestores.view_dict_table (`id` UInt64) ENGINE = Dictionary('omni_testdumprestores.view_dict')
: string
|
: string
|
: string
|
< row
: string
| z_src_mv
: string
| MaterializedView
: uint64 0
: uint64 0
: string
| CREATE MATERIALIZED VIEW omni_testdumprestores.z_src_mv TO omni_testdumprestores.target (`id` UInt64, `name` String) AS SELECT ` omni_testdumprestores.a_mv_dict_table `.id, 'm' AS name FROM omni_testdumprestores.src AS ` omni_testdumprestores.a_mv_dict_table `
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
| .events
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| .innerXid.events
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_after
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| a_mv_dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
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
| dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| line
| break
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
| mv_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| nl
| mv
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
> query ^
: string
| omni_testdumprestores
: string
| view_dict
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| view_dict_table
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdumprestores
: string
| z_src_mv
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
