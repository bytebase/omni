// Bulk.find().upsert() - Set upsert flag for update operations

// Upsert with updateOne
db.users.initializeUnorderedBulkOp().find({ email: "alice@example.com" }).upsert().updateOne({ $set: { name: "Alice", email: "alice@example.com" } })

// Upsert with update
db.users.initializeOrderedBulkOp().find({ sku: "ABC123" }).upsert().update({ $set: { sku: "ABC123", qty: 100 } })

// Upsert with replaceOne
db.users.initializeUnorderedBulkOp().find({ email: "bob@example.com" }).upsert().replaceOne({ name: "Bob", email: "bob@example.com", role: "user" })

// Upsert and execute
db.products.initializeOrderedBulkOp().find({ sku: "XYZ789" }).upsert().updateOne({ $set: { price: 29.99 } }).execute()

// Multiple upserts
db.inventory.initializeUnorderedBulkOp().find({ item: "apple" }).upsert().updateOne({ $set: { qty: 50 } }).find({ item: "banana" }).upsert().updateOne({ $set: { qty: 30 } })
