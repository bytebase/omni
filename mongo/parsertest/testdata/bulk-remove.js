// Bulk.find().remove() - Queue a remove operation for matching documents

// Remove all matching documents
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).remove()

// Remove with complex query
db.users.initializeOrderedBulkOp().find({ age: { $lt: 18 } }).remove()

// Remove and execute
db.users.initializeUnorderedBulkOp().find({ deleted: true }).remove().execute()

// Multiple removes
db.users.initializeOrderedBulkOp().find({ role: "guest" }).remove().find({ expired: true }).remove()

// Remove with date condition
db.logs.initializeUnorderedBulkOp().find({ createdAt: { $lt: ISODate("2023-01-01") } }).remove()
