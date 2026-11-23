-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDump default#1
> exec
| -- The schema the dump tests read: an external table, a name ending in "__", two
| -- tables Hive translated to external and then made managed, one of them
| -- partitioned, a view with a comment and a materialized view over it, beside
| -- two databases that read each other's views.
| CREATE EXTERNAL TABLE omni_dump.ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_users'
< ok - -
> exec
| -- A name ending in "__", which Hive 4 has no index to give one to: every name
| -- SHOW TABLES answers is a table's.
| CREATE TABLE omni_dump.events__ (id INT) STORED AS PARQUET
< ok - -
> exec
| -- Another database, with a view named as this one's a_top is and one reading
| -- this database's z_base.
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
| -- A view of this database reading the other database's a_top — a name a_top has
| -- here too, which is not a read of this database's — and holding its reader's
| -- qualified name in a literal, which is not a read either: counting it would
| -- close a cycle with a_top's and leave the order to names.
| CREATE VIEW omni_dump.z_base AS SELECT id, '`omni_dump`.`a_top`' AS tag FROM omni_other.a_top
< ok - -
> exec
| -- A view of the other database reading z_base, and this database's a_top reading
| -- that view: creating a_top resolves the other database's view, so z_base has to
| -- be there first even though it is not a_top's stated parent. The backquoted
| -- alias holds a quote.
| CREATE VIEW omni_other.middle AS SELECT id FROM omni_dump.z_base
< ok - -
> exec
| CREATE VIEW omni_dump.a_top AS SELECT id AS `it's` FROM omni_other.middle
< ok - -
> exec
| -- A table Hive translated to external and then made managed: the only managed
| -- table the pinned server can make, and one a dump restores as managed by
| -- stating the ALTER beside the CREATE.
| CREATE TABLE omni_dump.orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_dump.orders SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO omni_dump.orders VALUES (1, 1.5, 'n1')
< ok - -
> exec
| -- A partitioned managed table with one partition and two rows, whose partition
| -- is a name and a location rather than a clause SHOW CREATE TABLE states.
| CREATE TABLE omni_dump.events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_dump.events SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-01')
< ok - -
> exec
| -- A partition that lives outside the table's own directory, which its name does
| -- not state and only the describe of the partition does.
| ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-02') LOCATION 'file:/tmp/omni/dump_events_dt_2024-01-02'
< ok - -
> exec
| -- A partition that overrides its table's file format and SerDe properties, which
| -- ADD PARTITION does not carry and only ALTER TABLE ... PARTITION (...) restores.
| ALTER TABLE omni_dump.events ADD PARTITION (dt='2024-01-03')
< ok - -
> exec
| ALTER TABLE omni_dump.events PARTITION (dt='2024-01-03') SET FILEFORMAT ORC
< ok - -
> exec
| ALTER TABLE omni_dump.events PARTITION (dt='2024-01-03') SET SERDEPROPERTIES ('orc.compress'='ZLIB')
< ok - -
> exec
| INSERT INTO omni_dump.events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')
< ok - -
> exec
| -- An external partitioned table whose partitions live at both default and custom
| -- locations, with one partition unsetting a table SerDe property.
| CREATE EXTERNAL TABLE omni_dump.ext_events (id INT) PARTITIONED BY (dt STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY '|' STORED AS TEXTFILE LOCATION '/tmp/omni/dump_ext_events'
< ok - -
> exec
| ALTER TABLE omni_dump.ext_events ADD PARTITION (dt='2024-01-01')
< ok - -
> exec
| ALTER TABLE omni_dump.ext_events ADD PARTITION (dt='2024-01-02') LOCATION 'file:/tmp/omni/dump_ext_events_dt_2024-01-02'
< ok - -
> exec
| ALTER TABLE omni_dump.ext_events PARTITION (dt='2024-01-02') UNSET SERDEPROPERTIES ('field.delim')
< ok - -
> exec
| -- Tables with named PRIMARY KEY, UNIQUE, and FOREIGN KEY constraints where the
| -- child sorts before the parent: SHOW CREATE TABLE appends ALTER TABLE lines
| -- without a terminating semicolon on CREATE TABLE, and the foreign key must run
| -- after the parent table and its primary key exist.
| CREATE TABLE omni_dump.z_parent (id INT, code STRING, CONSTRAINT pk_z_parent PRIMARY KEY (id) DISABLE NOVALIDATE RELY, CONSTRAINT uq_z_parent UNIQUE (code) DISABLE NOVALIDATE RELY) STORED AS PARQUET
< ok - -
> exec
| CREATE TABLE omni_dump.a_child (id INT, parent_id INT, CONSTRAINT fk_a_child_parent FOREIGN KEY (parent_id) REFERENCES omni_dump.z_parent(id) DISABLE NOVALIDATE RELY) STORED AS PARQUET
< ok - -
> exec
| -- A view with column comments, a view comment, and custom TBLPROPERTIES, all of
| -- which SHOW CREATE TABLE omits.
| CREATE VIEW omni_dump.v_docs (id COMMENT 'the id', label COMMENT 'it\'s a \\label') COMMENT 'documented view' TBLPROPERTIES ('team'='core') AS SELECT id, name FROM omni_dump.ext_users
< ok - -
> exec
| -- A view whose comment SHOW CREATE TABLE leaves out, and the quote in it, which
| -- a replay escapes rather than doubling. Its query names a column with an alias
| -- SHOW CREATE TABLE has to quote: the statement the server states is the one a
| -- replay takes, whatever the query could not be restated from.
| CREATE VIEW omni_dump.z_users COMMENT 'the users\' rows' AS SELECT id AS 1bad, name FROM omni_dump.ext_users
< ok - -
> exec
| -- A materialized view over that view, named so that name order alone would place
| -- it before the view it resolves.
| CREATE MATERIALIZED VIEW omni_dump.a_users_mv DISABLE REWRITE COMMENT 'the users materialized' AS SELECT `1bad`
| FROM omni_dump.z_users
< ok - -
> exec
| -- A partitioned materialized view: the statement states a partition as a table
| -- does, with the column's type, and a materialized view takes it without.
| CREATE MATERIALIZED VIEW omni_dump.p_events_mv DISABLE REWRITE COMMENT 'says PARTITIONED BY (old)' PARTITIONED ON (dt) AS SELECT id, dt FROM omni_dump.events
< ok - -
> query
| SHOW TABLES IN `OMNI_DUMP`
< result
< column "tab_name" "STRING"
< row
: string
| a_child
< row
: string
| a_top
< row
: string
| a_users_mv
< row
: string
| events
< row
: string
| events__
< row
: string
| ext_events
< row
: string
| ext_users
< row
: string
| orders
< row
: string
| p_events_mv
< row
: string
| v_docs
< row
: string
| z_base
< row
: string
| z_parent
< row
: string
| z_users
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`a_child`
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
| parent_id
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
| Location:           
: string
| file:/opt/hive/data/warehouse/omni_dump.db/a_child
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"id\":\"true\",\"parent_id\":\"true\"}}
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
< row
: string
|
: null
: null
< row
: string
| # Constraints
: null
: null
< row
: string
|
: null
: null
< row
: string
| # Foreign Keys
: null
: null
< row
: string
| Table:              
: string
| omni_dump.a_child   
: null
< row
: string
| Constraint Name:    
: string
| fk_a_child_parent   
: null
< row
: string
| Parent Column Name:omni_dump.z_parent.id
: string
| Column Name:parent_id
: string
| Key Sequence:1      
< row
: string
|
: null
: null
> query
| SHOW CREATE TABLE `omni_dump`.`a_child`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_dump`.`a_child`(
< row
: string
|   `id` int, 
< row
: string
|   `parent_id` int)
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
|   'file:/opt/hive/data/warehouse/omni_dump.db/a_child'
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
< row
: string
| ALTER TABLE `omni_dump`.`a_child` ADD CONSTRAINT `fk_a_child_parent` FOREIGN KEY (`parent_id`) REFERENCES `omni_dump`.`z_parent`(`id`) DISABLE NOVALIDATE RELY;
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
| SELECT id AS it's FROM omni_other.middle
: null
< row
: string
| Expanded Query:     
: string
| SELECT `middle`.`id` AS `it's` FROM `omni_other`.`middle`
: null
> query
| SHOW CREATE TABLE `omni_dump`.`a_top`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`a_top` AS SELECT `middle`.`id` AS `it's` FROM `omni_other`.`middle`
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`a_users_mv`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| 1bad
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
| Location:           
: string
| file:/opt/hive/data/warehouse/omni_dump.db/a_users_mv
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"1bad\":\"true\"}}
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
| the users materialized
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
| SELECT 1bad         
: null
< row
: string
|
: string
|                     
: string
| FROM omni_dump.z_users
< row
: string
| Expanded Query:     
: string
| SELECT `z_users`.`1bad`
: null
< row
: string
|
: string
|                     
: string
| FROM `omni_dump`.`z_users`
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
| omni_dump.ext_users
: string
| N/A
: null
> query
| SHOW CREATE TABLE `omni_dump`.`a_users_mv`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_dump`.`a_users_mv`(
< row
: string
|   `1bad` int)
< row
: string
| COMMENT 'the users materialized'
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
|   'file:/opt/hive/data/warehouse/omni_dump.db/a_users_mv'
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
| DESCRIBE FORMATTED `OMNI_DUMP`.`events`
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
| file:/opt/hive/data/warehouse/omni_dump.db/events
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
| 3                   
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
| SHOW CREATE TABLE `omni_dump`.`events`
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
| SHOW PARTITIONS `omni_dump`.`events`
< result
< column "partition" "STRING"
< row
: string
| dt=2024-01-01
< row
: string
| dt=2024-01-02
< row
: string
| dt=2024-01-03
> query
| DESCRIBE FORMATTED `omni_dump`.`events` PARTITION (`dt`='2024-01-01')
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
| # Detailed Partition Information
: null
: null
< row
: string
| Partition Value:    
: string
| [2024-01-01]        
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| Table:              
: string
| events              
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
| Location:           
: string
| file:/opt/hive/data/warehouse/omni_dump.db/events/dt=2024-01-01
: null
< row
: string
| Partition Parameters:
: null
: null
< row
: string
|
: string
| COLUMN_STATS_ACCURATE
: string
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"body\":\"true\",\"id\":\"true\"}}
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
| DESCRIBE FORMATTED `omni_dump`.`events` PARTITION (`dt`='2024-01-02')
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
| # Detailed Partition Information
: null
: null
< row
: string
| Partition Value:    
: string
| [2024-01-02]        
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| Table:              
: string
| events              
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
| Location:           
: string
| file:/tmp/omni/dump_events_dt_2024-01-02
: null
< row
: string
| Partition Parameters:
: null
: null
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
| DESCRIBE FORMATTED `omni_dump`.`events` PARTITION (`dt`='2024-01-03')
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
| # Detailed Partition Information
: null
: null
< row
: string
| Partition Value:    
: string
| [2024-01-03]        
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| Table:              
: string
| events              
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
| Location:           
: string
| file:/opt/hive/data/warehouse/omni_dump.db/events/dt=2024-01-03
: null
< row
: string
| Partition Parameters:
: null
: null
< row
: string
|
: string
| COLUMN_STATS_ACCURATE
: string
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"body\":\"true\",\"id\":\"true\"}}
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
| orc.compress        
: string
| ZLIB                
< row
: string
|
: string
| serialization.format
: string
| 1                   
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`events__`
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
| Location:           
: string
| file:/opt/hive/data/warehouse/omni_dump.db/events__
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
| SHOW CREATE TABLE `omni_dump`.`events__`
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
| DESCRIBE FORMATTED `OMNI_DUMP`.`ext_events`
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
| file:/tmp/omni/dump_ext_events
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
| numPartitions       
: string
| 2                   
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
| |                   
< row
: string
|
: string
| serialization.format
: string
| |                   
> query
| SHOW CREATE TABLE `omni_dump`.`ext_events`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_dump`.`ext_events`(
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
|   'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe' 
< row
: string
| WITH SERDEPROPERTIES ( 
< row
: string
|   'field.delim'='|', 
< row
: string
|   'serialization.format'='|') 
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
|   'file:/tmp/omni/dump_ext_events'
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
| SHOW PARTITIONS `omni_dump`.`ext_events`
< result
< column "partition" "STRING"
< row
: string
| dt=2024-01-01
< row
: string
| dt=2024-01-02
> query
| DESCRIBE FORMATTED `omni_dump`.`ext_events` PARTITION (`dt`='2024-01-01')
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
| # Detailed Partition Information
: null
: null
< row
: string
| Partition Value:    
: string
| [2024-01-01]        
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| Table:              
: string
| ext_events          
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
| Location:           
: string
| file:/tmp/omni/dump_ext_events/dt=2024-01-01
: null
< row
: string
| Partition Parameters:
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
| |                   
< row
: string
|
: string
| serialization.format
: string
| |                   
> query
| DESCRIBE FORMATTED `omni_dump`.`ext_events` PARTITION (`dt`='2024-01-02')
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
| # Detailed Partition Information
: null
: null
< row
: string
| Partition Value:    
: string
| [2024-01-02]        
: null
< row
: string
| Database:           
: string
| omni_dump           
: null
< row
: string
| Table:              
: string
| ext_events          
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
| Location:           
: string
| file:/tmp/omni/dump_ext_events_dt_2024-01-02
: null
< row
: string
| Partition Parameters:
: null
: null
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
| serialization.format
: string
| |                   
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
| DESCRIBE FORMATTED `OMNI_DUMP`.`orders`
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
| file:/opt/hive/data/warehouse/omni_dump.db/orders
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
| SHOW CREATE TABLE `omni_dump`.`orders`
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
| DESCRIBE FORMATTED `OMNI_DUMP`.`p_events_mv`
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
| file:/opt/hive/data/warehouse/omni_dump.db/p_events_mv
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
| {\"BASIC_STATS\":\"true\"}
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
| says PARTITIONED BY (old)
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
| 8                   
< row
: string
|
: string
| totalSize           
: string
| 197                 
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
| SELECT id, dt FROM omni_dump.events
: null
< row
: string
| Expanded Query:     
: string
| SELECT `id`, `dt` FROM (SELECT `events`.`id`, `events`.`dt` FROM `omni_dump`.`events`) `p_events_mv`
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
| omni_dump.events
: string
| N/A
: null
> query
| SHOW CREATE TABLE `omni_dump`.`p_events_mv`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE TABLE `omni_dump`.`p_events_mv`(
< row
: string
|   `id` int)
< row
: string
| COMMENT 'says PARTITIONED BY (old)'
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
|   'file:/opt/hive/data/warehouse/omni_dump.db/p_events_mv'
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
| DESCRIBE FORMATTED `OMNI_DUMP`.`v_docs`
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
| label
: string
| string
: string
| it's a \\label
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
| comment             
: string
| documented view     
< row
: string
|
: string
| team                
: string
| core                
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
| SELECT id, name FROM omni_dump.ext_users
: null
< row
: string
| Expanded Query:     
: string
| SELECT `id` AS `id`, `name` AS `label` FROM (SELECT `ext_users`.`id`, `ext_users`.`name` FROM `omni_dump`.`ext_users`) `v_docs`
: null
> query
| SHOW CREATE TABLE `omni_dump`.`v_docs`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`v_docs` AS SELECT `id` AS `id`, `name` AS `label` FROM (SELECT `ext_users`.`id`, `ext_users`.`name` FROM `omni_dump`.`ext_users`) `v_docs`
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
| SHOW CREATE TABLE `omni_dump`.`z_base`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_base` AS SELECT `a_top`.`id`, '`omni_dump`.`a_top`' AS `tag` FROM `omni_other`.`a_top`
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`z_parent`
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
| code
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
| file:/opt/hive/data/warehouse/omni_dump.db/z_parent
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
| {\"BASIC_STATS\":\"true\",\"COLUMN_STATS\":{\"code\":\"true\",\"id\":\"true\"}}
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
< row
: string
|
: null
: null
< row
: string
| # Constraints
: null
: null
< row
: string
|
: null
: null
< row
: string
| # Primary Key
: null
: null
< row
: string
| Table:              
: string
| omni_dump.z_parent  
: null
< row
: string
| Constraint Name:    
: string
| pk_z_parent         
: null
< row
: string
| Column Name:        
: string
| id                  
: null
< row
: string
|
: null
: null
< row
: string
| # Unique Constraints
: null
: null
< row
: string
| Table:              
: string
| omni_dump.z_parent  
: null
< row
: string
| Constraint Name:    
: string
| uq_z_parent         
: null
< row
: string
| Column Name:code    
: string
| Key Sequence:1      
: null
< row
: string
|
: null
: null
> query
| SHOW CREATE TABLE `omni_dump`.`z_parent`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE EXTERNAL TABLE `omni_dump`.`z_parent`(
< row
: string
|   `id` int, 
< row
: string
|   `code` string)
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
|   'file:/opt/hive/data/warehouse/omni_dump.db/z_parent'
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
< row
: string
| ALTER TABLE `omni_dump`.`z_parent` ADD CONSTRAINT `pk_z_parent` PRIMARY KEY (`id`) DISABLE NOVALIDATE RELY;
< row
: string
| ALTER TABLE `omni_dump`.`z_parent` ADD CONSTRAINT `uq_z_parent` UNIQUE (`code`) DISABLE NOVALIDATE RELY;
> query
| DESCRIBE FORMATTED `OMNI_DUMP`.`z_users`
< result
< column "col_name" "STRING"
< column "data_type" "STRING"
< column "comment" "STRING"
< row
: string
| 1bad
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
| comment             
: string
| the users' rows     
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
| SELECT id AS 1bad, name FROM omni_dump.ext_users
: null
< row
: string
| Expanded Query:     
: string
| SELECT `ext_users`.`id` AS `1bad`, `ext_users`.`name` FROM `omni_dump`.`ext_users`
: null
> query
| SHOW CREATE TABLE `omni_dump`.`z_users`
< result
< column "createtab_stmt" "STRING"
< row
: string
| CREATE VIEW `omni_dump`.`z_users` AS SELECT `ext_users`.`id` AS `1bad`, `ext_users`.`name` FROM `omni_dump`.`ext_users`
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`a_top`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@a_top","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_other@middle","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@a_top]"},{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_other@middle]"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`v_docs`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@v_docs","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@ext_users","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_dump@v_docs]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`z_base`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@z_base","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_other@a_top","tabletype":"VIRTUAL_VIEW","tableParents":"[omni_dump@z_base]"},{"tablename":"omni_other@t","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_other@a_top]"}],"input_partitions":[]}
> query
| EXPLAIN DEPENDENCY SELECT * FROM `OMNI_DUMP`.`z_users`
< result
< column "Explain" "STRING"
< row
: string
| {"input_tables":[{"tablename":"omni_dump@z_users","tabletype":"VIRTUAL_VIEW"},{"tablename":"omni_dump@ext_users","tabletype":"EXTERNAL_TABLE","tableParents":"[omni_dump@z_users]"}],"input_partitions":[]}
