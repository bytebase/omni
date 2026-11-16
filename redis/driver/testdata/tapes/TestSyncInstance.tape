-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestSyncInstance driver#1
> raw
| {"login":{},"op":"Info","section":"server"}
< raw
| {"info":"redis_version:7.4.12\r"}
> raw
| {"login":{},"op":"ConfigGet","parameter":"databases"}
< raw
| {"config":{"databases":"1"}}
