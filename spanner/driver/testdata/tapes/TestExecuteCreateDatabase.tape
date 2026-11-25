-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteCreateDatabase driver#1
> raw
| {"op":"createDatabase","instancePath":"projects/omni/instances/createdb","createStatement":"CREATE DATABASE `created`","extraStatements":["CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id)"]}
< raw
| {}
> raw
| {"op":"close"}
< raw
| {}

== TestExecuteCreateDatabase driver#2
> raw
| {"op":"getDatabaseDdl","database":"projects/omni/instances/createdb/databases/created"}
< raw
| {"ddl":["CREATE TABLE Written (\n  Id INT64 NOT NULL,\n) PRIMARY KEY(Id)"]}
> raw
| {"op":"close"}
< raw
| {}
