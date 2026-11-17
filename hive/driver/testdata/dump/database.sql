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

--
-- View structure for `z_base`
--
CREATE VIEW `omni_dump`.`z_base` AS SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`;

--
-- View structure for `a_top`
--
CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` FROM `omni_dump`.`z_base`;

