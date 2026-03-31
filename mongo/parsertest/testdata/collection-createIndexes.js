// db.collection.createIndexes() - Create multiple indexes on a collection

// Create multiple indexes at once
db.users.createIndexes([
    { key: { email: 1 }, unique: true },
    { key: { username: 1 } },
    { key: { createdAt: -1 } }
])

// Create indexes with options
db.orders.createIndexes([
    { key: { customerId: 1 } },
    { key: { status: 1 }, sparse: true },
    { key: { orderDate: -1 }, name: "order_date_idx" }
])

// Create compound indexes
db.products.createIndexes([
    { key: { category: 1, name: 1 } },
    { key: { price: 1, rating: -1 } }
])

// With write concern options
db.users.createIndexes(
    [{ key: { email: 1 } }],
    { commitQuorum: "majority" }
)

// Collection access patterns
db["users"].createIndexes([{ key: { email: 1 } }])
db.getCollection("users").createIndexes([{ key: { email: 1 } }])
