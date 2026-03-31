// cursor.toArray() - Return all documents as an array

// Basic usage
db.users.find().toArray()
db.users.find({}).toArray()

// With query filter
db.users.find({ status: "active" }).toArray()
db.users.find({ age: { $gte: 18, $lte: 65 } }).toArray()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).toArray()
db.users.find().limit(10).toArray()
db.users.find().skip(5).limit(5).toArray()

// With sort and limit
db.users.find({ status: "active" }).sort({ createdAt: -1 }).limit(100).toArray()

// With projection
db.users.find().projection({ name: 1, email: 1 }).toArray()
db.users.find().project({ name: 1, _id: 0 }).toArray()

// Complex queries
db.orders.find({ $or: [{ status: "pending" }, { priority: "high" }] }).toArray()
db.logs.find({ level: { $in: ["error", "warn"] } }).sort({ timestamp: -1 }).limit(50).toArray()

// With maxTimeMS for timeout
db.largeCollection.find().limit(1000).maxTimeMS(30000).toArray()
