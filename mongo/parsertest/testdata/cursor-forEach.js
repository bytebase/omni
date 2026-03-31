// cursor.forEach() - Iterate over each document with a callback function

// Basic usage with document argument
db.users.find().forEach({ handler: "print" })

// With query filter
db.users.find({ status: "active" }).forEach({ action: "log" })
db.orders.find({ total: { $gt: 100 } }).forEach({ type: "process" })

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).forEach({ output: "console" })
db.users.find().limit(10).forEach({ format: "json" })
db.users.find({ status: "active" }).sort({ createdAt: -1 }).limit(5).forEach({ mode: "verbose" })

// With projection
db.users.find().projection({ name: 1, email: 1 }).forEach({ display: true })

// With string argument representing callback name
db.users.find().forEach("printjson")
db.logs.find({ level: "error" }).forEach("processError")
