// sh.shardCollection() - Shard a collection

// Shard with hashed key
sh.shardCollection("mydb.users", { _id: "hashed" })

// Shard with range-based key
sh.shardCollection("mydb.orders", { customerId: 1 })

// Compound shard key
sh.shardCollection("mydb.events", { region: 1, timestamp: 1 })

// With unique constraint
sh.shardCollection("mydb.products", { sku: 1 }, { unique: true })

// With numInitialChunks
sh.shardCollection("mydb.logs", { timestamp: 1 }, { numInitialChunks: 4 })

// Hashed with presplit chunks
sh.shardCollection("mydb.telemetry", { deviceId: "hashed" }, { numInitialChunks: 10 })

// With collation
sh.shardCollection(
    "mydb.customers",
    { lastName: 1 },
    { collation: { locale: "en", strength: 2 } }
)

// With timeseries collection
sh.shardCollection(
    "mydb.metrics",
    { "metadata.sensorId": 1, timestamp: 1 },
    { timeseries: { timeField: "timestamp", metaField: "metadata" } }
)
