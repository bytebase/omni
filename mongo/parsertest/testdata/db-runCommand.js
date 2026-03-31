// db.runCommand() - Run a database command

// Basic commands
db.runCommand({ ping: 1 })
db.runCommand({ serverStatus: 1 })
db.runCommand({ buildInfo: 1 })

// Collection operations
db.runCommand({ create: "testCollection" })
db.runCommand({ drop: "testCollection" })
db.runCommand({ listCollections: 1 })

// Find command
db.runCommand({
    find: "users",
    filter: { age: { $gt: 21 } },
    limit: 10
})

// Aggregate command
db.runCommand({
    aggregate: "orders",
    pipeline: [
        { $match: { status: "active" } },
        { $group: { _id: "$customerId", total: { $sum: "$amount" } } }
    ],
    cursor: {}
})

// Index commands
db.runCommand({ createIndexes: "users", indexes: [{ key: { email: 1 }, name: "email_1" }] })
db.runCommand({ dropIndexes: "users", index: "email_1" })

// Admin commands
db.runCommand({ collStats: "users" })
db.runCommand({ dbStats: 1, scale: 1024 })
db.runCommand({ validate: "users" })

// User management
db.runCommand({
    createUser: "appUser",
    pwd: "password123",
    roles: [{ role: "readWrite", db: "mydb" }]
})

// Profiling
db.runCommand({ profile: 1, slowms: 100 })
db.runCommand({ profile: 0 })

// Replication commands
db.runCommand({ replSetGetStatus: 1 })
db.runCommand({ isMaster: 1 })
