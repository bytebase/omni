// sh.refineCollectionShardKey() - Refine the shard key of a sharded collection

// Basic usage - add suffix fields
sh.refineCollectionShardKey("mydb.users", { region: 1, email: 1 })

// Refine with additional timestamp
sh.refineCollectionShardKey("mydb.orders", { customerId: 1, orderId: 1, timestamp: 1 })

// Add more specificity to shard key
sh.refineCollectionShardKey("mydb.events", { category: 1, subCategory: 1, eventId: 1 })
