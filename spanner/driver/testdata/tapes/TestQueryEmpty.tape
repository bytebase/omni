-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryEmpty driver#1
> raw
| {"op":"read","statement":"SELECT Id, FirstName FROM Singers WHERE Id \u003c 0"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}},{"name":"FirstName","type":{"code":"STRING"}}]}
> raw
| {"op":"close"}
< raw
| {}
