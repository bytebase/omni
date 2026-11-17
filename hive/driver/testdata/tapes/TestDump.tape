-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDump default#1
> exec
| -- The schema the dump tests read. The external table, the table whose name ends
| -- in "__" and the two views are what TestDumpRestores restores; the managed
| -- table and the partitioned one carry what a dump leaves out, so each line of
| -- theirs ends with the gap marker.
| CREATE EXTERNAL TABLE omni_dump.ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_users'
< ok - -
> exec
| -- A name ending in "__", which Hive 4 has no index to give one to: every name
| -- SHOW TABLES answers is a table's.
| CREATE TABLE omni_dump.events__ (id INT) STORED AS PARQUET
< ok - -
> exec
| -- A view holding the other's qualified name in a literal, which is not a read:
| -- counting it would close a cycle with a_top's and leave the order to names.
| CREATE VIEW omni_dump.z_base AS SELECT id, '`omni_dump`.`a_top`' AS tag FROM omni_dump.ext_users
< ok - -
> exec
| -- Only the order of the views, not their names, can replay this one, and the
| -- backquoted alias holds a quote the scanner has to leave alone.
| CREATE VIEW omni_dump.a_top AS SELECT id AS `it's` FROM omni_dump.z_base
< ok - -
> exec
| -- A table Hive translated to external and then made managed: the only managed
| -- table the pinned server can make, and the one a dump cannot restore as
| -- managed.
| CREATE TABLE omni_dump.orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET -- not restorable: TestDumpGaps
< ok - -
> exec
| ALTER TABLE omni_dump.orders SET TBLPROPERTIES('EXTERNAL'='FALSE') -- not restorable: TestDumpGaps
< ok - -
> exec
| INSERT INTO omni_dump.orders VALUES (1, 1.5, 'n1') -- not restorable: TestDumpGaps
< ok - -
> exec
| -- A partitioned managed table with one partition and two rows.
| CREATE TABLE omni_dump.events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET -- not restorable: TestDumpGaps
< ok - -
> exec
| ALTER TABLE omni_dump.events SET TBLPROPERTIES('EXTERNAL'='FALSE') -- not restorable: TestDumpGaps
< ok - -
> exec
| ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-01') -- not restorable: TestDumpGaps
< ok - -
> exec
| INSERT INTO omni_dump.events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y') -- not restorable: TestDumpGaps
< ok - -
> query
| SHOW TABLES IN `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| events
< row
: string
| events__
< row
: string
| ext_users
< row
: string
| orders
< row
: string
| z_base
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`a_top`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` AS `it's` FROM `omni_dump`.`z_base`
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`events`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_dump`.`events`(
< row
: string
|   `id` int, 
< row
: string
|   `body` string)
< row
: string
| PARTITIONED BY ( 
< row
: string
|   `dt` string)
< row
: string
| ROW FORMAT SERDE 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
< row
: string
| STORED AS INPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
< row
: string
| OUTPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
< row
: string
| LOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_dump.db/events'
< row
: string
| TBLPROPERTIES (
< row
: string
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
< row
: string
|   'bucketing_version'='2', 
< row
: string
|   'external.table.purge'='TRUE', 
< row
: string
|   'last_modified_by'='hive', 
< row
: string
|   'last_modified_time'='1000000000', 
< row
: string
|   'transient_lastDdlTime'='1000000000')
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`events__`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_dump`.`events__`(
< row
: string
|   `id` int)
< row
: string
| ROW FORMAT SERDE 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
< row
: string
| STORED AS INPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
< row
: string
| OUTPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
< row
: string
| LOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_dump.db/events__'
< row
: string
| TBLPROPERTIES (
< row
: string
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
< row
: string
|   'bucketing_version'='2', 
< row
: string
|   'external.table.purge'='TRUE', 
< row
: string
|   'transient_lastDdlTime'='1000000000')
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`ext_users`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_dump`.`ext_users`(
< row
: string
|   `id` int COMMENT 'the id', 
< row
: string
|   `name` string)
< row
: string
| ROW FORMAT SERDE 
< row
: string
|   'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe' 
< row
: string
| WITH SERDEPROPERTIES ( 
< row
: string
|   'field.delim'=',', 
< row
: string
|   'serialization.format'=',') 
< row
: string
| STORED AS INPUTFORMAT 
< row
: string
|   'org.apache.hadoop.mapred.TextInputFormat' 
< row
: string
| OUTPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
< row
: string
| LOCATION
< row
: string
|   'file:/tmp/omni/dump_ext_users'
< row
: string
| TBLPROPERTIES (
< row
: string
|   'bucketing_version'='2', 
< row
: string
|   'transient_lastDdlTime'='1000000000')
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`orders`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_dump`.`orders`(
< row
: string
|   `id` bigint, 
< row
: string
|   `amount` decimal(10,2), 
< row
: string
|   `note` string COMMENT 'a note')
< row
: string
| ROW FORMAT SERDE 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
< row
: string
| STORED AS INPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
< row
: string
| OUTPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
< row
: string
| LOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_dump.db/orders'
< row
: string
| TBLPROPERTIES (
< row
: string
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
< row
: string
|   'bucketing_version'='2', 
< row
: string
|   'external.table.purge'='TRUE', 
< row
: string
|   'last_modified_by'='hive', 
< row
: string
|   'last_modified_time'='1000000000', 
< row
: string
|   'transient_lastDdlTime'='1000000000')
> query
| SHOW CREATE TABLE `OMNI_DUMP`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `ext_users`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_dump`.`ext_users`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`a_top`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@a_top","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@a_top]"},{"tablename":"omni_dump@ext_users","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_dump@z_base]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`z_base`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@ext_users","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_dump@z_base]"}],"input_partitions":[]}
