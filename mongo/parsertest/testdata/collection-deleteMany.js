// db.collection.deleteMany() - Delete multiple documents

// Basic delete
db.users.deleteMany({ status: "inactive" })

// Delete all documents
db.users.deleteMany({})

// Delete with comparison operators
db.users.deleteMany({ lastLogin: { $lt: ISODate("2023-01-01") } })

// Delete with logical operators
db.users.deleteMany({ $or: [{ deleted: true }, { banned: true }] })

// Delete with options
db.users.deleteMany({ temp: true }, { writeConcern: { w: "majority" } })

// Collection access patterns
db["logs"].deleteMany({ level: "debug" })
db.getCollection("sessions").deleteMany({ expired: true })
