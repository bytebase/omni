// db.setLogLevel() - Set the log verbosity level

// Set global log level
db.setLogLevel(1)
db.setLogLevel(2)
db.setLogLevel(0)

// Set component-specific log level
db.setLogLevel(1, "query")
db.setLogLevel(2, "replication")
db.setLogLevel(1, "storage")
db.setLogLevel(3, "network")
db.setLogLevel(2, "accessControl")
