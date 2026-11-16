-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteRollback default#1
> raw
| {"op":"batchUpdate","statements":["UPDATE Singers SET Rating = 1.0 WHERE Id = 1","UPDATE Missing SET a = 1"]}
< raw
| {"error":{"code":3,"text":"spanner: code = \"InvalidArgument\", desc = \"Table not found: Missing [at 1:8]\\nUPDATE Missing SET a = 1\\n       ^\""}}
> raw
| {"op":"read","statement":"SELECT Rating FROM Singers WHERE Id = 1"}
< raw
| {"columns":[{"name":"Rating","type":{"code":"FLOAT64"}}],"rows":[[{"type":{"code":"FLOAT64"},"value":4.5}]]}
> raw
| {"op":"close"}
< raw
| {}
