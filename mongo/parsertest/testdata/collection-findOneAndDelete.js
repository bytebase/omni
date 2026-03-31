// db.collection.findOneAndDelete() - Find and delete a document atomically

// Basic find and delete
db.users.findOneAndDelete({ name: "alice" })

// With projection
db.users.findOneAndDelete({ status: "deleted" }, { projection: { name: 1, email: 1 } })

// With sort (delete oldest)
db.users.findOneAndDelete({ status: "expired" }, { sort: { createdAt: 1 } })

// Delete by _id
db.users.findOneAndDelete({ _id: ObjectId("507f1f77bcf86cd799439011") })

// With maxTimeMS
db.users.findOneAndDelete({ temp: true }, { maxTimeMS: 5000 })

// Collection access patterns
db["users"].findOneAndDelete({ processed: true })
db.getCollection("queue").findOneAndDelete({ status: "pending" }, { sort: { priority: -1 } })
