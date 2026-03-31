// db.aggregate() - Run aggregation pipeline against the database

// Basic aggregation
db.aggregate([
    { $currentOp: {} }
])

// List all collections via aggregation
db.aggregate([
    { $listLocalSessions: {} }
])

// Admin database operations
db.aggregate([
    { $currentOp: { allUsers: true, idleSessions: true } }
])

// With options
db.aggregate([
    { $listLocalSessions: { allUsers: true } }
], { allowDiskUse: true })
