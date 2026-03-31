// db.collection.storageSize() - Get the storage size of a collection in bytes

// Basic usage
db.users.storageSize()
db.orders.storageSize()
db.products.storageSize()
db.logs.storageSize()

// Collection access patterns
db["users"].storageSize()
db.getCollection("users").storageSize()
db["large-collection"].storageSize()
db.getCollection("archived.orders").storageSize()
