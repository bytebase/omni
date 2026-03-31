// getDB() - Get a database reference from a connection

// Basic getDB
Mongo().getDB("test")
Mongo().getDB("mydb")
Mongo().getDB("admin")

// getDB with connection string
Mongo("localhost").getDB("test")
Mongo("mongodb://localhost:27017").getDB("mydb")

// getDB from db.getMongo()
db.getMongo().getDB("anotherdb")
db.getMongo().getDB("test")

// Chain database operations
Mongo().getDB("test").getCollection("users")
Mongo("localhost").getDB("mydb").getCollectionNames()
