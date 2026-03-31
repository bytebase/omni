// rs.add() - Add a member to the replica set

// Add member by hostname
rs.add("mongo4:27017")

// Add member with options document
rs.add({ host: "mongo4:27017" })

// Add member with priority
rs.add({ host: "mongo4:27017", priority: 0 })

// Add member as hidden
rs.add({ host: "mongo4:27017", hidden: true, priority: 0 })

// Add member with votes
rs.add({ host: "mongo4:27017", votes: 0, priority: 0 })

// Add member with build indexes disabled
rs.add({ host: "mongo4:27017", buildIndexes: false, priority: 0 })

// Add delayed member
rs.add({
    host: "mongo4:27017",
    priority: 0,
    hidden: true,
    secondaryDelaySecs: 3600
})
