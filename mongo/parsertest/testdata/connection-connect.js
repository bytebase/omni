// connect() - Connect to a MongoDB server and return a database object

// Basic connect
connect()
connect("localhost")
connect("localhost:27017")
connect("mongodb://localhost:27017")

// Connect with database name
connect("localhost/mydb")
connect("mongodb://localhost:27017/test")

// Connect with options
connect("mongodb://localhost:27017/test", { readPreference: "secondary" })
connect("mongodb://user:pass@localhost:27017/admin", { ssl: true })

// Connect with authentication
connect("mongodb://user:password@localhost:27017/mydb?authSource=admin")

// Connect and chain methods
connect("localhost").getDB("test")
connect("mongodb://localhost:27017").setReadPref("primaryPreferred")
