// Bulk.tojson() - Get JSON representation of bulk operation

// Basic tojson
db.users.initializeUnorderedBulkOp().tojson()

// Tojson after inserts
db.users.initializeOrderedBulkOp().insert({ name: "alice" }).tojson()

// Tojson after updates
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).update({ $set: { status: "archived" } }).tojson()

// Tojson for ordered bulk
db.products.initializeOrderedBulkOp().insert({ sku: "A001" }).insert({ sku: "A002" }).tojson()

// Tojson after mixed operations
db.orders.initializeUnorderedBulkOp().insert({ item: "apple" }).find({ qty: 0 }).remove().tojson()
