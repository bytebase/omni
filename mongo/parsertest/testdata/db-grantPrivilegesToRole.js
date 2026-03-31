// db.grantPrivilegesToRole() - Grant privileges to a role

// Grant single privilege
db.grantPrivilegesToRole("customRole", [
    { resource: { db: "mydb", collection: "users" }, actions: ["find"] }
])

// Grant multiple privileges
db.grantPrivilegesToRole("appRole", [
    { resource: { db: "app", collection: "logs" }, actions: ["find", "insert"] },
    { resource: { db: "app", collection: "metrics" }, actions: ["find"] }
])

// Grant database-wide privilege
db.grantPrivilegesToRole("dbRole", [
    { resource: { db: "reporting", collection: "" }, actions: ["find", "listCollections"] }
])

// With write concern
db.grantPrivilegesToRole("criticalRole", [
    { resource: { db: "critical", collection: "data" }, actions: ["insert", "update"] }
], { w: "majority" })
