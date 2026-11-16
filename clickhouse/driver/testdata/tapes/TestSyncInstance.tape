-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestSyncInstance default#1
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
| SELECT VERSION()
< result
< column "VERSION()" "String" "string"
< row
: string
| 26.8.38.9
> query
|
| 		SELECT
| 			name
| 		FROM system.users
| 		ORDER BY name
| 	
< result
< column "name" "String" "string"
< row
: string
| default
> query
|
| 		SELECT
| 			schema_name
| 		FROM information_schema.schemata
| 		WHERE schema_name NOT IN ('INFORMATION_SCHEMA', 'information_schema', 'system')
| 		ORDER BY schema_name
< result
< column "schema_name" "String" "string"
< row
: string
| default
< row
: string
| omni_testsyncinstance
