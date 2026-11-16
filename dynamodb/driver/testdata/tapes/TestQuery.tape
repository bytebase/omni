-- pin: default amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQuery default#1
> raw
| {"op":"executeStatement","statement":"SELECT * FROM \"people\""}
< raw
| {"items":[{"active":{"BOOL":true},"age":{"N":"33"},"email":{"S":"u1@example.com"},"flags":{"BS":["/g==","/w=="]},"id":{"S":"u1"},"list":{"L":[{"N":"7"},{"S":"s"},{"B":"/g=="}]},"meta":{"M":{"k":{"S":"v"},"raw":{"B":"/w=="}}},"missing":{"NULL":true},"name":{"S":"Ann"},"payload":{"B":"/wD+"},"scores":{"NS":["1","2"]},"tags":{"SS":["a","b"]}},{"age":{"N":"44"},"email":{"S":"u2@example.com"},"id":{"S":"u2"},"name":{"S":"Bob"}},{"email":{"S":"u3@example.com"},"id":{"S":"u3"},"late":{"S":"last"}}]}
