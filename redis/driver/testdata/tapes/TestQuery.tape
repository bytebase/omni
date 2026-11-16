-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestQuery driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["SET","greeting","hello world"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["RPUSH","queue","first","second"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":2}}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Session","session":2}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["GET","greeting"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"hello world"}}]}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["GET","missing"]]}
< raw
| {"commands":[{"error":"redis: nil"}]}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["LRANGE","queue","0","-1"]]}
< raw
| {"commands":[{"value":{"kind":"array","array":[{"kind":"string","string":"first"},{"kind":"string","string":"second"}]}}]}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["DBSIZE"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":2}}]}
> raw
| {"login":{},"op":"CloseSession","session":2}
< raw
| {}
