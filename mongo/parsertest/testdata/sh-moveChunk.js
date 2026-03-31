// sh.moveChunk() - Move a chunk to a different shard

// Move by find query
sh.moveChunk("mydb.users", { zipcode: "10001" }, "shard0001")

// Move chunk containing specific document
sh.moveChunk("test.orders", { _id: ObjectId("507f1f77bcf86cd799439011") }, "shard0002")

// Move by shard key value
sh.moveChunk("analytics.events", { region: "US", timestamp: ISODate("2024-01-01") }, "shard0003")
