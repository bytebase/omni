// The statements of span.tsv, one per line, each run in the database the
// last "// session:" line names, on the fixture the last "// fixture:" line
// names (fixtures/<name>.sql).
// fixture: base
// session: shop

// find and its relatives
db.users.find()
db.users.find({name: "a"})
db.users.find({}, {email: 1, _id: 0})
db.users.find({ssn: "1"}).sort({name: 1}).limit(5)
db.users.find().skip(1).limit(1)
db.users.find({$expr: {$eq: ["$name", "a"]}})
db.users.find({manager: null}, {name: 1}).sort({_id: -1})
db.users.find().explain()
db.users.find().explain("executionStats")
db.users.findOne({_id: 1})
db.users.findOne()
db.users.findOne({}, {email: 1})
db.users.countDocuments({})
db.users.countDocuments({name: "a"})
db.users.estimatedDocumentCount()
db.users.count()
db.users.distinct("name")
db.users.distinct("email", {name: "a"})

// names spelled other ways
db.getCollection("users").find()
db["users"].find()
db.getCollection("Mixed Case").find()
db.getCollection("audit.log").find()
db.getCollection("audit.log").aggregate([{$lookup: {from: "users", localField: "uid", foreignField: "_id", as: "user"}}])
db.empty.find()
db.nosuch.find()
db.nosuch.aggregate([{$lookup: {from: "nosuch2", localField: "a", foreignField: "b", as: "c"}}])

// pipelines that read one collection
db.users.aggregate([])
db.users.aggregate([{$match: {name: "a"}}])
db.users.aggregate([{$project: {name: 1}}])
db.users.aggregate([{$group: {_id: "$name", n: {$sum: 1}}}])
db.users.aggregate([{$sort: {name: 1}}, {$skip: 1}, {$limit: 2}])
db.users.aggregate([{$addFields: {x: 1}}, {$set: {y: 2}}, {$unset: "ssn"}])
db.users.aggregate([{$sortByCount: "$name"}])
db.users.aggregate([{$count: "n"}])
db.users.aggregate([{$replaceRoot: {newRoot: {n: "$name"}}}])
db.users.aggregate([{$setWindowFields: {sortBy: {_id: 1}, output: {r: {$rank: {}}}}}])

// $lookup
db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "orders"}}])
db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "orders"}}, {$unwind: "$orders"}, {$project: {name: 1, "orders.total": 1}}])
db.users.aggregate([{$lookup: {from: "orders", let: {u: "$_id"}, pipeline: [{$match: {$expr: {$eq: ["$uid", "$$u"]}}}], as: "orders"}}])
db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$match: {total: {$gt: 1}}}], as: "big"}}])
db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "uid", pipeline: [{$project: {total: 1}}], as: "orders"}}])
db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "orders"}}, {$lookup: {from: "archive", localField: "_id", foreignField: "_id", as: "old"}}])
db.users.aggregate([{$lookup: {from: "orders", let: {u: "$_id"}, pipeline: [{$match: {$expr: {$eq: ["$uid", "$$u"]}}}, {$lookup: {from: "products", localField: "sku", foreignField: "_id", as: "product"}}], as: "orders"}}])
db.users.aggregate([{$lookup: {from: "orders", pipeline: [{$lookup: {from: "products", pipeline: [{$lookup: {from: "reviews", localField: "_id", foreignField: "sku", as: "r"}}], as: "p"}}], as: "o"}}])
db.users.aggregate([{$lookup: {from: "users", localField: "manager", foreignField: "_id", as: "users"}}])
db.users.aggregate([{$lookup: {from: "active_users", localField: "_id", foreignField: "_id", as: "me"}}])
db.users.aggregate([{$lookup: {from: "user_orders", localField: "_id", foreignField: "_id", as: "mine"}}])
db.orders.aggregate([{$lookup: {from: "users", localField: "uid", foreignField: "_id", as: "orders"}}])
db.orders.aggregate([{$lookup: {from: "nosuch", localField: "uid", foreignField: "_id", as: "x"}}])
db.orders.aggregate([{$lookup: {from: {db: "hr", coll: "employees"}, localField: "uid", foreignField: "_id", as: "x"}}])
db.orders.aggregate([{$match: {total: 5}}, {$lookup: {from: "products", localField: "sku", foreignField: "_id", as: "product"}}, {$replaceRoot: {newRoot: {$first: "$product"}}}])

// $graphLookup
db.users.aggregate([{$graphLookup: {from: "users", startWith: "$manager", connectFromField: "manager", connectToField: "_id", as: "chain"}}])
db.users.aggregate([{$graphLookup: {from: "orders", startWith: "$_id", connectFromField: "uid", connectToField: "uid", as: "g"}}])
db.users.aggregate([{$graphLookup: {from: "active_users", startWith: "$_id", connectFromField: "_id", connectToField: "_id", as: "g"}}])

// $unionWith
db.users.aggregate([{$unionWith: "archive"}])
db.users.aggregate([{$unionWith: {coll: "archive"}}])
db.users.aggregate([{$unionWith: {coll: "orders", pipeline: [{$lookup: {from: "products", localField: "sku", foreignField: "_id", as: "p"}}]}}])
db.users.aggregate([{$unionWith: {coll: "archive", pipeline: [{$unionWith: "orders"}]}}])
db.users.aggregate([{$unionWith: "all_people"}])
db.users.aggregate([{$unionWith: "active_users"}, {$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "o"}}])

// $facet
db.users.aggregate([{$facet: {o: [{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "o"}}], n: [{$count: "n"}]}}])
db.users.aggregate([{$facet: {a: [{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "o"}}], b: [{$lookup: {from: "archive", localField: "_id", foreignField: "_id", as: "x"}}]}}])
db.users.aggregate([{$unionWith: {coll: "orders", pipeline: [{$facet: {p: [{$lookup: {from: "products", localField: "sku", foreignField: "_id", as: "p"}}]}}]}}])
db.users.aggregate([{$facet: {u: [{$unionWith: "archive"}]}}])

// $out and $merge write their target; they read the rest
db.users.aggregate([{$out: "report"}])
db.users.aggregate([{$out: {db: "hr", coll: "report"}}])
db.users.aggregate([{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "o"}}, {$out: "report"}])
db.users.aggregate([{$merge: {into: "report"}}])
db.users.aggregate([{$merge: {into: {db: "hr", coll: "report"}, on: "_id", whenMatched: "replace"}}])
db.users.aggregate([{$unionWith: "archive"}, {$merge: "report"}])

// views
db.active_users.find()
db.active_users.find({name: "a"}, {email: 1})
db.active_users.findOne()
db.active_users.countDocuments({})
db.active_users.distinct("name")
db.user_orders.find()
db.active_user_orders.find({})
db.all_people.find()
db.order_products.find()
db.user_graph.find()
db.order_facets.aggregate([])
db.active_users.aggregate([{$lookup: {from: "user_orders", localField: "_id", foreignField: "_id", as: "o"}}])
db.getCollection("Mixed View").find()

// other databases
db.getSiblingDB("hr").employees.find()
db.getSiblingDB("hr").employees.find({name: "ann"}).limit(1)
db.getSiblingDB("hr").getCollection("employees").find()
db.getSiblingDB("hr").employees.findOne()
db.getSiblingDB("hr").employees.countDocuments({})
db.getSiblingDB("hr").staff.find()
db.getSiblingDB("hr").staff_depts.find()
db.getSiblingDB("hr").employees.aggregate([{$lookup: {from: "departments", localField: "dept", foreignField: "_id", as: "d"}}])
db.getSiblingDB("hr").employees.aggregate([{$unionWith: "staff"}])
db.getSiblingDB("hr").users.find()

// session: hr
db.employees.find()
db.staff_depts.aggregate([{$lookup: {from: "employees", localField: "boss", foreignField: "_id", as: "boss"}}])
db.departments.aggregate([{$graphLookup: {from: "employees", startWith: "$_id", connectFromField: "dept", connectToField: "dept", as: "people"}}])
db.users.find()
