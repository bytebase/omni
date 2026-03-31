// db.grantRolesToUser() - Grant roles to a user

// Grant single role
db.grantRolesToUser("appUser", ["readWrite"])

// Grant multiple roles
db.grantRolesToUser("adminUser", ["dbAdmin", "userAdmin"])

// Grant roles from different databases
db.grantRolesToUser("crossDbUser", [
    { role: "read", db: "logs" },
    { role: "readWrite", db: "app" }
])

// With write concern
db.grantRolesToUser("criticalUser", ["dbOwner"], { w: "majority" })
