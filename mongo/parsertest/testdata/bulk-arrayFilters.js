// Bulk.find().arrayFilters() - Specify array filters for bulk update operations

// Basic arrayFilters with updateOne
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).arrayFilters([{ "elem.grade": { $gte: 85 } }]).updateOne({ $set: { "grades.$[elem].passed": true } })

// ArrayFilters with update (all matching)
db.students.initializeOrderedBulkOp().find({ semester: 1 }).arrayFilters([{ "x.grade": { $gte: 90 } }]).update({ $set: { "grades.$[x].honors": true } })

// Multiple array filters
db.orders.initializeUnorderedBulkOp().find({ status: "processing" }).arrayFilters([{ "item.qty": { $gt: 10 } }, { "item.price": { $lt: 50 } }]).updateOne({ $set: { "items.$[item].discount": 0.1 } })

// ArrayFilters with nested arrays
db.inventory.initializeOrderedBulkOp().find({ warehouse: "A" }).arrayFilters([{ "loc.qty": { $gt: 0 } }]).update({ $inc: { "locations.$[loc].qty": -1 } })

// ArrayFilters and execute
db.products.initializeUnorderedBulkOp().find({ category: "electronics" }).arrayFilters([{ "v.rating": { $gte: 4 } }]).updateOne({ $set: { "variants.$[v].featured": true } }).execute()
