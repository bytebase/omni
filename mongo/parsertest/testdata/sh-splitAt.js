// sh.splitAt() - Split a chunk at a specific shard key value

// Basic usage
sh.splitAt("mydb.users", { zipcode: "50000" })

// Split at specific point
sh.splitAt("test.orders", { orderId: 1000000 })

// Split with compound key
sh.splitAt("analytics.events", { region: "US", timestamp: ISODate("2024-06-01") })

// Split at ObjectId boundary
sh.splitAt("mydb.documents", { _id: ObjectId("507f1f77bcf86cd799439011") })
