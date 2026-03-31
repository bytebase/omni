// db.collection.watch() - Open a change stream cursor (unsupported - requires persistent connection)

// Basic watch
db.users.watch()
db.orders.watch()

// Watch with empty pipeline
db.users.watch([])

// Watch with pipeline stages
db.users.watch([{ $match: { operationType: "insert" } }])
db.orders.watch([{ $match: { "fullDocument.status": "shipped" } }])

// Watch with multiple pipeline stages
db.products.watch([
    { $match: { operationType: { $in: ["insert", "update", "replace"] } } },
    { $project: { fullDocument: 1, operationType: 1 } }
])

// Watch with options
db.users.watch([], { fullDocument: "updateLookup" })
db.orders.watch([], { fullDocument: "whenAvailable" })
db.products.watch([], { fullDocumentBeforeChange: "whenAvailable" })

// Watch with resumeAfter
db.users.watch([], { resumeAfter: { _data: "826..." } })

// Watch with startAfter
db.orders.watch([], { startAfter: { _data: "826..." } })

// Watch with startAtOperationTime
db.events.watch([], { startAtOperationTime: Timestamp(1609459200, 1) })

// Watch with batchSize
db.users.watch([], { batchSize: 100 })

// Watch with maxAwaitTimeMS
db.logs.watch([], { maxAwaitTimeMS: 5000 })

// Watch with collation
db.products.watch([{ $match: { "fullDocument.category": "Electronics" } }], { collation: { locale: "en", strength: 2 } })

// Watch with showExpandedEvents
db.users.watch([], { showExpandedEvents: true })

// Watch combined options
db.orders.watch([
    { $match: { operationType: "update" } }
], {
    fullDocument: "updateLookup",
    batchSize: 50,
    maxAwaitTimeMS: 10000
})

// Collection access patterns
db["users"].watch()
db["users"].watch([{ $match: { operationType: "delete" } }])
db.getCollection("users").watch()
db.getCollection("events").watch([], { fullDocument: "updateLookup" })
db["change-events"].watch([])
