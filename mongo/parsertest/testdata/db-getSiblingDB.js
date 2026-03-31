// db.getSiblingDB() - Returns another database without modifying db variable

// Basic usage - switch to different databases
db.getSiblingDB("admin")
db.getSiblingDB("test")
db.getSiblingDB("myapp")

// Common use case: access admin database
db.getSiblingDB("admin")

// Access database with special characters in name
db.getSiblingDB("my-app-db")
db.getSiblingDB("app_v2")
