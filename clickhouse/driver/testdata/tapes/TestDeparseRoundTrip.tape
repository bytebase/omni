-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDeparseRoundTrip default#1
> exec
|
| CREATE TABLE events (
|   id UInt64,
|   name String COMMENT 'the ''event'' name\\v1',
|   at DateTime64(3, 'UTC'),
|   amount Nullable(Decimal(12, 2)),
|   tags Array(String),
|   attrs Map(String, String),
|   level Enum8('info' = 1, 'warn' = 2) DEFAULT 'info',
|   seen Date DEFAULT '2026-01-02',
|   retries UInt32 DEFAULT 3,
|   INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| COMMENT 'events ''table''';
< ok 0 -
> exec
|
|
| CREATE TABLE `odd-name space` (
|   `select` Int32 DEFAULT 0,
|   `note col` Nullable(String)
| )
| ENGINE = Memory
| COMMENT 'line1\nline2';
< ok 0 -
> exec
|
|
| CREATE TABLE scratch (
|   k Int32
| )
| ENGINE = ReplacingMergeTree
| ORDER BY tuple();
< ok 0 -
> exec
|
|
| CREATE VIEW event_names AS SELECT id, name FROM events COMMENT 'active ''view''';
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
| omni_testdeparseroundtrip
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
| event_names
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
| event_names
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
| events
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
| events
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
| the 'event' name\v1
< row
: string
| events
: string
| at
: uint64 3
: string
|
: string
| 0
: string
| DateTime64(3, 'UTC')
: string
|
: string
|
: string
|
< row
: string
| events
: string
| amount
: uint64 4
: string
|
: string
| 1
: string
| Nullable(Decimal(12, 2))
: string
|
: string
|
: string
|
< row
: string
| events
: string
| tags
: uint64 5
: string
|
: string
| 0
: string
| Array(String)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| attrs
: uint64 6
: string
|
: string
| 0
: string
| Map(String, String)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| level
: uint64 7
: string
| 'info'
: string
| 0
: string
| Enum8('info' = 1, 'warn' = 2)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| seen
: uint64 8
: string
| '2026-01-02'
: string
| 0
: string
| Date
: string
|
: string
|
: string
|
< row
: string
| events
: string
| retries
: uint64 9
: string
| 3
: string
| 0
: string
| UInt32
: string
|
: string
|
: string
|
< row
: string
| odd-name space
: string
| select
: uint64 1
: string
| 0
: string
| 0
: string
| Int32
: string
|
: string
|
: string
|
< row
: string
| odd-name space
: string
| note col
: uint64 2
: string
|
: string
| 1
: string
| Nullable(String)
: string
|
: string
|
: string
|
< row
: string
| scratch
: string
| k
: uint64 1
: string
|
: string
| 0
: string
| Int32
: string
|
: string
|
: string
|
> query
|
| 		SELECT table, name, default_kind, default_expression
| 		FROM system.columns
| 		WHERE database = $1 AND default_kind IN ('ALIAS', 'MATERIALIZED')
| 		ORDER BY table, name
: string
| omni_testdeparseroundtrip
< result
< column "table" "String" "string"
< column "name" "String" "string"
< column "default_kind" "String" "string"
< column "default_expression" "String" "string"
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
| omni_testdeparseroundtrip
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
| event_names
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdeparseroundtrip.event_names (`id` UInt64, `name` String) COMMENT 'active \'view\'' AS SELECT id, name FROM omni_testdeparseroundtrip.events
: string
| active 'view'
: string
|
: string
|
< row
: string
| events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.events (`id` UInt64, `name` String COMMENT 'the \'event\' name\\v1', `at` DateTime64(3, 'UTC'), `amount` Nullable(Decimal(12, 2)), `tags` Array(String), `attrs` Map(String, String), `level` Enum8('info' = 1, 'warn' = 2) DEFAULT 'info', `seen` Date DEFAULT '2026-01-02', `retries` UInt32 DEFAULT 3, INDEX name_bloom name TYPE bloom_filter GRANULARITY 4) ENGINE = MergeTree PRIMARY KEY id ORDER BY (id, name) SETTINGS index_granularity = 8192 COMMENT 'events \'table\''
: string
| events 'table'
: string
| id, name
: string
| id
< row
: string
| odd-name space
: string
| Memory
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.`odd-name space` (`select` Int32 DEFAULT 0, `note col` Nullable(String)) ENGINE = Memory COMMENT 'line1\nline2'
: string
| line1
| line2
: string
|
: string
|
< row
: string
| scratch
: string
| ReplacingMergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.scratch (`k` Int32) ENGINE = ReplacingMergeTree ORDER BY tuple() SETTINGS index_granularity = 8192
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
| omni_testdeparseroundtrip
: string
| events
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
| omni_testdeparseroundtrip
: string
| odd-name space
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdeparseroundtrip
: string
| scratch
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query
|
| 		SELECT
| 			short_name,
| 			table
| 		FROM system.row_policies
| 		WHERE database = ?
| 		ORDER BY table, short_name
: string
| omni_testdeparseroundtrip
< result
< column "short_name" "String" "string"
< column "table" "String" "string"
> exec
| DROP DATABASE IF EXISTS `omni_testdeparseroundtrip`
< ok 0 -
> exec
| CREATE DATABASE `omni_testdeparseroundtrip`
< ok 0 -
> begin
< ok
> exec
| --
| -- Table structure for `events`
| --
| CREATE TABLE `events` (
|   `id` UInt64,
|   `name` String COMMENT 'the \'event\' name\\v1',
|   `at` DateTime64(3, 'UTC'),
|   `amount` Nullable(Decimal(12, 2)) DEFAULT NULL,
|   `tags` Array(String),
|   `attrs` Map(String, String),
|   `level` Enum8('info' = 1, 'warn' = 2) DEFAULT 'info',
|   `seen` Date DEFAULT '2026-01-02',
|   `retries` UInt32 DEFAULT 3,
|   INDEX `name_bloom` name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| COMMENT 'events \'table\'';
< ok 0 -
> exec
|
|
| --
| -- Table structure for `odd-name space`
| --
| CREATE TABLE `odd-name space` (
|   `select` Int32 DEFAULT 0,
|   `note col` Nullable(String) DEFAULT NULL
| )
| ENGINE = Memory
| COMMENT 'line1\nline2';
< ok 0 -
> exec
|
|
| --
| -- Table structure for `scratch`
| --
| CREATE TABLE `scratch` (
|   `k` Int32
| )
| ENGINE = ReplacingMergeTree
| ORDER BY tuple();
< ok 0 -
> exec
|
|
| --
| -- View structure for `event_names`
| --
| CREATE VIEW omni_testdeparseroundtrip.event_names (`id` UInt64, `name` String) COMMENT 'active \'view\'' AS SELECT id, name FROM omni_testdeparseroundtrip.events;
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
| omni_testdeparseroundtrip
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
| event_names
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
| event_names
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
| events
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
| events
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
| the 'event' name\v1
< row
: string
| events
: string
| at
: uint64 3
: string
|
: string
| 0
: string
| DateTime64(3, 'UTC')
: string
|
: string
|
: string
|
< row
: string
| events
: string
| amount
: uint64 4
: string
| NULL
: string
| 1
: string
| Nullable(Decimal(12, 2))
: string
|
: string
|
: string
|
< row
: string
| events
: string
| tags
: uint64 5
: string
|
: string
| 0
: string
| Array(String)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| attrs
: uint64 6
: string
|
: string
| 0
: string
| Map(String, String)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| level
: uint64 7
: string
| 'info'
: string
| 0
: string
| Enum8('info' = 1, 'warn' = 2)
: string
|
: string
|
: string
|
< row
: string
| events
: string
| seen
: uint64 8
: string
| '2026-01-02'
: string
| 0
: string
| Date
: string
|
: string
|
: string
|
< row
: string
| events
: string
| retries
: uint64 9
: string
| 3
: string
| 0
: string
| UInt32
: string
|
: string
|
: string
|
< row
: string
| odd-name space
: string
| select
: uint64 1
: string
| 0
: string
| 0
: string
| Int32
: string
|
: string
|
: string
|
< row
: string
| odd-name space
: string
| note col
: uint64 2
: string
| NULL
: string
| 1
: string
| Nullable(String)
: string
|
: string
|
: string
|
< row
: string
| scratch
: string
| k
: uint64 1
: string
|
: string
| 0
: string
| Int32
: string
|
: string
|
: string
|
> query
|
| 		SELECT table, name, default_kind, default_expression
| 		FROM system.columns
| 		WHERE database = $1 AND default_kind IN ('ALIAS', 'MATERIALIZED')
| 		ORDER BY table, name
: string
| omni_testdeparseroundtrip
< result
< column "table" "String" "string"
< column "name" "String" "string"
< column "default_kind" "String" "string"
< column "default_expression" "String" "string"
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
| omni_testdeparseroundtrip
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
| event_names
: string
| View
: uint64 0
: uint64 0
: string
| CREATE VIEW omni_testdeparseroundtrip.event_names (`id` UInt64, `name` String) COMMENT 'active \'view\'' AS SELECT id, name FROM omni_testdeparseroundtrip.events
: string
| active 'view'
: string
|
: string
|
< row
: string
| events
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.events (`id` UInt64, `name` String COMMENT 'the \'event\' name\\v1', `at` DateTime64(3, 'UTC'), `amount` Nullable(Decimal(12, 2)) DEFAULT NULL, `tags` Array(String), `attrs` Map(String, String), `level` Enum8('info' = 1, 'warn' = 2) DEFAULT 'info', `seen` Date DEFAULT '2026-01-02', `retries` UInt32 DEFAULT 3, INDEX name_bloom name TYPE bloom_filter GRANULARITY 4) ENGINE = MergeTree PRIMARY KEY id ORDER BY (id, name) SETTINGS index_granularity = 8192 COMMENT 'events \'table\''
: string
| events 'table'
: string
| id, name
: string
| id
< row
: string
| odd-name space
: string
| Memory
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.`odd-name space` (`select` Int32 DEFAULT 0, `note col` Nullable(String) DEFAULT NULL) ENGINE = Memory COMMENT 'line1\nline2'
: string
| line1
| line2
: string
|
: string
|
< row
: string
| scratch
: string
| ReplacingMergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testdeparseroundtrip.scratch (`k` Int32) ENGINE = ReplacingMergeTree ORDER BY tuple() SETTINGS index_granularity = 8192
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
| omni_testdeparseroundtrip
: string
| events
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
| omni_testdeparseroundtrip
: string
| odd-name space
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testdeparseroundtrip
: string
| scratch
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
