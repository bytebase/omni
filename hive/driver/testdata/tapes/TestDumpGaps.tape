-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDumpGaps/a_partition_whose_value_is_null default#1
> exec
| CREATE TABLE `omni_gap_null_partition`.p (id INT) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE `omni_gap_null_partition`.p SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO `omni_gap_null_partition`.p VALUES (1, null)
< ok - -
> query
| SHOW TABLES IN `omni_gap_null_partition`
< result
< column "tab_name" "STRING"
< row
: string
| p
> query
| DESCRIBE FORMATTED `omni_gap_null_partition`.`p`
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
|
: null
: null
< row
: string
| # Partition Information
: null
: null
< row
: string
| # col_name
: string
| data_type
: string
| comment
< row
: string
| dt
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
| omni_gap_null_partition
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
| file:/opt/hive/data/warehouse/omni_gap_null_partition.db/p
: null
< row
: string
| Table Type:         
: string
| MANAGED_TABLE       
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
| {\"BASIC_STATS\":\"true\"}
< row
: string
|
: string
| EXTERNAL            
: string
| FALSE               
< row
: string
|
: string
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
< row
: string
|
: string
| last_modified_by    
: string
| hive                
< row
: string
|
: string
| last_modified_time  
: string
| 1000000000          
< row
: string
|
: string
| numFiles            
: string
| 1                   
< row
: string
|
: string
| numPartitions       
: string
| 1                   
< row
: string
|
: string
| numRows             
: string
| 1                   
< row
: string
|
: string
| rawDataSize         
: string
| 33                  
< row
: string
|
: string
| totalSize           
: string
| 405                 
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> query
| SHOW CREATE TABLE `omni_gap_null_partition`.`p`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_gap_null_partition`.`p`(
< row
: string
|   `id` int)
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
|   'file:/opt/hive/data/warehouse/omni_gap_null_partition.db/p'
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
| SHOW PARTITIONS `omni_gap_null_partition`.`p`
< result
< column "partition" "STRING"
< row
: string
| dt=__HIVE_DEFAULT_PARTITION__
> query
| SHOW TABLES FROM `omni_gap_null_partition`
< result
< column "tab_name" "STRING"
< row
: string
| p
> query
| DESCRIBE FORMATTED `omni_gap_null_partition`.`p`
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
|
: null
: null
< row
: string
| # Partition Information
: null
: null
< row
: string
| # col_name
: string
| data_type
: string
| comment
< row
: string
| dt
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
| omni_gap_null_partition
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
| file:/opt/hive/data/warehouse/omni_gap_null_partition.db/p
: null
< row
: string
| Table Type:         
: string
| MANAGED_TABLE       
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
| {\"BASIC_STATS\":\"true\"}
< row
: string
|
: string
| EXTERNAL            
: string
| FALSE               
< row
: string
|
: string
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
< row
: string
|
: string
| last_modified_by    
: string
| hive                
< row
: string
|
: string
| last_modified_time  
: string
| 1000000000          
< row
: string
|
: string
| numFiles            
: string
| 1                   
< row
: string
|
: string
| numPartitions       
: string
| 1                   
< row
: string
|
: string
| numRows             
: string
| 1                   
< row
: string
|
: string
| rawDataSize         
: string
| 33                  
< row
: string
|
: string
| totalSize           
: string
| 405                 
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> query
| SHOW PARTITIONS `omni_gap_null_partition`.`p`
< result
< column "partition" "STRING"
< row
: string
| dt=__HIVE_DEFAULT_PARTITION__
> exec
| DROP DATABASE IF EXISTS `omni_gap_null_partition` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_null_partition`
< ok - -
> exec
| USE `omni_gap_null_partition`
< ok - -
> exec
| --
| -- Table structure for `p`
| --
| CREATE TABLE `omni_gap_null_partition`.`p`(
|   `id` int)
| PARTITIONED BY ( 
|   `dt` string)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_null_partition.db/p'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'last_modified_by'='hive', 
|   'last_modified_time'='1000000000', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> exec
|
|
| ALTER TABLE `omni_gap_null_partition`.`p` SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> query
| SHOW TABLES FROM `omni_gap_null_partition`
< result
< column "tab_name" "STRING"
< row
: string
| p
> query
| DESCRIBE FORMATTED `omni_gap_null_partition`.`p`
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
|
: null
: null
< row
: string
| # Partition Information
: null
: null
< row
: string
| # col_name
: string
| data_type
: string
| comment
< row
: string
| dt
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
| omni_gap_null_partition
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
| file:/opt/hive/data/warehouse/omni_gap_null_partition.db/p
: null
< row
: string
| Table Type:         
: string
| MANAGED_TABLE       
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
| {\"BASIC_STATS\":\"true\"}
< row
: string
|
: string
| EXTERNAL            
: string
| FALSE               
< row
: string
|
: string
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
< row
: string
|
: string
| last_modified_by    
: string
| hive                
< row
: string
|
: string
| last_modified_time  
: string
| 1000000000          
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
| numPartitions       
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> query
| SHOW PARTITIONS `omni_gap_null_partition`.`p`
< result
< column "partition" "STRING"
> exec
| ALTER TABLE `omni_gap_null_partition`.p ADD PARTITION (dt='__HIVE_DEFAULT_PARTITION__')
< error
| 10111
| Error while compiling statement: FAILED: SemanticException [Error 10111]: Partition value contains a reserved substring (User value: __HIVE_DEFAULT_PARTITION__ Reserved substring: __HIVE_DEFAULT_PARTITION__)
> exec
| ALTER TABLE `omni_gap_null_partition`.p ADD PARTITION (dt=null)
< error
| 10248
| Error while compiling statement: FAILED: SemanticException [Error 10248]: Cannot add partition column dt of type void as it cannot be converted to type string

== TestDumpGaps/a_view_whose_source_is_gone default#1
> exec
| CREATE TABLE `omni_gap_broken_view`.t (id INT) STORED AS PARQUET
< ok - -
> exec
| CREATE VIEW `omni_gap_broken_view`.v AS SELECT id FROM `omni_gap_broken_view`.t
< ok - -
> exec
| DROP TABLE `omni_gap_broken_view`.t
< ok - -
> query
| SHOW TABLES IN `omni_gap_broken_view`
< result
< column "tab_name" "STRING"
< row
: string
| v
> query
| DESCRIBE FORMATTED `omni_gap_broken_view`.`v`
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
| omni_gap_broken_view
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
| SELECT id FROM omni_gap_broken_view.t
: null
< row
: string
| Expanded Query:     
: string
| SELECT `t`.`id` FROM `omni_gap_broken_view`.`t`
: null
> query
| SHOW CREATE TABLE `omni_gap_broken_view`.`v`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_gap_broken_view`.`v` AS SELECT `t`.`id` FROM `omni_gap_broken_view`.`t`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `omni_gap_broken_view`.`v`
< error
| 40000
| Error while compiling statement: FAILED: SemanticException Line 1:21 Table not found 't' in definition of VIEW v [
| SELECT `t`.`id` FROM `omni_gap_broken_view`.`t`
| ] used as v at Line 1:33
> query
| SHOW TABLES FROM `omni_gap_broken_view`
< result
< column "tab_name" "STRING"
< row
: string
| v
> query
| DESCRIBE FORMATTED `omni_gap_broken_view`.`v`
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
| omni_gap_broken_view
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
| SELECT id FROM omni_gap_broken_view.t
: null
< row
: string
| Expanded Query:     
: string
| SELECT `t`.`id` FROM `omni_gap_broken_view`.`t`
: null
> exec
| DROP DATABASE IF EXISTS `omni_gap_broken_view` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_broken_view`
< ok - -
> exec
| USE `omni_gap_broken_view`
< ok - -
> exec
| --
| -- View structure for `v`
| --
| CREATE VIEW `omni_gap_broken_view`.`v` AS SELECT `t`.`id` FROM `omni_gap_broken_view`.`t`
< error
| 10001
| Error while compiling statement: FAILED: SemanticException [Error 10001]: Line 4:63 Table not found 't'

== TestDumpGaps/the_database_itself default#1
> exec
| DROP DATABASE IF EXISTS `omni_gap_database` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_database` COMMENT 'a comment the dump leaves out'
< ok - -
> exec
| CREATE TABLE `omni_gap_database`.t (id INT) STORED AS PARQUET
< ok - -
> query
| SHOW TABLES IN `omni_gap_database`
< result
< column "tab_name" "STRING"
< row
: string
| t
> query
| DESCRIBE FORMATTED `omni_gap_database`.`t`
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
| omni_gap_database   
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
| file:/opt/hive/data/warehouse/omni_gap_database.db/t
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\"}}
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
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> query
| SHOW CREATE TABLE `omni_gap_database`.`t`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_gap_database`.`t`(
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
|   'file:/opt/hive/data/warehouse/omni_gap_database.db/t'
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
| SHOW TABLES FROM `omni_gap_database`
< result
< column "tab_name" "STRING"
< row
: string
| t
> query
| DESCRIBE FORMATTED `omni_gap_database`.`t`
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
| omni_gap_database   
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
| file:/opt/hive/data/warehouse/omni_gap_database.db/t
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\"}}
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
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> exec
| DROP DATABASE IF EXISTS `omni_gap_database` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_database`
< ok - -
> exec
| USE `omni_gap_database`
< ok - -
> exec
| --
| -- Table structure for `t`
| --
| CREATE EXTERNAL TABLE `omni_gap_database`.`t`(
|   `id` int)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_database.db/t'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> query
| SHOW TABLES FROM `omni_gap_database`
< result
< column "tab_name" "STRING"
< row
: string
| t
> query
| DESCRIBE FORMATTED `omni_gap_database`.`t`
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
| omni_gap_database   
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
| file:/opt/hive/data/warehouse/omni_gap_database.db/t
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\"}}
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
| TRANSLATED_TO_EXTERNAL
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
| external.table.purge
: string
| TRUE                
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
| org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat
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
| serialization.format
: string
| 1                   
> query
| SHOW CREATE DATABASE `omni_gap_database`
< result
< column "createdb_stmt" "STRING"
< row
: string
| CREATE DATABASE `omni_gap_database`
< row
: string
| LOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_gap_database.db'
< row
: string
| MANAGEDLOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_gap_database.db'
