// Bulk.find().collation() - Specify collation for bulk operations

// Collation with update
db.users.initializeUnorderedBulkOp().find({ name: "cafe" }).collation({ locale: "fr", strength: 1 }).update({ $set: { found: true } })

// Collation with remove
db.users.initializeOrderedBulkOp().find({ city: "munchen" }).collation({ locale: "de" }).remove()

// Collation with updateOne
db.products.initializeUnorderedBulkOp().find({ name: "WIDGET" }).collation({ locale: "en", strength: 2 }).updateOne({ $set: { normalized: true } })

// Collation with case insensitive search
db.users.initializeOrderedBulkOp().find({ lastName: "smith" }).collation({ locale: "en", strength: 1 }).update({ $set: { processed: true } })

// Multiple operations with collation
db.inventory.initializeUnorderedBulkOp().find({ item: "abc" }).collation({ locale: "en", strength: 2 }).updateOne({ $inc: { qty: 1 } }).execute()
