// db.collection.getPlanCache() - Get the plan cache object for a collection

// Basic usage - get the plan cache object
db.users.getPlanCache()

// With collection access patterns
db["users"].getPlanCache()
db.getCollection("users").getPlanCache()

// Different collection names
db.orders.getPlanCache()
db.products.getPlanCache()
db.customers.getPlanCache()
