// getReadPrefMode() - Get the read preference mode for the connection

// Basic getReadPrefMode
Mongo().getReadPrefMode()
Mongo("localhost").getReadPrefMode()
Mongo("mongodb://localhost:27017").getReadPrefMode()

// getReadPrefMode from db.getMongo()
db.getMongo().getReadPrefMode()

// getReadPrefMode after setting read preference
Mongo().setReadPref("secondary").getReadPrefMode()
Mongo("localhost").setReadPref("primaryPreferred").getReadPrefMode()
