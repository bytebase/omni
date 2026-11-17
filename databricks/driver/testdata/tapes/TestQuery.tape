-- pin: driver dbc-3c689f3f-f754.cloud.databricks.com sha256:21c3596967a931ba19b528877357d45cb7c2ea34abc68315d0c4ca7907d9f408 4.2.0 warehouse_id=45847ad2fbd0b316 catalog=workspace

== TestQuery driver#1
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT 1 AS n\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000001\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"n\",\"type_text\":\"INT\",\"type_name\":\"INT\",\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"1\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"row_limit\":1,\"statement\":\"SELECT 1 AS n UNION ALL SELECT 2\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000002\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"n\",\"type_text\":\"INT\",\"type_name\":\"INT\",\"nullable\":false}]},\"truncated\":true},\"result\":{\"data_array\":[[\"1\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT CAST(NULL AS DECIMAL(12,2)) AS d, CAST(NULL AS INT) AS n, '' AS s\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000003\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":3,\"columns\":[{\"name\":\"d\",\"type_text\":\"DECIMAL(12,2)\",\"type_name\":\"DECIMAL\",\"nullable\":false},{\"name\":\"n\",\"type_text\":\"INT\",\"type_name\":\"INT\",\"position\":1,\"nullable\":false},{\"name\":\"s\",\"type_text\":\"STRING\",\"type_name\":\"STRING\",\"position\":2,\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"\",\"\",\"\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT 1 AS n UNION ALL SELECT 2 UNION ALL SELECT 3\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000004\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"n\",\"type_text\":\"INT\",\"type_name\":\"INT\",\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"1\"],[\"2\"],[\"3\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT 1 AS n WHERE false\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000005\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"n\",\"type_text\":\"INT\",\"type_name\":\"INT\",\"nullable\":false}]},\"truncated\":false},\"result\":{}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"EXPLAIN FORMATTED SELECT 1 AS n\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000006\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"plan\",\"type_text\":\"STRING\",\"type_name\":\"STRING\",\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"== Physical Plan ==\\nLocalTableScan (1)\\n\\n\\n(1) LocalTableScan\\nOutput [1]: [n#0]\\nArguments: [n#0]\\n\\n\\n== Photon Explanation ==\\nPhoton does not fully support the query because:\\n\\t\\tUnsupported node: LocalTableScan [n#0].\\n\\nReference node:\\n\\tLocalTableScan [n#0]\\n\"]]}}"}
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT * FROM omni_driver_missing_table\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000007\",\"status\":{\"state\":\"FAILED\",\"error\":{\"error_code\":\"BAD_REQUEST\",\"message\":\"[TABLE_OR_VIEW_NOT_FOUND] The table or view `omni_driver_missing_table` cannot be found. Verify the spelling and correctness of the schema and catalog.\\nSearch path: [`system`.`session`, `system`.`builtin`, `system`.`ai`, `workspace`.`default`].\\nIf you did not qualify the name with a schema, verify the current_schema() output, or qualify the name with the correct schema and catalog.\\nTo tolerate the error on drop use DROP VIEW IF EXISTS or DROP TABLE IF EXISTS. SQLSTATE: 42P01; line 1 pos 14\"},\"sql_state\":\"42P01\"}}"}
