-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestStatementTextKeepsComments driver#1
> raw
| {"op":"read","statement":"SELECT 1 AS a -- a note"}
< raw
| {"columns":[{"name":"a","type":{"code":"INT64"}}],"rows":[[{"type":{"code":"INT64"},"value":"1"}]]}
> raw
| {"op":"close"}
< raw
| {}
