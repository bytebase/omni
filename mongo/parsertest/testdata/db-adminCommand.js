// db.adminCommand() - Run an administrative command against the admin database

// Server status
db.adminCommand({ serverStatus: 1 })

// Current operations
db.adminCommand({ currentOp: 1 })
db.adminCommand({ currentOp: 1, active: true })

// Kill operations
db.adminCommand({ killOp: 1, op: 12345 })

// Shutdown
db.adminCommand({ shutdown: 1 })
db.adminCommand({ shutdown: 1, force: true })

// Log management
db.adminCommand({ getLog: "global" })
db.adminCommand({ getLog: "startupWarnings" })

// Replica set commands
db.adminCommand({ replSetGetStatus: 1 })
db.adminCommand({ replSetGetConfig: 1 })
db.adminCommand({ replSetStepDown: 60 })

// Sharding commands
db.adminCommand({ listShards: 1 })
db.adminCommand({ enableSharding: "mydb" })
db.adminCommand({ shardCollection: "mydb.users", key: { _id: "hashed" } })

// User management
db.adminCommand({ listDatabases: 1 })
db.adminCommand({ listDatabases: 1, nameOnly: true })

// Feature compatibility
db.adminCommand({ getParameter: 1, featureCompatibilityVersion: 1 })
db.adminCommand({ setFeatureCompatibilityVersion: "6.0" })

// Connection pool
db.adminCommand({ connPoolStats: 1 })

// fsync
db.adminCommand({ fsync: 1 })
db.adminCommand({ fsync: 1, lock: true })
