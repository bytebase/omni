// cursor.maxTimeMS() - Set maximum execution time for the query

// Basic usage
db.users.find().maxTimeMS(5000)
db.users.find().maxTimeMS(1000)
db.users.find().maxTimeMS(30000)

// With query filter
db.users.find({ status: "active" }).maxTimeMS(10000)
db.orders.find({ total: { $gt: 1000 } }).maxTimeMS(5000)

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).maxTimeMS(3000)
db.users.find().sort({ createdAt: -1 }).limit(100).maxTimeMS(5000)
db.users.find({ status: "active" }).skip(1000).limit(100).maxTimeMS(10000)

// For long-running queries
db.analytics.find({ date: { $gte: ISODate("2024-01-01") } }).maxTimeMS(60000)

// Short timeout for quick queries
db.users.find({ _id: ObjectId("507f1f77bcf86cd799439011") }).maxTimeMS(100)

// With aggregation-style queries
db.logs.find({ $text: { $search: "error" } }).maxTimeMS(15000)

// With explain
db.users.find({ status: "active" }).maxTimeMS(5000).explain()
