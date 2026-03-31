// cursor.tryNext() - Return next document without blocking (for tailable cursors)

// Basic usage
db.users.find().tryNext()
db.cappedCollection.find().tryNext()

// With tailable cursor
db.cappedCollection.find().tailable().tryNext()
db.oplog.find().tailable(true).tryNext()

// With query filter
db.events.find({ type: "notification" }).tailable().tryNext()
db.logs.find({ level: "error" }).tailable().tryNext()

// Chained with other cursor methods
db.cappedCollection.find().tailable().batchSize(10).tryNext()
db.events.find().tailable().maxAwaitTimeMS(1000).tryNext()

// For non-blocking iteration
db.logs.find({ timestamp: { $gte: ISODate("2024-01-01") } }).tailable().tryNext()

// With projection
db.events.find().tailable().projection({ message: 1 }).tryNext()

// Standard cursor (returns next or null)
db.users.find({ status: "active" }).tryNext()
