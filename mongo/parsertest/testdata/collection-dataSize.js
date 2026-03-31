// db.collection.dataSize() - Get the size of the collection's data in bytes

// Basic usage
db.users.dataSize()
db.orders.dataSize()
db.products.dataSize()
db.logs.dataSize()

// Collection access patterns
db["users"].dataSize()
db.getCollection("users").dataSize()
db["large-collection"].dataSize()
db.getCollection("archived.orders").dataSize()
