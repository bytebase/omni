// Bulk.execute() - Execute all bulk operations

// Basic execute
db.users.initializeUnorderedBulkOp().insert({ name: "alice" }).execute()

// Execute with write concern
db.users.initializeOrderedBulkOp().insert({ name: "bob" }).execute({ w: "majority" })

// Chained operations then execute
db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).remove().execute()

// Execute ordered operations
db.orders.initializeOrderedBulkOp().insert({ item: "apple" }).insert({ item: "banana" }).execute()

// Execute with detailed write concern
db.users.initializeOrderedBulkOp().insert({ name: "charlie" }).execute({ w: 1, j: true })
