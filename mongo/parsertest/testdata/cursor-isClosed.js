// cursor.isClosed() - Check if the cursor is closed

// Basic usage
db.users.find().isClosed()
db.users.find({}).isClosed()

// With query filter
db.users.find({ status: "active" }).isClosed()
db.orders.find({ total: { $gt: 100 } }).isClosed()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).isClosed()
db.users.find().limit(10).isClosed()
db.users.find().batchSize(100).isClosed()

// Check status after operations
db.users.find({ status: "pending" }).skip(50).limit(25).isClosed()

// With projection
db.users.find().projection({ name: 1 }).isClosed()
