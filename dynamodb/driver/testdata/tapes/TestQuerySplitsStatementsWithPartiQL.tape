-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQuerySplitsStatementsWithPartiQL driver#1
> raw
| {"op":"executeStatement","statement":"\nSELECT * FROM \"people\" WHERE id = 'u3'"}
< raw
| {"items":[{"email":{"S":"u3@example.com"},"id":{"S":"u3"},"late":{"S":"last"}}]}
