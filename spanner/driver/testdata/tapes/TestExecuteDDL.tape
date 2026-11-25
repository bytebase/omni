-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteDDL driver#1
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/execute/databases/ddl","statements":["CREATE TABLE Written (Id INT64 NOT NULL, Name STRING(64)) PRIMARY KEY (Id)","DROP TABLE Written"]}
< raw
| {}
> raw
| {"op":"close"}
< raw
| {}
