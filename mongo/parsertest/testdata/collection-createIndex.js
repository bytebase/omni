// db.collection.createIndex() - Create an index on a collection

// Basic single field index
db.users.createIndex({ email: 1 })
db.users.createIndex({ age: -1 })

// Compound index
db.users.createIndex({ lastName: 1, firstName: 1 })
db.orders.createIndex({ customerId: 1, createdAt: -1 })

// Text index
db.articles.createIndex({ content: "text" })
db.posts.createIndex({ title: "text", body: "text" })

// Unique index
db.users.createIndex({ email: 1 }, { unique: true })

// Sparse index
db.users.createIndex({ phone: 1 }, { sparse: true })

// TTL index
db.sessions.createIndex({ createdAt: 1 }, { expireAfterSeconds: 3600 })

// Partial index
db.orders.createIndex({ status: 1 }, { partialFilterExpression: { status: "pending" } })

// Index with options
db.users.createIndex({ username: 1 }, { unique: true, background: true })
db.products.createIndex({ name: 1, category: 1 }, { name: "name_category_idx" })

// Hashed index
db.users.createIndex({ hashedField: "hashed" })

// Geospatial index
db.places.createIndex({ location: "2dsphere" })
db.stores.createIndex({ coordinates: "2d" })

// Collection access patterns
db["users"].createIndex({ email: 1 })
db.getCollection("users").createIndex({ email: 1 })
