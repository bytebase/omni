// db.revokeRolesFromUser() - Revoke roles from a user

// Revoke single role
db.revokeRolesFromUser("appUser", ["readWrite"])

// Revoke multiple roles
db.revokeRolesFromUser("adminUser", ["dbAdmin", "userAdmin"])

// Revoke roles from different databases
db.revokeRolesFromUser("crossDbUser", [
    { role: "read", db: "logs" },
    { role: "readWrite", db: "archive" }
])

// With write concern
db.revokeRolesFromUser("sensitiveUser", ["dbOwner"], { w: "majority" })
