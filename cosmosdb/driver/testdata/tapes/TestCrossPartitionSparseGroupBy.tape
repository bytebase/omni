-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestCrossPartitionSparseGroupBy emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.v AS k, COUNT(1) AS n FROM c GROUP BY c.v"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.v AS k, COUNT(1) AS n FROM c GROUP BY c.v","engine":true}
< raw
| {"items":[{"n":2119},{"k":null,"n":1},{"k":false,"n":1},{"k":true,"n":1},{"k":1,"n":1},{"k":20,"n":1},{"k":"a","n":1},{"k":"b","n":1},{"k":[1,2],"n":1},{"k":[2],"n":1},{"k":{"a":1},"n":1},{"k":{"b":1},"n":1}],"pages":1}
