// setReadConcern() - Set the read concern for the connection

// Basic setReadConcern
Mongo().setReadConcern("local")
Mongo().setReadConcern("majority")
Mongo().setReadConcern("linearizable")
Mongo().setReadConcern("available")
Mongo().setReadConcern("snapshot")

// setReadConcern with connection string
Mongo("localhost").setReadConcern("majority")
Mongo("mongodb://localhost:27017").setReadConcern("local")

// setReadConcern from db.getMongo()
db.getMongo().setReadConcern("majority")
db.getMongo().setReadConcern("linearizable")

// Chain with other methods
Mongo().setReadConcern("majority").getDB("test")
Mongo("localhost").setReadConcern("local").setReadPref("secondary")
