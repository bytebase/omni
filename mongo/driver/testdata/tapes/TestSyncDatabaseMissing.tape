-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestSyncDatabaseMissing default#1
> raw
| {"op":"ListDatabaseNames"}
< raw
| {"value":{"names":["admin","config","local"]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
