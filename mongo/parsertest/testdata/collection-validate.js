// db.collection.validate() - Validate a collection's data and indexes

// Basic validation
db.users.validate()
db.orders.validate()

// Full validation (more thorough but slower)
db.users.validate({ full: true })

// Quick validation (faster, less thorough)
db.products.validate({ full: false })

// Validation with repair option
db.logs.validate({ repair: true })

// Validation with metadata check only
db.sessions.validate({ metadata: true })

// Combined options
db.largeCollection.validate({ full: true, background: true })

// Collection access patterns
db["users"].validate()
db.getCollection("users").validate()
db["important-data"].validate({ full: true })
db.getCollection("production.orders").validate({ repair: true })
