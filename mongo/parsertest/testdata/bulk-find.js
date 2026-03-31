// Bulk.find() - Specify a query for bulk update/remove operations

// Basic find for bulk update
db.users.initializeUnorderedBulkOp().find({ status: "inactive" })

// Find with complex query
db.users.initializeOrderedBulkOp().find({ age: { $gt: 18 } })

// Find with multiple conditions
db.users.initializeUnorderedBulkOp().find({ status: "active", role: "user" })

// Find with ObjectId
db.users.initializeOrderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") })

// Find with nested document query
db.orders.initializeUnorderedBulkOp().find({ "address.city": "NYC" })
