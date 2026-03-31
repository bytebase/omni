// Bulk.find().removeOne() - Queue a remove operation for one matching document

// Remove one matching document
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).removeOne()

// Remove one with complex query
db.users.initializeOrderedBulkOp().find({ age: { $lt: 18 } }).removeOne()

// Remove one and execute
db.users.initializeUnorderedBulkOp().find({ deleted: true }).removeOne().execute()

// Multiple remove ones
db.users.initializeOrderedBulkOp().find({ role: "guest" }).removeOne().find({ expired: true }).removeOne()

// Remove one with ObjectId
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).removeOne()
