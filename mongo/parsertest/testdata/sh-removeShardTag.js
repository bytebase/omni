// sh.removeShardTag() - Remove a tag from a shard (deprecated, use removeShardFromZone)

// Basic usage
sh.removeShardTag("shard0000", "NYC")

// Remove geographic tag
sh.removeShardTag("shard0001", "LAX")

// Remove data center tag
sh.removeShardTag("shard0002", "DC1")
