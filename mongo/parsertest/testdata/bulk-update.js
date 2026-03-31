// Bulk.find().update() - Queue an update operation for matching documents

// Update all matching documents
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).update({ $set: { status: "archived" } })

// Update with multiple operators
db.users.initializeOrderedBulkOp().find({ role: "user" }).update({ $set: { verified: true }, $inc: { loginCount: 1 } })

// Update and execute
db.users.initializeUnorderedBulkOp().find({ needsUpdate: true }).update({ $set: { updated: true } }).execute()

// Multiple updates
db.users.initializeOrderedBulkOp().find({ tier: "free" }).update({ $set: { tier: "basic" } }).find({ tier: "basic" }).update({ $set: { tier: "premium" } })

// Update with array operators
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).update({ $push: { tags: "active" } })
