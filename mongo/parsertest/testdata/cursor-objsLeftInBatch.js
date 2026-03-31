// cursor.objsLeftInBatch() - Return number of documents remaining in current batch

// Basic usage
db.users.find().objsLeftInBatch()
db.users.find({}).objsLeftInBatch()

// With query filter
db.users.find({ status: "active" }).objsLeftInBatch()
db.orders.find({ total: { $gt: 100 } }).objsLeftInBatch()

// Chained with batchSize
db.users.find().batchSize(100).objsLeftInBatch()
db.users.find().batchSize(50).objsLeftInBatch()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).objsLeftInBatch()
db.users.find().limit(1000).objsLeftInBatch()
db.users.find().skip(100).batchSize(25).objsLeftInBatch()

// For monitoring batch consumption
db.logs.find({ level: "error" }).batchSize(10).objsLeftInBatch()

// With projection
db.users.find().projection({ name: 1 }).objsLeftInBatch()
