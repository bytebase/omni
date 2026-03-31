// cursor.tailable() - Create a tailable cursor for capped collections

// Basic tailable cursor
db.cappedCollection.find().tailable()
db.oplog.find().tailable()

// With explicit true/false
db.cappedCollection.find().tailable(true)
db.cappedCollection.find().tailable(false)

// With query filter
db.logs.find({ level: "error" }).tailable()
db.events.find({ type: "notification" }).tailable(true)

// Tailable with awaitData (for blocking behavior)
db.cappedCollection.find().tailable().maxAwaitTimeMS(1000)
db.oplog.find().tailable(true).maxAwaitTimeMS(5000)

// Chained with other cursor methods
db.cappedCollection.find().tailable().batchSize(10)
db.events.find().tailable(true).batchSize(100).maxAwaitTimeMS(2000)

// For real-time log tailing
db.logs.find({ timestamp: { $gte: ISODate("2024-01-01") } }).tailable()

// With projection
db.events.find().tailable().projection({ message: 1, timestamp: 1 })

// Oplog tailing pattern
db.oplog.find({ ts: { $gt: Timestamp(0, 0) } }).tailable(true)
