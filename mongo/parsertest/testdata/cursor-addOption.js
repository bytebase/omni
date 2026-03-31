// cursor.addOption() - Add query flags using numeric option values

// Add tailable option (2)
db.cappedCollection.find().addOption(2)

// Add slaveOk option (4)
db.users.find().addOption(4)

// Add oplogReplay option (8)
db.oplog.find().addOption(8)

// Add noCursorTimeout option (16)
db.users.find().addOption(16)

// Add awaitData option (32)
db.cappedCollection.find().addOption(32)

// Add exhaust option (64)
db.users.find().addOption(64)

// Add partial option (128)
db.users.find().addOption(128)

// Combine multiple options
db.cappedCollection.find().addOption(2).addOption(32)
db.users.find().addOption(4).addOption(16)

// Chained with other cursor methods
db.users.find().addOption(16).sort({ name: 1 })
db.users.find().addOption(4).limit(100)

// With query filter
db.logs.find({ level: "error" }).addOption(16)
db.events.find({ type: "notification" }).addOption(2).addOption(32)
