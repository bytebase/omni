// cursor.close() - Close the cursor and free server resources

// Basic usage - close cursor after getting some results
db.users.find().close()

// With query filter
db.users.find({ status: "active" }).close()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).close()
db.users.find().limit(10).close()
db.users.find().batchSize(100).close()

// Close after iterating
db.users.find({ status: "pending" }).sort({ createdAt: -1 }).close()

// With projection
db.users.find().projection({ name: 1 }).close()
