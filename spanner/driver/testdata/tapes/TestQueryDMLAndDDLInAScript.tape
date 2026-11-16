-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryDMLAndDDLInAScript default#1
> raw
| {"op":"update","statement":"UPDATE Singers SET Rating = 5.0 WHERE Id = 1"}
< raw
| {"count":1}
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/queryscript/databases/queryscript","statements":["CREATE TABLE Written (Id INT64 NOT NULL) PRIMARY KEY (Id)"]}
< raw
| {}
> raw
| {"op":"update","statement":"INSERT INTO Written (Id) VALUES (1), (2)"}
< raw
| {"count":2}
> raw
| {"op":"read","statement":"SELECT Id FROM Written ORDER BY Id"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}}],"rows":[[{"type":{"code":"INT64"},"value":"1"}],[{"type":{"code":"INT64"},"value":"2"}]]}
> raw
| {"op":"close"}
< raw
| {}
