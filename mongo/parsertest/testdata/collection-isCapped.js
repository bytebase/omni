// db.collection.isCapped() - Check if a collection is a capped collection

// Basic usage
db.users.isCapped()
db.orders.isCapped()
db.logs.isCapped()
db.events.isCapped()

// Commonly used for checking log collections
db.systemLogs.isCapped()
db.auditTrail.isCapped()

// Collection access patterns
db["users"].isCapped()
db.getCollection("users").isCapped()
db["capped-logs"].isCapped()
db.getCollection("system.profile").isCapped()
