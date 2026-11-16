-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQueryError default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertOne({_id: 1});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.runCommand({nosuchcommand: 1})"}
< raw
| {"error":{"code":59,"name":"CommandNotFound","message":"no such command: 'nosuchcommand'"}}
> raw
| {"op":"Disconnect"}
< raw
| {}
