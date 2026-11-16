-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryExplainStopsAtError default#1
> raw
| {"op":"analyzeQuery","statement":"SELECT * FROM Missing"}
< raw
| {"error":{"code":3,"text":"spanner: code = \"InvalidArgument\", desc = \"Table not found: Missing [at 1:15]\\nSELECT * FROM Missing\\n              ^\", requestID = \"omni-request-1\""}}
> raw
| {"op":"close"}
< raw
| {}
