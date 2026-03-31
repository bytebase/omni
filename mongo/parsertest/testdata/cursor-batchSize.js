// cursor.batchSize() - Set batch size for cursor iteration

// Basic usage
db.users.find().batchSize(100)
db.users.find().batchSize(50)
db.users.find().batchSize(1000)

// With query filter
db.users.find({ status: "active" }).batchSize(50)
db.users.find({ age: { $gte: 18 } }).batchSize(200)

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).batchSize(25)
db.users.find().sort({ name: 1 }).batchSize(25).limit(100)
db.users.find().skip(100).limit(50).batchSize(10)

// Small batch size for memory-constrained environments
db.largeCollection.find().batchSize(10)

// Large batch size for faster iteration
db.users.find().batchSize(5000)

// With projection
db.users.find().projection({ name: 1 }).batchSize(100)
