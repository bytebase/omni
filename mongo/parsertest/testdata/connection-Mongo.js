// Mongo() - Create a new connection to a MongoDB server

// Basic connection
Mongo()
Mongo("localhost")
Mongo("localhost:27017")
Mongo("mongodb://localhost:27017")

// Connection with options
Mongo("mongodb://localhost:27017", { tls: true })
Mongo("mongodb://user:pass@localhost:27017/test?authSource=admin")
Mongo("mongodb://localhost:27017", { ssl: true, retryWrites: true })

// Connection with replica set
Mongo("mongodb://host1:27017,host2:27017,host3:27017/test?replicaSet=myRS")
Mongo("mongodb://localhost:27017", { replicaSet: "myRS" })

// Using connection to get database
Mongo().getDB("test")
Mongo("localhost").getDB("mydb")
Mongo("mongodb://localhost:27017").getDB("admin")

// Chaining multiple methods
Mongo("localhost").getDB("test").getCollection("users")
Mongo().setReadPref("secondary")
Mongo("localhost").setReadConcern("majority")
