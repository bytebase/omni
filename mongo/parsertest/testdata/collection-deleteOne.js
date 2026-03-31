// db.collection.deleteOne() - Delete a single document

// Basic delete
db.users.deleteOne({ name: "alice" })

// Delete by _id
db.users.deleteOne({ _id: ObjectId("507f1f77bcf86cd799439011") })

// Delete with comparison operators
db.users.deleteOne({ createdAt: { $lt: ISODate("2020-01-01") } })

// Delete with options
db.users.deleteOne({ status: "deleted" }, { writeConcern: { w: 1 } })

// Collection access patterns
db["users"].deleteOne({ temp: true })
db.getCollection("users").deleteOne({ expired: true })
