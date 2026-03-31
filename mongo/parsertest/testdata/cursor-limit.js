// cursor.limit() - Limit the number of documents returned

// Basic usage
db.users.find().limit(10)
db.users.find().limit(1)
db.users.find().limit(100)

// With query filter
db.users.find({ status: "active" }).limit(50)
db.users.find({ age: { $gte: 18 } }).limit(25)

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).limit(10)
db.users.find().skip(20).limit(10)
db.users.find({ status: "active" }).sort({ createdAt: -1 }).limit(5)

// Limit with projection
db.users.find().projection({ name: 1, email: 1 }).limit(10)

// Limit 0 returns all documents (no limit)
db.users.find().limit(0)
