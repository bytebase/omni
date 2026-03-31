// db.currentOp() - Returns information on current database operations

// Basic usage
db.currentOp()

// Show all operations including idle connections
db.currentOp(true)

// Filter by active operations
db.currentOp({ active: true })

// Filter by operation type
db.currentOp({ op: "query" })
db.currentOp({ op: "insert" })
db.currentOp({ op: "update" })

// Filter by namespace
db.currentOp({ ns: "mydb.users" })

// Filter long-running operations
db.currentOp({ secs_running: { $gt: 5 } })

// Combined filters
db.currentOp({
    active: true,
    secs_running: { $gt: 3 },
    op: { $ne: "none" }
})

// Filter by client
db.currentOp({ client: "192.168.1.100:12345" })

// Show operations waiting for lock
db.currentOp({ waitingForLock: true })
