// db.collection.distinct() - Find distinct values for a field

// Distinct with field only (required)
db.users.distinct("status")
db.users.distinct("country")
db.users.distinct("role")
db.orders.distinct("category")
db.products.distinct("brand")

// Distinct with nested field
db.users.distinct("address.city")
db.users.distinct("address.country")
db.users.distinct("profile.department")

// Distinct with field and empty query
db.users.distinct("status", {})
db.users.distinct("city", {})

// Distinct with field and filter query
db.users.distinct("city", { country: "USA" })
db.users.distinct("status", { active: true })
db.users.distinct("role", { department: "engineering" })
db.orders.distinct("productId", { status: "completed" })

// Distinct with comparison operators in query
db.users.distinct("city", { age: { $gt: 18 } })
db.users.distinct("status", { createdAt: { $gt: ISODate("2024-01-01") } })

// Distinct with logical operators in query
db.users.distinct("role", { $or: [{ active: true }, { verified: true }] })
db.users.distinct("department", { $and: [{ status: "active" }, { role: "employee" }] })

// Distinct with array operators in query
db.users.distinct("city", { tags: { $in: ["premium", "enterprise"] } })

// Distinct with field, query, and options (options passed to driver)
db.users.distinct("email", { status: "active" }, { collation: { locale: "en" } })
db.users.distinct("name", {}, { maxTimeMS: 5000 })

// Distinct with collection access patterns
db["users"].distinct("status")
db['audit-logs'].distinct("action")
db.getCollection("orders").distinct("status", { year: 2024 })
