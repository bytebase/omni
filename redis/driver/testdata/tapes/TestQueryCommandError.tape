-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestQueryCommandError driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["SET","counter","text"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Session","session":2}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["INCR","counter"]]}
< raw
| {"commands":[{"error":"ERR value is not an integer or out of range"}]}
> raw
| {"login":{},"op":"CloseSession","session":2}
< raw
| {}
