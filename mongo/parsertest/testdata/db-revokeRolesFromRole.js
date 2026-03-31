// db.revokeRolesFromRole() - Revoke roles from a role

// Revoke single role
db.revokeRolesFromRole("appRole", ["read"])

// Revoke multiple roles
db.revokeRolesFromRole("superRole", ["dbAdmin", "userAdmin"])

// Revoke roles from different databases
db.revokeRolesFromRole("crossDbRole", [
    { role: "read", db: "reporting" },
    { role: "readWrite", db: "logs" }
])

// With write concern
db.revokeRolesFromRole("criticalRole", ["dbOwner"], { w: "majority" })
