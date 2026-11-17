-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecutePostgres default#1
> raw
| {"op":"batchUpdate","statements":["UPDATE Singers SET rating = 5.0 WHERE id = 1;","UPDATE Singers SET rating = 4.0 WHERE id = 2"]}
< raw
| {"counts":[1,1]}
> raw
| {"op":"updateDatabaseDdl","database":"projects/omni/instances/pgexecute/databases/pgexecute","statements":["CREATE TABLE written (id bigint NOT NULL, name text, PRIMARY KEY (id));","DROP TABLE written"]}
< raw
| {}
> raw
| {"op":"close"}
< raw
| {}
