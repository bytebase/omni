-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestExecute default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertMany([{_id: 1, name: \"alice\"}]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$numberInt":"1"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.createIndex({name: 1}, {unique: true});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["name_1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.updateOne({_id: 1}, {$set: {name: \"alice\", age: 30}});"}
< raw
| {"value":{"operation":"OpUpdateOne","values":[{"acknowledged":true,"matchedCount":{"$numberLong":"1"},"modifiedCount":{"$numberLong":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"},"name":"alice","age":{"$numberInt":"30"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertOne({_id: 2});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"2"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertOne({_id: 1});"}
< raw
| {"error":{"message":"insertOne failed: write exception: write errors: [E11000 duplicate key error collection: omni_driver.users index: _id_ dup key: { _id: 1 }]"}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find()"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"1"},"name":"alice","age":{"$numberInt":"30"}},{"_id":{"$numberInt":"2"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertOne({_id: 9, name: \"dave\", modifiedCount: 100, deletedCount: 7});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"9"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.find({_id: 9});"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$numberInt":"9"},"name":"dave","modifiedCount":{"$numberInt":"100"},"deletedCount":{"$numberInt":"7"}}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
