-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDumpRestores default#1
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
| -- A view in another database named as a_top is, which z_base reads: a
| -- dependency outside the database dumped is not one of its views.
| DROP DATABASE IF EXISTS omni_other CASCADE
< ok - -
> exec
| CREATE DATABASE omni_other
< ok - -
> exec
| CREATE TABLE omni_other.t (id INT) STORED AS PARQUET
< ok - -
> exec
| CREATE VIEW omni_other.a_top AS SELECT id FROM omni_other.t
< ok - -
> exec
| -- A view holding the other's qualified name in a literal, which is not a read:
| -- counting it would close a cycle with a_top's and leave the order to names.
| CREATE VIEW omni_dump.z_base AS SELECT id, '`omni_dump`.`a_top`' AS tag FROM omni_other.a_top
< ok - -
> exec
| -- Only the order of the views, not their names, can replay this one, and the
| -- backquoted alias holds a quote the scanner has to leave alone.
| CREATE VIEW omni_dump.a_top AS SELECT id AS `it's` FROM omni_dump.z_base
< ok - -
> query
| SHOW TABLES FROM `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| events__
< row
: string
| ext_users
< row
: string
| z_base
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`a_top`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| it's
: string
| int
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Table Type:         
: string
| VIRTUAL_VIEW        
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| null                
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
|
: null
: null
< row
: string
| # View Information
: null
: null
< row
: string
| Original Query:     
: string
| SELECT id AS it's FROM omni_dump.z_base
: null
< row
: string
| Expanded Query:     
: string
| SELECT `z_base`.`id` AS `it's` FROM `omni_dump`.`z_base`
: null
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`ext_users`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| int
: string
| the id
< row
: string
| name
: string
| string
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Location:           
: string
| file:/tmp/omni/dump_ext_users
: null
< row
: string
| Table Type:         
: string
| EXTERNAL_TABLE      
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| COLUMN_STATS_ACCURATE
: string
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\",\"name\":\"true\"}}
< row
: string
|
: string
| EXTERNAL            
: string
| TRUE                
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| numFiles            
: string
| 0                   
< row
: string
|
: string
| numRows             
: string
| 0                   
< row
: string
|
: string
| rawDataSize         
: string
| 0                   
< row
: string
|
: string
| totalSize           
: string
| 0                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
| Storage Desc Params:
: null
: null
< row
: string
|
: string
| field.delim         
: string
| ,                   
< row
: string
|
: string
| serialization.format
: string
| ,                   
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`z_base`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| int
: string
|
< row
: string
| tag
: string
| string
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Table Type:         
: string
| VIRTUAL_VIEW        
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| null                
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
|
: null
: null
< row
: string
| # View Information
: null
: null
< row
: string
| Original Query:     
: string
| SELECT id, '`omni_dump`.`a_top`' AS tag FROM omni_other.a_top
: null
< row
: string
| Expanded Query:     
: string
| SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
: null
> query
| SHOW TABLES IN `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| events__
< row
: string
| ext_users
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
| SHOW CREATE TABLE `OMNI_DUMP`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`a_top`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@a_top","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@a_top]"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`z_base`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
> exec
| DROP DATABASE IF EXISTS `OMNI_DUMP` CASCADE
< ok - -
> exec
| CREATE DATABASE `OMNI_DUMP`
< ok - -
> exec
| USE `OMNI_DUMP`
< ok - -
> exec
| --
| -- Table structure for `events__`
| --
| CREATE EXTERNAL TABLE `omni_dump`.`events__`(
|   `id` int)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_dump.db/events__'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> exec
|
|
| --
| -- Table structure for `ext_users`
| --
| CREATE EXTERNAL TABLE `omni_dump`.`ext_users`(
|   `id` int COMMENT 'the id', 
|   `name` string)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe' 
| WITH SERDEPROPERTIES ( 
|   'field.delim'=',', 
|   'serialization.format'=',') 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.mapred.TextInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
| LOCATION
|   'file:/tmp/omni/dump_ext_users'
| TBLPROPERTIES (
|   'bucketing_version'='2', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> exec
|
|
| --
| -- View structure for `z_base`
| --
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
< ok - -
> exec
|
|
| --
| -- View structure for `a_top`
| --
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` AS `it's` FROM `omni_dump`.`z_base`
< ok - -
> query
| SHOW TABLES FROM `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| events__
< row
: string
| ext_users
< row
: string
| z_base
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`a_top`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| it's
: string
| int
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Table Type:         
: string
| VIRTUAL_VIEW        
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| null                
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
|
: null
: null
< row
: string
| # View Information
: null
: null
< row
: string
| Original Query:     
: string
| SELECT z_base.id AS it's FROM omni_dump.z_base
: null
< row
: string
| Expanded Query:     
: string
| SELECT `z_base`.`id` AS `it's` FROM `omni_dump`.`z_base`
: null
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`ext_users`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| int
: string
| the id
< row
: string
| name
: string
| string
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Location:           
: string
| file:/tmp/omni/dump_ext_users
: null
< row
: string
| Table Type:         
: string
| EXTERNAL_TABLE      
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| COLUMN_STATS_ACCURATE
: string
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\",\"name\":\"true\"}}
< row
: string
|
: string
| EXTERNAL            
: string
| TRUE                
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| numFiles            
: string
| 0                   
< row
: string
|
: string
| numRows             
: string
| 0                   
< row
: string
|
: string
| rawDataSize         
: string
| 0                   
< row
: string
|
: string
| totalSize           
: string
| 0                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
| Storage Desc Params:
: null
: null
< row
: string
|
: string
| field.delim         
: string
| ,                   
< row
: string
|
: string
| serialization.format
: string
| ,                   
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`z_base`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| int
: string
|
< row
: string
| tag
: string
| string
: string
|
< row
: string
|
: null
: null
< row
: string
| # Detailed Table Information
: null
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| OwnerType:          
: string
| USER                
: null
< row
: string
| Owner:              
: string
| hive                
: null
< row
: string
| CreateTime:         
: string
| Thu Jan  1 00:00:00 UTC 2026
: null
< row
: string
| LastAccessTime:     
: string
| UNKNOWN             
: null
< row
: string
| Retention:          
: string
| 0                   
: null
< row
: string
| Table Type:         
: string
| VIRTUAL_VIEW        
: null
< row
: string
| Table Parameters:
: null
: null
< row
: string
|
: string
| bucketing_version   
: string
| 2                   
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000000          
< row
: string
|
: null
: null
< row
: string
| # Storage Information
: null
: null
< row
: string
| SerDe Library:      
: string
| null                
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.mapred.TextInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat
: null
< row
: string
| Compressed:         
: string
| No                  
: null
< row
: string
| Num Buckets:        
: string
| -1                  
: null
< row
: string
| Bucket Columns:     
: string
| []                  
: null
< row
: string
| Sort Columns:       
: string
| []                  
: null
< row
: string
|
: null
: null
< row
: string
| # View Information
: null
: null
< row
: string
| Original Query:     
: string
| SELECT a_top.id, '`omni_dump`.`a_top`' AS tag FROM omni_other.a_top
: null
< row
: string
| Expanded Query:     
: string
| SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
: null
> query
| SHOW TABLES IN `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| events__
< row
: string
| ext_users
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
| SHOW CREATE TABLE `OMNI_DUMP`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`a_top`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@a_top","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@a_top]"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`z_base`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
