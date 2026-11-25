-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryExplain driver#1
> raw
| {"op":"analyzeQuery","statement":"SELECT * FROM Singers"}
< raw
| {"error":{"code":13,"text":"spanner: code = \"Internal\", desc = \"query plan unavailable\""}}
> raw
| {"op":"close"}
< raw
| {}
