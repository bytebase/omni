// cursor.showRecordId() - Include the $recordId field in result documents

// Show record IDs
db.users.find().showRecordId(true)
db.users.find({ status: "active" }).showRecordId(true)

// Hide record IDs (default behavior)
db.users.find().showRecordId(false)

// With query filter
db.users.find({ age: { $gt: 25 } }).showRecordId(true)
db.orders.find({ total: { $gt: 1000 } }).showRecordId(true)

// Chained with other cursor methods
db.users.find().showRecordId(true).sort({ name: 1 })
db.users.find().showRecordId(true).limit(10)
db.users.find().showRecordId(true).skip(5).limit(5)

// With projection
db.users.find().projection({ name: 1, email: 1 }).showRecordId(true)

// Useful for debugging and data inspection
db.collection.find({ corrupted: true }).showRecordId(true)

// With pretty for readable output
db.users.find().showRecordId(true).pretty()
