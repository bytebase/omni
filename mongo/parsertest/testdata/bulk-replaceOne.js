// Bulk.find().replaceOne() - Queue a replace operation for one matching document

// Basic replaceOne
db.users.initializeUnorderedBulkOp().find({ _id: ObjectId("507f1f77bcf86cd799439011") }).replaceOne({ name: "Alice Updated", email: "alice@example.com", age: 26 })

// ReplaceOne with query
db.users.initializeOrderedBulkOp().find({ email: "bob@example.com" }).replaceOne({ name: "Bob New", email: "bob@example.com", role: "admin" })

// ReplaceOne and execute
db.products.initializeUnorderedBulkOp().find({ sku: "ABC123" }).replaceOne({ sku: "ABC123", name: "Widget", price: 19.99 }).execute()

// Multiple replaceOnes
db.users.initializeOrderedBulkOp().find({ name: "alice" }).replaceOne({ name: "Alice", status: "active" }).find({ name: "bob" }).replaceOne({ name: "Bob", status: "inactive" })

// ReplaceOne with nested document
db.orders.initializeUnorderedBulkOp().find({ orderId: "ORD001" }).replaceOne({ orderId: "ORD001", items: [{ name: "apple", qty: 5 }], total: 25.00 })
