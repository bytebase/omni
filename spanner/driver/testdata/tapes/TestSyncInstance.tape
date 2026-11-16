-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestSyncInstance default#1
> raw
| {"op":"listDatabases"}
< raw
| {"databases":[{"Name":"projects/omni/instances/syncinstance/databases/fixture","Dialect":"GOOGLE_STANDARD_SQL"},{"Name":"projects/omni/instances/syncinstance/databases/pgdialect","Dialect":"POSTGRESQL"}]}
> raw
| {"op":"close"}
< raw
| {}
