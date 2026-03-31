// sh.addShardTag() - Associate a tag with a shard (deprecated, use addShardToZone)

// Basic usage
sh.addShardTag("shard0000", "NYC")

// Add tag for geographic distribution
sh.addShardTag("shard0001", "LAX")

// Add multiple tags to same shard
sh.addShardTag("shard0002", "EU")
sh.addShardTag("shard0002", "EMEA")
