// db.collection.updateOne() - Update a single document

// Basic update with $set
db.users.updateOne({ name: "alice" }, { $set: { age: 26 } })

// Update with multiple operators
db.users.updateOne({ _id: ObjectId("507f1f77bcf86cd799439011") }, { $set: { status: "active" }, $inc: { loginCount: 1 } })

// Update with upsert
db.users.updateOne({ email: "new@example.com" }, { $set: { name: "new user" } }, { upsert: true })

// Update with array filters
db.users.updateOne({ name: "alice" }, { $set: { "grades.$[elem].score": 100 } }, { arrayFilters: [{ "elem.grade": { $gte: 85 } }] })

// Collection access patterns
db["users"].updateOne({ name: "bob" }, { $set: { active: true } })
db.getCollection("users").updateOne({ name: "charlie" }, { $unset: { tempField: "" } })
