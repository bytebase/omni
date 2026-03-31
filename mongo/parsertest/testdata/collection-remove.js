// db.collection.remove() - Remove documents (deprecated - use deleteOne or deleteMany)

// Remove all documents
db.users.remove({})
db.logs.remove({})

// Remove with filter
db.users.remove({ status: "deleted" })
db.users.remove({ name: "alice" })
db.orders.remove({ status: "cancelled" })

// Remove with comparison operators
db.users.remove({ age: { $lt: 18 } })
db.logs.remove({ createdAt: { $lt: ISODate("2023-01-01") } })
db.sessions.remove({ expiresAt: { $lte: Date() } })

// Remove with complex filter
db.users.remove({ $or: [{ status: "deleted" }, { status: "banned" }] })
db.orders.remove({ status: "cancelled", createdAt: { $lt: ISODate("2024-01-01") } })

// Remove with array operators
db.users.remove({ tags: { $in: ["spam", "bot"] } })
db.products.remove({ categories: { $all: ["discontinued", "clearance"] } })

// Remove single document (justOne option)
db.users.remove({ status: "inactive" }, true)
db.users.remove({ status: "pending" }, { justOne: true })
db.queue.remove({ processed: true }, { justOne: true })

// Remove with write concern
db.users.remove({ temp: true }, { writeConcern: { w: "majority" } })
db.orders.remove({ status: "test" }, { writeConcern: { w: 1, j: true } })

// Remove with collation
db.products.remove({ name: "cafe" }, { collation: { locale: "fr", strength: 1 } })

// Remove by _id
db.users.remove({ _id: ObjectId("507f1f77bcf86cd799439011") })
db.users.remove({ _id: 1 })

// Remove with nested field
db.users.remove({ "profile.deleted": true })
db.orders.remove({ "shipping.status": "failed" })

// Remove with let option
db.users.remove({ $expr: { $eq: ["$_id", "$$targetId"] } }, { let: { targetId: ObjectId("507f1f77bcf86cd799439011") } })

// Collection access patterns
db["users"].remove({ status: "deleted" })
db["users"].remove({}, true)
db.getCollection("users").remove({ expired: true })
db.getCollection("temp-data").remove({})
db["user-sessions"].remove({ expiresAt: { $lt: Date() } })
