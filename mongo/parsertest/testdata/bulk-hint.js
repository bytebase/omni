// Bulk.find().hint() - Specify index hint for bulk operations

// Hint with index name
db.users.initializeUnorderedBulkOp().find({ status: "active" }).hint("status_1").update({ $set: { checked: true } })

// Hint with index document
db.users.initializeOrderedBulkOp().find({ email: "alice@example.com" }).hint({ email: 1 }).updateOne({ $set: { verified: true } })

// Hint with compound index
db.orders.initializeUnorderedBulkOp().find({ customerId: "C001", status: "pending" }).hint({ customerId: 1, status: 1 }).update({ $set: { reviewed: true } })

// Hint with remove
db.logs.initializeOrderedBulkOp().find({ level: "debug" }).hint("level_1_timestamp_1").remove()

// Hint and execute
db.products.initializeUnorderedBulkOp().find({ category: "electronics" }).hint({ category: 1 }).updateOne({ $inc: { viewCount: 1 } }).execute()
