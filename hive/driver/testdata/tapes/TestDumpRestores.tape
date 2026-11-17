-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDumpRestores default#1
> exec
| -- The schema the dump tests read. The external table and the two views are what
| -- TestDumpRestores restores; the managed table and the partitioned one carry
| -- what a dump leaves out, so each line of theirs ends with the gap marker.
| CREATE EXTERNAL TABLE omni_dump.ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_users'
< ok - -
> exec
| CREATE VIEW omni_dump.z_base AS SELECT id FROM omni_dump.ext_users
< ok - -
> exec
| -- Only the order of the views, not their names, can replay this one.
| CREATE VIEW omni_dump.a_top AS SELECT id FROM omni_dump.z_base
< ok - -
> query
| SHOW TABLES FROM `omni_dump`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| ext_users
< row
: string
| z_base
> query
| DESCRIBE FORMATTED `omni_dump`.`a_top`
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
| SELECT id FROM omni_dump.z_base
: null
< row
: string
| Expanded Query:     
: string
| SELECT `z_base`.`id` FROM `omni_dump`.`z_base`
: null
> query
| DESCRIBE FORMATTED `omni_dump`.`ext_users`
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
| DESCRIBE FORMATTED `omni_dump`.`z_base`
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
| SELECT id FROM omni_dump.ext_users
: null
< row
: string
| Expanded Query:     
: string
| SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`
: null
> query
| SHOW TABLES IN `omni_dump`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| ext_users
< row
: string
| z_base
> query
| SHOW CREATE TABLE `omni_dump`.`a_top`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` FROM `omni_dump`.`z_base`
> query
| SHOW CREATE TABLE `omni_dump`.`ext_users`
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
| SHOW CREATE TABLE `omni_dump`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`
> exec
| DROP DATABASE IF EXISTS `omni_dump` CASCADE
< ok - -
> exec
| CREATE DATABASE `omni_dump`
< ok - -
> exec
| USE `omni_dump`
< ok - -
> exec
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
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`
< ok - -
> exec
|
|
| --
| -- View structure for `a_top`
| --
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` FROM `omni_dump`.`z_base`
< ok - -
> query
| SHOW TABLES FROM `omni_dump`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| ext_users
< row
: string
| z_base
> query
| DESCRIBE FORMATTED `omni_dump`.`a_top`
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
| SELECT z_base.id FROM omni_dump.z_base
: null
< row
: string
| Expanded Query:     
: string
| SELECT `z_base`.`id` FROM `omni_dump`.`z_base`
: null
> query
| DESCRIBE FORMATTED `omni_dump`.`ext_users`
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
| DESCRIBE FORMATTED `omni_dump`.`z_base`
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
| SELECT ext_users.id FROM omni_dump.ext_users
: null
< row
: string
| Expanded Query:     
: string
| SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`
: null
> query
| SHOW TABLES IN `omni_dump`
< result
< column "tab_name" "STRING"
< row
: string
| a_top
< row
: string
| ext_users
< row
: string
| z_base
> query
| SHOW CREATE TABLE `omni_dump`.`a_top`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `z_base`.`id` FROM `omni_dump`.`z_base`
> query
| SHOW CREATE TABLE `omni_dump`.`ext_users`
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
| SHOW CREATE TABLE `omni_dump`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `ext_users`.`id` FROM `omni_dump`.`ext_users`
