// db.dropAllUsers() - Drop all users from the database

// Basic usage
db.dropAllUsers()

// With write concern
db.dropAllUsers({ w: "majority" })
db.dropAllUsers({ w: 1, j: true })
