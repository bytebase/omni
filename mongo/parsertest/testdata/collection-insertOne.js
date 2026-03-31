// db.collection.insertOne() - Insert a single document

// Basic insert
db.users.insertOne({ name: "alice", age: 25 })
db.users.insertOne({ _id: ObjectId(), name: "bob" })

// Insert with nested document
db.users.insertOne({ name: "charlie", address: { city: "NYC", zip: "10001" } })

// Insert with array
db.users.insertOne({ name: "dave", tags: ["admin", "user"] })

// Insert with helper functions
db.users.insertOne({ name: "eve", createdAt: ISODate(), id: UUID("550e8400-e29b-41d4-a716-446655440000") })

// Insert with options
db.users.insertOne({ name: "frank" }, { writeConcern: { w: "majority" } })

// Collection access patterns
db["users"].insertOne({ name: "grace" })
db.getCollection("users").insertOne({ name: "heidi" })
