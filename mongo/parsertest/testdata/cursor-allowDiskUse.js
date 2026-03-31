// cursor.allowDiskUse() - Allow query to use disk for large sorts

// Basic usage (enable disk use)
db.users.find().allowDiskUse()
db.users.find().allowDiskUse(true)

// Disable disk use
db.users.find().allowDiskUse(false)

// With query filter
db.users.find({ status: "active" }).allowDiskUse()
db.orders.find({ total: { $gt: 1000 } }).allowDiskUse(true)

// Essential for large sorts
db.users.find().sort({ name: 1 }).allowDiskUse()
db.largeCollection.find().sort({ createdAt: -1, name: 1 }).allowDiskUse(true)

// Chained with other cursor methods
db.users.find().sort({ score: -1 }).allowDiskUse().limit(1000)
db.users.find({ status: "active" }).sort({ name: 1 }).allowDiskUse(true).skip(100).limit(50)

// For queries exceeding memory limits
db.analytics.find({ date: { $gte: ISODate("2024-01-01") } }).sort({ date: -1 }).allowDiskUse()
db.logs.find().sort({ timestamp: -1, level: 1 }).allowDiskUse(true)

// With projection to reduce memory
db.users.find().projection({ name: 1, email: 1 }).sort({ name: 1 }).allowDiskUse()

// With explain
db.users.find().sort({ name: 1 }).allowDiskUse().explain("executionStats")
