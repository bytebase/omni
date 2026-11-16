-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestQueryAllNodes driver#1
> raw
| {"login":{},"cluster":true,"op":"Do","commands":[["SET","allnodes:key","allnodes-value"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"cluster":true,"op":"DoEachShard","commands":[["GET","allnodes:key"],["DBSIZE"]]}
< raw
| {"shards":[[{"value":{"kind":"string","string":"allnodes-value"}},{"value":{"kind":"int","int":1}}]]}
