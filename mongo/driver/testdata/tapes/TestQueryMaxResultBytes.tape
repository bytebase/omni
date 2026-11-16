-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQueryMaxResultBytes default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertMany([{_id: 1}, {_id: 2}, {_id: 3}]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$numberInt":"1"},{"$numberInt":"2"},{"$numberInt":"3"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"}},{"_id":{"$numberInt":"2"}},{"_id":{"$numberInt":"3"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"}},{"_id":{"$numberInt":"2"}},{"_id":{"$numberInt":"3"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"}},{"_id":{"$numberInt":"2"}},{"_id":{"$numberInt":"3"}}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
