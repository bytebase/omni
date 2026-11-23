--
-- Table structure for `a_child`
--
CREATE EXTERNAL TABLE `omni_dump`.`a_child`(
  `id` int, 
  `parent_id` int)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/a_child'
TBLPROPERTIES (
  'TRANSLATED_TO_EXTERNAL'='TRUE', 
  'bucketing_version'='2', 
  'external.table.purge'='TRUE', 
  'transient_lastDdlTime'='1000000000');

--
-- Table structure for `events`
--
CREATE TABLE `omni_dump`.`events`(
  `id` int, 
  `body` string)
PARTITIONED BY ( 
  `dt` string)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/events'
TBLPROPERTIES (
  'TRANSLATED_TO_EXTERNAL'='TRUE', 
  'bucketing_version'='2', 
  'external.table.purge'='TRUE', 
  'last_modified_by'='hive', 
  'last_modified_time'='1000000000', 
  'transient_lastDdlTime'='1000000000');

ALTER TABLE `omni_dump`.`events` SET TBLPROPERTIES('EXTERNAL'='FALSE');

ALTER TABLE `omni_dump`.`events` ADD PARTITION (`dt`='2024-01-01') LOCATION 'file:/opt/hive/data/warehouse/omni_dump.db/events/dt=2024-01-01';

ALTER TABLE `omni_dump`.`events` ADD PARTITION (`dt`='2024-01-02') LOCATION 'file:/tmp/omni/dump_events_dt_2024-01-02';

ALTER TABLE `omni_dump`.`events` ADD PARTITION (`dt`='2024-01-03') LOCATION 'file:/opt/hive/data/warehouse/omni_dump.db/events/dt=2024-01-03';

ALTER TABLE `omni_dump`.`events` PARTITION (`dt`='2024-01-03') SET FILEFORMAT INPUTFORMAT 'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat' SERDE 'org.apache.hadoop.hive.ql.io.orc.OrcSerde';

ALTER TABLE `omni_dump`.`events` PARTITION (`dt`='2024-01-03') SET SERDEPROPERTIES ('orc.compress'='ZLIB');

--
-- Table structure for `events__`
--
CREATE EXTERNAL TABLE `omni_dump`.`events__`(
  `id` int)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/events__'
TBLPROPERTIES (
  'TRANSLATED_TO_EXTERNAL'='TRUE', 
  'bucketing_version'='2', 
  'external.table.purge'='TRUE', 
  'transient_lastDdlTime'='1000000000');

--
-- Table structure for `ext_events`
--
CREATE EXTERNAL TABLE `omni_dump`.`ext_events`(
  `id` int)
PARTITIONED BY ( 
  `dt` string)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe' 
WITH SERDEPROPERTIES ( 
  'field.delim'='|', 
  'serialization.format'='|') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  'file:/tmp/omni/dump_ext_events'
TBLPROPERTIES (
  'bucketing_version'='2', 
  'transient_lastDdlTime'='1000000000');

ALTER TABLE `omni_dump`.`ext_events` ADD PARTITION (`dt`='2024-01-01') LOCATION 'file:/tmp/omni/dump_ext_events/dt=2024-01-01';

ALTER TABLE `omni_dump`.`ext_events` ADD PARTITION (`dt`='2024-01-02') LOCATION 'file:/tmp/omni/dump_ext_events_dt_2024-01-02';

ALTER TABLE `omni_dump`.`ext_events` PARTITION (`dt`='2024-01-02') UNSET SERDEPROPERTIES ('field.delim');

--
-- Table structure for `ext_users`
--
CREATE EXTERNAL TABLE `omni_dump`.`ext_users`(
  `id` int COMMENT 'the id', 
  `name` string)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe' 
WITH SERDEPROPERTIES ( 
  'field.delim'=',', 
  'serialization.format'=',') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  'file:/tmp/omni/dump_ext_users'
TBLPROPERTIES (
  'bucketing_version'='2', 
  'transient_lastDdlTime'='1000000000');

--
-- Table structure for `orders`
--
CREATE TABLE `omni_dump`.`orders`(
  `id` bigint, 
  `amount` decimal(10,2), 
  `note` string COMMENT 'a note')
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/orders'
TBLPROPERTIES (
  'TRANSLATED_TO_EXTERNAL'='TRUE', 
  'bucketing_version'='2', 
  'external.table.purge'='TRUE', 
  'last_modified_by'='hive', 
  'last_modified_time'='1000000000', 
  'transient_lastDdlTime'='1000000000');

ALTER TABLE `omni_dump`.`orders` SET TBLPROPERTIES('EXTERNAL'='FALSE');

--
-- Table structure for `z_parent`
--
CREATE EXTERNAL TABLE `omni_dump`.`z_parent`(
  `id` int, 
  `code` string)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/z_parent'
TBLPROPERTIES (
  'TRANSLATED_TO_EXTERNAL'='TRUE', 
  'bucketing_version'='2', 
  'external.table.purge'='TRUE', 
  'transient_lastDdlTime'='1000000000');

ALTER TABLE `omni_dump`.`z_parent` ADD CONSTRAINT `pk_z_parent` PRIMARY KEY (`id`) DISABLE NOVALIDATE RELY;

ALTER TABLE `omni_dump`.`z_parent` ADD CONSTRAINT `uq_z_parent` UNIQUE (`code`) DISABLE NOVALIDATE RELY;

--
-- Foreign key constraints for `a_child`
--
ALTER TABLE `omni_dump`.`a_child` ADD CONSTRAINT `fk_a_child_parent` FOREIGN KEY (`parent_id`) REFERENCES `omni_dump`.`z_parent`(`id`) DISABLE NOVALIDATE RELY;

--
-- View structure for `v_docs`
--
CREATE VIEW `omni_dump`.`v_docs` (`id` COMMENT 'the id', `label` COMMENT 'it\'s a \\label') AS SELECT `ext_users`.`id`, `ext_users`.`name` FROM `omni_dump`.`ext_users`;

ALTER VIEW `omni_dump`.`v_docs` SET TBLPROPERTIES('comment'='documented view', 'team'='core');

--
-- View structure for `z_base`
--
CREATE VIEW `omni_dump`.`z_base` AS SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`;

--
-- View structure for `z_users`
--
CREATE VIEW `omni_dump`.`z_users` AS SELECT `ext_users`.`id` AS `1bad`, `ext_users`.`name` FROM `omni_dump`.`ext_users`;

ALTER VIEW `omni_dump`.`z_users` SET TBLPROPERTIES('comment'='the users\' rows');

--
-- View structure for `a_top`
--
CREATE VIEW `omni_dump`.`a_top` AS SELECT `middle`.`id` AS `it's` FROM `omni_other`.`middle`;

--
-- Materialized view structure for `a_users_mv`
--
CREATE MATERIALIZED VIEW `omni_dump`.`a_users_mv` DISABLE REWRITE
COMMENT 'the users materialized'
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.orc.OrcSerde' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/a_users_mv'
TBLPROPERTIES (
  'bucketing_version'='2', 
  'transient_lastDdlTime'='1000000000')
AS SELECT `z_users`.`1bad`
FROM `omni_dump`.`z_users`;

--
-- Materialized view structure for `p_events_mv`
--
CREATE MATERIALIZED VIEW `omni_dump`.`p_events_mv` DISABLE REWRITE
COMMENT 'says PARTITIONED BY (old)'
PARTITIONED ON (`dt`)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.orc.OrcSerde' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat'
LOCATION
  'file:/opt/hive/data/warehouse/omni_dump.db/p_events_mv'
TBLPROPERTIES (
  'bucketing_version'='2', 
  'transient_lastDdlTime'='1000000000')
AS SELECT id, dt FROM omni_dump.events;

