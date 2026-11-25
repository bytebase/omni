// The starting databases of span.tsv: mongosh statements, one per line, each
// run in the database the last "// database:" line names. A view is made with
// the create command, since gomongo, which runs these lines as Bytebase runs a
// statement, has no db.createView.
// database: hr
db.employees.insertOne({_id: 1, name: "ann", salary: 100, dept: 1, boss: null})
db.departments.insertOne({_id: 1, title: "eng"})
db.runCommand({create: "staff", viewOn: "employees", pipeline: [{$project: {salary: 0}}]})
db.runCommand({create: "staff_depts", viewOn: "staff", pipeline: [{$lookup: {from: "departments", localField: "dept", foreignField: "_id", as: "dept"}}]})
// database: shop
db.users.insertOne({_id: 1, name: "a", email: "a@x", ssn: "1", manager: null})
db.orders.insertOne({_id: 1, uid: 1, total: 5, sku: "s1"})
db.products.insertOne({_id: "s1", name: "p"})
db.reviews.insertOne({_id: 1, sku: "s1", body: "ok"})
db.archive.insertOne({_id: 2, name: "old"})
db.createCollection("empty")
db.getCollection("Mixed Case").insertOne({_id: 1})
db.getCollection("audit.log").insertOne({_id: 1, uid: 1})
db.runCommand({create: "active_users", viewOn: "users", pipeline: [{$match: {name: {$ne: null}}}]})
db.runCommand({create: "user_orders", viewOn: "users", pipeline: [{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "orders"}}]})
db.runCommand({create: "active_user_orders", viewOn: "active_users", pipeline: [{$lookup: {from: "orders", localField: "_id", foreignField: "uid", as: "orders"}}]})
db.runCommand({create: "all_people", viewOn: "users", pipeline: [{$unionWith: "archive"}]})
db.runCommand({create: "order_products", viewOn: "orders", pipeline: [{$lookup: {from: "products", let: {s: "$sku"}, pipeline: [{$match: {$expr: {$eq: ["$_id", "$$s"]}}}, {$lookup: {from: "reviews", localField: "_id", foreignField: "sku", as: "reviews"}}], as: "product"}}]})
db.runCommand({create: "user_graph", viewOn: "users", pipeline: [{$graphLookup: {from: "user_orders", startWith: "$_id", connectFromField: "_id", connectToField: "_id", as: "graph"}}]})
db.runCommand({create: "order_facets", viewOn: "orders", pipeline: [{$facet: {buyers: [{$lookup: {from: "active_users", localField: "uid", foreignField: "_id", as: "buyer"}}]}}]})
db.runCommand({create: "Mixed View", viewOn: "Mixed Case", pipeline: []})
