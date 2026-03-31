// sh.reshardCollection() - Reshard a collection with a new shard key

// Basic usage - change to hashed shard key
sh.reshardCollection("mydb.users", { email: "hashed" })

// Reshard with range-based key
sh.reshardCollection("mydb.orders", { orderId: 1 })

// Reshard with compound key
sh.reshardCollection("mydb.events", { region: 1, timestamp: 1 })

// With options
sh.reshardCollection(
    "mydb.products",
    { category: 1, productId: 1 },
    { unique: false, numInitialChunks: 8 }
)

// With zone specification
sh.reshardCollection(
    "mydb.customers",
    { country: 1, customerId: 1 },
    { zones: [{ zone: "EU", min: { country: "AA" }, max: { country: "MZ" } }] }
)
