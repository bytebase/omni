// db.watch() - Open a change stream on the database

// Basic usage
db.watch()

// Watch with pipeline
db.watch([
    { $match: { operationType: "insert" } }
])

// Watch specific collections
db.watch([
    { $match: { "ns.coll": { $in: ["users", "orders"] } } }
])

// Watch for specific operation types
db.watch([
    { $match: { operationType: { $in: ["insert", "update", "delete"] } } }
])

// With options
db.watch([], { fullDocument: "updateLookup" })
db.watch([], { maxAwaitTimeMS: 5000 })

// Watch with resume token
db.watch([], {
    resumeAfter: { _data: "826..." }
})

// Watch with start time
db.watch([], {
    startAtOperationTime: Timestamp(1234567890, 1)
})

// Combined pipeline and options
db.watch([
    { $match: { operationType: "update" } },
    { $project: { fullDocument: 1, operationType: 1 } }
], {
    fullDocument: "updateLookup",
    maxAwaitTimeMS: 10000
})
