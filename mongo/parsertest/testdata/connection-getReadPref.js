// getReadPref() - Get the read preference for the connection

// Basic getReadPref
Mongo().getReadPref()
Mongo("localhost").getReadPref()
Mongo("mongodb://localhost:27017").getReadPref()

// getReadPref from db.getMongo()
db.getMongo().getReadPref()

// getReadPref after setting it
Mongo().setReadPref("secondary").getReadPref()
Mongo("localhost").setReadPref("nearest").getReadPref()
