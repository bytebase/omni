// db.collection.getPlanCache().help() - Display help information for plan cache methods

// Basic usage
db.users.getPlanCache().help()

// With collection access patterns
db["users"].getPlanCache().help()
db.getCollection("users").getPlanCache().help()

// Different collections
db.orders.getPlanCache().help()
db.products.getPlanCache().help()
