// cursor.size() - Return count of documents considering skip and limit

// Basic usage
db.users.find().size()
db.users.find({}).size()

// With query filter
db.users.find({ status: "active" }).size()
db.users.find({ age: { $gte: 18 } }).size()

// Size considers limit and skip (unlike count)
db.users.find().limit(10).size()
db.users.find().skip(5).size()
db.users.find().skip(10).limit(5).size()

// Compare with count behavior
db.users.find().limit(100).size()
db.users.find().skip(50).limit(25).size()

// With sort (doesn't affect size)
db.users.find().sort({ name: 1 }).size()

// Complex queries
db.orders.find({ status: "completed", total: { $gt: 100 } }).size()
db.logs.find({ level: { $in: ["error", "fatal"] } }).skip(10).size()

// With projection
db.users.find().projection({ _id: 1 }).size()
