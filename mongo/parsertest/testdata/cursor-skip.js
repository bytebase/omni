// cursor.skip() - Skip a number of documents in the result set

// Basic usage
db.users.find().skip(10)
db.users.find().skip(0)
db.users.find().skip(100)

// With query filter
db.users.find({ status: "active" }).skip(50)
db.users.find({ age: { $gte: 18 } }).skip(25)

// Chained with other cursor methods (pagination pattern)
db.users.find().skip(0).limit(10)
db.users.find().skip(10).limit(10)
db.users.find().skip(20).limit(10)

// With sort for consistent pagination
db.users.find().sort({ _id: 1 }).skip(100).limit(10)
db.users.find({ status: "active" }).sort({ createdAt: -1 }).skip(50).limit(25)

// Skip with projection
db.users.find().projection({ name: 1 }).skip(5).limit(5)
