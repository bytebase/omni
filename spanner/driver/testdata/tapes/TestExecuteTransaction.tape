-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteTransaction default#1
> raw
| {"op":"batchUpdate","statements":["UPDATE Singers SET Rating = 5.0 WHERE Id = 1","UPDATE Singers SET Rating = 4.0 WHERE Id = 2"]}
< raw
| {"counts":[1,1]}
> raw
| {"op":"close"}
< raw
| {}
