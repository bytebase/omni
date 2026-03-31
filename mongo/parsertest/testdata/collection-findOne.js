// db.collection.findOne() - Find a single document in a collection

// Basic findOne
db.users.findOne()
db.users.findOne({})

// FindOne with simple filter
db.users.findOne({ name: "alice" })
db.users.findOne({ _id: ObjectId("507f1f77bcf86cd799439011") })
db.users.findOne({ email: "alice@example.com" })

// FindOne with comparison operators
db.users.findOne({ age: { $gt: 18 } })
db.users.findOne({ age: { $gte: 21, $lt: 65 } })

// FindOne with logical operators
db.users.findOne({ $or: [{ name: "alice" }, { email: "alice@example.com" }] })
db.users.findOne({ $and: [{ status: "active" }, { verified: true }] })

// FindOne with nested documents
db.users.findOne({ "address.city": "New York" })
db.users.findOne({ "profile.settings.notifications": true })

// FindOne with array operators
db.users.findOne({ tags: { $in: ["admin", "moderator"] } })
db.users.findOne({ roles: { $all: ["read", "write"] } })

// FindOne with helper functions
db.users.findOne({ createdAt: { $gt: ISODate("2024-01-01") } })
db.sessions.findOne({ sessionId: UUID("550e8400-e29b-41d4-a716-446655440000") })
db.orders.findOne({ total: NumberDecimal("99.99") })

// FindOne with collection access patterns
db["users"].findOne({ name: "bob" })
db['audit-logs'].findOne({ action: "login" })
db.getCollection("users").findOne({ active: true })
