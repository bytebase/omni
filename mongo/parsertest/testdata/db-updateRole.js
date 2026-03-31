// db.updateRole() - Update a role

// Update privileges
db.updateRole("customRole", {
    privileges: [
        { resource: { db: "mydb", collection: "users" }, actions: ["find", "update"] }
    ]
})

// Update inherited roles
db.updateRole("appRole", {
    roles: [
        { role: "readWrite", db: "app" },
        { role: "read", db: "logs" }
    ]
})

// Update both privileges and roles
db.updateRole("fullRole", {
    privileges: [
        { resource: { db: "admin", collection: "system.users" }, actions: ["find"] }
    ],
    roles: ["readWriteAnyDatabase"]
})

// Update authentication restrictions
db.updateRole("restrictedRole", {
    authenticationRestrictions: [
        { clientSource: ["192.168.1.0/24"] }
    ]
})

// With write concern
db.updateRole("criticalRole", {
    privileges: [
        { resource: { db: "critical", collection: "" }, actions: ["find"] }
    ]
}, { w: "majority" })
