// rs.initiate() - Initialize replica set

// No config (auto-initiate single node)
rs.initiate()

// With configuration
rs.initiate({
    _id: "myReplicaSet",
    members: [
        { _id: 0, host: "mongo1:27017" },
        { _id: 1, host: "mongo2:27017" },
        { _id: 2, host: "mongo3:27017" }
    ]
})

// With configuration including priority and votes
rs.initiate({
    _id: "rs0",
    members: [
        { _id: 0, host: "primary:27017", priority: 2 },
        { _id: 1, host: "secondary1:27017", priority: 1 },
        { _id: 2, host: "secondary2:27017", priority: 1, votes: 1 }
    ]
})

// With arbiter configuration
rs.initiate({
    _id: "rs0",
    members: [
        { _id: 0, host: "mongo1:27017" },
        { _id: 1, host: "mongo2:27017" },
        { _id: 2, host: "arbiter:27017", arbiterOnly: true }
    ]
})
