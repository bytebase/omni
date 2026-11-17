-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDumpGaps/a_materialized_view's_inner_table default#1
> exec
| DROP DATABASE IF EXISTS `omni_gap_mv`
< ok 0 -
> exec
| CREATE DATABASE `omni_gap_mv`
< ok 0 -
> exec
| CREATE TABLE `omni_gap_mv`.src (id UInt64) ENGINE = MergeTree ORDER BY id
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW `omni_gap_mv`.mv ENGINE = MergeTree ORDER BY id AS SELECT id FROM `omni_gap_mv`.src
< ok 0 -
> exec
| SET format_display_secrets_in_show_and_select = 1
< ok 0 -
> query
|
| 		SELECT
| 			name
| 		FROM system.tables
| 		WHERE database = ? AND NOT match(name, ?)
| 		ORDER BY name
: string
| omni_gap_mv
: string
| ^[.]inner_id[.]([a-z_]+[.])?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
< result
< column "name" "String" "string"
< row
: string
| mv
< row
: string
| src
> query
| SHOW CREATE TABLE `omni_gap_mv`.`mv`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_gap_mv.mv
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_gap_mv.src
> query
| SHOW CREATE TABLE `omni_gap_mv`.`src`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_gap_mv.src
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_gap_mv.mv
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_gap_mv.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_gap_mv mv (children 5)
< row
: string
|  Identifier omni_gap_mv
< row
: string
|  Identifier mv
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
|        TableIdentifier omni_gap_mv.src
< row
: string
|  ViewTargets (children 1)
< row
: string
|   Storage definition (children 3)
< row
: string
|    Function MergeTree
< row
: string
|    Identifier id
< row
: string
|    Set
> exec
| DROP DATABASE IF EXISTS `omni_gap_mv`
< ok 0 -
> exec
| CREATE DATABASE `omni_gap_mv`
< ok 0 -
> begin
< ok
> exec
| --
| -- Table structure for `src`
| --
| CREATE TABLE omni_gap_mv.src
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
| -- Materialized view structure for `mv`
| --
| CREATE MATERIALIZED VIEW omni_gap_mv.mv
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_gap_mv.src;
< ok 0 -
> commit
< ok
> query
| SELECT count() FROM system.tables WHERE database = ? AND name NOT LIKE '.%'
: string
| omni_gap_mv
< result
< column "count()" "UInt64" "uint64"
< row
: uint64 2
> query
| SELECT count() FROM system.tables WHERE database = ? AND name LIKE '.%'
: string
| omni_gap_mv
< result
< column "count()" "UInt64" "uint64"
< row
: uint64 1

== TestDumpGaps/the_database_itself default#1
> exec
| DROP DATABASE IF EXISTS `omni_gap_database`
< ok 0 -
> exec
| CREATE DATABASE `omni_gap_database` ENGINE = Atomic COMMENT 'a comment the dump leaves out'
< ok 0 -
> exec
| CREATE TABLE `omni_gap_database`.t (id UInt64) ENGINE = MergeTree ORDER BY id
< ok 0 -
> exec
| SET format_display_secrets_in_show_and_select = 1
< ok 0 -
> query
|
| 		SELECT
| 			name
| 		FROM system.tables
| 		WHERE database = ? AND NOT match(name, ?)
| 		ORDER BY name
: string
| omni_gap_database
: string
| ^[.]inner_id[.]([a-z_]+[.])?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
< result
< column "name" "String" "string"
< row
: string
| t
> query
| SHOW CREATE TABLE `omni_gap_database`.`t`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_gap_database.t
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> exec
| DROP DATABASE IF EXISTS `omni_gap_database`
< ok 0 -
> exec
| CREATE DATABASE `omni_gap_database`
< ok 0 -
> begin
< ok
> exec
| --
| -- Table structure for `t`
| --
| CREATE TABLE omni_gap_database.t
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192;
< ok 0 -
> commit
< ok
> query
| SELECT comment FROM system.databases WHERE name = ?
: string
| omni_gap_database
< result
< column "comment" "String" "string"
< row
: string
|
