-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQueryMaxResultBytesStopsTheRead driver#1
> raw
| {"op":"executeStatement","limit":3,"statement":"SELECT * FROM \"paged\" WHERE owner = 'o1' AND label \u003c\u003e 'row-0' ORDER BY seq"}
< raw
| {"items":[{"label":{"S":"row-1"},"owner":{"S":"o1"},"seq":{"S":"s1"}},{"label":{"S":"row-2"},"owner":{"S":"o1"},"seq":{"S":"s2"}}],"nextToken":"b21uaS10b2tlbi0x"}
