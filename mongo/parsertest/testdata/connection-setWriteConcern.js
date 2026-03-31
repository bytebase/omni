// setWriteConcern() - Set the write concern for the connection

// Basic setWriteConcern
Mongo().setWriteConcern({ w: 1 })
Mongo().setWriteConcern({ w: "majority" })
Mongo().setWriteConcern({ w: 2 })

// setWriteConcern with journal
Mongo().setWriteConcern({ w: 1, j: true })
Mongo().setWriteConcern({ w: "majority", j: true })

// setWriteConcern with timeout
Mongo().setWriteConcern({ w: "majority", wtimeout: 5000 })
Mongo().setWriteConcern({ w: 2, j: true, wtimeout: 10000 })

// setWriteConcern with connection string
Mongo("localhost").setWriteConcern({ w: "majority" })
Mongo("mongodb://localhost:27017").setWriteConcern({ w: 1, j: true })

// setWriteConcern from db.getMongo()
db.getMongo().setWriteConcern({ w: "majority" })
db.getMongo().setWriteConcern({ w: 1, j: true, wtimeout: 5000 })

// Chain with other methods
Mongo().setWriteConcern({ w: "majority" }).getDB("test")
