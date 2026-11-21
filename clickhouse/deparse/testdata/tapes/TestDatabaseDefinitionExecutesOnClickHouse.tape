-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDatabaseDefinitionExecutesOnClickHouse default#1
> exec
| --
| -- Table structure for `events`
| --
| CREATE TABLE `events` (
|   `id` UInt64,
|   `name` String COMMENT 'the \'event\'\\name',
|   `at` DateTime64(3, 'UTC'),
|   `amount` Nullable(Decimal(12, 2)) DEFAULT NULL,
|   `tags` Array(String),
|   `attrs` Map(String, String),
|   `level` Enum8('info' = 1, 'warn' = 2) DEFAULT 'info',
|   `seen` Date DEFAULT '2026-01-02',
|   `retries` UInt32 DEFAULT 3,
|   `created_at` DateTime DEFAULT now(),
|   INDEX `name_bloom` name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| COMMENT 'events\' log\\v1'
< ok 0 -
> exec
| --
| -- Table structure for `scratch`
| --
| CREATE TABLE `scratch` (
|   `k` Int32
| )
| ENGINE = ReplacingMergeTree
| ORDER BY tuple()
< ok 0 -
> exec
| --
| -- Table structure for `odd-table name`
| --
| CREATE TABLE `odd-table name` (
|   `select` Int32 DEFAULT 0,
|   `col with space` Nullable(String) DEFAULT NULL
| )
| ENGINE = Memory
| COMMENT 'line1\nline2'
< ok 0 -
> exec
| --
| -- View structure for `synced_view`
| --
| CREATE VIEW omni_testdatabasedefinitionexecutesonclickhouse.synced_view (`id` UInt64, `name` String) AS SELECT id, name FROM omni_testdatabasedefinitionexecutesonclickhouse.events
< ok 0 -
> exec
| --
| -- View structure for `query view`
| --
| CREATE OR REPLACE VIEW `query view` AS SELECT id, seen FROM events WHERE retries > 0 COMMENT 'filtered \'events\''
< ok 0 -
> query
| SELECT count() FROM system.tables WHERE database = $1
: string
| omni_testdatabasedefinitionexecutesonclickhouse
< result
< column "count()" "UInt64" "uint64"
< row
: uint64 5
> query
| SELECT column_default FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'events' AND column_name = 'seen'
: string
| omni_testdatabasedefinitionexecutesonclickhouse
< result
< column "column_default" "String" "string"
< row
: string
| '2026-01-02'
> query
| SELECT comment FROM system.tables WHERE database = $1 AND name = 'events'
: string
| omni_testdatabasedefinitionexecutesonclickhouse
< result
< column "comment" "String" "string"
< row
: string
| events' log\v1
> query
| SELECT comment FROM system.tables WHERE database = $1 AND name = 'query view'
: string
| omni_testdatabasedefinitionexecutesonclickhouse
< result
< column "comment" "String" "string"
< row
: string
| filtered 'events'
