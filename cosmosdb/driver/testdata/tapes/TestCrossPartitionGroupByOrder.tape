-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestCrossPartitionGroupByOrder emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket","engine":true}
< raw
| {"items":[{"bucket":"big-1","n":1050},{"bucket":"big-2","n":1050},{"bucket":"small-1","n":5},{"bucket":"small-2","n":5},{"bucket":"small-3","n":5},{"bucket":"small-4","n":5},{"bucket":"small-5","n":5},{"bucket":"small-6","n":5}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket DESC"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket DESC","engine":true}
< raw
| {"items":[{"bucket":"small-6","n":5},{"bucket":"small-5","n":5},{"bucket":"small-4","n":5},{"bucket":"small-3","n":5},{"bucket":"small-2","n":5},{"bucket":"small-1","n":5},{"bucket":"big-2","n":1050},{"bucket":"big-1","n":1050}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket","engine":true}
< raw
| {"items":[{"bucket":"big-1","n":1050},{"bucket":"big-2","n":1050},{"bucket":"small-1","n":5},{"bucket":"small-2","n":5},{"bucket":"small-3","n":5},{"bucket":"small-4","n":5},{"bucket":"small-5","n":5},{"bucket":"small-6","n":5}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket","engine":true}
< raw
| {"items":[{"bucket":"big-1","n":1050},{"bucket":"big-2","n":1050},{"bucket":"small-1","n":5},{"bucket":"small-2","n":5},{"bucket":"small-3","n":5},{"bucket":"small-4","n":5},{"bucket":"small-5","n":5},{"bucket":"small-6","n":5}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket","engine":true}
< raw
| {"items":[{"bucket":"big-1","n":1050},{"bucket":"big-2","n":1050},{"bucket":"small-1","n":5},{"bucket":"small-2","n":5},{"bucket":"small-3","n":5},{"bucket":"small-4","n":5},{"bucket":"small-5","n":5},{"bucket":"small-6","n":5}],"pages":1}
