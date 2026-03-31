// cursor.comment() - Add a comment to the query for profiling/logging

// Basic usage
db.users.find().comment("User listing query")
db.users.find().comment("Admin dashboard query")

// With query filter
db.users.find({ status: "active" }).comment("Active users query")
db.orders.find({ total: { $gt: 1000 } }).comment("High value orders")

// Chained with other cursor methods
db.users.find().sort({ createdAt: -1 }).comment("Recent users").limit(10)
db.users.find({ role: "admin" }).comment("Admin lookup").projection({ name: 1, email: 1 })

// Detailed comments for debugging
db.users.find({ status: "pending" }).comment("Pending approval - ticket #12345")
db.analytics.find({ date: { $gte: ISODate("2024-01-01") } }).comment("Q1 2024 analytics export")

// Comment with special characters
db.logs.find({ level: "error" }).comment("Error logs - check /var/log")
