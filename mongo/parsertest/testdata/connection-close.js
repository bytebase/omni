// close() - Close the database connection

// Close on Mongo connection
Mongo().close()
Mongo("localhost").close()
Mongo("mongodb://localhost:27017").close()

// Close after operations
Mongo("localhost").getDB("test").close()

// Close from db.getMongo()
db.getMongo().close()
