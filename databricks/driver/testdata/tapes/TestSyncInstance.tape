-- pin: driver dbc-3c689f3f-f754.cloud.databricks.com sha256:21c3596967a931ba19b528877357d45cb7c2ea34abc68315d0c4ca7907d9f408 4.2.0 warehouse_id=45847ad2fbd0b316 catalog=workspace

== TestSyncInstance driver#1
> raw
| {"method":"POST","route":"/api/2.0/sql/statements","body":"{\"catalog\":\"workspace\",\"statement\":\"SELECT VERSION()\",\"warehouse_id\":\"45847ad2fbd0b316\"}"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"statement_id\":\"00000000-0000-0000-0000-000000000001\",\"status\":{\"state\":\"SUCCEEDED\"},\"manifest\":{\"format\":\"JSON_ARRAY\",\"schema\":{\"column_count\":1,\"columns\":[{\"name\":\"version()\",\"type_text\":\"STRING\",\"type_name\":\"STRING\",\"nullable\":false}]},\"truncated\":false},\"result\":{\"data_array\":[[\"4.2.0 0000000000000000000000000000000000000000\"]]}}"}
> raw
| {"method":"GET","route":"/api/2.1/unity-catalog/catalogs?max_results=0"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"catalogs\":[{\"name\":\"samples\"},{\"name\":\"system\"},{\"name\":\"workspace\"}]}"}
