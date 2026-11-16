-- pin: default amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQueryPages default#1
> raw
| {"op":"executeStatement","limit":2,"statement":"SELECT * FROM \"paged\" WHERE owner = 'o1' ORDER BY seq"}
< raw
| {"items":[{"label":{"S":"row-0"},"owner":{"S":"o1"},"seq":{"S":"s0"}},{"label":{"S":"row-1"},"owner":{"S":"o1"},"seq":{"S":"s1"}}],"nextToken":"b21uaS10b2tlbi0x"}
> raw
| {"op":"executeStatement","limit":2,"statement":"SELECT * FROM \"paged\" WHERE owner = 'o1' ORDER BY seq","nextToken":"b21uaS10b2tlbi0x"}
< raw
| {"items":[{"label":{"S":"row-2"},"owner":{"S":"o1"},"seq":{"S":"s2"}},{"label":{"S":"row-3"},"owner":{"S":"o1"},"seq":{"S":"s3"}}],"nextToken":"b21uaS10b2tlbi0y"}
> raw
| {"op":"executeStatement","limit":2,"statement":"SELECT * FROM \"paged\" WHERE owner = 'o1' ORDER BY seq","nextToken":"b21uaS10b2tlbi0y"}
< raw
| {"items":[{"label":{"S":"row-4"},"owner":{"S":"o1"},"seq":{"S":"s4"}}]}
