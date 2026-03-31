// startSession() - Start a new client session

// Basic startSession
Mongo().startSession()
Mongo("localhost").startSession()
Mongo("mongodb://localhost:27017").startSession()

// startSession with options
Mongo().startSession({})
Mongo().startSession({ causalConsistency: true })
Mongo().startSession({ causalConsistency: false })

// startSession with retryable writes
Mongo().startSession({ retryWrites: true })
Mongo("localhost").startSession({ retryWrites: true, causalConsistency: true })

// startSession with snapshot reads
Mongo().startSession({ snapshot: true })

// startSession from db.getMongo()
db.getMongo().startSession()
db.getMongo().startSession({ causalConsistency: true })
db.getMongo().startSession({ retryWrites: true })
