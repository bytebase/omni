// Bulk array insert pattern - Multiple inserts using bulk operations

// Multiple inserts in sequence
db.users.initializeUnorderedBulkOp().insert({ name: "alice" }).insert({ name: "bob" }).insert({ name: "charlie" })

// Insert array of documents pattern
db.products.initializeOrderedBulkOp().insert({ sku: "A001", name: "Widget" }).insert({ sku: "A002", name: "Gadget" }).insert({ sku: "A003", name: "Gizmo" }).execute()

// Insert with varying document structures
db.logs.initializeUnorderedBulkOp().insert({ level: "info", message: "Started" }).insert({ level: "warn", message: "Slow query", duration: 500 }).insert({ level: "error", message: "Failed", stack: "..." })

// Bulk insert with IDs
db.users.initializeOrderedBulkOp().insert({ _id: ObjectId(), name: "dave" }).insert({ _id: ObjectId(), name: "eve" })

// Insert documents with helper functions
db.events.initializeUnorderedBulkOp().insert({ type: "login", timestamp: ISODate() }).insert({ type: "logout", timestamp: ISODate() }).execute()
