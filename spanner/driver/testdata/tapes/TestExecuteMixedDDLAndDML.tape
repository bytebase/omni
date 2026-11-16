-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteMixedDDLAndDML default#1
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/execute/databases/mixed","statements":["CREATE TABLE Mixed (Id INT64 NOT NULL) PRIMARY KEY (Id)","UPDATE Singers SET Rating = 3.0 WHERE Id = 1"]}
< raw
| {"error":{"code":3,"text":"rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: UPDATE Singers SET Rating = 3.0 WHERE Id = 1 : Syntax error on line 1, column 1: Encountered 'UPDATE' while parsing: ddl_statement"}}
> raw
| {"op":"close"}
< raw
| {}
