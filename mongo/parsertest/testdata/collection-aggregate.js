// db.collection.aggregate() - Aggregation pipeline

// Empty pipeline
db.orders.aggregate([])

// Single stage pipelines
db.orders.aggregate([{ $match: { status: "completed" } }])
db.orders.aggregate([{ $group: { _id: "$category", count: { $sum: 1 } } }])
db.orders.aggregate([{ $sort: { createdAt: -1 } }])
db.orders.aggregate([{ $limit: 10 }])
db.orders.aggregate([{ $skip: 20 }])
db.orders.aggregate([{ $project: { name: 1, total: 1, _id: 0 } }])
db.orders.aggregate([{ $unwind: "$items" }])
db.orders.aggregate([{ $count: "totalOrders" }])

// $match stage variations
db.users.aggregate([{ $match: { age: { $gt: 18 } } }])
db.users.aggregate([{ $match: { status: { $in: ["active", "pending"] } } }])
db.users.aggregate([{ $match: { $or: [{ role: "admin" }, { role: "moderator" }] } }])
db.users.aggregate([{ $match: { "address.country": "USA" } }])

// $group stage variations
db.orders.aggregate([{ $group: { _id: "$customerId", total: { $sum: "$amount" } } }])
db.orders.aggregate([{ $group: { _id: "$category", avgPrice: { $avg: "$price" } } }])
db.orders.aggregate([{ $group: { _id: "$status", count: { $sum: 1 }, items: { $push: "$name" } } }])
db.orders.aggregate([{ $group: { _id: null, totalRevenue: { $sum: "$amount" } } }])
db.sales.aggregate([{ $group: { _id: { year: { $year: "$date" }, month: { $month: "$date" } }, total: { $sum: "$amount" } } }])

// $project stage variations
db.users.aggregate([{ $project: { name: 1, email: 1 } }])
db.users.aggregate([{ $project: { password: 0, ssn: 0 } }])
db.users.aggregate([{ $project: { fullName: { $concat: ["$firstName", " ", "$lastName"] } } }])
db.orders.aggregate([{ $project: { total: { $multiply: ["$price", "$quantity"] } } }])

// $sort stage variations
db.users.aggregate([{ $sort: { name: 1 } }])
db.users.aggregate([{ $sort: { createdAt: -1 } }])
db.users.aggregate([{ $sort: { lastName: 1, firstName: 1 } }])

// $lookup stage (join)
db.orders.aggregate([{ $lookup: { from: "users", localField: "customerId", foreignField: "_id", as: "customer" } }])
db.orders.aggregate([{ $lookup: { from: "products", localField: "productIds", foreignField: "_id", as: "products" } }])

// Multi-stage pipelines
db.orders.aggregate([
    { $match: { status: "completed" } },
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } },
    { $sort: { total: -1 } },
    { $limit: 10 }
])

db.sales.aggregate([
    { $match: { date: { $gte: ISODate("2024-01-01"), $lt: ISODate("2025-01-01") } } },
    { $group: {
        _id: { year: { $year: "$date" }, month: { $month: "$date" } },
        totalRevenue: { $sum: "$amount" },
        avgOrderValue: { $avg: "$amount" },
        orderCount: { $sum: 1 }
    } },
    { $sort: { "_id.year": 1, "_id.month": 1 } }
])

db.users.aggregate([
    { $match: { status: "active" } },
    { $lookup: { from: "orders", localField: "_id", foreignField: "customerId", as: "orders" } },
    { $project: { name: 1, email: 1, orderCount: { $size: "$orders" } } },
    { $sort: { orderCount: -1 } },
    { $limit: 100 }
])

// Pipeline with $addFields
db.orders.aggregate([
    { $addFields: { totalWithTax: { $multiply: ["$total", 1.1] } } }
])

// Pipeline with $set (alias for $addFields)
db.orders.aggregate([
    { $set: { processed: true, processedAt: ISODate() } }
])

// Pipeline with $unset
db.users.aggregate([
    { $unset: ["password", "ssn", "internalNotes"] }
])

// Pipeline with $replaceRoot
db.orders.aggregate([
    { $replaceRoot: { newRoot: "$shipping" } }
])

// Pipeline with $facet (multiple pipelines)
db.products.aggregate([
    { $facet: {
        categoryCounts: [{ $group: { _id: "$category", count: { $sum: 1 } } }],
        priceStats: [{ $group: { _id: null, avgPrice: { $avg: "$price" }, maxPrice: { $max: "$price" } } }]
    } }
])

// Aggregate with options (options passed to driver)
db.orders.aggregate([{ $match: { status: "completed" } }], { allowDiskUse: true })
db.orders.aggregate([{ $group: { _id: "$category" } }], { maxTimeMS: 60000 })
db.orders.aggregate([{ $sort: { total: -1 } }], { collation: { locale: "en" } })

// Aggregate with collection access patterns
db["orders"].aggregate([{ $match: { status: "pending" } }])
db['audit-logs'].aggregate([{ $group: { _id: "$action", count: { $sum: 1 } } }])
db.getCollection("sales").aggregate([{ $match: { year: 2024 } }])
