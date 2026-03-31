// db.collection.dropIndex() - Drop an index from a collection

// Drop by index name
db.users.dropIndex("email_1")
db.orders.dropIndex("customerId_1_createdAt_-1")

// Drop by index key pattern
db.users.dropIndex({ email: 1 })
db.orders.dropIndex({ customerId: 1, createdAt: -1 })

// Drop text index
db.articles.dropIndex("content_text")

// Drop geospatial index
db.places.dropIndex({ location: "2dsphere" })

// Collection access patterns
db["users"].dropIndex("email_1")
db.getCollection("users").dropIndex({ email: 1 })
