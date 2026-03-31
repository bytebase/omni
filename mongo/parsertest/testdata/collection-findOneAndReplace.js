// db.collection.findOneAndReplace() - Find and replace a document atomically

// Basic find and replace
db.users.findOneAndReplace({ name: "alice" }, { name: "alice", version: 2 })

// Return new document
db.users.findOneAndReplace({ name: "bob" }, { name: "bob", replaced: true }, { returnDocument: "after" })

// With upsert
db.users.findOneAndReplace({ email: "new@example.com" }, { email: "new@example.com", name: "new" }, { upsert: true })

// With projection
db.users.findOneAndReplace({ name: "charlie" }, { name: "charlie", active: true }, { projection: { name: 1 } })

// With sort
db.users.findOneAndReplace({ status: "old" }, { status: "new", migrated: true }, { sort: { createdAt: 1 } })

// Collection access patterns
db["users"].findOneAndReplace({ id: 1 }, { id: 1, data: "new" })
db.getCollection("users").findOneAndReplace({ key: "abc" }, { key: "abc", value: "xyz" })
