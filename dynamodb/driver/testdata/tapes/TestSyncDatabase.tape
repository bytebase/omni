-- pin: default amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestSyncDatabase default#1
> raw
| {"op":"listTables","limit":100}
< raw
| {"tables":["orders","paged","people","writes"]}
> raw
| {"op":"describeTable","table":"orders"}
< raw
| {"table":{"name":"orders","status":"ACTIVE","itemCount":2,"keySchema":[{"attributeName":"pk","keyType":"HASH"}]}}
> raw
| {"op":"describeTable","table":"paged"}
< raw
| {"table":{"name":"paged","status":"ACTIVE","itemCount":5,"keySchema":[{"attributeName":"owner","keyType":"HASH"},{"attributeName":"seq","keyType":"RANGE"}],"secondaryIndexes":[{"name":"by-label","global":true,"keySchema":[{"attributeName":"label","keyType":"HASH"}]},{"name":"by-owner-label","keySchema":[{"attributeName":"owner","keyType":"HASH"},{"attributeName":"label","keyType":"RANGE"}]}]}}
> raw
| {"op":"describeTable","table":"people"}
< raw
| {"table":{"name":"people","status":"ACTIVE","itemCount":3,"keySchema":[{"attributeName":"id","keyType":"HASH"},{"attributeName":"email","keyType":"RANGE"}]}}
> raw
| {"op":"describeTable","table":"writes"}
< raw
| {"table":{"name":"writes","status":"ACTIVE","keySchema":[{"attributeName":"id","keyType":"HASH"}]}}
