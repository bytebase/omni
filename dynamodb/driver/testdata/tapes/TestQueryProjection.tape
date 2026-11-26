-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQueryProjection driver#1
> raw
| {"op":"executeStatement","statement":"SELECT name, age FROM \"people\" WHERE id = 'u2'"}
< raw
| {"items":[{"age":{"N":"44"},"name":{"S":"Bob"}}]}
