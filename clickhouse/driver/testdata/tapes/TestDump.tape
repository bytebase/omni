-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDump default#1
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
> exec
|
|
| CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src; -- not restorable: TestDumpGaps
< ok 0 -
> query
|
| 		SELECT
| 			name,
| 			engine
| 		FROM system.tables
| 		WHERE database = ? AND name NOT LIKE '.%'
| 		ORDER BY name
: string
| omni_testdump
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
| mv_own
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
| SHOW CREATE TABLE omni_testdump.a_first
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.v_names
> query
| SHOW CREATE TABLE omni_testdump.a_top
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.a_top TO omni_testdump.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.mv_base
> query
| SHOW CREATE DICTIONARY omni_testdump.dict
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdump.dict
| (
|     `id` UInt64,
|     `name` String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE omni_testdump.mv_base
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.mv_base TO omni_testdump.target
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE omni_testdump.mv_own
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.mv_own
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE omni_testdump.mv_to
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.mv_to TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE omni_testdump.`odd-name`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.`odd-name`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE omni_testdump.src
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.src
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
| SHOW CREATE TABLE omni_testdump.target
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE omni_testdump.v_names
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.src
