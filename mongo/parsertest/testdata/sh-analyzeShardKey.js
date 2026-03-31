// sh.analyzeShardKey() - Analyze a potential shard key

// Basic usage
sh.analyzeShardKey("mydb.users", { email: 1 })

// Analyze hashed key
sh.analyzeShardKey("mydb.orders", { orderId: "hashed" })

// Analyze compound key
sh.analyzeShardKey("mydb.events", { region: 1, timestamp: 1 })

// With options
sh.analyzeShardKey("mydb.logs", { deviceId: 1 }, { sampleRate: 0.5 })
