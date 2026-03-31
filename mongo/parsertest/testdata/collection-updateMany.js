// db.collection.updateMany() - Update multiple documents

// Basic update
db.users.updateMany({ status: "inactive" }, { $set: { status: "active" } })

// Update with comparison operators
db.users.updateMany({ age: { $lt: 18 } }, { $set: { category: "minor" } })

// Update with multiple operators
db.users.updateMany({ verified: false }, { $set: { verified: true }, $currentDate: { verifiedAt: true } })

// Update with options
db.users.updateMany({ country: "US" }, { $set: { region: "NA" } }, { writeConcern: { w: "majority" } })

// Collection access patterns
db["users"].updateMany({}, { $set: { migrated: true } })
db.getCollection("users").updateMany({ role: "guest" }, { $set: { role: "user" } })
