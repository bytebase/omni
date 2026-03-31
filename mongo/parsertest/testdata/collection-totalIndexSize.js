// db.collection.totalIndexSize() - Get the total size of all indexes on a collection in bytes

// Basic usage
db.users.totalIndexSize()
db.orders.totalIndexSize()
db.products.totalIndexSize()
db.sessions.totalIndexSize()

// Collection access patterns
db["users"].totalIndexSize()
db.getCollection("users").totalIndexSize()
db["indexed-collection"].totalIndexSize()
db.getCollection("heavily.indexed").totalIndexSize()
