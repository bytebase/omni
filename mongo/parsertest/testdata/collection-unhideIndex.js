// db.collection.unhideIndex() - Unhide a hidden index to make it visible to the query planner

// Unhide index by name
db.users.unhideIndex("email_1")
db.orders.unhideIndex("customerId_1_createdAt_-1")
db.products.unhideIndex("name_text")

// Unhide index by key specification
db.users.unhideIndex({ email: 1 })
db.orders.unhideIndex({ customerId: 1, createdAt: -1 })
db.products.unhideIndex({ category: 1, price: -1 })

// Unhide compound index
db.users.unhideIndex({ lastName: 1, firstName: 1 })
db.orders.unhideIndex({ status: 1, priority: -1, createdAt: -1 })

// Unhide text index
db.articles.unhideIndex({ content: "text" })

// Unhide geospatial index
db.places.unhideIndex({ location: "2dsphere" })

// Unhide hashed index
db.users.unhideIndex({ hashedField: "hashed" })

// Collection access patterns
db["users"].unhideIndex("email_1")
db["users"].unhideIndex({ email: 1 })
db.getCollection("users").unhideIndex("status_1")
db.getCollection("users").unhideIndex({ status: 1 })
db["user-profiles"].unhideIndex("userId_1")
