// db.collection.hideIndex() - Hide an index from the query planner

// Hide index by name
db.users.hideIndex("email_1")
db.orders.hideIndex("customerId_1_createdAt_-1")
db.products.hideIndex("name_text")

// Hide index by key specification
db.users.hideIndex({ email: 1 })
db.orders.hideIndex({ customerId: 1, createdAt: -1 })
db.products.hideIndex({ category: 1, price: -1 })

// Hide compound index
db.users.hideIndex({ lastName: 1, firstName: 1 })
db.orders.hideIndex({ status: 1, priority: -1, createdAt: -1 })

// Hide text index
db.articles.hideIndex({ content: "text" })

// Hide geospatial index
db.places.hideIndex({ location: "2dsphere" })

// Hide hashed index
db.users.hideIndex({ hashedField: "hashed" })

// Collection access patterns
db["users"].hideIndex("email_1")
db["users"].hideIndex({ email: 1 })
db.getCollection("users").hideIndex("status_1")
db.getCollection("users").hideIndex({ status: 1 })
db["user-profiles"].hideIndex("userId_1")
