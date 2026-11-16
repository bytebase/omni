-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQueryDoubleNotation default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertOne({_id: 1, created_at: 1779696815227});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.aggregate([{$project: {_id: 0, big: {$divide: [\"$created_at\", 1]}, small: {$divide: [0.0000001, 1]}}}])"}
< raw
| {"value":{"operation":"OpAggregate","values":[{"big":{"$numberDouble":"1.779696815227E+12"},"small":{"$numberDouble":"1E-07"}}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
