// Bulk.getOperations() - Get the list of queued operations

// Get operations after inserts
db.users.initializeUnorderedBulkOp().insert({ name: "alice" }).getOperations()

// Get operations after updates
db.users.initializeOrderedBulkOp().find({ status: "inactive" }).update({ $set: { status: "archived" } }).getOperations()

// Get operations after mixed operations
db.users.initializeUnorderedBulkOp().insert({ name: "bob" }).find({ old: true }).remove().getOperations()

// Get operations for ordered bulk
db.products.initializeOrderedBulkOp().insert({ sku: "A001" }).insert({ sku: "A002" }).getOperations()

// Get operations before execute
db.orders.initializeUnorderedBulkOp().find({ status: "pending" }).updateOne({ $set: { status: "processing" } }).getOperations()
