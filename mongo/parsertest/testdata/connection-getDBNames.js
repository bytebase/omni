// getDBNames() - Get a list of all database names

// Basic getDBNames
Mongo().getDBNames()
Mongo("localhost").getDBNames()
Mongo("mongodb://localhost:27017").getDBNames()

// getDBNames from db.getMongo()
db.getMongo().getDBNames()

// getDBNames after connecting with credentials
Mongo("mongodb://user:pass@localhost:27017").getDBNames()
Mongo("mongodb://admin:password@localhost:27017/?authSource=admin").getDBNames()
