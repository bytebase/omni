// cursor.isExhausted() - Check if the cursor has no more documents and is closed

// Basic usage
db.users.find().isExhausted()
db.users.find({}).isExhausted()

// With query filter
db.users.find({ status: "active" }).isExhausted()
db.users.find({ role: "superadmin" }).isExhausted()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).isExhausted()
db.users.find().limit(10).isExhausted()
db.users.find().skip(1000000).isExhausted()

// Check exhaustion status
db.users.find({ status: "deleted" }).batchSize(10).isExhausted()
db.logs.find({ timestamp: { $lt: ISODate("2020-01-01") } }).isExhausted()

// With projection
db.users.find().projection({ _id: 1 }).isExhausted()
