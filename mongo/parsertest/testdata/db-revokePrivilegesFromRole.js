// db.revokePrivilegesFromRole() - Revoke privileges from a role

// Revoke single privilege
db.revokePrivilegesFromRole("customRole", [
    { resource: { db: "mydb", collection: "users" }, actions: ["insert"] }
])

// Revoke multiple privileges
db.revokePrivilegesFromRole("appRole", [
    { resource: { db: "app", collection: "logs" }, actions: ["insert", "update"] },
    { resource: { db: "app", collection: "metrics" }, actions: ["remove"] }
])

// With write concern
db.revokePrivilegesFromRole("sensitiveRole", [
    { resource: { db: "sensitive", collection: "data" }, actions: ["find"] }
], { w: "majority" })
