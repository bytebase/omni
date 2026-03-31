// cursor.returnKey() - Return only the index keys in the result documents

// Return only index keys
db.users.find().returnKey(true)
db.users.find({ status: "active" }).returnKey(true)

// Disable returnKey (return full documents)
db.users.find().returnKey(false)

// With index hint
db.users.find({ status: "active" }).hint({ status: 1 }).returnKey(true)
db.users.find({ email: "test@example.com" }).hint("email_1").returnKey(true)

// With compound index
db.users.find({ name: "alice", age: 25 }).hint({ name: 1, age: 1 }).returnKey(true)

// Chained with other cursor methods
db.users.find().hint({ age: 1 }).returnKey(true).sort({ age: 1 })
db.users.find().returnKey(true).limit(10)

// Useful for debugging index usage
db.orders.find({ customerId: ObjectId("507f1f77bcf86cd799439011") }).hint({ customerId: 1 }).returnKey(true)

// With explain to see index behavior
db.users.find({ status: "active" }).returnKey(true).explain()
