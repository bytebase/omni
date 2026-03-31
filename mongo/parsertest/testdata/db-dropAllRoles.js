// db.dropAllRoles() - Drop all user-defined roles from the database

// Basic usage
db.dropAllRoles()

// With write concern
db.dropAllRoles({ w: "majority" })
db.dropAllRoles({ w: 1, j: true })
