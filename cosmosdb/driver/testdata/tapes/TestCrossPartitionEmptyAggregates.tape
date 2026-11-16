-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestCrossPartitionEmptyAggregates emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MIN(c.seq) FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Min\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": MIN(c.seq), \\\\\\\"item2\\\\\\\": {\\\\\\\"min\\\\\\\": MIN(c.seq), \\\\\\\"count\\\\\\\": COUNT(c.seq)}}]\\\\nFROM c\\\\nWHERE (c.seq \u003e 99999999)\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MIN(c.seq) FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MAX(c.seq) FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Max\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": MAX(c.seq), \\\\\\\"item2\\\\\\\": {\\\\\\\"max\\\\\\\": MAX(c.seq), \\\\\\\"count\\\\\\\": COUNT(c.seq)}}]\\\\nFROM c\\\\nWHERE (c.seq \u003e 99999999)\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MAX(c.seq) FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE AVG(c.seq) FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Average\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": {\\\\\\\"sum\\\\\\\": SUM(c.seq), \\\\\\\"count\\\\\\\": COUNT(c.seq)}}]\\\\nFROM c\\\\nWHERE (c.seq \u003e 99999999)\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE AVG(c.seq) FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT MIN(c.seq) AS m FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT MIN(c.seq) AS m FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"items":[{}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT MIN(c.seq) AS lo, MAX(c.seq) AS hi FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT MIN(c.seq) AS lo, MAX(c.seq) AS hi FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"items":[{}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE COUNT(1) FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Count\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": COUNT(1)}]\\\\nFROM c\\\\nWHERE (c.seq \u003e 99999999)\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE COUNT(1) FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"items":[0],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT COUNT(1) AS n FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 Bad Request","statusCode":400,"errorCode":"400 Bad Request","body":"{\"code\":\"BadRequest\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Cross partition query only supports 'VALUE \u003cAggregateFunc\u003e' for aggregates.\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT COUNT(1) AS n FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"items":[{"n":0}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE SUM(c.seq) FROM c WHERE c.seq \u003e 99999999"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Sum\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": SUM(c.seq)}]\\\\nFROM c\\\\nWHERE (c.seq \u003e 99999999)\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE SUM(c.seq) FROM c WHERE c.seq \u003e 99999999","engine":true}
< raw
| {"items":[0],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MIN(c.seq) FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Min\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": MIN(c.seq), \\\\\\\"item2\\\\\\\": {\\\\\\\"min\\\\\\\": MIN(c.seq), \\\\\\\"count\\\\\\\": COUNT(c.seq)}}]\\\\nFROM c\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MIN(c.seq) FROM c","engine":true}
< raw
| {"items":[0],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MAX(c.seq) FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/paged/docs","status":"400 BadRequest","statusCode":400,"errorCode":"400 BadRequest","body":"{\"code\":\"BadRequest\",\"message\":\"The provided cross partition query can not be directly served by the gateway. This is a first chance (internal) exception that all newer clients will know how to handle gracefully. This exception is traced, but unless you see it bubble up as an exception (which only happens on older SDK clients), then you can safely ignore this message.\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Microsoft.Azure.Documents.Common/2.14.0\",\"additionalErrorInfo\":\"{\\\"partitionedQueryExecutionInfoVersion\\\":2,\\\"queryInfo\\\":{\\\"distinctType\\\":\\\"None\\\",\\\"top\\\":null,\\\"offset\\\":null,\\\"limit\\\":null,\\\"orderBy\\\":[],\\\"orderByExpressions\\\":[],\\\"groupByExpressions\\\":[],\\\"groupByAliases\\\":[],\\\"aggregates\\\":[\\\"Max\\\"],\\\"groupByAliasToAggregateType\\\":{},\\\"rewrittenQuery\\\":\\\"SELECT VALUE [{\\\\\\\"item\\\\\\\": MAX(c.seq), \\\\\\\"item2\\\\\\\": {\\\\\\\"max\\\\\\\": MAX(c.seq), \\\\\\\"count\\\\\\\": COUNT(c.seq)}}]\\\\nFROM c\\\",\\\"hasSelectValue\\\":true,\\\"dCountInfo\\\":null,\\\"hasNonStreamingOrderBy\\\":false},\\\"queryRanges\\\":[{\\\"min\\\":\\\"\\\",\\\"max\\\":\\\"FF\\\",\\\"isMinInclusive\\\":true,\\\"isMaxInclusive\\\":false}],\\\"hybridSearchQueryInfo\\\":null}\"}"}}}
> raw
| {"op":"Page","database":"cosmosfixture","container":"paged","statement":"SELECT VALUE MAX(c.seq) FROM c","engine":true}
< raw
| {"items":[2129],"pages":1}
