// cursor.readPref() - Specify read preference for replica set queries

// Primary read preference (default)
db.users.find().readPref("primary")

// Primary preferred
db.users.find().readPref("primaryPreferred")

// Secondary read preference
db.users.find().readPref("secondary")

// Secondary preferred
db.users.find().readPref("secondaryPreferred")

// Nearest read preference
db.users.find().readPref("nearest")

// With tag sets
db.users.find().readPref("secondary", [{ region: "east" }])
db.users.find().readPref("nearest", [{ dc: "us-east-1" }, { dc: "us-west-2" }])

// With maxStalenessSeconds
db.users.find().readPref("secondaryPreferred", [], { maxStalenessSeconds: 120 })

// With query filter
db.users.find({ status: "active" }).readPref("secondary")
db.analytics.find({ date: { $gte: ISODate("2024-01-01") } }).readPref("secondaryPreferred")

// Chained with other cursor methods
db.users.find().readPref("secondary").sort({ name: 1 })
db.users.find().readPref("nearest").limit(10)
db.reports.find().readPref("secondaryPreferred").maxTimeMS(30000)

// Read from specific data center
db.users.find().readPref("nearest", [{ datacenter: "NYC" }])
