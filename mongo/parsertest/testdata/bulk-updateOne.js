// Bulk.find().updateOne() - Queue an update operation for one matching document

// Update one matching document
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).updateOne({ $set: { status: "active" } })

// Update one with complex query
db.users.initializeOrderedBulkOp().find({ email: "alice@example.com" }).updateOne({ $set: { verified: true } })

// Update one and execute
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).updateOne({ $inc: { count: 1 } }).execute()

// Multiple update ones
db.users.initializeOrderedBulkOp().find({ name: "alice" }).updateOne({ $set: { active: true } }).find({ name: "bob" }).updateOne({ $set: { active: false } })

// Update one with unset
db.users.initializeUnorderedBulkOp().find({ tempField: { $exists: true } }).updateOne({ $unset: { tempField: "" } })
