// cursor.hasNext() - Check if there are more documents in the cursor

// Basic usage
db.users.find().hasNext()
db.users.find({}).hasNext()

// With query filter
db.users.find({ status: "active" }).hasNext()
db.users.find({ age: { $gt: 65 } }).hasNext()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).hasNext()
db.users.find().limit(10).hasNext()
db.users.find().skip(1000).hasNext()

// Check for existence of matching documents
db.users.find({ email: "admin@example.com" }).hasNext()
db.orders.find({ status: "pending", createdAt: { $lt: ISODate("2024-01-01") } }).hasNext()

// With projection
db.users.find().projection({ _id: 1 }).hasNext()
