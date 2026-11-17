-- pin: driver dbc-3c689f3f-f754.cloud.databricks.com sha256:21c3596967a931ba19b528877357d45cb7c2ea34abc68315d0c4ca7907d9f408 4.2.0 warehouse_id=45847ad2fbd0b316 catalog=workspace

== TestExecute driver#1
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"CREATE SCHEMA `workspace`.`omni_driver_probe`\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000001\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{},\"truncated\":false},\"result\":{}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"CREATE TABLE `workspace`.`omni_driver_probe`.`probe` (id INT)\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000002\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{},\"truncated\":false},\"result\":{}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"INSERT INTO `workspace`.`omni_driver_probe`.`probe` VALUES (1), (2)\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000003\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":2,\"columns\":[{\"name\":\"num_affected_rows\",\"type_text\":\"BIGINT\",\"type_name\":\"LONG\",\"nullable\":false},{\"name\":\"num_inserted_rows\",\"type_text\":\"BIGINT\",\"type_name\":\"LONG\",\"position\":1,\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"2\",\"2\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT COUNT(*) FROM `workspace`.`omni_driver_probe`.`probe`\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000004\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"count(1)\",\"type_text\":\"BIGINT\",\"type_name\":\"LONG\",\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"2\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"DROP TABLE IF EXISTS `workspace`.`omni_driver_probe`.`probe`\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000005\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{},\"truncated\":false},\"result\":{}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"DROP SCHEMA IF EXISTS `workspace`.`omni_driver_probe`\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000006\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{},\"truncated\":false},\"result\":{}}"}
