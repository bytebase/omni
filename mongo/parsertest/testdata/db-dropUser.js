// db.dropUser() - Drop a user from the database

// Basic usage
db.dropUser("tempUser")
db.dropUser("oldUser")
db.dropUser("testAccount")

// With write concern
db.dropUser("removedUser", { w: "majority" })
db.dropUser("deprecatedUser", { w: 1, j: true })
