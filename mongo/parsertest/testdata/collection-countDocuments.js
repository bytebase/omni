// db.collection.countDocuments() - Count documents matching a filter

// Count all documents
db.users.countDocuments()
db.users.countDocuments({})

// Count with simple filter
db.users.countDocuments({ status: "active" })
db.users.countDocuments({ verified: true })
db.users.countDocuments({ role: "admin" })

// Count with comparison operators
db.users.countDocuments({ age: { $gt: 18 } })
db.users.countDocuments({ age: { $gte: 21, $lt: 65 } })
db.users.countDocuments({ loginCount: { $gte: 10 } })

// Count with logical operators
db.users.countDocuments({ $or: [{ status: "active" }, { status: "pending" }] })
db.users.countDocuments({ $and: [{ verified: true }, { active: true }] })

// Count with array operators
db.users.countDocuments({ tags: { $in: ["premium", "enterprise"] } })
db.users.countDocuments({ roles: { $all: ["read", "write"] } })

// Count with existence check
db.users.countDocuments({ email: { $exists: true } })
db.users.countDocuments({ deletedAt: { $exists: false } })

// Count with nested documents
db.users.countDocuments({ "address.country": "USA" })
db.users.countDocuments({ "profile.verified": true })

// Count with helper functions
db.users.countDocuments({ createdAt: { $gt: ISODate("2024-01-01") } })
db.users.countDocuments({ lastLogin: { $lt: ISODate("2024-06-01") } })

// Count with options (options passed to driver)
db.users.countDocuments({ status: "active" }, { skip: 10, limit: 100 })
db.users.countDocuments({ verified: true }, { maxTimeMS: 5000 })
db.users.countDocuments({}, { hint: { status: 1 } })

// Count with collection access patterns
db["users"].countDocuments({ active: true })
db['audit-logs'].countDocuments({ level: "error" })
db.getCollection("orders").countDocuments({ status: "completed" })
