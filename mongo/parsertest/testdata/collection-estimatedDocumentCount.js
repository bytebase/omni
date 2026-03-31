// db.collection.estimatedDocumentCount() - Fast estimated count using collection metadata

// Basic estimated count (no filter - uses metadata)
db.users.estimatedDocumentCount()
db.orders.estimatedDocumentCount()
db.products.estimatedDocumentCount()

// Estimated count with options (options passed to driver)
db.users.estimatedDocumentCount({})
db.users.estimatedDocumentCount({ maxTimeMS: 1000 })
db.users.estimatedDocumentCount({ maxTimeMS: 5000 })

// Estimated count with collection access patterns
db["users"].estimatedDocumentCount()
db['audit-logs'].estimatedDocumentCount()
db.getCollection("orders").estimatedDocumentCount()
db.getCollection("large-collection").estimatedDocumentCount({ maxTimeMS: 10000 })
