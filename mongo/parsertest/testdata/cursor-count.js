// cursor.count() - Count the number of documents matching the query

// Basic usage
db.users.find().count()
db.users.find({}).count()

// With query filter
db.users.find({ status: "active" }).count()
db.users.find({ age: { $gte: 18 } }).count()
db.users.find({ "address.city": "New York" }).count()

// After sort (sort doesn't affect count)
db.users.find().sort({ name: 1 }).count()

// After limit/skip (count ignores limit/skip by default)
db.users.find().limit(10).count()
db.users.find().skip(5).count()
db.users.find().skip(5).limit(10).count()

// Complex query with count
db.users.find({ $or: [{ status: "active" }, { role: "admin" }] }).count()
db.orders.find({ total: { $gt: 100 }, status: "completed" }).count()
