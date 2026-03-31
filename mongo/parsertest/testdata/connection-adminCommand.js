// adminCommand() - Run an administrative command on the admin database

// Basic admin commands
Mongo().adminCommand({ listDatabases: 1 })
Mongo("localhost").adminCommand({ serverStatus: 1 })
Mongo("mongodb://localhost:27017").adminCommand({ ping: 1 })

// User management commands
Mongo().adminCommand({ createUser: "myuser", pwd: "password", roles: ["readWrite"] })
Mongo().adminCommand({ usersInfo: 1 })
Mongo().adminCommand({ dropUser: "olduser" })

// Replica set commands
Mongo().adminCommand({ replSetGetStatus: 1 })
Mongo().adminCommand({ replSetGetConfig: 1 })

// Sharding commands
Mongo().adminCommand({ listShards: 1 })
Mongo().adminCommand({ balancerStatus: 1 })

// Database commands
Mongo().adminCommand({ currentOp: 1 })
Mongo("localhost").adminCommand({ killOp: 1, op: 12345 })

// adminCommand from db.getMongo()
db.getMongo().adminCommand({ listDatabases: 1 })
db.getMongo().adminCommand({ serverStatus: 1 })
db.getMongo().adminCommand({ ping: 1 })
