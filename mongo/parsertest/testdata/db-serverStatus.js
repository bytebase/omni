// db.serverStatus() - Returns an overview of the database server state

// Basic server status
db.serverStatus()

// Server status with specific sections
db.serverStatus({ repl: 0, metrics: 0 })
db.serverStatus({ locks: 1 })
db.serverStatus({ opcounters: 1, connections: 1 })

// Exclude specific sections
db.serverStatus({ asserts: 0, connections: 0 })

// Server status with latchAnalysis
db.serverStatus({ latchAnalysis: 1 })

// Server status with mirroredReads
db.serverStatus({ mirroredReads: 1 })
