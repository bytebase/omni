// db.shutdownServer() - Shut down the MongoDB server

// Basic usage
db.shutdownServer()

// With options
db.shutdownServer({ force: true })
db.shutdownServer({ timeoutSecs: 60 })
db.shutdownServer({ force: true, timeoutSecs: 30 })
