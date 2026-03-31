// cursor.readConcern() - Specify read concern level for the query

// Local read concern (default)
db.users.find().readConcern({ level: "local" })

// Available read concern
db.users.find().readConcern({ level: "available" })

// Majority read concern
db.users.find().readConcern({ level: "majority" })

// Linearizable read concern
db.users.find().readConcern({ level: "linearizable" })

// Snapshot read concern
db.users.find().readConcern({ level: "snapshot" })

// With query filter
db.users.find({ status: "active" }).readConcern({ level: "majority" })
db.orders.find({ total: { $gt: 1000 } }).readConcern({ level: "local" })

// Chained with other cursor methods
db.users.find().readConcern({ level: "majority" }).sort({ name: 1 })
db.users.find().readConcern({ level: "majority" }).limit(10)
db.users.find({ status: "active" }).readConcern({ level: "majority" }).maxTimeMS(5000)

// For consistency-critical queries
db.inventory.find({ sku: "abc123" }).readConcern({ level: "linearizable" })
db.balances.find({ accountId: "12345" }).readConcern({ level: "majority" })
