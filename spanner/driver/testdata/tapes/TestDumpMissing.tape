-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestDumpMissing default#1
> raw
| {"op":"getDatabaseDdl","database":"projects/omni/instances/sync/databases/missing"}
< raw
| {"error":{"code":5,"text":"rpc error: code = NotFound desc = Database not found: projects/omni/instances/sync/databases/missing\nerror details: name = ResourceInfo type = type.googleapis.com/google.spanner.admin.database.v1.Database resourcename = projects/omni/instances/sync/databases/missing owner =  desc = Database does not exist."}}
> raw
| {"op":"close"}
< raw
| {}
