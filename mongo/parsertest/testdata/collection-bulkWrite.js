// db.collection.bulkWrite() - Bulk write operations (unsupported - use individual write commands)

// Basic bulk write with insert
db.users.bulkWrite([
    { insertOne: { document: { name: "alice", age: 25 } } }
])

// Bulk write with multiple operations
db.users.bulkWrite([
    { insertOne: { document: { name: "alice", age: 25 } } },
    { insertOne: { document: { name: "bob", age: 30 } } },
    { updateOne: { filter: { name: "charlie" }, update: { $set: { age: 35 } } } },
    { updateMany: { filter: { status: "pending" }, update: { $set: { status: "active" } } } },
    { deleteOne: { filter: { name: "dave" } } },
    { deleteMany: { filter: { status: "deleted" } } },
    { replaceOne: { filter: { name: "eve" }, replacement: { name: "eve", age: 28 } } }
])

// Bulk write with options
db.orders.bulkWrite([
    { insertOne: { document: { orderId: 1, status: "new" } } },
    { updateOne: { filter: { orderId: 2 }, update: { $set: { status: "shipped" } } } }
], { ordered: false })

db.products.bulkWrite([
    { deleteOne: { filter: { discontinued: true } } }
], { ordered: true })

// Bulk write with write concern
db.inventory.bulkWrite([
    { updateMany: { filter: { qty: { $lt: 10 } }, update: { $set: { reorder: true } } } }
], { writeConcern: { w: "majority" } })

// Bulk write with upsert in update operations
db.users.bulkWrite([
    { updateOne: { filter: { email: "test@example.com" }, update: { $set: { name: "Test User" } }, upsert: true } },
    { replaceOne: { filter: { code: "ABC" }, replacement: { code: "ABC", value: 100 }, upsert: true } }
])

// Bulk write with array filters
db.inventory.bulkWrite([
    { updateOne: { filter: { item: "abc123" }, update: { $set: { "sizes.$[size].qty": 0 } }, arrayFilters: [{ "size.type": "small" }] } }
])

// Bulk write with collation
db.products.bulkWrite([
    { deleteMany: { filter: { category: "electronics" } } }
], { ordered: false, collation: { locale: "en", strength: 2 } })

// Collection access patterns
db["users"].bulkWrite([{ insertOne: { document: { x: 1 } } }])
db.getCollection("users").bulkWrite([{ insertOne: { document: { y: 1 } } }])
db["user-data"].bulkWrite([{ deleteOne: { filter: { expired: true } } }], { ordered: false })
