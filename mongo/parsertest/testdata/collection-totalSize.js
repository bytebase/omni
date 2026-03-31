// db.collection.totalSize() - Get the total size of collection data plus indexes in bytes

// Basic usage
db.users.totalSize()
db.orders.totalSize()
db.products.totalSize()
db.logs.totalSize()

// Collection access patterns
db["users"].totalSize()
db.getCollection("users").totalSize()
db["large-collection"].totalSize()
db.getCollection("archived.data").totalSize()
