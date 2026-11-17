-- pin: driver dbc-3c689f3f-f754.cloud.databricks.com sha256:21c3596967a931ba19b528877357d45cb7c2ea34abc68315d0c4ca7907d9f408 4.2.0 warehouse_id=45847ad2fbd0b316 catalog=workspace

== TestPing driver#1
> raw
| {"method":"GET","route":"/api/2.0/workspace/list?path=%2F"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"objects\":[{\"path\":\"/Repos\",\"object_type\":\"DIRECTORY\"},{\"path\":\"/Shared\",\"object_type\":\"DIRECTORY\"},{\"path\":\"/Users\",\"object_type\":\"DIRECTORY\"}]}"}
