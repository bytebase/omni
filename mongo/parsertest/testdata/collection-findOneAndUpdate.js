// db.collection.findOneAndUpdate() - Find and update a document atomically

// Basic find and update
db.users.findOneAndUpdate({ name: "alice" }, { $set: { age: 26 } })

// Return new document
db.users.findOneAndUpdate({ name: "bob" }, { $inc: { score: 10 } }, { returnDocument: "after" })

// With upsert
db.users.findOneAndUpdate({ email: "new@example.com" }, { $set: { name: "new" } }, { upsert: true, returnDocument: "after" })

// With projection
db.users.findOneAndUpdate({ name: "charlie" }, { $set: { active: true } }, { projection: { name: 1, active: 1 } })

// With sort
db.users.findOneAndUpdate({ status: "pending" }, { $set: { status: "processing" } }, { sort: { createdAt: 1 } })

// Collection access patterns
db["users"].findOneAndUpdate({ id: 1 }, { $set: { processed: true } })
db.getCollection("users").findOneAndUpdate({ key: "abc" }, { $set: { value: "xyz" } })
