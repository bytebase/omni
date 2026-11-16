-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestExecuteRefusedCommand driver#1
> raw
| {"login":{},"op":"Do","commands":[["SET","key","text"],["LPUSH","key","x"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}},{"error":"WRONGTYPE Operation against a key holding the wrong kind of value"}]}
