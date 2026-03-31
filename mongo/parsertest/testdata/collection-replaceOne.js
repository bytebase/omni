// db.collection.replaceOne() - Replace a single document

// Basic replace
db.users.replaceOne({ name: "alice" }, { name: "alice", age: 30, status: "active" })

// Replace with upsert
db.users.replaceOne({ email: "new@example.com" }, { email: "new@example.com", name: "New User" }, { upsert: true })

// Replace by _id
db.users.replaceOne({ _id: ObjectId("507f1f77bcf86cd799439011") }, { _id: ObjectId("507f1f77bcf86cd799439011"), name: "replaced", version: 2 })

// Replace with options
db.users.replaceOne({ name: "bob" }, { name: "bob", migrated: true }, { writeConcern: { w: 1 } })

// Collection access patterns
db["users"].replaceOne({ id: 1 }, { id: 1, data: "new" })
db.getCollection("users").replaceOne({ key: "abc" }, { key: "abc", value: "xyz" })
