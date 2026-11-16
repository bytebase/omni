-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryGrant default#1
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/querygrant/databases/querygrant","statements":["CREATE ROLE reviewer"]}
< raw
| {}
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/querygrant/databases/querygrant","statements":["GRANT SELECT ON TABLE Singers TO ROLE reviewer"]}
< raw
| {}
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/querygrant/databases/querygrant","statements":["REVOKE SELECT ON TABLE Singers FROM ROLE reviewer"]}
< raw
| {}
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/querygrant/databases/querygrant","statements":["GRANT SELECT ON TABLE Singers TO ROLE reviewer"]}
< raw
| {}
> raw
| {"op":"close"}
< raw
| {}
