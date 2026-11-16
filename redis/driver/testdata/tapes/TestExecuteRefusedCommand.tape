-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestExecuteRefusedCommand driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["SET","key","text"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["LPUSH","key","x"]]}
< raw
| {"commands":[{"error":"WRONGTYPE Operation against a key holding the wrong kind of value"}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Session","session":2}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["GET","after"]]}
< raw
| {"commands":[{"error":"redis: nil"}]}
> raw
| {"login":{},"op":"CloseSession","session":2}
< raw
| {}
