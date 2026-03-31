// cursor.explain() - Return query execution information

// Basic explain (default verbosity)
db.users.find().explain()
db.users.find({}).explain()

// Explain with verbosity levels
db.users.find().explain("queryPlanner")
db.users.find().explain("executionStats")
db.users.find().explain("allPlansExecution")

// With query filter
db.users.find({ status: "active" }).explain()
db.users.find({ age: { $gt: 25 } }).explain("executionStats")

// Explain sorted queries
db.users.find().sort({ name: 1 }).explain()
db.users.find({ status: "active" }).sort({ createdAt: -1 }).explain("executionStats")

// Explain queries with indexes
db.users.find({ email: "test@example.com" }).explain("allPlansExecution")
db.orders.find({ customerId: ObjectId("507f1f77bcf86cd799439011") }).explain()

// Explain complex queries
db.users.find({ $or: [{ status: "active" }, { role: "admin" }] }).explain("executionStats")

// Explain with limit/skip
db.users.find().skip(100).limit(10).explain()
