// getWriteConcern() - Get the write concern for the connection

// Basic getWriteConcern
Mongo().getWriteConcern()
Mongo("localhost").getWriteConcern()
Mongo("mongodb://localhost:27017").getWriteConcern()

// getWriteConcern from db.getMongo()
db.getMongo().getWriteConcern()

// getWriteConcern after setting it
Mongo().setWriteConcern({ w: "majority" }).getWriteConcern()
Mongo("localhost").setWriteConcern({ w: 1, j: true }).getWriteConcern()
