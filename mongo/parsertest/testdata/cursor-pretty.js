// cursor.pretty() - Format output for readability in the shell

// Basic usage
db.users.find().pretty()
db.users.find({}).pretty()

// With query filter
db.users.find({ status: "active" }).pretty()
db.users.find({ age: { $gt: 25 } }).pretty()

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).pretty()
db.users.find().limit(10).pretty()
db.users.find().skip(5).limit(5).pretty()

// Complex queries with pretty output
db.users.find({ $or: [{ status: "active" }, { role: "admin" }] }).pretty()
db.orders.find({ items: { $elemMatch: { price: { $gt: 100 } } } }).pretty()

// With projection
db.users.find().projection({ name: 1, email: 1, address: 1 }).pretty()

// Nested document display
db.users.find({ "address.city": "New York" }).pretty()
