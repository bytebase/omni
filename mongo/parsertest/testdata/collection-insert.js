// db.collection.insert() - Insert documents (deprecated - use insertOne or insertMany)

// Insert single document
db.users.insert({ name: "alice", age: 25 })
db.users.insert({ name: "bob", email: "bob@example.com", status: "active" })

// Insert with _id
db.users.insert({ _id: 1, name: "charlie" })
db.users.insert({ _id: ObjectId("507f1f77bcf86cd799439011"), name: "dave" })

// Insert with nested document
db.users.insert({
    name: "eve",
    address: {
        street: "123 Main St",
        city: "New York",
        zip: "10001"
    }
})

// Insert with array
db.users.insert({
    name: "frank",
    tags: ["admin", "user", "moderator"],
    scores: [85, 90, 78]
})

// Insert multiple documents (array)
db.users.insert([
    { name: "grace", age: 30 },
    { name: "henry", age: 35 },
    { name: "ivy", age: 28 }
])

// Insert with write concern
db.users.insert({ name: "jack" }, { writeConcern: { w: "majority" } })
db.orders.insert({ orderId: 1 }, { writeConcern: { w: 1, j: true } })

// Insert with ordered option
db.users.insert([
    { name: "kate" },
    { name: "leo" },
    { name: "mike" }
], { ordered: false })

// Insert complex document with dates
db.events.insert({
    type: "login",
    userId: ObjectId("507f1f77bcf86cd799439011"),
    timestamp: Date(),
    metadata: { ip: "192.168.1.1", userAgent: "Mozilla/5.0" }
})

// Insert with ISODate
db.logs.insert({
    message: "Application started",
    level: "info",
    createdAt: ISODate("2024-01-15T10:30:00Z")
})

// Insert with various data types
db.mixed.insert({
    stringField: "text",
    numberField: 42,
    floatField: 3.14,
    boolField: true,
    nullField: null,
    arrayField: [1, 2, 3],
    objectField: { nested: "value" }
})

// Collection access patterns
db["users"].insert({ name: "nina" })
db["users"].insert([{ a: 1 }, { b: 2 }])
db.getCollection("users").insert({ name: "oscar" })
db.getCollection("user-data").insert({ type: "profile", data: {} })
db["event-log"].insert({ event: "click", target: "button" })
