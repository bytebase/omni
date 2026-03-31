// db.collection.getShardVersion() - Get the shard version for a collection

// Basic shard version
db.users.getShardVersion()
db.orders.getShardVersion()
db.products.getShardVersion()

// Get version for various collections
db.events.getShardVersion()
db.logs.getShardVersion()
db.analytics.getShardVersion()

// Get version for sharded collections
db.transactions.getShardVersion()
db.sessions.getShardVersion()
db.inventory.getShardVersion()

// Collection access patterns
db["users"].getShardVersion()
db.getCollection("users").getShardVersion()
db["user-data"].getShardVersion()
db.getCollection("order.items").getShardVersion()
