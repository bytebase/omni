// setReadPref() - Set the read preference for the connection

// Basic setReadPref with mode only
Mongo().setReadPref("primary")
Mongo().setReadPref("secondary")
Mongo().setReadPref("primaryPreferred")
Mongo().setReadPref("secondaryPreferred")
Mongo().setReadPref("nearest")

// setReadPref with connection string
Mongo("localhost").setReadPref("secondary")
Mongo("mongodb://localhost:27017").setReadPref("nearest")

// setReadPref with tag set
Mongo().setReadPref("secondary", [{ dc: "east" }])
Mongo().setReadPref("nearest", [{ region: "us-east-1" }])
Mongo("localhost").setReadPref("secondaryPreferred", [{ dc: "east" }, { dc: "west" }])

// setReadPref with options
Mongo().setReadPref("secondary", [{ dc: "east" }], { maxStalenessSeconds: 120 })

// setReadPref from db.getMongo()
db.getMongo().setReadPref("secondary")
db.getMongo().setReadPref("nearest", [{ region: "us" }])
