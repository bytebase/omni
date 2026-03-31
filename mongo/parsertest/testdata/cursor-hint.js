// cursor.hint() - Force the query optimizer to use a specific index

// Hint with index name
db.users.find({ status: "active" }).hint("status_1")
db.users.find({ email: "test@example.com" }).hint("email_1")

// Hint with index specification document
db.users.find({ status: "active" }).hint({ status: 1 })
db.users.find({ name: "alice", age: 25 }).hint({ name: 1, age: 1 })

// Force collection scan (no index)
db.users.find({ status: "active" }).hint({ $natural: 1 })
db.users.find().hint({ $natural: -1 })

// Chained with other cursor methods
db.users.find({ status: "active" }).hint("status_1").sort({ createdAt: -1 })
db.users.find({ status: "active" }).hint({ status: 1 }).limit(10)
db.orders.find({ customerId: ObjectId("507f1f77bcf86cd799439011") }).hint("customerId_1").explain()

// Compound index hint
db.products.find({ category: "electronics", price: { $lt: 500 } }).hint({ category: 1, price: 1 })

// With projection
db.users.find({ status: "active" }).hint("status_1").projection({ name: 1, email: 1 })
