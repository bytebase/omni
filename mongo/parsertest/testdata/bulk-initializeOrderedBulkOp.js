// Bulk.initializeOrderedBulkOp() - Initialize ordered bulk operation

// Basic initialization
db.users.initializeOrderedBulkOp()

// Initialize and assign to variable
db.users.initializeOrderedBulkOp()

// Initialize with bracket notation
db["users"].initializeOrderedBulkOp()

// Initialize with getCollection
db.getCollection("users").initializeOrderedBulkOp()

// Ordered bulk operations execute in order
db.orders.initializeOrderedBulkOp()
