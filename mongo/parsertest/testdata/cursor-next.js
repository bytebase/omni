// cursor.next() - Return the next document in the cursor

// Basic usage
db.users.find().next()
db.users.find({}).next()

// With query filter
db.users.find({ status: "active" }).next()
db.users.find({ role: "admin" }).next()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).next()
db.users.find().sort({ createdAt: -1 }).next()
db.users.find().skip(10).next()

// Get first matching document
db.users.find({ email: "admin@example.com" }).next()
db.orders.find({ status: "pending" }).sort({ createdAt: 1 }).next()

// With projection
db.users.find().projection({ name: 1, email: 1 }).next()

// With limit (still returns just one document)
db.users.find().limit(100).next()
