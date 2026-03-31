// db.collection.getIndexes() - Get all indexes on a collection

// Basic getIndexes
db.users.getIndexes()
db.orders.getIndexes()
db.products.getIndexes()
db.sessions.getIndexes()

// GetIndexes with collection access patterns
db["users"].getIndexes()
db['audit-logs'].getIndexes()
db["user-sessions"].getIndexes()
db.getCollection("users").getIndexes()
db.getCollection("my.collection").getIndexes()
db.getCollection("special-collection").getIndexes()
