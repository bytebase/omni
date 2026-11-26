-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQueryLimitCapsRows driver#1
> raw
| {"op":"executeStatement","limit":2,"statement":"SELECT * FROM \"paged\" WHERE owner = 'o1' ORDER BY seq"}
< raw
| {"items":[{"label":{"S":"row-0"},"owner":{"S":"o1"},"seq":{"S":"s0"}},{"label":{"S":"row-1"},"owner":{"S":"o1"},"seq":{"S":"s1"}}],"nextToken":"b21uaS10b2tlbi0x"}
