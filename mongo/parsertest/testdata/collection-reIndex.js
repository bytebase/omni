// db.collection.reIndex() - Rebuild all indexes on a collection (deprecated)

// Basic reIndex
db.users.reIndex()
db.orders.reIndex()
db.products.reIndex()

// reIndex after heavy write operations
db.logs.reIndex()
db.events.reIndex()
db.sessions.reIndex()

// reIndex on large collections
db.analytics.reIndex()
db.transactions.reIndex()

// Collection access patterns
db["users"].reIndex()
db.getCollection("users").reIndex()
db["user-data"].reIndex()
db.getCollection("order.items").reIndex()
