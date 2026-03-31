// sh.addShardToZone() - Associate a shard with a zone

// Basic usage
sh.addShardToZone("shard0000", "NYC")

// Add shard to geographic zone
sh.addShardToZone("shard0001", "EU-WEST")

// Add shard to data center zone
sh.addShardToZone("shard0002", "DC1")

// Multiple shards in same zone
sh.addShardToZone("shard0003", "US-EAST")
sh.addShardToZone("shard0004", "US-EAST")
