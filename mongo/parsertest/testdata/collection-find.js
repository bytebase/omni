// db.collection.find() - Query documents in a collection

// Basic find
db.users.find()
db.users.find({})

// Find with simple filter
db.users.find({ name: "alice" })
db.users.find({ age: 25 })
db.users.find({ status: "active" })

// Find with comparison operators
db.users.find({ age: { $gt: 25 } })
db.users.find({ age: { $gte: 18 } })
db.users.find({ age: { $lt: 65 } })
db.users.find({ age: { $lte: 30 } })
db.users.find({ age: { $ne: 0 } })
db.users.find({ age: { $gte: 18, $lt: 65 } })

// Find with array operators
db.users.find({ status: { $in: ["active", "pending"] } })
db.users.find({ status: { $nin: ["deleted", "banned"] } })
db.users.find({ tags: { $all: ["mongodb", "database"] } })
db.users.find({ scores: { $elemMatch: { $gt: 80, $lt: 90 } } })

// Find with logical operators
db.users.find({ $or: [{ name: "alice" }, { name: "bob" }] })
db.users.find({ $and: [{ age: { $gt: 18 } }, { status: "active" }] })
db.users.find({ age: { $not: { $lt: 18 } } })
db.users.find({ $nor: [{ status: "deleted" }, { status: "banned" }] })

// Find with nested documents
db.users.find({ "address.city": "New York" })
db.users.find({ "profile.settings.theme": "dark" })
db.users.find({ profile: { name: "test", active: true } })

// Find with array fields
db.users.find({ tags: "mongodb" })
db.users.find({ "tags.0": "primary" })

// Find with existence check
db.users.find({ email: { $exists: true } })
db.users.find({ deletedAt: { $exists: false } })

// Find with type check
db.users.find({ age: { $type: "number" } })
db.users.find({ name: { $type: "string" } })

// Find with regex
db.users.find({ name: /^alice/i })
db.users.find({ email: { $regex: /.*@example\.com$/ } })

// Find with helper functions
db.users.find({ _id: ObjectId("507f1f77bcf86cd799439011") })
db.users.find({ createdAt: { $gt: ISODate("2024-01-01T00:00:00Z") } })
db.users.find({ sessionId: UUID("550e8400-e29b-41d4-a716-446655440000") })

// Find with cursor modifiers
db.users.find().sort({ age: -1 })
db.users.find().limit(10)
db.users.find().skip(5)
db.users.find().count()
db.users.find().projection({ name: 1, age: 1 })
db.users.find().project({ name: 1, email: 1, _id: 0 })

// Find with chained cursor modifiers
db.users.find().sort({ age: -1 }).limit(10)
db.users.find().sort({ createdAt: -1 }).skip(20).limit(10)
db.users.find({ status: "active" }).sort({ name: 1 }).limit(100).skip(0)
db.users.find({ status: "active" }).sort({ name: 1 }).limit(100).count()

// Find with collection access patterns
db["users"].find({ name: "alice" })
db['user-logs'].find({ level: "error" })
db.getCollection("users").find({ active: true })
db.getCollection("my.collection").find()
