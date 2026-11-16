-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestSyncDatabase driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["SET","greeting","hello world"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["HSET","user:1","name","Alice"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":1}}]}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["RPUSH","queue","first","second"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":2}}]}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["INCR","counter"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":1}}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}
