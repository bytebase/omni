// rs.reconfig() - Reconfigure the replica set

// Basic reconfig with new configuration
rs.reconfig({
    _id: "rs0",
    version: 2,
    members: [
        { _id: 0, host: "mongo1:27017" },
        { _id: 1, host: "mongo2:27017" },
        { _id: 2, host: "mongo3:27017" }
    ]
})

// Reconfig with force option
rs.reconfig({
    _id: "rs0",
    version: 3,
    members: [
        { _id: 0, host: "mongo1:27017" },
        { _id: 1, host: "mongo2:27017" }
    ]
}, { force: true })

// Reconfig changing member priority
rs.reconfig({
    _id: "rs0",
    version: 4,
    members: [
        { _id: 0, host: "mongo1:27017", priority: 2 },
        { _id: 1, host: "mongo2:27017", priority: 1 },
        { _id: 2, host: "mongo3:27017", priority: 0 }
    ]
})

// Reconfig adding hidden member
rs.reconfig({
    _id: "rs0",
    version: 5,
    members: [
        { _id: 0, host: "mongo1:27017" },
        { _id: 1, host: "mongo2:27017" },
        { _id: 2, host: "mongo3:27017", hidden: true, priority: 0 }
    ]
})
