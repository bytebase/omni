-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestExecuteDelimiterAndDDLFlows default#1
> exec
| USE `omni_qe_ddl_flows`
< ok 0 0
> begin
< ok
> exec
| 		CREATE TABLE delimiter_one (id INT)
< ok 0 0
> exec
|
| 		CREATE TABLE delimiter_two (id INT)
< ok 0 0
> commit
< ok
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'omni_qe_ddl_flows' AND table_name = 'delimiter_one'
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 1
> query
| SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'omni_qe_ddl_flows' AND table_name = 'delimiter_two'
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 1
> begin
< ok
> exec
|
| 		USE omni_qe_ddl_flows;
< ok 0 0
> exec
|
| 		CREATE TABLE omni_qe_ddl_flows.execute_create_index_in_transaction (id INT);
< ok 0 0
> exec
|
| 		BEGIN;
< ok 0 0
> exec
|
| 		CREATE INDEX idx_execute_create_index_in_transaction ON omni_qe_ddl_flows.execute_create_index_in_transaction(id);
< ok 0 0
> exec
|
| 		COMMIT;
< ok 0 0
> commit
< ok
> query
|
| 		SELECT COUNT(*) FROM information_schema.tidb_indexes
| 		WHERE table_schema = 'omni_qe_ddl_flows' AND table_name = 'execute_create_index_in_transaction' AND key_name = 'idx_execute_create_index_in_transaction'
| 	
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 1
> begin
< ok
> exec
|
| 		USE omni_qe_ddl_flows;
< ok 0 0
> exec
|
| 		CREATE TABLE omni_qe_ddl_flows.prepare_statement_flow (id INT);
< ok 0 0
> exec
|
| 		SET @sql := 'CREATE INDEX idx_prepare_statement_flow ON omni_qe_ddl_flows.prepare_statement_flow(id)';
< ok 0 0
> exec
|
| 		PREPARE stmt FROM @sql;
< ok 0 0
> exec
|
| 		EXECUTE stmt;
< ok 0 0
> exec
|
| 		DEALLOCATE PREPARE stmt;
< ok 0 0
> commit
< ok
> query
|
| 		SELECT COUNT(*) FROM information_schema.tidb_indexes
| 		WHERE table_schema = 'omni_qe_ddl_flows' AND table_name = 'prepare_statement_flow' AND key_name = 'idx_prepare_statement_flow'
| 	
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 1
