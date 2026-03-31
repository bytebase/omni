// db.collection.getShardDistribution() - Print shard distribution statistics

// Basic shard distribution
db.users.getShardDistribution()
db.orders.getShardDistribution()
db.products.getShardDistribution()

// Get distribution for large collections
db.events.getShardDistribution()
db.logs.getShardDistribution()
db.analytics.getShardDistribution()

// Get distribution for sharded collections
db.transactions.getShardDistribution()
db.sessions.getShardDistribution()
db.inventory.getShardDistribution()

// Collection access patterns
db["users"].getShardDistribution()
db.getCollection("users").getShardDistribution()
db["user-data"].getShardDistribution()
db.getCollection("order.items").getShardDistribution()
