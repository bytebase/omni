// db.collection.getPlanCache().clearPlansByQuery() - Clear cached plans for a specific query shape

// Clear by query shape
db.users.getPlanCache().clearPlansByQuery({ status: "active" })

// Clear with projection and sort
db.users.getPlanCache().clearPlansByQuery({ status: "active" }, { name: 1 }, { age: -1 })

// Clear with only query and projection
db.orders.getPlanCache().clearPlansByQuery({ customer: "alice" }, { total: 1 })

// Collection access patterns
db["orders"].getPlanCache().clearPlansByQuery({ customer: "alice" })
db.getCollection("products").getPlanCache().clearPlansByQuery({ category: "electronics" })

// Complex query shapes
db.users.getPlanCache().clearPlansByQuery({
    status: "active",
    age: { $gte: 18 }
})

db.orders.getPlanCache().clearPlansByQuery(
    { status: "pending", total: { $gt: 100 } },
    { _id: 0, orderId: 1, total: 1 },
    { createdAt: -1 }
)
