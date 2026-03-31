// Bulk.initializeUnorderedBulkOp() - Initialize unordered bulk operation

// Basic initialization
db.users.initializeUnorderedBulkOp()

// Initialize and assign to variable
db.users.initializeUnorderedBulkOp()

// Initialize with bracket notation
db["users"].initializeUnorderedBulkOp()

// Initialize with getCollection
db.getCollection("users").initializeUnorderedBulkOp()

// Unordered bulk operations may execute in any order for efficiency
db.orders.initializeUnorderedBulkOp()
