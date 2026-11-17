-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecutePostgresAutoCommit default#1
> raw
| {"op":"partitionedUpdate","statement":"UPDATE Singers SET active = TRUE WHERE id = 1"}
< raw
| {"count":1}
> raw
| {"op":"close"}
< raw
| {}
