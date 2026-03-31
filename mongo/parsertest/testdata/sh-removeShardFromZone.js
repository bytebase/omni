// sh.removeShardFromZone() - Remove a shard from a zone

// Basic usage
sh.removeShardFromZone("shard0000", "NYC")

// Remove from geographic zone
sh.removeShardFromZone("shard0001", "EU-WEST")

// Remove from data center zone
sh.removeShardFromZone("shard0002", "DC1")
