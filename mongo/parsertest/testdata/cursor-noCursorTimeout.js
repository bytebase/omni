// cursor.noCursorTimeout() - Prevent cursor from timing out on the server

// Basic usage
db.users.find().noCursorTimeout()
db.users.find({}).noCursorTimeout()

// With query filter
db.users.find({ status: "active" }).noCursorTimeout()
db.largeCollection.find({ processed: false }).noCursorTimeout()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).noCursorTimeout()
db.users.find().batchSize(100).noCursorTimeout()
db.users.find().limit(10000).noCursorTimeout()

// For long-running batch operations
db.logs.find({ level: "error" }).noCursorTimeout().batchSize(50)
db.analytics.find({ date: { $gte: ISODate("2024-01-01") } }).noCursorTimeout()

// Chained after noCursorTimeout
db.documents.find({ needsProcessing: true }).noCursorTimeout().toArray()

// With projection
db.users.find().projection({ _id: 1, data: 1 }).noCursorTimeout()
