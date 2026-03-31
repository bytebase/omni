// db.createCollection() - Create a new collection with options

// Basic collection creation
db.createCollection("users")
db.createCollection("orders")
db.createCollection("products")

// Create capped collection
db.createCollection("logs", { capped: true, size: 10000 })
db.createCollection("events", { capped: true, size: 100000, max: 5000 })

// Create collection with validation
db.createCollection("contacts", {
    validator: { $jsonSchema: {
        bsonType: "object",
        required: ["name", "email"],
        properties: {
            name: { bsonType: "string" },
            email: { bsonType: "string" }
        }
    }}
})

// Create collection with storage engine options
db.createCollection("archive", {
    storageEngine: { wiredTiger: { configString: "block_compressor=zstd" } }
})

// Create time series collection
db.createCollection("weather", {
    timeseries: {
        timeField: "timestamp",
        metaField: "metadata",
        granularity: "hours"
    }
})

// Create collection with collation
db.createCollection("users_intl", {
    collation: { locale: "en_US", strength: 2 }
})

// Create clustered collection
db.createCollection("orders_clustered", {
    clusteredIndex: {
        key: { _id: 1 },
        unique: true
    }
})
