// watch() - Watch for changes on the entire cluster

// Basic watch
Mongo().watch()
Mongo("localhost").watch()
Mongo("mongodb://localhost:27017").watch()

// watch with empty pipeline
Mongo().watch([])

// watch with pipeline stages
Mongo().watch([{ $match: { operationType: "insert" } }])
Mongo("localhost").watch([{ $match: { "fullDocument.status": "active" } }])
Mongo().watch([
    { $match: { operationType: { $in: ["insert", "update", "replace"] } } },
    { $project: { documentKey: 1, fullDocument: 1 } }
])

// watch with options
Mongo().watch([], { fullDocument: "updateLookup" })
Mongo("localhost").watch([], { maxAwaitTimeMS: 1000 })
Mongo().watch([], { batchSize: 100, fullDocument: "updateLookup" })

// watch from db.getMongo()
db.getMongo().watch()
db.getMongo().watch([{ $match: { operationType: "insert" } }])
db.getMongo().watch([], { fullDocument: "updateLookup" })
