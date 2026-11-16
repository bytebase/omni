-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestSyncDatabase default#1
> exec
| CREATE EXTERNAL TABLE ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/ext_users'
< ok - -
> exec
| CREATE TABLE orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE orders SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO orders VALUES (1, 1.5, 'n1')
< ok - -
> exec
| CREATE TABLE events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE events SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| ALTER TABLE events ADD PARTITION (dt='2024-01-01')
< ok - -
> exec
| INSERT INTO events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')
< ok - -
> exec
| CREATE VIEW v_users AS SELECT id, name FROM ext_users
< ok - -
> query
| SHOW TABLES FROM `omni_sync`
< result
< column "tab_name" "STRING"
< row
: string
| events
< row
: string
| ext_users
< row
: string
| orders
< row
: string
| v_users
> query
| DESCRIBE FORMATTED `omni_sync`.`events`
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
| body
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
| omni_sync           
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
| Thu Jan  1 00:00:01 UTC 2026
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
| file:/opt/hive/data/warehouse/omni_sync.db/events
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
| 1000000002          
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
| 2                   
< row
: string
|
: string
| rawDataSize         
: string
| 76                  
< row
: string
|
: string
| totalSize           
: string
| 557                 
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000002          
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
| SHOW PARTITIONS `omni_sync`.`events`
< result
< column "partition" "STRING"
< row
: string
| dt=2024-01-01
> query
| DESCRIBE FORMATTED `omni_sync`.`ext_users`
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
| omni_sync           
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
| Thu Jan  1 00:00:03 UTC 2026
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
| file:/tmp/omni/ext_users
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
| 1000000004          
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
| DESCRIBE FORMATTED `omni_sync`.`orders`
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
| amount
: string
| decimal(10,2)
: string
|
< row
: string
| note
: string
| string
: string
| a note
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
| omni_sync           
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
| Thu Jan  1 00:00:05 UTC 2026
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
| file:/opt/hive/data/warehouse/omni_sync.db/orders
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"amount\":\"true\",\"id\":\"true\",\"note\":\"true\"}}
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
| 1000000006          
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
| 1                   
< row
: string
|
: string
| rawDataSize         
: string
| 105                 
< row
: string
|
: string
| totalSize           
: string
| 774                 
< row
: string
|
: string
| transient_lastDdlTime
: string
| 1000000002          
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
| SHOW PARTITIONS `omni_sync`.`orders`
< error
| 10241
| Error while compiling statement: FAILED: Execution Error, return code 10241 from org.apache.hadoop.hive.ql.ddl.DDLTask. Table omni_sync.orders is not a partitioned table
> query
| DESCRIBE FORMATTED `omni_sync`.`v_users`
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
| omni_sync           
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
| Thu Jan  1 00:00:07 UTC 2026
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
| 1000000008          
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
| SELECT id, name FROM ext_users
: null
< row
: string
| Expanded Query:     
: string
| SELECT `ext_users`.`id`, `ext_users`.`name` FROM `omni_sync`.`ext_users`
: null
