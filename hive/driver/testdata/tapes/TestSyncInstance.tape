-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "_c0" "STRING"
< row
: string
| 4.0.1 r3af4517eb8cfd9407ad34ed78a0b48b57dfaa264
> query
| SHOW DATABASES
< result
< column "database_name" "STRING"
< row
: string
| default
< row
: string
| sync_inst_db
