-- pin: driver redis:7.4.12-alpine sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe 7.4.12

== TestSyncInstanceFromKeyspace driver#1
> raw
| {"login":{},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{},"op":"Do","session":1,"commands":[["ACL","SETUSER","noconfig","on","\u003enoconfig-pw","~*","+@all","-config"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{},"op":"CloseSession","session":1}
< raw
| {}

== TestSyncInstanceFromKeyspace driver#2
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Session","session":1}
< raw
| {}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Do","session":1,"commands":[["SET","keyspace:key","value"]]}
< raw
| {"commands":[{"value":{"kind":"string","string":"OK"}}]}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"CloseSession","session":1}
< raw
| {}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Session","session":2}
< raw
| {}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Do","session":2,"commands":[["CONFIG","GET","databases"]]}
< raw
| {"commands":[{"error":"NOPERM User noconfig has no permissions to run the 'config|get' command"}]}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"CloseSession","session":2}
< raw
| {}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Info","section":"server"}
< raw
| {"info":"redis_version:7.4.12\r"}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"ConfigGet","parameter":"databases"}
< raw
| {"error":"NOPERM User noconfig has no permissions to run the 'config|get' command"}
> raw
| {"login":{"user":"noconfig","password":"noconfig-pw"},"op":"Info","section":"keyspace"}
< raw
| {"info":"db0:keys=1,expires=0,avg_ttl=0,subexpiry=0\r"}
