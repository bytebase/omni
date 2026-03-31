// sh.enableSharding() - Enable sharding for a database

// Basic usage
sh.enableSharding("myDatabase")

// With primaryShard option
sh.enableSharding("myDatabase", { primaryShard: "shard0001" })

// Enable sharding for test database
sh.enableSharding("test")

// Enable with specific primary shard
sh.enableSharding("analytics", { primaryShard: "shard0002" })
