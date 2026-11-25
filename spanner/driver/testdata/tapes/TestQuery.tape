-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestQuery driver#1
> raw
| {"op":"read","statement":"SELECT Id, FirstName, Active, Rating, Photo, Born, UpdatedAt, Meta, Scores FROM Singers ORDER BY Id"}
< raw
| {"columns":[{"name":"Id","type":{"code":"INT64"}},{"name":"FirstName","type":{"code":"STRING"}},{"name":"Active","type":{"code":"BOOL"}},{"name":"Rating","type":{"code":"FLOAT64"}},{"name":"Photo","type":{"code":"BYTES"}},{"name":"Born","type":{"code":"DATE"}},{"name":"UpdatedAt","type":{"code":"TIMESTAMP"}},{"name":"Meta","type":{"code":"JSON"}},{"name":"Scores","type":{"code":"ARRAY","arrayElementType":{"code":"INT64"}}}],"rows":[[{"type":{"code":"INT64"},"value":"1"},{"type":{"code":"STRING"},"value":"Ann"},{"type":{"code":"BOOL"},"value":true},{"type":{"code":"FLOAT64"},"value":4.5},{"type":{"code":"BYTES"},"value":"AQID"},{"type":{"code":"DATE"},"value":"2024-01-15"},{"type":{"code":"TIMESTAMP"},"value":"2024-01-15T10:30:45.123456Z"},{"type":{"code":"JSON"},"value":"{\"a\":1,\"b\":[2,3]}"},{"type":{"code":"ARRAY","arrayElementType":{"code":"INT64"}},"value":["1","2","3"]}],[{"type":{"code":"INT64"},"value":"2"},{"type":{"code":"STRING"},"value":"Bob"},{"type":{"code":"BOOL"},"value":null},{"type":{"code":"FLOAT64"},"value":null},{"type":{"code":"BYTES"},"value":null},{"type":{"code":"DATE"},"value":null},{"type":{"code":"TIMESTAMP"},"value":null},{"type":{"code":"JSON"},"value":null},{"type":{"code":"ARRAY","arrayElementType":{"code":"INT64"}},"value":null}],[{"type":{"code":"INT64"},"value":"3"},{"type":{"code":"STRING"},"value":"Cid"},{"type":{"code":"BOOL"},"value":false},{"type":{"code":"FLOAT64"},"value":3.25},{"type":{"code":"BYTES"},"value":null},{"type":{"code":"DATE"},"value":null},{"type":{"code":"TIMESTAMP"},"value":null},{"type":{"code":"JSON"},"value":null},{"type":{"code":"ARRAY","arrayElementType":{"code":"INT64"}},"value":null}]]}
> raw
| {"op":"close"}
< raw
| {}
