// cursor.itcount() - Count documents by iterating through the cursor

// Basic usage
db.users.find().itcount()
db.users.find({}).itcount()

// With query filter
db.users.find({ status: "active" }).itcount()
db.users.find({ age: { $gte: 18, $lte: 65 } }).itcount()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).itcount()
db.users.find().limit(100).itcount()
db.users.find().skip(50).itcount()

// Unlike count(), itcount() respects limit and skip
db.users.find().skip(10).limit(5).itcount()

// Complex queries
db.orders.find({ status: "completed", total: { $gt: 500 } }).itcount()
db.logs.find({ level: { $in: ["error", "fatal"] } }).itcount()

// With projection
db.users.find().projection({ _id: 1 }).itcount()
