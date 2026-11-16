-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryLimit default#1
> raw
| {"op":"read","statement":"SELECT Id FROM Singers ORDER BY Id LIMIT 2"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}}],"rows":[[{"type":{"code":"INT64"},"value":"1"}],[{"type":{"code":"INT64"},"value":"2"}]]}
> raw
| {"op":"read","statement":"SELECT Id FROM Singers ORDER BY Id LIMIT CAST(3 AS INT64)"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}}],"rows":[[{"type":{"code":"INT64"},"value":"1"}],[{"type":{"code":"INT64"},"value":"2"}],[{"type":{"code":"INT64"},"value":"3"}]]}
> raw
| {"op":"close"}
< raw
| {}
