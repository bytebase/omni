// Bulk.find().hint().updateOne() - UpdateOne with index hint

// UpdateOne with string hint
db.users.initializeUnorderedBulkOp().find({ email: "alice@example.com" }).hint("email_1").updateOne({ $set: { lastLogin: ISODate() } })

// UpdateOne with document hint
db.orders.initializeOrderedBulkOp().find({ orderId: "ORD001" }).hint({ orderId: 1 }).updateOne({ $set: { status: "shipped" } })

// Chained find with hint and updateOne
db.products.initializeUnorderedBulkOp().find({ sku: "ABC123" }).hint({ sku: 1, warehouse: 1 }).updateOne({ $inc: { stock: -1 } })

// Multiple updateOnes with hints
db.inventory.initializeOrderedBulkOp().find({ item: "apple" }).hint("item_1").updateOne({ $set: { fresh: true } }).find({ item: "banana" }).hint("item_1").updateOne({ $set: { fresh: true } })

// UpdateOne with hint and execute
db.logs.initializeUnorderedBulkOp().find({ level: "error", timestamp: { $lt: ISODate() } }).hint("level_1_timestamp_-1").updateOne({ $set: { archived: true } }).execute()
