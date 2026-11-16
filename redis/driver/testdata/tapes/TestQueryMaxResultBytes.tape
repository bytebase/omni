-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestQueryMaxResultBytes driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["RPUSH","list","aaaaaaaaaa","bbbbbbbbbb","cccccccccc","dddddddddd","eeeeeeeeee","ffffffffff","gggggggggg","hhhhhhhhhh"]]}
< raw
| {"commands":[{"value":{"kind":"int","int":8}}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Session","session":2}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":2,"commands":[["LRANGE","list","0","-1"]]}
< raw
| {"commands":[{"value":{"kind":"array","array":[{"kind":"string","string":"aaaaaaaaaa"},{"kind":"string","string":"bbbbbbbbbb"},{"kind":"string","string":"cccccccccc"},{"kind":"string","string":"dddddddddd"},{"kind":"string","string":"eeeeeeeeee"},{"kind":"string","string":"ffffffffff"},{"kind":"string","string":"gggggggggg"},{"kind":"string","string":"hhhhhhhhhh"}]}}]}
> raw
| {"login":{},"op":"CloseSession","session":2}
< raw
| {}
