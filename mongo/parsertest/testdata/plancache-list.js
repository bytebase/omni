// db.collection.getPlanCache().list() - List query plan cache entries

// Basic usage
db.users.getPlanCache().list()

// With collection access patterns
db["users"].getPlanCache().list()
db.getCollection("users").getPlanCache().list()

// Different collections
db.orders.getPlanCache().list()
db.products.getPlanCache().list()

// With options (pipeline stages)
db.users.getPlanCache().list([{ $match: { isActive: true } }])

db.orders.getPlanCache().list([
    { $match: { planCacheKey: "abc123" } }
])

// List with projection
db.users.getPlanCache().list([
    { $project: { queryHash: 1, planCacheKey: 1 } }
])
