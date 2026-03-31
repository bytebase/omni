// Bulk.find().arrayFilters().updateOne() - UpdateOne with array filters

// Basic updateOne with arrayFilters
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).arrayFilters([{ "elem.status": "pending" }]).updateOne({ $set: { "tasks.$[elem].status": "completed" } })

// UpdateOne with multiple array filters
db.orders.initializeOrderedBulkOp().find({ customerId: "C001" }).arrayFilters([{ "item.qty": { $gt: 5 } }, { "item.price": { $gte: 100 } }]).updateOne({ $set: { "items.$[item].priority": "high" } })

// UpdateOne with nested array filter
db.inventory.initializeUnorderedBulkOp().find({ warehouse: "main" }).arrayFilters([{ "bin.count": { $lt: 10 } }]).updateOne({ $set: { "bins.$[bin].needsRestock": true } })

// Chained operations with arrayFilters
db.students.initializeOrderedBulkOp().find({ class: "Math101" }).arrayFilters([{ "g.score": { $gte: 90 } }]).updateOne({ $set: { "grades.$[g].letter": "A" } }).find({ class: "Math101" }).arrayFilters([{ "g.score": { $lt: 60 } }]).updateOne({ $set: { "grades.$[g].letter": "F" } })

// UpdateOne with arrayFilters and execute
db.products.initializeUnorderedBulkOp().find({ category: "clothing" }).arrayFilters([{ "size.stock": 0 }]).updateOne({ $set: { "sizes.$[size].outOfStock": true } }).execute()
