// db.collection.analyzeShardKey() - Analyze a shard key for a collection

// Analyze simple shard key
db.users.analyzeShardKey({ email: 1 })
db.orders.analyzeShardKey({ customerId: 1 })
db.products.analyzeShardKey({ category: 1 })

// Analyze compound shard key
db.orders.analyzeShardKey({ customerId: 1, createdAt: 1 })
db.events.analyzeShardKey({ tenantId: 1, timestamp: 1 })
db.logs.analyzeShardKey({ source: 1, level: 1, timestamp: 1 })

// Analyze hashed shard key
db.users.analyzeShardKey({ _id: "hashed" })
db.sessions.analyzeShardKey({ sessionId: "hashed" })

// Analyze with options
db.orders.analyzeShardKey({ customerId: 1 }, { keyCharacteristics: true })
db.products.analyzeShardKey({ sku: 1 }, { readWriteDistribution: true })

// Analyze with sampleRate
db.largeCollection.analyzeShardKey({ userId: 1 }, { sampleRate: 0.1 })
db.analytics.analyzeShardKey({ eventType: 1 }, { sampleRate: 0.05 })

// Analyze with sampleSize
db.orders.analyzeShardKey({ region: 1 }, { sampleSize: 10000 })

// Analyze with all options
db.transactions.analyzeShardKey(
    { accountId: 1, transactionDate: 1 },
    {
        keyCharacteristics: true,
        readWriteDistribution: true,
        sampleRate: 0.2
    }
)

// Collection access patterns
db["users"].analyzeShardKey({ userId: 1 })
db["users"].analyzeShardKey({ email: 1 }, { keyCharacteristics: true })
db.getCollection("users").analyzeShardKey({ tenantId: 1 })
db.getCollection("orders").analyzeShardKey({ customerId: 1, orderId: 1 })
db["event-logs"].analyzeShardKey({ source: 1, timestamp: 1 })
