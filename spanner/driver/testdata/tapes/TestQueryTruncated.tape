-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQueryTruncated default#1
> raw
| {"op":"read","statement":"SELECT Id, FirstName FROM Singers ORDER BY Id"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}},{"name":"FirstName","type":{"code":"STRING"}}],"rows":[[{"type":{"code":"INT64"},"value":"1"},{"type":{"code":"STRING"},"value":"Ann"}],[{"type":{"code":"INT64"},"value":"2"},{"type":{"code":"STRING"},"value":"Bob"}],[{"type":{"code":"INT64"},"value":"3"},{"type":{"code":"STRING"},"value":"Cid"}]]}
> raw
| {"op":"close"}
< raw
| {}
