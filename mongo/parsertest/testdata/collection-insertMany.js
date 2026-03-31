// db.collection.insertMany() - Insert multiple documents

// Basic insert
db.users.insertMany([{ name: "alice" }, { name: "bob" }])

// Insert with various documents
db.users.insertMany([
    { name: "charlie", age: 25 },
    { name: "dave", age: 30 },
    { name: "eve", age: 35 }
])

// Insert with options
db.users.insertMany([{ name: "frank" }], { ordered: false })
db.users.insertMany([{ name: "grace" }], { writeConcern: { w: 1 } })

// Collection access patterns
db["users"].insertMany([{ name: "heidi" }])
db.getCollection("users").insertMany([{ name: "ivan" }])
