// db.collection.explain() - Return explain plan for query operations (unsupported)

// Basic explain
db.users.explain()
db.orders.explain()

// Explain with verbosity modes
db.users.explain("queryPlanner")
db.users.explain("executionStats")
db.users.explain("allPlansExecution")

// Explain find operations
db.users.explain().find({ status: "active" })
db.users.explain("executionStats").find({ age: { $gt: 25 } })
db.orders.explain("allPlansExecution").find({ customerId: 123, status: "pending" })

// Explain find with sort and limit
db.users.explain().find({ status: "active" }).sort({ createdAt: -1 })
db.users.explain("executionStats").find({}).sort({ age: 1 }).limit(10)

// Explain aggregate operations
db.users.explain().aggregate([{ $match: { status: "active" } }, { $group: { _id: "$role", count: { $sum: 1 } } }])
db.orders.explain("executionStats").aggregate([{ $match: { status: "shipped" } }, { $sort: { createdAt: -1 } }])

// Explain count operations
db.users.explain().count({ status: "active" })
db.orders.explain("executionStats").count({ createdAt: { $gte: ISODate("2024-01-01") } })

// Explain distinct operations
db.users.explain().distinct("status")
db.products.explain("executionStats").distinct("category", { inStock: true })

// Explain update operations
db.users.explain().update({ name: "alice" }, { $set: { age: 26 } })
db.orders.explain("executionStats").updateMany({ status: "pending" }, { $set: { status: "cancelled" } })

// Explain delete operations
db.users.explain().remove({ status: "deleted" })
db.logs.explain("executionStats").deleteMany({ createdAt: { $lt: ISODate("2023-01-01") } })

// Explain findAndModify operations
db.users.explain().findAndModify({ query: { name: "alice" }, update: { $inc: { loginCount: 1 } } })

// Collection access patterns
db["users"].explain()
db["users"].explain("executionStats")
db.getCollection("users").explain()
db.getCollection("orders").explain("allPlansExecution").find({ status: "active" })
db["user-logs"].explain().find({ level: "error" })
