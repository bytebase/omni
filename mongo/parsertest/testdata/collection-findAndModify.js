// db.collection.findAndModify() - Atomically find and modify a document (unsupported - use findOneAndUpdate/Delete/Replace)

// Basic findAndModify with update
db.users.findAndModify({
    query: { name: "alice" },
    update: { $set: { status: "active" } }
})

// findAndModify with sort
db.queue.findAndModify({
    query: { status: "pending" },
    sort: { priority: -1, createdAt: 1 },
    update: { $set: { status: "processing" } }
})

// findAndModify returning new document
db.users.findAndModify({
    query: { _id: ObjectId("507f1f77bcf86cd799439011") },
    update: { $inc: { score: 10 } },
    new: true
})

// findAndModify with upsert
db.users.findAndModify({
    query: { email: "test@example.com" },
    update: { $set: { name: "Test User", active: true } },
    upsert: true,
    new: true
})

// findAndModify with projection
db.users.findAndModify({
    query: { name: "bob" },
    update: { $set: { lastLogin: Date() } },
    fields: { name: 1, email: 1, lastLogin: 1 }
})

// findAndModify for removal
db.users.findAndModify({
    query: { status: "deleted" },
    remove: true
})

// findAndModify with sort for removal
db.tasks.findAndModify({
    query: { completed: true },
    sort: { completedAt: 1 },
    remove: true
})

// findAndModify with bypassDocumentValidation
db.users.findAndModify({
    query: { _id: 1 },
    update: { $set: { invalidField: "value" } },
    bypassDocumentValidation: true
})

// findAndModify with writeConcern
db.orders.findAndModify({
    query: { orderId: 12345 },
    update: { $set: { status: "confirmed" } },
    writeConcern: { w: "majority", wtimeout: 5000 }
})

// findAndModify with collation
db.products.findAndModify({
    query: { name: "cafe" },
    update: { $set: { available: true } },
    collation: { locale: "fr", strength: 1 }
})

// findAndModify with arrayFilters
db.inventory.findAndModify({
    query: { item: "abc123" },
    update: { $set: { "sizes.$[elem].qty": 0 } },
    arrayFilters: [{ "elem.size": "small" }]
})

// findAndModify with maxTimeMS
db.largeCollection.findAndModify({
    query: { type: "temp" },
    update: { $set: { processed: true } },
    maxTimeMS: 5000
})

// findAndModify with let variables
db.users.findAndModify({
    query: { $expr: { $eq: ["$_id", "$$targetId"] } },
    update: { $set: { found: true } },
    let: { targetId: ObjectId("507f1f77bcf86cd799439011") }
})

// Collection access patterns
db["users"].findAndModify({ query: { x: 1 }, update: { $set: { y: 2 } } })
db.getCollection("users").findAndModify({ query: { a: 1 }, remove: true })
db["task-queue"].findAndModify({ query: { status: "new" }, sort: { createdAt: 1 }, update: { $set: { status: "taken" } } })
