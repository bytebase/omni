// db.collection.count() - Count documents (deprecated - use countDocuments or estimatedDocumentCount)

// Basic count
db.users.count()
db.orders.count()

// Count with filter
db.users.count({ status: "active" })
db.users.count({ age: { $gte: 18 } })
db.orders.count({ status: "pending", customerId: ObjectId("507f1f77bcf86cd799439011") })

// Count with complex filter
db.users.count({ $or: [{ status: "active" }, { role: "admin" }] })
db.products.count({ price: { $gte: 10, $lte: 100 }, inStock: true })

// Count with options
db.users.count({ status: "active" }, { limit: 1000 })
db.orders.count({ createdAt: { $gte: ISODate("2024-01-01") } }, { skip: 10 })
db.products.count({}, { limit: 100, skip: 50 })

// Count with hint
db.users.count({ email: { $exists: true } }, { hint: { email: 1 } })
db.orders.count({ customerId: 123 }, { hint: "customerId_1" })

// Count with maxTimeMS
db.largeCollection.count({ type: "log" }, { maxTimeMS: 5000 })

// Count with readConcern
db.users.count({ active: true }, { readConcern: { level: "majority" } })

// Collection access patterns
db["users"].count()
db["users"].count({ status: "active" })
db.getCollection("users").count()
db.getCollection("users").count({ role: "admin" })
db["user-events"].count({ type: "login" })
