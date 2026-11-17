-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDumpGaps/a_materialized_view's_definition default#1
> exec
| CREATE TABLE `omni_gap_materialized_view`.src (id INT) STORED AS PARQUET
< ok - -
> exec
| CREATE MATERIALIZED VIEW `omni_gap_materialized_view`.mv DISABLE REWRITE AS SELECT id FROM `omni_gap_materialized_view`.src
< ok - -
> query
| SHOW TABLES IN `omni_gap_materialized_view`
< result
< column "tab_name" "STRING"
< row
: string
| mv
< row
: string
| src
> query
| SHOW CREATE TABLE `omni_gap_materialized_view`.`mv`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_gap_materialized_view`.`mv`(
< row
: string
|   `id` int)
< row
: string
| ROW FORMAT SERDE 
< row
: string
|   'org.apache.hadoop.hive.ql.io.orc.OrcSerde' 
< row
: string
| STORED AS INPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' 
< row
: string
| OUTPUTFORMAT 
< row
: string
|   'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat'
< row
: string
| LOCATION
< row
: string
|   'file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/mv'
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
| SHOW CREATE TABLE `omni_gap_materialized_view`.`src`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_gap_materialized_view`.`src`(
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
|   'file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/src'
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
| SHOW TABLES FROM `omni_gap_materialized_view`
< result
< column "tab_name" "STRING"
< row
: string
| mv
< row
: string
| src
> query
| DESCRIBE FORMATTED `omni_gap_materialized_view`.`mv`
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
| omni_gap_materialized_view
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
| file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/mv
: null
< row
: string
| Table Type:         
: string
| MATERIALIZED_VIEW   
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
| bucketing_version   
: string
| 2                   
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
| org.apache.hadoop.hive.ql.io.orc.OrcSerde
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.orc.OrcInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat
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
| # Materialized View Information
: null
: null
< row
: string
| Original Query:     
: string
| SELECT id FROM omni_gap_materialized_view.src
: null
< row
: string
| Expanded Query:     
: string
| SELECT `src`.`id` FROM `omni_gap_materialized_view`.`src`
: null
< row
: string
| Rewrite Enabled:    
: string
| No                  
: null
< row
: string
| Outdated for Rewriting:
: string
| Unknown             
: null
< row
: string
|
: null
: null
< row
: string
| # Materialized View Source table information
: null
: null
< row
: string
| Table name
: string
| Snapshot
: null
< row
: string
| omni_gap_materialized_view.src
: string
| N/A
: null
> query
| DESCRIBE FORMATTED `omni_gap_materialized_view`.`src`
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
| omni_gap_materialized_view
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
| file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/src
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
| DROP DATABASE IF EXISTS `omni_gap_materialized_view` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_materialized_view`
< ok - -
> exec
| USE `omni_gap_materialized_view`
< ok - -
> exec
| --
| -- Table structure for `mv`
| --
| CREATE TABLE `omni_gap_materialized_view`.`mv`(
|   `id` int)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.orc.OrcSerde' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/mv'
| TBLPROPERTIES (
|   'bucketing_version'='2', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> exec
|
|
| --
| -- Table structure for `src`
| --
| CREATE EXTERNAL TABLE `omni_gap_materialized_view`.`src`(
|   `id` int)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/src'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> query
| SHOW TABLES FROM `omni_gap_materialized_view`
< result
< column "tab_name" "STRING"
< row
: string
| mv
< row
: string
| src
> query
| DESCRIBE FORMATTED `omni_gap_materialized_view`.`mv`
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
| omni_gap_materialized_view
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
| file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/mv
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
| org.apache.hadoop.hive.ql.io.orc.OrcSerde
: null
< row
: string
| InputFormat:        
: string
| org.apache.hadoop.hive.ql.io.orc.OrcInputFormat
: null
< row
: string
| OutputFormat:       
: string
| org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat
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
| DESCRIBE FORMATTED `omni_gap_materialized_view`.`src`
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
| omni_gap_materialized_view
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
| file:/opt/hive/data/warehouse/omni_gap_materialized_view.db/src
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

== TestDumpGaps/a_partitioned_table's_partitions default#1
> exec
| CREATE TABLE `omni_gap_partition`.events (id INT) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE `omni_gap_partition`.events SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| ALTER TABLE `omni_gap_partition`.events ADD PARTITION (dt='2024-01-01')
< ok - -
> query
| SHOW TABLES IN `omni_gap_partition`
< result
< column "tab_name" "STRING"
< row
: string
| events
> query
| SHOW CREATE TABLE `omni_gap_partition`.`events`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_gap_partition`.`events`(
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
|   'file:/opt/hive/data/warehouse/omni_gap_partition.db/events'
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
| SHOW TABLES FROM `omni_gap_partition`
< result
< column "tab_name" "STRING"
< row
: string
| events
> query
| DESCRIBE FORMATTED `omni_gap_partition`.`events`
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
| omni_gap_partition  
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
| file:/opt/hive/data/warehouse/omni_gap_partition.db/events
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
| 1                   
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
| SHOW PARTITIONS `omni_gap_partition`.`events`
< result
< column "partition" "STRING"
< row
: string
| dt=2024-01-01
> exec
| DROP DATABASE IF EXISTS `omni_gap_partition` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_partition`
< ok - -
> exec
| USE `omni_gap_partition`
< ok - -
> exec
| --
| -- Table structure for `events`
| --
| CREATE TABLE `omni_gap_partition`.`events`(
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
|   'file:/opt/hive/data/warehouse/omni_gap_partition.db/events'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'last_modified_by'='hive', 
|   'last_modified_time'='1000000000', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> query
| SHOW TABLES FROM `omni_gap_partition`
< result
< column "tab_name" "STRING"
< row
: string
| events
> query
| DESCRIBE FORMATTED `omni_gap_partition`.`events`
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
| omni_gap_partition  
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
| file:/opt/hive/data/warehouse/omni_gap_partition.db/events
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
| {\"BASIC_STATS\":\"true\"}
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

== TestDumpGaps/a_table_Hive_made_managed_by_an_ALTER default#1
> exec
| CREATE TABLE `omni_gap_managed`.orders (id BIGINT) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE `omni_gap_managed`.orders SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> query
| SHOW TABLES IN `omni_gap_managed`
< result
< column "tab_name" "STRING"
< row
: string
| orders
> query
| SHOW CREATE TABLE `omni_gap_managed`.`orders`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_gap_managed`.`orders`(
< row
: string
|   `id` bigint)
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
|   'file:/opt/hive/data/warehouse/omni_gap_managed.db/orders'
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
| SHOW TABLES FROM `omni_gap_managed`
< result
< column "tab_name" "STRING"
< row
: string
| orders
> query
| DESCRIBE FORMATTED `omni_gap_managed`.`orders`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| bigint
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
| omni_gap_managed    
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
| file:/opt/hive/data/warehouse/omni_gap_managed.db/orders
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\"}}
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
| SHOW PARTITIONS `omni_gap_managed`.`orders`
< error
| 10241
| Error while compiling statement: FAILED: Execution Error, return code 10241 from org.apache.hadoop.hive.ql.ddl.DDLTask. Table omni_gap_managed.orders is not a partitioned table
> exec
| DROP DATABASE IF EXISTS `omni_gap_managed` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_managed`
< ok - -
> exec
| USE `omni_gap_managed`
< ok - -
> exec
| --
| -- Table structure for `orders`
| --
| CREATE TABLE `omni_gap_managed`.`orders`(
|   `id` bigint)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_managed.db/orders'
| TBLPROPERTIES (
|   'TRANSLATED_TO_EXTERNAL'='TRUE', 
|   'bucketing_version'='2', 
|   'external.table.purge'='TRUE', 
|   'last_modified_by'='hive', 
|   'last_modified_time'='1000000000', 
|   'transient_lastDdlTime'='1000000000')
< ok - -
> query
| SHOW TABLES FROM `omni_gap_managed`
< result
< column "tab_name" "STRING"
< row
: string
| orders
> query
| DESCRIBE FORMATTED `omni_gap_managed`.`orders`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| id
: string
| bigint
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
| omni_gap_managed    
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
| file:/opt/hive/data/warehouse/omni_gap_managed.db/orders
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

== TestDumpGaps/a_view's_comment default#1
> exec
| CREATE TABLE `omni_gap_view_comment`.t (id INT) STORED AS PARQUET
< ok - -
> exec
| CREATE VIEW `omni_gap_view_comment`.v COMMENT 'a comment the dump leaves out' AS SELECT id FROM `omni_gap_view_comment`.t
< ok - -
> query
| SHOW TABLES IN `omni_gap_view_comment`
< result
< column "tab_name" "STRING"
< row
: string
| t
< row
: string
| v
> query
| SHOW CREATE TABLE `omni_gap_view_comment`.`t`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_gap_view_comment`.`t`(
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
|   'file:/opt/hive/data/warehouse/omni_gap_view_comment.db/t'
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
| SHOW CREATE TABLE `omni_gap_view_comment`.`v`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_gap_view_comment`.`v` AS SELECT `t`.`id` FROM `omni_gap_view_comment`.`t`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `omni_gap_view_comment`.`v`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_gap_view_comment@v","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_gap_view_comment@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_gap_view_comment@v]"}],"input_partitions":[]}
> query
| SHOW TABLES FROM `omni_gap_view_comment`
< result
< column "tab_name" "STRING"
< row
: string
| t
< row
: string
| v
> query
| DESCRIBE FORMATTED `omni_gap_view_comment`.`t`
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
| omni_gap_view_comment
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
| file:/opt/hive/data/warehouse/omni_gap_view_comment.db/t
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
| DESCRIBE FORMATTED `omni_gap_view_comment`.`v`
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
| omni_gap_view_comment
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
| comment             
: string
| a comment the dump leaves out
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
| SELECT id FROM omni_gap_view_comment.t
: null
< row
: string
| Expanded Query:     
: string
| SELECT `t`.`id` FROM `omni_gap_view_comment`.`t`
: null
> exec
| DROP DATABASE IF EXISTS `omni_gap_view_comment` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_gap_view_comment`
< ok - -
> exec
| USE `omni_gap_view_comment`
< ok - -
> exec
| --
| -- Table structure for `t`
| --
| CREATE EXTERNAL TABLE `omni_gap_view_comment`.`t`(
|   `id` int)
| ROW FORMAT SERDE 
|   'org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe' 
| STORED AS INPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat' 
| OUTPUTFORMAT 
|   'org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat'
| LOCATION
|   'file:/opt/hive/data/warehouse/omni_gap_view_comment.db/t'
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
| -- View structure for `v`
| --
| CREATE VIEW `omni_gap_view_comment`.`v` AS SELECT `t`.`id` FROM `omni_gap_view_comment`.`t`
< ok - -
> query
| SHOW TABLES FROM `omni_gap_view_comment`
< result
< column "tab_name" "STRING"
< row
: string
| t
< row
: string
| v
> query
| DESCRIBE FORMATTED `omni_gap_view_comment`.`t`
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
| omni_gap_view_comment
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
| file:/opt/hive/data/warehouse/omni_gap_view_comment.db/t
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
| DESCRIBE FORMATTED `omni_gap_view_comment`.`v`
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
| omni_gap_view_comment
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
| SELECT t.id FROM omni_gap_view_comment.t
: null
< row
: string
| Expanded Query:     
: string
| SELECT `t`.`id` FROM `omni_gap_view_comment`.`t`
: null

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
