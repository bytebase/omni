// getReadPrefTagSet() - Get the read preference tag set for the connection

// Basic getReadPrefTagSet
Mongo().getReadPrefTagSet()
Mongo("localhost").getReadPrefTagSet()
Mongo("mongodb://localhost:27017").getReadPrefTagSet()

// getReadPrefTagSet from db.getMongo()
db.getMongo().getReadPrefTagSet()

// getReadPrefTagSet after setting read preference with tags
Mongo().setReadPref("secondary", [{ dc: "east" }]).getReadPrefTagSet()
Mongo("localhost").setReadPref("nearest", [{ region: "us" }]).getReadPrefTagSet()
