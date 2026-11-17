-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteCreateDatabaseRefreshesTheDialect default#1
> raw
| {"op":"createDatabase","instancePath":"projects/omni/instances/pgfresh","createStatement":"CREATE DATABASE \"pgfresh\"","dialect":"POSTGRESQL"}
< raw
| {}
> raw
| {"op":"read","statement":"SELECT $$a;b$$ AS x"}
< raw
| {"columns":[{"name":"x","type":{"code":"STRING"}}],"rows":[[{"type":{"code":"STRING"},"value":"a;b"}]]}
> raw
| {"op":"close"}
< raw
| {}
