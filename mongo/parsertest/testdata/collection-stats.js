// db.collection.stats() - Get collection statistics

// Basic stats
db.users.stats()
db.orders.stats()
db.products.stats()

// Stats with scale factor (convert bytes to KB, MB, etc.)
db.users.stats({ scale: 1024 })
db.largeCollection.stats({ scale: 1048576 })

// Stats with index details
db.users.stats({ indexDetails: true })

// Stats with multiple options
db.orders.stats({ scale: 1024, indexDetails: true })

// Collection access patterns
db["users"].stats()
db.getCollection("users").stats()
db["large-collection"].stats({ scale: 1024 })
db.getCollection("archived.data").stats({ indexDetails: true })
