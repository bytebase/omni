// db.collection.getPlanCache().clear() - Clear all cached query plans for a collection

// Basic usage
db.users.getPlanCache().clear()

// With collection access patterns
db["users"].getPlanCache().clear()
db.getCollection("users").getPlanCache().clear()

// Different collections
db.orders.getPlanCache().clear()
db.products.getPlanCache().clear()
db.sessions.getPlanCache().clear()
