// sh.unshardCollection() - Unshard a collection (move all data to a single shard)

// Basic usage
sh.unshardCollection("mydb.users")

// Unshard to specific shard
sh.unshardCollection("mydb.orders", "shard0001")

// With options
sh.unshardCollection("mydb.logs", { toShard: "shard0002" })
