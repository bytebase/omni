// getReadConcern() - Get the read concern for the connection

// Basic getReadConcern
Mongo().getReadConcern()
Mongo("localhost").getReadConcern()
Mongo("mongodb://localhost:27017").getReadConcern()

// getReadConcern from db.getMongo()
db.getMongo().getReadConcern()

// getReadConcern after setting it
Mongo().setReadConcern("majority").getReadConcern()
