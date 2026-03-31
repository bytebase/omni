// Bulk.toString() - Get string representation of bulk operation

// Basic toString
db.users.initializeUnorderedBulkOp().toString()

// ToString after inserts
db.users.initializeOrderedBulkOp().insert({ name: "alice" }).toString()

// ToString after updates
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).update({ $set: { status: "archived" } }).toString()

// ToString for ordered bulk
db.products.initializeOrderedBulkOp().insert({ sku: "A001" }).insert({ sku: "A002" }).toString()

// ToString after mixed operations
db.orders.initializeUnorderedBulkOp().insert({ item: "apple" }).find({ qty: 0 }).remove().toString()
