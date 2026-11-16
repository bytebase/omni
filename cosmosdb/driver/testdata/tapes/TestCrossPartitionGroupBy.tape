-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestCrossPartitionGroupBy emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT c.country, COUNT(1) AS cities, SUM(c.population) AS pop FROM c GROUP BY c.country"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT c.country, COUNT(1) AS cities, SUM(c.population) AS pop FROM c GROUP BY c.country","engine":true}
< raw
| {"items":[{"cities":5,"country":"AE","pop":1500000},{"cities":5,"country":"BR","pop":1500000},{"cities":5,"country":"CN","pop":1500000},{"cities":5,"country":"DE","pop":1500000},{"cities":5,"country":"FR","pop":1500000},{"cities":5,"country":"IN","pop":1500000},{"cities":5,"country":"JP","pop":1500000},{"cities":5,"country":"US","pop":1500000}],"pages":1}
