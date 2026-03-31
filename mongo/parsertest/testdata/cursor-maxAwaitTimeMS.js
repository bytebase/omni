// cursor.maxAwaitTimeMS() - Set maximum time for tailable cursor getMore operations

// Basic usage
db.oplog.find().maxAwaitTimeMS(1000)
db.oplog.find().maxAwaitTimeMS(5000)

// With tailable cursor
db.capped.find().tailable().maxAwaitTimeMS(2000)
db.events.find().tailable(true).maxAwaitTimeMS(3000)

// Different timeout values
db.oplog.find().maxAwaitTimeMS(100)
db.oplog.find().maxAwaitTimeMS(10000)
db.oplog.find().maxAwaitTimeMS(30000)

// Chained with other cursor methods
db.logs.find({ level: "error" }).tailable().maxAwaitTimeMS(5000)
db.events.find().tailable().batchSize(100).maxAwaitTimeMS(2000)

// For change streams and tailable cursors
db.changes.find().tailable(true).maxAwaitTimeMS(1000).batchSize(10)
