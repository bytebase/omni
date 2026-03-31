// db.createView() - Create a view

// Basic view
db.createView("activeUsers", "users", [
    { $match: { status: "active" } }
])

// View with multiple pipeline stages
db.createView("orderSummary", "orders", [
    { $match: { status: "completed" } },
    { $group: { _id: "$customerId", totalAmount: { $sum: "$amount" } } },
    { $sort: { totalAmount: -1 } }
])

// View with lookup
db.createView("userOrders", "users", [
    { $lookup: {
        from: "orders",
        localField: "_id",
        foreignField: "userId",
        as: "orders"
    }}
])

// View with collation
db.createView("sortedProducts", "products", [
    { $sort: { name: 1 } }
], {
    collation: { locale: "en", strength: 2 }
})

// View with projection
db.createView("publicUsers", "users", [
    { $project: {
        name: 1,
        email: 1,
        createdAt: 1,
        _id: 0
    }}
])
