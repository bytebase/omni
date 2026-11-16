-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQueryLimitDoesNotCapACount default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertMany([\n\t\t{_id: 1}, {_id: 2}, {_id: 3}, {_id: 4}, {_id: 5}\n\t]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$numberInt":"1"},{"$numberInt":"2"},{"$numberInt":"3"},{"$numberInt":"4"},{"$numberInt":"5"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.countDocuments({})"}
< raw
| {"value":{"operation":"OpCountDocuments","values":[{"$numberLong":"5"}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.estimatedDocumentCount()"}
< raw
| {"value":{"operation":"OpEstimatedDocumentCount","values":[{"$numberLong":"5"}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
