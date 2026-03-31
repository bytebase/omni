// cursor.map() - Apply a function to each document and return array of results

// Basic usage with document argument
db.users.find().map({ field: "name" })

// With query filter
db.users.find({ status: "active" }).map({ extract: "email" })
db.orders.find({ total: { $gt: 100 } }).map({ compute: "tax" })

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).map({ field: "name" })
db.users.find().limit(10).map({ transform: "user" })
db.users.find({ status: "active" }).sort({ createdAt: -1 }).limit(5).map({ format: "output" })

// With projection
db.users.find().projection({ name: 1, email: 1 }).map({ concat: "contact" })

// With string argument representing transformation name
db.products.find({ inStock: true }).map("applyDiscount")
db.logs.find({ level: "error" }).map("extractMessage")
