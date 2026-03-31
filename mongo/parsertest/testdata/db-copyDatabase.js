// db.copyDatabase() - Deprecated: Copy a database

// Basic usage (deprecated in MongoDB 4.2+)
db.copyDatabase("source", "target")
db.copyDatabase("production", "staging")
db.copyDatabase("olddb", "newdb", "remote.host.com")

// With authentication
db.copyDatabase("source", "target", "remote.host.com", "username", "password")
