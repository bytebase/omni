// Bulk mixed operations - Combining insert, update, and remove in one bulk

// Mixed insert and update
db.users.initializeOrderedBulkOp().insert({ name: "alice", status: "new" }).find({ status: "inactive" }).update({ $set: { status: "archived" } })

// Mixed insert, update, and remove
db.products.initializeUnorderedBulkOp().insert({ sku: "NEW001", name: "New Product" }).find({ discontinued: true }).remove().find({ price: { $lt: 10 } }).updateOne({ $set: { onSale: true } })

// Complex mixed operations with execute
db.orders.initializeOrderedBulkOp().insert({ orderId: "ORD100", status: "pending" }).find({ status: "shipped" }).updateOne({ $set: { status: "delivered" } }).find({ status: "cancelled" }).remove().execute()

// Multiple operations of each type
db.inventory.initializeUnorderedBulkOp().insert({ item: "apple", qty: 100 }).insert({ item: "banana", qty: 50 }).find({ qty: 0 }).remove().find({ qty: { $lt: 10 } }).update({ $set: { lowStock: true } }).execute()

// Mixed with upsert
db.users.initializeOrderedBulkOp().find({ email: "new@example.com" }).upsert().updateOne({ $set: { name: "New User", email: "new@example.com" } }).insert({ email: "another@example.com", name: "Another User" }).find({ deleted: true }).remove()

// Mixed ordered operations (execute in sequence)
db.logs.initializeOrderedBulkOp().insert({ level: "info", message: "Process started" }).find({ level: "debug" }).remove().insert({ level: "info", message: "Process completed" }).execute({ w: "majority" })
