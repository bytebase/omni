-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestSyncDatabaseDefaultDatabase default#1
> exec
| -- The schema the sync tests read. Nothing is inserted: a table's row count and
| -- data size are the server's approximations, and an empty table makes both
| -- zero in every recording.
| CREATE TABLE events (
|   id UInt64,
|   name String COMMENT 'the event name',
|   at DateTime64(3, 'UTC'),
|   amount Nullable(Decimal(12, 2)),
|   tags Array(String),
|   attrs Map(String, String),
|   fixed FixedString(4),
|   flag Bool,
|   level Enum8('info' = 1, 'warn' = 2),
|   seen Date DEFAULT '2026-01-02'
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| COMMENT 'events the sync reads';
< ok 0 -
> exec
|
|
| CREATE TABLE skipped (
|   url String,
|   body String,
|   INDEX url_bloom url TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY url;
< ok 0 -
> exec
|
|
| CREATE VIEW event_names AS SELECT id, name FROM events;
< ok 0 -
> query
| SELECT currentDatabase()
< result
< column "currentDatabase()" "String" "string"
< row
: string
| omni_testsyncdatabasedefaultdatabase
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
| omni_testsyncdatabasedefaultdatabase
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
| the event name
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
| fixed
: uint64 7
: string
|
: string
| 0
: string
| FixedString(4)
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
| flag
: uint64 8
: string
|
: string
| 0
: string
| Bool
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
: uint64 9
: string
|
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
: uint64 10
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
| skipped
: string
| url
: uint64 1
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
| skipped
: string
| body
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
| omni_testsyncdatabasedefaultdatabase
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
| CREATE VIEW omni_testsyncdatabasedefaultdatabase.event_names (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testsyncdatabasedefaultdatabase.events
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
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testsyncdatabasedefaultdatabase.events (`id` UInt64, `name` String COMMENT 'the event name', `at` DateTime64(3, 'UTC'), `amount` Nullable(Decimal(12, 2)), `tags` Array(String), `attrs` Map(String, String), `fixed` FixedString(4), `flag` Bool, `level` Enum8('info' = 1, 'warn' = 2), `seen` Date DEFAULT '2026-01-02') ENGINE = MergeTree PRIMARY KEY id ORDER BY (id, name) SETTINGS index_granularity = 8192 COMMENT 'events the sync reads'
: string
| events the sync reads
: string
| id, name
: string
| id
< row
: string
| skipped
: string
| MergeTree
: uint64 0
: uint64 0
: string
| CREATE TABLE omni_testsyncdatabasedefaultdatabase.skipped (`url` String, `body` String, INDEX url_bloom url TYPE bloom_filter GRANULARITY 4) ENGINE = MergeTree ORDER BY url SETTINGS index_granularity = 8192
: string
|
: string
| url
: string
| url
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
| omni_testsyncdatabasedefaultdatabase
: string
| events
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
> query ^
: string
| omni_testsyncdatabasedefaultdatabase
: string
| skipped
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "expr" "String" "string"
< column "granularity" "UInt64" "uint64"
< row
: string
| url_bloom
: string
| bloom_filter
: string
| url
: uint64 4
