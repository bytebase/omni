-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestCrossPartitionAggregates emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE COUNT(1) FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Count\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": COUNT(1)}]\\\\nFROM c\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE COUNT(1) FROM c","engine":true}
< raw
| {"items":[40],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT COUNT(1) AS total FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT COUNT(1) AS total FROM c","engine":true}
< raw
| {"items":[{"total":40}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE SUM(c.population) FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Sum\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": SUM(c.population)}]\\\\nFROM c\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE SUM(c.population) FROM c","engine":true}
< raw
| {"items":[12000000],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE AVG(c.population) FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Average\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": {\\\\\\\"sum\\\\\\\": SUM(c.population), \\\\\\\"count\\\\\\\": COUNT(c.population)}}]\\\\nFROM c\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE AVG(c.population) FROM c","engine":true}
< raw
| {"items":[300000],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT MIN(c.population) AS lo, MAX(c.population) AS hi FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT MIN(c.population) AS lo, MAX(c.population) AS hi FROM c","engine":true}
< raw
| {"items":[{"hi":500000,"lo":100000}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT COUNT(1) AS total, SUM(c.population) AS pop FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/cities/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT COUNT(1) AS total, SUM(c.population) AS pop FROM c","engine":true}
< raw
| {"items":[{"pop":12000000,"total":40}],"pages":1}
