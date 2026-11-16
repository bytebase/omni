-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestQuery driver#1
> raw
| {"login":{},"op":"Do","commands":[["SET","greeting","hello world"],["RPUSH","queue","first","second"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}},{"value":{"kind":"int","int":2}}]}
> raw
| {"login":{},"op":"Do","commands":[["GET","greeting"],["GET","missing"],["LRANGE","queue","0","-1"],["DBSIZE"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"hello world"}},{"error":"redis: nil"},{"value":{"kind":"array","array":[{"kind":"string","string":"first"},{"kind":"string","string":"second"}]}},{"value":{"kind":"int","int":2}}]}
