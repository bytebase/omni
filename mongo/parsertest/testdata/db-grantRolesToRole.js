// db.grantRolesToRole() - Grant roles to a role

// Grant single role
db.grantRolesToRole("appRole", ["read"])

// Grant multiple roles
db.grantRolesToRole("superRole", ["readWrite", "dbAdmin"])

// Grant roles from different databases
db.grantRolesToRole("crossDbRole", [
    { role: "read", db: "reporting" },
    { role: "readWrite", db: "app" }
])

// With write concern
db.grantRolesToRole("criticalRole", ["dbAdmin"], { w: "majority" })
